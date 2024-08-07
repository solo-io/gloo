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

package v1

import (
	"context"
	json "encoding/json"
	"fmt"
	"time"

	v1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/apis/gateway.solo.io/v1"
	gatewaysoloiov1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/client/applyconfiguration/gateway.solo.io/v1"
	scheme "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/client/clientset/versioned/scheme"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// MatchableTcpGatewaysGetter has a method to return a MatchableTcpGatewayInterface.
// A group's client should implement this interface.
type MatchableTcpGatewaysGetter interface {
	MatchableTcpGateways(namespace string) MatchableTcpGatewayInterface
}

// MatchableTcpGatewayInterface has methods to work with MatchableTcpGateway resources.
type MatchableTcpGatewayInterface interface {
	Create(ctx context.Context, matchableTcpGateway *v1.MatchableTcpGateway, opts metav1.CreateOptions) (*v1.MatchableTcpGateway, error)
	Update(ctx context.Context, matchableTcpGateway *v1.MatchableTcpGateway, opts metav1.UpdateOptions) (*v1.MatchableTcpGateway, error)
	UpdateStatus(ctx context.Context, matchableTcpGateway *v1.MatchableTcpGateway, opts metav1.UpdateOptions) (*v1.MatchableTcpGateway, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.MatchableTcpGateway, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.MatchableTcpGatewayList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.MatchableTcpGateway, err error)
	Apply(ctx context.Context, matchableTcpGateway *gatewaysoloiov1.MatchableTcpGatewayApplyConfiguration, opts metav1.ApplyOptions) (result *v1.MatchableTcpGateway, err error)
	ApplyStatus(ctx context.Context, matchableTcpGateway *gatewaysoloiov1.MatchableTcpGatewayApplyConfiguration, opts metav1.ApplyOptions) (result *v1.MatchableTcpGateway, err error)
	MatchableTcpGatewayExpansion
}

// matchableTcpGateways implements MatchableTcpGatewayInterface
type matchableTcpGateways struct {
	client rest.Interface
	ns     string
}

// newMatchableTcpGateways returns a MatchableTcpGateways
func newMatchableTcpGateways(c *GatewayV1Client, namespace string) *matchableTcpGateways {
	return &matchableTcpGateways{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the matchableTcpGateway, and returns the corresponding matchableTcpGateway object, and an error if there is any.
func (c *matchableTcpGateways) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.MatchableTcpGateway, err error) {
	result = &v1.MatchableTcpGateway{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("tcpgateways").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of MatchableTcpGateways that match those selectors.
func (c *matchableTcpGateways) List(ctx context.Context, opts metav1.ListOptions) (result *v1.MatchableTcpGatewayList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.MatchableTcpGatewayList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("tcpgateways").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested matchableTcpGateways.
func (c *matchableTcpGateways) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("tcpgateways").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a matchableTcpGateway and creates it.  Returns the server's representation of the matchableTcpGateway, and an error, if there is any.
func (c *matchableTcpGateways) Create(ctx context.Context, matchableTcpGateway *v1.MatchableTcpGateway, opts metav1.CreateOptions) (result *v1.MatchableTcpGateway, err error) {
	result = &v1.MatchableTcpGateway{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("tcpgateways").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(matchableTcpGateway).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a matchableTcpGateway and updates it. Returns the server's representation of the matchableTcpGateway, and an error, if there is any.
func (c *matchableTcpGateways) Update(ctx context.Context, matchableTcpGateway *v1.MatchableTcpGateway, opts metav1.UpdateOptions) (result *v1.MatchableTcpGateway, err error) {
	result = &v1.MatchableTcpGateway{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("tcpgateways").
		Name(matchableTcpGateway.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(matchableTcpGateway).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *matchableTcpGateways) UpdateStatus(ctx context.Context, matchableTcpGateway *v1.MatchableTcpGateway, opts metav1.UpdateOptions) (result *v1.MatchableTcpGateway, err error) {
	result = &v1.MatchableTcpGateway{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("tcpgateways").
		Name(matchableTcpGateway.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(matchableTcpGateway).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the matchableTcpGateway and deletes it. Returns an error if one occurs.
func (c *matchableTcpGateways) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("tcpgateways").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *matchableTcpGateways) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("tcpgateways").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched matchableTcpGateway.
func (c *matchableTcpGateways) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.MatchableTcpGateway, err error) {
	result = &v1.MatchableTcpGateway{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("tcpgateways").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

// Apply takes the given apply declarative configuration, applies it and returns the applied matchableTcpGateway.
func (c *matchableTcpGateways) Apply(ctx context.Context, matchableTcpGateway *gatewaysoloiov1.MatchableTcpGatewayApplyConfiguration, opts metav1.ApplyOptions) (result *v1.MatchableTcpGateway, err error) {
	if matchableTcpGateway == nil {
		return nil, fmt.Errorf("matchableTcpGateway provided to Apply must not be nil")
	}
	patchOpts := opts.ToPatchOptions()
	data, err := json.Marshal(matchableTcpGateway)
	if err != nil {
		return nil, err
	}
	name := matchableTcpGateway.Name
	if name == nil {
		return nil, fmt.Errorf("matchableTcpGateway.Name must be provided to Apply")
	}
	result = &v1.MatchableTcpGateway{}
	err = c.client.Patch(types.ApplyPatchType).
		Namespace(c.ns).
		Resource("tcpgateways").
		Name(*name).
		VersionedParams(&patchOpts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

// ApplyStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating ApplyStatus().
func (c *matchableTcpGateways) ApplyStatus(ctx context.Context, matchableTcpGateway *gatewaysoloiov1.MatchableTcpGatewayApplyConfiguration, opts metav1.ApplyOptions) (result *v1.MatchableTcpGateway, err error) {
	if matchableTcpGateway == nil {
		return nil, fmt.Errorf("matchableTcpGateway provided to Apply must not be nil")
	}
	patchOpts := opts.ToPatchOptions()
	data, err := json.Marshal(matchableTcpGateway)
	if err != nil {
		return nil, err
	}

	name := matchableTcpGateway.Name
	if name == nil {
		return nil, fmt.Errorf("matchableTcpGateway.Name must be provided to Apply")
	}

	result = &v1.MatchableTcpGateway{}
	err = c.client.Patch(types.ApplyPatchType).
		Namespace(c.ns).
		Resource("tcpgateways").
		Name(*name).
		SubResource("status").
		VersionedParams(&patchOpts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
