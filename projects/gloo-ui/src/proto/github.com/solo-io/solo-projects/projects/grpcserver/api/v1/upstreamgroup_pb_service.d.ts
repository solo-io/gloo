// package: glooeeapi.solo.io
// file: github.com/solo-io/solo-projects/projects/grpcserver/api/v1/upstreamgroup.proto

import * as github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstreamgroup_pb from "../../../../../../../github.com/solo-io/solo-projects/projects/grpcserver/api/v1/upstreamgroup_pb";
import {grpc} from "@improbable-eng/grpc-web";

type UpstreamGroupApiGetUpstreamGroup = {
  readonly methodName: string;
  readonly service: typeof UpstreamGroupApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstreamgroup_pb.GetUpstreamGroupRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstreamgroup_pb.GetUpstreamGroupResponse;
};

type UpstreamGroupApiListUpstreamGroups = {
  readonly methodName: string;
  readonly service: typeof UpstreamGroupApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstreamgroup_pb.ListUpstreamGroupsRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstreamgroup_pb.ListUpstreamGroupsResponse;
};

type UpstreamGroupApiCreateUpstreamGroup = {
  readonly methodName: string;
  readonly service: typeof UpstreamGroupApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstreamgroup_pb.CreateUpstreamGroupRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstreamgroup_pb.CreateUpstreamGroupResponse;
};

type UpstreamGroupApiUpdateUpstreamGroup = {
  readonly methodName: string;
  readonly service: typeof UpstreamGroupApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstreamgroup_pb.UpdateUpstreamGroupRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstreamgroup_pb.UpdateUpstreamGroupResponse;
};

type UpstreamGroupApiDeleteUpstreamGroup = {
  readonly methodName: string;
  readonly service: typeof UpstreamGroupApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstreamgroup_pb.DeleteUpstreamGroupRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstreamgroup_pb.DeleteUpstreamGroupResponse;
};

export class UpstreamGroupApi {
  static readonly serviceName: string;
  static readonly GetUpstreamGroup: UpstreamGroupApiGetUpstreamGroup;
  static readonly ListUpstreamGroups: UpstreamGroupApiListUpstreamGroups;
  static readonly CreateUpstreamGroup: UpstreamGroupApiCreateUpstreamGroup;
  static readonly UpdateUpstreamGroup: UpstreamGroupApiUpdateUpstreamGroup;
  static readonly DeleteUpstreamGroup: UpstreamGroupApiDeleteUpstreamGroup;
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

export class UpstreamGroupApiClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  getUpstreamGroup(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstreamgroup_pb.GetUpstreamGroupRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstreamgroup_pb.GetUpstreamGroupResponse|null) => void
  ): UnaryResponse;
  getUpstreamGroup(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstreamgroup_pb.GetUpstreamGroupRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstreamgroup_pb.GetUpstreamGroupResponse|null) => void
  ): UnaryResponse;
  listUpstreamGroups(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstreamgroup_pb.ListUpstreamGroupsRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstreamgroup_pb.ListUpstreamGroupsResponse|null) => void
  ): UnaryResponse;
  listUpstreamGroups(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstreamgroup_pb.ListUpstreamGroupsRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstreamgroup_pb.ListUpstreamGroupsResponse|null) => void
  ): UnaryResponse;
  createUpstreamGroup(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstreamgroup_pb.CreateUpstreamGroupRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstreamgroup_pb.CreateUpstreamGroupResponse|null) => void
  ): UnaryResponse;
  createUpstreamGroup(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstreamgroup_pb.CreateUpstreamGroupRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstreamgroup_pb.CreateUpstreamGroupResponse|null) => void
  ): UnaryResponse;
  updateUpstreamGroup(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstreamgroup_pb.UpdateUpstreamGroupRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstreamgroup_pb.UpdateUpstreamGroupResponse|null) => void
  ): UnaryResponse;
  updateUpstreamGroup(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstreamgroup_pb.UpdateUpstreamGroupRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstreamgroup_pb.UpdateUpstreamGroupResponse|null) => void
  ): UnaryResponse;
  deleteUpstreamGroup(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstreamgroup_pb.DeleteUpstreamGroupRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstreamgroup_pb.DeleteUpstreamGroupResponse|null) => void
  ): UnaryResponse;
  deleteUpstreamGroup(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstreamgroup_pb.DeleteUpstreamGroupRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstreamgroup_pb.DeleteUpstreamGroupResponse|null) => void
  ): UnaryResponse;
}

