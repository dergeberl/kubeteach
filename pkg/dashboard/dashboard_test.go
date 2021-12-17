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
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/dergeberl/kubeteach/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("dashboard tests", func() {
	Context("Run checks", func() {
		timeout, retry := time.Second*5, time.Millisecond*100
		ctx := context.Background()
		var dashboard1 Config
		dashboard1listen := "localhost:8090"
		webterminalBasicAuthUser := "webterminaltestuser"
		webterminalBasicAuthPass := "webterminaltestpw"
		webterminalListen := "localhost:8079"

		var dashboard2 Config
		dashboard2listen := "localhost:8091"

		var dashboard3 Config
		dashboard3listen := "localhost:8092"
		basicAuthUser := "testuser"
		basicAuthPass := "testpw"

		task1 := v1alpha1.TaskDefinition{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test1",
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
		task2 := task1
		task2.Name = "test2"
		taskState := "active"

		It("apply tasksDefinition", func() {
			Expect(k8sClient.Create(ctx, &task1)).Should(Succeed())
			Expect(k8sClient.Create(ctx, &task2)).Should(Succeed())
		})

		It("apply tasks status", func() {
			task1.Status.State = &taskState
			task2.Status.State = &taskState
			Expect(k8sClient.Status().Update(ctx, &task1)).Should(Succeed())
			Expect(k8sClient.Status().Update(ctx, &task2)).Should(Succeed())
		})

		It("create shell dummy", func() {
			r := chi.NewRouter()
			r.Use(middleware.BasicAuth("", map[string]string{webterminalBasicAuthUser: webterminalBasicAuthPass}))
			r.Route("/shell", func(r chi.Router) {
				r.HandleFunc("/*", func(writer http.ResponseWriter, request *http.Request) {
					_, _ = fmt.Fprint(writer, "OK")
				})
			})
			go func() {
				err := http.ListenAndServe(webterminalListen, r)
				Expect(err).ToNot(HaveOccurred())
			}()
		})

		It("create dashboard1", func() {
			dashboard1 = New(k8sClient,
				dashboard1listen,
				"../../dashboard/dist/",
				"",
				"",
				true,
				"localhost",
				"8079",
				webterminalBasicAuthUser+":"+webterminalBasicAuthPass)
			go func() {
				err := dashboard1.Run()
				Expect(err).ToNot(HaveOccurred())
			}()
		})

		It("get index", func() {
			var resp *http.Response
			var err error
			Eventually(func() error {
				resp, err = http.Get("http://" + dashboard1listen + "/index.html")
				return err
			}, timeout, retry).Should(BeNil())
			_, err = ioutil.ReadAll(resp.Body)
			Expect(err).Should(BeNil())
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))

		})

		It("get tasks", func() {
			var resp *http.Response
			var err error
			Eventually(func() error {
				resp, err = http.Get("http://" + dashboard1listen + "/api/tasks")
				return err
			}, timeout, retry).Should(BeNil())
			data, err := ioutil.ReadAll(resp.Body)
			Expect(err).Should(BeNil())
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))
			Expect(string(data)).
				Should(Equal("[{\"name\":\"test1\",\"namespace\":\"default\",\"title\":\"test\",\"description\":\"test\",\"uid\":\"" + string(task1.UID) + "\"}," + //nolint:lll
					"{\"name\":\"test2\",\"namespace\":\"default\",\"title\":\"test\",\"description\":\"test\",\"uid\":\"" + string(task2.UID) + "\"}]")) //nolint:lll
		})

		It("get tasks status", func() {
			var resp *http.Response
			var err error
			Eventually(func() error {
				resp, err = http.Get("http://" + dashboard1listen + "/api/taskstatus/" + string(task1.UID))
				return err
			}, timeout, retry).Should(BeNil())
			data, err := ioutil.ReadAll(resp.Body)
			Expect(err).Should(BeNil())
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))
			Expect(string(data)).Should(Equal("{\"status\":\"active\"}"))
		})

		It("get tasks status - fail no task found", func() {
			var resp *http.Response
			var err error
			Eventually(func() error {
				resp, err = http.Get("http://" + dashboard1listen + "/api/taskstatus/wrong-id")
				return err
			}, timeout, retry).Should(BeNil())
			_, err = ioutil.ReadAll(resp.Body)
			Expect(err).Should(BeNil())
			Expect(resp.StatusCode).Should(Equal(http.StatusNotFound))
		})

		It("get shell endpoint", func() {
			var resp *http.Response
			var err error
			Eventually(func() error {
				resp, err = http.Get("http://" + dashboard1listen + "/shell/")
				return err
			}, timeout, retry).Should(BeNil())
			_, err = ioutil.ReadAll(resp.Body)
			Expect(err).Should(BeNil())
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))

		})

		It("create dashboard2 without k8s client", func() {
			dashboard2 = New(nil,
				dashboard2listen,
				"./dashboard/dist/",
				"",
				"",
				false,
				"",
				"",
				"")
			go func() {
				err := dashboard2.Run()
				Expect(err).ToNot(HaveOccurred())
			}()
		})

		It("get tasks - fail", func() {
			var resp *http.Response
			var err error
			Eventually(func() error {
				resp, err = http.Get("http://" + dashboard2listen + "/api/tasks")
				return err
			}, timeout, retry).Should(BeNil())
			_, err = ioutil.ReadAll(resp.Body)
			Expect(err).Should(BeNil())
			Expect(resp.StatusCode).Should(Equal(http.StatusInternalServerError))
		})
		It("get taskstatus - fail", func() {
			var resp *http.Response
			var err error
			Eventually(func() error {
				resp, err = http.Get("http://" + dashboard2listen + "/api/taskstatus/" + string(task1.UID))
				return err
			}, timeout, retry).Should(BeNil())
			_, err = ioutil.ReadAll(resp.Body)
			Expect(err).Should(BeNil())
			Expect(resp.StatusCode).Should(Equal(http.StatusInternalServerError))
		})

		It("create dashboard3 with basic auth", func() {
			dashboard3 = New(k8sClient,
				dashboard3listen,
				"./dashboard/dist/",
				basicAuthUser,
				basicAuthPass,
				false,
				"",
				"",
				"")
			go func() {
				err := dashboard3.Run()
				Expect(err).ToNot(HaveOccurred())
			}()
		})
		It("get tasks - fail 401", func() {
			var resp *http.Response
			var err error
			Eventually(func() error {
				resp, err = http.Get("http://" + dashboard3listen + "/api/tasks")
				return err
			}, timeout, retry).Should(BeNil())
			_, err = ioutil.ReadAll(resp.Body)
			Expect(err).Should(BeNil())
			Expect(resp.StatusCode).Should(Equal(http.StatusUnauthorized))
		})

		It("get tasks - with basic auth", func() {
			req, err := http.NewRequest("GET", "http://"+dashboard3listen+"/api/tasks", nil)
			Expect(err).Should(BeNil())
			req.SetBasicAuth(basicAuthUser, basicAuthPass)
			var resp *http.Response
			Eventually(func() error {
				resp, err = (&http.Client{Timeout: time.Second * 4}).Do(req)
				return err
			}, timeout, retry).Should(BeNil())
			_, err = ioutil.ReadAll(resp.Body)
			Expect(err).Should(BeNil())
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))
		})
	})
})
