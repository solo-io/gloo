/* eslint-disable */
// package: grpc_web.options.gloo.solo.io
// file: github.com/solo-io/solo-apis/api/gloo/gloo/v1/options/grpc_web/grpc_web.proto

import * as jspb from "google-protobuf";
import * as extproto_ext_pb from "../../../../../../../../../extproto/ext_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_transformation_parameters_pb from "../../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/options/transformation/parameters_pb";

export class GrpcWeb extends jspb.Message {
  getDisable(): boolean;
  setDisable(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GrpcWeb.AsObject;
  static toObject(includeInstance: boolean, msg: GrpcWeb): GrpcWeb.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GrpcWeb, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GrpcWeb;
  static deserializeBinaryFromReader(message: GrpcWeb, reader: jspb.BinaryReader): GrpcWeb;
}

export namespace GrpcWeb {
  export type AsObject = {
    disable: boolean,
  }
}
