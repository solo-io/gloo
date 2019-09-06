// package: gateway.solo.io
// file: github.com/solo-io/gloo/projects/gateway/api/v1/virtual_service.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../../../../../gogoproto/gogo_pb";
import * as github_com_solo_io_solo_kit_api_v1_metadata_pb from "../../../../../../../github.com/solo-io/solo-kit/api/v1/metadata_pb";
import * as github_com_solo_io_solo_kit_api_v1_status_pb from "../../../../../../../github.com/solo-io/solo-kit/api/v1/status_pb";
import * as github_com_solo_io_solo_kit_api_v1_ref_pb from "../../../../../../../github.com/solo-io/solo-kit/api/v1/ref_pb";
import * as github_com_solo_io_solo_kit_api_v1_solo_kit_pb from "../../../../../../../github.com/solo-io/solo-kit/api/v1/solo-kit_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_v1_ssl_pb from "../../../../../../../github.com/solo-io/gloo/projects/gloo/api/v1/ssl_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_v1_proxy_pb from "../../../../../../../github.com/solo-io/gloo/projects/gloo/api/v1/proxy_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_v1_plugins_pb from "../../../../../../../github.com/solo-io/gloo/projects/gloo/api/v1/plugins_pb";

export class VirtualService extends jspb.Message {
  hasVirtualHost(): boolean;
  clearVirtualHost(): void;
  getVirtualHost(): VirtualHost | undefined;
  setVirtualHost(value?: VirtualHost): void;

  hasSslConfig(): boolean;
  clearSslConfig(): void;
  getSslConfig(): github_com_solo_io_gloo_projects_gloo_api_v1_ssl_pb.SslConfig | undefined;
  setSslConfig(value?: github_com_solo_io_gloo_projects_gloo_api_v1_ssl_pb.SslConfig): void;

  getDisplayName(): string;
  setDisplayName(value: string): void;

  hasStatus(): boolean;
  clearStatus(): void;
  getStatus(): github_com_solo_io_solo_kit_api_v1_status_pb.Status | undefined;
  setStatus(value?: github_com_solo_io_solo_kit_api_v1_status_pb.Status): void;

  hasMetadata(): boolean;
  clearMetadata(): void;
  getMetadata(): github_com_solo_io_solo_kit_api_v1_metadata_pb.Metadata | undefined;
  setMetadata(value?: github_com_solo_io_solo_kit_api_v1_metadata_pb.Metadata): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VirtualService.AsObject;
  static toObject(includeInstance: boolean, msg: VirtualService): VirtualService.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: VirtualService, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VirtualService;
  static deserializeBinaryFromReader(message: VirtualService, reader: jspb.BinaryReader): VirtualService;
}

export namespace VirtualService {
  export type AsObject = {
    virtualHost?: VirtualHost.AsObject,
    sslConfig?: github_com_solo_io_gloo_projects_gloo_api_v1_ssl_pb.SslConfig.AsObject,
    displayName: string,
    status?: github_com_solo_io_solo_kit_api_v1_status_pb.Status.AsObject,
    metadata?: github_com_solo_io_solo_kit_api_v1_metadata_pb.Metadata.AsObject,
  }
}

export class VirtualHost extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  clearDomainsList(): void;
  getDomainsList(): Array<string>;
  setDomainsList(value: Array<string>): void;
  addDomains(value: string, index?: number): string;

  clearRoutesList(): void;
  getRoutesList(): Array<Route>;
  setRoutesList(value: Array<Route>): void;
  addRoutes(value?: Route, index?: number): Route;

  hasVirtualHostPlugins(): boolean;
  clearVirtualHostPlugins(): void;
  getVirtualHostPlugins(): github_com_solo_io_gloo_projects_gloo_api_v1_plugins_pb.VirtualHostPlugins | undefined;
  setVirtualHostPlugins(value?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_pb.VirtualHostPlugins): void;

  hasCorsPolicy(): boolean;
  clearCorsPolicy(): void;
  getCorsPolicy(): github_com_solo_io_gloo_projects_gloo_api_v1_proxy_pb.CorsPolicy | undefined;
  setCorsPolicy(value?: github_com_solo_io_gloo_projects_gloo_api_v1_proxy_pb.CorsPolicy): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VirtualHost.AsObject;
  static toObject(includeInstance: boolean, msg: VirtualHost): VirtualHost.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: VirtualHost, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VirtualHost;
  static deserializeBinaryFromReader(message: VirtualHost, reader: jspb.BinaryReader): VirtualHost;
}

