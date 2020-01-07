// package: gateway.solo.io
// file: gloo/projects/gateway/api/v1/virtual_service.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../../../gogoproto/gogo_pb";
import * as extproto_ext_pb from "../../../../../extproto/ext_pb";
import * as solo_kit_api_v1_metadata_pb from "../../../../../solo-kit/api/v1/metadata_pb";
import * as solo_kit_api_v1_status_pb from "../../../../../solo-kit/api/v1/status_pb";
import * as solo_kit_api_v1_ref_pb from "../../../../../solo-kit/api/v1/ref_pb";
import * as solo_kit_api_v1_solo_kit_pb from "../../../../../solo-kit/api/v1/solo-kit_pb";
import * as gloo_projects_gloo_api_v1_ssl_pb from "../../../../../gloo/projects/gloo/api/v1/ssl_pb";
import * as gloo_projects_gloo_api_v1_proxy_pb from "../../../../../gloo/projects/gloo/api/v1/proxy_pb";
import * as gloo_projects_gloo_api_v1_options_pb from "../../../../../gloo/projects/gloo/api/v1/options_pb";
import * as gloo_projects_gloo_api_v1_core_matchers_matchers_pb from "../../../../../gloo/projects/gloo/api/v1/core/matchers/matchers_pb";

export class VirtualService extends jspb.Message {
  hasVirtualHost(): boolean;
  clearVirtualHost(): void;
  getVirtualHost(): VirtualHost | undefined;
  setVirtualHost(value?: VirtualHost): void;

  hasSslConfig(): boolean;
  clearSslConfig(): void;
  getSslConfig(): gloo_projects_gloo_api_v1_ssl_pb.SslConfig | undefined;
  setSslConfig(value?: gloo_projects_gloo_api_v1_ssl_pb.SslConfig): void;

  getDisplayName(): string;
  setDisplayName(value: string): void;

  hasStatus(): boolean;
  clearStatus(): void;
  getStatus(): solo_kit_api_v1_status_pb.Status | undefined;
  setStatus(value?: solo_kit_api_v1_status_pb.Status): void;

  hasMetadata(): boolean;
  clearMetadata(): void;
  getMetadata(): solo_kit_api_v1_metadata_pb.Metadata | undefined;
  setMetadata(value?: solo_kit_api_v1_metadata_pb.Metadata): void;

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
    sslConfig?: gloo_projects_gloo_api_v1_ssl_pb.SslConfig.AsObject,
    displayName: string,
    status?: solo_kit_api_v1_status_pb.Status.AsObject,
    metadata?: solo_kit_api_v1_metadata_pb.Metadata.AsObject,
  }
}

export class VirtualHost extends jspb.Message {
  clearDomainsList(): void;
  getDomainsList(): Array<string>;
  setDomainsList(value: Array<string>): void;
  addDomains(value: string, index?: number): string;

  clearRoutesList(): void;
  getRoutesList(): Array<Route>;
  setRoutesList(value: Array<Route>): void;
  addRoutes(value?: Route, index?: number): Route;

  hasOptions(): boolean;
  clearOptions(): void;
  getOptions(): gloo_projects_gloo_api_v1_options_pb.VirtualHostOptions | undefined;
  setOptions(value?: gloo_projects_gloo_api_v1_options_pb.VirtualHostOptions): void;

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
    domainsList: Array<string>,
    routesList: Array<Route.AsObject>,
    options?: gloo_projects_gloo_api_v1_options_pb.VirtualHostOptions.AsObject,
  }
}

export class Route extends jspb.Message {
  clearMatchersList(): void;
  getMatchersList(): Array<gloo_projects_gloo_api_v1_core_matchers_matchers_pb.Matcher>;
  setMatchersList(value: Array<gloo_projects_gloo_api_v1_core_matchers_matchers_pb.Matcher>): void;
  addMatchers(value?: gloo_projects_gloo_api_v1_core_matchers_matchers_pb.Matcher, index?: number): gloo_projects_gloo_api_v1_core_matchers_matchers_pb.Matcher;

