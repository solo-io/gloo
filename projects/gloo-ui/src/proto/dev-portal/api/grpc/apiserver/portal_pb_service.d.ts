// package: apiserver.devportal.solo.io
// file: dev-portal/api/grpc/apiserver/portal.proto

import * as dev_portal_api_grpc_apiserver_portal_pb from "../../../../dev-portal/api/grpc/apiserver/portal_pb";
import {grpc} from "@improbable-eng/grpc-web";

type PortalApiGetPortal = {
  readonly methodName: string;
  readonly service: typeof PortalApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof dev_portal_api_grpc_apiserver_portal_pb.GetPortalRequest;
  readonly responseType: typeof dev_portal_api_grpc_apiserver_portal_pb.GetPortalResponse;
};

type PortalApiListPortals = {
  readonly methodName: string;
  readonly service: typeof PortalApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof dev_portal_api_grpc_apiserver_portal_pb.ListPortalsRequest;
  readonly responseType: typeof dev_portal_api_grpc_apiserver_portal_pb.ListPortalsResponse;
};

type PortalApiCreatePortal = {
  readonly methodName: string;
  readonly service: typeof PortalApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof dev_portal_api_grpc_apiserver_portal_pb.CreatePortalRequest;
  readonly responseType: typeof dev_portal_api_grpc_apiserver_portal_pb.CreatePortalResponse;
};

type PortalApiUpdatePortal = {
  readonly methodName: string;
  readonly service: typeof PortalApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof dev_portal_api_grpc_apiserver_portal_pb.UpdatePortalRequest;
  readonly responseType: typeof dev_portal_api_grpc_apiserver_portal_pb.UpdatePortalResponse;
};

type PortalApiDeletePortal = {
  readonly methodName: string;
  readonly service: typeof PortalApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof dev_portal_api_grpc_apiserver_portal_pb.CreatePortalRequest;
  readonly responseType: typeof dev_portal_api_grpc_apiserver_portal_pb.DeletePortalResponse;
};

export class PortalApi {
  static readonly serviceName: string;
  static readonly GetPortal: PortalApiGetPortal;
  static readonly ListPortals: PortalApiListPortals;
  static readonly CreatePortal: PortalApiCreatePortal;
  static readonly UpdatePortal: PortalApiUpdatePortal;
  static readonly DeletePortal: PortalApiDeletePortal;
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

export class PortalApiClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  getPortal(
    requestMessage: dev_portal_api_grpc_apiserver_portal_pb.GetPortalRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: dev_portal_api_grpc_apiserver_portal_pb.GetPortalResponse|null) => void
  ): UnaryResponse;
  getPortal(
    requestMessage: dev_portal_api_grpc_apiserver_portal_pb.GetPortalRequest,
    callback: (error: ServiceError|null, responseMessage: dev_portal_api_grpc_apiserver_portal_pb.GetPortalResponse|null) => void
  ): UnaryResponse;
  listPortals(
    requestMessage: dev_portal_api_grpc_apiserver_portal_pb.ListPortalsRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: dev_portal_api_grpc_apiserver_portal_pb.ListPortalsResponse|null) => void
  ): UnaryResponse;
  listPortals(
    requestMessage: dev_portal_api_grpc_apiserver_portal_pb.ListPortalsRequest,
    callback: (error: ServiceError|null, responseMessage: dev_portal_api_grpc_apiserver_portal_pb.ListPortalsResponse|null) => void
  ): UnaryResponse;
  createPortal(
    requestMessage: dev_portal_api_grpc_apiserver_portal_pb.CreatePortalRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: dev_portal_api_grpc_apiserver_portal_pb.CreatePortalResponse|null) => void
  ): UnaryResponse;
  createPortal(
    requestMessage: dev_portal_api_grpc_apiserver_portal_pb.CreatePortalRequest,
    callback: (error: ServiceError|null, responseMessage: dev_portal_api_grpc_apiserver_portal_pb.CreatePortalResponse|null) => void
  ): UnaryResponse;
  updatePortal(
    requestMessage: dev_portal_api_grpc_apiserver_portal_pb.UpdatePortalRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: dev_portal_api_grpc_apiserver_portal_pb.UpdatePortalResponse|null) => void
  ): UnaryResponse;
  updatePortal(
    requestMessage: dev_portal_api_grpc_apiserver_portal_pb.UpdatePortalRequest,
    callback: (error: ServiceError|null, responseMessage: dev_portal_api_grpc_apiserver_portal_pb.UpdatePortalResponse|null) => void
  ): UnaryResponse;
  deletePortal(
    requestMessage: dev_portal_api_grpc_apiserver_portal_pb.CreatePortalRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: dev_portal_api_grpc_apiserver_portal_pb.DeletePortalResponse|null) => void
  ): UnaryResponse;
  deletePortal(
    requestMessage: dev_portal_api_grpc_apiserver_portal_pb.CreatePortalRequest,
    callback: (error: ServiceError|null, responseMessage: dev_portal_api_grpc_apiserver_portal_pb.DeletePortalResponse|null) => void
  ): UnaryResponse;
}

