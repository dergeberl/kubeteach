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
	teachv1alpha1 "kubeteach/api/v1alpha1"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
)

type testData struct {
	state          string
	taskDefinition teachv1alpha1.TaskDefinition
}

var _ = Describe("TaskConditions ApplyChecks", func() {
	Context("Run checks in checkItems", func() {
		It("run testcases", func() {
			ctx := context.Background()
			testObjects := []runtime.Object{
				&v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "test2"}},
			}
			By("deploy testdata objects")
			for _, obj := range testObjects {
				Expect(k8sClient.Create(ctx, obj)).Should(Succeed())

			}

			By("Apply testdata task definitions")
			requireTask := "task1"
			tests := []testData{{
				state: stateActive,
				taskDefinition: teachv1alpha1.TaskDefinition{

					ObjectMeta: metav1.ObjectMeta{Name: "task1", Namespace: "default"},
					Spec: teachv1alpha1.TaskDefinitionSpec{
						TaskSpec: teachv1alpha1.TaskSpec{
							Title:       "task1",
							Description: "Task1 description",
							HelpURL:     "HelpURL",
						},
						TaskConditions: []teachv1alpha1.TaskCondition{{
							APIVersion: "v1",
							Kind:       "Namespace",
							APIGroup:   "",
							Name:       "test1",
							ResourceCondition: []teachv1alpha1.ResourceCondition{{
								Field:    "metadata.name",
								Operator: "eq",
								Value:    "test1",
							},
							},
						}},
						RequiredTaskName: nil,
					},
				},
			}, {
				state: stateSuccessful,
				taskDefinition: teachv1alpha1.TaskDefinition{

					ObjectMeta: metav1.ObjectMeta{Name: "task2", Namespace: "default"},
					Spec: teachv1alpha1.TaskDefinitionSpec{
						TaskSpec: teachv1alpha1.TaskSpec{
							Title:       "task2",
							Description: "Task2 description",
							HelpURL:     "HelpURL",
						},
						TaskConditions: []teachv1alpha1.TaskCondition{{
							APIVersion: "v1",
							Kind:       "Namespace",
							APIGroup:   "",
							Name:       "test2",
							ResourceCondition: []teachv1alpha1.ResourceCondition{{
								Field:    "metadata.name",
								Operator: "eq",
								Value:    "test2",
							},
							},
						}},
						RequiredTaskName: nil,
					},
				},
			}, {
				state: statePending,
				taskDefinition: teachv1alpha1.TaskDefinition{

					ObjectMeta: metav1.ObjectMeta{Name: "task3", Namespace: "default"},
					Spec: teachv1alpha1.TaskDefinitionSpec{
						TaskSpec: teachv1alpha1.TaskSpec{
							Title:       "task3",
							Description: "Task3 description",
							HelpURL:     "HelpURL",
						},
						TaskConditions: []teachv1alpha1.TaskCondition{{
							APIVersion: "v1",
							Kind:       "Namespace",
							APIGroup:   "",
							Name:       "test3",
							ResourceCondition: []teachv1alpha1.ResourceCondition{{
								Field:    "metadata.name",
								Operator: "eq",
								Value:    "test3",
							},
							},
						}},
						RequiredTaskName: &requireTask,
					},
				},
			},
			}
			for _, test := range tests {
				Expect(k8sClient.Create(ctx, &test.taskDefinition)).Should(Succeed())
			}
			for _, testdata := range tests {
				curTask := &teachv1alpha1.Task{}
				Eventually(func() error {
					err := k8sClient.Get(ctx, types.NamespacedName{Name: testdata.taskDefinition.Name, Namespace: testdata.taskDefinition.Namespace}, curTask) // nolint:lll
					if err != nil {
						return err
					}
					if curTask.Status.State != nil && *curTask.Status.State == testdata.state {
						return nil
					}
					if curTask.Status.State != nil {
						return errors.New("not expected state " + *curTask.Status.State + " in task " + testdata.taskDefinition.Name)
					}
					return errors.New("task has no status in task " + testdata.taskDefinition.Name)
				}, time.Second*5, time.Second*1).Should(Succeed())
			}

			By("update task spec")
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
			}, time.Second*15, time.Second*1).Should(Succeed())

			testObjects2 := []runtime.Object{
				&v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "test1"}},
				&v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "test3"}},
			}
			By("deploy testdata objects")
			for _, obj := range testObjects2 {
				Expect(k8sClient.Create(ctx, obj)).Should(Succeed())
			}
			for _, test := range tests {
				curTask := &teachv1alpha1.Task{}
				Eventually(func() error {
					err := k8sClient.Get(ctx, types.NamespacedName{Name: test.taskDefinition.Name, Namespace: test.taskDefinition.Namespace}, curTask) // nolint:lll
					if err != nil {
						return err
					}
					if curTask.Status.State != nil && *curTask.Status.State == stateSuccessful {
						return nil
					}
					return errors.New("not expected state")
				}, time.Second*5, time.Second*1).Should(Succeed())
			}

			By("delete Taskdefinitions")

			for _, test := range tests {
				Expect(k8sClient.Delete(ctx, &test.taskDefinition)).Should(Succeed())
			}
		})
	})
})
