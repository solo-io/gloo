// package: glooeeapi.solo.io
// file: github.com/solo-io/solo-projects/projects/grpcserver/api/v1/virtualservice.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../../../../../gogoproto/gogo_pb";
import * as github_com_solo_io_gloo_projects_gateway_api_v1_virtual_service_pb from "../../../../../../../github.com/solo-io/gloo/projects/gateway/api/v1/virtual_service_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_v1_proxy_pb from "../../../../../../../github.com/solo-io/gloo/projects/gloo/api/v1/proxy_pb";
import * as github_com_solo_io_solo_projects_projects_gloo_api_v1_plugins_extauth_extauth_pb from "../../../../../../../github.com/solo-io/solo-projects/projects/gloo/api/v1/plugins/extauth/extauth_pb";
import * as github_com_solo_io_solo_projects_projects_gloo_api_v1_plugins_ratelimit_ratelimit_pb from "../../../../../../../github.com/solo-io/solo-projects/projects/gloo/api/v1/plugins/ratelimit/ratelimit_pb";
import * as github_com_solo_io_solo_kit_api_v1_ref_pb from "../../../../../../../github.com/solo-io/solo-kit/api/v1/ref_pb";

export class GetVirtualServiceRequest extends jspb.Message {
  hasRef(): boolean;
  clearRef(): void;
  getRef(): github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef | undefined;
  setRef(value?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetVirtualServiceRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetVirtualServiceRequest): GetVirtualServiceRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetVirtualServiceRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetVirtualServiceRequest;
  static deserializeBinaryFromReader(message: GetVirtualServiceRequest, reader: jspb.BinaryReader): GetVirtualServiceRequest;
}

export namespace GetVirtualServiceRequest {
  export type AsObject = {
    ref?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject,
  }
}

export class GetVirtualServiceResponse extends jspb.Message {
  hasVirtualService(): boolean;
  clearVirtualService(): void;
  getVirtualService(): github_com_solo_io_gloo_projects_gateway_api_v1_virtual_service_pb.VirtualService | undefined;
  setVirtualService(value?: github_com_solo_io_gloo_projects_gateway_api_v1_virtual_service_pb.VirtualService): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetVirtualServiceResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetVirtualServiceResponse): GetVirtualServiceResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetVirtualServiceResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetVirtualServiceResponse;
  static deserializeBinaryFromReader(message: GetVirtualServiceResponse, reader: jspb.BinaryReader): GetVirtualServiceResponse;
}

export namespace GetVirtualServiceResponse {
  export type AsObject = {
    virtualService?: github_com_solo_io_gloo_projects_gateway_api_v1_virtual_service_pb.VirtualService.AsObject,
  }
}

export class ListVirtualServicesRequest extends jspb.Message {
  clearNamespaceListList(): void;
  getNamespaceListList(): Array<string>;
  setNamespaceListList(value: Array<string>): void;
  addNamespaceList(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListVirtualServicesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListVirtualServicesRequest): ListVirtualServicesRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListVirtualServicesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListVirtualServicesRequest;
  static deserializeBinaryFromReader(message: ListVirtualServicesRequest, reader: jspb.BinaryReader): ListVirtualServicesRequest;
}

export namespace ListVirtualServicesRequest {
  export type AsObject = {
    namespaceListList: Array<string>,
  }
}

export class ListVirtualServicesResponse extends jspb.Message {
  clearVirtualServiceList(): void;
  getVirtualServiceList(): Array<github_com_solo_io_gloo_projects_gateway_api_v1_virtual_service_pb.VirtualService>;
  setVirtualServiceList(value: Array<github_com_solo_io_gloo_projects_gateway_api_v1_virtual_service_pb.VirtualService>): void;
  addVirtualService(value?: github_com_solo_io_gloo_projects_gateway_api_v1_virtual_service_pb.VirtualService, index?: number): github_com_solo_io_gloo_projects_gateway_api_v1_virtual_service_pb.VirtualService;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListVirtualServicesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListVirtualServicesResponse): ListVirtualServicesResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListVirtualServicesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListVirtualServicesResponse;
  static deserializeBinaryFromReader(message: ListVirtualServicesResponse, reader: jspb.BinaryReader): ListVirtualServicesResponse;
}

