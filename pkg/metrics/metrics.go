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

// Package metrics is a prometheus exporter for kubeteach
package metrics

import (
	"context"
	"errors"

	kubeteachv1alpha1 "github.com/dergeberl/kubeteach/api/v1alpha1"
	"github.com/dergeberl/kubeteach/controllers"
	"github.com/go-logr/logr"
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	metricPrefix = "kubeteach"
)

// Exporter represents a prometheus exporter for kubeteach
type Exporter struct {
	log                           logr.Logger
	k8sClient                     client.Client
	ExerciseSetPointsTotal        *prometheus.Desc
	ExerciseSetPointsAchieved     *prometheus.Desc
	ExerciseSetTasksWithoutPoints *prometheus.Desc
	ExerciseSetTasks              *prometheus.Desc
	ExerciseSetActiveTasks        *prometheus.Desc
	ExerciseSetPendingTasks       *prometheus.Desc
	ExerciseSetSuccessfulTasks    *prometheus.Desc
	ExerciseSetUnknownTasks       *prometheus.Desc
	TaskState                     *prometheus.Desc
}

// New generate a new Exporter object
func New(k8sClient client.Client, log logr.Logger) *Exporter {
	labels := []string{"name", "namespace"}
	return &Exporter{
		log:       log,
		k8sClient: k8sClient,
		ExerciseSetPointsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(metricPrefix, "exerciseset", "points"),
			"total points of exerciseset",
			labels,
			nil),
		ExerciseSetPointsAchieved: prometheus.NewDesc(
			prometheus.BuildFQName(metricPrefix, "exerciseset", "points_achieved"),
			"achieved points of exerciseset",
			labels,
			nil),
		ExerciseSetTasksWithoutPoints: prometheus.NewDesc(
			prometheus.BuildFQName(metricPrefix, "exerciseset", "tasks_without_points"),
			"task without points of exerciseset",
			labels,
			nil),
		ExerciseSetTasks: prometheus.NewDesc(
			prometheus.BuildFQName(metricPrefix, "exerciseset", "tasks"),
			"total tasks of exerciseset",
			labels,
			nil),
		ExerciseSetActiveTasks: prometheus.NewDesc(
			prometheus.BuildFQName(metricPrefix, "exerciseset", "active_tasks"),
			"active points of exerciseset",
			labels,
			nil),
		ExerciseSetPendingTasks: prometheus.NewDesc(
			prometheus.BuildFQName(metricPrefix, "exerciseset", "pending_tasks"),
			"total pending of exerciseset",
			labels,
			nil),
		ExerciseSetSuccessfulTasks: prometheus.NewDesc(
			prometheus.BuildFQName(metricPrefix, "exerciseset", "successful_tasks"),
			"successful points of exerciseset",
			labels,
			nil),
		ExerciseSetUnknownTasks: prometheus.NewDesc(
			prometheus.BuildFQName(metricPrefix, "exerciseset", "unknown_tasks"),
			"unknown points of exerciseset",
			labels,
			nil),
		TaskState: prometheus.NewDesc(
			prometheus.BuildFQName(metricPrefix, "task", "state"),
			"state of task (1 = successful, 2 = active, 3 = pending, 4 = unknown)",
			labels,
			nil),
	}
}

// Describe describes metrics for kubeteach
func (e *Exporter) Describe(descs chan<- *prometheus.Desc) {
	descs <- e.ExerciseSetPointsTotal
	descs <- e.ExerciseSetPointsAchieved
	descs <- e.ExerciseSetTasksWithoutPoints
	descs <- e.ExerciseSetTasks
	descs <- e.ExerciseSetActiveTasks
	descs <- e.ExerciseSetPendingTasks
	descs <- e.ExerciseSetSuccessfulTasks
	descs <- e.ExerciseSetUnknownTasks
	descs <- e.TaskState
}

