// package: fed.rpc.solo.io
// file: github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/failover_scheme.proto

import * as github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_failover_scheme_pb from "../../../../../../../../github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/failover_scheme_pb";
import {grpc} from "@improbable-eng/grpc-web";

type FailoverSchemeApiGetFailoverScheme = {
  readonly methodName: string;
  readonly service: typeof FailoverSchemeApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_failover_scheme_pb.GetFailoverSchemeRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_failover_scheme_pb.GetFailoverSchemeResponse;
};

type FailoverSchemeApiGetFailoverSchemeYaml = {
  readonly methodName: string;
  readonly service: typeof FailoverSchemeApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_failover_scheme_pb.GetFailoverSchemeYamlRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_failover_scheme_pb.GetFailoverSchemeYamlResponse;
};

export class FailoverSchemeApi {
  static readonly serviceName: string;
  static readonly GetFailoverScheme: FailoverSchemeApiGetFailoverScheme;
  static readonly GetFailoverSchemeYaml: FailoverSchemeApiGetFailoverSchemeYaml;
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

export class FailoverSchemeApiClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  getFailoverScheme(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_failover_scheme_pb.GetFailoverSchemeRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_failover_scheme_pb.GetFailoverSchemeResponse|null) => void
  ): UnaryResponse;
  getFailoverScheme(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_failover_scheme_pb.GetFailoverSchemeRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_failover_scheme_pb.GetFailoverSchemeResponse|null) => void
  ): UnaryResponse;
  getFailoverSchemeYaml(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_failover_scheme_pb.GetFailoverSchemeYamlRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_failover_scheme_pb.GetFailoverSchemeYamlResponse|null) => void
  ): UnaryResponse;
  getFailoverSchemeYaml(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_failover_scheme_pb.GetFailoverSchemeYamlRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_failover_scheme_pb.GetFailoverSchemeYamlResponse|null) => void
  ): UnaryResponse;
}

