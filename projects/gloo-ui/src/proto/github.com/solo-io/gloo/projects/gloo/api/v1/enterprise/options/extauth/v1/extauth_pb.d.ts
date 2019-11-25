// package: enterprise.gloo.solo.io
// file: github.com/solo-io/gloo/projects/gloo/api/v1/enterprise/options/extauth/v1/extauth.proto

import * as jspb from "google-protobuf";
import * as github_com_solo_io_solo_kit_api_v1_ref_pb from "../../../../../../../../../../../github.com/solo-io/solo-kit/api/v1/ref_pb";
import * as gogoproto_gogo_pb from "../../../../../../../../../../../gogoproto/gogo_pb";
import * as github_com_solo_io_solo_kit_api_v1_metadata_pb from "../../../../../../../../../../../github.com/solo-io/solo-kit/api/v1/metadata_pb";
import * as github_com_solo_io_solo_kit_api_v1_status_pb from "../../../../../../../../../../../github.com/solo-io/solo-kit/api/v1/status_pb";
import * as github_com_solo_io_solo_kit_api_v1_solo_kit_pb from "../../../../../../../../../../../github.com/solo-io/solo-kit/api/v1/solo-kit_pb";
import * as envoy_api_v2_discovery_pb from "../../../../../../../../../../../envoy/api/v2/discovery_pb";
import * as google_api_annotations_pb from "../../../../../../../../../../../google/api/annotations_pb";
import * as google_protobuf_duration_pb from "google-protobuf/google/protobuf/duration_pb";
import * as google_protobuf_struct_pb from "google-protobuf/google/protobuf/struct_pb";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";

export class AuthConfig extends jspb.Message {
  hasStatus(): boolean;
  clearStatus(): void;
  getStatus(): github_com_solo_io_solo_kit_api_v1_status_pb.Status | undefined;
  setStatus(value?: github_com_solo_io_solo_kit_api_v1_status_pb.Status): void;

  hasMetadata(): boolean;
  clearMetadata(): void;
  getMetadata(): github_com_solo_io_solo_kit_api_v1_metadata_pb.Metadata | undefined;
  setMetadata(value?: github_com_solo_io_solo_kit_api_v1_metadata_pb.Metadata): void;

  clearConfigsList(): void;
  getConfigsList(): Array<AuthConfig.Config>;
  setConfigsList(value: Array<AuthConfig.Config>): void;
  addConfigs(value?: AuthConfig.Config, index?: number): AuthConfig.Config;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AuthConfig.AsObject;
  static toObject(includeInstance: boolean, msg: AuthConfig): AuthConfig.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: AuthConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AuthConfig;
  static deserializeBinaryFromReader(message: AuthConfig, reader: jspb.BinaryReader): AuthConfig;
}

export namespace AuthConfig {
  export type AsObject = {
    status?: github_com_solo_io_solo_kit_api_v1_status_pb.Status.AsObject,
    metadata?: github_com_solo_io_solo_kit_api_v1_metadata_pb.Metadata.AsObject,
    configsList: Array<AuthConfig.Config.AsObject>,
  }

  export class Config extends jspb.Message {
    hasBasicAuth(): boolean;
    clearBasicAuth(): void;
    getBasicAuth(): BasicAuth | undefined;
    setBasicAuth(value?: BasicAuth): void;

    hasOauth(): boolean;
    clearOauth(): void;
    getOauth(): OAuth | undefined;
    setOauth(value?: OAuth): void;

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
      basicAuth?: BasicAuth.AsObject,
      oauth?: OAuth.AsObject,
      apiKeyAuth?: ApiKeyAuth.AsObject,
      pluginAuth?: AuthPlugin.AsObject,
      opaAuth?: OpaAuth.AsObject,
      ldap?: Ldap.AsObject,
    }

    export enum AuthConfigCase {
      AUTH_CONFIG_NOT_SET = 0,
      BASIC_AUTH = 1,
      OAUTH = 2,
      API_KEY_AUTH = 4,
      PLUGIN_AUTH = 5,
      OPA_AUTH = 6,
      LDAP = 7,
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
    userIdHeader: string,
    requestTimeout?: google_protobuf_duration_pb.Duration.AsObject,
    failureModeAllow: boolean,
    requestBody?: BufferSettings.AsObject,
    clearRouteCache: boolean,
    statusOnError: number,
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
    }
  }
}

export class BufferSettings extends jspb.Message {
  getMaxRequestBytes(): number;
  setMaxRequestBytes(value: number): void;

  getAllowPartialMessage(): boolean;
  setAllowPartialMessage(value: boolean): void;

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
  }
}

export class CustomAuth extends jspb.Message {
  getContextExtensionsMap(): jspb.Map<string, string>;
  clearContextExtensionsMap(): void;
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
    appUrl: string,
    callbackPath: string,
    scopesList: Array<string>,
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
  }
}

export class OpaAuth extends jspb.Message {
  clearModulesList(): void;
  getModulesList(): Array<github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef>;
  setModulesList(value: Array<github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef>): void;
  addModules(value?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef, index?: number): github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef;

  getQuery(): string;
  setQuery(value: string): void;

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

export class ExtAuthConfig extends jspb.Message {
  getAuthConfigRefName(): string;
  setAuthConfigRefName(value: string): void;

  clearConfigsList(): void;
  getConfigsList(): Array<ExtAuthConfig.Config>;
  setConfigsList(value: Array<ExtAuthConfig.Config>): void;
  addConfigs(value?: ExtAuthConfig.Config, index?: number): ExtAuthConfig.Config;

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
  }

  export class OAuthConfig extends jspb.Message {
    getClientId(): string;
    setClientId(value: string): void;

    getClientSecret(): string;
    setClientSecret(value: string): void;

    getIssuerUrl(): string;
    setIssuerUrl(value: string): void;

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
      appUrl: string,
      callbackPath: string,
      scopesList: Array<string>,
    }
  }

  export class ApiKeyAuthConfig extends jspb.Message {
    getValidApiKeyAndUserMap(): jspb.Map<string, string>;
    clearValidApiKeyAndUserMap(): void;
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
      validApiKeyAndUserMap: Array<[string, string]>,
    }
  }

  export class OpaAuthConfig extends jspb.Message {
    getModulesMap(): jspb.Map<string, string>;
    clearModulesMap(): void;
    getQuery(): string;
    setQuery(value: string): void;

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
    }
  }

  export class Config extends jspb.Message {
    hasOauth(): boolean;
    clearOauth(): void;
    getOauth(): ExtAuthConfig.OAuthConfig | undefined;
    setOauth(value?: ExtAuthConfig.OAuthConfig): void;

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
      oauth?: ExtAuthConfig.OAuthConfig.AsObject,
      basicAuth?: BasicAuth.AsObject,
      apiKeyAuth?: ExtAuthConfig.ApiKeyAuthConfig.AsObject,
      pluginAuth?: AuthPlugin.AsObject,
      opaAuth?: ExtAuthConfig.OpaAuthConfig.AsObject,
      ldap?: Ldap.AsObject,
    }

    export enum AuthConfigCase {
      AUTH_CONFIG_NOT_SET = 0,
      OAUTH = 3,
      BASIC_AUTH = 4,
      API_KEY_AUTH = 5,
      PLUGIN_AUTH = 6,
      OPA_AUTH = 7,
      LDAP = 8,
    }
  }
}

