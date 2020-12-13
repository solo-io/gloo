/* eslint-disable */
// package: solo.io.envoy.config.route.v3
// file: github.com/solo-io/gloo/projects/gloo/api/external/envoy/config/route/v3/route_components.proto

import * as jspb from "google-protobuf";
import * as github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb from "../../../../../../../../../../../github.com/solo-io/gloo/projects/gloo/api/external/envoy/config/core/v3/base_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_extension_pb from "../../../../../../../../../../../github.com/solo-io/gloo/projects/gloo/api/external/envoy/config/core/v3/extension_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_proxy_protocol_pb from "../../../../../../../../../../../github.com/solo-io/gloo/projects/gloo/api/external/envoy/config/core/v3/proxy_protocol_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_matcher_v3_regex_pb from "../../../../../../../../../../../github.com/solo-io/gloo/projects/gloo/api/external/envoy/type/matcher/v3/regex_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_matcher_v3_string_pb from "../../../../../../../../../../../github.com/solo-io/gloo/projects/gloo/api/external/envoy/type/matcher/v3/string_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_metadata_v3_metadata_pb from "../../../../../../../../../../../github.com/solo-io/gloo/projects/gloo/api/external/envoy/type/metadata/v3/metadata_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_tracing_v3_custom_tag_pb from "../../../../../../../../../../../github.com/solo-io/gloo/projects/gloo/api/external/envoy/type/tracing/v3/custom_tag_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_v3_percent_pb from "../../../../../../../../../../../github.com/solo-io/gloo/projects/gloo/api/external/envoy/type/v3/percent_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_v3_range_pb from "../../../../../../../../../../../github.com/solo-io/gloo/projects/gloo/api/external/envoy/type/v3/range_pb";
import * as google_protobuf_any_pb from "google-protobuf/google/protobuf/any_pb";
import * as google_protobuf_duration_pb from "google-protobuf/google/protobuf/duration_pb";
import * as google_protobuf_struct_pb from "google-protobuf/google/protobuf/struct_pb";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";
import * as envoy_annotations_deprecation_pb from "../../../../../../../../../../../envoy/annotations/deprecation_pb";
import * as udpa_annotations_migrate_pb from "../../../../../../../../../../../udpa/annotations/migrate_pb";
import * as udpa_annotations_status_pb from "../../../../../../../../../../../udpa/annotations/status_pb";
import * as udpa_annotations_versioning_pb from "../../../../../../../../../../../udpa/annotations/versioning_pb";
import * as validate_validate_pb from "../../../../../../../../../../../validate/validate_pb";
import * as extproto_ext_pb from "../../../../../../../../../../../protoc-gen-ext/extproto/ext_pb";

export class VirtualHost extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  clearDomainsList(): void;
  getDomainsList(): Array<string>;
  setDomainsList(value: Array<string>): void;
  addDomains(value: string, index?: number): string;

  clearRoutesList(): void;
  getRoutesList(): Array<Route>;
  setRoutesList(value: Array<Route>): void;
  addRoutes(value?: Route, index?: number): Route;

  getRequireTls(): VirtualHost.TlsRequirementTypeMap[keyof VirtualHost.TlsRequirementTypeMap];
  setRequireTls(value: VirtualHost.TlsRequirementTypeMap[keyof VirtualHost.TlsRequirementTypeMap]): void;

  clearVirtualClustersList(): void;
  getVirtualClustersList(): Array<VirtualCluster>;
  setVirtualClustersList(value: Array<VirtualCluster>): void;
  addVirtualClusters(value?: VirtualCluster, index?: number): VirtualCluster;

  clearRateLimitsList(): void;
  getRateLimitsList(): Array<RateLimit>;
  setRateLimitsList(value: Array<RateLimit>): void;
  addRateLimits(value?: RateLimit, index?: number): RateLimit;

  clearRequestHeadersToAddList(): void;
  getRequestHeadersToAddList(): Array<github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb.HeaderValueOption>;
  setRequestHeadersToAddList(value: Array<github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb.HeaderValueOption>): void;
  addRequestHeadersToAdd(value?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb.HeaderValueOption, index?: number): github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb.HeaderValueOption;

  clearRequestHeadersToRemoveList(): void;
  getRequestHeadersToRemoveList(): Array<string>;
  setRequestHeadersToRemoveList(value: Array<string>): void;
  addRequestHeadersToRemove(value: string, index?: number): string;

  clearResponseHeadersToAddList(): void;
  getResponseHeadersToAddList(): Array<github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb.HeaderValueOption>;
  setResponseHeadersToAddList(value: Array<github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb.HeaderValueOption>): void;
  addResponseHeadersToAdd(value?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb.HeaderValueOption, index?: number): github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb.HeaderValueOption;

  clearResponseHeadersToRemoveList(): void;
  getResponseHeadersToRemoveList(): Array<string>;
  setResponseHeadersToRemoveList(value: Array<string>): void;
  addResponseHeadersToRemove(value: string, index?: number): string;

  hasCors(): boolean;
  clearCors(): void;
  getCors(): CorsPolicy | undefined;
  setCors(value?: CorsPolicy): void;

  getTypedPerFilterConfigMap(): jspb.Map<string, google_protobuf_any_pb.Any>;
  clearTypedPerFilterConfigMap(): void;
  getIncludeRequestAttemptCount(): boolean;
  setIncludeRequestAttemptCount(value: boolean): void;

  getIncludeAttemptCountInResponse(): boolean;
  setIncludeAttemptCountInResponse(value: boolean): void;

  hasRetryPolicy(): boolean;
  clearRetryPolicy(): void;
  getRetryPolicy(): RetryPolicy | undefined;
  setRetryPolicy(value?: RetryPolicy): void;

  hasRetryPolicyTypedConfig(): boolean;
  clearRetryPolicyTypedConfig(): void;
  getRetryPolicyTypedConfig(): google_protobuf_any_pb.Any | undefined;
  setRetryPolicyTypedConfig(value?: google_protobuf_any_pb.Any): void;

  hasHedgePolicy(): boolean;
  clearHedgePolicy(): void;
  getHedgePolicy(): HedgePolicy | undefined;
  setHedgePolicy(value?: HedgePolicy): void;

  hasPerRequestBufferLimitBytes(): boolean;
  clearPerRequestBufferLimitBytes(): void;
  getPerRequestBufferLimitBytes(): google_protobuf_wrappers_pb.UInt32Value | undefined;
  setPerRequestBufferLimitBytes(value?: google_protobuf_wrappers_pb.UInt32Value): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VirtualHost.AsObject;
  static toObject(includeInstance: boolean, msg: VirtualHost): VirtualHost.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: VirtualHost, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VirtualHost;
  static deserializeBinaryFromReader(message: VirtualHost, reader: jspb.BinaryReader): VirtualHost;
}

export namespace VirtualHost {
  export type AsObject = {
    name: string,
    domainsList: Array<string>,
    routesList: Array<Route.AsObject>,
    requireTls: VirtualHost.TlsRequirementTypeMap[keyof VirtualHost.TlsRequirementTypeMap],
    virtualClustersList: Array<VirtualCluster.AsObject>,
    rateLimitsList: Array<RateLimit.AsObject>,
    requestHeadersToAddList: Array<github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb.HeaderValueOption.AsObject>,
    requestHeadersToRemoveList: Array<string>,
    responseHeadersToAddList: Array<github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb.HeaderValueOption.AsObject>,
    responseHeadersToRemoveList: Array<string>,
    cors?: CorsPolicy.AsObject,
    typedPerFilterConfigMap: Array<[string, google_protobuf_any_pb.Any.AsObject]>,
    includeRequestAttemptCount: boolean,
    includeAttemptCountInResponse: boolean,
    retryPolicy?: RetryPolicy.AsObject,
    retryPolicyTypedConfig?: google_protobuf_any_pb.Any.AsObject,
    hedgePolicy?: HedgePolicy.AsObject,
    perRequestBufferLimitBytes?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
  }

  export interface TlsRequirementTypeMap {
    NONE: 0;
    EXTERNAL_ONLY: 1;
    ALL: 2;
  }

  export const TlsRequirementType: TlsRequirementTypeMap;
}

export class FilterAction extends jspb.Message {
  hasAction(): boolean;
  clearAction(): void;
  getAction(): google_protobuf_any_pb.Any | undefined;
  setAction(value?: google_protobuf_any_pb.Any): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FilterAction.AsObject;
  static toObject(includeInstance: boolean, msg: FilterAction): FilterAction.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: FilterAction, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FilterAction;
  static deserializeBinaryFromReader(message: FilterAction, reader: jspb.BinaryReader): FilterAction;
}

export namespace FilterAction {
  export type AsObject = {
    action?: google_protobuf_any_pb.Any.AsObject,
  }
}

