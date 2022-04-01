/* eslint-disable */
// package: rpc.edge.gloo.solo.io
// file: github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/graphql.proto

import * as jspb from "google-protobuf";
import * as extproto_ext_pb from "../../../../../../../../extproto/ext_pb";
import * as github_com_solo_io_skv2_api_core_v1_core_pb from "../../../../../../../../github.com/solo-io/skv2/api/core/v1/core_pb";
import * as github_com_solo_io_solo_apis_api_gloo_graphql_gloo_v1alpha1_graphql_pb from "../../../../../../../../github.com/solo-io/solo-apis/api/gloo/graphql.gloo/v1alpha1/graphql_pb";
import * as github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb from "../../../../../../../../github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/common_pb";

export class GraphqlApi extends jspb.Message {
  hasMetadata(): boolean;
  clearMetadata(): void;
  getMetadata(): github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ObjectMeta | undefined;
  setMetadata(value?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ObjectMeta): void;

  hasSpec(): boolean;
  clearSpec(): void;
  getSpec(): github_com_solo_io_solo_apis_api_gloo_graphql_gloo_v1alpha1_graphql_pb.GraphQLApiSpec | undefined;
  setSpec(value?: github_com_solo_io_solo_apis_api_gloo_graphql_gloo_v1alpha1_graphql_pb.GraphQLApiSpec): void;

  hasStatus(): boolean;
  clearStatus(): void;
  getStatus(): github_com_solo_io_solo_apis_api_gloo_graphql_gloo_v1alpha1_graphql_pb.GraphQLApiStatus | undefined;
  setStatus(value?: github_com_solo_io_solo_apis_api_gloo_graphql_gloo_v1alpha1_graphql_pb.GraphQLApiStatus): void;

  hasGlooInstance(): boolean;
  clearGlooInstance(): void;
  getGlooInstance(): github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef | undefined;
  setGlooInstance(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GraphqlApi.AsObject;
  static toObject(includeInstance: boolean, msg: GraphqlApi): GraphqlApi.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GraphqlApi, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GraphqlApi;
  static deserializeBinaryFromReader(message: GraphqlApi, reader: jspb.BinaryReader): GraphqlApi;
}

export namespace GraphqlApi {
  export type AsObject = {
    metadata?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ObjectMeta.AsObject,
    spec?: github_com_solo_io_solo_apis_api_gloo_graphql_gloo_v1alpha1_graphql_pb.GraphQLApiSpec.AsObject,
    status?: github_com_solo_io_solo_apis_api_gloo_graphql_gloo_v1alpha1_graphql_pb.GraphQLApiStatus.AsObject,
    glooInstance?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef.AsObject,
  }
}

export class GraphqlApiSummary extends jspb.Message {
  hasMetadata(): boolean;
  clearMetadata(): void;
  getMetadata(): github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ObjectMeta | undefined;
  setMetadata(value?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ObjectMeta): void;

  hasStatus(): boolean;
  clearStatus(): void;
  getStatus(): github_com_solo_io_solo_apis_api_gloo_graphql_gloo_v1alpha1_graphql_pb.GraphQLApiStatus | undefined;
  setStatus(value?: github_com_solo_io_solo_apis_api_gloo_graphql_gloo_v1alpha1_graphql_pb.GraphQLApiStatus): void;

  hasGlooInstance(): boolean;
  clearGlooInstance(): void;
  getGlooInstance(): github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef | undefined;
  setGlooInstance(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef): void;

  hasExecutable(): boolean;
  clearExecutable(): void;
  getExecutable(): GraphqlApiSummary.ExecutableSchemaSummary | undefined;
  setExecutable(value?: GraphqlApiSummary.ExecutableSchemaSummary): void;

  hasStitched(): boolean;
  clearStitched(): void;
  getStitched(): GraphqlApiSummary.StitchedSchemaSummary | undefined;
  setStitched(value?: GraphqlApiSummary.StitchedSchemaSummary): void;

