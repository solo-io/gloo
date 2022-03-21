/* eslint-disable */
// package: envoy.config.resolver.stitching.v2
// file: github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/extensions/graphql/stitching.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_any_pb from "google-protobuf/google/protobuf/any_pb";
import * as validate_validate_pb from "../../../../../../../../../../validate/validate_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_extensions_graphql_graphql_pb from "../../../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/extensions/graphql/graphql_pb";

export class FieldNode extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FieldNode.AsObject;
  static toObject(includeInstance: boolean, msg: FieldNode): FieldNode.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: FieldNode, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FieldNode;
  static deserializeBinaryFromReader(message: FieldNode, reader: jspb.BinaryReader): FieldNode;
}

export namespace FieldNode {
  export type AsObject = {
    name: string,
  }
}

export class FieldNodeMap extends jspb.Message {
  getNodesMap(): jspb.Map<string, FieldNodes>;
  clearNodesMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FieldNodeMap.AsObject;
  static toObject(includeInstance: boolean, msg: FieldNodeMap): FieldNodeMap.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: FieldNodeMap, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FieldNodeMap;
  static deserializeBinaryFromReader(message: FieldNodeMap, reader: jspb.BinaryReader): FieldNodeMap;
}

export namespace FieldNodeMap {
  export type AsObject = {
    nodesMap: Array<[string, FieldNodes.AsObject]>,
  }
}

export class FieldNodes extends jspb.Message {
  clearFieldNodesList(): void;
  getFieldNodesList(): Array<FieldNode>;
  setFieldNodesList(value: Array<FieldNode>): void;
  addFieldNodes(value?: FieldNode, index?: number): FieldNode;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FieldNodes.AsObject;
  static toObject(includeInstance: boolean, msg: FieldNodes): FieldNodes.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: FieldNodes, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FieldNodes;
  static deserializeBinaryFromReader(message: FieldNodes, reader: jspb.BinaryReader): FieldNodes;
}

export namespace FieldNodes {
  export type AsObject = {
    fieldNodesList: Array<FieldNode.AsObject>,
  }
}

export class ResolverConfig extends jspb.Message {
  getSelectionSet(): string;
  setSelectionSet(value: string): void;

  getFieldName(): string;
  setFieldName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResolverConfig.AsObject;
  static toObject(includeInstance: boolean, msg: ResolverConfig): ResolverConfig.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ResolverConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResolverConfig;
  static deserializeBinaryFromReader(message: ResolverConfig, reader: jspb.BinaryReader): ResolverConfig;
}

export namespace ResolverConfig {
  export type AsObject = {
    selectionSet: string,
    fieldName: string,
  }
}

export class Schemas extends jspb.Message {
  clearSchemasList(): void;
  getSchemasList(): Array<string>;
  setSchemasList(value: Array<string>): void;
  addSchemas(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Schemas.AsObject;
  static toObject(includeInstance: boolean, msg: Schemas): Schemas.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Schemas, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Schemas;
  static deserializeBinaryFromReader(message: Schemas, reader: jspb.BinaryReader): Schemas;
}

export namespace Schemas {
  export type AsObject = {
    schemasList: Array<string>,
  }
}

export class ArgPath extends jspb.Message {
  clearSetterPathList(): void;
  getSetterPathList(): Array<string>;
  setSetterPathList(value: Array<string>): void;
  addSetterPath(value: string, index?: number): string;

  clearExtractionPathList(): void;
  getExtractionPathList(): Array<string>;
  setExtractionPathList(value: Array<string>): void;
  addExtractionPath(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ArgPath.AsObject;
  static toObject(includeInstance: boolean, msg: ArgPath): ArgPath.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ArgPath, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ArgPath;
  static deserializeBinaryFromReader(message: ArgPath, reader: jspb.BinaryReader): ArgPath;
}

export namespace ArgPath {
  export type AsObject = {
    setterPathList: Array<string>,
    extractionPathList: Array<string>,
  }
}

export class ResolverInfo extends jspb.Message {
  getFieldName(): string;
  setFieldName(value: string): void;

  clearArgsList(): void;
  getArgsList(): Array<ArgPath>;
  setArgsList(value: Array<ArgPath>): void;
  addArgs(value?: ArgPath, index?: number): ArgPath;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResolverInfo.AsObject;
  static toObject(includeInstance: boolean, msg: ResolverInfo): ResolverInfo.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ResolverInfo, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResolverInfo;
  static deserializeBinaryFromReader(message: ResolverInfo, reader: jspb.BinaryReader): ResolverInfo;
}

export namespace ResolverInfo {
  export type AsObject = {
    fieldName: string,
    argsList: Array<ArgPath.AsObject>,
  }
}

export class MergedTypeConfig extends jspb.Message {
  getTypeName(): string;
  setTypeName(value: string): void;

