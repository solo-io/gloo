import styled from '@emotion/styled';
import { graphqlConfigApi } from 'API/graphql';
import {
  useGetGraphqlApiDetails,
  useGetGraphqlApiYaml,
  useListUpstreams,
} from 'API/hooks';
import AreaHeader from 'Components/Common/AreaHeader';
import ConfirmationModal from 'Components/Common/ConfirmationModal';
import ErrorModal from 'Components/Common/ErrorModal';
import { ClusterObjectRef } from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';
import { Upstream } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/gloo_resources_pb';
import React from 'react';
import { useNavigate, useParams } from 'react-router';
import { SoloNegativeButton } from 'Styles/StyledComponents/button';
import { useDeleteAPI } from 'utils/hooks';
import ResolversTable from './ResolversTable';

const ConfigArea = styled.div`
  margin-bottom: 20px;
`;

export const GraphqlApiDetails: React.FC<{
  apiRef: ClusterObjectRef.AsObject;
}> = ({ apiRef }) => {
  const navigate = useNavigate();

  // api hooks
  const {
    data: graphqlApi,
    error: graphqlApiError,
    mutate,
  } = useGetGraphqlApiDetails(apiRef);
  const { data: graphqlApiYaml, error: graphqlApiYamlError } =
    useGetGraphqlApiYaml(apiRef);
  const {
    isDeleting,
    triggerDelete,
    cancelDelete,
    closeErrorModal,
    errorModalIsOpen,
    errorDeleteModalProps,
    deleteFn,
  } = useDeleteAPI({ revalidate: mutate });
  const { data: upstreams, error: upstreamsError } = useListUpstreams();

  const loadYaml = async () => {
    if (!apiRef.name || !apiRef.namespace) {
      return '';
    }
    try {
      const yaml = await graphqlConfigApi.getGraphqlApiYaml(apiRef);
      return yaml;
    } catch (error) {
      console.error(error);
    }
    return '';
  };

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
      <ConfigArea className='-mt-1'>
        <AreaHeader
          title='Configuration'
          contentTitle={`${apiRef.namespace}--${apiRef.name}.yaml`}
          yaml={graphqlApiYaml}
          onLoadContent={loadYaml}
        />
        <ResolversTable apiRef={apiRef} />
      </ConfigArea>
      <ConfigArea>
        <div className='flex p-4 mb-5 bg-gray-100 border border-gray-300 rounded-lg'>
          <div className='w-1/5 mr-5'>
            <div className='mb-2 text-lg font-medium'>Upstreams</div>
            {resolverUpstreams?.map(resolverUpstream => {
              const glooInstNamespace =
                resolverUpstream.glooInstance?.namespace;
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
      </ConfigArea>
      <div>
        <SoloNegativeButton
          data-testid='delete-api'
          onClick={() => triggerDelete(apiRef)}>
          Delete API
        </SoloNegativeButton>
      </div>
      <ConfirmationModal
        visible={isDeleting}
        confirmPrompt='delete this API'
        confirmButtonText='Delete'
        goForIt={deleteFn}
        cancel={cancelDelete}
        isNegative
      />
      <ErrorModal
        {...errorDeleteModalProps}
        cancel={closeErrorModal}
        visible={errorModalIsOpen}
      />
    </>
  );
};