  getApiTypeSummaryCase(): GraphqlApiSummary.ApiTypeSummaryCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GraphqlApiSummary.AsObject;
  static toObject(includeInstance: boolean, msg: GraphqlApiSummary): GraphqlApiSummary.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GraphqlApiSummary, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GraphqlApiSummary;
  static deserializeBinaryFromReader(message: GraphqlApiSummary, reader: jspb.BinaryReader): GraphqlApiSummary;
}

export namespace GraphqlApiSummary {
  export type AsObject = {
    metadata?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ObjectMeta.AsObject,
    status?: github_com_solo_io_solo_apis_api_gloo_graphql_gloo_v1alpha1_graphql_pb.GraphQLApiStatus.AsObject,
    glooInstance?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef.AsObject,
    executable?: GraphqlApiSummary.ExecutableSchemaSummary.AsObject,
    stitched?: GraphqlApiSummary.StitchedSchemaSummary.AsObject,
  }

  export class ExecutableSchemaSummary extends jspb.Message {
    getNumResolvers(): number;
    setNumResolvers(value: number): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ExecutableSchemaSummary.AsObject;
    static toObject(includeInstance: boolean, msg: ExecutableSchemaSummary): ExecutableSchemaSummary.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: ExecutableSchemaSummary, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ExecutableSchemaSummary;
    static deserializeBinaryFromReader(message: ExecutableSchemaSummary, reader: jspb.BinaryReader): ExecutableSchemaSummary;
  }

  export namespace ExecutableSchemaSummary {
    export type AsObject = {
      numResolvers: number,
    }
  }

  export class StitchedSchemaSummary extends jspb.Message {
    getNumSubschemas(): number;
    setNumSubschemas(value: number): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): StitchedSchemaSummary.AsObject;
    static toObject(includeInstance: boolean, msg: StitchedSchemaSummary): StitchedSchemaSummary.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: StitchedSchemaSummary, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): StitchedSchemaSummary;
    static deserializeBinaryFromReader(message: StitchedSchemaSummary, reader: jspb.BinaryReader): StitchedSchemaSummary;
  }

  export namespace StitchedSchemaSummary {
    export type AsObject = {
      numSubschemas: number,
    }
  }

  export enum ApiTypeSummaryCase {
    API_TYPE_SUMMARY_NOT_SET = 0,
    EXECUTABLE = 4,
    STITCHED = 5,
  }
}

export class GetGraphqlApiRequest extends jspb.Message {
  hasGraphqlApiRef(): boolean;
  clearGraphqlApiRef(): void;
  getGraphqlApiRef(): github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef | undefined;
  setGraphqlApiRef(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetGraphqlApiRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetGraphqlApiRequest): GetGraphqlApiRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetGraphqlApiRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetGraphqlApiRequest;
  static deserializeBinaryFromReader(message: GetGraphqlApiRequest, reader: jspb.BinaryReader): GetGraphqlApiRequest;
}

export namespace GetGraphqlApiRequest {
  export type AsObject = {
    graphqlApiRef?: github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef.AsObject,
  }
}

export class GetGraphqlApiResponse extends jspb.Message {
  hasGraphqlApi(): boolean;
  clearGraphqlApi(): void;
  getGraphqlApi(): GraphqlApi | undefined;
  setGraphqlApi(value?: GraphqlApi): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetGraphqlApiResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetGraphqlApiResponse): GetGraphqlApiResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetGraphqlApiResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetGraphqlApiResponse;
  static deserializeBinaryFromReader(message: GetGraphqlApiResponse, reader: jspb.BinaryReader): GetGraphqlApiResponse;
}

export namespace GetGraphqlApiResponse {
  export type AsObject = {
    graphqlApi?: GraphqlApi.AsObject,
  }
}

export class ListGraphqlApisRequest extends jspb.Message {
  hasGlooInstanceRef(): boolean;
  clearGlooInstanceRef(): void;
  getGlooInstanceRef(): github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef | undefined;
  setGlooInstanceRef(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListGraphqlApisRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListGraphqlApisRequest): ListGraphqlApisRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListGraphqlApisRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListGraphqlApisRequest;
  static deserializeBinaryFromReader(message: ListGraphqlApisRequest, reader: jspb.BinaryReader): ListGraphqlApisRequest;
}

export namespace ListGraphqlApisRequest {
  export type AsObject = {
    glooInstanceRef?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef.AsObject,
  }
}

