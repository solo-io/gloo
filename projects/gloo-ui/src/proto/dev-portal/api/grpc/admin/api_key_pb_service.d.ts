// package: admin.devportal.solo.io
// file: dev-portal/api/grpc/admin/api_key.proto

import * as dev_portal_api_grpc_admin_api_key_pb from "../../../../dev-portal/api/grpc/admin/api_key_pb";
import * as google_protobuf_empty_pb from "google-protobuf/google/protobuf/empty_pb";
import * as dev_portal_api_dev_portal_v1_common_pb from "../../../../dev-portal/api/dev-portal/v1/common_pb";
import {grpc} from "@improbable-eng/grpc-web";

type ApiKeyApiListApiKeys = {
  readonly methodName: string;
  readonly service: typeof ApiKeyApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof google_protobuf_empty_pb.Empty;
  readonly responseType: typeof dev_portal_api_grpc_admin_api_key_pb.ApiKeyList;
};

type ApiKeyApiDeleteApiKey = {
  readonly methodName: string;
  readonly service: typeof ApiKeyApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof dev_portal_api_dev_portal_v1_common_pb.ObjectRef;
  readonly responseType: typeof google_protobuf_empty_pb.Empty;
};

export class ApiKeyApi {
  static readonly serviceName: string;
  static readonly ListApiKeys: ApiKeyApiListApiKeys;
  static readonly DeleteApiKey: ApiKeyApiDeleteApiKey;
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

export class ApiKeyApiClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  listApiKeys(
    requestMessage: google_protobuf_empty_pb.Empty,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: dev_portal_api_grpc_admin_api_key_pb.ApiKeyList|null) => void
  ): UnaryResponse;
  listApiKeys(
    requestMessage: google_protobuf_empty_pb.Empty,
    callback: (error: ServiceError|null, responseMessage: dev_portal_api_grpc_admin_api_key_pb.ApiKeyList|null) => void
  ): UnaryResponse;
  deleteApiKey(
    requestMessage: dev_portal_api_dev_portal_v1_common_pb.ObjectRef,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: google_protobuf_empty_pb.Empty|null) => void
  ): UnaryResponse;
  deleteApiKey(
    requestMessage: dev_portal_api_dev_portal_v1_common_pb.ObjectRef,
    callback: (error: ServiceError|null, responseMessage: google_protobuf_empty_pb.Empty|null) => void
  ): UnaryResponse;
}

