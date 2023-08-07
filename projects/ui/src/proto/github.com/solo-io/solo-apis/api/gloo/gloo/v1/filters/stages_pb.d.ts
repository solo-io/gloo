/* eslint-disable */
// package: filters.gloo.solo.io
// file: github.com/solo-io/solo-apis/api/gloo/gloo/v1/filters/stages.proto

import * as jspb from "google-protobuf";
import * as extproto_ext_pb from "../../../../../../../../extproto/ext_pb";

export class FilterStage extends jspb.Message {
  getStage(): FilterStage.StageMap[keyof FilterStage.StageMap];
  setStage(value: FilterStage.StageMap[keyof FilterStage.StageMap]): void;

  getPredicate(): FilterStage.PredicateMap[keyof FilterStage.PredicateMap];
  setPredicate(value: FilterStage.PredicateMap[keyof FilterStage.PredicateMap]): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FilterStage.AsObject;
  static toObject(includeInstance: boolean, msg: FilterStage): FilterStage.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: FilterStage, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FilterStage;
  static deserializeBinaryFromReader(message: FilterStage, reader: jspb.BinaryReader): FilterStage;
}

export namespace FilterStage {
  export type AsObject = {
    stage: FilterStage.StageMap[keyof FilterStage.StageMap],
    predicate: FilterStage.PredicateMap[keyof FilterStage.PredicateMap],
  }

  export interface StageMap {
    FAULTSTAGE: 0;
    CORSSTAGE: 1;
    WAFSTAGE: 2;
    AUTHNSTAGE: 3;
    AUTHZSTAGE: 4;
    RATELIMITSTAGE: 5;
    ACCEPTEDSTAGE: 6;
    OUTAUTHSTAGE: 7;
    ROUTESTAGE: 8;
  }

  export const Stage: StageMap;

  export interface PredicateMap {
    DURING: 0;
    BEFORE: 1;
    AFTER: 2;
  }

  export const Predicate: PredicateMap;
}
