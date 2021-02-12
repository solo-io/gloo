/* eslint-disable */
// package: fed.rpc.solo.io
// file: github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/federated_enterprise_gloo_resources.proto

import * as jspb from "google-protobuf";
import * as extproto_ext_pb from "../../../../../../../../extproto/ext_pb";
import * as github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_enterprise_gloo_v1_auth_config_pb from "../../../../../../../../github.com/solo-io/solo-projects/projects/gloo-fed/api/fed.enterprise.gloo/v1/auth_config_pb";
import * as github_com_solo_io_skv2_api_core_v1_core_pb from "../../../../../../../../github.com/solo-io/skv2/api/core/v1/core_pb";
import * as github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb from "../../../../../../../../github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/common_pb";

export class FederatedAuthConfig extends jspb.Message {
  hasMetadata(): boolean;
  clearMetadata(): void;
  getMetadata(): github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb.ObjectMeta | undefined;
  setMetadata(value?: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb.ObjectMeta): void;

  hasSpec(): boolean;
  clearSpec(): void;
  getSpec(): github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_enterprise_gloo_v1_auth_config_pb.FederatedAuthConfigSpec | undefined;
  setSpec(value?: github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_enterprise_gloo_v1_auth_config_pb.FederatedAuthConfigSpec): void;

  hasStatus(): boolean;
  clearStatus(): void;
  getStatus(): github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_enterprise_gloo_v1_auth_config_pb.FederatedAuthConfigStatus | undefined;
  setStatus(value?: github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_enterprise_gloo_v1_auth_config_pb.FederatedAuthConfigStatus): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FederatedAuthConfig.AsObject;
  static toObject(includeInstance: boolean, msg: FederatedAuthConfig): FederatedAuthConfig.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: FederatedAuthConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FederatedAuthConfig;
  static deserializeBinaryFromReader(message: FederatedAuthConfig, reader: jspb.BinaryReader): FederatedAuthConfig;
}

export namespace FederatedAuthConfig {
  export type AsObject = {
    metadata?: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb.ObjectMeta.AsObject,
    spec?: github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_enterprise_gloo_v1_auth_config_pb.FederatedAuthConfigSpec.AsObject,
    status?: github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_enterprise_gloo_v1_auth_config_pb.FederatedAuthConfigStatus.AsObject,
  }
}

export class ListFederatedAuthConfigsRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListFederatedAuthConfigsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListFederatedAuthConfigsRequest): ListFederatedAuthConfigsRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListFederatedAuthConfigsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListFederatedAuthConfigsRequest;
  static deserializeBinaryFromReader(message: ListFederatedAuthConfigsRequest, reader: jspb.BinaryReader): ListFederatedAuthConfigsRequest;
}

export namespace ListFederatedAuthConfigsRequest {
  export type AsObject = {
  }
}

export class ListFederatedAuthConfigsResponse extends jspb.Message {
  clearFederatedAuthConfigsList(): void;
  getFederatedAuthConfigsList(): Array<FederatedAuthConfig>;
  setFederatedAuthConfigsList(value: Array<FederatedAuthConfig>): void;
  addFederatedAuthConfigs(value?: FederatedAuthConfig, index?: number): FederatedAuthConfig;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListFederatedAuthConfigsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListFederatedAuthConfigsResponse): ListFederatedAuthConfigsResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListFederatedAuthConfigsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListFederatedAuthConfigsResponse;
  static deserializeBinaryFromReader(message: ListFederatedAuthConfigsResponse, reader: jspb.BinaryReader): ListFederatedAuthConfigsResponse;
}

export namespace ListFederatedAuthConfigsResponse {
  export type AsObject = {
    federatedAuthConfigsList: Array<FederatedAuthConfig.AsObject>,
  }
}

export class GetFederatedAuthConfigYamlRequest extends jspb.Message {
  hasFederatedAuthConfigRef(): boolean;
  clearFederatedAuthConfigRef(): void;
  getFederatedAuthConfigRef(): github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef | undefined;
  setFederatedAuthConfigRef(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetFederatedAuthConfigYamlRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetFederatedAuthConfigYamlRequest): GetFederatedAuthConfigYamlRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetFederatedAuthConfigYamlRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetFederatedAuthConfigYamlRequest;
  static deserializeBinaryFromReader(message: GetFederatedAuthConfigYamlRequest, reader: jspb.BinaryReader): GetFederatedAuthConfigYamlRequest;
}

export namespace GetFederatedAuthConfigYamlRequest {
  export type AsObject = {
    federatedAuthConfigRef?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef.AsObject,
  }
}

export class GetFederatedAuthConfigYamlResponse extends jspb.Message {
  hasYamlData(): boolean;
  clearYamlData(): void;
  getYamlData(): github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb.ResourceYaml | undefined;
  setYamlData(value?: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb.ResourceYaml): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetFederatedAuthConfigYamlResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetFederatedAuthConfigYamlResponse): GetFederatedAuthConfigYamlResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetFederatedAuthConfigYamlResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetFederatedAuthConfigYamlResponse;
  static deserializeBinaryFromReader(message: GetFederatedAuthConfigYamlResponse, reader: jspb.BinaryReader): GetFederatedAuthConfigYamlResponse;
}

export namespace GetFederatedAuthConfigYamlResponse {
  export type AsObject = {
    yamlData?: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb.ResourceYaml.AsObject,
  }
}
