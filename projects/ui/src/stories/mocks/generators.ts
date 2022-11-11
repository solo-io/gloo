import { faker } from '@faker-js/faker';
import { Struct, Value } from 'google-protobuf/google/protobuf/struct_pb';
import {
  BoolValue,
  UInt32Value,
} from 'google-protobuf/google/protobuf/wrappers_pb';
import { ClusterObjectRef } from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';
import {
  GatewaySpec,
  TcpGateway,
} from 'proto/github.com/solo-io/solo-apis/api/gloo/gateway/v1/gateway_pb';
import {
  HttpGateway,
  VirtualServiceSelectorExpressions,
} from 'proto/github.com/solo-io/solo-apis/api/gloo/gateway/v1/http_gateway_pb';
import { Extensions } from 'proto/github.com/solo-io/solo-apis/api/gloo/gloo/v1/extensions_pb';
import {
  AccessLog,
  AccessLoggingService,
  FileSink,
  GrpcService,
} from 'proto/github.com/solo-io/solo-apis/api/gloo/gloo/v1/options/als/als_pb';
import { ProxyProtocol } from 'proto/github.com/solo-io/solo-apis/api/gloo/gloo/v1/options/proxy_protocol/proxy_protocol_pb';
import {
  HttpListenerOptions,
  ListenerOptions,
} from 'proto/github.com/solo-io/solo-apis/api/gloo/gloo/v1/options_pb';
import {
  Listener,
  ProxySpec,
  TcpHost,
} from 'proto/github.com/solo-io/solo-apis/api/gloo/gloo/v1/proxy_pb';
import { SslConfig } from 'proto/github.com/solo-io/solo-apis/api/gloo/gloo/v1/ssl_pb';
import {
  UpstreamSpec,
  UpstreamStatus,
} from 'proto/github.com/solo-io/solo-apis/api/gloo/gloo/v1/upstream_pb';
import {
  ExecutableSchema,
  GraphQLApiSpec,
  GraphQLApiStatus,
  StitchedSchema,
} from 'proto/github.com/solo-io/solo-apis/api/gloo/graphql.gloo/v1beta1/graphql_pb';
import { SocketOption } from 'proto/github.com/solo-io/solo-kit/api/external/envoy/api/v2/core/socket_option_pb';
import { ResourceRef } from 'proto/github.com/solo-io/solo-kit/api/v1/ref_pb';
import { FederatedAuthConfig } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/federated_enterprise_gloo_resources_pb';
import {
  FederatedGateway,
  FederatedRouteTable,
  FederatedVirtualService,
} from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/federated_gateway_resources_pb';
import {
  FederatedSettings,
  FederatedUpstream,
  FederatedUpstreamGroup,
} from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/federated_gloo_resources_pb';
import { FederatedRateLimitConfig } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/federated_ratelimit_resources_pb';
import {
  ObjectMeta,
  Time,
} from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/common_pb';
import { Gateway } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/gateway_resources_pb';
import { Placement } from 'proto/github.com/solo-io/solo-projects/projects/gloo-fed/api/multicluster/v1alpha1/multicluster_pb';

import {
  ClusterDetails,
  ConfigDump,
  GlooInstance,
} from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/glooinstance_pb';
import {
  Proxy,
  Upstream,
} from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/gloo_resources_pb';
import { GraphqlApi } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/graphql_pb';
import { FederatedAuthConfigSpec } from 'proto/github.com/solo-io/solo-projects/projects/gloo-fed/api/fed.enterprise.gloo/v1/auth_config_pb';
import { FederatedGatewaySpec } from 'proto/github.com/solo-io/solo-projects/projects/gloo-fed/api/fed.gateway/v1/gateway_pb';
import { FederatedRouteTableSpec } from 'proto/github.com/solo-io/solo-projects/projects/gloo-fed/api/fed.gateway/v1/route_table_pb';
import {
  FederatedVirtualServiceSpec,
  FederatedVirtualServiceStatus,
} from 'proto/github.com/solo-io/solo-projects/projects/gloo-fed/api/fed.gateway/v1/virtual_service_pb';
import { FederatedSettingsSpec } from 'proto/github.com/solo-io/solo-projects/projects/gloo-fed/api/fed.gloo/v1/settings_pb';
import { FederatedUpstreamGroupSpec } from 'proto/github.com/solo-io/solo-projects/projects/gloo-fed/api/fed.gloo/v1/upstream_group_pb';
import { FederatedUpstreamSpec } from 'proto/github.com/solo-io/solo-projects/projects/gloo-fed/api/fed.gloo/v1/upstream_pb';
import { FederatedRateLimitConfigSpec } from 'proto/github.com/solo-io/solo-projects/projects/gloo-fed/api/fed.ratelimit/v1alpha1/rate_limit_config_pb';
import {
  PlacementStatus,
  TemplateMetadata,
} from 'proto/github.com/solo-io/solo-projects/projects/gloo-fed/api/fed/core/v1/placement_pb';

