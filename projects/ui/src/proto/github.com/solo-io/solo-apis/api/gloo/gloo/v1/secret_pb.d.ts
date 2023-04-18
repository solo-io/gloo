/* eslint-disable */
// package: gloo.solo.io
// file: github.com/solo-io/solo-apis/api/gloo/gloo/v1/secret.proto

import * as jspb from "google-protobuf";
import * as extproto_ext_pb from "../../../../../../../extproto/ext_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_extensions_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/extensions_pb";
import * as github_com_solo_io_solo_apis_api_gloo_enterprise_gloo_v1_auth_config_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/enterprise.gloo/v1/auth_config_pb";
import * as github_com_solo_io_solo_kit_api_v1_metadata_pb from "../../../../../../../github.com/solo-io/solo-kit/api/v1/metadata_pb";
import * as github_com_solo_io_solo_kit_api_v1_solo_kit_pb from "../../../../../../../github.com/solo-io/solo-kit/api/v1/solo-kit_pb";

export class Secret extends jspb.Message {
  hasAws(): boolean;
  clearAws(): void;
  getAws(): AwsSecret | undefined;
  setAws(value?: AwsSecret): void;

  hasAzure(): boolean;
  clearAzure(): void;
  getAzure(): AzureSecret | undefined;
  setAzure(value?: AzureSecret): void;

  hasTls(): boolean;
  clearTls(): void;
  getTls(): TlsSecret | undefined;
  setTls(value?: TlsSecret): void;

  hasOauth(): boolean;
  clearOauth(): void;
  getOauth(): github_com_solo_io_solo_apis_api_gloo_enterprise_gloo_v1_auth_config_pb.OauthSecret | undefined;
  setOauth(value?: github_com_solo_io_solo_apis_api_gloo_enterprise_gloo_v1_auth_config_pb.OauthSecret): void;

  hasApiKey(): boolean;
  clearApiKey(): void;
  getApiKey(): github_com_solo_io_solo_apis_api_gloo_enterprise_gloo_v1_auth_config_pb.ApiKey | undefined;
  setApiKey(value?: github_com_solo_io_solo_apis_api_gloo_enterprise_gloo_v1_auth_config_pb.ApiKey): void;

  hasHeader(): boolean;
  clearHeader(): void;
  getHeader(): HeaderSecret | undefined;
  setHeader(value?: HeaderSecret): void;

  hasCredentials(): boolean;
  clearCredentials(): void;
  getCredentials(): AccountCredentialsSecret | undefined;
  setCredentials(value?: AccountCredentialsSecret): void;

  hasEncryption(): boolean;
  clearEncryption(): void;
  getEncryption(): EncryptionKeySecret | undefined;
  setEncryption(value?: EncryptionKeySecret): void;

  hasExtensions(): boolean;
  clearExtensions(): void;
  getExtensions(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_extensions_pb.Extensions | undefined;
  setExtensions(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_extensions_pb.Extensions): void;

  hasMetadata(): boolean;
  clearMetadata(): void;
  getMetadata(): github_com_solo_io_solo_kit_api_v1_metadata_pb.Metadata | undefined;
  setMetadata(value?: github_com_solo_io_solo_kit_api_v1_metadata_pb.Metadata): void;

  getKindCase(): Secret.KindCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Secret.AsObject;
  static toObject(includeInstance: boolean, msg: Secret): Secret.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Secret, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Secret;
  static deserializeBinaryFromReader(message: Secret, reader: jspb.BinaryReader): Secret;
}

export namespace Secret {
  export type AsObject = {
    aws?: AwsSecret.AsObject,
    azure?: AzureSecret.AsObject,
    tls?: TlsSecret.AsObject,
    oauth?: github_com_solo_io_solo_apis_api_gloo_enterprise_gloo_v1_auth_config_pb.OauthSecret.AsObject,
    apiKey?: github_com_solo_io_solo_apis_api_gloo_enterprise_gloo_v1_auth_config_pb.ApiKey.AsObject,
    header?: HeaderSecret.AsObject,
    credentials?: AccountCredentialsSecret.AsObject,
    encryption?: EncryptionKeySecret.AsObject,
    extensions?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_extensions_pb.Extensions.AsObject,
    metadata?: github_com_solo_io_solo_kit_api_v1_metadata_pb.Metadata.AsObject,
  }

