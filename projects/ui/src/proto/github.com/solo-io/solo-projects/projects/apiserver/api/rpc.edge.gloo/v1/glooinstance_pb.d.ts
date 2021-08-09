/* eslint-disable */
// package: rpc.edge.gloo.solo.io
// file: github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/glooinstance.proto

import * as jspb from "google-protobuf";
import * as extproto_ext_pb from "../../../../../../../../extproto/ext_pb";
import * as github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb from "../../../../../../../../github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/common_pb";
import * as github_com_solo_io_skv2_api_core_v1_core_pb from "../../../../../../../../github.com/solo-io/skv2/api/core/v1/core_pb";

export class GlooInstance extends jspb.Message {
  hasMetadata(): boolean;
  clearMetadata(): void;
  getMetadata(): github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ObjectMeta | undefined;
  setMetadata(value?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ObjectMeta): void;

  hasSpec(): boolean;
  clearSpec(): void;
  getSpec(): GlooInstance.GlooInstanceSpec | undefined;
  setSpec(value?: GlooInstance.GlooInstanceSpec): void;

  hasStatus(): boolean;
  clearStatus(): void;
  getStatus(): GlooInstance.GlooInstanceStatus | undefined;
  setStatus(value?: GlooInstance.GlooInstanceStatus): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GlooInstance.AsObject;
  static toObject(includeInstance: boolean, msg: GlooInstance): GlooInstance.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GlooInstance, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GlooInstance;
  static deserializeBinaryFromReader(message: GlooInstance, reader: jspb.BinaryReader): GlooInstance;
}

export namespace GlooInstance {
  export type AsObject = {
    metadata?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ObjectMeta.AsObject,
    spec?: GlooInstance.GlooInstanceSpec.AsObject,
    status?: GlooInstance.GlooInstanceStatus.AsObject,
  }

  export class GlooInstanceSpec extends jspb.Message {
    getCluster(): string;
    setCluster(value: string): void;

    getIsEnterprise(): boolean;
    setIsEnterprise(value: boolean): void;

    hasControlPlane(): boolean;
    clearControlPlane(): void;
    getControlPlane(): GlooInstance.GlooInstanceSpec.ControlPlane | undefined;
    setControlPlane(value?: GlooInstance.GlooInstanceSpec.ControlPlane): void;

    clearProxiesList(): void;
    getProxiesList(): Array<GlooInstance.GlooInstanceSpec.Proxy>;
    setProxiesList(value: Array<GlooInstance.GlooInstanceSpec.Proxy>): void;
    addProxies(value?: GlooInstance.GlooInstanceSpec.Proxy, index?: number): GlooInstance.GlooInstanceSpec.Proxy;

    getRegion(): string;
    setRegion(value: string): void;

