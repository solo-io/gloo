/* eslint-disable */
// package: tap.options.gloo.solo.io
// file: github.com/solo-io/solo-apis/api/gloo/gloo/v1/enterprise/options/tap/tap.proto

import * as jspb from "google-protobuf";
import * as github_com_solo_io_solo_kit_api_v1_ref_pb from "../../../../../../../../../../github.com/solo-io/solo-kit/api/v1/ref_pb";
import * as validate_validate_pb from "../../../../../../../../../../validate/validate_pb";
import * as google_protobuf_duration_pb from "google-protobuf/google/protobuf/duration_pb";

export class Tap extends jspb.Message {
  clearSinksList(): void;
  getSinksList(): Array<Sink>;
  setSinksList(value: Array<Sink>): void;
  addSinks(value?: Sink, index?: number): Sink;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Tap.AsObject;
  static toObject(includeInstance: boolean, msg: Tap): Tap.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Tap, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Tap;
  static deserializeBinaryFromReader(message: Tap, reader: jspb.BinaryReader): Tap;
}

export namespace Tap {
  export type AsObject = {
    sinksList: Array<Sink.AsObject>,
  }
}

export class Sink extends jspb.Message {
  hasGrpcService(): boolean;
  clearGrpcService(): void;
  getGrpcService(): GrpcService | undefined;
  setGrpcService(value?: GrpcService): void;

  hasHttpService(): boolean;
  clearHttpService(): void;
  getHttpService(): HttpService | undefined;
  setHttpService(value?: HttpService): void;

  getSinktypeCase(): Sink.SinktypeCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Sink.AsObject;
  static toObject(includeInstance: boolean, msg: Sink): Sink.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Sink, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Sink;
  static deserializeBinaryFromReader(message: Sink, reader: jspb.BinaryReader): Sink;
}

export namespace Sink {
  export type AsObject = {
    grpcService?: GrpcService.AsObject,
    httpService?: HttpService.AsObject,
  }

  export enum SinktypeCase {
    SINKTYPE_NOT_SET = 0,
    GRPC_SERVICE = 1,
    HTTP_SERVICE = 2,
  }
}

export class GrpcService extends jspb.Message {
  hasTapServer(): boolean;
  clearTapServer(): void;
  getTapServer(): github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef | undefined;
  setTapServer(value?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef): void;

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
    tapServer?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject,
  }
}

export class HttpService extends jspb.Message {
  hasTapServer(): boolean;
  clearTapServer(): void;
  getTapServer(): github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef | undefined;
  setTapServer(value?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef): void;

  hasTimeout(): boolean;
  clearTimeout(): void;
  getTimeout(): google_protobuf_duration_pb.Duration | undefined;
  setTimeout(value?: google_protobuf_duration_pb.Duration): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HttpService.AsObject;
  static toObject(includeInstance: boolean, msg: HttpService): HttpService.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: HttpService, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HttpService;
  static deserializeBinaryFromReader(message: HttpService, reader: jspb.BinaryReader): HttpService;
}

export namespace HttpService {
  export type AsObject = {
    tapServer?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject,
    timeout?: google_protobuf_duration_pb.Duration.AsObject,
  }
}
