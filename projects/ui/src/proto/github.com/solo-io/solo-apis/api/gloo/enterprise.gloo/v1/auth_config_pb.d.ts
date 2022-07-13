/* eslint-disable */
// package: enterprise.gloo.solo.io
// file: github.com/solo-io/solo-apis/api/gloo/enterprise.gloo/v1/auth_config.proto

import * as jspb from "google-protobuf";
import * as github_com_solo_io_solo_kit_api_v1_ref_pb from "../../../../../../../github.com/solo-io/solo-kit/api/v1/ref_pb";
import * as extproto_ext_pb from "../../../../../../../extproto/ext_pb";
import * as github_com_solo_io_solo_kit_api_v1_solo_kit_pb from "../../../../../../../github.com/solo-io/solo-kit/api/v1/solo-kit_pb";
import * as github_com_solo_io_solo_kit_api_external_envoy_api_v2_discovery_pb from "../../../../../../../github.com/solo-io/solo-kit/api/external/envoy/api/v2/discovery_pb";
import * as google_api_annotations_pb from "../../../../../../../google/api/annotations_pb";
import * as google_protobuf_duration_pb from "google-protobuf/google/protobuf/duration_pb";
import * as google_protobuf_struct_pb from "google-protobuf/google/protobuf/struct_pb";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";
import * as google_protobuf_empty_pb from "google-protobuf/google/protobuf/empty_pb";

export class AuthConfigSpec extends jspb.Message {
  clearConfigsList(): void;
  getConfigsList(): Array<AuthConfigSpec.Config>;
  setConfigsList(value: Array<AuthConfigSpec.Config>): void;
  addConfigs(value?: AuthConfigSpec.Config, index?: number): AuthConfigSpec.Config;

  hasBooleanExpr(): boolean;
  clearBooleanExpr(): void;
  getBooleanExpr(): google_protobuf_wrappers_pb.StringValue | undefined;
  setBooleanExpr(value?: google_protobuf_wrappers_pb.StringValue): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AuthConfigSpec.AsObject;
  static toObject(includeInstance: boolean, msg: AuthConfigSpec): AuthConfigSpec.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: AuthConfigSpec, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AuthConfigSpec;
  static deserializeBinaryFromReader(message: AuthConfigSpec, reader: jspb.BinaryReader): AuthConfigSpec;
}

export namespace AuthConfigSpec {
  export type AsObject = {
    configsList: Array<AuthConfigSpec.Config.AsObject>,
    booleanExpr?: google_protobuf_wrappers_pb.StringValue.AsObject,
  }

  export class Config extends jspb.Message {
    hasName(): boolean;
    clearName(): void;
    getName(): google_protobuf_wrappers_pb.StringValue | undefined;
    setName(value?: google_protobuf_wrappers_pb.StringValue): void;

    hasBasicAuth(): boolean;
    clearBasicAuth(): void;
    getBasicAuth(): BasicAuth | undefined;
    setBasicAuth(value?: BasicAuth): void;

    hasOauth(): boolean;
    clearOauth(): void;
    getOauth(): OAuth | undefined;
    setOauth(value?: OAuth): void;

    hasOauth2(): boolean;
    clearOauth2(): void;
    getOauth2(): OAuth2 | undefined;
    setOauth2(value?: OAuth2): void;

    hasApiKeyAuth(): boolean;
    clearApiKeyAuth(): void;
    getApiKeyAuth(): ApiKeyAuth | undefined;
    setApiKeyAuth(value?: ApiKeyAuth): void;

    hasPluginAuth(): boolean;
    clearPluginAuth(): void;
    getPluginAuth(): AuthPlugin | undefined;
    setPluginAuth(value?: AuthPlugin): void;

    hasOpaAuth(): boolean;
    clearOpaAuth(): void;
    getOpaAuth(): OpaAuth | undefined;
    setOpaAuth(value?: OpaAuth): void;

    hasLdap(): boolean;
    clearLdap(): void;
    getLdap(): Ldap | undefined;
    setLdap(value?: Ldap): void;

    hasJwt(): boolean;
    clearJwt(): void;
    getJwt(): google_protobuf_empty_pb.Empty | undefined;
    setJwt(value?: google_protobuf_empty_pb.Empty): void;

    hasPassThroughAuth(): boolean;
    clearPassThroughAuth(): void;
    getPassThroughAuth(): PassThroughAuth | undefined;
    setPassThroughAuth(value?: PassThroughAuth): void;

    getAuthConfigCase(): Config.AuthConfigCase;
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Config.AsObject;
    static toObject(includeInstance: boolean, msg: Config): Config.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Config, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Config;
    static deserializeBinaryFromReader(message: Config, reader: jspb.BinaryReader): Config;
  }

  export namespace Config {
    export type AsObject = {
      name?: google_protobuf_wrappers_pb.StringValue.AsObject,
      basicAuth?: BasicAuth.AsObject,
      oauth?: OAuth.AsObject,
      oauth2?: OAuth2.AsObject,
      apiKeyAuth?: ApiKeyAuth.AsObject,
      pluginAuth?: AuthPlugin.AsObject,
      opaAuth?: OpaAuth.AsObject,
      ldap?: Ldap.AsObject,
      jwt?: google_protobuf_empty_pb.Empty.AsObject,
      passThroughAuth?: PassThroughAuth.AsObject,
    }

    export enum AuthConfigCase {
      AUTH_CONFIG_NOT_SET = 0,
      BASIC_AUTH = 1,
      OAUTH = 2,
      OAUTH2 = 8,
      API_KEY_AUTH = 4,
      PLUGIN_AUTH = 5,
      OPA_AUTH = 6,
      LDAP = 7,
      JWT = 11,
      PASS_THROUGH_AUTH = 12,
    }
  }
}

export class ExtAuthExtension extends jspb.Message {
  hasDisable(): boolean;
  clearDisable(): void;
  getDisable(): boolean;
  setDisable(value: boolean): void;

  hasConfigRef(): boolean;
  clearConfigRef(): void;
  getConfigRef(): github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef | undefined;
  setConfigRef(value?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef): void;

  hasCustomAuth(): boolean;
  clearCustomAuth(): void;
  getCustomAuth(): CustomAuth | undefined;
  setCustomAuth(value?: CustomAuth): void;

  getSpecCase(): ExtAuthExtension.SpecCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ExtAuthExtension.AsObject;
  static toObject(includeInstance: boolean, msg: ExtAuthExtension): ExtAuthExtension.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ExtAuthExtension, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ExtAuthExtension;
  static deserializeBinaryFromReader(message: ExtAuthExtension, reader: jspb.BinaryReader): ExtAuthExtension;
}

export namespace ExtAuthExtension {
  export type AsObject = {
    disable: boolean,
    configRef?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject,
    customAuth?: CustomAuth.AsObject,
  }

  export enum SpecCase {
    SPEC_NOT_SET = 0,
    DISABLE = 1,
    CONFIG_REF = 2,
    CUSTOM_AUTH = 3,
  }
}

export class Settings extends jspb.Message {
  hasExtauthzServerRef(): boolean;
  clearExtauthzServerRef(): void;
  getExtauthzServerRef(): github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef | undefined;
  setExtauthzServerRef(value?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef): void;

  hasHttpService(): boolean;
  clearHttpService(): void;
  getHttpService(): HttpService | undefined;
  setHttpService(value?: HttpService): void;

  hasGrpcService(): boolean;
  clearGrpcService(): void;
  getGrpcService(): GrpcService | undefined;
  setGrpcService(value?: GrpcService): void;

  getUserIdHeader(): string;
  setUserIdHeader(value: string): void;

  hasRequestTimeout(): boolean;
  clearRequestTimeout(): void;
  getRequestTimeout(): google_protobuf_duration_pb.Duration | undefined;
  setRequestTimeout(value?: google_protobuf_duration_pb.Duration): void;

  getFailureModeAllow(): boolean;
  setFailureModeAllow(value: boolean): void;

  hasRequestBody(): boolean;
  clearRequestBody(): void;
  getRequestBody(): BufferSettings | undefined;
  setRequestBody(value?: BufferSettings): void;

  getClearRouteCache(): boolean;
  setClearRouteCache(value: boolean): void;

  getStatusOnError(): number;
  setStatusOnError(value: number): void;

  getTransportApiVersion(): Settings.ApiVersionMap[keyof Settings.ApiVersionMap];
  setTransportApiVersion(value: Settings.ApiVersionMap[keyof Settings.ApiVersionMap]): void;

  getStatPrefix(): string;
  setStatPrefix(value: string): void;

  getServiceTypeCase(): Settings.ServiceTypeCase;
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
    extauthzServerRef?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject,
    httpService?: HttpService.AsObject,
    grpcService?: GrpcService.AsObject,
    userIdHeader: string,
    requestTimeout?: google_protobuf_duration_pb.Duration.AsObject,
    failureModeAllow: boolean,
    requestBody?: BufferSettings.AsObject,
    clearRouteCache: boolean,
    statusOnError: number,
    transportApiVersion: Settings.ApiVersionMap[keyof Settings.ApiVersionMap],
    statPrefix: string,
  }

  export interface ApiVersionMap {
    V3: 0;
  }

  export const ApiVersion: ApiVersionMap;

  export enum ServiceTypeCase {
    SERVICE_TYPE_NOT_SET = 0,
    HTTP_SERVICE = 2,
    GRPC_SERVICE = 11,
  }
}

export class GrpcService extends jspb.Message {
  getAuthority(): string;
  setAuthority(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GrpcService.AsObject;
  static toObject(includeInstance: boolean, msg: GrpcService): GrpcService.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GrpcService, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GrpcService;
  static deserializeBinaryFromReader(message: GrpcService, reader: jspb.BinaryReader): GrpcService;
}

export namespace GrpcService {
  export type AsObject = {
    authority: string,
  }
}

export class HttpService extends jspb.Message {
  getPathPrefix(): string;
  setPathPrefix(value: string): void;

