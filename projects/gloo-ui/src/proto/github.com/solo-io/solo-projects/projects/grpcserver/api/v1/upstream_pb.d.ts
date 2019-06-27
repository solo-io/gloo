// package: glooeeapi.solo.io
// file: github.com/solo-io/solo-projects/projects/grpcserver/api/v1/upstream.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../../../../../gogoproto/gogo_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_v1_plugins_aws_aws_pb from "../../../../../../../github.com/solo-io/gloo/projects/gloo/api/v1/plugins/aws/aws_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_v1_plugins_static_static_pb from "../../../../../../../github.com/solo-io/gloo/projects/gloo/api/v1/plugins/static/static_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_v1_plugins_azure_azure_pb from "../../../../../../../github.com/solo-io/gloo/projects/gloo/api/v1/plugins/azure/azure_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_v1_plugins_consul_consul_pb from "../../../../../../../github.com/solo-io/gloo/projects/gloo/api/v1/plugins/consul/consul_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_v1_plugins_kubernetes_kubernetes_pb from "../../../../../../../github.com/solo-io/gloo/projects/gloo/api/v1/plugins/kubernetes/kubernetes_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_v1_upstream_pb from "../../../../../../../github.com/solo-io/gloo/projects/gloo/api/v1/upstream_pb";
import * as github_com_solo_io_solo_kit_api_v1_ref_pb from "../../../../../../../github.com/solo-io/solo-kit/api/v1/ref_pb";

export class GetUpstreamRequest extends jspb.Message {
  hasRef(): boolean;
  clearRef(): void;
  getRef(): github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef | undefined;
  setRef(value?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetUpstreamRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetUpstreamRequest): GetUpstreamRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetUpstreamRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetUpstreamRequest;
  static deserializeBinaryFromReader(message: GetUpstreamRequest, reader: jspb.BinaryReader): GetUpstreamRequest;
}

export namespace GetUpstreamRequest {
  export type AsObject = {
    ref?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject,
  }
}

export class GetUpstreamResponse extends jspb.Message {
  hasUpstream(): boolean;
  clearUpstream(): void;
  getUpstream(): github_com_solo_io_gloo_projects_gloo_api_v1_upstream_pb.Upstream | undefined;
  setUpstream(value?: github_com_solo_io_gloo_projects_gloo_api_v1_upstream_pb.Upstream): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetUpstreamResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetUpstreamResponse): GetUpstreamResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetUpstreamResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetUpstreamResponse;
  static deserializeBinaryFromReader(message: GetUpstreamResponse, reader: jspb.BinaryReader): GetUpstreamResponse;
}

export namespace GetUpstreamResponse {
  export type AsObject = {
    upstream?: github_com_solo_io_gloo_projects_gloo_api_v1_upstream_pb.Upstream.AsObject,
  }
}

export class ListUpstreamsRequest extends jspb.Message {
  clearNamespaceListList(): void;
  getNamespaceListList(): Array<string>;
  setNamespaceListList(value: Array<string>): void;
  addNamespaceList(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListUpstreamsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListUpstreamsRequest): ListUpstreamsRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListUpstreamsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListUpstreamsRequest;
  static deserializeBinaryFromReader(message: ListUpstreamsRequest, reader: jspb.BinaryReader): ListUpstreamsRequest;
}

export namespace ListUpstreamsRequest {
  export type AsObject = {
    namespaceListList: Array<string>,
  }
}

export class ListUpstreamsResponse extends jspb.Message {
  clearUpstreamList(): void;
  getUpstreamList(): Array<github_com_solo_io_gloo_projects_gloo_api_v1_upstream_pb.Upstream>;
  setUpstreamList(value: Array<github_com_solo_io_gloo_projects_gloo_api_v1_upstream_pb.Upstream>): void;
  addUpstream(value?: github_com_solo_io_gloo_projects_gloo_api_v1_upstream_pb.Upstream, index?: number): github_com_solo_io_gloo_projects_gloo_api_v1_upstream_pb.Upstream;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListUpstreamsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListUpstreamsResponse): ListUpstreamsResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListUpstreamsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListUpstreamsResponse;
  static deserializeBinaryFromReader(message: ListUpstreamsResponse, reader: jspb.BinaryReader): ListUpstreamsResponse;
}

