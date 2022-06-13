import {
  ClusterObjectRef,
  ObjectRef,
} from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';
import { GlooInstance } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/glooinstance_pb';
import { DescribeWasmFilterRequest } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/wasm_pb';
import { useEffect, useState } from 'react';
import { useParams } from 'react-router';
import { bootstrapApi } from './bootstrap';
import { failoverSchemeApi } from './failover-scheme';
import { federatedEnterpriseGlooResourceApi } from './federated-enterprise-gloo';
import { federatedGatewayResourceApi } from './federated-gateway';
import { federatedGlooResourceApi } from './federated-gloo';
import { gatewayResourceApi } from './gateway-resources';
import { glooInstanceApi } from './gloo-instance';
import { glooResourceApi } from './gloo-resource';
import { graphqlConfigApi } from './graphql';
import { useRequest } from './helpers';
import { routeTablesSelectorApi } from './virtual-service-routes';
import { wasmFilterApi } from './wasm-filter';

// -- glooInstanceApi
export const useListGlooInstances = () =>
  useRequest(glooInstanceApi.listGlooInstances, []);
export const useListClusterDetails = () =>
  useRequest(glooInstanceApi.listClusterDetails, []);
export const useGetConfigDumps = (glooInstanceRef: ObjectRef.AsObject) =>
  useRequest(glooInstanceApi.getConfigDumps, [glooInstanceRef]);
export const useGetUpstreamHosts = (glooInstanceRef: ObjectRef.AsObject) =>
  useRequest(glooInstanceApi.getUpstreamHosts, [glooInstanceRef]);

// -- gatewayResourceApi
export const useListVirtualServices = (
  ...args: Parameters<typeof gatewayResourceApi.listVirtualServices>
) => useRequest(gatewayResourceApi.listVirtualServices, args);
export const useListRouteTables = (
  ...args: Parameters<typeof gatewayResourceApi.listRouteTables>
) => useRequest(gatewayResourceApi.listRouteTables, args);
export const useGetGatewayYaml = (
  gatewayClusterObjectRef: ClusterObjectRef.AsObject
) => useRequest(gatewayResourceApi.getGatewayYAML, [gatewayClusterObjectRef]);

// -- routeTablesSelectorApi
export const useGetSubroutesForVirtualService = (
  ref?: ClusterObjectRef.AsObject
) => useRequest(routeTablesSelectorApi.getSubroutesForVirtualService, [ref]);

// -- gatewayResourceApi
export const useListGateways = (ref?: ObjectRef.AsObject) =>
  useRequest(gatewayResourceApi.listGateways, [ref]);

// -- glooResourceApi
export const useListSettings = (ref?: ObjectRef.AsObject) =>
  useRequest(glooResourceApi.listSettings, [ref]);
export const useListProxies = (ref?: ObjectRef.AsObject) =>
  useRequest(glooResourceApi.listProxies, [ref]);
export const useListUpstreams = (
  ...args: Parameters<typeof glooResourceApi.listUpstreams>
) => useRequest(glooResourceApi.listUpstreams, args);
export const useListUpstreamGroups = (
  ...args: Parameters<typeof glooResourceApi.listUpstreamGroups>
) => useRequest(glooResourceApi.listUpstreamGroups, args);

export const useGetUpstreamYaml = (upstreamRef: ClusterObjectRef.AsObject) =>
  useRequest(glooResourceApi.getUpstreamYAML, [upstreamRef]);
export const useGetUpstreamGroupYaml = (
  upstreamGroupRef: ClusterObjectRef.AsObject
) => useRequest(glooResourceApi.getUpstreamGroupYAML, [upstreamGroupRef]);

export const useGetUpstreamDetails = (upstreamRef: ClusterObjectRef.AsObject) =>
  useRequest(glooResourceApi.getUpstreamDetails, [upstreamRef]);
export const useGetUpstreamGroupDetails = (
  upstreamGroupRef: ClusterObjectRef.AsObject
) => useRequest(glooResourceApi.getUpstreamGroupDetails, [upstreamGroupRef]);