  getSelectionSetsMap(): jspb.Map<string, string>;
  clearSelectionSetsMap(): void;
  getUniqueFieldsToSubschemaNameMap(): jspb.Map<string, string>;
  clearUniqueFieldsToSubschemaNameMap(): void;
  getNonUniqueFieldsToSubschemaNamesMap(): jspb.Map<string, Schemas>;
  clearNonUniqueFieldsToSubschemaNamesMap(): void;
  getDeclarativeTargetSubschemasMap(): jspb.Map<string, Schemas>;
  clearDeclarativeTargetSubschemasMap(): void;
  getSubschemaNameToResolverInfoMap(): jspb.Map<string, ResolverInfo>;
  clearSubschemaNameToResolverInfoMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MergedTypeConfig.AsObject;
  static toObject(includeInstance: boolean, msg: MergedTypeConfig): MergedTypeConfig.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: MergedTypeConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MergedTypeConfig;
  static deserializeBinaryFromReader(message: MergedTypeConfig, reader: jspb.BinaryReader): MergedTypeConfig;
}

export namespace MergedTypeConfig {
  export type AsObject = {
    typeName: string,
    selectionSetsMap: Array<[string, string]>,
    uniqueFieldsToSubschemaNameMap: Array<[string, string]>,
    nonUniqueFieldsToSubschemaNamesMap: Array<[string, Schemas.AsObject]>,
    declarativeTargetSubschemasMap: Array<[string, Schemas.AsObject]>,
    subschemaNameToResolverInfoMap: Array<[string, ResolverInfo.AsObject]>,
  }
}

export class StitchingInfo extends jspb.Message {
  getFieldNodesByTypeMap(): jspb.Map<string, FieldNodes>;
  clearFieldNodesByTypeMap(): void;
  getFieldNodesByFieldMap(): jspb.Map<string, FieldNodeMap>;
  clearFieldNodesByFieldMap(): void;
  getMergedTypesMap(): jspb.Map<string, MergedTypeConfig>;
  clearMergedTypesMap(): void;
  getSubschemaNameToSubschemaConfigMap(): jspb.Map<string, StitchingInfo.SubschemaConfig>;
  clearSubschemaNameToSubschemaConfigMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StitchingInfo.AsObject;
  static toObject(includeInstance: boolean, msg: StitchingInfo): StitchingInfo.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: StitchingInfo, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StitchingInfo;
  static deserializeBinaryFromReader(message: StitchingInfo, reader: jspb.BinaryReader): StitchingInfo;
}

export namespace StitchingInfo {
  export type AsObject = {
    fieldNodesByTypeMap: Array<[string, FieldNodes.AsObject]>,
    fieldNodesByFieldMap: Array<[string, FieldNodeMap.AsObject]>,
    mergedTypesMap: Array<[string, MergedTypeConfig.AsObject]>,
    subschemaNameToSubschemaConfigMap: Array<[string, StitchingInfo.SubschemaConfig.AsObject]>,
  }

  export class SubschemaConfig extends jspb.Message {
    hasExecutableSchema(): boolean;
    clearExecutableSchema(): void;
    getExecutableSchema(): github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_extensions_graphql_graphql_pb.ExecutableSchema | undefined;
    setExecutableSchema(value?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_extensions_graphql_graphql_pb.ExecutableSchema): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): SubschemaConfig.AsObject;
    static toObject(includeInstance: boolean, msg: SubschemaConfig): SubschemaConfig.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: SubschemaConfig, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): SubschemaConfig;
    static deserializeBinaryFromReader(message: SubschemaConfig, reader: jspb.BinaryReader): SubschemaConfig;
  }

  export namespace SubschemaConfig {
    export type AsObject = {
      executableSchema?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_extensions_graphql_graphql_pb.ExecutableSchema.AsObject,
    }
  }
}

export class StitchingResolver extends jspb.Message {
  getSubschemaName(): string;
  setSubschemaName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StitchingResolver.AsObject;
  static toObject(includeInstance: boolean, msg: StitchingResolver): StitchingResolver.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: StitchingResolver, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StitchingResolver;
  static deserializeBinaryFromReader(message: StitchingResolver, reader: jspb.BinaryReader): StitchingResolver;
}

export namespace StitchingResolver {
  export type AsObject = {
    subschemaName: string,
  }
}
