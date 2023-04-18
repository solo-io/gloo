/* eslint-disable */
// package: envoy.config.filter.http.graphql.v2
// file: github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/extensions/graphql/graphql.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_duration_pb from "google-protobuf/google/protobuf/duration_pb";
import * as google_protobuf_any_pb from "google-protobuf/google/protobuf/any_pb";
import * as validate_validate_pb from "../../../../../../../../../../validate/validate_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_http_uri_pb from "../../../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/config/core/v3/http_uri_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_extension_pb from "../../../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/config/core/v3/extension_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_base_pb from "../../../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/config/core/v3/base_pb";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";

export class PathSegment extends jspb.Message {
  hasKey(): boolean;
  clearKey(): void;
  getKey(): string;
  setKey(value: string): void;

  hasIndex(): boolean;
  clearIndex(): void;
  getIndex(): number;
  setIndex(value: number): void;

  hasAll(): boolean;
  clearAll(): void;
  getAll(): boolean;
  setAll(value: boolean): void;

  getSegmentCase(): PathSegment.SegmentCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PathSegment.AsObject;
  static toObject(includeInstance: boolean, msg: PathSegment): PathSegment.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: PathSegment, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PathSegment;
  static deserializeBinaryFromReader(message: PathSegment, reader: jspb.BinaryReader): PathSegment;
}

export namespace PathSegment {
  export type AsObject = {
    key: string,
    index: number,
    all: boolean,
  }

  export enum SegmentCase {
    SEGMENT_NOT_SET = 0,
    KEY = 1,
    INDEX = 2,
    ALL = 3,
  }
}

export class Path extends jspb.Message {
  clearSegmentsList(): void;
  getSegmentsList(): Array<PathSegment>;
  setSegmentsList(value: Array<PathSegment>): void;
  addSegments(value?: PathSegment, index?: number): PathSegment;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Path.AsObject;
  static toObject(includeInstance: boolean, msg: Path): Path.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Path, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Path;
  static deserializeBinaryFromReader(message: Path, reader: jspb.BinaryReader): Path;
}

export namespace Path {
  export type AsObject = {
    segmentsList: Array<PathSegment.AsObject>,
  }
}

export class TemplatedPath extends jspb.Message {
  getPathTemplate(): string;
  setPathTemplate(value: string): void;

  getNamedPathsMap(): jspb.Map<string, Path>;
  clearNamedPathsMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TemplatedPath.AsObject;
  static toObject(includeInstance: boolean, msg: TemplatedPath): TemplatedPath.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: TemplatedPath, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TemplatedPath;
  static deserializeBinaryFromReader(message: TemplatedPath, reader: jspb.BinaryReader): TemplatedPath;
}

export namespace TemplatedPath {
  export type AsObject = {
    pathTemplate: string,
    namedPathsMap: Array<[string, Path.AsObject]>,
  }
}

export class ValueProvider extends jspb.Message {
  getProvidersMap(): jspb.Map<string, ValueProvider.Provider>;
  clearProvidersMap(): void;
  getProviderTemplate(): string;
  setProviderTemplate(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ValueProvider.AsObject;
  static toObject(includeInstance: boolean, msg: ValueProvider): ValueProvider.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ValueProvider, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ValueProvider;
  static deserializeBinaryFromReader(message: ValueProvider, reader: jspb.BinaryReader): ValueProvider;
}

export namespace ValueProvider {
  export type AsObject = {
    providersMap: Array<[string, ValueProvider.Provider.AsObject]>,
    providerTemplate: string,
  }

  export class GraphQLArgExtraction extends jspb.Message {
    getArgName(): string;
    setArgName(value: string): void;

