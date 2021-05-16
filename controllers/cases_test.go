package controllers

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	teachv1alpha1 "github.com/dergeberl/kubeteach/api/v1alpha1"
)

type testDataTaskDefinition struct {
	state          string
	initialDeploy  client.Object
	solution       client.Object
	taskDefinition teachv1alpha1.TaskDefinition
}

type testDataExerciseSet struct {
	status        teachv1alpha1.ExerciseSetStatus
	initialDeploy []client.Object
	exerciseSet   teachv1alpha1.ExerciseSet
}

var requireTask1 = "task1"
var requireTask4 = "task4-require"

var testsTaskDefinition = []testDataTaskDefinition{
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
var requiredTaskNameExerciseSet1 = "exerciseset1-3"
var testsExerciseSet = testDataExerciseSet{
	status: teachv1alpha1.ExerciseSetStatus{
		NumberOfTasks:              6,
		NumberOfActiveTasks:        3,
		NumberOfPendingTasks:       1,
		NumberOfSuccessfulTasks:    2,
		NumberOfUnknownTasks:       0,
		NumberOfTasksWithoutPoints: 2,
		PointsTotal:                10,
		PointsAchieved:             3,
	},
	initialDeploy: []client.Object{
		&v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "exerciseset1-1"}},
		&v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "exerciseset1-2"}},
	},
	exerciseSet: teachv1alpha1.ExerciseSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "exerciseset1",
			Namespace: "default",
		},
		Spec: teachv1alpha1.ExerciseSetSpec{
			TaskDefinitions: []teachv1alpha1.ExerciseSetSpecTaskDefinitions{
				{
					Name: "exerciseset1-1",
					TaskDefinitionSpec: teachv1alpha1.TaskDefinitionSpec{
						TaskSpec: teachv1alpha1.TaskSpec{
							Title:       "exerciseset1-1",
							Description: "exerciseset1-1",
						},
						TaskConditions: []teachv1alpha1.TaskCondition{{
							APIVersion: "v1",
							Kind:       "Namespace",
							APIGroup:   "",
							Name:       "exerciseset1-1",
							ResourceCondition: []teachv1alpha1.ResourceCondition{{
								Field:    "metadata.name",
								Operator: "eq",
								Value:    "exerciseset1-1",
							},
							},
						}},
						Points: 1,
					},
				}, {
					Name: "exerciseset1-2",
					TaskDefinitionSpec: teachv1alpha1.TaskDefinitionSpec{
						TaskSpec: teachv1alpha1.TaskSpec{
							Title:       "exerciseset1-2",
							Description: "exerciseset1-2",
						},
						TaskConditions: []teachv1alpha1.TaskCondition{{
							APIVersion: "v1",
							Kind:       "Namespace",
							APIGroup:   "",
							Name:       "exerciseset1-2",
							ResourceCondition: []teachv1alpha1.ResourceCondition{{
								Field:    "metadata.name",
								Operator: "eq",
								Value:    "exerciseset1-2",
							},
							},
						}},
						Points: 2,
					},
				}, {
					Name: "exerciseset1-3",
					TaskDefinitionSpec: teachv1alpha1.TaskDefinitionSpec{
						TaskSpec: teachv1alpha1.TaskSpec{
							Title:       "exerciseset1-3",
							Description: "exerciseset1-3",
						},
						TaskConditions: []teachv1alpha1.TaskCondition{{
							APIVersion: "v1",
							Kind:       "Namespace",
							APIGroup:   "",
							Name:       "exerciseset1-3",
							ResourceCondition: []teachv1alpha1.ResourceCondition{{
								Field:    "metadata.name",
								Operator: "eq",
								Value:    "exerciseset1-3",
							},
							},
						}},
						Points: 3,
					},
				}, {
					Name: "exerciseset1-4",
					TaskDefinitionSpec: teachv1alpha1.TaskDefinitionSpec{
						TaskSpec: teachv1alpha1.TaskSpec{
							Title:       "exerciseset1-4",
							Description: "exerciseset1-4",
						},
						TaskConditions: []teachv1alpha1.TaskCondition{{
							APIVersion: "v1",
							Kind:       "Namespace",
							APIGroup:   "",
							Name:       "exerciseset1-4",
							ResourceCondition: []teachv1alpha1.ResourceCondition{{
								Field:    "metadata.name",
								Operator: "eq",
								Value:    "exerciseset1-4",
							},
							},
						}},
						Points:           4,
						RequiredTaskName: &requiredTaskNameExerciseSet1,
					},
				}, {
					Name: "exerciseset1-5",
					TaskDefinitionSpec: teachv1alpha1.TaskDefinitionSpec{
						TaskSpec: teachv1alpha1.TaskSpec{
							Title:       "exerciseset1-5",
							Description: "exerciseset1-5",
						},
						TaskConditions: []teachv1alpha1.TaskCondition{{
							APIVersion: "v1",
							Kind:       "Namespace",
							APIGroup:   "",
							Name:       "exerciseset1-5",
							ResourceCondition: []teachv1alpha1.ResourceCondition{{
								Field:    "metadata.name",
								Operator: "eq",
								Value:    "exerciseset1-5",
							},
							},
						}},
						Points: 0,
					},
				}, {
					Name: "exerciseset1-6",
					TaskDefinitionSpec: teachv1alpha1.TaskDefinitionSpec{
						TaskSpec: teachv1alpha1.TaskSpec{
							Title:       "exerciseset1-6",
							Description: "exerciseset1-6",
						},
						TaskConditions: []teachv1alpha1.TaskCondition{{
							APIVersion: "v1",
							Kind:       "Namespace",
							APIGroup:   "",
							Name:       "exerciseset1-6",
							ResourceCondition: []teachv1alpha1.ResourceCondition{{
								Field:    "metadata.name",
								Operator: "eq",
								Value:    "exerciseset1-6",
							},
							},
						}},
						Points: 0,
					},
				},
			},
		},
	},
}
