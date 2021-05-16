/*
Copyright 2021 Maximilian Geberl.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"encoding/json"
	"reflect"
	"time"

	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	kubeteachv1alpha1 "github.com/dergeberl/kubeteach/api/v1alpha1"
)

// ExerciseSetReconciler reconciles a ExerciseSet object
type ExerciseSetReconciler struct {
	client.Client
	Log         logr.Logger
	Scheme      *runtime.Scheme
	RequeueTime time.Duration
}

//+kubebuilder:rbac:groups=kubeteach.geberl.io,resources=exercisesets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=kubeteach.geberl.io,resources=exercisesets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=kubeteach.geberl.io,resources=exercisesets/finalizers,verbs=update

// Reconcile handles reconcile of an ExersiceSet
func (r *ExerciseSetReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = r.Log.WithValues("exerciseset", req.NamespacedName)

	var exerciseSet kubeteachv1alpha1.ExerciseSet
	err := r.Client.Get(ctx, req.NamespacedName, &exerciseSet)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	var newExerciseSetStatus kubeteachv1alpha1.ExerciseSetStatus

	for _, taskDefinition := range exerciseSet.Spec.TaskDefinitions {
		var taskDefinitionObject kubeteachv1alpha1.TaskDefinition
		err = r.Client.Get(ctx, client.ObjectKey{Name: taskDefinition.Name, Namespace: req.Namespace}, &taskDefinitionObject)
		if err != nil {
			if client.IgnoreNotFound(err) != nil {
				return ctrl.Result{}, err
			}
			taskDefinitionObject = kubeteachv1alpha1.TaskDefinition{
				ObjectMeta: metav1.ObjectMeta{
					Name:      taskDefinition.Name,
					Namespace: req.Namespace,
					OwnerReferences: []metav1.OwnerReference{{
						APIVersion: exerciseSet.APIVersion,
						Kind:       exerciseSet.Kind,
						Name:       exerciseSet.Name,
						UID:        exerciseSet.UID,
					}},
				},
				Spec: taskDefinition.TaskDefinitionSpec,
			}
			err = r.Client.Create(ctx, &taskDefinitionObject)
			if err != nil {
				return ctrl.Result{}, err
			}
		}

		// update TaskDefinition if needed
		if !reflect.DeepEqual(taskDefinitionObject.Spec, taskDefinition.TaskDefinitionSpec) {
			taskDefinitionObject.Spec = taskDefinition.TaskDefinitionSpec
			err = r.Client.Update(ctx, &taskDefinitionObject)
			if err != nil {
				return ctrl.Result{}, err
			}
		}

		if !reflect.DeepEqual(taskDefinitionObject.OwnerReferences, []metav1.OwnerReference{{
			APIVersion: exerciseSet.APIVersion,
			Kind:       exerciseSet.Kind,
			Name:       exerciseSet.Name,
			UID:        exerciseSet.UID,
		}}) {
			taskDefinitionObject.OwnerReferences = []metav1.OwnerReference{{
				APIVersion: exerciseSet.APIVersion,
				Kind:       exerciseSet.Kind,
				Name:       exerciseSet.Name,
				UID:        exerciseSet.UID,
			}}
			err = r.Client.Update(ctx, &taskDefinitionObject)
			if err != nil {
				return ctrl.Result{}, err
			}
		}

		// count status infos
		newExerciseSetStatus.NumberOfTasks++
		if taskDefinitionObject.Status.State != nil {
			switch *taskDefinitionObject.Status.State {
			case stateActive:
				newExerciseSetStatus.NumberOfActiveTasks++
			case statePending:
				newExerciseSetStatus.NumberOfPendingTasks++
			case stateSuccessful:
				newExerciseSetStatus.NumberOfSuccessfulTasks++
			}
		} else {
			newExerciseSetStatus.NumberOfUnknownTasks++
		}

		newExerciseSetStatus.PointsTotal += taskDefinition.TaskDefinitionSpec.Points
		if taskDefinition.TaskDefinitionSpec.Points == 0 {
			newExerciseSetStatus.NumberOfTasksWithoutPoints++
		}
		if taskDefinitionObject.Status.State != nil &&
			*taskDefinitionObject.Status.State == stateSuccessful {
			newExerciseSetStatus.PointsAchieved += taskDefinition.TaskDefinitionSpec.Points
		}
	}

	// update status if needed
	if !reflect.DeepEqual(exerciseSet.Status, newExerciseSetStatus) {
		var statusJason []byte
		statusJason, err = json.Marshal(newExerciseSetStatus)
		if err != nil {
			return ctrl.Result{}, err
		}
		patch := []byte(`{"status":` + string(statusJason) + `}`)
		err = r.Client.Status().Patch(ctx, &exerciseSet, client.RawPatch(types.MergePatchType, patch))
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{RequeueAfter: r.RequeueTime}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ExerciseSetReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kubeteachv1alpha1.ExerciseSet{}).
		Complete(r)
}
