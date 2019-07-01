// package: validate
// file: github.com/solo-io/solo-kit/api/external/validate/validate.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_descriptor_pb from "google-protobuf/google/protobuf/descriptor_pb";
import * as google_protobuf_duration_pb from "google-protobuf/google/protobuf/duration_pb";
import * as google_protobuf_timestamp_pb from "google-protobuf/google/protobuf/timestamp_pb";

export class FieldRules extends jspb.Message {
  hasFloat(): boolean;
  clearFloat(): void;
  getFloat(): FloatRules | undefined;
  setFloat(value?: FloatRules): void;

  hasDouble(): boolean;
  clearDouble(): void;
  getDouble(): DoubleRules | undefined;
  setDouble(value?: DoubleRules): void;

  hasInt32(): boolean;
  clearInt32(): void;
  getInt32(): Int32Rules | undefined;
  setInt32(value?: Int32Rules): void;

  hasInt64(): boolean;
  clearInt64(): void;
  getInt64(): Int64Rules | undefined;
  setInt64(value?: Int64Rules): void;

  hasUint32(): boolean;
  clearUint32(): void;
  getUint32(): UInt32Rules | undefined;
  setUint32(value?: UInt32Rules): void;

  hasUint64(): boolean;
  clearUint64(): void;
  getUint64(): UInt64Rules | undefined;
  setUint64(value?: UInt64Rules): void;

  hasSint32(): boolean;
  clearSint32(): void;
  getSint32(): SInt32Rules | undefined;
  setSint32(value?: SInt32Rules): void;

  hasSint64(): boolean;
  clearSint64(): void;
  getSint64(): SInt64Rules | undefined;
  setSint64(value?: SInt64Rules): void;

  hasFixed32(): boolean;
  clearFixed32(): void;
  getFixed32(): Fixed32Rules | undefined;
  setFixed32(value?: Fixed32Rules): void;

  hasFixed64(): boolean;
  clearFixed64(): void;
  getFixed64(): Fixed64Rules | undefined;
  setFixed64(value?: Fixed64Rules): void;

  hasSfixed32(): boolean;
  clearSfixed32(): void;
  getSfixed32(): SFixed32Rules | undefined;
  setSfixed32(value?: SFixed32Rules): void;

  hasSfixed64(): boolean;
  clearSfixed64(): void;
  getSfixed64(): SFixed64Rules | undefined;
  setSfixed64(value?: SFixed64Rules): void;

  hasBool(): boolean;
  clearBool(): void;
  getBool(): BoolRules | undefined;
  setBool(value?: BoolRules): void;

  hasString(): boolean;
  clearString(): void;
  getString(): StringRules | undefined;
  setString(value?: StringRules): void;

  hasBytes(): boolean;
  clearBytes(): void;
  getBytes(): BytesRules | undefined;
  setBytes(value?: BytesRules): void;

  hasEnum(): boolean;
  clearEnum(): void;
  getEnum(): EnumRules | undefined;
  setEnum(value?: EnumRules): void;

  hasMessage(): boolean;
  clearMessage(): void;
  getMessage(): MessageRules | undefined;
  setMessage(value?: MessageRules): void;

  hasRepeated(): boolean;
  clearRepeated(): void;
  getRepeated(): RepeatedRules | undefined;
  setRepeated(value?: RepeatedRules): void;

  hasMap(): boolean;
  clearMap(): void;
  getMap(): MapRules | undefined;
  setMap(value?: MapRules): void;

  hasAny(): boolean;
  clearAny(): void;
  getAny(): AnyRules | undefined;
  setAny(value?: AnyRules): void;

  hasDuration(): boolean;
  clearDuration(): void;
  getDuration(): DurationRules | undefined;
  setDuration(value?: DurationRules): void;

  hasTimestamp(): boolean;
  clearTimestamp(): void;
  getTimestamp(): TimestampRules | undefined;
  setTimestamp(value?: TimestampRules): void;

  getTypeCase(): FieldRules.TypeCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FieldRules.AsObject;
  static toObject(includeInstance: boolean, msg: FieldRules): FieldRules.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: FieldRules, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FieldRules;
  static deserializeBinaryFromReader(message: FieldRules, reader: jspb.BinaryReader): FieldRules;
}

export namespace FieldRules {
  export type AsObject = {
    pb_float?: FloatRules.AsObject,
    pb_double?: DoubleRules.AsObject,
    int32?: Int32Rules.AsObject,
    int64?: Int64Rules.AsObject,
    uint32?: UInt32Rules.AsObject,
    uint64?: UInt64Rules.AsObject,
    sint32?: SInt32Rules.AsObject,
    sint64?: SInt64Rules.AsObject,
    fixed32?: Fixed32Rules.AsObject,
    fixed64?: Fixed64Rules.AsObject,
    sfixed32?: SFixed32Rules.AsObject,
    sfixed64?: SFixed64Rules.AsObject,
    bool?: BoolRules.AsObject,
    string?: StringRules.AsObject,
    bytes?: BytesRules.AsObject,
    pb_enum?: EnumRules.AsObject,
    message?: MessageRules.AsObject,
    repeated?: RepeatedRules.AsObject,
    map?: MapRules.AsObject,
    any?: AnyRules.AsObject,
    duration?: DurationRules.AsObject,
    timestamp?: TimestampRules.AsObject,
  }

