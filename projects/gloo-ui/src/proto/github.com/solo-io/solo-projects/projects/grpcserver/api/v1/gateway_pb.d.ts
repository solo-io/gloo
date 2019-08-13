// package: glooeeapi.solo.io
// file: github.com/solo-io/solo-projects/projects/grpcserver/api/v1/gateway.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../../../../../gogoproto/gogo_pb";
import * as github_com_solo_io_gloo_projects_gateway_api_v2_gateway_pb from "../../../../../../../github.com/solo-io/gloo/projects/gateway/api/v2/gateway_pb";
import * as github_com_solo_io_solo_kit_api_v1_ref_pb from "../../../../../../../github.com/solo-io/solo-kit/api/v1/ref_pb";
import * as github_com_solo_io_solo_projects_projects_grpcserver_api_v1_types_pb from "../../../../../../../github.com/solo-io/solo-projects/projects/grpcserver/api/v1/types_pb";

export class GatewayDetails extends jspb.Message {
  hasGateway(): boolean;
  clearGateway(): void;
  getGateway(): github_com_solo_io_gloo_projects_gateway_api_v2_gateway_pb.Gateway | undefined;
  setGateway(value?: github_com_solo_io_gloo_projects_gateway_api_v2_gateway_pb.Gateway): void;

  hasRaw(): boolean;
  clearRaw(): void;
  getRaw(): github_com_solo_io_solo_projects_projects_grpcserver_api_v1_types_pb.Raw | undefined;
  setRaw(value?: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_types_pb.Raw): void;

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
    gateway?: github_com_solo_io_gloo_projects_gateway_api_v2_gateway_pb.Gateway.AsObject,
    raw?: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_types_pb.Raw.AsObject,
  }
}

export class GetGatewayRequest extends jspb.Message {
  hasRef(): boolean;
  clearRef(): void;
  getRef(): github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef | undefined;
  setRef(value?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef): void;

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
    ref?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject,
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
  clearNamespacesList(): void;
  getNamespacesList(): Array<string>;
  setNamespacesList(value: Array<string>): void;
  addNamespaces(value: string, index?: number): string;

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
    namespacesList: Array<string>,
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

