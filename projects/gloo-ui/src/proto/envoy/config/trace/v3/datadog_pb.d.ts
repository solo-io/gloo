/* eslint-disable */
// package: envoy.config.trace.v3
// file: envoy/config/trace/v3/datadog.proto

import * as jspb from "google-protobuf";
import * as udpa_annotations_migrate_pb from "../../../../udpa/annotations/migrate_pb";
import * as udpa_annotations_status_pb from "../../../../udpa/annotations/status_pb";
import * as udpa_annotations_versioning_pb from "../../../../udpa/annotations/versioning_pb";
import * as validate_validate_pb from "../../../../validate/validate_pb";
import * as solo_kit_api_v1_ref_pb from "../../../../solo-kit/api/v1/ref_pb";
import * as gogoproto_gogo_pb from "../../../../gogoproto/gogo_pb";

export class DatadogConfig extends jspb.Message {
  hasCollectorUpstreamRef(): boolean;
  clearCollectorUpstreamRef(): void;
  getCollectorUpstreamRef(): solo_kit_api_v1_ref_pb.ResourceRef | undefined;
  setCollectorUpstreamRef(value?: solo_kit_api_v1_ref_pb.ResourceRef): void;

  getServiceName(): string;
  setServiceName(value: string): void;

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
    collectorUpstreamRef?: solo_kit_api_v1_ref_pb.ResourceRef.AsObject,
    serviceName: string,
  }
}