export class ListGraphqlApisResponse extends jspb.Message {
  clearGraphqlApisList(): void;
  getGraphqlApisList(): Array<GraphqlApiSummary>;
  setGraphqlApisList(value: Array<GraphqlApiSummary>): void;
  addGraphqlApis(value?: GraphqlApiSummary, index?: number): GraphqlApiSummary;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListGraphqlApisResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListGraphqlApisResponse): ListGraphqlApisResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListGraphqlApisResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListGraphqlApisResponse;
  static deserializeBinaryFromReader(message: ListGraphqlApisResponse, reader: jspb.BinaryReader): ListGraphqlApisResponse;
}

export namespace ListGraphqlApisResponse {
  export type AsObject = {
    graphqlApisList: Array<GraphqlApiSummary.AsObject>,
  }
}

export class GetGraphqlApiYamlRequest extends jspb.Message {
  hasGraphqlApiRef(): boolean;
  clearGraphqlApiRef(): void;
  getGraphqlApiRef(): github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef | undefined;
  setGraphqlApiRef(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetGraphqlApiYamlRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetGraphqlApiYamlRequest): GetGraphqlApiYamlRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetGraphqlApiYamlRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetGraphqlApiYamlRequest;
  static deserializeBinaryFromReader(message: GetGraphqlApiYamlRequest, reader: jspb.BinaryReader): GetGraphqlApiYamlRequest;
}

export namespace GetGraphqlApiYamlRequest {
  export type AsObject = {
    graphqlApiRef?: github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef.AsObject,
  }
}

export class GetGraphqlApiYamlResponse extends jspb.Message {
  hasYamlData(): boolean;
  clearYamlData(): void;
  getYamlData(): github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ResourceYaml | undefined;
  setYamlData(value?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ResourceYaml): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetGraphqlApiYamlResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetGraphqlApiYamlResponse): GetGraphqlApiYamlResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetGraphqlApiYamlResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetGraphqlApiYamlResponse;
  static deserializeBinaryFromReader(message: GetGraphqlApiYamlResponse, reader: jspb.BinaryReader): GetGraphqlApiYamlResponse;
}

export namespace GetGraphqlApiYamlResponse {
  export type AsObject = {
    yamlData?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ResourceYaml.AsObject,
  }
}

export class CreateGraphqlApiRequest extends jspb.Message {
  hasGraphqlApiRef(): boolean;
  clearGraphqlApiRef(): void;
  getGraphqlApiRef(): github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef | undefined;
  setGraphqlApiRef(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef): void;

  hasSpec(): boolean;
  clearSpec(): void;
  getSpec(): github_com_solo_io_solo_apis_api_gloo_graphql_gloo_v1alpha1_graphql_pb.GraphQLApiSpec | undefined;
  setSpec(value?: github_com_solo_io_solo_apis_api_gloo_graphql_gloo_v1alpha1_graphql_pb.GraphQLApiSpec): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateGraphqlApiRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateGraphqlApiRequest): CreateGraphqlApiRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateGraphqlApiRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateGraphqlApiRequest;
  static deserializeBinaryFromReader(message: CreateGraphqlApiRequest, reader: jspb.BinaryReader): CreateGraphqlApiRequest;
}

export namespace CreateGraphqlApiRequest {
  export type AsObject = {
    graphqlApiRef?: github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef.AsObject,
    spec?: github_com_solo_io_solo_apis_api_gloo_graphql_gloo_v1alpha1_graphql_pb.GraphQLApiSpec.AsObject,
  }
}

export class CreateGraphqlApiResponse extends jspb.Message {
  hasGraphqlApi(): boolean;
  clearGraphqlApi(): void;
  getGraphqlApi(): GraphqlApi | undefined;
  setGraphqlApi(value?: GraphqlApi): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateGraphqlApiResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CreateGraphqlApiResponse): CreateGraphqlApiResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateGraphqlApiResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateGraphqlApiResponse;
  static deserializeBinaryFromReader(message: CreateGraphqlApiResponse, reader: jspb.BinaryReader): CreateGraphqlApiResponse;
}

export namespace CreateGraphqlApiResponse {
  export type AsObject = {
    graphqlApi?: GraphqlApi.AsObject,
  }
}

