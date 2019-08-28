import { Spin } from 'antd';
import React from 'react';

interface LoadingProps {
  message?: string;
  children?: React.ReactChild;
  loading?: boolean;
  offset?: boolean;
  center?: boolean;
}

export const Loading = ({
  message,
  children,
  center = false,
  loading = true,
  offset = false
}: LoadingProps) => {
  const centering = center
    ? {
        display: 'flex',
        justifyContent: 'center',
        alignContent: 'center'
      }
    : {};

  return (
    <div style={{ width: '100%', ...centering }}>
      <Spin
        size='large'
        tip={message ? message : ''}
        spinning={loading}
        style={{ width: '100%', marginTop: `${offset ? '100px' : ''}` }}>
        {children}
      </Spin>
    </div>
  );
};
