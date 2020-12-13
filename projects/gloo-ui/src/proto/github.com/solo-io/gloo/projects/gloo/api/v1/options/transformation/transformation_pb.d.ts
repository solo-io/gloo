/* eslint-disable */
// package: transformation.options.gloo.solo.io
// file: github.com/solo-io/gloo/projects/gloo/api/v1/options/transformation/transformation.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_v1_core_matchers_matchers_pb from "../../../../../../../../../github.com/solo-io/gloo/projects/gloo/api/v1/core/matchers/matchers_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_transformation_pb from "../../../../../../../../../github.com/solo-io/gloo/projects/gloo/api/external/envoy/extensions/transformation/transformation_pb";
import * as extproto_ext_pb from "../../../../../../../../../protoc-gen-ext/extproto/ext_pb";

export class ResponseMatch extends jspb.Message {
  clearMatchersList(): void;
  getMatchersList(): Array<github_com_solo_io_gloo_projects_gloo_api_v1_core_matchers_matchers_pb.HeaderMatcher>;
  setMatchersList(value: Array<github_com_solo_io_gloo_projects_gloo_api_v1_core_matchers_matchers_pb.HeaderMatcher>): void;
  addMatchers(value?: github_com_solo_io_gloo_projects_gloo_api_v1_core_matchers_matchers_pb.HeaderMatcher, index?: number): github_com_solo_io_gloo_projects_gloo_api_v1_core_matchers_matchers_pb.HeaderMatcher;

  getResponseCodeDetails(): string;
  setResponseCodeDetails(value: string): void;

  hasResponseTransformation(): boolean;
  clearResponseTransformation(): void;
  getResponseTransformation(): github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_transformation_pb.Transformation | undefined;
  setResponseTransformation(value?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_transformation_pb.Transformation): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResponseMatch.AsObject;
  static toObject(includeInstance: boolean, msg: ResponseMatch): ResponseMatch.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ResponseMatch, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResponseMatch;
  static deserializeBinaryFromReader(message: ResponseMatch, reader: jspb.BinaryReader): ResponseMatch;
}

export namespace ResponseMatch {
  export type AsObject = {
    matchersList: Array<github_com_solo_io_gloo_projects_gloo_api_v1_core_matchers_matchers_pb.HeaderMatcher.AsObject>,
    responseCodeDetails: string,
    responseTransformation?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_transformation_pb.Transformation.AsObject,
  }
}

export class RequestMatch extends jspb.Message {
  hasMatcher(): boolean;
  clearMatcher(): void;
  getMatcher(): github_com_solo_io_gloo_projects_gloo_api_v1_core_matchers_matchers_pb.Matcher | undefined;
  setMatcher(value?: github_com_solo_io_gloo_projects_gloo_api_v1_core_matchers_matchers_pb.Matcher): void;

  getClearRouteCache(): boolean;
  setClearRouteCache(value: boolean): void;

  hasRequestTransformation(): boolean;
  clearRequestTransformation(): void;
  getRequestTransformation(): github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_transformation_pb.Transformation | undefined;
  setRequestTransformation(value?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_transformation_pb.Transformation): void;

  hasResponseTransformation(): boolean;
  clearResponseTransformation(): void;
  getResponseTransformation(): github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_transformation_pb.Transformation | undefined;
  setResponseTransformation(value?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_transformation_pb.Transformation): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RequestMatch.AsObject;
  static toObject(includeInstance: boolean, msg: RequestMatch): RequestMatch.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RequestMatch, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RequestMatch;
  static deserializeBinaryFromReader(message: RequestMatch, reader: jspb.BinaryReader): RequestMatch;
}

export namespace RequestMatch {
  export type AsObject = {
    matcher?: github_com_solo_io_gloo_projects_gloo_api_v1_core_matchers_matchers_pb.Matcher.AsObject,
    clearRouteCache: boolean,
    requestTransformation?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_transformation_pb.Transformation.AsObject,
    responseTransformation?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_transformation_pb.Transformation.AsObject,
  }
}

export class Transformations extends jspb.Message {
  hasRequestTransformation(): boolean;
  clearRequestTransformation(): void;
  getRequestTransformation(): github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_transformation_pb.Transformation | undefined;
  setRequestTransformation(value?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_transformation_pb.Transformation): void;

  getClearRouteCache(): boolean;
  setClearRouteCache(value: boolean): void;

  hasResponseTransformation(): boolean;
  clearResponseTransformation(): void;
  getResponseTransformation(): github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_transformation_pb.Transformation | undefined;
  setResponseTransformation(value?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_transformation_pb.Transformation): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Transformations.AsObject;
  static toObject(includeInstance: boolean, msg: Transformations): Transformations.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Transformations, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Transformations;
  static deserializeBinaryFromReader(message: Transformations, reader: jspb.BinaryReader): Transformations;
}

export namespace Transformations {
  export type AsObject = {
    requestTransformation?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_transformation_pb.Transformation.AsObject,
    clearRouteCache: boolean,
    responseTransformation?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_transformation_pb.Transformation.AsObject,
  }
}

export class RequestResponseTransformations extends jspb.Message {
  clearRequestTransformsList(): void;
  getRequestTransformsList(): Array<RequestMatch>;
  setRequestTransformsList(value: Array<RequestMatch>): void;
  addRequestTransforms(value?: RequestMatch, index?: number): RequestMatch;

  clearResponseTransformsList(): void;
  getResponseTransformsList(): Array<ResponseMatch>;
  setResponseTransformsList(value: Array<ResponseMatch>): void;
  addResponseTransforms(value?: ResponseMatch, index?: number): ResponseMatch;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RequestResponseTransformations.AsObject;
  static toObject(includeInstance: boolean, msg: RequestResponseTransformations): RequestResponseTransformations.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RequestResponseTransformations, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RequestResponseTransformations;
  static deserializeBinaryFromReader(message: RequestResponseTransformations, reader: jspb.BinaryReader): RequestResponseTransformations;
}

export namespace RequestResponseTransformations {
  export type AsObject = {
    requestTransformsList: Array<RequestMatch.AsObject>,
    responseTransformsList: Array<ResponseMatch.AsObject>,
  }
}

export class TransformationStages extends jspb.Message {
  hasEarly(): boolean;
  clearEarly(): void;
  getEarly(): RequestResponseTransformations | undefined;
  setEarly(value?: RequestResponseTransformations): void;

  hasRegular(): boolean;
  clearRegular(): void;
  getRegular(): RequestResponseTransformations | undefined;
  setRegular(value?: RequestResponseTransformations): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TransformationStages.AsObject;
  static toObject(includeInstance: boolean, msg: TransformationStages): TransformationStages.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: TransformationStages, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TransformationStages;
  static deserializeBinaryFromReader(message: TransformationStages, reader: jspb.BinaryReader): TransformationStages;
}

export namespace TransformationStages {
  export type AsObject = {
    early?: RequestResponseTransformations.AsObject,
    regular?: RequestResponseTransformations.AsObject,
  }
}
