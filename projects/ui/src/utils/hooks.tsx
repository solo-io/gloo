import { graphqlApi } from 'API/graphql';
import {
  useListClusterDetails,
  useListGlooInstances,
  useListGraphqlSchemas,
} from 'API/hooks';
import { ClusterObjectRef } from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';
import { GraphqlSchema } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/graphql_pb';
import React from 'react';
import { useNavigate } from 'react-router';
import { KeyedMutator } from 'swr';

interface DeleteConfig {
  optimistic?: boolean;
  revalidate:
    | KeyedMutator<GraphqlSchema.AsObject[]>
    | KeyedMutator<GraphqlSchema.AsObject>;
}

// TODO: make reusable
export function useDeleteAPI(config: DeleteConfig) {
  let { revalidate, optimistic } = config;

  const [isDeleting, setIsDeleting] = React.useState(false);
  const [errorModal, setErrorModal] = React.useState(false);
  const [errorMessage, setErrorMessage] = React.useState('');
  const [errorDescription, setErrorDescription] = React.useState('');
  const [schemaRefToDelete, setSchemaRefToDelete] =
    React.useState<ClusterObjectRef.AsObject>();
  const navigate = useNavigate();
  const { data: graphqlSchemas, error: graphqlSchemaError } =
    useListGraphqlSchemas();
  const { data: glooInstances, error: instancesError } = useListGlooInstances();

  const { data: clusterDetailsList, error: cError } = useListClusterDetails();

  const triggerDelete = (schemaRef: ClusterObjectRef.AsObject) => {
    setSchemaRefToDelete(schemaRef);
    setIsDeleting(true);
  };
  const cancelDelete = () => {
    setIsDeleting(false);
  };
  const deleteAPI = async () => {
    if (schemaRefToDelete) {
      if (optimistic) {
        setTimeout(() => {
          revalidate(
            graphqlSchemas?.filter(
              gqlSchema =>
                gqlSchema.metadata?.name !== schemaRefToDelete?.name &&
                gqlSchema.metadata?.namespace !== schemaRefToDelete?.namespace
            ),
            false
          );
        }, 300);
      }
      try {
        let deletedSchemaRef = await graphqlApi.deleteGraphqlSchema({
          name: schemaRefToDelete?.name,
          namespace: schemaRefToDelete?.namespace,
          clusterName: schemaRefToDelete?.clusterName,
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
