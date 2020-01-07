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

var google_protobuf_struct_pb = require('google-protobuf/google/protobuf/struct_pb.js');
var gogoproto_gogo_pb = require('../../../../../gogoproto/gogo_pb.js');
var extproto_ext_pb = require('../../../../../protoc-gen-ext/extproto/ext_pb.js');
var gloo_projects_gloo_api_v1_extensions_pb = require('../../../../../gloo/projects/gloo/api/v1/extensions_pb.js');
var gloo_projects_gloo_api_v1_options_cors_cors_pb = require('../../../../../gloo/projects/gloo/api/v1/options/cors/cors_pb.js');
var gloo_projects_gloo_api_v1_options_rest_rest_pb = require('../../../../../gloo/projects/gloo/api/v1/options/rest/rest_pb.js');
var gloo_projects_gloo_api_v1_options_grpc_grpc_pb = require('../../../../../gloo/projects/gloo/api/v1/options/grpc/grpc_pb.js');
var gloo_projects_gloo_api_v1_options_als_als_pb = require('../../../../../gloo/projects/gloo/api/v1/options/als/als_pb.js');
var gloo_projects_gloo_api_v1_options_grpc_web_grpc_web_pb = require('../../../../../gloo/projects/gloo/api/v1/options/grpc_web/grpc_web_pb.js');
var gloo_projects_gloo_api_v1_options_hcm_hcm_pb = require('../../../../../gloo/projects/gloo/api/v1/options/hcm/hcm_pb.js');
var gloo_projects_gloo_api_v1_options_lbhash_lbhash_pb = require('../../../../../gloo/projects/gloo/api/v1/options/lbhash/lbhash_pb.js');
var gloo_projects_gloo_api_v1_options_shadowing_shadowing_pb = require('../../../../../gloo/projects/gloo/api/v1/options/shadowing/shadowing_pb.js');
var gloo_projects_gloo_api_v1_options_tcp_tcp_pb = require('../../../../../gloo/projects/gloo/api/v1/options/tcp/tcp_pb.js');
var gloo_projects_gloo_api_v1_options_tracing_tracing_pb = require('../../../../../gloo/projects/gloo/api/v1/options/tracing/tracing_pb.js');
var gloo_projects_gloo_api_v1_options_retries_retries_pb = require('../../../../../gloo/projects/gloo/api/v1/options/retries/retries_pb.js');
var gloo_projects_gloo_api_v1_options_stats_stats_pb = require('../../../../../gloo/projects/gloo/api/v1/options/stats/stats_pb.js');
var gloo_projects_gloo_api_v1_options_faultinjection_fault_pb = require('../../../../../gloo/projects/gloo/api/v1/options/faultinjection/fault_pb.js');
var gloo_projects_gloo_api_v1_options_headers_headers_pb = require('../../../../../gloo/projects/gloo/api/v1/options/headers/headers_pb.js');
var gloo_projects_gloo_api_v1_options_aws_aws_pb = require('../../../../../gloo/projects/gloo/api/v1/options/aws/aws_pb.js');
var gloo_projects_gloo_api_v1_options_wasm_wasm_pb = require('../../../../../gloo/projects/gloo/api/v1/options/wasm/wasm_pb.js');
var gloo_projects_gloo_api_v1_options_azure_azure_pb = require('../../../../../gloo/projects/gloo/api/v1/options/azure/azure_pb.js');
var gloo_projects_gloo_api_v1_options_healthcheck_healthcheck_pb = require('../../../../../gloo/projects/gloo/api/v1/options/healthcheck/healthcheck_pb.js');
var gloo_projects_gloo_api_v1_options_protocol_upgrade_protocol_upgrade_pb = require('../../../../../gloo/projects/gloo/api/v1/options/protocol_upgrade/protocol_upgrade_pb.js');
var gloo_projects_gloo_api_external_envoy_extensions_transformation_transformation_pb = require('../../../../../gloo/projects/gloo/api/external/envoy/extensions/transformation/transformation_pb.js');
var gloo_projects_gloo_api_v1_enterprise_options_extauth_v1_extauth_pb = require('../../../../../gloo/projects/gloo/api/v1/enterprise/options/extauth/v1/extauth_pb.js');
var gloo_projects_gloo_api_v1_enterprise_options_jwt_jwt_pb = require('../../../../../gloo/projects/gloo/api/v1/enterprise/options/jwt/jwt_pb.js');
var gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_pb = require('../../../../../gloo/projects/gloo/api/v1/enterprise/options/ratelimit/ratelimit_pb.js');
var gloo_projects_gloo_api_v1_enterprise_options_rbac_rbac_pb = require('../../../../../gloo/projects/gloo/api/v1/enterprise/options/rbac/rbac_pb.js');
var gloo_projects_gloo_api_v1_enterprise_options_waf_waf_pb = require('../../../../../gloo/projects/gloo/api/v1/enterprise/options/waf/waf_pb.js');
var gloo_projects_gloo_api_v1_enterprise_options_dlp_dlp_pb = require('../../../../../gloo/projects/gloo/api/v1/enterprise/options/dlp/dlp_pb.js');
var google_protobuf_duration_pb = require('google-protobuf/google/protobuf/duration_pb.js');
var google_protobuf_wrappers_pb = require('google-protobuf/google/protobuf/wrappers_pb.js');
goog.exportSymbol('proto.gloo.solo.io.DestinationSpec', null, global);
goog.exportSymbol('proto.gloo.solo.io.HttpListenerOptions', null, global);
goog.exportSymbol('proto.gloo.solo.io.ListenerOptions', null, global);
goog.exportSymbol('proto.gloo.solo.io.RouteOptions', null, global);
goog.exportSymbol('proto.gloo.solo.io.TcpListenerOptions', null, global);
goog.exportSymbol('proto.gloo.solo.io.VirtualHostOptions', null, global);
goog.exportSymbol('proto.gloo.solo.io.WeightedDestinationOptions', null, global);

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
proto.gloo.solo.io.ListenerOptions = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.gloo.solo.io.ListenerOptions, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.gloo.solo.io.ListenerOptions.displayName = 'proto.gloo.solo.io.ListenerOptions';
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
proto.gloo.solo.io.ListenerOptions.prototype.toObject = function(opt_includeInstance) {
  return proto.gloo.solo.io.ListenerOptions.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.gloo.solo.io.ListenerOptions} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.gloo.solo.io.ListenerOptions.toObject = function(includeInstance, msg) {
  var f, obj = {
    accessLoggingService: (f = msg.getAccessLoggingService()) && gloo_projects_gloo_api_v1_options_als_als_pb.AccessLoggingService.toObject(includeInstance, f),
    extensions: (f = msg.getExtensions()) && gloo_projects_gloo_api_v1_extensions_pb.Extensions.toObject(includeInstance, f)
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
 * @return {!proto.gloo.solo.io.ListenerOptions}
 */
proto.gloo.solo.io.ListenerOptions.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.gloo.solo.io.ListenerOptions;
  return proto.gloo.solo.io.ListenerOptions.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.gloo.solo.io.ListenerOptions} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.gloo.solo.io.ListenerOptions}
 */
proto.gloo.solo.io.ListenerOptions.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new gloo_projects_gloo_api_v1_options_als_als_pb.AccessLoggingService;
      reader.readMessage(value,gloo_projects_gloo_api_v1_options_als_als_pb.AccessLoggingService.deserializeBinaryFromReader);
      msg.setAccessLoggingService(value);
      break;
    case 2:
      var value = new gloo_projects_gloo_api_v1_extensions_pb.Extensions;
      reader.readMessage(value,gloo_projects_gloo_api_v1_extensions_pb.Extensions.deserializeBinaryFromReader);
      msg.setExtensions(value);
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
proto.gloo.solo.io.ListenerOptions.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.gloo.solo.io.ListenerOptions.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.gloo.solo.io.ListenerOptions} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.gloo.solo.io.ListenerOptions.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getAccessLoggingService();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      gloo_projects_gloo_api_v1_options_als_als_pb.AccessLoggingService.serializeBinaryToWriter
    );
  }
  f = message.getExtensions();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      gloo_projects_gloo_api_v1_extensions_pb.Extensions.serializeBinaryToWriter
    );
  }
};


/**
 * optional als.options.gloo.solo.io.AccessLoggingService access_logging_service = 1;
 * @return {?proto.als.options.gloo.solo.io.AccessLoggingService}
 */
proto.gloo.solo.io.ListenerOptions.prototype.getAccessLoggingService = function() {
  return /** @type{?proto.als.options.gloo.solo.io.AccessLoggingService} */ (
    jspb.Message.getWrapperField(this, gloo_projects_gloo_api_v1_options_als_als_pb.AccessLoggingService, 1));
};


/** @param {?proto.als.options.gloo.solo.io.AccessLoggingService|undefined} value */
proto.gloo.solo.io.ListenerOptions.prototype.setAccessLoggingService = function(value) {
  jspb.Message.setWrapperField(this, 1, value);
};


proto.gloo.solo.io.ListenerOptions.prototype.clearAccessLoggingService = function() {
  this.setAccessLoggingService(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.gloo.solo.io.ListenerOptions.prototype.hasAccessLoggingService = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional Extensions extensions = 2;
 * @return {?proto.gloo.solo.io.Extensions}
 */
proto.gloo.solo.io.ListenerOptions.prototype.getExtensions = function() {
  return /** @type{?proto.gloo.solo.io.Extensions} */ (
    jspb.Message.getWrapperField(this, gloo_projects_gloo_api_v1_extensions_pb.Extensions, 2));
};


/** @param {?proto.gloo.solo.io.Extensions|undefined} value */
proto.gloo.solo.io.ListenerOptions.prototype.setExtensions = function(value) {
  jspb.Message.setWrapperField(this, 2, value);
};


proto.gloo.solo.io.ListenerOptions.prototype.clearExtensions = function() {
  this.setExtensions(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.gloo.solo.io.ListenerOptions.prototype.hasExtensions = function() {
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
proto.gloo.solo.io.HttpListenerOptions = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.gloo.solo.io.HttpListenerOptions, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.gloo.solo.io.HttpListenerOptions.displayName = 'proto.gloo.solo.io.HttpListenerOptions';
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
proto.gloo.solo.io.HttpListenerOptions.prototype.toObject = function(opt_includeInstance) {
  return proto.gloo.solo.io.HttpListenerOptions.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.gloo.solo.io.HttpListenerOptions} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.gloo.solo.io.HttpListenerOptions.toObject = function(includeInstance, msg) {
  var f, obj = {
    grpcWeb: (f = msg.getGrpcWeb()) && gloo_projects_gloo_api_v1_options_grpc_web_grpc_web_pb.GrpcWeb.toObject(includeInstance, f),
    httpConnectionManagerSettings: (f = msg.getHttpConnectionManagerSettings()) && gloo_projects_gloo_api_v1_options_hcm_hcm_pb.HttpConnectionManagerSettings.toObject(includeInstance, f),
    healthCheck: (f = msg.getHealthCheck()) && gloo_projects_gloo_api_v1_options_healthcheck_healthcheck_pb.HealthCheck.toObject(includeInstance, f),
    extensions: (f = msg.getExtensions()) && gloo_projects_gloo_api_v1_extensions_pb.Extensions.toObject(includeInstance, f),
    waf: (f = msg.getWaf()) && gloo_projects_gloo_api_v1_enterprise_options_waf_waf_pb.Settings.toObject(includeInstance, f),
    dlp: (f = msg.getDlp()) && gloo_projects_gloo_api_v1_enterprise_options_dlp_dlp_pb.FilterConfig.toObject(includeInstance, f),
    wasm: (f = msg.getWasm()) && gloo_projects_gloo_api_v1_options_wasm_wasm_pb.PluginSource.toObject(includeInstance, f)
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
 * @return {!proto.gloo.solo.io.HttpListenerOptions}
 */
proto.gloo.solo.io.HttpListenerOptions.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.gloo.solo.io.HttpListenerOptions;
  return proto.gloo.solo.io.HttpListenerOptions.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.gloo.solo.io.HttpListenerOptions} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.gloo.solo.io.HttpListenerOptions}
 */
proto.gloo.solo.io.HttpListenerOptions.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new gloo_projects_gloo_api_v1_options_grpc_web_grpc_web_pb.GrpcWeb;
      reader.readMessage(value,gloo_projects_gloo_api_v1_options_grpc_web_grpc_web_pb.GrpcWeb.deserializeBinaryFromReader);
      msg.setGrpcWeb(value);
      break;
    case 2:
      var value = new gloo_projects_gloo_api_v1_options_hcm_hcm_pb.HttpConnectionManagerSettings;
      reader.readMessage(value,gloo_projects_gloo_api_v1_options_hcm_hcm_pb.HttpConnectionManagerSettings.deserializeBinaryFromReader);
      msg.setHttpConnectionManagerSettings(value);
      break;
    case 4:
      var value = new gloo_projects_gloo_api_v1_options_healthcheck_healthcheck_pb.HealthCheck;
      reader.readMessage(value,gloo_projects_gloo_api_v1_options_healthcheck_healthcheck_pb.HealthCheck.deserializeBinaryFromReader);
      msg.setHealthCheck(value);
      break;
    case 3:
      var value = new gloo_projects_gloo_api_v1_extensions_pb.Extensions;
      reader.readMessage(value,gloo_projects_gloo_api_v1_extensions_pb.Extensions.deserializeBinaryFromReader);
      msg.setExtensions(value);
      break;
    case 5:
      var value = new gloo_projects_gloo_api_v1_enterprise_options_waf_waf_pb.Settings;
      reader.readMessage(value,gloo_projects_gloo_api_v1_enterprise_options_waf_waf_pb.Settings.deserializeBinaryFromReader);
      msg.setWaf(value);
      break;
    case 6:
      var value = new gloo_projects_gloo_api_v1_enterprise_options_dlp_dlp_pb.FilterConfig;
      reader.readMessage(value,gloo_projects_gloo_api_v1_enterprise_options_dlp_dlp_pb.FilterConfig.deserializeBinaryFromReader);
      msg.setDlp(value);
      break;
    case 7:
      var value = new gloo_projects_gloo_api_v1_options_wasm_wasm_pb.PluginSource;
      reader.readMessage(value,gloo_projects_gloo_api_v1_options_wasm_wasm_pb.PluginSource.deserializeBinaryFromReader);
      msg.setWasm(value);
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
proto.gloo.solo.io.HttpListenerOptions.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.gloo.solo.io.HttpListenerOptions.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.gloo.solo.io.HttpListenerOptions} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.gloo.solo.io.HttpListenerOptions.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getGrpcWeb();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      gloo_projects_gloo_api_v1_options_grpc_web_grpc_web_pb.GrpcWeb.serializeBinaryToWriter
    );
  }
  f = message.getHttpConnectionManagerSettings();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      gloo_projects_gloo_api_v1_options_hcm_hcm_pb.HttpConnectionManagerSettings.serializeBinaryToWriter
    );
  }
  f = message.getHealthCheck();
  if (f != null) {
    writer.writeMessage(
      4,
      f,
      gloo_projects_gloo_api_v1_options_healthcheck_healthcheck_pb.HealthCheck.serializeBinaryToWriter
    );
  }
  f = message.getExtensions();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      gloo_projects_gloo_api_v1_extensions_pb.Extensions.serializeBinaryToWriter
    );
  }
  f = message.getWaf();
  if (f != null) {
    writer.writeMessage(
      5,
      f,
      gloo_projects_gloo_api_v1_enterprise_options_waf_waf_pb.Settings.serializeBinaryToWriter
    );
  }
  f = message.getDlp();
  if (f != null) {
    writer.writeMessage(
      6,
      f,
      gloo_projects_gloo_api_v1_enterprise_options_dlp_dlp_pb.FilterConfig.serializeBinaryToWriter
    );
  }
  f = message.getWasm();
  if (f != null) {
    writer.writeMessage(
      7,
      f,
      gloo_projects_gloo_api_v1_options_wasm_wasm_pb.PluginSource.serializeBinaryToWriter
    );
  }
};


