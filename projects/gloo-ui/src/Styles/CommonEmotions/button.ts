import { css, keyframes } from '@emotion/core';
import { colors } from '../colors';
import styled from '@emotion/styled/macro';
import { soloConstants } from 'Styles/constants';
import blueProgressMask from 'assets/primary-progress-mask.svg';
import orangeProgressMask from 'assets/warning-progress-mask.svg';
import { Button } from 'antd';

const slide = keyframes`
  from { background-position: -281px 0; }
  to { background-position: 0 0; }
`;

export const SoloButtonCSS = css`
  position: relative;
  display: inline-block;
  padding: 10px 20px 13px;
  font-size: 16px;
  line-height: 14px;
  background: ${colors.seaBlue};
  color: white;
  border: none;
  outline: none;
  border-radius: ${soloConstants.smallRadius}px;
  min-width: 125px;
  cursor: pointer;

  transition: min-width ${soloConstants.transitionTime};

  &:hover,
  &:focus {
    color: white;
    background: ${colors.lakeBlue};
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
export const SoloButtonStyledComponent = styled(Button)`
  ${SoloButtonCSS};

  ${props =>
    // @ts-ignore
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
  ${SoloButtonCSS};
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
  ${SoloButtonCSS};
  background: ${colors.grapefruitOrange};

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

  ${props => {
    // @ts-ignore
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

export const InverseButtonCSS = css`
  ${SoloButtonCSS};
  background: white;
  color: ${colors.seaBlue};

  &:hover,
  &:focus {
    color: ${colors.lakeBlue};
    background: white;
  }

  &:active {
    background: ${colors.seaBlue};
    background: white;
  }
`;

export const SmallButtonCSS = css`
  padding: 7px 15px;
`;

export const IconButton = styled.button`
  display: inline-flex;
  cursor: pointer;
  border: none;
  outline: none;
  background: transparent;
  justify-content: center;
  align-items: center;

  &:disabled {
    opacity: 0.3;
    pointer-events: none;
    cursor: default;
  }
`;

export const ConfigurationButtonEmo = css`
  ${SoloButtonCSS};
  ${SmallButtonCSS};
  height: 43px;
  line-height: 29px;
  padding: 7px 20px;
  display: flex;
  margin-top: 15px;
  width: 280px;
  transition: all ${soloConstants.transitionTime};
`;
export const ConfigurationButtonDisabledEmo = css`
  ${ConfigurationButtonEmo};
  background: ${colors.januaryGrey};
  cursor: default;

  &:hover,
  &:focus,
  &:active {
    background: ${colors.januaryGrey};
    color: white;
  }
`;