  export enum KindCase {
    KIND_NOT_SET = 0,
    AWS = 1,
    AZURE = 2,
    TLS = 3,
    OAUTH = 5,
    API_KEY = 6,
    HEADER = 8,
    CREDENTIALS = 9,
    ENCRYPTION = 10,
    EXTENSIONS = 4,
  }
}

export class AwsSecret extends jspb.Message {
  getAccessKey(): string;
  setAccessKey(value: string): void;

  getSecretKey(): string;
  setSecretKey(value: string): void;

  getSessionToken(): string;
  setSessionToken(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AwsSecret.AsObject;
  static toObject(includeInstance: boolean, msg: AwsSecret): AwsSecret.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: AwsSecret, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AwsSecret;
  static deserializeBinaryFromReader(message: AwsSecret, reader: jspb.BinaryReader): AwsSecret;
}

export namespace AwsSecret {
  export type AsObject = {
    accessKey: string,
    secretKey: string,
    sessionToken: string,
  }
}

export class AzureSecret extends jspb.Message {
  getApiKeysMap(): jspb.Map<string, string>;
  clearApiKeysMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AzureSecret.AsObject;
  static toObject(includeInstance: boolean, msg: AzureSecret): AzureSecret.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: AzureSecret, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AzureSecret;
  static deserializeBinaryFromReader(message: AzureSecret, reader: jspb.BinaryReader): AzureSecret;
}

export namespace AzureSecret {
  export type AsObject = {
    apiKeysMap: Array<[string, string]>,
  }
}

export class TlsSecret extends jspb.Message {
  getCertChain(): string;
  setCertChain(value: string): void;

  getPrivateKey(): string;
  setPrivateKey(value: string): void;

  getRootCa(): string;
  setRootCa(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TlsSecret.AsObject;
  static toObject(includeInstance: boolean, msg: TlsSecret): TlsSecret.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: TlsSecret, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TlsSecret;
  static deserializeBinaryFromReader(message: TlsSecret, reader: jspb.BinaryReader): TlsSecret;
}

export namespace TlsSecret {
  export type AsObject = {
    certChain: string,
    privateKey: string,
    rootCa: string,
  }
}

export class HeaderSecret extends jspb.Message {
  getHeadersMap(): jspb.Map<string, string>;
  clearHeadersMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HeaderSecret.AsObject;
  static toObject(includeInstance: boolean, msg: HeaderSecret): HeaderSecret.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: HeaderSecret, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HeaderSecret;
  static deserializeBinaryFromReader(message: HeaderSecret, reader: jspb.BinaryReader): HeaderSecret;
}

export namespace HeaderSecret {
  export type AsObject = {
    headersMap: Array<[string, string]>,
  }
}

export class AccountCredentialsSecret extends jspb.Message {
  getUsername(): string;
  setUsername(value: string): void;

  getPassword(): string;
  setPassword(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AccountCredentialsSecret.AsObject;
  static toObject(includeInstance: boolean, msg: AccountCredentialsSecret): AccountCredentialsSecret.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: AccountCredentialsSecret, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AccountCredentialsSecret;
  static deserializeBinaryFromReader(message: AccountCredentialsSecret, reader: jspb.BinaryReader): AccountCredentialsSecret;
}

export namespace AccountCredentialsSecret {
  export type AsObject = {
    username: string,
    password: string,
  }
}

export class EncryptionKeySecret extends jspb.Message {
  getKey(): string;
  setKey(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EncryptionKeySecret.AsObject;
  static toObject(includeInstance: boolean, msg: EncryptionKeySecret): EncryptionKeySecret.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: EncryptionKeySecret, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EncryptionKeySecret;
  static deserializeBinaryFromReader(message: EncryptionKeySecret, reader: jspb.BinaryReader): EncryptionKeySecret;
}

export namespace EncryptionKeySecret {
  export type AsObject = {
    key: string,
  }
}
