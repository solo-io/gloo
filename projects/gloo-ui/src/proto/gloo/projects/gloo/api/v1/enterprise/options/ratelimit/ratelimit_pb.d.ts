/* eslint-disable */
// package: ratelimit.options.gloo.solo.io
// file: gloo/projects/gloo/api/v1/enterprise/options/ratelimit/ratelimit.proto

import * as jspb from "google-protobuf";
import * as solo_apis_api_rate_limiter_v1alpha1_ratelimit_pb from "../../../../../../../../solo-apis/api/rate-limiter/v1alpha1/ratelimit_pb";
import * as solo_kit_api_v1_ref_pb from "../../../../../../../../solo-kit/api/v1/ref_pb";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";
import * as google_protobuf_duration_pb from "google-protobuf/google/protobuf/duration_pb";
import * as gogoproto_gogo_pb from "../../../../../../../../gogoproto/gogo_pb";
import * as extproto_ext_pb from "../../../../../../../../protoc-gen-ext/extproto/ext_pb";

export class IngressRateLimit extends jspb.Message {
  hasAuthorizedLimits(): boolean;
  clearAuthorizedLimits(): void;
  getAuthorizedLimits(): solo_apis_api_rate_limiter_v1alpha1_ratelimit_pb.RateLimit | undefined;
  setAuthorizedLimits(value?: solo_apis_api_rate_limiter_v1alpha1_ratelimit_pb.RateLimit): void;

  hasAnonymousLimits(): boolean;
  clearAnonymousLimits(): void;
  getAnonymousLimits(): solo_apis_api_rate_limiter_v1alpha1_ratelimit_pb.RateLimit | undefined;
  setAnonymousLimits(value?: solo_apis_api_rate_limiter_v1alpha1_ratelimit_pb.RateLimit): void;

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
    authorizedLimits?: solo_apis_api_rate_limiter_v1alpha1_ratelimit_pb.RateLimit.AsObject,
    anonymousLimits?: solo_apis_api_rate_limiter_v1alpha1_ratelimit_pb.RateLimit.AsObject,
  }
}

export class Settings extends jspb.Message {
  hasRatelimitServerRef(): boolean;
  clearRatelimitServerRef(): void;
  getRatelimitServerRef(): solo_kit_api_v1_ref_pb.ResourceRef | undefined;
  setRatelimitServerRef(value?: solo_kit_api_v1_ref_pb.ResourceRef): void;

  hasRequestTimeout(): boolean;
  clearRequestTimeout(): void;
  getRequestTimeout(): google_protobuf_duration_pb.Duration | undefined;
  setRequestTimeout(value?: google_protobuf_duration_pb.Duration): void;

  getDenyOnFail(): boolean;
  setDenyOnFail(value: boolean): void;

  getRateLimitBeforeAuth(): boolean;
  setRateLimitBeforeAuth(value: boolean): void;

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
    ratelimitServerRef?: solo_kit_api_v1_ref_pb.ResourceRef.AsObject,
    requestTimeout?: google_protobuf_duration_pb.Duration.AsObject,
    denyOnFail: boolean,
    rateLimitBeforeAuth: boolean,
  }
}

export class ServiceSettings extends jspb.Message {
  clearDescriptorsList(): void;
  getDescriptorsList(): Array<solo_apis_api_rate_limiter_v1alpha1_ratelimit_pb.Descriptor>;
  setDescriptorsList(value: Array<solo_apis_api_rate_limiter_v1alpha1_ratelimit_pb.Descriptor>): void;
  addDescriptors(value?: solo_apis_api_rate_limiter_v1alpha1_ratelimit_pb.Descriptor, index?: number): solo_apis_api_rate_limiter_v1alpha1_ratelimit_pb.Descriptor;

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
    descriptorsList: Array<solo_apis_api_rate_limiter_v1alpha1_ratelimit_pb.Descriptor.AsObject>,
  }
}

export class RateLimitConfigRefs extends jspb.Message {
  clearRefsList(): void;
  getRefsList(): Array<RateLimitConfigRef>;
  setRefsList(value: Array<RateLimitConfigRef>): void;
  addRefs(value?: RateLimitConfigRef, index?: number): RateLimitConfigRef;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RateLimitConfigRefs.AsObject;
  static toObject(includeInstance: boolean, msg: RateLimitConfigRefs): RateLimitConfigRefs.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RateLimitConfigRefs, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RateLimitConfigRefs;
  static deserializeBinaryFromReader(message: RateLimitConfigRefs, reader: jspb.BinaryReader): RateLimitConfigRefs;
}

export namespace RateLimitConfigRefs {
  export type AsObject = {
    refsList: Array<RateLimitConfigRef.AsObject>,
  }
}

export class RateLimitConfigRef extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getNamespace(): string;
  setNamespace(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RateLimitConfigRef.AsObject;
  static toObject(includeInstance: boolean, msg: RateLimitConfigRef): RateLimitConfigRef.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RateLimitConfigRef, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RateLimitConfigRef;
  static deserializeBinaryFromReader(message: RateLimitConfigRef, reader: jspb.BinaryReader): RateLimitConfigRef;
}

export namespace RateLimitConfigRef {
  export type AsObject = {
    name: string,
    namespace: string,
  }
}

export class RateLimitVhostExtension extends jspb.Message {
  clearRateLimitsList(): void;
  getRateLimitsList(): Array<solo_apis_api_rate_limiter_v1alpha1_ratelimit_pb.RateLimitActions>;
  setRateLimitsList(value: Array<solo_apis_api_rate_limiter_v1alpha1_ratelimit_pb.RateLimitActions>): void;
  addRateLimits(value?: solo_apis_api_rate_limiter_v1alpha1_ratelimit_pb.RateLimitActions, index?: number): solo_apis_api_rate_limiter_v1alpha1_ratelimit_pb.RateLimitActions;

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
    rateLimitsList: Array<solo_apis_api_rate_limiter_v1alpha1_ratelimit_pb.RateLimitActions.AsObject>,
  }
}

export class RateLimitRouteExtension extends jspb.Message {
  getIncludeVhRateLimits(): boolean;
  setIncludeVhRateLimits(value: boolean): void;

  clearRateLimitsList(): void;
  getRateLimitsList(): Array<solo_apis_api_rate_limiter_v1alpha1_ratelimit_pb.RateLimitActions>;
  setRateLimitsList(value: Array<solo_apis_api_rate_limiter_v1alpha1_ratelimit_pb.RateLimitActions>): void;
  addRateLimits(value?: solo_apis_api_rate_limiter_v1alpha1_ratelimit_pb.RateLimitActions, index?: number): solo_apis_api_rate_limiter_v1alpha1_ratelimit_pb.RateLimitActions;

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
    rateLimitsList: Array<solo_apis_api_rate_limiter_v1alpha1_ratelimit_pb.RateLimitActions.AsObject>,
  }
}
