// package: gateway.solo.io
// file: github.com/solo-io/gloo/projects/gateway/api/v1/gateway.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";
import * as gogoproto_gogo_pb from "../../../../../../../gogoproto/gogo_pb";
import * as github_com_solo_io_solo_kit_api_v1_metadata_pb from "../../../../../../../github.com/solo-io/solo-kit/api/v1/metadata_pb";
import * as github_com_solo_io_solo_kit_api_v1_status_pb from "../../../../../../../github.com/solo-io/solo-kit/api/v1/status_pb";
import * as github_com_solo_io_solo_kit_api_v1_ref_pb from "../../../../../../../github.com/solo-io/solo-kit/api/v1/ref_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_v1_proxy_pb from "../../../../../../../github.com/solo-io/gloo/projects/gloo/api/v1/proxy_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_v1_plugins_pb from "../../../../../../../github.com/solo-io/gloo/projects/gloo/api/v1/plugins_pb";

export class Gateway extends jspb.Message {
  getSsl(): boolean;
  setSsl(value: boolean): void;

  clearVirtualServicesList(): void;
  getVirtualServicesList(): Array<github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef>;
  setVirtualServicesList(value: Array<github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef>): void;
  addVirtualServices(value?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef, index?: number): github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef;

  getBindAddress(): string;
  setBindAddress(value: string): void;

  getBindPort(): number;
  setBindPort(value: number): void;

  hasPlugins(): boolean;
  clearPlugins(): void;
  getPlugins(): github_com_solo_io_gloo_projects_gloo_api_v1_plugins_pb.ListenerPlugins | undefined;
  setPlugins(value?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_pb.ListenerPlugins): void;

  hasStatus(): boolean;
  clearStatus(): void;
  getStatus(): github_com_solo_io_solo_kit_api_v1_status_pb.Status | undefined;
  setStatus(value?: github_com_solo_io_solo_kit_api_v1_status_pb.Status): void;

  hasMetadata(): boolean;
  clearMetadata(): void;
  getMetadata(): github_com_solo_io_solo_kit_api_v1_metadata_pb.Metadata | undefined;
  setMetadata(value?: github_com_solo_io_solo_kit_api_v1_metadata_pb.Metadata): void;

  hasUseProxyProto(): boolean;
  clearUseProxyProto(): void;
  getUseProxyProto(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setUseProxyProto(value?: google_protobuf_wrappers_pb.BoolValue): void;

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
    virtualServicesList: Array<github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject>,
    bindAddress: string,
    bindPort: number,
    plugins?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_pb.ListenerPlugins.AsObject,
    status?: github_com_solo_io_solo_kit_api_v1_status_pb.Status.AsObject,
    metadata?: github_com_solo_io_solo_kit_api_v1_metadata_pb.Metadata.AsObject,
    useProxyProto?: google_protobuf_wrappers_pb.BoolValue.AsObject,
  }
}

