// package: gloo.solo.io
// file: github.com/solo-io/gloo/projects/gloo/api/v1/upstream.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../../../../../gogoproto/gogo_pb";
import * as github_com_solo_io_solo_kit_api_v1_metadata_pb from "../../../../../../../github.com/solo-io/solo-kit/api/v1/metadata_pb";
import * as github_com_solo_io_solo_kit_api_v1_status_pb from "../../../../../../../github.com/solo-io/solo-kit/api/v1/status_pb";
import * as github_com_solo_io_solo_kit_api_v1_solo_kit_pb from "../../../../../../../github.com/solo-io/solo-kit/api/v1/solo-kit_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_v1_plugins_pb from "../../../../../../../github.com/solo-io/gloo/projects/gloo/api/v1/plugins_pb";

export class Upstream extends jspb.Message {
  hasUpstreamSpec(): boolean;
  clearUpstreamSpec(): void;
  getUpstreamSpec(): github_com_solo_io_gloo_projects_gloo_api_v1_plugins_pb.UpstreamSpec | undefined;
  setUpstreamSpec(value?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_pb.UpstreamSpec): void;

  hasStatus(): boolean;
  clearStatus(): void;
  getStatus(): github_com_solo_io_solo_kit_api_v1_status_pb.Status | undefined;
  setStatus(value?: github_com_solo_io_solo_kit_api_v1_status_pb.Status): void;

  hasMetadata(): boolean;
  clearMetadata(): void;
  getMetadata(): github_com_solo_io_solo_kit_api_v1_metadata_pb.Metadata | undefined;
  setMetadata(value?: github_com_solo_io_solo_kit_api_v1_metadata_pb.Metadata): void;

  hasDiscoveryMetadata(): boolean;
  clearDiscoveryMetadata(): void;
  getDiscoveryMetadata(): DiscoveryMetadata | undefined;
  setDiscoveryMetadata(value?: DiscoveryMetadata): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Upstream.AsObject;
  static toObject(includeInstance: boolean, msg: Upstream): Upstream.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Upstream, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Upstream;
  static deserializeBinaryFromReader(message: Upstream, reader: jspb.BinaryReader): Upstream;
}

export namespace Upstream {
  export type AsObject = {
    upstreamSpec?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_pb.UpstreamSpec.AsObject,
    status?: github_com_solo_io_solo_kit_api_v1_status_pb.Status.AsObject,
    metadata?: github_com_solo_io_solo_kit_api_v1_metadata_pb.Metadata.AsObject,
    discoveryMetadata?: DiscoveryMetadata.AsObject,
  }
}

export class DiscoveryMetadata extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DiscoveryMetadata.AsObject;
  static toObject(includeInstance: boolean, msg: DiscoveryMetadata): DiscoveryMetadata.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DiscoveryMetadata, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DiscoveryMetadata;
  static deserializeBinaryFromReader(message: DiscoveryMetadata, reader: jspb.BinaryReader): DiscoveryMetadata;
}

export namespace DiscoveryMetadata {
  export type AsObject = {
  }
}

