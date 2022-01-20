import styled from '@emotion/styled/macro';
import * as React from 'react';
import { colors } from 'Styles/colors';
import { soloConstants } from 'Styles/StyledComponents/button';

export const Label = styled.label`
  display: block;
  color: ${colors.novemberGrey};
  font-size: 16px;
  margin-bottom: 10px;
  font-weight: 500;
`;

type TextareaStyleProps = {
  error?: boolean;
  borderless?: boolean;
  hideBottomMargin?: boolean;
};

const Textarea = styled.textarea`
  width: 100%;
  padding: 9px 15px 9px 11px;
  border: 1px solid ${colors.aprilGrey};
  border-radius: ${soloConstants.smallRadius}px;
  margin-bottom: ${(props: TextareaStyleProps) =>
    props.hideBottomMargin ? '0' : '15px'};
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

export interface TextareaProps
  extends Partial<
    React.DetailedHTMLProps<
      React.TextareaHTMLAttributes<HTMLTextAreaElement>,
      HTMLTextAreaElement
    >
  > {
  name?: string;
  title?: string;
  placeholder?: string;
  value: string | number;
  rows?: number;
  hideBottomMargin?: boolean;
  disabled?: boolean;
  onChange: (e: React.ChangeEvent<HTMLTextAreaElement>) => any;
  onBlur?: (e: React.ChangeEvent<HTMLTextAreaElement>) => any;
  borderless?: boolean;
  error?: boolean;
  requiredField?: {
    satisfied: boolean;
  };
  testId?: string;
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
    borderless,
    requiredField,
    testId,
    hideBottomMargin,
    ...rest
  } = props;

  return (
    <div>
      {title && (
        <Label>
          {title}
          {!!requiredField && (
            <span className={requiredField.satisfied ? '' : 'text-red-700'}>
              {' '}
              *
            </span>
          )}
        </Label>
      )}
      <Textarea
        data-testid={testId}
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
        hideBottomMargin={hideBottomMargin}
        {...rest}
      />
    </div>
  );
};
