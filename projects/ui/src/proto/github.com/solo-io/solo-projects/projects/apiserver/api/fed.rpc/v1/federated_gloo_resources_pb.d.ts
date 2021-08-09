/* eslint-disable */
// package: fed.rpc.solo.io
// file: github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/federated_gloo_resources.proto

import * as jspb from "google-protobuf";
import * as extproto_ext_pb from "../../../../../../../../extproto/ext_pb";
import * as github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_gloo_v1_upstream_pb from "../../../../../../../../github.com/solo-io/solo-projects/projects/gloo-fed/api/fed.gloo/v1/upstream_pb";
import * as github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_gloo_v1_upstream_group_pb from "../../../../../../../../github.com/solo-io/solo-projects/projects/gloo-fed/api/fed.gloo/v1/upstream_group_pb";
import * as github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_gloo_v1_settings_pb from "../../../../../../../../github.com/solo-io/solo-projects/projects/gloo-fed/api/fed.gloo/v1/settings_pb";
import * as github_com_solo_io_skv2_api_core_v1_core_pb from "../../../../../../../../github.com/solo-io/skv2/api/core/v1/core_pb";
import * as github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb from "../../../../../../../../github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/common_pb";

export class FederatedUpstream extends jspb.Message {
  hasMetadata(): boolean;
  clearMetadata(): void;
  getMetadata(): github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ObjectMeta | undefined;
  setMetadata(value?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ObjectMeta): void;

  hasSpec(): boolean;
  clearSpec(): void;
  getSpec(): github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_gloo_v1_upstream_pb.FederatedUpstreamSpec | undefined;
  setSpec(value?: github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_gloo_v1_upstream_pb.FederatedUpstreamSpec): void;

  hasStatus(): boolean;
  clearStatus(): void;
  getStatus(): github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_gloo_v1_upstream_pb.FederatedUpstreamStatus | undefined;
  setStatus(value?: github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_gloo_v1_upstream_pb.FederatedUpstreamStatus): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FederatedUpstream.AsObject;
  static toObject(includeInstance: boolean, msg: FederatedUpstream): FederatedUpstream.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: FederatedUpstream, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FederatedUpstream;
  static deserializeBinaryFromReader(message: FederatedUpstream, reader: jspb.BinaryReader): FederatedUpstream;
}

export namespace FederatedUpstream {
  export type AsObject = {
    metadata?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ObjectMeta.AsObject,
    spec?: github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_gloo_v1_upstream_pb.FederatedUpstreamSpec.AsObject,
    status?: github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_gloo_v1_upstream_pb.FederatedUpstreamStatus.AsObject,
  }
}

export class FederatedUpstreamGroup extends jspb.Message {
  hasMetadata(): boolean;
  clearMetadata(): void;
  getMetadata(): github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ObjectMeta | undefined;
  setMetadata(value?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ObjectMeta): void;

  hasSpec(): boolean;
  clearSpec(): void;
  getSpec(): github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_gloo_v1_upstream_group_pb.FederatedUpstreamGroupSpec | undefined;
  setSpec(value?: github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_gloo_v1_upstream_group_pb.FederatedUpstreamGroupSpec): void;

  hasStatus(): boolean;
  clearStatus(): void;
  getStatus(): github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_gloo_v1_upstream_group_pb.FederatedUpstreamGroupStatus | undefined;
  setStatus(value?: github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_gloo_v1_upstream_group_pb.FederatedUpstreamGroupStatus): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FederatedUpstreamGroup.AsObject;
  static toObject(includeInstance: boolean, msg: FederatedUpstreamGroup): FederatedUpstreamGroup.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: FederatedUpstreamGroup, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FederatedUpstreamGroup;
  static deserializeBinaryFromReader(message: FederatedUpstreamGroup, reader: jspb.BinaryReader): FederatedUpstreamGroup;
}

export namespace FederatedUpstreamGroup {
  export type AsObject = {
    metadata?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ObjectMeta.AsObject,
    spec?: github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_gloo_v1_upstream_group_pb.FederatedUpstreamGroupSpec.AsObject,
    status?: github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_gloo_v1_upstream_group_pb.FederatedUpstreamGroupStatus.AsObject,
  }
}

export class FederatedSettings extends jspb.Message {
  hasMetadata(): boolean;
  clearMetadata(): void;
  getMetadata(): github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ObjectMeta | undefined;
  setMetadata(value?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ObjectMeta): void;

  hasSpec(): boolean;
  clearSpec(): void;
  getSpec(): github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_gloo_v1_settings_pb.FederatedSettingsSpec | undefined;
  setSpec(value?: github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_gloo_v1_settings_pb.FederatedSettingsSpec): void;

