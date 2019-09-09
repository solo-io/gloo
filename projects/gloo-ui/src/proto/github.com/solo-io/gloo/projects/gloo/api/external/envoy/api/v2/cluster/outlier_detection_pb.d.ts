// package: envoy.api.v2.cluster
// file: github.com/solo-io/gloo/projects/gloo/api/external/envoy/api/v2/cluster/outlier_detection.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_duration_pb from "google-protobuf/google/protobuf/duration_pb";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";
import * as validate_validate_pb from "../../../../../../../../../../../validate/validate_pb";
import * as gogoproto_gogo_pb from "../../../../../../../../../../../gogoproto/gogo_pb";

export class OutlierDetection extends jspb.Message {
  hasConsecutive5xx(): boolean;
  clearConsecutive5xx(): void;
  getConsecutive5xx(): google_protobuf_wrappers_pb.UInt32Value | undefined;
  setConsecutive5xx(value?: google_protobuf_wrappers_pb.UInt32Value): void;

  hasInterval(): boolean;
  clearInterval(): void;
  getInterval(): google_protobuf_duration_pb.Duration | undefined;
  setInterval(value?: google_protobuf_duration_pb.Duration): void;

  hasBaseEjectionTime(): boolean;
  clearBaseEjectionTime(): void;
  getBaseEjectionTime(): google_protobuf_duration_pb.Duration | undefined;
  setBaseEjectionTime(value?: google_protobuf_duration_pb.Duration): void;

  hasMaxEjectionPercent(): boolean;
  clearMaxEjectionPercent(): void;
  getMaxEjectionPercent(): google_protobuf_wrappers_pb.UInt32Value | undefined;
  setMaxEjectionPercent(value?: google_protobuf_wrappers_pb.UInt32Value): void;

  hasEnforcingConsecutive5xx(): boolean;
  clearEnforcingConsecutive5xx(): void;
  getEnforcingConsecutive5xx(): google_protobuf_wrappers_pb.UInt32Value | undefined;
  setEnforcingConsecutive5xx(value?: google_protobuf_wrappers_pb.UInt32Value): void;

  hasEnforcingSuccessRate(): boolean;
  clearEnforcingSuccessRate(): void;
  getEnforcingSuccessRate(): google_protobuf_wrappers_pb.UInt32Value | undefined;
  setEnforcingSuccessRate(value?: google_protobuf_wrappers_pb.UInt32Value): void;

  hasSuccessRateMinimumHosts(): boolean;
  clearSuccessRateMinimumHosts(): void;
  getSuccessRateMinimumHosts(): google_protobuf_wrappers_pb.UInt32Value | undefined;
  setSuccessRateMinimumHosts(value?: google_protobuf_wrappers_pb.UInt32Value): void;

  hasSuccessRateRequestVolume(): boolean;
  clearSuccessRateRequestVolume(): void;
  getSuccessRateRequestVolume(): google_protobuf_wrappers_pb.UInt32Value | undefined;
  setSuccessRateRequestVolume(value?: google_protobuf_wrappers_pb.UInt32Value): void;

  hasSuccessRateStdevFactor(): boolean;
  clearSuccessRateStdevFactor(): void;
  getSuccessRateStdevFactor(): google_protobuf_wrappers_pb.UInt32Value | undefined;
  setSuccessRateStdevFactor(value?: google_protobuf_wrappers_pb.UInt32Value): void;

  hasConsecutiveGatewayFailure(): boolean;
  clearConsecutiveGatewayFailure(): void;
  getConsecutiveGatewayFailure(): google_protobuf_wrappers_pb.UInt32Value | undefined;
  setConsecutiveGatewayFailure(value?: google_protobuf_wrappers_pb.UInt32Value): void;

  hasEnforcingConsecutiveGatewayFailure(): boolean;
  clearEnforcingConsecutiveGatewayFailure(): void;
  getEnforcingConsecutiveGatewayFailure(): google_protobuf_wrappers_pb.UInt32Value | undefined;
  setEnforcingConsecutiveGatewayFailure(value?: google_protobuf_wrappers_pb.UInt32Value): void;

  getSplitExternalLocalOriginErrors(): boolean;
  setSplitExternalLocalOriginErrors(value: boolean): void;

  hasConsecutiveLocalOriginFailure(): boolean;
  clearConsecutiveLocalOriginFailure(): void;
  getConsecutiveLocalOriginFailure(): google_protobuf_wrappers_pb.UInt32Value | undefined;
  setConsecutiveLocalOriginFailure(value?: google_protobuf_wrappers_pb.UInt32Value): void;

  hasEnforcingConsecutiveLocalOriginFailure(): boolean;
  clearEnforcingConsecutiveLocalOriginFailure(): void;
  getEnforcingConsecutiveLocalOriginFailure(): google_protobuf_wrappers_pb.UInt32Value | undefined;
  setEnforcingConsecutiveLocalOriginFailure(value?: google_protobuf_wrappers_pb.UInt32Value): void;

  hasEnforcingLocalOriginSuccessRate(): boolean;
  clearEnforcingLocalOriginSuccessRate(): void;
  getEnforcingLocalOriginSuccessRate(): google_protobuf_wrappers_pb.UInt32Value | undefined;
  setEnforcingLocalOriginSuccessRate(value?: google_protobuf_wrappers_pb.UInt32Value): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OutlierDetection.AsObject;
  static toObject(includeInstance: boolean, msg: OutlierDetection): OutlierDetection.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: OutlierDetection, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OutlierDetection;
  static deserializeBinaryFromReader(message: OutlierDetection, reader: jspb.BinaryReader): OutlierDetection;
}

export namespace OutlierDetection {
  export type AsObject = {
    consecutive5xx?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
    interval?: google_protobuf_duration_pb.Duration.AsObject,
    baseEjectionTime?: google_protobuf_duration_pb.Duration.AsObject,
    maxEjectionPercent?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
    enforcingConsecutive5xx?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
    enforcingSuccessRate?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
    successRateMinimumHosts?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
    successRateRequestVolume?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
    successRateStdevFactor?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
    consecutiveGatewayFailure?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
    enforcingConsecutiveGatewayFailure?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
    splitExternalLocalOriginErrors: boolean,
    consecutiveLocalOriginFailure?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
    enforcingConsecutiveLocalOriginFailure?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
    enforcingLocalOriginSuccessRate?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
  }
}

