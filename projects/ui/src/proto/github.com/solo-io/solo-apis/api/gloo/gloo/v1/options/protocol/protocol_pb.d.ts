/* eslint-disable */
// package: protocol.options.gloo.solo.io
// file: github.com/solo-io/solo-apis/api/gloo/gloo/v1/options/protocol/protocol.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_duration_pb from "google-protobuf/google/protobuf/duration_pb";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";
import * as extproto_ext_pb from "../../../../../../../../../extproto/ext_pb";

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

  getHeadersWithUnderscoresAction(): HttpProtocolOptions.HeadersWithUnderscoresActionMap[keyof HttpProtocolOptions.HeadersWithUnderscoresActionMap];
  setHeadersWithUnderscoresAction(value: HttpProtocolOptions.HeadersWithUnderscoresActionMap[keyof HttpProtocolOptions.HeadersWithUnderscoresActionMap]): void;

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
    headersWithUnderscoresAction: HttpProtocolOptions.HeadersWithUnderscoresActionMap[keyof HttpProtocolOptions.HeadersWithUnderscoresActionMap],
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

export class Http2ProtocolOptions extends jspb.Message {
  hasMaxConcurrentStreams(): boolean;
  clearMaxConcurrentStreams(): void;
  getMaxConcurrentStreams(): google_protobuf_wrappers_pb.UInt32Value | undefined;
  setMaxConcurrentStreams(value?: google_protobuf_wrappers_pb.UInt32Value): void;

  hasInitialStreamWindowSize(): boolean;
  clearInitialStreamWindowSize(): void;
  getInitialStreamWindowSize(): google_protobuf_wrappers_pb.UInt32Value | undefined;
  setInitialStreamWindowSize(value?: google_protobuf_wrappers_pb.UInt32Value): void;

  hasInitialConnectionWindowSize(): boolean;
  clearInitialConnectionWindowSize(): void;
  getInitialConnectionWindowSize(): google_protobuf_wrappers_pb.UInt32Value | undefined;
  setInitialConnectionWindowSize(value?: google_protobuf_wrappers_pb.UInt32Value): void;

  hasOverrideStreamErrorOnInvalidHttpMessage(): boolean;
  clearOverrideStreamErrorOnInvalidHttpMessage(): void;
  getOverrideStreamErrorOnInvalidHttpMessage(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setOverrideStreamErrorOnInvalidHttpMessage(value?: google_protobuf_wrappers_pb.BoolValue): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Http2ProtocolOptions.AsObject;
  static toObject(includeInstance: boolean, msg: Http2ProtocolOptions): Http2ProtocolOptions.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Http2ProtocolOptions, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Http2ProtocolOptions;
  static deserializeBinaryFromReader(message: Http2ProtocolOptions, reader: jspb.BinaryReader): Http2ProtocolOptions;
}

export namespace Http2ProtocolOptions {
  export type AsObject = {
    maxConcurrentStreams?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
    initialStreamWindowSize?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
    initialConnectionWindowSize?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
    overrideStreamErrorOnInvalidHttpMessage?: google_protobuf_wrappers_pb.BoolValue.AsObject,
  }
}
