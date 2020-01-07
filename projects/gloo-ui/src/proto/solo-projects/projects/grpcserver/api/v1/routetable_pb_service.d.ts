// package: glooeeapi.solo.io
// file: solo-projects/projects/grpcserver/api/v1/routetable.proto

import * as solo_projects_projects_grpcserver_api_v1_routetable_pb from "../../../../../solo-projects/projects/grpcserver/api/v1/routetable_pb";
import {grpc} from "@improbable-eng/grpc-web";

type RouteTableApiGetRouteTable = {
  readonly methodName: string;
  readonly service: typeof RouteTableApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof solo_projects_projects_grpcserver_api_v1_routetable_pb.GetRouteTableRequest;
  readonly responseType: typeof solo_projects_projects_grpcserver_api_v1_routetable_pb.GetRouteTableResponse;
};

type RouteTableApiListRouteTables = {
  readonly methodName: string;
  readonly service: typeof RouteTableApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof solo_projects_projects_grpcserver_api_v1_routetable_pb.ListRouteTablesRequest;
  readonly responseType: typeof solo_projects_projects_grpcserver_api_v1_routetable_pb.ListRouteTablesResponse;
};

type RouteTableApiCreateRouteTable = {
  readonly methodName: string;
  readonly service: typeof RouteTableApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof solo_projects_projects_grpcserver_api_v1_routetable_pb.CreateRouteTableRequest;
  readonly responseType: typeof solo_projects_projects_grpcserver_api_v1_routetable_pb.CreateRouteTableResponse;
};

type RouteTableApiUpdateRouteTable = {
  readonly methodName: string;
  readonly service: typeof RouteTableApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof solo_projects_projects_grpcserver_api_v1_routetable_pb.UpdateRouteTableRequest;
  readonly responseType: typeof solo_projects_projects_grpcserver_api_v1_routetable_pb.UpdateRouteTableResponse;
};

type RouteTableApiUpdateRouteTableYaml = {
  readonly methodName: string;
  readonly service: typeof RouteTableApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof solo_projects_projects_grpcserver_api_v1_routetable_pb.UpdateRouteTableYamlRequest;
  readonly responseType: typeof solo_projects_projects_grpcserver_api_v1_routetable_pb.UpdateRouteTableResponse;
};

type RouteTableApiDeleteRouteTable = {
  readonly methodName: string;
  readonly service: typeof RouteTableApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof solo_projects_projects_grpcserver_api_v1_routetable_pb.DeleteRouteTableRequest;
  readonly responseType: typeof solo_projects_projects_grpcserver_api_v1_routetable_pb.DeleteRouteTableResponse;
};

export class RouteTableApi {
  static readonly serviceName: string;
  static readonly GetRouteTable: RouteTableApiGetRouteTable;
  static readonly ListRouteTables: RouteTableApiListRouteTables;
  static readonly CreateRouteTable: RouteTableApiCreateRouteTable;
  static readonly UpdateRouteTable: RouteTableApiUpdateRouteTable;
  static readonly UpdateRouteTableYaml: RouteTableApiUpdateRouteTableYaml;
  static readonly DeleteRouteTable: RouteTableApiDeleteRouteTable;
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

export class RouteTableApiClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  getRouteTable(
    requestMessage: solo_projects_projects_grpcserver_api_v1_routetable_pb.GetRouteTableRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: solo_projects_projects_grpcserver_api_v1_routetable_pb.GetRouteTableResponse|null) => void
  ): UnaryResponse;
  getRouteTable(
    requestMessage: solo_projects_projects_grpcserver_api_v1_routetable_pb.GetRouteTableRequest,
    callback: (error: ServiceError|null, responseMessage: solo_projects_projects_grpcserver_api_v1_routetable_pb.GetRouteTableResponse|null) => void
  ): UnaryResponse;
  listRouteTables(
    requestMessage: solo_projects_projects_grpcserver_api_v1_routetable_pb.ListRouteTablesRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: solo_projects_projects_grpcserver_api_v1_routetable_pb.ListRouteTablesResponse|null) => void
  ): UnaryResponse;
  listRouteTables(
    requestMessage: solo_projects_projects_grpcserver_api_v1_routetable_pb.ListRouteTablesRequest,
    callback: (error: ServiceError|null, responseMessage: solo_projects_projects_grpcserver_api_v1_routetable_pb.ListRouteTablesResponse|null) => void
  ): UnaryResponse;
  createRouteTable(
    requestMessage: solo_projects_projects_grpcserver_api_v1_routetable_pb.CreateRouteTableRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: solo_projects_projects_grpcserver_api_v1_routetable_pb.CreateRouteTableResponse|null) => void
  ): UnaryResponse;
  createRouteTable(
    requestMessage: solo_projects_projects_grpcserver_api_v1_routetable_pb.CreateRouteTableRequest,
    callback: (error: ServiceError|null, responseMessage: solo_projects_projects_grpcserver_api_v1_routetable_pb.CreateRouteTableResponse|null) => void
  ): UnaryResponse;
  updateRouteTable(
    requestMessage: solo_projects_projects_grpcserver_api_v1_routetable_pb.UpdateRouteTableRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: solo_projects_projects_grpcserver_api_v1_routetable_pb.UpdateRouteTableResponse|null) => void
  ): UnaryResponse;
  updateRouteTable(
    requestMessage: solo_projects_projects_grpcserver_api_v1_routetable_pb.UpdateRouteTableRequest,
    callback: (error: ServiceError|null, responseMessage: solo_projects_projects_grpcserver_api_v1_routetable_pb.UpdateRouteTableResponse|null) => void
  ): UnaryResponse;
  updateRouteTableYaml(
    requestMessage: solo_projects_projects_grpcserver_api_v1_routetable_pb.UpdateRouteTableYamlRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: solo_projects_projects_grpcserver_api_v1_routetable_pb.UpdateRouteTableResponse|null) => void
  ): UnaryResponse;
  updateRouteTableYaml(
    requestMessage: solo_projects_projects_grpcserver_api_v1_routetable_pb.UpdateRouteTableYamlRequest,
    callback: (error: ServiceError|null, responseMessage: solo_projects_projects_grpcserver_api_v1_routetable_pb.UpdateRouteTableResponse|null) => void
  ): UnaryResponse;
  deleteRouteTable(
    requestMessage: solo_projects_projects_grpcserver_api_v1_routetable_pb.DeleteRouteTableRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: solo_projects_projects_grpcserver_api_v1_routetable_pb.DeleteRouteTableResponse|null) => void
  ): UnaryResponse;
  deleteRouteTable(
    requestMessage: solo_projects_projects_grpcserver_api_v1_routetable_pb.DeleteRouteTableRequest,
    callback: (error: ServiceError|null, responseMessage: solo_projects_projects_grpcserver_api_v1_routetable_pb.DeleteRouteTableResponse|null) => void
  ): UnaryResponse;
}

