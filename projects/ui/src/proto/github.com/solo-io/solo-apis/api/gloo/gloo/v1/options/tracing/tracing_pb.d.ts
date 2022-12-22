/* eslint-disable */
// package: tracing.options.gloo.solo.io
// file: github.com/solo-io/solo-apis/api/gloo/gloo/v1/options/tracing/tracing.proto

import * as jspb from "google-protobuf";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_trace_v3_zipkin_pb from "../../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/config/trace/v3/zipkin_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_trace_v3_datadog_pb from "../../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/config/trace/v3/datadog_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_trace_v3_opentelemetry_pb from "../../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/config/trace/v3/opentelemetry_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_trace_v3_opencensus_pb from "../../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/config/trace/v3/opencensus_pb";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";
import * as google_protobuf_any_pb from "google-protobuf/google/protobuf/any_pb";
import * as extproto_ext_pb from "../../../../../../../../../extproto/ext_pb";
import * as github_com_solo_io_solo_kit_api_v1_ref_pb from "../../../../../../../../../github.com/solo-io/solo-kit/api/v1/ref_pb";

export class ListenerTracingSettings extends jspb.Message {
  clearRequestHeadersForTagsList(): void;
  getRequestHeadersForTagsList(): Array<google_protobuf_wrappers_pb.StringValue>;
  setRequestHeadersForTagsList(value: Array<google_protobuf_wrappers_pb.StringValue>): void;
  addRequestHeadersForTags(value?: google_protobuf_wrappers_pb.StringValue, index?: number): google_protobuf_wrappers_pb.StringValue;

  hasVerbose(): boolean;
  clearVerbose(): void;
  getVerbose(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setVerbose(value?: google_protobuf_wrappers_pb.BoolValue): void;

  hasTracePercentages(): boolean;
  clearTracePercentages(): void;
  getTracePercentages(): TracePercentages | undefined;
  setTracePercentages(value?: TracePercentages): void;

  hasZipkinConfig(): boolean;
  clearZipkinConfig(): void;
  getZipkinConfig(): github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_trace_v3_zipkin_pb.ZipkinConfig | undefined;
  setZipkinConfig(value?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_trace_v3_zipkin_pb.ZipkinConfig): void;

  hasDatadogConfig(): boolean;
  clearDatadogConfig(): void;
  getDatadogConfig(): github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_trace_v3_datadog_pb.DatadogConfig | undefined;
  setDatadogConfig(value?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_trace_v3_datadog_pb.DatadogConfig): void;

  hasOpenTelemetryConfig(): boolean;
  clearOpenTelemetryConfig(): void;
  getOpenTelemetryConfig(): github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_trace_v3_opentelemetry_pb.OpenTelemetryConfig | undefined;
  setOpenTelemetryConfig(value?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_trace_v3_opentelemetry_pb.OpenTelemetryConfig): void;

  hasOpenCensusConfig(): boolean;
  clearOpenCensusConfig(): void;
  getOpenCensusConfig(): github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_trace_v3_opencensus_pb.OpenCensusConfig | undefined;
  setOpenCensusConfig(value?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_trace_v3_opencensus_pb.OpenCensusConfig): void;

  clearEnvironmentVariablesForTagsList(): void;
  getEnvironmentVariablesForTagsList(): Array<TracingTagEnvironmentVariable>;
  setEnvironmentVariablesForTagsList(value: Array<TracingTagEnvironmentVariable>): void;
  addEnvironmentVariablesForTags(value?: TracingTagEnvironmentVariable, index?: number): TracingTagEnvironmentVariable;

  clearLiteralsForTagsList(): void;
  getLiteralsForTagsList(): Array<TracingTagLiteral>;
  setLiteralsForTagsList(value: Array<TracingTagLiteral>): void;
  addLiteralsForTags(value?: TracingTagLiteral, index?: number): TracingTagLiteral;

  getProviderConfigCase(): ListenerTracingSettings.ProviderConfigCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListenerTracingSettings.AsObject;
  static toObject(includeInstance: boolean, msg: ListenerTracingSettings): ListenerTracingSettings.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListenerTracingSettings, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListenerTracingSettings;
  static deserializeBinaryFromReader(message: ListenerTracingSettings, reader: jspb.BinaryReader): ListenerTracingSettings;
}

export namespace ListenerTracingSettings {
  export type AsObject = {
    requestHeadersForTagsList: Array<google_protobuf_wrappers_pb.StringValue.AsObject>,
    verbose?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    tracePercentages?: TracePercentages.AsObject,
    zipkinConfig?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_trace_v3_zipkin_pb.ZipkinConfig.AsObject,
    datadogConfig?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_trace_v3_datadog_pb.DatadogConfig.AsObject,
    openTelemetryConfig?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_trace_v3_opentelemetry_pb.OpenTelemetryConfig.AsObject,
    openCensusConfig?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_trace_v3_opencensus_pb.OpenCensusConfig.AsObject,
    environmentVariablesForTagsList: Array<TracingTagEnvironmentVariable.AsObject>,
    literalsForTagsList: Array<TracingTagLiteral.AsObject>,
  }