    clearPathList(): void;
    getPathList(): Array<PathSegment>;
    setPathList(value: Array<PathSegment>): void;
    addPath(value?: PathSegment, index?: number): PathSegment;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): GraphQLArgExtraction.AsObject;
    static toObject(includeInstance: boolean, msg: GraphQLArgExtraction): GraphQLArgExtraction.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: GraphQLArgExtraction, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): GraphQLArgExtraction;
    static deserializeBinaryFromReader(message: GraphQLArgExtraction, reader: jspb.BinaryReader): GraphQLArgExtraction;
  }

  export namespace GraphQLArgExtraction {
    export type AsObject = {
      argName: string,
      pathList: Array<PathSegment.AsObject>,
    }
  }

  export class GraphQLParentExtraction extends jspb.Message {
    clearPathList(): void;
    getPathList(): Array<PathSegment>;
    setPathList(value: Array<PathSegment>): void;
    addPath(value?: PathSegment, index?: number): PathSegment;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): GraphQLParentExtraction.AsObject;
    static toObject(includeInstance: boolean, msg: GraphQLParentExtraction): GraphQLParentExtraction.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: GraphQLParentExtraction, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): GraphQLParentExtraction;
    static deserializeBinaryFromReader(message: GraphQLParentExtraction, reader: jspb.BinaryReader): GraphQLParentExtraction;
  }

  export namespace GraphQLParentExtraction {
    export type AsObject = {
      pathList: Array<PathSegment.AsObject>,
    }
  }

  export class TypedValueProvider extends jspb.Message {
    getType(): ValueProvider.TypedValueProvider.TypeMap[keyof ValueProvider.TypedValueProvider.TypeMap];
    setType(value: ValueProvider.TypedValueProvider.TypeMap[keyof ValueProvider.TypedValueProvider.TypeMap]): void;

    hasHeader(): boolean;
    clearHeader(): void;
    getHeader(): string;
    setHeader(value: string): void;

    hasValue(): boolean;
    clearValue(): void;
    getValue(): string;
    setValue(value: string): void;

    getValProviderCase(): TypedValueProvider.ValProviderCase;
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): TypedValueProvider.AsObject;
    static toObject(includeInstance: boolean, msg: TypedValueProvider): TypedValueProvider.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: TypedValueProvider, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): TypedValueProvider;
    static deserializeBinaryFromReader(message: TypedValueProvider, reader: jspb.BinaryReader): TypedValueProvider;
  }

  export namespace TypedValueProvider {
    export type AsObject = {
      type: ValueProvider.TypedValueProvider.TypeMap[keyof ValueProvider.TypedValueProvider.TypeMap],
      header: string,
      value: string,
    }

    export interface TypeMap {
      STRING: 0;
      INT: 1;
      FLOAT: 2;
      BOOLEAN: 3;
    }

    export const Type: TypeMap;

    export enum ValProviderCase {
      VAL_PROVIDER_NOT_SET = 0,
      HEADER = 2,
      VALUE = 3,
    }
  }

  export class Provider extends jspb.Message {
    hasGraphqlArg(): boolean;
    clearGraphqlArg(): void;
    getGraphqlArg(): ValueProvider.GraphQLArgExtraction | undefined;
    setGraphqlArg(value?: ValueProvider.GraphQLArgExtraction): void;

    hasTypedProvider(): boolean;
    clearTypedProvider(): void;
    getTypedProvider(): ValueProvider.TypedValueProvider | undefined;
    setTypedProvider(value?: ValueProvider.TypedValueProvider): void;

    hasGraphqlParent(): boolean;
    clearGraphqlParent(): void;
    getGraphqlParent(): ValueProvider.GraphQLParentExtraction | undefined;
    setGraphqlParent(value?: ValueProvider.GraphQLParentExtraction): void;

    getProviderCase(): Provider.ProviderCase;
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Provider.AsObject;
    static toObject(includeInstance: boolean, msg: Provider): Provider.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Provider, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Provider;
    static deserializeBinaryFromReader(message: Provider, reader: jspb.BinaryReader): Provider;
  }

  export namespace Provider {
    export type AsObject = {
      graphqlArg?: ValueProvider.GraphQLArgExtraction.AsObject,
      typedProvider?: ValueProvider.TypedValueProvider.AsObject,
      graphqlParent?: ValueProvider.GraphQLParentExtraction.AsObject,
    }

    export enum ProviderCase {
      PROVIDER_NOT_SET = 0,
      GRAPHQL_ARG = 1,
      TYPED_PROVIDER = 2,
      GRAPHQL_PARENT = 3,
    }
  }
}

