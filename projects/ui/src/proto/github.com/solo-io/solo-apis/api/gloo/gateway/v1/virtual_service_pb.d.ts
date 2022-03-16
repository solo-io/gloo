/* eslint-disable */
// package: gateway.solo.io
// file: github.com/solo-io/solo-apis/api/gloo/gateway/v1/virtual_service.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_struct_pb from "google-protobuf/google/protobuf/struct_pb";
import * as extproto_ext_pb from "../../../../../../../extproto/ext_pb";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";
import * as github_com_solo_io_solo_kit_api_v1_ref_pb from "../../../../../../../github.com/solo-io/solo-kit/api/v1/ref_pb";
import * as github_com_solo_io_solo_kit_api_v1_solo_kit_pb from "../../../../../../../github.com/solo-io/solo-kit/api/v1/solo-kit_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_ssl_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/ssl_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_proxy_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/proxy_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/options_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_core_matchers_matchers_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/core/matchers/matchers_pb";

export class VirtualServiceSpec extends jspb.Message {
  hasVirtualHost(): boolean;
  clearVirtualHost(): void;
  getVirtualHost(): VirtualHost | undefined;
  setVirtualHost(value?: VirtualHost): void;

  hasSslConfig(): boolean;
  clearSslConfig(): void;
  getSslConfig(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_ssl_pb.SslConfig | undefined;
  setSslConfig(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_ssl_pb.SslConfig): void;

  getDisplayName(): string;
  setDisplayName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VirtualServiceSpec.AsObject;
  static toObject(includeInstance: boolean, msg: VirtualServiceSpec): VirtualServiceSpec.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: VirtualServiceSpec, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VirtualServiceSpec;
  static deserializeBinaryFromReader(message: VirtualServiceSpec, reader: jspb.BinaryReader): VirtualServiceSpec;
}

export namespace VirtualServiceSpec {
  export type AsObject = {
    virtualHost?: VirtualHost.AsObject,
    sslConfig?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_ssl_pb.SslConfig.AsObject,
    displayName: string,
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
  getOptions(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.VirtualHostOptions | undefined;
  setOptions(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.VirtualHostOptions): void;

  hasOptionsConfigRefs(): boolean;
  clearOptionsConfigRefs(): void;
  getOptionsConfigRefs(): DelegateOptionsRefs | undefined;
  setOptionsConfigRefs(value?: DelegateOptionsRefs): void;

  getExternalOptionsConfigCase(): VirtualHost.ExternalOptionsConfigCase;
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
    options?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.VirtualHostOptions.AsObject,
    optionsConfigRefs?: DelegateOptionsRefs.AsObject,
  }

  export enum ExternalOptionsConfigCase {
    EXTERNAL_OPTIONS_CONFIG_NOT_SET = 0,
    OPTIONS_CONFIG_REFS = 5,
  }
}

export class Route extends jspb.Message {
  clearMatchersList(): void;
  getMatchersList(): Array<github_com_solo_io_solo_apis_api_gloo_gloo_v1_core_matchers_matchers_pb.Matcher>;
  setMatchersList(value: Array<github_com_solo_io_solo_apis_api_gloo_gloo_v1_core_matchers_matchers_pb.Matcher>): void;
  addMatchers(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_core_matchers_matchers_pb.Matcher, index?: number): github_com_solo_io_solo_apis_api_gloo_gloo_v1_core_matchers_matchers_pb.Matcher;

