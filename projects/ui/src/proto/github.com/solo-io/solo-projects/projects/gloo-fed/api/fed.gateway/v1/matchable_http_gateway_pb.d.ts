/* eslint-disable */
// package: fed.gateway.solo.io
// file: github.com/solo-io/solo-projects/projects/gloo-fed/api/fed.gateway/v1/matchable_http_gateway.proto

import * as jspb from "google-protobuf";
import * as extproto_ext_pb from "../../../../../../../../extproto/ext_pb";
import * as github_com_solo_io_solo_projects_projects_gloo_fed_api_multicluster_v1alpha1_multicluster_pb from "../../../../../../../../github.com/solo-io/solo-projects/projects/gloo-fed/api/multicluster/v1alpha1/multicluster_pb";
import * as github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_core_v1_placement_pb from "../../../../../../../../github.com/solo-io/solo-projects/projects/gloo-fed/api/fed/core/v1/placement_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gateway_v1_matchable_http_gateway_pb from "../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gateway/v1/matchable_http_gateway_pb";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";

export class FederatedMatchableHttpGatewaySpec extends jspb.Message {
  hasTemplate(): boolean;
  clearTemplate(): void;
  getTemplate(): FederatedMatchableHttpGatewaySpec.Template | undefined;
  setTemplate(value?: FederatedMatchableHttpGatewaySpec.Template): void;

  hasPlacement(): boolean;
  clearPlacement(): void;
  getPlacement(): github_com_solo_io_solo_projects_projects_gloo_fed_api_multicluster_v1alpha1_multicluster_pb.Placement | undefined;
  setPlacement(value?: github_com_solo_io_solo_projects_projects_gloo_fed_api_multicluster_v1alpha1_multicluster_pb.Placement): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FederatedMatchableHttpGatewaySpec.AsObject;
  static toObject(includeInstance: boolean, msg: FederatedMatchableHttpGatewaySpec): FederatedMatchableHttpGatewaySpec.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: FederatedMatchableHttpGatewaySpec, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FederatedMatchableHttpGatewaySpec;
  static deserializeBinaryFromReader(message: FederatedMatchableHttpGatewaySpec, reader: jspb.BinaryReader): FederatedMatchableHttpGatewaySpec;
}

export namespace FederatedMatchableHttpGatewaySpec {
  export type AsObject = {
    template?: FederatedMatchableHttpGatewaySpec.Template.AsObject,
    placement?: github_com_solo_io_solo_projects_projects_gloo_fed_api_multicluster_v1alpha1_multicluster_pb.Placement.AsObject,
  }

  export class Template extends jspb.Message {
    hasSpec(): boolean;
    clearSpec(): void;
    getSpec(): github_com_solo_io_solo_apis_api_gloo_gateway_v1_matchable_http_gateway_pb.MatchableHttpGatewaySpec | undefined;
    setSpec(value?: github_com_solo_io_solo_apis_api_gloo_gateway_v1_matchable_http_gateway_pb.MatchableHttpGatewaySpec): void;

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
      spec?: github_com_solo_io_solo_apis_api_gloo_gateway_v1_matchable_http_gateway_pb.MatchableHttpGatewaySpec.AsObject,
      metadata?: github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_core_v1_placement_pb.TemplateMetadata.AsObject,
    }
  }
}

export class FederatedMatchableHttpGatewayStatus extends jspb.Message {
  hasPlacementStatus(): boolean;
  clearPlacementStatus(): void;
  getPlacementStatus(): github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_core_v1_placement_pb.PlacementStatus | undefined;
  setPlacementStatus(value?: github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_core_v1_placement_pb.PlacementStatus): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FederatedMatchableHttpGatewayStatus.AsObject;
  static toObject(includeInstance: boolean, msg: FederatedMatchableHttpGatewayStatus): FederatedMatchableHttpGatewayStatus.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: FederatedMatchableHttpGatewayStatus, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FederatedMatchableHttpGatewayStatus;
  static deserializeBinaryFromReader(message: FederatedMatchableHttpGatewayStatus, reader: jspb.BinaryReader): FederatedMatchableHttpGatewayStatus;
}

export namespace FederatedMatchableHttpGatewayStatus {
  export type AsObject = {
    placementStatus?: github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_core_v1_placement_pb.PlacementStatus.AsObject,
  }
}
