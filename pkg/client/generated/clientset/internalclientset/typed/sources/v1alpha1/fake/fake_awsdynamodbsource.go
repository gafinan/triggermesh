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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1alpha1 "github.com/triggermesh/triggermesh/pkg/apis/sources/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeAWSDynamoDBSources implements AWSDynamoDBSourceInterface
type FakeAWSDynamoDBSources struct {
	Fake *FakeSourcesV1alpha1
	ns   string
}

var awsdynamodbsourcesResource = schema.GroupVersionResource{Group: "sources.triggermesh.io", Version: "v1alpha1", Resource: "awsdynamodbsources"}

var awsdynamodbsourcesKind = schema.GroupVersionKind{Group: "sources.triggermesh.io", Version: "v1alpha1", Kind: "AWSDynamoDBSource"}

// Get takes name of the aWSDynamoDBSource, and returns the corresponding aWSDynamoDBSource object, and an error if there is any.
func (c *FakeAWSDynamoDBSources) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.AWSDynamoDBSource, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(awsdynamodbsourcesResource, c.ns, name), &v1alpha1.AWSDynamoDBSource{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.AWSDynamoDBSource), err
}

// List takes label and field selectors, and returns the list of AWSDynamoDBSources that match those selectors.
func (c *FakeAWSDynamoDBSources) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.AWSDynamoDBSourceList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(awsdynamodbsourcesResource, awsdynamodbsourcesKind, c.ns, opts), &v1alpha1.AWSDynamoDBSourceList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.AWSDynamoDBSourceList{ListMeta: obj.(*v1alpha1.AWSDynamoDBSourceList).ListMeta}
	for _, item := range obj.(*v1alpha1.AWSDynamoDBSourceList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested aWSDynamoDBSources.
func (c *FakeAWSDynamoDBSources) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(awsdynamodbsourcesResource, c.ns, opts))

}

// Create takes the representation of a aWSDynamoDBSource and creates it.  Returns the server's representation of the aWSDynamoDBSource, and an error, if there is any.
func (c *FakeAWSDynamoDBSources) Create(ctx context.Context, aWSDynamoDBSource *v1alpha1.AWSDynamoDBSource, opts v1.CreateOptions) (result *v1alpha1.AWSDynamoDBSource, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(awsdynamodbsourcesResource, c.ns, aWSDynamoDBSource), &v1alpha1.AWSDynamoDBSource{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.AWSDynamoDBSource), err
}

// Update takes the representation of a aWSDynamoDBSource and updates it. Returns the server's representation of the aWSDynamoDBSource, and an error, if there is any.
func (c *FakeAWSDynamoDBSources) Update(ctx context.Context, aWSDynamoDBSource *v1alpha1.AWSDynamoDBSource, opts v1.UpdateOptions) (result *v1alpha1.AWSDynamoDBSource, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(awsdynamodbsourcesResource, c.ns, aWSDynamoDBSource), &v1alpha1.AWSDynamoDBSource{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.AWSDynamoDBSource), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeAWSDynamoDBSources) UpdateStatus(ctx context.Context, aWSDynamoDBSource *v1alpha1.AWSDynamoDBSource, opts v1.UpdateOptions) (*v1alpha1.AWSDynamoDBSource, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(awsdynamodbsourcesResource, "status", c.ns, aWSDynamoDBSource), &v1alpha1.AWSDynamoDBSource{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.AWSDynamoDBSource), err
}

// Delete takes name of the aWSDynamoDBSource and deletes it. Returns an error if one occurs.
func (c *FakeAWSDynamoDBSources) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(awsdynamodbsourcesResource, c.ns, name), &v1alpha1.AWSDynamoDBSource{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeAWSDynamoDBSources) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(awsdynamodbsourcesResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.AWSDynamoDBSourceList{})
	return err
}

// Patch applies the patch and returns the patched aWSDynamoDBSource.
func (c *FakeAWSDynamoDBSources) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.AWSDynamoDBSource, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(awsdynamodbsourcesResource, c.ns, name, pt, data, subresources...), &v1alpha1.AWSDynamoDBSource{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.AWSDynamoDBSource), err
}
