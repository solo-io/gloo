// package: glooeeapi.solo.io
// file: github.com/solo-io/solo-projects/projects/grpcserver/api/v1/proxy.proto

var github_com_solo_io_solo_projects_projects_grpcserver_api_v1_proxy_pb = require("../../../../../../../github.com/solo-io/solo-projects/projects/grpcserver/api/v1/proxy_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var ProxyApi = (function () {
  function ProxyApi() {}
  ProxyApi.serviceName = "glooeeapi.solo.io.ProxyApi";
  return ProxyApi;
}());

ProxyApi.GetProxy = {
  methodName: "GetProxy",
  service: ProxyApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_proxy_pb.GetProxyRequest,
  responseType: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_proxy_pb.GetProxyResponse
};

ProxyApi.ListProxies = {
  methodName: "ListProxies",
  service: ProxyApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_proxy_pb.ListProxiesRequest,
  responseType: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_proxy_pb.ListProxiesResponse
};

exports.ProxyApi = ProxyApi;

function ProxyApiClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

ProxyApiClient.prototype.getProxy = function getProxy(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(ProxyApi.GetProxy, {
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

ProxyApiClient.prototype.listProxies = function listProxies(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(ProxyApi.ListProxies, {
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

exports.ProxyApiClient = ProxyApiClient;

