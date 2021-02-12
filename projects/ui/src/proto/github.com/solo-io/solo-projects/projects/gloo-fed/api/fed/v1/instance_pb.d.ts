/* eslint-disable */
// package: fed.solo.io
// file: github.com/solo-io/solo-projects/projects/gloo-fed/api/fed/v1/instance.proto

import * as jspb from "google-protobuf";
import * as github_com_solo_io_skv2_api_core_v1_core_pb from "../../../../../../../../github.com/solo-io/skv2/api/core/v1/core_pb";
import * as extproto_ext_pb from "../../../../../../../../extproto/ext_pb";

export class GlooInstanceSpec extends jspb.Message {
  getCluster(): string;
  setCluster(value: string): void;

  getIsEnterprise(): boolean;
  setIsEnterprise(value: boolean): void;

  hasControlPlane(): boolean;
  clearControlPlane(): void;
  getControlPlane(): GlooInstanceSpec.ControlPlane | undefined;
  setControlPlane(value?: GlooInstanceSpec.ControlPlane): void;

  clearProxiesList(): void;
  getProxiesList(): Array<GlooInstanceSpec.Proxy>;
  setProxiesList(value: Array<GlooInstanceSpec.Proxy>): void;
  addProxies(value?: GlooInstanceSpec.Proxy, index?: number): GlooInstanceSpec.Proxy;

  getRegion(): string;
  setRegion(value: string): void;

  hasAdmin(): boolean;
  clearAdmin(): void;
  getAdmin(): GlooInstanceSpec.Admin | undefined;
  setAdmin(value?: GlooInstanceSpec.Admin): void;

  hasCheck(): boolean;
  clearCheck(): void;
  getCheck(): GlooInstanceSpec.Check | undefined;
  setCheck(value?: GlooInstanceSpec.Check): void;

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
    controlPlane?: GlooInstanceSpec.ControlPlane.AsObject,
    proxiesList: Array<GlooInstanceSpec.Proxy.AsObject>,
    region: string,
    admin?: GlooInstanceSpec.Admin.AsObject,
    check?: GlooInstanceSpec.Check.AsObject,
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

    getWorkloadControllerType(): GlooInstanceSpec.Proxy.WorkloadControllerMap[keyof GlooInstanceSpec.Proxy.WorkloadControllerMap];
    setWorkloadControllerType(value: GlooInstanceSpec.Proxy.WorkloadControllerMap[keyof GlooInstanceSpec.Proxy.WorkloadControllerMap]): void;

    clearZonesList(): void;
    getZonesList(): Array<string>;
    setZonesList(value: Array<string>): void;
    addZones(value: string, index?: number): string;

