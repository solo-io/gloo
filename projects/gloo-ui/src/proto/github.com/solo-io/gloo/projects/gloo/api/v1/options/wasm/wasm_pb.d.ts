// package: wasm.options.gloo.solo.io
// file: github.com/solo-io/gloo/projects/gloo/api/v1/options/wasm/wasm.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../../../../../../../gogoproto/gogo_pb";

export class PluginSource extends jspb.Message {
  getImage(): string;
  setImage(value: string): void;

  getConfig(): string;
  setConfig(value: string): void;

  getName(): string;
  setName(value: string): void;

  getRootId(): string;
  setRootId(value: string): void;

  getVmType(): PluginSource.VmTypeMap[keyof PluginSource.VmTypeMap];
  setVmType(value: PluginSource.VmTypeMap[keyof PluginSource.VmTypeMap]): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PluginSource.AsObject;
  static toObject(includeInstance: boolean, msg: PluginSource): PluginSource.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: PluginSource, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PluginSource;
  static deserializeBinaryFromReader(message: PluginSource, reader: jspb.BinaryReader): PluginSource;
}

export namespace PluginSource {
  export type AsObject = {
    image: string,
    config: string,
    name: string,
    rootId: string,
    vmType: PluginSource.VmTypeMap[keyof PluginSource.VmTypeMap],
  }

  export interface VmTypeMap {
    V8: 0;
    WAVM: 1;
  }

  export const VmType: VmTypeMap;
}

