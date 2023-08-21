/* eslint-disable */
// package: connection_limit.options.gloo.solo.io
// file: github.com/solo-io/solo-apis/api/gloo/gloo/v1/options/connection_limit/connection_limit.proto

import * as jspb from "google-protobuf";
import * as extproto_ext_pb from "../../../../../../../../../extproto/ext_pb";
import * as github_com_solo_io_solo_kit_api_external_envoy_api_v2_core_base_pb from "../../../../../../../../../github.com/solo-io/solo-kit/api/external/envoy/api/v2/core/base_pb";
import * as google_protobuf_duration_pb from "google-protobuf/google/protobuf/duration_pb";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";

export class ConnectionLimit extends jspb.Message {
  hasMaxActiveConnections(): boolean;
  clearMaxActiveConnections(): void;
  getMaxActiveConnections(): google_protobuf_wrappers_pb.UInt32Value | undefined;
  setMaxActiveConnections(value?: google_protobuf_wrappers_pb.UInt32Value): void;

  hasDelayBeforeClose(): boolean;
  clearDelayBeforeClose(): void;
  getDelayBeforeClose(): google_protobuf_duration_pb.Duration | undefined;
  setDelayBeforeClose(value?: google_protobuf_duration_pb.Duration): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ConnectionLimit.AsObject;
  static toObject(includeInstance: boolean, msg: ConnectionLimit): ConnectionLimit.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ConnectionLimit, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ConnectionLimit;
  static deserializeBinaryFromReader(message: ConnectionLimit, reader: jspb.BinaryReader): ConnectionLimit;
}

export namespace ConnectionLimit {
  export type AsObject = {
    maxActiveConnections?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
    delayBeforeClose?: google_protobuf_duration_pb.Duration.AsObject,
  }
}
