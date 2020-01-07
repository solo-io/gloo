// package: glooeeapi.solo.io
// file: solo-projects/projects/grpcserver/api/v1/artifact.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../../../gogoproto/gogo_pb";
import * as gloo_projects_gloo_api_v1_artifact_pb from "../../../../../gloo/projects/gloo/api/v1/artifact_pb";
import * as solo_kit_api_v1_ref_pb from "../../../../../solo-kit/api/v1/ref_pb";

export class GetArtifactRequest extends jspb.Message {
  hasRef(): boolean;
  clearRef(): void;
  getRef(): solo_kit_api_v1_ref_pb.ResourceRef | undefined;
  setRef(value?: solo_kit_api_v1_ref_pb.ResourceRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetArtifactRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetArtifactRequest): GetArtifactRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetArtifactRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetArtifactRequest;
  static deserializeBinaryFromReader(message: GetArtifactRequest, reader: jspb.BinaryReader): GetArtifactRequest;
}

export namespace GetArtifactRequest {
  export type AsObject = {
    ref?: solo_kit_api_v1_ref_pb.ResourceRef.AsObject,
  }
}

export class GetArtifactResponse extends jspb.Message {
  hasArtifact(): boolean;
  clearArtifact(): void;
  getArtifact(): gloo_projects_gloo_api_v1_artifact_pb.Artifact | undefined;
  setArtifact(value?: gloo_projects_gloo_api_v1_artifact_pb.Artifact): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetArtifactResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetArtifactResponse): GetArtifactResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetArtifactResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetArtifactResponse;
  static deserializeBinaryFromReader(message: GetArtifactResponse, reader: jspb.BinaryReader): GetArtifactResponse;
}

export namespace GetArtifactResponse {
  export type AsObject = {
    artifact?: gloo_projects_gloo_api_v1_artifact_pb.Artifact.AsObject,
  }
}

export class ListArtifactsRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListArtifactsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListArtifactsRequest): ListArtifactsRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListArtifactsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListArtifactsRequest;
  static deserializeBinaryFromReader(message: ListArtifactsRequest, reader: jspb.BinaryReader): ListArtifactsRequest;
}

export namespace ListArtifactsRequest {
  export type AsObject = {
  }
}

export class ListArtifactsResponse extends jspb.Message {
  clearArtifactsList(): void;
  getArtifactsList(): Array<gloo_projects_gloo_api_v1_artifact_pb.Artifact>;
  setArtifactsList(value: Array<gloo_projects_gloo_api_v1_artifact_pb.Artifact>): void;
  addArtifacts(value?: gloo_projects_gloo_api_v1_artifact_pb.Artifact, index?: number): gloo_projects_gloo_api_v1_artifact_pb.Artifact;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListArtifactsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListArtifactsResponse): ListArtifactsResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListArtifactsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListArtifactsResponse;
  static deserializeBinaryFromReader(message: ListArtifactsResponse, reader: jspb.BinaryReader): ListArtifactsResponse;
}

export namespace ListArtifactsResponse {
  export type AsObject = {
    artifactsList: Array<gloo_projects_gloo_api_v1_artifact_pb.Artifact.AsObject>,
  }
}

export class CreateArtifactRequest extends jspb.Message {
  hasArtifact(): boolean;
  clearArtifact(): void;
  getArtifact(): gloo_projects_gloo_api_v1_artifact_pb.Artifact | undefined;
  setArtifact(value?: gloo_projects_gloo_api_v1_artifact_pb.Artifact): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateArtifactRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateArtifactRequest): CreateArtifactRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateArtifactRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateArtifactRequest;
  static deserializeBinaryFromReader(message: CreateArtifactRequest, reader: jspb.BinaryReader): CreateArtifactRequest;
}

export namespace CreateArtifactRequest {
  export type AsObject = {
    artifact?: gloo_projects_gloo_api_v1_artifact_pb.Artifact.AsObject,
  }
}

export class CreateArtifactResponse extends jspb.Message {
  hasArtifact(): boolean;
  clearArtifact(): void;
  getArtifact(): gloo_projects_gloo_api_v1_artifact_pb.Artifact | undefined;
  setArtifact(value?: gloo_projects_gloo_api_v1_artifact_pb.Artifact): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateArtifactResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CreateArtifactResponse): CreateArtifactResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateArtifactResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateArtifactResponse;
  static deserializeBinaryFromReader(message: CreateArtifactResponse, reader: jspb.BinaryReader): CreateArtifactResponse;
}

export namespace CreateArtifactResponse {
  export type AsObject = {
    artifact?: gloo_projects_gloo_api_v1_artifact_pb.Artifact.AsObject,
  }
}

export class UpdateArtifactRequest extends jspb.Message {
  hasArtifact(): boolean;
  clearArtifact(): void;
  getArtifact(): gloo_projects_gloo_api_v1_artifact_pb.Artifact | undefined;
  setArtifact(value?: gloo_projects_gloo_api_v1_artifact_pb.Artifact): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateArtifactRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateArtifactRequest): UpdateArtifactRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateArtifactRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateArtifactRequest;
  static deserializeBinaryFromReader(message: UpdateArtifactRequest, reader: jspb.BinaryReader): UpdateArtifactRequest;
}

export namespace UpdateArtifactRequest {
  export type AsObject = {
    artifact?: gloo_projects_gloo_api_v1_artifact_pb.Artifact.AsObject,
  }
}

export class UpdateArtifactResponse extends jspb.Message {
  hasArtifact(): boolean;
  clearArtifact(): void;
  getArtifact(): gloo_projects_gloo_api_v1_artifact_pb.Artifact | undefined;
  setArtifact(value?: gloo_projects_gloo_api_v1_artifact_pb.Artifact): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateArtifactResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateArtifactResponse): UpdateArtifactResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateArtifactResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateArtifactResponse;
  static deserializeBinaryFromReader(message: UpdateArtifactResponse, reader: jspb.BinaryReader): UpdateArtifactResponse;
}

export namespace UpdateArtifactResponse {
  export type AsObject = {
    artifact?: gloo_projects_gloo_api_v1_artifact_pb.Artifact.AsObject,
  }
}

export class DeleteArtifactRequest extends jspb.Message {
  hasRef(): boolean;
  clearRef(): void;
  getRef(): solo_kit_api_v1_ref_pb.ResourceRef | undefined;
  setRef(value?: solo_kit_api_v1_ref_pb.ResourceRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteArtifactRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteArtifactRequest): DeleteArtifactRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteArtifactRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteArtifactRequest;
  static deserializeBinaryFromReader(message: DeleteArtifactRequest, reader: jspb.BinaryReader): DeleteArtifactRequest;
}

export namespace DeleteArtifactRequest {
  export type AsObject = {
    ref?: solo_kit_api_v1_ref_pb.ResourceRef.AsObject,
  }
}

export class DeleteArtifactResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteArtifactResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteArtifactResponse): DeleteArtifactResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteArtifactResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteArtifactResponse;
  static deserializeBinaryFromReader(message: DeleteArtifactResponse, reader: jspb.BinaryReader): DeleteArtifactResponse;
}

export namespace DeleteArtifactResponse {
  export type AsObject = {
  }
}

