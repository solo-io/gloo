// package: gloo.solo.io
// file: github.com/solo-io/gloo/projects/gloo/api/v1/load_balancer.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../../../../../gogoproto/gogo_pb";
import * as google_protobuf_duration_pb from "google-protobuf/google/protobuf/duration_pb";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";

export class LoadBalancerConfig extends jspb.Message {
  hasHealthyPanicThreshold(): boolean;
  clearHealthyPanicThreshold(): void;
  getHealthyPanicThreshold(): google_protobuf_wrappers_pb.DoubleValue | undefined;
  setHealthyPanicThreshold(value?: google_protobuf_wrappers_pb.DoubleValue): void;

  hasUpdateMergeWindow(): boolean;
  clearUpdateMergeWindow(): void;
  getUpdateMergeWindow(): google_protobuf_duration_pb.Duration | undefined;
  setUpdateMergeWindow(value?: google_protobuf_duration_pb.Duration): void;

  hasRoundRobin(): boolean;
  clearRoundRobin(): void;
  getRoundRobin(): LoadBalancerConfig.RoundRobin | undefined;
  setRoundRobin(value?: LoadBalancerConfig.RoundRobin): void;

  hasLeastRequest(): boolean;
  clearLeastRequest(): void;
  getLeastRequest(): LoadBalancerConfig.LeastRequest | undefined;
  setLeastRequest(value?: LoadBalancerConfig.LeastRequest): void;

  hasRandom(): boolean;
  clearRandom(): void;
  getRandom(): LoadBalancerConfig.Random | undefined;
  setRandom(value?: LoadBalancerConfig.Random): void;

  getTypeCase(): LoadBalancerConfig.TypeCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LoadBalancerConfig.AsObject;
  static toObject(includeInstance: boolean, msg: LoadBalancerConfig): LoadBalancerConfig.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: LoadBalancerConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LoadBalancerConfig;
  static deserializeBinaryFromReader(message: LoadBalancerConfig, reader: jspb.BinaryReader): LoadBalancerConfig;
}

export namespace LoadBalancerConfig {
  export type AsObject = {
    healthyPanicThreshold?: google_protobuf_wrappers_pb.DoubleValue.AsObject,
    updateMergeWindow?: google_protobuf_duration_pb.Duration.AsObject,
    roundRobin?: LoadBalancerConfig.RoundRobin.AsObject,
    leastRequest?: LoadBalancerConfig.LeastRequest.AsObject,
    random?: LoadBalancerConfig.Random.AsObject,
  }

  export class RoundRobin extends jspb.Message {
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): RoundRobin.AsObject;
    static toObject(includeInstance: boolean, msg: RoundRobin): RoundRobin.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: RoundRobin, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): RoundRobin;
    static deserializeBinaryFromReader(message: RoundRobin, reader: jspb.BinaryReader): RoundRobin;
  }

  export namespace RoundRobin {
    export type AsObject = {
    }
  }

  export class LeastRequest extends jspb.Message {
    getChoiceCount(): number;
    setChoiceCount(value: number): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): LeastRequest.AsObject;
    static toObject(includeInstance: boolean, msg: LeastRequest): LeastRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: LeastRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): LeastRequest;
    static deserializeBinaryFromReader(message: LeastRequest, reader: jspb.BinaryReader): LeastRequest;
  }

  export namespace LeastRequest {
    export type AsObject = {
      choiceCount: number,
    }
  }

  export class Random extends jspb.Message {
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Random.AsObject;
    static toObject(includeInstance: boolean, msg: Random): Random.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Random, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Random;
    static deserializeBinaryFromReader(message: Random, reader: jspb.BinaryReader): Random;
  }

  export namespace Random {
    export type AsObject = {
    }
  }

  export enum TypeCase {
    TYPE_NOT_SET = 0,
    ROUND_ROBIN = 3,
    LEAST_REQUEST = 4,
    RANDOM = 5,
  }
}

