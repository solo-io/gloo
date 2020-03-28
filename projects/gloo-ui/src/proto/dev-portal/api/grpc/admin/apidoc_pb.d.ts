/* eslint-disable */
// package: admin.devportal.solo.io
// file: dev-portal/api/grpc/admin/apidoc.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../../gogoproto/gogo_pb";
import * as google_protobuf_empty_pb from "google-protobuf/google/protobuf/empty_pb";
import * as dev_portal_api_dev_portal_v1_apidoc_pb from "../../../../dev-portal/api/dev-portal/v1/apidoc_pb";
import * as dev_portal_api_dev_portal_v1_common_pb from "../../../../dev-portal/api/dev-portal/v1/common_pb";
import * as dev_portal_api_grpc_common_common_pb from "../../../../dev-portal/api/grpc/common/common_pb";

export class ApiDoc extends jspb.Message {
  hasMetadata(): boolean;
  clearMetadata(): void;
  getMetadata(): dev_portal_api_grpc_common_common_pb.ObjectMeta | undefined;
  setMetadata(value?: dev_portal_api_grpc_common_common_pb.ObjectMeta): void;

  hasSpec(): boolean;
  clearSpec(): void;
  getSpec(): dev_portal_api_dev_portal_v1_apidoc_pb.ApiDocSpec | undefined;
  setSpec(value?: dev_portal_api_dev_portal_v1_apidoc_pb.ApiDocSpec): void;

  hasStatus(): boolean;
  clearStatus(): void;
  getStatus(): dev_portal_api_dev_portal_v1_apidoc_pb.ApiDocStatus | undefined;
  setStatus(value?: dev_portal_api_dev_portal_v1_apidoc_pb.ApiDocStatus): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ApiDoc.AsObject;
  static toObject(includeInstance: boolean, msg: ApiDoc): ApiDoc.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ApiDoc, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ApiDoc;
  static deserializeBinaryFromReader(message: ApiDoc, reader: jspb.BinaryReader): ApiDoc;
}

export namespace ApiDoc {
  export type AsObject = {
    metadata?: dev_portal_api_grpc_common_common_pb.ObjectMeta.AsObject,
    spec?: dev_portal_api_dev_portal_v1_apidoc_pb.ApiDocSpec.AsObject,
    status?: dev_portal_api_dev_portal_v1_apidoc_pb.ApiDocStatus.AsObject,
  }
}

export class ApiDocList extends jspb.Message {
  clearApidocsList(): void;
  getApidocsList(): Array<ApiDoc>;
  setApidocsList(value: Array<ApiDoc>): void;
  addApidocs(value?: ApiDoc, index?: number): ApiDoc;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ApiDocList.AsObject;
  static toObject(includeInstance: boolean, msg: ApiDocList): ApiDocList.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ApiDocList, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ApiDocList;
  static deserializeBinaryFromReader(message: ApiDocList, reader: jspb.BinaryReader): ApiDocList;
}

export namespace ApiDocList {
  export type AsObject = {
    apidocsList: Array<ApiDoc.AsObject>,
  }
}

export class ApiDocGetRequest extends jspb.Message {
  hasApidoc(): boolean;
  clearApidoc(): void;
  getApidoc(): dev_portal_api_dev_portal_v1_common_pb.ObjectRef | undefined;
  setApidoc(value?: dev_portal_api_dev_portal_v1_common_pb.ObjectRef): void;

  getWithassets(): boolean;
  setWithassets(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ApiDocGetRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ApiDocGetRequest): ApiDocGetRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ApiDocGetRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ApiDocGetRequest;
  static deserializeBinaryFromReader(message: ApiDocGetRequest, reader: jspb.BinaryReader): ApiDocGetRequest;
}

export namespace ApiDocGetRequest {
  export type AsObject = {
    apidoc?: dev_portal_api_dev_portal_v1_common_pb.ObjectRef.AsObject,
    withassets: boolean,
  }
}

export class ApiDocWriteRequest extends jspb.Message {
  hasApidoc(): boolean;
  clearApidoc(): void;
  getApidoc(): ApiDoc | undefined;
  setApidoc(value?: ApiDoc): void;

  clearPortalsList(): void;
  getPortalsList(): Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef>;
  setPortalsList(value: Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef>): void;
  addPortals(value?: dev_portal_api_dev_portal_v1_common_pb.ObjectRef, index?: number): dev_portal_api_dev_portal_v1_common_pb.ObjectRef;

  clearUsersList(): void;
  getUsersList(): Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef>;
  setUsersList(value: Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef>): void;
  addUsers(value?: dev_portal_api_dev_portal_v1_common_pb.ObjectRef, index?: number): dev_portal_api_dev_portal_v1_common_pb.ObjectRef;

  clearGroupsList(): void;
  getGroupsList(): Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef>;
  setGroupsList(value: Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef>): void;
  addGroups(value?: dev_portal_api_dev_portal_v1_common_pb.ObjectRef, index?: number): dev_portal_api_dev_portal_v1_common_pb.ObjectRef;

  getApiDocOnly(): boolean;
  setApiDocOnly(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ApiDocWriteRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ApiDocWriteRequest): ApiDocWriteRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ApiDocWriteRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ApiDocWriteRequest;
  static deserializeBinaryFromReader(message: ApiDocWriteRequest, reader: jspb.BinaryReader): ApiDocWriteRequest;
}

export namespace ApiDocWriteRequest {
  export type AsObject = {
    apidoc?: ApiDoc.AsObject,
    portalsList: Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef.AsObject>,
    usersList: Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef.AsObject>,
    groupsList: Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef.AsObject>,
    apiDocOnly: boolean,
  }
}

export class ApiDocFilter extends jspb.Message {
  clearPortalsList(): void;
  getPortalsList(): Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef>;
  setPortalsList(value: Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef>): void;
  addPortals(value?: dev_portal_api_dev_portal_v1_common_pb.ObjectRef, index?: number): dev_portal_api_dev_portal_v1_common_pb.ObjectRef;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ApiDocFilter.AsObject;
  static toObject(includeInstance: boolean, msg: ApiDocFilter): ApiDocFilter.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ApiDocFilter, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ApiDocFilter;
  static deserializeBinaryFromReader(message: ApiDocFilter, reader: jspb.BinaryReader): ApiDocFilter;
}

export namespace ApiDocFilter {
  export type AsObject = {
    portalsList: Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef.AsObject>,
  }
}
