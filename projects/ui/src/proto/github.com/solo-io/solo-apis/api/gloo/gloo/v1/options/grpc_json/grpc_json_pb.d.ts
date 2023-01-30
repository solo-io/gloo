/* eslint-disable */
// package: grpc_json.options.gloo.solo.io
// file: github.com/solo-io/solo-apis/api/gloo/gloo/v1/options/grpc_json/grpc_json.proto

import * as jspb from "google-protobuf";
import * as github_com_solo_io_solo_kit_api_v1_ref_pb from "../../../../../../../../../github.com/solo-io/solo-kit/api/v1/ref_pb";
import * as extproto_ext_pb from "../../../../../../../../../extproto/ext_pb";
import * as validate_validate_pb from "../../../../../../../../../validate/validate_pb";

export class GrpcJsonTranscoder extends jspb.Message {
  hasProtoDescriptor(): boolean;
  clearProtoDescriptor(): void;
  getProtoDescriptor(): string;
  setProtoDescriptor(value: string): void;

  hasProtoDescriptorBin(): boolean;
  clearProtoDescriptorBin(): void;
  getProtoDescriptorBin(): Uint8Array | string;
  getProtoDescriptorBin_asU8(): Uint8Array;
  getProtoDescriptorBin_asB64(): string;
  setProtoDescriptorBin(value: Uint8Array | string): void;

  hasProtoDescriptorConfigMap(): boolean;
  clearProtoDescriptorConfigMap(): void;
  getProtoDescriptorConfigMap(): GrpcJsonTranscoder.DescriptorConfigMap | undefined;
  setProtoDescriptorConfigMap(value?: GrpcJsonTranscoder.DescriptorConfigMap): void;

  clearServicesList(): void;
  getServicesList(): Array<string>;
  setServicesList(value: Array<string>): void;
  addServices(value: string, index?: number): string;

  hasPrintOptions(): boolean;
  clearPrintOptions(): void;
  getPrintOptions(): GrpcJsonTranscoder.PrintOptions | undefined;
  setPrintOptions(value?: GrpcJsonTranscoder.PrintOptions): void;

  getMatchIncomingRequestRoute(): boolean;
  setMatchIncomingRequestRoute(value: boolean): void;

  clearIgnoredQueryParametersList(): void;
  getIgnoredQueryParametersList(): Array<string>;
  setIgnoredQueryParametersList(value: Array<string>): void;
  addIgnoredQueryParameters(value: string, index?: number): string;

  getAutoMapping(): boolean;
  setAutoMapping(value: boolean): void;

  getIgnoreUnknownQueryParameters(): boolean;
  setIgnoreUnknownQueryParameters(value: boolean): void;

  getConvertGrpcStatus(): boolean;
  setConvertGrpcStatus(value: boolean): void;

  getDescriptorSetCase(): GrpcJsonTranscoder.DescriptorSetCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GrpcJsonTranscoder.AsObject;
  static toObject(includeInstance: boolean, msg: GrpcJsonTranscoder): GrpcJsonTranscoder.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GrpcJsonTranscoder, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GrpcJsonTranscoder;
  static deserializeBinaryFromReader(message: GrpcJsonTranscoder, reader: jspb.BinaryReader): GrpcJsonTranscoder;
}

export namespace GrpcJsonTranscoder {
  export type AsObject = {
    protoDescriptor: string,
    protoDescriptorBin: Uint8Array | string,
    protoDescriptorConfigMap?: GrpcJsonTranscoder.DescriptorConfigMap.AsObject,
    servicesList: Array<string>,
    printOptions?: GrpcJsonTranscoder.PrintOptions.AsObject,
    matchIncomingRequestRoute: boolean,
    ignoredQueryParametersList: Array<string>,
    autoMapping: boolean,
    ignoreUnknownQueryParameters: boolean,
    convertGrpcStatus: boolean,
  }

  export class PrintOptions extends jspb.Message {
    getAddWhitespace(): boolean;
    setAddWhitespace(value: boolean): void;

    getAlwaysPrintPrimitiveFields(): boolean;
    setAlwaysPrintPrimitiveFields(value: boolean): void;

    getAlwaysPrintEnumsAsInts(): boolean;
    setAlwaysPrintEnumsAsInts(value: boolean): void;

    getPreserveProtoFieldNames(): boolean;
    setPreserveProtoFieldNames(value: boolean): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): PrintOptions.AsObject;
    static toObject(includeInstance: boolean, msg: PrintOptions): PrintOptions.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: PrintOptions, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): PrintOptions;
    static deserializeBinaryFromReader(message: PrintOptions, reader: jspb.BinaryReader): PrintOptions;
  }

  export namespace PrintOptions {
    export type AsObject = {
      addWhitespace: boolean,
      alwaysPrintPrimitiveFields: boolean,
      alwaysPrintEnumsAsInts: boolean,
      preserveProtoFieldNames: boolean,
    }
  }

  export class DescriptorConfigMap extends jspb.Message {
    hasConfigMapRef(): boolean;
    clearConfigMapRef(): void;
    getConfigMapRef(): github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef | undefined;
    setConfigMapRef(value?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef): void;

    getKey(): string;
    setKey(value: string): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): DescriptorConfigMap.AsObject;
    static toObject(includeInstance: boolean, msg: DescriptorConfigMap): DescriptorConfigMap.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: DescriptorConfigMap, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): DescriptorConfigMap;
    static deserializeBinaryFromReader(message: DescriptorConfigMap, reader: jspb.BinaryReader): DescriptorConfigMap;
  }

  export namespace DescriptorConfigMap {
    export type AsObject = {
      configMapRef?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject,
      key: string,
    }
  }

  export enum DescriptorSetCase {
    DESCRIPTOR_SET_NOT_SET = 0,
    PROTO_DESCRIPTOR = 1,
    PROTO_DESCRIPTOR_BIN = 4,
    PROTO_DESCRIPTOR_CONFIG_MAP = 10,
  }
}
