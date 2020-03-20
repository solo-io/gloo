/* eslint-disable */
// package: devportal.solo.io
// file: dev-portal/api/dev-portal/v1/common.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../../gogoproto/gogo_pb";
import * as extproto_ext_pb from "../../../../protoc-gen-ext/extproto/ext_pb";

export class Selector extends jspb.Message {
  getMatchLabelsMap(): jspb.Map<string, string>;
  clearMatchLabelsMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Selector.AsObject;
  static toObject(includeInstance: boolean, msg: Selector): Selector.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Selector, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Selector;
  static deserializeBinaryFromReader(message: Selector, reader: jspb.BinaryReader): Selector;
}

export namespace Selector {
  export type AsObject = {
    matchLabelsMap: Array<[string, string]>,
  }
}

export class DataSource extends jspb.Message {
  hasInlineBytes(): boolean;
  clearInlineBytes(): void;
  getInlineBytes(): Uint8Array | string;
  getInlineBytes_asU8(): Uint8Array;
  getInlineBytes_asB64(): string;
  setInlineBytes(value: Uint8Array | string): void;

  hasFetchUrl(): boolean;
  clearFetchUrl(): void;
  getFetchUrl(): string;
  setFetchUrl(value: string): void;

  hasConfigMap(): boolean;
  clearConfigMap(): void;
  getConfigMap(): DataSource.ConfigMapData | undefined;
  setConfigMap(value?: DataSource.ConfigMapData): void;

  getSourceTypeCase(): DataSource.SourceTypeCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DataSource.AsObject;
  static toObject(includeInstance: boolean, msg: DataSource): DataSource.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DataSource, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DataSource;
  static deserializeBinaryFromReader(message: DataSource, reader: jspb.BinaryReader): DataSource;
}

export namespace DataSource {
  export type AsObject = {
    inlineBytes: Uint8Array | string,
    fetchUrl: string,
    configMap?: DataSource.ConfigMapData.AsObject,
  }

  export class ConfigMapData extends jspb.Message {
    getName(): string;
    setName(value: string): void;

    getNamespace(): string;
    setNamespace(value: string): void;

    getKey(): string;
    setKey(value: string): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ConfigMapData.AsObject;
    static toObject(includeInstance: boolean, msg: ConfigMapData): ConfigMapData.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: ConfigMapData, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ConfigMapData;
    static deserializeBinaryFromReader(message: ConfigMapData, reader: jspb.BinaryReader): ConfigMapData;
  }

  export namespace ConfigMapData {
    export type AsObject = {
      name: string,
      namespace: string,
      key: string,
    }
  }

  export enum SourceTypeCase {
    SOURCE_TYPE_NOT_SET = 0,
    INLINE_BYTES = 1,
    FETCH_URL = 2,
    CONFIG_MAP = 3,
  }
}

export class ObjectRef extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getNamespace(): string;
  setNamespace(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ObjectRef.AsObject;
  static toObject(includeInstance: boolean, msg: ObjectRef): ObjectRef.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ObjectRef, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ObjectRef;
  static deserializeBinaryFromReader(message: ObjectRef, reader: jspb.BinaryReader): ObjectRef;
}

export namespace ObjectRef {
  export type AsObject = {
    name: string,
    namespace: string,
  }
}
