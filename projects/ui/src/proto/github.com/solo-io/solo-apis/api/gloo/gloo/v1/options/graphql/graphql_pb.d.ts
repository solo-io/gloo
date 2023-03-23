/* eslint-disable */
// package: graphql.options.gloo.solo.io
// file: github.com/solo-io/solo-apis/api/gloo/gloo/v1/options/graphql/graphql.proto

import * as jspb from "google-protobuf";
import * as extproto_ext_pb from "../../../../../../../../../extproto/ext_pb";

export class ServiceSpec extends jspb.Message {
  hasEndpoint(): boolean;
  clearEndpoint(): void;
  getEndpoint(): ServiceSpec.Endpoint | undefined;
  setEndpoint(value?: ServiceSpec.Endpoint): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ServiceSpec.AsObject;
  static toObject(includeInstance: boolean, msg: ServiceSpec): ServiceSpec.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ServiceSpec, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ServiceSpec;
  static deserializeBinaryFromReader(message: ServiceSpec, reader: jspb.BinaryReader): ServiceSpec;
}

export namespace ServiceSpec {
  export type AsObject = {
    endpoint?: ServiceSpec.Endpoint.AsObject,
  }

  export class Endpoint extends jspb.Message {
    getUrl(): string;
    setUrl(value: string): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Endpoint.AsObject;
    static toObject(includeInstance: boolean, msg: Endpoint): Endpoint.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Endpoint, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Endpoint;
    static deserializeBinaryFromReader(message: Endpoint, reader: jspb.BinaryReader): Endpoint;
  }

  export namespace Endpoint {
    export type AsObject = {
      url: string,
    }
  }
}
