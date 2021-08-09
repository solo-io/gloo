// package: rpc.edge.gloo.solo.io
// file: github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/gateway_resources.proto

import * as github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gateway_resources_pb from "../../../../../../../../github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/gateway_resources_pb";
import {grpc} from "@improbable-eng/grpc-web";

type GatewayResourceApiListGateways = {
  readonly methodName: string;
  readonly service: typeof GatewayResourceApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gateway_resources_pb.ListGatewaysRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gateway_resources_pb.ListGatewaysResponse;
};

type GatewayResourceApiGetGatewayYaml = {
  readonly methodName: string;
  readonly service: typeof GatewayResourceApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gateway_resources_pb.GetGatewayYamlRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gateway_resources_pb.GetGatewayYamlResponse;
};

type GatewayResourceApiListVirtualServices = {
  readonly methodName: string;
  readonly service: typeof GatewayResourceApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gateway_resources_pb.ListVirtualServicesRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gateway_resources_pb.ListVirtualServicesResponse;
};

type GatewayResourceApiGetVirtualServiceYaml = {
  readonly methodName: string;
  readonly service: typeof GatewayResourceApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gateway_resources_pb.GetVirtualServiceYamlRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gateway_resources_pb.GetVirtualServiceYamlResponse;
};

type GatewayResourceApiListRouteTables = {
  readonly methodName: string;
  readonly service: typeof GatewayResourceApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gateway_resources_pb.ListRouteTablesRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gateway_resources_pb.ListRouteTablesResponse;
};

type GatewayResourceApiGetRouteTableYaml = {
  readonly methodName: string;
  readonly service: typeof GatewayResourceApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gateway_resources_pb.GetRouteTableYamlRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gateway_resources_pb.GetRouteTableYamlResponse;
};

export class GatewayResourceApi {
  static readonly serviceName: string;
  static readonly ListGateways: GatewayResourceApiListGateways;
  static readonly GetGatewayYaml: GatewayResourceApiGetGatewayYaml;
  static readonly ListVirtualServices: GatewayResourceApiListVirtualServices;
  static readonly GetVirtualServiceYaml: GatewayResourceApiGetVirtualServiceYaml;
  static readonly ListRouteTables: GatewayResourceApiListRouteTables;
  static readonly GetRouteTableYaml: GatewayResourceApiGetRouteTableYaml;
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

export class GatewayResourceApiClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  listGateways(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gateway_resources_pb.ListGatewaysRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gateway_resources_pb.ListGatewaysResponse|null) => void
  ): UnaryResponse;
  listGateways(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gateway_resources_pb.ListGatewaysRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gateway_resources_pb.ListGatewaysResponse|null) => void
  ): UnaryResponse;
  getGatewayYaml(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gateway_resources_pb.GetGatewayYamlRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gateway_resources_pb.GetGatewayYamlResponse|null) => void
  ): UnaryResponse;
  getGatewayYaml(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gateway_resources_pb.GetGatewayYamlRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gateway_resources_pb.GetGatewayYamlResponse|null) => void
  ): UnaryResponse;
  listVirtualServices(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gateway_resources_pb.ListVirtualServicesRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gateway_resources_pb.ListVirtualServicesResponse|null) => void
  ): UnaryResponse;
  listVirtualServices(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gateway_resources_pb.ListVirtualServicesRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gateway_resources_pb.ListVirtualServicesResponse|null) => void
  ): UnaryResponse;
  getVirtualServiceYaml(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gateway_resources_pb.GetVirtualServiceYamlRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gateway_resources_pb.GetVirtualServiceYamlResponse|null) => void
  ): UnaryResponse;
  getVirtualServiceYaml(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gateway_resources_pb.GetVirtualServiceYamlRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gateway_resources_pb.GetVirtualServiceYamlResponse|null) => void
  ): UnaryResponse;
  listRouteTables(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gateway_resources_pb.ListRouteTablesRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gateway_resources_pb.ListRouteTablesResponse|null) => void
  ): UnaryResponse;
  listRouteTables(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gateway_resources_pb.ListRouteTablesRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gateway_resources_pb.ListRouteTablesResponse|null) => void
  ): UnaryResponse;
  getRouteTableYaml(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gateway_resources_pb.GetRouteTableYamlRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gateway_resources_pb.GetRouteTableYamlResponse|null) => void
  ): UnaryResponse;
  getRouteTableYaml(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gateway_resources_pb.GetRouteTableYamlRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_gateway_resources_pb.GetRouteTableYamlResponse|null) => void
  ): UnaryResponse;
}