  export enum ProviderConfigCase {
    PROVIDER_CONFIG_NOT_SET = 0,
    ZIPKIN_CONFIG = 4,
    DATADOG_CONFIG = 5,
    OPEN_TELEMETRY_CONFIG = 8,
    OPEN_CENSUS_CONFIG = 9,
  }
}

export class RouteTracingSettings extends jspb.Message {
  getRouteDescriptor(): string;
  setRouteDescriptor(value: string): void;

  hasTracePercentages(): boolean;
  clearTracePercentages(): void;
  getTracePercentages(): TracePercentages | undefined;
  setTracePercentages(value?: TracePercentages): void;

  hasPropagate(): boolean;
  clearPropagate(): void;
  getPropagate(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setPropagate(value?: google_protobuf_wrappers_pb.BoolValue): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RouteTracingSettings.AsObject;
  static toObject(includeInstance: boolean, msg: RouteTracingSettings): RouteTracingSettings.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RouteTracingSettings, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RouteTracingSettings;
  static deserializeBinaryFromReader(message: RouteTracingSettings, reader: jspb.BinaryReader): RouteTracingSettings;
}

export namespace RouteTracingSettings {
  export type AsObject = {
    routeDescriptor: string,
    tracePercentages?: TracePercentages.AsObject,
    propagate?: google_protobuf_wrappers_pb.BoolValue.AsObject,
  }
}

export class TracePercentages extends jspb.Message {
  hasClientSamplePercentage(): boolean;
  clearClientSamplePercentage(): void;
  getClientSamplePercentage(): google_protobuf_wrappers_pb.FloatValue | undefined;
  setClientSamplePercentage(value?: google_protobuf_wrappers_pb.FloatValue): void;

  hasRandomSamplePercentage(): boolean;
  clearRandomSamplePercentage(): void;
  getRandomSamplePercentage(): google_protobuf_wrappers_pb.FloatValue | undefined;
  setRandomSamplePercentage(value?: google_protobuf_wrappers_pb.FloatValue): void;

  hasOverallSamplePercentage(): boolean;
  clearOverallSamplePercentage(): void;
  getOverallSamplePercentage(): google_protobuf_wrappers_pb.FloatValue | undefined;
  setOverallSamplePercentage(value?: google_protobuf_wrappers_pb.FloatValue): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TracePercentages.AsObject;
  static toObject(includeInstance: boolean, msg: TracePercentages): TracePercentages.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: TracePercentages, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TracePercentages;
  static deserializeBinaryFromReader(message: TracePercentages, reader: jspb.BinaryReader): TracePercentages;
}

export namespace TracePercentages {
  export type AsObject = {
    clientSamplePercentage?: google_protobuf_wrappers_pb.FloatValue.AsObject,
    randomSamplePercentage?: google_protobuf_wrappers_pb.FloatValue.AsObject,
    overallSamplePercentage?: google_protobuf_wrappers_pb.FloatValue.AsObject,
  }
}

export class TracingTagEnvironmentVariable extends jspb.Message {
  hasTag(): boolean;
  clearTag(): void;
  getTag(): google_protobuf_wrappers_pb.StringValue | undefined;
  setTag(value?: google_protobuf_wrappers_pb.StringValue): void;

  hasName(): boolean;
  clearName(): void;
  getName(): google_protobuf_wrappers_pb.StringValue | undefined;
  setName(value?: google_protobuf_wrappers_pb.StringValue): void;

  hasDefaultValue(): boolean;
  clearDefaultValue(): void;
  getDefaultValue(): google_protobuf_wrappers_pb.StringValue | undefined;
  setDefaultValue(value?: google_protobuf_wrappers_pb.StringValue): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TracingTagEnvironmentVariable.AsObject;
  static toObject(includeInstance: boolean, msg: TracingTagEnvironmentVariable): TracingTagEnvironmentVariable.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: TracingTagEnvironmentVariable, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TracingTagEnvironmentVariable;
  static deserializeBinaryFromReader(message: TracingTagEnvironmentVariable, reader: jspb.BinaryReader): TracingTagEnvironmentVariable;
}

export namespace TracingTagEnvironmentVariable {
  export type AsObject = {
    tag?: google_protobuf_wrappers_pb.StringValue.AsObject,
    name?: google_protobuf_wrappers_pb.StringValue.AsObject,
    defaultValue?: google_protobuf_wrappers_pb.StringValue.AsObject,
  }
}

export class TracingTagLiteral extends jspb.Message {
  hasTag(): boolean;
  clearTag(): void;
  getTag(): google_protobuf_wrappers_pb.StringValue | undefined;
  setTag(value?: google_protobuf_wrappers_pb.StringValue): void;

  hasValue(): boolean;
  clearValue(): void;
  getValue(): google_protobuf_wrappers_pb.StringValue | undefined;
  setValue(value?: google_protobuf_wrappers_pb.StringValue): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TracingTagLiteral.AsObject;
  static toObject(includeInstance: boolean, msg: TracingTagLiteral): TracingTagLiteral.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: TracingTagLiteral, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TracingTagLiteral;
  static deserializeBinaryFromReader(message: TracingTagLiteral, reader: jspb.BinaryReader): TracingTagLiteral;
}

export namespace TracingTagLiteral {
  export type AsObject = {
    tag?: google_protobuf_wrappers_pb.StringValue.AsObject,
    value?: google_protobuf_wrappers_pb.StringValue.AsObject,
  }
}
