import React from 'react';
import { ReactComponent as LoadingHexagon } from 'assets/Loader.svg';
import styled from '@emotion/styled';
import { keyframes } from '@emotion/core';

type RotaterProps = {
  degrees: number;
};
const Rotater = styled.div`
  transform: rotate(${(props: RotaterProps) => props.degrees}deg);
  width: 141px;
  height: 141px;

  svg {
    width: 141px;
    height: 141px;
  }
`;

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
  const [degrees, setDegrees] = React.useState(0);

  React.useEffect(() => {
    const spinterval = setInterval(() => {
      setDegrees(oldDegree => (oldDegree + 60) % 360);
    }, 200);

    return () => clearInterval(spinterval);
  }, []);

  const centering = center
    ? {
        display: 'flex',
        justifyContent: 'center',
        alignContent: 'center'
      }
    : {};

  return (
    <div
      style={{
        width: '100%',
        ...centering
      }}>
      <div style={{ width: '141px', height: '141px' }}>
        <Rotater degrees={degrees}>
          <LoadingHexagon />
        </Rotater>
      </div>
    </div>
  );
};
