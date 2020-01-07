// package: transformation.options.gloo.solo.io
// file: gloo/projects/gloo/api/v1/options/transformation/parameters.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";
import * as gogoproto_gogo_pb from "../../../../../../../gogoproto/gogo_pb";
import * as extproto_ext_pb from "../../../../../../../extproto/ext_pb";

export class Parameters extends jspb.Message {
  getHeadersMap(): jspb.Map<string, string>;
  clearHeadersMap(): void;
  hasPath(): boolean;
  clearPath(): void;
  getPath(): google_protobuf_wrappers_pb.StringValue | undefined;
  setPath(value?: google_protobuf_wrappers_pb.StringValue): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Parameters.AsObject;
  static toObject(includeInstance: boolean, msg: Parameters): Parameters.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Parameters, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Parameters;
  static deserializeBinaryFromReader(message: Parameters, reader: jspb.BinaryReader): Parameters;
}

export namespace Parameters {
  export type AsObject = {
    headersMap: Array<[string, string]>,
    path?: google_protobuf_wrappers_pb.StringValue.AsObject,
  }
}

