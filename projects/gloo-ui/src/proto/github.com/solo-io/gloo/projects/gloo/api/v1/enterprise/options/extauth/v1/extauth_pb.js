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

var github_com_solo$io_solo$kit_api_v1_ref_pb = require('../../../../../../../../../../../github.com/solo-io/solo-kit/api/v1/ref_pb.js');
var extproto_ext_pb = require('../../../../../../../../../../../extproto/ext_pb.js');
var github_com_solo$io_solo$kit_api_v1_metadata_pb = require('../../../../../../../../../../../github.com/solo-io/solo-kit/api/v1/metadata_pb.js');
var github_com_solo$io_solo$kit_api_v1_status_pb = require('../../../../../../../../../../../github.com/solo-io/solo-kit/api/v1/status_pb.js');
var github_com_solo$io_solo$kit_api_v1_solo$kit_pb = require('../../../../../../../../../../../github.com/solo-io/solo-kit/api/v1/solo-kit_pb.js');
var github_com_solo$io_solo$kit_api_external_envoy_api_v2_discovery_pb = require('../../../../../../../../../../../github.com/solo-io/solo-kit/api/external/envoy/api/v2/discovery_pb.js');
var google_api_annotations_pb = require('../../../../../../../../../../../github.com/solo-io/solo-kit/api/external/google/api/annotations_pb.js');
var google_protobuf_duration_pb = require('google-protobuf/google/protobuf/duration_pb.js');
var google_protobuf_struct_pb = require('google-protobuf/google/protobuf/struct_pb.js');
var google_protobuf_wrappers_pb = require('google-protobuf/google/protobuf/wrappers_pb.js');
var google_protobuf_empty_pb = require('google-protobuf/google/protobuf/empty_pb.js');
goog.exportSymbol('proto.enterprise.gloo.solo.io.AccessTokenValidation', null, global);
goog.exportSymbol('proto.enterprise.gloo.solo.io.AccessTokenValidation.ScopeList', null, global);
goog.exportSymbol('proto.enterprise.gloo.solo.io.ApiKeyAuth', null, global);
goog.exportSymbol('proto.enterprise.gloo.solo.io.ApiKeyAuth.SecretKey', null, global);
goog.exportSymbol('proto.enterprise.gloo.solo.io.ApiKeySecret', null, global);
goog.exportSymbol('proto.enterprise.gloo.solo.io.AuthConfig', null, global);
goog.exportSymbol('proto.enterprise.gloo.solo.io.AuthConfig.Config', null, global);
goog.exportSymbol('proto.enterprise.gloo.solo.io.AuthPlugin', null, global);
goog.exportSymbol('proto.enterprise.gloo.solo.io.BasicAuth', null, global);
goog.exportSymbol('proto.enterprise.gloo.solo.io.BasicAuth.Apr', null, global);
goog.exportSymbol('proto.enterprise.gloo.solo.io.BasicAuth.Apr.SaltedHashedPassword', null, global);
goog.exportSymbol('proto.enterprise.gloo.solo.io.BufferSettings', null, global);
goog.exportSymbol('proto.enterprise.gloo.solo.io.CustomAuth', null, global);
goog.exportSymbol('proto.enterprise.gloo.solo.io.DiscoveryOverride', null, global);
goog.exportSymbol('proto.enterprise.gloo.solo.io.ExtAuthConfig', null, global);
goog.exportSymbol('proto.enterprise.gloo.solo.io.ExtAuthConfig.ApiKeyAuthConfig', null, global);
goog.exportSymbol('proto.enterprise.gloo.solo.io.ExtAuthConfig.ApiKeyAuthConfig.KeyMetadata', null, global);
goog.exportSymbol('proto.enterprise.gloo.solo.io.ExtAuthConfig.Config', null, global);
goog.exportSymbol('proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuth2Config', null, global);
goog.exportSymbol('proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuthConfig', null, global);
goog.exportSymbol('proto.enterprise.gloo.solo.io.ExtAuthConfig.OidcAuthorizationCodeConfig', null, global);
goog.exportSymbol('proto.enterprise.gloo.solo.io.ExtAuthConfig.OpaAuthConfig', null, global);
goog.exportSymbol('proto.enterprise.gloo.solo.io.ExtAuthExtension', null, global);
goog.exportSymbol('proto.enterprise.gloo.solo.io.HeaderConfiguration', null, global);
goog.exportSymbol('proto.enterprise.gloo.solo.io.HttpService', null, global);
goog.exportSymbol('proto.enterprise.gloo.solo.io.HttpService.Request', null, global);
goog.exportSymbol('proto.enterprise.gloo.solo.io.HttpService.Response', null, global);
goog.exportSymbol('proto.enterprise.gloo.solo.io.Ldap', null, global);
goog.exportSymbol('proto.enterprise.gloo.solo.io.Ldap.ConnectionPool', null, global);
goog.exportSymbol('proto.enterprise.gloo.solo.io.OAuth', null, global);
goog.exportSymbol('proto.enterprise.gloo.solo.io.OAuth2', null, global);
goog.exportSymbol('proto.enterprise.gloo.solo.io.OauthSecret', null, global);
goog.exportSymbol('proto.enterprise.gloo.solo.io.OidcAuthorizationCode', null, global);
goog.exportSymbol('proto.enterprise.gloo.solo.io.OpaAuth', null, global);
goog.exportSymbol('proto.enterprise.gloo.solo.io.PassThroughAuth', null, global);
goog.exportSymbol('proto.enterprise.gloo.solo.io.PassThroughGrpc', null, global);
goog.exportSymbol('proto.enterprise.gloo.solo.io.RedisOptions', null, global);
goog.exportSymbol('proto.enterprise.gloo.solo.io.Settings', null, global);
goog.exportSymbol('proto.enterprise.gloo.solo.io.Settings.ApiVersion', null, global);
goog.exportSymbol('proto.enterprise.gloo.solo.io.UserSession', null, global);
goog.exportSymbol('proto.enterprise.gloo.solo.io.UserSession.CookieOptions', null, global);
goog.exportSymbol('proto.enterprise.gloo.solo.io.UserSession.InternalSession', null, global);
goog.exportSymbol('proto.enterprise.gloo.solo.io.UserSession.RedisSession', null, global);

/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.enterprise.gloo.solo.io.AuthConfig = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.enterprise.gloo.solo.io.AuthConfig.repeatedFields_, null);
};
goog.inherits(proto.enterprise.gloo.solo.io.AuthConfig, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.enterprise.gloo.solo.io.AuthConfig.displayName = 'proto.enterprise.gloo.solo.io.AuthConfig';
}
/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.enterprise.gloo.solo.io.AuthConfig.repeatedFields_ = [3];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.enterprise.gloo.solo.io.AuthConfig.prototype.toObject = function(opt_includeInstance) {
  return proto.enterprise.gloo.solo.io.AuthConfig.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.enterprise.gloo.solo.io.AuthConfig} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.AuthConfig.toObject = function(includeInstance, msg) {
  var f, obj = {
    status: (f = msg.getStatus()) && github_com_solo$io_solo$kit_api_v1_status_pb.Status.toObject(includeInstance, f),
    metadata: (f = msg.getMetadata()) && github_com_solo$io_solo$kit_api_v1_metadata_pb.Metadata.toObject(includeInstance, f),
    configsList: jspb.Message.toObjectList(msg.getConfigsList(),
    proto.enterprise.gloo.solo.io.AuthConfig.Config.toObject, includeInstance),
    booleanExpr: (f = msg.getBooleanExpr()) && google_protobuf_wrappers_pb.StringValue.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.enterprise.gloo.solo.io.AuthConfig}
 */
proto.enterprise.gloo.solo.io.AuthConfig.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.enterprise.gloo.solo.io.AuthConfig;
  return proto.enterprise.gloo.solo.io.AuthConfig.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.enterprise.gloo.solo.io.AuthConfig} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.enterprise.gloo.solo.io.AuthConfig}
 */
proto.enterprise.gloo.solo.io.AuthConfig.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new github_com_solo$io_solo$kit_api_v1_status_pb.Status;
      reader.readMessage(value,github_com_solo$io_solo$kit_api_v1_status_pb.Status.deserializeBinaryFromReader);
      msg.setStatus(value);
      break;
    case 2:
      var value = new github_com_solo$io_solo$kit_api_v1_metadata_pb.Metadata;
      reader.readMessage(value,github_com_solo$io_solo$kit_api_v1_metadata_pb.Metadata.deserializeBinaryFromReader);
      msg.setMetadata(value);
      break;
    case 3:
      var value = new proto.enterprise.gloo.solo.io.AuthConfig.Config;
      reader.readMessage(value,proto.enterprise.gloo.solo.io.AuthConfig.Config.deserializeBinaryFromReader);
      msg.addConfigs(value);
      break;
    case 10:
      var value = new google_protobuf_wrappers_pb.StringValue;
      reader.readMessage(value,google_protobuf_wrappers_pb.StringValue.deserializeBinaryFromReader);
      msg.setBooleanExpr(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.enterprise.gloo.solo.io.AuthConfig.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.enterprise.gloo.solo.io.AuthConfig.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.enterprise.gloo.solo.io.AuthConfig} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.AuthConfig.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getStatus();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      github_com_solo$io_solo$kit_api_v1_status_pb.Status.serializeBinaryToWriter
    );
  }
  f = message.getMetadata();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      github_com_solo$io_solo$kit_api_v1_metadata_pb.Metadata.serializeBinaryToWriter
    );
  }
  f = message.getConfigsList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      3,
      f,
      proto.enterprise.gloo.solo.io.AuthConfig.Config.serializeBinaryToWriter
    );
  }
  f = message.getBooleanExpr();
  if (f != null) {
    writer.writeMessage(
      10,
      f,
      google_protobuf_wrappers_pb.StringValue.serializeBinaryToWriter
    );
  }
};



/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.enterprise.gloo.solo.io.AuthConfig.Config = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, proto.enterprise.gloo.solo.io.AuthConfig.Config.oneofGroups_);
};
goog.inherits(proto.enterprise.gloo.solo.io.AuthConfig.Config, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.enterprise.gloo.solo.io.AuthConfig.Config.displayName = 'proto.enterprise.gloo.solo.io.AuthConfig.Config';
}
/**
 * Oneof group definitions for this message. Each group defines the field
 * numbers belonging to that group. When of these fields' value is set, all
 * other fields in the group are cleared. During deserialization, if multiple
 * fields are encountered for a group, only the last value seen will be kept.
 * @private {!Array<!Array<number>>}
 * @const
 */
proto.enterprise.gloo.solo.io.AuthConfig.Config.oneofGroups_ = [[1,2,8,4,5,6,7,11,12]];

/**
 * @enum {number}
 */
proto.enterprise.gloo.solo.io.AuthConfig.Config.AuthConfigCase = {
  AUTH_CONFIG_NOT_SET: 0,
  BASIC_AUTH: 1,
  OAUTH: 2,
  OAUTH2: 8,
  API_KEY_AUTH: 4,
  PLUGIN_AUTH: 5,
  OPA_AUTH: 6,
  LDAP: 7,
  JWT: 11,
  PASS_THROUGH_AUTH: 12
};

/**
 * @return {proto.enterprise.gloo.solo.io.AuthConfig.Config.AuthConfigCase}
 */
proto.enterprise.gloo.solo.io.AuthConfig.Config.prototype.getAuthConfigCase = function() {
  return /** @type {proto.enterprise.gloo.solo.io.AuthConfig.Config.AuthConfigCase} */(jspb.Message.computeOneofCase(this, proto.enterprise.gloo.solo.io.AuthConfig.Config.oneofGroups_[0]));
};



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.enterprise.gloo.solo.io.AuthConfig.Config.prototype.toObject = function(opt_includeInstance) {
  return proto.enterprise.gloo.solo.io.AuthConfig.Config.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.enterprise.gloo.solo.io.AuthConfig.Config} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.AuthConfig.Config.toObject = function(includeInstance, msg) {
  var f, obj = {
    name: (f = msg.getName()) && google_protobuf_wrappers_pb.StringValue.toObject(includeInstance, f),
    basicAuth: (f = msg.getBasicAuth()) && proto.enterprise.gloo.solo.io.BasicAuth.toObject(includeInstance, f),
    oauth: (f = msg.getOauth()) && proto.enterprise.gloo.solo.io.OAuth.toObject(includeInstance, f),
    oauth2: (f = msg.getOauth2()) && proto.enterprise.gloo.solo.io.OAuth2.toObject(includeInstance, f),
    apiKeyAuth: (f = msg.getApiKeyAuth()) && proto.enterprise.gloo.solo.io.ApiKeyAuth.toObject(includeInstance, f),
    pluginAuth: (f = msg.getPluginAuth()) && proto.enterprise.gloo.solo.io.AuthPlugin.toObject(includeInstance, f),
    opaAuth: (f = msg.getOpaAuth()) && proto.enterprise.gloo.solo.io.OpaAuth.toObject(includeInstance, f),
    ldap: (f = msg.getLdap()) && proto.enterprise.gloo.solo.io.Ldap.toObject(includeInstance, f),
    jwt: (f = msg.getJwt()) && google_protobuf_empty_pb.Empty.toObject(includeInstance, f),
    passThroughAuth: (f = msg.getPassThroughAuth()) && proto.enterprise.gloo.solo.io.PassThroughAuth.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.enterprise.gloo.solo.io.AuthConfig.Config}
 */
proto.enterprise.gloo.solo.io.AuthConfig.Config.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.enterprise.gloo.solo.io.AuthConfig.Config;
  return proto.enterprise.gloo.solo.io.AuthConfig.Config.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.enterprise.gloo.solo.io.AuthConfig.Config} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.enterprise.gloo.solo.io.AuthConfig.Config}
 */
proto.enterprise.gloo.solo.io.AuthConfig.Config.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 9:
      var value = new google_protobuf_wrappers_pb.StringValue;
      reader.readMessage(value,google_protobuf_wrappers_pb.StringValue.deserializeBinaryFromReader);
      msg.setName(value);
      break;
    case 1:
      var value = new proto.enterprise.gloo.solo.io.BasicAuth;
      reader.readMessage(value,proto.enterprise.gloo.solo.io.BasicAuth.deserializeBinaryFromReader);
      msg.setBasicAuth(value);
      break;
    case 2:
      var value = new proto.enterprise.gloo.solo.io.OAuth;
      reader.readMessage(value,proto.enterprise.gloo.solo.io.OAuth.deserializeBinaryFromReader);
      msg.setOauth(value);
      break;
    case 8:
      var value = new proto.enterprise.gloo.solo.io.OAuth2;
      reader.readMessage(value,proto.enterprise.gloo.solo.io.OAuth2.deserializeBinaryFromReader);
      msg.setOauth2(value);
      break;
    case 4:
      var value = new proto.enterprise.gloo.solo.io.ApiKeyAuth;
      reader.readMessage(value,proto.enterprise.gloo.solo.io.ApiKeyAuth.deserializeBinaryFromReader);
      msg.setApiKeyAuth(value);
      break;
    case 5:
      var value = new proto.enterprise.gloo.solo.io.AuthPlugin;
      reader.readMessage(value,proto.enterprise.gloo.solo.io.AuthPlugin.deserializeBinaryFromReader);
      msg.setPluginAuth(value);
      break;
    case 6:
      var value = new proto.enterprise.gloo.solo.io.OpaAuth;
      reader.readMessage(value,proto.enterprise.gloo.solo.io.OpaAuth.deserializeBinaryFromReader);
      msg.setOpaAuth(value);
      break;
    case 7:
      var value = new proto.enterprise.gloo.solo.io.Ldap;
      reader.readMessage(value,proto.enterprise.gloo.solo.io.Ldap.deserializeBinaryFromReader);
      msg.setLdap(value);
      break;
    case 11:
      var value = new google_protobuf_empty_pb.Empty;
      reader.readMessage(value,google_protobuf_empty_pb.Empty.deserializeBinaryFromReader);
      msg.setJwt(value);
      break;
    case 12:
      var value = new proto.enterprise.gloo.solo.io.PassThroughAuth;
      reader.readMessage(value,proto.enterprise.gloo.solo.io.PassThroughAuth.deserializeBinaryFromReader);
      msg.setPassThroughAuth(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.enterprise.gloo.solo.io.AuthConfig.Config.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.enterprise.gloo.solo.io.AuthConfig.Config.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.enterprise.gloo.solo.io.AuthConfig.Config} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.AuthConfig.Config.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getName();
  if (f != null) {
    writer.writeMessage(
      9,
      f,
      google_protobuf_wrappers_pb.StringValue.serializeBinaryToWriter
    );
  }
  f = message.getBasicAuth();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto.enterprise.gloo.solo.io.BasicAuth.serializeBinaryToWriter
    );
  }
  f = message.getOauth();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.enterprise.gloo.solo.io.OAuth.serializeBinaryToWriter
    );
  }
  f = message.getOauth2();
  if (f != null) {
    writer.writeMessage(
      8,
      f,
      proto.enterprise.gloo.solo.io.OAuth2.serializeBinaryToWriter
    );
  }
  f = message.getApiKeyAuth();
  if (f != null) {
    writer.writeMessage(
      4,
      f,
      proto.enterprise.gloo.solo.io.ApiKeyAuth.serializeBinaryToWriter
    );
  }
  f = message.getPluginAuth();
  if (f != null) {
    writer.writeMessage(
      5,
      f,
      proto.enterprise.gloo.solo.io.AuthPlugin.serializeBinaryToWriter
    );
  }
  f = message.getOpaAuth();
  if (f != null) {
    writer.writeMessage(
      6,
      f,
      proto.enterprise.gloo.solo.io.OpaAuth.serializeBinaryToWriter
    );
  }
  f = message.getLdap();
  if (f != null) {
    writer.writeMessage(
      7,
      f,
      proto.enterprise.gloo.solo.io.Ldap.serializeBinaryToWriter
    );
  }
  f = message.getJwt();
  if (f != null) {
    writer.writeMessage(
      11,
      f,
      google_protobuf_empty_pb.Empty.serializeBinaryToWriter
    );
  }
  f = message.getPassThroughAuth();
  if (f != null) {
    writer.writeMessage(
      12,
      f,
      proto.enterprise.gloo.solo.io.PassThroughAuth.serializeBinaryToWriter
    );
  }
};


/**
 * optional google.protobuf.StringValue name = 9;
 * @return {?proto.google.protobuf.StringValue}
 */
proto.enterprise.gloo.solo.io.AuthConfig.Config.prototype.getName = function() {
  return /** @type{?proto.google.protobuf.StringValue} */ (
    jspb.Message.getWrapperField(this, google_protobuf_wrappers_pb.StringValue, 9));
};


/** @param {?proto.google.protobuf.StringValue|undefined} value */
proto.enterprise.gloo.solo.io.AuthConfig.Config.prototype.setName = function(value) {
  jspb.Message.setWrapperField(this, 9, value);
};


