// package: glooeeapi.solo.io
// file: github.com/solo-io/solo-projects/projects/grpcserver/api/v1/artifact.proto

import * as github_com_solo_io_solo_projects_projects_grpcserver_api_v1_artifact_pb from "../../../../../../../github.com/solo-io/solo-projects/projects/grpcserver/api/v1/artifact_pb";
import {grpc} from "@improbable-eng/grpc-web";

type ArtifactApiGetArtifact = {
  readonly methodName: string;
  readonly service: typeof ArtifactApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_artifact_pb.GetArtifactRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_artifact_pb.GetArtifactResponse;
};

type ArtifactApiListArtifacts = {
  readonly methodName: string;
  readonly service: typeof ArtifactApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_artifact_pb.ListArtifactsRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_artifact_pb.ListArtifactsResponse;
};

type ArtifactApiCreateArtifact = {
  readonly methodName: string;
  readonly service: typeof ArtifactApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_artifact_pb.CreateArtifactRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_artifact_pb.CreateArtifactResponse;
};

type ArtifactApiUpdateArtifact = {
  readonly methodName: string;
  readonly service: typeof ArtifactApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_artifact_pb.UpdateArtifactRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_artifact_pb.UpdateArtifactResponse;
};

type ArtifactApiDeleteArtifact = {
  readonly methodName: string;
  readonly service: typeof ArtifactApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_artifact_pb.DeleteArtifactRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_grpcserver_api_v1_artifact_pb.DeleteArtifactResponse;
};

export class ArtifactApi {
  static readonly serviceName: string;
  static readonly GetArtifact: ArtifactApiGetArtifact;
  static readonly ListArtifacts: ArtifactApiListArtifacts;
  static readonly CreateArtifact: ArtifactApiCreateArtifact;
  static readonly UpdateArtifact: ArtifactApiUpdateArtifact;
  static readonly DeleteArtifact: ArtifactApiDeleteArtifact;
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

export class ArtifactApiClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  getArtifact(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_artifact_pb.GetArtifactRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_artifact_pb.GetArtifactResponse|null) => void
  ): UnaryResponse;
  getArtifact(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_artifact_pb.GetArtifactRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_artifact_pb.GetArtifactResponse|null) => void
  ): UnaryResponse;
  listArtifacts(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_artifact_pb.ListArtifactsRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_artifact_pb.ListArtifactsResponse|null) => void
  ): UnaryResponse;
  listArtifacts(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_artifact_pb.ListArtifactsRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_artifact_pb.ListArtifactsResponse|null) => void
  ): UnaryResponse;
  createArtifact(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_artifact_pb.CreateArtifactRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_artifact_pb.CreateArtifactResponse|null) => void
  ): UnaryResponse;
  createArtifact(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_artifact_pb.CreateArtifactRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_artifact_pb.CreateArtifactResponse|null) => void
  ): UnaryResponse;
  updateArtifact(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_artifact_pb.UpdateArtifactRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_artifact_pb.UpdateArtifactResponse|null) => void
  ): UnaryResponse;
  updateArtifact(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_artifact_pb.UpdateArtifactRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_artifact_pb.UpdateArtifactResponse|null) => void
  ): UnaryResponse;
  deleteArtifact(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_artifact_pb.DeleteArtifactRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_artifact_pb.DeleteArtifactResponse|null) => void
  ): UnaryResponse;
  deleteArtifact(
    requestMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_artifact_pb.DeleteArtifactRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_artifact_pb.DeleteArtifactResponse|null) => void
  ): UnaryResponse;
}

