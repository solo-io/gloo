// package: glooe.solo.io
// file: github.com/solo-io/solo-projects/projects/gloo/api/v1/ratelimit.proto

import * as jspb from "google-protobuf";
import * as envoy_api_v2_discovery_pb from "../../../../../../../envoy/api/v2/discovery_pb";
import * as google_api_annotations_pb from "../../../../../../../google/api/annotations_pb";
import * as github_com_solo_io_solo_projects_projects_gloo_api_v1_plugins_ratelimit_ratelimit_pb from "../../../../../../../github.com/solo-io/solo-projects/projects/gloo/api/v1/plugins/ratelimit/ratelimit_pb";
import * as gogoproto_gogo_pb from "../../../../../../../gogoproto/gogo_pb";

export class RateLimitConfig extends jspb.Message {
  getDomain(): string;
  setDomain(value: string): void;

  clearDescriptorsList(): void;
  getDescriptorsList(): Array<github_com_solo_io_solo_projects_projects_gloo_api_v1_plugins_ratelimit_ratelimit_pb.Descriptor>;
  setDescriptorsList(value: Array<github_com_solo_io_solo_projects_projects_gloo_api_v1_plugins_ratelimit_ratelimit_pb.Descriptor>): void;
  addDescriptors(value?: github_com_solo_io_solo_projects_projects_gloo_api_v1_plugins_ratelimit_ratelimit_pb.Descriptor, index?: number): github_com_solo_io_solo_projects_projects_gloo_api_v1_plugins_ratelimit_ratelimit_pb.Descriptor;

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
    descriptorsList: Array<github_com_solo_io_solo_projects_projects_gloo_api_v1_plugins_ratelimit_ratelimit_pb.Descriptor.AsObject>,
  }
}

