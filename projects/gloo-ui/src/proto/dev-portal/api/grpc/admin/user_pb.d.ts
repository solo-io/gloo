/* eslint-disable */
// package: admin.devportal.solo.io
// file: dev-portal/api/grpc/admin/user.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../../gogoproto/gogo_pb";
import * as google_protobuf_empty_pb from "google-protobuf/google/protobuf/empty_pb";
import * as dev_portal_api_dev_portal_v1_user_pb from "../../../../dev-portal/api/dev-portal/v1/user_pb";
import * as dev_portal_api_dev_portal_v1_common_pb from "../../../../dev-portal/api/dev-portal/v1/common_pb";
import * as dev_portal_api_grpc_common_common_pb from "../../../../dev-portal/api/grpc/common/common_pb";

export class User extends jspb.Message {
  hasMetadata(): boolean;
  clearMetadata(): void;
  getMetadata(): dev_portal_api_grpc_common_common_pb.ObjectMeta | undefined;
  setMetadata(value?: dev_portal_api_grpc_common_common_pb.ObjectMeta): void;

  hasSpec(): boolean;
  clearSpec(): void;
  getSpec(): dev_portal_api_dev_portal_v1_user_pb.UserSpec | undefined;
  setSpec(value?: dev_portal_api_dev_portal_v1_user_pb.UserSpec): void;

  hasStatus(): boolean;
  clearStatus(): void;
  getStatus(): dev_portal_api_dev_portal_v1_user_pb.UserStatus | undefined;
  setStatus(value?: dev_portal_api_dev_portal_v1_user_pb.UserStatus): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): User.AsObject;
  static toObject(includeInstance: boolean, msg: User): User.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: User, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): User;
  static deserializeBinaryFromReader(message: User, reader: jspb.BinaryReader): User;
}

export namespace User {
  export type AsObject = {
    metadata?: dev_portal_api_grpc_common_common_pb.ObjectMeta.AsObject,
    spec?: dev_portal_api_dev_portal_v1_user_pb.UserSpec.AsObject,
    status?: dev_portal_api_dev_portal_v1_user_pb.UserStatus.AsObject,
  }
}

export class UserList extends jspb.Message {
  clearUsersList(): void;
  getUsersList(): Array<User>;
  setUsersList(value: Array<User>): void;
  addUsers(value?: User, index?: number): User;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserList.AsObject;
  static toObject(includeInstance: boolean, msg: UserList): UserList.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UserList, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserList;
  static deserializeBinaryFromReader(message: UserList, reader: jspb.BinaryReader): UserList;
}

export namespace UserList {
  export type AsObject = {
    usersList: Array<User.AsObject>,
  }
}

export class UserWriteRequest extends jspb.Message {
  hasUser(): boolean;
  clearUser(): void;
  getUser(): User | undefined;
  setUser(value?: User): void;

  getPassword(): string;
  setPassword(value: string): void;

  clearApiDocsList(): void;
  getApiDocsList(): Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef>;
  setApiDocsList(value: Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef>): void;
  addApiDocs(value?: dev_portal_api_dev_portal_v1_common_pb.ObjectRef, index?: number): dev_portal_api_dev_portal_v1_common_pb.ObjectRef;

  clearPortalsList(): void;
  getPortalsList(): Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef>;
  setPortalsList(value: Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef>): void;
  addPortals(value?: dev_portal_api_dev_portal_v1_common_pb.ObjectRef, index?: number): dev_portal_api_dev_portal_v1_common_pb.ObjectRef;

  clearGroupsList(): void;
  getGroupsList(): Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef>;
  setGroupsList(value: Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef>): void;
  addGroups(value?: dev_portal_api_dev_portal_v1_common_pb.ObjectRef, index?: number): dev_portal_api_dev_portal_v1_common_pb.ObjectRef;

  getUserOnly(): boolean;
  setUserOnly(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserWriteRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UserWriteRequest): UserWriteRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UserWriteRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserWriteRequest;
  static deserializeBinaryFromReader(message: UserWriteRequest, reader: jspb.BinaryReader): UserWriteRequest;
}

export namespace UserWriteRequest {
  export type AsObject = {
    user?: User.AsObject,
    password: string,
    apiDocsList: Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef.AsObject>,
    portalsList: Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef.AsObject>,
    groupsList: Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef.AsObject>,
    userOnly: boolean,
  }
}

export class UserFilter extends jspb.Message {
  clearApiDocsList(): void;
  getApiDocsList(): Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef>;
  setApiDocsList(value: Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef>): void;
  addApiDocs(value?: dev_portal_api_dev_portal_v1_common_pb.ObjectRef, index?: number): dev_portal_api_dev_portal_v1_common_pb.ObjectRef;

  clearPortalsList(): void;
  getPortalsList(): Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef>;
  setPortalsList(value: Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef>): void;
  addPortals(value?: dev_portal_api_dev_portal_v1_common_pb.ObjectRef, index?: number): dev_portal_api_dev_portal_v1_common_pb.ObjectRef;

  clearGroupsList(): void;
  getGroupsList(): Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef>;
  setGroupsList(value: Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef>): void;
  addGroups(value?: dev_portal_api_dev_portal_v1_common_pb.ObjectRef, index?: number): dev_portal_api_dev_portal_v1_common_pb.ObjectRef;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserFilter.AsObject;
  static toObject(includeInstance: boolean, msg: UserFilter): UserFilter.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UserFilter, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserFilter;
  static deserializeBinaryFromReader(message: UserFilter, reader: jspb.BinaryReader): UserFilter;
}

export namespace UserFilter {
  export type AsObject = {
    apiDocsList: Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef.AsObject>,
    portalsList: Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef.AsObject>,
    groupsList: Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef.AsObject>,
  }
}
