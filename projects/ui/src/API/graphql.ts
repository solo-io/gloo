import { grpc } from '@improbable-eng/grpc-web';
import { Value } from 'google-protobuf/google/protobuf/struct_pb';
import { StringValue } from 'google-protobuf/google/protobuf/wrappers_pb';
import { FieldDefinitionNode } from 'graphql';
import isEmpty from 'lodash/isEmpty';
import {
  ClusterObjectRef,
  ObjectRef,
} from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';
import {
  ExecutableSchema,
  Executor,
  GraphQLApiSpec,
  GrpcDescriptorRegistry,
  GrpcRequestTemplate,
  GrpcResolver,
  RequestTemplate,
  Resolution,
  ResponseTemplate,
  RESTResolver,
  StitchedSchema,
} from 'proto/github.com/solo-io/solo-apis/api/gloo/graphql.gloo/v1beta1/graphql_pb';
import { ResourceRef } from 'proto/github.com/solo-io/solo-kit/api/v1/ref_pb';
import { ServiceError } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/gateway_resources_pb_service';
import {
  CreateGraphqlApiRequest,
  DeleteGraphqlApiRequest,
  GetGraphqlApiRequest,
  GetGraphqlApiYamlRequest,
  GetStitchedSchemaDefinitionRequest,
  GraphqlApi,
  ListGraphqlApisRequest,
  UpdateGraphqlApiRequest,
  UpdateGraphqlApiResponse,
  ValidateResolverYamlRequest,
  ValidateResolverYamlResponse,
  ValidateSchemaDefinitionRequest,
  ValidateSchemaDefinitionResponse,
} from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/graphql_pb';
import { GraphqlConfigApiClient } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/graphql_pb_service';
import { updateSchemaAndResolutionMap } from './graphql-resolutions';
import {
  getClusterRefClassFromClusterRefObj,
  getObjectRefClassFromRefObj,
  host,
} from './helpers';
const graphqlApiClient = new GraphqlConfigApiClient(host, {
  transport: grpc.CrossBrowserHttpTransport({ withCredentials: false }),
  debug: true,
});

export type ResolverItem = {
  field: FieldDefinitionNode;
  resolverType?: 'REST' | 'gRPC';
  request?: RequestTemplate.AsObject;
  response?: ResponseTemplate.AsObject;
  grpcRequest?: GrpcRequestTemplate.AsObject;
  spanName?: string;
  upstreamRef?: ObjectRef.AsObject;
  isNewResolution: boolean;
  fieldReturnType: string;
  objectType: string;
};

export const graphqlConfigApi = {
  listGraphqlApis,
  getGraphqlApi,
  getGraphqlApiYaml,
  createGraphqlApi,
  validateResolverYaml,
  validateSchema,
  updateGraphqlApi,
  updateGraphqlApiIntrospection,
  deleteGraphqlApi,
  updateGraphqlApiResolver,
  getGraphqlApiWithResolver,
  getStitchedSchemaDefinition,
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
      if (error !== null) rejectWithError(reject, error);
      else {
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
      if (error !== null) rejectWithError(reject, error);
      else {
        resolve(data!.toObject()!.graphqlApi!);
      }
    });
  });
}

