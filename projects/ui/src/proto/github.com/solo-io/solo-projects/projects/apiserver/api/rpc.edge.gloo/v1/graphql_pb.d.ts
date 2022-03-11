/* eslint-disable */
// package: rpc.edge.gloo.solo.io
// file: github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/graphql.proto

import * as jspb from "google-protobuf";
import * as extproto_ext_pb from "../../../../../../../../extproto/ext_pb";
import * as github_com_solo_io_skv2_api_core_v1_core_pb from "../../../../../../../../github.com/solo-io/skv2/api/core/v1/core_pb";
import * as github_com_solo_io_solo_apis_api_gloo_graphql_gloo_v1alpha1_graphql_pb from "../../../../../../../../github.com/solo-io/solo-apis/api/gloo/graphql.gloo/v1alpha1/graphql_pb";
import * as github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb from "../../../../../../../../github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/common_pb";

export class GraphqlSchema extends jspb.Message {
  hasMetadata(): boolean;
  clearMetadata(): void;
  getMetadata(): github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ObjectMeta | undefined;
  setMetadata(value?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ObjectMeta): void;

  hasSpec(): boolean;
  clearSpec(): void;
  getSpec(): github_com_solo_io_solo_apis_api_gloo_graphql_gloo_v1alpha1_graphql_pb.GraphQLSchemaSpec | undefined;
  setSpec(value?: github_com_solo_io_solo_apis_api_gloo_graphql_gloo_v1alpha1_graphql_pb.GraphQLSchemaSpec): void;

  hasStatus(): boolean;
  clearStatus(): void;
  getStatus(): github_com_solo_io_solo_apis_api_gloo_graphql_gloo_v1alpha1_graphql_pb.GraphQLSchemaStatus | undefined;
  setStatus(value?: github_com_solo_io_solo_apis_api_gloo_graphql_gloo_v1alpha1_graphql_pb.GraphQLSchemaStatus): void;

  hasGlooInstance(): boolean;
  clearGlooInstance(): void;
  getGlooInstance(): github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef | undefined;
  setGlooInstance(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GraphqlSchema.AsObject;
  static toObject(includeInstance: boolean, msg: GraphqlSchema): GraphqlSchema.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GraphqlSchema, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GraphqlSchema;
  static deserializeBinaryFromReader(message: GraphqlSchema, reader: jspb.BinaryReader): GraphqlSchema;
}

export namespace GraphqlSchema {
  export type AsObject = {
    metadata?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ObjectMeta.AsObject,
    spec?: github_com_solo_io_solo_apis_api_gloo_graphql_gloo_v1alpha1_graphql_pb.GraphQLSchemaSpec.AsObject,
    status?: github_com_solo_io_solo_apis_api_gloo_graphql_gloo_v1alpha1_graphql_pb.GraphQLSchemaStatus.AsObject,
    glooInstance?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef.AsObject,
  }
}

export class GetGraphqlSchemaRequest extends jspb.Message {
  hasGraphqlSchemaRef(): boolean;
  clearGraphqlSchemaRef(): void;
  getGraphqlSchemaRef(): github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef | undefined;
  setGraphqlSchemaRef(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetGraphqlSchemaRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetGraphqlSchemaRequest): GetGraphqlSchemaRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetGraphqlSchemaRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetGraphqlSchemaRequest;
  static deserializeBinaryFromReader(message: GetGraphqlSchemaRequest, reader: jspb.BinaryReader): GetGraphqlSchemaRequest;
}

export namespace GetGraphqlSchemaRequest {
  export type AsObject = {
    graphqlSchemaRef?: github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef.AsObject,
  }
}

export class GetGraphqlSchemaResponse extends jspb.Message {
  hasGraphqlSchema(): boolean;
  clearGraphqlSchema(): void;
  getGraphqlSchema(): GraphqlSchema | undefined;
  setGraphqlSchema(value?: GraphqlSchema): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetGraphqlSchemaResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetGraphqlSchemaResponse): GetGraphqlSchemaResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetGraphqlSchemaResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetGraphqlSchemaResponse;
  static deserializeBinaryFromReader(message: GetGraphqlSchemaResponse, reader: jspb.BinaryReader): GetGraphqlSchemaResponse;
}

export namespace GetGraphqlSchemaResponse {
  export type AsObject = {
    graphqlSchema?: GraphqlSchema.AsObject,
  }
}

