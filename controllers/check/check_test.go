/*


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

package check

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeteachv1 "kubeteach/api/v1"
)

var _ = Describe("TaskConditions ApplyChecks", func() {
	Context("Run checks in checkItems", func() {
		It("Should get the expected results for the check", func() {
			ctx := context.Background()
			By("ResourceCondition if a Namespace equal metadata.name")
			object := &v1.Namespace{
				ObjectMeta: metav1.ObjectMeta{Name: "kubeteach"},
				Spec:       v1.NamespaceSpec{},
			}
			Expect(k8sClient.Create(ctx, object)).Should(Succeed())
			c := CheckController{
				Client: k8sClient,
				Log:    nil,
				CheckItems: []kubeteachv1.TaskCondition{{
					ApiVersion: "v1",
					Kind:       "Namespace",
					MatchAll:   false,
					ResourceCondition: []kubeteachv1.ResourceCondition{{
						Field:    "metadata.name",
						Operator: "eq",
						Value:    "kubeteach",
					},
					},
				},
				},
			}
			Expect(c.ApplyChecks(ctx)).Should(BeTrue())

			By("ResourceCondition if a Namespace contains metadata.name")
			object = &v1.Namespace{
				ObjectMeta: metav1.ObjectMeta{Name: "kubeteach-contains"},
				Spec:       v1.NamespaceSpec{},
			}
			Expect(k8sClient.Create(ctx, object)).Should(Succeed())
			c = CheckController{
				Client: k8sClient,
				Log:    nil,
				CheckItems: []kubeteachv1.TaskCondition{{
					ApiVersion: "v1",
					Kind:       "Namespace",
					MatchAll:   false,
					ResourceCondition: []kubeteachv1.ResourceCondition{{
						Field:    "metadata.name",
						Operator: "contains",
						Value:    "contains",
					},
					},
				},
				},
			}
			Expect(c.ApplyChecks(ctx)).Should(BeTrue())

			By("ResourceCondition if a Namespace not exists")
			c = CheckController{
				Client: k8sClient,
				Log:    nil,
				CheckItems: []kubeteachv1.TaskCondition{{
					ApiVersion: "v1",
					Kind:       "Namespace",
					MatchAll:   true,
					ResourceCondition: []kubeteachv1.ResourceCondition{{
						Field:    "metadata.name",
						Operator: "neq",
						Value:    "not-found",
					},
					},
				},
				},
			}
			Expect(c.ApplyChecks(ctx)).Should(BeTrue())

			By("ResourceCondition if a greater than")
			object = &v1.Namespace{
				ObjectMeta: metav1.ObjectMeta{Name: "kubeteach-finalizers", Finalizers: []string{"test.domain/finalizer1", "test.domain/finalizer2", "test.domain/finalizer3"}},
				Spec:       v1.NamespaceSpec{},
			}
			Expect(k8sClient.Create(ctx, object)).Should(Succeed())
			c = CheckController{
				Client: k8sClient,
				Log:    nil,
				CheckItems: []kubeteachv1.TaskCondition{{
					ApiVersion: "v1",
					Kind:       "Namespace",
					MatchAll:   false,
					ResourceCondition: []kubeteachv1.ResourceCondition{{
						Field:    "metadata.finalizers.#",
						Operator: "gt",
						Value:    "2",
					}, {
						Field:    "metadata.name",
						Operator: "eq",
						Value:    "kubeteach-finalizers",
					},
					},
				},
				},
			}
			Expect(c.ApplyChecks(ctx)).Should(BeTrue())

			By("ResourceCondition if a less than")
			c = CheckController{
				Client: k8sClient,
				Log:    nil,
				CheckItems: []kubeteachv1.TaskCondition{{
					ApiVersion: "v1",
					Kind:       "Namespace",
					MatchAll:   false,
					ResourceCondition: []kubeteachv1.ResourceCondition{{
						Field:    "metadata.finalizers.#",
						Operator: "lt",
						Value:    "5",
					}, {
						Field:    "metadata.name",
						Operator: "eq",
						Value:    "kubeteach-finalizers",
					},
					},
				},
				},
			}
			Expect(c.ApplyChecks(ctx)).Should(BeTrue())

			By("ResourceCondition if a less than no integer - fail")
			c = CheckController{
				Client: k8sClient,
				Log:    nil,
				CheckItems: []kubeteachv1.TaskCondition{{
					ApiVersion: "v1",
					Kind:       "Namespace",
					MatchAll:   false,
					ResourceCondition: []kubeteachv1.ResourceCondition{{
						Field:    "metadata.finalizers.#",
						Operator: "lt",
						Value:    "asdf",
					},
					},
				},
				},
			}
			_, err := c.ApplyChecks(ctx)
			Expect(err).ShouldNot(Succeed())

			By("ResourceCondition if a greater than no integer - fail")
			c = CheckController{
				Client: k8sClient,
				Log:    nil,
				CheckItems: []kubeteachv1.TaskCondition{{
					ApiVersion: "v1",
					Kind:       "Namespace",
					MatchAll:   false,
					ResourceCondition: []kubeteachv1.ResourceCondition{{
						Field:    "metadata.finalizers.#",
						Operator: "gt",
						Value:    "asdf",
					},
					},
				},
				},
			}
			_, err = c.ApplyChecks(ctx)
			Expect(err).ShouldNot(Succeed())

			By("ResourceCondition if a Namespace not exists - check fail")
			c = CheckController{
				Client: k8sClient,
				Log:    nil,
				CheckItems: []kubeteachv1.TaskCondition{{
					ApiVersion: "v1",
					Kind:       "Namespace",
					MatchAll:   false,
					ResourceCondition: []kubeteachv1.ResourceCondition{{
						Field:    "metadata.name",
						Operator: "eq",
						Value:    "not-found",
					},
					},
				},
				},
			}
			Expect(c.ApplyChecks(ctx)).Should(BeFalse())

			By("Error Invalid Kind")
			c = CheckController{
				Client: k8sClient,
				Log:    nil,
				CheckItems: []kubeteachv1.TaskCondition{{
					ApiVersion: "v1",
					Kind:       "InvalidKind",
					MatchAll:   false,
					ResourceCondition: []kubeteachv1.ResourceCondition{{
						Field:    "metadata.name",
						Operator: "eq",
						Value:    "not-found",
					},
					},
				},
				},
			}
			_, err = c.ApplyChecks(ctx)
			Expect(err).ShouldNot(Succeed())

			By("Error 0 TaskConditions")
			c = CheckController{
				Client:     k8sClient,
				Log:        nil,
				CheckItems: nil,
			}
			_, err = c.ApplyChecks(ctx)
			Expect(err).ShouldNot(Succeed())
		})
	})
})
