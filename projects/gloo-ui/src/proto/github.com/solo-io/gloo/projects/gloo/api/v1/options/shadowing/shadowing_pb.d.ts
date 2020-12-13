/* eslint-disable */
// package: shadowing.options.gloo.solo.io
// file: github.com/solo-io/gloo/projects/gloo/api/v1/options/shadowing/shadowing.proto

import * as jspb from "google-protobuf";
import * as github_com_solo_io_solo_kit_api_v1_ref_pb from "../../../../../../../../../github.com/solo-io/solo-kit/api/v1/ref_pb";
import * as extproto_ext_pb from "../../../../../../../../../protoc-gen-ext/extproto/ext_pb";

export class RouteShadowing extends jspb.Message {
  hasUpstream(): boolean;
  clearUpstream(): void;
  getUpstream(): github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef | undefined;
  setUpstream(value?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef): void;

  getPercentage(): number;
  setPercentage(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RouteShadowing.AsObject;
  static toObject(includeInstance: boolean, msg: RouteShadowing): RouteShadowing.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RouteShadowing, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RouteShadowing;
  static deserializeBinaryFromReader(message: RouteShadowing, reader: jspb.BinaryReader): RouteShadowing;
}

export namespace RouteShadowing {
  export type AsObject = {
    upstream?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject,
    percentage: number,
  }
}
