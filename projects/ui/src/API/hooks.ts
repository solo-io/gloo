import { GatewayResourceApi } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/gateway_resources_pb_service';
import useSWR from 'swr';
import {
  VirtualService,
  Gateway,
  RouteTable,
} from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/gateway_resources_pb';
import { gatewayResourceApi } from './gateway-resources';
import { GlooResourceApi } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/gloo_resources_pb_service';
import {
  Upstream,
  UpstreamGroup,
  Proxy,
  Settings,
} from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/gloo_resources_pb';
import { glooResourceApi } from './gloo-resource';
import { FailoverSchemeApi } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/failover_scheme_pb_service';
import { FailoverScheme } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/failover_scheme_pb';
import { failoverSchemeApi } from './failover-scheme';
import {
  ObjectRef,
  ClusterObjectRef,
} from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';
import { GlooInstanceApi } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/glooinstance_pb_service';
import { glooInstanceApi } from './gloo-instance';
import {
  GlooInstance,
  ClusterDetails,
  ConfigDump,
  HostList,
} from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/glooinstance_pb';
import { federatedGatewayResourceApi } from './federated-gateway';
import { FederatedGatewayResourceApi } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/federated_gateway_resources_pb_service';
import {
  FederatedVirtualService,
  FederatedGateway,
  FederatedRouteTable,
} from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/federated_gateway_resources_pb';
import {
  FederatedUpstream,
  FederatedUpstreamGroup,
  FederatedSettings,
} from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/federated_gloo_resources_pb';
import { FederatedGlooResourceApi } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/federated_gloo_resources_pb_service';
import { federatedGlooResourceApi } from './federated-gloo';
import { FederatedAuthConfig } from '../proto/github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/federated_enterprise_gloo_resources_pb';
import { FederatedEnterpriseGlooResourceApi } from '../proto/github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/federated_enterprise_gloo_resources_pb_service';
import { FederatedRatelimitResourceApi } from '../proto/github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/federated_ratelimit_resources_pb_service';
import { federatedEnterpriseGlooResourceApi } from './federated-enterprise-gloo';
import { SubRouteTableRow } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/rt_selector_pb';
import { routeTablesSelectorApi } from './virtual-service-routes';
import {
  DescribeWasmFilterRequest,
  WasmFilter,
} from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/wasm_pb';
import { WasmFilterApi } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/wasm_pb_service';
import { wasmFilterApi } from './wasm-filter';
import { VirtualServiceRoutesApi } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/rt_selector_pb_service';
import { FederatedRateLimitConfig } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/federated_ratelimit_resources_pb';
import { BootstrapApi } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/bootstrap_pb_service';
import { bootstrapApi } from './bootstrap';
import { GlooFedCheckResponse } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/bootstrap_pb';
import { graphqlConfigApi } from './graphql';
import { GraphqlConfigApi } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/graphql_pb_service';
import { GraphqlApi } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/graphql_pb';
import { useState, useEffect } from 'react';
import { useParams } from 'react-router';

const normalRefreshInterval = 10000;

export function useListGlooInstances() {
  return useSWR<GlooInstance.AsObject[]>(
    GlooInstanceApi.ListGlooInstances.methodName,
    () => glooInstanceApi.listGlooInstances(),
    { refreshInterval: normalRefreshInterval }
  );
}
export function useListClusterDetails() {
  return useSWR<ClusterDetails.AsObject[]>(
    GlooInstanceApi.ListClusterDetails.methodName,
    () => glooInstanceApi.listClusterDetails(),
    { refreshInterval: normalRefreshInterval }
  );
}

export function useListVirtualServices(ref?: ObjectRef.AsObject) {
  let key = !!ref
    ? !!ref.name && !!ref.namespace
      ? `${GatewayResourceApi.ListVirtualServices.methodName}/${ref.namespace}/${ref.name}`
      : null
    : GatewayResourceApi.ListVirtualServices.methodName;

  return useSWR<VirtualService.AsObject[]>(
    key,
    () => gatewayResourceApi.listVirtualServices(ref),
    { refreshInterval: normalRefreshInterval }
  );
}

export function useGetSubroutesForVirtualService(
  ref?: ClusterObjectRef.AsObject
) {
  let key = !!ref
    ? `${VirtualServiceRoutesApi.GetVirtualServiceRoutes.methodName}/${ref.clusterName}/${ref.namespace}/${ref.name}`
    : null;

  return useSWR<SubRouteTableRow.AsObject[]>(
    key,
    () => routeTablesSelectorApi.getSubroutesForVirtualService(ref),
    { refreshInterval: normalRefreshInterval }
  );
}

