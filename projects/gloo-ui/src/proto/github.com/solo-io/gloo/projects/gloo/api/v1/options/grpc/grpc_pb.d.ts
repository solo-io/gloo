// package: grpc.options.gloo.solo.io
// file: github.com/solo-io/gloo/projects/gloo/api/v1/options/grpc/grpc.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../../../../../../../gogoproto/gogo_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_v1_options_transformation_parameters_pb from "../../../../../../../../../github.com/solo-io/gloo/projects/gloo/api/v1/options/transformation/parameters_pb";

export class ServiceSpec extends jspb.Message {
  getDescriptors(): Uint8Array | string;
  getDescriptors_asU8(): Uint8Array;
  getDescriptors_asB64(): string;
  setDescriptors(value: Uint8Array | string): void;

  clearGrpcServicesList(): void;
  getGrpcServicesList(): Array<ServiceSpec.GrpcService>;
  setGrpcServicesList(value: Array<ServiceSpec.GrpcService>): void;
  addGrpcServices(value?: ServiceSpec.GrpcService, index?: number): ServiceSpec.GrpcService;

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
    descriptors: Uint8Array | string,
    grpcServicesList: Array<ServiceSpec.GrpcService.AsObject>,
  }

  export class GrpcService extends jspb.Message {
    getPackageName(): string;
    setPackageName(value: string): void;

    getServiceName(): string;
    setServiceName(value: string): void;

    clearFunctionNamesList(): void;
    getFunctionNamesList(): Array<string>;
    setFunctionNamesList(value: Array<string>): void;
    addFunctionNames(value: string, index?: number): string;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): GrpcService.AsObject;
    static toObject(includeInstance: boolean, msg: GrpcService): GrpcService.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: GrpcService, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): GrpcService;
    static deserializeBinaryFromReader(message: GrpcService, reader: jspb.BinaryReader): GrpcService;
  }

  export namespace GrpcService {
    export type AsObject = {
      packageName: string,
      serviceName: string,
      functionNamesList: Array<string>,
    }
  }
}

export class DestinationSpec extends jspb.Message {
  getPackage(): string;
  setPackage(value: string): void;

  getService(): string;
  setService(value: string): void;

  getFunction(): string;
  setFunction(value: string): void;

  hasParameters(): boolean;
  clearParameters(): void;
  getParameters(): github_com_solo_io_gloo_projects_gloo_api_v1_options_transformation_parameters_pb.Parameters | undefined;
  setParameters(value?: github_com_solo_io_gloo_projects_gloo_api_v1_options_transformation_parameters_pb.Parameters): void;

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
    pb_package: string,
    service: string,
    pb_function: string,
    parameters?: github_com_solo_io_gloo_projects_gloo_api_v1_options_transformation_parameters_pb.Parameters.AsObject,
  }
}

