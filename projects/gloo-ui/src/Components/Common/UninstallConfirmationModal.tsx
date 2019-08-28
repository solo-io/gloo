import styled from '@emotion/styled';
import { Modal } from 'antd';
import { ModalProps } from 'antd/lib/modal';
import { ReactComponent as WarningExclamation } from 'assets/big-warning-exclamation.svg';
import * as React from 'react';
import { colors } from 'Styles';
import {
  ButtonNegativeProgress,
  SoloCancelButton,
  SoloNegativeButton
} from 'Styles/CommonEmotions/button';

const maskStyle = {
  background: 'transparent'
};
const floatingStyle = {
  borderRadius: '10px'
};
const bodyStyle = {
  borderRadius: '10px'
};

const ContentContainer = styled.div`
  width: 200px;
  margin: 25px auto 50px;
  text-align: center;
`;

const WarningCircle = styled.div`
  display: inline-flex;
  justify-content: center;
  align-items: center;
  width: 128px;
  height: 128px;
  border-radius: 100%;
  background: ${colors.flashlightGold};
  border: 2px solid ${colors.sunGold};
`;
const ContentText = styled.div`
  margin-top: 30px;
  font-size: 22px;
  color: ${colors.novemberGrey};
  width: 100%;
`;

const ButtonGroup = styled.div`
  display: flex;
  margin-top: 15px;
  justify-content: space-between;

  > button {
    min-width: 0;
  }
`;

interface Props extends ModalProps {
  visible: boolean;
  confirmationText?: string;
  doUninstall: () => any;
  cancel: () => any;
  uninstalling: boolean;
}

export const UninstallConfirmationModal = (props: Props) => {
  const closeModal = (): void => {
    props.cancel();
  };

  const { confirmationText, doUninstall } = props;

  return (
    <React.Fragment>
      <Modal
        visible={props.visible}
        footer={null}
        onCancel={closeModal}
        width={360}
        maskStyle={maskStyle}
        style={floatingStyle}
        bodyStyle={bodyStyle}>
        <ContentContainer>
          <WarningCircle>
            <WarningExclamation />
          </WarningCircle>
          <ContentText>
            {confirmationText
              ? confirmationText
              : 'Are you sure you want to uninstall this?'}
          </ContentText>

          <ButtonGroup>
            {/*
            // @ts-ignore*/}
            <SoloNegativeButton
              onClick={doUninstall}
              inProgress={props.uninstalling}>
              <ButtonNegativeProgress />
              Uninstall
            </SoloNegativeButton>
            <SoloCancelButton onClick={closeModal}>Cancel</SoloCancelButton>
          </ButtonGroup>
        </ContentContainer>
      </Modal>
    </React.Fragment>
  );
};
