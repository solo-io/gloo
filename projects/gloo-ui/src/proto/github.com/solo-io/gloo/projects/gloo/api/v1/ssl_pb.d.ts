// package: gloo.solo.io
// file: github.com/solo-io/gloo/projects/gloo/api/v1/ssl.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../../../../../gogoproto/gogo_pb";
import * as github_com_solo_io_solo_kit_api_v1_ref_pb from "../../../../../../../github.com/solo-io/solo-kit/api/v1/ref_pb";

export class SslConfig extends jspb.Message {
  hasSecretRef(): boolean;
  clearSecretRef(): void;
  getSecretRef(): github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef | undefined;
  setSecretRef(value?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef): void;

  hasSslFiles(): boolean;
  clearSslFiles(): void;
  getSslFiles(): SSLFiles | undefined;
  setSslFiles(value?: SSLFiles): void;

  hasSds(): boolean;
  clearSds(): void;
  getSds(): SDSConfig | undefined;
  setSds(value?: SDSConfig): void;

  clearSniDomainsList(): void;
  getSniDomainsList(): Array<string>;
  setSniDomainsList(value: Array<string>): void;
  addSniDomains(value: string, index?: number): string;

  clearVerifySubjectAltNameList(): void;
  getVerifySubjectAltNameList(): Array<string>;
  setVerifySubjectAltNameList(value: Array<string>): void;
  addVerifySubjectAltName(value: string, index?: number): string;

  hasParameters(): boolean;
  clearParameters(): void;
  getParameters(): SslParameters | undefined;
  setParameters(value?: SslParameters): void;

  getSslSecretsCase(): SslConfig.SslSecretsCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SslConfig.AsObject;
  static toObject(includeInstance: boolean, msg: SslConfig): SslConfig.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: SslConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SslConfig;
  static deserializeBinaryFromReader(message: SslConfig, reader: jspb.BinaryReader): SslConfig;
}

export namespace SslConfig {
  export type AsObject = {
    secretRef?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject,
    sslFiles?: SSLFiles.AsObject,
    sds?: SDSConfig.AsObject,
    sniDomainsList: Array<string>,
    verifySubjectAltNameList: Array<string>,
    parameters?: SslParameters.AsObject,
  }

  export enum SslSecretsCase {
    SSL_SECRETS_NOT_SET = 0,
    SECRET_REF = 1,
    SSL_FILES = 2,
    SDS = 4,
  }
}

export class SSLFiles extends jspb.Message {
  getTlsCert(): string;
  setTlsCert(value: string): void;

  getTlsKey(): string;
  setTlsKey(value: string): void;

  getRootCa(): string;
  setRootCa(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SSLFiles.AsObject;
  static toObject(includeInstance: boolean, msg: SSLFiles): SSLFiles.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: SSLFiles, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SSLFiles;
  static deserializeBinaryFromReader(message: SSLFiles, reader: jspb.BinaryReader): SSLFiles;
}

export namespace SSLFiles {
  export type AsObject = {
    tlsCert: string,
    tlsKey: string,
    rootCa: string,
  }
}

export class UpstreamSslConfig extends jspb.Message {
  hasSecretRef(): boolean;
  clearSecretRef(): void;
  getSecretRef(): github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef | undefined;
  setSecretRef(value?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef): void;

  hasSslFiles(): boolean;
  clearSslFiles(): void;
  getSslFiles(): SSLFiles | undefined;
  setSslFiles(value?: SSLFiles): void;

  hasSds(): boolean;
  clearSds(): void;
  getSds(): SDSConfig | undefined;
  setSds(value?: SDSConfig): void;

  getSni(): string;
  setSni(value: string): void;

  clearVerifySubjectAltNameList(): void;
  getVerifySubjectAltNameList(): Array<string>;
  setVerifySubjectAltNameList(value: Array<string>): void;
  addVerifySubjectAltName(value: string, index?: number): string;

  hasParameters(): boolean;
  clearParameters(): void;
  getParameters(): SslParameters | undefined;
  setParameters(value?: SslParameters): void;

