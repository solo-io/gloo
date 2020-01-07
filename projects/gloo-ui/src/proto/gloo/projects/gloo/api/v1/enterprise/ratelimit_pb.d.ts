// package: glooe.solo.io
// file: gloo/projects/gloo/api/v1/enterprise/ratelimit.proto

import * as jspb from "google-protobuf";
import * as envoy_api_v2_discovery_pb from "../../../../../../envoy/api/v2/discovery_pb";
import * as google_api_annotations_pb from "../../../../../../google/api/annotations_pb";
import * as gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_pb from "../../../../../../gloo/projects/gloo/api/v1/enterprise/options/ratelimit/ratelimit_pb";
import * as gogoproto_gogo_pb from "../../../../../../gogoproto/gogo_pb";
import * as extproto_ext_pb from "../../../../../../extproto/ext_pb";

export class RateLimitConfig extends jspb.Message {
  getDomain(): string;
  setDomain(value: string): void;

  clearDescriptorsList(): void;
  getDescriptorsList(): Array<gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_pb.Descriptor>;
  setDescriptorsList(value: Array<gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_pb.Descriptor>): void;
  addDescriptors(value?: gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_pb.Descriptor, index?: number): gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_pb.Descriptor;

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
    descriptorsList: Array<gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_pb.Descriptor.AsObject>,
  }
}

