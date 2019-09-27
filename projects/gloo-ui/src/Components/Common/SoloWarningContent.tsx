import { INVALID_LICENSE_ERROR_ID } from '../../store/config/actions';
import * as React from 'react';
import { Modal } from 'antd';
const { warning } = Modal;

export interface SoloWarningContentProps {
  content?: string;
}

export const SoloWarningContent = (
  props: SoloWarningContentProps
): React.ReactNode => {
  const { content } = props;

  switch (content) {
    case INVALID_LICENSE_ERROR_ID:
      return (
        <>
          This feature requires an Enterprise Gloo license. <br />
          <a href='http://www.solo.io/gloo-trial'>
            Click here to request a trial license
          </a>
          .
        </>
      );
    default:
      return <>{content}</>;
  }
};

export const SoloWarning = (title: string, error: Error): void => {
  warning({
    title: title,
    content: SoloWarningContent({ content: error.message })
  });
};
