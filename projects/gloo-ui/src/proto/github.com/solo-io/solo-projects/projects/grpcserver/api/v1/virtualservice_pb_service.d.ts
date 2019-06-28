// package: glooeeapi.solo.io
// file: github.com/solo-io/solo-projects/projects/grpcserver/api/v1/virtualservice.proto

import * as github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb from "../../../../../../../github.com/solo-io/solo-projects/projects/grpcserver/api/v1/virtualservice_pb";
import {grpc} from "@improbable-eng/grpc-web";

type VirtualServiceApiGetVirtualService = {
  readonly methodName: string;
  readonly service: typeof VirtualServiceApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.GetVirtualServiceRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.GetVirtualServiceResponse;
};

type VirtualServiceApiListVirtualServices = {
  readonly methodName: string;
  readonly service: typeof VirtualServiceApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.ListVirtualServicesRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.ListVirtualServicesResponse;
};

type VirtualServiceApiStreamVirtualServiceList = {
  readonly methodName: string;
  readonly service: typeof VirtualServiceApi;
  readonly requestStream: false;
  readonly responseStream: true;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.StreamVirtualServiceListRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.StreamVirtualServiceListResponse;
};

type VirtualServiceApiCreateVirtualService = {
  readonly methodName: string;
  readonly service: typeof VirtualServiceApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.CreateVirtualServiceRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.CreateVirtualServiceResponse;
};

type VirtualServiceApiUpdateVirtualService = {
  readonly methodName: string;
  readonly service: typeof VirtualServiceApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.UpdateVirtualServiceRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.UpdateVirtualServiceResponse;
};

type VirtualServiceApiDeleteVirtualService = {
  readonly methodName: string;
  readonly service: typeof VirtualServiceApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.DeleteVirtualServiceRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.DeleteVirtualServiceResponse;
};

type VirtualServiceApiCreateRoute = {
  readonly methodName: string;
  readonly service: typeof VirtualServiceApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.CreateRouteRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.CreateRouteResponse;
};

type VirtualServiceApiUpdateRoute = {
  readonly methodName: string;
  readonly service: typeof VirtualServiceApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.UpdateRouteRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.UpdateRouteResponse;
};

type VirtualServiceApiDeleteRoute = {
  readonly methodName: string;
  readonly service: typeof VirtualServiceApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.DeleteRouteRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.DeleteRouteResponse;
};

type VirtualServiceApiSwapRoutes = {
  readonly methodName: string;
  readonly service: typeof VirtualServiceApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.SwapRoutesRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.SwapRoutesResponse;
};

type VirtualServiceApiShiftRoutes = {
  readonly methodName: string;
  readonly service: typeof VirtualServiceApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.ShiftRoutesRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.ShiftRoutesResponse;
};

export class VirtualServiceApi {
  static readonly serviceName: string;
  static readonly GetVirtualService: VirtualServiceApiGetVirtualService;
  static readonly ListVirtualServices: VirtualServiceApiListVirtualServices;
  static readonly StreamVirtualServiceList: VirtualServiceApiStreamVirtualServiceList;
  static readonly CreateVirtualService: VirtualServiceApiCreateVirtualService;
  static readonly UpdateVirtualService: VirtualServiceApiUpdateVirtualService;
  static readonly DeleteVirtualService: VirtualServiceApiDeleteVirtualService;
  static readonly CreateRoute: VirtualServiceApiCreateRoute;
  static readonly UpdateRoute: VirtualServiceApiUpdateRoute;
  static readonly DeleteRoute: VirtualServiceApiDeleteRoute;
  static readonly SwapRoutes: VirtualServiceApiSwapRoutes;
  static readonly ShiftRoutes: VirtualServiceApiShiftRoutes;
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

export class VirtualServiceApiClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  getVirtualService(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.GetVirtualServiceRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.GetVirtualServiceResponse|null) => void
  ): UnaryResponse;
  getVirtualService(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.GetVirtualServiceRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.GetVirtualServiceResponse|null) => void
  ): UnaryResponse;
  listVirtualServices(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.ListVirtualServicesRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.ListVirtualServicesResponse|null) => void
  ): UnaryResponse;
  listVirtualServices(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.ListVirtualServicesRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.ListVirtualServicesResponse|null) => void
  ): UnaryResponse;
  streamVirtualServiceList(requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.StreamVirtualServiceListRequest, metadata?: grpc.Metadata): ResponseStream<github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.StreamVirtualServiceListResponse>;
  createVirtualService(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.CreateVirtualServiceRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.CreateVirtualServiceResponse|null) => void
  ): UnaryResponse;
  createVirtualService(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.CreateVirtualServiceRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.CreateVirtualServiceResponse|null) => void
  ): UnaryResponse;
  updateVirtualService(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.UpdateVirtualServiceRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.UpdateVirtualServiceResponse|null) => void
  ): UnaryResponse;
  updateVirtualService(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.UpdateVirtualServiceRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.UpdateVirtualServiceResponse|null) => void
  ): UnaryResponse;
  deleteVirtualService(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.DeleteVirtualServiceRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.DeleteVirtualServiceResponse|null) => void
  ): UnaryResponse;
  deleteVirtualService(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.DeleteVirtualServiceRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.DeleteVirtualServiceResponse|null) => void
  ): UnaryResponse;
  createRoute(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.CreateRouteRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.CreateRouteResponse|null) => void
  ): UnaryResponse;
  createRoute(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.CreateRouteRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.CreateRouteResponse|null) => void
  ): UnaryResponse;
  updateRoute(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.UpdateRouteRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.UpdateRouteResponse|null) => void
  ): UnaryResponse;
  updateRoute(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.UpdateRouteRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.UpdateRouteResponse|null) => void
  ): UnaryResponse;
  deleteRoute(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.DeleteRouteRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.DeleteRouteResponse|null) => void
  ): UnaryResponse;
  deleteRoute(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.DeleteRouteRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.DeleteRouteResponse|null) => void
  ): UnaryResponse;
  swapRoutes(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.SwapRoutesRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.SwapRoutesResponse|null) => void
  ): UnaryResponse;
  swapRoutes(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.SwapRoutesRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.SwapRoutesResponse|null) => void
  ): UnaryResponse;
  shiftRoutes(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.ShiftRoutesRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.ShiftRoutesResponse|null) => void
  ): UnaryResponse;
  shiftRoutes(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.ShiftRoutesRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.ShiftRoutesResponse|null) => void
  ): UnaryResponse;
}