export class Route extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  hasMatch(): boolean;
  clearMatch(): void;
  getMatch(): RouteMatch | undefined;
  setMatch(value?: RouteMatch): void;

  hasRoute(): boolean;
  clearRoute(): void;
  getRoute(): RouteAction | undefined;
  setRoute(value?: RouteAction): void;

  hasRedirect(): boolean;
  clearRedirect(): void;
  getRedirect(): RedirectAction | undefined;
  setRedirect(value?: RedirectAction): void;

  hasDirectResponse(): boolean;
  clearDirectResponse(): void;
  getDirectResponse(): DirectResponseAction | undefined;
  setDirectResponse(value?: DirectResponseAction): void;

  hasFilterAction(): boolean;
  clearFilterAction(): void;
  getFilterAction(): FilterAction | undefined;
  setFilterAction(value?: FilterAction): void;

  hasMetadata(): boolean;
  clearMetadata(): void;
  getMetadata(): github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb.Metadata | undefined;
  setMetadata(value?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb.Metadata): void;

  hasDecorator(): boolean;
  clearDecorator(): void;
  getDecorator(): Decorator | undefined;
  setDecorator(value?: Decorator): void;

  getTypedPerFilterConfigMap(): jspb.Map<string, google_protobuf_any_pb.Any>;
  clearTypedPerFilterConfigMap(): void;
  clearRequestHeadersToAddList(): void;
  getRequestHeadersToAddList(): Array<github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb.HeaderValueOption>;
  setRequestHeadersToAddList(value: Array<github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb.HeaderValueOption>): void;
  addRequestHeadersToAdd(value?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb.HeaderValueOption, index?: number): github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb.HeaderValueOption;

  clearRequestHeadersToRemoveList(): void;
  getRequestHeadersToRemoveList(): Array<string>;
  setRequestHeadersToRemoveList(value: Array<string>): void;
  addRequestHeadersToRemove(value: string, index?: number): string;

  clearResponseHeadersToAddList(): void;
  getResponseHeadersToAddList(): Array<github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb.HeaderValueOption>;
  setResponseHeadersToAddList(value: Array<github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb.HeaderValueOption>): void;
  addResponseHeadersToAdd(value?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb.HeaderValueOption, index?: number): github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb.HeaderValueOption;

  clearResponseHeadersToRemoveList(): void;
  getResponseHeadersToRemoveList(): Array<string>;
  setResponseHeadersToRemoveList(value: Array<string>): void;
  addResponseHeadersToRemove(value: string, index?: number): string;

  hasTracing(): boolean;
  clearTracing(): void;
  getTracing(): Tracing | undefined;
  setTracing(value?: Tracing): void;

  hasPerRequestBufferLimitBytes(): boolean;
  clearPerRequestBufferLimitBytes(): void;
  getPerRequestBufferLimitBytes(): google_protobuf_wrappers_pb.UInt32Value | undefined;
  setPerRequestBufferLimitBytes(value?: google_protobuf_wrappers_pb.UInt32Value): void;

  getActionCase(): Route.ActionCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Route.AsObject;
  static toObject(includeInstance: boolean, msg: Route): Route.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Route, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Route;
  static deserializeBinaryFromReader(message: Route, reader: jspb.BinaryReader): Route;
}

export namespace Route {
  export type AsObject = {
    name: string,
    match?: RouteMatch.AsObject,
    route?: RouteAction.AsObject,
    redirect?: RedirectAction.AsObject,
    directResponse?: DirectResponseAction.AsObject,
    filterAction?: FilterAction.AsObject,
    metadata?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb.Metadata.AsObject,
    decorator?: Decorator.AsObject,
    typedPerFilterConfigMap: Array<[string, google_protobuf_any_pb.Any.AsObject]>,
    requestHeadersToAddList: Array<github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb.HeaderValueOption.AsObject>,
    requestHeadersToRemoveList: Array<string>,
    responseHeadersToAddList: Array<github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb.HeaderValueOption.AsObject>,
    responseHeadersToRemoveList: Array<string>,
    tracing?: Tracing.AsObject,
    perRequestBufferLimitBytes?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
  }

  export enum ActionCase {
    ACTION_NOT_SET = 0,
    ROUTE = 2,
    REDIRECT = 3,
    DIRECT_RESPONSE = 7,
    FILTER_ACTION = 17,
  }
}

export class WeightedCluster extends jspb.Message {
  clearClustersList(): void;
  getClustersList(): Array<WeightedCluster.ClusterWeight>;
  setClustersList(value: Array<WeightedCluster.ClusterWeight>): void;
  addClusters(value?: WeightedCluster.ClusterWeight, index?: number): WeightedCluster.ClusterWeight;

  hasTotalWeight(): boolean;
  clearTotalWeight(): void;
  getTotalWeight(): google_protobuf_wrappers_pb.UInt32Value | undefined;
  setTotalWeight(value?: google_protobuf_wrappers_pb.UInt32Value): void;

  getRuntimeKeyPrefix(): string;
  setRuntimeKeyPrefix(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): WeightedCluster.AsObject;
  static toObject(includeInstance: boolean, msg: WeightedCluster): WeightedCluster.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: WeightedCluster, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): WeightedCluster;
  static deserializeBinaryFromReader(message: WeightedCluster, reader: jspb.BinaryReader): WeightedCluster;
}

export namespace WeightedCluster {
  export type AsObject = {
    clustersList: Array<WeightedCluster.ClusterWeight.AsObject>,
    totalWeight?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
    runtimeKeyPrefix: string,
  }

  export class ClusterWeight extends jspb.Message {
    getName(): string;
    setName(value: string): void;

    hasWeight(): boolean;
    clearWeight(): void;
    getWeight(): google_protobuf_wrappers_pb.UInt32Value | undefined;
    setWeight(value?: google_protobuf_wrappers_pb.UInt32Value): void;

    hasMetadataMatch(): boolean;
    clearMetadataMatch(): void;
    getMetadataMatch(): github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb.Metadata | undefined;
    setMetadataMatch(value?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb.Metadata): void;

    clearRequestHeadersToAddList(): void;
    getRequestHeadersToAddList(): Array<github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb.HeaderValueOption>;
    setRequestHeadersToAddList(value: Array<github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb.HeaderValueOption>): void;
    addRequestHeadersToAdd(value?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb.HeaderValueOption, index?: number): github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb.HeaderValueOption;

    clearRequestHeadersToRemoveList(): void;
    getRequestHeadersToRemoveList(): Array<string>;
    setRequestHeadersToRemoveList(value: Array<string>): void;
    addRequestHeadersToRemove(value: string, index?: number): string;

    clearResponseHeadersToAddList(): void;
    getResponseHeadersToAddList(): Array<github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb.HeaderValueOption>;
    setResponseHeadersToAddList(value: Array<github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb.HeaderValueOption>): void;
    addResponseHeadersToAdd(value?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb.HeaderValueOption, index?: number): github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb.HeaderValueOption;

    clearResponseHeadersToRemoveList(): void;
    getResponseHeadersToRemoveList(): Array<string>;
    setResponseHeadersToRemoveList(value: Array<string>): void;
    addResponseHeadersToRemove(value: string, index?: number): string;

    getTypedPerFilterConfigMap(): jspb.Map<string, google_protobuf_any_pb.Any>;
    clearTypedPerFilterConfigMap(): void;
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ClusterWeight.AsObject;
    static toObject(includeInstance: boolean, msg: ClusterWeight): ClusterWeight.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: ClusterWeight, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ClusterWeight;
    static deserializeBinaryFromReader(message: ClusterWeight, reader: jspb.BinaryReader): ClusterWeight;
  }

  export namespace ClusterWeight {
    export type AsObject = {
      name: string,
      weight?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
      metadataMatch?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb.Metadata.AsObject,
      requestHeadersToAddList: Array<github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb.HeaderValueOption.AsObject>,
      requestHeadersToRemoveList: Array<string>,
      responseHeadersToAddList: Array<github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb.HeaderValueOption.AsObject>,
      responseHeadersToRemoveList: Array<string>,
      typedPerFilterConfigMap: Array<[string, google_protobuf_any_pb.Any.AsObject]>,
    }
  }
}

export class RouteMatch extends jspb.Message {
  hasPrefix(): boolean;
  clearPrefix(): void;
  getPrefix(): string;
  setPrefix(value: string): void;

  hasPath(): boolean;
  clearPath(): void;
  getPath(): string;
  setPath(value: string): void;

  hasSafeRegex(): boolean;
  clearSafeRegex(): void;
  getSafeRegex(): github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_matcher_v3_regex_pb.RegexMatcher | undefined;
  setSafeRegex(value?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_matcher_v3_regex_pb.RegexMatcher): void;

  hasConnectMatcher(): boolean;
  clearConnectMatcher(): void;
  getConnectMatcher(): RouteMatch.ConnectMatcher | undefined;
  setConnectMatcher(value?: RouteMatch.ConnectMatcher): void;

  hasCaseSensitive(): boolean;
  clearCaseSensitive(): void;
  getCaseSensitive(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setCaseSensitive(value?: google_protobuf_wrappers_pb.BoolValue): void;

  hasRuntimeFraction(): boolean;
  clearRuntimeFraction(): void;
  getRuntimeFraction(): github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb.RuntimeFractionalPercent | undefined;
  setRuntimeFraction(value?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb.RuntimeFractionalPercent): void;

  clearHeadersList(): void;
  getHeadersList(): Array<HeaderMatcher>;
  setHeadersList(value: Array<HeaderMatcher>): void;
  addHeaders(value?: HeaderMatcher, index?: number): HeaderMatcher;

  clearQueryParametersList(): void;
  getQueryParametersList(): Array<QueryParameterMatcher>;
  setQueryParametersList(value: Array<QueryParameterMatcher>): void;
  addQueryParameters(value?: QueryParameterMatcher, index?: number): QueryParameterMatcher;

  hasGrpc(): boolean;
  clearGrpc(): void;
  getGrpc(): RouteMatch.GrpcRouteMatchOptions | undefined;
  setGrpc(value?: RouteMatch.GrpcRouteMatchOptions): void;

  hasTlsContext(): boolean;
  clearTlsContext(): void;
  getTlsContext(): RouteMatch.TlsContextMatchOptions | undefined;
  setTlsContext(value?: RouteMatch.TlsContextMatchOptions): void;

  getPathSpecifierCase(): RouteMatch.PathSpecifierCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RouteMatch.AsObject;
  static toObject(includeInstance: boolean, msg: RouteMatch): RouteMatch.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RouteMatch, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RouteMatch;
  static deserializeBinaryFromReader(message: RouteMatch, reader: jspb.BinaryReader): RouteMatch;
}

export namespace RouteMatch {
  export type AsObject = {
    prefix: string,
    path: string,
    safeRegex?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_matcher_v3_regex_pb.RegexMatcher.AsObject,
    connectMatcher?: RouteMatch.ConnectMatcher.AsObject,
    caseSensitive?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    runtimeFraction?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb.RuntimeFractionalPercent.AsObject,
    headersList: Array<HeaderMatcher.AsObject>,
    queryParametersList: Array<QueryParameterMatcher.AsObject>,
    grpc?: RouteMatch.GrpcRouteMatchOptions.AsObject,
    tlsContext?: RouteMatch.TlsContextMatchOptions.AsObject,
  }

