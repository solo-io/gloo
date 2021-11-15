// package: rpc.edge.gloo.solo.io
// file: github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/glooinstance.proto

import * as github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_glooinstance_pb from "../../../../../../../../github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/glooinstance_pb";
import {grpc} from "@improbable-eng/grpc-web";

type GlooInstanceApiListGlooInstances = {
  readonly methodName: string;
  readonly service: typeof GlooInstanceApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_glooinstance_pb.ListGlooInstancesRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_glooinstance_pb.ListGlooInstancesResponse;
};

type GlooInstanceApiListClusterDetails = {
  readonly methodName: string;
  readonly service: typeof GlooInstanceApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_glooinstance_pb.ListClusterDetailsRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_glooinstance_pb.ListClusterDetailsResponse;
};

type GlooInstanceApiGetConfigDumps = {
  readonly methodName: string;
  readonly service: typeof GlooInstanceApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_glooinstance_pb.GetConfigDumpsRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_glooinstance_pb.GetConfigDumpsResponse;
};

type GlooInstanceApiGetUpstreamHosts = {
  readonly methodName: string;
  readonly service: typeof GlooInstanceApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_glooinstance_pb.GetUpstreamHostsRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_glooinstance_pb.GetUpstreamHostsResponse;
};

export class GlooInstanceApi {
  static readonly serviceName: string;
  static readonly ListGlooInstances: GlooInstanceApiListGlooInstances;
  static readonly ListClusterDetails: GlooInstanceApiListClusterDetails;
  static readonly GetConfigDumps: GlooInstanceApiGetConfigDumps;
  static readonly GetUpstreamHosts: GlooInstanceApiGetUpstreamHosts;
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

export class GlooInstanceApiClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  listGlooInstances(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_glooinstance_pb.ListGlooInstancesRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_glooinstance_pb.ListGlooInstancesResponse|null) => void
  ): UnaryResponse;
  listGlooInstances(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_glooinstance_pb.ListGlooInstancesRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_glooinstance_pb.ListGlooInstancesResponse|null) => void
  ): UnaryResponse;
  listClusterDetails(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_glooinstance_pb.ListClusterDetailsRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_glooinstance_pb.ListClusterDetailsResponse|null) => void
  ): UnaryResponse;
  listClusterDetails(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_glooinstance_pb.ListClusterDetailsRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_glooinstance_pb.ListClusterDetailsResponse|null) => void
  ): UnaryResponse;
  getConfigDumps(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_glooinstance_pb.GetConfigDumpsRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_glooinstance_pb.GetConfigDumpsResponse|null) => void
  ): UnaryResponse;
  getConfigDumps(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_glooinstance_pb.GetConfigDumpsRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_glooinstance_pb.GetConfigDumpsResponse|null) => void
  ): UnaryResponse;
  getUpstreamHosts(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_glooinstance_pb.GetUpstreamHostsRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_glooinstance_pb.GetUpstreamHostsResponse|null) => void
  ): UnaryResponse;
  getUpstreamHosts(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_glooinstance_pb.GetUpstreamHostsRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_glooinstance_pb.GetUpstreamHostsResponse|null) => void
  ): UnaryResponse;
}

