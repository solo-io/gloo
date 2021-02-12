/* eslint-disable */
// package: core.fed.solo.io
// file: github.com/solo-io/solo-projects/projects/gloo-fed/api/fed/core/v1/placement.proto

import * as jspb from "google-protobuf";
import * as extproto_ext_pb from "../../../../../../../../../extproto/ext_pb";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";

export class TemplateMetadata extends jspb.Message {
  getAnnotationsMap(): jspb.Map<string, string>;
  clearAnnotationsMap(): void;
  getLabelsMap(): jspb.Map<string, string>;
  clearLabelsMap(): void;
  getName(): string;
  setName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TemplateMetadata.AsObject;
  static toObject(includeInstance: boolean, msg: TemplateMetadata): TemplateMetadata.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: TemplateMetadata, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TemplateMetadata;
  static deserializeBinaryFromReader(message: TemplateMetadata, reader: jspb.BinaryReader): TemplateMetadata;
}

export namespace TemplateMetadata {
  export type AsObject = {
    annotationsMap: Array<[string, string]>,
    labelsMap: Array<[string, string]>,
    name: string,
  }
}

export class PlacementStatus extends jspb.Message {
  getClustersMap(): jspb.Map<string, PlacementStatus.Cluster>;
  clearClustersMap(): void;
  getState(): PlacementStatus.StateMap[keyof PlacementStatus.StateMap];
  setState(value: PlacementStatus.StateMap[keyof PlacementStatus.StateMap]): void;

  getMessage(): string;
  setMessage(value: string): void;

  getObservedGeneration(): number;
  setObservedGeneration(value: number): void;

  getWrittenBy(): string;
  setWrittenBy(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PlacementStatus.AsObject;
  static toObject(includeInstance: boolean, msg: PlacementStatus): PlacementStatus.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: PlacementStatus, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PlacementStatus;
  static deserializeBinaryFromReader(message: PlacementStatus, reader: jspb.BinaryReader): PlacementStatus;
}

export namespace PlacementStatus {
  export type AsObject = {
    clustersMap: Array<[string, PlacementStatus.Cluster.AsObject]>,
    state: PlacementStatus.StateMap[keyof PlacementStatus.StateMap],
    message: string,
    observedGeneration: number,
    writtenBy: string,
  }

  export class Namespace extends jspb.Message {
    getState(): PlacementStatus.StateMap[keyof PlacementStatus.StateMap];
    setState(value: PlacementStatus.StateMap[keyof PlacementStatus.StateMap]): void;

    getMessage(): string;
    setMessage(value: string): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Namespace.AsObject;
    static toObject(includeInstance: boolean, msg: Namespace): Namespace.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Namespace, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Namespace;
    static deserializeBinaryFromReader(message: Namespace, reader: jspb.BinaryReader): Namespace;
  }

  export namespace Namespace {
    export type AsObject = {
      state: PlacementStatus.StateMap[keyof PlacementStatus.StateMap],
      message: string,
    }
  }

  export class Cluster extends jspb.Message {
    getNamespacesMap(): jspb.Map<string, PlacementStatus.Namespace>;
    clearNamespacesMap(): void;
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Cluster.AsObject;
    static toObject(includeInstance: boolean, msg: Cluster): Cluster.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Cluster, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Cluster;
    static deserializeBinaryFromReader(message: Cluster, reader: jspb.BinaryReader): Cluster;
  }

  export namespace Cluster {
    export type AsObject = {
      namespacesMap: Array<[string, PlacementStatus.Namespace.AsObject]>,
    }
  }

  export interface StateMap {
    UNKNOWN: 0;
    PLACED: 1;
    FAILED: 2;
    STALE: 3;
    INVALID: 4;
    PENDING: 5;
  }

  export const State: StateMap;
}
