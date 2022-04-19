import SoloAddButton from 'Components/Common/SoloAddButton';
import { GlooInstance } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/glooinstance_pb';
import React from 'react';
import { NewApiModal } from './NewApiModal';

export const NewApiButton: React.FC<{
  glooInstance: GlooInstance.AsObject;
}> = ({ glooInstance }) => {
  const [isModalVisible, setIsModalVisible] = React.useState(false);

  return (
    <>
      <SoloAddButton
        data-testid='landing-create-api'
        onClick={() => setIsModalVisible(true)}
        className='absolute right-0 -top-8 '>
        Create API
      </SoloAddButton>

      <NewApiModal
        glooInstance={glooInstance}
        show={isModalVisible}
        onClose={() => setIsModalVisible(false)}
      />
    </>
  );
};