  export class GrpcRouteMatchOptions extends jspb.Message {
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): GrpcRouteMatchOptions.AsObject;
    static toObject(includeInstance: boolean, msg: GrpcRouteMatchOptions): GrpcRouteMatchOptions.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: GrpcRouteMatchOptions, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): GrpcRouteMatchOptions;
    static deserializeBinaryFromReader(message: GrpcRouteMatchOptions, reader: jspb.BinaryReader): GrpcRouteMatchOptions;
  }

  export namespace GrpcRouteMatchOptions {
    export type AsObject = {
    }
  }

  export class TlsContextMatchOptions extends jspb.Message {
    hasPresented(): boolean;
    clearPresented(): void;
    getPresented(): google_protobuf_wrappers_pb.BoolValue | undefined;
    setPresented(value?: google_protobuf_wrappers_pb.BoolValue): void;

    hasValidated(): boolean;
    clearValidated(): void;
    getValidated(): google_protobuf_wrappers_pb.BoolValue | undefined;
    setValidated(value?: google_protobuf_wrappers_pb.BoolValue): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): TlsContextMatchOptions.AsObject;
    static toObject(includeInstance: boolean, msg: TlsContextMatchOptions): TlsContextMatchOptions.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: TlsContextMatchOptions, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): TlsContextMatchOptions;
    static deserializeBinaryFromReader(message: TlsContextMatchOptions, reader: jspb.BinaryReader): TlsContextMatchOptions;
  }

  export namespace TlsContextMatchOptions {
    export type AsObject = {
      presented?: google_protobuf_wrappers_pb.BoolValue.AsObject,
      validated?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    }
  }

  export class ConnectMatcher extends jspb.Message {
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ConnectMatcher.AsObject;
    static toObject(includeInstance: boolean, msg: ConnectMatcher): ConnectMatcher.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: ConnectMatcher, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ConnectMatcher;
    static deserializeBinaryFromReader(message: ConnectMatcher, reader: jspb.BinaryReader): ConnectMatcher;
  }

  export namespace ConnectMatcher {
    export type AsObject = {
    }
  }

  export enum PathSpecifierCase {
    PATH_SPECIFIER_NOT_SET = 0,
    PREFIX = 1,
    PATH = 2,
    SAFE_REGEX = 10,
    CONNECT_MATCHER = 12,
  }
}

export class CorsPolicy extends jspb.Message {
  clearAllowOriginStringMatchList(): void;
  getAllowOriginStringMatchList(): Array<github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_matcher_v3_string_pb.StringMatcher>;
  setAllowOriginStringMatchList(value: Array<github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_matcher_v3_string_pb.StringMatcher>): void;
  addAllowOriginStringMatch(value?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_matcher_v3_string_pb.StringMatcher, index?: number): github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_matcher_v3_string_pb.StringMatcher;

  getAllowMethods(): string;
  setAllowMethods(value: string): void;

  getAllowHeaders(): string;
  setAllowHeaders(value: string): void;

  getExposeHeaders(): string;
  setExposeHeaders(value: string): void;

  getMaxAge(): string;
  setMaxAge(value: string): void;

  hasAllowCredentials(): boolean;
  clearAllowCredentials(): void;
  getAllowCredentials(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setAllowCredentials(value?: google_protobuf_wrappers_pb.BoolValue): void;

  hasFilterEnabled(): boolean;
  clearFilterEnabled(): void;
  getFilterEnabled(): github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb.RuntimeFractionalPercent | undefined;
  setFilterEnabled(value?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb.RuntimeFractionalPercent): void;

  hasShadowEnabled(): boolean;
  clearShadowEnabled(): void;
  getShadowEnabled(): github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb.RuntimeFractionalPercent | undefined;
  setShadowEnabled(value?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb.RuntimeFractionalPercent): void;

  getEnabledSpecifierCase(): CorsPolicy.EnabledSpecifierCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CorsPolicy.AsObject;
  static toObject(includeInstance: boolean, msg: CorsPolicy): CorsPolicy.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CorsPolicy, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CorsPolicy;
  static deserializeBinaryFromReader(message: CorsPolicy, reader: jspb.BinaryReader): CorsPolicy;
}

export namespace CorsPolicy {
  export type AsObject = {
    allowOriginStringMatchList: Array<github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_matcher_v3_string_pb.StringMatcher.AsObject>,
    allowMethods: string,
    allowHeaders: string,
    exposeHeaders: string,
    maxAge: string,
    allowCredentials?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    filterEnabled?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb.RuntimeFractionalPercent.AsObject,
    shadowEnabled?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb.RuntimeFractionalPercent.AsObject,
  }

  export enum EnabledSpecifierCase {
    ENABLED_SPECIFIER_NOT_SET = 0,
    FILTER_ENABLED = 9,
  }
}

export class RouteAction extends jspb.Message {
  hasCluster(): boolean;
  clearCluster(): void;
  getCluster(): string;
  setCluster(value: string): void;

  hasClusterHeader(): boolean;
  clearClusterHeader(): void;
  getClusterHeader(): string;
  setClusterHeader(value: string): void;

  hasWeightedClusters(): boolean;
  clearWeightedClusters(): void;
  getWeightedClusters(): WeightedCluster | undefined;
  setWeightedClusters(value?: WeightedCluster): void;

  getClusterNotFoundResponseCode(): RouteAction.ClusterNotFoundResponseCodeMap[keyof RouteAction.ClusterNotFoundResponseCodeMap];
  setClusterNotFoundResponseCode(value: RouteAction.ClusterNotFoundResponseCodeMap[keyof RouteAction.ClusterNotFoundResponseCodeMap]): void;

  hasMetadataMatch(): boolean;
  clearMetadataMatch(): void;
  getMetadataMatch(): github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb.Metadata | undefined;
  setMetadataMatch(value?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb.Metadata): void;

  getPrefixRewrite(): string;
  setPrefixRewrite(value: string): void;

  hasRegexRewrite(): boolean;
  clearRegexRewrite(): void;
  getRegexRewrite(): github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_matcher_v3_regex_pb.RegexMatchAndSubstitute | undefined;
  setRegexRewrite(value?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_matcher_v3_regex_pb.RegexMatchAndSubstitute): void;

  hasHostRewriteLiteral(): boolean;
  clearHostRewriteLiteral(): void;
  getHostRewriteLiteral(): string;
  setHostRewriteLiteral(value: string): void;

