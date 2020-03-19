/* eslint-disable */
// package: devportal.solo.io
// file: dev-portal/api/dev-portal/v1/portal.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_timestamp_pb from "google-protobuf/google/protobuf/timestamp_pb";
import * as dev_portal_api_dev_portal_v1_common_pb from "../../../../dev-portal/api/dev-portal/v1/common_pb";
import * as gogoproto_gogo_pb from "../../../../gogoproto/gogo_pb";
import * as extproto_ext_pb from "../../../../protoc-gen-ext/extproto/ext_pb";

export class PortalSpec extends jspb.Message {
  getDisplayname(): string;
  setDisplayname(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  clearDomainsList(): void;
  getDomainsList(): Array<string>;
  setDomainsList(value: Array<string>): void;
  addDomains(value: string, index?: number): string;

  hasPrimarylogo(): boolean;
  clearPrimarylogo(): void;
  getPrimarylogo(): dev_portal_api_dev_portal_v1_common_pb.DataSource | undefined;
  setPrimarylogo(value?: dev_portal_api_dev_portal_v1_common_pb.DataSource): void;

  hasFavicon(): boolean;
  clearFavicon(): void;
  getFavicon(): dev_portal_api_dev_portal_v1_common_pb.DataSource | undefined;
  setFavicon(value?: dev_portal_api_dev_portal_v1_common_pb.DataSource): void;

  hasBackgroundimage(): boolean;
  clearBackgroundimage(): void;
  getBackgroundimage(): dev_portal_api_dev_portal_v1_common_pb.DataSource | undefined;
  setBackgroundimage(value?: dev_portal_api_dev_portal_v1_common_pb.DataSource): void;

  hasCustomstyling(): boolean;
  clearCustomstyling(): void;
  getCustomstyling(): CustomStyling | undefined;
  setCustomstyling(value?: CustomStyling): void;

  clearStaticpagesList(): void;
  getStaticpagesList(): Array<StaticPage>;
  setStaticpagesList(value: Array<StaticPage>): void;
  addStaticpages(value?: StaticPage, index?: number): StaticPage;

  hasPublishapidocs(): boolean;
  clearPublishapidocs(): void;
  getPublishapidocs(): dev_portal_api_dev_portal_v1_common_pb.Selector | undefined;
  setPublishapidocs(value?: dev_portal_api_dev_portal_v1_common_pb.Selector): void;

  clearKeyscopesList(): void;
  getKeyscopesList(): Array<KeyScope>;
  setKeyscopesList(value: Array<KeyScope>): void;
  addKeyscopes(value?: KeyScope, index?: number): KeyScope;

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
    displayname: string,
    description: string,
    domainsList: Array<string>,
    primarylogo?: dev_portal_api_dev_portal_v1_common_pb.DataSource.AsObject,
    favicon?: dev_portal_api_dev_portal_v1_common_pb.DataSource.AsObject,
    backgroundimage?: dev_portal_api_dev_portal_v1_common_pb.DataSource.AsObject,
    customstyling?: CustomStyling.AsObject,
    staticpagesList: Array<StaticPage.AsObject>,
    publishapidocs?: dev_portal_api_dev_portal_v1_common_pb.Selector.AsObject,
    keyscopesList: Array<KeyScope.AsObject>,
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

  clearKeyscopesList(): void;
  getKeyscopesList(): Array<KeyScopeStatus>;
  setKeyscopesList(value: Array<KeyScopeStatus>): void;
  addKeyscopes(value?: KeyScopeStatus, index?: number): KeyScopeStatus;

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
    keyscopesList: Array<KeyScopeStatus.AsObject>,
  }

  export interface StateMap {
    PENDING: 0;
    PUBLISHED: 1;
    INVALID: 2;
    FAILED: 3;
  }

  export const State: StateMap;
}

export class CustomStyling extends jspb.Message {
  getPrimarycolor(): string;
  setPrimarycolor(value: string): void;

