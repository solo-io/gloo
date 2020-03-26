/* eslint-disable */
// package: devportal.solo.io
// file: dev-portal/api/dev-portal/v1/group.proto

import * as jspb from "google-protobuf";
import * as dev_portal_api_dev_portal_v1_common_pb from "../../../../dev-portal/api/dev-portal/v1/common_pb";
import * as dev_portal_api_dev_portal_v1_access_level_pb from "../../../../dev-portal/api/dev-portal/v1/access_level_pb";
import * as gogoproto_gogo_pb from "../../../../gogoproto/gogo_pb";
import * as extproto_ext_pb from "../../../../protoc-gen-ext/extproto/ext_pb";

export class GroupSpec extends jspb.Message {
  getDisplayName(): string;
  setDisplayName(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  hasUserSelector(): boolean;
  clearUserSelector(): void;
  getUserSelector(): dev_portal_api_dev_portal_v1_common_pb.Selector | undefined;
  setUserSelector(value?: dev_portal_api_dev_portal_v1_common_pb.Selector): void;

  hasAccessLevel(): boolean;
  clearAccessLevel(): void;
  getAccessLevel(): dev_portal_api_dev_portal_v1_access_level_pb.AccessLevel | undefined;
  setAccessLevel(value?: dev_portal_api_dev_portal_v1_access_level_pb.AccessLevel): void;

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
    displayName: string,
    description: string,
    userSelector?: dev_portal_api_dev_portal_v1_common_pb.Selector.AsObject,
    accessLevel?: dev_portal_api_dev_portal_v1_access_level_pb.AccessLevel.AsObject,
  }
}

export class GroupStatus extends jspb.Message {
  getObservedGeneration(): number;
  setObservedGeneration(value: number): void;

  clearUsersList(): void;
  getUsersList(): Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef>;
  setUsersList(value: Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef>): void;
  addUsers(value?: dev_portal_api_dev_portal_v1_common_pb.ObjectRef, index?: number): dev_portal_api_dev_portal_v1_common_pb.ObjectRef;

  hasAccessLevel(): boolean;
  clearAccessLevel(): void;
  getAccessLevel(): dev_portal_api_dev_portal_v1_access_level_pb.AccessLevelStatus | undefined;
  setAccessLevel(value?: dev_portal_api_dev_portal_v1_access_level_pb.AccessLevelStatus): void;

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
    observedGeneration: number,
    usersList: Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef.AsObject>,
    accessLevel?: dev_portal_api_dev_portal_v1_access_level_pb.AccessLevelStatus.AsObject,
  }
}
