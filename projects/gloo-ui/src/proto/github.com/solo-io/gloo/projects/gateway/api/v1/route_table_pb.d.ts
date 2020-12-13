/* eslint-disable */
// package: gateway.solo.io
// file: github.com/solo-io/gloo/projects/gateway/api/v1/route_table.proto

import * as jspb from "google-protobuf";
import * as extproto_ext_pb from "../../../../../../../protoc-gen-ext/extproto/ext_pb";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";
import * as github_com_solo_io_solo_kit_api_v1_metadata_pb from "../../../../../../../github.com/solo-io/solo-kit/api/v1/metadata_pb";
import * as github_com_solo_io_solo_kit_api_v1_status_pb from "../../../../../../../github.com/solo-io/solo-kit/api/v1/status_pb";
import * as github_com_solo_io_solo_kit_api_v1_solo_kit_pb from "../../../../../../../github.com/solo-io/solo-kit/api/v1/solo-kit_pb";
import * as github_com_solo_io_gloo_projects_gateway_api_v1_virtual_service_pb from "../../../../../../../github.com/solo-io/gloo/projects/gateway/api/v1/virtual_service_pb";

export class RouteTable extends jspb.Message {
  clearRoutesList(): void;
  getRoutesList(): Array<github_com_solo_io_gloo_projects_gateway_api_v1_virtual_service_pb.Route>;
  setRoutesList(value: Array<github_com_solo_io_gloo_projects_gateway_api_v1_virtual_service_pb.Route>): void;
  addRoutes(value?: github_com_solo_io_gloo_projects_gateway_api_v1_virtual_service_pb.Route, index?: number): github_com_solo_io_gloo_projects_gateway_api_v1_virtual_service_pb.Route;

  hasWeight(): boolean;
  clearWeight(): void;
  getWeight(): google_protobuf_wrappers_pb.Int32Value | undefined;
  setWeight(value?: google_protobuf_wrappers_pb.Int32Value): void;

  hasStatus(): boolean;
  clearStatus(): void;
  getStatus(): github_com_solo_io_solo_kit_api_v1_status_pb.Status | undefined;
  setStatus(value?: github_com_solo_io_solo_kit_api_v1_status_pb.Status): void;

  hasMetadata(): boolean;
  clearMetadata(): void;
  getMetadata(): github_com_solo_io_solo_kit_api_v1_metadata_pb.Metadata | undefined;
  setMetadata(value?: github_com_solo_io_solo_kit_api_v1_metadata_pb.Metadata): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RouteTable.AsObject;
  static toObject(includeInstance: boolean, msg: RouteTable): RouteTable.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RouteTable, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RouteTable;
  static deserializeBinaryFromReader(message: RouteTable, reader: jspb.BinaryReader): RouteTable;
}

export namespace RouteTable {
  export type AsObject = {
    routesList: Array<github_com_solo_io_gloo_projects_gateway_api_v1_virtual_service_pb.Route.AsObject>,
    weight?: google_protobuf_wrappers_pb.Int32Value.AsObject,
    status?: github_com_solo_io_solo_kit_api_v1_status_pb.Status.AsObject,
    metadata?: github_com_solo_io_solo_kit_api_v1_metadata_pb.Metadata.AsObject,
  }
}
