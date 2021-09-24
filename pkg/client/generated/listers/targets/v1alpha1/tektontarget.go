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

// TektonTargetLister helps list TektonTargets.
// All objects returned here must be treated as read-only.
type TektonTargetLister interface {
	// List lists all TektonTargets in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.TektonTarget, err error)
	// TektonTargets returns an object that can list and get TektonTargets.
	TektonTargets(namespace string) TektonTargetNamespaceLister
	TektonTargetListerExpansion
}

// tektonTargetLister implements the TektonTargetLister interface.
type tektonTargetLister struct {
	indexer cache.Indexer
}

// NewTektonTargetLister returns a new TektonTargetLister.
func NewTektonTargetLister(indexer cache.Indexer) TektonTargetLister {
	return &tektonTargetLister{indexer: indexer}
}

// List lists all TektonTargets in the indexer.
func (s *tektonTargetLister) List(selector labels.Selector) (ret []*v1alpha1.TektonTarget, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.TektonTarget))
	})
	return ret, err
}

// TektonTargets returns an object that can list and get TektonTargets.
func (s *tektonTargetLister) TektonTargets(namespace string) TektonTargetNamespaceLister {
	return tektonTargetNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// TektonTargetNamespaceLister helps list and get TektonTargets.
// All objects returned here must be treated as read-only.
type TektonTargetNamespaceLister interface {
	// List lists all TektonTargets in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.TektonTarget, err error)
	// Get retrieves the TektonTarget from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.TektonTarget, error)
	TektonTargetNamespaceListerExpansion
}

// tektonTargetNamespaceLister implements the TektonTargetNamespaceLister
// interface.
type tektonTargetNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all TektonTargets in the indexer for a given namespace.
func (s tektonTargetNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.TektonTarget, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.TektonTarget))
	})
	return ret, err
}

// Get retrieves the TektonTarget from the indexer for a given namespace and name.
func (s tektonTargetNamespaceLister) Get(name string) (*v1alpha1.TektonTarget, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("tektontarget"), name)
	}
	return obj.(*v1alpha1.TektonTarget), nil
}