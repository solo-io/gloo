/* eslint-disable */
// package: extproc.options.gloo.solo.io
// file: github.com/solo-io/solo-apis/api/gloo/gloo/v1/enterprise/options/extproc/extproc.proto

import * as jspb from "google-protobuf";
import * as extproto_ext_pb from "../../../../../../../../../../extproto/ext_pb";
import * as validate_validate_pb from "../../../../../../../../../../validate/validate_pb";
import * as github_com_solo_io_solo_kit_api_v1_ref_pb from "../../../../../../../../../../github.com/solo-io/solo-kit/api/v1/ref_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_filters_stages_pb from "../../../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/filters/stages_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_common_mutation_rules_v3_mutation_rules_pb from "../../../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/config/common/mutation_rules/v3/mutation_rules_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_base_pb from "../../../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/config/core/v3/base_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_extensions_filters_http_ext_proc_v3_processing_mode_pb from "../../../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/extensions/filters/http/ext_proc/v3/processing_mode_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_type_matcher_v3_string_pb from "../../../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/type/matcher/v3/string_pb";
import * as google_protobuf_duration_pb from "google-protobuf/google/protobuf/duration_pb";
import * as google_protobuf_struct_pb from "google-protobuf/google/protobuf/struct_pb";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";

export class Settings extends jspb.Message {
  hasGrpcService(): boolean;
  clearGrpcService(): void;
  getGrpcService(): GrpcService | undefined;
  setGrpcService(value?: GrpcService): void;

