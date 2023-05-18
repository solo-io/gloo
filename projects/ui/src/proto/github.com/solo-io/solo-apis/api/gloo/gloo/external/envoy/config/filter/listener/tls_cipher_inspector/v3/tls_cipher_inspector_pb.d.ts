/* eslint-disable */
// package: envoy.config.filter.listener.tls_cipher_inspector.v3
// file: github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/config/filter/listener/tls_cipher_inspector/v3/tls_cipher_inspector.proto

import * as jspb from "google-protobuf";
import * as udpa_annotations_status_pb from "../../../../../../../../../../../../../udpa/annotations/status_pb";
import * as extproto_ext_pb from "../../../../../../../../../../../../../extproto/ext_pb";

export class TlsCipherInspector extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TlsCipherInspector.AsObject;
  static toObject(includeInstance: boolean, msg: TlsCipherInspector): TlsCipherInspector.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: TlsCipherInspector, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TlsCipherInspector;
  static deserializeBinaryFromReader(message: TlsCipherInspector, reader: jspb.BinaryReader): TlsCipherInspector;
}

export namespace TlsCipherInspector {
  export type AsObject = {
  }
}