export class JsonValueList extends jspb.Message {
  clearValuesList(): void;
  getValuesList(): Array<JsonValue>;
  setValuesList(value: Array<JsonValue>): void;
  addValues(value?: JsonValue, index?: number): JsonValue;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): JsonValueList.AsObject;
  static toObject(includeInstance: boolean, msg: JsonValueList): JsonValueList.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: JsonValueList, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): JsonValueList;
  static deserializeBinaryFromReader(message: JsonValueList, reader: jspb.BinaryReader): JsonValueList;
}

export namespace JsonValueList {
  export type AsObject = {
    valuesList: Array<JsonValue.AsObject>,
  }
}

export class JsonValue extends jspb.Message {
  hasNode(): boolean;
  clearNode(): void;
  getNode(): JsonNode | undefined;
  setNode(value?: JsonNode): void;

  hasValueProvider(): boolean;
  clearValueProvider(): void;
  getValueProvider(): ValueProvider | undefined;
  setValueProvider(value?: ValueProvider): void;

  hasList(): boolean;
  clearList(): void;
  getList(): JsonValueList | undefined;
  setList(value?: JsonValueList): void;

  getJsonValCase(): JsonValue.JsonValCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): JsonValue.AsObject;
  static toObject(includeInstance: boolean, msg: JsonValue): JsonValue.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: JsonValue, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): JsonValue;
  static deserializeBinaryFromReader(message: JsonValue, reader: jspb.BinaryReader): JsonValue;
}

export namespace JsonValue {
  export type AsObject = {
    node?: JsonNode.AsObject,
    valueProvider?: ValueProvider.AsObject,
    list?: JsonValueList.AsObject,
  }

  export enum JsonValCase {
    JSON_VAL_NOT_SET = 0,
    NODE = 1,
    VALUE_PROVIDER = 2,
    LIST = 3,
  }
}

export class JsonKeyValue extends jspb.Message {
  getKey(): string;
  setKey(value: string): void;

  hasValue(): boolean;
  clearValue(): void;
  getValue(): JsonValue | undefined;
  setValue(value?: JsonValue): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): JsonKeyValue.AsObject;
  static toObject(includeInstance: boolean, msg: JsonKeyValue): JsonKeyValue.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: JsonKeyValue, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): JsonKeyValue;
  static deserializeBinaryFromReader(message: JsonKeyValue, reader: jspb.BinaryReader): JsonKeyValue;
}

export namespace JsonKeyValue {
  export type AsObject = {
    key: string,
    value?: JsonValue.AsObject,
  }
}

export class JsonNode extends jspb.Message {
  clearKeyValuesList(): void;
  getKeyValuesList(): Array<JsonKeyValue>;
  setKeyValuesList(value: Array<JsonKeyValue>): void;
  addKeyValues(value?: JsonKeyValue, index?: number): JsonKeyValue;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): JsonNode.AsObject;
  static toObject(includeInstance: boolean, msg: JsonNode): JsonNode.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: JsonNode, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): JsonNode;
  static deserializeBinaryFromReader(message: JsonNode, reader: jspb.BinaryReader): JsonNode;
}

export namespace JsonNode {
  export type AsObject = {
    keyValuesList: Array<JsonKeyValue.AsObject>,
  }
}

export class RequestTemplate extends jspb.Message {
  getHeadersMap(): jspb.Map<string, ValueProvider>;
  clearHeadersMap(): void;
  getQueryParamsMap(): jspb.Map<string, ValueProvider>;
  clearQueryParamsMap(): void;
  hasOutgoingBody(): boolean;
  clearOutgoingBody(): void;
  getOutgoingBody(): JsonValue | undefined;
  setOutgoingBody(value?: JsonValue): void;

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
    headersMap: Array<[string, ValueProvider.AsObject]>,
    queryParamsMap: Array<[string, ValueProvider.AsObject]>,
    outgoingBody?: JsonValue.AsObject,
  }
}

export class ResponseTemplate extends jspb.Message {
  clearResultRootList(): void;
  getResultRootList(): Array<PathSegment>;
  setResultRootList(value: Array<PathSegment>): void;
  addResultRoot(value?: PathSegment, index?: number): PathSegment;

  getSettersMap(): jspb.Map<string, TemplatedPath>;
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
    resultRootList: Array<PathSegment.AsObject>,
    settersMap: Array<[string, TemplatedPath.AsObject]>,
  }
}

