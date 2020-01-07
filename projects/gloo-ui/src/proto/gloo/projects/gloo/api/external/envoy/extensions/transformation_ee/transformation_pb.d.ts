// package: envoy.config.filter.http.transformation_ee.v2
// file: gloo/projects/gloo/api/external/envoy/extensions/transformation_ee/transformation.proto

import * as jspb from "google-protobuf";
import * as validate_validate_pb from "../../../../../../../../validate/validate_pb";
import * as envoy_api_v2_route_route_pb from "../../../../../../../../envoy/api/v2/route/route_pb";
import * as envoy_type_percent_pb from "../../../../../../../../envoy/type/percent_pb";

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
  hasDlpTransformation(): boolean;
  clearDlpTransformation(): void;
  getDlpTransformation(): DlpTransformation | undefined;
  setDlpTransformation(value?: DlpTransformation): void;

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
    dlpTransformation?: DlpTransformation.AsObject,
  }

  export enum TransformationTypeCase {
    TRANSFORMATION_TYPE_NOT_SET = 0,
    DLP_TRANSFORMATION = 1,
  }
}

export class DlpTransformation extends jspb.Message {
  clearActionsList(): void;
  getActionsList(): Array<Action>;
  setActionsList(value: Array<Action>): void;
  addActions(value?: Action, index?: number): Action;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DlpTransformation.AsObject;
  static toObject(includeInstance: boolean, msg: DlpTransformation): DlpTransformation.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DlpTransformation, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DlpTransformation;
  static deserializeBinaryFromReader(message: DlpTransformation, reader: jspb.BinaryReader): DlpTransformation;
}

export namespace DlpTransformation {
  export type AsObject = {
    actionsList: Array<Action.AsObject>,
  }
}

export class Action extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  clearRegexList(): void;
  getRegexList(): Array<string>;
  setRegexList(value: Array<string>): void;
  addRegex(value: string, index?: number): string;

  getShadow(): boolean;
  setShadow(value: boolean): void;

  hasPercent(): boolean;
  clearPercent(): void;
  getPercent(): envoy_type_percent_pb.Percent | undefined;
  setPercent(value?: envoy_type_percent_pb.Percent): void;

  getMaskChar(): string;
  setMaskChar(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Action.AsObject;
  static toObject(includeInstance: boolean, msg: Action): Action.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Action, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Action;
  static deserializeBinaryFromReader(message: Action, reader: jspb.BinaryReader): Action;
}

export namespace Action {
  export type AsObject = {
    name: string,
    regexList: Array<string>,
    shadow: boolean,
    percent?: envoy_type_percent_pb.Percent.AsObject,
    maskChar: string,
  }
}

