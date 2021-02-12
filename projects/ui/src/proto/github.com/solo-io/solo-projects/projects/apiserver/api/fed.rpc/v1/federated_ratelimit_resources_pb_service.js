// package: fed.rpc.solo.io
// file: github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/federated_ratelimit_resources.proto

var github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_ratelimit_resources_pb = require("../../../../../../../../github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/federated_ratelimit_resources_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var FederatedRatelimitResourceApi = (function () {
  function FederatedRatelimitResourceApi() {}
  FederatedRatelimitResourceApi.serviceName = "fed.rpc.solo.io.FederatedRatelimitResourceApi";
  return FederatedRatelimitResourceApi;
}());

FederatedRatelimitResourceApi.ListFederatedRateLimitConfigs = {
  methodName: "ListFederatedRateLimitConfigs",
  service: FederatedRatelimitResourceApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_ratelimit_resources_pb.ListFederatedRateLimitConfigsRequest,
  responseType: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_ratelimit_resources_pb.ListFederatedRateLimitConfigsResponse
};

FederatedRatelimitResourceApi.GetFederatedRateLimitConfigYaml = {
  methodName: "GetFederatedRateLimitConfigYaml",
  service: FederatedRatelimitResourceApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_ratelimit_resources_pb.GetFederatedRateLimitConfigYamlRequest,
  responseType: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_ratelimit_resources_pb.GetFederatedRateLimitConfigYamlResponse
};

exports.FederatedRatelimitResourceApi = FederatedRatelimitResourceApi;

function FederatedRatelimitResourceApiClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

FederatedRatelimitResourceApiClient.prototype.listFederatedRateLimitConfigs = function listFederatedRateLimitConfigs(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(FederatedRatelimitResourceApi.ListFederatedRateLimitConfigs, {
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

FederatedRatelimitResourceApiClient.prototype.getFederatedRateLimitConfigYaml = function getFederatedRateLimitConfigYaml(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(FederatedRatelimitResourceApi.GetFederatedRateLimitConfigYaml, {
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

exports.FederatedRatelimitResourceApiClient = FederatedRatelimitResourceApiClient;