  hasRequest(): boolean;
  clearRequest(): void;
  getRequest(): HttpService.Request | undefined;
  setRequest(value?: HttpService.Request): void;

  hasResponse(): boolean;
  clearResponse(): void;
  getResponse(): HttpService.Response | undefined;
  setResponse(value?: HttpService.Response): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HttpService.AsObject;
  static toObject(includeInstance: boolean, msg: HttpService): HttpService.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: HttpService, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HttpService;
  static deserializeBinaryFromReader(message: HttpService, reader: jspb.BinaryReader): HttpService;
}

export namespace HttpService {
  export type AsObject = {
    pathPrefix: string,
    request?: HttpService.Request.AsObject,
    response?: HttpService.Response.AsObject,
  }

  export class Request extends jspb.Message {
    clearAllowedHeadersList(): void;
    getAllowedHeadersList(): Array<string>;
    setAllowedHeadersList(value: Array<string>): void;
    addAllowedHeaders(value: string, index?: number): string;

    getHeadersToAddMap(): jspb.Map<string, string>;
    clearHeadersToAddMap(): void;
    clearAllowedHeadersRegexList(): void;
    getAllowedHeadersRegexList(): Array<string>;
    setAllowedHeadersRegexList(value: Array<string>): void;
    addAllowedHeadersRegex(value: string, index?: number): string;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Request.AsObject;
    static toObject(includeInstance: boolean, msg: Request): Request.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Request, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Request;
    static deserializeBinaryFromReader(message: Request, reader: jspb.BinaryReader): Request;
  }

  export namespace Request {
    export type AsObject = {
      allowedHeadersList: Array<string>,
      headersToAddMap: Array<[string, string]>,
      allowedHeadersRegexList: Array<string>,
    }
  }

  export class Response extends jspb.Message {
    clearAllowedUpstreamHeadersList(): void;
    getAllowedUpstreamHeadersList(): Array<string>;
    setAllowedUpstreamHeadersList(value: Array<string>): void;
    addAllowedUpstreamHeaders(value: string, index?: number): string;

    clearAllowedClientHeadersList(): void;
    getAllowedClientHeadersList(): Array<string>;
    setAllowedClientHeadersList(value: Array<string>): void;
    addAllowedClientHeaders(value: string, index?: number): string;

    clearAllowedUpstreamHeadersToAppendList(): void;
    getAllowedUpstreamHeadersToAppendList(): Array<string>;
    setAllowedUpstreamHeadersToAppendList(value: Array<string>): void;
    addAllowedUpstreamHeadersToAppend(value: string, index?: number): string;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Response.AsObject;
    static toObject(includeInstance: boolean, msg: Response): Response.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Response, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Response;
    static deserializeBinaryFromReader(message: Response, reader: jspb.BinaryReader): Response;
  }

  export namespace Response {
    export type AsObject = {
      allowedUpstreamHeadersList: Array<string>,
      allowedClientHeadersList: Array<string>,
      allowedUpstreamHeadersToAppendList: Array<string>,
    }
  }
}

export class BufferSettings extends jspb.Message {
  getMaxRequestBytes(): number;
  setMaxRequestBytes(value: number): void;

  getAllowPartialMessage(): boolean;
  setAllowPartialMessage(value: boolean): void;

  getPackAsBytes(): boolean;
  setPackAsBytes(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): BufferSettings.AsObject;
  static toObject(includeInstance: boolean, msg: BufferSettings): BufferSettings.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: BufferSettings, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): BufferSettings;
  static deserializeBinaryFromReader(message: BufferSettings, reader: jspb.BinaryReader): BufferSettings;
}

export namespace BufferSettings {
  export type AsObject = {
    maxRequestBytes: number,
    allowPartialMessage: boolean,
    packAsBytes: boolean,
  }
}

export class CustomAuth extends jspb.Message {
  getContextExtensionsMap(): jspb.Map<string, string>;
  clearContextExtensionsMap(): void;
  getName(): string;
  setName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CustomAuth.AsObject;
  static toObject(includeInstance: boolean, msg: CustomAuth): CustomAuth.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CustomAuth, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CustomAuth;
  static deserializeBinaryFromReader(message: CustomAuth, reader: jspb.BinaryReader): CustomAuth;
}

export namespace CustomAuth {
  export type AsObject = {
    contextExtensionsMap: Array<[string, string]>,
    name: string,
  }
}

export class AuthPlugin extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getPluginFileName(): string;
  setPluginFileName(value: string): void;

  getExportedSymbolName(): string;
  setExportedSymbolName(value: string): void;

  hasConfig(): boolean;
  clearConfig(): void;
  getConfig(): google_protobuf_struct_pb.Struct | undefined;
  setConfig(value?: google_protobuf_struct_pb.Struct): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AuthPlugin.AsObject;
  static toObject(includeInstance: boolean, msg: AuthPlugin): AuthPlugin.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: AuthPlugin, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AuthPlugin;
  static deserializeBinaryFromReader(message: AuthPlugin, reader: jspb.BinaryReader): AuthPlugin;
}

export namespace AuthPlugin {
  export type AsObject = {
    name: string,
    pluginFileName: string,
    exportedSymbolName: string,
    config?: google_protobuf_struct_pb.Struct.AsObject,
  }
}

export class BasicAuth extends jspb.Message {
  getRealm(): string;
  setRealm(value: string): void;

  hasApr(): boolean;
  clearApr(): void;
  getApr(): BasicAuth.Apr | undefined;
  setApr(value?: BasicAuth.Apr): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): BasicAuth.AsObject;
  static toObject(includeInstance: boolean, msg: BasicAuth): BasicAuth.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: BasicAuth, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): BasicAuth;
  static deserializeBinaryFromReader(message: BasicAuth, reader: jspb.BinaryReader): BasicAuth;
}

export namespace BasicAuth {
  export type AsObject = {
    realm: string,
    apr?: BasicAuth.Apr.AsObject,
  }

  export class Apr extends jspb.Message {
    getUsersMap(): jspb.Map<string, BasicAuth.Apr.SaltedHashedPassword>;
    clearUsersMap(): void;
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Apr.AsObject;
    static toObject(includeInstance: boolean, msg: Apr): Apr.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Apr, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Apr;
    static deserializeBinaryFromReader(message: Apr, reader: jspb.BinaryReader): Apr;
  }

  export namespace Apr {
    export type AsObject = {
      usersMap: Array<[string, BasicAuth.Apr.SaltedHashedPassword.AsObject]>,
    }

    export class SaltedHashedPassword extends jspb.Message {
      getSalt(): string;
      setSalt(value: string): void;

      getHashedPassword(): string;
      setHashedPassword(value: string): void;

      serializeBinary(): Uint8Array;
      toObject(includeInstance?: boolean): SaltedHashedPassword.AsObject;
      static toObject(includeInstance: boolean, msg: SaltedHashedPassword): SaltedHashedPassword.AsObject;
      static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
      static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
      static serializeBinaryToWriter(message: SaltedHashedPassword, writer: jspb.BinaryWriter): void;
      static deserializeBinary(bytes: Uint8Array): SaltedHashedPassword;
      static deserializeBinaryFromReader(message: SaltedHashedPassword, reader: jspb.BinaryReader): SaltedHashedPassword;
    }

    export namespace SaltedHashedPassword {
      export type AsObject = {
        salt: string,
        hashedPassword: string,
      }
    }
  }
}

export class OAuth extends jspb.Message {
  getClientId(): string;
  setClientId(value: string): void;

  hasClientSecretRef(): boolean;
  clearClientSecretRef(): void;
  getClientSecretRef(): github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef | undefined;
  setClientSecretRef(value?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef): void;

  getIssuerUrl(): string;
  setIssuerUrl(value: string): void;

  getAuthEndpointQueryParamsMap(): jspb.Map<string, string>;
  clearAuthEndpointQueryParamsMap(): void;
  getAppUrl(): string;
  setAppUrl(value: string): void;

  getCallbackPath(): string;
  setCallbackPath(value: string): void;

  clearScopesList(): void;
  getScopesList(): Array<string>;
  setScopesList(value: Array<string>): void;
  addScopes(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OAuth.AsObject;
  static toObject(includeInstance: boolean, msg: OAuth): OAuth.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: OAuth, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OAuth;
  static deserializeBinaryFromReader(message: OAuth, reader: jspb.BinaryReader): OAuth;
}

export namespace OAuth {
  export type AsObject = {
    clientId: string,
    clientSecretRef?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject,
    issuerUrl: string,
    authEndpointQueryParamsMap: Array<[string, string]>,
    appUrl: string,
    callbackPath: string,
    scopesList: Array<string>,
  }
}

export class OAuth2 extends jspb.Message {
  hasOidcAuthorizationCode(): boolean;
  clearOidcAuthorizationCode(): void;
  getOidcAuthorizationCode(): OidcAuthorizationCode | undefined;
  setOidcAuthorizationCode(value?: OidcAuthorizationCode): void;

  hasAccessTokenValidation(): boolean;
  clearAccessTokenValidation(): void;
  getAccessTokenValidation(): AccessTokenValidation | undefined;
  setAccessTokenValidation(value?: AccessTokenValidation): void;

  getOauthTypeCase(): OAuth2.OauthTypeCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OAuth2.AsObject;
  static toObject(includeInstance: boolean, msg: OAuth2): OAuth2.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: OAuth2, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OAuth2;
  static deserializeBinaryFromReader(message: OAuth2, reader: jspb.BinaryReader): OAuth2;
}

export namespace OAuth2 {
  export type AsObject = {
    oidcAuthorizationCode?: OidcAuthorizationCode.AsObject,
    accessTokenValidation?: AccessTokenValidation.AsObject,
  }

