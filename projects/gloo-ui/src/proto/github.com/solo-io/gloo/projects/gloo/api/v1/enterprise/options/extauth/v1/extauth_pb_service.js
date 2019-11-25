// package: enterprise.gloo.solo.io
// file: github.com/solo-io/gloo/projects/gloo/api/v1/enterprise/options/extauth/v1/extauth.proto

var github_com_solo_io_gloo_projects_gloo_api_v1_enterprise_options_extauth_v1_extauth_pb = require("../../../../../../../../../../../github.com/solo-io/gloo/projects/gloo/api/v1/enterprise/options/extauth/v1/extauth_pb");
var envoy_api_v2_discovery_pb = require("../../../../../../../../../../../envoy/api/v2/discovery_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var ExtAuthDiscoveryService = (function () {
  function ExtAuthDiscoveryService() {}
  ExtAuthDiscoveryService.serviceName = "enterprise.gloo.solo.io.ExtAuthDiscoveryService";
  return ExtAuthDiscoveryService;
}());

ExtAuthDiscoveryService.StreamExtAuthConfig = {
  methodName: "StreamExtAuthConfig",
  service: ExtAuthDiscoveryService,
  requestStream: true,
  responseStream: true,
  requestType: envoy_api_v2_discovery_pb.DiscoveryRequest,
  responseType: envoy_api_v2_discovery_pb.DiscoveryResponse
};

ExtAuthDiscoveryService.DeltaExtAuthConfig = {
  methodName: "DeltaExtAuthConfig",
  service: ExtAuthDiscoveryService,
  requestStream: true,
  responseStream: true,
  requestType: envoy_api_v2_discovery_pb.DeltaDiscoveryRequest,
  responseType: envoy_api_v2_discovery_pb.DeltaDiscoveryResponse
};

ExtAuthDiscoveryService.FetchExtAuthConfig = {
  methodName: "FetchExtAuthConfig",
  service: ExtAuthDiscoveryService,
  requestStream: false,
  responseStream: false,
  requestType: envoy_api_v2_discovery_pb.DiscoveryRequest,
  responseType: envoy_api_v2_discovery_pb.DiscoveryResponse
};

exports.ExtAuthDiscoveryService = ExtAuthDiscoveryService;

function ExtAuthDiscoveryServiceClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

ExtAuthDiscoveryServiceClient.prototype.streamExtAuthConfig = function streamExtAuthConfig(metadata) {
  var listeners = {
    data: [],
    end: [],
    status: []
  };
  var client = grpc.client(ExtAuthDiscoveryService.StreamExtAuthConfig, {
    host: this.serviceHost,
    metadata: metadata,
    transport: this.options.transport
  });
  client.onEnd(function (status, statusMessage, trailers) {
    listeners.status.forEach(function (handler) {
      handler({ code: status, details: statusMessage, metadata: trailers });
    });
    listeners.end.forEach(function (handler) {
      handler({ code: status, details: statusMessage, metadata: trailers });
    });
    listeners = null;
  });
  client.onMessage(function (message) {
    listeners.data.forEach(function (handler) {
      handler(message);
    })
  });
  client.start(metadata);
  return {
    on: function (type, handler) {
      listeners[type].push(handler);
      return this;
    },
    write: function (requestMessage) {
      client.send(requestMessage);
      return this;
    },
    end: function () {
      client.finishSend();
    },
    cancel: function () {
      listeners = null;
      client.close();
    }
  };
};

ExtAuthDiscoveryServiceClient.prototype.deltaExtAuthConfig = function deltaExtAuthConfig(metadata) {
  var listeners = {
    data: [],
    end: [],
    status: []
  };
  var client = grpc.client(ExtAuthDiscoveryService.DeltaExtAuthConfig, {
    host: this.serviceHost,
    metadata: metadata,
    transport: this.options.transport
  });
  client.onEnd(function (status, statusMessage, trailers) {
    listeners.status.forEach(function (handler) {
      handler({ code: status, details: statusMessage, metadata: trailers });
    });
    listeners.end.forEach(function (handler) {
      handler({ code: status, details: statusMessage, metadata: trailers });
    });
    listeners = null;
  });
  client.onMessage(function (message) {
    listeners.data.forEach(function (handler) {
      handler(message);
    })
  });
  client.start(metadata);
  return {
    on: function (type, handler) {
      listeners[type].push(handler);
      return this;
    },
    write: function (requestMessage) {
      client.send(requestMessage);
      return this;
    },
    end: function () {
      client.finishSend();
    },
    cancel: function () {
      listeners = null;
      client.close();
    }
  };
};

ExtAuthDiscoveryServiceClient.prototype.fetchExtAuthConfig = function fetchExtAuthConfig(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(ExtAuthDiscoveryService.FetchExtAuthConfig, {
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

exports.ExtAuthDiscoveryServiceClient = ExtAuthDiscoveryServiceClient;

