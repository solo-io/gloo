import * as React from 'react';
import styled from '@emotion/styled';
import { colors } from 'Styles/colors';
import { CloseOutlined } from '@ant-design/icons';

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

export const Input = styled.input<InputStyledProps>`
  width: 100%;
  padding: 7px 15px 7px 11px;
  border-radius: 8px;
  line-height: 16px;
  outline: none;
  padding-right: 41px;

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

const ClearSearchButton = styled.button`
  position: absolute;
  right: 0px;
  top: 0px;
  bottom: 0px;
  width: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.1rem;
  border-top-right-radius: 8px;
  border-bottom-right-radius: 8px;
  background-color: ${colors.dropBlue};
  color: ${colors.mayGrey};
  border: 1px solid ${colors.mayGrey};
  &:hover {
    color: ${colors.pondBlue};
    background-color: ${colors.splashBlue};
  }
  &:active {
    color: ${colors.seaBlue};
    background-color: ${colors.puddleBlue};
  }
`;
const InputContainer = styled.div`
  position: relative;
`;

export interface InputProps
  extends Partial<
    React.DetailedHTMLProps<
      React.InputHTMLAttributes<HTMLInputElement>,
      HTMLInputElement
    >
  > {
  name?: string;
  title?: string;
  label?: React.ReactNode;
  placeholder?: string;
  value: string | number;
  disabled?: boolean;
  onChange?: (e: React.ChangeEvent<HTMLInputElement>) => any;
  onBlur?: (e: React.ChangeEvent<HTMLInputElement>) => any;
  onKeyPress?: (e: React.KeyboardEvent) => void;
  borderless?: boolean;
  error?: boolean;
  password?: boolean;
  file?: boolean;
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
    file,
    label,
    ...rest
  } = props;

  let type = 'text';
  if (!!password) {
    type = 'password';
  }
  if (!!file) {
    type = 'file';
  }

  return (
    <div>
      {title && !label && <Label>{title}</Label>}
      {label && <Label>{label}</Label>}
      <InputContainer>
        <Input
          type={type}
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
          {...rest}
        />
        {value !== '' && (
          <ClearSearchButton
            onClick={e => {
              if (onChange)
                onChange({
                  ...e,
                  target: { ...e.target, name, value: '' },
                  currentTarget: { ...e.currentTarget, name, value: '' },
                } as any);
            }}>
            <CloseOutlined />
          </ClearSearchButton>
        )}
      </InputContainer>
    </div>
  );
};
