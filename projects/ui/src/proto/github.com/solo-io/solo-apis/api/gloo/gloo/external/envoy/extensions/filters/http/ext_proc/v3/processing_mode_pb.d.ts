/* eslint-disable */
// package: solo.io.envoy.extensions.filters.http.ext_proc.v3
// file: github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/extensions/filters/http/ext_proc/v3/processing_mode.proto

import * as jspb from "google-protobuf";
import * as udpa_annotations_status_pb from "../../../../../../../../../../../../../udpa/annotations/status_pb";
import * as validate_validate_pb from "../../../../../../../../../../../../../validate/validate_pb";
import * as extproto_ext_pb from "../../../../../../../../../../../../../extproto/ext_pb";

export class ProcessingMode extends jspb.Message {
  getRequestHeaderMode(): ProcessingMode.HeaderSendModeMap[keyof ProcessingMode.HeaderSendModeMap];
  setRequestHeaderMode(value: ProcessingMode.HeaderSendModeMap[keyof ProcessingMode.HeaderSendModeMap]): void;

  getResponseHeaderMode(): ProcessingMode.HeaderSendModeMap[keyof ProcessingMode.HeaderSendModeMap];
  setResponseHeaderMode(value: ProcessingMode.HeaderSendModeMap[keyof ProcessingMode.HeaderSendModeMap]): void;

  getRequestBodyMode(): ProcessingMode.BodySendModeMap[keyof ProcessingMode.BodySendModeMap];
  setRequestBodyMode(value: ProcessingMode.BodySendModeMap[keyof ProcessingMode.BodySendModeMap]): void;

  getResponseBodyMode(): ProcessingMode.BodySendModeMap[keyof ProcessingMode.BodySendModeMap];
  setResponseBodyMode(value: ProcessingMode.BodySendModeMap[keyof ProcessingMode.BodySendModeMap]): void;

  getRequestTrailerMode(): ProcessingMode.HeaderSendModeMap[keyof ProcessingMode.HeaderSendModeMap];
  setRequestTrailerMode(value: ProcessingMode.HeaderSendModeMap[keyof ProcessingMode.HeaderSendModeMap]): void;

  getResponseTrailerMode(): ProcessingMode.HeaderSendModeMap[keyof ProcessingMode.HeaderSendModeMap];
  setResponseTrailerMode(value: ProcessingMode.HeaderSendModeMap[keyof ProcessingMode.HeaderSendModeMap]): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProcessingMode.AsObject;
  static toObject(includeInstance: boolean, msg: ProcessingMode): ProcessingMode.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ProcessingMode, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProcessingMode;
  static deserializeBinaryFromReader(message: ProcessingMode, reader: jspb.BinaryReader): ProcessingMode;
}

export namespace ProcessingMode {
  export type AsObject = {
    requestHeaderMode: ProcessingMode.HeaderSendModeMap[keyof ProcessingMode.HeaderSendModeMap],
    responseHeaderMode: ProcessingMode.HeaderSendModeMap[keyof ProcessingMode.HeaderSendModeMap],
    requestBodyMode: ProcessingMode.BodySendModeMap[keyof ProcessingMode.BodySendModeMap],
    responseBodyMode: ProcessingMode.BodySendModeMap[keyof ProcessingMode.BodySendModeMap],
    requestTrailerMode: ProcessingMode.HeaderSendModeMap[keyof ProcessingMode.HeaderSendModeMap],
    responseTrailerMode: ProcessingMode.HeaderSendModeMap[keyof ProcessingMode.HeaderSendModeMap],
  }

  export interface HeaderSendModeMap {
    DEFAULT: 0;
    SEND: 1;
    SKIP: 2;
  }

  export const HeaderSendMode: HeaderSendModeMap;

  export interface BodySendModeMap {
    NONE: 0;
    STREAMED: 1;
    BUFFERED: 2;
    BUFFERED_PARTIAL: 3;
  }

  export const BodySendMode: BodySendModeMap;
}
