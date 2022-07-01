/* eslint-disable */
// package: gloo.solo.io
// file: github.com/solo-io/solo-apis/api/gloo/gloo/v1/proxy.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";
import * as google_protobuf_empty_pb from "google-protobuf/google/protobuf/empty_pb";
import * as google_protobuf_struct_pb from "google-protobuf/google/protobuf/struct_pb";
import * as extproto_ext_pb from "../../../../../../../extproto/ext_pb";
import * as github_com_solo_io_solo_kit_api_v1_ref_pb from "../../../../../../../github.com/solo-io/solo-kit/api/v1/ref_pb";
import * as github_com_solo_io_solo_kit_api_v1_solo_kit_pb from "../../../../../../../github.com/solo-io/solo-kit/api/v1/solo-kit_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_ssl_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/ssl_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_subset_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/subset_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/options_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_core_matchers_matchers_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/core/matchers/matchers_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_address_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/config/core/v3/address_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_dynamic_forward_proxy_dynamic_forward_proxy_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/options/dynamic_forward_proxy/dynamic_forward_proxy_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_type_matcher_v3_regex_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/type/matcher/v3/regex_pb";

export class ProxySpec extends jspb.Message {
  getCompressedspec(): string;
  setCompressedspec(value: string): void;

  clearListenersList(): void;
  getListenersList(): Array<Listener>;
  setListenersList(value: Array<Listener>): void;
  addListeners(value?: Listener, index?: number): Listener;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProxySpec.AsObject;
  static toObject(includeInstance: boolean, msg: ProxySpec): ProxySpec.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ProxySpec, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProxySpec;
  static deserializeBinaryFromReader(message: ProxySpec, reader: jspb.BinaryReader): ProxySpec;
}

export namespace ProxySpec {
  export type AsObject = {
    compressedspec: string,
    listenersList: Array<Listener.AsObject>,
  }
}

