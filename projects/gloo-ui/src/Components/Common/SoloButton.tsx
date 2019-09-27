import { css, SerializedStyles } from '@emotion/core';
import { BaseButtonProps } from 'antd/lib/button/button';
import * as React from 'react';
import {
  ButtonProgress,
  SoloButtonCSS,
  SoloButtonStyledComponent
} from 'Styles/CommonEmotions/button';

export interface SoloButtonProps extends BaseButtonProps {
  text: string;
  onClick: (e: React.MouseEvent<any, MouseEvent>) => void;
  inProgressText?: string;
  errorText?: string;
  error?: boolean;
  disabled?: boolean;
  uniqueCss?: SerializedStyles;
}

export const SoloButton = (props: SoloButtonProps) => {
  const {
    loading,
    text,
    onClick,
    inProgressText,
    errorText,
    error,
    uniqueCss,
    ...rest
  } = props;

  return (
    <>
      <SoloButtonStyledComponent
        {...rest}
        loading={loading}
        css={css`
          ${SoloButtonCSS};
          ${uniqueCss || {}};
        `}
        onClick={onClick}>
        <ButtonProgress />
        {loading ? inProgressText : error ? errorText : text}
      </SoloButtonStyledComponent>
    </>
  );
};
