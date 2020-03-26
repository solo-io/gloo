/* eslint-disable */
// package: devportal.solo.io
// file: dev-portal/api/dev-portal/v1/user.proto

import * as jspb from "google-protobuf";
import * as dev_portal_api_dev_portal_v1_access_level_pb from "../../../../dev-portal/api/dev-portal/v1/access_level_pb";
import * as gogoproto_gogo_pb from "../../../../gogoproto/gogo_pb";
import * as extproto_ext_pb from "../../../../protoc-gen-ext/extproto/ext_pb";

export class UserSpec extends jspb.Message {
  getUsername(): string;
  setUsername(value: string): void;

  getEmail(): string;
  setEmail(value: string): void;

  hasBasicAuth(): boolean;
  clearBasicAuth(): void;
  getBasicAuth(): UserSpec.BasicAuth | undefined;
  setBasicAuth(value?: UserSpec.BasicAuth): void;

  hasAccessLevel(): boolean;
  clearAccessLevel(): void;
  getAccessLevel(): dev_portal_api_dev_portal_v1_access_level_pb.AccessLevel | undefined;
  setAccessLevel(value?: dev_portal_api_dev_portal_v1_access_level_pb.AccessLevel): void;

  getAuthMethodCase(): UserSpec.AuthMethodCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserSpec.AsObject;
  static toObject(includeInstance: boolean, msg: UserSpec): UserSpec.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UserSpec, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserSpec;
  static deserializeBinaryFromReader(message: UserSpec, reader: jspb.BinaryReader): UserSpec;
}

export namespace UserSpec {
  export type AsObject = {
    username: string,
    email: string,
    basicAuth?: UserSpec.BasicAuth.AsObject,
    accessLevel?: dev_portal_api_dev_portal_v1_access_level_pb.AccessLevel.AsObject,
  }

  export class BasicAuth extends jspb.Message {
    getPasswordSecretName(): string;
    setPasswordSecretName(value: string): void;

    getPasswordSecretNamespace(): string;
    setPasswordSecretNamespace(value: string): void;

    getPasswordSecretKey(): string;
    setPasswordSecretKey(value: string): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): BasicAuth.AsObject;
    static toObject(includeInstance: boolean, msg: BasicAuth): BasicAuth.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: BasicAuth, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): BasicAuth;
    static deserializeBinaryFromReader(message: BasicAuth, reader: jspb.BinaryReader): BasicAuth;
  }

  export namespace BasicAuth {
    export type AsObject = {
      passwordSecretName: string,
      passwordSecretNamespace: string,
      passwordSecretKey: string,
    }
  }

  export enum AuthMethodCase {
    AUTH_METHOD_NOT_SET = 0,
    BASIC_AUTH = 3,
  }
}

export class UserStatus extends jspb.Message {
  getObservedGeneration(): number;
  setObservedGeneration(value: number): void;

  getState(): UserStatus.StateMap[keyof UserStatus.StateMap];
  setState(value: UserStatus.StateMap[keyof UserStatus.StateMap]): void;

  hasAccessLevel(): boolean;
  clearAccessLevel(): void;
  getAccessLevel(): dev_portal_api_dev_portal_v1_access_level_pb.AccessLevelStatus | undefined;
  setAccessLevel(value?: dev_portal_api_dev_portal_v1_access_level_pb.AccessLevelStatus): void;

  getHasLoggedIn(): boolean;
  setHasLoggedIn(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserStatus.AsObject;
  static toObject(includeInstance: boolean, msg: UserStatus): UserStatus.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UserStatus, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserStatus;
  static deserializeBinaryFromReader(message: UserStatus, reader: jspb.BinaryReader): UserStatus;
}

export namespace UserStatus {
  export type AsObject = {
    observedGeneration: number,
    state: UserStatus.StateMap[keyof UserStatus.StateMap],
    accessLevel?: dev_portal_api_dev_portal_v1_access_level_pb.AccessLevelStatus.AsObject,
    hasLoggedIn: boolean,
  }

  export interface StateMap {
    PENDING: 0;
    ACCEPTED: 1;
    INVALID: 2;
  }

  export const State: StateMap;
}