  export enum TypeCase {
    TYPE_NOT_SET = 0,
    FLOAT = 1,
    DOUBLE = 2,
    INT32 = 3,
    INT64 = 4,
    UINT32 = 5,
    UINT64 = 6,
    SINT32 = 7,
    SINT64 = 8,
    FIXED32 = 9,
    FIXED64 = 10,
    SFIXED32 = 11,
    SFIXED64 = 12,
    BOOL = 13,
    STRING = 14,
    BYTES = 15,
    ENUM = 16,
    MESSAGE = 17,
    REPEATED = 18,
    MAP = 19,
    ANY = 20,
    DURATION = 21,
    TIMESTAMP = 22,
  }
}

export class FloatRules extends jspb.Message {
  hasConst(): boolean;
  clearConst(): void;
  getConst(): number | undefined;
  setConst(value: number): void;

  hasLt(): boolean;
  clearLt(): void;
  getLt(): number | undefined;
  setLt(value: number): void;

  hasLte(): boolean;
  clearLte(): void;
  getLte(): number | undefined;
  setLte(value: number): void;

  hasGt(): boolean;
  clearGt(): void;
  getGt(): number | undefined;
  setGt(value: number): void;

  hasGte(): boolean;
  clearGte(): void;
  getGte(): number | undefined;
  setGte(value: number): void;

  clearInList(): void;
  getInList(): Array<number>;
  setInList(value: Array<number>): void;
  addIn(value: number, index?: number): number;

  clearNotInList(): void;
  getNotInList(): Array<number>;
  setNotInList(value: Array<number>): void;
  addNotIn(value: number, index?: number): number;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FloatRules.AsObject;
  static toObject(includeInstance: boolean, msg: FloatRules): FloatRules.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: FloatRules, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FloatRules;
  static deserializeBinaryFromReader(message: FloatRules, reader: jspb.BinaryReader): FloatRules;
}

export namespace FloatRules {
  export type AsObject = {
    pb_const?: number,
    lt?: number,
    lte?: number,
    gt?: number,
    gte?: number,
    inList: Array<number>,
    notInList: Array<number>,
  }
}

export class DoubleRules extends jspb.Message {
  hasConst(): boolean;
  clearConst(): void;
  getConst(): number | undefined;
  setConst(value: number): void;

  hasLt(): boolean;
  clearLt(): void;
  getLt(): number | undefined;
  setLt(value: number): void;

  hasLte(): boolean;
  clearLte(): void;
  getLte(): number | undefined;
  setLte(value: number): void;

  hasGt(): boolean;
  clearGt(): void;
  getGt(): number | undefined;
  setGt(value: number): void;

  hasGte(): boolean;
  clearGte(): void;
  getGte(): number | undefined;
  setGte(value: number): void;

  clearInList(): void;
  getInList(): Array<number>;
  setInList(value: Array<number>): void;
  addIn(value: number, index?: number): number;

  clearNotInList(): void;
  getNotInList(): Array<number>;
  setNotInList(value: Array<number>): void;
  addNotIn(value: number, index?: number): number;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DoubleRules.AsObject;
  static toObject(includeInstance: boolean, msg: DoubleRules): DoubleRules.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DoubleRules, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DoubleRules;
  static deserializeBinaryFromReader(message: DoubleRules, reader: jspb.BinaryReader): DoubleRules;
}

export namespace DoubleRules {
  export type AsObject = {
    pb_const?: number,
    lt?: number,
    lte?: number,
    gt?: number,
    gte?: number,
    inList: Array<number>,
    notInList: Array<number>,
  }
}

export class Int32Rules extends jspb.Message {
  hasConst(): boolean;
  clearConst(): void;
  getConst(): number | undefined;
  setConst(value: number): void;

  hasLt(): boolean;
  clearLt(): void;
  getLt(): number | undefined;
  setLt(value: number): void;

  hasLte(): boolean;
  clearLte(): void;
  getLte(): number | undefined;
  setLte(value: number): void;

  hasGt(): boolean;
  clearGt(): void;
  getGt(): number | undefined;
  setGt(value: number): void;

  hasGte(): boolean;
  clearGte(): void;
  getGte(): number | undefined;
  setGte(value: number): void;

  clearInList(): void;
  getInList(): Array<number>;
  setInList(value: Array<number>): void;
  addIn(value: number, index?: number): number;

  clearNotInList(): void;
  getNotInList(): Array<number>;
  setNotInList(value: Array<number>): void;
  addNotIn(value: number, index?: number): number;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Int32Rules.AsObject;
  static toObject(includeInstance: boolean, msg: Int32Rules): Int32Rules.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Int32Rules, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Int32Rules;
  static deserializeBinaryFromReader(message: Int32Rules, reader: jspb.BinaryReader): Int32Rules;
}

export namespace Int32Rules {
  export type AsObject = {
    pb_const?: number,
    lt?: number,
    lte?: number,
    gt?: number,
    gte?: number,
    inList: Array<number>,
    notInList: Array<number>,
  }
}

