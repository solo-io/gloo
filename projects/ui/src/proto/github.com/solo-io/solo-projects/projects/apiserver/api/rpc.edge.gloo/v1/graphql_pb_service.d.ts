// package: rpc.edge.gloo.solo.io
// file: github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/graphql.proto

import * as github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb from "../../../../../../../../github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/graphql_pb";
import {grpc} from "@improbable-eng/grpc-web";

type GraphqlConfigApiGetGraphqlApi = {
  readonly methodName: string;
  readonly service: typeof GraphqlConfigApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.GetGraphqlApiRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.GetGraphqlApiResponse;
};

type GraphqlConfigApiListGraphqlApis = {
  readonly methodName: string;
  readonly service: typeof GraphqlConfigApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.ListGraphqlApisRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.ListGraphqlApisResponse;
};

type GraphqlConfigApiGetGraphqlApiYaml = {
  readonly methodName: string;
  readonly service: typeof GraphqlConfigApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.GetGraphqlApiYamlRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.GetGraphqlApiYamlResponse;
};

type GraphqlConfigApiCreateGraphqlApi = {
  readonly methodName: string;
  readonly service: typeof GraphqlConfigApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.CreateGraphqlApiRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.CreateGraphqlApiResponse;
};

type GraphqlConfigApiUpdateGraphqlApi = {
  readonly methodName: string;
  readonly service: typeof GraphqlConfigApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.UpdateGraphqlApiRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.UpdateGraphqlApiResponse;
};

type GraphqlConfigApiDeleteGraphqlApi = {
  readonly methodName: string;
  readonly service: typeof GraphqlConfigApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.DeleteGraphqlApiRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.DeleteGraphqlApiResponse;
};

type GraphqlConfigApiValidateResolverYaml = {
  readonly methodName: string;
  readonly service: typeof GraphqlConfigApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.ValidateResolverYamlRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.ValidateResolverYamlResponse;
};

type GraphqlConfigApiValidateSchemaDefinition = {
  readonly methodName: string;
  readonly service: typeof GraphqlConfigApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.ValidateSchemaDefinitionRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.ValidateSchemaDefinitionResponse;
};

type GraphqlConfigApiGetStitchedSchemaDefinition = {
  readonly methodName: string;
  readonly service: typeof GraphqlConfigApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.GetStitchedSchemaDefinitionRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.GetStitchedSchemaDefinitionResponse;
};

export class GraphqlConfigApi {
  static readonly serviceName: string;
  static readonly GetGraphqlApi: GraphqlConfigApiGetGraphqlApi;
  static readonly ListGraphqlApis: GraphqlConfigApiListGraphqlApis;
  static readonly GetGraphqlApiYaml: GraphqlConfigApiGetGraphqlApiYaml;
  static readonly CreateGraphqlApi: GraphqlConfigApiCreateGraphqlApi;
  static readonly UpdateGraphqlApi: GraphqlConfigApiUpdateGraphqlApi;
  static readonly DeleteGraphqlApi: GraphqlConfigApiDeleteGraphqlApi;
  static readonly ValidateResolverYaml: GraphqlConfigApiValidateResolverYaml;
  static readonly ValidateSchemaDefinition: GraphqlConfigApiValidateSchemaDefinition;
  static readonly GetStitchedSchemaDefinition: GraphqlConfigApiGetStitchedSchemaDefinition;
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

export class GraphqlConfigApiClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  getGraphqlApi(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.GetGraphqlApiRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.GetGraphqlApiResponse|null) => void
  ): UnaryResponse;
  getGraphqlApi(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.GetGraphqlApiRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.GetGraphqlApiResponse|null) => void
  ): UnaryResponse;
  listGraphqlApis(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.ListGraphqlApisRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.ListGraphqlApisResponse|null) => void
  ): UnaryResponse;
  listGraphqlApis(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.ListGraphqlApisRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.ListGraphqlApisResponse|null) => void
  ): UnaryResponse;
  getGraphqlApiYaml(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.GetGraphqlApiYamlRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.GetGraphqlApiYamlResponse|null) => void
  ): UnaryResponse;
  getGraphqlApiYaml(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.GetGraphqlApiYamlRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.GetGraphqlApiYamlResponse|null) => void
  ): UnaryResponse;
  createGraphqlApi(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.CreateGraphqlApiRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.CreateGraphqlApiResponse|null) => void
  ): UnaryResponse;
  createGraphqlApi(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.CreateGraphqlApiRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.CreateGraphqlApiResponse|null) => void
  ): UnaryResponse;
  updateGraphqlApi(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.UpdateGraphqlApiRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.UpdateGraphqlApiResponse|null) => void
  ): UnaryResponse;
  updateGraphqlApi(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.UpdateGraphqlApiRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.UpdateGraphqlApiResponse|null) => void
  ): UnaryResponse;
  deleteGraphqlApi(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.DeleteGraphqlApiRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.DeleteGraphqlApiResponse|null) => void
  ): UnaryResponse;
  deleteGraphqlApi(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.DeleteGraphqlApiRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.DeleteGraphqlApiResponse|null) => void
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
  getStitchedSchemaDefinition(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.GetStitchedSchemaDefinitionRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.GetStitchedSchemaDefinitionResponse|null) => void
  ): UnaryResponse;
  getStitchedSchemaDefinition(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.GetStitchedSchemaDefinitionRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.GetStitchedSchemaDefinitionResponse|null) => void
  ): UnaryResponse;
}