proto.enterprise.gloo.solo.io.AuthConfig.Config.prototype.clearName = function() {
  this.setName(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.enterprise.gloo.solo.io.AuthConfig.Config.prototype.hasName = function() {
  return jspb.Message.getField(this, 9) != null;
};


/**
 * optional BasicAuth basic_auth = 1;
 * @return {?proto.enterprise.gloo.solo.io.BasicAuth}
 */
proto.enterprise.gloo.solo.io.AuthConfig.Config.prototype.getBasicAuth = function() {
  return /** @type{?proto.enterprise.gloo.solo.io.BasicAuth} */ (
    jspb.Message.getWrapperField(this, proto.enterprise.gloo.solo.io.BasicAuth, 1));
};


/** @param {?proto.enterprise.gloo.solo.io.BasicAuth|undefined} value */
proto.enterprise.gloo.solo.io.AuthConfig.Config.prototype.setBasicAuth = function(value) {
  jspb.Message.setOneofWrapperField(this, 1, proto.enterprise.gloo.solo.io.AuthConfig.Config.oneofGroups_[0], value);
};


proto.enterprise.gloo.solo.io.AuthConfig.Config.prototype.clearBasicAuth = function() {
  this.setBasicAuth(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.enterprise.gloo.solo.io.AuthConfig.Config.prototype.hasBasicAuth = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional OAuth oauth = 2;
 * @return {?proto.enterprise.gloo.solo.io.OAuth}
 */
proto.enterprise.gloo.solo.io.AuthConfig.Config.prototype.getOauth = function() {
  return /** @type{?proto.enterprise.gloo.solo.io.OAuth} */ (
    jspb.Message.getWrapperField(this, proto.enterprise.gloo.solo.io.OAuth, 2));
};


/** @param {?proto.enterprise.gloo.solo.io.OAuth|undefined} value */
proto.enterprise.gloo.solo.io.AuthConfig.Config.prototype.setOauth = function(value) {
  jspb.Message.setOneofWrapperField(this, 2, proto.enterprise.gloo.solo.io.AuthConfig.Config.oneofGroups_[0], value);
};


proto.enterprise.gloo.solo.io.AuthConfig.Config.prototype.clearOauth = function() {
  this.setOauth(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.enterprise.gloo.solo.io.AuthConfig.Config.prototype.hasOauth = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional OAuth2 oauth2 = 8;
 * @return {?proto.enterprise.gloo.solo.io.OAuth2}
 */
proto.enterprise.gloo.solo.io.AuthConfig.Config.prototype.getOauth2 = function() {
  return /** @type{?proto.enterprise.gloo.solo.io.OAuth2} */ (
    jspb.Message.getWrapperField(this, proto.enterprise.gloo.solo.io.OAuth2, 8));
};


/** @param {?proto.enterprise.gloo.solo.io.OAuth2|undefined} value */
proto.enterprise.gloo.solo.io.AuthConfig.Config.prototype.setOauth2 = function(value) {
  jspb.Message.setOneofWrapperField(this, 8, proto.enterprise.gloo.solo.io.AuthConfig.Config.oneofGroups_[0], value);
};


proto.enterprise.gloo.solo.io.AuthConfig.Config.prototype.clearOauth2 = function() {
  this.setOauth2(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.enterprise.gloo.solo.io.AuthConfig.Config.prototype.hasOauth2 = function() {
  return jspb.Message.getField(this, 8) != null;
};


/**
 * optional ApiKeyAuth api_key_auth = 4;
 * @return {?proto.enterprise.gloo.solo.io.ApiKeyAuth}
 */
proto.enterprise.gloo.solo.io.AuthConfig.Config.prototype.getApiKeyAuth = function() {
  return /** @type{?proto.enterprise.gloo.solo.io.ApiKeyAuth} */ (
    jspb.Message.getWrapperField(this, proto.enterprise.gloo.solo.io.ApiKeyAuth, 4));
};


/** @param {?proto.enterprise.gloo.solo.io.ApiKeyAuth|undefined} value */
proto.enterprise.gloo.solo.io.AuthConfig.Config.prototype.setApiKeyAuth = function(value) {
  jspb.Message.setOneofWrapperField(this, 4, proto.enterprise.gloo.solo.io.AuthConfig.Config.oneofGroups_[0], value);
};


proto.enterprise.gloo.solo.io.AuthConfig.Config.prototype.clearApiKeyAuth = function() {
  this.setApiKeyAuth(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.enterprise.gloo.solo.io.AuthConfig.Config.prototype.hasApiKeyAuth = function() {
  return jspb.Message.getField(this, 4) != null;
};


/**
 * optional AuthPlugin plugin_auth = 5;
 * @return {?proto.enterprise.gloo.solo.io.AuthPlugin}
 */
proto.enterprise.gloo.solo.io.AuthConfig.Config.prototype.getPluginAuth = function() {
  return /** @type{?proto.enterprise.gloo.solo.io.AuthPlugin} */ (
    jspb.Message.getWrapperField(this, proto.enterprise.gloo.solo.io.AuthPlugin, 5));
};


/** @param {?proto.enterprise.gloo.solo.io.AuthPlugin|undefined} value */
proto.enterprise.gloo.solo.io.AuthConfig.Config.prototype.setPluginAuth = function(value) {
  jspb.Message.setOneofWrapperField(this, 5, proto.enterprise.gloo.solo.io.AuthConfig.Config.oneofGroups_[0], value);
};


proto.enterprise.gloo.solo.io.AuthConfig.Config.prototype.clearPluginAuth = function() {
  this.setPluginAuth(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.enterprise.gloo.solo.io.AuthConfig.Config.prototype.hasPluginAuth = function() {
  return jspb.Message.getField(this, 5) != null;
};


/**
 * optional OpaAuth opa_auth = 6;
 * @return {?proto.enterprise.gloo.solo.io.OpaAuth}
 */
proto.enterprise.gloo.solo.io.AuthConfig.Config.prototype.getOpaAuth = function() {
  return /** @type{?proto.enterprise.gloo.solo.io.OpaAuth} */ (
    jspb.Message.getWrapperField(this, proto.enterprise.gloo.solo.io.OpaAuth, 6));
};


/** @param {?proto.enterprise.gloo.solo.io.OpaAuth|undefined} value */
proto.enterprise.gloo.solo.io.AuthConfig.Config.prototype.setOpaAuth = function(value) {
  jspb.Message.setOneofWrapperField(this, 6, proto.enterprise.gloo.solo.io.AuthConfig.Config.oneofGroups_[0], value);
};


proto.enterprise.gloo.solo.io.AuthConfig.Config.prototype.clearOpaAuth = function() {
  this.setOpaAuth(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.enterprise.gloo.solo.io.AuthConfig.Config.prototype.hasOpaAuth = function() {
  return jspb.Message.getField(this, 6) != null;
};


/**
 * optional Ldap ldap = 7;
 * @return {?proto.enterprise.gloo.solo.io.Ldap}
 */
proto.enterprise.gloo.solo.io.AuthConfig.Config.prototype.getLdap = function() {
  return /** @type{?proto.enterprise.gloo.solo.io.Ldap} */ (
    jspb.Message.getWrapperField(this, proto.enterprise.gloo.solo.io.Ldap, 7));
};


/** @param {?proto.enterprise.gloo.solo.io.Ldap|undefined} value */
proto.enterprise.gloo.solo.io.AuthConfig.Config.prototype.setLdap = function(value) {
  jspb.Message.setOneofWrapperField(this, 7, proto.enterprise.gloo.solo.io.AuthConfig.Config.oneofGroups_[0], value);
};


proto.enterprise.gloo.solo.io.AuthConfig.Config.prototype.clearLdap = function() {
  this.setLdap(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.enterprise.gloo.solo.io.AuthConfig.Config.prototype.hasLdap = function() {
  return jspb.Message.getField(this, 7) != null;
};


/**
 * optional google.protobuf.Empty jwt = 11;
 * @return {?proto.google.protobuf.Empty}
 */
proto.enterprise.gloo.solo.io.AuthConfig.Config.prototype.getJwt = function() {
  return /** @type{?proto.google.protobuf.Empty} */ (
    jspb.Message.getWrapperField(this, google_protobuf_empty_pb.Empty, 11));
};


/** @param {?proto.google.protobuf.Empty|undefined} value */
proto.enterprise.gloo.solo.io.AuthConfig.Config.prototype.setJwt = function(value) {
  jspb.Message.setOneofWrapperField(this, 11, proto.enterprise.gloo.solo.io.AuthConfig.Config.oneofGroups_[0], value);
};


proto.enterprise.gloo.solo.io.AuthConfig.Config.prototype.clearJwt = function() {
  this.setJwt(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.enterprise.gloo.solo.io.AuthConfig.Config.prototype.hasJwt = function() {
  return jspb.Message.getField(this, 11) != null;
};


/**
 * optional PassThroughAuth pass_through_auth = 12;
 * @return {?proto.enterprise.gloo.solo.io.PassThroughAuth}
 */
proto.enterprise.gloo.solo.io.AuthConfig.Config.prototype.getPassThroughAuth = function() {
  return /** @type{?proto.enterprise.gloo.solo.io.PassThroughAuth} */ (
    jspb.Message.getWrapperField(this, proto.enterprise.gloo.solo.io.PassThroughAuth, 12));
};


/** @param {?proto.enterprise.gloo.solo.io.PassThroughAuth|undefined} value */
proto.enterprise.gloo.solo.io.AuthConfig.Config.prototype.setPassThroughAuth = function(value) {
  jspb.Message.setOneofWrapperField(this, 12, proto.enterprise.gloo.solo.io.AuthConfig.Config.oneofGroups_[0], value);
};


proto.enterprise.gloo.solo.io.AuthConfig.Config.prototype.clearPassThroughAuth = function() {
  this.setPassThroughAuth(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.enterprise.gloo.solo.io.AuthConfig.Config.prototype.hasPassThroughAuth = function() {
  return jspb.Message.getField(this, 12) != null;
};


/**
 * optional core.solo.io.Status status = 1;
 * @return {?proto.core.solo.io.Status}
 */
proto.enterprise.gloo.solo.io.AuthConfig.prototype.getStatus = function() {
  return /** @type{?proto.core.solo.io.Status} */ (
    jspb.Message.getWrapperField(this, github_com_solo$io_solo$kit_api_v1_status_pb.Status, 1));
};


/** @param {?proto.core.solo.io.Status|undefined} value */
proto.enterprise.gloo.solo.io.AuthConfig.prototype.setStatus = function(value) {
  jspb.Message.setWrapperField(this, 1, value);
};


proto.enterprise.gloo.solo.io.AuthConfig.prototype.clearStatus = function() {
  this.setStatus(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.enterprise.gloo.solo.io.AuthConfig.prototype.hasStatus = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional core.solo.io.Metadata metadata = 2;
 * @return {?proto.core.solo.io.Metadata}
 */
proto.enterprise.gloo.solo.io.AuthConfig.prototype.getMetadata = function() {
  return /** @type{?proto.core.solo.io.Metadata} */ (
    jspb.Message.getWrapperField(this, github_com_solo$io_solo$kit_api_v1_metadata_pb.Metadata, 2));
};


/** @param {?proto.core.solo.io.Metadata|undefined} value */
proto.enterprise.gloo.solo.io.AuthConfig.prototype.setMetadata = function(value) {
  jspb.Message.setWrapperField(this, 2, value);
};


proto.enterprise.gloo.solo.io.AuthConfig.prototype.clearMetadata = function() {
  this.setMetadata(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.enterprise.gloo.solo.io.AuthConfig.prototype.hasMetadata = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * repeated Config configs = 3;
 * @return {!Array<!proto.enterprise.gloo.solo.io.AuthConfig.Config>}
 */
proto.enterprise.gloo.solo.io.AuthConfig.prototype.getConfigsList = function() {
  return /** @type{!Array<!proto.enterprise.gloo.solo.io.AuthConfig.Config>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.enterprise.gloo.solo.io.AuthConfig.Config, 3));
};


/** @param {!Array<!proto.enterprise.gloo.solo.io.AuthConfig.Config>} value */
proto.enterprise.gloo.solo.io.AuthConfig.prototype.setConfigsList = function(value) {
  jspb.Message.setRepeatedWrapperField(this, 3, value);
};


/**
 * @param {!proto.enterprise.gloo.solo.io.AuthConfig.Config=} opt_value
 * @param {number=} opt_index
 * @return {!proto.enterprise.gloo.solo.io.AuthConfig.Config}
 */
proto.enterprise.gloo.solo.io.AuthConfig.prototype.addConfigs = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 3, opt_value, proto.enterprise.gloo.solo.io.AuthConfig.Config, opt_index);
};


proto.enterprise.gloo.solo.io.AuthConfig.prototype.clearConfigsList = function() {
  this.setConfigsList([]);
};


/**
 * optional google.protobuf.StringValue boolean_expr = 10;
 * @return {?proto.google.protobuf.StringValue}
 */
proto.enterprise.gloo.solo.io.AuthConfig.prototype.getBooleanExpr = function() {
  return /** @type{?proto.google.protobuf.StringValue} */ (
    jspb.Message.getWrapperField(this, google_protobuf_wrappers_pb.StringValue, 10));
};


/** @param {?proto.google.protobuf.StringValue|undefined} value */
proto.enterprise.gloo.solo.io.AuthConfig.prototype.setBooleanExpr = function(value) {
  jspb.Message.setWrapperField(this, 10, value);
};


proto.enterprise.gloo.solo.io.AuthConfig.prototype.clearBooleanExpr = function() {
  this.setBooleanExpr(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.enterprise.gloo.solo.io.AuthConfig.prototype.hasBooleanExpr = function() {
  return jspb.Message.getField(this, 10) != null;
};



/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.enterprise.gloo.solo.io.ExtAuthExtension = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, proto.enterprise.gloo.solo.io.ExtAuthExtension.oneofGroups_);
};
goog.inherits(proto.enterprise.gloo.solo.io.ExtAuthExtension, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.enterprise.gloo.solo.io.ExtAuthExtension.displayName = 'proto.enterprise.gloo.solo.io.ExtAuthExtension';
}
/**
 * Oneof group definitions for this message. Each group defines the field
 * numbers belonging to that group. When of these fields' value is set, all
 * other fields in the group are cleared. During deserialization, if multiple
 * fields are encountered for a group, only the last value seen will be kept.
 * @private {!Array<!Array<number>>}
 * @const
 */
proto.enterprise.gloo.solo.io.ExtAuthExtension.oneofGroups_ = [[1,2,3]];

/**
 * @enum {number}
 */
proto.enterprise.gloo.solo.io.ExtAuthExtension.SpecCase = {
  SPEC_NOT_SET: 0,
  DISABLE: 1,
  CONFIG_REF: 2,
  CUSTOM_AUTH: 3
};

/**
 * @return {proto.enterprise.gloo.solo.io.ExtAuthExtension.SpecCase}
 */
proto.enterprise.gloo.solo.io.ExtAuthExtension.prototype.getSpecCase = function() {
  return /** @type {proto.enterprise.gloo.solo.io.ExtAuthExtension.SpecCase} */(jspb.Message.computeOneofCase(this, proto.enterprise.gloo.solo.io.ExtAuthExtension.oneofGroups_[0]));
};



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.enterprise.gloo.solo.io.ExtAuthExtension.prototype.toObject = function(opt_includeInstance) {
  return proto.enterprise.gloo.solo.io.ExtAuthExtension.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.enterprise.gloo.solo.io.ExtAuthExtension} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.ExtAuthExtension.toObject = function(includeInstance, msg) {
  var f, obj = {
    disable: jspb.Message.getFieldWithDefault(msg, 1, false),
    configRef: (f = msg.getConfigRef()) && github_com_solo$io_solo$kit_api_v1_ref_pb.ResourceRef.toObject(includeInstance, f),
    customAuth: (f = msg.getCustomAuth()) && proto.enterprise.gloo.solo.io.CustomAuth.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.enterprise.gloo.solo.io.ExtAuthExtension}
 */
proto.enterprise.gloo.solo.io.ExtAuthExtension.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.enterprise.gloo.solo.io.ExtAuthExtension;
  return proto.enterprise.gloo.solo.io.ExtAuthExtension.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.enterprise.gloo.solo.io.ExtAuthExtension} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.enterprise.gloo.solo.io.ExtAuthExtension}
 */
proto.enterprise.gloo.solo.io.ExtAuthExtension.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setDisable(value);
      break;
    case 2:
      var value = new github_com_solo$io_solo$kit_api_v1_ref_pb.ResourceRef;
      reader.readMessage(value,github_com_solo$io_solo$kit_api_v1_ref_pb.ResourceRef.deserializeBinaryFromReader);
      msg.setConfigRef(value);
      break;
    case 3:
      var value = new proto.enterprise.gloo.solo.io.CustomAuth;
      reader.readMessage(value,proto.enterprise.gloo.solo.io.CustomAuth.deserializeBinaryFromReader);
      msg.setCustomAuth(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.enterprise.gloo.solo.io.ExtAuthExtension.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.enterprise.gloo.solo.io.ExtAuthExtension.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.enterprise.gloo.solo.io.ExtAuthExtension} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.ExtAuthExtension.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = /** @type {boolean} */ (jspb.Message.getField(message, 1));
  if (f != null) {
    writer.writeBool(
      1,
      f
    );
  }
  f = message.getConfigRef();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      github_com_solo$io_solo$kit_api_v1_ref_pb.ResourceRef.serializeBinaryToWriter
    );
  }
  f = message.getCustomAuth();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      proto.enterprise.gloo.solo.io.CustomAuth.serializeBinaryToWriter
    );
  }
};


/**
 * optional bool disable = 1;
 * Note that Boolean fields may be set to 0/1 when serialized from a Java server.
 * You should avoid comparisons like {@code val === true/false} in those cases.
 * @return {boolean}
 */
proto.enterprise.gloo.solo.io.ExtAuthExtension.prototype.getDisable = function() {
  return /** @type {boolean} */ (jspb.Message.getFieldWithDefault(this, 1, false));
};


/** @param {boolean} value */
proto.enterprise.gloo.solo.io.ExtAuthExtension.prototype.setDisable = function(value) {
  jspb.Message.setOneofField(this, 1, proto.enterprise.gloo.solo.io.ExtAuthExtension.oneofGroups_[0], value);
};


proto.enterprise.gloo.solo.io.ExtAuthExtension.prototype.clearDisable = function() {
  jspb.Message.setOneofField(this, 1, proto.enterprise.gloo.solo.io.ExtAuthExtension.oneofGroups_[0], undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.enterprise.gloo.solo.io.ExtAuthExtension.prototype.hasDisable = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional core.solo.io.ResourceRef config_ref = 2;
 * @return {?proto.core.solo.io.ResourceRef}
 */
proto.enterprise.gloo.solo.io.ExtAuthExtension.prototype.getConfigRef = function() {
  return /** @type{?proto.core.solo.io.ResourceRef} */ (
    jspb.Message.getWrapperField(this, github_com_solo$io_solo$kit_api_v1_ref_pb.ResourceRef, 2));
};


/** @param {?proto.core.solo.io.ResourceRef|undefined} value */
proto.enterprise.gloo.solo.io.ExtAuthExtension.prototype.setConfigRef = function(value) {
  jspb.Message.setOneofWrapperField(this, 2, proto.enterprise.gloo.solo.io.ExtAuthExtension.oneofGroups_[0], value);
};


proto.enterprise.gloo.solo.io.ExtAuthExtension.prototype.clearConfigRef = function() {
  this.setConfigRef(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.enterprise.gloo.solo.io.ExtAuthExtension.prototype.hasConfigRef = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional CustomAuth custom_auth = 3;
 * @return {?proto.enterprise.gloo.solo.io.CustomAuth}
 */
proto.enterprise.gloo.solo.io.ExtAuthExtension.prototype.getCustomAuth = function() {
  return /** @type{?proto.enterprise.gloo.solo.io.CustomAuth} */ (
    jspb.Message.getWrapperField(this, proto.enterprise.gloo.solo.io.CustomAuth, 3));
};


/** @param {?proto.enterprise.gloo.solo.io.CustomAuth|undefined} value */
proto.enterprise.gloo.solo.io.ExtAuthExtension.prototype.setCustomAuth = function(value) {
  jspb.Message.setOneofWrapperField(this, 3, proto.enterprise.gloo.solo.io.ExtAuthExtension.oneofGroups_[0], value);
};


proto.enterprise.gloo.solo.io.ExtAuthExtension.prototype.clearCustomAuth = function() {
  this.setCustomAuth(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.enterprise.gloo.solo.io.ExtAuthExtension.prototype.hasCustomAuth = function() {
  return jspb.Message.getField(this, 3) != null;
};



/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.enterprise.gloo.solo.io.Settings = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.enterprise.gloo.solo.io.Settings, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.enterprise.gloo.solo.io.Settings.displayName = 'proto.enterprise.gloo.solo.io.Settings';
}


if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.enterprise.gloo.solo.io.Settings.prototype.toObject = function(opt_includeInstance) {
  return proto.enterprise.gloo.solo.io.Settings.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.enterprise.gloo.solo.io.Settings} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.Settings.toObject = function(includeInstance, msg) {
  var f, obj = {
    extauthzServerRef: (f = msg.getExtauthzServerRef()) && github_com_solo$io_solo$kit_api_v1_ref_pb.ResourceRef.toObject(includeInstance, f),
    httpService: (f = msg.getHttpService()) && proto.enterprise.gloo.solo.io.HttpService.toObject(includeInstance, f),
    userIdHeader: jspb.Message.getFieldWithDefault(msg, 3, ""),
    requestTimeout: (f = msg.getRequestTimeout()) && google_protobuf_duration_pb.Duration.toObject(includeInstance, f),
    failureModeAllow: jspb.Message.getFieldWithDefault(msg, 5, false),
    requestBody: (f = msg.getRequestBody()) && proto.enterprise.gloo.solo.io.BufferSettings.toObject(includeInstance, f),
    clearRouteCache: jspb.Message.getFieldWithDefault(msg, 7, false),
    statusOnError: jspb.Message.getFieldWithDefault(msg, 8, 0),
    transportApiVersion: jspb.Message.getFieldWithDefault(msg, 9, 0)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.enterprise.gloo.solo.io.Settings}
 */
proto.enterprise.gloo.solo.io.Settings.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.enterprise.gloo.solo.io.Settings;
  return proto.enterprise.gloo.solo.io.Settings.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.enterprise.gloo.solo.io.Settings} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.enterprise.gloo.solo.io.Settings}
 */
proto.enterprise.gloo.solo.io.Settings.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new github_com_solo$io_solo$kit_api_v1_ref_pb.ResourceRef;
      reader.readMessage(value,github_com_solo$io_solo$kit_api_v1_ref_pb.ResourceRef.deserializeBinaryFromReader);
      msg.setExtauthzServerRef(value);
      break;
    case 2:
      var value = new proto.enterprise.gloo.solo.io.HttpService;
      reader.readMessage(value,proto.enterprise.gloo.solo.io.HttpService.deserializeBinaryFromReader);
      msg.setHttpService(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setUserIdHeader(value);
      break;
    case 4:
      var value = new google_protobuf_duration_pb.Duration;
      reader.readMessage(value,google_protobuf_duration_pb.Duration.deserializeBinaryFromReader);
      msg.setRequestTimeout(value);
      break;
    case 5:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setFailureModeAllow(value);
      break;
    case 6:
      var value = new proto.enterprise.gloo.solo.io.BufferSettings;
      reader.readMessage(value,proto.enterprise.gloo.solo.io.BufferSettings.deserializeBinaryFromReader);
      msg.setRequestBody(value);
      break;
    case 7:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setClearRouteCache(value);
      break;
    case 8:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setStatusOnError(value);
      break;
    case 9:
      var value = /** @type {!proto.enterprise.gloo.solo.io.Settings.ApiVersion} */ (reader.readEnum());
      msg.setTransportApiVersion(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.enterprise.gloo.solo.io.Settings.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.enterprise.gloo.solo.io.Settings.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.enterprise.gloo.solo.io.Settings} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.Settings.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getExtauthzServerRef();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      github_com_solo$io_solo$kit_api_v1_ref_pb.ResourceRef.serializeBinaryToWriter
    );
  }
  f = message.getHttpService();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.enterprise.gloo.solo.io.HttpService.serializeBinaryToWriter
    );
  }
  f = message.getUserIdHeader();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getRequestTimeout();
  if (f != null) {
    writer.writeMessage(
      4,
      f,
      google_protobuf_duration_pb.Duration.serializeBinaryToWriter
    );
  }
  f = message.getFailureModeAllow();
  if (f) {
    writer.writeBool(
      5,
      f
    );
  }
  f = message.getRequestBody();
  if (f != null) {
    writer.writeMessage(
      6,
      f,
      proto.enterprise.gloo.solo.io.BufferSettings.serializeBinaryToWriter
    );
  }
  f = message.getClearRouteCache();
  if (f) {
    writer.writeBool(
      7,
      f
    );
  }
  f = message.getStatusOnError();
  if (f !== 0) {
    writer.writeUint32(
      8,
      f
    );
  }
  f = message.getTransportApiVersion();
  if (f !== 0.0) {
    writer.writeEnum(
      9,
      f
    );
  }
};


/**
 * @enum {number}
 */
proto.enterprise.gloo.solo.io.Settings.ApiVersion = {
  V3: 0
};

/**
 * optional core.solo.io.ResourceRef extauthz_server_ref = 1;
 * @return {?proto.core.solo.io.ResourceRef}
 */
proto.enterprise.gloo.solo.io.Settings.prototype.getExtauthzServerRef = function() {
  return /** @type{?proto.core.solo.io.ResourceRef} */ (
    jspb.Message.getWrapperField(this, github_com_solo$io_solo$kit_api_v1_ref_pb.ResourceRef, 1));
};


/** @param {?proto.core.solo.io.ResourceRef|undefined} value */
proto.enterprise.gloo.solo.io.Settings.prototype.setExtauthzServerRef = function(value) {
  jspb.Message.setWrapperField(this, 1, value);
};


proto.enterprise.gloo.solo.io.Settings.prototype.clearExtauthzServerRef = function() {
  this.setExtauthzServerRef(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.enterprise.gloo.solo.io.Settings.prototype.hasExtauthzServerRef = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional HttpService http_service = 2;
 * @return {?proto.enterprise.gloo.solo.io.HttpService}
 */
proto.enterprise.gloo.solo.io.Settings.prototype.getHttpService = function() {
  return /** @type{?proto.enterprise.gloo.solo.io.HttpService} */ (
    jspb.Message.getWrapperField(this, proto.enterprise.gloo.solo.io.HttpService, 2));
};


/** @param {?proto.enterprise.gloo.solo.io.HttpService|undefined} value */
proto.enterprise.gloo.solo.io.Settings.prototype.setHttpService = function(value) {
  jspb.Message.setWrapperField(this, 2, value);
};


proto.enterprise.gloo.solo.io.Settings.prototype.clearHttpService = function() {
  this.setHttpService(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.enterprise.gloo.solo.io.Settings.prototype.hasHttpService = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional string user_id_header = 3;
 * @return {string}
 */
proto.enterprise.gloo.solo.io.Settings.prototype.getUserIdHeader = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/** @param {string} value */
proto.enterprise.gloo.solo.io.Settings.prototype.setUserIdHeader = function(value) {
  jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional google.protobuf.Duration request_timeout = 4;
 * @return {?proto.google.protobuf.Duration}
 */
proto.enterprise.gloo.solo.io.Settings.prototype.getRequestTimeout = function() {
  return /** @type{?proto.google.protobuf.Duration} */ (
    jspb.Message.getWrapperField(this, google_protobuf_duration_pb.Duration, 4));
};


/** @param {?proto.google.protobuf.Duration|undefined} value */
proto.enterprise.gloo.solo.io.Settings.prototype.setRequestTimeout = function(value) {
  jspb.Message.setWrapperField(this, 4, value);
};


proto.enterprise.gloo.solo.io.Settings.prototype.clearRequestTimeout = function() {
  this.setRequestTimeout(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.enterprise.gloo.solo.io.Settings.prototype.hasRequestTimeout = function() {
  return jspb.Message.getField(this, 4) != null;
};


/**
 * optional bool failure_mode_allow = 5;
 * Note that Boolean fields may be set to 0/1 when serialized from a Java server.
 * You should avoid comparisons like {@code val === true/false} in those cases.
 * @return {boolean}
 */
proto.enterprise.gloo.solo.io.Settings.prototype.getFailureModeAllow = function() {
  return /** @type {boolean} */ (jspb.Message.getFieldWithDefault(this, 5, false));
};


/** @param {boolean} value */
proto.enterprise.gloo.solo.io.Settings.prototype.setFailureModeAllow = function(value) {
  jspb.Message.setProto3BooleanField(this, 5, value);
};


/**
 * optional BufferSettings request_body = 6;
 * @return {?proto.enterprise.gloo.solo.io.BufferSettings}
 */
proto.enterprise.gloo.solo.io.Settings.prototype.getRequestBody = function() {
  return /** @type{?proto.enterprise.gloo.solo.io.BufferSettings} */ (
    jspb.Message.getWrapperField(this, proto.enterprise.gloo.solo.io.BufferSettings, 6));
};


/** @param {?proto.enterprise.gloo.solo.io.BufferSettings|undefined} value */
proto.enterprise.gloo.solo.io.Settings.prototype.setRequestBody = function(value) {
  jspb.Message.setWrapperField(this, 6, value);
};


proto.enterprise.gloo.solo.io.Settings.prototype.clearRequestBody = function() {
  this.setRequestBody(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.enterprise.gloo.solo.io.Settings.prototype.hasRequestBody = function() {
  return jspb.Message.getField(this, 6) != null;
};


/**
 * optional bool clear_route_cache = 7;
 * Note that Boolean fields may be set to 0/1 when serialized from a Java server.
 * You should avoid comparisons like {@code val === true/false} in those cases.
 * @return {boolean}
 */
proto.enterprise.gloo.solo.io.Settings.prototype.getClearRouteCache = function() {
  return /** @type {boolean} */ (jspb.Message.getFieldWithDefault(this, 7, false));
};


/** @param {boolean} value */
proto.enterprise.gloo.solo.io.Settings.prototype.setClearRouteCache = function(value) {
  jspb.Message.setProto3BooleanField(this, 7, value);
};


/**
 * optional uint32 status_on_error = 8;
 * @return {number}
 */
proto.enterprise.gloo.solo.io.Settings.prototype.getStatusOnError = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 8, 0));
};


/** @param {number} value */
proto.enterprise.gloo.solo.io.Settings.prototype.setStatusOnError = function(value) {
  jspb.Message.setProto3IntField(this, 8, value);
};


/**
 * optional ApiVersion transport_api_version = 9;
 * @return {!proto.enterprise.gloo.solo.io.Settings.ApiVersion}
 */
proto.enterprise.gloo.solo.io.Settings.prototype.getTransportApiVersion = function() {
  return /** @type {!proto.enterprise.gloo.solo.io.Settings.ApiVersion} */ (jspb.Message.getFieldWithDefault(this, 9, 0));
};


/** @param {!proto.enterprise.gloo.solo.io.Settings.ApiVersion} value */
proto.enterprise.gloo.solo.io.Settings.prototype.setTransportApiVersion = function(value) {
  jspb.Message.setProto3EnumField(this, 9, value);
};



/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.enterprise.gloo.solo.io.HttpService = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.enterprise.gloo.solo.io.HttpService, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.enterprise.gloo.solo.io.HttpService.displayName = 'proto.enterprise.gloo.solo.io.HttpService';
}


if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.enterprise.gloo.solo.io.HttpService.prototype.toObject = function(opt_includeInstance) {
  return proto.enterprise.gloo.solo.io.HttpService.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.enterprise.gloo.solo.io.HttpService} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.HttpService.toObject = function(includeInstance, msg) {
  var f, obj = {
    pathPrefix: jspb.Message.getFieldWithDefault(msg, 1, ""),
    request: (f = msg.getRequest()) && proto.enterprise.gloo.solo.io.HttpService.Request.toObject(includeInstance, f),
    response: (f = msg.getResponse()) && proto.enterprise.gloo.solo.io.HttpService.Response.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.enterprise.gloo.solo.io.HttpService}
 */
proto.enterprise.gloo.solo.io.HttpService.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.enterprise.gloo.solo.io.HttpService;
  return proto.enterprise.gloo.solo.io.HttpService.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.enterprise.gloo.solo.io.HttpService} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.enterprise.gloo.solo.io.HttpService}
 */
proto.enterprise.gloo.solo.io.HttpService.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setPathPrefix(value);
      break;
    case 2:
      var value = new proto.enterprise.gloo.solo.io.HttpService.Request;
      reader.readMessage(value,proto.enterprise.gloo.solo.io.HttpService.Request.deserializeBinaryFromReader);
      msg.setRequest(value);
      break;
    case 3:
      var value = new proto.enterprise.gloo.solo.io.HttpService.Response;
      reader.readMessage(value,proto.enterprise.gloo.solo.io.HttpService.Response.deserializeBinaryFromReader);
      msg.setResponse(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.enterprise.gloo.solo.io.HttpService.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.enterprise.gloo.solo.io.HttpService.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.enterprise.gloo.solo.io.HttpService} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.HttpService.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getPathPrefix();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getRequest();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.enterprise.gloo.solo.io.HttpService.Request.serializeBinaryToWriter
    );
  }
  f = message.getResponse();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      proto.enterprise.gloo.solo.io.HttpService.Response.serializeBinaryToWriter
    );
  }
};



/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.enterprise.gloo.solo.io.HttpService.Request = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.enterprise.gloo.solo.io.HttpService.Request.repeatedFields_, null);
};
goog.inherits(proto.enterprise.gloo.solo.io.HttpService.Request, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.enterprise.gloo.solo.io.HttpService.Request.displayName = 'proto.enterprise.gloo.solo.io.HttpService.Request';
}
/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.enterprise.gloo.solo.io.HttpService.Request.repeatedFields_ = [1];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.enterprise.gloo.solo.io.HttpService.Request.prototype.toObject = function(opt_includeInstance) {
  return proto.enterprise.gloo.solo.io.HttpService.Request.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.enterprise.gloo.solo.io.HttpService.Request} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.HttpService.Request.toObject = function(includeInstance, msg) {
  var f, obj = {
    allowedHeadersList: jspb.Message.getRepeatedField(msg, 1),
    headersToAddMap: (f = msg.getHeadersToAddMap()) ? f.toObject(includeInstance, undefined) : []
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.enterprise.gloo.solo.io.HttpService.Request}
 */
proto.enterprise.gloo.solo.io.HttpService.Request.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.enterprise.gloo.solo.io.HttpService.Request;
  return proto.enterprise.gloo.solo.io.HttpService.Request.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.enterprise.gloo.solo.io.HttpService.Request} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.enterprise.gloo.solo.io.HttpService.Request}
 */
proto.enterprise.gloo.solo.io.HttpService.Request.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.addAllowedHeaders(value);
      break;
    case 2:
      var value = msg.getHeadersToAddMap();
      reader.readMessage(value, function(message, reader) {
        jspb.Map.deserializeBinary(message, reader, jspb.BinaryReader.prototype.readString, jspb.BinaryReader.prototype.readString, null, "");
         });
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.enterprise.gloo.solo.io.HttpService.Request.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.enterprise.gloo.solo.io.HttpService.Request.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.enterprise.gloo.solo.io.HttpService.Request} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.HttpService.Request.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getAllowedHeadersList();
  if (f.length > 0) {
    writer.writeRepeatedString(
      1,
      f
    );
  }
  f = message.getHeadersToAddMap(true);
  if (f && f.getLength() > 0) {
    f.serializeBinary(2, writer, jspb.BinaryWriter.prototype.writeString, jspb.BinaryWriter.prototype.writeString);
  }
};


/**
 * repeated string allowed_headers = 1;
 * @return {!Array<string>}
 */
proto.enterprise.gloo.solo.io.HttpService.Request.prototype.getAllowedHeadersList = function() {
  return /** @type {!Array<string>} */ (jspb.Message.getRepeatedField(this, 1));
};


/** @param {!Array<string>} value */
proto.enterprise.gloo.solo.io.HttpService.Request.prototype.setAllowedHeadersList = function(value) {
  jspb.Message.setField(this, 1, value || []);
};


/**
 * @param {!string} value
 * @param {number=} opt_index
 */
proto.enterprise.gloo.solo.io.HttpService.Request.prototype.addAllowedHeaders = function(value, opt_index) {
  jspb.Message.addToRepeatedField(this, 1, value, opt_index);
};


proto.enterprise.gloo.solo.io.HttpService.Request.prototype.clearAllowedHeadersList = function() {
  this.setAllowedHeadersList([]);
};


/**
 * map<string, string> headers_to_add = 2;
 * @param {boolean=} opt_noLazyCreate Do not create the map if
 * empty, instead returning `undefined`
 * @return {!jspb.Map<string,string>}
 */
proto.enterprise.gloo.solo.io.HttpService.Request.prototype.getHeadersToAddMap = function(opt_noLazyCreate) {
  return /** @type {!jspb.Map<string,string>} */ (
      jspb.Message.getMapField(this, 2, opt_noLazyCreate,
      null));
};


proto.enterprise.gloo.solo.io.HttpService.Request.prototype.clearHeadersToAddMap = function() {
  this.getHeadersToAddMap().clear();
};



/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.enterprise.gloo.solo.io.HttpService.Response = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.enterprise.gloo.solo.io.HttpService.Response.repeatedFields_, null);
};
goog.inherits(proto.enterprise.gloo.solo.io.HttpService.Response, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.enterprise.gloo.solo.io.HttpService.Response.displayName = 'proto.enterprise.gloo.solo.io.HttpService.Response';
}
/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.enterprise.gloo.solo.io.HttpService.Response.repeatedFields_ = [1,2];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.enterprise.gloo.solo.io.HttpService.Response.prototype.toObject = function(opt_includeInstance) {
  return proto.enterprise.gloo.solo.io.HttpService.Response.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.enterprise.gloo.solo.io.HttpService.Response} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.HttpService.Response.toObject = function(includeInstance, msg) {
  var f, obj = {
    allowedUpstreamHeadersList: jspb.Message.getRepeatedField(msg, 1),
    allowedClientHeadersList: jspb.Message.getRepeatedField(msg, 2)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.enterprise.gloo.solo.io.HttpService.Response}
 */
proto.enterprise.gloo.solo.io.HttpService.Response.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.enterprise.gloo.solo.io.HttpService.Response;
  return proto.enterprise.gloo.solo.io.HttpService.Response.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.enterprise.gloo.solo.io.HttpService.Response} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.enterprise.gloo.solo.io.HttpService.Response}
 */
proto.enterprise.gloo.solo.io.HttpService.Response.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.addAllowedUpstreamHeaders(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.addAllowedClientHeaders(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.enterprise.gloo.solo.io.HttpService.Response.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.enterprise.gloo.solo.io.HttpService.Response.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.enterprise.gloo.solo.io.HttpService.Response} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.HttpService.Response.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getAllowedUpstreamHeadersList();
  if (f.length > 0) {
    writer.writeRepeatedString(
      1,
      f
    );
  }
  f = message.getAllowedClientHeadersList();
  if (f.length > 0) {
    writer.writeRepeatedString(
      2,
      f
    );
  }
};


/**
 * repeated string allowed_upstream_headers = 1;
 * @return {!Array<string>}
 */
proto.enterprise.gloo.solo.io.HttpService.Response.prototype.getAllowedUpstreamHeadersList = function() {
  return /** @type {!Array<string>} */ (jspb.Message.getRepeatedField(this, 1));
};


/** @param {!Array<string>} value */
proto.enterprise.gloo.solo.io.HttpService.Response.prototype.setAllowedUpstreamHeadersList = function(value) {
  jspb.Message.setField(this, 1, value || []);
};


/**
 * @param {!string} value
 * @param {number=} opt_index
 */
proto.enterprise.gloo.solo.io.HttpService.Response.prototype.addAllowedUpstreamHeaders = function(value, opt_index) {
  jspb.Message.addToRepeatedField(this, 1, value, opt_index);
};


proto.enterprise.gloo.solo.io.HttpService.Response.prototype.clearAllowedUpstreamHeadersList = function() {
  this.setAllowedUpstreamHeadersList([]);
};


/**
 * repeated string allowed_client_headers = 2;
 * @return {!Array<string>}
 */
proto.enterprise.gloo.solo.io.HttpService.Response.prototype.getAllowedClientHeadersList = function() {
  return /** @type {!Array<string>} */ (jspb.Message.getRepeatedField(this, 2));
};


/** @param {!Array<string>} value */
proto.enterprise.gloo.solo.io.HttpService.Response.prototype.setAllowedClientHeadersList = function(value) {
  jspb.Message.setField(this, 2, value || []);
};


/**
 * @param {!string} value
 * @param {number=} opt_index
 */
proto.enterprise.gloo.solo.io.HttpService.Response.prototype.addAllowedClientHeaders = function(value, opt_index) {
  jspb.Message.addToRepeatedField(this, 2, value, opt_index);
};


proto.enterprise.gloo.solo.io.HttpService.Response.prototype.clearAllowedClientHeadersList = function() {
  this.setAllowedClientHeadersList([]);
};


/**
 * optional string path_prefix = 1;
 * @return {string}
 */
proto.enterprise.gloo.solo.io.HttpService.prototype.getPathPrefix = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.enterprise.gloo.solo.io.HttpService.prototype.setPathPrefix = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional Request request = 2;
 * @return {?proto.enterprise.gloo.solo.io.HttpService.Request}
 */
proto.enterprise.gloo.solo.io.HttpService.prototype.getRequest = function() {
  return /** @type{?proto.enterprise.gloo.solo.io.HttpService.Request} */ (
    jspb.Message.getWrapperField(this, proto.enterprise.gloo.solo.io.HttpService.Request, 2));
};


/** @param {?proto.enterprise.gloo.solo.io.HttpService.Request|undefined} value */
proto.enterprise.gloo.solo.io.HttpService.prototype.setRequest = function(value) {
  jspb.Message.setWrapperField(this, 2, value);
};


proto.enterprise.gloo.solo.io.HttpService.prototype.clearRequest = function() {
  this.setRequest(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.enterprise.gloo.solo.io.HttpService.prototype.hasRequest = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional Response response = 3;
 * @return {?proto.enterprise.gloo.solo.io.HttpService.Response}
 */
proto.enterprise.gloo.solo.io.HttpService.prototype.getResponse = function() {
  return /** @type{?proto.enterprise.gloo.solo.io.HttpService.Response} */ (
    jspb.Message.getWrapperField(this, proto.enterprise.gloo.solo.io.HttpService.Response, 3));
};


/** @param {?proto.enterprise.gloo.solo.io.HttpService.Response|undefined} value */
proto.enterprise.gloo.solo.io.HttpService.prototype.setResponse = function(value) {
  jspb.Message.setWrapperField(this, 3, value);
};


proto.enterprise.gloo.solo.io.HttpService.prototype.clearResponse = function() {
  this.setResponse(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.enterprise.gloo.solo.io.HttpService.prototype.hasResponse = function() {
  return jspb.Message.getField(this, 3) != null;
};



/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.enterprise.gloo.solo.io.BufferSettings = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.enterprise.gloo.solo.io.BufferSettings, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.enterprise.gloo.solo.io.BufferSettings.displayName = 'proto.enterprise.gloo.solo.io.BufferSettings';
}


if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.enterprise.gloo.solo.io.BufferSettings.prototype.toObject = function(opt_includeInstance) {
  return proto.enterprise.gloo.solo.io.BufferSettings.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.enterprise.gloo.solo.io.BufferSettings} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.BufferSettings.toObject = function(includeInstance, msg) {
  var f, obj = {
    maxRequestBytes: jspb.Message.getFieldWithDefault(msg, 1, 0),
    allowPartialMessage: jspb.Message.getFieldWithDefault(msg, 2, false),
    packAsBytes: jspb.Message.getFieldWithDefault(msg, 3, false)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.enterprise.gloo.solo.io.BufferSettings}
 */
proto.enterprise.gloo.solo.io.BufferSettings.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.enterprise.gloo.solo.io.BufferSettings;
  return proto.enterprise.gloo.solo.io.BufferSettings.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.enterprise.gloo.solo.io.BufferSettings} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.enterprise.gloo.solo.io.BufferSettings}
 */
proto.enterprise.gloo.solo.io.BufferSettings.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setMaxRequestBytes(value);
      break;
    case 2:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setAllowPartialMessage(value);
      break;
    case 3:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setPackAsBytes(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.enterprise.gloo.solo.io.BufferSettings.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.enterprise.gloo.solo.io.BufferSettings.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.enterprise.gloo.solo.io.BufferSettings} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.BufferSettings.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getMaxRequestBytes();
  if (f !== 0) {
    writer.writeUint32(
      1,
      f
    );
  }
  f = message.getAllowPartialMessage();
  if (f) {
    writer.writeBool(
      2,
      f
    );
  }
  f = message.getPackAsBytes();
  if (f) {
    writer.writeBool(
      3,
      f
    );
  }
};


/**
 * optional uint32 max_request_bytes = 1;
 * @return {number}
 */
proto.enterprise.gloo.solo.io.BufferSettings.prototype.getMaxRequestBytes = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/** @param {number} value */
proto.enterprise.gloo.solo.io.BufferSettings.prototype.setMaxRequestBytes = function(value) {
  jspb.Message.setProto3IntField(this, 1, value);
};


/**
 * optional bool allow_partial_message = 2;
 * Note that Boolean fields may be set to 0/1 when serialized from a Java server.
 * You should avoid comparisons like {@code val === true/false} in those cases.
 * @return {boolean}
 */
proto.enterprise.gloo.solo.io.BufferSettings.prototype.getAllowPartialMessage = function() {
  return /** @type {boolean} */ (jspb.Message.getFieldWithDefault(this, 2, false));
};


/** @param {boolean} value */
proto.enterprise.gloo.solo.io.BufferSettings.prototype.setAllowPartialMessage = function(value) {
  jspb.Message.setProto3BooleanField(this, 2, value);
};


/**
 * optional bool pack_as_bytes = 3;
 * Note that Boolean fields may be set to 0/1 when serialized from a Java server.
 * You should avoid comparisons like {@code val === true/false} in those cases.
 * @return {boolean}
 */
proto.enterprise.gloo.solo.io.BufferSettings.prototype.getPackAsBytes = function() {
  return /** @type {boolean} */ (jspb.Message.getFieldWithDefault(this, 3, false));
};


/** @param {boolean} value */
proto.enterprise.gloo.solo.io.BufferSettings.prototype.setPackAsBytes = function(value) {
  jspb.Message.setProto3BooleanField(this, 3, value);
};



/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.enterprise.gloo.solo.io.CustomAuth = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.enterprise.gloo.solo.io.CustomAuth, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.enterprise.gloo.solo.io.CustomAuth.displayName = 'proto.enterprise.gloo.solo.io.CustomAuth';
}


if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.enterprise.gloo.solo.io.CustomAuth.prototype.toObject = function(opt_includeInstance) {
  return proto.enterprise.gloo.solo.io.CustomAuth.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.enterprise.gloo.solo.io.CustomAuth} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.CustomAuth.toObject = function(includeInstance, msg) {
  var f, obj = {
    contextExtensionsMap: (f = msg.getContextExtensionsMap()) ? f.toObject(includeInstance, undefined) : []
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.enterprise.gloo.solo.io.CustomAuth}
 */
proto.enterprise.gloo.solo.io.CustomAuth.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.enterprise.gloo.solo.io.CustomAuth;
  return proto.enterprise.gloo.solo.io.CustomAuth.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.enterprise.gloo.solo.io.CustomAuth} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.enterprise.gloo.solo.io.CustomAuth}
 */
proto.enterprise.gloo.solo.io.CustomAuth.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = msg.getContextExtensionsMap();
      reader.readMessage(value, function(message, reader) {
        jspb.Map.deserializeBinary(message, reader, jspb.BinaryReader.prototype.readString, jspb.BinaryReader.prototype.readString, null, "");
         });
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.enterprise.gloo.solo.io.CustomAuth.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.enterprise.gloo.solo.io.CustomAuth.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.enterprise.gloo.solo.io.CustomAuth} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.CustomAuth.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getContextExtensionsMap(true);
  if (f && f.getLength() > 0) {
    f.serializeBinary(1, writer, jspb.BinaryWriter.prototype.writeString, jspb.BinaryWriter.prototype.writeString);
  }
};


/**
 * map<string, string> context_extensions = 1;
 * @param {boolean=} opt_noLazyCreate Do not create the map if
 * empty, instead returning `undefined`
 * @return {!jspb.Map<string,string>}
 */
proto.enterprise.gloo.solo.io.CustomAuth.prototype.getContextExtensionsMap = function(opt_noLazyCreate) {
  return /** @type {!jspb.Map<string,string>} */ (
      jspb.Message.getMapField(this, 1, opt_noLazyCreate,
      null));
};


proto.enterprise.gloo.solo.io.CustomAuth.prototype.clearContextExtensionsMap = function() {
  this.getContextExtensionsMap().clear();
};



/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.enterprise.gloo.solo.io.AuthPlugin = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.enterprise.gloo.solo.io.AuthPlugin, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.enterprise.gloo.solo.io.AuthPlugin.displayName = 'proto.enterprise.gloo.solo.io.AuthPlugin';
}


if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.enterprise.gloo.solo.io.AuthPlugin.prototype.toObject = function(opt_includeInstance) {
  return proto.enterprise.gloo.solo.io.AuthPlugin.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.enterprise.gloo.solo.io.AuthPlugin} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.AuthPlugin.toObject = function(includeInstance, msg) {
  var f, obj = {
    name: jspb.Message.getFieldWithDefault(msg, 1, ""),
    pluginFileName: jspb.Message.getFieldWithDefault(msg, 2, ""),
    exportedSymbolName: jspb.Message.getFieldWithDefault(msg, 3, ""),
    config: (f = msg.getConfig()) && google_protobuf_struct_pb.Struct.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.enterprise.gloo.solo.io.AuthPlugin}
 */
proto.enterprise.gloo.solo.io.AuthPlugin.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.enterprise.gloo.solo.io.AuthPlugin;
  return proto.enterprise.gloo.solo.io.AuthPlugin.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.enterprise.gloo.solo.io.AuthPlugin} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.enterprise.gloo.solo.io.AuthPlugin}
 */
proto.enterprise.gloo.solo.io.AuthPlugin.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setName(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setPluginFileName(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setExportedSymbolName(value);
      break;
    case 4:
      var value = new google_protobuf_struct_pb.Struct;
      reader.readMessage(value,google_protobuf_struct_pb.Struct.deserializeBinaryFromReader);
      msg.setConfig(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.enterprise.gloo.solo.io.AuthPlugin.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.enterprise.gloo.solo.io.AuthPlugin.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.enterprise.gloo.solo.io.AuthPlugin} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.AuthPlugin.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getName();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getPluginFileName();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getExportedSymbolName();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getConfig();
  if (f != null) {
    writer.writeMessage(
      4,
      f,
      google_protobuf_struct_pb.Struct.serializeBinaryToWriter
    );
  }
};


/**
 * optional string name = 1;
 * @return {string}
 */
proto.enterprise.gloo.solo.io.AuthPlugin.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.enterprise.gloo.solo.io.AuthPlugin.prototype.setName = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string plugin_file_name = 2;
 * @return {string}
 */
proto.enterprise.gloo.solo.io.AuthPlugin.prototype.getPluginFileName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.enterprise.gloo.solo.io.AuthPlugin.prototype.setPluginFileName = function(value) {
  jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string exported_symbol_name = 3;
 * @return {string}
 */
proto.enterprise.gloo.solo.io.AuthPlugin.prototype.getExportedSymbolName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/** @param {string} value */
proto.enterprise.gloo.solo.io.AuthPlugin.prototype.setExportedSymbolName = function(value) {
  jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional google.protobuf.Struct config = 4;
 * @return {?proto.google.protobuf.Struct}
 */
proto.enterprise.gloo.solo.io.AuthPlugin.prototype.getConfig = function() {
  return /** @type{?proto.google.protobuf.Struct} */ (
    jspb.Message.getWrapperField(this, google_protobuf_struct_pb.Struct, 4));
};


/** @param {?proto.google.protobuf.Struct|undefined} value */
proto.enterprise.gloo.solo.io.AuthPlugin.prototype.setConfig = function(value) {
  jspb.Message.setWrapperField(this, 4, value);
};


proto.enterprise.gloo.solo.io.AuthPlugin.prototype.clearConfig = function() {
  this.setConfig(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.enterprise.gloo.solo.io.AuthPlugin.prototype.hasConfig = function() {
  return jspb.Message.getField(this, 4) != null;
};



/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.enterprise.gloo.solo.io.BasicAuth = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.enterprise.gloo.solo.io.BasicAuth, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.enterprise.gloo.solo.io.BasicAuth.displayName = 'proto.enterprise.gloo.solo.io.BasicAuth';
}


if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.enterprise.gloo.solo.io.BasicAuth.prototype.toObject = function(opt_includeInstance) {
  return proto.enterprise.gloo.solo.io.BasicAuth.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.enterprise.gloo.solo.io.BasicAuth} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.BasicAuth.toObject = function(includeInstance, msg) {
  var f, obj = {
    realm: jspb.Message.getFieldWithDefault(msg, 1, ""),
    apr: (f = msg.getApr()) && proto.enterprise.gloo.solo.io.BasicAuth.Apr.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.enterprise.gloo.solo.io.BasicAuth}
 */
proto.enterprise.gloo.solo.io.BasicAuth.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.enterprise.gloo.solo.io.BasicAuth;
  return proto.enterprise.gloo.solo.io.BasicAuth.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.enterprise.gloo.solo.io.BasicAuth} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.enterprise.gloo.solo.io.BasicAuth}
 */
proto.enterprise.gloo.solo.io.BasicAuth.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setRealm(value);
      break;
    case 2:
      var value = new proto.enterprise.gloo.solo.io.BasicAuth.Apr;
      reader.readMessage(value,proto.enterprise.gloo.solo.io.BasicAuth.Apr.deserializeBinaryFromReader);
      msg.setApr(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.enterprise.gloo.solo.io.BasicAuth.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.enterprise.gloo.solo.io.BasicAuth.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.enterprise.gloo.solo.io.BasicAuth} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.BasicAuth.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getRealm();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getApr();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.enterprise.gloo.solo.io.BasicAuth.Apr.serializeBinaryToWriter
    );
  }
};



/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.enterprise.gloo.solo.io.BasicAuth.Apr = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.enterprise.gloo.solo.io.BasicAuth.Apr, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.enterprise.gloo.solo.io.BasicAuth.Apr.displayName = 'proto.enterprise.gloo.solo.io.BasicAuth.Apr';
}


if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.enterprise.gloo.solo.io.BasicAuth.Apr.prototype.toObject = function(opt_includeInstance) {
  return proto.enterprise.gloo.solo.io.BasicAuth.Apr.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.enterprise.gloo.solo.io.BasicAuth.Apr} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.BasicAuth.Apr.toObject = function(includeInstance, msg) {
  var f, obj = {
    usersMap: (f = msg.getUsersMap()) ? f.toObject(includeInstance, proto.enterprise.gloo.solo.io.BasicAuth.Apr.SaltedHashedPassword.toObject) : []
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.enterprise.gloo.solo.io.BasicAuth.Apr}
 */
proto.enterprise.gloo.solo.io.BasicAuth.Apr.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.enterprise.gloo.solo.io.BasicAuth.Apr;
  return proto.enterprise.gloo.solo.io.BasicAuth.Apr.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.enterprise.gloo.solo.io.BasicAuth.Apr} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.enterprise.gloo.solo.io.BasicAuth.Apr}
 */
proto.enterprise.gloo.solo.io.BasicAuth.Apr.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 2:
      var value = msg.getUsersMap();
      reader.readMessage(value, function(message, reader) {
        jspb.Map.deserializeBinary(message, reader, jspb.BinaryReader.prototype.readString, jspb.BinaryReader.prototype.readMessage, proto.enterprise.gloo.solo.io.BasicAuth.Apr.SaltedHashedPassword.deserializeBinaryFromReader, "");
         });
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.enterprise.gloo.solo.io.BasicAuth.Apr.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.enterprise.gloo.solo.io.BasicAuth.Apr.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.enterprise.gloo.solo.io.BasicAuth.Apr} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.BasicAuth.Apr.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getUsersMap(true);
  if (f && f.getLength() > 0) {
    f.serializeBinary(2, writer, jspb.BinaryWriter.prototype.writeString, jspb.BinaryWriter.prototype.writeMessage, proto.enterprise.gloo.solo.io.BasicAuth.Apr.SaltedHashedPassword.serializeBinaryToWriter);
  }
};



