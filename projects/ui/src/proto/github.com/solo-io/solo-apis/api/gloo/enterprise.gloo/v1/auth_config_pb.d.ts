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

  getFailOnRedirect(): boolean;
  setFailOnRedirect(value: boolean): void;

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
    failOnRedirect: boolean,
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

    hasHmacAuth(): boolean;
    clearHmacAuth(): void;
    getHmacAuth(): HmacAuth | undefined;
    setHmacAuth(value?: HmacAuth): void;

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
      hmacAuth?: HmacAuth.AsObject,
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
      HMAC_AUTH = 13,
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

export class HmacAuth extends jspb.Message {
  hasSecretRefs(): boolean;
  clearSecretRefs(): void;
  getSecretRefs(): SecretRefList | undefined;
  setSecretRefs(value?: SecretRefList): void;

  hasParametersInHeaders(): boolean;
  clearParametersInHeaders(): void;
  getParametersInHeaders(): HmacParametersInHeaders | undefined;
  setParametersInHeaders(value?: HmacParametersInHeaders): void;

  getSecretStorageCase(): HmacAuth.SecretStorageCase;
  getImplementationTypeCase(): HmacAuth.ImplementationTypeCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HmacAuth.AsObject;
  static toObject(includeInstance: boolean, msg: HmacAuth): HmacAuth.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: HmacAuth, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HmacAuth;
  static deserializeBinaryFromReader(message: HmacAuth, reader: jspb.BinaryReader): HmacAuth;
}

export namespace HmacAuth {
  export type AsObject = {
    secretRefs?: SecretRefList.AsObject,
    parametersInHeaders?: HmacParametersInHeaders.AsObject,
  }

  export enum SecretStorageCase {
    SECRET_STORAGE_NOT_SET = 0,
    SECRET_REFS = 1,
  }

  export enum ImplementationTypeCase {
    IMPLEMENTATION_TYPE_NOT_SET = 0,
    PARAMETERS_IN_HEADERS = 2,
  }
}

export class SecretRefList extends jspb.Message {
  clearSecretRefsList(): void;
  getSecretRefsList(): Array<github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef>;
  setSecretRefsList(value: Array<github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef>): void;
  addSecretRefs(value?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef, index?: number): github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SecretRefList.AsObject;
  static toObject(includeInstance: boolean, msg: SecretRefList): SecretRefList.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: SecretRefList, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SecretRefList;
  static deserializeBinaryFromReader(message: SecretRefList, reader: jspb.BinaryReader): SecretRefList;
}

export namespace SecretRefList {
  export type AsObject = {
    secretRefsList: Array<github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject>,
  }
}

export class HmacParametersInHeaders extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HmacParametersInHeaders.AsObject;
  static toObject(includeInstance: boolean, msg: HmacParametersInHeaders): HmacParametersInHeaders.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: HmacParametersInHeaders, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HmacParametersInHeaders;
  static deserializeBinaryFromReader(message: HmacParametersInHeaders, reader: jspb.BinaryReader): HmacParametersInHeaders;
}

export namespace HmacParametersInHeaders {
  export type AsObject = {
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

  hasOauth2(): boolean;
  clearOauth2(): void;
  getOauth2(): PlainOAuth2 | undefined;
  setOauth2(value?: PlainOAuth2): void;

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
    oauth2?: PlainOAuth2.AsObject,
  }

  export enum OauthTypeCase {
    OAUTH_TYPE_NOT_SET = 0,
    OIDC_AUTHORIZATION_CODE = 1,
    ACCESS_TOKEN_VALIDATION = 2,
    OAUTH2 = 3,
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

  hasCipherConfig(): boolean;
  clearCipherConfig(): void;
  getCipherConfig(): UserSession.CipherConfig | undefined;
  setCipherConfig(value?: UserSession.CipherConfig): void;

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
    cipherConfig?: UserSession.CipherConfig.AsObject,
  }

  export class InternalSession extends jspb.Message {
    hasAllowRefreshing(): boolean;
    clearAllowRefreshing(): void;
    getAllowRefreshing(): google_protobuf_wrappers_pb.BoolValue | undefined;
    setAllowRefreshing(value?: google_protobuf_wrappers_pb.BoolValue): void;

    getKeyPrefix(): string;
    setKeyPrefix(value: string): void;

    getTargetDomain(): string;
    setTargetDomain(value: string): void;

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
      allowRefreshing?: google_protobuf_wrappers_pb.BoolValue.AsObject,
      keyPrefix: string,
      targetDomain: string,
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

    getTargetDomain(): string;
    setTargetDomain(value: string): void;

    getHeaderName(): string;
    setHeaderName(value: string): void;

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
      targetDomain: string,
      headerName: string,
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

  export class CipherConfig extends jspb.Message {
    hasKeyRef(): boolean;
    clearKeyRef(): void;
    getKeyRef(): github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef | undefined;
    setKeyRef(value?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef): void;