  export enum OauthTypeCase {
    OAUTH_TYPE_NOT_SET = 0,
    OIDC_AUTHORIZATION_CODE = 1,
    ACCESS_TOKEN_VALIDATION = 2,
  }
}

export class RedisOptions extends jspb.Message {
  getHost(): string;
  setHost(value: string): void;

  getDb(): number;
  setDb(value: number): void;

  getPoolSize(): number;
  setPoolSize(value: number): void;

  getTlsCertMountPath(): string;
  setTlsCertMountPath(value: string): void;

  getSocketType(): RedisOptions.SocketTypeMap[keyof RedisOptions.SocketTypeMap];
  setSocketType(value: RedisOptions.SocketTypeMap[keyof RedisOptions.SocketTypeMap]): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RedisOptions.AsObject;
  static toObject(includeInstance: boolean, msg: RedisOptions): RedisOptions.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RedisOptions, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RedisOptions;
  static deserializeBinaryFromReader(message: RedisOptions, reader: jspb.BinaryReader): RedisOptions;
}

export namespace RedisOptions {
  export type AsObject = {
    host: string,
    db: number,
    poolSize: number,
    tlsCertMountPath: string,
    socketType: RedisOptions.SocketTypeMap[keyof RedisOptions.SocketTypeMap],
  }

  export interface SocketTypeMap {
    TCP: 0;
    TLS: 1;
  }

  export const SocketType: SocketTypeMap;
}

export class UserSession extends jspb.Message {
  getFailOnFetchFailure(): boolean;
  setFailOnFetchFailure(value: boolean): void;

  hasCookieOptions(): boolean;
  clearCookieOptions(): void;
  getCookieOptions(): UserSession.CookieOptions | undefined;
  setCookieOptions(value?: UserSession.CookieOptions): void;

  hasCookie(): boolean;
  clearCookie(): void;
  getCookie(): UserSession.InternalSession | undefined;
  setCookie(value?: UserSession.InternalSession): void;

  hasRedis(): boolean;
  clearRedis(): void;
  getRedis(): UserSession.RedisSession | undefined;
  setRedis(value?: UserSession.RedisSession): void;

  getSessionCase(): UserSession.SessionCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserSession.AsObject;
  static toObject(includeInstance: boolean, msg: UserSession): UserSession.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UserSession, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserSession;
  static deserializeBinaryFromReader(message: UserSession, reader: jspb.BinaryReader): UserSession;
}

export namespace UserSession {
  export type AsObject = {
    failOnFetchFailure: boolean,
    cookieOptions?: UserSession.CookieOptions.AsObject,
    cookie?: UserSession.InternalSession.AsObject,
    redis?: UserSession.RedisSession.AsObject,
  }

  export class InternalSession extends jspb.Message {
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): InternalSession.AsObject;
    static toObject(includeInstance: boolean, msg: InternalSession): InternalSession.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: InternalSession, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): InternalSession;
    static deserializeBinaryFromReader(message: InternalSession, reader: jspb.BinaryReader): InternalSession;
  }

  export namespace InternalSession {
    export type AsObject = {
    }
  }

  export class RedisSession extends jspb.Message {
    hasOptions(): boolean;
    clearOptions(): void;
    getOptions(): RedisOptions | undefined;
    setOptions(value?: RedisOptions): void;

    getKeyPrefix(): string;
    setKeyPrefix(value: string): void;

    getCookieName(): string;
    setCookieName(value: string): void;

    hasAllowRefreshing(): boolean;
    clearAllowRefreshing(): void;
    getAllowRefreshing(): google_protobuf_wrappers_pb.BoolValue | undefined;
    setAllowRefreshing(value?: google_protobuf_wrappers_pb.BoolValue): void;

    hasPreExpiryBuffer(): boolean;
    clearPreExpiryBuffer(): void;
    getPreExpiryBuffer(): google_protobuf_duration_pb.Duration | undefined;
    setPreExpiryBuffer(value?: google_protobuf_duration_pb.Duration): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): RedisSession.AsObject;
    static toObject(includeInstance: boolean, msg: RedisSession): RedisSession.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: RedisSession, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): RedisSession;
    static deserializeBinaryFromReader(message: RedisSession, reader: jspb.BinaryReader): RedisSession;
  }

  export namespace RedisSession {
    export type AsObject = {
      options?: RedisOptions.AsObject,
      keyPrefix: string,
      cookieName: string,
      allowRefreshing?: google_protobuf_wrappers_pb.BoolValue.AsObject,
      preExpiryBuffer?: google_protobuf_duration_pb.Duration.AsObject,
    }
  }

  export class CookieOptions extends jspb.Message {
    hasMaxAge(): boolean;
    clearMaxAge(): void;
    getMaxAge(): google_protobuf_wrappers_pb.UInt32Value | undefined;
    setMaxAge(value?: google_protobuf_wrappers_pb.UInt32Value): void;

    getNotSecure(): boolean;
    setNotSecure(value: boolean): void;

    hasHttpOnly(): boolean;
    clearHttpOnly(): void;
    getHttpOnly(): google_protobuf_wrappers_pb.BoolValue | undefined;
    setHttpOnly(value?: google_protobuf_wrappers_pb.BoolValue): void;

    hasPath(): boolean;
    clearPath(): void;
    getPath(): google_protobuf_wrappers_pb.StringValue | undefined;
    setPath(value?: google_protobuf_wrappers_pb.StringValue): void;

    getSameSite(): UserSession.CookieOptions.SameSiteMap[keyof UserSession.CookieOptions.SameSiteMap];
    setSameSite(value: UserSession.CookieOptions.SameSiteMap[keyof UserSession.CookieOptions.SameSiteMap]): void;

    getDomain(): string;
    setDomain(value: string): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): CookieOptions.AsObject;
    static toObject(includeInstance: boolean, msg: CookieOptions): CookieOptions.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: CookieOptions, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): CookieOptions;
    static deserializeBinaryFromReader(message: CookieOptions, reader: jspb.BinaryReader): CookieOptions;
  }

  export namespace CookieOptions {
    export type AsObject = {
      maxAge?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
      notSecure: boolean,
      httpOnly?: google_protobuf_wrappers_pb.BoolValue.AsObject,
      path?: google_protobuf_wrappers_pb.StringValue.AsObject,
      sameSite: UserSession.CookieOptions.SameSiteMap[keyof UserSession.CookieOptions.SameSiteMap],
      domain: string,
    }

    export interface SameSiteMap {
      DEFAULTMODE: 0;
      LAXMODE: 1;
      STRICTMODE: 2;
      NONEMODE: 3;
    }

    export const SameSite: SameSiteMap;
  }

  export enum SessionCase {
    SESSION_NOT_SET = 0,
    COOKIE = 3,
    REDIS = 4,
  }
}

export class HeaderConfiguration extends jspb.Message {
  getIdTokenHeader(): string;
  setIdTokenHeader(value: string): void;

  getAccessTokenHeader(): string;
  setAccessTokenHeader(value: string): void;

  hasUseBearerSchemaForAuthorization(): boolean;
  clearUseBearerSchemaForAuthorization(): void;
  getUseBearerSchemaForAuthorization(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setUseBearerSchemaForAuthorization(value?: google_protobuf_wrappers_pb.BoolValue): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HeaderConfiguration.AsObject;
  static toObject(includeInstance: boolean, msg: HeaderConfiguration): HeaderConfiguration.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: HeaderConfiguration, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HeaderConfiguration;
  static deserializeBinaryFromReader(message: HeaderConfiguration, reader: jspb.BinaryReader): HeaderConfiguration;
}

export namespace HeaderConfiguration {
  export type AsObject = {
    idTokenHeader: string,
    accessTokenHeader: string,
    useBearerSchemaForAuthorization?: google_protobuf_wrappers_pb.BoolValue.AsObject,
  }
}

export class DiscoveryOverride extends jspb.Message {
  getAuthEndpoint(): string;
  setAuthEndpoint(value: string): void;

  getTokenEndpoint(): string;
  setTokenEndpoint(value: string): void;

  getJwksUri(): string;
  setJwksUri(value: string): void;

  clearScopesList(): void;
  getScopesList(): Array<string>;
  setScopesList(value: Array<string>): void;
  addScopes(value: string, index?: number): string;

  clearResponseTypesList(): void;
  getResponseTypesList(): Array<string>;
  setResponseTypesList(value: Array<string>): void;
  addResponseTypes(value: string, index?: number): string;

  clearSubjectsList(): void;
  getSubjectsList(): Array<string>;
  setSubjectsList(value: Array<string>): void;
  addSubjects(value: string, index?: number): string;

  clearIdTokenAlgsList(): void;
  getIdTokenAlgsList(): Array<string>;
  setIdTokenAlgsList(value: Array<string>): void;
  addIdTokenAlgs(value: string, index?: number): string;

  clearAuthMethodsList(): void;
  getAuthMethodsList(): Array<string>;
  setAuthMethodsList(value: Array<string>): void;
  addAuthMethods(value: string, index?: number): string;

  clearClaimsList(): void;
  getClaimsList(): Array<string>;
  setClaimsList(value: Array<string>): void;
  addClaims(value: string, index?: number): string;

  getRevocationEndpoint(): string;
  setRevocationEndpoint(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DiscoveryOverride.AsObject;
  static toObject(includeInstance: boolean, msg: DiscoveryOverride): DiscoveryOverride.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DiscoveryOverride, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DiscoveryOverride;
  static deserializeBinaryFromReader(message: DiscoveryOverride, reader: jspb.BinaryReader): DiscoveryOverride;
}

export namespace DiscoveryOverride {
  export type AsObject = {
    authEndpoint: string,
    tokenEndpoint: string,
    jwksUri: string,
    scopesList: Array<string>,
    responseTypesList: Array<string>,
    subjectsList: Array<string>,
    idTokenAlgsList: Array<string>,
    authMethodsList: Array<string>,
    claimsList: Array<string>,
    revocationEndpoint: string,
  }
}

export class JwksOnDemandCacheRefreshPolicy extends jspb.Message {
  hasNever(): boolean;
  clearNever(): void;
  getNever(): google_protobuf_empty_pb.Empty | undefined;
  setNever(value?: google_protobuf_empty_pb.Empty): void;

