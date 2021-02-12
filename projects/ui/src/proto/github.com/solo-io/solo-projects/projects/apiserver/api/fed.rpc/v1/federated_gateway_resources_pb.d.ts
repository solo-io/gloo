/* eslint-disable */
// package: fed.rpc.solo.io
// file: github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/federated_gateway_resources.proto

import * as jspb from "google-protobuf";
import * as extproto_ext_pb from "../../../../../../../../extproto/ext_pb";
import * as github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_gateway_v1_gateway_pb from "../../../../../../../../github.com/solo-io/solo-projects/projects/gloo-fed/api/fed.gateway/v1/gateway_pb";
import * as github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_gateway_v1_virtual_service_pb from "../../../../../../../../github.com/solo-io/solo-projects/projects/gloo-fed/api/fed.gateway/v1/virtual_service_pb";
import * as github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_gateway_v1_route_table_pb from "../../../../../../../../github.com/solo-io/solo-projects/projects/gloo-fed/api/fed.gateway/v1/route_table_pb";
import * as github_com_solo_io_skv2_api_core_v1_core_pb from "../../../../../../../../github.com/solo-io/skv2/api/core/v1/core_pb";
import * as github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb from "../../../../../../../../github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/common_pb";

export class FederatedGateway extends jspb.Message {
  hasMetadata(): boolean;
  clearMetadata(): void;
  getMetadata(): github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb.ObjectMeta | undefined;
  setMetadata(value?: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb.ObjectMeta): void;

  hasSpec(): boolean;
  clearSpec(): void;
  getSpec(): github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_gateway_v1_gateway_pb.FederatedGatewaySpec | undefined;
  setSpec(value?: github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_gateway_v1_gateway_pb.FederatedGatewaySpec): void;

  hasStatus(): boolean;
  clearStatus(): void;
  getStatus(): github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_gateway_v1_gateway_pb.FederatedGatewayStatus | undefined;
  setStatus(value?: github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_gateway_v1_gateway_pb.FederatedGatewayStatus): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FederatedGateway.AsObject;
  static toObject(includeInstance: boolean, msg: FederatedGateway): FederatedGateway.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: FederatedGateway, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FederatedGateway;
  static deserializeBinaryFromReader(message: FederatedGateway, reader: jspb.BinaryReader): FederatedGateway;
}

export namespace FederatedGateway {
  export type AsObject = {
    metadata?: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb.ObjectMeta.AsObject,
    spec?: github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_gateway_v1_gateway_pb.FederatedGatewaySpec.AsObject,
    status?: github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_gateway_v1_gateway_pb.FederatedGatewayStatus.AsObject,
  }
}

export class FederatedVirtualService extends jspb.Message {
  hasMetadata(): boolean;
  clearMetadata(): void;
  getMetadata(): github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb.ObjectMeta | undefined;
  setMetadata(value?: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb.ObjectMeta): void;

  hasSpec(): boolean;
  clearSpec(): void;
  getSpec(): github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_gateway_v1_virtual_service_pb.FederatedVirtualServiceSpec | undefined;
  setSpec(value?: github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_gateway_v1_virtual_service_pb.FederatedVirtualServiceSpec): void;

  hasStatus(): boolean;
  clearStatus(): void;
  getStatus(): github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_gateway_v1_virtual_service_pb.FederatedVirtualServiceStatus | undefined;
  setStatus(value?: github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_gateway_v1_virtual_service_pb.FederatedVirtualServiceStatus): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FederatedVirtualService.AsObject;
  static toObject(includeInstance: boolean, msg: FederatedVirtualService): FederatedVirtualService.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: FederatedVirtualService, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FederatedVirtualService;
  static deserializeBinaryFromReader(message: FederatedVirtualService, reader: jspb.BinaryReader): FederatedVirtualService;
}

export namespace FederatedVirtualService {
  export type AsObject = {
    metadata?: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb.ObjectMeta.AsObject,
    spec?: github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_gateway_v1_virtual_service_pb.FederatedVirtualServiceSpec.AsObject,
    status?: github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_gateway_v1_virtual_service_pb.FederatedVirtualServiceStatus.AsObject,
  }
}

export class FederatedRouteTable extends jspb.Message {
  hasMetadata(): boolean;
  clearMetadata(): void;
  getMetadata(): github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb.ObjectMeta | undefined;
  setMetadata(value?: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb.ObjectMeta): void;

  hasSpec(): boolean;
  clearSpec(): void;
  getSpec(): github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_gateway_v1_route_table_pb.FederatedRouteTableSpec | undefined;
  setSpec(value?: github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_gateway_v1_route_table_pb.FederatedRouteTableSpec): void;

