// package: envoy.config.wasm.v2
// file: gloo/projects/gloo/api/external/envoy/api/v2/config/wasm/wasm.proto

import * as jspb from "google-protobuf";
import * as envoy_api_v2_core_base_pb from "../../../../../../../../../../envoy/api/v2/core/base_pb";
import * as validate_validate_pb from "../../../../../../../../../../validate/validate_pb";
import * as extproto_ext_pb from "../../../../../../../../../../extproto/ext_pb";

export class VmConfig extends jspb.Message {
  getVmId(): string;
  setVmId(value: string): void;

  getRuntime(): string;
  setRuntime(value: string): void;

  hasCode(): boolean;
  clearCode(): void;
  getCode(): envoy_api_v2_core_base_pb.AsyncDataSource | undefined;
  setCode(value?: envoy_api_v2_core_base_pb.AsyncDataSource): void;

  getConfiguration(): string;
  setConfiguration(value: string): void;

  getAllowPrecompiled(): boolean;
  setAllowPrecompiled(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VmConfig.AsObject;
  static toObject(includeInstance: boolean, msg: VmConfig): VmConfig.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: VmConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VmConfig;
  static deserializeBinaryFromReader(message: VmConfig, reader: jspb.BinaryReader): VmConfig;
}

export namespace VmConfig {
  export type AsObject = {
    vmId: string,
    runtime: string,
    code?: envoy_api_v2_core_base_pb.AsyncDataSource.AsObject,
    configuration: string,
    allowPrecompiled: boolean,
  }
}

export class PluginConfig extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getRootId(): string;
  setRootId(value: string): void;

  hasVmConfig(): boolean;
  clearVmConfig(): void;
  getVmConfig(): VmConfig | undefined;
  setVmConfig(value?: VmConfig): void;

  getConfiguration(): string;
  setConfiguration(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PluginConfig.AsObject;
  static toObject(includeInstance: boolean, msg: PluginConfig): PluginConfig.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: PluginConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PluginConfig;
  static deserializeBinaryFromReader(message: PluginConfig, reader: jspb.BinaryReader): PluginConfig;
}

export namespace PluginConfig {
  export type AsObject = {
    name: string,
    rootId: string,
    vmConfig?: VmConfig.AsObject,
    configuration: string,
  }
}

export class WasmService extends jspb.Message {
  hasConfig(): boolean;
  clearConfig(): void;
  getConfig(): PluginConfig | undefined;
  setConfig(value?: PluginConfig): void;

  getSingleton(): boolean;
  setSingleton(value: boolean): void;

  getStatPrefix(): string;
  setStatPrefix(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): WasmService.AsObject;
  static toObject(includeInstance: boolean, msg: WasmService): WasmService.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: WasmService, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): WasmService;
  static deserializeBinaryFromReader(message: WasmService, reader: jspb.BinaryReader): WasmService;
}

export namespace WasmService {
  export type AsObject = {
    config?: PluginConfig.AsObject,
    singleton: boolean,
    statPrefix: string,
  }
}

