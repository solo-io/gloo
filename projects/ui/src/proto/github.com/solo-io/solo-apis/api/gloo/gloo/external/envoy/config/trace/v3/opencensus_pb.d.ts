/* eslint-disable */
// package: solo.io.envoy.config.trace.v3
// file: github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/config/trace/v3/opencensus.proto

import * as jspb from "google-protobuf";
import * as udpa_annotations_migrate_pb from "../../../../../../../../../../../udpa/annotations/migrate_pb";
import * as udpa_annotations_status_pb from "../../../../../../../../../../../udpa/annotations/status_pb";
import * as udpa_annotations_versioning_pb from "../../../../../../../../../../../udpa/annotations/versioning_pb";
import * as validate_validate_pb from "../../../../../../../../../../../validate/validate_pb";
import * as extproto_ext_pb from "../../../../../../../../../../../extproto/ext_pb";

export class OpenCensusConfig extends jspb.Message {
  hasTraceConfig(): boolean;
  clearTraceConfig(): void;
  getTraceConfig(): TraceConfig | undefined;
  setTraceConfig(value?: TraceConfig): void;

  getOcagentExporterEnabled(): boolean;
  setOcagentExporterEnabled(value: boolean): void;

  hasHttpAddress(): boolean;
  clearHttpAddress(): void;
  getHttpAddress(): string;
  setHttpAddress(value: string): void;

  hasGrpcAddress(): boolean;
  clearGrpcAddress(): void;
  getGrpcAddress(): OpenCensusConfig.OcagentGrpcAddress | undefined;
  setGrpcAddress(value?: OpenCensusConfig.OcagentGrpcAddress): void;

  clearIncomingTraceContextList(): void;
  getIncomingTraceContextList(): Array<OpenCensusConfig.TraceContextMap[keyof OpenCensusConfig.TraceContextMap]>;
  setIncomingTraceContextList(value: Array<OpenCensusConfig.TraceContextMap[keyof OpenCensusConfig.TraceContextMap]>): void;
  addIncomingTraceContext(value: OpenCensusConfig.TraceContextMap[keyof OpenCensusConfig.TraceContextMap], index?: number): OpenCensusConfig.TraceContextMap[keyof OpenCensusConfig.TraceContextMap];

  clearOutgoingTraceContextList(): void;
  getOutgoingTraceContextList(): Array<OpenCensusConfig.TraceContextMap[keyof OpenCensusConfig.TraceContextMap]>;
  setOutgoingTraceContextList(value: Array<OpenCensusConfig.TraceContextMap[keyof OpenCensusConfig.TraceContextMap]>): void;
  addOutgoingTraceContext(value: OpenCensusConfig.TraceContextMap[keyof OpenCensusConfig.TraceContextMap], index?: number): OpenCensusConfig.TraceContextMap[keyof OpenCensusConfig.TraceContextMap];

  getOcagentAddressCase(): OpenCensusConfig.OcagentAddressCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OpenCensusConfig.AsObject;
  static toObject(includeInstance: boolean, msg: OpenCensusConfig): OpenCensusConfig.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: OpenCensusConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OpenCensusConfig;
  static deserializeBinaryFromReader(message: OpenCensusConfig, reader: jspb.BinaryReader): OpenCensusConfig;
}

export namespace OpenCensusConfig {
  export type AsObject = {
    traceConfig?: TraceConfig.AsObject,
    ocagentExporterEnabled: boolean,
    httpAddress: string,
    grpcAddress?: OpenCensusConfig.OcagentGrpcAddress.AsObject,
    incomingTraceContextList: Array<OpenCensusConfig.TraceContextMap[keyof OpenCensusConfig.TraceContextMap]>,
    outgoingTraceContextList: Array<OpenCensusConfig.TraceContextMap[keyof OpenCensusConfig.TraceContextMap]>,
  }

  export class OcagentGrpcAddress extends jspb.Message {
    getTargetUri(): string;
    setTargetUri(value: string): void;

    getStatPrefix(): string;
    setStatPrefix(value: string): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): OcagentGrpcAddress.AsObject;
    static toObject(includeInstance: boolean, msg: OcagentGrpcAddress): OcagentGrpcAddress.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: OcagentGrpcAddress, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): OcagentGrpcAddress;
    static deserializeBinaryFromReader(message: OcagentGrpcAddress, reader: jspb.BinaryReader): OcagentGrpcAddress;
  }

  export namespace OcagentGrpcAddress {
    export type AsObject = {
      targetUri: string,
      statPrefix: string,
    }
  }

  export interface TraceContextMap {
    NONE: 0;
    TRACE_CONTEXT: 1;
    GRPC_TRACE_BIN: 2;
    CLOUD_TRACE_CONTEXT: 3;
    B3: 4;
  }

  export const TraceContext: TraceContextMap;

  export enum OcagentAddressCase {
    OCAGENT_ADDRESS_NOT_SET = 0,
    HTTP_ADDRESS = 3,
    GRPC_ADDRESS = 4,
  }
}

