/* eslint-disable */
// package: fed.rpc.solo.io
// file: github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/enterprise_gloo_resources.proto

import * as jspb from "google-protobuf";
import * as extproto_ext_pb from "../../../../../../../../extproto/ext_pb";
import * as github_com_solo_io_solo_apis_api_gloo_enterprise_gloo_v1_auth_config_pb from "../../../../../../../../github.com/solo-io/solo-apis/api/gloo/enterprise.gloo/v1/auth_config_pb";
import * as github_com_solo_io_skv2_api_core_v1_core_pb from "../../../../../../../../github.com/solo-io/skv2/api/core/v1/core_pb";
import * as github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb from "../../../../../../../../github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/common_pb";

export class AuthConfig extends jspb.Message {
  hasMetadata(): boolean;
  clearMetadata(): void;
  getMetadata(): github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb.ObjectMeta | undefined;
  setMetadata(value?: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb.ObjectMeta): void;

  hasSpec(): boolean;
  clearSpec(): void;
  getSpec(): github_com_solo_io_solo_apis_api_gloo_enterprise_gloo_v1_auth_config_pb.AuthConfigSpec | undefined;
  setSpec(value?: github_com_solo_io_solo_apis_api_gloo_enterprise_gloo_v1_auth_config_pb.AuthConfigSpec): void;

  hasStatus(): boolean;
  clearStatus(): void;
  getStatus(): github_com_solo_io_solo_apis_api_gloo_enterprise_gloo_v1_auth_config_pb.AuthConfigStatus | undefined;
  setStatus(value?: github_com_solo_io_solo_apis_api_gloo_enterprise_gloo_v1_auth_config_pb.AuthConfigStatus): void;

  hasGlooInstance(): boolean;
  clearGlooInstance(): void;
  getGlooInstance(): github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef | undefined;
  setGlooInstance(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AuthConfig.AsObject;
  static toObject(includeInstance: boolean, msg: AuthConfig): AuthConfig.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: AuthConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AuthConfig;
  static deserializeBinaryFromReader(message: AuthConfig, reader: jspb.BinaryReader): AuthConfig;
}

export namespace AuthConfig {
  export type AsObject = {
    metadata?: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb.ObjectMeta.AsObject,
    spec?: github_com_solo_io_solo_apis_api_gloo_enterprise_gloo_v1_auth_config_pb.AuthConfigSpec.AsObject,
    status?: github_com_solo_io_solo_apis_api_gloo_enterprise_gloo_v1_auth_config_pb.AuthConfigStatus.AsObject,
    glooInstance?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef.AsObject,
  }
}

export class ListAuthConfigsRequest extends jspb.Message {
  hasGlooInstanceRef(): boolean;
  clearGlooInstanceRef(): void;
  getGlooInstanceRef(): github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef | undefined;
  setGlooInstanceRef(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListAuthConfigsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListAuthConfigsRequest): ListAuthConfigsRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListAuthConfigsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListAuthConfigsRequest;
  static deserializeBinaryFromReader(message: ListAuthConfigsRequest, reader: jspb.BinaryReader): ListAuthConfigsRequest;
}

export namespace ListAuthConfigsRequest {
  export type AsObject = {
    glooInstanceRef?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef.AsObject,
  }
}

export class ListAuthConfigsResponse extends jspb.Message {
  clearAuthConfigsList(): void;
  getAuthConfigsList(): Array<AuthConfig>;
  setAuthConfigsList(value: Array<AuthConfig>): void;
  addAuthConfigs(value?: AuthConfig, index?: number): AuthConfig;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListAuthConfigsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListAuthConfigsResponse): ListAuthConfigsResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListAuthConfigsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListAuthConfigsResponse;
  static deserializeBinaryFromReader(message: ListAuthConfigsResponse, reader: jspb.BinaryReader): ListAuthConfigsResponse;
}

export namespace ListAuthConfigsResponse {
  export type AsObject = {
    authConfigsList: Array<AuthConfig.AsObject>,
  }
}

export class GetAuthConfigYamlRequest extends jspb.Message {
  hasAuthConfigRef(): boolean;
  clearAuthConfigRef(): void;
  getAuthConfigRef(): github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef | undefined;
  setAuthConfigRef(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetAuthConfigYamlRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetAuthConfigYamlRequest): GetAuthConfigYamlRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetAuthConfigYamlRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetAuthConfigYamlRequest;
  static deserializeBinaryFromReader(message: GetAuthConfigYamlRequest, reader: jspb.BinaryReader): GetAuthConfigYamlRequest;
}

export namespace GetAuthConfigYamlRequest {
  export type AsObject = {
    authConfigRef?: github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef.AsObject,
  }
}

export class GetAuthConfigYamlResponse extends jspb.Message {
  hasYamlData(): boolean;
  clearYamlData(): void;
  getYamlData(): github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb.ResourceYaml | undefined;
  setYamlData(value?: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb.ResourceYaml): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetAuthConfigYamlResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetAuthConfigYamlResponse): GetAuthConfigYamlResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetAuthConfigYamlResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetAuthConfigYamlResponse;
  static deserializeBinaryFromReader(message: GetAuthConfigYamlResponse, reader: jspb.BinaryReader): GetAuthConfigYamlResponse;
}

export namespace GetAuthConfigYamlResponse {
  export type AsObject = {
    yamlData?: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb.ResourceYaml.AsObject,
  }
}
