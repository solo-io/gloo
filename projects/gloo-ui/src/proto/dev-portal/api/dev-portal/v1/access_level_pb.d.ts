/* eslint-disable */
// package: devportal.solo.io
// file: dev-portal/api/dev-portal/v1/access_level.proto

import * as jspb from "google-protobuf";
import * as dev_portal_api_dev_portal_v1_common_pb from "../../../../dev-portal/api/dev-portal/v1/common_pb";
import * as gogoproto_gogo_pb from "../../../../gogoproto/gogo_pb";
import * as extproto_ext_pb from "../../../../protoc-gen-ext/extproto/ext_pb";

export class AccessLevel extends jspb.Message {
  hasPortalselector(): boolean;
  clearPortalselector(): void;
  getPortalselector(): dev_portal_api_dev_portal_v1_common_pb.Selector | undefined;
  setPortalselector(value?: dev_portal_api_dev_portal_v1_common_pb.Selector): void;

  hasApidocselector(): boolean;
  clearApidocselector(): void;
  getApidocselector(): dev_portal_api_dev_portal_v1_common_pb.Selector | undefined;
  setApidocselector(value?: dev_portal_api_dev_portal_v1_common_pb.Selector): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AccessLevel.AsObject;
  static toObject(includeInstance: boolean, msg: AccessLevel): AccessLevel.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: AccessLevel, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AccessLevel;
  static deserializeBinaryFromReader(message: AccessLevel, reader: jspb.BinaryReader): AccessLevel;
}

export namespace AccessLevel {
  export type AsObject = {
    portalselector?: dev_portal_api_dev_portal_v1_common_pb.Selector.AsObject,
    apidocselector?: dev_portal_api_dev_portal_v1_common_pb.Selector.AsObject,
  }
}

export class AccessLevelStatus extends jspb.Message {
  clearPortalsList(): void;
  getPortalsList(): Array<string>;
  setPortalsList(value: Array<string>): void;
  addPortals(value: string, index?: number): string;

  clearApidocsList(): void;
  getApidocsList(): Array<string>;
  setApidocsList(value: Array<string>): void;
  addApidocs(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AccessLevelStatus.AsObject;
  static toObject(includeInstance: boolean, msg: AccessLevelStatus): AccessLevelStatus.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: AccessLevelStatus, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AccessLevelStatus;
  static deserializeBinaryFromReader(message: AccessLevelStatus, reader: jspb.BinaryReader): AccessLevelStatus;
}

export namespace AccessLevelStatus {
  export type AsObject = {
    portalsList: Array<string>,
    apidocsList: Array<string>,
  }
}
