// package: waf.plugins.gloo.solo.io
// file: github.com/solo-io/gloo/projects/gloo/api/v1/enterprise/plugins/waf/waf.proto

import * as jspb from "google-protobuf";
import * as github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_waf_waf_pb from "../../../../../../../../../../github.com/solo-io/gloo/projects/gloo/api/external/envoy/extensions/waf/waf_pb";
import * as gogoproto_gogo_pb from "../../../../../../../../../../gogoproto/gogo_pb";

export class Settings extends jspb.Message {
  getDisabled(): boolean;
  setDisabled(value: boolean): void;

  hasCoreRuleSet(): boolean;
  clearCoreRuleSet(): void;
  getCoreRuleSet(): CoreRuleSet | undefined;
  setCoreRuleSet(value?: CoreRuleSet): void;

  clearRuleSetsList(): void;
  getRuleSetsList(): Array<github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_waf_waf_pb.RuleSet>;
  setRuleSetsList(value: Array<github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_waf_waf_pb.RuleSet>): void;
  addRuleSets(value?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_waf_waf_pb.RuleSet, index?: number): github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_waf_waf_pb.RuleSet;

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
    coreRuleSet?: CoreRuleSet.AsObject,
    ruleSetsList: Array<github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_waf_waf_pb.RuleSet.AsObject>,
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

export class VhostSettings extends jspb.Message {
  getDisabled(): boolean;
  setDisabled(value: boolean): void;

  hasSettings(): boolean;
  clearSettings(): void;
  getSettings(): Settings | undefined;
  setSettings(value?: Settings): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VhostSettings.AsObject;
  static toObject(includeInstance: boolean, msg: VhostSettings): VhostSettings.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: VhostSettings, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VhostSettings;
  static deserializeBinaryFromReader(message: VhostSettings, reader: jspb.BinaryReader): VhostSettings;
}

export namespace VhostSettings {
  export type AsObject = {
    disabled: boolean,
    settings?: Settings.AsObject,
  }
}

export class RouteSettings extends jspb.Message {
  getDisabled(): boolean;
  setDisabled(value: boolean): void;

  hasSettings(): boolean;
  clearSettings(): void;
  getSettings(): Settings | undefined;
  setSettings(value?: Settings): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RouteSettings.AsObject;
  static toObject(includeInstance: boolean, msg: RouteSettings): RouteSettings.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RouteSettings, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RouteSettings;
  static deserializeBinaryFromReader(message: RouteSettings, reader: jspb.BinaryReader): RouteSettings;
}

export namespace RouteSettings {
  export type AsObject = {
    disabled: boolean,
    settings?: Settings.AsObject,
  }
}

