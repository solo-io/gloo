/* eslint-disable */
// package: fed.rpc.solo.io
// file: github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/wasm.proto

import * as jspb from "google-protobuf";
import * as extproto_ext_pb from "../../../../../../../../extproto/ext_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_core_matchers_matchers_pb from "../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/core/matchers/matchers_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_wasm_wasm_pb from "../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/options/wasm/wasm_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gateway_v1_gateway_pb from "../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gateway/v1/gateway_pb";
import * as github_com_solo_io_solo_projects_projects_gloo_fed_api_fed_v1_instance_pb from "../../../../../../../../github.com/solo-io/solo-projects/projects/gloo-fed/api/fed/v1/instance_pb";
import * as github_com_solo_io_skv2_api_core_v1_core_pb from "../../../../../../../../github.com/solo-io/skv2/api/core/v1/core_pb";

export class WasmFilter extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getRootId(): string;
  setRootId(value: string): void;

  getSource(): string;
  setSource(value: string): void;

  getConfig(): string;
  setConfig(value: string): void;

  clearLocationsList(): void;
  getLocationsList(): Array<WasmFilter.Location>;
  setLocationsList(value: Array<WasmFilter.Location>): void;
  addLocations(value?: WasmFilter.Location, index?: number): WasmFilter.Location;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): WasmFilter.AsObject;
  static toObject(includeInstance: boolean, msg: WasmFilter): WasmFilter.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: WasmFilter, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): WasmFilter;
  static deserializeBinaryFromReader(message: WasmFilter, reader: jspb.BinaryReader): WasmFilter;
}

export namespace WasmFilter {
  export type AsObject = {
    name: string,
    rootId: string,
    source: string,
    config: string,
    locationsList: Array<WasmFilter.Location.AsObject>,
  }

  export class Location extends jspb.Message {
    hasGatewayRef(): boolean;
    clearGatewayRef(): void;
    getGatewayRef(): github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef | undefined;
    setGatewayRef(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef): void;

    hasGatewayStatus(): boolean;
    clearGatewayStatus(): void;
    getGatewayStatus(): github_com_solo_io_solo_apis_api_gloo_gateway_v1_gateway_pb.GatewayStatus | undefined;
    setGatewayStatus(value?: github_com_solo_io_solo_apis_api_gloo_gateway_v1_gateway_pb.GatewayStatus): void;

    hasGlooInstanceRef(): boolean;
    clearGlooInstanceRef(): void;
    getGlooInstanceRef(): github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef | undefined;
    setGlooInstanceRef(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Location.AsObject;
    static toObject(includeInstance: boolean, msg: Location): Location.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Location, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Location;
    static deserializeBinaryFromReader(message: Location, reader: jspb.BinaryReader): Location;
  }

  export namespace Location {
    export type AsObject = {
      gatewayRef?: github_com_solo_io_skv2_api_core_v1_core_pb.ClusterObjectRef.AsObject,
      gatewayStatus?: github_com_solo_io_solo_apis_api_gloo_gateway_v1_gateway_pb.GatewayStatus.AsObject,
      glooInstanceRef?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef.AsObject,
    }
  }
}

export class ListWasmFiltersRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListWasmFiltersRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListWasmFiltersRequest): ListWasmFiltersRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListWasmFiltersRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListWasmFiltersRequest;
  static deserializeBinaryFromReader(message: ListWasmFiltersRequest, reader: jspb.BinaryReader): ListWasmFiltersRequest;
}

export namespace ListWasmFiltersRequest {
  export type AsObject = {
  }
}

export class ListWasmFiltersResponse extends jspb.Message {
  clearWasmFiltersList(): void;
  getWasmFiltersList(): Array<WasmFilter>;
  setWasmFiltersList(value: Array<WasmFilter>): void;
  addWasmFilters(value?: WasmFilter, index?: number): WasmFilter;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListWasmFiltersResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListWasmFiltersResponse): ListWasmFiltersResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListWasmFiltersResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListWasmFiltersResponse;
  static deserializeBinaryFromReader(message: ListWasmFiltersResponse, reader: jspb.BinaryReader): ListWasmFiltersResponse;
}

export namespace ListWasmFiltersResponse {
  export type AsObject = {
    wasmFiltersList: Array<WasmFilter.AsObject>,
  }
}

export class DescribeWasmFilterRequest extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getRootId(): string;
  setRootId(value: string): void;

  hasGatewayRef(): boolean;
  clearGatewayRef(): void;
  getGatewayRef(): github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef | undefined;
  setGatewayRef(value?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DescribeWasmFilterRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DescribeWasmFilterRequest): DescribeWasmFilterRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DescribeWasmFilterRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DescribeWasmFilterRequest;
  static deserializeBinaryFromReader(message: DescribeWasmFilterRequest, reader: jspb.BinaryReader): DescribeWasmFilterRequest;
}

export namespace DescribeWasmFilterRequest {
  export type AsObject = {
    name: string,
    rootId: string,
    gatewayRef?: github_com_solo_io_skv2_api_core_v1_core_pb.ObjectRef.AsObject,
  }
}

export class DescribeWasmFilterResponse extends jspb.Message {
  hasWasmFilter(): boolean;
  clearWasmFilter(): void;
  getWasmFilter(): WasmFilter | undefined;
  setWasmFilter(value?: WasmFilter): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DescribeWasmFilterResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DescribeWasmFilterResponse): DescribeWasmFilterResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DescribeWasmFilterResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DescribeWasmFilterResponse;
  static deserializeBinaryFromReader(message: DescribeWasmFilterResponse, reader: jspb.BinaryReader): DescribeWasmFilterResponse;
}

export namespace DescribeWasmFilterResponse {
  export type AsObject = {
    wasmFilter?: WasmFilter.AsObject,
  }
}
