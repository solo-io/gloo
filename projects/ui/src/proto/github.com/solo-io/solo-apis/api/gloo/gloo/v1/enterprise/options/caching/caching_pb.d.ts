/* eslint-disable */
// package: caching.options.gloo.solo.io
// file: github.com/solo-io/solo-apis/api/gloo/gloo/v1/enterprise/options/caching/caching.proto

import * as jspb from "google-protobuf";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_type_matcher_v3_string_pb from "../../../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/type/matcher/v3/string_pb";
import * as github_com_solo_io_solo_kit_api_v1_ref_pb from "../../../../../../../../../../github.com/solo-io/solo-kit/api/v1/ref_pb";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";
import * as google_protobuf_duration_pb from "google-protobuf/google/protobuf/duration_pb";
import * as extproto_ext_pb from "../../../../../../../../../../extproto/ext_pb";

export class Settings extends jspb.Message {
  hasCachingServiceRef(): boolean;
  clearCachingServiceRef(): void;
  getCachingServiceRef(): github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef | undefined;
  setCachingServiceRef(value?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef): void;

  clearAllowedVaryHeadersList(): void;
  getAllowedVaryHeadersList(): Array<github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_type_matcher_v3_string_pb.StringMatcher>;
  setAllowedVaryHeadersList(value: Array<github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_type_matcher_v3_string_pb.StringMatcher>): void;
  addAllowedVaryHeaders(value?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_type_matcher_v3_string_pb.StringMatcher, index?: number): github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_type_matcher_v3_string_pb.StringMatcher;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Settings.AsObject;
  static toObject(includeInstance: boolean, msg: Settings): Settings.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Settings, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Settings;
  static deserializeBinaryFromReader(message: Settings, reader: jspb.BinaryReader): Settings;
}

export namespace Settings {
  export type AsObject = {
    cachingServiceRef?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject,
    allowedVaryHeadersList: Array<github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_type_matcher_v3_string_pb.StringMatcher.AsObject>,
  }
}
