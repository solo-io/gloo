// package: envoy.config.filter.http.solo_jwt_authn.v2
// file: gloo/projects/gloo/api/external/envoy/extensions/jwt/solo_jwt_authn.proto

import * as jspb from "google-protobuf";
import * as validate_validate_pb from "../../../../../../../../validate/validate_pb";

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

