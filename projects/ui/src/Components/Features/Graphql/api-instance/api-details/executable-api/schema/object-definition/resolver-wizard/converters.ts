import { ResolverItem } from 'API/graphql';
import { Resolution } from 'proto/github.com/solo-io/solo-apis/api/gloo/graphql.gloo/v1beta1/graphql_pb';
import YAML from 'yaml';
import cloneDeep from 'lodash/cloneDeep';
import { getDefaultConfigFromType } from './ResolverConfigSection';
import { FieldDefinitionNode } from 'graphql';

export const removeNulls = (obj: any) => {
  const isArray = Array.isArray(obj);
  for (const k of Object.keys(obj)) {
    if (obj[k] === null || obj[k] === undefined || obj[k] === '') {
      if (isArray) {
        obj.splice(Number(k), 1);
      } else {
        delete obj[k];
      }
    } else if (typeof obj[k] === 'object') {
      removeNulls(obj[k]);
    }
    if (isArray && obj.length === Number(k)) {
      removeNulls(obj);
    }
    if (typeof obj[k] === 'object' && Object.keys(obj[k]).length === 0) {
      delete obj[k];
    }
  }
  return obj;
};

export const createResolverItem = (
  resolverConfig: string,
  resolverType: 'gRPC' | 'REST',
  field: FieldDefinitionNode,
  upstream: string,
  extras: Record<string, any> = {}
): ResolverItem => {
  /*
     `parsedResolverConfig` can be formatted in different ways:
     - `restResolver.[request | response | spanName | ...]`....
     - `grpcResolver.[request | response | spanName | ...]`...
     - `[request | response | spanName | ...]`...
  */
  let parsedResolverConfig;
  try {
    parsedResolverConfig = removeNulls(YAML.parse(resolverConfig));
  } catch (err: any) {
    throw err;
  }

  let headersMap: [string, string][] = [];
  let queryParamsMap: [string, string][] = [];
  let settersMap: [string, string][] = [];
  let requestMetadataMap: [string, string][] = [];
  let serviceName = '';
  let methodName = '';
  let outgoingMessageJson;
  let body;

  let resultRoot =
    resolverType === 'gRPC' && parsedResolverConfig?.grpcResolver?.resultRoot
      ? parsedResolverConfig?.grpcResolver?.resultRoot
      : parsedResolverConfig?.resultRoot;
  let spanName;

  if (parsedResolverConfig?.grpcResolver?.spanName) {
    spanName = parsedResolverConfig.grpcResolver.spanName;
  } else if (parsedResolverConfig?.restResolver?.spanName) {
    spanName = parsedResolverConfig.restResolver.spanName;
  } else if (parsedResolverConfig.spanName) {
    spanName = parsedResolverConfig.spanName;
  }

  let requestTransform =
    resolverType === 'gRPC' &&
    parsedResolverConfig?.grpcResolver?.requestTransform
      ? parsedResolverConfig.grpcResolver.requestTransform
      : parsedResolverConfig.requestTransform;
  let request =
    resolverType === 'REST' && parsedResolverConfig?.restResolver?.request
      ? parsedResolverConfig.restResolver.request
      : parsedResolverConfig.request;
  let response =
    resolverType === 'REST' && parsedResolverConfig?.restResolver?.response
      ? parsedResolverConfig.restResolver.response
      : parsedResolverConfig.response;

  if (resolverType === 'REST') {
    if (request) {
      const headerValues = request?.headers ?? {};

      headersMap = Object.keys(headerValues).reduce(
        (acc: [string, string][], curr: string) => {
          acc.push([curr, headerValues[curr]]);
          return acc;
        },
        []
      );

      const queryValues = request?.queryParams ?? {};

      queryParamsMap = Object.keys(queryValues).reduce(
        (acc: [string, string][], curr: string) => {
          acc.push([curr, queryValues[curr]]);
          return acc;
        },
        []
      );

      body = request?.body;
    }
    if (response) {
      const values = response?.settersMap ?? response?.setters ?? {};
      settersMap = Object.keys(values).reduce(
        (acc: [string, string][], curr: string) => {
          acc.push([curr, values[curr]]);
          return acc;
        },
        []
      );
      resultRoot = response?.resultRoot;
    }
  } else {
    if (resolverType === 'gRPC' && requestTransform) {
      const requestMetaValues = requestTransform?.requestMetadataMap ?? {};
      requestMetadataMap = Object.keys(requestMetaValues).reduce(
        (acc: [string, string][], curr: string) => {
          acc.push([curr, requestMetaValues[curr]]);
          return acc;
        },
        []
      );
      serviceName = requestTransform?.serviceName;
      methodName = requestTransform?.methodName;
      outgoingMessageJson = requestTransform?.outgoingMessageJson;
    }
  }
  let [upstreamName, upstreamNamespace] = upstream.split('::');

  const resolverItem = {
    upstreamRef: {
      name: upstreamName,
      namespace: upstreamNamespace,
    },
    //@ts-ignore
    ...(request && {
      request: {
        headersMap,
        queryParamsMap,
        body,
      },
    }),
    field,
    //@ts-ignore
    ...(requestTransform && {
      grpcRequest: {
        methodName,
        requestMetadataMap,
        serviceName,
        outgoingMessageJson,
      },
    }),
    resolverType: resolverType,
    //@ts-ignore
    ...(response && { response: { resultRoot, settersMap } }),
    spanName,
    ...extras,
  } as ResolverItem;
  return resolverItem;
};