const placementStatusStates = [
  PlacementStatus.State.UNKNOWN,
  PlacementStatus.State.PLACED,
  PlacementStatus.State.FAILED,
  PlacementStatus.State.STALE,
  PlacementStatus.State.INVALID,
  PlacementStatus.State.PENDING,
];

// functions for creating a fake status for a federated resource ===>
export const createFederatedResourceStatus = (
  meta: Partial<FederatedVirtualServiceStatus.AsObject> = {}
): FederatedVirtualServiceStatus.AsObject => {
  return {
    placementStatus: meta.placementStatus ?? createPlacementStatus(),
    namespacedPlacementStatusesMap:
      meta.namespacedPlacementStatusesMap ??
      Array.from({ length: Number(faker.random.numeric()) }).map(() => {
        return [faker.random.word(), createPlacementStatus()];
      }),
  };
};

export const createPlacementStatus = (
  meta: Partial<PlacementStatus.AsObject> = {}
): PlacementStatus.AsObject => {
  return {
    clustersMap:
      meta.clustersMap ??
      Array.from({ length: Number(faker.random.numeric()) }).map(() => {
        return [faker.random.word(), createCluster()];
      }),
    state: meta.state ?? faker.helpers.arrayElement(placementStatusStates),
    message: meta.message ?? faker.random.word(),
    observedGeneration:
      meta.observedGeneration ?? Number(faker.random.numeric()),
    //processingTime: meta.processingTime ?? createTimestamp(), // TODO add back when we add processingTime
    writtenBy: meta.writtenBy ?? faker.random.word(),
  };
};

export const createCluster = (
  meta: Partial<PlacementStatus.Cluster.AsObject> = {}
): PlacementStatus.Cluster.AsObject => {
  return {
    namespacesMap:
      meta.namespacesMap ??
      Array.from({ length: Number(faker.random.numeric()) }).map(() => {
        return [faker.random.word(), createNamespace()];
      }),
  };
};

export const createNamespace = (
  meta: Partial<PlacementStatus.Namespace.AsObject> = {}
): PlacementStatus.Namespace.AsObject => {
  return {
    state: meta.state ?? faker.helpers.arrayElement(placementStatusStates),
    message: meta.message ?? faker.random.word(),
  };
};
// <=== functions for creating a fake status for a federated resource

export const createInstanceStatus = (
  meta: Partial<GlooInstance.GlooInstanceStatus.AsObject> = {}
): GlooInstance.GlooInstanceStatus.AsObject => {
  return faker.helpers.arrayElement(placementStatusStates);
};

export const createTimestamp = (
  meta: Partial<Time.AsObject> = {}
): Time.AsObject => {
  return {
    seconds: meta.seconds ?? Number(faker.random.numeric(2)),
    nanos: meta.nanos ?? Number(faker.random.numeric(3)),
    ...meta,
  };
};

export const createObjMeta = (
  meta: Partial<ObjectMeta.AsObject> = {}
): ObjectMeta.AsObject => {
  return {
    name: meta.name ?? faker.random.word(),
    namespace: meta.namespace ?? faker.random.word(),
    uid: meta.uid ?? faker.datatype.uuid(),
    resourceVersion: meta.resourceVersion ?? faker.random.numeric(2),
    creationTimestamp: meta.creationTimestamp ?? createTimestamp(),
    labelsMap:
      meta.labelsMap ??
      Array.from({ length: Number(faker.random.numeric(1)) }).map(() => {
        return [faker.random.word(), faker.random.word()];
      }),
    annotationsMap:
      meta.annotationsMap ??
      Array.from({ length: Number(faker.random.numeric(1)) }).map(() => {
        return [faker.random.word(), faker.random.word()];
      }),
    clusterName: meta.clusterName ?? faker.random.word(),
    ...meta,
  };
};

export const createControlPlane = (
  meta: Partial<GlooInstance.GlooInstanceSpec.ControlPlane.AsObject> = {}
): GlooInstance.GlooInstanceSpec.ControlPlane.AsObject => {
  return {
    version: meta.version ?? faker.random.numeric(1),
    namespace: meta.namespace ?? faker.random.word(),
    watchedNamespacesList:
      meta.watchedNamespacesList ??
      Array.from({ length: Number(faker.random.numeric(2)) }).map(() => {
        return faker.random.word();
      }),
    ...meta,
  };
};

const createWorkloadControllerMap =
  (): GlooInstance.GlooInstanceSpec.Proxy.WorkloadControllerMap => {
    return {
      UNDEFINED: 0,
      DEPLOYMENT: 1,
      DAEMON_SET: 2,
    };
  };

const createPortsList = (
  meta: Partial<GlooInstance.GlooInstanceSpec.Proxy.IngressEndpoint.Port.AsObject> = {}
): GlooInstance.GlooInstanceSpec.Proxy.IngressEndpoint.Port.AsObject => {
  return {
    port: meta.port ?? Number(faker.random.numeric(2)),
    name: meta.name ?? faker.random.word(),
    ...meta,
  };
};

