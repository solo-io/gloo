// package: gloo.solo.io
// file: github.com/solo-io/gloo/projects/gloo/api/v1/circuit_breaker.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../../../../../gogoproto/gogo_pb";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";

export class CircuitBreakerConfig extends jspb.Message {
  hasMaxConnections(): boolean;
  clearMaxConnections(): void;
  getMaxConnections(): google_protobuf_wrappers_pb.UInt32Value | undefined;
  setMaxConnections(value?: google_protobuf_wrappers_pb.UInt32Value): void;

  hasMaxPendingRequests(): boolean;
  clearMaxPendingRequests(): void;
  getMaxPendingRequests(): google_protobuf_wrappers_pb.UInt32Value | undefined;
  setMaxPendingRequests(value?: google_protobuf_wrappers_pb.UInt32Value): void;

  hasMaxRequests(): boolean;
  clearMaxRequests(): void;
  getMaxRequests(): google_protobuf_wrappers_pb.UInt32Value | undefined;
  setMaxRequests(value?: google_protobuf_wrappers_pb.UInt32Value): void;

  hasMaxRetries(): boolean;
  clearMaxRetries(): void;
  getMaxRetries(): google_protobuf_wrappers_pb.UInt32Value | undefined;
  setMaxRetries(value?: google_protobuf_wrappers_pb.UInt32Value): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CircuitBreakerConfig.AsObject;
  static toObject(includeInstance: boolean, msg: CircuitBreakerConfig): CircuitBreakerConfig.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CircuitBreakerConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CircuitBreakerConfig;
  static deserializeBinaryFromReader(message: CircuitBreakerConfig, reader: jspb.BinaryReader): CircuitBreakerConfig;
}

export namespace CircuitBreakerConfig {
  export type AsObject = {
    maxConnections?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
    maxPendingRequests?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
    maxRequests?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
    maxRetries?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
  }
}