export class RESTResolver extends jspb.Message {
  hasServerUri(): boolean;
  clearServerUri(): void;
  getServerUri(): github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_http_uri_pb.HttpUri | undefined;
  setServerUri(value?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_http_uri_pb.HttpUri): void;

  hasRequestTransform(): boolean;
  clearRequestTransform(): void;
  getRequestTransform(): RequestTemplate | undefined;
  setRequestTransform(value?: RequestTemplate): void;

  hasPreExecutionTransform(): boolean;
  clearPreExecutionTransform(): void;
  getPreExecutionTransform(): ResponseTemplate | undefined;
  setPreExecutionTransform(value?: ResponseTemplate): void;

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
    serverUri?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_http_uri_pb.HttpUri.AsObject,
    requestTransform?: RequestTemplate.AsObject,
    preExecutionTransform?: ResponseTemplate.AsObject,
    spanName: string,
  }
}

export class GrpcRequestTemplate extends jspb.Message {
  hasOutgoingMessageJson(): boolean;
  clearOutgoingMessageJson(): void;
  getOutgoingMessageJson(): JsonValue | undefined;
  setOutgoingMessageJson(value?: JsonValue): void;

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
    outgoingMessageJson?: JsonValue.AsObject,
    serviceName: string,
    methodName: string,
    requestMetadataMap: Array<[string, string]>,
  }
}

export class GrpcDescriptorRegistry extends jspb.Message {
  hasProtoDescriptors(): boolean;
  clearProtoDescriptors(): void;
  getProtoDescriptors(): github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_base_pb.DataSource | undefined;
  setProtoDescriptors(value?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_base_pb.DataSource): void;

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
    protoDescriptors?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_base_pb.DataSource.AsObject,
  }
}

export class GrpcResolver extends jspb.Message {
  hasServerUri(): boolean;
  clearServerUri(): void;
  getServerUri(): github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_http_uri_pb.HttpUri | undefined;
  setServerUri(value?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_http_uri_pb.HttpUri): void;

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
    serverUri?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_http_uri_pb.HttpUri.AsObject,
    requestTransform?: GrpcRequestTemplate.AsObject,
    spanName: string,
  }
}

export class StaticResolver extends jspb.Message {
  hasSyncResponse(): boolean;
  clearSyncResponse(): void;
  getSyncResponse(): string;
  setSyncResponse(value: string): void;

  hasAsyncResponse(): boolean;
  clearAsyncResponse(): void;
  getAsyncResponse(): StaticResolver.AsyncResponse | undefined;
  setAsyncResponse(value?: StaticResolver.AsyncResponse): void;

  hasErrorResponse(): boolean;
  clearErrorResponse(): void;
  getErrorResponse(): string;
  setErrorResponse(value: string): void;

  getResponseCase(): StaticResolver.ResponseCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StaticResolver.AsObject;
  static toObject(includeInstance: boolean, msg: StaticResolver): StaticResolver.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: StaticResolver, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StaticResolver;
  static deserializeBinaryFromReader(message: StaticResolver, reader: jspb.BinaryReader): StaticResolver;
}

export namespace StaticResolver {
  export type AsObject = {
    syncResponse: string,
    asyncResponse?: StaticResolver.AsyncResponse.AsObject,
    errorResponse: string,
  }

  export class AsyncResponse extends jspb.Message {
    getResponse(): string;
    setResponse(value: string): void;

    getDelayMs(): number;
    setDelayMs(value: number): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): AsyncResponse.AsObject;
    static toObject(includeInstance: boolean, msg: AsyncResponse): AsyncResponse.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: AsyncResponse, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): AsyncResponse;
    static deserializeBinaryFromReader(message: AsyncResponse, reader: jspb.BinaryReader): AsyncResponse;
  }

  export namespace AsyncResponse {
    export type AsObject = {
      response: string,
      delayMs: number,
    }
  }

  export enum ResponseCase {
    RESPONSE_NOT_SET = 0,
    SYNC_RESPONSE = 1,
    ASYNC_RESPONSE = 2,
    ERROR_RESPONSE = 3,
  }
}

