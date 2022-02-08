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
