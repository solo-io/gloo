import styled from '@emotion/styled';
import { ReactComponent as CloseX } from 'assets/close-x.svg';
import * as React from 'react';
import { colors, soloConstants } from 'Styles';

const ModalWindow = styled.div`
  position: fixed;
  left: 0;
  right: 0;
  top: 0;
  bottom: 0;
  display: grid;
  justify-content: center;
  align-content: center;
  background: rgba(0, 0, 0, 0.1);
  /* overflow: auto; */
  line-height: 19px;
  z-index: 100;
`;

const BlockHolder = styled.div`
  max-height: 80vh;
  height: 60vh;
`;

type ModalBlockProps = {
  width: number | string;
};
const ModalBlock = styled.div`
  position: relative;
  max-width: 100%;
  width: ${(props: ModalBlockProps) =>
    props.width === 'auto' ? props.width : `${props.width}px`};
  border-radius: 10px;
  background: white;
`;

const Title = styled.div`
  font-size: 22px;
  font-weight: bold;
  color: ${colors.novemberGrey};
  line-height: normal;
  padding: ${soloConstants.largeBuffer}px ${soloConstants.smallBuffer}px 13px;
`;

const CloseXContainer = styled.div`
  position: absolute;
  display: flex;
  right: 16px;
  top: 16px;
  z-index: 2;

  > svg {
    width: 16px;
    height: 16px;
    cursor: pointer;
  }
`;

interface ContentProps {
  noPadding?: boolean;
}
const Content = styled.div`
  padding: 0
    ${(props: ContentProps) =>
      !!props.noPadding
        ? '0'
        : `${soloConstants.smallBuffer}px ${soloConstants.largeBuffer}px`};
`;

interface ModalProps {
  visible: boolean;
  width: number;
  title?: string | React.ReactNode;
  children: React.ReactChild;
  onClose?: () => any;
  noPadding?: boolean;
}

export const SoloModal = (props: ModalProps) => {
  const { visible, width, title, children, onClose, noPadding } = props;

  if (!visible) {
    document.body.style.overflow = 'auto';
    return null;
  }

  document.body.style.overflow = 'hidden';

  return (
    <ModalWindow>
      <BlockHolder
        onClick={(evt: React.SyntheticEvent) => evt.stopPropagation()}>
        <ModalBlock width={width}>
          {!!onClose && (
            <CloseXContainer onClick={onClose}>
              <CloseX />
            </CloseXContainer>
          )}
          {!!title && <Title>{title}</Title>}
          <Content noPadding={noPadding}>{children}</Content>
        </ModalBlock>
      </BlockHolder>
    </ModalWindow>
  );
};
