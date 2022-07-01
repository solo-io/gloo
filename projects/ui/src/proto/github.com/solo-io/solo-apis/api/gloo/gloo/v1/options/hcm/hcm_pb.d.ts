/* eslint-disable */
// package: hcm.options.gloo.solo.io
// file: github.com/solo-io/solo-apis/api/gloo/gloo/v1/options/hcm/hcm.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";
import * as google_protobuf_duration_pb from "google-protobuf/google/protobuf/duration_pb";
import * as google_protobuf_any_pb from "google-protobuf/google/protobuf/any_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_tracing_tracing_pb from "../../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/options/tracing/tracing_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_protocol_upgrade_protocol_upgrade_pb from "../../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/options/protocol_upgrade/protocol_upgrade_pb";
import * as extproto_ext_pb from "../../../../../../../../../extproto/ext_pb";

export class HttpConnectionManagerSettings extends jspb.Message {
  getSkipXffAppend(): boolean;
  setSkipXffAppend(value: boolean): void;

  getVia(): string;
  setVia(value: string): void;

  getXffNumTrustedHops(): number;
  setXffNumTrustedHops(value: number): void;

  hasUseRemoteAddress(): boolean;
  clearUseRemoteAddress(): void;
  getUseRemoteAddress(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setUseRemoteAddress(value?: google_protobuf_wrappers_pb.BoolValue): void;

  hasGenerateRequestId(): boolean;
  clearGenerateRequestId(): void;
  getGenerateRequestId(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setGenerateRequestId(value?: google_protobuf_wrappers_pb.BoolValue): void;

  getProxy100Continue(): boolean;
  setProxy100Continue(value: boolean): void;

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

  getServerName(): string;
  setServerName(value: string): void;

  getStripAnyHostPort(): boolean;
  setStripAnyHostPort(value: boolean): void;

  getAcceptHttp10(): boolean;
  setAcceptHttp10(value: boolean): void;

  getDefaultHostForHttp10(): string;
  setDefaultHostForHttp10(value: string): void;

  getAllowChunkedLength(): boolean;
  setAllowChunkedLength(value: boolean): void;

  getEnableTrailers(): boolean;
  setEnableTrailers(value: boolean): void;

  hasProperCaseHeaderKeyFormat(): boolean;
  clearProperCaseHeaderKeyFormat(): void;
  getProperCaseHeaderKeyFormat(): boolean;
  setProperCaseHeaderKeyFormat(value: boolean): void;

  hasPreserveCaseHeaderKeyFormat(): boolean;
  clearPreserveCaseHeaderKeyFormat(): void;
  getPreserveCaseHeaderKeyFormat(): boolean;
  setPreserveCaseHeaderKeyFormat(value: boolean): void;

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

  getPreserveExternalRequestId(): boolean;
  setPreserveExternalRequestId(value: boolean): void;

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

  getMergeSlashes(): boolean;
  setMergeSlashes(value: boolean): void;

  hasNormalizePath(): boolean;
  clearNormalizePath(): void;
  getNormalizePath(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setNormalizePath(value?: google_protobuf_wrappers_pb.BoolValue): void;

  hasUuidRequestIdConfig(): boolean;
  clearUuidRequestIdConfig(): void;
  getUuidRequestIdConfig(): HttpConnectionManagerSettings.UuidRequestIdConfigSettings | undefined;
  setUuidRequestIdConfig(value?: HttpConnectionManagerSettings.UuidRequestIdConfigSettings): void;

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
    skipXffAppend: boolean,
    via: string,
    xffNumTrustedHops: number,
    useRemoteAddress?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    generateRequestId?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    proxy100Continue: boolean,
    streamIdleTimeout?: google_protobuf_duration_pb.Duration.AsObject,
    idleTimeout?: google_protobuf_duration_pb.Duration.AsObject,
    maxRequestHeadersKb?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
    requestTimeout?: google_protobuf_duration_pb.Duration.AsObject,
    requestHeadersTimeout?: google_protobuf_duration_pb.Duration.AsObject,
    drainTimeout?: google_protobuf_duration_pb.Duration.AsObject,
    delayedCloseTimeout?: google_protobuf_duration_pb.Duration.AsObject,
    serverName: string,
    stripAnyHostPort: boolean,
    acceptHttp10: boolean,
    defaultHostForHttp10: string,
    allowChunkedLength: boolean,
    enableTrailers: boolean,
    properCaseHeaderKeyFormat: boolean,
    preserveCaseHeaderKeyFormat: boolean,
    tracing?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_tracing_tracing_pb.ListenerTracingSettings.AsObject,
    forwardClientCertDetails: HttpConnectionManagerSettings.ForwardClientCertDetailsMap[keyof HttpConnectionManagerSettings.ForwardClientCertDetailsMap],
    setCurrentClientCertDetails?: HttpConnectionManagerSettings.SetCurrentClientCertDetails.AsObject,
    preserveExternalRequestId: boolean,
    upgradesList: Array<github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_protocol_upgrade_protocol_upgrade_pb.ProtocolUpgradeConfig.AsObject>,
    maxConnectionDuration?: google_protobuf_duration_pb.Duration.AsObject,
    maxStreamDuration?: google_protobuf_duration_pb.Duration.AsObject,
    maxHeadersCount?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
    headersWithUnderscoresAction: HttpConnectionManagerSettings.HeadersWithUnderscoreActionMap[keyof HttpConnectionManagerSettings.HeadersWithUnderscoreActionMap],
    maxRequestsPerConnection?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
    serverHeaderTransformation: HttpConnectionManagerSettings.ServerHeaderTransformationMap[keyof HttpConnectionManagerSettings.ServerHeaderTransformationMap],
    pathWithEscapedSlashesAction: HttpConnectionManagerSettings.PathWithEscapedSlashesActionMap[keyof HttpConnectionManagerSettings.PathWithEscapedSlashesActionMap],
    codecType: HttpConnectionManagerSettings.CodecTypeMap[keyof HttpConnectionManagerSettings.CodecTypeMap],
    mergeSlashes: boolean,
    normalizePath?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    uuidRequestIdConfig?: HttpConnectionManagerSettings.UuidRequestIdConfigSettings.AsObject,
  }

  export class SetCurrentClientCertDetails extends jspb.Message {
    hasSubject(): boolean;
    clearSubject(): void;
    getSubject(): google_protobuf_wrappers_pb.BoolValue | undefined;
    setSubject(value?: google_protobuf_wrappers_pb.BoolValue): void;

    getCert(): boolean;
    setCert(value: boolean): void;

    getChain(): boolean;
    setChain(value: boolean): void;

    getDns(): boolean;
    setDns(value: boolean): void;

    getUri(): boolean;
    setUri(value: boolean): void;

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
      cert: boolean,
      chain: boolean,
      dns: boolean,
      uri: boolean,
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
