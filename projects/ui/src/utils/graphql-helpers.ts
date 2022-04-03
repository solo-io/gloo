import {
  DocumentNode,
  EnumTypeDefinitionNode,
  Kind,
  ObjectTypeDefinitionNode,
} from 'graphql';
import gql from 'graphql-tag';
import lodash from 'lodash';
import { ClusterObjectRef } from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';
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
  !!graphqlApi.spec?.executableSchema;

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

export type supportedDefinitionTypes =
  | ObjectTypeDefinitionNode
  | EnumTypeDefinitionNode;

export const parseSchema = (schemaDefinition?: string) => {
  const emptyDoc = {
    kind: Kind.DOCUMENT,
    definitions: [] as supportedDefinitionTypes[],
  };
  if (!schemaDefinition) return emptyDoc;
  // Try to parse the serialized GraphQL schema definition to JSON (using gql`...`).
  let query: DocumentNode;
  try {
    query = gql`
      ${schemaDefinition}
    `;
  } catch {
    return emptyDoc;
  }
  if (!query) return emptyDoc;
  // We support enum and object type definitions here.
  const definitions = lodash.cloneDeep(
    query.definitions.filter(
      d =>
        d.kind === Kind.ENUM_TYPE_DEFINITION ||
        d.kind === Kind.OBJECT_TYPE_DEFINITION
    )
  ) as supportedDefinitionTypes[];
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
  return { kind: Kind.DOCUMENT, definitions };
};

export const objectToArrayMap = (map: { [key: string]: any }) =>
  Object.keys(map).map(k => [k, map[k]] as [string, string]);
export const arrayMapToObject = <T>(map: [keyof T, any][]) => {
  const objMap = {} as T;
  map.forEach(pair => (objMap[pair[0]] = pair[1]));
  return objMap;
};
