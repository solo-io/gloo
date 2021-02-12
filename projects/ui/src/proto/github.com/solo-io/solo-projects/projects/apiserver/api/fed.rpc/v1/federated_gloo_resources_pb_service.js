// package: fed.rpc.solo.io
// file: github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/federated_gloo_resources.proto

var github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gloo_resources_pb = require("../../../../../../../../github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/federated_gloo_resources_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var FederatedGlooResourceApi = (function () {
  function FederatedGlooResourceApi() {}
  FederatedGlooResourceApi.serviceName = "fed.rpc.solo.io.FederatedGlooResourceApi";
  return FederatedGlooResourceApi;
}());

FederatedGlooResourceApi.ListFederatedUpstreams = {
  methodName: "ListFederatedUpstreams",
  service: FederatedGlooResourceApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gloo_resources_pb.ListFederatedUpstreamsRequest,
  responseType: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gloo_resources_pb.ListFederatedUpstreamsResponse
};

FederatedGlooResourceApi.GetFederatedUpstreamYaml = {
  methodName: "GetFederatedUpstreamYaml",
  service: FederatedGlooResourceApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gloo_resources_pb.GetFederatedUpstreamYamlRequest,
  responseType: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gloo_resources_pb.GetFederatedUpstreamYamlResponse
};

FederatedGlooResourceApi.ListFederatedUpstreamGroups = {
  methodName: "ListFederatedUpstreamGroups",
  service: FederatedGlooResourceApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gloo_resources_pb.ListFederatedUpstreamGroupsRequest,
  responseType: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gloo_resources_pb.ListFederatedUpstreamGroupsResponse
};

FederatedGlooResourceApi.GetFederatedUpstreamGroupYaml = {
  methodName: "GetFederatedUpstreamGroupYaml",
  service: FederatedGlooResourceApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gloo_resources_pb.GetFederatedUpstreamGroupYamlRequest,
  responseType: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gloo_resources_pb.GetFederatedUpstreamGroupYamlResponse
};

FederatedGlooResourceApi.ListFederatedSettings = {
  methodName: "ListFederatedSettings",
  service: FederatedGlooResourceApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gloo_resources_pb.ListFederatedSettingsRequest,
  responseType: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gloo_resources_pb.ListFederatedSettingsResponse
};

FederatedGlooResourceApi.GetFederatedSettingsYaml = {
  methodName: "GetFederatedSettingsYaml",
  service: FederatedGlooResourceApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gloo_resources_pb.GetFederatedSettingsYamlRequest,
  responseType: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gloo_resources_pb.GetFederatedSettingsYamlResponse
};

exports.FederatedGlooResourceApi = FederatedGlooResourceApi;

function FederatedGlooResourceApiClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

FederatedGlooResourceApiClient.prototype.listFederatedUpstreams = function listFederatedUpstreams(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(FederatedGlooResourceApi.ListFederatedUpstreams, {
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

FederatedGlooResourceApiClient.prototype.getFederatedUpstreamYaml = function getFederatedUpstreamYaml(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(FederatedGlooResourceApi.GetFederatedUpstreamYaml, {
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

FederatedGlooResourceApiClient.prototype.listFederatedUpstreamGroups = function listFederatedUpstreamGroups(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(FederatedGlooResourceApi.ListFederatedUpstreamGroups, {
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

FederatedGlooResourceApiClient.prototype.getFederatedUpstreamGroupYaml = function getFederatedUpstreamGroupYaml(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(FederatedGlooResourceApi.GetFederatedUpstreamGroupYaml, {
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

FederatedGlooResourceApiClient.prototype.listFederatedSettings = function listFederatedSettings(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(FederatedGlooResourceApi.ListFederatedSettings, {
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

FederatedGlooResourceApiClient.prototype.getFederatedSettingsYaml = function getFederatedSettingsYaml(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(FederatedGlooResourceApi.GetFederatedSettingsYaml, {
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

exports.FederatedGlooResourceApiClient = FederatedGlooResourceApiClient;

