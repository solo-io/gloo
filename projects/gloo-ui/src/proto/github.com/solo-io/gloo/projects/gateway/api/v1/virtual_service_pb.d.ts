// package: gateway.solo.io
// file: github.com/solo-io/gloo/projects/gateway/api/v1/virtual_service.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../../../../../gogoproto/gogo_pb";
import * as github_com_solo_io_solo_kit_api_v1_metadata_pb from "../../../../../../../github.com/solo-io/solo-kit/api/v1/metadata_pb";
import * as github_com_solo_io_solo_kit_api_v1_status_pb from "../../../../../../../github.com/solo-io/solo-kit/api/v1/status_pb";
import * as github_com_solo_io_solo_kit_api_v1_solo_kit_pb from "../../../../../../../github.com/solo-io/solo-kit/api/v1/solo-kit_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_v1_ssl_pb from "../../../../../../../github.com/solo-io/gloo/projects/gloo/api/v1/ssl_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_v1_proxy_pb from "../../../../../../../github.com/solo-io/gloo/projects/gloo/api/v1/proxy_pb";

export class VirtualService extends jspb.Message {
  hasVirtualHost(): boolean;
  clearVirtualHost(): void;
  getVirtualHost(): github_com_solo_io_gloo_projects_gloo_api_v1_proxy_pb.VirtualHost | undefined;
  setVirtualHost(value?: github_com_solo_io_gloo_projects_gloo_api_v1_proxy_pb.VirtualHost): void;

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
    virtualHost?: github_com_solo_io_gloo_projects_gloo_api_v1_proxy_pb.VirtualHost.AsObject,
    sslConfig?: github_com_solo_io_gloo_projects_gloo_api_v1_ssl_pb.SslConfig.AsObject,
    displayName: string,
    status?: github_com_solo_io_solo_kit_api_v1_status_pb.Status.AsObject,
    metadata?: github_com_solo_io_solo_kit_api_v1_metadata_pb.Metadata.AsObject,
  }
}

