// package: envoy.api.v2.core
// file: github.com/solo-io/solo-kit/api/external/envoy/api/v2/core/http_uri.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_duration_pb from "google-protobuf/google/protobuf/duration_pb";
import * as validate_validate_pb from "../../../../../../../../../validate/validate_pb";
import * as gogoproto_gogo_pb from "../../../../../../../../../gogoproto/gogo_pb";

export class HttpUri extends jspb.Message {
  getUri(): string;
  setUri(value: string): void;

  hasCluster(): boolean;
  clearCluster(): void;
  getCluster(): string;
  setCluster(value: string): void;

  hasTimeout(): boolean;
  clearTimeout(): void;
  getTimeout(): google_protobuf_duration_pb.Duration | undefined;
  setTimeout(value?: google_protobuf_duration_pb.Duration): void;

  getHttpUpstreamTypeCase(): HttpUri.HttpUpstreamTypeCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HttpUri.AsObject;
  static toObject(includeInstance: boolean, msg: HttpUri): HttpUri.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: HttpUri, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HttpUri;
  static deserializeBinaryFromReader(message: HttpUri, reader: jspb.BinaryReader): HttpUri;
}

export namespace HttpUri {
  export type AsObject = {
    uri: string,
    cluster: string,
    timeout?: google_protobuf_duration_pb.Duration.AsObject,
  }

  export enum HttpUpstreamTypeCase {
    HTTP_UPSTREAM_TYPE_NOT_SET = 0,
    CLUSTER = 2,
  }
}

