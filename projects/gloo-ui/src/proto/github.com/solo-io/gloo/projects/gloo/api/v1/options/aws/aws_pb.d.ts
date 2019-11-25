// package: aws.options.gloo.solo.io
// file: github.com/solo-io/gloo/projects/gloo/api/v1/options/aws/aws.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../../../../../../../gogoproto/gogo_pb";
import * as github_com_solo_io_solo_kit_api_v1_ref_pb from "../../../../../../../../../github.com/solo-io/solo-kit/api/v1/ref_pb";

export class UpstreamSpec extends jspb.Message {
  getRegion(): string;
  setRegion(value: string): void;

  hasSecretRef(): boolean;
  clearSecretRef(): void;
  getSecretRef(): github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef | undefined;
  setSecretRef(value?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef): void;

  clearLambdaFunctionsList(): void;
  getLambdaFunctionsList(): Array<LambdaFunctionSpec>;
  setLambdaFunctionsList(value: Array<LambdaFunctionSpec>): void;
  addLambdaFunctions(value?: LambdaFunctionSpec, index?: number): LambdaFunctionSpec;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpstreamSpec.AsObject;
  static toObject(includeInstance: boolean, msg: UpstreamSpec): UpstreamSpec.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpstreamSpec, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpstreamSpec;
  static deserializeBinaryFromReader(message: UpstreamSpec, reader: jspb.BinaryReader): UpstreamSpec;
}

export namespace UpstreamSpec {
  export type AsObject = {
    region: string,
    secretRef?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject,
    lambdaFunctionsList: Array<LambdaFunctionSpec.AsObject>,
  }
}

export class LambdaFunctionSpec extends jspb.Message {
  getLogicalName(): string;
  setLogicalName(value: string): void;

  getLambdaFunctionName(): string;
  setLambdaFunctionName(value: string): void;

  getQualifier(): string;
  setQualifier(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LambdaFunctionSpec.AsObject;
  static toObject(includeInstance: boolean, msg: LambdaFunctionSpec): LambdaFunctionSpec.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: LambdaFunctionSpec, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LambdaFunctionSpec;
  static deserializeBinaryFromReader(message: LambdaFunctionSpec, reader: jspb.BinaryReader): LambdaFunctionSpec;
}

export namespace LambdaFunctionSpec {
  export type AsObject = {
    logicalName: string,
    lambdaFunctionName: string,
    qualifier: string,
  }
}

export class DestinationSpec extends jspb.Message {
  getLogicalName(): string;
  setLogicalName(value: string): void;

  getInvocationStyle(): DestinationSpec.InvocationStyleMap[keyof DestinationSpec.InvocationStyleMap];
  setInvocationStyle(value: DestinationSpec.InvocationStyleMap[keyof DestinationSpec.InvocationStyleMap]): void;

  getResponseTransformation(): boolean;
  setResponseTransformation(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DestinationSpec.AsObject;
  static toObject(includeInstance: boolean, msg: DestinationSpec): DestinationSpec.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DestinationSpec, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DestinationSpec;
  static deserializeBinaryFromReader(message: DestinationSpec, reader: jspb.BinaryReader): DestinationSpec;
}

export namespace DestinationSpec {
  export type AsObject = {
    logicalName: string,
    invocationStyle: DestinationSpec.InvocationStyleMap[keyof DestinationSpec.InvocationStyleMap],
    responseTransformation: boolean,
  }

  export interface InvocationStyleMap {
    SYNC: 0;
    ASYNC: 1;
  }

  export const InvocationStyle: InvocationStyleMap;
}

