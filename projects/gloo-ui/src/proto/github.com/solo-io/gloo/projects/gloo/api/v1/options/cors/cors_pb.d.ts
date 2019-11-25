// package: cors.options.gloo.solo.io
// file: github.com/solo-io/gloo/projects/gloo/api/v1/options/cors/cors.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../../../../../../../gogoproto/gogo_pb";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";
import * as google_protobuf_duration_pb from "google-protobuf/google/protobuf/duration_pb";

export class CorsPolicy extends jspb.Message {
  clearAllowOriginList(): void;
  getAllowOriginList(): Array<string>;
  setAllowOriginList(value: Array<string>): void;
  addAllowOrigin(value: string, index?: number): string;

  clearAllowOriginRegexList(): void;
  getAllowOriginRegexList(): Array<string>;
  setAllowOriginRegexList(value: Array<string>): void;
  addAllowOriginRegex(value: string, index?: number): string;

  clearAllowMethodsList(): void;
  getAllowMethodsList(): Array<string>;
  setAllowMethodsList(value: Array<string>): void;
  addAllowMethods(value: string, index?: number): string;

  clearAllowHeadersList(): void;
  getAllowHeadersList(): Array<string>;
  setAllowHeadersList(value: Array<string>): void;
  addAllowHeaders(value: string, index?: number): string;

  clearExposeHeadersList(): void;
  getExposeHeadersList(): Array<string>;
  setExposeHeadersList(value: Array<string>): void;
  addExposeHeaders(value: string, index?: number): string;

  getMaxAge(): string;
  setMaxAge(value: string): void;

  getAllowCredentials(): boolean;
  setAllowCredentials(value: boolean): void;

  getDisableForRoute(): boolean;
  setDisableForRoute(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CorsPolicy.AsObject;
  static toObject(includeInstance: boolean, msg: CorsPolicy): CorsPolicy.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CorsPolicy, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CorsPolicy;
  static deserializeBinaryFromReader(message: CorsPolicy, reader: jspb.BinaryReader): CorsPolicy;
}

export namespace CorsPolicy {
  export type AsObject = {
    allowOriginList: Array<string>,
    allowOriginRegexList: Array<string>,
    allowMethodsList: Array<string>,
    allowHeadersList: Array<string>,
    exposeHeadersList: Array<string>,
    maxAge: string,
    allowCredentials: boolean,
    disableForRoute: boolean,
  }
}

