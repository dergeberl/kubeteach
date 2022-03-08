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
	"errors"
	"fmt"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	teachv1alpha1 "github.com/dergeberl/kubeteach/api/v1alpha1"
)

var _ = Describe("TaskConditions tests", func() {
	timeout, retry := time.Second*5, time.Millisecond*300
	Context("Run checks in checkItems", func() {
		It("apply taskDefinitions", func() {
			for _, test := range testsTaskDefinition {
				Expect(k8sClient.Create(ctx, &test.taskDefinition)).Should(Succeed())
			}
		})

		It("apply initial test objects", func() {
			for _, test := range testsTaskDefinition {
				if test.initialDeploy == nil {
					continue
				}
				Expect(k8sClient.Create(ctx, test.initialDeploy)).Should(Succeed())
			}
		})

		It("check state after initial objects", func() {
			for _, testdata := range testsTaskDefinition {
				curTask := &teachv1alpha1.Task{}
				Eventually(func() error {
					err := k8sClient.Get(ctx,
						types.NamespacedName{Name: testdata.taskDefinition.Name, Namespace: testdata.taskDefinition.Namespace},
						curTask)
					if err != nil {
						return err
					}
					if curTask.Status.State != nil && *curTask.Status.State == testdata.state {
						return nil
					}
					if curTask.Status.State != nil {
						return fmt.Errorf("got state %v but want %v in task %v", curTask.Status.State, testdata.state, curTask.Name)
					}
					return fmt.Errorf("got no state but want %v in task %v", testdata.state, curTask.Name)
				}, timeout, retry).Should(Succeed())
			}
		})

		It("test taskSpec update", func() {
			task1 := &teachv1alpha1.TaskDefinition{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: "task1", Namespace: "default"}, task1)).Should(Succeed())
			task1.Spec.TaskSpec.HelpURL = "newURL"
			Expect(k8sClient.Update(ctx, task1)).Should(Succeed())

			Eventually(func() error {
				task := &teachv1alpha1.Task{}
				err := k8sClient.Get(ctx, types.NamespacedName{Name: "task1", Namespace: "default"}, task)
				if err != nil {
					return err
				}
				if task.Spec.HelpURL == task1.Spec.TaskSpec.HelpURL {
					return nil
				}
				return errors.New("no update")
			}, timeout, retry).Should(Succeed())
		})

		It("recreate task", func() {
			task1 := &teachv1alpha1.Task{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: "task1", Namespace: "default"}, task1)).Should(Succeed())
			Expect(k8sClient.Delete(ctx, task1)).Should(Succeed())
			Eventually(func() error {
				task := &teachv1alpha1.Task{}
				return k8sClient.Get(ctx, types.NamespacedName{Name: "task1", Namespace: "default"}, task)
			}, timeout, retry).Should(Succeed())
		})

		It("apply solutions", func() {
			for _, test := range testsTaskDefinition {
				if test.solution != nil {
					Expect(k8sClient.Create(ctx, test.solution)).Should(Succeed())
				}
			}
		})

		It("check if every test is successful", func() {
			for _, test := range testsTaskDefinition {
				By(test.taskDefinition.Name)
				curTask := &teachv1alpha1.TaskDefinition{}
				Eventually(func() error {
					err := k8sClient.Get(ctx,
						types.NamespacedName{Name: test.taskDefinition.Name, Namespace: test.taskDefinition.Namespace},
						curTask)
					if err != nil {
						return err
					}
					if curTask.Status.State != nil && *curTask.Status.State == StateSuccessful {
						return nil
					}
					if curTask.Status.State != nil {
						return fmt.Errorf("got state %v but want %v in task %v", *curTask.Status.State, StateSuccessful, curTask)
					}
					return fmt.Errorf("got no state but want %v in task %v", StateSuccessful, curTask)
				}, timeout, retry).Should(Succeed())
			}

		})

		It("delete all tests", func() {
			for _, test := range testsTaskDefinition {
				Expect(k8sClient.Delete(ctx, &test.taskDefinition)).Should(Succeed())
			}
		})

		It("check deletion", func() {
			for _, test := range testsTaskDefinition {
				Eventually(func() error {
					curTask := &teachv1alpha1.TaskDefinition{}
					err := k8sClient.Get(ctx,
						types.NamespacedName{Name: test.taskDefinition.Name, Namespace: test.taskDefinition.Namespace},
						curTask)
					if err == nil {
						return fmt.Errorf("taskdefinition still exists %v", test.taskDefinition.Name)
					}
					return client.IgnoreNotFound(err)
				}, timeout, retry).Should(Succeed())
			}
		})

		It("check failed condition event", func() {
			taskDefinition := teachv1alpha1.TaskDefinition{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "failed-condition",
					Namespace: "default",
				},
				Spec: teachv1alpha1.TaskDefinitionSpec{
					TaskSpec: teachv1alpha1.TaskSpec{
						Title:       "failed-condition",
						Description: "failed-condition",
					},
					TaskConditions: []teachv1alpha1.TaskCondition{
						{
							APIVersion: "v1",
							Kind:       "WrongKind",
							APIGroup:   "",
							Name:       "failed-condition",
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, &taskDefinition)).Should(Succeed())
			Eventually(func() error {
				eventList := &v1.EventList{}
				err := k8sClient.List(ctx, eventList)
				if err != nil {
					return err
				}
				for _, event := range eventList.Items {
					if event.InvolvedObject.Name == taskDefinition.Name {
						if event.Type == v1.EventTypeWarning &&
							strings.Contains(event.Message, "Conditions apply fail with error") {
							return nil
						}
					}
				}

				return errors.New("failed event not found")
			}, timeout, retry).Should(Succeed())
			Expect(k8sClient.Delete(ctx, &taskDefinition)).Should(Succeed())
		})

	})
})
