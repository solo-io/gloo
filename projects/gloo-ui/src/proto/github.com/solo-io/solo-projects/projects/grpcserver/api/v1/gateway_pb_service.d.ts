// package: glooeeapi.solo.io
// file: github.com/solo-io/solo-projects/projects/grpcserver/api/v1/gateway.proto

import * as github_com_solo_io_solo_projects_projects_grpcserver_api_v1_gateway_pb from "../../../../../../../github.com/solo-io/solo-projects/projects/grpcserver/api/v1/gateway_pb";
import {grpc} from "@improbable-eng/grpc-web";

type GatewayApiGetGateway = {
  readonly methodName: string;
  readonly service: typeof GatewayApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_gateway_pb.GetGatewayRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_gateway_pb.GetGatewayResponse;
};

type GatewayApiListGateways = {
  readonly methodName: string;
  readonly service: typeof GatewayApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_gateway_pb.ListGatewaysRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_gateway_pb.ListGatewaysResponse;
};

type GatewayApiUpdateGateway = {
  readonly methodName: string;
  readonly service: typeof GatewayApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_gateway_pb.UpdateGatewayRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_gateway_pb.UpdateGatewayResponse;
};

export class GatewayApi {
  static readonly serviceName: string;
  static readonly GetGateway: GatewayApiGetGateway;
  static readonly ListGateways: GatewayApiListGateways;
  static readonly UpdateGateway: GatewayApiUpdateGateway;
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

export class GatewayApiClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  getGateway(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_gateway_pb.GetGatewayRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_gateway_pb.GetGatewayResponse|null) => void
  ): UnaryResponse;
  getGateway(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_gateway_pb.GetGatewayRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_gateway_pb.GetGatewayResponse|null) => void
  ): UnaryResponse;
  listGateways(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_gateway_pb.ListGatewaysRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_gateway_pb.ListGatewaysResponse|null) => void
  ): UnaryResponse;
  listGateways(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_gateway_pb.ListGatewaysRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_gateway_pb.ListGatewaysResponse|null) => void
  ): UnaryResponse;
  updateGateway(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_gateway_pb.UpdateGatewayRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_gateway_pb.UpdateGatewayResponse|null) => void
  ): UnaryResponse;
  updateGateway(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_gateway_pb.UpdateGatewayRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_gateway_pb.UpdateGatewayResponse|null) => void
  ): UnaryResponse;
}

