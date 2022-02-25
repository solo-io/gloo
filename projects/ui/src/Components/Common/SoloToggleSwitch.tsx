import styled from '@emotion/styled/macro';
import * as React from 'react';
import { colors } from 'Styles/colors';
import { Label as InputLabel } from './SoloInput';

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

export interface SoloToggleSwitchProps {
  checked: boolean;
  onChange: (checked: boolean) => void;
  checkedLabel?: string;
  uncheckedLabel?: string;
  labelOnRight?: boolean;
  small?: boolean;
  disabled?: boolean;
  label?: string;
}

export const SoloToggleSwitch = (props: SoloToggleSwitchProps) => {
  const {
    checked,
    checkedLabel,
    uncheckedLabel,
    onChange,
    labelOnRight,
    small,
    disabled,
    label,
  } = props;

  const onClick = (): void => {
    onChange(!checked);
  };
  let labelSection: React.ReactNode = null;

  if (!!label) {
    labelSection = <InputLabel>{label}</InputLabel>;
  } else if (!!checkedLabel && checked) {
    labelSection = (
      <Label small={false} labelOnRight={labelOnRight}>
        {checkedLabel}
      </Label>
    );
  } else if (!!uncheckedLabel && !checked) {
    labelSection = (
      <Label small={false} labelOnRight={labelOnRight}>
        {uncheckedLabel}
      </Label>
    );
  }

  return (
    <>
      {!labelOnRight && labelSection}
      <button
        type='button'
        disabled={disabled}
        onClick={onClick}
        aria-pressed='false'
        className={`
        ${disabled ? 'opacity-60 cursor-not-allowed' : 'opacity-100'}
        ${
          checked ? ' bg-blue-500gloo' : 'bg-gray-200'
        }  relative inline-flex flex-shrink-0 h-6 w-11 border-2 border-transparent rounded-full cursor-pointer transition-colors ease-in-out duration-200 focus:outline-none `}>
        <span className='sr-only'>{labelSection}</span>
        <span
          aria-hidden='true'
          className={`pointer-events-none ${
            checked ? ' translate-x-5' : 'translate-x-0 '
          } inline-block h-5 w-5 rounded-full bg-white shadow transform ring-0 transition ease-in-out duration-200`}></span>
      </button>
      {!!labelOnRight && labelSection}
    </>
  );
};
