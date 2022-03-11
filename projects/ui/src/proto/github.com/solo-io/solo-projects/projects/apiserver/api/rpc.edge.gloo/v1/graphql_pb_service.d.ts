// package: rpc.edge.gloo.solo.io
// file: github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/graphql.proto

import * as github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb from "../../../../../../../../github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/graphql_pb";
import {grpc} from "@improbable-eng/grpc-web";

type GraphqlApiGetGraphqlSchema = {
  readonly methodName: string;
  readonly service: typeof GraphqlApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.GetGraphqlSchemaRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.GetGraphqlSchemaResponse;
};

type GraphqlApiListGraphqlSchemas = {
  readonly methodName: string;
  readonly service: typeof GraphqlApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.ListGraphqlSchemasRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.ListGraphqlSchemasResponse;
};

type GraphqlApiGetGraphqlSchemaYaml = {
  readonly methodName: string;
  readonly service: typeof GraphqlApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.GetGraphqlSchemaYamlRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.GetGraphqlSchemaYamlResponse;
};

type GraphqlApiCreateGraphqlSchema = {
  readonly methodName: string;
  readonly service: typeof GraphqlApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.CreateGraphqlSchemaRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.CreateGraphqlSchemaResponse;
};

type GraphqlApiUpdateGraphqlSchema = {
  readonly methodName: string;
  readonly service: typeof GraphqlApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.UpdateGraphqlSchemaRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.UpdateGraphqlSchemaResponse;
};

type GraphqlApiDeleteGraphqlSchema = {
  readonly methodName: string;
  readonly service: typeof GraphqlApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.DeleteGraphqlSchemaRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.DeleteGraphqlSchemaResponse;
};

type GraphqlApiValidateResolverYaml = {
  readonly methodName: string;
  readonly service: typeof GraphqlApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.ValidateResolverYamlRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.ValidateResolverYamlResponse;
};

type GraphqlApiValidateSchemaDefinition = {
  readonly methodName: string;
  readonly service: typeof GraphqlApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.ValidateSchemaDefinitionRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.ValidateSchemaDefinitionResponse;
};

export class GraphqlApi {
  static readonly serviceName: string;
  static readonly GetGraphqlSchema: GraphqlApiGetGraphqlSchema;
  static readonly ListGraphqlSchemas: GraphqlApiListGraphqlSchemas;
  static readonly GetGraphqlSchemaYaml: GraphqlApiGetGraphqlSchemaYaml;
  static readonly CreateGraphqlSchema: GraphqlApiCreateGraphqlSchema;
  static readonly UpdateGraphqlSchema: GraphqlApiUpdateGraphqlSchema;
  static readonly DeleteGraphqlSchema: GraphqlApiDeleteGraphqlSchema;
  static readonly ValidateResolverYaml: GraphqlApiValidateResolverYaml;
  static readonly ValidateSchemaDefinition: GraphqlApiValidateSchemaDefinition;
}

export type ServiceError = { message: string, code: number; metadata: grpc.Metadata }
export type Status = { details: string, code: number; metadata: grpc.Metadata }

interface UnaryResponse {
  cancel(): void;
}
interface ResponseStream<T> {
  cancel(): void;
  on(type: 'data', handler: (message: T) => void): ResponseStream<T>;
  on(type: 'end', handler: (status?: Status) => void): ResponseStream<T>;
  on(type: 'status', handler: (status: Status) => void): ResponseStream<T>;
}
interface RequestStream<T> {
  write(message: T): RequestStream<T>;
  end(): void;
  cancel(): void;
  on(type: 'end', handler: (status?: Status) => void): RequestStream<T>;
  on(type: 'status', handler: (status: Status) => void): RequestStream<T>;
}
interface BidirectionalStream<ReqT, ResT> {
  write(message: ReqT): BidirectionalStream<ReqT, ResT>;
  end(): void;
  cancel(): void;
  on(type: 'data', handler: (message: ResT) => void): BidirectionalStream<ReqT, ResT>;
  on(type: 'end', handler: (status?: Status) => void): BidirectionalStream<ReqT, ResT>;
  on(type: 'status', handler: (status: Status) => void): BidirectionalStream<ReqT, ResT>;
}

export class GraphqlApiClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  getGraphqlSchema(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.GetGraphqlSchemaRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.GetGraphqlSchemaResponse|null) => void
  ): UnaryResponse;
  getGraphqlSchema(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.GetGraphqlSchemaRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.GetGraphqlSchemaResponse|null) => void
  ): UnaryResponse;
  listGraphqlSchemas(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.ListGraphqlSchemasRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.ListGraphqlSchemasResponse|null) => void
  ): UnaryResponse;
  listGraphqlSchemas(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.ListGraphqlSchemasRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.ListGraphqlSchemasResponse|null) => void
  ): UnaryResponse;
  getGraphqlSchemaYaml(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.GetGraphqlSchemaYamlRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.GetGraphqlSchemaYamlResponse|null) => void
  ): UnaryResponse;
  getGraphqlSchemaYaml(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.GetGraphqlSchemaYamlRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.GetGraphqlSchemaYamlResponse|null) => void
  ): UnaryResponse;
  createGraphqlSchema(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.CreateGraphqlSchemaRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.CreateGraphqlSchemaResponse|null) => void
  ): UnaryResponse;
  createGraphqlSchema(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.CreateGraphqlSchemaRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.CreateGraphqlSchemaResponse|null) => void
  ): UnaryResponse;
  updateGraphqlSchema(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.UpdateGraphqlSchemaRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.UpdateGraphqlSchemaResponse|null) => void
  ): UnaryResponse;
  updateGraphqlSchema(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.UpdateGraphqlSchemaRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.UpdateGraphqlSchemaResponse|null) => void
  ): UnaryResponse;
  deleteGraphqlSchema(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.DeleteGraphqlSchemaRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.DeleteGraphqlSchemaResponse|null) => void
  ): UnaryResponse;
  deleteGraphqlSchema(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.DeleteGraphqlSchemaRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.DeleteGraphqlSchemaResponse|null) => void
  ): UnaryResponse;
  validateResolverYaml(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.ValidateResolverYamlRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.ValidateResolverYamlResponse|null) => void
  ): UnaryResponse;
  validateResolverYaml(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.ValidateResolverYamlRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.ValidateResolverYamlResponse|null) => void
  ): UnaryResponse;
  validateSchemaDefinition(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.ValidateSchemaDefinitionRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.ValidateSchemaDefinitionResponse|null) => void
  ): UnaryResponse;
  validateSchemaDefinition(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.ValidateSchemaDefinitionRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.ValidateSchemaDefinitionResponse|null) => void
  ): UnaryResponse;
}

