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
import { GraphqlConfigApiClient } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/graphql_pb_service';
import {
  GetGraphqlApiRequest,
  GraphqlApi,
  ListGraphqlApisRequest,
  GetGraphqlApiYamlRequest,
  ValidateResolverYamlRequest,
  ValidateResolverYamlResponse,
  CreateGraphqlApiRequest,
  CreateGraphqlApiResponse,
  UpdateGraphqlApiRequest,
  UpdateGraphqlApiResponse,
  DeleteGraphqlApiResponse,
  DeleteGraphqlApiRequest,
} from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/graphql_pb';
import {
  ExecutableSchema,
  Executor,
  GraphQLApiSpec,
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
const graphqlApiClient = new GraphqlConfigApiClient(host, {
  transport: grpc.CrossBrowserHttpTransport({ withCredentials: false }),
  debug: true,
});

export const graphqlConfigApi = {
  listGraphqlApis,
  getGraphqlApi,
  getGraphqlApiYaml,
  createGraphqlApi,
  validateResolverYaml,
  updateGraphqlApi,
  updateGraphqlApiIntrospection,
  deleteGraphqlApi,
  updateGraphqlApiResolver,
};

function listGraphqlApis(
  glooInstanceRef?: ObjectRef.AsObject
): Promise<GraphqlApi.AsObject[]> {
  let request = new ListGraphqlApisRequest();
  if (glooInstanceRef) {
    request.setGlooInstanceRef(getObjectRefClassFromRefObj(glooInstanceRef));
  }
  return new Promise((resolve, reject) => {
    graphqlApiClient.listGraphqlApis(request, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        resolve(data!.toObject().graphqlApisList);
      }
    });
  });
}

function getGraphqlApi(
  graphqlApiRef: ClusterObjectRef.AsObject
): Promise<GraphqlApi.AsObject> {
  let request = new GetGraphqlApiRequest();
  request.setGraphqlApiRef(getClusterRefClassFromClusterRefObj(graphqlApiRef));

  return new Promise((resolve, reject) => {
    graphqlApiClient.getGraphqlApi(request, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        resolve(data!.toObject()!.graphqlApi!);
      }
    });
  });
}

function getGraphqlApiPb(
  graphqlApiRef: ClusterObjectRef.AsObject
): Promise<GraphqlApi> {
  let request = new GetGraphqlApiRequest();
  request.setGraphqlApiRef(getClusterRefClassFromClusterRefObj(graphqlApiRef));

  return new Promise((resolve, reject) => {
    graphqlApiClient.getGraphqlApi(request, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        resolve(data!.getGraphqlApi()!);
      }
    });
  });
}

