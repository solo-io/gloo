package apiserverutils

import (
	"sort"

	"github.com/solo-io/skv2/contrib/pkg/sets"
	rpc_edge_v1 "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/rpc.edge.gloo/v1"
	fedv1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1"
	fedv1types "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1/types"
)

// Converts a Fed GlooInstance into an Apiserver GlooInstance
func ConvertToRpcGlooInstance(glooInstance *fedv1.GlooInstance) *rpc_edge_v1.GlooInstance {
	spec := glooInstance.Spec
	rpcSpec := &rpc_edge_v1.GlooInstance_GlooInstanceSpec{
		Cluster:      spec.GetCluster(),
		IsEnterprise: spec.GetIsEnterprise(),
		ControlPlane: &rpc_edge_v1.GlooInstance_GlooInstanceSpec_ControlPlane{
			Version:           spec.GetControlPlane().GetVersion(),
			Namespace:         spec.GetControlPlane().GetNamespace(),
			WatchedNamespaces: spec.GetControlPlane().GetWatchedNamespaces(),
		},
		Proxies: convertToRpcProxies(spec.GetProxies()),
		Region:  spec.GetRegion(),
		Check:   convertToRpcCheck(spec.GetCheck()),
	}
	rpcStatus := &rpc_edge_v1.GlooInstance_GlooInstanceStatus{}

	return &rpc_edge_v1.GlooInstance{
		Metadata: ToMetadata(glooInstance.ObjectMeta),
		Spec:     rpcSpec,
		Status:   rpcStatus,
	}
}

func convertToRpcProxies(proxies []*fedv1types.GlooInstanceSpec_Proxy) []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Proxy {
	rpcProxies := make([]*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Proxy, 0, len(proxies))
	for _, proxy := range proxies {
		rpcProxies = append(rpcProxies, &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Proxy{
			Replicas:                      proxy.GetReplicas(),
			AvailableReplicas:             proxy.GetAvailableReplicas(),
			ReadyReplicas:                 proxy.GetReadyReplicas(),
			WasmEnabled:                   proxy.GetWasmEnabled(),
			ReadConfigMulticlusterEnabled: proxy.GetReadConfigMulticlusterEnabled(),
			Version:                       proxy.GetVersion(),
			Name:                          proxy.GetName(),
			Namespace:                     proxy.GetNamespace(),
			WorkloadControllerType:        convertToRpcWorkloadControllerType(proxy.GetWorkloadControllerType()),
			Zones:                         proxy.GetZones(),
			IngressEndpoints:              convertToRpcIngressEndpoints(proxy.GetIngressEndpoints()),
		})
	}
	return rpcProxies
}

func convertToRpcWorkloadControllerType(workloadControllerType fedv1types.GlooInstanceSpec_Proxy_WorkloadController) rpc_edge_v1.GlooInstance_GlooInstanceSpec_Proxy_WorkloadController {
	switch workloadControllerType {
	case fedv1types.GlooInstanceSpec_Proxy_DEPLOYMENT:
		return rpc_edge_v1.GlooInstance_GlooInstanceSpec_Proxy_DEPLOYMENT
	case fedv1types.GlooInstanceSpec_Proxy_DAEMON_SET:
		return rpc_edge_v1.GlooInstance_GlooInstanceSpec_Proxy_DAEMON_SET
	}
	return rpc_edge_v1.GlooInstance_GlooInstanceSpec_Proxy_UNDEFINED
}

func convertToRpcIngressEndpoints(ingressEndpoints []*fedv1types.GlooInstanceSpec_Proxy_IngressEndpoint) []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Proxy_IngressEndpoint {
	rpcIngressEndpoints := make([]*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Proxy_IngressEndpoint, 0, len(ingressEndpoints))
	for _, ingressEndpoint := range ingressEndpoints {
		rpcIngressEndpoints = append(rpcIngressEndpoints, &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Proxy_IngressEndpoint{
			Address:     ingressEndpoint.GetAddress(),
			Ports:       convertToRpcPorts(ingressEndpoint.GetPorts()),
			ServiceName: ingressEndpoint.GetServiceName(),
		})
	}
	return rpcIngressEndpoints
}

func convertToRpcPorts(ports []*fedv1types.GlooInstanceSpec_Proxy_IngressEndpoint_Port) []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Proxy_IngressEndpoint_Port {
	rpcPorts := make([]*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Proxy_IngressEndpoint_Port, 0, len(ports))
	for _, port := range ports {
		rpcPorts = append(rpcPorts, &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Proxy_IngressEndpoint_Port{
			Port: port.GetPort(),
			Name: port.GetName(),
		})
	}
	return rpcPorts
}

func convertToRpcCheck(check *fedv1types.GlooInstanceSpec_Check) *rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check {
	return &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check{
		Gateways:              convertToRpcSummary(check.GetGateways()),
		MatchableHttpGateways: convertToRpcSummary(check.GetMatchableHttpGateways()),
		RateLimitConfigs:      convertToRpcSummary(check.GetRateLimitConfigs()),
		VirtualServices:       convertToRpcSummary(check.GetVirtualServices()),
		RouteTables:           convertToRpcSummary(check.GetRouteTables()),
		AuthConfigs:           convertToRpcSummary(check.GetAuthConfigs()),
		Settings:              convertToRpcSummary(check.GetSettings()),
		Upstreams:             convertToRpcSummary(check.GetUpstreams()),
		UpstreamGroups:        convertToRpcSummary(check.GetUpstreamGroups()),
		Proxies:               convertToRpcSummary(check.GetProxies()),
		Deployments:           convertToRpcSummary(check.GetDeployments()),
		Pods:                  convertToRpcSummary(check.GetPods()),
	}
}

func convertToRpcSummary(summary *fedv1types.GlooInstanceSpec_Check_Summary) *rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary {
	return &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary{
		Total:    summary.GetTotal(),
		Errors:   convertToRpcResourceReports(summary.GetErrors()),
		Warnings: convertToRpcResourceReports(summary.GetWarnings()),
	}
}

func convertToRpcResourceReports(resourceReports []*fedv1types.GlooInstanceSpec_Check_Summary_ResourceReport) []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport {
	rpcResourceReports := make([]*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport, 0, len(resourceReports))
	for _, resourceReport := range resourceReports {
		rpcResourceReports = append(rpcResourceReports, &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{
			Ref:     resourceReport.GetRef(),
			Message: resourceReport.GetMessage(),
		})
	}
	return rpcResourceReports
}

// Performs the same function as `projects/gloo-fed/pkg/discovery/translator/summarize/sort.go`
// but does not depend on Fed APIs
func SortCheckSummaryLists(summary *rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary) {
	sort.SliceStable(summary.Errors, func(i, j int) bool {
		return sets.Key(summary.Errors[i].GetRef()) < sets.Key(summary.Errors[j].GetRef())
	})

	sort.SliceStable(summary.Warnings, func(i, j int) bool {
		return sets.Key(summary.Warnings[i].GetRef()) < sets.Key(summary.Warnings[j].GetRef())
	})
}