/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.enterprise.gloo.solo.io.BasicAuth.Apr.SaltedHashedPassword = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.enterprise.gloo.solo.io.BasicAuth.Apr.SaltedHashedPassword, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.enterprise.gloo.solo.io.BasicAuth.Apr.SaltedHashedPassword.displayName = 'proto.enterprise.gloo.solo.io.BasicAuth.Apr.SaltedHashedPassword';
}


if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.enterprise.gloo.solo.io.BasicAuth.Apr.SaltedHashedPassword.prototype.toObject = function(opt_includeInstance) {
  return proto.enterprise.gloo.solo.io.BasicAuth.Apr.SaltedHashedPassword.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.enterprise.gloo.solo.io.BasicAuth.Apr.SaltedHashedPassword} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.BasicAuth.Apr.SaltedHashedPassword.toObject = function(includeInstance, msg) {
  var f, obj = {
    salt: jspb.Message.getFieldWithDefault(msg, 1, ""),
    hashedPassword: jspb.Message.getFieldWithDefault(msg, 2, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.enterprise.gloo.solo.io.BasicAuth.Apr.SaltedHashedPassword}
 */
proto.enterprise.gloo.solo.io.BasicAuth.Apr.SaltedHashedPassword.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.enterprise.gloo.solo.io.BasicAuth.Apr.SaltedHashedPassword;
  return proto.enterprise.gloo.solo.io.BasicAuth.Apr.SaltedHashedPassword.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.enterprise.gloo.solo.io.BasicAuth.Apr.SaltedHashedPassword} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.enterprise.gloo.solo.io.BasicAuth.Apr.SaltedHashedPassword}
 */
proto.enterprise.gloo.solo.io.BasicAuth.Apr.SaltedHashedPassword.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setSalt(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setHashedPassword(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.enterprise.gloo.solo.io.BasicAuth.Apr.SaltedHashedPassword.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.enterprise.gloo.solo.io.BasicAuth.Apr.SaltedHashedPassword.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.enterprise.gloo.solo.io.BasicAuth.Apr.SaltedHashedPassword} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.BasicAuth.Apr.SaltedHashedPassword.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getSalt();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getHashedPassword();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
};


/**
 * optional string salt = 1;
 * @return {string}
 */
proto.enterprise.gloo.solo.io.BasicAuth.Apr.SaltedHashedPassword.prototype.getSalt = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.enterprise.gloo.solo.io.BasicAuth.Apr.SaltedHashedPassword.prototype.setSalt = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string hashed_password = 2;
 * @return {string}
 */
proto.enterprise.gloo.solo.io.BasicAuth.Apr.SaltedHashedPassword.prototype.getHashedPassword = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.enterprise.gloo.solo.io.BasicAuth.Apr.SaltedHashedPassword.prototype.setHashedPassword = function(value) {
  jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * map<string, SaltedHashedPassword> users = 2;
 * @param {boolean=} opt_noLazyCreate Do not create the map if
 * empty, instead returning `undefined`
 * @return {!jspb.Map<string,!proto.enterprise.gloo.solo.io.BasicAuth.Apr.SaltedHashedPassword>}
 */
proto.enterprise.gloo.solo.io.BasicAuth.Apr.prototype.getUsersMap = function(opt_noLazyCreate) {
  return /** @type {!jspb.Map<string,!proto.enterprise.gloo.solo.io.BasicAuth.Apr.SaltedHashedPassword>} */ (
      jspb.Message.getMapField(this, 2, opt_noLazyCreate,
      proto.enterprise.gloo.solo.io.BasicAuth.Apr.SaltedHashedPassword));
};


proto.enterprise.gloo.solo.io.BasicAuth.Apr.prototype.clearUsersMap = function() {
  this.getUsersMap().clear();
};


/**
 * optional string realm = 1;
 * @return {string}
 */
proto.enterprise.gloo.solo.io.BasicAuth.prototype.getRealm = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.enterprise.gloo.solo.io.BasicAuth.prototype.setRealm = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional Apr apr = 2;
 * @return {?proto.enterprise.gloo.solo.io.BasicAuth.Apr}
 */
proto.enterprise.gloo.solo.io.BasicAuth.prototype.getApr = function() {
  return /** @type{?proto.enterprise.gloo.solo.io.BasicAuth.Apr} */ (
    jspb.Message.getWrapperField(this, proto.enterprise.gloo.solo.io.BasicAuth.Apr, 2));
};


/** @param {?proto.enterprise.gloo.solo.io.BasicAuth.Apr|undefined} value */
proto.enterprise.gloo.solo.io.BasicAuth.prototype.setApr = function(value) {
  jspb.Message.setWrapperField(this, 2, value);
};


proto.enterprise.gloo.solo.io.BasicAuth.prototype.clearApr = function() {
  this.setApr(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.enterprise.gloo.solo.io.BasicAuth.prototype.hasApr = function() {
  return jspb.Message.getField(this, 2) != null;
};



/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.enterprise.gloo.solo.io.OAuth = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.enterprise.gloo.solo.io.OAuth.repeatedFields_, null);
};
goog.inherits(proto.enterprise.gloo.solo.io.OAuth, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.enterprise.gloo.solo.io.OAuth.displayName = 'proto.enterprise.gloo.solo.io.OAuth';
}
/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.enterprise.gloo.solo.io.OAuth.repeatedFields_ = [6];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.enterprise.gloo.solo.io.OAuth.prototype.toObject = function(opt_includeInstance) {
  return proto.enterprise.gloo.solo.io.OAuth.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.enterprise.gloo.solo.io.OAuth} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.OAuth.toObject = function(includeInstance, msg) {
  var f, obj = {
    clientId: jspb.Message.getFieldWithDefault(msg, 1, ""),
    clientSecretRef: (f = msg.getClientSecretRef()) && github_com_solo$io_solo$kit_api_v1_ref_pb.ResourceRef.toObject(includeInstance, f),
    issuerUrl: jspb.Message.getFieldWithDefault(msg, 3, ""),
    authEndpointQueryParamsMap: (f = msg.getAuthEndpointQueryParamsMap()) ? f.toObject(includeInstance, undefined) : [],
    appUrl: jspb.Message.getFieldWithDefault(msg, 4, ""),
    callbackPath: jspb.Message.getFieldWithDefault(msg, 5, ""),
    scopesList: jspb.Message.getRepeatedField(msg, 6)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.enterprise.gloo.solo.io.OAuth}
 */
proto.enterprise.gloo.solo.io.OAuth.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.enterprise.gloo.solo.io.OAuth;
  return proto.enterprise.gloo.solo.io.OAuth.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.enterprise.gloo.solo.io.OAuth} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.enterprise.gloo.solo.io.OAuth}
 */
proto.enterprise.gloo.solo.io.OAuth.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setClientId(value);
      break;
    case 2:
      var value = new github_com_solo$io_solo$kit_api_v1_ref_pb.ResourceRef;
      reader.readMessage(value,github_com_solo$io_solo$kit_api_v1_ref_pb.ResourceRef.deserializeBinaryFromReader);
      msg.setClientSecretRef(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setIssuerUrl(value);
      break;
    case 7:
      var value = msg.getAuthEndpointQueryParamsMap();
      reader.readMessage(value, function(message, reader) {
        jspb.Map.deserializeBinary(message, reader, jspb.BinaryReader.prototype.readString, jspb.BinaryReader.prototype.readString, null, "");
         });
      break;
    case 4:
      var value = /** @type {string} */ (reader.readString());
      msg.setAppUrl(value);
      break;
    case 5:
      var value = /** @type {string} */ (reader.readString());
      msg.setCallbackPath(value);
      break;
    case 6:
      var value = /** @type {string} */ (reader.readString());
      msg.addScopes(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.enterprise.gloo.solo.io.OAuth.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.enterprise.gloo.solo.io.OAuth.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.enterprise.gloo.solo.io.OAuth} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.OAuth.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getClientId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getClientSecretRef();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      github_com_solo$io_solo$kit_api_v1_ref_pb.ResourceRef.serializeBinaryToWriter
    );
  }
  f = message.getIssuerUrl();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getAuthEndpointQueryParamsMap(true);
  if (f && f.getLength() > 0) {
    f.serializeBinary(7, writer, jspb.BinaryWriter.prototype.writeString, jspb.BinaryWriter.prototype.writeString);
  }
  f = message.getAppUrl();
  if (f.length > 0) {
    writer.writeString(
      4,
      f
    );
  }
  f = message.getCallbackPath();
  if (f.length > 0) {
    writer.writeString(
      5,
      f
    );
  }
  f = message.getScopesList();
  if (f.length > 0) {
    writer.writeRepeatedString(
      6,
      f
    );
  }
};


/**
 * optional string client_id = 1;
 * @return {string}
 */
proto.enterprise.gloo.solo.io.OAuth.prototype.getClientId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.enterprise.gloo.solo.io.OAuth.prototype.setClientId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional core.solo.io.ResourceRef client_secret_ref = 2;
 * @return {?proto.core.solo.io.ResourceRef}
 */
proto.enterprise.gloo.solo.io.OAuth.prototype.getClientSecretRef = function() {
  return /** @type{?proto.core.solo.io.ResourceRef} */ (
    jspb.Message.getWrapperField(this, github_com_solo$io_solo$kit_api_v1_ref_pb.ResourceRef, 2));
};


/** @param {?proto.core.solo.io.ResourceRef|undefined} value */
proto.enterprise.gloo.solo.io.OAuth.prototype.setClientSecretRef = function(value) {
  jspb.Message.setWrapperField(this, 2, value);
};


proto.enterprise.gloo.solo.io.OAuth.prototype.clearClientSecretRef = function() {
  this.setClientSecretRef(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.enterprise.gloo.solo.io.OAuth.prototype.hasClientSecretRef = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional string issuer_url = 3;
 * @return {string}
 */
proto.enterprise.gloo.solo.io.OAuth.prototype.getIssuerUrl = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/** @param {string} value */
proto.enterprise.gloo.solo.io.OAuth.prototype.setIssuerUrl = function(value) {
  jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * map<string, string> auth_endpoint_query_params = 7;
 * @param {boolean=} opt_noLazyCreate Do not create the map if
 * empty, instead returning `undefined`
 * @return {!jspb.Map<string,string>}
 */
proto.enterprise.gloo.solo.io.OAuth.prototype.getAuthEndpointQueryParamsMap = function(opt_noLazyCreate) {
  return /** @type {!jspb.Map<string,string>} */ (
      jspb.Message.getMapField(this, 7, opt_noLazyCreate,
      null));
};


proto.enterprise.gloo.solo.io.OAuth.prototype.clearAuthEndpointQueryParamsMap = function() {
  this.getAuthEndpointQueryParamsMap().clear();
};


/**
 * optional string app_url = 4;
 * @return {string}
 */
proto.enterprise.gloo.solo.io.OAuth.prototype.getAppUrl = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 4, ""));
};


/** @param {string} value */
proto.enterprise.gloo.solo.io.OAuth.prototype.setAppUrl = function(value) {
  jspb.Message.setProto3StringField(this, 4, value);
};


/**
 * optional string callback_path = 5;
 * @return {string}
 */
proto.enterprise.gloo.solo.io.OAuth.prototype.getCallbackPath = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 5, ""));
};


/** @param {string} value */
proto.enterprise.gloo.solo.io.OAuth.prototype.setCallbackPath = function(value) {
  jspb.Message.setProto3StringField(this, 5, value);
};


/**
 * repeated string scopes = 6;
 * @return {!Array<string>}
 */
proto.enterprise.gloo.solo.io.OAuth.prototype.getScopesList = function() {
  return /** @type {!Array<string>} */ (jspb.Message.getRepeatedField(this, 6));
};


/** @param {!Array<string>} value */
proto.enterprise.gloo.solo.io.OAuth.prototype.setScopesList = function(value) {
  jspb.Message.setField(this, 6, value || []);
};


/**
 * @param {!string} value
 * @param {number=} opt_index
 */
proto.enterprise.gloo.solo.io.OAuth.prototype.addScopes = function(value, opt_index) {
  jspb.Message.addToRepeatedField(this, 6, value, opt_index);
};


proto.enterprise.gloo.solo.io.OAuth.prototype.clearScopesList = function() {
  this.setScopesList([]);
};



/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.enterprise.gloo.solo.io.OAuth2 = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, proto.enterprise.gloo.solo.io.OAuth2.oneofGroups_);
};
goog.inherits(proto.enterprise.gloo.solo.io.OAuth2, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.enterprise.gloo.solo.io.OAuth2.displayName = 'proto.enterprise.gloo.solo.io.OAuth2';
}
/**
 * Oneof group definitions for this message. Each group defines the field
 * numbers belonging to that group. When of these fields' value is set, all
 * other fields in the group are cleared. During deserialization, if multiple
 * fields are encountered for a group, only the last value seen will be kept.
 * @private {!Array<!Array<number>>}
 * @const
 */
proto.enterprise.gloo.solo.io.OAuth2.oneofGroups_ = [[1,2]];

/**
 * @enum {number}
 */
proto.enterprise.gloo.solo.io.OAuth2.OauthTypeCase = {
  OAUTH_TYPE_NOT_SET: 0,
  OIDC_AUTHORIZATION_CODE: 1,
  ACCESS_TOKEN_VALIDATION: 2
};

/**
 * @return {proto.enterprise.gloo.solo.io.OAuth2.OauthTypeCase}
 */
proto.enterprise.gloo.solo.io.OAuth2.prototype.getOauthTypeCase = function() {
  return /** @type {proto.enterprise.gloo.solo.io.OAuth2.OauthTypeCase} */(jspb.Message.computeOneofCase(this, proto.enterprise.gloo.solo.io.OAuth2.oneofGroups_[0]));
};



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.enterprise.gloo.solo.io.OAuth2.prototype.toObject = function(opt_includeInstance) {
  return proto.enterprise.gloo.solo.io.OAuth2.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.enterprise.gloo.solo.io.OAuth2} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.OAuth2.toObject = function(includeInstance, msg) {
  var f, obj = {
    oidcAuthorizationCode: (f = msg.getOidcAuthorizationCode()) && proto.enterprise.gloo.solo.io.OidcAuthorizationCode.toObject(includeInstance, f),
    accessTokenValidation: (f = msg.getAccessTokenValidation()) && proto.enterprise.gloo.solo.io.AccessTokenValidation.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.enterprise.gloo.solo.io.OAuth2}
 */
proto.enterprise.gloo.solo.io.OAuth2.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.enterprise.gloo.solo.io.OAuth2;
  return proto.enterprise.gloo.solo.io.OAuth2.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.enterprise.gloo.solo.io.OAuth2} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.enterprise.gloo.solo.io.OAuth2}
 */
proto.enterprise.gloo.solo.io.OAuth2.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.enterprise.gloo.solo.io.OidcAuthorizationCode;
      reader.readMessage(value,proto.enterprise.gloo.solo.io.OidcAuthorizationCode.deserializeBinaryFromReader);
      msg.setOidcAuthorizationCode(value);
      break;
    case 2:
      var value = new proto.enterprise.gloo.solo.io.AccessTokenValidation;
      reader.readMessage(value,proto.enterprise.gloo.solo.io.AccessTokenValidation.deserializeBinaryFromReader);
      msg.setAccessTokenValidation(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.enterprise.gloo.solo.io.OAuth2.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.enterprise.gloo.solo.io.OAuth2.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.enterprise.gloo.solo.io.OAuth2} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.OAuth2.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getOidcAuthorizationCode();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto.enterprise.gloo.solo.io.OidcAuthorizationCode.serializeBinaryToWriter
    );
  }
  f = message.getAccessTokenValidation();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.enterprise.gloo.solo.io.AccessTokenValidation.serializeBinaryToWriter
    );
  }
};


