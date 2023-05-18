/* eslint-disable */
// package: envoy.config.matching.custom_matchers.server_name.v3
// file: github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/config/matching/custom_matchers/server_name/v3/server_name_matcher.proto

import * as jspb from "google-protobuf";
import * as validate_validate_pb from "../../../../../../../../../../../../../validate/validate_pb";
import * as xds_type_matcher_v3_matcher_pb from "../../../../../../../../../../../../../xds/type/matcher/v3/matcher_pb";
import * as extproto_ext_pb from "../../../../../../../../../../../../../extproto/ext_pb";

export class ServerNameMatcher extends jspb.Message {
  clearServerNameMatchersList(): void;
  getServerNameMatchersList(): Array<ServerNameMatcher.ServerNameSetMatcher>;
  setServerNameMatchersList(value: Array<ServerNameMatcher.ServerNameSetMatcher>): void;
  addServerNameMatchers(value?: ServerNameMatcher.ServerNameSetMatcher, index?: number): ServerNameMatcher.ServerNameSetMatcher;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ServerNameMatcher.AsObject;
  static toObject(includeInstance: boolean, msg: ServerNameMatcher): ServerNameMatcher.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ServerNameMatcher, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ServerNameMatcher;
  static deserializeBinaryFromReader(message: ServerNameMatcher, reader: jspb.BinaryReader): ServerNameMatcher;
}

export namespace ServerNameMatcher {
  export type AsObject = {
    serverNameMatchersList: Array<ServerNameMatcher.ServerNameSetMatcher.AsObject>,
  }

  export class ServerNameSetMatcher extends jspb.Message {
    clearServerNamesList(): void;
    getServerNamesList(): Array<string>;
    setServerNamesList(value: Array<string>): void;
    addServerNames(value: string, index?: number): string;

    hasOnMatch(): boolean;
    clearOnMatch(): void;
    getOnMatch(): xds_type_matcher_v3_matcher_pb.Matcher.OnMatch | undefined;
    setOnMatch(value?: xds_type_matcher_v3_matcher_pb.Matcher.OnMatch): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ServerNameSetMatcher.AsObject;
    static toObject(includeInstance: boolean, msg: ServerNameSetMatcher): ServerNameSetMatcher.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: ServerNameSetMatcher, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ServerNameSetMatcher;
    static deserializeBinaryFromReader(message: ServerNameSetMatcher, reader: jspb.BinaryReader): ServerNameSetMatcher;
  }

  export namespace ServerNameSetMatcher {
    export type AsObject = {
      serverNamesList: Array<string>,
      onMatch?: xds_type_matcher_v3_matcher_pb.Matcher.OnMatch.AsObject,
    }
  }
}