export function useListRouteTables(ref?: ObjectRef.AsObject) {
  let key = !!ref
    ? !!ref.name && !!ref.namespace
      ? `${GatewayResourceApi.ListRouteTables.methodName}/${ref.namespace}/${ref.name}`
      : null
    : GatewayResourceApi.ListRouteTables.methodName;

  return useSWR<RouteTable.AsObject[]>(
    key,
    () => gatewayResourceApi.listRouteTables(ref),
    { refreshInterval: normalRefreshInterval }
  );
}

export function useListGateways(ref?: ObjectRef.AsObject) {
  let key = !!ref
    ? `${GatewayResourceApi.ListGateways.methodName}/${ref.namespace}/${ref.name}`
    : GatewayResourceApi.ListGateways.methodName;

  return useSWR<Gateway.AsObject[]>(
    key,
    () => gatewayResourceApi.listGateways(ref),
    { refreshInterval: normalRefreshInterval }
  );
}

export function useListSettings(ref?: ObjectRef.AsObject) {
  let key = !!ref
    ? `${GlooResourceApi.ListSettings.methodName}/${ref.namespace}/${ref.name}`
    : null;

  return useSWR<Settings.AsObject[]>(
    key,
    () => glooResourceApi.listSettings(ref),
    { refreshInterval: normalRefreshInterval }
  );
}

export function useListProxies(ref?: ObjectRef.AsObject) {
  let key = !!ref
    ? `${GlooResourceApi.ListProxies.methodName}/${ref.namespace}/${ref.name}`
    : null;

  return useSWR<Proxy.AsObject[]>(key, () => glooResourceApi.listProxies(ref), {
    refreshInterval: normalRefreshInterval,
  });
}

export function useListUpstreams(ref?: ObjectRef.AsObject) {
  let key = !!ref
    ? !!ref.name && !!ref.namespace
      ? `${GlooResourceApi.ListUpstreams.methodName}/${ref.namespace}/${ref.name}`
      : null
    : GlooResourceApi.ListUpstreams.methodName;

  return useSWR<Upstream.AsObject[]>(
    key,
    () => glooResourceApi.listUpstreams(ref),
    { refreshInterval: normalRefreshInterval }
  );
}

export function useGetUpstreamDetails(
  glooInstRef: ObjectRef.AsObject,
  upstreamRef: ClusterObjectRef.AsObject
) {
  const key = `get-${GlooResourceApi.ListUpstreams.methodName}/${glooInstRef.namespace}/${glooInstRef.name}/${upstreamRef.clusterName}/${upstreamRef.namespace}/${upstreamRef.name}`;

  return useSWR<Upstream.AsObject>(
    key,
    () => glooResourceApi.getUpstream(glooInstRef, upstreamRef),
    { refreshInterval: normalRefreshInterval }
  );
}

export function useGetUpstreamYaml(upstreamRef: ClusterObjectRef.AsObject) {
  const key = `${GlooResourceApi.GetUpstreamYaml.methodName}/${upstreamRef.clusterName}/${upstreamRef.namespace}/${upstreamRef.name}`;

  return useSWR<string>(
    key,
    () => glooResourceApi.getUpstreamYAML(upstreamRef),
    { refreshInterval: normalRefreshInterval }
  );
}

export function useGetUpstreamGroupYaml(
  upstreamGroupRef: ClusterObjectRef.AsObject
) {
  const key = `${GlooResourceApi.GetUpstreamGroupYaml.methodName}/${upstreamGroupRef.clusterName}/${upstreamGroupRef.namespace}/${upstreamGroupRef.name}`;

  return useSWR<string>(
    key,
    () => glooResourceApi.getUpstreamGroupYAML(upstreamGroupRef),
    { refreshInterval: normalRefreshInterval }
  );
}

export function useGetFailoverSchemeYaml(
  failoverSchemeRef: ObjectRef.AsObject
) {
  const key = `${FailoverSchemeApi.GetFailoverSchemeYaml.methodName}/${failoverSchemeRef.namespace}/${failoverSchemeRef.name}`;

  return useSWR<string>(
    key,
    () => failoverSchemeApi.getFailoverSchemeYAML(failoverSchemeRef),
    { refreshInterval: normalRefreshInterval }
  );
}

export function useGetUpstreamGroupDetails(
  glooInstRef: ObjectRef.AsObject,
  upstreamGroupRef: ClusterObjectRef.AsObject
) {
  const key = `get-${GlooResourceApi.ListUpstreamGroups.methodName}/${glooInstRef.namespace}/${glooInstRef.name}/${upstreamGroupRef.clusterName}/${upstreamGroupRef.namespace}/${upstreamGroupRef.name}`;

  return useSWR<UpstreamGroup.AsObject>(
    key,
    () => glooResourceApi.getUpstreamGroup(glooInstRef, upstreamGroupRef),
    { refreshInterval: normalRefreshInterval }
  );
}

