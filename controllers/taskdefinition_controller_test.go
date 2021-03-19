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
	"errors"
	"fmt"
	teachv1alpha1 "kubeteach/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("TaskConditions tests", func() {
	Context("Run checks in checkItems", func() {
		ctx := context.Background()

		It("apply taskDefinitions", func() {
			for _, test := range tests {
				Expect(k8sClient.Create(ctx, &test.taskDefinition)).Should(Succeed())
			}
		})

		It("apply initial test objects", func() {
			for _, test := range tests {
				if test.initialDeploy == nil {
					continue
				}
				Expect(k8sClient.Create(ctx, test.initialDeploy)).Should(Succeed())
			}
		})

		It("check state after initial objects", func() {
			for _, testdata := range tests {
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
				}, time.Second*5, time.Second*1).Should(Succeed())
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
			}, time.Second*5, time.Second*1).Should(Succeed())
		})

		It("recreate task", func() {
			task1 := &teachv1alpha1.Task{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: "task1", Namespace: "default"}, task1)).Should(Succeed())
			Expect(k8sClient.Delete(ctx, task1)).Should(Succeed())

			Eventually(func() error {
				task := &teachv1alpha1.Task{}
				err := k8sClient.Get(ctx, types.NamespacedName{Name: "task1", Namespace: "default"}, task)
				if err != nil {
					return err
				}
				return nil
			}, time.Second*5, time.Second*1).Should(Succeed())
		})

		It("apply solutions", func() {
			for _, test := range tests {
				if test.solution != nil {
					Expect(k8sClient.Create(ctx, test.solution)).Should(Succeed())
				}
			}
		})

		It("check if every test is successful", func() {
			for _, test := range tests {
				By(test.taskDefinition.Name)
				curTask := &teachv1alpha1.TaskDefinition{}
				Eventually(func() error {
					err := k8sClient.Get(ctx,
						types.NamespacedName{Name: test.taskDefinition.Name, Namespace: test.taskDefinition.Namespace},
						curTask)
					if err != nil {
						return err
					}
					if curTask.Status.State != nil && *curTask.Status.State == stateSuccessful {
						return nil
					}
					if curTask.Status.State != nil {
						return fmt.Errorf("got state %v but want %v in task %v", *curTask.Status.State, stateSuccessful, curTask)
					}
					return fmt.Errorf("got no state but want %v in task %v", stateSuccessful, curTask)
				}, time.Second*5, time.Second*1).Should(Succeed())
			}

		})

		It("delete all tests", func() {
			for _, test := range tests {
				Expect(k8sClient.Delete(ctx, &test.taskDefinition)).Should(Succeed())
			}
		})

		It("check deletion", func() {
			for _, test := range tests {
				Eventually(func() error {
					curTask := &teachv1alpha1.TaskDefinition{}
					err := k8sClient.Get(ctx,
						types.NamespacedName{Name: test.taskDefinition.Name, Namespace: test.taskDefinition.Namespace},
						curTask)
					if err == nil {
						return fmt.Errorf("taskdefinition still exists %v", test.taskDefinition.Name)
					}
					return client.IgnoreNotFound(err)
				}, time.Second*5, time.Second*1).Should(Succeed())
			}
		})

	})
})
