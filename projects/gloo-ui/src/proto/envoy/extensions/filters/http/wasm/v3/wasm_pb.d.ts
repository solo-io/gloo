/* eslint-disable */
// package: envoy.extensions.filters.http.wasm.v3
// file: envoy/extensions/filters/http/wasm/v3/wasm.proto

import * as jspb from "google-protobuf";
import * as envoy_extensions_wasm_v3_wasm_pb from "../../../../../../envoy/extensions/wasm/v3/wasm_pb";
import * as validate_validate_pb from "../../../../../../validate/validate_pb";
import * as gogoproto_gogo_pb from "../../../../../../gogoproto/gogo_pb";
import * as extproto_ext_pb from "../../../../../../protoc-gen-ext/extproto/ext_pb";

export class Wasm extends jspb.Message {
  hasConfig(): boolean;
  clearConfig(): void;
  getConfig(): envoy_extensions_wasm_v3_wasm_pb.PluginConfig | undefined;
  setConfig(value?: envoy_extensions_wasm_v3_wasm_pb.PluginConfig): void;

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
    config?: envoy_extensions_wasm_v3_wasm_pb.PluginConfig.AsObject,
  }
}
