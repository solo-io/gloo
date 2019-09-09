// package: envoy.api.v2.core
// file: github.com/solo-io/gloo/projects/gloo/api/external/envoy/api/v2/core/health_check.proto

import * as jspb from "google-protobuf";
import * as envoy_api_v2_core_base_pb from "../../../../../../../../../../../envoy/api/v2/core/base_pb";
import * as envoy_type_range_pb from "../../../../../../../../../../../envoy/type/range_pb";
import * as google_protobuf_any_pb from "google-protobuf/google/protobuf/any_pb";
import * as google_protobuf_duration_pb from "google-protobuf/google/protobuf/duration_pb";
import * as google_protobuf_struct_pb from "google-protobuf/google/protobuf/struct_pb";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";
import * as validate_validate_pb from "../../../../../../../../../../../validate/validate_pb";
import * as gogoproto_gogo_pb from "../../../../../../../../../../../gogoproto/gogo_pb";

export class HealthCheck extends jspb.Message {
  hasTimeout(): boolean;
  clearTimeout(): void;
  getTimeout(): google_protobuf_duration_pb.Duration | undefined;
  setTimeout(value?: google_protobuf_duration_pb.Duration): void;

  hasInterval(): boolean;
  clearInterval(): void;
  getInterval(): google_protobuf_duration_pb.Duration | undefined;
  setInterval(value?: google_protobuf_duration_pb.Duration): void;

  hasInitialJitter(): boolean;
  clearInitialJitter(): void;
  getInitialJitter(): google_protobuf_duration_pb.Duration | undefined;
  setInitialJitter(value?: google_protobuf_duration_pb.Duration): void;

  hasIntervalJitter(): boolean;
  clearIntervalJitter(): void;
  getIntervalJitter(): google_protobuf_duration_pb.Duration | undefined;
  setIntervalJitter(value?: google_protobuf_duration_pb.Duration): void;

  getIntervalJitterPercent(): number;
  setIntervalJitterPercent(value: number): void;

  hasUnhealthyThreshold(): boolean;
  clearUnhealthyThreshold(): void;
  getUnhealthyThreshold(): google_protobuf_wrappers_pb.UInt32Value | undefined;
  setUnhealthyThreshold(value?: google_protobuf_wrappers_pb.UInt32Value): void;

  hasHealthyThreshold(): boolean;
  clearHealthyThreshold(): void;
  getHealthyThreshold(): google_protobuf_wrappers_pb.UInt32Value | undefined;
  setHealthyThreshold(value?: google_protobuf_wrappers_pb.UInt32Value): void;