export const getResolverFromConfig = (resolver?: Resolution.AsObject) => {
  if (resolver?.restResolver || resolver?.grpcResolver) {
    // TODO:  This conversion doesn't quite work.
    // conversion: Resolution.AsObject -> protobufjs.Message -> proto3 JsonValue -> string
    let parsed: Record<string, any> = cloneDeep(resolver);
    if (parsed.restResolver) {
      parsed = cloneDeep(parsed.restResolver);
    } else if (parsed.grpcResolver) {
      parsed = cloneDeep(parsed.grpcResolver);
    }
    delete parsed.upstreamRef;

    if (parsed?.request?.headersMap) {
      parsed.request.headers = parsed.request.headersMap.reduce(
        (acc: any, curr: [string, string]) => {
          const [key, value] = curr;
          return {
            ...acc,
            [key]: value,
          };
        },
        {}
      );
      if (!parsed.request.headersMap.length) {
        delete parsed.request.headers;
      }
      delete parsed.request.headersMap;
    }

    if (parsed?.request?.queryParamsMap) {
      parsed.request.queryParams = parsed.request.queryParamsMap.reduce(
        (acc: any, curr: [string, string]) => {
          const [key, value] = curr;
          return {
            ...acc,
            [key]: value,
          };
        },
        {}
      );
      if (!parsed.request.queryParamsMap.length) {
        delete parsed.request.queryParams;
      }
      delete parsed.request.queryParamsMap;
    }
    if (parsed?.response?.settersMap) {
      parsed.response.setters = parsed.response.settersMap.reduce(
        (acc: any, curr: [string, string]) => {
          const [key, value] = curr;
          return {
            ...acc,
            [key]: value,
          };
        },
        {}
      );
      if (!parsed.response.settersMap.length) {
        delete parsed.response.setters;
      }
      delete parsed.response.settersMap;
    }

    if (parsed?.requestTransform?.requestMetadataMap) {
      parsed.requestTransform.requestMetadata =
        parsed.requestTransform.requestMetadataMap.reduce(
          (acc: any, curr: any) => {
            const [key, value] = curr;
            return {
              ...acc,
              [key]: value,
            };
          },
          {}
        );
      if (!parsed.requestTransform.requestMetadataMap.length) {
        delete parsed.requestTransform.requestMetadata;
      }
      delete parsed.requestTransform.requestMetadataMap;
    }

    parsed = removeNulls(parsed);
    if (!Object.keys(parsed).length) {
      let type = '';
      if (resolver.restResolver) {
        type = 'REST';
      } else if (resolver.grpcResolver) {
        type = 'gRPC';
      }
      if (type) {
        return getDefaultConfigFromType(type as any);
      }
    }

    YAML.scalarOptions.null.nullStr = '';
    return YAML.stringify(parsed, { simpleKeys: true });
  }
  return '';
};
