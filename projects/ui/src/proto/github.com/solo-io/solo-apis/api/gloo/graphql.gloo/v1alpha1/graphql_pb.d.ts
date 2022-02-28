/* eslint-disable */
// package: graphql.gloo.solo.io
// file: github.com/solo-io/solo-apis/api/gloo/graphql.gloo/v1alpha1/graphql.proto

import * as jspb from "google-protobuf";
import * as extproto_ext_pb from "../../../../../../../extproto/ext_pb";
import * as github_com_solo_io_solo_kit_api_v1_solo_kit_pb from "../../../../../../../github.com/solo-io/solo-kit/api/v1/solo-kit_pb";
import * as github_com_solo_io_solo_kit_api_v1_ref_pb from "../../../../../../../github.com/solo-io/solo-kit/api/v1/ref_pb";
import * as validate_validate_pb from "../../../../../../../validate/validate_pb";
import * as google_protobuf_empty_pb from "google-protobuf/google/protobuf/empty_pb";
import * as google_protobuf_struct_pb from "google-protobuf/google/protobuf/struct_pb";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";

export class RequestTemplate extends jspb.Message {
  getHeadersMap(): jspb.Map<string, string>;
  clearHeadersMap(): void;
  getQueryParamsMap(): jspb.Map<string, string>;
  clearQueryParamsMap(): void;
  hasBody(): boolean;
  clearBody(): void;
  getBody(): google_protobuf_struct_pb.Value | undefined;
  setBody(value?: google_protobuf_struct_pb.Value): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RequestTemplate.AsObject;
  static toObject(includeInstance: boolean, msg: RequestTemplate): RequestTemplate.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RequestTemplate, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RequestTemplate;
  static deserializeBinaryFromReader(message: RequestTemplate, reader: jspb.BinaryReader): RequestTemplate;
}

export namespace RequestTemplate {
  export type AsObject = {
    headersMap: Array<[string, string]>,
    queryParamsMap: Array<[string, string]>,
    body?: google_protobuf_struct_pb.Value.AsObject,
  }
}

export class ResponseTemplate extends jspb.Message {
  getResultRoot(): string;
  setResultRoot(value: string): void;

  getSettersMap(): jspb.Map<string, string>;
  clearSettersMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResponseTemplate.AsObject;
  static toObject(includeInstance: boolean, msg: ResponseTemplate): ResponseTemplate.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ResponseTemplate, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResponseTemplate;
  static deserializeBinaryFromReader(message: ResponseTemplate, reader: jspb.BinaryReader): ResponseTemplate;
}

export namespace ResponseTemplate {
  export type AsObject = {
    resultRoot: string,
    settersMap: Array<[string, string]>,
  }
}

export class GrpcRequestTemplate extends jspb.Message {
  hasOutgoingMessageJson(): boolean;
  clearOutgoingMessageJson(): void;
  getOutgoingMessageJson(): google_protobuf_struct_pb.Value | undefined;
  setOutgoingMessageJson(value?: google_protobuf_struct_pb.Value): void;

  getServiceName(): string;
  setServiceName(value: string): void;

  getMethodName(): string;
  setMethodName(value: string): void;

  getRequestMetadataMap(): jspb.Map<string, string>;
  clearRequestMetadataMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GrpcRequestTemplate.AsObject;
  static toObject(includeInstance: boolean, msg: GrpcRequestTemplate): GrpcRequestTemplate.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GrpcRequestTemplate, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GrpcRequestTemplate;
  static deserializeBinaryFromReader(message: GrpcRequestTemplate, reader: jspb.BinaryReader): GrpcRequestTemplate;
}

export namespace GrpcRequestTemplate {
  export type AsObject = {
    outgoingMessageJson?: google_protobuf_struct_pb.Value.AsObject,
    serviceName: string,
    methodName: string,
    requestMetadataMap: Array<[string, string]>,
  }
}

export class RESTResolver extends jspb.Message {
  hasUpstreamRef(): boolean;
  clearUpstreamRef(): void;
  getUpstreamRef(): github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef | undefined;
  setUpstreamRef(value?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef): void;