  hasFilterStage(): boolean;
  clearFilterStage(): void;
  getFilterStage(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_filters_stages_pb.FilterStage | undefined;
  setFilterStage(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_filters_stages_pb.FilterStage): void;

  hasFailureModeAllow(): boolean;
  clearFailureModeAllow(): void;
  getFailureModeAllow(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setFailureModeAllow(value?: google_protobuf_wrappers_pb.BoolValue): void;

  hasProcessingMode(): boolean;
  clearProcessingMode(): void;
  getProcessingMode(): github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_extensions_filters_http_ext_proc_v3_processing_mode_pb.ProcessingMode | undefined;
  setProcessingMode(value?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_extensions_filters_http_ext_proc_v3_processing_mode_pb.ProcessingMode): void;

  hasAsyncMode(): boolean;
  clearAsyncMode(): void;
  getAsyncMode(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setAsyncMode(value?: google_protobuf_wrappers_pb.BoolValue): void;

  clearRequestAttributesList(): void;
  getRequestAttributesList(): Array<string>;
  setRequestAttributesList(value: Array<string>): void;
  addRequestAttributes(value: string, index?: number): string;

  clearResponseAttributesList(): void;
  getResponseAttributesList(): Array<string>;
  setResponseAttributesList(value: Array<string>): void;
  addResponseAttributes(value: string, index?: number): string;

  hasMessageTimeout(): boolean;
  clearMessageTimeout(): void;
  getMessageTimeout(): google_protobuf_duration_pb.Duration | undefined;
  setMessageTimeout(value?: google_protobuf_duration_pb.Duration): void;

  hasStatPrefix(): boolean;
  clearStatPrefix(): void;
  getStatPrefix(): google_protobuf_wrappers_pb.StringValue | undefined;
  setStatPrefix(value?: google_protobuf_wrappers_pb.StringValue): void;

  hasMutationRules(): boolean;
  clearMutationRules(): void;
  getMutationRules(): github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_common_mutation_rules_v3_mutation_rules_pb.HeaderMutationRules | undefined;
  setMutationRules(value?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_common_mutation_rules_v3_mutation_rules_pb.HeaderMutationRules): void;

  hasMaxMessageTimeout(): boolean;
  clearMaxMessageTimeout(): void;
  getMaxMessageTimeout(): google_protobuf_duration_pb.Duration | undefined;
  setMaxMessageTimeout(value?: google_protobuf_duration_pb.Duration): void;

  hasDisableClearRouteCache(): boolean;
  clearDisableClearRouteCache(): void;
  getDisableClearRouteCache(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setDisableClearRouteCache(value?: google_protobuf_wrappers_pb.BoolValue): void;

  hasForwardRules(): boolean;
  clearForwardRules(): void;
  getForwardRules(): HeaderForwardingRules | undefined;
  setForwardRules(value?: HeaderForwardingRules): void;

  hasFilterMetadata(): boolean;
  clearFilterMetadata(): void;
  getFilterMetadata(): google_protobuf_struct_pb.Struct | undefined;
  setFilterMetadata(value?: google_protobuf_struct_pb.Struct): void;

  hasAllowModeOverride(): boolean;
  clearAllowModeOverride(): void;
  getAllowModeOverride(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setAllowModeOverride(value?: google_protobuf_wrappers_pb.BoolValue): void;

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
    grpcService?: GrpcService.AsObject,
    filterStage?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_filters_stages_pb.FilterStage.AsObject,
    failureModeAllow?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    processingMode?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_extensions_filters_http_ext_proc_v3_processing_mode_pb.ProcessingMode.AsObject,
    asyncMode?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    requestAttributesList: Array<string>,
    responseAttributesList: Array<string>,
    messageTimeout?: google_protobuf_duration_pb.Duration.AsObject,
    statPrefix?: google_protobuf_wrappers_pb.StringValue.AsObject,
    mutationRules?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_common_mutation_rules_v3_mutation_rules_pb.HeaderMutationRules.AsObject,
    maxMessageTimeout?: google_protobuf_duration_pb.Duration.AsObject,
    disableClearRouteCache?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    forwardRules?: HeaderForwardingRules.AsObject,
    filterMetadata?: google_protobuf_struct_pb.Struct.AsObject,
    allowModeOverride?: google_protobuf_wrappers_pb.BoolValue.AsObject,
  }
}

export class RouteSettings extends jspb.Message {
  hasDisabled(): boolean;
  clearDisabled(): void;
  getDisabled(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setDisabled(value?: google_protobuf_wrappers_pb.BoolValue): void;

  hasOverrides(): boolean;
  clearOverrides(): void;
  getOverrides(): Overrides | undefined;
  setOverrides(value?: Overrides): void;

  getOverrideCase(): RouteSettings.OverrideCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RouteSettings.AsObject;
  static toObject(includeInstance: boolean, msg: RouteSettings): RouteSettings.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RouteSettings, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RouteSettings;
  static deserializeBinaryFromReader(message: RouteSettings, reader: jspb.BinaryReader): RouteSettings;
}

export namespace RouteSettings {
  export type AsObject = {
    disabled?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    overrides?: Overrides.AsObject,
  }

  export enum OverrideCase {
    OVERRIDE_NOT_SET = 0,
    DISABLED = 1,
    OVERRIDES = 2,
  }
}

export class GrpcService extends jspb.Message {
  hasExtProcServerRef(): boolean;
  clearExtProcServerRef(): void;
  getExtProcServerRef(): github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef | undefined;
  setExtProcServerRef(value?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef): void;

  hasAuthority(): boolean;
  clearAuthority(): void;
  getAuthority(): google_protobuf_wrappers_pb.StringValue | undefined;
  setAuthority(value?: google_protobuf_wrappers_pb.StringValue): void;

  hasRetryPolicy(): boolean;
  clearRetryPolicy(): void;
  getRetryPolicy(): github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_base_pb.RetryPolicy | undefined;
  setRetryPolicy(value?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_base_pb.RetryPolicy): void;

  hasTimeout(): boolean;
  clearTimeout(): void;
  getTimeout(): google_protobuf_duration_pb.Duration | undefined;
  setTimeout(value?: google_protobuf_duration_pb.Duration): void;

  clearInitialMetadataList(): void;
  getInitialMetadataList(): Array<github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_base_pb.HeaderValue>;
  setInitialMetadataList(value: Array<github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_base_pb.HeaderValue>): void;
  addInitialMetadata(value?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_base_pb.HeaderValue, index?: number): github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_base_pb.HeaderValue;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GrpcService.AsObject;
  static toObject(includeInstance: boolean, msg: GrpcService): GrpcService.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GrpcService, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GrpcService;
  static deserializeBinaryFromReader(message: GrpcService, reader: jspb.BinaryReader): GrpcService;
}

export namespace GrpcService {
  export type AsObject = {
    extProcServerRef?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject,
    authority?: google_protobuf_wrappers_pb.StringValue.AsObject,
    retryPolicy?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_base_pb.RetryPolicy.AsObject,
    timeout?: google_protobuf_duration_pb.Duration.AsObject,
    initialMetadataList: Array<github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_base_pb.HeaderValue.AsObject>,
  }
}

export class Overrides extends jspb.Message {
  hasProcessingMode(): boolean;
  clearProcessingMode(): void;
  getProcessingMode(): github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_extensions_filters_http_ext_proc_v3_processing_mode_pb.ProcessingMode | undefined;
  setProcessingMode(value?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_extensions_filters_http_ext_proc_v3_processing_mode_pb.ProcessingMode): void;

  hasAsyncMode(): boolean;
  clearAsyncMode(): void;
  getAsyncMode(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setAsyncMode(value?: google_protobuf_wrappers_pb.BoolValue): void;

  clearRequestAttributesList(): void;
  getRequestAttributesList(): Array<string>;
  setRequestAttributesList(value: Array<string>): void;
  addRequestAttributes(value: string, index?: number): string;

  clearResponseAttributesList(): void;
  getResponseAttributesList(): Array<string>;
  setResponseAttributesList(value: Array<string>): void;
  addResponseAttributes(value: string, index?: number): string;

  hasGrpcService(): boolean;
  clearGrpcService(): void;
  getGrpcService(): GrpcService | undefined;
  setGrpcService(value?: GrpcService): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Overrides.AsObject;
  static toObject(includeInstance: boolean, msg: Overrides): Overrides.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Overrides, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Overrides;
  static deserializeBinaryFromReader(message: Overrides, reader: jspb.BinaryReader): Overrides;
}

export namespace Overrides {
  export type AsObject = {
    processingMode?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_extensions_filters_http_ext_proc_v3_processing_mode_pb.ProcessingMode.AsObject,
    asyncMode?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    requestAttributesList: Array<string>,
    responseAttributesList: Array<string>,
    grpcService?: GrpcService.AsObject,
  }
}

export class HeaderForwardingRules extends jspb.Message {
  hasAllowedHeaders(): boolean;
  clearAllowedHeaders(): void;
  getAllowedHeaders(): github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_type_matcher_v3_string_pb.ListStringMatcher | undefined;
  setAllowedHeaders(value?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_type_matcher_v3_string_pb.ListStringMatcher): void;

  hasDisallowedHeaders(): boolean;
  clearDisallowedHeaders(): void;
  getDisallowedHeaders(): github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_type_matcher_v3_string_pb.ListStringMatcher | undefined;
  setDisallowedHeaders(value?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_type_matcher_v3_string_pb.ListStringMatcher): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HeaderForwardingRules.AsObject;
  static toObject(includeInstance: boolean, msg: HeaderForwardingRules): HeaderForwardingRules.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: HeaderForwardingRules, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HeaderForwardingRules;
  static deserializeBinaryFromReader(message: HeaderForwardingRules, reader: jspb.BinaryReader): HeaderForwardingRules;
}

export namespace HeaderForwardingRules {
  export type AsObject = {
    allowedHeaders?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_type_matcher_v3_string_pb.ListStringMatcher.AsObject,
    disallowedHeaders?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_type_matcher_v3_string_pb.ListStringMatcher.AsObject,
  }
}
