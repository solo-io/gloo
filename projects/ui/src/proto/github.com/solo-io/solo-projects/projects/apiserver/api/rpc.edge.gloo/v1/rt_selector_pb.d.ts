/* eslint-disable */
// package: rpc.edge.gloo.solo.io
// file: github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/rt_selector.proto

import * as jspb from "google-protobuf";
import * as extproto_ext_pb from "../../../../../../../../extproto/ext_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_core_matchers_matchers_pb from "../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/core/matchers/matchers_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_proxy_pb from "../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/proxy_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gateway_v1_virtual_service_pb from "../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gateway/v1/virtual_service_pb";
import * as github_com_solo_io_skv2_api_core_v1_core_pb from "../../../../../../../../github.com/solo-io/skv2/api/core/v1/core_pb";

export class SubRouteTableRow extends jspb.Message {
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
  getDelegateAction(): github_com_solo_io_solo_apis_api_gloo_gateway_v1_virtual_service_pb.DelegateAction | undefined;
  setDelegateAction(value?: github_com_solo_io_solo_apis_api_gloo_gateway_v1_virtual_service_pb.DelegateAction): void;

  getMatcher(): string;
  setMatcher(value: string): void;

  getMatchType(): string;
  setMatchType(value: string): void;

  clearMethodsList(): void;
  getMethodsList(): Array<string>;
  setMethodsList(value: Array<string>): void;
  addMethods(value: string, index?: number): string;

  clearHeadersList(): void;
  getHeadersList(): Array<github_com_solo_io_solo_apis_api_gloo_gloo_v1_core_matchers_matchers_pb.HeaderMatcher>;
  setHeadersList(value: Array<github_com_solo_io_solo_apis_api_gloo_gloo_v1_core_matchers_matchers_pb.HeaderMatcher>): void;
  addHeaders(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_core_matchers_matchers_pb.HeaderMatcher, index?: number): github_com_solo_io_solo_apis_api_gloo_gloo_v1_core_matchers_matchers_pb.HeaderMatcher;

  clearQueryParametersList(): void;
  getQueryParametersList(): Array<github_com_solo_io_solo_apis_api_gloo_gloo_v1_core_matchers_matchers_pb.QueryParameterMatcher>;
  setQueryParametersList(value: Array<github_com_solo_io_solo_apis_api_gloo_gloo_v1_core_matchers_matchers_pb.QueryParameterMatcher>): void;
  addQueryParameters(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_core_matchers_matchers_pb.QueryParameterMatcher, index?: number): github_com_solo_io_solo_apis_api_gloo_gloo_v1_core_matchers_matchers_pb.QueryParameterMatcher;

  clearRtRoutesList(): void;
  getRtRoutesList(): Array<SubRouteTableRow>;
  setRtRoutesList(value: Array<SubRouteTableRow>): void;
  addRtRoutes(value?: SubRouteTableRow, index?: number): SubRouteTableRow;

  getActionCase(): SubRouteTableRow.ActionCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SubRouteTableRow.AsObject;
  static toObject(includeInstance: boolean, msg: SubRouteTableRow): SubRouteTableRow.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: SubRouteTableRow, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SubRouteTableRow;
  static deserializeBinaryFromReader(message: SubRouteTableRow, reader: jspb.BinaryReader): SubRouteTableRow;
}

export namespace SubRouteTableRow {
  export type AsObject = {
    routeAction?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_proxy_pb.RouteAction.AsObject,
    redirectAction?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_proxy_pb.RedirectAction.AsObject,
    directResponseAction?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_proxy_pb.DirectResponseAction.AsObject,
    delegateAction?: github_com_solo_io_solo_apis_api_gloo_gateway_v1_virtual_service_pb.DelegateAction.AsObject,
    matcher: string,
    matchType: string,
    methodsList: Array<string>,
    headersList: Array<github_com_solo_io_solo_apis_api_gloo_gloo_v1_core_matchers_matchers_pb.HeaderMatcher.AsObject>,
    queryParametersList: Array<github_com_solo_io_solo_apis_api_gloo_gloo_v1_core_matchers_matchers_pb.QueryParameterMatcher.AsObject>,
    rtRoutesList: Array<SubRouteTableRow.AsObject>,
  }

  export enum ActionCase {
    ACTION_NOT_SET = 0,
    ROUTE_ACTION = 1,
    REDIRECT_ACTION = 2,
    DIRECT_RESPONSE_ACTION = 3,
    DELEGATE_ACTION = 4,
  }
}

export class GetVirtualServiceRoutesRequest extends jspb.Message {
  hasVirtualServiceRef(): boolean;
  clearVirtualServiceRef(): void;
  getVirtualServiceRef(): github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef | undefined;
  setVirtualServiceRef(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetVirtualServiceRoutesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetVirtualServiceRoutesRequest): GetVirtualServiceRoutesRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetVirtualServiceRoutesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetVirtualServiceRoutesRequest;
  static deserializeBinaryFromReader(message: GetVirtualServiceRoutesRequest, reader: jspb.BinaryReader): GetVirtualServiceRoutesRequest;
}

export namespace GetVirtualServiceRoutesRequest {
  export type AsObject = {
    virtualServiceRef?: github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef.AsObject,
  }
}

export class GetVirtualServiceRoutesResponse extends jspb.Message {
  clearVsRoutesList(): void;
  getVsRoutesList(): Array<SubRouteTableRow>;
  setVsRoutesList(value: Array<SubRouteTableRow>): void;
  addVsRoutes(value?: SubRouteTableRow, index?: number): SubRouteTableRow;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetVirtualServiceRoutesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetVirtualServiceRoutesResponse): GetVirtualServiceRoutesResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetVirtualServiceRoutesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetVirtualServiceRoutesResponse;
  static deserializeBinaryFromReader(message: GetVirtualServiceRoutesResponse, reader: jspb.BinaryReader): GetVirtualServiceRoutesResponse;
}

export namespace GetVirtualServiceRoutesResponse {
  export type AsObject = {
    vsRoutesList: Array<SubRouteTableRow.AsObject>,
  }
}
