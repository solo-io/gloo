// package: core.solo.io
// file: github.com/solo-io/solo-kit/api/v1/status.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../../../gogoproto/gogo_pb";

export class Status extends jspb.Message {
  getState(): Status.StateMap[keyof Status.StateMap];
  setState(value: Status.StateMap[keyof Status.StateMap]): void;

  getReason(): string;
  setReason(value: string): void;

  getReportedBy(): string;
  setReportedBy(value: string): void;

  getSubresourceStatusesMap(): jspb.Map<string, Status>;
  clearSubresourceStatusesMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Status.AsObject;
  static toObject(includeInstance: boolean, msg: Status): Status.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Status, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Status;
  static deserializeBinaryFromReader(message: Status, reader: jspb.BinaryReader): Status;
}

export namespace Status {
  export type AsObject = {
    state: Status.StateMap[keyof Status.StateMap],
    reason: string,
    reportedBy: string,
    subresourceStatusesMap: Array<[string, Status.AsObject]>,
  }

  export interface StateMap {
    PENDING: 0;
    ACCEPTED: 1;
    REJECTED: 2;
  }

  export const State: StateMap;
}

