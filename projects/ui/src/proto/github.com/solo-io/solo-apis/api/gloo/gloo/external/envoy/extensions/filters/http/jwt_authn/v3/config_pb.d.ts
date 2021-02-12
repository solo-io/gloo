/* eslint-disable */
// package: solo.io.envoy.extensions.filters.http.jwt_authn.v3
// file: github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/extensions/filters/http/jwt_authn/v3/config.proto

import * as jspb from "google-protobuf";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_base_pb from "../../../../../../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/config/core/v3/base_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_http_uri_pb from "../../../../../../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/config/core/v3/http_uri_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_route_v3_route_components_pb from "../../../../../../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/config/route/v3/route_components_pb";
import * as google_protobuf_duration_pb from "google-protobuf/google/protobuf/duration_pb";
import * as google_protobuf_empty_pb from "google-protobuf/google/protobuf/empty_pb";
import * as udpa_annotations_status_pb from "../../../../../../../../../../../../../udpa/annotations/status_pb";
import * as udpa_annotations_versioning_pb from "../../../../../../../../../../../../../udpa/annotations/versioning_pb";
import * as validate_validate_pb from "../../../../../../../../../../../../../validate/validate_pb";

export class JwtProvider extends jspb.Message {
  getIssuer(): string;
  setIssuer(value: string): void;

  clearAudiencesList(): void;
  getAudiencesList(): Array<string>;
  setAudiencesList(value: Array<string>): void;
  addAudiences(value: string, index?: number): string;

  hasRemoteJwks(): boolean;
  clearRemoteJwks(): void;
  getRemoteJwks(): RemoteJwks | undefined;
  setRemoteJwks(value?: RemoteJwks): void;

  hasLocalJwks(): boolean;
  clearLocalJwks(): void;
  getLocalJwks(): github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_base_pb.DataSource | undefined;
  setLocalJwks(value?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_base_pb.DataSource): void;

  getForward(): boolean;
  setForward(value: boolean): void;

  clearFromHeadersList(): void;
  getFromHeadersList(): Array<JwtHeader>;
  setFromHeadersList(value: Array<JwtHeader>): void;
  addFromHeaders(value?: JwtHeader, index?: number): JwtHeader;

  clearFromParamsList(): void;
  getFromParamsList(): Array<string>;
  setFromParamsList(value: Array<string>): void;
  addFromParams(value: string, index?: number): string;

  getForwardPayloadHeader(): string;
  setForwardPayloadHeader(value: string): void;

  getPayloadInMetadata(): string;
  setPayloadInMetadata(value: string): void;

  getClockSkewSeconds(): number;
  setClockSkewSeconds(value: number): void;

  getJwksSourceSpecifierCase(): JwtProvider.JwksSourceSpecifierCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): JwtProvider.AsObject;
  static toObject(includeInstance: boolean, msg: JwtProvider): JwtProvider.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: JwtProvider, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): JwtProvider;
  static deserializeBinaryFromReader(message: JwtProvider, reader: jspb.BinaryReader): JwtProvider;
}

export namespace JwtProvider {
  export type AsObject = {
    issuer: string,
    audiencesList: Array<string>,
    remoteJwks?: RemoteJwks.AsObject,
    localJwks?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_base_pb.DataSource.AsObject,
    forward: boolean,
    fromHeadersList: Array<JwtHeader.AsObject>,
    fromParamsList: Array<string>,
    forwardPayloadHeader: string,
    payloadInMetadata: string,
    clockSkewSeconds: number,
  }

  export enum JwksSourceSpecifierCase {
    JWKS_SOURCE_SPECIFIER_NOT_SET = 0,
    REMOTE_JWKS = 3,
    LOCAL_JWKS = 4,
  }
}

export class RemoteJwks extends jspb.Message {
  hasHttpUri(): boolean;
  clearHttpUri(): void;
  getHttpUri(): github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_http_uri_pb.HttpUri | undefined;
  setHttpUri(value?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_http_uri_pb.HttpUri): void;

  hasCacheDuration(): boolean;
  clearCacheDuration(): void;
  getCacheDuration(): google_protobuf_duration_pb.Duration | undefined;
  setCacheDuration(value?: google_protobuf_duration_pb.Duration): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoteJwks.AsObject;
  static toObject(includeInstance: boolean, msg: RemoteJwks): RemoteJwks.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RemoteJwks, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoteJwks;
  static deserializeBinaryFromReader(message: RemoteJwks, reader: jspb.BinaryReader): RemoteJwks;
}

export namespace RemoteJwks {
  export type AsObject = {
    httpUri?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_http_uri_pb.HttpUri.AsObject,
    cacheDuration?: google_protobuf_duration_pb.Duration.AsObject,
  }
}