export class UpdateGraphqlApiRequest extends jspb.Message {
  hasGraphqlApiRef(): boolean;
  clearGraphqlApiRef(): void;
  getGraphqlApiRef(): github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef | undefined;
  setGraphqlApiRef(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef): void;

  hasSpec(): boolean;
  clearSpec(): void;
  getSpec(): github_com_solo_io_solo_apis_api_gloo_graphql_gloo_v1alpha1_graphql_pb.GraphQLApiSpec | undefined;
  setSpec(value?: github_com_solo_io_solo_apis_api_gloo_graphql_gloo_v1alpha1_graphql_pb.GraphQLApiSpec): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateGraphqlApiRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateGraphqlApiRequest): UpdateGraphqlApiRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateGraphqlApiRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateGraphqlApiRequest;
  static deserializeBinaryFromReader(message: UpdateGraphqlApiRequest, reader: jspb.BinaryReader): UpdateGraphqlApiRequest;
}

export namespace UpdateGraphqlApiRequest {
  export type AsObject = {
    graphqlApiRef?: github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef.AsObject,
    spec?: github_com_solo_io_solo_apis_api_gloo_graphql_gloo_v1alpha1_graphql_pb.GraphQLApiSpec.AsObject,
  }
}

export class UpdateGraphqlApiResponse extends jspb.Message {
  hasGraphqlApi(): boolean;
  clearGraphqlApi(): void;
  getGraphqlApi(): GraphqlApi | undefined;
  setGraphqlApi(value?: GraphqlApi): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateGraphqlApiResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateGraphqlApiResponse): UpdateGraphqlApiResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateGraphqlApiResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateGraphqlApiResponse;
  static deserializeBinaryFromReader(message: UpdateGraphqlApiResponse, reader: jspb.BinaryReader): UpdateGraphqlApiResponse;
}

export namespace UpdateGraphqlApiResponse {
  export type AsObject = {
    graphqlApi?: GraphqlApi.AsObject,
  }
}

export class DeleteGraphqlApiRequest extends jspb.Message {
  hasGraphqlApiRef(): boolean;
  clearGraphqlApiRef(): void;
  getGraphqlApiRef(): github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef | undefined;
  setGraphqlApiRef(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteGraphqlApiRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteGraphqlApiRequest): DeleteGraphqlApiRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteGraphqlApiRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteGraphqlApiRequest;
  static deserializeBinaryFromReader(message: DeleteGraphqlApiRequest, reader: jspb.BinaryReader): DeleteGraphqlApiRequest;
}

export namespace DeleteGraphqlApiRequest {
  export type AsObject = {
    graphqlApiRef?: github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef.AsObject,
  }
}

export class DeleteGraphqlApiResponse extends jspb.Message {
  hasGraphqlApiRef(): boolean;
  clearGraphqlApiRef(): void;
  getGraphqlApiRef(): github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef | undefined;
  setGraphqlApiRef(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteGraphqlApiResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteGraphqlApiResponse): DeleteGraphqlApiResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteGraphqlApiResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteGraphqlApiResponse;
  static deserializeBinaryFromReader(message: DeleteGraphqlApiResponse, reader: jspb.BinaryReader): DeleteGraphqlApiResponse;
}

export namespace DeleteGraphqlApiResponse {
  export type AsObject = {
    graphqlApiRef?: github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef.AsObject,
  }
}

export class ValidateResolverYamlRequest extends jspb.Message {
  getYaml(): string;
  setYaml(value: string): void;

  getResolverType(): ValidateResolverYamlRequest.ResolverTypeMap[keyof ValidateResolverYamlRequest.ResolverTypeMap];
  setResolverType(value: ValidateResolverYamlRequest.ResolverTypeMap[keyof ValidateResolverYamlRequest.ResolverTypeMap]): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ValidateResolverYamlRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ValidateResolverYamlRequest): ValidateResolverYamlRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ValidateResolverYamlRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ValidateResolverYamlRequest;
  static deserializeBinaryFromReader(message: ValidateResolverYamlRequest, reader: jspb.BinaryReader): ValidateResolverYamlRequest;
}

export namespace ValidateResolverYamlRequest {
  export type AsObject = {
    yaml: string,
    resolverType: ValidateResolverYamlRequest.ResolverTypeMap[keyof ValidateResolverYamlRequest.ResolverTypeMap],
  }

