/* eslint-disable */
// package: envoy.type.v3
// file: envoy/type/v3/semantic_version.proto

import * as jspb from "google-protobuf";
import * as udpa_annotations_status_pb from "../../../udpa/annotations/status_pb";
import * as udpa_annotations_versioning_pb from "../../../udpa/annotations/versioning_pb";
import * as gogoproto_gogo_pb from "../../../gogoproto/gogo_pb";

export class SemanticVersion extends jspb.Message {
  getMajorNumber(): number;
  setMajorNumber(value: number): void;

  getMinorNumber(): number;
  setMinorNumber(value: number): void;

  getPatch(): number;
  setPatch(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SemanticVersion.AsObject;
  static toObject(includeInstance: boolean, msg: SemanticVersion): SemanticVersion.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: SemanticVersion, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SemanticVersion;
  static deserializeBinaryFromReader(message: SemanticVersion, reader: jspb.BinaryReader): SemanticVersion;
}

export namespace SemanticVersion {
  export type AsObject = {
    majorNumber: number,
    minorNumber: number,
    patch: number,
  }
}
