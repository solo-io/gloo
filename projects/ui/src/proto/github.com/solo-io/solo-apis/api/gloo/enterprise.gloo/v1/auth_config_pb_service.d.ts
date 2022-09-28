// package: enterprise.gloo.solo.io
// file: github.com/solo-io/solo-apis/api/gloo/enterprise.gloo/v1/auth_config.proto

import * as github_com_solo_io_solo_apis_api_gloo_enterprise_gloo_v1_auth_config_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/enterprise.gloo/v1/auth_config_pb";
import * as github_com_solo_io_solo_kit_api_external_envoy_api_v2_discovery_pb from "../../../../../../../github.com/solo-io/solo-kit/api/external/envoy/api/v2/discovery_pb";
import {grpc} from "@improbable-eng/grpc-web";

type ExtAuthDiscoveryServiceStreamExtAuthConfig = {
  readonly methodName: string;
  readonly service: typeof ExtAuthDiscoveryService;
  readonly requestStream: true;
  readonly responseStream: true;
  readonly requestType: typeof github_com_solo_io_solo_kit_api_external_envoy_api_v2_discovery_pb.DiscoveryRequest;
  readonly responseType: typeof github_com_solo_io_solo_kit_api_external_envoy_api_v2_discovery_pb.DiscoveryResponse;
};

type ExtAuthDiscoveryServiceDeltaExtAuthConfig = {
  readonly methodName: string;
  readonly service: typeof ExtAuthDiscoveryService;
  readonly requestStream: true;
  readonly responseStream: true;
  readonly requestType: typeof github_com_solo_io_solo_kit_api_external_envoy_api_v2_discovery_pb.DeltaDiscoveryRequest;
  readonly responseType: typeof github_com_solo_io_solo_kit_api_external_envoy_api_v2_discovery_pb.DeltaDiscoveryResponse;
};

type ExtAuthDiscoveryServiceFetchExtAuthConfig = {
  readonly methodName: string;
  readonly service: typeof ExtAuthDiscoveryService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_kit_api_external_envoy_api_v2_discovery_pb.DiscoveryRequest;
  readonly responseType: typeof github_com_solo_io_solo_kit_api_external_envoy_api_v2_discovery_pb.DiscoveryResponse;
};

export class ExtAuthDiscoveryService {
  static readonly serviceName: string;
  static readonly StreamExtAuthConfig: ExtAuthDiscoveryServiceStreamExtAuthConfig;
  static readonly DeltaExtAuthConfig: ExtAuthDiscoveryServiceDeltaExtAuthConfig;
  static readonly FetchExtAuthConfig: ExtAuthDiscoveryServiceFetchExtAuthConfig;
}

type ApiKeyServiceCreate = {
  readonly methodName: string;
  readonly service: typeof ApiKeyService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_apis_api_gloo_enterprise_gloo_v1_auth_config_pb.ApiKeyCreateRequest;
  readonly responseType: typeof github_com_solo_io_solo_apis_api_gloo_enterprise_gloo_v1_auth_config_pb.ApiKeyCreateResponse;
};

type ApiKeyServiceRead = {
  readonly methodName: string;
  readonly service: typeof ApiKeyService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_apis_api_gloo_enterprise_gloo_v1_auth_config_pb.ApiKeyReadRequest;
  readonly responseType: typeof github_com_solo_io_solo_apis_api_gloo_enterprise_gloo_v1_auth_config_pb.ApiKeyReadResponse;
};

type ApiKeyServiceUpdate = {
  readonly methodName: string;
  readonly service: typeof ApiKeyService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_apis_api_gloo_enterprise_gloo_v1_auth_config_pb.ApiKeyUpdateRequest;
  readonly responseType: typeof github_com_solo_io_solo_apis_api_gloo_enterprise_gloo_v1_auth_config_pb.ApiKeyUpdateResponse;
};

type ApiKeyServiceDelete = {
  readonly methodName: string;
  readonly service: typeof ApiKeyService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_apis_api_gloo_enterprise_gloo_v1_auth_config_pb.ApiKeyDeleteRequest;
  readonly responseType: typeof github_com_solo_io_solo_apis_api_gloo_enterprise_gloo_v1_auth_config_pb.ApiKeyDeleteResponse;
};

