// package: admin.devportal.solo.io
// file: dev-portal/api/grpc/admin/user.proto

import * as dev_portal_api_grpc_admin_user_pb from "../../../../dev-portal/api/grpc/admin/user_pb";
import * as google_protobuf_empty_pb from "google-protobuf/google/protobuf/empty_pb";
import * as dev_portal_api_dev_portal_v1_common_pb from "../../../../dev-portal/api/dev-portal/v1/common_pb";
import {grpc} from "@improbable-eng/grpc-web";

type UserApiGetUser = {
  readonly methodName: string;
  readonly service: typeof UserApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof dev_portal_api_dev_portal_v1_common_pb.ObjectRef;
  readonly responseType: typeof dev_portal_api_grpc_admin_user_pb.User;
};

type UserApiListUsers = {
  readonly methodName: string;
  readonly service: typeof UserApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof dev_portal_api_grpc_admin_user_pb.UserFilter;
  readonly responseType: typeof dev_portal_api_grpc_admin_user_pb.UserList;
};

type UserApiCreateUser = {
  readonly methodName: string;
  readonly service: typeof UserApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof dev_portal_api_grpc_admin_user_pb.UserWriteRequest;
  readonly responseType: typeof dev_portal_api_grpc_admin_user_pb.User;
};

type UserApiUpdateUser = {
  readonly methodName: string;
  readonly service: typeof UserApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof dev_portal_api_grpc_admin_user_pb.UserWriteRequest;
  readonly responseType: typeof dev_portal_api_grpc_admin_user_pb.User;
};

type UserApiDeleteUser = {
  readonly methodName: string;
  readonly service: typeof UserApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof dev_portal_api_dev_portal_v1_common_pb.ObjectRef;
  readonly responseType: typeof google_protobuf_empty_pb.Empty;
};

export class UserApi {
  static readonly serviceName: string;
  static readonly GetUser: UserApiGetUser;
  static readonly ListUsers: UserApiListUsers;
  static readonly CreateUser: UserApiCreateUser;
  static readonly UpdateUser: UserApiUpdateUser;
  static readonly DeleteUser: UserApiDeleteUser;
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

export class UserApiClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  getUser(
    requestMessage: dev_portal_api_dev_portal_v1_common_pb.ObjectRef,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: dev_portal_api_grpc_admin_user_pb.User|null) => void
  ): UnaryResponse;
  getUser(
    requestMessage: dev_portal_api_dev_portal_v1_common_pb.ObjectRef,
    callback: (error: ServiceError|null, responseMessage: dev_portal_api_grpc_admin_user_pb.User|null) => void
  ): UnaryResponse;
  listUsers(
    requestMessage: dev_portal_api_grpc_admin_user_pb.UserFilter,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: dev_portal_api_grpc_admin_user_pb.UserList|null) => void
  ): UnaryResponse;
  listUsers(
    requestMessage: dev_portal_api_grpc_admin_user_pb.UserFilter,
    callback: (error: ServiceError|null, responseMessage: dev_portal_api_grpc_admin_user_pb.UserList|null) => void
  ): UnaryResponse;
  createUser(
    requestMessage: dev_portal_api_grpc_admin_user_pb.UserWriteRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: dev_portal_api_grpc_admin_user_pb.User|null) => void
  ): UnaryResponse;
  createUser(
    requestMessage: dev_portal_api_grpc_admin_user_pb.UserWriteRequest,
    callback: (error: ServiceError|null, responseMessage: dev_portal_api_grpc_admin_user_pb.User|null) => void
  ): UnaryResponse;
  updateUser(
    requestMessage: dev_portal_api_grpc_admin_user_pb.UserWriteRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: dev_portal_api_grpc_admin_user_pb.User|null) => void
  ): UnaryResponse;
  updateUser(
    requestMessage: dev_portal_api_grpc_admin_user_pb.UserWriteRequest,
    callback: (error: ServiceError|null, responseMessage: dev_portal_api_grpc_admin_user_pb.User|null) => void
  ): UnaryResponse;
  deleteUser(
    requestMessage: dev_portal_api_dev_portal_v1_common_pb.ObjectRef,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: google_protobuf_empty_pb.Empty|null) => void
  ): UnaryResponse;
  deleteUser(
    requestMessage: dev_portal_api_dev_portal_v1_common_pb.ObjectRef,
    callback: (error: ServiceError|null, responseMessage: google_protobuf_empty_pb.Empty|null) => void
  ): UnaryResponse;
}

