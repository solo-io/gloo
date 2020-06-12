/* eslint-disable */
// package: envoy.extensions.filters.http.buffer.v3
// file: gloo/projects/gloo/api/external/envoy/extensions/filters/http/buffer/v3/buffer.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";
import * as validate_validate_pb from "../../../../../../../../../../../validate/validate_pb";
import * as gogoproto_gogo_pb from "../../../../../../../../../../../gogoproto/gogo_pb";
import * as extproto_ext_pb from "../../../../../../../../../../../protoc-gen-ext/extproto/ext_pb";

export class Buffer extends jspb.Message {
  hasMaxRequestBytes(): boolean;
  clearMaxRequestBytes(): void;
  getMaxRequestBytes(): google_protobuf_wrappers_pb.UInt32Value | undefined;
  setMaxRequestBytes(value?: google_protobuf_wrappers_pb.UInt32Value): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Buffer.AsObject;
  static toObject(includeInstance: boolean, msg: Buffer): Buffer.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Buffer, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Buffer;
  static deserializeBinaryFromReader(message: Buffer, reader: jspb.BinaryReader): Buffer;
}

export namespace Buffer {
  export type AsObject = {
    maxRequestBytes?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
  }
}

export class BufferPerRoute extends jspb.Message {
  hasDisabled(): boolean;
  clearDisabled(): void;
  getDisabled(): boolean;
  setDisabled(value: boolean): void;

  hasBuffer(): boolean;
  clearBuffer(): void;
  getBuffer(): Buffer | undefined;
  setBuffer(value?: Buffer): void;

  getOverrideCase(): BufferPerRoute.OverrideCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): BufferPerRoute.AsObject;
  static toObject(includeInstance: boolean, msg: BufferPerRoute): BufferPerRoute.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: BufferPerRoute, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): BufferPerRoute;
  static deserializeBinaryFromReader(message: BufferPerRoute, reader: jspb.BinaryReader): BufferPerRoute;
}

export namespace BufferPerRoute {
  export type AsObject = {
    disabled: boolean,
    buffer?: Buffer.AsObject,
  }

  export enum OverrideCase {
    OVERRIDE_NOT_SET = 0,
    DISABLED = 1,
    BUFFER = 2,
  }
}
