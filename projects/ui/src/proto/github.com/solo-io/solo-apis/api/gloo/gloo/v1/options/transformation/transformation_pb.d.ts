/* eslint-disable */
// package: transformation.options.gloo.solo.io
// file: github.com/solo-io/solo-apis/api/gloo/gloo/v1/options/transformation/transformation.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_empty_pb from "google-protobuf/google/protobuf/empty_pb";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_core_matchers_matchers_pb from "../../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/core/matchers/matchers_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_extensions_transformers_xslt_xslt_transformer_pb from "../../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/extensions/transformers/xslt/xslt_transformer_pb";
import * as extproto_ext_pb from "../../../../../../../../../extproto/ext_pb";

export class ResponseMatch extends jspb.Message {
  clearMatchersList(): void;
  getMatchersList(): Array<github_com_solo_io_solo_apis_api_gloo_gloo_v1_core_matchers_matchers_pb.HeaderMatcher>;
  setMatchersList(value: Array<github_com_solo_io_solo_apis_api_gloo_gloo_v1_core_matchers_matchers_pb.HeaderMatcher>): void;
  addMatchers(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_core_matchers_matchers_pb.HeaderMatcher, index?: number): github_com_solo_io_solo_apis_api_gloo_gloo_v1_core_matchers_matchers_pb.HeaderMatcher;

  getResponseCodeDetails(): string;
  setResponseCodeDetails(value: string): void;

  hasResponseTransformation(): boolean;
  clearResponseTransformation(): void;
  getResponseTransformation(): Transformation | undefined;
  setResponseTransformation(value?: Transformation): void;

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
    matchersList: Array<github_com_solo_io_solo_apis_api_gloo_gloo_v1_core_matchers_matchers_pb.HeaderMatcher.AsObject>,
    responseCodeDetails: string,
    responseTransformation?: Transformation.AsObject,
  }
}

export class RequestMatch extends jspb.Message {
  hasMatcher(): boolean;
  clearMatcher(): void;
  getMatcher(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_core_matchers_matchers_pb.Matcher | undefined;
  setMatcher(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_core_matchers_matchers_pb.Matcher): void;

  getClearRouteCache(): boolean;
  setClearRouteCache(value: boolean): void;

  hasRequestTransformation(): boolean;
  clearRequestTransformation(): void;
  getRequestTransformation(): Transformation | undefined;
  setRequestTransformation(value?: Transformation): void;

  hasResponseTransformation(): boolean;
  clearResponseTransformation(): void;
  getResponseTransformation(): Transformation | undefined;
  setResponseTransformation(value?: Transformation): void;

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
    matcher?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_core_matchers_matchers_pb.Matcher.AsObject,
    clearRouteCache: boolean,
    requestTransformation?: Transformation.AsObject,
    responseTransformation?: Transformation.AsObject,
  }
}

export class Transformations extends jspb.Message {
  hasRequestTransformation(): boolean;
  clearRequestTransformation(): void;
  getRequestTransformation(): Transformation | undefined;
  setRequestTransformation(value?: Transformation): void;

  getClearRouteCache(): boolean;
  setClearRouteCache(value: boolean): void;

  hasResponseTransformation(): boolean;
  clearResponseTransformation(): void;
  getResponseTransformation(): Transformation | undefined;
  setResponseTransformation(value?: Transformation): void;

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
    requestTransformation?: Transformation.AsObject,
    clearRouteCache: boolean,
    responseTransformation?: Transformation.AsObject,
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

  getInheritTransformation(): boolean;
  setInheritTransformation(value: boolean): void;

  hasLogRequestResponseInfo(): boolean;
  clearLogRequestResponseInfo(): void;
  getLogRequestResponseInfo(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setLogRequestResponseInfo(value?: google_protobuf_wrappers_pb.BoolValue): void;

  hasEscapeCharacters(): boolean;
  clearEscapeCharacters(): void;
  getEscapeCharacters(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setEscapeCharacters(value?: google_protobuf_wrappers_pb.BoolValue): void;

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
    inheritTransformation: boolean,
    logRequestResponseInfo?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    escapeCharacters?: google_protobuf_wrappers_pb.BoolValue.AsObject,
  }
}

export class Transformation extends jspb.Message {
  hasTransformationTemplate(): boolean;
  clearTransformationTemplate(): void;
  getTransformationTemplate(): TransformationTemplate | undefined;
  setTransformationTemplate(value?: TransformationTemplate): void;

  hasHeaderBodyTransform(): boolean;
  clearHeaderBodyTransform(): void;
  getHeaderBodyTransform(): HeaderBodyTransform | undefined;
  setHeaderBodyTransform(value?: HeaderBodyTransform): void;

  hasXsltTransformation(): boolean;
  clearXsltTransformation(): void;
  getXsltTransformation(): github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_extensions_transformers_xslt_xslt_transformer_pb.XsltTransformation | undefined;
  setXsltTransformation(value?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_extensions_transformers_xslt_xslt_transformer_pb.XsltTransformation): void;

