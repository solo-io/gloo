// package: glooeeapi.solo.io
// file: github.com/solo-io/solo-projects/projects/grpcserver/api/v1/types.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../../../../../gogoproto/gogo_pb";

export class Raw extends jspb.Message {
  getFileName(): string;
  setFileName(value: string): void;

  getContent(): string;
  setContent(value: string): void;

  getContentRenderError(): string;
  setContentRenderError(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Raw.AsObject;
  static toObject(includeInstance: boolean, msg: Raw): Raw.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Raw, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Raw;
  static deserializeBinaryFromReader(message: Raw, reader: jspb.BinaryReader): Raw;
}

export namespace Raw {
  export type AsObject = {
    fileName: string,
    content: string,
    contentRenderError: string,
  }
}

