// package: envoy.type
// file: github.com/solo-io/solo-kit/api/external/envoy/type/percent.proto

import * as jspb from "google-protobuf";
import * as validate_validate_pb from "../../../../../../../validate/validate_pb";
import * as gogoproto_gogo_pb from "../../../../../../../gogoproto/gogo_pb";

export class Percent extends jspb.Message {
  getValue(): number;
  setValue(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Percent.AsObject;
  static toObject(includeInstance: boolean, msg: Percent): Percent.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Percent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Percent;
  static deserializeBinaryFromReader(message: Percent, reader: jspb.BinaryReader): Percent;
}

export namespace Percent {
  export type AsObject = {
    value: number,
  }
}

export class FractionalPercent extends jspb.Message {
  getNumerator(): number;
  setNumerator(value: number): void;

  getDenominator(): FractionalPercent.DenominatorTypeMap[keyof FractionalPercent.DenominatorTypeMap];
  setDenominator(value: FractionalPercent.DenominatorTypeMap[keyof FractionalPercent.DenominatorTypeMap]): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FractionalPercent.AsObject;
  static toObject(includeInstance: boolean, msg: FractionalPercent): FractionalPercent.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: FractionalPercent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FractionalPercent;
  static deserializeBinaryFromReader(message: FractionalPercent, reader: jspb.BinaryReader): FractionalPercent;
}

export namespace FractionalPercent {
  export type AsObject = {
    numerator: number,
    denominator: FractionalPercent.DenominatorTypeMap[keyof FractionalPercent.DenominatorTypeMap],
  }

  export interface DenominatorTypeMap {
    HUNDRED: 0;
    TEN_THOUSAND: 1;
    MILLION: 2;
  }

  export const DenominatorType: DenominatorTypeMap;
}

