// package: grpc_web.plugins.gloo.solo.io
// file: github.com/solo-io/gloo/projects/gloo/api/v1/plugins/grpc_web/grpc_web.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../../../../../../../gogoproto/gogo_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_v1_plugins_transformation_parameters_pb from "../../../../../../../../../github.com/solo-io/gloo/projects/gloo/api/v1/plugins/transformation/parameters_pb";

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

