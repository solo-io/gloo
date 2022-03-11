/* eslint-disable */
// package: dfp.options.gloo.solo.io
// file: github.com/solo-io/solo-apis/api/gloo/gloo/v1/options/dynamic_forward_proxy/dynamic_forward_proxy.proto

import * as jspb from "google-protobuf";
import * as extproto_ext_pb from "../../../../../../../../../extproto/ext_pb";
import * as validate_validate_pb from "../../../../../../../../../validate/validate_pb";
import * as google_protobuf_duration_pb from "google-protobuf/google/protobuf/duration_pb";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_address_pb from "../../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/config/core/v3/address_pb";

export class FilterConfig extends jspb.Message {
  hasDnsCacheConfig(): boolean;
  clearDnsCacheConfig(): void;
  getDnsCacheConfig(): DnsCacheConfig | undefined;
  setDnsCacheConfig(value?: DnsCacheConfig): void;

  getSaveUpstreamAddress(): boolean;
  setSaveUpstreamAddress(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FilterConfig.AsObject;
  static toObject(includeInstance: boolean, msg: FilterConfig): FilterConfig.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: FilterConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FilterConfig;
  static deserializeBinaryFromReader(message: FilterConfig, reader: jspb.BinaryReader): FilterConfig;
}

export namespace FilterConfig {
  export type AsObject = {
    dnsCacheConfig?: DnsCacheConfig.AsObject,
    saveUpstreamAddress: boolean,
  }
}

export class DnsCacheCircuitBreakers extends jspb.Message {
  hasMaxPendingRequests(): boolean;
  clearMaxPendingRequests(): void;
  getMaxPendingRequests(): google_protobuf_wrappers_pb.UInt32Value | undefined;
  setMaxPendingRequests(value?: google_protobuf_wrappers_pb.UInt32Value): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DnsCacheCircuitBreakers.AsObject;
  static toObject(includeInstance: boolean, msg: DnsCacheCircuitBreakers): DnsCacheCircuitBreakers.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DnsCacheCircuitBreakers, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DnsCacheCircuitBreakers;
  static deserializeBinaryFromReader(message: DnsCacheCircuitBreakers, reader: jspb.BinaryReader): DnsCacheCircuitBreakers;
}

export namespace DnsCacheCircuitBreakers {
  export type AsObject = {
    maxPendingRequests?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
  }
}

export class DnsCacheConfig extends jspb.Message {
  getDnsLookupFamily(): DnsLookupFamilyMap[keyof DnsLookupFamilyMap];
  setDnsLookupFamily(value: DnsLookupFamilyMap[keyof DnsLookupFamilyMap]): void;

  hasDnsRefreshRate(): boolean;
  clearDnsRefreshRate(): void;
  getDnsRefreshRate(): google_protobuf_duration_pb.Duration | undefined;
  setDnsRefreshRate(value?: google_protobuf_duration_pb.Duration): void;

  hasHostTtl(): boolean;
  clearHostTtl(): void;
  getHostTtl(): google_protobuf_duration_pb.Duration | undefined;
  setHostTtl(value?: google_protobuf_duration_pb.Duration): void;

  hasMaxHosts(): boolean;
  clearMaxHosts(): void;
  getMaxHosts(): google_protobuf_wrappers_pb.UInt32Value | undefined;
  setMaxHosts(value?: google_protobuf_wrappers_pb.UInt32Value): void;

  hasDnsFailureRefreshRate(): boolean;
  clearDnsFailureRefreshRate(): void;
  getDnsFailureRefreshRate(): RefreshRate | undefined;
  setDnsFailureRefreshRate(value?: RefreshRate): void;

  hasDnsCacheCircuitBreaker(): boolean;
  clearDnsCacheCircuitBreaker(): void;
  getDnsCacheCircuitBreaker(): DnsCacheCircuitBreakers | undefined;
  setDnsCacheCircuitBreaker(value?: DnsCacheCircuitBreakers): void;

  hasCaresDns(): boolean;
  clearCaresDns(): void;
  getCaresDns(): CaresDnsResolverConfig | undefined;
  setCaresDns(value?: CaresDnsResolverConfig): void;

  hasAppleDns(): boolean;
  clearAppleDns(): void;
  getAppleDns(): AppleDnsResolverConfig | undefined;
  setAppleDns(value?: AppleDnsResolverConfig): void;

  clearPreresolveHostnamesList(): void;
  getPreresolveHostnamesList(): Array<github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_address_pb.SocketAddress>;
  setPreresolveHostnamesList(value: Array<github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_address_pb.SocketAddress>): void;
  addPreresolveHostnames(value?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_address_pb.SocketAddress, index?: number): github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_address_pb.SocketAddress;

  hasDnsQueryTimeout(): boolean;
  clearDnsQueryTimeout(): void;
  getDnsQueryTimeout(): google_protobuf_duration_pb.Duration | undefined;
  setDnsQueryTimeout(value?: google_protobuf_duration_pb.Duration): void;

