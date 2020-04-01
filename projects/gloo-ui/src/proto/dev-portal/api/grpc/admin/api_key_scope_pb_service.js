// package: admin.devportal.solo.io
// file: dev-portal/api/grpc/admin/api_key_scope.proto

var dev_portal_api_grpc_admin_api_key_scope_pb = require("../../../../dev-portal/api/grpc/admin/api_key_scope_pb");
var google_protobuf_empty_pb = require("google-protobuf/google/protobuf/empty_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var ApiKeyScopeApi = (function () {
  function ApiKeyScopeApi() {}
  ApiKeyScopeApi.serviceName = "admin.devportal.solo.io.ApiKeyScopeApi";
  return ApiKeyScopeApi;
}());

ApiKeyScopeApi.ListApiKeyScopes = {
  methodName: "ListApiKeyScopes",
  service: ApiKeyScopeApi,
  requestStream: false,
  responseStream: false,
  requestType: dev_portal_api_grpc_admin_api_key_scope_pb.ApiKeyScopeFilter,
  responseType: dev_portal_api_grpc_admin_api_key_scope_pb.ApiKeyScopeList
};

ApiKeyScopeApi.CreateApiKeyScope = {
  methodName: "CreateApiKeyScope",
  service: ApiKeyScopeApi,
  requestStream: false,
  responseStream: false,
  requestType: dev_portal_api_grpc_admin_api_key_scope_pb.ApiKeyScopeWriteRequest,
  responseType: dev_portal_api_grpc_admin_api_key_scope_pb.ApiKeyScope
};

ApiKeyScopeApi.UpdateApiKeyScope = {
  methodName: "UpdateApiKeyScope",
  service: ApiKeyScopeApi,
  requestStream: false,
  responseStream: false,
  requestType: dev_portal_api_grpc_admin_api_key_scope_pb.ApiKeyScopeWriteRequest,
  responseType: dev_portal_api_grpc_admin_api_key_scope_pb.ApiKeyScope
};

ApiKeyScopeApi.DeleteApiKeyScope = {
  methodName: "DeleteApiKeyScope",
  service: ApiKeyScopeApi,
  requestStream: false,
  responseStream: false,
  requestType: dev_portal_api_grpc_admin_api_key_scope_pb.ApiKeyScopeRef,
  responseType: google_protobuf_empty_pb.Empty
};

exports.ApiKeyScopeApi = ApiKeyScopeApi;

function ApiKeyScopeApiClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

ApiKeyScopeApiClient.prototype.listApiKeyScopes = function listApiKeyScopes(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(ApiKeyScopeApi.ListApiKeyScopes, {
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

ApiKeyScopeApiClient.prototype.createApiKeyScope = function createApiKeyScope(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(ApiKeyScopeApi.CreateApiKeyScope, {
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

ApiKeyScopeApiClient.prototype.updateApiKeyScope = function updateApiKeyScope(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(ApiKeyScopeApi.UpdateApiKeyScope, {
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

ApiKeyScopeApiClient.prototype.deleteApiKeyScope = function deleteApiKeyScope(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(ApiKeyScopeApi.DeleteApiKeyScope, {
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

exports.ApiKeyScopeApiClient = ApiKeyScopeApiClient;

