/* eslint-disable */
// package: static.options.gloo.solo.io
// file: gloo/projects/gloo/api/v1/options/static/static.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../../../../../gogoproto/gogo_pb";
import * as extproto_ext_pb from "../../../../../../../protoc-gen-ext/extproto/ext_pb";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";
import * as gloo_projects_gloo_api_v1_options_service_spec_pb from "../../../../../../../gloo/projects/gloo/api/v1/options/service_spec_pb";

export class UpstreamSpec extends jspb.Message {
  clearHostsList(): void;
  getHostsList(): Array<Host>;
  setHostsList(value: Array<Host>): void;
  addHosts(value?: Host, index?: number): Host;

  getUseTls(): boolean;
  setUseTls(value: boolean): void;

  hasServiceSpec(): boolean;
  clearServiceSpec(): void;
  getServiceSpec(): gloo_projects_gloo_api_v1_options_service_spec_pb.ServiceSpec | undefined;
  setServiceSpec(value?: gloo_projects_gloo_api_v1_options_service_spec_pb.ServiceSpec): void;

  hasAutoSniRewrite(): boolean;
  clearAutoSniRewrite(): void;
  getAutoSniRewrite(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setAutoSniRewrite(value?: google_protobuf_wrappers_pb.BoolValue): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpstreamSpec.AsObject;
  static toObject(includeInstance: boolean, msg: UpstreamSpec): UpstreamSpec.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpstreamSpec, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpstreamSpec;
  static deserializeBinaryFromReader(message: UpstreamSpec, reader: jspb.BinaryReader): UpstreamSpec;
}

export namespace UpstreamSpec {
  export type AsObject = {
    hostsList: Array<Host.AsObject>,
    useTls: boolean,
    serviceSpec?: gloo_projects_gloo_api_v1_options_service_spec_pb.ServiceSpec.AsObject,
    autoSniRewrite?: google_protobuf_wrappers_pb.BoolValue.AsObject,
  }
}

export class Host extends jspb.Message {
  getAddr(): string;
  setAddr(value: string): void;

  getPort(): number;
  setPort(value: number): void;

  getSniAddr(): string;
  setSniAddr(value: string): void;

  hasHealthCheckConfig(): boolean;
  clearHealthCheckConfig(): void;
  getHealthCheckConfig(): Host.HealthCheckConfig | undefined;
  setHealthCheckConfig(value?: Host.HealthCheckConfig): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Host.AsObject;
  static toObject(includeInstance: boolean, msg: Host): Host.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Host, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Host;
  static deserializeBinaryFromReader(message: Host, reader: jspb.BinaryReader): Host;
}

export namespace Host {
  export type AsObject = {
    addr: string,
    port: number,
    sniAddr: string,
    healthCheckConfig?: Host.HealthCheckConfig.AsObject,
  }

  export class HealthCheckConfig extends jspb.Message {
    getPath(): string;
    setPath(value: string): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): HealthCheckConfig.AsObject;
    static toObject(includeInstance: boolean, msg: HealthCheckConfig): HealthCheckConfig.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: HealthCheckConfig, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): HealthCheckConfig;
    static deserializeBinaryFromReader(message: HealthCheckConfig, reader: jspb.BinaryReader): HealthCheckConfig;
  }

  export namespace HealthCheckConfig {
    export type AsObject = {
      path: string,
    }
  }
}
