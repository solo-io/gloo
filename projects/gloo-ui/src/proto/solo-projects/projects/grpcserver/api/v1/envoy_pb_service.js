// package: glooeeapi.solo.io
// file: solo-projects/projects/grpcserver/api/v1/envoy.proto

var solo_projects_projects_grpcserver_api_v1_envoy_pb = require("../../../../../solo-projects/projects/grpcserver/api/v1/envoy_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var EnvoyApi = (function () {
  function EnvoyApi() {}
  EnvoyApi.serviceName = "glooeeapi.solo.io.EnvoyApi";
  return EnvoyApi;
}());

EnvoyApi.ListEnvoyDetails = {
  methodName: "ListEnvoyDetails",
  service: EnvoyApi,
  requestStream: false,
  responseStream: false,
  requestType: solo_projects_projects_grpcserver_api_v1_envoy_pb.ListEnvoyDetailsRequest,
  responseType: solo_projects_projects_grpcserver_api_v1_envoy_pb.ListEnvoyDetailsResponse
};

exports.EnvoyApi = EnvoyApi;

function EnvoyApiClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

EnvoyApiClient.prototype.listEnvoyDetails = function listEnvoyDetails(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(EnvoyApi.ListEnvoyDetails, {
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

exports.EnvoyApiClient = EnvoyApiClient;

