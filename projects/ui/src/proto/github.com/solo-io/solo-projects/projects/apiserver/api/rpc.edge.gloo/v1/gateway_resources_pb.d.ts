/* eslint-disable */
// package: rpc.edge.gloo.solo.io
// file: github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/gateway_resources.proto

import * as jspb from "google-protobuf";
import * as extproto_ext_pb from "../../../../../../../../extproto/ext_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gateway_v1_gateway_pb from "../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gateway/v1/gateway_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gateway_v1_matchable_http_gateway_pb from "../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gateway/v1/matchable_http_gateway_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gateway_v1_virtual_service_pb from "../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gateway/v1/virtual_service_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gateway_v1_route_table_pb from "../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gateway/v1/route_table_pb";
import * as github_com_solo_io_skv2_api_core_v1_core_pb from "../../../../../../../../github.com/solo-io/skv2/api/core/v1/core_pb";
import * as github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb from "../../../../../../../../github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/common_pb";

export class Gateway extends jspb.Message {
  hasMetadata(): boolean;
  clearMetadata(): void;
  getMetadata(): github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ObjectMeta | undefined;
  setMetadata(value?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ObjectMeta): void;

  hasSpec(): boolean;
  clearSpec(): void;
  getSpec(): github_com_solo_io_solo_apis_api_gloo_gateway_v1_gateway_pb.GatewaySpec | undefined;
  setSpec(value?: github_com_solo_io_solo_apis_api_gloo_gateway_v1_gateway_pb.GatewaySpec): void;

  hasStatus(): boolean;
  clearStatus(): void;
  getStatus(): github_com_solo_io_solo_apis_api_gloo_gateway_v1_gateway_pb.GatewayStatus | undefined;
  setStatus(value?: github_com_solo_io_solo_apis_api_gloo_gateway_v1_gateway_pb.GatewayStatus): void;

  hasGlooInstance(): boolean;
  clearGlooInstance(): void;
  getGlooInstance(): github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef | undefined;
  setGlooInstance(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Gateway.AsObject;
  static toObject(includeInstance: boolean, msg: Gateway): Gateway.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Gateway, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Gateway;
  static deserializeBinaryFromReader(message: Gateway, reader: jspb.BinaryReader): Gateway;
}

export namespace Gateway {
  export type AsObject = {
    metadata?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ObjectMeta.AsObject,
    spec?: github_com_solo_io_solo_apis_api_gloo_gateway_v1_gateway_pb.GatewaySpec.AsObject,
    status?: github_com_solo_io_solo_apis_api_gloo_gateway_v1_gateway_pb.GatewayStatus.AsObject,
    glooInstance?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef.AsObject,
  }
}

export class MatchableHttpGateway extends jspb.Message {
  hasMetadata(): boolean;
  clearMetadata(): void;
  getMetadata(): github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ObjectMeta | undefined;
  setMetadata(value?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ObjectMeta): void;

  hasSpec(): boolean;
  clearSpec(): void;
  getSpec(): github_com_solo_io_solo_apis_api_gloo_gateway_v1_matchable_http_gateway_pb.MatchableHttpGatewaySpec | undefined;
  setSpec(value?: github_com_solo_io_solo_apis_api_gloo_gateway_v1_matchable_http_gateway_pb.MatchableHttpGatewaySpec): void;

  hasStatus(): boolean;
  clearStatus(): void;
  getStatus(): github_com_solo_io_solo_apis_api_gloo_gateway_v1_matchable_http_gateway_pb.MatchableHttpGatewayStatus | undefined;
  setStatus(value?: github_com_solo_io_solo_apis_api_gloo_gateway_v1_matchable_http_gateway_pb.MatchableHttpGatewayStatus): void;

