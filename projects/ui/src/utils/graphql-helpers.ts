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
