/* eslint-disable */
// package: solo.io.envoy.extensions.filters.http.wasm.v3
// file: github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/extensions/filters/http/wasm/v3/wasm.proto

import * as jspb from "google-protobuf";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_extensions_wasm_v3_wasm_pb from "../../../../../../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/extensions/wasm/v3/wasm_pb";
import * as validate_validate_pb from "../../../../../../../../../../../../../validate/validate_pb";
import * as extproto_ext_pb from "../../../../../../../../../../../../../extproto/ext_pb";

export class Wasm extends jspb.Message {
  hasConfig(): boolean;
  clearConfig(): void;
  getConfig(): github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_extensions_wasm_v3_wasm_pb.PluginConfig | undefined;
  setConfig(value?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_extensions_wasm_v3_wasm_pb.PluginConfig): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Wasm.AsObject;
  static toObject(includeInstance: boolean, msg: Wasm): Wasm.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Wasm, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Wasm;
  static deserializeBinaryFromReader(message: Wasm, reader: jspb.BinaryReader): Wasm;
}

export namespace Wasm {
  export type AsObject = {
    config?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_extensions_wasm_v3_wasm_pb.PluginConfig.AsObject,
  }
}
