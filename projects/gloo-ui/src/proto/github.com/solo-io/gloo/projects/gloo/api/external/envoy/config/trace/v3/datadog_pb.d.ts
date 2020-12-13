/* eslint-disable */
// package: solo.io.envoy.config.trace.v3
// file: github.com/solo-io/gloo/projects/gloo/api/external/envoy/config/trace/v3/datadog.proto

import * as jspb from "google-protobuf";
import * as udpa_annotations_migrate_pb from "../../../../../../../../../../../udpa/annotations/migrate_pb";
import * as udpa_annotations_status_pb from "../../../../../../../../../../../udpa/annotations/status_pb";
import * as udpa_annotations_versioning_pb from "../../../../../../../../../../../udpa/annotations/versioning_pb";
import * as validate_validate_pb from "../../../../../../../../../../../validate/validate_pb";
import * as github_com_solo_io_solo_kit_api_v1_ref_pb from "../../../../../../../../../../../github.com/solo-io/solo-kit/api/v1/ref_pb";
import * as extproto_ext_pb from "../../../../../../../../../../../protoc-gen-ext/extproto/ext_pb";

export class DatadogConfig extends jspb.Message {
  hasCollectorUpstreamRef(): boolean;
  clearCollectorUpstreamRef(): void;
  getCollectorUpstreamRef(): github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef | undefined;
  setCollectorUpstreamRef(value?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef): void;

  hasClusterName(): boolean;
  clearClusterName(): void;
  getClusterName(): string;
  setClusterName(value: string): void;

  getServiceName(): string;
  setServiceName(value: string): void;

  getCollectorClusterCase(): DatadogConfig.CollectorClusterCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DatadogConfig.AsObject;
  static toObject(includeInstance: boolean, msg: DatadogConfig): DatadogConfig.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DatadogConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DatadogConfig;
  static deserializeBinaryFromReader(message: DatadogConfig, reader: jspb.BinaryReader): DatadogConfig;
}

export namespace DatadogConfig {
  export type AsObject = {
    collectorUpstreamRef?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject,
    clusterName: string,
    serviceName: string,
  }

  export enum CollectorClusterCase {
    COLLECTOR_CLUSTER_NOT_SET = 0,
    COLLECTOR_UPSTREAM_REF = 1,
    CLUSTER_NAME = 3,
  }
}
