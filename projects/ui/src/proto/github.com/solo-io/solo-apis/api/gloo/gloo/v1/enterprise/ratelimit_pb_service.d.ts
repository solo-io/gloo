// package: glooe.solo.io
// file: github.com/solo-io/solo-apis/api/gloo/gloo/v1/enterprise/ratelimit.proto

import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_ratelimit_pb from "../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/enterprise/ratelimit_pb";
import * as github_com_solo_io_solo_kit_api_external_envoy_api_v2_discovery_pb from "../../../../../../../../github.com/solo-io/solo-kit/api/external/envoy/api/v2/discovery_pb";
import {grpc} from "@improbable-eng/grpc-web";

type RateLimitDiscoveryServiceStreamRateLimitConfig = {
  readonly methodName: string;
  readonly service: typeof RateLimitDiscoveryService;
  readonly requestStream: true;
  readonly responseStream: true;
  readonly requestType: typeof github_com_solo_io_solo_kit_api_external_envoy_api_v2_discovery_pb.DiscoveryRequest;
  readonly responseType: typeof github_com_solo_io_solo_kit_api_external_envoy_api_v2_discovery_pb.DiscoveryResponse;
};

type RateLimitDiscoveryServiceDeltaRateLimitConfig = {
  readonly methodName: string;
  readonly service: typeof RateLimitDiscoveryService;
  readonly requestStream: true;
  readonly responseStream: true;
  readonly requestType: typeof github_com_solo_io_solo_kit_api_external_envoy_api_v2_discovery_pb.DeltaDiscoveryRequest;
  readonly responseType: typeof github_com_solo_io_solo_kit_api_external_envoy_api_v2_discovery_pb.DeltaDiscoveryResponse;
};

type RateLimitDiscoveryServiceFetchRateLimitConfig = {
  readonly methodName: string;
  readonly service: typeof RateLimitDiscoveryService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_kit_api_external_envoy_api_v2_discovery_pb.DiscoveryRequest;
  readonly responseType: typeof github_com_solo_io_solo_kit_api_external_envoy_api_v2_discovery_pb.DiscoveryResponse;
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
  streamRateLimitConfig(metadata?: grpc.Metadata): BidirectionalStream<github_com_solo_io_solo_kit_api_external_envoy_api_v2_discovery_pb.DiscoveryRequest, github_com_solo_io_solo_kit_api_external_envoy_api_v2_discovery_pb.DiscoveryResponse>;
  deltaRateLimitConfig(metadata?: grpc.Metadata): BidirectionalStream<github_com_solo_io_solo_kit_api_external_envoy_api_v2_discovery_pb.DeltaDiscoveryRequest, github_com_solo_io_solo_kit_api_external_envoy_api_v2_discovery_pb.DeltaDiscoveryResponse>;
  fetchRateLimitConfig(
    requestMessage: github_com_solo_io_solo_kit_api_external_envoy_api_v2_discovery_pb.DiscoveryRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_kit_api_external_envoy_api_v2_discovery_pb.DiscoveryResponse|null) => void
  ): UnaryResponse;
  fetchRateLimitConfig(
    requestMessage: github_com_solo_io_solo_kit_api_external_envoy_api_v2_discovery_pb.DiscoveryRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_kit_api_external_envoy_api_v2_discovery_pb.DiscoveryResponse|null) => void
  ): UnaryResponse;
}

