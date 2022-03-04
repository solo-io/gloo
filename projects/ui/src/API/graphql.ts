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
  Executor,
  GraphQLSchemaSpec,
  GrpcDescriptorRegistry,
  GrpcRequestTemplate,
  RequestTemplate,
  Resolution,
  ResponseTemplate,
  RESTResolver,
} from 'proto/github.com/solo-io/solo-apis/api/gloo/graphql.gloo/v1alpha1/graphql_pb';
import { StringValue } from 'google-protobuf/google/protobuf/wrappers_pb';
import { Struct, Value } from 'google-protobuf/google/protobuf/struct_pb';
import { ResourceRef } from 'proto/github.com/solo-io/solo-kit/api/v1/ref_pb';
import { struct } from 'pb-util';
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
  updateGraphqlSchemaIntrospection,
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
  let local = new Executor.Local();
  local.setEnableIntrospection(true);
  let executor = new Executor();
  executor.setLocal(local);
  executableSchema.setExecutor(executor);
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

function schemaSpecFromObject(
  schemaSpec: GraphQLSchemaSpec.AsObject,
  schemaSpecToUpdate = new GraphQLSchemaSpec()
): GraphQLSchemaSpec {
  let { executableSchema, statPrefix } = schemaSpec;
  if (statPrefix !== undefined) {
    let { value } = statPrefix;
    let newStatPrefix = schemaSpecToUpdate.getStatPrefix() ?? new StringValue();
    newStatPrefix.setValue(value);
    schemaSpecToUpdate.setStatPrefix(newStatPrefix);
  }

  if (executableSchema !== undefined) {
    let { schemaDefinition, executor, grpcDescriptorRegistry } =
      executableSchema;

    let newExecutableSchema =
      schemaSpecToUpdate.getExecutableSchema() ?? new ExecutableSchema();

    if (!!schemaDefinition) {
      newExecutableSchema.setSchemaDefinition(schemaDefinition);
    }

    if (executor !== undefined) {
      let { local } = executor;
      let newExecutor = newExecutableSchema.getExecutor() ?? new Executor();

      if (local !== undefined) {
        let { enableIntrospection, resolutionsMap } = local;
        let newLocal = newExecutor.getLocal() ?? new Executor.Local();

        if (enableIntrospection !== undefined) {
          newLocal.setEnableIntrospection(enableIntrospection);
        }

        if (resolutionsMap !== undefined) {
          let newResolutionsMap = newLocal.getResolutionsMap();
          newResolutionsMap.forEach((resolution, resolutionName) => {
            newResolutionsMap.set(resolutionName, resolution);
          });
        }

        newExecutor.setLocal(newLocal);
      }

      newExecutableSchema.setExecutor(newExecutor);
    }

    if (grpcDescriptorRegistry !== undefined) {
      let { protoDescriptor, protoDescriptorBin } = grpcDescriptorRegistry;
      let newGrpcDescriptorRegistry =
        newExecutableSchema?.getGrpcDescriptorRegistry() ??
        new GrpcDescriptorRegistry();
      newGrpcDescriptorRegistry.setProtoDescriptor(protoDescriptor);
      newGrpcDescriptorRegistry.setProtoDescriptorBin(protoDescriptorBin);

      newExecutableSchema.setGrpcDescriptorRegistry(newGrpcDescriptorRegistry);
    }

    schemaSpecToUpdate.setExecutableSchema(newExecutableSchema);
  }
  return schemaSpecToUpdate;
}

