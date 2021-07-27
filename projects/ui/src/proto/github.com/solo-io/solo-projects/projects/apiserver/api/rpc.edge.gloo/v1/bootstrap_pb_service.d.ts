// package: rpc.edge.gloo.solo.io
// file: github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/bootstrap.proto

import * as github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_bootstrap_pb from "../../../../../../../../github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/bootstrap_pb";
import {grpc} from "@improbable-eng/grpc-web";

type BootstrapApiIsGlooFedEnabled = {
  readonly methodName: string;
  readonly service: typeof BootstrapApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_bootstrap_pb.GlooFedCheckRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_bootstrap_pb.GlooFedCheckResponse;
};

export class BootstrapApi {
  static readonly serviceName: string;
  static readonly IsGlooFedEnabled: BootstrapApiIsGlooFedEnabled;
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

export class BootstrapApiClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  isGlooFedEnabled(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_bootstrap_pb.GlooFedCheckRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_bootstrap_pb.GlooFedCheckResponse|null) => void
  ): UnaryResponse;
  isGlooFedEnabled(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_bootstrap_pb.GlooFedCheckRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_bootstrap_pb.GlooFedCheckResponse|null) => void
  ): UnaryResponse;
}

