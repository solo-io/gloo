/* eslint-disable */
// package: gateway.solo.io
// file: github.com/solo-io/solo-apis/api/gloo/gateway/v1/matchable_tcp_gateway.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_struct_pb from "google-protobuf/google/protobuf/struct_pb";
import * as extproto_ext_pb from "../../../../../../../extproto/ext_pb";
import * as github_com_solo_io_solo_kit_api_v1_ref_pb from "../../../../../../../github.com/solo-io/solo-kit/api/v1/ref_pb";
import * as github_com_solo_io_solo_kit_api_v1_solo_kit_pb from "../../../../../../../github.com/solo-io/solo-kit/api/v1/solo-kit_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gateway_v1_gateway_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gateway/v1/gateway_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_ssl_ssl_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/ssl/ssl_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_address_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/config/core/v3/address_pb";

export class MatchableTcpGatewaySpec extends jspb.Message {
  hasMatcher(): boolean;
  clearMatcher(): void;
  getMatcher(): MatchableTcpGatewaySpec.Matcher | undefined;
  setMatcher(value?: MatchableTcpGatewaySpec.Matcher): void;

  hasTcpGateway(): boolean;
  clearTcpGateway(): void;
  getTcpGateway(): github_com_solo_io_solo_apis_api_gloo_gateway_v1_gateway_pb.TcpGateway | undefined;
  setTcpGateway(value?: github_com_solo_io_solo_apis_api_gloo_gateway_v1_gateway_pb.TcpGateway): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MatchableTcpGatewaySpec.AsObject;
  static toObject(includeInstance: boolean, msg: MatchableTcpGatewaySpec): MatchableTcpGatewaySpec.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: MatchableTcpGatewaySpec, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MatchableTcpGatewaySpec;
  static deserializeBinaryFromReader(message: MatchableTcpGatewaySpec, reader: jspb.BinaryReader): MatchableTcpGatewaySpec;
}

export namespace MatchableTcpGatewaySpec {
  export type AsObject = {
    matcher?: MatchableTcpGatewaySpec.Matcher.AsObject,
    tcpGateway?: github_com_solo_io_solo_apis_api_gloo_gateway_v1_gateway_pb.TcpGateway.AsObject,
  }

  export class Matcher extends jspb.Message {
    clearSourcePrefixRangesList(): void;
    getSourcePrefixRangesList(): Array<github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_address_pb.CidrRange>;
    setSourcePrefixRangesList(value: Array<github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_address_pb.CidrRange>): void;
    addSourcePrefixRanges(value?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_address_pb.CidrRange, index?: number): github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_address_pb.CidrRange;

    hasSslConfig(): boolean;
    clearSslConfig(): void;
    getSslConfig(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_ssl_ssl_pb.SslConfig | undefined;
    setSslConfig(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_ssl_ssl_pb.SslConfig): void;

    clearPassthroughCipherSuitesList(): void;
    getPassthroughCipherSuitesList(): Array<string>;
    setPassthroughCipherSuitesList(value: Array<string>): void;
    addPassthroughCipherSuites(value: string, index?: number): string;

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
      sslConfig?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_ssl_ssl_pb.SslConfig.AsObject,
      passthroughCipherSuitesList: Array<string>,
    }
  }
}

export class MatchableTcpGatewayStatus extends jspb.Message {
  getState(): MatchableTcpGatewayStatus.StateMap[keyof MatchableTcpGatewayStatus.StateMap];
  setState(value: MatchableTcpGatewayStatus.StateMap[keyof MatchableTcpGatewayStatus.StateMap]): void;

  getReason(): string;
  setReason(value: string): void;

  getReportedBy(): string;
  setReportedBy(value: string): void;

  getSubresourceStatusesMap(): jspb.Map<string, MatchableTcpGatewayStatus>;
  clearSubresourceStatusesMap(): void;
  hasDetails(): boolean;
  clearDetails(): void;
  getDetails(): google_protobuf_struct_pb.Struct | undefined;
  setDetails(value?: google_protobuf_struct_pb.Struct): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MatchableTcpGatewayStatus.AsObject;
  static toObject(includeInstance: boolean, msg: MatchableTcpGatewayStatus): MatchableTcpGatewayStatus.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: MatchableTcpGatewayStatus, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MatchableTcpGatewayStatus;
  static deserializeBinaryFromReader(message: MatchableTcpGatewayStatus, reader: jspb.BinaryReader): MatchableTcpGatewayStatus;
}

export namespace MatchableTcpGatewayStatus {
  export type AsObject = {
    state: MatchableTcpGatewayStatus.StateMap[keyof MatchableTcpGatewayStatus.StateMap],
    reason: string,
    reportedBy: string,
    subresourceStatusesMap: Array<[string, MatchableTcpGatewayStatus.AsObject]>,
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

export class MatchableTcpGatewayNamespacedStatuses extends jspb.Message {
  getStatusesMap(): jspb.Map<string, MatchableTcpGatewayStatus>;
  clearStatusesMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MatchableTcpGatewayNamespacedStatuses.AsObject;
  static toObject(includeInstance: boolean, msg: MatchableTcpGatewayNamespacedStatuses): MatchableTcpGatewayNamespacedStatuses.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: MatchableTcpGatewayNamespacedStatuses, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MatchableTcpGatewayNamespacedStatuses;
  static deserializeBinaryFromReader(message: MatchableTcpGatewayNamespacedStatuses, reader: jspb.BinaryReader): MatchableTcpGatewayNamespacedStatuses;
}

export namespace MatchableTcpGatewayNamespacedStatuses {
  export type AsObject = {
    statusesMap: Array<[string, MatchableTcpGatewayStatus.AsObject]>,
  }
}
