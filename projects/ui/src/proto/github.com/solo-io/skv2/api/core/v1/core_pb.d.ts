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

export class ObjectRefList extends jspb.Message {
  clearRefsList(): void;
  getRefsList(): Array<ObjectRef>;
  setRefsList(value: Array<ObjectRef>): void;
  addRefs(value?: ObjectRef, index?: number): ObjectRef;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ObjectRefList.AsObject;
  static toObject(includeInstance: boolean, msg: ObjectRefList): ObjectRefList.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ObjectRefList, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ObjectRefList;
  static deserializeBinaryFromReader(message: ObjectRefList, reader: jspb.BinaryReader): ObjectRefList;
}

export namespace ObjectRefList {
  export type AsObject = {
    refsList: Array<ObjectRef.AsObject>,
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

export class ObjectSelector extends jspb.Message {
  clearNamespacesList(): void;
  getNamespacesList(): Array<string>;
  setNamespacesList(value: Array<string>): void;
  addNamespaces(value: string, index?: number): string;

  getLabelsMap(): jspb.Map<string, string>;
  clearLabelsMap(): void;
  clearExpressionsList(): void;
  getExpressionsList(): Array<ObjectSelector.Expression>;
  setExpressionsList(value: Array<ObjectSelector.Expression>): void;
  addExpressions(value?: ObjectSelector.Expression, index?: number): ObjectSelector.Expression;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ObjectSelector.AsObject;
  static toObject(includeInstance: boolean, msg: ObjectSelector): ObjectSelector.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ObjectSelector, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ObjectSelector;
  static deserializeBinaryFromReader(message: ObjectSelector, reader: jspb.BinaryReader): ObjectSelector;
}

export namespace ObjectSelector {
  export type AsObject = {
    namespacesList: Array<string>,
    labelsMap: Array<[string, string]>,
    expressionsList: Array<ObjectSelector.Expression.AsObject>,
  }

  export class Expression extends jspb.Message {
    getKey(): string;
    setKey(value: string): void;

    getOperator(): ObjectSelector.Expression.OperatorMap[keyof ObjectSelector.Expression.OperatorMap];
    setOperator(value: ObjectSelector.Expression.OperatorMap[keyof ObjectSelector.Expression.OperatorMap]): void;

    clearValuesList(): void;
    getValuesList(): Array<string>;
    setValuesList(value: Array<string>): void;
    addValues(value: string, index?: number): string;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Expression.AsObject;
    static toObject(includeInstance: boolean, msg: Expression): Expression.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Expression, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Expression;
    static deserializeBinaryFromReader(message: Expression, reader: jspb.BinaryReader): Expression;
  }

  export namespace Expression {
    export type AsObject = {
      key: string,
      operator: ObjectSelector.Expression.OperatorMap[keyof ObjectSelector.Expression.OperatorMap],
      valuesList: Array<string>,
    }

    export interface OperatorMap {
      EQUALS: 0;
      DOUBLEEQUALS: 1;
      NOTEQUALS: 2;
      IN: 3;
      NOTIN: 4;
      EXISTS: 5;
      DOESNOTEXIST: 6;
      GREATERTHAN: 7;
      LESSTHAN: 8;
    }

    export const Operator: OperatorMap;
  }
}
