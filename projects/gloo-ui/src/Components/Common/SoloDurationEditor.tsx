import * as React from 'react';
/** @jsx jsx */
import { css, jsx } from '@emotion/core';
import { colors, soloConstants } from '../../Styles';
import styled from '@emotion/styled/macro';
import { Duration } from 'google-protobuf/google/protobuf/duration_pb';

export const Label = styled.label`
  display: block;
  color: ${colors.novemberGrey};
  font-size: 16px;
  margin-bottom: 10px;
  font-weight: 500;
`;

const InputHolder = styled<'div', { leftSide?: boolean }>('div')`
  display: inline-block;
  width: calc(50% - 4px);
  ${props => props.leftSide && 'margin-right: 8px;'};
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

export interface DurationProps {
  title?: string;
  value: Duration.AsObject | undefined;
  disabled?: boolean;
  onChange?: (newDuration: Duration.AsObject) => any;
  onBlur?: (newDuration: Duration.AsObject) => any;
  borderless?: boolean;
  error?: boolean;
}

export const SoloDurationEditor = (props: DurationProps) => {
  const { title, value, onChange, onBlur, disabled, error, borderless } = props;

  const onChangeSeconds = (e: React.ChangeEvent<HTMLInputElement>): void => {
    const seconds = parseInt(e.target.value);
    if (!!seconds) {
      onChange!({ nanos: value ? value.nanos : 0, seconds });
    }
  };
  const onChangeNanos = (e: React.ChangeEvent<HTMLInputElement>): void => {
    const nanos = parseInt(e.target.value);
    if (!!nanos) {
      onChange!({ seconds: value ? value.seconds : 0, nanos });
    }
  };
  const onBlurSeconds = (e: React.ChangeEvent<HTMLInputElement>): void => {
    if (!!onBlur) {
      const seconds = parseInt(e.target.value);
      if (!!seconds) {
        onBlur({ nanos: value ? value.nanos : 0, seconds });
      }
    }
  };
  const onBlurNanos = (e: React.ChangeEvent<HTMLInputElement>): void => {
    if (!!onBlur) {
      const nanos = parseInt(e.target.value);
      if (!!nanos) {
        onBlur({ seconds: value ? value.seconds : 0, nanos });
      }
    }
  };

  return (
    <div>
      {title && <Label>{title} [secs | nano]</Label>}
      <InputHolder leftSide>
        <Input
          borderless={borderless}
          placeholder={'##'}
          title={'Seconds'}
          type={'number'}
          value={value ? value.seconds : undefined}
          onChange={onChangeSeconds}
          onBlur={onBlurSeconds}
          disabled={disabled}
          error={error}
        />
      </InputHolder>
      <InputHolder>
        <Input
          borderless={borderless}
          placeholder={'##'}
          title={'Nanos'}
          type={'number'}
          value={value ? value.nanos : undefined}
          onChange={onChangeNanos}
          onBlur={onBlurNanos}
          disabled={disabled}
        />
      </InputHolder>
    </div>
  );
};