  getLogRequestResponseInfo(): boolean;
  setLogRequestResponseInfo(value: boolean): void;

  getTransformationTypeCase(): Transformation.TransformationTypeCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Transformation.AsObject;
  static toObject(includeInstance: boolean, msg: Transformation): Transformation.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Transformation, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Transformation;
  static deserializeBinaryFromReader(message: Transformation, reader: jspb.BinaryReader): Transformation;
}

export namespace Transformation {
  export type AsObject = {
    transformationTemplate?: TransformationTemplate.AsObject,
    headerBodyTransform?: HeaderBodyTransform.AsObject,
    xsltTransformation?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_extensions_transformers_xslt_xslt_transformer_pb.XsltTransformation.AsObject,
    logRequestResponseInfo: boolean,
  }

  export enum TransformationTypeCase {
    TRANSFORMATION_TYPE_NOT_SET = 0,
    TRANSFORMATION_TEMPLATE = 1,
    HEADER_BODY_TRANSFORM = 2,
    XSLT_TRANSFORMATION = 3,
  }
}

export class Extraction extends jspb.Message {
  hasHeader(): boolean;
  clearHeader(): void;
  getHeader(): string;
  setHeader(value: string): void;

  hasBody(): boolean;
  clearBody(): void;
  getBody(): google_protobuf_empty_pb.Empty | undefined;
  setBody(value?: google_protobuf_empty_pb.Empty): void;

  getRegex(): string;
  setRegex(value: string): void;

  getSubgroup(): number;
  setSubgroup(value: number): void;

  getSourceCase(): Extraction.SourceCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Extraction.AsObject;
  static toObject(includeInstance: boolean, msg: Extraction): Extraction.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Extraction, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Extraction;
  static deserializeBinaryFromReader(message: Extraction, reader: jspb.BinaryReader): Extraction;
}

export namespace Extraction {
  export type AsObject = {
    header: string,
    body?: google_protobuf_empty_pb.Empty.AsObject,
    regex: string,
    subgroup: number,
  }

  export enum SourceCase {
    SOURCE_NOT_SET = 0,
    HEADER = 1,
    BODY = 4,
  }
}

export class TransformationTemplate extends jspb.Message {
  getAdvancedTemplates(): boolean;
  setAdvancedTemplates(value: boolean): void;

  getExtractorsMap(): jspb.Map<string, Extraction>;
  clearExtractorsMap(): void;
  getHeadersMap(): jspb.Map<string, InjaTemplate>;
  clearHeadersMap(): void;
  clearHeadersToAppendList(): void;
  getHeadersToAppendList(): Array<TransformationTemplate.HeaderToAppend>;
  setHeadersToAppendList(value: Array<TransformationTemplate.HeaderToAppend>): void;
  addHeadersToAppend(value?: TransformationTemplate.HeaderToAppend, index?: number): TransformationTemplate.HeaderToAppend;

  clearHeadersToRemoveList(): void;
  getHeadersToRemoveList(): Array<string>;
  setHeadersToRemoveList(value: Array<string>): void;
  addHeadersToRemove(value: string, index?: number): string;

  hasBody(): boolean;
  clearBody(): void;
  getBody(): InjaTemplate | undefined;
  setBody(value?: InjaTemplate): void;

  hasPassthrough(): boolean;
  clearPassthrough(): void;
  getPassthrough(): Passthrough | undefined;
  setPassthrough(value?: Passthrough): void;

  hasMergeExtractorsToBody(): boolean;
  clearMergeExtractorsToBody(): void;
  getMergeExtractorsToBody(): MergeExtractorsToBody | undefined;
  setMergeExtractorsToBody(value?: MergeExtractorsToBody): void;

  getParseBodyBehavior(): TransformationTemplate.RequestBodyParseMap[keyof TransformationTemplate.RequestBodyParseMap];
  setParseBodyBehavior(value: TransformationTemplate.RequestBodyParseMap[keyof TransformationTemplate.RequestBodyParseMap]): void;

  getIgnoreErrorOnParse(): boolean;
  setIgnoreErrorOnParse(value: boolean): void;

  clearDynamicMetadataValuesList(): void;
  getDynamicMetadataValuesList(): Array<TransformationTemplate.DynamicMetadataValue>;
  setDynamicMetadataValuesList(value: Array<TransformationTemplate.DynamicMetadataValue>): void;
  addDynamicMetadataValues(value?: TransformationTemplate.DynamicMetadataValue, index?: number): TransformationTemplate.DynamicMetadataValue;

