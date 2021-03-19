package controllers

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	teachv1alpha1 "kubeteach/api/v1alpha1"
)

type testData struct {
	state          string
	initialDeploy  runtime.Object
	solution       runtime.Object
	taskDefinition teachv1alpha1.TaskDefinition
}

var requireTask1 = "task1"
var requireTask4 = "task4-require"

var tests = []testData{
	{
		solution:      &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "test1"}},
		initialDeploy: nil,
		state:         stateActive,
		taskDefinition: teachv1alpha1.TaskDefinition{
			ObjectMeta: metav1.ObjectMeta{Name: "task1", Namespace: "default"},
			Spec: teachv1alpha1.TaskDefinitionSpec{
				TaskSpec: teachv1alpha1.TaskSpec{
					Title:       "task1",
					Description: "Task1 description",
					HelpURL:     "HelpURL",
				},
				TaskConditions: []teachv1alpha1.TaskCondition{{
					APIVersion: "v1",
					Kind:       "Namespace",
					APIGroup:   "",
					Name:       "test1",
					ResourceCondition: []teachv1alpha1.ResourceCondition{{
						Field:    "metadata.name",
						Operator: "eq",
						Value:    "test1",
					},
					},
				}},
				RequiredTaskName: nil,
			},
		},
	}, {
		state:         stateSuccessful,
		solution:      nil,
		initialDeploy: &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "test2"}},
		taskDefinition: teachv1alpha1.TaskDefinition{
			ObjectMeta: metav1.ObjectMeta{Name: "task2", Namespace: "default"},
			Spec: teachv1alpha1.TaskDefinitionSpec{
				TaskSpec: teachv1alpha1.TaskSpec{
					Title:       "task2",
					Description: "Task2 description",
					HelpURL:     "HelpURL",
				},
				TaskConditions: []teachv1alpha1.TaskCondition{{
					APIVersion: "v1",
					Kind:       "Namespace",
					APIGroup:   "",
					Name:       "test2",
					ResourceCondition: []teachv1alpha1.ResourceCondition{{
						Field:    "metadata.name",
						Operator: "eq",
						Value:    "test2",
					},
					},
				}},
				RequiredTaskName: nil,
			},
		},
	}, {
		state:         statePending,
		solution:      &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "test3"}},
		initialDeploy: nil,
		taskDefinition: teachv1alpha1.TaskDefinition{
			ObjectMeta: metav1.ObjectMeta{Name: "task3", Namespace: "default"},
			Spec: teachv1alpha1.TaskDefinitionSpec{
				TaskSpec: teachv1alpha1.TaskSpec{
					Title:       "task3",
					Description: "Task3 description",
					HelpURL:     "HelpURL",
				},
				TaskConditions: []teachv1alpha1.TaskCondition{{
					APIVersion: "v1",
					Kind:       "Namespace",
					APIGroup:   "",
					Name:       "test3",
					ResourceCondition: []teachv1alpha1.ResourceCondition{{
						Field:    "metadata.name",
						Operator: "eq",
						Value:    "test3",
					},
					},
				}},
				RequiredTaskName: &requireTask1,
			},
		},
	}, {
		state: statePending,
		solution: &teachv1alpha1.TaskDefinition{
			ObjectMeta: metav1.ObjectMeta{Name: "task4-require", Namespace: "default"},
			Spec: teachv1alpha1.TaskDefinitionSpec{
				TaskSpec: teachv1alpha1.TaskSpec{
					Title:       "task4-require",
					Description: "task4-require description",
					HelpURL:     "HelpURL",
				},
				TaskConditions: []teachv1alpha1.TaskCondition{{
					APIVersion: "v1",
					Kind:       "Namespace",
					APIGroup:   "",
					Name:       "default",
					ResourceCondition: []teachv1alpha1.ResourceCondition{{
						Field:    "metadata.name",
						Operator: "eq",
						Value:    "default",
					},
					},
				}},
			},
		},
		initialDeploy: nil,
		taskDefinition: teachv1alpha1.TaskDefinition{
			ObjectMeta: metav1.ObjectMeta{Name: "task4", Namespace: "default"},
			Spec: teachv1alpha1.TaskDefinitionSpec{
				TaskSpec: teachv1alpha1.TaskSpec{
					Title:       "task4",
					Description: "Task4 description",
					HelpURL:     "HelpURL",
				},
				TaskConditions: []teachv1alpha1.TaskCondition{{
					APIVersion: "v1",
					Kind:       "Namespace",
					APIGroup:   "",
					Name:       "default",
					ResourceCondition: []teachv1alpha1.ResourceCondition{{
						Field:    "metadata.name",
						Operator: "eq",
						Value:    "default",
					},
					},
				}},
				RequiredTaskName: &requireTask4,
			},
		},
	},
}
