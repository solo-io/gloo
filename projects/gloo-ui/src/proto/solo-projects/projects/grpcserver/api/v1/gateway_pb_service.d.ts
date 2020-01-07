// package: glooeeapi.solo.io
// file: solo-projects/projects/grpcserver/api/v1/gateway.proto

import * as solo_projects_projects_grpcserver_api_v1_gateway_pb from "../../../../../solo-projects/projects/grpcserver/api/v1/gateway_pb";
import {grpc} from "@improbable-eng/grpc-web";

type GatewayApiGetGateway = {
  readonly methodName: string;
  readonly service: typeof GatewayApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof solo_projects_projects_grpcserver_api_v1_gateway_pb.GetGatewayRequest;
  readonly responseType: typeof solo_projects_projects_grpcserver_api_v1_gateway_pb.GetGatewayResponse;
};

type GatewayApiListGateways = {
  readonly methodName: string;
  readonly service: typeof GatewayApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof solo_projects_projects_grpcserver_api_v1_gateway_pb.ListGatewaysRequest;
  readonly responseType: typeof solo_projects_projects_grpcserver_api_v1_gateway_pb.ListGatewaysResponse;
};

type GatewayApiUpdateGateway = {
  readonly methodName: string;
  readonly service: typeof GatewayApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof solo_projects_projects_grpcserver_api_v1_gateway_pb.UpdateGatewayRequest;
  readonly responseType: typeof solo_projects_projects_grpcserver_api_v1_gateway_pb.UpdateGatewayResponse;
};

type GatewayApiUpdateGatewayYaml = {
  readonly methodName: string;
  readonly service: typeof GatewayApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof solo_projects_projects_grpcserver_api_v1_gateway_pb.UpdateGatewayYamlRequest;
  readonly responseType: typeof solo_projects_projects_grpcserver_api_v1_gateway_pb.UpdateGatewayResponse;
};

export class GatewayApi {
  static readonly serviceName: string;
  static readonly GetGateway: GatewayApiGetGateway;
  static readonly ListGateways: GatewayApiListGateways;
  static readonly UpdateGateway: GatewayApiUpdateGateway;
  static readonly UpdateGatewayYaml: GatewayApiUpdateGatewayYaml;
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
    requestMessage: solo_projects_projects_grpcserver_api_v1_gateway_pb.GetGatewayRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: solo_projects_projects_grpcserver_api_v1_gateway_pb.GetGatewayResponse|null) => void
  ): UnaryResponse;
  getGateway(
    requestMessage: solo_projects_projects_grpcserver_api_v1_gateway_pb.GetGatewayRequest,
    callback: (error: ServiceError|null, responseMessage: solo_projects_projects_grpcserver_api_v1_gateway_pb.GetGatewayResponse|null) => void
  ): UnaryResponse;
  listGateways(
    requestMessage: solo_projects_projects_grpcserver_api_v1_gateway_pb.ListGatewaysRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: solo_projects_projects_grpcserver_api_v1_gateway_pb.ListGatewaysResponse|null) => void
  ): UnaryResponse;
  listGateways(
    requestMessage: solo_projects_projects_grpcserver_api_v1_gateway_pb.ListGatewaysRequest,
    callback: (error: ServiceError|null, responseMessage: solo_projects_projects_grpcserver_api_v1_gateway_pb.ListGatewaysResponse|null) => void
  ): UnaryResponse;
  updateGateway(
    requestMessage: solo_projects_projects_grpcserver_api_v1_gateway_pb.UpdateGatewayRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: solo_projects_projects_grpcserver_api_v1_gateway_pb.UpdateGatewayResponse|null) => void
  ): UnaryResponse;
  updateGateway(
    requestMessage: solo_projects_projects_grpcserver_api_v1_gateway_pb.UpdateGatewayRequest,
    callback: (error: ServiceError|null, responseMessage: solo_projects_projects_grpcserver_api_v1_gateway_pb.UpdateGatewayResponse|null) => void
  ): UnaryResponse;
  updateGatewayYaml(
    requestMessage: solo_projects_projects_grpcserver_api_v1_gateway_pb.UpdateGatewayYamlRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: solo_projects_projects_grpcserver_api_v1_gateway_pb.UpdateGatewayResponse|null) => void
  ): UnaryResponse;
  updateGatewayYaml(
    requestMessage: solo_projects_projects_grpcserver_api_v1_gateway_pb.UpdateGatewayYamlRequest,
    callback: (error: ServiceError|null, responseMessage: solo_projects_projects_grpcserver_api_v1_gateway_pb.UpdateGatewayResponse|null) => void
  ): UnaryResponse;
}

