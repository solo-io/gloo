/* eslint-disable */
// package: envoy.extensions.cache.grpc.v2
// file: github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/extensions/cache/grpc/config.proto

import * as jspb from "google-protobuf";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_base_pb from "../../../../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/config/core/v3/base_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_grpc_service_pb from "../../../../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/config/core/v3/grpc_service_pb";
import * as udpa_annotations_status_pb from "../../../../../../../../../../../udpa/annotations/status_pb";
import * as google_protobuf_duration_pb from "google-protobuf/google/protobuf/duration_pb";
import * as validate_validate_pb from "../../../../../../../../../../../validate/validate_pb";

export class GrpcCacheConfig extends jspb.Message {
  hasService(): boolean;
  clearService(): void;
  getService(): github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_grpc_service_pb.GrpcService | undefined;
  setService(value?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_grpc_service_pb.GrpcService): void;

  hasTimeout(): boolean;
  clearTimeout(): void;
  getTimeout(): google_protobuf_duration_pb.Duration | undefined;
  setTimeout(value?: google_protobuf_duration_pb.Duration): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GrpcCacheConfig.AsObject;
  static toObject(includeInstance: boolean, msg: GrpcCacheConfig): GrpcCacheConfig.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GrpcCacheConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GrpcCacheConfig;
  static deserializeBinaryFromReader(message: GrpcCacheConfig, reader: jspb.BinaryReader): GrpcCacheConfig;
}

export namespace GrpcCacheConfig {
  export type AsObject = {
    service?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_grpc_service_pb.GrpcService.AsObject,
    timeout?: google_protobuf_duration_pb.Duration.AsObject,
  }
}