export class AbstractTypeResolver extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AbstractTypeResolver.AsObject;
  static toObject(includeInstance: boolean, msg: AbstractTypeResolver): AbstractTypeResolver.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: AbstractTypeResolver, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AbstractTypeResolver;
  static deserializeBinaryFromReader(message: AbstractTypeResolver, reader: jspb.BinaryReader): AbstractTypeResolver;
}

export namespace AbstractTypeResolver {
  export type AsObject = {
  }
}

export class QueryMatcher extends jspb.Message {
  hasFieldMatcher(): boolean;
  clearFieldMatcher(): void;
  getFieldMatcher(): QueryMatcher.FieldMatcher | undefined;
  setFieldMatcher(value?: QueryMatcher.FieldMatcher): void;

  getMatchCase(): QueryMatcher.MatchCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): QueryMatcher.AsObject;
  static toObject(includeInstance: boolean, msg: QueryMatcher): QueryMatcher.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: QueryMatcher, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): QueryMatcher;
  static deserializeBinaryFromReader(message: QueryMatcher, reader: jspb.BinaryReader): QueryMatcher;
}

export namespace QueryMatcher {
  export type AsObject = {
    fieldMatcher?: QueryMatcher.FieldMatcher.AsObject,
  }

  export class FieldMatcher extends jspb.Message {
    getType(): string;
    setType(value: string): void;

    getField(): string;
    setField(value: string): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): FieldMatcher.AsObject;
    static toObject(includeInstance: boolean, msg: FieldMatcher): FieldMatcher.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: FieldMatcher, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): FieldMatcher;
    static deserializeBinaryFromReader(message: FieldMatcher, reader: jspb.BinaryReader): FieldMatcher;
  }

  export namespace FieldMatcher {
    export type AsObject = {
      type: string,
      field: string,
    }
  }

  export enum MatchCase {
    MATCH_NOT_SET = 0,
    FIELD_MATCHER = 1,
  }
}

export class Resolution extends jspb.Message {
  hasMatcher(): boolean;
  clearMatcher(): void;
  getMatcher(): QueryMatcher | undefined;
  setMatcher(value?: QueryMatcher): void;

  hasResolver(): boolean;
  clearResolver(): void;
  getResolver(): github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_extension_pb.TypedExtensionConfig | undefined;
  setResolver(value?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_extension_pb.TypedExtensionConfig): void;

  getStatPrefix(): string;
  setStatPrefix(value: string): void;

  hasCacheControl(): boolean;
  clearCacheControl(): void;
  getCacheControl(): CacheControl | undefined;
  setCacheControl(value?: CacheControl): void;

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
    matcher?: QueryMatcher.AsObject,
    resolver?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_extension_pb.TypedExtensionConfig.AsObject,
    statPrefix: string,
    cacheControl?: CacheControl.AsObject,
  }
}

export class CacheControl extends jspb.Message {
  hasMaxAge(): boolean;
  clearMaxAge(): void;
  getMaxAge(): google_protobuf_wrappers_pb.UInt32Value | undefined;
  setMaxAge(value?: google_protobuf_wrappers_pb.UInt32Value): void;

  getScope(): CacheControl.CacheControlScopeMap[keyof CacheControl.CacheControlScopeMap];
  setScope(value: CacheControl.CacheControlScopeMap[keyof CacheControl.CacheControlScopeMap]): void;

  getInheritMaxAge(): boolean;
  setInheritMaxAge(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CacheControl.AsObject;
  static toObject(includeInstance: boolean, msg: CacheControl): CacheControl.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CacheControl, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CacheControl;
  static deserializeBinaryFromReader(message: CacheControl, reader: jspb.BinaryReader): CacheControl;
}

export namespace CacheControl {
  export type AsObject = {
    maxAge?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
    scope: CacheControl.CacheControlScopeMap[keyof CacheControl.CacheControlScopeMap],
    inheritMaxAge: boolean,
  }

  export interface CacheControlScopeMap {
    UNSET: 0;
    PUBLIC: 1;
    PRIVATE: 2;
  }

