/* eslint-disable */
// package: hcm.options.gloo.solo.io
// file: github.com/solo-io/solo-apis/api/gloo/gloo/v1/options/hcm/hcm.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";
import * as google_protobuf_duration_pb from "google-protobuf/google/protobuf/duration_pb";
import * as google_protobuf_any_pb from "google-protobuf/google/protobuf/any_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_tracing_tracing_pb from "../../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/options/tracing/tracing_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_protocol_upgrade_protocol_upgrade_pb from "../../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/options/protocol_upgrade/protocol_upgrade_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_protocol_protocol_pb from "../../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/options/protocol/protocol_pb";
import * as extproto_ext_pb from "../../../../../../../../../extproto/ext_pb";

export class HttpConnectionManagerSettings extends jspb.Message {
  hasSkipXffAppend(): boolean;
  clearSkipXffAppend(): void;
  getSkipXffAppend(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setSkipXffAppend(value?: google_protobuf_wrappers_pb.BoolValue): void;

  hasVia(): boolean;
  clearVia(): void;
  getVia(): google_protobuf_wrappers_pb.StringValue | undefined;
  setVia(value?: google_protobuf_wrappers_pb.StringValue): void;

  hasXffNumTrustedHops(): boolean;
  clearXffNumTrustedHops(): void;
  getXffNumTrustedHops(): google_protobuf_wrappers_pb.UInt32Value | undefined;
  setXffNumTrustedHops(value?: google_protobuf_wrappers_pb.UInt32Value): void;

  hasUseRemoteAddress(): boolean;
  clearUseRemoteAddress(): void;
  getUseRemoteAddress(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setUseRemoteAddress(value?: google_protobuf_wrappers_pb.BoolValue): void;

  hasGenerateRequestId(): boolean;
  clearGenerateRequestId(): void;
  getGenerateRequestId(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setGenerateRequestId(value?: google_protobuf_wrappers_pb.BoolValue): void;

  hasProxy100Continue(): boolean;
  clearProxy100Continue(): void;
  getProxy100Continue(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setProxy100Continue(value?: google_protobuf_wrappers_pb.BoolValue): void;

  hasStreamIdleTimeout(): boolean;
  clearStreamIdleTimeout(): void;
  getStreamIdleTimeout(): google_protobuf_duration_pb.Duration | undefined;
  setStreamIdleTimeout(value?: google_protobuf_duration_pb.Duration): void;

  hasIdleTimeout(): boolean;
  clearIdleTimeout(): void;
  getIdleTimeout(): google_protobuf_duration_pb.Duration | undefined;
  setIdleTimeout(value?: google_protobuf_duration_pb.Duration): void;

  hasMaxRequestHeadersKb(): boolean;
  clearMaxRequestHeadersKb(): void;
  getMaxRequestHeadersKb(): google_protobuf_wrappers_pb.UInt32Value | undefined;
  setMaxRequestHeadersKb(value?: google_protobuf_wrappers_pb.UInt32Value): void;

  hasRequestTimeout(): boolean;
  clearRequestTimeout(): void;
  getRequestTimeout(): google_protobuf_duration_pb.Duration | undefined;
  setRequestTimeout(value?: google_protobuf_duration_pb.Duration): void;

  hasRequestHeadersTimeout(): boolean;
  clearRequestHeadersTimeout(): void;
  getRequestHeadersTimeout(): google_protobuf_duration_pb.Duration | undefined;
  setRequestHeadersTimeout(value?: google_protobuf_duration_pb.Duration): void;

  hasDrainTimeout(): boolean;
  clearDrainTimeout(): void;
  getDrainTimeout(): google_protobuf_duration_pb.Duration | undefined;
  setDrainTimeout(value?: google_protobuf_duration_pb.Duration): void;

  hasDelayedCloseTimeout(): boolean;
  clearDelayedCloseTimeout(): void;
  getDelayedCloseTimeout(): google_protobuf_duration_pb.Duration | undefined;
  setDelayedCloseTimeout(value?: google_protobuf_duration_pb.Duration): void;

  hasServerName(): boolean;
  clearServerName(): void;
  getServerName(): google_protobuf_wrappers_pb.StringValue | undefined;
  setServerName(value?: google_protobuf_wrappers_pb.StringValue): void;

  hasStripAnyHostPort(): boolean;
  clearStripAnyHostPort(): void;
  getStripAnyHostPort(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setStripAnyHostPort(value?: google_protobuf_wrappers_pb.BoolValue): void;

  hasAcceptHttp10(): boolean;
  clearAcceptHttp10(): void;
  getAcceptHttp10(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setAcceptHttp10(value?: google_protobuf_wrappers_pb.BoolValue): void;

  hasDefaultHostForHttp10(): boolean;
  clearDefaultHostForHttp10(): void;
  getDefaultHostForHttp10(): google_protobuf_wrappers_pb.StringValue | undefined;
  setDefaultHostForHttp10(value?: google_protobuf_wrappers_pb.StringValue): void;

  hasAllowChunkedLength(): boolean;
  clearAllowChunkedLength(): void;
  getAllowChunkedLength(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setAllowChunkedLength(value?: google_protobuf_wrappers_pb.BoolValue): void;

  hasEnableTrailers(): boolean;
  clearEnableTrailers(): void;
  getEnableTrailers(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setEnableTrailers(value?: google_protobuf_wrappers_pb.BoolValue): void;

  hasProperCaseHeaderKeyFormat(): boolean;
  clearProperCaseHeaderKeyFormat(): void;
  getProperCaseHeaderKeyFormat(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setProperCaseHeaderKeyFormat(value?: google_protobuf_wrappers_pb.BoolValue): void;

  hasPreserveCaseHeaderKeyFormat(): boolean;
  clearPreserveCaseHeaderKeyFormat(): void;
  getPreserveCaseHeaderKeyFormat(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setPreserveCaseHeaderKeyFormat(value?: google_protobuf_wrappers_pb.BoolValue): void;

  hasTracing(): boolean;
  clearTracing(): void;
  getTracing(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_tracing_tracing_pb.ListenerTracingSettings | undefined;
  setTracing(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_tracing_tracing_pb.ListenerTracingSettings): void;

  getForwardClientCertDetails(): HttpConnectionManagerSettings.ForwardClientCertDetailsMap[keyof HttpConnectionManagerSettings.ForwardClientCertDetailsMap];
  setForwardClientCertDetails(value: HttpConnectionManagerSettings.ForwardClientCertDetailsMap[keyof HttpConnectionManagerSettings.ForwardClientCertDetailsMap]): void;

  hasSetCurrentClientCertDetails(): boolean;
  clearSetCurrentClientCertDetails(): void;
  getSetCurrentClientCertDetails(): HttpConnectionManagerSettings.SetCurrentClientCertDetails | undefined;
  setSetCurrentClientCertDetails(value?: HttpConnectionManagerSettings.SetCurrentClientCertDetails): void;

  hasPreserveExternalRequestId(): boolean;
  clearPreserveExternalRequestId(): void;
  getPreserveExternalRequestId(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setPreserveExternalRequestId(value?: google_protobuf_wrappers_pb.BoolValue): void;

  clearUpgradesList(): void;
  getUpgradesList(): Array<github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_protocol_upgrade_protocol_upgrade_pb.ProtocolUpgradeConfig>;
  setUpgradesList(value: Array<github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_protocol_upgrade_protocol_upgrade_pb.ProtocolUpgradeConfig>): void;
  addUpgrades(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_protocol_upgrade_protocol_upgrade_pb.ProtocolUpgradeConfig, index?: number): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_protocol_upgrade_protocol_upgrade_pb.ProtocolUpgradeConfig;

  hasMaxConnectionDuration(): boolean;
  clearMaxConnectionDuration(): void;
  getMaxConnectionDuration(): google_protobuf_duration_pb.Duration | undefined;
  setMaxConnectionDuration(value?: google_protobuf_duration_pb.Duration): void;

  hasMaxStreamDuration(): boolean;
  clearMaxStreamDuration(): void;
  getMaxStreamDuration(): google_protobuf_duration_pb.Duration | undefined;
  setMaxStreamDuration(value?: google_protobuf_duration_pb.Duration): void;

  hasMaxHeadersCount(): boolean;
  clearMaxHeadersCount(): void;
  getMaxHeadersCount(): google_protobuf_wrappers_pb.UInt32Value | undefined;
  setMaxHeadersCount(value?: google_protobuf_wrappers_pb.UInt32Value): void;

  getHeadersWithUnderscoresAction(): HttpConnectionManagerSettings.HeadersWithUnderscoreActionMap[keyof HttpConnectionManagerSettings.HeadersWithUnderscoreActionMap];
  setHeadersWithUnderscoresAction(value: HttpConnectionManagerSettings.HeadersWithUnderscoreActionMap[keyof HttpConnectionManagerSettings.HeadersWithUnderscoreActionMap]): void;

  hasMaxRequestsPerConnection(): boolean;
  clearMaxRequestsPerConnection(): void;
  getMaxRequestsPerConnection(): google_protobuf_wrappers_pb.UInt32Value | undefined;
  setMaxRequestsPerConnection(value?: google_protobuf_wrappers_pb.UInt32Value): void;

  getServerHeaderTransformation(): HttpConnectionManagerSettings.ServerHeaderTransformationMap[keyof HttpConnectionManagerSettings.ServerHeaderTransformationMap];
  setServerHeaderTransformation(value: HttpConnectionManagerSettings.ServerHeaderTransformationMap[keyof HttpConnectionManagerSettings.ServerHeaderTransformationMap]): void;

  getPathWithEscapedSlashesAction(): HttpConnectionManagerSettings.PathWithEscapedSlashesActionMap[keyof HttpConnectionManagerSettings.PathWithEscapedSlashesActionMap];
  setPathWithEscapedSlashesAction(value: HttpConnectionManagerSettings.PathWithEscapedSlashesActionMap[keyof HttpConnectionManagerSettings.PathWithEscapedSlashesActionMap]): void;

  getCodecType(): HttpConnectionManagerSettings.CodecTypeMap[keyof HttpConnectionManagerSettings.CodecTypeMap];
  setCodecType(value: HttpConnectionManagerSettings.CodecTypeMap[keyof HttpConnectionManagerSettings.CodecTypeMap]): void;

  hasMergeSlashes(): boolean;
  clearMergeSlashes(): void;
  getMergeSlashes(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setMergeSlashes(value?: google_protobuf_wrappers_pb.BoolValue): void;

  hasNormalizePath(): boolean;
  clearNormalizePath(): void;
  getNormalizePath(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setNormalizePath(value?: google_protobuf_wrappers_pb.BoolValue): void;

  hasUuidRequestIdConfig(): boolean;
  clearUuidRequestIdConfig(): void;
  getUuidRequestIdConfig(): HttpConnectionManagerSettings.UuidRequestIdConfigSettings | undefined;
  setUuidRequestIdConfig(value?: HttpConnectionManagerSettings.UuidRequestIdConfigSettings): void;

  hasHttp2ProtocolOptions(): boolean;
  clearHttp2ProtocolOptions(): void;
  getHttp2ProtocolOptions(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_protocol_protocol_pb.Http2ProtocolOptions | undefined;
  setHttp2ProtocolOptions(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_protocol_protocol_pb.Http2ProtocolOptions): void;

  hasInternalAddressConfig(): boolean;
  clearInternalAddressConfig(): void;
  getInternalAddressConfig(): HttpConnectionManagerSettings.InternalAddressConfig | undefined;
  setInternalAddressConfig(value?: HttpConnectionManagerSettings.InternalAddressConfig): void;

  getHeaderFormatCase(): HttpConnectionManagerSettings.HeaderFormatCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HttpConnectionManagerSettings.AsObject;
  static toObject(includeInstance: boolean, msg: HttpConnectionManagerSettings): HttpConnectionManagerSettings.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: HttpConnectionManagerSettings, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HttpConnectionManagerSettings;
  static deserializeBinaryFromReader(message: HttpConnectionManagerSettings, reader: jspb.BinaryReader): HttpConnectionManagerSettings;
}

export namespace HttpConnectionManagerSettings {
  export type AsObject = {
    skipXffAppend?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    via?: google_protobuf_wrappers_pb.StringValue.AsObject,
    xffNumTrustedHops?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
    useRemoteAddress?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    generateRequestId?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    proxy100Continue?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    streamIdleTimeout?: google_protobuf_duration_pb.Duration.AsObject,
    idleTimeout?: google_protobuf_duration_pb.Duration.AsObject,
    maxRequestHeadersKb?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
    requestTimeout?: google_protobuf_duration_pb.Duration.AsObject,
    requestHeadersTimeout?: google_protobuf_duration_pb.Duration.AsObject,
    drainTimeout?: google_protobuf_duration_pb.Duration.AsObject,
    delayedCloseTimeout?: google_protobuf_duration_pb.Duration.AsObject,
    serverName?: google_protobuf_wrappers_pb.StringValue.AsObject,
    stripAnyHostPort?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    acceptHttp10?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    defaultHostForHttp10?: google_protobuf_wrappers_pb.StringValue.AsObject,
    allowChunkedLength?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    enableTrailers?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    properCaseHeaderKeyFormat?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    preserveCaseHeaderKeyFormat?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    tracing?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_tracing_tracing_pb.ListenerTracingSettings.AsObject,
    forwardClientCertDetails: HttpConnectionManagerSettings.ForwardClientCertDetailsMap[keyof HttpConnectionManagerSettings.ForwardClientCertDetailsMap],
    setCurrentClientCertDetails?: HttpConnectionManagerSettings.SetCurrentClientCertDetails.AsObject,
    preserveExternalRequestId?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    upgradesList: Array<github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_protocol_upgrade_protocol_upgrade_pb.ProtocolUpgradeConfig.AsObject>,
    maxConnectionDuration?: google_protobuf_duration_pb.Duration.AsObject,
    maxStreamDuration?: google_protobuf_duration_pb.Duration.AsObject,
    maxHeadersCount?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
    headersWithUnderscoresAction: HttpConnectionManagerSettings.HeadersWithUnderscoreActionMap[keyof HttpConnectionManagerSettings.HeadersWithUnderscoreActionMap],
    maxRequestsPerConnection?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
    serverHeaderTransformation: HttpConnectionManagerSettings.ServerHeaderTransformationMap[keyof HttpConnectionManagerSettings.ServerHeaderTransformationMap],
    pathWithEscapedSlashesAction: HttpConnectionManagerSettings.PathWithEscapedSlashesActionMap[keyof HttpConnectionManagerSettings.PathWithEscapedSlashesActionMap],
    codecType: HttpConnectionManagerSettings.CodecTypeMap[keyof HttpConnectionManagerSettings.CodecTypeMap],
    mergeSlashes?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    normalizePath?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    uuidRequestIdConfig?: HttpConnectionManagerSettings.UuidRequestIdConfigSettings.AsObject,
    http2ProtocolOptions?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_protocol_protocol_pb.Http2ProtocolOptions.AsObject,
    internalAddressConfig?: HttpConnectionManagerSettings.InternalAddressConfig.AsObject,
  }

  export class SetCurrentClientCertDetails extends jspb.Message {
    hasSubject(): boolean;
    clearSubject(): void;
    getSubject(): google_protobuf_wrappers_pb.BoolValue | undefined;
    setSubject(value?: google_protobuf_wrappers_pb.BoolValue): void;

    hasCert(): boolean;
    clearCert(): void;
    getCert(): google_protobuf_wrappers_pb.BoolValue | undefined;
    setCert(value?: google_protobuf_wrappers_pb.BoolValue): void;

    hasChain(): boolean;
    clearChain(): void;
    getChain(): google_protobuf_wrappers_pb.BoolValue | undefined;
    setChain(value?: google_protobuf_wrappers_pb.BoolValue): void;

    hasDns(): boolean;
    clearDns(): void;
    getDns(): google_protobuf_wrappers_pb.BoolValue | undefined;
    setDns(value?: google_protobuf_wrappers_pb.BoolValue): void;

    hasUri(): boolean;
    clearUri(): void;
    getUri(): google_protobuf_wrappers_pb.BoolValue | undefined;
    setUri(value?: google_protobuf_wrappers_pb.BoolValue): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): SetCurrentClientCertDetails.AsObject;
    static toObject(includeInstance: boolean, msg: SetCurrentClientCertDetails): SetCurrentClientCertDetails.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: SetCurrentClientCertDetails, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): SetCurrentClientCertDetails;
    static deserializeBinaryFromReader(message: SetCurrentClientCertDetails, reader: jspb.BinaryReader): SetCurrentClientCertDetails;
  }

  export namespace SetCurrentClientCertDetails {
    export type AsObject = {
      subject?: google_protobuf_wrappers_pb.BoolValue.AsObject,
      cert?: google_protobuf_wrappers_pb.BoolValue.AsObject,
      chain?: google_protobuf_wrappers_pb.BoolValue.AsObject,
      dns?: google_protobuf_wrappers_pb.BoolValue.AsObject,
      uri?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    }
  }

  export class UuidRequestIdConfigSettings extends jspb.Message {
    hasPackTraceReason(): boolean;
    clearPackTraceReason(): void;
    getPackTraceReason(): google_protobuf_wrappers_pb.BoolValue | undefined;
    setPackTraceReason(value?: google_protobuf_wrappers_pb.BoolValue): void;

    hasUseRequestIdForTraceSampling(): boolean;
    clearUseRequestIdForTraceSampling(): void;
    getUseRequestIdForTraceSampling(): google_protobuf_wrappers_pb.BoolValue | undefined;
    setUseRequestIdForTraceSampling(value?: google_protobuf_wrappers_pb.BoolValue): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): UuidRequestIdConfigSettings.AsObject;
    static toObject(includeInstance: boolean, msg: UuidRequestIdConfigSettings): UuidRequestIdConfigSettings.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: UuidRequestIdConfigSettings, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): UuidRequestIdConfigSettings;
    static deserializeBinaryFromReader(message: UuidRequestIdConfigSettings, reader: jspb.BinaryReader): UuidRequestIdConfigSettings;
  }

  export namespace UuidRequestIdConfigSettings {
    export type AsObject = {
      packTraceReason?: google_protobuf_wrappers_pb.BoolValue.AsObject,
      useRequestIdForTraceSampling?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    }
  }

  export class CidrRange extends jspb.Message {
    getAddressPrefix(): string;
    setAddressPrefix(value: string): void;

    hasPrefixLen(): boolean;
    clearPrefixLen(): void;
    getPrefixLen(): google_protobuf_wrappers_pb.UInt32Value | undefined;
    setPrefixLen(value?: google_protobuf_wrappers_pb.UInt32Value): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): CidrRange.AsObject;
    static toObject(includeInstance: boolean, msg: CidrRange): CidrRange.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: CidrRange, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): CidrRange;
    static deserializeBinaryFromReader(message: CidrRange, reader: jspb.BinaryReader): CidrRange;
  }

  export namespace CidrRange {
    export type AsObject = {
      addressPrefix: string,
      prefixLen?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
    }
  }

  export class InternalAddressConfig extends jspb.Message {
    hasUnixSockets(): boolean;
    clearUnixSockets(): void;
    getUnixSockets(): google_protobuf_wrappers_pb.BoolValue | undefined;
    setUnixSockets(value?: google_protobuf_wrappers_pb.BoolValue): void;

    clearCidrRangesList(): void;
    getCidrRangesList(): Array<HttpConnectionManagerSettings.CidrRange>;
    setCidrRangesList(value: Array<HttpConnectionManagerSettings.CidrRange>): void;
    addCidrRanges(value?: HttpConnectionManagerSettings.CidrRange, index?: number): HttpConnectionManagerSettings.CidrRange;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): InternalAddressConfig.AsObject;
    static toObject(includeInstance: boolean, msg: InternalAddressConfig): InternalAddressConfig.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: InternalAddressConfig, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): InternalAddressConfig;
    static deserializeBinaryFromReader(message: InternalAddressConfig, reader: jspb.BinaryReader): InternalAddressConfig;
  }

  export namespace InternalAddressConfig {
    export type AsObject = {
      unixSockets?: google_protobuf_wrappers_pb.BoolValue.AsObject,
      cidrRangesList: Array<HttpConnectionManagerSettings.CidrRange.AsObject>,
    }
  }

  export interface ForwardClientCertDetailsMap {
    SANITIZE: 0;
    FORWARD_ONLY: 1;
    APPEND_FORWARD: 2;
    SANITIZE_SET: 3;
    ALWAYS_FORWARD_ONLY: 4;
  }

  export const ForwardClientCertDetails: ForwardClientCertDetailsMap;

  export interface ServerHeaderTransformationMap {
    OVERWRITE: 0;
    APPEND_IF_ABSENT: 1;
    PASS_THROUGH: 2;
  }

  export const ServerHeaderTransformation: ServerHeaderTransformationMap;

  export interface HeadersWithUnderscoreActionMap {
    ALLOW: 0;
    REJECT_CLIENT_REQUEST: 1;
    DROP_HEADER: 2;
  }

  export const HeadersWithUnderscoreAction: HeadersWithUnderscoreActionMap;

  export interface PathWithEscapedSlashesActionMap {
    IMPLEMENTATION_SPECIFIC_DEFAULT: 0;
    KEEP_UNCHANGED: 1;
    REJECT_REQUEST: 2;
    UNESCAPE_AND_REDIRECT: 3;
    UNESCAPE_AND_FORWARD: 4;
  }

  export const PathWithEscapedSlashesAction: PathWithEscapedSlashesActionMap;

  export interface CodecTypeMap {
    AUTO: 0;
    HTTP1: 1;
    HTTP2: 2;
  }

  export const CodecType: CodecTypeMap;

  export enum HeaderFormatCase {
    HEADER_FORMAT_NOT_SET = 0,
    PROPER_CASE_HEADER_KEY_FORMAT = 22,
    PRESERVE_CASE_HEADER_KEY_FORMAT = 31,
  }
}
