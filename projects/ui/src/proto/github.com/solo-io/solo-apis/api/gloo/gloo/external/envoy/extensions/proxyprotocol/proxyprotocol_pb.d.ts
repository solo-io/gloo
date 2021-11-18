/* eslint-disable */
// package: envoy.config.listener.proxy_protocol.v3
// file: github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/extensions/proxyprotocol/proxyprotocol.proto

import * as jspb from "google-protobuf";
import * as udpa_annotations_status_pb from "../../../../../../../../../../udpa/annotations/status_pb";
import * as validate_validate_pb from "../../../../../../../../../../validate/validate_pb";

export class CustomProxyProtocol extends jspb.Message {
  clearRulesList(): void;
  getRulesList(): Array<CustomProxyProtocol.Rule>;
  setRulesList(value: Array<CustomProxyProtocol.Rule>): void;
  addRules(value?: CustomProxyProtocol.Rule, index?: number): CustomProxyProtocol.Rule;

  getAllowRequestsWithoutProxyProtocol(): boolean;
  setAllowRequestsWithoutProxyProtocol(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CustomProxyProtocol.AsObject;
  static toObject(includeInstance: boolean, msg: CustomProxyProtocol): CustomProxyProtocol.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CustomProxyProtocol, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CustomProxyProtocol;
  static deserializeBinaryFromReader(message: CustomProxyProtocol, reader: jspb.BinaryReader): CustomProxyProtocol;
}

export namespace CustomProxyProtocol {
  export type AsObject = {
    rulesList: Array<CustomProxyProtocol.Rule.AsObject>,
    allowRequestsWithoutProxyProtocol: boolean,
  }

  export class KeyValuePair extends jspb.Message {
    getMetadataNamespace(): string;
    setMetadataNamespace(value: string): void;

    getKey(): string;
    setKey(value: string): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): KeyValuePair.AsObject;
    static toObject(includeInstance: boolean, msg: KeyValuePair): KeyValuePair.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: KeyValuePair, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): KeyValuePair;
    static deserializeBinaryFromReader(message: KeyValuePair, reader: jspb.BinaryReader): KeyValuePair;
  }

  export namespace KeyValuePair {
    export type AsObject = {
      metadataNamespace: string,
      key: string,
    }
  }

  export class Rule extends jspb.Message {
    getTlvType(): number;
    setTlvType(value: number): void;

    hasOnTlvPresent(): boolean;
    clearOnTlvPresent(): void;
    getOnTlvPresent(): CustomProxyProtocol.KeyValuePair | undefined;
    setOnTlvPresent(value?: CustomProxyProtocol.KeyValuePair): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Rule.AsObject;
    static toObject(includeInstance: boolean, msg: Rule): Rule.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Rule, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Rule;
    static deserializeBinaryFromReader(message: Rule, reader: jspb.BinaryReader): Rule;
  }

  export namespace Rule {
    export type AsObject = {
      tlvType: number,
      onTlvPresent?: CustomProxyProtocol.KeyValuePair.AsObject,
    }
  }
}
