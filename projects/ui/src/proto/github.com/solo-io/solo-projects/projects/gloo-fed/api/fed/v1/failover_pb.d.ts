/* eslint-disable */
// package: fed.solo.io
// file: github.com/solo-io/solo-projects/projects/gloo-fed/api/fed/v1/failover.proto

import * as jspb from "google-protobuf";
import * as extproto_ext_pb from "../../../../../../../../extproto/ext_pb";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";
import * as google_protobuf_timestamp_pb from "google-protobuf/google/protobuf/timestamp_pb";
import * as github_com_solo_io_skv2_api_core_v1_core_pb from "../../../../../../../../github.com/solo-io/skv2/api/core/v1/core_pb";

export class FailoverSchemeSpec extends jspb.Message {
  hasPrimary(): boolean;
  clearPrimary(): void;
  getPrimary(): github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef | undefined;
  setPrimary(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef): void;

  clearFailoverGroupsList(): void;
  getFailoverGroupsList(): Array<FailoverSchemeSpec.FailoverEndpoints>;
  setFailoverGroupsList(value: Array<FailoverSchemeSpec.FailoverEndpoints>): void;
  addFailoverGroups(value?: FailoverSchemeSpec.FailoverEndpoints, index?: number): FailoverSchemeSpec.FailoverEndpoints;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FailoverSchemeSpec.AsObject;
  static toObject(includeInstance: boolean, msg: FailoverSchemeSpec): FailoverSchemeSpec.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: FailoverSchemeSpec, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FailoverSchemeSpec;
  static deserializeBinaryFromReader(message: FailoverSchemeSpec, reader: jspb.BinaryReader): FailoverSchemeSpec;
}

export namespace FailoverSchemeSpec {
  export type AsObject = {
    primary?: github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef.AsObject,
    failoverGroupsList: Array<FailoverSchemeSpec.FailoverEndpoints.AsObject>,
  }

  export class FailoverEndpoints extends jspb.Message {
    clearPriorityGroupList(): void;
    getPriorityGroupList(): Array<FailoverSchemeSpec.FailoverEndpoints.LocalityLbTargets>;
    setPriorityGroupList(value: Array<FailoverSchemeSpec.FailoverEndpoints.LocalityLbTargets>): void;
    addPriorityGroup(value?: FailoverSchemeSpec.FailoverEndpoints.LocalityLbTargets, index?: number): FailoverSchemeSpec.FailoverEndpoints.LocalityLbTargets;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): FailoverEndpoints.AsObject;
    static toObject(includeInstance: boolean, msg: FailoverEndpoints): FailoverEndpoints.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: FailoverEndpoints, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): FailoverEndpoints;
    static deserializeBinaryFromReader(message: FailoverEndpoints, reader: jspb.BinaryReader): FailoverEndpoints;
  }

  export namespace FailoverEndpoints {
    export type AsObject = {
      priorityGroupList: Array<FailoverSchemeSpec.FailoverEndpoints.LocalityLbTargets.AsObject>,
    }

    export class LocalityLbTargets extends jspb.Message {
      getCluster(): string;
      setCluster(value: string): void;

      clearUpstreamsList(): void;
      getUpstreamsList(): Array<github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef>;
      setUpstreamsList(value: Array<github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef>): void;
      addUpstreams(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef, index?: number): github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef;

      hasLocalityWeight(): boolean;
      clearLocalityWeight(): void;
      getLocalityWeight(): google_protobuf_wrappers_pb.UInt32Value | undefined;
      setLocalityWeight(value?: google_protobuf_wrappers_pb.UInt32Value): void;

      serializeBinary(): Uint8Array;
      toObject(includeInstance?: boolean): LocalityLbTargets.AsObject;
      static toObject(includeInstance: boolean, msg: LocalityLbTargets): LocalityLbTargets.AsObject;
      static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
      static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
      static serializeBinaryToWriter(message: LocalityLbTargets, writer: jspb.BinaryWriter): void;
      static deserializeBinary(bytes: Uint8Array): LocalityLbTargets;
      static deserializeBinaryFromReader(message: LocalityLbTargets, reader: jspb.BinaryReader): LocalityLbTargets;
    }

    export namespace LocalityLbTargets {
      export type AsObject = {
        cluster: string,
        upstreamsList: Array<github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef.AsObject>,
        localityWeight?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
      }
    }
  }
}

export class FailoverSchemeStatus extends jspb.Message {
  getState(): FailoverSchemeStatus.StateMap[keyof FailoverSchemeStatus.StateMap];
  setState(value: FailoverSchemeStatus.StateMap[keyof FailoverSchemeStatus.StateMap]): void;

  getMessage(): string;
  setMessage(value: string): void;

  getObservedGeneration(): number;
  setObservedGeneration(value: number): void;

  hasProcessingTime(): boolean;
  clearProcessingTime(): void;
  getProcessingTime(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setProcessingTime(value?: google_protobuf_timestamp_pb.Timestamp): void;

  getNamespacedStatusesMap(): jspb.Map<string, FailoverSchemeStatus.Status>;
  clearNamespacedStatusesMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FailoverSchemeStatus.AsObject;
  static toObject(includeInstance: boolean, msg: FailoverSchemeStatus): FailoverSchemeStatus.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: FailoverSchemeStatus, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FailoverSchemeStatus;
  static deserializeBinaryFromReader(message: FailoverSchemeStatus, reader: jspb.BinaryReader): FailoverSchemeStatus;
}

export namespace FailoverSchemeStatus {
  export type AsObject = {
    state: FailoverSchemeStatus.StateMap[keyof FailoverSchemeStatus.StateMap],
    message: string,
    observedGeneration: number,
    processingTime?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    namespacedStatusesMap: Array<[string, FailoverSchemeStatus.Status.AsObject]>,
  }

  export class Status extends jspb.Message {
    getState(): FailoverSchemeStatus.StateMap[keyof FailoverSchemeStatus.StateMap];
    setState(value: FailoverSchemeStatus.StateMap[keyof FailoverSchemeStatus.StateMap]): void;

    getMessage(): string;
    setMessage(value: string): void;

    getObservedGeneration(): number;
    setObservedGeneration(value: number): void;

    hasProcessingTime(): boolean;
    clearProcessingTime(): void;
    getProcessingTime(): google_protobuf_timestamp_pb.Timestamp | undefined;
    setProcessingTime(value?: google_protobuf_timestamp_pb.Timestamp): void;

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
      state: FailoverSchemeStatus.StateMap[keyof FailoverSchemeStatus.StateMap],
      message: string,
      observedGeneration: number,
      processingTime?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    }
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
