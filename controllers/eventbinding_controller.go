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
	servingv1alpha1 "github.com/knative/serving/pkg/apis/serving/v1alpha1"
	tektonexperimentalv1alpha1 "github.com/vincent-pli/tekton-listener/api/v1alpha1"
	resources "github.com/vincent-pli/tekton-listener/controllers/resources"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// EventBindingReconciler reconciles a EventBinding object
type EventBindingReconciler struct {
	client.Client
	Log                  logr.Logger
	ListenerAdapterImage string
	Scheme               *runtime.Scheme
	Recorder             record.EventRecorder
}

// +kubebuilder:rbac:groups=tektonexperimental.vincent-pli.com,resources=eventbindings,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=tektonexperimental.vincent-pli.com,resources=eventbindings/status,verbs=get;update;patch

func (r *EventBindingReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
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

	var reconcileErr error
	if eventBinding.DeletionTimestamp == nil {
		reconcileErr = r.reconcile(ctx, &eventBinding)
	} else {
		reconcileErr = r.finalize(ctx, &eventBinding)
	}

	return ctrl.Result{}, reconcileErr
}

func (r *EventBindingReconciler) reconcile(ctx context.Context, source *tektonexperimentalv1alpha1.EventBinding) (ctrl.Result, error) {
	ksvc, err := r.getOwnedService(ctx, source)
	if err != nil {
		if apierrors.IsNotFound(err) {
			ksvc = resources.MakeService(source, r.ListenerAdapterImage)
			if err = controllerutil.SetControllerReference(source, ksvc, r.Scheme); err != nil {
				return err
			}
			if err = r.Client.Create(ctx, ksvc); err != nil {
				return err
			}
			r.Recorder.Eventf(source, corev1.EventTypeNormal, "ServiceCreated", "Created Service %q", ksvc.Name)
			// TODO: Mark Deploying for the ksvc
			// Wait for the Service to get a status
			return nil
		}
		// Error was something other than NotFound
		return err
	}

	routeCondition := ksvc.Status.GetCondition(servingv1alpha1.ServiceConditionRoutesReady)
	receiveAdapterDomain := "fake"
	if routeCondition != nil && routeCondition.Status == corev1.ConditionTrue && receiveAdapterDomain != "" {
		// TODO: Mark Deployed for the ksvc
		// TODO: Mark some condition for the webhook status?
		r.addFinalizer(source)
		fmt.Println("ksvc is ready...")
	} else {
		return ctrl.Result{Requeue: true, RequeueAfter: 10}, err
	}
	return nil
}

func (r *EventBindingReconciler) finalize(ctx context.Context, source *tektonexperimentalv1alpha1.EventBinding) (ctrl.Result, error) {
	// Always remove the finalizer. If there's a failure cleaning up, an event
	// will be recorded allowing the webhook to be removed manually by the
	// operator.
	r.removeFinalizer(source)
	return ctrl.Result{}, nil
}

func (r *EventBindingReconciler) addFinalizer(s *tektonexperimentalv1alpha1.EventBinding) {
	finalizers := sets.NewString(s.Finalizers...)
	finalizers.Insert("tekton-listener-source-controlle")
	s.Finalizers = finalizers.List()
}

func (r *EventBindingReconciler) removeFinalizer(s *tektonexperimentalv1alpha1.EventBinding) {
	finalizers := sets.NewString(s.Finalizers...)
	finalizers.Delete("tekton-listener-source-controlle")
	s.Finalizers = finalizers.List()
}

func (r *EventBindingReconciler) getOwnedService(ctx context.Context, source *tektonexperimentalv1alpha1.EventBinding) (*servingv1alpha1.Service, error) {
	list := &servingv1alpha1.ServiceList{}
	err := r.Client.List(ctx, list, client.InNamespace(source.Namespace))
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