export function useListUpstreamGroups(ref?: ObjectRef.AsObject) {
  let key = !!ref
    ? !!ref.name && !!ref.namespace
      ? `${GlooResourceApi.ListUpstreamGroups.methodName}/${ref.namespace}/${ref.name}`
      : null
    : GlooResourceApi.ListUpstreamGroups.methodName;

  return useSWR<UpstreamGroup.AsObject[]>(
    key,
    () => glooResourceApi.listUpstreamGroups(ref),
    { refreshInterval: normalRefreshInterval }
  );
}

export function useGetFailoverScheme(upstreamRef: ClusterObjectRef.AsObject) {
  const key = `${FailoverSchemeApi.GetFailoverScheme.methodName}/${upstreamRef.clusterName}/${upstreamRef.namespace}/${upstreamRef.name}`;

  return useSWR<FailoverScheme.AsObject>(
    key,
    () => failoverSchemeApi.getFailoverScheme(upstreamRef),
    { refreshInterval: normalRefreshInterval }
  );
}

/**
 * FEDERATED SECTION
 */
export function useListFederatedVirtualServices() {
  return useSWR<FederatedVirtualService.AsObject[]>(
    FederatedGatewayResourceApi.ListFederatedVirtualServices.methodName,
    () => federatedGatewayResourceApi.listFederatedVirtualServices(),
    { refreshInterval: normalRefreshInterval }
  );
}
export function useListFederatedGateways() {
  return useSWR<FederatedGateway.AsObject[]>(
    FederatedGatewayResourceApi.ListFederatedGateways.methodName,
    () => federatedGatewayResourceApi.listFederatedGateways(),
    { refreshInterval: normalRefreshInterval }
  );
}
export function useListFederatedRouteTables() {
  return useSWR<FederatedRouteTable.AsObject[]>(
    FederatedGatewayResourceApi.ListFederatedRouteTables.methodName,
    () => federatedGatewayResourceApi.listFederatedRouteTables(),
    { refreshInterval: normalRefreshInterval }
  );
}

export function useListFederatedUpstreams() {
  return useSWR<FederatedUpstream.AsObject[]>(
    FederatedGlooResourceApi.ListFederatedUpstreams.methodName,
    () => federatedGlooResourceApi.listFederatedUpstreams(),
    { refreshInterval: normalRefreshInterval }
  );
}
export function useListFederatedUpstreamGroups() {
  return useSWR<FederatedUpstreamGroup.AsObject[]>(
    FederatedGlooResourceApi.ListFederatedUpstreamGroups.methodName,
    () => federatedGlooResourceApi.listFederatedUpstreamGroups(),
    { refreshInterval: normalRefreshInterval }
  );
}
export function useListFederatedSettings() {
  return useSWR<FederatedSettings.AsObject[]>(
    FederatedGlooResourceApi.ListFederatedSettings.methodName,
    () => federatedGlooResourceApi.listFederatedSettings(),
    { refreshInterval: normalRefreshInterval }
  );
}
export function useListFederatedAuthConfigs() {
  return useSWR<FederatedAuthConfig.AsObject[]>(
    FederatedEnterpriseGlooResourceApi.ListFederatedAuthConfigs.methodName,
    () => federatedEnterpriseGlooResourceApi.listFederatedAuthConfigs(),
    { refreshInterval: normalRefreshInterval }
  );
}
export function useListFederatedRateLimits() {
  return useSWR<FederatedRateLimitConfig.AsObject[]>(
    FederatedRatelimitResourceApi.ListFederatedRateLimitConfigs.methodName,
    () => federatedEnterpriseGlooResourceApi.listFederatedRateLimitConfigs(),
    { refreshInterval: normalRefreshInterval }
  );
}

export function useGetWasmFilter(
  wasmFilterRequestRef: DescribeWasmFilterRequest.AsObject
) {
  let key = !!wasmFilterRequestRef?.gatewayRef
    ? `${WasmFilterApi.DescribeWasmFilter.methodName}/${wasmFilterRequestRef.gatewayRef.namespace}/${wasmFilterRequestRef.gatewayRef.name}/${wasmFilterRequestRef.name}`
    : null;

  return useSWR<WasmFilter.AsObject>(
    key,
    () => wasmFilterApi.getWasmFilter(wasmFilterRequestRef),
    { refreshInterval: normalRefreshInterval }
  );
}

export function useListWasmFilters() {
  return useSWR<WasmFilter.AsObject[]>(
    WasmFilterApi.ListWasmFilters.methodName,
    () => wasmFilterApi.listWasmFilters(),
    { refreshInterval: normalRefreshInterval }
  );
}

