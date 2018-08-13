package bootstrap

import (
	"k8s.io/client-go/kubernetes"
)

type Config interface {
	KubeClient() kubernetes.Interface
}
