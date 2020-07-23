/* eslint-disable */
// package: envoy.type.v3
// file: envoy/type/v3/range.proto

import * as jspb from "google-protobuf";
import * as udpa_annotations_status_pb from "../../../udpa/annotations/status_pb";
import * as udpa_annotations_versioning_pb from "../../../udpa/annotations/versioning_pb";
import * as gogoproto_gogo_pb from "../../../gogoproto/gogo_pb";

export class Int64Range extends jspb.Message {
  getStart(): number;
  setStart(value: number): void;

  getEnd(): number;
  setEnd(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Int64Range.AsObject;
  static toObject(includeInstance: boolean, msg: Int64Range): Int64Range.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Int64Range, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Int64Range;
  static deserializeBinaryFromReader(message: Int64Range, reader: jspb.BinaryReader): Int64Range;
}

export namespace Int64Range {
  export type AsObject = {
    start: number,
    end: number,
  }
}

export class Int32Range extends jspb.Message {
  getStart(): number;
  setStart(value: number): void;

  getEnd(): number;
  setEnd(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Int32Range.AsObject;
  static toObject(includeInstance: boolean, msg: Int32Range): Int32Range.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Int32Range, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Int32Range;
  static deserializeBinaryFromReader(message: Int32Range, reader: jspb.BinaryReader): Int32Range;
}

export namespace Int32Range {
  export type AsObject = {
    start: number,
    end: number,
  }
}

export class DoubleRange extends jspb.Message {
  getStart(): number;
  setStart(value: number): void;

  getEnd(): number;
  setEnd(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DoubleRange.AsObject;
  static toObject(includeInstance: boolean, msg: DoubleRange): DoubleRange.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DoubleRange, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DoubleRange;
  static deserializeBinaryFromReader(message: DoubleRange, reader: jspb.BinaryReader): DoubleRange;
}

export namespace DoubleRange {
  export type AsObject = {
    start: number,
    end: number,
  }
}
