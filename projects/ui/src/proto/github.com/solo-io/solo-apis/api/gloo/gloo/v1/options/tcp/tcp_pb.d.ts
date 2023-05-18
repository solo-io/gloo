/* eslint-disable */
// package: tcp.options.gloo.solo.io
// file: github.com/solo-io/solo-apis/api/gloo/gloo/v1/options/tcp/tcp.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";
import * as google_protobuf_duration_pb from "google-protobuf/google/protobuf/duration_pb";
import * as github_com_solo_io_solo_kit_api_external_envoy_api_v2_core_base_pb from "../../../../../../../../../github.com/solo-io/solo-kit/api/external/envoy/api/v2/core/base_pb";
import * as extproto_ext_pb from "../../../../../../../../../extproto/ext_pb";

export class TcpProxySettings extends jspb.Message {
  hasMaxConnectAttempts(): boolean;
  clearMaxConnectAttempts(): void;
  getMaxConnectAttempts(): google_protobuf_wrappers_pb.UInt32Value | undefined;
  setMaxConnectAttempts(value?: google_protobuf_wrappers_pb.UInt32Value): void;

  hasIdleTimeout(): boolean;
  clearIdleTimeout(): void;
  getIdleTimeout(): google_protobuf_duration_pb.Duration | undefined;
  setIdleTimeout(value?: google_protobuf_duration_pb.Duration): void;

  hasTunnelingConfig(): boolean;
  clearTunnelingConfig(): void;
  getTunnelingConfig(): TcpProxySettings.TunnelingConfig | undefined;
  setTunnelingConfig(value?: TcpProxySettings.TunnelingConfig): void;

  hasAccessLogFlushInterval(): boolean;
  clearAccessLogFlushInterval(): void;
  getAccessLogFlushInterval(): google_protobuf_duration_pb.Duration | undefined;
  setAccessLogFlushInterval(value?: google_protobuf_duration_pb.Duration): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TcpProxySettings.AsObject;
  static toObject(includeInstance: boolean, msg: TcpProxySettings): TcpProxySettings.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: TcpProxySettings, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TcpProxySettings;
  static deserializeBinaryFromReader(message: TcpProxySettings, reader: jspb.BinaryReader): TcpProxySettings;
}

export namespace TcpProxySettings {
  export type AsObject = {
    maxConnectAttempts?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
    idleTimeout?: google_protobuf_duration_pb.Duration.AsObject,
    tunnelingConfig?: TcpProxySettings.TunnelingConfig.AsObject,
    accessLogFlushInterval?: google_protobuf_duration_pb.Duration.AsObject,
  }

  export class TunnelingConfig extends jspb.Message {
    getHostname(): string;
    setHostname(value: string): void;

    clearHeadersToAddList(): void;
    getHeadersToAddList(): Array<HeaderValueOption>;
    setHeadersToAddList(value: Array<HeaderValueOption>): void;
    addHeadersToAdd(value?: HeaderValueOption, index?: number): HeaderValueOption;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): TunnelingConfig.AsObject;
    static toObject(includeInstance: boolean, msg: TunnelingConfig): TunnelingConfig.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: TunnelingConfig, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): TunnelingConfig;
    static deserializeBinaryFromReader(message: TunnelingConfig, reader: jspb.BinaryReader): TunnelingConfig;
  }

  export namespace TunnelingConfig {
    export type AsObject = {
      hostname: string,
      headersToAddList: Array<HeaderValueOption.AsObject>,
    }
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
