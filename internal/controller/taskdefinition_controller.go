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

// Package controllers is a package of kubeteach and used for reconcile logic of kubernetes CRDs
package controller

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	teachv1alpha1 "github.com/dergeberl/kubeteach/api/v1alpha1"
	"github.com/dergeberl/kubeteach/internal/controller/condition"
)

// const for state field
const (
	StateActive     = "active"
	StateSuccessful = "successful"
	StatePending    = "pending"
)

// TaskDefinitionReconciler reconciles a TaskDefinition object
type TaskDefinitionReconciler struct {
	client.Client
	Scheme      *runtime.Scheme
	Recorder    record.EventRecorder
	RequeueTime time.Duration
}

// +kubebuilder:rbac:groups=kubeteach.geberl.io,resources=taskdefinitions,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=kubeteach.geberl.io,resources=taskdefinitions/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=kubeteach.geberl.io,resources=tasks,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=kubeteach.geberl.io,resources=tasks/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=kubeteach.geberl.io,resources=tasks/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch

// Reconcile handles all about taskdefinitions and tasks
func (r *TaskDefinitionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// get current taskDefinition
	taskDefinition := teachv1alpha1.TaskDefinition{}
	err := r.Client.Get(ctx, req.NamespacedName, &taskDefinition)
	if err != nil {
		// ignore taskdefinitons that dose not exists
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// skip delete objects
	if !taskDefinition.ObjectMeta.DeletionTimestamp.IsZero() {
		return ctrl.Result{}, nil
	}

	// set status if empty to StatePending
	if taskDefinition.Status.State == nil {
		err = r.setState(ctx, StatePending, &taskDefinition)
		if err != nil {
			return ctrl.Result{}, err
		}
		err = r.notifyExerciseSet(ctx, taskDefinition)
		if err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	// skip if status is already StateSuccessful
	if *taskDefinition.Status.State == StateSuccessful {
		return ctrl.Result{}, nil
	}

	// create or update task for taskdefinition
	task, err := r.createOrUpdateTask(ctx, &taskDefinition)
	if err != nil {
		return ctrl.Result{}, err
	}

	// check pending state
	if *taskDefinition.Status.State == StatePending {
		return r.checkPending(ctx, req, &taskDefinition, &task)
	}

	// run ConditionChecks checks
	ConditionChecks := condition.Checks{
		Client: r.Client,
	}
	status, err := ConditionChecks.ApplyChecks(ctx, taskDefinition.Spec.TaskConditions)
	if err != nil {
		r.Recorder.Event(&taskDefinition, "Warning", "Error", fmt.Sprintf("Conditions apply fail with error: %v", err))
		return ctrl.Result{}, err
	}

	// check status
	if status {
		err = r.setState(ctx, StateSuccessful, &taskDefinition, &task)
		if err != nil {
			return ctrl.Result{}, err
		}
		err = r.notifyExerciseSet(ctx, taskDefinition)
		if err != nil {
			return ctrl.Result{}, err
		}
		r.Recorder.Event(&task, "Normal", "Successful", "Task is successfully completed")
		return ctrl.Result{}, nil
	}
	return ctrl.Result{RequeueAfter: r.RequeueTime}, nil
}

// checkPending check if task is still in pending or the required task is already done
func (r *TaskDefinitionReconciler) checkPending(
	ctx context.Context,
	req ctrl.Request,
	taskDefinition *teachv1alpha1.TaskDefinition,
	task *teachv1alpha1.Task,
) (ctrl.Result, error) {
	if taskDefinition.Spec.RequiredTaskName != nil {
		// get pre required taskdefiniton
		reqTask := teachv1alpha1.TaskDefinition{}
		err := r.Client.Get(ctx, client.ObjectKey{
			Name:      *taskDefinition.Spec.RequiredTaskName,
			Namespace: req.Namespace},
			&reqTask)
		if err != nil {
			if errors.IsNotFound(err) {
				return ctrl.Result{RequeueAfter: r.RequeueTime}, nil
			}
			return ctrl.Result{}, err
		}

		// set state to active if pre required task is successful
		if reqTask.Status.State != nil && *reqTask.Status.State == StateSuccessful {
			r.Recorder.Event(task, "Normal", "Active", "Pre required task is successful, task is now active")
			err = r.setState(ctx, StateActive, taskDefinition, task)
			if err != nil {
				return ctrl.Result{}, err
			}
			err = r.notifyExerciseSet(ctx, *taskDefinition)
			if err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{Requeue: true}, nil
		}

		// requeue after RequeueTime to check again
		return ctrl.Result{RequeueAfter: r.RequeueTime}, nil
	}

	// set state to active no pre required task is defined
	r.Recorder.Event(task, "Normal", "Active", "Task has no pre required task, task is now active")
	err := r.setState(ctx, StateActive, taskDefinition, task)
	if err != nil {
		return ctrl.Result{}, err
	}
	err = r.notifyExerciseSet(ctx, *taskDefinition)
	if err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{Requeue: true}, nil
}

// createOrUpdateTask creates task fot taskdefinition if needed and update task if something changed.
func (r *TaskDefinitionReconciler) createOrUpdateTask(
	ctx context.Context,
	taskDefinition *teachv1alpha1.TaskDefinition,
) (teachv1alpha1.Task, error) {
	// search if task already exists
	taskList := teachv1alpha1.TaskList{}
	err := r.Client.List(ctx, &taskList)
	if err != nil {
		return teachv1alpha1.Task{}, err
	}
	var task *teachv1alpha1.Task
	for i, taskTtem := range taskList.Items {
		if taskTtem.OwnerReferences[0].UID == taskDefinition.UID {
			// found task
			task = &taskList.Items[i]
			break
		}
	}

	// create task if not found
	if task == nil {
		task = &teachv1alpha1.Task{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Task",
				APIVersion: "kubeteach.geberl.io/v1alpha1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      taskDefinition.ObjectMeta.Name,
				Namespace: taskDefinition.ObjectMeta.Namespace,
				OwnerReferences: []metav1.OwnerReference{{
					APIVersion: taskDefinition.APIVersion,
					Kind:       taskDefinition.Kind,
					Name:       taskDefinition.Name,
					UID:        taskDefinition.UID,
				}},
			},
			Spec:   taskDefinition.Spec.TaskSpec,
			Status: teachv1alpha1.TaskStatus{State: taskDefinition.Status.State},
		}
		err = r.Client.Create(ctx, task)
		if err != nil {
			return teachv1alpha1.Task{}, err
		}
		r.Recorder.Event(taskDefinition, "Normal", "Created", "Task created")

		return *task, nil
	}

	// sync spec if task.Spec != taskDefinition.Spec.TaskSpec.
	if !reflect.DeepEqual(task.Spec, taskDefinition.Spec.TaskSpec) {
		task.Spec = taskDefinition.Spec.TaskSpec
		err := r.Update(ctx, task)
		if err != nil {
			return teachv1alpha1.Task{}, err
		}
		r.Recorder.Event(taskDefinition, "Normal", "Update", "Task updated")
	}

	// sync status if status.state is not the same
	if taskDefinition.Status.State != task.Status.State {
		if err := r.setState(ctx, *taskDefinition.Status.State, task); err != nil {
			return teachv1alpha1.Task{}, err
		}
		r.Recorder.Event(taskDefinition, "Normal", "Update", "Task Status updated")
	}
	return *task, nil
}

// setState stets a the status.state field in all objects that are given
func (r *TaskDefinitionReconciler) setState(
	ctx context.Context,
	state string,
	objects ...client.Object,
) error {
	patch := []byte(`{"status":{"state":"` + state + `"}}`)
	for _, object := range objects {
		err := r.Status().Patch(ctx, object, client.RawPatch(types.MergePatchType, patch))
		if err != nil {
			return err
		}
	}
	return nil
}

// notifyExerciseSet chanes an annotation of the exerciseSet to trigger an reconcile
func (r *TaskDefinitionReconciler) notifyExerciseSet(
	ctx context.Context,
	taskDefinition teachv1alpha1.TaskDefinition,
) error {
	for _, owner := range taskDefinition.OwnerReferences {
		if owner.Kind == "ExerciseSet" &&
			owner.Name != "" {
			var exerciseSet teachv1alpha1.ExerciseSet
			err := r.Client.Get(ctx, client.ObjectKey{Name: owner.Name, Namespace: taskDefinition.Namespace}, &exerciseSet)
			if err != nil {
				return err
			}
			patch := []byte(`{"metadata": { "annotations": {"geberl.io/kubeteach-trigger": "` + fmt.Sprint(time.Now().UnixNano()) + `"}}}`)
			err = r.Client.Patch(ctx, &exerciseSet, client.RawPatch(types.MergePatchType, patch))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *TaskDefinitionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&teachv1alpha1.TaskDefinition{}).
		Complete(r)
}