/**
 * optional grpc_web.options.gloo.solo.io.GrpcWeb grpc_web = 1;
 * @return {?proto.grpc_web.options.gloo.solo.io.GrpcWeb}
 */
proto.gloo.solo.io.HttpListenerOptions.prototype.getGrpcWeb = function() {
  return /** @type{?proto.grpc_web.options.gloo.solo.io.GrpcWeb} */ (
    jspb.Message.getWrapperField(this, gloo_projects_gloo_api_v1_options_grpc_web_grpc_web_pb.GrpcWeb, 1));
};


/** @param {?proto.grpc_web.options.gloo.solo.io.GrpcWeb|undefined} value */
proto.gloo.solo.io.HttpListenerOptions.prototype.setGrpcWeb = function(value) {
  jspb.Message.setWrapperField(this, 1, value);
};


proto.gloo.solo.io.HttpListenerOptions.prototype.clearGrpcWeb = function() {
  this.setGrpcWeb(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.gloo.solo.io.HttpListenerOptions.prototype.hasGrpcWeb = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional hcm.options.gloo.solo.io.HttpConnectionManagerSettings http_connection_manager_settings = 2;
 * @return {?proto.hcm.options.gloo.solo.io.HttpConnectionManagerSettings}
 */
proto.gloo.solo.io.HttpListenerOptions.prototype.getHttpConnectionManagerSettings = function() {
  return /** @type{?proto.hcm.options.gloo.solo.io.HttpConnectionManagerSettings} */ (
    jspb.Message.getWrapperField(this, gloo_projects_gloo_api_v1_options_hcm_hcm_pb.HttpConnectionManagerSettings, 2));
};


/** @param {?proto.hcm.options.gloo.solo.io.HttpConnectionManagerSettings|undefined} value */
proto.gloo.solo.io.HttpListenerOptions.prototype.setHttpConnectionManagerSettings = function(value) {
  jspb.Message.setWrapperField(this, 2, value);
};


proto.gloo.solo.io.HttpListenerOptions.prototype.clearHttpConnectionManagerSettings = function() {
  this.setHttpConnectionManagerSettings(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.gloo.solo.io.HttpListenerOptions.prototype.hasHttpConnectionManagerSettings = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional healthcheck.options.gloo.solo.io.HealthCheck health_check = 4;
 * @return {?proto.healthcheck.options.gloo.solo.io.HealthCheck}
 */
proto.gloo.solo.io.HttpListenerOptions.prototype.getHealthCheck = function() {
  return /** @type{?proto.healthcheck.options.gloo.solo.io.HealthCheck} */ (
    jspb.Message.getWrapperField(this, gloo_projects_gloo_api_v1_options_healthcheck_healthcheck_pb.HealthCheck, 4));
};


/** @param {?proto.healthcheck.options.gloo.solo.io.HealthCheck|undefined} value */
proto.gloo.solo.io.HttpListenerOptions.prototype.setHealthCheck = function(value) {
  jspb.Message.setWrapperField(this, 4, value);
};


proto.gloo.solo.io.HttpListenerOptions.prototype.clearHealthCheck = function() {
  this.setHealthCheck(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.gloo.solo.io.HttpListenerOptions.prototype.hasHealthCheck = function() {
  return jspb.Message.getField(this, 4) != null;
};


/**
 * optional Extensions extensions = 3;
 * @return {?proto.gloo.solo.io.Extensions}
 */
proto.gloo.solo.io.HttpListenerOptions.prototype.getExtensions = function() {
  return /** @type{?proto.gloo.solo.io.Extensions} */ (
    jspb.Message.getWrapperField(this, gloo_projects_gloo_api_v1_extensions_pb.Extensions, 3));
};


/** @param {?proto.gloo.solo.io.Extensions|undefined} value */
proto.gloo.solo.io.HttpListenerOptions.prototype.setExtensions = function(value) {
  jspb.Message.setWrapperField(this, 3, value);
};


proto.gloo.solo.io.HttpListenerOptions.prototype.clearExtensions = function() {
  this.setExtensions(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.gloo.solo.io.HttpListenerOptions.prototype.hasExtensions = function() {
  return jspb.Message.getField(this, 3) != null;
};


/**
 * optional waf.options.gloo.solo.io.Settings waf = 5;
 * @return {?proto.waf.options.gloo.solo.io.Settings}
 */
proto.gloo.solo.io.HttpListenerOptions.prototype.getWaf = function() {
  return /** @type{?proto.waf.options.gloo.solo.io.Settings} */ (
    jspb.Message.getWrapperField(this, gloo_projects_gloo_api_v1_enterprise_options_waf_waf_pb.Settings, 5));
};


/** @param {?proto.waf.options.gloo.solo.io.Settings|undefined} value */
proto.gloo.solo.io.HttpListenerOptions.prototype.setWaf = function(value) {
  jspb.Message.setWrapperField(this, 5, value);
};


proto.gloo.solo.io.HttpListenerOptions.prototype.clearWaf = function() {
  this.setWaf(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.gloo.solo.io.HttpListenerOptions.prototype.hasWaf = function() {
  return jspb.Message.getField(this, 5) != null;
};


/**
 * optional dlp.options.gloo.solo.io.FilterConfig dlp = 6;
 * @return {?proto.dlp.options.gloo.solo.io.FilterConfig}
 */
proto.gloo.solo.io.HttpListenerOptions.prototype.getDlp = function() {
  return /** @type{?proto.dlp.options.gloo.solo.io.FilterConfig} */ (
    jspb.Message.getWrapperField(this, gloo_projects_gloo_api_v1_enterprise_options_dlp_dlp_pb.FilterConfig, 6));
};


/** @param {?proto.dlp.options.gloo.solo.io.FilterConfig|undefined} value */
proto.gloo.solo.io.HttpListenerOptions.prototype.setDlp = function(value) {
  jspb.Message.setWrapperField(this, 6, value);
};


proto.gloo.solo.io.HttpListenerOptions.prototype.clearDlp = function() {
  this.setDlp(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.gloo.solo.io.HttpListenerOptions.prototype.hasDlp = function() {
  return jspb.Message.getField(this, 6) != null;
};


/**
 * optional wasm.options.gloo.solo.io.PluginSource wasm = 7;
 * @return {?proto.wasm.options.gloo.solo.io.PluginSource}
 */
proto.gloo.solo.io.HttpListenerOptions.prototype.getWasm = function() {
  return /** @type{?proto.wasm.options.gloo.solo.io.PluginSource} */ (
    jspb.Message.getWrapperField(this, gloo_projects_gloo_api_v1_options_wasm_wasm_pb.PluginSource, 7));
};


/** @param {?proto.wasm.options.gloo.solo.io.PluginSource|undefined} value */
proto.gloo.solo.io.HttpListenerOptions.prototype.setWasm = function(value) {
  jspb.Message.setWrapperField(this, 7, value);
};


proto.gloo.solo.io.HttpListenerOptions.prototype.clearWasm = function() {
  this.setWasm(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.gloo.solo.io.HttpListenerOptions.prototype.hasWasm = function() {
  return jspb.Message.getField(this, 7) != null;
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
proto.gloo.solo.io.TcpListenerOptions = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.gloo.solo.io.TcpListenerOptions, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.gloo.solo.io.TcpListenerOptions.displayName = 'proto.gloo.solo.io.TcpListenerOptions';
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
proto.gloo.solo.io.TcpListenerOptions.prototype.toObject = function(opt_includeInstance) {
  return proto.gloo.solo.io.TcpListenerOptions.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.gloo.solo.io.TcpListenerOptions} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.gloo.solo.io.TcpListenerOptions.toObject = function(includeInstance, msg) {
  var f, obj = {
    tcpProxySettings: (f = msg.getTcpProxySettings()) && gloo_projects_gloo_api_v1_options_tcp_tcp_pb.TcpProxySettings.toObject(includeInstance, f)
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
 * @return {!proto.gloo.solo.io.TcpListenerOptions}
 */
proto.gloo.solo.io.TcpListenerOptions.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.gloo.solo.io.TcpListenerOptions;
  return proto.gloo.solo.io.TcpListenerOptions.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.gloo.solo.io.TcpListenerOptions} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.gloo.solo.io.TcpListenerOptions}
 */
proto.gloo.solo.io.TcpListenerOptions.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 3:
      var value = new gloo_projects_gloo_api_v1_options_tcp_tcp_pb.TcpProxySettings;
      reader.readMessage(value,gloo_projects_gloo_api_v1_options_tcp_tcp_pb.TcpProxySettings.deserializeBinaryFromReader);
      msg.setTcpProxySettings(value);
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
proto.gloo.solo.io.TcpListenerOptions.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.gloo.solo.io.TcpListenerOptions.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.gloo.solo.io.TcpListenerOptions} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.gloo.solo.io.TcpListenerOptions.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getTcpProxySettings();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      gloo_projects_gloo_api_v1_options_tcp_tcp_pb.TcpProxySettings.serializeBinaryToWriter
    );
  }
};


/**
 * optional tcp.options.gloo.solo.io.TcpProxySettings tcp_proxy_settings = 3;
 * @return {?proto.tcp.options.gloo.solo.io.TcpProxySettings}
 */
proto.gloo.solo.io.TcpListenerOptions.prototype.getTcpProxySettings = function() {
  return /** @type{?proto.tcp.options.gloo.solo.io.TcpProxySettings} */ (
    jspb.Message.getWrapperField(this, gloo_projects_gloo_api_v1_options_tcp_tcp_pb.TcpProxySettings, 3));
};


/** @param {?proto.tcp.options.gloo.solo.io.TcpProxySettings|undefined} value */
proto.gloo.solo.io.TcpListenerOptions.prototype.setTcpProxySettings = function(value) {
  jspb.Message.setWrapperField(this, 3, value);
};


proto.gloo.solo.io.TcpListenerOptions.prototype.clearTcpProxySettings = function() {
  this.setTcpProxySettings(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.gloo.solo.io.TcpListenerOptions.prototype.hasTcpProxySettings = function() {
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
proto.gloo.solo.io.VirtualHostOptions = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.gloo.solo.io.VirtualHostOptions, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.gloo.solo.io.VirtualHostOptions.displayName = 'proto.gloo.solo.io.VirtualHostOptions';
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
proto.gloo.solo.io.VirtualHostOptions.prototype.toObject = function(opt_includeInstance) {
  return proto.gloo.solo.io.VirtualHostOptions.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.gloo.solo.io.VirtualHostOptions} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.gloo.solo.io.VirtualHostOptions.toObject = function(includeInstance, msg) {
  var f, obj = {
    extensions: (f = msg.getExtensions()) && gloo_projects_gloo_api_v1_extensions_pb.Extensions.toObject(includeInstance, f),
    retries: (f = msg.getRetries()) && gloo_projects_gloo_api_v1_options_retries_retries_pb.RetryPolicy.toObject(includeInstance, f),
    stats: (f = msg.getStats()) && gloo_projects_gloo_api_v1_options_stats_stats_pb.Stats.toObject(includeInstance, f),
    headerManipulation: (f = msg.getHeaderManipulation()) && gloo_projects_gloo_api_v1_options_headers_headers_pb.HeaderManipulation.toObject(includeInstance, f),
    cors: (f = msg.getCors()) && gloo_projects_gloo_api_v1_options_cors_cors_pb.CorsPolicy.toObject(includeInstance, f),
    transformations: (f = msg.getTransformations()) && gloo_projects_gloo_api_external_envoy_extensions_transformation_transformation_pb.RouteTransformations.toObject(includeInstance, f),
    ratelimitBasic: (f = msg.getRatelimitBasic()) && gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_pb.IngressRateLimit.toObject(includeInstance, f),
    ratelimit: (f = msg.getRatelimit()) && gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_pb.RateLimitVhostExtension.toObject(includeInstance, f),
    waf: (f = msg.getWaf()) && gloo_projects_gloo_api_v1_enterprise_options_waf_waf_pb.Settings.toObject(includeInstance, f),
    jwt: (f = msg.getJwt()) && gloo_projects_gloo_api_v1_enterprise_options_jwt_jwt_pb.VhostExtension.toObject(includeInstance, f),
    rbac: (f = msg.getRbac()) && gloo_projects_gloo_api_v1_enterprise_options_rbac_rbac_pb.ExtensionSettings.toObject(includeInstance, f),
    extauth: (f = msg.getExtauth()) && gloo_projects_gloo_api_v1_enterprise_options_extauth_v1_extauth_pb.ExtAuthExtension.toObject(includeInstance, f),
    dlp: (f = msg.getDlp()) && gloo_projects_gloo_api_v1_enterprise_options_dlp_dlp_pb.Config.toObject(includeInstance, f)
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
 * @return {!proto.gloo.solo.io.VirtualHostOptions}
 */
proto.gloo.solo.io.VirtualHostOptions.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.gloo.solo.io.VirtualHostOptions;
  return proto.gloo.solo.io.VirtualHostOptions.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.gloo.solo.io.VirtualHostOptions} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.gloo.solo.io.VirtualHostOptions}
 */
proto.gloo.solo.io.VirtualHostOptions.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new gloo_projects_gloo_api_v1_extensions_pb.Extensions;
      reader.readMessage(value,gloo_projects_gloo_api_v1_extensions_pb.Extensions.deserializeBinaryFromReader);
      msg.setExtensions(value);
      break;
    case 5:
      var value = new gloo_projects_gloo_api_v1_options_retries_retries_pb.RetryPolicy;
      reader.readMessage(value,gloo_projects_gloo_api_v1_options_retries_retries_pb.RetryPolicy.deserializeBinaryFromReader);
      msg.setRetries(value);
      break;
    case 10:
      var value = new gloo_projects_gloo_api_v1_options_stats_stats_pb.Stats;
      reader.readMessage(value,gloo_projects_gloo_api_v1_options_stats_stats_pb.Stats.deserializeBinaryFromReader);
      msg.setStats(value);
      break;
    case 2:
      var value = new gloo_projects_gloo_api_v1_options_headers_headers_pb.HeaderManipulation;
      reader.readMessage(value,gloo_projects_gloo_api_v1_options_headers_headers_pb.HeaderManipulation.deserializeBinaryFromReader);
      msg.setHeaderManipulation(value);
      break;
    case 3:
      var value = new gloo_projects_gloo_api_v1_options_cors_cors_pb.CorsPolicy;
      reader.readMessage(value,gloo_projects_gloo_api_v1_options_cors_cors_pb.CorsPolicy.deserializeBinaryFromReader);
      msg.setCors(value);
      break;
    case 4:
      var value = new gloo_projects_gloo_api_external_envoy_extensions_transformation_transformation_pb.RouteTransformations;
      reader.readMessage(value,gloo_projects_gloo_api_external_envoy_extensions_transformation_transformation_pb.RouteTransformations.deserializeBinaryFromReader);
      msg.setTransformations(value);
      break;
    case 6:
      var value = new gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_pb.IngressRateLimit;
      reader.readMessage(value,gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_pb.IngressRateLimit.deserializeBinaryFromReader);
      msg.setRatelimitBasic(value);
      break;
    case 7:
      var value = new gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_pb.RateLimitVhostExtension;
      reader.readMessage(value,gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_pb.RateLimitVhostExtension.deserializeBinaryFromReader);
      msg.setRatelimit(value);
      break;
    case 8:
      var value = new gloo_projects_gloo_api_v1_enterprise_options_waf_waf_pb.Settings;
      reader.readMessage(value,gloo_projects_gloo_api_v1_enterprise_options_waf_waf_pb.Settings.deserializeBinaryFromReader);
      msg.setWaf(value);
      break;
    case 9:
      var value = new gloo_projects_gloo_api_v1_enterprise_options_jwt_jwt_pb.VhostExtension;
      reader.readMessage(value,gloo_projects_gloo_api_v1_enterprise_options_jwt_jwt_pb.VhostExtension.deserializeBinaryFromReader);
      msg.setJwt(value);
      break;
    case 11:
      var value = new gloo_projects_gloo_api_v1_enterprise_options_rbac_rbac_pb.ExtensionSettings;
      reader.readMessage(value,gloo_projects_gloo_api_v1_enterprise_options_rbac_rbac_pb.ExtensionSettings.deserializeBinaryFromReader);
      msg.setRbac(value);
      break;
    case 12:
      var value = new gloo_projects_gloo_api_v1_enterprise_options_extauth_v1_extauth_pb.ExtAuthExtension;
      reader.readMessage(value,gloo_projects_gloo_api_v1_enterprise_options_extauth_v1_extauth_pb.ExtAuthExtension.deserializeBinaryFromReader);
      msg.setExtauth(value);
      break;
    case 13:
      var value = new gloo_projects_gloo_api_v1_enterprise_options_dlp_dlp_pb.Config;
      reader.readMessage(value,gloo_projects_gloo_api_v1_enterprise_options_dlp_dlp_pb.Config.deserializeBinaryFromReader);
      msg.setDlp(value);
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
proto.gloo.solo.io.VirtualHostOptions.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.gloo.solo.io.VirtualHostOptions.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.gloo.solo.io.VirtualHostOptions} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.gloo.solo.io.VirtualHostOptions.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getExtensions();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      gloo_projects_gloo_api_v1_extensions_pb.Extensions.serializeBinaryToWriter
    );
  }
  f = message.getRetries();
  if (f != null) {
    writer.writeMessage(
      5,
      f,
      gloo_projects_gloo_api_v1_options_retries_retries_pb.RetryPolicy.serializeBinaryToWriter
    );
  }
  f = message.getStats();
  if (f != null) {
    writer.writeMessage(
      10,
      f,
      gloo_projects_gloo_api_v1_options_stats_stats_pb.Stats.serializeBinaryToWriter
    );
  }
  f = message.getHeaderManipulation();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      gloo_projects_gloo_api_v1_options_headers_headers_pb.HeaderManipulation.serializeBinaryToWriter
    );
  }
  f = message.getCors();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      gloo_projects_gloo_api_v1_options_cors_cors_pb.CorsPolicy.serializeBinaryToWriter
    );
  }
  f = message.getTransformations();
  if (f != null) {
    writer.writeMessage(
      4,
      f,
      gloo_projects_gloo_api_external_envoy_extensions_transformation_transformation_pb.RouteTransformations.serializeBinaryToWriter
    );
  }
  f = message.getRatelimitBasic();
  if (f != null) {
    writer.writeMessage(
      6,
      f,
      gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_pb.IngressRateLimit.serializeBinaryToWriter
    );
  }
  f = message.getRatelimit();
  if (f != null) {
    writer.writeMessage(
      7,
      f,
      gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_pb.RateLimitVhostExtension.serializeBinaryToWriter
    );
  }
  f = message.getWaf();
  if (f != null) {
    writer.writeMessage(
      8,
      f,
      gloo_projects_gloo_api_v1_enterprise_options_waf_waf_pb.Settings.serializeBinaryToWriter
    );
  }
  f = message.getJwt();
  if (f != null) {
    writer.writeMessage(
      9,
      f,
      gloo_projects_gloo_api_v1_enterprise_options_jwt_jwt_pb.VhostExtension.serializeBinaryToWriter
    );
  }
  f = message.getRbac();
  if (f != null) {
    writer.writeMessage(
      11,
      f,
      gloo_projects_gloo_api_v1_enterprise_options_rbac_rbac_pb.ExtensionSettings.serializeBinaryToWriter
    );
  }
  f = message.getExtauth();
  if (f != null) {
    writer.writeMessage(
      12,
      f,
      gloo_projects_gloo_api_v1_enterprise_options_extauth_v1_extauth_pb.ExtAuthExtension.serializeBinaryToWriter
    );
  }
  f = message.getDlp();
  if (f != null) {
    writer.writeMessage(
      13,
      f,
      gloo_projects_gloo_api_v1_enterprise_options_dlp_dlp_pb.Config.serializeBinaryToWriter
    );
  }
};


/**
 * optional Extensions extensions = 1;
 * @return {?proto.gloo.solo.io.Extensions}
 */
proto.gloo.solo.io.VirtualHostOptions.prototype.getExtensions = function() {
  return /** @type{?proto.gloo.solo.io.Extensions} */ (
    jspb.Message.getWrapperField(this, gloo_projects_gloo_api_v1_extensions_pb.Extensions, 1));
};


/** @param {?proto.gloo.solo.io.Extensions|undefined} value */
proto.gloo.solo.io.VirtualHostOptions.prototype.setExtensions = function(value) {
  jspb.Message.setWrapperField(this, 1, value);
};


proto.gloo.solo.io.VirtualHostOptions.prototype.clearExtensions = function() {
  this.setExtensions(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.gloo.solo.io.VirtualHostOptions.prototype.hasExtensions = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional retries.options.gloo.solo.io.RetryPolicy retries = 5;
 * @return {?proto.retries.options.gloo.solo.io.RetryPolicy}
 */
proto.gloo.solo.io.VirtualHostOptions.prototype.getRetries = function() {
  return /** @type{?proto.retries.options.gloo.solo.io.RetryPolicy} */ (
    jspb.Message.getWrapperField(this, gloo_projects_gloo_api_v1_options_retries_retries_pb.RetryPolicy, 5));
};


/** @param {?proto.retries.options.gloo.solo.io.RetryPolicy|undefined} value */
proto.gloo.solo.io.VirtualHostOptions.prototype.setRetries = function(value) {
  jspb.Message.setWrapperField(this, 5, value);
};


proto.gloo.solo.io.VirtualHostOptions.prototype.clearRetries = function() {
  this.setRetries(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.gloo.solo.io.VirtualHostOptions.prototype.hasRetries = function() {
  return jspb.Message.getField(this, 5) != null;
};


/**
 * optional stats.options.gloo.solo.io.Stats stats = 10;
 * @return {?proto.stats.options.gloo.solo.io.Stats}
 */
proto.gloo.solo.io.VirtualHostOptions.prototype.getStats = function() {
  return /** @type{?proto.stats.options.gloo.solo.io.Stats} */ (
    jspb.Message.getWrapperField(this, gloo_projects_gloo_api_v1_options_stats_stats_pb.Stats, 10));
};


/** @param {?proto.stats.options.gloo.solo.io.Stats|undefined} value */
proto.gloo.solo.io.VirtualHostOptions.prototype.setStats = function(value) {
  jspb.Message.setWrapperField(this, 10, value);
};


proto.gloo.solo.io.VirtualHostOptions.prototype.clearStats = function() {
  this.setStats(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.gloo.solo.io.VirtualHostOptions.prototype.hasStats = function() {
  return jspb.Message.getField(this, 10) != null;
};


/**
 * optional headers.options.gloo.solo.io.HeaderManipulation header_manipulation = 2;
 * @return {?proto.headers.options.gloo.solo.io.HeaderManipulation}
 */
proto.gloo.solo.io.VirtualHostOptions.prototype.getHeaderManipulation = function() {
  return /** @type{?proto.headers.options.gloo.solo.io.HeaderManipulation} */ (
    jspb.Message.getWrapperField(this, gloo_projects_gloo_api_v1_options_headers_headers_pb.HeaderManipulation, 2));
};


/** @param {?proto.headers.options.gloo.solo.io.HeaderManipulation|undefined} value */
proto.gloo.solo.io.VirtualHostOptions.prototype.setHeaderManipulation = function(value) {
  jspb.Message.setWrapperField(this, 2, value);
};


proto.gloo.solo.io.VirtualHostOptions.prototype.clearHeaderManipulation = function() {
  this.setHeaderManipulation(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.gloo.solo.io.VirtualHostOptions.prototype.hasHeaderManipulation = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional cors.options.gloo.solo.io.CorsPolicy cors = 3;
 * @return {?proto.cors.options.gloo.solo.io.CorsPolicy}
 */
proto.gloo.solo.io.VirtualHostOptions.prototype.getCors = function() {
  return /** @type{?proto.cors.options.gloo.solo.io.CorsPolicy} */ (
    jspb.Message.getWrapperField(this, gloo_projects_gloo_api_v1_options_cors_cors_pb.CorsPolicy, 3));
};


/** @param {?proto.cors.options.gloo.solo.io.CorsPolicy|undefined} value */
proto.gloo.solo.io.VirtualHostOptions.prototype.setCors = function(value) {
  jspb.Message.setWrapperField(this, 3, value);
};


proto.gloo.solo.io.VirtualHostOptions.prototype.clearCors = function() {
  this.setCors(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.gloo.solo.io.VirtualHostOptions.prototype.hasCors = function() {
  return jspb.Message.getField(this, 3) != null;
};


/**
 * optional envoy.api.v2.filter.http.RouteTransformations transformations = 4;
 * @return {?proto.envoy.api.v2.filter.http.RouteTransformations}
 */
proto.gloo.solo.io.VirtualHostOptions.prototype.getTransformations = function() {
  return /** @type{?proto.envoy.api.v2.filter.http.RouteTransformations} */ (
    jspb.Message.getWrapperField(this, gloo_projects_gloo_api_external_envoy_extensions_transformation_transformation_pb.RouteTransformations, 4));
};


/** @param {?proto.envoy.api.v2.filter.http.RouteTransformations|undefined} value */
proto.gloo.solo.io.VirtualHostOptions.prototype.setTransformations = function(value) {
  jspb.Message.setWrapperField(this, 4, value);
};


proto.gloo.solo.io.VirtualHostOptions.prototype.clearTransformations = function() {
  this.setTransformations(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.gloo.solo.io.VirtualHostOptions.prototype.hasTransformations = function() {
  return jspb.Message.getField(this, 4) != null;
};


/**
 * optional ratelimit.options.gloo.solo.io.IngressRateLimit ratelimit_basic = 6;
 * @return {?proto.ratelimit.options.gloo.solo.io.IngressRateLimit}
 */
proto.gloo.solo.io.VirtualHostOptions.prototype.getRatelimitBasic = function() {
  return /** @type{?proto.ratelimit.options.gloo.solo.io.IngressRateLimit} */ (
    jspb.Message.getWrapperField(this, gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_pb.IngressRateLimit, 6));
};


/** @param {?proto.ratelimit.options.gloo.solo.io.IngressRateLimit|undefined} value */
proto.gloo.solo.io.VirtualHostOptions.prototype.setRatelimitBasic = function(value) {
  jspb.Message.setWrapperField(this, 6, value);
};


proto.gloo.solo.io.VirtualHostOptions.prototype.clearRatelimitBasic = function() {
  this.setRatelimitBasic(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.gloo.solo.io.VirtualHostOptions.prototype.hasRatelimitBasic = function() {
  return jspb.Message.getField(this, 6) != null;
};


/**
 * optional ratelimit.options.gloo.solo.io.RateLimitVhostExtension ratelimit = 7;
 * @return {?proto.ratelimit.options.gloo.solo.io.RateLimitVhostExtension}
 */
proto.gloo.solo.io.VirtualHostOptions.prototype.getRatelimit = function() {
  return /** @type{?proto.ratelimit.options.gloo.solo.io.RateLimitVhostExtension} */ (
    jspb.Message.getWrapperField(this, gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_pb.RateLimitVhostExtension, 7));
};


/** @param {?proto.ratelimit.options.gloo.solo.io.RateLimitVhostExtension|undefined} value */
proto.gloo.solo.io.VirtualHostOptions.prototype.setRatelimit = function(value) {
  jspb.Message.setWrapperField(this, 7, value);
};


proto.gloo.solo.io.VirtualHostOptions.prototype.clearRatelimit = function() {
  this.setRatelimit(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.gloo.solo.io.VirtualHostOptions.prototype.hasRatelimit = function() {
  return jspb.Message.getField(this, 7) != null;
};


/**
 * optional waf.options.gloo.solo.io.Settings waf = 8;
 * @return {?proto.waf.options.gloo.solo.io.Settings}
 */
proto.gloo.solo.io.VirtualHostOptions.prototype.getWaf = function() {
  return /** @type{?proto.waf.options.gloo.solo.io.Settings} */ (
    jspb.Message.getWrapperField(this, gloo_projects_gloo_api_v1_enterprise_options_waf_waf_pb.Settings, 8));
};


/** @param {?proto.waf.options.gloo.solo.io.Settings|undefined} value */
proto.gloo.solo.io.VirtualHostOptions.prototype.setWaf = function(value) {
  jspb.Message.setWrapperField(this, 8, value);
};


proto.gloo.solo.io.VirtualHostOptions.prototype.clearWaf = function() {
  this.setWaf(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.gloo.solo.io.VirtualHostOptions.prototype.hasWaf = function() {
  return jspb.Message.getField(this, 8) != null;
};


/**
 * optional jwt.options.gloo.solo.io.VhostExtension jwt = 9;
 * @return {?proto.jwt.options.gloo.solo.io.VhostExtension}
 */
proto.gloo.solo.io.VirtualHostOptions.prototype.getJwt = function() {
  return /** @type{?proto.jwt.options.gloo.solo.io.VhostExtension} */ (
    jspb.Message.getWrapperField(this, gloo_projects_gloo_api_v1_enterprise_options_jwt_jwt_pb.VhostExtension, 9));
};


/** @param {?proto.jwt.options.gloo.solo.io.VhostExtension|undefined} value */
proto.gloo.solo.io.VirtualHostOptions.prototype.setJwt = function(value) {
  jspb.Message.setWrapperField(this, 9, value);
};


proto.gloo.solo.io.VirtualHostOptions.prototype.clearJwt = function() {
  this.setJwt(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.gloo.solo.io.VirtualHostOptions.prototype.hasJwt = function() {
  return jspb.Message.getField(this, 9) != null;
};


/**
 * optional rbac.options.gloo.solo.io.ExtensionSettings rbac = 11;
 * @return {?proto.rbac.options.gloo.solo.io.ExtensionSettings}
 */
proto.gloo.solo.io.VirtualHostOptions.prototype.getRbac = function() {
  return /** @type{?proto.rbac.options.gloo.solo.io.ExtensionSettings} */ (
    jspb.Message.getWrapperField(this, gloo_projects_gloo_api_v1_enterprise_options_rbac_rbac_pb.ExtensionSettings, 11));
};


/** @param {?proto.rbac.options.gloo.solo.io.ExtensionSettings|undefined} value */
proto.gloo.solo.io.VirtualHostOptions.prototype.setRbac = function(value) {
  jspb.Message.setWrapperField(this, 11, value);
};


proto.gloo.solo.io.VirtualHostOptions.prototype.clearRbac = function() {
  this.setRbac(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.gloo.solo.io.VirtualHostOptions.prototype.hasRbac = function() {
  return jspb.Message.getField(this, 11) != null;
};


/**
 * optional enterprise.gloo.solo.io.ExtAuthExtension extauth = 12;
 * @return {?proto.enterprise.gloo.solo.io.ExtAuthExtension}
 */
proto.gloo.solo.io.VirtualHostOptions.prototype.getExtauth = function() {
  return /** @type{?proto.enterprise.gloo.solo.io.ExtAuthExtension} */ (
    jspb.Message.getWrapperField(this, gloo_projects_gloo_api_v1_enterprise_options_extauth_v1_extauth_pb.ExtAuthExtension, 12));
};


/** @param {?proto.enterprise.gloo.solo.io.ExtAuthExtension|undefined} value */
proto.gloo.solo.io.VirtualHostOptions.prototype.setExtauth = function(value) {
  jspb.Message.setWrapperField(this, 12, value);
};


proto.gloo.solo.io.VirtualHostOptions.prototype.clearExtauth = function() {
  this.setExtauth(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.gloo.solo.io.VirtualHostOptions.prototype.hasExtauth = function() {
  return jspb.Message.getField(this, 12) != null;
};


/**
 * optional dlp.options.gloo.solo.io.Config dlp = 13;
 * @return {?proto.dlp.options.gloo.solo.io.Config}
 */
proto.gloo.solo.io.VirtualHostOptions.prototype.getDlp = function() {
  return /** @type{?proto.dlp.options.gloo.solo.io.Config} */ (
    jspb.Message.getWrapperField(this, gloo_projects_gloo_api_v1_enterprise_options_dlp_dlp_pb.Config, 13));
};


/** @param {?proto.dlp.options.gloo.solo.io.Config|undefined} value */
proto.gloo.solo.io.VirtualHostOptions.prototype.setDlp = function(value) {
  jspb.Message.setWrapperField(this, 13, value);
};


proto.gloo.solo.io.VirtualHostOptions.prototype.clearDlp = function() {
  this.setDlp(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.gloo.solo.io.VirtualHostOptions.prototype.hasDlp = function() {
  return jspb.Message.getField(this, 13) != null;
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
proto.gloo.solo.io.RouteOptions = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.gloo.solo.io.RouteOptions.repeatedFields_, proto.gloo.solo.io.RouteOptions.oneofGroups_);
};
goog.inherits(proto.gloo.solo.io.RouteOptions, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.gloo.solo.io.RouteOptions.displayName = 'proto.gloo.solo.io.RouteOptions';
}
/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.gloo.solo.io.RouteOptions.repeatedFields_ = [21];

/**
 * Oneof group definitions for this message. Each group defines the field
 * numbers belonging to that group. When of these fields' value is set, all
 * other fields in the group are cleared. During deserialization, if multiple
 * fields are encountered for a group, only the last value seen will be kept.
 * @private {!Array<!Array<number>>}
 * @const
 */
proto.gloo.solo.io.RouteOptions.oneofGroups_ = [[10,19]];

/**
 * @enum {number}
 */
proto.gloo.solo.io.RouteOptions.HostRewriteTypeCase = {
  HOST_REWRITE_TYPE_NOT_SET: 0,
  HOST_REWRITE: 10,
  AUTO_HOST_REWRITE: 19
};

/**
 * @return {proto.gloo.solo.io.RouteOptions.HostRewriteTypeCase}
 */
proto.gloo.solo.io.RouteOptions.prototype.getHostRewriteTypeCase = function() {
  return /** @type {proto.gloo.solo.io.RouteOptions.HostRewriteTypeCase} */(jspb.Message.computeOneofCase(this, proto.gloo.solo.io.RouteOptions.oneofGroups_[0]));
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
proto.gloo.solo.io.RouteOptions.prototype.toObject = function(opt_includeInstance) {
  return proto.gloo.solo.io.RouteOptions.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.gloo.solo.io.RouteOptions} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.gloo.solo.io.RouteOptions.toObject = function(includeInstance, msg) {
  var f, obj = {
    transformations: (f = msg.getTransformations()) && gloo_projects_gloo_api_external_envoy_extensions_transformation_transformation_pb.RouteTransformations.toObject(includeInstance, f),
    faults: (f = msg.getFaults()) && gloo_projects_gloo_api_v1_options_faultinjection_fault_pb.RouteFaults.toObject(includeInstance, f),
    prefixRewrite: (f = msg.getPrefixRewrite()) && google_protobuf_wrappers_pb.StringValue.toObject(includeInstance, f),
    timeout: (f = msg.getTimeout()) && google_protobuf_duration_pb.Duration.toObject(includeInstance, f),
    retries: (f = msg.getRetries()) && gloo_projects_gloo_api_v1_options_retries_retries_pb.RetryPolicy.toObject(includeInstance, f),
    extensions: (f = msg.getExtensions()) && gloo_projects_gloo_api_v1_extensions_pb.Extensions.toObject(includeInstance, f),
    tracing: (f = msg.getTracing()) && gloo_projects_gloo_api_v1_options_tracing_tracing_pb.RouteTracingSettings.toObject(includeInstance, f),
    shadowing: (f = msg.getShadowing()) && gloo_projects_gloo_api_v1_options_shadowing_shadowing_pb.RouteShadowing.toObject(includeInstance, f),
    headerManipulation: (f = msg.getHeaderManipulation()) && gloo_projects_gloo_api_v1_options_headers_headers_pb.HeaderManipulation.toObject(includeInstance, f),
    hostRewrite: jspb.Message.getFieldWithDefault(msg, 10, ""),
    autoHostRewrite: (f = msg.getAutoHostRewrite()) && google_protobuf_wrappers_pb.BoolValue.toObject(includeInstance, f),
    cors: (f = msg.getCors()) && gloo_projects_gloo_api_v1_options_cors_cors_pb.CorsPolicy.toObject(includeInstance, f),
    lbHash: (f = msg.getLbHash()) && gloo_projects_gloo_api_v1_options_lbhash_lbhash_pb.RouteActionHashConfig.toObject(includeInstance, f),
    upgradesList: jspb.Message.toObjectList(msg.getUpgradesList(),
    gloo_projects_gloo_api_v1_options_protocol_upgrade_protocol_upgrade_pb.ProtocolUpgradeConfig.toObject, includeInstance),
    ratelimitBasic: (f = msg.getRatelimitBasic()) && gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_pb.IngressRateLimit.toObject(includeInstance, f),
    ratelimit: (f = msg.getRatelimit()) && gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_pb.RateLimitRouteExtension.toObject(includeInstance, f),
    waf: (f = msg.getWaf()) && gloo_projects_gloo_api_v1_enterprise_options_waf_waf_pb.Settings.toObject(includeInstance, f),
    jwt: (f = msg.getJwt()) && gloo_projects_gloo_api_v1_enterprise_options_jwt_jwt_pb.RouteExtension.toObject(includeInstance, f),
    rbac: (f = msg.getRbac()) && gloo_projects_gloo_api_v1_enterprise_options_rbac_rbac_pb.ExtensionSettings.toObject(includeInstance, f),
    extauth: (f = msg.getExtauth()) && gloo_projects_gloo_api_v1_enterprise_options_extauth_v1_extauth_pb.ExtAuthExtension.toObject(includeInstance, f),
    dlp: (f = msg.getDlp()) && gloo_projects_gloo_api_v1_enterprise_options_dlp_dlp_pb.Config.toObject(includeInstance, f)
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
 * @return {!proto.gloo.solo.io.RouteOptions}
 */
proto.gloo.solo.io.RouteOptions.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.gloo.solo.io.RouteOptions;
  return proto.gloo.solo.io.RouteOptions.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.gloo.solo.io.RouteOptions} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.gloo.solo.io.RouteOptions}
 */
proto.gloo.solo.io.RouteOptions.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new gloo_projects_gloo_api_external_envoy_extensions_transformation_transformation_pb.RouteTransformations;
      reader.readMessage(value,gloo_projects_gloo_api_external_envoy_extensions_transformation_transformation_pb.RouteTransformations.deserializeBinaryFromReader);
      msg.setTransformations(value);
      break;
    case 2:
      var value = new gloo_projects_gloo_api_v1_options_faultinjection_fault_pb.RouteFaults;
      reader.readMessage(value,gloo_projects_gloo_api_v1_options_faultinjection_fault_pb.RouteFaults.deserializeBinaryFromReader);
      msg.setFaults(value);
      break;
    case 3:
      var value = new google_protobuf_wrappers_pb.StringValue;
      reader.readMessage(value,google_protobuf_wrappers_pb.StringValue.deserializeBinaryFromReader);
      msg.setPrefixRewrite(value);
      break;
    case 4:
      var value = new google_protobuf_duration_pb.Duration;
      reader.readMessage(value,google_protobuf_duration_pb.Duration.deserializeBinaryFromReader);
      msg.setTimeout(value);
      break;
    case 5:
      var value = new gloo_projects_gloo_api_v1_options_retries_retries_pb.RetryPolicy;
      reader.readMessage(value,gloo_projects_gloo_api_v1_options_retries_retries_pb.RetryPolicy.deserializeBinaryFromReader);
      msg.setRetries(value);
      break;
    case 6:
      var value = new gloo_projects_gloo_api_v1_extensions_pb.Extensions;
      reader.readMessage(value,gloo_projects_gloo_api_v1_extensions_pb.Extensions.deserializeBinaryFromReader);
      msg.setExtensions(value);
      break;
    case 7:
      var value = new gloo_projects_gloo_api_v1_options_tracing_tracing_pb.RouteTracingSettings;
      reader.readMessage(value,gloo_projects_gloo_api_v1_options_tracing_tracing_pb.RouteTracingSettings.deserializeBinaryFromReader);
      msg.setTracing(value);
      break;
    case 8:
      var value = new gloo_projects_gloo_api_v1_options_shadowing_shadowing_pb.RouteShadowing;
      reader.readMessage(value,gloo_projects_gloo_api_v1_options_shadowing_shadowing_pb.RouteShadowing.deserializeBinaryFromReader);
      msg.setShadowing(value);
      break;
    case 9:
      var value = new gloo_projects_gloo_api_v1_options_headers_headers_pb.HeaderManipulation;
      reader.readMessage(value,gloo_projects_gloo_api_v1_options_headers_headers_pb.HeaderManipulation.deserializeBinaryFromReader);
      msg.setHeaderManipulation(value);
      break;
    case 10:
      var value = /** @type {string} */ (reader.readString());
      msg.setHostRewrite(value);
      break;
    case 19:
      var value = new google_protobuf_wrappers_pb.BoolValue;
      reader.readMessage(value,google_protobuf_wrappers_pb.BoolValue.deserializeBinaryFromReader);
      msg.setAutoHostRewrite(value);
      break;
    case 11:
      var value = new gloo_projects_gloo_api_v1_options_cors_cors_pb.CorsPolicy;
      reader.readMessage(value,gloo_projects_gloo_api_v1_options_cors_cors_pb.CorsPolicy.deserializeBinaryFromReader);
      msg.setCors(value);
      break;
    case 12:
      var value = new gloo_projects_gloo_api_v1_options_lbhash_lbhash_pb.RouteActionHashConfig;
      reader.readMessage(value,gloo_projects_gloo_api_v1_options_lbhash_lbhash_pb.RouteActionHashConfig.deserializeBinaryFromReader);
      msg.setLbHash(value);
      break;
    case 21:
      var value = new gloo_projects_gloo_api_v1_options_protocol_upgrade_protocol_upgrade_pb.ProtocolUpgradeConfig;
      reader.readMessage(value,gloo_projects_gloo_api_v1_options_protocol_upgrade_protocol_upgrade_pb.ProtocolUpgradeConfig.deserializeBinaryFromReader);
      msg.addUpgrades(value);
      break;
    case 13:
      var value = new gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_pb.IngressRateLimit;
      reader.readMessage(value,gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_pb.IngressRateLimit.deserializeBinaryFromReader);
      msg.setRatelimitBasic(value);
      break;
    case 14:
      var value = new gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_pb.RateLimitRouteExtension;
      reader.readMessage(value,gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_pb.RateLimitRouteExtension.deserializeBinaryFromReader);
      msg.setRatelimit(value);
      break;
    case 15:
      var value = new gloo_projects_gloo_api_v1_enterprise_options_waf_waf_pb.Settings;
      reader.readMessage(value,gloo_projects_gloo_api_v1_enterprise_options_waf_waf_pb.Settings.deserializeBinaryFromReader);
      msg.setWaf(value);
      break;
    case 16:
      var value = new gloo_projects_gloo_api_v1_enterprise_options_jwt_jwt_pb.RouteExtension;
      reader.readMessage(value,gloo_projects_gloo_api_v1_enterprise_options_jwt_jwt_pb.RouteExtension.deserializeBinaryFromReader);
      msg.setJwt(value);
      break;
    case 17:
      var value = new gloo_projects_gloo_api_v1_enterprise_options_rbac_rbac_pb.ExtensionSettings;
      reader.readMessage(value,gloo_projects_gloo_api_v1_enterprise_options_rbac_rbac_pb.ExtensionSettings.deserializeBinaryFromReader);
      msg.setRbac(value);
      break;
    case 18:
      var value = new gloo_projects_gloo_api_v1_enterprise_options_extauth_v1_extauth_pb.ExtAuthExtension;
      reader.readMessage(value,gloo_projects_gloo_api_v1_enterprise_options_extauth_v1_extauth_pb.ExtAuthExtension.deserializeBinaryFromReader);
      msg.setExtauth(value);
      break;
    case 20:
      var value = new gloo_projects_gloo_api_v1_enterprise_options_dlp_dlp_pb.Config;
      reader.readMessage(value,gloo_projects_gloo_api_v1_enterprise_options_dlp_dlp_pb.Config.deserializeBinaryFromReader);
      msg.setDlp(value);
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
proto.gloo.solo.io.RouteOptions.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.gloo.solo.io.RouteOptions.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.gloo.solo.io.RouteOptions} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.gloo.solo.io.RouteOptions.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getTransformations();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      gloo_projects_gloo_api_external_envoy_extensions_transformation_transformation_pb.RouteTransformations.serializeBinaryToWriter
    );
  }
  f = message.getFaults();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      gloo_projects_gloo_api_v1_options_faultinjection_fault_pb.RouteFaults.serializeBinaryToWriter
    );
  }
  f = message.getPrefixRewrite();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      google_protobuf_wrappers_pb.StringValue.serializeBinaryToWriter
    );
  }
  f = message.getTimeout();
  if (f != null) {
    writer.writeMessage(
      4,
      f,
      google_protobuf_duration_pb.Duration.serializeBinaryToWriter
    );
  }
  f = message.getRetries();
  if (f != null) {
    writer.writeMessage(
      5,
      f,
      gloo_projects_gloo_api_v1_options_retries_retries_pb.RetryPolicy.serializeBinaryToWriter
    );
  }
  f = message.getExtensions();
  if (f != null) {
    writer.writeMessage(
      6,
      f,
      gloo_projects_gloo_api_v1_extensions_pb.Extensions.serializeBinaryToWriter
    );
  }
  f = message.getTracing();
  if (f != null) {
    writer.writeMessage(
      7,
      f,
      gloo_projects_gloo_api_v1_options_tracing_tracing_pb.RouteTracingSettings.serializeBinaryToWriter
    );
  }
  f = message.getShadowing();
  if (f != null) {
    writer.writeMessage(
      8,
      f,
      gloo_projects_gloo_api_v1_options_shadowing_shadowing_pb.RouteShadowing.serializeBinaryToWriter
    );
  }
  f = message.getHeaderManipulation();
  if (f != null) {
    writer.writeMessage(
      9,
      f,
      gloo_projects_gloo_api_v1_options_headers_headers_pb.HeaderManipulation.serializeBinaryToWriter
    );
  }
  f = /** @type {string} */ (jspb.Message.getField(message, 10));
  if (f != null) {
    writer.writeString(
      10,
      f
    );
  }
  f = message.getAutoHostRewrite();
  if (f != null) {
    writer.writeMessage(
      19,
      f,
      google_protobuf_wrappers_pb.BoolValue.serializeBinaryToWriter
    );
  }
  f = message.getCors();
  if (f != null) {
    writer.writeMessage(
      11,
      f,
      gloo_projects_gloo_api_v1_options_cors_cors_pb.CorsPolicy.serializeBinaryToWriter
    );
  }
  f = message.getLbHash();
  if (f != null) {
    writer.writeMessage(
      12,
      f,
      gloo_projects_gloo_api_v1_options_lbhash_lbhash_pb.RouteActionHashConfig.serializeBinaryToWriter
    );
  }
  f = message.getUpgradesList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      21,
      f,
      gloo_projects_gloo_api_v1_options_protocol_upgrade_protocol_upgrade_pb.ProtocolUpgradeConfig.serializeBinaryToWriter
    );
  }
  f = message.getRatelimitBasic();
  if (f != null) {
    writer.writeMessage(
      13,
      f,
      gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_pb.IngressRateLimit.serializeBinaryToWriter
    );
  }
  f = message.getRatelimit();
  if (f != null) {
    writer.writeMessage(
      14,
      f,
      gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_pb.RateLimitRouteExtension.serializeBinaryToWriter
    );
  }
  f = message.getWaf();
  if (f != null) {
    writer.writeMessage(
      15,
      f,
      gloo_projects_gloo_api_v1_enterprise_options_waf_waf_pb.Settings.serializeBinaryToWriter
    );
  }
  f = message.getJwt();
  if (f != null) {
    writer.writeMessage(
      16,
      f,
      gloo_projects_gloo_api_v1_enterprise_options_jwt_jwt_pb.RouteExtension.serializeBinaryToWriter
    );
  }
  f = message.getRbac();
  if (f != null) {
    writer.writeMessage(
      17,
      f,
      gloo_projects_gloo_api_v1_enterprise_options_rbac_rbac_pb.ExtensionSettings.serializeBinaryToWriter
    );
  }
  f = message.getExtauth();
  if (f != null) {
    writer.writeMessage(
      18,
      f,
      gloo_projects_gloo_api_v1_enterprise_options_extauth_v1_extauth_pb.ExtAuthExtension.serializeBinaryToWriter
    );
  }
  f = message.getDlp();
  if (f != null) {
    writer.writeMessage(
      20,
      f,
      gloo_projects_gloo_api_v1_enterprise_options_dlp_dlp_pb.Config.serializeBinaryToWriter
    );
  }
};


/**
 * optional envoy.api.v2.filter.http.RouteTransformations transformations = 1;
 * @return {?proto.envoy.api.v2.filter.http.RouteTransformations}
 */
proto.gloo.solo.io.RouteOptions.prototype.getTransformations = function() {
  return /** @type{?proto.envoy.api.v2.filter.http.RouteTransformations} */ (
    jspb.Message.getWrapperField(this, gloo_projects_gloo_api_external_envoy_extensions_transformation_transformation_pb.RouteTransformations, 1));
};


/** @param {?proto.envoy.api.v2.filter.http.RouteTransformations|undefined} value */
proto.gloo.solo.io.RouteOptions.prototype.setTransformations = function(value) {
  jspb.Message.setWrapperField(this, 1, value);
};


proto.gloo.solo.io.RouteOptions.prototype.clearTransformations = function() {
  this.setTransformations(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.gloo.solo.io.RouteOptions.prototype.hasTransformations = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional fault.options.gloo.solo.io.RouteFaults faults = 2;
 * @return {?proto.fault.options.gloo.solo.io.RouteFaults}
 */
proto.gloo.solo.io.RouteOptions.prototype.getFaults = function() {
  return /** @type{?proto.fault.options.gloo.solo.io.RouteFaults} */ (
    jspb.Message.getWrapperField(this, gloo_projects_gloo_api_v1_options_faultinjection_fault_pb.RouteFaults, 2));
};


/** @param {?proto.fault.options.gloo.solo.io.RouteFaults|undefined} value */
proto.gloo.solo.io.RouteOptions.prototype.setFaults = function(value) {
  jspb.Message.setWrapperField(this, 2, value);
};


proto.gloo.solo.io.RouteOptions.prototype.clearFaults = function() {
  this.setFaults(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.gloo.solo.io.RouteOptions.prototype.hasFaults = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional google.protobuf.StringValue prefix_rewrite = 3;
 * @return {?proto.google.protobuf.StringValue}
 */
proto.gloo.solo.io.RouteOptions.prototype.getPrefixRewrite = function() {
  return /** @type{?proto.google.protobuf.StringValue} */ (
    jspb.Message.getWrapperField(this, google_protobuf_wrappers_pb.StringValue, 3));
};


/** @param {?proto.google.protobuf.StringValue|undefined} value */
proto.gloo.solo.io.RouteOptions.prototype.setPrefixRewrite = function(value) {
  jspb.Message.setWrapperField(this, 3, value);
};


proto.gloo.solo.io.RouteOptions.prototype.clearPrefixRewrite = function() {
  this.setPrefixRewrite(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.gloo.solo.io.RouteOptions.prototype.hasPrefixRewrite = function() {
  return jspb.Message.getField(this, 3) != null;
};


/**
 * optional google.protobuf.Duration timeout = 4;
 * @return {?proto.google.protobuf.Duration}
 */
proto.gloo.solo.io.RouteOptions.prototype.getTimeout = function() {
  return /** @type{?proto.google.protobuf.Duration} */ (
    jspb.Message.getWrapperField(this, google_protobuf_duration_pb.Duration, 4));
};


/** @param {?proto.google.protobuf.Duration|undefined} value */
proto.gloo.solo.io.RouteOptions.prototype.setTimeout = function(value) {
  jspb.Message.setWrapperField(this, 4, value);
};


proto.gloo.solo.io.RouteOptions.prototype.clearTimeout = function() {
  this.setTimeout(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.gloo.solo.io.RouteOptions.prototype.hasTimeout = function() {
  return jspb.Message.getField(this, 4) != null;
};


/**
 * optional retries.options.gloo.solo.io.RetryPolicy retries = 5;
 * @return {?proto.retries.options.gloo.solo.io.RetryPolicy}
 */
proto.gloo.solo.io.RouteOptions.prototype.getRetries = function() {
  return /** @type{?proto.retries.options.gloo.solo.io.RetryPolicy} */ (
    jspb.Message.getWrapperField(this, gloo_projects_gloo_api_v1_options_retries_retries_pb.RetryPolicy, 5));
};


/** @param {?proto.retries.options.gloo.solo.io.RetryPolicy|undefined} value */
proto.gloo.solo.io.RouteOptions.prototype.setRetries = function(value) {
  jspb.Message.setWrapperField(this, 5, value);
};


proto.gloo.solo.io.RouteOptions.prototype.clearRetries = function() {
  this.setRetries(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.gloo.solo.io.RouteOptions.prototype.hasRetries = function() {
  return jspb.Message.getField(this, 5) != null;
};


/**
 * optional Extensions extensions = 6;
 * @return {?proto.gloo.solo.io.Extensions}
 */
proto.gloo.solo.io.RouteOptions.prototype.getExtensions = function() {
  return /** @type{?proto.gloo.solo.io.Extensions} */ (
    jspb.Message.getWrapperField(this, gloo_projects_gloo_api_v1_extensions_pb.Extensions, 6));
};


/** @param {?proto.gloo.solo.io.Extensions|undefined} value */
proto.gloo.solo.io.RouteOptions.prototype.setExtensions = function(value) {
  jspb.Message.setWrapperField(this, 6, value);
};


proto.gloo.solo.io.RouteOptions.prototype.clearExtensions = function() {
  this.setExtensions(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.gloo.solo.io.RouteOptions.prototype.hasExtensions = function() {
  return jspb.Message.getField(this, 6) != null;
};


/**
 * optional tracing.options.gloo.solo.io.RouteTracingSettings tracing = 7;
 * @return {?proto.tracing.options.gloo.solo.io.RouteTracingSettings}
 */
proto.gloo.solo.io.RouteOptions.prototype.getTracing = function() {
  return /** @type{?proto.tracing.options.gloo.solo.io.RouteTracingSettings} */ (
    jspb.Message.getWrapperField(this, gloo_projects_gloo_api_v1_options_tracing_tracing_pb.RouteTracingSettings, 7));
};


/** @param {?proto.tracing.options.gloo.solo.io.RouteTracingSettings|undefined} value */
proto.gloo.solo.io.RouteOptions.prototype.setTracing = function(value) {
  jspb.Message.setWrapperField(this, 7, value);
};


proto.gloo.solo.io.RouteOptions.prototype.clearTracing = function() {
  this.setTracing(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.gloo.solo.io.RouteOptions.prototype.hasTracing = function() {
  return jspb.Message.getField(this, 7) != null;
};


/**
 * optional shadowing.options.gloo.solo.io.RouteShadowing shadowing = 8;
 * @return {?proto.shadowing.options.gloo.solo.io.RouteShadowing}
 */
proto.gloo.solo.io.RouteOptions.prototype.getShadowing = function() {
  return /** @type{?proto.shadowing.options.gloo.solo.io.RouteShadowing} */ (
    jspb.Message.getWrapperField(this, gloo_projects_gloo_api_v1_options_shadowing_shadowing_pb.RouteShadowing, 8));
};


/** @param {?proto.shadowing.options.gloo.solo.io.RouteShadowing|undefined} value */
proto.gloo.solo.io.RouteOptions.prototype.setShadowing = function(value) {
  jspb.Message.setWrapperField(this, 8, value);
};


proto.gloo.solo.io.RouteOptions.prototype.clearShadowing = function() {
  this.setShadowing(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.gloo.solo.io.RouteOptions.prototype.hasShadowing = function() {
  return jspb.Message.getField(this, 8) != null;
};


/**
 * optional headers.options.gloo.solo.io.HeaderManipulation header_manipulation = 9;
 * @return {?proto.headers.options.gloo.solo.io.HeaderManipulation}
 */
proto.gloo.solo.io.RouteOptions.prototype.getHeaderManipulation = function() {
  return /** @type{?proto.headers.options.gloo.solo.io.HeaderManipulation} */ (
    jspb.Message.getWrapperField(this, gloo_projects_gloo_api_v1_options_headers_headers_pb.HeaderManipulation, 9));
};


/** @param {?proto.headers.options.gloo.solo.io.HeaderManipulation|undefined} value */
proto.gloo.solo.io.RouteOptions.prototype.setHeaderManipulation = function(value) {
  jspb.Message.setWrapperField(this, 9, value);
};


proto.gloo.solo.io.RouteOptions.prototype.clearHeaderManipulation = function() {
  this.setHeaderManipulation(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.gloo.solo.io.RouteOptions.prototype.hasHeaderManipulation = function() {
  return jspb.Message.getField(this, 9) != null;
};


/**
 * optional string host_rewrite = 10;
 * @return {string}
 */
proto.gloo.solo.io.RouteOptions.prototype.getHostRewrite = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 10, ""));
};


/** @param {string} value */
proto.gloo.solo.io.RouteOptions.prototype.setHostRewrite = function(value) {
  jspb.Message.setOneofField(this, 10, proto.gloo.solo.io.RouteOptions.oneofGroups_[0], value);
};


proto.gloo.solo.io.RouteOptions.prototype.clearHostRewrite = function() {
  jspb.Message.setOneofField(this, 10, proto.gloo.solo.io.RouteOptions.oneofGroups_[0], undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.gloo.solo.io.RouteOptions.prototype.hasHostRewrite = function() {
  return jspb.Message.getField(this, 10) != null;
};


/**
 * optional google.protobuf.BoolValue auto_host_rewrite = 19;
 * @return {?proto.google.protobuf.BoolValue}
 */
proto.gloo.solo.io.RouteOptions.prototype.getAutoHostRewrite = function() {
  return /** @type{?proto.google.protobuf.BoolValue} */ (
    jspb.Message.getWrapperField(this, google_protobuf_wrappers_pb.BoolValue, 19));
};


/** @param {?proto.google.protobuf.BoolValue|undefined} value */
proto.gloo.solo.io.RouteOptions.prototype.setAutoHostRewrite = function(value) {
  jspb.Message.setOneofWrapperField(this, 19, proto.gloo.solo.io.RouteOptions.oneofGroups_[0], value);
};


proto.gloo.solo.io.RouteOptions.prototype.clearAutoHostRewrite = function() {
  this.setAutoHostRewrite(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.gloo.solo.io.RouteOptions.prototype.hasAutoHostRewrite = function() {
  return jspb.Message.getField(this, 19) != null;
};


/**
 * optional cors.options.gloo.solo.io.CorsPolicy cors = 11;
 * @return {?proto.cors.options.gloo.solo.io.CorsPolicy}
 */
proto.gloo.solo.io.RouteOptions.prototype.getCors = function() {
  return /** @type{?proto.cors.options.gloo.solo.io.CorsPolicy} */ (
    jspb.Message.getWrapperField(this, gloo_projects_gloo_api_v1_options_cors_cors_pb.CorsPolicy, 11));
};


/** @param {?proto.cors.options.gloo.solo.io.CorsPolicy|undefined} value */
proto.gloo.solo.io.RouteOptions.prototype.setCors = function(value) {
  jspb.Message.setWrapperField(this, 11, value);
};


proto.gloo.solo.io.RouteOptions.prototype.clearCors = function() {
  this.setCors(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.gloo.solo.io.RouteOptions.prototype.hasCors = function() {
  return jspb.Message.getField(this, 11) != null;
};


/**
 * optional lbhash.options.gloo.solo.io.RouteActionHashConfig lb_hash = 12;
 * @return {?proto.lbhash.options.gloo.solo.io.RouteActionHashConfig}
 */
proto.gloo.solo.io.RouteOptions.prototype.getLbHash = function() {
  return /** @type{?proto.lbhash.options.gloo.solo.io.RouteActionHashConfig} */ (
    jspb.Message.getWrapperField(this, gloo_projects_gloo_api_v1_options_lbhash_lbhash_pb.RouteActionHashConfig, 12));
};


/** @param {?proto.lbhash.options.gloo.solo.io.RouteActionHashConfig|undefined} value */
proto.gloo.solo.io.RouteOptions.prototype.setLbHash = function(value) {
  jspb.Message.setWrapperField(this, 12, value);
};


proto.gloo.solo.io.RouteOptions.prototype.clearLbHash = function() {
  this.setLbHash(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.gloo.solo.io.RouteOptions.prototype.hasLbHash = function() {
  return jspb.Message.getField(this, 12) != null;
};


/**
 * repeated protocol_upgrade.options.gloo.solo.io.ProtocolUpgradeConfig upgrades = 21;
 * @return {!Array<!proto.protocol_upgrade.options.gloo.solo.io.ProtocolUpgradeConfig>}
 */
proto.gloo.solo.io.RouteOptions.prototype.getUpgradesList = function() {
  return /** @type{!Array<!proto.protocol_upgrade.options.gloo.solo.io.ProtocolUpgradeConfig>} */ (
    jspb.Message.getRepeatedWrapperField(this, gloo_projects_gloo_api_v1_options_protocol_upgrade_protocol_upgrade_pb.ProtocolUpgradeConfig, 21));
};


/** @param {!Array<!proto.protocol_upgrade.options.gloo.solo.io.ProtocolUpgradeConfig>} value */
proto.gloo.solo.io.RouteOptions.prototype.setUpgradesList = function(value) {
  jspb.Message.setRepeatedWrapperField(this, 21, value);
};


/**
 * @param {!proto.protocol_upgrade.options.gloo.solo.io.ProtocolUpgradeConfig=} opt_value
 * @param {number=} opt_index
 * @return {!proto.protocol_upgrade.options.gloo.solo.io.ProtocolUpgradeConfig}
 */
proto.gloo.solo.io.RouteOptions.prototype.addUpgrades = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 21, opt_value, proto.protocol_upgrade.options.gloo.solo.io.ProtocolUpgradeConfig, opt_index);
};


proto.gloo.solo.io.RouteOptions.prototype.clearUpgradesList = function() {
  this.setUpgradesList([]);
};


/**
 * optional ratelimit.options.gloo.solo.io.IngressRateLimit ratelimit_basic = 13;
 * @return {?proto.ratelimit.options.gloo.solo.io.IngressRateLimit}
 */
proto.gloo.solo.io.RouteOptions.prototype.getRatelimitBasic = function() {
  return /** @type{?proto.ratelimit.options.gloo.solo.io.IngressRateLimit} */ (
    jspb.Message.getWrapperField(this, gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_pb.IngressRateLimit, 13));
};


/** @param {?proto.ratelimit.options.gloo.solo.io.IngressRateLimit|undefined} value */
proto.gloo.solo.io.RouteOptions.prototype.setRatelimitBasic = function(value) {
  jspb.Message.setWrapperField(this, 13, value);
};


proto.gloo.solo.io.RouteOptions.prototype.clearRatelimitBasic = function() {
  this.setRatelimitBasic(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.gloo.solo.io.RouteOptions.prototype.hasRatelimitBasic = function() {
  return jspb.Message.getField(this, 13) != null;
};


/**
 * optional ratelimit.options.gloo.solo.io.RateLimitRouteExtension ratelimit = 14;
 * @return {?proto.ratelimit.options.gloo.solo.io.RateLimitRouteExtension}
 */
proto.gloo.solo.io.RouteOptions.prototype.getRatelimit = function() {
  return /** @type{?proto.ratelimit.options.gloo.solo.io.RateLimitRouteExtension} */ (
    jspb.Message.getWrapperField(this, gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_pb.RateLimitRouteExtension, 14));
};


/** @param {?proto.ratelimit.options.gloo.solo.io.RateLimitRouteExtension|undefined} value */
proto.gloo.solo.io.RouteOptions.prototype.setRatelimit = function(value) {
  jspb.Message.setWrapperField(this, 14, value);
};


proto.gloo.solo.io.RouteOptions.prototype.clearRatelimit = function() {
  this.setRatelimit(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.gloo.solo.io.RouteOptions.prototype.hasRatelimit = function() {
  return jspb.Message.getField(this, 14) != null;
};


/**
 * optional waf.options.gloo.solo.io.Settings waf = 15;
 * @return {?proto.waf.options.gloo.solo.io.Settings}
 */
proto.gloo.solo.io.RouteOptions.prototype.getWaf = function() {
  return /** @type{?proto.waf.options.gloo.solo.io.Settings} */ (
    jspb.Message.getWrapperField(this, gloo_projects_gloo_api_v1_enterprise_options_waf_waf_pb.Settings, 15));
};


/** @param {?proto.waf.options.gloo.solo.io.Settings|undefined} value */
proto.gloo.solo.io.RouteOptions.prototype.setWaf = function(value) {
  jspb.Message.setWrapperField(this, 15, value);
};


proto.gloo.solo.io.RouteOptions.prototype.clearWaf = function() {
  this.setWaf(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.gloo.solo.io.RouteOptions.prototype.hasWaf = function() {
  return jspb.Message.getField(this, 15) != null;
};


/**
 * optional jwt.options.gloo.solo.io.RouteExtension jwt = 16;
 * @return {?proto.jwt.options.gloo.solo.io.RouteExtension}
 */
proto.gloo.solo.io.RouteOptions.prototype.getJwt = function() {
  return /** @type{?proto.jwt.options.gloo.solo.io.RouteExtension} */ (
    jspb.Message.getWrapperField(this, gloo_projects_gloo_api_v1_enterprise_options_jwt_jwt_pb.RouteExtension, 16));
};


/** @param {?proto.jwt.options.gloo.solo.io.RouteExtension|undefined} value */
proto.gloo.solo.io.RouteOptions.prototype.setJwt = function(value) {
  jspb.Message.setWrapperField(this, 16, value);
};


proto.gloo.solo.io.RouteOptions.prototype.clearJwt = function() {
  this.setJwt(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.gloo.solo.io.RouteOptions.prototype.hasJwt = function() {
  return jspb.Message.getField(this, 16) != null;
};


/**
 * optional rbac.options.gloo.solo.io.ExtensionSettings rbac = 17;
 * @return {?proto.rbac.options.gloo.solo.io.ExtensionSettings}
 */
proto.gloo.solo.io.RouteOptions.prototype.getRbac = function() {
  return /** @type{?proto.rbac.options.gloo.solo.io.ExtensionSettings} */ (
    jspb.Message.getWrapperField(this, gloo_projects_gloo_api_v1_enterprise_options_rbac_rbac_pb.ExtensionSettings, 17));
};


/** @param {?proto.rbac.options.gloo.solo.io.ExtensionSettings|undefined} value */
proto.gloo.solo.io.RouteOptions.prototype.setRbac = function(value) {
  jspb.Message.setWrapperField(this, 17, value);
};


proto.gloo.solo.io.RouteOptions.prototype.clearRbac = function() {
  this.setRbac(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.gloo.solo.io.RouteOptions.prototype.hasRbac = function() {
  return jspb.Message.getField(this, 17) != null;
};


/**
 * optional enterprise.gloo.solo.io.ExtAuthExtension extauth = 18;
 * @return {?proto.enterprise.gloo.solo.io.ExtAuthExtension}
 */
proto.gloo.solo.io.RouteOptions.prototype.getExtauth = function() {
  return /** @type{?proto.enterprise.gloo.solo.io.ExtAuthExtension} */ (
    jspb.Message.getWrapperField(this, gloo_projects_gloo_api_v1_enterprise_options_extauth_v1_extauth_pb.ExtAuthExtension, 18));
};


/** @param {?proto.enterprise.gloo.solo.io.ExtAuthExtension|undefined} value */
proto.gloo.solo.io.RouteOptions.prototype.setExtauth = function(value) {
  jspb.Message.setWrapperField(this, 18, value);
};


proto.gloo.solo.io.RouteOptions.prototype.clearExtauth = function() {
  this.setExtauth(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.gloo.solo.io.RouteOptions.prototype.hasExtauth = function() {
  return jspb.Message.getField(this, 18) != null;
};


/**
 * optional dlp.options.gloo.solo.io.Config dlp = 20;
 * @return {?proto.dlp.options.gloo.solo.io.Config}
 */
proto.gloo.solo.io.RouteOptions.prototype.getDlp = function() {
  return /** @type{?proto.dlp.options.gloo.solo.io.Config} */ (
    jspb.Message.getWrapperField(this, gloo_projects_gloo_api_v1_enterprise_options_dlp_dlp_pb.Config, 20));
};


/** @param {?proto.dlp.options.gloo.solo.io.Config|undefined} value */
proto.gloo.solo.io.RouteOptions.prototype.setDlp = function(value) {
  jspb.Message.setWrapperField(this, 20, value);
};


proto.gloo.solo.io.RouteOptions.prototype.clearDlp = function() {
  this.setDlp(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.gloo.solo.io.RouteOptions.prototype.hasDlp = function() {
  return jspb.Message.getField(this, 20) != null;
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
proto.gloo.solo.io.DestinationSpec = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, proto.gloo.solo.io.DestinationSpec.oneofGroups_);
};
goog.inherits(proto.gloo.solo.io.DestinationSpec, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.gloo.solo.io.DestinationSpec.displayName = 'proto.gloo.solo.io.DestinationSpec';
}
/**
 * Oneof group definitions for this message. Each group defines the field
 * numbers belonging to that group. When of these fields' value is set, all
 * other fields in the group are cleared. During deserialization, if multiple
 * fields are encountered for a group, only the last value seen will be kept.
 * @private {!Array<!Array<number>>}
 * @const
 */
proto.gloo.solo.io.DestinationSpec.oneofGroups_ = [[1,2,3,4]];

/**
 * @enum {number}
 */
proto.gloo.solo.io.DestinationSpec.DestinationTypeCase = {
  DESTINATION_TYPE_NOT_SET: 0,
  AWS: 1,
  AZURE: 2,
  REST: 3,
  GRPC: 4
};

/**
 * @return {proto.gloo.solo.io.DestinationSpec.DestinationTypeCase}
 */
proto.gloo.solo.io.DestinationSpec.prototype.getDestinationTypeCase = function() {
  return /** @type {proto.gloo.solo.io.DestinationSpec.DestinationTypeCase} */(jspb.Message.computeOneofCase(this, proto.gloo.solo.io.DestinationSpec.oneofGroups_[0]));
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
proto.gloo.solo.io.DestinationSpec.prototype.toObject = function(opt_includeInstance) {
  return proto.gloo.solo.io.DestinationSpec.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.gloo.solo.io.DestinationSpec} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.gloo.solo.io.DestinationSpec.toObject = function(includeInstance, msg) {
  var f, obj = {
    aws: (f = msg.getAws()) && gloo_projects_gloo_api_v1_options_aws_aws_pb.DestinationSpec.toObject(includeInstance, f),
    azure: (f = msg.getAzure()) && gloo_projects_gloo_api_v1_options_azure_azure_pb.DestinationSpec.toObject(includeInstance, f),
    rest: (f = msg.getRest()) && gloo_projects_gloo_api_v1_options_rest_rest_pb.DestinationSpec.toObject(includeInstance, f),
    grpc: (f = msg.getGrpc()) && gloo_projects_gloo_api_v1_options_grpc_grpc_pb.DestinationSpec.toObject(includeInstance, f)
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
 * @return {!proto.gloo.solo.io.DestinationSpec}
 */
proto.gloo.solo.io.DestinationSpec.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.gloo.solo.io.DestinationSpec;
  return proto.gloo.solo.io.DestinationSpec.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.gloo.solo.io.DestinationSpec} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.gloo.solo.io.DestinationSpec}
 */
proto.gloo.solo.io.DestinationSpec.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new gloo_projects_gloo_api_v1_options_aws_aws_pb.DestinationSpec;
      reader.readMessage(value,gloo_projects_gloo_api_v1_options_aws_aws_pb.DestinationSpec.deserializeBinaryFromReader);
      msg.setAws(value);
      break;
    case 2:
      var value = new gloo_projects_gloo_api_v1_options_azure_azure_pb.DestinationSpec;
      reader.readMessage(value,gloo_projects_gloo_api_v1_options_azure_azure_pb.DestinationSpec.deserializeBinaryFromReader);
      msg.setAzure(value);
      break;
    case 3:
      var value = new gloo_projects_gloo_api_v1_options_rest_rest_pb.DestinationSpec;
      reader.readMessage(value,gloo_projects_gloo_api_v1_options_rest_rest_pb.DestinationSpec.deserializeBinaryFromReader);
      msg.setRest(value);
      break;
    case 4:
      var value = new gloo_projects_gloo_api_v1_options_grpc_grpc_pb.DestinationSpec;
      reader.readMessage(value,gloo_projects_gloo_api_v1_options_grpc_grpc_pb.DestinationSpec.deserializeBinaryFromReader);
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
proto.gloo.solo.io.DestinationSpec.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.gloo.solo.io.DestinationSpec.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.gloo.solo.io.DestinationSpec} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.gloo.solo.io.DestinationSpec.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getAws();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      gloo_projects_gloo_api_v1_options_aws_aws_pb.DestinationSpec.serializeBinaryToWriter
    );
  }
  f = message.getAzure();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      gloo_projects_gloo_api_v1_options_azure_azure_pb.DestinationSpec.serializeBinaryToWriter
    );
  }
  f = message.getRest();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      gloo_projects_gloo_api_v1_options_rest_rest_pb.DestinationSpec.serializeBinaryToWriter
    );
  }
  f = message.getGrpc();
  if (f != null) {
    writer.writeMessage(
      4,
      f,
      gloo_projects_gloo_api_v1_options_grpc_grpc_pb.DestinationSpec.serializeBinaryToWriter
    );
  }
};


/**
 * optional aws.options.gloo.solo.io.DestinationSpec aws = 1;
 * @return {?proto.aws.options.gloo.solo.io.DestinationSpec}
 */
proto.gloo.solo.io.DestinationSpec.prototype.getAws = function() {
  return /** @type{?proto.aws.options.gloo.solo.io.DestinationSpec} */ (
    jspb.Message.getWrapperField(this, gloo_projects_gloo_api_v1_options_aws_aws_pb.DestinationSpec, 1));
};


/** @param {?proto.aws.options.gloo.solo.io.DestinationSpec|undefined} value */
proto.gloo.solo.io.DestinationSpec.prototype.setAws = function(value) {
  jspb.Message.setOneofWrapperField(this, 1, proto.gloo.solo.io.DestinationSpec.oneofGroups_[0], value);
};


proto.gloo.solo.io.DestinationSpec.prototype.clearAws = function() {
  this.setAws(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.gloo.solo.io.DestinationSpec.prototype.hasAws = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional azure.options.gloo.solo.io.DestinationSpec azure = 2;
 * @return {?proto.azure.options.gloo.solo.io.DestinationSpec}
 */
proto.gloo.solo.io.DestinationSpec.prototype.getAzure = function() {
  return /** @type{?proto.azure.options.gloo.solo.io.DestinationSpec} */ (
    jspb.Message.getWrapperField(this, gloo_projects_gloo_api_v1_options_azure_azure_pb.DestinationSpec, 2));
};


/** @param {?proto.azure.options.gloo.solo.io.DestinationSpec|undefined} value */
proto.gloo.solo.io.DestinationSpec.prototype.setAzure = function(value) {
  jspb.Message.setOneofWrapperField(this, 2, proto.gloo.solo.io.DestinationSpec.oneofGroups_[0], value);
};


proto.gloo.solo.io.DestinationSpec.prototype.clearAzure = function() {
  this.setAzure(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.gloo.solo.io.DestinationSpec.prototype.hasAzure = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional rest.options.gloo.solo.io.DestinationSpec rest = 3;
 * @return {?proto.rest.options.gloo.solo.io.DestinationSpec}
 */
proto.gloo.solo.io.DestinationSpec.prototype.getRest = function() {
  return /** @type{?proto.rest.options.gloo.solo.io.DestinationSpec} */ (
    jspb.Message.getWrapperField(this, gloo_projects_gloo_api_v1_options_rest_rest_pb.DestinationSpec, 3));
};


/** @param {?proto.rest.options.gloo.solo.io.DestinationSpec|undefined} value */
proto.gloo.solo.io.DestinationSpec.prototype.setRest = function(value) {
  jspb.Message.setOneofWrapperField(this, 3, proto.gloo.solo.io.DestinationSpec.oneofGroups_[0], value);
};


proto.gloo.solo.io.DestinationSpec.prototype.clearRest = function() {
  this.setRest(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.gloo.solo.io.DestinationSpec.prototype.hasRest = function() {
  return jspb.Message.getField(this, 3) != null;
};


/**
 * optional grpc.options.gloo.solo.io.DestinationSpec grpc = 4;
 * @return {?proto.grpc.options.gloo.solo.io.DestinationSpec}
 */
proto.gloo.solo.io.DestinationSpec.prototype.getGrpc = function() {
  return /** @type{?proto.grpc.options.gloo.solo.io.DestinationSpec} */ (
    jspb.Message.getWrapperField(this, gloo_projects_gloo_api_v1_options_grpc_grpc_pb.DestinationSpec, 4));
};


/** @param {?proto.grpc.options.gloo.solo.io.DestinationSpec|undefined} value */
proto.gloo.solo.io.DestinationSpec.prototype.setGrpc = function(value) {
  jspb.Message.setOneofWrapperField(this, 4, proto.gloo.solo.io.DestinationSpec.oneofGroups_[0], value);
};


proto.gloo.solo.io.DestinationSpec.prototype.clearGrpc = function() {
  this.setGrpc(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.gloo.solo.io.DestinationSpec.prototype.hasGrpc = function() {
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
proto.gloo.solo.io.WeightedDestinationOptions = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.gloo.solo.io.WeightedDestinationOptions, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.gloo.solo.io.WeightedDestinationOptions.displayName = 'proto.gloo.solo.io.WeightedDestinationOptions';
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
proto.gloo.solo.io.WeightedDestinationOptions.prototype.toObject = function(opt_includeInstance) {
  return proto.gloo.solo.io.WeightedDestinationOptions.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.gloo.solo.io.WeightedDestinationOptions} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.gloo.solo.io.WeightedDestinationOptions.toObject = function(includeInstance, msg) {
  var f, obj = {
    headerManipulation: (f = msg.getHeaderManipulation()) && gloo_projects_gloo_api_v1_options_headers_headers_pb.HeaderManipulation.toObject(includeInstance, f),
    transformations: (f = msg.getTransformations()) && gloo_projects_gloo_api_external_envoy_extensions_transformation_transformation_pb.RouteTransformations.toObject(includeInstance, f),
    extensions: (f = msg.getExtensions()) && gloo_projects_gloo_api_v1_extensions_pb.Extensions.toObject(includeInstance, f),
    extauth: (f = msg.getExtauth()) && gloo_projects_gloo_api_v1_enterprise_options_extauth_v1_extauth_pb.ExtAuthExtension.toObject(includeInstance, f)
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
 * @return {!proto.gloo.solo.io.WeightedDestinationOptions}
 */
proto.gloo.solo.io.WeightedDestinationOptions.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.gloo.solo.io.WeightedDestinationOptions;
  return proto.gloo.solo.io.WeightedDestinationOptions.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.gloo.solo.io.WeightedDestinationOptions} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.gloo.solo.io.WeightedDestinationOptions}
 */
proto.gloo.solo.io.WeightedDestinationOptions.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new gloo_projects_gloo_api_v1_options_headers_headers_pb.HeaderManipulation;
      reader.readMessage(value,gloo_projects_gloo_api_v1_options_headers_headers_pb.HeaderManipulation.deserializeBinaryFromReader);
      msg.setHeaderManipulation(value);
      break;
    case 2:
      var value = new gloo_projects_gloo_api_external_envoy_extensions_transformation_transformation_pb.RouteTransformations;
      reader.readMessage(value,gloo_projects_gloo_api_external_envoy_extensions_transformation_transformation_pb.RouteTransformations.deserializeBinaryFromReader);
      msg.setTransformations(value);
      break;
    case 3:
      var value = new gloo_projects_gloo_api_v1_extensions_pb.Extensions;
      reader.readMessage(value,gloo_projects_gloo_api_v1_extensions_pb.Extensions.deserializeBinaryFromReader);
      msg.setExtensions(value);
      break;
    case 4:
      var value = new gloo_projects_gloo_api_v1_enterprise_options_extauth_v1_extauth_pb.ExtAuthExtension;
      reader.readMessage(value,gloo_projects_gloo_api_v1_enterprise_options_extauth_v1_extauth_pb.ExtAuthExtension.deserializeBinaryFromReader);
      msg.setExtauth(value);
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
proto.gloo.solo.io.WeightedDestinationOptions.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.gloo.solo.io.WeightedDestinationOptions.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.gloo.solo.io.WeightedDestinationOptions} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.gloo.solo.io.WeightedDestinationOptions.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getHeaderManipulation();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      gloo_projects_gloo_api_v1_options_headers_headers_pb.HeaderManipulation.serializeBinaryToWriter
    );
  }
  f = message.getTransformations();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      gloo_projects_gloo_api_external_envoy_extensions_transformation_transformation_pb.RouteTransformations.serializeBinaryToWriter
    );
  }
  f = message.getExtensions();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      gloo_projects_gloo_api_v1_extensions_pb.Extensions.serializeBinaryToWriter
    );
  }
  f = message.getExtauth();
  if (f != null) {
    writer.writeMessage(
      4,
      f,
      gloo_projects_gloo_api_v1_enterprise_options_extauth_v1_extauth_pb.ExtAuthExtension.serializeBinaryToWriter
    );
  }
};


/**
 * optional headers.options.gloo.solo.io.HeaderManipulation header_manipulation = 1;
 * @return {?proto.headers.options.gloo.solo.io.HeaderManipulation}
 */
proto.gloo.solo.io.WeightedDestinationOptions.prototype.getHeaderManipulation = function() {
  return /** @type{?proto.headers.options.gloo.solo.io.HeaderManipulation} */ (
    jspb.Message.getWrapperField(this, gloo_projects_gloo_api_v1_options_headers_headers_pb.HeaderManipulation, 1));
};


/** @param {?proto.headers.options.gloo.solo.io.HeaderManipulation|undefined} value */
proto.gloo.solo.io.WeightedDestinationOptions.prototype.setHeaderManipulation = function(value) {
  jspb.Message.setWrapperField(this, 1, value);
};


proto.gloo.solo.io.WeightedDestinationOptions.prototype.clearHeaderManipulation = function() {
  this.setHeaderManipulation(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.gloo.solo.io.WeightedDestinationOptions.prototype.hasHeaderManipulation = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional envoy.api.v2.filter.http.RouteTransformations transformations = 2;
 * @return {?proto.envoy.api.v2.filter.http.RouteTransformations}
 */
proto.gloo.solo.io.WeightedDestinationOptions.prototype.getTransformations = function() {
  return /** @type{?proto.envoy.api.v2.filter.http.RouteTransformations} */ (
    jspb.Message.getWrapperField(this, gloo_projects_gloo_api_external_envoy_extensions_transformation_transformation_pb.RouteTransformations, 2));
};


/** @param {?proto.envoy.api.v2.filter.http.RouteTransformations|undefined} value */
proto.gloo.solo.io.WeightedDestinationOptions.prototype.setTransformations = function(value) {
  jspb.Message.setWrapperField(this, 2, value);
};


proto.gloo.solo.io.WeightedDestinationOptions.prototype.clearTransformations = function() {
  this.setTransformations(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.gloo.solo.io.WeightedDestinationOptions.prototype.hasTransformations = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional Extensions extensions = 3;
 * @return {?proto.gloo.solo.io.Extensions}
 */
proto.gloo.solo.io.WeightedDestinationOptions.prototype.getExtensions = function() {
  return /** @type{?proto.gloo.solo.io.Extensions} */ (
    jspb.Message.getWrapperField(this, gloo_projects_gloo_api_v1_extensions_pb.Extensions, 3));
};


/** @param {?proto.gloo.solo.io.Extensions|undefined} value */
proto.gloo.solo.io.WeightedDestinationOptions.prototype.setExtensions = function(value) {
  jspb.Message.setWrapperField(this, 3, value);
};


proto.gloo.solo.io.WeightedDestinationOptions.prototype.clearExtensions = function() {
  this.setExtensions(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.gloo.solo.io.WeightedDestinationOptions.prototype.hasExtensions = function() {
  return jspb.Message.getField(this, 3) != null;
};


/**
 * optional enterprise.gloo.solo.io.ExtAuthExtension extauth = 4;
 * @return {?proto.enterprise.gloo.solo.io.ExtAuthExtension}
 */
proto.gloo.solo.io.WeightedDestinationOptions.prototype.getExtauth = function() {
  return /** @type{?proto.enterprise.gloo.solo.io.ExtAuthExtension} */ (
    jspb.Message.getWrapperField(this, gloo_projects_gloo_api_v1_enterprise_options_extauth_v1_extauth_pb.ExtAuthExtension, 4));
};


/** @param {?proto.enterprise.gloo.solo.io.ExtAuthExtension|undefined} value */
proto.gloo.solo.io.WeightedDestinationOptions.prototype.setExtauth = function(value) {
  jspb.Message.setWrapperField(this, 4, value);
};


proto.gloo.solo.io.WeightedDestinationOptions.prototype.clearExtauth = function() {
  this.setExtauth(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.gloo.solo.io.WeightedDestinationOptions.prototype.hasExtauth = function() {
  return jspb.Message.getField(this, 4) != null;
};


goog.object.extend(exports, proto.gloo.solo.io);