    getKeyCase(): CipherConfig.KeyCase;
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): CipherConfig.AsObject;
    static toObject(includeInstance: boolean, msg: CipherConfig): CipherConfig.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: CipherConfig, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): CipherConfig;
    static deserializeBinaryFromReader(message: CipherConfig, reader: jspb.BinaryReader): CipherConfig;
  }

  export namespace CipherConfig {
    export type AsObject = {
      keyRef?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject,
    }

    export enum KeyCase {
      KEY_NOT_SET = 0,
      KEY_REF = 1,
    }
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

  getEndSessionEndpoint(): string;
  setEndSessionEndpoint(value: string): void;

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
    endSessionEndpoint: string,
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

export class EndSessionProperties extends jspb.Message {
  getMethodtype(): EndSessionProperties.MethodTypeMap[keyof EndSessionProperties.MethodTypeMap];
  setMethodtype(value: EndSessionProperties.MethodTypeMap[keyof EndSessionProperties.MethodTypeMap]): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EndSessionProperties.AsObject;
  static toObject(includeInstance: boolean, msg: EndSessionProperties): EndSessionProperties.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: EndSessionProperties, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EndSessionProperties;
  static deserializeBinaryFromReader(message: EndSessionProperties, reader: jspb.BinaryReader): EndSessionProperties;
}

export namespace EndSessionProperties {
  export type AsObject = {
    methodtype: EndSessionProperties.MethodTypeMap[keyof EndSessionProperties.MethodTypeMap],
  }

  export interface MethodTypeMap {
    GETMETHOD: 0;
    POSTMETHOD: 1;
  }

  export const MethodType: MethodTypeMap;
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

  hasEndSessionProperties(): boolean;
  clearEndSessionProperties(): void;
  getEndSessionProperties(): EndSessionProperties | undefined;
  setEndSessionProperties(value?: EndSessionProperties): void;

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
    endSessionProperties?: EndSessionProperties.AsObject,
  }
}

export class PlainOAuth2 extends jspb.Message {
  getClientId(): string;
  setClientId(value: string): void;

  hasClientSecretRef(): boolean;
  clearClientSecretRef(): void;
  getClientSecretRef(): github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef | undefined;
  setClientSecretRef(value?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef): void;

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

  hasSession(): boolean;
  clearSession(): void;
  getSession(): UserSession | undefined;
  setSession(value?: UserSession): void;

  getLogoutPath(): string;
  setLogoutPath(value: string): void;

  getTokenEndpointQueryParamsMap(): jspb.Map<string, string>;
  clearTokenEndpointQueryParamsMap(): void;
  getAfterLogoutUrl(): string;
  setAfterLogoutUrl(value: string): void;

  getAuthEndpoint(): string;
  setAuthEndpoint(value: string): void;

  getTokenEndpoint(): string;
  setTokenEndpoint(value: string): void;

  getRevocationEndpoint(): string;
  setRevocationEndpoint(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PlainOAuth2.AsObject;
  static toObject(includeInstance: boolean, msg: PlainOAuth2): PlainOAuth2.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: PlainOAuth2, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PlainOAuth2;
  static deserializeBinaryFromReader(message: PlainOAuth2, reader: jspb.BinaryReader): PlainOAuth2;
}

export namespace PlainOAuth2 {
  export type AsObject = {
    clientId: string,
    clientSecretRef?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject,
    authEndpointQueryParamsMap: Array<[string, string]>,
    appUrl: string,
    callbackPath: string,
    scopesList: Array<string>,
    session?: UserSession.AsObject,
    logoutPath: string,
    tokenEndpointQueryParamsMap: Array<[string, string]>,
    afterLogoutUrl: string,
    authEndpoint: string,
    tokenEndpoint: string,
    revocationEndpoint: string,
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
  getHeadersFromMetadataEntryMap(): jspb.Map<string, ApiKeyAuth.MetadataEntry>;
  clearHeadersFromMetadataEntryMap(): void;
  hasK8sSecretApikeyStorage(): boolean;
  clearK8sSecretApikeyStorage(): void;
  getK8sSecretApikeyStorage(): K8sSecretApiKeyStorage | undefined;
  setK8sSecretApikeyStorage(value?: K8sSecretApiKeyStorage): void;

  hasAerospikeApikeyStorage(): boolean;
  clearAerospikeApikeyStorage(): void;
  getAerospikeApikeyStorage(): AerospikeApiKeyStorage | undefined;
  setAerospikeApikeyStorage(value?: AerospikeApiKeyStorage): void;

  getStorageBackendCase(): ApiKeyAuth.StorageBackendCase;
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
    headersFromMetadataEntryMap: Array<[string, ApiKeyAuth.MetadataEntry.AsObject]>,
    k8sSecretApikeyStorage?: K8sSecretApiKeyStorage.AsObject,
    aerospikeApikeyStorage?: AerospikeApiKeyStorage.AsObject,
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

  export class MetadataEntry extends jspb.Message {
    getName(): string;
    setName(value: string): void;

    getRequired(): boolean;
    setRequired(value: boolean): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): MetadataEntry.AsObject;
    static toObject(includeInstance: boolean, msg: MetadataEntry): MetadataEntry.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: MetadataEntry, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): MetadataEntry;
    static deserializeBinaryFromReader(message: MetadataEntry, reader: jspb.BinaryReader): MetadataEntry;
  }

