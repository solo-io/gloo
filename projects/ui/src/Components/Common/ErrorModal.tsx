import styled from '@emotion/styled/macro';
import { Modal } from 'antd';
import { ModalProps } from 'antd/lib/modal';
import { ReactComponent as CloseIcon } from 'assets/icon-close.svg';
import * as React from 'react';
import { colors } from 'Styles/colors';

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
  background: #f4c9d1;
  border: 2px solid #b05464;
  position: relative;
`;
const ContentText = styled.div`
  margin-top: 30px;
  font-size: 22px;
  color: ${colors.novemberGrey};
  width: 100%;
`;

const DescriptionText = styled.div`
  margin-top: 30px;
  font-size: 16px;
  color: ${colors.juneGrey};
  width: 100%;
`;

const ErrorX = styled(CloseIcon)`
  .secondary {
    fill: #b05464;
  }
`;

export interface ErrorModalProps extends ModalProps {
  errorMessage?: string;
  visible?: boolean;
  isNegative?: boolean;
  errorDescription?: string;
  cancel: () => void;
}

const ErrorModal = (props: ErrorModalProps) => {
  const { errorMessage, errorDescription } = props;
  const closeModal = (): void => {
    props.cancel();
  };
  return (
    <>
      <Modal
        visible={props.visible}
        footer={null}
        width={360}
        onCancel={closeModal}
        maskStyle={maskStyle}
        style={floatingStyle}
        bodyStyle={bodyStyle}
      >
        <ContentContainer>
          <WarningCircle>
            <ErrorX className='w-24 h-24 ' />
          </WarningCircle>
          <ContentText>
            {!!errorMessage ? errorMessage : 'There was an error'}
          </ContentText>
          <DescriptionText>
            {!!errorDescription ? errorDescription : ''}
          </DescriptionText>{' '}
        </ContentContainer>
      </Modal>
    </>
  );
};

export { ErrorModal as default };
