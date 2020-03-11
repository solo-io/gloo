/* eslint-disable */
// package: common.devportal.solo.io
// file: dev-portal/api/grpc/common/common.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../../gogoproto/gogo_pb";
import * as extproto_ext_pb from "../../../../protoc-gen-ext/extproto/ext_pb";

export class ObjectReference extends jspb.Message {
  getKind(): string;
  setKind(value: string): void;

  getNamespace(): string;
  setNamespace(value: string): void;

  getName(): string;
  setName(value: string): void;

  getUid(): string;
  setUid(value: string): void;

  getApiversion(): string;
  setApiversion(value: string): void;

  getResourceversion(): string;
  setResourceversion(value: string): void;

  getFieldpath(): string;
  setFieldpath(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ObjectReference.AsObject;
  static toObject(includeInstance: boolean, msg: ObjectReference): ObjectReference.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ObjectReference, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ObjectReference;
  static deserializeBinaryFromReader(message: ObjectReference, reader: jspb.BinaryReader): ObjectReference;
}

export namespace ObjectReference {
  export type AsObject = {
    kind: string,
    namespace: string,
    name: string,
    uid: string,
    apiversion: string,
    resourceversion: string,
    fieldpath: string,
  }
}

export class ObjectMeta extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getGeneratename(): string;
  setGeneratename(value: string): void;

  getNamespace(): string;
  setNamespace(value: string): void;

  getSelflink(): string;
  setSelflink(value: string): void;

  getUid(): string;
  setUid(value: string): void;

  getResourceversion(): string;
  setResourceversion(value: string): void;

  getGeneration(): number;
  setGeneration(value: number): void;

  hasCreationtimestamp(): boolean;
  clearCreationtimestamp(): void;
  getCreationtimestamp(): Time | undefined;
  setCreationtimestamp(value?: Time): void;

  hasDeletiontimestamp(): boolean;
  clearDeletiontimestamp(): void;
  getDeletiontimestamp(): Time | undefined;
  setDeletiontimestamp(value?: Time): void;

  getDeletiongraceperiodseconds(): number;
  setDeletiongraceperiodseconds(value: number): void;

  getLabelsMap(): jspb.Map<string, string>;
  clearLabelsMap(): void;
  getAnnotationsMap(): jspb.Map<string, string>;
  clearAnnotationsMap(): void;
  clearOwnerreferencesList(): void;
  getOwnerreferencesList(): Array<OwnerReference>;
  setOwnerreferencesList(value: Array<OwnerReference>): void;
  addOwnerreferences(value?: OwnerReference, index?: number): OwnerReference;

  clearFinalizersList(): void;
  getFinalizersList(): Array<string>;
  setFinalizersList(value: Array<string>): void;
  addFinalizers(value: string, index?: number): string;

  getClustername(): string;
  setClustername(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ObjectMeta.AsObject;
  static toObject(includeInstance: boolean, msg: ObjectMeta): ObjectMeta.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ObjectMeta, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ObjectMeta;
  static deserializeBinaryFromReader(message: ObjectMeta, reader: jspb.BinaryReader): ObjectMeta;
}

export namespace ObjectMeta {
  export type AsObject = {
    name: string,
    generatename: string,
    namespace: string,
    selflink: string,
    uid: string,
    resourceversion: string,
    generation: number,
    creationtimestamp?: Time.AsObject,
    deletiontimestamp?: Time.AsObject,
    deletiongraceperiodseconds: number,
    labelsMap: Array<[string, string]>,
    annotationsMap: Array<[string, string]>,
    ownerreferencesList: Array<OwnerReference.AsObject>,
    finalizersList: Array<string>,
    clustername: string,
  }
}

export class OwnerReference extends jspb.Message {
  getApiversion(): string;
  setApiversion(value: string): void;

  getKind(): string;
  setKind(value: string): void;

  getName(): string;
  setName(value: string): void;

  getUid(): string;
  setUid(value: string): void;

  getController(): boolean;
  setController(value: boolean): void;

  getBlockownerdeletion(): boolean;
  setBlockownerdeletion(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OwnerReference.AsObject;
  static toObject(includeInstance: boolean, msg: OwnerReference): OwnerReference.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: OwnerReference, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OwnerReference;
  static deserializeBinaryFromReader(message: OwnerReference, reader: jspb.BinaryReader): OwnerReference;
}

export namespace OwnerReference {
  export type AsObject = {
    apiversion: string,
    kind: string,
    name: string,
    uid: string,
    controller: boolean,
    blockownerdeletion: boolean,
  }
}

export class Time extends jspb.Message {
  getSeconds(): number;
  setSeconds(value: number): void;

  getNanos(): number;
  setNanos(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Time.AsObject;
  static toObject(includeInstance: boolean, msg: Time): Time.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Time, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Time;
  static deserializeBinaryFromReader(message: Time, reader: jspb.BinaryReader): Time;
}

export namespace Time {
  export type AsObject = {
    seconds: number,
    nanos: number,
  }
}
