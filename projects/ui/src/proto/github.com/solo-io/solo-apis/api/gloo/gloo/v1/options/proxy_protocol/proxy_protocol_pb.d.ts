/* eslint-disable */
// package: proxy_protocol.options.gloo.solo.io
// file: github.com/solo-io/solo-apis/api/gloo/gloo/v1/options/proxy_protocol/proxy_protocol.proto

import * as jspb from "google-protobuf";
import * as extproto_ext_pb from "../../../../../../../../../extproto/ext_pb";
import * as validate_validate_pb from "../../../../../../../../../validate/validate_pb";

export class ProxyProtocol extends jspb.Message {
  clearRulesList(): void;
  getRulesList(): Array<ProxyProtocol.Rule>;
  setRulesList(value: Array<ProxyProtocol.Rule>): void;
  addRules(value?: ProxyProtocol.Rule, index?: number): ProxyProtocol.Rule;

  getAllowRequestsWithoutProxyProtocol(): boolean;
  setAllowRequestsWithoutProxyProtocol(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProxyProtocol.AsObject;
  static toObject(includeInstance: boolean, msg: ProxyProtocol): ProxyProtocol.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ProxyProtocol, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProxyProtocol;
  static deserializeBinaryFromReader(message: ProxyProtocol, reader: jspb.BinaryReader): ProxyProtocol;
}

export namespace ProxyProtocol {
  export type AsObject = {
    rulesList: Array<ProxyProtocol.Rule.AsObject>,
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
    getOnTlvPresent(): ProxyProtocol.KeyValuePair | undefined;
    setOnTlvPresent(value?: ProxyProtocol.KeyValuePair): void;

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
      onTlvPresent?: ProxyProtocol.KeyValuePair.AsObject,
    }
  }
}
