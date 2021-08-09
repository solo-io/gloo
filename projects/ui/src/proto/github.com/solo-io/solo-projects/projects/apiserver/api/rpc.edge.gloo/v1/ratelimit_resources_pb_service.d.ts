// package: rpc.edge.gloo.solo.io
// file: github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/ratelimit_resources.proto

import * as github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_ratelimit_resources_pb from "../../../../../../../../github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/ratelimit_resources_pb";
import {grpc} from "@improbable-eng/grpc-web";

type RatelimitResourceApiListRateLimitConfigs = {
  readonly methodName: string;
  readonly service: typeof RatelimitResourceApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_ratelimit_resources_pb.ListRateLimitConfigsRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_ratelimit_resources_pb.ListRateLimitConfigsResponse;
};

type RatelimitResourceApiGetRateLimitConfigYaml = {
  readonly methodName: string;
  readonly service: typeof RatelimitResourceApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_ratelimit_resources_pb.GetRateLimitConfigYamlRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_ratelimit_resources_pb.GetRateLimitConfigYamlResponse;
};

export class RatelimitResourceApi {
  static readonly serviceName: string;
  static readonly ListRateLimitConfigs: RatelimitResourceApiListRateLimitConfigs;
  static readonly GetRateLimitConfigYaml: RatelimitResourceApiGetRateLimitConfigYaml;
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

export class RatelimitResourceApiClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  listRateLimitConfigs(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_ratelimit_resources_pb.ListRateLimitConfigsRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_ratelimit_resources_pb.ListRateLimitConfigsResponse|null) => void
  ): UnaryResponse;
  listRateLimitConfigs(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_ratelimit_resources_pb.ListRateLimitConfigsRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_ratelimit_resources_pb.ListRateLimitConfigsResponse|null) => void
  ): UnaryResponse;
  getRateLimitConfigYaml(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_ratelimit_resources_pb.GetRateLimitConfigYamlRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_ratelimit_resources_pb.GetRateLimitConfigYamlResponse|null) => void
  ): UnaryResponse;
  getRateLimitConfigYaml(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_ratelimit_resources_pb.GetRateLimitConfigYamlRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_ratelimit_resources_pb.GetRateLimitConfigYamlResponse|null) => void
  ): UnaryResponse;
}