export namespace ListVirtualServicesResponse {
  export type AsObject = {
    virtualServiceList: Array<github_com_solo_io_gloo_projects_gateway_api_v1_virtual_service_pb.VirtualService.AsObject>,
  }
}

export class StreamVirtualServiceListRequest extends jspb.Message {
  getNamespace(): string;
  setNamespace(value: string): void;

  getSelectorMap(): jspb.Map<string, string>;
  clearSelectorMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StreamVirtualServiceListRequest.AsObject;
  static toObject(includeInstance: boolean, msg: StreamVirtualServiceListRequest): StreamVirtualServiceListRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: StreamVirtualServiceListRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StreamVirtualServiceListRequest;
  static deserializeBinaryFromReader(message: StreamVirtualServiceListRequest, reader: jspb.BinaryReader): StreamVirtualServiceListRequest;
}

export namespace StreamVirtualServiceListRequest {
  export type AsObject = {
    namespace: string,
    selectorMap: Array<[string, string]>,
  }
}

export class StreamVirtualServiceListResponse extends jspb.Message {
  hasVirtualService(): boolean;
  clearVirtualService(): void;
  getVirtualService(): github_com_solo_io_gloo_projects_gateway_api_v1_virtual_service_pb.VirtualService | undefined;
  setVirtualService(value?: github_com_solo_io_gloo_projects_gateway_api_v1_virtual_service_pb.VirtualService): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StreamVirtualServiceListResponse.AsObject;
  static toObject(includeInstance: boolean, msg: StreamVirtualServiceListResponse): StreamVirtualServiceListResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: StreamVirtualServiceListResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StreamVirtualServiceListResponse;
  static deserializeBinaryFromReader(message: StreamVirtualServiceListResponse, reader: jspb.BinaryReader): StreamVirtualServiceListResponse;
}

export namespace StreamVirtualServiceListResponse {
  export type AsObject = {
    virtualService?: github_com_solo_io_gloo_projects_gateway_api_v1_virtual_service_pb.VirtualService.AsObject,
  }
}

export class VirtualServiceInput extends jspb.Message {
  hasRef(): boolean;
  clearRef(): void;
  getRef(): github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef | undefined;
  setRef(value?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef): void;

  getDisplayName(): string;
  setDisplayName(value: string): void;

  clearDomainsList(): void;
  getDomainsList(): Array<string>;
  setDomainsList(value: Array<string>): void;
  addDomains(value: string, index?: number): string;

  clearRoutesList(): void;
  getRoutesList(): Array<github_com_solo_io_gloo_projects_gloo_api_v1_proxy_pb.Route>;
  setRoutesList(value: Array<github_com_solo_io_gloo_projects_gloo_api_v1_proxy_pb.Route>): void;
  addRoutes(value?: github_com_solo_io_gloo_projects_gloo_api_v1_proxy_pb.Route, index?: number): github_com_solo_io_gloo_projects_gloo_api_v1_proxy_pb.Route;

  hasSecretRef(): boolean;
  clearSecretRef(): void;
  getSecretRef(): github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef | undefined;
  setSecretRef(value?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef): void;

  hasRateLimitConfig(): boolean;
  clearRateLimitConfig(): void;
  getRateLimitConfig(): github_com_solo_io_solo_projects_projects_gloo_api_v1_plugins_ratelimit_ratelimit_pb.IngressRateLimit | undefined;
  setRateLimitConfig(value?: github_com_solo_io_solo_projects_projects_gloo_api_v1_plugins_ratelimit_ratelimit_pb.IngressRateLimit): void;

  hasBasicAuth(): boolean;
  clearBasicAuth(): void;
  getBasicAuth(): VirtualServiceInput.BasicAuthInput | undefined;
  setBasicAuth(value?: VirtualServiceInput.BasicAuthInput): void;

