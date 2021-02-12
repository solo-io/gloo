/* eslint-disable */
// package: wasm.options.gloo.solo.io
// file: github.com/solo-io/solo-apis/api/gloo/gloo/v1/options/wasm/wasm.proto

import * as jspb from "google-protobuf";
import * as extproto_ext_pb from "../../../../../../../../../extproto/ext_pb";
import * as google_protobuf_any_pb from "google-protobuf/google/protobuf/any_pb";

export class PluginSource extends jspb.Message {
  clearFiltersList(): void;
  getFiltersList(): Array<WasmFilter>;
  setFiltersList(value: Array<WasmFilter>): void;
  addFilters(value?: WasmFilter, index?: number): WasmFilter;

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
    filtersList: Array<WasmFilter.AsObject>,
  }
}

export class WasmFilter extends jspb.Message {
  hasImage(): boolean;
  clearImage(): void;
  getImage(): string;
  setImage(value: string): void;

  hasFilePath(): boolean;
  clearFilePath(): void;
  getFilePath(): string;
  setFilePath(value: string): void;

  hasConfig(): boolean;
  clearConfig(): void;
  getConfig(): google_protobuf_any_pb.Any | undefined;
  setConfig(value?: google_protobuf_any_pb.Any): void;

  hasFilterStage(): boolean;
  clearFilterStage(): void;
  getFilterStage(): FilterStage | undefined;
  setFilterStage(value?: FilterStage): void;

  getName(): string;
  setName(value: string): void;

  getRootId(): string;
  setRootId(value: string): void;

  getVmType(): WasmFilter.VmTypeMap[keyof WasmFilter.VmTypeMap];
  setVmType(value: WasmFilter.VmTypeMap[keyof WasmFilter.VmTypeMap]): void;

  getSrcCase(): WasmFilter.SrcCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): WasmFilter.AsObject;
  static toObject(includeInstance: boolean, msg: WasmFilter): WasmFilter.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: WasmFilter, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): WasmFilter;
  static deserializeBinaryFromReader(message: WasmFilter, reader: jspb.BinaryReader): WasmFilter;
}

export namespace WasmFilter {
  export type AsObject = {
    image: string,
    filePath: string,
    config?: google_protobuf_any_pb.Any.AsObject,
    filterStage?: FilterStage.AsObject,
    name: string,
    rootId: string,
    vmType: WasmFilter.VmTypeMap[keyof WasmFilter.VmTypeMap],
  }

  export interface VmTypeMap {
    V8: 0;
    WAVM: 1;
  }

  export const VmType: VmTypeMap;

  export enum SrcCase {
    SRC_NOT_SET = 0,
    IMAGE = 2,
    FILE_PATH = 8,
  }
}

export class FilterStage extends jspb.Message {
  getStage(): FilterStage.StageMap[keyof FilterStage.StageMap];
  setStage(value: FilterStage.StageMap[keyof FilterStage.StageMap]): void;

  getPredicate(): FilterStage.PredicateMap[keyof FilterStage.PredicateMap];
  setPredicate(value: FilterStage.PredicateMap[keyof FilterStage.PredicateMap]): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FilterStage.AsObject;
  static toObject(includeInstance: boolean, msg: FilterStage): FilterStage.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: FilterStage, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FilterStage;
  static deserializeBinaryFromReader(message: FilterStage, reader: jspb.BinaryReader): FilterStage;
}

export namespace FilterStage {
  export type AsObject = {
    stage: FilterStage.StageMap[keyof FilterStage.StageMap],
    predicate: FilterStage.PredicateMap[keyof FilterStage.PredicateMap],
  }

  export interface StageMap {
    FAULTSTAGE: 0;
    CORSSTAGE: 1;
    WAFSTAGE: 2;
    AUTHNSTAGE: 3;
    AUTHZSTAGE: 4;
    RATELIMITSTAGE: 5;
    ACCEPTEDSTAGE: 6;
    OUTAUTHSTAGE: 7;
    ROUTESTAGE: 8;
  }

  export const Stage: StageMap;

  export interface PredicateMap {
    DURING: 0;
    BEFORE: 1;
    AFTER: 2;
  }

  export const Predicate: PredicateMap;
}
