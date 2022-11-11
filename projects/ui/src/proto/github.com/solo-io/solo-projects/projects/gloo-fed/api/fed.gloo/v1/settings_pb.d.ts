/* eslint-disable */
// package: fed.gloo.solo.io
// file: github.com/solo-io/solo-projects/projects/gloo-fed/api/fed.gloo/v1/settings.proto

import * as jspb from "google-protobuf";
import * as extproto_ext_pb from "../../../../../../../../extproto/ext_pb";
import * as github_com_solo_io_solo_projects_projects_gloo_fed_api_multicluster_v1alpha1_multicluster_pb from "../../../../../../../../github.com/solo-io/solo-projects/projects/gloo-fed/api/multicluster/v1alpha1/multicluster_pb";
import * as github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_core_v1_placement_pb from "../../../../../../../../github.com/solo-io/solo-projects/projects/gloo-fed/api/fed/core/v1/placement_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_settings_pb from "../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/settings_pb";

export class FederatedSettingsSpec extends jspb.Message {
  hasTemplate(): boolean;
  clearTemplate(): void;
  getTemplate(): FederatedSettingsSpec.Template | undefined;
  setTemplate(value?: FederatedSettingsSpec.Template): void;

  hasPlacement(): boolean;
  clearPlacement(): void;
  getPlacement(): github_com_solo_io_solo_projects_projects_gloo_fed_api_multicluster_v1alpha1_multicluster_pb.Placement | undefined;
  setPlacement(value?: github_com_solo_io_solo_projects_projects_gloo_fed_api_multicluster_v1alpha1_multicluster_pb.Placement): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FederatedSettingsSpec.AsObject;
  static toObject(includeInstance: boolean, msg: FederatedSettingsSpec): FederatedSettingsSpec.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: FederatedSettingsSpec, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FederatedSettingsSpec;
  static deserializeBinaryFromReader(message: FederatedSettingsSpec, reader: jspb.BinaryReader): FederatedSettingsSpec;
}

export namespace FederatedSettingsSpec {
  export type AsObject = {
    template?: FederatedSettingsSpec.Template.AsObject,
    placement?: github_com_solo_io_solo_projects_projects_gloo_fed_api_multicluster_v1alpha1_multicluster_pb.Placement.AsObject,
  }

  export class Template extends jspb.Message {
    hasSpec(): boolean;
    clearSpec(): void;
    getSpec(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_settings_pb.SettingsSpec | undefined;
    setSpec(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_settings_pb.SettingsSpec): void;

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
      spec?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_settings_pb.SettingsSpec.AsObject,
      metadata?: github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_core_v1_placement_pb.TemplateMetadata.AsObject,
    }
  }
}

export class FederatedSettingsStatus extends jspb.Message {
  hasPlacementStatus(): boolean;
  clearPlacementStatus(): void;
  getPlacementStatus(): github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_core_v1_placement_pb.PlacementStatus | undefined;
  setPlacementStatus(value?: github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_core_v1_placement_pb.PlacementStatus): void;

  getNamespacedPlacementStatusesMap(): jspb.Map<string, github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_core_v1_placement_pb.PlacementStatus>;
  clearNamespacedPlacementStatusesMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FederatedSettingsStatus.AsObject;
  static toObject(includeInstance: boolean, msg: FederatedSettingsStatus): FederatedSettingsStatus.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: FederatedSettingsStatus, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FederatedSettingsStatus;
  static deserializeBinaryFromReader(message: FederatedSettingsStatus, reader: jspb.BinaryReader): FederatedSettingsStatus;
}

export namespace FederatedSettingsStatus {
  export type AsObject = {
    placementStatus?: github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_core_v1_placement_pb.PlacementStatus.AsObject,
    namespacedPlacementStatusesMap: Array<[string, github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_core_v1_placement_pb.PlacementStatus.AsObject]>,
  }
}
