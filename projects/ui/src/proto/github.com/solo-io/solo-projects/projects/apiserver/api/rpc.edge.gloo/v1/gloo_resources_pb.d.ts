/* eslint-disable */
// package: rpc.edge.gloo.solo.io
// file: github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/gloo_resources.proto

import * as jspb from "google-protobuf";
import * as extproto_ext_pb from "../../../../../../../../extproto/ext_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_upstream_pb from "../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/upstream_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_settings_pb from "../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/settings_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_proxy_pb from "../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/proxy_pb";
import * as github_com_solo_io_skv2_api_core_v1_core_pb from "../../../../../../../../github.com/solo-io/skv2/api/core/v1/core_pb";
import * as github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb from "../../../../../../../../github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/common_pb";

export class Upstream extends jspb.Message {
  hasMetadata(): boolean;
  clearMetadata(): void;
  getMetadata(): github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ObjectMeta | undefined;
  setMetadata(value?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ObjectMeta): void;

  hasSpec(): boolean;
  clearSpec(): void;
  getSpec(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_upstream_pb.UpstreamSpec | undefined;
  setSpec(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_upstream_pb.UpstreamSpec): void;

  hasStatus(): boolean;
  clearStatus(): void;
  getStatus(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_upstream_pb.UpstreamStatus | undefined;
  setStatus(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_upstream_pb.UpstreamStatus): void;

  hasGlooInstance(): boolean;
  clearGlooInstance(): void;
  getGlooInstance(): github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef | undefined;
  setGlooInstance(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Upstream.AsObject;
  static toObject(includeInstance: boolean, msg: Upstream): Upstream.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Upstream, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Upstream;
  static deserializeBinaryFromReader(message: Upstream, reader: jspb.BinaryReader): Upstream;
}

export namespace Upstream {
  export type AsObject = {
    metadata?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ObjectMeta.AsObject,
    spec?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_upstream_pb.UpstreamSpec.AsObject,
    status?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_upstream_pb.UpstreamStatus.AsObject,
    glooInstance?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef.AsObject,
  }
}

export class UpstreamGroup extends jspb.Message {
  hasMetadata(): boolean;
  clearMetadata(): void;
  getMetadata(): github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ObjectMeta | undefined;
  setMetadata(value?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ObjectMeta): void;

  hasSpec(): boolean;
  clearSpec(): void;
  getSpec(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_proxy_pb.UpstreamGroupSpec | undefined;
  setSpec(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_proxy_pb.UpstreamGroupSpec): void;

  hasStatus(): boolean;
  clearStatus(): void;
  getStatus(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_proxy_pb.UpstreamGroupStatus | undefined;
  setStatus(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_proxy_pb.UpstreamGroupStatus): void;

  hasGlooInstance(): boolean;
  clearGlooInstance(): void;
  getGlooInstance(): github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef | undefined;
  setGlooInstance(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpstreamGroup.AsObject;
  static toObject(includeInstance: boolean, msg: UpstreamGroup): UpstreamGroup.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpstreamGroup, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpstreamGroup;
  static deserializeBinaryFromReader(message: UpstreamGroup, reader: jspb.BinaryReader): UpstreamGroup;
}

export namespace UpstreamGroup {
  export type AsObject = {
    metadata?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ObjectMeta.AsObject,
    spec?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_proxy_pb.UpstreamGroupSpec.AsObject,
    status?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_proxy_pb.UpstreamGroupStatus.AsObject,
    glooInstance?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef.AsObject,
  }
}

export class Settings extends jspb.Message {
  hasMetadata(): boolean;
  clearMetadata(): void;
  getMetadata(): github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ObjectMeta | undefined;
  setMetadata(value?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ObjectMeta): void;

  hasSpec(): boolean;
  clearSpec(): void;
  getSpec(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_settings_pb.SettingsSpec | undefined;
  setSpec(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_settings_pb.SettingsSpec): void;

  hasStatus(): boolean;
  clearStatus(): void;
  getStatus(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_settings_pb.SettingsStatus | undefined;
  setStatus(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_settings_pb.SettingsStatus): void;

  hasGlooInstance(): boolean;
  clearGlooInstance(): void;
  getGlooInstance(): github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef | undefined;
  setGlooInstance(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef): void;

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
    metadata?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ObjectMeta.AsObject,
    spec?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_settings_pb.SettingsSpec.AsObject,
    status?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_settings_pb.SettingsStatus.AsObject,
    glooInstance?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef.AsObject,
  }
}

export class Proxy extends jspb.Message {
  hasMetadata(): boolean;
  clearMetadata(): void;
  getMetadata(): github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ObjectMeta | undefined;
  setMetadata(value?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ObjectMeta): void;

  hasSpec(): boolean;
  clearSpec(): void;
  getSpec(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_proxy_pb.ProxySpec | undefined;
  setSpec(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_proxy_pb.ProxySpec): void;

  hasStatus(): boolean;
  clearStatus(): void;
  getStatus(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_proxy_pb.ProxyStatus | undefined;
  setStatus(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_proxy_pb.ProxyStatus): void;

  hasGlooInstance(): boolean;
  clearGlooInstance(): void;
  getGlooInstance(): github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef | undefined;
  setGlooInstance(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Proxy.AsObject;
  static toObject(includeInstance: boolean, msg: Proxy): Proxy.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Proxy, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Proxy;
  static deserializeBinaryFromReader(message: Proxy, reader: jspb.BinaryReader): Proxy;
}

export namespace Proxy {
  export type AsObject = {
    metadata?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ObjectMeta.AsObject,
    spec?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_proxy_pb.ProxySpec.AsObject,
    status?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_proxy_pb.ProxyStatus.AsObject,
    glooInstance?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef.AsObject,
  }
}

export class ListUpstreamsRequest extends jspb.Message {
  hasGlooInstanceRef(): boolean;
  clearGlooInstanceRef(): void;
  getGlooInstanceRef(): github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef | undefined;
  setGlooInstanceRef(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef): void;

  hasPagination(): boolean;
  clearPagination(): void;
  getPagination(): github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.Pagination | undefined;
  setPagination(value?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.Pagination): void;

  getQueryString(): string;
  setQueryString(value: string): void;

  hasStatusFilter(): boolean;
  clearStatusFilter(): void;
  getStatusFilter(): github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.StatusFilter | undefined;
  setStatusFilter(value?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.StatusFilter): void;

  hasSortOptions(): boolean;
  clearSortOptions(): void;
  getSortOptions(): github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.SortOptions | undefined;
  setSortOptions(value?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.SortOptions): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListUpstreamsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListUpstreamsRequest): ListUpstreamsRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListUpstreamsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListUpstreamsRequest;
  static deserializeBinaryFromReader(message: ListUpstreamsRequest, reader: jspb.BinaryReader): ListUpstreamsRequest;
}

export namespace ListUpstreamsRequest {
  export type AsObject = {
    glooInstanceRef?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef.AsObject,
    pagination?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.Pagination.AsObject,
    queryString: string,
    statusFilter?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.StatusFilter.AsObject,
    sortOptions?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.SortOptions.AsObject,
  }
}

export class ListUpstreamsResponse extends jspb.Message {
  clearUpstreamsList(): void;
  getUpstreamsList(): Array<Upstream>;
  setUpstreamsList(value: Array<Upstream>): void;
  addUpstreams(value?: Upstream, index?: number): Upstream;

  getTotal(): number;
  setTotal(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListUpstreamsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListUpstreamsResponse): ListUpstreamsResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListUpstreamsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListUpstreamsResponse;
  static deserializeBinaryFromReader(message: ListUpstreamsResponse, reader: jspb.BinaryReader): ListUpstreamsResponse;
}

export namespace ListUpstreamsResponse {
  export type AsObject = {
    upstreamsList: Array<Upstream.AsObject>,
    total: number,
  }
}

export class GetUpstreamYamlRequest extends jspb.Message {
  hasUpstreamRef(): boolean;
  clearUpstreamRef(): void;
  getUpstreamRef(): github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef | undefined;
  setUpstreamRef(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetUpstreamYamlRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetUpstreamYamlRequest): GetUpstreamYamlRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetUpstreamYamlRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetUpstreamYamlRequest;
  static deserializeBinaryFromReader(message: GetUpstreamYamlRequest, reader: jspb.BinaryReader): GetUpstreamYamlRequest;
}

export namespace GetUpstreamYamlRequest {
  export type AsObject = {
    upstreamRef?: github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef.AsObject,
  }
}

export class GetUpstreamYamlResponse extends jspb.Message {
  hasYamlData(): boolean;
  clearYamlData(): void;
  getYamlData(): github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ResourceYaml | undefined;
  setYamlData(value?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ResourceYaml): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetUpstreamYamlResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetUpstreamYamlResponse): GetUpstreamYamlResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetUpstreamYamlResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetUpstreamYamlResponse;
  static deserializeBinaryFromReader(message: GetUpstreamYamlResponse, reader: jspb.BinaryReader): GetUpstreamYamlResponse;
}

export namespace GetUpstreamYamlResponse {
  export type AsObject = {
    yamlData?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ResourceYaml.AsObject,
  }
}

export class GetUpstreamDetailsRequest extends jspb.Message {
  hasUpstreamRef(): boolean;
  clearUpstreamRef(): void;
  getUpstreamRef(): github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef | undefined;
  setUpstreamRef(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetUpstreamDetailsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetUpstreamDetailsRequest): GetUpstreamDetailsRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetUpstreamDetailsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetUpstreamDetailsRequest;
  static deserializeBinaryFromReader(message: GetUpstreamDetailsRequest, reader: jspb.BinaryReader): GetUpstreamDetailsRequest;
}

export namespace GetUpstreamDetailsRequest {
  export type AsObject = {
    upstreamRef?: github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef.AsObject,
  }
}

export class GetUpstreamDetailsResponse extends jspb.Message {
  hasUpstream(): boolean;
  clearUpstream(): void;
  getUpstream(): Upstream | undefined;
  setUpstream(value?: Upstream): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetUpstreamDetailsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetUpstreamDetailsResponse): GetUpstreamDetailsResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetUpstreamDetailsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetUpstreamDetailsResponse;
  static deserializeBinaryFromReader(message: GetUpstreamDetailsResponse, reader: jspb.BinaryReader): GetUpstreamDetailsResponse;
}

export namespace GetUpstreamDetailsResponse {
  export type AsObject = {
    upstream?: Upstream.AsObject,
  }
}

export class ListUpstreamGroupsRequest extends jspb.Message {
  hasGlooInstanceRef(): boolean;
  clearGlooInstanceRef(): void;
  getGlooInstanceRef(): github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef | undefined;
  setGlooInstanceRef(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef): void;

  hasPagination(): boolean;
  clearPagination(): void;
  getPagination(): github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.Pagination | undefined;
  setPagination(value?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.Pagination): void;

  getQueryString(): string;
  setQueryString(value: string): void;

  hasStatusFilter(): boolean;
  clearStatusFilter(): void;
  getStatusFilter(): github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.StatusFilter | undefined;
  setStatusFilter(value?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.StatusFilter): void;

  hasSortOptions(): boolean;
  clearSortOptions(): void;
  getSortOptions(): github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.SortOptions | undefined;
  setSortOptions(value?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.SortOptions): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListUpstreamGroupsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListUpstreamGroupsRequest): ListUpstreamGroupsRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListUpstreamGroupsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListUpstreamGroupsRequest;
  static deserializeBinaryFromReader(message: ListUpstreamGroupsRequest, reader: jspb.BinaryReader): ListUpstreamGroupsRequest;
}

export namespace ListUpstreamGroupsRequest {
  export type AsObject = {
    glooInstanceRef?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef.AsObject,
    pagination?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.Pagination.AsObject,
    queryString: string,
    statusFilter?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.StatusFilter.AsObject,
    sortOptions?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.SortOptions.AsObject,
  }
}

export class ListUpstreamGroupsResponse extends jspb.Message {
  clearUpstreamGroupsList(): void;
  getUpstreamGroupsList(): Array<UpstreamGroup>;
  setUpstreamGroupsList(value: Array<UpstreamGroup>): void;
  addUpstreamGroups(value?: UpstreamGroup, index?: number): UpstreamGroup;

  getTotal(): number;
  setTotal(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListUpstreamGroupsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListUpstreamGroupsResponse): ListUpstreamGroupsResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListUpstreamGroupsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListUpstreamGroupsResponse;
  static deserializeBinaryFromReader(message: ListUpstreamGroupsResponse, reader: jspb.BinaryReader): ListUpstreamGroupsResponse;
}

export namespace ListUpstreamGroupsResponse {
  export type AsObject = {
    upstreamGroupsList: Array<UpstreamGroup.AsObject>,
    total: number,
  }
}

export class GetUpstreamGroupYamlRequest extends jspb.Message {
  hasUpstreamGroupRef(): boolean;
  clearUpstreamGroupRef(): void;
  getUpstreamGroupRef(): github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef | undefined;
  setUpstreamGroupRef(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetUpstreamGroupYamlRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetUpstreamGroupYamlRequest): GetUpstreamGroupYamlRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetUpstreamGroupYamlRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetUpstreamGroupYamlRequest;
  static deserializeBinaryFromReader(message: GetUpstreamGroupYamlRequest, reader: jspb.BinaryReader): GetUpstreamGroupYamlRequest;
}

export namespace GetUpstreamGroupYamlRequest {
  export type AsObject = {
    upstreamGroupRef?: github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef.AsObject,
  }
}

export class GetUpstreamGroupYamlResponse extends jspb.Message {
  hasYamlData(): boolean;
  clearYamlData(): void;
  getYamlData(): github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ResourceYaml | undefined;
  setYamlData(value?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ResourceYaml): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetUpstreamGroupYamlResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetUpstreamGroupYamlResponse): GetUpstreamGroupYamlResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetUpstreamGroupYamlResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetUpstreamGroupYamlResponse;
  static deserializeBinaryFromReader(message: GetUpstreamGroupYamlResponse, reader: jspb.BinaryReader): GetUpstreamGroupYamlResponse;
}

export namespace GetUpstreamGroupYamlResponse {
  export type AsObject = {
    yamlData?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ResourceYaml.AsObject,
  }
}

export class GetUpstreamGroupDetailsRequest extends jspb.Message {
  hasUpstreamGroupRef(): boolean;
  clearUpstreamGroupRef(): void;
  getUpstreamGroupRef(): github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef | undefined;
  setUpstreamGroupRef(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetUpstreamGroupDetailsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetUpstreamGroupDetailsRequest): GetUpstreamGroupDetailsRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetUpstreamGroupDetailsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetUpstreamGroupDetailsRequest;
  static deserializeBinaryFromReader(message: GetUpstreamGroupDetailsRequest, reader: jspb.BinaryReader): GetUpstreamGroupDetailsRequest;
}

export namespace GetUpstreamGroupDetailsRequest {
  export type AsObject = {
    upstreamGroupRef?: github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef.AsObject,
  }
}

export class GetUpstreamGroupDetailsResponse extends jspb.Message {
  hasUpstreamGroup(): boolean;
  clearUpstreamGroup(): void;
  getUpstreamGroup(): UpstreamGroup | undefined;
  setUpstreamGroup(value?: UpstreamGroup): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetUpstreamGroupDetailsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetUpstreamGroupDetailsResponse): GetUpstreamGroupDetailsResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetUpstreamGroupDetailsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetUpstreamGroupDetailsResponse;
  static deserializeBinaryFromReader(message: GetUpstreamGroupDetailsResponse, reader: jspb.BinaryReader): GetUpstreamGroupDetailsResponse;
}

export namespace GetUpstreamGroupDetailsResponse {
  export type AsObject = {
    upstreamGroup?: UpstreamGroup.AsObject,
  }
}

export class ListSettingsRequest extends jspb.Message {
  hasGlooInstanceRef(): boolean;
  clearGlooInstanceRef(): void;
  getGlooInstanceRef(): github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef | undefined;
  setGlooInstanceRef(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef): void;

  hasPagination(): boolean;
  clearPagination(): void;
  getPagination(): github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.Pagination | undefined;
  setPagination(value?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.Pagination): void;

  getQueryString(): string;
  setQueryString(value: string): void;

  hasStatusFilter(): boolean;
  clearStatusFilter(): void;
  getStatusFilter(): github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.StatusFilter | undefined;
  setStatusFilter(value?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.StatusFilter): void;

  hasSortOptions(): boolean;
  clearSortOptions(): void;
  getSortOptions(): github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.SortOptions | undefined;
  setSortOptions(value?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.SortOptions): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListSettingsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListSettingsRequest): ListSettingsRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListSettingsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListSettingsRequest;
  static deserializeBinaryFromReader(message: ListSettingsRequest, reader: jspb.BinaryReader): ListSettingsRequest;
}

export namespace ListSettingsRequest {
  export type AsObject = {
    glooInstanceRef?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef.AsObject,
    pagination?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.Pagination.AsObject,
    queryString: string,
    statusFilter?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.StatusFilter.AsObject,
    sortOptions?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.SortOptions.AsObject,
  }
}

export class ListSettingsResponse extends jspb.Message {
  clearSettingsList(): void;
  getSettingsList(): Array<Settings>;
  setSettingsList(value: Array<Settings>): void;
  addSettings(value?: Settings, index?: number): Settings;

  getTotal(): number;
  setTotal(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListSettingsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListSettingsResponse): ListSettingsResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListSettingsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListSettingsResponse;
  static deserializeBinaryFromReader(message: ListSettingsResponse, reader: jspb.BinaryReader): ListSettingsResponse;
}

export namespace ListSettingsResponse {
  export type AsObject = {
    settingsList: Array<Settings.AsObject>,
    total: number,
  }
}

export class GetSettingsYamlRequest extends jspb.Message {
  hasSettingsRef(): boolean;
  clearSettingsRef(): void;
  getSettingsRef(): github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef | undefined;
  setSettingsRef(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetSettingsYamlRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetSettingsYamlRequest): GetSettingsYamlRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetSettingsYamlRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetSettingsYamlRequest;
  static deserializeBinaryFromReader(message: GetSettingsYamlRequest, reader: jspb.BinaryReader): GetSettingsYamlRequest;
}

export namespace GetSettingsYamlRequest {
  export type AsObject = {
    settingsRef?: github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef.AsObject,
  }
}

export class GetSettingsYamlResponse extends jspb.Message {
  hasYamlData(): boolean;
  clearYamlData(): void;
  getYamlData(): github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ResourceYaml | undefined;
  setYamlData(value?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ResourceYaml): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetSettingsYamlResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetSettingsYamlResponse): GetSettingsYamlResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetSettingsYamlResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetSettingsYamlResponse;
  static deserializeBinaryFromReader(message: GetSettingsYamlResponse, reader: jspb.BinaryReader): GetSettingsYamlResponse;
}

export namespace GetSettingsYamlResponse {
  export type AsObject = {
    yamlData?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ResourceYaml.AsObject,
  }
}

export class GetSettingsDetailsRequest extends jspb.Message {
  hasSettingsRef(): boolean;
  clearSettingsRef(): void;
  getSettingsRef(): github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef | undefined;
  setSettingsRef(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetSettingsDetailsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetSettingsDetailsRequest): GetSettingsDetailsRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetSettingsDetailsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetSettingsDetailsRequest;
  static deserializeBinaryFromReader(message: GetSettingsDetailsRequest, reader: jspb.BinaryReader): GetSettingsDetailsRequest;
}

export namespace GetSettingsDetailsRequest {
  export type AsObject = {
    settingsRef?: github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef.AsObject,
  }
}

export class GetSettingsDetailsResponse extends jspb.Message {
  hasSettings(): boolean;
  clearSettings(): void;
  getSettings(): Settings | undefined;
  setSettings(value?: Settings): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetSettingsDetailsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetSettingsDetailsResponse): GetSettingsDetailsResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetSettingsDetailsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetSettingsDetailsResponse;
  static deserializeBinaryFromReader(message: GetSettingsDetailsResponse, reader: jspb.BinaryReader): GetSettingsDetailsResponse;
}

export namespace GetSettingsDetailsResponse {
  export type AsObject = {
    settings?: Settings.AsObject,
  }
}

export class ListProxiesRequest extends jspb.Message {
  hasGlooInstanceRef(): boolean;
  clearGlooInstanceRef(): void;
  getGlooInstanceRef(): github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef | undefined;
  setGlooInstanceRef(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef): void;

