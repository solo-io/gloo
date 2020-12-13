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
var extproto_ext_pb = require('../../../../../../../../../extproto/ext_pb.js');
goog.exportSymbol('proto.solo.io.envoy.annotations.disallowedByDefault', null, global);
goog.exportSymbol('proto.solo.io.envoy.annotations.disallowedByDefaultEnum', null, global);

/**
 * A tuple of {field number, class constructor} for the extension
 * field named `disallowedByDefault`.
 * @type {!jspb.ExtensionFieldInfo<boolean>}
 */
proto.solo.io.envoy.annotations.disallowedByDefault = new jspb.ExtensionFieldInfo(
    189503208,
    {disallowedByDefault: 0},
    null,
     /** @type {?function((boolean|undefined),!jspb.Message=): !Object} */ (
         null),
    0);

google_protobuf_descriptor_pb.FieldOptions.extensionsBinary[189503208] = new jspb.ExtensionFieldBinaryInfo(
    proto.solo.io.envoy.annotations.disallowedByDefault,
    jspb.BinaryReader.prototype.readBool,
    jspb.BinaryWriter.prototype.writeBool,
    undefined,
    undefined,
    false);
// This registers the extension field with the extended class, so that
// toObject() will function correctly.
google_protobuf_descriptor_pb.FieldOptions.extensions[189503208] = proto.solo.io.envoy.annotations.disallowedByDefault;


/**
 * A tuple of {field number, class constructor} for the extension
 * field named `disallowedByDefaultEnum`.
 * @type {!jspb.ExtensionFieldInfo<boolean>}
 */
proto.solo.io.envoy.annotations.disallowedByDefaultEnum = new jspb.ExtensionFieldInfo(
    70100854,
    {disallowedByDefaultEnum: 0},
    null,
     /** @type {?function((boolean|undefined),!jspb.Message=): !Object} */ (
         null),
    0);

google_protobuf_descriptor_pb.EnumValueOptions.extensionsBinary[70100854] = new jspb.ExtensionFieldBinaryInfo(
    proto.solo.io.envoy.annotations.disallowedByDefaultEnum,
    jspb.BinaryReader.prototype.readBool,
    jspb.BinaryWriter.prototype.writeBool,
    undefined,
    undefined,
    false);
// This registers the extension field with the extended class, so that
// toObject() will function correctly.
google_protobuf_descriptor_pb.EnumValueOptions.extensions[70100854] = proto.solo.io.envoy.annotations.disallowedByDefaultEnum;

goog.object.extend(exports, proto.solo.io.envoy.annotations);