  hasOauth(): boolean;
  clearOauth(): void;
  getOauth(): github_com_solo_io_solo_projects_projects_gloo_api_v1_plugins_extauth_extauth_pb.OAuth | undefined;
  setOauth(value?: github_com_solo_io_solo_projects_projects_gloo_api_v1_plugins_extauth_extauth_pb.OAuth): void;

  hasCustomAuth(): boolean;
  clearCustomAuth(): void;
  getCustomAuth(): github_com_solo_io_solo_projects_projects_gloo_api_v1_plugins_extauth_extauth_pb.CustomAuth | undefined;
  setCustomAuth(value?: github_com_solo_io_solo_projects_projects_gloo_api_v1_plugins_extauth_extauth_pb.CustomAuth): void;

  getExtAuthConfigCase(): VirtualServiceInput.ExtAuthConfigCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VirtualServiceInput.AsObject;
  static toObject(includeInstance: boolean, msg: VirtualServiceInput): VirtualServiceInput.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: VirtualServiceInput, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VirtualServiceInput;
  static deserializeBinaryFromReader(message: VirtualServiceInput, reader: jspb.BinaryReader): VirtualServiceInput;
}

export namespace VirtualServiceInput {
  export type AsObject = {
    ref?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject,
    displayName: string,
    domainsList: Array<string>,
    routesList: Array<github_com_solo_io_gloo_projects_gloo_api_v1_proxy_pb.Route.AsObject>,
    secretRef?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject,
    rateLimitConfig?: github_com_solo_io_solo_projects_projects_gloo_api_v1_plugins_ratelimit_ratelimit_pb.IngressRateLimit.AsObject,
    basicAuth?: VirtualServiceInput.BasicAuthInput.AsObject,
    oauth?: github_com_solo_io_solo_projects_projects_gloo_api_v1_plugins_extauth_extauth_pb.OAuth.AsObject,
    customAuth?: github_com_solo_io_solo_projects_projects_gloo_api_v1_plugins_extauth_extauth_pb.CustomAuth.AsObject,
  }

  export class BasicAuthInput extends jspb.Message {
    getRealm(): string;
    setRealm(value: string): void;

    getSpecCsv(): string;
    setSpecCsv(value: string): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): BasicAuthInput.AsObject;
    static toObject(includeInstance: boolean, msg: BasicAuthInput): BasicAuthInput.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: BasicAuthInput, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): BasicAuthInput;
    static deserializeBinaryFromReader(message: BasicAuthInput, reader: jspb.BinaryReader): BasicAuthInput;
  }

  export namespace BasicAuthInput {
    export type AsObject = {
      realm: string,
      specCsv: string,
    }
  }

  export enum ExtAuthConfigCase {
    EXT_AUTH_CONFIG_NOT_SET = 0,
    BASIC_AUTH = 7,
    OAUTH = 8,
    CUSTOM_AUTH = 9,
  }
}

export class CreateVirtualServiceRequest extends jspb.Message {
  hasInput(): boolean;
  clearInput(): void;
  getInput(): VirtualServiceInput | undefined;
  setInput(value?: VirtualServiceInput): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateVirtualServiceRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateVirtualServiceRequest): CreateVirtualServiceRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateVirtualServiceRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateVirtualServiceRequest;
  static deserializeBinaryFromReader(message: CreateVirtualServiceRequest, reader: jspb.BinaryReader): CreateVirtualServiceRequest;
}

export namespace CreateVirtualServiceRequest {
  export type AsObject = {
    input?: VirtualServiceInput.AsObject,
  }
}

export class CreateVirtualServiceResponse extends jspb.Message {
  hasVirtualService(): boolean;
  clearVirtualService(): void;
  getVirtualService(): github_com_solo_io_gloo_projects_gateway_api_v1_virtual_service_pb.VirtualService | undefined;
  setVirtualService(value?: github_com_solo_io_gloo_projects_gateway_api_v1_virtual_service_pb.VirtualService): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateVirtualServiceResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CreateVirtualServiceResponse): CreateVirtualServiceResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateVirtualServiceResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateVirtualServiceResponse;
  static deserializeBinaryFromReader(message: CreateVirtualServiceResponse, reader: jspb.BinaryReader): CreateVirtualServiceResponse;
}

