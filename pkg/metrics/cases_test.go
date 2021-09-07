package metrics

import (
	"io"
	"strings"

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

var expectExerciseSet = map[string]io.Reader{
	"kubeteach_exerciseset_active_tasks": strings.NewReader(`
# HELP kubeteach_exerciseset_active_tasks active points of exerciseset
# TYPE kubeteach_exerciseset_active_tasks gauge
kubeteach_exerciseset_active_tasks{name="exerciseset1",namespace="default"} 3
`),
	"kubeteach_exerciseset_pending_tasks": strings.NewReader(`
# HELP kubeteach_exerciseset_pending_tasks total pending of exerciseset
# TYPE kubeteach_exerciseset_pending_tasks gauge
kubeteach_exerciseset_pending_tasks{name="exerciseset1",namespace="default"} 1
`),
	"kubeteach_exerciseset_points_achieved": strings.NewReader(`
# HELP kubeteach_exerciseset_points_achieved achieved points of exerciseset
# TYPE kubeteach_exerciseset_points_achieved gauge
kubeteach_exerciseset_points_achieved{name="exerciseset1",namespace="default"} 3
`),
	"kubeteach_exerciseset_points": strings.NewReader(`
# HELP kubeteach_exerciseset_points total points of exerciseset
# TYPE kubeteach_exerciseset_points gauge
kubeteach_exerciseset_points{name="exerciseset1",namespace="default"} 10
`),
	"kubeteach_exerciseset_successful_tasks": strings.NewReader(`
# HELP kubeteach_exerciseset_successful_tasks successful points of exerciseset
# TYPE kubeteach_exerciseset_successful_tasks gauge
kubeteach_exerciseset_successful_tasks{name="exerciseset1",namespace="default"} 2
`),
	"kubeteach_exerciseset_tasks": strings.NewReader(`
# HELP kubeteach_exerciseset_tasks total tasks of exerciseset
# TYPE kubeteach_exerciseset_tasks gauge
kubeteach_exerciseset_tasks{name="exerciseset1",namespace="default"} 6
`),
	"kubeteach_exerciseset_tasks_without_points": strings.NewReader(`
# HELP kubeteach_exerciseset_tasks_without_points task without points of exerciseset
# TYPE kubeteach_exerciseset_tasks_without_points gauge
kubeteach_exerciseset_tasks_without_points{name="exerciseset1",namespace="default"} 2
`),
	"kubeteach_exerciseset_unknown_tasks": strings.NewReader(`
# HELP kubeteach_exerciseset_unknown_tasks unknown points of exerciseset
# TYPE kubeteach_exerciseset_unknown_tasks gauge
kubeteach_exerciseset_unknown_tasks{name="exerciseset1",namespace="default"} 0
`),
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

var expectTask = map[string]io.Reader{
	"kubeteach_task_state": strings.NewReader(`
# HELP kubeteach_task_state state of task (1 = successful, 2 = active, 3 = pending, 4 = unknown)
# TYPE kubeteach_task_state gauge
kubeteach_task_state{name="task1",namespace="default"} 1
kubeteach_task_state{name="task2",namespace="default"} 2
kubeteach_task_state{name="task3",namespace="default"} 3
kubeteach_task_state{name="task4",namespace="default"} 4
`),
}
