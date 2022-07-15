package statuses

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	errors "github.com/rotisserie/eris"
	gwv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/api/v2/reporter"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func NewHTTPReporter(
	statusBuilder reporter.StatusBuilder,
	statusClient resources.StatusClient,
	statusChan StatusReportChan,
) *httpReporter {
	return &httpReporter{
		statusBuilder: statusBuilder,
		statusClient:  statusClient,
		statusChan:    statusChan,
	}
}

type httpReporter struct {
	statusBuilder reporter.StatusBuilder

	statusChan StatusReportChan

	statusClient resources.StatusClient
}

func (m *httpReporter) WriteReports(
	ctx context.Context,
	errs reporter.ResourceReports,
	subresourceStatuses reporter.SubResourceStatuses,
) error {
	resourceByKind := make(StatusReports)
	for resource, report := range errs {
		status := m.statusBuilder.StatusFromReport(report, subresourceStatuses[resource])
		gvk, err := resourceToGVK(resource)
		if err != nil {
			contextutils.LoggerFrom(ctx).Error(err)
			continue
		}
		kind := gvkToString(gvk)
		if _, ok := resourceByKind[kind]; !ok {
			resourceByKind[kind] = make(map[string]*core.Status)
		}
		resourceByKind[kind][m.buildKey(resource)] = status
	}

	select {
	case <-ctx.Done():
		return nil
	case <-time.After(time.Millisecond * 100):
		return nil
	case m.statusChan <- resourceByKind:
		return nil
	}
}

func (m *httpReporter) buildKey(resource resources.InputResource) string {
	return fmt.Sprintf("%s.%s", resource.GetMetadata().GetNamespace(), resource.GetMetadata().GetName())
}

type StatusReports map[string]map[string]*core.Status
type StatusReportChan chan StatusReports

func NewStatusHandler(ctx context.Context) (*statusHandler, StatusReportChan) {
	statusChan := make(StatusReportChan, 5)
	handler := &statusHandler{
		statuses:   StatusReports{},
		statusChan: statusChan,
	}
	go handler.receiveStatusesForever(ctx)
	return handler, statusChan

}

type statusHandler struct {
	lock sync.RWMutex
	// map of KIND -> namespace.name -> *core.Status
	statuses map[string]map[string]*core.Status
	// Chan used to send data to this handler to display
	statusChan StatusReportChan
}

func (m *statusHandler) receiveStatusesForever(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
		case newMap, ok := <-m.statusChan:
			if !ok {
				return
			}
			m.lock.Lock()
			for key, val := range newMap {
				m.statuses[key] = val
			}
			m.lock.Unlock()
		}
	}
}

func (m *statusHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	byt, err := m.getData(request.Context(), request.URL.Query())
	if err != nil {
		writer.Write([]byte(err.Error()))
		return
	}
	writer.Write(byt)
}

func (m *statusHandler) getData(ctx context.Context, data url.Values) ([]byte, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	contextutils.LoggerFrom(ctx).Debugf("url values %+v", data)
	if !data.Has("gvk") {
		return json.Marshal(m.statuses)
	}

	gvkParam := strings.ToLower(data.Get("gvk"))

	gvkData, ok := m.statuses[gvkParam]
	if !ok {
		return nil, errors.Errorf("No data available for GVK: %s", gvkParam)
	}

	if !data.Has("ref") {
		return json.Marshal(gvkData)
	}

	result := map[string]*core.Status{}
	refs := strings.Split(data.Get("ref"), ",")
	for _, ref := range refs {
		status, ok := gvkData[ref]
		if !ok {
			continue
		}
		result[ref] = status
	}
	return json.Marshal(result)
}

func gvkToString(kind schema.GroupVersionKind) string {
	return fmt.Sprintf("%s.%s", strings.ToLower(kind.Kind), kind.GroupVersion().String())
}

func resourceToGVK(resource resources.Resource) (schema.GroupVersionKind, error) {
	switch resource.(type) {
	// Gateway resources
	case *gwv1.Gateway:
		return gwv1.GatewayGVK, nil
	case *gwv1.RouteTable:
		return gwv1.RouteTableGVK, nil
	case *gwv1.VirtualService:
		return gwv1.VirtualServiceGVK, nil
	// Gloo resources
	case *gloov1.Proxy:
		return gloov1.ProxyGVK, nil
	case *gloov1.Secret:
		return gloov1.SecretGVK, nil
	case *gloov1.Upstream:
		return gloov1.UpstreamGVK, nil
	case *gloov1.UpstreamGroup:
		return gloov1.UpstreamGroupGVK, nil
	default:
		return schema.GroupVersionKind{}, errors.Errorf("status reporting is not supported for resource type: %T", resource)
	}
}
