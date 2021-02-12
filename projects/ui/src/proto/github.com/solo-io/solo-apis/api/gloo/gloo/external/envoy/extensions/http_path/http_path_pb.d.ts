/* eslint-disable */
// package: envoy.config.health_checker.http_path.v2
// file: github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/extensions/http_path/http_path.proto

import * as jspb from "google-protobuf";
import * as udpa_annotations_status_pb from "../../../../../../../../../../udpa/annotations/status_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_health_check_pb from "../../../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/config/core/v3/health_check_pb";
import * as extproto_ext_pb from "../../../../../../../../../../extproto/ext_pb";

export class HttpPath extends jspb.Message {
  hasHttpHealthCheck(): boolean;
  clearHttpHealthCheck(): void;
  getHttpHealthCheck(): github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_health_check_pb.HealthCheck.HttpHealthCheck | undefined;
  setHttpHealthCheck(value?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_health_check_pb.HealthCheck.HttpHealthCheck): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HttpPath.AsObject;
  static toObject(includeInstance: boolean, msg: HttpPath): HttpPath.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: HttpPath, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HttpPath;
  static deserializeBinaryFromReader(message: HttpPath, reader: jspb.BinaryReader): HttpPath;
}

export namespace HttpPath {
  export type AsObject = {
    httpHealthCheck?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_health_check_pb.HealthCheck.HttpHealthCheck.AsObject,
  }
}
