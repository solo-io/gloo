import * as React from 'react';
/** @jsx jsx */
import { css, jsx, SerializedStyles } from '@emotion/core';
import {
  SoloButtonCSS,
  ButtonProgress,
  SoloButtonStyledComponent
} from 'Styles/CommonEmotions/button';

import { BaseButtonProps } from 'antd/lib/button/button';

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
    <React.Fragment>
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
    </React.Fragment>
  );
};