const createIngressEndpointsList = (
  meta: Partial<GlooInstance.GlooInstanceSpec.Proxy.IngressEndpoint.AsObject> = {}
): GlooInstance.GlooInstanceSpec.Proxy.IngressEndpoint.AsObject => {
  return {
    address: meta.address ?? faker.address.streetAddress(),
    serviceName: meta.serviceName ?? faker.random.word(),
    portsList:
      meta.portsList ??
      Array.from({ length: Number(faker.random.numeric(1)) }).map(() => {
        return createPortsList();
      }),
    ...meta,
  };
};

export const createGlooProxyObject = (
  meta: Partial<GlooInstance.GlooInstanceSpec.Proxy.AsObject> = {}
): GlooInstance.GlooInstanceSpec.Proxy.AsObject => {
  return {
    replicas: meta.replicas ?? Number(faker.random.numeric(1)),
    availableReplicas:
      meta.availableReplicas ?? Number(faker.random.numeric(1)),
    readyReplicas: meta.readyReplicas ?? Number(faker.random.numeric(1)),
    wasmEnabled: meta.wasmEnabled ?? faker.datatype.boolean(),
    readConfigMulticlusterEnabled:
      meta.readConfigMulticlusterEnabled ?? faker.datatype.boolean(),
    version: meta.version ?? faker.random.word(),
    name: meta.name ?? faker.random.word(),
    namespace: meta.namespace ?? faker.random.word(),
    workloadControllerType:
      meta.workloadControllerType ??
      faker.helpers.arrayElement([
        createWorkloadControllerMap().DAEMON_SET,
        createWorkloadControllerMap().DEPLOYMENT,
        createWorkloadControllerMap().UNDEFINED,
      ]),
    zonesList:
      meta.zonesList ??
      Array.from({ length: Number(faker.random.numeric(1)) }).map(() => {
        return faker.random.word();
      }),
    ingressEndpointsList:
      meta.ingressEndpointsList ??
      Array.from({
        length: faker.datatype.number(2),
      }).map(() => {
        return createIngressEndpointsList();
      }),
    ...meta,
  } as GlooInstance.GlooInstanceSpec.Proxy.AsObject;
};

interface ObjectRef {
  AsObject: {
    name: string;
    namespace: string;
  };
}

export const createObjRef = (
  meta: Partial<ObjectRef['AsObject']> = {}
): ObjectRef['AsObject'] => {
  return {
    name: meta.name ?? faker.random.word(),
    namespace: meta.namespace ?? faker.random.word(),
    ...meta,
  };
};

export const createResourceReport = (
  meta: Partial<GlooInstance.GlooInstanceSpec.Check.Summary.ResourceReport.AsObject> = {}
): GlooInstance.GlooInstanceSpec.Check.Summary.ResourceReport.AsObject => {
  return {
    ref: meta.ref ?? createObjRef(),
    message: meta.message ?? faker.random.word(),
    ...meta,
  };
};

export const createCheckSummary = (
  meta: Partial<GlooInstance.GlooInstanceSpec.Check.Summary.AsObject> = {}
): GlooInstance.GlooInstanceSpec.Check.Summary.AsObject => {
  return {
    total: meta.total ?? faker.datatype.number(5),
    errorsList:
      meta.errorsList ??
      Array.from({ length: faker.datatype.number(3) }).map(() => {
        return createResourceReport();
      }),
    warningsList:
      meta.warningsList ??
      Array.from({ length: faker.datatype.number(3) }).map(() => {
        return createResourceReport();
      }),
  };
};

export const createGlooCheck = (
  meta: Partial<GlooInstance.GlooInstanceSpec.Check.AsObject> = {}
): GlooInstance.GlooInstanceSpec.Check.AsObject => {
  return {
    gateways: meta.gateways ?? createCheckSummary(),
    virtualServices: meta.virtualServices ?? createCheckSummary(),
    routeTables: meta.routeTables ?? createCheckSummary(),
    authConfigs: meta.authConfigs ?? createCheckSummary(),
    settings: meta.settings ?? createCheckSummary(),
    upstreams: meta.upstreams ?? createCheckSummary(),
    upstreamGroups: meta.upstreamGroups ?? createCheckSummary(),
    proxies: meta.proxies ?? createCheckSummary(),
    rateLimitConfigs: meta.rateLimitConfigs ?? createCheckSummary(),
    matchableHttpGateways: meta.matchableHttpGateways ?? createCheckSummary(),
    deployments: meta.deployments ?? createCheckSummary(),
    pods: meta.pods ?? createCheckSummary(),
    ...meta,
  };
};

