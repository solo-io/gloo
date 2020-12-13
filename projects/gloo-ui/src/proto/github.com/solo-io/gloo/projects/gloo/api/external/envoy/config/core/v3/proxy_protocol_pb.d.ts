/* eslint-disable */
// package: solo.io.envoy.config.core.v3
// file: github.com/solo-io/gloo/projects/gloo/api/external/envoy/config/core/v3/proxy_protocol.proto

import * as jspb from "google-protobuf";
import * as udpa_annotations_status_pb from "../../../../../../../../../../../udpa/annotations/status_pb";
import * as validate_validate_pb from "../../../../../../../../../../../validate/validate_pb";
import * as extproto_ext_pb from "../../../../../../../../../../../protoc-gen-ext/extproto/ext_pb";

export class ProxyProtocolConfig extends jspb.Message {
  getVersion(): ProxyProtocolConfig.VersionMap[keyof ProxyProtocolConfig.VersionMap];
  setVersion(value: ProxyProtocolConfig.VersionMap[keyof ProxyProtocolConfig.VersionMap]): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProxyProtocolConfig.AsObject;
  static toObject(includeInstance: boolean, msg: ProxyProtocolConfig): ProxyProtocolConfig.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ProxyProtocolConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProxyProtocolConfig;
  static deserializeBinaryFromReader(message: ProxyProtocolConfig, reader: jspb.BinaryReader): ProxyProtocolConfig;
}

export namespace ProxyProtocolConfig {
  export type AsObject = {
    version: ProxyProtocolConfig.VersionMap[keyof ProxyProtocolConfig.VersionMap],
  }

  export interface VersionMap {
    V1: 0;
    V2: 1;
  }

  export const Version: VersionMap;
}
