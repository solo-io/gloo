import { ResolverItem } from 'API/graphql';
import GQLJsonDescriptor from 'Components/Features/Graphql/data/graphql.json';
import { FieldDefinitionNode } from 'graphql';
import cloneDeep from 'lodash/cloneDeep';
import { Resolution } from 'proto/github.com/solo-io/solo-apis/api/gloo/graphql.gloo/v1beta1/graphql_pb';
import protobuf from 'protobufjs';
import { arrayMapToObject, objectToArrayMap } from 'utils/graphql-helpers';
import YAML from 'yaml';
import { getDefaultConfigFromType } from './ResolverConfigSection';
import { ResolverType } from './ResolverWizard';

/**
 * This is the root of the generated GraphQL protobuffer type descriptor.
 */
const jsonRoot = protobuf.Root.fromJSON(GQLJsonDescriptor);

export function removeNulls(obj: any) {
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
}

/**
 * Returns type information about which resolver to use.
 * @param resolver
 */
function getResolverTypeDataFromResolution(
  resolver: Resolution.AsObject,
  resolverTypeOverride?: ResolverType
) {
  let resolverPBType: string | undefined;
  let resolverKey: keyof typeof resolver | undefined;
  let resolverType: ResolverType | undefined;
  if (!!resolver?.grpcResolver || resolverTypeOverride === 'gRPC') {
    resolverPBType = 'graphql.gloo.solo.io.GrpcResolver';
    resolverKey = 'grpcResolver';
    resolverType = 'gRPC';
  } else if (!!resolver?.mockResolver || resolverTypeOverride === 'Mock') {
    resolverPBType = 'graphql.gloo.solo.io.MockResolver';
    resolverKey = 'mockResolver';
    resolverType = 'Mock';
  } else if (!!resolver?.restResolver || resolverTypeOverride === 'REST') {
    resolverPBType = 'graphql.gloo.solo.io.RESTResolver';
    resolverKey = 'restResolver';
    resolverType = 'REST';
  }
  return { resolverPBType, resolverKey, resolverType };
}

/**
 * This is called when creating the parameters to UPDATE the api.
 * @param resolverConfig
 * @param resolverType
 * @param field
 * @param upstream
 * @param extras
 * @returns
 */
export function createResolverItem(
  resolverConfig: string,
  resolverType: ResolverType,
  field: FieldDefinitionNode,
  upstream: string,
  extras: Record<string, any> = {}
) {
  //
  // Parses the resolver config YAML string.
  let resolver: any;
  try {
    YAML.scalarOptions.null.nullStr = '';
    resolver = removeNulls(YAML.parse(resolverConfig));
  } catch (err: any) {
    throw err;
  }
  let parsedResolver = {} as any;
  let pbTypeObj: protobuf.Type;
  //
  // Gets type information for the resolver.
  const { resolverKey, resolverPBType } = getResolverTypeDataFromResolution(
    resolver,
    resolverType
  );
  if (!!resolverKey && !!resolverPBType) {
    //
    // e.g. If `resolver === {restResolver: {...}}`,
    // this sets: `resolver = resolver.restResolver;`
    if (resolver[resolverKey]) resolver = resolver[resolverKey];
    //
    // Set default values if appropriate.
    // e.g. If a mockResolver doesn't have the response type specified,
    // this defaults it to syncResponse.
    if (
      resolverType === 'Mock' &&
      !resolver?.syncResponse &&
      !resolver?.errorResponse &&
      !resolver?.asyncResponse
    )
      resolver = { syncResponse: resolver };
    if (resolver.upstreamRef) delete resolver.upstreamRef;
    //
    // Converts using the generated proto type descriptor.
    pbTypeObj = jsonRoot.lookupType(resolverPBType);
    parsedResolver = preMarshallProtoValues(resolver, pbTypeObj.toJSON());
  }
  if (
    !parsedResolver ||
    typeof parsedResolver !== 'object' ||
    Object.keys(parsedResolver).length === 0
  ) {
    let errorMsg = 'Invalid configuration.';
    if (!!pbTypeObj!)
      errorMsg +=
        ' Start with these root properties: "' +
        Object.keys(pbTypeObj.fields)
          // TODO: Removing upstreamRef is kind of a hack. This should be passed as a separate object and not saved in the string.
          .filter(f => f !== 'upstreamRef')
          .join('" "') +
        '"';
    throw new Error(errorMsg);
  }
  //
  // Gets the upstream dropdown value.
  let [upstreamName, upstreamNamespace] = upstream.split('::');
  const upstreamRef = {
    name: upstreamName,
    namespace: upstreamNamespace,
  };
  return {
    // TODO: Could improve the ResolverItem type so the resolution objects are passed more explicitly.
    // This might look like: { RESTResolver?: {...} | GrpcResolver?: {...} | MockResolver?: {...} }
    // Or it could be simplified. But at this point, the parsedResolver is converted
    // into a RESTResolver.AsObject | Grpc.AsObject | MockResolver.AsObject for use in API/graphql.ts.
    ...(resolverType === 'REST' && parsedResolver),
    ...(resolverType === 'gRPC' && {
      grpcRequest: parsedResolver.requestTransform,
      spanName: parsedResolver.spanName,
    }),
    ...(resolverType === 'Mock' && { mockResolver: parsedResolver }),
    upstreamRef,
    field,
    resolverType,
    ...extras,
  } as ResolverItem;
}

