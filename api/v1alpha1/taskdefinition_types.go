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

// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Status",type=string,JSONPath=`.status.state`
// +kubebuilder:subresource:status

// TaskDefinition is the Schema for the taskdefinitions API
type TaskDefinition struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TaskDefinitionSpec   `json:"spec,omitempty"`
	Status TaskDefinitionStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// TaskDefinitionList contains a list of TaskDefinition
type TaskDefinitionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TaskDefinition `json:"items"`
}

// TaskDefinitionSpec defines the desired state of TaskDefinition.
type TaskDefinitionSpec struct {
	// TaskSpec represents spec of the task that is creating for this TaskDefinition.
	// +kubebuilder:validation:Required
	TaskSpec TaskSpec `json:"taskSpec"`
	// TaskConditions defines a list of conditions for a object that must be true to complete the task.
	// +kubebuilder:validation:MinItems=1
	TaskConditions []TaskCondition `json:"taskCondition"`
	// RequiredTaskName defines a TaskDefinition Name that have to be done before.
	// Useful for example if in task1 a object should be created and in task2 the object should be deleted again.
	//  +optional
	RequiredTaskName *string `json:"requiredTaskName,omitempty"`
}

// TaskCondition defines a list of conditions for a object that must be true to complete the task.
type TaskCondition struct {
	// APIVersion is used of the object that should be match this conditions
	// +kubebuilder:validation:MinLength=1
	APIVersion string `json:"apiVersion"`
	// Kind is used of the object that should be match this conditions
	// +kubebuilder:validation:MinLength=1
	Kind string `json:"kind"`
	// APIGroup is used of the object that should be match this conditions
	//  +optional
	APIGroup string `json:"apiGroup,omitempty"`
	// Name defines the name of the object that must apply to this conditions
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`
	// Namespace is used to find the object if it is namespaced
	//  +optional
	Namespace string `json:"namespace,omitempty"`
	// NotExists if set to true, all ResourceCondition are ignored and the TaskCondition is true if object do not exists
	//  +optional
	NotExists bool `json:"notExists,omitempty"`
	// ResourceCondition describe the conditions that must be apply to success this TaskCondition
	// If no ResourceCondition is set this TaskCondition just check if object exits
	//  +optional
	ResourceCondition []ResourceCondition `json:"resourceCondition,omitempty"`
}

// ResourceCondition describe the conditions that must be apply to success this TaskCondition
type ResourceCondition struct {
	// Field is the json search string for this condition.
	// Example: metadata.name
	// For more details have a look into gjson docs: https://github.com/tidwall/gjson
	// +kubebuilder:validation:MinLength=1
	Field string `json:"field"`
	// Operator is for the condition.
	// Valid operators are eq, neq, lt, gt, nil, notnil contains.
	// +kubebuilder:validation:Enum=eq;neq;lt;gt;contains;nil;notnil
	Operator string `json:"operator"`
	// Value contains the value which the Operator must match.
	// Must be a string but for lt and gt only numbers are allowed in this string.
	// Value is ignored by Operator nil and notnil
	//  +optional
	Value string `json:"value,omitempty"`
}

// TaskDefinitionStatus defines the observed state of TaskDefinition
type TaskDefinitionStatus struct {
	// State represent the status of this task
	// Can be pending, active, successful, error
	//  +optional
	State *string `json:"state"`
}

func init() {
	SchemeBuilder.Register(&TaskDefinition{}, &TaskDefinitionList{})
}
