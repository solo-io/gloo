/* eslint-disable */
// package: fed.rpc.solo.io
// file: github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/gateway_resources.proto

import * as jspb from "google-protobuf";
import * as extproto_ext_pb from "../../../../../../../../extproto/ext_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gateway_v1_gateway_pb from "../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gateway/v1/gateway_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gateway_v1_virtual_service_pb from "../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gateway/v1/virtual_service_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gateway_v1_route_table_pb from "../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gateway/v1/route_table_pb";
import * as github_com_solo_io_skv2_api_core_v1_core_pb from "../../../../../../../../github.com/solo-io/skv2/api/core/v1/core_pb";
import * as github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb from "../../../../../../../../github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/common_pb";

export class Gateway extends jspb.Message {
  hasMetadata(): boolean;
  clearMetadata(): void;
  getMetadata(): github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb.ObjectMeta | undefined;
  setMetadata(value?: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb.ObjectMeta): void;

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
    metadata?: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb.ObjectMeta.AsObject,
    spec?: github_com_solo_io_solo_apis_api_gloo_gateway_v1_gateway_pb.GatewaySpec.AsObject,
    status?: github_com_solo_io_solo_apis_api_gloo_gateway_v1_gateway_pb.GatewayStatus.AsObject,
    glooInstance?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef.AsObject,
  }
}

export class VirtualService extends jspb.Message {
  hasMetadata(): boolean;
  clearMetadata(): void;
  getMetadata(): github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb.ObjectMeta | undefined;
  setMetadata(value?: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb.ObjectMeta): void;

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
    metadata?: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb.ObjectMeta.AsObject,
    spec?: github_com_solo_io_solo_apis_api_gloo_gateway_v1_virtual_service_pb.VirtualServiceSpec.AsObject,
    status?: github_com_solo_io_solo_apis_api_gloo_gateway_v1_virtual_service_pb.VirtualServiceStatus.AsObject,
    glooInstance?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef.AsObject,
  }
}

export class RouteTable extends jspb.Message {
  hasMetadata(): boolean;
  clearMetadata(): void;
  getMetadata(): github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb.ObjectMeta | undefined;
  setMetadata(value?: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb.ObjectMeta): void;

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
    metadata?: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb.ObjectMeta.AsObject,
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
  }
}

export class ListGatewaysResponse extends jspb.Message {
  clearGatewaysList(): void;
  getGatewaysList(): Array<Gateway>;
  setGatewaysList(value: Array<Gateway>): void;
  addGateways(value?: Gateway, index?: number): Gateway;

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
  getYamlData(): github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb.ResourceYaml | undefined;
  setYamlData(value?: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb.ResourceYaml): void;

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
    yamlData?: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb.ResourceYaml.AsObject,
  }
}

export class ListVirtualServicesRequest extends jspb.Message {
  hasGlooInstanceRef(): boolean;
  clearGlooInstanceRef(): void;
  getGlooInstanceRef(): github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef | undefined;
  setGlooInstanceRef(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef): void;

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
  }
}

export class ListVirtualServicesResponse extends jspb.Message {
  clearVirtualServicesList(): void;
  getVirtualServicesList(): Array<VirtualService>;
  setVirtualServicesList(value: Array<VirtualService>): void;
  addVirtualServices(value?: VirtualService, index?: number): VirtualService;

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
  getYamlData(): github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb.ResourceYaml | undefined;
  setYamlData(value?: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb.ResourceYaml): void;

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
    yamlData?: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb.ResourceYaml.AsObject,
  }
}

export class ListRouteTablesRequest extends jspb.Message {
  hasGlooInstanceRef(): boolean;
  clearGlooInstanceRef(): void;
  getGlooInstanceRef(): github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef | undefined;
  setGlooInstanceRef(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef): void;

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
  }
}

export class ListRouteTablesResponse extends jspb.Message {
  clearRouteTablesList(): void;
  getRouteTablesList(): Array<RouteTable>;
  setRouteTablesList(value: Array<RouteTable>): void;
  addRouteTables(value?: RouteTable, index?: number): RouteTable;

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
  getYamlData(): github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb.ResourceYaml | undefined;
  setYamlData(value?: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb.ResourceYaml): void;

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
    yamlData?: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb.ResourceYaml.AsObject,
  }
}
