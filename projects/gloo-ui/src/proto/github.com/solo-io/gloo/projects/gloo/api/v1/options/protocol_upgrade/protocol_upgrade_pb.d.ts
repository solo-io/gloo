// package: protocol_upgrade.options.gloo.solo.io
// file: github.com/solo-io/gloo/projects/gloo/api/v1/options/protocol_upgrade/protocol_upgrade.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";
import * as gogoproto_gogo_pb from "../../../../../../../../../gogoproto/gogo_pb";

export class ProtocolUpgradeConfig extends jspb.Message {
  hasWebsocket(): boolean;
  clearWebsocket(): void;
  getWebsocket(): ProtocolUpgradeConfig.ProtocolUpgradeSpec | undefined;
  setWebsocket(value?: ProtocolUpgradeConfig.ProtocolUpgradeSpec): void;

  getUpgradeTypeCase(): ProtocolUpgradeConfig.UpgradeTypeCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProtocolUpgradeConfig.AsObject;
  static toObject(includeInstance: boolean, msg: ProtocolUpgradeConfig): ProtocolUpgradeConfig.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ProtocolUpgradeConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProtocolUpgradeConfig;
  static deserializeBinaryFromReader(message: ProtocolUpgradeConfig, reader: jspb.BinaryReader): ProtocolUpgradeConfig;
}

export namespace ProtocolUpgradeConfig {
  export type AsObject = {
    websocket?: ProtocolUpgradeConfig.ProtocolUpgradeSpec.AsObject,
  }

  export class ProtocolUpgradeSpec extends jspb.Message {
    hasEnabled(): boolean;
    clearEnabled(): void;
    getEnabled(): google_protobuf_wrappers_pb.BoolValue | undefined;
    setEnabled(value?: google_protobuf_wrappers_pb.BoolValue): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ProtocolUpgradeSpec.AsObject;
    static toObject(includeInstance: boolean, msg: ProtocolUpgradeSpec): ProtocolUpgradeSpec.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: ProtocolUpgradeSpec, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ProtocolUpgradeSpec;
    static deserializeBinaryFromReader(message: ProtocolUpgradeSpec, reader: jspb.BinaryReader): ProtocolUpgradeSpec;
  }

  export namespace ProtocolUpgradeSpec {
    export type AsObject = {
      enabled?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    }
  }

  export enum UpgradeTypeCase {
    UPGRADE_TYPE_NOT_SET = 0,
    WEBSOCKET = 1,
  }
}

