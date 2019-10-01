// package: ratelimit.plugins.gloo.solo.io
// file: github.com/solo-io/gloo/projects/gloo/api/v1/enterprise/plugins/ratelimit/ratelimit.proto

import * as jspb from "google-protobuf";
import * as github_com_solo_io_solo_kit_api_v1_ref_pb from "../../../../../../../../../../github.com/solo-io/solo-kit/api/v1/ref_pb";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";
import * as google_protobuf_duration_pb from "google-protobuf/google/protobuf/duration_pb";
import * as gogoproto_gogo_pb from "../../../../../../../../../../gogoproto/gogo_pb";

export class Descriptor extends jspb.Message {
  getKey(): string;
  setKey(value: string): void;

  getValue(): string;
  setValue(value: string): void;

  hasRateLimit(): boolean;
  clearRateLimit(): void;
  getRateLimit(): RateLimit | undefined;
  setRateLimit(value?: RateLimit): void;

  clearDescriptorsList(): void;
  getDescriptorsList(): Array<Descriptor>;
  setDescriptorsList(value: Array<Descriptor>): void;
  addDescriptors(value?: Descriptor, index?: number): Descriptor;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Descriptor.AsObject;
  static toObject(includeInstance: boolean, msg: Descriptor): Descriptor.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Descriptor, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Descriptor;
  static deserializeBinaryFromReader(message: Descriptor, reader: jspb.BinaryReader): Descriptor;
}

export namespace Descriptor {
  export type AsObject = {
    key: string,
    value: string,
    rateLimit?: RateLimit.AsObject,
    descriptorsList: Array<Descriptor.AsObject>,
  }
}

export class RateLimit extends jspb.Message {
  getUnit(): RateLimit.UnitMap[keyof RateLimit.UnitMap];
  setUnit(value: RateLimit.UnitMap[keyof RateLimit.UnitMap]): void;

  getRequestsPerUnit(): number;
  setRequestsPerUnit(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RateLimit.AsObject;
  static toObject(includeInstance: boolean, msg: RateLimit): RateLimit.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RateLimit, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RateLimit;
  static deserializeBinaryFromReader(message: RateLimit, reader: jspb.BinaryReader): RateLimit;
}

export namespace RateLimit {
  export type AsObject = {
    unit: RateLimit.UnitMap[keyof RateLimit.UnitMap],
    requestsPerUnit: number,
  }

  export interface UnitMap {
    UNKNOWN: 0;
    SECOND: 1;
    MINUTE: 2;
    HOUR: 3;
    DAY: 4;
  }

  export const Unit: UnitMap;
}

export class IngressRateLimit extends jspb.Message {
  hasAuthorizedLimits(): boolean;
  clearAuthorizedLimits(): void;
  getAuthorizedLimits(): RateLimit | undefined;
  setAuthorizedLimits(value?: RateLimit): void;

  hasAnonymousLimits(): boolean;
  clearAnonymousLimits(): void;
  getAnonymousLimits(): RateLimit | undefined;
  setAnonymousLimits(value?: RateLimit): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): IngressRateLimit.AsObject;
  static toObject(includeInstance: boolean, msg: IngressRateLimit): IngressRateLimit.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: IngressRateLimit, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): IngressRateLimit;
  static deserializeBinaryFromReader(message: IngressRateLimit, reader: jspb.BinaryReader): IngressRateLimit;
}

export namespace IngressRateLimit {
  export type AsObject = {
    authorizedLimits?: RateLimit.AsObject,
    anonymousLimits?: RateLimit.AsObject,
  }
}

export class Settings extends jspb.Message {
  hasRatelimitServerRef(): boolean;
  clearRatelimitServerRef(): void;
  getRatelimitServerRef(): github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef | undefined;
  setRatelimitServerRef(value?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef): void;

  hasRequestTimeout(): boolean;
  clearRequestTimeout(): void;
  getRequestTimeout(): google_protobuf_duration_pb.Duration | undefined;
  setRequestTimeout(value?: google_protobuf_duration_pb.Duration): void;