  hasRouteAction(): boolean;
  clearRouteAction(): void;
  getRouteAction(): gloo_projects_gloo_api_v1_proxy_pb.RouteAction | undefined;
  setRouteAction(value?: gloo_projects_gloo_api_v1_proxy_pb.RouteAction): void;

  hasRedirectAction(): boolean;
  clearRedirectAction(): void;
  getRedirectAction(): gloo_projects_gloo_api_v1_proxy_pb.RedirectAction | undefined;
  setRedirectAction(value?: gloo_projects_gloo_api_v1_proxy_pb.RedirectAction): void;

  hasDirectResponseAction(): boolean;
  clearDirectResponseAction(): void;
  getDirectResponseAction(): gloo_projects_gloo_api_v1_proxy_pb.DirectResponseAction | undefined;
  setDirectResponseAction(value?: gloo_projects_gloo_api_v1_proxy_pb.DirectResponseAction): void;

  hasDelegateAction(): boolean;
  clearDelegateAction(): void;
  getDelegateAction(): DelegateAction | undefined;
  setDelegateAction(value?: DelegateAction): void;

  hasOptions(): boolean;
  clearOptions(): void;
  getOptions(): gloo_projects_gloo_api_v1_options_pb.RouteOptions | undefined;
  setOptions(value?: gloo_projects_gloo_api_v1_options_pb.RouteOptions): void;

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
    matchersList: Array<gloo_projects_gloo_api_v1_core_matchers_matchers_pb.Matcher.AsObject>,
    routeAction?: gloo_projects_gloo_api_v1_proxy_pb.RouteAction.AsObject,
    redirectAction?: gloo_projects_gloo_api_v1_proxy_pb.RedirectAction.AsObject,
    directResponseAction?: gloo_projects_gloo_api_v1_proxy_pb.DirectResponseAction.AsObject,
    delegateAction?: DelegateAction.AsObject,
    options?: gloo_projects_gloo_api_v1_options_pb.RouteOptions.AsObject,
  }

  export enum ActionCase {
    ACTION_NOT_SET = 0,
    ROUTE_ACTION = 2,
    REDIRECT_ACTION = 3,
    DIRECT_RESPONSE_ACTION = 4,
    DELEGATE_ACTION = 5,
  }
}

export class DelegateAction extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getNamespace(): string;
  setNamespace(value: string): void;

  hasRef(): boolean;
  clearRef(): void;
  getRef(): solo_kit_api_v1_ref_pb.ResourceRef | undefined;
  setRef(value?: solo_kit_api_v1_ref_pb.ResourceRef): void;

  hasSelector(): boolean;
  clearSelector(): void;
  getSelector(): RouteTableSelector | undefined;
  setSelector(value?: RouteTableSelector): void;

  getDelegationTypeCase(): DelegateAction.DelegationTypeCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DelegateAction.AsObject;
  static toObject(includeInstance: boolean, msg: DelegateAction): DelegateAction.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DelegateAction, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DelegateAction;
  static deserializeBinaryFromReader(message: DelegateAction, reader: jspb.BinaryReader): DelegateAction;
}

export namespace DelegateAction {
  export type AsObject = {
    name: string,
    namespace: string,
    ref?: solo_kit_api_v1_ref_pb.ResourceRef.AsObject,
    selector?: RouteTableSelector.AsObject,
  }

  export enum DelegationTypeCase {
    DELEGATION_TYPE_NOT_SET = 0,
    REF = 3,
    SELECTOR = 4,
  }
}

export class RouteTableSelector extends jspb.Message {
  clearNamespacesList(): void;
  getNamespacesList(): Array<string>;
  setNamespacesList(value: Array<string>): void;
  addNamespaces(value: string, index?: number): string;

  getLabelsMap(): jspb.Map<string, string>;
  clearLabelsMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RouteTableSelector.AsObject;
  static toObject(includeInstance: boolean, msg: RouteTableSelector): RouteTableSelector.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RouteTableSelector, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RouteTableSelector;
  static deserializeBinaryFromReader(message: RouteTableSelector, reader: jspb.BinaryReader): RouteTableSelector;
}

export namespace RouteTableSelector {
  export type AsObject = {
    namespacesList: Array<string>,
    labelsMap: Array<[string, string]>,
  }
}

