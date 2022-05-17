/* eslint-disable */
// package: gloo.solo.io
// file: github.com/solo-io/solo-apis/api/gloo/gloo/v1/options.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_struct_pb from "google-protobuf/google/protobuf/struct_pb";
import * as extproto_ext_pb from "../../../../../../../extproto/ext_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_extensions_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/extensions_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_cors_cors_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/options/cors/cors_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_rest_rest_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/options/rest/rest_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_grpc_grpc_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/options/grpc/grpc_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_als_als_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/options/als/als_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_proxy_protocol_proxy_protocol_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/options/proxy_protocol/proxy_protocol_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_grpc_web_grpc_web_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/options/grpc_web/grpc_web_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_grpc_json_grpc_json_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/options/grpc_json/grpc_json_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_hcm_hcm_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/options/hcm/hcm_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_lbhash_lbhash_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/options/lbhash/lbhash_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_shadowing_shadowing_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/options/shadowing/shadowing_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_tcp_tcp_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/options/tcp/tcp_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_tracing_tracing_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/options/tracing/tracing_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_retries_retries_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/options/retries/retries_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_stats_stats_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/options/stats/stats_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_faultinjection_fault_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/options/faultinjection/fault_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_headers_headers_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/options/headers/headers_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_aws_aws_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/options/aws/aws_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_wasm_wasm_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/options/wasm/wasm_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_azure_azure_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/options/azure/azure_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_healthcheck_healthcheck_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/options/healthcheck/healthcheck_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_protocol_upgrade_protocol_upgrade_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/options/protocol_upgrade/protocol_upgrade_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_extensions_proxylatency_proxylatency_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/extensions/proxylatency/proxylatency_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_extensions_filters_http_buffer_v3_buffer_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/extensions/filters/http/buffer/v3/buffer_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_extensions_filters_http_csrf_v3_csrf_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/extensions/filters/http/csrf/v3/csrf_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_filter_http_gzip_v2_gzip_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/config/filter/http/gzip/v2/gzip_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_type_matcher_v3_regex_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/type/matcher/v3/regex_pb";
import * as github_com_solo_io_solo_apis_api_gloo_enterprise_gloo_v1_auth_config_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/enterprise.gloo/v1/auth_config_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_jwt_jwt_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/enterprise/options/jwt/jwt_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_ratelimit_ratelimit_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/enterprise/options/ratelimit/ratelimit_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_caching_caching_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/enterprise/options/caching/caching_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_rbac_rbac_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/enterprise/options/rbac/rbac_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_waf_waf_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/enterprise/options/waf/waf_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_dlp_dlp_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/enterprise/options/dlp/dlp_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_transformation_transformation_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/options/transformation/transformation_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_dynamic_forward_proxy_dynamic_forward_proxy_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/options/dynamic_forward_proxy/dynamic_forward_proxy_pb";
import * as github_com_solo_io_solo_kit_api_external_envoy_api_v2_core_base_pb from "../../../../../../../github.com/solo-io/solo-kit/api/external/envoy/api/v2/core/base_pb";
import * as github_com_solo_io_solo_kit_api_external_envoy_api_v2_core_socket_option_pb from "../../../../../../../github.com/solo-io/solo-kit/api/external/envoy/api/v2/core/socket_option_pb";
import * as google_protobuf_duration_pb from "google-protobuf/google/protobuf/duration_pb";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";