  hasReuseConnection(): boolean;
  clearReuseConnection(): void;
  getReuseConnection(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setReuseConnection(value?: google_protobuf_wrappers_pb.BoolValue): void;

  hasHttpHealthCheck(): boolean;
  clearHttpHealthCheck(): void;
  getHttpHealthCheck(): HealthCheck.HttpHealthCheck | undefined;
  setHttpHealthCheck(value?: HealthCheck.HttpHealthCheck): void;

  hasTcpHealthCheck(): boolean;
  clearTcpHealthCheck(): void;
  getTcpHealthCheck(): HealthCheck.TcpHealthCheck | undefined;
  setTcpHealthCheck(value?: HealthCheck.TcpHealthCheck): void;

  hasGrpcHealthCheck(): boolean;
  clearGrpcHealthCheck(): void;
  getGrpcHealthCheck(): HealthCheck.GrpcHealthCheck | undefined;
  setGrpcHealthCheck(value?: HealthCheck.GrpcHealthCheck): void;

  hasCustomHealthCheck(): boolean;
  clearCustomHealthCheck(): void;
  getCustomHealthCheck(): HealthCheck.CustomHealthCheck | undefined;
  setCustomHealthCheck(value?: HealthCheck.CustomHealthCheck): void;

  hasNoTrafficInterval(): boolean;
  clearNoTrafficInterval(): void;
  getNoTrafficInterval(): google_protobuf_duration_pb.Duration | undefined;
  setNoTrafficInterval(value?: google_protobuf_duration_pb.Duration): void;

  hasUnhealthyInterval(): boolean;
  clearUnhealthyInterval(): void;
  getUnhealthyInterval(): google_protobuf_duration_pb.Duration | undefined;
  setUnhealthyInterval(value?: google_protobuf_duration_pb.Duration): void;

  hasUnhealthyEdgeInterval(): boolean;
  clearUnhealthyEdgeInterval(): void;
  getUnhealthyEdgeInterval(): google_protobuf_duration_pb.Duration | undefined;
  setUnhealthyEdgeInterval(value?: google_protobuf_duration_pb.Duration): void;

  hasHealthyEdgeInterval(): boolean;
  clearHealthyEdgeInterval(): void;
  getHealthyEdgeInterval(): google_protobuf_duration_pb.Duration | undefined;
  setHealthyEdgeInterval(value?: google_protobuf_duration_pb.Duration): void;

  getEventLogPath(): string;
  setEventLogPath(value: string): void;

  getAlwaysLogHealthCheckFailures(): boolean;
  setAlwaysLogHealthCheckFailures(value: boolean): void;

  getHealthCheckerCase(): HealthCheck.HealthCheckerCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HealthCheck.AsObject;
  static toObject(includeInstance: boolean, msg: HealthCheck): HealthCheck.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: HealthCheck, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HealthCheck;
  static deserializeBinaryFromReader(message: HealthCheck, reader: jspb.BinaryReader): HealthCheck;
}

export namespace HealthCheck {
  export type AsObject = {
    timeout?: google_protobuf_duration_pb.Duration.AsObject,
    interval?: google_protobuf_duration_pb.Duration.AsObject,
    initialJitter?: google_protobuf_duration_pb.Duration.AsObject,
    intervalJitter?: google_protobuf_duration_pb.Duration.AsObject,
    intervalJitterPercent: number,
    unhealthyThreshold?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
    healthyThreshold?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
    reuseConnection?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    httpHealthCheck?: HealthCheck.HttpHealthCheck.AsObject,
    tcpHealthCheck?: HealthCheck.TcpHealthCheck.AsObject,
    grpcHealthCheck?: HealthCheck.GrpcHealthCheck.AsObject,
    customHealthCheck?: HealthCheck.CustomHealthCheck.AsObject,
    noTrafficInterval?: google_protobuf_duration_pb.Duration.AsObject,
    unhealthyInterval?: google_protobuf_duration_pb.Duration.AsObject,
    unhealthyEdgeInterval?: google_protobuf_duration_pb.Duration.AsObject,
    healthyEdgeInterval?: google_protobuf_duration_pb.Duration.AsObject,
    eventLogPath: string,
    alwaysLogHealthCheckFailures: boolean,
  }

  export class Payload extends jspb.Message {
    hasText(): boolean;
    clearText(): void;
    getText(): string;
    setText(value: string): void;

    getPayloadCase(): Payload.PayloadCase;
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Payload.AsObject;
    static toObject(includeInstance: boolean, msg: Payload): Payload.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Payload, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Payload;
    static deserializeBinaryFromReader(message: Payload, reader: jspb.BinaryReader): Payload;
  }

  export namespace Payload {
    export type AsObject = {
      text: string,
    }

    export enum PayloadCase {
      PAYLOAD_NOT_SET = 0,
      TEXT = 1,
    }
  }

  export class HttpHealthCheck extends jspb.Message {
    getHost(): string;
    setHost(value: string): void;

    getPath(): string;
    setPath(value: string): void;

    getServiceName(): string;
    setServiceName(value: string): void;

    clearRequestHeadersToAddList(): void;
    getRequestHeadersToAddList(): Array<envoy_api_v2_core_base_pb.HeaderValueOption>;
    setRequestHeadersToAddList(value: Array<envoy_api_v2_core_base_pb.HeaderValueOption>): void;
    addRequestHeadersToAdd(value?: envoy_api_v2_core_base_pb.HeaderValueOption, index?: number): envoy_api_v2_core_base_pb.HeaderValueOption;

    clearRequestHeadersToRemoveList(): void;
    getRequestHeadersToRemoveList(): Array<string>;
    setRequestHeadersToRemoveList(value: Array<string>): void;
    addRequestHeadersToRemove(value: string, index?: number): string;

    getUseHttp2(): boolean;
    setUseHttp2(value: boolean): void;

