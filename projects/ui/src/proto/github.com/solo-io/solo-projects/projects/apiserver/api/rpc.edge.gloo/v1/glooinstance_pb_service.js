// package: rpc.edge.gloo.solo.io
// file: github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/glooinstance.proto

var github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_glooinstance_pb = require("../../../../../../../../github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/glooinstance_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var GlooInstanceApi = (function () {
  function GlooInstanceApi() {}
  GlooInstanceApi.serviceName = "rpc.edge.gloo.solo.io.GlooInstanceApi";
  return GlooInstanceApi;
}());

GlooInstanceApi.ListGlooInstances = {
  methodName: "ListGlooInstances",
  service: GlooInstanceApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_glooinstance_pb.ListGlooInstancesRequest,
  responseType: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_glooinstance_pb.ListGlooInstancesResponse
};

GlooInstanceApi.ListClusterDetails = {
  methodName: "ListClusterDetails",
  service: GlooInstanceApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_glooinstance_pb.ListClusterDetailsRequest,
  responseType: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_glooinstance_pb.ListClusterDetailsResponse
};

GlooInstanceApi.GetConfigDumps = {
  methodName: "GetConfigDumps",
  service: GlooInstanceApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_glooinstance_pb.GetConfigDumpsRequest,
  responseType: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_glooinstance_pb.GetConfigDumpsResponse
};

exports.GlooInstanceApi = GlooInstanceApi;

function GlooInstanceApiClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

GlooInstanceApiClient.prototype.listGlooInstances = function listGlooInstances(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(GlooInstanceApi.ListGlooInstances, {
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

GlooInstanceApiClient.prototype.listClusterDetails = function listClusterDetails(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(GlooInstanceApi.ListClusterDetails, {
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

GlooInstanceApiClient.prototype.getConfigDumps = function getConfigDumps(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(GlooInstanceApi.GetConfigDumps, {
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

exports.GlooInstanceApiClient = GlooInstanceApiClient;

