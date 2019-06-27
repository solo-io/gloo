// package: glooeeapi.solo.io
// file: github.com/solo-io/solo-projects/projects/grpcserver/api/v1/secret.proto

var github_com_solo_io_solo_projects_projects_grpcserver_api_v1_secret_pb = require("../../../../../../../github.com/solo-io/solo-projects/projects/grpcserver/api/v1/secret_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var SecretApi = (function () {
  function SecretApi() {}
  SecretApi.serviceName = "glooeeapi.solo.io.SecretApi";
  return SecretApi;
}());

SecretApi.GetSecret = {
  methodName: "GetSecret",
  service: SecretApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_secret_pb.GetSecretRequest,
  responseType: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_secret_pb.GetSecretResponse
};

SecretApi.ListSecrets = {
  methodName: "ListSecrets",
  service: SecretApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_secret_pb.ListSecretsRequest,
  responseType: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_secret_pb.ListSecretsResponse
};

SecretApi.CreateSecret = {
  methodName: "CreateSecret",
  service: SecretApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_secret_pb.CreateSecretRequest,
  responseType: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_secret_pb.CreateSecretResponse
};

SecretApi.UpdateSecret = {
  methodName: "UpdateSecret",
  service: SecretApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_secret_pb.UpdateSecretRequest,
  responseType: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_secret_pb.UpdateSecretResponse
};

SecretApi.DeleteSecret = {
  methodName: "DeleteSecret",
  service: SecretApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_secret_pb.DeleteSecretRequest,
  responseType: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_secret_pb.DeleteSecretResponse
};

exports.SecretApi = SecretApi;

function SecretApiClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

SecretApiClient.prototype.getSecret = function getSecret(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(SecretApi.GetSecret, {
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

SecretApiClient.prototype.listSecrets = function listSecrets(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(SecretApi.ListSecrets, {
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

SecretApiClient.prototype.createSecret = function createSecret(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(SecretApi.CreateSecret, {
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

SecretApiClient.prototype.updateSecret = function updateSecret(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(SecretApi.UpdateSecret, {
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

SecretApiClient.prototype.deleteSecret = function deleteSecret(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(SecretApi.DeleteSecret, {
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

exports.SecretApiClient = SecretApiClient;

