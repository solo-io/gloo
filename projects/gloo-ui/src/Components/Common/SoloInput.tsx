import * as React from 'react';
/** @jsx jsx */
import { css, jsx } from '@emotion/core';
import { colors, soloConstants } from '../../Styles';
import styled from '@emotion/styled/macro';

export const Label = styled.label`
  display: block;
  color: ${colors.novemberGrey};
  font-size: 16px;
  margin-bottom: 10px;
  font-weight: 500;
`;

const Input = styled<'input', { error?: boolean; borderless?: boolean }>(
  'input'
)`
  width: 100%;
  padding: 9px 15px 9px 11px;
  border-radius: ${soloConstants.smallRadius}px;

  line-height: 16px;
  outline: none;

  border: 1px solid
    ${props =>
      props.error
        ? colors.grapefruitOrange
        : props.borderless
        ? 'none'
        : colors.aprilGrey};

  &:disabled {
    background: ${colors.aprilGrey};
  }

  &::placeholder {
    color: ${colors.juneGrey};
  }
`;

export interface InputProps {
  name?: string;
  title?: string;
  type?: string;
  placeholder?: string;
  value: string | number;
  disabled?: boolean;
  onChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
  onBlur?: (e: React.ChangeEvent<HTMLInputElement>) => void;
  borderless?: boolean;
  error?: boolean;
  onKeyPress?: (e: React.KeyboardEvent) => void;
}

export const SoloInput: React.FC<InputProps> = props => {
  const {
    name,
    title,
    placeholder,
    value,
    onChange,
    onBlur,
    disabled,
    error,
    borderless,
    type = 'text',
    onKeyPress
  } = props;

  return (
    <div>
      {title && <Label>{title}</Label>}
      <Input
        borderless={borderless}
        name={name}
        placeholder={placeholder}
        title={title}
        type={type}
        value={value}
        onChange={onChange}
        onBlur={onBlur}
        disabled={disabled}
        error={error}
        onKeyPress={onKeyPress}
      />
    </div>
  );
};