    hasCheck(): boolean;
    clearCheck(): void;
    getCheck(): GlooInstance.GlooInstanceSpec.Check | undefined;
    setCheck(value?: GlooInstance.GlooInstanceSpec.Check): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): GlooInstanceSpec.AsObject;
    static toObject(includeInstance: boolean, msg: GlooInstanceSpec): GlooInstanceSpec.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: GlooInstanceSpec, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): GlooInstanceSpec;
    static deserializeBinaryFromReader(message: GlooInstanceSpec, reader: jspb.BinaryReader): GlooInstanceSpec;
  }

  export namespace GlooInstanceSpec {
    export type AsObject = {
      cluster: string,
      isEnterprise: boolean,
      controlPlane?: GlooInstance.GlooInstanceSpec.ControlPlane.AsObject,
      proxiesList: Array<GlooInstance.GlooInstanceSpec.Proxy.AsObject>,
      region: string,
      check?: GlooInstance.GlooInstanceSpec.Check.AsObject,
    }

    export class ControlPlane extends jspb.Message {
      getVersion(): string;
      setVersion(value: string): void;

      getNamespace(): string;
      setNamespace(value: string): void;

      clearWatchedNamespacesList(): void;
      getWatchedNamespacesList(): Array<string>;
      setWatchedNamespacesList(value: Array<string>): void;
      addWatchedNamespaces(value: string, index?: number): string;

      serializeBinary(): Uint8Array;
      toObject(includeInstance?: boolean): ControlPlane.AsObject;
      static toObject(includeInstance: boolean, msg: ControlPlane): ControlPlane.AsObject;
      static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
      static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
      static serializeBinaryToWriter(message: ControlPlane, writer: jspb.BinaryWriter): void;
      static deserializeBinary(bytes: Uint8Array): ControlPlane;
      static deserializeBinaryFromReader(message: ControlPlane, reader: jspb.BinaryReader): ControlPlane;
    }

    export namespace ControlPlane {
      export type AsObject = {
        version: string,
        namespace: string,
        watchedNamespacesList: Array<string>,
      }
    }

    export class Proxy extends jspb.Message {
      getReplicas(): number;
      setReplicas(value: number): void;

      getAvailableReplicas(): number;
      setAvailableReplicas(value: number): void;

      getReadyReplicas(): number;
      setReadyReplicas(value: number): void;

      getWasmEnabled(): boolean;
      setWasmEnabled(value: boolean): void;

      getReadConfigMulticlusterEnabled(): boolean;
      setReadConfigMulticlusterEnabled(value: boolean): void;

      getVersion(): string;
      setVersion(value: string): void;

      getName(): string;
      setName(value: string): void;

      getNamespace(): string;
      setNamespace(value: string): void;

      getWorkloadControllerType(): GlooInstance.GlooInstanceSpec.Proxy.WorkloadControllerMap[keyof GlooInstance.GlooInstanceSpec.Proxy.WorkloadControllerMap];
      setWorkloadControllerType(value: GlooInstance.GlooInstanceSpec.Proxy.WorkloadControllerMap[keyof GlooInstance.GlooInstanceSpec.Proxy.WorkloadControllerMap]): void;

      clearZonesList(): void;
      getZonesList(): Array<string>;
      setZonesList(value: Array<string>): void;
      addZones(value: string, index?: number): string;

      clearIngressEndpointsList(): void;
      getIngressEndpointsList(): Array<GlooInstance.GlooInstanceSpec.Proxy.IngressEndpoint>;
      setIngressEndpointsList(value: Array<GlooInstance.GlooInstanceSpec.Proxy.IngressEndpoint>): void;
      addIngressEndpoints(value?: GlooInstance.GlooInstanceSpec.Proxy.IngressEndpoint, index?: number): GlooInstance.GlooInstanceSpec.Proxy.IngressEndpoint;

      serializeBinary(): Uint8Array;
      toObject(includeInstance?: boolean): Proxy.AsObject;
      static toObject(includeInstance: boolean, msg: Proxy): Proxy.AsObject;
      static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
      static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
      static serializeBinaryToWriter(message: Proxy, writer: jspb.BinaryWriter): void;
      static deserializeBinary(bytes: Uint8Array): Proxy;
      static deserializeBinaryFromReader(message: Proxy, reader: jspb.BinaryReader): Proxy;
    }

    export namespace Proxy {
      export type AsObject = {
        replicas: number,
        availableReplicas: number,
        readyReplicas: number,
        wasmEnabled: boolean,
        readConfigMulticlusterEnabled: boolean,
        version: string,
        name: string,
        namespace: string,
        workloadControllerType: GlooInstance.GlooInstanceSpec.Proxy.WorkloadControllerMap[keyof GlooInstance.GlooInstanceSpec.Proxy.WorkloadControllerMap],
        zonesList: Array<string>,
        ingressEndpointsList: Array<GlooInstance.GlooInstanceSpec.Proxy.IngressEndpoint.AsObject>,
      }

      export class IngressEndpoint extends jspb.Message {
        getAddress(): string;
        setAddress(value: string): void;

        clearPortsList(): void;
        getPortsList(): Array<GlooInstance.GlooInstanceSpec.Proxy.IngressEndpoint.Port>;
        setPortsList(value: Array<GlooInstance.GlooInstanceSpec.Proxy.IngressEndpoint.Port>): void;
        addPorts(value?: GlooInstance.GlooInstanceSpec.Proxy.IngressEndpoint.Port, index?: number): GlooInstance.GlooInstanceSpec.Proxy.IngressEndpoint.Port;

        getServiceName(): string;
        setServiceName(value: string): void;

        serializeBinary(): Uint8Array;
        toObject(includeInstance?: boolean): IngressEndpoint.AsObject;
        static toObject(includeInstance: boolean, msg: IngressEndpoint): IngressEndpoint.AsObject;
        static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
        static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
        static serializeBinaryToWriter(message: IngressEndpoint, writer: jspb.BinaryWriter): void;
        static deserializeBinary(bytes: Uint8Array): IngressEndpoint;
        static deserializeBinaryFromReader(message: IngressEndpoint, reader: jspb.BinaryReader): IngressEndpoint;
      }

      export namespace IngressEndpoint {
        export type AsObject = {
          address: string,
          portsList: Array<GlooInstance.GlooInstanceSpec.Proxy.IngressEndpoint.Port.AsObject>,
          serviceName: string,
        }

        export class Port extends jspb.Message {
          getPort(): number;
          setPort(value: number): void;

          getName(): string;
          setName(value: string): void;

          serializeBinary(): Uint8Array;
          toObject(includeInstance?: boolean): Port.AsObject;
          static toObject(includeInstance: boolean, msg: Port): Port.AsObject;
          static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
          static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
          static serializeBinaryToWriter(message: Port, writer: jspb.BinaryWriter): void;
          static deserializeBinary(bytes: Uint8Array): Port;
          static deserializeBinaryFromReader(message: Port, reader: jspb.BinaryReader): Port;
        }

        export namespace Port {
          export type AsObject = {
            port: number,
            name: string,
          }
        }
      }

      export interface WorkloadControllerMap {
        UNDEFINED: 0;
        DEPLOYMENT: 1;
        DAEMON_SET: 2;
      }

      export const WorkloadController: WorkloadControllerMap;
    }

    export class Check extends jspb.Message {
      hasGateways(): boolean;
      clearGateways(): void;
      getGateways(): GlooInstance.GlooInstanceSpec.Check.Summary | undefined;
      setGateways(value?: GlooInstance.GlooInstanceSpec.Check.Summary): void;

      hasVirtualServices(): boolean;
      clearVirtualServices(): void;
      getVirtualServices(): GlooInstance.GlooInstanceSpec.Check.Summary | undefined;
      setVirtualServices(value?: GlooInstance.GlooInstanceSpec.Check.Summary): void;

      hasRouteTables(): boolean;
      clearRouteTables(): void;
      getRouteTables(): GlooInstance.GlooInstanceSpec.Check.Summary | undefined;
      setRouteTables(value?: GlooInstance.GlooInstanceSpec.Check.Summary): void;

      hasAuthConfigs(): boolean;
      clearAuthConfigs(): void;
      getAuthConfigs(): GlooInstance.GlooInstanceSpec.Check.Summary | undefined;
      setAuthConfigs(value?: GlooInstance.GlooInstanceSpec.Check.Summary): void;

      hasSettings(): boolean;
      clearSettings(): void;
      getSettings(): GlooInstance.GlooInstanceSpec.Check.Summary | undefined;
      setSettings(value?: GlooInstance.GlooInstanceSpec.Check.Summary): void;

      hasUpstreams(): boolean;
      clearUpstreams(): void;
      getUpstreams(): GlooInstance.GlooInstanceSpec.Check.Summary | undefined;
      setUpstreams(value?: GlooInstance.GlooInstanceSpec.Check.Summary): void;

      hasUpstreamGroups(): boolean;
      clearUpstreamGroups(): void;
      getUpstreamGroups(): GlooInstance.GlooInstanceSpec.Check.Summary | undefined;
      setUpstreamGroups(value?: GlooInstance.GlooInstanceSpec.Check.Summary): void;

      hasProxies(): boolean;
      clearProxies(): void;
      getProxies(): GlooInstance.GlooInstanceSpec.Check.Summary | undefined;
      setProxies(value?: GlooInstance.GlooInstanceSpec.Check.Summary): void;

      hasDeployments(): boolean;
      clearDeployments(): void;
      getDeployments(): GlooInstance.GlooInstanceSpec.Check.Summary | undefined;
      setDeployments(value?: GlooInstance.GlooInstanceSpec.Check.Summary): void;

      hasPods(): boolean;
      clearPods(): void;
      getPods(): GlooInstance.GlooInstanceSpec.Check.Summary | undefined;
      setPods(value?: GlooInstance.GlooInstanceSpec.Check.Summary): void;

      serializeBinary(): Uint8Array;
      toObject(includeInstance?: boolean): Check.AsObject;
      static toObject(includeInstance: boolean, msg: Check): Check.AsObject;
      static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
      static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
      static serializeBinaryToWriter(message: Check, writer: jspb.BinaryWriter): void;
      static deserializeBinary(bytes: Uint8Array): Check;
      static deserializeBinaryFromReader(message: Check, reader: jspb.BinaryReader): Check;
    }

    export namespace Check {
      export type AsObject = {
        gateways?: GlooInstance.GlooInstanceSpec.Check.Summary.AsObject,
        virtualServices?: GlooInstance.GlooInstanceSpec.Check.Summary.AsObject,
        routeTables?: GlooInstance.GlooInstanceSpec.Check.Summary.AsObject,
        authConfigs?: GlooInstance.GlooInstanceSpec.Check.Summary.AsObject,
        settings?: GlooInstance.GlooInstanceSpec.Check.Summary.AsObject,
        upstreams?: GlooInstance.GlooInstanceSpec.Check.Summary.AsObject,
        upstreamGroups?: GlooInstance.GlooInstanceSpec.Check.Summary.AsObject,
        proxies?: GlooInstance.GlooInstanceSpec.Check.Summary.AsObject,
        deployments?: GlooInstance.GlooInstanceSpec.Check.Summary.AsObject,
        pods?: GlooInstance.GlooInstanceSpec.Check.Summary.AsObject,
      }

      export class Summary extends jspb.Message {
        getTotal(): number;
        setTotal(value: number): void;

        clearErrorsList(): void;
        getErrorsList(): Array<GlooInstance.GlooInstanceSpec.Check.Summary.ResourceReport>;
        setErrorsList(value: Array<GlooInstance.GlooInstanceSpec.Check.Summary.ResourceReport>): void;
        addErrors(value?: GlooInstance.GlooInstanceSpec.Check.Summary.ResourceReport, index?: number): GlooInstance.GlooInstanceSpec.Check.Summary.ResourceReport;

        clearWarningsList(): void;
        getWarningsList(): Array<GlooInstance.GlooInstanceSpec.Check.Summary.ResourceReport>;
        setWarningsList(value: Array<GlooInstance.GlooInstanceSpec.Check.Summary.ResourceReport>): void;
        addWarnings(value?: GlooInstance.GlooInstanceSpec.Check.Summary.ResourceReport, index?: number): GlooInstance.GlooInstanceSpec.Check.Summary.ResourceReport;

        serializeBinary(): Uint8Array;
        toObject(includeInstance?: boolean): Summary.AsObject;
        static toObject(includeInstance: boolean, msg: Summary): Summary.AsObject;
        static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
        static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
        static serializeBinaryToWriter(message: Summary, writer: jspb.BinaryWriter): void;
        static deserializeBinary(bytes: Uint8Array): Summary;
        static deserializeBinaryFromReader(message: Summary, reader: jspb.BinaryReader): Summary;
      }

      export namespace Summary {
        export type AsObject = {
          total: number,
          errorsList: Array<GlooInstance.GlooInstanceSpec.Check.Summary.ResourceReport.AsObject>,
          warningsList: Array<GlooInstance.GlooInstanceSpec.Check.Summary.ResourceReport.AsObject>,
        }

        export class ResourceReport extends jspb.Message {
          hasRef(): boolean;
          clearRef(): void;
          getRef(): github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef | undefined;
          setRef(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef): void;

          getMessage(): string;
          setMessage(value: string): void;

          serializeBinary(): Uint8Array;
          toObject(includeInstance?: boolean): ResourceReport.AsObject;
          static toObject(includeInstance: boolean, msg: ResourceReport): ResourceReport.AsObject;
          static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
          static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
          static serializeBinaryToWriter(message: ResourceReport, writer: jspb.BinaryWriter): void;
          static deserializeBinary(bytes: Uint8Array): ResourceReport;
          static deserializeBinaryFromReader(message: ResourceReport, reader: jspb.BinaryReader): ResourceReport;
        }

        export namespace ResourceReport {
          export type AsObject = {
            ref?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef.AsObject,
            message: string,
          }
        }
      }
    }
  }

  export class GlooInstanceStatus extends jspb.Message {
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): GlooInstanceStatus.AsObject;
    static toObject(includeInstance: boolean, msg: GlooInstanceStatus): GlooInstanceStatus.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: GlooInstanceStatus, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): GlooInstanceStatus;
    static deserializeBinaryFromReader(message: GlooInstanceStatus, reader: jspb.BinaryReader): GlooInstanceStatus;
  }

  export namespace GlooInstanceStatus {
    export type AsObject = {
    }
  }
}

