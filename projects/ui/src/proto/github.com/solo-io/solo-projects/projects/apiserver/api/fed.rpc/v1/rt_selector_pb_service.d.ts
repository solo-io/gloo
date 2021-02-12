// package: fed.rpc.solo.io
// file: github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/rt_selector.proto

import * as github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_rt_selector_pb from "../../../../../../../../github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/rt_selector_pb";
import {grpc} from "@improbable-eng/grpc-web";

type VirtualServiceRoutesApiGetVirtualServiceRoutes = {
  readonly methodName: string;
  readonly service: typeof VirtualServiceRoutesApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_rt_selector_pb.GetVirtualServiceRoutesRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_rt_selector_pb.GetVirtualServiceRoutesResponse;
};

export class VirtualServiceRoutesApi {
  static readonly serviceName: string;
  static readonly GetVirtualServiceRoutes: VirtualServiceRoutesApiGetVirtualServiceRoutes;
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

export class VirtualServiceRoutesApiClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  getVirtualServiceRoutes(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_rt_selector_pb.GetVirtualServiceRoutesRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_rt_selector_pb.GetVirtualServiceRoutesResponse|null) => void
  ): UnaryResponse;
  getVirtualServiceRoutes(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_rt_selector_pb.GetVirtualServiceRoutesRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_rt_selector_pb.GetVirtualServiceRoutesResponse|null) => void
  ): UnaryResponse;
}

