// package: glooeeapi.solo.io
// file: github.com/solo-io/solo-projects/projects/grpcserver/api/v1/virtualservice.proto

var github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb = require("../../../../../../../github.com/solo-io/solo-projects/projects/grpcserver/api/v1/virtualservice_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var VirtualServiceApi = (function () {
  function VirtualServiceApi() {}
  VirtualServiceApi.serviceName = "glooeeapi.solo.io.VirtualServiceApi";
  return VirtualServiceApi;
}());

VirtualServiceApi.GetVirtualService = {
  methodName: "GetVirtualService",
  service: VirtualServiceApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.GetVirtualServiceRequest,
  responseType: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.GetVirtualServiceResponse
};

VirtualServiceApi.ListVirtualServices = {
  methodName: "ListVirtualServices",
  service: VirtualServiceApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.ListVirtualServicesRequest,
  responseType: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.ListVirtualServicesResponse
};

VirtualServiceApi.StreamVirtualServiceList = {
  methodName: "StreamVirtualServiceList",
  service: VirtualServiceApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.StreamVirtualServiceListRequest,
  responseType: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.StreamVirtualServiceListResponse
};

VirtualServiceApi.CreateVirtualService = {
  methodName: "CreateVirtualService",
  service: VirtualServiceApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.CreateVirtualServiceRequest,
  responseType: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.CreateVirtualServiceResponse
};

VirtualServiceApi.UpdateVirtualService = {
  methodName: "UpdateVirtualService",
  service: VirtualServiceApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.UpdateVirtualServiceRequest,
  responseType: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.UpdateVirtualServiceResponse
};

VirtualServiceApi.DeleteVirtualService = {
  methodName: "DeleteVirtualService",
  service: VirtualServiceApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.DeleteVirtualServiceRequest,
  responseType: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.DeleteVirtualServiceResponse
};

VirtualServiceApi.CreateRoute = {
  methodName: "CreateRoute",
  service: VirtualServiceApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.CreateRouteRequest,
  responseType: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.CreateRouteResponse
};

VirtualServiceApi.UpdateRoute = {
  methodName: "UpdateRoute",
  service: VirtualServiceApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.UpdateRouteRequest,
  responseType: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.UpdateRouteResponse
};

VirtualServiceApi.DeleteRoute = {
  methodName: "DeleteRoute",
  service: VirtualServiceApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.DeleteRouteRequest,
  responseType: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.DeleteRouteResponse
};

VirtualServiceApi.SwapRoutes = {
  methodName: "SwapRoutes",
  service: VirtualServiceApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.SwapRoutesRequest,
  responseType: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.SwapRoutesResponse
};

VirtualServiceApi.ShiftRoutes = {
  methodName: "ShiftRoutes",
  service: VirtualServiceApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.ShiftRoutesRequest,
  responseType: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_virtualservice_pb.ShiftRoutesResponse
};

exports.VirtualServiceApi = VirtualServiceApi;

function VirtualServiceApiClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

VirtualServiceApiClient.prototype.getVirtualService = function getVirtualService(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(VirtualServiceApi.GetVirtualService, {
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

VirtualServiceApiClient.prototype.listVirtualServices = function listVirtualServices(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(VirtualServiceApi.ListVirtualServices, {
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

VirtualServiceApiClient.prototype.streamVirtualServiceList = function streamVirtualServiceList(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(VirtualServiceApi.StreamVirtualServiceList, {
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

VirtualServiceApiClient.prototype.createVirtualService = function createVirtualService(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(VirtualServiceApi.CreateVirtualService, {
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

VirtualServiceApiClient.prototype.updateVirtualService = function updateVirtualService(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(VirtualServiceApi.UpdateVirtualService, {
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

VirtualServiceApiClient.prototype.deleteVirtualService = function deleteVirtualService(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(VirtualServiceApi.DeleteVirtualService, {
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

VirtualServiceApiClient.prototype.createRoute = function createRoute(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(VirtualServiceApi.CreateRoute, {
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

VirtualServiceApiClient.prototype.updateRoute = function updateRoute(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(VirtualServiceApi.UpdateRoute, {
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

VirtualServiceApiClient.prototype.deleteRoute = function deleteRoute(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(VirtualServiceApi.DeleteRoute, {
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

VirtualServiceApiClient.prototype.swapRoutes = function swapRoutes(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(VirtualServiceApi.SwapRoutes, {
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

VirtualServiceApiClient.prototype.shiftRoutes = function shiftRoutes(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(VirtualServiceApi.ShiftRoutes, {
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

exports.VirtualServiceApiClient = VirtualServiceApiClient;

