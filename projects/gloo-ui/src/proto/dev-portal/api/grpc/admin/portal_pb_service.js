// package: admin.devportal.solo.io
// file: dev-portal/api/grpc/admin/portal.proto

var dev_portal_api_grpc_admin_portal_pb = require("../../../../dev-portal/api/grpc/admin/portal_pb");
var google_protobuf_empty_pb = require("google-protobuf/google/protobuf/empty_pb");
var dev_portal_api_dev_portal_v1_common_pb = require("../../../../dev-portal/api/dev-portal/v1/common_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var PortalApi = (function () {
  function PortalApi() {}
  PortalApi.serviceName = "admin.devportal.solo.io.PortalApi";
  return PortalApi;
}());

PortalApi.GetPortal = {
  methodName: "GetPortal",
  service: PortalApi,
  requestStream: false,
  responseStream: false,
  requestType: dev_portal_api_dev_portal_v1_common_pb.ObjectRef,
  responseType: dev_portal_api_grpc_admin_portal_pb.Portal
};

PortalApi.GetPortalWithAssets = {
  methodName: "GetPortalWithAssets",
  service: PortalApi,
  requestStream: false,
  responseStream: false,
  requestType: dev_portal_api_dev_portal_v1_common_pb.ObjectRef,
  responseType: dev_portal_api_grpc_admin_portal_pb.Portal
};

PortalApi.ListPortals = {
  methodName: "ListPortals",
  service: PortalApi,
  requestStream: false,
  responseStream: false,
  requestType: google_protobuf_empty_pb.Empty,
  responseType: dev_portal_api_grpc_admin_portal_pb.PortalList
};

PortalApi.CreatePortal = {
  methodName: "CreatePortal",
  service: PortalApi,
  requestStream: false,
  responseStream: false,
  requestType: dev_portal_api_grpc_admin_portal_pb.PortalWriteRequest,
  responseType: dev_portal_api_grpc_admin_portal_pb.Portal
};

PortalApi.UpdatePortal = {
  methodName: "UpdatePortal",
  service: PortalApi,
  requestStream: false,
  responseStream: false,
  requestType: dev_portal_api_grpc_admin_portal_pb.PortalWriteRequest,
  responseType: dev_portal_api_grpc_admin_portal_pb.Portal
};

PortalApi.DeletePortal = {
  methodName: "DeletePortal",
  service: PortalApi,
  requestStream: false,
  responseStream: false,
  requestType: dev_portal_api_dev_portal_v1_common_pb.ObjectRef,
  responseType: google_protobuf_empty_pb.Empty
};

exports.PortalApi = PortalApi;

function PortalApiClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

PortalApiClient.prototype.getPortal = function getPortal(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(PortalApi.GetPortal, {
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

PortalApiClient.prototype.getPortalWithAssets = function getPortalWithAssets(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(PortalApi.GetPortalWithAssets, {
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

PortalApiClient.prototype.listPortals = function listPortals(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(PortalApi.ListPortals, {
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

PortalApiClient.prototype.createPortal = function createPortal(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(PortalApi.CreatePortal, {
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

PortalApiClient.prototype.updatePortal = function updatePortal(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(PortalApi.UpdatePortal, {
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

PortalApiClient.prototype.deletePortal = function deletePortal(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(PortalApi.DeletePortal, {
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

exports.PortalApiClient = PortalApiClient;