export class ListenerOptions extends jspb.Message {
  hasAccessLoggingService(): boolean;
  clearAccessLoggingService(): void;
  getAccessLoggingService(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_als_als_pb.AccessLoggingService | undefined;
  setAccessLoggingService(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_als_als_pb.AccessLoggingService): void;

  hasExtensions(): boolean;
  clearExtensions(): void;
  getExtensions(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_extensions_pb.Extensions | undefined;
  setExtensions(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_extensions_pb.Extensions): void;

  hasPerConnectionBufferLimitBytes(): boolean;
  clearPerConnectionBufferLimitBytes(): void;
  getPerConnectionBufferLimitBytes(): google_protobuf_wrappers_pb.UInt32Value | undefined;
  setPerConnectionBufferLimitBytes(value?: google_protobuf_wrappers_pb.UInt32Value): void;

  clearSocketOptionsList(): void;
  getSocketOptionsList(): Array<github_com_solo_io_solo_kit_api_external_envoy_api_v2_core_socket_option_pb.SocketOption>;
  setSocketOptionsList(value: Array<github_com_solo_io_solo_kit_api_external_envoy_api_v2_core_socket_option_pb.SocketOption>): void;
  addSocketOptions(value?: github_com_solo_io_solo_kit_api_external_envoy_api_v2_core_socket_option_pb.SocketOption, index?: number): github_com_solo_io_solo_kit_api_external_envoy_api_v2_core_socket_option_pb.SocketOption;

  hasProxyProtocol(): boolean;
  clearProxyProtocol(): void;
  getProxyProtocol(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_proxy_protocol_proxy_protocol_pb.ProxyProtocol | undefined;
  setProxyProtocol(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_proxy_protocol_proxy_protocol_pb.ProxyProtocol): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListenerOptions.AsObject;
  static toObject(includeInstance: boolean, msg: ListenerOptions): ListenerOptions.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListenerOptions, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListenerOptions;
  static deserializeBinaryFromReader(message: ListenerOptions, reader: jspb.BinaryReader): ListenerOptions;
}

export namespace ListenerOptions {
  export type AsObject = {
    accessLoggingService?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_als_als_pb.AccessLoggingService.AsObject,
    extensions?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_extensions_pb.Extensions.AsObject,
    perConnectionBufferLimitBytes?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
    socketOptionsList: Array<github_com_solo_io_solo_kit_api_external_envoy_api_v2_core_socket_option_pb.SocketOption.AsObject>,
    proxyProtocol?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_proxy_protocol_proxy_protocol_pb.ProxyProtocol.AsObject,
  }
}

export class RouteConfigurationOptions extends jspb.Message {
  hasMaxDirectResponseBodySizeBytes(): boolean;
  clearMaxDirectResponseBodySizeBytes(): void;
  getMaxDirectResponseBodySizeBytes(): google_protobuf_wrappers_pb.UInt32Value | undefined;
  setMaxDirectResponseBodySizeBytes(value?: google_protobuf_wrappers_pb.UInt32Value): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RouteConfigurationOptions.AsObject;
  static toObject(includeInstance: boolean, msg: RouteConfigurationOptions): RouteConfigurationOptions.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RouteConfigurationOptions, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RouteConfigurationOptions;
  static deserializeBinaryFromReader(message: RouteConfigurationOptions, reader: jspb.BinaryReader): RouteConfigurationOptions;
}

export namespace RouteConfigurationOptions {
  export type AsObject = {
    maxDirectResponseBodySizeBytes?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
  }
}

export class HttpListenerOptions extends jspb.Message {
  hasGrpcWeb(): boolean;
  clearGrpcWeb(): void;
  getGrpcWeb(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_grpc_web_grpc_web_pb.GrpcWeb | undefined;
  setGrpcWeb(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_grpc_web_grpc_web_pb.GrpcWeb): void;

  hasHttpConnectionManagerSettings(): boolean;
  clearHttpConnectionManagerSettings(): void;
  getHttpConnectionManagerSettings(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_hcm_hcm_pb.HttpConnectionManagerSettings | undefined;
  setHttpConnectionManagerSettings(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_hcm_hcm_pb.HttpConnectionManagerSettings): void;

  hasHealthCheck(): boolean;
  clearHealthCheck(): void;
  getHealthCheck(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_healthcheck_healthcheck_pb.HealthCheck | undefined;
  setHealthCheck(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_healthcheck_healthcheck_pb.HealthCheck): void;

  hasExtensions(): boolean;
  clearExtensions(): void;
  getExtensions(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_extensions_pb.Extensions | undefined;
  setExtensions(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_extensions_pb.Extensions): void;

  hasWaf(): boolean;
  clearWaf(): void;
  getWaf(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_waf_waf_pb.Settings | undefined;
  setWaf(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_waf_waf_pb.Settings): void;

  hasDlp(): boolean;
  clearDlp(): void;
  getDlp(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_dlp_dlp_pb.FilterConfig | undefined;
  setDlp(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_dlp_dlp_pb.FilterConfig): void;

  hasWasm(): boolean;
  clearWasm(): void;
  getWasm(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_wasm_wasm_pb.PluginSource | undefined;
  setWasm(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_wasm_wasm_pb.PluginSource): void;

  hasExtauth(): boolean;
  clearExtauth(): void;
  getExtauth(): github_com_solo_io_solo_apis_api_gloo_enterprise_gloo_v1_auth_config_pb.Settings | undefined;
  setExtauth(value?: github_com_solo_io_solo_apis_api_gloo_enterprise_gloo_v1_auth_config_pb.Settings): void;

  hasRatelimitServer(): boolean;
  clearRatelimitServer(): void;
  getRatelimitServer(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_ratelimit_ratelimit_pb.Settings | undefined;
  setRatelimitServer(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_ratelimit_ratelimit_pb.Settings): void;

  hasCaching(): boolean;
  clearCaching(): void;
  getCaching(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_caching_caching_pb.Settings | undefined;
  setCaching(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_caching_caching_pb.Settings): void;

  hasGzip(): boolean;
  clearGzip(): void;
  getGzip(): github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_filter_http_gzip_v2_gzip_pb.Gzip | undefined;
  setGzip(value?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_filter_http_gzip_v2_gzip_pb.Gzip): void;

  hasProxyLatency(): boolean;
  clearProxyLatency(): void;
  getProxyLatency(): github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_extensions_proxylatency_proxylatency_pb.ProxyLatency | undefined;
  setProxyLatency(value?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_extensions_proxylatency_proxylatency_pb.ProxyLatency): void;

  hasBuffer(): boolean;
  clearBuffer(): void;
  getBuffer(): github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_extensions_filters_http_buffer_v3_buffer_pb.Buffer | undefined;
  setBuffer(value?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_extensions_filters_http_buffer_v3_buffer_pb.Buffer): void;

  hasCsrf(): boolean;
  clearCsrf(): void;
  getCsrf(): github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_extensions_filters_http_csrf_v3_csrf_pb.CsrfPolicy | undefined;
  setCsrf(value?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_extensions_filters_http_csrf_v3_csrf_pb.CsrfPolicy): void;

  hasGrpcJsonTranscoder(): boolean;
  clearGrpcJsonTranscoder(): void;
  getGrpcJsonTranscoder(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_grpc_json_grpc_json_pb.GrpcJsonTranscoder | undefined;
  setGrpcJsonTranscoder(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_grpc_json_grpc_json_pb.GrpcJsonTranscoder): void;

  hasSanitizeClusterHeader(): boolean;
  clearSanitizeClusterHeader(): void;
  getSanitizeClusterHeader(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setSanitizeClusterHeader(value?: google_protobuf_wrappers_pb.BoolValue): void;

  hasLeftmostXffAddress(): boolean;
  clearLeftmostXffAddress(): void;
  getLeftmostXffAddress(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setLeftmostXffAddress(value?: google_protobuf_wrappers_pb.BoolValue): void;

  hasDynamicForwardProxy(): boolean;
  clearDynamicForwardProxy(): void;
  getDynamicForwardProxy(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_dynamic_forward_proxy_dynamic_forward_proxy_pb.FilterConfig | undefined;
  setDynamicForwardProxy(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_dynamic_forward_proxy_dynamic_forward_proxy_pb.FilterConfig): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HttpListenerOptions.AsObject;
  static toObject(includeInstance: boolean, msg: HttpListenerOptions): HttpListenerOptions.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: HttpListenerOptions, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HttpListenerOptions;
  static deserializeBinaryFromReader(message: HttpListenerOptions, reader: jspb.BinaryReader): HttpListenerOptions;
}

export namespace HttpListenerOptions {
  export type AsObject = {
    grpcWeb?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_grpc_web_grpc_web_pb.GrpcWeb.AsObject,
    httpConnectionManagerSettings?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_hcm_hcm_pb.HttpConnectionManagerSettings.AsObject,
    healthCheck?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_healthcheck_healthcheck_pb.HealthCheck.AsObject,
    extensions?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_extensions_pb.Extensions.AsObject,
    waf?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_waf_waf_pb.Settings.AsObject,
    dlp?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_dlp_dlp_pb.FilterConfig.AsObject,
    wasm?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_wasm_wasm_pb.PluginSource.AsObject,
    extauth?: github_com_solo_io_solo_apis_api_gloo_enterprise_gloo_v1_auth_config_pb.Settings.AsObject,
    ratelimitServer?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_ratelimit_ratelimit_pb.Settings.AsObject,
    caching?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_caching_caching_pb.Settings.AsObject,
    gzip?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_filter_http_gzip_v2_gzip_pb.Gzip.AsObject,
    proxyLatency?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_extensions_proxylatency_proxylatency_pb.ProxyLatency.AsObject,
    buffer?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_extensions_filters_http_buffer_v3_buffer_pb.Buffer.AsObject,
    csrf?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_extensions_filters_http_csrf_v3_csrf_pb.CsrfPolicy.AsObject,
    grpcJsonTranscoder?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_grpc_json_grpc_json_pb.GrpcJsonTranscoder.AsObject,
    sanitizeClusterHeader?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    leftmostXffAddress?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    dynamicForwardProxy?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_dynamic_forward_proxy_dynamic_forward_proxy_pb.FilterConfig.AsObject,
  }
}

export class TcpListenerOptions extends jspb.Message {
  hasTcpProxySettings(): boolean;
  clearTcpProxySettings(): void;
  getTcpProxySettings(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_tcp_tcp_pb.TcpProxySettings | undefined;
  setTcpProxySettings(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_tcp_tcp_pb.TcpProxySettings): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TcpListenerOptions.AsObject;
  static toObject(includeInstance: boolean, msg: TcpListenerOptions): TcpListenerOptions.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: TcpListenerOptions, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TcpListenerOptions;
  static deserializeBinaryFromReader(message: TcpListenerOptions, reader: jspb.BinaryReader): TcpListenerOptions;
}

export namespace TcpListenerOptions {
  export type AsObject = {
    tcpProxySettings?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_tcp_tcp_pb.TcpProxySettings.AsObject,
  }
}

export class VirtualHostOptions extends jspb.Message {
  hasExtensions(): boolean;
  clearExtensions(): void;
  getExtensions(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_extensions_pb.Extensions | undefined;
  setExtensions(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_extensions_pb.Extensions): void;

  hasRetries(): boolean;
  clearRetries(): void;
  getRetries(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_retries_retries_pb.RetryPolicy | undefined;
  setRetries(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_retries_retries_pb.RetryPolicy): void;

  hasStats(): boolean;
  clearStats(): void;
  getStats(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_stats_stats_pb.Stats | undefined;
  setStats(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_stats_stats_pb.Stats): void;

  hasHeaderManipulation(): boolean;
  clearHeaderManipulation(): void;
  getHeaderManipulation(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_headers_headers_pb.HeaderManipulation | undefined;
  setHeaderManipulation(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_headers_headers_pb.HeaderManipulation): void;

  hasCors(): boolean;
  clearCors(): void;
  getCors(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_cors_cors_pb.CorsPolicy | undefined;
  setCors(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_cors_cors_pb.CorsPolicy): void;

  hasTransformations(): boolean;
  clearTransformations(): void;
  getTransformations(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_transformation_transformation_pb.Transformations | undefined;
  setTransformations(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_transformation_transformation_pb.Transformations): void;

  hasRatelimitBasic(): boolean;
  clearRatelimitBasic(): void;
  getRatelimitBasic(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_ratelimit_ratelimit_pb.IngressRateLimit | undefined;
  setRatelimitBasic(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_ratelimit_ratelimit_pb.IngressRateLimit): void;

  hasRatelimitEarly(): boolean;
  clearRatelimitEarly(): void;
  getRatelimitEarly(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_ratelimit_ratelimit_pb.RateLimitVhostExtension | undefined;
  setRatelimitEarly(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_ratelimit_ratelimit_pb.RateLimitVhostExtension): void;

  hasRateLimitEarlyConfigs(): boolean;
  clearRateLimitEarlyConfigs(): void;
  getRateLimitEarlyConfigs(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_ratelimit_ratelimit_pb.RateLimitConfigRefs | undefined;
  setRateLimitEarlyConfigs(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_ratelimit_ratelimit_pb.RateLimitConfigRefs): void;

  hasRatelimit(): boolean;
  clearRatelimit(): void;
  getRatelimit(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_ratelimit_ratelimit_pb.RateLimitVhostExtension | undefined;
  setRatelimit(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_ratelimit_ratelimit_pb.RateLimitVhostExtension): void;

  hasRateLimitConfigs(): boolean;
  clearRateLimitConfigs(): void;
  getRateLimitConfigs(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_ratelimit_ratelimit_pb.RateLimitConfigRefs | undefined;
  setRateLimitConfigs(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_ratelimit_ratelimit_pb.RateLimitConfigRefs): void;

  hasRatelimitRegular(): boolean;
  clearRatelimitRegular(): void;
  getRatelimitRegular(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_ratelimit_ratelimit_pb.RateLimitVhostExtension | undefined;
  setRatelimitRegular(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_ratelimit_ratelimit_pb.RateLimitVhostExtension): void;

  hasRateLimitRegularConfigs(): boolean;
  clearRateLimitRegularConfigs(): void;
  getRateLimitRegularConfigs(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_ratelimit_ratelimit_pb.RateLimitConfigRefs | undefined;
  setRateLimitRegularConfigs(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_ratelimit_ratelimit_pb.RateLimitConfigRefs): void;

  hasWaf(): boolean;
  clearWaf(): void;
  getWaf(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_waf_waf_pb.Settings | undefined;
  setWaf(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_waf_waf_pb.Settings): void;

  hasJwt(): boolean;
  clearJwt(): void;
  getJwt(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_jwt_jwt_pb.VhostExtension | undefined;
  setJwt(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_jwt_jwt_pb.VhostExtension): void;

  hasJwtStaged(): boolean;
  clearJwtStaged(): void;
  getJwtStaged(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_jwt_jwt_pb.JwtStagedVhostExtension | undefined;
  setJwtStaged(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_jwt_jwt_pb.JwtStagedVhostExtension): void;

  hasRbac(): boolean;
  clearRbac(): void;
  getRbac(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_rbac_rbac_pb.ExtensionSettings | undefined;
  setRbac(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_rbac_rbac_pb.ExtensionSettings): void;

  hasExtauth(): boolean;
  clearExtauth(): void;
  getExtauth(): github_com_solo_io_solo_apis_api_gloo_enterprise_gloo_v1_auth_config_pb.ExtAuthExtension | undefined;
  setExtauth(value?: github_com_solo_io_solo_apis_api_gloo_enterprise_gloo_v1_auth_config_pb.ExtAuthExtension): void;

  hasDlp(): boolean;
  clearDlp(): void;
  getDlp(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_dlp_dlp_pb.Config | undefined;
  setDlp(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_dlp_dlp_pb.Config): void;

  hasBufferPerRoute(): boolean;
  clearBufferPerRoute(): void;
  getBufferPerRoute(): github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_extensions_filters_http_buffer_v3_buffer_pb.BufferPerRoute | undefined;
  setBufferPerRoute(value?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_extensions_filters_http_buffer_v3_buffer_pb.BufferPerRoute): void;

  hasCsrf(): boolean;
  clearCsrf(): void;
  getCsrf(): github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_extensions_filters_http_csrf_v3_csrf_pb.CsrfPolicy | undefined;
  setCsrf(value?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_extensions_filters_http_csrf_v3_csrf_pb.CsrfPolicy): void;

  hasIncludeRequestAttemptCount(): boolean;
  clearIncludeRequestAttemptCount(): void;
  getIncludeRequestAttemptCount(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setIncludeRequestAttemptCount(value?: google_protobuf_wrappers_pb.BoolValue): void;

  hasIncludeAttemptCountInResponse(): boolean;
  clearIncludeAttemptCountInResponse(): void;
  getIncludeAttemptCountInResponse(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setIncludeAttemptCountInResponse(value?: google_protobuf_wrappers_pb.BoolValue): void;

  hasStagedTransformations(): boolean;
  clearStagedTransformations(): void;
  getStagedTransformations(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_transformation_transformation_pb.TransformationStages | undefined;
  setStagedTransformations(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_transformation_transformation_pb.TransformationStages): void;

  getRateLimitEarlyConfigTypeCase(): VirtualHostOptions.RateLimitEarlyConfigTypeCase;
  getRateLimitConfigTypeCase(): VirtualHostOptions.RateLimitConfigTypeCase;
  getRateLimitRegularConfigTypeCase(): VirtualHostOptions.RateLimitRegularConfigTypeCase;
  getJwtConfigCase(): VirtualHostOptions.JwtConfigCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VirtualHostOptions.AsObject;
  static toObject(includeInstance: boolean, msg: VirtualHostOptions): VirtualHostOptions.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: VirtualHostOptions, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VirtualHostOptions;
  static deserializeBinaryFromReader(message: VirtualHostOptions, reader: jspb.BinaryReader): VirtualHostOptions;
}

export namespace VirtualHostOptions {
  export type AsObject = {
    extensions?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_extensions_pb.Extensions.AsObject,
    retries?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_retries_retries_pb.RetryPolicy.AsObject,
    stats?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_stats_stats_pb.Stats.AsObject,
    headerManipulation?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_headers_headers_pb.HeaderManipulation.AsObject,
    cors?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_cors_cors_pb.CorsPolicy.AsObject,
    transformations?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_transformation_transformation_pb.Transformations.AsObject,
    ratelimitBasic?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_ratelimit_ratelimit_pb.IngressRateLimit.AsObject,
    ratelimitEarly?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_ratelimit_ratelimit_pb.RateLimitVhostExtension.AsObject,
    rateLimitEarlyConfigs?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_ratelimit_ratelimit_pb.RateLimitConfigRefs.AsObject,
    ratelimit?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_ratelimit_ratelimit_pb.RateLimitVhostExtension.AsObject,
    rateLimitConfigs?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_ratelimit_ratelimit_pb.RateLimitConfigRefs.AsObject,
    ratelimitRegular?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_ratelimit_ratelimit_pb.RateLimitVhostExtension.AsObject,
    rateLimitRegularConfigs?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_ratelimit_ratelimit_pb.RateLimitConfigRefs.AsObject,
    waf?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_waf_waf_pb.Settings.AsObject,
    jwt?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_jwt_jwt_pb.VhostExtension.AsObject,
    jwtStaged?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_jwt_jwt_pb.JwtStagedVhostExtension.AsObject,
    rbac?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_rbac_rbac_pb.ExtensionSettings.AsObject,
    extauth?: github_com_solo_io_solo_apis_api_gloo_enterprise_gloo_v1_auth_config_pb.ExtAuthExtension.AsObject,
    dlp?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_dlp_dlp_pb.Config.AsObject,
    bufferPerRoute?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_extensions_filters_http_buffer_v3_buffer_pb.BufferPerRoute.AsObject,
    csrf?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_extensions_filters_http_csrf_v3_csrf_pb.CsrfPolicy.AsObject,
    includeRequestAttemptCount?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    includeAttemptCountInResponse?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    stagedTransformations?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_transformation_transformation_pb.TransformationStages.AsObject,
  }

  export enum RateLimitEarlyConfigTypeCase {
    RATE_LIMIT_EARLY_CONFIG_TYPE_NOT_SET = 0,
    RATELIMIT_EARLY = 72,
    RATE_LIMIT_EARLY_CONFIGS = 73,
  }

  export enum RateLimitConfigTypeCase {
    RATE_LIMIT_CONFIG_TYPE_NOT_SET = 0,
    RATELIMIT = 70,
    RATE_LIMIT_CONFIGS = 71,
  }

  export enum RateLimitRegularConfigTypeCase {
    RATE_LIMIT_REGULAR_CONFIG_TYPE_NOT_SET = 0,
    RATELIMIT_REGULAR = 74,
    RATE_LIMIT_REGULAR_CONFIGS = 75,
  }

  export enum JwtConfigCase {
    JWT_CONFIG_NOT_SET = 0,
    JWT = 9,
    JWT_STAGED = 19,
  }
}

export class RouteOptions extends jspb.Message {
  hasTransformations(): boolean;
  clearTransformations(): void;
  getTransformations(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_transformation_transformation_pb.Transformations | undefined;
  setTransformations(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_transformation_transformation_pb.Transformations): void;

  hasFaults(): boolean;
  clearFaults(): void;
  getFaults(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_faultinjection_fault_pb.RouteFaults | undefined;
  setFaults(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_faultinjection_fault_pb.RouteFaults): void;

  hasPrefixRewrite(): boolean;
  clearPrefixRewrite(): void;
  getPrefixRewrite(): google_protobuf_wrappers_pb.StringValue | undefined;
  setPrefixRewrite(value?: google_protobuf_wrappers_pb.StringValue): void;

  hasTimeout(): boolean;
  clearTimeout(): void;
  getTimeout(): google_protobuf_duration_pb.Duration | undefined;
  setTimeout(value?: google_protobuf_duration_pb.Duration): void;

  hasRetries(): boolean;
  clearRetries(): void;
  getRetries(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_retries_retries_pb.RetryPolicy | undefined;
  setRetries(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_retries_retries_pb.RetryPolicy): void;

  hasExtensions(): boolean;
  clearExtensions(): void;
  getExtensions(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_extensions_pb.Extensions | undefined;
  setExtensions(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_extensions_pb.Extensions): void;

  hasTracing(): boolean;
  clearTracing(): void;
  getTracing(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_tracing_tracing_pb.RouteTracingSettings | undefined;
  setTracing(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_tracing_tracing_pb.RouteTracingSettings): void;

  hasShadowing(): boolean;
  clearShadowing(): void;
  getShadowing(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_shadowing_shadowing_pb.RouteShadowing | undefined;
  setShadowing(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_shadowing_shadowing_pb.RouteShadowing): void;

  hasHeaderManipulation(): boolean;
  clearHeaderManipulation(): void;
  getHeaderManipulation(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_headers_headers_pb.HeaderManipulation | undefined;
  setHeaderManipulation(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_headers_headers_pb.HeaderManipulation): void;

  hasHostRewrite(): boolean;
  clearHostRewrite(): void;
  getHostRewrite(): string;
  setHostRewrite(value: string): void;

  hasAutoHostRewrite(): boolean;
  clearAutoHostRewrite(): void;
  getAutoHostRewrite(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setAutoHostRewrite(value?: google_protobuf_wrappers_pb.BoolValue): void;

  hasCors(): boolean;
  clearCors(): void;
  getCors(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_cors_cors_pb.CorsPolicy | undefined;
  setCors(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_cors_cors_pb.CorsPolicy): void;

  hasLbHash(): boolean;
  clearLbHash(): void;
  getLbHash(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_lbhash_lbhash_pb.RouteActionHashConfig | undefined;
  setLbHash(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_lbhash_lbhash_pb.RouteActionHashConfig): void;

  clearUpgradesList(): void;
  getUpgradesList(): Array<github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_protocol_upgrade_protocol_upgrade_pb.ProtocolUpgradeConfig>;
  setUpgradesList(value: Array<github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_protocol_upgrade_protocol_upgrade_pb.ProtocolUpgradeConfig>): void;
  addUpgrades(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_protocol_upgrade_protocol_upgrade_pb.ProtocolUpgradeConfig, index?: number): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_protocol_upgrade_protocol_upgrade_pb.ProtocolUpgradeConfig;

  hasRatelimitBasic(): boolean;
  clearRatelimitBasic(): void;
  getRatelimitBasic(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_ratelimit_ratelimit_pb.IngressRateLimit | undefined;
  setRatelimitBasic(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_ratelimit_ratelimit_pb.IngressRateLimit): void;

  hasRatelimitEarly(): boolean;
  clearRatelimitEarly(): void;
  getRatelimitEarly(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_ratelimit_ratelimit_pb.RateLimitRouteExtension | undefined;
  setRatelimitEarly(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_ratelimit_ratelimit_pb.RateLimitRouteExtension): void;

  hasRateLimitEarlyConfigs(): boolean;
  clearRateLimitEarlyConfigs(): void;
  getRateLimitEarlyConfigs(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_ratelimit_ratelimit_pb.RateLimitConfigRefs | undefined;
  setRateLimitEarlyConfigs(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_ratelimit_ratelimit_pb.RateLimitConfigRefs): void;

  hasRatelimit(): boolean;
  clearRatelimit(): void;
  getRatelimit(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_ratelimit_ratelimit_pb.RateLimitRouteExtension | undefined;
  setRatelimit(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_ratelimit_ratelimit_pb.RateLimitRouteExtension): void;

  hasRateLimitConfigs(): boolean;
  clearRateLimitConfigs(): void;
  getRateLimitConfigs(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_ratelimit_ratelimit_pb.RateLimitConfigRefs | undefined;
  setRateLimitConfigs(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_ratelimit_ratelimit_pb.RateLimitConfigRefs): void;

  hasRatelimitRegular(): boolean;
  clearRatelimitRegular(): void;
  getRatelimitRegular(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_ratelimit_ratelimit_pb.RateLimitRouteExtension | undefined;
  setRatelimitRegular(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_ratelimit_ratelimit_pb.RateLimitRouteExtension): void;

  hasRateLimitRegularConfigs(): boolean;
  clearRateLimitRegularConfigs(): void;
  getRateLimitRegularConfigs(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_ratelimit_ratelimit_pb.RateLimitConfigRefs | undefined;
  setRateLimitRegularConfigs(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_ratelimit_ratelimit_pb.RateLimitConfigRefs): void;

  hasWaf(): boolean;
  clearWaf(): void;
  getWaf(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_waf_waf_pb.Settings | undefined;
  setWaf(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_waf_waf_pb.Settings): void;

  hasJwt(): boolean;
  clearJwt(): void;
  getJwt(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_jwt_jwt_pb.RouteExtension | undefined;
  setJwt(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_jwt_jwt_pb.RouteExtension): void;

  hasJwtStaged(): boolean;
  clearJwtStaged(): void;
  getJwtStaged(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_jwt_jwt_pb.JwtStagedRouteExtension | undefined;
  setJwtStaged(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_jwt_jwt_pb.JwtStagedRouteExtension): void;

  hasRbac(): boolean;
  clearRbac(): void;
  getRbac(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_rbac_rbac_pb.ExtensionSettings | undefined;
  setRbac(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_rbac_rbac_pb.ExtensionSettings): void;

  hasExtauth(): boolean;
  clearExtauth(): void;
  getExtauth(): github_com_solo_io_solo_apis_api_gloo_enterprise_gloo_v1_auth_config_pb.ExtAuthExtension | undefined;
  setExtauth(value?: github_com_solo_io_solo_apis_api_gloo_enterprise_gloo_v1_auth_config_pb.ExtAuthExtension): void;

  hasDlp(): boolean;
  clearDlp(): void;
  getDlp(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_dlp_dlp_pb.Config | undefined;
  setDlp(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_dlp_dlp_pb.Config): void;

  hasBufferPerRoute(): boolean;
  clearBufferPerRoute(): void;
  getBufferPerRoute(): github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_extensions_filters_http_buffer_v3_buffer_pb.BufferPerRoute | undefined;
  setBufferPerRoute(value?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_extensions_filters_http_buffer_v3_buffer_pb.BufferPerRoute): void;

  hasCsrf(): boolean;
  clearCsrf(): void;
  getCsrf(): github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_extensions_filters_http_csrf_v3_csrf_pb.CsrfPolicy | undefined;
  setCsrf(value?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_extensions_filters_http_csrf_v3_csrf_pb.CsrfPolicy): void;

  hasStagedTransformations(): boolean;
  clearStagedTransformations(): void;
  getStagedTransformations(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_transformation_transformation_pb.TransformationStages | undefined;
  setStagedTransformations(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_transformation_transformation_pb.TransformationStages): void;

  getEnvoyMetadataMap(): jspb.Map<string, google_protobuf_struct_pb.Struct>;
  clearEnvoyMetadataMap(): void;
  hasRegexRewrite(): boolean;
  clearRegexRewrite(): void;
  getRegexRewrite(): github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_type_matcher_v3_regex_pb.RegexMatchAndSubstitute | undefined;
  setRegexRewrite(value?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_type_matcher_v3_regex_pb.RegexMatchAndSubstitute): void;

  getHostRewriteTypeCase(): RouteOptions.HostRewriteTypeCase;
  getRateLimitEarlyConfigTypeCase(): RouteOptions.RateLimitEarlyConfigTypeCase;
  getRateLimitConfigTypeCase(): RouteOptions.RateLimitConfigTypeCase;
  getRateLimitRegularConfigTypeCase(): RouteOptions.RateLimitRegularConfigTypeCase;
  getJwtConfigCase(): RouteOptions.JwtConfigCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RouteOptions.AsObject;
  static toObject(includeInstance: boolean, msg: RouteOptions): RouteOptions.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RouteOptions, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RouteOptions;
  static deserializeBinaryFromReader(message: RouteOptions, reader: jspb.BinaryReader): RouteOptions;
}

export namespace RouteOptions {
  export type AsObject = {
    transformations?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_transformation_transformation_pb.Transformations.AsObject,
    faults?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_faultinjection_fault_pb.RouteFaults.AsObject,
    prefixRewrite?: google_protobuf_wrappers_pb.StringValue.AsObject,
    timeout?: google_protobuf_duration_pb.Duration.AsObject,
    retries?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_retries_retries_pb.RetryPolicy.AsObject,
    extensions?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_extensions_pb.Extensions.AsObject,
    tracing?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_tracing_tracing_pb.RouteTracingSettings.AsObject,
    shadowing?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_shadowing_shadowing_pb.RouteShadowing.AsObject,
    headerManipulation?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_headers_headers_pb.HeaderManipulation.AsObject,
    hostRewrite: string,
    autoHostRewrite?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    cors?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_cors_cors_pb.CorsPolicy.AsObject,
    lbHash?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_lbhash_lbhash_pb.RouteActionHashConfig.AsObject,
    upgradesList: Array<github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_protocol_upgrade_protocol_upgrade_pb.ProtocolUpgradeConfig.AsObject>,
    ratelimitBasic?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_ratelimit_ratelimit_pb.IngressRateLimit.AsObject,
    ratelimitEarly?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_ratelimit_ratelimit_pb.RateLimitRouteExtension.AsObject,
    rateLimitEarlyConfigs?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_ratelimit_ratelimit_pb.RateLimitConfigRefs.AsObject,
    ratelimit?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_ratelimit_ratelimit_pb.RateLimitRouteExtension.AsObject,
    rateLimitConfigs?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_ratelimit_ratelimit_pb.RateLimitConfigRefs.AsObject,
    ratelimitRegular?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_ratelimit_ratelimit_pb.RateLimitRouteExtension.AsObject,
    rateLimitRegularConfigs?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_ratelimit_ratelimit_pb.RateLimitConfigRefs.AsObject,
    waf?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_waf_waf_pb.Settings.AsObject,
    jwt?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_jwt_jwt_pb.RouteExtension.AsObject,
    jwtStaged?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_jwt_jwt_pb.JwtStagedRouteExtension.AsObject,
    rbac?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_rbac_rbac_pb.ExtensionSettings.AsObject,
    extauth?: github_com_solo_io_solo_apis_api_gloo_enterprise_gloo_v1_auth_config_pb.ExtAuthExtension.AsObject,
    dlp?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_dlp_dlp_pb.Config.AsObject,
    bufferPerRoute?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_extensions_filters_http_buffer_v3_buffer_pb.BufferPerRoute.AsObject,
    csrf?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_extensions_filters_http_csrf_v3_csrf_pb.CsrfPolicy.AsObject,
    stagedTransformations?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_transformation_transformation_pb.TransformationStages.AsObject,
    envoyMetadataMap: Array<[string, google_protobuf_struct_pb.Struct.AsObject]>,
    regexRewrite?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_type_matcher_v3_regex_pb.RegexMatchAndSubstitute.AsObject,
  }

  export enum HostRewriteTypeCase {
    HOST_REWRITE_TYPE_NOT_SET = 0,
    HOST_REWRITE = 10,
    AUTO_HOST_REWRITE = 19,
  }

  export enum RateLimitEarlyConfigTypeCase {
    RATE_LIMIT_EARLY_CONFIG_TYPE_NOT_SET = 0,
    RATELIMIT_EARLY = 142,
    RATE_LIMIT_EARLY_CONFIGS = 143,
  }

  export enum RateLimitConfigTypeCase {
    RATE_LIMIT_CONFIG_TYPE_NOT_SET = 0,
    RATELIMIT = 140,
    RATE_LIMIT_CONFIGS = 141,
  }

  export enum RateLimitRegularConfigTypeCase {
    RATE_LIMIT_REGULAR_CONFIG_TYPE_NOT_SET = 0,
    RATELIMIT_REGULAR = 144,
    RATE_LIMIT_REGULAR_CONFIGS = 145,
  }

  export enum JwtConfigCase {
    JWT_CONFIG_NOT_SET = 0,
    JWT = 16,
    JWT_STAGED = 25,
  }
}

export class DestinationSpec extends jspb.Message {
  hasAws(): boolean;
  clearAws(): void;
  getAws(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_aws_aws_pb.DestinationSpec | undefined;
  setAws(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_aws_aws_pb.DestinationSpec): void;

  hasAzure(): boolean;
  clearAzure(): void;
  getAzure(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_azure_azure_pb.DestinationSpec | undefined;
  setAzure(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_azure_azure_pb.DestinationSpec): void;

  hasRest(): boolean;
  clearRest(): void;
  getRest(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_rest_rest_pb.DestinationSpec | undefined;
  setRest(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_rest_rest_pb.DestinationSpec): void;

  hasGrpc(): boolean;
  clearGrpc(): void;
  getGrpc(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_grpc_grpc_pb.DestinationSpec | undefined;
  setGrpc(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_grpc_grpc_pb.DestinationSpec): void;

  getDestinationTypeCase(): DestinationSpec.DestinationTypeCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DestinationSpec.AsObject;
  static toObject(includeInstance: boolean, msg: DestinationSpec): DestinationSpec.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DestinationSpec, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DestinationSpec;
  static deserializeBinaryFromReader(message: DestinationSpec, reader: jspb.BinaryReader): DestinationSpec;
}

export namespace DestinationSpec {
  export type AsObject = {
    aws?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_aws_aws_pb.DestinationSpec.AsObject,
    azure?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_azure_azure_pb.DestinationSpec.AsObject,
    rest?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_rest_rest_pb.DestinationSpec.AsObject,
    grpc?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_grpc_grpc_pb.DestinationSpec.AsObject,
  }

  export enum DestinationTypeCase {
    DESTINATION_TYPE_NOT_SET = 0,
    AWS = 1,
    AZURE = 2,
    REST = 3,
    GRPC = 4,
  }
}

export class WeightedDestinationOptions extends jspb.Message {
  hasHeaderManipulation(): boolean;
  clearHeaderManipulation(): void;
  getHeaderManipulation(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_headers_headers_pb.HeaderManipulation | undefined;
  setHeaderManipulation(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_headers_headers_pb.HeaderManipulation): void;

  hasTransformations(): boolean;
  clearTransformations(): void;
  getTransformations(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_transformation_transformation_pb.Transformations | undefined;
  setTransformations(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_transformation_transformation_pb.Transformations): void;

  hasExtensions(): boolean;
  clearExtensions(): void;
  getExtensions(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_extensions_pb.Extensions | undefined;
  setExtensions(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_extensions_pb.Extensions): void;

  hasExtauth(): boolean;
  clearExtauth(): void;
  getExtauth(): github_com_solo_io_solo_apis_api_gloo_enterprise_gloo_v1_auth_config_pb.ExtAuthExtension | undefined;
  setExtauth(value?: github_com_solo_io_solo_apis_api_gloo_enterprise_gloo_v1_auth_config_pb.ExtAuthExtension): void;

  hasBufferPerRoute(): boolean;
  clearBufferPerRoute(): void;
  getBufferPerRoute(): github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_extensions_filters_http_buffer_v3_buffer_pb.BufferPerRoute | undefined;
  setBufferPerRoute(value?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_extensions_filters_http_buffer_v3_buffer_pb.BufferPerRoute): void;

  hasCsrf(): boolean;
  clearCsrf(): void;
  getCsrf(): github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_extensions_filters_http_csrf_v3_csrf_pb.CsrfPolicy | undefined;
  setCsrf(value?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_extensions_filters_http_csrf_v3_csrf_pb.CsrfPolicy): void;

  hasStagedTransformations(): boolean;
  clearStagedTransformations(): void;
  getStagedTransformations(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_transformation_transformation_pb.TransformationStages | undefined;
  setStagedTransformations(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_transformation_transformation_pb.TransformationStages): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): WeightedDestinationOptions.AsObject;
  static toObject(includeInstance: boolean, msg: WeightedDestinationOptions): WeightedDestinationOptions.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: WeightedDestinationOptions, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): WeightedDestinationOptions;
  static deserializeBinaryFromReader(message: WeightedDestinationOptions, reader: jspb.BinaryReader): WeightedDestinationOptions;
}

export namespace WeightedDestinationOptions {
  export type AsObject = {
    headerManipulation?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_headers_headers_pb.HeaderManipulation.AsObject,
    transformations?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_transformation_transformation_pb.Transformations.AsObject,
    extensions?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_extensions_pb.Extensions.AsObject,
    extauth?: github_com_solo_io_solo_apis_api_gloo_enterprise_gloo_v1_auth_config_pb.ExtAuthExtension.AsObject,
    bufferPerRoute?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_extensions_filters_http_buffer_v3_buffer_pb.BufferPerRoute.AsObject,
    csrf?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_extensions_filters_http_csrf_v3_csrf_pb.CsrfPolicy.AsObject,
    stagedTransformations?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_transformation_transformation_pb.TransformationStages.AsObject,
  }
}