  hasGlooInstance(): boolean;
  clearGlooInstance(): void;
  getGlooInstance(): github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef | undefined;
  setGlooInstance(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MatchableHttpGateway.AsObject;
  static toObject(includeInstance: boolean, msg: MatchableHttpGateway): MatchableHttpGateway.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: MatchableHttpGateway, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MatchableHttpGateway;
  static deserializeBinaryFromReader(message: MatchableHttpGateway, reader: jspb.BinaryReader): MatchableHttpGateway;
}

export namespace MatchableHttpGateway {
  export type AsObject = {
    metadata?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ObjectMeta.AsObject,
    spec?: github_com_solo_io_solo_apis_api_gloo_gateway_v1_matchable_http_gateway_pb.MatchableHttpGatewaySpec.AsObject,
    status?: github_com_solo_io_solo_apis_api_gloo_gateway_v1_matchable_http_gateway_pb.MatchableHttpGatewayStatus.AsObject,
    glooInstance?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef.AsObject,
  }
}

export class VirtualService extends jspb.Message {
  hasMetadata(): boolean;
  clearMetadata(): void;
  getMetadata(): github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ObjectMeta | undefined;
  setMetadata(value?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ObjectMeta): void;

  hasSpec(): boolean;
  clearSpec(): void;
  getSpec(): github_com_solo_io_solo_apis_api_gloo_gateway_v1_virtual_service_pb.VirtualServiceSpec | undefined;
  setSpec(value?: github_com_solo_io_solo_apis_api_gloo_gateway_v1_virtual_service_pb.VirtualServiceSpec): void;

  hasStatus(): boolean;
  clearStatus(): void;
  getStatus(): github_com_solo_io_solo_apis_api_gloo_gateway_v1_virtual_service_pb.VirtualServiceStatus | undefined;
  setStatus(value?: github_com_solo_io_solo_apis_api_gloo_gateway_v1_virtual_service_pb.VirtualServiceStatus): void;

  hasGlooInstance(): boolean;
  clearGlooInstance(): void;
  getGlooInstance(): github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef | undefined;
  setGlooInstance(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VirtualService.AsObject;
  static toObject(includeInstance: boolean, msg: VirtualService): VirtualService.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: VirtualService, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VirtualService;
  static deserializeBinaryFromReader(message: VirtualService, reader: jspb.BinaryReader): VirtualService;
}

export namespace VirtualService {
  export type AsObject = {
    metadata?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ObjectMeta.AsObject,
    spec?: github_com_solo_io_solo_apis_api_gloo_gateway_v1_virtual_service_pb.VirtualServiceSpec.AsObject,
    status?: github_com_solo_io_solo_apis_api_gloo_gateway_v1_virtual_service_pb.VirtualServiceStatus.AsObject,
    glooInstance?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef.AsObject,
  }
}

export class RouteTable extends jspb.Message {
  hasMetadata(): boolean;
  clearMetadata(): void;
  getMetadata(): github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ObjectMeta | undefined;
  setMetadata(value?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ObjectMeta): void;

  hasSpec(): boolean;
  clearSpec(): void;
  getSpec(): github_com_solo_io_solo_apis_api_gloo_gateway_v1_route_table_pb.RouteTableSpec | undefined;
  setSpec(value?: github_com_solo_io_solo_apis_api_gloo_gateway_v1_route_table_pb.RouteTableSpec): void;

  hasStatus(): boolean;
  clearStatus(): void;
  getStatus(): github_com_solo_io_solo_apis_api_gloo_gateway_v1_route_table_pb.RouteTableStatus | undefined;
  setStatus(value?: github_com_solo_io_solo_apis_api_gloo_gateway_v1_route_table_pb.RouteTableStatus): void;

