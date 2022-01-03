/* eslint-disable */
// package: gateway.solo.io
// file: github.com/solo-io/solo-apis/api/gloo/gateway/v1/gateway.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_struct_pb from "google-protobuf/google/protobuf/struct_pb";
import * as extproto_ext_pb from "../../../../../../../extproto/ext_pb";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";
import * as github_com_solo_io_solo_kit_api_v1_ref_pb from "../../../../../../../github.com/solo-io/solo-kit/api/v1/ref_pb";
import * as github_com_solo_io_solo_kit_api_v1_solo_kit_pb from "../../../../../../../github.com/solo-io/solo-kit/api/v1/solo-kit_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_proxy_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/proxy_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/options_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_ssl_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/ssl_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_address_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/config/core/v3/address_pb";

export class GatewaySpec extends jspb.Message {
  getSsl(): boolean;
  setSsl(value: boolean): void;

  getBindAddress(): string;
  setBindAddress(value: string): void;

  getBindPort(): number;
  setBindPort(value: number): void;

  hasOptions(): boolean;
  clearOptions(): void;
  getOptions(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.ListenerOptions | undefined;
  setOptions(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.ListenerOptions): void;

  hasUseProxyProto(): boolean;
  clearUseProxyProto(): void;
  getUseProxyProto(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setUseProxyProto(value?: google_protobuf_wrappers_pb.BoolValue): void;

  hasHttpGateway(): boolean;
  clearHttpGateway(): void;
  getHttpGateway(): HttpGateway | undefined;
  setHttpGateway(value?: HttpGateway): void;

  hasTcpGateway(): boolean;
  clearTcpGateway(): void;
  getTcpGateway(): TcpGateway | undefined;
  setTcpGateway(value?: TcpGateway): void;

  hasHybridGateway(): boolean;
  clearHybridGateway(): void;
  getHybridGateway(): HybridGateway | undefined;
  setHybridGateway(value?: HybridGateway): void;

  clearProxyNamesList(): void;
  getProxyNamesList(): Array<string>;
  setProxyNamesList(value: Array<string>): void;
  addProxyNames(value: string, index?: number): string;

  hasRouteOptions(): boolean;
  clearRouteOptions(): void;
  getRouteOptions(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.RouteConfigurationOptions | undefined;
  setRouteOptions(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.RouteConfigurationOptions): void;

  getGatewaytypeCase(): GatewaySpec.GatewaytypeCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GatewaySpec.AsObject;
  static toObject(includeInstance: boolean, msg: GatewaySpec): GatewaySpec.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GatewaySpec, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GatewaySpec;
  static deserializeBinaryFromReader(message: GatewaySpec, reader: jspb.BinaryReader): GatewaySpec;
}

export namespace GatewaySpec {
  export type AsObject = {
    ssl: boolean,
    bindAddress: string,
    bindPort: number,
    options?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.ListenerOptions.AsObject,
    useProxyProto?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    httpGateway?: HttpGateway.AsObject,
    tcpGateway?: TcpGateway.AsObject,
    hybridGateway?: HybridGateway.AsObject,
    proxyNamesList: Array<string>,
    routeOptions?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.RouteConfigurationOptions.AsObject,
  }

  export enum GatewaytypeCase {
    GATEWAYTYPE_NOT_SET = 0,
    HTTP_GATEWAY = 9,
    TCP_GATEWAY = 10,
    HYBRID_GATEWAY = 11,
  }
}

export class HttpGateway extends jspb.Message {
  clearVirtualServicesList(): void;
  getVirtualServicesList(): Array<github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef>;
  setVirtualServicesList(value: Array<github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef>): void;
  addVirtualServices(value?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef, index?: number): github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef;

  getVirtualServiceSelectorMap(): jspb.Map<string, string>;
  clearVirtualServiceSelectorMap(): void;
  hasVirtualServiceExpressions(): boolean;
  clearVirtualServiceExpressions(): void;
  getVirtualServiceExpressions(): VirtualServiceSelectorExpressions | undefined;
  setVirtualServiceExpressions(value?: VirtualServiceSelectorExpressions): void;

  clearVirtualServiceNamespacesList(): void;
  getVirtualServiceNamespacesList(): Array<string>;
  setVirtualServiceNamespacesList(value: Array<string>): void;
  addVirtualServiceNamespaces(value: string, index?: number): string;

  hasOptions(): boolean;
  clearOptions(): void;
  getOptions(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.HttpListenerOptions | undefined;
  setOptions(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.HttpListenerOptions): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HttpGateway.AsObject;
  static toObject(includeInstance: boolean, msg: HttpGateway): HttpGateway.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: HttpGateway, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HttpGateway;
  static deserializeBinaryFromReader(message: HttpGateway, reader: jspb.BinaryReader): HttpGateway;
}

export namespace HttpGateway {
  export type AsObject = {
    virtualServicesList: Array<github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject>,
    virtualServiceSelectorMap: Array<[string, string]>,
    virtualServiceExpressions?: VirtualServiceSelectorExpressions.AsObject,
    virtualServiceNamespacesList: Array<string>,
    options?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.HttpListenerOptions.AsObject,
  }
}

export class TcpGateway extends jspb.Message {
  clearTcpHostsList(): void;
  getTcpHostsList(): Array<github_com_solo_io_solo_apis_api_gloo_gloo_v1_proxy_pb.TcpHost>;
  setTcpHostsList(value: Array<github_com_solo_io_solo_apis_api_gloo_gloo_v1_proxy_pb.TcpHost>): void;
  addTcpHosts(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_proxy_pb.TcpHost, index?: number): github_com_solo_io_solo_apis_api_gloo_gloo_v1_proxy_pb.TcpHost;

  hasOptions(): boolean;
  clearOptions(): void;
  getOptions(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.TcpListenerOptions | undefined;
  setOptions(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.TcpListenerOptions): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TcpGateway.AsObject;
  static toObject(includeInstance: boolean, msg: TcpGateway): TcpGateway.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: TcpGateway, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TcpGateway;
  static deserializeBinaryFromReader(message: TcpGateway, reader: jspb.BinaryReader): TcpGateway;
}

export namespace TcpGateway {
  export type AsObject = {
    tcpHostsList: Array<github_com_solo_io_solo_apis_api_gloo_gloo_v1_proxy_pb.TcpHost.AsObject>,
    options?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.TcpListenerOptions.AsObject,
  }
}

export class HybridGateway extends jspb.Message {
  clearMatchedGatewaysList(): void;
  getMatchedGatewaysList(): Array<MatchedGateway>;
  setMatchedGatewaysList(value: Array<MatchedGateway>): void;
  addMatchedGateways(value?: MatchedGateway, index?: number): MatchedGateway;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HybridGateway.AsObject;
  static toObject(includeInstance: boolean, msg: HybridGateway): HybridGateway.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: HybridGateway, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HybridGateway;
  static deserializeBinaryFromReader(message: HybridGateway, reader: jspb.BinaryReader): HybridGateway;
}

export namespace HybridGateway {
  export type AsObject = {
    matchedGatewaysList: Array<MatchedGateway.AsObject>,
  }
}

export class MatchedGateway extends jspb.Message {
  hasMatcher(): boolean;
  clearMatcher(): void;
  getMatcher(): Matcher | undefined;
  setMatcher(value?: Matcher): void;

  hasHttpGateway(): boolean;
  clearHttpGateway(): void;
  getHttpGateway(): HttpGateway | undefined;
  setHttpGateway(value?: HttpGateway): void;

  hasTcpGateway(): boolean;
  clearTcpGateway(): void;
  getTcpGateway(): TcpGateway | undefined;
  setTcpGateway(value?: TcpGateway): void;

  getGatewaytypeCase(): MatchedGateway.GatewaytypeCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MatchedGateway.AsObject;
  static toObject(includeInstance: boolean, msg: MatchedGateway): MatchedGateway.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: MatchedGateway, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MatchedGateway;
  static deserializeBinaryFromReader(message: MatchedGateway, reader: jspb.BinaryReader): MatchedGateway;
}

export namespace MatchedGateway {
  export type AsObject = {
    matcher?: Matcher.AsObject,
    httpGateway?: HttpGateway.AsObject,
    tcpGateway?: TcpGateway.AsObject,
  }

  export enum GatewaytypeCase {
    GATEWAYTYPE_NOT_SET = 0,
    HTTP_GATEWAY = 2,
    TCP_GATEWAY = 3,
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

export class VirtualServiceSelectorExpressions extends jspb.Message {
  clearExpressionsList(): void;
  getExpressionsList(): Array<VirtualServiceSelectorExpressions.Expression>;
  setExpressionsList(value: Array<VirtualServiceSelectorExpressions.Expression>): void;
  addExpressions(value?: VirtualServiceSelectorExpressions.Expression, index?: number): VirtualServiceSelectorExpressions.Expression;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VirtualServiceSelectorExpressions.AsObject;
  static toObject(includeInstance: boolean, msg: VirtualServiceSelectorExpressions): VirtualServiceSelectorExpressions.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: VirtualServiceSelectorExpressions, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VirtualServiceSelectorExpressions;
  static deserializeBinaryFromReader(message: VirtualServiceSelectorExpressions, reader: jspb.BinaryReader): VirtualServiceSelectorExpressions;
}

export namespace VirtualServiceSelectorExpressions {
  export type AsObject = {
    expressionsList: Array<VirtualServiceSelectorExpressions.Expression.AsObject>,
  }

  export class Expression extends jspb.Message {
    getKey(): string;
    setKey(value: string): void;

    getOperator(): VirtualServiceSelectorExpressions.Expression.OperatorMap[keyof VirtualServiceSelectorExpressions.Expression.OperatorMap];
    setOperator(value: VirtualServiceSelectorExpressions.Expression.OperatorMap[keyof VirtualServiceSelectorExpressions.Expression.OperatorMap]): void;

    clearValuesList(): void;
    getValuesList(): Array<string>;
    setValuesList(value: Array<string>): void;
    addValues(value: string, index?: number): string;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Expression.AsObject;
    static toObject(includeInstance: boolean, msg: Expression): Expression.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Expression, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Expression;
    static deserializeBinaryFromReader(message: Expression, reader: jspb.BinaryReader): Expression;
  }

  export namespace Expression {
    export type AsObject = {
      key: string,
      operator: VirtualServiceSelectorExpressions.Expression.OperatorMap[keyof VirtualServiceSelectorExpressions.Expression.OperatorMap],
      valuesList: Array<string>,
    }

    export interface OperatorMap {
      EQUALS: 0;
      DOUBLEEQUALS: 1;
      NOTEQUALS: 2;
      IN: 3;
      NOTIN: 4;
      EXISTS: 5;
      DOESNOTEXIST: 6;
      GREATERTHAN: 7;
      LESSTHAN: 8;
    }

    export const Operator: OperatorMap;
  }
}

export class GatewayStatus extends jspb.Message {
  getState(): GatewayStatus.StateMap[keyof GatewayStatus.StateMap];
  setState(value: GatewayStatus.StateMap[keyof GatewayStatus.StateMap]): void;

  getReason(): string;
  setReason(value: string): void;

  getReportedBy(): string;
  setReportedBy(value: string): void;

  getSubresourceStatusesMap(): jspb.Map<string, GatewayStatus>;
  clearSubresourceStatusesMap(): void;
  hasDetails(): boolean;
  clearDetails(): void;
  getDetails(): google_protobuf_struct_pb.Struct | undefined;
  setDetails(value?: google_protobuf_struct_pb.Struct): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GatewayStatus.AsObject;
  static toObject(includeInstance: boolean, msg: GatewayStatus): GatewayStatus.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GatewayStatus, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GatewayStatus;
  static deserializeBinaryFromReader(message: GatewayStatus, reader: jspb.BinaryReader): GatewayStatus;
}

export namespace GatewayStatus {
  export type AsObject = {
    state: GatewayStatus.StateMap[keyof GatewayStatus.StateMap],
    reason: string,
    reportedBy: string,
    subresourceStatusesMap: Array<[string, GatewayStatus.AsObject]>,
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

export class GatewayNamespacedStatuses extends jspb.Message {
  getStatusesMap(): jspb.Map<string, GatewayStatus>;
  clearStatusesMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GatewayNamespacedStatuses.AsObject;
  static toObject(includeInstance: boolean, msg: GatewayNamespacedStatuses): GatewayNamespacedStatuses.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GatewayNamespacedStatuses, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GatewayNamespacedStatuses;
  static deserializeBinaryFromReader(message: GatewayNamespacedStatuses, reader: jspb.BinaryReader): GatewayNamespacedStatuses;
}

export namespace GatewayNamespacedStatuses {
  export type AsObject = {
    statusesMap: Array<[string, GatewayStatus.AsObject]>,
  }
}