/**
 * optional OidcAuthorizationCode oidc_authorization_code = 1;
 * @return {?proto.enterprise.gloo.solo.io.OidcAuthorizationCode}
 */
proto.enterprise.gloo.solo.io.OAuth2.prototype.getOidcAuthorizationCode = function() {
  return /** @type{?proto.enterprise.gloo.solo.io.OidcAuthorizationCode} */ (
    jspb.Message.getWrapperField(this, proto.enterprise.gloo.solo.io.OidcAuthorizationCode, 1));
};


/** @param {?proto.enterprise.gloo.solo.io.OidcAuthorizationCode|undefined} value */
proto.enterprise.gloo.solo.io.OAuth2.prototype.setOidcAuthorizationCode = function(value) {
  jspb.Message.setOneofWrapperField(this, 1, proto.enterprise.gloo.solo.io.OAuth2.oneofGroups_[0], value);
};


proto.enterprise.gloo.solo.io.OAuth2.prototype.clearOidcAuthorizationCode = function() {
  this.setOidcAuthorizationCode(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.enterprise.gloo.solo.io.OAuth2.prototype.hasOidcAuthorizationCode = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional AccessTokenValidation access_token_validation = 2;
 * @return {?proto.enterprise.gloo.solo.io.AccessTokenValidation}
 */
proto.enterprise.gloo.solo.io.OAuth2.prototype.getAccessTokenValidation = function() {
  return /** @type{?proto.enterprise.gloo.solo.io.AccessTokenValidation} */ (
    jspb.Message.getWrapperField(this, proto.enterprise.gloo.solo.io.AccessTokenValidation, 2));
};


/** @param {?proto.enterprise.gloo.solo.io.AccessTokenValidation|undefined} value */
proto.enterprise.gloo.solo.io.OAuth2.prototype.setAccessTokenValidation = function(value) {
  jspb.Message.setOneofWrapperField(this, 2, proto.enterprise.gloo.solo.io.OAuth2.oneofGroups_[0], value);
};


proto.enterprise.gloo.solo.io.OAuth2.prototype.clearAccessTokenValidation = function() {
  this.setAccessTokenValidation(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.enterprise.gloo.solo.io.OAuth2.prototype.hasAccessTokenValidation = function() {
  return jspb.Message.getField(this, 2) != null;
};



/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.enterprise.gloo.solo.io.RedisOptions = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.enterprise.gloo.solo.io.RedisOptions, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.enterprise.gloo.solo.io.RedisOptions.displayName = 'proto.enterprise.gloo.solo.io.RedisOptions';
}


if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.enterprise.gloo.solo.io.RedisOptions.prototype.toObject = function(opt_includeInstance) {
  return proto.enterprise.gloo.solo.io.RedisOptions.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.enterprise.gloo.solo.io.RedisOptions} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.RedisOptions.toObject = function(includeInstance, msg) {
  var f, obj = {
    host: jspb.Message.getFieldWithDefault(msg, 1, ""),
    db: jspb.Message.getFieldWithDefault(msg, 2, 0),
    poolSize: jspb.Message.getFieldWithDefault(msg, 3, 0)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.enterprise.gloo.solo.io.RedisOptions}
 */
proto.enterprise.gloo.solo.io.RedisOptions.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.enterprise.gloo.solo.io.RedisOptions;
  return proto.enterprise.gloo.solo.io.RedisOptions.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.enterprise.gloo.solo.io.RedisOptions} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.enterprise.gloo.solo.io.RedisOptions}
 */
proto.enterprise.gloo.solo.io.RedisOptions.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setHost(value);
      break;
    case 2:
      var value = /** @type {number} */ (reader.readInt32());
      msg.setDb(value);
      break;
    case 3:
      var value = /** @type {number} */ (reader.readInt32());
      msg.setPoolSize(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.enterprise.gloo.solo.io.RedisOptions.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.enterprise.gloo.solo.io.RedisOptions.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.enterprise.gloo.solo.io.RedisOptions} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.RedisOptions.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getHost();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getDb();
  if (f !== 0) {
    writer.writeInt32(
      2,
      f
    );
  }
  f = message.getPoolSize();
  if (f !== 0) {
    writer.writeInt32(
      3,
      f
    );
  }
};


/**
 * optional string host = 1;
 * @return {string}
 */
proto.enterprise.gloo.solo.io.RedisOptions.prototype.getHost = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.enterprise.gloo.solo.io.RedisOptions.prototype.setHost = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional int32 db = 2;
 * @return {number}
 */
proto.enterprise.gloo.solo.io.RedisOptions.prototype.getDb = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/** @param {number} value */
proto.enterprise.gloo.solo.io.RedisOptions.prototype.setDb = function(value) {
  jspb.Message.setProto3IntField(this, 2, value);
};


/**
 * optional int32 pool_size = 3;
 * @return {number}
 */
proto.enterprise.gloo.solo.io.RedisOptions.prototype.getPoolSize = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
};


/** @param {number} value */
proto.enterprise.gloo.solo.io.RedisOptions.prototype.setPoolSize = function(value) {
  jspb.Message.setProto3IntField(this, 3, value);
};



/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.enterprise.gloo.solo.io.UserSession = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, proto.enterprise.gloo.solo.io.UserSession.oneofGroups_);
};
goog.inherits(proto.enterprise.gloo.solo.io.UserSession, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.enterprise.gloo.solo.io.UserSession.displayName = 'proto.enterprise.gloo.solo.io.UserSession';
}
/**
 * Oneof group definitions for this message. Each group defines the field
 * numbers belonging to that group. When of these fields' value is set, all
 * other fields in the group are cleared. During deserialization, if multiple
 * fields are encountered for a group, only the last value seen will be kept.
 * @private {!Array<!Array<number>>}
 * @const
 */
proto.enterprise.gloo.solo.io.UserSession.oneofGroups_ = [[3,4]];

/**
 * @enum {number}
 */
proto.enterprise.gloo.solo.io.UserSession.SessionCase = {
  SESSION_NOT_SET: 0,
  COOKIE: 3,
  REDIS: 4
};

/**
 * @return {proto.enterprise.gloo.solo.io.UserSession.SessionCase}
 */
proto.enterprise.gloo.solo.io.UserSession.prototype.getSessionCase = function() {
  return /** @type {proto.enterprise.gloo.solo.io.UserSession.SessionCase} */(jspb.Message.computeOneofCase(this, proto.enterprise.gloo.solo.io.UserSession.oneofGroups_[0]));
};



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.enterprise.gloo.solo.io.UserSession.prototype.toObject = function(opt_includeInstance) {
  return proto.enterprise.gloo.solo.io.UserSession.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.enterprise.gloo.solo.io.UserSession} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.UserSession.toObject = function(includeInstance, msg) {
  var f, obj = {
    failOnFetchFailure: jspb.Message.getFieldWithDefault(msg, 1, false),
    cookieOptions: (f = msg.getCookieOptions()) && proto.enterprise.gloo.solo.io.UserSession.CookieOptions.toObject(includeInstance, f),
    cookie: (f = msg.getCookie()) && proto.enterprise.gloo.solo.io.UserSession.InternalSession.toObject(includeInstance, f),
    redis: (f = msg.getRedis()) && proto.enterprise.gloo.solo.io.UserSession.RedisSession.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.enterprise.gloo.solo.io.UserSession}
 */
proto.enterprise.gloo.solo.io.UserSession.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.enterprise.gloo.solo.io.UserSession;
  return proto.enterprise.gloo.solo.io.UserSession.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.enterprise.gloo.solo.io.UserSession} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.enterprise.gloo.solo.io.UserSession}
 */
proto.enterprise.gloo.solo.io.UserSession.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setFailOnFetchFailure(value);
      break;
    case 2:
      var value = new proto.enterprise.gloo.solo.io.UserSession.CookieOptions;
      reader.readMessage(value,proto.enterprise.gloo.solo.io.UserSession.CookieOptions.deserializeBinaryFromReader);
      msg.setCookieOptions(value);
      break;
    case 3:
      var value = new proto.enterprise.gloo.solo.io.UserSession.InternalSession;
      reader.readMessage(value,proto.enterprise.gloo.solo.io.UserSession.InternalSession.deserializeBinaryFromReader);
      msg.setCookie(value);
      break;
    case 4:
      var value = new proto.enterprise.gloo.solo.io.UserSession.RedisSession;
      reader.readMessage(value,proto.enterprise.gloo.solo.io.UserSession.RedisSession.deserializeBinaryFromReader);
      msg.setRedis(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.enterprise.gloo.solo.io.UserSession.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.enterprise.gloo.solo.io.UserSession.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.enterprise.gloo.solo.io.UserSession} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.UserSession.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getFailOnFetchFailure();
  if (f) {
    writer.writeBool(
      1,
      f
    );
  }
  f = message.getCookieOptions();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.enterprise.gloo.solo.io.UserSession.CookieOptions.serializeBinaryToWriter
    );
  }
  f = message.getCookie();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      proto.enterprise.gloo.solo.io.UserSession.InternalSession.serializeBinaryToWriter
    );
  }
  f = message.getRedis();
  if (f != null) {
    writer.writeMessage(
      4,
      f,
      proto.enterprise.gloo.solo.io.UserSession.RedisSession.serializeBinaryToWriter
    );
  }
};



/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.enterprise.gloo.solo.io.UserSession.InternalSession = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.enterprise.gloo.solo.io.UserSession.InternalSession, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.enterprise.gloo.solo.io.UserSession.InternalSession.displayName = 'proto.enterprise.gloo.solo.io.UserSession.InternalSession';
}


if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.enterprise.gloo.solo.io.UserSession.InternalSession.prototype.toObject = function(opt_includeInstance) {
  return proto.enterprise.gloo.solo.io.UserSession.InternalSession.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.enterprise.gloo.solo.io.UserSession.InternalSession} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.UserSession.InternalSession.toObject = function(includeInstance, msg) {
  var f, obj = {

  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.enterprise.gloo.solo.io.UserSession.InternalSession}
 */
proto.enterprise.gloo.solo.io.UserSession.InternalSession.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.enterprise.gloo.solo.io.UserSession.InternalSession;
  return proto.enterprise.gloo.solo.io.UserSession.InternalSession.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.enterprise.gloo.solo.io.UserSession.InternalSession} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.enterprise.gloo.solo.io.UserSession.InternalSession}
 */
proto.enterprise.gloo.solo.io.UserSession.InternalSession.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.enterprise.gloo.solo.io.UserSession.InternalSession.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.enterprise.gloo.solo.io.UserSession.InternalSession.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.enterprise.gloo.solo.io.UserSession.InternalSession} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.UserSession.InternalSession.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
};



/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.enterprise.gloo.solo.io.UserSession.RedisSession = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.enterprise.gloo.solo.io.UserSession.RedisSession, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.enterprise.gloo.solo.io.UserSession.RedisSession.displayName = 'proto.enterprise.gloo.solo.io.UserSession.RedisSession';
}


if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.enterprise.gloo.solo.io.UserSession.RedisSession.prototype.toObject = function(opt_includeInstance) {
  return proto.enterprise.gloo.solo.io.UserSession.RedisSession.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.enterprise.gloo.solo.io.UserSession.RedisSession} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.UserSession.RedisSession.toObject = function(includeInstance, msg) {
  var f, obj = {
    options: (f = msg.getOptions()) && proto.enterprise.gloo.solo.io.RedisOptions.toObject(includeInstance, f),
    keyPrefix: jspb.Message.getFieldWithDefault(msg, 2, ""),
    cookieName: jspb.Message.getFieldWithDefault(msg, 3, ""),
    allowRefreshing: (f = msg.getAllowRefreshing()) && google_protobuf_wrappers_pb.BoolValue.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.enterprise.gloo.solo.io.UserSession.RedisSession}
 */
proto.enterprise.gloo.solo.io.UserSession.RedisSession.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.enterprise.gloo.solo.io.UserSession.RedisSession;
  return proto.enterprise.gloo.solo.io.UserSession.RedisSession.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.enterprise.gloo.solo.io.UserSession.RedisSession} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.enterprise.gloo.solo.io.UserSession.RedisSession}
 */
proto.enterprise.gloo.solo.io.UserSession.RedisSession.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.enterprise.gloo.solo.io.RedisOptions;
      reader.readMessage(value,proto.enterprise.gloo.solo.io.RedisOptions.deserializeBinaryFromReader);
      msg.setOptions(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setKeyPrefix(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setCookieName(value);
      break;
    case 4:
      var value = new google_protobuf_wrappers_pb.BoolValue;
      reader.readMessage(value,google_protobuf_wrappers_pb.BoolValue.deserializeBinaryFromReader);
      msg.setAllowRefreshing(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.enterprise.gloo.solo.io.UserSession.RedisSession.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.enterprise.gloo.solo.io.UserSession.RedisSession.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.enterprise.gloo.solo.io.UserSession.RedisSession} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.UserSession.RedisSession.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getOptions();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto.enterprise.gloo.solo.io.RedisOptions.serializeBinaryToWriter
    );
  }
  f = message.getKeyPrefix();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getCookieName();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getAllowRefreshing();
  if (f != null) {
    writer.writeMessage(
      4,
      f,
      google_protobuf_wrappers_pb.BoolValue.serializeBinaryToWriter
    );
  }
};


/**
 * optional RedisOptions options = 1;
 * @return {?proto.enterprise.gloo.solo.io.RedisOptions}
 */
proto.enterprise.gloo.solo.io.UserSession.RedisSession.prototype.getOptions = function() {
  return /** @type{?proto.enterprise.gloo.solo.io.RedisOptions} */ (
    jspb.Message.getWrapperField(this, proto.enterprise.gloo.solo.io.RedisOptions, 1));
};


/** @param {?proto.enterprise.gloo.solo.io.RedisOptions|undefined} value */
proto.enterprise.gloo.solo.io.UserSession.RedisSession.prototype.setOptions = function(value) {
  jspb.Message.setWrapperField(this, 1, value);
};


proto.enterprise.gloo.solo.io.UserSession.RedisSession.prototype.clearOptions = function() {
  this.setOptions(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.enterprise.gloo.solo.io.UserSession.RedisSession.prototype.hasOptions = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional string key_prefix = 2;
 * @return {string}
 */
proto.enterprise.gloo.solo.io.UserSession.RedisSession.prototype.getKeyPrefix = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.enterprise.gloo.solo.io.UserSession.RedisSession.prototype.setKeyPrefix = function(value) {
  jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string cookie_name = 3;
 * @return {string}
 */
proto.enterprise.gloo.solo.io.UserSession.RedisSession.prototype.getCookieName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/** @param {string} value */
proto.enterprise.gloo.solo.io.UserSession.RedisSession.prototype.setCookieName = function(value) {
  jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional google.protobuf.BoolValue allow_refreshing = 4;
 * @return {?proto.google.protobuf.BoolValue}
 */
proto.enterprise.gloo.solo.io.UserSession.RedisSession.prototype.getAllowRefreshing = function() {
  return /** @type{?proto.google.protobuf.BoolValue} */ (
    jspb.Message.getWrapperField(this, google_protobuf_wrappers_pb.BoolValue, 4));
};


/** @param {?proto.google.protobuf.BoolValue|undefined} value */
proto.enterprise.gloo.solo.io.UserSession.RedisSession.prototype.setAllowRefreshing = function(value) {
  jspb.Message.setWrapperField(this, 4, value);
};


proto.enterprise.gloo.solo.io.UserSession.RedisSession.prototype.clearAllowRefreshing = function() {
  this.setAllowRefreshing(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.enterprise.gloo.solo.io.UserSession.RedisSession.prototype.hasAllowRefreshing = function() {
  return jspb.Message.getField(this, 4) != null;
};



/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.enterprise.gloo.solo.io.UserSession.CookieOptions = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.enterprise.gloo.solo.io.UserSession.CookieOptions, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.enterprise.gloo.solo.io.UserSession.CookieOptions.displayName = 'proto.enterprise.gloo.solo.io.UserSession.CookieOptions';
}


if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.enterprise.gloo.solo.io.UserSession.CookieOptions.prototype.toObject = function(opt_includeInstance) {
  return proto.enterprise.gloo.solo.io.UserSession.CookieOptions.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.enterprise.gloo.solo.io.UserSession.CookieOptions} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.UserSession.CookieOptions.toObject = function(includeInstance, msg) {
  var f, obj = {
    maxAge: (f = msg.getMaxAge()) && google_protobuf_wrappers_pb.UInt32Value.toObject(includeInstance, f),
    notSecure: jspb.Message.getFieldWithDefault(msg, 2, false),
    path: (f = msg.getPath()) && google_protobuf_wrappers_pb.StringValue.toObject(includeInstance, f),
    domain: jspb.Message.getFieldWithDefault(msg, 4, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.enterprise.gloo.solo.io.UserSession.CookieOptions}
 */
proto.enterprise.gloo.solo.io.UserSession.CookieOptions.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.enterprise.gloo.solo.io.UserSession.CookieOptions;
  return proto.enterprise.gloo.solo.io.UserSession.CookieOptions.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.enterprise.gloo.solo.io.UserSession.CookieOptions} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.enterprise.gloo.solo.io.UserSession.CookieOptions}
 */
proto.enterprise.gloo.solo.io.UserSession.CookieOptions.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new google_protobuf_wrappers_pb.UInt32Value;
      reader.readMessage(value,google_protobuf_wrappers_pb.UInt32Value.deserializeBinaryFromReader);
      msg.setMaxAge(value);
      break;
    case 2:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setNotSecure(value);
      break;
    case 3:
      var value = new google_protobuf_wrappers_pb.StringValue;
      reader.readMessage(value,google_protobuf_wrappers_pb.StringValue.deserializeBinaryFromReader);
      msg.setPath(value);
      break;
    case 4:
      var value = /** @type {string} */ (reader.readString());
      msg.setDomain(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.enterprise.gloo.solo.io.UserSession.CookieOptions.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.enterprise.gloo.solo.io.UserSession.CookieOptions.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.enterprise.gloo.solo.io.UserSession.CookieOptions} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.UserSession.CookieOptions.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getMaxAge();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      google_protobuf_wrappers_pb.UInt32Value.serializeBinaryToWriter
    );
  }
  f = message.getNotSecure();
  if (f) {
    writer.writeBool(
      2,
      f
    );
  }
  f = message.getPath();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      google_protobuf_wrappers_pb.StringValue.serializeBinaryToWriter
    );
  }
  f = message.getDomain();
  if (f.length > 0) {
    writer.writeString(
      4,
      f
    );
  }
};


/**
 * optional google.protobuf.UInt32Value max_age = 1;
 * @return {?proto.google.protobuf.UInt32Value}
 */
proto.enterprise.gloo.solo.io.UserSession.CookieOptions.prototype.getMaxAge = function() {
  return /** @type{?proto.google.protobuf.UInt32Value} */ (
    jspb.Message.getWrapperField(this, google_protobuf_wrappers_pb.UInt32Value, 1));
};


/** @param {?proto.google.protobuf.UInt32Value|undefined} value */
proto.enterprise.gloo.solo.io.UserSession.CookieOptions.prototype.setMaxAge = function(value) {
  jspb.Message.setWrapperField(this, 1, value);
};