  hasStatus(): boolean;
  clearStatus(): void;
  getStatus(): github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_gloo_v1_settings_pb.FederatedSettingsStatus | undefined;
  setStatus(value?: github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_gloo_v1_settings_pb.FederatedSettingsStatus): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FederatedSettings.AsObject;
  static toObject(includeInstance: boolean, msg: FederatedSettings): FederatedSettings.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: FederatedSettings, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FederatedSettings;
  static deserializeBinaryFromReader(message: FederatedSettings, reader: jspb.BinaryReader): FederatedSettings;
}

export namespace FederatedSettings {
  export type AsObject = {
    metadata?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ObjectMeta.AsObject,
    spec?: github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_gloo_v1_settings_pb.FederatedSettingsSpec.AsObject,
    status?: github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_gloo_v1_settings_pb.FederatedSettingsStatus.AsObject,
  }
}

export class ListFederatedUpstreamsRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListFederatedUpstreamsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListFederatedUpstreamsRequest): ListFederatedUpstreamsRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListFederatedUpstreamsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListFederatedUpstreamsRequest;
  static deserializeBinaryFromReader(message: ListFederatedUpstreamsRequest, reader: jspb.BinaryReader): ListFederatedUpstreamsRequest;
}

export namespace ListFederatedUpstreamsRequest {
  export type AsObject = {
  }
}

export class ListFederatedUpstreamsResponse extends jspb.Message {
  clearFederatedUpstreamsList(): void;
  getFederatedUpstreamsList(): Array<FederatedUpstream>;
  setFederatedUpstreamsList(value: Array<FederatedUpstream>): void;
  addFederatedUpstreams(value?: FederatedUpstream, index?: number): FederatedUpstream;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListFederatedUpstreamsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListFederatedUpstreamsResponse): ListFederatedUpstreamsResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListFederatedUpstreamsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListFederatedUpstreamsResponse;
  static deserializeBinaryFromReader(message: ListFederatedUpstreamsResponse, reader: jspb.BinaryReader): ListFederatedUpstreamsResponse;
}

export namespace ListFederatedUpstreamsResponse {
  export type AsObject = {
    federatedUpstreamsList: Array<FederatedUpstream.AsObject>,
  }
}

export class GetFederatedUpstreamYamlRequest extends jspb.Message {
  hasFederatedUpstreamRef(): boolean;
  clearFederatedUpstreamRef(): void;
  getFederatedUpstreamRef(): github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef | undefined;
  setFederatedUpstreamRef(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetFederatedUpstreamYamlRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetFederatedUpstreamYamlRequest): GetFederatedUpstreamYamlRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetFederatedUpstreamYamlRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetFederatedUpstreamYamlRequest;
  static deserializeBinaryFromReader(message: GetFederatedUpstreamYamlRequest, reader: jspb.BinaryReader): GetFederatedUpstreamYamlRequest;
}

export namespace GetFederatedUpstreamYamlRequest {
  export type AsObject = {
    federatedUpstreamRef?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef.AsObject,
  }
}

export class GetFederatedUpstreamYamlResponse extends jspb.Message {
  hasYamlData(): boolean;
  clearYamlData(): void;
  getYamlData(): github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ResourceYaml | undefined;
  setYamlData(value?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ResourceYaml): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetFederatedUpstreamYamlResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetFederatedUpstreamYamlResponse): GetFederatedUpstreamYamlResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetFederatedUpstreamYamlResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetFederatedUpstreamYamlResponse;
  static deserializeBinaryFromReader(message: GetFederatedUpstreamYamlResponse, reader: jspb.BinaryReader): GetFederatedUpstreamYamlResponse;
}

export namespace GetFederatedUpstreamYamlResponse {
  export type AsObject = {
    yamlData?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ResourceYaml.AsObject,
  }
}

export class ListFederatedUpstreamGroupsRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListFederatedUpstreamGroupsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListFederatedUpstreamGroupsRequest): ListFederatedUpstreamGroupsRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListFederatedUpstreamGroupsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListFederatedUpstreamGroupsRequest;
  static deserializeBinaryFromReader(message: ListFederatedUpstreamGroupsRequest, reader: jspb.BinaryReader): ListFederatedUpstreamGroupsRequest;
}

export namespace ListFederatedUpstreamGroupsRequest {
  export type AsObject = {
  }
}

export class ListFederatedUpstreamGroupsResponse extends jspb.Message {
  clearFederatedUpstreamGroupsList(): void;
  getFederatedUpstreamGroupsList(): Array<FederatedUpstreamGroup>;
  setFederatedUpstreamGroupsList(value: Array<FederatedUpstreamGroup>): void;
  addFederatedUpstreamGroups(value?: FederatedUpstreamGroup, index?: number): FederatedUpstreamGroup;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListFederatedUpstreamGroupsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListFederatedUpstreamGroupsResponse): ListFederatedUpstreamGroupsResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListFederatedUpstreamGroupsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListFederatedUpstreamGroupsResponse;
  static deserializeBinaryFromReader(message: ListFederatedUpstreamGroupsResponse, reader: jspb.BinaryReader): ListFederatedUpstreamGroupsResponse;
}

