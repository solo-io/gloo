// package: envoy.config.filter.http.aws.v2
// file: github.com/solo-io/gloo/projects/gloo/api/external/envoy/extensions/aws/filter.proto

import * as jspb from "google-protobuf";
import * as validate_validate_pb from "../../../../../../../../../../validate/validate_pb";

export class LambdaPerRoute extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getQualifier(): string;
  setQualifier(value: string): void;

  getAsync(): boolean;
  setAsync(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LambdaPerRoute.AsObject;
  static toObject(includeInstance: boolean, msg: LambdaPerRoute): LambdaPerRoute.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: LambdaPerRoute, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LambdaPerRoute;
  static deserializeBinaryFromReader(message: LambdaPerRoute, reader: jspb.BinaryReader): LambdaPerRoute;
}

export namespace LambdaPerRoute {
  export type AsObject = {
    name: string,
    qualifier: string,
    async: boolean,
  }
}

export class LambdaProtocolExtension extends jspb.Message {
  getHost(): string;
  setHost(value: string): void;

  getRegion(): string;
  setRegion(value: string): void;

  getAccessKey(): string;
  setAccessKey(value: string): void;

  getSecretKey(): string;
  setSecretKey(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LambdaProtocolExtension.AsObject;
  static toObject(includeInstance: boolean, msg: LambdaProtocolExtension): LambdaProtocolExtension.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: LambdaProtocolExtension, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LambdaProtocolExtension;
  static deserializeBinaryFromReader(message: LambdaProtocolExtension, reader: jspb.BinaryReader): LambdaProtocolExtension;
}

export namespace LambdaProtocolExtension {
  export type AsObject = {
    host: string,
    region: string,
    accessKey: string,
    secretKey: string,
  }
}

