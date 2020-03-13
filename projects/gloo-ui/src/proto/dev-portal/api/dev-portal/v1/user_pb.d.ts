/* eslint-disable */
// package: devportal.solo.io
// file: dev-portal/api/dev-portal/v1/user.proto

import * as jspb from "google-protobuf";
import * as dev_portal_api_dev_portal_v1_common_pb from "../../../../dev-portal/api/dev-portal/v1/common_pb";
import * as gogoproto_gogo_pb from "../../../../gogoproto/gogo_pb";
import * as extproto_ext_pb from "../../../../protoc-gen-ext/extproto/ext_pb";

export class UserSpec extends jspb.Message {
  getUsername(): string;
  setUsername(value: string): void;

  getEmail(): string;
  setEmail(value: string): void;

  hasBasicauth(): boolean;
  clearBasicauth(): void;
  getBasicauth(): UserSpec.BasicAuth | undefined;
  setBasicauth(value?: UserSpec.BasicAuth): void;

  getAuthmethodCase(): UserSpec.AuthmethodCase;
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
    basicauth?: UserSpec.BasicAuth.AsObject,
  }

  export class BasicAuth extends jspb.Message {
    getPasswordsecretname(): string;
    setPasswordsecretname(value: string): void;

    getPasswordsecretnamespace(): string;
    setPasswordsecretnamespace(value: string): void;

    getPasswordsecretkey(): string;
    setPasswordsecretkey(value: string): void;

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
      passwordsecretname: string,
      passwordsecretnamespace: string,
      passwordsecretkey: string,
    }
  }

  export enum AuthmethodCase {
    AUTHMETHOD_NOT_SET = 0,
    BASICAUTH = 3,
  }
}

export class UserStatus extends jspb.Message {
  getObservedgeneration(): number;
  setObservedgeneration(value: number): void;

  clearUsersList(): void;
  getUsersList(): Array<string>;
  setUsersList(value: Array<string>): void;
  addUsers(value: string, index?: number): string;

  clearApidocsList(): void;
  getApidocsList(): Array<string>;
  setApidocsList(value: Array<string>): void;
  addApidocs(value: string, index?: number): string;

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
    observedgeneration: number,
    usersList: Array<string>,
    apidocsList: Array<string>,
  }
}
