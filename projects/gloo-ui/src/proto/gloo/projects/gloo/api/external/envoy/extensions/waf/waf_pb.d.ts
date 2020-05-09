/* eslint-disable */
// package: envoy.config.filter.http.modsecurity.v2
// file: gloo/projects/gloo/api/external/envoy/extensions/waf/waf.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../../../../../../gogoproto/gogo_pb";

export class AuditLogging extends jspb.Message {
  getAction(): AuditLogging.AuditLogActionMap[keyof AuditLogging.AuditLogActionMap];
  setAction(value: AuditLogging.AuditLogActionMap[keyof AuditLogging.AuditLogActionMap]): void;

  getLocation(): AuditLogging.AuditLogLocationMap[keyof AuditLogging.AuditLogLocationMap];
  setLocation(value: AuditLogging.AuditLogLocationMap[keyof AuditLogging.AuditLogLocationMap]): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AuditLogging.AsObject;
  static toObject(includeInstance: boolean, msg: AuditLogging): AuditLogging.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: AuditLogging, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AuditLogging;
  static deserializeBinaryFromReader(message: AuditLogging, reader: jspb.BinaryReader): AuditLogging;
}

export namespace AuditLogging {
  export type AsObject = {
    action: AuditLogging.AuditLogActionMap[keyof AuditLogging.AuditLogActionMap],
    location: AuditLogging.AuditLogLocationMap[keyof AuditLogging.AuditLogLocationMap],
  }

  export interface AuditLogActionMap {
    NEVER: 0;
    RELEVANT_ONLY: 1;
    ALWAYS: 2;
  }

  export const AuditLogAction: AuditLogActionMap;

  export interface AuditLogLocationMap {
    FILTER_STATE: 0;
    DYNAMIC_METADATA: 1;
  }

  export const AuditLogLocation: AuditLogLocationMap;
}

export class ModSecurity extends jspb.Message {
  getDisabled(): boolean;
  setDisabled(value: boolean): void;

  clearRuleSetsList(): void;
  getRuleSetsList(): Array<RuleSet>;
  setRuleSetsList(value: Array<RuleSet>): void;
  addRuleSets(value?: RuleSet, index?: number): RuleSet;

  getCustomInterventionMessage(): string;
  setCustomInterventionMessage(value: string): void;

  hasAuditLogging(): boolean;
  clearAuditLogging(): void;
  getAuditLogging(): AuditLogging | undefined;
  setAuditLogging(value?: AuditLogging): void;

  getRegressionLogs(): boolean;
  setRegressionLogs(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ModSecurity.AsObject;
  static toObject(includeInstance: boolean, msg: ModSecurity): ModSecurity.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ModSecurity, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ModSecurity;
  static deserializeBinaryFromReader(message: ModSecurity, reader: jspb.BinaryReader): ModSecurity;
}

export namespace ModSecurity {
  export type AsObject = {
    disabled: boolean,
    ruleSetsList: Array<RuleSet.AsObject>,
    customInterventionMessage: string,
    auditLogging?: AuditLogging.AsObject,
    regressionLogs: boolean,
  }
}

export class RuleSet extends jspb.Message {
  getRuleStr(): string;
  setRuleStr(value: string): void;

  clearFilesList(): void;
  getFilesList(): Array<string>;
  setFilesList(value: Array<string>): void;
  addFiles(value: string, index?: number): string;

  getDirectory(): string;
  setDirectory(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RuleSet.AsObject;
  static toObject(includeInstance: boolean, msg: RuleSet): RuleSet.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RuleSet, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RuleSet;
  static deserializeBinaryFromReader(message: RuleSet, reader: jspb.BinaryReader): RuleSet;
}

export namespace RuleSet {
  export type AsObject = {
    ruleStr: string,
    filesList: Array<string>,
    directory: string,
  }
}

export class ModSecurityPerRoute extends jspb.Message {
  getDisabled(): boolean;
  setDisabled(value: boolean): void;

  clearRuleSetsList(): void;
  getRuleSetsList(): Array<RuleSet>;
  setRuleSetsList(value: Array<RuleSet>): void;
  addRuleSets(value?: RuleSet, index?: number): RuleSet;

  getCustomInterventionMessage(): string;
  setCustomInterventionMessage(value: string): void;

  hasAuditLogging(): boolean;
  clearAuditLogging(): void;
  getAuditLogging(): AuditLogging | undefined;
  setAuditLogging(value?: AuditLogging): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ModSecurityPerRoute.AsObject;
  static toObject(includeInstance: boolean, msg: ModSecurityPerRoute): ModSecurityPerRoute.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ModSecurityPerRoute, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ModSecurityPerRoute;
  static deserializeBinaryFromReader(message: ModSecurityPerRoute, reader: jspb.BinaryReader): ModSecurityPerRoute;
}

export namespace ModSecurityPerRoute {
  export type AsObject = {
    disabled: boolean,
    ruleSetsList: Array<RuleSet.AsObject>,
    customInterventionMessage: string,
    auditLogging?: AuditLogging.AsObject,
  }
}
