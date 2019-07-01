// package: retries.plugins.gloo.solo.io
// file: github.com/solo-io/gloo/projects/gloo/api/v1/plugins/retries/retries.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_duration_pb from "google-protobuf/google/protobuf/duration_pb";
import * as gogoproto_gogo_pb from "../../../../../../../../../gogoproto/gogo_pb";

export class RetryPolicy extends jspb.Message {
  getRetryOn(): string;
  setRetryOn(value: string): void;

  getNumRetries(): number;
  setNumRetries(value: number): void;

  hasPerTryTimeout(): boolean;
  clearPerTryTimeout(): void;
  getPerTryTimeout(): google_protobuf_duration_pb.Duration | undefined;
  setPerTryTimeout(value?: google_protobuf_duration_pb.Duration): void;

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
  }
}