export class TraceConfig extends jspb.Message {
  hasProbabilitySampler(): boolean;
  clearProbabilitySampler(): void;
  getProbabilitySampler(): ProbabilitySampler | undefined;
  setProbabilitySampler(value?: ProbabilitySampler): void;

  hasConstantSampler(): boolean;
  clearConstantSampler(): void;
  getConstantSampler(): ConstantSampler | undefined;
  setConstantSampler(value?: ConstantSampler): void;

  hasRateLimitingSampler(): boolean;
  clearRateLimitingSampler(): void;
  getRateLimitingSampler(): RateLimitingSampler | undefined;
  setRateLimitingSampler(value?: RateLimitingSampler): void;

  getMaxNumberOfAttributes(): number;
  setMaxNumberOfAttributes(value: number): void;

  getMaxNumberOfAnnotations(): number;
  setMaxNumberOfAnnotations(value: number): void;

  getMaxNumberOfMessageEvents(): number;
  setMaxNumberOfMessageEvents(value: number): void;

  getMaxNumberOfLinks(): number;
  setMaxNumberOfLinks(value: number): void;

  getSamplerCase(): TraceConfig.SamplerCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TraceConfig.AsObject;
  static toObject(includeInstance: boolean, msg: TraceConfig): TraceConfig.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: TraceConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TraceConfig;
  static deserializeBinaryFromReader(message: TraceConfig, reader: jspb.BinaryReader): TraceConfig;
}

export namespace TraceConfig {
  export type AsObject = {
    probabilitySampler?: ProbabilitySampler.AsObject,
    constantSampler?: ConstantSampler.AsObject,
    rateLimitingSampler?: RateLimitingSampler.AsObject,
    maxNumberOfAttributes: number,
    maxNumberOfAnnotations: number,
    maxNumberOfMessageEvents: number,
    maxNumberOfLinks: number,
  }

  export enum SamplerCase {
    SAMPLER_NOT_SET = 0,
    PROBABILITY_SAMPLER = 1,
    CONSTANT_SAMPLER = 2,
    RATE_LIMITING_SAMPLER = 3,
  }
}

export class ProbabilitySampler extends jspb.Message {
  getSamplingprobability(): number;
  setSamplingprobability(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProbabilitySampler.AsObject;
  static toObject(includeInstance: boolean, msg: ProbabilitySampler): ProbabilitySampler.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ProbabilitySampler, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProbabilitySampler;
  static deserializeBinaryFromReader(message: ProbabilitySampler, reader: jspb.BinaryReader): ProbabilitySampler;
}

export namespace ProbabilitySampler {
  export type AsObject = {
    samplingprobability: number,
  }
}

export class ConstantSampler extends jspb.Message {
  getDecision(): ConstantSampler.ConstantDecisionMap[keyof ConstantSampler.ConstantDecisionMap];
  setDecision(value: ConstantSampler.ConstantDecisionMap[keyof ConstantSampler.ConstantDecisionMap]): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ConstantSampler.AsObject;
  static toObject(includeInstance: boolean, msg: ConstantSampler): ConstantSampler.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ConstantSampler, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ConstantSampler;
  static deserializeBinaryFromReader(message: ConstantSampler, reader: jspb.BinaryReader): ConstantSampler;
}

export namespace ConstantSampler {
  export type AsObject = {
    decision: ConstantSampler.ConstantDecisionMap[keyof ConstantSampler.ConstantDecisionMap],
  }

  export interface ConstantDecisionMap {
    ALWAYS_OFF: 0;
    ALWAYS_ON: 1;
    ALWAYS_PARENT: 2;
  }

  export const ConstantDecision: ConstantDecisionMap;
}

export class RateLimitingSampler extends jspb.Message {
  getQps(): number;
  setQps(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RateLimitingSampler.AsObject;
  static toObject(includeInstance: boolean, msg: RateLimitingSampler): RateLimitingSampler.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RateLimitingSampler, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RateLimitingSampler;
  static deserializeBinaryFromReader(message: RateLimitingSampler, reader: jspb.BinaryReader): RateLimitingSampler;
}

export namespace RateLimitingSampler {
  export type AsObject = {
    qps: number,
  }
}
