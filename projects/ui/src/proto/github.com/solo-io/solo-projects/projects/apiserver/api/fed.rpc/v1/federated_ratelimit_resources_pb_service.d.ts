// package: fed.rpc.solo.io
// file: github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/federated_ratelimit_resources.proto

import * as github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_ratelimit_resources_pb from "../../../../../../../../github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/federated_ratelimit_resources_pb";
import {grpc} from "@improbable-eng/grpc-web";

type FederatedRatelimitResourceApiListFederatedRateLimitConfigs = {
  readonly methodName: string;
  readonly service: typeof FederatedRatelimitResourceApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_ratelimit_resources_pb.ListFederatedRateLimitConfigsRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_ratelimit_resources_pb.ListFederatedRateLimitConfigsResponse;
};

type FederatedRatelimitResourceApiGetFederatedRateLimitConfigYaml = {
  readonly methodName: string;
  readonly service: typeof FederatedRatelimitResourceApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_ratelimit_resources_pb.GetFederatedRateLimitConfigYamlRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_ratelimit_resources_pb.GetFederatedRateLimitConfigYamlResponse;
};

export class FederatedRatelimitResourceApi {
  static readonly serviceName: string;
  static readonly ListFederatedRateLimitConfigs: FederatedRatelimitResourceApiListFederatedRateLimitConfigs;
  static readonly GetFederatedRateLimitConfigYaml: FederatedRatelimitResourceApiGetFederatedRateLimitConfigYaml;
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

export class FederatedRatelimitResourceApiClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  listFederatedRateLimitConfigs(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_ratelimit_resources_pb.ListFederatedRateLimitConfigsRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_ratelimit_resources_pb.ListFederatedRateLimitConfigsResponse|null) => void
  ): UnaryResponse;
  listFederatedRateLimitConfigs(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_ratelimit_resources_pb.ListFederatedRateLimitConfigsRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_ratelimit_resources_pb.ListFederatedRateLimitConfigsResponse|null) => void
  ): UnaryResponse;
  getFederatedRateLimitConfigYaml(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_ratelimit_resources_pb.GetFederatedRateLimitConfigYamlRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_ratelimit_resources_pb.GetFederatedRateLimitConfigYamlResponse|null) => void
  ): UnaryResponse;
  getFederatedRateLimitConfigYaml(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_ratelimit_resources_pb.GetFederatedRateLimitConfigYamlRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_ratelimit_resources_pb.GetFederatedRateLimitConfigYamlResponse|null) => void
  ): UnaryResponse;
}

