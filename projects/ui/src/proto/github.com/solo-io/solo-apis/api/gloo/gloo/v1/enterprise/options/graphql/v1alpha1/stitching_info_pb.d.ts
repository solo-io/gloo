/* eslint-disable */
// package: graphql.gloo.solo.io
// file: github.com/solo-io/solo-apis/api/gloo/gloo/v1/enterprise/options/graphql/v1alpha1/stitching_info.proto

import * as jspb from "google-protobuf";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_extensions_graphql_stitching_pb from "../../../../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/extensions/graphql/stitching_pb";

export class GraphQLToolsStitchingInput extends jspb.Message {
  clearSubschemasList(): void;
  getSubschemasList(): Array<GraphQLToolsStitchingInput.Schema>;
  setSubschemasList(value: Array<GraphQLToolsStitchingInput.Schema>): void;
  addSubschemas(value?: GraphQLToolsStitchingInput.Schema, index?: number): GraphQLToolsStitchingInput.Schema;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GraphQLToolsStitchingInput.AsObject;
  static toObject(includeInstance: boolean, msg: GraphQLToolsStitchingInput): GraphQLToolsStitchingInput.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GraphQLToolsStitchingInput, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GraphQLToolsStitchingInput;
  static deserializeBinaryFromReader(message: GraphQLToolsStitchingInput, reader: jspb.BinaryReader): GraphQLToolsStitchingInput;
}

export namespace GraphQLToolsStitchingInput {
  export type AsObject = {
    subschemasList: Array<GraphQLToolsStitchingInput.Schema.AsObject>,
  }

  export class Schema extends jspb.Message {
    getName(): string;
    setName(value: string): void;

    getSchema(): string;
    setSchema(value: string): void;

    getTypeMergeConfigMap(): jspb.Map<string, GraphQLToolsStitchingInput.Schema.TypeMergeConfig>;
    clearTypeMergeConfigMap(): void;
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Schema.AsObject;
    static toObject(includeInstance: boolean, msg: Schema): Schema.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Schema, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Schema;
    static deserializeBinaryFromReader(message: Schema, reader: jspb.BinaryReader): Schema;
  }

  export namespace Schema {
    export type AsObject = {
      name: string,
      schema: string,
      typeMergeConfigMap: Array<[string, GraphQLToolsStitchingInput.Schema.TypeMergeConfig.AsObject]>,
    }

    export class TypeMergeConfig extends jspb.Message {
      getSelectionSet(): string;
      setSelectionSet(value: string): void;

      getFieldName(): string;
      setFieldName(value: string): void;

      getArgsMap(): jspb.Map<string, string>;
      clearArgsMap(): void;
      serializeBinary(): Uint8Array;
      toObject(includeInstance?: boolean): TypeMergeConfig.AsObject;
      static toObject(includeInstance: boolean, msg: TypeMergeConfig): TypeMergeConfig.AsObject;
      static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
      static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
      static serializeBinaryToWriter(message: TypeMergeConfig, writer: jspb.BinaryWriter): void;
      static deserializeBinary(bytes: Uint8Array): TypeMergeConfig;
      static deserializeBinaryFromReader(message: TypeMergeConfig, reader: jspb.BinaryReader): TypeMergeConfig;
    }

    export namespace TypeMergeConfig {
      export type AsObject = {
        selectionSet: string,
        fieldName: string,
        argsMap: Array<[string, string]>,
      }
    }
  }
}

export class GraphQlToolsStitchingOutput extends jspb.Message {
  getFieldNodesByTypeMap(): jspb.Map<string, github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_extensions_graphql_stitching_pb.FieldNodes>;
  clearFieldNodesByTypeMap(): void;
  getFieldNodesByFieldMap(): jspb.Map<string, github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_extensions_graphql_stitching_pb.FieldNodeMap>;
  clearFieldNodesByFieldMap(): void;
  getMergedTypesMap(): jspb.Map<string, github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_extensions_graphql_stitching_pb.MergedTypeConfig>;
  clearMergedTypesMap(): void;
  getStitchedSchema(): string;
  setStitchedSchema(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GraphQlToolsStitchingOutput.AsObject;
  static toObject(includeInstance: boolean, msg: GraphQlToolsStitchingOutput): GraphQlToolsStitchingOutput.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GraphQlToolsStitchingOutput, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GraphQlToolsStitchingOutput;
  static deserializeBinaryFromReader(message: GraphQlToolsStitchingOutput, reader: jspb.BinaryReader): GraphQlToolsStitchingOutput;
}

export namespace GraphQlToolsStitchingOutput {
  export type AsObject = {
    fieldNodesByTypeMap: Array<[string, github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_extensions_graphql_stitching_pb.FieldNodes.AsObject]>,
    fieldNodesByFieldMap: Array<[string, github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_extensions_graphql_stitching_pb.FieldNodeMap.AsObject]>,
    mergedTypesMap: Array<[string, github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_extensions_graphql_stitching_pb.MergedTypeConfig.AsObject]>,
    stitchedSchema: string,
  }
}
