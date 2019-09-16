package envoydetails

import (
	"context"
	"fmt"

	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"

	"github.com/solo-io/go-utils/contextutils"
	v1 "github.com/solo-io/solo-projects/projects/grpcserver/api/v1"
	"go.uber.org/zap"
	kubev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	podsGetter        corev1.PodsGetter
	httpGetter        HttpGetter
	proxyStatusGetter ProxyStatusGetter
}

var _ Client = &client{}

func NewClient(podsGetter corev1.PodsGetter, httpGetter HttpGetter, proxyStatusGetter ProxyStatusGetter) Client {
	return &client{
		podsGetter:        podsGetter,
		httpGetter:        httpGetter,
		proxyStatusGetter: proxyStatusGetter,
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
		details := c.getEnvoyDetails(ctx, pod)

		envoyDetailsList = append(envoyDetailsList, details)
	}

	return envoyDetailsList, nil
}

func (c *client) getEnvoyDetails(ctx context.Context, pod kubev1.Pod) *v1.EnvoyDetails {

	envoyName := getName(pod)
	dumpString, contentRenderError := c.getEnvoyConfig(ctx, pod)
	details := &v1.EnvoyDetails{
		Name: envoyName,
		Raw: &v1.Raw{
			FileName:           envoyName + ".json",
			Content:            dumpString,
			ContentRenderError: contentRenderError,
		},
		Status: c.proxyStatusGetter.GetProxyStatus(ctx, pod),
	}
	return details
}

func (c *client) getEnvoyConfig(ctx context.Context, pod kubev1.Pod) (string, string) {
	port, ok := pod.Annotations[ReadConfigPortAnnotationKey]
	if !ok {
		contextutils.LoggerFrom(ctx).Infow(
			fmt.Sprintf("missing admin port label for gateway proxy pod %v.%v, will not show config dump", pod.Namespace, pod.Name),
			zap.String("namespace", pod.Namespace),
			zap.String("name", pod.Name))
		return "", ""
	}

	path, ok := pod.Annotations[ReadConfigConfigDumpAnnotationKey]
	if !ok {
		contextutils.LoggerFrom(ctx).Infow(
			fmt.Sprintf("missing admin config_dump path label for gateway proxy pod %v.%v, will not show config dump", pod.Namespace, pod.Name),
			zap.String("namespace", pod.Namespace),
			zap.String("name", pod.Name))
		return "", ""
	}

	dumpString, err := c.httpGetter.Get(pod.Status.PodIP, port, path)
	if err != nil {
		contentRenderError := FailedToGetEnvoyConfig(pod.Namespace, pod.Name)
		contextutils.LoggerFrom(ctx).Warnw(contentRenderError,
			zap.Error(err),
			zap.String("namespace", pod.Namespace),
			zap.String("name", pod.Name))
		return "", contentRenderError
	}
	return dumpString, ""
}
