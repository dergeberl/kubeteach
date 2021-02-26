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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TaskDefinitionSpec defines the desired state of TaskDefinition
type TaskDefinitionSpec struct {
	TaskSpec         TaskSpec        `json:"taskSpec,omitempty"`
	TaskConditions   []TaskCondition `json:"taskConditions,omitempty"`
	RequiredTaskName *string         `json:"requiredTaskName,omitempty"`
}

type TaskCondition struct {
	ApiVersion        string              `json:"apiVersion,omitempty"`
	Kind              string              `json:"kind,omitempty"`
	ApiGroup          string              `json:"apiGroup,omitempty"`
	MatchAll          bool                `json:"matchAll,omitempty"`
	ResourceCondition []ResourceCondition `json:"resourceCondition,omitempty"`
}

type ResourceCondition struct {
	Field    string `json:"field,omitempty"`
	Operator string `json:"operator,omitempty"`
	Value    string `json:"value,omitempty"`
}

// TaskDefinitionStatus defines the observed state of TaskDefinition
type TaskDefinitionStatus struct {
	State string `json:"state"`
}

// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Status",type=string,JSONPath=`.status.state`

// TaskDefinition is the Schema for the taskdefinitions API
type TaskDefinition struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TaskDefinitionSpec   `json:"spec,omitempty"`
	Status TaskDefinitionStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// TaskDefinitionList contains a list of TaskDefinition
type TaskDefinitionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TaskDefinition `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TaskDefinition{}, &TaskDefinitionList{})
}