export namespace ListUpstreamsResponse {
  export type AsObject = {
    upstreamList: Array<github_com_solo_io_gloo_projects_gloo_api_v1_upstream_pb.Upstream.AsObject>,
  }
}

export class StreamUpstreamListRequest extends jspb.Message {
  getNamespace(): string;
  setNamespace(value: string): void;

  getSelectorMap(): jspb.Map<string, string>;
  clearSelectorMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StreamUpstreamListRequest.AsObject;
  static toObject(includeInstance: boolean, msg: StreamUpstreamListRequest): StreamUpstreamListRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: StreamUpstreamListRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StreamUpstreamListRequest;
  static deserializeBinaryFromReader(message: StreamUpstreamListRequest, reader: jspb.BinaryReader): StreamUpstreamListRequest;
}

export namespace StreamUpstreamListRequest {
  export type AsObject = {
    namespace: string,
    selectorMap: Array<[string, string]>,
  }
}

export class StreamUpstreamListResponse extends jspb.Message {
  hasUpstream(): boolean;
  clearUpstream(): void;
  getUpstream(): github_com_solo_io_gloo_projects_gloo_api_v1_upstream_pb.Upstream | undefined;
  setUpstream(value?: github_com_solo_io_gloo_projects_gloo_api_v1_upstream_pb.Upstream): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StreamUpstreamListResponse.AsObject;
  static toObject(includeInstance: boolean, msg: StreamUpstreamListResponse): StreamUpstreamListResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: StreamUpstreamListResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StreamUpstreamListResponse;
  static deserializeBinaryFromReader(message: StreamUpstreamListResponse, reader: jspb.BinaryReader): StreamUpstreamListResponse;
}

export namespace StreamUpstreamListResponse {
  export type AsObject = {
    upstream?: github_com_solo_io_gloo_projects_gloo_api_v1_upstream_pb.Upstream.AsObject,
  }
}

export class UpstreamInput extends jspb.Message {
  hasRef(): boolean;
  clearRef(): void;
  getRef(): github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef | undefined;
  setRef(value?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef): void;

  hasKube(): boolean;
  clearKube(): void;
  getKube(): github_com_solo_io_gloo_projects_gloo_api_v1_plugins_kubernetes_kubernetes_pb.UpstreamSpec | undefined;
  setKube(value?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_kubernetes_kubernetes_pb.UpstreamSpec): void;

  hasStatic(): boolean;
  clearStatic(): void;
  getStatic(): github_com_solo_io_gloo_projects_gloo_api_v1_plugins_static_static_pb.UpstreamSpec | undefined;
  setStatic(value?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_static_static_pb.UpstreamSpec): void;

  hasAws(): boolean;
  clearAws(): void;
  getAws(): github_com_solo_io_gloo_projects_gloo_api_v1_plugins_aws_aws_pb.UpstreamSpec | undefined;
  setAws(value?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_aws_aws_pb.UpstreamSpec): void;

  hasAzure(): boolean;
  clearAzure(): void;
  getAzure(): github_com_solo_io_gloo_projects_gloo_api_v1_plugins_azure_azure_pb.UpstreamSpec | undefined;
  setAzure(value?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_azure_azure_pb.UpstreamSpec): void;

  hasConsul(): boolean;
  clearConsul(): void;
  getConsul(): github_com_solo_io_gloo_projects_gloo_api_v1_plugins_consul_consul_pb.UpstreamSpec | undefined;
  setConsul(value?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_consul_consul_pb.UpstreamSpec): void;

  getSpecCase(): UpstreamInput.SpecCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpstreamInput.AsObject;
  static toObject(includeInstance: boolean, msg: UpstreamInput): UpstreamInput.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpstreamInput, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpstreamInput;
  static deserializeBinaryFromReader(message: UpstreamInput, reader: jspb.BinaryReader): UpstreamInput;
}

export namespace UpstreamInput {
  export type AsObject = {
    ref?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject,
    kube?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_kubernetes_kubernetes_pb.UpstreamSpec.AsObject,
    pb_static?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_static_static_pb.UpstreamSpec.AsObject,
    aws?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_aws_aws_pb.UpstreamSpec.AsObject,
    azure?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_azure_azure_pb.UpstreamSpec.AsObject,
    consul?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_consul_consul_pb.UpstreamSpec.AsObject,
  }

