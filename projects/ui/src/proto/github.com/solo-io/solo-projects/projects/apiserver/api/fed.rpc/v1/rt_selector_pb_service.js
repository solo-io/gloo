// package: fed.rpc.solo.io
// file: github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/rt_selector.proto

var github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_rt_selector_pb = require("../../../../../../../../github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/rt_selector_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var VirtualServiceRoutesApi = (function () {
  function VirtualServiceRoutesApi() {}
  VirtualServiceRoutesApi.serviceName = "fed.rpc.solo.io.VirtualServiceRoutesApi";
  return VirtualServiceRoutesApi;
}());

VirtualServiceRoutesApi.GetVirtualServiceRoutes = {
  methodName: "GetVirtualServiceRoutes",
  service: VirtualServiceRoutesApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_rt_selector_pb.GetVirtualServiceRoutesRequest,
  responseType: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_rt_selector_pb.GetVirtualServiceRoutesResponse
};

exports.VirtualServiceRoutesApi = VirtualServiceRoutesApi;

function VirtualServiceRoutesApiClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

VirtualServiceRoutesApiClient.prototype.getVirtualServiceRoutes = function getVirtualServiceRoutes(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(VirtualServiceRoutesApi.GetVirtualServiceRoutes, {
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

exports.VirtualServiceRoutesApiClient = VirtualServiceRoutesApiClient;

