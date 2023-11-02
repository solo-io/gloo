/* eslint-disable */
// package: envoy.config.tap.output_sink.v3
// file: github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/config/tap/output_sink/v3/grpc_output_sink.proto

import * as jspb from "google-protobuf";
import * as validate_validate_pb from "../../../../../../../../../../../../validate/validate_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_grpc_service_pb from "../../../../../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/config/core/v3/grpc_service_pb";

export class GrpcOutputSink extends jspb.Message {
  hasGrpcService(): boolean;
  clearGrpcService(): void;
  getGrpcService(): github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_grpc_service_pb.GrpcService | undefined;
  setGrpcService(value?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_grpc_service_pb.GrpcService): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GrpcOutputSink.AsObject;
  static toObject(includeInstance: boolean, msg: GrpcOutputSink): GrpcOutputSink.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GrpcOutputSink, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GrpcOutputSink;
  static deserializeBinaryFromReader(message: GrpcOutputSink, reader: jspb.BinaryReader): GrpcOutputSink;
}

export namespace GrpcOutputSink {
  export type AsObject = {
    grpcService?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_grpc_service_pb.GrpcService.AsObject,
  }
}
