/* eslint-disable */
// package: envoy.config.tap.output_sink.v3
// file: github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/config/tap/output_sink/v3/http_output_sink.proto

import * as jspb from "google-protobuf";
import * as validate_validate_pb from "../../../../../../../../../../../../validate/validate_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_http_uri_pb from "../../../../../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/config/core/v3/http_uri_pb";

export class HttpOutputSink extends jspb.Message {
  hasServerUri(): boolean;
  clearServerUri(): void;
  getServerUri(): github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_http_uri_pb.HttpUri | undefined;
  setServerUri(value?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_http_uri_pb.HttpUri): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HttpOutputSink.AsObject;
  static toObject(includeInstance: boolean, msg: HttpOutputSink): HttpOutputSink.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: HttpOutputSink, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HttpOutputSink;
  static deserializeBinaryFromReader(message: HttpOutputSink, reader: jspb.BinaryReader): HttpOutputSink;
}

export namespace HttpOutputSink {
  export type AsObject = {
    serverUri?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_http_uri_pb.HttpUri.AsObject,
  }
}
