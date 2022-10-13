/* eslint-disable */
// package: multicluster.solo.io
// file: github.com/solo-io/solo-projects/projects/gloo-fed/api/multicluster/v1alpha1/multicluster.proto

import * as jspb from "google-protobuf";
import * as extproto_ext_pb from "../../../../../../../../extproto/ext_pb";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";
import * as github_com_solo_io_skv2_api_core_v1_core_pb from "../../../../../../../../github.com/solo-io/skv2/api/core/v1/core_pb";

export class Placement extends jspb.Message {
  clearNamespacesList(): void;
  getNamespacesList(): Array<string>;
  setNamespacesList(value: Array<string>): void;
  addNamespaces(value: string, index?: number): string;

  clearClustersList(): void;
  getClustersList(): Array<string>;
  setClustersList(value: Array<string>): void;
  addClusters(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Placement.AsObject;
  static toObject(includeInstance: boolean, msg: Placement): Placement.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Placement, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Placement;
  static deserializeBinaryFromReader(message: Placement, reader: jspb.BinaryReader): Placement;
}

export namespace Placement {
  export type AsObject = {
    namespacesList: Array<string>,
    clustersList: Array<string>,
  }
}

export class MultiClusterRoleSpec extends jspb.Message {
  clearRulesList(): void;
  getRulesList(): Array<MultiClusterRoleSpec.Rule>;
  setRulesList(value: Array<MultiClusterRoleSpec.Rule>): void;
  addRules(value?: MultiClusterRoleSpec.Rule, index?: number): MultiClusterRoleSpec.Rule;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MultiClusterRoleSpec.AsObject;
  static toObject(includeInstance: boolean, msg: MultiClusterRoleSpec): MultiClusterRoleSpec.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: MultiClusterRoleSpec, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MultiClusterRoleSpec;
  static deserializeBinaryFromReader(message: MultiClusterRoleSpec, reader: jspb.BinaryReader): MultiClusterRoleSpec;
}

export namespace MultiClusterRoleSpec {
  export type AsObject = {
    rulesList: Array<MultiClusterRoleSpec.Rule.AsObject>,
  }

  export class Rule extends jspb.Message {
    getApiGroup(): string;
    setApiGroup(value: string): void;

    hasKind(): boolean;
    clearKind(): void;
    getKind(): google_protobuf_wrappers_pb.StringValue | undefined;
    setKind(value?: google_protobuf_wrappers_pb.StringValue): void;

    getAction(): MultiClusterRoleSpec.Rule.ActionMap[keyof MultiClusterRoleSpec.Rule.ActionMap];
    setAction(value: MultiClusterRoleSpec.Rule.ActionMap[keyof MultiClusterRoleSpec.Rule.ActionMap]): void;

    clearPlacementsList(): void;
    getPlacementsList(): Array<Placement>;
    setPlacementsList(value: Array<Placement>): void;
    addPlacements(value?: Placement, index?: number): Placement;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Rule.AsObject;
    static toObject(includeInstance: boolean, msg: Rule): Rule.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Rule, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Rule;
    static deserializeBinaryFromReader(message: Rule, reader: jspb.BinaryReader): Rule;
  }

  export namespace Rule {
    export type AsObject = {
      apiGroup: string,
      kind?: google_protobuf_wrappers_pb.StringValue.AsObject,
      action: MultiClusterRoleSpec.Rule.ActionMap[keyof MultiClusterRoleSpec.Rule.ActionMap],
      placementsList: Array<Placement.AsObject>,
    }

    export interface ActionMap {
      ANY: 0;
      CREATE: 1;
      UPDATE: 2;
      DELETE: 3;
    }

    export const Action: ActionMap;
  }
}

export class MultiClusterRoleStatus extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MultiClusterRoleStatus.AsObject;
  static toObject(includeInstance: boolean, msg: MultiClusterRoleStatus): MultiClusterRoleStatus.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: MultiClusterRoleStatus, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MultiClusterRoleStatus;
  static deserializeBinaryFromReader(message: MultiClusterRoleStatus, reader: jspb.BinaryReader): MultiClusterRoleStatus;
}

export namespace MultiClusterRoleStatus {
  export type AsObject = {
  }
}

export class MultiClusterRoleBindingSpec extends jspb.Message {
  clearSubjectsList(): void;
  getSubjectsList(): Array<github_com_solo_io_skv2_api_core_v1_core_pb.TypedObjectRef>;
  setSubjectsList(value: Array<github_com_solo_io_skv2_api_core_v1_core_pb.TypedObjectRef>): void;
  addSubjects(value?: github_com_solo_io_skv2_api_core_v1_core_pb.TypedObjectRef, index?: number): github_com_solo_io_skv2_api_core_v1_core_pb.TypedObjectRef;

  hasRoleRef(): boolean;
  clearRoleRef(): void;
  getRoleRef(): github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef | undefined;
  setRoleRef(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MultiClusterRoleBindingSpec.AsObject;
  static toObject(includeInstance: boolean, msg: MultiClusterRoleBindingSpec): MultiClusterRoleBindingSpec.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: MultiClusterRoleBindingSpec, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MultiClusterRoleBindingSpec;
  static deserializeBinaryFromReader(message: MultiClusterRoleBindingSpec, reader: jspb.BinaryReader): MultiClusterRoleBindingSpec;
}

export namespace MultiClusterRoleBindingSpec {
  export type AsObject = {
    subjectsList: Array<github_com_solo_io_skv2_api_core_v1_core_pb.TypedObjectRef.AsObject>,
    roleRef?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef.AsObject,
  }
}

export class MultiClusterRoleBindingStatus extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MultiClusterRoleBindingStatus.AsObject;
  static toObject(includeInstance: boolean, msg: MultiClusterRoleBindingStatus): MultiClusterRoleBindingStatus.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: MultiClusterRoleBindingStatus, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MultiClusterRoleBindingStatus;
  static deserializeBinaryFromReader(message: MultiClusterRoleBindingStatus, reader: jspb.BinaryReader): MultiClusterRoleBindingStatus;
}

export namespace MultiClusterRoleBindingStatus {
  export type AsObject = {
  }
}
