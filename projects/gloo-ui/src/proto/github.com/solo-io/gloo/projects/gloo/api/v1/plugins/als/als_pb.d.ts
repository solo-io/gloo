// package: als.plugins.gloo.solo.io
// file: github.com/solo-io/gloo/projects/gloo/api/v1/plugins/als/als.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../../../../../../../gogoproto/gogo_pb";
import * as google_protobuf_struct_pb from "google-protobuf/google/protobuf/struct_pb";

export class AccessLoggingService extends jspb.Message {
  clearAccessLogList(): void;
  getAccessLogList(): Array<AccessLog>;
  setAccessLogList(value: Array<AccessLog>): void;
  addAccessLog(value?: AccessLog, index?: number): AccessLog;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AccessLoggingService.AsObject;
  static toObject(includeInstance: boolean, msg: AccessLoggingService): AccessLoggingService.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: AccessLoggingService, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AccessLoggingService;
  static deserializeBinaryFromReader(message: AccessLoggingService, reader: jspb.BinaryReader): AccessLoggingService;
}

export namespace AccessLoggingService {
  export type AsObject = {
    accessLogList: Array<AccessLog.AsObject>,
  }
}

export class AccessLog extends jspb.Message {
  hasFileSink(): boolean;
  clearFileSink(): void;
  getFileSink(): FileSink | undefined;
  setFileSink(value?: FileSink): void;

  getOutputdestinationCase(): AccessLog.OutputdestinationCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AccessLog.AsObject;
  static toObject(includeInstance: boolean, msg: AccessLog): AccessLog.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: AccessLog, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AccessLog;
  static deserializeBinaryFromReader(message: AccessLog, reader: jspb.BinaryReader): AccessLog;
}

export namespace AccessLog {
  export type AsObject = {
    fileSink?: FileSink.AsObject,
  }

  export enum OutputdestinationCase {
    OUTPUTDESTINATION_NOT_SET = 0,
    FILE_SINK = 2,
  }
}

export class FileSink extends jspb.Message {
  getPath(): string;
  setPath(value: string): void;

  hasStringFormat(): boolean;
  clearStringFormat(): void;
  getStringFormat(): string;
  setStringFormat(value: string): void;

  hasJsonFormat(): boolean;
  clearJsonFormat(): void;
  getJsonFormat(): google_protobuf_struct_pb.Struct | undefined;
  setJsonFormat(value?: google_protobuf_struct_pb.Struct): void;

  getOutputFormatCase(): FileSink.OutputFormatCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FileSink.AsObject;
  static toObject(includeInstance: boolean, msg: FileSink): FileSink.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: FileSink, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FileSink;
  static deserializeBinaryFromReader(message: FileSink, reader: jspb.BinaryReader): FileSink;
}

export namespace FileSink {
  export type AsObject = {
    path: string,
    stringFormat: string,
    jsonFormat?: google_protobuf_struct_pb.Struct.AsObject,
  }

  export enum OutputFormatCase {
    OUTPUT_FORMAT_NOT_SET = 0,
    STRING_FORMAT = 2,
    JSON_FORMAT = 3,
  }
}

