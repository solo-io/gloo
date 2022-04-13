/* eslint-disable */
// package: enterprise.graphql.gloo.solo.io
// file: github.com/solo-io/solo-projects/projects/gloo/api/enterprise/graphql/v1/diff.proto

import * as jspb from "google-protobuf";

export class GraphQLInspectorDiffInput extends jspb.Message {
  getOldSchema(): string;
  setOldSchema(value: string): void;

  getNewSchema(): string;
  setNewSchema(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GraphQLInspectorDiffInput.AsObject;
  static toObject(includeInstance: boolean, msg: GraphQLInspectorDiffInput): GraphQLInspectorDiffInput.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GraphQLInspectorDiffInput, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GraphQLInspectorDiffInput;
  static deserializeBinaryFromReader(message: GraphQLInspectorDiffInput, reader: jspb.BinaryReader): GraphQLInspectorDiffInput;
}

export namespace GraphQLInspectorDiffInput {
  export type AsObject = {
    oldSchema: string,
    newSchema: string,
  }
}

export class GraphQLInspectorDiffOutput extends jspb.Message {
  clearChangesList(): void;
  getChangesList(): Array<GraphQLInspectorDiffOutput.Change>;
  setChangesList(value: Array<GraphQLInspectorDiffOutput.Change>): void;
  addChanges(value?: GraphQLInspectorDiffOutput.Change, index?: number): GraphQLInspectorDiffOutput.Change;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GraphQLInspectorDiffOutput.AsObject;
  static toObject(includeInstance: boolean, msg: GraphQLInspectorDiffOutput): GraphQLInspectorDiffOutput.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GraphQLInspectorDiffOutput, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GraphQLInspectorDiffOutput;
  static deserializeBinaryFromReader(message: GraphQLInspectorDiffOutput, reader: jspb.BinaryReader): GraphQLInspectorDiffOutput;
}

export namespace GraphQLInspectorDiffOutput {
  export type AsObject = {
    changesList: Array<GraphQLInspectorDiffOutput.Change.AsObject>,
  }

  export class Change extends jspb.Message {
    getMessage(): string;
    setMessage(value: string): void;

    getPath(): string;
    setPath(value: string): void;

    getChangeType(): string;
    setChangeType(value: string): void;

    hasCriticality(): boolean;
    clearCriticality(): void;
    getCriticality(): GraphQLInspectorDiffOutput.Criticality | undefined;
    setCriticality(value?: GraphQLInspectorDiffOutput.Criticality): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Change.AsObject;
    static toObject(includeInstance: boolean, msg: Change): Change.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Change, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Change;
    static deserializeBinaryFromReader(message: Change, reader: jspb.BinaryReader): Change;
  }

  export namespace Change {
    export type AsObject = {
      message: string,
      path: string,
      changeType: string,
      criticality?: GraphQLInspectorDiffOutput.Criticality.AsObject,
    }
  }

  export class Criticality extends jspb.Message {
    getLevel(): GraphQLInspectorDiffOutput.CriticalityLevelMap[keyof GraphQLInspectorDiffOutput.CriticalityLevelMap];
    setLevel(value: GraphQLInspectorDiffOutput.CriticalityLevelMap[keyof GraphQLInspectorDiffOutput.CriticalityLevelMap]): void;

    getReason(): string;
    setReason(value: string): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Criticality.AsObject;
    static toObject(includeInstance: boolean, msg: Criticality): Criticality.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Criticality, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Criticality;
    static deserializeBinaryFromReader(message: Criticality, reader: jspb.BinaryReader): Criticality;
  }

  export namespace Criticality {
    export type AsObject = {
      level: GraphQLInspectorDiffOutput.CriticalityLevelMap[keyof GraphQLInspectorDiffOutput.CriticalityLevelMap],
      reason: string,
    }
  }

  export interface CriticalityLevelMap {
    NON_BREAKING: 0;
    DANGEROUS: 1;
    BREAKING: 2;
  }

  export const CriticalityLevel: CriticalityLevelMap;
}
