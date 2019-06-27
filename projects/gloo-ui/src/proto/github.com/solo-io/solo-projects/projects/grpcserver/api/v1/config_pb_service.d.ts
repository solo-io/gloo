// package: glooeeapi.solo.io
// file: github.com/solo-io/solo-projects/projects/grpcserver/api/v1/config.proto

import * as github_com_solo_io_solo_projects_projects_grpcserver_api_v1_config_pb from "../../../../../../../github.com/solo-io/solo-projects/projects/grpcserver/api/v1/config_pb";
import {grpc} from "@improbable-eng/grpc-web";

type ConfigApiGetVersion = {
  readonly methodName: string;
  readonly service: typeof ConfigApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_config_pb.GetVersionRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_config_pb.GetVersionResponse;
};

type ConfigApiGetOAuthEndpoint = {
  readonly methodName: string;
  readonly service: typeof ConfigApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_config_pb.GetOAuthEndpointRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_config_pb.GetOAuthEndpointResponse;
};

type ConfigApiGetIsLicenseValid = {
  readonly methodName: string;
  readonly service: typeof ConfigApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_config_pb.GetIsLicenseValidRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_config_pb.GetIsLicenseValidResponse;
};

type ConfigApiGetSettings = {
  readonly methodName: string;
  readonly service: typeof ConfigApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_config_pb.GetSettingsRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_config_pb.GetSettingsResponse;
};

type ConfigApiUpdateSettings = {
  readonly methodName: string;
  readonly service: typeof ConfigApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_config_pb.UpdateSettingsRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_config_pb.UpdateSettingsResponse;
};

type ConfigApiListNamespaces = {
  readonly methodName: string;
  readonly service: typeof ConfigApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_config_pb.ListNamespacesRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_config_pb.ListNamespacesResponse;
};

export class ConfigApi {
  static readonly serviceName: string;
  static readonly GetVersion: ConfigApiGetVersion;
  static readonly GetOAuthEndpoint: ConfigApiGetOAuthEndpoint;
  static readonly GetIsLicenseValid: ConfigApiGetIsLicenseValid;
  static readonly GetSettings: ConfigApiGetSettings;
  static readonly UpdateSettings: ConfigApiUpdateSettings;
  static readonly ListNamespaces: ConfigApiListNamespaces;
}

export type ServiceError = { message: string, code: number; metadata: grpc.Metadata }
export type Status = { details: string, code: number; metadata: grpc.Metadata }

interface UnaryResponse {
  cancel(): void;
}
interface ResponseStream<T> {
  cancel(): void;
  on(type: 'data', handler: (message: T) => void): ResponseStream<T>;
  on(type: 'end', handler: () => void): ResponseStream<T>;
  on(type: 'status', handler: (status: Status) => void): ResponseStream<T>;
}
interface RequestStream<T> {
  write(message: T): RequestStream<T>;
  end(): void;
  cancel(): void;
  on(type: 'end', handler: () => void): RequestStream<T>;
  on(type: 'status', handler: (status: Status) => void): RequestStream<T>;
}
interface BidirectionalStream<ReqT, ResT> {
  write(message: ReqT): BidirectionalStream<ReqT, ResT>;
  end(): void;
  cancel(): void;
  on(type: 'data', handler: (message: ResT) => void): BidirectionalStream<ReqT, ResT>;
  on(type: 'end', handler: () => void): BidirectionalStream<ReqT, ResT>;
  on(type: 'status', handler: (status: Status) => void): BidirectionalStream<ReqT, ResT>;
}

export class ConfigApiClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  getVersion(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_config_pb.GetVersionRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_config_pb.GetVersionResponse|null) => void
  ): UnaryResponse;
  getVersion(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_config_pb.GetVersionRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_config_pb.GetVersionResponse|null) => void
  ): UnaryResponse;
  getOAuthEndpoint(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_config_pb.GetOAuthEndpointRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_config_pb.GetOAuthEndpointResponse|null) => void
  ): UnaryResponse;
  getOAuthEndpoint(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_config_pb.GetOAuthEndpointRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_config_pb.GetOAuthEndpointResponse|null) => void
  ): UnaryResponse;
  getIsLicenseValid(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_config_pb.GetIsLicenseValidRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_config_pb.GetIsLicenseValidResponse|null) => void
  ): UnaryResponse;
  getIsLicenseValid(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_config_pb.GetIsLicenseValidRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_config_pb.GetIsLicenseValidResponse|null) => void
  ): UnaryResponse;
  getSettings(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_config_pb.GetSettingsRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_config_pb.GetSettingsResponse|null) => void
  ): UnaryResponse;
  getSettings(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_config_pb.GetSettingsRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_config_pb.GetSettingsResponse|null) => void
  ): UnaryResponse;
  updateSettings(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_config_pb.UpdateSettingsRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_config_pb.UpdateSettingsResponse|null) => void
  ): UnaryResponse;
  updateSettings(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_config_pb.UpdateSettingsRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_config_pb.UpdateSettingsResponse|null) => void
  ): UnaryResponse;
  listNamespaces(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_config_pb.ListNamespacesRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_config_pb.ListNamespacesResponse|null) => void
  ): UnaryResponse;
  listNamespaces(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_config_pb.ListNamespacesRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_config_pb.ListNamespacesResponse|null) => void
  ): UnaryResponse;
}