export class Listener extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getBindAddress(): string;
  setBindAddress(value: string): void;

  getBindPort(): number;
  setBindPort(value: number): void;

  hasHttpListener(): boolean;
  clearHttpListener(): void;
  getHttpListener(): HttpListener | undefined;
  setHttpListener(value?: HttpListener): void;

  hasTcpListener(): boolean;
  clearTcpListener(): void;
  getTcpListener(): TcpListener | undefined;
  setTcpListener(value?: TcpListener): void;

  hasHybridListener(): boolean;
  clearHybridListener(): void;
  getHybridListener(): HybridListener | undefined;
  setHybridListener(value?: HybridListener): void;

  hasAggregateListener(): boolean;
  clearAggregateListener(): void;
  getAggregateListener(): AggregateListener | undefined;
  setAggregateListener(value?: AggregateListener): void;

  clearSslConfigurationsList(): void;
  getSslConfigurationsList(): Array<github_com_solo_io_solo_apis_api_gloo_gloo_v1_ssl_pb.SslConfig>;
  setSslConfigurationsList(value: Array<github_com_solo_io_solo_apis_api_gloo_gloo_v1_ssl_pb.SslConfig>): void;
  addSslConfigurations(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_ssl_pb.SslConfig, index?: number): github_com_solo_io_solo_apis_api_gloo_gloo_v1_ssl_pb.SslConfig;

  hasUseProxyProto(): boolean;
  clearUseProxyProto(): void;
  getUseProxyProto(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setUseProxyProto(value?: google_protobuf_wrappers_pb.BoolValue): void;

  hasOptions(): boolean;
  clearOptions(): void;
  getOptions(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.ListenerOptions | undefined;
  setOptions(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.ListenerOptions): void;

  hasMetadata(): boolean;
  clearMetadata(): void;
  getMetadata(): google_protobuf_struct_pb.Struct | undefined;
  setMetadata(value?: google_protobuf_struct_pb.Struct): void;

  hasMetadataStatic(): boolean;
  clearMetadataStatic(): void;
  getMetadataStatic(): SourceMetadata | undefined;
  setMetadataStatic(value?: SourceMetadata): void;

  hasRouteOptions(): boolean;
  clearRouteOptions(): void;
  getRouteOptions(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.RouteConfigurationOptions | undefined;
  setRouteOptions(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.RouteConfigurationOptions): void;

  getListenertypeCase(): Listener.ListenertypeCase;
  getOpaqueMetadataCase(): Listener.OpaqueMetadataCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Listener.AsObject;
  static toObject(includeInstance: boolean, msg: Listener): Listener.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Listener, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Listener;
  static deserializeBinaryFromReader(message: Listener, reader: jspb.BinaryReader): Listener;
}

export namespace Listener {
  export type AsObject = {
    name: string,
    bindAddress: string,
    bindPort: number,
    httpListener?: HttpListener.AsObject,
    tcpListener?: TcpListener.AsObject,
    hybridListener?: HybridListener.AsObject,
    aggregateListener?: AggregateListener.AsObject,
    sslConfigurationsList: Array<github_com_solo_io_solo_apis_api_gloo_gloo_v1_ssl_pb.SslConfig.AsObject>,
    useProxyProto?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    options?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.ListenerOptions.AsObject,
    metadata?: google_protobuf_struct_pb.Struct.AsObject,
    metadataStatic?: SourceMetadata.AsObject,
    routeOptions?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.RouteConfigurationOptions.AsObject,
  }

  export enum ListenertypeCase {
    LISTENERTYPE_NOT_SET = 0,
    HTTP_LISTENER = 4,
    TCP_LISTENER = 5,
    HYBRID_LISTENER = 11,
    AGGREGATE_LISTENER = 13,
  }

  export enum OpaqueMetadataCase {
    OPAQUE_METADATA_NOT_SET = 0,
    METADATA = 9,
    METADATA_STATIC = 12,
  }
}

export class TcpListener extends jspb.Message {
  clearTcpHostsList(): void;
  getTcpHostsList(): Array<TcpHost>;
  setTcpHostsList(value: Array<TcpHost>): void;
  addTcpHosts(value?: TcpHost, index?: number): TcpHost;

  hasOptions(): boolean;
  clearOptions(): void;
  getOptions(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.TcpListenerOptions | undefined;
  setOptions(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.TcpListenerOptions): void;

  getStatPrefix(): string;
  setStatPrefix(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TcpListener.AsObject;
  static toObject(includeInstance: boolean, msg: TcpListener): TcpListener.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: TcpListener, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TcpListener;
  static deserializeBinaryFromReader(message: TcpListener, reader: jspb.BinaryReader): TcpListener;
}

export namespace TcpListener {
  export type AsObject = {
    tcpHostsList: Array<TcpHost.AsObject>,
    options?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.TcpListenerOptions.AsObject,
    statPrefix: string,
  }
}

export class TcpHost extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  hasSslConfig(): boolean;
  clearSslConfig(): void;
  getSslConfig(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_ssl_pb.SslConfig | undefined;
  setSslConfig(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_ssl_pb.SslConfig): void;

  hasDestination(): boolean;
  clearDestination(): void;
  getDestination(): TcpHost.TcpAction | undefined;
  setDestination(value?: TcpHost.TcpAction): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TcpHost.AsObject;
  static toObject(includeInstance: boolean, msg: TcpHost): TcpHost.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: TcpHost, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TcpHost;
  static deserializeBinaryFromReader(message: TcpHost, reader: jspb.BinaryReader): TcpHost;
}

export namespace TcpHost {
  export type AsObject = {
    name: string,
    sslConfig?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_ssl_pb.SslConfig.AsObject,
    destination?: TcpHost.TcpAction.AsObject,
  }

  export class TcpAction extends jspb.Message {
    hasSingle(): boolean;
    clearSingle(): void;
    getSingle(): Destination | undefined;
    setSingle(value?: Destination): void;

    hasMulti(): boolean;
    clearMulti(): void;
    getMulti(): MultiDestination | undefined;
    setMulti(value?: MultiDestination): void;

    hasUpstreamGroup(): boolean;
    clearUpstreamGroup(): void;
    getUpstreamGroup(): github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef | undefined;
    setUpstreamGroup(value?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef): void;

    hasForwardSniClusterName(): boolean;
    clearForwardSniClusterName(): void;
    getForwardSniClusterName(): google_protobuf_empty_pb.Empty | undefined;
    setForwardSniClusterName(value?: google_protobuf_empty_pb.Empty): void;

    getDestinationCase(): TcpAction.DestinationCase;
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): TcpAction.AsObject;
    static toObject(includeInstance: boolean, msg: TcpAction): TcpAction.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: TcpAction, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): TcpAction;
    static deserializeBinaryFromReader(message: TcpAction, reader: jspb.BinaryReader): TcpAction;
  }

  export namespace TcpAction {
    export type AsObject = {
      single?: Destination.AsObject,
      multi?: MultiDestination.AsObject,
      upstreamGroup?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject,
      forwardSniClusterName?: google_protobuf_empty_pb.Empty.AsObject,
    }

    export enum DestinationCase {
      DESTINATION_NOT_SET = 0,
      SINGLE = 1,
      MULTI = 2,
      UPSTREAM_GROUP = 3,
      FORWARD_SNI_CLUSTER_NAME = 4,
    }
  }
}

export class HttpListener extends jspb.Message {
  clearVirtualHostsList(): void;
  getVirtualHostsList(): Array<VirtualHost>;
  setVirtualHostsList(value: Array<VirtualHost>): void;
  addVirtualHosts(value?: VirtualHost, index?: number): VirtualHost;

  hasOptions(): boolean;
  clearOptions(): void;
  getOptions(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.HttpListenerOptions | undefined;
  setOptions(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.HttpListenerOptions): void;

  getStatPrefix(): string;
  setStatPrefix(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HttpListener.AsObject;
  static toObject(includeInstance: boolean, msg: HttpListener): HttpListener.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: HttpListener, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HttpListener;
  static deserializeBinaryFromReader(message: HttpListener, reader: jspb.BinaryReader): HttpListener;
}

export namespace HttpListener {
  export type AsObject = {
    virtualHostsList: Array<VirtualHost.AsObject>,
    options?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.HttpListenerOptions.AsObject,
    statPrefix: string,
  }
}

export class HybridListener extends jspb.Message {
  clearMatchedListenersList(): void;
  getMatchedListenersList(): Array<MatchedListener>;
  setMatchedListenersList(value: Array<MatchedListener>): void;
  addMatchedListeners(value?: MatchedListener, index?: number): MatchedListener;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HybridListener.AsObject;
  static toObject(includeInstance: boolean, msg: HybridListener): HybridListener.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: HybridListener, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HybridListener;
  static deserializeBinaryFromReader(message: HybridListener, reader: jspb.BinaryReader): HybridListener;
}

export namespace HybridListener {
  export type AsObject = {
    matchedListenersList: Array<MatchedListener.AsObject>,
  }
}

export class MatchedListener extends jspb.Message {
  hasMatcher(): boolean;
  clearMatcher(): void;
  getMatcher(): Matcher | undefined;
  setMatcher(value?: Matcher): void;

  hasHttpListener(): boolean;
  clearHttpListener(): void;
  getHttpListener(): HttpListener | undefined;
  setHttpListener(value?: HttpListener): void;

  hasTcpListener(): boolean;
  clearTcpListener(): void;
  getTcpListener(): TcpListener | undefined;
  setTcpListener(value?: TcpListener): void;

  clearSslConfigurationsList(): void;
  getSslConfigurationsList(): Array<github_com_solo_io_solo_apis_api_gloo_gloo_v1_ssl_pb.SslConfig>;
  setSslConfigurationsList(value: Array<github_com_solo_io_solo_apis_api_gloo_gloo_v1_ssl_pb.SslConfig>): void;
  addSslConfigurations(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_ssl_pb.SslConfig, index?: number): github_com_solo_io_solo_apis_api_gloo_gloo_v1_ssl_pb.SslConfig;

  getListenertypeCase(): MatchedListener.ListenertypeCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MatchedListener.AsObject;
  static toObject(includeInstance: boolean, msg: MatchedListener): MatchedListener.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: MatchedListener, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MatchedListener;
  static deserializeBinaryFromReader(message: MatchedListener, reader: jspb.BinaryReader): MatchedListener;
}

export namespace MatchedListener {
  export type AsObject = {
    matcher?: Matcher.AsObject,
    httpListener?: HttpListener.AsObject,
    tcpListener?: TcpListener.AsObject,
    sslConfigurationsList: Array<github_com_solo_io_solo_apis_api_gloo_gloo_v1_ssl_pb.SslConfig.AsObject>,
  }

  export enum ListenertypeCase {
    LISTENERTYPE_NOT_SET = 0,
    HTTP_LISTENER = 2,
    TCP_LISTENER = 3,
  }
}

export class Matcher extends jspb.Message {
  hasSslConfig(): boolean;
  clearSslConfig(): void;
  getSslConfig(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_ssl_pb.SslConfig | undefined;
  setSslConfig(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_ssl_pb.SslConfig): void;

  clearSourcePrefixRangesList(): void;
  getSourcePrefixRangesList(): Array<github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_address_pb.CidrRange>;
  setSourcePrefixRangesList(value: Array<github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_address_pb.CidrRange>): void;
  addSourcePrefixRanges(value?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_address_pb.CidrRange, index?: number): github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_address_pb.CidrRange;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Matcher.AsObject;
  static toObject(includeInstance: boolean, msg: Matcher): Matcher.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Matcher, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Matcher;
  static deserializeBinaryFromReader(message: Matcher, reader: jspb.BinaryReader): Matcher;
}

export namespace Matcher {
  export type AsObject = {
    sslConfig?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_ssl_pb.SslConfig.AsObject,
    sourcePrefixRangesList: Array<github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_address_pb.CidrRange.AsObject>,
  }
}

export class AggregateListener extends jspb.Message {
  hasHttpResources(): boolean;
  clearHttpResources(): void;
  getHttpResources(): AggregateListener.HttpResources | undefined;
  setHttpResources(value?: AggregateListener.HttpResources): void;

  clearHttpFilterChainsList(): void;
  getHttpFilterChainsList(): Array<AggregateListener.HttpFilterChain>;
  setHttpFilterChainsList(value: Array<AggregateListener.HttpFilterChain>): void;
  addHttpFilterChains(value?: AggregateListener.HttpFilterChain, index?: number): AggregateListener.HttpFilterChain;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AggregateListener.AsObject;
  static toObject(includeInstance: boolean, msg: AggregateListener): AggregateListener.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: AggregateListener, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AggregateListener;
  static deserializeBinaryFromReader(message: AggregateListener, reader: jspb.BinaryReader): AggregateListener;
}

export namespace AggregateListener {
  export type AsObject = {
    httpResources?: AggregateListener.HttpResources.AsObject,
    httpFilterChainsList: Array<AggregateListener.HttpFilterChain.AsObject>,
  }

  export class HttpResources extends jspb.Message {
    getVirtualHostsMap(): jspb.Map<string, VirtualHost>;
    clearVirtualHostsMap(): void;
    getHttpOptionsMap(): jspb.Map<string, github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.HttpListenerOptions>;
    clearHttpOptionsMap(): void;
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): HttpResources.AsObject;
    static toObject(includeInstance: boolean, msg: HttpResources): HttpResources.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: HttpResources, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): HttpResources;
    static deserializeBinaryFromReader(message: HttpResources, reader: jspb.BinaryReader): HttpResources;
  }

  export namespace HttpResources {
    export type AsObject = {
      virtualHostsMap: Array<[string, VirtualHost.AsObject]>,
      httpOptionsMap: Array<[string, github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.HttpListenerOptions.AsObject]>,
    }
  }

  export class HttpFilterChain extends jspb.Message {
    hasMatcher(): boolean;
    clearMatcher(): void;
    getMatcher(): Matcher | undefined;
    setMatcher(value?: Matcher): void;

    getHttpOptionsRef(): string;
    setHttpOptionsRef(value: string): void;

    clearVirtualHostRefsList(): void;
    getVirtualHostRefsList(): Array<string>;
    setVirtualHostRefsList(value: Array<string>): void;
    addVirtualHostRefs(value: string, index?: number): string;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): HttpFilterChain.AsObject;
    static toObject(includeInstance: boolean, msg: HttpFilterChain): HttpFilterChain.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: HttpFilterChain, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): HttpFilterChain;
    static deserializeBinaryFromReader(message: HttpFilterChain, reader: jspb.BinaryReader): HttpFilterChain;
  }

  export namespace HttpFilterChain {
    export type AsObject = {
      matcher?: Matcher.AsObject,
      httpOptionsRef: string,
      virtualHostRefsList: Array<string>,
    }
  }
}

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

  hasOptions(): boolean;
  clearOptions(): void;
  getOptions(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.VirtualHostOptions | undefined;
  setOptions(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.VirtualHostOptions): void;

  hasMetadata(): boolean;
  clearMetadata(): void;
  getMetadata(): google_protobuf_struct_pb.Struct | undefined;
  setMetadata(value?: google_protobuf_struct_pb.Struct): void;

  hasMetadataStatic(): boolean;
  clearMetadataStatic(): void;
  getMetadataStatic(): SourceMetadata | undefined;
  setMetadataStatic(value?: SourceMetadata): void;

  getOpaqueMetadataCase(): VirtualHost.OpaqueMetadataCase;
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
    options?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.VirtualHostOptions.AsObject,
    metadata?: google_protobuf_struct_pb.Struct.AsObject,
    metadataStatic?: SourceMetadata.AsObject,
  }

  export enum OpaqueMetadataCase {
    OPAQUE_METADATA_NOT_SET = 0,
    METADATA = 6,
    METADATA_STATIC = 7,
  }
}

