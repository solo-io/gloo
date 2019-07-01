// package: envoy.api.v2.core
// file: github.com/solo-io/solo-kit/api/external/envoy/api/v2/core/base.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_any_pb from "google-protobuf/google/protobuf/any_pb";
import * as google_protobuf_struct_pb from "google-protobuf/google/protobuf/struct_pb";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";
import * as validate_validate_pb from "../../../../../../../../../validate/validate_pb";
import * as gogoproto_gogo_pb from "../../../../../../../../../gogoproto/gogo_pb";
import * as envoy_type_percent_pb from "../../../../../../../../../envoy/type/percent_pb";

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

  hasAppend(): boolean;
  clearAppend(): void;
  getAppend(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setAppend(value?: google_protobuf_wrappers_pb.BoolValue): void;

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
    append?: google_protobuf_wrappers_pb.BoolValue.AsObject,
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

export class SocketOption extends jspb.Message {
  getDescription(): string;
  setDescription(value: string): void;

  getLevel(): number;
  setLevel(value: number): void;

  getName(): number;
  setName(value: number): void;

  hasIntValue(): boolean;
  clearIntValue(): void;
  getIntValue(): number;
  setIntValue(value: number): void;

  hasBufValue(): boolean;
  clearBufValue(): void;
  getBufValue(): Uint8Array | string;
  getBufValue_asU8(): Uint8Array;
  getBufValue_asB64(): string;
  setBufValue(value: Uint8Array | string): void;

  getState(): SocketOption.SocketState;
  setState(value: SocketOption.SocketState): void;

  getValueCase(): SocketOption.ValueCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SocketOption.AsObject;
  static toObject(includeInstance: boolean, msg: SocketOption): SocketOption.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: SocketOption, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SocketOption;
  static deserializeBinaryFromReader(message: SocketOption, reader: jspb.BinaryReader): SocketOption;
}

export namespace SocketOption {
  export type AsObject = {
    description: string,
    level: number,
    name: number,
    intValue: number,
    bufValue: Uint8Array | string,
    state: SocketOption.SocketState,
  }

  export enum SocketState {
    STATE_PREBIND = 0,
    STATE_BOUND = 1,
    STATE_LISTENING = 2,
  }

  export enum ValueCase {
    VALUE_NOT_SET = 0,
    INT_VALUE = 4,
    BUF_VALUE = 5,
  }
}

export class RuntimeFractionalPercent extends jspb.Message {
  hasDefaultValue(): boolean;
  clearDefaultValue(): void;
  getDefaultValue(): envoy_type_percent_pb.FractionalPercent | undefined;
  setDefaultValue(value?: envoy_type_percent_pb.FractionalPercent): void;

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
    defaultValue?: envoy_type_percent_pb.FractionalPercent.AsObject,
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

export enum RoutingPriority {
  DEFAULT = 0,
  HIGH = 1,
}

export enum RequestMethod {
  METHOD_UNSPECIFIED = 0,
  GET = 1,
  HEAD = 2,
  POST = 3,
  PUT = 4,
  DELETE = 5,
  CONNECT = 6,
  OPTIONS = 7,
  TRACE = 8,
}

