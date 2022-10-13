/* eslint-disable */
// package: fed.enterprise.gloo.solo.io
// file: github.com/solo-io/solo-projects/projects/gloo-fed/api/fed.enterprise.gloo/v1/auth_config.proto

import * as jspb from "google-protobuf";
import * as extproto_ext_pb from "../../../../../../../../extproto/ext_pb";
import * as github_com_solo_io_solo_projects_projects_gloo_fed_api_multicluster_v1alpha1_multicluster_pb from "../../../../../../../../github.com/solo-io/solo-projects/projects/gloo-fed/api/multicluster/v1alpha1/multicluster_pb";
import * as github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_core_v1_placement_pb from "../../../../../../../../github.com/solo-io/solo-projects/projects/gloo-fed/api/fed/core/v1/placement_pb";
import * as github_com_solo_io_solo_apis_api_gloo_enterprise_gloo_v1_auth_config_pb from "../../../../../../../../github.com/solo-io/solo-apis/api/gloo/enterprise.gloo/v1/auth_config_pb";

export class FederatedAuthConfigSpec extends jspb.Message {
  hasTemplate(): boolean;
  clearTemplate(): void;
  getTemplate(): FederatedAuthConfigSpec.Template | undefined;
  setTemplate(value?: FederatedAuthConfigSpec.Template): void;

  hasPlacement(): boolean;
  clearPlacement(): void;
  getPlacement(): github_com_solo_io_solo_projects_projects_gloo_fed_api_multicluster_v1alpha1_multicluster_pb.Placement | undefined;
  setPlacement(value?: github_com_solo_io_solo_projects_projects_gloo_fed_api_multicluster_v1alpha1_multicluster_pb.Placement): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FederatedAuthConfigSpec.AsObject;
  static toObject(includeInstance: boolean, msg: FederatedAuthConfigSpec): FederatedAuthConfigSpec.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: FederatedAuthConfigSpec, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FederatedAuthConfigSpec;
  static deserializeBinaryFromReader(message: FederatedAuthConfigSpec, reader: jspb.BinaryReader): FederatedAuthConfigSpec;
}

export namespace FederatedAuthConfigSpec {
  export type AsObject = {
    template?: FederatedAuthConfigSpec.Template.AsObject,
    placement?: github_com_solo_io_solo_projects_projects_gloo_fed_api_multicluster_v1alpha1_multicluster_pb.Placement.AsObject,
  }

  export class Template extends jspb.Message {
    hasSpec(): boolean;
    clearSpec(): void;
    getSpec(): github_com_solo_io_solo_apis_api_gloo_enterprise_gloo_v1_auth_config_pb.AuthConfigSpec | undefined;
    setSpec(value?: github_com_solo_io_solo_apis_api_gloo_enterprise_gloo_v1_auth_config_pb.AuthConfigSpec): void;

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
      spec?: github_com_solo_io_solo_apis_api_gloo_enterprise_gloo_v1_auth_config_pb.AuthConfigSpec.AsObject,
      metadata?: github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_core_v1_placement_pb.TemplateMetadata.AsObject,
    }
  }
}

export class FederatedAuthConfigStatus extends jspb.Message {
  hasPlacementStatus(): boolean;
  clearPlacementStatus(): void;
  getPlacementStatus(): github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_core_v1_placement_pb.PlacementStatus | undefined;
  setPlacementStatus(value?: github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_core_v1_placement_pb.PlacementStatus): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FederatedAuthConfigStatus.AsObject;
  static toObject(includeInstance: boolean, msg: FederatedAuthConfigStatus): FederatedAuthConfigStatus.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: FederatedAuthConfigStatus, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FederatedAuthConfigStatus;
  static deserializeBinaryFromReader(message: FederatedAuthConfigStatus, reader: jspb.BinaryReader): FederatedAuthConfigStatus;
}

export namespace FederatedAuthConfigStatus {
  export type AsObject = {
    placementStatus?: github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_core_v1_placement_pb.PlacementStatus.AsObject,
  }
}
