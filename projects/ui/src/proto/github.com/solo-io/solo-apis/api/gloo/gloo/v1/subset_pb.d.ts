/* eslint-disable */
// package: gloo.solo.io
// file: github.com/solo-io/solo-apis/api/gloo/gloo/v1/subset.proto

import * as jspb from "google-protobuf";
import * as extproto_ext_pb from "../../../../../../../extproto/ext_pb";

export class Subset extends jspb.Message {
  getValuesMap(): jspb.Map<string, string>;
  clearValuesMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Subset.AsObject;
  static toObject(includeInstance: boolean, msg: Subset): Subset.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Subset, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Subset;
  static deserializeBinaryFromReader(message: Subset, reader: jspb.BinaryReader): Subset;
}

export namespace Subset {
  export type AsObject = {
    valuesMap: Array<[string, string]>,
  }
}
