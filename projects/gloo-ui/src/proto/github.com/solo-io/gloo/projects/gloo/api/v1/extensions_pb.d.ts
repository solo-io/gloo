// package: gloo.solo.io
// file: github.com/solo-io/gloo/projects/gloo/api/v1/extensions.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_struct_pb from "google-protobuf/google/protobuf/struct_pb";
import * as gogoproto_gogo_pb from "../../../../../../../gogoproto/gogo_pb";

export class Extensions extends jspb.Message {
  getConfigsMap(): jspb.Map<string, google_protobuf_struct_pb.Struct>;
  clearConfigsMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Extensions.AsObject;
  static toObject(includeInstance: boolean, msg: Extensions): Extensions.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Extensions, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Extensions;
  static deserializeBinaryFromReader(message: Extensions, reader: jspb.BinaryReader): Extensions;
}

export namespace Extensions {
  export type AsObject = {
    configsMap: Array<[string, google_protobuf_struct_pb.Struct.AsObject]>,
  }
}

export class Extension extends jspb.Message {
  hasConfig(): boolean;
  clearConfig(): void;
  getConfig(): google_protobuf_struct_pb.Struct | undefined;
  setConfig(value?: google_protobuf_struct_pb.Struct): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Extension.AsObject;
  static toObject(includeInstance: boolean, msg: Extension): Extension.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Extension, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Extension;
  static deserializeBinaryFromReader(message: Extension, reader: jspb.BinaryReader): Extension;
}

export namespace Extension {
  export type AsObject = {
    config?: google_protobuf_struct_pb.Struct.AsObject,
  }
}

