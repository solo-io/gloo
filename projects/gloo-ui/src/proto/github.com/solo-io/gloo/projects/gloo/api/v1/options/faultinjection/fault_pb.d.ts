// package: fault.options.gloo.solo.io
// file: github.com/solo-io/gloo/projects/gloo/api/v1/options/faultinjection/fault.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_duration_pb from "google-protobuf/google/protobuf/duration_pb";
import * as validate_validate_pb from "../../../../../../../../../validate/validate_pb";
import * as gogoproto_gogo_pb from "../../../../../../../../../gogoproto/gogo_pb";

export class RouteAbort extends jspb.Message {
  getPercentage(): number;
  setPercentage(value: number): void;

  getHttpStatus(): number;
  setHttpStatus(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RouteAbort.AsObject;
  static toObject(includeInstance: boolean, msg: RouteAbort): RouteAbort.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RouteAbort, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RouteAbort;
  static deserializeBinaryFromReader(message: RouteAbort, reader: jspb.BinaryReader): RouteAbort;
}

export namespace RouteAbort {
  export type AsObject = {
    percentage: number,
    httpStatus: number,
  }
}

export class RouteDelay extends jspb.Message {
  getPercentage(): number;
  setPercentage(value: number): void;

  hasFixedDelay(): boolean;
  clearFixedDelay(): void;
  getFixedDelay(): google_protobuf_duration_pb.Duration | undefined;
  setFixedDelay(value?: google_protobuf_duration_pb.Duration): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RouteDelay.AsObject;
  static toObject(includeInstance: boolean, msg: RouteDelay): RouteDelay.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RouteDelay, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RouteDelay;
  static deserializeBinaryFromReader(message: RouteDelay, reader: jspb.BinaryReader): RouteDelay;
}

export namespace RouteDelay {
  export type AsObject = {
    percentage: number,
    fixedDelay?: google_protobuf_duration_pb.Duration.AsObject,
  }
}

export class RouteFaults extends jspb.Message {
  hasAbort(): boolean;
  clearAbort(): void;
  getAbort(): RouteAbort | undefined;
  setAbort(value?: RouteAbort): void;

  hasDelay(): boolean;
  clearDelay(): void;
  getDelay(): RouteDelay | undefined;
  setDelay(value?: RouteDelay): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RouteFaults.AsObject;
  static toObject(includeInstance: boolean, msg: RouteFaults): RouteFaults.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RouteFaults, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RouteFaults;
  static deserializeBinaryFromReader(message: RouteFaults, reader: jspb.BinaryReader): RouteFaults;
}

export namespace RouteFaults {
  export type AsObject = {
    abort?: RouteAbort.AsObject,
    delay?: RouteDelay.AsObject,
  }
}

