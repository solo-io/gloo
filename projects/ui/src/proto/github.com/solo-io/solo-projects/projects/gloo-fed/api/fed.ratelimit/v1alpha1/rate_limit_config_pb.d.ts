/* eslint-disable */
// package: fed.ratelimit.solo.io
// file: github.com/solo-io/solo-projects/projects/gloo-fed/api/fed.ratelimit/v1alpha1/rate_limit_config.proto

import * as jspb from "google-protobuf";
import * as extproto_ext_pb from "../../../../../../../../extproto/ext_pb";
import * as github_com_solo_io_skv2_enterprise_multicluster_admission_webhook_api_multicluster_v1alpha1_multicluster_pb from "../../../../../../../../github.com/solo-io/skv2-enterprise/multicluster-admission-webhook/api/multicluster/v1alpha1/multicluster_pb";
import * as github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_core_v1_placement_pb from "../../../../../../../../github.com/solo-io/solo-projects/projects/gloo-fed/api/fed/core/v1/placement_pb";
import * as github_com_solo_io_solo_apis_api_rate_limiter_v1alpha1_ratelimit_pb from "../../../../../../../../github.com/solo-io/solo-apis/api/rate-limiter/v1alpha1/ratelimit_pb";

export class FederatedRateLimitConfigSpec extends jspb.Message {
  hasTemplate(): boolean;
  clearTemplate(): void;
  getTemplate(): FederatedRateLimitConfigSpec.Template | undefined;
  setTemplate(value?: FederatedRateLimitConfigSpec.Template): void;

  hasPlacement(): boolean;
  clearPlacement(): void;
  getPlacement(): github_com_solo_io_skv2_enterprise_multicluster_admission_webhook_api_multicluster_v1alpha1_multicluster_pb.Placement | undefined;
  setPlacement(value?: github_com_solo_io_skv2_enterprise_multicluster_admission_webhook_api_multicluster_v1alpha1_multicluster_pb.Placement): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FederatedRateLimitConfigSpec.AsObject;
  static toObject(includeInstance: boolean, msg: FederatedRateLimitConfigSpec): FederatedRateLimitConfigSpec.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: FederatedRateLimitConfigSpec, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FederatedRateLimitConfigSpec;
  static deserializeBinaryFromReader(message: FederatedRateLimitConfigSpec, reader: jspb.BinaryReader): FederatedRateLimitConfigSpec;
}

export namespace FederatedRateLimitConfigSpec {
  export type AsObject = {
    template?: FederatedRateLimitConfigSpec.Template.AsObject,
    placement?: github_com_solo_io_skv2_enterprise_multicluster_admission_webhook_api_multicluster_v1alpha1_multicluster_pb.Placement.AsObject,
  }

  export class Template extends jspb.Message {
    hasSpec(): boolean;
    clearSpec(): void;
    getSpec(): github_com_solo_io_solo_apis_api_rate_limiter_v1alpha1_ratelimit_pb.RateLimitConfigSpec | undefined;
    setSpec(value?: github_com_solo_io_solo_apis_api_rate_limiter_v1alpha1_ratelimit_pb.RateLimitConfigSpec): void;

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
      spec?: github_com_solo_io_solo_apis_api_rate_limiter_v1alpha1_ratelimit_pb.RateLimitConfigSpec.AsObject,
      metadata?: github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_core_v1_placement_pb.TemplateMetadata.AsObject,
    }
  }
}

export class FederatedRateLimitConfigStatus extends jspb.Message {
  hasPlacementStatus(): boolean;
  clearPlacementStatus(): void;
  getPlacementStatus(): github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_core_v1_placement_pb.PlacementStatus | undefined;
  setPlacementStatus(value?: github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_core_v1_placement_pb.PlacementStatus): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FederatedRateLimitConfigStatus.AsObject;
  static toObject(includeInstance: boolean, msg: FederatedRateLimitConfigStatus): FederatedRateLimitConfigStatus.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: FederatedRateLimitConfigStatus, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FederatedRateLimitConfigStatus;
  static deserializeBinaryFromReader(message: FederatedRateLimitConfigStatus, reader: jspb.BinaryReader): FederatedRateLimitConfigStatus;
}

export namespace FederatedRateLimitConfigStatus {
  export type AsObject = {
    placementStatus?: github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_core_v1_placement_pb.PlacementStatus.AsObject,
  }
}
