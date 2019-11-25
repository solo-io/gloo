// package: aws_ec2.options.gloo.solo.io
// file: github.com/solo-io/gloo/projects/gloo/api/v1/options/aws/ec2/aws_ec2.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../../../../../../../../gogoproto/gogo_pb";
import * as github_com_solo_io_solo_kit_api_v1_ref_pb from "../../../../../../../../../../github.com/solo-io/solo-kit/api/v1/ref_pb";

export class UpstreamSpec extends jspb.Message {
  getRegion(): string;
  setRegion(value: string): void;

  hasSecretRef(): boolean;
  clearSecretRef(): void;
  getSecretRef(): github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef | undefined;
  setSecretRef(value?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef): void;

  getRoleArn(): string;
  setRoleArn(value: string): void;

  clearFiltersList(): void;
  getFiltersList(): Array<TagFilter>;
  setFiltersList(value: Array<TagFilter>): void;
  addFilters(value?: TagFilter, index?: number): TagFilter;

  getPublicIp(): boolean;
  setPublicIp(value: boolean): void;

  getPort(): number;
  setPort(value: number): void;

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
    region: string,
    secretRef?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject,
    roleArn: string,
    filtersList: Array<TagFilter.AsObject>,
    publicIp: boolean,
    port: number,
  }
}

export class TagFilter extends jspb.Message {
  hasKey(): boolean;
  clearKey(): void;
  getKey(): string;
  setKey(value: string): void;

  hasKvPair(): boolean;
  clearKvPair(): void;
  getKvPair(): TagFilter.KvPair | undefined;
  setKvPair(value?: TagFilter.KvPair): void;

  getSpecCase(): TagFilter.SpecCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TagFilter.AsObject;
  static toObject(includeInstance: boolean, msg: TagFilter): TagFilter.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: TagFilter, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TagFilter;
  static deserializeBinaryFromReader(message: TagFilter, reader: jspb.BinaryReader): TagFilter;
}

export namespace TagFilter {
  export type AsObject = {
    key: string,
    kvPair?: TagFilter.KvPair.AsObject,
  }

  export class KvPair extends jspb.Message {
    getKey(): string;
    setKey(value: string): void;

    getValue(): string;
    setValue(value: string): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): KvPair.AsObject;
    static toObject(includeInstance: boolean, msg: KvPair): KvPair.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: KvPair, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): KvPair;
    static deserializeBinaryFromReader(message: KvPair, reader: jspb.BinaryReader): KvPair;
  }

  export namespace KvPair {
    export type AsObject = {
      key: string,
      value: string,
    }
  }

  export enum SpecCase {
    SPEC_NOT_SET = 0,
    KEY = 1,
    KV_PAIR = 2,
  }
}

