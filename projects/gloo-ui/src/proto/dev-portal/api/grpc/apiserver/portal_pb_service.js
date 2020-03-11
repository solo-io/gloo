// package: apiserver.devportal.solo.io
// file: dev-portal/api/grpc/apiserver/portal.proto

var dev_portal_api_grpc_apiserver_portal_pb = require("../../../../dev-portal/api/grpc/apiserver/portal_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var PortalApi = (function () {
  function PortalApi() {}
  PortalApi.serviceName = "apiserver.devportal.solo.io.PortalApi";
  return PortalApi;
}());

PortalApi.GetPortal = {
  methodName: "GetPortal",
  service: PortalApi,
  requestStream: false,
  responseStream: false,
  requestType: dev_portal_api_grpc_apiserver_portal_pb.GetPortalRequest,
  responseType: dev_portal_api_grpc_apiserver_portal_pb.GetPortalResponse
};

PortalApi.ListPortals = {
  methodName: "ListPortals",
  service: PortalApi,
  requestStream: false,
  responseStream: false,
  requestType: dev_portal_api_grpc_apiserver_portal_pb.ListPortalsRequest,
  responseType: dev_portal_api_grpc_apiserver_portal_pb.ListPortalsResponse
};

PortalApi.CreatePortal = {
  methodName: "CreatePortal",
  service: PortalApi,
  requestStream: false,
  responseStream: false,
  requestType: dev_portal_api_grpc_apiserver_portal_pb.CreatePortalRequest,
  responseType: dev_portal_api_grpc_apiserver_portal_pb.CreatePortalResponse
};

PortalApi.UpdatePortal = {
  methodName: "UpdatePortal",
  service: PortalApi,
  requestStream: false,
  responseStream: false,
  requestType: dev_portal_api_grpc_apiserver_portal_pb.UpdatePortalRequest,
  responseType: dev_portal_api_grpc_apiserver_portal_pb.UpdatePortalResponse
};

PortalApi.DeletePortal = {
  methodName: "DeletePortal",
  service: PortalApi,
  requestStream: false,
  responseStream: false,
  requestType: dev_portal_api_grpc_apiserver_portal_pb.CreatePortalRequest,
  responseType: dev_portal_api_grpc_apiserver_portal_pb.DeletePortalResponse
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

