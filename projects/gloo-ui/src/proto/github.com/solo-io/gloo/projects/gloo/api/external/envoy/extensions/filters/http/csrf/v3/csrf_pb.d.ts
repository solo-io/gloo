/* eslint-disable */
// package: solo.io.envoy.extensions.filters.http.csrf.v3
// file: github.com/solo-io/gloo/projects/gloo/api/external/envoy/extensions/filters/http/csrf/v3/csrf.proto

import * as jspb from "google-protobuf";
import * as github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb from "../../../../../../../../../../../../../github.com/solo-io/gloo/projects/gloo/api/external/envoy/config/core/v3/base_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_matcher_v3_string_pb from "../../../../../../../../../../../../../github.com/solo-io/gloo/projects/gloo/api/external/envoy/type/matcher/v3/string_pb";
import * as validate_validate_pb from "../../../../../../../../../../../../../validate/validate_pb";
import * as extproto_ext_pb from "../../../../../../../../../../../../../protoc-gen-ext/extproto/ext_pb";

export class CsrfPolicy extends jspb.Message {
  hasFilterEnabled(): boolean;
  clearFilterEnabled(): void;
  getFilterEnabled(): github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb.RuntimeFractionalPercent | undefined;
  setFilterEnabled(value?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb.RuntimeFractionalPercent): void;

  hasShadowEnabled(): boolean;
  clearShadowEnabled(): void;
  getShadowEnabled(): github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb.RuntimeFractionalPercent | undefined;
  setShadowEnabled(value?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb.RuntimeFractionalPercent): void;

  clearAdditionalOriginsList(): void;
  getAdditionalOriginsList(): Array<github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_matcher_v3_string_pb.StringMatcher>;
  setAdditionalOriginsList(value: Array<github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_matcher_v3_string_pb.StringMatcher>): void;
  addAdditionalOrigins(value?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_matcher_v3_string_pb.StringMatcher, index?: number): github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_matcher_v3_string_pb.StringMatcher;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CsrfPolicy.AsObject;
  static toObject(includeInstance: boolean, msg: CsrfPolicy): CsrfPolicy.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CsrfPolicy, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CsrfPolicy;
  static deserializeBinaryFromReader(message: CsrfPolicy, reader: jspb.BinaryReader): CsrfPolicy;
}

export namespace CsrfPolicy {
  export type AsObject = {
    filterEnabled?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb.RuntimeFractionalPercent.AsObject,
    shadowEnabled?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb.RuntimeFractionalPercent.AsObject,
    additionalOriginsList: Array<github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_matcher_v3_string_pb.StringMatcher.AsObject>,
  }
}