  hasRequest(): boolean;
  clearRequest(): void;
  getRequest(): RequestTemplate | undefined;
  setRequest(value?: RequestTemplate): void;

  hasResponse(): boolean;
  clearResponse(): void;
  getResponse(): ResponseTemplate | undefined;
  setResponse(value?: ResponseTemplate): void;

  getSpanName(): string;
  setSpanName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RESTResolver.AsObject;
  static toObject(includeInstance: boolean, msg: RESTResolver): RESTResolver.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RESTResolver, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RESTResolver;
  static deserializeBinaryFromReader(message: RESTResolver, reader: jspb.BinaryReader): RESTResolver;
}

export namespace RESTResolver {
  export type AsObject = {
    upstreamRef?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject,
    request?: RequestTemplate.AsObject,
    response?: ResponseTemplate.AsObject,
    spanName: string,
  }
}

export class GrpcDescriptorRegistry extends jspb.Message {
  hasProtoDescriptor(): boolean;
  clearProtoDescriptor(): void;
  getProtoDescriptor(): string;
  setProtoDescriptor(value: string): void;

  hasProtoDescriptorBin(): boolean;
  clearProtoDescriptorBin(): void;
  getProtoDescriptorBin(): Uint8Array | string;
  getProtoDescriptorBin_asU8(): Uint8Array;
  getProtoDescriptorBin_asB64(): string;
  setProtoDescriptorBin(value: Uint8Array | string): void;

  getDescriptorSetCase(): GrpcDescriptorRegistry.DescriptorSetCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GrpcDescriptorRegistry.AsObject;
  static toObject(includeInstance: boolean, msg: GrpcDescriptorRegistry): GrpcDescriptorRegistry.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GrpcDescriptorRegistry, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GrpcDescriptorRegistry;
  static deserializeBinaryFromReader(message: GrpcDescriptorRegistry, reader: jspb.BinaryReader): GrpcDescriptorRegistry;
}

export namespace GrpcDescriptorRegistry {
  export type AsObject = {
    protoDescriptor: string,
    protoDescriptorBin: Uint8Array | string,
  }

  export enum DescriptorSetCase {
    DESCRIPTOR_SET_NOT_SET = 0,
    PROTO_DESCRIPTOR = 1,
    PROTO_DESCRIPTOR_BIN = 2,
  }
}

export class GrpcResolver extends jspb.Message {
  hasUpstreamRef(): boolean;
  clearUpstreamRef(): void;
  getUpstreamRef(): github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef | undefined;
  setUpstreamRef(value?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef): void;

  hasRequestTransform(): boolean;
  clearRequestTransform(): void;
  getRequestTransform(): GrpcRequestTemplate | undefined;
  setRequestTransform(value?: GrpcRequestTemplate): void;

  getSpanName(): string;
  setSpanName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GrpcResolver.AsObject;
  static toObject(includeInstance: boolean, msg: GrpcResolver): GrpcResolver.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GrpcResolver, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GrpcResolver;
  static deserializeBinaryFromReader(message: GrpcResolver, reader: jspb.BinaryReader): GrpcResolver;
}

export namespace GrpcResolver {
  export type AsObject = {
    upstreamRef?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject,
    requestTransform?: GrpcRequestTemplate.AsObject,
    spanName: string,
  }
}

export class Resolution extends jspb.Message {
  hasRestResolver(): boolean;
  clearRestResolver(): void;
  getRestResolver(): RESTResolver | undefined;
  setRestResolver(value?: RESTResolver): void;

  hasGrpcResolver(): boolean;
  clearGrpcResolver(): void;
  getGrpcResolver(): GrpcResolver | undefined;
  setGrpcResolver(value?: GrpcResolver): void;

  hasStatPrefix(): boolean;
  clearStatPrefix(): void;
  getStatPrefix(): google_protobuf_wrappers_pb.StringValue | undefined;
  setStatPrefix(value?: google_protobuf_wrappers_pb.StringValue): void;

