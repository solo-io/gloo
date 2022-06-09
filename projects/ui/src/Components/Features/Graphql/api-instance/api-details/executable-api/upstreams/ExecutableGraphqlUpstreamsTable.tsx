import { useGetGraphqlApiDetails, useListUpstreams } from 'API/hooks';
import { ClusterObjectRef } from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';
import { Upstream } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/gloo_resources_pb';
import React, { useMemo } from 'react';
import { useNavigate } from 'react-router';

const ExecutableGraphqlUpstreamsTable: React.FC<{
  apiRef: ClusterObjectRef.AsObject;
}> = ({ apiRef }) => {
  const navigate = useNavigate();

  // api hooks
  const { data: graphqlApi } = useGetGraphqlApiDetails(apiRef);
  // TODO: The referenced upstreams should be returned as part of the graphql details response.
  const maxUpstreams = 500;
  const { data: upstreamsResponse } = useListUpstreams(undefined, {
    limit: maxUpstreams,
    offset: 0,
  });
  const upstreams = upstreamsResponse?.upstreamsList;

  const resolverUpstreams = useMemo<Upstream.AsObject[]>(() => {
    let resUpstreams =
      graphqlApi?.spec?.executableSchema?.executor?.local?.resolutionsMap
        .filter(
          ([_rName, r], index, arr) =>
            index ===
            arr?.findIndex(
              ([_n, rr]) =>
                rr?.restResolver?.upstreamRef?.name ===
                r.restResolver?.upstreamRef?.name
            )
        )
        .map(([_resolveName, resolver]) => resolver.restResolver?.upstreamRef);
    let fullUpstreams = upstreams?.filter(
      upstream =>
        !!resUpstreams?.find(
          rU =>
            rU?.name === upstream.metadata?.name &&
            rU?.namespace === upstream.metadata?.namespace
        )
    );
    return !!fullUpstreams ? fullUpstreams : [];
    /* eslint-disable-next-line react-hooks/exhaustive-deps */
  }, [
    !!graphqlApi,
    !!upstreams,
    graphqlApi?.spec?.executableSchema?.executor?.local?.resolutionsMap?.length,
  ]);

  return (
    <div className='mb-8'>
      {resolverUpstreams?.map(resolverUpstream => {
        const glooInstNamespace = resolverUpstream.glooInstance?.namespace;
        const glooInstName = resolverUpstream.glooInstance?.name;
        const upstreamCluster = resolverUpstream.metadata?.clusterName ?? '';
        const upstreamNamespace = resolverUpstream.metadata?.namespace ?? '';
        const upstreamName = resolverUpstream.metadata?.name ?? '';
        const link = !!upstreamCluster
          ? `/gloo-instances/${glooInstNamespace}/${glooInstName}/upstreams/${upstreamCluster}/${upstreamNamespace}/${upstreamName}`
          : `/gloo-instances/${glooInstNamespace}/${glooInstName}/upstreams/${upstreamNamespace}/${upstreamName}`;
        return (
          <div key={link} className='mb-2'>
            <a
              className={'cursor-pointer text-blue-500gloo text-base'}
              onClick={() => navigate(link)}>
              {upstreamName}
            </a>
          </div>
        );
      })}
    </div>
  );
};

export default ExecutableGraphqlUpstreamsTable;
