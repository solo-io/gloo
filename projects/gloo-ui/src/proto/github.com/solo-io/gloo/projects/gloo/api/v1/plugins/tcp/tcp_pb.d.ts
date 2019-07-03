// package: tcp.plugins.gloo.solo.io
// file: github.com/solo-io/gloo/projects/gloo/api/v1/plugins/tcp/tcp.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../../../../../../../gogoproto/gogo_pb";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";
import * as google_protobuf_duration_pb from "google-protobuf/google/protobuf/duration_pb";

export class TcpProxySettings extends jspb.Message {
  hasMaxConnectAttempts(): boolean;
  clearMaxConnectAttempts(): void;
  getMaxConnectAttempts(): google_protobuf_wrappers_pb.UInt32Value | undefined;
  setMaxConnectAttempts(value?: google_protobuf_wrappers_pb.UInt32Value): void;

  hasIdleTimeout(): boolean;
  clearIdleTimeout(): void;
  getIdleTimeout(): google_protobuf_duration_pb.Duration | undefined;
  setIdleTimeout(value?: google_protobuf_duration_pb.Duration): void;

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
  }
}

