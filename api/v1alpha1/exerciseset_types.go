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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//+kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="PointsAchieved",type=string,JSONPath=`.status.pointsAchieved`
// +kubebuilder:printcolumn:name="PointsTotal",type=string,JSONPath=`.status.pointsTotal`
// +kubebuilder:printcolumn:name="Tasks",type=string,JSONPath=`.status.numberOfTasks`
// +kubebuilder:printcolumn:name="Successful",type=string,JSONPath=`.status.numberOfSuccessfulTasks`
// +kubebuilder:printcolumn:name="Active",type=string,JSONPath=`.status.numberOfActiveTasks`
// +kubebuilder:printcolumn:name="Pending",type=string,JSONPath=`.status.numberOfPendingTasks`
//+kubebuilder:subresource:status

// ExerciseSet is the Schema for the exercisesets API
type ExerciseSet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ExerciseSetSpec   `json:"spec,omitempty"`
	Status ExerciseSetStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ExerciseSetList contains a list of ExerciseSet
type ExerciseSetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ExerciseSet `json:"items"`
}

// ExerciseSetSpec defines the desired state of ExerciseSet
type ExerciseSetSpec struct {
	// TaskDefinitionSpec represents the Spec of an TaskDefinition
	// +kubebuilder:validation:Required
	TaskDefinitions []ExerciseSetSpecTaskDefinitions `json:"taskDefinitions,omitempty"`
}

// ExerciseSetSpecTaskDefinitions defines the desired state of ExerciseSet
type ExerciseSetSpecTaskDefinitions struct {
	// Name is the name of the TaskDefinition
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`
	// TaskDefinitionSpec represents the Spec of an TaskDefinition
	// +kubebuilder:validation:Required
	TaskDefinitionSpec TaskDefinitionSpec `json:"taskDefinitionSpec"`
}

// ExerciseSetStatus defines the observed state of ExerciseSet
type ExerciseSetStatus struct {
	// NumberOfTasks is the number of total tasks of this ExerciseSet
	// +optional
	NumberOfTasks int `json:"numberOfTasks"`
	// NumberOfActiveTasks is the number of active tasks of this ExerciseSet
	// +optional
	NumberOfActiveTasks int `json:"numberOfActiveTasks"`
	// NumberOfPendingTasks is the number of pending tasks of this ExerciseSet
	// +optional
	NumberOfPendingTasks int `json:"numberOfPendingTasks"`
	// NumberOfSuccessfulTasks is the number of successful tasks of this ExerciseSet
	// +optional
	NumberOfSuccessfulTasks int `json:"numberOfSuccessfulTasks"`
	// NumberOfUnknownTasks is the number of tasks with an unknown state of this ExerciseSet
	// +optional
	NumberOfUnknownTasks int `json:"numberOfUnknownTasks"`
	// NumberOfTasksWithoutPoints is the number of tasks that have no points defined of this ExerciseSet
	// +optional
	NumberOfTasksWithoutPoints int `json:"numberOfTasksWithoutPoints"`
	// PointsTotal is the total sum of points for all tasks of this ExerciseSet
	// +optional
	PointsTotal int `json:"pointsTotal"`
	// PointsAchieved is the total sum of points for all tasks that are successful of this ExerciseSet
	// +optional
	PointsAchieved int `json:"pointsAchieved"`
}

func init() {
	SchemeBuilder.Register(&ExerciseSet{}, &ExerciseSetList{})
}