async function updateGraphqlSchema(
  updateGraphqlSchemaRequest: Partial<UpdateGraphqlSchemaRequest.AsObject>
): Promise<GraphqlSchema.AsObject> {
  let { graphqlSchemaRef, spec } = updateGraphqlSchemaRequest;
  let currentGraphqlSchema = await getGraphqlSchemaPb(graphqlSchemaRef!);

  let request = new UpdateGraphqlSchemaRequest();
  request.setGraphqlSchemaRef(
    getClusterRefClassFromClusterRefObj(graphqlSchemaRef!)
  );

  let graphqlSchemaSpec = schemaSpecFromObject(
    spec!,
    currentGraphqlSchema?.getSpec()
  );
  request.setSpec(graphqlSchemaSpec);

  return new Promise((resolve, reject) => {
    graphqlApiClient.updateGraphqlSchema(request, (error, data) => {
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

async function updateGraphqlSchemaIntrospection(
  graphqlSchemaRef: ClusterObjectRef.AsObject,
  introspectionEnabled: boolean
): Promise<UpdateGraphqlSchemaResponse.AsObject> {
  let currentGraphqlSchema = await getGraphqlSchemaPb(graphqlSchemaRef!);

  // currentResolverMap.forEach(([key, value]) => newMetadata.getLabelsMap().set(key, value));
  let request = new UpdateGraphqlSchemaRequest();
  let graphqlSchemaSpec = currentGraphqlSchema?.getSpec();

  request.setGraphqlSchemaRef(
    getClusterRefClassFromClusterRefObj(graphqlSchemaRef!)
  );
  graphqlSchemaSpec
    ?.getExecutableSchema()
    ?.getExecutor()
    ?.getLocal()
    ?.setEnableIntrospection(introspectionEnabled);
  request.setSpec(graphqlSchemaSpec);

  return new Promise((resolve, reject) => {
    graphqlApiClient.updateGraphqlSchema(request, (error, data) => {
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
  resolverItem: {
    resolverName: string;
    resolverType?: 'REST' | 'gRPC';
    request?: RequestTemplate.AsObject;
    response?: ResponseTemplate.AsObject;
    grpcRequest?: GrpcRequestTemplate.AsObject;
    spanName?: string;
    upstreamRef?: ObjectRef.AsObject;
    hasDirective?: boolean;
    fieldWithDirective?: string;
    fieldWithoutDirective?: string;
  },
  isRemove?: boolean
): Promise<GraphqlSchema.AsObject> {
  let currentGraphqlSchema = await getGraphqlSchemaPb(graphqlSchemaRef!);

  let currentSpec = currentGraphqlSchema?.getSpec();

  if (currentSpec === undefined) {
    currentSpec = new GraphQLSchemaSpec();
  }
  let currentExSchema = currentSpec?.getExecutableSchema();
  if (currentExSchema === undefined) {
    currentExSchema = new ExecutableSchema();
  }
  let currExecutor = currentExSchema?.getExecutor();
  if (currExecutor === undefined) {
    currExecutor = new Executor();
  }
  let currLocal = currExecutor.getLocal();
  if (currLocal === undefined) {
    currLocal = new Executor.Local();
  }

  let currResolMap = currLocal?.getResolutionsMap();

  let newResolution = new Resolution();

  if (resolverItem.resolverType === 'REST') {
    let newRestResolver =
      currResolMap?.get(resolverItem.resolverName)?.getRestResolver() ??
      new RESTResolver();
    let usRef = newRestResolver?.getUpstreamRef() ?? new ResourceRef();
    if (!usRef?.toObject().name || !usRef?.toObject().namespace) {
      usRef.setName(resolverItem?.upstreamRef?.name!);
      usRef.setNamespace(resolverItem?.upstreamRef?.namespace!);
    }

    newRestResolver.setUpstreamRef(usRef);
    if (resolverItem.request !== undefined) {
      let { headersMap, queryParamsMap, body } = resolverItem.request;
      let newReq = newRestResolver?.getRequest() ?? new RequestTemplate();

      if (body !== undefined) {
        let bodyVal = new Value();
        bodyVal.setStringValue(body.stringValue);
        newReq.setBody(bodyVal);
      } else {
        newReq.clearBody();
        newReq.setBody(undefined);
      }

      if (headersMap?.length > 0) {
        let newHeadersMap = newReq.getHeadersMap();
        headersMap.forEach(([val, key]) => {
          newHeadersMap.set(val, key);
        });
      } else {
        newReq.clearHeadersMap();
      }

      if (queryParamsMap?.length > 0) {
        let qParamsMap = newReq.getQueryParamsMap();
        queryParamsMap.forEach(([val, key]) => {
          qParamsMap.set(val, key);
        });
      } else {
        newReq.clearQueryParamsMap();
      }
      newRestResolver.setRequest(newReq);
    } else {
      newRestResolver?.clearRequest();
    }

    if (resolverItem.response !== undefined) {
      let { resultRoot, settersMap } = resolverItem.response;
      let newRes = newRestResolver.getResponse() ?? new ResponseTemplate();
      if (resultRoot !== undefined && resultRoot !== '') {
        newRes.setResultRoot(resultRoot);
      }
      if (settersMap?.length > 0) {
        let newSettersMap = newRes.getSettersMap();
        settersMap.forEach(([key, val]) => {
          newSettersMap.set(val, key);
        });
      }
      newRestResolver.setResponse(newRes);
    } else {
      newRestResolver.clearResponse();
    }
    newResolution.setRestResolver(newRestResolver);
  }

  let request = new UpdateGraphqlSchemaRequest();

  request.setGraphqlSchemaRef(
    getClusterRefClassFromClusterRefObj(graphqlSchemaRef!)
  );

  currExecutor.setLocal(currLocal);
  currentExSchema.setExecutor(currExecutor);
  let currentSchemaDef = currentExSchema.getSchemaDefinition();
  // TODO: find a better way to do this
  let { fieldWithDirective, fieldWithoutDirective } = resolverItem;
  if (
    !isRemove &&
    !resolverItem.hasDirective &&
    fieldWithDirective &&
    fieldWithoutDirective
  ) {
    currentExSchema.setSchemaDefinition(
      currentSchemaDef.replace(fieldWithoutDirective, fieldWithDirective)
    );
  }
  if (isRemove) {
    currResolMap.del(resolverItem.resolverName);
    if (!!fieldWithDirective && fieldWithoutDirective) {
      currentExSchema.setSchemaDefinition(
        currentSchemaDef.replace(fieldWithDirective, fieldWithoutDirective)
      );
    }
  } else {
    currResolMap.set(resolverItem.resolverName, newResolution);
  }

  currentSpec.setExecutableSchema(currentExSchema);
  request.setSpec(currentSpec);

  return new Promise((resolve, reject) => {
    graphqlApiClient.updateGraphqlSchema(request, (error, data) => {
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
