import * as React from 'react';
/** @jsx jsx */
import { css, jsx } from '@emotion/core';
import { colors, soloConstants } from '../../Styles';
import styled from '@emotion/styled/macro';
import { Checkbox } from 'antd';
import { CheckboxChangeEvent } from 'antd/lib/checkbox';
import { Label } from './SoloInput';

const InputStyling = css`
  border-radius: 10px;
  width: 190px;
  padding: 10px 16px;
  margin-bottom: 20px;
  background: white;
  border: 1px solid ${colors.juneGrey};
`;

const CheckboxStyling = css`
  .ant-checkbox {
    .ant-checkbox-inner {
      background: ${colors.januaryGrey};
      border: 1px solid ${colors.juneGrey};
      border-radius: 5px;
      width: 18px;
      height: 18px;
    }

    &.ant-checkbox-checked {
      .ant-checkbox-inner {
        background: ${colors.puddleBlue};
        border-color: ${colors.seaBlue};

        &:after {
          border-color: ${colors.seaBlue};
          border-width: 1px;
          transform: rotate(45deg) scale(1) translate(-37%, -66%);
          height: 9px;
        }
      }
    }
  }
`;

const OnlyCheckbox = styled.span`
${CheckboxStyling}
color: ${colors.septemberGrey};
`;

const CheckboxWrapper = styled<'div', { checked?: boolean }>('div')`
  ${InputStyling}
  display: flex;
  justify-content: space-between;
  padding: 7px 7px 7px 16px;
  color: ${colors.septemberGrey};
  transition: background ${soloConstants.transitionTime},
    border ${soloConstants.transitionTime};

  ${props =>
    !!props.checked
      ? `background: ${colors.dropBlue};
        border-color: ${colors.seaBlue};`
      : ``}

  ${CheckboxStyling};
`;

export interface CheckboxProps {
  checked: boolean;
  disabled?: boolean;
  onChange: (e: CheckboxChangeEvent) => void;
  title?: string;
  withWrapper?: boolean;
  label?: boolean;
}

export const SoloCheckbox: React.FC<CheckboxProps> = props => {
  const { title, checked, onChange, withWrapper, label, disabled } = props;

  if (!!withWrapper) {
    return (
      <CheckboxWrapper checked={checked}>
        {label ? <Label>{title}</Label> : title}
        <Checkbox disabled={disabled} checked={checked} onChange={onChange} />
      </CheckboxWrapper>
    );
  }

  return (
    <OnlyCheckbox>
      {label ? <Label>{title}</Label> : title}
      <Checkbox disabled={disabled} checked={checked} onChange={onChange} />
    </OnlyCheckbox>
  );
};