  hasAutoHostRewrite(): boolean;
  clearAutoHostRewrite(): void;
  getAutoHostRewrite(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setAutoHostRewrite(value?: google_protobuf_wrappers_pb.BoolValue): void;

  hasHostRewriteHeader(): boolean;
  clearHostRewriteHeader(): void;
  getHostRewriteHeader(): string;
  setHostRewriteHeader(value: string): void;

  hasTimeout(): boolean;
  clearTimeout(): void;
  getTimeout(): google_protobuf_duration_pb.Duration | undefined;
  setTimeout(value?: google_protobuf_duration_pb.Duration): void;

  hasIdleTimeout(): boolean;
  clearIdleTimeout(): void;
  getIdleTimeout(): google_protobuf_duration_pb.Duration | undefined;
  setIdleTimeout(value?: google_protobuf_duration_pb.Duration): void;

  hasRetryPolicy(): boolean;
  clearRetryPolicy(): void;
  getRetryPolicy(): RetryPolicy | undefined;
  setRetryPolicy(value?: RetryPolicy): void;

  hasRetryPolicyTypedConfig(): boolean;
  clearRetryPolicyTypedConfig(): void;
  getRetryPolicyTypedConfig(): google_protobuf_any_pb.Any | undefined;
  setRetryPolicyTypedConfig(value?: google_protobuf_any_pb.Any): void;

  clearRequestMirrorPoliciesList(): void;
  getRequestMirrorPoliciesList(): Array<RouteAction.RequestMirrorPolicy>;
  setRequestMirrorPoliciesList(value: Array<RouteAction.RequestMirrorPolicy>): void;
  addRequestMirrorPolicies(value?: RouteAction.RequestMirrorPolicy, index?: number): RouteAction.RequestMirrorPolicy;

  getPriority(): github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb.RoutingPriorityMap[keyof github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb.RoutingPriorityMap];
  setPriority(value: github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb.RoutingPriorityMap[keyof github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb.RoutingPriorityMap]): void;

  clearRateLimitsList(): void;
  getRateLimitsList(): Array<RateLimit>;
  setRateLimitsList(value: Array<RateLimit>): void;
  addRateLimits(value?: RateLimit, index?: number): RateLimit;

  hasIncludeVhRateLimits(): boolean;
  clearIncludeVhRateLimits(): void;
  getIncludeVhRateLimits(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setIncludeVhRateLimits(value?: google_protobuf_wrappers_pb.BoolValue): void;

  clearHashPolicyList(): void;
  getHashPolicyList(): Array<RouteAction.HashPolicy>;
  setHashPolicyList(value: Array<RouteAction.HashPolicy>): void;
  addHashPolicy(value?: RouteAction.HashPolicy, index?: number): RouteAction.HashPolicy;

  hasCors(): boolean;
  clearCors(): void;
  getCors(): CorsPolicy | undefined;
  setCors(value?: CorsPolicy): void;

  hasMaxGrpcTimeout(): boolean;
  clearMaxGrpcTimeout(): void;
  getMaxGrpcTimeout(): google_protobuf_duration_pb.Duration | undefined;
  setMaxGrpcTimeout(value?: google_protobuf_duration_pb.Duration): void;

  hasGrpcTimeoutOffset(): boolean;
  clearGrpcTimeoutOffset(): void;
  getGrpcTimeoutOffset(): google_protobuf_duration_pb.Duration | undefined;
  setGrpcTimeoutOffset(value?: google_protobuf_duration_pb.Duration): void;

  clearUpgradeConfigsList(): void;
  getUpgradeConfigsList(): Array<RouteAction.UpgradeConfig>;
  setUpgradeConfigsList(value: Array<RouteAction.UpgradeConfig>): void;
  addUpgradeConfigs(value?: RouteAction.UpgradeConfig, index?: number): RouteAction.UpgradeConfig;

  hasInternalRedirectPolicy(): boolean;
  clearInternalRedirectPolicy(): void;
  getInternalRedirectPolicy(): InternalRedirectPolicy | undefined;
  setInternalRedirectPolicy(value?: InternalRedirectPolicy): void;

  getInternalRedirectAction(): RouteAction.InternalRedirectActionMap[keyof RouteAction.InternalRedirectActionMap];
  setInternalRedirectAction(value: RouteAction.InternalRedirectActionMap[keyof RouteAction.InternalRedirectActionMap]): void;

  hasMaxInternalRedirects(): boolean;
  clearMaxInternalRedirects(): void;
  getMaxInternalRedirects(): google_protobuf_wrappers_pb.UInt32Value | undefined;
  setMaxInternalRedirects(value?: google_protobuf_wrappers_pb.UInt32Value): void;

  hasHedgePolicy(): boolean;
  clearHedgePolicy(): void;
  getHedgePolicy(): HedgePolicy | undefined;
  setHedgePolicy(value?: HedgePolicy): void;

  getClusterSpecifierCase(): RouteAction.ClusterSpecifierCase;
  getHostRewriteSpecifierCase(): RouteAction.HostRewriteSpecifierCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RouteAction.AsObject;
  static toObject(includeInstance: boolean, msg: RouteAction): RouteAction.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RouteAction, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RouteAction;
  static deserializeBinaryFromReader(message: RouteAction, reader: jspb.BinaryReader): RouteAction;
}

export namespace RouteAction {
  export type AsObject = {
    cluster: string,
    clusterHeader: string,
    weightedClusters?: WeightedCluster.AsObject,
    clusterNotFoundResponseCode: RouteAction.ClusterNotFoundResponseCodeMap[keyof RouteAction.ClusterNotFoundResponseCodeMap],
    metadataMatch?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb.Metadata.AsObject,
    prefixRewrite: string,
    regexRewrite?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_matcher_v3_regex_pb.RegexMatchAndSubstitute.AsObject,
    hostRewriteLiteral: string,
    autoHostRewrite?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    hostRewriteHeader: string,
    timeout?: google_protobuf_duration_pb.Duration.AsObject,
    idleTimeout?: google_protobuf_duration_pb.Duration.AsObject,
    retryPolicy?: RetryPolicy.AsObject,
    retryPolicyTypedConfig?: google_protobuf_any_pb.Any.AsObject,
    requestMirrorPoliciesList: Array<RouteAction.RequestMirrorPolicy.AsObject>,
    priority: github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb.RoutingPriorityMap[keyof github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb.RoutingPriorityMap],
    rateLimitsList: Array<RateLimit.AsObject>,
    includeVhRateLimits?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    hashPolicyList: Array<RouteAction.HashPolicy.AsObject>,
    cors?: CorsPolicy.AsObject,
    maxGrpcTimeout?: google_protobuf_duration_pb.Duration.AsObject,
    grpcTimeoutOffset?: google_protobuf_duration_pb.Duration.AsObject,
    upgradeConfigsList: Array<RouteAction.UpgradeConfig.AsObject>,
    internalRedirectPolicy?: InternalRedirectPolicy.AsObject,
    internalRedirectAction: RouteAction.InternalRedirectActionMap[keyof RouteAction.InternalRedirectActionMap],
    maxInternalRedirects?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
    hedgePolicy?: HedgePolicy.AsObject,
  }

  export class RequestMirrorPolicy extends jspb.Message {
    getCluster(): string;
    setCluster(value: string): void;

    hasRuntimeFraction(): boolean;
    clearRuntimeFraction(): void;
    getRuntimeFraction(): github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb.RuntimeFractionalPercent | undefined;
    setRuntimeFraction(value?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb.RuntimeFractionalPercent): void;

    hasTraceSampled(): boolean;
    clearTraceSampled(): void;
    getTraceSampled(): google_protobuf_wrappers_pb.BoolValue | undefined;
    setTraceSampled(value?: google_protobuf_wrappers_pb.BoolValue): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): RequestMirrorPolicy.AsObject;
    static toObject(includeInstance: boolean, msg: RequestMirrorPolicy): RequestMirrorPolicy.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: RequestMirrorPolicy, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): RequestMirrorPolicy;
    static deserializeBinaryFromReader(message: RequestMirrorPolicy, reader: jspb.BinaryReader): RequestMirrorPolicy;
  }

  export namespace RequestMirrorPolicy {
    export type AsObject = {
      cluster: string,
      runtimeFraction?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb.RuntimeFractionalPercent.AsObject,
      traceSampled?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    }
  }

  export class HashPolicy extends jspb.Message {
    hasHeader(): boolean;
    clearHeader(): void;
    getHeader(): RouteAction.HashPolicy.Header | undefined;
    setHeader(value?: RouteAction.HashPolicy.Header): void;

    hasCookie(): boolean;
    clearCookie(): void;
    getCookie(): RouteAction.HashPolicy.Cookie | undefined;
    setCookie(value?: RouteAction.HashPolicy.Cookie): void;

    hasConnectionProperties(): boolean;
    clearConnectionProperties(): void;
    getConnectionProperties(): RouteAction.HashPolicy.ConnectionProperties | undefined;
    setConnectionProperties(value?: RouteAction.HashPolicy.ConnectionProperties): void;

    hasQueryParameter(): boolean;
    clearQueryParameter(): void;
    getQueryParameter(): RouteAction.HashPolicy.QueryParameter | undefined;
    setQueryParameter(value?: RouteAction.HashPolicy.QueryParameter): void;

    hasFilterState(): boolean;
    clearFilterState(): void;
    getFilterState(): RouteAction.HashPolicy.FilterState | undefined;
    setFilterState(value?: RouteAction.HashPolicy.FilterState): void;

    getTerminal(): boolean;
    setTerminal(value: boolean): void;

