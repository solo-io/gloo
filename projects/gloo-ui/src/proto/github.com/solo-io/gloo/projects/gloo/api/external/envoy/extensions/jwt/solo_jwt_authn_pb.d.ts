/* eslint-disable */
// package: envoy.config.filter.http.solo_jwt_authn.v2
// file: github.com/solo-io/gloo/projects/gloo/api/external/envoy/extensions/jwt/solo_jwt_authn.proto

import * as jspb from "google-protobuf";
import * as validate_validate_pb from "../../../../../../../../../../validate/validate_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_filters_http_jwt_authn_v3_config_pb from "../../../../../../../../../../github.com/solo-io/gloo/projects/gloo/api/external/envoy/extensions/filters/http/jwt_authn/v3/config_pb";

export class JwtWithStage extends jspb.Message {
  hasJwtAuthn(): boolean;
  clearJwtAuthn(): void;
  getJwtAuthn(): github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_filters_http_jwt_authn_v3_config_pb.JwtAuthentication | undefined;
  setJwtAuthn(value?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_filters_http_jwt_authn_v3_config_pb.JwtAuthentication): void;

  getStage(): number;
  setStage(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): JwtWithStage.AsObject;
  static toObject(includeInstance: boolean, msg: JwtWithStage): JwtWithStage.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: JwtWithStage, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): JwtWithStage;
  static deserializeBinaryFromReader(message: JwtWithStage, reader: jspb.BinaryReader): JwtWithStage;
}

export namespace JwtWithStage {
  export type AsObject = {
    jwtAuthn?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_filters_http_jwt_authn_v3_config_pb.JwtAuthentication.AsObject,
    stage: number,
  }
}

export class SoloJwtAuthnPerRoute extends jspb.Message {
  getRequirement(): string;
  setRequirement(value: string): void;

  getClaimsToHeadersMap(): jspb.Map<string, SoloJwtAuthnPerRoute.ClaimToHeaders>;
  clearClaimsToHeadersMap(): void;
  getClearRouteCache(): boolean;
  setClearRouteCache(value: boolean): void;

  getPayloadInMetadata(): string;
  setPayloadInMetadata(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SoloJwtAuthnPerRoute.AsObject;
  static toObject(includeInstance: boolean, msg: SoloJwtAuthnPerRoute): SoloJwtAuthnPerRoute.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: SoloJwtAuthnPerRoute, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SoloJwtAuthnPerRoute;
  static deserializeBinaryFromReader(message: SoloJwtAuthnPerRoute, reader: jspb.BinaryReader): SoloJwtAuthnPerRoute;
}

export namespace SoloJwtAuthnPerRoute {
  export type AsObject = {
    requirement: string,
    claimsToHeadersMap: Array<[string, SoloJwtAuthnPerRoute.ClaimToHeaders.AsObject]>,
    clearRouteCache: boolean,
    payloadInMetadata: string,
  }

  export class ClaimToHeader extends jspb.Message {
    getClaim(): string;
    setClaim(value: string): void;

    getHeader(): string;
    setHeader(value: string): void;

    getAppend(): boolean;
    setAppend(value: boolean): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ClaimToHeader.AsObject;
    static toObject(includeInstance: boolean, msg: ClaimToHeader): ClaimToHeader.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: ClaimToHeader, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ClaimToHeader;
    static deserializeBinaryFromReader(message: ClaimToHeader, reader: jspb.BinaryReader): ClaimToHeader;
  }

  export namespace ClaimToHeader {
    export type AsObject = {
      claim: string,
      header: string,
      append: boolean,
    }
  }

  export class ClaimToHeaders extends jspb.Message {
    clearClaimsList(): void;
    getClaimsList(): Array<SoloJwtAuthnPerRoute.ClaimToHeader>;
    setClaimsList(value: Array<SoloJwtAuthnPerRoute.ClaimToHeader>): void;
    addClaims(value?: SoloJwtAuthnPerRoute.ClaimToHeader, index?: number): SoloJwtAuthnPerRoute.ClaimToHeader;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ClaimToHeaders.AsObject;
    static toObject(includeInstance: boolean, msg: ClaimToHeaders): ClaimToHeaders.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: ClaimToHeaders, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ClaimToHeaders;
    static deserializeBinaryFromReader(message: ClaimToHeaders, reader: jspb.BinaryReader): ClaimToHeaders;
  }

  export namespace ClaimToHeaders {
    export type AsObject = {
      claimsList: Array<SoloJwtAuthnPerRoute.ClaimToHeader.AsObject>,
    }
  }
}

export class StagedJwtAuthnPerRoute extends jspb.Message {
  getJwtConfigsMap(): jspb.Map<number, SoloJwtAuthnPerRoute>;
  clearJwtConfigsMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StagedJwtAuthnPerRoute.AsObject;
  static toObject(includeInstance: boolean, msg: StagedJwtAuthnPerRoute): StagedJwtAuthnPerRoute.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: StagedJwtAuthnPerRoute, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StagedJwtAuthnPerRoute;
  static deserializeBinaryFromReader(message: StagedJwtAuthnPerRoute, reader: jspb.BinaryReader): StagedJwtAuthnPerRoute;
}

export namespace StagedJwtAuthnPerRoute {
  export type AsObject = {
    jwtConfigsMap: Array<[number, SoloJwtAuthnPerRoute.AsObject]>,
  }
}
