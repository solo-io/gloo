/* eslint-disable */
// package: envoy.config.transformer.xslt.v2
// file: github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/extensions/transformers/xslt/xslt_transformer.proto

import * as jspb from "google-protobuf";

export class XsltTransformation extends jspb.Message {
  getXslt(): string;
  setXslt(value: string): void;

  getSetContentType(): string;
  setSetContentType(value: string): void;

  getNonXmlTransform(): boolean;
  setNonXmlTransform(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): XsltTransformation.AsObject;
  static toObject(includeInstance: boolean, msg: XsltTransformation): XsltTransformation.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: XsltTransformation, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): XsltTransformation;
  static deserializeBinaryFromReader(message: XsltTransformation, reader: jspb.BinaryReader): XsltTransformation;
}

export namespace XsltTransformation {
  export type AsObject = {
    xslt: string,
    setContentType: string,
    nonXmlTransform: boolean,
  }
}