  hasStatus(): boolean;
  clearStatus(): void;
  getStatus(): github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_gateway_v1_route_table_pb.FederatedRouteTableStatus | undefined;
  setStatus(value?: github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_gateway_v1_route_table_pb.FederatedRouteTableStatus): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FederatedRouteTable.AsObject;
  static toObject(includeInstance: boolean, msg: FederatedRouteTable): FederatedRouteTable.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: FederatedRouteTable, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FederatedRouteTable;
  static deserializeBinaryFromReader(message: FederatedRouteTable, reader: jspb.BinaryReader): FederatedRouteTable;
}

export namespace FederatedRouteTable {
  export type AsObject = {
    metadata?: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb.ObjectMeta.AsObject,
    spec?: github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_gateway_v1_route_table_pb.FederatedRouteTableSpec.AsObject,
    status?: github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_gateway_v1_route_table_pb.FederatedRouteTableStatus.AsObject,
  }
}

export class ListFederatedGatewaysRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListFederatedGatewaysRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListFederatedGatewaysRequest): ListFederatedGatewaysRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListFederatedGatewaysRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListFederatedGatewaysRequest;
  static deserializeBinaryFromReader(message: ListFederatedGatewaysRequest, reader: jspb.BinaryReader): ListFederatedGatewaysRequest;
}

export namespace ListFederatedGatewaysRequest {
  export type AsObject = {
  }
}

export class ListFederatedGatewaysResponse extends jspb.Message {
  clearFederatedGatewaysList(): void;
  getFederatedGatewaysList(): Array<FederatedGateway>;
  setFederatedGatewaysList(value: Array<FederatedGateway>): void;
  addFederatedGateways(value?: FederatedGateway, index?: number): FederatedGateway;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListFederatedGatewaysResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListFederatedGatewaysResponse): ListFederatedGatewaysResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListFederatedGatewaysResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListFederatedGatewaysResponse;
  static deserializeBinaryFromReader(message: ListFederatedGatewaysResponse, reader: jspb.BinaryReader): ListFederatedGatewaysResponse;
}

export namespace ListFederatedGatewaysResponse {
  export type AsObject = {
    federatedGatewaysList: Array<FederatedGateway.AsObject>,
  }
}

export class GetFederatedGatewayYamlRequest extends jspb.Message {
  hasFederatedGatewayRef(): boolean;
  clearFederatedGatewayRef(): void;
  getFederatedGatewayRef(): github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef | undefined;
  setFederatedGatewayRef(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetFederatedGatewayYamlRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetFederatedGatewayYamlRequest): GetFederatedGatewayYamlRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetFederatedGatewayYamlRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetFederatedGatewayYamlRequest;
  static deserializeBinaryFromReader(message: GetFederatedGatewayYamlRequest, reader: jspb.BinaryReader): GetFederatedGatewayYamlRequest;
}

export namespace GetFederatedGatewayYamlRequest {
  export type AsObject = {
    federatedGatewayRef?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef.AsObject,
  }
}

export class GetFederatedGatewayYamlResponse extends jspb.Message {
  hasYamlData(): boolean;
  clearYamlData(): void;
  getYamlData(): github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb.ResourceYaml | undefined;
  setYamlData(value?: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb.ResourceYaml): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetFederatedGatewayYamlResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetFederatedGatewayYamlResponse): GetFederatedGatewayYamlResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetFederatedGatewayYamlResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetFederatedGatewayYamlResponse;
  static deserializeBinaryFromReader(message: GetFederatedGatewayYamlResponse, reader: jspb.BinaryReader): GetFederatedGatewayYamlResponse;
}

export namespace GetFederatedGatewayYamlResponse {
  export type AsObject = {
    yamlData?: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb.ResourceYaml.AsObject,
  }
}

export class ListFederatedVirtualServicesRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListFederatedVirtualServicesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListFederatedVirtualServicesRequest): ListFederatedVirtualServicesRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListFederatedVirtualServicesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListFederatedVirtualServicesRequest;
  static deserializeBinaryFromReader(message: ListFederatedVirtualServicesRequest, reader: jspb.BinaryReader): ListFederatedVirtualServicesRequest;
}

export namespace ListFederatedVirtualServicesRequest {
  export type AsObject = {
  }
}

export class ListFederatedVirtualServicesResponse extends jspb.Message {
  clearFederatedVirtualServicesList(): void;
  getFederatedVirtualServicesList(): Array<FederatedVirtualService>;
  setFederatedVirtualServicesList(value: Array<FederatedVirtualService>): void;
  addFederatedVirtualServices(value?: FederatedVirtualService, index?: number): FederatedVirtualService;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListFederatedVirtualServicesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListFederatedVirtualServicesResponse): ListFederatedVirtualServicesResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListFederatedVirtualServicesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListFederatedVirtualServicesResponse;
  static deserializeBinaryFromReader(message: ListFederatedVirtualServicesResponse, reader: jspb.BinaryReader): ListFederatedVirtualServicesResponse;
}

