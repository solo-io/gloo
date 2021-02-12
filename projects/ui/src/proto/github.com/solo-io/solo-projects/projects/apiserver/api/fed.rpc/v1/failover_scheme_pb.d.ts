/* eslint-disable */
// package: fed.rpc.solo.io
// file: github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/failover_scheme.proto

import * as jspb from "google-protobuf";
import * as extproto_ext_pb from "../../../../../../../../extproto/ext_pb";
import * as github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_v1_failover_pb from "../../../../../../../../github.com/solo-io/solo-projects/projects/gloo-fed/api/fed/v1/failover_pb";
import * as github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb from "../../../../../../../../github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/common_pb";
import * as github_com_solo_io_skv2_api_core_v1_core_pb from "../../../../../../../../github.com/solo-io/skv2/api/core/v1/core_pb";

export class FailoverScheme extends jspb.Message {
  hasMetadata(): boolean;
  clearMetadata(): void;
  getMetadata(): github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb.ObjectMeta | undefined;
  setMetadata(value?: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb.ObjectMeta): void;

  hasSpec(): boolean;
  clearSpec(): void;
  getSpec(): github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_v1_failover_pb.FailoverSchemeSpec | undefined;
  setSpec(value?: github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_v1_failover_pb.FailoverSchemeSpec): void;

  hasStatus(): boolean;
  clearStatus(): void;
  getStatus(): github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_v1_failover_pb.FailoverSchemeStatus | undefined;
  setStatus(value?: github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_v1_failover_pb.FailoverSchemeStatus): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FailoverScheme.AsObject;
  static toObject(includeInstance: boolean, msg: FailoverScheme): FailoverScheme.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: FailoverScheme, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FailoverScheme;
  static deserializeBinaryFromReader(message: FailoverScheme, reader: jspb.BinaryReader): FailoverScheme;
}

export namespace FailoverScheme {
  export type AsObject = {
    metadata?: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb.ObjectMeta.AsObject,
    spec?: github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_v1_failover_pb.FailoverSchemeSpec.AsObject,
    status?: github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_v1_failover_pb.FailoverSchemeStatus.AsObject,
  }
}

export class GetFailoverSchemeRequest extends jspb.Message {
  hasUpstreamRef(): boolean;
  clearUpstreamRef(): void;
  getUpstreamRef(): github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef | undefined;
  setUpstreamRef(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetFailoverSchemeRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetFailoverSchemeRequest): GetFailoverSchemeRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetFailoverSchemeRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetFailoverSchemeRequest;
  static deserializeBinaryFromReader(message: GetFailoverSchemeRequest, reader: jspb.BinaryReader): GetFailoverSchemeRequest;
}

export namespace GetFailoverSchemeRequest {
  export type AsObject = {
    upstreamRef?: github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef.AsObject,
  }
}

export class GetFailoverSchemeResponse extends jspb.Message {
  hasFailoverScheme(): boolean;
  clearFailoverScheme(): void;
  getFailoverScheme(): FailoverScheme | undefined;
  setFailoverScheme(value?: FailoverScheme): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetFailoverSchemeResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetFailoverSchemeResponse): GetFailoverSchemeResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetFailoverSchemeResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetFailoverSchemeResponse;
  static deserializeBinaryFromReader(message: GetFailoverSchemeResponse, reader: jspb.BinaryReader): GetFailoverSchemeResponse;
}

export namespace GetFailoverSchemeResponse {
  export type AsObject = {
    failoverScheme?: FailoverScheme.AsObject,
  }
}

export class GetFailoverSchemeYamlRequest extends jspb.Message {
  hasFailoverSchemeRef(): boolean;
  clearFailoverSchemeRef(): void;
  getFailoverSchemeRef(): github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef | undefined;
  setFailoverSchemeRef(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetFailoverSchemeYamlRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetFailoverSchemeYamlRequest): GetFailoverSchemeYamlRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetFailoverSchemeYamlRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetFailoverSchemeYamlRequest;
  static deserializeBinaryFromReader(message: GetFailoverSchemeYamlRequest, reader: jspb.BinaryReader): GetFailoverSchemeYamlRequest;
}

export namespace GetFailoverSchemeYamlRequest {
  export type AsObject = {
    failoverSchemeRef?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef.AsObject,
  }
}

export class GetFailoverSchemeYamlResponse extends jspb.Message {
  hasYamlData(): boolean;
  clearYamlData(): void;
  getYamlData(): github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb.ResourceYaml | undefined;
  setYamlData(value?: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb.ResourceYaml): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetFailoverSchemeYamlResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetFailoverSchemeYamlResponse): GetFailoverSchemeYamlResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetFailoverSchemeYamlResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetFailoverSchemeYamlResponse;
  static deserializeBinaryFromReader(message: GetFailoverSchemeYamlResponse, reader: jspb.BinaryReader): GetFailoverSchemeYamlResponse;
}

export namespace GetFailoverSchemeYamlResponse {
  export type AsObject = {
    yamlData?: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_common_pb.ResourceYaml.AsObject,
  }
}
