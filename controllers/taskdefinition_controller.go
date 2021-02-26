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
	"github.com/go-logr/logr"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	kubeteachv1 "kubeteach/api/v1"
	"kubeteach/controllers/check"
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
	requeueTime     = time.Duration(2) * time.Second
)

// TaskDefinitionReconciler reconciles a TaskDefinition object
type TaskDefinitionReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=kubeteach.geberl.io,resources=taskdefinitions,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=kubeteach.geberl.io,resources=taskdefinitions/status,verbs=get;update;patch

func (r *TaskDefinitionReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	_ = r.Log.WithValues("taskdefinition", req.NamespacedName)

	// get taskDefinition
	taskDefinition := kubeteachv1.TaskDefinition{}
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

	//TODO kann weg
	//if taskDefinition.Status.State == nil {
	//	if taskDefinition.Spec.RequiredTaskName != nil {
	//		err = r.setState(ctx, statePending, &taskDefinition.ObjectMeta, nil)
	//		if err != nil {
	//			return reconcile.Result{}, err
	//		}
	//	} else {
	//		err = r.setState(ctx, stateActive, &taskDefinition.ObjectMeta, nil)
	//		if err != nil {
	//			return reconcile.Result{}, err
	//		}
	//	}
	//	return ctrl.Result{}, nil
	//}

	//skip successfull objects
	if taskDefinition.Status.State == stateSuccessful {
		return ctrl.Result{}, nil
	}

	//create or update Task
	taskMeta, err := r.createOrUpdateTask(ctx, &taskDefinition)

	//check pending state
	if taskDefinition.Status.State == statePending || taskDefinition.Status.State == "" {
		if taskDefinition.Spec.RequiredTaskName != nil {
			reqTask := kubeteachv1.TaskDefinition{}
			err = r.Client.Get(ctx, client.ObjectKey{Name: *taskDefinition.Spec.RequiredTaskName, Namespace: req.Namespace}, &reqTask)
			if err != nil {
				if errors.IsNotFound(err) {
					return ctrl.Result{}, nil
				}
				return ctrl.Result{}, err
			}
			if reqTask.Status.State == stateSuccessful {
				err = r.setState(ctx, stateActive, &taskDefinition.ObjectMeta, &taskMeta)
				if err != nil {
					return reconcile.Result{}, err
				}
			}
		} else {
			err = r.setState(ctx, stateActive, &taskDefinition.ObjectMeta, &taskMeta)
			if err != nil {
				return reconcile.Result{}, err
			}
		}
		if taskDefinition.Status.State == "" {
			err = r.setState(ctx, statePending, &taskDefinition.ObjectMeta, nil)
			if err != nil {
				return reconcile.Result{}, err
			}
		}
		return ctrl.Result{RequeueAfter: requeueTime}, nil
	}

	checkController := check.CheckController{
		Client:     r.Client,
		Log:        r.Log,
		CheckItems: taskDefinition.Spec.TaskConditions,
	}

	status, err := checkController.ApplyChecks(ctx)
	if err != nil {
		return ctrl.Result{}, err
	}

	if status {
		err = r.setState(ctx, stateSuccessful, &taskDefinition.ObjectMeta, &taskMeta)
		if err != nil {
			return reconcile.Result{}, err
		}
		err = r.Client.Create(ctx, &v1.Event{
			ObjectMeta: metav1.ObjectMeta{
				Namespace:    taskMeta.Namespace,
				GenerateName: taskMeta.Name,
			},
			InvolvedObject: v1.ObjectReference{
				Kind:       "Task",
				Namespace:  taskMeta.Namespace,
				Name:       taskMeta.Name,
				UID:        taskMeta.UID,
				APIVersion: "v1",
			},
			Reason:              "success",
			Message:             "Task successfully done",
			Source:              v1.EventSource{Component: "kubeteach"},
			Type:                v1.EventTypeNormal,
			FirstTimestamp:      metav1.Now(),
			Series:              nil,
			Action:              "",
			Related:             nil,
			ReportingController: "",
			ReportingInstance:   "",
		})
		if err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	return ctrl.Result{RequeueAfter: time.Duration(5) * time.Second}, nil
}

func (r *TaskDefinitionReconciler) createOrUpdateTask(ctx context.Context, taskDefinition *kubeteachv1.TaskDefinition) (metav1.ObjectMeta, error) {

	taskList := kubeteachv1.TaskList{}
	r.Client.List(ctx, &taskList)
	var task *kubeteachv1.Task
	for _, taskTtem := range taskList.Items {
		if taskTtem.OwnerReferences[0].UID == taskDefinition.UID {
			//found task
			task = &taskTtem
			break
		}
	}

	//create task if not found
	if task == nil {
		task = &kubeteachv1.Task{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Task",
				APIVersion: "kubeteach.geberl.io/v1",
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
			Status: kubeteachv1.TaskStatus{State: &taskDefinition.Status.State},
		}
		err := r.Client.Create(ctx, task)
		if err != nil {
			return metav1.ObjectMeta{}, err
		}
		return task.ObjectMeta, nil
	}

	//sync spec
	if !reflect.DeepEqual(task.Spec, taskDefinition.Spec.TaskSpec) {
		task.Spec = taskDefinition.Spec.TaskSpec
		err := r.Update(ctx, task)
		if err != nil {
			return metav1.ObjectMeta{}, err
		}
	}

	//sync status
	if taskDefinition.Status.State != *task.Status.State {
		patch := []byte(`{"status":{"state":"` + taskDefinition.Status.State + `"}}`)
		if err := r.Client.Patch(ctx, task, client.RawPatch(types.MergePatchType, patch)); err != nil {
			return metav1.ObjectMeta{}, err
		}
	}
	return task.ObjectMeta, nil
}

func (r *TaskDefinitionReconciler) setState(ctx context.Context, state string, taskDefinition, task *metav1.ObjectMeta) error {
	patch := []byte(`{"status":{"state":"` + state + `"}}`)

	if taskDefinition != nil {
		err := r.Client.Patch(ctx, &kubeteachv1.TaskDefinition{ObjectMeta: *taskDefinition}, client.RawPatch(types.MergePatchType, patch))
		if err != nil {
			return err
		}
	}

	if task != nil {
		err := r.Client.Patch(ctx, &kubeteachv1.Task{ObjectMeta: *task}, client.RawPatch(types.MergePatchType, patch))
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *TaskDefinitionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kubeteachv1.TaskDefinition{}).
		Complete(r)
}
