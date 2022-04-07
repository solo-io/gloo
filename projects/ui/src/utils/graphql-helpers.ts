import {
  DocumentNode,
  EnumTypeDefinitionNode,
  FieldDefinitionNode,
  Kind,
  ObjectTypeDefinitionNode,
  parse,
} from 'graphql';
import lodash from 'lodash';
import { ClusterObjectRef } from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';
import { Resolution } from 'proto/github.com/solo-io/solo-apis/api/gloo/graphql.gloo/v1beta1/graphql_pb';
import { ResourceRef } from 'proto/github.com/solo-io/solo-kit/api/v1/ref_pb';
import { GraphqlApi } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/graphql_pb';

export function isElementInView(el: HTMLElement | null) {
  if (!el) return false;
  var rect = el.getBoundingClientRect();
  return (
    rect.top >= 0 &&
    rect.left >= 0 &&
    rect.bottom <=
      (window.innerHeight || document.documentElement.clientHeight) &&
    rect.right <= (window.innerWidth || document.documentElement.clientWidth)
  );
}

export const makeSchemaDefinitionId = (
  apiRef: ClusterObjectRef.AsObject,
  d: { name: { value: string } }
) => `${apiRef.namespace}-${apiRef.name}-${d.name.value.replace(/-|\s/g, '_')}`;

export const isExecutableAPI = (graphqlApi: GraphqlApi.AsObject) =>
  !!graphqlApi.spec?.executableSchema || !!(graphqlApi as any)?.executable;
export const isStitchedAPI = (graphqlApi: GraphqlApi.AsObject) =>
  !!graphqlApi.spec?.stitchedSchema || !!(graphqlApi as any)?.stitched;

export const makeGraphqlApiLink = (
  apiName?: string,
  apiNamespace?: string,
  apiCluster?: string,
  glooName?: string,
  glooNamespace?: string,
  isGlooFedEnabled?: boolean
) => {
  return !!isGlooFedEnabled
    ? `/gloo-instances/${glooNamespace ?? ''}/${glooName ?? ''}/apis/${
        apiCluster ?? ''
      }/${apiNamespace ?? ''}/${apiName ?? ''}/`
    : `/gloo-instances/${glooNamespace ?? ''}/${glooName ?? ''}/apis/${
        apiNamespace ?? ''
      }/${apiName ?? ''}/`;
};
// -------------------------------------------------------------------- //
// -------------------------------------------------------------------- //
// -------------------------------------------------------------------- //

export const makeGraphqlApiRef = (api: GraphqlApi.AsObject) => ({
  name: api?.metadata?.name ?? '',
  namespace: api?.metadata?.namespace ?? '',
  clusterName: api?.metadata?.clusterName ?? '',
});

export type SupportedDefinitionNode =
  | ObjectTypeDefinitionNode
  | EnumTypeDefinitionNode;
export interface SupportedDocumentNode extends DocumentNode {
  definitions: SupportedDefinitionNode[];
}
/**
 *
 * @param api
 * @returns The parsed schema with only the supported object definitions
 * (enum type or object type), sorted in the order: query, mutation, everything else.
 */
export const getParsedExecutableApiSchema = (
  api: GraphqlApi.AsObject | undefined
) => parseSchemaString(api?.spec?.executableSchema?.schemaDefinition);
/**
 *
 * @param schemaString
 * @returns The parsed schema with only the supported object definitions
 * (enum type or object type), sorted in the order: query, mutation, everything else.
 */
export const parseSchemaString = (schemaString: string | undefined) => {
  const emptySchema = {
    kind: Kind.DOCUMENT,
    definitions: [],
  } as SupportedDocumentNode;
  if (!schemaString) return emptySchema;
  try {
    // Parse and return it.
    const parsedSchema = parse(schemaString) as SupportedDocumentNode;
    // We support enum and object type definitions here.
    const definitions = lodash.cloneDeep(
      parsedSchema.definitions.filter(
        d =>
          d.kind === Kind.ENUM_TYPE_DEFINITION ||
          d.kind === Kind.OBJECT_TYPE_DEFINITION
      )
    ) as SupportedDefinitionNode[];
    // ? Uncomment this push(...mockEnumDefinitions) line for testing enums:
    // definitions.push(...mockEnumDefinitions);
    // We can sort the definitions here, and any filtering will keep it sorted.
    definitions.sort((a, b) => {
      // Ordering: Query, mutation, Everything else.
      if (a.name.value === 'Query') return -1;
      else if (b.name.value === 'Query') return 1;
      if (a.name.value === 'Mutation') return -1;
      else if (b.name.value === 'Mutation') return 1;
      else return 0;
    });
    return { ...parsedSchema, definitions };
  } catch (_) {
    return emptySchema;
  }
};
// -------------------------------------------------------------------- //
// -------------------------------------------------------------------- //
// -------------------------------------------------------------------- //

/**
 * Other useful functions: `getResolutionName(field)` or
 * `getResolutionNameFromDirective(directive)`.
 * @param field
 * @returns The resolve directive object for a field, or undefined if it does not exist.
 */