  export const CacheControlScope: CacheControlScopeMap;
}

export class GraphQLConfig extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GraphQLConfig.AsObject;
  static toObject(includeInstance: boolean, msg: GraphQLConfig): GraphQLConfig.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GraphQLConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GraphQLConfig;
  static deserializeBinaryFromReader(message: GraphQLConfig, reader: jspb.BinaryReader): GraphQLConfig;
}

export namespace GraphQLConfig {
  export type AsObject = {
  }
}

export class GraphQLRouteConfig extends jspb.Message {
  hasExecutableSchema(): boolean;
  clearExecutableSchema(): void;
  getExecutableSchema(): ExecutableSchema | undefined;
  setExecutableSchema(value?: ExecutableSchema): void;

  getStatPrefix(): string;
  setStatPrefix(value: string): void;

  hasPersistedQueryCacheConfig(): boolean;
  clearPersistedQueryCacheConfig(): void;
  getPersistedQueryCacheConfig(): PersistedQueryCacheConfig | undefined;
  setPersistedQueryCacheConfig(value?: PersistedQueryCacheConfig): void;

  clearAllowedQueryHashesList(): void;
  getAllowedQueryHashesList(): Array<string>;
  setAllowedQueryHashesList(value: Array<string>): void;
  addAllowedQueryHashes(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GraphQLRouteConfig.AsObject;
  static toObject(includeInstance: boolean, msg: GraphQLRouteConfig): GraphQLRouteConfig.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GraphQLRouteConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GraphQLRouteConfig;
  static deserializeBinaryFromReader(message: GraphQLRouteConfig, reader: jspb.BinaryReader): GraphQLRouteConfig;
}

export namespace GraphQLRouteConfig {
  export type AsObject = {
    executableSchema?: ExecutableSchema.AsObject,
    statPrefix: string,
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
  hasSchemaDefinition(): boolean;
  clearSchemaDefinition(): void;
  getSchemaDefinition(): github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_base_pb.DataSource | undefined;
  setSchemaDefinition(value?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_base_pb.DataSource): void;

  hasExecutor(): boolean;
  clearExecutor(): void;
  getExecutor(): Executor | undefined;
  setExecutor(value?: Executor): void;

  getExtensionsMap(): jspb.Map<string, google_protobuf_any_pb.Any>;
  clearExtensionsMap(): void;
  getLogRequestResponseInfo(): boolean;
  setLogRequestResponseInfo(value: boolean): void;

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
    schemaDefinition?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_base_pb.DataSource.AsObject,
    executor?: Executor.AsObject,
    extensionsMap: Array<[string, google_protobuf_any_pb.Any.AsObject]>,
    logRequestResponseInfo: boolean,
  }
}

export class Executor extends jspb.Message {
  hasLocal(): boolean;
  clearLocal(): void;
  getLocal(): Executor.Local | undefined;
  setLocal(value?: Executor.Local): void;

  hasRemote(): boolean;
  clearRemote(): void;
  getRemote(): Executor.Remote | undefined;
  setRemote(value?: Executor.Remote): void;

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
    remote?: Executor.Remote.AsObject,
  }

  export class Local extends jspb.Message {
    clearResolutionsList(): void;
    getResolutionsList(): Array<Resolution>;
    setResolutionsList(value: Array<Resolution>): void;
    addResolutions(value?: Resolution, index?: number): Resolution;

    getEnableIntrospection(): boolean;
    setEnableIntrospection(value: boolean): void;