export const createObjSpec = (
  meta: Partial<GlooInstance.GlooInstanceSpec.AsObject> = {}
): GlooInstance.GlooInstanceSpec.AsObject => {
  return {
    cluster: meta.cluster ?? faker.random.word(),
    isEnterprise: meta.isEnterprise ?? faker.datatype.boolean(),
    controlPlane: meta.controlPlane ?? createControlPlane(),
    proxiesList:
      meta.proxiesList ??
      Array.from({ length: Number(faker.random.numeric(1)) }).map(() =>
        createGlooProxyObject()
      ),
    region: meta.region ?? faker.random.word(),
    check: meta.check ?? createGlooCheck(),
    ...meta,
  };
};

export const createGlooInstanceObj = (
  meta: Partial<GlooInstance.AsObject> = {}
): GlooInstance.AsObject => {
  return {
    metadata: createObjMeta(),
    spec: createObjSpec(),
    status: createInstanceStatus(),
    ...meta,
  };
};

export const createGlooInstance = () => {
  const instance = new GlooInstance();
  instance.setMetadata();
  instance.setStatus();
  return instance;
};

export const createClusterDetails = (): ClusterDetails => {
  const details = new ClusterDetails();
  details.setCluster(faker.random.word());
  const glooInstances = Array.from({ length: faker.datatype.number(2) }).map(
    () => {
      return createGlooInstance();
    }
  );
  details.setGlooInstancesList(glooInstances);
  return details;
};

export const createClusterDetailsObj = (
  meta: Partial<
    ClusterDetails.AsObject & { glooInstance: Partial<GlooInstance.AsObject> }
  > = {}
): ClusterDetails.AsObject => {
  return {
    cluster: meta.cluster ?? faker.random.word(),
    glooInstancesList: Array.from({
      length: faker.datatype.number(4),
    }).map(() => {
      return createGlooInstanceObj(meta.glooInstance);
    }),
    ...meta,
  };
};
// spec?: github_com_solo_io_solo_apis_api_gloo_gateway_v1_virtual_service_pb.VirtualServiceSpec.AsObject,
//       metadata?: github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_core_v1_placement_pb.TemplateMetadata.AsObject,
export const createFederatedVirtualServiceSpecTemplate =
  (): FederatedVirtualServiceSpec.Template.AsObject => {
    return {};
  };
// TemplateMetadata.AsObject
export const createTemplateMetadata = (
  meta: Partial<TemplateMetadata.AsObject> = {}
): TemplateMetadata.AsObject => {
  return {
    annotationsMap:
      meta.annotationsMap ??
      Array.from({ length: faker.datatype.number(2) }).map(() => {
        return [faker.random.word(), faker.random.word()];
      }),
    labelsMap:
      meta.labelsMap ??
      Array.from({ length: faker.datatype.number(2) }).map(() => {
        return [faker.random.word(), faker.random.word()];
      }),
    name: meta.name ?? faker.random.word(),
  };
};

export const createPlacement = (
  meta: Partial<Placement.AsObject> = {}
): Placement.AsObject => {
  return {
    namespacesList:
      meta.namespacesList ??
      Array.from({ length: faker.datatype.number(4) }).map(() => {
        return faker.random.word();
      }),
    clustersList:
      meta.clustersList ??
      Array.from({ length: faker.datatype.number(4) }).map(() => {
        return faker.random.word();
      }),
  };
};

export const createFederatedVirtualServiceSpec = (
  meta: Partial<FederatedVirtualServiceSpec.AsObject> = {}
): FederatedVirtualServiceSpec.AsObject => {
  return {
    template: meta.template ?? undefined,
    placement: meta.placement ?? createPlacement(),
  };
};

export const createFederatedVirtualService = (
  meta: Partial<FederatedVirtualService.AsObject> = {}
): FederatedVirtualService.AsObject => {
  return {
    metadata: meta.metadata ?? createObjMeta(),
    spec: meta.spec ?? createFederatedVirtualServiceSpec(),
    status: meta.status ?? createFederatedResourceStatus(),
  };
};

export const createFederatedRouteTableSpec = (
  meta: Partial<FederatedRouteTableSpec.AsObject> = {}
): FederatedRouteTableSpec.AsObject => {
  return {
    template: meta.template ?? undefined,
    placement: meta.placement ?? createPlacement(),
  };
};

export const createListFederatedRouteTables = (
  meta: Partial<FederatedRouteTable.AsObject> = {}
): FederatedRouteTable.AsObject => {
  return {
    metadata: meta.metadata ?? createObjMeta(),
    spec: meta.spec ?? createFederatedRouteTableSpec(),
    status: meta.status ?? createFederatedResourceStatus(),
  };
};

export const createFederatedUpstreamSpec = (
  meta: Partial<FederatedUpstreamSpec.AsObject> = {}
): FederatedUpstreamSpec.AsObject => {
  return {
    template: meta.template,
    placement: meta.placement ?? createPlacement(),
  };
};

