/* eslint-disable */
// package: devportal.solo.io
// file: dev-portal/api/dev-portal/v1/apidoc.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_timestamp_pb from "google-protobuf/google/protobuf/timestamp_pb";
import * as dev_portal_api_dev_portal_v1_common_pb from "../../../../dev-portal/api/dev-portal/v1/common_pb";
import * as gogoproto_gogo_pb from "../../../../gogoproto/gogo_pb";
import * as extproto_ext_pb from "../../../../protoc-gen-ext/extproto/ext_pb";

export class ApiDocSpec extends jspb.Message {
  getDisplayname(): string;
  setDisplayname(value: string): void;

  getVersion(): string;
  setVersion(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  hasDatasource(): boolean;
  clearDatasource(): void;
  getDatasource(): dev_portal_api_dev_portal_v1_common_pb.DataSource | undefined;
  setDatasource(value?: dev_portal_api_dev_portal_v1_common_pb.DataSource): void;

  hasOpenapi(): boolean;
  clearOpenapi(): void;
  getOpenapi(): ApiDocSpec.OpenApi | undefined;
  setOpenapi(value?: ApiDocSpec.OpenApi): void;

  getApidoctypeCase(): ApiDocSpec.ApidoctypeCase;
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
    displayname: string,
    version: string,
    description: string,
    datasource?: dev_portal_api_dev_portal_v1_common_pb.DataSource.AsObject,
    openapi?: ApiDocSpec.OpenApi.AsObject,
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

  export enum ApidoctypeCase {
    APIDOCTYPE_NOT_SET = 0,
    OPENAPI = 5,
  }
}

export class ApiDocStatus extends jspb.Message {
  getObservedgeneration(): number;
  setObservedgeneration(value: number): void;

  getState(): ApiDocStatus.StateMap[keyof ApiDocStatus.StateMap];
  setState(value: ApiDocStatus.StateMap[keyof ApiDocStatus.StateMap]): void;

  getReason(): string;
  setReason(value: string): void;

  hasModifieddate(): boolean;
  clearModifieddate(): void;
  getModifieddate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setModifieddate(value?: google_protobuf_timestamp_pb.Timestamp): void;

  getNumberofendpoints(): number;
  setNumberofendpoints(value: number): void;

  getBasepath(): string;
  setBasepath(value: string): void;

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
    observedgeneration: number,
    state: ApiDocStatus.StateMap[keyof ApiDocStatus.StateMap],
    reason: string,
    modifieddate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    numberofendpoints: number,
    basepath: string,
  }

  export interface StateMap {
    REJECTED: 0;
    ACCEPTED: 1;
  }

  export const State: StateMap;
}
