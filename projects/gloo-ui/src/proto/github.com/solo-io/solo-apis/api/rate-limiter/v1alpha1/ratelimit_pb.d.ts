/* eslint-disable */
// package: ratelimit.api.solo.io
// file: github.com/solo-io/solo-apis/api/rate-limiter/v1alpha1/ratelimit.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";
import * as extproto_ext_pb from "../../../../../../protoc-gen-ext/extproto/ext_pb";

export class RateLimitConfigSpec extends jspb.Message {
  hasRaw(): boolean;
  clearRaw(): void;
  getRaw(): RateLimitConfigSpec.Raw | undefined;
  setRaw(value?: RateLimitConfigSpec.Raw): void;

  getConfigTypeCase(): RateLimitConfigSpec.ConfigTypeCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RateLimitConfigSpec.AsObject;
  static toObject(includeInstance: boolean, msg: RateLimitConfigSpec): RateLimitConfigSpec.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RateLimitConfigSpec, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RateLimitConfigSpec;
  static deserializeBinaryFromReader(message: RateLimitConfigSpec, reader: jspb.BinaryReader): RateLimitConfigSpec;
}

export namespace RateLimitConfigSpec {
  export type AsObject = {
    raw?: RateLimitConfigSpec.Raw.AsObject,
  }

  export class Raw extends jspb.Message {
    clearDescriptorsList(): void;
    getDescriptorsList(): Array<Descriptor>;
    setDescriptorsList(value: Array<Descriptor>): void;
    addDescriptors(value?: Descriptor, index?: number): Descriptor;

    clearRateLimitsList(): void;
    getRateLimitsList(): Array<RateLimitActions>;
    setRateLimitsList(value: Array<RateLimitActions>): void;
    addRateLimits(value?: RateLimitActions, index?: number): RateLimitActions;

    clearSetDescriptorsList(): void;
    getSetDescriptorsList(): Array<SetDescriptor>;
    setSetDescriptorsList(value: Array<SetDescriptor>): void;
    addSetDescriptors(value?: SetDescriptor, index?: number): SetDescriptor;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Raw.AsObject;
    static toObject(includeInstance: boolean, msg: Raw): Raw.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Raw, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Raw;
    static deserializeBinaryFromReader(message: Raw, reader: jspb.BinaryReader): Raw;
  }

  export namespace Raw {
    export type AsObject = {
      descriptorsList: Array<Descriptor.AsObject>,
      rateLimitsList: Array<RateLimitActions.AsObject>,
      setDescriptorsList: Array<SetDescriptor.AsObject>,
    }
  }

  export enum ConfigTypeCase {
    CONFIG_TYPE_NOT_SET = 0,
    RAW = 1,
  }
}

export class RateLimitConfigStatus extends jspb.Message {
  getState(): RateLimitConfigStatus.StateMap[keyof RateLimitConfigStatus.StateMap];
  setState(value: RateLimitConfigStatus.StateMap[keyof RateLimitConfigStatus.StateMap]): void;

  getMessage(): string;
  setMessage(value: string): void;

  getObservedGeneration(): number;
  setObservedGeneration(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RateLimitConfigStatus.AsObject;
  static toObject(includeInstance: boolean, msg: RateLimitConfigStatus): RateLimitConfigStatus.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RateLimitConfigStatus, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RateLimitConfigStatus;
  static deserializeBinaryFromReader(message: RateLimitConfigStatus, reader: jspb.BinaryReader): RateLimitConfigStatus;
}

export namespace RateLimitConfigStatus {
  export type AsObject = {
    state: RateLimitConfigStatus.StateMap[keyof RateLimitConfigStatus.StateMap],
    message: string,
    observedGeneration: number,
  }

  export interface StateMap {
    PENDING: 0;
    ACCEPTED: 1;
    REJECTED: 2;
  }

  export const State: StateMap;
}

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

  getWeight(): number;
  setWeight(value: number): void;

  getAlwaysApply(): boolean;
  setAlwaysApply(value: boolean): void;

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
    weight: number,
    alwaysApply: boolean,
  }
}

export class SetDescriptor extends jspb.Message {
  clearSimpleDescriptorsList(): void;
  getSimpleDescriptorsList(): Array<SimpleDescriptor>;
  setSimpleDescriptorsList(value: Array<SimpleDescriptor>): void;
  addSimpleDescriptors(value?: SimpleDescriptor, index?: number): SimpleDescriptor;

  hasRateLimit(): boolean;
  clearRateLimit(): void;
  getRateLimit(): RateLimit | undefined;
  setRateLimit(value?: RateLimit): void;

  getAlwaysApply(): boolean;
  setAlwaysApply(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetDescriptor.AsObject;
  static toObject(includeInstance: boolean, msg: SetDescriptor): SetDescriptor.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: SetDescriptor, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetDescriptor;
  static deserializeBinaryFromReader(message: SetDescriptor, reader: jspb.BinaryReader): SetDescriptor;
}

export namespace SetDescriptor {
  export type AsObject = {
    simpleDescriptorsList: Array<SimpleDescriptor.AsObject>,
    rateLimit?: RateLimit.AsObject,
    alwaysApply: boolean,
  }
}

export class SimpleDescriptor extends jspb.Message {
  getKey(): string;
  setKey(value: string): void;

