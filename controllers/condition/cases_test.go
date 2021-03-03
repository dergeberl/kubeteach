package condition

import (
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	teachv1alpha1 "kubeteach/api/v1alpha1"
)

type conditionTest struct {
	name          string
	obj           runtime.Object
	taskCondition []teachv1alpha1.TaskCondition
	state         types.GomegaMatcher
	err           types.GomegaMatcher
}

//nolint:lll
var testCases = []conditionTest{
	{
		name:          "error - invalid objectType",
		obj:           nil,
		taskCondition: []teachv1alpha1.TaskCondition{{APIVersion: "", Kind: "", APIGroup: "", MatchAll: false, ResourceCondition: nil}},
		state:         BeFalse(),
		err:           Not(BeNil()),
	}, {
		name:          "error - no TaskCondition",
		obj:           nil,
		taskCondition: []teachv1alpha1.TaskCondition{},
		state:         BeFalse(),
		err:           Not(BeNil()),
	}, {
		name:          "error - invalid operator",
		obj:           nil,
		taskCondition: []teachv1alpha1.TaskCondition{{APIVersion: "v1", Kind: "Namespace", APIGroup: "", MatchAll: false, ResourceCondition: []teachv1alpha1.ResourceCondition{{Field: "metadata.name", Operator: "invalid", Value: ""}}}},
		state:         BeFalse(),
		err:           Not(BeNil()),
	}, {
		name:          "error - no ResourceCondition set",
		obj:           nil,
		taskCondition: []teachv1alpha1.TaskCondition{{APIVersion: "v1", Kind: "Namespace", APIGroup: "", MatchAll: false, ResourceCondition: nil}},
		state:         BeFalse(),
		err:           Not(BeNil()),
	}, {
		name:          "true - test eq",
		obj:           &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "test1-eq"}},
		taskCondition: []teachv1alpha1.TaskCondition{{APIVersion: "v1", Kind: "Namespace", APIGroup: "", MatchAll: false, ResourceCondition: []teachv1alpha1.ResourceCondition{{Field: "metadata.name", Operator: "eq", Value: "test1-eq"}}}},
		state:         BeTrue(),
		err:           BeNil(),
	}, {
		name:          "false - test eq",
		obj:           nil,
		taskCondition: []teachv1alpha1.TaskCondition{{APIVersion: "v1", Kind: "Namespace", APIGroup: "", MatchAll: false, ResourceCondition: []teachv1alpha1.ResourceCondition{{Field: "metadata.name", Operator: "eq", Value: "test2-eq"}}}},
		state:         BeFalse(),
		err:           BeNil(),
	}, {
		name:          "true - test neq",
		obj:           nil,
		taskCondition: []teachv1alpha1.TaskCondition{{APIVersion: "v1", Kind: "Namespace", APIGroup: "", MatchAll: true, ResourceCondition: []teachv1alpha1.ResourceCondition{{Field: "metadata.name", Operator: "neq", Value: "test1-neq"}}}},
		state:         BeTrue(),
		err:           BeNil(),
	}, {
		name:          "false - test neq",
		obj:           &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "test2-neq"}},
		taskCondition: []teachv1alpha1.TaskCondition{{APIVersion: "v1", Kind: "Namespace", APIGroup: "", MatchAll: true, ResourceCondition: []teachv1alpha1.ResourceCondition{{Field: "metadata.name", Operator: "neq", Value: "test2-neq"}}}},
		state:         BeFalse(),
		err:           BeNil(),
	}, {
		name:          "true - test lt",
		obj:           &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "test1-lt", Finalizers: []string{"test.domain/finalizer1", "test.domain/finalizer2", "test.domain/finalizer3"}}},
		taskCondition: []teachv1alpha1.TaskCondition{{APIVersion: "v1", Kind: "Namespace", APIGroup: "", MatchAll: false, ResourceCondition: []teachv1alpha1.ResourceCondition{{Field: "metadata.name", Operator: "eq", Value: "test1-lt"}, {Field: "metadata.finalizers.#", Operator: "lt", Value: "5"}}}},
		state:         BeTrue(),
		err:           BeNil(),
	}, {
		name:          "false - test lt",
		obj:           &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "test2-lt", Finalizers: []string{"test.domain/finalizer1", "test.domain/finalizer2", "test.domain/finalizer3"}}},
		taskCondition: []teachv1alpha1.TaskCondition{{APIVersion: "v1", Kind: "Namespace", APIGroup: "", MatchAll: false, ResourceCondition: []teachv1alpha1.ResourceCondition{{Field: "metadata.name", Operator: "eq", Value: "test2-lt"}, {Field: "metadata.finalizers.#", Operator: "lt", Value: "1"}}}},
		state:         BeFalse(),
		err:           BeNil(),
	}, {
		name:          "error - test lt no int in value",
		obj:           nil,
		taskCondition: []teachv1alpha1.TaskCondition{{APIVersion: "v1", Kind: "Namespace", APIGroup: "", MatchAll: false, ResourceCondition: []teachv1alpha1.ResourceCondition{{Field: "", Operator: "lt", Value: "noInt"}}}},
		state:         BeFalse(),
		err:           Not(BeNil()),
	}, {
		name:          "true - test gt",
		obj:           &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "test1-gt", Finalizers: []string{"test.domain/finalizer1", "test.domain/finalizer2", "test.domain/finalizer3"}}},
		taskCondition: []teachv1alpha1.TaskCondition{{APIVersion: "v1", Kind: "Namespace", APIGroup: "", MatchAll: false, ResourceCondition: []teachv1alpha1.ResourceCondition{{Field: "metadata.name", Operator: "eq", Value: "test1-gt"}, {Field: "metadata.finalizers.#", Operator: "gt", Value: "1"}}}},
		state:         BeTrue(),
		err:           BeNil(),
	}, {
		name:          "false - test gt",
		obj:           &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "test2-gt", Finalizers: []string{"test.domain/finalizer1", "test.domain/finalizer2", "test.domain/finalizer3"}}},
		taskCondition: []teachv1alpha1.TaskCondition{{APIVersion: "v1", Kind: "Namespace", APIGroup: "", MatchAll: false, ResourceCondition: []teachv1alpha1.ResourceCondition{{Field: "metadata.name", Operator: "eq", Value: "test2-gt"}, {Field: "metadata.finalizers.#", Operator: "gt", Value: "5"}}}},
		state:         BeFalse(),
		err:           BeNil(),
	}, {
		name:          "error - test gt no int in value",
		obj:           nil,
		taskCondition: []teachv1alpha1.TaskCondition{{APIVersion: "v1", Kind: "Namespace", APIGroup: "", MatchAll: false, ResourceCondition: []teachv1alpha1.ResourceCondition{{Field: "", Operator: "gt", Value: "noInt"}}}},
		state:         BeFalse(),
		err:           Not(BeNil()),
	}, {
		name:          "true - test contains",
		obj:           &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "test1-contains1"}},
		taskCondition: []teachv1alpha1.TaskCondition{{APIVersion: "v1", Kind: "Namespace", APIGroup: "", MatchAll: false, ResourceCondition: []teachv1alpha1.ResourceCondition{{Field: "metadata.name", Operator: "contains", Value: "contains1"}}}},
		state:         BeTrue(),
		err:           BeNil(),
	}, {
		name:          "false - test contains",
		obj:           nil,
		taskCondition: []teachv1alpha1.TaskCondition{{APIVersion: "v1", Kind: "Namespace", APIGroup: "", MatchAll: false, ResourceCondition: []teachv1alpha1.ResourceCondition{{Field: "metadata.name", Operator: "contains", Value: "contains2"}}}},
		state:         BeFalse(),
		err:           BeNil(),
	}, {
		name:          "true - test nil",
		obj:           &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "test1-nil"}},
		taskCondition: []teachv1alpha1.TaskCondition{{APIVersion: "v1", Kind: "Namespace", APIGroup: "", MatchAll: false, ResourceCondition: []teachv1alpha1.ResourceCondition{{Field: "metadata.name", Operator: "eq", Value: "test1-nil"}, {Field: "metadata.deletionTimestamp", Operator: "nil"}}}},
		state:         BeTrue(),
		err:           BeNil(),
	}, {
		name:          "false - test nil",
		obj:           &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "test2-nil"}},
		taskCondition: []teachv1alpha1.TaskCondition{{APIVersion: "v1", Kind: "Namespace", APIGroup: "", MatchAll: false, ResourceCondition: []teachv1alpha1.ResourceCondition{{Field: "metadata.name", Operator: "eq", Value: "test2-nil"}, {Field: "metadata.name", Operator: "nil"}}}},
		state:         BeFalse(),
		err:           BeNil(),
	}, {
		name:          "true - test notnil",
		obj:           &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "test1-notnil"}},
		taskCondition: []teachv1alpha1.TaskCondition{{APIVersion: "v1", Kind: "Namespace", APIGroup: "", MatchAll: false, ResourceCondition: []teachv1alpha1.ResourceCondition{{Field: "metadata.name", Operator: "eq", Value: "test1-notnil"}, {Field: "metadata.name", Operator: "notnil"}}}},
		state:         BeTrue(),
		err:           BeNil(),
	}, {
		name:          "false - test notnil",
		obj:           &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "test2-notnil"}},
		taskCondition: []teachv1alpha1.TaskCondition{{APIVersion: "v1", Kind: "Namespace", APIGroup: "", MatchAll: false, ResourceCondition: []teachv1alpha1.ResourceCondition{{Field: "metadata.name", Operator: "eq", Value: "test2-notnil"}, {Field: "metadata.deletionTimestamp", Operator: "notnil"}}}},
		state:         BeFalse(),
		err:           BeNil(),
	},
}
