// package: gateway.solo.io
// file: gloo/projects/gateway/api/v1/gateway.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../../../gogoproto/gogo_pb";
import * as extproto_ext_pb from "../../../../../extproto/ext_pb";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";
import * as solo_kit_api_v1_metadata_pb from "../../../../../solo-kit/api/v1/metadata_pb";
import * as solo_kit_api_v1_status_pb from "../../../../../solo-kit/api/v1/status_pb";
import * as solo_kit_api_v1_ref_pb from "../../../../../solo-kit/api/v1/ref_pb";
import * as solo_kit_api_v1_solo_kit_pb from "../../../../../solo-kit/api/v1/solo-kit_pb";
import * as gloo_projects_gloo_api_v1_proxy_pb from "../../../../../gloo/projects/gloo/api/v1/proxy_pb";
import * as gloo_projects_gloo_api_v1_options_pb from "../../../../../gloo/projects/gloo/api/v1/options_pb";

export class Gateway extends jspb.Message {
  getSsl(): boolean;
  setSsl(value: boolean): void;

  getBindAddress(): string;
  setBindAddress(value: string): void;

  getBindPort(): number;
  setBindPort(value: number): void;

  hasOptions(): boolean;
  clearOptions(): void;
  getOptions(): gloo_projects_gloo_api_v1_options_pb.ListenerOptions | undefined;
  setOptions(value?: gloo_projects_gloo_api_v1_options_pb.ListenerOptions): void;

  hasStatus(): boolean;
  clearStatus(): void;
  getStatus(): solo_kit_api_v1_status_pb.Status | undefined;
  setStatus(value?: solo_kit_api_v1_status_pb.Status): void;

  hasMetadata(): boolean;
  clearMetadata(): void;
  getMetadata(): solo_kit_api_v1_metadata_pb.Metadata | undefined;
  setMetadata(value?: solo_kit_api_v1_metadata_pb.Metadata): void;

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

  getGatewaytypeCase(): Gateway.GatewaytypeCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Gateway.AsObject;
  static toObject(includeInstance: boolean, msg: Gateway): Gateway.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Gateway, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Gateway;
  static deserializeBinaryFromReader(message: Gateway, reader: jspb.BinaryReader): Gateway;
}

export namespace Gateway {
  export type AsObject = {
    ssl: boolean,
    bindAddress: string,
    bindPort: number,
    options?: gloo_projects_gloo_api_v1_options_pb.ListenerOptions.AsObject,
    status?: solo_kit_api_v1_status_pb.Status.AsObject,
    metadata?: solo_kit_api_v1_metadata_pb.Metadata.AsObject,
    useProxyProto?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    httpGateway?: HttpGateway.AsObject,
    tcpGateway?: TcpGateway.AsObject,
    proxyNamesList: Array<string>,
  }

  export enum GatewaytypeCase {
    GATEWAYTYPE_NOT_SET = 0,
    HTTP_GATEWAY = 9,
    TCP_GATEWAY = 10,
  }
}

export class HttpGateway extends jspb.Message {
  clearVirtualServicesList(): void;
  getVirtualServicesList(): Array<solo_kit_api_v1_ref_pb.ResourceRef>;
  setVirtualServicesList(value: Array<solo_kit_api_v1_ref_pb.ResourceRef>): void;
  addVirtualServices(value?: solo_kit_api_v1_ref_pb.ResourceRef, index?: number): solo_kit_api_v1_ref_pb.ResourceRef;

  getVirtualServiceSelectorMap(): jspb.Map<string, string>;
  clearVirtualServiceSelectorMap(): void;
  clearVirtualServiceNamespacesList(): void;
  getVirtualServiceNamespacesList(): Array<string>;
  setVirtualServiceNamespacesList(value: Array<string>): void;
  addVirtualServiceNamespaces(value: string, index?: number): string;

  hasOptions(): boolean;
  clearOptions(): void;
  getOptions(): gloo_projects_gloo_api_v1_options_pb.HttpListenerOptions | undefined;
  setOptions(value?: gloo_projects_gloo_api_v1_options_pb.HttpListenerOptions): void;

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
    virtualServicesList: Array<solo_kit_api_v1_ref_pb.ResourceRef.AsObject>,
    virtualServiceSelectorMap: Array<[string, string]>,
    virtualServiceNamespacesList: Array<string>,
    options?: gloo_projects_gloo_api_v1_options_pb.HttpListenerOptions.AsObject,
  }
}

export class TcpGateway extends jspb.Message {
  clearTcpHostsList(): void;
  getTcpHostsList(): Array<gloo_projects_gloo_api_v1_proxy_pb.TcpHost>;
  setTcpHostsList(value: Array<gloo_projects_gloo_api_v1_proxy_pb.TcpHost>): void;
  addTcpHosts(value?: gloo_projects_gloo_api_v1_proxy_pb.TcpHost, index?: number): gloo_projects_gloo_api_v1_proxy_pb.TcpHost;

  hasOptions(): boolean;
  clearOptions(): void;
  getOptions(): gloo_projects_gloo_api_v1_options_pb.TcpListenerOptions | undefined;
  setOptions(value?: gloo_projects_gloo_api_v1_options_pb.TcpListenerOptions): void;

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
    tcpHostsList: Array<gloo_projects_gloo_api_v1_proxy_pb.TcpHost.AsObject>,
    options?: gloo_projects_gloo_api_v1_options_pb.TcpListenerOptions.AsObject,
  }
}

