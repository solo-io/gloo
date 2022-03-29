// package: rpc.edge.gloo.solo.io
// file: github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/bootstrap.proto

var github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_bootstrap_pb = require("../../../../../../../../github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/bootstrap_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var BootstrapApi = (function () {
  function BootstrapApi() {}
  BootstrapApi.serviceName = "rpc.edge.gloo.solo.io.BootstrapApi";
  return BootstrapApi;
}());

BootstrapApi.IsGlooFedEnabled = {
  methodName: "IsGlooFedEnabled",
  service: BootstrapApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_bootstrap_pb.GlooFedCheckRequest,
  responseType: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_bootstrap_pb.GlooFedCheckResponse
};

BootstrapApi.IsGraphqlEnabled = {
  methodName: "IsGraphqlEnabled",
  service: BootstrapApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_bootstrap_pb.GraphqlCheckRequest,
  responseType: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_bootstrap_pb.GraphqlCheckResponse
};

BootstrapApi.GetConsoleOptions = {
  methodName: "GetConsoleOptions",
  service: BootstrapApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_bootstrap_pb.GetConsoleOptionsRequest,
  responseType: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_bootstrap_pb.GetConsoleOptionsResponse
};

exports.BootstrapApi = BootstrapApi;

function BootstrapApiClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

BootstrapApiClient.prototype.isGlooFedEnabled = function isGlooFedEnabled(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(BootstrapApi.IsGlooFedEnabled, {
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

BootstrapApiClient.prototype.isGraphqlEnabled = function isGraphqlEnabled(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(BootstrapApi.IsGraphqlEnabled, {
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

BootstrapApiClient.prototype.getConsoleOptions = function getConsoleOptions(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(BootstrapApi.GetConsoleOptions, {
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

exports.BootstrapApiClient = BootstrapApiClient;