proto.enterprise.gloo.solo.io.UserSession.CookieOptions.prototype.clearMaxAge = function() {
  this.setMaxAge(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.enterprise.gloo.solo.io.UserSession.CookieOptions.prototype.hasMaxAge = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional bool not_secure = 2;
 * Note that Boolean fields may be set to 0/1 when serialized from a Java server.
 * You should avoid comparisons like {@code val === true/false} in those cases.
 * @return {boolean}
 */
proto.enterprise.gloo.solo.io.UserSession.CookieOptions.prototype.getNotSecure = function() {
  return /** @type {boolean} */ (jspb.Message.getFieldWithDefault(this, 2, false));
};


/** @param {boolean} value */
proto.enterprise.gloo.solo.io.UserSession.CookieOptions.prototype.setNotSecure = function(value) {
  jspb.Message.setProto3BooleanField(this, 2, value);
};


/**
 * optional google.protobuf.StringValue path = 3;
 * @return {?proto.google.protobuf.StringValue}
 */
proto.enterprise.gloo.solo.io.UserSession.CookieOptions.prototype.getPath = function() {
  return /** @type{?proto.google.protobuf.StringValue} */ (
    jspb.Message.getWrapperField(this, google_protobuf_wrappers_pb.StringValue, 3));
};


/** @param {?proto.google.protobuf.StringValue|undefined} value */
proto.enterprise.gloo.solo.io.UserSession.CookieOptions.prototype.setPath = function(value) {
  jspb.Message.setWrapperField(this, 3, value);
};


proto.enterprise.gloo.solo.io.UserSession.CookieOptions.prototype.clearPath = function() {
  this.setPath(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.enterprise.gloo.solo.io.UserSession.CookieOptions.prototype.hasPath = function() {
  return jspb.Message.getField(this, 3) != null;
};


/**
 * optional string domain = 4;
 * @return {string}
 */
proto.enterprise.gloo.solo.io.UserSession.CookieOptions.prototype.getDomain = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 4, ""));
};


/** @param {string} value */
proto.enterprise.gloo.solo.io.UserSession.CookieOptions.prototype.setDomain = function(value) {
  jspb.Message.setProto3StringField(this, 4, value);
};


/**
 * optional bool fail_on_fetch_failure = 1;
 * Note that Boolean fields may be set to 0/1 when serialized from a Java server.
 * You should avoid comparisons like {@code val === true/false} in those cases.
 * @return {boolean}
 */
proto.enterprise.gloo.solo.io.UserSession.prototype.getFailOnFetchFailure = function() {
  return /** @type {boolean} */ (jspb.Message.getFieldWithDefault(this, 1, false));
};


/** @param {boolean} value */
proto.enterprise.gloo.solo.io.UserSession.prototype.setFailOnFetchFailure = function(value) {
  jspb.Message.setProto3BooleanField(this, 1, value);
};


/**
 * optional CookieOptions cookie_options = 2;
 * @return {?proto.enterprise.gloo.solo.io.UserSession.CookieOptions}
 */
proto.enterprise.gloo.solo.io.UserSession.prototype.getCookieOptions = function() {
  return /** @type{?proto.enterprise.gloo.solo.io.UserSession.CookieOptions} */ (
    jspb.Message.getWrapperField(this, proto.enterprise.gloo.solo.io.UserSession.CookieOptions, 2));
};


/** @param {?proto.enterprise.gloo.solo.io.UserSession.CookieOptions|undefined} value */
proto.enterprise.gloo.solo.io.UserSession.prototype.setCookieOptions = function(value) {
  jspb.Message.setWrapperField(this, 2, value);
};


proto.enterprise.gloo.solo.io.UserSession.prototype.clearCookieOptions = function() {
  this.setCookieOptions(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.enterprise.gloo.solo.io.UserSession.prototype.hasCookieOptions = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional InternalSession cookie = 3;
 * @return {?proto.enterprise.gloo.solo.io.UserSession.InternalSession}
 */
proto.enterprise.gloo.solo.io.UserSession.prototype.getCookie = function() {
  return /** @type{?proto.enterprise.gloo.solo.io.UserSession.InternalSession} */ (
    jspb.Message.getWrapperField(this, proto.enterprise.gloo.solo.io.UserSession.InternalSession, 3));
};


/** @param {?proto.enterprise.gloo.solo.io.UserSession.InternalSession|undefined} value */
proto.enterprise.gloo.solo.io.UserSession.prototype.setCookie = function(value) {
  jspb.Message.setOneofWrapperField(this, 3, proto.enterprise.gloo.solo.io.UserSession.oneofGroups_[0], value);
};


proto.enterprise.gloo.solo.io.UserSession.prototype.clearCookie = function() {
  this.setCookie(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.enterprise.gloo.solo.io.UserSession.prototype.hasCookie = function() {
  return jspb.Message.getField(this, 3) != null;
};


/**
 * optional RedisSession redis = 4;
 * @return {?proto.enterprise.gloo.solo.io.UserSession.RedisSession}
 */
proto.enterprise.gloo.solo.io.UserSession.prototype.getRedis = function() {
  return /** @type{?proto.enterprise.gloo.solo.io.UserSession.RedisSession} */ (
    jspb.Message.getWrapperField(this, proto.enterprise.gloo.solo.io.UserSession.RedisSession, 4));
};


/** @param {?proto.enterprise.gloo.solo.io.UserSession.RedisSession|undefined} value */
proto.enterprise.gloo.solo.io.UserSession.prototype.setRedis = function(value) {
  jspb.Message.setOneofWrapperField(this, 4, proto.enterprise.gloo.solo.io.UserSession.oneofGroups_[0], value);
};


proto.enterprise.gloo.solo.io.UserSession.prototype.clearRedis = function() {
  this.setRedis(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.enterprise.gloo.solo.io.UserSession.prototype.hasRedis = function() {
  return jspb.Message.getField(this, 4) != null;
};



/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.enterprise.gloo.solo.io.HeaderConfiguration = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.enterprise.gloo.solo.io.HeaderConfiguration, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.enterprise.gloo.solo.io.HeaderConfiguration.displayName = 'proto.enterprise.gloo.solo.io.HeaderConfiguration';
}


if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.enterprise.gloo.solo.io.HeaderConfiguration.prototype.toObject = function(opt_includeInstance) {
  return proto.enterprise.gloo.solo.io.HeaderConfiguration.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.enterprise.gloo.solo.io.HeaderConfiguration} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.HeaderConfiguration.toObject = function(includeInstance, msg) {
  var f, obj = {
    idTokenHeader: jspb.Message.getFieldWithDefault(msg, 1, ""),
    accessTokenHeader: jspb.Message.getFieldWithDefault(msg, 2, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.enterprise.gloo.solo.io.HeaderConfiguration}
 */
proto.enterprise.gloo.solo.io.HeaderConfiguration.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.enterprise.gloo.solo.io.HeaderConfiguration;
  return proto.enterprise.gloo.solo.io.HeaderConfiguration.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.enterprise.gloo.solo.io.HeaderConfiguration} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.enterprise.gloo.solo.io.HeaderConfiguration}
 */
proto.enterprise.gloo.solo.io.HeaderConfiguration.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setIdTokenHeader(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setAccessTokenHeader(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.enterprise.gloo.solo.io.HeaderConfiguration.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.enterprise.gloo.solo.io.HeaderConfiguration.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.enterprise.gloo.solo.io.HeaderConfiguration} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.HeaderConfiguration.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getIdTokenHeader();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getAccessTokenHeader();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
};


/**
 * optional string id_token_header = 1;
 * @return {string}
 */
proto.enterprise.gloo.solo.io.HeaderConfiguration.prototype.getIdTokenHeader = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.enterprise.gloo.solo.io.HeaderConfiguration.prototype.setIdTokenHeader = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string access_token_header = 2;
 * @return {string}
 */
proto.enterprise.gloo.solo.io.HeaderConfiguration.prototype.getAccessTokenHeader = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.enterprise.gloo.solo.io.HeaderConfiguration.prototype.setAccessTokenHeader = function(value) {
  jspb.Message.setProto3StringField(this, 2, value);
};



/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.enterprise.gloo.solo.io.DiscoveryOverride = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.enterprise.gloo.solo.io.DiscoveryOverride.repeatedFields_, null);
};
goog.inherits(proto.enterprise.gloo.solo.io.DiscoveryOverride, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.enterprise.gloo.solo.io.DiscoveryOverride.displayName = 'proto.enterprise.gloo.solo.io.DiscoveryOverride';
}
/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.enterprise.gloo.solo.io.DiscoveryOverride.repeatedFields_ = [4,5,6,7,8,9];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.enterprise.gloo.solo.io.DiscoveryOverride.prototype.toObject = function(opt_includeInstance) {
  return proto.enterprise.gloo.solo.io.DiscoveryOverride.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.enterprise.gloo.solo.io.DiscoveryOverride} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.DiscoveryOverride.toObject = function(includeInstance, msg) {
  var f, obj = {
    authEndpoint: jspb.Message.getFieldWithDefault(msg, 1, ""),
    tokenEndpoint: jspb.Message.getFieldWithDefault(msg, 2, ""),
    jwksUri: jspb.Message.getFieldWithDefault(msg, 3, ""),
    scopesList: jspb.Message.getRepeatedField(msg, 4),
    responseTypesList: jspb.Message.getRepeatedField(msg, 5),
    subjectsList: jspb.Message.getRepeatedField(msg, 6),
    idTokenAlgsList: jspb.Message.getRepeatedField(msg, 7),
    authMethodsList: jspb.Message.getRepeatedField(msg, 8),
    claimsList: jspb.Message.getRepeatedField(msg, 9)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.enterprise.gloo.solo.io.DiscoveryOverride}
 */
proto.enterprise.gloo.solo.io.DiscoveryOverride.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.enterprise.gloo.solo.io.DiscoveryOverride;
  return proto.enterprise.gloo.solo.io.DiscoveryOverride.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.enterprise.gloo.solo.io.DiscoveryOverride} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.enterprise.gloo.solo.io.DiscoveryOverride}
 */
proto.enterprise.gloo.solo.io.DiscoveryOverride.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setAuthEndpoint(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setTokenEndpoint(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setJwksUri(value);
      break;
    case 4:
      var value = /** @type {string} */ (reader.readString());
      msg.addScopes(value);
      break;
    case 5:
      var value = /** @type {string} */ (reader.readString());
      msg.addResponseTypes(value);
      break;
    case 6:
      var value = /** @type {string} */ (reader.readString());
      msg.addSubjects(value);
      break;
    case 7:
      var value = /** @type {string} */ (reader.readString());
      msg.addIdTokenAlgs(value);
      break;
    case 8:
      var value = /** @type {string} */ (reader.readString());
      msg.addAuthMethods(value);
      break;
    case 9:
      var value = /** @type {string} */ (reader.readString());
      msg.addClaims(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.enterprise.gloo.solo.io.DiscoveryOverride.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.enterprise.gloo.solo.io.DiscoveryOverride.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.enterprise.gloo.solo.io.DiscoveryOverride} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.DiscoveryOverride.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getAuthEndpoint();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getTokenEndpoint();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getJwksUri();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getScopesList();
  if (f.length > 0) {
    writer.writeRepeatedString(
      4,
      f
    );
  }
  f = message.getResponseTypesList();
  if (f.length > 0) {
    writer.writeRepeatedString(
      5,
      f
    );
  }
  f = message.getSubjectsList();
  if (f.length > 0) {
    writer.writeRepeatedString(
      6,
      f
    );
  }
  f = message.getIdTokenAlgsList();
  if (f.length > 0) {
    writer.writeRepeatedString(
      7,
      f
    );
  }
  f = message.getAuthMethodsList();
  if (f.length > 0) {
    writer.writeRepeatedString(
      8,
      f
    );
  }
  f = message.getClaimsList();
  if (f.length > 0) {
    writer.writeRepeatedString(
      9,
      f
    );
  }
};


/**
 * optional string auth_endpoint = 1;
 * @return {string}
 */
proto.enterprise.gloo.solo.io.DiscoveryOverride.prototype.getAuthEndpoint = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.enterprise.gloo.solo.io.DiscoveryOverride.prototype.setAuthEndpoint = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string token_endpoint = 2;
 * @return {string}
 */
proto.enterprise.gloo.solo.io.DiscoveryOverride.prototype.getTokenEndpoint = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.enterprise.gloo.solo.io.DiscoveryOverride.prototype.setTokenEndpoint = function(value) {
  jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string jwks_uri = 3;
 * @return {string}
 */
proto.enterprise.gloo.solo.io.DiscoveryOverride.prototype.getJwksUri = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/** @param {string} value */
proto.enterprise.gloo.solo.io.DiscoveryOverride.prototype.setJwksUri = function(value) {
  jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * repeated string scopes = 4;
 * @return {!Array<string>}
 */
proto.enterprise.gloo.solo.io.DiscoveryOverride.prototype.getScopesList = function() {
  return /** @type {!Array<string>} */ (jspb.Message.getRepeatedField(this, 4));
};


/** @param {!Array<string>} value */
proto.enterprise.gloo.solo.io.DiscoveryOverride.prototype.setScopesList = function(value) {
  jspb.Message.setField(this, 4, value || []);
};


/**
 * @param {!string} value
 * @param {number=} opt_index
 */
proto.enterprise.gloo.solo.io.DiscoveryOverride.prototype.addScopes = function(value, opt_index) {
  jspb.Message.addToRepeatedField(this, 4, value, opt_index);
};


proto.enterprise.gloo.solo.io.DiscoveryOverride.prototype.clearScopesList = function() {
  this.setScopesList([]);
};


/**
 * repeated string response_types = 5;
 * @return {!Array<string>}
 */
proto.enterprise.gloo.solo.io.DiscoveryOverride.prototype.getResponseTypesList = function() {
  return /** @type {!Array<string>} */ (jspb.Message.getRepeatedField(this, 5));
};


/** @param {!Array<string>} value */
proto.enterprise.gloo.solo.io.DiscoveryOverride.prototype.setResponseTypesList = function(value) {
  jspb.Message.setField(this, 5, value || []);
};


/**
 * @param {!string} value
 * @param {number=} opt_index
 */
proto.enterprise.gloo.solo.io.DiscoveryOverride.prototype.addResponseTypes = function(value, opt_index) {
  jspb.Message.addToRepeatedField(this, 5, value, opt_index);
};


proto.enterprise.gloo.solo.io.DiscoveryOverride.prototype.clearResponseTypesList = function() {
  this.setResponseTypesList([]);
};


/**
 * repeated string subjects = 6;
 * @return {!Array<string>}
 */
proto.enterprise.gloo.solo.io.DiscoveryOverride.prototype.getSubjectsList = function() {
  return /** @type {!Array<string>} */ (jspb.Message.getRepeatedField(this, 6));
};


/** @param {!Array<string>} value */
proto.enterprise.gloo.solo.io.DiscoveryOverride.prototype.setSubjectsList = function(value) {
  jspb.Message.setField(this, 6, value || []);
};


/**
 * @param {!string} value
 * @param {number=} opt_index
 */
proto.enterprise.gloo.solo.io.DiscoveryOverride.prototype.addSubjects = function(value, opt_index) {
  jspb.Message.addToRepeatedField(this, 6, value, opt_index);
};


proto.enterprise.gloo.solo.io.DiscoveryOverride.prototype.clearSubjectsList = function() {
  this.setSubjectsList([]);
};


/**
 * repeated string id_token_algs = 7;
 * @return {!Array<string>}
 */
proto.enterprise.gloo.solo.io.DiscoveryOverride.prototype.getIdTokenAlgsList = function() {
  return /** @type {!Array<string>} */ (jspb.Message.getRepeatedField(this, 7));
};


/** @param {!Array<string>} value */
proto.enterprise.gloo.solo.io.DiscoveryOverride.prototype.setIdTokenAlgsList = function(value) {
  jspb.Message.setField(this, 7, value || []);
};


/**
 * @param {!string} value
 * @param {number=} opt_index
 */
proto.enterprise.gloo.solo.io.DiscoveryOverride.prototype.addIdTokenAlgs = function(value, opt_index) {
  jspb.Message.addToRepeatedField(this, 7, value, opt_index);
};


proto.enterprise.gloo.solo.io.DiscoveryOverride.prototype.clearIdTokenAlgsList = function() {
  this.setIdTokenAlgsList([]);
};


/**
 * repeated string auth_methods = 8;
 * @return {!Array<string>}
 */
proto.enterprise.gloo.solo.io.DiscoveryOverride.prototype.getAuthMethodsList = function() {
  return /** @type {!Array<string>} */ (jspb.Message.getRepeatedField(this, 8));
};


/** @param {!Array<string>} value */
proto.enterprise.gloo.solo.io.DiscoveryOverride.prototype.setAuthMethodsList = function(value) {
  jspb.Message.setField(this, 8, value || []);
};


/**
 * @param {!string} value
 * @param {number=} opt_index
 */
proto.enterprise.gloo.solo.io.DiscoveryOverride.prototype.addAuthMethods = function(value, opt_index) {
  jspb.Message.addToRepeatedField(this, 8, value, opt_index);
};


proto.enterprise.gloo.solo.io.DiscoveryOverride.prototype.clearAuthMethodsList = function() {
  this.setAuthMethodsList([]);
};


/**
 * repeated string claims = 9;
 * @return {!Array<string>}
 */
proto.enterprise.gloo.solo.io.DiscoveryOverride.prototype.getClaimsList = function() {
  return /** @type {!Array<string>} */ (jspb.Message.getRepeatedField(this, 9));
};


/** @param {!Array<string>} value */
proto.enterprise.gloo.solo.io.DiscoveryOverride.prototype.setClaimsList = function(value) {
  jspb.Message.setField(this, 9, value || []);
};


/**
 * @param {!string} value
 * @param {number=} opt_index
 */
proto.enterprise.gloo.solo.io.DiscoveryOverride.prototype.addClaims = function(value, opt_index) {
  jspb.Message.addToRepeatedField(this, 9, value, opt_index);
};


proto.enterprise.gloo.solo.io.DiscoveryOverride.prototype.clearClaimsList = function() {
  this.setClaimsList([]);
};



/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.enterprise.gloo.solo.io.OidcAuthorizationCode = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.enterprise.gloo.solo.io.OidcAuthorizationCode.repeatedFields_, null);
};
goog.inherits(proto.enterprise.gloo.solo.io.OidcAuthorizationCode, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.enterprise.gloo.solo.io.OidcAuthorizationCode.displayName = 'proto.enterprise.gloo.solo.io.OidcAuthorizationCode';
}
/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.enterprise.gloo.solo.io.OidcAuthorizationCode.repeatedFields_ = [7];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.enterprise.gloo.solo.io.OidcAuthorizationCode.prototype.toObject = function(opt_includeInstance) {
  return proto.enterprise.gloo.solo.io.OidcAuthorizationCode.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.enterprise.gloo.solo.io.OidcAuthorizationCode} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.OidcAuthorizationCode.toObject = function(includeInstance, msg) {
  var f, obj = {
    clientId: jspb.Message.getFieldWithDefault(msg, 1, ""),
    clientSecretRef: (f = msg.getClientSecretRef()) && github_com_solo$io_solo$kit_api_v1_ref_pb.ResourceRef.toObject(includeInstance, f),
    issuerUrl: jspb.Message.getFieldWithDefault(msg, 3, ""),
    authEndpointQueryParamsMap: (f = msg.getAuthEndpointQueryParamsMap()) ? f.toObject(includeInstance, undefined) : [],
    appUrl: jspb.Message.getFieldWithDefault(msg, 5, ""),
    callbackPath: jspb.Message.getFieldWithDefault(msg, 6, ""),
    logoutPath: jspb.Message.getFieldWithDefault(msg, 9, ""),
    scopesList: jspb.Message.getRepeatedField(msg, 7),
    session: (f = msg.getSession()) && proto.enterprise.gloo.solo.io.UserSession.toObject(includeInstance, f),
    headers: (f = msg.getHeaders()) && proto.enterprise.gloo.solo.io.HeaderConfiguration.toObject(includeInstance, f),
    discoveryOverride: (f = msg.getDiscoveryOverride()) && proto.enterprise.gloo.solo.io.DiscoveryOverride.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.enterprise.gloo.solo.io.OidcAuthorizationCode}
 */
proto.enterprise.gloo.solo.io.OidcAuthorizationCode.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.enterprise.gloo.solo.io.OidcAuthorizationCode;
  return proto.enterprise.gloo.solo.io.OidcAuthorizationCode.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.enterprise.gloo.solo.io.OidcAuthorizationCode} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.enterprise.gloo.solo.io.OidcAuthorizationCode}
 */
proto.enterprise.gloo.solo.io.OidcAuthorizationCode.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setClientId(value);
      break;
    case 2:
      var value = new github_com_solo$io_solo$kit_api_v1_ref_pb.ResourceRef;
      reader.readMessage(value,github_com_solo$io_solo$kit_api_v1_ref_pb.ResourceRef.deserializeBinaryFromReader);
      msg.setClientSecretRef(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setIssuerUrl(value);
      break;
    case 4:
      var value = msg.getAuthEndpointQueryParamsMap();
      reader.readMessage(value, function(message, reader) {
        jspb.Map.deserializeBinary(message, reader, jspb.BinaryReader.prototype.readString, jspb.BinaryReader.prototype.readString, null, "");
         });
      break;
    case 5:
      var value = /** @type {string} */ (reader.readString());
      msg.setAppUrl(value);
      break;
    case 6:
      var value = /** @type {string} */ (reader.readString());
      msg.setCallbackPath(value);
      break;
    case 9:
      var value = /** @type {string} */ (reader.readString());
      msg.setLogoutPath(value);
      break;
    case 7:
      var value = /** @type {string} */ (reader.readString());
      msg.addScopes(value);
      break;
    case 8:
      var value = new proto.enterprise.gloo.solo.io.UserSession;
      reader.readMessage(value,proto.enterprise.gloo.solo.io.UserSession.deserializeBinaryFromReader);
      msg.setSession(value);
      break;
    case 10:
      var value = new proto.enterprise.gloo.solo.io.HeaderConfiguration;
      reader.readMessage(value,proto.enterprise.gloo.solo.io.HeaderConfiguration.deserializeBinaryFromReader);
      msg.setHeaders(value);
      break;
    case 11:
      var value = new proto.enterprise.gloo.solo.io.DiscoveryOverride;
      reader.readMessage(value,proto.enterprise.gloo.solo.io.DiscoveryOverride.deserializeBinaryFromReader);
      msg.setDiscoveryOverride(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.enterprise.gloo.solo.io.OidcAuthorizationCode.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.enterprise.gloo.solo.io.OidcAuthorizationCode.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.enterprise.gloo.solo.io.OidcAuthorizationCode} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.OidcAuthorizationCode.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getClientId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getClientSecretRef();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      github_com_solo$io_solo$kit_api_v1_ref_pb.ResourceRef.serializeBinaryToWriter
    );
  }
  f = message.getIssuerUrl();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getAuthEndpointQueryParamsMap(true);
  if (f && f.getLength() > 0) {
    f.serializeBinary(4, writer, jspb.BinaryWriter.prototype.writeString, jspb.BinaryWriter.prototype.writeString);
  }
  f = message.getAppUrl();
  if (f.length > 0) {
    writer.writeString(
      5,
      f
    );
  }
  f = message.getCallbackPath();
  if (f.length > 0) {
    writer.writeString(
      6,
      f
    );
  }
  f = message.getLogoutPath();
  if (f.length > 0) {
    writer.writeString(
      9,
      f
    );
  }
  f = message.getScopesList();
  if (f.length > 0) {
    writer.writeRepeatedString(
      7,
      f
    );
  }
  f = message.getSession();
  if (f != null) {
    writer.writeMessage(
      8,
      f,
      proto.enterprise.gloo.solo.io.UserSession.serializeBinaryToWriter
    );
  }
  f = message.getHeaders();
  if (f != null) {
    writer.writeMessage(
      10,
      f,
      proto.enterprise.gloo.solo.io.HeaderConfiguration.serializeBinaryToWriter
    );
  }
  f = message.getDiscoveryOverride();
  if (f != null) {
    writer.writeMessage(
      11,
      f,
      proto.enterprise.gloo.solo.io.DiscoveryOverride.serializeBinaryToWriter
    );
  }
};


/**
 * optional string client_id = 1;
 * @return {string}
 */
proto.enterprise.gloo.solo.io.OidcAuthorizationCode.prototype.getClientId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.enterprise.gloo.solo.io.OidcAuthorizationCode.prototype.setClientId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional core.solo.io.ResourceRef client_secret_ref = 2;
 * @return {?proto.core.solo.io.ResourceRef}
 */
proto.enterprise.gloo.solo.io.OidcAuthorizationCode.prototype.getClientSecretRef = function() {
  return /** @type{?proto.core.solo.io.ResourceRef} */ (
    jspb.Message.getWrapperField(this, github_com_solo$io_solo$kit_api_v1_ref_pb.ResourceRef, 2));
};


/** @param {?proto.core.solo.io.ResourceRef|undefined} value */
proto.enterprise.gloo.solo.io.OidcAuthorizationCode.prototype.setClientSecretRef = function(value) {
  jspb.Message.setWrapperField(this, 2, value);
};


proto.enterprise.gloo.solo.io.OidcAuthorizationCode.prototype.clearClientSecretRef = function() {
  this.setClientSecretRef(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.enterprise.gloo.solo.io.OidcAuthorizationCode.prototype.hasClientSecretRef = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional string issuer_url = 3;
 * @return {string}
 */
proto.enterprise.gloo.solo.io.OidcAuthorizationCode.prototype.getIssuerUrl = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/** @param {string} value */
proto.enterprise.gloo.solo.io.OidcAuthorizationCode.prototype.setIssuerUrl = function(value) {
  jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * map<string, string> auth_endpoint_query_params = 4;
 * @param {boolean=} opt_noLazyCreate Do not create the map if
 * empty, instead returning `undefined`
 * @return {!jspb.Map<string,string>}
 */
proto.enterprise.gloo.solo.io.OidcAuthorizationCode.prototype.getAuthEndpointQueryParamsMap = function(opt_noLazyCreate) {
  return /** @type {!jspb.Map<string,string>} */ (
      jspb.Message.getMapField(this, 4, opt_noLazyCreate,
      null));
};


proto.enterprise.gloo.solo.io.OidcAuthorizationCode.prototype.clearAuthEndpointQueryParamsMap = function() {
  this.getAuthEndpointQueryParamsMap().clear();
};


/**
 * optional string app_url = 5;
 * @return {string}
 */
proto.enterprise.gloo.solo.io.OidcAuthorizationCode.prototype.getAppUrl = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 5, ""));
};


/** @param {string} value */
proto.enterprise.gloo.solo.io.OidcAuthorizationCode.prototype.setAppUrl = function(value) {
  jspb.Message.setProto3StringField(this, 5, value);
};


/**
 * optional string callback_path = 6;
 * @return {string}
 */
proto.enterprise.gloo.solo.io.OidcAuthorizationCode.prototype.getCallbackPath = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 6, ""));
};


/** @param {string} value */
proto.enterprise.gloo.solo.io.OidcAuthorizationCode.prototype.setCallbackPath = function(value) {
  jspb.Message.setProto3StringField(this, 6, value);
};


/**
 * optional string logout_path = 9;
 * @return {string}
 */
proto.enterprise.gloo.solo.io.OidcAuthorizationCode.prototype.getLogoutPath = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 9, ""));
};


/** @param {string} value */
proto.enterprise.gloo.solo.io.OidcAuthorizationCode.prototype.setLogoutPath = function(value) {
  jspb.Message.setProto3StringField(this, 9, value);
};


/**
 * repeated string scopes = 7;
 * @return {!Array<string>}
 */
proto.enterprise.gloo.solo.io.OidcAuthorizationCode.prototype.getScopesList = function() {
  return /** @type {!Array<string>} */ (jspb.Message.getRepeatedField(this, 7));
};


/** @param {!Array<string>} value */
proto.enterprise.gloo.solo.io.OidcAuthorizationCode.prototype.setScopesList = function(value) {
  jspb.Message.setField(this, 7, value || []);
};


/**
 * @param {!string} value
 * @param {number=} opt_index
 */
proto.enterprise.gloo.solo.io.OidcAuthorizationCode.prototype.addScopes = function(value, opt_index) {
  jspb.Message.addToRepeatedField(this, 7, value, opt_index);
};


proto.enterprise.gloo.solo.io.OidcAuthorizationCode.prototype.clearScopesList = function() {
  this.setScopesList([]);
};


/**
 * optional UserSession session = 8;
 * @return {?proto.enterprise.gloo.solo.io.UserSession}
 */
proto.enterprise.gloo.solo.io.OidcAuthorizationCode.prototype.getSession = function() {
  return /** @type{?proto.enterprise.gloo.solo.io.UserSession} */ (
    jspb.Message.getWrapperField(this, proto.enterprise.gloo.solo.io.UserSession, 8));
};


/** @param {?proto.enterprise.gloo.solo.io.UserSession|undefined} value */
proto.enterprise.gloo.solo.io.OidcAuthorizationCode.prototype.setSession = function(value) {
  jspb.Message.setWrapperField(this, 8, value);
};


proto.enterprise.gloo.solo.io.OidcAuthorizationCode.prototype.clearSession = function() {
  this.setSession(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.enterprise.gloo.solo.io.OidcAuthorizationCode.prototype.hasSession = function() {
  return jspb.Message.getField(this, 8) != null;
};


/**
 * optional HeaderConfiguration headers = 10;
 * @return {?proto.enterprise.gloo.solo.io.HeaderConfiguration}
 */
proto.enterprise.gloo.solo.io.OidcAuthorizationCode.prototype.getHeaders = function() {
  return /** @type{?proto.enterprise.gloo.solo.io.HeaderConfiguration} */ (
    jspb.Message.getWrapperField(this, proto.enterprise.gloo.solo.io.HeaderConfiguration, 10));
};


/** @param {?proto.enterprise.gloo.solo.io.HeaderConfiguration|undefined} value */
proto.enterprise.gloo.solo.io.OidcAuthorizationCode.prototype.setHeaders = function(value) {
  jspb.Message.setWrapperField(this, 10, value);
};


proto.enterprise.gloo.solo.io.OidcAuthorizationCode.prototype.clearHeaders = function() {
  this.setHeaders(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.enterprise.gloo.solo.io.OidcAuthorizationCode.prototype.hasHeaders = function() {
  return jspb.Message.getField(this, 10) != null;
};


/**
 * optional DiscoveryOverride discovery_override = 11;
 * @return {?proto.enterprise.gloo.solo.io.DiscoveryOverride}
 */
proto.enterprise.gloo.solo.io.OidcAuthorizationCode.prototype.getDiscoveryOverride = function() {
  return /** @type{?proto.enterprise.gloo.solo.io.DiscoveryOverride} */ (
    jspb.Message.getWrapperField(this, proto.enterprise.gloo.solo.io.DiscoveryOverride, 11));
};


/** @param {?proto.enterprise.gloo.solo.io.DiscoveryOverride|undefined} value */
proto.enterprise.gloo.solo.io.OidcAuthorizationCode.prototype.setDiscoveryOverride = function(value) {
  jspb.Message.setWrapperField(this, 11, value);
};


proto.enterprise.gloo.solo.io.OidcAuthorizationCode.prototype.clearDiscoveryOverride = function() {
  this.setDiscoveryOverride(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.enterprise.gloo.solo.io.OidcAuthorizationCode.prototype.hasDiscoveryOverride = function() {
  return jspb.Message.getField(this, 11) != null;
};



/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.enterprise.gloo.solo.io.AccessTokenValidation = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, proto.enterprise.gloo.solo.io.AccessTokenValidation.oneofGroups_);
};
goog.inherits(proto.enterprise.gloo.solo.io.AccessTokenValidation, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.enterprise.gloo.solo.io.AccessTokenValidation.displayName = 'proto.enterprise.gloo.solo.io.AccessTokenValidation';
}
/**
 * Oneof group definitions for this message. Each group defines the field
 * numbers belonging to that group. When of these fields' value is set, all
 * other fields in the group are cleared. During deserialization, if multiple
 * fields are encountered for a group, only the last value seen will be kept.
 * @private {!Array<!Array<number>>}
 * @const
 */
proto.enterprise.gloo.solo.io.AccessTokenValidation.oneofGroups_ = [[1],[6]];

/**
 * @enum {number}
 */
proto.enterprise.gloo.solo.io.AccessTokenValidation.ValidationTypeCase = {
  VALIDATION_TYPE_NOT_SET: 0,
  INTROSPECTION_URL: 1
};

/**
 * @return {proto.enterprise.gloo.solo.io.AccessTokenValidation.ValidationTypeCase}
 */
proto.enterprise.gloo.solo.io.AccessTokenValidation.prototype.getValidationTypeCase = function() {
  return /** @type {proto.enterprise.gloo.solo.io.AccessTokenValidation.ValidationTypeCase} */(jspb.Message.computeOneofCase(this, proto.enterprise.gloo.solo.io.AccessTokenValidation.oneofGroups_[0]));
};

/**
 * @enum {number}
 */
proto.enterprise.gloo.solo.io.AccessTokenValidation.ScopeValidationCase = {
  SCOPE_VALIDATION_NOT_SET: 0,
  REQUIRED_SCOPES: 6
};

/**
 * @return {proto.enterprise.gloo.solo.io.AccessTokenValidation.ScopeValidationCase}
 */
proto.enterprise.gloo.solo.io.AccessTokenValidation.prototype.getScopeValidationCase = function() {
  return /** @type {proto.enterprise.gloo.solo.io.AccessTokenValidation.ScopeValidationCase} */(jspb.Message.computeOneofCase(this, proto.enterprise.gloo.solo.io.AccessTokenValidation.oneofGroups_[1]));
};



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.enterprise.gloo.solo.io.AccessTokenValidation.prototype.toObject = function(opt_includeInstance) {
  return proto.enterprise.gloo.solo.io.AccessTokenValidation.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.enterprise.gloo.solo.io.AccessTokenValidation} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.AccessTokenValidation.toObject = function(includeInstance, msg) {
  var f, obj = {
    introspectionUrl: jspb.Message.getFieldWithDefault(msg, 1, ""),
    userinfoUrl: jspb.Message.getFieldWithDefault(msg, 4, ""),
    cacheTimeout: (f = msg.getCacheTimeout()) && google_protobuf_duration_pb.Duration.toObject(includeInstance, f),
    requiredScopes: (f = msg.getRequiredScopes()) && proto.enterprise.gloo.solo.io.AccessTokenValidation.ScopeList.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.enterprise.gloo.solo.io.AccessTokenValidation}
 */
proto.enterprise.gloo.solo.io.AccessTokenValidation.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.enterprise.gloo.solo.io.AccessTokenValidation;
  return proto.enterprise.gloo.solo.io.AccessTokenValidation.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.enterprise.gloo.solo.io.AccessTokenValidation} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.enterprise.gloo.solo.io.AccessTokenValidation}
 */
proto.enterprise.gloo.solo.io.AccessTokenValidation.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setIntrospectionUrl(value);
      break;
    case 4:
      var value = /** @type {string} */ (reader.readString());
      msg.setUserinfoUrl(value);
      break;
    case 5:
      var value = new google_protobuf_duration_pb.Duration;
      reader.readMessage(value,google_protobuf_duration_pb.Duration.deserializeBinaryFromReader);
      msg.setCacheTimeout(value);
      break;
    case 6:
      var value = new proto.enterprise.gloo.solo.io.AccessTokenValidation.ScopeList;
      reader.readMessage(value,proto.enterprise.gloo.solo.io.AccessTokenValidation.ScopeList.deserializeBinaryFromReader);
      msg.setRequiredScopes(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.enterprise.gloo.solo.io.AccessTokenValidation.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.enterprise.gloo.solo.io.AccessTokenValidation.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.enterprise.gloo.solo.io.AccessTokenValidation} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.AccessTokenValidation.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = /** @type {string} */ (jspb.Message.getField(message, 1));
  if (f != null) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getUserinfoUrl();
  if (f.length > 0) {
    writer.writeString(
      4,
      f
    );
  }
  f = message.getCacheTimeout();
  if (f != null) {
    writer.writeMessage(
      5,
      f,
      google_protobuf_duration_pb.Duration.serializeBinaryToWriter
    );
  }
  f = message.getRequiredScopes();
  if (f != null) {
    writer.writeMessage(
      6,
      f,
      proto.enterprise.gloo.solo.io.AccessTokenValidation.ScopeList.serializeBinaryToWriter
    );
  }
};



/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.enterprise.gloo.solo.io.AccessTokenValidation.ScopeList = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.enterprise.gloo.solo.io.AccessTokenValidation.ScopeList.repeatedFields_, null);
};
goog.inherits(proto.enterprise.gloo.solo.io.AccessTokenValidation.ScopeList, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.enterprise.gloo.solo.io.AccessTokenValidation.ScopeList.displayName = 'proto.enterprise.gloo.solo.io.AccessTokenValidation.ScopeList';
}
/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.enterprise.gloo.solo.io.AccessTokenValidation.ScopeList.repeatedFields_ = [1];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.enterprise.gloo.solo.io.AccessTokenValidation.ScopeList.prototype.toObject = function(opt_includeInstance) {
  return proto.enterprise.gloo.solo.io.AccessTokenValidation.ScopeList.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.enterprise.gloo.solo.io.AccessTokenValidation.ScopeList} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.AccessTokenValidation.ScopeList.toObject = function(includeInstance, msg) {
  var f, obj = {
    scopeList: jspb.Message.getRepeatedField(msg, 1)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.enterprise.gloo.solo.io.AccessTokenValidation.ScopeList}
 */
proto.enterprise.gloo.solo.io.AccessTokenValidation.ScopeList.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.enterprise.gloo.solo.io.AccessTokenValidation.ScopeList;
  return proto.enterprise.gloo.solo.io.AccessTokenValidation.ScopeList.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.enterprise.gloo.solo.io.AccessTokenValidation.ScopeList} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.enterprise.gloo.solo.io.AccessTokenValidation.ScopeList}
 */
proto.enterprise.gloo.solo.io.AccessTokenValidation.ScopeList.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.addScope(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.enterprise.gloo.solo.io.AccessTokenValidation.ScopeList.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.enterprise.gloo.solo.io.AccessTokenValidation.ScopeList.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.enterprise.gloo.solo.io.AccessTokenValidation.ScopeList} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.AccessTokenValidation.ScopeList.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getScopeList();
  if (f.length > 0) {
    writer.writeRepeatedString(
      1,
      f
    );
  }
};


/**
 * repeated string scope = 1;
 * @return {!Array<string>}
 */
proto.enterprise.gloo.solo.io.AccessTokenValidation.ScopeList.prototype.getScopeList = function() {
  return /** @type {!Array<string>} */ (jspb.Message.getRepeatedField(this, 1));
};


/** @param {!Array<string>} value */
proto.enterprise.gloo.solo.io.AccessTokenValidation.ScopeList.prototype.setScopeList = function(value) {
  jspb.Message.setField(this, 1, value || []);
};


/**
 * @param {!string} value
 * @param {number=} opt_index
 */
proto.enterprise.gloo.solo.io.AccessTokenValidation.ScopeList.prototype.addScope = function(value, opt_index) {
  jspb.Message.addToRepeatedField(this, 1, value, opt_index);
};


proto.enterprise.gloo.solo.io.AccessTokenValidation.ScopeList.prototype.clearScopeList = function() {
  this.setScopeList([]);
};


/**
 * optional string introspection_url = 1;
 * @return {string}
 */
proto.enterprise.gloo.solo.io.AccessTokenValidation.prototype.getIntrospectionUrl = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.enterprise.gloo.solo.io.AccessTokenValidation.prototype.setIntrospectionUrl = function(value) {
  jspb.Message.setOneofField(this, 1, proto.enterprise.gloo.solo.io.AccessTokenValidation.oneofGroups_[0], value);
};


proto.enterprise.gloo.solo.io.AccessTokenValidation.prototype.clearIntrospectionUrl = function() {
  jspb.Message.setOneofField(this, 1, proto.enterprise.gloo.solo.io.AccessTokenValidation.oneofGroups_[0], undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.enterprise.gloo.solo.io.AccessTokenValidation.prototype.hasIntrospectionUrl = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional string userinfo_url = 4;
 * @return {string}
 */
proto.enterprise.gloo.solo.io.AccessTokenValidation.prototype.getUserinfoUrl = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 4, ""));
};


/** @param {string} value */
proto.enterprise.gloo.solo.io.AccessTokenValidation.prototype.setUserinfoUrl = function(value) {
  jspb.Message.setProto3StringField(this, 4, value);
};


/**
 * optional google.protobuf.Duration cache_timeout = 5;
 * @return {?proto.google.protobuf.Duration}
 */
proto.enterprise.gloo.solo.io.AccessTokenValidation.prototype.getCacheTimeout = function() {
  return /** @type{?proto.google.protobuf.Duration} */ (
    jspb.Message.getWrapperField(this, google_protobuf_duration_pb.Duration, 5));
};


