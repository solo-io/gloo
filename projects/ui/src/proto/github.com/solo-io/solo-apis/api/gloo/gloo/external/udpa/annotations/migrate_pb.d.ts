/* eslint-disable */
// package: solo.io.udpa.annotations
// file: github.com/solo-io/solo-apis/api/gloo/gloo/external/udpa/annotations/migrate.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_descriptor_pb from "google-protobuf/google/protobuf/descriptor_pb";

export class MigrateAnnotation extends jspb.Message {
  getRename(): string;
  setRename(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MigrateAnnotation.AsObject;
  static toObject(includeInstance: boolean, msg: MigrateAnnotation): MigrateAnnotation.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: MigrateAnnotation, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MigrateAnnotation;
  static deserializeBinaryFromReader(message: MigrateAnnotation, reader: jspb.BinaryReader): MigrateAnnotation;
}

export namespace MigrateAnnotation {
  export type AsObject = {
    rename: string,
  }
}

export class FieldMigrateAnnotation extends jspb.Message {
  getRename(): string;
  setRename(value: string): void;

  getOneofPromotion(): string;
  setOneofPromotion(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FieldMigrateAnnotation.AsObject;
  static toObject(includeInstance: boolean, msg: FieldMigrateAnnotation): FieldMigrateAnnotation.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: FieldMigrateAnnotation, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FieldMigrateAnnotation;
  static deserializeBinaryFromReader(message: FieldMigrateAnnotation, reader: jspb.BinaryReader): FieldMigrateAnnotation;
}

export namespace FieldMigrateAnnotation {
  export type AsObject = {
    rename: string,
    oneofPromotion: string,
  }
}

export class FileMigrateAnnotation extends jspb.Message {
  getMoveToPackage(): string;
  setMoveToPackage(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FileMigrateAnnotation.AsObject;
  static toObject(includeInstance: boolean, msg: FileMigrateAnnotation): FileMigrateAnnotation.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: FileMigrateAnnotation, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FileMigrateAnnotation;
  static deserializeBinaryFromReader(message: FileMigrateAnnotation, reader: jspb.BinaryReader): FileMigrateAnnotation;
}

export namespace FileMigrateAnnotation {
  export type AsObject = {
    moveToPackage: string,
  }
}

  export const messageMigrate: jspb.ExtensionFieldInfo<MigrateAnnotation>;

  export const fieldMigrate: jspb.ExtensionFieldInfo<FieldMigrateAnnotation>;

  export const enumMigrate: jspb.ExtensionFieldInfo<MigrateAnnotation>;

  export const enumValueMigrate: jspb.ExtensionFieldInfo<MigrateAnnotation>;

  export const fileMigrate: jspb.ExtensionFieldInfo<FileMigrateAnnotation>;
