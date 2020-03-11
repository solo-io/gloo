/* eslint-disable */
// package: apiserver.devportal.solo.io
// file: dev-portal/api/grpc/apiserver/portal.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../../gogoproto/gogo_pb";
import * as dev_portal_api_dev_portal_v1_portal_pb from "../../../../dev-portal/api/dev-portal/v1/portal_pb";
import * as dev_portal_api_grpc_common_common_pb from "../../../../dev-portal/api/grpc/common/common_pb";

export class Portal extends jspb.Message {
  hasMetadata(): boolean;
  clearMetadata(): void;
  getMetadata(): dev_portal_api_grpc_common_common_pb.ObjectMeta | undefined;
  setMetadata(value?: dev_portal_api_grpc_common_common_pb.ObjectMeta): void;

  hasSpec(): boolean;
  clearSpec(): void;
  getSpec(): dev_portal_api_dev_portal_v1_portal_pb.PortalSpec | undefined;
  setSpec(value?: dev_portal_api_dev_portal_v1_portal_pb.PortalSpec): void;

  hasStatus(): boolean;
  clearStatus(): void;
  getStatus(): dev_portal_api_dev_portal_v1_portal_pb.PortalStatus | undefined;
  setStatus(value?: dev_portal_api_dev_portal_v1_portal_pb.PortalStatus): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Portal.AsObject;
  static toObject(includeInstance: boolean, msg: Portal): Portal.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Portal, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Portal;
  static deserializeBinaryFromReader(message: Portal, reader: jspb.BinaryReader): Portal;
}

export namespace Portal {
  export type AsObject = {
    metadata?: dev_portal_api_grpc_common_common_pb.ObjectMeta.AsObject,
    spec?: dev_portal_api_dev_portal_v1_portal_pb.PortalSpec.AsObject,
    status?: dev_portal_api_dev_portal_v1_portal_pb.PortalStatus.AsObject,
  }
}

export class GetPortalRequest extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getNamespace(): string;
  setNamespace(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetPortalRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetPortalRequest): GetPortalRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetPortalRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetPortalRequest;
  static deserializeBinaryFromReader(message: GetPortalRequest, reader: jspb.BinaryReader): GetPortalRequest;
}

export namespace GetPortalRequest {
  export type AsObject = {
    name: string,
    namespace: string,
  }
}

export class GetPortalResponse extends jspb.Message {
  hasPortal(): boolean;
  clearPortal(): void;
  getPortal(): Portal | undefined;
  setPortal(value?: Portal): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetPortalResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetPortalResponse): GetPortalResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetPortalResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetPortalResponse;
  static deserializeBinaryFromReader(message: GetPortalResponse, reader: jspb.BinaryReader): GetPortalResponse;
}

export namespace GetPortalResponse {
  export type AsObject = {
    portal?: Portal.AsObject,
  }
}

export class ListPortalsRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListPortalsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListPortalsRequest): ListPortalsRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListPortalsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListPortalsRequest;
  static deserializeBinaryFromReader(message: ListPortalsRequest, reader: jspb.BinaryReader): ListPortalsRequest;
}

export namespace ListPortalsRequest {
  export type AsObject = {
  }
}

export class ListPortalsResponse extends jspb.Message {
  clearPortalsList(): void;
  getPortalsList(): Array<Portal>;
  setPortalsList(value: Array<Portal>): void;
  addPortals(value?: Portal, index?: number): Portal;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListPortalsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListPortalsResponse): ListPortalsResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListPortalsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListPortalsResponse;
  static deserializeBinaryFromReader(message: ListPortalsResponse, reader: jspb.BinaryReader): ListPortalsResponse;
}

export namespace ListPortalsResponse {
  export type AsObject = {
    portalsList: Array<Portal.AsObject>,
  }
}

export class CreatePortalRequest extends jspb.Message {
  hasPortal(): boolean;
  clearPortal(): void;
  getPortal(): Portal | undefined;
  setPortal(value?: Portal): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreatePortalRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreatePortalRequest): CreatePortalRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreatePortalRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreatePortalRequest;
  static deserializeBinaryFromReader(message: CreatePortalRequest, reader: jspb.BinaryReader): CreatePortalRequest;
}

export namespace CreatePortalRequest {
  export type AsObject = {
    portal?: Portal.AsObject,
  }
}

export class CreatePortalResponse extends jspb.Message {
  hasPortal(): boolean;
  clearPortal(): void;
  getPortal(): Portal | undefined;
  setPortal(value?: Portal): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreatePortalResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CreatePortalResponse): CreatePortalResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreatePortalResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreatePortalResponse;
  static deserializeBinaryFromReader(message: CreatePortalResponse, reader: jspb.BinaryReader): CreatePortalResponse;
}

export namespace CreatePortalResponse {
  export type AsObject = {
    portal?: Portal.AsObject,
  }
}

export class UpdatePortalRequest extends jspb.Message {
  hasPortal(): boolean;
  clearPortal(): void;
  getPortal(): Portal | undefined;
  setPortal(value?: Portal): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdatePortalRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdatePortalRequest): UpdatePortalRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdatePortalRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdatePortalRequest;
  static deserializeBinaryFromReader(message: UpdatePortalRequest, reader: jspb.BinaryReader): UpdatePortalRequest;
}

export namespace UpdatePortalRequest {
  export type AsObject = {
    portal?: Portal.AsObject,
  }
}

export class UpdatePortalResponse extends jspb.Message {
  hasPortal(): boolean;
  clearPortal(): void;
  getPortal(): Portal | undefined;
  setPortal(value?: Portal): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdatePortalResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdatePortalResponse): UpdatePortalResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdatePortalResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdatePortalResponse;
  static deserializeBinaryFromReader(message: UpdatePortalResponse, reader: jspb.BinaryReader): UpdatePortalResponse;
}

export namespace UpdatePortalResponse {
  export type AsObject = {
    portal?: Portal.AsObject,
  }
}

export class DeletePortalRequest extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getNamespace(): string;
  setNamespace(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeletePortalRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeletePortalRequest): DeletePortalRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeletePortalRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeletePortalRequest;
  static deserializeBinaryFromReader(message: DeletePortalRequest, reader: jspb.BinaryReader): DeletePortalRequest;
}

export namespace DeletePortalRequest {
  export type AsObject = {
    name: string,
    namespace: string,
  }
}

export class DeletePortalResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeletePortalResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeletePortalResponse): DeletePortalResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeletePortalResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeletePortalResponse;
  static deserializeBinaryFromReader(message: DeletePortalResponse, reader: jspb.BinaryReader): DeletePortalResponse;
}

export namespace DeletePortalResponse {
  export type AsObject = {
  }
}