  export namespace MetadataEntry {
    export type AsObject = {
      name: string,
      required: boolean,
    }
  }

  export enum StorageBackendCase {
    STORAGE_BACKEND_NOT_SET = 0,
    K8S_SECRET_APIKEY_STORAGE = 6,
    AEROSPIKE_APIKEY_STORAGE = 7,
  }
}

export class K8sSecretApiKeyStorage extends jspb.Message {
  getLabelSelectorMap(): jspb.Map<string, string>;
  clearLabelSelectorMap(): void;
  clearApiKeySecretRefsList(): void;
  getApiKeySecretRefsList(): Array<github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef>;
  setApiKeySecretRefsList(value: Array<github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef>): void;
  addApiKeySecretRefs(value?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef, index?: number): github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): K8sSecretApiKeyStorage.AsObject;
  static toObject(includeInstance: boolean, msg: K8sSecretApiKeyStorage): K8sSecretApiKeyStorage.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: K8sSecretApiKeyStorage, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): K8sSecretApiKeyStorage;
  static deserializeBinaryFromReader(message: K8sSecretApiKeyStorage, reader: jspb.BinaryReader): K8sSecretApiKeyStorage;
}

export namespace K8sSecretApiKeyStorage {
  export type AsObject = {
    labelSelectorMap: Array<[string, string]>,
    apiKeySecretRefsList: Array<github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject>,
  }
}

export class AerospikeApiKeyStorage extends jspb.Message {
  getHostname(): string;
  setHostname(value: string): void;

  getNamespace(): string;
  setNamespace(value: string): void;

  getSet(): string;
  setSet(value: string): void;

  getPort(): number;
  setPort(value: number): void;

  getBatchSize(): number;
  setBatchSize(value: number): void;

  hasCommitAll(): boolean;
  clearCommitAll(): void;
  getCommitAll(): number;
  setCommitAll(value: number): void;

  hasCommitMaster(): boolean;
  clearCommitMaster(): void;
  getCommitMaster(): number;
  setCommitMaster(value: number): void;

  hasReadModeSc(): boolean;
  clearReadModeSc(): void;
  getReadModeSc(): AerospikeApiKeyStorage.readModeSc | undefined;
  setReadModeSc(value?: AerospikeApiKeyStorage.readModeSc): void;

  hasReadModeAp(): boolean;
  clearReadModeAp(): void;
  getReadModeAp(): AerospikeApiKeyStorage.readModeAp | undefined;
  setReadModeAp(value?: AerospikeApiKeyStorage.readModeAp): void;

  getNodeTlsName(): string;
  setNodeTlsName(value: string): void;

  getCertPath(): string;
  setCertPath(value: string): void;

  getKeyPath(): string;
  setKeyPath(value: string): void;

  getAllowInsecure(): boolean;
  setAllowInsecure(value: boolean): void;

  getRootCaPath(): string;
  setRootCaPath(value: string): void;

  getTlsVersion(): string;
  setTlsVersion(value: string): void;

  clearTlsCurveGroupsList(): void;
  getTlsCurveGroupsList(): Array<AerospikeApiKeyStorage.tlsCurveID>;
  setTlsCurveGroupsList(value: Array<AerospikeApiKeyStorage.tlsCurveID>): void;
  addTlsCurveGroups(value?: AerospikeApiKeyStorage.tlsCurveID, index?: number): AerospikeApiKeyStorage.tlsCurveID;

  getCommitLevelCase(): AerospikeApiKeyStorage.CommitLevelCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AerospikeApiKeyStorage.AsObject;
  static toObject(includeInstance: boolean, msg: AerospikeApiKeyStorage): AerospikeApiKeyStorage.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: AerospikeApiKeyStorage, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AerospikeApiKeyStorage;
  static deserializeBinaryFromReader(message: AerospikeApiKeyStorage, reader: jspb.BinaryReader): AerospikeApiKeyStorage;
}

export namespace AerospikeApiKeyStorage {
  export type AsObject = {
    hostname: string,
    namespace: string,
    set: string,
    port: number,
    batchSize: number,
    commitAll: number,
    commitMaster: number,
    readModeSc?: AerospikeApiKeyStorage.readModeSc.AsObject,
    readModeAp?: AerospikeApiKeyStorage.readModeAp.AsObject,
    nodeTlsName: string,
    certPath: string,
    keyPath: string,
    allowInsecure: boolean,
    rootCaPath: string,
    tlsVersion: string,
    tlsCurveGroupsList: Array<AerospikeApiKeyStorage.tlsCurveID.AsObject>,
  }

  export class readModeSc extends jspb.Message {
    hasReadModeScSession(): boolean;
    clearReadModeScSession(): void;
    getReadModeScSession(): number;
    setReadModeScSession(value: number): void;

    hasReadModeScLinearize(): boolean;
    clearReadModeScLinearize(): void;
    getReadModeScLinearize(): number;
    setReadModeScLinearize(value: number): void;

    hasReadModeScReplica(): boolean;
    clearReadModeScReplica(): void;
    getReadModeScReplica(): number;
    setReadModeScReplica(value: number): void;

