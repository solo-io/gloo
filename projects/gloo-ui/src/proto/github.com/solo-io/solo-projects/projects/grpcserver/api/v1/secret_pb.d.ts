// package: glooeeapi.solo.io
// file: github.com/solo-io/solo-projects/projects/grpcserver/api/v1/secret.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../../../../../gogoproto/gogo_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_v1_extensions_pb from "../../../../../../../github.com/solo-io/gloo/projects/gloo/api/v1/extensions_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_v1_secret_pb from "../../../../../../../github.com/solo-io/gloo/projects/gloo/api/v1/secret_pb";
import * as github_com_solo_io_solo_kit_api_v1_ref_pb from "../../../../../../../github.com/solo-io/solo-kit/api/v1/ref_pb";

export class GetSecretRequest extends jspb.Message {
  hasRef(): boolean;
  clearRef(): void;
  getRef(): github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef | undefined;
  setRef(value?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetSecretRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetSecretRequest): GetSecretRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetSecretRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetSecretRequest;
  static deserializeBinaryFromReader(message: GetSecretRequest, reader: jspb.BinaryReader): GetSecretRequest;
}

export namespace GetSecretRequest {
  export type AsObject = {
    ref?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject,
  }
}

export class GetSecretResponse extends jspb.Message {
  hasSecret(): boolean;
  clearSecret(): void;
  getSecret(): github_com_solo_io_gloo_projects_gloo_api_v1_secret_pb.Secret | undefined;
  setSecret(value?: github_com_solo_io_gloo_projects_gloo_api_v1_secret_pb.Secret): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetSecretResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetSecretResponse): GetSecretResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetSecretResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetSecretResponse;
  static deserializeBinaryFromReader(message: GetSecretResponse, reader: jspb.BinaryReader): GetSecretResponse;
}

export namespace GetSecretResponse {
  export type AsObject = {
    secret?: github_com_solo_io_gloo_projects_gloo_api_v1_secret_pb.Secret.AsObject,
  }
}

export class ListSecretsRequest extends jspb.Message {
  clearNamespacesList(): void;
  getNamespacesList(): Array<string>;
  setNamespacesList(value: Array<string>): void;
  addNamespaces(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListSecretsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListSecretsRequest): ListSecretsRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListSecretsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListSecretsRequest;
  static deserializeBinaryFromReader(message: ListSecretsRequest, reader: jspb.BinaryReader): ListSecretsRequest;
}

export namespace ListSecretsRequest {
  export type AsObject = {
    namespacesList: Array<string>,
  }
}

export class ListSecretsResponse extends jspb.Message {
  clearSecretsList(): void;
  getSecretsList(): Array<github_com_solo_io_gloo_projects_gloo_api_v1_secret_pb.Secret>;
  setSecretsList(value: Array<github_com_solo_io_gloo_projects_gloo_api_v1_secret_pb.Secret>): void;
  addSecrets(value?: github_com_solo_io_gloo_projects_gloo_api_v1_secret_pb.Secret, index?: number): github_com_solo_io_gloo_projects_gloo_api_v1_secret_pb.Secret;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListSecretsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListSecretsResponse): ListSecretsResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListSecretsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListSecretsResponse;
  static deserializeBinaryFromReader(message: ListSecretsResponse, reader: jspb.BinaryReader): ListSecretsResponse;
}

export namespace ListSecretsResponse {
  export type AsObject = {
    secretsList: Array<github_com_solo_io_gloo_projects_gloo_api_v1_secret_pb.Secret.AsObject>,
  }
}

export class CreateSecretRequest extends jspb.Message {
  hasRef(): boolean;
  clearRef(): void;
  getRef(): github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef | undefined;
  setRef(value?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef): void;

  hasAws(): boolean;
  clearAws(): void;
  getAws(): github_com_solo_io_gloo_projects_gloo_api_v1_secret_pb.AwsSecret | undefined;
  setAws(value?: github_com_solo_io_gloo_projects_gloo_api_v1_secret_pb.AwsSecret): void;