    getPolicySpecifierCase(): HashPolicy.PolicySpecifierCase;
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): HashPolicy.AsObject;
    static toObject(includeInstance: boolean, msg: HashPolicy): HashPolicy.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: HashPolicy, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): HashPolicy;
    static deserializeBinaryFromReader(message: HashPolicy, reader: jspb.BinaryReader): HashPolicy;
  }

  export namespace HashPolicy {
    export type AsObject = {
      header?: RouteAction.HashPolicy.Header.AsObject,
      cookie?: RouteAction.HashPolicy.Cookie.AsObject,
      connectionProperties?: RouteAction.HashPolicy.ConnectionProperties.AsObject,
      queryParameter?: RouteAction.HashPolicy.QueryParameter.AsObject,
      filterState?: RouteAction.HashPolicy.FilterState.AsObject,
      terminal: boolean,
    }

    export class Header extends jspb.Message {
      getHeaderName(): string;
      setHeaderName(value: string): void;

      hasRegexRewrite(): boolean;
      clearRegexRewrite(): void;
      getRegexRewrite(): github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_matcher_v3_regex_pb.RegexMatchAndSubstitute | undefined;
      setRegexRewrite(value?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_matcher_v3_regex_pb.RegexMatchAndSubstitute): void;

      serializeBinary(): Uint8Array;
      toObject(includeInstance?: boolean): Header.AsObject;
      static toObject(includeInstance: boolean, msg: Header): Header.AsObject;
      static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
      static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
      static serializeBinaryToWriter(message: Header, writer: jspb.BinaryWriter): void;
      static deserializeBinary(bytes: Uint8Array): Header;
      static deserializeBinaryFromReader(message: Header, reader: jspb.BinaryReader): Header;
    }

    export namespace Header {
      export type AsObject = {
        headerName: string,
        regexRewrite?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_matcher_v3_regex_pb.RegexMatchAndSubstitute.AsObject,
      }
    }

    export class Cookie extends jspb.Message {
      getName(): string;
      setName(value: string): void;

      hasTtl(): boolean;
      clearTtl(): void;
      getTtl(): google_protobuf_duration_pb.Duration | undefined;
      setTtl(value?: google_protobuf_duration_pb.Duration): void;

      getPath(): string;
      setPath(value: string): void;

      serializeBinary(): Uint8Array;
      toObject(includeInstance?: boolean): Cookie.AsObject;
      static toObject(includeInstance: boolean, msg: Cookie): Cookie.AsObject;
      static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
      static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
      static serializeBinaryToWriter(message: Cookie, writer: jspb.BinaryWriter): void;
      static deserializeBinary(bytes: Uint8Array): Cookie;
      static deserializeBinaryFromReader(message: Cookie, reader: jspb.BinaryReader): Cookie;
    }

    export namespace Cookie {
      export type AsObject = {
        name: string,
        ttl?: google_protobuf_duration_pb.Duration.AsObject,
        path: string,
      }
    }

    export class ConnectionProperties extends jspb.Message {
      getSourceIp(): boolean;
      setSourceIp(value: boolean): void;

      serializeBinary(): Uint8Array;
      toObject(includeInstance?: boolean): ConnectionProperties.AsObject;
      static toObject(includeInstance: boolean, msg: ConnectionProperties): ConnectionProperties.AsObject;
      static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
      static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
      static serializeBinaryToWriter(message: ConnectionProperties, writer: jspb.BinaryWriter): void;
      static deserializeBinary(bytes: Uint8Array): ConnectionProperties;
      static deserializeBinaryFromReader(message: ConnectionProperties, reader: jspb.BinaryReader): ConnectionProperties;
    }

    export namespace ConnectionProperties {
      export type AsObject = {
        sourceIp: boolean,
      }
    }

    export class QueryParameter extends jspb.Message {
      getName(): string;
      setName(value: string): void;

      serializeBinary(): Uint8Array;
      toObject(includeInstance?: boolean): QueryParameter.AsObject;
      static toObject(includeInstance: boolean, msg: QueryParameter): QueryParameter.AsObject;
      static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
      static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
      static serializeBinaryToWriter(message: QueryParameter, writer: jspb.BinaryWriter): void;
      static deserializeBinary(bytes: Uint8Array): QueryParameter;
      static deserializeBinaryFromReader(message: QueryParameter, reader: jspb.BinaryReader): QueryParameter;
    }

    export namespace QueryParameter {
      export type AsObject = {
        name: string,
      }
    }

    export class FilterState extends jspb.Message {
      getKey(): string;
      setKey(value: string): void;

      serializeBinary(): Uint8Array;
      toObject(includeInstance?: boolean): FilterState.AsObject;
      static toObject(includeInstance: boolean, msg: FilterState): FilterState.AsObject;
      static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
      static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
      static serializeBinaryToWriter(message: FilterState, writer: jspb.BinaryWriter): void;
      static deserializeBinary(bytes: Uint8Array): FilterState;
      static deserializeBinaryFromReader(message: FilterState, reader: jspb.BinaryReader): FilterState;
    }

    export namespace FilterState {
      export type AsObject = {
        key: string,
      }
    }

    export enum PolicySpecifierCase {
      POLICY_SPECIFIER_NOT_SET = 0,
      HEADER = 1,
      COOKIE = 2,
      CONNECTION_PROPERTIES = 3,
      QUERY_PARAMETER = 5,
      FILTER_STATE = 6,
    }
  }

  export class UpgradeConfig extends jspb.Message {
    getUpgradeType(): string;
    setUpgradeType(value: string): void;

    hasEnabled(): boolean;
    clearEnabled(): void;
    getEnabled(): google_protobuf_wrappers_pb.BoolValue | undefined;
    setEnabled(value?: google_protobuf_wrappers_pb.BoolValue): void;

    hasConnectConfig(): boolean;
    clearConnectConfig(): void;
    getConnectConfig(): RouteAction.UpgradeConfig.ConnectConfig | undefined;
    setConnectConfig(value?: RouteAction.UpgradeConfig.ConnectConfig): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): UpgradeConfig.AsObject;
    static toObject(includeInstance: boolean, msg: UpgradeConfig): UpgradeConfig.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: UpgradeConfig, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): UpgradeConfig;
    static deserializeBinaryFromReader(message: UpgradeConfig, reader: jspb.BinaryReader): UpgradeConfig;
  }

  export namespace UpgradeConfig {
    export type AsObject = {
      upgradeType: string,
      enabled?: google_protobuf_wrappers_pb.BoolValue.AsObject,
      connectConfig?: RouteAction.UpgradeConfig.ConnectConfig.AsObject,
    }

    export class ConnectConfig extends jspb.Message {
      hasProxyProtocolConfig(): boolean;
      clearProxyProtocolConfig(): void;
      getProxyProtocolConfig(): github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_proxy_protocol_pb.ProxyProtocolConfig | undefined;
      setProxyProtocolConfig(value?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_proxy_protocol_pb.ProxyProtocolConfig): void;

      serializeBinary(): Uint8Array;
      toObject(includeInstance?: boolean): ConnectConfig.AsObject;
      static toObject(includeInstance: boolean, msg: ConnectConfig): ConnectConfig.AsObject;
      static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
      static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
      static serializeBinaryToWriter(message: ConnectConfig, writer: jspb.BinaryWriter): void;
      static deserializeBinary(bytes: Uint8Array): ConnectConfig;
      static deserializeBinaryFromReader(message: ConnectConfig, reader: jspb.BinaryReader): ConnectConfig;
    }

    export namespace ConnectConfig {
      export type AsObject = {
        proxyProtocolConfig?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_proxy_protocol_pb.ProxyProtocolConfig.AsObject,
      }
    }
  }

  export interface ClusterNotFoundResponseCodeMap {
    SERVICE_UNAVAILABLE: 0;
    NOT_FOUND: 1;
  }

  export const ClusterNotFoundResponseCode: ClusterNotFoundResponseCodeMap;

  export interface InternalRedirectActionMap {
    PASS_THROUGH_INTERNAL_REDIRECT: 0;
    HANDLE_INTERNAL_REDIRECT: 1;
  }

  export const InternalRedirectAction: InternalRedirectActionMap;

  export enum ClusterSpecifierCase {
    CLUSTER_SPECIFIER_NOT_SET = 0,
    CLUSTER = 1,
    CLUSTER_HEADER = 2,
    WEIGHTED_CLUSTERS = 3,
  }

  export enum HostRewriteSpecifierCase {
    HOST_REWRITE_SPECIFIER_NOT_SET = 0,
    HOST_REWRITE_LITERAL = 6,
    AUTO_HOST_REWRITE = 7,
    HOST_REWRITE_HEADER = 29,
  }
}

export class RetryPolicy extends jspb.Message {
  getRetryOn(): string;
  setRetryOn(value: string): void;

  hasNumRetries(): boolean;
  clearNumRetries(): void;
  getNumRetries(): google_protobuf_wrappers_pb.UInt32Value | undefined;
  setNumRetries(value?: google_protobuf_wrappers_pb.UInt32Value): void;

  hasPerTryTimeout(): boolean;
  clearPerTryTimeout(): void;
  getPerTryTimeout(): google_protobuf_duration_pb.Duration | undefined;
  setPerTryTimeout(value?: google_protobuf_duration_pb.Duration): void;

  hasRetryPriority(): boolean;
  clearRetryPriority(): void;
  getRetryPriority(): RetryPolicy.RetryPriority | undefined;
  setRetryPriority(value?: RetryPolicy.RetryPriority): void;

  clearRetryHostPredicateList(): void;
  getRetryHostPredicateList(): Array<RetryPolicy.RetryHostPredicate>;
  setRetryHostPredicateList(value: Array<RetryPolicy.RetryHostPredicate>): void;
  addRetryHostPredicate(value?: RetryPolicy.RetryHostPredicate, index?: number): RetryPolicy.RetryHostPredicate;

  getHostSelectionRetryMaxAttempts(): number;
  setHostSelectionRetryMaxAttempts(value: number): void;

  clearRetriableStatusCodesList(): void;
  getRetriableStatusCodesList(): Array<number>;
  setRetriableStatusCodesList(value: Array<number>): void;
  addRetriableStatusCodes(value: number, index?: number): number;

  hasRetryBackOff(): boolean;
  clearRetryBackOff(): void;
  getRetryBackOff(): RetryPolicy.RetryBackOff | undefined;
  setRetryBackOff(value?: RetryPolicy.RetryBackOff): void;

  clearRetriableHeadersList(): void;
  getRetriableHeadersList(): Array<HeaderMatcher>;
  setRetriableHeadersList(value: Array<HeaderMatcher>): void;
  addRetriableHeaders(value?: HeaderMatcher, index?: number): HeaderMatcher;

  clearRetriableRequestHeadersList(): void;
  getRetriableRequestHeadersList(): Array<HeaderMatcher>;
  setRetriableRequestHeadersList(value: Array<HeaderMatcher>): void;
  addRetriableRequestHeaders(value?: HeaderMatcher, index?: number): HeaderMatcher;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RetryPolicy.AsObject;
  static toObject(includeInstance: boolean, msg: RetryPolicy): RetryPolicy.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RetryPolicy, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RetryPolicy;
  static deserializeBinaryFromReader(message: RetryPolicy, reader: jspb.BinaryReader): RetryPolicy;
}

export namespace RetryPolicy {
  export type AsObject = {
    retryOn: string,
    numRetries?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
    perTryTimeout?: google_protobuf_duration_pb.Duration.AsObject,
    retryPriority?: RetryPolicy.RetryPriority.AsObject,
    retryHostPredicateList: Array<RetryPolicy.RetryHostPredicate.AsObject>,
    hostSelectionRetryMaxAttempts: number,
    retriableStatusCodesList: Array<number>,
    retryBackOff?: RetryPolicy.RetryBackOff.AsObject,
    retriableHeadersList: Array<HeaderMatcher.AsObject>,
    retriableRequestHeadersList: Array<HeaderMatcher.AsObject>,
  }

  export class RetryPriority extends jspb.Message {
    getName(): string;
    setName(value: string): void;