    hasReadModeScAllowUnavailable(): boolean;
    clearReadModeScAllowUnavailable(): void;
    getReadModeScAllowUnavailable(): number;
    setReadModeScAllowUnavailable(value: number): void;

    getReadModeScCase(): readModeSc.ReadModeScCase;
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): readModeSc.AsObject;
    static toObject(includeInstance: boolean, msg: readModeSc): readModeSc.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: readModeSc, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): readModeSc;
    static deserializeBinaryFromReader(message: readModeSc, reader: jspb.BinaryReader): readModeSc;
  }

  export namespace readModeSc {
    export type AsObject = {
      readModeScSession: number,
      readModeScLinearize: number,
      readModeScReplica: number,
      readModeScAllowUnavailable: number,
    }

    export enum ReadModeScCase {
      READ_MODE_SC_NOT_SET = 0,
      READ_MODE_SC_SESSION = 1,
      READ_MODE_SC_LINEARIZE = 2,
      READ_MODE_SC_REPLICA = 3,
      READ_MODE_SC_ALLOW_UNAVAILABLE = 4,
    }
  }

  export class readModeAp extends jspb.Message {
    hasReadModeApOne(): boolean;
    clearReadModeApOne(): void;
    getReadModeApOne(): number;
    setReadModeApOne(value: number): void;

    hasReadModeApAll(): boolean;
    clearReadModeApAll(): void;
    getReadModeApAll(): number;
    setReadModeApAll(value: number): void;

    getReadModeApCase(): readModeAp.ReadModeApCase;
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): readModeAp.AsObject;
    static toObject(includeInstance: boolean, msg: readModeAp): readModeAp.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: readModeAp, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): readModeAp;
    static deserializeBinaryFromReader(message: readModeAp, reader: jspb.BinaryReader): readModeAp;
  }

  export namespace readModeAp {
    export type AsObject = {
      readModeApOne: number,
      readModeApAll: number,
    }

    export enum ReadModeApCase {
      READ_MODE_AP_NOT_SET = 0,
      READ_MODE_AP_ONE = 1,
      READ_MODE_AP_ALL = 2,
    }
  }

  export class tlsCurveID extends jspb.Message {
    hasCurveP256(): boolean;
    clearCurveP256(): void;
    getCurveP256(): number;
    setCurveP256(value: number): void;

    hasCurveP384(): boolean;
    clearCurveP384(): void;
    getCurveP384(): number;
    setCurveP384(value: number): void;

    hasCurveP521(): boolean;
    clearCurveP521(): void;
    getCurveP521(): number;
    setCurveP521(value: number): void;

    hasX25519(): boolean;
    clearX25519(): void;
    getX25519(): number;
    setX25519(value: number): void;

    getCurveIdCase(): tlsCurveID.CurveIdCase;
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): tlsCurveID.AsObject;
    static toObject(includeInstance: boolean, msg: tlsCurveID): tlsCurveID.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: tlsCurveID, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): tlsCurveID;
    static deserializeBinaryFromReader(message: tlsCurveID, reader: jspb.BinaryReader): tlsCurveID;
  }

  export namespace tlsCurveID {
    export type AsObject = {
      curveP256: number,
      curveP384: number,
      curveP521: number,
      x25519: number,
    }

    export enum CurveIdCase {
      CURVE_ID_NOT_SET = 0,
      CURVE_P256 = 1,
      CURVE_P384 = 2,
      CURVE_P521 = 3,
      X_25519 = 4,
    }
  }

  export enum CommitLevelCase {
    COMMIT_LEVEL_NOT_SET = 0,
    COMMIT_ALL = 6,
    COMMIT_MASTER = 7,
  }
}

export class ApiKey extends jspb.Message {
  getApiKey(): string;
  setApiKey(value: string): void;

  clearLabelsList(): void;
  getLabelsList(): Array<string>;
  setLabelsList(value: Array<string>): void;
  addLabels(value: string, index?: number): string;

  getMetadataMap(): jspb.Map<string, string>;
  clearMetadataMap(): void;
  getUuid(): string;
  setUuid(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ApiKey.AsObject;
  static toObject(includeInstance: boolean, msg: ApiKey): ApiKey.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ApiKey, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ApiKey;
  static deserializeBinaryFromReader(message: ApiKey, reader: jspb.BinaryReader): ApiKey;
}

export namespace ApiKey {
  export type AsObject = {
    apiKey: string,
    labelsList: Array<string>,
    metadataMap: Array<[string, string]>,
    uuid: string,
  }
}

export class ApiKeySecret extends jspb.Message {
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

  getDisableGroupChecking(): boolean;
  setDisableGroupChecking(value: boolean): void;

  hasGroupLookupSettings(): boolean;
  clearGroupLookupSettings(): void;
  getGroupLookupSettings(): LdapServiceAccount | undefined;
  setGroupLookupSettings(value?: LdapServiceAccount): void;

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
    disableGroupChecking: boolean,
    groupLookupSettings?: LdapServiceAccount.AsObject,
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

export class LdapServiceAccount extends jspb.Message {
  hasCredentialsSecretRef(): boolean;
  clearCredentialsSecretRef(): void;
  getCredentialsSecretRef(): github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef | undefined;
  setCredentialsSecretRef(value?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef): void;