export class Route extends jspb.Message {
  clearMatchersList(): void;
  getMatchersList(): Array<github_com_solo_io_solo_apis_api_gloo_gloo_v1_core_matchers_matchers_pb.Matcher>;
  setMatchersList(value: Array<github_com_solo_io_solo_apis_api_gloo_gloo_v1_core_matchers_matchers_pb.Matcher>): void;
  addMatchers(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_core_matchers_matchers_pb.Matcher, index?: number): github_com_solo_io_solo_apis_api_gloo_gloo_v1_core_matchers_matchers_pb.Matcher;

  hasRouteAction(): boolean;
  clearRouteAction(): void;
  getRouteAction(): RouteAction | undefined;
  setRouteAction(value?: RouteAction): void;

  hasRedirectAction(): boolean;
  clearRedirectAction(): void;
  getRedirectAction(): RedirectAction | undefined;
  setRedirectAction(value?: RedirectAction): void;

  hasDirectResponseAction(): boolean;
  clearDirectResponseAction(): void;
  getDirectResponseAction(): DirectResponseAction | undefined;
  setDirectResponseAction(value?: DirectResponseAction): void;

  hasGraphqlApiRef(): boolean;
  clearGraphqlApiRef(): void;
  getGraphqlApiRef(): github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef | undefined;
  setGraphqlApiRef(value?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef): void;

