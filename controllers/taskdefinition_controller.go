/*
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
	"fmt"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	teachv1alpha1 "kubeteach/api/v1alpha1"
	"kubeteach/controllers/condition"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"time"
)

const (
	stateActive     = "active"
	stateSuccessful = "successful"
	statePending    = "pending"
	stateError      = "error"
	//requeueTime     = time.Duration(5) * time.Second
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

// Reconcile handles all about taskdefinitions ans tasks
func (r *TaskDefinitionReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	_ = r.Log.WithValues("taskdefinition", req.NamespacedName)

	// get taskDefinition
	taskDefinition := teachv1alpha1.TaskDefinition{}
	err := r.Client.Get(ctx, req.NamespacedName, &taskDefinition)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	//skip delete objects
	if !taskDefinition.ObjectMeta.DeletionTimestamp.IsZero() {
		return ctrl.Result{}, nil
	}

	// set status if empty
	if taskDefinition.Status.State == nil {
		err = r.setState(ctx, statePending, taskDefinition.DeepCopyObject())
		if err != nil {
			return reconcile.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	// set state back to pending when it was error to try a reconcile
	if *taskDefinition.Status.State == stateError {
		*taskDefinition.Status.State = stateActive
	}

	//skip successful objects
	if *taskDefinition.Status.State == stateSuccessful {
		return ctrl.Result{}, nil
	}

	//create or update Task
	task, err := r.createOrUpdateTask(ctx, &taskDefinition)
	if err != nil {
		return ctrl.Result{}, err
	}
	//check pending state
	if *taskDefinition.Status.State == statePending {
		return r.checkPending(ctx, req, taskDefinition, task)
	}

	//run checks
	ConditionChecks := condition.Checks{
		Client: r.Client,
		Log:    r.Log,
	}
	status, err := ConditionChecks.ApplyChecks(ctx, taskDefinition.Spec.TaskConditions)
	if err != nil {
		r.Recorder.Event(&taskDefinition, "Warning", "Error", fmt.Sprintf("Conditions apply fail with error: %v", err))
		return r.errorRequeueAfter(ctx, err, taskDefinition, &task)
	}

	// check status
	if status {
		err = r.setState(ctx, stateSuccessful, taskDefinition.DeepCopyObject(), task.DeepCopyObject())
		if err != nil {
			return reconcile.Result{}, err
		}
		r.Recorder.Event(&task, "Normal", "Successful", "Task is successfully completed")
		return ctrl.Result{}, nil
	}
	return ctrl.Result{RequeueAfter: r.RequeueTime}, nil
}

func (r *TaskDefinitionReconciler) checkPending(ctx context.Context, req ctrl.Request, taskDefinition teachv1alpha1.TaskDefinition, task teachv1alpha1.Task) (ctrl.Result, error) {
	if taskDefinition.Spec.RequiredTaskName != nil {
		reqTask := teachv1alpha1.TaskDefinition{}
		err := r.Client.Get(ctx, client.ObjectKey{Name: *taskDefinition.Spec.RequiredTaskName, Namespace: req.Namespace}, &reqTask)
		if err != nil {
			if errors.IsNotFound(err) {
				return ctrl.Result{}, nil
			}
			return ctrl.Result{}, err
		}
		if *reqTask.Status.State == stateSuccessful {
			r.Recorder.Event(&task, "Normal", "Active", "Pre required task is now successful, applying pre required objects")
			err = r.setState(ctx, stateActive, taskDefinition.DeepCopyObject(), task.DeepCopyObject())
			if err != nil {
				return reconcile.Result{}, err
			}
			return ctrl.Result{Requeue: true}, nil
		}
		return ctrl.Result{RequeueAfter: r.RequeueTime}, nil
	}
	r.Recorder.Event(&task, "Normal", "Active", "Task has no pre required task, applying pre required objects")
	err := r.setState(ctx, stateActive, taskDefinition.DeepCopyObject(), task.DeepCopyObject())
	if err != nil {
		return reconcile.Result{}, err
	}
	return ctrl.Result{Requeue: true}, nil

}

func (r *TaskDefinitionReconciler) createOrUpdateTask(ctx context.Context, taskDefinition *teachv1alpha1.TaskDefinition) (teachv1alpha1.Task, error) {

	taskList := teachv1alpha1.TaskList{}
	err := r.Client.List(ctx, &taskList)
	if err != nil {
		return teachv1alpha1.Task{}, err
	}
	var task *teachv1alpha1.Task
	for _, taskTtem := range taskList.Items {
		if taskTtem.OwnerReferences[0].UID == taskDefinition.UID {
			//found task
			task = &taskTtem
			break
		}
	}
	//TODO check if something is still in deletion
	//create task if not found
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
		err := r.Client.Create(ctx, task)
		if err != nil {
			return teachv1alpha1.Task{}, err
		}
		r.Recorder.Event(taskDefinition, "Normal", "Created", "Task created")

		return *task, nil
	}

	//sync spec
	if !reflect.DeepEqual(task.Spec, taskDefinition.Spec.TaskSpec) {
		task.Spec = taskDefinition.Spec.TaskSpec
		err := r.Update(ctx, task)
		if err != nil {
			return teachv1alpha1.Task{}, err
		}
		r.Recorder.Event(taskDefinition, "Normal", "Update", "Task updated")
	}

	//sync status
	if taskDefinition.Status.State != task.Status.State {
		if err := r.setState(ctx, *taskDefinition.Status.State, task.DeepCopyObject()); err != nil {
			return teachv1alpha1.Task{}, err
		}
		r.Recorder.Event(taskDefinition, "Normal", "Update", "Task Status updated")

	}
	return *task, nil
}

func (r *TaskDefinitionReconciler) setState(ctx context.Context, state string, objects ...runtime.Object) error {
	patch := []byte(`{"status":{"state":"` + state + `"}}`)
	for _, object := range objects {
		err := r.Status().Patch(ctx, object, client.RawPatch(types.MergePatchType, patch))
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *TaskDefinitionReconciler) errorRequeueAfter(ctx context.Context, err error, taskDefinition teachv1alpha1.TaskDefinition, objects ...runtime.Object) (ctrl.Result, error) {
	_ = r.setState(ctx, stateError, &taskDefinition)
	_ = r.setState(ctx, stateError, objects...)
	var errCount int
	if taskDefinition.Status.ErrorCount != nil {
		errCount = *taskDefinition.Status.ErrorCount
	}
	errCount++
	patch := []byte(`{"status":{"errorCount":"` + fmt.Sprint(errCount) + `"}}`)
	_ = r.Status().Patch(ctx, &taskDefinition, client.RawPatch(types.MergePatchType, patch))

	return ctrl.Result{RequeueAfter: time.Duration(errCount) * time.Second * 5}, err
}

// SetupWithManager is used by kubebuilder to init the controller loop
func (r *TaskDefinitionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&teachv1alpha1.TaskDefinition{}).
		Complete(r)
}
