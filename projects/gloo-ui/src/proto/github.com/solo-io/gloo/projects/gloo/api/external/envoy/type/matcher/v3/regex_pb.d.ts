/* eslint-disable */
// package: solo.io.envoy.type.matcher.v3
// file: github.com/solo-io/gloo/projects/gloo/api/external/envoy/type/matcher/v3/regex.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";
import * as udpa_annotations_status_pb from "../../../../../../../../../../../udpa/annotations/status_pb";
import * as udpa_annotations_versioning_pb from "../../../../../../../../../../../udpa/annotations/versioning_pb";
import * as validate_validate_pb from "../../../../../../../../../../../validate/validate_pb";
import * as extproto_ext_pb from "../../../../../../../../../../../protoc-gen-ext/extproto/ext_pb";

export class RegexMatcher extends jspb.Message {
  hasGoogleRe2(): boolean;
  clearGoogleRe2(): void;
  getGoogleRe2(): RegexMatcher.GoogleRE2 | undefined;
  setGoogleRe2(value?: RegexMatcher.GoogleRE2): void;

  getRegex(): string;
  setRegex(value: string): void;

  getEngineTypeCase(): RegexMatcher.EngineTypeCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RegexMatcher.AsObject;
  static toObject(includeInstance: boolean, msg: RegexMatcher): RegexMatcher.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RegexMatcher, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RegexMatcher;
  static deserializeBinaryFromReader(message: RegexMatcher, reader: jspb.BinaryReader): RegexMatcher;
}

export namespace RegexMatcher {
  export type AsObject = {
    googleRe2?: RegexMatcher.GoogleRE2.AsObject,
    regex: string,
  }

  export class GoogleRE2 extends jspb.Message {
    hasMaxProgramSize(): boolean;
    clearMaxProgramSize(): void;
    getMaxProgramSize(): google_protobuf_wrappers_pb.UInt32Value | undefined;
    setMaxProgramSize(value?: google_protobuf_wrappers_pb.UInt32Value): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): GoogleRE2.AsObject;
    static toObject(includeInstance: boolean, msg: GoogleRE2): GoogleRE2.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: GoogleRE2, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): GoogleRE2;
    static deserializeBinaryFromReader(message: GoogleRE2, reader: jspb.BinaryReader): GoogleRE2;
  }

  export namespace GoogleRE2 {
    export type AsObject = {
      maxProgramSize?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
    }
  }

  export enum EngineTypeCase {
    ENGINE_TYPE_NOT_SET = 0,
    GOOGLE_RE2 = 1,
  }
}

export class RegexMatchAndSubstitute extends jspb.Message {
  hasPattern(): boolean;
  clearPattern(): void;
  getPattern(): RegexMatcher | undefined;
  setPattern(value?: RegexMatcher): void;

  getSubstitution(): string;
  setSubstitution(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RegexMatchAndSubstitute.AsObject;
  static toObject(includeInstance: boolean, msg: RegexMatchAndSubstitute): RegexMatchAndSubstitute.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RegexMatchAndSubstitute, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RegexMatchAndSubstitute;
  static deserializeBinaryFromReader(message: RegexMatchAndSubstitute, reader: jspb.BinaryReader): RegexMatchAndSubstitute;
}

export namespace RegexMatchAndSubstitute {
  export type AsObject = {
    pattern?: RegexMatcher.AsObject,
    substitution: string,
  }
}
