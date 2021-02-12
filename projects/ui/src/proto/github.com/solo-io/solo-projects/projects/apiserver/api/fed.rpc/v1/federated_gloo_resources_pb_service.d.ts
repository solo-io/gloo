// package: fed.rpc.solo.io
// file: github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/federated_gloo_resources.proto

import * as github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gloo_resources_pb from "../../../../../../../../github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/federated_gloo_resources_pb";
import {grpc} from "@improbable-eng/grpc-web";

type FederatedGlooResourceApiListFederatedUpstreams = {
  readonly methodName: string;
  readonly service: typeof FederatedGlooResourceApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gloo_resources_pb.ListFederatedUpstreamsRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gloo_resources_pb.ListFederatedUpstreamsResponse;
};

type FederatedGlooResourceApiGetFederatedUpstreamYaml = {
  readonly methodName: string;
  readonly service: typeof FederatedGlooResourceApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gloo_resources_pb.GetFederatedUpstreamYamlRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gloo_resources_pb.GetFederatedUpstreamYamlResponse;
};

type FederatedGlooResourceApiListFederatedUpstreamGroups = {
  readonly methodName: string;
  readonly service: typeof FederatedGlooResourceApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gloo_resources_pb.ListFederatedUpstreamGroupsRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gloo_resources_pb.ListFederatedUpstreamGroupsResponse;
};

type FederatedGlooResourceApiGetFederatedUpstreamGroupYaml = {
  readonly methodName: string;
  readonly service: typeof FederatedGlooResourceApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gloo_resources_pb.GetFederatedUpstreamGroupYamlRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gloo_resources_pb.GetFederatedUpstreamGroupYamlResponse;
};

type FederatedGlooResourceApiListFederatedSettings = {
  readonly methodName: string;
  readonly service: typeof FederatedGlooResourceApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gloo_resources_pb.ListFederatedSettingsRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gloo_resources_pb.ListFederatedSettingsResponse;
};

type FederatedGlooResourceApiGetFederatedSettingsYaml = {
  readonly methodName: string;
  readonly service: typeof FederatedGlooResourceApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gloo_resources_pb.GetFederatedSettingsYamlRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gloo_resources_pb.GetFederatedSettingsYamlResponse;
};

export class FederatedGlooResourceApi {
  static readonly serviceName: string;
  static readonly ListFederatedUpstreams: FederatedGlooResourceApiListFederatedUpstreams;
  static readonly GetFederatedUpstreamYaml: FederatedGlooResourceApiGetFederatedUpstreamYaml;
  static readonly ListFederatedUpstreamGroups: FederatedGlooResourceApiListFederatedUpstreamGroups;
  static readonly GetFederatedUpstreamGroupYaml: FederatedGlooResourceApiGetFederatedUpstreamGroupYaml;
  static readonly ListFederatedSettings: FederatedGlooResourceApiListFederatedSettings;
  static readonly GetFederatedSettingsYaml: FederatedGlooResourceApiGetFederatedSettingsYaml;
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

export class FederatedGlooResourceApiClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  listFederatedUpstreams(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gloo_resources_pb.ListFederatedUpstreamsRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gloo_resources_pb.ListFederatedUpstreamsResponse|null) => void
  ): UnaryResponse;
  listFederatedUpstreams(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gloo_resources_pb.ListFederatedUpstreamsRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gloo_resources_pb.ListFederatedUpstreamsResponse|null) => void
  ): UnaryResponse;
  getFederatedUpstreamYaml(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gloo_resources_pb.GetFederatedUpstreamYamlRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gloo_resources_pb.GetFederatedUpstreamYamlResponse|null) => void
  ): UnaryResponse;
  getFederatedUpstreamYaml(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gloo_resources_pb.GetFederatedUpstreamYamlRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gloo_resources_pb.GetFederatedUpstreamYamlResponse|null) => void
  ): UnaryResponse;
  listFederatedUpstreamGroups(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gloo_resources_pb.ListFederatedUpstreamGroupsRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gloo_resources_pb.ListFederatedUpstreamGroupsResponse|null) => void
  ): UnaryResponse;
  listFederatedUpstreamGroups(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gloo_resources_pb.ListFederatedUpstreamGroupsRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gloo_resources_pb.ListFederatedUpstreamGroupsResponse|null) => void
  ): UnaryResponse;
  getFederatedUpstreamGroupYaml(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gloo_resources_pb.GetFederatedUpstreamGroupYamlRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gloo_resources_pb.GetFederatedUpstreamGroupYamlResponse|null) => void
  ): UnaryResponse;
  getFederatedUpstreamGroupYaml(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gloo_resources_pb.GetFederatedUpstreamGroupYamlRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gloo_resources_pb.GetFederatedUpstreamGroupYamlResponse|null) => void
  ): UnaryResponse;
  listFederatedSettings(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gloo_resources_pb.ListFederatedSettingsRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gloo_resources_pb.ListFederatedSettingsResponse|null) => void
  ): UnaryResponse;
  listFederatedSettings(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gloo_resources_pb.ListFederatedSettingsRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gloo_resources_pb.ListFederatedSettingsResponse|null) => void
  ): UnaryResponse;
  getFederatedSettingsYaml(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gloo_resources_pb.GetFederatedSettingsYamlRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gloo_resources_pb.GetFederatedSettingsYamlResponse|null) => void
  ): UnaryResponse;
  getFederatedSettingsYaml(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gloo_resources_pb.GetFederatedSettingsYamlRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gloo_resources_pb.GetFederatedSettingsYamlResponse|null) => void
  ): UnaryResponse;
}

