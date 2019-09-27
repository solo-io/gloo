// package: glooeeapi.solo.io
// file: github.com/solo-io/solo-projects/projects/grpcserver/api/v1/routetable.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../../../../../gogoproto/gogo_pb";
import * as github_com_solo_io_gloo_projects_gateway_api_v1_route_table_pb from "../../../../../../../github.com/solo-io/gloo/projects/gateway/api/v1/route_table_pb";
import * as github_com_solo_io_solo_projects_projects_grpcserver_api_v1_types_pb from "../../../../../../../github.com/solo-io/solo-projects/projects/grpcserver/api/v1/types_pb";
import * as github_com_solo_io_solo_kit_api_v1_ref_pb from "../../../../../../../github.com/solo-io/solo-kit/api/v1/ref_pb";

export class RouteTableDetails extends jspb.Message {
  hasRouteTable(): boolean;
  clearRouteTable(): void;
  getRouteTable(): github_com_solo_io_gloo_projects_gateway_api_v1_route_table_pb.RouteTable | undefined;
  setRouteTable(value?: github_com_solo_io_gloo_projects_gateway_api_v1_route_table_pb.RouteTable): void;

  hasRaw(): boolean;
  clearRaw(): void;
  getRaw(): github_com_solo_io_solo_projects_projects_grpcserver_api_v1_types_pb.Raw | undefined;
  setRaw(value?: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_types_pb.Raw): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RouteTableDetails.AsObject;
  static toObject(includeInstance: boolean, msg: RouteTableDetails): RouteTableDetails.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RouteTableDetails, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RouteTableDetails;
  static deserializeBinaryFromReader(message: RouteTableDetails, reader: jspb.BinaryReader): RouteTableDetails;
}

export namespace RouteTableDetails {
  export type AsObject = {
    routeTable?: github_com_solo_io_gloo_projects_gateway_api_v1_route_table_pb.RouteTable.AsObject,
    raw?: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_types_pb.Raw.AsObject,
  }
}

export class GetRouteTableRequest extends jspb.Message {
  hasRef(): boolean;
  clearRef(): void;
  getRef(): github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef | undefined;
  setRef(value?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetRouteTableRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetRouteTableRequest): GetRouteTableRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetRouteTableRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetRouteTableRequest;
  static deserializeBinaryFromReader(message: GetRouteTableRequest, reader: jspb.BinaryReader): GetRouteTableRequest;
}

export namespace GetRouteTableRequest {
  export type AsObject = {
    ref?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject,
  }
}

export class GetRouteTableResponse extends jspb.Message {
  hasRouteTableDetails(): boolean;
  clearRouteTableDetails(): void;
  getRouteTableDetails(): RouteTableDetails | undefined;
  setRouteTableDetails(value?: RouteTableDetails): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetRouteTableResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetRouteTableResponse): GetRouteTableResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetRouteTableResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetRouteTableResponse;
  static deserializeBinaryFromReader(message: GetRouteTableResponse, reader: jspb.BinaryReader): GetRouteTableResponse;
}

export namespace GetRouteTableResponse {
  export type AsObject = {
    routeTableDetails?: RouteTableDetails.AsObject,
  }
}

export class ListRouteTablesRequest extends jspb.Message {
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
  }
}

export class ListRouteTablesResponse extends jspb.Message {
  clearRouteTableDetailsList(): void;
  getRouteTableDetailsList(): Array<RouteTableDetails>;
  setRouteTableDetailsList(value: Array<RouteTableDetails>): void;
  addRouteTableDetails(value?: RouteTableDetails, index?: number): RouteTableDetails;

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
    routeTableDetailsList: Array<RouteTableDetails.AsObject>,
  }
}

export class CreateRouteTableRequest extends jspb.Message {
  hasRouteTable(): boolean;
  clearRouteTable(): void;
  getRouteTable(): github_com_solo_io_gloo_projects_gateway_api_v1_route_table_pb.RouteTable | undefined;
  setRouteTable(value?: github_com_solo_io_gloo_projects_gateway_api_v1_route_table_pb.RouteTable): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateRouteTableRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateRouteTableRequest): CreateRouteTableRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateRouteTableRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateRouteTableRequest;
  static deserializeBinaryFromReader(message: CreateRouteTableRequest, reader: jspb.BinaryReader): CreateRouteTableRequest;
}

