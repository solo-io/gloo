// package: fed.rpc.solo.io
// file: github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/gloo_resources.proto

import * as github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_gloo_resources_pb from "../../../../../../../../github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/gloo_resources_pb";
import {grpc} from "@improbable-eng/grpc-web";

type GlooResourceApiListUpstreams = {
  readonly methodName: string;
  readonly service: typeof GlooResourceApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_gloo_resources_pb.ListUpstreamsRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_gloo_resources_pb.ListUpstreamsResponse;
};

type GlooResourceApiGetUpstreamYaml = {
  readonly methodName: string;
  readonly service: typeof GlooResourceApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_gloo_resources_pb.GetUpstreamYamlRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_gloo_resources_pb.GetUpstreamYamlResponse;
};

type GlooResourceApiListUpstreamGroups = {
  readonly methodName: string;
  readonly service: typeof GlooResourceApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_gloo_resources_pb.ListUpstreamGroupsRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_gloo_resources_pb.ListUpstreamGroupsResponse;
};

type GlooResourceApiGetUpstreamGroupYaml = {
  readonly methodName: string;
  readonly service: typeof GlooResourceApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_gloo_resources_pb.GetUpstreamGroupYamlRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_gloo_resources_pb.GetUpstreamGroupYamlResponse;
};

type GlooResourceApiListSettings = {
  readonly methodName: string;
  readonly service: typeof GlooResourceApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_gloo_resources_pb.ListSettingsRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_gloo_resources_pb.ListSettingsResponse;
};

type GlooResourceApiGetSettingsYaml = {
  readonly methodName: string;
  readonly service: typeof GlooResourceApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_gloo_resources_pb.GetSettingsYamlRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_gloo_resources_pb.GetSettingsYamlResponse;
};

type GlooResourceApiListProxies = {
  readonly methodName: string;
  readonly service: typeof GlooResourceApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_gloo_resources_pb.ListProxiesRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_gloo_resources_pb.ListProxiesResponse;
};

type GlooResourceApiGetProxyYaml = {
  readonly methodName: string;
  readonly service: typeof GlooResourceApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_gloo_resources_pb.GetProxyYamlRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_gloo_resources_pb.GetProxyYamlResponse;
};

export class GlooResourceApi {
  static readonly serviceName: string;
  static readonly ListUpstreams: GlooResourceApiListUpstreams;
  static readonly GetUpstreamYaml: GlooResourceApiGetUpstreamYaml;
  static readonly ListUpstreamGroups: GlooResourceApiListUpstreamGroups;
  static readonly GetUpstreamGroupYaml: GlooResourceApiGetUpstreamGroupYaml;
  static readonly ListSettings: GlooResourceApiListSettings;
  static readonly GetSettingsYaml: GlooResourceApiGetSettingsYaml;
  static readonly ListProxies: GlooResourceApiListProxies;
  static readonly GetProxyYaml: GlooResourceApiGetProxyYaml;
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
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_gloo_resources_pb.ListUpstreamsRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_gloo_resources_pb.ListUpstreamsResponse|null) => void
  ): UnaryResponse;
  listUpstreams(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_gloo_resources_pb.ListUpstreamsRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_gloo_resources_pb.ListUpstreamsResponse|null) => void
  ): UnaryResponse;
  getUpstreamYaml(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_gloo_resources_pb.GetUpstreamYamlRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_gloo_resources_pb.GetUpstreamYamlResponse|null) => void
  ): UnaryResponse;
  getUpstreamYaml(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_gloo_resources_pb.GetUpstreamYamlRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_gloo_resources_pb.GetUpstreamYamlResponse|null) => void
  ): UnaryResponse;
  listUpstreamGroups(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_gloo_resources_pb.ListUpstreamGroupsRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_gloo_resources_pb.ListUpstreamGroupsResponse|null) => void
  ): UnaryResponse;
  listUpstreamGroups(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_gloo_resources_pb.ListUpstreamGroupsRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_gloo_resources_pb.ListUpstreamGroupsResponse|null) => void
  ): UnaryResponse;
  getUpstreamGroupYaml(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_gloo_resources_pb.GetUpstreamGroupYamlRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_gloo_resources_pb.GetUpstreamGroupYamlResponse|null) => void
  ): UnaryResponse;
  getUpstreamGroupYaml(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_gloo_resources_pb.GetUpstreamGroupYamlRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_gloo_resources_pb.GetUpstreamGroupYamlResponse|null) => void
  ): UnaryResponse;
  listSettings(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_gloo_resources_pb.ListSettingsRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_gloo_resources_pb.ListSettingsResponse|null) => void
  ): UnaryResponse;
  listSettings(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_gloo_resources_pb.ListSettingsRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_gloo_resources_pb.ListSettingsResponse|null) => void
  ): UnaryResponse;
  getSettingsYaml(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_gloo_resources_pb.GetSettingsYamlRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_gloo_resources_pb.GetSettingsYamlResponse|null) => void
  ): UnaryResponse;
  getSettingsYaml(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_gloo_resources_pb.GetSettingsYamlRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_gloo_resources_pb.GetSettingsYamlResponse|null) => void
  ): UnaryResponse;
  listProxies(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_gloo_resources_pb.ListProxiesRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_gloo_resources_pb.ListProxiesResponse|null) => void
  ): UnaryResponse;
  listProxies(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_gloo_resources_pb.ListProxiesRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_gloo_resources_pb.ListProxiesResponse|null) => void
  ): UnaryResponse;
  getProxyYaml(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_gloo_resources_pb.GetProxyYamlRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_gloo_resources_pb.GetProxyYamlResponse|null) => void
  ): UnaryResponse;
  getProxyYaml(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_gloo_resources_pb.GetProxyYamlRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_gloo_resources_pb.GetProxyYamlResponse|null) => void
  ): UnaryResponse;
}

