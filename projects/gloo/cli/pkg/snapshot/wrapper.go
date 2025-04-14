package snapshot

import (
	"fmt"

	gatewaykube "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/apis/gateway.solo.io/v1"
	"github.com/solo-io/gloo/projects/gateway2/api/v1alpha1"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1/kube/apis/enterprise.gloo.solo.io/v1"
	glookube "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/kube/apis/gloo.solo.io/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
	apixv1a1 "sigs.k8s.io/gateway-api/apisx/v1alpha1"
)

type Wrapper interface {
	GetName() string
	GetNamespace() string
	Index() string
	GetLabels() map[string]string
	GetObjectKind() schema.ObjectKind
	FileOrigin() string
}

var _ Wrapper = &YAMLWrapper{}

type YAMLWrapper struct {
	runtime.Object
	fileOrigin string
	Yaml       string
}

func (w *YAMLWrapper) ObjectKind() schema.ObjectKind {
	if w.Object == nil {
		return nil
	}
	return w.Object.GetObjectKind()
}

func NewYAMLWrapper(obj runtime.Object, fileOrigin string) *YAMLWrapper {
	return &YAMLWrapper{
		Object:     obj,
		fileOrigin: fileOrigin,
	}
}

func (w *YAMLWrapper) GetLabels() map[string]string {
	return map[string]string{}
}
func (w *YAMLWrapper) GetName() string {
	return "unknown"
}
func (w *YAMLWrapper) GetNamespace() string {
	return "unknown"
}
func (w *YAMLWrapper) FileOrigin() string {
	return w.fileOrigin
}
func (w *YAMLWrapper) HasFileOrigin(fileOrigin string) *YAMLWrapper {
	w.fileOrigin = fileOrigin
	return w
}
func (w *YAMLWrapper) Index() string {
	return "unknown-unknown"
}

type HTTPRouteWrapper struct {
	*gwv1.HTTPRoute
	fileOrigin string
}

func NewHTTPRouteWrapper(httpRoute *gwv1.HTTPRoute, fileOrigin string) *HTTPRouteWrapper {
	return &HTTPRouteWrapper{
		HTTPRoute:  httpRoute,
		fileOrigin: fileOrigin,
	}
}

func (w *HTTPRouteWrapper) Index() string {
	return fmt.Sprintf("%s/%s", w.Namespace, w.Name)
}

func (w *HTTPRouteWrapper) FileOrigin() string {
	return w.fileOrigin
}
func (w *HTTPRouteWrapper) HasFileOrigin(fileOrigin string) *HTTPRouteWrapper {
	w.fileOrigin = fileOrigin
	return w
}

type SettingsWrapper struct {
	*glookube.Settings
	fileOrigin string
}

func NewSettingsWrapper(settings *glookube.Settings, fileOrigin string) *SettingsWrapper {
	return &SettingsWrapper{
		Settings:   settings,
		fileOrigin: fileOrigin,
	}
}

func (w *SettingsWrapper) Index() string {
	return fmt.Sprintf("%s/%s", w.Namespace, w.Name)
}
func (w *SettingsWrapper) FileOrigin() string {
	return w.fileOrigin
}
func (w *SettingsWrapper) HasFileOrigin(fileOrigin string) *SettingsWrapper {
	w.fileOrigin = fileOrigin
	return w
}

var _ Wrapper = &RouteOptionWrapper{}

type RouteOptionWrapper struct {
	*gatewaykube.RouteOption
	fileOrigin string
}

func NewRouteOptionWrapper(opt *gatewaykube.RouteOption, fileOrigin string) *RouteOptionWrapper {
	return &RouteOptionWrapper{
		RouteOption: opt,
		fileOrigin:  fileOrigin,
	}
}

func (w *RouteOptionWrapper) Index() string {
	return fmt.Sprintf("%s/%s", w.Namespace, w.Name)
}

func (w *RouteOptionWrapper) FileOrigin() string {
	return w.fileOrigin
}
func (w *RouteOptionWrapper) HasFileOrigin(fileOrigin string) *RouteOptionWrapper {
	w.fileOrigin = fileOrigin
	return w
}

type VirtualHostOptionWrapper struct {
	*gatewaykube.VirtualHostOption
	fileOrigin string
}

func NewVirtualHostOptionWrapper(opt *gatewaykube.VirtualHostOption, fileOrigin string) *VirtualHostOptionWrapper {
	return &VirtualHostOptionWrapper{
		VirtualHostOption: opt,
		fileOrigin:        fileOrigin,
	}
}

