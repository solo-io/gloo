/* eslint-disable */
// package: envoy.extensions.wasm.v3
// file: envoy/extensions/wasm/v3/wasm.proto

import * as jspb from "google-protobuf";
import * as envoy_config_core_v3_base_pb from "../../../../envoy/config/core/v3/base_pb";
import * as google_protobuf_any_pb from "google-protobuf/google/protobuf/any_pb";
import * as validate_validate_pb from "../../../../validate/validate_pb";
import * as gogoproto_gogo_pb from "../../../../gogoproto/gogo_pb";

export class VmConfig extends jspb.Message {
  getVmId(): string;
  setVmId(value: string): void;

  getRuntime(): string;
  setRuntime(value: string): void;

  hasCode(): boolean;
  clearCode(): void;
  getCode(): envoy_config_core_v3_base_pb.AsyncDataSource | undefined;
  setCode(value?: envoy_config_core_v3_base_pb.AsyncDataSource): void;

  hasConfiguration(): boolean;
  clearConfiguration(): void;
  getConfiguration(): google_protobuf_any_pb.Any | undefined;
  setConfiguration(value?: google_protobuf_any_pb.Any): void;

  getAllowPrecompiled(): boolean;
  setAllowPrecompiled(value: boolean): void;

  getNackOnCodeCacheMiss(): boolean;
  setNackOnCodeCacheMiss(value: boolean): void;

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
    code?: envoy_config_core_v3_base_pb.AsyncDataSource.AsObject,
    configuration?: google_protobuf_any_pb.Any.AsObject,
    allowPrecompiled: boolean,
    nackOnCodeCacheMiss: boolean,
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

  hasConfiguration(): boolean;
  clearConfiguration(): void;
  getConfiguration(): google_protobuf_any_pb.Any | undefined;
  setConfiguration(value?: google_protobuf_any_pb.Any): void;

  getFailOpen(): boolean;
  setFailOpen(value: boolean): void;

  getVmCase(): PluginConfig.VmCase;
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
    configuration?: google_protobuf_any_pb.Any.AsObject,
    failOpen: boolean,
  }

  export enum VmCase {
    VM_NOT_SET = 0,
    VM_CONFIG = 3,
  }
}

export class WasmService extends jspb.Message {
  hasConfig(): boolean;
  clearConfig(): void;
  getConfig(): PluginConfig | undefined;
  setConfig(value?: PluginConfig): void;

  getSingleton(): boolean;
  setSingleton(value: boolean): void;

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
  }
}