// Collect collects metrics for kubeteach
func (e *Exporter) Collect(metrics chan<- prometheus.Metric) {
	if e.k8sClient == nil {
		if e.log != nil {
			e.log.Error(errors.New("no k8s client found"), "failed to get data from kubernetes")
		}
		return
	}
	err := e.collectExerciseSet(metrics)
	if err != nil {
		if e.log != nil {
			e.log.Error(err, "failed to get ExerciseSet metrics")
		}
	}
	err = e.collectTask(metrics)
	if err != nil {
		if e.log != nil {
			e.log.Error(err, "failed to get task metrics")
		}
	}
}

func (e *Exporter) collectExerciseSet(metrics chan<- prometheus.Metric) error {
	ctx := context.Background()
	exerciseSetList := &kubeteachv1alpha1.ExerciseSetList{}
	err := e.k8sClient.List(ctx, exerciseSetList)
	if err != nil {
		return err
	}
	for _, exerciseSet := range exerciseSetList.Items {
		metrics <- prometheus.MustNewConstMetric(
			e.ExerciseSetPointsTotal,
			prometheus.GaugeValue,
			float64(exerciseSet.Status.PointsTotal),
			exerciseSet.Name, exerciseSet.Namespace,
		)
		metrics <- prometheus.MustNewConstMetric(
			e.ExerciseSetPointsAchieved,
			prometheus.GaugeValue,
			float64(exerciseSet.Status.PointsAchieved),
			exerciseSet.Name, exerciseSet.Namespace,
		)
		metrics <- prometheus.MustNewConstMetric(
			e.ExerciseSetTasksWithoutPoints,
			prometheus.GaugeValue,
			float64(exerciseSet.Status.NumberOfTasksWithoutPoints),
			exerciseSet.Name, exerciseSet.Namespace,
		)
		metrics <- prometheus.MustNewConstMetric(
			e.ExerciseSetTasks,
			prometheus.GaugeValue,
			float64(exerciseSet.Status.NumberOfTasks),
			exerciseSet.Name, exerciseSet.Namespace,
		)
		metrics <- prometheus.MustNewConstMetric(
			e.ExerciseSetActiveTasks,
			prometheus.GaugeValue,
			float64(exerciseSet.Status.NumberOfActiveTasks),
			exerciseSet.Name, exerciseSet.Namespace,
		)
		metrics <- prometheus.MustNewConstMetric(
			e.ExerciseSetPendingTasks,
			prometheus.GaugeValue,
			float64(exerciseSet.Status.NumberOfPendingTasks),
			exerciseSet.Name, exerciseSet.Namespace,
		)
		metrics <- prometheus.MustNewConstMetric(
			e.ExerciseSetSuccessfulTasks,
			prometheus.GaugeValue,
			float64(exerciseSet.Status.NumberOfSuccessfulTasks),
			exerciseSet.Name, exerciseSet.Namespace,
		)
		metrics <- prometheus.MustNewConstMetric(
			e.ExerciseSetUnknownTasks,
			prometheus.GaugeValue,
			float64(exerciseSet.Status.NumberOfUnknownTasks),
			exerciseSet.Name, exerciseSet.Namespace,
		)
	}

	return nil
}

func (e *Exporter) collectTask(metrics chan<- prometheus.Metric) error {
	ctx := context.Background()
	taskList := &kubeteachv1alpha1.TaskList{}
	err := e.k8sClient.List(ctx, taskList)
	if err != nil {
		return err
	}
	for _, task := range taskList.Items {
		var state string
		if task.Status.State != nil {
			state = *task.Status.State
		}
		// 1 = successful
		// 2 = active
		// 3 = pending
		// 4 = unknown
		stateInt := 4
		switch state {
		case controllers.StateSuccessful:
			stateInt = 1
		case controllers.StateActive:
			stateInt = 2
		case controllers.StatePending:
			stateInt = 3
		}
		metrics <- prometheus.MustNewConstMetric(
			e.TaskState,
			prometheus.GaugeValue,
			float64(stateInt),
			task.Name, task.Namespace,
		)
	}
	return nil
}