  getResolverCase(): Resolution.ResolverCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Resolution.AsObject;
  static toObject(includeInstance: boolean, msg: Resolution): Resolution.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Resolution, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Resolution;
  static deserializeBinaryFromReader(message: Resolution, reader: jspb.BinaryReader): Resolution;
}

export namespace Resolution {
  export type AsObject = {
    restResolver?: RESTResolver.AsObject,
    grpcResolver?: GrpcResolver.AsObject,
    statPrefix?: google_protobuf_wrappers_pb.StringValue.AsObject,
  }

  export enum ResolverCase {
    RESOLVER_NOT_SET = 0,
    REST_RESOLVER = 1,
    GRPC_RESOLVER = 2,
  }
}

export class GraphQLSchemaSpec extends jspb.Message {
  hasExecutableSchema(): boolean;
  clearExecutableSchema(): void;
  getExecutableSchema(): ExecutableSchema | undefined;
  setExecutableSchema(value?: ExecutableSchema): void;

  hasStatPrefix(): boolean;
  clearStatPrefix(): void;
  getStatPrefix(): google_protobuf_wrappers_pb.StringValue | undefined;
  setStatPrefix(value?: google_protobuf_wrappers_pb.StringValue): void;

  hasPersistedQueryCacheConfig(): boolean;
  clearPersistedQueryCacheConfig(): void;
  getPersistedQueryCacheConfig(): PersistedQueryCacheConfig | undefined;
  setPersistedQueryCacheConfig(value?: PersistedQueryCacheConfig): void;

  clearAllowedQueryHashesList(): void;
  getAllowedQueryHashesList(): Array<string>;
  setAllowedQueryHashesList(value: Array<string>): void;
  addAllowedQueryHashes(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GraphQLSchemaSpec.AsObject;
  static toObject(includeInstance: boolean, msg: GraphQLSchemaSpec): GraphQLSchemaSpec.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GraphQLSchemaSpec, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GraphQLSchemaSpec;
  static deserializeBinaryFromReader(message: GraphQLSchemaSpec, reader: jspb.BinaryReader): GraphQLSchemaSpec;
}

export namespace GraphQLSchemaSpec {
  export type AsObject = {
    executableSchema?: ExecutableSchema.AsObject,
    statPrefix?: google_protobuf_wrappers_pb.StringValue.AsObject,
    persistedQueryCacheConfig?: PersistedQueryCacheConfig.AsObject,
    allowedQueryHashesList: Array<string>,
  }
}

export class PersistedQueryCacheConfig extends jspb.Message {
  getCacheSize(): number;
  setCacheSize(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PersistedQueryCacheConfig.AsObject;
  static toObject(includeInstance: boolean, msg: PersistedQueryCacheConfig): PersistedQueryCacheConfig.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: PersistedQueryCacheConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PersistedQueryCacheConfig;
  static deserializeBinaryFromReader(message: PersistedQueryCacheConfig, reader: jspb.BinaryReader): PersistedQueryCacheConfig;
}

export namespace PersistedQueryCacheConfig {
  export type AsObject = {
    cacheSize: number,
  }
}

export class ExecutableSchema extends jspb.Message {
  getSchemaDefinition(): string;
  setSchemaDefinition(value: string): void;

  hasExecutor(): boolean;
  clearExecutor(): void;
  getExecutor(): Executor | undefined;
  setExecutor(value?: Executor): void;

  hasGrpcDescriptorRegistry(): boolean;
  clearGrpcDescriptorRegistry(): void;
  getGrpcDescriptorRegistry(): GrpcDescriptorRegistry | undefined;
  setGrpcDescriptorRegistry(value?: GrpcDescriptorRegistry): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ExecutableSchema.AsObject;
  static toObject(includeInstance: boolean, msg: ExecutableSchema): ExecutableSchema.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ExecutableSchema, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ExecutableSchema;
  static deserializeBinaryFromReader(message: ExecutableSchema, reader: jspb.BinaryReader): ExecutableSchema;
}

export namespace ExecutableSchema {
  export type AsObject = {
    schemaDefinition: string,
    executor?: Executor.AsObject,
    grpcDescriptorRegistry?: GrpcDescriptorRegistry.AsObject,
  }
}

export class Executor extends jspb.Message {
  hasLocal(): boolean;
  clearLocal(): void;
  getLocal(): Executor.Local | undefined;
  setLocal(value?: Executor.Local): void;

