// package: transformation.plugins.gloo.solo.io
// file: github.com/solo-io/gloo/projects/gloo/api/v1/plugins/transformation/prefix_rewrite.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../../../../../../../gogoproto/gogo_pb";

export class PrefixRewrite extends jspb.Message {
  getPrefixRewrite(): string;
  setPrefixRewrite(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PrefixRewrite.AsObject;
  static toObject(includeInstance: boolean, msg: PrefixRewrite): PrefixRewrite.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: PrefixRewrite, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PrefixRewrite;
  static deserializeBinaryFromReader(message: PrefixRewrite, reader: jspb.BinaryReader): PrefixRewrite;
}

export namespace PrefixRewrite {
  export type AsObject = {
    prefixRewrite: string,
  }
}

