/* eslint-disable */
// package: solo.io.envoy.api.v2.core
// file: github.com/solo-io/solo-kit/api/external/envoy/api/v2/core/base.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_any_pb from "google-protobuf/google/protobuf/any_pb";
import * as google_protobuf_struct_pb from "google-protobuf/google/protobuf/struct_pb";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";
import * as github_com_solo_io_solo_kit_api_external_envoy_api_v2_core_address_pb from "../../../../../../../../../github.com/solo-io/solo-kit/api/external/envoy/api/v2/core/address_pb";
import * as github_com_solo_io_solo_kit_api_external_envoy_api_v2_core_http_uri_pb from "../../../../../../../../../github.com/solo-io/solo-kit/api/external/envoy/api/v2/core/http_uri_pb";
import * as github_com_solo_io_solo_kit_api_external_envoy_type_percent_pb from "../../../../../../../../../github.com/solo-io/solo-kit/api/external/envoy/type/percent_pb";
import * as github_com_solo_io_solo_kit_api_external_envoy_type_semantic_version_pb from "../../../../../../../../../github.com/solo-io/solo-kit/api/external/envoy/type/semantic_version_pb";
import * as validate_validate_pb from "../../../../../../../../../validate/validate_pb";
import * as extproto_ext_pb from "../../../../../../../../../extproto/ext_pb";
import * as github_com_solo_io_solo_kit_api_v1_ref_pb from "../../../../../../../../../github.com/solo-io/solo-kit/api/v1/ref_pb";

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

export class BuildVersion extends jspb.Message {
  hasVersion(): boolean;
  clearVersion(): void;
  getVersion(): github_com_solo_io_solo_kit_api_external_envoy_type_semantic_version_pb.SemanticVersion | undefined;
  setVersion(value?: github_com_solo_io_solo_kit_api_external_envoy_type_semantic_version_pb.SemanticVersion): void;

  hasMetadata(): boolean;
  clearMetadata(): void;
  getMetadata(): google_protobuf_struct_pb.Struct | undefined;
  setMetadata(value?: google_protobuf_struct_pb.Struct): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): BuildVersion.AsObject;
  static toObject(includeInstance: boolean, msg: BuildVersion): BuildVersion.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: BuildVersion, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): BuildVersion;
  static deserializeBinaryFromReader(message: BuildVersion, reader: jspb.BinaryReader): BuildVersion;
}

export namespace BuildVersion {
  export type AsObject = {
    version?: github_com_solo_io_solo_kit_api_external_envoy_type_semantic_version_pb.SemanticVersion.AsObject,
    metadata?: google_protobuf_struct_pb.Struct.AsObject,
  }
}

export class Extension extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getCategory(): string;
  setCategory(value: string): void;

  getTypeDescriptor(): string;
  setTypeDescriptor(value: string): void;

  hasVersion(): boolean;
  clearVersion(): void;
  getVersion(): BuildVersion | undefined;
  setVersion(value?: BuildVersion): void;

  getDisabled(): boolean;
  setDisabled(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Extension.AsObject;
  static toObject(includeInstance: boolean, msg: Extension): Extension.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Extension, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Extension;
  static deserializeBinaryFromReader(message: Extension, reader: jspb.BinaryReader): Extension;
}

export namespace Extension {
  export type AsObject = {
    name: string,
    category: string,
    typeDescriptor: string,
    version?: BuildVersion.AsObject,
    disabled: boolean,
  }
}

