/* eslint-disable */
// package: gateway.solo.io
// file: github.com/solo-io/solo-apis/api/gloo/gateway/v1/gateway.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_struct_pb from "google-protobuf/google/protobuf/struct_pb";
import * as extproto_ext_pb from "../../../../../../../extproto/ext_pb";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";
import * as github_com_solo_io_solo_kit_api_v1_ref_pb from "../../../../../../../github.com/solo-io/solo-kit/api/v1/ref_pb";
import * as github_com_solo_io_solo_kit_api_v1_solo_kit_pb from "../../../../../../../github.com/solo-io/solo-kit/api/v1/solo-kit_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_proxy_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/proxy_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/options_pb";

export class GatewaySpec extends jspb.Message {
  getSsl(): boolean;
  setSsl(value: boolean): void;

  getBindAddress(): string;
  setBindAddress(value: string): void;

  getBindPort(): number;
  setBindPort(value: number): void;

  hasOptions(): boolean;
  clearOptions(): void;
  getOptions(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.ListenerOptions | undefined;
  setOptions(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.ListenerOptions): void;

  hasUseProxyProto(): boolean;
  clearUseProxyProto(): void;
  getUseProxyProto(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setUseProxyProto(value?: google_protobuf_wrappers_pb.BoolValue): void;

  hasHttpGateway(): boolean;
  clearHttpGateway(): void;
  getHttpGateway(): HttpGateway | undefined;
  setHttpGateway(value?: HttpGateway): void;

  hasTcpGateway(): boolean;
  clearTcpGateway(): void;
  getTcpGateway(): TcpGateway | undefined;
  setTcpGateway(value?: TcpGateway): void;

  clearProxyNamesList(): void;
  getProxyNamesList(): Array<string>;
  setProxyNamesList(value: Array<string>): void;
  addProxyNames(value: string, index?: number): string;

  hasRouteOptions(): boolean;
  clearRouteOptions(): void;
  getRouteOptions(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.RouteConfigurationOptions | undefined;
  setRouteOptions(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.RouteConfigurationOptions): void;

  getGatewaytypeCase(): GatewaySpec.GatewaytypeCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GatewaySpec.AsObject;
  static toObject(includeInstance: boolean, msg: GatewaySpec): GatewaySpec.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GatewaySpec, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GatewaySpec;
  static deserializeBinaryFromReader(message: GatewaySpec, reader: jspb.BinaryReader): GatewaySpec;
}

export namespace GatewaySpec {
  export type AsObject = {
    ssl: boolean,
    bindAddress: string,
    bindPort: number,
    options?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.ListenerOptions.AsObject,
    useProxyProto?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    httpGateway?: HttpGateway.AsObject,
    tcpGateway?: TcpGateway.AsObject,
    proxyNamesList: Array<string>,
    routeOptions?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.RouteConfigurationOptions.AsObject,
  }

  export enum GatewaytypeCase {
    GATEWAYTYPE_NOT_SET = 0,
    HTTP_GATEWAY = 9,
    TCP_GATEWAY = 10,
  }
}

export class HttpGateway extends jspb.Message {
  clearVirtualServicesList(): void;
  getVirtualServicesList(): Array<github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef>;
  setVirtualServicesList(value: Array<github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef>): void;
  addVirtualServices(value?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef, index?: number): github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef;

  getVirtualServiceSelectorMap(): jspb.Map<string, string>;
  clearVirtualServiceSelectorMap(): void;
  clearVirtualServiceNamespacesList(): void;
  getVirtualServiceNamespacesList(): Array<string>;
  setVirtualServiceNamespacesList(value: Array<string>): void;
  addVirtualServiceNamespaces(value: string, index?: number): string;

  hasOptions(): boolean;
  clearOptions(): void;
  getOptions(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.HttpListenerOptions | undefined;
  setOptions(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.HttpListenerOptions): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HttpGateway.AsObject;
  static toObject(includeInstance: boolean, msg: HttpGateway): HttpGateway.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: HttpGateway, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HttpGateway;
  static deserializeBinaryFromReader(message: HttpGateway, reader: jspb.BinaryReader): HttpGateway;
}

export namespace HttpGateway {
  export type AsObject = {
    virtualServicesList: Array<github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject>,
    virtualServiceSelectorMap: Array<[string, string]>,
    virtualServiceNamespacesList: Array<string>,
    options?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.HttpListenerOptions.AsObject,
  }
}

export class TcpGateway extends jspb.Message {
  clearTcpHostsList(): void;
  getTcpHostsList(): Array<github_com_solo_io_solo_apis_api_gloo_gloo_v1_proxy_pb.TcpHost>;
  setTcpHostsList(value: Array<github_com_solo_io_solo_apis_api_gloo_gloo_v1_proxy_pb.TcpHost>): void;
  addTcpHosts(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_proxy_pb.TcpHost, index?: number): github_com_solo_io_solo_apis_api_gloo_gloo_v1_proxy_pb.TcpHost;

  hasOptions(): boolean;
  clearOptions(): void;
  getOptions(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.TcpListenerOptions | undefined;
  setOptions(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.TcpListenerOptions): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TcpGateway.AsObject;
  static toObject(includeInstance: boolean, msg: TcpGateway): TcpGateway.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: TcpGateway, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TcpGateway;
  static deserializeBinaryFromReader(message: TcpGateway, reader: jspb.BinaryReader): TcpGateway;
}

export namespace TcpGateway {
  export type AsObject = {
    tcpHostsList: Array<github_com_solo_io_solo_apis_api_gloo_gloo_v1_proxy_pb.TcpHost.AsObject>,
    options?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.TcpListenerOptions.AsObject,
  }
}

export class GatewayStatus extends jspb.Message {
  getState(): GatewayStatus.StateMap[keyof GatewayStatus.StateMap];
  setState(value: GatewayStatus.StateMap[keyof GatewayStatus.StateMap]): void;

  getReason(): string;
  setReason(value: string): void;

  getReportedBy(): string;
  setReportedBy(value: string): void;

  getSubresourceStatusesMap(): jspb.Map<string, GatewayStatus>;
  clearSubresourceStatusesMap(): void;
  hasDetails(): boolean;
  clearDetails(): void;
  getDetails(): google_protobuf_struct_pb.Struct | undefined;
  setDetails(value?: google_protobuf_struct_pb.Struct): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GatewayStatus.AsObject;
  static toObject(includeInstance: boolean, msg: GatewayStatus): GatewayStatus.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GatewayStatus, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GatewayStatus;
  static deserializeBinaryFromReader(message: GatewayStatus, reader: jspb.BinaryReader): GatewayStatus;
}

export namespace GatewayStatus {
  export type AsObject = {
    state: GatewayStatus.StateMap[keyof GatewayStatus.StateMap],
    reason: string,
    reportedBy: string,
    subresourceStatusesMap: Array<[string, GatewayStatus.AsObject]>,
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
