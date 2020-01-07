// package: core.solo.io
// file: solo-kit/api/v1/ref.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../gogoproto/gogo_pb";
import * as extproto_ext_pb from "../../../extproto/ext_pb";

export class ResourceRef extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getNamespace(): string;
  setNamespace(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResourceRef.AsObject;
  static toObject(includeInstance: boolean, msg: ResourceRef): ResourceRef.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ResourceRef, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResourceRef;
  static deserializeBinaryFromReader(message: ResourceRef, reader: jspb.BinaryReader): ResourceRef;
}

export namespace ResourceRef {
  export type AsObject = {
    name: string,
    namespace: string,
  }
}