  hasOptions(): boolean;
  clearOptions(): void;
  getOptions(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.RouteOptions | undefined;
  setOptions(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.RouteOptions): void;

  hasMetadata(): boolean;
  clearMetadata(): void;
  getMetadata(): google_protobuf_struct_pb.Struct | undefined;
  setMetadata(value?: google_protobuf_struct_pb.Struct): void;

  hasMetadataStatic(): boolean;
  clearMetadataStatic(): void;
  getMetadataStatic(): SourceMetadata | undefined;
  setMetadataStatic(value?: SourceMetadata): void;

  getName(): string;
  setName(value: string): void;

  getActionCase(): Route.ActionCase;
  getOpaqueMetadataCase(): Route.OpaqueMetadataCase;
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
    matchersList: Array<github_com_solo_io_solo_apis_api_gloo_gloo_v1_core_matchers_matchers_pb.Matcher.AsObject>,
    routeAction?: RouteAction.AsObject,
    redirectAction?: RedirectAction.AsObject,
    directResponseAction?: DirectResponseAction.AsObject,
    graphqlApiRef?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject,
    options?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.RouteOptions.AsObject,
    metadata?: google_protobuf_struct_pb.Struct.AsObject,
    metadataStatic?: SourceMetadata.AsObject,
    name: string,
  }

