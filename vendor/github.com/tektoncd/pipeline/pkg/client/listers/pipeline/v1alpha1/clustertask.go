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
package v1alpha1

import (
	v1alpha1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// ClusterTaskLister helps list ClusterTasks.
type ClusterTaskLister interface {
	// List lists all ClusterTasks in the indexer.
	List(selector labels.Selector) (ret []*v1alpha1.ClusterTask, err error)
	// Get retrieves the ClusterTask from the index for a given name.
	Get(name string) (*v1alpha1.ClusterTask, error)
	ClusterTaskListerExpansion
}

// clusterTaskLister implements the ClusterTaskLister interface.
type clusterTaskLister struct {
	indexer cache.Indexer
}

// NewClusterTaskLister returns a new ClusterTaskLister.
func NewClusterTaskLister(indexer cache.Indexer) ClusterTaskLister {
	return &clusterTaskLister{indexer: indexer}
}

// List lists all ClusterTasks in the indexer.
func (s *clusterTaskLister) List(selector labels.Selector) (ret []*v1alpha1.ClusterTask, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.ClusterTask))
	})
	return ret, err
}

// Get retrieves the ClusterTask from the index for a given name.
func (s *clusterTaskLister) Get(name string) (*v1alpha1.ClusterTask, error) {
	obj, exists, err := s.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("clustertask"), name)
	}
	return obj.(*v1alpha1.ClusterTask), nil
}
