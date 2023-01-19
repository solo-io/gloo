/* eslint-disable */
// package: gloo.solo.io
// file: github.com/solo-io/solo-apis/api/gloo/gloo/v1/options/router/router.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";

export class Router extends jspb.Message {
  hasSuppressEnvoyHeaders(): boolean;
  clearSuppressEnvoyHeaders(): void;
  getSuppressEnvoyHeaders(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setSuppressEnvoyHeaders(value?: google_protobuf_wrappers_pb.BoolValue): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Router.AsObject;
  static toObject(includeInstance: boolean, msg: Router): Router.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Router, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Router;
  static deserializeBinaryFromReader(message: Router, reader: jspb.BinaryReader): Router;
}

export namespace Router {
  export type AsObject = {
    suppressEnvoyHeaders?: google_protobuf_wrappers_pb.BoolValue.AsObject,
  }
}
