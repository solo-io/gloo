import styled from '@emotion/styled';
import { Modal } from 'antd';
import { ModalProps } from 'antd/lib/modal';
import { ReactComponent as SuccessCheckmark } from 'assets/big-successful-checkmark.svg';
import * as React from 'react';
import { colors } from 'Styles';

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
`;
const ContentText = styled.div`
  margin-top: 30px;
  font-size: 22px;
  color: ${colors.novemberGrey};
  width: 100%;
`;

interface Props extends ModalProps {
  visible: boolean;
  successMessage: string;
}

export const SuccessModal = (props: Props) => {
  const { successMessage, visible } = props;
  const [showModal, setShowModal] = React.useState(visible);

  React.useEffect(() => {
    if (visible) {
      setShowModal(true);
      setTimeout(() => {
        setShowModal(false);
      }, 2000);
    }
  }, [visible]);

  return (
    <React.Fragment>
      <Modal
        visible={showModal}
        footer={null}
        width={360}
        maskStyle={maskStyle}
        style={floatingStyle}
        bodyStyle={bodyStyle}>
        <ContentContainer>
          <WarningCircle>
            <SuccessCheckmark />
          </WarningCircle>
          <ContentText>{successMessage}</ContentText>
        </ContentContainer>
      </Modal>
    </React.Fragment>
  );
};