export class Int64Rules extends jspb.Message {
  hasConst(): boolean;
  clearConst(): void;
  getConst(): number | undefined;
  setConst(value: number): void;

  hasLt(): boolean;
  clearLt(): void;
  getLt(): number | undefined;
  setLt(value: number): void;

  hasLte(): boolean;
  clearLte(): void;
  getLte(): number | undefined;
  setLte(value: number): void;

  hasGt(): boolean;
  clearGt(): void;
  getGt(): number | undefined;
  setGt(value: number): void;

  hasGte(): boolean;
  clearGte(): void;
  getGte(): number | undefined;
  setGte(value: number): void;

  clearInList(): void;
  getInList(): Array<number>;
  setInList(value: Array<number>): void;
  addIn(value: number, index?: number): number;

  clearNotInList(): void;
  getNotInList(): Array<number>;
  setNotInList(value: Array<number>): void;
  addNotIn(value: number, index?: number): number;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Int64Rules.AsObject;
  static toObject(includeInstance: boolean, msg: Int64Rules): Int64Rules.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Int64Rules, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Int64Rules;
  static deserializeBinaryFromReader(message: Int64Rules, reader: jspb.BinaryReader): Int64Rules;
}

export namespace Int64Rules {
  export type AsObject = {
    pb_const?: number,
    lt?: number,
    lte?: number,
    gt?: number,
    gte?: number,
    inList: Array<number>,
    notInList: Array<number>,
  }
}

export class UInt32Rules extends jspb.Message {
  hasConst(): boolean;
  clearConst(): void;
  getConst(): number | undefined;
  setConst(value: number): void;

  hasLt(): boolean;
  clearLt(): void;
  getLt(): number | undefined;
  setLt(value: number): void;

  hasLte(): boolean;
  clearLte(): void;
  getLte(): number | undefined;
  setLte(value: number): void;

  hasGt(): boolean;
  clearGt(): void;
  getGt(): number | undefined;
  setGt(value: number): void;

  hasGte(): boolean;
  clearGte(): void;
  getGte(): number | undefined;
  setGte(value: number): void;

  clearInList(): void;
  getInList(): Array<number>;
  setInList(value: Array<number>): void;
  addIn(value: number, index?: number): number;

  clearNotInList(): void;
  getNotInList(): Array<number>;
  setNotInList(value: Array<number>): void;
  addNotIn(value: number, index?: number): number;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UInt32Rules.AsObject;
  static toObject(includeInstance: boolean, msg: UInt32Rules): UInt32Rules.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UInt32Rules, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UInt32Rules;
  static deserializeBinaryFromReader(message: UInt32Rules, reader: jspb.BinaryReader): UInt32Rules;
}

export namespace UInt32Rules {
  export type AsObject = {
    pb_const?: number,
    lt?: number,
    lte?: number,
    gt?: number,
    gte?: number,
    inList: Array<number>,
    notInList: Array<number>,
  }
}

export class UInt64Rules extends jspb.Message {
  hasConst(): boolean;
  clearConst(): void;
  getConst(): number | undefined;
  setConst(value: number): void;

  hasLt(): boolean;
  clearLt(): void;
  getLt(): number | undefined;
  setLt(value: number): void;

  hasLte(): boolean;
  clearLte(): void;
  getLte(): number | undefined;
  setLte(value: number): void;

  hasGt(): boolean;
  clearGt(): void;
  getGt(): number | undefined;
  setGt(value: number): void;

  hasGte(): boolean;
  clearGte(): void;
  getGte(): number | undefined;
  setGte(value: number): void;

  clearInList(): void;
  getInList(): Array<number>;
  setInList(value: Array<number>): void;
  addIn(value: number, index?: number): number;

  clearNotInList(): void;
  getNotInList(): Array<number>;
  setNotInList(value: Array<number>): void;
  addNotIn(value: number, index?: number): number;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UInt64Rules.AsObject;
  static toObject(includeInstance: boolean, msg: UInt64Rules): UInt64Rules.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UInt64Rules, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UInt64Rules;
  static deserializeBinaryFromReader(message: UInt64Rules, reader: jspb.BinaryReader): UInt64Rules;
}

export namespace UInt64Rules {
  export type AsObject = {
    pb_const?: number,
    lt?: number,
    lte?: number,
    gt?: number,
    gte?: number,
    inList: Array<number>,
    notInList: Array<number>,
  }
}

export class SInt32Rules extends jspb.Message {
  hasConst(): boolean;
  clearConst(): void;
  getConst(): number | undefined;
  setConst(value: number): void;

  hasLt(): boolean;
  clearLt(): void;
  getLt(): number | undefined;
  setLt(value: number): void;

  hasLte(): boolean;
  clearLte(): void;
  getLte(): number | undefined;
  setLte(value: number): void;

  hasGt(): boolean;
  clearGt(): void;
  getGt(): number | undefined;
  setGt(value: number): void;

  hasGte(): boolean;
  clearGte(): void;
  getGte(): number | undefined;
  setGte(value: number): void;

