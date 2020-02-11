// package: envoy.config.filter.http.gzip.v2
// file: gloo/projects/gloo/api/external/envoy/config/filter/http/gzip/v2/gzip.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../../../../../../../../../gogoproto/gogo_pb";
import * as extproto_ext_pb from "../../../../../../../../../../../extproto/ext_pb";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";
import * as validate_validate_pb from "../../../../../../../../../../../validate/validate_pb";

export class Gzip extends jspb.Message {
  hasMemoryLevel(): boolean;
  clearMemoryLevel(): void;
  getMemoryLevel(): google_protobuf_wrappers_pb.UInt32Value | undefined;
  setMemoryLevel(value?: google_protobuf_wrappers_pb.UInt32Value): void;

  hasContentLength(): boolean;
  clearContentLength(): void;
  getContentLength(): google_protobuf_wrappers_pb.UInt32Value | undefined;
  setContentLength(value?: google_protobuf_wrappers_pb.UInt32Value): void;

  getCompressionLevel(): Gzip.CompressionLevel.EnumMap[keyof Gzip.CompressionLevel.EnumMap];
  setCompressionLevel(value: Gzip.CompressionLevel.EnumMap[keyof Gzip.CompressionLevel.EnumMap]): void;

  getCompressionStrategy(): Gzip.CompressionStrategyMap[keyof Gzip.CompressionStrategyMap];
  setCompressionStrategy(value: Gzip.CompressionStrategyMap[keyof Gzip.CompressionStrategyMap]): void;

  clearContentTypeList(): void;
  getContentTypeList(): Array<string>;
  setContentTypeList(value: Array<string>): void;
  addContentType(value: string, index?: number): string;

  getDisableOnEtagHeader(): boolean;
  setDisableOnEtagHeader(value: boolean): void;

  getRemoveAcceptEncodingHeader(): boolean;
  setRemoveAcceptEncodingHeader(value: boolean): void;

  hasWindowBits(): boolean;
  clearWindowBits(): void;
  getWindowBits(): google_protobuf_wrappers_pb.UInt32Value | undefined;
  setWindowBits(value?: google_protobuf_wrappers_pb.UInt32Value): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Gzip.AsObject;
  static toObject(includeInstance: boolean, msg: Gzip): Gzip.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Gzip, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Gzip;
  static deserializeBinaryFromReader(message: Gzip, reader: jspb.BinaryReader): Gzip;
}

export namespace Gzip {
  export type AsObject = {
    memoryLevel?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
    contentLength?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
    compressionLevel: Gzip.CompressionLevel.EnumMap[keyof Gzip.CompressionLevel.EnumMap],
    compressionStrategy: Gzip.CompressionStrategyMap[keyof Gzip.CompressionStrategyMap],
    contentTypeList: Array<string>,
    disableOnEtagHeader: boolean,
    removeAcceptEncodingHeader: boolean,
    windowBits?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
  }

  export class CompressionLevel extends jspb.Message {
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): CompressionLevel.AsObject;
    static toObject(includeInstance: boolean, msg: CompressionLevel): CompressionLevel.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: CompressionLevel, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): CompressionLevel;
    static deserializeBinaryFromReader(message: CompressionLevel, reader: jspb.BinaryReader): CompressionLevel;
  }

  export namespace CompressionLevel {
    export type AsObject = {
    }

    export interface EnumMap {
      DEFAULT: 0;
      BEST: 1;
      SPEED: 2;
    }

    export const Enum: EnumMap;
  }

  export interface CompressionStrategyMap {
    DEFAULT: 0;
    FILTERED: 1;
    HUFFMAN: 2;
    RLE: 3;
  }

  export const CompressionStrategy: CompressionStrategyMap;
}

