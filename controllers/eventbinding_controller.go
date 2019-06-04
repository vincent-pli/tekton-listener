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

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	tektonexperimentalv1alpha1 "github.com/vincent-pli/tekton-listener/api/v1alpha1"
)

// EventBindingReconciler reconciles a EventBinding object
type EventBindingReconciler struct {
	client.Client
	Log logr.Logger
}

// +kubebuilder:rbac:groups=tektonexperimental.vincent-pli.com,resources=eventbindings,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=tektonexperimental.vincent-pli.com,resources=eventbindings/status,verbs=get;update;patch

func (r *EventBindingReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("eventbinding", req.NamespacedName)

	// your logic here

	return ctrl.Result{}, nil
}

func (r *EventBindingReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&tektonexperimentalv1alpha1.EventBinding{}).
		Complete(r)
}
