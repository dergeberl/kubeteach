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

package metrics

import (
	teachv1alpha1 "github.com/dergeberl/kubeteach/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type testDataExerciseSet struct {
	status      teachv1alpha1.ExerciseSetStatus
	exerciseSet teachv1alpha1.ExerciseSet
}

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
							APIVersion:        "v1",
							Kind:              "Namespace",
							APIGroup:          "",
							Name:              "exerciseset1-1",
							ResourceCondition: []teachv1alpha1.ResourceCondition{},
						},
						},
					},
				},
			},
		},
	},
}

var testTasks1 = teachv1alpha1.Task{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "task1",
		Namespace: "default",
	},
	Spec: teachv1alpha1.TaskSpec{
		Title:       "task1",
		Description: "task1",
	},
}
var testTasks2 = teachv1alpha1.Task{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "task2",
		Namespace: "default",
	},
	Spec: teachv1alpha1.TaskSpec{
		Title:       "task2",
		Description: "task2",
	},
}
var testTasks3 = teachv1alpha1.Task{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "task3",
		Namespace: "default",
	},
	Spec: teachv1alpha1.TaskSpec{
		Title:       "task3",
		Description: "task3",
	},
}
var testTasks4 = teachv1alpha1.Task{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "task4",
		Namespace: "default",
	},
	Spec: teachv1alpha1.TaskSpec{
		Title:       "task4",
		Description: "task4",
	},
}
