// package: fed.rpc.solo.io
// file: github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/enterprise_gloo_resources.proto

var github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_enterprise_gloo_resources_pb = require("../../../../../../../../github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/enterprise_gloo_resources_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var EnterpriseGlooResourceApi = (function () {
  function EnterpriseGlooResourceApi() {}
  EnterpriseGlooResourceApi.serviceName = "fed.rpc.solo.io.EnterpriseGlooResourceApi";
  return EnterpriseGlooResourceApi;
}());

EnterpriseGlooResourceApi.ListAuthConfigs = {
  methodName: "ListAuthConfigs",
  service: EnterpriseGlooResourceApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_enterprise_gloo_resources_pb.ListAuthConfigsRequest,
  responseType: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_enterprise_gloo_resources_pb.ListAuthConfigsResponse
};

EnterpriseGlooResourceApi.GetAuthConfigYaml = {
  methodName: "GetAuthConfigYaml",
  service: EnterpriseGlooResourceApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_enterprise_gloo_resources_pb.GetAuthConfigYamlRequest,
  responseType: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_enterprise_gloo_resources_pb.GetAuthConfigYamlResponse
};

exports.EnterpriseGlooResourceApi = EnterpriseGlooResourceApi;

function EnterpriseGlooResourceApiClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

EnterpriseGlooResourceApiClient.prototype.listAuthConfigs = function listAuthConfigs(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(EnterpriseGlooResourceApi.ListAuthConfigs, {
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

EnterpriseGlooResourceApiClient.prototype.getAuthConfigYaml = function getAuthConfigYaml(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(EnterpriseGlooResourceApi.GetAuthConfigYaml, {
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

exports.EnterpriseGlooResourceApiClient = EnterpriseGlooResourceApiClient;

