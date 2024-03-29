/*
Copyright 2019 The Tekton Authors

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

package reconciler

import (
	"time"

	cachingclientset "github.com/knative/caching/pkg/client/clientset/versioned"
	sharedclientset "github.com/knative/pkg/client/clientset/versioned"
	"github.com/knative/pkg/configmap"
	"github.com/knative/pkg/logging/logkey"
	clientset "github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	pipelineScheme "github.com/tektoncd/pipeline/pkg/client/clientset/versioned/scheme"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/record"
)

// Options defines the common reconciler options.
// We define this to reduce the boilerplate argument list when
// creating our controllers.
type Options struct {
	KubeClientSet     kubernetes.Interface
	SharedClientSet   sharedclientset.Interface
	PipelineClientSet clientset.Interface
	CachingClientSet  cachingclientset.Interface

	ConfigMapWatcher configmap.Watcher
	Logger           *zap.SugaredLogger
	Recorder         record.EventRecorder

	ResyncPeriod time.Duration
}

// GetTrackerLease returns a multiple of the resync period to use as the
// duration for tracker leases. This attempts to ensure that resyncs happen to
// refresh leases frequently enough that we don't miss updates to tracked
// objects.
func (o Options) GetTrackerLease() time.Duration {
	return o.ResyncPeriod * 3
}

// Base implements the core controller logic, given a Reconciler.
type Base struct {
	// KubeClientSet allows us to talk to the k8s for core APIs
	KubeClientSet kubernetes.Interface

	// SharedClientSet allows us to configure shared objects
	SharedClientSet sharedclientset.Interface

	// PipelineClientSet allows us to configure pipeline objects
	PipelineClientSet clientset.Interface

	// CachingClientSet allows us to instantiate Image objects
	CachingClientSet cachingclientset.Interface

	// ConfigMapWatcher allows us to watch for ConfigMap changes.
	ConfigMapWatcher configmap.Watcher

	// Recorder is an event recorder for recording Event resources to the
	// Kubernetes API.
	Recorder record.EventRecorder

	// Sugared logger is easier to use but is not as performant as the
	// raw logger. In performance critical paths, call logger.Desugar()
	// and use the returned raw logger instead. In addition to the
	// performance benefits, raw logger also preserves type-safety at
	// the expense of slightly greater verbosity.
	Logger *zap.SugaredLogger
}

// NewBase instantiates a new instance of Base implementing
// the common & boilerplate code between our reconcilers.
func NewBase(opt Options, controllerAgentName string) *Base {
	// Enrich the logs with controller name
	logger := opt.Logger.Named(controllerAgentName).With(zap.String(logkey.ControllerType, controllerAgentName))

	// Use recorder provided in options if presents.   Otherwise, create a new one.
	recorder := opt.Recorder

	if recorder == nil {
		// Create event broadcaster
		logger.Debug("Creating event broadcaster")
		eventBroadcaster := record.NewBroadcaster()
		eventBroadcaster.StartLogging(logger.Named("event-broadcaster").Infof)
		eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: opt.KubeClientSet.CoreV1().Events("")})

		recorder = eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})
	} else {
		logger.Debug("Using recorder from option")
	}

	base := &Base{
		KubeClientSet:     opt.KubeClientSet,
		SharedClientSet:   opt.SharedClientSet,
		PipelineClientSet: opt.PipelineClientSet,
		CachingClientSet:  opt.CachingClientSet,
		ConfigMapWatcher:  opt.ConfigMapWatcher,
		Recorder:          recorder,
		Logger:            logger,
	}

	return base
}

func init() {
	// Add pipeline types to the default Kubernetes Scheme so Events can be
	// logged for pipeline types.
	if err := pipelineScheme.AddToScheme(scheme.Scheme); err != nil {
		panic(err)
	}
}
