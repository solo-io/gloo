package wasmfilter_handler

import (
	"context"
	"sort"
	"strings"

	"google.golang.org/protobuf/encoding/prototext"

	"github.com/rotisserie/eris"
	"github.com/solo-io/go-utils/contextutils"
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/client"

	gateway_solo_io_v1 "github.com/solo-io/solo-apis/pkg/api/gateway.solo.io/v1"
	"github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1/options/wasm"
	rpc_edge_v1 "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/rpc.edge.gloo/v1"
	"github.com/solo-io/solo-projects/projects/apiserver/server/apiserverutils"
	fedv1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1"
)

func NewWasmFilterHandler(
	glooInstanceClient fedv1.GlooInstanceClient,
	mcGatewayCRDClientset gateway_solo_io_v1.MulticlusterClientset,
) rpc_edge_v1.WasmFilterApiServer {
	return &wasmFilterHandler{
		glooInstanceClient: glooInstanceClient,
		gatewayMCClientset: mcGatewayCRDClientset,
	}
}

type wasmFilterHandler struct {
	glooInstanceClient fedv1.GlooInstanceClient
	gatewayMCClientset gateway_solo_io_v1.MulticlusterClientset
}

func (k *wasmFilterHandler) DescribeWasmFilter(ctx context.Context, request *rpc_edge_v1.DescribeWasmFilterRequest) (*rpc_edge_v1.DescribeWasmFilterResponse, error) {
	if request.GetName() == "" || request.GetGatewayRef() == nil {
		return nil, eris.Errorf("invalid request, %v", request)
	}
	var wasmFilter *rpc_edge_v1.WasmFilter
	glooInstanceList, err := k.glooInstanceClient.ListGlooInstance(ctx)
	if err != nil {
		wrapped := eris.Wrapf(err, "Failed to get list gloo instances")
		contextutils.LoggerFrom(ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}
	for _, glooInstance := range glooInstanceList.Items {
		glooInstanceRef := apiserverutils.ToObjectRef(glooInstance.GetName(), glooInstance.GetNamespace())
		cluster := glooInstance.Spec.GetCluster()
		gatewayClient, err := k.gatewayMCClientset.Cluster(cluster)
		if err != nil {
			wrapped := eris.Wrapf(err, "Failed to get multicluster gateway client set")
			contextutils.LoggerFrom(ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
			return nil, wrapped
		}
		gateway, err := gatewayClient.Gateways().GetGateway(ctx, client.ObjectKey{
			Name:      request.GetGatewayRef().GetName(),
			Namespace: request.GetGatewayRef().GetNamespace(),
		})
		if err != nil {
			wrapped := eris.Wrapf(err, "Failed to get get gateway")
			contextutils.LoggerFrom(ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
			return nil, wrapped
		}
		gatewayRef := apiserverutils.ToClusterObjectRef(gateway.GetName(), gateway.GetNamespace(), cluster)
		for _, filter := range gateway.Spec.GetHttpGateway().GetOptions().GetWasm().GetFilters() {
			if request.GetName() == filter.GetName() {
				if wasmFilter == nil {
					wasmFilter = BuildRpcWasmFilter(filter,
						&rpc_edge_v1.WasmFilter_Location{
							GatewayRef:      &gatewayRef,
							GatewayStatus:   &gateway.Status,
							GlooInstanceRef: &glooInstanceRef,
						})
				} else {
					wasmFilter.Locations = append(wasmFilter.Locations, &rpc_edge_v1.WasmFilter_Location{
						GatewayRef:      &gatewayRef,
						GatewayStatus:   &gateway.Status,
						GlooInstanceRef: &glooInstanceRef,
					})
				}
			}
		}
	}
	return &rpc_edge_v1.DescribeWasmFilterResponse{
		WasmFilter: wasmFilter,
	}, nil
}

func (k *wasmFilterHandler) ListWasmFilters(ctx context.Context, request *rpc_edge_v1.ListWasmFiltersRequest) (*rpc_edge_v1.ListWasmFiltersResponse, error) {
	glooInstanceList, err := k.glooInstanceClient.ListGlooInstance(ctx)
	if err != nil {
		wrapped := eris.Wrapf(err, "Failed to get list gloo instances")
		contextutils.LoggerFrom(ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}
	wasmFilterMap := map[string]*rpc_edge_v1.WasmFilter{}
	for _, glooInstance := range glooInstanceList.Items {
		glooInstanceRef := apiserverutils.ToObjectRef(glooInstance.GetName(), glooInstance.GetNamespace())
		cluster := glooInstance.Spec.GetCluster()
		gatewayClient, err := k.gatewayMCClientset.Cluster(cluster)
		if err != nil {
			wrapped := eris.Wrapf(err, "Failed to get multicluster gateway client set")
			contextutils.LoggerFrom(ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
			return nil, wrapped
		}
		gatewayList, err := gatewayClient.Gateways().ListGateway(ctx)
		if err != nil {
			wrapped := eris.Wrapf(err, "Failed to get list gateways")
			contextutils.LoggerFrom(ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
			return nil, wrapped
		}
		for _, gateway := range gatewayList.Items {
			gatewayRef := apiserverutils.ToClusterObjectRef(gateway.GetName(), gateway.GetNamespace(), cluster)
			for _, filter := range gateway.Spec.GetHttpGateway().GetOptions().GetWasm().GetFilters() {
				filterKey := key(filter.GetName(), gateway.GetName(), gateway.GetNamespace())
				if existingRPCWasmFilter, ok := wasmFilterMap[filterKey]; ok {
					existingRPCWasmFilter.Locations = append(existingRPCWasmFilter.Locations, &rpc_edge_v1.WasmFilter_Location{
						GatewayRef:      &gatewayRef,
						GatewayStatus:   &gateway.Status,
						GlooInstanceRef: &glooInstanceRef,
					})
				} else {
					wasmFilterMap[filterKey] = BuildRpcWasmFilter(filter,
						&rpc_edge_v1.WasmFilter_Location{
							GatewayRef:      &gatewayRef,
							GatewayStatus:   &gateway.Status,
							GlooInstanceRef: &glooInstanceRef,
						})
				}
			}
		}
	}
	var wasmFilters []*rpc_edge_v1.WasmFilter
	for _, filter := range wasmFilterMap {
		wasmFilters = append(wasmFilters, filter)
	}
	sortWasmFilters(wasmFilters)

	return &rpc_edge_v1.ListWasmFiltersResponse{
		WasmFilters: wasmFilters,
	}, nil
}

func key(filterName, gatewayName, gatewayNamespace string) string {
	return strings.Join([]string{filterName, gatewayName, gatewayNamespace}, ".")
}

func sortWasmFilters(wasmFilters []*rpc_edge_v1.WasmFilter) {
	sort.Slice(wasmFilters, func(i, j int) bool {
		x := wasmFilters[i]
		y := wasmFilters[j]
		return x.GetName() < y.GetName()
	})
}

func BuildRpcWasmFilter(wasmFilter *wasm.WasmFilter, location *rpc_edge_v1.WasmFilter_Location) *rpc_edge_v1.WasmFilter {
	marshaler := &prototext.MarshalOptions{Indent: "  "}
	bytes, err := marshaler.Marshal(wasmFilter.GetConfig())
	if err != nil {
		return nil
	}
	configJson := string(bytes)

	source := wasmFilter.GetFilePath()
	if source == "" {
		source = wasmFilter.GetImage()
	}

	return &rpc_edge_v1.WasmFilter{
		Name:      wasmFilter.GetName(),
		RootId:    wasmFilter.GetRootId(),
		Source:    source,
		Config:    configJson,
		Locations: []*rpc_edge_v1.WasmFilter_Location{location},
	}
}
