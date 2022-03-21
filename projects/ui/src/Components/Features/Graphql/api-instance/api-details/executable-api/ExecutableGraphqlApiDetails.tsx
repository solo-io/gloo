import { useGetGraphqlApiDetails, useListUpstreams } from 'API/hooks';
import { ClusterObjectRef } from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';
import { Upstream } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/gloo_resources_pb';
import React from 'react';
import { useNavigate } from 'react-router';
import GraphqlApiConfigurationHeader from '../GraphqlApiConfigurationHeader';
import GraphqlDeleteApiButton from '../GraphqlDeleteApiButton';
import ResolversTable from './ResolversTable';

export const ExecutableGraphqlApiDetails: React.FC<{
  apiRef: ClusterObjectRef.AsObject;
}> = ({ apiRef }) => {
  const navigate = useNavigate();

  // api hooks
  const { data: graphqlApi } = useGetGraphqlApiDetails(apiRef);
  const { data: upstreams } = useListUpstreams();

  // These resolver upstreams are listed in the <ResolversTable> component.
  const [resolverUpstreams, setResolverUpstreams] = React.useState<
    Upstream.AsObject[]
  >([]);
  React.useEffect(() => {
    let resolverUpstreams =
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
        !!resolverUpstreams?.find(
          rU =>
            rU?.name === upstream.metadata?.name &&
            rU?.namespace === upstream.metadata?.namespace
        )
    );
    if (!!fullUpstreams) {
      setResolverUpstreams(fullUpstreams);
    }
    /* eslint-disable-next-line react-hooks/exhaustive-deps */
  }, [
    !!graphqlApi,
    !!upstreams,
    graphqlApi?.spec?.executableSchema?.executor?.local?.resolutionsMap?.length,
  ]);

  return (
    <>
      <GraphqlApiConfigurationHeader apiRef={apiRef} />

      <ResolversTable apiRef={apiRef} />

      <div className='flex p-4 mb-5 bg-gray-100 border border-gray-300 rounded-lg'>
        <div className='w-1/5 mr-5'>
          <div className='mb-2 text-lg font-medium'>Upstreams</div>
          {resolverUpstreams?.map(resolverUpstream => {
            const glooInstNamespace = resolverUpstream.glooInstance?.namespace;
            const glooInstName = resolverUpstream.glooInstance?.name;
            const upstreamCluster =
              resolverUpstream.metadata?.clusterName ?? '';
            const upstreamNamespace =
              resolverUpstream.metadata?.namespace ?? '';
            const upstreamName = resolverUpstream.metadata?.name ?? '';
            const link = !!upstreamCluster
              ? `/gloo-instances/${glooInstNamespace}/${glooInstName}/upstreams/${upstreamCluster}/${upstreamNamespace}/${upstreamName}`
              : `/gloo-instances/${glooInstNamespace}/${glooInstName}/upstreams/${upstreamNamespace}/${upstreamName}`;
            return (
              <div key={link}>
                <div
                  className={'cursor-pointer text-blue-500gloo text-base'}
                  onClick={() => {
                    navigate(link);
                  }}>
                  {upstreamName}
                </div>
              </div>
            );
          })}
        </div>
      </div>

      <GraphqlDeleteApiButton apiRef={apiRef} />
    </>
  );
};