export class ListGlooInstancesRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListGlooInstancesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListGlooInstancesRequest): ListGlooInstancesRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListGlooInstancesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListGlooInstancesRequest;
  static deserializeBinaryFromReader(message: ListGlooInstancesRequest, reader: jspb.BinaryReader): ListGlooInstancesRequest;
}

export namespace ListGlooInstancesRequest {
  export type AsObject = {
  }
}

export class ListGlooInstancesResponse extends jspb.Message {
  clearGlooInstancesList(): void;
  getGlooInstancesList(): Array<GlooInstance>;
  setGlooInstancesList(value: Array<GlooInstance>): void;
  addGlooInstances(value?: GlooInstance, index?: number): GlooInstance;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListGlooInstancesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListGlooInstancesResponse): ListGlooInstancesResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListGlooInstancesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListGlooInstancesResponse;
  static deserializeBinaryFromReader(message: ListGlooInstancesResponse, reader: jspb.BinaryReader): ListGlooInstancesResponse;
}

export namespace ListGlooInstancesResponse {
  export type AsObject = {
    glooInstancesList: Array<GlooInstance.AsObject>,
  }
}

export class ClusterDetails extends jspb.Message {
  getCluster(): string;
  setCluster(value: string): void;

  clearGlooInstancesList(): void;
  getGlooInstancesList(): Array<GlooInstance>;
  setGlooInstancesList(value: Array<GlooInstance>): void;
  addGlooInstances(value?: GlooInstance, index?: number): GlooInstance;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ClusterDetails.AsObject;
  static toObject(includeInstance: boolean, msg: ClusterDetails): ClusterDetails.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ClusterDetails, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ClusterDetails;
  static deserializeBinaryFromReader(message: ClusterDetails, reader: jspb.BinaryReader): ClusterDetails;
}