export class ListGraphqlSchemasRequest extends jspb.Message {
  hasGlooInstanceRef(): boolean;
  clearGlooInstanceRef(): void;
  getGlooInstanceRef(): github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef | undefined;
  setGlooInstanceRef(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListGraphqlSchemasRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListGraphqlSchemasRequest): ListGraphqlSchemasRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListGraphqlSchemasRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListGraphqlSchemasRequest;
  static deserializeBinaryFromReader(message: ListGraphqlSchemasRequest, reader: jspb.BinaryReader): ListGraphqlSchemasRequest;
}

export namespace ListGraphqlSchemasRequest {
  export type AsObject = {
    glooInstanceRef?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef.AsObject,
  }
}

export class ListGraphqlSchemasResponse extends jspb.Message {
  clearGraphqlSchemasList(): void;
  getGraphqlSchemasList(): Array<GraphqlSchema>;
  setGraphqlSchemasList(value: Array<GraphqlSchema>): void;
  addGraphqlSchemas(value?: GraphqlSchema, index?: number): GraphqlSchema;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListGraphqlSchemasResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListGraphqlSchemasResponse): ListGraphqlSchemasResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListGraphqlSchemasResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListGraphqlSchemasResponse;
  static deserializeBinaryFromReader(message: ListGraphqlSchemasResponse, reader: jspb.BinaryReader): ListGraphqlSchemasResponse;
}

export namespace ListGraphqlSchemasResponse {
  export type AsObject = {
    graphqlSchemasList: Array<GraphqlSchema.AsObject>,
  }
}

export class GetGraphqlSchemaYamlRequest extends jspb.Message {
  hasGraphqlSchemaRef(): boolean;
  clearGraphqlSchemaRef(): void;
  getGraphqlSchemaRef(): github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef | undefined;
  setGraphqlSchemaRef(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetGraphqlSchemaYamlRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetGraphqlSchemaYamlRequest): GetGraphqlSchemaYamlRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetGraphqlSchemaYamlRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetGraphqlSchemaYamlRequest;
  static deserializeBinaryFromReader(message: GetGraphqlSchemaYamlRequest, reader: jspb.BinaryReader): GetGraphqlSchemaYamlRequest;
}

export namespace GetGraphqlSchemaYamlRequest {
  export type AsObject = {
    graphqlSchemaRef?: github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef.AsObject,
  }
}

export class GetGraphqlSchemaYamlResponse extends jspb.Message {
  hasYamlData(): boolean;
  clearYamlData(): void;
  getYamlData(): github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ResourceYaml | undefined;
  setYamlData(value?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ResourceYaml): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetGraphqlSchemaYamlResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetGraphqlSchemaYamlResponse): GetGraphqlSchemaYamlResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetGraphqlSchemaYamlResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetGraphqlSchemaYamlResponse;
  static deserializeBinaryFromReader(message: GetGraphqlSchemaYamlResponse, reader: jspb.BinaryReader): GetGraphqlSchemaYamlResponse;
}

export namespace GetGraphqlSchemaYamlResponse {
  export type AsObject = {
    yamlData?: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_common_pb.ResourceYaml.AsObject,
  }
}

export class CreateGraphqlSchemaRequest extends jspb.Message {
  hasGraphqlSchemaRef(): boolean;
  clearGraphqlSchemaRef(): void;
  getGraphqlSchemaRef(): github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef | undefined;
  setGraphqlSchemaRef(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef): void;

  hasSpec(): boolean;
  clearSpec(): void;
  getSpec(): github_com_solo_io_solo_apis_api_gloo_graphql_gloo_v1alpha1_graphql_pb.GraphQLSchemaSpec | undefined;
  setSpec(value?: github_com_solo_io_solo_apis_api_gloo_graphql_gloo_v1alpha1_graphql_pb.GraphQLSchemaSpec): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateGraphqlSchemaRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateGraphqlSchemaRequest): CreateGraphqlSchemaRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateGraphqlSchemaRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateGraphqlSchemaRequest;
  static deserializeBinaryFromReader(message: CreateGraphqlSchemaRequest, reader: jspb.BinaryReader): CreateGraphqlSchemaRequest;
}

export namespace CreateGraphqlSchemaRequest {
  export type AsObject = {
    graphqlSchemaRef?: github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef.AsObject,
    spec?: github_com_solo_io_solo_apis_api_gloo_graphql_gloo_v1alpha1_graphql_pb.GraphQLSchemaSpec.AsObject,
  }
}

export class CreateGraphqlSchemaResponse extends jspb.Message {
  hasGraphqlSchema(): boolean;
  clearGraphqlSchema(): void;
  getGraphqlSchema(): GraphqlSchema | undefined;
  setGraphqlSchema(value?: GraphqlSchema): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateGraphqlSchemaResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CreateGraphqlSchemaResponse): CreateGraphqlSchemaResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateGraphqlSchemaResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateGraphqlSchemaResponse;
  static deserializeBinaryFromReader(message: CreateGraphqlSchemaResponse, reader: jspb.BinaryReader): CreateGraphqlSchemaResponse;
}

