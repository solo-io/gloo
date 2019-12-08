// package: rbac.options.gloo.solo.io
// file: github.com/solo-io/gloo/projects/gloo/api/v1/enterprise/options/rbac/rbac.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../../../../../../../../gogoproto/gogo_pb";

export class Settings extends jspb.Message {
  getRequireRbac(): boolean;
  setRequireRbac(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Settings.AsObject;
  static toObject(includeInstance: boolean, msg: Settings): Settings.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Settings, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Settings;
  static deserializeBinaryFromReader(message: Settings, reader: jspb.BinaryReader): Settings;
}

export namespace Settings {
  export type AsObject = {
    requireRbac: boolean,
  }
}

export class ExtensionSettings extends jspb.Message {
  getDisable(): boolean;
  setDisable(value: boolean): void;

  getPoliciesMap(): jspb.Map<string, Policy>;
  clearPoliciesMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ExtensionSettings.AsObject;
  static toObject(includeInstance: boolean, msg: ExtensionSettings): ExtensionSettings.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ExtensionSettings, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ExtensionSettings;
  static deserializeBinaryFromReader(message: ExtensionSettings, reader: jspb.BinaryReader): ExtensionSettings;
}

export namespace ExtensionSettings {
  export type AsObject = {
    disable: boolean,
    policiesMap: Array<[string, Policy.AsObject]>,
  }
}

export class Policy extends jspb.Message {
  clearPrincipalsList(): void;
  getPrincipalsList(): Array<Principal>;
  setPrincipalsList(value: Array<Principal>): void;
  addPrincipals(value?: Principal, index?: number): Principal;

  hasPermissions(): boolean;
  clearPermissions(): void;
  getPermissions(): Permissions | undefined;
  setPermissions(value?: Permissions): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Policy.AsObject;
  static toObject(includeInstance: boolean, msg: Policy): Policy.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Policy, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Policy;
  static deserializeBinaryFromReader(message: Policy, reader: jspb.BinaryReader): Policy;
}

export namespace Policy {
  export type AsObject = {
    principalsList: Array<Principal.AsObject>,
    permissions?: Permissions.AsObject,
  }
}

export class Principal extends jspb.Message {
  hasJwtPrincipal(): boolean;
  clearJwtPrincipal(): void;
  getJwtPrincipal(): JWTPrincipal | undefined;
  setJwtPrincipal(value?: JWTPrincipal): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Principal.AsObject;
  static toObject(includeInstance: boolean, msg: Principal): Principal.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Principal, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Principal;
  static deserializeBinaryFromReader(message: Principal, reader: jspb.BinaryReader): Principal;
}

export namespace Principal {
  export type AsObject = {
    jwtPrincipal?: JWTPrincipal.AsObject,
  }
}

export class JWTPrincipal extends jspb.Message {
  getClaimsMap(): jspb.Map<string, string>;
  clearClaimsMap(): void;
  getProvider(): string;
  setProvider(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): JWTPrincipal.AsObject;
  static toObject(includeInstance: boolean, msg: JWTPrincipal): JWTPrincipal.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: JWTPrincipal, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): JWTPrincipal;
  static deserializeBinaryFromReader(message: JWTPrincipal, reader: jspb.BinaryReader): JWTPrincipal;
}

export namespace JWTPrincipal {
  export type AsObject = {
    claimsMap: Array<[string, string]>,
    provider: string,
  }
}

export class Permissions extends jspb.Message {
  getPathPrefix(): string;
  setPathPrefix(value: string): void;

  clearMethodsList(): void;
  getMethodsList(): Array<string>;
  setMethodsList(value: Array<string>): void;
  addMethods(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Permissions.AsObject;
  static toObject(includeInstance: boolean, msg: Permissions): Permissions.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Permissions, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Permissions;
  static deserializeBinaryFromReader(message: Permissions, reader: jspb.BinaryReader): Permissions;
}

export namespace Permissions {
  export type AsObject = {
    pathPrefix: string,
    methodsList: Array<string>,
  }
}

