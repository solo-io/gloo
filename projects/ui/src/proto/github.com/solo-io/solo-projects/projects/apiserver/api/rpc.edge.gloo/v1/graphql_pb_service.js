// package: rpc.edge.gloo.solo.io
// file: github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/graphql.proto

var github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb = require("../../../../../../../../github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/graphql_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var GraphqlConfigApi = (function () {
  function GraphqlConfigApi() {}
  GraphqlConfigApi.serviceName = "rpc.edge.gloo.solo.io.GraphqlConfigApi";
  return GraphqlConfigApi;
}());

GraphqlConfigApi.GetGraphqlApi = {
  methodName: "GetGraphqlApi",
  service: GraphqlConfigApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.GetGraphqlApiRequest,
  responseType: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.GetGraphqlApiResponse
};

GraphqlConfigApi.ListGraphqlApis = {
  methodName: "ListGraphqlApis",
  service: GraphqlConfigApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.ListGraphqlApisRequest,
  responseType: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.ListGraphqlApisResponse
};

GraphqlConfigApi.GetGraphqlApiYaml = {
  methodName: "GetGraphqlApiYaml",
  service: GraphqlConfigApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.GetGraphqlApiYamlRequest,
  responseType: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.GetGraphqlApiYamlResponse
};

GraphqlConfigApi.CreateGraphqlApi = {
  methodName: "CreateGraphqlApi",
  service: GraphqlConfigApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.CreateGraphqlApiRequest,
  responseType: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.CreateGraphqlApiResponse
};

GraphqlConfigApi.UpdateGraphqlApi = {
  methodName: "UpdateGraphqlApi",
  service: GraphqlConfigApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.UpdateGraphqlApiRequest,
  responseType: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.UpdateGraphqlApiResponse
};

GraphqlConfigApi.DeleteGraphqlApi = {
  methodName: "DeleteGraphqlApi",
  service: GraphqlConfigApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.DeleteGraphqlApiRequest,
  responseType: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.DeleteGraphqlApiResponse
};

GraphqlConfigApi.ValidateResolverYaml = {
  methodName: "ValidateResolverYaml",
  service: GraphqlConfigApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.ValidateResolverYamlRequest,
  responseType: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.ValidateResolverYamlResponse
};

GraphqlConfigApi.ValidateSchemaDefinition = {
  methodName: "ValidateSchemaDefinition",
  service: GraphqlConfigApi,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.ValidateSchemaDefinitionRequest,
  responseType: github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_graphql_pb.ValidateSchemaDefinitionResponse
};

exports.GraphqlConfigApi = GraphqlConfigApi;

function GraphqlConfigApiClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

GraphqlConfigApiClient.prototype.getGraphqlApi = function getGraphqlApi(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(GraphqlConfigApi.GetGraphqlApi, {
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

GraphqlConfigApiClient.prototype.listGraphqlApis = function listGraphqlApis(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(GraphqlConfigApi.ListGraphqlApis, {
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

GraphqlConfigApiClient.prototype.getGraphqlApiYaml = function getGraphqlApiYaml(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(GraphqlConfigApi.GetGraphqlApiYaml, {
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

GraphqlConfigApiClient.prototype.createGraphqlApi = function createGraphqlApi(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(GraphqlConfigApi.CreateGraphqlApi, {
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

GraphqlConfigApiClient.prototype.updateGraphqlApi = function updateGraphqlApi(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(GraphqlConfigApi.UpdateGraphqlApi, {
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

GraphqlConfigApiClient.prototype.deleteGraphqlApi = function deleteGraphqlApi(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(GraphqlConfigApi.DeleteGraphqlApi, {
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

GraphqlConfigApiClient.prototype.validateResolverYaml = function validateResolverYaml(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(GraphqlConfigApi.ValidateResolverYaml, {
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

GraphqlConfigApiClient.prototype.validateSchemaDefinition = function validateSchemaDefinition(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(GraphqlConfigApi.ValidateSchemaDefinition, {
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

exports.GraphqlConfigApiClient = GraphqlConfigApiClient;