  clearInList(): void;
  getInList(): Array<number>;
  setInList(value: Array<number>): void;
  addIn(value: number, index?: number): number;

  clearNotInList(): void;
  getNotInList(): Array<number>;
  setNotInList(value: Array<number>): void;
  addNotIn(value: number, index?: number): number;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SInt32Rules.AsObject;
  static toObject(includeInstance: boolean, msg: SInt32Rules): SInt32Rules.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: SInt32Rules, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SInt32Rules;
  static deserializeBinaryFromReader(message: SInt32Rules, reader: jspb.BinaryReader): SInt32Rules;
}

export namespace SInt32Rules {
  export type AsObject = {
    pb_const?: number,
    lt?: number,
    lte?: number,
    gt?: number,
    gte?: number,
    inList: Array<number>,
    notInList: Array<number>,
  }
}

export class SInt64Rules extends jspb.Message {
  hasConst(): boolean;
  clearConst(): void;
  getConst(): number | undefined;
  setConst(value: number): void;

  hasLt(): boolean;
  clearLt(): void;
  getLt(): number | undefined;
  setLt(value: number): void;

  hasLte(): boolean;
  clearLte(): void;
  getLte(): number | undefined;
  setLte(value: number): void;

  hasGt(): boolean;
  clearGt(): void;
  getGt(): number | undefined;
  setGt(value: number): void;

  hasGte(): boolean;
  clearGte(): void;
  getGte(): number | undefined;
  setGte(value: number): void;

  clearInList(): void;
  getInList(): Array<number>;
  setInList(value: Array<number>): void;
  addIn(value: number, index?: number): number;

  clearNotInList(): void;
  getNotInList(): Array<number>;
  setNotInList(value: Array<number>): void;
  addNotIn(value: number, index?: number): number;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SInt64Rules.AsObject;
  static toObject(includeInstance: boolean, msg: SInt64Rules): SInt64Rules.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: SInt64Rules, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SInt64Rules;
  static deserializeBinaryFromReader(message: SInt64Rules, reader: jspb.BinaryReader): SInt64Rules;
}

export namespace SInt64Rules {
  export type AsObject = {
    pb_const?: number,
    lt?: number,
    lte?: number,
    gt?: number,
    gte?: number,
    inList: Array<number>,
    notInList: Array<number>,
  }
}

export class Fixed32Rules extends jspb.Message {
  hasConst(): boolean;
  clearConst(): void;
  getConst(): number | undefined;
  setConst(value: number): void;

  hasLt(): boolean;
  clearLt(): void;
  getLt(): number | undefined;
  setLt(value: number): void;

  hasLte(): boolean;
  clearLte(): void;
  getLte(): number | undefined;
  setLte(value: number): void;

  hasGt(): boolean;
  clearGt(): void;
  getGt(): number | undefined;
  setGt(value: number): void;

  hasGte(): boolean;
  clearGte(): void;
  getGte(): number | undefined;
  setGte(value: number): void;

  clearInList(): void;
  getInList(): Array<number>;
  setInList(value: Array<number>): void;
  addIn(value: number, index?: number): number;

  clearNotInList(): void;
  getNotInList(): Array<number>;
  setNotInList(value: Array<number>): void;
  addNotIn(value: number, index?: number): number;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Fixed32Rules.AsObject;
  static toObject(includeInstance: boolean, msg: Fixed32Rules): Fixed32Rules.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Fixed32Rules, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Fixed32Rules;
  static deserializeBinaryFromReader(message: Fixed32Rules, reader: jspb.BinaryReader): Fixed32Rules;
}

export namespace Fixed32Rules {
  export type AsObject = {
    pb_const?: number,
    lt?: number,
    lte?: number,
    gt?: number,
    gte?: number,
    inList: Array<number>,
    notInList: Array<number>,
  }
}

export class Fixed64Rules extends jspb.Message {
  hasConst(): boolean;
  clearConst(): void;
  getConst(): number | undefined;
  setConst(value: number): void;

  hasLt(): boolean;
  clearLt(): void;
  getLt(): number | undefined;
  setLt(value: number): void;

  hasLte(): boolean;
  clearLte(): void;
  getLte(): number | undefined;
  setLte(value: number): void;

  hasGt(): boolean;
  clearGt(): void;
  getGt(): number | undefined;
  setGt(value: number): void;

  hasGte(): boolean;
  clearGte(): void;
  getGte(): number | undefined;
  setGte(value: number): void;

  clearInList(): void;
  getInList(): Array<number>;
  setInList(value: Array<number>): void;
  addIn(value: number, index?: number): number;

  clearNotInList(): void;
  getNotInList(): Array<number>;
  setNotInList(value: Array<number>): void;
  addNotIn(value: number, index?: number): number;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Fixed64Rules.AsObject;
  static toObject(includeInstance: boolean, msg: Fixed64Rules): Fixed64Rules.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Fixed64Rules, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Fixed64Rules;
  static deserializeBinaryFromReader(message: Fixed64Rules, reader: jspb.BinaryReader): Fixed64Rules;
}

