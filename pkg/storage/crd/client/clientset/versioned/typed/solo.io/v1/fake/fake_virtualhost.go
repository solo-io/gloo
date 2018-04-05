/*
Copyright 2018 The Kubernetes Authors.

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

package fake

import (
	solo_io_v1 "github.com/solo-io/gloo-storage/crd/solo.io/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeVirtualHosts implements VirtualHostInterface
type FakeVirtualHosts struct {
	Fake *FakeGlooV1
	ns   string
}

var virtualhostsResource = schema.GroupVersionResource{Group: "gloo.solo.io", Version: "v1", Resource: "virtualhosts"}

var virtualhostsKind = schema.GroupVersionKind{Group: "gloo.solo.io", Version: "v1", Kind: "VirtualHost"}

// Get takes name of the virtualHost, and returns the corresponding virtualHost object, and an error if there is any.
func (c *FakeVirtualHosts) Get(name string, options v1.GetOptions) (result *solo_io_v1.VirtualHost, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(virtualhostsResource, c.ns, name), &solo_io_v1.VirtualHost{})

	if obj == nil {
		return nil, err
	}
	return obj.(*solo_io_v1.VirtualHost), err
}

// List takes label and field selectors, and returns the list of VirtualHosts that match those selectors.
func (c *FakeVirtualHosts) List(opts v1.ListOptions) (result *solo_io_v1.VirtualHostList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(virtualhostsResource, virtualhostsKind, c.ns, opts), &solo_io_v1.VirtualHostList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &solo_io_v1.VirtualHostList{}
	for _, item := range obj.(*solo_io_v1.VirtualHostList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested virtualHosts.
func (c *FakeVirtualHosts) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(virtualhostsResource, c.ns, opts))

}

// Create takes the representation of a virtualHost and creates it.  Returns the server's representation of the virtualHost, and an error, if there is any.
func (c *FakeVirtualHosts) Create(virtualHost *solo_io_v1.VirtualHost) (result *solo_io_v1.VirtualHost, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(virtualhostsResource, c.ns, virtualHost), &solo_io_v1.VirtualHost{})

	if obj == nil {
		return nil, err
	}
	return obj.(*solo_io_v1.VirtualHost), err
}

// Update takes the representation of a virtualHost and updates it. Returns the server's representation of the virtualHost, and an error, if there is any.
func (c *FakeVirtualHosts) Update(virtualHost *solo_io_v1.VirtualHost) (result *solo_io_v1.VirtualHost, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(virtualhostsResource, c.ns, virtualHost), &solo_io_v1.VirtualHost{})

	if obj == nil {
		return nil, err
	}
	return obj.(*solo_io_v1.VirtualHost), err
}

// Delete takes name of the virtualHost and deletes it. Returns an error if one occurs.
func (c *FakeVirtualHosts) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(virtualhostsResource, c.ns, name), &solo_io_v1.VirtualHost{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeVirtualHosts) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(virtualhostsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &solo_io_v1.VirtualHostList{})
	return err
}

// Patch applies the patch and returns the patched virtualHost.
func (c *FakeVirtualHosts) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *solo_io_v1.VirtualHost, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(virtualhostsResource, c.ns, name, data, subresources...), &solo_io_v1.VirtualHost{})

	if obj == nil {
		return nil, err
	}
	return obj.(*solo_io_v1.VirtualHost), err
}
