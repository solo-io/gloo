// package: fed.rpc.solo.io
// file: github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/federated_gateway_resources.proto

var github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gateway_resources_pb = require("../../../../../../../../github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/federated_gateway_resources_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var FederatedGatewayResourceApi = (function () {
  function FederatedGatewayResourceApi() {}
  FederatedGatewayResourceApi.serviceName = "fed.rpc.solo.io.FederatedGatewayResourceApi";
  return FederatedGatewayResourceApi;
}());

FederatedGatewayResourceApi.ListFederatedGateways = {
  methodName: "ListFederatedGateways",
  service: FederatedGatewayResourceApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gateway_resources_pb.ListFederatedGatewaysRequest,
  responseType: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gateway_resources_pb.ListFederatedGatewaysResponse
};

FederatedGatewayResourceApi.GetFederatedGatewayYaml = {
  methodName: "GetFederatedGatewayYaml",
  service: FederatedGatewayResourceApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gateway_resources_pb.GetFederatedGatewayYamlRequest,
  responseType: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gateway_resources_pb.GetFederatedGatewayYamlResponse
};

FederatedGatewayResourceApi.ListFederatedVirtualServices = {
  methodName: "ListFederatedVirtualServices",
  service: FederatedGatewayResourceApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gateway_resources_pb.ListFederatedVirtualServicesRequest,
  responseType: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gateway_resources_pb.ListFederatedVirtualServicesResponse
};

FederatedGatewayResourceApi.GetFederatedVirtualServiceYaml = {
  methodName: "GetFederatedVirtualServiceYaml",
  service: FederatedGatewayResourceApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gateway_resources_pb.GetFederatedVirtualServiceYamlRequest,
  responseType: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gateway_resources_pb.GetFederatedVirtualServiceYamlResponse
};

FederatedGatewayResourceApi.ListFederatedRouteTables = {
  methodName: "ListFederatedRouteTables",
  service: FederatedGatewayResourceApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gateway_resources_pb.ListFederatedRouteTablesRequest,
  responseType: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gateway_resources_pb.ListFederatedRouteTablesResponse
};

FederatedGatewayResourceApi.GetFederatedRouteTableYaml = {
  methodName: "GetFederatedRouteTableYaml",
  service: FederatedGatewayResourceApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gateway_resources_pb.GetFederatedRouteTableYamlRequest,
  responseType: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_federated_gateway_resources_pb.GetFederatedRouteTableYamlResponse
};

exports.FederatedGatewayResourceApi = FederatedGatewayResourceApi;

function FederatedGatewayResourceApiClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

FederatedGatewayResourceApiClient.prototype.listFederatedGateways = function listFederatedGateways(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(FederatedGatewayResourceApi.ListFederatedGateways, {
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

FederatedGatewayResourceApiClient.prototype.getFederatedGatewayYaml = function getFederatedGatewayYaml(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(FederatedGatewayResourceApi.GetFederatedGatewayYaml, {
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

FederatedGatewayResourceApiClient.prototype.listFederatedVirtualServices = function listFederatedVirtualServices(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(FederatedGatewayResourceApi.ListFederatedVirtualServices, {
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

FederatedGatewayResourceApiClient.prototype.getFederatedVirtualServiceYaml = function getFederatedVirtualServiceYaml(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(FederatedGatewayResourceApi.GetFederatedVirtualServiceYaml, {
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

FederatedGatewayResourceApiClient.prototype.listFederatedRouteTables = function listFederatedRouteTables(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(FederatedGatewayResourceApi.ListFederatedRouteTables, {
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

FederatedGatewayResourceApiClient.prototype.getFederatedRouteTableYaml = function getFederatedRouteTableYaml(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(FederatedGatewayResourceApi.GetFederatedRouteTableYaml, {
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

exports.FederatedGatewayResourceApiClient = FederatedGatewayResourceApiClient;

