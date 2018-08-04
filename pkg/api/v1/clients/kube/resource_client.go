package kube

import (
	"reflect"
	"sort"

	"github.com/gogo/protobuf/proto"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/crd"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/crd/client/clientset/versioned"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/protoutils"
	apiexts "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	kubewatch "k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/rest"
)

type ResourceClient struct {
	crd          crd.Crd
	apiexts      apiexts.Interface
	kube         versioned.Interface
	ownerLabel   string
	resourceName string
	resourceType resources.InputResource
}

func NewResourceClient(crd crd.Crd, cfg *rest.Config, resourceType resources.InputResource) (*ResourceClient, error) {
	apiExts, err := apiexts.NewForConfig(cfg)
	if err != nil {
		return nil, errors.Wrapf(err, "creating api extensions client")
	}
	crdClient, err := versioned.NewForConfig(cfg, crd)
	if err != nil {
		return nil, errors.Wrapf(err, "creating crd client")
	}
	return &ResourceClient{
		crd:          crd,
		apiexts:      apiExts,
		kube:         crdClient,
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
	return rc.crd.Register(rc.apiexts)
}

func (rc *ResourceClient) Read(namespace, name string, opts clients.ReadOpts) (resources.Resource, error) {
	if err := resources.ValidateName(name); err != nil {
		return nil, errors.Wrapf(err, "validation error")
	}
	opts = opts.WithDefaults()
	namespace = clients.DefaultNamespaceIfEmpty(namespace)

	resourceCrd, err := rc.kube.ResourcesV1().Resources(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil, errors.NewNotExistErr(namespace, name, err)
		}
		return nil, errors.Wrapf(err, "reading resource from kubernetes")
	}
	resource := rc.NewResource()
	if resourceCrd.Spec != nil {
		if err := protoutils.UnmarshalMap(*resourceCrd.Spec, resource); err != nil {
			return nil, errors.Wrapf(err, "reading crd spec into %v", rc.resourceName)
		}
	}
	resources.UpdateMetadata(resource, func(meta *core.Metadata) {
		meta.ResourceVersion = resourceCrd.ResourceVersion
	})
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
	clone := proto.Clone(resource).(resources.InputResource)
	clone.SetMetadata(meta)
	resourceCrd := rc.crd.KubeResource(clone)

	if rc.exist(meta.Namespace, meta.Name) {
		if !opts.OverwriteExisting {
			return nil, errors.NewExistErr(meta)
		}
		if _, err := rc.kube.ResourcesV1().Resources(meta.Namespace).Update(resourceCrd); err != nil {
			return nil, errors.Wrapf(err, "updating kube resource %v", resourceCrd.Name)
		}
	} else {
		if _, err := rc.kube.ResourcesV1().Resources(meta.Namespace).Create(resourceCrd); err != nil {
			return nil, errors.Wrapf(err, "creating kube resource %v", resourceCrd.Name)
		}
	}

	// return a read object to update the resource version
	return rc.Read(meta.Namespace, meta.Name, clients.ReadOpts{Ctx: opts.Ctx})
}

func (rc *ResourceClient) Delete(namespace, name string, opts clients.DeleteOpts) error {
	opts = opts.WithDefaults()
	if !rc.exist(namespace, name) {
		if !opts.IgnoreNotExist {
			return errors.NewNotExistErr(namespace, name)
		}
		return nil
	}

	if err := rc.kube.ResourcesV1().Resources(namespace).Delete(name, nil); err != nil {
		return errors.Wrapf(err, "deleting resource %v", name)
	}
	return nil
}

func (rc *ResourceClient) List(namespace string, opts clients.ListOpts) ([]resources.Resource, error) {
	opts = opts.WithDefaults()
	namespace = clients.DefaultNamespaceIfEmpty(namespace)

	resourceCrdList, err := rc.kube.ResourcesV1().Resources(namespace).List(metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(opts.Selector).String(),
	})
	if err != nil {
		return nil, errors.Wrapf(err, "listing resources in %v", namespace)
	}
	var resourceList []resources.Resource
	for _, resourceCrd := range resourceCrdList.Items {
		resource := rc.NewResource()
		if resourceCrd.Spec != nil {
			if err := protoutils.UnmarshalMap(*resourceCrd.Spec, resource); err != nil {
				return nil, errors.Wrapf(err, "reading crd spec into %v", rc.resourceName)
			}
		}
		resources.UpdateMetadata(resource, func(meta *core.Metadata) {
			meta.Namespace = resourceCrd.Namespace
			meta.Name = resourceCrd.Name
			meta.ResourceVersion = resourceCrd.ResourceVersion
		})
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
	watch, err := rc.kube.ResourcesV1().Resources(namespace).Watch(metav1.ListOptions{
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
			}
		}
	}()

	return resourcesChan, errs, nil
}

func (rc *ResourceClient) exist(namespace, name string) bool {
	_, err := rc.kube.ResourcesV1().Resources(namespace).Get(name, metav1.GetOptions{})
	return err == nil

}