  hasPagination(): boolean;
  clearPagination(): void;
  getPagination(): github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.Pagination | undefined;
  setPagination(value?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.Pagination): void;

  getQueryString(): string;
  setQueryString(value: string): void;

  hasStatusFilter(): boolean;
  clearStatusFilter(): void;
  getStatusFilter(): github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.StatusFilter | undefined;
  setStatusFilter(value?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.StatusFilter): void;

  hasSortOptions(): boolean;
  clearSortOptions(): void;
  getSortOptions(): github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.SortOptions | undefined;
  setSortOptions(value?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.SortOptions): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListProxiesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListProxiesRequest): ListProxiesRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListProxiesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListProxiesRequest;
  static deserializeBinaryFromReader(message: ListProxiesRequest, reader: jspb.BinaryReader): ListProxiesRequest;
}

export namespace ListProxiesRequest {
  export type AsObject = {
    glooInstanceRef?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef.AsObject,
    pagination?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.Pagination.AsObject,
    queryString: string,
    statusFilter?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.StatusFilter.AsObject,
    sortOptions?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.SortOptions.AsObject,
  }
}

export class ListProxiesResponse extends jspb.Message {
  clearProxiesList(): void;
  getProxiesList(): Array<Proxy>;
  setProxiesList(value: Array<Proxy>): void;
  addProxies(value?: Proxy, index?: number): Proxy;

