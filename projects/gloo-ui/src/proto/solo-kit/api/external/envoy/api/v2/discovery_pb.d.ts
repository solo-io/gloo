// package: envoy.api.v2
// file: solo-kit/api/external/envoy/api/v2/discovery.proto

import * as jspb from "google-protobuf";
import * as envoy_api_v2_core_base_pb from "../../../../../../envoy/api/v2/core/base_pb";
import * as google_protobuf_any_pb from "google-protobuf/google/protobuf/any_pb";
import * as google_rpc_status_pb from "../../../../../../google/rpc/status_pb";
import * as gogoproto_gogo_pb from "../../../../../../gogoproto/gogo_pb";
import * as extproto_ext_pb from "../../../../../../extproto/ext_pb";

export class DiscoveryRequest extends jspb.Message {
  getVersionInfo(): string;
  setVersionInfo(value: string): void;

  hasNode(): boolean;
  clearNode(): void;
  getNode(): envoy_api_v2_core_base_pb.Node | undefined;
  setNode(value?: envoy_api_v2_core_base_pb.Node): void;

  clearResourceNamesList(): void;
  getResourceNamesList(): Array<string>;
  setResourceNamesList(value: Array<string>): void;
  addResourceNames(value: string, index?: number): string;

  getTypeUrl(): string;
  setTypeUrl(value: string): void;

  getResponseNonce(): string;
  setResponseNonce(value: string): void;

  hasErrorDetail(): boolean;
  clearErrorDetail(): void;
  getErrorDetail(): google_rpc_status_pb.Status | undefined;
  setErrorDetail(value?: google_rpc_status_pb.Status): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DiscoveryRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DiscoveryRequest): DiscoveryRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DiscoveryRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DiscoveryRequest;
  static deserializeBinaryFromReader(message: DiscoveryRequest, reader: jspb.BinaryReader): DiscoveryRequest;
}

export namespace DiscoveryRequest {
  export type AsObject = {
    versionInfo: string,
    node?: envoy_api_v2_core_base_pb.Node.AsObject,
    resourceNamesList: Array<string>,
    typeUrl: string,
    responseNonce: string,
    errorDetail?: google_rpc_status_pb.Status.AsObject,
  }
}

export class DiscoveryResponse extends jspb.Message {
  getVersionInfo(): string;
  setVersionInfo(value: string): void;

  clearResourcesList(): void;
  getResourcesList(): Array<google_protobuf_any_pb.Any>;
  setResourcesList(value: Array<google_protobuf_any_pb.Any>): void;
  addResources(value?: google_protobuf_any_pb.Any, index?: number): google_protobuf_any_pb.Any;

  getCanary(): boolean;
  setCanary(value: boolean): void;

  getTypeUrl(): string;
  setTypeUrl(value: string): void;

  getNonce(): string;
  setNonce(value: string): void;

  hasControlPlane(): boolean;
  clearControlPlane(): void;
  getControlPlane(): envoy_api_v2_core_base_pb.ControlPlane | undefined;
  setControlPlane(value?: envoy_api_v2_core_base_pb.ControlPlane): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DiscoveryResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DiscoveryResponse): DiscoveryResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DiscoveryResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DiscoveryResponse;
  static deserializeBinaryFromReader(message: DiscoveryResponse, reader: jspb.BinaryReader): DiscoveryResponse;
}

export namespace DiscoveryResponse {
  export type AsObject = {
    versionInfo: string,
    resourcesList: Array<google_protobuf_any_pb.Any.AsObject>,
    canary: boolean,
    typeUrl: string,
    nonce: string,
    controlPlane?: envoy_api_v2_core_base_pb.ControlPlane.AsObject,
  }
}

export class DeltaDiscoveryRequest extends jspb.Message {
  hasNode(): boolean;
  clearNode(): void;
  getNode(): envoy_api_v2_core_base_pb.Node | undefined;
  setNode(value?: envoy_api_v2_core_base_pb.Node): void;

