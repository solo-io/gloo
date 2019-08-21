package envoydetails

import (
	"context"
	"fmt"

	"github.com/solo-io/go-utils/contextutils"
	v1 "github.com/solo-io/solo-projects/projects/grpcserver/api/v1"
	"go.uber.org/zap"
	kubev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

//go:generate mockgen -destination mocks/mock_client.go -package mocks github.com/solo-io/solo-projects/projects/grpcserver/server/service/envoysvc/envoydetails Client
//go:generate mockgen -destination mocks/mock_pods_getter.go -package mocks k8s.io/client-go/kubernetes/typed/core/v1 PodsGetter,PodInterface

const (
	GatewayProxyLabelSelector = "gloo=gateway-proxy"
	// TODO export in gloo
	ReadConfigPortAnnotationKey       = "readconfig-port"
	ReadConfigConfigDumpAnnotationKey = "readconfig-config_dump"
	GatewayProxyIdLabel               = "gateway-proxy-id"
)

type Client interface {
	List(ctx context.Context, namespace string) ([]*v1.EnvoyDetails, error)
}

type client struct {
	// Use pods getter rather than namespaced pod interface
	// to support accessing pods across namespaces in the future
	podsGetter corev1.PodsGetter
	httpGetter HttpGetter
}

var _ Client = &client{}

func NewClient(podsGetter corev1.PodsGetter, httpGetter HttpGetter) Client {
	return &client{
		podsGetter: podsGetter,
		httpGetter: httpGetter,
	}
}

func (c *client) List(ctx context.Context, namespace string) ([]*v1.EnvoyDetails, error) {
	podList, err := c.podsGetter.Pods(namespace).List(metav1.ListOptions{LabelSelector: GatewayProxyLabelSelector})
	if err != nil {
		wrapped := FailedToListPodsError(err, namespace, GatewayProxyLabelSelector)
		contextutils.LoggerFrom(ctx).Errorw(wrapped.Error(),
			zap.Error(err),
			zap.String("namespace", namespace),
			zap.String("labelSelector", GatewayProxyLabelSelector))
		return nil, wrapped
	}

	envoyDetailsList := make([]*v1.EnvoyDetails, 0, len(podList.Items))
	for _, pod := range podList.Items {
		port, ok := pod.Annotations[ReadConfigPortAnnotationKey]
		if !ok {
			contextutils.LoggerFrom(ctx).Infow(
				fmt.Sprintf("Missing admin port label for gateway proxy pod %v.%v, excluding from details list", pod.Namespace, pod.Name),
				zap.String("namespace", pod.Namespace),
				zap.String("name", pod.Name))
			continue
		}

		path, ok := pod.Annotations[ReadConfigConfigDumpAnnotationKey]
		if !ok {
			contextutils.LoggerFrom(ctx).Infow(
				fmt.Sprintf("Missing admin config_dump path label for gateway proxy pod %v.%v, excluding from details list", pod.Namespace, pod.Name),
				zap.String("namespace", pod.Namespace),
				zap.String("name", pod.Name))
			continue
		}

		dumpString, err := c.httpGetter.Get(pod.Status.PodIP, port, path)
		var contentRenderError string
		if err != nil {
			contentRenderError = FailedToGetEnvoyConfig(pod.Namespace, pod.Name)
			contextutils.LoggerFrom(ctx).Warnw(contentRenderError,
				zap.Error(err),
				zap.String("namespace", pod.Namespace),
				zap.String("name", pod.Name))
		}

		envoyName := getName(pod)
		details := &v1.EnvoyDetails{
			Name: envoyName,
			Raw: &v1.Raw{
				FileName:           envoyName + ".json",
				Content:            dumpString,
				ContentRenderError: contentRenderError,
			},
		}
		envoyDetailsList = append(envoyDetailsList, details)
	}

	return envoyDetailsList, nil
}

func getName(pod kubev1.Pod) string {
	if id, ok := pod.Labels[GatewayProxyIdLabel]; ok {
		return id
	}
	return pod.Name
}
