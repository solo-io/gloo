/* eslint-disable */
// package: admin.devportal.solo.io
// file: dev-portal/api/grpc/admin/api_key_scope.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../../gogoproto/gogo_pb";
import * as google_protobuf_empty_pb from "google-protobuf/google/protobuf/empty_pb";
import * as dev_portal_api_dev_portal_v1_portal_pb from "../../../../dev-portal/api/dev-portal/v1/portal_pb";
import * as dev_portal_api_dev_portal_v1_common_pb from "../../../../dev-portal/api/dev-portal/v1/common_pb";
import * as dev_portal_api_grpc_admin_apidoc_pb from "../../../../dev-portal/api/grpc/admin/apidoc_pb";

export class ApiKeyScope extends jspb.Message {
  hasSpec(): boolean;
  clearSpec(): void;
  getSpec(): dev_portal_api_dev_portal_v1_portal_pb.KeyScope | undefined;
  setSpec(value?: dev_portal_api_dev_portal_v1_portal_pb.KeyScope): void;

  hasStatus(): boolean;
  clearStatus(): void;
  getStatus(): dev_portal_api_dev_portal_v1_portal_pb.KeyScopeStatus | undefined;
  setStatus(value?: dev_portal_api_dev_portal_v1_portal_pb.KeyScopeStatus): void;

  hasPortal(): boolean;
  clearPortal(): void;
  getPortal(): dev_portal_api_dev_portal_v1_common_pb.ObjectRef | undefined;
  setPortal(value?: dev_portal_api_dev_portal_v1_common_pb.ObjectRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ApiKeyScope.AsObject;
  static toObject(includeInstance: boolean, msg: ApiKeyScope): ApiKeyScope.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ApiKeyScope, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ApiKeyScope;
  static deserializeBinaryFromReader(message: ApiKeyScope, reader: jspb.BinaryReader): ApiKeyScope;
}

export namespace ApiKeyScope {
  export type AsObject = {
    spec?: dev_portal_api_dev_portal_v1_portal_pb.KeyScope.AsObject,
    status?: dev_portal_api_dev_portal_v1_portal_pb.KeyScopeStatus.AsObject,
    portal?: dev_portal_api_dev_portal_v1_common_pb.ObjectRef.AsObject,
  }
}

export class ApiKeyScopeResponse extends jspb.Message {
  hasApiKeyScope(): boolean;
  clearApiKeyScope(): void;
  getApiKeyScope(): ApiKeyScope | undefined;
  setApiKeyScope(value?: ApiKeyScope): void;

  clearApiDocsList(): void;
  getApiDocsList(): Array<dev_portal_api_grpc_admin_apidoc_pb.ApiDoc>;
  setApiDocsList(value: Array<dev_portal_api_grpc_admin_apidoc_pb.ApiDoc>): void;
  addApiDocs(value?: dev_portal_api_grpc_admin_apidoc_pb.ApiDoc, index?: number): dev_portal_api_grpc_admin_apidoc_pb.ApiDoc;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ApiKeyScopeResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ApiKeyScopeResponse): ApiKeyScopeResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ApiKeyScopeResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ApiKeyScopeResponse;
  static deserializeBinaryFromReader(message: ApiKeyScopeResponse, reader: jspb.BinaryReader): ApiKeyScopeResponse;
}

export namespace ApiKeyScopeResponse {
  export type AsObject = {
    apiKeyScope?: ApiKeyScope.AsObject,
    apiDocsList: Array<dev_portal_api_grpc_admin_apidoc_pb.ApiDoc.AsObject>,
  }
}

export class ApiKeyScopeRef extends jspb.Message {
  clearApiKeyScopeList(): void;
  getApiKeyScopeList(): Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef>;
  setApiKeyScopeList(value: Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef>): void;
  addApiKeyScope(value?: dev_portal_api_dev_portal_v1_common_pb.ObjectRef, index?: number): dev_portal_api_dev_portal_v1_common_pb.ObjectRef;

