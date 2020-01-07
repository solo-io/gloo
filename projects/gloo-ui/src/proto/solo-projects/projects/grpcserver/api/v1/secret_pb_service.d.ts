// package: glooeeapi.solo.io
// file: solo-projects/projects/grpcserver/api/v1/secret.proto

import * as solo_projects_projects_grpcserver_api_v1_secret_pb from "../../../../../solo-projects/projects/grpcserver/api/v1/secret_pb";
import {grpc} from "@improbable-eng/grpc-web";

type SecretApiGetSecret = {
  readonly methodName: string;
  readonly service: typeof SecretApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof solo_projects_projects_grpcserver_api_v1_secret_pb.GetSecretRequest;
  readonly responseType: typeof solo_projects_projects_grpcserver_api_v1_secret_pb.GetSecretResponse;
};

type SecretApiListSecrets = {
  readonly methodName: string;
  readonly service: typeof SecretApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof solo_projects_projects_grpcserver_api_v1_secret_pb.ListSecretsRequest;
  readonly responseType: typeof solo_projects_projects_grpcserver_api_v1_secret_pb.ListSecretsResponse;
};

type SecretApiCreateSecret = {
  readonly methodName: string;
  readonly service: typeof SecretApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof solo_projects_projects_grpcserver_api_v1_secret_pb.CreateSecretRequest;
  readonly responseType: typeof solo_projects_projects_grpcserver_api_v1_secret_pb.CreateSecretResponse;
};

type SecretApiUpdateSecret = {
  readonly methodName: string;
  readonly service: typeof SecretApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof solo_projects_projects_grpcserver_api_v1_secret_pb.UpdateSecretRequest;
  readonly responseType: typeof solo_projects_projects_grpcserver_api_v1_secret_pb.UpdateSecretResponse;
};

type SecretApiDeleteSecret = {
  readonly methodName: string;
  readonly service: typeof SecretApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof solo_projects_projects_grpcserver_api_v1_secret_pb.DeleteSecretRequest;
  readonly responseType: typeof solo_projects_projects_grpcserver_api_v1_secret_pb.DeleteSecretResponse;
};

export class SecretApi {
  static readonly serviceName: string;
  static readonly GetSecret: SecretApiGetSecret;
  static readonly ListSecrets: SecretApiListSecrets;
  static readonly CreateSecret: SecretApiCreateSecret;
  static readonly UpdateSecret: SecretApiUpdateSecret;
  static readonly DeleteSecret: SecretApiDeleteSecret;
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

export class SecretApiClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  getSecret(
    requestMessage: solo_projects_projects_grpcserver_api_v1_secret_pb.GetSecretRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: solo_projects_projects_grpcserver_api_v1_secret_pb.GetSecretResponse|null) => void
  ): UnaryResponse;
  getSecret(
    requestMessage: solo_projects_projects_grpcserver_api_v1_secret_pb.GetSecretRequest,
    callback: (error: ServiceError|null, responseMessage: solo_projects_projects_grpcserver_api_v1_secret_pb.GetSecretResponse|null) => void
  ): UnaryResponse;
  listSecrets(
    requestMessage: solo_projects_projects_grpcserver_api_v1_secret_pb.ListSecretsRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: solo_projects_projects_grpcserver_api_v1_secret_pb.ListSecretsResponse|null) => void
  ): UnaryResponse;
  listSecrets(
    requestMessage: solo_projects_projects_grpcserver_api_v1_secret_pb.ListSecretsRequest,
    callback: (error: ServiceError|null, responseMessage: solo_projects_projects_grpcserver_api_v1_secret_pb.ListSecretsResponse|null) => void
  ): UnaryResponse;
  createSecret(
    requestMessage: solo_projects_projects_grpcserver_api_v1_secret_pb.CreateSecretRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: solo_projects_projects_grpcserver_api_v1_secret_pb.CreateSecretResponse|null) => void
  ): UnaryResponse;
  createSecret(
    requestMessage: solo_projects_projects_grpcserver_api_v1_secret_pb.CreateSecretRequest,
    callback: (error: ServiceError|null, responseMessage: solo_projects_projects_grpcserver_api_v1_secret_pb.CreateSecretResponse|null) => void
  ): UnaryResponse;
  updateSecret(
    requestMessage: solo_projects_projects_grpcserver_api_v1_secret_pb.UpdateSecretRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: solo_projects_projects_grpcserver_api_v1_secret_pb.UpdateSecretResponse|null) => void
  ): UnaryResponse;
  updateSecret(
    requestMessage: solo_projects_projects_grpcserver_api_v1_secret_pb.UpdateSecretRequest,
    callback: (error: ServiceError|null, responseMessage: solo_projects_projects_grpcserver_api_v1_secret_pb.UpdateSecretResponse|null) => void
  ): UnaryResponse;
  deleteSecret(
    requestMessage: solo_projects_projects_grpcserver_api_v1_secret_pb.DeleteSecretRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: solo_projects_projects_grpcserver_api_v1_secret_pb.DeleteSecretResponse|null) => void
  ): UnaryResponse;
  deleteSecret(
    requestMessage: solo_projects_projects_grpcserver_api_v1_secret_pb.DeleteSecretRequest,
    callback: (error: ServiceError|null, responseMessage: solo_projects_projects_grpcserver_api_v1_secret_pb.DeleteSecretResponse|null) => void
  ): UnaryResponse;
}

