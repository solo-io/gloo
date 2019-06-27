// package: extauth.plugins.gloo.solo.io
// file: github.com/solo-io/solo-projects/projects/gloo/api/v1/plugins/extauth/extauth.proto

import * as jspb from "google-protobuf";
import * as github_com_solo_io_solo_kit_api_v1_ref_pb from "../../../../../../../../../github.com/solo-io/solo-kit/api/v1/ref_pb";
import * as gogoproto_gogo_pb from "../../../../../../../../../gogoproto/gogo_pb";
import * as envoy_api_v2_discovery_pb from "../../../../../../../../../envoy/api/v2/discovery_pb";
import * as google_api_annotations_pb from "../../../../../../../../../google/api/annotations_pb";
import * as google_protobuf_duration_pb from "google-protobuf/google/protobuf/duration_pb";

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

export class VhostExtension extends jspb.Message {
  hasBasicAuth(): boolean;
  clearBasicAuth(): void;
  getBasicAuth(): BasicAuth | undefined;
  setBasicAuth(value?: BasicAuth): void;

  hasOauth(): boolean;
  clearOauth(): void;
  getOauth(): OAuth | undefined;
  setOauth(value?: OAuth): void;

  hasCustomAuth(): boolean;
  clearCustomAuth(): void;
  getCustomAuth(): CustomAuth | undefined;
  setCustomAuth(value?: CustomAuth): void;

  getAuthConfigCase(): VhostExtension.AuthConfigCase;
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
    basicAuth?: BasicAuth.AsObject,
    oauth?: OAuth.AsObject,
    customAuth?: CustomAuth.AsObject,
  }

  export enum AuthConfigCase {
    AUTH_CONFIG_NOT_SET = 0,
    BASIC_AUTH = 1,
    OAUTH = 2,
    CUSTOM_AUTH = 3,
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

export class ExtAuthConfig extends jspb.Message {
  getVhost(): string;
  setVhost(value: string): void;

  hasOauth(): boolean;
  clearOauth(): void;
  getOauth(): ExtAuthConfig.OAuthConfig | undefined;
  setOauth(value?: ExtAuthConfig.OAuthConfig): void;

  hasBasicAuth(): boolean;
  clearBasicAuth(): void;
  getBasicAuth(): BasicAuth | undefined;
  setBasicAuth(value?: BasicAuth): void;

  getAuthConfigCase(): ExtAuthConfig.AuthConfigCase;
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
    vhost: string,
    oauth?: ExtAuthConfig.OAuthConfig.AsObject,
    basicAuth?: BasicAuth.AsObject,
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
    }
  }

  export enum AuthConfigCase {
    AUTH_CONFIG_NOT_SET = 0,
    OAUTH = 3,
    BASIC_AUTH = 4,
  }
}