/** @param {?proto.google.protobuf.Duration|undefined} value */
proto.enterprise.gloo.solo.io.AccessTokenValidation.prototype.setCacheTimeout = function(value) {
  jspb.Message.setWrapperField(this, 5, value);
};


proto.enterprise.gloo.solo.io.AccessTokenValidation.prototype.clearCacheTimeout = function() {
  this.setCacheTimeout(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.enterprise.gloo.solo.io.AccessTokenValidation.prototype.hasCacheTimeout = function() {
  return jspb.Message.getField(this, 5) != null;
};


/**
 * optional ScopeList required_scopes = 6;
 * @return {?proto.enterprise.gloo.solo.io.AccessTokenValidation.ScopeList}
 */
proto.enterprise.gloo.solo.io.AccessTokenValidation.prototype.getRequiredScopes = function() {
  return /** @type{?proto.enterprise.gloo.solo.io.AccessTokenValidation.ScopeList} */ (
    jspb.Message.getWrapperField(this, proto.enterprise.gloo.solo.io.AccessTokenValidation.ScopeList, 6));
};


/** @param {?proto.enterprise.gloo.solo.io.AccessTokenValidation.ScopeList|undefined} value */
proto.enterprise.gloo.solo.io.AccessTokenValidation.prototype.setRequiredScopes = function(value) {
  jspb.Message.setOneofWrapperField(this, 6, proto.enterprise.gloo.solo.io.AccessTokenValidation.oneofGroups_[1], value);
};


proto.enterprise.gloo.solo.io.AccessTokenValidation.prototype.clearRequiredScopes = function() {
  this.setRequiredScopes(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.enterprise.gloo.solo.io.AccessTokenValidation.prototype.hasRequiredScopes = function() {
  return jspb.Message.getField(this, 6) != null;
};



/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.enterprise.gloo.solo.io.OauthSecret = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.enterprise.gloo.solo.io.OauthSecret, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.enterprise.gloo.solo.io.OauthSecret.displayName = 'proto.enterprise.gloo.solo.io.OauthSecret';
}


if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.enterprise.gloo.solo.io.OauthSecret.prototype.toObject = function(opt_includeInstance) {
  return proto.enterprise.gloo.solo.io.OauthSecret.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.enterprise.gloo.solo.io.OauthSecret} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.OauthSecret.toObject = function(includeInstance, msg) {
  var f, obj = {
    clientSecret: jspb.Message.getFieldWithDefault(msg, 1, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.enterprise.gloo.solo.io.OauthSecret}
 */
proto.enterprise.gloo.solo.io.OauthSecret.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.enterprise.gloo.solo.io.OauthSecret;
  return proto.enterprise.gloo.solo.io.OauthSecret.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.enterprise.gloo.solo.io.OauthSecret} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.enterprise.gloo.solo.io.OauthSecret}
 */
proto.enterprise.gloo.solo.io.OauthSecret.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setClientSecret(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.enterprise.gloo.solo.io.OauthSecret.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.enterprise.gloo.solo.io.OauthSecret.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.enterprise.gloo.solo.io.OauthSecret} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.OauthSecret.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getClientSecret();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string client_secret = 1;
 * @return {string}
 */
proto.enterprise.gloo.solo.io.OauthSecret.prototype.getClientSecret = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.enterprise.gloo.solo.io.OauthSecret.prototype.setClientSecret = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};



/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.enterprise.gloo.solo.io.ApiKeyAuth = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.enterprise.gloo.solo.io.ApiKeyAuth.repeatedFields_, null);
};
goog.inherits(proto.enterprise.gloo.solo.io.ApiKeyAuth, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.enterprise.gloo.solo.io.ApiKeyAuth.displayName = 'proto.enterprise.gloo.solo.io.ApiKeyAuth';
}
/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.enterprise.gloo.solo.io.ApiKeyAuth.repeatedFields_ = [2];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.enterprise.gloo.solo.io.ApiKeyAuth.prototype.toObject = function(opt_includeInstance) {
  return proto.enterprise.gloo.solo.io.ApiKeyAuth.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.enterprise.gloo.solo.io.ApiKeyAuth} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.ApiKeyAuth.toObject = function(includeInstance, msg) {
  var f, obj = {
    labelSelectorMap: (f = msg.getLabelSelectorMap()) ? f.toObject(includeInstance, undefined) : [],
    apiKeySecretRefsList: jspb.Message.toObjectList(msg.getApiKeySecretRefsList(),
    github_com_solo$io_solo$kit_api_v1_ref_pb.ResourceRef.toObject, includeInstance),
    headerName: jspb.Message.getFieldWithDefault(msg, 3, ""),
    headersFromMetadataMap: (f = msg.getHeadersFromMetadataMap()) ? f.toObject(includeInstance, proto.enterprise.gloo.solo.io.ApiKeyAuth.SecretKey.toObject) : []
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.enterprise.gloo.solo.io.ApiKeyAuth}
 */
proto.enterprise.gloo.solo.io.ApiKeyAuth.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.enterprise.gloo.solo.io.ApiKeyAuth;
  return proto.enterprise.gloo.solo.io.ApiKeyAuth.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.enterprise.gloo.solo.io.ApiKeyAuth} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.enterprise.gloo.solo.io.ApiKeyAuth}
 */
proto.enterprise.gloo.solo.io.ApiKeyAuth.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = msg.getLabelSelectorMap();
      reader.readMessage(value, function(message, reader) {
        jspb.Map.deserializeBinary(message, reader, jspb.BinaryReader.prototype.readString, jspb.BinaryReader.prototype.readString, null, "");
         });
      break;
    case 2:
      var value = new github_com_solo$io_solo$kit_api_v1_ref_pb.ResourceRef;
      reader.readMessage(value,github_com_solo$io_solo$kit_api_v1_ref_pb.ResourceRef.deserializeBinaryFromReader);
      msg.addApiKeySecretRefs(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setHeaderName(value);
      break;
    case 4:
      var value = msg.getHeadersFromMetadataMap();
      reader.readMessage(value, function(message, reader) {
        jspb.Map.deserializeBinary(message, reader, jspb.BinaryReader.prototype.readString, jspb.BinaryReader.prototype.readMessage, proto.enterprise.gloo.solo.io.ApiKeyAuth.SecretKey.deserializeBinaryFromReader, "");
         });
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.enterprise.gloo.solo.io.ApiKeyAuth.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.enterprise.gloo.solo.io.ApiKeyAuth.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.enterprise.gloo.solo.io.ApiKeyAuth} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.ApiKeyAuth.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getLabelSelectorMap(true);
  if (f && f.getLength() > 0) {
    f.serializeBinary(1, writer, jspb.BinaryWriter.prototype.writeString, jspb.BinaryWriter.prototype.writeString);
  }
  f = message.getApiKeySecretRefsList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      2,
      f,
      github_com_solo$io_solo$kit_api_v1_ref_pb.ResourceRef.serializeBinaryToWriter
    );
  }
  f = message.getHeaderName();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getHeadersFromMetadataMap(true);
  if (f && f.getLength() > 0) {
    f.serializeBinary(4, writer, jspb.BinaryWriter.prototype.writeString, jspb.BinaryWriter.prototype.writeMessage, proto.enterprise.gloo.solo.io.ApiKeyAuth.SecretKey.serializeBinaryToWriter);
  }
};



/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.enterprise.gloo.solo.io.ApiKeyAuth.SecretKey = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.enterprise.gloo.solo.io.ApiKeyAuth.SecretKey, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.enterprise.gloo.solo.io.ApiKeyAuth.SecretKey.displayName = 'proto.enterprise.gloo.solo.io.ApiKeyAuth.SecretKey';
}


if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.enterprise.gloo.solo.io.ApiKeyAuth.SecretKey.prototype.toObject = function(opt_includeInstance) {
  return proto.enterprise.gloo.solo.io.ApiKeyAuth.SecretKey.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.enterprise.gloo.solo.io.ApiKeyAuth.SecretKey} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.ApiKeyAuth.SecretKey.toObject = function(includeInstance, msg) {
  var f, obj = {
    name: jspb.Message.getFieldWithDefault(msg, 1, ""),
    required: jspb.Message.getFieldWithDefault(msg, 2, false)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.enterprise.gloo.solo.io.ApiKeyAuth.SecretKey}
 */
proto.enterprise.gloo.solo.io.ApiKeyAuth.SecretKey.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.enterprise.gloo.solo.io.ApiKeyAuth.SecretKey;
  return proto.enterprise.gloo.solo.io.ApiKeyAuth.SecretKey.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.enterprise.gloo.solo.io.ApiKeyAuth.SecretKey} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.enterprise.gloo.solo.io.ApiKeyAuth.SecretKey}
 */
proto.enterprise.gloo.solo.io.ApiKeyAuth.SecretKey.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setName(value);
      break;
    case 2:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setRequired(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.enterprise.gloo.solo.io.ApiKeyAuth.SecretKey.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.enterprise.gloo.solo.io.ApiKeyAuth.SecretKey.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.enterprise.gloo.solo.io.ApiKeyAuth.SecretKey} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.ApiKeyAuth.SecretKey.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getName();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getRequired();
  if (f) {
    writer.writeBool(
      2,
      f
    );
  }
};


/**
 * optional string name = 1;
 * @return {string}
 */
proto.enterprise.gloo.solo.io.ApiKeyAuth.SecretKey.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.enterprise.gloo.solo.io.ApiKeyAuth.SecretKey.prototype.setName = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional bool required = 2;
 * Note that Boolean fields may be set to 0/1 when serialized from a Java server.
 * You should avoid comparisons like {@code val === true/false} in those cases.
 * @return {boolean}
 */
proto.enterprise.gloo.solo.io.ApiKeyAuth.SecretKey.prototype.getRequired = function() {
  return /** @type {boolean} */ (jspb.Message.getFieldWithDefault(this, 2, false));
};


/** @param {boolean} value */
proto.enterprise.gloo.solo.io.ApiKeyAuth.SecretKey.prototype.setRequired = function(value) {
  jspb.Message.setProto3BooleanField(this, 2, value);
};


/**
 * map<string, string> label_selector = 1;
 * @param {boolean=} opt_noLazyCreate Do not create the map if
 * empty, instead returning `undefined`
 * @return {!jspb.Map<string,string>}
 */
proto.enterprise.gloo.solo.io.ApiKeyAuth.prototype.getLabelSelectorMap = function(opt_noLazyCreate) {
  return /** @type {!jspb.Map<string,string>} */ (
      jspb.Message.getMapField(this, 1, opt_noLazyCreate,
      null));
};


proto.enterprise.gloo.solo.io.ApiKeyAuth.prototype.clearLabelSelectorMap = function() {
  this.getLabelSelectorMap().clear();
};


/**
 * repeated core.solo.io.ResourceRef api_key_secret_refs = 2;
 * @return {!Array<!proto.core.solo.io.ResourceRef>}
 */
proto.enterprise.gloo.solo.io.ApiKeyAuth.prototype.getApiKeySecretRefsList = function() {
  return /** @type{!Array<!proto.core.solo.io.ResourceRef>} */ (
    jspb.Message.getRepeatedWrapperField(this, github_com_solo$io_solo$kit_api_v1_ref_pb.ResourceRef, 2));
};


/** @param {!Array<!proto.core.solo.io.ResourceRef>} value */
proto.enterprise.gloo.solo.io.ApiKeyAuth.prototype.setApiKeySecretRefsList = function(value) {
  jspb.Message.setRepeatedWrapperField(this, 2, value);
};


/**
 * @param {!proto.core.solo.io.ResourceRef=} opt_value
 * @param {number=} opt_index
 * @return {!proto.core.solo.io.ResourceRef}
 */
proto.enterprise.gloo.solo.io.ApiKeyAuth.prototype.addApiKeySecretRefs = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 2, opt_value, proto.core.solo.io.ResourceRef, opt_index);
};


proto.enterprise.gloo.solo.io.ApiKeyAuth.prototype.clearApiKeySecretRefsList = function() {
  this.setApiKeySecretRefsList([]);
};


/**
 * optional string header_name = 3;
 * @return {string}
 */
proto.enterprise.gloo.solo.io.ApiKeyAuth.prototype.getHeaderName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/** @param {string} value */
proto.enterprise.gloo.solo.io.ApiKeyAuth.prototype.setHeaderName = function(value) {
  jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * map<string, SecretKey> headers_from_metadata = 4;
 * @param {boolean=} opt_noLazyCreate Do not create the map if
 * empty, instead returning `undefined`
 * @return {!jspb.Map<string,!proto.enterprise.gloo.solo.io.ApiKeyAuth.SecretKey>}
 */
proto.enterprise.gloo.solo.io.ApiKeyAuth.prototype.getHeadersFromMetadataMap = function(opt_noLazyCreate) {
  return /** @type {!jspb.Map<string,!proto.enterprise.gloo.solo.io.ApiKeyAuth.SecretKey>} */ (
      jspb.Message.getMapField(this, 4, opt_noLazyCreate,
      proto.enterprise.gloo.solo.io.ApiKeyAuth.SecretKey));
};


proto.enterprise.gloo.solo.io.ApiKeyAuth.prototype.clearHeadersFromMetadataMap = function() {
  this.getHeadersFromMetadataMap().clear();
};



/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.enterprise.gloo.solo.io.ApiKeySecret = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.enterprise.gloo.solo.io.ApiKeySecret.repeatedFields_, null);
};
goog.inherits(proto.enterprise.gloo.solo.io.ApiKeySecret, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.enterprise.gloo.solo.io.ApiKeySecret.displayName = 'proto.enterprise.gloo.solo.io.ApiKeySecret';
}
/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.enterprise.gloo.solo.io.ApiKeySecret.repeatedFields_ = [3];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.enterprise.gloo.solo.io.ApiKeySecret.prototype.toObject = function(opt_includeInstance) {
  return proto.enterprise.gloo.solo.io.ApiKeySecret.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.enterprise.gloo.solo.io.ApiKeySecret} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.ApiKeySecret.toObject = function(includeInstance, msg) {
  var f, obj = {
    generateApiKey: jspb.Message.getFieldWithDefault(msg, 1, false),
    apiKey: jspb.Message.getFieldWithDefault(msg, 2, ""),
    labelsList: jspb.Message.getRepeatedField(msg, 3),
    metadataMap: (f = msg.getMetadataMap()) ? f.toObject(includeInstance, undefined) : []
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.enterprise.gloo.solo.io.ApiKeySecret}
 */
proto.enterprise.gloo.solo.io.ApiKeySecret.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.enterprise.gloo.solo.io.ApiKeySecret;
  return proto.enterprise.gloo.solo.io.ApiKeySecret.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.enterprise.gloo.solo.io.ApiKeySecret} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.enterprise.gloo.solo.io.ApiKeySecret}
 */
proto.enterprise.gloo.solo.io.ApiKeySecret.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setGenerateApiKey(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setApiKey(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.addLabels(value);
      break;
    case 4:
      var value = msg.getMetadataMap();
      reader.readMessage(value, function(message, reader) {
        jspb.Map.deserializeBinary(message, reader, jspb.BinaryReader.prototype.readString, jspb.BinaryReader.prototype.readString, null, "");
         });
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.enterprise.gloo.solo.io.ApiKeySecret.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.enterprise.gloo.solo.io.ApiKeySecret.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.enterprise.gloo.solo.io.ApiKeySecret} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.ApiKeySecret.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getGenerateApiKey();
  if (f) {
    writer.writeBool(
      1,
      f
    );
  }
  f = message.getApiKey();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getLabelsList();
  if (f.length > 0) {
    writer.writeRepeatedString(
      3,
      f
    );
  }
  f = message.getMetadataMap(true);
  if (f && f.getLength() > 0) {
    f.serializeBinary(4, writer, jspb.BinaryWriter.prototype.writeString, jspb.BinaryWriter.prototype.writeString);
  }
};


/**
 * optional bool generate_api_key = 1;
 * Note that Boolean fields may be set to 0/1 when serialized from a Java server.
 * You should avoid comparisons like {@code val === true/false} in those cases.
 * @return {boolean}
 */
proto.enterprise.gloo.solo.io.ApiKeySecret.prototype.getGenerateApiKey = function() {
  return /** @type {boolean} */ (jspb.Message.getFieldWithDefault(this, 1, false));
};


/** @param {boolean} value */
proto.enterprise.gloo.solo.io.ApiKeySecret.prototype.setGenerateApiKey = function(value) {
  jspb.Message.setProto3BooleanField(this, 1, value);
};


/**
 * optional string api_key = 2;
 * @return {string}
 */
proto.enterprise.gloo.solo.io.ApiKeySecret.prototype.getApiKey = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.enterprise.gloo.solo.io.ApiKeySecret.prototype.setApiKey = function(value) {
  jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * repeated string labels = 3;
 * @return {!Array<string>}
 */
proto.enterprise.gloo.solo.io.ApiKeySecret.prototype.getLabelsList = function() {
  return /** @type {!Array<string>} */ (jspb.Message.getRepeatedField(this, 3));
};


/** @param {!Array<string>} value */
proto.enterprise.gloo.solo.io.ApiKeySecret.prototype.setLabelsList = function(value) {
  jspb.Message.setField(this, 3, value || []);
};


/**
 * @param {!string} value
 * @param {number=} opt_index
 */
proto.enterprise.gloo.solo.io.ApiKeySecret.prototype.addLabels = function(value, opt_index) {
  jspb.Message.addToRepeatedField(this, 3, value, opt_index);
};


proto.enterprise.gloo.solo.io.ApiKeySecret.prototype.clearLabelsList = function() {
  this.setLabelsList([]);
};


/**
 * map<string, string> metadata = 4;
 * @param {boolean=} opt_noLazyCreate Do not create the map if
 * empty, instead returning `undefined`
 * @return {!jspb.Map<string,string>}
 */
proto.enterprise.gloo.solo.io.ApiKeySecret.prototype.getMetadataMap = function(opt_noLazyCreate) {
  return /** @type {!jspb.Map<string,string>} */ (
      jspb.Message.getMapField(this, 4, opt_noLazyCreate,
      null));
};


proto.enterprise.gloo.solo.io.ApiKeySecret.prototype.clearMetadataMap = function() {
  this.getMetadataMap().clear();
};



/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.enterprise.gloo.solo.io.OpaAuth = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.enterprise.gloo.solo.io.OpaAuth.repeatedFields_, null);
};
goog.inherits(proto.enterprise.gloo.solo.io.OpaAuth, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.enterprise.gloo.solo.io.OpaAuth.displayName = 'proto.enterprise.gloo.solo.io.OpaAuth';
}
/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.enterprise.gloo.solo.io.OpaAuth.repeatedFields_ = [1];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.enterprise.gloo.solo.io.OpaAuth.prototype.toObject = function(opt_includeInstance) {
  return proto.enterprise.gloo.solo.io.OpaAuth.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.enterprise.gloo.solo.io.OpaAuth} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.OpaAuth.toObject = function(includeInstance, msg) {
  var f, obj = {
    modulesList: jspb.Message.toObjectList(msg.getModulesList(),
    github_com_solo$io_solo$kit_api_v1_ref_pb.ResourceRef.toObject, includeInstance),
    query: jspb.Message.getFieldWithDefault(msg, 2, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.enterprise.gloo.solo.io.OpaAuth}
 */
proto.enterprise.gloo.solo.io.OpaAuth.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.enterprise.gloo.solo.io.OpaAuth;
  return proto.enterprise.gloo.solo.io.OpaAuth.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.enterprise.gloo.solo.io.OpaAuth} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.enterprise.gloo.solo.io.OpaAuth}
 */
proto.enterprise.gloo.solo.io.OpaAuth.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new github_com_solo$io_solo$kit_api_v1_ref_pb.ResourceRef;
      reader.readMessage(value,github_com_solo$io_solo$kit_api_v1_ref_pb.ResourceRef.deserializeBinaryFromReader);
      msg.addModules(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setQuery(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.enterprise.gloo.solo.io.OpaAuth.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.enterprise.gloo.solo.io.OpaAuth.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.enterprise.gloo.solo.io.OpaAuth} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.OpaAuth.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getModulesList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      1,
      f,
      github_com_solo$io_solo$kit_api_v1_ref_pb.ResourceRef.serializeBinaryToWriter
    );
  }
  f = message.getQuery();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
};


/**
 * repeated core.solo.io.ResourceRef modules = 1;
 * @return {!Array<!proto.core.solo.io.ResourceRef>}
 */
proto.enterprise.gloo.solo.io.OpaAuth.prototype.getModulesList = function() {
  return /** @type{!Array<!proto.core.solo.io.ResourceRef>} */ (
    jspb.Message.getRepeatedWrapperField(this, github_com_solo$io_solo$kit_api_v1_ref_pb.ResourceRef, 1));
};


/** @param {!Array<!proto.core.solo.io.ResourceRef>} value */
proto.enterprise.gloo.solo.io.OpaAuth.prototype.setModulesList = function(value) {
  jspb.Message.setRepeatedWrapperField(this, 1, value);
};


/**
 * @param {!proto.core.solo.io.ResourceRef=} opt_value
 * @param {number=} opt_index
 * @return {!proto.core.solo.io.ResourceRef}
 */
proto.enterprise.gloo.solo.io.OpaAuth.prototype.addModules = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 1, opt_value, proto.core.solo.io.ResourceRef, opt_index);
};


proto.enterprise.gloo.solo.io.OpaAuth.prototype.clearModulesList = function() {
  this.setModulesList([]);
};


/**
 * optional string query = 2;
 * @return {string}
 */
proto.enterprise.gloo.solo.io.OpaAuth.prototype.getQuery = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.enterprise.gloo.solo.io.OpaAuth.prototype.setQuery = function(value) {
  jspb.Message.setProto3StringField(this, 2, value);
};



/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.enterprise.gloo.solo.io.Ldap = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.enterprise.gloo.solo.io.Ldap.repeatedFields_, null);
};
goog.inherits(proto.enterprise.gloo.solo.io.Ldap, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.enterprise.gloo.solo.io.Ldap.displayName = 'proto.enterprise.gloo.solo.io.Ldap';
}
/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.enterprise.gloo.solo.io.Ldap.repeatedFields_ = [4];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.enterprise.gloo.solo.io.Ldap.prototype.toObject = function(opt_includeInstance) {
  return proto.enterprise.gloo.solo.io.Ldap.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.enterprise.gloo.solo.io.Ldap} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.Ldap.toObject = function(includeInstance, msg) {
  var f, obj = {
    address: jspb.Message.getFieldWithDefault(msg, 1, ""),
    userdntemplate: jspb.Message.getFieldWithDefault(msg, 2, ""),
    membershipattributename: jspb.Message.getFieldWithDefault(msg, 3, ""),
    allowedgroupsList: jspb.Message.getRepeatedField(msg, 4),
    pool: (f = msg.getPool()) && proto.enterprise.gloo.solo.io.Ldap.ConnectionPool.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.enterprise.gloo.solo.io.Ldap}
 */
proto.enterprise.gloo.solo.io.Ldap.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.enterprise.gloo.solo.io.Ldap;
  return proto.enterprise.gloo.solo.io.Ldap.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.enterprise.gloo.solo.io.Ldap} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.enterprise.gloo.solo.io.Ldap}
 */
proto.enterprise.gloo.solo.io.Ldap.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setAddress(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setUserdntemplate(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setMembershipattributename(value);
      break;
    case 4:
      var value = /** @type {string} */ (reader.readString());
      msg.addAllowedgroups(value);
      break;
    case 5:
      var value = new proto.enterprise.gloo.solo.io.Ldap.ConnectionPool;
      reader.readMessage(value,proto.enterprise.gloo.solo.io.Ldap.ConnectionPool.deserializeBinaryFromReader);
      msg.setPool(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.enterprise.gloo.solo.io.Ldap.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.enterprise.gloo.solo.io.Ldap.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.enterprise.gloo.solo.io.Ldap} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.Ldap.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getAddress();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getUserdntemplate();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getMembershipattributename();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getAllowedgroupsList();
  if (f.length > 0) {
    writer.writeRepeatedString(
      4,
      f
    );
  }
  f = message.getPool();
  if (f != null) {
    writer.writeMessage(
      5,
      f,
      proto.enterprise.gloo.solo.io.Ldap.ConnectionPool.serializeBinaryToWriter
    );
  }
};



/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.enterprise.gloo.solo.io.Ldap.ConnectionPool = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.enterprise.gloo.solo.io.Ldap.ConnectionPool, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.enterprise.gloo.solo.io.Ldap.ConnectionPool.displayName = 'proto.enterprise.gloo.solo.io.Ldap.ConnectionPool';
}


if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.enterprise.gloo.solo.io.Ldap.ConnectionPool.prototype.toObject = function(opt_includeInstance) {
  return proto.enterprise.gloo.solo.io.Ldap.ConnectionPool.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.enterprise.gloo.solo.io.Ldap.ConnectionPool} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.Ldap.ConnectionPool.toObject = function(includeInstance, msg) {
  var f, obj = {
    maxsize: (f = msg.getMaxsize()) && google_protobuf_wrappers_pb.UInt32Value.toObject(includeInstance, f),
    initialsize: (f = msg.getInitialsize()) && google_protobuf_wrappers_pb.UInt32Value.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.enterprise.gloo.solo.io.Ldap.ConnectionPool}
 */
proto.enterprise.gloo.solo.io.Ldap.ConnectionPool.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.enterprise.gloo.solo.io.Ldap.ConnectionPool;
  return proto.enterprise.gloo.solo.io.Ldap.ConnectionPool.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.enterprise.gloo.solo.io.Ldap.ConnectionPool} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.enterprise.gloo.solo.io.Ldap.ConnectionPool}
 */
proto.enterprise.gloo.solo.io.Ldap.ConnectionPool.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new google_protobuf_wrappers_pb.UInt32Value;
      reader.readMessage(value,google_protobuf_wrappers_pb.UInt32Value.deserializeBinaryFromReader);
      msg.setMaxsize(value);
      break;
    case 2:
      var value = new google_protobuf_wrappers_pb.UInt32Value;
      reader.readMessage(value,google_protobuf_wrappers_pb.UInt32Value.deserializeBinaryFromReader);
      msg.setInitialsize(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.enterprise.gloo.solo.io.Ldap.ConnectionPool.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.enterprise.gloo.solo.io.Ldap.ConnectionPool.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.enterprise.gloo.solo.io.Ldap.ConnectionPool} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.Ldap.ConnectionPool.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getMaxsize();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      google_protobuf_wrappers_pb.UInt32Value.serializeBinaryToWriter
    );
  }
  f = message.getInitialsize();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      google_protobuf_wrappers_pb.UInt32Value.serializeBinaryToWriter
    );
  }
};


/**
 * optional google.protobuf.UInt32Value maxSize = 1;
 * @return {?proto.google.protobuf.UInt32Value}
 */
proto.enterprise.gloo.solo.io.Ldap.ConnectionPool.prototype.getMaxsize = function() {
  return /** @type{?proto.google.protobuf.UInt32Value} */ (
    jspb.Message.getWrapperField(this, google_protobuf_wrappers_pb.UInt32Value, 1));
};


/** @param {?proto.google.protobuf.UInt32Value|undefined} value */
proto.enterprise.gloo.solo.io.Ldap.ConnectionPool.prototype.setMaxsize = function(value) {
  jspb.Message.setWrapperField(this, 1, value);
};