  getTypeUrl(): string;
  setTypeUrl(value: string): void;

  clearResourceNamesSubscribeList(): void;
  getResourceNamesSubscribeList(): Array<string>;
  setResourceNamesSubscribeList(value: Array<string>): void;
  addResourceNamesSubscribe(value: string, index?: number): string;

  clearResourceNamesUnsubscribeList(): void;
  getResourceNamesUnsubscribeList(): Array<string>;
  setResourceNamesUnsubscribeList(value: Array<string>): void;
  addResourceNamesUnsubscribe(value: string, index?: number): string;

  getInitialResourceVersionsMap(): jspb.Map<string, string>;
  clearInitialResourceVersionsMap(): void;
  getResponseNonce(): string;
  setResponseNonce(value: string): void;

  hasErrorDetail(): boolean;
  clearErrorDetail(): void;
  getErrorDetail(): google_rpc_status_pb.Status | undefined;
  setErrorDetail(value?: google_rpc_status_pb.Status): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeltaDiscoveryRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeltaDiscoveryRequest): DeltaDiscoveryRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeltaDiscoveryRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeltaDiscoveryRequest;
  static deserializeBinaryFromReader(message: DeltaDiscoveryRequest, reader: jspb.BinaryReader): DeltaDiscoveryRequest;
}

export namespace DeltaDiscoveryRequest {
  export type AsObject = {
    node?: envoy_api_v2_core_base_pb.Node.AsObject,
    typeUrl: string,
    resourceNamesSubscribeList: Array<string>,
    resourceNamesUnsubscribeList: Array<string>,
    initialResourceVersionsMap: Array<[string, string]>,
    responseNonce: string,
    errorDetail?: google_rpc_status_pb.Status.AsObject,
  }
}

export class DeltaDiscoveryResponse extends jspb.Message {
  getSystemVersionInfo(): string;
  setSystemVersionInfo(value: string): void;

  clearResourcesList(): void;
  getResourcesList(): Array<Resource>;
  setResourcesList(value: Array<Resource>): void;
  addResources(value?: Resource, index?: number): Resource;

  getTypeUrl(): string;
  setTypeUrl(value: string): void;

  clearRemovedResourcesList(): void;
  getRemovedResourcesList(): Array<string>;
  setRemovedResourcesList(value: Array<string>): void;
  addRemovedResources(value: string, index?: number): string;

  getNonce(): string;
  setNonce(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeltaDiscoveryResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeltaDiscoveryResponse): DeltaDiscoveryResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeltaDiscoveryResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeltaDiscoveryResponse;
  static deserializeBinaryFromReader(message: DeltaDiscoveryResponse, reader: jspb.BinaryReader): DeltaDiscoveryResponse;
}

export namespace DeltaDiscoveryResponse {
  export type AsObject = {
    systemVersionInfo: string,
    resourcesList: Array<Resource.AsObject>,
    typeUrl: string,
    removedResourcesList: Array<string>,
    nonce: string,
  }
}

export class Resource extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  clearAliasesList(): void;
  getAliasesList(): Array<string>;
  setAliasesList(value: Array<string>): void;
  addAliases(value: string, index?: number): string;

  getVersion(): string;
  setVersion(value: string): void;

  hasResource(): boolean;
  clearResource(): void;
  getResource(): google_protobuf_any_pb.Any | undefined;
  setResource(value?: google_protobuf_any_pb.Any): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Resource.AsObject;
  static toObject(includeInstance: boolean, msg: Resource): Resource.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Resource, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Resource;
  static deserializeBinaryFromReader(message: Resource, reader: jspb.BinaryReader): Resource;
}

export namespace Resource {
  export type AsObject = {
    name: string,
    aliasesList: Array<string>,
    version: string,
    resource?: google_protobuf_any_pb.Any.AsObject,
  }
}