export namespace Fixed64Rules {
  export type AsObject = {
    pb_const?: number,
    lt?: number,
    lte?: number,
    gt?: number,
    gte?: number,
    inList: Array<number>,
    notInList: Array<number>,
  }
}

export class SFixed32Rules extends jspb.Message {
  hasConst(): boolean;
  clearConst(): void;
  getConst(): number | undefined;
  setConst(value: number): void;

  hasLt(): boolean;
  clearLt(): void;
  getLt(): number | undefined;
  setLt(value: number): void;

  hasLte(): boolean;
  clearLte(): void;
  getLte(): number | undefined;
  setLte(value: number): void;

  hasGt(): boolean;
  clearGt(): void;
  getGt(): number | undefined;
  setGt(value: number): void;

  hasGte(): boolean;
  clearGte(): void;
  getGte(): number | undefined;
  setGte(value: number): void;

  clearInList(): void;
  getInList(): Array<number>;
  setInList(value: Array<number>): void;
  addIn(value: number, index?: number): number;

  clearNotInList(): void;
  getNotInList(): Array<number>;
  setNotInList(value: Array<number>): void;
  addNotIn(value: number, index?: number): number;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SFixed32Rules.AsObject;
  static toObject(includeInstance: boolean, msg: SFixed32Rules): SFixed32Rules.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: SFixed32Rules, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SFixed32Rules;
  static deserializeBinaryFromReader(message: SFixed32Rules, reader: jspb.BinaryReader): SFixed32Rules;
}

export namespace SFixed32Rules {
  export type AsObject = {
    pb_const?: number,
    lt?: number,
    lte?: number,
    gt?: number,
    gte?: number,
    inList: Array<number>,
    notInList: Array<number>,
  }
}

export class SFixed64Rules extends jspb.Message {
  hasConst(): boolean;
  clearConst(): void;
  getConst(): number | undefined;
  setConst(value: number): void;

  hasLt(): boolean;
  clearLt(): void;
  getLt(): number | undefined;
  setLt(value: number): void;

  hasLte(): boolean;
  clearLte(): void;
  getLte(): number | undefined;
  setLte(value: number): void;

  hasGt(): boolean;
  clearGt(): void;
  getGt(): number | undefined;
  setGt(value: number): void;

  hasGte(): boolean;
  clearGte(): void;
  getGte(): number | undefined;
  setGte(value: number): void;

  clearInList(): void;
  getInList(): Array<number>;
  setInList(value: Array<number>): void;
  addIn(value: number, index?: number): number;

  clearNotInList(): void;
  getNotInList(): Array<number>;
  setNotInList(value: Array<number>): void;
  addNotIn(value: number, index?: number): number;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SFixed64Rules.AsObject;
  static toObject(includeInstance: boolean, msg: SFixed64Rules): SFixed64Rules.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: SFixed64Rules, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SFixed64Rules;
  static deserializeBinaryFromReader(message: SFixed64Rules, reader: jspb.BinaryReader): SFixed64Rules;
}

export namespace SFixed64Rules {
  export type AsObject = {
    pb_const?: number,
    lt?: number,
    lte?: number,
    gt?: number,
    gte?: number,
    inList: Array<number>,
    notInList: Array<number>,
  }
}

export class BoolRules extends jspb.Message {
  hasConst(): boolean;
  clearConst(): void;
  getConst(): boolean | undefined;
  setConst(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): BoolRules.AsObject;
  static toObject(includeInstance: boolean, msg: BoolRules): BoolRules.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: BoolRules, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): BoolRules;
  static deserializeBinaryFromReader(message: BoolRules, reader: jspb.BinaryReader): BoolRules;
}

export namespace BoolRules {
  export type AsObject = {
    pb_const?: boolean,
  }
}

export class StringRules extends jspb.Message {
  hasConst(): boolean;
  clearConst(): void;
  getConst(): string | undefined;
  setConst(value: string): void;

  hasMinLen(): boolean;
  clearMinLen(): void;
  getMinLen(): number | undefined;
  setMinLen(value: number): void;

  hasMaxLen(): boolean;
  clearMaxLen(): void;
  getMaxLen(): number | undefined;
  setMaxLen(value: number): void;

  hasMinBytes(): boolean;
  clearMinBytes(): void;
  getMinBytes(): number | undefined;
  setMinBytes(value: number): void;

  hasMaxBytes(): boolean;
  clearMaxBytes(): void;
  getMaxBytes(): number | undefined;
  setMaxBytes(value: number): void;

  hasPattern(): boolean;
  clearPattern(): void;
  getPattern(): string | undefined;
  setPattern(value: string): void;

  hasPrefix(): boolean;
  clearPrefix(): void;
  getPrefix(): string | undefined;
  setPrefix(value: string): void;

  hasSuffix(): boolean;
  clearSuffix(): void;
  getSuffix(): string | undefined;
  setSuffix(value: string): void;

  hasContains(): boolean;
  clearContains(): void;
  getContains(): string | undefined;
  setContains(value: string): void;

  clearInList(): void;
  getInList(): Array<string>;
  setInList(value: Array<string>): void;
  addIn(value: string, index?: number): string;

  clearNotInList(): void;
  getNotInList(): Array<string>;
  setNotInList(value: Array<string>): void;
  addNotIn(value: string, index?: number): string;

