// package: rbac.plugins.gloo.solo.io
// file: github.com/solo-io/solo-projects/projects/gloo/api/v1/plugins/rbac/rbac.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../../../../../../../gogoproto/gogo_pb";

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

export class Config extends jspb.Message {
  getPoliciesMap(): jspb.Map<string, Policy>;
  clearPoliciesMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Config.AsObject;
  static toObject(includeInstance: boolean, msg: Config): Config.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Config, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Config;
  static deserializeBinaryFromReader(message: Config, reader: jspb.BinaryReader): Config;
}

export namespace Config {
  export type AsObject = {
    policiesMap: Array<[string, Policy.AsObject]>,
  }
}

export class VhostExtension extends jspb.Message {
  hasConfig(): boolean;
  clearConfig(): void;
  getConfig(): Config | undefined;
  setConfig(value?: Config): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VhostExtension.AsObject;
  static toObject(includeInstance: boolean, msg: VhostExtension): VhostExtension.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: VhostExtension, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VhostExtension;
  static deserializeBinaryFromReader(message: VhostExtension, reader: jspb.BinaryReader): VhostExtension;
}

export namespace VhostExtension {
  export type AsObject = {
    config?: Config.AsObject,
  }
}

export class RouteExtension extends jspb.Message {
  hasDisable(): boolean;
  clearDisable(): void;
  getDisable(): boolean;
  setDisable(value: boolean): void;

  hasConfig(): boolean;
  clearConfig(): void;
  getConfig(): Config | undefined;
  setConfig(value?: Config): void;

  getRouteCase(): RouteExtension.RouteCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RouteExtension.AsObject;
  static toObject(includeInstance: boolean, msg: RouteExtension): RouteExtension.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RouteExtension, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RouteExtension;
  static deserializeBinaryFromReader(message: RouteExtension, reader: jspb.BinaryReader): RouteExtension;
}

export namespace RouteExtension {
  export type AsObject = {
    disable: boolean,
    config?: Config.AsObject,
  }

  export enum RouteCase {
    ROUTE_NOT_SET = 0,
    DISABLE = 1,
    CONFIG = 2,
  }
}

