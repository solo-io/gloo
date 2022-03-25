import { graphqlConfigApi } from 'API/graphql';
import { ClusterObjectRef } from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';
import React from 'react';
import {
  SoloButtonStyledComponent,
  SoloCancelButton,
} from 'Styles/StyledComponents/button';
import { UpdateApiModal } from '../../update-api-modal/UpdateApiModal';

const GraphqlEditApiButton: React.FC<{
  apiRef: ClusterObjectRef.AsObject;
}> = ({ apiRef }) => {
  const [showEditModal, setShowEditModal] = React.useState(false);
  const [schemaString, setSchemaString] = React.useState('');

  const toggleModal = () => {
    setShowEditModal(!showEditModal);
  };

  React.useEffect(() => {
    graphqlConfigApi.getGraphqlApi(apiRef).then(res => {
      const schemaDef = res.spec!.executableSchema!.schemaDefinition!;
      setSchemaString(schemaDef);
    });
  }, [graphqlConfigApi.getGraphqlApi, apiRef]);

  return (
    <>
      <SoloButtonStyledComponent
        data-testid='edit-api'
        onClick={() => toggleModal()}>
        Update Schema
      </SoloButtonStyledComponent>
      <UpdateApiModal
        schemaString={schemaString}
        apiRef={apiRef}
        show={showEditModal}
        onClose={toggleModal}
      />
    </>
  );
};

export default GraphqlEditApiButton;