export namespace CreateVirtualServiceResponse {
  export type AsObject = {
    virtualService?: github_com_solo_io_gloo_projects_gateway_api_v1_virtual_service_pb.VirtualService.AsObject,
  }
}

export class UpdateVirtualServiceRequest extends jspb.Message {
  hasInput(): boolean;
  clearInput(): void;
  getInput(): VirtualServiceInput | undefined;
  setInput(value?: VirtualServiceInput): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateVirtualServiceRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateVirtualServiceRequest): UpdateVirtualServiceRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateVirtualServiceRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateVirtualServiceRequest;
  static deserializeBinaryFromReader(message: UpdateVirtualServiceRequest, reader: jspb.BinaryReader): UpdateVirtualServiceRequest;
}

export namespace UpdateVirtualServiceRequest {
  export type AsObject = {
    input?: VirtualServiceInput.AsObject,
  }
}

export class UpdateVirtualServiceResponse extends jspb.Message {
  hasVirtualService(): boolean;
  clearVirtualService(): void;
  getVirtualService(): github_com_solo_io_gloo_projects_gateway_api_v1_virtual_service_pb.VirtualService | undefined;
  setVirtualService(value?: github_com_solo_io_gloo_projects_gateway_api_v1_virtual_service_pb.VirtualService): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateVirtualServiceResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateVirtualServiceResponse): UpdateVirtualServiceResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateVirtualServiceResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateVirtualServiceResponse;
  static deserializeBinaryFromReader(message: UpdateVirtualServiceResponse, reader: jspb.BinaryReader): UpdateVirtualServiceResponse;
}

export namespace UpdateVirtualServiceResponse {
  export type AsObject = {
    virtualService?: github_com_solo_io_gloo_projects_gateway_api_v1_virtual_service_pb.VirtualService.AsObject,
  }
}

export class DeleteVirtualServiceRequest extends jspb.Message {
  hasRef(): boolean;
  clearRef(): void;
  getRef(): github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef | undefined;
  setRef(value?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteVirtualServiceRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteVirtualServiceRequest): DeleteVirtualServiceRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteVirtualServiceRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteVirtualServiceRequest;
  static deserializeBinaryFromReader(message: DeleteVirtualServiceRequest, reader: jspb.BinaryReader): DeleteVirtualServiceRequest;
}

export namespace DeleteVirtualServiceRequest {
  export type AsObject = {
    ref?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject,
  }
}

export class DeleteVirtualServiceResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteVirtualServiceResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteVirtualServiceResponse): DeleteVirtualServiceResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteVirtualServiceResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteVirtualServiceResponse;
  static deserializeBinaryFromReader(message: DeleteVirtualServiceResponse, reader: jspb.BinaryReader): DeleteVirtualServiceResponse;
}

export namespace DeleteVirtualServiceResponse {
  export type AsObject = {
  }
}

export class RouteInput extends jspb.Message {
  hasVirtualServiceRef(): boolean;
  clearVirtualServiceRef(): void;
  getVirtualServiceRef(): github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef | undefined;
  setVirtualServiceRef(value?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef): void;

  getIndex(): number;
  setIndex(value: number): void;

  hasRoute(): boolean;
  clearRoute(): void;
  getRoute(): github_com_solo_io_gloo_projects_gloo_api_v1_proxy_pb.Route | undefined;
  setRoute(value?: github_com_solo_io_gloo_projects_gloo_api_v1_proxy_pb.Route): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RouteInput.AsObject;
  static toObject(includeInstance: boolean, msg: RouteInput): RouteInput.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RouteInput, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RouteInput;
  static deserializeBinaryFromReader(message: RouteInput, reader: jspb.BinaryReader): RouteInput;
}

