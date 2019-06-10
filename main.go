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

package main

import (
	"flag"
	"fmt"
	"os"

	servingv1alpha1 "github.com/knative/serving/pkg/apis/serving/v1alpha1"
	tektonexperimentalv1alpha1 "github.com/vincent-pli/tekton-listener/api/v1alpha1"
	"github.com/vincent-pli/tekton-listener/controllers"
	"k8s.io/apimachinery/pkg/runtime"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	// +kubebuilder:scaffold:imports
)

var (
	scheme              = runtime.NewScheme()
	setupLog            = ctrl.Log.WithName("setup")
	lsImageEnvVar       = "LS_IMAGE"
	controllerAgentName = "tekton-listener-source-controller"
)

func init() {

	tektonexperimentalv1alpha1.AddToScheme(scheme)
	servingv1alpha1.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.Parse()

	ctrl.SetLogger(zap.Logger(true))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{Scheme: scheme, MetricsBindAddress: metricsAddr})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	err = (&controllers.ListenerTemplateReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("ListenerTemplate"),
	}).SetupWithManager(mgr)
	if err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "ListenerTemplate")
		os.Exit(1)
	}

	listenerAdapterImage, defined := os.LookupEnv(lsImageEnvVar)
	if !defined {
		err = fmt.Errorf("required environment variable %q not defined", lsImageEnvVar)
	}

	err = (&controllers.EventBindingReconciler{
		Client:               mgr.GetClient(),
		Log:                  ctrl.Log.WithName("controllers").WithName("EventBinding"),
		ListenerAdapterImage: listenerAdapterImage,
		Scheme:               scheme,
		Recorder:             mgr.GetEventRecorderFor(controllerAgentName),
	}).SetupWithManager(mgr)
	if err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "EventBinding")
		os.Exit(1)
	}
	// +kubebuilder:scaffold:builder

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
