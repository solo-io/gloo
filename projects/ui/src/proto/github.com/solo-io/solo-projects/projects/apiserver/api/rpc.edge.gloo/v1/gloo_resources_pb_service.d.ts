// package: rpc.edge.gloo.solo.io
// file: github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/gloo_resources.proto

import * as github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb from "../../../../../../../../github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/gloo_resources_pb";
import {grpc} from "@improbable-eng/grpc-web";

type GlooResourceApiListUpstreams = {
  readonly methodName: string;
  readonly service: typeof GlooResourceApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.ListUpstreamsRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.ListUpstreamsResponse;
};

type GlooResourceApiGetUpstreamYaml = {
  readonly methodName: string;
  readonly service: typeof GlooResourceApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.GetUpstreamYamlRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.GetUpstreamYamlResponse;
};

type GlooResourceApiGetUpstreamDetails = {
  readonly methodName: string;
  readonly service: typeof GlooResourceApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.GetUpstreamDetailsRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.GetUpstreamDetailsResponse;
};

type GlooResourceApiListUpstreamGroups = {
  readonly methodName: string;
  readonly service: typeof GlooResourceApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.ListUpstreamGroupsRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.ListUpstreamGroupsResponse;
};

type GlooResourceApiGetUpstreamGroupYaml = {
  readonly methodName: string;
  readonly service: typeof GlooResourceApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.GetUpstreamGroupYamlRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.GetUpstreamGroupYamlResponse;
};

type GlooResourceApiGetUpstreamGroupDetails = {
  readonly methodName: string;
  readonly service: typeof GlooResourceApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.GetUpstreamGroupDetailsRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.GetUpstreamGroupDetailsResponse;
};

type GlooResourceApiListSettings = {
  readonly methodName: string;
  readonly service: typeof GlooResourceApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.ListSettingsRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.ListSettingsResponse;
};

type GlooResourceApiGetSettingsYaml = {
  readonly methodName: string;
  readonly service: typeof GlooResourceApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.GetSettingsYamlRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.GetSettingsYamlResponse;
};

type GlooResourceApiGetSettingsDetails = {
  readonly methodName: string;
  readonly service: typeof GlooResourceApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.GetSettingsDetailsRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.GetSettingsDetailsResponse;
};

type GlooResourceApiListProxies = {
  readonly methodName: string;
  readonly service: typeof GlooResourceApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.ListProxiesRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.ListProxiesResponse;
};

type GlooResourceApiGetProxyYaml = {
  readonly methodName: string;
  readonly service: typeof GlooResourceApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.GetProxyYamlRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.GetProxyYamlResponse;
};

type GlooResourceApiGetProxyDetails = {
  readonly methodName: string;
  readonly service: typeof GlooResourceApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.GetProxyDetailsRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.GetProxyDetailsResponse;
};

export class GlooResourceApi {
  static readonly serviceName: string;
  static readonly ListUpstreams: GlooResourceApiListUpstreams;
  static readonly GetUpstreamYaml: GlooResourceApiGetUpstreamYaml;
  static readonly GetUpstreamDetails: GlooResourceApiGetUpstreamDetails;
  static readonly ListUpstreamGroups: GlooResourceApiListUpstreamGroups;
  static readonly GetUpstreamGroupYaml: GlooResourceApiGetUpstreamGroupYaml;
  static readonly GetUpstreamGroupDetails: GlooResourceApiGetUpstreamGroupDetails;
  static readonly ListSettings: GlooResourceApiListSettings;
  static readonly GetSettingsYaml: GlooResourceApiGetSettingsYaml;
  static readonly GetSettingsDetails: GlooResourceApiGetSettingsDetails;
  static readonly ListProxies: GlooResourceApiListProxies;
  static readonly GetProxyYaml: GlooResourceApiGetProxyYaml;
  static readonly GetProxyDetails: GlooResourceApiGetProxyDetails;
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

export class GlooResourceApiClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  listUpstreams(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.ListUpstreamsRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.ListUpstreamsResponse|null) => void
  ): UnaryResponse;
  listUpstreams(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.ListUpstreamsRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.ListUpstreamsResponse|null) => void
  ): UnaryResponse;
  getUpstreamYaml(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.GetUpstreamYamlRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.GetUpstreamYamlResponse|null) => void
  ): UnaryResponse;
  getUpstreamYaml(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.GetUpstreamYamlRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.GetUpstreamYamlResponse|null) => void
  ): UnaryResponse;
  getUpstreamDetails(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.GetUpstreamDetailsRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.GetUpstreamDetailsResponse|null) => void
  ): UnaryResponse;
  getUpstreamDetails(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.GetUpstreamDetailsRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.GetUpstreamDetailsResponse|null) => void
  ): UnaryResponse;
  listUpstreamGroups(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.ListUpstreamGroupsRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.ListUpstreamGroupsResponse|null) => void
  ): UnaryResponse;
  listUpstreamGroups(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.ListUpstreamGroupsRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.ListUpstreamGroupsResponse|null) => void
  ): UnaryResponse;
  getUpstreamGroupYaml(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.GetUpstreamGroupYamlRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.GetUpstreamGroupYamlResponse|null) => void
  ): UnaryResponse;
  getUpstreamGroupYaml(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.GetUpstreamGroupYamlRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.GetUpstreamGroupYamlResponse|null) => void
  ): UnaryResponse;
  getUpstreamGroupDetails(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.GetUpstreamGroupDetailsRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.GetUpstreamGroupDetailsResponse|null) => void
  ): UnaryResponse;
  getUpstreamGroupDetails(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.GetUpstreamGroupDetailsRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.GetUpstreamGroupDetailsResponse|null) => void
  ): UnaryResponse;
  listSettings(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.ListSettingsRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.ListSettingsResponse|null) => void
  ): UnaryResponse;
  listSettings(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.ListSettingsRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.ListSettingsResponse|null) => void
  ): UnaryResponse;
  getSettingsYaml(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.GetSettingsYamlRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.GetSettingsYamlResponse|null) => void
  ): UnaryResponse;
  getSettingsYaml(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.GetSettingsYamlRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.GetSettingsYamlResponse|null) => void
  ): UnaryResponse;
  getSettingsDetails(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.GetSettingsDetailsRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.GetSettingsDetailsResponse|null) => void
  ): UnaryResponse;
  getSettingsDetails(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.GetSettingsDetailsRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.GetSettingsDetailsResponse|null) => void
  ): UnaryResponse;
  listProxies(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.ListProxiesRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.ListProxiesResponse|null) => void
  ): UnaryResponse;
  listProxies(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.ListProxiesRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.ListProxiesResponse|null) => void
  ): UnaryResponse;
  getProxyYaml(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.GetProxyYamlRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.GetProxyYamlResponse|null) => void
  ): UnaryResponse;
  getProxyYaml(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.GetProxyYamlRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.GetProxyYamlResponse|null) => void
  ): UnaryResponse;
  getProxyDetails(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.GetProxyDetailsRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.GetProxyDetailsResponse|null) => void
  ): UnaryResponse;
  getProxyDetails(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.GetProxyDetailsRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gloo_resources_pb.GetProxyDetailsResponse|null) => void
  ): UnaryResponse;
}

