// package: glooeeapi.solo.io
// file: github.com/solo-io/solo-projects/projects/grpcserver/api/v1/upstream.proto

import * as github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstream_pb from "../../../../../../../github.com/solo-io/solo-projects/projects/grpcserver/api/v1/upstream_pb";
import {grpc} from "@improbable-eng/grpc-web";

type UpstreamApiGetUpstream = {
  readonly methodName: string;
  readonly service: typeof UpstreamApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstream_pb.GetUpstreamRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstream_pb.GetUpstreamResponse;
};

type UpstreamApiListUpstreams = {
  readonly methodName: string;
  readonly service: typeof UpstreamApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstream_pb.ListUpstreamsRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstream_pb.ListUpstreamsResponse;
};

type UpstreamApiStreamUpstreamList = {
  readonly methodName: string;
  readonly service: typeof UpstreamApi;
  readonly requestStream: false;
  readonly responseStream: true;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstream_pb.StreamUpstreamListRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstream_pb.StreamUpstreamListResponse;
};

type UpstreamApiCreateUpstream = {
  readonly methodName: string;
  readonly service: typeof UpstreamApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstream_pb.CreateUpstreamRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstream_pb.CreateUpstreamResponse;
};

type UpstreamApiUpdateUpstream = {
  readonly methodName: string;
  readonly service: typeof UpstreamApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstream_pb.UpdateUpstreamRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstream_pb.UpdateUpstreamResponse;
};

type UpstreamApiDeleteUpstream = {
  readonly methodName: string;
  readonly service: typeof UpstreamApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstream_pb.DeleteUpstreamRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstream_pb.DeleteUpstreamResponse;
};

export class UpstreamApi {
  static readonly serviceName: string;
  static readonly GetUpstream: UpstreamApiGetUpstream;
  static readonly ListUpstreams: UpstreamApiListUpstreams;
  static readonly StreamUpstreamList: UpstreamApiStreamUpstreamList;
  static readonly CreateUpstream: UpstreamApiCreateUpstream;
  static readonly UpdateUpstream: UpstreamApiUpdateUpstream;
  static readonly DeleteUpstream: UpstreamApiDeleteUpstream;
}

export type ServiceError = { message: string, code: number; metadata: grpc.Metadata }
export type Status = { details: string, code: number; metadata: grpc.Metadata }

interface UnaryResponse {
  cancel(): void;
}
interface ResponseStream<T> {
  cancel(): void;
  on(type: 'data', handler: (message: T) => void): ResponseStream<T>;
  on(type: 'end', handler: () => void): ResponseStream<T>;
  on(type: 'status', handler: (status: Status) => void): ResponseStream<T>;
}
interface RequestStream<T> {
  write(message: T): RequestStream<T>;
  end(): void;
  cancel(): void;
  on(type: 'end', handler: () => void): RequestStream<T>;
  on(type: 'status', handler: (status: Status) => void): RequestStream<T>;
}
interface BidirectionalStream<ReqT, ResT> {
  write(message: ReqT): BidirectionalStream<ReqT, ResT>;
  end(): void;
  cancel(): void;
  on(type: 'data', handler: (message: ResT) => void): BidirectionalStream<ReqT, ResT>;
  on(type: 'end', handler: () => void): BidirectionalStream<ReqT, ResT>;
  on(type: 'status', handler: (status: Status) => void): BidirectionalStream<ReqT, ResT>;
}

export class UpstreamApiClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  getUpstream(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstream_pb.GetUpstreamRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstream_pb.GetUpstreamResponse|null) => void
  ): UnaryResponse;
  getUpstream(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstream_pb.GetUpstreamRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstream_pb.GetUpstreamResponse|null) => void
  ): UnaryResponse;
  listUpstreams(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstream_pb.ListUpstreamsRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstream_pb.ListUpstreamsResponse|null) => void
  ): UnaryResponse;
  listUpstreams(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstream_pb.ListUpstreamsRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstream_pb.ListUpstreamsResponse|null) => void
  ): UnaryResponse;
  streamUpstreamList(requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstream_pb.StreamUpstreamListRequest, metadata?: grpc.Metadata): ResponseStream<github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstream_pb.StreamUpstreamListResponse>;
  createUpstream(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstream_pb.CreateUpstreamRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstream_pb.CreateUpstreamResponse|null) => void
  ): UnaryResponse;
  createUpstream(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstream_pb.CreateUpstreamRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstream_pb.CreateUpstreamResponse|null) => void
  ): UnaryResponse;
  updateUpstream(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstream_pb.UpdateUpstreamRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstream_pb.UpdateUpstreamResponse|null) => void
  ): UnaryResponse;
  updateUpstream(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstream_pb.UpdateUpstreamRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstream_pb.UpdateUpstreamResponse|null) => void
  ): UnaryResponse;
  deleteUpstream(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstream_pb.DeleteUpstreamRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstream_pb.DeleteUpstreamResponse|null) => void
  ): UnaryResponse;
  deleteUpstream(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstream_pb.DeleteUpstreamRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstream_pb.DeleteUpstreamResponse|null) => void
  ): UnaryResponse;
}

