// package: glooeeapi.solo.io
// file: solo-projects/projects/grpcserver/api/v1/upstreamgroup.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../../../gogoproto/gogo_pb";
import * as gloo_projects_gloo_api_v1_proxy_pb from "../../../../../gloo/projects/gloo/api/v1/proxy_pb";
import * as solo_projects_projects_grpcserver_api_v1_types_pb from "../../../../../solo-projects/projects/grpcserver/api/v1/types_pb";
import * as solo_kit_api_v1_ref_pb from "../../../../../solo-kit/api/v1/ref_pb";

export class UpstreamGroupDetails extends jspb.Message {
  hasUpstreamGroup(): boolean;
  clearUpstreamGroup(): void;
  getUpstreamGroup(): gloo_projects_gloo_api_v1_proxy_pb.UpstreamGroup | undefined;
  setUpstreamGroup(value?: gloo_projects_gloo_api_v1_proxy_pb.UpstreamGroup): void;

  hasRaw(): boolean;
  clearRaw(): void;
  getRaw(): solo_projects_projects_grpcserver_api_v1_types_pb.Raw | undefined;
  setRaw(value?: solo_projects_projects_grpcserver_api_v1_types_pb.Raw): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpstreamGroupDetails.AsObject;
  static toObject(includeInstance: boolean, msg: UpstreamGroupDetails): UpstreamGroupDetails.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpstreamGroupDetails, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpstreamGroupDetails;
  static deserializeBinaryFromReader(message: UpstreamGroupDetails, reader: jspb.BinaryReader): UpstreamGroupDetails;
}

export namespace UpstreamGroupDetails {
  export type AsObject = {
    upstreamGroup?: gloo_projects_gloo_api_v1_proxy_pb.UpstreamGroup.AsObject,
    raw?: solo_projects_projects_grpcserver_api_v1_types_pb.Raw.AsObject,
  }
}

export class GetUpstreamGroupRequest extends jspb.Message {
  hasRef(): boolean;
  clearRef(): void;
  getRef(): solo_kit_api_v1_ref_pb.ResourceRef | undefined;
  setRef(value?: solo_kit_api_v1_ref_pb.ResourceRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetUpstreamGroupRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetUpstreamGroupRequest): GetUpstreamGroupRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetUpstreamGroupRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetUpstreamGroupRequest;
  static deserializeBinaryFromReader(message: GetUpstreamGroupRequest, reader: jspb.BinaryReader): GetUpstreamGroupRequest;
}

export namespace GetUpstreamGroupRequest {
  export type AsObject = {
    ref?: solo_kit_api_v1_ref_pb.ResourceRef.AsObject,
  }
}

export class GetUpstreamGroupResponse extends jspb.Message {
  hasUpstreamGroupDetails(): boolean;
  clearUpstreamGroupDetails(): void;
  getUpstreamGroupDetails(): UpstreamGroupDetails | undefined;
  setUpstreamGroupDetails(value?: UpstreamGroupDetails): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetUpstreamGroupResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetUpstreamGroupResponse): GetUpstreamGroupResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetUpstreamGroupResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetUpstreamGroupResponse;
  static deserializeBinaryFromReader(message: GetUpstreamGroupResponse, reader: jspb.BinaryReader): GetUpstreamGroupResponse;
}

export namespace GetUpstreamGroupResponse {
  export type AsObject = {
    upstreamGroupDetails?: UpstreamGroupDetails.AsObject,
  }
}

export class ListUpstreamGroupsRequest extends jspb.Message {
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
  }
}

export class ListUpstreamGroupsResponse extends jspb.Message {
  clearUpstreamGroupDetailsList(): void;
  getUpstreamGroupDetailsList(): Array<UpstreamGroupDetails>;
  setUpstreamGroupDetailsList(value: Array<UpstreamGroupDetails>): void;
  addUpstreamGroupDetails(value?: UpstreamGroupDetails, index?: number): UpstreamGroupDetails;

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
    upstreamGroupDetailsList: Array<UpstreamGroupDetails.AsObject>,
  }
}

export class CreateUpstreamGroupRequest extends jspb.Message {
  hasUpstreamGroup(): boolean;
  clearUpstreamGroup(): void;
  getUpstreamGroup(): gloo_projects_gloo_api_v1_proxy_pb.UpstreamGroup | undefined;
  setUpstreamGroup(value?: gloo_projects_gloo_api_v1_proxy_pb.UpstreamGroup): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateUpstreamGroupRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateUpstreamGroupRequest): CreateUpstreamGroupRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateUpstreamGroupRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateUpstreamGroupRequest;
  static deserializeBinaryFromReader(message: CreateUpstreamGroupRequest, reader: jspb.BinaryReader): CreateUpstreamGroupRequest;
}

export namespace CreateUpstreamGroupRequest {
  export type AsObject = {
    upstreamGroup?: gloo_projects_gloo_api_v1_proxy_pb.UpstreamGroup.AsObject,
  }
}

