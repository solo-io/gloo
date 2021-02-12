/* eslint-disable */
// package: fed.rpc.solo.io
// file: github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/glooinstance.proto

import * as jspb from "google-protobuf";
import * as extproto_ext_pb from "../../../../../../../../extproto/ext_pb";
import * as github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_v1_instance_pb from "../../../../../../../../github.com/solo-io/solo-projects/projects/gloo-fed/api/fed/v1/instance_pb";
import * as github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb from "../../../../../../../../github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/common_pb";
import * as github_com_solo_io_skv2_api_core_v1_core_pb from "../../../../../../../../github.com/solo-io/skv2/api/core/v1/core_pb";

export class GlooInstance extends jspb.Message {
  hasMetadata(): boolean;
  clearMetadata(): void;
  getMetadata(): github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb.ObjectMeta | undefined;
  setMetadata(value?: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb.ObjectMeta): void;

  hasSpec(): boolean;
  clearSpec(): void;
  getSpec(): github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_v1_instance_pb.GlooInstanceSpec | undefined;
  setSpec(value?: github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_v1_instance_pb.GlooInstanceSpec): void;

  hasStatus(): boolean;
  clearStatus(): void;
  getStatus(): github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_v1_instance_pb.GlooInstanceStatus | undefined;
  setStatus(value?: github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_v1_instance_pb.GlooInstanceStatus): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GlooInstance.AsObject;
  static toObject(includeInstance: boolean, msg: GlooInstance): GlooInstance.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GlooInstance, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GlooInstance;
  static deserializeBinaryFromReader(message: GlooInstance, reader: jspb.BinaryReader): GlooInstance;
}

export namespace GlooInstance {
  export type AsObject = {
    metadata?: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb.ObjectMeta.AsObject,
    spec?: github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_v1_instance_pb.GlooInstanceSpec.AsObject,
    status?: github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_v1_instance_pb.GlooInstanceStatus.AsObject,
  }
}

export class ListGlooInstancesRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListGlooInstancesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListGlooInstancesRequest): ListGlooInstancesRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListGlooInstancesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListGlooInstancesRequest;
  static deserializeBinaryFromReader(message: ListGlooInstancesRequest, reader: jspb.BinaryReader): ListGlooInstancesRequest;
}

export namespace ListGlooInstancesRequest {
  export type AsObject = {
  }
}

export class ListGlooInstancesResponse extends jspb.Message {
  clearGlooInstancesList(): void;
  getGlooInstancesList(): Array<GlooInstance>;
  setGlooInstancesList(value: Array<GlooInstance>): void;
  addGlooInstances(value?: GlooInstance, index?: number): GlooInstance;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListGlooInstancesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListGlooInstancesResponse): ListGlooInstancesResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListGlooInstancesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListGlooInstancesResponse;
  static deserializeBinaryFromReader(message: ListGlooInstancesResponse, reader: jspb.BinaryReader): ListGlooInstancesResponse;
}

export namespace ListGlooInstancesResponse {
  export type AsObject = {
    glooInstancesList: Array<GlooInstance.AsObject>,
  }
}

export class ClusterDetails extends jspb.Message {
  getCluster(): string;
  setCluster(value: string): void;

  clearGlooInstancesList(): void;
  getGlooInstancesList(): Array<GlooInstance>;
  setGlooInstancesList(value: Array<GlooInstance>): void;
  addGlooInstances(value?: GlooInstance, index?: number): GlooInstance;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ClusterDetails.AsObject;
  static toObject(includeInstance: boolean, msg: ClusterDetails): ClusterDetails.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ClusterDetails, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ClusterDetails;
  static deserializeBinaryFromReader(message: ClusterDetails, reader: jspb.BinaryReader): ClusterDetails;
}

export namespace ClusterDetails {
  export type AsObject = {
    cluster: string,
    glooInstancesList: Array<GlooInstance.AsObject>,
  }
}

export class ListClusterDetailsRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListClusterDetailsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListClusterDetailsRequest): ListClusterDetailsRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListClusterDetailsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListClusterDetailsRequest;
  static deserializeBinaryFromReader(message: ListClusterDetailsRequest, reader: jspb.BinaryReader): ListClusterDetailsRequest;
}

export namespace ListClusterDetailsRequest {
  export type AsObject = {
  }
}

export class ListClusterDetailsResponse extends jspb.Message {
  clearClusterDetailsList(): void;
  getClusterDetailsList(): Array<ClusterDetails>;
  setClusterDetailsList(value: Array<ClusterDetails>): void;
  addClusterDetails(value?: ClusterDetails, index?: number): ClusterDetails;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListClusterDetailsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListClusterDetailsResponse): ListClusterDetailsResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListClusterDetailsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListClusterDetailsResponse;
  static deserializeBinaryFromReader(message: ListClusterDetailsResponse, reader: jspb.BinaryReader): ListClusterDetailsResponse;
}

export namespace ListClusterDetailsResponse {
  export type AsObject = {
    clusterDetailsList: Array<ClusterDetails.AsObject>,
  }
}

export class ConfigDump extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getRaw(): string;
  setRaw(value: string): void;

  getError(): string;
  setError(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ConfigDump.AsObject;
  static toObject(includeInstance: boolean, msg: ConfigDump): ConfigDump.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ConfigDump, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ConfigDump;
  static deserializeBinaryFromReader(message: ConfigDump, reader: jspb.BinaryReader): ConfigDump;
}

export namespace ConfigDump {
  export type AsObject = {
    name: string,
    raw: string,
    error: string,
  }
}

export class GetConfigDumpsRequest extends jspb.Message {
  hasGlooInstanceRef(): boolean;
  clearGlooInstanceRef(): void;
  getGlooInstanceRef(): github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef | undefined;
  setGlooInstanceRef(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetConfigDumpsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetConfigDumpsRequest): GetConfigDumpsRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetConfigDumpsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetConfigDumpsRequest;
  static deserializeBinaryFromReader(message: GetConfigDumpsRequest, reader: jspb.BinaryReader): GetConfigDumpsRequest;
}

export namespace GetConfigDumpsRequest {
  export type AsObject = {
    glooInstanceRef?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef.AsObject,
  }
}

export class GetConfigDumpsResponse extends jspb.Message {
  clearConfigDumpsList(): void;
  getConfigDumpsList(): Array<ConfigDump>;
  setConfigDumpsList(value: Array<ConfigDump>): void;
  addConfigDumps(value?: ConfigDump, index?: number): ConfigDump;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetConfigDumpsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetConfigDumpsResponse): GetConfigDumpsResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetConfigDumpsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetConfigDumpsResponse;
  static deserializeBinaryFromReader(message: GetConfigDumpsResponse, reader: jspb.BinaryReader): GetConfigDumpsResponse;
}

export namespace GetConfigDumpsResponse {
  export type AsObject = {
    configDumpsList: Array<ConfigDump.AsObject>,
  }
}
