/* eslint-disable */
// package: gateway.solo.io
// file: github.com/solo-io/solo-apis/api/gloo/gateway/v1/external_options.proto

import * as jspb from "google-protobuf";
import * as extproto_ext_pb from "../../../../../../../extproto/ext_pb";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";
import * as github_com_solo_io_solo_kit_api_v1_metadata_pb from "../../../../../../../github.com/solo-io/solo-kit/api/v1/metadata_pb";
import * as github_com_solo_io_solo_kit_api_v1_status_pb from "../../../../../../../github.com/solo-io/solo-kit/api/v1/status_pb";
import * as github_com_solo_io_solo_kit_api_v1_ref_pb from "../../../../../../../github.com/solo-io/solo-kit/api/v1/ref_pb";
import * as github_com_solo_io_solo_kit_api_v1_solo_kit_pb from "../../../../../../../github.com/solo-io/solo-kit/api/v1/solo-kit_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/options_pb";

export class VirtualHostOption extends jspb.Message {
  hasStatus(): boolean;
  clearStatus(): void;
  getStatus(): github_com_solo_io_solo_kit_api_v1_status_pb.Status | undefined;
  setStatus(value?: github_com_solo_io_solo_kit_api_v1_status_pb.Status): void;

  hasMetadata(): boolean;
  clearMetadata(): void;
  getMetadata(): github_com_solo_io_solo_kit_api_v1_metadata_pb.Metadata | undefined;
  setMetadata(value?: github_com_solo_io_solo_kit_api_v1_metadata_pb.Metadata): void;

  hasOptions(): boolean;
  clearOptions(): void;
  getOptions(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.VirtualHostOptions | undefined;
  setOptions(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.VirtualHostOptions): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VirtualHostOption.AsObject;
  static toObject(includeInstance: boolean, msg: VirtualHostOption): VirtualHostOption.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: VirtualHostOption, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VirtualHostOption;
  static deserializeBinaryFromReader(message: VirtualHostOption, reader: jspb.BinaryReader): VirtualHostOption;
}

export namespace VirtualHostOption {
  export type AsObject = {
    status?: github_com_solo_io_solo_kit_api_v1_status_pb.Status.AsObject,
    metadata?: github_com_solo_io_solo_kit_api_v1_metadata_pb.Metadata.AsObject,
    options?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.VirtualHostOptions.AsObject,
  }
}

export class RouteOption extends jspb.Message {
  hasStatus(): boolean;
  clearStatus(): void;
  getStatus(): github_com_solo_io_solo_kit_api_v1_status_pb.Status | undefined;
  setStatus(value?: github_com_solo_io_solo_kit_api_v1_status_pb.Status): void;

  hasMetadata(): boolean;
  clearMetadata(): void;
  getMetadata(): github_com_solo_io_solo_kit_api_v1_metadata_pb.Metadata | undefined;
  setMetadata(value?: github_com_solo_io_solo_kit_api_v1_metadata_pb.Metadata): void;

  hasOptions(): boolean;
  clearOptions(): void;
  getOptions(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.RouteOptions | undefined;
  setOptions(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.RouteOptions): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RouteOption.AsObject;
  static toObject(includeInstance: boolean, msg: RouteOption): RouteOption.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RouteOption, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RouteOption;
  static deserializeBinaryFromReader(message: RouteOption, reader: jspb.BinaryReader): RouteOption;
}

export namespace RouteOption {
  export type AsObject = {
    status?: github_com_solo_io_solo_kit_api_v1_status_pb.Status.AsObject,
    metadata?: github_com_solo_io_solo_kit_api_v1_metadata_pb.Metadata.AsObject,
    options?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.RouteOptions.AsObject,
  }
}