  hasAlways(): boolean;
  clearAlways(): void;
  getAlways(): google_protobuf_empty_pb.Empty | undefined;
  setAlways(value?: google_protobuf_empty_pb.Empty): void;

  hasMaxIdpReqPerPollingInterval(): boolean;
  clearMaxIdpReqPerPollingInterval(): void;
  getMaxIdpReqPerPollingInterval(): number;
  setMaxIdpReqPerPollingInterval(value: number): void;

  getPolicyCase(): JwksOnDemandCacheRefreshPolicy.PolicyCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): JwksOnDemandCacheRefreshPolicy.AsObject;
  static toObject(includeInstance: boolean, msg: JwksOnDemandCacheRefreshPolicy): JwksOnDemandCacheRefreshPolicy.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: JwksOnDemandCacheRefreshPolicy, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): JwksOnDemandCacheRefreshPolicy;
  static deserializeBinaryFromReader(message: JwksOnDemandCacheRefreshPolicy, reader: jspb.BinaryReader): JwksOnDemandCacheRefreshPolicy;
}

export namespace JwksOnDemandCacheRefreshPolicy {
  export type AsObject = {
    never?: google_protobuf_empty_pb.Empty.AsObject,
    always?: google_protobuf_empty_pb.Empty.AsObject,
    maxIdpReqPerPollingInterval: number,
  }

  export enum PolicyCase {
    POLICY_NOT_SET = 0,
    NEVER = 1,
    ALWAYS = 2,
    MAX_IDP_REQ_PER_POLLING_INTERVAL = 3,
  }
}

export class AutoMapFromMetadata extends jspb.Message {
  getNamespace(): string;
  setNamespace(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AutoMapFromMetadata.AsObject;
  static toObject(includeInstance: boolean, msg: AutoMapFromMetadata): AutoMapFromMetadata.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: AutoMapFromMetadata, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AutoMapFromMetadata;
  static deserializeBinaryFromReader(message: AutoMapFromMetadata, reader: jspb.BinaryReader): AutoMapFromMetadata;
}

export namespace AutoMapFromMetadata {
  export type AsObject = {
    namespace: string,
  }
}

export class OidcAuthorizationCode extends jspb.Message {
  getClientId(): string;
  setClientId(value: string): void;

  hasClientSecretRef(): boolean;
  clearClientSecretRef(): void;
  getClientSecretRef(): github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef | undefined;
  setClientSecretRef(value?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef): void;

  getIssuerUrl(): string;
  setIssuerUrl(value: string): void;

  getAuthEndpointQueryParamsMap(): jspb.Map<string, string>;
  clearAuthEndpointQueryParamsMap(): void;
  getTokenEndpointQueryParamsMap(): jspb.Map<string, string>;
  clearTokenEndpointQueryParamsMap(): void;
  getAppUrl(): string;
  setAppUrl(value: string): void;

  getCallbackPath(): string;
  setCallbackPath(value: string): void;

  getLogoutPath(): string;
  setLogoutPath(value: string): void;

  getAfterLogoutUrl(): string;
  setAfterLogoutUrl(value: string): void;

  clearScopesList(): void;
  getScopesList(): Array<string>;
  setScopesList(value: Array<string>): void;
  addScopes(value: string, index?: number): string;

  hasSession(): boolean;
  clearSession(): void;
  getSession(): UserSession | undefined;
  setSession(value?: UserSession): void;

  hasHeaders(): boolean;
  clearHeaders(): void;
  getHeaders(): HeaderConfiguration | undefined;
  setHeaders(value?: HeaderConfiguration): void;

  hasDiscoveryOverride(): boolean;
  clearDiscoveryOverride(): void;
  getDiscoveryOverride(): DiscoveryOverride | undefined;
  setDiscoveryOverride(value?: DiscoveryOverride): void;

  hasDiscoveryPollInterval(): boolean;
  clearDiscoveryPollInterval(): void;
  getDiscoveryPollInterval(): google_protobuf_duration_pb.Duration | undefined;
  setDiscoveryPollInterval(value?: google_protobuf_duration_pb.Duration): void;

  hasJwksCacheRefreshPolicy(): boolean;
  clearJwksCacheRefreshPolicy(): void;
  getJwksCacheRefreshPolicy(): JwksOnDemandCacheRefreshPolicy | undefined;
  setJwksCacheRefreshPolicy(value?: JwksOnDemandCacheRefreshPolicy): void;

  getSessionIdHeaderName(): string;
  setSessionIdHeaderName(value: string): void;

  getParseCallbackPathAsRegex(): boolean;
  setParseCallbackPathAsRegex(value: boolean): void;

  hasAutoMapFromMetadata(): boolean;
  clearAutoMapFromMetadata(): void;
  getAutoMapFromMetadata(): AutoMapFromMetadata | undefined;
  setAutoMapFromMetadata(value?: AutoMapFromMetadata): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OidcAuthorizationCode.AsObject;
  static toObject(includeInstance: boolean, msg: OidcAuthorizationCode): OidcAuthorizationCode.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: OidcAuthorizationCode, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OidcAuthorizationCode;
  static deserializeBinaryFromReader(message: OidcAuthorizationCode, reader: jspb.BinaryReader): OidcAuthorizationCode;
}

export namespace OidcAuthorizationCode {
  export type AsObject = {
    clientId: string,
    clientSecretRef?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject,
    issuerUrl: string,
    authEndpointQueryParamsMap: Array<[string, string]>,
    tokenEndpointQueryParamsMap: Array<[string, string]>,
    appUrl: string,
    callbackPath: string,
    logoutPath: string,
    afterLogoutUrl: string,
    scopesList: Array<string>,
    session?: UserSession.AsObject,
    headers?: HeaderConfiguration.AsObject,
    discoveryOverride?: DiscoveryOverride.AsObject,
    discoveryPollInterval?: google_protobuf_duration_pb.Duration.AsObject,
    jwksCacheRefreshPolicy?: JwksOnDemandCacheRefreshPolicy.AsObject,
    sessionIdHeaderName: string,
    parseCallbackPathAsRegex: boolean,
    autoMapFromMetadata?: AutoMapFromMetadata.AsObject,
  }
}

export class JwtValidation extends jspb.Message {
  hasRemoteJwks(): boolean;
  clearRemoteJwks(): void;
  getRemoteJwks(): JwtValidation.RemoteJwks | undefined;
  setRemoteJwks(value?: JwtValidation.RemoteJwks): void;

  hasLocalJwks(): boolean;
  clearLocalJwks(): void;
  getLocalJwks(): JwtValidation.LocalJwks | undefined;
  setLocalJwks(value?: JwtValidation.LocalJwks): void;

  getIssuer(): string;
  setIssuer(value: string): void;

  getJwksSourceSpecifierCase(): JwtValidation.JwksSourceSpecifierCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): JwtValidation.AsObject;
  static toObject(includeInstance: boolean, msg: JwtValidation): JwtValidation.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: JwtValidation, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): JwtValidation;
  static deserializeBinaryFromReader(message: JwtValidation, reader: jspb.BinaryReader): JwtValidation;
}

export namespace JwtValidation {
  export type AsObject = {
    remoteJwks?: JwtValidation.RemoteJwks.AsObject,
    localJwks?: JwtValidation.LocalJwks.AsObject,
    issuer: string,
  }

  export class RemoteJwks extends jspb.Message {
    getUrl(): string;
    setUrl(value: string): void;

    hasRefreshInterval(): boolean;
    clearRefreshInterval(): void;
    getRefreshInterval(): google_protobuf_duration_pb.Duration | undefined;
    setRefreshInterval(value?: google_protobuf_duration_pb.Duration): void;

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
      refreshInterval?: google_protobuf_duration_pb.Duration.AsObject,
    }
  }

  export class LocalJwks extends jspb.Message {
    getInlineString(): string;
    setInlineString(value: string): void;

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
      inlineString: string,
    }
  }

  export enum JwksSourceSpecifierCase {
    JWKS_SOURCE_SPECIFIER_NOT_SET = 0,
    REMOTE_JWKS = 1,
    LOCAL_JWKS = 2,
  }
}

export class IntrospectionValidation extends jspb.Message {
  getIntrospectionUrl(): string;
  setIntrospectionUrl(value: string): void;

  getClientId(): string;
  setClientId(value: string): void;

  hasClientSecretRef(): boolean;
  clearClientSecretRef(): void;
  getClientSecretRef(): github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef | undefined;
  setClientSecretRef(value?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef): void;

  getUserIdAttributeName(): string;
  setUserIdAttributeName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): IntrospectionValidation.AsObject;
  static toObject(includeInstance: boolean, msg: IntrospectionValidation): IntrospectionValidation.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: IntrospectionValidation, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): IntrospectionValidation;
  static deserializeBinaryFromReader(message: IntrospectionValidation, reader: jspb.BinaryReader): IntrospectionValidation;
}

export namespace IntrospectionValidation {
  export type AsObject = {
    introspectionUrl: string,
    clientId: string,
    clientSecretRef?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject,
    userIdAttributeName: string,
  }
}

export class AccessTokenValidation extends jspb.Message {
  hasIntrospectionUrl(): boolean;
  clearIntrospectionUrl(): void;
  getIntrospectionUrl(): string;
  setIntrospectionUrl(value: string): void;

  hasJwt(): boolean;
  clearJwt(): void;
  getJwt(): JwtValidation | undefined;
  setJwt(value?: JwtValidation): void;

