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

package dashboard

import (
	"context"
	"io/ioutil"
	"net/http"

	"github.com/dergeberl/kubeteach/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("metrics tests", func() {
	Context("Run checks", func() {
		ctx := context.Background()
		task := v1alpha1.TaskDefinition{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test",
				Namespace: "default",
			},
			Spec: v1alpha1.TaskDefinitionSpec{
				TaskSpec: v1alpha1.TaskSpec{
					Title:       "test",
					Description: "test",
				},
				TaskConditions: []v1alpha1.TaskCondition{
					{
						APIVersion: "v1",
						Name:       "test",
						Kind:       "Namespace",
					},
				},
			},
		}
		taskState := "active"
		It("apply tasksDefinition", func() {
			Expect(k8sClient.Create(ctx, &task)).Should(Succeed())
		})

		It("apply tasks status", func() {
			task.Status.State = &taskState
			Expect(k8sClient.Status().Update(ctx, &task)).Should(Succeed())
		})

		It("get index", func() {
			resp, err := http.Get("http://localhost:8090/")
			Expect(err).Should(BeNil())
			_, err = ioutil.ReadAll(resp.Body)
			Expect(err).Should(BeNil())
		})

		It("apply tasks", func() {
			resp, err := http.Get("http://localhost:8090/api/tasks")
			Expect(err).Should(BeNil())
			data, err := ioutil.ReadAll(resp.Body)
			Expect(err).Should(BeNil())
			Expect(string(data)).
				Should(Equal("[{\"name\":\"test\",\"namespace\":\"default\",\"title\":\"test\",\"description\":\"test\",\"uid\":\"" + string(task.UID) + "\"}]")) //nolint:lll
		})

		It("apply tasks status", func() {
			resp, err := http.Get("http://localhost:8090/api/taskstatus/" + string(task.UID))
			Expect(err).Should(BeNil())
			data, err := ioutil.ReadAll(resp.Body)
			Expect(err).Should(BeNil())
			Expect(string(data)).Should(Equal("{\"status\":\"active\"}"))
		})

		It("shell endpint", func() {
			resp, err := http.Get("http://localhost:8090/shell/" + string(task.UID))
			Expect(err).Should(BeNil())
			_, err = ioutil.ReadAll(resp.Body)
			Expect(err).Should(BeNil())
		})

	})
})
