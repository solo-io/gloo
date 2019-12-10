import styled from '@emotion/styled';
import { Divider } from 'antd';
import { ReactComponent as GreenPlus } from 'assets/small-green-plus.svg';
import { SoloModal } from 'Components/Common/SoloModal';
import * as React from 'react';
import { colors } from 'Styles';
import { CreateRouteTableForm } from './CreateRouteTableForm';

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

export const CreateRouteTableModal = (props: Props) => {
  const [showModal, setShowModal] = React.useState(false);

  return (
    <ModalContainer data-testid='create-route-table-modal'>
      {props.withoutDivider ? (
        <ModalTrigger onClick={() => setShowModal(s => !s)}>
          <>
            {props.promptText || 'Create a Route Table'}
            <StyledGreenPlus style={{ marginLeft: '7px', width: '18px' }} />
          </>
        </ModalTrigger>
      ) : (
        <ModalTrigger onClick={() => setShowModal(s => !s)}>
          <>
            <StyledGreenPlus />
            {props.promptText || 'Create a Route Table'}
          </>
          <Divider type='vertical' style={{ height: '1.5em' }} />
        </ModalTrigger>
      )}
      <SoloModal
        visible={showModal}
        width={650}
        title='Create a Route Table'
        onClose={() => setShowModal(false)}>
        <>
          <Legend>
            The RouteTable is a child Routing object for the Gloo Gateway. A
            RouteTable gets built into the complete routing configuration when
            it is referenced by a delegateAction, either in a parent
            VirtualService or another RouteTable. Routes specified in a
            RouteTable must have their paths start with the prefix provided in
            the parent’s matcher.
          </Legend>
          <CreateRouteTableForm toggleModal={setShowModal} />
        </>
      </SoloModal>
    </ModalContainer>
  );
};