  getExecutorCase(): Executor.ExecutorCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Executor.AsObject;
  static toObject(includeInstance: boolean, msg: Executor): Executor.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Executor, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Executor;
  static deserializeBinaryFromReader(message: Executor, reader: jspb.BinaryReader): Executor;
}

export namespace Executor {
  export type AsObject = {
    local?: Executor.Local.AsObject,
  }

  export class Local extends jspb.Message {
    getResolutionsMap(): jspb.Map<string, Resolution>;
    clearResolutionsMap(): void;
    getEnableIntrospection(): boolean;
    setEnableIntrospection(value: boolean): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Local.AsObject;
    static toObject(includeInstance: boolean, msg: Local): Local.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Local, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Local;
    static deserializeBinaryFromReader(message: Local, reader: jspb.BinaryReader): Local;
  }

  export namespace Local {
    export type AsObject = {
      resolutionsMap: Array<[string, Resolution.AsObject]>,
      enableIntrospection: boolean,
    }
  }

  export enum ExecutorCase {
    EXECUTOR_NOT_SET = 0,
    LOCAL = 1,
  }
}

export class GraphQLSchemaStatus extends jspb.Message {
  getState(): GraphQLSchemaStatus.StateMap[keyof GraphQLSchemaStatus.StateMap];
  setState(value: GraphQLSchemaStatus.StateMap[keyof GraphQLSchemaStatus.StateMap]): void;

  getReason(): string;
  setReason(value: string): void;

  getReportedBy(): string;
  setReportedBy(value: string): void;

  getSubresourceStatusesMap(): jspb.Map<string, GraphQLSchemaStatus>;
  clearSubresourceStatusesMap(): void;
  hasDetails(): boolean;
  clearDetails(): void;
  getDetails(): google_protobuf_struct_pb.Struct | undefined;
  setDetails(value?: google_protobuf_struct_pb.Struct): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GraphQLSchemaStatus.AsObject;
  static toObject(includeInstance: boolean, msg: GraphQLSchemaStatus): GraphQLSchemaStatus.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GraphQLSchemaStatus, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GraphQLSchemaStatus;
  static deserializeBinaryFromReader(message: GraphQLSchemaStatus, reader: jspb.BinaryReader): GraphQLSchemaStatus;
}

export namespace GraphQLSchemaStatus {
  export type AsObject = {
    state: GraphQLSchemaStatus.StateMap[keyof GraphQLSchemaStatus.StateMap],
    reason: string,
    reportedBy: string,
    subresourceStatusesMap: Array<[string, GraphQLSchemaStatus.AsObject]>,
    details?: google_protobuf_struct_pb.Struct.AsObject,
  }

  export interface StateMap {
    PENDING: 0;
    ACCEPTED: 1;
    REJECTED: 2;
    WARNING: 3;
  }

  export const State: StateMap;
}

export class GraphQLSchemaNamespacedStatuses extends jspb.Message {
  getStatusesMap(): jspb.Map<string, GraphQLSchemaStatus>;
  clearStatusesMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GraphQLSchemaNamespacedStatuses.AsObject;
  static toObject(includeInstance: boolean, msg: GraphQLSchemaNamespacedStatuses): GraphQLSchemaNamespacedStatuses.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GraphQLSchemaNamespacedStatuses, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GraphQLSchemaNamespacedStatuses;
  static deserializeBinaryFromReader(message: GraphQLSchemaNamespacedStatuses, reader: jspb.BinaryReader): GraphQLSchemaNamespacedStatuses;
}

export namespace GraphQLSchemaNamespacedStatuses {
  export type AsObject = {
    statusesMap: Array<[string, GraphQLSchemaStatus.AsObject]>,
  }
}