    clearIngressEndpointsList(): void;
    getIngressEndpointsList(): Array<GlooInstanceSpec.Proxy.IngressEndpoint>;
    setIngressEndpointsList(value: Array<GlooInstanceSpec.Proxy.IngressEndpoint>): void;
    addIngressEndpoints(value?: GlooInstanceSpec.Proxy.IngressEndpoint, index?: number): GlooInstanceSpec.Proxy.IngressEndpoint;

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
      workloadControllerType: GlooInstanceSpec.Proxy.WorkloadControllerMap[keyof GlooInstanceSpec.Proxy.WorkloadControllerMap],
      zonesList: Array<string>,
      ingressEndpointsList: Array<GlooInstanceSpec.Proxy.IngressEndpoint.AsObject>,
    }

    export class IngressEndpoint extends jspb.Message {
      getAddress(): string;
      setAddress(value: string): void;

      clearPortsList(): void;
      getPortsList(): Array<GlooInstanceSpec.Proxy.IngressEndpoint.Port>;
      setPortsList(value: Array<GlooInstanceSpec.Proxy.IngressEndpoint.Port>): void;
      addPorts(value?: GlooInstanceSpec.Proxy.IngressEndpoint.Port, index?: number): GlooInstanceSpec.Proxy.IngressEndpoint.Port;

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
        portsList: Array<GlooInstanceSpec.Proxy.IngressEndpoint.Port.AsObject>,
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

  export class Admin extends jspb.Message {
    getWriteNamespace(): string;
    setWriteNamespace(value: string): void;

    hasProxyId(): boolean;
    clearProxyId(): void;
    getProxyId(): github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef | undefined;
    setProxyId(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Admin.AsObject;
    static toObject(includeInstance: boolean, msg: Admin): Admin.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Admin, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Admin;
    static deserializeBinaryFromReader(message: Admin, reader: jspb.BinaryReader): Admin;
  }

  export namespace Admin {
    export type AsObject = {
      writeNamespace: string,
      proxyId?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef.AsObject,
    }
  }

  export class Check extends jspb.Message {
    hasGateways(): boolean;
    clearGateways(): void;
    getGateways(): GlooInstanceSpec.Check.Summary | undefined;
    setGateways(value?: GlooInstanceSpec.Check.Summary): void;

    hasVirtualServices(): boolean;
    clearVirtualServices(): void;
    getVirtualServices(): GlooInstanceSpec.Check.Summary | undefined;
    setVirtualServices(value?: GlooInstanceSpec.Check.Summary): void;

    hasRouteTables(): boolean;
    clearRouteTables(): void;
    getRouteTables(): GlooInstanceSpec.Check.Summary | undefined;
    setRouteTables(value?: GlooInstanceSpec.Check.Summary): void;

    hasAuthConfigs(): boolean;
    clearAuthConfigs(): void;
    getAuthConfigs(): GlooInstanceSpec.Check.Summary | undefined;
    setAuthConfigs(value?: GlooInstanceSpec.Check.Summary): void;

    hasSettings(): boolean;
    clearSettings(): void;
    getSettings(): GlooInstanceSpec.Check.Summary | undefined;
    setSettings(value?: GlooInstanceSpec.Check.Summary): void;

    hasUpstreams(): boolean;
    clearUpstreams(): void;
    getUpstreams(): GlooInstanceSpec.Check.Summary | undefined;
    setUpstreams(value?: GlooInstanceSpec.Check.Summary): void;

    hasUpstreamGroups(): boolean;
    clearUpstreamGroups(): void;
    getUpstreamGroups(): GlooInstanceSpec.Check.Summary | undefined;
    setUpstreamGroups(value?: GlooInstanceSpec.Check.Summary): void;

    hasProxies(): boolean;
    clearProxies(): void;
    getProxies(): GlooInstanceSpec.Check.Summary | undefined;
    setProxies(value?: GlooInstanceSpec.Check.Summary): void;

    hasDeployments(): boolean;
    clearDeployments(): void;
    getDeployments(): GlooInstanceSpec.Check.Summary | undefined;
    setDeployments(value?: GlooInstanceSpec.Check.Summary): void;

    hasPods(): boolean;
    clearPods(): void;
    getPods(): GlooInstanceSpec.Check.Summary | undefined;
    setPods(value?: GlooInstanceSpec.Check.Summary): void;

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
      gateways?: GlooInstanceSpec.Check.Summary.AsObject,
      virtualServices?: GlooInstanceSpec.Check.Summary.AsObject,
      routeTables?: GlooInstanceSpec.Check.Summary.AsObject,
      authConfigs?: GlooInstanceSpec.Check.Summary.AsObject,
      settings?: GlooInstanceSpec.Check.Summary.AsObject,
      upstreams?: GlooInstanceSpec.Check.Summary.AsObject,
      upstreamGroups?: GlooInstanceSpec.Check.Summary.AsObject,
      proxies?: GlooInstanceSpec.Check.Summary.AsObject,
      deployments?: GlooInstanceSpec.Check.Summary.AsObject,
      pods?: GlooInstanceSpec.Check.Summary.AsObject,
    }

    export class Summary extends jspb.Message {
      getTotal(): number;
      setTotal(value: number): void;

      clearErrorsList(): void;
      getErrorsList(): Array<GlooInstanceSpec.Check.Summary.ResourceReport>;
      setErrorsList(value: Array<GlooInstanceSpec.Check.Summary.ResourceReport>): void;
      addErrors(value?: GlooInstanceSpec.Check.Summary.ResourceReport, index?: number): GlooInstanceSpec.Check.Summary.ResourceReport;

      clearWarningsList(): void;
      getWarningsList(): Array<GlooInstanceSpec.Check.Summary.ResourceReport>;
      setWarningsList(value: Array<GlooInstanceSpec.Check.Summary.ResourceReport>): void;
      addWarnings(value?: GlooInstanceSpec.Check.Summary.ResourceReport, index?: number): GlooInstanceSpec.Check.Summary.ResourceReport;

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
        errorsList: Array<GlooInstanceSpec.Check.Summary.ResourceReport.AsObject>,
        warningsList: Array<GlooInstanceSpec.Check.Summary.ResourceReport.AsObject>,
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
