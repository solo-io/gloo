/* eslint-disable */
// package: gateway.solo.io
// file: github.com/solo-io/solo-apis/api/gloo/gateway/v1/http_gateway.proto

import * as jspb from "google-protobuf";
import * as extproto_ext_pb from "../../../../../../../extproto/ext_pb";
import * as github_com_solo_io_solo_kit_api_v1_ref_pb from "../../../../../../../github.com/solo-io/solo-kit/api/v1/ref_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_proxy_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/proxy_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/options_pb";

export class HttpGateway extends jspb.Message {
  clearVirtualServicesList(): void;
  getVirtualServicesList(): Array<github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef>;
  setVirtualServicesList(value: Array<github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef>): void;
  addVirtualServices(value?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef, index?: number): github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef;

  getVirtualServiceSelectorMap(): jspb.Map<string, string>;
  clearVirtualServiceSelectorMap(): void;
  hasVirtualServiceExpressions(): boolean;
  clearVirtualServiceExpressions(): void;
  getVirtualServiceExpressions(): VirtualServiceSelectorExpressions | undefined;
  setVirtualServiceExpressions(value?: VirtualServiceSelectorExpressions): void;

  clearVirtualServiceNamespacesList(): void;
  getVirtualServiceNamespacesList(): Array<string>;
  setVirtualServiceNamespacesList(value: Array<string>): void;
  addVirtualServiceNamespaces(value: string, index?: number): string;

  hasOptions(): boolean;
  clearOptions(): void;
  getOptions(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.HttpListenerOptions | undefined;
  setOptions(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.HttpListenerOptions): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HttpGateway.AsObject;
  static toObject(includeInstance: boolean, msg: HttpGateway): HttpGateway.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: HttpGateway, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HttpGateway;
  static deserializeBinaryFromReader(message: HttpGateway, reader: jspb.BinaryReader): HttpGateway;
}

export namespace HttpGateway {
  export type AsObject = {
    virtualServicesList: Array<github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject>,
    virtualServiceSelectorMap: Array<[string, string]>,
    virtualServiceExpressions?: VirtualServiceSelectorExpressions.AsObject,
    virtualServiceNamespacesList: Array<string>,
    options?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.HttpListenerOptions.AsObject,
  }
}

export class VirtualServiceSelectorExpressions extends jspb.Message {
  clearExpressionsList(): void;
  getExpressionsList(): Array<VirtualServiceSelectorExpressions.Expression>;
  setExpressionsList(value: Array<VirtualServiceSelectorExpressions.Expression>): void;
  addExpressions(value?: VirtualServiceSelectorExpressions.Expression, index?: number): VirtualServiceSelectorExpressions.Expression;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VirtualServiceSelectorExpressions.AsObject;
  static toObject(includeInstance: boolean, msg: VirtualServiceSelectorExpressions): VirtualServiceSelectorExpressions.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: VirtualServiceSelectorExpressions, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VirtualServiceSelectorExpressions;
  static deserializeBinaryFromReader(message: VirtualServiceSelectorExpressions, reader: jspb.BinaryReader): VirtualServiceSelectorExpressions;
}

export namespace VirtualServiceSelectorExpressions {
  export type AsObject = {
    expressionsList: Array<VirtualServiceSelectorExpressions.Expression.AsObject>,
  }

  export class Expression extends jspb.Message {
    getKey(): string;
    setKey(value: string): void;

    getOperator(): VirtualServiceSelectorExpressions.Expression.OperatorMap[keyof VirtualServiceSelectorExpressions.Expression.OperatorMap];
    setOperator(value: VirtualServiceSelectorExpressions.Expression.OperatorMap[keyof VirtualServiceSelectorExpressions.Expression.OperatorMap]): void;

    clearValuesList(): void;
    getValuesList(): Array<string>;
    setValuesList(value: Array<string>): void;
    addValues(value: string, index?: number): string;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Expression.AsObject;
    static toObject(includeInstance: boolean, msg: Expression): Expression.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Expression, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Expression;
    static deserializeBinaryFromReader(message: Expression, reader: jspb.BinaryReader): Expression;
  }

  export namespace Expression {
    export type AsObject = {
      key: string,
      operator: VirtualServiceSelectorExpressions.Expression.OperatorMap[keyof VirtualServiceSelectorExpressions.Expression.OperatorMap],
      valuesList: Array<string>,
    }

    export interface OperatorMap {
      EQUALS: 0;
      DOUBLEEQUALS: 1;
      NOTEQUALS: 2;
      IN: 3;
      NOTIN: 4;
      EXISTS: 5;
      DOESNOTEXIST: 6;
      GREATERTHAN: 7;
      LESSTHAN: 8;
    }

    export const Operator: OperatorMap;
  }
}