export class CreateUpstreamGroupResponse extends jspb.Message {
  hasUpstreamGroupDetails(): boolean;
  clearUpstreamGroupDetails(): void;
  getUpstreamGroupDetails(): UpstreamGroupDetails | undefined;
  setUpstreamGroupDetails(value?: UpstreamGroupDetails): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateUpstreamGroupResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CreateUpstreamGroupResponse): CreateUpstreamGroupResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateUpstreamGroupResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateUpstreamGroupResponse;
  static deserializeBinaryFromReader(message: CreateUpstreamGroupResponse, reader: jspb.BinaryReader): CreateUpstreamGroupResponse;
}

export namespace CreateUpstreamGroupResponse {
  export type AsObject = {
    upstreamGroupDetails?: UpstreamGroupDetails.AsObject,
  }
}

export class UpdateUpstreamGroupRequest extends jspb.Message {
  hasUpstreamGroup(): boolean;
  clearUpstreamGroup(): void;
  getUpstreamGroup(): gloo_projects_gloo_api_v1_proxy_pb.UpstreamGroup | undefined;
  setUpstreamGroup(value?: gloo_projects_gloo_api_v1_proxy_pb.UpstreamGroup): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateUpstreamGroupRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateUpstreamGroupRequest): UpdateUpstreamGroupRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateUpstreamGroupRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateUpstreamGroupRequest;
  static deserializeBinaryFromReader(message: UpdateUpstreamGroupRequest, reader: jspb.BinaryReader): UpdateUpstreamGroupRequest;
}

export namespace UpdateUpstreamGroupRequest {
  export type AsObject = {
    upstreamGroup?: gloo_projects_gloo_api_v1_proxy_pb.UpstreamGroup.AsObject,
  }
}

export class UpdateUpstreamGroupYamlRequest extends jspb.Message {
  hasEditedYamlData(): boolean;
  clearEditedYamlData(): void;
  getEditedYamlData(): solo_projects_projects_grpcserver_api_v1_types_pb.EditedResourceYaml | undefined;
  setEditedYamlData(value?: solo_projects_projects_grpcserver_api_v1_types_pb.EditedResourceYaml): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateUpstreamGroupYamlRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateUpstreamGroupYamlRequest): UpdateUpstreamGroupYamlRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateUpstreamGroupYamlRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateUpstreamGroupYamlRequest;
  static deserializeBinaryFromReader(message: UpdateUpstreamGroupYamlRequest, reader: jspb.BinaryReader): UpdateUpstreamGroupYamlRequest;
}

export namespace UpdateUpstreamGroupYamlRequest {
  export type AsObject = {
    editedYamlData?: solo_projects_projects_grpcserver_api_v1_types_pb.EditedResourceYaml.AsObject,
  }
}

export class UpdateUpstreamGroupResponse extends jspb.Message {
  hasUpstreamGroupDetails(): boolean;
  clearUpstreamGroupDetails(): void;
  getUpstreamGroupDetails(): UpstreamGroupDetails | undefined;
  setUpstreamGroupDetails(value?: UpstreamGroupDetails): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateUpstreamGroupResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateUpstreamGroupResponse): UpdateUpstreamGroupResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateUpstreamGroupResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateUpstreamGroupResponse;
  static deserializeBinaryFromReader(message: UpdateUpstreamGroupResponse, reader: jspb.BinaryReader): UpdateUpstreamGroupResponse;
}

export namespace UpdateUpstreamGroupResponse {
  export type AsObject = {
    upstreamGroupDetails?: UpstreamGroupDetails.AsObject,
  }
}

export class DeleteUpstreamGroupRequest extends jspb.Message {
  hasRef(): boolean;
  clearRef(): void;
  getRef(): solo_kit_api_v1_ref_pb.ResourceRef | undefined;
  setRef(value?: solo_kit_api_v1_ref_pb.ResourceRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteUpstreamGroupRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteUpstreamGroupRequest): DeleteUpstreamGroupRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteUpstreamGroupRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteUpstreamGroupRequest;
  static deserializeBinaryFromReader(message: DeleteUpstreamGroupRequest, reader: jspb.BinaryReader): DeleteUpstreamGroupRequest;
}

export namespace DeleteUpstreamGroupRequest {
  export type AsObject = {
    ref?: solo_kit_api_v1_ref_pb.ResourceRef.AsObject,
  }
}

export class DeleteUpstreamGroupResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteUpstreamGroupResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteUpstreamGroupResponse): DeleteUpstreamGroupResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteUpstreamGroupResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteUpstreamGroupResponse;
  static deserializeBinaryFromReader(message: DeleteUpstreamGroupResponse, reader: jspb.BinaryReader): DeleteUpstreamGroupResponse;
}

export namespace DeleteUpstreamGroupResponse {
  export type AsObject = {
  }
}