  hasIntrospection(): boolean;
  clearIntrospection(): void;
  getIntrospection(): IntrospectionValidation | undefined;
  setIntrospection(value?: IntrospectionValidation): void;

  getUserinfoUrl(): string;
  setUserinfoUrl(value: string): void;

  hasCacheTimeout(): boolean;
  clearCacheTimeout(): void;
  getCacheTimeout(): google_protobuf_duration_pb.Duration | undefined;
  setCacheTimeout(value?: google_protobuf_duration_pb.Duration): void;

  hasRequiredScopes(): boolean;
  clearRequiredScopes(): void;
  getRequiredScopes(): AccessTokenValidation.ScopeList | undefined;
  setRequiredScopes(value?: AccessTokenValidation.ScopeList): void;

  getValidationTypeCase(): AccessTokenValidation.ValidationTypeCase;
  getScopeValidationCase(): AccessTokenValidation.ScopeValidationCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AccessTokenValidation.AsObject;
  static toObject(includeInstance: boolean, msg: AccessTokenValidation): AccessTokenValidation.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: AccessTokenValidation, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AccessTokenValidation;
  static deserializeBinaryFromReader(message: AccessTokenValidation, reader: jspb.BinaryReader): AccessTokenValidation;
}

export namespace AccessTokenValidation {
  export type AsObject = {
    introspectionUrl: string,
    jwt?: JwtValidation.AsObject,
    introspection?: IntrospectionValidation.AsObject,
    userinfoUrl: string,
    cacheTimeout?: google_protobuf_duration_pb.Duration.AsObject,
    requiredScopes?: AccessTokenValidation.ScopeList.AsObject,
  }

  export class ScopeList extends jspb.Message {
    clearScopeList(): void;
    getScopeList(): Array<string>;
    setScopeList(value: Array<string>): void;
    addScope(value: string, index?: number): string;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ScopeList.AsObject;
    static toObject(includeInstance: boolean, msg: ScopeList): ScopeList.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: ScopeList, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ScopeList;
    static deserializeBinaryFromReader(message: ScopeList, reader: jspb.BinaryReader): ScopeList;
  }

  export namespace ScopeList {
    export type AsObject = {
      scopeList: Array<string>,
    }
  }

  export enum ValidationTypeCase {
    VALIDATION_TYPE_NOT_SET = 0,
    INTROSPECTION_URL = 1,
    JWT = 2,
    INTROSPECTION = 3,
  }

  export enum ScopeValidationCase {
    SCOPE_VALIDATION_NOT_SET = 0,
    REQUIRED_SCOPES = 6,
  }
}

export class OauthSecret extends jspb.Message {
  getClientSecret(): string;
  setClientSecret(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OauthSecret.AsObject;
  static toObject(includeInstance: boolean, msg: OauthSecret): OauthSecret.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: OauthSecret, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OauthSecret;
  static deserializeBinaryFromReader(message: OauthSecret, reader: jspb.BinaryReader): OauthSecret;
}

export namespace OauthSecret {
  export type AsObject = {
    clientSecret: string,
  }
}

export class ApiKeyAuth extends jspb.Message {
  getLabelSelectorMap(): jspb.Map<string, string>;
  clearLabelSelectorMap(): void;
  clearApiKeySecretRefsList(): void;
  getApiKeySecretRefsList(): Array<github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef>;
  setApiKeySecretRefsList(value: Array<github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef>): void;
  addApiKeySecretRefs(value?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef, index?: number): github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef;

  getHeaderName(): string;
  setHeaderName(value: string): void;

  getHeadersFromMetadataMap(): jspb.Map<string, ApiKeyAuth.SecretKey>;
  clearHeadersFromMetadataMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ApiKeyAuth.AsObject;
  static toObject(includeInstance: boolean, msg: ApiKeyAuth): ApiKeyAuth.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ApiKeyAuth, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ApiKeyAuth;
  static deserializeBinaryFromReader(message: ApiKeyAuth, reader: jspb.BinaryReader): ApiKeyAuth;
}

export namespace ApiKeyAuth {
  export type AsObject = {
    labelSelectorMap: Array<[string, string]>,
    apiKeySecretRefsList: Array<github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject>,
    headerName: string,
    headersFromMetadataMap: Array<[string, ApiKeyAuth.SecretKey.AsObject]>,
  }

  export class SecretKey extends jspb.Message {
    getName(): string;
    setName(value: string): void;

    getRequired(): boolean;
    setRequired(value: boolean): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): SecretKey.AsObject;
    static toObject(includeInstance: boolean, msg: SecretKey): SecretKey.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: SecretKey, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): SecretKey;
    static deserializeBinaryFromReader(message: SecretKey, reader: jspb.BinaryReader): SecretKey;
  }

  export namespace SecretKey {
    export type AsObject = {
      name: string,
      required: boolean,
    }
  }
}

export class ApiKeySecret extends jspb.Message {
  getGenerateApiKey(): boolean;
  setGenerateApiKey(value: boolean): void;

  getApiKey(): string;
  setApiKey(value: string): void;

  clearLabelsList(): void;
  getLabelsList(): Array<string>;
  setLabelsList(value: Array<string>): void;
  addLabels(value: string, index?: number): string;

  getMetadataMap(): jspb.Map<string, string>;
  clearMetadataMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ApiKeySecret.AsObject;
  static toObject(includeInstance: boolean, msg: ApiKeySecret): ApiKeySecret.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ApiKeySecret, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ApiKeySecret;
  static deserializeBinaryFromReader(message: ApiKeySecret, reader: jspb.BinaryReader): ApiKeySecret;
}

export namespace ApiKeySecret {
  export type AsObject = {
    generateApiKey: boolean,
    apiKey: string,
    labelsList: Array<string>,
    metadataMap: Array<[string, string]>,
  }
}

export class OpaAuth extends jspb.Message {
  clearModulesList(): void;
  getModulesList(): Array<github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef>;
  setModulesList(value: Array<github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef>): void;
  addModules(value?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef, index?: number): github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef;

  getQuery(): string;
  setQuery(value: string): void;

  hasOptions(): boolean;
  clearOptions(): void;
  getOptions(): OpaAuthOptions | undefined;
  setOptions(value?: OpaAuthOptions): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OpaAuth.AsObject;
  static toObject(includeInstance: boolean, msg: OpaAuth): OpaAuth.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: OpaAuth, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OpaAuth;
  static deserializeBinaryFromReader(message: OpaAuth, reader: jspb.BinaryReader): OpaAuth;
}

export namespace OpaAuth {
  export type AsObject = {
    modulesList: Array<github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject>,
    query: string,
    options?: OpaAuthOptions.AsObject,
  }
}

export class OpaAuthOptions extends jspb.Message {
  getFastInputConversion(): boolean;
  setFastInputConversion(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OpaAuthOptions.AsObject;
  static toObject(includeInstance: boolean, msg: OpaAuthOptions): OpaAuthOptions.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: OpaAuthOptions, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OpaAuthOptions;
  static deserializeBinaryFromReader(message: OpaAuthOptions, reader: jspb.BinaryReader): OpaAuthOptions;
}

export namespace OpaAuthOptions {
  export type AsObject = {
    fastInputConversion: boolean,
  }
}

export class Ldap extends jspb.Message {
  getAddress(): string;
  setAddress(value: string): void;

  getUserdntemplate(): string;
  setUserdntemplate(value: string): void;

  getMembershipattributename(): string;
  setMembershipattributename(value: string): void;

  clearAllowedgroupsList(): void;
  getAllowedgroupsList(): Array<string>;
  setAllowedgroupsList(value: Array<string>): void;
  addAllowedgroups(value: string, index?: number): string;

  hasPool(): boolean;
  clearPool(): void;
  getPool(): Ldap.ConnectionPool | undefined;
  setPool(value?: Ldap.ConnectionPool): void;

  getSearchfilter(): string;
  setSearchfilter(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Ldap.AsObject;
  static toObject(includeInstance: boolean, msg: Ldap): Ldap.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Ldap, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Ldap;
  static deserializeBinaryFromReader(message: Ldap, reader: jspb.BinaryReader): Ldap;
}

export namespace Ldap {
  export type AsObject = {
    address: string,
    userdntemplate: string,
    membershipattributename: string,
    allowedgroupsList: Array<string>,
    pool?: Ldap.ConnectionPool.AsObject,
    searchfilter: string,
  }

  export class ConnectionPool extends jspb.Message {
    hasMaxsize(): boolean;
    clearMaxsize(): void;
    getMaxsize(): google_protobuf_wrappers_pb.UInt32Value | undefined;
    setMaxsize(value?: google_protobuf_wrappers_pb.UInt32Value): void;

    hasInitialsize(): boolean;
    clearInitialsize(): void;
    getInitialsize(): google_protobuf_wrappers_pb.UInt32Value | undefined;
    setInitialsize(value?: google_protobuf_wrappers_pb.UInt32Value): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ConnectionPool.AsObject;
    static toObject(includeInstance: boolean, msg: ConnectionPool): ConnectionPool.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: ConnectionPool, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ConnectionPool;
    static deserializeBinaryFromReader(message: ConnectionPool, reader: jspb.BinaryReader): ConnectionPool;
  }

  export namespace ConnectionPool {
    export type AsObject = {
      maxsize?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
      initialsize?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
    }
  }
}

export class PassThroughAuth extends jspb.Message {
  hasGrpc(): boolean;
  clearGrpc(): void;
  getGrpc(): PassThroughGrpc | undefined;
  setGrpc(value?: PassThroughGrpc): void;

  hasHttp(): boolean;
  clearHttp(): void;
  getHttp(): PassThroughHttp | undefined;
  setHttp(value?: PassThroughHttp): void;

  hasConfig(): boolean;
  clearConfig(): void;
  getConfig(): google_protobuf_struct_pb.Struct | undefined;
  setConfig(value?: google_protobuf_struct_pb.Struct): void;

