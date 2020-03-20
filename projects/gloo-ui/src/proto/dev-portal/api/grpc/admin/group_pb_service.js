// package: admin.devportal.solo.io
// file: dev-portal/api/grpc/admin/group.proto

var dev_portal_api_grpc_admin_group_pb = require("../../../../dev-portal/api/grpc/admin/group_pb");
var google_protobuf_empty_pb = require("google-protobuf/google/protobuf/empty_pb");
var dev_portal_api_dev_portal_v1_common_pb = require("../../../../dev-portal/api/dev-portal/v1/common_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var GroupApi = (function () {
  function GroupApi() {}
  GroupApi.serviceName = "admin.devportal.solo.io.GroupApi";
  return GroupApi;
}());

GroupApi.GetGroup = {
  methodName: "GetGroup",
  service: GroupApi,
  requestStream: false,
  responseStream: false,
  requestType: dev_portal_api_dev_portal_v1_common_pb.ObjectRef,
  responseType: dev_portal_api_grpc_admin_group_pb.Group
};

GroupApi.ListGroups = {
  methodName: "ListGroups",
  service: GroupApi,
  requestStream: false,
  responseStream: false,
  requestType: dev_portal_api_grpc_admin_group_pb.GroupFilter,
  responseType: dev_portal_api_grpc_admin_group_pb.GroupList
};

GroupApi.CreateGroup = {
  methodName: "CreateGroup",
  service: GroupApi,
  requestStream: false,
  responseStream: false,
  requestType: dev_portal_api_grpc_admin_group_pb.GroupWriteRequest,
  responseType: dev_portal_api_grpc_admin_group_pb.Group
};

GroupApi.UpdateGroup = {
  methodName: "UpdateGroup",
  service: GroupApi,
  requestStream: false,
  responseStream: false,
  requestType: dev_portal_api_grpc_admin_group_pb.GroupWriteRequest,
  responseType: dev_portal_api_grpc_admin_group_pb.Group
};

GroupApi.DeleteGroup = {
  methodName: "DeleteGroup",
  service: GroupApi,
  requestStream: false,
  responseStream: false,
  requestType: dev_portal_api_dev_portal_v1_common_pb.ObjectRef,
  responseType: google_protobuf_empty_pb.Empty
};

exports.GroupApi = GroupApi;

function GroupApiClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

GroupApiClient.prototype.getGroup = function getGroup(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(GroupApi.GetGroup, {
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

GroupApiClient.prototype.listGroups = function listGroups(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(GroupApi.ListGroups, {
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

GroupApiClient.prototype.createGroup = function createGroup(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(GroupApi.CreateGroup, {
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

GroupApiClient.prototype.updateGroup = function updateGroup(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(GroupApi.UpdateGroup, {
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

GroupApiClient.prototype.deleteGroup = function deleteGroup(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(GroupApi.DeleteGroup, {
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

exports.GroupApiClient = GroupApiClient;