  hasAzure(): boolean;
  clearAzure(): void;
  getAzure(): github_com_solo_io_gloo_projects_gloo_api_v1_secret_pb.AzureSecret | undefined;
  setAzure(value?: github_com_solo_io_gloo_projects_gloo_api_v1_secret_pb.AzureSecret): void;

  hasTls(): boolean;
  clearTls(): void;
  getTls(): github_com_solo_io_gloo_projects_gloo_api_v1_secret_pb.TlsSecret | undefined;
  setTls(value?: github_com_solo_io_gloo_projects_gloo_api_v1_secret_pb.TlsSecret): void;

  hasExtension(): boolean;
  clearExtension(): void;
  getExtension(): github_com_solo_io_gloo_projects_gloo_api_v1_extensions_pb.Extension | undefined;
  setExtension(value?: github_com_solo_io_gloo_projects_gloo_api_v1_extensions_pb.Extension): void;

  hasSecret(): boolean;
  clearSecret(): void;
  getSecret(): github_com_solo_io_gloo_projects_gloo_api_v1_secret_pb.Secret | undefined;
  setSecret(value?: github_com_solo_io_gloo_projects_gloo_api_v1_secret_pb.Secret): void;

  getKindCase(): CreateSecretRequest.KindCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateSecretRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateSecretRequest): CreateSecretRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateSecretRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateSecretRequest;
  static deserializeBinaryFromReader(message: CreateSecretRequest, reader: jspb.BinaryReader): CreateSecretRequest;
}

export namespace CreateSecretRequest {
  export type AsObject = {
    ref?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject,
    aws?: github_com_solo_io_gloo_projects_gloo_api_v1_secret_pb.AwsSecret.AsObject,
    azure?: github_com_solo_io_gloo_projects_gloo_api_v1_secret_pb.AzureSecret.AsObject,
    tls?: github_com_solo_io_gloo_projects_gloo_api_v1_secret_pb.TlsSecret.AsObject,
    extension?: github_com_solo_io_gloo_projects_gloo_api_v1_extensions_pb.Extension.AsObject,
    secret?: github_com_solo_io_gloo_projects_gloo_api_v1_secret_pb.Secret.AsObject,
  }

  export enum KindCase {
    KIND_NOT_SET = 0,
    AWS = 2,
    AZURE = 3,
    TLS = 4,
    EXTENSION = 5,
  }
}

export class CreateSecretResponse extends jspb.Message {
  hasSecret(): boolean;
  clearSecret(): void;
  getSecret(): github_com_solo_io_gloo_projects_gloo_api_v1_secret_pb.Secret | undefined;
  setSecret(value?: github_com_solo_io_gloo_projects_gloo_api_v1_secret_pb.Secret): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateSecretResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CreateSecretResponse): CreateSecretResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateSecretResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateSecretResponse;
  static deserializeBinaryFromReader(message: CreateSecretResponse, reader: jspb.BinaryReader): CreateSecretResponse;
}

export namespace CreateSecretResponse {
  export type AsObject = {
    secret?: github_com_solo_io_gloo_projects_gloo_api_v1_secret_pb.Secret.AsObject,
  }
}

export class UpdateSecretRequest extends jspb.Message {
  hasRef(): boolean;
  clearRef(): void;
  getRef(): github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef | undefined;
  setRef(value?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef): void;

  hasAws(): boolean;
  clearAws(): void;
  getAws(): github_com_solo_io_gloo_projects_gloo_api_v1_secret_pb.AwsSecret | undefined;
  setAws(value?: github_com_solo_io_gloo_projects_gloo_api_v1_secret_pb.AwsSecret): void;

