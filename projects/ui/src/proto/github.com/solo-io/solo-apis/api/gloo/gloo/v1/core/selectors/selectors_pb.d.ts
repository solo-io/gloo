/* eslint-disable */
// package: selectors.core.gloo.solo.io
// file: github.com/solo-io/solo-apis/api/gloo/gloo/v1/core/selectors/selectors.proto

import * as jspb from "google-protobuf";
import * as extproto_ext_pb from "../../../../../../../../../extproto/ext_pb";

export class Selector extends jspb.Message {
  clearNamespacesList(): void;
  getNamespacesList(): Array<string>;
  setNamespacesList(value: Array<string>): void;
  addNamespaces(value: string, index?: number): string;

  getLabelsMap(): jspb.Map<string, string>;
  clearLabelsMap(): void;
  clearExpressionsList(): void;
  getExpressionsList(): Array<Selector.Expression>;
  setExpressionsList(value: Array<Selector.Expression>): void;
  addExpressions(value?: Selector.Expression, index?: number): Selector.Expression;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Selector.AsObject;
  static toObject(includeInstance: boolean, msg: Selector): Selector.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Selector, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Selector;
  static deserializeBinaryFromReader(message: Selector, reader: jspb.BinaryReader): Selector;
}

export namespace Selector {
  export type AsObject = {
    namespacesList: Array<string>,
    labelsMap: Array<[string, string]>,
    expressionsList: Array<Selector.Expression.AsObject>,
  }

  export class Expression extends jspb.Message {
    getKey(): string;
    setKey(value: string): void;

    getOperator(): Selector.Expression.OperatorMap[keyof Selector.Expression.OperatorMap];
    setOperator(value: Selector.Expression.OperatorMap[keyof Selector.Expression.OperatorMap]): void;

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
      operator: Selector.Expression.OperatorMap[keyof Selector.Expression.OperatorMap],
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