func (w *VirtualHostOptionWrapper) Index() string {
	return fmt.Sprintf("%s/%s", w.Namespace, w.Name)
}
func (w *VirtualHostOptionWrapper) FileOrigin() string {
	return w.fileOrigin
}
func (w *VirtualHostOptionWrapper) HasFileOrigin(fileOrigin string) *VirtualHostOptionWrapper {
	w.fileOrigin = fileOrigin
	return w
}

type GatewayParametersWrapper struct {
	*v1alpha1.GatewayParameters
	fileOrigin string
}

func NewGatewayParametersWrapper(params *v1alpha1.GatewayParameters, fileOrigin string) *GatewayParametersWrapper {
	return &GatewayParametersWrapper{
		GatewayParameters: params,
		fileOrigin:        fileOrigin,
	}
}

func (w *GatewayParametersWrapper) Index() string {
	return fmt.Sprintf("%s/%s", w.Namespace, w.Name)
}
func (w *GatewayParametersWrapper) FileOrigin() string {
	return w.fileOrigin
}
func (w *GatewayParametersWrapper) HasFileOrigin(fileOrigin string) *GatewayParametersWrapper {
	w.fileOrigin = fileOrigin
	return w
}

type ListenerOptionWrapper struct {
	*gatewaykube.ListenerOption
	fileOrigin string
}

func NewListenerOptionWrapper(opt *gatewaykube.ListenerOption, fileOrigin string) *ListenerOptionWrapper {
	return &ListenerOptionWrapper{
		ListenerOption: opt,
		fileOrigin:     fileOrigin,
	}
}

func (w *ListenerOptionWrapper) Index() string {
	return fmt.Sprintf("%s/%s", w.Namespace, w.Name)
}
func (w *ListenerOptionWrapper) FileOrigin() string {
	return w.fileOrigin
}
func (w *ListenerOptionWrapper) HasFileOrigin(fileOrigin string) *ListenerOptionWrapper {
	w.fileOrigin = fileOrigin
	return w
}

type HTTPListenerOptionWrapper struct {
	*gatewaykube.HttpListenerOption
	fileOrigin string
}

func NewHTTPListenerOptionWrapper(opt *gatewaykube.HttpListenerOption, fileOrigin string) *HTTPListenerOptionWrapper {
	return &HTTPListenerOptionWrapper{
		HttpListenerOption: opt,
		fileOrigin:         fileOrigin,
	}
}
func (w *HTTPListenerOptionWrapper) Index() string {
	return fmt.Sprintf("%s/%s", w.Namespace, w.Name)
}
func (w *HTTPListenerOptionWrapper) FileOrigin() string {
	return w.fileOrigin
}
func (w *HTTPListenerOptionWrapper) HasFileOrigin(fileOrigin string) *HTTPListenerOptionWrapper {
	w.fileOrigin = fileOrigin
	return w
}

type UpstreamWrapper struct {
	*glookube.Upstream
	fileOrigin string
}

func NewUpstreamWrapper(upstream *glookube.Upstream, fileOrigin string) *UpstreamWrapper {
	return &UpstreamWrapper{
		Upstream:   upstream,
		fileOrigin: fileOrigin,
	}
}

func (w *UpstreamWrapper) Index() string {
	return fmt.Sprintf("%s/%s", w.Namespace, w.Name)
}
func (w *UpstreamWrapper) FileOrigin() string {
	return w.fileOrigin
}

func (w *UpstreamWrapper) HasFileOrigin(fileOrigin string) *UpstreamWrapper {
	w.fileOrigin = fileOrigin
	return w
}

type AuthConfigWrapper struct {
	*v1.AuthConfig
	fileOrigin string
}

func NewAuthConfigWrapper(authConfig *v1.AuthConfig, fileOrigin string) *AuthConfigWrapper {
	return &AuthConfigWrapper{
		AuthConfig: authConfig,
		fileOrigin: fileOrigin,
	}
}

func (w *AuthConfigWrapper) Index() string {
	return fmt.Sprintf("%s/%s", w.Namespace, w.Name)
}
func (w *AuthConfigWrapper) FileOrigin() string {
	return w.fileOrigin
}
func (w *AuthConfigWrapper) HasFileOrigin(fileOrigin string) *AuthConfigWrapper {
	w.fileOrigin = fileOrigin
	return w
}

type GatewayWrapper struct {
	*gwv1.Gateway
	fileOrigin string
}