  export enum ActionCase {
    ACTION_NOT_SET = 0,
    ROUTE_ACTION = 2,
    REDIRECT_ACTION = 3,
    DIRECT_RESPONSE_ACTION = 4,
    GRAPHQL_API_REF = 8,
  }

  export enum OpaqueMetadataCase {
    OPAQUE_METADATA_NOT_SET = 0,
    METADATA = 6,
    METADATA_STATIC = 9,
  }
}

export class RouteAction extends jspb.Message {
  hasSingle(): boolean;
  clearSingle(): void;
  getSingle(): Destination | undefined;
  setSingle(value?: Destination): void;

  hasMulti(): boolean;
  clearMulti(): void;
  getMulti(): MultiDestination | undefined;
  setMulti(value?: MultiDestination): void;

  hasUpstreamGroup(): boolean;
  clearUpstreamGroup(): void;
  getUpstreamGroup(): github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef | undefined;
  setUpstreamGroup(value?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef): void;

  hasClusterHeader(): boolean;
  clearClusterHeader(): void;
  getClusterHeader(): string;
  setClusterHeader(value: string): void;

  hasDynamicForwardProxy(): boolean;
  clearDynamicForwardProxy(): void;
  getDynamicForwardProxy(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_dynamic_forward_proxy_dynamic_forward_proxy_pb.PerRouteConfig | undefined;
  setDynamicForwardProxy(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_dynamic_forward_proxy_dynamic_forward_proxy_pb.PerRouteConfig): void;

  getDestinationCase(): RouteAction.DestinationCase;
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
    single?: Destination.AsObject,
    multi?: MultiDestination.AsObject,
    upstreamGroup?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject,
    clusterHeader: string,
    dynamicForwardProxy?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_dynamic_forward_proxy_dynamic_forward_proxy_pb.PerRouteConfig.AsObject,
  }

  export enum DestinationCase {
    DESTINATION_NOT_SET = 0,
    SINGLE = 1,
    MULTI = 2,
    UPSTREAM_GROUP = 3,
    CLUSTER_HEADER = 4,
    DYNAMIC_FORWARD_PROXY = 5,
  }
}

export class Destination extends jspb.Message {
  hasUpstream(): boolean;
  clearUpstream(): void;
  getUpstream(): github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef | undefined;
  setUpstream(value?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef): void;

