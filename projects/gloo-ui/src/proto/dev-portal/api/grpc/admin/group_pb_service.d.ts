// package: admin.devportal.solo.io
// file: dev-portal/api/grpc/admin/group.proto

import * as dev_portal_api_grpc_admin_group_pb from "../../../../dev-portal/api/grpc/admin/group_pb";
import * as google_protobuf_empty_pb from "google-protobuf/google/protobuf/empty_pb";
import * as dev_portal_api_dev_portal_v1_common_pb from "../../../../dev-portal/api/dev-portal/v1/common_pb";
import {grpc} from "@improbable-eng/grpc-web";

type GroupApiGetGroup = {
  readonly methodName: string;
  readonly service: typeof GroupApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof dev_portal_api_dev_portal_v1_common_pb.ObjectRef;
  readonly responseType: typeof dev_portal_api_grpc_admin_group_pb.Group;
};

type GroupApiListGroups = {
  readonly methodName: string;
  readonly service: typeof GroupApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof dev_portal_api_grpc_admin_group_pb.GroupFilter;
  readonly responseType: typeof dev_portal_api_grpc_admin_group_pb.GroupList;
};

type GroupApiCreateGroup = {
  readonly methodName: string;
  readonly service: typeof GroupApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof dev_portal_api_grpc_admin_group_pb.GroupWriteRequest;
  readonly responseType: typeof dev_portal_api_grpc_admin_group_pb.Group;
};

type GroupApiUpdateGroup = {
  readonly methodName: string;
  readonly service: typeof GroupApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof dev_portal_api_grpc_admin_group_pb.GroupWriteRequest;
  readonly responseType: typeof dev_portal_api_grpc_admin_group_pb.Group;
};

type GroupApiDeleteGroup = {
  readonly methodName: string;
  readonly service: typeof GroupApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof dev_portal_api_dev_portal_v1_common_pb.ObjectRef;
  readonly responseType: typeof google_protobuf_empty_pb.Empty;
};

export class GroupApi {
  static readonly serviceName: string;
  static readonly GetGroup: GroupApiGetGroup;
  static readonly ListGroups: GroupApiListGroups;
  static readonly CreateGroup: GroupApiCreateGroup;
  static readonly UpdateGroup: GroupApiUpdateGroup;
  static readonly DeleteGroup: GroupApiDeleteGroup;
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

export class GroupApiClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  getGroup(
    requestMessage: dev_portal_api_dev_portal_v1_common_pb.ObjectRef,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: dev_portal_api_grpc_admin_group_pb.Group|null) => void
  ): UnaryResponse;
  getGroup(
    requestMessage: dev_portal_api_dev_portal_v1_common_pb.ObjectRef,
    callback: (error: ServiceError|null, responseMessage: dev_portal_api_grpc_admin_group_pb.Group|null) => void
  ): UnaryResponse;
  listGroups(
    requestMessage: dev_portal_api_grpc_admin_group_pb.GroupFilter,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: dev_portal_api_grpc_admin_group_pb.GroupList|null) => void
  ): UnaryResponse;
  listGroups(
    requestMessage: dev_portal_api_grpc_admin_group_pb.GroupFilter,
    callback: (error: ServiceError|null, responseMessage: dev_portal_api_grpc_admin_group_pb.GroupList|null) => void
  ): UnaryResponse;
  createGroup(
    requestMessage: dev_portal_api_grpc_admin_group_pb.GroupWriteRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: dev_portal_api_grpc_admin_group_pb.Group|null) => void
  ): UnaryResponse;
  createGroup(
    requestMessage: dev_portal_api_grpc_admin_group_pb.GroupWriteRequest,
    callback: (error: ServiceError|null, responseMessage: dev_portal_api_grpc_admin_group_pb.Group|null) => void
  ): UnaryResponse;
  updateGroup(
    requestMessage: dev_portal_api_grpc_admin_group_pb.GroupWriteRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: dev_portal_api_grpc_admin_group_pb.Group|null) => void
  ): UnaryResponse;
  updateGroup(
    requestMessage: dev_portal_api_grpc_admin_group_pb.GroupWriteRequest,
    callback: (error: ServiceError|null, responseMessage: dev_portal_api_grpc_admin_group_pb.Group|null) => void
  ): UnaryResponse;
  deleteGroup(
    requestMessage: dev_portal_api_dev_portal_v1_common_pb.ObjectRef,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: google_protobuf_empty_pb.Empty|null) => void
  ): UnaryResponse;
  deleteGroup(
    requestMessage: dev_portal_api_dev_portal_v1_common_pb.ObjectRef,
    callback: (error: ServiceError|null, responseMessage: google_protobuf_empty_pb.Empty|null) => void
  ): UnaryResponse;
}

