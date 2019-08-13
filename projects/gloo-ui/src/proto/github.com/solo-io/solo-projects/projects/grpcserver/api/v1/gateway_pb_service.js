// package: glooeeapi.solo.io
// file: github.com/solo-io/solo-projects/projects/grpcserver/api/v1/gateway.proto

var github_com_solo_io_solo_projects_projects_grpcserver_api_v1_gateway_pb = require("../../../../../../../github.com/solo-io/solo-projects/projects/grpcserver/api/v1/gateway_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var GatewayApi = (function () {
  function GatewayApi() {}
  GatewayApi.serviceName = "glooeeapi.solo.io.GatewayApi";
  return GatewayApi;
}());

GatewayApi.GetGateway = {
  methodName: "GetGateway",
  service: GatewayApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_gateway_pb.GetGatewayRequest,
  responseType: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_gateway_pb.GetGatewayResponse
};

GatewayApi.ListGateways = {
  methodName: "ListGateways",
  service: GatewayApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_gateway_pb.ListGatewaysRequest,
  responseType: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_gateway_pb.ListGatewaysResponse
};

exports.GatewayApi = GatewayApi;

function GatewayApiClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

GatewayApiClient.prototype.getGateway = function getGateway(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(GatewayApi.GetGateway, {
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

GatewayApiClient.prototype.listGateways = function listGateways(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(GatewayApi.ListGateways, {
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

exports.GatewayApiClient = GatewayApiClient;

