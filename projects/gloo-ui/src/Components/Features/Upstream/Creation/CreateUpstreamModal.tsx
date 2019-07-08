import * as React from 'react';
/** @jsx jsx */
import { jsx } from '@emotion/core';

import styled from '@emotion/styled/macro';
import { colors } from 'Styles';
import { Divider } from 'antd';
import { SoloModal } from 'Components/Common/SoloModal';
import { CreateUpstreamForm } from './CreateUpstreamForm';
import { ReactComponent as GreenPlus } from 'assets/small-green-plus.svg';

interface Props {}

const StyledGreenPlus = styled(GreenPlus)`
  cursor: pointer;
  margin-right: 7px;
  .a {
    fill: ${colors.forestGreen};
  }
`;

const ModalContainer = styled.div`
  display: flex;
  flex-direction: row;
  align-content: center;
`;
const Legend = styled.div`
  background-color: ${colors.januaryGrey};
`;

// TODO: use spec font
const ModalTrigger = styled.div`
  cursor: pointer;
  display: flex;
  align-items: center;
  padding: 0 10px;
  font-size: 14px;
`;
export const CreateUpstreamModal = (props: Props) => {
  const [showModal, setShowModal] = React.useState(false);

  return (
    <ModalContainer>
      <ModalTrigger onClick={() => setShowModal(s => !s)}>
        <React.Fragment>
          <StyledGreenPlus />
          Create Upstream
        </React.Fragment>
        <Divider type='vertical' style={{ height: '1.5em' }} />
      </ModalTrigger>
      <SoloModal
        visible={showModal}
        width={650}
        title='Create an Upstream'
        onClose={() => setShowModal(false)}>
        <React.Fragment>
          <Legend>
            Lorem ipsum dolor sit amet consectetur adipisicing elit. Aspernatur
            officia deleniti ullam hic nostrum quod explicabo optio accusantium,
            maiores cumque asperiores! Consectetur illum omnis eum qui
            reprehenderit in eaque doloremque!
          </Legend>
          <CreateUpstreamForm />
        </React.Fragment>
      </SoloModal>
    </ModalContainer>
  );
};
