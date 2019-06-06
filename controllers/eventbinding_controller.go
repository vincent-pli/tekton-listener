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
	"fmt"
	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	servingv1alpha1 "github.com/knative/serving/pkg/apis/serving/v1alpha1"
	tektonexperimentalv1alpha1 "github.com/vincent-pli/tekton-listener/api/v1alpha1"
	resources "github.com/vincent-pli/tekton-listener/controllers/resources"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"k8s.io/apimachinery/pkg/runtime"
)

// EventBindingReconciler reconciles a EventBinding object
type EventBindingReconciler struct {
	client.Client
	Log logr.Logger
	listenerAdapterImage string
	scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=tektonexperimental.vincent-pli.com,resources=eventbindings,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=tektonexperimental.vincent-pli.com,resources=eventbindings/status,verbs=get;update;patch

func (r *EventBindingReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	cxt := context.Background()
	log := r.Log.WithValues("eventbinding", req.NamespacedName)

	// your logic here
	var eventBinding tektonexperimentalv1alpha1.EventBinding
	if err := r.Get(ctx, req.NamespacedName, &eventBinding); err != nil {
		log.Error(err, "unable to fetch EventBinding")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, ignoreNotFound(err)
	}
	log.Info("yyyyyyyyyyyy")
	log.Info(eventBinding.Name)
	
	if eventBinding.DeletionTimestamp == nil {
		return r.reconcile(ctx, eventBinding)
	} else {
		return r.finalize(ctx, eventBinding)
	}

	return ctrl.Result{}, nil
}

func (r *EventBindingReconciler) reconcile(cxt ctx context.Context, source tektonexperimentalv1alpha1.EventBinding) (ctrl.Result, error)) {
	ksvc, err := r.getOwnedService(ctx, source)
	if err != nil {
		if apierrors.IsNotFound(err) {
			ksvc = resources.MakeService(source, r.listenerAdapterImage)
			if err = controllerutil.SetControllerReference(source, ksvc, r.scheme); err != nil {
				return ctrl.Result{}, err
			}
			if err = r.client.Create(ctx, ksvc); err != nil {
				return ctrl.Result{}, err
			}
			r.recorder.Eventf(source, corev1.EventTypeNormal, "ServiceCreated", "Created Service %q", ksvc.Name)
			// TODO: Mark Deploying for the ksvc
			// Wait for the Service to get a status
			return ctrl.Result{}, nil
		}
		// Error was something other than NotFound
		return ctrl.Result{}, err
	}

	routeCondition := ksvc.Status.GetCondition(servingv1alpha1.ServiceConditionRoutesReady)
	receiveAdapterDomain := ksvc.Status.Domain
	if routeCondition != nil && routeCondition.Status == corev1.ConditionTrue && receiveAdapterDomain != "" {
		// TODO: Mark Deployed for the ksvc
		// TODO: Mark some condition for the webhook status?
		r.addFinalizer(source)
		fmt.Println("ksvc is ready...");
	} else {
		return ctrl.Result{true, 10}, err
	}
	return ctrl.Result{} nil
}

func (r *EventBindingReconciler) finalize(cxt ctx context.Context, source tektonexperimentalv1alpha1.EventBinding) (ctrl.Result, error) {
	
}

func (r *reconciler) getOwnedService(ctx context.Context, source tektonexperimentalv1alpha1.EventBinding) (*servingv1alpha1.Service, error) {
	list := &servingv1alpha1.ServiceList{}
	err := r.client.List(ctx, &client.ListOptions{
		Namespace:     source.Namespace,
		LabelSelector: labels.Everything(),
		// TODO this is here because the fake client needs it.
		// Remove this when it's no longer needed.
		Raw: &metav1.ListOptions{
			TypeMeta: metav1.TypeMeta{
				APIVersion: servingv1alpha1.SchemeGroupVersion.String(),
				Kind:       "Service",
			},
		},
	},
		list)
	if err != nil {
		return nil, err
	}
	for _, ksvc := range list.Items {
		if metav1.IsControlledBy(&ksvc, source) {
			//TODO if there are >1 controlled, delete all but first?
			return &ksvc, nil
		}
	}
	return nil, apierrors.NewNotFound(servingv1alpha1.Resource("services"), "")
}

func (r *EventBindingReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&tektonexperimentalv1alpha1.EventBinding{}).
		Complete(r)
}
