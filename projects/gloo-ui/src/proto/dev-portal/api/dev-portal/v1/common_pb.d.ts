/* eslint-disable */
// package: devportal.solo.io
// file: dev-portal/api/dev-portal/v1/common.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../../gogoproto/gogo_pb";
import * as extproto_ext_pb from "../../../../protoc-gen-ext/extproto/ext_pb";

export class DataSource extends jspb.Message {
  hasInlineString(): boolean;
  clearInlineString(): void;
  getInlineString(): string;
  setInlineString(value: string): void;

  hasFetchUrl(): boolean;
  clearFetchUrl(): void;
  getFetchUrl(): string;
  setFetchUrl(value: string): void;

  getSourcetypeCase(): DataSource.SourcetypeCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DataSource.AsObject;
  static toObject(includeInstance: boolean, msg: DataSource): DataSource.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DataSource, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DataSource;
  static deserializeBinaryFromReader(message: DataSource, reader: jspb.BinaryReader): DataSource;
}

export namespace DataSource {
  export type AsObject = {
    inlineString: string,
    fetchUrl: string,
  }

  export enum SourcetypeCase {
    SOURCETYPE_NOT_SET = 0,
    INLINE_STRING = 1,
    FETCH_URL = 2,
  }
}
