import { Alert } from 'antd';
import * as React from 'react';

interface Props {
  className?: string;
  message: string;
}

export const WarningMessage: React.FC<Props> = ({ className, message }) => {
  return Boolean(message) ? (
    <Alert
      showIcon
      type='warning'
      className={!!className ? className : 'p-2 mb-3 mt-3'}
      message={' '}
      description={message}
    />
  ) : null;
};

export default WarningMessage;