export const createFederatedUpstream = (
  meta: Partial<FederatedUpstream.AsObject> = {}
): FederatedUpstream.AsObject => {
  return {
    metadata: meta.metadata ?? createObjMeta(),
    spec: meta.spec ?? createFederatedUpstreamSpec(),
    status: meta.status ?? createFederatedResourceStatus(),
  };
};

export const createFederatedUpstreamGroupSpec = (
  meta: Partial<FederatedUpstreamGroupSpec.AsObject> = {}
): FederatedUpstreamGroupSpec.AsObject => {
  return {
    template: meta.template,
    placement: meta.placement ?? createPlacement(),
  };
};

export const createFederatedUpstreamGroup = (
  meta: Partial<FederatedUpstreamGroup.AsObject> = {}
): FederatedUpstreamGroup.AsObject => {
  return {
    metadata: meta.metadata ?? createObjMeta(),
    spec: meta.spec ?? createFederatedUpstreamGroupSpec(),
    status: meta.status ?? createFederatedResourceStatus(),
  };
};

export const createFederatedAuthConfigSpec = (
  meta: Partial<FederatedAuthConfigSpec.AsObject> = {}
): FederatedAuthConfigSpec.AsObject => {
  return {
    template: meta.template,
    placement: meta.placement ?? createPlacement(),
  };
};

export const createFederatedAuthConfig = (
  meta: Partial<FederatedAuthConfig.AsObject> = {}
): FederatedAuthConfig.AsObject => {
  return {
    metadata: meta.metadata ?? createObjMeta(),
    spec: meta.spec ?? createFederatedAuthConfigSpec(),
    status: meta.status ?? createFederatedResourceStatus(),
  };
};

export const createFederatedGatewaySpec = (
  meta: Partial<FederatedGatewaySpec.AsObject> = {}
): FederatedGatewaySpec.AsObject => {
  return {
    template: meta.template,
    placement: meta.placement ?? createPlacement(),
  };
};

export const createFederatedGateway = (
  meta: Partial<FederatedGateway.AsObject> = {}
): FederatedGateway.AsObject => {
  return {
    metadata: meta.metadata ?? createObjMeta(),
    spec: meta.spec ?? createFederatedGatewaySpec(),
    status: meta.status ?? createFederatedResourceStatus(),
  };
};

export const createFederatedSettingsSpec = (
  meta: Partial<FederatedSettingsSpec.AsObject> = {}
): FederatedSettingsSpec.AsObject => {
  return {
    template: meta.template,
    placement: meta.placement ?? createPlacement(),
  };
};

export const createFederatedSettings = (
  meta: Partial<FederatedSettings.AsObject> = {}
): FederatedSettings.AsObject => {
  return {
    metadata: meta.metadata ?? createObjMeta(),
    spec: meta.spec ?? createFederatedSettingsSpec(),
    status: meta.status ?? createFederatedResourceStatus(),
  };
};

export const createStruct = (
  meta: Partial<Struct.AsObject> = {}
): Struct.AsObject => {
  return {
    fieldsMap:
      meta.fieldsMap ??
      Array.from({
        length: faker.datatype.number(5),
      }).map(() => {
        return [faker.random.word(), new Value().toObject()];
      }),
  };
};

export const createFileSink = (
  meta: Partial<FileSink.AsObject> = {}
): FileSink.AsObject => {
  return {
    path: meta.path ?? faker.random.word(),
    stringFormat: meta.stringFormat ?? faker.random.word(),
    jsonFormat: meta.jsonFormat ?? createStruct(),
  };
};

export const createGrpcService = (
  meta: Partial<GrpcService.AsObject> = {}
): GrpcService.AsObject => {
  return {
    logName: meta.logName ?? faker.random.word(),
    staticClusterName: meta.staticClusterName ?? faker.random.word(),
    additionalRequestHeadersToLogList:
      meta.additionalRequestHeadersToLogList ??
      Array.from({
        length: faker.datatype.number(2),
      }).map(() => {
        return faker.random.word();
      }),
    additionalResponseHeadersToLogList:
      meta.additionalResponseHeadersToLogList ??
      Array.from({
        length: faker.datatype.number(2),
      }).map(() => {
        return faker.random.word();
      }),
    additionalResponseTrailersToLogList:
      meta.additionalResponseTrailersToLogList ??
      Array.from({
        length: faker.datatype.number(2),
      }).map(() => {
        return faker.random.word();
      }),
  };
};

export const createAccessLog = (
  meta: Partial<AccessLog.AsObject> = {}
): AccessLog.AsObject => {
  return {
    fileSink: meta.fileSink ?? createFileSink(),
    grpcService: meta.grpcService ?? createGrpcService(),
  };
};

export const createAccessLoggingService = (
  meta: Partial<AccessLoggingService.AsObject> = {}
): AccessLoggingService.AsObject => {
  return {
    accessLogList:
      meta.accessLogList ??
      Array.from({
        length: faker.datatype.number(4),
      }).map(() => {
        return createAccessLog();
      }),
  };
};