export namespace CreateGraphqlSchemaResponse {
  export type AsObject = {
    graphqlSchema?: GraphqlSchema.AsObject,
  }
}

export class UpdateGraphqlSchemaRequest extends jspb.Message {
  hasGraphqlSchemaRef(): boolean;
  clearGraphqlSchemaRef(): void;
  getGraphqlSchemaRef(): github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef | undefined;
  setGraphqlSchemaRef(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef): void;

  hasSpec(): boolean;
  clearSpec(): void;
  getSpec(): github_com_solo_io_solo_apis_api_gloo_graphql_gloo_v1alpha1_graphql_pb.GraphQLSchemaSpec | undefined;
  setSpec(value?: github_com_solo_io_solo_apis_api_gloo_graphql_gloo_v1alpha1_graphql_pb.GraphQLSchemaSpec): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateGraphqlSchemaRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateGraphqlSchemaRequest): UpdateGraphqlSchemaRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateGraphqlSchemaRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateGraphqlSchemaRequest;
  static deserializeBinaryFromReader(message: UpdateGraphqlSchemaRequest, reader: jspb.BinaryReader): UpdateGraphqlSchemaRequest;
}

export namespace UpdateGraphqlSchemaRequest {
  export type AsObject = {
    graphqlSchemaRef?: github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef.AsObject,
    spec?: github_com_solo_io_solo_apis_api_gloo_graphql_gloo_v1alpha1_graphql_pb.GraphQLSchemaSpec.AsObject,
  }
}

export class UpdateGraphqlSchemaResponse extends jspb.Message {
  hasGraphqlSchema(): boolean;
  clearGraphqlSchema(): void;
  getGraphqlSchema(): GraphqlSchema | undefined;
  setGraphqlSchema(value?: GraphqlSchema): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateGraphqlSchemaResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateGraphqlSchemaResponse): UpdateGraphqlSchemaResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateGraphqlSchemaResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateGraphqlSchemaResponse;
  static deserializeBinaryFromReader(message: UpdateGraphqlSchemaResponse, reader: jspb.BinaryReader): UpdateGraphqlSchemaResponse;
}

export namespace UpdateGraphqlSchemaResponse {
  export type AsObject = {
    graphqlSchema?: GraphqlSchema.AsObject,
  }
}

export class DeleteGraphqlSchemaRequest extends jspb.Message {
  hasGraphqlSchemaRef(): boolean;
  clearGraphqlSchemaRef(): void;
  getGraphqlSchemaRef(): github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef | undefined;
  setGraphqlSchemaRef(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteGraphqlSchemaRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteGraphqlSchemaRequest): DeleteGraphqlSchemaRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteGraphqlSchemaRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteGraphqlSchemaRequest;
  static deserializeBinaryFromReader(message: DeleteGraphqlSchemaRequest, reader: jspb.BinaryReader): DeleteGraphqlSchemaRequest;
}

export namespace DeleteGraphqlSchemaRequest {
  export type AsObject = {
    graphqlSchemaRef?: github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef.AsObject,
  }
}

export class DeleteGraphqlSchemaResponse extends jspb.Message {
  hasGraphqlSchemaRef(): boolean;
  clearGraphqlSchemaRef(): void;
  getGraphqlSchemaRef(): github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef | undefined;
  setGraphqlSchemaRef(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteGraphqlSchemaResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteGraphqlSchemaResponse): DeleteGraphqlSchemaResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteGraphqlSchemaResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteGraphqlSchemaResponse;
  static deserializeBinaryFromReader(message: DeleteGraphqlSchemaResponse, reader: jspb.BinaryReader): DeleteGraphqlSchemaResponse;
}

export namespace DeleteGraphqlSchemaResponse {
  export type AsObject = {
    graphqlSchemaRef?: github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef.AsObject,
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
  getSpec(): github_com_solo_io_solo_apis_api_gloo_graphql_gloo_v1alpha1_graphql_pb.GraphQLSchemaSpec | undefined;
  setSpec(value?: github_com_solo_io_solo_apis_api_gloo_graphql_gloo_v1alpha1_graphql_pb.GraphQLSchemaSpec): void;

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
    spec?: github_com_solo_io_solo_apis_api_gloo_graphql_gloo_v1alpha1_graphql_pb.GraphQLSchemaSpec.AsObject,
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
