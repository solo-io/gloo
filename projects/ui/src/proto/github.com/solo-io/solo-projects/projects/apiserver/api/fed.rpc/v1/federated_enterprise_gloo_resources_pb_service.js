// package: fed.rpc.solo.io
// file: github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/federated_enterprise_gloo_resources.proto

var github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_enterprise_gloo_resources_pb = require("../../../../../../../../github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/federated_enterprise_gloo_resources_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var FederatedEnterpriseGlooResourceApi = (function () {
  function FederatedEnterpriseGlooResourceApi() {}
  FederatedEnterpriseGlooResourceApi.serviceName = "fed.rpc.solo.io.FederatedEnterpriseGlooResourceApi";
  return FederatedEnterpriseGlooResourceApi;
}());

FederatedEnterpriseGlooResourceApi.ListFederatedAuthConfigs = {
  methodName: "ListFederatedAuthConfigs",
  service: FederatedEnterpriseGlooResourceApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_enterprise_gloo_resources_pb.ListFederatedAuthConfigsRequest,
  responseType: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_enterprise_gloo_resources_pb.ListFederatedAuthConfigsResponse
};

FederatedEnterpriseGlooResourceApi.GetFederatedAuthConfigYaml = {
  methodName: "GetFederatedAuthConfigYaml",
  service: FederatedEnterpriseGlooResourceApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_enterprise_gloo_resources_pb.GetFederatedAuthConfigYamlRequest,
  responseType: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_enterprise_gloo_resources_pb.GetFederatedAuthConfigYamlResponse
};

exports.FederatedEnterpriseGlooResourceApi = FederatedEnterpriseGlooResourceApi;

function FederatedEnterpriseGlooResourceApiClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

FederatedEnterpriseGlooResourceApiClient.prototype.listFederatedAuthConfigs = function listFederatedAuthConfigs(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(FederatedEnterpriseGlooResourceApi.ListFederatedAuthConfigs, {
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

FederatedEnterpriseGlooResourceApiClient.prototype.getFederatedAuthConfigYaml = function getFederatedAuthConfigYaml(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(FederatedEnterpriseGlooResourceApi.GetFederatedAuthConfigYaml, {
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

exports.FederatedEnterpriseGlooResourceApiClient = FederatedEnterpriseGlooResourceApiClient;