  getSecondarycolor(): string;
  setSecondarycolor(value: string): void;

  getBackgroundcolor(): string;
  setBackgroundcolor(value: string): void;

  getNavigationlinkscoloroverride(): string;
  setNavigationlinkscoloroverride(value: string): void;

  getButtoncoloroverride(): string;
  setButtoncoloroverride(value: string): void;

  getDefaulttextcolor(): string;
  setDefaulttextcolor(value: string): void;

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
    primarycolor: string,
    secondarycolor: string,
    backgroundcolor: string,
    navigationlinkscoloroverride: string,
    buttoncoloroverride: string,
    defaulttextcolor: string,
  }
}

export class StaticPage extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  getPath(): string;
  setPath(value: string): void;

  getNavigationlinkname(): string;
  setNavigationlinkname(value: string): void;

  hasContent(): boolean;
  clearContent(): void;
  getContent(): dev_portal_api_dev_portal_v1_common_pb.DataSource | undefined;
  setContent(value?: dev_portal_api_dev_portal_v1_common_pb.DataSource): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StaticPage.AsObject;
  static toObject(includeInstance: boolean, msg: StaticPage): StaticPage.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: StaticPage, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StaticPage;
  static deserializeBinaryFromReader(message: StaticPage, reader: jspb.BinaryReader): StaticPage;
}

export namespace StaticPage {
  export type AsObject = {
    name: string,
    description: string,
    path: string,
    navigationlinkname: string,
    content?: dev_portal_api_dev_portal_v1_common_pb.DataSource.AsObject,
  }
}

export class KeyScope extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getNamespace(): string;
  setNamespace(value: string): void;

  hasApidocs(): boolean;
  clearApidocs(): void;
  getApidocs(): dev_portal_api_dev_portal_v1_common_pb.Selector | undefined;
  setApidocs(value?: dev_portal_api_dev_portal_v1_common_pb.Selector): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): KeyScope.AsObject;
  static toObject(includeInstance: boolean, msg: KeyScope): KeyScope.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: KeyScope, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): KeyScope;
  static deserializeBinaryFromReader(message: KeyScope, reader: jspb.BinaryReader): KeyScope;
}

export namespace KeyScope {
  export type AsObject = {
    name: string,
    namespace: string,
    apidocs?: dev_portal_api_dev_portal_v1_common_pb.Selector.AsObject,
  }
}

export class KeyScopeStatus extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  hasCreateddate(): boolean;
  clearCreateddate(): void;
  getCreateddate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setCreateddate(value?: google_protobuf_timestamp_pb.Timestamp): void;

  clearAccessibleapidocsList(): void;
  getAccessibleapidocsList(): Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef>;
  setAccessibleapidocsList(value: Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef>): void;
  addAccessibleapidocs(value?: dev_portal_api_dev_portal_v1_common_pb.ObjectRef, index?: number): dev_portal_api_dev_portal_v1_common_pb.ObjectRef;

  clearProvisionedkeysList(): void;
  getProvisionedkeysList(): Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef>;
  setProvisionedkeysList(value: Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef>): void;
  addProvisionedkeys(value?: dev_portal_api_dev_portal_v1_common_pb.ObjectRef, index?: number): dev_portal_api_dev_portal_v1_common_pb.ObjectRef;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): KeyScopeStatus.AsObject;
  static toObject(includeInstance: boolean, msg: KeyScopeStatus): KeyScopeStatus.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: KeyScopeStatus, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): KeyScopeStatus;
  static deserializeBinaryFromReader(message: KeyScopeStatus, reader: jspb.BinaryReader): KeyScopeStatus;
}

export namespace KeyScopeStatus {
  export type AsObject = {
    name: string,
    createddate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    accessibleapidocsList: Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef.AsObject>,
    provisionedkeysList: Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef.AsObject>,
  }
}
