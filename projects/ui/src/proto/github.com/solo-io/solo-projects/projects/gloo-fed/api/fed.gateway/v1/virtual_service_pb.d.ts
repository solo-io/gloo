/* eslint-disable */
// package: fed.gateway.solo.io
// file: github.com/solo-io/solo-projects/projects/gloo-fed/api/fed.gateway/v1/virtual_service.proto

import * as jspb from "google-protobuf";
import * as extproto_ext_pb from "../../../../../../../../extproto/ext_pb";
import * as github_com_solo_io_solo_projects_projects_gloo_fed_api_multicluster_v1alpha1_multicluster_pb from "../../../../../../../../github.com/solo-io/solo-projects/projects/gloo-fed/api/multicluster/v1alpha1/multicluster_pb";
import * as github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_core_v1_placement_pb from "../../../../../../../../github.com/solo-io/solo-projects/projects/gloo-fed/api/fed/core/v1/placement_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gateway_v1_virtual_service_pb from "../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gateway/v1/virtual_service_pb";

export class FederatedVirtualServiceSpec extends jspb.Message {
  hasTemplate(): boolean;
  clearTemplate(): void;
  getTemplate(): FederatedVirtualServiceSpec.Template | undefined;
  setTemplate(value?: FederatedVirtualServiceSpec.Template): void;

  hasPlacement(): boolean;
  clearPlacement(): void;
  getPlacement(): github_com_solo_io_solo_projects_projects_gloo_fed_api_multicluster_v1alpha1_multicluster_pb.Placement | undefined;
  setPlacement(value?: github_com_solo_io_solo_projects_projects_gloo_fed_api_multicluster_v1alpha1_multicluster_pb.Placement): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FederatedVirtualServiceSpec.AsObject;
  static toObject(includeInstance: boolean, msg: FederatedVirtualServiceSpec): FederatedVirtualServiceSpec.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: FederatedVirtualServiceSpec, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FederatedVirtualServiceSpec;
  static deserializeBinaryFromReader(message: FederatedVirtualServiceSpec, reader: jspb.BinaryReader): FederatedVirtualServiceSpec;
}

export namespace FederatedVirtualServiceSpec {
  export type AsObject = {
    template?: FederatedVirtualServiceSpec.Template.AsObject,
    placement?: github_com_solo_io_solo_projects_projects_gloo_fed_api_multicluster_v1alpha1_multicluster_pb.Placement.AsObject,
  }

  export class Template extends jspb.Message {
    hasSpec(): boolean;
    clearSpec(): void;
    getSpec(): github_com_solo_io_solo_apis_api_gloo_gateway_v1_virtual_service_pb.VirtualServiceSpec | undefined;
    setSpec(value?: github_com_solo_io_solo_apis_api_gloo_gateway_v1_virtual_service_pb.VirtualServiceSpec): void;

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
      spec?: github_com_solo_io_solo_apis_api_gloo_gateway_v1_virtual_service_pb.VirtualServiceSpec.AsObject,
      metadata?: github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_core_v1_placement_pb.TemplateMetadata.AsObject,
    }
  }
}

export class FederatedVirtualServiceStatus extends jspb.Message {
  hasPlacementStatus(): boolean;
  clearPlacementStatus(): void;
  getPlacementStatus(): github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_core_v1_placement_pb.PlacementStatus | undefined;
  setPlacementStatus(value?: github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_core_v1_placement_pb.PlacementStatus): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FederatedVirtualServiceStatus.AsObject;
  static toObject(includeInstance: boolean, msg: FederatedVirtualServiceStatus): FederatedVirtualServiceStatus.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: FederatedVirtualServiceStatus, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FederatedVirtualServiceStatus;
  static deserializeBinaryFromReader(message: FederatedVirtualServiceStatus, reader: jspb.BinaryReader): FederatedVirtualServiceStatus;
}

export namespace FederatedVirtualServiceStatus {
  export type AsObject = {
    placementStatus?: github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_core_v1_placement_pb.PlacementStatus.AsObject,
  }
}
