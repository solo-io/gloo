import { css, keyframes } from '@emotion/core';
import styled from '@emotion/styled';
import { Button } from 'antd';
import blueProgressMask from 'assets/primary-progress-mask.svg';
import orangeProgressMask from 'assets/warning-progress-mask.svg';
import { ButtonProps } from 'antd/lib/button';
import { colors } from '../colors';

export const soloConstants = {
  smallBuffer: 18,
  buffer: 20,
  largeBuffer: 23,

  smallRadius: 8,
  radius: 10,
  largeRadius: 16,

  transitionTime: '.3s',
};
const slide = keyframes`
  from { background-position: -281px 0; }
  to { background-position: 0 0; }
`;

const ProgressSliderCSS = css`
  display: none;
  position: absolute;
  left: 0;
  right: 0;
  top: 0;
  bottom: 0;
  border-radius: ${soloConstants.smallRadius}px;
  background: url(${blueProgressMask}) repeat 0 0;
  animation: ${slide} 8s linear infinite;
  z-index: 1;
`;

export const ButtonProgress = styled.div`
  ${ProgressSliderCSS};
`;

interface SoloButtonProps extends Omit<ButtonProps, 'onClick'> {
  inProgress?: boolean;
  green?: boolean;
  disable?: boolean;
  ['data-testid']?: string;
}

export const SoloButtonStyledComponent = styled(Button)<SoloButtonProps>`
  position: relative;
  display: inline-block;
  padding: 10px 20px 13px;
  font-size: 16px;
  line-height: 14px;
  background: ${colors.pondBlue};
  color: white;
  border: none;
  outline: none;
  border-radius: ${soloConstants.smallRadius}px;
  min-width: 125px;
  height: auto;
  cursor: pointer;

  transition: min-width ${soloConstants.transitionTime};

  &:hover,
  &:focus {
    color: white;
    background: ${colors.splashBlue};
    outline: none;
  }

  &:active {
    background: ${colors.seaBlue};
  }

  &:disabled {
    opacity: 0.3;
    pointer-events: none;
    cursor: default;
    background: ${colors.seaBlue};
    color: white;
  }

  > span {
    position: relative;
    z-index: 2;
  }

  ${props =>
    props.green
      ? `background: ${colors.forestGreen};
    
    
        &:hover,
        &:focus {
            background: ${colors.standGreen};
          }
          
          &:active {
          background: ${colors.forestGreen};
        }`
      : ''};

  ${(props: SoloButtonProps) =>
    props.inProgress
      ? `
        pointer-events: none;
        cursor: default;
        min-width: 100%;
        
        > ${ButtonProgress} {
          display: block;
        }`
      : ''};
`;

export const SoloCancelButton = styled.button`
  position: relative;
  display: inline-block;
  padding: 10px 20px 13px;
  font-size: 16px;
  line-height: 14px;
  background: ${colors.pondBlue};
  color: white;
  border: none;
  outline: none;
  border-radius: ${soloConstants.smallRadius}px;
  min-width: 125px;
  height: auto;
  cursor: pointer;

  transition: min-width ${soloConstants.transitionTime};

  &:hover,
  &:focus {
    color: white;
    background: ${colors.puddleBlue};
    outline: none;
  }

  &:active {
    background: ${colors.seaBlue};
  }

  &:disabled {
    opacity: 0.3;
    pointer-events: none;
    cursor: default;
    background: ${colors.seaBlue};
    color: white;
  }

  > span {
    position: relative;
    z-index: 2;
  }
  background: ${colors.juneGrey};

  &:hover,
  &:focus {
    background: ${colors.mayGrey};
  }

  &:active {
    background: ${colors.juneGrey};
  }

  &:disabled {
    background: ${colors.juneGrey};
  }
`;

export const ButtonNegativeProgress = styled.div`
  ${ProgressSliderCSS};
  background: url(${orangeProgressMask}) repeat 0 0;
`;

export const SoloNegativeButton = styled.button`
  position: relative;
  display: inline-block;
  padding: 10px 20px 13px;
  font-size: 16px;
  line-height: 14px;
  background: ${colors.pondBlue};
  color: white;
  border: none;
  outline: none;
  border-radius: ${soloConstants.smallRadius}px;
  min-width: 125px;
  height: auto;
  cursor: pointer;

  transition: min-width ${soloConstants.transitionTime};

  &:hover,
  &:focus {
    color: white;
    background: ${colors.puddleBlue};
    outline: none;
  }

  &:active {
    background: ${colors.seaBlue};
  }

  &:disabled {
    opacity: 0.3;
    pointer-events: none;
    cursor: default;
    background: ${colors.seaBlue};
    color: white;
  }

  > span {
    position: relative;
    z-index: 2;
  }
  background: #b05464;

  &:hover,
  &:focus {
    background: ${colors.peachOrange};
  }

  &:active {
    background: ${colors.grapefruitOrange};
  }

  &:disabled {
    background: ${colors.grapefruitOrange};
  }

  ${(props: SoloButtonProps) => {
    return props.inProgress
      ? `
        min-width: 100%;
        pointer-events: none;
        cursor: default; 
        
        > div {
          display: block;
        }`
      : '';
  }};
`;
