// package: glooeeapi.solo.io
// file: github.com/solo-io/solo-projects/projects/grpcserver/api/v1/envoy.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../../../../../gogoproto/gogo_pb";
import * as github_com_solo_io_solo_projects_projects_grpcserver_api_v1_types_pb from "../../../../../../../github.com/solo-io/solo-projects/projects/grpcserver/api/v1/types_pb";

export class EnvoyDetails extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  hasRaw(): boolean;
  clearRaw(): void;
  getRaw(): github_com_solo_io_solo_projects_projects_grpcserver_api_v1_types_pb.Raw | undefined;
  setRaw(value?: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_types_pb.Raw): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EnvoyDetails.AsObject;
  static toObject(includeInstance: boolean, msg: EnvoyDetails): EnvoyDetails.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: EnvoyDetails, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EnvoyDetails;
  static deserializeBinaryFromReader(message: EnvoyDetails, reader: jspb.BinaryReader): EnvoyDetails;
}

export namespace EnvoyDetails {
  export type AsObject = {
    name: string,
    raw?: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_types_pb.Raw.AsObject,
  }
}

export class ListEnvoyDetailsRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListEnvoyDetailsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListEnvoyDetailsRequest): ListEnvoyDetailsRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListEnvoyDetailsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListEnvoyDetailsRequest;
  static deserializeBinaryFromReader(message: ListEnvoyDetailsRequest, reader: jspb.BinaryReader): ListEnvoyDetailsRequest;
}

export namespace ListEnvoyDetailsRequest {
  export type AsObject = {
  }
}

export class ListEnvoyDetailsResponse extends jspb.Message {
  clearEnvoyDetailsList(): void;
  getEnvoyDetailsList(): Array<EnvoyDetails>;
  setEnvoyDetailsList(value: Array<EnvoyDetails>): void;
  addEnvoyDetails(value?: EnvoyDetails, index?: number): EnvoyDetails;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListEnvoyDetailsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListEnvoyDetailsResponse): ListEnvoyDetailsResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListEnvoyDetailsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListEnvoyDetailsResponse;
  static deserializeBinaryFromReader(message: ListEnvoyDetailsResponse, reader: jspb.BinaryReader): ListEnvoyDetailsResponse;
}

export namespace ListEnvoyDetailsResponse {
  export type AsObject = {
    envoyDetailsList: Array<EnvoyDetails.AsObject>,
  }
}

