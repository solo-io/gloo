/* eslint-disable */
// package: envoy.config.filter.http.aws_lambda.v2
// file: github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/extensions/aws/filter.proto

import * as jspb from "google-protobuf";
import * as extproto_ext_pb from "../../../../../../../../../../extproto/ext_pb";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";
import * as google_protobuf_duration_pb from "google-protobuf/google/protobuf/duration_pb";
import * as validate_validate_pb from "../../../../../../../../../../validate/validate_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_extension_pb from "../../../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/config/core/v3/extension_pb";

export class AWSLambdaPerRoute extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getQualifier(): string;
  setQualifier(value: string): void;

  getAsync(): boolean;
  setAsync(value: boolean): void;

  hasEmptyBodyOverride(): boolean;
  clearEmptyBodyOverride(): void;
  getEmptyBodyOverride(): google_protobuf_wrappers_pb.StringValue | undefined;
  setEmptyBodyOverride(value?: google_protobuf_wrappers_pb.StringValue): void;

  getUnwrapAsAlb(): boolean;
  setUnwrapAsAlb(value: boolean): void;

  hasTransformerConfig(): boolean;
  clearTransformerConfig(): void;
  getTransformerConfig(): github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_extension_pb.TypedExtensionConfig | undefined;
  setTransformerConfig(value?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_extension_pb.TypedExtensionConfig): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AWSLambdaPerRoute.AsObject;
  static toObject(includeInstance: boolean, msg: AWSLambdaPerRoute): AWSLambdaPerRoute.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: AWSLambdaPerRoute, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AWSLambdaPerRoute;
  static deserializeBinaryFromReader(message: AWSLambdaPerRoute, reader: jspb.BinaryReader): AWSLambdaPerRoute;
}

export namespace AWSLambdaPerRoute {
  export type AsObject = {
    name: string,
    qualifier: string,
    async: boolean,
    emptyBodyOverride?: google_protobuf_wrappers_pb.StringValue.AsObject,
    unwrapAsAlb: boolean,
    transformerConfig?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_config_core_v3_extension_pb.TypedExtensionConfig.AsObject,
  }
}

export class AWSLambdaProtocolExtension extends jspb.Message {
  getHost(): string;
  setHost(value: string): void;

  getRegion(): string;
  setRegion(value: string): void;

  getAccessKey(): string;
  setAccessKey(value: string): void;

  getSecretKey(): string;
  setSecretKey(value: string): void;

  getSessionToken(): string;
  setSessionToken(value: string): void;

  getRoleArn(): string;
  setRoleArn(value: string): void;

  getDisableRoleChaining(): boolean;
  setDisableRoleChaining(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AWSLambdaProtocolExtension.AsObject;
  static toObject(includeInstance: boolean, msg: AWSLambdaProtocolExtension): AWSLambdaProtocolExtension.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: AWSLambdaProtocolExtension, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AWSLambdaProtocolExtension;
  static deserializeBinaryFromReader(message: AWSLambdaProtocolExtension, reader: jspb.BinaryReader): AWSLambdaProtocolExtension;
}

export namespace AWSLambdaProtocolExtension {
  export type AsObject = {
    host: string,
    region: string,
    accessKey: string,
    secretKey: string,
    sessionToken: string,
    roleArn: string,
    disableRoleChaining: boolean,
  }
}

export class AWSLambdaConfig extends jspb.Message {
  hasUseDefaultCredentials(): boolean;
  clearUseDefaultCredentials(): void;
  getUseDefaultCredentials(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setUseDefaultCredentials(value?: google_protobuf_wrappers_pb.BoolValue): void;

  hasServiceAccountCredentials(): boolean;
  clearServiceAccountCredentials(): void;
  getServiceAccountCredentials(): AWSLambdaConfig.ServiceAccountCredentials | undefined;
  setServiceAccountCredentials(value?: AWSLambdaConfig.ServiceAccountCredentials): void;

  getPropagateOriginalRouting(): boolean;
  setPropagateOriginalRouting(value: boolean): void;

  hasCredentialRefreshDelay(): boolean;
  clearCredentialRefreshDelay(): void;
  getCredentialRefreshDelay(): google_protobuf_duration_pb.Duration | undefined;
  setCredentialRefreshDelay(value?: google_protobuf_duration_pb.Duration): void;

  getCredentialsFetcherCase(): AWSLambdaConfig.CredentialsFetcherCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AWSLambdaConfig.AsObject;
  static toObject(includeInstance: boolean, msg: AWSLambdaConfig): AWSLambdaConfig.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: AWSLambdaConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AWSLambdaConfig;
  static deserializeBinaryFromReader(message: AWSLambdaConfig, reader: jspb.BinaryReader): AWSLambdaConfig;
}

export namespace AWSLambdaConfig {
  export type AsObject = {
    useDefaultCredentials?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    serviceAccountCredentials?: AWSLambdaConfig.ServiceAccountCredentials.AsObject,
    propagateOriginalRouting: boolean,
    credentialRefreshDelay?: google_protobuf_duration_pb.Duration.AsObject,
  }

  export class ServiceAccountCredentials extends jspb.Message {
    getCluster(): string;
    setCluster(value: string): void;

    getUri(): string;
    setUri(value: string): void;

    hasTimeout(): boolean;
    clearTimeout(): void;
    getTimeout(): google_protobuf_duration_pb.Duration | undefined;
    setTimeout(value?: google_protobuf_duration_pb.Duration): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ServiceAccountCredentials.AsObject;
    static toObject(includeInstance: boolean, msg: ServiceAccountCredentials): ServiceAccountCredentials.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: ServiceAccountCredentials, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ServiceAccountCredentials;
    static deserializeBinaryFromReader(message: ServiceAccountCredentials, reader: jspb.BinaryReader): ServiceAccountCredentials;
  }

  export namespace ServiceAccountCredentials {
    export type AsObject = {
      cluster: string,
      uri: string,
      timeout?: google_protobuf_duration_pb.Duration.AsObject,
    }
  }

  export enum CredentialsFetcherCase {
    CREDENTIALS_FETCHER_NOT_SET = 0,
    USE_DEFAULT_CREDENTIALS = 1,
    SERVICE_ACCOUNT_CREDENTIALS = 2,
  }
}

export class ApiGatewayTransformation extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ApiGatewayTransformation.AsObject;
  static toObject(includeInstance: boolean, msg: ApiGatewayTransformation): ApiGatewayTransformation.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ApiGatewayTransformation, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ApiGatewayTransformation;
  static deserializeBinaryFromReader(message: ApiGatewayTransformation, reader: jspb.BinaryReader): ApiGatewayTransformation;
}

export namespace ApiGatewayTransformation {
  export type AsObject = {
  }
}
