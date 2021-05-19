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

var _ = Describe("Test exerciseSet api with creation and deletion on k8s api", func() {
	Context("Task Type tests", func() {
		ctx := context.Background()

		It("test validation", func() {
			for _, test := range exerciseSetCases {
				Expect(k8sClient.Create(ctx, test.obj)).Should(test.err)
			}
		})
		It("test deepcopy exerciseSet", func() {
			exerciseSet := &ExerciseSet{
				ObjectMeta: metav1.ObjectMeta{Name: "exerciseSet", Namespace: "default"},
				Spec: ExerciseSetSpec{
					TaskDefinitions: []ExerciseSetSpecTaskDefinitions{
						{
							TaskDefinitionSpec: TaskDefinitionSpec{
								TaskSpec: TaskSpec{
									Title:           "Test1",
									Description:     "Test1",
									LongDescription: "Test1",
									HelpURL:         "Test1",
								},
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
						},
					},
				},
				Status: ExerciseSetStatus{
					NumberOfTasks:              1,
					NumberOfActiveTasks:        1,
					NumberOfPendingTasks:       1,
					NumberOfSuccessfulTasks:    1,
					NumberOfUnknownTasks:       1,
					NumberOfTasksWithoutPoints: 1,
					PointsTotal:                1,
					PointsAchieved:             1,
				},
			}
			Expect(reflect.DeepEqual(exerciseSet, exerciseSet.DeepCopyObject())).Should(BeTrue())
			Expect(reflect.DeepEqual(exerciseSet, exerciseSet.DeepCopy())).Should(BeTrue())
			Expect(reflect.DeepEqual(exerciseSet.Spec, *exerciseSet.Spec.DeepCopy())).Should(BeTrue())
			Expect(reflect.DeepEqual(exerciseSet.Status, *exerciseSet.Status.DeepCopy())).Should(BeTrue())
			Expect(reflect.DeepEqual(exerciseSet.Spec.TaskDefinitions[0], *exerciseSet.Spec.TaskDefinitions[0].DeepCopy())).Should(BeTrue())
			exerciseSet = nil
			Expect(reflect.DeepEqual(nil, exerciseSet.DeepCopyObject())).Should(BeTrue())
		})
		It("test deepcopy list exerciseSetList", func() {
			exerciseSetList := &ExerciseSetList{}
			Expect(k8sClient.List(ctx, exerciseSetList)).Should(Succeed())
			Expect(reflect.DeepEqual(exerciseSetList, exerciseSetList.DeepCopyObject())).Should(BeTrue())
			Expect(reflect.DeepEqual(*exerciseSetList, *exerciseSetList.DeepCopy())).Should(BeTrue())
			exerciseSetList = nil
			Expect(reflect.DeepEqual(nil, exerciseSetList.DeepCopyObject())).Should(BeTrue())
		})

	})
})
