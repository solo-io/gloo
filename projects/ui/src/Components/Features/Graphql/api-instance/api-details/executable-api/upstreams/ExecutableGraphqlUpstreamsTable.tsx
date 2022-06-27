import { useGetGraphqlApiDetails, useListUpstreams } from 'API/hooks';
import { ClusterObjectRef } from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';
import { ResourceRef } from 'proto/github.com/solo-io/solo-kit/api/v1/ref_pb';
import React, { useMemo } from 'react';
import { di } from 'react-magnetic-di/macro';
import { useNavigate } from 'react-router';
import { Spacer } from 'Styles/StyledComponents/spacer';

const ExecutableGraphqlUpstreamsTable: React.FC<{
  apiRef: ClusterObjectRef.AsObject;
}> = ({ apiRef }) => {
  di(useNavigate, useGetGraphqlApiDetails, useListUpstreams);
  const navigate = useNavigate();
  const { data: graphqlApi } = useGetGraphqlApiDetails(apiRef);

  // The set of upstream resourceRefs that this GraphQL API's spec references.
  const resolverUpstreams = useMemo(() => {
    //
    // Gets the resolutions map from the spec.
    const resolutionsMap =
      graphqlApi?.spec?.executableSchema?.executor?.local?.resolutionsMap;
    if (!resolutionsMap) return [];
    //
    // Gets the resolver upstream references from the map.
    return resolutionsMap
      .map(([_resolveName, resolver]) => resolver.restResolver?.upstreamRef)
      .filter(upRef => upRef !== undefined)
      .reduce((accum, cur) => {
        if (cur === undefined) return accum;
        const existingUpstream = accum.find(
          u => u.name === cur.name && u.namespace === cur.namespace
        );
        if (!existingUpstream) accum.push(cur);
        return accum;
      }, [] as ResourceRef.AsObject[]);
  }, [graphqlApi]);

  return (
    <div className='mb-8'>
      {resolverUpstreams?.map(upstream => {
        const upNamespace = upstream.namespace ?? '';
        const upName = upstream.name ?? '';
        const upCluster = graphqlApi?.metadata?.clusterName ?? '';
        const glooInstNamespace = graphqlApi?.glooInstance?.namespace;
        const glooInstName = graphqlApi?.glooInstance?.name;
        const upstreamDetailsUrl = !!upCluster
          ? `/gloo-instances/${glooInstNamespace}/${glooInstName}/upstreams/${upCluster}/${upNamespace}/${upName}`
          : `/gloo-instances/${glooInstNamespace}/${glooInstName}/upstreams/${upNamespace}/${upName}`;
        return (
          <Spacer mb={2} key={upstream.name + '::' + upstream.namespace}>
            <a
              className={'cursor-pointer text-blue-500gloo text-base'}
              onClick={() => navigate(upstreamDetailsUrl)}>
              {upName}
            </a>
          </Spacer>
        );
      })}
    </div>
  );
};

export default ExecutableGraphqlUpstreamsTable;
