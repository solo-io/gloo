/* eslint-disable */
// package: fed.gateway.solo.io
// file: github.com/solo-io/solo-projects/projects/gloo-fed/api/fed.gateway/v1/matchable_tcp_gateway.proto

import * as jspb from "google-protobuf";
import * as extproto_ext_pb from "../../../../../../../../extproto/ext_pb";
import * as github_com_solo_io_solo_projects_projects_gloo_fed_api_multicluster_v1alpha1_multicluster_pb from "../../../../../../../../github.com/solo-io/solo-projects/projects/gloo-fed/api/multicluster/v1alpha1/multicluster_pb";
import * as github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_core_v1_placement_pb from "../../../../../../../../github.com/solo-io/solo-projects/projects/gloo-fed/api/fed/core/v1/placement_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gateway_v1_matchable_tcp_gateway_pb from "../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gateway/v1/matchable_tcp_gateway_pb";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";

export class FederatedMatchableTcpGatewaySpec extends jspb.Message {
  hasTemplate(): boolean;
  clearTemplate(): void;
  getTemplate(): FederatedMatchableTcpGatewaySpec.Template | undefined;
  setTemplate(value?: FederatedMatchableTcpGatewaySpec.Template): void;

  hasPlacement(): boolean;
  clearPlacement(): void;
  getPlacement(): github_com_solo_io_solo_projects_projects_gloo_fed_api_multicluster_v1alpha1_multicluster_pb.Placement | undefined;
  setPlacement(value?: github_com_solo_io_solo_projects_projects_gloo_fed_api_multicluster_v1alpha1_multicluster_pb.Placement): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FederatedMatchableTcpGatewaySpec.AsObject;
  static toObject(includeInstance: boolean, msg: FederatedMatchableTcpGatewaySpec): FederatedMatchableTcpGatewaySpec.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: FederatedMatchableTcpGatewaySpec, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FederatedMatchableTcpGatewaySpec;
  static deserializeBinaryFromReader(message: FederatedMatchableTcpGatewaySpec, reader: jspb.BinaryReader): FederatedMatchableTcpGatewaySpec;
}

export namespace FederatedMatchableTcpGatewaySpec {
  export type AsObject = {
    template?: FederatedMatchableTcpGatewaySpec.Template.AsObject,
    placement?: github_com_solo_io_solo_projects_projects_gloo_fed_api_multicluster_v1alpha1_multicluster_pb.Placement.AsObject,
  }

  export class Template extends jspb.Message {
    hasSpec(): boolean;
    clearSpec(): void;
    getSpec(): github_com_solo_io_solo_apis_api_gloo_gateway_v1_matchable_tcp_gateway_pb.MatchableTcpGatewaySpec | undefined;
    setSpec(value?: github_com_solo_io_solo_apis_api_gloo_gateway_v1_matchable_tcp_gateway_pb.MatchableTcpGatewaySpec): void;

    hasMetadata(): boolean;
    clearMetadata(): void;
    getMetadata(): github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_core_v1_placement_pb.TemplateMetadata | undefined;
    setMetadata(value?: github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_core_v1_placement_pb.TemplateMetadata): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Template.AsObject;
    static toObject(includeInstance: boolean, msg: Template): Template.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Template, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Template;
    static deserializeBinaryFromReader(message: Template, reader: jspb.BinaryReader): Template;
  }

  export namespace Template {
    export type AsObject = {
      spec?: github_com_solo_io_solo_apis_api_gloo_gateway_v1_matchable_tcp_gateway_pb.MatchableTcpGatewaySpec.AsObject,
      metadata?: github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_core_v1_placement_pb.TemplateMetadata.AsObject,
    }
  }
}

export class FederatedMatchableTcpGatewayStatus extends jspb.Message {
  hasPlacementStatus(): boolean;
  clearPlacementStatus(): void;
  getPlacementStatus(): github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_core_v1_placement_pb.PlacementStatus | undefined;
  setPlacementStatus(value?: github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_core_v1_placement_pb.PlacementStatus): void;

  getNamespacedPlacementStatusesMap(): jspb.Map<string, github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_core_v1_placement_pb.PlacementStatus>;
  clearNamespacedPlacementStatusesMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FederatedMatchableTcpGatewayStatus.AsObject;
  static toObject(includeInstance: boolean, msg: FederatedMatchableTcpGatewayStatus): FederatedMatchableTcpGatewayStatus.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: FederatedMatchableTcpGatewayStatus, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FederatedMatchableTcpGatewayStatus;
  static deserializeBinaryFromReader(message: FederatedMatchableTcpGatewayStatus, reader: jspb.BinaryReader): FederatedMatchableTcpGatewayStatus;
}

export namespace FederatedMatchableTcpGatewayStatus {
  export type AsObject = {
    placementStatus?: github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_core_v1_placement_pb.PlacementStatus.AsObject,
    namespacedPlacementStatusesMap: Array<[string, github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_core_v1_placement_pb.PlacementStatus.AsObject]>,
  }
}