  getSslSecretsCase(): UpstreamSslConfig.SslSecretsCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpstreamSslConfig.AsObject;
  static toObject(includeInstance: boolean, msg: UpstreamSslConfig): UpstreamSslConfig.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpstreamSslConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpstreamSslConfig;
  static deserializeBinaryFromReader(message: UpstreamSslConfig, reader: jspb.BinaryReader): UpstreamSslConfig;
}

export namespace UpstreamSslConfig {
  export type AsObject = {
    secretRef?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject,
    sslFiles?: SSLFiles.AsObject,
    sds?: SDSConfig.AsObject,
    sni: string,
    verifySubjectAltNameList: Array<string>,
    parameters?: SslParameters.AsObject,
  }

  export enum SslSecretsCase {
    SSL_SECRETS_NOT_SET = 0,
    SECRET_REF = 1,
    SSL_FILES = 2,
    SDS = 4,
  }
}

export class SDSConfig extends jspb.Message {
  getTargetUri(): string;
  setTargetUri(value: string): void;

  hasCallCredentials(): boolean;
  clearCallCredentials(): void;
  getCallCredentials(): CallCredentials | undefined;
  setCallCredentials(value?: CallCredentials): void;

  getCertificatesSecretName(): string;
  setCertificatesSecretName(value: string): void;

  getValidationContextName(): string;
  setValidationContextName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SDSConfig.AsObject;
  static toObject(includeInstance: boolean, msg: SDSConfig): SDSConfig.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: SDSConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SDSConfig;
  static deserializeBinaryFromReader(message: SDSConfig, reader: jspb.BinaryReader): SDSConfig;
}

export namespace SDSConfig {
  export type AsObject = {
    targetUri: string,
    callCredentials?: CallCredentials.AsObject,
    certificatesSecretName: string,
    validationContextName: string,
  }
}

export class CallCredentials extends jspb.Message {
  hasFileCredentialSource(): boolean;
  clearFileCredentialSource(): void;
  getFileCredentialSource(): CallCredentials.FileCredentialSource | undefined;
  setFileCredentialSource(value?: CallCredentials.FileCredentialSource): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CallCredentials.AsObject;
  static toObject(includeInstance: boolean, msg: CallCredentials): CallCredentials.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CallCredentials, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CallCredentials;
  static deserializeBinaryFromReader(message: CallCredentials, reader: jspb.BinaryReader): CallCredentials;
}

export namespace CallCredentials {
  export type AsObject = {
    fileCredentialSource?: CallCredentials.FileCredentialSource.AsObject,
  }

  export class FileCredentialSource extends jspb.Message {
    getTokenFileName(): string;
    setTokenFileName(value: string): void;

    getHeader(): string;
    setHeader(value: string): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): FileCredentialSource.AsObject;
    static toObject(includeInstance: boolean, msg: FileCredentialSource): FileCredentialSource.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: FileCredentialSource, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): FileCredentialSource;
    static deserializeBinaryFromReader(message: FileCredentialSource, reader: jspb.BinaryReader): FileCredentialSource;
  }

  export namespace FileCredentialSource {
    export type AsObject = {
      tokenFileName: string,
      header: string,
    }
  }
}

export class SslParameters extends jspb.Message {
  getMinimumProtocolVersion(): SslParameters.ProtocolVersion;
  setMinimumProtocolVersion(value: SslParameters.ProtocolVersion): void;

  getMaximumProtocolVersion(): SslParameters.ProtocolVersion;
  setMaximumProtocolVersion(value: SslParameters.ProtocolVersion): void;

  clearCipherSuitesList(): void;
  getCipherSuitesList(): Array<string>;
  setCipherSuitesList(value: Array<string>): void;
  addCipherSuites(value: string, index?: number): string;

  clearEcdhCurvesList(): void;
  getEcdhCurvesList(): Array<string>;
  setEcdhCurvesList(value: Array<string>): void;
  addEcdhCurves(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SslParameters.AsObject;
  static toObject(includeInstance: boolean, msg: SslParameters): SslParameters.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: SslParameters, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SslParameters;
  static deserializeBinaryFromReader(message: SslParameters, reader: jspb.BinaryReader): SslParameters;
}

export namespace SslParameters {
  export type AsObject = {
    minimumProtocolVersion: SslParameters.ProtocolVersion,
    maximumProtocolVersion: SslParameters.ProtocolVersion,
    cipherSuitesList: Array<string>,
    ecdhCurvesList: Array<string>,
  }

  export enum ProtocolVersion {
    TLS_AUTO = 0,
    TLSV1_0 = 1,
    TLSV1_1 = 2,
    TLSV1_2 = 3,
    TLSV1_3 = 4,
  }
}

