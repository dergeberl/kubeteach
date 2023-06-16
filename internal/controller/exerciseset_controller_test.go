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

package controller

import (
	"errors"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	teachv1alpha1 "github.com/dergeberl/kubeteach/api/v1alpha1"
)

var _ = Describe("ExerciseSet tests", func() {
	timeout, retry := time.Second*10, time.Millisecond*300
	Context("Run checks", func() {
		It("apply ExerciseSet", func() {
			Expect(k8sClient.Create(ctx, &testsExerciseSet.exerciseSet)).Should(Succeed())
		})

		It("apply initial test objects", func() {
			for _, test := range testsExerciseSet.initialDeploy {
				Expect(k8sClient.Create(ctx, test)).Should(Succeed())
			}
		})

		It("check state after initial objects", func() {

			Eventually(func() error {
				curExerciseSet := &teachv1alpha1.ExerciseSet{}
				err := k8sClient.Get(ctx,
					types.NamespacedName{
						Name:      testsExerciseSet.exerciseSet.Name,
						Namespace: testsExerciseSet.exerciseSet.Namespace},
					curExerciseSet)
				if err != nil {
					return err
				}
				if curExerciseSet.Status.NumberOfTasks != testsExerciseSet.status.NumberOfTasks {
					return errors.New("NumberOfTasks in status is wrong")
				}
				if curExerciseSet.Status.NumberOfSuccessfulTasks != testsExerciseSet.status.NumberOfSuccessfulTasks {
					return errors.New("NumberOfSuccessfulTasks in status is wrong")
				}
				if curExerciseSet.Status.NumberOfActiveTasks != testsExerciseSet.status.NumberOfActiveTasks {
					return errors.New("NumberOfActiveTasks in status is wrong")
				}
				if curExerciseSet.Status.NumberOfPendingTasks != testsExerciseSet.status.NumberOfPendingTasks {
					return errors.New("NumberOfPendingTasks in status is wrong")
				}
				if curExerciseSet.Status.NumberOfUnknownTasks != testsExerciseSet.status.NumberOfUnknownTasks {
					return errors.New("NumberOfUnknownTasks in status is wrong")
				}
				if curExerciseSet.Status.NumberOfTasksWithoutPoints != testsExerciseSet.status.NumberOfTasksWithoutPoints {
					return errors.New("NumberOfTasksWithoutPoints in status is wrong")
				}
				if curExerciseSet.Status.PointsTotal != testsExerciseSet.status.PointsTotal {
					return errors.New("PointsTotal in status is wrong")
				}
				if curExerciseSet.Status.PointsAchieved != testsExerciseSet.status.PointsAchieved {
					return errors.New("PointsAchieved in status is wrong")
				}
				return nil
			}, timeout, retry).Should(Succeed())

		})

		It("test taskDefinitionSpec update", func() {
			taskDefinitionEdit := &teachv1alpha1.TaskDefinition{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: "exerciseset1-5", Namespace: "default"}, taskDefinitionEdit)).Should(Succeed())
			taskDefinitionEdit.Spec.TaskSpec.Title = "NewName"
			Expect(k8sClient.Update(ctx, taskDefinitionEdit)).Should(Succeed())

			Eventually(func() error {
				taskDefinition := &teachv1alpha1.TaskDefinition{}
				err := k8sClient.Get(ctx, types.NamespacedName{Name: "exerciseset1-5", Namespace: "default"}, taskDefinition)
				if err != nil {
					return err
				}
				if taskDefinition.Spec.TaskSpec.Title == "exerciseset1-5" {
					return nil
				}
				return errors.New("no update")
			}, timeout, retry).Should(Succeed())
		})

		It("test taskDefinition owner fix", func() {
			taskDefinitionEdit := &teachv1alpha1.TaskDefinition{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: "exerciseset1-6", Namespace: "default"}, taskDefinitionEdit)).Should(Succeed())
			taskDefinitionEdit.OwnerReferences = []v1.OwnerReference{}
			Expect(k8sClient.Update(ctx, taskDefinitionEdit)).Should(Succeed())

			Eventually(func() error {
				taskDefinition := &teachv1alpha1.TaskDefinition{}
				err := k8sClient.Get(ctx, types.NamespacedName{Name: "exerciseset1-6", Namespace: "default"}, taskDefinition)
				if err != nil {
					return err
				}
				if len(taskDefinition.OwnerReferences) == 1 {
					return nil
				}
				return errors.New("no update")
			}, timeout, retry).Should(Succeed())
		})

		It("test clean up", func() {
			Expect(k8sClient.Delete(ctx, &testsExerciseSet.exerciseSet)).Should(Succeed())
		})
	})
})
