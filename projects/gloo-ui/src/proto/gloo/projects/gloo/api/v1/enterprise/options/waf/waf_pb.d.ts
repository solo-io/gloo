// package: waf.options.gloo.solo.io
// file: gloo/projects/gloo/api/v1/enterprise/options/waf/waf.proto

import * as jspb from "google-protobuf";
import * as gloo_projects_gloo_api_external_envoy_extensions_waf_waf_pb from "../../../../../../../../gloo/projects/gloo/api/external/envoy/extensions/waf/waf_pb";
import * as gogoproto_gogo_pb from "../../../../../../../../gogoproto/gogo_pb";
import * as extproto_ext_pb from "../../../../../../../../extproto/ext_pb";

export class Settings extends jspb.Message {
  getDisabled(): boolean;
  setDisabled(value: boolean): void;

  getCustomInterventionMessage(): string;
  setCustomInterventionMessage(value: string): void;

  hasCoreRuleSet(): boolean;
  clearCoreRuleSet(): void;
  getCoreRuleSet(): CoreRuleSet | undefined;
  setCoreRuleSet(value?: CoreRuleSet): void;

  clearRuleSetsList(): void;
  getRuleSetsList(): Array<gloo_projects_gloo_api_external_envoy_extensions_waf_waf_pb.RuleSet>;
  setRuleSetsList(value: Array<gloo_projects_gloo_api_external_envoy_extensions_waf_waf_pb.RuleSet>): void;
  addRuleSets(value?: gloo_projects_gloo_api_external_envoy_extensions_waf_waf_pb.RuleSet, index?: number): gloo_projects_gloo_api_external_envoy_extensions_waf_waf_pb.RuleSet;

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
    disabled: boolean,
    customInterventionMessage: string,
    coreRuleSet?: CoreRuleSet.AsObject,
    ruleSetsList: Array<gloo_projects_gloo_api_external_envoy_extensions_waf_waf_pb.RuleSet.AsObject>,
  }
}

export class CoreRuleSet extends jspb.Message {
  hasCustomSettingsString(): boolean;
  clearCustomSettingsString(): void;
  getCustomSettingsString(): string;
  setCustomSettingsString(value: string): void;

  hasCustomSettingsFile(): boolean;
  clearCustomSettingsFile(): void;
  getCustomSettingsFile(): string;
  setCustomSettingsFile(value: string): void;

  getCustomsettingstypeCase(): CoreRuleSet.CustomsettingstypeCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CoreRuleSet.AsObject;
  static toObject(includeInstance: boolean, msg: CoreRuleSet): CoreRuleSet.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CoreRuleSet, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CoreRuleSet;
  static deserializeBinaryFromReader(message: CoreRuleSet, reader: jspb.BinaryReader): CoreRuleSet;
}

export namespace CoreRuleSet {
  export type AsObject = {
    customSettingsString: string,
    customSettingsFile: string,
  }

  export enum CustomsettingstypeCase {
    CUSTOMSETTINGSTYPE_NOT_SET = 0,
    CUSTOM_SETTINGS_STRING = 2,
    CUSTOM_SETTINGS_FILE = 3,
  }
}

