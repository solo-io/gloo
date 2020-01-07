// package: glooe.solo.io
// file: gloo/projects/gloo/api/v1/enterprise/ratelimit.proto

import * as gloo_projects_gloo_api_v1_enterprise_ratelimit_pb from "../../../../../../gloo/projects/gloo/api/v1/enterprise/ratelimit_pb";
import * as envoy_api_v2_discovery_pb from "../../../../../../envoy/api/v2/discovery_pb";
import {grpc} from "@improbable-eng/grpc-web";

type RateLimitDiscoveryServiceStreamRateLimitConfig = {
  readonly methodName: string;
  readonly service: typeof RateLimitDiscoveryService;
  readonly requestStream: true;
  readonly responseStream: true;
  readonly requestType: typeof envoy_api_v2_discovery_pb.DiscoveryRequest;
  readonly responseType: typeof envoy_api_v2_discovery_pb.DiscoveryResponse;
};

type RateLimitDiscoveryServiceDeltaRateLimitConfig = {
  readonly methodName: string;
  readonly service: typeof RateLimitDiscoveryService;
  readonly requestStream: true;
  readonly responseStream: true;
  readonly requestType: typeof envoy_api_v2_discovery_pb.DeltaDiscoveryRequest;
  readonly responseType: typeof envoy_api_v2_discovery_pb.DeltaDiscoveryResponse;
};

type RateLimitDiscoveryServiceFetchRateLimitConfig = {
  readonly methodName: string;
  readonly service: typeof RateLimitDiscoveryService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof envoy_api_v2_discovery_pb.DiscoveryRequest;
  readonly responseType: typeof envoy_api_v2_discovery_pb.DiscoveryResponse;
};

export class RateLimitDiscoveryService {
  static readonly serviceName: string;
  static readonly StreamRateLimitConfig: RateLimitDiscoveryServiceStreamRateLimitConfig;
  static readonly DeltaRateLimitConfig: RateLimitDiscoveryServiceDeltaRateLimitConfig;
  static readonly FetchRateLimitConfig: RateLimitDiscoveryServiceFetchRateLimitConfig;
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

export class RateLimitDiscoveryServiceClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  streamRateLimitConfig(metadata?: grpc.Metadata): BidirectionalStream<envoy_api_v2_discovery_pb.DiscoveryRequest, envoy_api_v2_discovery_pb.DiscoveryResponse>;
  deltaRateLimitConfig(metadata?: grpc.Metadata): BidirectionalStream<envoy_api_v2_discovery_pb.DeltaDiscoveryRequest, envoy_api_v2_discovery_pb.DeltaDiscoveryResponse>;
  fetchRateLimitConfig(
    requestMessage: envoy_api_v2_discovery_pb.DiscoveryRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: envoy_api_v2_discovery_pb.DiscoveryResponse|null) => void
  ): UnaryResponse;
  fetchRateLimitConfig(
    requestMessage: envoy_api_v2_discovery_pb.DiscoveryRequest,
    callback: (error: ServiceError|null, responseMessage: envoy_api_v2_discovery_pb.DiscoveryResponse|null) => void
  ): UnaryResponse;
}

