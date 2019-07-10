/*
Copyright 2019 vincent-pli.

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
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// EventBindingSpec defines the desired state of EventBinding
type EventBindingSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// ServiceAccountName holds the name of the Kubernetes service account
	// as which the underlying K8s resources should be run. If unspecified
	// this will default to the "default" service account for the namespace
	// in which the GitLabSource exists.
	// +optional
	ServiceAccountName string      `json:"serviceAccountName,omitempty"`
	TemplateRef        TemplateRef `json:"templateRef"`
	Event              Event       `json:"event"`
	Params             []Param     `json:"params,omitempty"`
}

// Event use to define a cloud event.
type Event struct {
	// Class of Cloudevent
	Class string `json:"class,omitempty"`
	// Type of Cloudevent
	Type string `json:"type,omitempty"`
}

// TemplateRef can be used to refer to a specific instance of a eventBinding.
type TemplateRef struct {
	// Name of the referent
	Name string `json:"name,omitempty"`
	// API version of the referent
	// +optional
	APIVersion string `json:"apiVersion,omitempty"`
}

// Param declares a value to use for the Param called Name.
type Param struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// EventBindingStatus defines the observed state of EventBinding
type EventBindingStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// EventBinding is the Schema for the eventbindings API
type EventBinding struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EventBindingSpec   `json:"spec,omitempty"`
	Status EventBindingStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true

// EventBindingList contains a list of EventBinding
type EventBindingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EventBinding `json:"items"`
}

// SetDefaults for pipelinerun
func (eb *EventBinding) SetDefaults(ctx context.Context) {}


func init() {
	SchemeBuilder.Register(&EventBinding{}, &EventBindingList{})
}