  getCheckGroupsWithServiceAccount(): boolean;
  setCheckGroupsWithServiceAccount(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LdapServiceAccount.AsObject;
  static toObject(includeInstance: boolean, msg: LdapServiceAccount): LdapServiceAccount.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: LdapServiceAccount, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LdapServiceAccount;
  static deserializeBinaryFromReader(message: LdapServiceAccount, reader: jspb.BinaryReader): LdapServiceAccount;
}

export namespace LdapServiceAccount {
  export type AsObject = {
    credentialsSecretRef?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject,
    checkGroupsWithServiceAccount: boolean,
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

  getFailureModeAllow(): boolean;
  setFailureModeAllow(value: boolean): void;

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
    failureModeAllow: boolean,
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

  hasTlsconfig(): boolean;
  clearTlsconfig(): void;
  getTlsconfig(): PassThroughGrpcTLSConfig | undefined;
  setTlsconfig(value?: PassThroughGrpcTLSConfig): void;

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
    tlsconfig?: PassThroughGrpcTLSConfig.AsObject,
  }
}

export class PassThroughGrpcTLSConfig extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PassThroughGrpcTLSConfig.AsObject;
  static toObject(includeInstance: boolean, msg: PassThroughGrpcTLSConfig): PassThroughGrpcTLSConfig.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: PassThroughGrpcTLSConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PassThroughGrpcTLSConfig;
  static deserializeBinaryFromReader(message: PassThroughGrpcTLSConfig, reader: jspb.BinaryReader): PassThroughGrpcTLSConfig;
}

export namespace PassThroughGrpcTLSConfig {
  export type AsObject = {
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

  getFailOnRedirect(): boolean;
  setFailOnRedirect(value: boolean): void;

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
    failOnRedirect: boolean,
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

  export class UserSessionConfig extends jspb.Message {
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

    hasCipherConfig(): boolean;
    clearCipherConfig(): void;
    getCipherConfig(): ExtAuthConfig.UserSessionConfig.CipherConfig | undefined;
    setCipherConfig(value?: ExtAuthConfig.UserSessionConfig.CipherConfig): void;

    getSessionCase(): UserSessionConfig.SessionCase;
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): UserSessionConfig.AsObject;
    static toObject(includeInstance: boolean, msg: UserSessionConfig): UserSessionConfig.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: UserSessionConfig, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): UserSessionConfig;
    static deserializeBinaryFromReader(message: UserSessionConfig, reader: jspb.BinaryReader): UserSessionConfig;
  }

  export namespace UserSessionConfig {
    export type AsObject = {
      failOnFetchFailure: boolean,
      cookieOptions?: UserSession.CookieOptions.AsObject,
      cookie?: UserSession.InternalSession.AsObject,
      redis?: UserSession.RedisSession.AsObject,
      cipherConfig?: ExtAuthConfig.UserSessionConfig.CipherConfig.AsObject,
    }

    export class CipherConfig extends jspb.Message {
      getKey(): string;
      setKey(value: string): void;

      serializeBinary(): Uint8Array;
      toObject(includeInstance?: boolean): CipherConfig.AsObject;
      static toObject(includeInstance: boolean, msg: CipherConfig): CipherConfig.AsObject;
      static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
      static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
      static serializeBinaryToWriter(message: CipherConfig, writer: jspb.BinaryWriter): void;
      static deserializeBinary(bytes: Uint8Array): CipherConfig;
      static deserializeBinaryFromReader(message: CipherConfig, reader: jspb.BinaryReader): CipherConfig;
    }

    export namespace CipherConfig {
      export type AsObject = {
        key: string,
      }
    }

    export enum SessionCase {
      SESSION_NOT_SET = 0,
      COOKIE = 3,
      REDIS = 4,
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

    hasEndSessionProperties(): boolean;
    clearEndSessionProperties(): void;
    getEndSessionProperties(): EndSessionProperties | undefined;
    setEndSessionProperties(value?: EndSessionProperties): void;

    hasUserSession(): boolean;
    clearUserSession(): void;
    getUserSession(): ExtAuthConfig.UserSessionConfig | undefined;
    setUserSession(value?: ExtAuthConfig.UserSessionConfig): void;

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
      endSessionProperties?: EndSessionProperties.AsObject,
      userSession?: ExtAuthConfig.UserSessionConfig.AsObject,
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

  export class PlainOAuth2Config extends jspb.Message {
    getClientId(): string;
    setClientId(value: string): void;

    getClientSecret(): string;
    setClientSecret(value: string): void;

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

    hasSession(): boolean;
    clearSession(): void;
    getSession(): UserSession | undefined;
    setSession(value?: UserSession): void;

    getLogoutPath(): string;
    setLogoutPath(value: string): void;

    getTokenEndpointQueryParamsMap(): jspb.Map<string, string>;
    clearTokenEndpointQueryParamsMap(): void;
    getAfterLogoutUrl(): string;
    setAfterLogoutUrl(value: string): void;

    getAuthEndpoint(): string;
    setAuthEndpoint(value: string): void;

