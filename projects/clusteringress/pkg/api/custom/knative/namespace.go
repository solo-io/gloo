package knative

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
)

var _ resources.ResourceNamespaceLister = &knativeResourceNamespaceLister{}

func KnativeResourceNamespaceLister() resources.ResourceNamespaceLister {
	return &knativeResourceNamespaceLister{}
}

type knativeResourceNamespaceLister struct {
}

func (kns *knativeResourceNamespaceLister) GetResourceNamespaceList(opts resources.ResourceNamespaceListOptions, filtered resources.ResourceNamespaceList) (resources.ResourceNamespaceList, error) {
	panic("the client does not support Getting Resource Namespace Lists")
}

func (kns *knativeResourceNamespaceLister) GetResourceNamespaceWatch(opts resources.ResourceNamespaceWatchOptions, filtered resources.ResourceNamespaceList) (chan resources.ResourceNamespaceList, <-chan error, error) {
	panic("the client does not support Getting Resource Namespace Watch")
}
