// package: glooeeapi.solo.io
// file: solo-projects/projects/grpcserver/api/v1/gateway.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../../../gogoproto/gogo_pb";
import * as gloo_projects_gateway_api_v1_gateway_pb from "../../../../../gloo/projects/gateway/api/v1/gateway_pb";
import * as solo_kit_api_v1_ref_pb from "../../../../../solo-kit/api/v1/ref_pb";
import * as solo_projects_projects_grpcserver_api_v1_types_pb from "../../../../../solo-projects/projects/grpcserver/api/v1/types_pb";

export class GatewayDetails extends jspb.Message {
  hasGateway(): boolean;
  clearGateway(): void;
  getGateway(): gloo_projects_gateway_api_v1_gateway_pb.Gateway | undefined;
  setGateway(value?: gloo_projects_gateway_api_v1_gateway_pb.Gateway): void;

  hasRaw(): boolean;
  clearRaw(): void;
  getRaw(): solo_projects_projects_grpcserver_api_v1_types_pb.Raw | undefined;
  setRaw(value?: solo_projects_projects_grpcserver_api_v1_types_pb.Raw): void;

  hasStatus(): boolean;
  clearStatus(): void;
  getStatus(): solo_projects_projects_grpcserver_api_v1_types_pb.Status | undefined;
  setStatus(value?: solo_projects_projects_grpcserver_api_v1_types_pb.Status): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GatewayDetails.AsObject;
  static toObject(includeInstance: boolean, msg: GatewayDetails): GatewayDetails.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GatewayDetails, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GatewayDetails;
  static deserializeBinaryFromReader(message: GatewayDetails, reader: jspb.BinaryReader): GatewayDetails;
}

export namespace GatewayDetails {
  export type AsObject = {
    gateway?: gloo_projects_gateway_api_v1_gateway_pb.Gateway.AsObject,
    raw?: solo_projects_projects_grpcserver_api_v1_types_pb.Raw.AsObject,
    status?: solo_projects_projects_grpcserver_api_v1_types_pb.Status.AsObject,
  }
}

export class GetGatewayRequest extends jspb.Message {
  hasRef(): boolean;
  clearRef(): void;
  getRef(): solo_kit_api_v1_ref_pb.ResourceRef | undefined;
  setRef(value?: solo_kit_api_v1_ref_pb.ResourceRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetGatewayRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetGatewayRequest): GetGatewayRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetGatewayRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetGatewayRequest;
  static deserializeBinaryFromReader(message: GetGatewayRequest, reader: jspb.BinaryReader): GetGatewayRequest;
}

export namespace GetGatewayRequest {
  export type AsObject = {
    ref?: solo_kit_api_v1_ref_pb.ResourceRef.AsObject,
  }
}

export class GetGatewayResponse extends jspb.Message {
  hasGatewayDetails(): boolean;
  clearGatewayDetails(): void;
  getGatewayDetails(): GatewayDetails | undefined;
  setGatewayDetails(value?: GatewayDetails): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetGatewayResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetGatewayResponse): GetGatewayResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetGatewayResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetGatewayResponse;
  static deserializeBinaryFromReader(message: GetGatewayResponse, reader: jspb.BinaryReader): GetGatewayResponse;
}

export namespace GetGatewayResponse {
  export type AsObject = {
    gatewayDetails?: GatewayDetails.AsObject,
  }
}

export class ListGatewaysRequest extends jspb.Message {
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
  }
}

export class ListGatewaysResponse extends jspb.Message {
  clearGatewayDetailsList(): void;
  getGatewayDetailsList(): Array<GatewayDetails>;
  setGatewayDetailsList(value: Array<GatewayDetails>): void;
  addGatewayDetails(value?: GatewayDetails, index?: number): GatewayDetails;

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
    gatewayDetailsList: Array<GatewayDetails.AsObject>,
  }
}

export class UpdateGatewayRequest extends jspb.Message {
  hasGateway(): boolean;
  clearGateway(): void;
  getGateway(): gloo_projects_gateway_api_v1_gateway_pb.Gateway | undefined;
  setGateway(value?: gloo_projects_gateway_api_v1_gateway_pb.Gateway): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateGatewayRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateGatewayRequest): UpdateGatewayRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateGatewayRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateGatewayRequest;
  static deserializeBinaryFromReader(message: UpdateGatewayRequest, reader: jspb.BinaryReader): UpdateGatewayRequest;
}

export namespace UpdateGatewayRequest {
  export type AsObject = {
    gateway?: gloo_projects_gateway_api_v1_gateway_pb.Gateway.AsObject,
  }
}

export class UpdateGatewayYamlRequest extends jspb.Message {
  hasEditedYamlData(): boolean;
  clearEditedYamlData(): void;
  getEditedYamlData(): solo_projects_projects_grpcserver_api_v1_types_pb.EditedResourceYaml | undefined;
  setEditedYamlData(value?: solo_projects_projects_grpcserver_api_v1_types_pb.EditedResourceYaml): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateGatewayYamlRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateGatewayYamlRequest): UpdateGatewayYamlRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateGatewayYamlRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateGatewayYamlRequest;
  static deserializeBinaryFromReader(message: UpdateGatewayYamlRequest, reader: jspb.BinaryReader): UpdateGatewayYamlRequest;
}

export namespace UpdateGatewayYamlRequest {
  export type AsObject = {
    editedYamlData?: solo_projects_projects_grpcserver_api_v1_types_pb.EditedResourceYaml.AsObject,
  }
}

export class UpdateGatewayResponse extends jspb.Message {
  hasGatewayDetails(): boolean;
  clearGatewayDetails(): void;
  getGatewayDetails(): GatewayDetails | undefined;
  setGatewayDetails(value?: GatewayDetails): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateGatewayResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateGatewayResponse): UpdateGatewayResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateGatewayResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateGatewayResponse;
  static deserializeBinaryFromReader(message: UpdateGatewayResponse, reader: jspb.BinaryReader): UpdateGatewayResponse;
}

export namespace UpdateGatewayResponse {
  export type AsObject = {
    gatewayDetails?: GatewayDetails.AsObject,
  }
}

