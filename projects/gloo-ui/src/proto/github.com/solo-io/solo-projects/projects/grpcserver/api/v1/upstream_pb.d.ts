// package: glooeeapi.solo.io
// file: github.com/solo-io/solo-projects/projects/grpcserver/api/v1/upstream.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../../../../../gogoproto/gogo_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_v1_upstream_pb from "../../../../../../../github.com/solo-io/gloo/projects/gloo/api/v1/upstream_pb";
import * as github_com_solo_io_solo_kit_api_v1_ref_pb from "../../../../../../../github.com/solo-io/solo-kit/api/v1/ref_pb";
import * as github_com_solo_io_solo_projects_projects_grpcserver_api_v1_types_pb from "../../../../../../../github.com/solo-io/solo-projects/projects/grpcserver/api/v1/types_pb";

export class UpstreamDetails extends jspb.Message {
  hasUpstream(): boolean;
  clearUpstream(): void;
  getUpstream(): github_com_solo_io_gloo_projects_gloo_api_v1_upstream_pb.Upstream | undefined;
  setUpstream(value?: github_com_solo_io_gloo_projects_gloo_api_v1_upstream_pb.Upstream): void;

  hasRaw(): boolean;
  clearRaw(): void;
  getRaw(): github_com_solo_io_solo_projects_projects_grpcserver_api_v1_types_pb.Raw | undefined;
  setRaw(value?: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_types_pb.Raw): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpstreamDetails.AsObject;
  static toObject(includeInstance: boolean, msg: UpstreamDetails): UpstreamDetails.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpstreamDetails, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpstreamDetails;
  static deserializeBinaryFromReader(message: UpstreamDetails, reader: jspb.BinaryReader): UpstreamDetails;
}

export namespace UpstreamDetails {
  export type AsObject = {
    upstream?: github_com_solo_io_gloo_projects_gloo_api_v1_upstream_pb.Upstream.AsObject,
    raw?: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_types_pb.Raw.AsObject,
  }
}

export class GetUpstreamRequest extends jspb.Message {
  hasRef(): boolean;
  clearRef(): void;
  getRef(): github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef | undefined;
  setRef(value?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetUpstreamRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetUpstreamRequest): GetUpstreamRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetUpstreamRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetUpstreamRequest;
  static deserializeBinaryFromReader(message: GetUpstreamRequest, reader: jspb.BinaryReader): GetUpstreamRequest;
}

export namespace GetUpstreamRequest {
  export type AsObject = {
    ref?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject,
  }
}

export class GetUpstreamResponse extends jspb.Message {
  hasUpstreamDetails(): boolean;
  clearUpstreamDetails(): void;
  getUpstreamDetails(): UpstreamDetails | undefined;
  setUpstreamDetails(value?: UpstreamDetails): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetUpstreamResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetUpstreamResponse): GetUpstreamResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetUpstreamResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetUpstreamResponse;
  static deserializeBinaryFromReader(message: GetUpstreamResponse, reader: jspb.BinaryReader): GetUpstreamResponse;
}

export namespace GetUpstreamResponse {
  export type AsObject = {
    upstreamDetails?: UpstreamDetails.AsObject,
  }
}

export class ListUpstreamsRequest extends jspb.Message {
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
  }
}

export class ListUpstreamsResponse extends jspb.Message {
  clearUpstreamDetailsList(): void;
  getUpstreamDetailsList(): Array<UpstreamDetails>;
  setUpstreamDetailsList(value: Array<UpstreamDetails>): void;
  addUpstreamDetails(value?: UpstreamDetails, index?: number): UpstreamDetails;

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
    upstreamDetailsList: Array<UpstreamDetails.AsObject>,
  }
}

