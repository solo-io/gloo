/* eslint-disable */
// package: envoy.config.core.v3
// file: envoy/config/core/v3/grpc_service.proto

import * as jspb from "google-protobuf";
import * as envoy_config_core_v3_base_pb from "../../../../envoy/config/core/v3/base_pb";
import * as google_protobuf_any_pb from "google-protobuf/google/protobuf/any_pb";
import * as google_protobuf_duration_pb from "google-protobuf/google/protobuf/duration_pb";
import * as google_protobuf_empty_pb from "google-protobuf/google/protobuf/empty_pb";
import * as google_protobuf_struct_pb from "google-protobuf/google/protobuf/struct_pb";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";
import * as udpa_annotations_sensitive_pb from "../../../../udpa/annotations/sensitive_pb";
import * as udpa_annotations_status_pb from "../../../../udpa/annotations/status_pb";
import * as udpa_annotations_versioning_pb from "../../../../udpa/annotations/versioning_pb";
import * as validate_validate_pb from "../../../../validate/validate_pb";
import * as gogoproto_gogo_pb from "../../../../gogoproto/gogo_pb";

export class GrpcService extends jspb.Message {
  hasEnvoyGrpc(): boolean;
  clearEnvoyGrpc(): void;
  getEnvoyGrpc(): GrpcService.EnvoyGrpc | undefined;
  setEnvoyGrpc(value?: GrpcService.EnvoyGrpc): void;

  hasGoogleGrpc(): boolean;
  clearGoogleGrpc(): void;
  getGoogleGrpc(): GrpcService.GoogleGrpc | undefined;
  setGoogleGrpc(value?: GrpcService.GoogleGrpc): void;

  hasTimeout(): boolean;
  clearTimeout(): void;
  getTimeout(): google_protobuf_duration_pb.Duration | undefined;
  setTimeout(value?: google_protobuf_duration_pb.Duration): void;

  clearInitialMetadataList(): void;
  getInitialMetadataList(): Array<envoy_config_core_v3_base_pb.HeaderValue>;
  setInitialMetadataList(value: Array<envoy_config_core_v3_base_pb.HeaderValue>): void;
  addInitialMetadata(value?: envoy_config_core_v3_base_pb.HeaderValue, index?: number): envoy_config_core_v3_base_pb.HeaderValue;

