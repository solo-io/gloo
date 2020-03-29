/* eslint-disable */
// package: admin.devportal.solo.io
// file: dev-portal/api/grpc/admin/api_key.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../../gogoproto/gogo_pb";
import * as google_protobuf_empty_pb from "google-protobuf/google/protobuf/empty_pb";
import * as dev_portal_api_dev_portal_v1_common_pb from "../../../../dev-portal/api/dev-portal/v1/common_pb";
import * as dev_portal_api_grpc_common_common_pb from "../../../../dev-portal/api/grpc/common/common_pb";

export class ApiKey extends jspb.Message {
  hasMetadata(): boolean;
  clearMetadata(): void;
  getMetadata(): dev_portal_api_grpc_common_common_pb.ObjectMeta | undefined;
  setMetadata(value?: dev_portal_api_grpc_common_common_pb.ObjectMeta): void;

  getValue(): string;
  setValue(value: string): void;

  hasUser(): boolean;
  clearUser(): void;
  getUser(): dev_portal_api_dev_portal_v1_common_pb.ObjectRef | undefined;
  setUser(value?: dev_portal_api_dev_portal_v1_common_pb.ObjectRef): void;

  hasKeyScope(): boolean;
  clearKeyScope(): void;
  getKeyScope(): dev_portal_api_dev_portal_v1_common_pb.ObjectRef | undefined;
  setKeyScope(value?: dev_portal_api_dev_portal_v1_common_pb.ObjectRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ApiKey.AsObject;
  static toObject(includeInstance: boolean, msg: ApiKey): ApiKey.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ApiKey, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ApiKey;
  static deserializeBinaryFromReader(message: ApiKey, reader: jspb.BinaryReader): ApiKey;
}

export namespace ApiKey {
  export type AsObject = {
    metadata?: dev_portal_api_grpc_common_common_pb.ObjectMeta.AsObject,
    value: string,
    user?: dev_portal_api_dev_portal_v1_common_pb.ObjectRef.AsObject,
    keyScope?: dev_portal_api_dev_portal_v1_common_pb.ObjectRef.AsObject,
  }
}

export class ApiKeyList extends jspb.Message {
  clearApiKeysList(): void;
  getApiKeysList(): Array<ApiKey>;
  setApiKeysList(value: Array<ApiKey>): void;
  addApiKeys(value?: ApiKey, index?: number): ApiKey;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ApiKeyList.AsObject;
  static toObject(includeInstance: boolean, msg: ApiKeyList): ApiKeyList.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ApiKeyList, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ApiKeyList;
  static deserializeBinaryFromReader(message: ApiKeyList, reader: jspb.BinaryReader): ApiKeyList;
}

export namespace ApiKeyList {
  export type AsObject = {
    apiKeysList: Array<ApiKey.AsObject>,
  }
}