export class CreateUpstreamRequest extends jspb.Message {
  hasUpstreamInput(): boolean;
  clearUpstreamInput(): void;
  getUpstreamInput(): github_com_solo_io_gloo_projects_gloo_api_v1_upstream_pb.Upstream | undefined;
  setUpstreamInput(value?: github_com_solo_io_gloo_projects_gloo_api_v1_upstream_pb.Upstream): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateUpstreamRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateUpstreamRequest): CreateUpstreamRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateUpstreamRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateUpstreamRequest;
  static deserializeBinaryFromReader(message: CreateUpstreamRequest, reader: jspb.BinaryReader): CreateUpstreamRequest;
}

export namespace CreateUpstreamRequest {
  export type AsObject = {
    upstreamInput?: github_com_solo_io_gloo_projects_gloo_api_v1_upstream_pb.Upstream.AsObject,
  }
}

export class CreateUpstreamResponse extends jspb.Message {
  hasUpstreamDetails(): boolean;
  clearUpstreamDetails(): void;
  getUpstreamDetails(): UpstreamDetails | undefined;
  setUpstreamDetails(value?: UpstreamDetails): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateUpstreamResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CreateUpstreamResponse): CreateUpstreamResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateUpstreamResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateUpstreamResponse;
  static deserializeBinaryFromReader(message: CreateUpstreamResponse, reader: jspb.BinaryReader): CreateUpstreamResponse;
}

export namespace CreateUpstreamResponse {
  export type AsObject = {
    upstreamDetails?: UpstreamDetails.AsObject,
  }
}

export class UpdateUpstreamRequest extends jspb.Message {
  hasUpstreamInput(): boolean;
  clearUpstreamInput(): void;
  getUpstreamInput(): github_com_solo_io_gloo_projects_gloo_api_v1_upstream_pb.Upstream | undefined;
  setUpstreamInput(value?: github_com_solo_io_gloo_projects_gloo_api_v1_upstream_pb.Upstream): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateUpstreamRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateUpstreamRequest): UpdateUpstreamRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateUpstreamRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateUpstreamRequest;
  static deserializeBinaryFromReader(message: UpdateUpstreamRequest, reader: jspb.BinaryReader): UpdateUpstreamRequest;
}

export namespace UpdateUpstreamRequest {
  export type AsObject = {
    upstreamInput?: github_com_solo_io_gloo_projects_gloo_api_v1_upstream_pb.Upstream.AsObject,
  }
}

export class UpdateUpstreamResponse extends jspb.Message {
  hasUpstreamDetails(): boolean;
  clearUpstreamDetails(): void;
  getUpstreamDetails(): UpstreamDetails | undefined;
  setUpstreamDetails(value?: UpstreamDetails): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateUpstreamResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateUpstreamResponse): UpdateUpstreamResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateUpstreamResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateUpstreamResponse;
  static deserializeBinaryFromReader(message: UpdateUpstreamResponse, reader: jspb.BinaryReader): UpdateUpstreamResponse;
}

export namespace UpdateUpstreamResponse {
  export type AsObject = {
    upstreamDetails?: UpstreamDetails.AsObject,
  }
}

export class DeleteUpstreamRequest extends jspb.Message {
  hasRef(): boolean;
  clearRef(): void;
  getRef(): github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef | undefined;
  setRef(value?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteUpstreamRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteUpstreamRequest): DeleteUpstreamRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteUpstreamRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteUpstreamRequest;
  static deserializeBinaryFromReader(message: DeleteUpstreamRequest, reader: jspb.BinaryReader): DeleteUpstreamRequest;
}

export namespace DeleteUpstreamRequest {
  export type AsObject = {
    ref?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject,
  }
}

export class DeleteUpstreamResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteUpstreamResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteUpstreamResponse): DeleteUpstreamResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteUpstreamResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteUpstreamResponse;
  static deserializeBinaryFromReader(message: DeleteUpstreamResponse, reader: jspb.BinaryReader): DeleteUpstreamResponse;
}

export namespace DeleteUpstreamResponse {
  export type AsObject = {
  }
}

