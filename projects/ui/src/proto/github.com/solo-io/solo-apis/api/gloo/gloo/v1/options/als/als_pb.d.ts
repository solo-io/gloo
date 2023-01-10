/* eslint-disable */
// package: als.options.gloo.solo.io
// file: github.com/solo-io/solo-apis/api/gloo/gloo/v1/options/als/als.proto

import * as jspb from "google-protobuf";
import * as extproto_ext_pb from "../../../../../../../../../extproto/ext_pb";
import * as github_com_solo_io_solo_kit_api_v1_ref_pb from "../../../../../../../../../github.com/solo-io/solo-kit/api/v1/ref_pb";
import * as validate_validate_pb from "../../../../../../../../../validate/validate_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_base_pb from "../../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/config/core/v3/base_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_type_v3_percent_pb from "../../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/type/v3/percent_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_route_v3_route_components_pb from "../../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/config/route/v3/route_components_pb";
import * as google_protobuf_struct_pb from "google-protobuf/google/protobuf/struct_pb";

export class AccessLoggingService extends jspb.Message {
  clearAccessLogList(): void;
  getAccessLogList(): Array<AccessLog>;
  setAccessLogList(value: Array<AccessLog>): void;
  addAccessLog(value?: AccessLog, index?: number): AccessLog;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AccessLoggingService.AsObject;
  static toObject(includeInstance: boolean, msg: AccessLoggingService): AccessLoggingService.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: AccessLoggingService, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AccessLoggingService;
  static deserializeBinaryFromReader(message: AccessLoggingService, reader: jspb.BinaryReader): AccessLoggingService;
}

export namespace AccessLoggingService {
  export type AsObject = {
    accessLogList: Array<AccessLog.AsObject>,
  }
}

export class AccessLog extends jspb.Message {
  hasFileSink(): boolean;
  clearFileSink(): void;
  getFileSink(): FileSink | undefined;
  setFileSink(value?: FileSink): void;

  hasGrpcService(): boolean;
  clearGrpcService(): void;
  getGrpcService(): GrpcService | undefined;
  setGrpcService(value?: GrpcService): void;

  hasFilter(): boolean;
  clearFilter(): void;
  getFilter(): AccessLogFilter | undefined;
  setFilter(value?: AccessLogFilter): void;

  getOutputdestinationCase(): AccessLog.OutputdestinationCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AccessLog.AsObject;
  static toObject(includeInstance: boolean, msg: AccessLog): AccessLog.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: AccessLog, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AccessLog;
  static deserializeBinaryFromReader(message: AccessLog, reader: jspb.BinaryReader): AccessLog;
}

export namespace AccessLog {
  export type AsObject = {
    fileSink?: FileSink.AsObject,
    grpcService?: GrpcService.AsObject,
    filter?: AccessLogFilter.AsObject,
  }

  export enum OutputdestinationCase {
    OUTPUTDESTINATION_NOT_SET = 0,
    FILE_SINK = 2,
    GRPC_SERVICE = 3,
  }
}

export class FileSink extends jspb.Message {
  getPath(): string;
  setPath(value: string): void;

  hasStringFormat(): boolean;
  clearStringFormat(): void;
  getStringFormat(): string;
  setStringFormat(value: string): void;

  hasJsonFormat(): boolean;
  clearJsonFormat(): void;
  getJsonFormat(): google_protobuf_struct_pb.Struct | undefined;
  setJsonFormat(value?: google_protobuf_struct_pb.Struct): void;

  getOutputFormatCase(): FileSink.OutputFormatCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FileSink.AsObject;
  static toObject(includeInstance: boolean, msg: FileSink): FileSink.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: FileSink, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FileSink;
  static deserializeBinaryFromReader(message: FileSink, reader: jspb.BinaryReader): FileSink;
}

export namespace FileSink {
  export type AsObject = {
    path: string,
    stringFormat: string,
    jsonFormat?: google_protobuf_struct_pb.Struct.AsObject,
  }

  export enum OutputFormatCase {
    OUTPUT_FORMAT_NOT_SET = 0,
    STRING_FORMAT = 2,
    JSON_FORMAT = 3,
  }
}

export class GrpcService extends jspb.Message {
  getLogName(): string;
  setLogName(value: string): void;

  hasStaticClusterName(): boolean;
  clearStaticClusterName(): void;
  getStaticClusterName(): string;
  setStaticClusterName(value: string): void;

  clearAdditionalRequestHeadersToLogList(): void;
  getAdditionalRequestHeadersToLogList(): Array<string>;
  setAdditionalRequestHeadersToLogList(value: Array<string>): void;
  addAdditionalRequestHeadersToLog(value: string, index?: number): string;

