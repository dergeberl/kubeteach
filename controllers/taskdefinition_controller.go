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
package controllers

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	teachv1alpha1 "github.com/dergeberl/kubeteach/api/v1alpha1"
	"github.com/dergeberl/kubeteach/controllers/condition"
)

const (
	stateActive     = "active"
	stateSuccessful = "successful"
	statePending    = "pending"
)

// TaskDefinitionReconciler reconciles a TaskDefinition object
type TaskDefinitionReconciler struct {
	client.Client
	Log         logr.Logger
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
	_ = r.Log.WithValues("taskdefinition", req.NamespacedName)

	_ = r.Log.WithValues("taskdefinition", req.NamespacedName)

	// get current taskDefinition
	taskDefinition := teachv1alpha1.TaskDefinition{}
	err := r.Client.Get(ctx, req.NamespacedName, &taskDefinition)
	if err != nil {
		// ignore taskdefinitons that dose not exists
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// skip delete objects
	if !taskDefinition.ObjectMeta.DeletionTimestamp.IsZero() {
		return ctrl.Result{}, nil
	}

	// set status if empty to statePending
	if taskDefinition.Status.State == nil {
		err = r.setState(ctx, statePending, &taskDefinition)
		if err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	// skip if status is already stateSuccessful
	if *taskDefinition.Status.State == stateSuccessful {
		return ctrl.Result{}, nil
	}

	// create or update task for taskdefinition
	task, err := r.createOrUpdateTask(ctx, &taskDefinition)
	if err != nil {
		return ctrl.Result{}, err
	}

	// check pending state
	if *taskDefinition.Status.State == statePending {
		return r.checkPending(ctx, req, &taskDefinition, &task)
	}

	// run ConditionChecks checks
	ConditionChecks := condition.Checks{
		Client: r.Client,
		Log:    r.Log,
	}
	status, err := ConditionChecks.ApplyChecks(ctx, taskDefinition.Spec.TaskConditions)
	if err != nil {
		r.Recorder.Event(&taskDefinition, "Warning", "Error", fmt.Sprintf("Conditions apply fail with error: %v", err))
		return ctrl.Result{}, err
	}

	// check status
	if status {
		err = r.setState(ctx, stateSuccessful, &taskDefinition, &task)
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
		if reqTask.Status.State != nil && *reqTask.Status.State == stateSuccessful {
			r.Recorder.Event(task, "Normal", "Active", "Pre required task is successful, task is now active")
			err = r.setState(ctx, stateActive, taskDefinition, task)
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
	err := r.setState(ctx, stateActive, taskDefinition, task)
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

// SetupWithManager sets up the controller with the Manager.
func (r *TaskDefinitionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&teachv1alpha1.TaskDefinition{}).
		Complete(r)
}
