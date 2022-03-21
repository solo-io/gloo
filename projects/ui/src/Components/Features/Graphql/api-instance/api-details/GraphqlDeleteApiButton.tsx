import { useGetGraphqlApiDetails } from 'API/hooks';
import ConfirmationModal from 'Components/Common/ConfirmationModal';
import ErrorModal from 'Components/Common/ErrorModal';
import { ClusterObjectRef } from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';
import React from 'react';
import { SoloNegativeButton } from 'Styles/StyledComponents/button';
import { useDeleteAPI } from 'utils/hooks';

const GraphqlDeleteApiButton: React.FC<{
  apiRef: ClusterObjectRef.AsObject;
}> = ({ apiRef }) => {
  const { mutate } = useGetGraphqlApiDetails(apiRef);
  const {
    isDeleting,
    triggerDelete,
    cancelDelete,
    closeErrorModal,
    errorModalIsOpen,
    errorDeleteModalProps,
    deleteFn,
  } = useDeleteAPI({ revalidate: mutate });

  return (
    <>
      <SoloNegativeButton
        data-testid='delete-api'
        onClick={() => triggerDelete(apiRef)}>
        Delete API
      </SoloNegativeButton>
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

export default GraphqlDeleteApiButton;
