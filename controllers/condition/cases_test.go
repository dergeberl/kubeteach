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
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	teachv1alpha1 "github.com/dergeberl/kubeteach/api/v1alpha1"
)

type conditionTest struct {
	name          string
	obj           []runtime.Object
	taskCondition []teachv1alpha1.TaskCondition
	state         types.GomegaMatcher
	err           types.GomegaMatcher
}

//nolint:lll
var testCases = []conditionTest{
	{
		name:          "error - invalid objectType",
		obj:           nil,
		taskCondition: []teachv1alpha1.TaskCondition{{APIVersion: "", Kind: "", APIGroup: "", ResourceCondition: nil}},
		state:         BeFalse(),
		err:           Not(BeNil()),
	}, {
		name:          "error - no TaskCondition",
		obj:           nil,
		taskCondition: []teachv1alpha1.TaskCondition{},
		state:         BeFalse(),
		err:           Not(BeNil()),
	}, {
		name: "error - invalid operator",
		obj: []runtime.Object{
			&v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "error-no-operator"}},
		},
		taskCondition: []teachv1alpha1.TaskCondition{{APIVersion: "v1", Kind: "Namespace", Name: "error-no-operator", APIGroup: "", ResourceCondition: []teachv1alpha1.ResourceCondition{{Field: "metadata.name", Operator: "invalid", Value: ""}}}},
		state:         BeFalse(),
		err:           Not(BeNil()),
	}, {
		name: "true - no ResourceCondition set but object exists",
		obj: []runtime.Object{
			&v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "no-resource-condition"}},
		},
		taskCondition: []teachv1alpha1.TaskCondition{{APIVersion: "v1", Kind: "Namespace", APIGroup: "", Name: "no-resource-condition", ResourceCondition: nil}},
		state:         BeTrue(),
		err:           BeNil(),
	}, {
		name:          "false - object-not-found",
		obj:           nil,
		taskCondition: []teachv1alpha1.TaskCondition{{APIVersion: "v1", Kind: "Namespace", APIGroup: "", Name: "object-not-found", NotExists: false, ResourceCondition: nil}},
		state:         BeFalse(),
		err:           BeNil(),
	}, {
		name:          "true - notExists",
		obj:           nil,
		taskCondition: []teachv1alpha1.TaskCondition{{APIVersion: "v1", Kind: "Namespace", APIGroup: "", Name: "test1-notexists", NotExists: true, ResourceCondition: nil}},
		state:         BeTrue(),
		err:           BeNil(),
	}, {
		name: "false - notExists",
		obj: []runtime.Object{
			&v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "test2-notexists"}},
		},
		taskCondition: []teachv1alpha1.TaskCondition{{APIVersion: "v1", Kind: "Namespace", APIGroup: "", Name: "test2-notexists", NotExists: true, ResourceCondition: nil}},
		state:         BeFalse(),
		err:           BeNil(),
	}, {
		name: "true - test eq",
		obj: []runtime.Object{
			&v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "test1-eq"}},
		},
		taskCondition: []teachv1alpha1.TaskCondition{{APIVersion: "v1", Kind: "Namespace", APIGroup: "", Name: "test1-eq", ResourceCondition: []teachv1alpha1.ResourceCondition{{Field: "metadata.name", Operator: "eq", Value: "test1-eq"}}}},
		state:         BeTrue(),
		err:           BeNil(),
	}, {
		name:          "false - test eq",
		obj:           nil,
		taskCondition: []teachv1alpha1.TaskCondition{{APIVersion: "v1", Kind: "Namespace", APIGroup: "", Name: "test1-eq", ResourceCondition: []teachv1alpha1.ResourceCondition{{Field: "metadata.name", Operator: "eq", Value: "test2-eq"}}}},
		state:         BeFalse(),
		err:           BeNil(),
	}, {
		name: "true - test neq",
		obj: []runtime.Object{
			&v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "test1-neq"}},
		},
		taskCondition: []teachv1alpha1.TaskCondition{{APIVersion: "v1", Kind: "Namespace", APIGroup: "", Name: "test1-neq", ResourceCondition: []teachv1alpha1.ResourceCondition{{Field: "metadata.name", Operator: "neq", Value: "neq"}}}},
		state:         BeTrue(),
		err:           BeNil(),
	}, {
		name: "false - test neq",
		obj: []runtime.Object{
			&v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "test2-neq"}},
		},
		taskCondition: []teachv1alpha1.TaskCondition{{APIVersion: "v1", Kind: "Namespace", APIGroup: "", Name: "test2-neq", ResourceCondition: []teachv1alpha1.ResourceCondition{{Field: "metadata.name", Operator: "neq", Value: "test2-neq"}}}},
		state:         BeFalse(),
		err:           BeNil(),
	}, {
		name: "true - test lt",
		obj: []runtime.Object{
			&v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "test1-lt", Finalizers: []string{"test.domain/finalizer1", "test.domain/finalizer2", "test.domain/finalizer3"}}},
		},
		taskCondition: []teachv1alpha1.TaskCondition{{APIVersion: "v1", Kind: "Namespace", APIGroup: "", Name: "test1-lt", ResourceCondition: []teachv1alpha1.ResourceCondition{{Field: "metadata.finalizers.#", Operator: "lt", Value: "5"}}}},
		state:         BeTrue(),
		err:           BeNil(),
	}, {
		name: "false - test lt",
		obj: []runtime.Object{
			&v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "test2-lt", Finalizers: []string{"test.domain/finalizer1", "test.domain/finalizer2", "test.domain/finalizer3"}}},
		},
		taskCondition: []teachv1alpha1.TaskCondition{{APIVersion: "v1", Kind: "Namespace", APIGroup: "", Name: "test2-lt", ResourceCondition: []teachv1alpha1.ResourceCondition{{Field: "metadata.finalizers.#", Operator: "lt", Value: "1"}}}},
		state:         BeFalse(),
		err:           BeNil(),
	}, {
		name: "error - test lt no int in value",
		obj: []runtime.Object{
			&v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "test3-lt", Finalizers: []string{"test.domain/finalizer1", "test.domain/finalizer2", "test.domain/finalizer3"}}},
		},
		taskCondition: []teachv1alpha1.TaskCondition{{APIVersion: "v1", Kind: "Namespace", APIGroup: "", Name: "test3-lt", ResourceCondition: []teachv1alpha1.ResourceCondition{{Field: "", Operator: "lt", Value: "noInt"}}}},
		state:         BeFalse(),
		err:           Not(BeNil()),
	}, {
		name: "true - test gt",
		obj: []runtime.Object{
			&v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "test1-gt", Finalizers: []string{"test.domain/finalizer1", "test.domain/finalizer2", "test.domain/finalizer3"}}},
		},
		taskCondition: []teachv1alpha1.TaskCondition{{APIVersion: "v1", Kind: "Namespace", APIGroup: "", Name: "test1-gt", ResourceCondition: []teachv1alpha1.ResourceCondition{{Field: "metadata.finalizers.#", Operator: "gt", Value: "1"}}}},
		state:         BeTrue(),
		err:           BeNil(),
	}, {
		name: "false - test gt",
		obj: []runtime.Object{
			&v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "test2-gt", Finalizers: []string{"test.domain/finalizer1", "test.domain/finalizer2", "test.domain/finalizer3"}}},
		},
		taskCondition: []teachv1alpha1.TaskCondition{{APIVersion: "v1", Kind: "Namespace", APIGroup: "", Name: "test2-gt", ResourceCondition: []teachv1alpha1.ResourceCondition{{Field: "metadata.finalizers.#", Operator: "gt", Value: "5"}}}},
		state:         BeFalse(),
		err:           BeNil(),
	}, {
		name: "error - test gt no int in value",
		obj: []runtime.Object{
			&v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "test3-gt", Finalizers: []string{"test.domain/finalizer1", "test.domain/finalizer2", "test.domain/finalizer3"}}},
		},
		taskCondition: []teachv1alpha1.TaskCondition{{APIVersion: "v1", Kind: "Namespace", APIGroup: "", Name: "test3-gt", ResourceCondition: []teachv1alpha1.ResourceCondition{{Field: "", Operator: "gt", Value: "noInt"}}}},
		state:         BeFalse(),
		err:           Not(BeNil()),
	}, {
		name: "true - test contains",
		obj: []runtime.Object{
			&v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "test1-contains1"}},
		},
		taskCondition: []teachv1alpha1.TaskCondition{{APIVersion: "v1", Kind: "Namespace", APIGroup: "", Name: "test1-contains1", ResourceCondition: []teachv1alpha1.ResourceCondition{{Field: "metadata.name", Operator: "contains", Value: "contains1"}}}},
		state:         BeTrue(),
		err:           BeNil(),
	}, {
		name: "false - test contains",
		obj: []runtime.Object{
			&v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "test2-contains1"}},
		},
		taskCondition: []teachv1alpha1.TaskCondition{{APIVersion: "v1", Kind: "Namespace", APIGroup: "", Name: "test2-contains1", ResourceCondition: []teachv1alpha1.ResourceCondition{{Field: "metadata.name", Operator: "contains", Value: "contains2"}}}},
		state:         BeFalse(),
		err:           BeNil(),
	}, {
		name: "true - test nil",
		obj: []runtime.Object{
			&v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "test1-nil"}},
		},
		taskCondition: []teachv1alpha1.TaskCondition{{APIVersion: "v1", Kind: "Namespace", APIGroup: "", Name: "test1-nil", ResourceCondition: []teachv1alpha1.ResourceCondition{{Field: "metadata.deletionTimestamp", Operator: "nil"}}}},
		state:         BeTrue(),
		err:           BeNil(),
	}, {
		name: "false - test nil",
		obj: []runtime.Object{
			&v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "test2-nil"}},
		},
		taskCondition: []teachv1alpha1.TaskCondition{{APIVersion: "v1", Kind: "Namespace", APIGroup: "", Name: "test2-nil", ResourceCondition: []teachv1alpha1.ResourceCondition{{Field: "metadata.name", Operator: "nil"}}}},
		state:         BeFalse(),
		err:           BeNil(),
	}, {
		name: "true - test notnil",
		obj: []runtime.Object{
			&v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "test1-notnil"}},
		},
		taskCondition: []teachv1alpha1.TaskCondition{{APIVersion: "v1", Kind: "Namespace", APIGroup: "", Name: "test1-notnil", ResourceCondition: []teachv1alpha1.ResourceCondition{{Field: "metadata.name", Operator: "notnil"}}}},
		state:         BeTrue(),
		err:           BeNil(),
	}, {
		name: "false - test notnil",
		obj: []runtime.Object{
			&v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "test2-notnil"}},
		},
		taskCondition: []teachv1alpha1.TaskCondition{{APIVersion: "v1", Kind: "Namespace", APIGroup: "", Name: "test2-notnil", ResourceCondition: []teachv1alpha1.ResourceCondition{{Field: "metadata.deletionTimestamp", Operator: "notnil"}}}},
		state:         BeFalse(),
		err:           BeNil(),
	}, {
		name: "false - test multi conditions",
		obj: []runtime.Object{
			&v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "test1-multiconditions1"}},
		},
		taskCondition: []teachv1alpha1.TaskCondition{
			{APIVersion: "v1", Kind: "Namespace", APIGroup: "", Name: "test1-multiconditions1", ResourceCondition: []teachv1alpha1.ResourceCondition{{Field: "metadata.name", Operator: "eq", Value: "test1-multiconditions1"}}},
			{APIVersion: "v1", Kind: "Namespace", APIGroup: "", Name: "test1-multiconditions2", ResourceCondition: []teachv1alpha1.ResourceCondition{{Field: "metadata.name", Operator: "eq", Value: "test1-multiconditions2"}}},
		},
		state: BeFalse(),
		err:   BeNil(),
	}, {
		name: "true - test multi conditions",
		obj: []runtime.Object{
			&v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "test2-multiconditions1"}},
			&v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "test2-multiconditions2"}},
		},
		taskCondition: []teachv1alpha1.TaskCondition{
			{APIVersion: "v1", Kind: "Namespace", APIGroup: "", Name: "test2-multiconditions1", ResourceCondition: []teachv1alpha1.ResourceCondition{{Field: "metadata.name", Operator: "eq", Value: "test2-multiconditions1"}}},
			{APIVersion: "v1", Kind: "Namespace", APIGroup: "", Name: "test2-multiconditions2", ResourceCondition: []teachv1alpha1.ResourceCondition{{Field: "metadata.name", Operator: "eq", Value: "test2-multiconditions2"}}},
		},
		state: BeTrue(),
		err:   BeNil(),
	}, {
		name: "false - error - test multi conditions",
		obj: []runtime.Object{
			&v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "test3-multiconditions1"}},
			&v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "test3-multiconditions2"}},
		},
		taskCondition: []teachv1alpha1.TaskCondition{
			{APIVersion: "v1", Kind: "Namespace", APIGroup: "", Name: "test3-multiconditions1", ResourceCondition: []teachv1alpha1.ResourceCondition{{Field: "metadata.name", Operator: "eq", Value: "test3-multiconditions1"}}},
			{APIVersion: "v1", Kind: "Namespace", APIGroup: "", Name: "test3-multiconditions2", ResourceCondition: []teachv1alpha1.ResourceCondition{{Field: "metadata.name", Operator: "ERROR", Value: "test3-multiconditions2"}}},
		},
		state: BeFalse(),
		err:   Not(BeNil()),
	},
}
