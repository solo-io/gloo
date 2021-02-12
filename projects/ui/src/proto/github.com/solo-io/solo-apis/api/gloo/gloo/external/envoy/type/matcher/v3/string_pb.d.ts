/* eslint-disable */
// package: solo.io.envoy.type.matcher.v3
// file: github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/type/matcher/v3/string.proto

import * as jspb from "google-protobuf";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_type_matcher_v3_regex_pb from "../../../../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/type/matcher/v3/regex_pb";
import * as envoy_annotations_deprecation_pb from "../../../../../../../../../../../envoy/annotations/deprecation_pb";
import * as udpa_annotations_status_pb from "../../../../../../../../../../../udpa/annotations/status_pb";
import * as udpa_annotations_versioning_pb from "../../../../../../../../../../../udpa/annotations/versioning_pb";
import * as validate_validate_pb from "../../../../../../../../../../../validate/validate_pb";
import * as extproto_ext_pb from "../../../../../../../../../../../extproto/ext_pb";

export class StringMatcher extends jspb.Message {
  hasExact(): boolean;
  clearExact(): void;
  getExact(): string;
  setExact(value: string): void;

  hasPrefix(): boolean;
  clearPrefix(): void;
  getPrefix(): string;
  setPrefix(value: string): void;

  hasSuffix(): boolean;
  clearSuffix(): void;
  getSuffix(): string;
  setSuffix(value: string): void;

  hasSafeRegex(): boolean;
  clearSafeRegex(): void;
  getSafeRegex(): github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_type_matcher_v3_regex_pb.RegexMatcher | undefined;
  setSafeRegex(value?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_type_matcher_v3_regex_pb.RegexMatcher): void;

  getIgnoreCase(): boolean;
  setIgnoreCase(value: boolean): void;

  getMatchPatternCase(): StringMatcher.MatchPatternCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StringMatcher.AsObject;
  static toObject(includeInstance: boolean, msg: StringMatcher): StringMatcher.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: StringMatcher, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StringMatcher;
  static deserializeBinaryFromReader(message: StringMatcher, reader: jspb.BinaryReader): StringMatcher;
}

export namespace StringMatcher {
  export type AsObject = {
    exact: string,
    prefix: string,
    suffix: string,
    safeRegex?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_type_matcher_v3_regex_pb.RegexMatcher.AsObject,
    ignoreCase: boolean,
  }

  export enum MatchPatternCase {
    MATCH_PATTERN_NOT_SET = 0,
    EXACT = 1,
    PREFIX = 2,
    SUFFIX = 3,
    SAFE_REGEX = 5,
  }
}

export class ListStringMatcher extends jspb.Message {
  clearPatternsList(): void;
  getPatternsList(): Array<StringMatcher>;
  setPatternsList(value: Array<StringMatcher>): void;
  addPatterns(value?: StringMatcher, index?: number): StringMatcher;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListStringMatcher.AsObject;
  static toObject(includeInstance: boolean, msg: ListStringMatcher): ListStringMatcher.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListStringMatcher, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListStringMatcher;
  static deserializeBinaryFromReader(message: ListStringMatcher, reader: jspb.BinaryReader): ListStringMatcher;
}

export namespace ListStringMatcher {
  export type AsObject = {
    patternsList: Array<StringMatcher.AsObject>,
  }
}