export const createExtensions = (
  meta: Partial<Extensions.AsObject> = {}
): Extensions.AsObject => {
  return {
    configsMap:
      meta.configsMap ??
      Array.from({
        length: faker.datatype.number(5),
      }).map(() => {
        return [faker.random.word(), createStruct()];
      }),
  };
};

export const createSocketOption = (
  meta: Partial<SocketOption.AsObject> = {}
): SocketOption.AsObject => {
  return {
    description: meta.description ?? faker.random.word(),
    level: meta.level ?? faker.datatype.number(3),
    name: meta.name ?? faker.datatype.number(5),
    intValue: meta.intValue ?? faker.datatype.number(3),
    bufValue: meta.bufValue ?? faker.random.word(),
    state: meta.state ?? faker.helpers.arrayElement([0, 1, 2]),
  };
};

export const createProxyProtocolKeyValuePair = (
  meta: Partial<ProxyProtocol.KeyValuePair.AsObject> = {}
): ProxyProtocol.KeyValuePair.AsObject => {
  return {
    metadataNamespace: meta.metadataNamespace ?? faker.random.word(),
    key: meta.key ?? faker.random.word(),
  };
};

export const createProxyProtocolRule = (
  meta: Partial<ProxyProtocol.Rule.AsObject> = {}
): ProxyProtocol.Rule.AsObject => {
  return {
    tlvType: meta.tlvType ?? faker.datatype.number(10),
    onTlvPresent: meta.onTlvPresent ?? createProxyProtocolKeyValuePair(),
  };
};

export const createProxyProtocol = (
  meta: Partial<ProxyProtocol.AsObject> = {}
): ProxyProtocol.AsObject => {
  return {
    rulesList:
      meta.rulesList ??
      Array.from({
        length: faker.datatype.number(4),
      }).map(() => {
        return createProxyProtocolRule();
      }),
    allowRequestsWithoutProxyProtocol:
      meta.allowRequestsWithoutProxyProtocol ?? faker.datatype.boolean(),
  };
};

export const createListenerOptions = (
  meta: Partial<ListenerOptions.AsObject> = {}
): ListenerOptions.AsObject => {
  return {
    accessLoggingService:
      meta.accessLoggingService ?? createAccessLoggingService(),
    extensions: meta.extensions ?? createExtensions(),
    perConnectionBufferLimitBytes:
      meta.perConnectionBufferLimitBytes ??
      new UInt32Value().setValue(faker.datatype.number(10)).toObject(),
    socketOptionsList:
      meta.socketOptionsList ??
      Array.from({ length: faker.datatype.number(7) }).map(() => {
        return createSocketOption();
      }),
    proxyProtocol: meta.proxyProtocol ?? createProxyProtocol(),
  };
};

export const createBoolValue = (meta?: boolean) => {
  const boolVal = new BoolValue();
  boolVal.setValue(meta ?? faker.datatype.boolean());
  return boolVal.toObject();
};

export const createResourceRef = (meta: Partial<ResourceRef.AsObject> = {}) => {
  return {
    name: meta.name ?? faker.random.word(),
    namespace: meta.namespace ?? faker.random.word(),
  };
};

export const createVirtualServiceSelectorExpressionsExpression = (
  meta: Partial<VirtualServiceSelectorExpressions.Expression.AsObject> = {}
): VirtualServiceSelectorExpressions.Expression.AsObject => {
  return {
    key: meta.key ?? faker.random.word(),
    operator:
      meta.operator ?? faker.helpers.arrayElement([0, 1, 2, 3, 4, 5, 6, 7, 8]),
    valuesList:
      meta.valuesList ??
      Array.from({ length: faker.datatype.number(1) }).map(() => {
        return faker.random.word();
      }),
  };
};

export const createVirtualServiceSelectorExpressions = (
  meta: Partial<VirtualServiceSelectorExpressions.AsObject> = {}
): VirtualServiceSelectorExpressions.AsObject => {
  return {
    expressionsList:
      meta.expressionsList ??
      Array.from({ length: faker.datatype.number() }).map(() => {
        return createVirtualServiceSelectorExpressionsExpression();
      }),
  };
};

export const createHttpListenerOptions = (
  meta: Partial<HttpListenerOptions.AsObject> = {}
): HttpListenerOptions.AsObject => {
  return {};
};

export const createHttpGateway = (meta: Partial<HttpGateway.AsObject> = {}) => {
  return {
    virtualServicesList:
      meta.virtualServicesList ??
      Array.from({ length: faker.datatype.number(3) }).map(() => {
        return createResourceRef();
      }),
    virtualServiceSelectorMap:
      meta.virtualServiceSelectorMap ??
      Array.from({ length: faker.datatype.number(3) }).map(() => {
        return [faker.random.word(), faker.random.word()];
      }),
    virtualServiceExpressions:
      meta.virtualServiceExpressions ??
      createVirtualServiceSelectorExpressions(),
    virtualServiceNamespacesList:
      meta.virtualServiceNamespacesList ??
      Array.from({
        length: faker.datatype.number(5),
      }).map(() => {
        return faker.random.word();
      }),
    options: meta.options ?? createHttpListenerOptions(),
  };
};