export namespace RouteInput {
  export type AsObject = {
    virtualServiceRef?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject,
    index: number,
    route?: github_com_solo_io_gloo_projects_gloo_api_v1_proxy_pb.Route.AsObject,
  }
}

export class CreateRouteRequest extends jspb.Message {
  hasInput(): boolean;
  clearInput(): void;
  getInput(): RouteInput | undefined;
  setInput(value?: RouteInput): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateRouteRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateRouteRequest): CreateRouteRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateRouteRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateRouteRequest;
  static deserializeBinaryFromReader(message: CreateRouteRequest, reader: jspb.BinaryReader): CreateRouteRequest;
}

export namespace CreateRouteRequest {
  export type AsObject = {
    input?: RouteInput.AsObject,
  }
}

export class CreateRouteResponse extends jspb.Message {
  hasVirtualService(): boolean;
  clearVirtualService(): void;
  getVirtualService(): github_com_solo_io_gloo_projects_gateway_api_v1_virtual_service_pb.VirtualService | undefined;
  setVirtualService(value?: github_com_solo_io_gloo_projects_gateway_api_v1_virtual_service_pb.VirtualService): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateRouteResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CreateRouteResponse): CreateRouteResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateRouteResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateRouteResponse;
  static deserializeBinaryFromReader(message: CreateRouteResponse, reader: jspb.BinaryReader): CreateRouteResponse;
}

export namespace CreateRouteResponse {
  export type AsObject = {
    virtualService?: github_com_solo_io_gloo_projects_gateway_api_v1_virtual_service_pb.VirtualService.AsObject,
  }
}

export class UpdateRouteRequest extends jspb.Message {
  hasInput(): boolean;
  clearInput(): void;
  getInput(): RouteInput | undefined;
  setInput(value?: RouteInput): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateRouteRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateRouteRequest): UpdateRouteRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateRouteRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateRouteRequest;
  static deserializeBinaryFromReader(message: UpdateRouteRequest, reader: jspb.BinaryReader): UpdateRouteRequest;
}

export namespace UpdateRouteRequest {
  export type AsObject = {
    input?: RouteInput.AsObject,
  }
}

export class UpdateRouteResponse extends jspb.Message {
  hasVirtualService(): boolean;
  clearVirtualService(): void;
  getVirtualService(): github_com_solo_io_gloo_projects_gateway_api_v1_virtual_service_pb.VirtualService | undefined;
  setVirtualService(value?: github_com_solo_io_gloo_projects_gateway_api_v1_virtual_service_pb.VirtualService): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateRouteResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateRouteResponse): UpdateRouteResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateRouteResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateRouteResponse;
  static deserializeBinaryFromReader(message: UpdateRouteResponse, reader: jspb.BinaryReader): UpdateRouteResponse;
}

export namespace UpdateRouteResponse {
  export type AsObject = {
    virtualService?: github_com_solo_io_gloo_projects_gateway_api_v1_virtual_service_pb.VirtualService.AsObject,
  }
}

export class DeleteRouteRequest extends jspb.Message {
  hasVirtualServiceRef(): boolean;
  clearVirtualServiceRef(): void;
  getVirtualServiceRef(): github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef | undefined;
  setVirtualServiceRef(value?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef): void;

  getIndex(): number;
  setIndex(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteRouteRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteRouteRequest): DeleteRouteRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteRouteRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteRouteRequest;
  static deserializeBinaryFromReader(message: DeleteRouteRequest, reader: jspb.BinaryReader): DeleteRouteRequest;
}

export namespace DeleteRouteRequest {
  export type AsObject = {
    virtualServiceRef?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject,
    index: number,
  }
}