export namespace ClusterDetails {
  export type AsObject = {
    cluster: string,
    glooInstancesList: Array<GlooInstance.AsObject>,
  }
}

export class ListClusterDetailsRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListClusterDetailsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListClusterDetailsRequest): ListClusterDetailsRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListClusterDetailsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListClusterDetailsRequest;
  static deserializeBinaryFromReader(message: ListClusterDetailsRequest, reader: jspb.BinaryReader): ListClusterDetailsRequest;
}

export namespace ListClusterDetailsRequest {
  export type AsObject = {
  }
}

export class ListClusterDetailsResponse extends jspb.Message {
  clearClusterDetailsList(): void;
  getClusterDetailsList(): Array<ClusterDetails>;
  setClusterDetailsList(value: Array<ClusterDetails>): void;
  addClusterDetails(value?: ClusterDetails, index?: number): ClusterDetails;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListClusterDetailsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListClusterDetailsResponse): ListClusterDetailsResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListClusterDetailsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListClusterDetailsResponse;
  static deserializeBinaryFromReader(message: ListClusterDetailsResponse, reader: jspb.BinaryReader): ListClusterDetailsResponse;
}

export namespace ListClusterDetailsResponse {
  export type AsObject = {
    clusterDetailsList: Array<ClusterDetails.AsObject>,
  }
}

