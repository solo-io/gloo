// package: fed.rpc.solo.io
// file: github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/failover_scheme.proto

var github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_failover_scheme_pb = require("../../../../../../../../github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/failover_scheme_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var FailoverSchemeApi = (function () {
  function FailoverSchemeApi() {}
  FailoverSchemeApi.serviceName = "fed.rpc.solo.io.FailoverSchemeApi";
  return FailoverSchemeApi;
}());

FailoverSchemeApi.GetFailoverScheme = {
  methodName: "GetFailoverScheme",
  service: FailoverSchemeApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_failover_scheme_pb.GetFailoverSchemeRequest,
  responseType: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_failover_scheme_pb.GetFailoverSchemeResponse
};

FailoverSchemeApi.GetFailoverSchemeYaml = {
  methodName: "GetFailoverSchemeYaml",
  service: FailoverSchemeApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_failover_scheme_pb.GetFailoverSchemeYamlRequest,
  responseType: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_failover_scheme_pb.GetFailoverSchemeYamlResponse
};

exports.FailoverSchemeApi = FailoverSchemeApi;

function FailoverSchemeApiClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

FailoverSchemeApiClient.prototype.getFailoverScheme = function getFailoverScheme(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(FailoverSchemeApi.GetFailoverScheme, {
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

FailoverSchemeApiClient.prototype.getFailoverSchemeYaml = function getFailoverSchemeYaml(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(FailoverSchemeApi.GetFailoverSchemeYaml, {
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

exports.FailoverSchemeApiClient = FailoverSchemeApiClient;

