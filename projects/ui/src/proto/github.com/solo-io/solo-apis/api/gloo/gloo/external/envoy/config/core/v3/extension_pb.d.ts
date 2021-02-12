/* eslint-disable */
// package: solo.io.envoy.config.core.v3
// file: github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/config/core/v3/extension.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_any_pb from "google-protobuf/google/protobuf/any_pb";
import * as udpa_annotations_status_pb from "../../../../../../../../../../../udpa/annotations/status_pb";
import * as validate_validate_pb from "../../../../../../../../../../../validate/validate_pb";
import * as extproto_ext_pb from "../../../../../../../../../../../extproto/ext_pb";

export class TypedExtensionConfig extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  hasTypedConfig(): boolean;
  clearTypedConfig(): void;
  getTypedConfig(): google_protobuf_any_pb.Any | undefined;
  setTypedConfig(value?: google_protobuf_any_pb.Any): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TypedExtensionConfig.AsObject;
  static toObject(includeInstance: boolean, msg: TypedExtensionConfig): TypedExtensionConfig.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: TypedExtensionConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TypedExtensionConfig;
  static deserializeBinaryFromReader(message: TypedExtensionConfig, reader: jspb.BinaryReader): TypedExtensionConfig;
}

export namespace TypedExtensionConfig {
  export type AsObject = {
    name: string,
    typedConfig?: google_protobuf_any_pb.Any.AsObject,
  }
}