func NewGatewayWrapper(gw *gwv1.Gateway, fileOrigin string) *GatewayWrapper {
	return &GatewayWrapper{
		Gateway:    gw,
		fileOrigin: fileOrigin,
	}
}

func (w *GatewayWrapper) Index() string {
	return fmt.Sprintf("%s/%s", w.Namespace, w.Name)
}
func (w *GatewayWrapper) FileOrigin() string {
	return w.fileOrigin
}
func (w *GatewayWrapper) HasFileOrigin(fileOrigin string) *GatewayWrapper {
	w.fileOrigin = fileOrigin
	return w
}

type ListenerSetWrapper struct {
	*apixv1a1.XListenerSet
	fileOrigin string
}

func NewListenerSetWrapper(listenerSet *apixv1a1.XListenerSet, fileOrigin string) *ListenerSetWrapper {
	return &ListenerSetWrapper{
		XListenerSet: listenerSet,
		fileOrigin:   fileOrigin,
	}
}

func (w *ListenerSetWrapper) Index() string {
	return fmt.Sprintf("%s/%s", w.Namespace, w.Name)
}
func (w *ListenerSetWrapper) FileOrigin() string {
	return w.fileOrigin
}
func (w *ListenerSetWrapper) HasFileOrigin(fileOrigin string) *ListenerSetWrapper {
	w.fileOrigin = fileOrigin
	return w
}

type RouteTableWrapper struct {
	*gatewaykube.RouteTable
	fileOrigin string
}

func newRouteTableWrapper(routeTable *gatewaykube.RouteTable, fileOrigin string) *RouteTableWrapper {
	return &RouteTableWrapper{
		RouteTable: routeTable,
		fileOrigin: fileOrigin,
	}
}

func (w *RouteTableWrapper) Index() string {
	return fmt.Sprintf("%s/%s", w.Namespace, w.Name)
}
func (w *RouteTableWrapper) FileOrigin() string {
	return w.fileOrigin
}
func (w *RouteTableWrapper) HasFileOrigin(fileOrigin string) *RouteTableWrapper {
	w.fileOrigin = fileOrigin
	return w
}

var _ Wrapper = &VirtualServiceWrapper{}

type VirtualServiceWrapper struct {
	*gatewaykube.VirtualService
	fileOrigin string
}

func NewVirtualServiceWrapper(virtualService *gatewaykube.VirtualService, fileOrigin string) *VirtualServiceWrapper {
	return &VirtualServiceWrapper{
		VirtualService: virtualService,
		fileOrigin:     fileOrigin,
	}
}

func (w *VirtualServiceWrapper) Index() string {
	return fmt.Sprintf("%s/%s", w.GetNamespace(), w.GetName())
}
func (w *VirtualServiceWrapper) FileOrigin() string {
	return w.fileOrigin
}
func (w *VirtualServiceWrapper) HasFileOrigin(fileOrigin string) *VirtualServiceWrapper {
	w.fileOrigin = fileOrigin
	return w
}

var _ Wrapper = &GlooGatewayWrapper{}

type GlooGatewayWrapper struct {
	*gatewaykube.Gateway
	fileOrigin string
}

func NewGlooGatewayWrapper(gateway *gatewaykube.Gateway, fileOrigin string) *GlooGatewayWrapper {
	return &GlooGatewayWrapper{
		Gateway:    gateway,
		fileOrigin: fileOrigin,
	}
}

func (w *GlooGatewayWrapper) Index() string {
	return fmt.Sprintf("%s/%s", w.Namespace, w.Name)
}
func (w *GlooGatewayWrapper) FileOrigin() string {
	return w.fileOrigin
}
func (w *GlooGatewayWrapper) HasFileOrigin(fileOrigin string) *GlooGatewayWrapper {
	w.fileOrigin = fileOrigin
	return w
}

type DirectResponseWrapper struct {
	*v1alpha1.DirectResponse
	fileOrigin string
}

func NewDirectResponseWrapper(resp *v1alpha1.DirectResponse, fileOrigin string) *DirectResponseWrapper {
	return &DirectResponseWrapper{
		DirectResponse: resp,
		fileOrigin:     fileOrigin,
	}
}

func (w *DirectResponseWrapper) Index() string {
	return fmt.Sprintf("%s/%s", w.Namespace, w.Name)
}
func (w *DirectResponseWrapper) FileOrigin() string {
	return w.fileOrigin
}
func (w *DirectResponseWrapper) HasFileOrigin(fileOrigin string) *DirectResponseWrapper {
	w.fileOrigin = fileOrigin
	return w
}
