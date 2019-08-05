// package: hostrewrite.plugins.gloo.solo.io
// file: github.com/solo-io/gloo/projects/gloo/api/v1/plugins/hostrewrite/hostrewrite.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../../../../../../../gogoproto/gogo_pb";

export class HostRewrite extends jspb.Message {
  hasHostRewrite(): boolean;
  clearHostRewrite(): void;
  getHostRewrite(): string;
  setHostRewrite(value: string): void;

  getHostRewriteTypeCase(): HostRewrite.HostRewriteTypeCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HostRewrite.AsObject;
  static toObject(includeInstance: boolean, msg: HostRewrite): HostRewrite.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: HostRewrite, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HostRewrite;
  static deserializeBinaryFromReader(message: HostRewrite, reader: jspb.BinaryReader): HostRewrite;
}

export namespace HostRewrite {
  export type AsObject = {
    hostRewrite: string,
  }

  export enum HostRewriteTypeCase {
    HOST_REWRITE_TYPE_NOT_SET = 0,
    HOST_REWRITE = 1,
  }
}