  getValue(): string;
  setValue(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SimpleDescriptor.AsObject;
  static toObject(includeInstance: boolean, msg: SimpleDescriptor): SimpleDescriptor.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: SimpleDescriptor, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SimpleDescriptor;
  static deserializeBinaryFromReader(message: SimpleDescriptor, reader: jspb.BinaryReader): SimpleDescriptor;
}

export namespace SimpleDescriptor {
  export type AsObject = {
    key: string,
    value: string,
  }
}

export class RateLimitActions extends jspb.Message {
  clearActionsList(): void;
  getActionsList(): Array<Action>;
  setActionsList(value: Array<Action>): void;
  addActions(value?: Action, index?: number): Action;

  clearSetActionsList(): void;
  getSetActionsList(): Array<Action>;
  setSetActionsList(value: Array<Action>): void;
  addSetActions(value?: Action, index?: number): Action;

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
    setActionsList: Array<Action.AsObject>,
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

  hasMetadata(): boolean;
  clearMetadata(): void;
  getMetadata(): Action.MetaData | undefined;
  setMetadata(value?: Action.MetaData): void;

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
    metadata?: Action.MetaData.AsObject,
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
    getHeadersList(): Array<Action.HeaderValueMatch.HeaderMatcher>;
    setHeadersList(value: Array<Action.HeaderValueMatch.HeaderMatcher>): void;
    addHeaders(value?: Action.HeaderValueMatch.HeaderMatcher, index?: number): Action.HeaderValueMatch.HeaderMatcher;

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
      headersList: Array<Action.HeaderValueMatch.HeaderMatcher.AsObject>,
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
      getRangeMatch(): Action.HeaderValueMatch.HeaderMatcher.Int64Range | undefined;
      setRangeMatch(value?: Action.HeaderValueMatch.HeaderMatcher.Int64Range): void;

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
        rangeMatch?: Action.HeaderValueMatch.HeaderMatcher.Int64Range.AsObject,
        presentMatch: boolean,
        prefixMatch: string,
        suffixMatch: string,
        invertMatch: boolean,
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
  }

  export class MetaData extends jspb.Message {
    getDescriptorKey(): string;
    setDescriptorKey(value: string): void;

    hasMetadataKey(): boolean;
    clearMetadataKey(): void;
    getMetadataKey(): Action.MetaData.MetadataKey | undefined;
    setMetadataKey(value?: Action.MetaData.MetadataKey): void;

    getDefaultValue(): string;
    setDefaultValue(value: string): void;

    getSource(): Action.MetaData.SourceMap[keyof Action.MetaData.SourceMap];
    setSource(value: Action.MetaData.SourceMap[keyof Action.MetaData.SourceMap]): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): MetaData.AsObject;
    static toObject(includeInstance: boolean, msg: MetaData): MetaData.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: MetaData, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): MetaData;
    static deserializeBinaryFromReader(message: MetaData, reader: jspb.BinaryReader): MetaData;
  }

  export namespace MetaData {
    export type AsObject = {
      descriptorKey: string,
      metadataKey?: Action.MetaData.MetadataKey.AsObject,
      defaultValue: string,
      source: Action.MetaData.SourceMap[keyof Action.MetaData.SourceMap],
    }

    export class MetadataKey extends jspb.Message {
      getKey(): string;
      setKey(value: string): void;

      clearPathList(): void;
      getPathList(): Array<Action.MetaData.MetadataKey.PathSegment>;
      setPathList(value: Array<Action.MetaData.MetadataKey.PathSegment>): void;
      addPath(value?: Action.MetaData.MetadataKey.PathSegment, index?: number): Action.MetaData.MetadataKey.PathSegment;

      serializeBinary(): Uint8Array;
      toObject(includeInstance?: boolean): MetadataKey.AsObject;
      static toObject(includeInstance: boolean, msg: MetadataKey): MetadataKey.AsObject;
      static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
      static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
      static serializeBinaryToWriter(message: MetadataKey, writer: jspb.BinaryWriter): void;
      static deserializeBinary(bytes: Uint8Array): MetadataKey;
      static deserializeBinaryFromReader(message: MetadataKey, reader: jspb.BinaryReader): MetadataKey;
    }

    export namespace MetadataKey {
      export type AsObject = {
        key: string,
        pathList: Array<Action.MetaData.MetadataKey.PathSegment.AsObject>,
      }

      export class PathSegment extends jspb.Message {
        hasKey(): boolean;
        clearKey(): void;
        getKey(): string;
        setKey(value: string): void;

        getSegmentCase(): PathSegment.SegmentCase;
        serializeBinary(): Uint8Array;
        toObject(includeInstance?: boolean): PathSegment.AsObject;
        static toObject(includeInstance: boolean, msg: PathSegment): PathSegment.AsObject;
        static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
        static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
        static serializeBinaryToWriter(message: PathSegment, writer: jspb.BinaryWriter): void;
        static deserializeBinary(bytes: Uint8Array): PathSegment;
        static deserializeBinaryFromReader(message: PathSegment, reader: jspb.BinaryReader): PathSegment;
      }

      export namespace PathSegment {
        export type AsObject = {
          key: string,
        }

        export enum SegmentCase {
          SEGMENT_NOT_SET = 0,
          KEY = 1,
        }
      }
    }

    export interface SourceMap {
      DYNAMIC: 0;
      ROUTE_ENTRY: 1;
    }

    export const Source: SourceMap;
  }

  export enum ActionSpecifierCase {
    ACTION_SPECIFIER_NOT_SET = 0,
    SOURCE_CLUSTER = 1,
    DESTINATION_CLUSTER = 2,
    REQUEST_HEADERS = 3,
    REMOTE_ADDRESS = 4,
    GENERIC_KEY = 5,
    HEADER_VALUE_MATCH = 6,
    METADATA = 8,
  }
}
