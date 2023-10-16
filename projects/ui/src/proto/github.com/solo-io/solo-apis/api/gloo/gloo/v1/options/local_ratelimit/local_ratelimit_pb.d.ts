/* eslint-disable */
// package: local_ratelimit.options.gloo.solo.io
// file: github.com/solo-io/solo-apis/api/gloo/gloo/v1/options/local_ratelimit/local_ratelimit.proto

import * as jspb from "google-protobuf";
import * as extproto_ext_pb from "../../../../../../../../../extproto/ext_pb";
import * as github_com_solo_io_solo_kit_api_external_envoy_api_v2_core_base_pb from "../../../../../../../../../github.com/solo-io/solo-kit/api/external/envoy/api/v2/core/base_pb";
import * as google_protobuf_duration_pb from "google-protobuf/google/protobuf/duration_pb";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";

export class TokenBucket extends jspb.Message {
  getMaxTokens(): number;
  setMaxTokens(value: number): void;

  hasTokensPerFill(): boolean;
  clearTokensPerFill(): void;
  getTokensPerFill(): google_protobuf_wrappers_pb.UInt32Value | undefined;
  setTokensPerFill(value?: google_protobuf_wrappers_pb.UInt32Value): void;

  hasFillInterval(): boolean;
  clearFillInterval(): void;
  getFillInterval(): google_protobuf_duration_pb.Duration | undefined;
  setFillInterval(value?: google_protobuf_duration_pb.Duration): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TokenBucket.AsObject;
  static toObject(includeInstance: boolean, msg: TokenBucket): TokenBucket.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: TokenBucket, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TokenBucket;
  static deserializeBinaryFromReader(message: TokenBucket, reader: jspb.BinaryReader): TokenBucket;
}

export namespace TokenBucket {
  export type AsObject = {
    maxTokens: number,
    tokensPerFill?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
    fillInterval?: google_protobuf_duration_pb.Duration.AsObject,
  }
}

export class Settings extends jspb.Message {
  hasDefaultLimit(): boolean;
  clearDefaultLimit(): void;
  getDefaultLimit(): TokenBucket | undefined;
  setDefaultLimit(value?: TokenBucket): void;

  hasLocalRateLimitPerDownstreamConnection(): boolean;
  clearLocalRateLimitPerDownstreamConnection(): void;
  getLocalRateLimitPerDownstreamConnection(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setLocalRateLimitPerDownstreamConnection(value?: google_protobuf_wrappers_pb.BoolValue): void;

  hasEnableXRatelimitHeaders(): boolean;
  clearEnableXRatelimitHeaders(): void;
  getEnableXRatelimitHeaders(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setEnableXRatelimitHeaders(value?: google_protobuf_wrappers_pb.BoolValue): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Settings.AsObject;
  static toObject(includeInstance: boolean, msg: Settings): Settings.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Settings, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Settings;
  static deserializeBinaryFromReader(message: Settings, reader: jspb.BinaryReader): Settings;
}

export namespace Settings {
  export type AsObject = {
    defaultLimit?: TokenBucket.AsObject,
    localRateLimitPerDownstreamConnection?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    enableXRatelimitHeaders?: google_protobuf_wrappers_pb.BoolValue.AsObject,
  }
}