    getTokenEndpoint(): string;
    setTokenEndpoint(value: string): void;

    getRevocationEndpoint(): string;
    setRevocationEndpoint(value: string): void;

    hasUserSession(): boolean;
    clearUserSession(): void;
    getUserSession(): ExtAuthConfig.UserSessionConfig | undefined;
    setUserSession(value?: ExtAuthConfig.UserSessionConfig): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): PlainOAuth2Config.AsObject;
    static toObject(includeInstance: boolean, msg: PlainOAuth2Config): PlainOAuth2Config.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: PlainOAuth2Config, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): PlainOAuth2Config;
    static deserializeBinaryFromReader(message: PlainOAuth2Config, reader: jspb.BinaryReader): PlainOAuth2Config;
  }

  export namespace PlainOAuth2Config {
    export type AsObject = {
      clientId: string,
      clientSecret: string,
      authEndpointQueryParamsMap: Array<[string, string]>,
      appUrl: string,
      callbackPath: string,
      scopesList: Array<string>,
      session?: UserSession.AsObject,
      logoutPath: string,
      tokenEndpointQueryParamsMap: Array<[string, string]>,
      afterLogoutUrl: string,
      authEndpoint: string,
      tokenEndpoint: string,
      revocationEndpoint: string,
      userSession?: ExtAuthConfig.UserSessionConfig.AsObject,
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

    hasOauth2Config(): boolean;
    clearOauth2Config(): void;
    getOauth2Config(): ExtAuthConfig.PlainOAuth2Config | undefined;
    setOauth2Config(value?: ExtAuthConfig.PlainOAuth2Config): void;

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
      oauth2Config?: ExtAuthConfig.PlainOAuth2Config.AsObject,
    }

    export enum OauthTypeCase {
      OAUTH_TYPE_NOT_SET = 0,
      OIDC_AUTHORIZATION_CODE = 1,
      ACCESS_TOKEN_VALIDATION_CONFIG = 3,
      OAUTH2_CONFIG = 4,
    }
  }

  export class ApiKeyAuthConfig extends jspb.Message {
    getValidApiKeysMap(): jspb.Map<string, ExtAuthConfig.ApiKeyAuthConfig.KeyMetadata>;
    clearValidApiKeysMap(): void;
    getHeaderName(): string;
    setHeaderName(value: string): void;

    getHeadersFromKeyMetadataMap(): jspb.Map<string, string>;
    clearHeadersFromKeyMetadataMap(): void;
    hasK8sSecretApikeyStorage(): boolean;
    clearK8sSecretApikeyStorage(): void;
    getK8sSecretApikeyStorage(): K8sSecretApiKeyStorage | undefined;
    setK8sSecretApikeyStorage(value?: K8sSecretApiKeyStorage): void;

    hasAerospikeApikeyStorage(): boolean;
    clearAerospikeApikeyStorage(): void;
    getAerospikeApikeyStorage(): AerospikeApiKeyStorage | undefined;
    setAerospikeApikeyStorage(value?: AerospikeApiKeyStorage): void;

    getStorageBackendCase(): ApiKeyAuthConfig.StorageBackendCase;
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
      k8sSecretApikeyStorage?: K8sSecretApiKeyStorage.AsObject,
      aerospikeApikeyStorage?: AerospikeApiKeyStorage.AsObject,
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

    export enum StorageBackendCase {
      STORAGE_BACKEND_NOT_SET = 0,
      K8S_SECRET_APIKEY_STORAGE = 4,
      AEROSPIKE_APIKEY_STORAGE = 5,
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

  export class LdapConfig extends jspb.Message {
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

    getDisableGroupChecking(): boolean;
    setDisableGroupChecking(value: boolean): void;

    hasGroupLookupSettings(): boolean;
    clearGroupLookupSettings(): void;
    getGroupLookupSettings(): ExtAuthConfig.LdapServiceAccountConfig | undefined;
    setGroupLookupSettings(value?: ExtAuthConfig.LdapServiceAccountConfig): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): LdapConfig.AsObject;
    static toObject(includeInstance: boolean, msg: LdapConfig): LdapConfig.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: LdapConfig, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): LdapConfig;
    static deserializeBinaryFromReader(message: LdapConfig, reader: jspb.BinaryReader): LdapConfig;
  }

  export namespace LdapConfig {
    export type AsObject = {
      address: string,
      userdntemplate: string,
      membershipattributename: string,
      allowedgroupsList: Array<string>,
      pool?: Ldap.ConnectionPool.AsObject,
      searchfilter: string,
      disableGroupChecking: boolean,
      groupLookupSettings?: ExtAuthConfig.LdapServiceAccountConfig.AsObject,
    }
  }

  export class LdapServiceAccountConfig extends jspb.Message {
    getUsername(): string;
    setUsername(value: string): void;

    getPassword(): string;
    setPassword(value: string): void;

