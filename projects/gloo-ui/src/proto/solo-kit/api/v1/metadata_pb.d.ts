// package: core.solo.io
// file: solo-kit/api/v1/metadata.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../gogoproto/gogo_pb";
import * as extproto_ext_pb from "../../../extproto/ext_pb";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";

export class Metadata extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getNamespace(): string;
  setNamespace(value: string): void;

  getCluster(): string;
  setCluster(value: string): void;

  getResourceVersion(): string;
  setResourceVersion(value: string): void;

  getLabelsMap(): jspb.Map<string, string>;
  clearLabelsMap(): void;
  getAnnotationsMap(): jspb.Map<string, string>;
  clearAnnotationsMap(): void;
  getGeneration(): number;
  setGeneration(value: number): void;

  clearOwnerReferencesList(): void;
  getOwnerReferencesList(): Array<Metadata.OwnerReference>;
  setOwnerReferencesList(value: Array<Metadata.OwnerReference>): void;
  addOwnerReferences(value?: Metadata.OwnerReference, index?: number): Metadata.OwnerReference;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Metadata.AsObject;
  static toObject(includeInstance: boolean, msg: Metadata): Metadata.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Metadata, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Metadata;
  static deserializeBinaryFromReader(message: Metadata, reader: jspb.BinaryReader): Metadata;
}

export namespace Metadata {
  export type AsObject = {
    name: string,
    namespace: string,
    cluster: string,
    resourceVersion: string,
    labelsMap: Array<[string, string]>,
    annotationsMap: Array<[string, string]>,
    generation: number,
    ownerReferencesList: Array<Metadata.OwnerReference.AsObject>,
  }

  export class OwnerReference extends jspb.Message {
    getApiVersion(): string;
    setApiVersion(value: string): void;

    hasBlockOwnerDeletion(): boolean;
    clearBlockOwnerDeletion(): void;
    getBlockOwnerDeletion(): google_protobuf_wrappers_pb.BoolValue | undefined;
    setBlockOwnerDeletion(value?: google_protobuf_wrappers_pb.BoolValue): void;

    hasController(): boolean;
    clearController(): void;
    getController(): google_protobuf_wrappers_pb.BoolValue | undefined;
    setController(value?: google_protobuf_wrappers_pb.BoolValue): void;

    getKind(): string;
    setKind(value: string): void;

    getName(): string;
    setName(value: string): void;

    getUid(): string;
    setUid(value: string): void;

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
      apiVersion: string,
      blockOwnerDeletion?: google_protobuf_wrappers_pb.BoolValue.AsObject,
      controller?: google_protobuf_wrappers_pb.BoolValue.AsObject,
      kind: string,
      name: string,
      uid: string,
    }
  }
}

