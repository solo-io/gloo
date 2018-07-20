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
	kubewatch "k8s.io/apimachinery/pkg/watch"
)

type ResourceClient struct {
	crd          crd.Crd
	apiexts      apiexts.Interface
	kube         versioned.Interface
	ownerLabel   string
	resourceName string
	resourceType resources.Resource
}

func NewResourceClient(crd crd.Crd, apiexts apiexts.Interface, kube versioned.Interface, resourceType resources.Resource) *ResourceClient {
	return &ResourceClient{
		crd:          crd,
		apiexts:      apiexts,
		kube:         kube,
		resourceName: reflect.TypeOf(resourceType).Name(),
		resourceType: resourceType,
	}
}

var _ clients.ResourceClient = &ResourceClient{}

func (rc *ResourceClient) Register() error {
	return rc.crd.Register(rc.apiexts)
}

func (rc *ResourceClient) Read(name string, opts clients.ReadOpts) (resources.Resource, error) {
	if err := resources.ValidateName(name); err != nil {
		return nil, errors.Wrapf(err, "validation error")
	}
	opts = opts.WithDefaults()

	resourceCrd, err := rc.kube.ResourcesV1().Resources(opts.Namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil, errors.NewNotExistErr(opts.Namespace, name, err)
		}
		return nil, errors.Wrapf(err, "reading resource from kubernetes")
	}
	resource := resources.Clone(rc.resourceType)
	if resourceCrd.Spec != nil {
		if err := protoutils.UnmarshalMap(*resourceCrd.Spec, resource); err != nil {
			return nil, errors.Wrapf(err, "reading crd spec into %v", rc.resourceName)
		}
	}
	resources.UpdateMetadata(resource, func(meta core.Metadata) core.Metadata {
		meta.ResourceVersion = resourceCrd.ResourceVersion
		return meta
	})
	return resource, nil
}

func (rc *ResourceClient) Write(resource resources.Resource, opts clients.WriteOpts) (resources.Resource, error) {
	opts = opts.WithDefaults()
	if err := resources.Validate(resource); err != nil {
		return nil, errors.Wrapf(err, "validation error")
	}
	meta := resource.GetMetadata()
	if meta.Namespace == "" {
		meta.Namespace = clients.DefaultNamespace
	}

	// mutate and return clone
	clone := proto.Clone(resource).(resources.Resource)
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
	return rc.Read(meta.Name, clients.ReadOpts{Ctx: opts.Ctx, Namespace: meta.Namespace})
}

func (rc *ResourceClient) Delete(name string, opts clients.DeleteOpts) error {
	opts = opts.WithDefaults()
	if !rc.exist(opts.Namespace, name) {
		if !opts.IgnoreNotExist {
			return errors.NewNotExistErr(opts.Namespace, name)
		}
		return nil
	}

	if err := rc.kube.ResourcesV1().Resources(opts.Namespace).Delete(name, nil); err != nil {
		return errors.Wrapf(err, "deleting resource %v", name)
	}
	return nil
}

func (rc *ResourceClient) List(opts clients.ListOpts) ([]resources.Resource, error) {
	opts = opts.WithDefaults()

	resourceCrdList, err := rc.kube.ResourcesV1().Resources(opts.Namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, errors.Wrapf(err, "deleting listing resources in %v", opts.Namespace)
	}
	var resourceList []resources.Resource
	for _, resourceCrd := range resourceCrdList.Items {
		resource := resources.Clone(rc.resourceType)
		if resourceCrd.Spec != nil {
			if err := protoutils.UnmarshalMap(*resourceCrd.Spec, resource); err != nil {
				return nil, errors.Wrapf(err, "reading crd spec into %v", rc.resourceName)
			}
		}
		resources.UpdateMetadata(resource, func(meta core.Metadata) core.Metadata {
			meta.Namespace = resourceCrd.Namespace
			meta.Name = resourceCrd.Name
			meta.ResourceVersion = resourceCrd.ResourceVersion
			return meta
		})
		resourceList = append(resourceList, resource)
	}

	sort.SliceStable(resourceList, func(i, j int) bool {
		return resourceList[i].GetMetadata().Name < resourceList[j].GetMetadata().Name
	})

	return resourceList, nil
}

func (rc *ResourceClient) Watch(opts clients.WatchOpts) (<-chan []resources.Resource, <-chan error, error) {
	opts = opts.WithDefaults()
	watch, err := rc.kube.ResourcesV1().Resources(opts.Namespace).Watch(metav1.ListOptions{})
	if err != nil {
		return nil, nil, errors.Wrapf(err, "initiating kube watch in %v", opts.Namespace)
	}
	resourcesChan := make(chan []resources.Resource)
	errs := make(chan error)
	updateResourceList := func() {
		list, err := rc.List(clients.ListOpts{
			Ctx:       opts.Ctx,
			Selector:  opts.Selector,
			Namespace: opts.Namespace,
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
