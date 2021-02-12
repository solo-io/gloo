/* eslint-disable */
// package: solo.io.envoy.config.trace.v3
// file: github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/config/trace/v3/zipkin.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";
import * as envoy_annotations_deprecation_pb from "../../../../../../../../../../../envoy/annotations/deprecation_pb";
import * as udpa_annotations_migrate_pb from "../../../../../../../../../../../udpa/annotations/migrate_pb";
import * as udpa_annotations_status_pb from "../../../../../../../../../../../udpa/annotations/status_pb";
import * as udpa_annotations_versioning_pb from "../../../../../../../../../../../udpa/annotations/versioning_pb";
import * as validate_validate_pb from "../../../../../../../../../../../validate/validate_pb";
import * as github_com_solo_io_solo_kit_api_v1_ref_pb from "../../../../../../../../../../../github.com/solo-io/solo-kit/api/v1/ref_pb";
import * as extproto_ext_pb from "../../../../../../../../../../../extproto/ext_pb";

export class ZipkinConfig extends jspb.Message {
  hasCollectorUpstreamRef(): boolean;
  clearCollectorUpstreamRef(): void;
  getCollectorUpstreamRef(): github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef | undefined;
  setCollectorUpstreamRef(value?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef): void;

  hasClusterName(): boolean;
  clearClusterName(): void;
  getClusterName(): string;
  setClusterName(value: string): void;

  getCollectorEndpoint(): string;
  setCollectorEndpoint(value: string): void;

  getTraceId128bit(): boolean;
  setTraceId128bit(value: boolean): void;

  hasSharedSpanContext(): boolean;
  clearSharedSpanContext(): void;
  getSharedSpanContext(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setSharedSpanContext(value?: google_protobuf_wrappers_pb.BoolValue): void;

  getCollectorEndpointVersion(): ZipkinConfig.CollectorEndpointVersionMap[keyof ZipkinConfig.CollectorEndpointVersionMap];
  setCollectorEndpointVersion(value: ZipkinConfig.CollectorEndpointVersionMap[keyof ZipkinConfig.CollectorEndpointVersionMap]): void;

  getCollectorClusterCase(): ZipkinConfig.CollectorClusterCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ZipkinConfig.AsObject;
  static toObject(includeInstance: boolean, msg: ZipkinConfig): ZipkinConfig.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ZipkinConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ZipkinConfig;
  static deserializeBinaryFromReader(message: ZipkinConfig, reader: jspb.BinaryReader): ZipkinConfig;
}

export namespace ZipkinConfig {
  export type AsObject = {
    collectorUpstreamRef?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject,
    clusterName: string,
    collectorEndpoint: string,
    traceId128bit: boolean,
    sharedSpanContext?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    collectorEndpointVersion: ZipkinConfig.CollectorEndpointVersionMap[keyof ZipkinConfig.CollectorEndpointVersionMap],
  }

  export interface CollectorEndpointVersionMap {
    DEPRECATED_AND_UNAVAILABLE_DO_NOT_USE: 0;
    HTTP_JSON: 1;
    HTTP_PROTO: 2;
  }

  export const CollectorEndpointVersion: CollectorEndpointVersionMap;

  export enum CollectorClusterCase {
    COLLECTOR_CLUSTER_NOT_SET = 0,
    COLLECTOR_UPSTREAM_REF = 1,
    CLUSTER_NAME = 6,
  }
}