  hasEmail(): boolean;
  clearEmail(): void;
  getEmail(): boolean | undefined;
  setEmail(value: boolean): void;

  hasHostname(): boolean;
  clearHostname(): void;
  getHostname(): boolean | undefined;
  setHostname(value: boolean): void;

  hasIp(): boolean;
  clearIp(): void;
  getIp(): boolean | undefined;
  setIp(value: boolean): void;

  hasIpv4(): boolean;
  clearIpv4(): void;
  getIpv4(): boolean | undefined;
  setIpv4(value: boolean): void;

  hasIpv6(): boolean;
  clearIpv6(): void;
  getIpv6(): boolean | undefined;
  setIpv6(value: boolean): void;

  hasUri(): boolean;
  clearUri(): void;
  getUri(): boolean | undefined;
  setUri(value: boolean): void;

  hasUriRef(): boolean;
  clearUriRef(): void;
  getUriRef(): boolean | undefined;
  setUriRef(value: boolean): void;

  getWellKnownCase(): StringRules.WellKnownCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StringRules.AsObject;
  static toObject(includeInstance: boolean, msg: StringRules): StringRules.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: StringRules, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StringRules;
  static deserializeBinaryFromReader(message: StringRules, reader: jspb.BinaryReader): StringRules;
}

export namespace StringRules {
  export type AsObject = {
    pb_const?: string,
    minLen?: number,
    maxLen?: number,
    minBytes?: number,
    maxBytes?: number,
    pattern?: string,
    prefix?: string,
    suffix?: string,
    contains?: string,
    inList: Array<string>,
    notInList: Array<string>,
    email?: boolean,
    hostname?: boolean,
    ip?: boolean,
    ipv4?: boolean,
    ipv6?: boolean,
    uri?: boolean,
    uriRef?: boolean,
  }

  export enum WellKnownCase {
    WELL_KNOWN_NOT_SET = 0,
    EMAIL = 12,
    HOSTNAME = 13,
    IP = 14,
    IPV4 = 15,
    IPV6 = 16,
    URI = 17,
    URI_REF = 18,
  }
}

export class BytesRules extends jspb.Message {
  hasConst(): boolean;
  clearConst(): void;
  getConst(): Uint8Array | string;
  getConst_asU8(): Uint8Array;
  getConst_asB64(): string;
  setConst(value: Uint8Array | string): void;

  hasMinLen(): boolean;
  clearMinLen(): void;
  getMinLen(): number | undefined;
  setMinLen(value: number): void;

  hasMaxLen(): boolean;
  clearMaxLen(): void;
  getMaxLen(): number | undefined;
  setMaxLen(value: number): void;

  hasPattern(): boolean;
  clearPattern(): void;
  getPattern(): string | undefined;
  setPattern(value: string): void;

  hasPrefix(): boolean;
  clearPrefix(): void;
  getPrefix(): Uint8Array | string;
  getPrefix_asU8(): Uint8Array;
  getPrefix_asB64(): string;
  setPrefix(value: Uint8Array | string): void;

  hasSuffix(): boolean;
  clearSuffix(): void;
  getSuffix(): Uint8Array | string;
  getSuffix_asU8(): Uint8Array;
  getSuffix_asB64(): string;
  setSuffix(value: Uint8Array | string): void;

  hasContains(): boolean;
  clearContains(): void;
  getContains(): Uint8Array | string;
  getContains_asU8(): Uint8Array;
  getContains_asB64(): string;
  setContains(value: Uint8Array | string): void;

  clearInList(): void;
  getInList(): Array<Uint8Array | string>;
  getInList_asU8(): Array<Uint8Array>;
  getInList_asB64(): Array<string>;
  setInList(value: Array<Uint8Array | string>): void;
  addIn(value: Uint8Array | string, index?: number): Uint8Array | string;

  clearNotInList(): void;
  getNotInList(): Array<Uint8Array | string>;
  getNotInList_asU8(): Array<Uint8Array>;
  getNotInList_asB64(): Array<string>;
  setNotInList(value: Array<Uint8Array | string>): void;
  addNotIn(value: Uint8Array | string, index?: number): Uint8Array | string;

  hasIp(): boolean;
  clearIp(): void;
  getIp(): boolean | undefined;
  setIp(value: boolean): void;

  hasIpv4(): boolean;
  clearIpv4(): void;
  getIpv4(): boolean | undefined;
  setIpv4(value: boolean): void;

  hasIpv6(): boolean;
  clearIpv6(): void;
  getIpv6(): boolean | undefined;
  setIpv6(value: boolean): void;

  getWellKnownCase(): BytesRules.WellKnownCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): BytesRules.AsObject;
  static toObject(includeInstance: boolean, msg: BytesRules): BytesRules.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: BytesRules, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): BytesRules;
  static deserializeBinaryFromReader(message: BytesRules, reader: jspb.BinaryReader): BytesRules;
}

