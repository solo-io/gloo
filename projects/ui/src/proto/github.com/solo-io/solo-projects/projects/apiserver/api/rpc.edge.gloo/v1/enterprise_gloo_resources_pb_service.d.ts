// package: rpc.edge.gloo.solo.io
// file: github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/enterprise_gloo_resources.proto

import * as github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_enterprise_gloo_resources_pb from "../../../../../../../../github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/enterprise_gloo_resources_pb";
import {grpc} from "@improbable-eng/grpc-web";

type EnterpriseGlooResourceApiListAuthConfigs = {
  readonly methodName: string;
  readonly service: typeof EnterpriseGlooResourceApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_enterprise_gloo_resources_pb.ListAuthConfigsRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_enterprise_gloo_resources_pb.ListAuthConfigsResponse;
};

type EnterpriseGlooResourceApiGetAuthConfigYaml = {
  readonly methodName: string;
  readonly service: typeof EnterpriseGlooResourceApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_enterprise_gloo_resources_pb.GetAuthConfigYamlRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_enterprise_gloo_resources_pb.GetAuthConfigYamlResponse;
};

export class EnterpriseGlooResourceApi {
  static readonly serviceName: string;
  static readonly ListAuthConfigs: EnterpriseGlooResourceApiListAuthConfigs;
  static readonly GetAuthConfigYaml: EnterpriseGlooResourceApiGetAuthConfigYaml;
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

export class EnterpriseGlooResourceApiClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  listAuthConfigs(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_enterprise_gloo_resources_pb.ListAuthConfigsRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_enterprise_gloo_resources_pb.ListAuthConfigsResponse|null) => void
  ): UnaryResponse;
  listAuthConfigs(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_enterprise_gloo_resources_pb.ListAuthConfigsRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_enterprise_gloo_resources_pb.ListAuthConfigsResponse|null) => void
  ): UnaryResponse;
  getAuthConfigYaml(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_enterprise_gloo_resources_pb.GetAuthConfigYamlRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_enterprise_gloo_resources_pb.GetAuthConfigYamlResponse|null) => void
  ): UnaryResponse;
  getAuthConfigYaml(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_enterprise_gloo_resources_pb.GetAuthConfigYamlRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_enterprise_gloo_resources_pb.GetAuthConfigYamlResponse|null) => void
  ): UnaryResponse;
}

