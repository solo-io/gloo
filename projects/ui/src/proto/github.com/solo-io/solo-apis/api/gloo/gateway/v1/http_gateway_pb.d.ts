/* eslint-disable */
// package: gateway.solo.io
// file: github.com/solo-io/solo-apis/api/gloo/gateway/v1/http_gateway.proto

import * as jspb from "google-protobuf";
import * as extproto_ext_pb from "../../../../../../../extproto/ext_pb";
import * as github_com_solo_io_solo_kit_api_v1_metadata_pb from "../../../../../../../github.com/solo-io/solo-kit/api/v1/metadata_pb";
import * as github_com_solo_io_solo_kit_api_v1_status_pb from "../../../../../../../github.com/solo-io/solo-kit/api/v1/status_pb";
import * as github_com_solo_io_solo_kit_api_v1_ref_pb from "../../../../../../../github.com/solo-io/solo-kit/api/v1/ref_pb";
import * as github_com_solo_io_solo_kit_api_v1_solo_kit_pb from "../../../../../../../github.com/solo-io/solo-kit/api/v1/solo-kit_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_proxy_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/proxy_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/options_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_ssl_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/ssl_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_address_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/config/core/v3/address_pb";

export class MatchableHttpGateway extends jspb.Message {
  hasNamespacedStatuses(): boolean;
  clearNamespacedStatuses(): void;
  getNamespacedStatuses(): github_com_solo_io_solo_kit_api_v1_status_pb.NamespacedStatuses | undefined;
  setNamespacedStatuses(value?: github_com_solo_io_solo_kit_api_v1_status_pb.NamespacedStatuses): void;

  hasMetadata(): boolean;
  clearMetadata(): void;
  getMetadata(): github_com_solo_io_solo_kit_api_v1_metadata_pb.Metadata | undefined;
  setMetadata(value?: github_com_solo_io_solo_kit_api_v1_metadata_pb.Metadata): void;

  hasMatcher(): boolean;
  clearMatcher(): void;
  getMatcher(): MatchableHttpGateway.Matcher | undefined;
  setMatcher(value?: MatchableHttpGateway.Matcher): void;

  hasHttpGateway(): boolean;
  clearHttpGateway(): void;
  getHttpGateway(): HttpGateway | undefined;
  setHttpGateway(value?: HttpGateway): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MatchableHttpGateway.AsObject;
  static toObject(includeInstance: boolean, msg: MatchableHttpGateway): MatchableHttpGateway.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: MatchableHttpGateway, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MatchableHttpGateway;
  static deserializeBinaryFromReader(message: MatchableHttpGateway, reader: jspb.BinaryReader): MatchableHttpGateway;
}

export namespace MatchableHttpGateway {
  export type AsObject = {
    namespacedStatuses?: github_com_solo_io_solo_kit_api_v1_status_pb.NamespacedStatuses.AsObject,
    metadata?: github_com_solo_io_solo_kit_api_v1_metadata_pb.Metadata.AsObject,
    matcher?: MatchableHttpGateway.Matcher.AsObject,
    httpGateway?: HttpGateway.AsObject,
  }

  export class Matcher extends jspb.Message {
    clearSourcePrefixRangesList(): void;
    getSourcePrefixRangesList(): Array<github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_address_pb.CidrRange>;
    setSourcePrefixRangesList(value: Array<github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_address_pb.CidrRange>): void;
    addSourcePrefixRanges(value?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_address_pb.CidrRange, index?: number): github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_address_pb.CidrRange;

    hasSslConfig(): boolean;
    clearSslConfig(): void;
    getSslConfig(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_ssl_pb.SslConfig | undefined;
    setSslConfig(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_ssl_pb.SslConfig): void;

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
      sourcePrefixRangesList: Array<github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_address_pb.CidrRange.AsObject>,
      sslConfig?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_ssl_pb.SslConfig.AsObject,
    }
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
