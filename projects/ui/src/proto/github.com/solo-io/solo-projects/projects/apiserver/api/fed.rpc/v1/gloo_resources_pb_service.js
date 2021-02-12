// package: fed.rpc.solo.io
// file: github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/gloo_resources.proto

var github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_gloo_resources_pb = require("../../../../../../../../github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/gloo_resources_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var GlooResourceApi = (function () {
  function GlooResourceApi() {}
  GlooResourceApi.serviceName = "fed.rpc.solo.io.GlooResourceApi";
  return GlooResourceApi;
}());

GlooResourceApi.ListUpstreams = {
  methodName: "ListUpstreams",
  service: GlooResourceApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_gloo_resources_pb.ListUpstreamsRequest,
  responseType: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_gloo_resources_pb.ListUpstreamsResponse
};

GlooResourceApi.GetUpstreamYaml = {
  methodName: "GetUpstreamYaml",
  service: GlooResourceApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_gloo_resources_pb.GetUpstreamYamlRequest,
  responseType: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_gloo_resources_pb.GetUpstreamYamlResponse
};

GlooResourceApi.ListUpstreamGroups = {
  methodName: "ListUpstreamGroups",
  service: GlooResourceApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_gloo_resources_pb.ListUpstreamGroupsRequest,
  responseType: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_gloo_resources_pb.ListUpstreamGroupsResponse
};

GlooResourceApi.GetUpstreamGroupYaml = {
  methodName: "GetUpstreamGroupYaml",
  service: GlooResourceApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_gloo_resources_pb.GetUpstreamGroupYamlRequest,
  responseType: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_gloo_resources_pb.GetUpstreamGroupYamlResponse
};

GlooResourceApi.ListSettings = {
  methodName: "ListSettings",
  service: GlooResourceApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_gloo_resources_pb.ListSettingsRequest,
  responseType: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_gloo_resources_pb.ListSettingsResponse
};

GlooResourceApi.GetSettingsYaml = {
  methodName: "GetSettingsYaml",
  service: GlooResourceApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_gloo_resources_pb.GetSettingsYamlRequest,
  responseType: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_gloo_resources_pb.GetSettingsYamlResponse
};

GlooResourceApi.ListProxies = {
  methodName: "ListProxies",
  service: GlooResourceApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_gloo_resources_pb.ListProxiesRequest,
  responseType: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_gloo_resources_pb.ListProxiesResponse
};

GlooResourceApi.GetProxyYaml = {
  methodName: "GetProxyYaml",
  service: GlooResourceApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_gloo_resources_pb.GetProxyYamlRequest,
  responseType: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_gloo_resources_pb.GetProxyYamlResponse
};

exports.GlooResourceApi = GlooResourceApi;

function GlooResourceApiClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

GlooResourceApiClient.prototype.listUpstreams = function listUpstreams(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(GlooResourceApi.ListUpstreams, {
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

GlooResourceApiClient.prototype.getUpstreamYaml = function getUpstreamYaml(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(GlooResourceApi.GetUpstreamYaml, {
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

GlooResourceApiClient.prototype.listUpstreamGroups = function listUpstreamGroups(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(GlooResourceApi.ListUpstreamGroups, {
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

GlooResourceApiClient.prototype.getUpstreamGroupYaml = function getUpstreamGroupYaml(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(GlooResourceApi.GetUpstreamGroupYaml, {
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

GlooResourceApiClient.prototype.listSettings = function listSettings(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(GlooResourceApi.ListSettings, {
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

GlooResourceApiClient.prototype.getSettingsYaml = function getSettingsYaml(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(GlooResourceApi.GetSettingsYaml, {
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

GlooResourceApiClient.prototype.listProxies = function listProxies(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(GlooResourceApi.ListProxies, {
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

GlooResourceApiClient.prototype.getProxyYaml = function getProxyYaml(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(GlooResourceApi.GetProxyYaml, {
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

exports.GlooResourceApiClient = GlooResourceApiClient;

