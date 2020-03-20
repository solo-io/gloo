// package: admin.devportal.solo.io
// file: dev-portal/api/grpc/admin/user.proto

var dev_portal_api_grpc_admin_user_pb = require("../../../../dev-portal/api/grpc/admin/user_pb");
var google_protobuf_empty_pb = require("google-protobuf/google/protobuf/empty_pb");
var dev_portal_api_dev_portal_v1_common_pb = require("../../../../dev-portal/api/dev-portal/v1/common_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var UserApi = (function () {
  function UserApi() {}
  UserApi.serviceName = "admin.devportal.solo.io.UserApi";
  return UserApi;
}());

UserApi.GetUser = {
  methodName: "GetUser",
  service: UserApi,
  requestStream: false,
  responseStream: false,
  requestType: dev_portal_api_dev_portal_v1_common_pb.ObjectRef,
  responseType: dev_portal_api_grpc_admin_user_pb.User
};

UserApi.ListUsers = {
  methodName: "ListUsers",
  service: UserApi,
  requestStream: false,
  responseStream: false,
  requestType: dev_portal_api_grpc_admin_user_pb.UserFilter,
  responseType: dev_portal_api_grpc_admin_user_pb.UserList
};

UserApi.CreateUser = {
  methodName: "CreateUser",
  service: UserApi,
  requestStream: false,
  responseStream: false,
  requestType: dev_portal_api_grpc_admin_user_pb.UserWriteRequest,
  responseType: dev_portal_api_grpc_admin_user_pb.User
};

UserApi.UpdateUser = {
  methodName: "UpdateUser",
  service: UserApi,
  requestStream: false,
  responseStream: false,
  requestType: dev_portal_api_grpc_admin_user_pb.UserWriteRequest,
  responseType: dev_portal_api_grpc_admin_user_pb.User
};

UserApi.DeleteUser = {
  methodName: "DeleteUser",
  service: UserApi,
  requestStream: false,
  responseStream: false,
  requestType: dev_portal_api_dev_portal_v1_common_pb.ObjectRef,
  responseType: google_protobuf_empty_pb.Empty
};

exports.UserApi = UserApi;

function UserApiClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

UserApiClient.prototype.getUser = function getUser(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(UserApi.GetUser, {
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

UserApiClient.prototype.listUsers = function listUsers(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(UserApi.ListUsers, {
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

UserApiClient.prototype.createUser = function createUser(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(UserApi.CreateUser, {
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

UserApiClient.prototype.updateUser = function updateUser(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(UserApi.UpdateUser, {
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

UserApiClient.prototype.deleteUser = function deleteUser(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(UserApi.DeleteUser, {
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

exports.UserApiClient = UserApiClient;

