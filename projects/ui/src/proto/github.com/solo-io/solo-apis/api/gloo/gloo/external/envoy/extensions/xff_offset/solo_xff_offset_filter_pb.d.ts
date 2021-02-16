/* eslint-disable */
// package: envoy.config.filter.http.solo_xff_offset.v2
// file: github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/extensions/xff_offset/solo_xff_offset_filter.proto

import * as jspb from "google-protobuf";

export class SoloXffOffset extends jspb.Message {
  getOffset(): number;
  setOffset(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SoloXffOffset.AsObject;
  static toObject(includeInstance: boolean, msg: SoloXffOffset): SoloXffOffset.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: SoloXffOffset, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SoloXffOffset;
  static deserializeBinaryFromReader(message: SoloXffOffset, reader: jspb.BinaryReader): SoloXffOffset;
}

export namespace SoloXffOffset {
  export type AsObject = {
    offset: number,
  }
}
