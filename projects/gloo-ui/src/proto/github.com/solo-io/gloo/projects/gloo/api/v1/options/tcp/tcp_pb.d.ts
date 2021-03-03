/* eslint-disable */
// package: tcp.options.gloo.solo.io
// file: github.com/solo-io/gloo/projects/gloo/api/v1/options/tcp/tcp.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";
import * as google_protobuf_duration_pb from "google-protobuf/google/protobuf/duration_pb";
import * as github_com_solo_io_solo_kit_api_external_envoy_api_v2_core_base_pb from "../../../../../../../../../github.com/solo-io/solo-kit/api/external/envoy/api/v2/core/base_pb";
import * as extproto_ext_pb from "../../../../../../../../../protoc-gen-ext/extproto/ext_pb";

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
  }

  export class TunnelingConfig extends jspb.Message {
    getHostname(): string;
    setHostname(value: string): void;

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
    }
  }
}