  hasGlooInstance(): boolean;
  clearGlooInstance(): void;
  getGlooInstance(): github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef | undefined;
  setGlooInstance(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RouteTable.AsObject;
  static toObject(includeInstance: boolean, msg: RouteTable): RouteTable.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RouteTable, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RouteTable;
  static deserializeBinaryFromReader(message: RouteTable, reader: jspb.BinaryReader): RouteTable;
}

export namespace RouteTable {
  export type AsObject = {
    metadata?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ObjectMeta.AsObject,
    spec?: github_com_solo_io_solo_apis_api_gloo_gateway_v1_route_table_pb.RouteTableSpec.AsObject,
    status?: github_com_solo_io_solo_apis_api_gloo_gateway_v1_route_table_pb.RouteTableStatus.AsObject,
    glooInstance?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef.AsObject,
  }
}

export class ListGatewaysRequest extends jspb.Message {
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
  toObject(includeInstance?: boolean): ListGatewaysRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListGatewaysRequest): ListGatewaysRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListGatewaysRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListGatewaysRequest;
  static deserializeBinaryFromReader(message: ListGatewaysRequest, reader: jspb.BinaryReader): ListGatewaysRequest;
}

export namespace ListGatewaysRequest {
  export type AsObject = {
    glooInstanceRef?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef.AsObject,
    pagination?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.Pagination.AsObject,
    queryString: string,
    statusFilter?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.StatusFilter.AsObject,
    sortOptions?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.SortOptions.AsObject,
  }
}

export class ListGatewaysResponse extends jspb.Message {
  clearGatewaysList(): void;
  getGatewaysList(): Array<Gateway>;
  setGatewaysList(value: Array<Gateway>): void;
  addGateways(value?: Gateway, index?: number): Gateway;

  getTotal(): number;
  setTotal(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListGatewaysResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListGatewaysResponse): ListGatewaysResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListGatewaysResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListGatewaysResponse;
  static deserializeBinaryFromReader(message: ListGatewaysResponse, reader: jspb.BinaryReader): ListGatewaysResponse;
}

export namespace ListGatewaysResponse {
  export type AsObject = {
    gatewaysList: Array<Gateway.AsObject>,
    total: number,
  }
}

export class GetGatewayYamlRequest extends jspb.Message {
  hasGatewayRef(): boolean;
  clearGatewayRef(): void;
  getGatewayRef(): github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef | undefined;
  setGatewayRef(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetGatewayYamlRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetGatewayYamlRequest): GetGatewayYamlRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetGatewayYamlRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetGatewayYamlRequest;
  static deserializeBinaryFromReader(message: GetGatewayYamlRequest, reader: jspb.BinaryReader): GetGatewayYamlRequest;
}

export namespace GetGatewayYamlRequest {
  export type AsObject = {
    gatewayRef?: github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef.AsObject,
  }
}

export class GetGatewayYamlResponse extends jspb.Message {
  hasYamlData(): boolean;
  clearYamlData(): void;
  getYamlData(): github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ResourceYaml | undefined;
  setYamlData(value?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ResourceYaml): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetGatewayYamlResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetGatewayYamlResponse): GetGatewayYamlResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetGatewayYamlResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetGatewayYamlResponse;
  static deserializeBinaryFromReader(message: GetGatewayYamlResponse, reader: jspb.BinaryReader): GetGatewayYamlResponse;
}

export namespace GetGatewayYamlResponse {
  export type AsObject = {
    yamlData?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ResourceYaml.AsObject,
  }
}

export class ListMatchableHttpGatewaysRequest extends jspb.Message {
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
  toObject(includeInstance?: boolean): ListMatchableHttpGatewaysRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListMatchableHttpGatewaysRequest): ListMatchableHttpGatewaysRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListMatchableHttpGatewaysRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListMatchableHttpGatewaysRequest;
  static deserializeBinaryFromReader(message: ListMatchableHttpGatewaysRequest, reader: jspb.BinaryReader): ListMatchableHttpGatewaysRequest;
}

export namespace ListMatchableHttpGatewaysRequest {
  export type AsObject = {
    glooInstanceRef?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef.AsObject,
    pagination?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.Pagination.AsObject,
    queryString: string,
    statusFilter?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.StatusFilter.AsObject,
    sortOptions?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.SortOptions.AsObject,
  }
}

export class ListMatchableHttpGatewaysResponse extends jspb.Message {
  clearMatchableHttpGatewaysList(): void;
  getMatchableHttpGatewaysList(): Array<MatchableHttpGateway>;
  setMatchableHttpGatewaysList(value: Array<MatchableHttpGateway>): void;
  addMatchableHttpGateways(value?: MatchableHttpGateway, index?: number): MatchableHttpGateway;

