// package: glooeeapi.solo.io
// file: solo-projects/projects/grpcserver/api/v1/config.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../../../gogoproto/gogo_pb";
import * as gloo_projects_gloo_api_v1_settings_pb from "../../../../../gloo/projects/gloo/api/v1/settings_pb";
import * as solo_projects_projects_grpcserver_api_v1_types_pb from "../../../../../solo-projects/projects/grpcserver/api/v1/types_pb";

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

export class SettingsDetails extends jspb.Message {
  hasSettings(): boolean;
  clearSettings(): void;
  getSettings(): gloo_projects_gloo_api_v1_settings_pb.Settings | undefined;
  setSettings(value?: gloo_projects_gloo_api_v1_settings_pb.Settings): void;

  hasRaw(): boolean;
  clearRaw(): void;
  getRaw(): solo_projects_projects_grpcserver_api_v1_types_pb.Raw | undefined;
  setRaw(value?: solo_projects_projects_grpcserver_api_v1_types_pb.Raw): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SettingsDetails.AsObject;
  static toObject(includeInstance: boolean, msg: SettingsDetails): SettingsDetails.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: SettingsDetails, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SettingsDetails;
  static deserializeBinaryFromReader(message: SettingsDetails, reader: jspb.BinaryReader): SettingsDetails;
}

export namespace SettingsDetails {
  export type AsObject = {
    settings?: gloo_projects_gloo_api_v1_settings_pb.Settings.AsObject,
    raw?: solo_projects_projects_grpcserver_api_v1_types_pb.Raw.AsObject,
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
  hasSettingsDetails(): boolean;
  clearSettingsDetails(): void;
  getSettingsDetails(): SettingsDetails | undefined;
  setSettingsDetails(value?: SettingsDetails): void;

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
    settingsDetails?: SettingsDetails.AsObject,
  }
}

export class UpdateSettingsRequest extends jspb.Message {
  hasSettings(): boolean;
  clearSettings(): void;
  getSettings(): gloo_projects_gloo_api_v1_settings_pb.Settings | undefined;
  setSettings(value?: gloo_projects_gloo_api_v1_settings_pb.Settings): void;

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
    settings?: gloo_projects_gloo_api_v1_settings_pb.Settings.AsObject,
  }
}

export class UpdateSettingsResponse extends jspb.Message {
  hasSettingsDetails(): boolean;
  clearSettingsDetails(): void;
  getSettingsDetails(): SettingsDetails | undefined;
  setSettingsDetails(value?: SettingsDetails): void;

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
    settingsDetails?: SettingsDetails.AsObject,
  }
}

export class UpdateSettingsYamlRequest extends jspb.Message {
  hasEditedYamlData(): boolean;
  clearEditedYamlData(): void;
  getEditedYamlData(): solo_projects_projects_grpcserver_api_v1_types_pb.EditedResourceYaml | undefined;
  setEditedYamlData(value?: solo_projects_projects_grpcserver_api_v1_types_pb.EditedResourceYaml): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateSettingsYamlRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateSettingsYamlRequest): UpdateSettingsYamlRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateSettingsYamlRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateSettingsYamlRequest;
  static deserializeBinaryFromReader(message: UpdateSettingsYamlRequest, reader: jspb.BinaryReader): UpdateSettingsYamlRequest;
}

export namespace UpdateSettingsYamlRequest {
  export type AsObject = {
    editedYamlData?: solo_projects_projects_grpcserver_api_v1_types_pb.EditedResourceYaml.AsObject,
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

  getInvalidReason(): string;
  setInvalidReason(value: string): void;

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
    invalidReason: string,
  }
}

export class GetPodNamespaceRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetPodNamespaceRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetPodNamespaceRequest): GetPodNamespaceRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetPodNamespaceRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetPodNamespaceRequest;
  static deserializeBinaryFromReader(message: GetPodNamespaceRequest, reader: jspb.BinaryReader): GetPodNamespaceRequest;
}

export namespace GetPodNamespaceRequest {
  export type AsObject = {
  }
}

export class GetPodNamespaceResponse extends jspb.Message {
  getNamespace(): string;
  setNamespace(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetPodNamespaceResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetPodNamespaceResponse): GetPodNamespaceResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetPodNamespaceResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetPodNamespaceResponse;
  static deserializeBinaryFromReader(message: GetPodNamespaceResponse, reader: jspb.BinaryReader): GetPodNamespaceResponse;
}

export namespace GetPodNamespaceResponse {
  export type AsObject = {
    namespace: string,
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

