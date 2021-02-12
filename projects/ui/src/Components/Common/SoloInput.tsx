import * as React from 'react';
import styled from '@emotion/styled';
import { colors } from 'Styles/colors';

export const Label = styled.label`
  display: block;
  color: ${colors.novemberGrey};
  font-size: 16px;
  margin-bottom: 10px;
  font-weight: 600;
`;

type InputStyledProps = {
  error?: boolean;
  borderless?: boolean;
};

const Input = styled.input<InputStyledProps>`
  width: 100%;
  padding: 7px 15px 7px 11px;
  border-radius: 8px;
  line-height: 16px;
  outline: none;

  border: 1px solid
    ${(props: InputStyledProps) =>
      props.error
        ? colors.grapefruitOrange
        : props.borderless
        ? 'none'
        : colors.aprilGrey};

  &:disabled {
    background: ${colors.februaryGrey};
    color: ${colors.septemberGrey};
    border: 1px solid ${colors.mayGrey};
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
  onChange?: (e: React.ChangeEvent<HTMLInputElement>) => any;
  onBlur?: (e: React.ChangeEvent<HTMLInputElement>) => any;
  onKeyPress?: (e: React.KeyboardEvent) => void;
  borderless?: boolean;
  error?: boolean;
  password?: boolean;
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
    borderless,
    onKeyPress,
    password,
  } = props;

  return (
    <div>
      {title && <Label>{title}</Label>}
      <Input
        type={!!password ? 'password' : 'text'}
        borderless={borderless}
        name={name}
        placeholder={placeholder}
        title={title}
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