  getTotal(): number;
  setTotal(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListMatchableHttpGatewaysResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListMatchableHttpGatewaysResponse): ListMatchableHttpGatewaysResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListMatchableHttpGatewaysResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListMatchableHttpGatewaysResponse;
  static deserializeBinaryFromReader(message: ListMatchableHttpGatewaysResponse, reader: jspb.BinaryReader): ListMatchableHttpGatewaysResponse;
}

export namespace ListMatchableHttpGatewaysResponse {
  export type AsObject = {
    matchableHttpGatewaysList: Array<MatchableHttpGateway.AsObject>,
    total: number,
  }
}

export class GetMatchableHttpGatewayYamlRequest extends jspb.Message {
  hasMatchableHttpGatewayRef(): boolean;
  clearMatchableHttpGatewayRef(): void;
  getMatchableHttpGatewayRef(): github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef | undefined;
  setMatchableHttpGatewayRef(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetMatchableHttpGatewayYamlRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetMatchableHttpGatewayYamlRequest): GetMatchableHttpGatewayYamlRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetMatchableHttpGatewayYamlRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetMatchableHttpGatewayYamlRequest;
  static deserializeBinaryFromReader(message: GetMatchableHttpGatewayYamlRequest, reader: jspb.BinaryReader): GetMatchableHttpGatewayYamlRequest;
}

export namespace GetMatchableHttpGatewayYamlRequest {
  export type AsObject = {
    matchableHttpGatewayRef?: github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef.AsObject,
  }
}

export class GetMatchableHttpGatewayYamlResponse extends jspb.Message {
  hasYamlData(): boolean;
  clearYamlData(): void;
  getYamlData(): github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ResourceYaml | undefined;
  setYamlData(value?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ResourceYaml): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetMatchableHttpGatewayYamlResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetMatchableHttpGatewayYamlResponse): GetMatchableHttpGatewayYamlResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetMatchableHttpGatewayYamlResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetMatchableHttpGatewayYamlResponse;
  static deserializeBinaryFromReader(message: GetMatchableHttpGatewayYamlResponse, reader: jspb.BinaryReader): GetMatchableHttpGatewayYamlResponse;
}

export namespace GetMatchableHttpGatewayYamlResponse {
  export type AsObject = {
    yamlData?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ResourceYaml.AsObject,
  }
}

export class ListVirtualServicesRequest extends jspb.Message {
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
  toObject(includeInstance?: boolean): ListVirtualServicesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListVirtualServicesRequest): ListVirtualServicesRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListVirtualServicesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListVirtualServicesRequest;
  static deserializeBinaryFromReader(message: ListVirtualServicesRequest, reader: jspb.BinaryReader): ListVirtualServicesRequest;
}

export namespace ListVirtualServicesRequest {
  export type AsObject = {
    glooInstanceRef?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef.AsObject,
    pagination?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.Pagination.AsObject,
    queryString: string,
    statusFilter?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.StatusFilter.AsObject,
    sortOptions?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.SortOptions.AsObject,
  }
}

export class ListVirtualServicesResponse extends jspb.Message {
  clearVirtualServicesList(): void;
  getVirtualServicesList(): Array<VirtualService>;
  setVirtualServicesList(value: Array<VirtualService>): void;
  addVirtualServices(value?: VirtualService, index?: number): VirtualService;

  getTotal(): number;
  setTotal(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListVirtualServicesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListVirtualServicesResponse): ListVirtualServicesResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListVirtualServicesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListVirtualServicesResponse;
  static deserializeBinaryFromReader(message: ListVirtualServicesResponse, reader: jspb.BinaryReader): ListVirtualServicesResponse;
}

export namespace ListVirtualServicesResponse {
  export type AsObject = {
    virtualServicesList: Array<VirtualService.AsObject>,
    total: number,
  }
}