export function useGetGatewayYaml(
  gatewayClusterObjectRef: ClusterObjectRef.AsObject
) {
  let key = `${GatewayResourceApi.GetGatewayYaml.methodName}/${gatewayClusterObjectRef.namespace}/${gatewayClusterObjectRef.name}/${gatewayClusterObjectRef.clusterName}`;

  return useSWR<string>(
    key,
    () => gatewayResourceApi.getGatewayYAML(gatewayClusterObjectRef),
    {
      refreshInterval: normalRefreshInterval,
    }
  );
}

export function useGetConfigDumps(glooInstanceRef: ObjectRef.AsObject) {
  let key = `${GlooInstanceApi.GetConfigDumps.methodName}/${glooInstanceRef.namespace}/${glooInstanceRef.name}`;

  return useSWR<ConfigDump.AsObject[]>(
    key,
    () => glooInstanceApi.getConfigDumps(glooInstanceRef),
    {
      refreshInterval: normalRefreshInterval,
    }
  );
}

export function useGetUpstreamHosts(glooInstanceRef: ObjectRef.AsObject) {
  let key = `${GlooInstanceApi.GetUpstreamHosts.methodName}/${glooInstanceRef.namespace}/${glooInstanceRef.name}`;

  return useSWR<Map<string, HostList.AsObject>>(
    key,
    () => glooInstanceApi.getUpstreamHosts(glooInstanceRef),
    {
      refreshInterval: normalRefreshInterval,
    }
  );
}

export function useIsGlooFedEnabled() {
  return useSWR<GlooFedCheckResponse.AsObject>(
    BootstrapApi.IsGlooFedEnabled.methodName,
    () => bootstrapApi.isGlooFedEnabled(),
    { refreshInterval: normalRefreshInterval }
  );
}

export function useIsGraphqlEnabled() {
  return useSWR<boolean>(
    BootstrapApi.IsGraphqlEnabled.methodName,
    () => bootstrapApi.isGraphqlEnabled(),
    { refreshInterval: normalRefreshInterval }
  );
}

export const LIST_GRAPHQL_APIS_KEY = (glooInstanceRef?: ObjectRef.AsObject) =>
  !!glooInstanceRef
    ? `${GraphqlConfigApi.ListGraphqlApis.methodName}/${glooInstanceRef.namespace}/${glooInstanceRef.name}`
    : GraphqlConfigApi.ListGraphqlApis.methodName;
export function useListGraphqlApis(glooInstanceRef?: ObjectRef.AsObject) {
  return useSWR<GraphqlApi.AsObject[]>(
    GraphqlConfigApi.ListGraphqlApis.methodName,
    () => graphqlConfigApi.listGraphqlApis(glooInstanceRef),
    { refreshInterval: normalRefreshInterval }
  );
}

export function useGetGraphqlApiDetails(
  graphqlApiRef: ClusterObjectRef.AsObject
) {
  const key = `${GraphqlConfigApi.GetGraphqlApi.methodName}/${graphqlApiRef.namespace}/${graphqlApiRef.name}/${graphqlApiRef.clusterName}`;

  return useSWR<GraphqlApi.AsObject>(
    key,
    () => graphqlConfigApi.getGraphqlApi(graphqlApiRef),
    { refreshInterval: normalRefreshInterval }
  );
}

export function useGetGraphqlApiYaml(graphqlApiRef: ClusterObjectRef.AsObject) {
  const key = `${GraphqlConfigApi.GetGraphqlApiYaml.methodName}/${graphqlApiRef.namespace}/${graphqlApiRef.name}/${graphqlApiRef.clusterName}`;

  return useSWR<string>(
    key,
    () => graphqlConfigApi.getGraphqlApiYaml(graphqlApiRef),
    { refreshInterval: normalRefreshInterval }
  );
}

export function usePageGlooInstance() {
  // URL parameters (if on /apis/ then name='', namespace='')
  // Gets replaced by the default/initial gloo instance.
  const { name = '', namespace = '' } = useParams();
  const { data: glooInstances, error: instancesError } = useListGlooInstances();
  const [glooInstance, setGlooInstance] = useState<GlooInstance.AsObject>();
  useEffect(() => {
    if (!!glooInstances) {
      if (glooInstances.length === 1 && name == '' && namespace === '') {
        setGlooInstance(glooInstances[0]);
      } else {
        setGlooInstance(
          glooInstances.find(
            instance =>
              instance.metadata?.name === name &&
              instance.metadata?.namespace === namespace
          )
        );
      }
    } else {
      setGlooInstance(undefined);
    }
  }, [name, namespace, glooInstances]);
  return [glooInstance, glooInstances, instancesError] as [
    typeof glooInstance,
    typeof glooInstances,
    typeof instancesError
  ];
}
