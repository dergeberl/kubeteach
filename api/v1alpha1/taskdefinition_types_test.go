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

package v1alpha1

import (
	"context"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Test task api with creation and deletion on k8s api", func() {
	Context("Task Type tests", func() {
		ctx := context.Background()

		It("test validation", func() {
			for _, test := range taskDefinitionCases {
				Expect(k8sClient.Create(ctx, test.obj)).Should(test.err)
			}
		})
		It("test deepcopy taskDefinition", func() {
			requiredTaskName := "task1"
			status := "pending"
			taskDefinition := &TaskDefinition{
				ObjectMeta: metav1.ObjectMeta{Name: "task", Namespace: "default"},
				Spec: TaskDefinitionSpec{
					TaskSpec: TaskSpec{
						Title:           "Test1",
						Description:     "Test1",
						LongDescription: "Test1",
						HelpURL:         "Test1",
					},
					RequiredTaskName: &requiredTaskName,
					TaskConditions: []TaskCondition{
						{
							APIVersion: "v1",
							Kind:       "Namespace",
							Name:       "default",
							Namespace:  "",
							NotExists:  false,
							ResourceCondition: []ResourceCondition{
								{
									Field:    "meta.name",
									Operator: "eq",
									Value:    "name",
								},
							},
						},
					},
				},
				Status: TaskDefinitionStatus{
					State: &status,
				},
			}
			Expect(reflect.DeepEqual(taskDefinition, taskDefinition.DeepCopyObject())).Should(BeTrue())
			Expect(reflect.DeepEqual(taskDefinition, taskDefinition.DeepCopy())).Should(BeTrue())
			Expect(reflect.DeepEqual(taskDefinition.Spec, *taskDefinition.Spec.DeepCopy())).Should(BeTrue())
			Expect(reflect.DeepEqual(taskDefinition.Status, *taskDefinition.Status.DeepCopy())).Should(BeTrue())
			Expect(reflect.DeepEqual(taskDefinition.Spec.TaskConditions[0], *taskDefinition.Spec.TaskConditions[0].DeepCopy())).Should(BeTrue())
			Expect(reflect.DeepEqual(
				taskDefinition.Spec.TaskConditions[0].ResourceCondition[0],
				*taskDefinition.Spec.TaskConditions[0].ResourceCondition[0].DeepCopy())).Should(BeTrue())
			taskDefinition = nil
			Expect(reflect.DeepEqual(nil, taskDefinition.DeepCopyObject())).Should(BeTrue())

		})
		It("test deepcopy list taskDefinition", func() {
			taskDefinitionList := &TaskDefinitionList{}
			Expect(k8sClient.List(ctx, taskDefinitionList)).Should(Succeed())
			Expect(reflect.DeepEqual(taskDefinitionList, taskDefinitionList.DeepCopyObject())).Should(BeTrue())
			Expect(reflect.DeepEqual(*taskDefinitionList, *taskDefinitionList.DeepCopy())).Should(BeTrue())
		})

	})
})
