/* eslint-disable */
// package: solo.io.envoy.api.v2.core
// file: github.com/solo-io/solo-kit/api/external/envoy/api/v2/core/address.proto

import * as jspb from "google-protobuf";
import * as github_com_solo_io_solo_kit_api_external_envoy_api_v2_core_socket_option_pb from "../../../../../../../../../github.com/solo-io/solo-kit/api/external/envoy/api/v2/core/socket_option_pb";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";
import * as validate_validate_pb from "../../../../../../../../../validate/validate_pb";
import * as extproto_ext_pb from "../../../../../../../../../extproto/ext_pb";

export class Pipe extends jspb.Message {
  getPath(): string;
  setPath(value: string): void;

  getMode(): number;
  setMode(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Pipe.AsObject;
  static toObject(includeInstance: boolean, msg: Pipe): Pipe.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Pipe, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Pipe;
  static deserializeBinaryFromReader(message: Pipe, reader: jspb.BinaryReader): Pipe;
}

export namespace Pipe {
  export type AsObject = {
    path: string,
    mode: number,
  }
}

export class SocketAddress extends jspb.Message {
  getProtocol(): SocketAddress.ProtocolMap[keyof SocketAddress.ProtocolMap];
  setProtocol(value: SocketAddress.ProtocolMap[keyof SocketAddress.ProtocolMap]): void;

  getAddress(): string;
  setAddress(value: string): void;

  hasPortValue(): boolean;
  clearPortValue(): void;
  getPortValue(): number;
  setPortValue(value: number): void;

  hasNamedPort(): boolean;
  clearNamedPort(): void;
  getNamedPort(): string;
  setNamedPort(value: string): void;

  getResolverName(): string;
  setResolverName(value: string): void;

  getIpv4Compat(): boolean;
  setIpv4Compat(value: boolean): void;

  getPortSpecifierCase(): SocketAddress.PortSpecifierCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SocketAddress.AsObject;
  static toObject(includeInstance: boolean, msg: SocketAddress): SocketAddress.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: SocketAddress, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SocketAddress;
  static deserializeBinaryFromReader(message: SocketAddress, reader: jspb.BinaryReader): SocketAddress;
}

export namespace SocketAddress {
  export type AsObject = {
    protocol: SocketAddress.ProtocolMap[keyof SocketAddress.ProtocolMap],
    address: string,
    portValue: number,
    namedPort: string,
    resolverName: string,
    ipv4Compat: boolean,
  }

  export interface ProtocolMap {
    TCP: 0;
    UDP: 1;
  }

  export const Protocol: ProtocolMap;

  export enum PortSpecifierCase {
    PORT_SPECIFIER_NOT_SET = 0,
    PORT_VALUE = 3,
    NAMED_PORT = 4,
  }
}

export class TcpKeepalive extends jspb.Message {
  hasKeepaliveProbes(): boolean;
  clearKeepaliveProbes(): void;
  getKeepaliveProbes(): google_protobuf_wrappers_pb.UInt32Value | undefined;
  setKeepaliveProbes(value?: google_protobuf_wrappers_pb.UInt32Value): void;

  hasKeepaliveTime(): boolean;
  clearKeepaliveTime(): void;
  getKeepaliveTime(): google_protobuf_wrappers_pb.UInt32Value | undefined;
  setKeepaliveTime(value?: google_protobuf_wrappers_pb.UInt32Value): void;

  hasKeepaliveInterval(): boolean;
  clearKeepaliveInterval(): void;
  getKeepaliveInterval(): google_protobuf_wrappers_pb.UInt32Value | undefined;
  setKeepaliveInterval(value?: google_protobuf_wrappers_pb.UInt32Value): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TcpKeepalive.AsObject;
  static toObject(includeInstance: boolean, msg: TcpKeepalive): TcpKeepalive.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: TcpKeepalive, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TcpKeepalive;
  static deserializeBinaryFromReader(message: TcpKeepalive, reader: jspb.BinaryReader): TcpKeepalive;
}

export namespace TcpKeepalive {
  export type AsObject = {
    keepaliveProbes?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
    keepaliveTime?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
    keepaliveInterval?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
  }
}

export class BindConfig extends jspb.Message {
  hasSourceAddress(): boolean;
  clearSourceAddress(): void;
  getSourceAddress(): SocketAddress | undefined;
  setSourceAddress(value?: SocketAddress): void;

  hasFreebind(): boolean;
  clearFreebind(): void;
  getFreebind(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setFreebind(value?: google_protobuf_wrappers_pb.BoolValue): void;

  clearSocketOptionsList(): void;
  getSocketOptionsList(): Array<github_com_solo_io_solo_kit_api_external_envoy_api_v2_core_socket_option_pb.SocketOption>;
  setSocketOptionsList(value: Array<github_com_solo_io_solo_kit_api_external_envoy_api_v2_core_socket_option_pb.SocketOption>): void;
  addSocketOptions(value?: github_com_solo_io_solo_kit_api_external_envoy_api_v2_core_socket_option_pb.SocketOption, index?: number): github_com_solo_io_solo_kit_api_external_envoy_api_v2_core_socket_option_pb.SocketOption;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): BindConfig.AsObject;
  static toObject(includeInstance: boolean, msg: BindConfig): BindConfig.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: BindConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): BindConfig;
  static deserializeBinaryFromReader(message: BindConfig, reader: jspb.BinaryReader): BindConfig;
}

export namespace BindConfig {
  export type AsObject = {
    sourceAddress?: SocketAddress.AsObject,
    freebind?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    socketOptionsList: Array<github_com_solo_io_solo_kit_api_external_envoy_api_v2_core_socket_option_pb.SocketOption.AsObject>,
  }
}

export class Address extends jspb.Message {
  hasSocketAddress(): boolean;
  clearSocketAddress(): void;
  getSocketAddress(): SocketAddress | undefined;
  setSocketAddress(value?: SocketAddress): void;

  hasPipe(): boolean;
  clearPipe(): void;
  getPipe(): Pipe | undefined;
  setPipe(value?: Pipe): void;

  getAddressCase(): Address.AddressCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Address.AsObject;
  static toObject(includeInstance: boolean, msg: Address): Address.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Address, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Address;
  static deserializeBinaryFromReader(message: Address, reader: jspb.BinaryReader): Address;
}

export namespace Address {
  export type AsObject = {
    socketAddress?: SocketAddress.AsObject,
    pipe?: Pipe.AsObject,
  }

  export enum AddressCase {
    ADDRESS_NOT_SET = 0,
    SOCKET_ADDRESS = 1,
    PIPE = 2,
  }
}

export class CidrRange extends jspb.Message {
  getAddressPrefix(): string;
  setAddressPrefix(value: string): void;

  hasPrefixLen(): boolean;
  clearPrefixLen(): void;
  getPrefixLen(): google_protobuf_wrappers_pb.UInt32Value | undefined;
  setPrefixLen(value?: google_protobuf_wrappers_pb.UInt32Value): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CidrRange.AsObject;
  static toObject(includeInstance: boolean, msg: CidrRange): CidrRange.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CidrRange, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CidrRange;
  static deserializeBinaryFromReader(message: CidrRange, reader: jspb.BinaryReader): CidrRange;
}

export namespace CidrRange {
  export type AsObject = {
    addressPrefix: string,
    prefixLen?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
  }
}
