import SoloAddButton from 'Components/Common/SoloAddButton';
import React from 'react';
import { NewApiModal } from './NewApiModal';

export const NewApiButton: React.FC = () => {
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
        show={isModalVisible}
        onClose={() => setIsModalVisible(false)}
      />
    </>
  );
};
