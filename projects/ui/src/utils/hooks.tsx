import { graphqlConfigApi } from 'API/graphql';
import {
  useListClusterDetails,
  useListGlooInstances,
  useListGraphqlApis,
} from 'API/hooks';
import { ClusterObjectRef } from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';
import { GraphqlApi } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/graphql_pb';
import React from 'react';
import { useNavigate } from 'react-router';
import { KeyedMutator } from 'swr';

interface DeleteConfig {
  optimistic?: boolean;
  revalidate:
    | KeyedMutator<GraphqlApi.AsObject[]>
    | KeyedMutator<GraphqlApi.AsObject>;
}

// TODO: make reusable
export function useDeleteAPI(config: DeleteConfig) {
  let { revalidate, optimistic } = config;

  const [isDeleting, setIsDeleting] = React.useState(false);
  const [errorModal, setErrorModal] = React.useState(false);
  const [errorMessage, setErrorMessage] = React.useState('');
  const [errorDescription, setErrorDescription] = React.useState('');
  const [apiRefToDelete, setApiRefToDelete] =
    React.useState<ClusterObjectRef.AsObject>();
  const navigate = useNavigate();
  const { data: graphqlApis, error: graphqlApiError } = useListGraphqlApis();
  const { data: glooInstances, error: instancesError } = useListGlooInstances();

  const { data: clusterDetailsList, error: cError } = useListClusterDetails();

  const triggerDelete = (apiRef: ClusterObjectRef.AsObject) => {
    setApiRefToDelete(apiRef);
    setIsDeleting(true);
  };
  const cancelDelete = () => {
    setIsDeleting(false);
  };
  const deleteAPI = async () => {
    if (apiRefToDelete) {
      if (optimistic) {
        setTimeout(() => {
          revalidate(
            graphqlApis?.filter(
              gqlApi =>
                gqlApi.metadata?.name !== apiRefToDelete?.name &&
                gqlApi.metadata?.namespace !== apiRefToDelete?.namespace
            ),
            false
          );
        }, 300);
      }
      try {
        let deletedApiRef = await graphqlConfigApi.deleteGraphqlApi({
          name: apiRefToDelete?.name,
          namespace: apiRefToDelete?.namespace,
          clusterName: apiRefToDelete?.clusterName,
        });
        cancelDelete();
        if (optimistic) {
          setTimeout(() => {
            revalidate();
          }, 300);
        }

        navigate(
          clusterDetailsList?.length === 1 && glooInstances?.length === 1
            ? `/gloo-instances/${
                clusterDetailsList[0]!.glooInstancesList[0].metadata?.namespace
              }/${
                clusterDetailsList[0]!.glooInstancesList[0].metadata?.name
              }/apis/`
            : '/apis/'
        );
      } catch (error: any) {
        cancelDelete();
        setErrorMessage('API deletion failed');
        setErrorDescription(error.message ?? '');
        setErrorModal(true);
      }
    }
  };

  return {
    isDeleting,
    triggerDelete,
    cancelDelete,
    closeErrorModal: () => setErrorModal(false),
    errorModalIsOpen: errorModal,
    deleteFn: deleteAPI,
    errorDeleteModalProps: {
      cancel: () => setErrorModal(false),
      visible: errorModal,
      errorDescription,
      errorMessage,
      isNegative: true,
    },
  };
}
