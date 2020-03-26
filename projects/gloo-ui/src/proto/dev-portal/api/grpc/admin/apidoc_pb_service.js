// package: admin.devportal.solo.io
// file: dev-portal/api/grpc/admin/apidoc.proto

var dev_portal_api_grpc_admin_apidoc_pb = require("../../../../dev-portal/api/grpc/admin/apidoc_pb");
var google_protobuf_empty_pb = require("google-protobuf/google/protobuf/empty_pb");
var dev_portal_api_dev_portal_v1_common_pb = require("../../../../dev-portal/api/dev-portal/v1/common_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var ApiDocApi = (function () {
  function ApiDocApi() {}
  ApiDocApi.serviceName = "admin.devportal.solo.io.ApiDocApi";
  return ApiDocApi;
}());

ApiDocApi.GetApiDoc = {
  methodName: "GetApiDoc",
  service: ApiDocApi,
  requestStream: false,
  responseStream: false,
  requestType: dev_portal_api_grpc_admin_apidoc_pb.ApiDocGetRequest,
  responseType: dev_portal_api_grpc_admin_apidoc_pb.ApiDoc
};

ApiDocApi.ListApiDocs = {
  methodName: "ListApiDocs",
  service: ApiDocApi,
  requestStream: false,
  responseStream: false,
  requestType: dev_portal_api_grpc_admin_apidoc_pb.ApiDocFilter,
  responseType: dev_portal_api_grpc_admin_apidoc_pb.ApiDocList
};

ApiDocApi.CreateApiDoc = {
  methodName: "CreateApiDoc",
  service: ApiDocApi,
  requestStream: false,
  responseStream: false,
  requestType: dev_portal_api_grpc_admin_apidoc_pb.ApiDocWriteRequest,
  responseType: dev_portal_api_grpc_admin_apidoc_pb.ApiDoc
};

ApiDocApi.UpdateApiDoc = {
  methodName: "UpdateApiDoc",
  service: ApiDocApi,
  requestStream: false,
  responseStream: false,
  requestType: dev_portal_api_grpc_admin_apidoc_pb.ApiDocWriteRequest,
  responseType: dev_portal_api_grpc_admin_apidoc_pb.ApiDoc
};

ApiDocApi.DeleteApiDoc = {
  methodName: "DeleteApiDoc",
  service: ApiDocApi,
  requestStream: false,
  responseStream: false,
  requestType: dev_portal_api_dev_portal_v1_common_pb.ObjectRef,
  responseType: google_protobuf_empty_pb.Empty
};

exports.ApiDocApi = ApiDocApi;

function ApiDocApiClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

ApiDocApiClient.prototype.getApiDoc = function getApiDoc(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(ApiDocApi.GetApiDoc, {
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

ApiDocApiClient.prototype.listApiDocs = function listApiDocs(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(ApiDocApi.ListApiDocs, {
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

ApiDocApiClient.prototype.createApiDoc = function createApiDoc(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(ApiDocApi.CreateApiDoc, {
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

ApiDocApiClient.prototype.updateApiDoc = function updateApiDoc(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(ApiDocApi.UpdateApiDoc, {
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

ApiDocApiClient.prototype.deleteApiDoc = function deleteApiDoc(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(ApiDocApi.DeleteApiDoc, {
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

exports.ApiDocApiClient = ApiDocApiClient;

