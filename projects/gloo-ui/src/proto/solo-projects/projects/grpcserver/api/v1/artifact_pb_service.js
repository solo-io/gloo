// package: glooeeapi.solo.io
// file: solo-projects/projects/grpcserver/api/v1/artifact.proto

var solo_projects_projects_grpcserver_api_v1_artifact_pb = require("../../../../../solo-projects/projects/grpcserver/api/v1/artifact_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var ArtifactApi = (function () {
  function ArtifactApi() {}
  ArtifactApi.serviceName = "glooeeapi.solo.io.ArtifactApi";
  return ArtifactApi;
}());

ArtifactApi.GetArtifact = {
  methodName: "GetArtifact",
  service: ArtifactApi,
  requestStream: false,
  responseStream: false,
  requestType: solo_projects_projects_grpcserver_api_v1_artifact_pb.GetArtifactRequest,
  responseType: solo_projects_projects_grpcserver_api_v1_artifact_pb.GetArtifactResponse
};

ArtifactApi.ListArtifacts = {
  methodName: "ListArtifacts",
  service: ArtifactApi,
  requestStream: false,
  responseStream: false,
  requestType: solo_projects_projects_grpcserver_api_v1_artifact_pb.ListArtifactsRequest,
  responseType: solo_projects_projects_grpcserver_api_v1_artifact_pb.ListArtifactsResponse
};

ArtifactApi.CreateArtifact = {
  methodName: "CreateArtifact",
  service: ArtifactApi,
  requestStream: false,
  responseStream: false,
  requestType: solo_projects_projects_grpcserver_api_v1_artifact_pb.CreateArtifactRequest,
  responseType: solo_projects_projects_grpcserver_api_v1_artifact_pb.CreateArtifactResponse
};

ArtifactApi.UpdateArtifact = {
  methodName: "UpdateArtifact",
  service: ArtifactApi,
  requestStream: false,
  responseStream: false,
  requestType: solo_projects_projects_grpcserver_api_v1_artifact_pb.UpdateArtifactRequest,
  responseType: solo_projects_projects_grpcserver_api_v1_artifact_pb.UpdateArtifactResponse
};

ArtifactApi.DeleteArtifact = {
  methodName: "DeleteArtifact",
  service: ArtifactApi,
  requestStream: false,
  responseStream: false,
  requestType: solo_projects_projects_grpcserver_api_v1_artifact_pb.DeleteArtifactRequest,
  responseType: solo_projects_projects_grpcserver_api_v1_artifact_pb.DeleteArtifactResponse
};

exports.ArtifactApi = ArtifactApi;

function ArtifactApiClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

ArtifactApiClient.prototype.getArtifact = function getArtifact(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(ArtifactApi.GetArtifact, {
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

ArtifactApiClient.prototype.listArtifacts = function listArtifacts(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(ArtifactApi.ListArtifacts, {
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

ArtifactApiClient.prototype.createArtifact = function createArtifact(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(ArtifactApi.CreateArtifact, {
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

ArtifactApiClient.prototype.updateArtifact = function updateArtifact(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(ArtifactApi.UpdateArtifact, {
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

ArtifactApiClient.prototype.deleteArtifact = function deleteArtifact(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(ArtifactApi.DeleteArtifact, {
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

exports.ArtifactApiClient = ArtifactApiClient;

