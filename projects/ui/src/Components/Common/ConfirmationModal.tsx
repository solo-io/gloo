import styled from '@emotion/styled/macro';
import { Modal } from 'antd';
import { ModalProps } from 'antd/lib/modal';
import { ReactComponent as WarningExclamation } from 'assets/big-warning-exclamation.svg';
import * as React from 'react';
import { colors } from 'Styles/colors';
import {
  SoloButtonStyledComponent,
  SoloCancelButton,
  SoloNegativeButton,
} from 'Styles/StyledComponents/button';

const maskStyle = {
  background: 'rgba(255, 255, 255, 0)',
};
const floatingStyle = {
  borderRadius: '10px',
};
const bodyStyle = {
  borderRadius: '10px',
};

const ContentContainer = styled.div`
  width: 100%;
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
  justify-content: center;

  > button {
    min-width: 0;

    &:first-of-type {
      margin-right: 10px;
    }
  }
`;

interface Props extends ModalProps {
  confirmationTopic?: string;
  confirmText?: string;
  goForIt: () => any;
  cancel: () => any;
  visible?: boolean;
  isNegative?: boolean;
  confirmTestId?: string;
  cancelTestId?: string;
}

// TODO: make this workflow into a reusable hook
const ConfirmationModal = (props: Props) => {
  const closeModal = (): void => {
    props.cancel();
  };

  const { confirmationTopic, confirmText, goForIt, isNegative } = props;

  return (
    <>
      <Modal
        visible={props.visible}
        footer={null}
        onCancel={closeModal}
        width={360}
        maskStyle={maskStyle}
        style={floatingStyle}
        bodyStyle={bodyStyle}
      >
        <ContentContainer>
          <WarningCircle>
            <WarningExclamation className='w-16 h-16' />
          </WarningCircle>
          <ContentText>
            Are you sure you want to{' '}
            {confirmationTopic ? confirmationTopic : 'remove this'}?
          </ContentText>

          <ButtonGroup>
            {isNegative ? (
              <SoloNegativeButton
                data-testid={props.confirmTestId}
                onClick={goForIt}
              >
                {confirmText ? confirmText : 'Confirm'}
              </SoloNegativeButton>
            ) : (
              <SoloButtonStyledComponent
                data-testid={props.confirmTestId}
                onClick={goForIt}
              >
                {confirmText ? confirmText : 'Confirm'}
              </SoloButtonStyledComponent>
            )}
            <SoloCancelButton
              data-testid={props.cancelTestId}
              onClick={closeModal}
            >
              Cancel
            </SoloCancelButton>
          </ButtonGroup>
        </ContentContainer>
      </Modal>
    </>
  );
};

export { ConfirmationModal as default };
