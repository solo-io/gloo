/* eslint-disable */
// package: glooe.solo.io
// file: github.com/solo-io/gloo/projects/gloo/api/v1/enterprise/ratelimit.proto

import * as jspb from "google-protobuf";
import * as github_com_solo_io_solo_kit_api_external_envoy_api_v2_discovery_pb from "../../../../../../../../github.com/solo-io/solo-kit/api/external/envoy/api/v2/discovery_pb";
import * as google_api_annotations_pb from "../../../../../../../../google/api/annotations_pb";
import * as github_com_solo_io_solo_apis_api_rate_limiter_v1alpha1_ratelimit_pb from "../../../../../../../../github.com/solo-io/solo-apis/api/rate-limiter/v1alpha1/ratelimit_pb";
import * as extproto_ext_pb from "../../../../../../../../protoc-gen-ext/extproto/ext_pb";

export class RateLimitConfig extends jspb.Message {
  getDomain(): string;
  setDomain(value: string): void;

  clearDescriptorsList(): void;
  getDescriptorsList(): Array<github_com_solo_io_solo_apis_api_rate_limiter_v1alpha1_ratelimit_pb.Descriptor>;
  setDescriptorsList(value: Array<github_com_solo_io_solo_apis_api_rate_limiter_v1alpha1_ratelimit_pb.Descriptor>): void;
  addDescriptors(value?: github_com_solo_io_solo_apis_api_rate_limiter_v1alpha1_ratelimit_pb.Descriptor, index?: number): github_com_solo_io_solo_apis_api_rate_limiter_v1alpha1_ratelimit_pb.Descriptor;

  clearSetDescriptorsList(): void;
  getSetDescriptorsList(): Array<github_com_solo_io_solo_apis_api_rate_limiter_v1alpha1_ratelimit_pb.SetDescriptor>;
  setSetDescriptorsList(value: Array<github_com_solo_io_solo_apis_api_rate_limiter_v1alpha1_ratelimit_pb.SetDescriptor>): void;
  addSetDescriptors(value?: github_com_solo_io_solo_apis_api_rate_limiter_v1alpha1_ratelimit_pb.SetDescriptor, index?: number): github_com_solo_io_solo_apis_api_rate_limiter_v1alpha1_ratelimit_pb.SetDescriptor;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RateLimitConfig.AsObject;
  static toObject(includeInstance: boolean, msg: RateLimitConfig): RateLimitConfig.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RateLimitConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RateLimitConfig;
  static deserializeBinaryFromReader(message: RateLimitConfig, reader: jspb.BinaryReader): RateLimitConfig;
}

export namespace RateLimitConfig {
  export type AsObject = {
    domain: string,
    descriptorsList: Array<github_com_solo_io_solo_apis_api_rate_limiter_v1alpha1_ratelimit_pb.Descriptor.AsObject>,
    setDescriptorsList: Array<github_com_solo_io_solo_apis_api_rate_limiter_v1alpha1_ratelimit_pb.SetDescriptor.AsObject>,
  }
}