export namespace BytesRules {
  export type AsObject = {
    const: Uint8Array | string,
    minLen?: number,
    maxLen?: number,
    pattern?: string,
    prefix: Uint8Array | string,
    suffix: Uint8Array | string,
    contains: Uint8Array | string,
    inList: Array<Uint8Array | string>,
    notInList: Array<Uint8Array | string>,
    ip?: boolean,
    ipv4?: boolean,
    ipv6?: boolean,
  }

  export enum WellKnownCase {
    WELL_KNOWN_NOT_SET = 0,
    IP = 10,
    IPV4 = 11,
    IPV6 = 12,
  }
}

export class EnumRules extends jspb.Message {
  hasConst(): boolean;
  clearConst(): void;
  getConst(): number | undefined;
  setConst(value: number): void;

  hasDefinedOnly(): boolean;
  clearDefinedOnly(): void;
  getDefinedOnly(): boolean | undefined;
  setDefinedOnly(value: boolean): void;

  clearInList(): void;
  getInList(): Array<number>;
  setInList(value: Array<number>): void;
  addIn(value: number, index?: number): number;

  clearNotInList(): void;
  getNotInList(): Array<number>;
  setNotInList(value: Array<number>): void;
  addNotIn(value: number, index?: number): number;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EnumRules.AsObject;
  static toObject(includeInstance: boolean, msg: EnumRules): EnumRules.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: EnumRules, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EnumRules;
  static deserializeBinaryFromReader(message: EnumRules, reader: jspb.BinaryReader): EnumRules;
}

export namespace EnumRules {
  export type AsObject = {
    pb_const?: number,
    definedOnly?: boolean,
    inList: Array<number>,
    notInList: Array<number>,
  }
}

export class MessageRules extends jspb.Message {
  hasSkip(): boolean;
  clearSkip(): void;
  getSkip(): boolean | undefined;
  setSkip(value: boolean): void;

  hasRequired(): boolean;
  clearRequired(): void;
  getRequired(): boolean | undefined;
  setRequired(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MessageRules.AsObject;
  static toObject(includeInstance: boolean, msg: MessageRules): MessageRules.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: MessageRules, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MessageRules;
  static deserializeBinaryFromReader(message: MessageRules, reader: jspb.BinaryReader): MessageRules;
}

export namespace MessageRules {
  export type AsObject = {
    skip?: boolean,
    required?: boolean,
  }
}

export class RepeatedRules extends jspb.Message {
  hasMinItems(): boolean;
  clearMinItems(): void;
  getMinItems(): number | undefined;
  setMinItems(value: number): void;

  hasMaxItems(): boolean;
  clearMaxItems(): void;
  getMaxItems(): number | undefined;
  setMaxItems(value: number): void;

  hasUnique(): boolean;
  clearUnique(): void;
  getUnique(): boolean | undefined;
  setUnique(value: boolean): void;

  hasItems(): boolean;
  clearItems(): void;
  getItems(): FieldRules | undefined;
  setItems(value?: FieldRules): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RepeatedRules.AsObject;
  static toObject(includeInstance: boolean, msg: RepeatedRules): RepeatedRules.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RepeatedRules, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RepeatedRules;
  static deserializeBinaryFromReader(message: RepeatedRules, reader: jspb.BinaryReader): RepeatedRules;
}

export namespace RepeatedRules {
  export type AsObject = {
    minItems?: number,
    maxItems?: number,
    unique?: boolean,
    items?: FieldRules.AsObject,
  }
}

export class MapRules extends jspb.Message {
  hasMinPairs(): boolean;
  clearMinPairs(): void;
  getMinPairs(): number | undefined;
  setMinPairs(value: number): void;

  hasMaxPairs(): boolean;
  clearMaxPairs(): void;
  getMaxPairs(): number | undefined;
  setMaxPairs(value: number): void;

  hasNoSparse(): boolean;
  clearNoSparse(): void;
  getNoSparse(): boolean | undefined;
  setNoSparse(value: boolean): void;

  hasKeys(): boolean;
  clearKeys(): void;
  getKeys(): FieldRules | undefined;
  setKeys(value?: FieldRules): void;

  hasValues(): boolean;
  clearValues(): void;
  getValues(): FieldRules | undefined;
  setValues(value?: FieldRules): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MapRules.AsObject;
  static toObject(includeInstance: boolean, msg: MapRules): MapRules.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: MapRules, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MapRules;
  static deserializeBinaryFromReader(message: MapRules, reader: jspb.BinaryReader): MapRules;
}

export namespace MapRules {
  export type AsObject = {
    minPairs?: number,
    maxPairs?: number,
    noSparse?: boolean,
    keys?: FieldRules.AsObject,
    values?: FieldRules.AsObject,
  }
}

export class AnyRules extends jspb.Message {
  hasRequired(): boolean;
  clearRequired(): void;
  getRequired(): boolean | undefined;
  setRequired(value: boolean): void;

  clearInList(): void;
  getInList(): Array<string>;
  setInList(value: Array<string>): void;
  addIn(value: string, index?: number): string;

