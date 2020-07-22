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

var envoy_config_core_v3_base_pb = require('../../../../envoy/config/core/v3/base_pb.js');
var envoy_config_core_v3_extension_pb = require('../../../../envoy/config/core/v3/extension_pb.js');
var envoy_config_core_v3_proxy_protocol_pb = require('../../../../envoy/config/core/v3/proxy_protocol_pb.js');
var envoy_type_matcher_v3_regex_pb = require('../../../../envoy/type/matcher/v3/regex_pb.js');
var envoy_type_matcher_v3_string_pb = require('../../../../envoy/type/matcher/v3/string_pb.js');
var envoy_type_metadata_v3_metadata_pb = require('../../../../envoy/type/metadata/v3/metadata_pb.js');
var envoy_type_tracing_v3_custom_tag_pb = require('../../../../envoy/type/tracing/v3/custom_tag_pb.js');
var envoy_type_v3_percent_pb = require('../../../../envoy/type/v3/percent_pb.js');
var envoy_type_v3_range_pb = require('../../../../envoy/type/v3/range_pb.js');
var google_protobuf_any_pb = require('google-protobuf/google/protobuf/any_pb.js');
var google_protobuf_duration_pb = require('google-protobuf/google/protobuf/duration_pb.js');
var google_protobuf_struct_pb = require('google-protobuf/google/protobuf/struct_pb.js');
var google_protobuf_wrappers_pb = require('google-protobuf/google/protobuf/wrappers_pb.js');
var envoy_annotations_deprecation_pb = require('../../../../envoy/annotations/deprecation_pb.js');
var udpa_annotations_migrate_pb = require('../../../../udpa/annotations/migrate_pb.js');
var udpa_annotations_status_pb = require('../../../../udpa/annotations/status_pb.js');
var udpa_annotations_versioning_pb = require('../../../../udpa/annotations/versioning_pb.js');
var validate_validate_pb = require('../../../../validate/validate_pb.js');
var gogoproto_gogo_pb = require('../../../../gogoproto/gogo_pb.js');
goog.exportSymbol('proto.envoy.config.route.v3.CorsPolicy', null, global);
goog.exportSymbol('proto.envoy.config.route.v3.Decorator', null, global);
goog.exportSymbol('proto.envoy.config.route.v3.DirectResponseAction', null, global);
goog.exportSymbol('proto.envoy.config.route.v3.FilterAction', null, global);
goog.exportSymbol('proto.envoy.config.route.v3.HeaderMatcher', null, global);
goog.exportSymbol('proto.envoy.config.route.v3.HedgePolicy', null, global);
goog.exportSymbol('proto.envoy.config.route.v3.InternalRedirectPolicy', null, global);
goog.exportSymbol('proto.envoy.config.route.v3.QueryParameterMatcher', null, global);
goog.exportSymbol('proto.envoy.config.route.v3.RateLimit', null, global);
goog.exportSymbol('proto.envoy.config.route.v3.RateLimit.Action', null, global);
goog.exportSymbol('proto.envoy.config.route.v3.RateLimit.Action.DestinationCluster', null, global);
goog.exportSymbol('proto.envoy.config.route.v3.RateLimit.Action.DynamicMetaData', null, global);
goog.exportSymbol('proto.envoy.config.route.v3.RateLimit.Action.GenericKey', null, global);
goog.exportSymbol('proto.envoy.config.route.v3.RateLimit.Action.HeaderValueMatch', null, global);
goog.exportSymbol('proto.envoy.config.route.v3.RateLimit.Action.RemoteAddress', null, global);
goog.exportSymbol('proto.envoy.config.route.v3.RateLimit.Action.RequestHeaders', null, global);
goog.exportSymbol('proto.envoy.config.route.v3.RateLimit.Action.SourceCluster', null, global);
goog.exportSymbol('proto.envoy.config.route.v3.RateLimit.Override', null, global);
goog.exportSymbol('proto.envoy.config.route.v3.RateLimit.Override.DynamicMetadata', null, global);
goog.exportSymbol('proto.envoy.config.route.v3.RedirectAction', null, global);
goog.exportSymbol('proto.envoy.config.route.v3.RedirectAction.RedirectResponseCode', null, global);
goog.exportSymbol('proto.envoy.config.route.v3.RetryPolicy', null, global);
goog.exportSymbol('proto.envoy.config.route.v3.RetryPolicy.RetryBackOff', null, global);
goog.exportSymbol('proto.envoy.config.route.v3.RetryPolicy.RetryHostPredicate', null, global);
goog.exportSymbol('proto.envoy.config.route.v3.RetryPolicy.RetryPriority', null, global);
goog.exportSymbol('proto.envoy.config.route.v3.Route', null, global);
goog.exportSymbol('proto.envoy.config.route.v3.RouteAction', null, global);
goog.exportSymbol('proto.envoy.config.route.v3.RouteAction.ClusterNotFoundResponseCode', null, global);
goog.exportSymbol('proto.envoy.config.route.v3.RouteAction.HashPolicy', null, global);
goog.exportSymbol('proto.envoy.config.route.v3.RouteAction.HashPolicy.ConnectionProperties', null, global);
goog.exportSymbol('proto.envoy.config.route.v3.RouteAction.HashPolicy.Cookie', null, global);
goog.exportSymbol('proto.envoy.config.route.v3.RouteAction.HashPolicy.FilterState', null, global);
goog.exportSymbol('proto.envoy.config.route.v3.RouteAction.HashPolicy.Header', null, global);
goog.exportSymbol('proto.envoy.config.route.v3.RouteAction.HashPolicy.QueryParameter', null, global);
goog.exportSymbol('proto.envoy.config.route.v3.RouteAction.InternalRedirectAction', null, global);
goog.exportSymbol('proto.envoy.config.route.v3.RouteAction.RequestMirrorPolicy', null, global);
goog.exportSymbol('proto.envoy.config.route.v3.RouteAction.UpgradeConfig', null, global);
goog.exportSymbol('proto.envoy.config.route.v3.RouteAction.UpgradeConfig.ConnectConfig', null, global);
goog.exportSymbol('proto.envoy.config.route.v3.RouteMatch', null, global);
goog.exportSymbol('proto.envoy.config.route.v3.RouteMatch.ConnectMatcher', null, global);
goog.exportSymbol('proto.envoy.config.route.v3.RouteMatch.GrpcRouteMatchOptions', null, global);
goog.exportSymbol('proto.envoy.config.route.v3.RouteMatch.TlsContextMatchOptions', null, global);
goog.exportSymbol('proto.envoy.config.route.v3.Tracing', null, global);
goog.exportSymbol('proto.envoy.config.route.v3.VirtualCluster', null, global);
goog.exportSymbol('proto.envoy.config.route.v3.VirtualHost', null, global);
goog.exportSymbol('proto.envoy.config.route.v3.VirtualHost.TlsRequirementType', null, global);
goog.exportSymbol('proto.envoy.config.route.v3.WeightedCluster', null, global);
goog.exportSymbol('proto.envoy.config.route.v3.WeightedCluster.ClusterWeight', null, global);

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
proto.envoy.config.route.v3.VirtualHost = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.envoy.config.route.v3.VirtualHost.repeatedFields_, null);
};
goog.inherits(proto.envoy.config.route.v3.VirtualHost, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.envoy.config.route.v3.VirtualHost.displayName = 'proto.envoy.config.route.v3.VirtualHost';
}
/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.envoy.config.route.v3.VirtualHost.repeatedFields_ = [2,3,5,6,7,13,10,11];



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
proto.envoy.config.route.v3.VirtualHost.prototype.toObject = function(opt_includeInstance) {
  return proto.envoy.config.route.v3.VirtualHost.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.envoy.config.route.v3.VirtualHost} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.VirtualHost.toObject = function(includeInstance, msg) {
  var f, obj = {
    name: jspb.Message.getFieldWithDefault(msg, 1, ""),
    domainsList: jspb.Message.getRepeatedField(msg, 2),
    routesList: jspb.Message.toObjectList(msg.getRoutesList(),
    proto.envoy.config.route.v3.Route.toObject, includeInstance),
    requireTls: jspb.Message.getFieldWithDefault(msg, 4, 0),
    virtualClustersList: jspb.Message.toObjectList(msg.getVirtualClustersList(),
    proto.envoy.config.route.v3.VirtualCluster.toObject, includeInstance),
    rateLimitsList: jspb.Message.toObjectList(msg.getRateLimitsList(),
    proto.envoy.config.route.v3.RateLimit.toObject, includeInstance),
    requestHeadersToAddList: jspb.Message.toObjectList(msg.getRequestHeadersToAddList(),
    envoy_config_core_v3_base_pb.HeaderValueOption.toObject, includeInstance),
    requestHeadersToRemoveList: jspb.Message.getRepeatedField(msg, 13),
    responseHeadersToAddList: jspb.Message.toObjectList(msg.getResponseHeadersToAddList(),
    envoy_config_core_v3_base_pb.HeaderValueOption.toObject, includeInstance),
    responseHeadersToRemoveList: jspb.Message.getRepeatedField(msg, 11),
    cors: (f = msg.getCors()) && proto.envoy.config.route.v3.CorsPolicy.toObject(includeInstance, f),
    typedPerFilterConfigMap: (f = msg.getTypedPerFilterConfigMap()) ? f.toObject(includeInstance, proto.google.protobuf.Any.toObject) : [],
    includeRequestAttemptCount: jspb.Message.getFieldWithDefault(msg, 14, false),
    includeAttemptCountInResponse: jspb.Message.getFieldWithDefault(msg, 19, false),
    retryPolicy: (f = msg.getRetryPolicy()) && proto.envoy.config.route.v3.RetryPolicy.toObject(includeInstance, f),
    retryPolicyTypedConfig: (f = msg.getRetryPolicyTypedConfig()) && google_protobuf_any_pb.Any.toObject(includeInstance, f),
    hedgePolicy: (f = msg.getHedgePolicy()) && proto.envoy.config.route.v3.HedgePolicy.toObject(includeInstance, f),
    perRequestBufferLimitBytes: (f = msg.getPerRequestBufferLimitBytes()) && google_protobuf_wrappers_pb.UInt32Value.toObject(includeInstance, f)
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
 * @return {!proto.envoy.config.route.v3.VirtualHost}
 */
proto.envoy.config.route.v3.VirtualHost.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.envoy.config.route.v3.VirtualHost;
  return proto.envoy.config.route.v3.VirtualHost.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.envoy.config.route.v3.VirtualHost} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.envoy.config.route.v3.VirtualHost}
 */
proto.envoy.config.route.v3.VirtualHost.deserializeBinaryFromReader = function(msg, reader) {
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
      msg.addDomains(value);
      break;
    case 3:
      var value = new proto.envoy.config.route.v3.Route;
      reader.readMessage(value,proto.envoy.config.route.v3.Route.deserializeBinaryFromReader);
      msg.addRoutes(value);
      break;
    case 4:
      var value = /** @type {!proto.envoy.config.route.v3.VirtualHost.TlsRequirementType} */ (reader.readEnum());
      msg.setRequireTls(value);
      break;
    case 5:
      var value = new proto.envoy.config.route.v3.VirtualCluster;
      reader.readMessage(value,proto.envoy.config.route.v3.VirtualCluster.deserializeBinaryFromReader);
      msg.addVirtualClusters(value);
      break;
    case 6:
      var value = new proto.envoy.config.route.v3.RateLimit;
      reader.readMessage(value,proto.envoy.config.route.v3.RateLimit.deserializeBinaryFromReader);
      msg.addRateLimits(value);
      break;
    case 7:
      var value = new envoy_config_core_v3_base_pb.HeaderValueOption;
      reader.readMessage(value,envoy_config_core_v3_base_pb.HeaderValueOption.deserializeBinaryFromReader);
      msg.addRequestHeadersToAdd(value);
      break;
    case 13:
      var value = /** @type {string} */ (reader.readString());
      msg.addRequestHeadersToRemove(value);
      break;
    case 10:
      var value = new envoy_config_core_v3_base_pb.HeaderValueOption;
      reader.readMessage(value,envoy_config_core_v3_base_pb.HeaderValueOption.deserializeBinaryFromReader);
      msg.addResponseHeadersToAdd(value);
      break;
    case 11:
      var value = /** @type {string} */ (reader.readString());
      msg.addResponseHeadersToRemove(value);
      break;
    case 8:
      var value = new proto.envoy.config.route.v3.CorsPolicy;
      reader.readMessage(value,proto.envoy.config.route.v3.CorsPolicy.deserializeBinaryFromReader);
      msg.setCors(value);
      break;
    case 15:
      var value = msg.getTypedPerFilterConfigMap();
      reader.readMessage(value, function(message, reader) {
        jspb.Map.deserializeBinary(message, reader, jspb.BinaryReader.prototype.readString, jspb.BinaryReader.prototype.readMessage, proto.google.protobuf.Any.deserializeBinaryFromReader, "");
         });
      break;
    case 14:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setIncludeRequestAttemptCount(value);
      break;
    case 19:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setIncludeAttemptCountInResponse(value);
      break;
    case 16:
      var value = new proto.envoy.config.route.v3.RetryPolicy;
      reader.readMessage(value,proto.envoy.config.route.v3.RetryPolicy.deserializeBinaryFromReader);
      msg.setRetryPolicy(value);
      break;
    case 20:
      var value = new google_protobuf_any_pb.Any;
      reader.readMessage(value,google_protobuf_any_pb.Any.deserializeBinaryFromReader);
      msg.setRetryPolicyTypedConfig(value);
      break;
    case 17:
      var value = new proto.envoy.config.route.v3.HedgePolicy;
      reader.readMessage(value,proto.envoy.config.route.v3.HedgePolicy.deserializeBinaryFromReader);
      msg.setHedgePolicy(value);
      break;
    case 18:
      var value = new google_protobuf_wrappers_pb.UInt32Value;
      reader.readMessage(value,google_protobuf_wrappers_pb.UInt32Value.deserializeBinaryFromReader);
      msg.setPerRequestBufferLimitBytes(value);
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
proto.envoy.config.route.v3.VirtualHost.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.envoy.config.route.v3.VirtualHost.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.envoy.config.route.v3.VirtualHost} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.VirtualHost.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getName();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getDomainsList();
  if (f.length > 0) {
    writer.writeRepeatedString(
      2,
      f
    );
  }
  f = message.getRoutesList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      3,
      f,
      proto.envoy.config.route.v3.Route.serializeBinaryToWriter
    );
  }
  f = message.getRequireTls();
  if (f !== 0.0) {
    writer.writeEnum(
      4,
      f
    );
  }
  f = message.getVirtualClustersList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      5,
      f,
      proto.envoy.config.route.v3.VirtualCluster.serializeBinaryToWriter
    );
  }
  f = message.getRateLimitsList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      6,
      f,
      proto.envoy.config.route.v3.RateLimit.serializeBinaryToWriter
    );
  }
  f = message.getRequestHeadersToAddList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      7,
      f,
      envoy_config_core_v3_base_pb.HeaderValueOption.serializeBinaryToWriter
    );
  }
  f = message.getRequestHeadersToRemoveList();
  if (f.length > 0) {
    writer.writeRepeatedString(
      13,
      f
    );
  }
  f = message.getResponseHeadersToAddList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      10,
      f,
      envoy_config_core_v3_base_pb.HeaderValueOption.serializeBinaryToWriter
    );
  }
  f = message.getResponseHeadersToRemoveList();
  if (f.length > 0) {
    writer.writeRepeatedString(
      11,
      f
    );
  }
  f = message.getCors();
  if (f != null) {
    writer.writeMessage(
      8,
      f,
      proto.envoy.config.route.v3.CorsPolicy.serializeBinaryToWriter
    );
  }
  f = message.getTypedPerFilterConfigMap(true);
  if (f && f.getLength() > 0) {
    f.serializeBinary(15, writer, jspb.BinaryWriter.prototype.writeString, jspb.BinaryWriter.prototype.writeMessage, proto.google.protobuf.Any.serializeBinaryToWriter);
  }
  f = message.getIncludeRequestAttemptCount();
  if (f) {
    writer.writeBool(
      14,
      f
    );
  }
  f = message.getIncludeAttemptCountInResponse();
  if (f) {
    writer.writeBool(
      19,
      f
    );
  }
  f = message.getRetryPolicy();
  if (f != null) {
    writer.writeMessage(
      16,
      f,
      proto.envoy.config.route.v3.RetryPolicy.serializeBinaryToWriter
    );
  }
  f = message.getRetryPolicyTypedConfig();
  if (f != null) {
    writer.writeMessage(
      20,
      f,
      google_protobuf_any_pb.Any.serializeBinaryToWriter
    );
  }
  f = message.getHedgePolicy();
  if (f != null) {
    writer.writeMessage(
      17,
      f,
      proto.envoy.config.route.v3.HedgePolicy.serializeBinaryToWriter
    );
  }
  f = message.getPerRequestBufferLimitBytes();
  if (f != null) {
    writer.writeMessage(
      18,
      f,
      google_protobuf_wrappers_pb.UInt32Value.serializeBinaryToWriter
    );
  }
};


/**
 * @enum {number}
 */
proto.envoy.config.route.v3.VirtualHost.TlsRequirementType = {
  NONE: 0,
  EXTERNAL_ONLY: 1,
  ALL: 2
};

/**
 * optional string name = 1;
 * @return {string}
 */
proto.envoy.config.route.v3.VirtualHost.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.envoy.config.route.v3.VirtualHost.prototype.setName = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * repeated string domains = 2;
 * @return {!Array<string>}
 */
proto.envoy.config.route.v3.VirtualHost.prototype.getDomainsList = function() {
  return /** @type {!Array<string>} */ (jspb.Message.getRepeatedField(this, 2));
};


/** @param {!Array<string>} value */
proto.envoy.config.route.v3.VirtualHost.prototype.setDomainsList = function(value) {
  jspb.Message.setField(this, 2, value || []);
};


/**
 * @param {!string} value
 * @param {number=} opt_index
 */
proto.envoy.config.route.v3.VirtualHost.prototype.addDomains = function(value, opt_index) {
  jspb.Message.addToRepeatedField(this, 2, value, opt_index);
};


proto.envoy.config.route.v3.VirtualHost.prototype.clearDomainsList = function() {
  this.setDomainsList([]);
};


/**
 * repeated Route routes = 3;
 * @return {!Array<!proto.envoy.config.route.v3.Route>}
 */
proto.envoy.config.route.v3.VirtualHost.prototype.getRoutesList = function() {
  return /** @type{!Array<!proto.envoy.config.route.v3.Route>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.envoy.config.route.v3.Route, 3));
};


/** @param {!Array<!proto.envoy.config.route.v3.Route>} value */
proto.envoy.config.route.v3.VirtualHost.prototype.setRoutesList = function(value) {
  jspb.Message.setRepeatedWrapperField(this, 3, value);
};


/**
 * @param {!proto.envoy.config.route.v3.Route=} opt_value
 * @param {number=} opt_index
 * @return {!proto.envoy.config.route.v3.Route}
 */
proto.envoy.config.route.v3.VirtualHost.prototype.addRoutes = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 3, opt_value, proto.envoy.config.route.v3.Route, opt_index);
};


proto.envoy.config.route.v3.VirtualHost.prototype.clearRoutesList = function() {
  this.setRoutesList([]);
};


/**
 * optional TlsRequirementType require_tls = 4;
 * @return {!proto.envoy.config.route.v3.VirtualHost.TlsRequirementType}
 */
proto.envoy.config.route.v3.VirtualHost.prototype.getRequireTls = function() {
  return /** @type {!proto.envoy.config.route.v3.VirtualHost.TlsRequirementType} */ (jspb.Message.getFieldWithDefault(this, 4, 0));
};


/** @param {!proto.envoy.config.route.v3.VirtualHost.TlsRequirementType} value */
proto.envoy.config.route.v3.VirtualHost.prototype.setRequireTls = function(value) {
  jspb.Message.setProto3EnumField(this, 4, value);
};


/**
 * repeated VirtualCluster virtual_clusters = 5;
 * @return {!Array<!proto.envoy.config.route.v3.VirtualCluster>}
 */
proto.envoy.config.route.v3.VirtualHost.prototype.getVirtualClustersList = function() {
  return /** @type{!Array<!proto.envoy.config.route.v3.VirtualCluster>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.envoy.config.route.v3.VirtualCluster, 5));
};


/** @param {!Array<!proto.envoy.config.route.v3.VirtualCluster>} value */
proto.envoy.config.route.v3.VirtualHost.prototype.setVirtualClustersList = function(value) {
  jspb.Message.setRepeatedWrapperField(this, 5, value);
};


/**
 * @param {!proto.envoy.config.route.v3.VirtualCluster=} opt_value
 * @param {number=} opt_index
 * @return {!proto.envoy.config.route.v3.VirtualCluster}
 */
proto.envoy.config.route.v3.VirtualHost.prototype.addVirtualClusters = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 5, opt_value, proto.envoy.config.route.v3.VirtualCluster, opt_index);
};


proto.envoy.config.route.v3.VirtualHost.prototype.clearVirtualClustersList = function() {
  this.setVirtualClustersList([]);
};


/**
 * repeated RateLimit rate_limits = 6;
 * @return {!Array<!proto.envoy.config.route.v3.RateLimit>}
 */
proto.envoy.config.route.v3.VirtualHost.prototype.getRateLimitsList = function() {
  return /** @type{!Array<!proto.envoy.config.route.v3.RateLimit>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.envoy.config.route.v3.RateLimit, 6));
};


/** @param {!Array<!proto.envoy.config.route.v3.RateLimit>} value */
proto.envoy.config.route.v3.VirtualHost.prototype.setRateLimitsList = function(value) {
  jspb.Message.setRepeatedWrapperField(this, 6, value);
};


/**
 * @param {!proto.envoy.config.route.v3.RateLimit=} opt_value
 * @param {number=} opt_index
 * @return {!proto.envoy.config.route.v3.RateLimit}
 */
proto.envoy.config.route.v3.VirtualHost.prototype.addRateLimits = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 6, opt_value, proto.envoy.config.route.v3.RateLimit, opt_index);
};


proto.envoy.config.route.v3.VirtualHost.prototype.clearRateLimitsList = function() {
  this.setRateLimitsList([]);
};


/**
 * repeated envoy.config.core.v3.HeaderValueOption request_headers_to_add = 7;
 * @return {!Array<!proto.envoy.config.core.v3.HeaderValueOption>}
 */
proto.envoy.config.route.v3.VirtualHost.prototype.getRequestHeadersToAddList = function() {
  return /** @type{!Array<!proto.envoy.config.core.v3.HeaderValueOption>} */ (
    jspb.Message.getRepeatedWrapperField(this, envoy_config_core_v3_base_pb.HeaderValueOption, 7));
};


/** @param {!Array<!proto.envoy.config.core.v3.HeaderValueOption>} value */
proto.envoy.config.route.v3.VirtualHost.prototype.setRequestHeadersToAddList = function(value) {
  jspb.Message.setRepeatedWrapperField(this, 7, value);
};


/**
 * @param {!proto.envoy.config.core.v3.HeaderValueOption=} opt_value
 * @param {number=} opt_index
 * @return {!proto.envoy.config.core.v3.HeaderValueOption}
 */
proto.envoy.config.route.v3.VirtualHost.prototype.addRequestHeadersToAdd = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 7, opt_value, proto.envoy.config.core.v3.HeaderValueOption, opt_index);
};


proto.envoy.config.route.v3.VirtualHost.prototype.clearRequestHeadersToAddList = function() {
  this.setRequestHeadersToAddList([]);
};


/**
 * repeated string request_headers_to_remove = 13;
 * @return {!Array<string>}
 */
proto.envoy.config.route.v3.VirtualHost.prototype.getRequestHeadersToRemoveList = function() {
  return /** @type {!Array<string>} */ (jspb.Message.getRepeatedField(this, 13));
};


/** @param {!Array<string>} value */
proto.envoy.config.route.v3.VirtualHost.prototype.setRequestHeadersToRemoveList = function(value) {
  jspb.Message.setField(this, 13, value || []);
};


/**
 * @param {!string} value
 * @param {number=} opt_index
 */
proto.envoy.config.route.v3.VirtualHost.prototype.addRequestHeadersToRemove = function(value, opt_index) {
  jspb.Message.addToRepeatedField(this, 13, value, opt_index);
};


proto.envoy.config.route.v3.VirtualHost.prototype.clearRequestHeadersToRemoveList = function() {
  this.setRequestHeadersToRemoveList([]);
};


/**
 * repeated envoy.config.core.v3.HeaderValueOption response_headers_to_add = 10;
 * @return {!Array<!proto.envoy.config.core.v3.HeaderValueOption>}
 */
proto.envoy.config.route.v3.VirtualHost.prototype.getResponseHeadersToAddList = function() {
  return /** @type{!Array<!proto.envoy.config.core.v3.HeaderValueOption>} */ (
    jspb.Message.getRepeatedWrapperField(this, envoy_config_core_v3_base_pb.HeaderValueOption, 10));
};


/** @param {!Array<!proto.envoy.config.core.v3.HeaderValueOption>} value */
proto.envoy.config.route.v3.VirtualHost.prototype.setResponseHeadersToAddList = function(value) {
  jspb.Message.setRepeatedWrapperField(this, 10, value);
};


/**
 * @param {!proto.envoy.config.core.v3.HeaderValueOption=} opt_value
 * @param {number=} opt_index
 * @return {!proto.envoy.config.core.v3.HeaderValueOption}
 */
proto.envoy.config.route.v3.VirtualHost.prototype.addResponseHeadersToAdd = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 10, opt_value, proto.envoy.config.core.v3.HeaderValueOption, opt_index);
};


proto.envoy.config.route.v3.VirtualHost.prototype.clearResponseHeadersToAddList = function() {
  this.setResponseHeadersToAddList([]);
};


/**
 * repeated string response_headers_to_remove = 11;
 * @return {!Array<string>}
 */
proto.envoy.config.route.v3.VirtualHost.prototype.getResponseHeadersToRemoveList = function() {
  return /** @type {!Array<string>} */ (jspb.Message.getRepeatedField(this, 11));
};


/** @param {!Array<string>} value */
proto.envoy.config.route.v3.VirtualHost.prototype.setResponseHeadersToRemoveList = function(value) {
  jspb.Message.setField(this, 11, value || []);
};


/**
 * @param {!string} value
 * @param {number=} opt_index
 */
proto.envoy.config.route.v3.VirtualHost.prototype.addResponseHeadersToRemove = function(value, opt_index) {
  jspb.Message.addToRepeatedField(this, 11, value, opt_index);
};


proto.envoy.config.route.v3.VirtualHost.prototype.clearResponseHeadersToRemoveList = function() {
  this.setResponseHeadersToRemoveList([]);
};


/**
 * optional CorsPolicy cors = 8;
 * @return {?proto.envoy.config.route.v3.CorsPolicy}
 */
proto.envoy.config.route.v3.VirtualHost.prototype.getCors = function() {
  return /** @type{?proto.envoy.config.route.v3.CorsPolicy} */ (
    jspb.Message.getWrapperField(this, proto.envoy.config.route.v3.CorsPolicy, 8));
};


/** @param {?proto.envoy.config.route.v3.CorsPolicy|undefined} value */
proto.envoy.config.route.v3.VirtualHost.prototype.setCors = function(value) {
  jspb.Message.setWrapperField(this, 8, value);
};


proto.envoy.config.route.v3.VirtualHost.prototype.clearCors = function() {
  this.setCors(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.VirtualHost.prototype.hasCors = function() {
  return jspb.Message.getField(this, 8) != null;
};


/**
 * map<string, google.protobuf.Any> typed_per_filter_config = 15;
 * @param {boolean=} opt_noLazyCreate Do not create the map if
 * empty, instead returning `undefined`
 * @return {!jspb.Map<string,!proto.google.protobuf.Any>}
 */
proto.envoy.config.route.v3.VirtualHost.prototype.getTypedPerFilterConfigMap = function(opt_noLazyCreate) {
  return /** @type {!jspb.Map<string,!proto.google.protobuf.Any>} */ (
      jspb.Message.getMapField(this, 15, opt_noLazyCreate,
      proto.google.protobuf.Any));
};


proto.envoy.config.route.v3.VirtualHost.prototype.clearTypedPerFilterConfigMap = function() {
  this.getTypedPerFilterConfigMap().clear();
};


/**
 * optional bool include_request_attempt_count = 14;
 * Note that Boolean fields may be set to 0/1 when serialized from a Java server.
 * You should avoid comparisons like {@code val === true/false} in those cases.
 * @return {boolean}
 */
proto.envoy.config.route.v3.VirtualHost.prototype.getIncludeRequestAttemptCount = function() {
  return /** @type {boolean} */ (jspb.Message.getFieldWithDefault(this, 14, false));
};


/** @param {boolean} value */
proto.envoy.config.route.v3.VirtualHost.prototype.setIncludeRequestAttemptCount = function(value) {
  jspb.Message.setProto3BooleanField(this, 14, value);
};


/**
 * optional bool include_attempt_count_in_response = 19;
 * Note that Boolean fields may be set to 0/1 when serialized from a Java server.
 * You should avoid comparisons like {@code val === true/false} in those cases.
 * @return {boolean}
 */
proto.envoy.config.route.v3.VirtualHost.prototype.getIncludeAttemptCountInResponse = function() {
  return /** @type {boolean} */ (jspb.Message.getFieldWithDefault(this, 19, false));
};


/** @param {boolean} value */
proto.envoy.config.route.v3.VirtualHost.prototype.setIncludeAttemptCountInResponse = function(value) {
  jspb.Message.setProto3BooleanField(this, 19, value);
};


/**
 * optional RetryPolicy retry_policy = 16;
 * @return {?proto.envoy.config.route.v3.RetryPolicy}
 */
proto.envoy.config.route.v3.VirtualHost.prototype.getRetryPolicy = function() {
  return /** @type{?proto.envoy.config.route.v3.RetryPolicy} */ (
    jspb.Message.getWrapperField(this, proto.envoy.config.route.v3.RetryPolicy, 16));
};


/** @param {?proto.envoy.config.route.v3.RetryPolicy|undefined} value */
proto.envoy.config.route.v3.VirtualHost.prototype.setRetryPolicy = function(value) {
  jspb.Message.setWrapperField(this, 16, value);
};


proto.envoy.config.route.v3.VirtualHost.prototype.clearRetryPolicy = function() {
  this.setRetryPolicy(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.VirtualHost.prototype.hasRetryPolicy = function() {
  return jspb.Message.getField(this, 16) != null;
};


/**
 * optional google.protobuf.Any retry_policy_typed_config = 20;
 * @return {?proto.google.protobuf.Any}
 */
proto.envoy.config.route.v3.VirtualHost.prototype.getRetryPolicyTypedConfig = function() {
  return /** @type{?proto.google.protobuf.Any} */ (
    jspb.Message.getWrapperField(this, google_protobuf_any_pb.Any, 20));
};


/** @param {?proto.google.protobuf.Any|undefined} value */
proto.envoy.config.route.v3.VirtualHost.prototype.setRetryPolicyTypedConfig = function(value) {
  jspb.Message.setWrapperField(this, 20, value);
};


proto.envoy.config.route.v3.VirtualHost.prototype.clearRetryPolicyTypedConfig = function() {
  this.setRetryPolicyTypedConfig(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.VirtualHost.prototype.hasRetryPolicyTypedConfig = function() {
  return jspb.Message.getField(this, 20) != null;
};


/**
 * optional HedgePolicy hedge_policy = 17;
 * @return {?proto.envoy.config.route.v3.HedgePolicy}
 */
proto.envoy.config.route.v3.VirtualHost.prototype.getHedgePolicy = function() {
  return /** @type{?proto.envoy.config.route.v3.HedgePolicy} */ (
    jspb.Message.getWrapperField(this, proto.envoy.config.route.v3.HedgePolicy, 17));
};


/** @param {?proto.envoy.config.route.v3.HedgePolicy|undefined} value */
proto.envoy.config.route.v3.VirtualHost.prototype.setHedgePolicy = function(value) {
  jspb.Message.setWrapperField(this, 17, value);
};


proto.envoy.config.route.v3.VirtualHost.prototype.clearHedgePolicy = function() {
  this.setHedgePolicy(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.VirtualHost.prototype.hasHedgePolicy = function() {
  return jspb.Message.getField(this, 17) != null;
};


/**
 * optional google.protobuf.UInt32Value per_request_buffer_limit_bytes = 18;
 * @return {?proto.google.protobuf.UInt32Value}
 */
proto.envoy.config.route.v3.VirtualHost.prototype.getPerRequestBufferLimitBytes = function() {
  return /** @type{?proto.google.protobuf.UInt32Value} */ (
    jspb.Message.getWrapperField(this, google_protobuf_wrappers_pb.UInt32Value, 18));
};


/** @param {?proto.google.protobuf.UInt32Value|undefined} value */
proto.envoy.config.route.v3.VirtualHost.prototype.setPerRequestBufferLimitBytes = function(value) {
  jspb.Message.setWrapperField(this, 18, value);
};


proto.envoy.config.route.v3.VirtualHost.prototype.clearPerRequestBufferLimitBytes = function() {
  this.setPerRequestBufferLimitBytes(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.VirtualHost.prototype.hasPerRequestBufferLimitBytes = function() {
  return jspb.Message.getField(this, 18) != null;
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
proto.envoy.config.route.v3.FilterAction = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.envoy.config.route.v3.FilterAction, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.envoy.config.route.v3.FilterAction.displayName = 'proto.envoy.config.route.v3.FilterAction';
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
proto.envoy.config.route.v3.FilterAction.prototype.toObject = function(opt_includeInstance) {
  return proto.envoy.config.route.v3.FilterAction.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.envoy.config.route.v3.FilterAction} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.FilterAction.toObject = function(includeInstance, msg) {
  var f, obj = {
    action: (f = msg.getAction()) && google_protobuf_any_pb.Any.toObject(includeInstance, f)
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
 * @return {!proto.envoy.config.route.v3.FilterAction}
 */
proto.envoy.config.route.v3.FilterAction.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.envoy.config.route.v3.FilterAction;
  return proto.envoy.config.route.v3.FilterAction.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.envoy.config.route.v3.FilterAction} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.envoy.config.route.v3.FilterAction}
 */
proto.envoy.config.route.v3.FilterAction.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new google_protobuf_any_pb.Any;
      reader.readMessage(value,google_protobuf_any_pb.Any.deserializeBinaryFromReader);
      msg.setAction(value);
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
proto.envoy.config.route.v3.FilterAction.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.envoy.config.route.v3.FilterAction.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.envoy.config.route.v3.FilterAction} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.FilterAction.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getAction();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      google_protobuf_any_pb.Any.serializeBinaryToWriter
    );
  }
};


/**
 * optional google.protobuf.Any action = 1;
 * @return {?proto.google.protobuf.Any}
 */
proto.envoy.config.route.v3.FilterAction.prototype.getAction = function() {
  return /** @type{?proto.google.protobuf.Any} */ (
    jspb.Message.getWrapperField(this, google_protobuf_any_pb.Any, 1));
};


/** @param {?proto.google.protobuf.Any|undefined} value */
proto.envoy.config.route.v3.FilterAction.prototype.setAction = function(value) {
  jspb.Message.setWrapperField(this, 1, value);
};


proto.envoy.config.route.v3.FilterAction.prototype.clearAction = function() {
  this.setAction(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.FilterAction.prototype.hasAction = function() {
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
proto.envoy.config.route.v3.Route = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.envoy.config.route.v3.Route.repeatedFields_, proto.envoy.config.route.v3.Route.oneofGroups_);
};
goog.inherits(proto.envoy.config.route.v3.Route, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.envoy.config.route.v3.Route.displayName = 'proto.envoy.config.route.v3.Route';
}
/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.envoy.config.route.v3.Route.repeatedFields_ = [9,12,10,11];

/**
 * Oneof group definitions for this message. Each group defines the field
 * numbers belonging to that group. When of these fields' value is set, all
 * other fields in the group are cleared. During deserialization, if multiple
 * fields are encountered for a group, only the last value seen will be kept.
 * @private {!Array<!Array<number>>}
 * @const
 */
proto.envoy.config.route.v3.Route.oneofGroups_ = [[2,3,7,17]];

/**
 * @enum {number}
 */
proto.envoy.config.route.v3.Route.ActionCase = {
  ACTION_NOT_SET: 0,
  ROUTE: 2,
  REDIRECT: 3,
  DIRECT_RESPONSE: 7,
  FILTER_ACTION: 17
};

/**
 * @return {proto.envoy.config.route.v3.Route.ActionCase}
 */
proto.envoy.config.route.v3.Route.prototype.getActionCase = function() {
  return /** @type {proto.envoy.config.route.v3.Route.ActionCase} */(jspb.Message.computeOneofCase(this, proto.envoy.config.route.v3.Route.oneofGroups_[0]));
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
proto.envoy.config.route.v3.Route.prototype.toObject = function(opt_includeInstance) {
  return proto.envoy.config.route.v3.Route.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.envoy.config.route.v3.Route} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.Route.toObject = function(includeInstance, msg) {
  var f, obj = {
    name: jspb.Message.getFieldWithDefault(msg, 14, ""),
    match: (f = msg.getMatch()) && proto.envoy.config.route.v3.RouteMatch.toObject(includeInstance, f),
    route: (f = msg.getRoute()) && proto.envoy.config.route.v3.RouteAction.toObject(includeInstance, f),
    redirect: (f = msg.getRedirect()) && proto.envoy.config.route.v3.RedirectAction.toObject(includeInstance, f),
    directResponse: (f = msg.getDirectResponse()) && proto.envoy.config.route.v3.DirectResponseAction.toObject(includeInstance, f),
    filterAction: (f = msg.getFilterAction()) && proto.envoy.config.route.v3.FilterAction.toObject(includeInstance, f),
    metadata: (f = msg.getMetadata()) && envoy_config_core_v3_base_pb.Metadata.toObject(includeInstance, f),
    decorator: (f = msg.getDecorator()) && proto.envoy.config.route.v3.Decorator.toObject(includeInstance, f),
    typedPerFilterConfigMap: (f = msg.getTypedPerFilterConfigMap()) ? f.toObject(includeInstance, proto.google.protobuf.Any.toObject) : [],
    requestHeadersToAddList: jspb.Message.toObjectList(msg.getRequestHeadersToAddList(),
    envoy_config_core_v3_base_pb.HeaderValueOption.toObject, includeInstance),
    requestHeadersToRemoveList: jspb.Message.getRepeatedField(msg, 12),
    responseHeadersToAddList: jspb.Message.toObjectList(msg.getResponseHeadersToAddList(),
    envoy_config_core_v3_base_pb.HeaderValueOption.toObject, includeInstance),
    responseHeadersToRemoveList: jspb.Message.getRepeatedField(msg, 11),
    tracing: (f = msg.getTracing()) && proto.envoy.config.route.v3.Tracing.toObject(includeInstance, f),
    perRequestBufferLimitBytes: (f = msg.getPerRequestBufferLimitBytes()) && google_protobuf_wrappers_pb.UInt32Value.toObject(includeInstance, f)
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
 * @return {!proto.envoy.config.route.v3.Route}
 */
proto.envoy.config.route.v3.Route.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.envoy.config.route.v3.Route;
  return proto.envoy.config.route.v3.Route.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.envoy.config.route.v3.Route} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.envoy.config.route.v3.Route}
 */
proto.envoy.config.route.v3.Route.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 14:
      var value = /** @type {string} */ (reader.readString());
      msg.setName(value);
      break;
    case 1:
      var value = new proto.envoy.config.route.v3.RouteMatch;
      reader.readMessage(value,proto.envoy.config.route.v3.RouteMatch.deserializeBinaryFromReader);
      msg.setMatch(value);
      break;
    case 2:
      var value = new proto.envoy.config.route.v3.RouteAction;
      reader.readMessage(value,proto.envoy.config.route.v3.RouteAction.deserializeBinaryFromReader);
      msg.setRoute(value);
      break;
    case 3:
      var value = new proto.envoy.config.route.v3.RedirectAction;
      reader.readMessage(value,proto.envoy.config.route.v3.RedirectAction.deserializeBinaryFromReader);
      msg.setRedirect(value);
      break;
    case 7:
      var value = new proto.envoy.config.route.v3.DirectResponseAction;
      reader.readMessage(value,proto.envoy.config.route.v3.DirectResponseAction.deserializeBinaryFromReader);
      msg.setDirectResponse(value);
      break;
    case 17:
      var value = new proto.envoy.config.route.v3.FilterAction;
      reader.readMessage(value,proto.envoy.config.route.v3.FilterAction.deserializeBinaryFromReader);
      msg.setFilterAction(value);
      break;
    case 4:
      var value = new envoy_config_core_v3_base_pb.Metadata;
      reader.readMessage(value,envoy_config_core_v3_base_pb.Metadata.deserializeBinaryFromReader);
      msg.setMetadata(value);
      break;
    case 5:
      var value = new proto.envoy.config.route.v3.Decorator;
      reader.readMessage(value,proto.envoy.config.route.v3.Decorator.deserializeBinaryFromReader);
      msg.setDecorator(value);
      break;
    case 13:
      var value = msg.getTypedPerFilterConfigMap();
      reader.readMessage(value, function(message, reader) {
        jspb.Map.deserializeBinary(message, reader, jspb.BinaryReader.prototype.readString, jspb.BinaryReader.prototype.readMessage, proto.google.protobuf.Any.deserializeBinaryFromReader, "");
         });
      break;
    case 9:
      var value = new envoy_config_core_v3_base_pb.HeaderValueOption;
      reader.readMessage(value,envoy_config_core_v3_base_pb.HeaderValueOption.deserializeBinaryFromReader);
      msg.addRequestHeadersToAdd(value);
      break;
    case 12:
      var value = /** @type {string} */ (reader.readString());
      msg.addRequestHeadersToRemove(value);
      break;
    case 10:
      var value = new envoy_config_core_v3_base_pb.HeaderValueOption;
      reader.readMessage(value,envoy_config_core_v3_base_pb.HeaderValueOption.deserializeBinaryFromReader);
      msg.addResponseHeadersToAdd(value);
      break;
    case 11:
      var value = /** @type {string} */ (reader.readString());
      msg.addResponseHeadersToRemove(value);
      break;
    case 15:
      var value = new proto.envoy.config.route.v3.Tracing;
      reader.readMessage(value,proto.envoy.config.route.v3.Tracing.deserializeBinaryFromReader);
      msg.setTracing(value);
      break;
    case 16:
      var value = new google_protobuf_wrappers_pb.UInt32Value;
      reader.readMessage(value,google_protobuf_wrappers_pb.UInt32Value.deserializeBinaryFromReader);
      msg.setPerRequestBufferLimitBytes(value);
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
proto.envoy.config.route.v3.Route.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.envoy.config.route.v3.Route.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.envoy.config.route.v3.Route} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.Route.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getName();
  if (f.length > 0) {
    writer.writeString(
      14,
      f
    );
  }
  f = message.getMatch();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto.envoy.config.route.v3.RouteMatch.serializeBinaryToWriter
    );
  }
  f = message.getRoute();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.envoy.config.route.v3.RouteAction.serializeBinaryToWriter
    );
  }
  f = message.getRedirect();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      proto.envoy.config.route.v3.RedirectAction.serializeBinaryToWriter
    );
  }
  f = message.getDirectResponse();
  if (f != null) {
    writer.writeMessage(
      7,
      f,
      proto.envoy.config.route.v3.DirectResponseAction.serializeBinaryToWriter
    );
  }
  f = message.getFilterAction();
  if (f != null) {
    writer.writeMessage(
      17,
      f,
      proto.envoy.config.route.v3.FilterAction.serializeBinaryToWriter
    );
  }
  f = message.getMetadata();
  if (f != null) {
    writer.writeMessage(
      4,
      f,
      envoy_config_core_v3_base_pb.Metadata.serializeBinaryToWriter
    );
  }
  f = message.getDecorator();
  if (f != null) {
    writer.writeMessage(
      5,
      f,
      proto.envoy.config.route.v3.Decorator.serializeBinaryToWriter
    );
  }
  f = message.getTypedPerFilterConfigMap(true);
  if (f && f.getLength() > 0) {
    f.serializeBinary(13, writer, jspb.BinaryWriter.prototype.writeString, jspb.BinaryWriter.prototype.writeMessage, proto.google.protobuf.Any.serializeBinaryToWriter);
  }
  f = message.getRequestHeadersToAddList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      9,
      f,
      envoy_config_core_v3_base_pb.HeaderValueOption.serializeBinaryToWriter
    );
  }
  f = message.getRequestHeadersToRemoveList();
  if (f.length > 0) {
    writer.writeRepeatedString(
      12,
      f
    );
  }
  f = message.getResponseHeadersToAddList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      10,
      f,
      envoy_config_core_v3_base_pb.HeaderValueOption.serializeBinaryToWriter
    );
  }
  f = message.getResponseHeadersToRemoveList();
  if (f.length > 0) {
    writer.writeRepeatedString(
      11,
      f
    );
  }
  f = message.getTracing();
  if (f != null) {
    writer.writeMessage(
      15,
      f,
      proto.envoy.config.route.v3.Tracing.serializeBinaryToWriter
    );
  }
  f = message.getPerRequestBufferLimitBytes();
  if (f != null) {
    writer.writeMessage(
      16,
      f,
      google_protobuf_wrappers_pb.UInt32Value.serializeBinaryToWriter
    );
  }
};


/**
 * optional string name = 14;
 * @return {string}
 */
proto.envoy.config.route.v3.Route.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 14, ""));
};


/** @param {string} value */
proto.envoy.config.route.v3.Route.prototype.setName = function(value) {
  jspb.Message.setProto3StringField(this, 14, value);
};


/**
 * optional RouteMatch match = 1;
 * @return {?proto.envoy.config.route.v3.RouteMatch}
 */
proto.envoy.config.route.v3.Route.prototype.getMatch = function() {
  return /** @type{?proto.envoy.config.route.v3.RouteMatch} */ (
    jspb.Message.getWrapperField(this, proto.envoy.config.route.v3.RouteMatch, 1));
};


/** @param {?proto.envoy.config.route.v3.RouteMatch|undefined} value */
proto.envoy.config.route.v3.Route.prototype.setMatch = function(value) {
  jspb.Message.setWrapperField(this, 1, value);
};


proto.envoy.config.route.v3.Route.prototype.clearMatch = function() {
  this.setMatch(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.Route.prototype.hasMatch = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional RouteAction route = 2;
 * @return {?proto.envoy.config.route.v3.RouteAction}
 */
proto.envoy.config.route.v3.Route.prototype.getRoute = function() {
  return /** @type{?proto.envoy.config.route.v3.RouteAction} */ (
    jspb.Message.getWrapperField(this, proto.envoy.config.route.v3.RouteAction, 2));
};


/** @param {?proto.envoy.config.route.v3.RouteAction|undefined} value */
proto.envoy.config.route.v3.Route.prototype.setRoute = function(value) {
  jspb.Message.setOneofWrapperField(this, 2, proto.envoy.config.route.v3.Route.oneofGroups_[0], value);
};


proto.envoy.config.route.v3.Route.prototype.clearRoute = function() {
  this.setRoute(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.Route.prototype.hasRoute = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional RedirectAction redirect = 3;
 * @return {?proto.envoy.config.route.v3.RedirectAction}
 */
proto.envoy.config.route.v3.Route.prototype.getRedirect = function() {
  return /** @type{?proto.envoy.config.route.v3.RedirectAction} */ (
    jspb.Message.getWrapperField(this, proto.envoy.config.route.v3.RedirectAction, 3));
};


/** @param {?proto.envoy.config.route.v3.RedirectAction|undefined} value */
proto.envoy.config.route.v3.Route.prototype.setRedirect = function(value) {
  jspb.Message.setOneofWrapperField(this, 3, proto.envoy.config.route.v3.Route.oneofGroups_[0], value);
};


proto.envoy.config.route.v3.Route.prototype.clearRedirect = function() {
  this.setRedirect(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.Route.prototype.hasRedirect = function() {
  return jspb.Message.getField(this, 3) != null;
};


/**
 * optional DirectResponseAction direct_response = 7;
 * @return {?proto.envoy.config.route.v3.DirectResponseAction}
 */
proto.envoy.config.route.v3.Route.prototype.getDirectResponse = function() {
  return /** @type{?proto.envoy.config.route.v3.DirectResponseAction} */ (
    jspb.Message.getWrapperField(this, proto.envoy.config.route.v3.DirectResponseAction, 7));
};


/** @param {?proto.envoy.config.route.v3.DirectResponseAction|undefined} value */
proto.envoy.config.route.v3.Route.prototype.setDirectResponse = function(value) {
  jspb.Message.setOneofWrapperField(this, 7, proto.envoy.config.route.v3.Route.oneofGroups_[0], value);
};


proto.envoy.config.route.v3.Route.prototype.clearDirectResponse = function() {
  this.setDirectResponse(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.Route.prototype.hasDirectResponse = function() {
  return jspb.Message.getField(this, 7) != null;
};


/**
 * optional FilterAction filter_action = 17;
 * @return {?proto.envoy.config.route.v3.FilterAction}
 */
proto.envoy.config.route.v3.Route.prototype.getFilterAction = function() {
  return /** @type{?proto.envoy.config.route.v3.FilterAction} */ (
    jspb.Message.getWrapperField(this, proto.envoy.config.route.v3.FilterAction, 17));
};


/** @param {?proto.envoy.config.route.v3.FilterAction|undefined} value */
proto.envoy.config.route.v3.Route.prototype.setFilterAction = function(value) {
  jspb.Message.setOneofWrapperField(this, 17, proto.envoy.config.route.v3.Route.oneofGroups_[0], value);
};


proto.envoy.config.route.v3.Route.prototype.clearFilterAction = function() {
  this.setFilterAction(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.Route.prototype.hasFilterAction = function() {
  return jspb.Message.getField(this, 17) != null;
};


/**
 * optional envoy.config.core.v3.Metadata metadata = 4;
 * @return {?proto.envoy.config.core.v3.Metadata}
 */
proto.envoy.config.route.v3.Route.prototype.getMetadata = function() {
  return /** @type{?proto.envoy.config.core.v3.Metadata} */ (
    jspb.Message.getWrapperField(this, envoy_config_core_v3_base_pb.Metadata, 4));
};


/** @param {?proto.envoy.config.core.v3.Metadata|undefined} value */
proto.envoy.config.route.v3.Route.prototype.setMetadata = function(value) {
  jspb.Message.setWrapperField(this, 4, value);
};


proto.envoy.config.route.v3.Route.prototype.clearMetadata = function() {
  this.setMetadata(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.Route.prototype.hasMetadata = function() {
  return jspb.Message.getField(this, 4) != null;
};


/**
 * optional Decorator decorator = 5;
 * @return {?proto.envoy.config.route.v3.Decorator}
 */
proto.envoy.config.route.v3.Route.prototype.getDecorator = function() {
  return /** @type{?proto.envoy.config.route.v3.Decorator} */ (
    jspb.Message.getWrapperField(this, proto.envoy.config.route.v3.Decorator, 5));
};


/** @param {?proto.envoy.config.route.v3.Decorator|undefined} value */
proto.envoy.config.route.v3.Route.prototype.setDecorator = function(value) {
  jspb.Message.setWrapperField(this, 5, value);
};


proto.envoy.config.route.v3.Route.prototype.clearDecorator = function() {
  this.setDecorator(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.Route.prototype.hasDecorator = function() {
  return jspb.Message.getField(this, 5) != null;
};


/**
 * map<string, google.protobuf.Any> typed_per_filter_config = 13;
 * @param {boolean=} opt_noLazyCreate Do not create the map if
 * empty, instead returning `undefined`
 * @return {!jspb.Map<string,!proto.google.protobuf.Any>}
 */
proto.envoy.config.route.v3.Route.prototype.getTypedPerFilterConfigMap = function(opt_noLazyCreate) {
  return /** @type {!jspb.Map<string,!proto.google.protobuf.Any>} */ (
      jspb.Message.getMapField(this, 13, opt_noLazyCreate,
      proto.google.protobuf.Any));
};


proto.envoy.config.route.v3.Route.prototype.clearTypedPerFilterConfigMap = function() {
  this.getTypedPerFilterConfigMap().clear();
};


/**
 * repeated envoy.config.core.v3.HeaderValueOption request_headers_to_add = 9;
 * @return {!Array<!proto.envoy.config.core.v3.HeaderValueOption>}
 */
proto.envoy.config.route.v3.Route.prototype.getRequestHeadersToAddList = function() {
  return /** @type{!Array<!proto.envoy.config.core.v3.HeaderValueOption>} */ (
    jspb.Message.getRepeatedWrapperField(this, envoy_config_core_v3_base_pb.HeaderValueOption, 9));
};


/** @param {!Array<!proto.envoy.config.core.v3.HeaderValueOption>} value */
proto.envoy.config.route.v3.Route.prototype.setRequestHeadersToAddList = function(value) {
  jspb.Message.setRepeatedWrapperField(this, 9, value);
};


/**
 * @param {!proto.envoy.config.core.v3.HeaderValueOption=} opt_value
 * @param {number=} opt_index
 * @return {!proto.envoy.config.core.v3.HeaderValueOption}
 */
proto.envoy.config.route.v3.Route.prototype.addRequestHeadersToAdd = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 9, opt_value, proto.envoy.config.core.v3.HeaderValueOption, opt_index);
};


proto.envoy.config.route.v3.Route.prototype.clearRequestHeadersToAddList = function() {
  this.setRequestHeadersToAddList([]);
};


/**
 * repeated string request_headers_to_remove = 12;
 * @return {!Array<string>}
 */
proto.envoy.config.route.v3.Route.prototype.getRequestHeadersToRemoveList = function() {
  return /** @type {!Array<string>} */ (jspb.Message.getRepeatedField(this, 12));
};


/** @param {!Array<string>} value */
proto.envoy.config.route.v3.Route.prototype.setRequestHeadersToRemoveList = function(value) {
  jspb.Message.setField(this, 12, value || []);
};


/**
 * @param {!string} value
 * @param {number=} opt_index
 */
proto.envoy.config.route.v3.Route.prototype.addRequestHeadersToRemove = function(value, opt_index) {
  jspb.Message.addToRepeatedField(this, 12, value, opt_index);
};


proto.envoy.config.route.v3.Route.prototype.clearRequestHeadersToRemoveList = function() {
  this.setRequestHeadersToRemoveList([]);
};


/**
 * repeated envoy.config.core.v3.HeaderValueOption response_headers_to_add = 10;
 * @return {!Array<!proto.envoy.config.core.v3.HeaderValueOption>}
 */
proto.envoy.config.route.v3.Route.prototype.getResponseHeadersToAddList = function() {
  return /** @type{!Array<!proto.envoy.config.core.v3.HeaderValueOption>} */ (
    jspb.Message.getRepeatedWrapperField(this, envoy_config_core_v3_base_pb.HeaderValueOption, 10));
};


/** @param {!Array<!proto.envoy.config.core.v3.HeaderValueOption>} value */
proto.envoy.config.route.v3.Route.prototype.setResponseHeadersToAddList = function(value) {
  jspb.Message.setRepeatedWrapperField(this, 10, value);
};


/**
 * @param {!proto.envoy.config.core.v3.HeaderValueOption=} opt_value
 * @param {number=} opt_index
 * @return {!proto.envoy.config.core.v3.HeaderValueOption}
 */
proto.envoy.config.route.v3.Route.prototype.addResponseHeadersToAdd = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 10, opt_value, proto.envoy.config.core.v3.HeaderValueOption, opt_index);
};


proto.envoy.config.route.v3.Route.prototype.clearResponseHeadersToAddList = function() {
  this.setResponseHeadersToAddList([]);
};


/**
 * repeated string response_headers_to_remove = 11;
 * @return {!Array<string>}
 */
proto.envoy.config.route.v3.Route.prototype.getResponseHeadersToRemoveList = function() {
  return /** @type {!Array<string>} */ (jspb.Message.getRepeatedField(this, 11));
};


/** @param {!Array<string>} value */
proto.envoy.config.route.v3.Route.prototype.setResponseHeadersToRemoveList = function(value) {
  jspb.Message.setField(this, 11, value || []);
};


/**
 * @param {!string} value
 * @param {number=} opt_index
 */
proto.envoy.config.route.v3.Route.prototype.addResponseHeadersToRemove = function(value, opt_index) {
  jspb.Message.addToRepeatedField(this, 11, value, opt_index);
};


proto.envoy.config.route.v3.Route.prototype.clearResponseHeadersToRemoveList = function() {
  this.setResponseHeadersToRemoveList([]);
};


/**
 * optional Tracing tracing = 15;
 * @return {?proto.envoy.config.route.v3.Tracing}
 */
proto.envoy.config.route.v3.Route.prototype.getTracing = function() {
  return /** @type{?proto.envoy.config.route.v3.Tracing} */ (
    jspb.Message.getWrapperField(this, proto.envoy.config.route.v3.Tracing, 15));
};


/** @param {?proto.envoy.config.route.v3.Tracing|undefined} value */
proto.envoy.config.route.v3.Route.prototype.setTracing = function(value) {
  jspb.Message.setWrapperField(this, 15, value);
};


proto.envoy.config.route.v3.Route.prototype.clearTracing = function() {
  this.setTracing(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.Route.prototype.hasTracing = function() {
  return jspb.Message.getField(this, 15) != null;
};


/**
 * optional google.protobuf.UInt32Value per_request_buffer_limit_bytes = 16;
 * @return {?proto.google.protobuf.UInt32Value}
 */
proto.envoy.config.route.v3.Route.prototype.getPerRequestBufferLimitBytes = function() {
  return /** @type{?proto.google.protobuf.UInt32Value} */ (
    jspb.Message.getWrapperField(this, google_protobuf_wrappers_pb.UInt32Value, 16));
};


/** @param {?proto.google.protobuf.UInt32Value|undefined} value */
proto.envoy.config.route.v3.Route.prototype.setPerRequestBufferLimitBytes = function(value) {
  jspb.Message.setWrapperField(this, 16, value);
};


proto.envoy.config.route.v3.Route.prototype.clearPerRequestBufferLimitBytes = function() {
  this.setPerRequestBufferLimitBytes(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.Route.prototype.hasPerRequestBufferLimitBytes = function() {
  return jspb.Message.getField(this, 16) != null;
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
proto.envoy.config.route.v3.WeightedCluster = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.envoy.config.route.v3.WeightedCluster.repeatedFields_, null);
};
goog.inherits(proto.envoy.config.route.v3.WeightedCluster, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.envoy.config.route.v3.WeightedCluster.displayName = 'proto.envoy.config.route.v3.WeightedCluster';
}
/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.envoy.config.route.v3.WeightedCluster.repeatedFields_ = [1];



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
proto.envoy.config.route.v3.WeightedCluster.prototype.toObject = function(opt_includeInstance) {
  return proto.envoy.config.route.v3.WeightedCluster.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.envoy.config.route.v3.WeightedCluster} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.WeightedCluster.toObject = function(includeInstance, msg) {
  var f, obj = {
    clustersList: jspb.Message.toObjectList(msg.getClustersList(),
    proto.envoy.config.route.v3.WeightedCluster.ClusterWeight.toObject, includeInstance),
    totalWeight: (f = msg.getTotalWeight()) && google_protobuf_wrappers_pb.UInt32Value.toObject(includeInstance, f),
    runtimeKeyPrefix: jspb.Message.getFieldWithDefault(msg, 2, "")
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
 * @return {!proto.envoy.config.route.v3.WeightedCluster}
 */
proto.envoy.config.route.v3.WeightedCluster.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.envoy.config.route.v3.WeightedCluster;
  return proto.envoy.config.route.v3.WeightedCluster.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.envoy.config.route.v3.WeightedCluster} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.envoy.config.route.v3.WeightedCluster}
 */
proto.envoy.config.route.v3.WeightedCluster.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.envoy.config.route.v3.WeightedCluster.ClusterWeight;
      reader.readMessage(value,proto.envoy.config.route.v3.WeightedCluster.ClusterWeight.deserializeBinaryFromReader);
      msg.addClusters(value);
      break;
    case 3:
      var value = new google_protobuf_wrappers_pb.UInt32Value;
      reader.readMessage(value,google_protobuf_wrappers_pb.UInt32Value.deserializeBinaryFromReader);
      msg.setTotalWeight(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setRuntimeKeyPrefix(value);
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
proto.envoy.config.route.v3.WeightedCluster.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.envoy.config.route.v3.WeightedCluster.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.envoy.config.route.v3.WeightedCluster} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.WeightedCluster.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getClustersList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      1,
      f,
      proto.envoy.config.route.v3.WeightedCluster.ClusterWeight.serializeBinaryToWriter
    );
  }
  f = message.getTotalWeight();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      google_protobuf_wrappers_pb.UInt32Value.serializeBinaryToWriter
    );
  }
  f = message.getRuntimeKeyPrefix();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
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
proto.envoy.config.route.v3.WeightedCluster.ClusterWeight = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.envoy.config.route.v3.WeightedCluster.ClusterWeight.repeatedFields_, null);
};
goog.inherits(proto.envoy.config.route.v3.WeightedCluster.ClusterWeight, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.envoy.config.route.v3.WeightedCluster.ClusterWeight.displayName = 'proto.envoy.config.route.v3.WeightedCluster.ClusterWeight';
}
/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.envoy.config.route.v3.WeightedCluster.ClusterWeight.repeatedFields_ = [4,9,5,6];



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
proto.envoy.config.route.v3.WeightedCluster.ClusterWeight.prototype.toObject = function(opt_includeInstance) {
  return proto.envoy.config.route.v3.WeightedCluster.ClusterWeight.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.envoy.config.route.v3.WeightedCluster.ClusterWeight} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.WeightedCluster.ClusterWeight.toObject = function(includeInstance, msg) {
  var f, obj = {
    name: jspb.Message.getFieldWithDefault(msg, 1, ""),
    weight: (f = msg.getWeight()) && google_protobuf_wrappers_pb.UInt32Value.toObject(includeInstance, f),
    metadataMatch: (f = msg.getMetadataMatch()) && envoy_config_core_v3_base_pb.Metadata.toObject(includeInstance, f),
    requestHeadersToAddList: jspb.Message.toObjectList(msg.getRequestHeadersToAddList(),
    envoy_config_core_v3_base_pb.HeaderValueOption.toObject, includeInstance),
    requestHeadersToRemoveList: jspb.Message.getRepeatedField(msg, 9),
    responseHeadersToAddList: jspb.Message.toObjectList(msg.getResponseHeadersToAddList(),
    envoy_config_core_v3_base_pb.HeaderValueOption.toObject, includeInstance),
    responseHeadersToRemoveList: jspb.Message.getRepeatedField(msg, 6),
    typedPerFilterConfigMap: (f = msg.getTypedPerFilterConfigMap()) ? f.toObject(includeInstance, proto.google.protobuf.Any.toObject) : []
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
 * @return {!proto.envoy.config.route.v3.WeightedCluster.ClusterWeight}
 */
proto.envoy.config.route.v3.WeightedCluster.ClusterWeight.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.envoy.config.route.v3.WeightedCluster.ClusterWeight;
  return proto.envoy.config.route.v3.WeightedCluster.ClusterWeight.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.envoy.config.route.v3.WeightedCluster.ClusterWeight} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.envoy.config.route.v3.WeightedCluster.ClusterWeight}
 */
proto.envoy.config.route.v3.WeightedCluster.ClusterWeight.deserializeBinaryFromReader = function(msg, reader) {
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
      var value = new google_protobuf_wrappers_pb.UInt32Value;
      reader.readMessage(value,google_protobuf_wrappers_pb.UInt32Value.deserializeBinaryFromReader);
      msg.setWeight(value);
      break;
    case 3:
      var value = new envoy_config_core_v3_base_pb.Metadata;
      reader.readMessage(value,envoy_config_core_v3_base_pb.Metadata.deserializeBinaryFromReader);
      msg.setMetadataMatch(value);
      break;
    case 4:
      var value = new envoy_config_core_v3_base_pb.HeaderValueOption;
      reader.readMessage(value,envoy_config_core_v3_base_pb.HeaderValueOption.deserializeBinaryFromReader);
      msg.addRequestHeadersToAdd(value);
      break;
    case 9:
      var value = /** @type {string} */ (reader.readString());
      msg.addRequestHeadersToRemove(value);
      break;
    case 5:
      var value = new envoy_config_core_v3_base_pb.HeaderValueOption;
      reader.readMessage(value,envoy_config_core_v3_base_pb.HeaderValueOption.deserializeBinaryFromReader);
      msg.addResponseHeadersToAdd(value);
      break;
    case 6:
      var value = /** @type {string} */ (reader.readString());
      msg.addResponseHeadersToRemove(value);
      break;
    case 10:
      var value = msg.getTypedPerFilterConfigMap();
      reader.readMessage(value, function(message, reader) {
        jspb.Map.deserializeBinary(message, reader, jspb.BinaryReader.prototype.readString, jspb.BinaryReader.prototype.readMessage, proto.google.protobuf.Any.deserializeBinaryFromReader, "");
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
proto.envoy.config.route.v3.WeightedCluster.ClusterWeight.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.envoy.config.route.v3.WeightedCluster.ClusterWeight.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.envoy.config.route.v3.WeightedCluster.ClusterWeight} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.WeightedCluster.ClusterWeight.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getName();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getWeight();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      google_protobuf_wrappers_pb.UInt32Value.serializeBinaryToWriter
    );
  }
  f = message.getMetadataMatch();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      envoy_config_core_v3_base_pb.Metadata.serializeBinaryToWriter
    );
  }
  f = message.getRequestHeadersToAddList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      4,
      f,
      envoy_config_core_v3_base_pb.HeaderValueOption.serializeBinaryToWriter
    );
  }
  f = message.getRequestHeadersToRemoveList();
  if (f.length > 0) {
    writer.writeRepeatedString(
      9,
      f
    );
  }
  f = message.getResponseHeadersToAddList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      5,
      f,
      envoy_config_core_v3_base_pb.HeaderValueOption.serializeBinaryToWriter
    );
  }
  f = message.getResponseHeadersToRemoveList();
  if (f.length > 0) {
    writer.writeRepeatedString(
      6,
      f
    );
  }
  f = message.getTypedPerFilterConfigMap(true);
  if (f && f.getLength() > 0) {
    f.serializeBinary(10, writer, jspb.BinaryWriter.prototype.writeString, jspb.BinaryWriter.prototype.writeMessage, proto.google.protobuf.Any.serializeBinaryToWriter);
  }
};


/**
 * optional string name = 1;
 * @return {string}
 */
proto.envoy.config.route.v3.WeightedCluster.ClusterWeight.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.envoy.config.route.v3.WeightedCluster.ClusterWeight.prototype.setName = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional google.protobuf.UInt32Value weight = 2;
 * @return {?proto.google.protobuf.UInt32Value}
 */
proto.envoy.config.route.v3.WeightedCluster.ClusterWeight.prototype.getWeight = function() {
  return /** @type{?proto.google.protobuf.UInt32Value} */ (
    jspb.Message.getWrapperField(this, google_protobuf_wrappers_pb.UInt32Value, 2));
};


/** @param {?proto.google.protobuf.UInt32Value|undefined} value */
proto.envoy.config.route.v3.WeightedCluster.ClusterWeight.prototype.setWeight = function(value) {
  jspb.Message.setWrapperField(this, 2, value);
};


proto.envoy.config.route.v3.WeightedCluster.ClusterWeight.prototype.clearWeight = function() {
  this.setWeight(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.WeightedCluster.ClusterWeight.prototype.hasWeight = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional envoy.config.core.v3.Metadata metadata_match = 3;
 * @return {?proto.envoy.config.core.v3.Metadata}
 */
proto.envoy.config.route.v3.WeightedCluster.ClusterWeight.prototype.getMetadataMatch = function() {
  return /** @type{?proto.envoy.config.core.v3.Metadata} */ (
    jspb.Message.getWrapperField(this, envoy_config_core_v3_base_pb.Metadata, 3));
};


/** @param {?proto.envoy.config.core.v3.Metadata|undefined} value */
proto.envoy.config.route.v3.WeightedCluster.ClusterWeight.prototype.setMetadataMatch = function(value) {
  jspb.Message.setWrapperField(this, 3, value);
};


proto.envoy.config.route.v3.WeightedCluster.ClusterWeight.prototype.clearMetadataMatch = function() {
  this.setMetadataMatch(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.WeightedCluster.ClusterWeight.prototype.hasMetadataMatch = function() {
  return jspb.Message.getField(this, 3) != null;
};


/**
 * repeated envoy.config.core.v3.HeaderValueOption request_headers_to_add = 4;
 * @return {!Array<!proto.envoy.config.core.v3.HeaderValueOption>}
 */
proto.envoy.config.route.v3.WeightedCluster.ClusterWeight.prototype.getRequestHeadersToAddList = function() {
  return /** @type{!Array<!proto.envoy.config.core.v3.HeaderValueOption>} */ (
    jspb.Message.getRepeatedWrapperField(this, envoy_config_core_v3_base_pb.HeaderValueOption, 4));
};


/** @param {!Array<!proto.envoy.config.core.v3.HeaderValueOption>} value */
proto.envoy.config.route.v3.WeightedCluster.ClusterWeight.prototype.setRequestHeadersToAddList = function(value) {
  jspb.Message.setRepeatedWrapperField(this, 4, value);
};


/**
 * @param {!proto.envoy.config.core.v3.HeaderValueOption=} opt_value
 * @param {number=} opt_index
 * @return {!proto.envoy.config.core.v3.HeaderValueOption}
 */
proto.envoy.config.route.v3.WeightedCluster.ClusterWeight.prototype.addRequestHeadersToAdd = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 4, opt_value, proto.envoy.config.core.v3.HeaderValueOption, opt_index);
};


proto.envoy.config.route.v3.WeightedCluster.ClusterWeight.prototype.clearRequestHeadersToAddList = function() {
  this.setRequestHeadersToAddList([]);
};


/**
 * repeated string request_headers_to_remove = 9;
 * @return {!Array<string>}
 */
proto.envoy.config.route.v3.WeightedCluster.ClusterWeight.prototype.getRequestHeadersToRemoveList = function() {
  return /** @type {!Array<string>} */ (jspb.Message.getRepeatedField(this, 9));
};


/** @param {!Array<string>} value */
proto.envoy.config.route.v3.WeightedCluster.ClusterWeight.prototype.setRequestHeadersToRemoveList = function(value) {
  jspb.Message.setField(this, 9, value || []);
};


/**
 * @param {!string} value
 * @param {number=} opt_index
 */
proto.envoy.config.route.v3.WeightedCluster.ClusterWeight.prototype.addRequestHeadersToRemove = function(value, opt_index) {
  jspb.Message.addToRepeatedField(this, 9, value, opt_index);
};


proto.envoy.config.route.v3.WeightedCluster.ClusterWeight.prototype.clearRequestHeadersToRemoveList = function() {
  this.setRequestHeadersToRemoveList([]);
};


/**
 * repeated envoy.config.core.v3.HeaderValueOption response_headers_to_add = 5;
 * @return {!Array<!proto.envoy.config.core.v3.HeaderValueOption>}
 */
proto.envoy.config.route.v3.WeightedCluster.ClusterWeight.prototype.getResponseHeadersToAddList = function() {
  return /** @type{!Array<!proto.envoy.config.core.v3.HeaderValueOption>} */ (
    jspb.Message.getRepeatedWrapperField(this, envoy_config_core_v3_base_pb.HeaderValueOption, 5));
};


/** @param {!Array<!proto.envoy.config.core.v3.HeaderValueOption>} value */
proto.envoy.config.route.v3.WeightedCluster.ClusterWeight.prototype.setResponseHeadersToAddList = function(value) {
  jspb.Message.setRepeatedWrapperField(this, 5, value);
};


/**
 * @param {!proto.envoy.config.core.v3.HeaderValueOption=} opt_value
 * @param {number=} opt_index
 * @return {!proto.envoy.config.core.v3.HeaderValueOption}
 */
proto.envoy.config.route.v3.WeightedCluster.ClusterWeight.prototype.addResponseHeadersToAdd = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 5, opt_value, proto.envoy.config.core.v3.HeaderValueOption, opt_index);
};


proto.envoy.config.route.v3.WeightedCluster.ClusterWeight.prototype.clearResponseHeadersToAddList = function() {
  this.setResponseHeadersToAddList([]);
};


/**
 * repeated string response_headers_to_remove = 6;
 * @return {!Array<string>}
 */
proto.envoy.config.route.v3.WeightedCluster.ClusterWeight.prototype.getResponseHeadersToRemoveList = function() {
  return /** @type {!Array<string>} */ (jspb.Message.getRepeatedField(this, 6));
};


/** @param {!Array<string>} value */
proto.envoy.config.route.v3.WeightedCluster.ClusterWeight.prototype.setResponseHeadersToRemoveList = function(value) {
  jspb.Message.setField(this, 6, value || []);
};


/**
 * @param {!string} value
 * @param {number=} opt_index
 */
proto.envoy.config.route.v3.WeightedCluster.ClusterWeight.prototype.addResponseHeadersToRemove = function(value, opt_index) {
  jspb.Message.addToRepeatedField(this, 6, value, opt_index);
};


proto.envoy.config.route.v3.WeightedCluster.ClusterWeight.prototype.clearResponseHeadersToRemoveList = function() {
  this.setResponseHeadersToRemoveList([]);
};


/**
 * map<string, google.protobuf.Any> typed_per_filter_config = 10;
 * @param {boolean=} opt_noLazyCreate Do not create the map if
 * empty, instead returning `undefined`
 * @return {!jspb.Map<string,!proto.google.protobuf.Any>}
 */
proto.envoy.config.route.v3.WeightedCluster.ClusterWeight.prototype.getTypedPerFilterConfigMap = function(opt_noLazyCreate) {
  return /** @type {!jspb.Map<string,!proto.google.protobuf.Any>} */ (
      jspb.Message.getMapField(this, 10, opt_noLazyCreate,
      proto.google.protobuf.Any));
};


proto.envoy.config.route.v3.WeightedCluster.ClusterWeight.prototype.clearTypedPerFilterConfigMap = function() {
  this.getTypedPerFilterConfigMap().clear();
};


/**
 * repeated ClusterWeight clusters = 1;
 * @return {!Array<!proto.envoy.config.route.v3.WeightedCluster.ClusterWeight>}
 */
proto.envoy.config.route.v3.WeightedCluster.prototype.getClustersList = function() {
  return /** @type{!Array<!proto.envoy.config.route.v3.WeightedCluster.ClusterWeight>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.envoy.config.route.v3.WeightedCluster.ClusterWeight, 1));
};


/** @param {!Array<!proto.envoy.config.route.v3.WeightedCluster.ClusterWeight>} value */
proto.envoy.config.route.v3.WeightedCluster.prototype.setClustersList = function(value) {
  jspb.Message.setRepeatedWrapperField(this, 1, value);
};


/**
 * @param {!proto.envoy.config.route.v3.WeightedCluster.ClusterWeight=} opt_value
 * @param {number=} opt_index
 * @return {!proto.envoy.config.route.v3.WeightedCluster.ClusterWeight}
 */
proto.envoy.config.route.v3.WeightedCluster.prototype.addClusters = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 1, opt_value, proto.envoy.config.route.v3.WeightedCluster.ClusterWeight, opt_index);
};


proto.envoy.config.route.v3.WeightedCluster.prototype.clearClustersList = function() {
  this.setClustersList([]);
};


/**
 * optional google.protobuf.UInt32Value total_weight = 3;
 * @return {?proto.google.protobuf.UInt32Value}
 */
proto.envoy.config.route.v3.WeightedCluster.prototype.getTotalWeight = function() {
  return /** @type{?proto.google.protobuf.UInt32Value} */ (
    jspb.Message.getWrapperField(this, google_protobuf_wrappers_pb.UInt32Value, 3));
};


/** @param {?proto.google.protobuf.UInt32Value|undefined} value */
proto.envoy.config.route.v3.WeightedCluster.prototype.setTotalWeight = function(value) {
  jspb.Message.setWrapperField(this, 3, value);
};


proto.envoy.config.route.v3.WeightedCluster.prototype.clearTotalWeight = function() {
  this.setTotalWeight(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.WeightedCluster.prototype.hasTotalWeight = function() {
  return jspb.Message.getField(this, 3) != null;
};


/**
 * optional string runtime_key_prefix = 2;
 * @return {string}
 */
proto.envoy.config.route.v3.WeightedCluster.prototype.getRuntimeKeyPrefix = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.envoy.config.route.v3.WeightedCluster.prototype.setRuntimeKeyPrefix = function(value) {
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
proto.envoy.config.route.v3.RouteMatch = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.envoy.config.route.v3.RouteMatch.repeatedFields_, proto.envoy.config.route.v3.RouteMatch.oneofGroups_);
};
goog.inherits(proto.envoy.config.route.v3.RouteMatch, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.envoy.config.route.v3.RouteMatch.displayName = 'proto.envoy.config.route.v3.RouteMatch';
}
/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.envoy.config.route.v3.RouteMatch.repeatedFields_ = [6,7];

/**
 * Oneof group definitions for this message. Each group defines the field
 * numbers belonging to that group. When of these fields' value is set, all
 * other fields in the group are cleared. During deserialization, if multiple
 * fields are encountered for a group, only the last value seen will be kept.
 * @private {!Array<!Array<number>>}
 * @const
 */
proto.envoy.config.route.v3.RouteMatch.oneofGroups_ = [[1,2,10,12]];

/**
 * @enum {number}
 */
proto.envoy.config.route.v3.RouteMatch.PathSpecifierCase = {
  PATH_SPECIFIER_NOT_SET: 0,
  PREFIX: 1,
  PATH: 2,
  SAFE_REGEX: 10,
  CONNECT_MATCHER: 12
};

/**
 * @return {proto.envoy.config.route.v3.RouteMatch.PathSpecifierCase}
 */
proto.envoy.config.route.v3.RouteMatch.prototype.getPathSpecifierCase = function() {
  return /** @type {proto.envoy.config.route.v3.RouteMatch.PathSpecifierCase} */(jspb.Message.computeOneofCase(this, proto.envoy.config.route.v3.RouteMatch.oneofGroups_[0]));
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
proto.envoy.config.route.v3.RouteMatch.prototype.toObject = function(opt_includeInstance) {
  return proto.envoy.config.route.v3.RouteMatch.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.envoy.config.route.v3.RouteMatch} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.RouteMatch.toObject = function(includeInstance, msg) {
  var f, obj = {
    prefix: jspb.Message.getFieldWithDefault(msg, 1, ""),
    path: jspb.Message.getFieldWithDefault(msg, 2, ""),
    safeRegex: (f = msg.getSafeRegex()) && envoy_type_matcher_v3_regex_pb.RegexMatcher.toObject(includeInstance, f),
    connectMatcher: (f = msg.getConnectMatcher()) && proto.envoy.config.route.v3.RouteMatch.ConnectMatcher.toObject(includeInstance, f),
    caseSensitive: (f = msg.getCaseSensitive()) && google_protobuf_wrappers_pb.BoolValue.toObject(includeInstance, f),
    runtimeFraction: (f = msg.getRuntimeFraction()) && envoy_config_core_v3_base_pb.RuntimeFractionalPercent.toObject(includeInstance, f),
    headersList: jspb.Message.toObjectList(msg.getHeadersList(),
    proto.envoy.config.route.v3.HeaderMatcher.toObject, includeInstance),
    queryParametersList: jspb.Message.toObjectList(msg.getQueryParametersList(),
    proto.envoy.config.route.v3.QueryParameterMatcher.toObject, includeInstance),
    grpc: (f = msg.getGrpc()) && proto.envoy.config.route.v3.RouteMatch.GrpcRouteMatchOptions.toObject(includeInstance, f),
    tlsContext: (f = msg.getTlsContext()) && proto.envoy.config.route.v3.RouteMatch.TlsContextMatchOptions.toObject(includeInstance, f)
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
 * @return {!proto.envoy.config.route.v3.RouteMatch}
 */
proto.envoy.config.route.v3.RouteMatch.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.envoy.config.route.v3.RouteMatch;
  return proto.envoy.config.route.v3.RouteMatch.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.envoy.config.route.v3.RouteMatch} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.envoy.config.route.v3.RouteMatch}
 */
proto.envoy.config.route.v3.RouteMatch.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setPrefix(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setPath(value);
      break;
    case 10:
      var value = new envoy_type_matcher_v3_regex_pb.RegexMatcher;
      reader.readMessage(value,envoy_type_matcher_v3_regex_pb.RegexMatcher.deserializeBinaryFromReader);
      msg.setSafeRegex(value);
      break;
    case 12:
      var value = new proto.envoy.config.route.v3.RouteMatch.ConnectMatcher;
      reader.readMessage(value,proto.envoy.config.route.v3.RouteMatch.ConnectMatcher.deserializeBinaryFromReader);
      msg.setConnectMatcher(value);
      break;
    case 4:
      var value = new google_protobuf_wrappers_pb.BoolValue;
      reader.readMessage(value,google_protobuf_wrappers_pb.BoolValue.deserializeBinaryFromReader);
      msg.setCaseSensitive(value);
      break;
    case 9:
      var value = new envoy_config_core_v3_base_pb.RuntimeFractionalPercent;
      reader.readMessage(value,envoy_config_core_v3_base_pb.RuntimeFractionalPercent.deserializeBinaryFromReader);
      msg.setRuntimeFraction(value);
      break;
    case 6:
      var value = new proto.envoy.config.route.v3.HeaderMatcher;
      reader.readMessage(value,proto.envoy.config.route.v3.HeaderMatcher.deserializeBinaryFromReader);
      msg.addHeaders(value);
      break;
    case 7:
      var value = new proto.envoy.config.route.v3.QueryParameterMatcher;
      reader.readMessage(value,proto.envoy.config.route.v3.QueryParameterMatcher.deserializeBinaryFromReader);
      msg.addQueryParameters(value);
      break;
    case 8:
      var value = new proto.envoy.config.route.v3.RouteMatch.GrpcRouteMatchOptions;
      reader.readMessage(value,proto.envoy.config.route.v3.RouteMatch.GrpcRouteMatchOptions.deserializeBinaryFromReader);
      msg.setGrpc(value);
      break;
    case 11:
      var value = new proto.envoy.config.route.v3.RouteMatch.TlsContextMatchOptions;
      reader.readMessage(value,proto.envoy.config.route.v3.RouteMatch.TlsContextMatchOptions.deserializeBinaryFromReader);
      msg.setTlsContext(value);
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
proto.envoy.config.route.v3.RouteMatch.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.envoy.config.route.v3.RouteMatch.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.envoy.config.route.v3.RouteMatch} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.RouteMatch.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = /** @type {string} */ (jspb.Message.getField(message, 1));
  if (f != null) {
    writer.writeString(
      1,
      f
    );
  }
  f = /** @type {string} */ (jspb.Message.getField(message, 2));
  if (f != null) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getSafeRegex();
  if (f != null) {
    writer.writeMessage(
      10,
      f,
      envoy_type_matcher_v3_regex_pb.RegexMatcher.serializeBinaryToWriter
    );
  }
  f = message.getConnectMatcher();
  if (f != null) {
    writer.writeMessage(
      12,
      f,
      proto.envoy.config.route.v3.RouteMatch.ConnectMatcher.serializeBinaryToWriter
    );
  }
  f = message.getCaseSensitive();
  if (f != null) {
    writer.writeMessage(
      4,
      f,
      google_protobuf_wrappers_pb.BoolValue.serializeBinaryToWriter
    );
  }
  f = message.getRuntimeFraction();
  if (f != null) {
    writer.writeMessage(
      9,
      f,
      envoy_config_core_v3_base_pb.RuntimeFractionalPercent.serializeBinaryToWriter
    );
  }
  f = message.getHeadersList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      6,
      f,
      proto.envoy.config.route.v3.HeaderMatcher.serializeBinaryToWriter
    );
  }
  f = message.getQueryParametersList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      7,
      f,
      proto.envoy.config.route.v3.QueryParameterMatcher.serializeBinaryToWriter
    );
  }
  f = message.getGrpc();
  if (f != null) {
    writer.writeMessage(
      8,
      f,
      proto.envoy.config.route.v3.RouteMatch.GrpcRouteMatchOptions.serializeBinaryToWriter
    );
  }
  f = message.getTlsContext();
  if (f != null) {
    writer.writeMessage(
      11,
      f,
      proto.envoy.config.route.v3.RouteMatch.TlsContextMatchOptions.serializeBinaryToWriter
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
proto.envoy.config.route.v3.RouteMatch.GrpcRouteMatchOptions = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.envoy.config.route.v3.RouteMatch.GrpcRouteMatchOptions, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.envoy.config.route.v3.RouteMatch.GrpcRouteMatchOptions.displayName = 'proto.envoy.config.route.v3.RouteMatch.GrpcRouteMatchOptions';
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
proto.envoy.config.route.v3.RouteMatch.GrpcRouteMatchOptions.prototype.toObject = function(opt_includeInstance) {
  return proto.envoy.config.route.v3.RouteMatch.GrpcRouteMatchOptions.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.envoy.config.route.v3.RouteMatch.GrpcRouteMatchOptions} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.RouteMatch.GrpcRouteMatchOptions.toObject = function(includeInstance, msg) {
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
 * @return {!proto.envoy.config.route.v3.RouteMatch.GrpcRouteMatchOptions}
 */
proto.envoy.config.route.v3.RouteMatch.GrpcRouteMatchOptions.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.envoy.config.route.v3.RouteMatch.GrpcRouteMatchOptions;
  return proto.envoy.config.route.v3.RouteMatch.GrpcRouteMatchOptions.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.envoy.config.route.v3.RouteMatch.GrpcRouteMatchOptions} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.envoy.config.route.v3.RouteMatch.GrpcRouteMatchOptions}
 */
proto.envoy.config.route.v3.RouteMatch.GrpcRouteMatchOptions.deserializeBinaryFromReader = function(msg, reader) {
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
proto.envoy.config.route.v3.RouteMatch.GrpcRouteMatchOptions.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.envoy.config.route.v3.RouteMatch.GrpcRouteMatchOptions.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.envoy.config.route.v3.RouteMatch.GrpcRouteMatchOptions} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.RouteMatch.GrpcRouteMatchOptions.serializeBinaryToWriter = function(message, writer) {
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
proto.envoy.config.route.v3.RouteMatch.TlsContextMatchOptions = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.envoy.config.route.v3.RouteMatch.TlsContextMatchOptions, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.envoy.config.route.v3.RouteMatch.TlsContextMatchOptions.displayName = 'proto.envoy.config.route.v3.RouteMatch.TlsContextMatchOptions';
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
proto.envoy.config.route.v3.RouteMatch.TlsContextMatchOptions.prototype.toObject = function(opt_includeInstance) {
  return proto.envoy.config.route.v3.RouteMatch.TlsContextMatchOptions.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.envoy.config.route.v3.RouteMatch.TlsContextMatchOptions} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.RouteMatch.TlsContextMatchOptions.toObject = function(includeInstance, msg) {
  var f, obj = {
    presented: (f = msg.getPresented()) && google_protobuf_wrappers_pb.BoolValue.toObject(includeInstance, f),
    validated: (f = msg.getValidated()) && google_protobuf_wrappers_pb.BoolValue.toObject(includeInstance, f)
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
 * @return {!proto.envoy.config.route.v3.RouteMatch.TlsContextMatchOptions}
 */
proto.envoy.config.route.v3.RouteMatch.TlsContextMatchOptions.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.envoy.config.route.v3.RouteMatch.TlsContextMatchOptions;
  return proto.envoy.config.route.v3.RouteMatch.TlsContextMatchOptions.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.envoy.config.route.v3.RouteMatch.TlsContextMatchOptions} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.envoy.config.route.v3.RouteMatch.TlsContextMatchOptions}
 */
proto.envoy.config.route.v3.RouteMatch.TlsContextMatchOptions.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new google_protobuf_wrappers_pb.BoolValue;
      reader.readMessage(value,google_protobuf_wrappers_pb.BoolValue.deserializeBinaryFromReader);
      msg.setPresented(value);
      break;
    case 2:
      var value = new google_protobuf_wrappers_pb.BoolValue;
      reader.readMessage(value,google_protobuf_wrappers_pb.BoolValue.deserializeBinaryFromReader);
      msg.setValidated(value);
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
proto.envoy.config.route.v3.RouteMatch.TlsContextMatchOptions.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.envoy.config.route.v3.RouteMatch.TlsContextMatchOptions.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.envoy.config.route.v3.RouteMatch.TlsContextMatchOptions} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.RouteMatch.TlsContextMatchOptions.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getPresented();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      google_protobuf_wrappers_pb.BoolValue.serializeBinaryToWriter
    );
  }
  f = message.getValidated();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      google_protobuf_wrappers_pb.BoolValue.serializeBinaryToWriter
    );
  }
};


/**
 * optional google.protobuf.BoolValue presented = 1;
 * @return {?proto.google.protobuf.BoolValue}
 */
proto.envoy.config.route.v3.RouteMatch.TlsContextMatchOptions.prototype.getPresented = function() {
  return /** @type{?proto.google.protobuf.BoolValue} */ (
    jspb.Message.getWrapperField(this, google_protobuf_wrappers_pb.BoolValue, 1));
};


/** @param {?proto.google.protobuf.BoolValue|undefined} value */
proto.envoy.config.route.v3.RouteMatch.TlsContextMatchOptions.prototype.setPresented = function(value) {
  jspb.Message.setWrapperField(this, 1, value);
};


proto.envoy.config.route.v3.RouteMatch.TlsContextMatchOptions.prototype.clearPresented = function() {
  this.setPresented(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.RouteMatch.TlsContextMatchOptions.prototype.hasPresented = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional google.protobuf.BoolValue validated = 2;
 * @return {?proto.google.protobuf.BoolValue}
 */
proto.envoy.config.route.v3.RouteMatch.TlsContextMatchOptions.prototype.getValidated = function() {
  return /** @type{?proto.google.protobuf.BoolValue} */ (
    jspb.Message.getWrapperField(this, google_protobuf_wrappers_pb.BoolValue, 2));
};


/** @param {?proto.google.protobuf.BoolValue|undefined} value */
proto.envoy.config.route.v3.RouteMatch.TlsContextMatchOptions.prototype.setValidated = function(value) {
  jspb.Message.setWrapperField(this, 2, value);
};


proto.envoy.config.route.v3.RouteMatch.TlsContextMatchOptions.prototype.clearValidated = function() {
  this.setValidated(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.RouteMatch.TlsContextMatchOptions.prototype.hasValidated = function() {
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
proto.envoy.config.route.v3.RouteMatch.ConnectMatcher = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.envoy.config.route.v3.RouteMatch.ConnectMatcher, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.envoy.config.route.v3.RouteMatch.ConnectMatcher.displayName = 'proto.envoy.config.route.v3.RouteMatch.ConnectMatcher';
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
proto.envoy.config.route.v3.RouteMatch.ConnectMatcher.prototype.toObject = function(opt_includeInstance) {
  return proto.envoy.config.route.v3.RouteMatch.ConnectMatcher.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.envoy.config.route.v3.RouteMatch.ConnectMatcher} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.RouteMatch.ConnectMatcher.toObject = function(includeInstance, msg) {
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
 * @return {!proto.envoy.config.route.v3.RouteMatch.ConnectMatcher}
 */
proto.envoy.config.route.v3.RouteMatch.ConnectMatcher.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.envoy.config.route.v3.RouteMatch.ConnectMatcher;
  return proto.envoy.config.route.v3.RouteMatch.ConnectMatcher.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.envoy.config.route.v3.RouteMatch.ConnectMatcher} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.envoy.config.route.v3.RouteMatch.ConnectMatcher}
 */
proto.envoy.config.route.v3.RouteMatch.ConnectMatcher.deserializeBinaryFromReader = function(msg, reader) {
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
proto.envoy.config.route.v3.RouteMatch.ConnectMatcher.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.envoy.config.route.v3.RouteMatch.ConnectMatcher.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.envoy.config.route.v3.RouteMatch.ConnectMatcher} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.RouteMatch.ConnectMatcher.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
};


/**
 * optional string prefix = 1;
 * @return {string}
 */
proto.envoy.config.route.v3.RouteMatch.prototype.getPrefix = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.envoy.config.route.v3.RouteMatch.prototype.setPrefix = function(value) {
  jspb.Message.setOneofField(this, 1, proto.envoy.config.route.v3.RouteMatch.oneofGroups_[0], value);
};


proto.envoy.config.route.v3.RouteMatch.prototype.clearPrefix = function() {
  jspb.Message.setOneofField(this, 1, proto.envoy.config.route.v3.RouteMatch.oneofGroups_[0], undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.RouteMatch.prototype.hasPrefix = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional string path = 2;
 * @return {string}
 */
proto.envoy.config.route.v3.RouteMatch.prototype.getPath = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.envoy.config.route.v3.RouteMatch.prototype.setPath = function(value) {
  jspb.Message.setOneofField(this, 2, proto.envoy.config.route.v3.RouteMatch.oneofGroups_[0], value);
};


proto.envoy.config.route.v3.RouteMatch.prototype.clearPath = function() {
  jspb.Message.setOneofField(this, 2, proto.envoy.config.route.v3.RouteMatch.oneofGroups_[0], undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.RouteMatch.prototype.hasPath = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional envoy.type.matcher.v3.RegexMatcher safe_regex = 10;
 * @return {?proto.envoy.type.matcher.v3.RegexMatcher}
 */
proto.envoy.config.route.v3.RouteMatch.prototype.getSafeRegex = function() {
  return /** @type{?proto.envoy.type.matcher.v3.RegexMatcher} */ (
    jspb.Message.getWrapperField(this, envoy_type_matcher_v3_regex_pb.RegexMatcher, 10));
};


/** @param {?proto.envoy.type.matcher.v3.RegexMatcher|undefined} value */
proto.envoy.config.route.v3.RouteMatch.prototype.setSafeRegex = function(value) {
  jspb.Message.setOneofWrapperField(this, 10, proto.envoy.config.route.v3.RouteMatch.oneofGroups_[0], value);
};


proto.envoy.config.route.v3.RouteMatch.prototype.clearSafeRegex = function() {
  this.setSafeRegex(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.RouteMatch.prototype.hasSafeRegex = function() {
  return jspb.Message.getField(this, 10) != null;
};


/**
 * optional ConnectMatcher connect_matcher = 12;
 * @return {?proto.envoy.config.route.v3.RouteMatch.ConnectMatcher}
 */
proto.envoy.config.route.v3.RouteMatch.prototype.getConnectMatcher = function() {
  return /** @type{?proto.envoy.config.route.v3.RouteMatch.ConnectMatcher} */ (
    jspb.Message.getWrapperField(this, proto.envoy.config.route.v3.RouteMatch.ConnectMatcher, 12));
};


/** @param {?proto.envoy.config.route.v3.RouteMatch.ConnectMatcher|undefined} value */
proto.envoy.config.route.v3.RouteMatch.prototype.setConnectMatcher = function(value) {
  jspb.Message.setOneofWrapperField(this, 12, proto.envoy.config.route.v3.RouteMatch.oneofGroups_[0], value);
};


proto.envoy.config.route.v3.RouteMatch.prototype.clearConnectMatcher = function() {
  this.setConnectMatcher(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.RouteMatch.prototype.hasConnectMatcher = function() {
  return jspb.Message.getField(this, 12) != null;
};


/**
 * optional google.protobuf.BoolValue case_sensitive = 4;
 * @return {?proto.google.protobuf.BoolValue}
 */
proto.envoy.config.route.v3.RouteMatch.prototype.getCaseSensitive = function() {
  return /** @type{?proto.google.protobuf.BoolValue} */ (
    jspb.Message.getWrapperField(this, google_protobuf_wrappers_pb.BoolValue, 4));
};


/** @param {?proto.google.protobuf.BoolValue|undefined} value */
proto.envoy.config.route.v3.RouteMatch.prototype.setCaseSensitive = function(value) {
  jspb.Message.setWrapperField(this, 4, value);
};


proto.envoy.config.route.v3.RouteMatch.prototype.clearCaseSensitive = function() {
  this.setCaseSensitive(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.RouteMatch.prototype.hasCaseSensitive = function() {
  return jspb.Message.getField(this, 4) != null;
};


/**
 * optional envoy.config.core.v3.RuntimeFractionalPercent runtime_fraction = 9;
 * @return {?proto.envoy.config.core.v3.RuntimeFractionalPercent}
 */
proto.envoy.config.route.v3.RouteMatch.prototype.getRuntimeFraction = function() {
  return /** @type{?proto.envoy.config.core.v3.RuntimeFractionalPercent} */ (
    jspb.Message.getWrapperField(this, envoy_config_core_v3_base_pb.RuntimeFractionalPercent, 9));
};


/** @param {?proto.envoy.config.core.v3.RuntimeFractionalPercent|undefined} value */
proto.envoy.config.route.v3.RouteMatch.prototype.setRuntimeFraction = function(value) {
  jspb.Message.setWrapperField(this, 9, value);
};


proto.envoy.config.route.v3.RouteMatch.prototype.clearRuntimeFraction = function() {
  this.setRuntimeFraction(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.RouteMatch.prototype.hasRuntimeFraction = function() {
  return jspb.Message.getField(this, 9) != null;
};


/**
 * repeated HeaderMatcher headers = 6;
 * @return {!Array<!proto.envoy.config.route.v3.HeaderMatcher>}
 */
proto.envoy.config.route.v3.RouteMatch.prototype.getHeadersList = function() {
  return /** @type{!Array<!proto.envoy.config.route.v3.HeaderMatcher>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.envoy.config.route.v3.HeaderMatcher, 6));
};


/** @param {!Array<!proto.envoy.config.route.v3.HeaderMatcher>} value */
proto.envoy.config.route.v3.RouteMatch.prototype.setHeadersList = function(value) {
  jspb.Message.setRepeatedWrapperField(this, 6, value);
};


/**
 * @param {!proto.envoy.config.route.v3.HeaderMatcher=} opt_value
 * @param {number=} opt_index
 * @return {!proto.envoy.config.route.v3.HeaderMatcher}
 */
proto.envoy.config.route.v3.RouteMatch.prototype.addHeaders = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 6, opt_value, proto.envoy.config.route.v3.HeaderMatcher, opt_index);
};


proto.envoy.config.route.v3.RouteMatch.prototype.clearHeadersList = function() {
  this.setHeadersList([]);
};


/**
 * repeated QueryParameterMatcher query_parameters = 7;
 * @return {!Array<!proto.envoy.config.route.v3.QueryParameterMatcher>}
 */
proto.envoy.config.route.v3.RouteMatch.prototype.getQueryParametersList = function() {
  return /** @type{!Array<!proto.envoy.config.route.v3.QueryParameterMatcher>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.envoy.config.route.v3.QueryParameterMatcher, 7));
};


/** @param {!Array<!proto.envoy.config.route.v3.QueryParameterMatcher>} value */
proto.envoy.config.route.v3.RouteMatch.prototype.setQueryParametersList = function(value) {
  jspb.Message.setRepeatedWrapperField(this, 7, value);
};


/**
 * @param {!proto.envoy.config.route.v3.QueryParameterMatcher=} opt_value
 * @param {number=} opt_index
 * @return {!proto.envoy.config.route.v3.QueryParameterMatcher}
 */
proto.envoy.config.route.v3.RouteMatch.prototype.addQueryParameters = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 7, opt_value, proto.envoy.config.route.v3.QueryParameterMatcher, opt_index);
};


proto.envoy.config.route.v3.RouteMatch.prototype.clearQueryParametersList = function() {
  this.setQueryParametersList([]);
};


/**
 * optional GrpcRouteMatchOptions grpc = 8;
 * @return {?proto.envoy.config.route.v3.RouteMatch.GrpcRouteMatchOptions}
 */
proto.envoy.config.route.v3.RouteMatch.prototype.getGrpc = function() {
  return /** @type{?proto.envoy.config.route.v3.RouteMatch.GrpcRouteMatchOptions} */ (
    jspb.Message.getWrapperField(this, proto.envoy.config.route.v3.RouteMatch.GrpcRouteMatchOptions, 8));
};


/** @param {?proto.envoy.config.route.v3.RouteMatch.GrpcRouteMatchOptions|undefined} value */
proto.envoy.config.route.v3.RouteMatch.prototype.setGrpc = function(value) {
  jspb.Message.setWrapperField(this, 8, value);
};


proto.envoy.config.route.v3.RouteMatch.prototype.clearGrpc = function() {
  this.setGrpc(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.RouteMatch.prototype.hasGrpc = function() {
  return jspb.Message.getField(this, 8) != null;
};


/**
 * optional TlsContextMatchOptions tls_context = 11;
 * @return {?proto.envoy.config.route.v3.RouteMatch.TlsContextMatchOptions}
 */
proto.envoy.config.route.v3.RouteMatch.prototype.getTlsContext = function() {
  return /** @type{?proto.envoy.config.route.v3.RouteMatch.TlsContextMatchOptions} */ (
    jspb.Message.getWrapperField(this, proto.envoy.config.route.v3.RouteMatch.TlsContextMatchOptions, 11));
};


/** @param {?proto.envoy.config.route.v3.RouteMatch.TlsContextMatchOptions|undefined} value */
proto.envoy.config.route.v3.RouteMatch.prototype.setTlsContext = function(value) {
  jspb.Message.setWrapperField(this, 11, value);
};


proto.envoy.config.route.v3.RouteMatch.prototype.clearTlsContext = function() {
  this.setTlsContext(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.RouteMatch.prototype.hasTlsContext = function() {
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
proto.envoy.config.route.v3.CorsPolicy = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.envoy.config.route.v3.CorsPolicy.repeatedFields_, proto.envoy.config.route.v3.CorsPolicy.oneofGroups_);
};
goog.inherits(proto.envoy.config.route.v3.CorsPolicy, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.envoy.config.route.v3.CorsPolicy.displayName = 'proto.envoy.config.route.v3.CorsPolicy';
}
/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.envoy.config.route.v3.CorsPolicy.repeatedFields_ = [11];

/**
 * Oneof group definitions for this message. Each group defines the field
 * numbers belonging to that group. When of these fields' value is set, all
 * other fields in the group are cleared. During deserialization, if multiple
 * fields are encountered for a group, only the last value seen will be kept.
 * @private {!Array<!Array<number>>}
 * @const
 */
proto.envoy.config.route.v3.CorsPolicy.oneofGroups_ = [[9]];

/**
 * @enum {number}
 */
proto.envoy.config.route.v3.CorsPolicy.EnabledSpecifierCase = {
  ENABLED_SPECIFIER_NOT_SET: 0,
  FILTER_ENABLED: 9
};

/**
 * @return {proto.envoy.config.route.v3.CorsPolicy.EnabledSpecifierCase}
 */
proto.envoy.config.route.v3.CorsPolicy.prototype.getEnabledSpecifierCase = function() {
  return /** @type {proto.envoy.config.route.v3.CorsPolicy.EnabledSpecifierCase} */(jspb.Message.computeOneofCase(this, proto.envoy.config.route.v3.CorsPolicy.oneofGroups_[0]));
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
proto.envoy.config.route.v3.CorsPolicy.prototype.toObject = function(opt_includeInstance) {
  return proto.envoy.config.route.v3.CorsPolicy.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.envoy.config.route.v3.CorsPolicy} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.CorsPolicy.toObject = function(includeInstance, msg) {
  var f, obj = {
    allowOriginStringMatchList: jspb.Message.toObjectList(msg.getAllowOriginStringMatchList(),
    envoy_type_matcher_v3_string_pb.StringMatcher.toObject, includeInstance),
    allowMethods: jspb.Message.getFieldWithDefault(msg, 2, ""),
    allowHeaders: jspb.Message.getFieldWithDefault(msg, 3, ""),
    exposeHeaders: jspb.Message.getFieldWithDefault(msg, 4, ""),
    maxAge: jspb.Message.getFieldWithDefault(msg, 5, ""),
    allowCredentials: (f = msg.getAllowCredentials()) && google_protobuf_wrappers_pb.BoolValue.toObject(includeInstance, f),
    filterEnabled: (f = msg.getFilterEnabled()) && envoy_config_core_v3_base_pb.RuntimeFractionalPercent.toObject(includeInstance, f),
    shadowEnabled: (f = msg.getShadowEnabled()) && envoy_config_core_v3_base_pb.RuntimeFractionalPercent.toObject(includeInstance, f)
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
 * @return {!proto.envoy.config.route.v3.CorsPolicy}
 */
proto.envoy.config.route.v3.CorsPolicy.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.envoy.config.route.v3.CorsPolicy;
  return proto.envoy.config.route.v3.CorsPolicy.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.envoy.config.route.v3.CorsPolicy} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.envoy.config.route.v3.CorsPolicy}
 */
proto.envoy.config.route.v3.CorsPolicy.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 11:
      var value = new envoy_type_matcher_v3_string_pb.StringMatcher;
      reader.readMessage(value,envoy_type_matcher_v3_string_pb.StringMatcher.deserializeBinaryFromReader);
      msg.addAllowOriginStringMatch(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setAllowMethods(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setAllowHeaders(value);
      break;
    case 4:
      var value = /** @type {string} */ (reader.readString());
      msg.setExposeHeaders(value);
      break;
    case 5:
      var value = /** @type {string} */ (reader.readString());
      msg.setMaxAge(value);
      break;
    case 6:
      var value = new google_protobuf_wrappers_pb.BoolValue;
      reader.readMessage(value,google_protobuf_wrappers_pb.BoolValue.deserializeBinaryFromReader);
      msg.setAllowCredentials(value);
      break;
    case 9:
      var value = new envoy_config_core_v3_base_pb.RuntimeFractionalPercent;
      reader.readMessage(value,envoy_config_core_v3_base_pb.RuntimeFractionalPercent.deserializeBinaryFromReader);
      msg.setFilterEnabled(value);
      break;
    case 10:
      var value = new envoy_config_core_v3_base_pb.RuntimeFractionalPercent;
      reader.readMessage(value,envoy_config_core_v3_base_pb.RuntimeFractionalPercent.deserializeBinaryFromReader);
      msg.setShadowEnabled(value);
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
proto.envoy.config.route.v3.CorsPolicy.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.envoy.config.route.v3.CorsPolicy.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.envoy.config.route.v3.CorsPolicy} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.CorsPolicy.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getAllowOriginStringMatchList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      11,
      f,
      envoy_type_matcher_v3_string_pb.StringMatcher.serializeBinaryToWriter
    );
  }
  f = message.getAllowMethods();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getAllowHeaders();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getExposeHeaders();
  if (f.length > 0) {
    writer.writeString(
      4,
      f
    );
  }
  f = message.getMaxAge();
  if (f.length > 0) {
    writer.writeString(
      5,
      f
    );
  }
  f = message.getAllowCredentials();
  if (f != null) {
    writer.writeMessage(
      6,
      f,
      google_protobuf_wrappers_pb.BoolValue.serializeBinaryToWriter
    );
  }
  f = message.getFilterEnabled();
  if (f != null) {
    writer.writeMessage(
      9,
      f,
      envoy_config_core_v3_base_pb.RuntimeFractionalPercent.serializeBinaryToWriter
    );
  }
  f = message.getShadowEnabled();
  if (f != null) {
    writer.writeMessage(
      10,
      f,
      envoy_config_core_v3_base_pb.RuntimeFractionalPercent.serializeBinaryToWriter
    );
  }
};


/**
 * repeated envoy.type.matcher.v3.StringMatcher allow_origin_string_match = 11;
 * @return {!Array<!proto.envoy.type.matcher.v3.StringMatcher>}
 */
proto.envoy.config.route.v3.CorsPolicy.prototype.getAllowOriginStringMatchList = function() {
  return /** @type{!Array<!proto.envoy.type.matcher.v3.StringMatcher>} */ (
    jspb.Message.getRepeatedWrapperField(this, envoy_type_matcher_v3_string_pb.StringMatcher, 11));
};


/** @param {!Array<!proto.envoy.type.matcher.v3.StringMatcher>} value */
proto.envoy.config.route.v3.CorsPolicy.prototype.setAllowOriginStringMatchList = function(value) {
  jspb.Message.setRepeatedWrapperField(this, 11, value);
};


/**
 * @param {!proto.envoy.type.matcher.v3.StringMatcher=} opt_value
 * @param {number=} opt_index
 * @return {!proto.envoy.type.matcher.v3.StringMatcher}
 */
proto.envoy.config.route.v3.CorsPolicy.prototype.addAllowOriginStringMatch = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 11, opt_value, proto.envoy.type.matcher.v3.StringMatcher, opt_index);
};


proto.envoy.config.route.v3.CorsPolicy.prototype.clearAllowOriginStringMatchList = function() {
  this.setAllowOriginStringMatchList([]);
};


/**
 * optional string allow_methods = 2;
 * @return {string}
 */
proto.envoy.config.route.v3.CorsPolicy.prototype.getAllowMethods = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.envoy.config.route.v3.CorsPolicy.prototype.setAllowMethods = function(value) {
  jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string allow_headers = 3;
 * @return {string}
 */
proto.envoy.config.route.v3.CorsPolicy.prototype.getAllowHeaders = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/** @param {string} value */
proto.envoy.config.route.v3.CorsPolicy.prototype.setAllowHeaders = function(value) {
  jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional string expose_headers = 4;
 * @return {string}
 */
proto.envoy.config.route.v3.CorsPolicy.prototype.getExposeHeaders = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 4, ""));
};


/** @param {string} value */
proto.envoy.config.route.v3.CorsPolicy.prototype.setExposeHeaders = function(value) {
  jspb.Message.setProto3StringField(this, 4, value);
};


/**
 * optional string max_age = 5;
 * @return {string}
 */
proto.envoy.config.route.v3.CorsPolicy.prototype.getMaxAge = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 5, ""));
};


/** @param {string} value */
proto.envoy.config.route.v3.CorsPolicy.prototype.setMaxAge = function(value) {
  jspb.Message.setProto3StringField(this, 5, value);
};


/**
 * optional google.protobuf.BoolValue allow_credentials = 6;
 * @return {?proto.google.protobuf.BoolValue}
 */
proto.envoy.config.route.v3.CorsPolicy.prototype.getAllowCredentials = function() {
  return /** @type{?proto.google.protobuf.BoolValue} */ (
    jspb.Message.getWrapperField(this, google_protobuf_wrappers_pb.BoolValue, 6));
};


/** @param {?proto.google.protobuf.BoolValue|undefined} value */
proto.envoy.config.route.v3.CorsPolicy.prototype.setAllowCredentials = function(value) {
  jspb.Message.setWrapperField(this, 6, value);
};


proto.envoy.config.route.v3.CorsPolicy.prototype.clearAllowCredentials = function() {
  this.setAllowCredentials(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.CorsPolicy.prototype.hasAllowCredentials = function() {
  return jspb.Message.getField(this, 6) != null;
};


/**
 * optional envoy.config.core.v3.RuntimeFractionalPercent filter_enabled = 9;
 * @return {?proto.envoy.config.core.v3.RuntimeFractionalPercent}
 */
proto.envoy.config.route.v3.CorsPolicy.prototype.getFilterEnabled = function() {
  return /** @type{?proto.envoy.config.core.v3.RuntimeFractionalPercent} */ (
    jspb.Message.getWrapperField(this, envoy_config_core_v3_base_pb.RuntimeFractionalPercent, 9));
};


/** @param {?proto.envoy.config.core.v3.RuntimeFractionalPercent|undefined} value */
proto.envoy.config.route.v3.CorsPolicy.prototype.setFilterEnabled = function(value) {
  jspb.Message.setOneofWrapperField(this, 9, proto.envoy.config.route.v3.CorsPolicy.oneofGroups_[0], value);
};


proto.envoy.config.route.v3.CorsPolicy.prototype.clearFilterEnabled = function() {
  this.setFilterEnabled(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.CorsPolicy.prototype.hasFilterEnabled = function() {
  return jspb.Message.getField(this, 9) != null;
};


/**
 * optional envoy.config.core.v3.RuntimeFractionalPercent shadow_enabled = 10;
 * @return {?proto.envoy.config.core.v3.RuntimeFractionalPercent}
 */
proto.envoy.config.route.v3.CorsPolicy.prototype.getShadowEnabled = function() {
  return /** @type{?proto.envoy.config.core.v3.RuntimeFractionalPercent} */ (
    jspb.Message.getWrapperField(this, envoy_config_core_v3_base_pb.RuntimeFractionalPercent, 10));
};


/** @param {?proto.envoy.config.core.v3.RuntimeFractionalPercent|undefined} value */
proto.envoy.config.route.v3.CorsPolicy.prototype.setShadowEnabled = function(value) {
  jspb.Message.setWrapperField(this, 10, value);
};


proto.envoy.config.route.v3.CorsPolicy.prototype.clearShadowEnabled = function() {
  this.setShadowEnabled(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.CorsPolicy.prototype.hasShadowEnabled = function() {
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
proto.envoy.config.route.v3.RouteAction = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.envoy.config.route.v3.RouteAction.repeatedFields_, proto.envoy.config.route.v3.RouteAction.oneofGroups_);
};
goog.inherits(proto.envoy.config.route.v3.RouteAction, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.envoy.config.route.v3.RouteAction.displayName = 'proto.envoy.config.route.v3.RouteAction';
}
/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.envoy.config.route.v3.RouteAction.repeatedFields_ = [30,13,15,25];

/**
 * Oneof group definitions for this message. Each group defines the field
 * numbers belonging to that group. When of these fields' value is set, all
 * other fields in the group are cleared. During deserialization, if multiple
 * fields are encountered for a group, only the last value seen will be kept.
 * @private {!Array<!Array<number>>}
 * @const
 */
proto.envoy.config.route.v3.RouteAction.oneofGroups_ = [[1,2,3],[6,7,29]];

/**
 * @enum {number}
 */
proto.envoy.config.route.v3.RouteAction.ClusterSpecifierCase = {
  CLUSTER_SPECIFIER_NOT_SET: 0,
  CLUSTER: 1,
  CLUSTER_HEADER: 2,
  WEIGHTED_CLUSTERS: 3
};

/**
 * @return {proto.envoy.config.route.v3.RouteAction.ClusterSpecifierCase}
 */
proto.envoy.config.route.v3.RouteAction.prototype.getClusterSpecifierCase = function() {
  return /** @type {proto.envoy.config.route.v3.RouteAction.ClusterSpecifierCase} */(jspb.Message.computeOneofCase(this, proto.envoy.config.route.v3.RouteAction.oneofGroups_[0]));
};

/**
 * @enum {number}
 */
proto.envoy.config.route.v3.RouteAction.HostRewriteSpecifierCase = {
  HOST_REWRITE_SPECIFIER_NOT_SET: 0,
  HOST_REWRITE_LITERAL: 6,
  AUTO_HOST_REWRITE: 7,
  HOST_REWRITE_HEADER: 29
};

/**
 * @return {proto.envoy.config.route.v3.RouteAction.HostRewriteSpecifierCase}
 */
proto.envoy.config.route.v3.RouteAction.prototype.getHostRewriteSpecifierCase = function() {
  return /** @type {proto.envoy.config.route.v3.RouteAction.HostRewriteSpecifierCase} */(jspb.Message.computeOneofCase(this, proto.envoy.config.route.v3.RouteAction.oneofGroups_[1]));
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
proto.envoy.config.route.v3.RouteAction.prototype.toObject = function(opt_includeInstance) {
  return proto.envoy.config.route.v3.RouteAction.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.envoy.config.route.v3.RouteAction} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.RouteAction.toObject = function(includeInstance, msg) {
  var f, obj = {
    cluster: jspb.Message.getFieldWithDefault(msg, 1, ""),
    clusterHeader: jspb.Message.getFieldWithDefault(msg, 2, ""),
    weightedClusters: (f = msg.getWeightedClusters()) && proto.envoy.config.route.v3.WeightedCluster.toObject(includeInstance, f),
    clusterNotFoundResponseCode: jspb.Message.getFieldWithDefault(msg, 20, 0),
    metadataMatch: (f = msg.getMetadataMatch()) && envoy_config_core_v3_base_pb.Metadata.toObject(includeInstance, f),
    prefixRewrite: jspb.Message.getFieldWithDefault(msg, 5, ""),
    regexRewrite: (f = msg.getRegexRewrite()) && envoy_type_matcher_v3_regex_pb.RegexMatchAndSubstitute.toObject(includeInstance, f),
    hostRewriteLiteral: jspb.Message.getFieldWithDefault(msg, 6, ""),
    autoHostRewrite: (f = msg.getAutoHostRewrite()) && google_protobuf_wrappers_pb.BoolValue.toObject(includeInstance, f),
    hostRewriteHeader: jspb.Message.getFieldWithDefault(msg, 29, ""),
    timeout: (f = msg.getTimeout()) && google_protobuf_duration_pb.Duration.toObject(includeInstance, f),
    idleTimeout: (f = msg.getIdleTimeout()) && google_protobuf_duration_pb.Duration.toObject(includeInstance, f),
    retryPolicy: (f = msg.getRetryPolicy()) && proto.envoy.config.route.v3.RetryPolicy.toObject(includeInstance, f),
    retryPolicyTypedConfig: (f = msg.getRetryPolicyTypedConfig()) && google_protobuf_any_pb.Any.toObject(includeInstance, f),
    requestMirrorPoliciesList: jspb.Message.toObjectList(msg.getRequestMirrorPoliciesList(),
    proto.envoy.config.route.v3.RouteAction.RequestMirrorPolicy.toObject, includeInstance),
    priority: jspb.Message.getFieldWithDefault(msg, 11, 0),
    rateLimitsList: jspb.Message.toObjectList(msg.getRateLimitsList(),
    proto.envoy.config.route.v3.RateLimit.toObject, includeInstance),
    includeVhRateLimits: (f = msg.getIncludeVhRateLimits()) && google_protobuf_wrappers_pb.BoolValue.toObject(includeInstance, f),
    hashPolicyList: jspb.Message.toObjectList(msg.getHashPolicyList(),
    proto.envoy.config.route.v3.RouteAction.HashPolicy.toObject, includeInstance),
    cors: (f = msg.getCors()) && proto.envoy.config.route.v3.CorsPolicy.toObject(includeInstance, f),
    maxGrpcTimeout: (f = msg.getMaxGrpcTimeout()) && google_protobuf_duration_pb.Duration.toObject(includeInstance, f),
    grpcTimeoutOffset: (f = msg.getGrpcTimeoutOffset()) && google_protobuf_duration_pb.Duration.toObject(includeInstance, f),
    upgradeConfigsList: jspb.Message.toObjectList(msg.getUpgradeConfigsList(),
    proto.envoy.config.route.v3.RouteAction.UpgradeConfig.toObject, includeInstance),
    internalRedirectPolicy: (f = msg.getInternalRedirectPolicy()) && proto.envoy.config.route.v3.InternalRedirectPolicy.toObject(includeInstance, f),
    internalRedirectAction: jspb.Message.getFieldWithDefault(msg, 26, 0),
    maxInternalRedirects: (f = msg.getMaxInternalRedirects()) && google_protobuf_wrappers_pb.UInt32Value.toObject(includeInstance, f),
    hedgePolicy: (f = msg.getHedgePolicy()) && proto.envoy.config.route.v3.HedgePolicy.toObject(includeInstance, f)
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
 * @return {!proto.envoy.config.route.v3.RouteAction}
 */
proto.envoy.config.route.v3.RouteAction.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.envoy.config.route.v3.RouteAction;
  return proto.envoy.config.route.v3.RouteAction.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.envoy.config.route.v3.RouteAction} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.envoy.config.route.v3.RouteAction}
 */
proto.envoy.config.route.v3.RouteAction.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setCluster(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setClusterHeader(value);
      break;
    case 3:
      var value = new proto.envoy.config.route.v3.WeightedCluster;
      reader.readMessage(value,proto.envoy.config.route.v3.WeightedCluster.deserializeBinaryFromReader);
      msg.setWeightedClusters(value);
      break;
    case 20:
      var value = /** @type {!proto.envoy.config.route.v3.RouteAction.ClusterNotFoundResponseCode} */ (reader.readEnum());
      msg.setClusterNotFoundResponseCode(value);
      break;
    case 4:
      var value = new envoy_config_core_v3_base_pb.Metadata;
      reader.readMessage(value,envoy_config_core_v3_base_pb.Metadata.deserializeBinaryFromReader);
      msg.setMetadataMatch(value);
      break;
    case 5:
      var value = /** @type {string} */ (reader.readString());
      msg.setPrefixRewrite(value);
      break;
    case 32:
      var value = new envoy_type_matcher_v3_regex_pb.RegexMatchAndSubstitute;
      reader.readMessage(value,envoy_type_matcher_v3_regex_pb.RegexMatchAndSubstitute.deserializeBinaryFromReader);
      msg.setRegexRewrite(value);
      break;
    case 6:
      var value = /** @type {string} */ (reader.readString());
      msg.setHostRewriteLiteral(value);
      break;
    case 7:
      var value = new google_protobuf_wrappers_pb.BoolValue;
      reader.readMessage(value,google_protobuf_wrappers_pb.BoolValue.deserializeBinaryFromReader);
      msg.setAutoHostRewrite(value);
      break;
    case 29:
      var value = /** @type {string} */ (reader.readString());
      msg.setHostRewriteHeader(value);
      break;
    case 8:
      var value = new google_protobuf_duration_pb.Duration;
      reader.readMessage(value,google_protobuf_duration_pb.Duration.deserializeBinaryFromReader);
      msg.setTimeout(value);
      break;
    case 24:
      var value = new google_protobuf_duration_pb.Duration;
      reader.readMessage(value,google_protobuf_duration_pb.Duration.deserializeBinaryFromReader);
      msg.setIdleTimeout(value);
      break;
    case 9:
      var value = new proto.envoy.config.route.v3.RetryPolicy;
      reader.readMessage(value,proto.envoy.config.route.v3.RetryPolicy.deserializeBinaryFromReader);
      msg.setRetryPolicy(value);
      break;
    case 33:
      var value = new google_protobuf_any_pb.Any;
      reader.readMessage(value,google_protobuf_any_pb.Any.deserializeBinaryFromReader);
      msg.setRetryPolicyTypedConfig(value);
      break;
    case 30:
      var value = new proto.envoy.config.route.v3.RouteAction.RequestMirrorPolicy;
      reader.readMessage(value,proto.envoy.config.route.v3.RouteAction.RequestMirrorPolicy.deserializeBinaryFromReader);
      msg.addRequestMirrorPolicies(value);
      break;
    case 11:
      var value = /** @type {!proto.envoy.config.core.v3.RoutingPriority} */ (reader.readEnum());
      msg.setPriority(value);
      break;
    case 13:
      var value = new proto.envoy.config.route.v3.RateLimit;
      reader.readMessage(value,proto.envoy.config.route.v3.RateLimit.deserializeBinaryFromReader);
      msg.addRateLimits(value);
      break;
    case 14:
      var value = new google_protobuf_wrappers_pb.BoolValue;
      reader.readMessage(value,google_protobuf_wrappers_pb.BoolValue.deserializeBinaryFromReader);
      msg.setIncludeVhRateLimits(value);
      break;
    case 15:
      var value = new proto.envoy.config.route.v3.RouteAction.HashPolicy;
      reader.readMessage(value,proto.envoy.config.route.v3.RouteAction.HashPolicy.deserializeBinaryFromReader);
      msg.addHashPolicy(value);
      break;
    case 17:
      var value = new proto.envoy.config.route.v3.CorsPolicy;
      reader.readMessage(value,proto.envoy.config.route.v3.CorsPolicy.deserializeBinaryFromReader);
      msg.setCors(value);
      break;
    case 23:
      var value = new google_protobuf_duration_pb.Duration;
      reader.readMessage(value,google_protobuf_duration_pb.Duration.deserializeBinaryFromReader);
      msg.setMaxGrpcTimeout(value);
      break;
    case 28:
      var value = new google_protobuf_duration_pb.Duration;
      reader.readMessage(value,google_protobuf_duration_pb.Duration.deserializeBinaryFromReader);
      msg.setGrpcTimeoutOffset(value);
      break;
    case 25:
      var value = new proto.envoy.config.route.v3.RouteAction.UpgradeConfig;
      reader.readMessage(value,proto.envoy.config.route.v3.RouteAction.UpgradeConfig.deserializeBinaryFromReader);
      msg.addUpgradeConfigs(value);
      break;
    case 34:
      var value = new proto.envoy.config.route.v3.InternalRedirectPolicy;
      reader.readMessage(value,proto.envoy.config.route.v3.InternalRedirectPolicy.deserializeBinaryFromReader);
      msg.setInternalRedirectPolicy(value);
      break;
    case 26:
      var value = /** @type {!proto.envoy.config.route.v3.RouteAction.InternalRedirectAction} */ (reader.readEnum());
      msg.setInternalRedirectAction(value);
      break;
    case 31:
      var value = new google_protobuf_wrappers_pb.UInt32Value;
      reader.readMessage(value,google_protobuf_wrappers_pb.UInt32Value.deserializeBinaryFromReader);
      msg.setMaxInternalRedirects(value);
      break;
    case 27:
      var value = new proto.envoy.config.route.v3.HedgePolicy;
      reader.readMessage(value,proto.envoy.config.route.v3.HedgePolicy.deserializeBinaryFromReader);
      msg.setHedgePolicy(value);
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
proto.envoy.config.route.v3.RouteAction.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.envoy.config.route.v3.RouteAction.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.envoy.config.route.v3.RouteAction} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.RouteAction.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = /** @type {string} */ (jspb.Message.getField(message, 1));
  if (f != null) {
    writer.writeString(
      1,
      f
    );
  }
  f = /** @type {string} */ (jspb.Message.getField(message, 2));
  if (f != null) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getWeightedClusters();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      proto.envoy.config.route.v3.WeightedCluster.serializeBinaryToWriter
    );
  }
  f = message.getClusterNotFoundResponseCode();
  if (f !== 0.0) {
    writer.writeEnum(
      20,
      f
    );
  }
  f = message.getMetadataMatch();
  if (f != null) {
    writer.writeMessage(
      4,
      f,
      envoy_config_core_v3_base_pb.Metadata.serializeBinaryToWriter
    );
  }
  f = message.getPrefixRewrite();
  if (f.length > 0) {
    writer.writeString(
      5,
      f
    );
  }
  f = message.getRegexRewrite();
  if (f != null) {
    writer.writeMessage(
      32,
      f,
      envoy_type_matcher_v3_regex_pb.RegexMatchAndSubstitute.serializeBinaryToWriter
    );
  }
  f = /** @type {string} */ (jspb.Message.getField(message, 6));
  if (f != null) {
    writer.writeString(
      6,
      f
    );
  }
  f = message.getAutoHostRewrite();
  if (f != null) {
    writer.writeMessage(
      7,
      f,
      google_protobuf_wrappers_pb.BoolValue.serializeBinaryToWriter
    );
  }
  f = /** @type {string} */ (jspb.Message.getField(message, 29));
  if (f != null) {
    writer.writeString(
      29,
      f
    );
  }
  f = message.getTimeout();
  if (f != null) {
    writer.writeMessage(
      8,
      f,
      google_protobuf_duration_pb.Duration.serializeBinaryToWriter
    );
  }
  f = message.getIdleTimeout();
  if (f != null) {
    writer.writeMessage(
      24,
      f,
      google_protobuf_duration_pb.Duration.serializeBinaryToWriter
    );
  }
  f = message.getRetryPolicy();
  if (f != null) {
    writer.writeMessage(
      9,
      f,
      proto.envoy.config.route.v3.RetryPolicy.serializeBinaryToWriter
    );
  }
  f = message.getRetryPolicyTypedConfig();
  if (f != null) {
    writer.writeMessage(
      33,
      f,
      google_protobuf_any_pb.Any.serializeBinaryToWriter
    );
  }
  f = message.getRequestMirrorPoliciesList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      30,
      f,
      proto.envoy.config.route.v3.RouteAction.RequestMirrorPolicy.serializeBinaryToWriter
    );
  }
  f = message.getPriority();
  if (f !== 0.0) {
    writer.writeEnum(
      11,
      f
    );
  }
  f = message.getRateLimitsList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      13,
      f,
      proto.envoy.config.route.v3.RateLimit.serializeBinaryToWriter
    );
  }
  f = message.getIncludeVhRateLimits();
  if (f != null) {
    writer.writeMessage(
      14,
      f,
      google_protobuf_wrappers_pb.BoolValue.serializeBinaryToWriter
    );
  }
  f = message.getHashPolicyList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      15,
      f,
      proto.envoy.config.route.v3.RouteAction.HashPolicy.serializeBinaryToWriter
    );
  }
  f = message.getCors();
  if (f != null) {
    writer.writeMessage(
      17,
      f,
      proto.envoy.config.route.v3.CorsPolicy.serializeBinaryToWriter
    );
  }
  f = message.getMaxGrpcTimeout();
  if (f != null) {
    writer.writeMessage(
      23,
      f,
      google_protobuf_duration_pb.Duration.serializeBinaryToWriter
    );
  }
  f = message.getGrpcTimeoutOffset();
  if (f != null) {
    writer.writeMessage(
      28,
      f,
      google_protobuf_duration_pb.Duration.serializeBinaryToWriter
    );
  }
  f = message.getUpgradeConfigsList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      25,
      f,
      proto.envoy.config.route.v3.RouteAction.UpgradeConfig.serializeBinaryToWriter
    );
  }
  f = message.getInternalRedirectPolicy();
  if (f != null) {
    writer.writeMessage(
      34,
      f,
      proto.envoy.config.route.v3.InternalRedirectPolicy.serializeBinaryToWriter
    );
  }
  f = message.getInternalRedirectAction();
  if (f !== 0.0) {
    writer.writeEnum(
      26,
      f
    );
  }
  f = message.getMaxInternalRedirects();
  if (f != null) {
    writer.writeMessage(
      31,
      f,
      google_protobuf_wrappers_pb.UInt32Value.serializeBinaryToWriter
    );
  }
  f = message.getHedgePolicy();
  if (f != null) {
    writer.writeMessage(
      27,
      f,
      proto.envoy.config.route.v3.HedgePolicy.serializeBinaryToWriter
    );
  }
};


/**
 * @enum {number}
 */
proto.envoy.config.route.v3.RouteAction.ClusterNotFoundResponseCode = {
  SERVICE_UNAVAILABLE: 0,
  NOT_FOUND: 1
};

/**
 * @enum {number}
 */
proto.envoy.config.route.v3.RouteAction.InternalRedirectAction = {
  PASS_THROUGH_INTERNAL_REDIRECT: 0,
  HANDLE_INTERNAL_REDIRECT: 1
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
proto.envoy.config.route.v3.RouteAction.RequestMirrorPolicy = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.envoy.config.route.v3.RouteAction.RequestMirrorPolicy, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.envoy.config.route.v3.RouteAction.RequestMirrorPolicy.displayName = 'proto.envoy.config.route.v3.RouteAction.RequestMirrorPolicy';
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
proto.envoy.config.route.v3.RouteAction.RequestMirrorPolicy.prototype.toObject = function(opt_includeInstance) {
  return proto.envoy.config.route.v3.RouteAction.RequestMirrorPolicy.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.envoy.config.route.v3.RouteAction.RequestMirrorPolicy} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.RouteAction.RequestMirrorPolicy.toObject = function(includeInstance, msg) {
  var f, obj = {
    cluster: jspb.Message.getFieldWithDefault(msg, 1, ""),
    runtimeFraction: (f = msg.getRuntimeFraction()) && envoy_config_core_v3_base_pb.RuntimeFractionalPercent.toObject(includeInstance, f),
    traceSampled: (f = msg.getTraceSampled()) && google_protobuf_wrappers_pb.BoolValue.toObject(includeInstance, f)
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
 * @return {!proto.envoy.config.route.v3.RouteAction.RequestMirrorPolicy}
 */
proto.envoy.config.route.v3.RouteAction.RequestMirrorPolicy.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.envoy.config.route.v3.RouteAction.RequestMirrorPolicy;
  return proto.envoy.config.route.v3.RouteAction.RequestMirrorPolicy.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.envoy.config.route.v3.RouteAction.RequestMirrorPolicy} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.envoy.config.route.v3.RouteAction.RequestMirrorPolicy}
 */
proto.envoy.config.route.v3.RouteAction.RequestMirrorPolicy.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setCluster(value);
      break;
    case 3:
      var value = new envoy_config_core_v3_base_pb.RuntimeFractionalPercent;
      reader.readMessage(value,envoy_config_core_v3_base_pb.RuntimeFractionalPercent.deserializeBinaryFromReader);
      msg.setRuntimeFraction(value);
      break;
    case 4:
      var value = new google_protobuf_wrappers_pb.BoolValue;
      reader.readMessage(value,google_protobuf_wrappers_pb.BoolValue.deserializeBinaryFromReader);
      msg.setTraceSampled(value);
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
proto.envoy.config.route.v3.RouteAction.RequestMirrorPolicy.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.envoy.config.route.v3.RouteAction.RequestMirrorPolicy.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.envoy.config.route.v3.RouteAction.RequestMirrorPolicy} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.RouteAction.RequestMirrorPolicy.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getCluster();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getRuntimeFraction();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      envoy_config_core_v3_base_pb.RuntimeFractionalPercent.serializeBinaryToWriter
    );
  }
  f = message.getTraceSampled();
  if (f != null) {
    writer.writeMessage(
      4,
      f,
      google_protobuf_wrappers_pb.BoolValue.serializeBinaryToWriter
    );
  }
};


/**
 * optional string cluster = 1;
 * @return {string}
 */
proto.envoy.config.route.v3.RouteAction.RequestMirrorPolicy.prototype.getCluster = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.envoy.config.route.v3.RouteAction.RequestMirrorPolicy.prototype.setCluster = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional envoy.config.core.v3.RuntimeFractionalPercent runtime_fraction = 3;
 * @return {?proto.envoy.config.core.v3.RuntimeFractionalPercent}
 */
proto.envoy.config.route.v3.RouteAction.RequestMirrorPolicy.prototype.getRuntimeFraction = function() {
  return /** @type{?proto.envoy.config.core.v3.RuntimeFractionalPercent} */ (
    jspb.Message.getWrapperField(this, envoy_config_core_v3_base_pb.RuntimeFractionalPercent, 3));
};


/** @param {?proto.envoy.config.core.v3.RuntimeFractionalPercent|undefined} value */
proto.envoy.config.route.v3.RouteAction.RequestMirrorPolicy.prototype.setRuntimeFraction = function(value) {
  jspb.Message.setWrapperField(this, 3, value);
};


proto.envoy.config.route.v3.RouteAction.RequestMirrorPolicy.prototype.clearRuntimeFraction = function() {
  this.setRuntimeFraction(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.RouteAction.RequestMirrorPolicy.prototype.hasRuntimeFraction = function() {
  return jspb.Message.getField(this, 3) != null;
};


/**
 * optional google.protobuf.BoolValue trace_sampled = 4;
 * @return {?proto.google.protobuf.BoolValue}
 */
proto.envoy.config.route.v3.RouteAction.RequestMirrorPolicy.prototype.getTraceSampled = function() {
  return /** @type{?proto.google.protobuf.BoolValue} */ (
    jspb.Message.getWrapperField(this, google_protobuf_wrappers_pb.BoolValue, 4));
};


/** @param {?proto.google.protobuf.BoolValue|undefined} value */
proto.envoy.config.route.v3.RouteAction.RequestMirrorPolicy.prototype.setTraceSampled = function(value) {
  jspb.Message.setWrapperField(this, 4, value);
};


proto.envoy.config.route.v3.RouteAction.RequestMirrorPolicy.prototype.clearTraceSampled = function() {
  this.setTraceSampled(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.RouteAction.RequestMirrorPolicy.prototype.hasTraceSampled = function() {
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
proto.envoy.config.route.v3.RouteAction.HashPolicy = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, proto.envoy.config.route.v3.RouteAction.HashPolicy.oneofGroups_);
};
goog.inherits(proto.envoy.config.route.v3.RouteAction.HashPolicy, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.envoy.config.route.v3.RouteAction.HashPolicy.displayName = 'proto.envoy.config.route.v3.RouteAction.HashPolicy';
}
/**
 * Oneof group definitions for this message. Each group defines the field
 * numbers belonging to that group. When of these fields' value is set, all
 * other fields in the group are cleared. During deserialization, if multiple
 * fields are encountered for a group, only the last value seen will be kept.
 * @private {!Array<!Array<number>>}
 * @const
 */
proto.envoy.config.route.v3.RouteAction.HashPolicy.oneofGroups_ = [[1,2,3,5,6]];

/**
 * @enum {number}
 */
proto.envoy.config.route.v3.RouteAction.HashPolicy.PolicySpecifierCase = {
  POLICY_SPECIFIER_NOT_SET: 0,
  HEADER: 1,
  COOKIE: 2,
  CONNECTION_PROPERTIES: 3,
  QUERY_PARAMETER: 5,
  FILTER_STATE: 6
};

/**
 * @return {proto.envoy.config.route.v3.RouteAction.HashPolicy.PolicySpecifierCase}
 */
proto.envoy.config.route.v3.RouteAction.HashPolicy.prototype.getPolicySpecifierCase = function() {
  return /** @type {proto.envoy.config.route.v3.RouteAction.HashPolicy.PolicySpecifierCase} */(jspb.Message.computeOneofCase(this, proto.envoy.config.route.v3.RouteAction.HashPolicy.oneofGroups_[0]));
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
proto.envoy.config.route.v3.RouteAction.HashPolicy.prototype.toObject = function(opt_includeInstance) {
  return proto.envoy.config.route.v3.RouteAction.HashPolicy.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.envoy.config.route.v3.RouteAction.HashPolicy} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.RouteAction.HashPolicy.toObject = function(includeInstance, msg) {
  var f, obj = {
    header: (f = msg.getHeader()) && proto.envoy.config.route.v3.RouteAction.HashPolicy.Header.toObject(includeInstance, f),
    cookie: (f = msg.getCookie()) && proto.envoy.config.route.v3.RouteAction.HashPolicy.Cookie.toObject(includeInstance, f),
    connectionProperties: (f = msg.getConnectionProperties()) && proto.envoy.config.route.v3.RouteAction.HashPolicy.ConnectionProperties.toObject(includeInstance, f),
    queryParameter: (f = msg.getQueryParameter()) && proto.envoy.config.route.v3.RouteAction.HashPolicy.QueryParameter.toObject(includeInstance, f),
    filterState: (f = msg.getFilterState()) && proto.envoy.config.route.v3.RouteAction.HashPolicy.FilterState.toObject(includeInstance, f),
    terminal: jspb.Message.getFieldWithDefault(msg, 4, false)
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
 * @return {!proto.envoy.config.route.v3.RouteAction.HashPolicy}
 */
proto.envoy.config.route.v3.RouteAction.HashPolicy.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.envoy.config.route.v3.RouteAction.HashPolicy;
  return proto.envoy.config.route.v3.RouteAction.HashPolicy.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.envoy.config.route.v3.RouteAction.HashPolicy} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.envoy.config.route.v3.RouteAction.HashPolicy}
 */
proto.envoy.config.route.v3.RouteAction.HashPolicy.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.envoy.config.route.v3.RouteAction.HashPolicy.Header;
      reader.readMessage(value,proto.envoy.config.route.v3.RouteAction.HashPolicy.Header.deserializeBinaryFromReader);
      msg.setHeader(value);
      break;
    case 2:
      var value = new proto.envoy.config.route.v3.RouteAction.HashPolicy.Cookie;
      reader.readMessage(value,proto.envoy.config.route.v3.RouteAction.HashPolicy.Cookie.deserializeBinaryFromReader);
      msg.setCookie(value);
      break;
    case 3:
      var value = new proto.envoy.config.route.v3.RouteAction.HashPolicy.ConnectionProperties;
      reader.readMessage(value,proto.envoy.config.route.v3.RouteAction.HashPolicy.ConnectionProperties.deserializeBinaryFromReader);
      msg.setConnectionProperties(value);
      break;
    case 5:
      var value = new proto.envoy.config.route.v3.RouteAction.HashPolicy.QueryParameter;
      reader.readMessage(value,proto.envoy.config.route.v3.RouteAction.HashPolicy.QueryParameter.deserializeBinaryFromReader);
      msg.setQueryParameter(value);
      break;
    case 6:
      var value = new proto.envoy.config.route.v3.RouteAction.HashPolicy.FilterState;
      reader.readMessage(value,proto.envoy.config.route.v3.RouteAction.HashPolicy.FilterState.deserializeBinaryFromReader);
      msg.setFilterState(value);
      break;
    case 4:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setTerminal(value);
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
proto.envoy.config.route.v3.RouteAction.HashPolicy.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.envoy.config.route.v3.RouteAction.HashPolicy.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.envoy.config.route.v3.RouteAction.HashPolicy} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.RouteAction.HashPolicy.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getHeader();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto.envoy.config.route.v3.RouteAction.HashPolicy.Header.serializeBinaryToWriter
    );
  }
  f = message.getCookie();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.envoy.config.route.v3.RouteAction.HashPolicy.Cookie.serializeBinaryToWriter
    );
  }
  f = message.getConnectionProperties();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      proto.envoy.config.route.v3.RouteAction.HashPolicy.ConnectionProperties.serializeBinaryToWriter
    );
  }
  f = message.getQueryParameter();
  if (f != null) {
    writer.writeMessage(
      5,
      f,
      proto.envoy.config.route.v3.RouteAction.HashPolicy.QueryParameter.serializeBinaryToWriter
    );
  }
  f = message.getFilterState();
  if (f != null) {
    writer.writeMessage(
      6,
      f,
      proto.envoy.config.route.v3.RouteAction.HashPolicy.FilterState.serializeBinaryToWriter
    );
  }
  f = message.getTerminal();
  if (f) {
    writer.writeBool(
      4,
      f
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
proto.envoy.config.route.v3.RouteAction.HashPolicy.Header = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.envoy.config.route.v3.RouteAction.HashPolicy.Header, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.envoy.config.route.v3.RouteAction.HashPolicy.Header.displayName = 'proto.envoy.config.route.v3.RouteAction.HashPolicy.Header';
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
proto.envoy.config.route.v3.RouteAction.HashPolicy.Header.prototype.toObject = function(opt_includeInstance) {
  return proto.envoy.config.route.v3.RouteAction.HashPolicy.Header.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.envoy.config.route.v3.RouteAction.HashPolicy.Header} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.RouteAction.HashPolicy.Header.toObject = function(includeInstance, msg) {
  var f, obj = {
    headerName: jspb.Message.getFieldWithDefault(msg, 1, ""),
    regexRewrite: (f = msg.getRegexRewrite()) && envoy_type_matcher_v3_regex_pb.RegexMatchAndSubstitute.toObject(includeInstance, f)
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
 * @return {!proto.envoy.config.route.v3.RouteAction.HashPolicy.Header}
 */
proto.envoy.config.route.v3.RouteAction.HashPolicy.Header.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.envoy.config.route.v3.RouteAction.HashPolicy.Header;
  return proto.envoy.config.route.v3.RouteAction.HashPolicy.Header.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.envoy.config.route.v3.RouteAction.HashPolicy.Header} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.envoy.config.route.v3.RouteAction.HashPolicy.Header}
 */
proto.envoy.config.route.v3.RouteAction.HashPolicy.Header.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setHeaderName(value);
      break;
    case 2:
      var value = new envoy_type_matcher_v3_regex_pb.RegexMatchAndSubstitute;
      reader.readMessage(value,envoy_type_matcher_v3_regex_pb.RegexMatchAndSubstitute.deserializeBinaryFromReader);
      msg.setRegexRewrite(value);
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
proto.envoy.config.route.v3.RouteAction.HashPolicy.Header.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.envoy.config.route.v3.RouteAction.HashPolicy.Header.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.envoy.config.route.v3.RouteAction.HashPolicy.Header} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.RouteAction.HashPolicy.Header.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getHeaderName();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getRegexRewrite();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      envoy_type_matcher_v3_regex_pb.RegexMatchAndSubstitute.serializeBinaryToWriter
    );
  }
};


/**
 * optional string header_name = 1;
 * @return {string}
 */
proto.envoy.config.route.v3.RouteAction.HashPolicy.Header.prototype.getHeaderName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.envoy.config.route.v3.RouteAction.HashPolicy.Header.prototype.setHeaderName = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional envoy.type.matcher.v3.RegexMatchAndSubstitute regex_rewrite = 2;
 * @return {?proto.envoy.type.matcher.v3.RegexMatchAndSubstitute}
 */
proto.envoy.config.route.v3.RouteAction.HashPolicy.Header.prototype.getRegexRewrite = function() {
  return /** @type{?proto.envoy.type.matcher.v3.RegexMatchAndSubstitute} */ (
    jspb.Message.getWrapperField(this, envoy_type_matcher_v3_regex_pb.RegexMatchAndSubstitute, 2));
};


/** @param {?proto.envoy.type.matcher.v3.RegexMatchAndSubstitute|undefined} value */
proto.envoy.config.route.v3.RouteAction.HashPolicy.Header.prototype.setRegexRewrite = function(value) {
  jspb.Message.setWrapperField(this, 2, value);
};


proto.envoy.config.route.v3.RouteAction.HashPolicy.Header.prototype.clearRegexRewrite = function() {
  this.setRegexRewrite(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.RouteAction.HashPolicy.Header.prototype.hasRegexRewrite = function() {
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
proto.envoy.config.route.v3.RouteAction.HashPolicy.Cookie = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.envoy.config.route.v3.RouteAction.HashPolicy.Cookie, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.envoy.config.route.v3.RouteAction.HashPolicy.Cookie.displayName = 'proto.envoy.config.route.v3.RouteAction.HashPolicy.Cookie';
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
proto.envoy.config.route.v3.RouteAction.HashPolicy.Cookie.prototype.toObject = function(opt_includeInstance) {
  return proto.envoy.config.route.v3.RouteAction.HashPolicy.Cookie.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.envoy.config.route.v3.RouteAction.HashPolicy.Cookie} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.RouteAction.HashPolicy.Cookie.toObject = function(includeInstance, msg) {
  var f, obj = {
    name: jspb.Message.getFieldWithDefault(msg, 1, ""),
    ttl: (f = msg.getTtl()) && google_protobuf_duration_pb.Duration.toObject(includeInstance, f),
    path: jspb.Message.getFieldWithDefault(msg, 3, "")
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
 * @return {!proto.envoy.config.route.v3.RouteAction.HashPolicy.Cookie}
 */
proto.envoy.config.route.v3.RouteAction.HashPolicy.Cookie.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.envoy.config.route.v3.RouteAction.HashPolicy.Cookie;
  return proto.envoy.config.route.v3.RouteAction.HashPolicy.Cookie.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.envoy.config.route.v3.RouteAction.HashPolicy.Cookie} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.envoy.config.route.v3.RouteAction.HashPolicy.Cookie}
 */
proto.envoy.config.route.v3.RouteAction.HashPolicy.Cookie.deserializeBinaryFromReader = function(msg, reader) {
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
      var value = new google_protobuf_duration_pb.Duration;
      reader.readMessage(value,google_protobuf_duration_pb.Duration.deserializeBinaryFromReader);
      msg.setTtl(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setPath(value);
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
proto.envoy.config.route.v3.RouteAction.HashPolicy.Cookie.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.envoy.config.route.v3.RouteAction.HashPolicy.Cookie.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.envoy.config.route.v3.RouteAction.HashPolicy.Cookie} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.RouteAction.HashPolicy.Cookie.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getName();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getTtl();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      google_protobuf_duration_pb.Duration.serializeBinaryToWriter
    );
  }
  f = message.getPath();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
};


/**
 * optional string name = 1;
 * @return {string}
 */
proto.envoy.config.route.v3.RouteAction.HashPolicy.Cookie.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.envoy.config.route.v3.RouteAction.HashPolicy.Cookie.prototype.setName = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional google.protobuf.Duration ttl = 2;
 * @return {?proto.google.protobuf.Duration}
 */
proto.envoy.config.route.v3.RouteAction.HashPolicy.Cookie.prototype.getTtl = function() {
  return /** @type{?proto.google.protobuf.Duration} */ (
    jspb.Message.getWrapperField(this, google_protobuf_duration_pb.Duration, 2));
};


/** @param {?proto.google.protobuf.Duration|undefined} value */
proto.envoy.config.route.v3.RouteAction.HashPolicy.Cookie.prototype.setTtl = function(value) {
  jspb.Message.setWrapperField(this, 2, value);
};


proto.envoy.config.route.v3.RouteAction.HashPolicy.Cookie.prototype.clearTtl = function() {
  this.setTtl(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.RouteAction.HashPolicy.Cookie.prototype.hasTtl = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional string path = 3;
 * @return {string}
 */
proto.envoy.config.route.v3.RouteAction.HashPolicy.Cookie.prototype.getPath = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/** @param {string} value */
proto.envoy.config.route.v3.RouteAction.HashPolicy.Cookie.prototype.setPath = function(value) {
  jspb.Message.setProto3StringField(this, 3, value);
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
proto.envoy.config.route.v3.RouteAction.HashPolicy.ConnectionProperties = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.envoy.config.route.v3.RouteAction.HashPolicy.ConnectionProperties, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.envoy.config.route.v3.RouteAction.HashPolicy.ConnectionProperties.displayName = 'proto.envoy.config.route.v3.RouteAction.HashPolicy.ConnectionProperties';
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
proto.envoy.config.route.v3.RouteAction.HashPolicy.ConnectionProperties.prototype.toObject = function(opt_includeInstance) {
  return proto.envoy.config.route.v3.RouteAction.HashPolicy.ConnectionProperties.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.envoy.config.route.v3.RouteAction.HashPolicy.ConnectionProperties} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.RouteAction.HashPolicy.ConnectionProperties.toObject = function(includeInstance, msg) {
  var f, obj = {
    sourceIp: jspb.Message.getFieldWithDefault(msg, 1, false)
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
 * @return {!proto.envoy.config.route.v3.RouteAction.HashPolicy.ConnectionProperties}
 */
proto.envoy.config.route.v3.RouteAction.HashPolicy.ConnectionProperties.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.envoy.config.route.v3.RouteAction.HashPolicy.ConnectionProperties;
  return proto.envoy.config.route.v3.RouteAction.HashPolicy.ConnectionProperties.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.envoy.config.route.v3.RouteAction.HashPolicy.ConnectionProperties} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.envoy.config.route.v3.RouteAction.HashPolicy.ConnectionProperties}
 */
proto.envoy.config.route.v3.RouteAction.HashPolicy.ConnectionProperties.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setSourceIp(value);
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
proto.envoy.config.route.v3.RouteAction.HashPolicy.ConnectionProperties.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.envoy.config.route.v3.RouteAction.HashPolicy.ConnectionProperties.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.envoy.config.route.v3.RouteAction.HashPolicy.ConnectionProperties} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.RouteAction.HashPolicy.ConnectionProperties.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getSourceIp();
  if (f) {
    writer.writeBool(
      1,
      f
    );
  }
};


/**
 * optional bool source_ip = 1;
 * Note that Boolean fields may be set to 0/1 when serialized from a Java server.
 * You should avoid comparisons like {@code val === true/false} in those cases.
 * @return {boolean}
 */
proto.envoy.config.route.v3.RouteAction.HashPolicy.ConnectionProperties.prototype.getSourceIp = function() {
  return /** @type {boolean} */ (jspb.Message.getFieldWithDefault(this, 1, false));
};


/** @param {boolean} value */
proto.envoy.config.route.v3.RouteAction.HashPolicy.ConnectionProperties.prototype.setSourceIp = function(value) {
  jspb.Message.setProto3BooleanField(this, 1, value);
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
proto.envoy.config.route.v3.RouteAction.HashPolicy.QueryParameter = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.envoy.config.route.v3.RouteAction.HashPolicy.QueryParameter, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.envoy.config.route.v3.RouteAction.HashPolicy.QueryParameter.displayName = 'proto.envoy.config.route.v3.RouteAction.HashPolicy.QueryParameter';
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
proto.envoy.config.route.v3.RouteAction.HashPolicy.QueryParameter.prototype.toObject = function(opt_includeInstance) {
  return proto.envoy.config.route.v3.RouteAction.HashPolicy.QueryParameter.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.envoy.config.route.v3.RouteAction.HashPolicy.QueryParameter} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.RouteAction.HashPolicy.QueryParameter.toObject = function(includeInstance, msg) {
  var f, obj = {
    name: jspb.Message.getFieldWithDefault(msg, 1, "")
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
 * @return {!proto.envoy.config.route.v3.RouteAction.HashPolicy.QueryParameter}
 */
proto.envoy.config.route.v3.RouteAction.HashPolicy.QueryParameter.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.envoy.config.route.v3.RouteAction.HashPolicy.QueryParameter;
  return proto.envoy.config.route.v3.RouteAction.HashPolicy.QueryParameter.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.envoy.config.route.v3.RouteAction.HashPolicy.QueryParameter} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.envoy.config.route.v3.RouteAction.HashPolicy.QueryParameter}
 */
proto.envoy.config.route.v3.RouteAction.HashPolicy.QueryParameter.deserializeBinaryFromReader = function(msg, reader) {
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
proto.envoy.config.route.v3.RouteAction.HashPolicy.QueryParameter.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.envoy.config.route.v3.RouteAction.HashPolicy.QueryParameter.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.envoy.config.route.v3.RouteAction.HashPolicy.QueryParameter} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.RouteAction.HashPolicy.QueryParameter.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getName();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string name = 1;
 * @return {string}
 */
proto.envoy.config.route.v3.RouteAction.HashPolicy.QueryParameter.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.envoy.config.route.v3.RouteAction.HashPolicy.QueryParameter.prototype.setName = function(value) {
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
proto.envoy.config.route.v3.RouteAction.HashPolicy.FilterState = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.envoy.config.route.v3.RouteAction.HashPolicy.FilterState, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.envoy.config.route.v3.RouteAction.HashPolicy.FilterState.displayName = 'proto.envoy.config.route.v3.RouteAction.HashPolicy.FilterState';
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
proto.envoy.config.route.v3.RouteAction.HashPolicy.FilterState.prototype.toObject = function(opt_includeInstance) {
  return proto.envoy.config.route.v3.RouteAction.HashPolicy.FilterState.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.envoy.config.route.v3.RouteAction.HashPolicy.FilterState} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.RouteAction.HashPolicy.FilterState.toObject = function(includeInstance, msg) {
  var f, obj = {
    key: jspb.Message.getFieldWithDefault(msg, 1, "")
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
 * @return {!proto.envoy.config.route.v3.RouteAction.HashPolicy.FilterState}
 */
proto.envoy.config.route.v3.RouteAction.HashPolicy.FilterState.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.envoy.config.route.v3.RouteAction.HashPolicy.FilterState;
  return proto.envoy.config.route.v3.RouteAction.HashPolicy.FilterState.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.envoy.config.route.v3.RouteAction.HashPolicy.FilterState} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.envoy.config.route.v3.RouteAction.HashPolicy.FilterState}
 */
proto.envoy.config.route.v3.RouteAction.HashPolicy.FilterState.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setKey(value);
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
proto.envoy.config.route.v3.RouteAction.HashPolicy.FilterState.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.envoy.config.route.v3.RouteAction.HashPolicy.FilterState.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.envoy.config.route.v3.RouteAction.HashPolicy.FilterState} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.RouteAction.HashPolicy.FilterState.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getKey();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string key = 1;
 * @return {string}
 */
proto.envoy.config.route.v3.RouteAction.HashPolicy.FilterState.prototype.getKey = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.envoy.config.route.v3.RouteAction.HashPolicy.FilterState.prototype.setKey = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional Header header = 1;
 * @return {?proto.envoy.config.route.v3.RouteAction.HashPolicy.Header}
 */
proto.envoy.config.route.v3.RouteAction.HashPolicy.prototype.getHeader = function() {
  return /** @type{?proto.envoy.config.route.v3.RouteAction.HashPolicy.Header} */ (
    jspb.Message.getWrapperField(this, proto.envoy.config.route.v3.RouteAction.HashPolicy.Header, 1));
};


/** @param {?proto.envoy.config.route.v3.RouteAction.HashPolicy.Header|undefined} value */
proto.envoy.config.route.v3.RouteAction.HashPolicy.prototype.setHeader = function(value) {
  jspb.Message.setOneofWrapperField(this, 1, proto.envoy.config.route.v3.RouteAction.HashPolicy.oneofGroups_[0], value);
};


proto.envoy.config.route.v3.RouteAction.HashPolicy.prototype.clearHeader = function() {
  this.setHeader(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.RouteAction.HashPolicy.prototype.hasHeader = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional Cookie cookie = 2;
 * @return {?proto.envoy.config.route.v3.RouteAction.HashPolicy.Cookie}
 */
proto.envoy.config.route.v3.RouteAction.HashPolicy.prototype.getCookie = function() {
  return /** @type{?proto.envoy.config.route.v3.RouteAction.HashPolicy.Cookie} */ (
    jspb.Message.getWrapperField(this, proto.envoy.config.route.v3.RouteAction.HashPolicy.Cookie, 2));
};


/** @param {?proto.envoy.config.route.v3.RouteAction.HashPolicy.Cookie|undefined} value */
proto.envoy.config.route.v3.RouteAction.HashPolicy.prototype.setCookie = function(value) {
  jspb.Message.setOneofWrapperField(this, 2, proto.envoy.config.route.v3.RouteAction.HashPolicy.oneofGroups_[0], value);
};


proto.envoy.config.route.v3.RouteAction.HashPolicy.prototype.clearCookie = function() {
  this.setCookie(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.RouteAction.HashPolicy.prototype.hasCookie = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional ConnectionProperties connection_properties = 3;
 * @return {?proto.envoy.config.route.v3.RouteAction.HashPolicy.ConnectionProperties}
 */
proto.envoy.config.route.v3.RouteAction.HashPolicy.prototype.getConnectionProperties = function() {
  return /** @type{?proto.envoy.config.route.v3.RouteAction.HashPolicy.ConnectionProperties} */ (
    jspb.Message.getWrapperField(this, proto.envoy.config.route.v3.RouteAction.HashPolicy.ConnectionProperties, 3));
};


/** @param {?proto.envoy.config.route.v3.RouteAction.HashPolicy.ConnectionProperties|undefined} value */
proto.envoy.config.route.v3.RouteAction.HashPolicy.prototype.setConnectionProperties = function(value) {
  jspb.Message.setOneofWrapperField(this, 3, proto.envoy.config.route.v3.RouteAction.HashPolicy.oneofGroups_[0], value);
};


proto.envoy.config.route.v3.RouteAction.HashPolicy.prototype.clearConnectionProperties = function() {
  this.setConnectionProperties(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.RouteAction.HashPolicy.prototype.hasConnectionProperties = function() {
  return jspb.Message.getField(this, 3) != null;
};


/**
 * optional QueryParameter query_parameter = 5;
 * @return {?proto.envoy.config.route.v3.RouteAction.HashPolicy.QueryParameter}
 */
proto.envoy.config.route.v3.RouteAction.HashPolicy.prototype.getQueryParameter = function() {
  return /** @type{?proto.envoy.config.route.v3.RouteAction.HashPolicy.QueryParameter} */ (
    jspb.Message.getWrapperField(this, proto.envoy.config.route.v3.RouteAction.HashPolicy.QueryParameter, 5));
};


/** @param {?proto.envoy.config.route.v3.RouteAction.HashPolicy.QueryParameter|undefined} value */
proto.envoy.config.route.v3.RouteAction.HashPolicy.prototype.setQueryParameter = function(value) {
  jspb.Message.setOneofWrapperField(this, 5, proto.envoy.config.route.v3.RouteAction.HashPolicy.oneofGroups_[0], value);
};


proto.envoy.config.route.v3.RouteAction.HashPolicy.prototype.clearQueryParameter = function() {
  this.setQueryParameter(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.RouteAction.HashPolicy.prototype.hasQueryParameter = function() {
  return jspb.Message.getField(this, 5) != null;
};


/**
 * optional FilterState filter_state = 6;
 * @return {?proto.envoy.config.route.v3.RouteAction.HashPolicy.FilterState}
 */
proto.envoy.config.route.v3.RouteAction.HashPolicy.prototype.getFilterState = function() {
  return /** @type{?proto.envoy.config.route.v3.RouteAction.HashPolicy.FilterState} */ (
    jspb.Message.getWrapperField(this, proto.envoy.config.route.v3.RouteAction.HashPolicy.FilterState, 6));
};


/** @param {?proto.envoy.config.route.v3.RouteAction.HashPolicy.FilterState|undefined} value */
proto.envoy.config.route.v3.RouteAction.HashPolicy.prototype.setFilterState = function(value) {
  jspb.Message.setOneofWrapperField(this, 6, proto.envoy.config.route.v3.RouteAction.HashPolicy.oneofGroups_[0], value);
};


proto.envoy.config.route.v3.RouteAction.HashPolicy.prototype.clearFilterState = function() {
  this.setFilterState(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.RouteAction.HashPolicy.prototype.hasFilterState = function() {
  return jspb.Message.getField(this, 6) != null;
};


/**
 * optional bool terminal = 4;
 * Note that Boolean fields may be set to 0/1 when serialized from a Java server.
 * You should avoid comparisons like {@code val === true/false} in those cases.
 * @return {boolean}
 */
proto.envoy.config.route.v3.RouteAction.HashPolicy.prototype.getTerminal = function() {
  return /** @type {boolean} */ (jspb.Message.getFieldWithDefault(this, 4, false));
};


/** @param {boolean} value */
proto.envoy.config.route.v3.RouteAction.HashPolicy.prototype.setTerminal = function(value) {
  jspb.Message.setProto3BooleanField(this, 4, value);
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
proto.envoy.config.route.v3.RouteAction.UpgradeConfig = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.envoy.config.route.v3.RouteAction.UpgradeConfig, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.envoy.config.route.v3.RouteAction.UpgradeConfig.displayName = 'proto.envoy.config.route.v3.RouteAction.UpgradeConfig';
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
proto.envoy.config.route.v3.RouteAction.UpgradeConfig.prototype.toObject = function(opt_includeInstance) {
  return proto.envoy.config.route.v3.RouteAction.UpgradeConfig.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.envoy.config.route.v3.RouteAction.UpgradeConfig} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.RouteAction.UpgradeConfig.toObject = function(includeInstance, msg) {
  var f, obj = {
    upgradeType: jspb.Message.getFieldWithDefault(msg, 1, ""),
    enabled: (f = msg.getEnabled()) && google_protobuf_wrappers_pb.BoolValue.toObject(includeInstance, f),
    connectConfig: (f = msg.getConnectConfig()) && proto.envoy.config.route.v3.RouteAction.UpgradeConfig.ConnectConfig.toObject(includeInstance, f)
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
 * @return {!proto.envoy.config.route.v3.RouteAction.UpgradeConfig}
 */
proto.envoy.config.route.v3.RouteAction.UpgradeConfig.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.envoy.config.route.v3.RouteAction.UpgradeConfig;
  return proto.envoy.config.route.v3.RouteAction.UpgradeConfig.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.envoy.config.route.v3.RouteAction.UpgradeConfig} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.envoy.config.route.v3.RouteAction.UpgradeConfig}
 */
proto.envoy.config.route.v3.RouteAction.UpgradeConfig.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setUpgradeType(value);
      break;
    case 2:
      var value = new google_protobuf_wrappers_pb.BoolValue;
      reader.readMessage(value,google_protobuf_wrappers_pb.BoolValue.deserializeBinaryFromReader);
      msg.setEnabled(value);
      break;
    case 3:
      var value = new proto.envoy.config.route.v3.RouteAction.UpgradeConfig.ConnectConfig;
      reader.readMessage(value,proto.envoy.config.route.v3.RouteAction.UpgradeConfig.ConnectConfig.deserializeBinaryFromReader);
      msg.setConnectConfig(value);
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
proto.envoy.config.route.v3.RouteAction.UpgradeConfig.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.envoy.config.route.v3.RouteAction.UpgradeConfig.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.envoy.config.route.v3.RouteAction.UpgradeConfig} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.RouteAction.UpgradeConfig.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getUpgradeType();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getEnabled();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      google_protobuf_wrappers_pb.BoolValue.serializeBinaryToWriter
    );
  }
  f = message.getConnectConfig();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      proto.envoy.config.route.v3.RouteAction.UpgradeConfig.ConnectConfig.serializeBinaryToWriter
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
proto.envoy.config.route.v3.RouteAction.UpgradeConfig.ConnectConfig = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.envoy.config.route.v3.RouteAction.UpgradeConfig.ConnectConfig, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.envoy.config.route.v3.RouteAction.UpgradeConfig.ConnectConfig.displayName = 'proto.envoy.config.route.v3.RouteAction.UpgradeConfig.ConnectConfig';
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
proto.envoy.config.route.v3.RouteAction.UpgradeConfig.ConnectConfig.prototype.toObject = function(opt_includeInstance) {
  return proto.envoy.config.route.v3.RouteAction.UpgradeConfig.ConnectConfig.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.envoy.config.route.v3.RouteAction.UpgradeConfig.ConnectConfig} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.RouteAction.UpgradeConfig.ConnectConfig.toObject = function(includeInstance, msg) {
  var f, obj = {
    proxyProtocolConfig: (f = msg.getProxyProtocolConfig()) && envoy_config_core_v3_proxy_protocol_pb.ProxyProtocolConfig.toObject(includeInstance, f)
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
 * @return {!proto.envoy.config.route.v3.RouteAction.UpgradeConfig.ConnectConfig}
 */
proto.envoy.config.route.v3.RouteAction.UpgradeConfig.ConnectConfig.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.envoy.config.route.v3.RouteAction.UpgradeConfig.ConnectConfig;
  return proto.envoy.config.route.v3.RouteAction.UpgradeConfig.ConnectConfig.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.envoy.config.route.v3.RouteAction.UpgradeConfig.ConnectConfig} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.envoy.config.route.v3.RouteAction.UpgradeConfig.ConnectConfig}
 */
proto.envoy.config.route.v3.RouteAction.UpgradeConfig.ConnectConfig.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new envoy_config_core_v3_proxy_protocol_pb.ProxyProtocolConfig;
      reader.readMessage(value,envoy_config_core_v3_proxy_protocol_pb.ProxyProtocolConfig.deserializeBinaryFromReader);
      msg.setProxyProtocolConfig(value);
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
proto.envoy.config.route.v3.RouteAction.UpgradeConfig.ConnectConfig.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.envoy.config.route.v3.RouteAction.UpgradeConfig.ConnectConfig.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.envoy.config.route.v3.RouteAction.UpgradeConfig.ConnectConfig} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.RouteAction.UpgradeConfig.ConnectConfig.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getProxyProtocolConfig();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      envoy_config_core_v3_proxy_protocol_pb.ProxyProtocolConfig.serializeBinaryToWriter
    );
  }
};


/**
 * optional envoy.config.core.v3.ProxyProtocolConfig proxy_protocol_config = 1;
 * @return {?proto.envoy.config.core.v3.ProxyProtocolConfig}
 */
proto.envoy.config.route.v3.RouteAction.UpgradeConfig.ConnectConfig.prototype.getProxyProtocolConfig = function() {
  return /** @type{?proto.envoy.config.core.v3.ProxyProtocolConfig} */ (
    jspb.Message.getWrapperField(this, envoy_config_core_v3_proxy_protocol_pb.ProxyProtocolConfig, 1));
};


/** @param {?proto.envoy.config.core.v3.ProxyProtocolConfig|undefined} value */
proto.envoy.config.route.v3.RouteAction.UpgradeConfig.ConnectConfig.prototype.setProxyProtocolConfig = function(value) {
  jspb.Message.setWrapperField(this, 1, value);
};


proto.envoy.config.route.v3.RouteAction.UpgradeConfig.ConnectConfig.prototype.clearProxyProtocolConfig = function() {
  this.setProxyProtocolConfig(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.RouteAction.UpgradeConfig.ConnectConfig.prototype.hasProxyProtocolConfig = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional string upgrade_type = 1;
 * @return {string}
 */
proto.envoy.config.route.v3.RouteAction.UpgradeConfig.prototype.getUpgradeType = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.envoy.config.route.v3.RouteAction.UpgradeConfig.prototype.setUpgradeType = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional google.protobuf.BoolValue enabled = 2;
 * @return {?proto.google.protobuf.BoolValue}
 */
proto.envoy.config.route.v3.RouteAction.UpgradeConfig.prototype.getEnabled = function() {
  return /** @type{?proto.google.protobuf.BoolValue} */ (
    jspb.Message.getWrapperField(this, google_protobuf_wrappers_pb.BoolValue, 2));
};


/** @param {?proto.google.protobuf.BoolValue|undefined} value */
proto.envoy.config.route.v3.RouteAction.UpgradeConfig.prototype.setEnabled = function(value) {
  jspb.Message.setWrapperField(this, 2, value);
};


proto.envoy.config.route.v3.RouteAction.UpgradeConfig.prototype.clearEnabled = function() {
  this.setEnabled(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.RouteAction.UpgradeConfig.prototype.hasEnabled = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional ConnectConfig connect_config = 3;
 * @return {?proto.envoy.config.route.v3.RouteAction.UpgradeConfig.ConnectConfig}
 */
proto.envoy.config.route.v3.RouteAction.UpgradeConfig.prototype.getConnectConfig = function() {
  return /** @type{?proto.envoy.config.route.v3.RouteAction.UpgradeConfig.ConnectConfig} */ (
    jspb.Message.getWrapperField(this, proto.envoy.config.route.v3.RouteAction.UpgradeConfig.ConnectConfig, 3));
};


/** @param {?proto.envoy.config.route.v3.RouteAction.UpgradeConfig.ConnectConfig|undefined} value */
proto.envoy.config.route.v3.RouteAction.UpgradeConfig.prototype.setConnectConfig = function(value) {
  jspb.Message.setWrapperField(this, 3, value);
};


proto.envoy.config.route.v3.RouteAction.UpgradeConfig.prototype.clearConnectConfig = function() {
  this.setConnectConfig(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.RouteAction.UpgradeConfig.prototype.hasConnectConfig = function() {
  return jspb.Message.getField(this, 3) != null;
};


/**
 * optional string cluster = 1;
 * @return {string}
 */
proto.envoy.config.route.v3.RouteAction.prototype.getCluster = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.envoy.config.route.v3.RouteAction.prototype.setCluster = function(value) {
  jspb.Message.setOneofField(this, 1, proto.envoy.config.route.v3.RouteAction.oneofGroups_[0], value);
};


proto.envoy.config.route.v3.RouteAction.prototype.clearCluster = function() {
  jspb.Message.setOneofField(this, 1, proto.envoy.config.route.v3.RouteAction.oneofGroups_[0], undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.RouteAction.prototype.hasCluster = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional string cluster_header = 2;
 * @return {string}
 */
proto.envoy.config.route.v3.RouteAction.prototype.getClusterHeader = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.envoy.config.route.v3.RouteAction.prototype.setClusterHeader = function(value) {
  jspb.Message.setOneofField(this, 2, proto.envoy.config.route.v3.RouteAction.oneofGroups_[0], value);
};


proto.envoy.config.route.v3.RouteAction.prototype.clearClusterHeader = function() {
  jspb.Message.setOneofField(this, 2, proto.envoy.config.route.v3.RouteAction.oneofGroups_[0], undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.RouteAction.prototype.hasClusterHeader = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional WeightedCluster weighted_clusters = 3;
 * @return {?proto.envoy.config.route.v3.WeightedCluster}
 */
proto.envoy.config.route.v3.RouteAction.prototype.getWeightedClusters = function() {
  return /** @type{?proto.envoy.config.route.v3.WeightedCluster} */ (
    jspb.Message.getWrapperField(this, proto.envoy.config.route.v3.WeightedCluster, 3));
};


/** @param {?proto.envoy.config.route.v3.WeightedCluster|undefined} value */
proto.envoy.config.route.v3.RouteAction.prototype.setWeightedClusters = function(value) {
  jspb.Message.setOneofWrapperField(this, 3, proto.envoy.config.route.v3.RouteAction.oneofGroups_[0], value);
};


proto.envoy.config.route.v3.RouteAction.prototype.clearWeightedClusters = function() {
  this.setWeightedClusters(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.RouteAction.prototype.hasWeightedClusters = function() {
  return jspb.Message.getField(this, 3) != null;
};


/**
 * optional ClusterNotFoundResponseCode cluster_not_found_response_code = 20;
 * @return {!proto.envoy.config.route.v3.RouteAction.ClusterNotFoundResponseCode}
 */
proto.envoy.config.route.v3.RouteAction.prototype.getClusterNotFoundResponseCode = function() {
  return /** @type {!proto.envoy.config.route.v3.RouteAction.ClusterNotFoundResponseCode} */ (jspb.Message.getFieldWithDefault(this, 20, 0));
};


/** @param {!proto.envoy.config.route.v3.RouteAction.ClusterNotFoundResponseCode} value */
proto.envoy.config.route.v3.RouteAction.prototype.setClusterNotFoundResponseCode = function(value) {
  jspb.Message.setProto3EnumField(this, 20, value);
};


/**
 * optional envoy.config.core.v3.Metadata metadata_match = 4;
 * @return {?proto.envoy.config.core.v3.Metadata}
 */
proto.envoy.config.route.v3.RouteAction.prototype.getMetadataMatch = function() {
  return /** @type{?proto.envoy.config.core.v3.Metadata} */ (
    jspb.Message.getWrapperField(this, envoy_config_core_v3_base_pb.Metadata, 4));
};


/** @param {?proto.envoy.config.core.v3.Metadata|undefined} value */
proto.envoy.config.route.v3.RouteAction.prototype.setMetadataMatch = function(value) {
  jspb.Message.setWrapperField(this, 4, value);
};


proto.envoy.config.route.v3.RouteAction.prototype.clearMetadataMatch = function() {
  this.setMetadataMatch(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.RouteAction.prototype.hasMetadataMatch = function() {
  return jspb.Message.getField(this, 4) != null;
};


/**
 * optional string prefix_rewrite = 5;
 * @return {string}
 */
proto.envoy.config.route.v3.RouteAction.prototype.getPrefixRewrite = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 5, ""));
};


/** @param {string} value */
proto.envoy.config.route.v3.RouteAction.prototype.setPrefixRewrite = function(value) {
  jspb.Message.setProto3StringField(this, 5, value);
};


/**
 * optional envoy.type.matcher.v3.RegexMatchAndSubstitute regex_rewrite = 32;
 * @return {?proto.envoy.type.matcher.v3.RegexMatchAndSubstitute}
 */
proto.envoy.config.route.v3.RouteAction.prototype.getRegexRewrite = function() {
  return /** @type{?proto.envoy.type.matcher.v3.RegexMatchAndSubstitute} */ (
    jspb.Message.getWrapperField(this, envoy_type_matcher_v3_regex_pb.RegexMatchAndSubstitute, 32));
};


/** @param {?proto.envoy.type.matcher.v3.RegexMatchAndSubstitute|undefined} value */
proto.envoy.config.route.v3.RouteAction.prototype.setRegexRewrite = function(value) {
  jspb.Message.setWrapperField(this, 32, value);
};


proto.envoy.config.route.v3.RouteAction.prototype.clearRegexRewrite = function() {
  this.setRegexRewrite(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.RouteAction.prototype.hasRegexRewrite = function() {
  return jspb.Message.getField(this, 32) != null;
};


/**
 * optional string host_rewrite_literal = 6;
 * @return {string}
 */
proto.envoy.config.route.v3.RouteAction.prototype.getHostRewriteLiteral = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 6, ""));
};


/** @param {string} value */
proto.envoy.config.route.v3.RouteAction.prototype.setHostRewriteLiteral = function(value) {
  jspb.Message.setOneofField(this, 6, proto.envoy.config.route.v3.RouteAction.oneofGroups_[1], value);
};


proto.envoy.config.route.v3.RouteAction.prototype.clearHostRewriteLiteral = function() {
  jspb.Message.setOneofField(this, 6, proto.envoy.config.route.v3.RouteAction.oneofGroups_[1], undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.RouteAction.prototype.hasHostRewriteLiteral = function() {
  return jspb.Message.getField(this, 6) != null;
};


/**
 * optional google.protobuf.BoolValue auto_host_rewrite = 7;
 * @return {?proto.google.protobuf.BoolValue}
 */
proto.envoy.config.route.v3.RouteAction.prototype.getAutoHostRewrite = function() {
  return /** @type{?proto.google.protobuf.BoolValue} */ (
    jspb.Message.getWrapperField(this, google_protobuf_wrappers_pb.BoolValue, 7));
};


/** @param {?proto.google.protobuf.BoolValue|undefined} value */
proto.envoy.config.route.v3.RouteAction.prototype.setAutoHostRewrite = function(value) {
  jspb.Message.setOneofWrapperField(this, 7, proto.envoy.config.route.v3.RouteAction.oneofGroups_[1], value);
};


proto.envoy.config.route.v3.RouteAction.prototype.clearAutoHostRewrite = function() {
  this.setAutoHostRewrite(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.RouteAction.prototype.hasAutoHostRewrite = function() {
  return jspb.Message.getField(this, 7) != null;
};


/**
 * optional string host_rewrite_header = 29;
 * @return {string}
 */
proto.envoy.config.route.v3.RouteAction.prototype.getHostRewriteHeader = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 29, ""));
};


/** @param {string} value */
proto.envoy.config.route.v3.RouteAction.prototype.setHostRewriteHeader = function(value) {
  jspb.Message.setOneofField(this, 29, proto.envoy.config.route.v3.RouteAction.oneofGroups_[1], value);
};


proto.envoy.config.route.v3.RouteAction.prototype.clearHostRewriteHeader = function() {
  jspb.Message.setOneofField(this, 29, proto.envoy.config.route.v3.RouteAction.oneofGroups_[1], undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.RouteAction.prototype.hasHostRewriteHeader = function() {
  return jspb.Message.getField(this, 29) != null;
};


/**
 * optional google.protobuf.Duration timeout = 8;
 * @return {?proto.google.protobuf.Duration}
 */
proto.envoy.config.route.v3.RouteAction.prototype.getTimeout = function() {
  return /** @type{?proto.google.protobuf.Duration} */ (
    jspb.Message.getWrapperField(this, google_protobuf_duration_pb.Duration, 8));
};


/** @param {?proto.google.protobuf.Duration|undefined} value */
proto.envoy.config.route.v3.RouteAction.prototype.setTimeout = function(value) {
  jspb.Message.setWrapperField(this, 8, value);
};


proto.envoy.config.route.v3.RouteAction.prototype.clearTimeout = function() {
  this.setTimeout(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.RouteAction.prototype.hasTimeout = function() {
  return jspb.Message.getField(this, 8) != null;
};


/**
 * optional google.protobuf.Duration idle_timeout = 24;
 * @return {?proto.google.protobuf.Duration}
 */
proto.envoy.config.route.v3.RouteAction.prototype.getIdleTimeout = function() {
  return /** @type{?proto.google.protobuf.Duration} */ (
    jspb.Message.getWrapperField(this, google_protobuf_duration_pb.Duration, 24));
};


/** @param {?proto.google.protobuf.Duration|undefined} value */
proto.envoy.config.route.v3.RouteAction.prototype.setIdleTimeout = function(value) {
  jspb.Message.setWrapperField(this, 24, value);
};


proto.envoy.config.route.v3.RouteAction.prototype.clearIdleTimeout = function() {
  this.setIdleTimeout(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.RouteAction.prototype.hasIdleTimeout = function() {
  return jspb.Message.getField(this, 24) != null;
};


/**
 * optional RetryPolicy retry_policy = 9;
 * @return {?proto.envoy.config.route.v3.RetryPolicy}
 */
proto.envoy.config.route.v3.RouteAction.prototype.getRetryPolicy = function() {
  return /** @type{?proto.envoy.config.route.v3.RetryPolicy} */ (
    jspb.Message.getWrapperField(this, proto.envoy.config.route.v3.RetryPolicy, 9));
};


/** @param {?proto.envoy.config.route.v3.RetryPolicy|undefined} value */
proto.envoy.config.route.v3.RouteAction.prototype.setRetryPolicy = function(value) {
  jspb.Message.setWrapperField(this, 9, value);
};


proto.envoy.config.route.v3.RouteAction.prototype.clearRetryPolicy = function() {
  this.setRetryPolicy(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.RouteAction.prototype.hasRetryPolicy = function() {
  return jspb.Message.getField(this, 9) != null;
};


/**
 * optional google.protobuf.Any retry_policy_typed_config = 33;
 * @return {?proto.google.protobuf.Any}
 */
proto.envoy.config.route.v3.RouteAction.prototype.getRetryPolicyTypedConfig = function() {
  return /** @type{?proto.google.protobuf.Any} */ (
    jspb.Message.getWrapperField(this, google_protobuf_any_pb.Any, 33));
};


/** @param {?proto.google.protobuf.Any|undefined} value */
proto.envoy.config.route.v3.RouteAction.prototype.setRetryPolicyTypedConfig = function(value) {
  jspb.Message.setWrapperField(this, 33, value);
};


proto.envoy.config.route.v3.RouteAction.prototype.clearRetryPolicyTypedConfig = function() {
  this.setRetryPolicyTypedConfig(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.RouteAction.prototype.hasRetryPolicyTypedConfig = function() {
  return jspb.Message.getField(this, 33) != null;
};


/**
 * repeated RequestMirrorPolicy request_mirror_policies = 30;
 * @return {!Array<!proto.envoy.config.route.v3.RouteAction.RequestMirrorPolicy>}
 */
proto.envoy.config.route.v3.RouteAction.prototype.getRequestMirrorPoliciesList = function() {
  return /** @type{!Array<!proto.envoy.config.route.v3.RouteAction.RequestMirrorPolicy>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.envoy.config.route.v3.RouteAction.RequestMirrorPolicy, 30));
};


/** @param {!Array<!proto.envoy.config.route.v3.RouteAction.RequestMirrorPolicy>} value */
proto.envoy.config.route.v3.RouteAction.prototype.setRequestMirrorPoliciesList = function(value) {
  jspb.Message.setRepeatedWrapperField(this, 30, value);
};


/**
 * @param {!proto.envoy.config.route.v3.RouteAction.RequestMirrorPolicy=} opt_value
 * @param {number=} opt_index
 * @return {!proto.envoy.config.route.v3.RouteAction.RequestMirrorPolicy}
 */
proto.envoy.config.route.v3.RouteAction.prototype.addRequestMirrorPolicies = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 30, opt_value, proto.envoy.config.route.v3.RouteAction.RequestMirrorPolicy, opt_index);
};


proto.envoy.config.route.v3.RouteAction.prototype.clearRequestMirrorPoliciesList = function() {
  this.setRequestMirrorPoliciesList([]);
};


/**
 * optional envoy.config.core.v3.RoutingPriority priority = 11;
 * @return {!proto.envoy.config.core.v3.RoutingPriority}
 */
proto.envoy.config.route.v3.RouteAction.prototype.getPriority = function() {
  return /** @type {!proto.envoy.config.core.v3.RoutingPriority} */ (jspb.Message.getFieldWithDefault(this, 11, 0));
};


/** @param {!proto.envoy.config.core.v3.RoutingPriority} value */
proto.envoy.config.route.v3.RouteAction.prototype.setPriority = function(value) {
  jspb.Message.setProto3EnumField(this, 11, value);
};


/**
 * repeated RateLimit rate_limits = 13;
 * @return {!Array<!proto.envoy.config.route.v3.RateLimit>}
 */
proto.envoy.config.route.v3.RouteAction.prototype.getRateLimitsList = function() {
  return /** @type{!Array<!proto.envoy.config.route.v3.RateLimit>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.envoy.config.route.v3.RateLimit, 13));
};


/** @param {!Array<!proto.envoy.config.route.v3.RateLimit>} value */
proto.envoy.config.route.v3.RouteAction.prototype.setRateLimitsList = function(value) {
  jspb.Message.setRepeatedWrapperField(this, 13, value);
};


/**
 * @param {!proto.envoy.config.route.v3.RateLimit=} opt_value
 * @param {number=} opt_index
 * @return {!proto.envoy.config.route.v3.RateLimit}
 */
proto.envoy.config.route.v3.RouteAction.prototype.addRateLimits = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 13, opt_value, proto.envoy.config.route.v3.RateLimit, opt_index);
};


proto.envoy.config.route.v3.RouteAction.prototype.clearRateLimitsList = function() {
  this.setRateLimitsList([]);
};


/**
 * optional google.protobuf.BoolValue include_vh_rate_limits = 14;
 * @return {?proto.google.protobuf.BoolValue}
 */
proto.envoy.config.route.v3.RouteAction.prototype.getIncludeVhRateLimits = function() {
  return /** @type{?proto.google.protobuf.BoolValue} */ (
    jspb.Message.getWrapperField(this, google_protobuf_wrappers_pb.BoolValue, 14));
};


/** @param {?proto.google.protobuf.BoolValue|undefined} value */
proto.envoy.config.route.v3.RouteAction.prototype.setIncludeVhRateLimits = function(value) {
  jspb.Message.setWrapperField(this, 14, value);
};


proto.envoy.config.route.v3.RouteAction.prototype.clearIncludeVhRateLimits = function() {
  this.setIncludeVhRateLimits(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.RouteAction.prototype.hasIncludeVhRateLimits = function() {
  return jspb.Message.getField(this, 14) != null;
};


/**
 * repeated HashPolicy hash_policy = 15;
 * @return {!Array<!proto.envoy.config.route.v3.RouteAction.HashPolicy>}
 */
proto.envoy.config.route.v3.RouteAction.prototype.getHashPolicyList = function() {
  return /** @type{!Array<!proto.envoy.config.route.v3.RouteAction.HashPolicy>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.envoy.config.route.v3.RouteAction.HashPolicy, 15));
};


/** @param {!Array<!proto.envoy.config.route.v3.RouteAction.HashPolicy>} value */
proto.envoy.config.route.v3.RouteAction.prototype.setHashPolicyList = function(value) {
  jspb.Message.setRepeatedWrapperField(this, 15, value);
};


/**
 * @param {!proto.envoy.config.route.v3.RouteAction.HashPolicy=} opt_value
 * @param {number=} opt_index
 * @return {!proto.envoy.config.route.v3.RouteAction.HashPolicy}
 */
proto.envoy.config.route.v3.RouteAction.prototype.addHashPolicy = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 15, opt_value, proto.envoy.config.route.v3.RouteAction.HashPolicy, opt_index);
};


proto.envoy.config.route.v3.RouteAction.prototype.clearHashPolicyList = function() {
  this.setHashPolicyList([]);
};


/**
 * optional CorsPolicy cors = 17;
 * @return {?proto.envoy.config.route.v3.CorsPolicy}
 */
proto.envoy.config.route.v3.RouteAction.prototype.getCors = function() {
  return /** @type{?proto.envoy.config.route.v3.CorsPolicy} */ (
    jspb.Message.getWrapperField(this, proto.envoy.config.route.v3.CorsPolicy, 17));
};


/** @param {?proto.envoy.config.route.v3.CorsPolicy|undefined} value */
proto.envoy.config.route.v3.RouteAction.prototype.setCors = function(value) {
  jspb.Message.setWrapperField(this, 17, value);
};


proto.envoy.config.route.v3.RouteAction.prototype.clearCors = function() {
  this.setCors(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.RouteAction.prototype.hasCors = function() {
  return jspb.Message.getField(this, 17) != null;
};


/**
 * optional google.protobuf.Duration max_grpc_timeout = 23;
 * @return {?proto.google.protobuf.Duration}
 */
proto.envoy.config.route.v3.RouteAction.prototype.getMaxGrpcTimeout = function() {
  return /** @type{?proto.google.protobuf.Duration} */ (
    jspb.Message.getWrapperField(this, google_protobuf_duration_pb.Duration, 23));
};


/** @param {?proto.google.protobuf.Duration|undefined} value */
proto.envoy.config.route.v3.RouteAction.prototype.setMaxGrpcTimeout = function(value) {
  jspb.Message.setWrapperField(this, 23, value);
};


proto.envoy.config.route.v3.RouteAction.prototype.clearMaxGrpcTimeout = function() {
  this.setMaxGrpcTimeout(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.RouteAction.prototype.hasMaxGrpcTimeout = function() {
  return jspb.Message.getField(this, 23) != null;
};


/**
 * optional google.protobuf.Duration grpc_timeout_offset = 28;
 * @return {?proto.google.protobuf.Duration}
 */
proto.envoy.config.route.v3.RouteAction.prototype.getGrpcTimeoutOffset = function() {
  return /** @type{?proto.google.protobuf.Duration} */ (
    jspb.Message.getWrapperField(this, google_protobuf_duration_pb.Duration, 28));
};


/** @param {?proto.google.protobuf.Duration|undefined} value */
proto.envoy.config.route.v3.RouteAction.prototype.setGrpcTimeoutOffset = function(value) {
  jspb.Message.setWrapperField(this, 28, value);
};


proto.envoy.config.route.v3.RouteAction.prototype.clearGrpcTimeoutOffset = function() {
  this.setGrpcTimeoutOffset(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.RouteAction.prototype.hasGrpcTimeoutOffset = function() {
  return jspb.Message.getField(this, 28) != null;
};


/**
 * repeated UpgradeConfig upgrade_configs = 25;
 * @return {!Array<!proto.envoy.config.route.v3.RouteAction.UpgradeConfig>}
 */
proto.envoy.config.route.v3.RouteAction.prototype.getUpgradeConfigsList = function() {
  return /** @type{!Array<!proto.envoy.config.route.v3.RouteAction.UpgradeConfig>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.envoy.config.route.v3.RouteAction.UpgradeConfig, 25));
};


/** @param {!Array<!proto.envoy.config.route.v3.RouteAction.UpgradeConfig>} value */
proto.envoy.config.route.v3.RouteAction.prototype.setUpgradeConfigsList = function(value) {
  jspb.Message.setRepeatedWrapperField(this, 25, value);
};


/**
 * @param {!proto.envoy.config.route.v3.RouteAction.UpgradeConfig=} opt_value
 * @param {number=} opt_index
 * @return {!proto.envoy.config.route.v3.RouteAction.UpgradeConfig}
 */
proto.envoy.config.route.v3.RouteAction.prototype.addUpgradeConfigs = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 25, opt_value, proto.envoy.config.route.v3.RouteAction.UpgradeConfig, opt_index);
};


proto.envoy.config.route.v3.RouteAction.prototype.clearUpgradeConfigsList = function() {
  this.setUpgradeConfigsList([]);
};


/**
 * optional InternalRedirectPolicy internal_redirect_policy = 34;
 * @return {?proto.envoy.config.route.v3.InternalRedirectPolicy}
 */
proto.envoy.config.route.v3.RouteAction.prototype.getInternalRedirectPolicy = function() {
  return /** @type{?proto.envoy.config.route.v3.InternalRedirectPolicy} */ (
    jspb.Message.getWrapperField(this, proto.envoy.config.route.v3.InternalRedirectPolicy, 34));
};


/** @param {?proto.envoy.config.route.v3.InternalRedirectPolicy|undefined} value */
proto.envoy.config.route.v3.RouteAction.prototype.setInternalRedirectPolicy = function(value) {
  jspb.Message.setWrapperField(this, 34, value);
};


proto.envoy.config.route.v3.RouteAction.prototype.clearInternalRedirectPolicy = function() {
  this.setInternalRedirectPolicy(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.RouteAction.prototype.hasInternalRedirectPolicy = function() {
  return jspb.Message.getField(this, 34) != null;
};


/**
 * optional InternalRedirectAction internal_redirect_action = 26;
 * @return {!proto.envoy.config.route.v3.RouteAction.InternalRedirectAction}
 */
proto.envoy.config.route.v3.RouteAction.prototype.getInternalRedirectAction = function() {
  return /** @type {!proto.envoy.config.route.v3.RouteAction.InternalRedirectAction} */ (jspb.Message.getFieldWithDefault(this, 26, 0));
};


/** @param {!proto.envoy.config.route.v3.RouteAction.InternalRedirectAction} value */
proto.envoy.config.route.v3.RouteAction.prototype.setInternalRedirectAction = function(value) {
  jspb.Message.setProto3EnumField(this, 26, value);
};


/**
 * optional google.protobuf.UInt32Value max_internal_redirects = 31;
 * @return {?proto.google.protobuf.UInt32Value}
 */
proto.envoy.config.route.v3.RouteAction.prototype.getMaxInternalRedirects = function() {
  return /** @type{?proto.google.protobuf.UInt32Value} */ (
    jspb.Message.getWrapperField(this, google_protobuf_wrappers_pb.UInt32Value, 31));
};


/** @param {?proto.google.protobuf.UInt32Value|undefined} value */
proto.envoy.config.route.v3.RouteAction.prototype.setMaxInternalRedirects = function(value) {
  jspb.Message.setWrapperField(this, 31, value);
};


proto.envoy.config.route.v3.RouteAction.prototype.clearMaxInternalRedirects = function() {
  this.setMaxInternalRedirects(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.RouteAction.prototype.hasMaxInternalRedirects = function() {
  return jspb.Message.getField(this, 31) != null;
};


/**
 * optional HedgePolicy hedge_policy = 27;
 * @return {?proto.envoy.config.route.v3.HedgePolicy}
 */
proto.envoy.config.route.v3.RouteAction.prototype.getHedgePolicy = function() {
  return /** @type{?proto.envoy.config.route.v3.HedgePolicy} */ (
    jspb.Message.getWrapperField(this, proto.envoy.config.route.v3.HedgePolicy, 27));
};


/** @param {?proto.envoy.config.route.v3.HedgePolicy|undefined} value */
proto.envoy.config.route.v3.RouteAction.prototype.setHedgePolicy = function(value) {
  jspb.Message.setWrapperField(this, 27, value);
};


proto.envoy.config.route.v3.RouteAction.prototype.clearHedgePolicy = function() {
  this.setHedgePolicy(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.RouteAction.prototype.hasHedgePolicy = function() {
  return jspb.Message.getField(this, 27) != null;
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
proto.envoy.config.route.v3.RetryPolicy = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.envoy.config.route.v3.RetryPolicy.repeatedFields_, null);
};
goog.inherits(proto.envoy.config.route.v3.RetryPolicy, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.envoy.config.route.v3.RetryPolicy.displayName = 'proto.envoy.config.route.v3.RetryPolicy';
}
/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.envoy.config.route.v3.RetryPolicy.repeatedFields_ = [5,7,9,10];



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
proto.envoy.config.route.v3.RetryPolicy.prototype.toObject = function(opt_includeInstance) {
  return proto.envoy.config.route.v3.RetryPolicy.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.envoy.config.route.v3.RetryPolicy} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.RetryPolicy.toObject = function(includeInstance, msg) {
  var f, obj = {
    retryOn: jspb.Message.getFieldWithDefault(msg, 1, ""),
    numRetries: (f = msg.getNumRetries()) && google_protobuf_wrappers_pb.UInt32Value.toObject(includeInstance, f),
    perTryTimeout: (f = msg.getPerTryTimeout()) && google_protobuf_duration_pb.Duration.toObject(includeInstance, f),
    retryPriority: (f = msg.getRetryPriority()) && proto.envoy.config.route.v3.RetryPolicy.RetryPriority.toObject(includeInstance, f),
    retryHostPredicateList: jspb.Message.toObjectList(msg.getRetryHostPredicateList(),
    proto.envoy.config.route.v3.RetryPolicy.RetryHostPredicate.toObject, includeInstance),
    hostSelectionRetryMaxAttempts: jspb.Message.getFieldWithDefault(msg, 6, 0),
    retriableStatusCodesList: jspb.Message.getRepeatedField(msg, 7),
    retryBackOff: (f = msg.getRetryBackOff()) && proto.envoy.config.route.v3.RetryPolicy.RetryBackOff.toObject(includeInstance, f),
    retriableHeadersList: jspb.Message.toObjectList(msg.getRetriableHeadersList(),
    proto.envoy.config.route.v3.HeaderMatcher.toObject, includeInstance),
    retriableRequestHeadersList: jspb.Message.toObjectList(msg.getRetriableRequestHeadersList(),
    proto.envoy.config.route.v3.HeaderMatcher.toObject, includeInstance)
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
 * @return {!proto.envoy.config.route.v3.RetryPolicy}
 */
proto.envoy.config.route.v3.RetryPolicy.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.envoy.config.route.v3.RetryPolicy;
  return proto.envoy.config.route.v3.RetryPolicy.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.envoy.config.route.v3.RetryPolicy} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.envoy.config.route.v3.RetryPolicy}
 */
proto.envoy.config.route.v3.RetryPolicy.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setRetryOn(value);
      break;
    case 2:
      var value = new google_protobuf_wrappers_pb.UInt32Value;
      reader.readMessage(value,google_protobuf_wrappers_pb.UInt32Value.deserializeBinaryFromReader);
      msg.setNumRetries(value);
      break;
    case 3:
      var value = new google_protobuf_duration_pb.Duration;
      reader.readMessage(value,google_protobuf_duration_pb.Duration.deserializeBinaryFromReader);
      msg.setPerTryTimeout(value);
      break;
    case 4:
      var value = new proto.envoy.config.route.v3.RetryPolicy.RetryPriority;
      reader.readMessage(value,proto.envoy.config.route.v3.RetryPolicy.RetryPriority.deserializeBinaryFromReader);
      msg.setRetryPriority(value);
      break;
    case 5:
      var value = new proto.envoy.config.route.v3.RetryPolicy.RetryHostPredicate;
      reader.readMessage(value,proto.envoy.config.route.v3.RetryPolicy.RetryHostPredicate.deserializeBinaryFromReader);
      msg.addRetryHostPredicate(value);
      break;
    case 6:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setHostSelectionRetryMaxAttempts(value);
      break;
    case 7:
      var value = /** @type {!Array<number>} */ (reader.readPackedUint32());
      msg.setRetriableStatusCodesList(value);
      break;
    case 8:
      var value = new proto.envoy.config.route.v3.RetryPolicy.RetryBackOff;
      reader.readMessage(value,proto.envoy.config.route.v3.RetryPolicy.RetryBackOff.deserializeBinaryFromReader);
      msg.setRetryBackOff(value);
      break;
    case 9:
      var value = new proto.envoy.config.route.v3.HeaderMatcher;
      reader.readMessage(value,proto.envoy.config.route.v3.HeaderMatcher.deserializeBinaryFromReader);
      msg.addRetriableHeaders(value);
      break;
    case 10:
      var value = new proto.envoy.config.route.v3.HeaderMatcher;
      reader.readMessage(value,proto.envoy.config.route.v3.HeaderMatcher.deserializeBinaryFromReader);
      msg.addRetriableRequestHeaders(value);
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
proto.envoy.config.route.v3.RetryPolicy.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.envoy.config.route.v3.RetryPolicy.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.envoy.config.route.v3.RetryPolicy} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.RetryPolicy.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getRetryOn();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getNumRetries();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      google_protobuf_wrappers_pb.UInt32Value.serializeBinaryToWriter
    );
  }
  f = message.getPerTryTimeout();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      google_protobuf_duration_pb.Duration.serializeBinaryToWriter
    );
  }
  f = message.getRetryPriority();
  if (f != null) {
    writer.writeMessage(
      4,
      f,
      proto.envoy.config.route.v3.RetryPolicy.RetryPriority.serializeBinaryToWriter
    );
  }
  f = message.getRetryHostPredicateList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      5,
      f,
      proto.envoy.config.route.v3.RetryPolicy.RetryHostPredicate.serializeBinaryToWriter
    );
  }
  f = message.getHostSelectionRetryMaxAttempts();
  if (f !== 0) {
    writer.writeInt64(
      6,
      f
    );
  }
  f = message.getRetriableStatusCodesList();
  if (f.length > 0) {
    writer.writePackedUint32(
      7,
      f
    );
  }
  f = message.getRetryBackOff();
  if (f != null) {
    writer.writeMessage(
      8,
      f,
      proto.envoy.config.route.v3.RetryPolicy.RetryBackOff.serializeBinaryToWriter
    );
  }
  f = message.getRetriableHeadersList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      9,
      f,
      proto.envoy.config.route.v3.HeaderMatcher.serializeBinaryToWriter
    );
  }
  f = message.getRetriableRequestHeadersList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      10,
      f,
      proto.envoy.config.route.v3.HeaderMatcher.serializeBinaryToWriter
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
proto.envoy.config.route.v3.RetryPolicy.RetryPriority = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, proto.envoy.config.route.v3.RetryPolicy.RetryPriority.oneofGroups_);
};
goog.inherits(proto.envoy.config.route.v3.RetryPolicy.RetryPriority, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.envoy.config.route.v3.RetryPolicy.RetryPriority.displayName = 'proto.envoy.config.route.v3.RetryPolicy.RetryPriority';
}
/**
 * Oneof group definitions for this message. Each group defines the field
 * numbers belonging to that group. When of these fields' value is set, all
 * other fields in the group are cleared. During deserialization, if multiple
 * fields are encountered for a group, only the last value seen will be kept.
 * @private {!Array<!Array<number>>}
 * @const
 */
proto.envoy.config.route.v3.RetryPolicy.RetryPriority.oneofGroups_ = [[3]];

/**
 * @enum {number}
 */
proto.envoy.config.route.v3.RetryPolicy.RetryPriority.ConfigTypeCase = {
  CONFIG_TYPE_NOT_SET: 0,
  TYPED_CONFIG: 3
};

/**
 * @return {proto.envoy.config.route.v3.RetryPolicy.RetryPriority.ConfigTypeCase}
 */
proto.envoy.config.route.v3.RetryPolicy.RetryPriority.prototype.getConfigTypeCase = function() {
  return /** @type {proto.envoy.config.route.v3.RetryPolicy.RetryPriority.ConfigTypeCase} */(jspb.Message.computeOneofCase(this, proto.envoy.config.route.v3.RetryPolicy.RetryPriority.oneofGroups_[0]));
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
proto.envoy.config.route.v3.RetryPolicy.RetryPriority.prototype.toObject = function(opt_includeInstance) {
  return proto.envoy.config.route.v3.RetryPolicy.RetryPriority.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.envoy.config.route.v3.RetryPolicy.RetryPriority} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.RetryPolicy.RetryPriority.toObject = function(includeInstance, msg) {
  var f, obj = {
    name: jspb.Message.getFieldWithDefault(msg, 1, ""),
    typedConfig: (f = msg.getTypedConfig()) && google_protobuf_any_pb.Any.toObject(includeInstance, f)
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
 * @return {!proto.envoy.config.route.v3.RetryPolicy.RetryPriority}
 */
proto.envoy.config.route.v3.RetryPolicy.RetryPriority.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.envoy.config.route.v3.RetryPolicy.RetryPriority;
  return proto.envoy.config.route.v3.RetryPolicy.RetryPriority.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.envoy.config.route.v3.RetryPolicy.RetryPriority} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.envoy.config.route.v3.RetryPolicy.RetryPriority}
 */
proto.envoy.config.route.v3.RetryPolicy.RetryPriority.deserializeBinaryFromReader = function(msg, reader) {
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
    case 3:
      var value = new google_protobuf_any_pb.Any;
      reader.readMessage(value,google_protobuf_any_pb.Any.deserializeBinaryFromReader);
      msg.setTypedConfig(value);
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
proto.envoy.config.route.v3.RetryPolicy.RetryPriority.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.envoy.config.route.v3.RetryPolicy.RetryPriority.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.envoy.config.route.v3.RetryPolicy.RetryPriority} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.RetryPolicy.RetryPriority.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getName();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getTypedConfig();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      google_protobuf_any_pb.Any.serializeBinaryToWriter
    );
  }
};


/**
 * optional string name = 1;
 * @return {string}
 */
proto.envoy.config.route.v3.RetryPolicy.RetryPriority.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.envoy.config.route.v3.RetryPolicy.RetryPriority.prototype.setName = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional google.protobuf.Any typed_config = 3;
 * @return {?proto.google.protobuf.Any}
 */
proto.envoy.config.route.v3.RetryPolicy.RetryPriority.prototype.getTypedConfig = function() {
  return /** @type{?proto.google.protobuf.Any} */ (
    jspb.Message.getWrapperField(this, google_protobuf_any_pb.Any, 3));
};


/** @param {?proto.google.protobuf.Any|undefined} value */
proto.envoy.config.route.v3.RetryPolicy.RetryPriority.prototype.setTypedConfig = function(value) {
  jspb.Message.setOneofWrapperField(this, 3, proto.envoy.config.route.v3.RetryPolicy.RetryPriority.oneofGroups_[0], value);
};


proto.envoy.config.route.v3.RetryPolicy.RetryPriority.prototype.clearTypedConfig = function() {
  this.setTypedConfig(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.RetryPolicy.RetryPriority.prototype.hasTypedConfig = function() {
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
proto.envoy.config.route.v3.RetryPolicy.RetryHostPredicate = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, proto.envoy.config.route.v3.RetryPolicy.RetryHostPredicate.oneofGroups_);
};
goog.inherits(proto.envoy.config.route.v3.RetryPolicy.RetryHostPredicate, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.envoy.config.route.v3.RetryPolicy.RetryHostPredicate.displayName = 'proto.envoy.config.route.v3.RetryPolicy.RetryHostPredicate';
}
/**
 * Oneof group definitions for this message. Each group defines the field
 * numbers belonging to that group. When of these fields' value is set, all
 * other fields in the group are cleared. During deserialization, if multiple
 * fields are encountered for a group, only the last value seen will be kept.
 * @private {!Array<!Array<number>>}
 * @const
 */
proto.envoy.config.route.v3.RetryPolicy.RetryHostPredicate.oneofGroups_ = [[3]];

/**
 * @enum {number}
 */
proto.envoy.config.route.v3.RetryPolicy.RetryHostPredicate.ConfigTypeCase = {
  CONFIG_TYPE_NOT_SET: 0,
  TYPED_CONFIG: 3
};

/**
 * @return {proto.envoy.config.route.v3.RetryPolicy.RetryHostPredicate.ConfigTypeCase}
 */
proto.envoy.config.route.v3.RetryPolicy.RetryHostPredicate.prototype.getConfigTypeCase = function() {
  return /** @type {proto.envoy.config.route.v3.RetryPolicy.RetryHostPredicate.ConfigTypeCase} */(jspb.Message.computeOneofCase(this, proto.envoy.config.route.v3.RetryPolicy.RetryHostPredicate.oneofGroups_[0]));
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
proto.envoy.config.route.v3.RetryPolicy.RetryHostPredicate.prototype.toObject = function(opt_includeInstance) {
  return proto.envoy.config.route.v3.RetryPolicy.RetryHostPredicate.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.envoy.config.route.v3.RetryPolicy.RetryHostPredicate} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.RetryPolicy.RetryHostPredicate.toObject = function(includeInstance, msg) {
  var f, obj = {
    name: jspb.Message.getFieldWithDefault(msg, 1, ""),
    typedConfig: (f = msg.getTypedConfig()) && google_protobuf_any_pb.Any.toObject(includeInstance, f)
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
 * @return {!proto.envoy.config.route.v3.RetryPolicy.RetryHostPredicate}
 */
proto.envoy.config.route.v3.RetryPolicy.RetryHostPredicate.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.envoy.config.route.v3.RetryPolicy.RetryHostPredicate;
  return proto.envoy.config.route.v3.RetryPolicy.RetryHostPredicate.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.envoy.config.route.v3.RetryPolicy.RetryHostPredicate} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.envoy.config.route.v3.RetryPolicy.RetryHostPredicate}
 */
proto.envoy.config.route.v3.RetryPolicy.RetryHostPredicate.deserializeBinaryFromReader = function(msg, reader) {
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
    case 3:
      var value = new google_protobuf_any_pb.Any;
      reader.readMessage(value,google_protobuf_any_pb.Any.deserializeBinaryFromReader);
      msg.setTypedConfig(value);
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
proto.envoy.config.route.v3.RetryPolicy.RetryHostPredicate.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.envoy.config.route.v3.RetryPolicy.RetryHostPredicate.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.envoy.config.route.v3.RetryPolicy.RetryHostPredicate} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.RetryPolicy.RetryHostPredicate.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getName();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getTypedConfig();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      google_protobuf_any_pb.Any.serializeBinaryToWriter
    );
  }
};


/**
 * optional string name = 1;
 * @return {string}
 */
proto.envoy.config.route.v3.RetryPolicy.RetryHostPredicate.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.envoy.config.route.v3.RetryPolicy.RetryHostPredicate.prototype.setName = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional google.protobuf.Any typed_config = 3;
 * @return {?proto.google.protobuf.Any}
 */
proto.envoy.config.route.v3.RetryPolicy.RetryHostPredicate.prototype.getTypedConfig = function() {
  return /** @type{?proto.google.protobuf.Any} */ (
    jspb.Message.getWrapperField(this, google_protobuf_any_pb.Any, 3));
};


/** @param {?proto.google.protobuf.Any|undefined} value */
proto.envoy.config.route.v3.RetryPolicy.RetryHostPredicate.prototype.setTypedConfig = function(value) {
  jspb.Message.setOneofWrapperField(this, 3, proto.envoy.config.route.v3.RetryPolicy.RetryHostPredicate.oneofGroups_[0], value);
};


proto.envoy.config.route.v3.RetryPolicy.RetryHostPredicate.prototype.clearTypedConfig = function() {
  this.setTypedConfig(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.RetryPolicy.RetryHostPredicate.prototype.hasTypedConfig = function() {
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
proto.envoy.config.route.v3.RetryPolicy.RetryBackOff = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.envoy.config.route.v3.RetryPolicy.RetryBackOff, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.envoy.config.route.v3.RetryPolicy.RetryBackOff.displayName = 'proto.envoy.config.route.v3.RetryPolicy.RetryBackOff';
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
proto.envoy.config.route.v3.RetryPolicy.RetryBackOff.prototype.toObject = function(opt_includeInstance) {
  return proto.envoy.config.route.v3.RetryPolicy.RetryBackOff.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.envoy.config.route.v3.RetryPolicy.RetryBackOff} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.RetryPolicy.RetryBackOff.toObject = function(includeInstance, msg) {
  var f, obj = {
    baseInterval: (f = msg.getBaseInterval()) && google_protobuf_duration_pb.Duration.toObject(includeInstance, f),
    maxInterval: (f = msg.getMaxInterval()) && google_protobuf_duration_pb.Duration.toObject(includeInstance, f)
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
 * @return {!proto.envoy.config.route.v3.RetryPolicy.RetryBackOff}
 */
proto.envoy.config.route.v3.RetryPolicy.RetryBackOff.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.envoy.config.route.v3.RetryPolicy.RetryBackOff;
  return proto.envoy.config.route.v3.RetryPolicy.RetryBackOff.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.envoy.config.route.v3.RetryPolicy.RetryBackOff} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.envoy.config.route.v3.RetryPolicy.RetryBackOff}
 */
proto.envoy.config.route.v3.RetryPolicy.RetryBackOff.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new google_protobuf_duration_pb.Duration;
      reader.readMessage(value,google_protobuf_duration_pb.Duration.deserializeBinaryFromReader);
      msg.setBaseInterval(value);
      break;
    case 2:
      var value = new google_protobuf_duration_pb.Duration;
      reader.readMessage(value,google_protobuf_duration_pb.Duration.deserializeBinaryFromReader);
      msg.setMaxInterval(value);
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
proto.envoy.config.route.v3.RetryPolicy.RetryBackOff.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.envoy.config.route.v3.RetryPolicy.RetryBackOff.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.envoy.config.route.v3.RetryPolicy.RetryBackOff} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.RetryPolicy.RetryBackOff.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getBaseInterval();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      google_protobuf_duration_pb.Duration.serializeBinaryToWriter
    );
  }
  f = message.getMaxInterval();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      google_protobuf_duration_pb.Duration.serializeBinaryToWriter
    );
  }
};


/**
 * optional google.protobuf.Duration base_interval = 1;
 * @return {?proto.google.protobuf.Duration}
 */
proto.envoy.config.route.v3.RetryPolicy.RetryBackOff.prototype.getBaseInterval = function() {
  return /** @type{?proto.google.protobuf.Duration} */ (
    jspb.Message.getWrapperField(this, google_protobuf_duration_pb.Duration, 1));
};


/** @param {?proto.google.protobuf.Duration|undefined} value */
proto.envoy.config.route.v3.RetryPolicy.RetryBackOff.prototype.setBaseInterval = function(value) {
  jspb.Message.setWrapperField(this, 1, value);
};


proto.envoy.config.route.v3.RetryPolicy.RetryBackOff.prototype.clearBaseInterval = function() {
  this.setBaseInterval(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.RetryPolicy.RetryBackOff.prototype.hasBaseInterval = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional google.protobuf.Duration max_interval = 2;
 * @return {?proto.google.protobuf.Duration}
 */
proto.envoy.config.route.v3.RetryPolicy.RetryBackOff.prototype.getMaxInterval = function() {
  return /** @type{?proto.google.protobuf.Duration} */ (
    jspb.Message.getWrapperField(this, google_protobuf_duration_pb.Duration, 2));
};


/** @param {?proto.google.protobuf.Duration|undefined} value */
proto.envoy.config.route.v3.RetryPolicy.RetryBackOff.prototype.setMaxInterval = function(value) {
  jspb.Message.setWrapperField(this, 2, value);
};


proto.envoy.config.route.v3.RetryPolicy.RetryBackOff.prototype.clearMaxInterval = function() {
  this.setMaxInterval(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.RetryPolicy.RetryBackOff.prototype.hasMaxInterval = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional string retry_on = 1;
 * @return {string}
 */
proto.envoy.config.route.v3.RetryPolicy.prototype.getRetryOn = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.envoy.config.route.v3.RetryPolicy.prototype.setRetryOn = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional google.protobuf.UInt32Value num_retries = 2;
 * @return {?proto.google.protobuf.UInt32Value}
 */
proto.envoy.config.route.v3.RetryPolicy.prototype.getNumRetries = function() {
  return /** @type{?proto.google.protobuf.UInt32Value} */ (
    jspb.Message.getWrapperField(this, google_protobuf_wrappers_pb.UInt32Value, 2));
};


/** @param {?proto.google.protobuf.UInt32Value|undefined} value */
proto.envoy.config.route.v3.RetryPolicy.prototype.setNumRetries = function(value) {
  jspb.Message.setWrapperField(this, 2, value);
};


proto.envoy.config.route.v3.RetryPolicy.prototype.clearNumRetries = function() {
  this.setNumRetries(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.RetryPolicy.prototype.hasNumRetries = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional google.protobuf.Duration per_try_timeout = 3;
 * @return {?proto.google.protobuf.Duration}
 */
proto.envoy.config.route.v3.RetryPolicy.prototype.getPerTryTimeout = function() {
  return /** @type{?proto.google.protobuf.Duration} */ (
    jspb.Message.getWrapperField(this, google_protobuf_duration_pb.Duration, 3));
};


/** @param {?proto.google.protobuf.Duration|undefined} value */
proto.envoy.config.route.v3.RetryPolicy.prototype.setPerTryTimeout = function(value) {
  jspb.Message.setWrapperField(this, 3, value);
};


proto.envoy.config.route.v3.RetryPolicy.prototype.clearPerTryTimeout = function() {
  this.setPerTryTimeout(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.RetryPolicy.prototype.hasPerTryTimeout = function() {
  return jspb.Message.getField(this, 3) != null;
};


/**
 * optional RetryPriority retry_priority = 4;
 * @return {?proto.envoy.config.route.v3.RetryPolicy.RetryPriority}
 */
proto.envoy.config.route.v3.RetryPolicy.prototype.getRetryPriority = function() {
  return /** @type{?proto.envoy.config.route.v3.RetryPolicy.RetryPriority} */ (
    jspb.Message.getWrapperField(this, proto.envoy.config.route.v3.RetryPolicy.RetryPriority, 4));
};


/** @param {?proto.envoy.config.route.v3.RetryPolicy.RetryPriority|undefined} value */
proto.envoy.config.route.v3.RetryPolicy.prototype.setRetryPriority = function(value) {
  jspb.Message.setWrapperField(this, 4, value);
};


proto.envoy.config.route.v3.RetryPolicy.prototype.clearRetryPriority = function() {
  this.setRetryPriority(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.RetryPolicy.prototype.hasRetryPriority = function() {
  return jspb.Message.getField(this, 4) != null;
};


/**
 * repeated RetryHostPredicate retry_host_predicate = 5;
 * @return {!Array<!proto.envoy.config.route.v3.RetryPolicy.RetryHostPredicate>}
 */
proto.envoy.config.route.v3.RetryPolicy.prototype.getRetryHostPredicateList = function() {
  return /** @type{!Array<!proto.envoy.config.route.v3.RetryPolicy.RetryHostPredicate>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.envoy.config.route.v3.RetryPolicy.RetryHostPredicate, 5));
};


/** @param {!Array<!proto.envoy.config.route.v3.RetryPolicy.RetryHostPredicate>} value */
proto.envoy.config.route.v3.RetryPolicy.prototype.setRetryHostPredicateList = function(value) {
  jspb.Message.setRepeatedWrapperField(this, 5, value);
};


/**
 * @param {!proto.envoy.config.route.v3.RetryPolicy.RetryHostPredicate=} opt_value
 * @param {number=} opt_index
 * @return {!proto.envoy.config.route.v3.RetryPolicy.RetryHostPredicate}
 */
proto.envoy.config.route.v3.RetryPolicy.prototype.addRetryHostPredicate = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 5, opt_value, proto.envoy.config.route.v3.RetryPolicy.RetryHostPredicate, opt_index);
};


proto.envoy.config.route.v3.RetryPolicy.prototype.clearRetryHostPredicateList = function() {
  this.setRetryHostPredicateList([]);
};


/**
 * optional int64 host_selection_retry_max_attempts = 6;
 * @return {number}
 */
proto.envoy.config.route.v3.RetryPolicy.prototype.getHostSelectionRetryMaxAttempts = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 6, 0));
};


/** @param {number} value */
proto.envoy.config.route.v3.RetryPolicy.prototype.setHostSelectionRetryMaxAttempts = function(value) {
  jspb.Message.setProto3IntField(this, 6, value);
};


/**
 * repeated uint32 retriable_status_codes = 7;
 * @return {!Array<number>}
 */
proto.envoy.config.route.v3.RetryPolicy.prototype.getRetriableStatusCodesList = function() {
  return /** @type {!Array<number>} */ (jspb.Message.getRepeatedField(this, 7));
};


/** @param {!Array<number>} value */
proto.envoy.config.route.v3.RetryPolicy.prototype.setRetriableStatusCodesList = function(value) {
  jspb.Message.setField(this, 7, value || []);
};


/**
 * @param {!number} value
 * @param {number=} opt_index
 */
proto.envoy.config.route.v3.RetryPolicy.prototype.addRetriableStatusCodes = function(value, opt_index) {
  jspb.Message.addToRepeatedField(this, 7, value, opt_index);
};


proto.envoy.config.route.v3.RetryPolicy.prototype.clearRetriableStatusCodesList = function() {
  this.setRetriableStatusCodesList([]);
};


/**
 * optional RetryBackOff retry_back_off = 8;
 * @return {?proto.envoy.config.route.v3.RetryPolicy.RetryBackOff}
 */
proto.envoy.config.route.v3.RetryPolicy.prototype.getRetryBackOff = function() {
  return /** @type{?proto.envoy.config.route.v3.RetryPolicy.RetryBackOff} */ (
    jspb.Message.getWrapperField(this, proto.envoy.config.route.v3.RetryPolicy.RetryBackOff, 8));
};


/** @param {?proto.envoy.config.route.v3.RetryPolicy.RetryBackOff|undefined} value */
proto.envoy.config.route.v3.RetryPolicy.prototype.setRetryBackOff = function(value) {
  jspb.Message.setWrapperField(this, 8, value);
};


proto.envoy.config.route.v3.RetryPolicy.prototype.clearRetryBackOff = function() {
  this.setRetryBackOff(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.RetryPolicy.prototype.hasRetryBackOff = function() {
  return jspb.Message.getField(this, 8) != null;
};


/**
 * repeated HeaderMatcher retriable_headers = 9;
 * @return {!Array<!proto.envoy.config.route.v3.HeaderMatcher>}
 */
proto.envoy.config.route.v3.RetryPolicy.prototype.getRetriableHeadersList = function() {
  return /** @type{!Array<!proto.envoy.config.route.v3.HeaderMatcher>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.envoy.config.route.v3.HeaderMatcher, 9));
};


/** @param {!Array<!proto.envoy.config.route.v3.HeaderMatcher>} value */
proto.envoy.config.route.v3.RetryPolicy.prototype.setRetriableHeadersList = function(value) {
  jspb.Message.setRepeatedWrapperField(this, 9, value);
};


/**
 * @param {!proto.envoy.config.route.v3.HeaderMatcher=} opt_value
 * @param {number=} opt_index
 * @return {!proto.envoy.config.route.v3.HeaderMatcher}
 */
proto.envoy.config.route.v3.RetryPolicy.prototype.addRetriableHeaders = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 9, opt_value, proto.envoy.config.route.v3.HeaderMatcher, opt_index);
};


proto.envoy.config.route.v3.RetryPolicy.prototype.clearRetriableHeadersList = function() {
  this.setRetriableHeadersList([]);
};


/**
 * repeated HeaderMatcher retriable_request_headers = 10;
 * @return {!Array<!proto.envoy.config.route.v3.HeaderMatcher>}
 */
proto.envoy.config.route.v3.RetryPolicy.prototype.getRetriableRequestHeadersList = function() {
  return /** @type{!Array<!proto.envoy.config.route.v3.HeaderMatcher>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.envoy.config.route.v3.HeaderMatcher, 10));
};


/** @param {!Array<!proto.envoy.config.route.v3.HeaderMatcher>} value */
proto.envoy.config.route.v3.RetryPolicy.prototype.setRetriableRequestHeadersList = function(value) {
  jspb.Message.setRepeatedWrapperField(this, 10, value);
};


/**
 * @param {!proto.envoy.config.route.v3.HeaderMatcher=} opt_value
 * @param {number=} opt_index
 * @return {!proto.envoy.config.route.v3.HeaderMatcher}
 */
proto.envoy.config.route.v3.RetryPolicy.prototype.addRetriableRequestHeaders = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 10, opt_value, proto.envoy.config.route.v3.HeaderMatcher, opt_index);
};


proto.envoy.config.route.v3.RetryPolicy.prototype.clearRetriableRequestHeadersList = function() {
  this.setRetriableRequestHeadersList([]);
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
proto.envoy.config.route.v3.HedgePolicy = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.envoy.config.route.v3.HedgePolicy, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.envoy.config.route.v3.HedgePolicy.displayName = 'proto.envoy.config.route.v3.HedgePolicy';
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
proto.envoy.config.route.v3.HedgePolicy.prototype.toObject = function(opt_includeInstance) {
  return proto.envoy.config.route.v3.HedgePolicy.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.envoy.config.route.v3.HedgePolicy} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.HedgePolicy.toObject = function(includeInstance, msg) {
  var f, obj = {
    initialRequests: (f = msg.getInitialRequests()) && google_protobuf_wrappers_pb.UInt32Value.toObject(includeInstance, f),
    additionalRequestChance: (f = msg.getAdditionalRequestChance()) && envoy_type_v3_percent_pb.FractionalPercent.toObject(includeInstance, f),
    hedgeOnPerTryTimeout: jspb.Message.getFieldWithDefault(msg, 3, false)
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
 * @return {!proto.envoy.config.route.v3.HedgePolicy}
 */
proto.envoy.config.route.v3.HedgePolicy.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.envoy.config.route.v3.HedgePolicy;
  return proto.envoy.config.route.v3.HedgePolicy.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.envoy.config.route.v3.HedgePolicy} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.envoy.config.route.v3.HedgePolicy}
 */
proto.envoy.config.route.v3.HedgePolicy.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new google_protobuf_wrappers_pb.UInt32Value;
      reader.readMessage(value,google_protobuf_wrappers_pb.UInt32Value.deserializeBinaryFromReader);
      msg.setInitialRequests(value);
      break;
    case 2:
      var value = new envoy_type_v3_percent_pb.FractionalPercent;
      reader.readMessage(value,envoy_type_v3_percent_pb.FractionalPercent.deserializeBinaryFromReader);
      msg.setAdditionalRequestChance(value);
      break;
    case 3:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setHedgeOnPerTryTimeout(value);
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
proto.envoy.config.route.v3.HedgePolicy.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.envoy.config.route.v3.HedgePolicy.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.envoy.config.route.v3.HedgePolicy} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.HedgePolicy.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getInitialRequests();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      google_protobuf_wrappers_pb.UInt32Value.serializeBinaryToWriter
    );
  }
  f = message.getAdditionalRequestChance();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      envoy_type_v3_percent_pb.FractionalPercent.serializeBinaryToWriter
    );
  }
  f = message.getHedgeOnPerTryTimeout();
  if (f) {
    writer.writeBool(
      3,
      f
    );
  }
};


/**
 * optional google.protobuf.UInt32Value initial_requests = 1;
 * @return {?proto.google.protobuf.UInt32Value}
 */
proto.envoy.config.route.v3.HedgePolicy.prototype.getInitialRequests = function() {
  return /** @type{?proto.google.protobuf.UInt32Value} */ (
    jspb.Message.getWrapperField(this, google_protobuf_wrappers_pb.UInt32Value, 1));
};


/** @param {?proto.google.protobuf.UInt32Value|undefined} value */
proto.envoy.config.route.v3.HedgePolicy.prototype.setInitialRequests = function(value) {
  jspb.Message.setWrapperField(this, 1, value);
};


proto.envoy.config.route.v3.HedgePolicy.prototype.clearInitialRequests = function() {
  this.setInitialRequests(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.HedgePolicy.prototype.hasInitialRequests = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional envoy.type.v3.FractionalPercent additional_request_chance = 2;
 * @return {?proto.envoy.type.v3.FractionalPercent}
 */
proto.envoy.config.route.v3.HedgePolicy.prototype.getAdditionalRequestChance = function() {
  return /** @type{?proto.envoy.type.v3.FractionalPercent} */ (
    jspb.Message.getWrapperField(this, envoy_type_v3_percent_pb.FractionalPercent, 2));
};


/** @param {?proto.envoy.type.v3.FractionalPercent|undefined} value */
proto.envoy.config.route.v3.HedgePolicy.prototype.setAdditionalRequestChance = function(value) {
  jspb.Message.setWrapperField(this, 2, value);
};


proto.envoy.config.route.v3.HedgePolicy.prototype.clearAdditionalRequestChance = function() {
  this.setAdditionalRequestChance(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.HedgePolicy.prototype.hasAdditionalRequestChance = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional bool hedge_on_per_try_timeout = 3;
 * Note that Boolean fields may be set to 0/1 when serialized from a Java server.
 * You should avoid comparisons like {@code val === true/false} in those cases.
 * @return {boolean}
 */
proto.envoy.config.route.v3.HedgePolicy.prototype.getHedgeOnPerTryTimeout = function() {
  return /** @type {boolean} */ (jspb.Message.getFieldWithDefault(this, 3, false));
};


/** @param {boolean} value */
proto.envoy.config.route.v3.HedgePolicy.prototype.setHedgeOnPerTryTimeout = function(value) {
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
proto.envoy.config.route.v3.RedirectAction = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, proto.envoy.config.route.v3.RedirectAction.oneofGroups_);
};
goog.inherits(proto.envoy.config.route.v3.RedirectAction, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.envoy.config.route.v3.RedirectAction.displayName = 'proto.envoy.config.route.v3.RedirectAction';
}
/**
 * Oneof group definitions for this message. Each group defines the field
 * numbers belonging to that group. When of these fields' value is set, all
 * other fields in the group are cleared. During deserialization, if multiple
 * fields are encountered for a group, only the last value seen will be kept.
 * @private {!Array<!Array<number>>}
 * @const
 */
proto.envoy.config.route.v3.RedirectAction.oneofGroups_ = [[4,7],[2,5]];

/**
 * @enum {number}
 */
proto.envoy.config.route.v3.RedirectAction.SchemeRewriteSpecifierCase = {
  SCHEME_REWRITE_SPECIFIER_NOT_SET: 0,
  HTTPS_REDIRECT: 4,
  SCHEME_REDIRECT: 7
};

/**
 * @return {proto.envoy.config.route.v3.RedirectAction.SchemeRewriteSpecifierCase}
 */
proto.envoy.config.route.v3.RedirectAction.prototype.getSchemeRewriteSpecifierCase = function() {
  return /** @type {proto.envoy.config.route.v3.RedirectAction.SchemeRewriteSpecifierCase} */(jspb.Message.computeOneofCase(this, proto.envoy.config.route.v3.RedirectAction.oneofGroups_[0]));
};

/**
 * @enum {number}
 */
proto.envoy.config.route.v3.RedirectAction.PathRewriteSpecifierCase = {
  PATH_REWRITE_SPECIFIER_NOT_SET: 0,
  PATH_REDIRECT: 2,
  PREFIX_REWRITE: 5
};

/**
 * @return {proto.envoy.config.route.v3.RedirectAction.PathRewriteSpecifierCase}
 */
proto.envoy.config.route.v3.RedirectAction.prototype.getPathRewriteSpecifierCase = function() {
  return /** @type {proto.envoy.config.route.v3.RedirectAction.PathRewriteSpecifierCase} */(jspb.Message.computeOneofCase(this, proto.envoy.config.route.v3.RedirectAction.oneofGroups_[1]));
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
proto.envoy.config.route.v3.RedirectAction.prototype.toObject = function(opt_includeInstance) {
  return proto.envoy.config.route.v3.RedirectAction.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.envoy.config.route.v3.RedirectAction} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.RedirectAction.toObject = function(includeInstance, msg) {
  var f, obj = {
    httpsRedirect: jspb.Message.getFieldWithDefault(msg, 4, false),
    schemeRedirect: jspb.Message.getFieldWithDefault(msg, 7, ""),
    hostRedirect: jspb.Message.getFieldWithDefault(msg, 1, ""),
    portRedirect: jspb.Message.getFieldWithDefault(msg, 8, 0),
    pathRedirect: jspb.Message.getFieldWithDefault(msg, 2, ""),
    prefixRewrite: jspb.Message.getFieldWithDefault(msg, 5, ""),
    responseCode: jspb.Message.getFieldWithDefault(msg, 3, 0),
    stripQuery: jspb.Message.getFieldWithDefault(msg, 6, false)
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
 * @return {!proto.envoy.config.route.v3.RedirectAction}
 */
proto.envoy.config.route.v3.RedirectAction.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.envoy.config.route.v3.RedirectAction;
  return proto.envoy.config.route.v3.RedirectAction.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.envoy.config.route.v3.RedirectAction} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.envoy.config.route.v3.RedirectAction}
 */
proto.envoy.config.route.v3.RedirectAction.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 4:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setHttpsRedirect(value);
      break;
    case 7:
      var value = /** @type {string} */ (reader.readString());
      msg.setSchemeRedirect(value);
      break;
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setHostRedirect(value);
      break;
    case 8:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setPortRedirect(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setPathRedirect(value);
      break;
    case 5:
      var value = /** @type {string} */ (reader.readString());
      msg.setPrefixRewrite(value);
      break;
    case 3:
      var value = /** @type {!proto.envoy.config.route.v3.RedirectAction.RedirectResponseCode} */ (reader.readEnum());
      msg.setResponseCode(value);
      break;
    case 6:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setStripQuery(value);
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
proto.envoy.config.route.v3.RedirectAction.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.envoy.config.route.v3.RedirectAction.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.envoy.config.route.v3.RedirectAction} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.RedirectAction.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = /** @type {boolean} */ (jspb.Message.getField(message, 4));
  if (f != null) {
    writer.writeBool(
      4,
      f
    );
  }
  f = /** @type {string} */ (jspb.Message.getField(message, 7));
  if (f != null) {
    writer.writeString(
      7,
      f
    );
  }
  f = message.getHostRedirect();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getPortRedirect();
  if (f !== 0) {
    writer.writeUint32(
      8,
      f
    );
  }
  f = /** @type {string} */ (jspb.Message.getField(message, 2));
  if (f != null) {
    writer.writeString(
      2,
      f
    );
  }
  f = /** @type {string} */ (jspb.Message.getField(message, 5));
  if (f != null) {
    writer.writeString(
      5,
      f
    );
  }
  f = message.getResponseCode();
  if (f !== 0.0) {
    writer.writeEnum(
      3,
      f
    );
  }
  f = message.getStripQuery();
  if (f) {
    writer.writeBool(
      6,
      f
    );
  }
};


/**
 * @enum {number}
 */
proto.envoy.config.route.v3.RedirectAction.RedirectResponseCode = {
  MOVED_PERMANENTLY: 0,
  FOUND: 1,
  SEE_OTHER: 2,
  TEMPORARY_REDIRECT: 3,
  PERMANENT_REDIRECT: 4
};

/**
 * optional bool https_redirect = 4;
 * Note that Boolean fields may be set to 0/1 when serialized from a Java server.
 * You should avoid comparisons like {@code val === true/false} in those cases.
 * @return {boolean}
 */
proto.envoy.config.route.v3.RedirectAction.prototype.getHttpsRedirect = function() {
  return /** @type {boolean} */ (jspb.Message.getFieldWithDefault(this, 4, false));
};


/** @param {boolean} value */
proto.envoy.config.route.v3.RedirectAction.prototype.setHttpsRedirect = function(value) {
  jspb.Message.setOneofField(this, 4, proto.envoy.config.route.v3.RedirectAction.oneofGroups_[0], value);
};


proto.envoy.config.route.v3.RedirectAction.prototype.clearHttpsRedirect = function() {
  jspb.Message.setOneofField(this, 4, proto.envoy.config.route.v3.RedirectAction.oneofGroups_[0], undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.RedirectAction.prototype.hasHttpsRedirect = function() {
  return jspb.Message.getField(this, 4) != null;
};


/**
 * optional string scheme_redirect = 7;
 * @return {string}
 */
proto.envoy.config.route.v3.RedirectAction.prototype.getSchemeRedirect = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 7, ""));
};


/** @param {string} value */
proto.envoy.config.route.v3.RedirectAction.prototype.setSchemeRedirect = function(value) {
  jspb.Message.setOneofField(this, 7, proto.envoy.config.route.v3.RedirectAction.oneofGroups_[0], value);
};


proto.envoy.config.route.v3.RedirectAction.prototype.clearSchemeRedirect = function() {
  jspb.Message.setOneofField(this, 7, proto.envoy.config.route.v3.RedirectAction.oneofGroups_[0], undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.RedirectAction.prototype.hasSchemeRedirect = function() {
  return jspb.Message.getField(this, 7) != null;
};


/**
 * optional string host_redirect = 1;
 * @return {string}
 */
proto.envoy.config.route.v3.RedirectAction.prototype.getHostRedirect = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.envoy.config.route.v3.RedirectAction.prototype.setHostRedirect = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional uint32 port_redirect = 8;
 * @return {number}
 */
proto.envoy.config.route.v3.RedirectAction.prototype.getPortRedirect = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 8, 0));
};


/** @param {number} value */
proto.envoy.config.route.v3.RedirectAction.prototype.setPortRedirect = function(value) {
  jspb.Message.setProto3IntField(this, 8, value);
};


/**
 * optional string path_redirect = 2;
 * @return {string}
 */
proto.envoy.config.route.v3.RedirectAction.prototype.getPathRedirect = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.envoy.config.route.v3.RedirectAction.prototype.setPathRedirect = function(value) {
  jspb.Message.setOneofField(this, 2, proto.envoy.config.route.v3.RedirectAction.oneofGroups_[1], value);
};


proto.envoy.config.route.v3.RedirectAction.prototype.clearPathRedirect = function() {
  jspb.Message.setOneofField(this, 2, proto.envoy.config.route.v3.RedirectAction.oneofGroups_[1], undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.RedirectAction.prototype.hasPathRedirect = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional string prefix_rewrite = 5;
 * @return {string}
 */
proto.envoy.config.route.v3.RedirectAction.prototype.getPrefixRewrite = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 5, ""));
};


/** @param {string} value */
proto.envoy.config.route.v3.RedirectAction.prototype.setPrefixRewrite = function(value) {
  jspb.Message.setOneofField(this, 5, proto.envoy.config.route.v3.RedirectAction.oneofGroups_[1], value);
};


proto.envoy.config.route.v3.RedirectAction.prototype.clearPrefixRewrite = function() {
  jspb.Message.setOneofField(this, 5, proto.envoy.config.route.v3.RedirectAction.oneofGroups_[1], undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.RedirectAction.prototype.hasPrefixRewrite = function() {
  return jspb.Message.getField(this, 5) != null;
};


/**
 * optional RedirectResponseCode response_code = 3;
 * @return {!proto.envoy.config.route.v3.RedirectAction.RedirectResponseCode}
 */
proto.envoy.config.route.v3.RedirectAction.prototype.getResponseCode = function() {
  return /** @type {!proto.envoy.config.route.v3.RedirectAction.RedirectResponseCode} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
};


/** @param {!proto.envoy.config.route.v3.RedirectAction.RedirectResponseCode} value */
proto.envoy.config.route.v3.RedirectAction.prototype.setResponseCode = function(value) {
  jspb.Message.setProto3EnumField(this, 3, value);
};


/**
 * optional bool strip_query = 6;
 * Note that Boolean fields may be set to 0/1 when serialized from a Java server.
 * You should avoid comparisons like {@code val === true/false} in those cases.
 * @return {boolean}
 */
proto.envoy.config.route.v3.RedirectAction.prototype.getStripQuery = function() {
  return /** @type {boolean} */ (jspb.Message.getFieldWithDefault(this, 6, false));
};


/** @param {boolean} value */
proto.envoy.config.route.v3.RedirectAction.prototype.setStripQuery = function(value) {
  jspb.Message.setProto3BooleanField(this, 6, value);
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
proto.envoy.config.route.v3.DirectResponseAction = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.envoy.config.route.v3.DirectResponseAction, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.envoy.config.route.v3.DirectResponseAction.displayName = 'proto.envoy.config.route.v3.DirectResponseAction';
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
proto.envoy.config.route.v3.DirectResponseAction.prototype.toObject = function(opt_includeInstance) {
  return proto.envoy.config.route.v3.DirectResponseAction.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.envoy.config.route.v3.DirectResponseAction} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.DirectResponseAction.toObject = function(includeInstance, msg) {
  var f, obj = {
    status: jspb.Message.getFieldWithDefault(msg, 1, 0),
    body: (f = msg.getBody()) && envoy_config_core_v3_base_pb.DataSource.toObject(includeInstance, f)
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
 * @return {!proto.envoy.config.route.v3.DirectResponseAction}
 */
proto.envoy.config.route.v3.DirectResponseAction.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.envoy.config.route.v3.DirectResponseAction;
  return proto.envoy.config.route.v3.DirectResponseAction.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.envoy.config.route.v3.DirectResponseAction} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.envoy.config.route.v3.DirectResponseAction}
 */
proto.envoy.config.route.v3.DirectResponseAction.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setStatus(value);
      break;
    case 2:
      var value = new envoy_config_core_v3_base_pb.DataSource;
      reader.readMessage(value,envoy_config_core_v3_base_pb.DataSource.deserializeBinaryFromReader);
      msg.setBody(value);
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
proto.envoy.config.route.v3.DirectResponseAction.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.envoy.config.route.v3.DirectResponseAction.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.envoy.config.route.v3.DirectResponseAction} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.DirectResponseAction.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getStatus();
  if (f !== 0) {
    writer.writeUint32(
      1,
      f
    );
  }
  f = message.getBody();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      envoy_config_core_v3_base_pb.DataSource.serializeBinaryToWriter
    );
  }
};


/**
 * optional uint32 status = 1;
 * @return {number}
 */
proto.envoy.config.route.v3.DirectResponseAction.prototype.getStatus = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/** @param {number} value */
proto.envoy.config.route.v3.DirectResponseAction.prototype.setStatus = function(value) {
  jspb.Message.setProto3IntField(this, 1, value);
};


/**
 * optional envoy.config.core.v3.DataSource body = 2;
 * @return {?proto.envoy.config.core.v3.DataSource}
 */
proto.envoy.config.route.v3.DirectResponseAction.prototype.getBody = function() {
  return /** @type{?proto.envoy.config.core.v3.DataSource} */ (
    jspb.Message.getWrapperField(this, envoy_config_core_v3_base_pb.DataSource, 2));
};


/** @param {?proto.envoy.config.core.v3.DataSource|undefined} value */
proto.envoy.config.route.v3.DirectResponseAction.prototype.setBody = function(value) {
  jspb.Message.setWrapperField(this, 2, value);
};


proto.envoy.config.route.v3.DirectResponseAction.prototype.clearBody = function() {
  this.setBody(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.DirectResponseAction.prototype.hasBody = function() {
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
proto.envoy.config.route.v3.Decorator = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.envoy.config.route.v3.Decorator, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.envoy.config.route.v3.Decorator.displayName = 'proto.envoy.config.route.v3.Decorator';
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
proto.envoy.config.route.v3.Decorator.prototype.toObject = function(opt_includeInstance) {
  return proto.envoy.config.route.v3.Decorator.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.envoy.config.route.v3.Decorator} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.Decorator.toObject = function(includeInstance, msg) {
  var f, obj = {
    operation: jspb.Message.getFieldWithDefault(msg, 1, ""),
    propagate: (f = msg.getPropagate()) && google_protobuf_wrappers_pb.BoolValue.toObject(includeInstance, f)
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
 * @return {!proto.envoy.config.route.v3.Decorator}
 */
proto.envoy.config.route.v3.Decorator.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.envoy.config.route.v3.Decorator;
  return proto.envoy.config.route.v3.Decorator.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.envoy.config.route.v3.Decorator} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.envoy.config.route.v3.Decorator}
 */
proto.envoy.config.route.v3.Decorator.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setOperation(value);
      break;
    case 2:
      var value = new google_protobuf_wrappers_pb.BoolValue;
      reader.readMessage(value,google_protobuf_wrappers_pb.BoolValue.deserializeBinaryFromReader);
      msg.setPropagate(value);
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
proto.envoy.config.route.v3.Decorator.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.envoy.config.route.v3.Decorator.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.envoy.config.route.v3.Decorator} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.Decorator.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getOperation();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getPropagate();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      google_protobuf_wrappers_pb.BoolValue.serializeBinaryToWriter
    );
  }
};


/**
 * optional string operation = 1;
 * @return {string}
 */
proto.envoy.config.route.v3.Decorator.prototype.getOperation = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.envoy.config.route.v3.Decorator.prototype.setOperation = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional google.protobuf.BoolValue propagate = 2;
 * @return {?proto.google.protobuf.BoolValue}
 */
proto.envoy.config.route.v3.Decorator.prototype.getPropagate = function() {
  return /** @type{?proto.google.protobuf.BoolValue} */ (
    jspb.Message.getWrapperField(this, google_protobuf_wrappers_pb.BoolValue, 2));
};


/** @param {?proto.google.protobuf.BoolValue|undefined} value */
proto.envoy.config.route.v3.Decorator.prototype.setPropagate = function(value) {
  jspb.Message.setWrapperField(this, 2, value);
};


proto.envoy.config.route.v3.Decorator.prototype.clearPropagate = function() {
  this.setPropagate(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.Decorator.prototype.hasPropagate = function() {
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
proto.envoy.config.route.v3.Tracing = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.envoy.config.route.v3.Tracing.repeatedFields_, null);
};
goog.inherits(proto.envoy.config.route.v3.Tracing, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.envoy.config.route.v3.Tracing.displayName = 'proto.envoy.config.route.v3.Tracing';
}
/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.envoy.config.route.v3.Tracing.repeatedFields_ = [4];



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
proto.envoy.config.route.v3.Tracing.prototype.toObject = function(opt_includeInstance) {
  return proto.envoy.config.route.v3.Tracing.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.envoy.config.route.v3.Tracing} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.Tracing.toObject = function(includeInstance, msg) {
  var f, obj = {
    clientSampling: (f = msg.getClientSampling()) && envoy_type_v3_percent_pb.FractionalPercent.toObject(includeInstance, f),
    randomSampling: (f = msg.getRandomSampling()) && envoy_type_v3_percent_pb.FractionalPercent.toObject(includeInstance, f),
    overallSampling: (f = msg.getOverallSampling()) && envoy_type_v3_percent_pb.FractionalPercent.toObject(includeInstance, f),
    customTagsList: jspb.Message.toObjectList(msg.getCustomTagsList(),
    envoy_type_tracing_v3_custom_tag_pb.CustomTag.toObject, includeInstance)
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
 * @return {!proto.envoy.config.route.v3.Tracing}
 */
proto.envoy.config.route.v3.Tracing.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.envoy.config.route.v3.Tracing;
  return proto.envoy.config.route.v3.Tracing.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.envoy.config.route.v3.Tracing} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.envoy.config.route.v3.Tracing}
 */
proto.envoy.config.route.v3.Tracing.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new envoy_type_v3_percent_pb.FractionalPercent;
      reader.readMessage(value,envoy_type_v3_percent_pb.FractionalPercent.deserializeBinaryFromReader);
      msg.setClientSampling(value);
      break;
    case 2:
      var value = new envoy_type_v3_percent_pb.FractionalPercent;
      reader.readMessage(value,envoy_type_v3_percent_pb.FractionalPercent.deserializeBinaryFromReader);
      msg.setRandomSampling(value);
      break;
    case 3:
      var value = new envoy_type_v3_percent_pb.FractionalPercent;
      reader.readMessage(value,envoy_type_v3_percent_pb.FractionalPercent.deserializeBinaryFromReader);
      msg.setOverallSampling(value);
      break;
    case 4:
      var value = new envoy_type_tracing_v3_custom_tag_pb.CustomTag;
      reader.readMessage(value,envoy_type_tracing_v3_custom_tag_pb.CustomTag.deserializeBinaryFromReader);
      msg.addCustomTags(value);
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
proto.envoy.config.route.v3.Tracing.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.envoy.config.route.v3.Tracing.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.envoy.config.route.v3.Tracing} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.Tracing.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getClientSampling();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      envoy_type_v3_percent_pb.FractionalPercent.serializeBinaryToWriter
    );
  }
  f = message.getRandomSampling();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      envoy_type_v3_percent_pb.FractionalPercent.serializeBinaryToWriter
    );
  }
  f = message.getOverallSampling();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      envoy_type_v3_percent_pb.FractionalPercent.serializeBinaryToWriter
    );
  }
  f = message.getCustomTagsList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      4,
      f,
      envoy_type_tracing_v3_custom_tag_pb.CustomTag.serializeBinaryToWriter
    );
  }
};


/**
 * optional envoy.type.v3.FractionalPercent client_sampling = 1;
 * @return {?proto.envoy.type.v3.FractionalPercent}
 */
proto.envoy.config.route.v3.Tracing.prototype.getClientSampling = function() {
  return /** @type{?proto.envoy.type.v3.FractionalPercent} */ (
    jspb.Message.getWrapperField(this, envoy_type_v3_percent_pb.FractionalPercent, 1));
};


/** @param {?proto.envoy.type.v3.FractionalPercent|undefined} value */
proto.envoy.config.route.v3.Tracing.prototype.setClientSampling = function(value) {
  jspb.Message.setWrapperField(this, 1, value);
};


proto.envoy.config.route.v3.Tracing.prototype.clearClientSampling = function() {
  this.setClientSampling(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.Tracing.prototype.hasClientSampling = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional envoy.type.v3.FractionalPercent random_sampling = 2;
 * @return {?proto.envoy.type.v3.FractionalPercent}
 */
proto.envoy.config.route.v3.Tracing.prototype.getRandomSampling = function() {
  return /** @type{?proto.envoy.type.v3.FractionalPercent} */ (
    jspb.Message.getWrapperField(this, envoy_type_v3_percent_pb.FractionalPercent, 2));
};


/** @param {?proto.envoy.type.v3.FractionalPercent|undefined} value */
proto.envoy.config.route.v3.Tracing.prototype.setRandomSampling = function(value) {
  jspb.Message.setWrapperField(this, 2, value);
};


proto.envoy.config.route.v3.Tracing.prototype.clearRandomSampling = function() {
  this.setRandomSampling(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.Tracing.prototype.hasRandomSampling = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional envoy.type.v3.FractionalPercent overall_sampling = 3;
 * @return {?proto.envoy.type.v3.FractionalPercent}
 */
proto.envoy.config.route.v3.Tracing.prototype.getOverallSampling = function() {
  return /** @type{?proto.envoy.type.v3.FractionalPercent} */ (
    jspb.Message.getWrapperField(this, envoy_type_v3_percent_pb.FractionalPercent, 3));
};


/** @param {?proto.envoy.type.v3.FractionalPercent|undefined} value */
proto.envoy.config.route.v3.Tracing.prototype.setOverallSampling = function(value) {
  jspb.Message.setWrapperField(this, 3, value);
};


proto.envoy.config.route.v3.Tracing.prototype.clearOverallSampling = function() {
  this.setOverallSampling(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.Tracing.prototype.hasOverallSampling = function() {
  return jspb.Message.getField(this, 3) != null;
};


/**
 * repeated envoy.type.tracing.v3.CustomTag custom_tags = 4;
 * @return {!Array<!proto.envoy.type.tracing.v3.CustomTag>}
 */
proto.envoy.config.route.v3.Tracing.prototype.getCustomTagsList = function() {
  return /** @type{!Array<!proto.envoy.type.tracing.v3.CustomTag>} */ (
    jspb.Message.getRepeatedWrapperField(this, envoy_type_tracing_v3_custom_tag_pb.CustomTag, 4));
};


/** @param {!Array<!proto.envoy.type.tracing.v3.CustomTag>} value */
proto.envoy.config.route.v3.Tracing.prototype.setCustomTagsList = function(value) {
  jspb.Message.setRepeatedWrapperField(this, 4, value);
};


/**
 * @param {!proto.envoy.type.tracing.v3.CustomTag=} opt_value
 * @param {number=} opt_index
 * @return {!proto.envoy.type.tracing.v3.CustomTag}
 */
proto.envoy.config.route.v3.Tracing.prototype.addCustomTags = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 4, opt_value, proto.envoy.type.tracing.v3.CustomTag, opt_index);
};


proto.envoy.config.route.v3.Tracing.prototype.clearCustomTagsList = function() {
  this.setCustomTagsList([]);
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
proto.envoy.config.route.v3.VirtualCluster = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.envoy.config.route.v3.VirtualCluster.repeatedFields_, null);
};
goog.inherits(proto.envoy.config.route.v3.VirtualCluster, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.envoy.config.route.v3.VirtualCluster.displayName = 'proto.envoy.config.route.v3.VirtualCluster';
}
/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.envoy.config.route.v3.VirtualCluster.repeatedFields_ = [4];



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
proto.envoy.config.route.v3.VirtualCluster.prototype.toObject = function(opt_includeInstance) {
  return proto.envoy.config.route.v3.VirtualCluster.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.envoy.config.route.v3.VirtualCluster} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.VirtualCluster.toObject = function(includeInstance, msg) {
  var f, obj = {
    headersList: jspb.Message.toObjectList(msg.getHeadersList(),
    proto.envoy.config.route.v3.HeaderMatcher.toObject, includeInstance),
    name: jspb.Message.getFieldWithDefault(msg, 2, "")
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
 * @return {!proto.envoy.config.route.v3.VirtualCluster}
 */
proto.envoy.config.route.v3.VirtualCluster.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.envoy.config.route.v3.VirtualCluster;
  return proto.envoy.config.route.v3.VirtualCluster.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.envoy.config.route.v3.VirtualCluster} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.envoy.config.route.v3.VirtualCluster}
 */
proto.envoy.config.route.v3.VirtualCluster.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 4:
      var value = new proto.envoy.config.route.v3.HeaderMatcher;
      reader.readMessage(value,proto.envoy.config.route.v3.HeaderMatcher.deserializeBinaryFromReader);
      msg.addHeaders(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setName(value);
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
proto.envoy.config.route.v3.VirtualCluster.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.envoy.config.route.v3.VirtualCluster.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.envoy.config.route.v3.VirtualCluster} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.VirtualCluster.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getHeadersList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      4,
      f,
      proto.envoy.config.route.v3.HeaderMatcher.serializeBinaryToWriter
    );
  }
  f = message.getName();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
};


/**
 * repeated HeaderMatcher headers = 4;
 * @return {!Array<!proto.envoy.config.route.v3.HeaderMatcher>}
 */
proto.envoy.config.route.v3.VirtualCluster.prototype.getHeadersList = function() {
  return /** @type{!Array<!proto.envoy.config.route.v3.HeaderMatcher>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.envoy.config.route.v3.HeaderMatcher, 4));
};


/** @param {!Array<!proto.envoy.config.route.v3.HeaderMatcher>} value */
proto.envoy.config.route.v3.VirtualCluster.prototype.setHeadersList = function(value) {
  jspb.Message.setRepeatedWrapperField(this, 4, value);
};


/**
 * @param {!proto.envoy.config.route.v3.HeaderMatcher=} opt_value
 * @param {number=} opt_index
 * @return {!proto.envoy.config.route.v3.HeaderMatcher}
 */
proto.envoy.config.route.v3.VirtualCluster.prototype.addHeaders = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 4, opt_value, proto.envoy.config.route.v3.HeaderMatcher, opt_index);
};


proto.envoy.config.route.v3.VirtualCluster.prototype.clearHeadersList = function() {
  this.setHeadersList([]);
};


/**
 * optional string name = 2;
 * @return {string}
 */
proto.envoy.config.route.v3.VirtualCluster.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.envoy.config.route.v3.VirtualCluster.prototype.setName = function(value) {
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
proto.envoy.config.route.v3.RateLimit = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.envoy.config.route.v3.RateLimit.repeatedFields_, null);
};
goog.inherits(proto.envoy.config.route.v3.RateLimit, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.envoy.config.route.v3.RateLimit.displayName = 'proto.envoy.config.route.v3.RateLimit';
}
/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.envoy.config.route.v3.RateLimit.repeatedFields_ = [3];



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
proto.envoy.config.route.v3.RateLimit.prototype.toObject = function(opt_includeInstance) {
  return proto.envoy.config.route.v3.RateLimit.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.envoy.config.route.v3.RateLimit} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.RateLimit.toObject = function(includeInstance, msg) {
  var f, obj = {
    stage: (f = msg.getStage()) && google_protobuf_wrappers_pb.UInt32Value.toObject(includeInstance, f),
    disableKey: jspb.Message.getFieldWithDefault(msg, 2, ""),
    actionsList: jspb.Message.toObjectList(msg.getActionsList(),
    proto.envoy.config.route.v3.RateLimit.Action.toObject, includeInstance),
    limit: (f = msg.getLimit()) && proto.envoy.config.route.v3.RateLimit.Override.toObject(includeInstance, f)
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
 * @return {!proto.envoy.config.route.v3.RateLimit}
 */
proto.envoy.config.route.v3.RateLimit.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.envoy.config.route.v3.RateLimit;
  return proto.envoy.config.route.v3.RateLimit.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.envoy.config.route.v3.RateLimit} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.envoy.config.route.v3.RateLimit}
 */
proto.envoy.config.route.v3.RateLimit.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new google_protobuf_wrappers_pb.UInt32Value;
      reader.readMessage(value,google_protobuf_wrappers_pb.UInt32Value.deserializeBinaryFromReader);
      msg.setStage(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setDisableKey(value);
      break;
    case 3:
      var value = new proto.envoy.config.route.v3.RateLimit.Action;
      reader.readMessage(value,proto.envoy.config.route.v3.RateLimit.Action.deserializeBinaryFromReader);
      msg.addActions(value);
      break;
    case 4:
      var value = new proto.envoy.config.route.v3.RateLimit.Override;
      reader.readMessage(value,proto.envoy.config.route.v3.RateLimit.Override.deserializeBinaryFromReader);
      msg.setLimit(value);
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
proto.envoy.config.route.v3.RateLimit.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.envoy.config.route.v3.RateLimit.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.envoy.config.route.v3.RateLimit} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.RateLimit.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getStage();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      google_protobuf_wrappers_pb.UInt32Value.serializeBinaryToWriter
    );
  }
  f = message.getDisableKey();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getActionsList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      3,
      f,
      proto.envoy.config.route.v3.RateLimit.Action.serializeBinaryToWriter
    );
  }
  f = message.getLimit();
  if (f != null) {
    writer.writeMessage(
      4,
      f,
      proto.envoy.config.route.v3.RateLimit.Override.serializeBinaryToWriter
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
proto.envoy.config.route.v3.RateLimit.Action = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, proto.envoy.config.route.v3.RateLimit.Action.oneofGroups_);
};
goog.inherits(proto.envoy.config.route.v3.RateLimit.Action, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.envoy.config.route.v3.RateLimit.Action.displayName = 'proto.envoy.config.route.v3.RateLimit.Action';
}
/**
 * Oneof group definitions for this message. Each group defines the field
 * numbers belonging to that group. When of these fields' value is set, all
 * other fields in the group are cleared. During deserialization, if multiple
 * fields are encountered for a group, only the last value seen will be kept.
 * @private {!Array<!Array<number>>}
 * @const
 */
proto.envoy.config.route.v3.RateLimit.Action.oneofGroups_ = [[1,2,3,4,5,6,7]];

/**
 * @enum {number}
 */
proto.envoy.config.route.v3.RateLimit.Action.ActionSpecifierCase = {
  ACTION_SPECIFIER_NOT_SET: 0,
  SOURCE_CLUSTER: 1,
  DESTINATION_CLUSTER: 2,
  REQUEST_HEADERS: 3,
  REMOTE_ADDRESS: 4,
  GENERIC_KEY: 5,
  HEADER_VALUE_MATCH: 6,
  DYNAMIC_METADATA: 7
};

/**
 * @return {proto.envoy.config.route.v3.RateLimit.Action.ActionSpecifierCase}
 */
proto.envoy.config.route.v3.RateLimit.Action.prototype.getActionSpecifierCase = function() {
  return /** @type {proto.envoy.config.route.v3.RateLimit.Action.ActionSpecifierCase} */(jspb.Message.computeOneofCase(this, proto.envoy.config.route.v3.RateLimit.Action.oneofGroups_[0]));
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
proto.envoy.config.route.v3.RateLimit.Action.prototype.toObject = function(opt_includeInstance) {
  return proto.envoy.config.route.v3.RateLimit.Action.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.envoy.config.route.v3.RateLimit.Action} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.RateLimit.Action.toObject = function(includeInstance, msg) {
  var f, obj = {
    sourceCluster: (f = msg.getSourceCluster()) && proto.envoy.config.route.v3.RateLimit.Action.SourceCluster.toObject(includeInstance, f),
    destinationCluster: (f = msg.getDestinationCluster()) && proto.envoy.config.route.v3.RateLimit.Action.DestinationCluster.toObject(includeInstance, f),
    requestHeaders: (f = msg.getRequestHeaders()) && proto.envoy.config.route.v3.RateLimit.Action.RequestHeaders.toObject(includeInstance, f),
    remoteAddress: (f = msg.getRemoteAddress()) && proto.envoy.config.route.v3.RateLimit.Action.RemoteAddress.toObject(includeInstance, f),
    genericKey: (f = msg.getGenericKey()) && proto.envoy.config.route.v3.RateLimit.Action.GenericKey.toObject(includeInstance, f),
    headerValueMatch: (f = msg.getHeaderValueMatch()) && proto.envoy.config.route.v3.RateLimit.Action.HeaderValueMatch.toObject(includeInstance, f),
    dynamicMetadata: (f = msg.getDynamicMetadata()) && proto.envoy.config.route.v3.RateLimit.Action.DynamicMetaData.toObject(includeInstance, f)
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
 * @return {!proto.envoy.config.route.v3.RateLimit.Action}
 */
proto.envoy.config.route.v3.RateLimit.Action.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.envoy.config.route.v3.RateLimit.Action;
  return proto.envoy.config.route.v3.RateLimit.Action.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.envoy.config.route.v3.RateLimit.Action} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.envoy.config.route.v3.RateLimit.Action}
 */
proto.envoy.config.route.v3.RateLimit.Action.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.envoy.config.route.v3.RateLimit.Action.SourceCluster;
      reader.readMessage(value,proto.envoy.config.route.v3.RateLimit.Action.SourceCluster.deserializeBinaryFromReader);
      msg.setSourceCluster(value);
      break;
    case 2:
      var value = new proto.envoy.config.route.v3.RateLimit.Action.DestinationCluster;
      reader.readMessage(value,proto.envoy.config.route.v3.RateLimit.Action.DestinationCluster.deserializeBinaryFromReader);
      msg.setDestinationCluster(value);
      break;
    case 3:
      var value = new proto.envoy.config.route.v3.RateLimit.Action.RequestHeaders;
      reader.readMessage(value,proto.envoy.config.route.v3.RateLimit.Action.RequestHeaders.deserializeBinaryFromReader);
      msg.setRequestHeaders(value);
      break;
    case 4:
      var value = new proto.envoy.config.route.v3.RateLimit.Action.RemoteAddress;
      reader.readMessage(value,proto.envoy.config.route.v3.RateLimit.Action.RemoteAddress.deserializeBinaryFromReader);
      msg.setRemoteAddress(value);
      break;
    case 5:
      var value = new proto.envoy.config.route.v3.RateLimit.Action.GenericKey;
      reader.readMessage(value,proto.envoy.config.route.v3.RateLimit.Action.GenericKey.deserializeBinaryFromReader);
      msg.setGenericKey(value);
      break;
    case 6:
      var value = new proto.envoy.config.route.v3.RateLimit.Action.HeaderValueMatch;
      reader.readMessage(value,proto.envoy.config.route.v3.RateLimit.Action.HeaderValueMatch.deserializeBinaryFromReader);
      msg.setHeaderValueMatch(value);
      break;
    case 7:
      var value = new proto.envoy.config.route.v3.RateLimit.Action.DynamicMetaData;
      reader.readMessage(value,proto.envoy.config.route.v3.RateLimit.Action.DynamicMetaData.deserializeBinaryFromReader);
      msg.setDynamicMetadata(value);
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
proto.envoy.config.route.v3.RateLimit.Action.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.envoy.config.route.v3.RateLimit.Action.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.envoy.config.route.v3.RateLimit.Action} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.RateLimit.Action.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getSourceCluster();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto.envoy.config.route.v3.RateLimit.Action.SourceCluster.serializeBinaryToWriter
    );
  }
  f = message.getDestinationCluster();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.envoy.config.route.v3.RateLimit.Action.DestinationCluster.serializeBinaryToWriter
    );
  }
  f = message.getRequestHeaders();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      proto.envoy.config.route.v3.RateLimit.Action.RequestHeaders.serializeBinaryToWriter
    );
  }
  f = message.getRemoteAddress();
  if (f != null) {
    writer.writeMessage(
      4,
      f,
      proto.envoy.config.route.v3.RateLimit.Action.RemoteAddress.serializeBinaryToWriter
    );
  }
  f = message.getGenericKey();
  if (f != null) {
    writer.writeMessage(
      5,
      f,
      proto.envoy.config.route.v3.RateLimit.Action.GenericKey.serializeBinaryToWriter
    );
  }
  f = message.getHeaderValueMatch();
  if (f != null) {
    writer.writeMessage(
      6,
      f,
      proto.envoy.config.route.v3.RateLimit.Action.HeaderValueMatch.serializeBinaryToWriter
    );
  }
  f = message.getDynamicMetadata();
  if (f != null) {
    writer.writeMessage(
      7,
      f,
      proto.envoy.config.route.v3.RateLimit.Action.DynamicMetaData.serializeBinaryToWriter
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
proto.envoy.config.route.v3.RateLimit.Action.SourceCluster = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.envoy.config.route.v3.RateLimit.Action.SourceCluster, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.envoy.config.route.v3.RateLimit.Action.SourceCluster.displayName = 'proto.envoy.config.route.v3.RateLimit.Action.SourceCluster';
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
proto.envoy.config.route.v3.RateLimit.Action.SourceCluster.prototype.toObject = function(opt_includeInstance) {
  return proto.envoy.config.route.v3.RateLimit.Action.SourceCluster.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.envoy.config.route.v3.RateLimit.Action.SourceCluster} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.RateLimit.Action.SourceCluster.toObject = function(includeInstance, msg) {
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
 * @return {!proto.envoy.config.route.v3.RateLimit.Action.SourceCluster}
 */
proto.envoy.config.route.v3.RateLimit.Action.SourceCluster.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.envoy.config.route.v3.RateLimit.Action.SourceCluster;
  return proto.envoy.config.route.v3.RateLimit.Action.SourceCluster.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.envoy.config.route.v3.RateLimit.Action.SourceCluster} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.envoy.config.route.v3.RateLimit.Action.SourceCluster}
 */
proto.envoy.config.route.v3.RateLimit.Action.SourceCluster.deserializeBinaryFromReader = function(msg, reader) {
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
proto.envoy.config.route.v3.RateLimit.Action.SourceCluster.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.envoy.config.route.v3.RateLimit.Action.SourceCluster.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.envoy.config.route.v3.RateLimit.Action.SourceCluster} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.RateLimit.Action.SourceCluster.serializeBinaryToWriter = function(message, writer) {
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
proto.envoy.config.route.v3.RateLimit.Action.DestinationCluster = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.envoy.config.route.v3.RateLimit.Action.DestinationCluster, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.envoy.config.route.v3.RateLimit.Action.DestinationCluster.displayName = 'proto.envoy.config.route.v3.RateLimit.Action.DestinationCluster';
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
proto.envoy.config.route.v3.RateLimit.Action.DestinationCluster.prototype.toObject = function(opt_includeInstance) {
  return proto.envoy.config.route.v3.RateLimit.Action.DestinationCluster.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.envoy.config.route.v3.RateLimit.Action.DestinationCluster} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.RateLimit.Action.DestinationCluster.toObject = function(includeInstance, msg) {
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
 * @return {!proto.envoy.config.route.v3.RateLimit.Action.DestinationCluster}
 */
proto.envoy.config.route.v3.RateLimit.Action.DestinationCluster.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.envoy.config.route.v3.RateLimit.Action.DestinationCluster;
  return proto.envoy.config.route.v3.RateLimit.Action.DestinationCluster.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.envoy.config.route.v3.RateLimit.Action.DestinationCluster} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.envoy.config.route.v3.RateLimit.Action.DestinationCluster}
 */
proto.envoy.config.route.v3.RateLimit.Action.DestinationCluster.deserializeBinaryFromReader = function(msg, reader) {
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
proto.envoy.config.route.v3.RateLimit.Action.DestinationCluster.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.envoy.config.route.v3.RateLimit.Action.DestinationCluster.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.envoy.config.route.v3.RateLimit.Action.DestinationCluster} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.RateLimit.Action.DestinationCluster.serializeBinaryToWriter = function(message, writer) {
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
proto.envoy.config.route.v3.RateLimit.Action.RequestHeaders = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.envoy.config.route.v3.RateLimit.Action.RequestHeaders, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.envoy.config.route.v3.RateLimit.Action.RequestHeaders.displayName = 'proto.envoy.config.route.v3.RateLimit.Action.RequestHeaders';
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
proto.envoy.config.route.v3.RateLimit.Action.RequestHeaders.prototype.toObject = function(opt_includeInstance) {
  return proto.envoy.config.route.v3.RateLimit.Action.RequestHeaders.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.envoy.config.route.v3.RateLimit.Action.RequestHeaders} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.RateLimit.Action.RequestHeaders.toObject = function(includeInstance, msg) {
  var f, obj = {
    headerName: jspb.Message.getFieldWithDefault(msg, 1, ""),
    descriptorKey: jspb.Message.getFieldWithDefault(msg, 2, ""),
    skipIfAbsent: jspb.Message.getFieldWithDefault(msg, 3, false)
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
 * @return {!proto.envoy.config.route.v3.RateLimit.Action.RequestHeaders}
 */
proto.envoy.config.route.v3.RateLimit.Action.RequestHeaders.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.envoy.config.route.v3.RateLimit.Action.RequestHeaders;
  return proto.envoy.config.route.v3.RateLimit.Action.RequestHeaders.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.envoy.config.route.v3.RateLimit.Action.RequestHeaders} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.envoy.config.route.v3.RateLimit.Action.RequestHeaders}
 */
proto.envoy.config.route.v3.RateLimit.Action.RequestHeaders.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setHeaderName(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setDescriptorKey(value);
      break;
    case 3:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setSkipIfAbsent(value);
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
proto.envoy.config.route.v3.RateLimit.Action.RequestHeaders.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.envoy.config.route.v3.RateLimit.Action.RequestHeaders.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.envoy.config.route.v3.RateLimit.Action.RequestHeaders} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.RateLimit.Action.RequestHeaders.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getHeaderName();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getDescriptorKey();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getSkipIfAbsent();
  if (f) {
    writer.writeBool(
      3,
      f
    );
  }
};


/**
 * optional string header_name = 1;
 * @return {string}
 */
proto.envoy.config.route.v3.RateLimit.Action.RequestHeaders.prototype.getHeaderName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.envoy.config.route.v3.RateLimit.Action.RequestHeaders.prototype.setHeaderName = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string descriptor_key = 2;
 * @return {string}
 */
proto.envoy.config.route.v3.RateLimit.Action.RequestHeaders.prototype.getDescriptorKey = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.envoy.config.route.v3.RateLimit.Action.RequestHeaders.prototype.setDescriptorKey = function(value) {
  jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional bool skip_if_absent = 3;
 * Note that Boolean fields may be set to 0/1 when serialized from a Java server.
 * You should avoid comparisons like {@code val === true/false} in those cases.
 * @return {boolean}
 */
proto.envoy.config.route.v3.RateLimit.Action.RequestHeaders.prototype.getSkipIfAbsent = function() {
  return /** @type {boolean} */ (jspb.Message.getFieldWithDefault(this, 3, false));
};


/** @param {boolean} value */
proto.envoy.config.route.v3.RateLimit.Action.RequestHeaders.prototype.setSkipIfAbsent = function(value) {
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
proto.envoy.config.route.v3.RateLimit.Action.RemoteAddress = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.envoy.config.route.v3.RateLimit.Action.RemoteAddress, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.envoy.config.route.v3.RateLimit.Action.RemoteAddress.displayName = 'proto.envoy.config.route.v3.RateLimit.Action.RemoteAddress';
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
proto.envoy.config.route.v3.RateLimit.Action.RemoteAddress.prototype.toObject = function(opt_includeInstance) {
  return proto.envoy.config.route.v3.RateLimit.Action.RemoteAddress.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.envoy.config.route.v3.RateLimit.Action.RemoteAddress} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.RateLimit.Action.RemoteAddress.toObject = function(includeInstance, msg) {
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
 * @return {!proto.envoy.config.route.v3.RateLimit.Action.RemoteAddress}
 */
proto.envoy.config.route.v3.RateLimit.Action.RemoteAddress.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.envoy.config.route.v3.RateLimit.Action.RemoteAddress;
  return proto.envoy.config.route.v3.RateLimit.Action.RemoteAddress.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.envoy.config.route.v3.RateLimit.Action.RemoteAddress} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.envoy.config.route.v3.RateLimit.Action.RemoteAddress}
 */
proto.envoy.config.route.v3.RateLimit.Action.RemoteAddress.deserializeBinaryFromReader = function(msg, reader) {
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
proto.envoy.config.route.v3.RateLimit.Action.RemoteAddress.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.envoy.config.route.v3.RateLimit.Action.RemoteAddress.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.envoy.config.route.v3.RateLimit.Action.RemoteAddress} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.RateLimit.Action.RemoteAddress.serializeBinaryToWriter = function(message, writer) {
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
proto.envoy.config.route.v3.RateLimit.Action.GenericKey = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.envoy.config.route.v3.RateLimit.Action.GenericKey, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.envoy.config.route.v3.RateLimit.Action.GenericKey.displayName = 'proto.envoy.config.route.v3.RateLimit.Action.GenericKey';
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
proto.envoy.config.route.v3.RateLimit.Action.GenericKey.prototype.toObject = function(opt_includeInstance) {
  return proto.envoy.config.route.v3.RateLimit.Action.GenericKey.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.envoy.config.route.v3.RateLimit.Action.GenericKey} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.RateLimit.Action.GenericKey.toObject = function(includeInstance, msg) {
  var f, obj = {
    descriptorValue: jspb.Message.getFieldWithDefault(msg, 1, "")
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
 * @return {!proto.envoy.config.route.v3.RateLimit.Action.GenericKey}
 */
proto.envoy.config.route.v3.RateLimit.Action.GenericKey.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.envoy.config.route.v3.RateLimit.Action.GenericKey;
  return proto.envoy.config.route.v3.RateLimit.Action.GenericKey.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.envoy.config.route.v3.RateLimit.Action.GenericKey} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.envoy.config.route.v3.RateLimit.Action.GenericKey}
 */
proto.envoy.config.route.v3.RateLimit.Action.GenericKey.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setDescriptorValue(value);
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
proto.envoy.config.route.v3.RateLimit.Action.GenericKey.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.envoy.config.route.v3.RateLimit.Action.GenericKey.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.envoy.config.route.v3.RateLimit.Action.GenericKey} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.RateLimit.Action.GenericKey.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getDescriptorValue();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string descriptor_value = 1;
 * @return {string}
 */
proto.envoy.config.route.v3.RateLimit.Action.GenericKey.prototype.getDescriptorValue = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.envoy.config.route.v3.RateLimit.Action.GenericKey.prototype.setDescriptorValue = function(value) {
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
proto.envoy.config.route.v3.RateLimit.Action.HeaderValueMatch = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.envoy.config.route.v3.RateLimit.Action.HeaderValueMatch.repeatedFields_, null);
};
goog.inherits(proto.envoy.config.route.v3.RateLimit.Action.HeaderValueMatch, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.envoy.config.route.v3.RateLimit.Action.HeaderValueMatch.displayName = 'proto.envoy.config.route.v3.RateLimit.Action.HeaderValueMatch';
}
/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.envoy.config.route.v3.RateLimit.Action.HeaderValueMatch.repeatedFields_ = [3];



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
proto.envoy.config.route.v3.RateLimit.Action.HeaderValueMatch.prototype.toObject = function(opt_includeInstance) {
  return proto.envoy.config.route.v3.RateLimit.Action.HeaderValueMatch.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.envoy.config.route.v3.RateLimit.Action.HeaderValueMatch} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.RateLimit.Action.HeaderValueMatch.toObject = function(includeInstance, msg) {
  var f, obj = {
    descriptorValue: jspb.Message.getFieldWithDefault(msg, 1, ""),
    expectMatch: (f = msg.getExpectMatch()) && google_protobuf_wrappers_pb.BoolValue.toObject(includeInstance, f),
    headersList: jspb.Message.toObjectList(msg.getHeadersList(),
    proto.envoy.config.route.v3.HeaderMatcher.toObject, includeInstance)
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
 * @return {!proto.envoy.config.route.v3.RateLimit.Action.HeaderValueMatch}
 */
proto.envoy.config.route.v3.RateLimit.Action.HeaderValueMatch.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.envoy.config.route.v3.RateLimit.Action.HeaderValueMatch;
  return proto.envoy.config.route.v3.RateLimit.Action.HeaderValueMatch.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.envoy.config.route.v3.RateLimit.Action.HeaderValueMatch} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.envoy.config.route.v3.RateLimit.Action.HeaderValueMatch}
 */
proto.envoy.config.route.v3.RateLimit.Action.HeaderValueMatch.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setDescriptorValue(value);
      break;
    case 2:
      var value = new google_protobuf_wrappers_pb.BoolValue;
      reader.readMessage(value,google_protobuf_wrappers_pb.BoolValue.deserializeBinaryFromReader);
      msg.setExpectMatch(value);
      break;
    case 3:
      var value = new proto.envoy.config.route.v3.HeaderMatcher;
      reader.readMessage(value,proto.envoy.config.route.v3.HeaderMatcher.deserializeBinaryFromReader);
      msg.addHeaders(value);
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
proto.envoy.config.route.v3.RateLimit.Action.HeaderValueMatch.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.envoy.config.route.v3.RateLimit.Action.HeaderValueMatch.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.envoy.config.route.v3.RateLimit.Action.HeaderValueMatch} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.RateLimit.Action.HeaderValueMatch.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getDescriptorValue();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getExpectMatch();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      google_protobuf_wrappers_pb.BoolValue.serializeBinaryToWriter
    );
  }
  f = message.getHeadersList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      3,
      f,
      proto.envoy.config.route.v3.HeaderMatcher.serializeBinaryToWriter
    );
  }
};


/**
 * optional string descriptor_value = 1;
 * @return {string}
 */
proto.envoy.config.route.v3.RateLimit.Action.HeaderValueMatch.prototype.getDescriptorValue = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.envoy.config.route.v3.RateLimit.Action.HeaderValueMatch.prototype.setDescriptorValue = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional google.protobuf.BoolValue expect_match = 2;
 * @return {?proto.google.protobuf.BoolValue}
 */
proto.envoy.config.route.v3.RateLimit.Action.HeaderValueMatch.prototype.getExpectMatch = function() {
  return /** @type{?proto.google.protobuf.BoolValue} */ (
    jspb.Message.getWrapperField(this, google_protobuf_wrappers_pb.BoolValue, 2));
};


/** @param {?proto.google.protobuf.BoolValue|undefined} value */
proto.envoy.config.route.v3.RateLimit.Action.HeaderValueMatch.prototype.setExpectMatch = function(value) {
  jspb.Message.setWrapperField(this, 2, value);
};


proto.envoy.config.route.v3.RateLimit.Action.HeaderValueMatch.prototype.clearExpectMatch = function() {
  this.setExpectMatch(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.RateLimit.Action.HeaderValueMatch.prototype.hasExpectMatch = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * repeated HeaderMatcher headers = 3;
 * @return {!Array<!proto.envoy.config.route.v3.HeaderMatcher>}
 */
proto.envoy.config.route.v3.RateLimit.Action.HeaderValueMatch.prototype.getHeadersList = function() {
  return /** @type{!Array<!proto.envoy.config.route.v3.HeaderMatcher>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.envoy.config.route.v3.HeaderMatcher, 3));
};


/** @param {!Array<!proto.envoy.config.route.v3.HeaderMatcher>} value */
proto.envoy.config.route.v3.RateLimit.Action.HeaderValueMatch.prototype.setHeadersList = function(value) {
  jspb.Message.setRepeatedWrapperField(this, 3, value);
};


/**
 * @param {!proto.envoy.config.route.v3.HeaderMatcher=} opt_value
 * @param {number=} opt_index
 * @return {!proto.envoy.config.route.v3.HeaderMatcher}
 */
proto.envoy.config.route.v3.RateLimit.Action.HeaderValueMatch.prototype.addHeaders = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 3, opt_value, proto.envoy.config.route.v3.HeaderMatcher, opt_index);
};


proto.envoy.config.route.v3.RateLimit.Action.HeaderValueMatch.prototype.clearHeadersList = function() {
  this.setHeadersList([]);
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
proto.envoy.config.route.v3.RateLimit.Action.DynamicMetaData = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.envoy.config.route.v3.RateLimit.Action.DynamicMetaData, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.envoy.config.route.v3.RateLimit.Action.DynamicMetaData.displayName = 'proto.envoy.config.route.v3.RateLimit.Action.DynamicMetaData';
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
proto.envoy.config.route.v3.RateLimit.Action.DynamicMetaData.prototype.toObject = function(opt_includeInstance) {
  return proto.envoy.config.route.v3.RateLimit.Action.DynamicMetaData.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.envoy.config.route.v3.RateLimit.Action.DynamicMetaData} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.RateLimit.Action.DynamicMetaData.toObject = function(includeInstance, msg) {
  var f, obj = {
    descriptorKey: jspb.Message.getFieldWithDefault(msg, 1, ""),
    metadataKey: (f = msg.getMetadataKey()) && envoy_type_metadata_v3_metadata_pb.MetadataKey.toObject(includeInstance, f)
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
 * @return {!proto.envoy.config.route.v3.RateLimit.Action.DynamicMetaData}
 */
proto.envoy.config.route.v3.RateLimit.Action.DynamicMetaData.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.envoy.config.route.v3.RateLimit.Action.DynamicMetaData;
  return proto.envoy.config.route.v3.RateLimit.Action.DynamicMetaData.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.envoy.config.route.v3.RateLimit.Action.DynamicMetaData} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.envoy.config.route.v3.RateLimit.Action.DynamicMetaData}
 */
proto.envoy.config.route.v3.RateLimit.Action.DynamicMetaData.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setDescriptorKey(value);
      break;
    case 2:
      var value = new envoy_type_metadata_v3_metadata_pb.MetadataKey;
      reader.readMessage(value,envoy_type_metadata_v3_metadata_pb.MetadataKey.deserializeBinaryFromReader);
      msg.setMetadataKey(value);
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
proto.envoy.config.route.v3.RateLimit.Action.DynamicMetaData.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.envoy.config.route.v3.RateLimit.Action.DynamicMetaData.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.envoy.config.route.v3.RateLimit.Action.DynamicMetaData} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.RateLimit.Action.DynamicMetaData.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getDescriptorKey();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getMetadataKey();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      envoy_type_metadata_v3_metadata_pb.MetadataKey.serializeBinaryToWriter
    );
  }
};


/**
 * optional string descriptor_key = 1;
 * @return {string}
 */
proto.envoy.config.route.v3.RateLimit.Action.DynamicMetaData.prototype.getDescriptorKey = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.envoy.config.route.v3.RateLimit.Action.DynamicMetaData.prototype.setDescriptorKey = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional envoy.type.metadata.v3.MetadataKey metadata_key = 2;
 * @return {?proto.envoy.type.metadata.v3.MetadataKey}
 */
proto.envoy.config.route.v3.RateLimit.Action.DynamicMetaData.prototype.getMetadataKey = function() {
  return /** @type{?proto.envoy.type.metadata.v3.MetadataKey} */ (
    jspb.Message.getWrapperField(this, envoy_type_metadata_v3_metadata_pb.MetadataKey, 2));
};


/** @param {?proto.envoy.type.metadata.v3.MetadataKey|undefined} value */
proto.envoy.config.route.v3.RateLimit.Action.DynamicMetaData.prototype.setMetadataKey = function(value) {
  jspb.Message.setWrapperField(this, 2, value);
};


proto.envoy.config.route.v3.RateLimit.Action.DynamicMetaData.prototype.clearMetadataKey = function() {
  this.setMetadataKey(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.RateLimit.Action.DynamicMetaData.prototype.hasMetadataKey = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional SourceCluster source_cluster = 1;
 * @return {?proto.envoy.config.route.v3.RateLimit.Action.SourceCluster}
 */
proto.envoy.config.route.v3.RateLimit.Action.prototype.getSourceCluster = function() {
  return /** @type{?proto.envoy.config.route.v3.RateLimit.Action.SourceCluster} */ (
    jspb.Message.getWrapperField(this, proto.envoy.config.route.v3.RateLimit.Action.SourceCluster, 1));
};


/** @param {?proto.envoy.config.route.v3.RateLimit.Action.SourceCluster|undefined} value */
proto.envoy.config.route.v3.RateLimit.Action.prototype.setSourceCluster = function(value) {
  jspb.Message.setOneofWrapperField(this, 1, proto.envoy.config.route.v3.RateLimit.Action.oneofGroups_[0], value);
};


proto.envoy.config.route.v3.RateLimit.Action.prototype.clearSourceCluster = function() {
  this.setSourceCluster(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.RateLimit.Action.prototype.hasSourceCluster = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional DestinationCluster destination_cluster = 2;
 * @return {?proto.envoy.config.route.v3.RateLimit.Action.DestinationCluster}
 */
proto.envoy.config.route.v3.RateLimit.Action.prototype.getDestinationCluster = function() {
  return /** @type{?proto.envoy.config.route.v3.RateLimit.Action.DestinationCluster} */ (
    jspb.Message.getWrapperField(this, proto.envoy.config.route.v3.RateLimit.Action.DestinationCluster, 2));
};


/** @param {?proto.envoy.config.route.v3.RateLimit.Action.DestinationCluster|undefined} value */
proto.envoy.config.route.v3.RateLimit.Action.prototype.setDestinationCluster = function(value) {
  jspb.Message.setOneofWrapperField(this, 2, proto.envoy.config.route.v3.RateLimit.Action.oneofGroups_[0], value);
};


proto.envoy.config.route.v3.RateLimit.Action.prototype.clearDestinationCluster = function() {
  this.setDestinationCluster(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.RateLimit.Action.prototype.hasDestinationCluster = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional RequestHeaders request_headers = 3;
 * @return {?proto.envoy.config.route.v3.RateLimit.Action.RequestHeaders}
 */
proto.envoy.config.route.v3.RateLimit.Action.prototype.getRequestHeaders = function() {
  return /** @type{?proto.envoy.config.route.v3.RateLimit.Action.RequestHeaders} */ (
    jspb.Message.getWrapperField(this, proto.envoy.config.route.v3.RateLimit.Action.RequestHeaders, 3));
};


/** @param {?proto.envoy.config.route.v3.RateLimit.Action.RequestHeaders|undefined} value */
proto.envoy.config.route.v3.RateLimit.Action.prototype.setRequestHeaders = function(value) {
  jspb.Message.setOneofWrapperField(this, 3, proto.envoy.config.route.v3.RateLimit.Action.oneofGroups_[0], value);
};


proto.envoy.config.route.v3.RateLimit.Action.prototype.clearRequestHeaders = function() {
  this.setRequestHeaders(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.RateLimit.Action.prototype.hasRequestHeaders = function() {
  return jspb.Message.getField(this, 3) != null;
};


/**
 * optional RemoteAddress remote_address = 4;
 * @return {?proto.envoy.config.route.v3.RateLimit.Action.RemoteAddress}
 */
proto.envoy.config.route.v3.RateLimit.Action.prototype.getRemoteAddress = function() {
  return /** @type{?proto.envoy.config.route.v3.RateLimit.Action.RemoteAddress} */ (
    jspb.Message.getWrapperField(this, proto.envoy.config.route.v3.RateLimit.Action.RemoteAddress, 4));
};


/** @param {?proto.envoy.config.route.v3.RateLimit.Action.RemoteAddress|undefined} value */
proto.envoy.config.route.v3.RateLimit.Action.prototype.setRemoteAddress = function(value) {
  jspb.Message.setOneofWrapperField(this, 4, proto.envoy.config.route.v3.RateLimit.Action.oneofGroups_[0], value);
};


proto.envoy.config.route.v3.RateLimit.Action.prototype.clearRemoteAddress = function() {
  this.setRemoteAddress(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.RateLimit.Action.prototype.hasRemoteAddress = function() {
  return jspb.Message.getField(this, 4) != null;
};


/**
 * optional GenericKey generic_key = 5;
 * @return {?proto.envoy.config.route.v3.RateLimit.Action.GenericKey}
 */
proto.envoy.config.route.v3.RateLimit.Action.prototype.getGenericKey = function() {
  return /** @type{?proto.envoy.config.route.v3.RateLimit.Action.GenericKey} */ (
    jspb.Message.getWrapperField(this, proto.envoy.config.route.v3.RateLimit.Action.GenericKey, 5));
};


/** @param {?proto.envoy.config.route.v3.RateLimit.Action.GenericKey|undefined} value */
proto.envoy.config.route.v3.RateLimit.Action.prototype.setGenericKey = function(value) {
  jspb.Message.setOneofWrapperField(this, 5, proto.envoy.config.route.v3.RateLimit.Action.oneofGroups_[0], value);
};


proto.envoy.config.route.v3.RateLimit.Action.prototype.clearGenericKey = function() {
  this.setGenericKey(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.RateLimit.Action.prototype.hasGenericKey = function() {
  return jspb.Message.getField(this, 5) != null;
};


/**
 * optional HeaderValueMatch header_value_match = 6;
 * @return {?proto.envoy.config.route.v3.RateLimit.Action.HeaderValueMatch}
 */
proto.envoy.config.route.v3.RateLimit.Action.prototype.getHeaderValueMatch = function() {
  return /** @type{?proto.envoy.config.route.v3.RateLimit.Action.HeaderValueMatch} */ (
    jspb.Message.getWrapperField(this, proto.envoy.config.route.v3.RateLimit.Action.HeaderValueMatch, 6));
};


/** @param {?proto.envoy.config.route.v3.RateLimit.Action.HeaderValueMatch|undefined} value */
proto.envoy.config.route.v3.RateLimit.Action.prototype.setHeaderValueMatch = function(value) {
  jspb.Message.setOneofWrapperField(this, 6, proto.envoy.config.route.v3.RateLimit.Action.oneofGroups_[0], value);
};


proto.envoy.config.route.v3.RateLimit.Action.prototype.clearHeaderValueMatch = function() {
  this.setHeaderValueMatch(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.RateLimit.Action.prototype.hasHeaderValueMatch = function() {
  return jspb.Message.getField(this, 6) != null;
};


/**
 * optional DynamicMetaData dynamic_metadata = 7;
 * @return {?proto.envoy.config.route.v3.RateLimit.Action.DynamicMetaData}
 */
proto.envoy.config.route.v3.RateLimit.Action.prototype.getDynamicMetadata = function() {
  return /** @type{?proto.envoy.config.route.v3.RateLimit.Action.DynamicMetaData} */ (
    jspb.Message.getWrapperField(this, proto.envoy.config.route.v3.RateLimit.Action.DynamicMetaData, 7));
};


/** @param {?proto.envoy.config.route.v3.RateLimit.Action.DynamicMetaData|undefined} value */
proto.envoy.config.route.v3.RateLimit.Action.prototype.setDynamicMetadata = function(value) {
  jspb.Message.setOneofWrapperField(this, 7, proto.envoy.config.route.v3.RateLimit.Action.oneofGroups_[0], value);
};


proto.envoy.config.route.v3.RateLimit.Action.prototype.clearDynamicMetadata = function() {
  this.setDynamicMetadata(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.RateLimit.Action.prototype.hasDynamicMetadata = function() {
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
proto.envoy.config.route.v3.RateLimit.Override = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, proto.envoy.config.route.v3.RateLimit.Override.oneofGroups_);
};
goog.inherits(proto.envoy.config.route.v3.RateLimit.Override, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.envoy.config.route.v3.RateLimit.Override.displayName = 'proto.envoy.config.route.v3.RateLimit.Override';
}
/**
 * Oneof group definitions for this message. Each group defines the field
 * numbers belonging to that group. When of these fields' value is set, all
 * other fields in the group are cleared. During deserialization, if multiple
 * fields are encountered for a group, only the last value seen will be kept.
 * @private {!Array<!Array<number>>}
 * @const
 */
proto.envoy.config.route.v3.RateLimit.Override.oneofGroups_ = [[1]];

/**
 * @enum {number}
 */
proto.envoy.config.route.v3.RateLimit.Override.OverrideSpecifierCase = {
  OVERRIDE_SPECIFIER_NOT_SET: 0,
  DYNAMIC_METADATA: 1
};

/**
 * @return {proto.envoy.config.route.v3.RateLimit.Override.OverrideSpecifierCase}
 */
proto.envoy.config.route.v3.RateLimit.Override.prototype.getOverrideSpecifierCase = function() {
  return /** @type {proto.envoy.config.route.v3.RateLimit.Override.OverrideSpecifierCase} */(jspb.Message.computeOneofCase(this, proto.envoy.config.route.v3.RateLimit.Override.oneofGroups_[0]));
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
proto.envoy.config.route.v3.RateLimit.Override.prototype.toObject = function(opt_includeInstance) {
  return proto.envoy.config.route.v3.RateLimit.Override.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.envoy.config.route.v3.RateLimit.Override} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.RateLimit.Override.toObject = function(includeInstance, msg) {
  var f, obj = {
    dynamicMetadata: (f = msg.getDynamicMetadata()) && proto.envoy.config.route.v3.RateLimit.Override.DynamicMetadata.toObject(includeInstance, f)
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
 * @return {!proto.envoy.config.route.v3.RateLimit.Override}
 */
proto.envoy.config.route.v3.RateLimit.Override.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.envoy.config.route.v3.RateLimit.Override;
  return proto.envoy.config.route.v3.RateLimit.Override.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.envoy.config.route.v3.RateLimit.Override} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.envoy.config.route.v3.RateLimit.Override}
 */
proto.envoy.config.route.v3.RateLimit.Override.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.envoy.config.route.v3.RateLimit.Override.DynamicMetadata;
      reader.readMessage(value,proto.envoy.config.route.v3.RateLimit.Override.DynamicMetadata.deserializeBinaryFromReader);
      msg.setDynamicMetadata(value);
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
proto.envoy.config.route.v3.RateLimit.Override.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.envoy.config.route.v3.RateLimit.Override.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.envoy.config.route.v3.RateLimit.Override} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.RateLimit.Override.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getDynamicMetadata();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto.envoy.config.route.v3.RateLimit.Override.DynamicMetadata.serializeBinaryToWriter
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
proto.envoy.config.route.v3.RateLimit.Override.DynamicMetadata = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.envoy.config.route.v3.RateLimit.Override.DynamicMetadata, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.envoy.config.route.v3.RateLimit.Override.DynamicMetadata.displayName = 'proto.envoy.config.route.v3.RateLimit.Override.DynamicMetadata';
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
proto.envoy.config.route.v3.RateLimit.Override.DynamicMetadata.prototype.toObject = function(opt_includeInstance) {
  return proto.envoy.config.route.v3.RateLimit.Override.DynamicMetadata.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.envoy.config.route.v3.RateLimit.Override.DynamicMetadata} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.RateLimit.Override.DynamicMetadata.toObject = function(includeInstance, msg) {
  var f, obj = {
    metadataKey: (f = msg.getMetadataKey()) && envoy_type_metadata_v3_metadata_pb.MetadataKey.toObject(includeInstance, f)
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
 * @return {!proto.envoy.config.route.v3.RateLimit.Override.DynamicMetadata}
 */
proto.envoy.config.route.v3.RateLimit.Override.DynamicMetadata.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.envoy.config.route.v3.RateLimit.Override.DynamicMetadata;
  return proto.envoy.config.route.v3.RateLimit.Override.DynamicMetadata.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.envoy.config.route.v3.RateLimit.Override.DynamicMetadata} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.envoy.config.route.v3.RateLimit.Override.DynamicMetadata}
 */
proto.envoy.config.route.v3.RateLimit.Override.DynamicMetadata.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new envoy_type_metadata_v3_metadata_pb.MetadataKey;
      reader.readMessage(value,envoy_type_metadata_v3_metadata_pb.MetadataKey.deserializeBinaryFromReader);
      msg.setMetadataKey(value);
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
proto.envoy.config.route.v3.RateLimit.Override.DynamicMetadata.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.envoy.config.route.v3.RateLimit.Override.DynamicMetadata.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.envoy.config.route.v3.RateLimit.Override.DynamicMetadata} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.RateLimit.Override.DynamicMetadata.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getMetadataKey();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      envoy_type_metadata_v3_metadata_pb.MetadataKey.serializeBinaryToWriter
    );
  }
};


/**
 * optional envoy.type.metadata.v3.MetadataKey metadata_key = 1;
 * @return {?proto.envoy.type.metadata.v3.MetadataKey}
 */
proto.envoy.config.route.v3.RateLimit.Override.DynamicMetadata.prototype.getMetadataKey = function() {
  return /** @type{?proto.envoy.type.metadata.v3.MetadataKey} */ (
    jspb.Message.getWrapperField(this, envoy_type_metadata_v3_metadata_pb.MetadataKey, 1));
};


/** @param {?proto.envoy.type.metadata.v3.MetadataKey|undefined} value */
proto.envoy.config.route.v3.RateLimit.Override.DynamicMetadata.prototype.setMetadataKey = function(value) {
  jspb.Message.setWrapperField(this, 1, value);
};


proto.envoy.config.route.v3.RateLimit.Override.DynamicMetadata.prototype.clearMetadataKey = function() {
  this.setMetadataKey(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.RateLimit.Override.DynamicMetadata.prototype.hasMetadataKey = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional DynamicMetadata dynamic_metadata = 1;
 * @return {?proto.envoy.config.route.v3.RateLimit.Override.DynamicMetadata}
 */
proto.envoy.config.route.v3.RateLimit.Override.prototype.getDynamicMetadata = function() {
  return /** @type{?proto.envoy.config.route.v3.RateLimit.Override.DynamicMetadata} */ (
    jspb.Message.getWrapperField(this, proto.envoy.config.route.v3.RateLimit.Override.DynamicMetadata, 1));
};


/** @param {?proto.envoy.config.route.v3.RateLimit.Override.DynamicMetadata|undefined} value */
proto.envoy.config.route.v3.RateLimit.Override.prototype.setDynamicMetadata = function(value) {
  jspb.Message.setOneofWrapperField(this, 1, proto.envoy.config.route.v3.RateLimit.Override.oneofGroups_[0], value);
};


proto.envoy.config.route.v3.RateLimit.Override.prototype.clearDynamicMetadata = function() {
  this.setDynamicMetadata(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.RateLimit.Override.prototype.hasDynamicMetadata = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional google.protobuf.UInt32Value stage = 1;
 * @return {?proto.google.protobuf.UInt32Value}
 */
proto.envoy.config.route.v3.RateLimit.prototype.getStage = function() {
  return /** @type{?proto.google.protobuf.UInt32Value} */ (
    jspb.Message.getWrapperField(this, google_protobuf_wrappers_pb.UInt32Value, 1));
};


/** @param {?proto.google.protobuf.UInt32Value|undefined} value */
proto.envoy.config.route.v3.RateLimit.prototype.setStage = function(value) {
  jspb.Message.setWrapperField(this, 1, value);
};


proto.envoy.config.route.v3.RateLimit.prototype.clearStage = function() {
  this.setStage(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.RateLimit.prototype.hasStage = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional string disable_key = 2;
 * @return {string}
 */
proto.envoy.config.route.v3.RateLimit.prototype.getDisableKey = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.envoy.config.route.v3.RateLimit.prototype.setDisableKey = function(value) {
  jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * repeated Action actions = 3;
 * @return {!Array<!proto.envoy.config.route.v3.RateLimit.Action>}
 */
proto.envoy.config.route.v3.RateLimit.prototype.getActionsList = function() {
  return /** @type{!Array<!proto.envoy.config.route.v3.RateLimit.Action>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.envoy.config.route.v3.RateLimit.Action, 3));
};


/** @param {!Array<!proto.envoy.config.route.v3.RateLimit.Action>} value */
proto.envoy.config.route.v3.RateLimit.prototype.setActionsList = function(value) {
  jspb.Message.setRepeatedWrapperField(this, 3, value);
};


/**
 * @param {!proto.envoy.config.route.v3.RateLimit.Action=} opt_value
 * @param {number=} opt_index
 * @return {!proto.envoy.config.route.v3.RateLimit.Action}
 */
proto.envoy.config.route.v3.RateLimit.prototype.addActions = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 3, opt_value, proto.envoy.config.route.v3.RateLimit.Action, opt_index);
};


proto.envoy.config.route.v3.RateLimit.prototype.clearActionsList = function() {
  this.setActionsList([]);
};


/**
 * optional Override limit = 4;
 * @return {?proto.envoy.config.route.v3.RateLimit.Override}
 */
proto.envoy.config.route.v3.RateLimit.prototype.getLimit = function() {
  return /** @type{?proto.envoy.config.route.v3.RateLimit.Override} */ (
    jspb.Message.getWrapperField(this, proto.envoy.config.route.v3.RateLimit.Override, 4));
};


/** @param {?proto.envoy.config.route.v3.RateLimit.Override|undefined} value */
proto.envoy.config.route.v3.RateLimit.prototype.setLimit = function(value) {
  jspb.Message.setWrapperField(this, 4, value);
};


proto.envoy.config.route.v3.RateLimit.prototype.clearLimit = function() {
  this.setLimit(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.RateLimit.prototype.hasLimit = function() {
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
proto.envoy.config.route.v3.HeaderMatcher = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, proto.envoy.config.route.v3.HeaderMatcher.oneofGroups_);
};
goog.inherits(proto.envoy.config.route.v3.HeaderMatcher, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.envoy.config.route.v3.HeaderMatcher.displayName = 'proto.envoy.config.route.v3.HeaderMatcher';
}
/**
 * Oneof group definitions for this message. Each group defines the field
 * numbers belonging to that group. When of these fields' value is set, all
 * other fields in the group are cleared. During deserialization, if multiple
 * fields are encountered for a group, only the last value seen will be kept.
 * @private {!Array<!Array<number>>}
 * @const
 */
proto.envoy.config.route.v3.HeaderMatcher.oneofGroups_ = [[4,11,6,7,9,10]];

/**
 * @enum {number}
 */
proto.envoy.config.route.v3.HeaderMatcher.HeaderMatchSpecifierCase = {
  HEADER_MATCH_SPECIFIER_NOT_SET: 0,
  EXACT_MATCH: 4,
  SAFE_REGEX_MATCH: 11,
  RANGE_MATCH: 6,
  PRESENT_MATCH: 7,
  PREFIX_MATCH: 9,
  SUFFIX_MATCH: 10
};

/**
 * @return {proto.envoy.config.route.v3.HeaderMatcher.HeaderMatchSpecifierCase}
 */
proto.envoy.config.route.v3.HeaderMatcher.prototype.getHeaderMatchSpecifierCase = function() {
  return /** @type {proto.envoy.config.route.v3.HeaderMatcher.HeaderMatchSpecifierCase} */(jspb.Message.computeOneofCase(this, proto.envoy.config.route.v3.HeaderMatcher.oneofGroups_[0]));
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
proto.envoy.config.route.v3.HeaderMatcher.prototype.toObject = function(opt_includeInstance) {
  return proto.envoy.config.route.v3.HeaderMatcher.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.envoy.config.route.v3.HeaderMatcher} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.HeaderMatcher.toObject = function(includeInstance, msg) {
  var f, obj = {
    name: jspb.Message.getFieldWithDefault(msg, 1, ""),
    exactMatch: jspb.Message.getFieldWithDefault(msg, 4, ""),
    safeRegexMatch: (f = msg.getSafeRegexMatch()) && envoy_type_matcher_v3_regex_pb.RegexMatcher.toObject(includeInstance, f),
    rangeMatch: (f = msg.getRangeMatch()) && envoy_type_v3_range_pb.Int64Range.toObject(includeInstance, f),
    presentMatch: jspb.Message.getFieldWithDefault(msg, 7, false),
    prefixMatch: jspb.Message.getFieldWithDefault(msg, 9, ""),
    suffixMatch: jspb.Message.getFieldWithDefault(msg, 10, ""),
    invertMatch: jspb.Message.getFieldWithDefault(msg, 8, false)
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
 * @return {!proto.envoy.config.route.v3.HeaderMatcher}
 */
proto.envoy.config.route.v3.HeaderMatcher.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.envoy.config.route.v3.HeaderMatcher;
  return proto.envoy.config.route.v3.HeaderMatcher.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.envoy.config.route.v3.HeaderMatcher} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.envoy.config.route.v3.HeaderMatcher}
 */
proto.envoy.config.route.v3.HeaderMatcher.deserializeBinaryFromReader = function(msg, reader) {
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
    case 4:
      var value = /** @type {string} */ (reader.readString());
      msg.setExactMatch(value);
      break;
    case 11:
      var value = new envoy_type_matcher_v3_regex_pb.RegexMatcher;
      reader.readMessage(value,envoy_type_matcher_v3_regex_pb.RegexMatcher.deserializeBinaryFromReader);
      msg.setSafeRegexMatch(value);
      break;
    case 6:
      var value = new envoy_type_v3_range_pb.Int64Range;
      reader.readMessage(value,envoy_type_v3_range_pb.Int64Range.deserializeBinaryFromReader);
      msg.setRangeMatch(value);
      break;
    case 7:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setPresentMatch(value);
      break;
    case 9:
      var value = /** @type {string} */ (reader.readString());
      msg.setPrefixMatch(value);
      break;
    case 10:
      var value = /** @type {string} */ (reader.readString());
      msg.setSuffixMatch(value);
      break;
    case 8:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setInvertMatch(value);
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
proto.envoy.config.route.v3.HeaderMatcher.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.envoy.config.route.v3.HeaderMatcher.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.envoy.config.route.v3.HeaderMatcher} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.HeaderMatcher.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getName();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = /** @type {string} */ (jspb.Message.getField(message, 4));
  if (f != null) {
    writer.writeString(
      4,
      f
    );
  }
  f = message.getSafeRegexMatch();
  if (f != null) {
    writer.writeMessage(
      11,
      f,
      envoy_type_matcher_v3_regex_pb.RegexMatcher.serializeBinaryToWriter
    );
  }
  f = message.getRangeMatch();
  if (f != null) {
    writer.writeMessage(
      6,
      f,
      envoy_type_v3_range_pb.Int64Range.serializeBinaryToWriter
    );
  }
  f = /** @type {boolean} */ (jspb.Message.getField(message, 7));
  if (f != null) {
    writer.writeBool(
      7,
      f
    );
  }
  f = /** @type {string} */ (jspb.Message.getField(message, 9));
  if (f != null) {
    writer.writeString(
      9,
      f
    );
  }
  f = /** @type {string} */ (jspb.Message.getField(message, 10));
  if (f != null) {
    writer.writeString(
      10,
      f
    );
  }
  f = message.getInvertMatch();
  if (f) {
    writer.writeBool(
      8,
      f
    );
  }
};


/**
 * optional string name = 1;
 * @return {string}
 */
proto.envoy.config.route.v3.HeaderMatcher.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.envoy.config.route.v3.HeaderMatcher.prototype.setName = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string exact_match = 4;
 * @return {string}
 */
proto.envoy.config.route.v3.HeaderMatcher.prototype.getExactMatch = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 4, ""));
};


/** @param {string} value */
proto.envoy.config.route.v3.HeaderMatcher.prototype.setExactMatch = function(value) {
  jspb.Message.setOneofField(this, 4, proto.envoy.config.route.v3.HeaderMatcher.oneofGroups_[0], value);
};


proto.envoy.config.route.v3.HeaderMatcher.prototype.clearExactMatch = function() {
  jspb.Message.setOneofField(this, 4, proto.envoy.config.route.v3.HeaderMatcher.oneofGroups_[0], undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.HeaderMatcher.prototype.hasExactMatch = function() {
  return jspb.Message.getField(this, 4) != null;
};


/**
 * optional envoy.type.matcher.v3.RegexMatcher safe_regex_match = 11;
 * @return {?proto.envoy.type.matcher.v3.RegexMatcher}
 */
proto.envoy.config.route.v3.HeaderMatcher.prototype.getSafeRegexMatch = function() {
  return /** @type{?proto.envoy.type.matcher.v3.RegexMatcher} */ (
    jspb.Message.getWrapperField(this, envoy_type_matcher_v3_regex_pb.RegexMatcher, 11));
};


/** @param {?proto.envoy.type.matcher.v3.RegexMatcher|undefined} value */
proto.envoy.config.route.v3.HeaderMatcher.prototype.setSafeRegexMatch = function(value) {
  jspb.Message.setOneofWrapperField(this, 11, proto.envoy.config.route.v3.HeaderMatcher.oneofGroups_[0], value);
};


proto.envoy.config.route.v3.HeaderMatcher.prototype.clearSafeRegexMatch = function() {
  this.setSafeRegexMatch(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.HeaderMatcher.prototype.hasSafeRegexMatch = function() {
  return jspb.Message.getField(this, 11) != null;
};


/**
 * optional envoy.type.v3.Int64Range range_match = 6;
 * @return {?proto.envoy.type.v3.Int64Range}
 */
proto.envoy.config.route.v3.HeaderMatcher.prototype.getRangeMatch = function() {
  return /** @type{?proto.envoy.type.v3.Int64Range} */ (
    jspb.Message.getWrapperField(this, envoy_type_v3_range_pb.Int64Range, 6));
};


/** @param {?proto.envoy.type.v3.Int64Range|undefined} value */
proto.envoy.config.route.v3.HeaderMatcher.prototype.setRangeMatch = function(value) {
  jspb.Message.setOneofWrapperField(this, 6, proto.envoy.config.route.v3.HeaderMatcher.oneofGroups_[0], value);
};


proto.envoy.config.route.v3.HeaderMatcher.prototype.clearRangeMatch = function() {
  this.setRangeMatch(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.HeaderMatcher.prototype.hasRangeMatch = function() {
  return jspb.Message.getField(this, 6) != null;
};


/**
 * optional bool present_match = 7;
 * Note that Boolean fields may be set to 0/1 when serialized from a Java server.
 * You should avoid comparisons like {@code val === true/false} in those cases.
 * @return {boolean}
 */
proto.envoy.config.route.v3.HeaderMatcher.prototype.getPresentMatch = function() {
  return /** @type {boolean} */ (jspb.Message.getFieldWithDefault(this, 7, false));
};


/** @param {boolean} value */
proto.envoy.config.route.v3.HeaderMatcher.prototype.setPresentMatch = function(value) {
  jspb.Message.setOneofField(this, 7, proto.envoy.config.route.v3.HeaderMatcher.oneofGroups_[0], value);
};


proto.envoy.config.route.v3.HeaderMatcher.prototype.clearPresentMatch = function() {
  jspb.Message.setOneofField(this, 7, proto.envoy.config.route.v3.HeaderMatcher.oneofGroups_[0], undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.HeaderMatcher.prototype.hasPresentMatch = function() {
  return jspb.Message.getField(this, 7) != null;
};


/**
 * optional string prefix_match = 9;
 * @return {string}
 */
proto.envoy.config.route.v3.HeaderMatcher.prototype.getPrefixMatch = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 9, ""));
};


/** @param {string} value */
proto.envoy.config.route.v3.HeaderMatcher.prototype.setPrefixMatch = function(value) {
  jspb.Message.setOneofField(this, 9, proto.envoy.config.route.v3.HeaderMatcher.oneofGroups_[0], value);
};


proto.envoy.config.route.v3.HeaderMatcher.prototype.clearPrefixMatch = function() {
  jspb.Message.setOneofField(this, 9, proto.envoy.config.route.v3.HeaderMatcher.oneofGroups_[0], undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.HeaderMatcher.prototype.hasPrefixMatch = function() {
  return jspb.Message.getField(this, 9) != null;
};


/**
 * optional string suffix_match = 10;
 * @return {string}
 */
proto.envoy.config.route.v3.HeaderMatcher.prototype.getSuffixMatch = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 10, ""));
};


/** @param {string} value */
proto.envoy.config.route.v3.HeaderMatcher.prototype.setSuffixMatch = function(value) {
  jspb.Message.setOneofField(this, 10, proto.envoy.config.route.v3.HeaderMatcher.oneofGroups_[0], value);
};


proto.envoy.config.route.v3.HeaderMatcher.prototype.clearSuffixMatch = function() {
  jspb.Message.setOneofField(this, 10, proto.envoy.config.route.v3.HeaderMatcher.oneofGroups_[0], undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.HeaderMatcher.prototype.hasSuffixMatch = function() {
  return jspb.Message.getField(this, 10) != null;
};


/**
 * optional bool invert_match = 8;
 * Note that Boolean fields may be set to 0/1 when serialized from a Java server.
 * You should avoid comparisons like {@code val === true/false} in those cases.
 * @return {boolean}
 */
proto.envoy.config.route.v3.HeaderMatcher.prototype.getInvertMatch = function() {
  return /** @type {boolean} */ (jspb.Message.getFieldWithDefault(this, 8, false));
};


/** @param {boolean} value */
proto.envoy.config.route.v3.HeaderMatcher.prototype.setInvertMatch = function(value) {
  jspb.Message.setProto3BooleanField(this, 8, value);
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
proto.envoy.config.route.v3.QueryParameterMatcher = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, proto.envoy.config.route.v3.QueryParameterMatcher.oneofGroups_);
};
goog.inherits(proto.envoy.config.route.v3.QueryParameterMatcher, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.envoy.config.route.v3.QueryParameterMatcher.displayName = 'proto.envoy.config.route.v3.QueryParameterMatcher';
}
/**
 * Oneof group definitions for this message. Each group defines the field
 * numbers belonging to that group. When of these fields' value is set, all
 * other fields in the group are cleared. During deserialization, if multiple
 * fields are encountered for a group, only the last value seen will be kept.
 * @private {!Array<!Array<number>>}
 * @const
 */
proto.envoy.config.route.v3.QueryParameterMatcher.oneofGroups_ = [[5,6]];

/**
 * @enum {number}
 */
proto.envoy.config.route.v3.QueryParameterMatcher.QueryParameterMatchSpecifierCase = {
  QUERY_PARAMETER_MATCH_SPECIFIER_NOT_SET: 0,
  STRING_MATCH: 5,
  PRESENT_MATCH: 6
};

/**
 * @return {proto.envoy.config.route.v3.QueryParameterMatcher.QueryParameterMatchSpecifierCase}
 */
proto.envoy.config.route.v3.QueryParameterMatcher.prototype.getQueryParameterMatchSpecifierCase = function() {
  return /** @type {proto.envoy.config.route.v3.QueryParameterMatcher.QueryParameterMatchSpecifierCase} */(jspb.Message.computeOneofCase(this, proto.envoy.config.route.v3.QueryParameterMatcher.oneofGroups_[0]));
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
proto.envoy.config.route.v3.QueryParameterMatcher.prototype.toObject = function(opt_includeInstance) {
  return proto.envoy.config.route.v3.QueryParameterMatcher.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.envoy.config.route.v3.QueryParameterMatcher} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.QueryParameterMatcher.toObject = function(includeInstance, msg) {
  var f, obj = {
    name: jspb.Message.getFieldWithDefault(msg, 1, ""),
    stringMatch: (f = msg.getStringMatch()) && envoy_type_matcher_v3_string_pb.StringMatcher.toObject(includeInstance, f),
    presentMatch: jspb.Message.getFieldWithDefault(msg, 6, false)
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
 * @return {!proto.envoy.config.route.v3.QueryParameterMatcher}
 */
proto.envoy.config.route.v3.QueryParameterMatcher.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.envoy.config.route.v3.QueryParameterMatcher;
  return proto.envoy.config.route.v3.QueryParameterMatcher.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.envoy.config.route.v3.QueryParameterMatcher} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.envoy.config.route.v3.QueryParameterMatcher}
 */
proto.envoy.config.route.v3.QueryParameterMatcher.deserializeBinaryFromReader = function(msg, reader) {
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
    case 5:
      var value = new envoy_type_matcher_v3_string_pb.StringMatcher;
      reader.readMessage(value,envoy_type_matcher_v3_string_pb.StringMatcher.deserializeBinaryFromReader);
      msg.setStringMatch(value);
      break;
    case 6:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setPresentMatch(value);
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
proto.envoy.config.route.v3.QueryParameterMatcher.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.envoy.config.route.v3.QueryParameterMatcher.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.envoy.config.route.v3.QueryParameterMatcher} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.QueryParameterMatcher.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getName();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getStringMatch();
  if (f != null) {
    writer.writeMessage(
      5,
      f,
      envoy_type_matcher_v3_string_pb.StringMatcher.serializeBinaryToWriter
    );
  }
  f = /** @type {boolean} */ (jspb.Message.getField(message, 6));
  if (f != null) {
    writer.writeBool(
      6,
      f
    );
  }
};


/**
 * optional string name = 1;
 * @return {string}
 */
proto.envoy.config.route.v3.QueryParameterMatcher.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.envoy.config.route.v3.QueryParameterMatcher.prototype.setName = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional envoy.type.matcher.v3.StringMatcher string_match = 5;
 * @return {?proto.envoy.type.matcher.v3.StringMatcher}
 */
proto.envoy.config.route.v3.QueryParameterMatcher.prototype.getStringMatch = function() {
  return /** @type{?proto.envoy.type.matcher.v3.StringMatcher} */ (
    jspb.Message.getWrapperField(this, envoy_type_matcher_v3_string_pb.StringMatcher, 5));
};


/** @param {?proto.envoy.type.matcher.v3.StringMatcher|undefined} value */
proto.envoy.config.route.v3.QueryParameterMatcher.prototype.setStringMatch = function(value) {
  jspb.Message.setOneofWrapperField(this, 5, proto.envoy.config.route.v3.QueryParameterMatcher.oneofGroups_[0], value);
};


proto.envoy.config.route.v3.QueryParameterMatcher.prototype.clearStringMatch = function() {
  this.setStringMatch(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.QueryParameterMatcher.prototype.hasStringMatch = function() {
  return jspb.Message.getField(this, 5) != null;
};


/**
 * optional bool present_match = 6;
 * Note that Boolean fields may be set to 0/1 when serialized from a Java server.
 * You should avoid comparisons like {@code val === true/false} in those cases.
 * @return {boolean}
 */
proto.envoy.config.route.v3.QueryParameterMatcher.prototype.getPresentMatch = function() {
  return /** @type {boolean} */ (jspb.Message.getFieldWithDefault(this, 6, false));
};


/** @param {boolean} value */
proto.envoy.config.route.v3.QueryParameterMatcher.prototype.setPresentMatch = function(value) {
  jspb.Message.setOneofField(this, 6, proto.envoy.config.route.v3.QueryParameterMatcher.oneofGroups_[0], value);
};


proto.envoy.config.route.v3.QueryParameterMatcher.prototype.clearPresentMatch = function() {
  jspb.Message.setOneofField(this, 6, proto.envoy.config.route.v3.QueryParameterMatcher.oneofGroups_[0], undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.QueryParameterMatcher.prototype.hasPresentMatch = function() {
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
proto.envoy.config.route.v3.InternalRedirectPolicy = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.envoy.config.route.v3.InternalRedirectPolicy.repeatedFields_, null);
};
goog.inherits(proto.envoy.config.route.v3.InternalRedirectPolicy, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.envoy.config.route.v3.InternalRedirectPolicy.displayName = 'proto.envoy.config.route.v3.InternalRedirectPolicy';
}
/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.envoy.config.route.v3.InternalRedirectPolicy.repeatedFields_ = [2,3];



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
proto.envoy.config.route.v3.InternalRedirectPolicy.prototype.toObject = function(opt_includeInstance) {
  return proto.envoy.config.route.v3.InternalRedirectPolicy.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.envoy.config.route.v3.InternalRedirectPolicy} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.InternalRedirectPolicy.toObject = function(includeInstance, msg) {
  var f, obj = {
    maxInternalRedirects: (f = msg.getMaxInternalRedirects()) && google_protobuf_wrappers_pb.UInt32Value.toObject(includeInstance, f),
    redirectResponseCodesList: jspb.Message.getRepeatedField(msg, 2),
    predicatesList: jspb.Message.toObjectList(msg.getPredicatesList(),
    envoy_config_core_v3_extension_pb.TypedExtensionConfig.toObject, includeInstance),
    allowCrossSchemeRedirect: jspb.Message.getFieldWithDefault(msg, 4, false)
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
 * @return {!proto.envoy.config.route.v3.InternalRedirectPolicy}
 */
proto.envoy.config.route.v3.InternalRedirectPolicy.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.envoy.config.route.v3.InternalRedirectPolicy;
  return proto.envoy.config.route.v3.InternalRedirectPolicy.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.envoy.config.route.v3.InternalRedirectPolicy} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.envoy.config.route.v3.InternalRedirectPolicy}
 */
proto.envoy.config.route.v3.InternalRedirectPolicy.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new google_protobuf_wrappers_pb.UInt32Value;
      reader.readMessage(value,google_protobuf_wrappers_pb.UInt32Value.deserializeBinaryFromReader);
      msg.setMaxInternalRedirects(value);
      break;
    case 2:
      var value = /** @type {!Array<number>} */ (reader.readPackedUint32());
      msg.setRedirectResponseCodesList(value);
      break;
    case 3:
      var value = new envoy_config_core_v3_extension_pb.TypedExtensionConfig;
      reader.readMessage(value,envoy_config_core_v3_extension_pb.TypedExtensionConfig.deserializeBinaryFromReader);
      msg.addPredicates(value);
      break;
    case 4:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setAllowCrossSchemeRedirect(value);
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
proto.envoy.config.route.v3.InternalRedirectPolicy.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.envoy.config.route.v3.InternalRedirectPolicy.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.envoy.config.route.v3.InternalRedirectPolicy} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.envoy.config.route.v3.InternalRedirectPolicy.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getMaxInternalRedirects();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      google_protobuf_wrappers_pb.UInt32Value.serializeBinaryToWriter
    );
  }
  f = message.getRedirectResponseCodesList();
  if (f.length > 0) {
    writer.writePackedUint32(
      2,
      f
    );
  }
  f = message.getPredicatesList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      3,
      f,
      envoy_config_core_v3_extension_pb.TypedExtensionConfig.serializeBinaryToWriter
    );
  }
  f = message.getAllowCrossSchemeRedirect();
  if (f) {
    writer.writeBool(
      4,
      f
    );
  }
};


/**
 * optional google.protobuf.UInt32Value max_internal_redirects = 1;
 * @return {?proto.google.protobuf.UInt32Value}
 */
proto.envoy.config.route.v3.InternalRedirectPolicy.prototype.getMaxInternalRedirects = function() {
  return /** @type{?proto.google.protobuf.UInt32Value} */ (
    jspb.Message.getWrapperField(this, google_protobuf_wrappers_pb.UInt32Value, 1));
};


/** @param {?proto.google.protobuf.UInt32Value|undefined} value */
proto.envoy.config.route.v3.InternalRedirectPolicy.prototype.setMaxInternalRedirects = function(value) {
  jspb.Message.setWrapperField(this, 1, value);
};


proto.envoy.config.route.v3.InternalRedirectPolicy.prototype.clearMaxInternalRedirects = function() {
  this.setMaxInternalRedirects(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.envoy.config.route.v3.InternalRedirectPolicy.prototype.hasMaxInternalRedirects = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * repeated uint32 redirect_response_codes = 2;
 * @return {!Array<number>}
 */
proto.envoy.config.route.v3.InternalRedirectPolicy.prototype.getRedirectResponseCodesList = function() {
  return /** @type {!Array<number>} */ (jspb.Message.getRepeatedField(this, 2));
};


/** @param {!Array<number>} value */
proto.envoy.config.route.v3.InternalRedirectPolicy.prototype.setRedirectResponseCodesList = function(value) {
  jspb.Message.setField(this, 2, value || []);
};


/**
 * @param {!number} value
 * @param {number=} opt_index
 */
proto.envoy.config.route.v3.InternalRedirectPolicy.prototype.addRedirectResponseCodes = function(value, opt_index) {
  jspb.Message.addToRepeatedField(this, 2, value, opt_index);
};


proto.envoy.config.route.v3.InternalRedirectPolicy.prototype.clearRedirectResponseCodesList = function() {
  this.setRedirectResponseCodesList([]);
};


/**
 * repeated envoy.config.core.v3.TypedExtensionConfig predicates = 3;
 * @return {!Array<!proto.envoy.config.core.v3.TypedExtensionConfig>}
 */
proto.envoy.config.route.v3.InternalRedirectPolicy.prototype.getPredicatesList = function() {
  return /** @type{!Array<!proto.envoy.config.core.v3.TypedExtensionConfig>} */ (
    jspb.Message.getRepeatedWrapperField(this, envoy_config_core_v3_extension_pb.TypedExtensionConfig, 3));
};


/** @param {!Array<!proto.envoy.config.core.v3.TypedExtensionConfig>} value */
proto.envoy.config.route.v3.InternalRedirectPolicy.prototype.setPredicatesList = function(value) {
  jspb.Message.setRepeatedWrapperField(this, 3, value);
};


/**
 * @param {!proto.envoy.config.core.v3.TypedExtensionConfig=} opt_value
 * @param {number=} opt_index
 * @return {!proto.envoy.config.core.v3.TypedExtensionConfig}
 */
proto.envoy.config.route.v3.InternalRedirectPolicy.prototype.addPredicates = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 3, opt_value, proto.envoy.config.core.v3.TypedExtensionConfig, opt_index);
};


proto.envoy.config.route.v3.InternalRedirectPolicy.prototype.clearPredicatesList = function() {
  this.setPredicatesList([]);
};


/**
 * optional bool allow_cross_scheme_redirect = 4;
 * Note that Boolean fields may be set to 0/1 when serialized from a Java server.
 * You should avoid comparisons like {@code val === true/false} in those cases.
 * @return {boolean}
 */
proto.envoy.config.route.v3.InternalRedirectPolicy.prototype.getAllowCrossSchemeRedirect = function() {
  return /** @type {boolean} */ (jspb.Message.getFieldWithDefault(this, 4, false));
};


/** @param {boolean} value */
proto.envoy.config.route.v3.InternalRedirectPolicy.prototype.setAllowCrossSchemeRedirect = function(value) {
  jspb.Message.setProto3BooleanField(this, 4, value);
};


goog.object.extend(exports, proto.envoy.config.route.v3);