  hasEscapeCharacters(): boolean;
  clearEscapeCharacters(): void;
  getEscapeCharacters(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setEscapeCharacters(value?: google_protobuf_wrappers_pb.BoolValue): void;

  getBodyTransformationCase(): TransformationTemplate.BodyTransformationCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TransformationTemplate.AsObject;
  static toObject(includeInstance: boolean, msg: TransformationTemplate): TransformationTemplate.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: TransformationTemplate, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TransformationTemplate;
  static deserializeBinaryFromReader(message: TransformationTemplate, reader: jspb.BinaryReader): TransformationTemplate;
}

export namespace TransformationTemplate {
  export type AsObject = {
    advancedTemplates: boolean,
    extractorsMap: Array<[string, Extraction.AsObject]>,
    headersMap: Array<[string, InjaTemplate.AsObject]>,
    headersToAppendList: Array<TransformationTemplate.HeaderToAppend.AsObject>,
    headersToRemoveList: Array<string>,
    body?: InjaTemplate.AsObject,
    passthrough?: Passthrough.AsObject,
    mergeExtractorsToBody?: MergeExtractorsToBody.AsObject,
    parseBodyBehavior: TransformationTemplate.RequestBodyParseMap[keyof TransformationTemplate.RequestBodyParseMap],
    ignoreErrorOnParse: boolean,
    dynamicMetadataValuesList: Array<TransformationTemplate.DynamicMetadataValue.AsObject>,
    escapeCharacters?: google_protobuf_wrappers_pb.BoolValue.AsObject,
  }

  export class HeaderToAppend extends jspb.Message {
    getKey(): string;
    setKey(value: string): void;

    hasValue(): boolean;
    clearValue(): void;
    getValue(): InjaTemplate | undefined;
    setValue(value?: InjaTemplate): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): HeaderToAppend.AsObject;
    static toObject(includeInstance: boolean, msg: HeaderToAppend): HeaderToAppend.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: HeaderToAppend, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): HeaderToAppend;
    static deserializeBinaryFromReader(message: HeaderToAppend, reader: jspb.BinaryReader): HeaderToAppend;
  }

  export namespace HeaderToAppend {
    export type AsObject = {
      key: string,
      value?: InjaTemplate.AsObject,
    }
  }

  export class DynamicMetadataValue extends jspb.Message {
    getMetadataNamespace(): string;
    setMetadataNamespace(value: string): void;

    getKey(): string;
    setKey(value: string): void;

    hasValue(): boolean;
    clearValue(): void;
    getValue(): InjaTemplate | undefined;
    setValue(value?: InjaTemplate): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): DynamicMetadataValue.AsObject;
    static toObject(includeInstance: boolean, msg: DynamicMetadataValue): DynamicMetadataValue.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: DynamicMetadataValue, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): DynamicMetadataValue;
    static deserializeBinaryFromReader(message: DynamicMetadataValue, reader: jspb.BinaryReader): DynamicMetadataValue;
  }

  export namespace DynamicMetadataValue {
    export type AsObject = {
      metadataNamespace: string,
      key: string,
      value?: InjaTemplate.AsObject,
    }
  }

  export interface RequestBodyParseMap {
    PARSEASJSON: 0;
    DONTPARSE: 1;
  }

  export const RequestBodyParse: RequestBodyParseMap;

  export enum BodyTransformationCase {
    BODY_TRANSFORMATION_NOT_SET = 0,
    BODY = 4,
    PASSTHROUGH = 5,
    MERGE_EXTRACTORS_TO_BODY = 6,
  }
}

export class InjaTemplate extends jspb.Message {
  getText(): string;
  setText(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): InjaTemplate.AsObject;
  static toObject(includeInstance: boolean, msg: InjaTemplate): InjaTemplate.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: InjaTemplate, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): InjaTemplate;
  static deserializeBinaryFromReader(message: InjaTemplate, reader: jspb.BinaryReader): InjaTemplate;
}

export namespace InjaTemplate {
  export type AsObject = {
    text: string,
  }
}

export class Passthrough extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Passthrough.AsObject;
  static toObject(includeInstance: boolean, msg: Passthrough): Passthrough.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Passthrough, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Passthrough;
  static deserializeBinaryFromReader(message: Passthrough, reader: jspb.BinaryReader): Passthrough;
}

export namespace Passthrough {
  export type AsObject = {
  }
}

export class MergeExtractorsToBody extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MergeExtractorsToBody.AsObject;
  static toObject(includeInstance: boolean, msg: MergeExtractorsToBody): MergeExtractorsToBody.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: MergeExtractorsToBody, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MergeExtractorsToBody;
  static deserializeBinaryFromReader(message: MergeExtractorsToBody, reader: jspb.BinaryReader): MergeExtractorsToBody;
}

export namespace MergeExtractorsToBody {
  export type AsObject = {
  }
}

export class HeaderBodyTransform extends jspb.Message {
  getAddRequestMetadata(): boolean;
  setAddRequestMetadata(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HeaderBodyTransform.AsObject;
  static toObject(includeInstance: boolean, msg: HeaderBodyTransform): HeaderBodyTransform.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: HeaderBodyTransform, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HeaderBodyTransform;
  static deserializeBinaryFromReader(message: HeaderBodyTransform, reader: jspb.BinaryReader): HeaderBodyTransform;
}

export namespace HeaderBodyTransform {
  export type AsObject = {
    addRequestMetadata: boolean,
  }
}