  getTotal(): number;
  setTotal(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListProxiesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListProxiesResponse): ListProxiesResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListProxiesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListProxiesResponse;
  static deserializeBinaryFromReader(message: ListProxiesResponse, reader: jspb.BinaryReader): ListProxiesResponse;
}

export namespace ListProxiesResponse {
  export type AsObject = {
    proxiesList: Array<Proxy.AsObject>,
    total: number,
  }
}

export class GetProxyYamlRequest extends jspb.Message {
  hasProxyRef(): boolean;
  clearProxyRef(): void;
  getProxyRef(): github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef | undefined;
  setProxyRef(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetProxyYamlRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetProxyYamlRequest): GetProxyYamlRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetProxyYamlRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetProxyYamlRequest;
  static deserializeBinaryFromReader(message: GetProxyYamlRequest, reader: jspb.BinaryReader): GetProxyYamlRequest;
}

export namespace GetProxyYamlRequest {
  export type AsObject = {
    proxyRef?: github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef.AsObject,
  }
}

export class GetProxyYamlResponse extends jspb.Message {
  hasYamlData(): boolean;
  clearYamlData(): void;
  getYamlData(): github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ResourceYaml | undefined;
  setYamlData(value?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ResourceYaml): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetProxyYamlResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetProxyYamlResponse): GetProxyYamlResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetProxyYamlResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetProxyYamlResponse;
  static deserializeBinaryFromReader(message: GetProxyYamlResponse, reader: jspb.BinaryReader): GetProxyYamlResponse;
}