    hasTypedConfig(): boolean;
    clearTypedConfig(): void;
    getTypedConfig(): google_protobuf_any_pb.Any | undefined;
    setTypedConfig(value?: google_protobuf_any_pb.Any): void;

    getConfigTypeCase(): RetryPriority.ConfigTypeCase;
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): RetryPriority.AsObject;
    static toObject(includeInstance: boolean, msg: RetryPriority): RetryPriority.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: RetryPriority, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): RetryPriority;
    static deserializeBinaryFromReader(message: RetryPriority, reader: jspb.BinaryReader): RetryPriority;
  }

  export namespace RetryPriority {
    export type AsObject = {
      name: string,
      typedConfig?: google_protobuf_any_pb.Any.AsObject,
    }

    export enum ConfigTypeCase {
      CONFIG_TYPE_NOT_SET = 0,
      TYPED_CONFIG = 3,
    }
  }

  export class RetryHostPredicate extends jspb.Message {
    getName(): string;
    setName(value: string): void;

    hasTypedConfig(): boolean;
    clearTypedConfig(): void;
    getTypedConfig(): google_protobuf_any_pb.Any | undefined;
    setTypedConfig(value?: google_protobuf_any_pb.Any): void;

    getConfigTypeCase(): RetryHostPredicate.ConfigTypeCase;
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): RetryHostPredicate.AsObject;
    static toObject(includeInstance: boolean, msg: RetryHostPredicate): RetryHostPredicate.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: RetryHostPredicate, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): RetryHostPredicate;
    static deserializeBinaryFromReader(message: RetryHostPredicate, reader: jspb.BinaryReader): RetryHostPredicate;
  }

  export namespace RetryHostPredicate {
    export type AsObject = {
      name: string,
      typedConfig?: google_protobuf_any_pb.Any.AsObject,
    }

    export enum ConfigTypeCase {
      CONFIG_TYPE_NOT_SET = 0,
      TYPED_CONFIG = 3,
    }
  }

  export class RetryBackOff extends jspb.Message {
    hasBaseInterval(): boolean;
    clearBaseInterval(): void;
    getBaseInterval(): google_protobuf_duration_pb.Duration | undefined;
    setBaseInterval(value?: google_protobuf_duration_pb.Duration): void;

    hasMaxInterval(): boolean;
    clearMaxInterval(): void;
    getMaxInterval(): google_protobuf_duration_pb.Duration | undefined;
    setMaxInterval(value?: google_protobuf_duration_pb.Duration): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): RetryBackOff.AsObject;
    static toObject(includeInstance: boolean, msg: RetryBackOff): RetryBackOff.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: RetryBackOff, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): RetryBackOff;
    static deserializeBinaryFromReader(message: RetryBackOff, reader: jspb.BinaryReader): RetryBackOff;
  }

  export namespace RetryBackOff {
    export type AsObject = {
      baseInterval?: google_protobuf_duration_pb.Duration.AsObject,
      maxInterval?: google_protobuf_duration_pb.Duration.AsObject,
    }
  }
}

export class HedgePolicy extends jspb.Message {
  hasInitialRequests(): boolean;
  clearInitialRequests(): void;
  getInitialRequests(): google_protobuf_wrappers_pb.UInt32Value | undefined;
  setInitialRequests(value?: google_protobuf_wrappers_pb.UInt32Value): void;

  hasAdditionalRequestChance(): boolean;
  clearAdditionalRequestChance(): void;
  getAdditionalRequestChance(): github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_v3_percent_pb.FractionalPercent | undefined;
  setAdditionalRequestChance(value?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_v3_percent_pb.FractionalPercent): void;

  getHedgeOnPerTryTimeout(): boolean;
  setHedgeOnPerTryTimeout(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HedgePolicy.AsObject;
  static toObject(includeInstance: boolean, msg: HedgePolicy): HedgePolicy.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: HedgePolicy, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HedgePolicy;
  static deserializeBinaryFromReader(message: HedgePolicy, reader: jspb.BinaryReader): HedgePolicy;
}

export namespace HedgePolicy {
  export type AsObject = {
    initialRequests?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
    additionalRequestChance?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_v3_percent_pb.FractionalPercent.AsObject,
    hedgeOnPerTryTimeout: boolean,
  }
}

export class RedirectAction extends jspb.Message {
  hasHttpsRedirect(): boolean;
  clearHttpsRedirect(): void;
  getHttpsRedirect(): boolean;
  setHttpsRedirect(value: boolean): void;

  hasSchemeRedirect(): boolean;
  clearSchemeRedirect(): void;
  getSchemeRedirect(): string;
  setSchemeRedirect(value: string): void;

  getHostRedirect(): string;
  setHostRedirect(value: string): void;

  getPortRedirect(): number;
  setPortRedirect(value: number): void;

  hasPathRedirect(): boolean;
  clearPathRedirect(): void;
  getPathRedirect(): string;
  setPathRedirect(value: string): void;

  hasPrefixRewrite(): boolean;
  clearPrefixRewrite(): void;
  getPrefixRewrite(): string;
  setPrefixRewrite(value: string): void;

  getResponseCode(): RedirectAction.RedirectResponseCodeMap[keyof RedirectAction.RedirectResponseCodeMap];
  setResponseCode(value: RedirectAction.RedirectResponseCodeMap[keyof RedirectAction.RedirectResponseCodeMap]): void;

  getStripQuery(): boolean;
  setStripQuery(value: boolean): void;

  getSchemeRewriteSpecifierCase(): RedirectAction.SchemeRewriteSpecifierCase;
  getPathRewriteSpecifierCase(): RedirectAction.PathRewriteSpecifierCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RedirectAction.AsObject;
  static toObject(includeInstance: boolean, msg: RedirectAction): RedirectAction.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RedirectAction, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RedirectAction;
  static deserializeBinaryFromReader(message: RedirectAction, reader: jspb.BinaryReader): RedirectAction;
}

export namespace RedirectAction {
  export type AsObject = {
    httpsRedirect: boolean,
    schemeRedirect: string,
    hostRedirect: string,
    portRedirect: number,
    pathRedirect: string,
    prefixRewrite: string,
    responseCode: RedirectAction.RedirectResponseCodeMap[keyof RedirectAction.RedirectResponseCodeMap],
    stripQuery: boolean,
  }

  export interface RedirectResponseCodeMap {
    MOVED_PERMANENTLY: 0;
    FOUND: 1;
    SEE_OTHER: 2;
    TEMPORARY_REDIRECT: 3;
    PERMANENT_REDIRECT: 4;
  }

  export const RedirectResponseCode: RedirectResponseCodeMap;

  export enum SchemeRewriteSpecifierCase {
    SCHEME_REWRITE_SPECIFIER_NOT_SET = 0,
    HTTPS_REDIRECT = 4,
    SCHEME_REDIRECT = 7,
  }

  export enum PathRewriteSpecifierCase {
    PATH_REWRITE_SPECIFIER_NOT_SET = 0,
    PATH_REDIRECT = 2,
    PREFIX_REWRITE = 5,
  }
}

export class DirectResponseAction extends jspb.Message {
  getStatus(): number;
  setStatus(value: number): void;

  hasBody(): boolean;
  clearBody(): void;
  getBody(): github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb.DataSource | undefined;
  setBody(value?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb.DataSource): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DirectResponseAction.AsObject;
  static toObject(includeInstance: boolean, msg: DirectResponseAction): DirectResponseAction.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DirectResponseAction, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DirectResponseAction;
  static deserializeBinaryFromReader(message: DirectResponseAction, reader: jspb.BinaryReader): DirectResponseAction;
}

export namespace DirectResponseAction {
  export type AsObject = {
    status: number,
    body?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_base_pb.DataSource.AsObject,
  }
}

export class Decorator extends jspb.Message {
  getOperation(): string;
  setOperation(value: string): void;

  hasPropagate(): boolean;
  clearPropagate(): void;
  getPropagate(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setPropagate(value?: google_protobuf_wrappers_pb.BoolValue): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Decorator.AsObject;
  static toObject(includeInstance: boolean, msg: Decorator): Decorator.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Decorator, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Decorator;
  static deserializeBinaryFromReader(message: Decorator, reader: jspb.BinaryReader): Decorator;
}

export namespace Decorator {
  export type AsObject = {
    operation: string,
    propagate?: google_protobuf_wrappers_pb.BoolValue.AsObject,
  }
}

export class Tracing extends jspb.Message {
  hasClientSampling(): boolean;
  clearClientSampling(): void;
  getClientSampling(): github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_v3_percent_pb.FractionalPercent | undefined;
  setClientSampling(value?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_v3_percent_pb.FractionalPercent): void;

  hasRandomSampling(): boolean;
  clearRandomSampling(): void;
  getRandomSampling(): github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_v3_percent_pb.FractionalPercent | undefined;
  setRandomSampling(value?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_v3_percent_pb.FractionalPercent): void;

  hasOverallSampling(): boolean;
  clearOverallSampling(): void;
  getOverallSampling(): github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_v3_percent_pb.FractionalPercent | undefined;
  setOverallSampling(value?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_v3_percent_pb.FractionalPercent): void;

