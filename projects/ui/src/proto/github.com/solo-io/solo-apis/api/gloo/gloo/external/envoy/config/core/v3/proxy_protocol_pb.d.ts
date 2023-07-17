/* eslint-disable */
// package: solo.io.envoy.config.core.v3
// file: github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/config/core/v3/proxy_protocol.proto

import * as jspb from "google-protobuf";
import * as udpa_annotations_status_pb from "../../../../../../../../../../../udpa/annotations/status_pb";
import * as validate_validate_pb from "../../../../../../../../../../../validate/validate_pb";
import * as extproto_ext_pb from "../../../../../../../../../../../extproto/ext_pb";

export class ProxyProtocolPassThroughTLVs extends jspb.Message {
  getMatchType(): ProxyProtocolPassThroughTLVs.PassTLVsMatchTypeMap[keyof ProxyProtocolPassThroughTLVs.PassTLVsMatchTypeMap];
  setMatchType(value: ProxyProtocolPassThroughTLVs.PassTLVsMatchTypeMap[keyof ProxyProtocolPassThroughTLVs.PassTLVsMatchTypeMap]): void;

  clearTlvTypeList(): void;
  getTlvTypeList(): Array<number>;
  setTlvTypeList(value: Array<number>): void;
  addTlvType(value: number, index?: number): number;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProxyProtocolPassThroughTLVs.AsObject;
  static toObject(includeInstance: boolean, msg: ProxyProtocolPassThroughTLVs): ProxyProtocolPassThroughTLVs.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ProxyProtocolPassThroughTLVs, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProxyProtocolPassThroughTLVs;
  static deserializeBinaryFromReader(message: ProxyProtocolPassThroughTLVs, reader: jspb.BinaryReader): ProxyProtocolPassThroughTLVs;
}

export namespace ProxyProtocolPassThroughTLVs {
  export type AsObject = {
    matchType: ProxyProtocolPassThroughTLVs.PassTLVsMatchTypeMap[keyof ProxyProtocolPassThroughTLVs.PassTLVsMatchTypeMap],
    tlvTypeList: Array<number>,
  }

  export interface PassTLVsMatchTypeMap {
    INCLUDE_ALL: 0;
    INCLUDE: 1;
  }

  export const PassTLVsMatchType: PassTLVsMatchTypeMap;
}

export class ProxyProtocolConfig extends jspb.Message {
  getVersion(): ProxyProtocolConfig.VersionMap[keyof ProxyProtocolConfig.VersionMap];
  setVersion(value: ProxyProtocolConfig.VersionMap[keyof ProxyProtocolConfig.VersionMap]): void;

  hasPassThroughTlvs(): boolean;
  clearPassThroughTlvs(): void;
  getPassThroughTlvs(): ProxyProtocolPassThroughTLVs | undefined;
  setPassThroughTlvs(value?: ProxyProtocolPassThroughTLVs): void;

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
    passThroughTlvs?: ProxyProtocolPassThroughTLVs.AsObject,
  }

  export interface VersionMap {
    V1: 0;
    V2: 1;
  }

  export const Version: VersionMap;
}
