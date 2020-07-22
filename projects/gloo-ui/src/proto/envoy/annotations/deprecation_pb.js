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
var gogoproto_gogo_pb = require('../../gogoproto/gogo_pb.js');
goog.exportSymbol('proto.envoy.annotations.disallowedByDefault', null, global);
goog.exportSymbol('proto.envoy.annotations.disallowedByDefaultEnum', null, global);

/**
 * A tuple of {field number, class constructor} for the extension
 * field named `disallowedByDefault`.
 * @type {!jspb.ExtensionFieldInfo<boolean>}
 */
proto.envoy.annotations.disallowedByDefault = new jspb.ExtensionFieldInfo(
    189503207,
    {disallowedByDefault: 0},
    null,
     /** @type {?function((boolean|undefined),!jspb.Message=): !Object} */ (
         null),
    0);

google_protobuf_descriptor_pb.FieldOptions.extensionsBinary[189503207] = new jspb.ExtensionFieldBinaryInfo(
    proto.envoy.annotations.disallowedByDefault,
    jspb.BinaryReader.prototype.readBool,
    jspb.BinaryWriter.prototype.writeBool,
    undefined,
    undefined,
    false);
// This registers the extension field with the extended class, so that
// toObject() will function correctly.
google_protobuf_descriptor_pb.FieldOptions.extensions[189503207] = proto.envoy.annotations.disallowedByDefault;


/**
 * A tuple of {field number, class constructor} for the extension
 * field named `disallowedByDefaultEnum`.
 * @type {!jspb.ExtensionFieldInfo<boolean>}
 */
proto.envoy.annotations.disallowedByDefaultEnum = new jspb.ExtensionFieldInfo(
    70100853,
    {disallowedByDefaultEnum: 0},
    null,
     /** @type {?function((boolean|undefined),!jspb.Message=): !Object} */ (
         null),
    0);

google_protobuf_descriptor_pb.EnumValueOptions.extensionsBinary[70100853] = new jspb.ExtensionFieldBinaryInfo(
    proto.envoy.annotations.disallowedByDefaultEnum,
    jspb.BinaryReader.prototype.readBool,
    jspb.BinaryWriter.prototype.writeBool,
    undefined,
    undefined,
    false);
// This registers the extension field with the extended class, so that
// toObject() will function correctly.
google_protobuf_descriptor_pb.EnumValueOptions.extensions[70100853] = proto.envoy.annotations.disallowedByDefaultEnum;

goog.object.extend(exports, proto.envoy.annotations);
