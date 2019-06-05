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
	pipelinev1alpha1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ListenerTemplateSpec defines the desired state of ListenerTemplate
type ListenerTemplateSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Params      []TemplateParam                         `json:"params,omitempty"`
	Resources   []pipelinev1alpha1.PipelineResourceSpec `json:"resources,omitempty"`
	PipelineRun pipelinev1alpha1.PipelineRunSpec        `json:"pipelinerun,omitempty"`
}

// TemplateParam defines arbitrary parameters needed by a Pipelinerun, Resource defined in the ListenerTemplate.
type TemplateParam struct {
	Name string `json:"name"`
	// +optional
	Description string `json:"description,omitempty"`
	// +optional
	Default string `json:"default,omitempty"`
}

// ListenerTemplateStatus defines the observed state of ListenerTemplate
type ListenerTemplateStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	AvailableReference int32 `json:"availableReference"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// ListenerTemplate is the Schema for the listenertemplates API
// +kubebuilder:printcolumn:name="reference",type="int32",JSONPath=".status.availableReference""
type ListenerTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ListenerTemplateSpec   `json:"spec,omitempty"`
	Status ListenerTemplateStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ListenerTemplateList contains a list of ListenerTemplate
type ListenerTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ListenerTemplate `json:"items"`
}

// HasReference returns true if AvailableReference in Status is not 0 .
func (lt *ListenerTemplate) HasReference() bool {
	return lt.Status.AvailableReference != 0
}

func init() {
	SchemeBuilder.Register(&ListenerTemplate{}, &ListenerTemplateList{})
}
