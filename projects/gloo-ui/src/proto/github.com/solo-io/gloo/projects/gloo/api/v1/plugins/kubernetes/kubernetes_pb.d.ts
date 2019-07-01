// package: kubernetes.plugins.gloo.solo.io
// file: github.com/solo-io/gloo/projects/gloo/api/v1/plugins/kubernetes/kubernetes.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../../../../../../../gogoproto/gogo_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_v1_plugins_service_spec_pb from "../../../../../../../../../github.com/solo-io/gloo/projects/gloo/api/v1/plugins/service_spec_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_v1_plugins_subset_spec_pb from "../../../../../../../../../github.com/solo-io/gloo/projects/gloo/api/v1/plugins/subset_spec_pb";

export class UpstreamSpec extends jspb.Message {
  getServiceName(): string;
  setServiceName(value: string): void;

  getServiceNamespace(): string;
  setServiceNamespace(value: string): void;

  getServicePort(): number;
  setServicePort(value: number): void;

  getSelectorMap(): jspb.Map<string, string>;
  clearSelectorMap(): void;
  hasServiceSpec(): boolean;
  clearServiceSpec(): void;
  getServiceSpec(): github_com_solo_io_gloo_projects_gloo_api_v1_plugins_service_spec_pb.ServiceSpec | undefined;
  setServiceSpec(value?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_service_spec_pb.ServiceSpec): void;

  hasSubsetSpec(): boolean;
  clearSubsetSpec(): void;
  getSubsetSpec(): github_com_solo_io_gloo_projects_gloo_api_v1_plugins_subset_spec_pb.SubsetSpec | undefined;
  setSubsetSpec(value?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_subset_spec_pb.SubsetSpec): void;

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
    serviceNamespace: string,
    servicePort: number,
    selectorMap: Array<[string, string]>,
    serviceSpec?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_service_spec_pb.ServiceSpec.AsObject,
    subsetSpec?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_subset_spec_pb.SubsetSpec.AsObject,
  }
}