  export interface ResolverTypeMap {
    RESOLVER_NOT_SET: 0;
    REST_RESOLVER: 1;
    GRPC_RESOLVER: 2;
  }

  export const ResolverType: ResolverTypeMap;
}

export class ValidateResolverYamlResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ValidateResolverYamlResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ValidateResolverYamlResponse): ValidateResolverYamlResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ValidateResolverYamlResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ValidateResolverYamlResponse;
  static deserializeBinaryFromReader(message: ValidateResolverYamlResponse, reader: jspb.BinaryReader): ValidateResolverYamlResponse;
}

export namespace ValidateResolverYamlResponse {
  export type AsObject = {
  }
}

export class ValidateSchemaDefinitionRequest extends jspb.Message {
  hasSchemaDefinition(): boolean;
  clearSchemaDefinition(): void;
  getSchemaDefinition(): string;
  setSchemaDefinition(value: string): void;

  hasSpec(): boolean;
  clearSpec(): void;
  getSpec(): github_com_solo_io_solo_apis_api_gloo_graphql_gloo_v1alpha1_graphql_pb.GraphQLApiSpec | undefined;
  setSpec(value?: github_com_solo_io_solo_apis_api_gloo_graphql_gloo_v1alpha1_graphql_pb.GraphQLApiSpec): void;

  getInputCase(): ValidateSchemaDefinitionRequest.InputCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ValidateSchemaDefinitionRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ValidateSchemaDefinitionRequest): ValidateSchemaDefinitionRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ValidateSchemaDefinitionRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ValidateSchemaDefinitionRequest;
  static deserializeBinaryFromReader(message: ValidateSchemaDefinitionRequest, reader: jspb.BinaryReader): ValidateSchemaDefinitionRequest;
}

export namespace ValidateSchemaDefinitionRequest {
  export type AsObject = {
    schemaDefinition: string,
    spec?: github_com_solo_io_solo_apis_api_gloo_graphql_gloo_v1alpha1_graphql_pb.GraphQLApiSpec.AsObject,
  }

  export enum InputCase {
    INPUT_NOT_SET = 0,
    SCHEMA_DEFINITION = 1,
    SPEC = 2,
  }
}

export class ValidateSchemaDefinitionResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ValidateSchemaDefinitionResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ValidateSchemaDefinitionResponse): ValidateSchemaDefinitionResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ValidateSchemaDefinitionResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ValidateSchemaDefinitionResponse;
  static deserializeBinaryFromReader(message: ValidateSchemaDefinitionResponse, reader: jspb.BinaryReader): ValidateSchemaDefinitionResponse;
}

export namespace ValidateSchemaDefinitionResponse {
  export type AsObject = {
  }
}

export class GetStitchedSchemaDefinitionRequest extends jspb.Message {
  hasStitchedSchemaApiRef(): boolean;
  clearStitchedSchemaApiRef(): void;
  getStitchedSchemaApiRef(): github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef | undefined;
  setStitchedSchemaApiRef(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetStitchedSchemaDefinitionRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetStitchedSchemaDefinitionRequest): GetStitchedSchemaDefinitionRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetStitchedSchemaDefinitionRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetStitchedSchemaDefinitionRequest;
  static deserializeBinaryFromReader(message: GetStitchedSchemaDefinitionRequest, reader: jspb.BinaryReader): GetStitchedSchemaDefinitionRequest;
}

export namespace GetStitchedSchemaDefinitionRequest {
  export type AsObject = {
    stitchedSchemaApiRef?: github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef.AsObject,
  }
}

export class GetStitchedSchemaDefinitionResponse extends jspb.Message {
  getStitchedSchemaSdl(): string;
  setStitchedSchemaSdl(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetStitchedSchemaDefinitionResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetStitchedSchemaDefinitionResponse): GetStitchedSchemaDefinitionResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetStitchedSchemaDefinitionResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetStitchedSchemaDefinitionResponse;
  static deserializeBinaryFromReader(message: GetStitchedSchemaDefinitionResponse, reader: jspb.BinaryReader): GetStitchedSchemaDefinitionResponse;
}

export namespace GetStitchedSchemaDefinitionResponse {
  export type AsObject = {
    stitchedSchemaSdl: string,
  }
}