  hasKube(): boolean;
  clearKube(): void;
  getKube(): KubernetesServiceDestination | undefined;
  setKube(value?: KubernetesServiceDestination): void;

  hasConsul(): boolean;
  clearConsul(): void;
  getConsul(): ConsulServiceDestination | undefined;
  setConsul(value?: ConsulServiceDestination): void;

  hasDestinationSpec(): boolean;
  clearDestinationSpec(): void;
  getDestinationSpec(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.DestinationSpec | undefined;
  setDestinationSpec(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.DestinationSpec): void;

  hasSubset(): boolean;
  clearSubset(): void;
  getSubset(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_subset_pb.Subset | undefined;
  setSubset(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_subset_pb.Subset): void;

  getDestinationTypeCase(): Destination.DestinationTypeCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Destination.AsObject;
  static toObject(includeInstance: boolean, msg: Destination): Destination.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Destination, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Destination;
  static deserializeBinaryFromReader(message: Destination, reader: jspb.BinaryReader): Destination;
}

export namespace Destination {
  export type AsObject = {
    upstream?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject,
    kube?: KubernetesServiceDestination.AsObject,
    consul?: ConsulServiceDestination.AsObject,
    destinationSpec?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.DestinationSpec.AsObject,
    subset?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_subset_pb.Subset.AsObject,
  }

  export enum DestinationTypeCase {
    DESTINATION_TYPE_NOT_SET = 0,
    UPSTREAM = 10,
    KUBE = 11,
    CONSUL = 12,
  }
}

export class KubernetesServiceDestination extends jspb.Message {
  hasRef(): boolean;
  clearRef(): void;
  getRef(): github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef | undefined;
  setRef(value?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef): void;

  getPort(): number;
  setPort(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): KubernetesServiceDestination.AsObject;
  static toObject(includeInstance: boolean, msg: KubernetesServiceDestination): KubernetesServiceDestination.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: KubernetesServiceDestination, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): KubernetesServiceDestination;
  static deserializeBinaryFromReader(message: KubernetesServiceDestination, reader: jspb.BinaryReader): KubernetesServiceDestination;
}

export namespace KubernetesServiceDestination {
  export type AsObject = {
    ref?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject,
    port: number,
  }
}

export class ConsulServiceDestination extends jspb.Message {
  getServiceName(): string;
  setServiceName(value: string): void;

  clearTagsList(): void;
  getTagsList(): Array<string>;
  setTagsList(value: Array<string>): void;
  addTags(value: string, index?: number): string;

  clearDataCentersList(): void;
  getDataCentersList(): Array<string>;
  setDataCentersList(value: Array<string>): void;
  addDataCenters(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ConsulServiceDestination.AsObject;
  static toObject(includeInstance: boolean, msg: ConsulServiceDestination): ConsulServiceDestination.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ConsulServiceDestination, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ConsulServiceDestination;
  static deserializeBinaryFromReader(message: ConsulServiceDestination, reader: jspb.BinaryReader): ConsulServiceDestination;
}

export namespace ConsulServiceDestination {
  export type AsObject = {
    serviceName: string,
    tagsList: Array<string>,
    dataCentersList: Array<string>,
  }
}

export class UpstreamGroupSpec extends jspb.Message {
  clearDestinationsList(): void;
  getDestinationsList(): Array<WeightedDestination>;
  setDestinationsList(value: Array<WeightedDestination>): void;
  addDestinations(value?: WeightedDestination, index?: number): WeightedDestination;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpstreamGroupSpec.AsObject;
  static toObject(includeInstance: boolean, msg: UpstreamGroupSpec): UpstreamGroupSpec.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpstreamGroupSpec, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpstreamGroupSpec;
  static deserializeBinaryFromReader(message: UpstreamGroupSpec, reader: jspb.BinaryReader): UpstreamGroupSpec;
}

export namespace UpstreamGroupSpec {
  export type AsObject = {
    destinationsList: Array<WeightedDestination.AsObject>,
  }
}

export class MultiDestination extends jspb.Message {
  clearDestinationsList(): void;
  getDestinationsList(): Array<WeightedDestination>;
  setDestinationsList(value: Array<WeightedDestination>): void;
  addDestinations(value?: WeightedDestination, index?: number): WeightedDestination;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MultiDestination.AsObject;
  static toObject(includeInstance: boolean, msg: MultiDestination): MultiDestination.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: MultiDestination, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MultiDestination;
  static deserializeBinaryFromReader(message: MultiDestination, reader: jspb.BinaryReader): MultiDestination;
}

export namespace MultiDestination {
  export type AsObject = {
    destinationsList: Array<WeightedDestination.AsObject>,
  }
}

export class WeightedDestination extends jspb.Message {
  hasDestination(): boolean;
  clearDestination(): void;
  getDestination(): Destination | undefined;
  setDestination(value?: Destination): void;

