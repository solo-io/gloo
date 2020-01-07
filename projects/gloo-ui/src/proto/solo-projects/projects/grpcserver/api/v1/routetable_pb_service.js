// package: glooeeapi.solo.io
// file: solo-projects/projects/grpcserver/api/v1/routetable.proto

var solo_projects_projects_grpcserver_api_v1_routetable_pb = require("../../../../../solo-projects/projects/grpcserver/api/v1/routetable_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var RouteTableApi = (function () {
  function RouteTableApi() {}
  RouteTableApi.serviceName = "glooeeapi.solo.io.RouteTableApi";
  return RouteTableApi;
}());

RouteTableApi.GetRouteTable = {
  methodName: "GetRouteTable",
  service: RouteTableApi,
  requestStream: false,
  responseStream: false,
  requestType: solo_projects_projects_grpcserver_api_v1_routetable_pb.GetRouteTableRequest,
  responseType: solo_projects_projects_grpcserver_api_v1_routetable_pb.GetRouteTableResponse
};

RouteTableApi.ListRouteTables = {
  methodName: "ListRouteTables",
  service: RouteTableApi,
  requestStream: false,
  responseStream: false,
  requestType: solo_projects_projects_grpcserver_api_v1_routetable_pb.ListRouteTablesRequest,
  responseType: solo_projects_projects_grpcserver_api_v1_routetable_pb.ListRouteTablesResponse
};

RouteTableApi.CreateRouteTable = {
  methodName: "CreateRouteTable",
  service: RouteTableApi,
  requestStream: false,
  responseStream: false,
  requestType: solo_projects_projects_grpcserver_api_v1_routetable_pb.CreateRouteTableRequest,
  responseType: solo_projects_projects_grpcserver_api_v1_routetable_pb.CreateRouteTableResponse
};

RouteTableApi.UpdateRouteTable = {
  methodName: "UpdateRouteTable",
  service: RouteTableApi,
  requestStream: false,
  responseStream: false,
  requestType: solo_projects_projects_grpcserver_api_v1_routetable_pb.UpdateRouteTableRequest,
  responseType: solo_projects_projects_grpcserver_api_v1_routetable_pb.UpdateRouteTableResponse
};

RouteTableApi.UpdateRouteTableYaml = {
  methodName: "UpdateRouteTableYaml",
  service: RouteTableApi,
  requestStream: false,
  responseStream: false,
  requestType: solo_projects_projects_grpcserver_api_v1_routetable_pb.UpdateRouteTableYamlRequest,
  responseType: solo_projects_projects_grpcserver_api_v1_routetable_pb.UpdateRouteTableResponse
};

RouteTableApi.DeleteRouteTable = {
  methodName: "DeleteRouteTable",
  service: RouteTableApi,
  requestStream: false,
  responseStream: false,
  requestType: solo_projects_projects_grpcserver_api_v1_routetable_pb.DeleteRouteTableRequest,
  responseType: solo_projects_projects_grpcserver_api_v1_routetable_pb.DeleteRouteTableResponse
};

exports.RouteTableApi = RouteTableApi;

function RouteTableApiClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

RouteTableApiClient.prototype.getRouteTable = function getRouteTable(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(RouteTableApi.GetRouteTable, {
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

RouteTableApiClient.prototype.listRouteTables = function listRouteTables(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(RouteTableApi.ListRouteTables, {
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

RouteTableApiClient.prototype.createRouteTable = function createRouteTable(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(RouteTableApi.CreateRouteTable, {
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

RouteTableApiClient.prototype.updateRouteTable = function updateRouteTable(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(RouteTableApi.UpdateRouteTable, {
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

RouteTableApiClient.prototype.updateRouteTableYaml = function updateRouteTableYaml(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(RouteTableApi.UpdateRouteTableYaml, {
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

RouteTableApiClient.prototype.deleteRouteTable = function deleteRouteTable(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(RouteTableApi.DeleteRouteTable, {
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

exports.RouteTableApiClient = RouteTableApiClient;