/**
 * This is called to parse the RETURNED api.
 */
export function getResolverFromConfig(resolver?: Resolution.AsObject) {
  //
  // Defaults to REST if the resolver is invalid or undefined.
  if (!resolver) return getDefaultConfigFromType('REST');
  //
  // Gets type information for the resolver.
  const { resolverKey, resolverPBType, resolverType } =
    getResolverTypeDataFromResolution(resolver);
  if (!resolverType) return getDefaultConfigFromType('REST');
  if (!resolverKey || !resolverPBType)
    return getDefaultConfigFromType(resolverType);
  let clonedResolver = cloneDeep(resolver[resolverKey]) as any;
  clonedResolver = removeNulls(clonedResolver);
  if (!clonedResolver || Object.keys(clonedResolver).length === 0)
    return getDefaultConfigFromType(resolverType);
  //
  // Remove upstream, so it doesn't show up in the editor.
  // We could do this for any other values that have a separate form control,
  // or keep upstream in for a full edit mode.
  if (clonedResolver.upstreamRef) delete clonedResolver.upstreamRef;
  //
  // Parses the object using the generated proto JSON descriptor.
  const pbTypeObj = jsonRoot.lookupType(resolverPBType);
  let parsedResolver: any;
  try {
    parsedResolver = postUnmarshallProtoValues(
      clonedResolver,
      pbTypeObj.toJSON()
    );
  } catch (e) {}
  //
  // Returns the stringified result.
  if (!parsedResolver || Object.keys(parsedResolver).length === 0)
    return getDefaultConfigFromType(resolverType);
  YAML.scalarOptions.null.nullStr = '';
  return YAML.stringify(parsedResolver, { simpleKeys: true });
}

export function isBase64(str: string) {
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
}

/**
 * @example
 * stringToBase64("foo"); // returns 'Zm9v'
 */
export function stringToBase64(str: string) {
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
}
/**
 * @example
 * base64ToString("Zm9v") // returns 'foo'
 */
