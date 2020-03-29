// package: admin.devportal.solo.io
// file: dev-portal/api/grpc/admin/api_key.proto

var dev_portal_api_grpc_admin_api_key_pb = require("../../../../dev-portal/api/grpc/admin/api_key_pb");
var google_protobuf_empty_pb = require("google-protobuf/google/protobuf/empty_pb");
var dev_portal_api_dev_portal_v1_common_pb = require("../../../../dev-portal/api/dev-portal/v1/common_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var ApiKeyApi = (function () {
  function ApiKeyApi() {}
  ApiKeyApi.serviceName = "admin.devportal.solo.io.ApiKeyApi";
  return ApiKeyApi;
}());

ApiKeyApi.ListApiKeys = {
  methodName: "ListApiKeys",
  service: ApiKeyApi,
  requestStream: false,
  responseStream: false,
  requestType: google_protobuf_empty_pb.Empty,
  responseType: dev_portal_api_grpc_admin_api_key_pb.ApiKeyList
};

ApiKeyApi.DeleteApiKey = {
  methodName: "DeleteApiKey",
  service: ApiKeyApi,
  requestStream: false,
  responseStream: false,
  requestType: dev_portal_api_dev_portal_v1_common_pb.ObjectRef,
  responseType: google_protobuf_empty_pb.Empty
};

exports.ApiKeyApi = ApiKeyApi;

function ApiKeyApiClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

ApiKeyApiClient.prototype.listApiKeys = function listApiKeys(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(ApiKeyApi.ListApiKeys, {
    request: requestMessage,
    host: this.serviceHost,
    metadata: metadata,
    transport: this.options.transport,
    debug: this.options.debug,
    onEnd: function (response) {
      if (callback) {
        if (response.status !== grpc.Code.OK) {
          var err = new Error(response.statusMessage);
          err.code = response.status;
          err.metadata = response.trailers;
          callback(err, null);
        } else {
          callback(null, response.message);
        }
      }
    }
  });
  return {
    cancel: function () {
      callback = null;
      client.close();
    }
  };
};

ApiKeyApiClient.prototype.deleteApiKey = function deleteApiKey(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(ApiKeyApi.DeleteApiKey, {
    request: requestMessage,
    host: this.serviceHost,
    metadata: metadata,
    transport: this.options.transport,
    debug: this.options.debug,
    onEnd: function (response) {
      if (callback) {
        if (response.status !== grpc.Code.OK) {
          var err = new Error(response.statusMessage);
          err.code = response.status;
          err.metadata = response.trailers;
          callback(err, null);
        } else {
          callback(null, response.message);
        }
      }
    }
  });
  return {
    cancel: function () {
      callback = null;
      client.close();
    }
  };
};

exports.ApiKeyApiClient = ApiKeyApiClient;

