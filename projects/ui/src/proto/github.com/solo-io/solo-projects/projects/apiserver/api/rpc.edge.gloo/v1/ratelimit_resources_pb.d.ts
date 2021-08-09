/* eslint-disable */
// package: rpc.edge.gloo.solo.io
// file: github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/ratelimit_resources.proto

import * as jspb from "google-protobuf";
import * as extproto_ext_pb from "../../../../../../../../extproto/ext_pb";
import * as github_com_solo_io_solo_apis_api_rate_limiter_v1alpha1_ratelimit_pb from "../../../../../../../../github.com/solo-io/solo-apis/api/rate-limiter/v1alpha1/ratelimit_pb";
import * as github_com_solo_io_skv2_api_core_v1_core_pb from "../../../../../../../../github.com/solo-io/skv2/api/core/v1/core_pb";
import * as github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb from "../../../../../../../../github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/common_pb";

export class RateLimitConfig extends jspb.Message {
  hasMetadata(): boolean;
  clearMetadata(): void;
  getMetadata(): github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ObjectMeta | undefined;
  setMetadata(value?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ObjectMeta): void;

  hasSpec(): boolean;
  clearSpec(): void;
  getSpec(): github_com_solo_io_solo_apis_api_rate_limiter_v1alpha1_ratelimit_pb.RateLimitConfigSpec | undefined;
  setSpec(value?: github_com_solo_io_solo_apis_api_rate_limiter_v1alpha1_ratelimit_pb.RateLimitConfigSpec): void;

  hasStatus(): boolean;
  clearStatus(): void;
  getStatus(): github_com_solo_io_solo_apis_api_rate_limiter_v1alpha1_ratelimit_pb.RateLimitConfigStatus | undefined;
  setStatus(value?: github_com_solo_io_solo_apis_api_rate_limiter_v1alpha1_ratelimit_pb.RateLimitConfigStatus): void;

  hasGlooInstance(): boolean;
  clearGlooInstance(): void;
  getGlooInstance(): github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef | undefined;
  setGlooInstance(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef): void;

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
    metadata?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ObjectMeta.AsObject,
    spec?: github_com_solo_io_solo_apis_api_rate_limiter_v1alpha1_ratelimit_pb.RateLimitConfigSpec.AsObject,
    status?: github_com_solo_io_solo_apis_api_rate_limiter_v1alpha1_ratelimit_pb.RateLimitConfigStatus.AsObject,
    glooInstance?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef.AsObject,
  }
}

export class ListRateLimitConfigsRequest extends jspb.Message {
  hasGlooInstanceRef(): boolean;
  clearGlooInstanceRef(): void;
  getGlooInstanceRef(): github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef | undefined;
  setGlooInstanceRef(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListRateLimitConfigsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListRateLimitConfigsRequest): ListRateLimitConfigsRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListRateLimitConfigsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListRateLimitConfigsRequest;
  static deserializeBinaryFromReader(message: ListRateLimitConfigsRequest, reader: jspb.BinaryReader): ListRateLimitConfigsRequest;
}

export namespace ListRateLimitConfigsRequest {
  export type AsObject = {
    glooInstanceRef?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef.AsObject,
  }
}

export class ListRateLimitConfigsResponse extends jspb.Message {
  clearRateLimitConfigsList(): void;
  getRateLimitConfigsList(): Array<RateLimitConfig>;
  setRateLimitConfigsList(value: Array<RateLimitConfig>): void;
  addRateLimitConfigs(value?: RateLimitConfig, index?: number): RateLimitConfig;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListRateLimitConfigsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListRateLimitConfigsResponse): ListRateLimitConfigsResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListRateLimitConfigsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListRateLimitConfigsResponse;
  static deserializeBinaryFromReader(message: ListRateLimitConfigsResponse, reader: jspb.BinaryReader): ListRateLimitConfigsResponse;
}

export namespace ListRateLimitConfigsResponse {
  export type AsObject = {
    rateLimitConfigsList: Array<RateLimitConfig.AsObject>,
  }
}

export class GetRateLimitConfigYamlRequest extends jspb.Message {
  hasRateLimitConfigRef(): boolean;
  clearRateLimitConfigRef(): void;
  getRateLimitConfigRef(): github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef | undefined;
  setRateLimitConfigRef(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetRateLimitConfigYamlRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetRateLimitConfigYamlRequest): GetRateLimitConfigYamlRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetRateLimitConfigYamlRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetRateLimitConfigYamlRequest;
  static deserializeBinaryFromReader(message: GetRateLimitConfigYamlRequest, reader: jspb.BinaryReader): GetRateLimitConfigYamlRequest;
}

export namespace GetRateLimitConfigYamlRequest {
  export type AsObject = {
    rateLimitConfigRef?: github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef.AsObject,
  }
}

export class GetRateLimitConfigYamlResponse extends jspb.Message {
  hasYamlData(): boolean;
  clearYamlData(): void;
  getYamlData(): github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ResourceYaml | undefined;
  setYamlData(value?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ResourceYaml): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetRateLimitConfigYamlResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetRateLimitConfigYamlResponse): GetRateLimitConfigYamlResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetRateLimitConfigYamlResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetRateLimitConfigYamlResponse;
  static deserializeBinaryFromReader(message: GetRateLimitConfigYamlResponse, reader: jspb.BinaryReader): GetRateLimitConfigYamlResponse;
}

export namespace GetRateLimitConfigYamlResponse {
  export type AsObject = {
    yamlData?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ResourceYaml.AsObject,
  }
}
