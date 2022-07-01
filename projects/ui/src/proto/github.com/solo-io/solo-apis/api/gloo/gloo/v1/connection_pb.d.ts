/* eslint-disable */
// package: gloo.solo.io
// file: github.com/solo-io/solo-apis/api/gloo/gloo/v1/connection.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_duration_pb from "google-protobuf/google/protobuf/duration_pb";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";
import * as extproto_ext_pb from "../../../../../../../extproto/ext_pb";

export class ConnectionConfig extends jspb.Message {
  getMaxRequestsPerConnection(): number;
  setMaxRequestsPerConnection(value: number): void;

  hasConnectTimeout(): boolean;
  clearConnectTimeout(): void;
  getConnectTimeout(): google_protobuf_duration_pb.Duration | undefined;
  setConnectTimeout(value?: google_protobuf_duration_pb.Duration): void;

  hasTcpKeepalive(): boolean;
  clearTcpKeepalive(): void;
  getTcpKeepalive(): ConnectionConfig.TcpKeepAlive | undefined;
  setTcpKeepalive(value?: ConnectionConfig.TcpKeepAlive): void;

  hasPerConnectionBufferLimitBytes(): boolean;
  clearPerConnectionBufferLimitBytes(): void;
  getPerConnectionBufferLimitBytes(): google_protobuf_wrappers_pb.UInt32Value | undefined;
  setPerConnectionBufferLimitBytes(value?: google_protobuf_wrappers_pb.UInt32Value): void;

  hasCommonHttpProtocolOptions(): boolean;
  clearCommonHttpProtocolOptions(): void;
  getCommonHttpProtocolOptions(): ConnectionConfig.HttpProtocolOptions | undefined;
  setCommonHttpProtocolOptions(value?: ConnectionConfig.HttpProtocolOptions): void;

  hasHttp1ProtocolOptions(): boolean;
  clearHttp1ProtocolOptions(): void;
  getHttp1ProtocolOptions(): ConnectionConfig.Http1ProtocolOptions | undefined;
  setHttp1ProtocolOptions(value?: ConnectionConfig.Http1ProtocolOptions): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ConnectionConfig.AsObject;
  static toObject(includeInstance: boolean, msg: ConnectionConfig): ConnectionConfig.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ConnectionConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ConnectionConfig;
  static deserializeBinaryFromReader(message: ConnectionConfig, reader: jspb.BinaryReader): ConnectionConfig;
}

export namespace ConnectionConfig {
  export type AsObject = {
    maxRequestsPerConnection: number,
    connectTimeout?: google_protobuf_duration_pb.Duration.AsObject,
    tcpKeepalive?: ConnectionConfig.TcpKeepAlive.AsObject,
    perConnectionBufferLimitBytes?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
    commonHttpProtocolOptions?: ConnectionConfig.HttpProtocolOptions.AsObject,
    http1ProtocolOptions?: ConnectionConfig.Http1ProtocolOptions.AsObject,
  }

  export class TcpKeepAlive extends jspb.Message {
    getKeepaliveProbes(): number;
    setKeepaliveProbes(value: number): void;

    hasKeepaliveTime(): boolean;
    clearKeepaliveTime(): void;
    getKeepaliveTime(): google_protobuf_duration_pb.Duration | undefined;
    setKeepaliveTime(value?: google_protobuf_duration_pb.Duration): void;