    getCheckGroupsWithServiceAccount(): boolean;
    setCheckGroupsWithServiceAccount(value: boolean): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): LdapServiceAccountConfig.AsObject;
    static toObject(includeInstance: boolean, msg: LdapServiceAccountConfig): LdapServiceAccountConfig.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: LdapServiceAccountConfig, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): LdapServiceAccountConfig;
    static deserializeBinaryFromReader(message: LdapServiceAccountConfig, reader: jspb.BinaryReader): LdapServiceAccountConfig;
  }

  export namespace LdapServiceAccountConfig {
    export type AsObject = {
      username: string,
      password: string,
      checkGroupsWithServiceAccount: boolean,
    }
  }

  export class HmacAuthConfig extends jspb.Message {
    hasSecretList(): boolean;
    clearSecretList(): void;
    getSecretList(): ExtAuthConfig.InMemorySecretList | undefined;
    setSecretList(value?: ExtAuthConfig.InMemorySecretList): void;

    hasParametersInHeaders(): boolean;
    clearParametersInHeaders(): void;
    getParametersInHeaders(): HmacParametersInHeaders | undefined;
    setParametersInHeaders(value?: HmacParametersInHeaders): void;

    getSecretStorageCase(): HmacAuthConfig.SecretStorageCase;
    getImplementationTypeCase(): HmacAuthConfig.ImplementationTypeCase;
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): HmacAuthConfig.AsObject;
    static toObject(includeInstance: boolean, msg: HmacAuthConfig): HmacAuthConfig.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: HmacAuthConfig, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): HmacAuthConfig;
    static deserializeBinaryFromReader(message: HmacAuthConfig, reader: jspb.BinaryReader): HmacAuthConfig;
  }

  export namespace HmacAuthConfig {
    export type AsObject = {
      secretList?: ExtAuthConfig.InMemorySecretList.AsObject,
      parametersInHeaders?: HmacParametersInHeaders.AsObject,
    }

    export enum SecretStorageCase {
      SECRET_STORAGE_NOT_SET = 0,
      SECRET_LIST = 1,
    }

    export enum ImplementationTypeCase {
      IMPLEMENTATION_TYPE_NOT_SET = 0,
      PARAMETERS_IN_HEADERS = 2,
    }
  }

  export class InMemorySecretList extends jspb.Message {
    getSecretListMap(): jspb.Map<string, string>;
    clearSecretListMap(): void;
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): InMemorySecretList.AsObject;
    static toObject(includeInstance: boolean, msg: InMemorySecretList): InMemorySecretList.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: InMemorySecretList, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): InMemorySecretList;
    static deserializeBinaryFromReader(message: InMemorySecretList, reader: jspb.BinaryReader): InMemorySecretList;
  }

  export namespace InMemorySecretList {
    export type AsObject = {
      secretListMap: Array<[string, string]>,
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

    hasLdapInternal(): boolean;
    clearLdapInternal(): void;
    getLdapInternal(): ExtAuthConfig.LdapConfig | undefined;
    setLdapInternal(value?: ExtAuthConfig.LdapConfig): void;

    hasJwt(): boolean;
    clearJwt(): void;
    getJwt(): google_protobuf_empty_pb.Empty | undefined;
    setJwt(value?: google_protobuf_empty_pb.Empty): void;

    hasPassThroughAuth(): boolean;
    clearPassThroughAuth(): void;
    getPassThroughAuth(): PassThroughAuth | undefined;
    setPassThroughAuth(value?: PassThroughAuth): void;

    hasHmacAuth(): boolean;
    clearHmacAuth(): void;
    getHmacAuth(): ExtAuthConfig.HmacAuthConfig | undefined;
    setHmacAuth(value?: ExtAuthConfig.HmacAuthConfig): void;

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
      ldapInternal?: ExtAuthConfig.LdapConfig.AsObject,
      jwt?: google_protobuf_empty_pb.Empty.AsObject,
      passThroughAuth?: PassThroughAuth.AsObject,
      hmacAuth?: ExtAuthConfig.HmacAuthConfig.AsObject,
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
      LDAP_INTERNAL = 14,
      JWT = 12,
      PASS_THROUGH_AUTH = 13,
      HMAC_AUTH = 15,
    }
  }
}

export class ApiKeyCreateRequest extends jspb.Message {
  clearApiKeysList(): void;
  getApiKeysList(): Array<ApiKey>;
  setApiKeysList(value: Array<ApiKey>): void;
  addApiKeys(value?: ApiKey, index?: number): ApiKey;

  clearRawApiKeysList(): void;
  getRawApiKeysList(): Array<string>;
  setRawApiKeysList(value: Array<string>): void;
  addRawApiKeys(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ApiKeyCreateRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ApiKeyCreateRequest): ApiKeyCreateRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ApiKeyCreateRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ApiKeyCreateRequest;
  static deserializeBinaryFromReader(message: ApiKeyCreateRequest, reader: jspb.BinaryReader): ApiKeyCreateRequest;
}

export namespace ApiKeyCreateRequest {
  export type AsObject = {
    apiKeysList: Array<ApiKey.AsObject>,
    rawApiKeysList: Array<string>,
  }
}

