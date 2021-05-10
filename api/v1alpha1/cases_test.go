package v1alpha1

import (
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type testCases struct {
	obj runtime.Object
	err types.GomegaMatcher
}

var taskCases = []testCases{
	{
		obj: &Task{
			ObjectMeta: metav1.ObjectMeta{Name: "task-valid1", Namespace: "default"},
			Spec: TaskSpec{
				Title:           "task",
				Description:     "task",
				LongDescription: "task",
				HelpURL:         "task",
			},
		},
		err: BeNil(),
	}, {
		obj: &Task{
			ObjectMeta: metav1.ObjectMeta{Name: "task-valid2", Namespace: "default"},
			Spec: TaskSpec{
				Title:       "task",
				Description: "task",
			},
		},
		err: BeNil(),
	}, {
		obj: &Task{
			ObjectMeta: metav1.ObjectMeta{Name: "task-invalid1", Namespace: "default"},
			Spec: TaskSpec{
				Description: "task",
			},
		},
		err: Not(BeNil()),
	}, {
		obj: &Task{
			ObjectMeta: metav1.ObjectMeta{Name: "task-invalid2", Namespace: "default"},
			Spec: TaskSpec{
				Title: "task",
			},
		},
		err: Not(BeNil()),
	}, {
		obj: &Task{
			ObjectMeta: metav1.ObjectMeta{Name: "task-invalid2", Namespace: "default"},
			Spec:       TaskSpec{},
		},
		err: Not(BeNil()),
	},
}

var taskDefinitionCases = []testCases{
	{
		obj: &TaskDefinition{
			ObjectMeta: metav1.ObjectMeta{Name: "valid1", Namespace: "default"},
			Spec: TaskDefinitionSpec{
				TaskSpec: TaskSpec{
					Title:       "Test1",
					Description: "Test1",
				},
				RequiredTaskName: nil,
				TaskConditions: []TaskCondition{
					{
						APIVersion:        "v1",
						Kind:              "Namespace",
						Name:              "default",
						ResourceCondition: nil,
					},
				},
			}},
		err: BeNil(),
	}, {
		obj: &TaskDefinition{
			ObjectMeta: metav1.ObjectMeta{Name: "valid2-all-operators", Namespace: "default"},
			Spec: TaskDefinitionSpec{
				TaskSpec: TaskSpec{
					Title:           "Test1",
					Description:     "Test1",
					LongDescription: "Test1",
					HelpURL:         "Test1",
				},
				RequiredTaskName: nil,
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
							}, {
								Field:    "meta.name",
								Operator: "neq",
								Value:    "name",
							}, {
								Field:    "meta.name",
								Operator: "lt",
								Value:    "name",
							}, {
								Field:    "meta.name",
								Operator: "gt",
								Value:    "name",
							}, {
								Field:    "meta.name",
								Operator: "contains",
								Value:    "name",
							}, {
								Field:    "meta.name",
								Operator: "nil",
								Value:    "name",
							}, {
								Field:    "meta.name",
								Operator: "notnil",
								Value:    "name",
							},
						},
					},
				},
			}},
		err: BeNil(),
	}, {
		obj: &TaskDefinition{
			ObjectMeta: metav1.ObjectMeta{Name: "invalid1-wrong-operator", Namespace: "default"},
			Spec: TaskDefinitionSpec{
				TaskSpec: TaskSpec{
					Title:           "Test1",
					Description:     "Test1",
					LongDescription: "Test1",
					HelpURL:         "Test1",
				},
				RequiredTaskName: nil,
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
								Operator: "error",
								Value:    "name",
							},
						},
					},
				},
			}},
		err: Not(BeNil()),
	}, {
		obj: &TaskDefinition{
			ObjectMeta: metav1.ObjectMeta{Name: "invalid2-no-operator", Namespace: "default"},
			Spec: TaskDefinitionSpec{
				TaskSpec: TaskSpec{
					Title:           "Test",
					Description:     "Test",
					LongDescription: "Test",
					HelpURL:         "Test",
				},
				RequiredTaskName: nil,
				TaskConditions: []TaskCondition{
					{
						APIVersion: "v1",
						Kind:       "Namespace",
						Name:       "default",
						Namespace:  "",
						NotExists:  false,
						ResourceCondition: []ResourceCondition{
							{
								Field: "meta.name",
								Value: "name",
							},
						},
					},
				},
			}},
		err: Not(BeNil()),
	}, {
		obj: &TaskDefinition{
			ObjectMeta: metav1.ObjectMeta{Name: "invalid3-no-field", Namespace: "default"},
			Spec: TaskDefinitionSpec{
				TaskSpec: TaskSpec{
					Title:           "Test",
					Description:     "Test",
					LongDescription: "Test",
					HelpURL:         "Test",
				},
				RequiredTaskName: nil,
				TaskConditions: []TaskCondition{
					{
						APIVersion: "v1",
						Kind:       "Namespace",
						Name:       "default",
						Namespace:  "",
						NotExists:  false,
						ResourceCondition: []ResourceCondition{
							{
								Operator: "eq",
								Value:    "name",
							},
						},
					},
				},
			}},
		err: Not(BeNil()),
	}, {
		obj: &TaskDefinition{
			ObjectMeta: metav1.ObjectMeta{Name: "invalid4-no-kind", Namespace: "default"},
			Spec: TaskDefinitionSpec{
				TaskSpec: TaskSpec{
					Title:           "Test",
					Description:     "Test",
					LongDescription: "Test",
					HelpURL:         "Test",
				},
				RequiredTaskName: nil,
				TaskConditions: []TaskCondition{
					{
						APIVersion: "v1",
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
			}},
		err: Not(BeNil()),
	}, {
		obj: &TaskDefinition{
			ObjectMeta: metav1.ObjectMeta{Name: "invalid5-no-name", Namespace: "default"},
			Spec: TaskDefinitionSpec{
				TaskSpec: TaskSpec{
					Title:           "Test",
					Description:     "Test",
					LongDescription: "Test",
					HelpURL:         "Test",
				},
				RequiredTaskName: nil,
				TaskConditions: []TaskCondition{
					{
						APIVersion: "v1",
						Kind:       "Namespace",
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
			}},
		err: Not(BeNil()),
	}, {
		obj: &TaskDefinition{
			ObjectMeta: metav1.ObjectMeta{Name: "invalid6-no-apiversion", Namespace: "default"},
			Spec: TaskDefinitionSpec{
				TaskSpec: TaskSpec{
					Title:           "Test",
					Description:     "Test",
					LongDescription: "Test",
					HelpURL:         "Test",
				},
				RequiredTaskName: nil,
				TaskConditions: []TaskCondition{
					{
						Name:      "default",
						Kind:      "Namespace",
						Namespace: "",
						NotExists: false,
						ResourceCondition: []ResourceCondition{
							{
								Field:    "meta.name",
								Operator: "eq",
								Value:    "name",
							},
						},
					},
				},
			}},
		err: Not(BeNil()),
	}, {
		obj: &TaskDefinition{
			ObjectMeta: metav1.ObjectMeta{Name: "invalid7-nil-taskconditions", Namespace: "default"},
			Spec: TaskDefinitionSpec{
				TaskSpec: TaskSpec{
					Title:           "Test",
					Description:     "Test",
					LongDescription: "Test",
					HelpURL:         "Test",
				},
				RequiredTaskName: nil,
			}},
		err: Not(BeNil()),
	}, {
		obj: &TaskDefinition{
			ObjectMeta: metav1.ObjectMeta{Name: "invalid8-no-taskconditions", Namespace: "default"},
			Spec: TaskDefinitionSpec{
				TaskSpec: TaskSpec{
					Title:           "Test",
					Description:     "Test",
					LongDescription: "Test",
					HelpURL:         "Test",
				},
				RequiredTaskName: nil,
				TaskConditions:   []TaskCondition{},
			}},
		err: Not(BeNil()),
	}, {
		obj: &TaskDefinition{
			ObjectMeta: metav1.ObjectMeta{Name: "invalid9-no-taskspec", Namespace: "default"},
			Spec: TaskDefinitionSpec{
				RequiredTaskName: nil,
				TaskConditions: []TaskCondition{
					{
						APIVersion:        "v1",
						Kind:              "Namespace",
						Name:              "default",
						ResourceCondition: nil,
					},
				},
			}},
		err: Not(BeNil()),
	}, {
		obj: &TaskDefinition{
			ObjectMeta: metav1.ObjectMeta{Name: "invalid10-invalid-taskspec", Namespace: "default"},
			Spec: TaskDefinitionSpec{
				TaskSpec: TaskSpec{
					Title: "Test1",
				},
				RequiredTaskName: nil,
				TaskConditions: []TaskCondition{
					{
						APIVersion:        "v1",
						Kind:              "Namespace",
						Name:              "default",
						ResourceCondition: nil,
					},
				},
			}},
		err: Not(BeNil()),
	},
}
