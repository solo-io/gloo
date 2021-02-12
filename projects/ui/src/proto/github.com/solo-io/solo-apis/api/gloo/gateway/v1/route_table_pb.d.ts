/* eslint-disable */
// package: gateway.solo.io
// file: github.com/solo-io/solo-apis/api/gloo/gateway/v1/route_table.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_struct_pb from "google-protobuf/google/protobuf/struct_pb";
import * as extproto_ext_pb from "../../../../../../../extproto/ext_pb";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";
import * as github_com_solo_io_solo_kit_api_v1_solo_kit_pb from "../../../../../../../github.com/solo-io/solo-kit/api/v1/solo-kit_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gateway_v1_virtual_service_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gateway/v1/virtual_service_pb";

export class RouteTableSpec extends jspb.Message {
  clearRoutesList(): void;
  getRoutesList(): Array<github_com_solo_io_solo_apis_api_gloo_gateway_v1_virtual_service_pb.Route>;
  setRoutesList(value: Array<github_com_solo_io_solo_apis_api_gloo_gateway_v1_virtual_service_pb.Route>): void;
  addRoutes(value?: github_com_solo_io_solo_apis_api_gloo_gateway_v1_virtual_service_pb.Route, index?: number): github_com_solo_io_solo_apis_api_gloo_gateway_v1_virtual_service_pb.Route;

  hasWeight(): boolean;
  clearWeight(): void;
  getWeight(): google_protobuf_wrappers_pb.Int32Value | undefined;
  setWeight(value?: google_protobuf_wrappers_pb.Int32Value): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RouteTableSpec.AsObject;
  static toObject(includeInstance: boolean, msg: RouteTableSpec): RouteTableSpec.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RouteTableSpec, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RouteTableSpec;
  static deserializeBinaryFromReader(message: RouteTableSpec, reader: jspb.BinaryReader): RouteTableSpec;
}

export namespace RouteTableSpec {
  export type AsObject = {
    routesList: Array<github_com_solo_io_solo_apis_api_gloo_gateway_v1_virtual_service_pb.Route.AsObject>,
    weight?: google_protobuf_wrappers_pb.Int32Value.AsObject,
  }
}

export class RouteTableStatus extends jspb.Message {
  getState(): RouteTableStatus.StateMap[keyof RouteTableStatus.StateMap];
  setState(value: RouteTableStatus.StateMap[keyof RouteTableStatus.StateMap]): void;

  getReason(): string;
  setReason(value: string): void;

  getReportedBy(): string;
  setReportedBy(value: string): void;

  getSubresourceStatusesMap(): jspb.Map<string, RouteTableStatus>;
  clearSubresourceStatusesMap(): void;
  hasDetails(): boolean;
  clearDetails(): void;
  getDetails(): google_protobuf_struct_pb.Struct | undefined;
  setDetails(value?: google_protobuf_struct_pb.Struct): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RouteTableStatus.AsObject;
  static toObject(includeInstance: boolean, msg: RouteTableStatus): RouteTableStatus.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RouteTableStatus, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RouteTableStatus;
  static deserializeBinaryFromReader(message: RouteTableStatus, reader: jspb.BinaryReader): RouteTableStatus;
}

export namespace RouteTableStatus {
  export type AsObject = {
    state: RouteTableStatus.StateMap[keyof RouteTableStatus.StateMap],
    reason: string,
    reportedBy: string,
    subresourceStatusesMap: Array<[string, RouteTableStatus.AsObject]>,
    details?: google_protobuf_struct_pb.Struct.AsObject,
  }

  export interface StateMap {
    PENDING: 0;
    ACCEPTED: 1;
    REJECTED: 2;
    WARNING: 3;
  }

  export const State: StateMap;
}
