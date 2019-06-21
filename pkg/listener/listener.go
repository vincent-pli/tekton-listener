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
	"log"

	tektoncdclientset "github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	//"github.com/knative/eventing-sources/pkg/kncloudevents"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"github.com/kubernetes-sigs/controller-runtime/pkg/client"
	tektonlistenerv1alpha1 "github.com/vincent-pli/tekton-listener/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	GLHeaderEvent    = "GitLab-Event1"
	GLHeaderDelivery = "GitLab-Delivery"
)

var (
	scheme              = runtime.NewScheme()
)

// Adapter converts incoming GitLab webhook events to CloudEvents
type Listener struct {
	tektonClient   tektoncdclientset.Interface
	client client.Client
	Log                  logr.Logger
}

// New creates an adapter to convert incoming GitLab webhook events to CloudEvents and
// then sends them to the specified Sink
func New() (*Listener, error) {
	tektonlistenerv1alpha1.AddToScheme(scheme)
	l := new(Listener)

	ctrl.SetLogger(zap.Logger(true))
        l.Log = ctrl.Log.WithName("kvs").WithName("Listener")

	// Get cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		l.Log.Error(err, "error getting in cluster config.")
		return nil, err
	}

	// Setup dynamic client
	l.client, err = client.New(ctrl.GetConfigOrDie(), client.Options{scheme, nil})
	if err != nil {
		l.Log.Error(err, "error create run-time client.")
                return nil, err
	}

	// Setup tektoncd client
	l.tektonClient, err = tektoncdclientset.NewForConfig(config)
	if err != nil {
		l.Log.Error(err, "error create tekton client.")
                return nil, err
	}

	return l, nil
}

// HandleEvent is invoked whenever an event comes in from GitLab
func (l *Listener) HandleEvent(payload interface{}) {
	err := l.handleEvent(payload)
	if err != nil {
		log.Printf("unexpected error handling GitLab event: %s", err)
	}
}

func (l *Listener) handleEvent(payload interface{}) error {
	l.Log.Info("xxxxxxxxxxxxxxxxxxxxxx")
	// parse event? should event be parsed in main or here?


	// get EventBinding and ListenerTemplate for event binding and Kind parse


	// Create PipelineResource and PipelineRun
	return nil
}
