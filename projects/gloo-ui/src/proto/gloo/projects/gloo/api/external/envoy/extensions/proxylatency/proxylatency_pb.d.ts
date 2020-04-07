/* eslint-disable */
// package: envoy.config.filter.http.proxylatency.v2
// file: gloo/projects/gloo/api/external/envoy/extensions/proxylatency/proxylatency.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";
import * as gogoproto_gogo_pb from "../../../../../../../../gogoproto/gogo_pb";

export class ProxyLatency extends jspb.Message {
  getRequest(): ProxyLatency.MeasurementMap[keyof ProxyLatency.MeasurementMap];
  setRequest(value: ProxyLatency.MeasurementMap[keyof ProxyLatency.MeasurementMap]): void;

  getResponse(): ProxyLatency.MeasurementMap[keyof ProxyLatency.MeasurementMap];
  setResponse(value: ProxyLatency.MeasurementMap[keyof ProxyLatency.MeasurementMap]): void;

  hasChargeClusterStat(): boolean;
  clearChargeClusterStat(): void;
  getChargeClusterStat(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setChargeClusterStat(value?: google_protobuf_wrappers_pb.BoolValue): void;

  hasChargeListenerStat(): boolean;
  clearChargeListenerStat(): void;
  getChargeListenerStat(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setChargeListenerStat(value?: google_protobuf_wrappers_pb.BoolValue): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProxyLatency.AsObject;
  static toObject(includeInstance: boolean, msg: ProxyLatency): ProxyLatency.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ProxyLatency, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProxyLatency;
  static deserializeBinaryFromReader(message: ProxyLatency, reader: jspb.BinaryReader): ProxyLatency;
}

export namespace ProxyLatency {
  export type AsObject = {
    request: ProxyLatency.MeasurementMap[keyof ProxyLatency.MeasurementMap],
    response: ProxyLatency.MeasurementMap[keyof ProxyLatency.MeasurementMap],
    chargeClusterStat?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    chargeListenerStat?: google_protobuf_wrappers_pb.BoolValue.AsObject,
  }

  export interface MeasurementMap {
    LAST_INCOMING_FIRST_OUTGOING: 0;
    FIRST_INCOMING_FIRST_OUTGOING: 1;
    LAST_INCOMING_LAST_OUTGOING: 2;
    FIRST_INCOMING_LAST_OUTGOING: 3;
  }

  export const Measurement: MeasurementMap;
}
