/* eslint-disable */
// package: solo.io.envoy.type.tracing.v3
// file: github.com/solo-io/gloo/projects/gloo/api/external/envoy/type/tracing/v3/custom_tag.proto

import * as jspb from "google-protobuf";
import * as github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_metadata_v3_metadata_pb from "../../../../../../../../../../../github.com/solo-io/gloo/projects/gloo/api/external/envoy/type/metadata/v3/metadata_pb";
import * as udpa_annotations_status_pb from "../../../../../../../../../../../udpa/annotations/status_pb";
import * as udpa_annotations_versioning_pb from "../../../../../../../../../../../udpa/annotations/versioning_pb";
import * as validate_validate_pb from "../../../../../../../../../../../validate/validate_pb";
import * as extproto_ext_pb from "../../../../../../../../../../../protoc-gen-ext/extproto/ext_pb";

export class CustomTag extends jspb.Message {
  getTag(): string;
  setTag(value: string): void;

  hasLiteral(): boolean;
  clearLiteral(): void;
  getLiteral(): CustomTag.Literal | undefined;
  setLiteral(value?: CustomTag.Literal): void;

  hasEnvironment(): boolean;
  clearEnvironment(): void;
  getEnvironment(): CustomTag.Environment | undefined;
  setEnvironment(value?: CustomTag.Environment): void;

  hasRequestHeader(): boolean;
  clearRequestHeader(): void;
  getRequestHeader(): CustomTag.Header | undefined;
  setRequestHeader(value?: CustomTag.Header): void;

  hasMetadata(): boolean;
  clearMetadata(): void;
  getMetadata(): CustomTag.Metadata | undefined;
  setMetadata(value?: CustomTag.Metadata): void;

  getTypeCase(): CustomTag.TypeCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CustomTag.AsObject;
  static toObject(includeInstance: boolean, msg: CustomTag): CustomTag.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CustomTag, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CustomTag;
  static deserializeBinaryFromReader(message: CustomTag, reader: jspb.BinaryReader): CustomTag;
}

export namespace CustomTag {
  export type AsObject = {
    tag: string,
    literal?: CustomTag.Literal.AsObject,
    environment?: CustomTag.Environment.AsObject,
    requestHeader?: CustomTag.Header.AsObject,
    metadata?: CustomTag.Metadata.AsObject,
  }

  export class Literal extends jspb.Message {
    getValue(): string;
    setValue(value: string): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Literal.AsObject;
    static toObject(includeInstance: boolean, msg: Literal): Literal.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Literal, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Literal;
    static deserializeBinaryFromReader(message: Literal, reader: jspb.BinaryReader): Literal;
  }

  export namespace Literal {
    export type AsObject = {
      value: string,
    }
  }

  export class Environment extends jspb.Message {
    getName(): string;
    setName(value: string): void;

    getDefaultValue(): string;
    setDefaultValue(value: string): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Environment.AsObject;
    static toObject(includeInstance: boolean, msg: Environment): Environment.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Environment, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Environment;
    static deserializeBinaryFromReader(message: Environment, reader: jspb.BinaryReader): Environment;
  }

  export namespace Environment {
    export type AsObject = {
      name: string,
      defaultValue: string,
    }
  }

  export class Header extends jspb.Message {
    getName(): string;
    setName(value: string): void;

    getDefaultValue(): string;
    setDefaultValue(value: string): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Header.AsObject;
    static toObject(includeInstance: boolean, msg: Header): Header.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Header, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Header;
    static deserializeBinaryFromReader(message: Header, reader: jspb.BinaryReader): Header;
  }

  export namespace Header {
    export type AsObject = {
      name: string,
      defaultValue: string,
    }
  }

  export class Metadata extends jspb.Message {
    hasKind(): boolean;
    clearKind(): void;
    getKind(): github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_metadata_v3_metadata_pb.MetadataKind | undefined;
    setKind(value?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_metadata_v3_metadata_pb.MetadataKind): void;

    hasMetadataKey(): boolean;
    clearMetadataKey(): void;
    getMetadataKey(): github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_metadata_v3_metadata_pb.MetadataKey | undefined;
    setMetadataKey(value?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_metadata_v3_metadata_pb.MetadataKey): void;

    getDefaultValue(): string;
    setDefaultValue(value: string): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Metadata.AsObject;
    static toObject(includeInstance: boolean, msg: Metadata): Metadata.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Metadata, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Metadata;
    static deserializeBinaryFromReader(message: Metadata, reader: jspb.BinaryReader): Metadata;
  }

  export namespace Metadata {
    export type AsObject = {
      kind?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_metadata_v3_metadata_pb.MetadataKind.AsObject,
      metadataKey?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_metadata_v3_metadata_pb.MetadataKey.AsObject,
      defaultValue: string,
    }
  }

  export enum TypeCase {
    TYPE_NOT_SET = 0,
    LITERAL = 2,
    ENVIRONMENT = 3,
    REQUEST_HEADER = 4,
    METADATA = 5,
  }
}