export namespace VirtualHost {
  export type AsObject = {
    name: string,
    domainsList: Array<string>,
    routesList: Array<Route.AsObject>,
    virtualHostPlugins?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_pb.VirtualHostPlugins.AsObject,
    corsPolicy?: github_com_solo_io_gloo_projects_gloo_api_v1_proxy_pb.CorsPolicy.AsObject,
  }
}

export class Route extends jspb.Message {
  hasMatcher(): boolean;
  clearMatcher(): void;
  getMatcher(): github_com_solo_io_gloo_projects_gloo_api_v1_proxy_pb.Matcher | undefined;
  setMatcher(value?: github_com_solo_io_gloo_projects_gloo_api_v1_proxy_pb.Matcher): void;

  hasRouteAction(): boolean;
  clearRouteAction(): void;
  getRouteAction(): github_com_solo_io_gloo_projects_gloo_api_v1_proxy_pb.RouteAction | undefined;
  setRouteAction(value?: github_com_solo_io_gloo_projects_gloo_api_v1_proxy_pb.RouteAction): void;

  hasRedirectAction(): boolean;
  clearRedirectAction(): void;
  getRedirectAction(): github_com_solo_io_gloo_projects_gloo_api_v1_proxy_pb.RedirectAction | undefined;
  setRedirectAction(value?: github_com_solo_io_gloo_projects_gloo_api_v1_proxy_pb.RedirectAction): void;

  hasDirectResponseAction(): boolean;
  clearDirectResponseAction(): void;
  getDirectResponseAction(): github_com_solo_io_gloo_projects_gloo_api_v1_proxy_pb.DirectResponseAction | undefined;
  setDirectResponseAction(value?: github_com_solo_io_gloo_projects_gloo_api_v1_proxy_pb.DirectResponseAction): void;

  hasDelegateAction(): boolean;
  clearDelegateAction(): void;
  getDelegateAction(): github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef | undefined;
  setDelegateAction(value?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef): void;

  hasRoutePlugins(): boolean;
  clearRoutePlugins(): void;
  getRoutePlugins(): github_com_solo_io_gloo_projects_gloo_api_v1_plugins_pb.RoutePlugins | undefined;
  setRoutePlugins(value?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_pb.RoutePlugins): void;

  getActionCase(): Route.ActionCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Route.AsObject;
  static toObject(includeInstance: boolean, msg: Route): Route.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Route, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Route;
  static deserializeBinaryFromReader(message: Route, reader: jspb.BinaryReader): Route;
}

export namespace Route {
  export type AsObject = {
    matcher?: github_com_solo_io_gloo_projects_gloo_api_v1_proxy_pb.Matcher.AsObject,
    routeAction?: github_com_solo_io_gloo_projects_gloo_api_v1_proxy_pb.RouteAction.AsObject,
    redirectAction?: github_com_solo_io_gloo_projects_gloo_api_v1_proxy_pb.RedirectAction.AsObject,
    directResponseAction?: github_com_solo_io_gloo_projects_gloo_api_v1_proxy_pb.DirectResponseAction.AsObject,
    delegateAction?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject,
    routePlugins?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_pb.RoutePlugins.AsObject,
  }

  export enum ActionCase {
    ACTION_NOT_SET = 0,
    ROUTE_ACTION = 2,
    REDIRECT_ACTION = 3,
    DIRECT_RESPONSE_ACTION = 4,
    DELEGATE_ACTION = 5,
  }
}

