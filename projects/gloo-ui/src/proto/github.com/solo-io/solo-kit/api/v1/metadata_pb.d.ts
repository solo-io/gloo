// package: core.solo.io
// file: github.com/solo-io/solo-kit/api/v1/metadata.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../../../gogoproto/gogo_pb";

export class Metadata extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getNamespace(): string;
  setNamespace(value: string): void;

  getCluster(): string;
  setCluster(value: string): void;

  getResourceVersion(): string;
  setResourceVersion(value: string): void;

  getLabelsMap(): jspb.Map<string, string>;
  clearLabelsMap(): void;
  getAnnotationsMap(): jspb.Map<string, string>;
  clearAnnotationsMap(): void;
  getGeneration(): number;
  setGeneration(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Metadata.AsObject;
  static toObject(includeInstance: boolean, msg: Metadata): Metadata.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Metadata, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Metadata;
  static deserializeBinaryFromReader(message: Metadata, reader: jspb.BinaryReader): Metadata;
}

export namespace Metadata {
  export type AsObject = {
    name: string,
    namespace: string,
    cluster: string,
    resourceVersion: string,
    labelsMap: Array<[string, string]>,
    annotationsMap: Array<[string, string]>,
    generation: number,
  }
}

