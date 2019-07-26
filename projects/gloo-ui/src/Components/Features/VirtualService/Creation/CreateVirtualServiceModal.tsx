import * as React from 'react';
/** @jsx jsx */
import { jsx } from '@emotion/core';

import styled from '@emotion/styled/macro';
import { withRouter, RouteComponentProps } from 'react-router';
import { colors } from 'Styles';
import { SoloModal } from 'Components/Common/SoloModal';
import { ReactComponent as GreenPlus } from 'assets/small-green-plus.svg';
import { Divider } from 'antd';
import { CreateVirtualServiceForm } from './CreateVirtualServiceForm';
interface Props {}

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
export const CreateVirtualServiceModal = (props: Props) => {
  const [showModal, setShowModal] = React.useState(false);

  return (
    <ModalContainer>
      <ModalTrigger onClick={() => setShowModal(s => !s)}>
        <React.Fragment>
          <StyledGreenPlus />
          Create Virtual Service
        </React.Fragment>
        <Divider type='vertical' style={{ height: '1.5em' }} />
      </ModalTrigger>
      <SoloModal
        visible={showModal}
        width={650}
        title='Create a Virtual Service'
        onClose={() => setShowModal(false)}>
        <React.Fragment>
          <Legend>
            Virtual Services define a set of route rules, an optional SNI
            configuration for a given domain or set of domains.
          </Legend>
          <CreateVirtualServiceForm />
        </React.Fragment>
      </SoloModal>
    </ModalContainer>
  );
};