  clearCustomTagsList(): void;
  getCustomTagsList(): Array<github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_tracing_v3_custom_tag_pb.CustomTag>;
  setCustomTagsList(value: Array<github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_tracing_v3_custom_tag_pb.CustomTag>): void;
  addCustomTags(value?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_tracing_v3_custom_tag_pb.CustomTag, index?: number): github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_tracing_v3_custom_tag_pb.CustomTag;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Tracing.AsObject;
  static toObject(includeInstance: boolean, msg: Tracing): Tracing.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Tracing, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Tracing;
  static deserializeBinaryFromReader(message: Tracing, reader: jspb.BinaryReader): Tracing;
}

export namespace Tracing {
  export type AsObject = {
    clientSampling?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_v3_percent_pb.FractionalPercent.AsObject,
    randomSampling?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_v3_percent_pb.FractionalPercent.AsObject,
    overallSampling?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_v3_percent_pb.FractionalPercent.AsObject,
    customTagsList: Array<github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_tracing_v3_custom_tag_pb.CustomTag.AsObject>,
  }
}

export class VirtualCluster extends jspb.Message {
  clearHeadersList(): void;
  getHeadersList(): Array<HeaderMatcher>;
  setHeadersList(value: Array<HeaderMatcher>): void;
  addHeaders(value?: HeaderMatcher, index?: number): HeaderMatcher;

  getName(): string;
  setName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VirtualCluster.AsObject;
  static toObject(includeInstance: boolean, msg: VirtualCluster): VirtualCluster.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: VirtualCluster, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VirtualCluster;
  static deserializeBinaryFromReader(message: VirtualCluster, reader: jspb.BinaryReader): VirtualCluster;
}

export namespace VirtualCluster {
  export type AsObject = {
    headersList: Array<HeaderMatcher.AsObject>,
    name: string,
  }
}

export class RateLimit extends jspb.Message {
  hasStage(): boolean;
  clearStage(): void;
  getStage(): google_protobuf_wrappers_pb.UInt32Value | undefined;
  setStage(value?: google_protobuf_wrappers_pb.UInt32Value): void;

  getDisableKey(): string;
  setDisableKey(value: string): void;

  clearActionsList(): void;
  getActionsList(): Array<RateLimit.Action>;
  setActionsList(value: Array<RateLimit.Action>): void;
  addActions(value?: RateLimit.Action, index?: number): RateLimit.Action;

  hasLimit(): boolean;
  clearLimit(): void;
  getLimit(): RateLimit.Override | undefined;
  setLimit(value?: RateLimit.Override): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RateLimit.AsObject;
  static toObject(includeInstance: boolean, msg: RateLimit): RateLimit.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RateLimit, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RateLimit;
  static deserializeBinaryFromReader(message: RateLimit, reader: jspb.BinaryReader): RateLimit;
}

export namespace RateLimit {
  export type AsObject = {
    stage?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
    disableKey: string,
    actionsList: Array<RateLimit.Action.AsObject>,
    limit?: RateLimit.Override.AsObject,
  }

  export class Action extends jspb.Message {
    hasSourceCluster(): boolean;
    clearSourceCluster(): void;
    getSourceCluster(): RateLimit.Action.SourceCluster | undefined;
    setSourceCluster(value?: RateLimit.Action.SourceCluster): void;

    hasDestinationCluster(): boolean;
    clearDestinationCluster(): void;
    getDestinationCluster(): RateLimit.Action.DestinationCluster | undefined;
    setDestinationCluster(value?: RateLimit.Action.DestinationCluster): void;

    hasRequestHeaders(): boolean;
    clearRequestHeaders(): void;
    getRequestHeaders(): RateLimit.Action.RequestHeaders | undefined;
    setRequestHeaders(value?: RateLimit.Action.RequestHeaders): void;

    hasRemoteAddress(): boolean;
    clearRemoteAddress(): void;
    getRemoteAddress(): RateLimit.Action.RemoteAddress | undefined;
    setRemoteAddress(value?: RateLimit.Action.RemoteAddress): void;

    hasGenericKey(): boolean;
    clearGenericKey(): void;
    getGenericKey(): RateLimit.Action.GenericKey | undefined;
    setGenericKey(value?: RateLimit.Action.GenericKey): void;

    hasHeaderValueMatch(): boolean;
    clearHeaderValueMatch(): void;
    getHeaderValueMatch(): RateLimit.Action.HeaderValueMatch | undefined;
    setHeaderValueMatch(value?: RateLimit.Action.HeaderValueMatch): void;

    hasDynamicMetadata(): boolean;
    clearDynamicMetadata(): void;
    getDynamicMetadata(): RateLimit.Action.DynamicMetaData | undefined;
    setDynamicMetadata(value?: RateLimit.Action.DynamicMetaData): void;

    getActionSpecifierCase(): Action.ActionSpecifierCase;
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Action.AsObject;
    static toObject(includeInstance: boolean, msg: Action): Action.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Action, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Action;
    static deserializeBinaryFromReader(message: Action, reader: jspb.BinaryReader): Action;
  }

  export namespace Action {
    export type AsObject = {
      sourceCluster?: RateLimit.Action.SourceCluster.AsObject,
      destinationCluster?: RateLimit.Action.DestinationCluster.AsObject,
      requestHeaders?: RateLimit.Action.RequestHeaders.AsObject,
      remoteAddress?: RateLimit.Action.RemoteAddress.AsObject,
      genericKey?: RateLimit.Action.GenericKey.AsObject,
      headerValueMatch?: RateLimit.Action.HeaderValueMatch.AsObject,
      dynamicMetadata?: RateLimit.Action.DynamicMetaData.AsObject,
    }

    export class SourceCluster extends jspb.Message {
      serializeBinary(): Uint8Array;
      toObject(includeInstance?: boolean): SourceCluster.AsObject;
      static toObject(includeInstance: boolean, msg: SourceCluster): SourceCluster.AsObject;
      static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
      static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
      static serializeBinaryToWriter(message: SourceCluster, writer: jspb.BinaryWriter): void;
      static deserializeBinary(bytes: Uint8Array): SourceCluster;
      static deserializeBinaryFromReader(message: SourceCluster, reader: jspb.BinaryReader): SourceCluster;
    }

    export namespace SourceCluster {
      export type AsObject = {
      }
    }

    export class DestinationCluster extends jspb.Message {
      serializeBinary(): Uint8Array;
      toObject(includeInstance?: boolean): DestinationCluster.AsObject;
      static toObject(includeInstance: boolean, msg: DestinationCluster): DestinationCluster.AsObject;
      static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
      static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
      static serializeBinaryToWriter(message: DestinationCluster, writer: jspb.BinaryWriter): void;
      static deserializeBinary(bytes: Uint8Array): DestinationCluster;
      static deserializeBinaryFromReader(message: DestinationCluster, reader: jspb.BinaryReader): DestinationCluster;
    }

    export namespace DestinationCluster {
      export type AsObject = {
      }
    }

    export class RequestHeaders extends jspb.Message {
      getHeaderName(): string;
      setHeaderName(value: string): void;

      getDescriptorKey(): string;
      setDescriptorKey(value: string): void;

      getSkipIfAbsent(): boolean;
      setSkipIfAbsent(value: boolean): void;

      serializeBinary(): Uint8Array;
      toObject(includeInstance?: boolean): RequestHeaders.AsObject;
      static toObject(includeInstance: boolean, msg: RequestHeaders): RequestHeaders.AsObject;
      static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
      static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
      static serializeBinaryToWriter(message: RequestHeaders, writer: jspb.BinaryWriter): void;
      static deserializeBinary(bytes: Uint8Array): RequestHeaders;
      static deserializeBinaryFromReader(message: RequestHeaders, reader: jspb.BinaryReader): RequestHeaders;
    }

    export namespace RequestHeaders {
      export type AsObject = {
        headerName: string,
        descriptorKey: string,
        skipIfAbsent: boolean,
      }
    }

    export class RemoteAddress extends jspb.Message {
      serializeBinary(): Uint8Array;
      toObject(includeInstance?: boolean): RemoteAddress.AsObject;
      static toObject(includeInstance: boolean, msg: RemoteAddress): RemoteAddress.AsObject;
      static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
      static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
      static serializeBinaryToWriter(message: RemoteAddress, writer: jspb.BinaryWriter): void;
      static deserializeBinary(bytes: Uint8Array): RemoteAddress;
      static deserializeBinaryFromReader(message: RemoteAddress, reader: jspb.BinaryReader): RemoteAddress;
    }

    export namespace RemoteAddress {
      export type AsObject = {
      }
    }

    export class GenericKey extends jspb.Message {
      getDescriptorValue(): string;
      setDescriptorValue(value: string): void;

      serializeBinary(): Uint8Array;
      toObject(includeInstance?: boolean): GenericKey.AsObject;
      static toObject(includeInstance: boolean, msg: GenericKey): GenericKey.AsObject;
      static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
      static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
      static serializeBinaryToWriter(message: GenericKey, writer: jspb.BinaryWriter): void;
      static deserializeBinary(bytes: Uint8Array): GenericKey;
      static deserializeBinaryFromReader(message: GenericKey, reader: jspb.BinaryReader): GenericKey;
    }

    export namespace GenericKey {
      export type AsObject = {
        descriptorValue: string,
      }
    }

    export class HeaderValueMatch extends jspb.Message {
      getDescriptorValue(): string;
      setDescriptorValue(value: string): void;

      hasExpectMatch(): boolean;
      clearExpectMatch(): void;
      getExpectMatch(): google_protobuf_wrappers_pb.BoolValue | undefined;
      setExpectMatch(value?: google_protobuf_wrappers_pb.BoolValue): void;

      clearHeadersList(): void;
      getHeadersList(): Array<HeaderMatcher>;
      setHeadersList(value: Array<HeaderMatcher>): void;
      addHeaders(value?: HeaderMatcher, index?: number): HeaderMatcher;

