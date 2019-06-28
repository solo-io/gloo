// package: glooeeapi.solo.io
// file: github.com/solo-io/solo-projects/projects/grpcserver/api/v1/config.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../../../../../gogoproto/gogo_pb";
import * as google_protobuf_duration_pb from "google-protobuf/google/protobuf/duration_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_v1_settings_pb from "../../../../../../../github.com/solo-io/gloo/projects/gloo/api/v1/settings_pb";
import * as github_com_solo_io_solo_kit_api_v1_ref_pb from "../../../../../../../github.com/solo-io/solo-kit/api/v1/ref_pb";

export class OAuthEndpoint extends jspb.Message {
  getUrl(): string;
  setUrl(value: string): void;

  getClientName(): string;
  setClientName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OAuthEndpoint.AsObject;
  static toObject(includeInstance: boolean, msg: OAuthEndpoint): OAuthEndpoint.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: OAuthEndpoint, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OAuthEndpoint;
  static deserializeBinaryFromReader(message: OAuthEndpoint, reader: jspb.BinaryReader): OAuthEndpoint;
}

export namespace OAuthEndpoint {
  export type AsObject = {
    url: string,
    clientName: string,
  }
}

export class GetVersionRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetVersionRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetVersionRequest): GetVersionRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetVersionRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetVersionRequest;
  static deserializeBinaryFromReader(message: GetVersionRequest, reader: jspb.BinaryReader): GetVersionRequest;
}

export namespace GetVersionRequest {
  export type AsObject = {
  }
}

export class GetVersionResponse extends jspb.Message {
  getVersion(): string;
  setVersion(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetVersionResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetVersionResponse): GetVersionResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetVersionResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetVersionResponse;
  static deserializeBinaryFromReader(message: GetVersionResponse, reader: jspb.BinaryReader): GetVersionResponse;
}

export namespace GetVersionResponse {
  export type AsObject = {
    version: string,
  }
}

export class GetOAuthEndpointRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetOAuthEndpointRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetOAuthEndpointRequest): GetOAuthEndpointRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetOAuthEndpointRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetOAuthEndpointRequest;
  static deserializeBinaryFromReader(message: GetOAuthEndpointRequest, reader: jspb.BinaryReader): GetOAuthEndpointRequest;
}

export namespace GetOAuthEndpointRequest {
  export type AsObject = {
  }
}

export class GetOAuthEndpointResponse extends jspb.Message {
  hasOAuthEndpoint(): boolean;
  clearOAuthEndpoint(): void;
  getOAuthEndpoint(): OAuthEndpoint | undefined;
  setOAuthEndpoint(value?: OAuthEndpoint): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetOAuthEndpointResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetOAuthEndpointResponse): GetOAuthEndpointResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetOAuthEndpointResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetOAuthEndpointResponse;
  static deserializeBinaryFromReader(message: GetOAuthEndpointResponse, reader: jspb.BinaryReader): GetOAuthEndpointResponse;
}

export namespace GetOAuthEndpointResponse {
  export type AsObject = {
    oAuthEndpoint?: OAuthEndpoint.AsObject,
  }
}

export class GetSettingsRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetSettingsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetSettingsRequest): GetSettingsRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetSettingsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetSettingsRequest;
  static deserializeBinaryFromReader(message: GetSettingsRequest, reader: jspb.BinaryReader): GetSettingsRequest;
}

export namespace GetSettingsRequest {
  export type AsObject = {
  }
}

export class GetSettingsResponse extends jspb.Message {
  hasSettings(): boolean;
  clearSettings(): void;
  getSettings(): github_com_solo_io_gloo_projects_gloo_api_v1_settings_pb.Settings | undefined;
  setSettings(value?: github_com_solo_io_gloo_projects_gloo_api_v1_settings_pb.Settings): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetSettingsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetSettingsResponse): GetSettingsResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetSettingsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetSettingsResponse;
  static deserializeBinaryFromReader(message: GetSettingsResponse, reader: jspb.BinaryReader): GetSettingsResponse;
}

