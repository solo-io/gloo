/* eslint-disable */
// package: solo.io.envoy.config.common.mutation_rules.v3
// file: github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/config/common/mutation_rules/v3/mutation_rules.proto

import * as jspb from "google-protobuf";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_base_pb from "../../../../../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/config/core/v3/base_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_type_matcher_v3_regex_pb from "../../../../../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/type/matcher/v3/regex_pb";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";
import * as udpa_annotations_status_pb from "../../../../../../../../../../../../udpa/annotations/status_pb";
import * as validate_validate_pb from "../../../../../../../../../../../../validate/validate_pb";
import * as extproto_ext_pb from "../../../../../../../../../../../../extproto/ext_pb";

export class HeaderMutationRules extends jspb.Message {
  hasAllowAllRouting(): boolean;
  clearAllowAllRouting(): void;
  getAllowAllRouting(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setAllowAllRouting(value?: google_protobuf_wrappers_pb.BoolValue): void;

  hasAllowEnvoy(): boolean;
  clearAllowEnvoy(): void;
  getAllowEnvoy(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setAllowEnvoy(value?: google_protobuf_wrappers_pb.BoolValue): void;

  hasDisallowSystem(): boolean;
  clearDisallowSystem(): void;
  getDisallowSystem(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setDisallowSystem(value?: google_protobuf_wrappers_pb.BoolValue): void;

  hasDisallowAll(): boolean;
  clearDisallowAll(): void;
  getDisallowAll(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setDisallowAll(value?: google_protobuf_wrappers_pb.BoolValue): void;

  hasAllowExpression(): boolean;
  clearAllowExpression(): void;
  getAllowExpression(): github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_type_matcher_v3_regex_pb.RegexMatcher | undefined;
  setAllowExpression(value?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_type_matcher_v3_regex_pb.RegexMatcher): void;

  hasDisallowExpression(): boolean;
  clearDisallowExpression(): void;
  getDisallowExpression(): github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_type_matcher_v3_regex_pb.RegexMatcher | undefined;
  setDisallowExpression(value?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_type_matcher_v3_regex_pb.RegexMatcher): void;

  hasDisallowIsError(): boolean;
  clearDisallowIsError(): void;
  getDisallowIsError(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setDisallowIsError(value?: google_protobuf_wrappers_pb.BoolValue): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HeaderMutationRules.AsObject;
  static toObject(includeInstance: boolean, msg: HeaderMutationRules): HeaderMutationRules.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: HeaderMutationRules, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HeaderMutationRules;
  static deserializeBinaryFromReader(message: HeaderMutationRules, reader: jspb.BinaryReader): HeaderMutationRules;
}

export namespace HeaderMutationRules {
  export type AsObject = {
    allowAllRouting?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    allowEnvoy?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    disallowSystem?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    disallowAll?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    allowExpression?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_type_matcher_v3_regex_pb.RegexMatcher.AsObject,
    disallowExpression?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_type_matcher_v3_regex_pb.RegexMatcher.AsObject,
    disallowIsError?: google_protobuf_wrappers_pb.BoolValue.AsObject,
  }
}

export class HeaderMutation extends jspb.Message {
  hasRemove(): boolean;
  clearRemove(): void;
  getRemove(): string;
  setRemove(value: string): void;

  hasAppend(): boolean;
  clearAppend(): void;
  getAppend(): github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_base_pb.HeaderValueOption | undefined;
  setAppend(value?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_base_pb.HeaderValueOption): void;

  getActionCase(): HeaderMutation.ActionCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HeaderMutation.AsObject;
  static toObject(includeInstance: boolean, msg: HeaderMutation): HeaderMutation.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: HeaderMutation, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HeaderMutation;
  static deserializeBinaryFromReader(message: HeaderMutation, reader: jspb.BinaryReader): HeaderMutation;
}

export namespace HeaderMutation {
  export type AsObject = {
    remove: string,
    append?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_base_pb.HeaderValueOption.AsObject,
  }

  export enum ActionCase {
    ACTION_NOT_SET = 0,
    REMOVE = 1,
    APPEND = 2,
  }
}
