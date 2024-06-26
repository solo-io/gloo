/*
Copyright The Kubernetes Authors.

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
	json "encoding/json"
	"fmt"

	v1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/apis/gateway.solo.io/v1"
	gatewaysoloiov1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/client/applyconfiguration/gateway.solo.io/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeHttpListenerOptions implements HttpListenerOptionInterface
type FakeHttpListenerOptions struct {
	Fake *FakeGatewayV1
	ns   string
}

var httplisteneroptionsResource = v1.SchemeGroupVersion.WithResource("httplisteneroptions")

var httplisteneroptionsKind = v1.SchemeGroupVersion.WithKind("HttpListenerOption")

// Get takes name of the httpListenerOption, and returns the corresponding httpListenerOption object, and an error if there is any.
func (c *FakeHttpListenerOptions) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.HttpListenerOption, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(httplisteneroptionsResource, c.ns, name), &v1.HttpListenerOption{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1.HttpListenerOption), err
}

// List takes label and field selectors, and returns the list of HttpListenerOptions that match those selectors.
func (c *FakeHttpListenerOptions) List(ctx context.Context, opts metav1.ListOptions) (result *v1.HttpListenerOptionList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(httplisteneroptionsResource, httplisteneroptionsKind, c.ns, opts), &v1.HttpListenerOptionList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1.HttpListenerOptionList{ListMeta: obj.(*v1.HttpListenerOptionList).ListMeta}
	for _, item := range obj.(*v1.HttpListenerOptionList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested httpListenerOptions.
func (c *FakeHttpListenerOptions) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(httplisteneroptionsResource, c.ns, opts))

}

// Create takes the representation of a httpListenerOption and creates it.  Returns the server's representation of the httpListenerOption, and an error, if there is any.
func (c *FakeHttpListenerOptions) Create(ctx context.Context, httpListenerOption *v1.HttpListenerOption, opts metav1.CreateOptions) (result *v1.HttpListenerOption, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(httplisteneroptionsResource, c.ns, httpListenerOption), &v1.HttpListenerOption{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1.HttpListenerOption), err
}

// Update takes the representation of a httpListenerOption and updates it. Returns the server's representation of the httpListenerOption, and an error, if there is any.
func (c *FakeHttpListenerOptions) Update(ctx context.Context, httpListenerOption *v1.HttpListenerOption, opts metav1.UpdateOptions) (result *v1.HttpListenerOption, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(httplisteneroptionsResource, c.ns, httpListenerOption), &v1.HttpListenerOption{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1.HttpListenerOption), err
}

// Delete takes name of the httpListenerOption and deletes it. Returns an error if one occurs.
func (c *FakeHttpListenerOptions) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(httplisteneroptionsResource, c.ns, name, opts), &v1.HttpListenerOption{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeHttpListenerOptions) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(httplisteneroptionsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1.HttpListenerOptionList{})
	return err
}

// Patch applies the patch and returns the patched httpListenerOption.
func (c *FakeHttpListenerOptions) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.HttpListenerOption, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(httplisteneroptionsResource, c.ns, name, pt, data, subresources...), &v1.HttpListenerOption{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1.HttpListenerOption), err
}

// Apply takes the given apply declarative configuration, applies it and returns the applied httpListenerOption.
func (c *FakeHttpListenerOptions) Apply(ctx context.Context, httpListenerOption *gatewaysoloiov1.HttpListenerOptionApplyConfiguration, opts metav1.ApplyOptions) (result *v1.HttpListenerOption, err error) {
	if httpListenerOption == nil {
		return nil, fmt.Errorf("httpListenerOption provided to Apply must not be nil")
	}
	data, err := json.Marshal(httpListenerOption)
	if err != nil {
		return nil, err
	}
	name := httpListenerOption.Name
	if name == nil {
		return nil, fmt.Errorf("httpListenerOption.Name must be provided to Apply")
	}
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(httplisteneroptionsResource, c.ns, *name, types.ApplyPatchType, data), &v1.HttpListenerOption{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1.HttpListenerOption), err
}
