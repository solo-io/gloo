// package: jwt.options.gloo.solo.io
// file: github.com/solo-io/gloo/projects/gloo/api/v1/enterprise/options/jwt/jwt.proto

import * as jspb from "google-protobuf";
import * as github_com_solo_io_solo_kit_api_v1_ref_pb from "../../../../../../../../../../github.com/solo-io/solo-kit/api/v1/ref_pb";
import * as gogoproto_gogo_pb from "../../../../../../../../../../gogoproto/gogo_pb";
import * as google_protobuf_duration_pb from "google-protobuf/google/protobuf/duration_pb";

export class VhostExtension extends jspb.Message {
  getProvidersMap(): jspb.Map<string, Provider>;
  clearProvidersMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VhostExtension.AsObject;
  static toObject(includeInstance: boolean, msg: VhostExtension): VhostExtension.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: VhostExtension, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VhostExtension;
  static deserializeBinaryFromReader(message: VhostExtension, reader: jspb.BinaryReader): VhostExtension;
}

export namespace VhostExtension {
  export type AsObject = {
    providersMap: Array<[string, Provider.AsObject]>,
  }
}

export class RouteExtension extends jspb.Message {
  getDisable(): boolean;
  setDisable(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RouteExtension.AsObject;
  static toObject(includeInstance: boolean, msg: RouteExtension): RouteExtension.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RouteExtension, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RouteExtension;
  static deserializeBinaryFromReader(message: RouteExtension, reader: jspb.BinaryReader): RouteExtension;
}

export namespace RouteExtension {
  export type AsObject = {
    disable: boolean,
  }
}

export class Provider extends jspb.Message {
  hasJwks(): boolean;
  clearJwks(): void;
  getJwks(): Jwks | undefined;
  setJwks(value?: Jwks): void;

  clearAudiencesList(): void;
  getAudiencesList(): Array<string>;
  setAudiencesList(value: Array<string>): void;
  addAudiences(value: string, index?: number): string;

  getIssuer(): string;
  setIssuer(value: string): void;

  hasTokenSource(): boolean;
  clearTokenSource(): void;
  getTokenSource(): TokenSource | undefined;
  setTokenSource(value?: TokenSource): void;

  getKeepToken(): boolean;
  setKeepToken(value: boolean): void;

  clearClaimsToHeadersList(): void;
  getClaimsToHeadersList(): Array<ClaimToHeader>;
  setClaimsToHeadersList(value: Array<ClaimToHeader>): void;
  addClaimsToHeaders(value?: ClaimToHeader, index?: number): ClaimToHeader;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Provider.AsObject;
  static toObject(includeInstance: boolean, msg: Provider): Provider.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Provider, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Provider;
  static deserializeBinaryFromReader(message: Provider, reader: jspb.BinaryReader): Provider;
}

export namespace Provider {
  export type AsObject = {
    jwks?: Jwks.AsObject,
    audiencesList: Array<string>,
    issuer: string,
    tokenSource?: TokenSource.AsObject,
    keepToken: boolean,
    claimsToHeadersList: Array<ClaimToHeader.AsObject>,
  }
}

export class Jwks extends jspb.Message {
  hasRemote(): boolean;
  clearRemote(): void;
  getRemote(): RemoteJwks | undefined;
  setRemote(value?: RemoteJwks): void;

  hasLocal(): boolean;
  clearLocal(): void;
  getLocal(): LocalJwks | undefined;
  setLocal(value?: LocalJwks): void;

  getJwksCase(): Jwks.JwksCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Jwks.AsObject;
  static toObject(includeInstance: boolean, msg: Jwks): Jwks.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Jwks, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Jwks;
  static deserializeBinaryFromReader(message: Jwks, reader: jspb.BinaryReader): Jwks;
}

export namespace Jwks {
  export type AsObject = {
    remote?: RemoteJwks.AsObject,
    local?: LocalJwks.AsObject,
  }

  export enum JwksCase {
    JWKS_NOT_SET = 0,
    REMOTE = 1,
    LOCAL = 2,
  }
}

export class RemoteJwks extends jspb.Message {
  getUrl(): string;
  setUrl(value: string): void;

  hasUpstreamRef(): boolean;
  clearUpstreamRef(): void;
  getUpstreamRef(): github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef | undefined;
  setUpstreamRef(value?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef): void;

  hasCacheDuration(): boolean;
  clearCacheDuration(): void;
  getCacheDuration(): google_protobuf_duration_pb.Duration | undefined;
  setCacheDuration(value?: google_protobuf_duration_pb.Duration): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoteJwks.AsObject;
  static toObject(includeInstance: boolean, msg: RemoteJwks): RemoteJwks.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RemoteJwks, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoteJwks;
  static deserializeBinaryFromReader(message: RemoteJwks, reader: jspb.BinaryReader): RemoteJwks;
}

export namespace RemoteJwks {
  export type AsObject = {
    url: string,
    upstreamRef?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject,
    cacheDuration?: google_protobuf_duration_pb.Duration.AsObject,
  }
}

export class LocalJwks extends jspb.Message {
  getKey(): string;
  setKey(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LocalJwks.AsObject;
  static toObject(includeInstance: boolean, msg: LocalJwks): LocalJwks.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: LocalJwks, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LocalJwks;
  static deserializeBinaryFromReader(message: LocalJwks, reader: jspb.BinaryReader): LocalJwks;
}

export namespace LocalJwks {
  export type AsObject = {
    key: string,
  }
}

export class TokenSource extends jspb.Message {
  clearHeadersList(): void;
  getHeadersList(): Array<TokenSource.HeaderSource>;
  setHeadersList(value: Array<TokenSource.HeaderSource>): void;
  addHeaders(value?: TokenSource.HeaderSource, index?: number): TokenSource.HeaderSource;

  clearQueryParamsList(): void;
  getQueryParamsList(): Array<string>;
  setQueryParamsList(value: Array<string>): void;
  addQueryParams(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TokenSource.AsObject;
  static toObject(includeInstance: boolean, msg: TokenSource): TokenSource.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: TokenSource, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TokenSource;
  static deserializeBinaryFromReader(message: TokenSource, reader: jspb.BinaryReader): TokenSource;
}

export namespace TokenSource {
  export type AsObject = {
    headersList: Array<TokenSource.HeaderSource.AsObject>,
    queryParamsList: Array<string>,
  }

  export class HeaderSource extends jspb.Message {
    getHeader(): string;
    setHeader(value: string): void;

    getPrefix(): string;
    setPrefix(value: string): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): HeaderSource.AsObject;
    static toObject(includeInstance: boolean, msg: HeaderSource): HeaderSource.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: HeaderSource, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): HeaderSource;
    static deserializeBinaryFromReader(message: HeaderSource, reader: jspb.BinaryReader): HeaderSource;
  }

  export namespace HeaderSource {
    export type AsObject = {
      header: string,
      prefix: string,
    }
  }
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