export class JwtHeader extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getValuePrefix(): string;
  setValuePrefix(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): JwtHeader.AsObject;
  static toObject(includeInstance: boolean, msg: JwtHeader): JwtHeader.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: JwtHeader, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): JwtHeader;
  static deserializeBinaryFromReader(message: JwtHeader, reader: jspb.BinaryReader): JwtHeader;
}

export namespace JwtHeader {
  export type AsObject = {
    name: string,
    valuePrefix: string,
  }
}

export class ProviderWithAudiences extends jspb.Message {
  getProviderName(): string;
  setProviderName(value: string): void;

  clearAudiencesList(): void;
  getAudiencesList(): Array<string>;
  setAudiencesList(value: Array<string>): void;
  addAudiences(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProviderWithAudiences.AsObject;
  static toObject(includeInstance: boolean, msg: ProviderWithAudiences): ProviderWithAudiences.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ProviderWithAudiences, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProviderWithAudiences;
  static deserializeBinaryFromReader(message: ProviderWithAudiences, reader: jspb.BinaryReader): ProviderWithAudiences;
}

export namespace ProviderWithAudiences {
  export type AsObject = {
    providerName: string,
    audiencesList: Array<string>,
  }
}

export class JwtRequirement extends jspb.Message {
  hasProviderName(): boolean;
  clearProviderName(): void;
  getProviderName(): string;
  setProviderName(value: string): void;

  hasProviderAndAudiences(): boolean;
  clearProviderAndAudiences(): void;
  getProviderAndAudiences(): ProviderWithAudiences | undefined;
  setProviderAndAudiences(value?: ProviderWithAudiences): void;

  hasRequiresAny(): boolean;
  clearRequiresAny(): void;
  getRequiresAny(): JwtRequirementOrList | undefined;
  setRequiresAny(value?: JwtRequirementOrList): void;

  hasRequiresAll(): boolean;
  clearRequiresAll(): void;
  getRequiresAll(): JwtRequirementAndList | undefined;
  setRequiresAll(value?: JwtRequirementAndList): void;

  hasAllowMissingOrFailed(): boolean;
  clearAllowMissingOrFailed(): void;
  getAllowMissingOrFailed(): google_protobuf_empty_pb.Empty | undefined;
  setAllowMissingOrFailed(value?: google_protobuf_empty_pb.Empty): void;

  hasAllowMissing(): boolean;
  clearAllowMissing(): void;
  getAllowMissing(): google_protobuf_empty_pb.Empty | undefined;
  setAllowMissing(value?: google_protobuf_empty_pb.Empty): void;

  getRequiresTypeCase(): JwtRequirement.RequiresTypeCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): JwtRequirement.AsObject;
  static toObject(includeInstance: boolean, msg: JwtRequirement): JwtRequirement.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: JwtRequirement, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): JwtRequirement;
  static deserializeBinaryFromReader(message: JwtRequirement, reader: jspb.BinaryReader): JwtRequirement;
}

export namespace JwtRequirement {
  export type AsObject = {
    providerName: string,
    providerAndAudiences?: ProviderWithAudiences.AsObject,
    requiresAny?: JwtRequirementOrList.AsObject,
    requiresAll?: JwtRequirementAndList.AsObject,
    allowMissingOrFailed?: google_protobuf_empty_pb.Empty.AsObject,
    allowMissing?: google_protobuf_empty_pb.Empty.AsObject,
  }

  export enum RequiresTypeCase {
    REQUIRES_TYPE_NOT_SET = 0,
    PROVIDER_NAME = 1,
    PROVIDER_AND_AUDIENCES = 2,
    REQUIRES_ANY = 3,
    REQUIRES_ALL = 4,
    ALLOW_MISSING_OR_FAILED = 5,
    ALLOW_MISSING = 6,
  }
}

export class JwtRequirementOrList extends jspb.Message {
  clearRequirementsList(): void;
  getRequirementsList(): Array<JwtRequirement>;
  setRequirementsList(value: Array<JwtRequirement>): void;
  addRequirements(value?: JwtRequirement, index?: number): JwtRequirement;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): JwtRequirementOrList.AsObject;
  static toObject(includeInstance: boolean, msg: JwtRequirementOrList): JwtRequirementOrList.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: JwtRequirementOrList, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): JwtRequirementOrList;
  static deserializeBinaryFromReader(message: JwtRequirementOrList, reader: jspb.BinaryReader): JwtRequirementOrList;
}

export namespace JwtRequirementOrList {
  export type AsObject = {
    requirementsList: Array<JwtRequirement.AsObject>,
  }
}

export class JwtRequirementAndList extends jspb.Message {
  clearRequirementsList(): void;
  getRequirementsList(): Array<JwtRequirement>;
  setRequirementsList(value: Array<JwtRequirement>): void;
  addRequirements(value?: JwtRequirement, index?: number): JwtRequirement;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): JwtRequirementAndList.AsObject;
  static toObject(includeInstance: boolean, msg: JwtRequirementAndList): JwtRequirementAndList.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: JwtRequirementAndList, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): JwtRequirementAndList;
  static deserializeBinaryFromReader(message: JwtRequirementAndList, reader: jspb.BinaryReader): JwtRequirementAndList;
}

