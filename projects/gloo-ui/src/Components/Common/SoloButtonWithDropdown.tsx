import * as React from 'react';
/** @jsx jsx */
import { css, jsx } from '@emotion/core';
import styled from '@emotion/styled/macro';

import { colors, soloConstants } from 'Styles';
import { ButtonProgress } from 'Styles/CommonEmotions/button';
import { SoloButton } from './SoloButton';

export enum InstallExtensionButtonAction {
  install,
  update,
  uninstall
}

const Container = styled<'div', { inProgress?: boolean }>('div')`
  display: flex;

  ${props =>
    props.inProgress
      ? ''
      : `> button {
    border-radius: ${soloConstants.smallRadius}px 0 0
      ${soloConstants.smallRadius}px;
  }`};
`;

const OptionsPrompt = styled<'div', { disabled?: boolean }>('div')`
  position: relative;
  height: 39px;
  width: 30px;
  font-size: 16px;
  background: ${colors.seaBlue};
  border-left: 1px solid ${colors.marchGrey};
  border-radius: 0 ${soloConstants.smallRadius}px ${soloConstants.smallRadius}px
    0;
  cursor: pointer;

  &:hover,
  &:focus {
    background: ${colors.lakeBlue};
  }

  &:active {
    background: ${colors.seaBlue};
  }

  ${props =>
    props.disabled
      ? `
    opacity: 0.3;
    pointer-events: none;
    cursor: default;
    `
      : ''}
`;
const Dots = styled.div`
  position: absolute;
  top: 5px;
  left: 0;
  right: 0;
  text-align: center;
  font-weight: 600;
  color: white;
`;

const Option = styled.div`
  position: absolute;
  right: 0;
  top: 100%;

  > * {
    margin-top: 5px;
  }
`;

interface ButtonProps {
  inProgress?: boolean;
  disabled?: boolean;
  children?: React.ReactChild;
  onClick: () => any;
  otherOptions: React.ReactNode[];
  onOpenOptionsClick: () => any;
}

export const SoloButtonWithDropdown = (props: ButtonProps) => {
  const [optionsOpen, setOptionsOpen] = React.useState(false);
  const containerRef = React.useRef();

  React.useEffect(() => {
    // add when mounted
    document.addEventListener('mousedown', handleClick);
    // return function to be called when unmounted
    return () => {
      document.removeEventListener('mousedown', handleClick);
    };
  }, []);

  const handleClick = (evt: any) => {
    console.log(containerRef);
    if (
      !!containerRef &&
      !!containerRef.current &&
      // @ts-ignore
      containerRef.current!.contains(evt.target)
    ) {
      // Clicked on menu option. Don't do anything.
    } else {
      setOptionsOpen(false);
    }
  };

  const {
    disabled,
    inProgress,
    onClick,
    onOpenOptionsClick,
    otherOptions
  } = props;

  const openOptions = (evt: React.MouseEvent) => {
    evt.preventDefault();
    onOpenOptionsClick();
    setOptionsOpen(s => !s);
  };

  return (
    // @ts-ignore
    <Container inProgress={inProgress} ref={containerRef}>
      {/*
    // @ts-ignore*/}
      <SoloButton inProgress={inProgress} onClick={onClick} disabled={disabled}>
        {props.children}
      </SoloButton>

      {!inProgress && (
        <React.Fragment>
          <OptionsPrompt onClick={openOptions} disabled={disabled}>
            <Dots>. . .</Dots>
            {optionsOpen &&
              otherOptions.map((option, ind) => {
                return <Option key={ind}>{option}</Option>;
              })}
          </OptionsPrompt>
        </React.Fragment>
      )}
    </Container>
  );
};