export class Node extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getCluster(): string;
  setCluster(value: string): void;

  hasMetadata(): boolean;
  clearMetadata(): void;
  getMetadata(): google_protobuf_struct_pb.Struct | undefined;
  setMetadata(value?: google_protobuf_struct_pb.Struct): void;

  hasLocality(): boolean;
  clearLocality(): void;
  getLocality(): Locality | undefined;
  setLocality(value?: Locality): void;

  getBuildVersion(): string;
  setBuildVersion(value: string): void;

  getUserAgentName(): string;
  setUserAgentName(value: string): void;

  hasUserAgentVersion(): boolean;
  clearUserAgentVersion(): void;
  getUserAgentVersion(): string;
  setUserAgentVersion(value: string): void;

  hasUserAgentBuildVersion(): boolean;
  clearUserAgentBuildVersion(): void;
  getUserAgentBuildVersion(): BuildVersion | undefined;
  setUserAgentBuildVersion(value?: BuildVersion): void;

  clearExtensionsList(): void;
  getExtensionsList(): Array<Extension>;
  setExtensionsList(value: Array<Extension>): void;
  addExtensions(value?: Extension, index?: number): Extension;

  clearClientFeaturesList(): void;
  getClientFeaturesList(): Array<string>;
  setClientFeaturesList(value: Array<string>): void;
  addClientFeatures(value: string, index?: number): string;

  clearListeningAddressesList(): void;
  getListeningAddressesList(): Array<github_com_solo_io_solo_kit_api_external_envoy_api_v2_core_address_pb.Address>;
  setListeningAddressesList(value: Array<github_com_solo_io_solo_kit_api_external_envoy_api_v2_core_address_pb.Address>): void;
  addListeningAddresses(value?: github_com_solo_io_solo_kit_api_external_envoy_api_v2_core_address_pb.Address, index?: number): github_com_solo_io_solo_kit_api_external_envoy_api_v2_core_address_pb.Address;

  getUserAgentVersionTypeCase(): Node.UserAgentVersionTypeCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Node.AsObject;
  static toObject(includeInstance: boolean, msg: Node): Node.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Node, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Node;
  static deserializeBinaryFromReader(message: Node, reader: jspb.BinaryReader): Node;
}

export namespace Node {
  export type AsObject = {
    id: string,
    cluster: string,
    metadata?: google_protobuf_struct_pb.Struct.AsObject,
    locality?: Locality.AsObject,
    buildVersion: string,
    userAgentName: string,
    userAgentVersion: string,
    userAgentBuildVersion?: BuildVersion.AsObject,
    extensionsList: Array<Extension.AsObject>,
    clientFeaturesList: Array<string>,
    listeningAddressesList: Array<github_com_solo_io_solo_kit_api_external_envoy_api_v2_core_address_pb.Address.AsObject>,
  }

  export enum UserAgentVersionTypeCase {
    USER_AGENT_VERSION_TYPE_NOT_SET = 0,
    USER_AGENT_VERSION = 7,
    USER_AGENT_BUILD_VERSION = 8,
  }
}

export class Metadata extends jspb.Message {
  getFilterMetadataMap(): jspb.Map<string, google_protobuf_struct_pb.Struct>;
  clearFilterMetadataMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Metadata.AsObject;
  static toObject(includeInstance: boolean, msg: Metadata): Metadata.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Metadata, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Metadata;
  static deserializeBinaryFromReader(message: Metadata, reader: jspb.BinaryReader): Metadata;
}

export namespace Metadata {
  export type AsObject = {
    filterMetadataMap: Array<[string, google_protobuf_struct_pb.Struct.AsObject]>,
  }
}

export class RuntimeUInt32 extends jspb.Message {
  getDefaultValue(): number;
  setDefaultValue(value: number): void;

