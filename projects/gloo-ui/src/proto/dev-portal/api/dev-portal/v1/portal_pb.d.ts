/* eslint-disable */
// package: devportal.solo.io
// file: dev-portal/api/dev-portal/v1/portal.proto

import * as jspb from "google-protobuf";
import * as dev_portal_api_dev_portal_v1_common_pb from "../../../../dev-portal/api/dev-portal/v1/common_pb";
import * as gogoproto_gogo_pb from "../../../../gogoproto/gogo_pb";
import * as extproto_ext_pb from "../../../../protoc-gen-ext/extproto/ext_pb";

export class PortalSpec extends jspb.Message {
  getTitle(): string;
  setTitle(value: string): void;

  getCompanyname(): string;
  setCompanyname(value: string): void;

  hasLogo(): boolean;
  clearLogo(): void;
  getLogo(): dev_portal_api_dev_portal_v1_common_pb.DataSource | undefined;
  setLogo(value?: dev_portal_api_dev_portal_v1_common_pb.DataSource): void;

  hasCustomstyling(): boolean;
  clearCustomstyling(): void;
  getCustomstyling(): CustomStyling | undefined;
  setCustomstyling(value?: CustomStyling): void;

  getApidoclabelsMap(): jspb.Map<string, string>;
  clearApidoclabelsMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PortalSpec.AsObject;
  static toObject(includeInstance: boolean, msg: PortalSpec): PortalSpec.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: PortalSpec, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PortalSpec;
  static deserializeBinaryFromReader(message: PortalSpec, reader: jspb.BinaryReader): PortalSpec;
}

export namespace PortalSpec {
  export type AsObject = {
    title: string,
    companyname: string,
    logo?: dev_portal_api_dev_portal_v1_common_pb.DataSource.AsObject,
    customstyling?: CustomStyling.AsObject,
    apidoclabelsMap: Array<[string, string]>,
  }
}

export class CustomStyling extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CustomStyling.AsObject;
  static toObject(includeInstance: boolean, msg: CustomStyling): CustomStyling.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CustomStyling, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CustomStyling;
  static deserializeBinaryFromReader(message: CustomStyling, reader: jspb.BinaryReader): CustomStyling;
}

export namespace CustomStyling {
  export type AsObject = {
  }
}

export class PortalStatus extends jspb.Message {
  getObservedgeneration(): number;
  setObservedgeneration(value: number): void;

  getState(): PortalStatus.StateMap[keyof PortalStatus.StateMap];
  setState(value: PortalStatus.StateMap[keyof PortalStatus.StateMap]): void;

  getReason(): string;
  setReason(value: string): void;

  getPublishurl(): string;
  setPublishurl(value: string): void;

  clearApidocsList(): void;
  getApidocsList(): Array<string>;
  setApidocsList(value: Array<string>): void;
  addApidocs(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PortalStatus.AsObject;
  static toObject(includeInstance: boolean, msg: PortalStatus): PortalStatus.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: PortalStatus, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PortalStatus;
  static deserializeBinaryFromReader(message: PortalStatus, reader: jspb.BinaryReader): PortalStatus;
}

export namespace PortalStatus {
  export type AsObject = {
    observedgeneration: number,
    state: PortalStatus.StateMap[keyof PortalStatus.StateMap],
    reason: string,
    publishurl: string,
    apidocsList: Array<string>,
  }

  export interface StateMap {
    PENDING: 0;
    PUBLISHED: 1;
    INVALID: 2;
    FAILED: 3;
  }

  export const State: StateMap;
}
