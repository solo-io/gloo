// package: envoy.config.filter.http.aws_lambda.v2
// file: github.com/solo-io/gloo/projects/gloo/api/external/envoy/extensions/aws/filter.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";
import * as validate_validate_pb from "../../../../../../../../../../validate/validate_pb";

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
  }
}

export class AWSLambdaConfig extends jspb.Message {
  hasUseDefaultCredentials(): boolean;
  clearUseDefaultCredentials(): void;
  getUseDefaultCredentials(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setUseDefaultCredentials(value?: google_protobuf_wrappers_pb.BoolValue): void;

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
  }
}

