/* eslint-disable */
// package: devportal.solo.io
// file: dev-portal/api/dev-portal/v1/apidoc.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_timestamp_pb from "google-protobuf/google/protobuf/timestamp_pb";
import * as dev_portal_api_dev_portal_v1_common_pb from "../../../../dev-portal/api/dev-portal/v1/common_pb";
import * as gogoproto_gogo_pb from "../../../../gogoproto/gogo_pb";
import * as extproto_ext_pb from "../../../../protoc-gen-ext/extproto/ext_pb";

export class ApiDocSpec extends jspb.Message {
  hasDataSource(): boolean;
  clearDataSource(): void;
  getDataSource(): dev_portal_api_dev_portal_v1_common_pb.DataSource | undefined;
  setDataSource(value?: dev_portal_api_dev_portal_v1_common_pb.DataSource): void;

  hasImage(): boolean;
  clearImage(): void;
  getImage(): dev_portal_api_dev_portal_v1_common_pb.DataSource | undefined;
  setImage(value?: dev_portal_api_dev_portal_v1_common_pb.DataSource): void;

  hasOpenApi(): boolean;
  clearOpenApi(): void;
  getOpenApi(): ApiDocSpec.OpenApi | undefined;
  setOpenApi(value?: ApiDocSpec.OpenApi): void;

  getApiDocTypeCase(): ApiDocSpec.ApiDocTypeCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ApiDocSpec.AsObject;
  static toObject(includeInstance: boolean, msg: ApiDocSpec): ApiDocSpec.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ApiDocSpec, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ApiDocSpec;
  static deserializeBinaryFromReader(message: ApiDocSpec, reader: jspb.BinaryReader): ApiDocSpec;
}

export namespace ApiDocSpec {
  export type AsObject = {
    dataSource?: dev_portal_api_dev_portal_v1_common_pb.DataSource.AsObject,
    image?: dev_portal_api_dev_portal_v1_common_pb.DataSource.AsObject,
    openApi?: ApiDocSpec.OpenApi.AsObject,
  }

  export class OpenApi extends jspb.Message {
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): OpenApi.AsObject;
    static toObject(includeInstance: boolean, msg: OpenApi): OpenApi.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: OpenApi, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): OpenApi;
    static deserializeBinaryFromReader(message: OpenApi, reader: jspb.BinaryReader): OpenApi;
  }

  export namespace OpenApi {
    export type AsObject = {
    }
  }

  export enum ApiDocTypeCase {
    API_DOC_TYPE_NOT_SET = 0,
    OPEN_API = 5,
  }
}

export class ApiDocStatus extends jspb.Message {
  getObservedGeneration(): number;
  setObservedGeneration(value: number): void;

  getState(): dev_portal_api_dev_portal_v1_common_pb.StateMap[keyof dev_portal_api_dev_portal_v1_common_pb.StateMap];
  setState(value: dev_portal_api_dev_portal_v1_common_pb.StateMap[keyof dev_portal_api_dev_portal_v1_common_pb.StateMap]): void;

  getReason(): string;
  setReason(value: string): void;

  hasModifiedDate(): boolean;
  clearModifiedDate(): void;
  getModifiedDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setModifiedDate(value?: google_protobuf_timestamp_pb.Timestamp): void;

  getDisplayName(): string;
  setDisplayName(value: string): void;

  getVersion(): string;
  setVersion(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  getNumberOfEndpoints(): number;
  setNumberOfEndpoints(value: number): void;

  getBasePath(): string;
  setBasePath(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ApiDocStatus.AsObject;
  static toObject(includeInstance: boolean, msg: ApiDocStatus): ApiDocStatus.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ApiDocStatus, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ApiDocStatus;
  static deserializeBinaryFromReader(message: ApiDocStatus, reader: jspb.BinaryReader): ApiDocStatus;
}

export namespace ApiDocStatus {
  export type AsObject = {
    observedGeneration: number,
    state: dev_portal_api_dev_portal_v1_common_pb.StateMap[keyof dev_portal_api_dev_portal_v1_common_pb.StateMap],
    reason: string,
    modifiedDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    displayName: string,
    version: string,
    description: string,
    numberOfEndpoints: number,
    basePath: string,
  }
}
