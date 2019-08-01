// package: tracing.plugins.gloo.solo.io
// file: github.com/solo-io/gloo/projects/gloo/api/v1/plugins/tracing/tracing.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../../../../../../../gogoproto/gogo_pb";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";
import * as google_protobuf_duration_pb from "google-protobuf/google/protobuf/duration_pb";

export class ListenerTracingSettings extends jspb.Message {
  clearRequestHeadersForTagsList(): void;
  getRequestHeadersForTagsList(): Array<string>;
  setRequestHeadersForTagsList(value: Array<string>): void;
  addRequestHeadersForTags(value: string, index?: number): string;

  getVerbose(): boolean;
  setVerbose(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListenerTracingSettings.AsObject;
  static toObject(includeInstance: boolean, msg: ListenerTracingSettings): ListenerTracingSettings.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListenerTracingSettings, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListenerTracingSettings;
  static deserializeBinaryFromReader(message: ListenerTracingSettings, reader: jspb.BinaryReader): ListenerTracingSettings;
}

export namespace ListenerTracingSettings {
  export type AsObject = {
    requestHeadersForTagsList: Array<string>,
    verbose: boolean,
  }
}

export class RouteTracingSettings extends jspb.Message {
  getRouteDescriptor(): string;
  setRouteDescriptor(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RouteTracingSettings.AsObject;
  static toObject(includeInstance: boolean, msg: RouteTracingSettings): RouteTracingSettings.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RouteTracingSettings, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RouteTracingSettings;
  static deserializeBinaryFromReader(message: RouteTracingSettings, reader: jspb.BinaryReader): RouteTracingSettings;
}

export namespace RouteTracingSettings {
  export type AsObject = {
    routeDescriptor: string,
  }
}

