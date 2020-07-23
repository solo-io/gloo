/* eslint-disable */
// package: envoy.config.core.v3
// file: envoy/config/core/v3/backoff.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_duration_pb from "google-protobuf/google/protobuf/duration_pb";
import * as udpa_annotations_status_pb from "../../../../udpa/annotations/status_pb";
import * as udpa_annotations_versioning_pb from "../../../../udpa/annotations/versioning_pb";
import * as validate_validate_pb from "../../../../validate/validate_pb";
import * as gogoproto_gogo_pb from "../../../../gogoproto/gogo_pb";

export class BackoffStrategy extends jspb.Message {
  hasBaseInterval(): boolean;
  clearBaseInterval(): void;
  getBaseInterval(): google_protobuf_duration_pb.Duration | undefined;
  setBaseInterval(value?: google_protobuf_duration_pb.Duration): void;

  hasMaxInterval(): boolean;
  clearMaxInterval(): void;
  getMaxInterval(): google_protobuf_duration_pb.Duration | undefined;
  setMaxInterval(value?: google_protobuf_duration_pb.Duration): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): BackoffStrategy.AsObject;
  static toObject(includeInstance: boolean, msg: BackoffStrategy): BackoffStrategy.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: BackoffStrategy, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): BackoffStrategy;
  static deserializeBinaryFromReader(message: BackoffStrategy, reader: jspb.BinaryReader): BackoffStrategy;
}

export namespace BackoffStrategy {
  export type AsObject = {
    baseInterval?: google_protobuf_duration_pb.Duration.AsObject,
    maxInterval?: google_protobuf_duration_pb.Duration.AsObject,
  }
}
