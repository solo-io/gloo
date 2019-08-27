package proxysvc

import (
	"context"

	"github.com/solo-io/solo-projects/projects/grpcserver/server/helpers/status"

	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	v1 "github.com/solo-io/solo-projects/projects/grpcserver/api/v1"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/helpers/rawgetter"
	"go.uber.org/zap"
)

type proxyGrpcService struct {
	ctx             context.Context
	proxyClient     gloov1.ProxyClient
	rawGetter       rawgetter.RawGetter
	statusConverter status.InputResourceStatusGetter
}

func NewProxyGrpcService(ctx context.Context, proxyClient gloov1.ProxyClient, rawGetter rawgetter.RawGetter, statusConverter status.InputResourceStatusGetter) v1.ProxyApiServer {
	return &proxyGrpcService{
		ctx:             ctx,
		proxyClient:     proxyClient,
		rawGetter:       rawGetter,
		statusConverter: statusConverter,
	}
}

func (s *proxyGrpcService) GetProxy(ctx context.Context, request *v1.GetProxyRequest) (*v1.GetProxyResponse, error) {
	proxy, err := s.proxyClient.Read(request.GetRef().GetNamespace(), request.GetRef().GetName(), clients.ReadOpts{Ctx: s.ctx})
	if err != nil {
		wrapped := FailedToGetProxyError(err, request.GetRef())
		contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}
	return &v1.GetProxyResponse{ProxyDetails: s.getDetails(proxy)}, nil
}

func (s *proxyGrpcService) ListProxies(ctx context.Context, request *v1.ListProxiesRequest) (*v1.ListProxiesResponse, error) {
	var proxyDetailsList []*v1.ProxyDetails
	for _, ns := range request.GetNamespaces() {
		proxiesInNamespace, err := s.proxyClient.List(ns, clients.ListOpts{Ctx: s.ctx})
		if err != nil {
			wrapped := FailedToListProxiesError(err, ns)
			contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
			return nil, wrapped
		}
		for _, p := range proxiesInNamespace {
			proxyDetailsList = append(proxyDetailsList, s.getDetails(p))
		}
	}
	return &v1.ListProxiesResponse{ProxyDetails: proxyDetailsList}, nil
}

func (s *proxyGrpcService) getDetails(proxy *gloov1.Proxy) *v1.ProxyDetails {
	return &v1.ProxyDetails{
		Proxy:  proxy,
		Raw:    s.rawGetter.GetRaw(s.ctx, proxy, gloov1.ProxyCrd),
		Status: s.statusConverter.GetApiStatusFromResource(proxy),
	}
}