export namespace GetSettingsResponse {
  export type AsObject = {
    settings?: github_com_solo_io_gloo_projects_gloo_api_v1_settings_pb.Settings.AsObject,
  }
}

export class UpdateSettingsRequest extends jspb.Message {
  hasRef(): boolean;
  clearRef(): void;
  getRef(): github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef | undefined;
  setRef(value?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef): void;

  clearWatchNamespacesList(): void;
  getWatchNamespacesList(): Array<string>;
  setWatchNamespacesList(value: Array<string>): void;
  addWatchNamespaces(value: string, index?: number): string;

  hasRefreshRate(): boolean;
  clearRefreshRate(): void;
  getRefreshRate(): google_protobuf_duration_pb.Duration | undefined;
  setRefreshRate(value?: google_protobuf_duration_pb.Duration): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateSettingsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateSettingsRequest): UpdateSettingsRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateSettingsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateSettingsRequest;
  static deserializeBinaryFromReader(message: UpdateSettingsRequest, reader: jspb.BinaryReader): UpdateSettingsRequest;
}

export namespace UpdateSettingsRequest {
  export type AsObject = {
    ref?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject,
    watchNamespacesList: Array<string>,
    refreshRate?: google_protobuf_duration_pb.Duration.AsObject,
  }
}

export class UpdateSettingsResponse extends jspb.Message {
  hasSettings(): boolean;
  clearSettings(): void;
  getSettings(): github_com_solo_io_gloo_projects_gloo_api_v1_settings_pb.Settings | undefined;
  setSettings(value?: github_com_solo_io_gloo_projects_gloo_api_v1_settings_pb.Settings): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateSettingsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateSettingsResponse): UpdateSettingsResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateSettingsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateSettingsResponse;
  static deserializeBinaryFromReader(message: UpdateSettingsResponse, reader: jspb.BinaryReader): UpdateSettingsResponse;
}

export namespace UpdateSettingsResponse {
  export type AsObject = {
    settings?: github_com_solo_io_gloo_projects_gloo_api_v1_settings_pb.Settings.AsObject,
  }
}

export class GetIsLicenseValidRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetIsLicenseValidRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetIsLicenseValidRequest): GetIsLicenseValidRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetIsLicenseValidRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetIsLicenseValidRequest;
  static deserializeBinaryFromReader(message: GetIsLicenseValidRequest, reader: jspb.BinaryReader): GetIsLicenseValidRequest;
}

export namespace GetIsLicenseValidRequest {
  export type AsObject = {
  }
}

export class GetIsLicenseValidResponse extends jspb.Message {
  getIsLicenseValid(): boolean;
  setIsLicenseValid(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetIsLicenseValidResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetIsLicenseValidResponse): GetIsLicenseValidResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetIsLicenseValidResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetIsLicenseValidResponse;
  static deserializeBinaryFromReader(message: GetIsLicenseValidResponse, reader: jspb.BinaryReader): GetIsLicenseValidResponse;
}

export namespace GetIsLicenseValidResponse {
  export type AsObject = {
    isLicenseValid: boolean,
  }
}

export class ListNamespacesRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListNamespacesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListNamespacesRequest): ListNamespacesRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListNamespacesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListNamespacesRequest;
  static deserializeBinaryFromReader(message: ListNamespacesRequest, reader: jspb.BinaryReader): ListNamespacesRequest;
}

export namespace ListNamespacesRequest {
  export type AsObject = {
  }
}

export class ListNamespacesResponse extends jspb.Message {
  clearNamespacesList(): void;
  getNamespacesList(): Array<string>;
  setNamespacesList(value: Array<string>): void;
  addNamespaces(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListNamespacesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListNamespacesResponse): ListNamespacesResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListNamespacesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListNamespacesResponse;
  static deserializeBinaryFromReader(message: ListNamespacesResponse, reader: jspb.BinaryReader): ListNamespacesResponse;
}

export namespace ListNamespacesResponse {
  export type AsObject = {
    namespacesList: Array<string>,
  }
}

