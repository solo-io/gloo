// package: fed.rpc.solo.io
// file: github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/federated_gateway_resources.proto

import * as github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gateway_resources_pb from "../../../../../../../../github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/federated_gateway_resources_pb";
import {grpc} from "@improbable-eng/grpc-web";

type FederatedGatewayResourceApiListFederatedGateways = {
  readonly methodName: string;
  readonly service: typeof FederatedGatewayResourceApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gateway_resources_pb.ListFederatedGatewaysRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gateway_resources_pb.ListFederatedGatewaysResponse;
};

type FederatedGatewayResourceApiGetFederatedGatewayYaml = {
  readonly methodName: string;
  readonly service: typeof FederatedGatewayResourceApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gateway_resources_pb.GetFederatedGatewayYamlRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gateway_resources_pb.GetFederatedGatewayYamlResponse;
};

type FederatedGatewayResourceApiListFederatedVirtualServices = {
  readonly methodName: string;
  readonly service: typeof FederatedGatewayResourceApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gateway_resources_pb.ListFederatedVirtualServicesRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gateway_resources_pb.ListFederatedVirtualServicesResponse;
};

type FederatedGatewayResourceApiGetFederatedVirtualServiceYaml = {
  readonly methodName: string;
  readonly service: typeof FederatedGatewayResourceApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gateway_resources_pb.GetFederatedVirtualServiceYamlRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gateway_resources_pb.GetFederatedVirtualServiceYamlResponse;
};

type FederatedGatewayResourceApiListFederatedRouteTables = {
  readonly methodName: string;
  readonly service: typeof FederatedGatewayResourceApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gateway_resources_pb.ListFederatedRouteTablesRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gateway_resources_pb.ListFederatedRouteTablesResponse;
};

type FederatedGatewayResourceApiGetFederatedRouteTableYaml = {
  readonly methodName: string;
  readonly service: typeof FederatedGatewayResourceApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gateway_resources_pb.GetFederatedRouteTableYamlRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gateway_resources_pb.GetFederatedRouteTableYamlResponse;
};

export class FederatedGatewayResourceApi {
  static readonly serviceName: string;
  static readonly ListFederatedGateways: FederatedGatewayResourceApiListFederatedGateways;
  static readonly GetFederatedGatewayYaml: FederatedGatewayResourceApiGetFederatedGatewayYaml;
  static readonly ListFederatedVirtualServices: FederatedGatewayResourceApiListFederatedVirtualServices;
  static readonly GetFederatedVirtualServiceYaml: FederatedGatewayResourceApiGetFederatedVirtualServiceYaml;
  static readonly ListFederatedRouteTables: FederatedGatewayResourceApiListFederatedRouteTables;
  static readonly GetFederatedRouteTableYaml: FederatedGatewayResourceApiGetFederatedRouteTableYaml;
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

export class FederatedGatewayResourceApiClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  listFederatedGateways(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gateway_resources_pb.ListFederatedGatewaysRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gateway_resources_pb.ListFederatedGatewaysResponse|null) => void
  ): UnaryResponse;
  listFederatedGateways(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gateway_resources_pb.ListFederatedGatewaysRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gateway_resources_pb.ListFederatedGatewaysResponse|null) => void
  ): UnaryResponse;
  getFederatedGatewayYaml(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gateway_resources_pb.GetFederatedGatewayYamlRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gateway_resources_pb.GetFederatedGatewayYamlResponse|null) => void
  ): UnaryResponse;
  getFederatedGatewayYaml(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gateway_resources_pb.GetFederatedGatewayYamlRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gateway_resources_pb.GetFederatedGatewayYamlResponse|null) => void
  ): UnaryResponse;
  listFederatedVirtualServices(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gateway_resources_pb.ListFederatedVirtualServicesRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gateway_resources_pb.ListFederatedVirtualServicesResponse|null) => void
  ): UnaryResponse;
  listFederatedVirtualServices(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gateway_resources_pb.ListFederatedVirtualServicesRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gateway_resources_pb.ListFederatedVirtualServicesResponse|null) => void
  ): UnaryResponse;
  getFederatedVirtualServiceYaml(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gateway_resources_pb.GetFederatedVirtualServiceYamlRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gateway_resources_pb.GetFederatedVirtualServiceYamlResponse|null) => void
  ): UnaryResponse;
  getFederatedVirtualServiceYaml(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gateway_resources_pb.GetFederatedVirtualServiceYamlRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gateway_resources_pb.GetFederatedVirtualServiceYamlResponse|null) => void
  ): UnaryResponse;
  listFederatedRouteTables(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gateway_resources_pb.ListFederatedRouteTablesRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gateway_resources_pb.ListFederatedRouteTablesResponse|null) => void
  ): UnaryResponse;
  listFederatedRouteTables(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gateway_resources_pb.ListFederatedRouteTablesRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gateway_resources_pb.ListFederatedRouteTablesResponse|null) => void
  ): UnaryResponse;
  getFederatedRouteTableYaml(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gateway_resources_pb.GetFederatedRouteTableYamlRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gateway_resources_pb.GetFederatedRouteTableYamlResponse|null) => void
  ): UnaryResponse;
  getFederatedRouteTableYaml(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gateway_resources_pb.GetFederatedRouteTableYamlRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gateway_resources_pb.GetFederatedRouteTableYamlResponse|null) => void
  ): UnaryResponse;
}

