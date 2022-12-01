import { useGetStitchedSchemaDefinition } from 'API/hooks';
import {
  buildSchema,
  DirectiveDefinitionNode,
  DocumentNode,
  EnumTypeDefinitionNode,
  FieldDefinitionNode,
  InputObjectTypeDefinitionNode,
  InterfaceTypeDefinitionNode,
  Kind,
  ObjectTypeDefinitionNode,
  OperationDefinitionNode,
  parse,
  TypeNode,
  UnionTypeDefinitionNode,
} from 'graphql';
import lodash from 'lodash';
import { ClusterObjectRef } from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';
import { Resolution } from 'proto/github.com/solo-io/solo-apis/api/gloo/graphql.gloo/v1beta1/graphql_pb';
import { ResourceRef } from 'proto/github.com/solo-io/solo-kit/api/v1/ref_pb';
import { GraphqlApi } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/graphql_pb';
import { useMemo } from 'react';

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
  d: SupportedDefinitionNode
) => {
  let name = d.kind + ':';
  if (d.name) name += d.name.value;
  return `${apiRef.namespace}-${apiRef.name}-${name.replace(/-|\s/g, '_')}`;
};

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
  | InterfaceTypeDefinitionNode
  | ObjectTypeDefinitionNode
  | EnumTypeDefinitionNode
  | InputObjectTypeDefinitionNode
  | UnionTypeDefinitionNode
  | OperationDefinitionNode
  | DirectiveDefinitionNode;
export interface SupportedDocumentNode extends DocumentNode {
  definitions: SupportedDefinitionNode[];
}

export const getKindTypeReadableName = (
  definitionNode: SupportedDefinitionNode
) => {
  switch (definitionNode.kind) {
    case Kind.OBJECT_TYPE_DEFINITION:
      // return 'Object';
      return 'type';
    case Kind.ENUM_TYPE_DEFINITION:
      return 'enum';
    case Kind.INTERFACE_TYPE_DEFINITION:
      return 'interface';
    case Kind.INPUT_OBJECT_TYPE_DEFINITION:
      return 'input';
    case Kind.DIRECTIVE_DEFINITION:
      return 'directive';
    case Kind.UNION_TYPE_DEFINITION:
      return 'union';
    case Kind.OPERATION_DEFINITION:
      // return lodash.capitalize(definitionNode.operation);
      return definitionNode.operation;
    default:
      return '';
  }
};

export const kindTypeSort = (
  a: SupportedDefinitionNode,
  b: SupportedDefinitionNode
) => {
  // - Objects
  const isAObjType = a.kind === Kind.OBJECT_TYPE_DEFINITION;
  const isBObjType = b.kind === Kind.OBJECT_TYPE_DEFINITION;
  if (isAObjType && a.name.value === 'Query') return -1;
  if (isBObjType && b.name.value === 'Query') return 1;
  if (isAObjType && a.name.value === 'Mutation') return -1;
  if (isBObjType && b.name.value === 'Mutation') return 1;
  if (isAObjType) return -1;
  if (isBObjType) return 1;
  // - Enums
  if (a.kind === Kind.ENUM_TYPE_DEFINITION) return -1;
  if (b.kind === Kind.ENUM_TYPE_DEFINITION) return 1;
  // - Interfaces
  if (a.kind === Kind.INTERFACE_TYPE_DEFINITION) return -1;
  if (b.kind === Kind.INTERFACE_TYPE_DEFINITION) return 1;
  // - Inputs
  if (a.kind === Kind.INPUT_OBJECT_TYPE_DEFINITION) return -1;
  if (b.kind === Kind.INPUT_OBJECT_TYPE_DEFINITION) return 1;
  // - Directives
  if (a.kind === Kind.DIRECTIVE_DEFINITION) return -1;
  if (b.kind === Kind.DIRECTIVE_DEFINITION) return 1;
  // - Unions
  if (a.kind === Kind.UNION_TYPE_DEFINITION) return -1;
  if (b.kind === Kind.UNION_TYPE_DEFINITION) return 1;
  // - Subscriptions
  if (a.kind === Kind.OPERATION_DEFINITION) return -1;
  if (b.kind === Kind.OPERATION_DEFINITION) return 1;
  else return 0;
};

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
          d.kind === Kind.INTERFACE_TYPE_DEFINITION ||
          d.kind === Kind.ENUM_TYPE_DEFINITION ||
          d.kind === Kind.OBJECT_TYPE_DEFINITION ||
          d.kind === Kind.INPUT_OBJECT_TYPE_DEFINITION ||
          d.kind === Kind.DIRECTIVE_DEFINITION ||
          d.kind === Kind.UNION_TYPE_DEFINITION ||
          d.kind === Kind.OPERATION_DEFINITION
      )
    ) as SupportedDefinitionNode[];
    // We can sort the definitions here, and any filtering will keep it sorted.
    definitions.sort(kindTypeSort);
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
  field: FieldDefinitionNode | undefined | null | TypeNode
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
  let typeBaseObj = 'type' in field ? field.type : (field as any);
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
  else if (
    typeBaseObj.kind === Kind.ARGUMENT &&
    typeBaseObj.value.kind === Kind.VARIABLE
  )
    return {
      fullType: '$' + typeBaseObj.value.name.value,
      parts: {
        prefix: '$',
        base: typeBaseObj.value.name.value,
        suffix: '',
      },
    };
  else return emptyType;
};

export const objectToArrayMap = (obj: { [key: string]: any }) =>
  Object.keys(obj).map(k => [k, obj[k]] as [string, string]);
export const arrayMapToObject = <T>(arr: [keyof T, any][]) => {
  const objMap = {} as T;
  arr.forEach(pair => (objMap[pair[0]] = pair[1]));
  return objMap;
};

/** This returns the `schemaText` and `parsedSchema` for any executable or stitched `graphqlApi`. */
export const useGetSchema = (graphqlApi: GraphqlApi.AsObject | undefined) => {
  // Try to get executable schema
  const executableSchemaText =
    graphqlApi?.spec?.executableSchema?.schemaDefinition ?? '';
  const executableSchema = useMemo(() => {
    try {
      return buildSchema(executableSchemaText, { assumeValidSDL: true });
    } catch {
      return undefined;
    }
  }, [executableSchemaText]);
  // Try to get stitched schema
  const { data: stitchedSchemaText } = useGetStitchedSchemaDefinition({
    name: graphqlApi?.metadata?.name ?? '',
    namespace: graphqlApi?.metadata?.namespace ?? '',
    clusterName: graphqlApi?.metadata?.clusterName ?? '',
  });
  const parsedStitchedSchema = useMemo(() => {
    try {
      return buildSchema(stitchedSchemaText ?? '', { assumeValidSDL: true });
    } catch {
      return undefined;
    }
  }, [stitchedSchemaText]);
  //
  // Return the correct schema
  if (!!graphqlApi && isStitchedAPI(graphqlApi)) {
    return {
      schemaText: stitchedSchemaText,
      parsedSchema: parsedStitchedSchema,
    };
  } else {
    return { schemaText: executableSchemaText, parsedSchema: executableSchema };
  }
};
