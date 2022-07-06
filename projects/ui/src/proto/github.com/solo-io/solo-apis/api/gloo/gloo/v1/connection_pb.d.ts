/* eslint-disable */
// package: gloo.solo.io
// file: github.com/solo-io/solo-apis/api/gloo/gloo/v1/connection.proto

import * as jspb from "google-protobuf";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_protocol_protocol_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/options/protocol/protocol_pb";
import * as google_protobuf_duration_pb from "google-protobuf/google/protobuf/duration_pb";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";
import * as extproto_ext_pb from "../../../../../../../extproto/ext_pb";

export class ConnectionConfig extends jspb.Message {
  getMaxRequestsPerConnection(): number;
  setMaxRequestsPerConnection(value: number): void;

  hasConnectTimeout(): boolean;
  clearConnectTimeout(): void;
  getConnectTimeout(): google_protobuf_duration_pb.Duration | undefined;
  setConnectTimeout(value?: google_protobuf_duration_pb.Duration): void;

  hasTcpKeepalive(): boolean;
  clearTcpKeepalive(): void;
  getTcpKeepalive(): ConnectionConfig.TcpKeepAlive | undefined;
  setTcpKeepalive(value?: ConnectionConfig.TcpKeepAlive): void;

  hasPerConnectionBufferLimitBytes(): boolean;
  clearPerConnectionBufferLimitBytes(): void;
  getPerConnectionBufferLimitBytes(): google_protobuf_wrappers_pb.UInt32Value | undefined;
  setPerConnectionBufferLimitBytes(value?: google_protobuf_wrappers_pb.UInt32Value): void;

  hasCommonHttpProtocolOptions(): boolean;
  clearCommonHttpProtocolOptions(): void;
  getCommonHttpProtocolOptions(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_protocol_protocol_pb.HttpProtocolOptions | undefined;
  setCommonHttpProtocolOptions(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_protocol_protocol_pb.HttpProtocolOptions): void;

  hasHttp1ProtocolOptions(): boolean;
  clearHttp1ProtocolOptions(): void;
  getHttp1ProtocolOptions(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_protocol_protocol_pb.Http1ProtocolOptions | undefined;
  setHttp1ProtocolOptions(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_protocol_protocol_pb.Http1ProtocolOptions): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ConnectionConfig.AsObject;
  static toObject(includeInstance: boolean, msg: ConnectionConfig): ConnectionConfig.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ConnectionConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ConnectionConfig;
  static deserializeBinaryFromReader(message: ConnectionConfig, reader: jspb.BinaryReader): ConnectionConfig;
}

export namespace ConnectionConfig {
  export type AsObject = {
    maxRequestsPerConnection: number,
    connectTimeout?: google_protobuf_duration_pb.Duration.AsObject,
    tcpKeepalive?: ConnectionConfig.TcpKeepAlive.AsObject,
    perConnectionBufferLimitBytes?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
    commonHttpProtocolOptions?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_protocol_protocol_pb.HttpProtocolOptions.AsObject,
    http1ProtocolOptions?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_protocol_protocol_pb.Http1ProtocolOptions.AsObject,
  }

  export class TcpKeepAlive extends jspb.Message {
    getKeepaliveProbes(): number;
    setKeepaliveProbes(value: number): void;

    hasKeepaliveTime(): boolean;
    clearKeepaliveTime(): void;
    getKeepaliveTime(): google_protobuf_duration_pb.Duration | undefined;
    setKeepaliveTime(value?: google_protobuf_duration_pb.Duration): void;

    hasKeepaliveInterval(): boolean;
    clearKeepaliveInterval(): void;
    getKeepaliveInterval(): google_protobuf_duration_pb.Duration | undefined;
    setKeepaliveInterval(value?: google_protobuf_duration_pb.Duration): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): TcpKeepAlive.AsObject;
    static toObject(includeInstance: boolean, msg: TcpKeepAlive): TcpKeepAlive.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: TcpKeepAlive, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): TcpKeepAlive;
    static deserializeBinaryFromReader(message: TcpKeepAlive, reader: jspb.BinaryReader): TcpKeepAlive;
  }

  export namespace TcpKeepAlive {
    export type AsObject = {
      keepaliveProbes: number,
      keepaliveTime?: google_protobuf_duration_pb.Duration.AsObject,
      keepaliveInterval?: google_protobuf_duration_pb.Duration.AsObject,
    }
  }
}