export namespace GetProxyYamlResponse {
  export type AsObject = {
    yamlData?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ResourceYaml.AsObject,
  }
}

export class GetProxyDetailsRequest extends jspb.Message {
  hasProxyRef(): boolean;
  clearProxyRef(): void;
  getProxyRef(): github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef | undefined;
  setProxyRef(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetProxyDetailsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetProxyDetailsRequest): GetProxyDetailsRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetProxyDetailsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetProxyDetailsRequest;
  static deserializeBinaryFromReader(message: GetProxyDetailsRequest, reader: jspb.BinaryReader): GetProxyDetailsRequest;
}

export namespace GetProxyDetailsRequest {
  export type AsObject = {
    proxyRef?: github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef.AsObject,
  }
}

export class GetProxyDetailsResponse extends jspb.Message {
  hasProxy(): boolean;
  clearProxy(): void;
  getProxy(): Proxy | undefined;
  setProxy(value?: Proxy): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetProxyDetailsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetProxyDetailsResponse): GetProxyDetailsResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetProxyDetailsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetProxyDetailsResponse;
  static deserializeBinaryFromReader(message: GetProxyDetailsResponse, reader: jspb.BinaryReader): GetProxyDetailsResponse;
}

export namespace GetProxyDetailsResponse {
  export type AsObject = {
    proxy?: Proxy.AsObject,
  }
}
