import {
  getClusterRefClassFromClusterRefObj,
  getObjectRefClassFromRefObj,
  host,
} from './helpers';
import { grpc } from '@improbable-eng/grpc-web';
import {
  ClusterObjectRef,
  ObjectRef,
} from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';
import { GraphqlApiClient } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/graphql_pb_service';
import {
  GetGraphqlSchemaRequest,
  GraphqlSchema,
  ListGraphqlSchemasRequest,
  GetGraphqlSchemaYamlRequest,
  ValidateResolverYamlRequest,
  ValidateResolverYamlResponse,
  CreateGraphqlSchemaRequest,
  CreateGraphqlSchemaResponse,
  UpdateGraphqlSchemaRequest,
  UpdateGraphqlSchemaResponse,
  DeleteGraphqlSchemaResponse,
  DeleteGraphqlSchemaRequest,
} from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/graphql_pb';
import {
  ExecutableSchema,
  GraphQLSchemaSpec,
  Resolution,
} from 'proto/github.com/solo-io/solo-apis/api/gloo/graphql.gloo/v1alpha1/graphql_pb';

const graphqlApiClient = new GraphqlApiClient(host, {
  transport: grpc.CrossBrowserHttpTransport({ withCredentials: false }),
  debug: true,
});

export const graphqlApi = {
  listGraphqlSchemas,
  getGraphqlSchema,
  getGraphqlSchemaYaml,
  createGraphqlSchema,
  validateResolverYaml,
  updateGraphqlSchema,
  deleteGraphqlSchema,
  updateGraphqlSchemaResolver,
};

function listGraphqlSchemas(
  glooInstanceRef?: ObjectRef.AsObject
): Promise<GraphqlSchema.AsObject[]> {
  let request = new ListGraphqlSchemasRequest();
  if (glooInstanceRef) {
    request.setGlooInstanceRef(getObjectRefClassFromRefObj(glooInstanceRef));
  }
  return new Promise((resolve, reject) => {
    graphqlApiClient.listGraphqlSchemas(request, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        resolve(data!.toObject().graphqlSchemasList);
      }
    });
  });
}

function getGraphqlSchema(
  graphqlSchemaRef: ClusterObjectRef.AsObject
): Promise<GraphqlSchema.AsObject> {
  let request = new GetGraphqlSchemaRequest();
  request.setGraphqlSchemaRef(
    getClusterRefClassFromClusterRefObj(graphqlSchemaRef)
  );

  return new Promise((resolve, reject) => {
    graphqlApiClient.getGraphqlSchema(request, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        resolve(data!.toObject()!.graphqlSchema!);
      }
    });
  });
}

function getGraphqlSchemaPb(
  graphqlSchemaRef: ClusterObjectRef.AsObject
): Promise<GraphqlSchema> {
  let request = new GetGraphqlSchemaRequest();
  request.setGraphqlSchemaRef(
    getClusterRefClassFromClusterRefObj(graphqlSchemaRef)
  );

  return new Promise((resolve, reject) => {
    graphqlApiClient.getGraphqlSchema(request, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        resolve(data!.getGraphqlSchema()!);
      }
    });
  });
}

function getGraphqlSchemaYaml(
  graphqlSchemaRef: ClusterObjectRef.AsObject
): Promise<string> {
  let request = new GetGraphqlSchemaYamlRequest();
  request.setGraphqlSchemaRef(
    getClusterRefClassFromClusterRefObj(graphqlSchemaRef)
  );

  return new Promise((resolve, reject) => {
    graphqlApiClient.getGraphqlSchemaYaml(request, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        resolve(data!.toObject().yamlData?.yaml ?? 'None');
      }
    });
  });
}