  getDnscachetypeCase(): DnsCacheConfig.DnscachetypeCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DnsCacheConfig.AsObject;
  static toObject(includeInstance: boolean, msg: DnsCacheConfig): DnsCacheConfig.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DnsCacheConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DnsCacheConfig;
  static deserializeBinaryFromReader(message: DnsCacheConfig, reader: jspb.BinaryReader): DnsCacheConfig;
}

export namespace DnsCacheConfig {
  export type AsObject = {
    dnsLookupFamily: DnsLookupFamilyMap[keyof DnsLookupFamilyMap],
    dnsRefreshRate?: google_protobuf_duration_pb.Duration.AsObject,
    hostTtl?: google_protobuf_duration_pb.Duration.AsObject,
    maxHosts?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
    dnsFailureRefreshRate?: RefreshRate.AsObject,
    dnsCacheCircuitBreaker?: DnsCacheCircuitBreakers.AsObject,
    caresDns?: CaresDnsResolverConfig.AsObject,
    appleDns?: AppleDnsResolverConfig.AsObject,
    preresolveHostnamesList: Array<github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_address_pb.SocketAddress.AsObject>,
    dnsQueryTimeout?: google_protobuf_duration_pb.Duration.AsObject,
  }

  export enum DnscachetypeCase {
    DNSCACHETYPE_NOT_SET = 0,
    CARES_DNS = 8,
    APPLE_DNS = 9,
  }
}

export class RefreshRate extends jspb.Message {
  hasBaseInterval(): boolean;
  clearBaseInterval(): void;
  getBaseInterval(): google_protobuf_duration_pb.Duration | undefined;
  setBaseInterval(value?: google_protobuf_duration_pb.Duration): void;

  hasMaxInterval(): boolean;
  clearMaxInterval(): void;
  getMaxInterval(): google_protobuf_duration_pb.Duration | undefined;
  setMaxInterval(value?: google_protobuf_duration_pb.Duration): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RefreshRate.AsObject;
  static toObject(includeInstance: boolean, msg: RefreshRate): RefreshRate.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RefreshRate, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RefreshRate;
  static deserializeBinaryFromReader(message: RefreshRate, reader: jspb.BinaryReader): RefreshRate;
}

export namespace RefreshRate {
  export type AsObject = {
    baseInterval?: google_protobuf_duration_pb.Duration.AsObject,
    maxInterval?: google_protobuf_duration_pb.Duration.AsObject,
  }
}

export class PerRouteConfig extends jspb.Message {
  hasHostRewrite(): boolean;
  clearHostRewrite(): void;
  getHostRewrite(): string;
  setHostRewrite(value: string): void;

  hasAutoHostRewriteHeader(): boolean;
  clearAutoHostRewriteHeader(): void;
  getAutoHostRewriteHeader(): string;
  setAutoHostRewriteHeader(value: string): void;

  getHostRewriteSpecifierCase(): PerRouteConfig.HostRewriteSpecifierCase;
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
    hostRewrite: string,
    autoHostRewriteHeader: string,
  }

  export enum HostRewriteSpecifierCase {
    HOST_REWRITE_SPECIFIER_NOT_SET = 0,
    HOST_REWRITE = 1,
    AUTO_HOST_REWRITE_HEADER = 2,
  }
}

export class DnsResolverOptions extends jspb.Message {
  getUseTcpForDnsLookups(): boolean;
  setUseTcpForDnsLookups(value: boolean): void;

  getNoDefaultSearchDomain(): boolean;
  setNoDefaultSearchDomain(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DnsResolverOptions.AsObject;
  static toObject(includeInstance: boolean, msg: DnsResolverOptions): DnsResolverOptions.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DnsResolverOptions, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DnsResolverOptions;
  static deserializeBinaryFromReader(message: DnsResolverOptions, reader: jspb.BinaryReader): DnsResolverOptions;
}

export namespace DnsResolverOptions {
  export type AsObject = {
    useTcpForDnsLookups: boolean,
    noDefaultSearchDomain: boolean,
  }
}

export class CaresDnsResolverConfig extends jspb.Message {
  clearResolversList(): void;
  getResolversList(): Array<github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_address_pb.Address>;
  setResolversList(value: Array<github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_address_pb.Address>): void;
  addResolvers(value?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_address_pb.Address, index?: number): github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_address_pb.Address;

  hasDnsResolverOptions(): boolean;
  clearDnsResolverOptions(): void;
  getDnsResolverOptions(): DnsResolverOptions | undefined;
  setDnsResolverOptions(value?: DnsResolverOptions): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CaresDnsResolverConfig.AsObject;
  static toObject(includeInstance: boolean, msg: CaresDnsResolverConfig): CaresDnsResolverConfig.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CaresDnsResolverConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CaresDnsResolverConfig;
  static deserializeBinaryFromReader(message: CaresDnsResolverConfig, reader: jspb.BinaryReader): CaresDnsResolverConfig;
}

export namespace CaresDnsResolverConfig {
  export type AsObject = {
    resolversList: Array<github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_address_pb.Address.AsObject>,
    dnsResolverOptions?: DnsResolverOptions.AsObject,
  }
}

export class AppleDnsResolverConfig extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AppleDnsResolverConfig.AsObject;
  static toObject(includeInstance: boolean, msg: AppleDnsResolverConfig): AppleDnsResolverConfig.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: AppleDnsResolverConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AppleDnsResolverConfig;
  static deserializeBinaryFromReader(message: AppleDnsResolverConfig, reader: jspb.BinaryReader): AppleDnsResolverConfig;
}

export namespace AppleDnsResolverConfig {
  export type AsObject = {
  }
}

export interface DnsLookupFamilyMap {
  V4_PREFERRED: 0;
  V4_ONLY: 1;
  V6_ONLY: 2;
  AUTO: 3;
  ALL: 4;
}

export const DnsLookupFamily: DnsLookupFamilyMap;
