import * as React from 'react';
/** @jsx jsx */
import { css, jsx, SerializedStyles } from '@emotion/core';
import { Button } from 'antd';
import { SoloButtonCSS } from 'Styles/CommonEmotions/button';

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
      <Button
        {...rest}
        loading={loading}
        css={css`
          ${SoloButtonCSS};
          ${uniqueCss || {}};
        `}
        onClick={onClick}>
        {loading ? inProgressText : error ? errorText : text}
      </Button>
    </React.Fragment>
  );
};