export function getGraphqlApiPb(
  graphqlApiRef: ClusterObjectRef.AsObject
): Promise<GraphqlApi> {
  let request = new GetGraphqlApiRequest();
  request.setGraphqlApiRef(getClusterRefClassFromClusterRefObj(graphqlApiRef));

  return new Promise((resolve, reject) => {
    graphqlApiClient.getGraphqlApi(request, (error, data) => {
      if (error !== null) rejectWithError(reject, error);
      else {
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
      if (error !== null) rejectWithError(reject, error);
      else {
        resolve(data!.toObject().yamlData?.yaml ?? 'None');
      }
    });
  });
}

function createGraphqlApi({
  graphqlApiRef,
  spec,
}: CreateGraphqlApiRequest.AsObject): Promise<GraphqlApi.AsObject> {
  let request = new CreateGraphqlApiRequest();
  let graphqlApiSpec = new GraphQLApiSpec();

  // Check API type
  const isStitchedApi = !!spec?.stitchedSchema;
  const isExecutableApi = !!spec?.executableSchema;
  if (!isStitchedApi && !isExecutableApi)
    return new Promise((resolve, reject) => reject('Invalid API type!'));

  if (isStitchedApi) {
    // -- Stitched
    let stitchedSchema = new StitchedSchema();
    // Uncomment the following lines to test schema persistence:
    // const mockSubschema = new StitchedSchema.SubschemaConfig();
    // mockSubschema.setName('test-name');
    // mockSubschema.setNamespace('test-namespace');
    // stitchedSchema.setSubschemasList([mockSubschema]);
    stitchedSchema.setSubschemasList([]);
    graphqlApiSpec.setStitchedSchema(stitchedSchema);
  } else if (isExecutableApi) {
    // -- Executable
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
  }

  request.setGraphqlApiRef(getClusterRefClassFromClusterRefObj(graphqlApiRef!));
  request.setSpec(graphqlApiSpec);

  return new Promise((resolve, reject) => {
    graphqlApiClient.createGraphqlApi(request, (error, data) => {
      if (error !== null) rejectWithError(reject, error);
      else {
        resolve(data!.toObject().graphqlApi!);
      }
    });
  });
}

function apiSpecFromObject(
  apiSpec: GraphQLApiSpec.AsObject,
  apiSpecToUpdate = new GraphQLApiSpec()
): GraphQLApiSpec {
  let { executableSchema, stitchedSchema, statPrefix } = apiSpec;
  if (statPrefix !== undefined) {
    let { value } = statPrefix;
    let newStatPrefix = apiSpecToUpdate.getStatPrefix() ?? new StringValue();
    newStatPrefix.setValue(value);
    apiSpecToUpdate.setStatPrefix(newStatPrefix);
  }

  // -- Stitched -- //
  if (stitchedSchema !== undefined) {
    let subschemasList = stitchedSchema.subschemasList;
    let newSchema = apiSpecToUpdate.getStitchedSchema() ?? new StitchedSchema();
    const newSubschemasList = subschemasList.map(subschema => {
      const newSubschema = new StitchedSchema.SubschemaConfig();
      // Set name and namespace.
      newSubschema.setName(subschema.name);
      newSubschema.setNamespace(subschema.namespace);
      // Set type merge map in place.
      const mergeMap = newSubschema.getTypeMergeMap();
      subschema.typeMergeMap.forEach(m => {
        // m[0] === typeName
        // m[1] === typeMergeConfig
        const mergeConfig =
          new StitchedSchema.SubschemaConfig.TypeMergeConfig();
        const argsMap = mergeConfig.getArgsMap();
        m[1].argsMap.forEach(arg => argsMap.set(arg[0], arg[1]));
        mergeConfig.setSelectionSet(m[1].selectionSet);
        mergeConfig.setQueryName(m[1].queryName);
        mergeMap.set(m[0], mergeConfig);
      });
      return newSubschema;
    });
    newSchema.setSubschemasList(newSubschemasList);
    apiSpecToUpdate.setStitchedSchema(newSchema);
  }

  // -- Executable -- //
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
        // TODO:  This doesn't actually update the resolutions map because that's in the apiSec
        //        Also, this is another weird flipped conversion.
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
    return graphqlApiClient.updateGraphqlApi(request, (error, data) => {
      if (error !== null) rejectWithError(reject, error);
      else {
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
    return graphqlApiClient.updateGraphqlApi(request, (error, data) => {
      if (error !== null) rejectWithError(reject, error);
      else {
        resolve(data!.toObject());
      }
    });
  });
}

async function updateGraphqlApiResolver(
  graphqlApiRef: ClusterObjectRef.AsObject,
  resolverItem: ResolverItem,
  shouldDelete?: boolean
): Promise<GraphqlApi.AsObject> {
  const fieldName = resolverItem.field.name.value;
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
      currResolMap?.get(fieldName)?.getRestResolver() ?? new RESTResolver();
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
        newHeadersMap.clear();
        headersMap.forEach(([key, val]) => {
          newHeadersMap.set(key, val);
        });
      } else {
        newReq.clearHeadersMap();
      }

      if (queryParamsMap?.length > 0) {
        let qParamsMap = newReq.getQueryParamsMap();
        queryParamsMap.forEach(([key, val]) => {
          qParamsMap.set(key, val);
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
        newSettersMap.clear();
        settersMap.forEach(([key, val]) => {
          newSettersMap.set(key, val);
        });
      }
      newRestResolver.setResponse(newRes);
    } else {
      newRestResolver.clearResponse();
    }
    newResolution.setRestResolver(newRestResolver);
  } else if (resolverItem.resolverType === 'gRPC') {
    let newGrpcResolver =
      currResolMap?.get(fieldName)?.getGrpcResolver() ?? new GrpcResolver();
    let usRef = newGrpcResolver?.getUpstreamRef() ?? new ResourceRef();
    if (!usRef?.toObject().name || !usRef?.toObject().namespace) {
      usRef.setName(resolverItem?.upstreamRef?.name!);
      usRef.setNamespace(resolverItem?.upstreamRef?.namespace!);
    }
    newGrpcResolver.setUpstreamRef(usRef);

    if (resolverItem.grpcRequest !== undefined) {
      let { methodName, requestMetadataMap, serviceName, outgoingMessageJson } =
        resolverItem.grpcRequest;
      let newReq =
        newGrpcResolver?.getRequestTransform() ?? new GrpcRequestTemplate();

      if (methodName !== undefined && methodName !== '') {
        newReq.setMethodName(methodName);
      }

      if (requestMetadataMap?.length > 0) {
        let newHeadersMap = newReq.getRequestMetadataMap();
        newHeadersMap.clear();
        requestMetadataMap.forEach(([key, val]) => {
          newHeadersMap.set(key, val);
        });
      } else {
        newReq.clearRequestMetadataMap();
      }

      if (serviceName !== undefined && serviceName !== '') {
        newReq.setServiceName(serviceName);
      }

      if (outgoingMessageJson !== undefined) {
        let outgoingJson = new Value();
        outgoingJson.setStringValue(outgoingMessageJson.stringValue);
        newReq.setOutgoingMessageJson(outgoingJson);
      } else {
        newReq.clearOutgoingMessageJson();
        newReq.setOutgoingMessageJson(undefined);
      }

      newGrpcResolver.setRequestTransform(newReq);
    } else {
      newGrpcResolver?.clearRequestTransform();
    }

    if (resolverItem.spanName !== undefined && resolverItem.spanName !== '') {
      newGrpcResolver.setSpanName(resolverItem.spanName);
    }
    newResolution.setGrpcResolver(newGrpcResolver);
  }

  let request = new UpdateGraphqlApiRequest();

  request.setGraphqlApiRef(getClusterRefClassFromClusterRefObj(graphqlApiRef!));

  currExecutor.setLocal(currLocal);
  currentExSchema.setExecutor(currExecutor);

  updateSchemaAndResolutionMap(
    resolverItem,
    newResolution,
    currentExSchema,
    currResolMap,
    shouldDelete
  );

  currentSpec.setExecutableSchema(currentExSchema);
  request.setSpec(currentSpec);

  return new Promise((resolve, reject) => {
    return graphqlApiClient.updateGraphqlApi(request, (error, data) => {
      if (error !== null) rejectWithError(reject, error);
      else {
        resolve(data!.toObject().graphqlApi!);
      }
    });
  });
}

async function getGraphqlApiWithResolver(
  graphqlApiRef: ClusterObjectRef.AsObject,
  resolverItem: ResolverItem,
  shouldDelete?: boolean
): Promise<GraphQLApiSpec> {
  const fieldName = resolverItem.field.name.value;
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
      currResolMap?.get(fieldName)?.getRestResolver() ?? new RESTResolver();
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
        newHeadersMap.clear();
        headersMap.forEach(([key, val]) => {
          newHeadersMap.set(key, val);
        });
      } else {
        newReq.clearHeadersMap();
      }

      if (queryParamsMap?.length > 0) {
        let qParamsMap = newReq.getQueryParamsMap();
        qParamsMap.clear();
        queryParamsMap.forEach(([key, val]) => {
          qParamsMap.set(key, val);
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
        newSettersMap.clear();
        settersMap.forEach(([key, val]) => {
          newSettersMap.set(key, val);
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

  updateSchemaAndResolutionMap(
    resolverItem,
    newResolution,
    currentExSchema,
    currResolMap,
    shouldDelete
  );

  currentSpec.setExecutableSchema(currentExSchema);
  request.setSpec(currentSpec);

  return request.getSpec()!;
}

function deleteGraphqlApi(
  graphqlApiRef: ClusterObjectRef.AsObject
): Promise<ClusterObjectRef.AsObject> {
  let request = new DeleteGraphqlApiRequest();
  request.setGraphqlApiRef(getClusterRefClassFromClusterRefObj(graphqlApiRef!));
  return new Promise((resolve, reject) => {
    graphqlApiClient.deleteGraphqlApi(request, (error, data) => {
      if (error !== null) rejectWithError(reject, error);
      else {
        resolve(data!.toObject().graphqlApiRef!);
      }
    });
  });
}

/**
 * When creating a new GraphQLSchema from scratch, the schema definition string should be passed in.
 * When editing an existing GraphQLSchema, the full GraphQLSchema spec should be passed in.
 *
 * An empty response is returned if validation succeeded. Otherwise, an error is returned.
 */
async function validateSchema(
  validationRequest: ValidateSchemaDefinitionRequest.AsObject & {
    apiRef?: ClusterObjectRef.AsObject;
    resolverItem?: ResolverItem;
  }
): Promise<ValidateSchemaDefinitionResponse.AsObject> {
  let request = new ValidateSchemaDefinitionRequest();
  const { schemaDefinition, spec, apiRef, resolverItem } = validationRequest;
  if (schemaDefinition) {
    request.setSchemaDefinition(schemaDefinition);
  }
  if (spec && resolverItem && apiRef) {
    const apiSpec = await getGraphqlApiWithResolver(apiRef!, resolverItem);
    request.setSpec(apiSpec);
  } else if (apiRef && spec) {
    const currentGraphqlApi = await getGraphqlApiPb(apiRef);
    const currentSpec = currentGraphqlApi!.getSpec()!;
    const executable = currentSpec.getExecutableSchema()!;
    executable.setSchemaDefinition(spec!.executableSchema!.schemaDefinition);
    currentSpec.setExecutableSchema(executable);
    request.setSpec(currentSpec);
  } else if (spec) {
    const newSpec = new GraphQLApiSpec();
    newSpec.setAllowedQueryHashesList(spec.allowedQueryHashesList);
    const executable = new ExecutableSchema();
    executable.setSchemaDefinition(spec!.executableSchema!.schemaDefinition);
    const executor = new Executor();
    const local = new Executor.Local();
    local.setEnableIntrospection(
      spec.executableSchema!.executor!.local!.enableIntrospection
    );
    executor.setLocal(local);
    executable.setExecutor(executor);
    newSpec.setExecutableSchema(executable);
    request.setSpec(newSpec);
  }

  return new Promise((resolve, reject) => {
    return graphqlApiClient.validateSchemaDefinition(request, (err, data) => {
      if (err) {
        return reject(err);
      }
      const dataObj = data!.toObject();

      if (isEmpty(dataObj)) {
        return resolve(dataObj);
      } else {
        return reject(dataObj);
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
      if (error !== null) rejectWithError(reject, error);
      else {
        resolve(data!.toObject());
      }
    });
  });
}

function getStitchedSchemaDefinition(
  graphqlApiRef: ClusterObjectRef.AsObject
): Promise<string> {
  let request = new GetStitchedSchemaDefinitionRequest();
  request.setStitchedSchemaApiRef(
    getClusterRefClassFromClusterRefObj(graphqlApiRef)
  );
  return new Promise((resolve, reject) => {
    graphqlApiClient.getStitchedSchemaDefinition(request, (error, data) => {
      if (error !== null) rejectWithError(reject, error);
      else resolve(data!.toObject()!.stitchedSchemaSdl!);
    });
  });
}

const rejectWithError = (
  reject: (reason: any) => void,
  error: ServiceError
) => {
  console.error('Error:', error.message);
  console.error('Code:', error.code);
  console.error('Metadata:', error.metadata);
  reject(error);
};
