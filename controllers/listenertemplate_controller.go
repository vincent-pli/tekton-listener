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
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	tektonexperimentalv1alpha1 "github.com/vincent-pli/tekton-listener/api/v1alpha1"
)

// ListenerTemplateReconciler reconciles a ListenerTemplate object
type ListenerTemplateReconciler struct {
	client.Client
	Log logr.Logger
}

// +kubebuilder:rbac:groups=tektonexperimental.vincent-pli.com,resources=listenertemplates,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=tektonexperimental.vincent-pli.com,resources=listenertemplates/status,verbs=get;update;patch

func (r *ListenerTemplateReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("listenertemplate", req.NamespacedName)

	// your logic here
	var listenerTemplate tektonexperimentalv1alpha1.ListenerTemplate
	if err := r.Get(ctx, req.NamespacedName, &listenerTemplate); err != nil {
		log.Error(err, "unable to fetch ListenerTemplate")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, ignoreNotFound(err)
	}
	log.Info("xxxxxxxxx")
	log.Info(listenerTemplate.Name)
	return ctrl.Result{}, nil
}

func ignoreNotFound(err error) error {
	if apierrs.IsNotFound(err) {
		return nil
	}
	return err
}

// SetupWithManager is
func (r *ListenerTemplateReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&tektonexperimentalv1alpha1.ListenerTemplate{}).
		Complete(r)
}
