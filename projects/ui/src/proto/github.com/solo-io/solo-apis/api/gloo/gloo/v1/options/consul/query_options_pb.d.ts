/* eslint-disable */
// package: consul.options.gloo.solo.io
// file: github.com/solo-io/solo-apis/api/gloo/gloo/v1/options/consul/query_options.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_duration_pb from "google-protobuf/google/protobuf/duration_pb";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";
import * as extproto_ext_pb from "../../../../../../../../../extproto/ext_pb";

export class QueryOptions extends jspb.Message {
  hasUseCache(): boolean;
  clearUseCache(): void;
  getUseCache(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setUseCache(value?: google_protobuf_wrappers_pb.BoolValue): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): QueryOptions.AsObject;
  static toObject(includeInstance: boolean, msg: QueryOptions): QueryOptions.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: QueryOptions, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): QueryOptions;
  static deserializeBinaryFromReader(message: QueryOptions, reader: jspb.BinaryReader): QueryOptions;
}

export namespace QueryOptions {
  export type AsObject = {
    useCache?: google_protobuf_wrappers_pb.BoolValue.AsObject,
  }
}

export interface ConsulConsistencyModesMap {
  DEFAULTMODE: 0;
  STALEMODE: 1;
  CONSISTENTMODE: 2;
}

export const ConsulConsistencyModes: ConsulConsistencyModesMap;