function getGraphqlApiYaml(
  graphqlApiRef: ClusterObjectRef.AsObject
): Promise<string> {
  let request = new GetGraphqlApiYamlRequest();
  request.setGraphqlApiRef(getClusterRefClassFromClusterRefObj(graphqlApiRef));

  return new Promise((resolve, reject) => {
    graphqlApiClient.getGraphqlApiYaml(request, (error, data) => {
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

function createGraphqlApi(
  createGraphqlApiRequest: CreateGraphqlApiRequest.AsObject
): Promise<GraphqlApi.AsObject> {
  let request = new CreateGraphqlApiRequest();
  let { graphqlApiRef, spec } = createGraphqlApiRequest;
  let graphqlApiSpec = new GraphQLApiSpec();
  let executableSchema = new ExecutableSchema();
  executableSchema.setSchemaDefinition(
    spec?.executableSchema?.schemaDefinition ?? ''
  );
  let local = new Executor.Local();
  local.setEnableIntrospection(true);
  let executor = new Executor();
  executor.setLocal(local);
  executableSchema.setExecutor(executor);
  graphqlApiSpec.setExecutableSchema(executableSchema);
  request.setGraphqlApiRef(getClusterRefClassFromClusterRefObj(graphqlApiRef!));

  request.setSpec(graphqlApiSpec);

  return new Promise((resolve, reject) => {
    graphqlApiClient.createGraphqlApi(request, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        resolve(data!.toObject().graphqlApi!);
      }
    });
  });
}

function apiSpecFromObject(
  apiSpec: GraphQLApiSpec.AsObject,
  apiSpecToUpdate = new GraphQLApiSpec()
): GraphQLApiSpec {
  let { executableSchema, statPrefix } = apiSpec;
  if (statPrefix !== undefined) {
    let { value } = statPrefix;
    let newStatPrefix = apiSpecToUpdate.getStatPrefix() ?? new StringValue();
    newStatPrefix.setValue(value);
    apiSpecToUpdate.setStatPrefix(newStatPrefix);
  }

  if (executableSchema !== undefined) {
    let { schemaDefinition, executor, grpcDescriptorRegistry } =
      executableSchema;

    let newExecutableSchema =
      apiSpecToUpdate.getExecutableSchema() ?? new ExecutableSchema();

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

    apiSpecToUpdate.setExecutableSchema(newExecutableSchema);
  }
  return apiSpecToUpdate;
}

async function updateGraphqlApi(
  updateGraphqlApiRequest: Partial<UpdateGraphqlApiRequest.AsObject>
): Promise<GraphqlApi.AsObject> {
  let { graphqlApiRef, spec } = updateGraphqlApiRequest;
  let currentGraphqlApi = await getGraphqlApiPb(graphqlApiRef!);

  let request = new UpdateGraphqlApiRequest();
  request.setGraphqlApiRef(getClusterRefClassFromClusterRefObj(graphqlApiRef!));

  let graphqlApiSpec = apiSpecFromObject(spec!, currentGraphqlApi?.getSpec());
  request.setSpec(graphqlApiSpec);

  return new Promise((resolve, reject) => {
    graphqlApiClient.updateGraphqlApi(request, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        resolve(data!.toObject().graphqlApi!);
      }
    });
  });
}

async function updateGraphqlApiIntrospection(
  graphqlApiRef: ClusterObjectRef.AsObject,
  introspectionEnabled: boolean
): Promise<UpdateGraphqlApiResponse.AsObject> {
  let currentGraphqlApi = await getGraphqlApiPb(graphqlApiRef!);

  // currentResolverMap.forEach(([key, value]) => newMetadata.getLabelsMap().set(key, value));
  let request = new UpdateGraphqlApiRequest();
  let graphqlApiSpec = currentGraphqlApi?.getSpec();

  request.setGraphqlApiRef(getClusterRefClassFromClusterRefObj(graphqlApiRef!));
  graphqlApiSpec
    ?.getExecutableSchema()
    ?.getExecutor()
    ?.getLocal()
    ?.setEnableIntrospection(introspectionEnabled);
  request.setSpec(graphqlApiSpec);

  return new Promise((resolve, reject) => {
    graphqlApiClient.updateGraphqlApi(request, (error, data) => {
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

async function updateGraphqlApiResolver(
  graphqlApiRef: ClusterObjectRef.AsObject,
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
): Promise<GraphqlApi.AsObject> {
  let currentGraphqlApi = await getGraphqlApiPb(graphqlApiRef!);

  let currentSpec = currentGraphqlApi?.getSpec();

  if (currentSpec === undefined) {
    currentSpec = new GraphQLApiSpec();
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

  let request = new UpdateGraphqlApiRequest();

  request.setGraphqlApiRef(getClusterRefClassFromClusterRefObj(graphqlApiRef!));

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
    graphqlApiClient.updateGraphqlApi(request, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        resolve(data!.toObject().graphqlApi!);
      }
    });
  });
}

function deleteGraphqlApi(
  graphqlApiRef: ClusterObjectRef.AsObject
): Promise<ClusterObjectRef.AsObject> {
  let request = new CreateGraphqlApiRequest();
  let graphqlApiSpec = new GraphQLApiSpec();

  request.setGraphqlApiRef(getClusterRefClassFromClusterRefObj(graphqlApiRef!));

  request.setSpec(graphqlApiSpec);

  return new Promise((resolve, reject) => {
    graphqlApiClient.deleteGraphqlApi(request, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        resolve(data!.toObject().graphqlApiRef!);
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
