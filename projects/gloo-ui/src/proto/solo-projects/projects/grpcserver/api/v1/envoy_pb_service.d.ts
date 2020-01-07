// package: glooeeapi.solo.io
// file: solo-projects/projects/grpcserver/api/v1/envoy.proto

import * as solo_projects_projects_grpcserver_api_v1_envoy_pb from "../../../../../solo-projects/projects/grpcserver/api/v1/envoy_pb";
import {grpc} from "@improbable-eng/grpc-web";

type EnvoyApiListEnvoyDetails = {
  readonly methodName: string;
  readonly service: typeof EnvoyApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof solo_projects_projects_grpcserver_api_v1_envoy_pb.ListEnvoyDetailsRequest;
  readonly responseType: typeof solo_projects_projects_grpcserver_api_v1_envoy_pb.ListEnvoyDetailsResponse;
};

export class EnvoyApi {
  static readonly serviceName: string;
  static readonly ListEnvoyDetails: EnvoyApiListEnvoyDetails;
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

export class EnvoyApiClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  listEnvoyDetails(
    requestMessage: solo_projects_projects_grpcserver_api_v1_envoy_pb.ListEnvoyDetailsRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: solo_projects_projects_grpcserver_api_v1_envoy_pb.ListEnvoyDetailsResponse|null) => void
  ): UnaryResponse;
  listEnvoyDetails(
    requestMessage: solo_projects_projects_grpcserver_api_v1_envoy_pb.ListEnvoyDetailsRequest,
    callback: (error: ServiceError|null, responseMessage: solo_projects_projects_grpcserver_api_v1_envoy_pb.ListEnvoyDetailsResponse|null) => void
  ): UnaryResponse;
}

