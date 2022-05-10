import { value } from 'pb-util';
import { ResolverItem } from 'API/graphql';
import GQLJsonDescriptor from 'Components/Features/Graphql/data/graphql.json';
import { FieldDefinitionNode } from 'graphql';
import cloneDeep from 'lodash/cloneDeep';
import { Resolution } from 'proto/github.com/solo-io/solo-apis/api/gloo/graphql.gloo/v1beta1/graphql_pb';
import protobuf from 'protobufjs';
import { arrayMapToObject, objectToArrayMap } from 'utils/graphql-helpers';
import YAML from 'yaml';
import { getDefaultConfigFromType } from './ResolverConfigSection';
import { ResolverType, ResolverWizardFormProps } from './ResolverWizard';

const jsonRoot = protobuf.Root.fromJSON(GQLJsonDescriptor);

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

/**
 * This is called when creating the parameters to UPDATE the api.
 * @param resolverConfig
 * @param resolverType
 * @param field
 * @param upstream
 * @param extras
 * @returns
 */
export const createResolverItem = (
  resolverConfig: string,
  resolverType: ResolverType,
  field: FieldDefinitionNode,
  upstream: string,
  extras: Record<string, any> = {}
): ResolverItem => {
  /*
     `parsedResolverConfig` can be formatted in different ways:
     - `restResolver.[request | response | spanName | ...]`....
     - `grpcResolver.[request | response | spanName | ...]`...
     - `mockResolver.[ syncResponse | ...]`...
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
      headersMap = objectToArrayMap(request?.headers ?? {});
      queryParamsMap = objectToArrayMap(request?.queryParams ?? {});
      body = request?.body;
    }
    if (response) {
      settersMap = objectToArrayMap(
        response?.settersMap ?? response?.setters ?? {}
      );
      resultRoot = response?.resultRoot;
    }
  } else {
    if (resolverType === 'gRPC' && requestTransform) {
      requestMetadataMap = objectToArrayMap(
        requestTransform?.requestMetadataMap ?? {}
      );
      serviceName = requestTransform?.serviceName;
      methodName = requestTransform?.methodName;
      outgoingMessageJson = requestTransform?.outgoingMessageJson;
    }
  }
  let [upstreamName, upstreamNamespace] = upstream.split('::');

  const mockPbType = jsonRoot.lookupType('graphql.gloo.solo.io.MockResolver');
  // const restPbType = jsonRoot.lookupType('graphql.gloo.solo.io.RESTResolver');
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
    mockResolver:
      parsedResolverConfig !== null && typeof parsedResolverConfig === 'object'
        ? preMarshallProtoValues(parsedResolverConfig, mockPbType.toJSON())
        : preMarshallProtoValues(
            { syncResponse: parsedResolverConfig },
            mockPbType.toJSON()
          ),
    resolverType: resolverType,
    // @ts-ignore
    ...(response && { response: { resultRoot, settersMap } }),
    spanName,

    //TODO: Get this working for rest, grpc resolvers.
    // ...(parsedResolverConfig !== null &&
    // typeof parsedResolverConfig === 'object'
    //   ? preMarshallProtoValues(parsedResolverConfig, restPbType.toJSON())
    //   : {}),

    ...extras,
  } as ResolverItem;
  console.log(resolverItem.mockResolver.syncResponse);
  return resolverItem;
};

/**
 * This is called to parse the RETURNED api.
 */
export const getResolverFromConfig = (resolver?: Resolution.AsObject) => {
  if (resolver?.restResolver || resolver?.grpcResolver) {
    // TODO:  This conversion doesn't quite work.
    // conversion: Resolution.AsObject -> protobufjs.Message -> proto3 JsonValue -> string
    let parsed: Record<string, any> = cloneDeep(resolver);
    if (!!parsed.restResolver) {
      parsed = cloneDeep(parsed.restResolver);
    } else if (!!parsed.grpcResolver) {
      parsed = cloneDeep(parsed.grpcResolver);
    }
    delete parsed.upstreamRef;

    if (parsed?.request?.headersMap) {
      parsed.request.headers = arrayMapToObject(parsed.request.headersMap);
      if (!parsed.request.headersMap.length) {
        delete parsed.request.headers;
      }
      delete parsed.request.headersMap;
    }

    if (parsed?.request?.queryParamsMap) {
      parsed.request.queryParams = arrayMapToObject(
        parsed.request.queryParamsMap
      );
      if (!parsed.request.queryParamsMap.length) {
        delete parsed.request.queryParams;
      }
      delete parsed.request.queryParamsMap;
    }
    if (parsed?.response?.settersMap) {
      parsed.response.setters = arrayMapToObject(parsed.response.settersMap);
      if (!parsed.response.settersMap.length) {
        delete parsed.response.setters;
      }
      delete parsed.response.settersMap;
    }

    if (parsed?.requestTransform?.requestMetadataMap) {
      parsed.requestTransform.requestMetadata = arrayMapToObject(
        parsed.requestTransform.requestMetadataMap
      );
      if (!parsed.requestTransform.requestMetadataMap.length) {
        delete parsed.requestTransform.requestMetadata;
      }
      delete parsed.requestTransform.requestMetadataMap;
    }

    parsed = removeNulls(parsed);
    if (!Object.keys(parsed).length) {
      let type = 'REST' as ResolverWizardFormProps['resolverType'];
      if (resolver.grpcResolver) type = 'gRPC';
      else if (resolver.mockResolver) type = 'Mock';
      return getDefaultConfigFromType(type);
    }

    YAML.scalarOptions.null.nullStr = '';
    return YAML.stringify(parsed, { simpleKeys: true });
  }

  if (resolver?.mockResolver) {
    let parsed = cloneDeep(resolver.mockResolver) as any;
    parsed = removeNulls(parsed);
    //
    // If not set, uses the default config.
    if (!parsed || !Object.keys(parsed).length)
      return getDefaultConfigFromType('Mock');
    //
    // Parses the object using the generated proto JSON descriptor.
    const pbTypeObj = jsonRoot.lookupType('graphql.gloo.solo.io.MockResolver');
    parsed = postUnmarshallProtoValues(parsed, pbTypeObj.toJSON());
    //
    // Returns the stringified result.
    YAML.scalarOptions.null.nullStr = '';
    return YAML.stringify(parsed, { simpleKeys: true });
  }
  return getDefaultConfigFromType('REST');
};

export const isBase64 = (str: string) => {
  // Technically, this is valid base64.
  // @see https://datatracker.ietf.org/doc/html/rfc4648#section-10
  if (!str || str.trim() === '') {
    return false;
  }
  try {
    return btoa(atob(str)) === str;
  } catch (err) {
    return false;
  }
};

/**
 * @example
 * stringToBase64("foo"); // returns 'Zm9v'
 */
export const stringToBase64 = (str: string) => {
  // First we escape the string using encodeURIComponent to get the UTF-8 encoding of the characters,
  // then we convert the percent encodings into raw bytes, and finally feed it to btoa() function.
  const utf8Bytes = encodeURIComponent(str).replace(
    /%([0-9A-F]{2})/g,
    (_match, p1) => {
      // @ts-ignore
      return String.fromCharCode('0x' + p1);
    }
  );
  return btoa(utf8Bytes);
};
/**
 * @example
 * base64ToString("Zm9v") // returns 'foo'
 */
export const base64ToString = (str: string) => {
  return decodeURIComponent(atob(str));
};

/**
 * When GETTING the values.
 * Recursively parses the resolver config value returned from the API,
 * using some special parsing cases to improve editor UX while using protobuffers.
 * For example, it formats google Value protobufs to be more readable.
 * Works with `preMarshallProtoValues(...)`.
 * @param obj The resolver config object returned from the API.
 * @param pbType The protobuffer type, created from the generated file:
 * `'Components/Features/Graphql/data/graphql.json'`
 * @returns A deep copy of obj, formatted for the resolver wizard config.
 */
function postUnmarshallProtoValues(obj: any, pbType: protobuf.IType) {
  if (obj === null || typeof obj !== 'object' || Object.keys(obj).length === 0)
    return cloneDeep(obj);
  //
  // -- RECURSE THROUGH `obj` -- //
  const parsedObj = {} as any;
  const { fields } = pbType;
  const pbFieldKeys = Object.keys(fields);
  Object.keys(obj).forEach(k => {
    // If this field is not in the proto type, ignore it.
    const pbFieldKey = pbFieldKeys.find(pbK => pbK === k);
    if (pbFieldKey === undefined) return;
    // ====================== //
    // SPECIAL PARSING CASES: //
    // ====================== //
    const fieldTypeName = fields[pbFieldKey].type;
    if (fieldTypeName === 'google.protobuf.Value') {
      // TODO: When the apiserver returns oneof nullValue, numberValue, boolValue, etc,
      // Then we can use pb-utils value.encode and decode. Unitl then we have to use
      // custom parsing to try and find the truthy value if it exists.
      // parsedObj[k] = obj[k];
      parsedObj[k] = convertFromPBValue(obj[k]);
      // parsedObj[k] = value.decode(obj[k]);
      return;
    }
    // ====================== //
    // DEFAULT TYPE HANDLING: //
    // ====================== //
    try {
      const nestedType = jsonRoot.lookupType(fieldTypeName) as protobuf.IType;
      parsedObj[k] = postUnmarshallProtoValues(obj[k], nestedType);
    } catch (e) {
      if (typeof fieldTypeName === 'string') parsedObj[k] = obj[k] + '';
      else if (typeof fieldTypeName === 'number') parsedObj[k] = Number(obj[k]);
      else if (typeof fieldTypeName === 'boolean')
        parsedObj[k] = Boolean(obj[k]);
      else parsedObj[k] = obj[k];
    }
  });
  return parsedObj;
}

/**
 * When SETTING the values.
 * Recursively processes the resolver config value using some
 * special parsing cases to improve editor UX while using protobuffers.
 * For example, it formats google Value protobufs to be more readable.
 * Works with `postUnmarshallProtoValues(...)`.
 * @param obj The resolver config object, parsed from the resolver config YAML string.
 * @param pbType The protobuffer type, created from the generated file:
 * `'Components/Features/Graphql/data/graphql.json'`
 * @returns A deep copy of obj, formatted to be marshalled and sent to the API.
 */
function preMarshallProtoValues(obj: any, pbType: protobuf.IType) {
  //
  // -- RECURSE THROUGH `obj` -- //
  const parsedObj = {} as any;
  const { fields } = pbType;
  const pbFieldKeys = Object.keys(fields);
  Object.keys(obj).forEach(k => {
    // If this field is not in the proto type, ignore it.
    const pbFieldKey = pbFieldKeys.find(pbK => pbK === k);
    if (pbFieldKey === undefined) return;
    // ====================== //
    // SPECIAL PARSING CASES: //
    // ====================== //
    const fieldTypeName = fields[pbFieldKey].type;
    if (fieldTypeName === 'google.protobuf.Value') {
      // TODO: When the apiserver returns oneof nullValue, numberValue, boolValue, etc,
      // Then we can use pb-utils value.encode and decode. Unitl then we have to use
      // custom parsing to try and find the truthy value if it exists.
      parsedObj[k] = convertToPBValue(obj[k]);
      // parsedObj[k] = value.encode(obj[k]);
      return;
    }
    // ====================== //
    // DEFAULT TYPE HANDLING: //
    // ====================== //
    try {
      const nestedType = jsonRoot.lookupType(fieldTypeName) as protobuf.IType;
      parsedObj[k] = preMarshallProtoValues(obj[k], nestedType);
    } catch (e) {
      if (typeof fieldTypeName === 'string') parsedObj[k] = obj[k] + '';
      else if (typeof fieldTypeName === 'number') parsedObj[k] = Number(obj[k]);
      else if (typeof fieldTypeName === 'boolean')
        parsedObj[k] = Boolean(obj[k]);
      else parsedObj[k] = obj[k];
    }
  });
  return parsedObj;
}

function convertFromPBValue(obj: any): any {
  if (!!obj.listValue?.valuesList) {
    return obj.listValue?.valuesList.map((value: any) =>
      convertFromPBValue(value)
    );
  } else if (!!obj.structValue?.fieldsMap) {
    return arrayMapToObject(
      (obj.structValue.fieldsMap as any[]).map(([key, value]) => [
        key,
        convertFromPBValue(value),
        // value,
      ])
    );
  } else {
    // TODO: There is no way to currently distinguish between 0, false, and NULL.
    // Setting 0 or false works fine. Setting "null" doesn't work.
    // So if there is a truthy key, take that one, otherwise return 0.
    const truthyKey = Object.keys(obj).find(k => !!obj[k]);
    if (truthyKey !== undefined) return obj[truthyKey];
    else return 0;
  }
}

function convertToPBValue(obj: any): any {
  if (obj === null) return { nullValue: 0 };
  if (typeof obj === 'string') return { stringValue: obj };
  if (typeof obj === 'number') return { numberValue: obj };
  if (typeof obj === 'boolean') return { boolValue: obj };
  if (Array.isArray(obj))
    return { listValue: { valuesList: obj.map(o => convertToPBValue(o)) } };
  if (typeof obj === 'object')
    return {
      structValue: {
        fieldsMap: objectToArrayMap(obj).map(([key, value]) => [key, value]),
      },
    };
  // if (typeof obj === 'object') return { structValue: obj };
  return null;
}
