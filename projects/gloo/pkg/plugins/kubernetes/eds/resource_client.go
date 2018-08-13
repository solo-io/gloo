package eds

import (
	"reflect"
	"sort"

	"github.com/gogo/protobuf/proto"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/errors"
	"k8s.io/api/core/v1"
	apiexts "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	kubewatch "k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

// TODO: change this becomes a proto
// Upstream Spec for Kubernetes Upstreams
// Kubernetes Upstreams represent a set of one or more addressable pods for a Kubernetes Service
// the Gloo Kubernetes Upstream maps to a single service port. Because Kubernetes Services support mulitple ports,
// Gloo requires that a different upstream be created for each port
// Kubernetes Upstreams are typically generated automatically by Gloo from the Kubernetes API
type UpstreamSpec struct {
	// The name of the Kubernetes Service
	ServiceName string `protobuf:"bytes,1,opt,name=service_name,json=serviceName,proto3" json:"service_name,omitempty"`
	// The namespace where the Service lives
	ServiceNamespace string `protobuf:"bytes,2,opt,name=service_namespace,json=serviceNamespace,proto3" json:"service_namespace,omitempty"`
	// The port where the Service is listening. If the service only has one port, this can be left empty
	ServicePort int32 `protobuf:"varint,3,opt,name=service_port,json=servicePort,proto3" json:"service_port,omitempty"`
	// Labels allow finer-grained filtering of pods for the Upstream. Gloo will select pods based on their labels if
	// any are provided here.
	// (see [Kubernetes labels and selectors](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/)
	Labels map[string]string `protobuf:"bytes,4,rep,name=labels" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func fromKubeMeta(meta metav1.ObjectMeta) core.Metadata {
	return core.Metadata{
		Name:            meta.Name,
		Namespace:       meta.Namespace,
		ResourceVersion: meta.ResourceVersion,
		Labels:          meta.Labels,
		Annotations:     meta.Annotations,
	}
}

func toKubeMeta(meta core.Metadata) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:            meta.Name,
		Namespace:       clients.DefaultNamespaceIfEmpty(meta.Namespace),
		ResourceVersion: meta.ResourceVersion,
		Labels:          meta.Labels,
		Annotations:     meta.Annotations,
	}
}

func fromKubeService(svc *v1.Service) *UpstreamSpec {
	into.SetMetadata(fromKubeMeta(svc.ObjectMeta))
	data := make(map[string]string)
	for k, v := range svc.Data {
		data[k] = string(v)
	}
	into.SetData(data)
}

func toKubeSecret(resource resources.DataResource) *v1.Secret {
	data := make(map[string][]byte)
	for k, v := range resource.GetData() {
		data[k] = []byte(v)
	}
	return &v1.Secret{
		ObjectMeta: toKubeMeta(resource.GetMetadata()),
		Data:       data,
	}
}

type ResourceClient struct {
	apiexts      apiexts.Interface
	kube         kubernetes.Interface
	ownerLabel   string
	resourceName string
	resourceType resources.DataResource
}

func NewResourceClient(kube kubernetes.Interface, resourceType resources.DataResource) (*ResourceClient, error) {
	return &ResourceClient{
		kube:         kube,
		resourceName: reflect.TypeOf(resourceType).Name(),
		resourceType: resourceType,
	}, nil
}

var _ clients.ResourceClient = &ResourceClient{}

func (rc *ResourceClient) Kind() string {
	return resources.Kind(rc.resourceType)
}

func (rc *ResourceClient) NewResource() resources.Resource {
	return resources.Clone(rc.resourceType)
}

func (rc *ResourceClient) Register() error {
	return nil
}

func (rc *ResourceClient) Read(namespace, name string, opts clients.ReadOpts) (resources.Resource, error) {
	if err := resources.ValidateName(name); err != nil {
		return nil, errors.Wrapf(err, "validation error")
	}
	opts = opts.WithDefaults()
	namespace = clients.DefaultNamespaceIfEmpty(namespace)

	secret, err := rc.kube.CoreV1().Secrets(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil, errors.NewNotExistErr(namespace, name, err)
		}
		return nil, errors.Wrapf(err, "reading secret from kubernetes")
	}
	resource := rc.NewResource().(resources.DataResource)
	fromKubeService(secret, resource)
	return resource, nil
}

func (rc *ResourceClient) Write(resource resources.Resource, opts clients.WriteOpts) (resources.Resource, error) {
	opts = opts.WithDefaults()
	if err := resources.Validate(resource); err != nil {
		return nil, errors.Wrapf(err, "validation error")
	}
	meta := resource.GetMetadata()
	meta.Namespace = clients.DefaultNamespaceIfEmpty(meta.Namespace)

	// mutate and return clone
	clone := proto.Clone(resource).(resources.Resource)
	clone.SetMetadata(meta)
	secret := toKubeSecret(resource.(resources.DataResource))

	original, err := rc.Read(meta.Namespace, meta.Name, clients.ReadOpts{
		Ctx: opts.Ctx,
	})
	if original != nil && err == nil {
		if !opts.OverwriteExisting {
			return nil, errors.NewExistErr(meta)
		}
		if meta.ResourceVersion != original.GetMetadata().ResourceVersion {
			return nil, errors.Errorf("resource version error. must update new resource version to match current")
		}
		if _, err := rc.kube.CoreV1().Secrets(secret.Namespace).Update(secret); err != nil {
			return nil, errors.Wrapf(err, "updating kube secret %v", secret.Name)
		}
	} else {
		if _, err := rc.kube.CoreV1().Secrets(secret.Namespace).Create(secret); err != nil {
			return nil, errors.Wrapf(err, "creating kube secret %v", secret.Name)
		}
	}

	// return a read object to update the resource version
	return rc.Read(secret.Namespace, secret.Name, clients.ReadOpts{Ctx: opts.Ctx})
}

func (rc *ResourceClient) Delete(namespace, name string, opts clients.DeleteOpts) error {
	opts = opts.WithDefaults()
	if !rc.exist(namespace, name) {
		if !opts.IgnoreNotExist {
			return errors.NewNotExistErr(namespace, name)
		}
		return nil
	}

	if err := rc.kube.CoreV1().Secrets(namespace).Delete(name, nil); err != nil {
		return errors.Wrapf(err, "deleting secret %v", name)
	}
	return nil
}

func (rc *ResourceClient) List(namespace string, opts clients.ListOpts) ([]resources.Resource, error) {
	opts = opts.WithDefaults()
	namespace = clients.DefaultNamespaceIfEmpty(namespace)

	secretList, err := rc.kube.CoreV1().Secrets(namespace).List(metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(opts.Selector).String(),
	})
	if err != nil {
		return nil, errors.Wrapf(err, "listing secrets in %v", namespace)
	}
	var resourceList []resources.Resource
	for _, secret := range secretList.Items {
		resource := rc.NewResource().(resources.DataResource)
		fromKubeService(&secret, resource)
		resourceList = append(resourceList, resource)
	}

	sort.SliceStable(resourceList, func(i, j int) bool {
		return resourceList[i].GetMetadata().Name < resourceList[j].GetMetadata().Name
	})

	return resourceList, nil
}

func (rc *ResourceClient) Watch(namespace string, opts clients.WatchOpts) (<-chan []resources.Resource, <-chan error, error) {
	opts = opts.WithDefaults()
	namespace = clients.DefaultNamespaceIfEmpty(namespace)
	watch, err := rc.kube.CoreV1().Secrets(namespace).Watch(metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(opts.Selector).String(),
	})
	if err != nil {
		return nil, nil, errors.Wrapf(err, "initiating kube watch in %v", namespace)
	}
	resourcesChan := make(chan []resources.Resource)
	errs := make(chan error)
	updateResourceList := func() {
		list, err := rc.List(namespace, clients.ListOpts{
			Ctx:      opts.Ctx,
			Selector: opts.Selector,
		})
		if err != nil {
			errs <- err
			return
		}
		resourcesChan <- list
	}
	// watch should open up with an initial read
	go updateResourceList()

	go func() {
		for {
			select {
			case event := <-watch.ResultChan():
				switch event.Type {
				case kubewatch.Error:
					errs <- errors.Errorf("error during watch: %v", event)
				default:
					updateResourceList()
				}
			case <-opts.Ctx.Done():
				watch.Stop()
				close(resourcesChan)
				close(errs)
				return
			}
		}
	}()

	return resourcesChan, errs, nil
}

func (rc *ResourceClient) exist(namespace, name string) bool {
	_, err := rc.kube.CoreV1().Secrets(namespace).Get(name, metav1.GetOptions{})
	return err == nil
}
