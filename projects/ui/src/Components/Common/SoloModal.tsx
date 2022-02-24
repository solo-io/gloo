import * as React from 'react';
import styled from '@emotion/styled';
import { ReactComponent as CloseX } from 'assets/close-x.svg';
import { colors } from 'Styles/colors';
import { useTheme } from '@emotion/react';
import { Tab, TabPanelProps, Tabs } from '@reach/tabs';
import { ReactComponent as Checkmark } from 'assets/success-checkmark.svg';

const StyledTabComponent = styled(Tab)<{ isSelected: boolean }>`
  ${props =>
    props.isSelected
      ? ` background-color: ${colors.seaBlue}; 
          border-right-width: 8px;
          border-color: ${colors.pondBlue}`
      : ` background-color: ${colors.oceanBlue};
        border-right-width: 8px;
        border-color: ${colors.oceanBlue}`}
`;

export const StyledModalTab = (
  props: {
    disabled?: boolean | undefined;
    isCompleted?: boolean | undefined;
    isSelected?: boolean | undefined;
  } & TabPanelProps
) => {
  const { isSelected, children, isCompleted, ...otherProps } = props;
  const theme = useTheme();
  return (
    <StyledTabComponent
      {...otherProps}
      isSelected={!!isSelected}
      className='flex justify-between p-1 pl-6 mb-2 text-left text-white focus:outline-none'>
      {children}

      <div
        style={{ backgroundColor: colors.pondBlue }}
        className={`w-4 h-4 rounded-full text-white justify-center items-center ${
          isCompleted ? 'flex' : 'hidden'
        }`}>
        <Checkmark className={'fill-current w-2 h-2'} />
      </div>
    </StyledTabComponent>
  );
};

export const StyledModalTabs = styled(Tabs)`
  display: grid;
  grid-template-columns: 190px 1fr;
`;

const ModalWindow = styled.div`
  position: fixed;
  left: 0;
  right: 0;
  top: 0;
  bottom: 0;
  display: grid;
  justify-content: center;
  padding-top: 100px;
  background: rgba(0, 0, 0, 0.1);
  overflow: auto;
  line-height: 19px;
  z-index: 100;
  border-radius: 10px;
`;

const BlockHolder = styled.div`
  max-height: 80vh;
  border-radius: 10px;
`;

interface ModalBlockProps {
  width: number | string;
}
const ModalBlock = styled.div<ModalBlockProps>`
  position: relative;
  width: ${(props: ModalBlockProps) =>
    props.width === 'auto' ? props.width : `${props.width}px`};
  max-width: calc(100vw - 50px);
  border-radius: 10px;
  background: white;
`;

const Title = styled.div`
  font-size: 22px;
  font-weight: bold;
  color: ${colors.novemberGrey};
  line-height: normal;
  padding: 23px 18px 13px;
`;

const CloseXContainer = styled.div`
  position: absolute;
  display: flex;
  right: 18px;
  top: 20px;
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
  padding: 0;
`;

interface ModalProps {
  visible: boolean;
  width: number;
  title?: React.ReactNode;
  children: React.ReactChild;
  onClose?: () => any;
}

export const SoloModal = (props: ModalProps) => {
  const { visible, width, title, children, onClose } = props;

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
          <Content>{children}</Content>
        </ModalBlock>
      </BlockHolder>
    </ModalWindow>
  );
};