  clearAdditionalResponseHeadersToLogList(): void;
  getAdditionalResponseHeadersToLogList(): Array<string>;
  setAdditionalResponseHeadersToLogList(value: Array<string>): void;
  addAdditionalResponseHeadersToLog(value: string, index?: number): string;

  clearAdditionalResponseTrailersToLogList(): void;
  getAdditionalResponseTrailersToLogList(): Array<string>;
  setAdditionalResponseTrailersToLogList(value: Array<string>): void;
  addAdditionalResponseTrailersToLog(value: string, index?: number): string;

  getServiceRefCase(): GrpcService.ServiceRefCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GrpcService.AsObject;
  static toObject(includeInstance: boolean, msg: GrpcService): GrpcService.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GrpcService, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GrpcService;
  static deserializeBinaryFromReader(message: GrpcService, reader: jspb.BinaryReader): GrpcService;
}

export namespace GrpcService {
  export type AsObject = {
    logName: string,
    staticClusterName: string,
    additionalRequestHeadersToLogList: Array<string>,
    additionalResponseHeadersToLogList: Array<string>,
    additionalResponseTrailersToLogList: Array<string>,
  }

  export enum ServiceRefCase {
    SERVICE_REF_NOT_SET = 0,
    STATIC_CLUSTER_NAME = 2,
  }
}

export class AccessLogFilter extends jspb.Message {
  hasStatusCodeFilter(): boolean;
  clearStatusCodeFilter(): void;
  getStatusCodeFilter(): StatusCodeFilter | undefined;
  setStatusCodeFilter(value?: StatusCodeFilter): void;

  hasDurationFilter(): boolean;
  clearDurationFilter(): void;
  getDurationFilter(): DurationFilter | undefined;
  setDurationFilter(value?: DurationFilter): void;

  hasNotHealthCheckFilter(): boolean;
  clearNotHealthCheckFilter(): void;
  getNotHealthCheckFilter(): NotHealthCheckFilter | undefined;
  setNotHealthCheckFilter(value?: NotHealthCheckFilter): void;

  hasTraceableFilter(): boolean;
  clearTraceableFilter(): void;
  getTraceableFilter(): TraceableFilter | undefined;
  setTraceableFilter(value?: TraceableFilter): void;

  hasRuntimeFilter(): boolean;
  clearRuntimeFilter(): void;
  getRuntimeFilter(): RuntimeFilter | undefined;
  setRuntimeFilter(value?: RuntimeFilter): void;

  hasAndFilter(): boolean;
  clearAndFilter(): void;
  getAndFilter(): AndFilter | undefined;
  setAndFilter(value?: AndFilter): void;

  hasOrFilter(): boolean;
  clearOrFilter(): void;
  getOrFilter(): OrFilter | undefined;
  setOrFilter(value?: OrFilter): void;

  hasHeaderFilter(): boolean;
  clearHeaderFilter(): void;
  getHeaderFilter(): HeaderFilter | undefined;
  setHeaderFilter(value?: HeaderFilter): void;

  hasResponseFlagFilter(): boolean;
  clearResponseFlagFilter(): void;
  getResponseFlagFilter(): ResponseFlagFilter | undefined;
  setResponseFlagFilter(value?: ResponseFlagFilter): void;

  hasGrpcStatusFilter(): boolean;
  clearGrpcStatusFilter(): void;
  getGrpcStatusFilter(): GrpcStatusFilter | undefined;
  setGrpcStatusFilter(value?: GrpcStatusFilter): void;

  getFilterSpecifierCase(): AccessLogFilter.FilterSpecifierCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AccessLogFilter.AsObject;
  static toObject(includeInstance: boolean, msg: AccessLogFilter): AccessLogFilter.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: AccessLogFilter, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AccessLogFilter;
  static deserializeBinaryFromReader(message: AccessLogFilter, reader: jspb.BinaryReader): AccessLogFilter;
}

export namespace AccessLogFilter {
  export type AsObject = {
    statusCodeFilter?: StatusCodeFilter.AsObject,
    durationFilter?: DurationFilter.AsObject,
    notHealthCheckFilter?: NotHealthCheckFilter.AsObject,
    traceableFilter?: TraceableFilter.AsObject,
    runtimeFilter?: RuntimeFilter.AsObject,
    andFilter?: AndFilter.AsObject,
    orFilter?: OrFilter.AsObject,
    headerFilter?: HeaderFilter.AsObject,
    responseFlagFilter?: ResponseFlagFilter.AsObject,
    grpcStatusFilter?: GrpcStatusFilter.AsObject,
  }