export const createTcpHost = (
  meta: Partial<TcpHost.AsObject> = {}
): TcpHost.AsObject => {
  return {
    name: meta.name ?? faker.random.word(),
    sslConfig: meta.sslConfig ?? undefined,
    destination: meta.destination ?? undefined,
  };
};

export const createTcpGateway = (
  meta: Partial<TcpGateway.AsObject> = {}
): TcpGateway.AsObject => {
  return {
    tcpHostsList:
      meta.tcpHostsList ??
      Array.from({
        length: faker.datatype.number(4),
      }).map(() => {
        return createTcpHost();
      }),
    options: meta.options ?? undefined,
  };
};

export const createGatewaySpec = (
  meta: Partial<GatewaySpec.AsObject> = {}
): GatewaySpec.AsObject => {
  return {
    ssl: meta.ssl != null ? meta.ssl : faker.datatype.boolean(),
    bindAddress: meta.bindAddress ?? faker.internet.ip(),
    bindPort: meta.bindPort ?? faker.datatype.number(1000),
    options: meta.options ?? createListenerOptions(),
    useProxyProto: meta.useProxyProto ?? createBoolValue(),
    httpGateway: meta.httpGateway ?? createHttpGateway(),
    tcpGateway: meta.tcpGateway ?? createTcpGateway(),
    hybridGateway: meta.httpGateway as any,
    proxyNamesList:
      meta.proxyNamesList ??
      Array.from({ length: faker.datatype.number(5) }).map(() => {
        return faker.random.word();
      }),
    routeOptions: meta.routeOptions ?? undefined,
  };
};

export const createGateway = (
  meta: Partial<Gateway.AsObject> = {}
): Gateway.AsObject => {
  return {
    metadata: meta.metadata ?? createObjMeta(),
    spec: meta.spec ?? createGatewaySpec(),
    status: meta.status ?? undefined,
    glooInstance: meta.glooInstance ?? undefined,
  };
};

export const createClusterObjectRef = (
  meta: Partial<ClusterObjectRef.AsObject> = {}
): ClusterObjectRef.AsObject => {
  return {
    name: meta.name ?? faker.random.word(),
    namespace: meta.namespace ?? faker.random.word(),
    clusterName: meta.clusterName ?? faker.random.word(),
  };
};

export const createUpstreamStatus = (
  meta: Partial<UpstreamStatus.AsObject> = {}
): UpstreamStatus.AsObject => {
  return {
    reason: 'Some reason',
    reportedBy: 'Fake status generator',
    state: UpstreamStatus.State.ACCEPTED,
    subresourceStatusesMap: [],
  };
};

export const createUpstreamSpec = (
  meta: Partial<UpstreamSpec.AsObject> = {}
): UpstreamSpec.AsObject => {
  return {
    healthChecksList: [],
    httpConnectHeadersList: [],
    protocolSelection:
      UpstreamSpec.ClusterProtocolSelection.USE_CONFIGURED_PROTOCOL,
    ...meta,
  };
};

export const createUpstream = (
  meta: Partial<Upstream.AsObject> = {}
): Upstream.AsObject => {
  return {
    metadata: meta.metadata ?? createObjMeta(),
    spec: meta.spec ?? createUpstreamSpec(),
    glooInstance: meta.glooInstance ?? undefined,
    status: meta.status ?? createUpstreamStatus(),
  };
};

export const createFederatedRateLimitConfigSpec = (
  meta: Partial<FederatedRateLimitConfigSpec.AsObject> = {}
): FederatedRateLimitConfigSpec.AsObject => {
  return {
    template: meta.template,
    placement: createPlacement(),
  };
};

export const createFederatedRateLimitConfig = (
  meta: Partial<FederatedRateLimitConfig.AsObject> = {}
): FederatedRateLimitConfig.AsObject => {
  return {
    metadata: meta.metadata ?? createObjMeta(),
    spec: meta.spec ?? createFederatedRateLimitConfigSpec(),
    status: meta.status ?? undefined,
  };
};

export const createConfigDump = (
  meta: Partial<ConfigDump.AsObject> = {}
): ConfigDump.AsObject => {
  return {
    name: meta.name ?? faker.random.word(),
    raw: meta.raw ?? faker.random.word(),
    error: meta.error ?? faker.random.word(),
  };
};