export class DeleteRouteResponse extends jspb.Message {
  hasVirtualService(): boolean;
  clearVirtualService(): void;
  getVirtualService(): github_com_solo_io_gloo_projects_gateway_api_v1_virtual_service_pb.VirtualService | undefined;
  setVirtualService(value?: github_com_solo_io_gloo_projects_gateway_api_v1_virtual_service_pb.VirtualService): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteRouteResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteRouteResponse): DeleteRouteResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteRouteResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteRouteResponse;
  static deserializeBinaryFromReader(message: DeleteRouteResponse, reader: jspb.BinaryReader): DeleteRouteResponse;
}

export namespace DeleteRouteResponse {
  export type AsObject = {
    virtualService?: github_com_solo_io_gloo_projects_gateway_api_v1_virtual_service_pb.VirtualService.AsObject,
  }
}

export class SwapRoutesRequest extends jspb.Message {
  hasVirtualServiceRef(): boolean;
  clearVirtualServiceRef(): void;
  getVirtualServiceRef(): github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef | undefined;
  setVirtualServiceRef(value?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef): void;

  getIndex1(): number;
  setIndex1(value: number): void;

  getIndex2(): number;
  setIndex2(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SwapRoutesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SwapRoutesRequest): SwapRoutesRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: SwapRoutesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SwapRoutesRequest;
  static deserializeBinaryFromReader(message: SwapRoutesRequest, reader: jspb.BinaryReader): SwapRoutesRequest;
}

export namespace SwapRoutesRequest {
  export type AsObject = {
    virtualServiceRef?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject,
    index1: number,
    index2: number,
  }
}

export class SwapRoutesResponse extends jspb.Message {
  hasVirtualService(): boolean;
  clearVirtualService(): void;
  getVirtualService(): github_com_solo_io_gloo_projects_gateway_api_v1_virtual_service_pb.VirtualService | undefined;
  setVirtualService(value?: github_com_solo_io_gloo_projects_gateway_api_v1_virtual_service_pb.VirtualService): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SwapRoutesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SwapRoutesResponse): SwapRoutesResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: SwapRoutesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SwapRoutesResponse;
  static deserializeBinaryFromReader(message: SwapRoutesResponse, reader: jspb.BinaryReader): SwapRoutesResponse;
}

export namespace SwapRoutesResponse {
  export type AsObject = {
    virtualService?: github_com_solo_io_gloo_projects_gateway_api_v1_virtual_service_pb.VirtualService.AsObject,
  }
}

export class ShiftRoutesRequest extends jspb.Message {
  hasVirtualServiceRef(): boolean;
  clearVirtualServiceRef(): void;
  getVirtualServiceRef(): github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef | undefined;
  setVirtualServiceRef(value?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef): void;

  getFromIndex(): number;
  setFromIndex(value: number): void;

  getToIndex(): number;
  setToIndex(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ShiftRoutesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ShiftRoutesRequest): ShiftRoutesRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ShiftRoutesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ShiftRoutesRequest;
  static deserializeBinaryFromReader(message: ShiftRoutesRequest, reader: jspb.BinaryReader): ShiftRoutesRequest;
}

export namespace ShiftRoutesRequest {
  export type AsObject = {
    virtualServiceRef?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject,
    fromIndex: number,
    toIndex: number,
  }
}

export class ShiftRoutesResponse extends jspb.Message {
  hasVirtualService(): boolean;
  clearVirtualService(): void;
  getVirtualService(): github_com_solo_io_gloo_projects_gateway_api_v1_virtual_service_pb.VirtualService | undefined;
  setVirtualService(value?: github_com_solo_io_gloo_projects_gateway_api_v1_virtual_service_pb.VirtualService): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ShiftRoutesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ShiftRoutesResponse): ShiftRoutesResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ShiftRoutesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ShiftRoutesResponse;
  static deserializeBinaryFromReader(message: ShiftRoutesResponse, reader: jspb.BinaryReader): ShiftRoutesResponse;
}

export namespace ShiftRoutesResponse {
  export type AsObject = {
    virtualService?: github_com_solo_io_gloo_projects_gateway_api_v1_virtual_service_pb.VirtualService.AsObject,
  }
}

