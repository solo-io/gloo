/* eslint-disable */
// package: udpa.annotations
// file: udpa/annotations/versioning.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_descriptor_pb from "google-protobuf/google/protobuf/descriptor_pb";

export class VersioningAnnotation extends jspb.Message {
  getPreviousMessageType(): string;
  setPreviousMessageType(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VersioningAnnotation.AsObject;
  static toObject(includeInstance: boolean, msg: VersioningAnnotation): VersioningAnnotation.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: VersioningAnnotation, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VersioningAnnotation;
  static deserializeBinaryFromReader(message: VersioningAnnotation, reader: jspb.BinaryReader): VersioningAnnotation;
}

export namespace VersioningAnnotation {
  export type AsObject = {
    previousMessageType: string,
  }
}

  export const versioning: jspb.ExtensionFieldInfo<VersioningAnnotation>;