  export enum FilterSpecifierCase {
    FILTER_SPECIFIER_NOT_SET = 0,
    STATUS_CODE_FILTER = 1,
    DURATION_FILTER = 2,
    NOT_HEALTH_CHECK_FILTER = 3,
    TRACEABLE_FILTER = 4,
    RUNTIME_FILTER = 5,
    AND_FILTER = 6,
    OR_FILTER = 7,
    HEADER_FILTER = 8,
    RESPONSE_FLAG_FILTER = 9,
    GRPC_STATUS_FILTER = 10,
  }
}

export class ComparisonFilter extends jspb.Message {
  getOp(): ComparisonFilter.OpMap[keyof ComparisonFilter.OpMap];
  setOp(value: ComparisonFilter.OpMap[keyof ComparisonFilter.OpMap]): void;

  hasValue(): boolean;
  clearValue(): void;
  getValue(): github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_base_pb.RuntimeUInt32 | undefined;
  setValue(value?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_base_pb.RuntimeUInt32): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ComparisonFilter.AsObject;
  static toObject(includeInstance: boolean, msg: ComparisonFilter): ComparisonFilter.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ComparisonFilter, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ComparisonFilter;
  static deserializeBinaryFromReader(message: ComparisonFilter, reader: jspb.BinaryReader): ComparisonFilter;
}

export namespace ComparisonFilter {
  export type AsObject = {
    op: ComparisonFilter.OpMap[keyof ComparisonFilter.OpMap],
    value?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_base_pb.RuntimeUInt32.AsObject,
  }

  export interface OpMap {
    EQ: 0;
    GE: 1;
    LE: 2;
  }

  export const Op: OpMap;
}

export class StatusCodeFilter extends jspb.Message {
  hasComparison(): boolean;
  clearComparison(): void;
  getComparison(): ComparisonFilter | undefined;
  setComparison(value?: ComparisonFilter): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StatusCodeFilter.AsObject;
  static toObject(includeInstance: boolean, msg: StatusCodeFilter): StatusCodeFilter.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: StatusCodeFilter, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StatusCodeFilter;
  static deserializeBinaryFromReader(message: StatusCodeFilter, reader: jspb.BinaryReader): StatusCodeFilter;
}

export namespace StatusCodeFilter {
  export type AsObject = {
    comparison?: ComparisonFilter.AsObject,
  }
}

export class DurationFilter extends jspb.Message {
  hasComparison(): boolean;
  clearComparison(): void;
  getComparison(): ComparisonFilter | undefined;
  setComparison(value?: ComparisonFilter): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DurationFilter.AsObject;
  static toObject(includeInstance: boolean, msg: DurationFilter): DurationFilter.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DurationFilter, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DurationFilter;
  static deserializeBinaryFromReader(message: DurationFilter, reader: jspb.BinaryReader): DurationFilter;
}

export namespace DurationFilter {
  export type AsObject = {
    comparison?: ComparisonFilter.AsObject,
  }
}

export class NotHealthCheckFilter extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): NotHealthCheckFilter.AsObject;
  static toObject(includeInstance: boolean, msg: NotHealthCheckFilter): NotHealthCheckFilter.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: NotHealthCheckFilter, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): NotHealthCheckFilter;
  static deserializeBinaryFromReader(message: NotHealthCheckFilter, reader: jspb.BinaryReader): NotHealthCheckFilter;
}

export namespace NotHealthCheckFilter {
  export type AsObject = {
  }
}

export class TraceableFilter extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TraceableFilter.AsObject;
  static toObject(includeInstance: boolean, msg: TraceableFilter): TraceableFilter.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: TraceableFilter, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TraceableFilter;
  static deserializeBinaryFromReader(message: TraceableFilter, reader: jspb.BinaryReader): TraceableFilter;
}

export namespace TraceableFilter {
  export type AsObject = {
  }
}

export class RuntimeFilter extends jspb.Message {
  getRuntimeKey(): string;
  setRuntimeKey(value: string): void;

  hasPercentSampled(): boolean;
  clearPercentSampled(): void;
  getPercentSampled(): github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_type_v3_percent_pb.FractionalPercent | undefined;
  setPercentSampled(value?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_type_v3_percent_pb.FractionalPercent): void;

