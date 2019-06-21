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

package listener

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	tektoncdclientset "github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	"github.com/knative/eventing-sources/pkg/kncloudevents"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	tektonlistenerv1alpha1 "github.com/vincent-pli/tekton-listener/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"github.com/go-logr/logr"
)

const (
	GLHeaderEvent    = "GitLab-Event"
	GLHeaderDelivery = "GitLab-Delivery"
)

var (
	scheme              = runtime.NewScheme()
)

// Adapter converts incoming GitLab webhook events to CloudEvents
type Listener struct {
	tektonClient   tektoncdclientset.Interface
	client client
	Log                  logr.Logger
}

// New creates an adapter to convert incoming GitLab webhook events to CloudEvents and
// then sends them to the specified Sink
func New() (*Listener, error) {
	tektonlistenerv1alpha1.AddToScheme(scheme)
	l := new(Listener)

	// Get cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		logging.Log.Errorf("error getting in cluster config: %s.", err.Error())
		return Resource{}, err
	}

	// Setup dynamic client
	l.client, err := client.New(config.GetConfigOrDie(), client.Options{scheme, {}})
	if err != nil {
		fmt.Println("failed to create client")
		os.Exit(1)
	}

	// Setup tektoncd client
	l.tektonClient, err := tektoncdclientset.NewForConfig(config)
	if err != nil {
		logging.Log.Errorf("error building tekton clientset: %s.", err.Error())
		return Resource{}, err
	}

	ctrl.SetLogger(zap.Logger(true))
	l.log = ctrl.Log.WithName("kvs").WithName("Listener")
	return l, nil
}

// HandleEvent is invoked whenever an event comes in from GitLab
func (l *Listener) HandleEvent(payload interface{}) {
	err := a.handleEvent(payload)
	if err != nil {
		log.Printf("unexpected error handling GitLab event: %s", err)
	}
}

func (l *Listener) handleEvent(payload interface{}) error {
	l.Log.Info("xxxxxxxxxxxxxxxxxxxxxx")
	l.Log.Info(payload)
	// parse event? should event be parsed in main or here?


	// get EventBinding and ListenerTemplate for event binding and Kind parse


	// Create PipelineResource and PipelineRun

}