  getProtocolCase(): PassThroughAuth.ProtocolCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PassThroughAuth.AsObject;
  static toObject(includeInstance: boolean, msg: PassThroughAuth): PassThroughAuth.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: PassThroughAuth, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PassThroughAuth;
  static deserializeBinaryFromReader(message: PassThroughAuth, reader: jspb.BinaryReader): PassThroughAuth;
}

export namespace PassThroughAuth {
  export type AsObject = {
    grpc?: PassThroughGrpc.AsObject,
    http?: PassThroughHttp.AsObject,
    config?: google_protobuf_struct_pb.Struct.AsObject,
  }

  export enum ProtocolCase {
    PROTOCOL_NOT_SET = 0,
    GRPC = 1,
    HTTP = 2,
  }
}

export class PassThroughGrpc extends jspb.Message {
  getAddress(): string;
  setAddress(value: string): void;

  hasConnectionTimeout(): boolean;
  clearConnectionTimeout(): void;
  getConnectionTimeout(): google_protobuf_duration_pb.Duration | undefined;
  setConnectionTimeout(value?: google_protobuf_duration_pb.Duration): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PassThroughGrpc.AsObject;
  static toObject(includeInstance: boolean, msg: PassThroughGrpc): PassThroughGrpc.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: PassThroughGrpc, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PassThroughGrpc;
  static deserializeBinaryFromReader(message: PassThroughGrpc, reader: jspb.BinaryReader): PassThroughGrpc;
}

export namespace PassThroughGrpc {
  export type AsObject = {
    address: string,
    connectionTimeout?: google_protobuf_duration_pb.Duration.AsObject,
  }
}

export class PassThroughHttp extends jspb.Message {
  getUrl(): string;
  setUrl(value: string): void;

  hasRequest(): boolean;
  clearRequest(): void;
  getRequest(): PassThroughHttp.Request | undefined;
  setRequest(value?: PassThroughHttp.Request): void;

  hasResponse(): boolean;
  clearResponse(): void;
  getResponse(): PassThroughHttp.Response | undefined;
  setResponse(value?: PassThroughHttp.Response): void;

  hasConnectionTimeout(): boolean;
  clearConnectionTimeout(): void;
  getConnectionTimeout(): google_protobuf_duration_pb.Duration | undefined;
  setConnectionTimeout(value?: google_protobuf_duration_pb.Duration): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PassThroughHttp.AsObject;
  static toObject(includeInstance: boolean, msg: PassThroughHttp): PassThroughHttp.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: PassThroughHttp, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PassThroughHttp;
  static deserializeBinaryFromReader(message: PassThroughHttp, reader: jspb.BinaryReader): PassThroughHttp;
}

export namespace PassThroughHttp {
  export type AsObject = {
    url: string,
    request?: PassThroughHttp.Request.AsObject,
    response?: PassThroughHttp.Response.AsObject,
    connectionTimeout?: google_protobuf_duration_pb.Duration.AsObject,
  }

  export class Request extends jspb.Message {
    clearAllowedHeadersList(): void;
    getAllowedHeadersList(): Array<string>;
    setAllowedHeadersList(value: Array<string>): void;
    addAllowedHeaders(value: string, index?: number): string;

    getHeadersToAddMap(): jspb.Map<string, string>;
    clearHeadersToAddMap(): void;
    getPassThroughState(): boolean;
    setPassThroughState(value: boolean): void;

    getPassThroughFilterMetadata(): boolean;
    setPassThroughFilterMetadata(value: boolean): void;

    getPassThroughBody(): boolean;
    setPassThroughBody(value: boolean): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Request.AsObject;
    static toObject(includeInstance: boolean, msg: Request): Request.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Request, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Request;
    static deserializeBinaryFromReader(message: Request, reader: jspb.BinaryReader): Request;
  }

  export namespace Request {
    export type AsObject = {
      allowedHeadersList: Array<string>,
      headersToAddMap: Array<[string, string]>,
      passThroughState: boolean,
      passThroughFilterMetadata: boolean,
      passThroughBody: boolean,
    }
  }

  export class Response extends jspb.Message {
    clearAllowedUpstreamHeadersList(): void;
    getAllowedUpstreamHeadersList(): Array<string>;
    setAllowedUpstreamHeadersList(value: Array<string>): void;
    addAllowedUpstreamHeaders(value: string, index?: number): string;

    clearAllowedClientHeadersOnDeniedList(): void;
    getAllowedClientHeadersOnDeniedList(): Array<string>;
    setAllowedClientHeadersOnDeniedList(value: Array<string>): void;
    addAllowedClientHeadersOnDenied(value: string, index?: number): string;

    getReadStateFromResponse(): boolean;
    setReadStateFromResponse(value: boolean): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Response.AsObject;
    static toObject(includeInstance: boolean, msg: Response): Response.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Response, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Response;
    static deserializeBinaryFromReader(message: Response, reader: jspb.BinaryReader): Response;
  }

  export namespace Response {
    export type AsObject = {
      allowedUpstreamHeadersList: Array<string>,
      allowedClientHeadersOnDeniedList: Array<string>,
      readStateFromResponse: boolean,
    }
  }
}

export class ExtAuthConfig extends jspb.Message {
  getAuthConfigRefName(): string;
  setAuthConfigRefName(value: string): void;

  clearConfigsList(): void;
  getConfigsList(): Array<ExtAuthConfig.Config>;
  setConfigsList(value: Array<ExtAuthConfig.Config>): void;
  addConfigs(value?: ExtAuthConfig.Config, index?: number): ExtAuthConfig.Config;

  hasBooleanExpr(): boolean;
  clearBooleanExpr(): void;
  getBooleanExpr(): google_protobuf_wrappers_pb.StringValue | undefined;
  setBooleanExpr(value?: google_protobuf_wrappers_pb.StringValue): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ExtAuthConfig.AsObject;
  static toObject(includeInstance: boolean, msg: ExtAuthConfig): ExtAuthConfig.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ExtAuthConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ExtAuthConfig;
  static deserializeBinaryFromReader(message: ExtAuthConfig, reader: jspb.BinaryReader): ExtAuthConfig;
}

export namespace ExtAuthConfig {
  export type AsObject = {
    authConfigRefName: string,
    configsList: Array<ExtAuthConfig.Config.AsObject>,
    booleanExpr?: google_protobuf_wrappers_pb.StringValue.AsObject,
  }

  export class OAuthConfig extends jspb.Message {
    getClientId(): string;
    setClientId(value: string): void;

    getClientSecret(): string;
    setClientSecret(value: string): void;

    getIssuerUrl(): string;
    setIssuerUrl(value: string): void;

    getAuthEndpointQueryParamsMap(): jspb.Map<string, string>;
    clearAuthEndpointQueryParamsMap(): void;
    getAppUrl(): string;
    setAppUrl(value: string): void;

    getCallbackPath(): string;
    setCallbackPath(value: string): void;

    clearScopesList(): void;
    getScopesList(): Array<string>;
    setScopesList(value: Array<string>): void;
    addScopes(value: string, index?: number): string;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): OAuthConfig.AsObject;
    static toObject(includeInstance: boolean, msg: OAuthConfig): OAuthConfig.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: OAuthConfig, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): OAuthConfig;
    static deserializeBinaryFromReader(message: OAuthConfig, reader: jspb.BinaryReader): OAuthConfig;
  }

  export namespace OAuthConfig {
    export type AsObject = {
      clientId: string,
      clientSecret: string,
      issuerUrl: string,
      authEndpointQueryParamsMap: Array<[string, string]>,
      appUrl: string,
      callbackPath: string,
      scopesList: Array<string>,
    }
  }

  export class OidcAuthorizationCodeConfig extends jspb.Message {
    getClientId(): string;
    setClientId(value: string): void;

    getClientSecret(): string;
    setClientSecret(value: string): void;

    getIssuerUrl(): string;
    setIssuerUrl(value: string): void;

    getAuthEndpointQueryParamsMap(): jspb.Map<string, string>;
    clearAuthEndpointQueryParamsMap(): void;
    getTokenEndpointQueryParamsMap(): jspb.Map<string, string>;
    clearTokenEndpointQueryParamsMap(): void;
    getAppUrl(): string;
    setAppUrl(value: string): void;

    getCallbackPath(): string;
    setCallbackPath(value: string): void;

    getLogoutPath(): string;
    setLogoutPath(value: string): void;

    getAfterLogoutUrl(): string;
    setAfterLogoutUrl(value: string): void;

    clearScopesList(): void;
    getScopesList(): Array<string>;
    setScopesList(value: Array<string>): void;
    addScopes(value: string, index?: number): string;

    hasSession(): boolean;
    clearSession(): void;
    getSession(): UserSession | undefined;
    setSession(value?: UserSession): void;

    hasHeaders(): boolean;
    clearHeaders(): void;
    getHeaders(): HeaderConfiguration | undefined;
    setHeaders(value?: HeaderConfiguration): void;

    hasDiscoveryOverride(): boolean;
    clearDiscoveryOverride(): void;
    getDiscoveryOverride(): DiscoveryOverride | undefined;
    setDiscoveryOverride(value?: DiscoveryOverride): void;

    hasDiscoveryPollInterval(): boolean;
    clearDiscoveryPollInterval(): void;
    getDiscoveryPollInterval(): google_protobuf_duration_pb.Duration | undefined;
    setDiscoveryPollInterval(value?: google_protobuf_duration_pb.Duration): void;

    hasJwksCacheRefreshPolicy(): boolean;
    clearJwksCacheRefreshPolicy(): void;
    getJwksCacheRefreshPolicy(): JwksOnDemandCacheRefreshPolicy | undefined;
    setJwksCacheRefreshPolicy(value?: JwksOnDemandCacheRefreshPolicy): void;

