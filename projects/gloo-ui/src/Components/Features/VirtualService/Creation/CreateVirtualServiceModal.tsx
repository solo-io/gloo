import styled from '@emotion/styled';
import { Divider } from 'antd';
import { ReactComponent as GreenPlus } from 'assets/small-green-plus.svg';
import { SoloModal } from 'Components/Common/SoloModal';
import * as React from 'react';
import { colors } from 'Styles';
import { CreateVirtualServiceForm } from './CreateVirtualServiceForm';

const StyledGreenPlus = styled(GreenPlus)`
  cursor: pointer;
  margin-right: 7px;
`;

const ModalContainer = styled.div`
  display: flex;
  flex-direction: row;
  align-content: center;
`;

const Legend = styled.div`
  background-color: ${colors.januaryGrey};
  padding: 13px 13px 13px 10px;
  margin-bottom: 23px;
`;
// TODO: use spec font
const ModalTrigger = styled.div`
  cursor: pointer;
  display: flex;
  align-items: center;
  padding: 0 10px;
  font-size: 14px;
`;

interface Props {
  withoutDivider?: boolean;
  promptText?: string;
}

export const CreateVirtualServiceModal = (props: Props) => {
  const [showModal, setShowModal] = React.useState(false);

  return (
    <ModalContainer data-testid='create-virtual-service-modal'>
      {props.withoutDivider ? (
        <ModalTrigger onClick={() => setShowModal(s => !s)}>
          <>
            {props.promptText || 'Create Virtual Service'}
            <StyledGreenPlus style={{ marginLeft: '7px', width: '18px' }} />
          </>
        </ModalTrigger>
      ) : (
        <ModalTrigger onClick={() => setShowModal(s => !s)}>
          <>
            <StyledGreenPlus />
            {props.promptText || 'Create Virtual Service'}
          </>
          <Divider type='vertical' style={{ height: '1.5em' }} />
        </ModalTrigger>
      )}
      <SoloModal
        visible={showModal}
        width={650}
        title='Create a Virtual Service'
        onClose={() => setShowModal(false)}>
        <>
          <Legend>
            Virtual Services define a set of route rules, an optional SNI
            configuration for a given domain or set of domains.
          </Legend>
          <CreateVirtualServiceForm toggleModal={setShowModal} />
        </>
      </SoloModal>
    </ModalContainer>
  );
};
