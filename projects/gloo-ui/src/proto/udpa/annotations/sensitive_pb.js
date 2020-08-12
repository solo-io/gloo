/* eslint-disable */
/**
 * @fileoverview
 * @enhanceable
 * @suppress {messageConventions} JS Compiler reports an error if a variable or
 *     field starts with 'MSG_' and isn't a translatable message.
 * @public
 */
// GENERATED CODE -- DO NOT EDIT!

var jspb = require('google-protobuf');
var goog = jspb;
var global = Function('return this')();

var google_protobuf_descriptor_pb = require('google-protobuf/google/protobuf/descriptor_pb.js');
goog.exportSymbol('proto.udpa.annotations.sensitive', null, global);

/**
 * A tuple of {field number, class constructor} for the extension
 * field named `sensitive`.
 * @type {!jspb.ExtensionFieldInfo<boolean>}
 */
proto.udpa.annotations.sensitive = new jspb.ExtensionFieldInfo(
    76569463,
    {sensitive: 0},
    null,
     /** @type {?function((boolean|undefined),!jspb.Message=): !Object} */ (
         null),
    0);

google_protobuf_descriptor_pb.FieldOptions.extensionsBinary[76569463] = new jspb.ExtensionFieldBinaryInfo(
    proto.udpa.annotations.sensitive,
    jspb.BinaryReader.prototype.readBool,
    jspb.BinaryWriter.prototype.writeBool,
    undefined,
    undefined,
    false);
// This registers the extension field with the extended class, so that
// toObject() will function correctly.
google_protobuf_descriptor_pb.FieldOptions.extensions[76569463] = proto.udpa.annotations.sensitive;

goog.object.extend(exports, proto.udpa.annotations);