export class ApiKeyService {
  static readonly serviceName: string;
  static readonly Create: ApiKeyServiceCreate;
  static readonly Read: ApiKeyServiceRead;
  static readonly Update: ApiKeyServiceUpdate;
  static readonly Delete: ApiKeyServiceDelete;
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
  streamExtAuthConfig(metadata?: grpc.Metadata): BidirectionalStream<github_com_solo_io_solo_kit_api_external_envoy_api_v2_discovery_pb.DiscoveryRequest, github_com_solo_io_solo_kit_api_external_envoy_api_v2_discovery_pb.DiscoveryResponse>;
  deltaExtAuthConfig(metadata?: grpc.Metadata): BidirectionalStream<github_com_solo_io_solo_kit_api_external_envoy_api_v2_discovery_pb.DeltaDiscoveryRequest, github_com_solo_io_solo_kit_api_external_envoy_api_v2_discovery_pb.DeltaDiscoveryResponse>;
  fetchExtAuthConfig(
    requestMessage: github_com_solo_io_solo_kit_api_external_envoy_api_v2_discovery_pb.DiscoveryRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_kit_api_external_envoy_api_v2_discovery_pb.DiscoveryResponse|null) => void
  ): UnaryResponse;
  fetchExtAuthConfig(
    requestMessage: github_com_solo_io_solo_kit_api_external_envoy_api_v2_discovery_pb.DiscoveryRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_kit_api_external_envoy_api_v2_discovery_pb.DiscoveryResponse|null) => void
  ): UnaryResponse;
}

export class ApiKeyServiceClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  create(
    requestMessage: github_com_solo_io_solo_apis_api_gloo_enterprise_gloo_v1_auth_config_pb.ApiKeyCreateRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_apis_api_gloo_enterprise_gloo_v1_auth_config_pb.ApiKeyCreateResponse|null) => void
  ): UnaryResponse;
  create(
    requestMessage: github_com_solo_io_solo_apis_api_gloo_enterprise_gloo_v1_auth_config_pb.ApiKeyCreateRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_apis_api_gloo_enterprise_gloo_v1_auth_config_pb.ApiKeyCreateResponse|null) => void
  ): UnaryResponse;
  read(
    requestMessage: github_com_solo_io_solo_apis_api_gloo_enterprise_gloo_v1_auth_config_pb.ApiKeyReadRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_apis_api_gloo_enterprise_gloo_v1_auth_config_pb.ApiKeyReadResponse|null) => void
  ): UnaryResponse;
  read(
    requestMessage: github_com_solo_io_solo_apis_api_gloo_enterprise_gloo_v1_auth_config_pb.ApiKeyReadRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_apis_api_gloo_enterprise_gloo_v1_auth_config_pb.ApiKeyReadResponse|null) => void
  ): UnaryResponse;
  update(
    requestMessage: github_com_solo_io_solo_apis_api_gloo_enterprise_gloo_v1_auth_config_pb.ApiKeyUpdateRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_apis_api_gloo_enterprise_gloo_v1_auth_config_pb.ApiKeyUpdateResponse|null) => void
  ): UnaryResponse;
  update(
    requestMessage: github_com_solo_io_solo_apis_api_gloo_enterprise_gloo_v1_auth_config_pb.ApiKeyUpdateRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_apis_api_gloo_enterprise_gloo_v1_auth_config_pb.ApiKeyUpdateResponse|null) => void
  ): UnaryResponse;
  delete(
    requestMessage: github_com_solo_io_solo_apis_api_gloo_enterprise_gloo_v1_auth_config_pb.ApiKeyDeleteRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_apis_api_gloo_enterprise_gloo_v1_auth_config_pb.ApiKeyDeleteResponse|null) => void
  ): UnaryResponse;
  delete(
    requestMessage: github_com_solo_io_solo_apis_api_gloo_enterprise_gloo_v1_auth_config_pb.ApiKeyDeleteRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_apis_api_gloo_enterprise_gloo_v1_auth_config_pb.ApiKeyDeleteResponse|null) => void
  ): UnaryResponse;
}

