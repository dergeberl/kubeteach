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

package metrics

import (
	"context"
	"time"

	"github.com/dergeberl/kubeteach/controllers"
	"k8s.io/utils/pointer"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/prometheus/client_golang/prometheus/testutil"
	ctrl "sigs.k8s.io/controller-runtime"
)

var _ = Describe("metrics tests", func() {
	timeout, retry := time.Second*10, time.Millisecond*300
	Context("Run checks", func() {
		ctx := context.Background()

		It("apply ExerciseSet", func() {
			Expect(k8sClient.Create(ctx, &testsExerciseSet.exerciseSet)).Should(Succeed())
		})

		It("apply ExerciseSet status", func() {
			testsExerciseSet.exerciseSet.Status = testsExerciseSet.status
			Expect(k8sClient.Status().Update(ctx, &testsExerciseSet.exerciseSet)).Should(Succeed())
		})

		It("test metrics count - must be 8", func() {
			Expect(
				testutil.CollectAndCount(
					New(k8sClient, ctrl.Log.WithName("metrics")),
				),
			).Should(BeEquivalentTo(8))
		})

		It("lint metrics", func() {
			Expect(
				testutil.CollectAndLint(
					New(k8sClient, ctrl.Log.WithName("metrics")),
				),
			).Should(BeNil())
		})

		It("check metrics", func() {
			Eventually(func() error {
				for metric, exp := range expectExerciseSet {
					err := testutil.CollectAndCompare(
						New(k8sClient, ctrl.Log.WithName("metrics")),
						exp,
						metric,
					)
					if err != nil {
						return err
					}
				}
				return nil
			}, timeout, retry).Should(Succeed())

		})

		It("clean up ExerciseSet", func() {
			Expect(k8sClient.Delete(ctx, &testsExerciseSet.exerciseSet)).Should(Succeed())
		})

		It("apply tasks", func() {
			Expect(k8sClient.Create(ctx, &testTasks1)).Should(Succeed())
			Expect(k8sClient.Create(ctx, &testTasks2)).Should(Succeed())
			Expect(k8sClient.Create(ctx, &testTasks3)).Should(Succeed())
			Expect(k8sClient.Create(ctx, &testTasks4)).Should(Succeed())
		})

		It("apply ExerciseSet status", func() {
			testTasks1.Status.State = pointer.StringPtr(controllers.StateSuccessful)
			testTasks2.Status.State = pointer.StringPtr(controllers.StateActive)
			testTasks3.Status.State = pointer.StringPtr(controllers.StatePending)
			Expect(k8sClient.Status().Update(ctx, &testTasks1)).Should(Succeed())
			Expect(k8sClient.Status().Update(ctx, &testTasks2)).Should(Succeed())
			Expect(k8sClient.Status().Update(ctx, &testTasks3)).Should(Succeed())
		})

		It("test metrics count - must be 4", func() {
			Expect(
				testutil.CollectAndCount(
					New(k8sClient, ctrl.Log.WithName("metrics")),
				),
			).Should(BeEquivalentTo(4))
		})

		It("lint metrics", func() {
			Expect(
				testutil.CollectAndLint(
					New(k8sClient, ctrl.Log.WithName("metrics")),
				),
			).Should(BeNil())
		})

		It("check metrics", func() {
			Eventually(func() error {
				for metric, exp := range expectTask {
					err := testutil.CollectAndCompare(
						New(k8sClient, ctrl.Log.WithName("metrics")),
						exp,
						metric,
					)
					if err != nil {
						return err
					}
				}
				return nil
			}, timeout, retry).Should(Succeed())
		})

		It("clean up task", func() {
			Expect(k8sClient.Delete(ctx, &testTasks1)).Should(Succeed())
			Expect(k8sClient.Delete(ctx, &testTasks2)).Should(Succeed())
			Expect(k8sClient.Delete(ctx, &testTasks3)).Should(Succeed())
			Expect(k8sClient.Delete(ctx, &testTasks4)).Should(Succeed())
		})

		It("check error", func() {
			testutil.CollectAndCount(
				New(nil, ctrl.Log.WithName("metrics")),
			)
			testutil.CollectAndCount(
				New(nil, nil),
			)
		})
	})
})
