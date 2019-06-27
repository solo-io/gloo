// package: gloo.solo.io
// file: github.com/solo-io/gloo/projects/gloo/api/v1/connection.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../../../../../gogoproto/gogo_pb";
import * as google_protobuf_duration_pb from "google-protobuf/google/protobuf/duration_pb";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";

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

