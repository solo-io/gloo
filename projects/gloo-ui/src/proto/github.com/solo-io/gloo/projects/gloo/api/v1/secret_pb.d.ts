// package: gloo.solo.io
// file: github.com/solo-io/gloo/projects/gloo/api/v1/secret.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../../../../../gogoproto/gogo_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_v1_extensions_pb from "../../../../../../../github.com/solo-io/gloo/projects/gloo/api/v1/extensions_pb";
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

  hasExtension(): boolean;
  clearExtension(): void;
  getExtension(): github_com_solo_io_gloo_projects_gloo_api_v1_extensions_pb.Extension | undefined;
  setExtension(value?: github_com_solo_io_gloo_projects_gloo_api_v1_extensions_pb.Extension): void;

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
    extension?: github_com_solo_io_gloo_projects_gloo_api_v1_extensions_pb.Extension.AsObject,
    metadata?: github_com_solo_io_solo_kit_api_v1_metadata_pb.Metadata.AsObject,
  }

  export enum KindCase {
    KIND_NOT_SET = 0,
    AWS = 1,
    AZURE = 2,
    TLS = 3,
    EXTENSION = 4,
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

