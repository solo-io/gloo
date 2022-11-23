/* eslint-disable */
// package: retries.options.gloo.solo.io
// file: github.com/solo-io/solo-apis/api/gloo/gloo/v1/options/retries/retries.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_duration_pb from "google-protobuf/google/protobuf/duration_pb";
import * as validate_validate_pb from "../../../../../../../../../validate/validate_pb";
import * as extproto_ext_pb from "../../../../../../../../../extproto/ext_pb";

export class RetryBackOff extends jspb.Message {
  hasBaseInterval(): boolean;
  clearBaseInterval(): void;
  getBaseInterval(): google_protobuf_duration_pb.Duration | undefined;
  setBaseInterval(value?: google_protobuf_duration_pb.Duration): void;

  hasMaxInterval(): boolean;
  clearMaxInterval(): void;
  getMaxInterval(): google_protobuf_duration_pb.Duration | undefined;
  setMaxInterval(value?: google_protobuf_duration_pb.Duration): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RetryBackOff.AsObject;
  static toObject(includeInstance: boolean, msg: RetryBackOff): RetryBackOff.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RetryBackOff, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RetryBackOff;
  static deserializeBinaryFromReader(message: RetryBackOff, reader: jspb.BinaryReader): RetryBackOff;
}

export namespace RetryBackOff {
  export type AsObject = {
    baseInterval?: google_protobuf_duration_pb.Duration.AsObject,
    maxInterval?: google_protobuf_duration_pb.Duration.AsObject,
  }
}

export class RetryPolicy extends jspb.Message {
  getRetryOn(): string;
  setRetryOn(value: string): void;

  getNumRetries(): number;
  setNumRetries(value: number): void;

  hasPerTryTimeout(): boolean;
  clearPerTryTimeout(): void;
  getPerTryTimeout(): google_protobuf_duration_pb.Duration | undefined;
  setPerTryTimeout(value?: google_protobuf_duration_pb.Duration): void;

  hasRetryBackOff(): boolean;
  clearRetryBackOff(): void;
  getRetryBackOff(): RetryBackOff | undefined;
  setRetryBackOff(value?: RetryBackOff): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RetryPolicy.AsObject;
  static toObject(includeInstance: boolean, msg: RetryPolicy): RetryPolicy.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RetryPolicy, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RetryPolicy;
  static deserializeBinaryFromReader(message: RetryPolicy, reader: jspb.BinaryReader): RetryPolicy;
}

export namespace RetryPolicy {
  export type AsObject = {
    retryOn: string,
    numRetries: number,
    perTryTimeout?: google_protobuf_duration_pb.Duration.AsObject,
    retryBackOff?: RetryBackOff.AsObject,
  }
}