  getUseIndependentRandomness(): boolean;
  setUseIndependentRandomness(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RuntimeFilter.AsObject;
  static toObject(includeInstance: boolean, msg: RuntimeFilter): RuntimeFilter.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RuntimeFilter, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RuntimeFilter;
  static deserializeBinaryFromReader(message: RuntimeFilter, reader: jspb.BinaryReader): RuntimeFilter;
}

export namespace RuntimeFilter {
  export type AsObject = {
    runtimeKey: string,
    percentSampled?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_type_v3_percent_pb.FractionalPercent.AsObject,
    useIndependentRandomness: boolean,
  }
}

export class AndFilter extends jspb.Message {
  clearFiltersList(): void;
  getFiltersList(): Array<AccessLogFilter>;
  setFiltersList(value: Array<AccessLogFilter>): void;
  addFilters(value?: AccessLogFilter, index?: number): AccessLogFilter;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AndFilter.AsObject;
  static toObject(includeInstance: boolean, msg: AndFilter): AndFilter.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: AndFilter, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AndFilter;
  static deserializeBinaryFromReader(message: AndFilter, reader: jspb.BinaryReader): AndFilter;
}

export namespace AndFilter {
  export type AsObject = {
    filtersList: Array<AccessLogFilter.AsObject>,
  }
}

export class OrFilter extends jspb.Message {
  clearFiltersList(): void;
  getFiltersList(): Array<AccessLogFilter>;
  setFiltersList(value: Array<AccessLogFilter>): void;
  addFilters(value?: AccessLogFilter, index?: number): AccessLogFilter;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OrFilter.AsObject;
  static toObject(includeInstance: boolean, msg: OrFilter): OrFilter.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: OrFilter, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OrFilter;
  static deserializeBinaryFromReader(message: OrFilter, reader: jspb.BinaryReader): OrFilter;
}

export namespace OrFilter {
  export type AsObject = {
    filtersList: Array<AccessLogFilter.AsObject>,
  }
}

export class HeaderFilter extends jspb.Message {
  hasHeader(): boolean;
  clearHeader(): void;
  getHeader(): github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_route_v3_route_components_pb.HeaderMatcher | undefined;
  setHeader(value?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_route_v3_route_components_pb.HeaderMatcher): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HeaderFilter.AsObject;
  static toObject(includeInstance: boolean, msg: HeaderFilter): HeaderFilter.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: HeaderFilter, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HeaderFilter;
  static deserializeBinaryFromReader(message: HeaderFilter, reader: jspb.BinaryReader): HeaderFilter;
}

export namespace HeaderFilter {
  export type AsObject = {
    header?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_route_v3_route_components_pb.HeaderMatcher.AsObject,
  }
}

export class ResponseFlagFilter extends jspb.Message {
  clearFlagsList(): void;
  getFlagsList(): Array<string>;
  setFlagsList(value: Array<string>): void;
  addFlags(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResponseFlagFilter.AsObject;
  static toObject(includeInstance: boolean, msg: ResponseFlagFilter): ResponseFlagFilter.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ResponseFlagFilter, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResponseFlagFilter;
  static deserializeBinaryFromReader(message: ResponseFlagFilter, reader: jspb.BinaryReader): ResponseFlagFilter;
}

export namespace ResponseFlagFilter {
  export type AsObject = {
    flagsList: Array<string>,
  }
}

export class GrpcStatusFilter extends jspb.Message {
  clearStatusesList(): void;
  getStatusesList(): Array<GrpcStatusFilter.StatusMap[keyof GrpcStatusFilter.StatusMap]>;
  setStatusesList(value: Array<GrpcStatusFilter.StatusMap[keyof GrpcStatusFilter.StatusMap]>): void;
  addStatuses(value: GrpcStatusFilter.StatusMap[keyof GrpcStatusFilter.StatusMap], index?: number): GrpcStatusFilter.StatusMap[keyof GrpcStatusFilter.StatusMap];

  getExclude(): boolean;
  setExclude(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GrpcStatusFilter.AsObject;
  static toObject(includeInstance: boolean, msg: GrpcStatusFilter): GrpcStatusFilter.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GrpcStatusFilter, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GrpcStatusFilter;
  static deserializeBinaryFromReader(message: GrpcStatusFilter, reader: jspb.BinaryReader): GrpcStatusFilter;
}

export namespace GrpcStatusFilter {
  export type AsObject = {
    statusesList: Array<GrpcStatusFilter.StatusMap[keyof GrpcStatusFilter.StatusMap]>,
    exclude: boolean,
  }

  export interface StatusMap {
    OK: 0;
    CANCELED: 1;
    UNKNOWN: 2;
    INVALID_ARGUMENT: 3;
    DEADLINE_EXCEEDED: 4;
    NOT_FOUND: 5;
    ALREADY_EXISTS: 6;
    PERMISSION_DENIED: 7;
    RESOURCE_EXHAUSTED: 8;
    FAILED_PRECONDITION: 9;
    ABORTED: 10;
    OUT_OF_RANGE: 11;
    UNIMPLEMENTED: 12;
    INTERNAL: 13;
    UNAVAILABLE: 14;
    DATA_LOSS: 15;
    UNAUTHENTICATED: 16;
  }

  export const Status: StatusMap;
}
