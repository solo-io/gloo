/* eslint-disable */
// package: admin.devportal.solo.io
// file: dev-portal/api/grpc/admin/portal.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../../gogoproto/gogo_pb";
import * as google_protobuf_empty_pb from "google-protobuf/google/protobuf/empty_pb";
import * as dev_portal_api_dev_portal_v1_portal_pb from "../../../../dev-portal/api/dev-portal/v1/portal_pb";
import * as dev_portal_api_dev_portal_v1_common_pb from "../../../../dev-portal/api/dev-portal/v1/common_pb";
import * as dev_portal_api_grpc_common_common_pb from "../../../../dev-portal/api/grpc/common/common_pb";

export class Portal extends jspb.Message {
  hasMetadata(): boolean;
  clearMetadata(): void;
  getMetadata(): dev_portal_api_grpc_common_common_pb.ObjectMeta | undefined;
  setMetadata(value?: dev_portal_api_grpc_common_common_pb.ObjectMeta): void;

  hasSpec(): boolean;
  clearSpec(): void;
  getSpec(): dev_portal_api_dev_portal_v1_portal_pb.PortalSpec | undefined;
  setSpec(value?: dev_portal_api_dev_portal_v1_portal_pb.PortalSpec): void;

  hasStatus(): boolean;
  clearStatus(): void;
  getStatus(): dev_portal_api_dev_portal_v1_portal_pb.PortalStatus | undefined;
  setStatus(value?: dev_portal_api_dev_portal_v1_portal_pb.PortalStatus): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Portal.AsObject;
  static toObject(includeInstance: boolean, msg: Portal): Portal.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Portal, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Portal;
  static deserializeBinaryFromReader(message: Portal, reader: jspb.BinaryReader): Portal;
}

export namespace Portal {
  export type AsObject = {
    metadata?: dev_portal_api_grpc_common_common_pb.ObjectMeta.AsObject,
    spec?: dev_portal_api_dev_portal_v1_portal_pb.PortalSpec.AsObject,
    status?: dev_portal_api_dev_portal_v1_portal_pb.PortalStatus.AsObject,
  }
}

export class PortalList extends jspb.Message {
  clearPortalsList(): void;
  getPortalsList(): Array<Portal>;
  setPortalsList(value: Array<Portal>): void;
  addPortals(value?: Portal, index?: number): Portal;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PortalList.AsObject;
  static toObject(includeInstance: boolean, msg: PortalList): PortalList.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: PortalList, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PortalList;
  static deserializeBinaryFromReader(message: PortalList, reader: jspb.BinaryReader): PortalList;
}

export namespace PortalList {
  export type AsObject = {
    portalsList: Array<Portal.AsObject>,
  }
}

export class PortalWriteRequest extends jspb.Message {
  hasPortal(): boolean;
  clearPortal(): void;
  getPortal(): Portal | undefined;
  setPortal(value?: Portal): void;

  clearApidocsList(): void;
  getApidocsList(): Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef>;
  setApidocsList(value: Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef>): void;
  addApidocs(value?: dev_portal_api_dev_portal_v1_common_pb.ObjectRef, index?: number): dev_portal_api_dev_portal_v1_common_pb.ObjectRef;

  clearUsersList(): void;
  getUsersList(): Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef>;
  setUsersList(value: Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef>): void;
  addUsers(value?: dev_portal_api_dev_portal_v1_common_pb.ObjectRef, index?: number): dev_portal_api_dev_portal_v1_common_pb.ObjectRef;

  clearGroupsList(): void;
  getGroupsList(): Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef>;
  setGroupsList(value: Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef>): void;
  addGroups(value?: dev_portal_api_dev_portal_v1_common_pb.ObjectRef, index?: number): dev_portal_api_dev_portal_v1_common_pb.ObjectRef;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PortalWriteRequest.AsObject;
  static toObject(includeInstance: boolean, msg: PortalWriteRequest): PortalWriteRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: PortalWriteRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PortalWriteRequest;
  static deserializeBinaryFromReader(message: PortalWriteRequest, reader: jspb.BinaryReader): PortalWriteRequest;
}

export namespace PortalWriteRequest {
  export type AsObject = {
    portal?: Portal.AsObject,
    apidocsList: Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef.AsObject>,
    usersList: Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef.AsObject>,
    groupsList: Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef.AsObject>,
  }
}