export function base64ToString(str: string) {
  return decodeURIComponent(atob(str));
}

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
function postUnmarshallProtoValues(
  obj: any,
  pbType: protobuf.IType,
  path = ''
) {
  const primitive = tryGetPrimitiveWithPbCheck(obj, pbType, path);
  if (!!primitive) return primitive.value;
  //
  // -- RECURSE THROUGH `obj` -- //
  const parsedObj = {} as any;
  const { fields } = pbType as protobuf.IType;
  const pbFieldKeys = Object.keys(fields);
  Object.keys(obj).forEach(k => {
    // If this field is not in the proto type, ignore it.
    let pbFieldKey = pbFieldKeys.find(pbK => pbK === k);
    // The field key might have changed to <fieldName> + "Map", so we can check for that.
    if (
      pbFieldKey === undefined &&
      k.substring(k.length - 3, k.length) === 'Map'
    ) {
      // Removes "Map" from the end of the key.
      k = k.substring(0, k.length - 3);
      pbFieldKey = pbFieldKeys.find(pbK => pbK === k);
    }
    if (pbFieldKey === undefined) {
      throwPbObjectKeyError(pbFieldKeys, k, path);
      return;
    }
    // ====================== //
    // SPECIAL PARSING CASES: //
    // ====================== //
    const fieldTypeName = fields[pbFieldKey].type;
    if (fieldTypeName === 'google.protobuf.Value') {
      // TODO: When the apiserver returns oneof nullValue, numberValue, boolValue, etc,
      // Then we can use pb-utils value.encode and decode. Unitl then we have to use
      // custom parsing to try and find the truthy value if it exists.
      parsedObj[k] = convertFromPBValue(obj[k]);
      // parsedObj[k] = value.decode(obj[k]);
      return;
    }
    // ====================== //
    // DEFAULT TYPE HANDLING: //
    // ====================== //
    const newPath = path === '' ? k : `${path}.${k}`;
    try {
      const nestedType = jsonRoot.lookupType(fieldTypeName) as protobuf.IType;
      parsedObj[k] = postUnmarshallProtoValues(obj[k], nestedType, newPath);
    } catch (e: any) {
      if (e?.message?.split(':')[0] === 'Parsing Error') throw e;
      if (!!(fields[pbFieldKey] as protobuf.IMapField).keyType) {
        // "Map" is added to the end of map field keys.
        parsedObj[k] = arrayMapToObject(obj[`${k}Map`]);
      } else {
        const primitive = tryGetPrimitiveWithPbCheck(
          // obj[`${k}Map`],
          obj[k],
          fields[pbFieldKey] as any,
          newPath
        );
        if (!!primitive) parsedObj[k] = primitive.value;
      }
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
function preMarshallProtoValues(obj: any, pbType: protobuf.IType, path = '') {
  const primitive = tryGetPrimitiveWithPbCheck(obj, pbType, path);
  if (!!primitive) return primitive.value;
  //
  // -- RECURSE THROUGH `obj` -- //
  const parsedObj = {} as any;
  const { fields } = pbType as protobuf.IType;
  const pbFieldKeys = Object.keys(fields);
  Object.keys(obj).forEach(k => {
    // If this field is not in the proto type, ignore it.
    const pbFieldKey = pbFieldKeys.find(pbK => pbK === k);
    if (pbFieldKey === undefined) {
      throwPbObjectKeyError(pbFieldKeys, k, path);
      return;
    }
    // ====================== //
    // SPECIAL PARSING CASES: //
    // ====================== //
    const fieldTypeName = fields[pbFieldKey].type;
    if (fieldTypeName === 'google.protobuf.Value') {
      // parsedObj[k] = value.encode(obj[k]);
      parsedObj[k] = convertToPBValue(obj[k]);
      return;
    }
    // ====================== //
    // DEFAULT TYPE HANDLING: //
    // ====================== //
    const newPath = path === '' ? k : `${path}.${k}`;
    try {
      const nestedType = jsonRoot.lookupType(fieldTypeName) as protobuf.IType;
      parsedObj[k] = preMarshallProtoValues(obj[k], nestedType, newPath);
    } catch (e: any) {
      if (e?.message?.split(':')[0] === 'Parsing Error') throw e;
      if (
        !!(fields[pbFieldKey] as protobuf.IMapField).keyType &&
        typeof obj[k] === 'object'
      ) {
        // "Map" is added to the end of map field keys.
        const keyName = `${k}Map`;
        parsedObj[keyName] = objectToArrayMapWithPbCheck(
          obj[k],
          fields[pbFieldKey],
          newPath
        );
      } else {
        const primitive = tryGetPrimitiveWithPbCheck(
          obj[k],
          fields[pbFieldKey],
          newPath
        );
        if (!!primitive) parsedObj[k] = primitive.value;
      }
    }
  });
  return parsedObj;
}

function objectToArrayMapWithPbCheck(obj: any, pbType: any, path: string) {
  const nonStringKeys = Object.keys(obj).filter(
    k => typeof obj[k] !== 'string'
  );
  if (pbType?.type === 'string' && nonStringKeys.length > 0)
    throw new Error(
      `Parsing Error: "${nonStringKeys.join('", "')}"${
        path && ' at "' + path + '"'
      } should be ${nonStringKeys.length > 1 ? 'strings' : 'a string'}.`
    );
  return objectToArrayMap(obj);
}

function tryGetPrimitiveWithPbCheck(obj: any, pbType: any, path: string) {
  if (
    obj === null ||
    typeof obj !== 'object' ||
    Object.keys(obj).length === 0
  ) {
    if (!!pbType?.fields)
      throwPbObjectValueError(Object.keys(pbType.fields), obj, path);
    if (!!pbType?.keyType) throwPbMapError(obj, path);
    if (pbType?.type === 'string' && typeof obj !== 'string')
      throw new Error(
        `Parsing Error: "${obj}"${path && ' at ' + path} should be a string.`
      );
    return { value: obj };
  }
  return null;
}

function getValidKeys(pbFieldKeys: string[]) {
  return (
    '"' +
    pbFieldKeys
      // TODO: Removing upstreamRef is kind of a hack. This should be passed as a separate object and not saved in the string.
      .filter(f => f !== 'upstreamRef')
      .join('", "') +
    '"'
  );
}

function throwPbObjectValueError(
  pbFieldKeys: string[],
  k: string,
  path: string
) {
  let errorMessage = 'Parsing Error: ';
  const validKeys = getValidKeys(pbFieldKeys);
  if (!!path)
    errorMessage += `"${path}" is an object, so it cannot be "${k}". Its properties include: ${validKeys}`;
  else
    errorMessage += `"${k}" is not an object. This configuration object must include the properties: ${validKeys}`;
  throw new Error(errorMessage);
}

function throwPbObjectKeyError(pbFieldKeys: string[], k: string, path: string) {
  let errorMessage = 'Parsing Error: ';
  if (pbFieldKeys.length === 0) errorMessage += `"${path}" is not an object.`;
  else {
    const validKeys = getValidKeys(pbFieldKeys);
    if (!!path)
      errorMessage += `"${k}" is not a property of "${path}". Valid properties include: ${validKeys}`;
    else
      errorMessage += `"${k}" is not a valid property. Valid properties include: ${validKeys}`;
  }
  throw new Error(errorMessage);
}

function throwPbMapError(k: string, path: string) {
  let errorMessage = 'Parsing Error: ';
  if (!!path) errorMessage += `"${path}" is an object, so it cannot be "${k}".`;
  else errorMessage += `"${k}" is not an object.`;
  throw new Error(errorMessage);
}

/**
 * This is similar to:
 * ```ts
 * import {value} from 'pb-utils';
 * ...
 * value.decode(obj)
 * ```
 * But `value.decode(obj)` will return `nullValue`, `numberValue`, and
 * `boolValue` since there is no way to currently distinguish between
 * 0, false, and NULL. So this function tries to find the truthy value
 * before parsing it, and otherwise returns 0. Also there are some naming
 * differences between the decoded value and the AsObject that we need. This
 * may be able to be fixed with some pb-utils options.
 * @param obj
 * @returns
 */
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
    // See function description for more info.
    const truthyKey = Object.keys(obj).find(k => !!obj[k]);
    if (truthyKey !== undefined) return obj[truthyKey];
    else return 0;
  }
}

/**
 * This is similar to:
 * ```ts
 * import {value} from 'pb-utils';
 * ...
 * value.encode(obj)
 * ```
 * But there are some naming differences between the encoded value
 * and the AsObject that we need. This may be able to be fixed with
 * some pb-utils options.
 * For example: An encoded "structValue" in our proto type has a "fieldsMap" property,
 * and the pb-utils encoded structValue has "fields" property.
 * For lists, we use `{"listValues": {"valuesList": [...] }}` and pb-utils uses
 * `{"listValues": {"values": [...]}}`. This conversion matters down the line when we're piecing the
 * protobufs together in the API/graphql.ts file.
 * @param obj
 * @returns
 */
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
  return null;
}