    getSessionIdHeaderName(): string;
    setSessionIdHeaderName(value: string): void;

    getParseCallbackPathAsRegex(): boolean;
    setParseCallbackPathAsRegex(value: boolean): void;

    hasAutoMapFromMetadata(): boolean;
    clearAutoMapFromMetadata(): void;
    getAutoMapFromMetadata(): AutoMapFromMetadata | undefined;
    setAutoMapFromMetadata(value?: AutoMapFromMetadata): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): OidcAuthorizationCodeConfig.AsObject;
    static toObject(includeInstance: boolean, msg: OidcAuthorizationCodeConfig): OidcAuthorizationCodeConfig.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: OidcAuthorizationCodeConfig, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): OidcAuthorizationCodeConfig;
    static deserializeBinaryFromReader(message: OidcAuthorizationCodeConfig, reader: jspb.BinaryReader): OidcAuthorizationCodeConfig;
  }

  export namespace OidcAuthorizationCodeConfig {
    export type AsObject = {
      clientId: string,
      clientSecret: string,
      issuerUrl: string,
      authEndpointQueryParamsMap: Array<[string, string]>,
      tokenEndpointQueryParamsMap: Array<[string, string]>,
      appUrl: string,
      callbackPath: string,
      logoutPath: string,
      afterLogoutUrl: string,
      scopesList: Array<string>,
      session?: UserSession.AsObject,
      headers?: HeaderConfiguration.AsObject,
      discoveryOverride?: DiscoveryOverride.AsObject,
      discoveryPollInterval?: google_protobuf_duration_pb.Duration.AsObject,
      jwksCacheRefreshPolicy?: JwksOnDemandCacheRefreshPolicy.AsObject,
      sessionIdHeaderName: string,
      parseCallbackPathAsRegex: boolean,
      autoMapFromMetadata?: AutoMapFromMetadata.AsObject,
    }
  }

  export class AccessTokenValidationConfig extends jspb.Message {
    hasIntrospectionUrl(): boolean;
    clearIntrospectionUrl(): void;
    getIntrospectionUrl(): string;
    setIntrospectionUrl(value: string): void;

    hasJwt(): boolean;
    clearJwt(): void;
    getJwt(): ExtAuthConfig.AccessTokenValidationConfig.JwtValidation | undefined;
    setJwt(value?: ExtAuthConfig.AccessTokenValidationConfig.JwtValidation): void;

    hasIntrospection(): boolean;
    clearIntrospection(): void;
    getIntrospection(): ExtAuthConfig.AccessTokenValidationConfig.IntrospectionValidation | undefined;
    setIntrospection(value?: ExtAuthConfig.AccessTokenValidationConfig.IntrospectionValidation): void;

    getUserinfoUrl(): string;
    setUserinfoUrl(value: string): void;

    hasCacheTimeout(): boolean;
    clearCacheTimeout(): void;
    getCacheTimeout(): google_protobuf_duration_pb.Duration | undefined;
    setCacheTimeout(value?: google_protobuf_duration_pb.Duration): void;

    hasRequiredScopes(): boolean;
    clearRequiredScopes(): void;
    getRequiredScopes(): ExtAuthConfig.AccessTokenValidationConfig.ScopeList | undefined;
    setRequiredScopes(value?: ExtAuthConfig.AccessTokenValidationConfig.ScopeList): void;

    getValidationTypeCase(): AccessTokenValidationConfig.ValidationTypeCase;
    getScopeValidationCase(): AccessTokenValidationConfig.ScopeValidationCase;
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): AccessTokenValidationConfig.AsObject;
    static toObject(includeInstance: boolean, msg: AccessTokenValidationConfig): AccessTokenValidationConfig.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: AccessTokenValidationConfig, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): AccessTokenValidationConfig;
    static deserializeBinaryFromReader(message: AccessTokenValidationConfig, reader: jspb.BinaryReader): AccessTokenValidationConfig;
  }

  export namespace AccessTokenValidationConfig {
    export type AsObject = {
      introspectionUrl: string,
      jwt?: ExtAuthConfig.AccessTokenValidationConfig.JwtValidation.AsObject,
      introspection?: ExtAuthConfig.AccessTokenValidationConfig.IntrospectionValidation.AsObject,
      userinfoUrl: string,
      cacheTimeout?: google_protobuf_duration_pb.Duration.AsObject,
      requiredScopes?: ExtAuthConfig.AccessTokenValidationConfig.ScopeList.AsObject,
    }

    export class JwtValidation extends jspb.Message {
      hasRemoteJwks(): boolean;
      clearRemoteJwks(): void;
      getRemoteJwks(): ExtAuthConfig.AccessTokenValidationConfig.JwtValidation.RemoteJwks | undefined;
      setRemoteJwks(value?: ExtAuthConfig.AccessTokenValidationConfig.JwtValidation.RemoteJwks): void;

      hasLocalJwks(): boolean;
      clearLocalJwks(): void;
      getLocalJwks(): ExtAuthConfig.AccessTokenValidationConfig.JwtValidation.LocalJwks | undefined;
      setLocalJwks(value?: ExtAuthConfig.AccessTokenValidationConfig.JwtValidation.LocalJwks): void;

      getIssuer(): string;
      setIssuer(value: string): void;

      getJwksSourceSpecifierCase(): JwtValidation.JwksSourceSpecifierCase;
      serializeBinary(): Uint8Array;
      toObject(includeInstance?: boolean): JwtValidation.AsObject;
      static toObject(includeInstance: boolean, msg: JwtValidation): JwtValidation.AsObject;
      static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
      static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
      static serializeBinaryToWriter(message: JwtValidation, writer: jspb.BinaryWriter): void;
      static deserializeBinary(bytes: Uint8Array): JwtValidation;
      static deserializeBinaryFromReader(message: JwtValidation, reader: jspb.BinaryReader): JwtValidation;
    }

    export namespace JwtValidation {
      export type AsObject = {
        remoteJwks?: ExtAuthConfig.AccessTokenValidationConfig.JwtValidation.RemoteJwks.AsObject,
        localJwks?: ExtAuthConfig.AccessTokenValidationConfig.JwtValidation.LocalJwks.AsObject,
        issuer: string,
      }

      export class RemoteJwks extends jspb.Message {
        getUrl(): string;
        setUrl(value: string): void;

        hasRefreshInterval(): boolean;
        clearRefreshInterval(): void;
        getRefreshInterval(): google_protobuf_duration_pb.Duration | undefined;
        setRefreshInterval(value?: google_protobuf_duration_pb.Duration): void;

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
          refreshInterval?: google_protobuf_duration_pb.Duration.AsObject,
        }
      }

      export class LocalJwks extends jspb.Message {
        getInlineString(): string;
        setInlineString(value: string): void;

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
          inlineString: string,
        }
      }

      export enum JwksSourceSpecifierCase {
        JWKS_SOURCE_SPECIFIER_NOT_SET = 0,
        REMOTE_JWKS = 1,
        LOCAL_JWKS = 2,
      }
    }

    export class IntrospectionValidation extends jspb.Message {
      getIntrospectionUrl(): string;
      setIntrospectionUrl(value: string): void;

      getClientId(): string;
      setClientId(value: string): void;

      getClientSecret(): string;
      setClientSecret(value: string): void;

      getUserIdAttributeName(): string;
      setUserIdAttributeName(value: string): void;

      serializeBinary(): Uint8Array;
      toObject(includeInstance?: boolean): IntrospectionValidation.AsObject;
      static toObject(includeInstance: boolean, msg: IntrospectionValidation): IntrospectionValidation.AsObject;
      static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
      static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
      static serializeBinaryToWriter(message: IntrospectionValidation, writer: jspb.BinaryWriter): void;
      static deserializeBinary(bytes: Uint8Array): IntrospectionValidation;
      static deserializeBinaryFromReader(message: IntrospectionValidation, reader: jspb.BinaryReader): IntrospectionValidation;
    }

    export namespace IntrospectionValidation {
      export type AsObject = {
        introspectionUrl: string,
        clientId: string,
        clientSecret: string,
        userIdAttributeName: string,
      }
    }

    export class ScopeList extends jspb.Message {
      clearScopeList(): void;
      getScopeList(): Array<string>;
      setScopeList(value: Array<string>): void;
      addScope(value: string, index?: number): string;

      serializeBinary(): Uint8Array;
      toObject(includeInstance?: boolean): ScopeList.AsObject;
      static toObject(includeInstance: boolean, msg: ScopeList): ScopeList.AsObject;
      static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
      static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
      static serializeBinaryToWriter(message: ScopeList, writer: jspb.BinaryWriter): void;
      static deserializeBinary(bytes: Uint8Array): ScopeList;
      static deserializeBinaryFromReader(message: ScopeList, reader: jspb.BinaryReader): ScopeList;
    }

    export namespace ScopeList {
      export type AsObject = {
        scopeList: Array<string>,
      }
    }

    export enum ValidationTypeCase {
      VALIDATION_TYPE_NOT_SET = 0,
      INTROSPECTION_URL = 1,
      JWT = 2,
      INTROSPECTION = 3,
    }

    export enum ScopeValidationCase {
      SCOPE_VALIDATION_NOT_SET = 0,
      REQUIRED_SCOPES = 6,
    }
  }

  export class OAuth2Config extends jspb.Message {
    hasOidcAuthorizationCode(): boolean;
    clearOidcAuthorizationCode(): void;
    getOidcAuthorizationCode(): ExtAuthConfig.OidcAuthorizationCodeConfig | undefined;
    setOidcAuthorizationCode(value?: ExtAuthConfig.OidcAuthorizationCodeConfig): void;

    hasAccessTokenValidationConfig(): boolean;
    clearAccessTokenValidationConfig(): void;
    getAccessTokenValidationConfig(): ExtAuthConfig.AccessTokenValidationConfig | undefined;
    setAccessTokenValidationConfig(value?: ExtAuthConfig.AccessTokenValidationConfig): void;

    getOauthTypeCase(): OAuth2Config.OauthTypeCase;
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): OAuth2Config.AsObject;
    static toObject(includeInstance: boolean, msg: OAuth2Config): OAuth2Config.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: OAuth2Config, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): OAuth2Config;
    static deserializeBinaryFromReader(message: OAuth2Config, reader: jspb.BinaryReader): OAuth2Config;
  }

  export namespace OAuth2Config {
    export type AsObject = {
      oidcAuthorizationCode?: ExtAuthConfig.OidcAuthorizationCodeConfig.AsObject,
      accessTokenValidationConfig?: ExtAuthConfig.AccessTokenValidationConfig.AsObject,
    }

    export enum OauthTypeCase {
      OAUTH_TYPE_NOT_SET = 0,
      OIDC_AUTHORIZATION_CODE = 1,
      ACCESS_TOKEN_VALIDATION_CONFIG = 3,
    }
  }

  export class ApiKeyAuthConfig extends jspb.Message {
    getValidApiKeysMap(): jspb.Map<string, ExtAuthConfig.ApiKeyAuthConfig.KeyMetadata>;
    clearValidApiKeysMap(): void;
    getHeaderName(): string;
    setHeaderName(value: string): void;

    getHeadersFromKeyMetadataMap(): jspb.Map<string, string>;
    clearHeadersFromKeyMetadataMap(): void;
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ApiKeyAuthConfig.AsObject;
    static toObject(includeInstance: boolean, msg: ApiKeyAuthConfig): ApiKeyAuthConfig.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: ApiKeyAuthConfig, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ApiKeyAuthConfig;
    static deserializeBinaryFromReader(message: ApiKeyAuthConfig, reader: jspb.BinaryReader): ApiKeyAuthConfig;
  }

  export namespace ApiKeyAuthConfig {
    export type AsObject = {
      validApiKeysMap: Array<[string, ExtAuthConfig.ApiKeyAuthConfig.KeyMetadata.AsObject]>,
      headerName: string,
      headersFromKeyMetadataMap: Array<[string, string]>,
    }

    export class KeyMetadata extends jspb.Message {
      getUsername(): string;
      setUsername(value: string): void;

      getMetadataMap(): jspb.Map<string, string>;
      clearMetadataMap(): void;
      serializeBinary(): Uint8Array;
      toObject(includeInstance?: boolean): KeyMetadata.AsObject;
      static toObject(includeInstance: boolean, msg: KeyMetadata): KeyMetadata.AsObject;
      static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
      static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
      static serializeBinaryToWriter(message: KeyMetadata, writer: jspb.BinaryWriter): void;
      static deserializeBinary(bytes: Uint8Array): KeyMetadata;
      static deserializeBinaryFromReader(message: KeyMetadata, reader: jspb.BinaryReader): KeyMetadata;
    }

    export namespace KeyMetadata {
      export type AsObject = {
        username: string,
        metadataMap: Array<[string, string]>,
      }
    }
  }

  export class OpaAuthConfig extends jspb.Message {
    getModulesMap(): jspb.Map<string, string>;
    clearModulesMap(): void;
    getQuery(): string;
    setQuery(value: string): void;

    hasOptions(): boolean;
    clearOptions(): void;
    getOptions(): OpaAuthOptions | undefined;
    setOptions(value?: OpaAuthOptions): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): OpaAuthConfig.AsObject;
    static toObject(includeInstance: boolean, msg: OpaAuthConfig): OpaAuthConfig.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: OpaAuthConfig, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): OpaAuthConfig;
    static deserializeBinaryFromReader(message: OpaAuthConfig, reader: jspb.BinaryReader): OpaAuthConfig;
  }

  export namespace OpaAuthConfig {
    export type AsObject = {
      modulesMap: Array<[string, string]>,
      query: string,
      options?: OpaAuthOptions.AsObject,
    }
  }

  export class Config extends jspb.Message {
    hasName(): boolean;
    clearName(): void;
    getName(): google_protobuf_wrappers_pb.StringValue | undefined;
    setName(value?: google_protobuf_wrappers_pb.StringValue): void;

    hasOauth(): boolean;
    clearOauth(): void;
    getOauth(): ExtAuthConfig.OAuthConfig | undefined;
    setOauth(value?: ExtAuthConfig.OAuthConfig): void;

    hasOauth2(): boolean;
    clearOauth2(): void;
    getOauth2(): ExtAuthConfig.OAuth2Config | undefined;
    setOauth2(value?: ExtAuthConfig.OAuth2Config): void;

    hasBasicAuth(): boolean;
    clearBasicAuth(): void;
    getBasicAuth(): BasicAuth | undefined;
    setBasicAuth(value?: BasicAuth): void;

    hasApiKeyAuth(): boolean;
    clearApiKeyAuth(): void;
    getApiKeyAuth(): ExtAuthConfig.ApiKeyAuthConfig | undefined;
    setApiKeyAuth(value?: ExtAuthConfig.ApiKeyAuthConfig): void;

    hasPluginAuth(): boolean;
    clearPluginAuth(): void;
    getPluginAuth(): AuthPlugin | undefined;
    setPluginAuth(value?: AuthPlugin): void;

    hasOpaAuth(): boolean;
    clearOpaAuth(): void;
    getOpaAuth(): ExtAuthConfig.OpaAuthConfig | undefined;
    setOpaAuth(value?: ExtAuthConfig.OpaAuthConfig): void;

    hasLdap(): boolean;
    clearLdap(): void;
    getLdap(): Ldap | undefined;
    setLdap(value?: Ldap): void;

    hasJwt(): boolean;
    clearJwt(): void;
    getJwt(): google_protobuf_empty_pb.Empty | undefined;
    setJwt(value?: google_protobuf_empty_pb.Empty): void;

    hasPassThroughAuth(): boolean;
    clearPassThroughAuth(): void;
    getPassThroughAuth(): PassThroughAuth | undefined;
    setPassThroughAuth(value?: PassThroughAuth): void;

    getAuthConfigCase(): Config.AuthConfigCase;
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Config.AsObject;
    static toObject(includeInstance: boolean, msg: Config): Config.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Config, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Config;
    static deserializeBinaryFromReader(message: Config, reader: jspb.BinaryReader): Config;
  }

  export namespace Config {
    export type AsObject = {
      name?: google_protobuf_wrappers_pb.StringValue.AsObject,
      oauth?: ExtAuthConfig.OAuthConfig.AsObject,
      oauth2?: ExtAuthConfig.OAuth2Config.AsObject,
      basicAuth?: BasicAuth.AsObject,
      apiKeyAuth?: ExtAuthConfig.ApiKeyAuthConfig.AsObject,
      pluginAuth?: AuthPlugin.AsObject,
      opaAuth?: ExtAuthConfig.OpaAuthConfig.AsObject,
      ldap?: Ldap.AsObject,
      jwt?: google_protobuf_empty_pb.Empty.AsObject,
      passThroughAuth?: PassThroughAuth.AsObject,
    }

    export enum AuthConfigCase {
      AUTH_CONFIG_NOT_SET = 0,
      OAUTH = 3,
      OAUTH2 = 9,
      BASIC_AUTH = 4,
      API_KEY_AUTH = 5,
      PLUGIN_AUTH = 6,
      OPA_AUTH = 7,
      LDAP = 8,
      JWT = 12,
      PASS_THROUGH_AUTH = 13,
    }
  }
}