  getDenyOnFail(): boolean;
  setDenyOnFail(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Settings.AsObject;
  static toObject(includeInstance: boolean, msg: Settings): Settings.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Settings, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Settings;
  static deserializeBinaryFromReader(message: Settings, reader: jspb.BinaryReader): Settings;
}

export namespace Settings {
  export type AsObject = {
    ratelimitServerRef?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject,
    requestTimeout?: google_protobuf_duration_pb.Duration.AsObject,
    denyOnFail: boolean,
  }
}

export class EnvoySettings extends jspb.Message {
  hasCustomConfig(): boolean;
  clearCustomConfig(): void;
  getCustomConfig(): EnvoySettings.RateLimitCustomConfig | undefined;
  setCustomConfig(value?: EnvoySettings.RateLimitCustomConfig): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EnvoySettings.AsObject;
  static toObject(includeInstance: boolean, msg: EnvoySettings): EnvoySettings.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: EnvoySettings, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EnvoySettings;
  static deserializeBinaryFromReader(message: EnvoySettings, reader: jspb.BinaryReader): EnvoySettings;
}

export namespace EnvoySettings {
  export type AsObject = {
    customConfig?: EnvoySettings.RateLimitCustomConfig.AsObject,
  }

  export class RateLimitCustomConfig extends jspb.Message {
    clearDescriptorsList(): void;
    getDescriptorsList(): Array<Descriptor>;
    setDescriptorsList(value: Array<Descriptor>): void;
    addDescriptors(value?: Descriptor, index?: number): Descriptor;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): RateLimitCustomConfig.AsObject;
    static toObject(includeInstance: boolean, msg: RateLimitCustomConfig): RateLimitCustomConfig.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: RateLimitCustomConfig, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): RateLimitCustomConfig;
    static deserializeBinaryFromReader(message: RateLimitCustomConfig, reader: jspb.BinaryReader): RateLimitCustomConfig;
  }

  export namespace RateLimitCustomConfig {
    export type AsObject = {
      descriptorsList: Array<Descriptor.AsObject>,
    }
  }
}

export class ServiceSettings extends jspb.Message {
  clearDescriptorsList(): void;
  getDescriptorsList(): Array<Descriptor>;
  setDescriptorsList(value: Array<Descriptor>): void;
  addDescriptors(value?: Descriptor, index?: number): Descriptor;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ServiceSettings.AsObject;
  static toObject(includeInstance: boolean, msg: ServiceSettings): ServiceSettings.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ServiceSettings, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ServiceSettings;
  static deserializeBinaryFromReader(message: ServiceSettings, reader: jspb.BinaryReader): ServiceSettings;
}

export namespace ServiceSettings {
  export type AsObject = {
    descriptorsList: Array<Descriptor.AsObject>,
  }
}

export class RateLimitActions extends jspb.Message {
  clearActionsList(): void;
  getActionsList(): Array<Action>;
  setActionsList(value: Array<Action>): void;
  addActions(value?: Action, index?: number): Action;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RateLimitActions.AsObject;
  static toObject(includeInstance: boolean, msg: RateLimitActions): RateLimitActions.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RateLimitActions, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RateLimitActions;
  static deserializeBinaryFromReader(message: RateLimitActions, reader: jspb.BinaryReader): RateLimitActions;
}

export namespace RateLimitActions {
  export type AsObject = {
    actionsList: Array<Action.AsObject>,
  }
}

export class RateLimitVhostExtension extends jspb.Message {
  clearRateLimitsList(): void;
  getRateLimitsList(): Array<RateLimitActions>;
  setRateLimitsList(value: Array<RateLimitActions>): void;
  addRateLimits(value?: RateLimitActions, index?: number): RateLimitActions;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RateLimitVhostExtension.AsObject;
  static toObject(includeInstance: boolean, msg: RateLimitVhostExtension): RateLimitVhostExtension.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RateLimitVhostExtension, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RateLimitVhostExtension;
  static deserializeBinaryFromReader(message: RateLimitVhostExtension, reader: jspb.BinaryReader): RateLimitVhostExtension;
}

export namespace RateLimitVhostExtension {
  export type AsObject = {
    rateLimitsList: Array<RateLimitActions.AsObject>,
  }
}

export class RateLimitRouteExtension extends jspb.Message {
  getIncludeVhRateLimits(): boolean;
  setIncludeVhRateLimits(value: boolean): void;

  clearRateLimitsList(): void;
  getRateLimitsList(): Array<RateLimitActions>;
  setRateLimitsList(value: Array<RateLimitActions>): void;
  addRateLimits(value?: RateLimitActions, index?: number): RateLimitActions;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RateLimitRouteExtension.AsObject;
  static toObject(includeInstance: boolean, msg: RateLimitRouteExtension): RateLimitRouteExtension.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RateLimitRouteExtension, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RateLimitRouteExtension;
  static deserializeBinaryFromReader(message: RateLimitRouteExtension, reader: jspb.BinaryReader): RateLimitRouteExtension;
}

export namespace RateLimitRouteExtension {
  export type AsObject = {
    includeVhRateLimits: boolean,
    rateLimitsList: Array<RateLimitActions.AsObject>,
  }
}