// -- failoverSchemeApi
export const useGetFailoverSchemeYaml = (
  failoverSchemeRef: ObjectRef.AsObject
) => useRequest(failoverSchemeApi.getFailoverSchemeYAML, [failoverSchemeRef]);
export const useGetFailoverScheme = (upstreamRef: ClusterObjectRef.AsObject) =>
  useRequest(failoverSchemeApi.getFailoverScheme, [upstreamRef]);

// -- wasmFilterApi
export const useGetWasmFilter = (
  wasmFilterRequestRef: DescribeWasmFilterRequest.AsObject
) => useRequest(wasmFilterApi.getWasmFilter, [wasmFilterRequestRef]);
export const useListWasmFilters = () =>
  useRequest(wasmFilterApi.listWasmFilters, []);

// -- bootstrapApi
export const useIsGlooFedEnabled = () =>
  useRequest(bootstrapApi.isGlooFedEnabled, []);
export const useIsGraphqlEnabled = () =>
  useRequest(bootstrapApi.isGraphqlEnabled, []);

// -- graphqlConfigApi
export const useListGraphqlApis = (glooInstanceRef?: ObjectRef.AsObject) =>
  useRequest(graphqlConfigApi.listGraphqlApis, [glooInstanceRef]);
export const useGetGraphqlApiDetails = (
  graphqlApiRef: ClusterObjectRef.AsObject
) => useRequest(graphqlConfigApi.getGraphqlApi, [graphqlApiRef]);
export const useGetGraphqlApiYaml = (
  graphqlApiRef: ClusterObjectRef.AsObject
) => useRequest(graphqlConfigApi.getGraphqlApiYaml, [graphqlApiRef]);
export const useGetStitchedSchemaDefinition = (
  graphqlApiRef: ClusterObjectRef.AsObject
) => useRequest(graphqlConfigApi.getStitchedSchemaDefinition, [graphqlApiRef]);

// ------------------------ //
//
// Federated Start
//
// -- federatedGatewayResourceApi
export const useListFederatedVirtualServices = () =>
  useRequest(federatedGatewayResourceApi.listFederatedVirtualServices, []);
export const useListFederatedGateways = () =>
  useRequest(federatedGatewayResourceApi.listFederatedGateways, []);
export const useListFederatedRouteTables = () =>
  useRequest(federatedGatewayResourceApi.listFederatedRouteTables, []);
//
// -- federatedGlooResourceApi
export const useListFederatedUpstreams = () =>
  useRequest(federatedGlooResourceApi.listFederatedUpstreams, []);
export const useListFederatedUpstreamGroups = () =>
  useRequest(federatedGlooResourceApi.listFederatedUpstreamGroups, []);
export const useListFederatedSettings = () =>
  useRequest(federatedGlooResourceApi.listFederatedSettings, []);
//
// -- federatedEnterpriseGlooResourceApi
export const useListFederatedAuthConfigs = () =>
  useRequest(federatedEnterpriseGlooResourceApi.listFederatedAuthConfigs, []);
export const useListFederatedRateLimits = () =>
  useRequest(
    federatedEnterpriseGlooResourceApi.listFederatedRateLimitConfigs,
    []
  );
//
// Federated End
//
// ------------------------ //

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
  return { glooInstance, glooInstances, instancesError };
}

export function usePageApiRef() {
  const {
    graphqlApiName = '',
    graphqlApiNamespace = '',
    graphqlApiClusterName = '',
  } = useParams();
  return {
    name: graphqlApiName,
    namespace: graphqlApiNamespace,
    clusterName: graphqlApiClusterName,
  } as ClusterObjectRef.AsObject;
}

export function useGetConsoleOptions() {
  const [readonly, setReadonly] = useState(true);
  const [apiExplorerEnabled, setApiExplorerEnabled] = useState(false);
  const [errorMessage, setErrorMessage] = useState('');
  const { data, error } = useRequest(bootstrapApi.getConsoleOptions, []);
  useEffect(() => {
    if (data?.options) {
      setReadonly(data.options.readOnly);
      setApiExplorerEnabled(data.options.apiExplorerEnabled);
    }
    if (error) {
      setErrorMessage(error?.message ?? error);
    }
  }, [data]);
  return { readonly, apiExplorerEnabled, errorMessage };
}