export class AuthConfigStatus extends jspb.Message {
  getState(): AuthConfigStatus.StateMap[keyof AuthConfigStatus.StateMap];
  setState(value: AuthConfigStatus.StateMap[keyof AuthConfigStatus.StateMap]): void;

  getReason(): string;
  setReason(value: string): void;

  getReportedBy(): string;
  setReportedBy(value: string): void;

  getSubresourceStatusesMap(): jspb.Map<string, AuthConfigStatus>;
  clearSubresourceStatusesMap(): void;
  hasDetails(): boolean;
  clearDetails(): void;
  getDetails(): google_protobuf_struct_pb.Struct | undefined;
  setDetails(value?: google_protobuf_struct_pb.Struct): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AuthConfigStatus.AsObject;
  static toObject(includeInstance: boolean, msg: AuthConfigStatus): AuthConfigStatus.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: AuthConfigStatus, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AuthConfigStatus;
  static deserializeBinaryFromReader(message: AuthConfigStatus, reader: jspb.BinaryReader): AuthConfigStatus;
}

export namespace AuthConfigStatus {
  export type AsObject = {
    state: AuthConfigStatus.StateMap[keyof AuthConfigStatus.StateMap],
    reason: string,
    reportedBy: string,
    subresourceStatusesMap: Array<[string, AuthConfigStatus.AsObject]>,
    details?: google_protobuf_struct_pb.Struct.AsObject,
  }

  export interface StateMap {
    PENDING: 0;
    ACCEPTED: 1;
    REJECTED: 2;
    WARNING: 3;
  }

  export const State: StateMap;
}

export class AuthConfigNamespacedStatuses extends jspb.Message {
  getStatusesMap(): jspb.Map<string, AuthConfigStatus>;
  clearStatusesMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AuthConfigNamespacedStatuses.AsObject;
  static toObject(includeInstance: boolean, msg: AuthConfigNamespacedStatuses): AuthConfigNamespacedStatuses.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: AuthConfigNamespacedStatuses, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AuthConfigNamespacedStatuses;
  static deserializeBinaryFromReader(message: AuthConfigNamespacedStatuses, reader: jspb.BinaryReader): AuthConfigNamespacedStatuses;
}

export namespace AuthConfigNamespacedStatuses {
  export type AsObject = {
    statusesMap: Array<[string, AuthConfigStatus.AsObject]>,
  }
}
