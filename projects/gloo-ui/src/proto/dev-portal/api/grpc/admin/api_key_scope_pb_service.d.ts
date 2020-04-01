// package: admin.devportal.solo.io
// file: dev-portal/api/grpc/admin/api_key_scope.proto

import * as dev_portal_api_grpc_admin_api_key_scope_pb from "../../../../dev-portal/api/grpc/admin/api_key_scope_pb";
import * as google_protobuf_empty_pb from "google-protobuf/google/protobuf/empty_pb";
import {grpc} from "@improbable-eng/grpc-web";

type ApiKeyScopeApiListApiKeyScopes = {
  readonly methodName: string;
  readonly service: typeof ApiKeyScopeApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof dev_portal_api_grpc_admin_api_key_scope_pb.ApiKeyScopeFilter;
  readonly responseType: typeof dev_portal_api_grpc_admin_api_key_scope_pb.ApiKeyScopeList;
};

type ApiKeyScopeApiCreateApiKeyScope = {
  readonly methodName: string;
  readonly service: typeof ApiKeyScopeApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof dev_portal_api_grpc_admin_api_key_scope_pb.ApiKeyScopeWriteRequest;
  readonly responseType: typeof dev_portal_api_grpc_admin_api_key_scope_pb.ApiKeyScope;
};

type ApiKeyScopeApiUpdateApiKeyScope = {
  readonly methodName: string;
  readonly service: typeof ApiKeyScopeApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof dev_portal_api_grpc_admin_api_key_scope_pb.ApiKeyScopeWriteRequest;
  readonly responseType: typeof dev_portal_api_grpc_admin_api_key_scope_pb.ApiKeyScope;
};

type ApiKeyScopeApiDeleteApiKeyScope = {
  readonly methodName: string;
  readonly service: typeof ApiKeyScopeApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof dev_portal_api_grpc_admin_api_key_scope_pb.ApiKeyScopeRef;
  readonly responseType: typeof google_protobuf_empty_pb.Empty;
};

export class ApiKeyScopeApi {
  static readonly serviceName: string;
  static readonly ListApiKeyScopes: ApiKeyScopeApiListApiKeyScopes;
  static readonly CreateApiKeyScope: ApiKeyScopeApiCreateApiKeyScope;
  static readonly UpdateApiKeyScope: ApiKeyScopeApiUpdateApiKeyScope;
  static readonly DeleteApiKeyScope: ApiKeyScopeApiDeleteApiKeyScope;
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

export class ApiKeyScopeApiClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  listApiKeyScopes(
    requestMessage: dev_portal_api_grpc_admin_api_key_scope_pb.ApiKeyScopeFilter,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: dev_portal_api_grpc_admin_api_key_scope_pb.ApiKeyScopeList|null) => void
  ): UnaryResponse;
  listApiKeyScopes(
    requestMessage: dev_portal_api_grpc_admin_api_key_scope_pb.ApiKeyScopeFilter,
    callback: (error: ServiceError|null, responseMessage: dev_portal_api_grpc_admin_api_key_scope_pb.ApiKeyScopeList|null) => void
  ): UnaryResponse;
  createApiKeyScope(
    requestMessage: dev_portal_api_grpc_admin_api_key_scope_pb.ApiKeyScopeWriteRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: dev_portal_api_grpc_admin_api_key_scope_pb.ApiKeyScope|null) => void
  ): UnaryResponse;
  createApiKeyScope(
    requestMessage: dev_portal_api_grpc_admin_api_key_scope_pb.ApiKeyScopeWriteRequest,
    callback: (error: ServiceError|null, responseMessage: dev_portal_api_grpc_admin_api_key_scope_pb.ApiKeyScope|null) => void
  ): UnaryResponse;
  updateApiKeyScope(
    requestMessage: dev_portal_api_grpc_admin_api_key_scope_pb.ApiKeyScopeWriteRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: dev_portal_api_grpc_admin_api_key_scope_pb.ApiKeyScope|null) => void
  ): UnaryResponse;
  updateApiKeyScope(
    requestMessage: dev_portal_api_grpc_admin_api_key_scope_pb.ApiKeyScopeWriteRequest,
    callback: (error: ServiceError|null, responseMessage: dev_portal_api_grpc_admin_api_key_scope_pb.ApiKeyScope|null) => void
  ): UnaryResponse;
  deleteApiKeyScope(
    requestMessage: dev_portal_api_grpc_admin_api_key_scope_pb.ApiKeyScopeRef,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: google_protobuf_empty_pb.Empty|null) => void
  ): UnaryResponse;
  deleteApiKeyScope(
    requestMessage: dev_portal_api_grpc_admin_api_key_scope_pb.ApiKeyScopeRef,
    callback: (error: ServiceError|null, responseMessage: google_protobuf_empty_pb.Empty|null) => void
  ): UnaryResponse;
}