export namespace JwtRequirementAndList {
  export type AsObject = {
    requirementsList: Array<JwtRequirement.AsObject>,
  }
}

export class RequirementRule extends jspb.Message {
  hasMatch(): boolean;
  clearMatch(): void;
  getMatch(): github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_route_v3_route_components_pb.RouteMatch | undefined;
  setMatch(value?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_route_v3_route_components_pb.RouteMatch): void;

  hasRequires(): boolean;
  clearRequires(): void;
  getRequires(): JwtRequirement | undefined;
  setRequires(value?: JwtRequirement): void;

  hasRequirementName(): boolean;
  clearRequirementName(): void;
  getRequirementName(): string;
  setRequirementName(value: string): void;

  getRequirementTypeCase(): RequirementRule.RequirementTypeCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RequirementRule.AsObject;
  static toObject(includeInstance: boolean, msg: RequirementRule): RequirementRule.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RequirementRule, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RequirementRule;
  static deserializeBinaryFromReader(message: RequirementRule, reader: jspb.BinaryReader): RequirementRule;
}

export namespace RequirementRule {
  export type AsObject = {
    match?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_route_v3_route_components_pb.RouteMatch.AsObject,
    requires?: JwtRequirement.AsObject,
    requirementName: string,
  }

  export enum RequirementTypeCase {
    REQUIREMENT_TYPE_NOT_SET = 0,
    REQUIRES = 2,
    REQUIREMENT_NAME = 3,
  }
}

export class FilterStateRule extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getRequiresMap(): jspb.Map<string, JwtRequirement>;
  clearRequiresMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FilterStateRule.AsObject;
  static toObject(includeInstance: boolean, msg: FilterStateRule): FilterStateRule.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: FilterStateRule, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FilterStateRule;
  static deserializeBinaryFromReader(message: FilterStateRule, reader: jspb.BinaryReader): FilterStateRule;
}

export namespace FilterStateRule {
  export type AsObject = {
    name: string,
    requiresMap: Array<[string, JwtRequirement.AsObject]>,
  }
}

export class JwtAuthentication extends jspb.Message {
  getProvidersMap(): jspb.Map<string, JwtProvider>;
  clearProvidersMap(): void;
  clearRulesList(): void;
  getRulesList(): Array<RequirementRule>;
  setRulesList(value: Array<RequirementRule>): void;
  addRules(value?: RequirementRule, index?: number): RequirementRule;

  hasFilterStateRules(): boolean;
  clearFilterStateRules(): void;
  getFilterStateRules(): FilterStateRule | undefined;
  setFilterStateRules(value?: FilterStateRule): void;

  getBypassCorsPreflight(): boolean;
  setBypassCorsPreflight(value: boolean): void;

  getRequirementMapMap(): jspb.Map<string, JwtRequirement>;
  clearRequirementMapMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): JwtAuthentication.AsObject;
  static toObject(includeInstance: boolean, msg: JwtAuthentication): JwtAuthentication.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: JwtAuthentication, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): JwtAuthentication;
  static deserializeBinaryFromReader(message: JwtAuthentication, reader: jspb.BinaryReader): JwtAuthentication;
}

export namespace JwtAuthentication {
  export type AsObject = {
    providersMap: Array<[string, JwtProvider.AsObject]>,
    rulesList: Array<RequirementRule.AsObject>,
    filterStateRules?: FilterStateRule.AsObject,
    bypassCorsPreflight: boolean,
    requirementMapMap: Array<[string, JwtRequirement.AsObject]>,
  }
}

export class PerRouteConfig extends jspb.Message {
  hasDisabled(): boolean;
  clearDisabled(): void;
  getDisabled(): boolean;
  setDisabled(value: boolean): void;

  hasRequirementName(): boolean;
  clearRequirementName(): void;
  getRequirementName(): string;
  setRequirementName(value: string): void;

  getRequirementSpecifierCase(): PerRouteConfig.RequirementSpecifierCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PerRouteConfig.AsObject;
  static toObject(includeInstance: boolean, msg: PerRouteConfig): PerRouteConfig.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: PerRouteConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PerRouteConfig;
  static deserializeBinaryFromReader(message: PerRouteConfig, reader: jspb.BinaryReader): PerRouteConfig;
}

export namespace PerRouteConfig {
  export type AsObject = {
    disabled: boolean,
    requirementName: string,
  }

  export enum RequirementSpecifierCase {
    REQUIREMENT_SPECIFIER_NOT_SET = 0,
    DISABLED = 1,
    REQUIREMENT_NAME = 2,
  }
}
