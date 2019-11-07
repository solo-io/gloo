// package: envoy.api.v2.filter.http
// file: github.com/solo-io/gloo/projects/gloo/api/external/envoy/extensions/transformation/transformation.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_empty_pb from "google-protobuf/google/protobuf/empty_pb";
import * as envoy_api_v2_route_route_pb from "../../../../../../../../../../envoy/api/v2/route/route_pb";
import * as validate_validate_pb from "../../../../../../../../../../validate/validate_pb";
import * as gogoproto_gogo_pb from "../../../../../../../../../../gogoproto/gogo_pb";

export class FilterTransformations extends jspb.Message {
  clearTransformationsList(): void;
  getTransformationsList(): Array<TransformationRule>;
  setTransformationsList(value: Array<TransformationRule>): void;
  addTransformations(value?: TransformationRule, index?: number): TransformationRule;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FilterTransformations.AsObject;
  static toObject(includeInstance: boolean, msg: FilterTransformations): FilterTransformations.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: FilterTransformations, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FilterTransformations;
  static deserializeBinaryFromReader(message: FilterTransformations, reader: jspb.BinaryReader): FilterTransformations;
}

export namespace FilterTransformations {
  export type AsObject = {
    transformationsList: Array<TransformationRule.AsObject>,
  }
}

export class TransformationRule extends jspb.Message {
  hasMatch(): boolean;
  clearMatch(): void;
  getMatch(): envoy_api_v2_route_route_pb.RouteMatch | undefined;
  setMatch(value?: envoy_api_v2_route_route_pb.RouteMatch): void;

  hasRouteTransformations(): boolean;
  clearRouteTransformations(): void;
  getRouteTransformations(): RouteTransformations | undefined;
  setRouteTransformations(value?: RouteTransformations): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TransformationRule.AsObject;
  static toObject(includeInstance: boolean, msg: TransformationRule): TransformationRule.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: TransformationRule, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TransformationRule;
  static deserializeBinaryFromReader(message: TransformationRule, reader: jspb.BinaryReader): TransformationRule;
}

export namespace TransformationRule {
  export type AsObject = {
    match?: envoy_api_v2_route_route_pb.RouteMatch.AsObject,
    routeTransformations?: RouteTransformations.AsObject,
  }
}

export class RouteTransformations extends jspb.Message {
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
  toObject(includeInstance?: boolean): RouteTransformations.AsObject;
  static toObject(includeInstance: boolean, msg: RouteTransformations): RouteTransformations.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RouteTransformations, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RouteTransformations;
  static deserializeBinaryFromReader(message: RouteTransformations, reader: jspb.BinaryReader): RouteTransformations;
}

export namespace RouteTransformations {
  export type AsObject = {
    requestTransformation?: Transformation.AsObject,
    clearRouteCache: boolean,
    responseTransformation?: Transformation.AsObject,
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
  }

  export enum TransformationTypeCase {
    TRANSFORMATION_TYPE_NOT_SET = 0,
    TRANSFORMATION_TEMPLATE = 1,
    HEADER_BODY_TRANSFORM = 2,
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
    body?: InjaTemplate.AsObject,
    passthrough?: Passthrough.AsObject,
    mergeExtractorsToBody?: MergeExtractorsToBody.AsObject,
    parseBodyBehavior: TransformationTemplate.RequestBodyParseMap[keyof TransformationTemplate.RequestBodyParseMap],
    ignoreErrorOnParse: boolean,
    dynamicMetadataValuesList: Array<TransformationTemplate.DynamicMetadataValue.AsObject>,
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
  }
}