export const getResolveDirective = (
  field: FieldDefinitionNode | undefined | null
) =>
  field?.directives?.find(
    d => d.kind === Kind.DIRECTIVE && d.name.value === 'resolve'
  );

/**
 *
 * @param field
 * @returns The resolve directive's name argument for a field
 * (e.g. a field with `@resolve(name: YourResolution)` would return `YourResolution`).
 */
export const getResolveDirectiveName = (
  field: FieldDefinitionNode | undefined | null
) => {
  //
  // Gets the resolve directive for this field.
  const resolveDirective = getResolveDirective(field);
  if (!resolveDirective) return undefined;
  //
  // Get the resolve directive's 'name' argument.
  // '@resolve(name: "...")'
  const resolveDirectiveArg = resolveDirective.arguments?.find(
    a => a.name.value === 'name'
  );
  if (!resolveDirectiveArg || resolveDirectiveArg.value.kind !== Kind.STRING)
    return undefined;
  return resolveDirectiveArg.value.value;
};

/**
 *
 * @param api
 * @param resolutionName
 * @returns The resolution object with the given resolutionName in the api spec.
 */
export const getResolution = (
  api: GraphqlApi.AsObject | undefined,
  resolutionName?: string
) => {
  if (!resolutionName || !api) return undefined;
  const resolutionMapItem =
    api?.spec?.executableSchema?.executor?.local?.resolutionsMap?.find(
      ([rN, _]) => rN === resolutionName
    );
  if (!resolutionMapItem) return undefined;
  return resolutionMapItem[1];
};

/**
 *
 * @param api
 * @param resolutionName
 * @returns Whether a resolution with the given resolutionName exists in the api spec.
 */
export const hasResolutionWithName = (
  api: GraphqlApi.AsObject | undefined,
  resolutionName: string | undefined
) => getResolution(api, resolutionName) !== undefined;

/**
 *
 * @param api
 * @param field
 * @returns Whether a field has a resolver directive that maps to a resolution in the api spec.
 */
export const hasResolutionForField = (
  api: GraphqlApi.AsObject | undefined,
  field: FieldDefinitionNode
) => getResolution(api, getResolveDirectiveName(field)) !== undefined;
// -------------------------------------------------------------------- //
// -------------------------------------------------------------------- //
// -------------------------------------------------------------------- //

/**
 *
 * @param upstreamRef
 * @returns A string id for the upstream ('name::namespace').
 */
export const getUpstreamId = (upstreamRef: ResourceRef.AsObject | undefined) =>
  upstreamRef ? `${upstreamRef?.name}::${upstreamRef?.namespace}` : '';

/**
 *
 * @param upstreamId
 * @returns An upstream ref object from the id value ('name::namespace').
 */
export const getUpstreamRefFromId = (upstreamId: string) => {
  try {
    const [name, namespace] = upstreamId.split('::');
    return { name, namespace } as ResourceRef.AsObject;
  } catch (err) {
    return { name: '', namespace: '' } as ResourceRef.AsObject;
  }
};

/**
 *
 * @param resolution
 * @returns The upstream ref object.
 */
export const getUpstreamRef = (resolution: Resolution.AsObject | undefined) =>
  resolution?.restResolver?.upstreamRef?.name
    ? resolution?.restResolver?.upstreamRef
    : resolution?.grpcResolver?.upstreamRef;
// -------------------------------------------------------------------- //
// -------------------------------------------------------------------- //
// -------------------------------------------------------------------- //

/**
 * Traverses the field definition node to build the string representation of its return type.
 * @returns [prefix, base-type, suffix]
 */
export const getFieldReturnType = (
  field: FieldDefinitionNode | undefined | null
) => {
  const emptyType = {
    fullType: '',
    parts: {
      prefix: '',
      base: '',
      suffix: '',
    },
  };
  if (!field) return emptyType;
  let typePrefix = '';
  let typeSuffix = '';
  let typeBaseObj = field.type;
  // The fieldDefinition could be nested.
  while (true) {
    if (typeBaseObj?.kind === Kind.NON_NULL_TYPE) {
      typeSuffix = '!' + typeSuffix;
    } else if (typeBaseObj?.kind === Kind.LIST_TYPE) {
      typePrefix = typePrefix + '[';
      typeSuffix = ']' + typeSuffix;
    } else break;
    typeBaseObj = typeBaseObj.type;
  }
  if (typeBaseObj.kind === Kind.NAMED_TYPE)
    return {
      fullType: typePrefix + typeBaseObj.name.value + typeSuffix,
      parts: {
        prefix: typePrefix,
        base: typeBaseObj.name.value,
        suffix: typeSuffix,
      },
    };
  else return emptyType;
};

export const objectToArrayMap = (map: { [key: string]: any }) =>
  Object.keys(map).map(k => [k, map[k]] as [string, string]);
export const arrayMapToObject = <T>(map: [keyof T, any][]) => {
  const objMap = {} as T;
  map.forEach(pair => (objMap[pair[0]] = pair[1]));
  return objMap;
};
