/* eslint-disable */
// package: devportal.solo.io
// file: dev-portal/api/dev-portal/v1/group.proto

import * as jspb from "google-protobuf";
import * as dev_portal_api_dev_portal_v1_common_pb from "../../../../dev-portal/api/dev-portal/v1/common_pb";
import * as gogoproto_gogo_pb from "../../../../gogoproto/gogo_pb";
import * as extproto_ext_pb from "../../../../protoc-gen-ext/extproto/ext_pb";

export class GroupSpec extends jspb.Message {
  getUserlabelsMap(): jspb.Map<string, string>;
  clearUserlabelsMap(): void;
  getApidoclabelsMap(): jspb.Map<string, string>;
  clearApidoclabelsMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GroupSpec.AsObject;
  static toObject(includeInstance: boolean, msg: GroupSpec): GroupSpec.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GroupSpec, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GroupSpec;
  static deserializeBinaryFromReader(message: GroupSpec, reader: jspb.BinaryReader): GroupSpec;
}

export namespace GroupSpec {
  export type AsObject = {
    userlabelsMap: Array<[string, string]>,
    apidoclabelsMap: Array<[string, string]>,
  }
}

export class GroupStatus extends jspb.Message {
  getObservedgeneration(): number;
  setObservedgeneration(value: number): void;

  clearUsersList(): void;
  getUsersList(): Array<string>;
  setUsersList(value: Array<string>): void;
  addUsers(value: string, index?: number): string;

  clearApidocsList(): void;
  getApidocsList(): Array<string>;
  setApidocsList(value: Array<string>): void;
  addApidocs(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GroupStatus.AsObject;
  static toObject(includeInstance: boolean, msg: GroupStatus): GroupStatus.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GroupStatus, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GroupStatus;
  static deserializeBinaryFromReader(message: GroupStatus, reader: jspb.BinaryReader): GroupStatus;
}

export namespace GroupStatus {
  export type AsObject = {
    observedgeneration: number,
    usersList: Array<string>,
    apidocsList: Array<string>,
  }
}