  hasAzure(): boolean;
  clearAzure(): void;
  getAzure(): github_com_solo_io_gloo_projects_gloo_api_v1_secret_pb.AzureSecret | undefined;
  setAzure(value?: github_com_solo_io_gloo_projects_gloo_api_v1_secret_pb.AzureSecret): void;

  hasTls(): boolean;
  clearTls(): void;
  getTls(): github_com_solo_io_gloo_projects_gloo_api_v1_secret_pb.TlsSecret | undefined;
  setTls(value?: github_com_solo_io_gloo_projects_gloo_api_v1_secret_pb.TlsSecret): void;

  hasExtension(): boolean;
  clearExtension(): void;
  getExtension(): github_com_solo_io_gloo_projects_gloo_api_v1_extensions_pb.Extension | undefined;
  setExtension(value?: github_com_solo_io_gloo_projects_gloo_api_v1_extensions_pb.Extension): void;

  hasSecret(): boolean;
  clearSecret(): void;
  getSecret(): github_com_solo_io_gloo_projects_gloo_api_v1_secret_pb.Secret | undefined;
  setSecret(value?: github_com_solo_io_gloo_projects_gloo_api_v1_secret_pb.Secret): void;

  getKindCase(): UpdateSecretRequest.KindCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateSecretRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateSecretRequest): UpdateSecretRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateSecretRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateSecretRequest;
  static deserializeBinaryFromReader(message: UpdateSecretRequest, reader: jspb.BinaryReader): UpdateSecretRequest;
}

export namespace UpdateSecretRequest {
  export type AsObject = {
    ref?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject,
    aws?: github_com_solo_io_gloo_projects_gloo_api_v1_secret_pb.AwsSecret.AsObject,
    azure?: github_com_solo_io_gloo_projects_gloo_api_v1_secret_pb.AzureSecret.AsObject,
    tls?: github_com_solo_io_gloo_projects_gloo_api_v1_secret_pb.TlsSecret.AsObject,
    extension?: github_com_solo_io_gloo_projects_gloo_api_v1_extensions_pb.Extension.AsObject,
    secret?: github_com_solo_io_gloo_projects_gloo_api_v1_secret_pb.Secret.AsObject,
  }

  export enum KindCase {
    KIND_NOT_SET = 0,
    AWS = 2,
    AZURE = 3,
    TLS = 4,
    EXTENSION = 5,
  }
}

export class UpdateSecretResponse extends jspb.Message {
  hasSecret(): boolean;
  clearSecret(): void;
  getSecret(): github_com_solo_io_gloo_projects_gloo_api_v1_secret_pb.Secret | undefined;
  setSecret(value?: github_com_solo_io_gloo_projects_gloo_api_v1_secret_pb.Secret): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateSecretResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateSecretResponse): UpdateSecretResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateSecretResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateSecretResponse;
  static deserializeBinaryFromReader(message: UpdateSecretResponse, reader: jspb.BinaryReader): UpdateSecretResponse;
}

export namespace UpdateSecretResponse {
  export type AsObject = {
    secret?: github_com_solo_io_gloo_projects_gloo_api_v1_secret_pb.Secret.AsObject,
  }
}

export class DeleteSecretRequest extends jspb.Message {
  hasRef(): boolean;
  clearRef(): void;
  getRef(): github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef | undefined;
  setRef(value?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteSecretRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteSecretRequest): DeleteSecretRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteSecretRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteSecretRequest;
  static deserializeBinaryFromReader(message: DeleteSecretRequest, reader: jspb.BinaryReader): DeleteSecretRequest;
}

export namespace DeleteSecretRequest {
  export type AsObject = {
    ref?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject,
  }
}

export class DeleteSecretResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteSecretResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteSecretResponse): DeleteSecretResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteSecretResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteSecretResponse;
  static deserializeBinaryFromReader(message: DeleteSecretResponse, reader: jspb.BinaryReader): DeleteSecretResponse;
}

export namespace DeleteSecretResponse {
  export type AsObject = {
  }
}

