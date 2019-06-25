import * as React from 'react';
/** @jsx jsx */
import { jsx } from '@emotion/core';
import styled from '@emotion/styled/macro';
import { Popover } from 'antd';
import { colors } from '../../Styles';

interface PopoverProps {
  title?: string | React.ReactNode;
  children?: React.ReactNode;
  content: string | React.ReactNode;
}

export const SoloPopover: React.FunctionComponent<PopoverProps> = (
  props: PopoverProps
) => {
  const [open, setOpen] = React.useState<boolean>(false);
  const { title, content } = props;

  const handleVisibleChange = (visible: boolean) => {
    setOpen(visible);
  };

  return (
    <Popover
      content={content}
      trigger='click'
      title={title}
      visible={open}
      onVisibleChange={handleVisibleChange}>
      <div onClick={() => setOpen(true)}>{props.children}</div>
    </Popover>
  );
};
