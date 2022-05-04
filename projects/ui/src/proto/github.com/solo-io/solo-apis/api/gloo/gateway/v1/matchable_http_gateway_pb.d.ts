/* eslint-disable */
// package: gateway.solo.io
// file: github.com/solo-io/solo-apis/api/gloo/gateway/v1/matchable_http_gateway.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_struct_pb from "google-protobuf/google/protobuf/struct_pb";
import * as extproto_ext_pb from "../../../../../../../extproto/ext_pb";
import * as github_com_solo_io_solo_kit_api_v1_ref_pb from "../../../../../../../github.com/solo-io/solo-kit/api/v1/ref_pb";
import * as github_com_solo_io_solo_kit_api_v1_solo_kit_pb from "../../../../../../../github.com/solo-io/solo-kit/api/v1/solo-kit_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gateway_v1_http_gateway_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gateway/v1/http_gateway_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_ssl_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/ssl_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_address_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/config/core/v3/address_pb";

export class MatchableHttpGatewaySpec extends jspb.Message {
  hasMatcher(): boolean;
  clearMatcher(): void;
  getMatcher(): MatchableHttpGatewaySpec.Matcher | undefined;
  setMatcher(value?: MatchableHttpGatewaySpec.Matcher): void;

  hasHttpGateway(): boolean;
  clearHttpGateway(): void;
  getHttpGateway(): github_com_solo_io_solo_apis_api_gloo_gateway_v1_http_gateway_pb.HttpGateway | undefined;
  setHttpGateway(value?: github_com_solo_io_solo_apis_api_gloo_gateway_v1_http_gateway_pb.HttpGateway): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MatchableHttpGatewaySpec.AsObject;
  static toObject(includeInstance: boolean, msg: MatchableHttpGatewaySpec): MatchableHttpGatewaySpec.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: MatchableHttpGatewaySpec, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MatchableHttpGatewaySpec;
  static deserializeBinaryFromReader(message: MatchableHttpGatewaySpec, reader: jspb.BinaryReader): MatchableHttpGatewaySpec;
}

export namespace MatchableHttpGatewaySpec {
  export type AsObject = {
    matcher?: MatchableHttpGatewaySpec.Matcher.AsObject,
    httpGateway?: github_com_solo_io_solo_apis_api_gloo_gateway_v1_http_gateway_pb.HttpGateway.AsObject,
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

export class MatchableHttpGatewayStatus extends jspb.Message {
  getState(): MatchableHttpGatewayStatus.StateMap[keyof MatchableHttpGatewayStatus.StateMap];
  setState(value: MatchableHttpGatewayStatus.StateMap[keyof MatchableHttpGatewayStatus.StateMap]): void;

  getReason(): string;
  setReason(value: string): void;

  getReportedBy(): string;
  setReportedBy(value: string): void;

  getSubresourceStatusesMap(): jspb.Map<string, MatchableHttpGatewayStatus>;
  clearSubresourceStatusesMap(): void;
  hasDetails(): boolean;
  clearDetails(): void;
  getDetails(): google_protobuf_struct_pb.Struct | undefined;
  setDetails(value?: google_protobuf_struct_pb.Struct): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MatchableHttpGatewayStatus.AsObject;
  static toObject(includeInstance: boolean, msg: MatchableHttpGatewayStatus): MatchableHttpGatewayStatus.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: MatchableHttpGatewayStatus, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MatchableHttpGatewayStatus;
  static deserializeBinaryFromReader(message: MatchableHttpGatewayStatus, reader: jspb.BinaryReader): MatchableHttpGatewayStatus;
}

export namespace MatchableHttpGatewayStatus {
  export type AsObject = {
    state: MatchableHttpGatewayStatus.StateMap[keyof MatchableHttpGatewayStatus.StateMap],
    reason: string,
    reportedBy: string,
    subresourceStatusesMap: Array<[string, MatchableHttpGatewayStatus.AsObject]>,
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

export class MatchableHttpGatewayNamespacedStatuses extends jspb.Message {
  getStatusesMap(): jspb.Map<string, MatchableHttpGatewayStatus>;
  clearStatusesMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MatchableHttpGatewayNamespacedStatuses.AsObject;
  static toObject(includeInstance: boolean, msg: MatchableHttpGatewayNamespacedStatuses): MatchableHttpGatewayNamespacedStatuses.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: MatchableHttpGatewayNamespacedStatuses, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MatchableHttpGatewayNamespacedStatuses;
  static deserializeBinaryFromReader(message: MatchableHttpGatewayNamespacedStatuses, reader: jspb.BinaryReader): MatchableHttpGatewayNamespacedStatuses;
}

export namespace MatchableHttpGatewayNamespacedStatuses {
  export type AsObject = {
    statusesMap: Array<[string, MatchableHttpGatewayStatus.AsObject]>,
  }
}
