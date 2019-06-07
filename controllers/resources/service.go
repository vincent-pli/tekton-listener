/*
Copyright 2019 The Knative Authors

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

package resources

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	servingv1alpha1 "github.com/knative/serving/pkg/apis/serving/v1alpha1"
	tektonexperimentalv1alpha1 "github.com/vincent-pli/tekton-listener/api/v1alpha1"
)

// MakeService generates, but does not create, a Service for the given
// GitLabSource.
func MakeService(source tektonexperimentalv1alpha1.EventBinding, receiveAdapterImage string) *servingv1alpha1.Service {
	labels := map[string]string{
		"listener-adapter": "normal",
	}

	return &servingv1alpha1.Service{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: fmt.Sprintf("%s-", source.Name),
			Namespace:    source.Namespace,
			Labels:       labels,
		},
		Spec: servingv1alpha1.ServiceSpec{
			ConfigurationSpec: &servingv1alpha1.ConfigurationSpec{
				Template: servingv1alpha1.RevisionTemplateSpec{
						Spec: servingv1alpha1.RevisionSpec{
							Container: corev1.Container{
								Image: receiveAdapterImage,
							},
						},
			},
		},
	}
}
