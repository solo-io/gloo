// package: consul.options.gloo.solo.io
// file: github.com/solo-io/gloo/projects/gloo/api/v1/options/consul/consul.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../../../../../../../gogoproto/gogo_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_v1_options_service_spec_pb from "../../../../../../../../../github.com/solo-io/gloo/projects/gloo/api/v1/options/service_spec_pb";

export class UpstreamSpec extends jspb.Message {
  getServiceName(): string;
  setServiceName(value: string): void;

  clearServiceTagsList(): void;
  getServiceTagsList(): Array<string>;
  setServiceTagsList(value: Array<string>): void;
  addServiceTags(value: string, index?: number): string;

  hasServiceSpec(): boolean;
  clearServiceSpec(): void;
  getServiceSpec(): github_com_solo_io_gloo_projects_gloo_api_v1_options_service_spec_pb.ServiceSpec | undefined;
  setServiceSpec(value?: github_com_solo_io_gloo_projects_gloo_api_v1_options_service_spec_pb.ServiceSpec): void;

  getConnectEnabled(): boolean;
  setConnectEnabled(value: boolean): void;

  clearDataCentersList(): void;
  getDataCentersList(): Array<string>;
  setDataCentersList(value: Array<string>): void;
  addDataCenters(value: string, index?: number): string;

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
    serviceName: string,
    serviceTagsList: Array<string>,
    serviceSpec?: github_com_solo_io_gloo_projects_gloo_api_v1_options_service_spec_pb.ServiceSpec.AsObject,
    connectEnabled: boolean,
    dataCentersList: Array<string>,
  }
}

