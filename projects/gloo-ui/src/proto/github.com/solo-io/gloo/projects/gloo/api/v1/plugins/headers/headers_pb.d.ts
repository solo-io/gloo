// package: headers.plugins.gloo.solo.io
// file: github.com/solo-io/gloo/projects/gloo/api/v1/plugins/headers/headers.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../../../../../../../gogoproto/gogo_pb";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";

export class HeaderManipulation extends jspb.Message {
  clearRequestHeadersToAddList(): void;
  getRequestHeadersToAddList(): Array<HeaderValueOption>;
  setRequestHeadersToAddList(value: Array<HeaderValueOption>): void;
  addRequestHeadersToAdd(value?: HeaderValueOption, index?: number): HeaderValueOption;

  clearRequestHeadersToRemoveList(): void;
  getRequestHeadersToRemoveList(): Array<string>;
  setRequestHeadersToRemoveList(value: Array<string>): void;
  addRequestHeadersToRemove(value: string, index?: number): string;

  clearResponseHeadersToAddList(): void;
  getResponseHeadersToAddList(): Array<HeaderValueOption>;
  setResponseHeadersToAddList(value: Array<HeaderValueOption>): void;
  addResponseHeadersToAdd(value?: HeaderValueOption, index?: number): HeaderValueOption;

  clearResponseHeadersToRemoveList(): void;
  getResponseHeadersToRemoveList(): Array<string>;
  setResponseHeadersToRemoveList(value: Array<string>): void;
  addResponseHeadersToRemove(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HeaderManipulation.AsObject;
  static toObject(includeInstance: boolean, msg: HeaderManipulation): HeaderManipulation.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: HeaderManipulation, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HeaderManipulation;
  static deserializeBinaryFromReader(message: HeaderManipulation, reader: jspb.BinaryReader): HeaderManipulation;
}

export namespace HeaderManipulation {
  export type AsObject = {
    requestHeadersToAddList: Array<HeaderValueOption.AsObject>,
    requestHeadersToRemoveList: Array<string>,
    responseHeadersToAddList: Array<HeaderValueOption.AsObject>,
    responseHeadersToRemoveList: Array<string>,
  }
}

export class HeaderValueOption extends jspb.Message {
  hasHeader(): boolean;
  clearHeader(): void;
  getHeader(): HeaderValue | undefined;
  setHeader(value?: HeaderValue): void;

  hasAppend(): boolean;
  clearAppend(): void;
  getAppend(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setAppend(value?: google_protobuf_wrappers_pb.BoolValue): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HeaderValueOption.AsObject;
  static toObject(includeInstance: boolean, msg: HeaderValueOption): HeaderValueOption.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: HeaderValueOption, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HeaderValueOption;
  static deserializeBinaryFromReader(message: HeaderValueOption, reader: jspb.BinaryReader): HeaderValueOption;
}

export namespace HeaderValueOption {
  export type AsObject = {
    header?: HeaderValue.AsObject,
    append?: google_protobuf_wrappers_pb.BoolValue.AsObject,
  }
}

export class HeaderValue extends jspb.Message {
  getKey(): string;
  setKey(value: string): void;

  getValue(): string;
  setValue(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HeaderValue.AsObject;
  static toObject(includeInstance: boolean, msg: HeaderValue): HeaderValue.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: HeaderValue, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HeaderValue;
  static deserializeBinaryFromReader(message: HeaderValue, reader: jspb.BinaryReader): HeaderValue;
}

export namespace HeaderValue {
  export type AsObject = {
    key: string,
    value: string,
  }
}