export class Action extends jspb.Message {
  hasSourceCluster(): boolean;
  clearSourceCluster(): void;
  getSourceCluster(): Action.SourceCluster | undefined;
  setSourceCluster(value?: Action.SourceCluster): void;

  hasDestinationCluster(): boolean;
  clearDestinationCluster(): void;
  getDestinationCluster(): Action.DestinationCluster | undefined;
  setDestinationCluster(value?: Action.DestinationCluster): void;

  hasRequestHeaders(): boolean;
  clearRequestHeaders(): void;
  getRequestHeaders(): Action.RequestHeaders | undefined;
  setRequestHeaders(value?: Action.RequestHeaders): void;

  hasRemoteAddress(): boolean;
  clearRemoteAddress(): void;
  getRemoteAddress(): Action.RemoteAddress | undefined;
  setRemoteAddress(value?: Action.RemoteAddress): void;

  hasGenericKey(): boolean;
  clearGenericKey(): void;
  getGenericKey(): Action.GenericKey | undefined;
  setGenericKey(value?: Action.GenericKey): void;

  hasHeaderValueMatch(): boolean;
  clearHeaderValueMatch(): void;
  getHeaderValueMatch(): Action.HeaderValueMatch | undefined;
  setHeaderValueMatch(value?: Action.HeaderValueMatch): void;

  getActionSpecifierCase(): Action.ActionSpecifierCase;
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
    sourceCluster?: Action.SourceCluster.AsObject,
    destinationCluster?: Action.DestinationCluster.AsObject,
    requestHeaders?: Action.RequestHeaders.AsObject,
    remoteAddress?: Action.RemoteAddress.AsObject,
    genericKey?: Action.GenericKey.AsObject,
    headerValueMatch?: Action.HeaderValueMatch.AsObject,
  }

  export class SourceCluster extends jspb.Message {
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): SourceCluster.AsObject;
    static toObject(includeInstance: boolean, msg: SourceCluster): SourceCluster.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: SourceCluster, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): SourceCluster;
    static deserializeBinaryFromReader(message: SourceCluster, reader: jspb.BinaryReader): SourceCluster;
  }

  export namespace SourceCluster {
    export type AsObject = {
    }
  }

  export class DestinationCluster extends jspb.Message {
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): DestinationCluster.AsObject;
    static toObject(includeInstance: boolean, msg: DestinationCluster): DestinationCluster.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: DestinationCluster, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): DestinationCluster;
    static deserializeBinaryFromReader(message: DestinationCluster, reader: jspb.BinaryReader): DestinationCluster;
  }

  export namespace DestinationCluster {
    export type AsObject = {
    }
  }

  export class RequestHeaders extends jspb.Message {
    getHeaderName(): string;
    setHeaderName(value: string): void;

    getDescriptorKey(): string;
    setDescriptorKey(value: string): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): RequestHeaders.AsObject;
    static toObject(includeInstance: boolean, msg: RequestHeaders): RequestHeaders.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: RequestHeaders, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): RequestHeaders;
    static deserializeBinaryFromReader(message: RequestHeaders, reader: jspb.BinaryReader): RequestHeaders;
  }

  export namespace RequestHeaders {
    export type AsObject = {
      headerName: string,
      descriptorKey: string,
    }
  }

  export class RemoteAddress extends jspb.Message {
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): RemoteAddress.AsObject;
    static toObject(includeInstance: boolean, msg: RemoteAddress): RemoteAddress.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: RemoteAddress, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): RemoteAddress;
    static deserializeBinaryFromReader(message: RemoteAddress, reader: jspb.BinaryReader): RemoteAddress;
  }

  export namespace RemoteAddress {
    export type AsObject = {
    }
  }

  export class GenericKey extends jspb.Message {
    getDescriptorValue(): string;
    setDescriptorValue(value: string): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): GenericKey.AsObject;
    static toObject(includeInstance: boolean, msg: GenericKey): GenericKey.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: GenericKey, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): GenericKey;
    static deserializeBinaryFromReader(message: GenericKey, reader: jspb.BinaryReader): GenericKey;
  }

  export namespace GenericKey {
    export type AsObject = {
      descriptorValue: string,
    }
  }

  export class HeaderValueMatch extends jspb.Message {
    getDescriptorValue(): string;
    setDescriptorValue(value: string): void;

    hasExpectMatch(): boolean;
    clearExpectMatch(): void;
    getExpectMatch(): google_protobuf_wrappers_pb.BoolValue | undefined;
    setExpectMatch(value?: google_protobuf_wrappers_pb.BoolValue): void;

    clearHeadersList(): void;
    getHeadersList(): Array<HeaderMatcher>;
    setHeadersList(value: Array<HeaderMatcher>): void;
    addHeaders(value?: HeaderMatcher, index?: number): HeaderMatcher;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): HeaderValueMatch.AsObject;
    static toObject(includeInstance: boolean, msg: HeaderValueMatch): HeaderValueMatch.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: HeaderValueMatch, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): HeaderValueMatch;
    static deserializeBinaryFromReader(message: HeaderValueMatch, reader: jspb.BinaryReader): HeaderValueMatch;
  }

  export namespace HeaderValueMatch {
    export type AsObject = {
      descriptorValue: string,
      expectMatch?: google_protobuf_wrappers_pb.BoolValue.AsObject,
      headersList: Array<HeaderMatcher.AsObject>,
    }
  }

  export enum ActionSpecifierCase {
    ACTION_SPECIFIER_NOT_SET = 0,
    SOURCE_CLUSTER = 1,
    DESTINATION_CLUSTER = 2,
    REQUEST_HEADERS = 3,
    REMOTE_ADDRESS = 4,
    GENERIC_KEY = 5,
    HEADER_VALUE_MATCH = 6,
  }
}