  clearPortalList(): void;
  getPortalList(): Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef>;
  setPortalList(value: Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef>): void;
  addPortal(value?: dev_portal_api_dev_portal_v1_common_pb.ObjectRef, index?: number): dev_portal_api_dev_portal_v1_common_pb.ObjectRef;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ApiKeyScopeRef.AsObject;
  static toObject(includeInstance: boolean, msg: ApiKeyScopeRef): ApiKeyScopeRef.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ApiKeyScopeRef, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ApiKeyScopeRef;
  static deserializeBinaryFromReader(message: ApiKeyScopeRef, reader: jspb.BinaryReader): ApiKeyScopeRef;
}

export namespace ApiKeyScopeRef {
  export type AsObject = {
    apiKeyScopeList: Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef.AsObject>,
    portalList: Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef.AsObject>,
  }
}

export class ApiKeyScopeList extends jspb.Message {
  clearApiKeyScopesList(): void;
  getApiKeyScopesList(): Array<ApiKeyScopeResponse>;
  setApiKeyScopesList(value: Array<ApiKeyScopeResponse>): void;
  addApiKeyScopes(value?: ApiKeyScopeResponse, index?: number): ApiKeyScopeResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ApiKeyScopeList.AsObject;
  static toObject(includeInstance: boolean, msg: ApiKeyScopeList): ApiKeyScopeList.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ApiKeyScopeList, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ApiKeyScopeList;
  static deserializeBinaryFromReader(message: ApiKeyScopeList, reader: jspb.BinaryReader): ApiKeyScopeList;
}

export namespace ApiKeyScopeList {
  export type AsObject = {
    apiKeyScopesList: Array<ApiKeyScopeResponse.AsObject>,
  }
}

export class ApiKeyScopeFilter extends jspb.Message {
  clearPortalsList(): void;
  getPortalsList(): Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef>;
  setPortalsList(value: Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef>): void;
  addPortals(value?: dev_portal_api_dev_portal_v1_common_pb.ObjectRef, index?: number): dev_portal_api_dev_portal_v1_common_pb.ObjectRef;

  clearApiDocsList(): void;
  getApiDocsList(): Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef>;
  setApiDocsList(value: Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef>): void;
  addApiDocs(value?: dev_portal_api_dev_portal_v1_common_pb.ObjectRef, index?: number): dev_portal_api_dev_portal_v1_common_pb.ObjectRef;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ApiKeyScopeFilter.AsObject;
  static toObject(includeInstance: boolean, msg: ApiKeyScopeFilter): ApiKeyScopeFilter.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ApiKeyScopeFilter, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ApiKeyScopeFilter;
  static deserializeBinaryFromReader(message: ApiKeyScopeFilter, reader: jspb.BinaryReader): ApiKeyScopeFilter;
}

export namespace ApiKeyScopeFilter {
  export type AsObject = {
    portalsList: Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef.AsObject>,
    apiDocsList: Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef.AsObject>,
  }
}

export class ApiKeyScopeWriteRequest extends jspb.Message {
  hasApiKeyScope(): boolean;
  clearApiKeyScope(): void;
  getApiKeyScope(): ApiKeyScope | undefined;
  setApiKeyScope(value?: ApiKeyScope): void;

  clearApiDocsList(): void;
  getApiDocsList(): Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef>;
  setApiDocsList(value: Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef>): void;
  addApiDocs(value?: dev_portal_api_dev_portal_v1_common_pb.ObjectRef, index?: number): dev_portal_api_dev_portal_v1_common_pb.ObjectRef;

  getApiKeyScopeOnly(): boolean;
  setApiKeyScopeOnly(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ApiKeyScopeWriteRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ApiKeyScopeWriteRequest): ApiKeyScopeWriteRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ApiKeyScopeWriteRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ApiKeyScopeWriteRequest;
  static deserializeBinaryFromReader(message: ApiKeyScopeWriteRequest, reader: jspb.BinaryReader): ApiKeyScopeWriteRequest;
}

export namespace ApiKeyScopeWriteRequest {
  export type AsObject = {
    apiKeyScope?: ApiKeyScope.AsObject,
    apiDocsList: Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef.AsObject>,
    apiKeyScopeOnly: boolean,
  }
}
