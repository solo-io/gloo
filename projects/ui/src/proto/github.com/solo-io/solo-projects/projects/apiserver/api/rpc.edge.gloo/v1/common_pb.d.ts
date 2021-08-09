/* eslint-disable */
// package: rpc.edge.gloo.solo.io
// file: github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/common.proto

import * as jspb from "google-protobuf";
import * as extproto_ext_pb from "../../../../../../../../extproto/ext_pb";

export class ObjectMeta extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getNamespace(): string;
  setNamespace(value: string): void;

  getUid(): string;
  setUid(value: string): void;

  getResourceVersion(): string;
  setResourceVersion(value: string): void;

  hasCreationTimestamp(): boolean;
  clearCreationTimestamp(): void;
  getCreationTimestamp(): Time | undefined;
  setCreationTimestamp(value?: Time): void;

  getLabelsMap(): jspb.Map<string, string>;
  clearLabelsMap(): void;
  getAnnotationsMap(): jspb.Map<string, string>;
  clearAnnotationsMap(): void;
  getClusterName(): string;
  setClusterName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ObjectMeta.AsObject;
  static toObject(includeInstance: boolean, msg: ObjectMeta): ObjectMeta.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ObjectMeta, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ObjectMeta;
  static deserializeBinaryFromReader(message: ObjectMeta, reader: jspb.BinaryReader): ObjectMeta;
}

export namespace ObjectMeta {
  export type AsObject = {
    name: string,
    namespace: string,
    uid: string,
    resourceVersion: string,
    creationTimestamp?: Time.AsObject,
    labelsMap: Array<[string, string]>,
    annotationsMap: Array<[string, string]>,
    clusterName: string,
  }
}

export class Time extends jspb.Message {
  getSeconds(): number;
  setSeconds(value: number): void;

  getNanos(): number;
  setNanos(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Time.AsObject;
  static toObject(includeInstance: boolean, msg: Time): Time.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Time, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Time;
  static deserializeBinaryFromReader(message: Time, reader: jspb.BinaryReader): Time;
}

export namespace Time {
  export type AsObject = {
    seconds: number,
    nanos: number,
  }
}

export class ResourceYaml extends jspb.Message {
  getYaml(): string;
  setYaml(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResourceYaml.AsObject;
  static toObject(includeInstance: boolean, msg: ResourceYaml): ResourceYaml.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ResourceYaml, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResourceYaml;
  static deserializeBinaryFromReader(message: ResourceYaml, reader: jspb.BinaryReader): ResourceYaml;
}

export namespace ResourceYaml {
  export type AsObject = {
    yaml: string,
  }
}