export class Int64Range extends jspb.Message {
  getStart(): number;
  setStart(value: number): void;

  getEnd(): number;
  setEnd(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Int64Range.AsObject;
  static toObject(includeInstance: boolean, msg: Int64Range): Int64Range.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Int64Range, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Int64Range;
  static deserializeBinaryFromReader(message: Int64Range, reader: jspb.BinaryReader): Int64Range;
}

export namespace Int64Range {
  export type AsObject = {
    start: number,
    end: number,
  }
}

export class HeaderMatcher extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  hasExactMatch(): boolean;
  clearExactMatch(): void;
  getExactMatch(): string;
  setExactMatch(value: string): void;

  hasRegexMatch(): boolean;
  clearRegexMatch(): void;
  getRegexMatch(): string;
  setRegexMatch(value: string): void;

  hasRangeMatch(): boolean;
  clearRangeMatch(): void;
  getRangeMatch(): Int64Range | undefined;
  setRangeMatch(value?: Int64Range): void;

  hasPresentMatch(): boolean;
  clearPresentMatch(): void;
  getPresentMatch(): boolean;
  setPresentMatch(value: boolean): void;

  hasPrefixMatch(): boolean;
  clearPrefixMatch(): void;
  getPrefixMatch(): string;
  setPrefixMatch(value: string): void;

  hasSuffixMatch(): boolean;
  clearSuffixMatch(): void;
  getSuffixMatch(): string;
  setSuffixMatch(value: string): void;

  getInvertMatch(): boolean;
  setInvertMatch(value: boolean): void;

  getHeaderMatchSpecifierCase(): HeaderMatcher.HeaderMatchSpecifierCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HeaderMatcher.AsObject;
  static toObject(includeInstance: boolean, msg: HeaderMatcher): HeaderMatcher.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: HeaderMatcher, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HeaderMatcher;
  static deserializeBinaryFromReader(message: HeaderMatcher, reader: jspb.BinaryReader): HeaderMatcher;
}

export namespace HeaderMatcher {
  export type AsObject = {
    name: string,
    exactMatch: string,
    regexMatch: string,
    rangeMatch?: Int64Range.AsObject,
    presentMatch: boolean,
    prefixMatch: string,
    suffixMatch: string,
    invertMatch: boolean,
  }

  export enum HeaderMatchSpecifierCase {
    HEADER_MATCH_SPECIFIER_NOT_SET = 0,
    EXACT_MATCH = 4,
    REGEX_MATCH = 5,
    RANGE_MATCH = 6,
    PRESENT_MATCH = 7,
    PREFIX_MATCH = 9,
    SUFFIX_MATCH = 10,
  }
}

export class QueryParameterMatcher extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getValue(): string;
  setValue(value: string): void;

  hasRegex(): boolean;
  clearRegex(): void;
  getRegex(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setRegex(value?: google_protobuf_wrappers_pb.BoolValue): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): QueryParameterMatcher.AsObject;
  static toObject(includeInstance: boolean, msg: QueryParameterMatcher): QueryParameterMatcher.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: QueryParameterMatcher, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): QueryParameterMatcher;
  static deserializeBinaryFromReader(message: QueryParameterMatcher, reader: jspb.BinaryReader): QueryParameterMatcher;
}

export namespace QueryParameterMatcher {
  export type AsObject = {
    name: string,
    value: string,
    regex?: google_protobuf_wrappers_pb.BoolValue.AsObject,
  }
}

