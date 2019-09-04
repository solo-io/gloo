// package: extauth.plugins.gloo.solo.io
// file: github.com/solo-io/gloo/projects/gloo/api/v1/enterprise/plugins/extauth/extauth.proto

import * as github_com_solo_io_gloo_projects_gloo_api_v1_enterprise_plugins_extauth_extauth_pb from "../../../../../../../../../../github.com/solo-io/gloo/projects/gloo/api/v1/enterprise/plugins/extauth/extauth_pb";
import * as envoy_api_v2_discovery_pb from "../../../../../../../../../../envoy/api/v2/discovery_pb";
import {grpc} from "@improbable-eng/grpc-web";

type ExtAuthDiscoveryServiceStreamExtAuthConfig = {
  readonly methodName: string;
  readonly service: typeof ExtAuthDiscoveryService;
  readonly requestStream: true;
  readonly responseStream: true;
  readonly requestType: typeof envoy_api_v2_discovery_pb.DiscoveryRequest;
  readonly responseType: typeof envoy_api_v2_discovery_pb.DiscoveryResponse;
};

type ExtAuthDiscoveryServiceDeltaExtAuthConfig = {
  readonly methodName: string;
  readonly service: typeof ExtAuthDiscoveryService;
  readonly requestStream: true;
  readonly responseStream: true;
  readonly requestType: typeof envoy_api_v2_discovery_pb.DeltaDiscoveryRequest;
  readonly responseType: typeof envoy_api_v2_discovery_pb.DeltaDiscoveryResponse;
};

type ExtAuthDiscoveryServiceFetchExtAuthConfig = {
  readonly methodName: string;
  readonly service: typeof ExtAuthDiscoveryService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof envoy_api_v2_discovery_pb.DiscoveryRequest;
  readonly responseType: typeof envoy_api_v2_discovery_pb.DiscoveryResponse;
};

export class ExtAuthDiscoveryService {
  static readonly serviceName: string;
  static readonly StreamExtAuthConfig: ExtAuthDiscoveryServiceStreamExtAuthConfig;
  static readonly DeltaExtAuthConfig: ExtAuthDiscoveryServiceDeltaExtAuthConfig;
  static readonly FetchExtAuthConfig: ExtAuthDiscoveryServiceFetchExtAuthConfig;
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

export class ExtAuthDiscoveryServiceClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  streamExtAuthConfig(metadata?: grpc.Metadata): BidirectionalStream<envoy_api_v2_discovery_pb.DiscoveryRequest, envoy_api_v2_discovery_pb.DiscoveryResponse>;
  deltaExtAuthConfig(metadata?: grpc.Metadata): BidirectionalStream<envoy_api_v2_discovery_pb.DeltaDiscoveryRequest, envoy_api_v2_discovery_pb.DeltaDiscoveryResponse>;
  fetchExtAuthConfig(
    requestMessage: envoy_api_v2_discovery_pb.DiscoveryRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: envoy_api_v2_discovery_pb.DiscoveryResponse|null) => void
  ): UnaryResponse;
  fetchExtAuthConfig(
    requestMessage: envoy_api_v2_discovery_pb.DiscoveryRequest,
    callback: (error: ServiceError|null, responseMessage: envoy_api_v2_discovery_pb.DiscoveryResponse|null) => void
  ): UnaryResponse;
}

