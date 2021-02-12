// package: fed.rpc.solo.io
// file: github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/federated_enterprise_gloo_resources.proto

import * as github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_enterprise_gloo_resources_pb from "../../../../../../../../github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/federated_enterprise_gloo_resources_pb";
import {grpc} from "@improbable-eng/grpc-web";

type FederatedEnterpriseGlooResourceApiListFederatedAuthConfigs = {
  readonly methodName: string;
  readonly service: typeof FederatedEnterpriseGlooResourceApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_enterprise_gloo_resources_pb.ListFederatedAuthConfigsRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_enterprise_gloo_resources_pb.ListFederatedAuthConfigsResponse;
};

type FederatedEnterpriseGlooResourceApiGetFederatedAuthConfigYaml = {
  readonly methodName: string;
  readonly service: typeof FederatedEnterpriseGlooResourceApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_enterprise_gloo_resources_pb.GetFederatedAuthConfigYamlRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_enterprise_gloo_resources_pb.GetFederatedAuthConfigYamlResponse;
};

export class FederatedEnterpriseGlooResourceApi {
  static readonly serviceName: string;
  static readonly ListFederatedAuthConfigs: FederatedEnterpriseGlooResourceApiListFederatedAuthConfigs;
  static readonly GetFederatedAuthConfigYaml: FederatedEnterpriseGlooResourceApiGetFederatedAuthConfigYaml;
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

export class FederatedEnterpriseGlooResourceApiClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  listFederatedAuthConfigs(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_enterprise_gloo_resources_pb.ListFederatedAuthConfigsRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_enterprise_gloo_resources_pb.ListFederatedAuthConfigsResponse|null) => void
  ): UnaryResponse;
  listFederatedAuthConfigs(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_enterprise_gloo_resources_pb.ListFederatedAuthConfigsRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_enterprise_gloo_resources_pb.ListFederatedAuthConfigsResponse|null) => void
  ): UnaryResponse;
  getFederatedAuthConfigYaml(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_enterprise_gloo_resources_pb.GetFederatedAuthConfigYamlRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_enterprise_gloo_resources_pb.GetFederatedAuthConfigYamlResponse|null) => void
  ): UnaryResponse;
  getFederatedAuthConfigYaml(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_enterprise_gloo_resources_pb.GetFederatedAuthConfigYamlRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_enterprise_gloo_resources_pb.GetFederatedAuthConfigYamlResponse|null) => void
  ): UnaryResponse;
}

