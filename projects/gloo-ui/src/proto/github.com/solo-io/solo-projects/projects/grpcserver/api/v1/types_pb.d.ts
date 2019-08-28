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

export class Status extends jspb.Message {
  getCode(): Status.CodeMap[keyof Status.CodeMap];
  setCode(value: Status.CodeMap[keyof Status.CodeMap]): void;

  getMessage(): string;
  setMessage(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Status.AsObject;
  static toObject(includeInstance: boolean, msg: Status): Status.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Status, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Status;
  static deserializeBinaryFromReader(message: Status, reader: jspb.BinaryReader): Status;
}

export namespace Status {
  export type AsObject = {
    code: Status.CodeMap[keyof Status.CodeMap],
    message: string,
  }

  export interface CodeMap {
    ERROR: 0;
    WARNING: 1;
    OK: 2;
  }

  export const Code: CodeMap;
}