      serializeBinary(): Uint8Array;
      toObject(includeInstance?: boolean): HeaderValueMatch.AsObject;
      static toObject(includeInstance: boolean, msg: HeaderValueMatch): HeaderValueMatch.AsObject;
      static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
      static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
      static serializeBinaryToWriter(message: HeaderValueMatch, writer: jspb.BinaryWriter): void;
      static deserializeBinary(bytes: Uint8Array): HeaderValueMatch;
      static deserializeBinaryFromReader(message: HeaderValueMatch, reader: jspb.BinaryReader): HeaderValueMatch;
    }

    export namespace HeaderValueMatch {
      export type AsObject = {
        descriptorValue: string,
        expectMatch?: google_protobuf_wrappers_pb.BoolValue.AsObject,
        headersList: Array<HeaderMatcher.AsObject>,
      }
    }

    export class DynamicMetaData extends jspb.Message {
      getDescriptorKey(): string;
      setDescriptorKey(value: string): void;

      hasMetadataKey(): boolean;
      clearMetadataKey(): void;
      getMetadataKey(): github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_metadata_v3_metadata_pb.MetadataKey | undefined;
      setMetadataKey(value?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_metadata_v3_metadata_pb.MetadataKey): void;

      serializeBinary(): Uint8Array;
      toObject(includeInstance?: boolean): DynamicMetaData.AsObject;
      static toObject(includeInstance: boolean, msg: DynamicMetaData): DynamicMetaData.AsObject;
      static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
      static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
      static serializeBinaryToWriter(message: DynamicMetaData, writer: jspb.BinaryWriter): void;
      static deserializeBinary(bytes: Uint8Array): DynamicMetaData;
      static deserializeBinaryFromReader(message: DynamicMetaData, reader: jspb.BinaryReader): DynamicMetaData;
    }

    export namespace DynamicMetaData {
      export type AsObject = {
        descriptorKey: string,
        metadataKey?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_metadata_v3_metadata_pb.MetadataKey.AsObject,
      }
    }

    export enum ActionSpecifierCase {
      ACTION_SPECIFIER_NOT_SET = 0,
      SOURCE_CLUSTER = 1,
      DESTINATION_CLUSTER = 2,
      REQUEST_HEADERS = 3,
      REMOTE_ADDRESS = 4,
      GENERIC_KEY = 5,
      HEADER_VALUE_MATCH = 6,
      DYNAMIC_METADATA = 7,
    }
  }

  export class Override extends jspb.Message {
    hasDynamicMetadata(): boolean;
    clearDynamicMetadata(): void;
    getDynamicMetadata(): RateLimit.Override.DynamicMetadata | undefined;
    setDynamicMetadata(value?: RateLimit.Override.DynamicMetadata): void;

    getOverrideSpecifierCase(): Override.OverrideSpecifierCase;
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Override.AsObject;
    static toObject(includeInstance: boolean, msg: Override): Override.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Override, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Override;
    static deserializeBinaryFromReader(message: Override, reader: jspb.BinaryReader): Override;
  }

  export namespace Override {
    export type AsObject = {
      dynamicMetadata?: RateLimit.Override.DynamicMetadata.AsObject,
    }

    export class DynamicMetadata extends jspb.Message {
      hasMetadataKey(): boolean;
      clearMetadataKey(): void;
      getMetadataKey(): github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_metadata_v3_metadata_pb.MetadataKey | undefined;
      setMetadataKey(value?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_metadata_v3_metadata_pb.MetadataKey): void;

      serializeBinary(): Uint8Array;
      toObject(includeInstance?: boolean): DynamicMetadata.AsObject;
      static toObject(includeInstance: boolean, msg: DynamicMetadata): DynamicMetadata.AsObject;
      static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
      static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
      static serializeBinaryToWriter(message: DynamicMetadata, writer: jspb.BinaryWriter): void;
      static deserializeBinary(bytes: Uint8Array): DynamicMetadata;
      static deserializeBinaryFromReader(message: DynamicMetadata, reader: jspb.BinaryReader): DynamicMetadata;
    }

    export namespace DynamicMetadata {
      export type AsObject = {
        metadataKey?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_metadata_v3_metadata_pb.MetadataKey.AsObject,
      }
    }

    export enum OverrideSpecifierCase {
      OVERRIDE_SPECIFIER_NOT_SET = 0,
      DYNAMIC_METADATA = 1,
    }
  }
}

export class HeaderMatcher extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  hasExactMatch(): boolean;
  clearExactMatch(): void;
  getExactMatch(): string;
  setExactMatch(value: string): void;

  hasSafeRegexMatch(): boolean;
  clearSafeRegexMatch(): void;
  getSafeRegexMatch(): github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_matcher_v3_regex_pb.RegexMatcher | undefined;
  setSafeRegexMatch(value?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_matcher_v3_regex_pb.RegexMatcher): void;

  hasRangeMatch(): boolean;
  clearRangeMatch(): void;
  getRangeMatch(): github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_v3_range_pb.Int64Range | undefined;
  setRangeMatch(value?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_v3_range_pb.Int64Range): void;

  hasPresentMatch(): boolean;
  clearPresentMatch(): void;
  getPresentMatch(): boolean;
  setPresentMatch(value: boolean): void;

  hasPrefixMatch(): boolean;
  clearPrefixMatch(): void;
  getPrefixMatch(): string;
  setPrefixMatch(value: string): void;

  hasSuffixMatch(): boolean;
  clearSuffixMatch(): void;
  getSuffixMatch(): string;
  setSuffixMatch(value: string): void;

  getInvertMatch(): boolean;
  setInvertMatch(value: boolean): void;

  getHeaderMatchSpecifierCase(): HeaderMatcher.HeaderMatchSpecifierCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HeaderMatcher.AsObject;
  static toObject(includeInstance: boolean, msg: HeaderMatcher): HeaderMatcher.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: HeaderMatcher, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HeaderMatcher;
  static deserializeBinaryFromReader(message: HeaderMatcher, reader: jspb.BinaryReader): HeaderMatcher;
}

export namespace HeaderMatcher {
  export type AsObject = {
    name: string,
    exactMatch: string,
    safeRegexMatch?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_matcher_v3_regex_pb.RegexMatcher.AsObject,
    rangeMatch?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_v3_range_pb.Int64Range.AsObject,
    presentMatch: boolean,
    prefixMatch: string,
    suffixMatch: string,
    invertMatch: boolean,
  }

  export enum HeaderMatchSpecifierCase {
    HEADER_MATCH_SPECIFIER_NOT_SET = 0,
    EXACT_MATCH = 4,
    SAFE_REGEX_MATCH = 11,
    RANGE_MATCH = 6,
    PRESENT_MATCH = 7,
    PREFIX_MATCH = 9,
    SUFFIX_MATCH = 10,
  }
}

export class QueryParameterMatcher extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  hasStringMatch(): boolean;
  clearStringMatch(): void;
  getStringMatch(): github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_matcher_v3_string_pb.StringMatcher | undefined;
  setStringMatch(value?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_matcher_v3_string_pb.StringMatcher): void;

  hasPresentMatch(): boolean;
  clearPresentMatch(): void;
  getPresentMatch(): boolean;
  setPresentMatch(value: boolean): void;

  getQueryParameterMatchSpecifierCase(): QueryParameterMatcher.QueryParameterMatchSpecifierCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): QueryParameterMatcher.AsObject;
  static toObject(includeInstance: boolean, msg: QueryParameterMatcher): QueryParameterMatcher.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: QueryParameterMatcher, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): QueryParameterMatcher;
  static deserializeBinaryFromReader(message: QueryParameterMatcher, reader: jspb.BinaryReader): QueryParameterMatcher;
}

export namespace QueryParameterMatcher {
  export type AsObject = {
    name: string,
    stringMatch?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_matcher_v3_string_pb.StringMatcher.AsObject,
    presentMatch: boolean,
  }

  export enum QueryParameterMatchSpecifierCase {
    QUERY_PARAMETER_MATCH_SPECIFIER_NOT_SET = 0,
    STRING_MATCH = 5,
    PRESENT_MATCH = 6,
  }
}

export class InternalRedirectPolicy extends jspb.Message {
  hasMaxInternalRedirects(): boolean;
  clearMaxInternalRedirects(): void;
  getMaxInternalRedirects(): google_protobuf_wrappers_pb.UInt32Value | undefined;
  setMaxInternalRedirects(value?: google_protobuf_wrappers_pb.UInt32Value): void;

  clearRedirectResponseCodesList(): void;
  getRedirectResponseCodesList(): Array<number>;
  setRedirectResponseCodesList(value: Array<number>): void;
  addRedirectResponseCodes(value: number, index?: number): number;

  clearPredicatesList(): void;
  getPredicatesList(): Array<github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_extension_pb.TypedExtensionConfig>;
  setPredicatesList(value: Array<github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_extension_pb.TypedExtensionConfig>): void;
  addPredicates(value?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_extension_pb.TypedExtensionConfig, index?: number): github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_extension_pb.TypedExtensionConfig;

  getAllowCrossSchemeRedirect(): boolean;
  setAllowCrossSchemeRedirect(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): InternalRedirectPolicy.AsObject;
  static toObject(includeInstance: boolean, msg: InternalRedirectPolicy): InternalRedirectPolicy.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: InternalRedirectPolicy, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): InternalRedirectPolicy;
  static deserializeBinaryFromReader(message: InternalRedirectPolicy, reader: jspb.BinaryReader): InternalRedirectPolicy;
}

export namespace InternalRedirectPolicy {
  export type AsObject = {
    maxInternalRedirects?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
    redirectResponseCodesList: Array<number>,
    predicatesList: Array<github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_extension_pb.TypedExtensionConfig.AsObject>,
    allowCrossSchemeRedirect: boolean,
  }
}
