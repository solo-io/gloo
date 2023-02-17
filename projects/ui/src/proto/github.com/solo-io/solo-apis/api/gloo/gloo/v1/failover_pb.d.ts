/* eslint-disable */
// package: gloo.solo.io
// file: github.com/solo-io/solo-apis/api/gloo/gloo/v1/failover.proto

import * as jspb from "google-protobuf";
import * as extproto_ext_pb from "../../../../../../../extproto/ext_pb";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_ssl_ssl_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/ssl/ssl_pb";
import * as validate_validate_pb from "../../../../../../../validate/validate_pb";

export class Failover extends jspb.Message {
  clearPrioritizedLocalitiesList(): void;
  getPrioritizedLocalitiesList(): Array<Failover.PrioritizedLocality>;
  setPrioritizedLocalitiesList(value: Array<Failover.PrioritizedLocality>): void;
  addPrioritizedLocalities(value?: Failover.PrioritizedLocality, index?: number): Failover.PrioritizedLocality;

  hasPolicy(): boolean;
  clearPolicy(): void;
  getPolicy(): Failover.Policy | undefined;
  setPolicy(value?: Failover.Policy): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Failover.AsObject;
  static toObject(includeInstance: boolean, msg: Failover): Failover.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Failover, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Failover;
  static deserializeBinaryFromReader(message: Failover, reader: jspb.BinaryReader): Failover;
}

export namespace Failover {
  export type AsObject = {
    prioritizedLocalitiesList: Array<Failover.PrioritizedLocality.AsObject>,
    policy?: Failover.Policy.AsObject,
  }

  export class PrioritizedLocality extends jspb.Message {
    clearLocalityEndpointsList(): void;
    getLocalityEndpointsList(): Array<LocalityLbEndpoints>;
    setLocalityEndpointsList(value: Array<LocalityLbEndpoints>): void;
    addLocalityEndpoints(value?: LocalityLbEndpoints, index?: number): LocalityLbEndpoints;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): PrioritizedLocality.AsObject;
    static toObject(includeInstance: boolean, msg: PrioritizedLocality): PrioritizedLocality.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: PrioritizedLocality, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): PrioritizedLocality;
    static deserializeBinaryFromReader(message: PrioritizedLocality, reader: jspb.BinaryReader): PrioritizedLocality;
  }

  export namespace PrioritizedLocality {
    export type AsObject = {
      localityEndpointsList: Array<LocalityLbEndpoints.AsObject>,
    }
  }

  export class Policy extends jspb.Message {
    hasOverprovisioningFactor(): boolean;
    clearOverprovisioningFactor(): void;
    getOverprovisioningFactor(): google_protobuf_wrappers_pb.UInt32Value | undefined;
    setOverprovisioningFactor(value?: google_protobuf_wrappers_pb.UInt32Value): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Policy.AsObject;
    static toObject(includeInstance: boolean, msg: Policy): Policy.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Policy, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Policy;
    static deserializeBinaryFromReader(message: Policy, reader: jspb.BinaryReader): Policy;
  }

  export namespace Policy {
    export type AsObject = {
      overprovisioningFactor?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
    }
  }
}

export class LocalityLbEndpoints extends jspb.Message {
  hasLocality(): boolean;
  clearLocality(): void;
  getLocality(): Locality | undefined;
  setLocality(value?: Locality): void;

  clearLbEndpointsList(): void;
  getLbEndpointsList(): Array<LbEndpoint>;
  setLbEndpointsList(value: Array<LbEndpoint>): void;
  addLbEndpoints(value?: LbEndpoint, index?: number): LbEndpoint;

  hasLoadBalancingWeight(): boolean;
  clearLoadBalancingWeight(): void;
  getLoadBalancingWeight(): google_protobuf_wrappers_pb.UInt32Value | undefined;
  setLoadBalancingWeight(value?: google_protobuf_wrappers_pb.UInt32Value): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LocalityLbEndpoints.AsObject;
  static toObject(includeInstance: boolean, msg: LocalityLbEndpoints): LocalityLbEndpoints.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: LocalityLbEndpoints, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LocalityLbEndpoints;
  static deserializeBinaryFromReader(message: LocalityLbEndpoints, reader: jspb.BinaryReader): LocalityLbEndpoints;
}

export namespace LocalityLbEndpoints {
  export type AsObject = {
    locality?: Locality.AsObject,
    lbEndpointsList: Array<LbEndpoint.AsObject>,
    loadBalancingWeight?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
  }
}

export class LbEndpoint extends jspb.Message {
  getAddress(): string;
  setAddress(value: string): void;

  getPort(): number;
  setPort(value: number): void;

  hasHealthCheckConfig(): boolean;
  clearHealthCheckConfig(): void;
  getHealthCheckConfig(): LbEndpoint.HealthCheckConfig | undefined;
  setHealthCheckConfig(value?: LbEndpoint.HealthCheckConfig): void;

  hasUpstreamSslConfig(): boolean;
  clearUpstreamSslConfig(): void;
  getUpstreamSslConfig(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_ssl_ssl_pb.UpstreamSslConfig | undefined;
  setUpstreamSslConfig(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_ssl_ssl_pb.UpstreamSslConfig): void;

  hasLoadBalancingWeight(): boolean;
  clearLoadBalancingWeight(): void;
  getLoadBalancingWeight(): google_protobuf_wrappers_pb.UInt32Value | undefined;
  setLoadBalancingWeight(value?: google_protobuf_wrappers_pb.UInt32Value): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LbEndpoint.AsObject;
  static toObject(includeInstance: boolean, msg: LbEndpoint): LbEndpoint.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: LbEndpoint, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LbEndpoint;
  static deserializeBinaryFromReader(message: LbEndpoint, reader: jspb.BinaryReader): LbEndpoint;
}

export namespace LbEndpoint {
  export type AsObject = {
    address: string,
    port: number,
    healthCheckConfig?: LbEndpoint.HealthCheckConfig.AsObject,
    upstreamSslConfig?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_ssl_ssl_pb.UpstreamSslConfig.AsObject,
    loadBalancingWeight?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
  }

  export class HealthCheckConfig extends jspb.Message {
    getPortValue(): number;
    setPortValue(value: number): void;

    getHostname(): string;
    setHostname(value: string): void;

    getPath(): string;
    setPath(value: string): void;

    getMethod(): string;
    setMethod(value: string): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): HealthCheckConfig.AsObject;
    static toObject(includeInstance: boolean, msg: HealthCheckConfig): HealthCheckConfig.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: HealthCheckConfig, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): HealthCheckConfig;
    static deserializeBinaryFromReader(message: HealthCheckConfig, reader: jspb.BinaryReader): HealthCheckConfig;
  }

  export namespace HealthCheckConfig {
    export type AsObject = {
      portValue: number,
      hostname: string,
      path: string,
      method: string,
    }
  }
}

export class Locality extends jspb.Message {
  getRegion(): string;
  setRegion(value: string): void;

  getZone(): string;
  setZone(value: string): void;

  getSubZone(): string;
  setSubZone(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Locality.AsObject;
  static toObject(includeInstance: boolean, msg: Locality): Locality.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Locality, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Locality;
  static deserializeBinaryFromReader(message: Locality, reader: jspb.BinaryReader): Locality;
}

export namespace Locality {
  export type AsObject = {
    region: string,
    zone: string,
    subZone: string,
  }
}