export namespace ListFederatedUpstreamGroupsResponse {
  export type AsObject = {
    federatedUpstreamGroupsList: Array<FederatedUpstreamGroup.AsObject>,
  }
}

export class GetFederatedUpstreamGroupYamlRequest extends jspb.Message {
  hasFederatedUpstreamGroupRef(): boolean;
  clearFederatedUpstreamGroupRef(): void;
  getFederatedUpstreamGroupRef(): github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef | undefined;
  setFederatedUpstreamGroupRef(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetFederatedUpstreamGroupYamlRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetFederatedUpstreamGroupYamlRequest): GetFederatedUpstreamGroupYamlRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetFederatedUpstreamGroupYamlRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetFederatedUpstreamGroupYamlRequest;
  static deserializeBinaryFromReader(message: GetFederatedUpstreamGroupYamlRequest, reader: jspb.BinaryReader): GetFederatedUpstreamGroupYamlRequest;
}

export namespace GetFederatedUpstreamGroupYamlRequest {
  export type AsObject = {
    federatedUpstreamGroupRef?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef.AsObject,
  }
}

export class GetFederatedUpstreamGroupYamlResponse extends jspb.Message {
  hasYamlData(): boolean;
  clearYamlData(): void;
  getYamlData(): github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ResourceYaml | undefined;
  setYamlData(value?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ResourceYaml): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetFederatedUpstreamGroupYamlResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetFederatedUpstreamGroupYamlResponse): GetFederatedUpstreamGroupYamlResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetFederatedUpstreamGroupYamlResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetFederatedUpstreamGroupYamlResponse;
  static deserializeBinaryFromReader(message: GetFederatedUpstreamGroupYamlResponse, reader: jspb.BinaryReader): GetFederatedUpstreamGroupYamlResponse;
}

export namespace GetFederatedUpstreamGroupYamlResponse {
  export type AsObject = {
    yamlData?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ResourceYaml.AsObject,
  }
}

export class ListFederatedSettingsRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListFederatedSettingsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListFederatedSettingsRequest): ListFederatedSettingsRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListFederatedSettingsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListFederatedSettingsRequest;
  static deserializeBinaryFromReader(message: ListFederatedSettingsRequest, reader: jspb.BinaryReader): ListFederatedSettingsRequest;
}

export namespace ListFederatedSettingsRequest {
  export type AsObject = {
  }
}

export class ListFederatedSettingsResponse extends jspb.Message {
  clearFederatedSettingsList(): void;
  getFederatedSettingsList(): Array<FederatedSettings>;
  setFederatedSettingsList(value: Array<FederatedSettings>): void;
  addFederatedSettings(value?: FederatedSettings, index?: number): FederatedSettings;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListFederatedSettingsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListFederatedSettingsResponse): ListFederatedSettingsResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListFederatedSettingsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListFederatedSettingsResponse;
  static deserializeBinaryFromReader(message: ListFederatedSettingsResponse, reader: jspb.BinaryReader): ListFederatedSettingsResponse;
}

export namespace ListFederatedSettingsResponse {
  export type AsObject = {
    federatedSettingsList: Array<FederatedSettings.AsObject>,
  }
}

export class GetFederatedSettingsYamlRequest extends jspb.Message {
  hasFederatedSettingsRef(): boolean;
  clearFederatedSettingsRef(): void;
  getFederatedSettingsRef(): github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef | undefined;
  setFederatedSettingsRef(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetFederatedSettingsYamlRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetFederatedSettingsYamlRequest): GetFederatedSettingsYamlRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetFederatedSettingsYamlRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetFederatedSettingsYamlRequest;
  static deserializeBinaryFromReader(message: GetFederatedSettingsYamlRequest, reader: jspb.BinaryReader): GetFederatedSettingsYamlRequest;
}

export namespace GetFederatedSettingsYamlRequest {
  export type AsObject = {
    federatedSettingsRef?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef.AsObject,
  }
}

export class GetFederatedSettingsYamlResponse extends jspb.Message {
  hasYamlData(): boolean;
  clearYamlData(): void;
  getYamlData(): github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ResourceYaml | undefined;
  setYamlData(value?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ResourceYaml): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetFederatedSettingsYamlResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetFederatedSettingsYamlResponse): GetFederatedSettingsYamlResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetFederatedSettingsYamlResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetFederatedSettingsYamlResponse;
  static deserializeBinaryFromReader(message: GetFederatedSettingsYamlResponse, reader: jspb.BinaryReader): GetFederatedSettingsYamlResponse;
}

export namespace GetFederatedSettingsYamlResponse {
  export type AsObject = {
    yamlData?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ResourceYaml.AsObject,
  }
}