  getRuntimeKey(): string;
  setRuntimeKey(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RuntimeUInt32.AsObject;
  static toObject(includeInstance: boolean, msg: RuntimeUInt32): RuntimeUInt32.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RuntimeUInt32, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RuntimeUInt32;
  static deserializeBinaryFromReader(message: RuntimeUInt32, reader: jspb.BinaryReader): RuntimeUInt32;
}

export namespace RuntimeUInt32 {
  export type AsObject = {
    defaultValue: number,
    runtimeKey: string,
  }
}

export class RuntimeFeatureFlag extends jspb.Message {
  hasDefaultValue(): boolean;
  clearDefaultValue(): void;
  getDefaultValue(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setDefaultValue(value?: google_protobuf_wrappers_pb.BoolValue): void;

  getRuntimeKey(): string;
  setRuntimeKey(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RuntimeFeatureFlag.AsObject;
  static toObject(includeInstance: boolean, msg: RuntimeFeatureFlag): RuntimeFeatureFlag.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RuntimeFeatureFlag, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RuntimeFeatureFlag;
  static deserializeBinaryFromReader(message: RuntimeFeatureFlag, reader: jspb.BinaryReader): RuntimeFeatureFlag;
}

export namespace RuntimeFeatureFlag {
  export type AsObject = {
    defaultValue?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    runtimeKey: string,
  }
}

export class HeaderValue extends jspb.Message {
  getKey(): string;
  setKey(value: string): void;

  getValue(): string;
  setValue(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HeaderValue.AsObject;
  static toObject(includeInstance: boolean, msg: HeaderValue): HeaderValue.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: HeaderValue, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HeaderValue;
  static deserializeBinaryFromReader(message: HeaderValue, reader: jspb.BinaryReader): HeaderValue;
}

export namespace HeaderValue {
  export type AsObject = {
    key: string,
    value: string,
  }
}

export class HeaderValueOption extends jspb.Message {
  hasHeader(): boolean;
  clearHeader(): void;
  getHeader(): HeaderValue | undefined;
  setHeader(value?: HeaderValue): void;

  hasHeaderSecretRef(): boolean;
  clearHeaderSecretRef(): void;
  getHeaderSecretRef(): github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef | undefined;
  setHeaderSecretRef(value?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef): void;

  hasAppend(): boolean;
  clearAppend(): void;
  getAppend(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setAppend(value?: google_protobuf_wrappers_pb.BoolValue): void;

  getHeaderOptionCase(): HeaderValueOption.HeaderOptionCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HeaderValueOption.AsObject;
  static toObject(includeInstance: boolean, msg: HeaderValueOption): HeaderValueOption.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: HeaderValueOption, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HeaderValueOption;
  static deserializeBinaryFromReader(message: HeaderValueOption, reader: jspb.BinaryReader): HeaderValueOption;
}

export namespace HeaderValueOption {
  export type AsObject = {
    header?: HeaderValue.AsObject,
    headerSecretRef?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject,
    append?: google_protobuf_wrappers_pb.BoolValue.AsObject,
  }

  export enum HeaderOptionCase {
    HEADER_OPTION_NOT_SET = 0,
    HEADER = 1,
    HEADER_SECRET_REF = 3,
  }
}

export class HeaderMap extends jspb.Message {
  clearHeadersList(): void;
  getHeadersList(): Array<HeaderValue>;
  setHeadersList(value: Array<HeaderValue>): void;
  addHeaders(value?: HeaderValue, index?: number): HeaderValue;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HeaderMap.AsObject;
  static toObject(includeInstance: boolean, msg: HeaderMap): HeaderMap.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: HeaderMap, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HeaderMap;
  static deserializeBinaryFromReader(message: HeaderMap, reader: jspb.BinaryReader): HeaderMap;
}

export namespace HeaderMap {
  export type AsObject = {
    headersList: Array<HeaderValue.AsObject>,
  }
}

export class DataSource extends jspb.Message {
  hasFilename(): boolean;
  clearFilename(): void;
  getFilename(): string;
  setFilename(value: string): void;

  hasInlineBytes(): boolean;
  clearInlineBytes(): void;
  getInlineBytes(): Uint8Array | string;
  getInlineBytes_asU8(): Uint8Array;
  getInlineBytes_asB64(): string;
  setInlineBytes(value: Uint8Array | string): void;

  hasInlineString(): boolean;
  clearInlineString(): void;
  getInlineString(): string;
  setInlineString(value: string): void;

  getSpecifierCase(): DataSource.SpecifierCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DataSource.AsObject;
  static toObject(includeInstance: boolean, msg: DataSource): DataSource.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DataSource, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DataSource;
  static deserializeBinaryFromReader(message: DataSource, reader: jspb.BinaryReader): DataSource;
}

export namespace DataSource {
  export type AsObject = {
    filename: string,
    inlineBytes: Uint8Array | string,
    inlineString: string,
  }

  export enum SpecifierCase {
    SPECIFIER_NOT_SET = 0,
    FILENAME = 1,
    INLINE_BYTES = 2,
    INLINE_STRING = 3,
  }
}

export class RemoteDataSource extends jspb.Message {
  hasHttpUri(): boolean;
  clearHttpUri(): void;
  getHttpUri(): github_com_solo_io_solo_kit_api_external_envoy_api_v2_core_http_uri_pb.HttpUri | undefined;
  setHttpUri(value?: github_com_solo_io_solo_kit_api_external_envoy_api_v2_core_http_uri_pb.HttpUri): void;

  getSha256(): string;
  setSha256(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoteDataSource.AsObject;
  static toObject(includeInstance: boolean, msg: RemoteDataSource): RemoteDataSource.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RemoteDataSource, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoteDataSource;
  static deserializeBinaryFromReader(message: RemoteDataSource, reader: jspb.BinaryReader): RemoteDataSource;
}

export namespace RemoteDataSource {
  export type AsObject = {
    httpUri?: github_com_solo_io_solo_kit_api_external_envoy_api_v2_core_http_uri_pb.HttpUri.AsObject,
    sha256: string,
  }
}

export class AsyncDataSource extends jspb.Message {
  hasLocal(): boolean;
  clearLocal(): void;
  getLocal(): DataSource | undefined;
  setLocal(value?: DataSource): void;

  hasRemote(): boolean;
  clearRemote(): void;
  getRemote(): RemoteDataSource | undefined;
  setRemote(value?: RemoteDataSource): void;

  getSpecifierCase(): AsyncDataSource.SpecifierCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AsyncDataSource.AsObject;
  static toObject(includeInstance: boolean, msg: AsyncDataSource): AsyncDataSource.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: AsyncDataSource, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AsyncDataSource;
  static deserializeBinaryFromReader(message: AsyncDataSource, reader: jspb.BinaryReader): AsyncDataSource;
}

export namespace AsyncDataSource {
  export type AsObject = {
    local?: DataSource.AsObject,
    remote?: RemoteDataSource.AsObject,
  }

  export enum SpecifierCase {
    SPECIFIER_NOT_SET = 0,
    LOCAL = 1,
    REMOTE = 2,
  }
}

export class TransportSocket extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  hasConfig(): boolean;
  clearConfig(): void;
  getConfig(): google_protobuf_struct_pb.Struct | undefined;
  setConfig(value?: google_protobuf_struct_pb.Struct): void;

  hasTypedConfig(): boolean;
  clearTypedConfig(): void;
  getTypedConfig(): google_protobuf_any_pb.Any | undefined;
  setTypedConfig(value?: google_protobuf_any_pb.Any): void;

  getConfigTypeCase(): TransportSocket.ConfigTypeCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TransportSocket.AsObject;
  static toObject(includeInstance: boolean, msg: TransportSocket): TransportSocket.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: TransportSocket, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TransportSocket;
  static deserializeBinaryFromReader(message: TransportSocket, reader: jspb.BinaryReader): TransportSocket;
}

export namespace TransportSocket {
  export type AsObject = {
    name: string,
    config?: google_protobuf_struct_pb.Struct.AsObject,
    typedConfig?: google_protobuf_any_pb.Any.AsObject,
  }

  export enum ConfigTypeCase {
    CONFIG_TYPE_NOT_SET = 0,
    CONFIG = 2,
    TYPED_CONFIG = 3,
  }
}

export class RuntimeFractionalPercent extends jspb.Message {
  hasDefaultValue(): boolean;
  clearDefaultValue(): void;
  getDefaultValue(): github_com_solo_io_solo_kit_api_external_envoy_type_percent_pb.FractionalPercent | undefined;
  setDefaultValue(value?: github_com_solo_io_solo_kit_api_external_envoy_type_percent_pb.FractionalPercent): void;

  getRuntimeKey(): string;
  setRuntimeKey(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RuntimeFractionalPercent.AsObject;
  static toObject(includeInstance: boolean, msg: RuntimeFractionalPercent): RuntimeFractionalPercent.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RuntimeFractionalPercent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RuntimeFractionalPercent;
  static deserializeBinaryFromReader(message: RuntimeFractionalPercent, reader: jspb.BinaryReader): RuntimeFractionalPercent;
}

export namespace RuntimeFractionalPercent {
  export type AsObject = {
    defaultValue?: github_com_solo_io_solo_kit_api_external_envoy_type_percent_pb.FractionalPercent.AsObject,
    runtimeKey: string,
  }
}

export class ControlPlane extends jspb.Message {
  getIdentifier(): string;
  setIdentifier(value: string): void;

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
    identifier: string,
  }
}

export interface RoutingPriorityMap {
  DEFAULT: 0;
  HIGH: 1;
}

export const RoutingPriority: RoutingPriorityMap;

export interface RequestMethodMap {
  METHOD_UNSPECIFIED: 0;
  GET: 1;
  HEAD: 2;
  POST: 3;
  PUT: 4;
  DELETE: 5;
  CONNECT: 6;
  OPTIONS: 7;
  TRACE: 8;
  PATCH: 9;
}

export const RequestMethod: RequestMethodMap;

export interface TrafficDirectionMap {
  UNSPECIFIED: 0;
  INBOUND: 1;
  OUTBOUND: 2;
}

export const TrafficDirection: TrafficDirectionMap;
