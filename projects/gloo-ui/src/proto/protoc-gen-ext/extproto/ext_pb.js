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
goog.exportSymbol('proto.extproto.equalAll', null, global);
goog.exportSymbol('proto.extproto.hashAll', null, global);
goog.exportSymbol('proto.extproto.skipHashing', null, global);

/**
 * A tuple of {field number, class constructor} for the extension
 * field named `hashAll`.
 * @type {!jspb.ExtensionFieldInfo<boolean>}
 */
proto.extproto.hashAll = new jspb.ExtensionFieldInfo(
    10071,
    {hashAll: 0},
    null,
     /** @type {?function((boolean|undefined),!jspb.Message=): !Object} */ (
         null),
    0);

google_protobuf_descriptor_pb.FileOptions.extensionsBinary[10071] = new jspb.ExtensionFieldBinaryInfo(
    proto.extproto.hashAll,
    jspb.BinaryReader.prototype.readBool,
    jspb.BinaryWriter.prototype.writeBool,
    undefined,
    undefined,
    false);
// This registers the extension field with the extended class, so that
// toObject() will function correctly.
google_protobuf_descriptor_pb.FileOptions.extensions[10071] = proto.extproto.hashAll;


/**
 * A tuple of {field number, class constructor} for the extension
 * field named `equalAll`.
 * @type {!jspb.ExtensionFieldInfo<boolean>}
 */
proto.extproto.equalAll = new jspb.ExtensionFieldInfo(
    10072,
    {equalAll: 0},
    null,
     /** @type {?function((boolean|undefined),!jspb.Message=): !Object} */ (
         null),
    0);

google_protobuf_descriptor_pb.FileOptions.extensionsBinary[10072] = new jspb.ExtensionFieldBinaryInfo(
    proto.extproto.equalAll,
    jspb.BinaryReader.prototype.readBool,
    jspb.BinaryWriter.prototype.writeBool,
    undefined,
    undefined,
    false);
// This registers the extension field with the extended class, so that
// toObject() will function correctly.
google_protobuf_descriptor_pb.FileOptions.extensions[10072] = proto.extproto.equalAll;


/**
 * A tuple of {field number, class constructor} for the extension
 * field named `skipHashing`.
 * @type {!jspb.ExtensionFieldInfo<boolean>}
 */
proto.extproto.skipHashing = new jspb.ExtensionFieldInfo(
    10071,
    {skipHashing: 0},
    null,
     /** @type {?function((boolean|undefined),!jspb.Message=): !Object} */ (
         null),
    0);

google_protobuf_descriptor_pb.FieldOptions.extensionsBinary[10071] = new jspb.ExtensionFieldBinaryInfo(
    proto.extproto.skipHashing,
    jspb.BinaryReader.prototype.readBool,
    jspb.BinaryWriter.prototype.writeBool,
    undefined,
    undefined,
    false);
// This registers the extension field with the extended class, so that
// toObject() will function correctly.
google_protobuf_descriptor_pb.FieldOptions.extensions[10071] = proto.extproto.skipHashing;

goog.object.extend(exports, proto.extproto);