export class ApiKeyCreateResponse extends jspb.Message {
  clearApiKeysList(): void;
  getApiKeysList(): Array<ApiKey>;
  setApiKeysList(value: Array<ApiKey>): void;
  addApiKeys(value?: ApiKey, index?: number): ApiKey;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ApiKeyCreateResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ApiKeyCreateResponse): ApiKeyCreateResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ApiKeyCreateResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ApiKeyCreateResponse;
  static deserializeBinaryFromReader(message: ApiKeyCreateResponse, reader: jspb.BinaryReader): ApiKeyCreateResponse;
}

export namespace ApiKeyCreateResponse {
  export type AsObject = {
    apiKeysList: Array<ApiKey.AsObject>,
  }
}

export class ApiKeyReadRequest extends jspb.Message {
  clearRawApiKeysList(): void;
  getRawApiKeysList(): Array<string>;
  setRawApiKeysList(value: Array<string>): void;
  addRawApiKeys(value: string, index?: number): string;

  clearLabelsList(): void;
  getLabelsList(): Array<string>;
  setLabelsList(value: Array<string>): void;
  addLabels(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ApiKeyReadRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ApiKeyReadRequest): ApiKeyReadRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ApiKeyReadRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ApiKeyReadRequest;
  static deserializeBinaryFromReader(message: ApiKeyReadRequest, reader: jspb.BinaryReader): ApiKeyReadRequest;
}

export namespace ApiKeyReadRequest {
  export type AsObject = {
    rawApiKeysList: Array<string>,
    labelsList: Array<string>,
  }
}

export class ApiKeyReadResponse extends jspb.Message {
  clearApiKeysList(): void;
  getApiKeysList(): Array<ApiKey>;
  setApiKeysList(value: Array<ApiKey>): void;
  addApiKeys(value?: ApiKey, index?: number): ApiKey;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ApiKeyReadResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ApiKeyReadResponse): ApiKeyReadResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ApiKeyReadResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ApiKeyReadResponse;
  static deserializeBinaryFromReader(message: ApiKeyReadResponse, reader: jspb.BinaryReader): ApiKeyReadResponse;
}

export namespace ApiKeyReadResponse {
  export type AsObject = {
    apiKeysList: Array<ApiKey.AsObject>,
  }
}

export class ApiKeyUpdateRequest extends jspb.Message {
  getUpsert(): boolean;
  setUpsert(value: boolean): void;

  clearApiKeysList(): void;
  getApiKeysList(): Array<ApiKey>;
  setApiKeysList(value: Array<ApiKey>): void;
  addApiKeys(value?: ApiKey, index?: number): ApiKey;

  clearRawApiKeysList(): void;
  getRawApiKeysList(): Array<string>;
  setRawApiKeysList(value: Array<string>): void;
  addRawApiKeys(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ApiKeyUpdateRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ApiKeyUpdateRequest): ApiKeyUpdateRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ApiKeyUpdateRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ApiKeyUpdateRequest;
  static deserializeBinaryFromReader(message: ApiKeyUpdateRequest, reader: jspb.BinaryReader): ApiKeyUpdateRequest;
}

export namespace ApiKeyUpdateRequest {
  export type AsObject = {
    upsert: boolean,
    apiKeysList: Array<ApiKey.AsObject>,
    rawApiKeysList: Array<string>,
  }
}

export class ApiKeyUpdateResponse extends jspb.Message {
  clearApiKeysList(): void;
  getApiKeysList(): Array<ApiKey>;
  setApiKeysList(value: Array<ApiKey>): void;
  addApiKeys(value?: ApiKey, index?: number): ApiKey;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ApiKeyUpdateResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ApiKeyUpdateResponse): ApiKeyUpdateResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ApiKeyUpdateResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ApiKeyUpdateResponse;
  static deserializeBinaryFromReader(message: ApiKeyUpdateResponse, reader: jspb.BinaryReader): ApiKeyUpdateResponse;
}

export namespace ApiKeyUpdateResponse {
  export type AsObject = {
    apiKeysList: Array<ApiKey.AsObject>,
  }
}

export class ApiKeyDeleteRequest extends jspb.Message {
  clearRawApiKeysList(): void;
  getRawApiKeysList(): Array<string>;
  setRawApiKeysList(value: Array<string>): void;
  addRawApiKeys(value: string, index?: number): string;

  clearLabelsList(): void;
  getLabelsList(): Array<string>;
  setLabelsList(value: Array<string>): void;
  addLabels(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ApiKeyDeleteRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ApiKeyDeleteRequest): ApiKeyDeleteRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ApiKeyDeleteRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ApiKeyDeleteRequest;
  static deserializeBinaryFromReader(message: ApiKeyDeleteRequest, reader: jspb.BinaryReader): ApiKeyDeleteRequest;
}

export namespace ApiKeyDeleteRequest {
  export type AsObject = {
    rawApiKeysList: Array<string>,
    labelsList: Array<string>,
  }
}

export class ApiKeyDeleteResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ApiKeyDeleteResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ApiKeyDeleteResponse): ApiKeyDeleteResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ApiKeyDeleteResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ApiKeyDeleteResponse;
  static deserializeBinaryFromReader(message: ApiKeyDeleteResponse, reader: jspb.BinaryReader): ApiKeyDeleteResponse;
}

export namespace ApiKeyDeleteResponse {
  export type AsObject = {
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
