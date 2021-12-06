/* eslint-disable */
// package: envoy.config.filter.http.graphql.v2
// file: github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/extensions/graphql/graphql.proto

import * as jspb from "google-protobuf";
import * as validate_validate_pb from "../../../../../../../../../../validate/validate_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_http_uri_pb from "../../../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/config/core/v3/http_uri_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_extension_pb from "../../../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/config/core/v3/extension_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_base_pb from "../../../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/config/core/v3/base_pb";

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

export class ValueProvider extends jspb.Message {
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

  getProviderTemplate(): string;
  setProviderTemplate(value: string): void;

  getProviderCase(): ValueProvider.ProviderCase;
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
    graphqlArg?: ValueProvider.GraphQLArgExtraction.AsObject,
    typedProvider?: ValueProvider.TypedValueProvider.AsObject,
    graphqlParent?: ValueProvider.GraphQLParentExtraction.AsObject,
    providerTemplate: string,
  }

  export class GraphQLArgExtraction extends jspb.Message {
    getArgName(): string;
    setArgName(value: string): void;

    clearPathList(): void;
    getPathList(): Array<PathSegment>;
    setPathList(value: Array<PathSegment>): void;
    addPath(value?: PathSegment, index?: number): PathSegment;

    getRequired(): boolean;
    setRequired(value: boolean): void;

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
      required: boolean,
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

  export enum ProviderCase {
    PROVIDER_NOT_SET = 0,
    GRAPHQL_ARG = 1,
    TYPED_PROVIDER = 2,
    GRAPHQL_PARENT = 3,
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
    spanName: string,
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

export class Query extends jspb.Message {
  hasQuery(): boolean;
  clearQuery(): void;
  getQuery(): github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_base_pb.DataSource | undefined;
  setQuery(value?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_base_pb.DataSource): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Query.AsObject;
  static toObject(includeInstance: boolean, msg: Query): Query.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Query, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Query;
  static deserializeBinaryFromReader(message: Query, reader: jspb.BinaryReader): Query;
}

export namespace Query {
  export type AsObject = {
    query?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_base_pb.DataSource.AsObject,
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
  }
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
  hasSchema(): boolean;
  clearSchema(): void;
  getSchema(): github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_base_pb.DataSource | undefined;
  setSchema(value?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_base_pb.DataSource): void;

  getEnableIntrospection(): boolean;
  setEnableIntrospection(value: boolean): void;

  clearResolutionsList(): void;
  getResolutionsList(): Array<Resolution>;
  setResolutionsList(value: Array<Resolution>): void;
  addResolutions(value?: Resolution, index?: number): Resolution;

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
    schema?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_base_pb.DataSource.AsObject,
    enableIntrospection: boolean,
    resolutionsList: Array<Resolution.AsObject>,
  }
}
