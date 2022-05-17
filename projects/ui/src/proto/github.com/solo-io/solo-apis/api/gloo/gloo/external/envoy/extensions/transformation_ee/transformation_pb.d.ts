/* eslint-disable */
// package: envoy.config.filter.http.transformation_ee.v2
// file: github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/extensions/transformation_ee/transformation.proto

import * as jspb from "google-protobuf";
import * as validate_validate_pb from "../../../../../../../../../../validate/validate_pb";
import * as envoy_api_v2_route_route_pb from "../../../../../../../../../../envoy/api/v2/route/route_pb";
import * as github_com_solo_io_solo_kit_api_external_envoy_type_percent_pb from "../../../../../../../../../../github.com/solo-io/solo-kit/api/external/envoy/type/percent_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_route_v3_route_components_pb from "../../../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/config/route/v3/route_components_pb";

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

  hasMatchV3(): boolean;
  clearMatchV3(): void;
  getMatchV3(): github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_route_v3_route_components_pb.RouteMatch | undefined;
  setMatchV3(value?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_route_v3_route_components_pb.RouteMatch): void;

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
    matchV3?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_route_v3_route_components_pb.RouteMatch.AsObject,
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

  hasOnStreamCompletionTransformation(): boolean;
  clearOnStreamCompletionTransformation(): void;
  getOnStreamCompletionTransformation(): Transformation | undefined;
  setOnStreamCompletionTransformation(value?: Transformation): void;

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
    onStreamCompletionTransformation?: Transformation.AsObject,
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

  getEnableHeaderTransformation(): boolean;
  setEnableHeaderTransformation(value: boolean): void;

  getEnableDynamicMetadataTransformation(): boolean;
  setEnableDynamicMetadataTransformation(value: boolean): void;

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
    enableHeaderTransformation: boolean,
    enableDynamicMetadataTransformation: boolean,
  }
}

export class Action extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  clearRegexList(): void;
  getRegexList(): Array<string>;
  setRegexList(value: Array<string>): void;
  addRegex(value: string, index?: number): string;

  clearRegexActionsList(): void;
  getRegexActionsList(): Array<RegexAction>;
  setRegexActionsList(value: Array<RegexAction>): void;
  addRegexActions(value?: RegexAction, index?: number): RegexAction;

  getShadow(): boolean;
  setShadow(value: boolean): void;

  hasPercent(): boolean;
  clearPercent(): void;
  getPercent(): github_com_solo_io_solo_kit_api_external_envoy_type_percent_pb.Percent | undefined;
  setPercent(value?: github_com_solo_io_solo_kit_api_external_envoy_type_percent_pb.Percent): void;

  getMaskChar(): string;
  setMaskChar(value: string): void;

  hasMatcher(): boolean;
  clearMatcher(): void;
  getMatcher(): Action.DlpMatcher | undefined;
  setMatcher(value?: Action.DlpMatcher): void;

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
    regexActionsList: Array<RegexAction.AsObject>,
    shadow: boolean,
    percent?: github_com_solo_io_solo_kit_api_external_envoy_type_percent_pb.Percent.AsObject,
    maskChar: string,
    matcher?: Action.DlpMatcher.AsObject,
  }

  export class RegexMatcher extends jspb.Message {
    clearRegexActionsList(): void;
    getRegexActionsList(): Array<RegexAction>;
    setRegexActionsList(value: Array<RegexAction>): void;
    addRegexActions(value?: RegexAction, index?: number): RegexAction;

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
      regexActionsList: Array<RegexAction.AsObject>,
    }
  }

  export class KeyValueMatcher extends jspb.Message {
    clearKeysList(): void;
    getKeysList(): Array<string>;
    setKeysList(value: Array<string>): void;
    addKeys(value: string, index?: number): string;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): KeyValueMatcher.AsObject;
    static toObject(includeInstance: boolean, msg: KeyValueMatcher): KeyValueMatcher.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: KeyValueMatcher, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): KeyValueMatcher;
    static deserializeBinaryFromReader(message: KeyValueMatcher, reader: jspb.BinaryReader): KeyValueMatcher;
  }

  export namespace KeyValueMatcher {
    export type AsObject = {
      keysList: Array<string>,
    }
  }

  export class DlpMatcher extends jspb.Message {
    hasRegexMatcher(): boolean;
    clearRegexMatcher(): void;
    getRegexMatcher(): Action.RegexMatcher | undefined;
    setRegexMatcher(value?: Action.RegexMatcher): void;

    hasKeyValueMatcher(): boolean;
    clearKeyValueMatcher(): void;
    getKeyValueMatcher(): Action.KeyValueMatcher | undefined;
    setKeyValueMatcher(value?: Action.KeyValueMatcher): void;

    getMatcherCase(): DlpMatcher.MatcherCase;
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): DlpMatcher.AsObject;
    static toObject(includeInstance: boolean, msg: DlpMatcher): DlpMatcher.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: DlpMatcher, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): DlpMatcher;
    static deserializeBinaryFromReader(message: DlpMatcher, reader: jspb.BinaryReader): DlpMatcher;
  }

  export namespace DlpMatcher {
    export type AsObject = {
      regexMatcher?: Action.RegexMatcher.AsObject,
      keyValueMatcher?: Action.KeyValueMatcher.AsObject,
    }

    export enum MatcherCase {
      MATCHER_NOT_SET = 0,
      REGEX_MATCHER = 1,
      KEY_VALUE_MATCHER = 2,
    }
  }
}

export class RegexAction extends jspb.Message {
  getRegex(): string;
  setRegex(value: string): void;

  getSubgroup(): number;
  setSubgroup(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RegexAction.AsObject;
  static toObject(includeInstance: boolean, msg: RegexAction): RegexAction.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RegexAction, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RegexAction;
  static deserializeBinaryFromReader(message: RegexAction, reader: jspb.BinaryReader): RegexAction;
}

export namespace RegexAction {
  export type AsObject = {
    regex: string,
    subgroup: number,
  }
}
