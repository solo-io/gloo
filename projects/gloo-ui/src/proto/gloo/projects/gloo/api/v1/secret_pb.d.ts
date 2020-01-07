// package: gloo.solo.io
// file: gloo/projects/gloo/api/v1/secret.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../../../gogoproto/gogo_pb";
import * as extproto_ext_pb from "../../../../../extproto/ext_pb";
import * as gloo_projects_gloo_api_v1_extensions_pb from "../../../../../gloo/projects/gloo/api/v1/extensions_pb";
import * as gloo_projects_gloo_api_v1_enterprise_options_extauth_v1_extauth_pb from "../../../../../gloo/projects/gloo/api/v1/enterprise/options/extauth/v1/extauth_pb";
import * as solo_kit_api_v1_metadata_pb from "../../../../../solo-kit/api/v1/metadata_pb";
import * as solo_kit_api_v1_solo_kit_pb from "../../../../../solo-kit/api/v1/solo-kit_pb";

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
  getOauth(): gloo_projects_gloo_api_v1_enterprise_options_extauth_v1_extauth_pb.OauthSecret | undefined;
  setOauth(value?: gloo_projects_gloo_api_v1_enterprise_options_extauth_v1_extauth_pb.OauthSecret): void;

  hasApiKey(): boolean;
  clearApiKey(): void;
  getApiKey(): gloo_projects_gloo_api_v1_enterprise_options_extauth_v1_extauth_pb.ApiKeySecret | undefined;
  setApiKey(value?: gloo_projects_gloo_api_v1_enterprise_options_extauth_v1_extauth_pb.ApiKeySecret): void;

  hasExtensions(): boolean;
  clearExtensions(): void;
  getExtensions(): gloo_projects_gloo_api_v1_extensions_pb.Extensions | undefined;
  setExtensions(value?: gloo_projects_gloo_api_v1_extensions_pb.Extensions): void;

  hasMetadata(): boolean;
  clearMetadata(): void;
  getMetadata(): solo_kit_api_v1_metadata_pb.Metadata | undefined;
  setMetadata(value?: solo_kit_api_v1_metadata_pb.Metadata): void;

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
    oauth?: gloo_projects_gloo_api_v1_enterprise_options_extauth_v1_extauth_pb.OauthSecret.AsObject,
    apiKey?: gloo_projects_gloo_api_v1_enterprise_options_extauth_v1_extauth_pb.ApiKeySecret.AsObject,
    extensions?: gloo_projects_gloo_api_v1_extensions_pb.Extensions.AsObject,
    metadata?: solo_kit_api_v1_metadata_pb.Metadata.AsObject,
  }

  export enum KindCase {
    KIND_NOT_SET = 0,
    AWS = 1,
    AZURE = 2,
    TLS = 3,
    OAUTH = 5,
    API_KEY = 6,
    EXTENSIONS = 4,
  }
}

export class AwsSecret extends jspb.Message {
  getAccessKey(): string;
  setAccessKey(value: string): void;

  getSecretKey(): string;
  setSecretKey(value: string): void;

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

