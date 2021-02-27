/*
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
	Spec              TaskDefinitionSpec   `json:"spec,omitempty"`
	Status            TaskDefinitionStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// TaskDefinitionList contains a list of TaskDefinition
type TaskDefinitionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TaskDefinition `json:"items"`
}

// TaskDefinitionSpec defines the desired state of TaskDefinition
type TaskDefinitionSpec struct {
	// TaskSpec TODO
	TaskSpec TaskSpec `json:"taskSpec,omitempty"`
	// TaskConditions TODO
	TaskConditions []TaskCondition `json:"taskConditions,omitempty"`
	// RequiredTaskName TODO
	RequiredTaskName *string `json:"requiredTaskName,omitempty"`
	// PreApply TODO
	PreApply *[]string `json:"preApply,omitempty"`
}

// TaskCondition TODO
type TaskCondition struct {
	// ApiVersion is used of the object that should be match this conditions
	ApiVersion string `json:"apiVersion,omitempty"`
	// Kind is used of the object that should be match this conditions
	Kind string `json:"kind,omitempty"`
	// ApiGroup is used of the object that should be match this conditions
	ApiGroup string `json:"apiGroup,omitempty"`
	// MatchAll it set to true, ResourceCondition must be successful on all objects of this type
	// Useful to check if a object is deleted
	MatchAll bool `json:"matchAll,omitempty"`
	// ResourceCondition describe the conditions that must be apply to success this TaskCondition
	ResourceCondition []ResourceCondition `json:"resourceCondition,omitempty"`
}

// ResourceCondition TODO
type ResourceCondition struct {
	// Field is the json search string for this condition.
	// Example: metadata.name
	// For more details have a look into gjson docs: https://github.com/tidwall/gjson
	Field string `json:"field,omitempty"`
	// Operator is for the condition.
	// Valid operators are eq, neq, lt, gt, nil, notnil contains.
	// +kubebuilder:validation:Enum=eq;neq;lt;gt;contains;nil;notnil
	Operator string `json:"operator,omitempty"`
	// Value contains the value which the operater must match.
	// Must be a string but for lt and gt only numbers are allowed in this string.
	// Value is ignored by operator nil and notnil
	Value string `json:"value,omitempty"`
}

// TaskDefinitionStatus defines the observed state of TaskDefinition
type TaskDefinitionStatus struct {
	// State represent the status of this task
	// Can be pending, active, successful
	State *string `json:"state"`
}

func init() {
	SchemeBuilder.Register(&TaskDefinition{}, &TaskDefinitionList{})
}