export namespace CreateRouteTableRequest {
  export type AsObject = {
    routeTable?: github_com_solo_io_gloo_projects_gateway_api_v1_route_table_pb.RouteTable.AsObject,
  }
}

export class CreateRouteTableResponse extends jspb.Message {
  hasRouteTableDetails(): boolean;
  clearRouteTableDetails(): void;
  getRouteTableDetails(): RouteTableDetails | undefined;
  setRouteTableDetails(value?: RouteTableDetails): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateRouteTableResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CreateRouteTableResponse): CreateRouteTableResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateRouteTableResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateRouteTableResponse;
  static deserializeBinaryFromReader(message: CreateRouteTableResponse, reader: jspb.BinaryReader): CreateRouteTableResponse;
}

export namespace CreateRouteTableResponse {
  export type AsObject = {
    routeTableDetails?: RouteTableDetails.AsObject,
  }
}

export class UpdateRouteTableRequest extends jspb.Message {
  hasRouteTable(): boolean;
  clearRouteTable(): void;
  getRouteTable(): github_com_solo_io_gloo_projects_gateway_api_v1_route_table_pb.RouteTable | undefined;
  setRouteTable(value?: github_com_solo_io_gloo_projects_gateway_api_v1_route_table_pb.RouteTable): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateRouteTableRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateRouteTableRequest): UpdateRouteTableRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateRouteTableRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateRouteTableRequest;
  static deserializeBinaryFromReader(message: UpdateRouteTableRequest, reader: jspb.BinaryReader): UpdateRouteTableRequest;
}

export namespace UpdateRouteTableRequest {
  export type AsObject = {
    routeTable?: github_com_solo_io_gloo_projects_gateway_api_v1_route_table_pb.RouteTable.AsObject,
  }
}

export class UpdateRouteTableYamlRequest extends jspb.Message {
  hasEditedYamlData(): boolean;
  clearEditedYamlData(): void;
  getEditedYamlData(): github_com_solo_io_solo_projects_projects_grpcserver_api_v1_types_pb.EditedResourceYaml | undefined;
  setEditedYamlData(value?: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_types_pb.EditedResourceYaml): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateRouteTableYamlRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateRouteTableYamlRequest): UpdateRouteTableYamlRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateRouteTableYamlRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateRouteTableYamlRequest;
  static deserializeBinaryFromReader(message: UpdateRouteTableYamlRequest, reader: jspb.BinaryReader): UpdateRouteTableYamlRequest;
}

export namespace UpdateRouteTableYamlRequest {
  export type AsObject = {
    editedYamlData?: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_types_pb.EditedResourceYaml.AsObject,
  }
}

export class UpdateRouteTableResponse extends jspb.Message {
  hasRouteTableDetails(): boolean;
  clearRouteTableDetails(): void;
  getRouteTableDetails(): RouteTableDetails | undefined;
  setRouteTableDetails(value?: RouteTableDetails): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateRouteTableResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateRouteTableResponse): UpdateRouteTableResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateRouteTableResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateRouteTableResponse;
  static deserializeBinaryFromReader(message: UpdateRouteTableResponse, reader: jspb.BinaryReader): UpdateRouteTableResponse;
}

export namespace UpdateRouteTableResponse {
  export type AsObject = {
    routeTableDetails?: RouteTableDetails.AsObject,
  }
}

export class DeleteRouteTableRequest extends jspb.Message {
  hasRef(): boolean;
  clearRef(): void;
  getRef(): github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef | undefined;
  setRef(value?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteRouteTableRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteRouteTableRequest): DeleteRouteTableRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteRouteTableRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteRouteTableRequest;
  static deserializeBinaryFromReader(message: DeleteRouteTableRequest, reader: jspb.BinaryReader): DeleteRouteTableRequest;
}

export namespace DeleteRouteTableRequest {
  export type AsObject = {
    ref?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject,
  }
}

export class DeleteRouteTableResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteRouteTableResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteRouteTableResponse): DeleteRouteTableResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteRouteTableResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteRouteTableResponse;
  static deserializeBinaryFromReader(message: DeleteRouteTableResponse, reader: jspb.BinaryReader): DeleteRouteTableResponse;
}

export namespace DeleteRouteTableResponse {
  export type AsObject = {
  }
}

