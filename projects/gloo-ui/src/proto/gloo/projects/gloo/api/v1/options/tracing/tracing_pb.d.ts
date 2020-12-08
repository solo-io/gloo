/* eslint-disable */
// package: tracing.options.gloo.solo.io
// file: gloo/projects/gloo/api/v1/options/tracing/tracing.proto

import * as jspb from "google-protobuf";
import * as gloo_projects_gloo_api_external_envoy_config_trace_v3_zipkin_pb from "../../../../../../../gloo/projects/gloo/api/external/envoy/config/trace/v3/zipkin_pb";
import * as gloo_projects_gloo_api_external_envoy_config_trace_v3_datadog_pb from "../../../../../../../gloo/projects/gloo/api/external/envoy/config/trace/v3/datadog_pb";
import * as gogoproto_gogo_pb from "../../../../../../../gogoproto/gogo_pb";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";
import * as google_protobuf_duration_pb from "google-protobuf/google/protobuf/duration_pb";
import * as google_protobuf_any_pb from "google-protobuf/google/protobuf/any_pb";
import * as extproto_ext_pb from "../../../../../../../protoc-gen-ext/extproto/ext_pb";
import * as solo_kit_api_v1_ref_pb from "../../../../../../../solo-kit/api/v1/ref_pb";

export class ListenerTracingSettings extends jspb.Message {
  clearRequestHeadersForTagsList(): void;
  getRequestHeadersForTagsList(): Array<string>;
  setRequestHeadersForTagsList(value: Array<string>): void;
  addRequestHeadersForTags(value: string, index?: number): string;

  getVerbose(): boolean;
  setVerbose(value: boolean): void;

  hasTracePercentages(): boolean;
  clearTracePercentages(): void;
  getTracePercentages(): TracePercentages | undefined;
  setTracePercentages(value?: TracePercentages): void;

  hasZipkinConfig(): boolean;
  clearZipkinConfig(): void;
  getZipkinConfig(): gloo_projects_gloo_api_external_envoy_config_trace_v3_zipkin_pb.ZipkinConfig | undefined;
  setZipkinConfig(value?: gloo_projects_gloo_api_external_envoy_config_trace_v3_zipkin_pb.ZipkinConfig): void;

  hasDatadogConfig(): boolean;
  clearDatadogConfig(): void;
  getDatadogConfig(): gloo_projects_gloo_api_external_envoy_config_trace_v3_datadog_pb.DatadogConfig | undefined;
  setDatadogConfig(value?: gloo_projects_gloo_api_external_envoy_config_trace_v3_datadog_pb.DatadogConfig): void;

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
    requestHeadersForTagsList: Array<string>,
    verbose: boolean,
    tracePercentages?: TracePercentages.AsObject,
    zipkinConfig?: gloo_projects_gloo_api_external_envoy_config_trace_v3_zipkin_pb.ZipkinConfig.AsObject,
    datadogConfig?: gloo_projects_gloo_api_external_envoy_config_trace_v3_datadog_pb.DatadogConfig.AsObject,
  }

  export enum ProviderConfigCase {
    PROVIDER_CONFIG_NOT_SET = 0,
    ZIPKIN_CONFIG = 4,
    DATADOG_CONFIG = 5,
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
