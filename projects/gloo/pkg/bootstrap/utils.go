package bootstrap

import (
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"path/filepath"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/crd"
	"github.com/solo-io/solo-kit/pkg/errors"
	"k8s.io/client-go/rest"
	"github.com/solo-io/solo-kit/pkg/utils/kubeutils"
	"k8s.io/client-go/kubernetes"
)

// sharedCache OR resourceCrd+cfg must be non-nil
func ConfigFactoryForSettings(settings *v1.Settings,
	sharedCache memory.InMemoryResourceCache,
	resourceCrd crd.Crd,
	cfg **rest.Config) (factory.ResourceClientFactory, error) {
	if settings.ConfigSource == nil {
		if sharedCache == nil {
			return nil, errors.Errorf("internal error: shared cache cannot be nil")
		}
		return &factory.MemoryResourceClientFactory{
			Cache: sharedCache,
		}, nil
	}

	switch source := settings.ConfigSource.(type) {
	// this is at trick to reuse the same cfg across multiple clients
	case *v1.Settings_KubernetesConfigSource:
		if cfg == nil {
			c, err := kubeutils.GetConfig("", "")
			if err != nil {
				return nil, err
			}
			*cfg = c
		}
		return &factory.KubeResourceClientFactory{
			Crd: resourceCrd,
			Cfg: *cfg,
		}, nil
	case *v1.Settings_DirectoryConfigSource:
		return &factory.FileResourceClientFactory{
			RootDir: filepath.Join(source.DirectoryConfigSource.Directory, resourceCrd.Plural),
		}, nil
	}
	return nil, errors.Errorf("invalid config source type")
}

// sharedCach OR resourceCrd+cfg must be non-nil
func SecretFactoryForSettings(pluralName string,
	settings *v1.Settings,
	sharedCache memory.InMemoryResourceCache,
	cfg **rest.Config,
	clientset *kubernetes.Interface) (factory.ResourceClientFactory, error) {
	if settings.SecretSource == nil {
		if sharedCache == nil {
			return nil, errors.Errorf("internal error: shared cache cannot be nil")
		}
		return &factory.MemoryResourceClientFactory{
			Cache: sharedCache,
		}, nil
	}

	switch source := settings.SecretSource.(type) {
	case *v1.Settings_KubernetesSecretSource:
		if cfg == nil {
			c, err := kubeutils.GetConfig("", "")
			if err != nil {
				return nil, err
			}
			*cfg = c
		}

		if clientset == nil {
			cs, err := kubernetes.NewForConfig(*cfg)
			if err != nil {
				return nil, err
			}
			*clientset = cs
		}
		return &factory.KubeSecretClientFactory{
			Clientset: clientset,
		}, nil
	case *v1.Settings_VaultSecretSource:
		return nil, errors.Errorf("vault configuration not implemented")
	case *v1.Settings_DirectorySecretSource:
		return &factory.FileResourceClientFactory{
			RootDir: filepath.Join(source.DirectorySecretSource.Directory, pluralName),
		}, nil
	}
	return nil, errors.Errorf("invalid config source type")
}
