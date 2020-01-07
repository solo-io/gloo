// package: glooeeapi.solo.io
// file: solo-projects/projects/grpcserver/api/v1/upstreamgroup.proto

var solo_projects_projects_grpcserver_api_v1_upstreamgroup_pb = require("../../../../../solo-projects/projects/grpcserver/api/v1/upstreamgroup_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var UpstreamGroupApi = (function () {
  function UpstreamGroupApi() {}
  UpstreamGroupApi.serviceName = "glooeeapi.solo.io.UpstreamGroupApi";
  return UpstreamGroupApi;
}());

UpstreamGroupApi.GetUpstreamGroup = {
  methodName: "GetUpstreamGroup",
  service: UpstreamGroupApi,
  requestStream: false,
  responseStream: false,
  requestType: solo_projects_projects_grpcserver_api_v1_upstreamgroup_pb.GetUpstreamGroupRequest,
  responseType: solo_projects_projects_grpcserver_api_v1_upstreamgroup_pb.GetUpstreamGroupResponse
};

UpstreamGroupApi.ListUpstreamGroups = {
  methodName: "ListUpstreamGroups",
  service: UpstreamGroupApi,
  requestStream: false,
  responseStream: false,
  requestType: solo_projects_projects_grpcserver_api_v1_upstreamgroup_pb.ListUpstreamGroupsRequest,
  responseType: solo_projects_projects_grpcserver_api_v1_upstreamgroup_pb.ListUpstreamGroupsResponse
};

UpstreamGroupApi.CreateUpstreamGroup = {
  methodName: "CreateUpstreamGroup",
  service: UpstreamGroupApi,
  requestStream: false,
  responseStream: false,
  requestType: solo_projects_projects_grpcserver_api_v1_upstreamgroup_pb.CreateUpstreamGroupRequest,
  responseType: solo_projects_projects_grpcserver_api_v1_upstreamgroup_pb.CreateUpstreamGroupResponse
};

UpstreamGroupApi.UpdateUpstreamGroup = {
  methodName: "UpdateUpstreamGroup",
  service: UpstreamGroupApi,
  requestStream: false,
  responseStream: false,
  requestType: solo_projects_projects_grpcserver_api_v1_upstreamgroup_pb.UpdateUpstreamGroupRequest,
  responseType: solo_projects_projects_grpcserver_api_v1_upstreamgroup_pb.UpdateUpstreamGroupResponse
};

UpstreamGroupApi.DeleteUpstreamGroup = {
  methodName: "DeleteUpstreamGroup",
  service: UpstreamGroupApi,
  requestStream: false,
  responseStream: false,
  requestType: solo_projects_projects_grpcserver_api_v1_upstreamgroup_pb.DeleteUpstreamGroupRequest,
  responseType: solo_projects_projects_grpcserver_api_v1_upstreamgroup_pb.DeleteUpstreamGroupResponse
};

exports.UpstreamGroupApi = UpstreamGroupApi;

function UpstreamGroupApiClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

UpstreamGroupApiClient.prototype.getUpstreamGroup = function getUpstreamGroup(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(UpstreamGroupApi.GetUpstreamGroup, {
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

UpstreamGroupApiClient.prototype.listUpstreamGroups = function listUpstreamGroups(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(UpstreamGroupApi.ListUpstreamGroups, {
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

UpstreamGroupApiClient.prototype.createUpstreamGroup = function createUpstreamGroup(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(UpstreamGroupApi.CreateUpstreamGroup, {
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

UpstreamGroupApiClient.prototype.updateUpstreamGroup = function updateUpstreamGroup(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(UpstreamGroupApi.UpdateUpstreamGroup, {
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

UpstreamGroupApiClient.prototype.deleteUpstreamGroup = function deleteUpstreamGroup(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(UpstreamGroupApi.DeleteUpstreamGroup, {
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

exports.UpstreamGroupApiClient = UpstreamGroupApiClient;

