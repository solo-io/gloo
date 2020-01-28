// package: glooeeapi.solo.io
// file: solo-projects/projects/grpcserver/api/v1/upstream.proto

var solo_projects_projects_grpcserver_api_v1_upstream_pb = require("../../../../../solo-projects/projects/grpcserver/api/v1/upstream_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var UpstreamApi = (function () {
  function UpstreamApi() {}
  UpstreamApi.serviceName = "glooeeapi.solo.io.UpstreamApi";
  return UpstreamApi;
}());

UpstreamApi.GetUpstream = {
  methodName: "GetUpstream",
  service: UpstreamApi,
  requestStream: false,
  responseStream: false,
  requestType: solo_projects_projects_grpcserver_api_v1_upstream_pb.GetUpstreamRequest,
  responseType: solo_projects_projects_grpcserver_api_v1_upstream_pb.GetUpstreamResponse
};

UpstreamApi.ListUpstreams = {
  methodName: "ListUpstreams",
  service: UpstreamApi,
  requestStream: false,
  responseStream: false,
  requestType: solo_projects_projects_grpcserver_api_v1_upstream_pb.ListUpstreamsRequest,
  responseType: solo_projects_projects_grpcserver_api_v1_upstream_pb.ListUpstreamsResponse
};

UpstreamApi.CreateUpstream = {
  methodName: "CreateUpstream",
  service: UpstreamApi,
  requestStream: false,
  responseStream: false,
  requestType: solo_projects_projects_grpcserver_api_v1_upstream_pb.CreateUpstreamRequest,
  responseType: solo_projects_projects_grpcserver_api_v1_upstream_pb.CreateUpstreamResponse
};

UpstreamApi.UpdateUpstream = {
  methodName: "UpdateUpstream",
  service: UpstreamApi,
  requestStream: false,
  responseStream: false,
  requestType: solo_projects_projects_grpcserver_api_v1_upstream_pb.UpdateUpstreamRequest,
  responseType: solo_projects_projects_grpcserver_api_v1_upstream_pb.UpdateUpstreamResponse
};

UpstreamApi.UpdateUpstreamYaml = {
  methodName: "UpdateUpstreamYaml",
  service: UpstreamApi,
  requestStream: false,
  responseStream: false,
  requestType: solo_projects_projects_grpcserver_api_v1_upstream_pb.UpdateUpstreamYamlRequest,
  responseType: solo_projects_projects_grpcserver_api_v1_upstream_pb.UpdateUpstreamResponse
};

UpstreamApi.DeleteUpstream = {
  methodName: "DeleteUpstream",
  service: UpstreamApi,
  requestStream: false,
  responseStream: false,
  requestType: solo_projects_projects_grpcserver_api_v1_upstream_pb.DeleteUpstreamRequest,
  responseType: solo_projects_projects_grpcserver_api_v1_upstream_pb.DeleteUpstreamResponse
};

exports.UpstreamApi = UpstreamApi;

function UpstreamApiClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

UpstreamApiClient.prototype.getUpstream = function getUpstream(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(UpstreamApi.GetUpstream, {
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

UpstreamApiClient.prototype.listUpstreams = function listUpstreams(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(UpstreamApi.ListUpstreams, {
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

UpstreamApiClient.prototype.createUpstream = function createUpstream(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(UpstreamApi.CreateUpstream, {
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

UpstreamApiClient.prototype.updateUpstream = function updateUpstream(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(UpstreamApi.UpdateUpstream, {
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

UpstreamApiClient.prototype.updateUpstreamYaml = function updateUpstreamYaml(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(UpstreamApi.UpdateUpstreamYaml, {
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

UpstreamApiClient.prototype.deleteUpstream = function deleteUpstream(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(UpstreamApi.DeleteUpstream, {
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

exports.UpstreamApiClient = UpstreamApiClient;