export const createSslConfigurations = (
  meta: Partial<SslConfig.AsObject> = {}
): SslConfig.AsObject => {
  return {
    secretRef: meta.secretRef,
    sslFiles: meta.sslFiles,
    sds: meta.sds,
    sniDomainsList:
      meta.sniDomainsList ??
      Array.from({ length: 1 }).map(() => {
        return faker.random.word();
      }),
    verifySubjectAltNameList:
      meta.verifySubjectAltNameList ??
      Array.from({ length: 1 }).map(() => {
        return faker.random.word();
      }),
    parameters: meta.parameters,
    alpnProtocolsList:
      meta.alpnProtocolsList ??
      Array.from({ length: 1 }).map(() => {
        return faker.random.word();
      }),
    oneWayTls: meta.oneWayTls,
    disableTlsSessionResumption: meta.disableTlsSessionResumption,
    transportSocketConnectTimeout: meta.transportSocketConnectTimeout,
  };
};

export const createListener = (
  meta: Partial<Listener.AsObject> = {}
): Listener.AsObject => {
  return {
    ...meta,
    name: meta.name ?? faker.random.word(),
    bindAddress: meta.bindAddress ?? faker.internet.ipv4(),
    bindPort: meta.bindPort ?? faker.datatype.number(1000),
    httpListener: meta.httpListener,
    tcpListener: meta.tcpListener,
    sslConfigurationsList:
      meta.sslConfigurationsList ??
      Array.from({ length: 1 }).map(() => {
        return createSslConfigurations();
      }),
  };
};

export const createProxySpec = (
  meta: Partial<ProxySpec.AsObject> = {}
): ProxySpec.AsObject => {
  return {
    compressedspec: meta.compressedspec ?? faker.random.word(),
    listenersList: Array.from({ length: 1 }).map(() => {
      return createListener();
    }),
  };
};

export const createProxy = (
  meta: Partial<Proxy.AsObject> = {}
): Proxy.AsObject => {
  return {
    metadata: meta.metadata ?? createObjMeta(),
    spec: meta.spec ?? createProxySpec(),
  };
};

export const createExecutableSchema = (
  meta: Partial<ExecutableSchema.AsObject> = {}
): ExecutableSchema.AsObject => {
  return {
    ...meta,
    schemaDefinition:
      meta.schemaDefinition ?? faker.random.words(faker.datatype.number(12)),
  };
};

export const createStitchedSchemaSubschemaConfigTypeMergeConfig = (
  meta: Partial<StitchedSchema.SubschemaConfig.TypeMergeConfig.AsObject> = {}
): StitchedSchema.SubschemaConfig.TypeMergeConfig.AsObject => {
  return {
    selectionSet: meta.selectionSet ?? faker.random.word(),
    queryName: meta.queryName ?? faker.random.word(),
    argsMap:
      meta.argsMap ??
      Array.from({ length: 2 }).map(() => {
        return [faker.random.word(), faker.random.word()];
      }),
  };
};

export const createStitchedSchemaSubschemaConfig = (
  meta: Partial<StitchedSchema.SubschemaConfig.AsObject> = {}
): StitchedSchema.SubschemaConfig.AsObject => {
  return {
    name: meta.name ?? faker.random.word(),
    namespace: meta.namespace ?? faker.random.word(),
    typeMergeMap:
      meta.typeMergeMap ??
      Array.from({ length: 2 }).map(() => {
        return [
          faker.random.word(),
          createStitchedSchemaSubschemaConfigTypeMergeConfig(),
        ];
      }),
  };
};

export const createStitchedSchema = (
  meta: Partial<StitchedSchema.AsObject> = {}
): StitchedSchema.AsObject => {
  return {
    subschemasList:
      meta.subschemasList ??
      Array.from({ length: 2 }).map(() => {
        return createStitchedSchemaSubschemaConfig();
      }),
  };
};

// TODO:  Add in more pieces.
export const createGraphqlApiSpec = (
  meta: Partial<GraphQLApiSpec.AsObject> = {}
): GraphQLApiSpec.AsObject => {
  return {
    executableSchema: meta.executableSchema ?? createExecutableSchema(),
    stitchedSchema: meta.stitchedSchema ?? createStitchedSchema(),
    allowedQueryHashesList:
      meta.allowedQueryHashesList ??
      Array.from({ length: 1 }).map(() => {
        return faker.random.word();
      }),
  };
};

export const createGraphQLApiStatus = (
  meta: Partial<GraphQLApiStatus.AsObject> = {}
): GraphQLApiStatus.AsObject => {
  return {
    state: meta.state ?? faker.helpers.arrayElement([0, 1, 2, 3]),
    reason: meta.reason ?? faker.random.word(),
    reportedBy: meta.reportedBy ?? faker.internet.userName(),
    subresourceStatusesMap: meta.subresourceStatusesMap ?? [],
    details: meta.details,
  };
};

export const createGraphqlApi = (
  meta: Partial<GraphqlApi.AsObject> = {}
): GraphqlApi.AsObject => {
  return {
    ...meta,
    metadata: meta.metadata ?? createObjMeta(),
    spec: meta.spec ?? createGraphqlApiSpec(),
    status: meta.status ?? createGraphQLApiStatus(),
    glooInstance: meta.glooInstance ?? createObjRef(),
  };
};