    clearExpectedStatusesList(): void;
    getExpectedStatusesList(): Array<envoy_type_range_pb.Int64Range>;
    setExpectedStatusesList(value: Array<envoy_type_range_pb.Int64Range>): void;
    addExpectedStatuses(value?: envoy_type_range_pb.Int64Range, index?: number): envoy_type_range_pb.Int64Range;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): HttpHealthCheck.AsObject;
    static toObject(includeInstance: boolean, msg: HttpHealthCheck): HttpHealthCheck.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: HttpHealthCheck, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): HttpHealthCheck;
    static deserializeBinaryFromReader(message: HttpHealthCheck, reader: jspb.BinaryReader): HttpHealthCheck;
  }

  export namespace HttpHealthCheck {
    export type AsObject = {
      host: string,
      path: string,
      serviceName: string,
      requestHeadersToAddList: Array<envoy_api_v2_core_base_pb.HeaderValueOption.AsObject>,
      requestHeadersToRemoveList: Array<string>,
      useHttp2: boolean,
      expectedStatusesList: Array<envoy_type_range_pb.Int64Range.AsObject>,
    }
  }

  export class TcpHealthCheck extends jspb.Message {
    hasSend(): boolean;
    clearSend(): void;
    getSend(): HealthCheck.Payload | undefined;
    setSend(value?: HealthCheck.Payload): void;

    clearReceiveList(): void;
    getReceiveList(): Array<HealthCheck.Payload>;
    setReceiveList(value: Array<HealthCheck.Payload>): void;
    addReceive(value?: HealthCheck.Payload, index?: number): HealthCheck.Payload;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): TcpHealthCheck.AsObject;
    static toObject(includeInstance: boolean, msg: TcpHealthCheck): TcpHealthCheck.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: TcpHealthCheck, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): TcpHealthCheck;
    static deserializeBinaryFromReader(message: TcpHealthCheck, reader: jspb.BinaryReader): TcpHealthCheck;
  }

  export namespace TcpHealthCheck {
    export type AsObject = {
      send?: HealthCheck.Payload.AsObject,
      receiveList: Array<HealthCheck.Payload.AsObject>,
    }
  }

  export class RedisHealthCheck extends jspb.Message {
    getKey(): string;
    setKey(value: string): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): RedisHealthCheck.AsObject;
    static toObject(includeInstance: boolean, msg: RedisHealthCheck): RedisHealthCheck.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: RedisHealthCheck, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): RedisHealthCheck;
    static deserializeBinaryFromReader(message: RedisHealthCheck, reader: jspb.BinaryReader): RedisHealthCheck;
  }

  export namespace RedisHealthCheck {
    export type AsObject = {
      key: string,
    }
  }

  export class GrpcHealthCheck extends jspb.Message {
    getServiceName(): string;
    setServiceName(value: string): void;

    getAuthority(): string;
    setAuthority(value: string): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): GrpcHealthCheck.AsObject;
    static toObject(includeInstance: boolean, msg: GrpcHealthCheck): GrpcHealthCheck.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: GrpcHealthCheck, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): GrpcHealthCheck;
    static deserializeBinaryFromReader(message: GrpcHealthCheck, reader: jspb.BinaryReader): GrpcHealthCheck;
  }

  export namespace GrpcHealthCheck {
    export type AsObject = {
      serviceName: string,
      authority: string,
    }
  }

  export class CustomHealthCheck extends jspb.Message {
    getName(): string;
    setName(value: string): void;

    hasConfig(): boolean;
    clearConfig(): void;
    getConfig(): google_protobuf_struct_pb.Struct | undefined;
    setConfig(value?: google_protobuf_struct_pb.Struct): void;

    hasTypedConfig(): boolean;
    clearTypedConfig(): void;
    getTypedConfig(): google_protobuf_any_pb.Any | undefined;
    setTypedConfig(value?: google_protobuf_any_pb.Any): void;

    getConfigTypeCase(): CustomHealthCheck.ConfigTypeCase;
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): CustomHealthCheck.AsObject;
    static toObject(includeInstance: boolean, msg: CustomHealthCheck): CustomHealthCheck.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: CustomHealthCheck, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): CustomHealthCheck;
    static deserializeBinaryFromReader(message: CustomHealthCheck, reader: jspb.BinaryReader): CustomHealthCheck;
  }

  export namespace CustomHealthCheck {
    export type AsObject = {
      name: string,
      config?: google_protobuf_struct_pb.Struct.AsObject,
      typedConfig?: google_protobuf_any_pb.Any.AsObject,
    }

    export enum ConfigTypeCase {
      CONFIG_TYPE_NOT_SET = 0,
      CONFIG = 2,
      TYPED_CONFIG = 3,
    }
  }

  export enum HealthCheckerCase {
    HEALTH_CHECKER_NOT_SET = 0,
    HTTP_HEALTH_CHECK = 8,
    TCP_HEALTH_CHECK = 9,
    GRPC_HEALTH_CHECK = 11,
    CUSTOM_HEALTH_CHECK = 13,
  }
}

export interface HealthStatusMap {
  UNKNOWN: 0;
  HEALTHY: 1;
  UNHEALTHY: 2;
  DRAINING: 3;
  TIMEOUT: 4;
  DEGRADED: 5;
}

export const HealthStatus: HealthStatusMap;

