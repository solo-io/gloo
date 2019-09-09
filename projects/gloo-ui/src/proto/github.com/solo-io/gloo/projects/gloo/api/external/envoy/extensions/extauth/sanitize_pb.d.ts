// package: envoy.config.filter.http.sanitize.v2
// file: github.com/solo-io/gloo/projects/gloo/api/external/envoy/extensions/extauth/sanitize.proto

import * as jspb from "google-protobuf";

export class Sanitize extends jspb.Message {
  clearHeadersToRemoveList(): void;
  getHeadersToRemoveList(): Array<string>;
  setHeadersToRemoveList(value: Array<string>): void;
  addHeadersToRemove(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Sanitize.AsObject;
  static toObject(includeInstance: boolean, msg: Sanitize): Sanitize.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Sanitize, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Sanitize;
  static deserializeBinaryFromReader(message: Sanitize, reader: jspb.BinaryReader): Sanitize;
}

export namespace Sanitize {
  export type AsObject = {
    headersToRemoveList: Array<string>,
  }
}

