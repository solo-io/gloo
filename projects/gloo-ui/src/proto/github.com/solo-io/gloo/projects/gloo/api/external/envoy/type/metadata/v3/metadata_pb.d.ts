/* eslint-disable */
// package: solo.io.envoy.type.metadata.v3
// file: github.com/solo-io/gloo/projects/gloo/api/external/envoy/type/metadata/v3/metadata.proto

import * as jspb from "google-protobuf";
import * as udpa_annotations_status_pb from "../../../../../../../../../../../udpa/annotations/status_pb";
import * as udpa_annotations_versioning_pb from "../../../../../../../../../../../udpa/annotations/versioning_pb";
import * as validate_validate_pb from "../../../../../../../../../../../validate/validate_pb";
import * as extproto_ext_pb from "../../../../../../../../../../../protoc-gen-ext/extproto/ext_pb";

export class MetadataKey extends jspb.Message {
  getKey(): string;
  setKey(value: string): void;

  clearPathList(): void;
  getPathList(): Array<MetadataKey.PathSegment>;
  setPathList(value: Array<MetadataKey.PathSegment>): void;
  addPath(value?: MetadataKey.PathSegment, index?: number): MetadataKey.PathSegment;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MetadataKey.AsObject;
  static toObject(includeInstance: boolean, msg: MetadataKey): MetadataKey.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: MetadataKey, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MetadataKey;
  static deserializeBinaryFromReader(message: MetadataKey, reader: jspb.BinaryReader): MetadataKey;
}

export namespace MetadataKey {
  export type AsObject = {
    key: string,
    pathList: Array<MetadataKey.PathSegment.AsObject>,
  }

  export class PathSegment extends jspb.Message {
    hasKey(): boolean;
    clearKey(): void;
    getKey(): string;
    setKey(value: string): void;

    getSegmentCase(): PathSegment.SegmentCase;
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): PathSegment.AsObject;
    static toObject(includeInstance: boolean, msg: PathSegment): PathSegment.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: PathSegment, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): PathSegment;
    static deserializeBinaryFromReader(message: PathSegment, reader: jspb.BinaryReader): PathSegment;
  }

  export namespace PathSegment {
    export type AsObject = {
      key: string,
    }

    export enum SegmentCase {
      SEGMENT_NOT_SET = 0,
      KEY = 1,
    }
  }
}

export class MetadataKind extends jspb.Message {
  hasRequest(): boolean;
  clearRequest(): void;
  getRequest(): MetadataKind.Request | undefined;
  setRequest(value?: MetadataKind.Request): void;

  hasRoute(): boolean;
  clearRoute(): void;
  getRoute(): MetadataKind.Route | undefined;
  setRoute(value?: MetadataKind.Route): void;

  hasCluster(): boolean;
  clearCluster(): void;
  getCluster(): MetadataKind.Cluster | undefined;
  setCluster(value?: MetadataKind.Cluster): void;

  hasHost(): boolean;
  clearHost(): void;
  getHost(): MetadataKind.Host | undefined;
  setHost(value?: MetadataKind.Host): void;

  getKindCase(): MetadataKind.KindCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MetadataKind.AsObject;
  static toObject(includeInstance: boolean, msg: MetadataKind): MetadataKind.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: MetadataKind, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MetadataKind;
  static deserializeBinaryFromReader(message: MetadataKind, reader: jspb.BinaryReader): MetadataKind;
}

export namespace MetadataKind {
  export type AsObject = {
    request?: MetadataKind.Request.AsObject,
    route?: MetadataKind.Route.AsObject,
    cluster?: MetadataKind.Cluster.AsObject,
    host?: MetadataKind.Host.AsObject,
  }

  export class Request extends jspb.Message {
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Request.AsObject;
    static toObject(includeInstance: boolean, msg: Request): Request.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Request, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Request;
    static deserializeBinaryFromReader(message: Request, reader: jspb.BinaryReader): Request;
  }

  export namespace Request {
    export type AsObject = {
    }
  }

  export class Route extends jspb.Message {
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Route.AsObject;
    static toObject(includeInstance: boolean, msg: Route): Route.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Route, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Route;
    static deserializeBinaryFromReader(message: Route, reader: jspb.BinaryReader): Route;
  }

  export namespace Route {
    export type AsObject = {
    }
  }

  export class Cluster extends jspb.Message {
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Cluster.AsObject;
    static toObject(includeInstance: boolean, msg: Cluster): Cluster.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Cluster, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Cluster;
    static deserializeBinaryFromReader(message: Cluster, reader: jspb.BinaryReader): Cluster;
  }

  export namespace Cluster {
    export type AsObject = {
    }
  }

  export class Host extends jspb.Message {
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Host.AsObject;
    static toObject(includeInstance: boolean, msg: Host): Host.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Host, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Host;
    static deserializeBinaryFromReader(message: Host, reader: jspb.BinaryReader): Host;
  }

  export namespace Host {
    export type AsObject = {
    }
  }

  export enum KindCase {
    KIND_NOT_SET = 0,
    REQUEST = 1,
    ROUTE = 2,
    CLUSTER = 3,
    HOST = 4,
  }
}
