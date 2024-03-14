package utils

import (
	"io"
	"os"
	"strconv"

	errors "github.com/rotisserie/eris"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/kube/apis/gloo.solo.io/v1"
	glooinstancev1 "github.com/solo-io/solo-apis/pkg/api/fed.solo.io/v1"
	soloapisv1 "github.com/solo-io/solo-apis/pkg/api/gateway.solo.io/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"
	ctrlzap "sigs.k8s.io/controller-runtime/pkg/log/zap"
	v1 "sigs.k8s.io/gateway-api/apis/v1"
	"sigs.k8s.io/gateway-api/apis/v1alpha2"
	"sigs.k8s.io/gateway-api/apis/v1beta1"
)

// gateway apis uses this to build test examples: https://github.com/kubernetes-sigs/gateway-api/blob/main/pkg/test/cel/main_test.go#L57
func PtrTo[T any](a T) *T {
	return &a
}

// TODO(npolshak): Should this be different from classic mode to k8s gateway mode?
func BuildClientScheme() (*runtime.Scheme, error) {
	clientScheme := runtime.NewScheme()

	// k8s resources
	err := corev1.AddToScheme(clientScheme)
	if err != nil {
		return nil, err
	}
	err = appsv1.AddToScheme(clientScheme)
	if err != nil {
		return nil, err
	}
	// k8s gateway resources
	err = v1alpha2.AddToScheme(clientScheme)
	if err != nil {
		return nil, err
	}
	err = v1beta1.AddToScheme(clientScheme)
	if err != nil {
		return nil, err
	}
	err = v1.AddToScheme(clientScheme)
	if err != nil {
		return nil, err
	}
	// gloo resources
	err = glooinstancev1.AddToScheme(clientScheme)
	if err != nil {
		return nil, err
	}
	//err = gloogatewayv1.AddToScheme(clientScheme)
	//if err != nil {
	//	return nil, err
	//}
	err = gloov1.AddToScheme(clientScheme)
	if err != nil {
		return nil, err
	}
	err = soloapisv1.AddToScheme(clientScheme)
	if err != nil {
		return nil, err

	}

	return clientScheme, err
}

func GetClient(kubeCtx string) (client.Client, error) {
	clientScheme, err := BuildClientScheme()
	if err != nil {
		return client.Client(nil), err
	}

	restCfg, err := GetClientConfig(kubeCtx)
	if err != nil {
		return client.Client(nil), err
	}

	cluster, err := client.New(restCfg, client.Options{Scheme: clientScheme})
	if err != nil {
		return client.Client(nil), err
	}

	// Note: controller-runtime v0.15+ requires the global logger to be explicitly set
	ctrllog.SetLogger(ctrlzap.New(ctrlzap.WriteTo(io.Discard)))

	return cluster, nil
}

func GetClientConfig(kubeCtx string) (*rest.Config, error) {
	// Let's avoid defaulting, require clients to be explicit.
	if kubeCtx == "" {
		return nil, errors.New("missing cluster name")
	}

	cfg, err := clientcmd.NewDefaultClientConfigLoadingRules().Load()
	if err != nil {
		return nil, err
	}

	config := clientcmd.NewNonInteractiveClientConfig(*cfg, kubeCtx, &clientcmd.ConfigOverrides{}, nil)
	restCfg, err := config.ClientConfig()
	if err != nil {
		return nil, err
	}

	// Lets speed up our client when running tests
	restCfg.QPS = 50
	if v := os.Getenv("K8S_CLIENT_QPS"); v != "" {
		qps, err := strconv.ParseFloat(v, 32)
		if err != nil {
			return nil, err
		}
		restCfg.QPS = float32(qps)
	}

	restCfg.Burst = 100
	if v := os.Getenv("K8S_CLIENT_BURST"); v != "" {
		burst, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		restCfg.Burst = burst
	}

	return restCfg, nil
}
