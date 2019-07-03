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
  font-weight: 500px;
`;

const Input = styled<'input', { error?: boolean; borderless?: boolean }>(
  'input'
)`
  width: 100%;
  padding: 9px 15px 9px 11px;
  border: 1px solid ${colors.aprilGrey};
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

interface InputProps {
  name?: string;
  title?: string;
  placeholder?: string;
  value: string | number;
  disabled?: boolean;
  onChange: (e: React.ChangeEvent<HTMLInputElement>) => any;
  onBlur?: (e: React.ChangeEvent<HTMLInputElement>) => any;
  borderless?: boolean;
  error?: boolean;
}

export const SoloInput = (props: InputProps) => {
  const {
    name,
    title,
    placeholder,
    value,
    onChange,
    onBlur,
    disabled,
    error,
    borderless
  } = props;

  return (
    <div>
      {title && <Label>{title}</Label>}
      {/*
                      // @ts-ignore*/}
      <Input
        borderless={borderless}
        name={name}
        placeholder={placeholder}
        title={title}
        value={value}
        onChange={onChange}
        onBlur={onBlur}
        disabled={disabled}
        error={error}
      />
    </div>
  );
};