  getWeight(): number;
  setWeight(value: number): void;

  hasOptions(): boolean;
  clearOptions(): void;
  getOptions(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.WeightedDestinationOptions | undefined;
  setOptions(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.WeightedDestinationOptions): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): WeightedDestination.AsObject;
  static toObject(includeInstance: boolean, msg: WeightedDestination): WeightedDestination.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: WeightedDestination, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): WeightedDestination;
  static deserializeBinaryFromReader(message: WeightedDestination, reader: jspb.BinaryReader): WeightedDestination;
}

export namespace WeightedDestination {
  export type AsObject = {
    destination?: Destination.AsObject,
    weight: number,
    options?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.WeightedDestinationOptions.AsObject,
  }
}

export class RedirectAction extends jspb.Message {
  getHostRedirect(): string;
  setHostRedirect(value: string): void;

  hasPathRedirect(): boolean;
  clearPathRedirect(): void;
  getPathRedirect(): string;
  setPathRedirect(value: string): void;

  hasPrefixRewrite(): boolean;
  clearPrefixRewrite(): void;
  getPrefixRewrite(): string;
  setPrefixRewrite(value: string): void;

  hasRegexRewrite(): boolean;
  clearRegexRewrite(): void;
  getRegexRewrite(): github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_type_matcher_v3_regex_pb.RegexMatchAndSubstitute | undefined;
  setRegexRewrite(value?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_type_matcher_v3_regex_pb.RegexMatchAndSubstitute): void;

  getResponseCode(): RedirectAction.RedirectResponseCodeMap[keyof RedirectAction.RedirectResponseCodeMap];
  setResponseCode(value: RedirectAction.RedirectResponseCodeMap[keyof RedirectAction.RedirectResponseCodeMap]): void;

  getHttpsRedirect(): boolean;
  setHttpsRedirect(value: boolean): void;

  getStripQuery(): boolean;
  setStripQuery(value: boolean): void;

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
    hostRedirect: string,
    pathRedirect: string,
    prefixRewrite: string,
    regexRewrite?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_type_matcher_v3_regex_pb.RegexMatchAndSubstitute.AsObject,
    responseCode: RedirectAction.RedirectResponseCodeMap[keyof RedirectAction.RedirectResponseCodeMap],
    httpsRedirect: boolean,
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

  export enum PathRewriteSpecifierCase {
    PATH_REWRITE_SPECIFIER_NOT_SET = 0,
    PATH_REDIRECT = 2,
    PREFIX_REWRITE = 5,
    REGEX_REWRITE = 32,
  }
}

export class DirectResponseAction extends jspb.Message {
  getStatus(): number;
  setStatus(value: number): void;

  getBody(): string;
  setBody(value: string): void;

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
    body: string,
  }
}

export class SourceMetadata extends jspb.Message {
  clearSourcesList(): void;
  getSourcesList(): Array<SourceMetadata.SourceRef>;
  setSourcesList(value: Array<SourceMetadata.SourceRef>): void;
  addSources(value?: SourceMetadata.SourceRef, index?: number): SourceMetadata.SourceRef;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SourceMetadata.AsObject;
  static toObject(includeInstance: boolean, msg: SourceMetadata): SourceMetadata.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: SourceMetadata, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SourceMetadata;
  static deserializeBinaryFromReader(message: SourceMetadata, reader: jspb.BinaryReader): SourceMetadata;
}

export namespace SourceMetadata {
  export type AsObject = {
    sourcesList: Array<SourceMetadata.SourceRef.AsObject>,
  }

  export class SourceRef extends jspb.Message {
    hasResourceRef(): boolean;
    clearResourceRef(): void;
    getResourceRef(): github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef | undefined;
    setResourceRef(value?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef): void;

    getResourceKind(): string;
    setResourceKind(value: string): void;

    getObservedGeneration(): number;
    setObservedGeneration(value: number): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): SourceRef.AsObject;
    static toObject(includeInstance: boolean, msg: SourceRef): SourceRef.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: SourceRef, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): SourceRef;
    static deserializeBinaryFromReader(message: SourceRef, reader: jspb.BinaryReader): SourceRef;
  }

  export namespace SourceRef {
    export type AsObject = {
      resourceRef?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject,
      resourceKind: string,
      observedGeneration: number,
    }
  }
}

