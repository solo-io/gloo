// package: admin.devportal.solo.io
// file: dev-portal/api/grpc/admin/apidoc.proto

import * as dev_portal_api_grpc_admin_apidoc_pb from "../../../../dev-portal/api/grpc/admin/apidoc_pb";
import * as google_protobuf_empty_pb from "google-protobuf/google/protobuf/empty_pb";
import * as dev_portal_api_dev_portal_v1_common_pb from "../../../../dev-portal/api/dev-portal/v1/common_pb";
import {grpc} from "@improbable-eng/grpc-web";

type ApiDocApiGetApiDoc = {
  readonly methodName: string;
  readonly service: typeof ApiDocApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof dev_portal_api_grpc_admin_apidoc_pb.ApiDocGetRequest;
  readonly responseType: typeof dev_portal_api_grpc_admin_apidoc_pb.ApiDoc;
};

type ApiDocApiListApiDocs = {
  readonly methodName: string;
  readonly service: typeof ApiDocApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof dev_portal_api_grpc_admin_apidoc_pb.ApiDocFilter;
  readonly responseType: typeof dev_portal_api_grpc_admin_apidoc_pb.ApiDocList;
};

type ApiDocApiCreateApiDoc = {
  readonly methodName: string;
  readonly service: typeof ApiDocApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof dev_portal_api_grpc_admin_apidoc_pb.ApiDocWriteRequest;
  readonly responseType: typeof dev_portal_api_grpc_admin_apidoc_pb.ApiDoc;
};

type ApiDocApiUpdateApiDoc = {
  readonly methodName: string;
  readonly service: typeof ApiDocApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof dev_portal_api_grpc_admin_apidoc_pb.ApiDocWriteRequest;
  readonly responseType: typeof dev_portal_api_grpc_admin_apidoc_pb.ApiDoc;
};

type ApiDocApiDeleteApiDoc = {
  readonly methodName: string;
  readonly service: typeof ApiDocApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof dev_portal_api_dev_portal_v1_common_pb.ObjectRef;
  readonly responseType: typeof google_protobuf_empty_pb.Empty;
};

export class ApiDocApi {
  static readonly serviceName: string;
  static readonly GetApiDoc: ApiDocApiGetApiDoc;
  static readonly ListApiDocs: ApiDocApiListApiDocs;
  static readonly CreateApiDoc: ApiDocApiCreateApiDoc;
  static readonly UpdateApiDoc: ApiDocApiUpdateApiDoc;
  static readonly DeleteApiDoc: ApiDocApiDeleteApiDoc;
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

export class ApiDocApiClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  getApiDoc(
    requestMessage: dev_portal_api_grpc_admin_apidoc_pb.ApiDocGetRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: dev_portal_api_grpc_admin_apidoc_pb.ApiDoc|null) => void
  ): UnaryResponse;
  getApiDoc(
    requestMessage: dev_portal_api_grpc_admin_apidoc_pb.ApiDocGetRequest,
    callback: (error: ServiceError|null, responseMessage: dev_portal_api_grpc_admin_apidoc_pb.ApiDoc|null) => void
  ): UnaryResponse;
  listApiDocs(
    requestMessage: dev_portal_api_grpc_admin_apidoc_pb.ApiDocFilter,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: dev_portal_api_grpc_admin_apidoc_pb.ApiDocList|null) => void
  ): UnaryResponse;
  listApiDocs(
    requestMessage: dev_portal_api_grpc_admin_apidoc_pb.ApiDocFilter,
    callback: (error: ServiceError|null, responseMessage: dev_portal_api_grpc_admin_apidoc_pb.ApiDocList|null) => void
  ): UnaryResponse;
  createApiDoc(
    requestMessage: dev_portal_api_grpc_admin_apidoc_pb.ApiDocWriteRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: dev_portal_api_grpc_admin_apidoc_pb.ApiDoc|null) => void
  ): UnaryResponse;
  createApiDoc(
    requestMessage: dev_portal_api_grpc_admin_apidoc_pb.ApiDocWriteRequest,
    callback: (error: ServiceError|null, responseMessage: dev_portal_api_grpc_admin_apidoc_pb.ApiDoc|null) => void
  ): UnaryResponse;
  updateApiDoc(
    requestMessage: dev_portal_api_grpc_admin_apidoc_pb.ApiDocWriteRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: dev_portal_api_grpc_admin_apidoc_pb.ApiDoc|null) => void
  ): UnaryResponse;
  updateApiDoc(
    requestMessage: dev_portal_api_grpc_admin_apidoc_pb.ApiDocWriteRequest,
    callback: (error: ServiceError|null, responseMessage: dev_portal_api_grpc_admin_apidoc_pb.ApiDoc|null) => void
  ): UnaryResponse;
  deleteApiDoc(
    requestMessage: dev_portal_api_dev_portal_v1_common_pb.ObjectRef,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: google_protobuf_empty_pb.Empty|null) => void
  ): UnaryResponse;
  deleteApiDoc(
    requestMessage: dev_portal_api_dev_portal_v1_common_pb.ObjectRef,
    callback: (error: ServiceError|null, responseMessage: google_protobuf_empty_pb.Empty|null) => void
  ): UnaryResponse;
}