    hasKeepaliveInterval(): boolean;
    clearKeepaliveInterval(): void;
    getKeepaliveInterval(): google_protobuf_duration_pb.Duration | undefined;
    setKeepaliveInterval(value?: google_protobuf_duration_pb.Duration): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): TcpKeepAlive.AsObject;
    static toObject(includeInstance: boolean, msg: TcpKeepAlive): TcpKeepAlive.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: TcpKeepAlive, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): TcpKeepAlive;
    static deserializeBinaryFromReader(message: TcpKeepAlive, reader: jspb.BinaryReader): TcpKeepAlive;
  }

  export namespace TcpKeepAlive {
    export type AsObject = {
      keepaliveProbes: number,
      keepaliveTime?: google_protobuf_duration_pb.Duration.AsObject,
      keepaliveInterval?: google_protobuf_duration_pb.Duration.AsObject,
    }
  }

  export class HttpProtocolOptions extends jspb.Message {
    hasIdleTimeout(): boolean;
    clearIdleTimeout(): void;
    getIdleTimeout(): google_protobuf_duration_pb.Duration | undefined;
    setIdleTimeout(value?: google_protobuf_duration_pb.Duration): void;

    getMaxHeadersCount(): number;
    setMaxHeadersCount(value: number): void;

    hasMaxStreamDuration(): boolean;
    clearMaxStreamDuration(): void;
    getMaxStreamDuration(): google_protobuf_duration_pb.Duration | undefined;
    setMaxStreamDuration(value?: google_protobuf_duration_pb.Duration): void;

    getHeadersWithUnderscoresAction(): ConnectionConfig.HttpProtocolOptions.HeadersWithUnderscoresActionMap[keyof ConnectionConfig.HttpProtocolOptions.HeadersWithUnderscoresActionMap];
    setHeadersWithUnderscoresAction(value: ConnectionConfig.HttpProtocolOptions.HeadersWithUnderscoresActionMap[keyof ConnectionConfig.HttpProtocolOptions.HeadersWithUnderscoresActionMap]): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): HttpProtocolOptions.AsObject;
    static toObject(includeInstance: boolean, msg: HttpProtocolOptions): HttpProtocolOptions.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: HttpProtocolOptions, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): HttpProtocolOptions;
    static deserializeBinaryFromReader(message: HttpProtocolOptions, reader: jspb.BinaryReader): HttpProtocolOptions;
  }

  export namespace HttpProtocolOptions {
    export type AsObject = {
      idleTimeout?: google_protobuf_duration_pb.Duration.AsObject,
      maxHeadersCount: number,
      maxStreamDuration?: google_protobuf_duration_pb.Duration.AsObject,
      headersWithUnderscoresAction: ConnectionConfig.HttpProtocolOptions.HeadersWithUnderscoresActionMap[keyof ConnectionConfig.HttpProtocolOptions.HeadersWithUnderscoresActionMap],
    }

    export interface HeadersWithUnderscoresActionMap {
      ALLOW: 0;
      REJECT_REQUEST: 1;
      DROP_HEADER: 2;
    }

    export const HeadersWithUnderscoresAction: HeadersWithUnderscoresActionMap;
  }

  export class Http1ProtocolOptions extends jspb.Message {
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

    hasOverrideStreamErrorOnInvalidHttpMessage(): boolean;
    clearOverrideStreamErrorOnInvalidHttpMessage(): void;
    getOverrideStreamErrorOnInvalidHttpMessage(): google_protobuf_wrappers_pb.BoolValue | undefined;
    setOverrideStreamErrorOnInvalidHttpMessage(value?: google_protobuf_wrappers_pb.BoolValue): void;

    getHeaderFormatCase(): Http1ProtocolOptions.HeaderFormatCase;
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Http1ProtocolOptions.AsObject;
    static toObject(includeInstance: boolean, msg: Http1ProtocolOptions): Http1ProtocolOptions.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Http1ProtocolOptions, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Http1ProtocolOptions;
    static deserializeBinaryFromReader(message: Http1ProtocolOptions, reader: jspb.BinaryReader): Http1ProtocolOptions;
  }

  export namespace Http1ProtocolOptions {
    export type AsObject = {
      enableTrailers: boolean,
      properCaseHeaderKeyFormat: boolean,
      preserveCaseHeaderKeyFormat: boolean,
      overrideStreamErrorOnInvalidHttpMessage?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    }

    export enum HeaderFormatCase {
      HEADER_FORMAT_NOT_SET = 0,
      PROPER_CASE_HEADER_KEY_FORMAT = 22,
      PRESERVE_CASE_HEADER_KEY_FORMAT = 31,
    }
  }
}
