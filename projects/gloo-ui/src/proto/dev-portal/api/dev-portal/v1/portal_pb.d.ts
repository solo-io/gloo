/* eslint-disable */
// package: devportal.solo.io
// file: dev-portal/api/dev-portal/v1/portal.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_timestamp_pb from "google-protobuf/google/protobuf/timestamp_pb";
import * as dev_portal_api_dev_portal_v1_common_pb from "../../../../dev-portal/api/dev-portal/v1/common_pb";
import * as gogoproto_gogo_pb from "../../../../gogoproto/gogo_pb";
import * as extproto_ext_pb from "../../../../protoc-gen-ext/extproto/ext_pb";

export class PortalSpec extends jspb.Message {
  getDisplayName(): string;
  setDisplayName(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  clearDomainsList(): void;
  getDomainsList(): Array<string>;
  setDomainsList(value: Array<string>): void;
  addDomains(value: string, index?: number): string;

  hasPrimaryLogo(): boolean;
  clearPrimaryLogo(): void;
  getPrimaryLogo(): dev_portal_api_dev_portal_v1_common_pb.DataSource | undefined;
  setPrimaryLogo(value?: dev_portal_api_dev_portal_v1_common_pb.DataSource): void;

  hasFavicon(): boolean;
  clearFavicon(): void;
  getFavicon(): dev_portal_api_dev_portal_v1_common_pb.DataSource | undefined;
  setFavicon(value?: dev_portal_api_dev_portal_v1_common_pb.DataSource): void;

  hasBanner(): boolean;
  clearBanner(): void;
  getBanner(): dev_portal_api_dev_portal_v1_common_pb.DataSource | undefined;
  setBanner(value?: dev_portal_api_dev_portal_v1_common_pb.DataSource): void;

  hasCustomStyling(): boolean;
  clearCustomStyling(): void;
  getCustomStyling(): CustomStyling | undefined;
  setCustomStyling(value?: CustomStyling): void;

  clearStaticPagesList(): void;
  getStaticPagesList(): Array<StaticPage>;
  setStaticPagesList(value: Array<StaticPage>): void;
  addStaticPages(value?: StaticPage, index?: number): StaticPage;

  hasPublishApiDocs(): boolean;
  clearPublishApiDocs(): void;
  getPublishApiDocs(): dev_portal_api_dev_portal_v1_common_pb.Selector | undefined;
  setPublishApiDocs(value?: dev_portal_api_dev_portal_v1_common_pb.Selector): void;

  clearKeyScopesList(): void;
  getKeyScopesList(): Array<KeyScope>;
  setKeyScopesList(value: Array<KeyScope>): void;
  addKeyScopes(value?: KeyScope, index?: number): KeyScope;

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
    displayName: string,
    description: string,
    domainsList: Array<string>,
    primaryLogo?: dev_portal_api_dev_portal_v1_common_pb.DataSource.AsObject,
    favicon?: dev_portal_api_dev_portal_v1_common_pb.DataSource.AsObject,
    banner?: dev_portal_api_dev_portal_v1_common_pb.DataSource.AsObject,
    customStyling?: CustomStyling.AsObject,
    staticPagesList: Array<StaticPage.AsObject>,
    publishApiDocs?: dev_portal_api_dev_portal_v1_common_pb.Selector.AsObject,
    keyScopesList: Array<KeyScope.AsObject>,
  }
}

export class PortalStatus extends jspb.Message {
  getObservedGeneration(): number;
  setObservedGeneration(value: number): void;

  getState(): PortalStatus.StateMap[keyof PortalStatus.StateMap];
  setState(value: PortalStatus.StateMap[keyof PortalStatus.StateMap]): void;

  getReason(): string;
  setReason(value: string): void;

  getPublishUrl(): string;
  setPublishUrl(value: string): void;

  clearApiDocsList(): void;
  getApiDocsList(): Array<string>;
  setApiDocsList(value: Array<string>): void;
  addApiDocs(value: string, index?: number): string;

  clearKeyScopesList(): void;
  getKeyScopesList(): Array<KeyScopeStatus>;
  setKeyScopesList(value: Array<KeyScopeStatus>): void;
  addKeyScopes(value?: KeyScopeStatus, index?: number): KeyScopeStatus;

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
    observedGeneration: number,
    state: PortalStatus.StateMap[keyof PortalStatus.StateMap],
    reason: string,
    publishUrl: string,
    apiDocsList: Array<string>,
    keyScopesList: Array<KeyScopeStatus.AsObject>,
  }

  export interface StateMap {
    PENDING: 0;
    PUBLISHING: 1;
    PUBLISHED: 2;
    INVALID: 3;
    FAILED: 4;
  }

  export const State: StateMap;
}

export class CustomStyling extends jspb.Message {
  getPrimaryColor(): string;
  setPrimaryColor(value: string): void;

  getSecondaryColor(): string;
  setSecondaryColor(value: string): void;

  getBackgroundColor(): string;
  setBackgroundColor(value: string): void;

  getNavigationLinksColorOverride(): string;
  setNavigationLinksColorOverride(value: string): void;

  getButtonColorOverride(): string;
  setButtonColorOverride(value: string): void;

  getDefaultTextColor(): string;
  setDefaultTextColor(value: string): void;

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
    primaryColor: string,
    secondaryColor: string,
    backgroundColor: string,
    navigationLinksColorOverride: string,
    buttonColorOverride: string,
    defaultTextColor: string,
  }
}

export class StaticPage extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  getPath(): string;
  setPath(value: string): void;

  getNavigationLinkName(): string;
  setNavigationLinkName(value: string): void;

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
    navigationLinkName: string,
    content?: dev_portal_api_dev_portal_v1_common_pb.DataSource.AsObject,
  }
}

export class KeyScope extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getNamespace(): string;
  setNamespace(value: string): void;

  hasApiDocs(): boolean;
  clearApiDocs(): void;
  getApiDocs(): dev_portal_api_dev_portal_v1_common_pb.Selector | undefined;
  setApiDocs(value?: dev_portal_api_dev_portal_v1_common_pb.Selector): void;

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
    apiDocs?: dev_portal_api_dev_portal_v1_common_pb.Selector.AsObject,
  }
}

export class KeyScopeStatus extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  hasCreatedDate(): boolean;
  clearCreatedDate(): void;
  getCreatedDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setCreatedDate(value?: google_protobuf_timestamp_pb.Timestamp): void;

  clearAccessibleApiDocsList(): void;
  getAccessibleApiDocsList(): Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef>;
  setAccessibleApiDocsList(value: Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef>): void;
  addAccessibleApiDocs(value?: dev_portal_api_dev_portal_v1_common_pb.ObjectRef, index?: number): dev_portal_api_dev_portal_v1_common_pb.ObjectRef;

  clearProvisionedKeysList(): void;
  getProvisionedKeysList(): Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef>;
  setProvisionedKeysList(value: Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef>): void;
  addProvisionedKeys(value?: dev_portal_api_dev_portal_v1_common_pb.ObjectRef, index?: number): dev_portal_api_dev_portal_v1_common_pb.ObjectRef;

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
    createdDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    accessibleApiDocsList: Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef.AsObject>,
    provisionedKeysList: Array<dev_portal_api_dev_portal_v1_common_pb.ObjectRef.AsObject>,
  }
}
