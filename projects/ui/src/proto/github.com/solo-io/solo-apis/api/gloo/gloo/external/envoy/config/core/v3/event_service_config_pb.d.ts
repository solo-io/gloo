/* eslint-disable */
// package: solo.io.envoy.config.core.v3
// file: github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/config/core/v3/event_service_config.proto

import * as jspb from "google-protobuf";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_grpc_service_pb from "../../../../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/config/core/v3/grpc_service_pb";
import * as udpa_annotations_status_pb from "../../../../../../../../../../../udpa/annotations/status_pb";
import * as udpa_annotations_versioning_pb from "../../../../../../../../../../../udpa/annotations/versioning_pb";
import * as validate_validate_pb from "../../../../../../../../../../../validate/validate_pb";
import * as extproto_ext_pb from "../../../../../../../../../../../extproto/ext_pb";

export class EventServiceConfig extends jspb.Message {
  hasGrpcService(): boolean;
  clearGrpcService(): void;
  getGrpcService(): github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_grpc_service_pb.GrpcService | undefined;
  setGrpcService(value?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_grpc_service_pb.GrpcService): void;

  getConfigSourceSpecifierCase(): EventServiceConfig.ConfigSourceSpecifierCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EventServiceConfig.AsObject;
  static toObject(includeInstance: boolean, msg: EventServiceConfig): EventServiceConfig.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: EventServiceConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EventServiceConfig;
  static deserializeBinaryFromReader(message: EventServiceConfig, reader: jspb.BinaryReader): EventServiceConfig;
}

export namespace EventServiceConfig {
  export type AsObject = {
    grpcService?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_grpc_service_pb.GrpcService.AsObject,
  }

  export enum ConfigSourceSpecifierCase {
    CONFIG_SOURCE_SPECIFIER_NOT_SET = 0,
    GRPC_SERVICE = 1,
  }
}
