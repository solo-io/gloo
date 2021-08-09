// package: rpc.edge.gloo.solo.io
// file: github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/wasm.proto

var github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_wasm_pb = require("../../../../../../../../github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/wasm_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var WasmFilterApi = (function () {
  function WasmFilterApi() {}
  WasmFilterApi.serviceName = "rpc.edge.gloo.solo.io.WasmFilterApi";
  return WasmFilterApi;
}());

WasmFilterApi.ListWasmFilters = {
  methodName: "ListWasmFilters",
  service: WasmFilterApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_wasm_pb.ListWasmFiltersRequest,
  responseType: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_wasm_pb.ListWasmFiltersResponse
};

WasmFilterApi.DescribeWasmFilter = {
  methodName: "DescribeWasmFilter",
  service: WasmFilterApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_wasm_pb.DescribeWasmFilterRequest,
  responseType: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_wasm_pb.DescribeWasmFilterResponse
};

exports.WasmFilterApi = WasmFilterApi;

function WasmFilterApiClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

WasmFilterApiClient.prototype.listWasmFilters = function listWasmFilters(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(WasmFilterApi.ListWasmFilters, {
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

WasmFilterApiClient.prototype.describeWasmFilter = function describeWasmFilter(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(WasmFilterApi.DescribeWasmFilter, {
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

exports.WasmFilterApiClient = WasmFilterApiClient;