export namespace ListFederatedVirtualServicesResponse {
  export type AsObject = {
    federatedVirtualServicesList: Array<FederatedVirtualService.AsObject>,
  }
}

export class GetFederatedVirtualServiceYamlRequest extends jspb.Message {
  hasFederatedVirtualServiceRef(): boolean;
  clearFederatedVirtualServiceRef(): void;
  getFederatedVirtualServiceRef(): github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef | undefined;
  setFederatedVirtualServiceRef(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetFederatedVirtualServiceYamlRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetFederatedVirtualServiceYamlRequest): GetFederatedVirtualServiceYamlRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetFederatedVirtualServiceYamlRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetFederatedVirtualServiceYamlRequest;
  static deserializeBinaryFromReader(message: GetFederatedVirtualServiceYamlRequest, reader: jspb.BinaryReader): GetFederatedVirtualServiceYamlRequest;
}

export namespace GetFederatedVirtualServiceYamlRequest {
  export type AsObject = {
    federatedVirtualServiceRef?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef.AsObject,
  }
}

export class GetFederatedVirtualServiceYamlResponse extends jspb.Message {
  hasYamlData(): boolean;
  clearYamlData(): void;
  getYamlData(): github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb.ResourceYaml | undefined;
  setYamlData(value?: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb.ResourceYaml): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetFederatedVirtualServiceYamlResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetFederatedVirtualServiceYamlResponse): GetFederatedVirtualServiceYamlResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetFederatedVirtualServiceYamlResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetFederatedVirtualServiceYamlResponse;
  static deserializeBinaryFromReader(message: GetFederatedVirtualServiceYamlResponse, reader: jspb.BinaryReader): GetFederatedVirtualServiceYamlResponse;
}

export namespace GetFederatedVirtualServiceYamlResponse {
  export type AsObject = {
    yamlData?: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb.ResourceYaml.AsObject,
  }
}

export class ListFederatedRouteTablesRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListFederatedRouteTablesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListFederatedRouteTablesRequest): ListFederatedRouteTablesRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListFederatedRouteTablesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListFederatedRouteTablesRequest;
  static deserializeBinaryFromReader(message: ListFederatedRouteTablesRequest, reader: jspb.BinaryReader): ListFederatedRouteTablesRequest;
}

export namespace ListFederatedRouteTablesRequest {
  export type AsObject = {
  }
}

export class ListFederatedRouteTablesResponse extends jspb.Message {
  clearFederatedRouteTablesList(): void;
  getFederatedRouteTablesList(): Array<FederatedRouteTable>;
  setFederatedRouteTablesList(value: Array<FederatedRouteTable>): void;
  addFederatedRouteTables(value?: FederatedRouteTable, index?: number): FederatedRouteTable;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListFederatedRouteTablesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListFederatedRouteTablesResponse): ListFederatedRouteTablesResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListFederatedRouteTablesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListFederatedRouteTablesResponse;
  static deserializeBinaryFromReader(message: ListFederatedRouteTablesResponse, reader: jspb.BinaryReader): ListFederatedRouteTablesResponse;
}

export namespace ListFederatedRouteTablesResponse {
  export type AsObject = {
    federatedRouteTablesList: Array<FederatedRouteTable.AsObject>,
  }
}

export class GetFederatedRouteTableYamlRequest extends jspb.Message {
  hasFederatedRouteTableRef(): boolean;
  clearFederatedRouteTableRef(): void;
  getFederatedRouteTableRef(): github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef | undefined;
  setFederatedRouteTableRef(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetFederatedRouteTableYamlRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetFederatedRouteTableYamlRequest): GetFederatedRouteTableYamlRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetFederatedRouteTableYamlRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetFederatedRouteTableYamlRequest;
  static deserializeBinaryFromReader(message: GetFederatedRouteTableYamlRequest, reader: jspb.BinaryReader): GetFederatedRouteTableYamlRequest;
}

export namespace GetFederatedRouteTableYamlRequest {
  export type AsObject = {
    federatedRouteTableRef?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef.AsObject,
  }
}

export class GetFederatedRouteTableYamlResponse extends jspb.Message {
  hasYamlData(): boolean;
  clearYamlData(): void;
  getYamlData(): github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb.ResourceYaml | undefined;
  setYamlData(value?: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb.ResourceYaml): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetFederatedRouteTableYamlResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetFederatedRouteTableYamlResponse): GetFederatedRouteTableYamlResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetFederatedRouteTableYamlResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetFederatedRouteTableYamlResponse;
  static deserializeBinaryFromReader(message: GetFederatedRouteTableYamlResponse, reader: jspb.BinaryReader): GetFederatedRouteTableYamlResponse;
}

export namespace GetFederatedRouteTableYamlResponse {
  export type AsObject = {
    yamlData?: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb.ResourceYaml.AsObject,
  }
}
