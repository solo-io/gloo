// package: rest.plugins.gloo.solo.io
// file: github.com/solo-io/gloo/projects/gloo/api/v1/plugins/rest/rest.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../../../../../../../gogoproto/gogo_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_transformation_pb from "../../../../../../../../../github.com/solo-io/gloo/projects/gloo/api/external/envoy/extensions/transformation/transformation_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_v1_plugins_transformation_parameters_pb from "../../../../../../../../../github.com/solo-io/gloo/projects/gloo/api/v1/plugins/transformation/parameters_pb";

export class ServiceSpec extends jspb.Message {
  getTransformationsMap(): jspb.Map<string, github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_transformation_pb.TransformationTemplate>;
  clearTransformationsMap(): void;
  hasSwaggerInfo(): boolean;
  clearSwaggerInfo(): void;
  getSwaggerInfo(): ServiceSpec.SwaggerInfo | undefined;
  setSwaggerInfo(value?: ServiceSpec.SwaggerInfo): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ServiceSpec.AsObject;
  static toObject(includeInstance: boolean, msg: ServiceSpec): ServiceSpec.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ServiceSpec, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ServiceSpec;
  static deserializeBinaryFromReader(message: ServiceSpec, reader: jspb.BinaryReader): ServiceSpec;
}

export namespace ServiceSpec {
  export type AsObject = {
    transformationsMap: Array<[string, github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_transformation_pb.TransformationTemplate.AsObject]>,
    swaggerInfo?: ServiceSpec.SwaggerInfo.AsObject,
  }

  export class SwaggerInfo extends jspb.Message {
    hasUrl(): boolean;
    clearUrl(): void;
    getUrl(): string;
    setUrl(value: string): void;

    hasInline(): boolean;
    clearInline(): void;
    getInline(): string;
    setInline(value: string): void;

    getSwaggerSpecCase(): SwaggerInfo.SwaggerSpecCase;
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): SwaggerInfo.AsObject;
    static toObject(includeInstance: boolean, msg: SwaggerInfo): SwaggerInfo.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: SwaggerInfo, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): SwaggerInfo;
    static deserializeBinaryFromReader(message: SwaggerInfo, reader: jspb.BinaryReader): SwaggerInfo;
  }

  export namespace SwaggerInfo {
    export type AsObject = {
      url: string,
      inline: string,
    }

    export enum SwaggerSpecCase {
      SWAGGER_SPEC_NOT_SET = 0,
      URL = 1,
      INLINE = 2,
    }
  }
}

export class DestinationSpec extends jspb.Message {
  getFunctionName(): string;
  setFunctionName(value: string): void;

  hasParameters(): boolean;
  clearParameters(): void;
  getParameters(): github_com_solo_io_gloo_projects_gloo_api_v1_plugins_transformation_parameters_pb.Parameters | undefined;
  setParameters(value?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_transformation_parameters_pb.Parameters): void;

  hasResponseTransformation(): boolean;
  clearResponseTransformation(): void;
  getResponseTransformation(): github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_transformation_pb.TransformationTemplate | undefined;
  setResponseTransformation(value?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_transformation_pb.TransformationTemplate): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DestinationSpec.AsObject;
  static toObject(includeInstance: boolean, msg: DestinationSpec): DestinationSpec.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DestinationSpec, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DestinationSpec;
  static deserializeBinaryFromReader(message: DestinationSpec, reader: jspb.BinaryReader): DestinationSpec;
}

export namespace DestinationSpec {
  export type AsObject = {
    functionName: string,
    parameters?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_transformation_parameters_pb.Parameters.AsObject,
    responseTransformation?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_transformation_pb.TransformationTemplate.AsObject,
  }
}

