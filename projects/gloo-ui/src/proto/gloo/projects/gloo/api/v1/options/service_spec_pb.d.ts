// package: options.gloo.solo.io
// file: gloo/projects/gloo/api/v1/options/service_spec.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../../../../gogoproto/gogo_pb";
import * as gloo_projects_gloo_api_v1_options_rest_rest_pb from "../../../../../../gloo/projects/gloo/api/v1/options/rest/rest_pb";
import * as gloo_projects_gloo_api_v1_options_grpc_grpc_pb from "../../../../../../gloo/projects/gloo/api/v1/options/grpc/grpc_pb";
import * as extproto_ext_pb from "../../../../../../extproto/ext_pb";

export class ServiceSpec extends jspb.Message {
  hasRest(): boolean;
  clearRest(): void;
  getRest(): gloo_projects_gloo_api_v1_options_rest_rest_pb.ServiceSpec | undefined;
  setRest(value?: gloo_projects_gloo_api_v1_options_rest_rest_pb.ServiceSpec): void;

  hasGrpc(): boolean;
  clearGrpc(): void;
  getGrpc(): gloo_projects_gloo_api_v1_options_grpc_grpc_pb.ServiceSpec | undefined;
  setGrpc(value?: gloo_projects_gloo_api_v1_options_grpc_grpc_pb.ServiceSpec): void;

  getPluginTypeCase(): ServiceSpec.PluginTypeCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ServiceSpec.AsObject;
  static toObject(includeInstance: boolean, msg: ServiceSpec): ServiceSpec.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ServiceSpec, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ServiceSpec;
  static deserializeBinaryFromReader(message: ServiceSpec, reader: jspb.BinaryReader): ServiceSpec;
}

export namespace ServiceSpec {
  export type AsObject = {
    rest?: gloo_projects_gloo_api_v1_options_rest_rest_pb.ServiceSpec.AsObject,
    grpc?: gloo_projects_gloo_api_v1_options_grpc_grpc_pb.ServiceSpec.AsObject,
  }

  export enum PluginTypeCase {
    PLUGIN_TYPE_NOT_SET = 0,
    REST = 1,
    GRPC = 2,
  }
}

