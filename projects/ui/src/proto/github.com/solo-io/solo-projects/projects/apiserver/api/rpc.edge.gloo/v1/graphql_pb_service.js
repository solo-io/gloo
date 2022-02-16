// package: rpc.edge.gloo.solo.io
// file: github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/graphql.proto

var github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb = require("../../../../../../../../github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/graphql_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var GraphqlApi = (function () {
  function GraphqlApi() {}
  GraphqlApi.serviceName = "rpc.edge.gloo.solo.io.GraphqlApi";
  return GraphqlApi;
}());

GraphqlApi.GetGraphqlSchema = {
  methodName: "GetGraphqlSchema",
  service: GraphqlApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.GetGraphqlSchemaRequest,
  responseType: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.GetGraphqlSchemaResponse
};

GraphqlApi.ListGraphqlSchemas = {
  methodName: "ListGraphqlSchemas",
  service: GraphqlApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.ListGraphqlSchemasRequest,
  responseType: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.ListGraphqlSchemasResponse
};

GraphqlApi.GetGraphqlSchemaYaml = {
  methodName: "GetGraphqlSchemaYaml",
  service: GraphqlApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.GetGraphqlSchemaYamlRequest,
  responseType: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.GetGraphqlSchemaYamlResponse
};

GraphqlApi.CreateGraphqlSchema = {
  methodName: "CreateGraphqlSchema",
  service: GraphqlApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.CreateGraphqlSchemaRequest,
  responseType: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.CreateGraphqlSchemaResponse
};

GraphqlApi.UpdateGraphqlSchema = {
  methodName: "UpdateGraphqlSchema",
  service: GraphqlApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.UpdateGraphqlSchemaRequest,
  responseType: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.UpdateGraphqlSchemaResponse
};

GraphqlApi.DeleteGraphqlSchema = {
  methodName: "DeleteGraphqlSchema",
  service: GraphqlApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.DeleteGraphqlSchemaRequest,
  responseType: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.DeleteGraphqlSchemaResponse
};

GraphqlApi.ValidateResolverYaml = {
  methodName: "ValidateResolverYaml",
  service: GraphqlApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.ValidateResolverYamlRequest,
  responseType: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.ValidateResolverYamlResponse
};

exports.GraphqlApi = GraphqlApi;

function GraphqlApiClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

GraphqlApiClient.prototype.getGraphqlSchema = function getGraphqlSchema(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(GraphqlApi.GetGraphqlSchema, {
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

GraphqlApiClient.prototype.listGraphqlSchemas = function listGraphqlSchemas(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(GraphqlApi.ListGraphqlSchemas, {
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

GraphqlApiClient.prototype.getGraphqlSchemaYaml = function getGraphqlSchemaYaml(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(GraphqlApi.GetGraphqlSchemaYaml, {
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

GraphqlApiClient.prototype.createGraphqlSchema = function createGraphqlSchema(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(GraphqlApi.CreateGraphqlSchema, {
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

GraphqlApiClient.prototype.updateGraphqlSchema = function updateGraphqlSchema(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(GraphqlApi.UpdateGraphqlSchema, {
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

GraphqlApiClient.prototype.deleteGraphqlSchema = function deleteGraphqlSchema(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(GraphqlApi.DeleteGraphqlSchema, {
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

GraphqlApiClient.prototype.validateResolverYaml = function validateResolverYaml(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(GraphqlApi.ValidateResolverYaml, {
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

exports.GraphqlApiClient = GraphqlApiClient;