export class ConfigDump extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getRaw(): string;
  setRaw(value: string): void;

  getError(): string;
  setError(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ConfigDump.AsObject;
  static toObject(includeInstance: boolean, msg: ConfigDump): ConfigDump.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ConfigDump, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ConfigDump;
  static deserializeBinaryFromReader(message: ConfigDump, reader: jspb.BinaryReader): ConfigDump;
}

export namespace ConfigDump {
  export type AsObject = {
    name: string,
    raw: string,
    error: string,
  }
}

export class GetConfigDumpsRequest extends jspb.Message {
  hasGlooInstanceRef(): boolean;
  clearGlooInstanceRef(): void;
  getGlooInstanceRef(): github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef | undefined;
  setGlooInstanceRef(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetConfigDumpsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetConfigDumpsRequest): GetConfigDumpsRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetConfigDumpsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetConfigDumpsRequest;
  static deserializeBinaryFromReader(message: GetConfigDumpsRequest, reader: jspb.BinaryReader): GetConfigDumpsRequest;
}

export namespace GetConfigDumpsRequest {
  export type AsObject = {
    glooInstanceRef?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef.AsObject,
  }
}

export class GetConfigDumpsResponse extends jspb.Message {
  clearConfigDumpsList(): void;
  getConfigDumpsList(): Array<ConfigDump>;
  setConfigDumpsList(value: Array<ConfigDump>): void;
  addConfigDumps(value?: ConfigDump, index?: number): ConfigDump;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetConfigDumpsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetConfigDumpsResponse): GetConfigDumpsResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetConfigDumpsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetConfigDumpsResponse;
  static deserializeBinaryFromReader(message: GetConfigDumpsResponse, reader: jspb.BinaryReader): GetConfigDumpsResponse;
}

export namespace GetConfigDumpsResponse {
  export type AsObject = {
    configDumpsList: Array<ConfigDump.AsObject>,
  }
}
