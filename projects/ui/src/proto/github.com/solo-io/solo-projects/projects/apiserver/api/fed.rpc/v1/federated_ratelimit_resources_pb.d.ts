/* eslint-disable */
// package: fed.rpc.solo.io
// file: github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/federated_ratelimit_resources.proto

import * as jspb from "google-protobuf";
import * as extproto_ext_pb from "../../../../../../../../extproto/ext_pb";
import * as github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_ratelimit_v1alpha1_rate_limit_config_pb from "../../../../../../../../github.com/solo-io/solo-projects/projects/gloo-fed/api/fed.ratelimit/v1alpha1/rate_limit_config_pb";
import * as github_com_solo_io_skv2_api_core_v1_core_pb from "../../../../../../../../github.com/solo-io/skv2/api/core/v1/core_pb";
import * as github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb from "../../../../../../../../github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/common_pb";

export class FederatedRateLimitConfig extends jspb.Message {
  hasMetadata(): boolean;
  clearMetadata(): void;
  getMetadata(): github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb.ObjectMeta | undefined;
  setMetadata(value?: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb.ObjectMeta): void;

  hasSpec(): boolean;
  clearSpec(): void;
  getSpec(): github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_ratelimit_v1alpha1_rate_limit_config_pb.FederatedRateLimitConfigSpec | undefined;
  setSpec(value?: github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_ratelimit_v1alpha1_rate_limit_config_pb.FederatedRateLimitConfigSpec): void;

  hasStatus(): boolean;
  clearStatus(): void;
  getStatus(): github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_ratelimit_v1alpha1_rate_limit_config_pb.FederatedRateLimitConfigStatus | undefined;
  setStatus(value?: github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_ratelimit_v1alpha1_rate_limit_config_pb.FederatedRateLimitConfigStatus): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FederatedRateLimitConfig.AsObject;
  static toObject(includeInstance: boolean, msg: FederatedRateLimitConfig): FederatedRateLimitConfig.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: FederatedRateLimitConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FederatedRateLimitConfig;
  static deserializeBinaryFromReader(message: FederatedRateLimitConfig, reader: jspb.BinaryReader): FederatedRateLimitConfig;
}

export namespace FederatedRateLimitConfig {
  export type AsObject = {
    metadata?: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb.ObjectMeta.AsObject,
    spec?: github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_ratelimit_v1alpha1_rate_limit_config_pb.FederatedRateLimitConfigSpec.AsObject,
    status?: github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_ratelimit_v1alpha1_rate_limit_config_pb.FederatedRateLimitConfigStatus.AsObject,
  }
}

export class ListFederatedRateLimitConfigsRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListFederatedRateLimitConfigsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListFederatedRateLimitConfigsRequest): ListFederatedRateLimitConfigsRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListFederatedRateLimitConfigsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListFederatedRateLimitConfigsRequest;
  static deserializeBinaryFromReader(message: ListFederatedRateLimitConfigsRequest, reader: jspb.BinaryReader): ListFederatedRateLimitConfigsRequest;
}

export namespace ListFederatedRateLimitConfigsRequest {
  export type AsObject = {
  }
}

export class ListFederatedRateLimitConfigsResponse extends jspb.Message {
  clearFederatedRateLimitConfigsList(): void;
  getFederatedRateLimitConfigsList(): Array<FederatedRateLimitConfig>;
  setFederatedRateLimitConfigsList(value: Array<FederatedRateLimitConfig>): void;
  addFederatedRateLimitConfigs(value?: FederatedRateLimitConfig, index?: number): FederatedRateLimitConfig;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListFederatedRateLimitConfigsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListFederatedRateLimitConfigsResponse): ListFederatedRateLimitConfigsResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListFederatedRateLimitConfigsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListFederatedRateLimitConfigsResponse;
  static deserializeBinaryFromReader(message: ListFederatedRateLimitConfigsResponse, reader: jspb.BinaryReader): ListFederatedRateLimitConfigsResponse;
}

export namespace ListFederatedRateLimitConfigsResponse {
  export type AsObject = {
    federatedRateLimitConfigsList: Array<FederatedRateLimitConfig.AsObject>,
  }
}

export class GetFederatedRateLimitConfigYamlRequest extends jspb.Message {
  hasFederatedRateLimitConfigRef(): boolean;
  clearFederatedRateLimitConfigRef(): void;
  getFederatedRateLimitConfigRef(): github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef | undefined;
  setFederatedRateLimitConfigRef(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetFederatedRateLimitConfigYamlRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetFederatedRateLimitConfigYamlRequest): GetFederatedRateLimitConfigYamlRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetFederatedRateLimitConfigYamlRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetFederatedRateLimitConfigYamlRequest;
  static deserializeBinaryFromReader(message: GetFederatedRateLimitConfigYamlRequest, reader: jspb.BinaryReader): GetFederatedRateLimitConfigYamlRequest;
}

export namespace GetFederatedRateLimitConfigYamlRequest {
  export type AsObject = {
    federatedRateLimitConfigRef?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef.AsObject,
  }
}

export class GetFederatedRateLimitConfigYamlResponse extends jspb.Message {
  hasYamlData(): boolean;
  clearYamlData(): void;
  getYamlData(): github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb.ResourceYaml | undefined;
  setYamlData(value?: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb.ResourceYaml): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetFederatedRateLimitConfigYamlResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetFederatedRateLimitConfigYamlResponse): GetFederatedRateLimitConfigYamlResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetFederatedRateLimitConfigYamlResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetFederatedRateLimitConfigYamlResponse;
  static deserializeBinaryFromReader(message: GetFederatedRateLimitConfigYamlResponse, reader: jspb.BinaryReader): GetFederatedRateLimitConfigYamlResponse;
}

export namespace GetFederatedRateLimitConfigYamlResponse {
  export type AsObject = {
    yamlData?: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb.ResourceYaml.AsObject,
  }
}
