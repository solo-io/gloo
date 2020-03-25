import * as React from 'react';
import styled from '@emotion/styled';
import { Switch } from 'antd';
import { colors } from 'Styles';

type ToggleSwitchContainerProps = {
  small?: boolean;
};
const ToggleSwitchContainer = styled.div`
  display: flex;
  align-items: center;
  height: 29px;
  cursor: pointer;

  .ant-switch {
    &::after {
      top: 2px;
    }

    &.ant-switch-checked {
      background: ${colors.seaBlue};

      &::after {
        left: 100%;
      }
    }
  }
  ${(props: ToggleSwitchContainerProps) => {
    return props.small
      ? `
        .ant-switch {
          height: 20px;
          width: 35px;

          &::after {
            width: 14px;
            height: 14px;
            left: 3px;
          }

          &.ant-switch-checked::after {
              margin-left: -3px;
            }
        }
        `
      : `
        .ant-switch {
          height: 29px;
          width: 58px;

          &::after {
            width: 23px;
            height: 23px;
            left: 2px;
          }

          &.ant-switch-checked::after {
              margin-left: -2px;
          }
        }
      `;
  }};
`;

type ToggleLabelProps = {
  small?: boolean;
  labelOnRight?: boolean;
};
const Label = styled.div`
  color: ${colors.novemberGrey};
  font-size: ${(props: ToggleLabelProps) => (props.small ? '18px' : '22px')};
  ${(props: ToggleLabelProps) =>
    props.labelOnRight ? 'margin-left: 7px;' : 'margin-right: 7px;'};
`;

interface ToggleSwitchProps {
  checked: boolean;
  onChange: (checked: boolean) => void;
  checkedLabel?: string;
  uncheckedLabel?: string;
  labelOnRight?: boolean;
  small?: boolean;
}

export const SoloToggleSwitch = (props: ToggleSwitchProps) => {
  const {
    checked,
    checkedLabel,
    uncheckedLabel,
    onChange,
    labelOnRight,
    small
  } = props;

  const onClick = (): void => {
    onChange(!checked);
  };

  const labelSection =
    !!checkedLabel && checked ? (
      <Label small={small} labelOnRight={labelOnRight}>
        {checkedLabel}
      </Label>
    ) : !!uncheckedLabel && !checked ? (
      <Label small={small} labelOnRight={labelOnRight}>
        {uncheckedLabel}
      </Label>
    ) : null;

  return (
    <ToggleSwitchContainer onClick={onClick} small={small}>
      {!labelOnRight && labelSection}
      <Switch checked={checked} />
      {!!labelOnRight && labelSection}
    </ToggleSwitchContainer>
  );
};
