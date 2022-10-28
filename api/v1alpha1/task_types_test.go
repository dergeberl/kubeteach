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

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Test task api with creation and deletion on k8s api", func() {
	Context("Task Type tests", func() {
		ctx := context.Background()

		It("test validation", func() {
			for _, test := range taskCases {
				Expect(k8sClient.Create(ctx, test.obj)).Should(test.err)
			}
		})

		It("test deepcopy task", func() {
			status := "pending"
			task := &Task{
				ObjectMeta: metav1.ObjectMeta{Name: "task", Namespace: "default"},
				Spec: TaskSpec{
					Title:           "Test1",
					Description:     "Test1",
					LongDescription: "Test1",
					HelpURL:         "Test1",
				},
				Status: TaskStatus{State: &status},
			}
			Expect(reflect.DeepEqual(task, task.DeepCopyObject())).Should(BeTrue())
			Expect(reflect.DeepEqual(*task, *task.DeepCopy())).Should(BeTrue())
			Expect(reflect.DeepEqual(task.Spec, *task.Spec.DeepCopy())).Should(BeTrue())
			Expect(reflect.DeepEqual(task.Status, *task.Status.DeepCopy())).Should(BeTrue())
			task = nil
			Expect(reflect.DeepEqual(nil, task.DeepCopyObject())).Should(BeTrue())
		})
		It("test deepcopy list task", func() {
			taskList := &TaskList{}
			Expect(k8sClient.List(ctx, taskList)).Should(Succeed())
			Expect(reflect.DeepEqual(taskList, taskList.DeepCopyObject())).Should(BeTrue())
			Expect(reflect.DeepEqual(*taskList, *taskList.DeepCopy())).Should(BeTrue())
		})
	})
})
