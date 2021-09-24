/*
Copyright 2021 TriggerMesh Inc.

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

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/triggermesh/triggermesh/pkg/apis/targets/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// GoogleCloudStorageTargetLister helps list GoogleCloudStorageTargets.
// All objects returned here must be treated as read-only.
type GoogleCloudStorageTargetLister interface {
	// List lists all GoogleCloudStorageTargets in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.GoogleCloudStorageTarget, err error)
	// GoogleCloudStorageTargets returns an object that can list and get GoogleCloudStorageTargets.
	GoogleCloudStorageTargets(namespace string) GoogleCloudStorageTargetNamespaceLister
	GoogleCloudStorageTargetListerExpansion
}

// googleCloudStorageTargetLister implements the GoogleCloudStorageTargetLister interface.
type googleCloudStorageTargetLister struct {
	indexer cache.Indexer
}

// NewGoogleCloudStorageTargetLister returns a new GoogleCloudStorageTargetLister.
func NewGoogleCloudStorageTargetLister(indexer cache.Indexer) GoogleCloudStorageTargetLister {
	return &googleCloudStorageTargetLister{indexer: indexer}
}

// List lists all GoogleCloudStorageTargets in the indexer.
func (s *googleCloudStorageTargetLister) List(selector labels.Selector) (ret []*v1alpha1.GoogleCloudStorageTarget, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.GoogleCloudStorageTarget))
	})
	return ret, err
}

// GoogleCloudStorageTargets returns an object that can list and get GoogleCloudStorageTargets.
func (s *googleCloudStorageTargetLister) GoogleCloudStorageTargets(namespace string) GoogleCloudStorageTargetNamespaceLister {
	return googleCloudStorageTargetNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// GoogleCloudStorageTargetNamespaceLister helps list and get GoogleCloudStorageTargets.
// All objects returned here must be treated as read-only.
type GoogleCloudStorageTargetNamespaceLister interface {
	// List lists all GoogleCloudStorageTargets in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.GoogleCloudStorageTarget, err error)
	// Get retrieves the GoogleCloudStorageTarget from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.GoogleCloudStorageTarget, error)
	GoogleCloudStorageTargetNamespaceListerExpansion
}

// googleCloudStorageTargetNamespaceLister implements the GoogleCloudStorageTargetNamespaceLister
// interface.
type googleCloudStorageTargetNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all GoogleCloudStorageTargets in the indexer for a given namespace.
func (s googleCloudStorageTargetNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.GoogleCloudStorageTarget, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.GoogleCloudStorageTarget))
	})
	return ret, err
}

// Get retrieves the GoogleCloudStorageTarget from the indexer for a given namespace and name.
func (s googleCloudStorageTargetNamespaceLister) Get(name string) (*v1alpha1.GoogleCloudStorageTarget, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("googlecloudstoragetarget"), name)
	}
	return obj.(*v1alpha1.GoogleCloudStorageTarget), nil
}