  getTargetSpecifierCase(): GrpcService.TargetSpecifierCase;
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
    envoyGrpc?: GrpcService.EnvoyGrpc.AsObject,
    googleGrpc?: GrpcService.GoogleGrpc.AsObject,
    timeout?: google_protobuf_duration_pb.Duration.AsObject,
    initialMetadataList: Array<envoy_config_core_v3_base_pb.HeaderValue.AsObject>,
  }

  export class EnvoyGrpc extends jspb.Message {
    getClusterName(): string;
    setClusterName(value: string): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): EnvoyGrpc.AsObject;
    static toObject(includeInstance: boolean, msg: EnvoyGrpc): EnvoyGrpc.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: EnvoyGrpc, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): EnvoyGrpc;
    static deserializeBinaryFromReader(message: EnvoyGrpc, reader: jspb.BinaryReader): EnvoyGrpc;
  }

  export namespace EnvoyGrpc {
    export type AsObject = {
      clusterName: string,
    }
  }

  export class GoogleGrpc extends jspb.Message {
    getTargetUri(): string;
    setTargetUri(value: string): void;

    hasChannelCredentials(): boolean;
    clearChannelCredentials(): void;
    getChannelCredentials(): GrpcService.GoogleGrpc.ChannelCredentials | undefined;
    setChannelCredentials(value?: GrpcService.GoogleGrpc.ChannelCredentials): void;

    clearCallCredentialsList(): void;
    getCallCredentialsList(): Array<GrpcService.GoogleGrpc.CallCredentials>;
    setCallCredentialsList(value: Array<GrpcService.GoogleGrpc.CallCredentials>): void;
    addCallCredentials(value?: GrpcService.GoogleGrpc.CallCredentials, index?: number): GrpcService.GoogleGrpc.CallCredentials;

    getStatPrefix(): string;
    setStatPrefix(value: string): void;

    getCredentialsFactoryName(): string;
    setCredentialsFactoryName(value: string): void;

    hasConfig(): boolean;
    clearConfig(): void;
    getConfig(): google_protobuf_struct_pb.Struct | undefined;
    setConfig(value?: google_protobuf_struct_pb.Struct): void;

    hasPerStreamBufferLimitBytes(): boolean;
    clearPerStreamBufferLimitBytes(): void;
    getPerStreamBufferLimitBytes(): google_protobuf_wrappers_pb.UInt32Value | undefined;
    setPerStreamBufferLimitBytes(value?: google_protobuf_wrappers_pb.UInt32Value): void;

    hasChannelArgs(): boolean;
    clearChannelArgs(): void;
    getChannelArgs(): GrpcService.GoogleGrpc.ChannelArgs | undefined;
    setChannelArgs(value?: GrpcService.GoogleGrpc.ChannelArgs): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): GoogleGrpc.AsObject;
    static toObject(includeInstance: boolean, msg: GoogleGrpc): GoogleGrpc.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: GoogleGrpc, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): GoogleGrpc;
    static deserializeBinaryFromReader(message: GoogleGrpc, reader: jspb.BinaryReader): GoogleGrpc;
  }

  export namespace GoogleGrpc {
    export type AsObject = {
      targetUri: string,
      channelCredentials?: GrpcService.GoogleGrpc.ChannelCredentials.AsObject,
      callCredentialsList: Array<GrpcService.GoogleGrpc.CallCredentials.AsObject>,
      statPrefix: string,
      credentialsFactoryName: string,
      config?: google_protobuf_struct_pb.Struct.AsObject,
      perStreamBufferLimitBytes?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
      channelArgs?: GrpcService.GoogleGrpc.ChannelArgs.AsObject,
    }

    export class SslCredentials extends jspb.Message {
      hasRootCerts(): boolean;
      clearRootCerts(): void;
      getRootCerts(): envoy_config_core_v3_base_pb.DataSource | undefined;
      setRootCerts(value?: envoy_config_core_v3_base_pb.DataSource): void;

      hasPrivateKey(): boolean;
      clearPrivateKey(): void;
      getPrivateKey(): envoy_config_core_v3_base_pb.DataSource | undefined;
      setPrivateKey(value?: envoy_config_core_v3_base_pb.DataSource): void;

      hasCertChain(): boolean;
      clearCertChain(): void;
      getCertChain(): envoy_config_core_v3_base_pb.DataSource | undefined;
      setCertChain(value?: envoy_config_core_v3_base_pb.DataSource): void;

      serializeBinary(): Uint8Array;
      toObject(includeInstance?: boolean): SslCredentials.AsObject;
      static toObject(includeInstance: boolean, msg: SslCredentials): SslCredentials.AsObject;
      static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
      static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
      static serializeBinaryToWriter(message: SslCredentials, writer: jspb.BinaryWriter): void;
      static deserializeBinary(bytes: Uint8Array): SslCredentials;
      static deserializeBinaryFromReader(message: SslCredentials, reader: jspb.BinaryReader): SslCredentials;
    }

    export namespace SslCredentials {
      export type AsObject = {
        rootCerts?: envoy_config_core_v3_base_pb.DataSource.AsObject,
        privateKey?: envoy_config_core_v3_base_pb.DataSource.AsObject,
        certChain?: envoy_config_core_v3_base_pb.DataSource.AsObject,
      }
    }

    export class GoogleLocalCredentials extends jspb.Message {
      serializeBinary(): Uint8Array;
      toObject(includeInstance?: boolean): GoogleLocalCredentials.AsObject;
      static toObject(includeInstance: boolean, msg: GoogleLocalCredentials): GoogleLocalCredentials.AsObject;
      static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
      static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
      static serializeBinaryToWriter(message: GoogleLocalCredentials, writer: jspb.BinaryWriter): void;
      static deserializeBinary(bytes: Uint8Array): GoogleLocalCredentials;
      static deserializeBinaryFromReader(message: GoogleLocalCredentials, reader: jspb.BinaryReader): GoogleLocalCredentials;
    }

    export namespace GoogleLocalCredentials {
      export type AsObject = {
      }
    }

    export class ChannelCredentials extends jspb.Message {
      hasSslCredentials(): boolean;
      clearSslCredentials(): void;
      getSslCredentials(): GrpcService.GoogleGrpc.SslCredentials | undefined;
      setSslCredentials(value?: GrpcService.GoogleGrpc.SslCredentials): void;

      hasGoogleDefault(): boolean;
      clearGoogleDefault(): void;
      getGoogleDefault(): google_protobuf_empty_pb.Empty | undefined;
      setGoogleDefault(value?: google_protobuf_empty_pb.Empty): void;

      hasLocalCredentials(): boolean;
      clearLocalCredentials(): void;
      getLocalCredentials(): GrpcService.GoogleGrpc.GoogleLocalCredentials | undefined;
      setLocalCredentials(value?: GrpcService.GoogleGrpc.GoogleLocalCredentials): void;

      getCredentialSpecifierCase(): ChannelCredentials.CredentialSpecifierCase;
      serializeBinary(): Uint8Array;
      toObject(includeInstance?: boolean): ChannelCredentials.AsObject;
      static toObject(includeInstance: boolean, msg: ChannelCredentials): ChannelCredentials.AsObject;
      static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
      static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
      static serializeBinaryToWriter(message: ChannelCredentials, writer: jspb.BinaryWriter): void;
      static deserializeBinary(bytes: Uint8Array): ChannelCredentials;
      static deserializeBinaryFromReader(message: ChannelCredentials, reader: jspb.BinaryReader): ChannelCredentials;
    }

    export namespace ChannelCredentials {
      export type AsObject = {
        sslCredentials?: GrpcService.GoogleGrpc.SslCredentials.AsObject,
        googleDefault?: google_protobuf_empty_pb.Empty.AsObject,
        localCredentials?: GrpcService.GoogleGrpc.GoogleLocalCredentials.AsObject,
      }

      export enum CredentialSpecifierCase {
        CREDENTIAL_SPECIFIER_NOT_SET = 0,
        SSL_CREDENTIALS = 1,
        GOOGLE_DEFAULT = 2,
        LOCAL_CREDENTIALS = 3,
      }
    }

    export class CallCredentials extends jspb.Message {
      hasAccessToken(): boolean;
      clearAccessToken(): void;
      getAccessToken(): string;
      setAccessToken(value: string): void;

      hasGoogleComputeEngine(): boolean;
      clearGoogleComputeEngine(): void;
      getGoogleComputeEngine(): google_protobuf_empty_pb.Empty | undefined;
      setGoogleComputeEngine(value?: google_protobuf_empty_pb.Empty): void;

      hasGoogleRefreshToken(): boolean;
      clearGoogleRefreshToken(): void;
      getGoogleRefreshToken(): string;
      setGoogleRefreshToken(value: string): void;

      hasServiceAccountJwtAccess(): boolean;
      clearServiceAccountJwtAccess(): void;
      getServiceAccountJwtAccess(): GrpcService.GoogleGrpc.CallCredentials.ServiceAccountJWTAccessCredentials | undefined;
      setServiceAccountJwtAccess(value?: GrpcService.GoogleGrpc.CallCredentials.ServiceAccountJWTAccessCredentials): void;

      hasGoogleIam(): boolean;
      clearGoogleIam(): void;
      getGoogleIam(): GrpcService.GoogleGrpc.CallCredentials.GoogleIAMCredentials | undefined;
      setGoogleIam(value?: GrpcService.GoogleGrpc.CallCredentials.GoogleIAMCredentials): void;

      hasFromPlugin(): boolean;
      clearFromPlugin(): void;
      getFromPlugin(): GrpcService.GoogleGrpc.CallCredentials.MetadataCredentialsFromPlugin | undefined;
      setFromPlugin(value?: GrpcService.GoogleGrpc.CallCredentials.MetadataCredentialsFromPlugin): void;

      hasStsService(): boolean;
      clearStsService(): void;
      getStsService(): GrpcService.GoogleGrpc.CallCredentials.StsService | undefined;
      setStsService(value?: GrpcService.GoogleGrpc.CallCredentials.StsService): void;

      getCredentialSpecifierCase(): CallCredentials.CredentialSpecifierCase;
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
        accessToken: string,
        googleComputeEngine?: google_protobuf_empty_pb.Empty.AsObject,
        googleRefreshToken: string,
        serviceAccountJwtAccess?: GrpcService.GoogleGrpc.CallCredentials.ServiceAccountJWTAccessCredentials.AsObject,
        googleIam?: GrpcService.GoogleGrpc.CallCredentials.GoogleIAMCredentials.AsObject,
        fromPlugin?: GrpcService.GoogleGrpc.CallCredentials.MetadataCredentialsFromPlugin.AsObject,
        stsService?: GrpcService.GoogleGrpc.CallCredentials.StsService.AsObject,
      }

      export class ServiceAccountJWTAccessCredentials extends jspb.Message {
        getJsonKey(): string;
        setJsonKey(value: string): void;

        getTokenLifetimeSeconds(): number;
        setTokenLifetimeSeconds(value: number): void;

        serializeBinary(): Uint8Array;
        toObject(includeInstance?: boolean): ServiceAccountJWTAccessCredentials.AsObject;
        static toObject(includeInstance: boolean, msg: ServiceAccountJWTAccessCredentials): ServiceAccountJWTAccessCredentials.AsObject;
        static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
        static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
        static serializeBinaryToWriter(message: ServiceAccountJWTAccessCredentials, writer: jspb.BinaryWriter): void;
        static deserializeBinary(bytes: Uint8Array): ServiceAccountJWTAccessCredentials;
        static deserializeBinaryFromReader(message: ServiceAccountJWTAccessCredentials, reader: jspb.BinaryReader): ServiceAccountJWTAccessCredentials;
      }

      export namespace ServiceAccountJWTAccessCredentials {
        export type AsObject = {
          jsonKey: string,
          tokenLifetimeSeconds: number,
        }
      }

      export class GoogleIAMCredentials extends jspb.Message {
        getAuthorizationToken(): string;
        setAuthorizationToken(value: string): void;

        getAuthoritySelector(): string;
        setAuthoritySelector(value: string): void;

        serializeBinary(): Uint8Array;
        toObject(includeInstance?: boolean): GoogleIAMCredentials.AsObject;
        static toObject(includeInstance: boolean, msg: GoogleIAMCredentials): GoogleIAMCredentials.AsObject;
        static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
        static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
        static serializeBinaryToWriter(message: GoogleIAMCredentials, writer: jspb.BinaryWriter): void;
        static deserializeBinary(bytes: Uint8Array): GoogleIAMCredentials;
        static deserializeBinaryFromReader(message: GoogleIAMCredentials, reader: jspb.BinaryReader): GoogleIAMCredentials;
      }

      export namespace GoogleIAMCredentials {
        export type AsObject = {
          authorizationToken: string,
          authoritySelector: string,
        }
      }

      export class MetadataCredentialsFromPlugin extends jspb.Message {
        getName(): string;
        setName(value: string): void;

        hasTypedConfig(): boolean;
        clearTypedConfig(): void;
        getTypedConfig(): google_protobuf_any_pb.Any | undefined;
        setTypedConfig(value?: google_protobuf_any_pb.Any): void;

        getConfigTypeCase(): MetadataCredentialsFromPlugin.ConfigTypeCase;
        serializeBinary(): Uint8Array;
        toObject(includeInstance?: boolean): MetadataCredentialsFromPlugin.AsObject;
        static toObject(includeInstance: boolean, msg: MetadataCredentialsFromPlugin): MetadataCredentialsFromPlugin.AsObject;
        static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
        static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
        static serializeBinaryToWriter(message: MetadataCredentialsFromPlugin, writer: jspb.BinaryWriter): void;
        static deserializeBinary(bytes: Uint8Array): MetadataCredentialsFromPlugin;
        static deserializeBinaryFromReader(message: MetadataCredentialsFromPlugin, reader: jspb.BinaryReader): MetadataCredentialsFromPlugin;
      }

      export namespace MetadataCredentialsFromPlugin {
        export type AsObject = {
          name: string,
          typedConfig?: google_protobuf_any_pb.Any.AsObject,
        }

        export enum ConfigTypeCase {
          CONFIG_TYPE_NOT_SET = 0,
          TYPED_CONFIG = 3,
        }
      }

      export class StsService extends jspb.Message {
        getTokenExchangeServiceUri(): string;
        setTokenExchangeServiceUri(value: string): void;

        getResource(): string;
        setResource(value: string): void;

        getAudience(): string;
        setAudience(value: string): void;

        getScope(): string;
        setScope(value: string): void;

        getRequestedTokenType(): string;
        setRequestedTokenType(value: string): void;

        getSubjectTokenPath(): string;
        setSubjectTokenPath(value: string): void;

        getSubjectTokenType(): string;
        setSubjectTokenType(value: string): void;

        getActorTokenPath(): string;
        setActorTokenPath(value: string): void;

        getActorTokenType(): string;
        setActorTokenType(value: string): void;

        serializeBinary(): Uint8Array;
        toObject(includeInstance?: boolean): StsService.AsObject;
        static toObject(includeInstance: boolean, msg: StsService): StsService.AsObject;
        static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
        static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
        static serializeBinaryToWriter(message: StsService, writer: jspb.BinaryWriter): void;
        static deserializeBinary(bytes: Uint8Array): StsService;
        static deserializeBinaryFromReader(message: StsService, reader: jspb.BinaryReader): StsService;
      }

      export namespace StsService {
        export type AsObject = {
          tokenExchangeServiceUri: string,
          resource: string,
          audience: string,
          scope: string,
          requestedTokenType: string,
          subjectTokenPath: string,
          subjectTokenType: string,
          actorTokenPath: string,
          actorTokenType: string,
        }
      }

      export enum CredentialSpecifierCase {
        CREDENTIAL_SPECIFIER_NOT_SET = 0,
        ACCESS_TOKEN = 1,
        GOOGLE_COMPUTE_ENGINE = 2,
        GOOGLE_REFRESH_TOKEN = 3,
        SERVICE_ACCOUNT_JWT_ACCESS = 4,
        GOOGLE_IAM = 5,
        FROM_PLUGIN = 6,
        STS_SERVICE = 7,
      }
    }

    export class ChannelArgs extends jspb.Message {
      getArgsMap(): jspb.Map<string, GrpcService.GoogleGrpc.ChannelArgs.Value>;
      clearArgsMap(): void;
      serializeBinary(): Uint8Array;
      toObject(includeInstance?: boolean): ChannelArgs.AsObject;
      static toObject(includeInstance: boolean, msg: ChannelArgs): ChannelArgs.AsObject;
      static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
      static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
      static serializeBinaryToWriter(message: ChannelArgs, writer: jspb.BinaryWriter): void;
      static deserializeBinary(bytes: Uint8Array): ChannelArgs;
      static deserializeBinaryFromReader(message: ChannelArgs, reader: jspb.BinaryReader): ChannelArgs;
    }

    export namespace ChannelArgs {
      export type AsObject = {
        argsMap: Array<[string, GrpcService.GoogleGrpc.ChannelArgs.Value.AsObject]>,
      }

      export class Value extends jspb.Message {
        hasStringValue(): boolean;
        clearStringValue(): void;
        getStringValue(): string;
        setStringValue(value: string): void;

        hasIntValue(): boolean;
        clearIntValue(): void;
        getIntValue(): number;
        setIntValue(value: number): void;

        getValueSpecifierCase(): Value.ValueSpecifierCase;
        serializeBinary(): Uint8Array;
        toObject(includeInstance?: boolean): Value.AsObject;
        static toObject(includeInstance: boolean, msg: Value): Value.AsObject;
        static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
        static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
        static serializeBinaryToWriter(message: Value, writer: jspb.BinaryWriter): void;
        static deserializeBinary(bytes: Uint8Array): Value;
        static deserializeBinaryFromReader(message: Value, reader: jspb.BinaryReader): Value;
      }

      export namespace Value {
        export type AsObject = {
          stringValue: string,
          intValue: number,
        }

        export enum ValueSpecifierCase {
          VALUE_SPECIFIER_NOT_SET = 0,
          STRING_VALUE = 1,
          INT_VALUE = 2,
        }
      }
    }
  }

  export enum TargetSpecifierCase {
    TARGET_SPECIFIER_NOT_SET = 0,
    ENVOY_GRPC = 1,
    GOOGLE_GRPC = 2,
  }
}
