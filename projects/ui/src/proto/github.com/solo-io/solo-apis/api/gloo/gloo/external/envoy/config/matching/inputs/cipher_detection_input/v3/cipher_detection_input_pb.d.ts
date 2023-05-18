/* eslint-disable */
// package: envoy.config.matching.cipher_detection_input.v3
// file: github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/config/matching/inputs/cipher_detection_input/v3/cipher_detection_input.proto

import * as jspb from "google-protobuf";
import * as udpa_annotations_status_pb from "../../../../../../../../../../../../../udpa/annotations/status_pb";
import * as validate_validate_pb from "../../../../../../../../../../../../../validate/validate_pb";
import * as extproto_ext_pb from "../../../../../../../../../../../../../extproto/ext_pb";

export class CipherDetectionInput extends jspb.Message {
  clearPassthroughCiphersList(): void;
  getPassthroughCiphersList(): Array<number>;
  setPassthroughCiphersList(value: Array<number>): void;
  addPassthroughCiphers(value: number, index?: number): number;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CipherDetectionInput.AsObject;
  static toObject(includeInstance: boolean, msg: CipherDetectionInput): CipherDetectionInput.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CipherDetectionInput, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CipherDetectionInput;
  static deserializeBinaryFromReader(message: CipherDetectionInput, reader: jspb.BinaryReader): CipherDetectionInput;
}

export namespace CipherDetectionInput {
  export type AsObject = {
    passthroughCiphersList: Array<number>,
  }
}
