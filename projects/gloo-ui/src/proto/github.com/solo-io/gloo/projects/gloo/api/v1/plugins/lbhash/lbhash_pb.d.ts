// package: lbhash.plugins.gloo.solo.io
// file: github.com/solo-io/gloo/projects/gloo/api/v1/plugins/lbhash/lbhash.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../../../../../../../gogoproto/gogo_pb";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";
import * as google_protobuf_duration_pb from "google-protobuf/google/protobuf/duration_pb";

export class RouteActionHashConfig extends jspb.Message {
  clearHashPoliciesList(): void;
  getHashPoliciesList(): Array<HashPolicy>;
  setHashPoliciesList(value: Array<HashPolicy>): void;
  addHashPolicies(value?: HashPolicy, index?: number): HashPolicy;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RouteActionHashConfig.AsObject;
  static toObject(includeInstance: boolean, msg: RouteActionHashConfig): RouteActionHashConfig.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RouteActionHashConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RouteActionHashConfig;
  static deserializeBinaryFromReader(message: RouteActionHashConfig, reader: jspb.BinaryReader): RouteActionHashConfig;
}

export namespace RouteActionHashConfig {
  export type AsObject = {
    hashPoliciesList: Array<HashPolicy.AsObject>,
  }
}

export class Cookie extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  hasTtl(): boolean;
  clearTtl(): void;
  getTtl(): google_protobuf_duration_pb.Duration | undefined;
  setTtl(value?: google_protobuf_duration_pb.Duration): void;

  getPath(): string;
  setPath(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Cookie.AsObject;
  static toObject(includeInstance: boolean, msg: Cookie): Cookie.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Cookie, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Cookie;
  static deserializeBinaryFromReader(message: Cookie, reader: jspb.BinaryReader): Cookie;
}

export namespace Cookie {
  export type AsObject = {
    name: string,
    ttl?: google_protobuf_duration_pb.Duration.AsObject,
    path: string,
  }
}

export class HashPolicy extends jspb.Message {
  hasHeader(): boolean;
  clearHeader(): void;
  getHeader(): string;
  setHeader(value: string): void;

  hasCookie(): boolean;
  clearCookie(): void;
  getCookie(): Cookie | undefined;
  setCookie(value?: Cookie): void;

  hasSourceIp(): boolean;
  clearSourceIp(): void;
  getSourceIp(): boolean;
  setSourceIp(value: boolean): void;

  getTerminal(): boolean;
  setTerminal(value: boolean): void;

  getKeytypeCase(): HashPolicy.KeytypeCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HashPolicy.AsObject;
  static toObject(includeInstance: boolean, msg: HashPolicy): HashPolicy.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: HashPolicy, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HashPolicy;
  static deserializeBinaryFromReader(message: HashPolicy, reader: jspb.BinaryReader): HashPolicy;
}

export namespace HashPolicy {
  export type AsObject = {
    header: string,
    cookie?: Cookie.AsObject,
    sourceIp: boolean,
    terminal: boolean,
  }

  export enum KeytypeCase {
    KEYTYPE_NOT_SET = 0,
    HEADER = 1,
    COOKIE = 2,
    SOURCE_IP = 3,
  }
}