export class GetVirtualServiceYamlRequest extends jspb.Message {
  hasVirtualServiceRef(): boolean;
  clearVirtualServiceRef(): void;
  getVirtualServiceRef(): github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef | undefined;
  setVirtualServiceRef(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetVirtualServiceYamlRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetVirtualServiceYamlRequest): GetVirtualServiceYamlRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetVirtualServiceYamlRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetVirtualServiceYamlRequest;
  static deserializeBinaryFromReader(message: GetVirtualServiceYamlRequest, reader: jspb.BinaryReader): GetVirtualServiceYamlRequest;
}

export namespace GetVirtualServiceYamlRequest {
  export type AsObject = {
    virtualServiceRef?: github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef.AsObject,
  }
}

export class GetVirtualServiceYamlResponse extends jspb.Message {
  hasYamlData(): boolean;
  clearYamlData(): void;
  getYamlData(): github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ResourceYaml | undefined;
  setYamlData(value?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ResourceYaml): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetVirtualServiceYamlResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetVirtualServiceYamlResponse): GetVirtualServiceYamlResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetVirtualServiceYamlResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetVirtualServiceYamlResponse;
  static deserializeBinaryFromReader(message: GetVirtualServiceYamlResponse, reader: jspb.BinaryReader): GetVirtualServiceYamlResponse;
}

export namespace GetVirtualServiceYamlResponse {
  export type AsObject = {
    yamlData?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ResourceYaml.AsObject,
  }
}

export class ListRouteTablesRequest extends jspb.Message {
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
  toObject(includeInstance?: boolean): ListRouteTablesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListRouteTablesRequest): ListRouteTablesRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListRouteTablesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListRouteTablesRequest;
  static deserializeBinaryFromReader(message: ListRouteTablesRequest, reader: jspb.BinaryReader): ListRouteTablesRequest;
}

export namespace ListRouteTablesRequest {
  export type AsObject = {
    glooInstanceRef?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef.AsObject,
    pagination?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.Pagination.AsObject,
    queryString: string,
    statusFilter?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.StatusFilter.AsObject,
    sortOptions?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.SortOptions.AsObject,
  }
}

export class ListRouteTablesResponse extends jspb.Message {
  clearRouteTablesList(): void;
  getRouteTablesList(): Array<RouteTable>;
  setRouteTablesList(value: Array<RouteTable>): void;
  addRouteTables(value?: RouteTable, index?: number): RouteTable;

  getTotal(): number;
  setTotal(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListRouteTablesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListRouteTablesResponse): ListRouteTablesResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListRouteTablesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListRouteTablesResponse;
  static deserializeBinaryFromReader(message: ListRouteTablesResponse, reader: jspb.BinaryReader): ListRouteTablesResponse;
}

export namespace ListRouteTablesResponse {
  export type AsObject = {
    routeTablesList: Array<RouteTable.AsObject>,
    total: number,
  }
}

export class GetRouteTableYamlRequest extends jspb.Message {
  hasRouteTableRef(): boolean;
  clearRouteTableRef(): void;
  getRouteTableRef(): github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef | undefined;
  setRouteTableRef(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetRouteTableYamlRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetRouteTableYamlRequest): GetRouteTableYamlRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetRouteTableYamlRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetRouteTableYamlRequest;
  static deserializeBinaryFromReader(message: GetRouteTableYamlRequest, reader: jspb.BinaryReader): GetRouteTableYamlRequest;
}

export namespace GetRouteTableYamlRequest {
  export type AsObject = {
    routeTableRef?: github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef.AsObject,
  }
}

export class GetRouteTableYamlResponse extends jspb.Message {
  hasYamlData(): boolean;
  clearYamlData(): void;
  getYamlData(): github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ResourceYaml | undefined;
  setYamlData(value?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ResourceYaml): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetRouteTableYamlResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetRouteTableYamlResponse): GetRouteTableYamlResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetRouteTableYamlResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetRouteTableYamlResponse;
  static deserializeBinaryFromReader(message: GetRouteTableYamlResponse, reader: jspb.BinaryReader): GetRouteTableYamlResponse;
}

export namespace GetRouteTableYamlResponse {
  export type AsObject = {
    yamlData?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ResourceYaml.AsObject,
  }
}
