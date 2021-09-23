/* eslint-disable */
// package: gateway.solo.io
// file: github.com/solo-io/solo-apis/api/gloo/gateway/v1/external_options.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_struct_pb from "google-protobuf/google/protobuf/struct_pb";
import * as extproto_ext_pb from "../../../../../../../extproto/ext_pb";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";
import * as github_com_solo_io_solo_kit_api_v1_ref_pb from "../../../../../../../github.com/solo-io/solo-kit/api/v1/ref_pb";
import * as github_com_solo_io_solo_kit_api_v1_solo_kit_pb from "../../../../../../../github.com/solo-io/solo-kit/api/v1/solo-kit_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/options_pb";

export class VirtualHostOptionSpec extends jspb.Message {
  hasOptions(): boolean;
  clearOptions(): void;
  getOptions(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.VirtualHostOptions | undefined;
  setOptions(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.VirtualHostOptions): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VirtualHostOptionSpec.AsObject;
  static toObject(includeInstance: boolean, msg: VirtualHostOptionSpec): VirtualHostOptionSpec.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: VirtualHostOptionSpec, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VirtualHostOptionSpec;
  static deserializeBinaryFromReader(message: VirtualHostOptionSpec, reader: jspb.BinaryReader): VirtualHostOptionSpec;
}

export namespace VirtualHostOptionSpec {
  export type AsObject = {
    options?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.VirtualHostOptions.AsObject,
  }
}

export class RouteOptionSpec extends jspb.Message {
  hasOptions(): boolean;
  clearOptions(): void;
  getOptions(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.RouteOptions | undefined;
  setOptions(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.RouteOptions): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RouteOptionSpec.AsObject;
  static toObject(includeInstance: boolean, msg: RouteOptionSpec): RouteOptionSpec.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RouteOptionSpec, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RouteOptionSpec;
  static deserializeBinaryFromReader(message: RouteOptionSpec, reader: jspb.BinaryReader): RouteOptionSpec;
}

export namespace RouteOptionSpec {
  export type AsObject = {
    options?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pb.RouteOptions.AsObject,
  }
}

export class VirtualHostOptionStatus extends jspb.Message {
  getState(): VirtualHostOptionStatus.StateMap[keyof VirtualHostOptionStatus.StateMap];
  setState(value: VirtualHostOptionStatus.StateMap[keyof VirtualHostOptionStatus.StateMap]): void;

  getReason(): string;
  setReason(value: string): void;

  getReportedBy(): string;
  setReportedBy(value: string): void;

  getSubresourceStatusesMap(): jspb.Map<string, VirtualHostOptionStatus>;
  clearSubresourceStatusesMap(): void;
  hasDetails(): boolean;
  clearDetails(): void;
  getDetails(): google_protobuf_struct_pb.Struct | undefined;
  setDetails(value?: google_protobuf_struct_pb.Struct): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VirtualHostOptionStatus.AsObject;
  static toObject(includeInstance: boolean, msg: VirtualHostOptionStatus): VirtualHostOptionStatus.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: VirtualHostOptionStatus, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VirtualHostOptionStatus;
  static deserializeBinaryFromReader(message: VirtualHostOptionStatus, reader: jspb.BinaryReader): VirtualHostOptionStatus;
}

export namespace VirtualHostOptionStatus {
  export type AsObject = {
    state: VirtualHostOptionStatus.StateMap[keyof VirtualHostOptionStatus.StateMap],
    reason: string,
    reportedBy: string,
    subresourceStatusesMap: Array<[string, VirtualHostOptionStatus.AsObject]>,
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

export class VirtualHostOptionNamespacedStatuses extends jspb.Message {
  getStatusesMap(): jspb.Map<string, VirtualHostOptionStatus>;
  clearStatusesMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VirtualHostOptionNamespacedStatuses.AsObject;
  static toObject(includeInstance: boolean, msg: VirtualHostOptionNamespacedStatuses): VirtualHostOptionNamespacedStatuses.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: VirtualHostOptionNamespacedStatuses, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VirtualHostOptionNamespacedStatuses;
  static deserializeBinaryFromReader(message: VirtualHostOptionNamespacedStatuses, reader: jspb.BinaryReader): VirtualHostOptionNamespacedStatuses;
}

export namespace VirtualHostOptionNamespacedStatuses {
  export type AsObject = {
    statusesMap: Array<[string, VirtualHostOptionStatus.AsObject]>,
  }
}

export class RouteOptionStatus extends jspb.Message {
  getState(): RouteOptionStatus.StateMap[keyof RouteOptionStatus.StateMap];
  setState(value: RouteOptionStatus.StateMap[keyof RouteOptionStatus.StateMap]): void;

  getReason(): string;
  setReason(value: string): void;

  getReportedBy(): string;
  setReportedBy(value: string): void;

  getSubresourceStatusesMap(): jspb.Map<string, RouteOptionStatus>;
  clearSubresourceStatusesMap(): void;
  hasDetails(): boolean;
  clearDetails(): void;
  getDetails(): google_protobuf_struct_pb.Struct | undefined;
  setDetails(value?: google_protobuf_struct_pb.Struct): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RouteOptionStatus.AsObject;
  static toObject(includeInstance: boolean, msg: RouteOptionStatus): RouteOptionStatus.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RouteOptionStatus, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RouteOptionStatus;
  static deserializeBinaryFromReader(message: RouteOptionStatus, reader: jspb.BinaryReader): RouteOptionStatus;
}

export namespace RouteOptionStatus {
  export type AsObject = {
    state: RouteOptionStatus.StateMap[keyof RouteOptionStatus.StateMap],
    reason: string,
    reportedBy: string,
    subresourceStatusesMap: Array<[string, RouteOptionStatus.AsObject]>,
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

export class RouteOptionNamespacedStatuses extends jspb.Message {
  getStatusesMap(): jspb.Map<string, RouteOptionStatus>;
  clearStatusesMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RouteOptionNamespacedStatuses.AsObject;
  static toObject(includeInstance: boolean, msg: RouteOptionNamespacedStatuses): RouteOptionNamespacedStatuses.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RouteOptionNamespacedStatuses, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RouteOptionNamespacedStatuses;
  static deserializeBinaryFromReader(message: RouteOptionNamespacedStatuses, reader: jspb.BinaryReader): RouteOptionNamespacedStatuses;
}

export namespace RouteOptionNamespacedStatuses {
  export type AsObject = {
    statusesMap: Array<[string, RouteOptionStatus.AsObject]>,
  }
}