    getMaxDepth(): number;
    setMaxDepth(value: number): void;

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
      resolutionsList: Array<Resolution.AsObject>,
      enableIntrospection: boolean,
      maxDepth: number,
    }
  }

  export class Remote extends jspb.Message {
    hasServerUri(): boolean;
    clearServerUri(): void;
    getServerUri(): github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_http_uri_pb.HttpUri | undefined;
    setServerUri(value?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_http_uri_pb.HttpUri): void;

    hasRequest(): boolean;
    clearRequest(): void;
    getRequest(): Executor.Remote.RemoteSchemaRequest | undefined;
    setRequest(value?: Executor.Remote.RemoteSchemaRequest): void;

    getSpanName(): string;
    setSpanName(value: string): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Remote.AsObject;
    static toObject(includeInstance: boolean, msg: Remote): Remote.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Remote, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Remote;
    static deserializeBinaryFromReader(message: Remote, reader: jspb.BinaryReader): Remote;
  }

  export namespace Remote {
    export type AsObject = {
      serverUri?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_http_uri_pb.HttpUri.AsObject,
      request?: Executor.Remote.RemoteSchemaRequest.AsObject,
      spanName: string,
    }

    export class Extraction extends jspb.Message {
      hasValue(): boolean;
      clearValue(): void;
      getValue(): string;
      setValue(value: string): void;

      hasHeader(): boolean;
      clearHeader(): void;
      getHeader(): string;
      setHeader(value: string): void;

      hasDynamicMetadata(): boolean;
      clearDynamicMetadata(): void;
      getDynamicMetadata(): Executor.Remote.Extraction.DynamicMetadataExtraction | undefined;
      setDynamicMetadata(value?: Executor.Remote.Extraction.DynamicMetadataExtraction): void;

      getExtractionTypeCase(): Extraction.ExtractionTypeCase;
      serializeBinary(): Uint8Array;
      toObject(includeInstance?: boolean): Extraction.AsObject;
      static toObject(includeInstance: boolean, msg: Extraction): Extraction.AsObject;
      static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
      static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
      static serializeBinaryToWriter(message: Extraction, writer: jspb.BinaryWriter): void;
      static deserializeBinary(bytes: Uint8Array): Extraction;
      static deserializeBinaryFromReader(message: Extraction, reader: jspb.BinaryReader): Extraction;
    }

    export namespace Extraction {
      export type AsObject = {
        value: string,
        header: string,
        dynamicMetadata?: Executor.Remote.Extraction.DynamicMetadataExtraction.AsObject,
      }

      export class DynamicMetadataExtraction extends jspb.Message {
        getMetadataNamespace(): string;
        setMetadataNamespace(value: string): void;

        getKey(): string;
        setKey(value: string): void;

        serializeBinary(): Uint8Array;
        toObject(includeInstance?: boolean): DynamicMetadataExtraction.AsObject;
        static toObject(includeInstance: boolean, msg: DynamicMetadataExtraction): DynamicMetadataExtraction.AsObject;
        static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
        static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
        static serializeBinaryToWriter(message: DynamicMetadataExtraction, writer: jspb.BinaryWriter): void;
        static deserializeBinary(bytes: Uint8Array): DynamicMetadataExtraction;
        static deserializeBinaryFromReader(message: DynamicMetadataExtraction, reader: jspb.BinaryReader): DynamicMetadataExtraction;
      }

      export namespace DynamicMetadataExtraction {
        export type AsObject = {
          metadataNamespace: string,
          key: string,
        }
      }

      export enum ExtractionTypeCase {
        EXTRACTION_TYPE_NOT_SET = 0,
        VALUE = 1,
        HEADER = 2,
        DYNAMIC_METADATA = 3,
      }
    }

    export class RemoteSchemaRequest extends jspb.Message {
      getHeadersMap(): jspb.Map<string, Executor.Remote.Extraction>;
      clearHeadersMap(): void;
      getQueryParamsMap(): jspb.Map<string, Executor.Remote.Extraction>;
      clearQueryParamsMap(): void;
      serializeBinary(): Uint8Array;
      toObject(includeInstance?: boolean): RemoteSchemaRequest.AsObject;
      static toObject(includeInstance: boolean, msg: RemoteSchemaRequest): RemoteSchemaRequest.AsObject;
      static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
      static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
      static serializeBinaryToWriter(message: RemoteSchemaRequest, writer: jspb.BinaryWriter): void;
      static deserializeBinary(bytes: Uint8Array): RemoteSchemaRequest;
      static deserializeBinaryFromReader(message: RemoteSchemaRequest, reader: jspb.BinaryReader): RemoteSchemaRequest;
    }

    export namespace RemoteSchemaRequest {
      export type AsObject = {
        headersMap: Array<[string, Executor.Remote.Extraction.AsObject]>,
        queryParamsMap: Array<[string, Executor.Remote.Extraction.AsObject]>,
      }
    }
  }

  export enum ExecutorCase {
    EXECUTOR_NOT_SET = 0,
    LOCAL = 1,
    REMOTE = 2,
  }
}
