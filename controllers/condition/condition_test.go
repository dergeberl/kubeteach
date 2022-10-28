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

package condition

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TaskConditions ApplyChecks", func() {
	Context("Run checks in checkItems", func() {
		It("run testcases", func() {
			ctx := context.Background()
			for _, test := range testCases {
				By(test.name)
				if test.obj != nil {
					for _, obj := range test.obj {
						Expect(k8sClient.Create(ctx, obj)).Should(Succeed())
					}
				}
				c := Checks{Client: k8sClient}
				got, gotErr := c.ApplyChecks(ctx, test.taskCondition)
				Expect(got).Should(test.state)
				Expect(gotErr).Should(test.err)
				if test.obj != nil {
					for _, obj := range test.obj {
						_ = k8sClient.Delete(ctx, obj)
					}
				}
			}
		})
	})
})
