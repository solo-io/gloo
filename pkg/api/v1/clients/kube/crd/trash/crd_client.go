package trash

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/rest"
)

var Scheme = runtime.NewScheme()
var Codecs = serializer.NewCodecFactory(Scheme)
var ParameterCodec = runtime.NewParameterCodec(Scheme)

type crdResult struct {
	Resource   resources.Resource
	ObjectKind schema.ObjectKind
}

func (r *crdResult) GetObjectKind() schema.ObjectKind {
	return r.ObjectKind
}

func (r *crdResult) DeepCopyObject() runtime.Object {
	return &crdResult{
		Resource:   resources.Clone(r.Resource),
		ObjectKind: r.ObjectKind,
	}
}

// ResourcesGetter has a method to return a ResourceInterface.
// A group's client should implement this interface.
type ResourcesGetter interface {
	Resources(namespace string) ResourceInterface
}

// ResourceInterface has methods to work with Resource resources.
type ResourceInterface interface {
	Create(resources.Resource) (resources.Resource, error)
	Update(resources.Resource) (resources.Resource, error)
	Delete(name string, options *meta_v1.DeleteOptions) error
	DeleteCollection(options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error
	Get(name string, options meta_v1.GetOptions) (resources.Resource, error)
	List(opts meta_v1.ListOptions) ([]resources.Resource, error)
	Watch(opts meta_v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result resources.Resource, err error)
}

// rscs implements ResourceInterface
type rscs struct {
	client       rest.Interface
	ns           string
	crdType      Crd
	resourceType resources.Resource
}

// newResources returns a Resources
func NewResources(c rest.Interface, namespace string, crdType Crd, resourceType resources.Resource) *rscs {
	return &rscs{
		client:       c,
		ns:           namespace,
		crdType:      crdType,
		resourceType: resourceType,
	}
}

func (c *rscs) newResult() *crdResult {
	return &crdResult{
		Resource:   resources.Clone(c.resourceType),
		ObjectKind: c.crdType.Type.GetObjectKind(),
	}
}

// Get takes name of the rsc, and returns the corresponding rsc object, and an error if there is any.
func (c *rscs) Get(name string, options meta_v1.GetOptions) (result resources.Resource, err error) {
	res := c.newResult()
	err = c.client.Get().
		Namespace(c.ns).
		Resource("rscs").
		Name(name).
		VersionedParams(&options, ParameterCodec).
		Do().
		Into(res)
	result = res.Resource
	return
}

// List takes label and field selectors, and returns the list of Resources that match those selectors.
func (c *rscs) List(opts meta_v1.ListOptions) (result []resources.Resource, err error) {
	res := c.newResult()
	err = c.client.Get().
		Namespace(c.ns).
		Resource("rscs").
		VersionedParams(&opts, ParameterCodec).
		Do().
		Into(res)
	result = res.Resource
	return
}

// Watch returns a watch.Interface that watches the requested rscs.
func (c *rscs) Watch(opts meta_v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("rscs").
		VersionedParams(&opts, ParameterCodec).
		Watch()
}

// Create takes the representation of a rsc and creates it.  Returns the server's representation of the rsc, and an error, if there is any.
func (c *rscs) Create(rsc resources.Resource) (result resources.Resource, err error) {
	res := c.newResult()
	err = c.client.Post().
		Namespace(c.ns).
		Resource("rscs").
		Body(rsc).
		Do().
		Into(res)
	result = res.Resource
	return
}

// Update takes the representation of a rsc and updates it. Returns the server's representation of the rsc, and an error, if there is any.
func (c *rscs) Update(rsc resources.Resource) (result resources.Resource, err error) {
	res := c.newResult()
	err = c.client.Put().
		Namespace(c.ns).
		Resource("rscs").
		Name(rsc.GetMetadata().Name).
		Body(rsc).
		Do().
		Into(res)
	result = res.Resource
	return
}

// Delete takes name of the rsc and deletes it. Returns an error if one occurs.
func (c *rscs) Delete(name string, options *meta_v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("rscs").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *rscs) DeleteCollection(options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("rscs").
		VersionedParams(&listOptions, ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched rsc.
func (c *rscs) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result resources.Resource, err error) {
	res := c.newResult()
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("rscs").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(res)
	result = res.Resource
	return
}