function createGraphqlSchema(
  createGraphqlSchemaRequest: CreateGraphqlSchemaRequest.AsObject
): Promise<GraphqlSchema.AsObject> {
  let request = new CreateGraphqlSchemaRequest();
  let { graphqlSchemaRef, spec } = createGraphqlSchemaRequest;
  let graphqlSchemaSpec = new GraphQLSchemaSpec();
  let executableSchema = new ExecutableSchema();
  executableSchema.setSchemaDefinition(
    spec?.executableSchema?.schemaDefinition ?? ''
  );
  graphqlSchemaSpec.setExecutableSchema(executableSchema);
  request.setGraphqlSchemaRef(
    getClusterRefClassFromClusterRefObj(graphqlSchemaRef!)
  );

  request.setSpec(graphqlSchemaSpec);

  return new Promise((resolve, reject) => {
    graphqlApiClient.createGraphqlSchema(request, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        resolve(data!.toObject().graphqlSchema!);
      }
    });
  });
}

async function updateGraphqlSchema(
  updateGraphqlSchemaRequest: UpdateGraphqlSchemaRequest.AsObject
): Promise<UpdateGraphqlSchemaResponse.AsObject> {
  let { graphqlSchemaRef, spec } = updateGraphqlSchemaRequest;
  let resolvers = spec?.executableSchema?.executor?.local?.resolutionsMap;

  let currentGraphqlSchema = await getGraphqlSchemaPb(graphqlSchemaRef!);

  let currentResolverMap = currentGraphqlSchema
    ?.getSpec()
    ?.getExecutableSchema()
    ?.getExecutor()
    ?.getLocal()
    ?.getResolutionsMap();

  // currentResolverMap.forEach(([key, value]) => newMetadata.getLabelsMap().set(key, value));
  let request = new CreateGraphqlSchemaRequest();
  let graphqlSchemaSpec = new GraphQLSchemaSpec();

  request.setGraphqlSchemaRef(
    getClusterRefClassFromClusterRefObj(graphqlSchemaRef!)
  );

  request.setSpec(graphqlSchemaSpec);

  return new Promise((resolve, reject) => {
    graphqlApiClient.createGraphqlSchema(request, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        resolve(data!.toObject());
      }
    });
  });
}

async function updateGraphqlSchemaResolver(
  graphqlSchemaRef: ClusterObjectRef.AsObject,
  resolverItem: [string, Resolution.AsObject]
): Promise<UpdateGraphqlSchemaResponse.AsObject> {
  let currentGraphqlSchema = await getGraphqlSchemaPb(graphqlSchemaRef!);

  let currentResolverMap = currentGraphqlSchema
    ?.getSpec()
    ?.getExecutableSchema()
    ?.getExecutor()
    ?.getLocal()
    ?.getResolutionsMap();

  // currentResolverMap.forEach(([key, value]) => newMetadata.getLabelsMap().set(key, value));
  let request = new CreateGraphqlSchemaRequest();
  let graphqlSchemaSpec = new GraphQLSchemaSpec();

  request.setGraphqlSchemaRef(
    getClusterRefClassFromClusterRefObj(graphqlSchemaRef!)
  );

  request.setSpec(graphqlSchemaSpec);

  return new Promise((resolve, reject) => {
    graphqlApiClient.createGraphqlSchema(request, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        resolve(data!.toObject());
      }
    });
  });
}

function deleteGraphqlSchema(
  graphqlSchemaRef: ClusterObjectRef.AsObject
): Promise<ClusterObjectRef.AsObject> {
  let request = new CreateGraphqlSchemaRequest();
  let graphqlSchemaSpec = new GraphQLSchemaSpec();

  request.setGraphqlSchemaRef(
    getClusterRefClassFromClusterRefObj(graphqlSchemaRef!)
  );

  request.setSpec(graphqlSchemaSpec);

  return new Promise((resolve, reject) => {
    graphqlApiClient.deleteGraphqlSchema(request, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        resolve(data!.toObject().graphqlSchemaRef!);
      }
    });
  });
}

function validateResolverYaml(
  validateResolverYamlRequest: ValidateResolverYamlRequest.AsObject
): Promise<ValidateResolverYamlResponse.AsObject> {
  let request = new ValidateResolverYamlRequest();
  let { resolverType, yaml } = validateResolverYamlRequest;
  request.setYaml(yaml);
  request.setResolverType(resolverType);

  return new Promise((resolve, reject) => {
    graphqlApiClient.validateResolverYaml(request, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        resolve(data!.toObject());
      }
    });
  });
}