export class UpstreamGroupStatus extends jspb.Message {
  getState(): UpstreamGroupStatus.StateMap[keyof UpstreamGroupStatus.StateMap];
  setState(value: UpstreamGroupStatus.StateMap[keyof UpstreamGroupStatus.StateMap]): void;

  getReason(): string;
  setReason(value: string): void;

  getReportedBy(): string;
  setReportedBy(value: string): void;

  getSubresourceStatusesMap(): jspb.Map<string, UpstreamGroupStatus>;
  clearSubresourceStatusesMap(): void;
  hasDetails(): boolean;
  clearDetails(): void;
  getDetails(): google_protobuf_struct_pb.Struct | undefined;
  setDetails(value?: google_protobuf_struct_pb.Struct): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpstreamGroupStatus.AsObject;
  static toObject(includeInstance: boolean, msg: UpstreamGroupStatus): UpstreamGroupStatus.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpstreamGroupStatus, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpstreamGroupStatus;
  static deserializeBinaryFromReader(message: UpstreamGroupStatus, reader: jspb.BinaryReader): UpstreamGroupStatus;
}

export namespace UpstreamGroupStatus {
  export type AsObject = {
    state: UpstreamGroupStatus.StateMap[keyof UpstreamGroupStatus.StateMap],
    reason: string,
    reportedBy: string,
    subresourceStatusesMap: Array<[string, UpstreamGroupStatus.AsObject]>,
    details?: google_protobuf_struct_pb.Struct.AsObject,
  }

  export interface StateMap {
    PENDING: 0;
    ACCEPTED: 1;
    REJECTED: 2;
    WARNING: 3;
  }

  export const State: StateMap;
}

export class UpstreamGroupNamespacedStatuses extends jspb.Message {
  getStatusesMap(): jspb.Map<string, UpstreamGroupStatus>;
  clearStatusesMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpstreamGroupNamespacedStatuses.AsObject;
  static toObject(includeInstance: boolean, msg: UpstreamGroupNamespacedStatuses): UpstreamGroupNamespacedStatuses.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpstreamGroupNamespacedStatuses, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpstreamGroupNamespacedStatuses;
  static deserializeBinaryFromReader(message: UpstreamGroupNamespacedStatuses, reader: jspb.BinaryReader): UpstreamGroupNamespacedStatuses;
}

export namespace UpstreamGroupNamespacedStatuses {
  export type AsObject = {
    statusesMap: Array<[string, UpstreamGroupStatus.AsObject]>,
  }
}

export class ProxyStatus extends jspb.Message {
  getState(): ProxyStatus.StateMap[keyof ProxyStatus.StateMap];
  setState(value: ProxyStatus.StateMap[keyof ProxyStatus.StateMap]): void;

  getReason(): string;
  setReason(value: string): void;

  getReportedBy(): string;
  setReportedBy(value: string): void;

  getSubresourceStatusesMap(): jspb.Map<string, ProxyStatus>;
  clearSubresourceStatusesMap(): void;
  hasDetails(): boolean;
  clearDetails(): void;
  getDetails(): google_protobuf_struct_pb.Struct | undefined;
  setDetails(value?: google_protobuf_struct_pb.Struct): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProxyStatus.AsObject;
  static toObject(includeInstance: boolean, msg: ProxyStatus): ProxyStatus.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ProxyStatus, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProxyStatus;
  static deserializeBinaryFromReader(message: ProxyStatus, reader: jspb.BinaryReader): ProxyStatus;
}

export namespace ProxyStatus {
  export type AsObject = {
    state: ProxyStatus.StateMap[keyof ProxyStatus.StateMap],
    reason: string,
    reportedBy: string,
    subresourceStatusesMap: Array<[string, ProxyStatus.AsObject]>,
    details?: google_protobuf_struct_pb.Struct.AsObject,
  }

  export interface StateMap {
    PENDING: 0;
    ACCEPTED: 1;
    REJECTED: 2;
    WARNING: 3;
  }

  export const State: StateMap;
}

export class ProxyNamespacedStatuses extends jspb.Message {
  getStatusesMap(): jspb.Map<string, ProxyStatus>;
  clearStatusesMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProxyNamespacedStatuses.AsObject;
  static toObject(includeInstance: boolean, msg: ProxyNamespacedStatuses): ProxyNamespacedStatuses.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ProxyNamespacedStatuses, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProxyNamespacedStatuses;
  static deserializeBinaryFromReader(message: ProxyNamespacedStatuses, reader: jspb.BinaryReader): ProxyNamespacedStatuses;
}

export namespace ProxyNamespacedStatuses {
  export type AsObject = {
    statusesMap: Array<[string, ProxyStatus.AsObject]>,
  }
}
