/* eslint-disable */
// package: core.skv2.solo.io
// file: github.com/solo-io/skv2/api/core/v1/core.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";
import * as google_protobuf_timestamp_pb from "google-protobuf/google/protobuf/timestamp_pb";
import * as extproto_ext_pb from "../../../../../../extproto/ext_pb";

export class ObjectRef extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getNamespace(): string;
  setNamespace(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ObjectRef.AsObject;
  static toObject(includeInstance: boolean, msg: ObjectRef): ObjectRef.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ObjectRef, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ObjectRef;
  static deserializeBinaryFromReader(message: ObjectRef, reader: jspb.BinaryReader): ObjectRef;
}

export namespace ObjectRef {
  export type AsObject = {
    name: string,
    namespace: string,
  }
}

export class ClusterObjectRef extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getNamespace(): string;
  setNamespace(value: string): void;

  getClusterName(): string;
  setClusterName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ClusterObjectRef.AsObject;
  static toObject(includeInstance: boolean, msg: ClusterObjectRef): ClusterObjectRef.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ClusterObjectRef, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ClusterObjectRef;
  static deserializeBinaryFromReader(message: ClusterObjectRef, reader: jspb.BinaryReader): ClusterObjectRef;
}

export namespace ClusterObjectRef {
  export type AsObject = {
    name: string,
    namespace: string,
    clusterName: string,
  }
}

export class TypedObjectRef extends jspb.Message {
  hasApiGroup(): boolean;
  clearApiGroup(): void;
  getApiGroup(): google_protobuf_wrappers_pb.StringValue | undefined;
  setApiGroup(value?: google_protobuf_wrappers_pb.StringValue): void;

  hasKind(): boolean;
  clearKind(): void;
  getKind(): google_protobuf_wrappers_pb.StringValue | undefined;
  setKind(value?: google_protobuf_wrappers_pb.StringValue): void;

  getName(): string;
  setName(value: string): void;

  getNamespace(): string;
  setNamespace(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TypedObjectRef.AsObject;
  static toObject(includeInstance: boolean, msg: TypedObjectRef): TypedObjectRef.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: TypedObjectRef, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TypedObjectRef;
  static deserializeBinaryFromReader(message: TypedObjectRef, reader: jspb.BinaryReader): TypedObjectRef;
}

export namespace TypedObjectRef {
  export type AsObject = {
    apiGroup?: google_protobuf_wrappers_pb.StringValue.AsObject,
    kind?: google_protobuf_wrappers_pb.StringValue.AsObject,
    name: string,
    namespace: string,
  }
}

export class TypedClusterObjectRef extends jspb.Message {
  hasApiGroup(): boolean;
  clearApiGroup(): void;
  getApiGroup(): google_protobuf_wrappers_pb.StringValue | undefined;
  setApiGroup(value?: google_protobuf_wrappers_pb.StringValue): void;

  hasKind(): boolean;
  clearKind(): void;
  getKind(): google_protobuf_wrappers_pb.StringValue | undefined;
  setKind(value?: google_protobuf_wrappers_pb.StringValue): void;

  getName(): string;
  setName(value: string): void;

  getNamespace(): string;
  setNamespace(value: string): void;

  getClusterName(): string;
  setClusterName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TypedClusterObjectRef.AsObject;
  static toObject(includeInstance: boolean, msg: TypedClusterObjectRef): TypedClusterObjectRef.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: TypedClusterObjectRef, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TypedClusterObjectRef;
  static deserializeBinaryFromReader(message: TypedClusterObjectRef, reader: jspb.BinaryReader): TypedClusterObjectRef;
}

export namespace TypedClusterObjectRef {
  export type AsObject = {
    apiGroup?: google_protobuf_wrappers_pb.StringValue.AsObject,
    kind?: google_protobuf_wrappers_pb.StringValue.AsObject,
    name: string,
    namespace: string,
    clusterName: string,
  }
}

export class Status extends jspb.Message {
  getState(): Status.StateMap[keyof Status.StateMap];
  setState(value: Status.StateMap[keyof Status.StateMap]): void;

  getMessage(): string;
  setMessage(value: string): void;

  getObservedGeneration(): number;
  setObservedGeneration(value: number): void;

  hasProcessingTime(): boolean;
  clearProcessingTime(): void;
  getProcessingTime(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setProcessingTime(value?: google_protobuf_timestamp_pb.Timestamp): void;

  hasOwner(): boolean;
  clearOwner(): void;
  getOwner(): google_protobuf_wrappers_pb.StringValue | undefined;
  setOwner(value?: google_protobuf_wrappers_pb.StringValue): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Status.AsObject;
  static toObject(includeInstance: boolean, msg: Status): Status.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Status, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Status;
  static deserializeBinaryFromReader(message: Status, reader: jspb.BinaryReader): Status;
}

export namespace Status {
  export type AsObject = {
    state: Status.StateMap[keyof Status.StateMap],
    message: string,
    observedGeneration: number,
    processingTime?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    owner?: google_protobuf_wrappers_pb.StringValue.AsObject,
  }

  export interface StateMap {
    PENDING: 0;
    PROCESSING: 1;
    INVALID: 2;
    FAILED: 3;
    ACCEPTED: 4;
  }

  export const State: StateMap;
}