  clearNotInList(): void;
  getNotInList(): Array<string>;
  setNotInList(value: Array<string>): void;
  addNotIn(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AnyRules.AsObject;
  static toObject(includeInstance: boolean, msg: AnyRules): AnyRules.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: AnyRules, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AnyRules;
  static deserializeBinaryFromReader(message: AnyRules, reader: jspb.BinaryReader): AnyRules;
}

export namespace AnyRules {
  export type AsObject = {
    required?: boolean,
    inList: Array<string>,
    notInList: Array<string>,
  }
}

export class DurationRules extends jspb.Message {
  hasRequired(): boolean;
  clearRequired(): void;
  getRequired(): boolean | undefined;
  setRequired(value: boolean): void;

  hasConst(): boolean;
  clearConst(): void;
  getConst(): google_protobuf_duration_pb.Duration | undefined;
  setConst(value?: google_protobuf_duration_pb.Duration): void;

  hasLt(): boolean;
  clearLt(): void;
  getLt(): google_protobuf_duration_pb.Duration | undefined;
  setLt(value?: google_protobuf_duration_pb.Duration): void;

  hasLte(): boolean;
  clearLte(): void;
  getLte(): google_protobuf_duration_pb.Duration | undefined;
  setLte(value?: google_protobuf_duration_pb.Duration): void;

  hasGt(): boolean;
  clearGt(): void;
  getGt(): google_protobuf_duration_pb.Duration | undefined;
  setGt(value?: google_protobuf_duration_pb.Duration): void;

  hasGte(): boolean;
  clearGte(): void;
  getGte(): google_protobuf_duration_pb.Duration | undefined;
  setGte(value?: google_protobuf_duration_pb.Duration): void;

  clearInList(): void;
  getInList(): Array<google_protobuf_duration_pb.Duration>;
  setInList(value: Array<google_protobuf_duration_pb.Duration>): void;
  addIn(value?: google_protobuf_duration_pb.Duration, index?: number): google_protobuf_duration_pb.Duration;

  clearNotInList(): void;
  getNotInList(): Array<google_protobuf_duration_pb.Duration>;
  setNotInList(value: Array<google_protobuf_duration_pb.Duration>): void;
  addNotIn(value?: google_protobuf_duration_pb.Duration, index?: number): google_protobuf_duration_pb.Duration;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DurationRules.AsObject;
  static toObject(includeInstance: boolean, msg: DurationRules): DurationRules.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DurationRules, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DurationRules;
  static deserializeBinaryFromReader(message: DurationRules, reader: jspb.BinaryReader): DurationRules;
}

export namespace DurationRules {
  export type AsObject = {
    required?: boolean,
    pb_const?: google_protobuf_duration_pb.Duration.AsObject,
    lt?: google_protobuf_duration_pb.Duration.AsObject,
    lte?: google_protobuf_duration_pb.Duration.AsObject,
    gt?: google_protobuf_duration_pb.Duration.AsObject,
    gte?: google_protobuf_duration_pb.Duration.AsObject,
    inList: Array<google_protobuf_duration_pb.Duration.AsObject>,
    notInList: Array<google_protobuf_duration_pb.Duration.AsObject>,
  }
}

export class TimestampRules extends jspb.Message {
  hasRequired(): boolean;
  clearRequired(): void;
  getRequired(): boolean | undefined;
  setRequired(value: boolean): void;

  hasConst(): boolean;
  clearConst(): void;
  getConst(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setConst(value?: google_protobuf_timestamp_pb.Timestamp): void;

  hasLt(): boolean;
  clearLt(): void;
  getLt(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setLt(value?: google_protobuf_timestamp_pb.Timestamp): void;

  hasLte(): boolean;
  clearLte(): void;
  getLte(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setLte(value?: google_protobuf_timestamp_pb.Timestamp): void;

  hasGt(): boolean;
  clearGt(): void;
  getGt(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setGt(value?: google_protobuf_timestamp_pb.Timestamp): void;

  hasGte(): boolean;
  clearGte(): void;
  getGte(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setGte(value?: google_protobuf_timestamp_pb.Timestamp): void;

  hasLtNow(): boolean;
  clearLtNow(): void;
  getLtNow(): boolean | undefined;
  setLtNow(value: boolean): void;

  hasGtNow(): boolean;
  clearGtNow(): void;
  getGtNow(): boolean | undefined;
  setGtNow(value: boolean): void;

  hasWithin(): boolean;
  clearWithin(): void;
  getWithin(): google_protobuf_duration_pb.Duration | undefined;
  setWithin(value?: google_protobuf_duration_pb.Duration): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TimestampRules.AsObject;
  static toObject(includeInstance: boolean, msg: TimestampRules): TimestampRules.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: TimestampRules, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TimestampRules;
  static deserializeBinaryFromReader(message: TimestampRules, reader: jspb.BinaryReader): TimestampRules;
}

export namespace TimestampRules {
  export type AsObject = {
    required?: boolean,
    pb_const?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    lt?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    lte?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    gt?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    gte?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    ltNow?: boolean,
    gtNow?: boolean,
    within?: google_protobuf_duration_pb.Duration.AsObject,
  }
}

  export const disabled: jspb.ExtensionFieldInfo<boolean>;

  export const required: jspb.ExtensionFieldInfo<boolean>;

  export const rules: jspb.ExtensionFieldInfo<FieldRules>;

