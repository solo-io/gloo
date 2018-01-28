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

package v1

import (
	scheme "github.com/solo-io/glue/config/watcher/crd/client/clientset/versioned/scheme"
	v1 "github.com/solo-io/glue/config/watcher/crd/solo.io/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// VirtualHostsGetter has a method to return a VirtualHostInterface.
// A group's client should implement this interface.
type VirtualHostsGetter interface {
	VirtualHosts(namespace string) VirtualHostInterface
}

// VirtualHostInterface has methods to work with VirtualHost resources.
type VirtualHostInterface interface {
	Create(*v1.VirtualHost) (*v1.VirtualHost, error)
	Update(*v1.VirtualHost) (*v1.VirtualHost, error)
	Delete(name string, options *meta_v1.DeleteOptions) error
	DeleteCollection(options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error
	Get(name string, options meta_v1.GetOptions) (*v1.VirtualHost, error)
	List(opts meta_v1.ListOptions) (*v1.VirtualHostList, error)
	Watch(opts meta_v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.VirtualHost, err error)
	VirtualHostExpansion
}

// virtualHosts implements VirtualHostInterface
type virtualHosts struct {
	client rest.Interface
	ns     string
}

// newVirtualHosts returns a VirtualHosts
func newVirtualHosts(c *GlueV1Client, namespace string) *virtualHosts {
	return &virtualHosts{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the virtualHost, and returns the corresponding virtualHost object, and an error if there is any.
func (c *virtualHosts) Get(name string, options meta_v1.GetOptions) (result *v1.VirtualHost, err error) {
	result = &v1.VirtualHost{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("virtualhosts").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of VirtualHosts that match those selectors.
func (c *virtualHosts) List(opts meta_v1.ListOptions) (result *v1.VirtualHostList, err error) {
	result = &v1.VirtualHostList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("virtualhosts").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested virtualHosts.
func (c *virtualHosts) Watch(opts meta_v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("virtualhosts").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a virtualHost and creates it.  Returns the server's representation of the virtualHost, and an error, if there is any.
func (c *virtualHosts) Create(virtualHost *v1.VirtualHost) (result *v1.VirtualHost, err error) {
	result = &v1.VirtualHost{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("virtualhosts").
		Body(virtualHost).
		Do().
		Into(result)
	return
}

// Update takes the representation of a virtualHost and updates it. Returns the server's representation of the virtualHost, and an error, if there is any.
func (c *virtualHosts) Update(virtualHost *v1.VirtualHost) (result *v1.VirtualHost, err error) {
	result = &v1.VirtualHost{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("virtualhosts").
		Name(virtualHost.Name).
		Body(virtualHost).
		Do().
		Into(result)
	return
}

// Delete takes name of the virtualHost and deletes it. Returns an error if one occurs.
func (c *virtualHosts) Delete(name string, options *meta_v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("virtualhosts").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *virtualHosts) DeleteCollection(options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("virtualhosts").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched virtualHost.
func (c *virtualHosts) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.VirtualHost, err error) {
	result = &v1.VirtualHost{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("virtualhosts").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