  export enum SpecCase {
    SPEC_NOT_SET = 0,
    KUBE = 2,
    STATIC = 3,
    AWS = 4,
    AZURE = 5,
    CONSUL = 6,
  }
}

export class CreateUpstreamRequest extends jspb.Message {
  hasInput(): boolean;
  clearInput(): void;
  getInput(): UpstreamInput | undefined;
  setInput(value?: UpstreamInput): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateUpstreamRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateUpstreamRequest): CreateUpstreamRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateUpstreamRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateUpstreamRequest;
  static deserializeBinaryFromReader(message: CreateUpstreamRequest, reader: jspb.BinaryReader): CreateUpstreamRequest;
}

export namespace CreateUpstreamRequest {
  export type AsObject = {
    input?: UpstreamInput.AsObject,
  }
}

export class CreateUpstreamResponse extends jspb.Message {
  hasUpstream(): boolean;
  clearUpstream(): void;
  getUpstream(): github_com_solo_io_gloo_projects_gloo_api_v1_upstream_pb.Upstream | undefined;
  setUpstream(value?: github_com_solo_io_gloo_projects_gloo_api_v1_upstream_pb.Upstream): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateUpstreamResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CreateUpstreamResponse): CreateUpstreamResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateUpstreamResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateUpstreamResponse;
  static deserializeBinaryFromReader(message: CreateUpstreamResponse, reader: jspb.BinaryReader): CreateUpstreamResponse;
}

export namespace CreateUpstreamResponse {
  export type AsObject = {
    upstream?: github_com_solo_io_gloo_projects_gloo_api_v1_upstream_pb.Upstream.AsObject,
  }
}

export class UpdateUpstreamRequest extends jspb.Message {
  hasInput(): boolean;
  clearInput(): void;
  getInput(): UpstreamInput | undefined;
  setInput(value?: UpstreamInput): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateUpstreamRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateUpstreamRequest): UpdateUpstreamRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateUpstreamRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateUpstreamRequest;
  static deserializeBinaryFromReader(message: UpdateUpstreamRequest, reader: jspb.BinaryReader): UpdateUpstreamRequest;
}

export namespace UpdateUpstreamRequest {
  export type AsObject = {
    input?: UpstreamInput.AsObject,
  }
}

export class UpdateUpstreamResponse extends jspb.Message {
  hasUpstream(): boolean;
  clearUpstream(): void;
  getUpstream(): github_com_solo_io_gloo_projects_gloo_api_v1_upstream_pb.Upstream | undefined;
  setUpstream(value?: github_com_solo_io_gloo_projects_gloo_api_v1_upstream_pb.Upstream): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateUpstreamResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateUpstreamResponse): UpdateUpstreamResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateUpstreamResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateUpstreamResponse;
  static deserializeBinaryFromReader(message: UpdateUpstreamResponse, reader: jspb.BinaryReader): UpdateUpstreamResponse;
}

export namespace UpdateUpstreamResponse {
  export type AsObject = {
    upstream?: github_com_solo_io_gloo_projects_gloo_api_v1_upstream_pb.Upstream.AsObject,
  }
}

export class DeleteUpstreamRequest extends jspb.Message {
  hasRef(): boolean;
  clearRef(): void;
  getRef(): github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef | undefined;
  setRef(value?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteUpstreamRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteUpstreamRequest): DeleteUpstreamRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteUpstreamRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteUpstreamRequest;
  static deserializeBinaryFromReader(message: DeleteUpstreamRequest, reader: jspb.BinaryReader): DeleteUpstreamRequest;
}

export namespace DeleteUpstreamRequest {
  export type AsObject = {
    ref?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject,
  }
}

export class DeleteUpstreamResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteUpstreamResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteUpstreamResponse): DeleteUpstreamResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteUpstreamResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteUpstreamResponse;
  static deserializeBinaryFromReader(message: DeleteUpstreamResponse, reader: jspb.BinaryReader): DeleteUpstreamResponse;
}

export namespace DeleteUpstreamResponse {
  export type AsObject = {
  }
}