  hasInheritableMatchers(): boolean;
  clearInheritableMatchers(): void;
  getInheritableMatchers(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setInheritableMatchers(value?: google_protobuf_wrappers_pb.BoolValue): void;

  hasInheritablePathMatchers(): boolean;
  clearInheritablePathMatchers(): void;
  getInheritablePathMatchers(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setInheritablePathMatchers(value?: google_protobuf_wrappers_pb.BoolValue): void;

  hasRouteAction(): boolean;
  clearRouteAction(): void;
  getRouteAction(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_proxy_pb.RouteAction | undefined;
  setRouteAction(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_proxy_pb.RouteAction): void;

  hasRedirectAction(): boolean;
  clearRedirectAction(): void;
  getRedirectAction(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_proxy_pb.RedirectAction | undefined;
  setRedirectAction(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_proxy_pb.RedirectAction): void;

  hasDirectResponseAction(): boolean;
  clearDirectResponseAction(): void;
  getDirectResponseAction(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_proxy_pb.DirectResponseAction | undefined;
  setDirectResponseAction(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_proxy_pb.DirectResponseAction): void;

  hasDelegateAction(): boolean;
  clearDelegateAction(): void;
  getDelegateAction(): DelegateAction | undefined;
  setDelegateAction(value?: DelegateAction): void;

  hasGraphqlApiRef(): boolean;
  clearGraphqlApiRef(): void;
  getGraphqlApiRef(): github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef | undefined;
  setGraphqlApiRef(value?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef): void;

  hasOptions(): boolean;
  clearOptions(): void;
  getOptions(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.RouteOptions | undefined;
  setOptions(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.RouteOptions): void;

  getName(): string;
  setName(value: string): void;

  hasOptionsConfigRefs(): boolean;
  clearOptionsConfigRefs(): void;
  getOptionsConfigRefs(): DelegateOptionsRefs | undefined;
  setOptionsConfigRefs(value?: DelegateOptionsRefs): void;

  getActionCase(): Route.ActionCase;
  getExternalOptionsConfigCase(): Route.ExternalOptionsConfigCase;
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
    matchersList: Array<github_com_solo_io_solo_apis_api_gloo_gloo_v1_core_matchers_matchers_pb.Matcher.AsObject>,
    inheritableMatchers?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    inheritablePathMatchers?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    routeAction?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_proxy_pb.RouteAction.AsObject,
    redirectAction?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_proxy_pb.RedirectAction.AsObject,
    directResponseAction?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_proxy_pb.DirectResponseAction.AsObject,
    delegateAction?: DelegateAction.AsObject,
    graphqlApiRef?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject,
    options?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.RouteOptions.AsObject,
    name: string,
    optionsConfigRefs?: DelegateOptionsRefs.AsObject,
  }

  export enum ActionCase {
    ACTION_NOT_SET = 0,
    ROUTE_ACTION = 2,
    REDIRECT_ACTION = 3,
    DIRECT_RESPONSE_ACTION = 4,
    DELEGATE_ACTION = 5,
    GRAPHQL_API_REF = 11,
  }

  export enum ExternalOptionsConfigCase {
    EXTERNAL_OPTIONS_CONFIG_NOT_SET = 0,
    OPTIONS_CONFIG_REFS = 10,
  }
}

export class DelegateOptionsRefs extends jspb.Message {
  clearDelegateOptionsList(): void;
  getDelegateOptionsList(): Array<github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef>;
  setDelegateOptionsList(value: Array<github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef>): void;
  addDelegateOptions(value?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef, index?: number): github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DelegateOptionsRefs.AsObject;
  static toObject(includeInstance: boolean, msg: DelegateOptionsRefs): DelegateOptionsRefs.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DelegateOptionsRefs, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DelegateOptionsRefs;
  static deserializeBinaryFromReader(message: DelegateOptionsRefs, reader: jspb.BinaryReader): DelegateOptionsRefs;
}

export namespace DelegateOptionsRefs {
  export type AsObject = {
    delegateOptionsList: Array<github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject>,
  }
}

export class DelegateAction extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getNamespace(): string;
  setNamespace(value: string): void;

  hasRef(): boolean;
  clearRef(): void;
  getRef(): github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef | undefined;
  setRef(value?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef): void;

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
    ref?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject,
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
  clearExpressionsList(): void;
  getExpressionsList(): Array<RouteTableSelector.Expression>;
  setExpressionsList(value: Array<RouteTableSelector.Expression>): void;
  addExpressions(value?: RouteTableSelector.Expression, index?: number): RouteTableSelector.Expression;

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
    expressionsList: Array<RouteTableSelector.Expression.AsObject>,
  }

  export class Expression extends jspb.Message {
    getKey(): string;
    setKey(value: string): void;

    getOperator(): RouteTableSelector.Expression.OperatorMap[keyof RouteTableSelector.Expression.OperatorMap];
    setOperator(value: RouteTableSelector.Expression.OperatorMap[keyof RouteTableSelector.Expression.OperatorMap]): void;

    clearValuesList(): void;
    getValuesList(): Array<string>;
    setValuesList(value: Array<string>): void;
    addValues(value: string, index?: number): string;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Expression.AsObject;
    static toObject(includeInstance: boolean, msg: Expression): Expression.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Expression, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Expression;
    static deserializeBinaryFromReader(message: Expression, reader: jspb.BinaryReader): Expression;
  }

  export namespace Expression {
    export type AsObject = {
      key: string,
      operator: RouteTableSelector.Expression.OperatorMap[keyof RouteTableSelector.Expression.OperatorMap],
      valuesList: Array<string>,
    }

    export interface OperatorMap {
      EQUALS: 0;
      DOUBLEEQUALS: 1;
      NOTEQUALS: 2;
      IN: 3;
      NOTIN: 4;
      EXISTS: 5;
      DOESNOTEXIST: 6;
      GREATERTHAN: 7;
      LESSTHAN: 8;
    }

    export const Operator: OperatorMap;
  }
}

export class VirtualServiceStatus extends jspb.Message {
  getState(): VirtualServiceStatus.StateMap[keyof VirtualServiceStatus.StateMap];
  setState(value: VirtualServiceStatus.StateMap[keyof VirtualServiceStatus.StateMap]): void;

  getReason(): string;
  setReason(value: string): void;

  getReportedBy(): string;
  setReportedBy(value: string): void;

  getSubresourceStatusesMap(): jspb.Map<string, VirtualServiceStatus>;
  clearSubresourceStatusesMap(): void;
  hasDetails(): boolean;
  clearDetails(): void;
  getDetails(): google_protobuf_struct_pb.Struct | undefined;
  setDetails(value?: google_protobuf_struct_pb.Struct): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VirtualServiceStatus.AsObject;
  static toObject(includeInstance: boolean, msg: VirtualServiceStatus): VirtualServiceStatus.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: VirtualServiceStatus, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VirtualServiceStatus;
  static deserializeBinaryFromReader(message: VirtualServiceStatus, reader: jspb.BinaryReader): VirtualServiceStatus;
}

export namespace VirtualServiceStatus {
  export type AsObject = {
    state: VirtualServiceStatus.StateMap[keyof VirtualServiceStatus.StateMap],
    reason: string,
    reportedBy: string,
    subresourceStatusesMap: Array<[string, VirtualServiceStatus.AsObject]>,
    details?: google_protobuf_struct_pb.Struct.AsObject,
  }

  export interface StateMap {
    PENDING: 0;
    ACCEPTED: 1;
    REJECTED: 2;
    WARNING: 3;
  }

  export const State: StateMap;
}

export class VirtualServiceNamespacedStatuses extends jspb.Message {
  getStatusesMap(): jspb.Map<string, VirtualServiceStatus>;
  clearStatusesMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VirtualServiceNamespacedStatuses.AsObject;
  static toObject(includeInstance: boolean, msg: VirtualServiceNamespacedStatuses): VirtualServiceNamespacedStatuses.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: VirtualServiceNamespacedStatuses, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VirtualServiceNamespacedStatuses;
  static deserializeBinaryFromReader(message: VirtualServiceNamespacedStatuses, reader: jspb.BinaryReader): VirtualServiceNamespacedStatuses;
}

export namespace VirtualServiceNamespacedStatuses {
  export type AsObject = {
    statusesMap: Array<[string, VirtualServiceStatus.AsObject]>,
  }
}
