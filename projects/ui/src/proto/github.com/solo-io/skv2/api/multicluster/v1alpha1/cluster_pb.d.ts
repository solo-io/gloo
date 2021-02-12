/* eslint-disable */
// package: multicluster.solo.io
// file: github.com/solo-io/skv2/api/multicluster/v1alpha1/cluster.proto

import * as jspb from "google-protobuf";
import * as extproto_ext_pb from "../../../../../../extproto/ext_pb";
import * as github_com_solo_io_skv2_api_core_v1_core_pb from "../../../../../../github.com/solo-io/skv2/api/core/v1/core_pb";

export class KubernetesClusterSpec extends jspb.Message {
  getSecretName(): string;
  setSecretName(value: string): void;

  getClusterDomain(): string;
  setClusterDomain(value: string): void;

  hasProviderInfo(): boolean;
  clearProviderInfo(): void;
  getProviderInfo(): KubernetesClusterSpec.ProviderInfo | undefined;
  setProviderInfo(value?: KubernetesClusterSpec.ProviderInfo): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): KubernetesClusterSpec.AsObject;
  static toObject(includeInstance: boolean, msg: KubernetesClusterSpec): KubernetesClusterSpec.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: KubernetesClusterSpec, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): KubernetesClusterSpec;
  static deserializeBinaryFromReader(message: KubernetesClusterSpec, reader: jspb.BinaryReader): KubernetesClusterSpec;
}

export namespace KubernetesClusterSpec {
  export type AsObject = {
    secretName: string,
    clusterDomain: string,
    providerInfo?: KubernetesClusterSpec.ProviderInfo.AsObject,
  }

  export class ProviderInfo extends jspb.Message {
    hasEks(): boolean;
    clearEks(): void;
    getEks(): KubernetesClusterSpec.Eks | undefined;
    setEks(value?: KubernetesClusterSpec.Eks): void;

    getProviderInfoTypeCase(): ProviderInfo.ProviderInfoTypeCase;
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ProviderInfo.AsObject;
    static toObject(includeInstance: boolean, msg: ProviderInfo): ProviderInfo.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: ProviderInfo, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ProviderInfo;
    static deserializeBinaryFromReader(message: ProviderInfo, reader: jspb.BinaryReader): ProviderInfo;
  }

  export namespace ProviderInfo {
    export type AsObject = {
      eks?: KubernetesClusterSpec.Eks.AsObject,
    }

    export enum ProviderInfoTypeCase {
      PROVIDER_INFO_TYPE_NOT_SET = 0,
      EKS = 1,
    }
  }

  export class Eks extends jspb.Message {
    getArn(): string;
    setArn(value: string): void;

    getAccountId(): string;
    setAccountId(value: string): void;

    getRegion(): string;
    setRegion(value: string): void;

    getName(): string;
    setName(value: string): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Eks.AsObject;
    static toObject(includeInstance: boolean, msg: Eks): Eks.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Eks, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Eks;
    static deserializeBinaryFromReader(message: Eks, reader: jspb.BinaryReader): Eks;
  }

  export namespace Eks {
    export type AsObject = {
      arn: string,
      accountId: string,
      region: string,
      name: string,
    }
  }
}

export class KubernetesClusterStatus extends jspb.Message {
  clearStatusList(): void;
  getStatusList(): Array<github_com_solo_io_skv2_api_core_v1_core_pb.Status>;
  setStatusList(value: Array<github_com_solo_io_skv2_api_core_v1_core_pb.Status>): void;
  addStatus(value?: github_com_solo_io_skv2_api_core_v1_core_pb.Status, index?: number): github_com_solo_io_skv2_api_core_v1_core_pb.Status;

  getNamespace(): string;
  setNamespace(value: string): void;

  clearPolicyRulesList(): void;
  getPolicyRulesList(): Array<PolicyRule>;
  setPolicyRulesList(value: Array<PolicyRule>): void;
  addPolicyRules(value?: PolicyRule, index?: number): PolicyRule;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): KubernetesClusterStatus.AsObject;
  static toObject(includeInstance: boolean, msg: KubernetesClusterStatus): KubernetesClusterStatus.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: KubernetesClusterStatus, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): KubernetesClusterStatus;
  static deserializeBinaryFromReader(message: KubernetesClusterStatus, reader: jspb.BinaryReader): KubernetesClusterStatus;
}

export namespace KubernetesClusterStatus {
  export type AsObject = {
    statusList: Array<github_com_solo_io_skv2_api_core_v1_core_pb.Status.AsObject>,
    namespace: string,
    policyRulesList: Array<PolicyRule.AsObject>,
  }
}

export class PolicyRule extends jspb.Message {
  clearVerbsList(): void;
  getVerbsList(): Array<string>;
  setVerbsList(value: Array<string>): void;
  addVerbs(value: string, index?: number): string;

  clearApiGroupsList(): void;
  getApiGroupsList(): Array<string>;
  setApiGroupsList(value: Array<string>): void;
  addApiGroups(value: string, index?: number): string;

  clearResourcesList(): void;
  getResourcesList(): Array<string>;
  setResourcesList(value: Array<string>): void;
  addResources(value: string, index?: number): string;

  clearResourceNamesList(): void;
  getResourceNamesList(): Array<string>;
  setResourceNamesList(value: Array<string>): void;
  addResourceNames(value: string, index?: number): string;

  clearNonResourceUrlsList(): void;
  getNonResourceUrlsList(): Array<string>;
  setNonResourceUrlsList(value: Array<string>): void;
  addNonResourceUrls(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PolicyRule.AsObject;
  static toObject(includeInstance: boolean, msg: PolicyRule): PolicyRule.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: PolicyRule, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PolicyRule;
  static deserializeBinaryFromReader(message: PolicyRule, reader: jspb.BinaryReader): PolicyRule;
}

export namespace PolicyRule {
  export type AsObject = {
    verbsList: Array<string>,
    apiGroupsList: Array<string>,
    resourcesList: Array<string>,
    resourceNamesList: Array<string>,
    nonResourceUrlsList: Array<string>,
  }
}