proto.enterprise.gloo.solo.io.Ldap.ConnectionPool.prototype.clearMaxsize = function() {
  this.setMaxsize(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.enterprise.gloo.solo.io.Ldap.ConnectionPool.prototype.hasMaxsize = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional google.protobuf.UInt32Value initialSize = 2;
 * @return {?proto.google.protobuf.UInt32Value}
 */
proto.enterprise.gloo.solo.io.Ldap.ConnectionPool.prototype.getInitialsize = function() {
  return /** @type{?proto.google.protobuf.UInt32Value} */ (
    jspb.Message.getWrapperField(this, google_protobuf_wrappers_pb.UInt32Value, 2));
};


/** @param {?proto.google.protobuf.UInt32Value|undefined} value */
proto.enterprise.gloo.solo.io.Ldap.ConnectionPool.prototype.setInitialsize = function(value) {
  jspb.Message.setWrapperField(this, 2, value);
};


proto.enterprise.gloo.solo.io.Ldap.ConnectionPool.prototype.clearInitialsize = function() {
  this.setInitialsize(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.enterprise.gloo.solo.io.Ldap.ConnectionPool.prototype.hasInitialsize = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional string address = 1;
 * @return {string}
 */
proto.enterprise.gloo.solo.io.Ldap.prototype.getAddress = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.enterprise.gloo.solo.io.Ldap.prototype.setAddress = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string userDnTemplate = 2;
 * @return {string}
 */
proto.enterprise.gloo.solo.io.Ldap.prototype.getUserdntemplate = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.enterprise.gloo.solo.io.Ldap.prototype.setUserdntemplate = function(value) {
  jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string membershipAttributeName = 3;
 * @return {string}
 */
proto.enterprise.gloo.solo.io.Ldap.prototype.getMembershipattributename = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/** @param {string} value */
proto.enterprise.gloo.solo.io.Ldap.prototype.setMembershipattributename = function(value) {
  jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * repeated string allowedGroups = 4;
 * @return {!Array<string>}
 */
proto.enterprise.gloo.solo.io.Ldap.prototype.getAllowedgroupsList = function() {
  return /** @type {!Array<string>} */ (jspb.Message.getRepeatedField(this, 4));
};


/** @param {!Array<string>} value */
proto.enterprise.gloo.solo.io.Ldap.prototype.setAllowedgroupsList = function(value) {
  jspb.Message.setField(this, 4, value || []);
};


/**
 * @param {!string} value
 * @param {number=} opt_index
 */
proto.enterprise.gloo.solo.io.Ldap.prototype.addAllowedgroups = function(value, opt_index) {
  jspb.Message.addToRepeatedField(this, 4, value, opt_index);
};


proto.enterprise.gloo.solo.io.Ldap.prototype.clearAllowedgroupsList = function() {
  this.setAllowedgroupsList([]);
};


/**
 * optional ConnectionPool pool = 5;
 * @return {?proto.enterprise.gloo.solo.io.Ldap.ConnectionPool}
 */
proto.enterprise.gloo.solo.io.Ldap.prototype.getPool = function() {
  return /** @type{?proto.enterprise.gloo.solo.io.Ldap.ConnectionPool} */ (
    jspb.Message.getWrapperField(this, proto.enterprise.gloo.solo.io.Ldap.ConnectionPool, 5));
};


/** @param {?proto.enterprise.gloo.solo.io.Ldap.ConnectionPool|undefined} value */
proto.enterprise.gloo.solo.io.Ldap.prototype.setPool = function(value) {
  jspb.Message.setWrapperField(this, 5, value);
};


proto.enterprise.gloo.solo.io.Ldap.prototype.clearPool = function() {
  this.setPool(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.enterprise.gloo.solo.io.Ldap.prototype.hasPool = function() {
  return jspb.Message.getField(this, 5) != null;
};



/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.enterprise.gloo.solo.io.PassThroughAuth = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, proto.enterprise.gloo.solo.io.PassThroughAuth.oneofGroups_);
};
goog.inherits(proto.enterprise.gloo.solo.io.PassThroughAuth, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.enterprise.gloo.solo.io.PassThroughAuth.displayName = 'proto.enterprise.gloo.solo.io.PassThroughAuth';
}
/**
 * Oneof group definitions for this message. Each group defines the field
 * numbers belonging to that group. When of these fields' value is set, all
 * other fields in the group are cleared. During deserialization, if multiple
 * fields are encountered for a group, only the last value seen will be kept.
 * @private {!Array<!Array<number>>}
 * @const
 */
proto.enterprise.gloo.solo.io.PassThroughAuth.oneofGroups_ = [[1]];

/**
 * @enum {number}
 */
proto.enterprise.gloo.solo.io.PassThroughAuth.ProtocolCase = {
  PROTOCOL_NOT_SET: 0,
  GRPC: 1
};

/**
 * @return {proto.enterprise.gloo.solo.io.PassThroughAuth.ProtocolCase}
 */
proto.enterprise.gloo.solo.io.PassThroughAuth.prototype.getProtocolCase = function() {
  return /** @type {proto.enterprise.gloo.solo.io.PassThroughAuth.ProtocolCase} */(jspb.Message.computeOneofCase(this, proto.enterprise.gloo.solo.io.PassThroughAuth.oneofGroups_[0]));
};



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.enterprise.gloo.solo.io.PassThroughAuth.prototype.toObject = function(opt_includeInstance) {
  return proto.enterprise.gloo.solo.io.PassThroughAuth.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.enterprise.gloo.solo.io.PassThroughAuth} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.PassThroughAuth.toObject = function(includeInstance, msg) {
  var f, obj = {
    grpc: (f = msg.getGrpc()) && proto.enterprise.gloo.solo.io.PassThroughGrpc.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.enterprise.gloo.solo.io.PassThroughAuth}
 */
proto.enterprise.gloo.solo.io.PassThroughAuth.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.enterprise.gloo.solo.io.PassThroughAuth;
  return proto.enterprise.gloo.solo.io.PassThroughAuth.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.enterprise.gloo.solo.io.PassThroughAuth} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.enterprise.gloo.solo.io.PassThroughAuth}
 */
proto.enterprise.gloo.solo.io.PassThroughAuth.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.enterprise.gloo.solo.io.PassThroughGrpc;
      reader.readMessage(value,proto.enterprise.gloo.solo.io.PassThroughGrpc.deserializeBinaryFromReader);
      msg.setGrpc(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.enterprise.gloo.solo.io.PassThroughAuth.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.enterprise.gloo.solo.io.PassThroughAuth.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.enterprise.gloo.solo.io.PassThroughAuth} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.PassThroughAuth.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getGrpc();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto.enterprise.gloo.solo.io.PassThroughGrpc.serializeBinaryToWriter
    );
  }
};


/**
 * optional PassThroughGrpc grpc = 1;
 * @return {?proto.enterprise.gloo.solo.io.PassThroughGrpc}
 */
proto.enterprise.gloo.solo.io.PassThroughAuth.prototype.getGrpc = function() {
  return /** @type{?proto.enterprise.gloo.solo.io.PassThroughGrpc} */ (
    jspb.Message.getWrapperField(this, proto.enterprise.gloo.solo.io.PassThroughGrpc, 1));
};


/** @param {?proto.enterprise.gloo.solo.io.PassThroughGrpc|undefined} value */
proto.enterprise.gloo.solo.io.PassThroughAuth.prototype.setGrpc = function(value) {
  jspb.Message.setOneofWrapperField(this, 1, proto.enterprise.gloo.solo.io.PassThroughAuth.oneofGroups_[0], value);
};


proto.enterprise.gloo.solo.io.PassThroughAuth.prototype.clearGrpc = function() {
  this.setGrpc(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.enterprise.gloo.solo.io.PassThroughAuth.prototype.hasGrpc = function() {
  return jspb.Message.getField(this, 1) != null;
};



/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.enterprise.gloo.solo.io.PassThroughGrpc = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.enterprise.gloo.solo.io.PassThroughGrpc, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.enterprise.gloo.solo.io.PassThroughGrpc.displayName = 'proto.enterprise.gloo.solo.io.PassThroughGrpc';
}


if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.enterprise.gloo.solo.io.PassThroughGrpc.prototype.toObject = function(opt_includeInstance) {
  return proto.enterprise.gloo.solo.io.PassThroughGrpc.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.enterprise.gloo.solo.io.PassThroughGrpc} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.PassThroughGrpc.toObject = function(includeInstance, msg) {
  var f, obj = {
    address: jspb.Message.getFieldWithDefault(msg, 1, ""),
    connectionTimeout: (f = msg.getConnectionTimeout()) && google_protobuf_duration_pb.Duration.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.enterprise.gloo.solo.io.PassThroughGrpc}
 */
proto.enterprise.gloo.solo.io.PassThroughGrpc.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.enterprise.gloo.solo.io.PassThroughGrpc;
  return proto.enterprise.gloo.solo.io.PassThroughGrpc.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.enterprise.gloo.solo.io.PassThroughGrpc} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.enterprise.gloo.solo.io.PassThroughGrpc}
 */
proto.enterprise.gloo.solo.io.PassThroughGrpc.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setAddress(value);
      break;
    case 2:
      var value = new google_protobuf_duration_pb.Duration;
      reader.readMessage(value,google_protobuf_duration_pb.Duration.deserializeBinaryFromReader);
      msg.setConnectionTimeout(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.enterprise.gloo.solo.io.PassThroughGrpc.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.enterprise.gloo.solo.io.PassThroughGrpc.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.enterprise.gloo.solo.io.PassThroughGrpc} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.PassThroughGrpc.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getAddress();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getConnectionTimeout();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      google_protobuf_duration_pb.Duration.serializeBinaryToWriter
    );
  }
};


/**
 * optional string address = 1;
 * @return {string}
 */
proto.enterprise.gloo.solo.io.PassThroughGrpc.prototype.getAddress = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.enterprise.gloo.solo.io.PassThroughGrpc.prototype.setAddress = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional google.protobuf.Duration connection_timeout = 2;
 * @return {?proto.google.protobuf.Duration}
 */
proto.enterprise.gloo.solo.io.PassThroughGrpc.prototype.getConnectionTimeout = function() {
  return /** @type{?proto.google.protobuf.Duration} */ (
    jspb.Message.getWrapperField(this, google_protobuf_duration_pb.Duration, 2));
};


/** @param {?proto.google.protobuf.Duration|undefined} value */
proto.enterprise.gloo.solo.io.PassThroughGrpc.prototype.setConnectionTimeout = function(value) {
  jspb.Message.setWrapperField(this, 2, value);
};


proto.enterprise.gloo.solo.io.PassThroughGrpc.prototype.clearConnectionTimeout = function() {
  this.setConnectionTimeout(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.enterprise.gloo.solo.io.PassThroughGrpc.prototype.hasConnectionTimeout = function() {
  return jspb.Message.getField(this, 2) != null;
};



/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.enterprise.gloo.solo.io.ExtAuthConfig.repeatedFields_, null);
};
goog.inherits(proto.enterprise.gloo.solo.io.ExtAuthConfig, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.enterprise.gloo.solo.io.ExtAuthConfig.displayName = 'proto.enterprise.gloo.solo.io.ExtAuthConfig';
}
/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.repeatedFields_ = [8];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.prototype.toObject = function(opt_includeInstance) {
  return proto.enterprise.gloo.solo.io.ExtAuthConfig.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.enterprise.gloo.solo.io.ExtAuthConfig} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.toObject = function(includeInstance, msg) {
  var f, obj = {
    authConfigRefName: jspb.Message.getFieldWithDefault(msg, 1, ""),
    configsList: jspb.Message.toObjectList(msg.getConfigsList(),
    proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.toObject, includeInstance),
    booleanExpr: (f = msg.getBooleanExpr()) && google_protobuf_wrappers_pb.StringValue.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.enterprise.gloo.solo.io.ExtAuthConfig}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.enterprise.gloo.solo.io.ExtAuthConfig;
  return proto.enterprise.gloo.solo.io.ExtAuthConfig.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.enterprise.gloo.solo.io.ExtAuthConfig} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.enterprise.gloo.solo.io.ExtAuthConfig}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setAuthConfigRefName(value);
      break;
    case 8:
      var value = new proto.enterprise.gloo.solo.io.ExtAuthConfig.Config;
      reader.readMessage(value,proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.deserializeBinaryFromReader);
      msg.addConfigs(value);
      break;
    case 10:
      var value = new google_protobuf_wrappers_pb.StringValue;
      reader.readMessage(value,google_protobuf_wrappers_pb.StringValue.deserializeBinaryFromReader);
      msg.setBooleanExpr(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.enterprise.gloo.solo.io.ExtAuthConfig.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.enterprise.gloo.solo.io.ExtAuthConfig} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getAuthConfigRefName();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getConfigsList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      8,
      f,
      proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.serializeBinaryToWriter
    );
  }
  f = message.getBooleanExpr();
  if (f != null) {
    writer.writeMessage(
      10,
      f,
      google_protobuf_wrappers_pb.StringValue.serializeBinaryToWriter
    );
  }
};



/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuthConfig = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuthConfig.repeatedFields_, null);
};
goog.inherits(proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuthConfig, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuthConfig.displayName = 'proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuthConfig';
}
/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuthConfig.repeatedFields_ = [6];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuthConfig.prototype.toObject = function(opt_includeInstance) {
  return proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuthConfig.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuthConfig} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuthConfig.toObject = function(includeInstance, msg) {
  var f, obj = {
    clientId: jspb.Message.getFieldWithDefault(msg, 1, ""),
    clientSecret: jspb.Message.getFieldWithDefault(msg, 2, ""),
    issuerUrl: jspb.Message.getFieldWithDefault(msg, 3, ""),
    authEndpointQueryParamsMap: (f = msg.getAuthEndpointQueryParamsMap()) ? f.toObject(includeInstance, undefined) : [],
    appUrl: jspb.Message.getFieldWithDefault(msg, 4, ""),
    callbackPath: jspb.Message.getFieldWithDefault(msg, 5, ""),
    scopesList: jspb.Message.getRepeatedField(msg, 6)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuthConfig}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuthConfig.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuthConfig;
  return proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuthConfig.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuthConfig} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuthConfig}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuthConfig.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setClientId(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setClientSecret(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setIssuerUrl(value);
      break;
    case 7:
      var value = msg.getAuthEndpointQueryParamsMap();
      reader.readMessage(value, function(message, reader) {
        jspb.Map.deserializeBinary(message, reader, jspb.BinaryReader.prototype.readString, jspb.BinaryReader.prototype.readString, null, "");
         });
      break;
    case 4:
      var value = /** @type {string} */ (reader.readString());
      msg.setAppUrl(value);
      break;
    case 5:
      var value = /** @type {string} */ (reader.readString());
      msg.setCallbackPath(value);
      break;
    case 6:
      var value = /** @type {string} */ (reader.readString());
      msg.addScopes(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuthConfig.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuthConfig.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuthConfig} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuthConfig.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getClientId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getClientSecret();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getIssuerUrl();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getAuthEndpointQueryParamsMap(true);
  if (f && f.getLength() > 0) {
    f.serializeBinary(7, writer, jspb.BinaryWriter.prototype.writeString, jspb.BinaryWriter.prototype.writeString);
  }
  f = message.getAppUrl();
  if (f.length > 0) {
    writer.writeString(
      4,
      f
    );
  }
  f = message.getCallbackPath();
  if (f.length > 0) {
    writer.writeString(
      5,
      f
    );
  }
  f = message.getScopesList();
  if (f.length > 0) {
    writer.writeRepeatedString(
      6,
      f
    );
  }
};


/**
 * optional string client_id = 1;
 * @return {string}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuthConfig.prototype.getClientId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuthConfig.prototype.setClientId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string client_secret = 2;
 * @return {string}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuthConfig.prototype.getClientSecret = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuthConfig.prototype.setClientSecret = function(value) {
  jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string issuer_url = 3;
 * @return {string}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuthConfig.prototype.getIssuerUrl = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/** @param {string} value */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuthConfig.prototype.setIssuerUrl = function(value) {
  jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * map<string, string> auth_endpoint_query_params = 7;
 * @param {boolean=} opt_noLazyCreate Do not create the map if
 * empty, instead returning `undefined`
 * @return {!jspb.Map<string,string>}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuthConfig.prototype.getAuthEndpointQueryParamsMap = function(opt_noLazyCreate) {
  return /** @type {!jspb.Map<string,string>} */ (
      jspb.Message.getMapField(this, 7, opt_noLazyCreate,
      null));
};


proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuthConfig.prototype.clearAuthEndpointQueryParamsMap = function() {
  this.getAuthEndpointQueryParamsMap().clear();
};


/**
 * optional string app_url = 4;
 * @return {string}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuthConfig.prototype.getAppUrl = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 4, ""));
};


/** @param {string} value */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuthConfig.prototype.setAppUrl = function(value) {
  jspb.Message.setProto3StringField(this, 4, value);
};


/**
 * optional string callback_path = 5;
 * @return {string}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuthConfig.prototype.getCallbackPath = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 5, ""));
};


/** @param {string} value */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuthConfig.prototype.setCallbackPath = function(value) {
  jspb.Message.setProto3StringField(this, 5, value);
};


/**
 * repeated string scopes = 6;
 * @return {!Array<string>}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuthConfig.prototype.getScopesList = function() {
  return /** @type {!Array<string>} */ (jspb.Message.getRepeatedField(this, 6));
};


/** @param {!Array<string>} value */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuthConfig.prototype.setScopesList = function(value) {
  jspb.Message.setField(this, 6, value || []);
};


/**
 * @param {!string} value
 * @param {number=} opt_index
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuthConfig.prototype.addScopes = function(value, opt_index) {
  jspb.Message.addToRepeatedField(this, 6, value, opt_index);
};


proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuthConfig.prototype.clearScopesList = function() {
  this.setScopesList([]);
};



/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OidcAuthorizationCodeConfig = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.enterprise.gloo.solo.io.ExtAuthConfig.OidcAuthorizationCodeConfig.repeatedFields_, null);
};
goog.inherits(proto.enterprise.gloo.solo.io.ExtAuthConfig.OidcAuthorizationCodeConfig, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.enterprise.gloo.solo.io.ExtAuthConfig.OidcAuthorizationCodeConfig.displayName = 'proto.enterprise.gloo.solo.io.ExtAuthConfig.OidcAuthorizationCodeConfig';
}
/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OidcAuthorizationCodeConfig.repeatedFields_ = [7];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OidcAuthorizationCodeConfig.prototype.toObject = function(opt_includeInstance) {
  return proto.enterprise.gloo.solo.io.ExtAuthConfig.OidcAuthorizationCodeConfig.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.enterprise.gloo.solo.io.ExtAuthConfig.OidcAuthorizationCodeConfig} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OidcAuthorizationCodeConfig.toObject = function(includeInstance, msg) {
  var f, obj = {
    clientId: jspb.Message.getFieldWithDefault(msg, 1, ""),
    clientSecret: jspb.Message.getFieldWithDefault(msg, 2, ""),
    issuerUrl: jspb.Message.getFieldWithDefault(msg, 3, ""),
    authEndpointQueryParamsMap: (f = msg.getAuthEndpointQueryParamsMap()) ? f.toObject(includeInstance, undefined) : [],
    appUrl: jspb.Message.getFieldWithDefault(msg, 5, ""),
    callbackPath: jspb.Message.getFieldWithDefault(msg, 6, ""),
    logoutPath: jspb.Message.getFieldWithDefault(msg, 9, ""),
    scopesList: jspb.Message.getRepeatedField(msg, 7),
    session: (f = msg.getSession()) && proto.enterprise.gloo.solo.io.UserSession.toObject(includeInstance, f),
    headers: (f = msg.getHeaders()) && proto.enterprise.gloo.solo.io.HeaderConfiguration.toObject(includeInstance, f),
    discoveryOverride: (f = msg.getDiscoveryOverride()) && proto.enterprise.gloo.solo.io.DiscoveryOverride.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.enterprise.gloo.solo.io.ExtAuthConfig.OidcAuthorizationCodeConfig}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OidcAuthorizationCodeConfig.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.enterprise.gloo.solo.io.ExtAuthConfig.OidcAuthorizationCodeConfig;
  return proto.enterprise.gloo.solo.io.ExtAuthConfig.OidcAuthorizationCodeConfig.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.enterprise.gloo.solo.io.ExtAuthConfig.OidcAuthorizationCodeConfig} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.enterprise.gloo.solo.io.ExtAuthConfig.OidcAuthorizationCodeConfig}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OidcAuthorizationCodeConfig.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setClientId(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setClientSecret(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setIssuerUrl(value);
      break;
    case 4:
      var value = msg.getAuthEndpointQueryParamsMap();
      reader.readMessage(value, function(message, reader) {
        jspb.Map.deserializeBinary(message, reader, jspb.BinaryReader.prototype.readString, jspb.BinaryReader.prototype.readString, null, "");
         });
      break;
    case 5:
      var value = /** @type {string} */ (reader.readString());
      msg.setAppUrl(value);
      break;
    case 6:
      var value = /** @type {string} */ (reader.readString());
      msg.setCallbackPath(value);
      break;
    case 9:
      var value = /** @type {string} */ (reader.readString());
      msg.setLogoutPath(value);
      break;
    case 7:
      var value = /** @type {string} */ (reader.readString());
      msg.addScopes(value);
      break;
    case 8:
      var value = new proto.enterprise.gloo.solo.io.UserSession;
      reader.readMessage(value,proto.enterprise.gloo.solo.io.UserSession.deserializeBinaryFromReader);
      msg.setSession(value);
      break;
    case 10:
      var value = new proto.enterprise.gloo.solo.io.HeaderConfiguration;
      reader.readMessage(value,proto.enterprise.gloo.solo.io.HeaderConfiguration.deserializeBinaryFromReader);
      msg.setHeaders(value);
      break;
    case 11:
      var value = new proto.enterprise.gloo.solo.io.DiscoveryOverride;
      reader.readMessage(value,proto.enterprise.gloo.solo.io.DiscoveryOverride.deserializeBinaryFromReader);
      msg.setDiscoveryOverride(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OidcAuthorizationCodeConfig.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.enterprise.gloo.solo.io.ExtAuthConfig.OidcAuthorizationCodeConfig.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.enterprise.gloo.solo.io.ExtAuthConfig.OidcAuthorizationCodeConfig} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OidcAuthorizationCodeConfig.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getClientId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getClientSecret();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getIssuerUrl();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getAuthEndpointQueryParamsMap(true);
  if (f && f.getLength() > 0) {
    f.serializeBinary(4, writer, jspb.BinaryWriter.prototype.writeString, jspb.BinaryWriter.prototype.writeString);
  }
  f = message.getAppUrl();
  if (f.length > 0) {
    writer.writeString(
      5,
      f
    );
  }
  f = message.getCallbackPath();
  if (f.length > 0) {
    writer.writeString(
      6,
      f
    );
  }
  f = message.getLogoutPath();
  if (f.length > 0) {
    writer.writeString(
      9,
      f
    );
  }
  f = message.getScopesList();
  if (f.length > 0) {
    writer.writeRepeatedString(
      7,
      f
    );
  }
  f = message.getSession();
  if (f != null) {
    writer.writeMessage(
      8,
      f,
      proto.enterprise.gloo.solo.io.UserSession.serializeBinaryToWriter
    );
  }
  f = message.getHeaders();
  if (f != null) {
    writer.writeMessage(
      10,
      f,
      proto.enterprise.gloo.solo.io.HeaderConfiguration.serializeBinaryToWriter
    );
  }
  f = message.getDiscoveryOverride();
  if (f != null) {
    writer.writeMessage(
      11,
      f,
      proto.enterprise.gloo.solo.io.DiscoveryOverride.serializeBinaryToWriter
    );
  }
};


/**
 * optional string client_id = 1;
 * @return {string}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OidcAuthorizationCodeConfig.prototype.getClientId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OidcAuthorizationCodeConfig.prototype.setClientId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string client_secret = 2;
 * @return {string}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OidcAuthorizationCodeConfig.prototype.getClientSecret = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OidcAuthorizationCodeConfig.prototype.setClientSecret = function(value) {
  jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string issuer_url = 3;
 * @return {string}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OidcAuthorizationCodeConfig.prototype.getIssuerUrl = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/** @param {string} value */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OidcAuthorizationCodeConfig.prototype.setIssuerUrl = function(value) {
  jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * map<string, string> auth_endpoint_query_params = 4;
 * @param {boolean=} opt_noLazyCreate Do not create the map if
 * empty, instead returning `undefined`
 * @return {!jspb.Map<string,string>}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OidcAuthorizationCodeConfig.prototype.getAuthEndpointQueryParamsMap = function(opt_noLazyCreate) {
  return /** @type {!jspb.Map<string,string>} */ (
      jspb.Message.getMapField(this, 4, opt_noLazyCreate,
      null));
};


proto.enterprise.gloo.solo.io.ExtAuthConfig.OidcAuthorizationCodeConfig.prototype.clearAuthEndpointQueryParamsMap = function() {
  this.getAuthEndpointQueryParamsMap().clear();
};


/**
 * optional string app_url = 5;
 * @return {string}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OidcAuthorizationCodeConfig.prototype.getAppUrl = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 5, ""));
};


/** @param {string} value */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OidcAuthorizationCodeConfig.prototype.setAppUrl = function(value) {
  jspb.Message.setProto3StringField(this, 5, value);
};


/**
 * optional string callback_path = 6;
 * @return {string}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OidcAuthorizationCodeConfig.prototype.getCallbackPath = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 6, ""));
};


/** @param {string} value */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OidcAuthorizationCodeConfig.prototype.setCallbackPath = function(value) {
  jspb.Message.setProto3StringField(this, 6, value);
};


/**
 * optional string logout_path = 9;
 * @return {string}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OidcAuthorizationCodeConfig.prototype.getLogoutPath = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 9, ""));
};


/** @param {string} value */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OidcAuthorizationCodeConfig.prototype.setLogoutPath = function(value) {
  jspb.Message.setProto3StringField(this, 9, value);
};


/**
 * repeated string scopes = 7;
 * @return {!Array<string>}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OidcAuthorizationCodeConfig.prototype.getScopesList = function() {
  return /** @type {!Array<string>} */ (jspb.Message.getRepeatedField(this, 7));
};


/** @param {!Array<string>} value */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OidcAuthorizationCodeConfig.prototype.setScopesList = function(value) {
  jspb.Message.setField(this, 7, value || []);
};


/**
 * @param {!string} value
 * @param {number=} opt_index
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OidcAuthorizationCodeConfig.prototype.addScopes = function(value, opt_index) {
  jspb.Message.addToRepeatedField(this, 7, value, opt_index);
};


proto.enterprise.gloo.solo.io.ExtAuthConfig.OidcAuthorizationCodeConfig.prototype.clearScopesList = function() {
  this.setScopesList([]);
};


/**
 * optional UserSession session = 8;
 * @return {?proto.enterprise.gloo.solo.io.UserSession}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OidcAuthorizationCodeConfig.prototype.getSession = function() {
  return /** @type{?proto.enterprise.gloo.solo.io.UserSession} */ (
    jspb.Message.getWrapperField(this, proto.enterprise.gloo.solo.io.UserSession, 8));
};


/** @param {?proto.enterprise.gloo.solo.io.UserSession|undefined} value */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OidcAuthorizationCodeConfig.prototype.setSession = function(value) {
  jspb.Message.setWrapperField(this, 8, value);
};


proto.enterprise.gloo.solo.io.ExtAuthConfig.OidcAuthorizationCodeConfig.prototype.clearSession = function() {
  this.setSession(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OidcAuthorizationCodeConfig.prototype.hasSession = function() {
  return jspb.Message.getField(this, 8) != null;
};


/**
 * optional HeaderConfiguration headers = 10;
 * @return {?proto.enterprise.gloo.solo.io.HeaderConfiguration}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OidcAuthorizationCodeConfig.prototype.getHeaders = function() {
  return /** @type{?proto.enterprise.gloo.solo.io.HeaderConfiguration} */ (
    jspb.Message.getWrapperField(this, proto.enterprise.gloo.solo.io.HeaderConfiguration, 10));
};


/** @param {?proto.enterprise.gloo.solo.io.HeaderConfiguration|undefined} value */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OidcAuthorizationCodeConfig.prototype.setHeaders = function(value) {
  jspb.Message.setWrapperField(this, 10, value);
};


proto.enterprise.gloo.solo.io.ExtAuthConfig.OidcAuthorizationCodeConfig.prototype.clearHeaders = function() {
  this.setHeaders(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OidcAuthorizationCodeConfig.prototype.hasHeaders = function() {
  return jspb.Message.getField(this, 10) != null;
};


/**
 * optional DiscoveryOverride discovery_override = 11;
 * @return {?proto.enterprise.gloo.solo.io.DiscoveryOverride}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OidcAuthorizationCodeConfig.prototype.getDiscoveryOverride = function() {
  return /** @type{?proto.enterprise.gloo.solo.io.DiscoveryOverride} */ (
    jspb.Message.getWrapperField(this, proto.enterprise.gloo.solo.io.DiscoveryOverride, 11));
};


/** @param {?proto.enterprise.gloo.solo.io.DiscoveryOverride|undefined} value */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OidcAuthorizationCodeConfig.prototype.setDiscoveryOverride = function(value) {
  jspb.Message.setWrapperField(this, 11, value);
};


proto.enterprise.gloo.solo.io.ExtAuthConfig.OidcAuthorizationCodeConfig.prototype.clearDiscoveryOverride = function() {
  this.setDiscoveryOverride(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OidcAuthorizationCodeConfig.prototype.hasDiscoveryOverride = function() {
  return jspb.Message.getField(this, 11) != null;
};



/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuth2Config = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuth2Config.oneofGroups_);
};
goog.inherits(proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuth2Config, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuth2Config.displayName = 'proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuth2Config';
}
/**
 * Oneof group definitions for this message. Each group defines the field
 * numbers belonging to that group. When of these fields' value is set, all
 * other fields in the group are cleared. During deserialization, if multiple
 * fields are encountered for a group, only the last value seen will be kept.
 * @private {!Array<!Array<number>>}
 * @const
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuth2Config.oneofGroups_ = [[1,2]];

/**
 * @enum {number}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuth2Config.OauthTypeCase = {
  OAUTH_TYPE_NOT_SET: 0,
  OIDC_AUTHORIZATION_CODE: 1,
  ACCESS_TOKEN_VALIDATION: 2
};

/**
 * @return {proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuth2Config.OauthTypeCase}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuth2Config.prototype.getOauthTypeCase = function() {
  return /** @type {proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuth2Config.OauthTypeCase} */(jspb.Message.computeOneofCase(this, proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuth2Config.oneofGroups_[0]));
};



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuth2Config.prototype.toObject = function(opt_includeInstance) {
  return proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuth2Config.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuth2Config} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuth2Config.toObject = function(includeInstance, msg) {
  var f, obj = {
    oidcAuthorizationCode: (f = msg.getOidcAuthorizationCode()) && proto.enterprise.gloo.solo.io.ExtAuthConfig.OidcAuthorizationCodeConfig.toObject(includeInstance, f),
    accessTokenValidation: (f = msg.getAccessTokenValidation()) && proto.enterprise.gloo.solo.io.AccessTokenValidation.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuth2Config}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuth2Config.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuth2Config;
  return proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuth2Config.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuth2Config} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuth2Config}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuth2Config.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.enterprise.gloo.solo.io.ExtAuthConfig.OidcAuthorizationCodeConfig;
      reader.readMessage(value,proto.enterprise.gloo.solo.io.ExtAuthConfig.OidcAuthorizationCodeConfig.deserializeBinaryFromReader);
      msg.setOidcAuthorizationCode(value);
      break;
    case 2:
      var value = new proto.enterprise.gloo.solo.io.AccessTokenValidation;
      reader.readMessage(value,proto.enterprise.gloo.solo.io.AccessTokenValidation.deserializeBinaryFromReader);
      msg.setAccessTokenValidation(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuth2Config.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuth2Config.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuth2Config} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuth2Config.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getOidcAuthorizationCode();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto.enterprise.gloo.solo.io.ExtAuthConfig.OidcAuthorizationCodeConfig.serializeBinaryToWriter
    );
  }
  f = message.getAccessTokenValidation();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.enterprise.gloo.solo.io.AccessTokenValidation.serializeBinaryToWriter
    );
  }
};


/**
 * optional OidcAuthorizationCodeConfig oidc_authorization_code = 1;
 * @return {?proto.enterprise.gloo.solo.io.ExtAuthConfig.OidcAuthorizationCodeConfig}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuth2Config.prototype.getOidcAuthorizationCode = function() {
  return /** @type{?proto.enterprise.gloo.solo.io.ExtAuthConfig.OidcAuthorizationCodeConfig} */ (
    jspb.Message.getWrapperField(this, proto.enterprise.gloo.solo.io.ExtAuthConfig.OidcAuthorizationCodeConfig, 1));
};


/** @param {?proto.enterprise.gloo.solo.io.ExtAuthConfig.OidcAuthorizationCodeConfig|undefined} value */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuth2Config.prototype.setOidcAuthorizationCode = function(value) {
  jspb.Message.setOneofWrapperField(this, 1, proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuth2Config.oneofGroups_[0], value);
};


proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuth2Config.prototype.clearOidcAuthorizationCode = function() {
  this.setOidcAuthorizationCode(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuth2Config.prototype.hasOidcAuthorizationCode = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional AccessTokenValidation access_token_validation = 2;
 * @return {?proto.enterprise.gloo.solo.io.AccessTokenValidation}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuth2Config.prototype.getAccessTokenValidation = function() {
  return /** @type{?proto.enterprise.gloo.solo.io.AccessTokenValidation} */ (
    jspb.Message.getWrapperField(this, proto.enterprise.gloo.solo.io.AccessTokenValidation, 2));
};


/** @param {?proto.enterprise.gloo.solo.io.AccessTokenValidation|undefined} value */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuth2Config.prototype.setAccessTokenValidation = function(value) {
  jspb.Message.setOneofWrapperField(this, 2, proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuth2Config.oneofGroups_[0], value);
};


proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuth2Config.prototype.clearAccessTokenValidation = function() {
  this.setAccessTokenValidation(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuth2Config.prototype.hasAccessTokenValidation = function() {
  return jspb.Message.getField(this, 2) != null;
};



/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.ApiKeyAuthConfig = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.enterprise.gloo.solo.io.ExtAuthConfig.ApiKeyAuthConfig, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.enterprise.gloo.solo.io.ExtAuthConfig.ApiKeyAuthConfig.displayName = 'proto.enterprise.gloo.solo.io.ExtAuthConfig.ApiKeyAuthConfig';
}


if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.ApiKeyAuthConfig.prototype.toObject = function(opt_includeInstance) {
  return proto.enterprise.gloo.solo.io.ExtAuthConfig.ApiKeyAuthConfig.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.enterprise.gloo.solo.io.ExtAuthConfig.ApiKeyAuthConfig} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.ApiKeyAuthConfig.toObject = function(includeInstance, msg) {
  var f, obj = {
    validApiKeysMap: (f = msg.getValidApiKeysMap()) ? f.toObject(includeInstance, proto.enterprise.gloo.solo.io.ExtAuthConfig.ApiKeyAuthConfig.KeyMetadata.toObject) : [],
    headerName: jspb.Message.getFieldWithDefault(msg, 2, ""),
    headersFromKeyMetadataMap: (f = msg.getHeadersFromKeyMetadataMap()) ? f.toObject(includeInstance, undefined) : []
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.enterprise.gloo.solo.io.ExtAuthConfig.ApiKeyAuthConfig}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.ApiKeyAuthConfig.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.enterprise.gloo.solo.io.ExtAuthConfig.ApiKeyAuthConfig;
  return proto.enterprise.gloo.solo.io.ExtAuthConfig.ApiKeyAuthConfig.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.enterprise.gloo.solo.io.ExtAuthConfig.ApiKeyAuthConfig} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.enterprise.gloo.solo.io.ExtAuthConfig.ApiKeyAuthConfig}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.ApiKeyAuthConfig.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = msg.getValidApiKeysMap();
      reader.readMessage(value, function(message, reader) {
        jspb.Map.deserializeBinary(message, reader, jspb.BinaryReader.prototype.readString, jspb.BinaryReader.prototype.readMessage, proto.enterprise.gloo.solo.io.ExtAuthConfig.ApiKeyAuthConfig.KeyMetadata.deserializeBinaryFromReader, "");
         });
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setHeaderName(value);
      break;
    case 3:
      var value = msg.getHeadersFromKeyMetadataMap();
      reader.readMessage(value, function(message, reader) {
        jspb.Map.deserializeBinary(message, reader, jspb.BinaryReader.prototype.readString, jspb.BinaryReader.prototype.readString, null, "");
         });
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.ApiKeyAuthConfig.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.enterprise.gloo.solo.io.ExtAuthConfig.ApiKeyAuthConfig.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.enterprise.gloo.solo.io.ExtAuthConfig.ApiKeyAuthConfig} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.ApiKeyAuthConfig.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getValidApiKeysMap(true);
  if (f && f.getLength() > 0) {
    f.serializeBinary(1, writer, jspb.BinaryWriter.prototype.writeString, jspb.BinaryWriter.prototype.writeMessage, proto.enterprise.gloo.solo.io.ExtAuthConfig.ApiKeyAuthConfig.KeyMetadata.serializeBinaryToWriter);
  }
  f = message.getHeaderName();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getHeadersFromKeyMetadataMap(true);
  if (f && f.getLength() > 0) {
    f.serializeBinary(3, writer, jspb.BinaryWriter.prototype.writeString, jspb.BinaryWriter.prototype.writeString);
  }
};



/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.ApiKeyAuthConfig.KeyMetadata = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.enterprise.gloo.solo.io.ExtAuthConfig.ApiKeyAuthConfig.KeyMetadata, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.enterprise.gloo.solo.io.ExtAuthConfig.ApiKeyAuthConfig.KeyMetadata.displayName = 'proto.enterprise.gloo.solo.io.ExtAuthConfig.ApiKeyAuthConfig.KeyMetadata';
}


if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.ApiKeyAuthConfig.KeyMetadata.prototype.toObject = function(opt_includeInstance) {
  return proto.enterprise.gloo.solo.io.ExtAuthConfig.ApiKeyAuthConfig.KeyMetadata.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.enterprise.gloo.solo.io.ExtAuthConfig.ApiKeyAuthConfig.KeyMetadata} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.ApiKeyAuthConfig.KeyMetadata.toObject = function(includeInstance, msg) {
  var f, obj = {
    username: jspb.Message.getFieldWithDefault(msg, 1, ""),
    metadataMap: (f = msg.getMetadataMap()) ? f.toObject(includeInstance, undefined) : []
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.enterprise.gloo.solo.io.ExtAuthConfig.ApiKeyAuthConfig.KeyMetadata}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.ApiKeyAuthConfig.KeyMetadata.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.enterprise.gloo.solo.io.ExtAuthConfig.ApiKeyAuthConfig.KeyMetadata;
  return proto.enterprise.gloo.solo.io.ExtAuthConfig.ApiKeyAuthConfig.KeyMetadata.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.enterprise.gloo.solo.io.ExtAuthConfig.ApiKeyAuthConfig.KeyMetadata} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.enterprise.gloo.solo.io.ExtAuthConfig.ApiKeyAuthConfig.KeyMetadata}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.ApiKeyAuthConfig.KeyMetadata.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setUsername(value);
      break;
    case 2:
      var value = msg.getMetadataMap();
      reader.readMessage(value, function(message, reader) {
        jspb.Map.deserializeBinary(message, reader, jspb.BinaryReader.prototype.readString, jspb.BinaryReader.prototype.readString, null, "");
         });
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.ApiKeyAuthConfig.KeyMetadata.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.enterprise.gloo.solo.io.ExtAuthConfig.ApiKeyAuthConfig.KeyMetadata.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.enterprise.gloo.solo.io.ExtAuthConfig.ApiKeyAuthConfig.KeyMetadata} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.ApiKeyAuthConfig.KeyMetadata.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getUsername();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getMetadataMap(true);
  if (f && f.getLength() > 0) {
    f.serializeBinary(2, writer, jspb.BinaryWriter.prototype.writeString, jspb.BinaryWriter.prototype.writeString);
  }
};


/**
 * optional string username = 1;
 * @return {string}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.ApiKeyAuthConfig.KeyMetadata.prototype.getUsername = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.enterprise.gloo.solo.io.ExtAuthConfig.ApiKeyAuthConfig.KeyMetadata.prototype.setUsername = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * map<string, string> metadata = 2;
 * @param {boolean=} opt_noLazyCreate Do not create the map if
 * empty, instead returning `undefined`
 * @return {!jspb.Map<string,string>}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.ApiKeyAuthConfig.KeyMetadata.prototype.getMetadataMap = function(opt_noLazyCreate) {
  return /** @type {!jspb.Map<string,string>} */ (
      jspb.Message.getMapField(this, 2, opt_noLazyCreate,
      null));
};


proto.enterprise.gloo.solo.io.ExtAuthConfig.ApiKeyAuthConfig.KeyMetadata.prototype.clearMetadataMap = function() {
  this.getMetadataMap().clear();
};


/**
 * map<string, KeyMetadata> valid_api_keys = 1;
 * @param {boolean=} opt_noLazyCreate Do not create the map if
 * empty, instead returning `undefined`
 * @return {!jspb.Map<string,!proto.enterprise.gloo.solo.io.ExtAuthConfig.ApiKeyAuthConfig.KeyMetadata>}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.ApiKeyAuthConfig.prototype.getValidApiKeysMap = function(opt_noLazyCreate) {
  return /** @type {!jspb.Map<string,!proto.enterprise.gloo.solo.io.ExtAuthConfig.ApiKeyAuthConfig.KeyMetadata>} */ (
      jspb.Message.getMapField(this, 1, opt_noLazyCreate,
      proto.enterprise.gloo.solo.io.ExtAuthConfig.ApiKeyAuthConfig.KeyMetadata));
};


proto.enterprise.gloo.solo.io.ExtAuthConfig.ApiKeyAuthConfig.prototype.clearValidApiKeysMap = function() {
  this.getValidApiKeysMap().clear();
};


/**
 * optional string header_name = 2;
 * @return {string}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.ApiKeyAuthConfig.prototype.getHeaderName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.enterprise.gloo.solo.io.ExtAuthConfig.ApiKeyAuthConfig.prototype.setHeaderName = function(value) {
  jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * map<string, string> headers_from_key_metadata = 3;
 * @param {boolean=} opt_noLazyCreate Do not create the map if
 * empty, instead returning `undefined`
 * @return {!jspb.Map<string,string>}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.ApiKeyAuthConfig.prototype.getHeadersFromKeyMetadataMap = function(opt_noLazyCreate) {
  return /** @type {!jspb.Map<string,string>} */ (
      jspb.Message.getMapField(this, 3, opt_noLazyCreate,
      null));
};


proto.enterprise.gloo.solo.io.ExtAuthConfig.ApiKeyAuthConfig.prototype.clearHeadersFromKeyMetadataMap = function() {
  this.getHeadersFromKeyMetadataMap().clear();
};



/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OpaAuthConfig = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.enterprise.gloo.solo.io.ExtAuthConfig.OpaAuthConfig, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.enterprise.gloo.solo.io.ExtAuthConfig.OpaAuthConfig.displayName = 'proto.enterprise.gloo.solo.io.ExtAuthConfig.OpaAuthConfig';
}


if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OpaAuthConfig.prototype.toObject = function(opt_includeInstance) {
  return proto.enterprise.gloo.solo.io.ExtAuthConfig.OpaAuthConfig.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.enterprise.gloo.solo.io.ExtAuthConfig.OpaAuthConfig} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OpaAuthConfig.toObject = function(includeInstance, msg) {
  var f, obj = {
    modulesMap: (f = msg.getModulesMap()) ? f.toObject(includeInstance, undefined) : [],
    query: jspb.Message.getFieldWithDefault(msg, 2, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.enterprise.gloo.solo.io.ExtAuthConfig.OpaAuthConfig}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OpaAuthConfig.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.enterprise.gloo.solo.io.ExtAuthConfig.OpaAuthConfig;
  return proto.enterprise.gloo.solo.io.ExtAuthConfig.OpaAuthConfig.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.enterprise.gloo.solo.io.ExtAuthConfig.OpaAuthConfig} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.enterprise.gloo.solo.io.ExtAuthConfig.OpaAuthConfig}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OpaAuthConfig.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = msg.getModulesMap();
      reader.readMessage(value, function(message, reader) {
        jspb.Map.deserializeBinary(message, reader, jspb.BinaryReader.prototype.readString, jspb.BinaryReader.prototype.readString, null, "");
         });
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setQuery(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OpaAuthConfig.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.enterprise.gloo.solo.io.ExtAuthConfig.OpaAuthConfig.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.enterprise.gloo.solo.io.ExtAuthConfig.OpaAuthConfig} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OpaAuthConfig.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getModulesMap(true);
  if (f && f.getLength() > 0) {
    f.serializeBinary(1, writer, jspb.BinaryWriter.prototype.writeString, jspb.BinaryWriter.prototype.writeString);
  }
  f = message.getQuery();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
};


/**
 * map<string, string> modules = 1;
 * @param {boolean=} opt_noLazyCreate Do not create the map if
 * empty, instead returning `undefined`
 * @return {!jspb.Map<string,string>}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OpaAuthConfig.prototype.getModulesMap = function(opt_noLazyCreate) {
  return /** @type {!jspb.Map<string,string>} */ (
      jspb.Message.getMapField(this, 1, opt_noLazyCreate,
      null));
};


proto.enterprise.gloo.solo.io.ExtAuthConfig.OpaAuthConfig.prototype.clearModulesMap = function() {
  this.getModulesMap().clear();
};


/**
 * optional string query = 2;
 * @return {string}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OpaAuthConfig.prototype.getQuery = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.enterprise.gloo.solo.io.ExtAuthConfig.OpaAuthConfig.prototype.setQuery = function(value) {
  jspb.Message.setProto3StringField(this, 2, value);
};



/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.Config = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.oneofGroups_);
};
goog.inherits(proto.enterprise.gloo.solo.io.ExtAuthConfig.Config, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.displayName = 'proto.enterprise.gloo.solo.io.ExtAuthConfig.Config';
}
/**
 * Oneof group definitions for this message. Each group defines the field
 * numbers belonging to that group. When of these fields' value is set, all
 * other fields in the group are cleared. During deserialization, if multiple
 * fields are encountered for a group, only the last value seen will be kept.
 * @private {!Array<!Array<number>>}
 * @const
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.oneofGroups_ = [[3,9,4,5,6,7,8,12,13]];

/**
 * @enum {number}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.AuthConfigCase = {
  AUTH_CONFIG_NOT_SET: 0,
  OAUTH: 3,
  OAUTH2: 9,
  BASIC_AUTH: 4,
  API_KEY_AUTH: 5,
  PLUGIN_AUTH: 6,
  OPA_AUTH: 7,
  LDAP: 8,
  JWT: 12,
  PASS_THROUGH_AUTH: 13
};

/**
 * @return {proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.AuthConfigCase}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.prototype.getAuthConfigCase = function() {
  return /** @type {proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.AuthConfigCase} */(jspb.Message.computeOneofCase(this, proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.oneofGroups_[0]));
};



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.prototype.toObject = function(opt_includeInstance) {
  return proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.enterprise.gloo.solo.io.ExtAuthConfig.Config} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.toObject = function(includeInstance, msg) {
  var f, obj = {
    name: (f = msg.getName()) && google_protobuf_wrappers_pb.StringValue.toObject(includeInstance, f),
    oauth: (f = msg.getOauth()) && proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuthConfig.toObject(includeInstance, f),
    oauth2: (f = msg.getOauth2()) && proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuth2Config.toObject(includeInstance, f),
    basicAuth: (f = msg.getBasicAuth()) && proto.enterprise.gloo.solo.io.BasicAuth.toObject(includeInstance, f),
    apiKeyAuth: (f = msg.getApiKeyAuth()) && proto.enterprise.gloo.solo.io.ExtAuthConfig.ApiKeyAuthConfig.toObject(includeInstance, f),
    pluginAuth: (f = msg.getPluginAuth()) && proto.enterprise.gloo.solo.io.AuthPlugin.toObject(includeInstance, f),
    opaAuth: (f = msg.getOpaAuth()) && proto.enterprise.gloo.solo.io.ExtAuthConfig.OpaAuthConfig.toObject(includeInstance, f),
    ldap: (f = msg.getLdap()) && proto.enterprise.gloo.solo.io.Ldap.toObject(includeInstance, f),
    jwt: (f = msg.getJwt()) && google_protobuf_empty_pb.Empty.toObject(includeInstance, f),
    passThroughAuth: (f = msg.getPassThroughAuth()) && proto.enterprise.gloo.solo.io.PassThroughAuth.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.enterprise.gloo.solo.io.ExtAuthConfig.Config}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.enterprise.gloo.solo.io.ExtAuthConfig.Config;
  return proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.enterprise.gloo.solo.io.ExtAuthConfig.Config} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.enterprise.gloo.solo.io.ExtAuthConfig.Config}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 11:
      var value = new google_protobuf_wrappers_pb.StringValue;
      reader.readMessage(value,google_protobuf_wrappers_pb.StringValue.deserializeBinaryFromReader);
      msg.setName(value);
      break;
    case 3:
      var value = new proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuthConfig;
      reader.readMessage(value,proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuthConfig.deserializeBinaryFromReader);
      msg.setOauth(value);
      break;
    case 9:
      var value = new proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuth2Config;
      reader.readMessage(value,proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuth2Config.deserializeBinaryFromReader);
      msg.setOauth2(value);
      break;
    case 4:
      var value = new proto.enterprise.gloo.solo.io.BasicAuth;
      reader.readMessage(value,proto.enterprise.gloo.solo.io.BasicAuth.deserializeBinaryFromReader);
      msg.setBasicAuth(value);
      break;
    case 5:
      var value = new proto.enterprise.gloo.solo.io.ExtAuthConfig.ApiKeyAuthConfig;
      reader.readMessage(value,proto.enterprise.gloo.solo.io.ExtAuthConfig.ApiKeyAuthConfig.deserializeBinaryFromReader);
      msg.setApiKeyAuth(value);
      break;
    case 6:
      var value = new proto.enterprise.gloo.solo.io.AuthPlugin;
      reader.readMessage(value,proto.enterprise.gloo.solo.io.AuthPlugin.deserializeBinaryFromReader);
      msg.setPluginAuth(value);
      break;
    case 7:
      var value = new proto.enterprise.gloo.solo.io.ExtAuthConfig.OpaAuthConfig;
      reader.readMessage(value,proto.enterprise.gloo.solo.io.ExtAuthConfig.OpaAuthConfig.deserializeBinaryFromReader);
      msg.setOpaAuth(value);
      break;
    case 8:
      var value = new proto.enterprise.gloo.solo.io.Ldap;
      reader.readMessage(value,proto.enterprise.gloo.solo.io.Ldap.deserializeBinaryFromReader);
      msg.setLdap(value);
      break;
    case 12:
      var value = new google_protobuf_empty_pb.Empty;
      reader.readMessage(value,google_protobuf_empty_pb.Empty.deserializeBinaryFromReader);
      msg.setJwt(value);
      break;
    case 13:
      var value = new proto.enterprise.gloo.solo.io.PassThroughAuth;
      reader.readMessage(value,proto.enterprise.gloo.solo.io.PassThroughAuth.deserializeBinaryFromReader);
      msg.setPassThroughAuth(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.enterprise.gloo.solo.io.ExtAuthConfig.Config} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getName();
  if (f != null) {
    writer.writeMessage(
      11,
      f,
      google_protobuf_wrappers_pb.StringValue.serializeBinaryToWriter
    );
  }
  f = message.getOauth();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuthConfig.serializeBinaryToWriter
    );
  }
  f = message.getOauth2();
  if (f != null) {
    writer.writeMessage(
      9,
      f,
      proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuth2Config.serializeBinaryToWriter
    );
  }
  f = message.getBasicAuth();
  if (f != null) {
    writer.writeMessage(
      4,
      f,
      proto.enterprise.gloo.solo.io.BasicAuth.serializeBinaryToWriter
    );
  }
  f = message.getApiKeyAuth();
  if (f != null) {
    writer.writeMessage(
      5,
      f,
      proto.enterprise.gloo.solo.io.ExtAuthConfig.ApiKeyAuthConfig.serializeBinaryToWriter
    );
  }
  f = message.getPluginAuth();
  if (f != null) {
    writer.writeMessage(
      6,
      f,
      proto.enterprise.gloo.solo.io.AuthPlugin.serializeBinaryToWriter
    );
  }
  f = message.getOpaAuth();
  if (f != null) {
    writer.writeMessage(
      7,
      f,
      proto.enterprise.gloo.solo.io.ExtAuthConfig.OpaAuthConfig.serializeBinaryToWriter
    );
  }
  f = message.getLdap();
  if (f != null) {
    writer.writeMessage(
      8,
      f,
      proto.enterprise.gloo.solo.io.Ldap.serializeBinaryToWriter
    );
  }
  f = message.getJwt();
  if (f != null) {
    writer.writeMessage(
      12,
      f,
      google_protobuf_empty_pb.Empty.serializeBinaryToWriter
    );
  }
  f = message.getPassThroughAuth();
  if (f != null) {
    writer.writeMessage(
      13,
      f,
      proto.enterprise.gloo.solo.io.PassThroughAuth.serializeBinaryToWriter
    );
  }
};


/**
 * optional google.protobuf.StringValue name = 11;
 * @return {?proto.google.protobuf.StringValue}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.prototype.getName = function() {
  return /** @type{?proto.google.protobuf.StringValue} */ (
    jspb.Message.getWrapperField(this, google_protobuf_wrappers_pb.StringValue, 11));
};


/** @param {?proto.google.protobuf.StringValue|undefined} value */
proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.prototype.setName = function(value) {
  jspb.Message.setWrapperField(this, 11, value);
};


proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.prototype.clearName = function() {
  this.setName(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.prototype.hasName = function() {
  return jspb.Message.getField(this, 11) != null;
};


/**
 * optional OAuthConfig oauth = 3;
 * @return {?proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuthConfig}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.prototype.getOauth = function() {
  return /** @type{?proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuthConfig} */ (
    jspb.Message.getWrapperField(this, proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuthConfig, 3));
};


/** @param {?proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuthConfig|undefined} value */
proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.prototype.setOauth = function(value) {
  jspb.Message.setOneofWrapperField(this, 3, proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.oneofGroups_[0], value);
};


proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.prototype.clearOauth = function() {
  this.setOauth(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.prototype.hasOauth = function() {
  return jspb.Message.getField(this, 3) != null;
};


/**
 * optional OAuth2Config oauth2 = 9;
 * @return {?proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuth2Config}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.prototype.getOauth2 = function() {
  return /** @type{?proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuth2Config} */ (
    jspb.Message.getWrapperField(this, proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuth2Config, 9));
};


/** @param {?proto.enterprise.gloo.solo.io.ExtAuthConfig.OAuth2Config|undefined} value */
proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.prototype.setOauth2 = function(value) {
  jspb.Message.setOneofWrapperField(this, 9, proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.oneofGroups_[0], value);
};


proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.prototype.clearOauth2 = function() {
  this.setOauth2(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.prototype.hasOauth2 = function() {
  return jspb.Message.getField(this, 9) != null;
};


/**
 * optional BasicAuth basic_auth = 4;
 * @return {?proto.enterprise.gloo.solo.io.BasicAuth}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.prototype.getBasicAuth = function() {
  return /** @type{?proto.enterprise.gloo.solo.io.BasicAuth} */ (
    jspb.Message.getWrapperField(this, proto.enterprise.gloo.solo.io.BasicAuth, 4));
};


/** @param {?proto.enterprise.gloo.solo.io.BasicAuth|undefined} value */
proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.prototype.setBasicAuth = function(value) {
  jspb.Message.setOneofWrapperField(this, 4, proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.oneofGroups_[0], value);
};


proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.prototype.clearBasicAuth = function() {
  this.setBasicAuth(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.prototype.hasBasicAuth = function() {
  return jspb.Message.getField(this, 4) != null;
};


/**
 * optional ApiKeyAuthConfig api_key_auth = 5;
 * @return {?proto.enterprise.gloo.solo.io.ExtAuthConfig.ApiKeyAuthConfig}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.prototype.getApiKeyAuth = function() {
  return /** @type{?proto.enterprise.gloo.solo.io.ExtAuthConfig.ApiKeyAuthConfig} */ (
    jspb.Message.getWrapperField(this, proto.enterprise.gloo.solo.io.ExtAuthConfig.ApiKeyAuthConfig, 5));
};


/** @param {?proto.enterprise.gloo.solo.io.ExtAuthConfig.ApiKeyAuthConfig|undefined} value */
proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.prototype.setApiKeyAuth = function(value) {
  jspb.Message.setOneofWrapperField(this, 5, proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.oneofGroups_[0], value);
};


proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.prototype.clearApiKeyAuth = function() {
  this.setApiKeyAuth(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.prototype.hasApiKeyAuth = function() {
  return jspb.Message.getField(this, 5) != null;
};


/**
 * optional AuthPlugin plugin_auth = 6;
 * @return {?proto.enterprise.gloo.solo.io.AuthPlugin}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.prototype.getPluginAuth = function() {
  return /** @type{?proto.enterprise.gloo.solo.io.AuthPlugin} */ (
    jspb.Message.getWrapperField(this, proto.enterprise.gloo.solo.io.AuthPlugin, 6));
};


/** @param {?proto.enterprise.gloo.solo.io.AuthPlugin|undefined} value */
proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.prototype.setPluginAuth = function(value) {
  jspb.Message.setOneofWrapperField(this, 6, proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.oneofGroups_[0], value);
};


proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.prototype.clearPluginAuth = function() {
  this.setPluginAuth(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.prototype.hasPluginAuth = function() {
  return jspb.Message.getField(this, 6) != null;
};


/**
 * optional OpaAuthConfig opa_auth = 7;
 * @return {?proto.enterprise.gloo.solo.io.ExtAuthConfig.OpaAuthConfig}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.prototype.getOpaAuth = function() {
  return /** @type{?proto.enterprise.gloo.solo.io.ExtAuthConfig.OpaAuthConfig} */ (
    jspb.Message.getWrapperField(this, proto.enterprise.gloo.solo.io.ExtAuthConfig.OpaAuthConfig, 7));
};


/** @param {?proto.enterprise.gloo.solo.io.ExtAuthConfig.OpaAuthConfig|undefined} value */
proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.prototype.setOpaAuth = function(value) {
  jspb.Message.setOneofWrapperField(this, 7, proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.oneofGroups_[0], value);
};


proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.prototype.clearOpaAuth = function() {
  this.setOpaAuth(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.prototype.hasOpaAuth = function() {
  return jspb.Message.getField(this, 7) != null;
};


/**
 * optional Ldap ldap = 8;
 * @return {?proto.enterprise.gloo.solo.io.Ldap}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.prototype.getLdap = function() {
  return /** @type{?proto.enterprise.gloo.solo.io.Ldap} */ (
    jspb.Message.getWrapperField(this, proto.enterprise.gloo.solo.io.Ldap, 8));
};


/** @param {?proto.enterprise.gloo.solo.io.Ldap|undefined} value */
proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.prototype.setLdap = function(value) {
  jspb.Message.setOneofWrapperField(this, 8, proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.oneofGroups_[0], value);
};


proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.prototype.clearLdap = function() {
  this.setLdap(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.prototype.hasLdap = function() {
  return jspb.Message.getField(this, 8) != null;
};


/**
 * optional google.protobuf.Empty jwt = 12;
 * @return {?proto.google.protobuf.Empty}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.prototype.getJwt = function() {
  return /** @type{?proto.google.protobuf.Empty} */ (
    jspb.Message.getWrapperField(this, google_protobuf_empty_pb.Empty, 12));
};


/** @param {?proto.google.protobuf.Empty|undefined} value */
proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.prototype.setJwt = function(value) {
  jspb.Message.setOneofWrapperField(this, 12, proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.oneofGroups_[0], value);
};


proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.prototype.clearJwt = function() {
  this.setJwt(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.prototype.hasJwt = function() {
  return jspb.Message.getField(this, 12) != null;
};


/**
 * optional PassThroughAuth pass_through_auth = 13;
 * @return {?proto.enterprise.gloo.solo.io.PassThroughAuth}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.prototype.getPassThroughAuth = function() {
  return /** @type{?proto.enterprise.gloo.solo.io.PassThroughAuth} */ (
    jspb.Message.getWrapperField(this, proto.enterprise.gloo.solo.io.PassThroughAuth, 13));
};


/** @param {?proto.enterprise.gloo.solo.io.PassThroughAuth|undefined} value */
proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.prototype.setPassThroughAuth = function(value) {
  jspb.Message.setOneofWrapperField(this, 13, proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.oneofGroups_[0], value);
};


proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.prototype.clearPassThroughAuth = function() {
  this.setPassThroughAuth(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.Config.prototype.hasPassThroughAuth = function() {
  return jspb.Message.getField(this, 13) != null;
};


/**
 * optional string auth_config_ref_name = 1;
 * @return {string}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.prototype.getAuthConfigRefName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.enterprise.gloo.solo.io.ExtAuthConfig.prototype.setAuthConfigRefName = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * repeated Config configs = 8;
 * @return {!Array<!proto.enterprise.gloo.solo.io.ExtAuthConfig.Config>}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.prototype.getConfigsList = function() {
  return /** @type{!Array<!proto.enterprise.gloo.solo.io.ExtAuthConfig.Config>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.enterprise.gloo.solo.io.ExtAuthConfig.Config, 8));
};


/** @param {!Array<!proto.enterprise.gloo.solo.io.ExtAuthConfig.Config>} value */
proto.enterprise.gloo.solo.io.ExtAuthConfig.prototype.setConfigsList = function(value) {
  jspb.Message.setRepeatedWrapperField(this, 8, value);
};


/**
 * @param {!proto.enterprise.gloo.solo.io.ExtAuthConfig.Config=} opt_value
 * @param {number=} opt_index
 * @return {!proto.enterprise.gloo.solo.io.ExtAuthConfig.Config}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.prototype.addConfigs = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 8, opt_value, proto.enterprise.gloo.solo.io.ExtAuthConfig.Config, opt_index);
};


proto.enterprise.gloo.solo.io.ExtAuthConfig.prototype.clearConfigsList = function() {
  this.setConfigsList([]);
};


/**
 * optional google.protobuf.StringValue boolean_expr = 10;
 * @return {?proto.google.protobuf.StringValue}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.prototype.getBooleanExpr = function() {
  return /** @type{?proto.google.protobuf.StringValue} */ (
    jspb.Message.getWrapperField(this, google_protobuf_wrappers_pb.StringValue, 10));
};


/** @param {?proto.google.protobuf.StringValue|undefined} value */
proto.enterprise.gloo.solo.io.ExtAuthConfig.prototype.setBooleanExpr = function(value) {
  jspb.Message.setWrapperField(this, 10, value);
};


proto.enterprise.gloo.solo.io.ExtAuthConfig.prototype.clearBooleanExpr = function() {
  this.setBooleanExpr(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.enterprise.gloo.solo.io.ExtAuthConfig.prototype.hasBooleanExpr = function() {
  return jspb.Message.getField(this, 10) != null;
};


goog.object.extend(exports, proto.enterprise.gloo.solo.io);
