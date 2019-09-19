// package: core.solo.io
// file: github.com/solo-io/solo-kit/api/v1/solo-kit.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_descriptor_pb from "google-protobuf/google/protobuf/descriptor_pb";

export class Resource extends jspb.Message {
  getShortName(): string;
  setShortName(value: string): void;

  getPluralName(): string;
  setPluralName(value: string): void;

  getClusterScoped(): boolean;
  setClusterScoped(value: boolean): void;

  getSkipDocsGen(): boolean;
  setSkipDocsGen(value: boolean): void;

  getSkipHashingAnnotations(): boolean;
  setSkipHashingAnnotations(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Resource.AsObject;
  static toObject(includeInstance: boolean, msg: Resource): Resource.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Resource, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Resource;
  static deserializeBinaryFromReader(message: Resource, reader: jspb.BinaryReader): Resource;
}

export namespace Resource {
  export type AsObject = {
    shortName: string,
    pluralName: string,
    clusterScoped: boolean,
    skipDocsGen: boolean,
    skipHashingAnnotations: boolean,
  }
}

  export const resource: jspb.ExtensionFieldInfo<Resource>;

  export const skipHashing: jspb.ExtensionFieldInfo<boolean>;

