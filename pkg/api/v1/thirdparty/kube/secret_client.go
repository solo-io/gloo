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

type secretClient struct {
	kube kubernetes.Interface
}

func NewSecretClient(kube kubernetes.Interface) thirdparty.ThirdPartyResourceClient {
	return &secretClient{
		kube: kube,
	}
}

func fromKube(secret *v1.Secret) *thirdparty.Secret {
	values := make(map[string]string)
	for k, v := range secret.Data {
		values[k] = string(v)
	}
	return &thirdparty.Secret{
		Data: thirdparty.Data{
			Metadata: fromKubeMeta(secret.ObjectMeta),
			Values:   values,
		},
	}
}

func toKube(secret thirdparty.ThirdPartyResource) *v1.Secret {
	data := make(map[string][]byte)
	for k, v := range secret.GetData() {
		data[k] = []byte(v)
	}
	return &v1.Secret{
		ObjectMeta: toKubeMeta(secret.GetMetadata()),
		Data:       data,
	}
}

func (rc *secretClient) Read(namespace, name string, opts clients.ReadOpts) (thirdparty.ThirdPartyResource, error) {
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
	return fromKube(secret), nil
}

func (rc *secretClient) Write(resource thirdparty.ThirdPartyResource, opts clients.WriteOpts) (thirdparty.ThirdPartyResource, error) {
	opts = opts.WithDefaults()
	if err := resources.ValidateName(resource.GetMetadata().Name); err != nil {
		return nil, errors.Wrapf(err, "validation error")
	}
	secret := toKube(resource)

	if rc.exist(secret.Namespace, secret.Name) {
		if !opts.OverwriteExisting {
			return nil, errors.NewExistErr(resource.GetMetadata())
		}
		if _, err := rc.kube.CoreV1().Secrets(secret.Namespace).Update(secret); err != nil {
			return nil, errors.Wrapf(err, "updating kube resource %v", secret.Name)
		}
	} else {
		if _, err := rc.kube.CoreV1().Secrets(secret.Namespace).Create(secret); err != nil {
			return nil, errors.Wrapf(err, "creating kube resource %v", secret.Name)
		}
	}

	// return a read object to update the resource version
	return rc.Read(secret.Namespace, secret.Name, clients.ReadOpts{Ctx: opts.Ctx})
}

func (rc *secretClient) Delete(namespace, name string, opts clients.DeleteOpts) error {
	opts = opts.WithDefaults()
	if !rc.exist(namespace, name) {
		if !opts.IgnoreNotExist {
			return errors.NewNotExistErr(namespace, name)
		}
		return nil
	}

	if err := rc.kube.CoreV1().Secrets(namespace).Delete(name, nil); err != nil {
		return errors.Wrapf(err, "deleting resource %v", name)
	}
	return nil
}

func (rc *secretClient) List(namespace string, opts clients.ListOpts) ([]thirdparty.ThirdPartyResource, error) {
	opts = opts.WithDefaults()
	namespace = clients.DefaultNamespaceIfEmpty(namespace)

	secretList, err := rc.kube.CoreV1().Secrets(namespace).List(metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(opts.Selector).String(),
	})
	if err != nil {
		return nil, errors.Wrapf(err, "listing resources in %v", namespace)
	}
	var resourceList []thirdparty.ThirdPartyResource
	for _, secret := range secretList.Items {
		resourceList = append(resourceList, fromKube(&secret))
	}

	sort.SliceStable(resourceList, func(i, j int) bool {
		return resourceList[i].GetMetadata().Name < resourceList[j].GetMetadata().Name
	})

	return resourceList, nil
}

func (rc *secretClient) Watch(namespace string, opts clients.WatchOpts) (<-chan []thirdparty.ThirdPartyResource, <-chan error, error) {
	opts = opts.WithDefaults()
	namespace = clients.DefaultNamespaceIfEmpty(namespace)
	watch, err := rc.kube.CoreV1().Secrets(namespace).Watch(metav1.ListOptions{})
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

func (rc *secretClient) exist(namespace, name string) bool {
	_, err := rc.kube.CoreV1().Secrets(namespace).Get(name, metav1.GetOptions{})
	return err == nil

}
