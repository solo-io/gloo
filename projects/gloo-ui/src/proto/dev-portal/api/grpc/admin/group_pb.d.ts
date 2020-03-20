/* eslint-disable */
// package: admin.devportal.solo.io
// file: dev-portal/api/grpc/admin/group.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../../gogoproto/gogo_pb";
import * as google_protobuf_empty_pb from "google-protobuf/google/protobuf/empty_pb";
import * as dev_portal_api_dev_portal_v1_group_pb from "../../../../dev-portal/api/dev-portal/v1/group_pb";
import * as dev_portal_api_dev_portal_v1_common_pb from "../../../../dev-portal/api/dev-portal/v1/common_pb";
import * as dev_portal_api_grpc_common_common_pb from "../../../../dev-portal/api/grpc/common/common_pb";

export class Group extends jspb.Message {
  hasMetadata(): boolean;
  clearMetadata(): void;
  getMetadata(): dev_portal_api_grpc_common_common_pb.ObjectMeta | undefined;
  setMetadata(value?: dev_portal_api_grpc_common_common_pb.ObjectMeta): void;

  hasSpec(): boolean;
  clearSpec(): void;
  getSpec(): dev_portal_api_dev_portal_v1_group_pb.GroupSpec | undefined;
  setSpec(value?: dev_portal_api_dev_portal_v1_group_pb.GroupSpec): void;

  hasStatus(): boolean;
  clearStatus(): void;
  getStatus(): dev_portal_api_dev_portal_v1_group_pb.GroupStatus | undefined;
  setStatus(value?: dev_portal_api_dev_portal_v1_group_pb.GroupStatus): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Group.AsObject;
  static toObject(includeInstance: boolean, msg: Group): Group.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Group, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Group;
  static deserializeBinaryFromReader(message: Group, reader: jspb.BinaryReader): Group;
}

export namespace Group {
  export type AsObject = {
    metadata?: dev_portal_api_grpc_common_common_pb.ObjectMeta.AsObject,
    spec?: dev_portal_api_dev_portal_v1_group_pb.GroupSpec.AsObject,
    status?: dev_portal_api_dev_portal_v1_group_pb.GroupStatus.AsObject,
  }
}

export class GroupList extends jspb.Message {
  clearGroupsList(): void;
  getGroupsList(): Array<Group>;
  setGroupsList(value: Array<Group>): void;
  addGroups(value?: Group, index?: number): Group;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GroupList.AsObject;
  static toObject(includeInstance: boolean, msg: GroupList): GroupList.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GroupList, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GroupList;
  static deserializeBinaryFromReader(message: GroupList, reader: jspb.BinaryReader): GroupList;
}

export namespace GroupList {
  export type AsObject = {
    groupsList: Array<Group.AsObject>,
  }
}

export class GroupWriteRequest extends jspb.Message {
  hasGroup(): boolean;
  clearGroup(): void;
  getGroup(): Group | undefined;
  setGroup(value?: Group): void;

  clearApiDocsList(): void;
  getApiDocsList(): Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef>;
  setApiDocsList(value: Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef>): void;
  addApiDocs(value?: dev_portal_api_dev_portal_v1_common_pb.ObjectRef, index?: number): dev_portal_api_dev_portal_v1_common_pb.ObjectRef;

  clearPortalsList(): void;
  getPortalsList(): Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef>;
  setPortalsList(value: Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef>): void;
  addPortals(value?: dev_portal_api_dev_portal_v1_common_pb.ObjectRef, index?: number): dev_portal_api_dev_portal_v1_common_pb.ObjectRef;

  clearUsersList(): void;
  getUsersList(): Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef>;
  setUsersList(value: Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef>): void;
  addUsers(value?: dev_portal_api_dev_portal_v1_common_pb.ObjectRef, index?: number): dev_portal_api_dev_portal_v1_common_pb.ObjectRef;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GroupWriteRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GroupWriteRequest): GroupWriteRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GroupWriteRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GroupWriteRequest;
  static deserializeBinaryFromReader(message: GroupWriteRequest, reader: jspb.BinaryReader): GroupWriteRequest;
}

export namespace GroupWriteRequest {
  export type AsObject = {
    group?: Group.AsObject,
    apiDocsList: Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef.AsObject>,
    portalsList: Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef.AsObject>,
    usersList: Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef.AsObject>,
  }
}

export class GroupFilter extends jspb.Message {
  clearApiDocsList(): void;
  getApiDocsList(): Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef>;
  setApiDocsList(value: Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef>): void;
  addApiDocs(value?: dev_portal_api_dev_portal_v1_common_pb.ObjectRef, index?: number): dev_portal_api_dev_portal_v1_common_pb.ObjectRef;

  clearPortalsList(): void;
  getPortalsList(): Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef>;
  setPortalsList(value: Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef>): void;
  addPortals(value?: dev_portal_api_dev_portal_v1_common_pb.ObjectRef, index?: number): dev_portal_api_dev_portal_v1_common_pb.ObjectRef;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GroupFilter.AsObject;
  static toObject(includeInstance: boolean, msg: GroupFilter): GroupFilter.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GroupFilter, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GroupFilter;
  static deserializeBinaryFromReader(message: GroupFilter, reader: jspb.BinaryReader): GroupFilter;
}

export namespace GroupFilter {
  export type AsObject = {
    apiDocsList: Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef.AsObject>,
    portalsList: Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef.AsObject>,
  }
}
