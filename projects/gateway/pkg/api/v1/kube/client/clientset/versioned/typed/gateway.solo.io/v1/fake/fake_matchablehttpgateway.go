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

	v1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/apis/gateway.solo.io/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeMatchableHttpGateways implements MatchableHttpGatewayInterface
type FakeMatchableHttpGateways struct {
	Fake *FakeGatewayV1
	ns   string
}

var matchablehttpgatewaysResource = v1.SchemeGroupVersion.WithResource("httpgateways")

var matchablehttpgatewaysKind = v1.SchemeGroupVersion.WithKind("MatchableHttpGateway")

// Get takes name of the matchableHttpGateway, and returns the corresponding matchableHttpGateway object, and an error if there is any.
func (c *FakeMatchableHttpGateways) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.MatchableHttpGateway, err error) {
	emptyResult := &v1.MatchableHttpGateway{}
	obj, err := c.Fake.
		Invokes(testing.NewGetActionWithOptions(matchablehttpgatewaysResource, c.ns, name, options), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1.MatchableHttpGateway), err
}

// List takes label and field selectors, and returns the list of MatchableHttpGateways that match those selectors.
func (c *FakeMatchableHttpGateways) List(ctx context.Context, opts metav1.ListOptions) (result *v1.MatchableHttpGatewayList, err error) {
	emptyResult := &v1.MatchableHttpGatewayList{}
	obj, err := c.Fake.
		Invokes(testing.NewListActionWithOptions(matchablehttpgatewaysResource, matchablehttpgatewaysKind, c.ns, opts), emptyResult)

	if obj == nil {
		return emptyResult, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1.MatchableHttpGatewayList{ListMeta: obj.(*v1.MatchableHttpGatewayList).ListMeta}
	for _, item := range obj.(*v1.MatchableHttpGatewayList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested matchableHttpGateways.
func (c *FakeMatchableHttpGateways) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchActionWithOptions(matchablehttpgatewaysResource, c.ns, opts))

}

// Create takes the representation of a matchableHttpGateway and creates it.  Returns the server's representation of the matchableHttpGateway, and an error, if there is any.
func (c *FakeMatchableHttpGateways) Create(ctx context.Context, matchableHttpGateway *v1.MatchableHttpGateway, opts metav1.CreateOptions) (result *v1.MatchableHttpGateway, err error) {
	emptyResult := &v1.MatchableHttpGateway{}
	obj, err := c.Fake.
		Invokes(testing.NewCreateActionWithOptions(matchablehttpgatewaysResource, c.ns, matchableHttpGateway, opts), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1.MatchableHttpGateway), err
}

// Update takes the representation of a matchableHttpGateway and updates it. Returns the server's representation of the matchableHttpGateway, and an error, if there is any.
func (c *FakeMatchableHttpGateways) Update(ctx context.Context, matchableHttpGateway *v1.MatchableHttpGateway, opts metav1.UpdateOptions) (result *v1.MatchableHttpGateway, err error) {
	emptyResult := &v1.MatchableHttpGateway{}
	obj, err := c.Fake.
		Invokes(testing.NewUpdateActionWithOptions(matchablehttpgatewaysResource, c.ns, matchableHttpGateway, opts), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1.MatchableHttpGateway), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeMatchableHttpGateways) UpdateStatus(ctx context.Context, matchableHttpGateway *v1.MatchableHttpGateway, opts metav1.UpdateOptions) (result *v1.MatchableHttpGateway, err error) {
	emptyResult := &v1.MatchableHttpGateway{}
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceActionWithOptions(matchablehttpgatewaysResource, "status", c.ns, matchableHttpGateway, opts), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1.MatchableHttpGateway), err
}

// Delete takes name of the matchableHttpGateway and deletes it. Returns an error if one occurs.
func (c *FakeMatchableHttpGateways) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(matchablehttpgatewaysResource, c.ns, name, opts), &v1.MatchableHttpGateway{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeMatchableHttpGateways) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	action := testing.NewDeleteCollectionActionWithOptions(matchablehttpgatewaysResource, c.ns, opts, listOpts)

	_, err := c.Fake.Invokes(action, &v1.MatchableHttpGatewayList{})
	return err
}

// Patch applies the patch and returns the patched matchableHttpGateway.
func (c *FakeMatchableHttpGateways) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.MatchableHttpGateway, err error) {
	emptyResult := &v1.MatchableHttpGateway{}
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceActionWithOptions(matchablehttpgatewaysResource, c.ns, name, pt, data, opts, subresources...), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1.MatchableHttpGateway), err
}
