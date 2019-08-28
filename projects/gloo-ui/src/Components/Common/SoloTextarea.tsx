import styled from '@emotion/styled';
import * as React from 'react';
import { colors, soloConstants } from 'Styles';

export const Label = styled.label`
  display: block;
  color: ${colors.novemberGrey};
  font-size: 16px;
  margin-bottom: 10px;
  font-weight: 500px;
`;
type TextareaStyleProps = { error?: boolean; borderless?: boolean };
const Textarea = styled.textarea`
  width: 100%;
  padding: 9px 15px 9px 11px;
  border: 1px solid ${colors.aprilGrey};
  border-radius: ${soloConstants.smallRadius}px;
  margin-bottom: 15px;
  line-height: 16px;
  outline: none;

  border: 1px solid
    ${(props: TextareaStyleProps) =>
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

interface TextareaProps {
  name?: string;
  title?: string;
  placeholder?: string;
  value: string | number;
  rows?: number;
  disabled?: boolean;
  onChange: (e: React.ChangeEvent<HTMLTextAreaElement>) => any;
  onBlur?: (e: React.ChangeEvent<HTMLTextAreaElement>) => any;
  borderless?: boolean;
  error?: boolean;
}

export const SoloTextarea = (props: TextareaProps) => {
  const {
    name,
    title,
    placeholder,
    value,
    rows,
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
      <Textarea
        borderless={borderless}
        name={name}
        rows={rows || 5}
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
