package kube

import (
	"sort"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/thirdparty"
	"github.com/solo-io/solo-kit/pkg/errors"
	"k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	kubewatch "k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

type configMapClient struct {
	kube kubernetes.Interface
}

func NewArtifactClient(kube kubernetes.Interface) thirdparty.ThirdPartyResourceClient {
	return &configMapClient{
		kube: kube,
	}
}

func fromKubeConfigMap(configMap *v1.ConfigMap) *thirdparty.Artifact {
	values := make(map[string]string)
	for k, v := range configMap.BinaryData {
		values[k] = string(v)
	}
	return &thirdparty.Artifact{
		Data: thirdparty.Data{
			Metadata: fromKubeMeta(configMap.ObjectMeta),
			Values:   values,
		},
	}
}

func toKubeConfigMap(resource thirdparty.ThirdPartyResource) *v1.ConfigMap {
	data := make(map[string][]byte)
	for k, v := range resource.GetData() {
		data[k] = []byte(v)
	}
	return &v1.ConfigMap{
		ObjectMeta: toKubeMeta(resource.GetMetadata()),
		BinaryData: data,
	}
}

func (rc *configMapClient) Read(namespace, name string, opts clients.ReadOpts) (thirdparty.ThirdPartyResource, error) {
	if err := resources.ValidateName(name); err != nil {
		return nil, errors.Wrapf(err, "validation error")
	}
	opts = opts.WithDefaults()
	namespace = clients.DefaultNamespaceIfEmpty(namespace)

	configMap, err := rc.kube.CoreV1().ConfigMaps(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil, errors.NewNotExistErr(namespace, name, err)
		}
		return nil, errors.Wrapf(err, "reading configMap from kubernetes")
	}
	return fromKubeConfigMap(configMap), nil
}

func (rc *configMapClient) Write(resource thirdparty.ThirdPartyResource, opts clients.WriteOpts) (thirdparty.ThirdPartyResource, error) {
	opts = opts.WithDefaults()
	if err := resources.ValidateName(resource.GetMetadata().Name); err != nil {
		return nil, errors.Wrapf(err, "validation error")
	}
	configMap := toKubeConfigMap(resource)

	if rc.exist(configMap.Namespace, configMap.Name) {
		if !opts.OverwriteExisting {
			return nil, errors.NewExistErr(resource.GetMetadata())
		}
		if _, err := rc.kube.CoreV1().ConfigMaps(configMap.Namespace).Update(configMap); err != nil {
			return nil, errors.Wrapf(err, "updating kube configMap %v", configMap.Name)
		}
	} else {
		if _, err := rc.kube.CoreV1().ConfigMaps(configMap.Namespace).Create(configMap); err != nil {
			return nil, errors.Wrapf(err, "creating kube configMap %v", configMap.Name)
		}
	}

	// return a read object to update the resource version
	return rc.Read(configMap.Namespace, configMap.Name, clients.ReadOpts{Ctx: opts.Ctx})
}

func (rc *configMapClient) Delete(namespace, name string, opts clients.DeleteOpts) error {
	opts = opts.WithDefaults()
	if !rc.exist(namespace, name) {
		if !opts.IgnoreNotExist {
			return errors.NewNotExistErr(namespace, name)
		}
		return nil
	}

	if err := rc.kube.CoreV1().ConfigMaps(namespace).Delete(name, nil); err != nil {
		return errors.Wrapf(err, "deleting configMap %v", name)
	}
	return nil
}

func (rc *configMapClient) List(namespace string, opts clients.ListOpts) ([]thirdparty.ThirdPartyResource, error) {
	opts = opts.WithDefaults()
	namespace = clients.DefaultNamespaceIfEmpty(namespace)

	configMapList, err := rc.kube.CoreV1().ConfigMaps(namespace).List(metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(opts.Selector).String(),
	})
	if err != nil {
		return nil, errors.Wrapf(err, "listing configMaps in %v", namespace)
	}
	var resourceList []thirdparty.ThirdPartyResource
	for _, configMap := range configMapList.Items {
		resourceList = append(resourceList, fromKubeConfigMap(&configMap))
	}

	sort.SliceStable(resourceList, func(i, j int) bool {
		return resourceList[i].GetMetadata().Name < resourceList[j].GetMetadata().Name
	})

	return resourceList, nil
}

func (rc *configMapClient) Watch(namespace string, opts clients.WatchOpts) (<-chan []thirdparty.ThirdPartyResource, <-chan error, error) {
	opts = opts.WithDefaults()
	namespace = clients.DefaultNamespaceIfEmpty(namespace)
	watch, err := rc.kube.CoreV1().ConfigMaps(namespace).Watch(metav1.ListOptions{})
	if err != nil {
		return nil, nil, errors.Wrapf(err, "initiating kube watch in %v", namespace)
	}
	resourcesChan := make(chan []thirdparty.ThirdPartyResource)
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

func (rc *configMapClient) exist(namespace, name string) bool {
	_, err := rc.kube.CoreV1().ConfigMaps(namespace).Get(name, metav1.GetOptions{})
	return err == nil

}
