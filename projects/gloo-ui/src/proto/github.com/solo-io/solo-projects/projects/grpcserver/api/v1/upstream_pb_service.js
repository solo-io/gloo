// package: glooeeapi.solo.io
// file: github.com/solo-io/solo-projects/projects/grpcserver/api/v1/upstream.proto

var github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstream_pb = require("../../../../../../../github.com/solo-io/solo-projects/projects/grpcserver/api/v1/upstream_pb");
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
  requestType: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstream_pb.GetUpstreamRequest,
  responseType: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstream_pb.GetUpstreamResponse
};

UpstreamApi.ListUpstreams = {
  methodName: "ListUpstreams",
  service: UpstreamApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstream_pb.ListUpstreamsRequest,
  responseType: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstream_pb.ListUpstreamsResponse
};

UpstreamApi.StreamUpstreamList = {
  methodName: "StreamUpstreamList",
  service: UpstreamApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstream_pb.StreamUpstreamListRequest,
  responseType: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstream_pb.StreamUpstreamListResponse
};

UpstreamApi.CreateUpstream = {
  methodName: "CreateUpstream",
  service: UpstreamApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstream_pb.CreateUpstreamRequest,
  responseType: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstream_pb.CreateUpstreamResponse
};

UpstreamApi.UpdateUpstream = {
  methodName: "UpdateUpstream",
  service: UpstreamApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstream_pb.UpdateUpstreamRequest,
  responseType: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstream_pb.UpdateUpstreamResponse
};

UpstreamApi.DeleteUpstream = {
  methodName: "DeleteUpstream",
  service: UpstreamApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstream_pb.DeleteUpstreamRequest,
  responseType: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_upstream_pb.DeleteUpstreamResponse
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

UpstreamApiClient.prototype.streamUpstreamList = function streamUpstreamList(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(UpstreamApi.StreamUpstreamList, {
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

