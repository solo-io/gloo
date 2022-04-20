import { ClusterObjectRef } from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';
import React from 'react';
import { SoloNegativeButton } from 'Styles/StyledComponents/button';
import { useDeleteApi } from 'utils/hooks';

const GraphqlDeleteApiButton: React.FC<{
  apiRef: ClusterObjectRef.AsObject;
}> = ({ apiRef }) => {
  const confirmDeleteApi = useDeleteApi();

  return (
    <>
      <SoloNegativeButton
        data-testid='delete-api'
        onClick={() => confirmDeleteApi(apiRef)}>
        Delete API
      </SoloNegativeButton>
    </>
  );
};

export default GraphqlDeleteApiButton;
