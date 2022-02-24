import React from 'react';
import styled from '@emotion/styled/macro';

type RotaterProps = {
  degrees: number;
};
const Rotater = styled.div<RotaterProps>`
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

export const Loading = ({ center = true, message }: LoadingProps) => {
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
        height: '100%',
        justifyContent: 'center',
        alignContent: 'center',
      }
    : {};

  return (
    <div
      style={{
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        width: '100%',
        padding: '20px 10px 35px',
        ...centering,
      }}>
      <div>
        <Rotater degrees={degrees}>
          <svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 141.735 161.207'>
            <g opacity='.695'>
              <path
                d='M124.544 109.2l17.191 9.96v-76.8l-17.191 9.861z'
                fill='#b1e7fe'
              />
              <path
                fill='#b1e7ff'
                d='M19.108 112.701l-17.243 9.893 66.88 38.613v-19.85z'
              />
              <path
                d='M17.19 51.918l-17.191-9.96v76.8l17.191-9.861z'
                fill='#6ac7f0'
              />
              <path
                d='M68.742 19.85V0L2.289 38.367l17.161 9.942z'
                fill='#54b7e3'
              />
              <path
                d='M122.55 48.461l17.246-9.893-66.8-38.569v19.85z'
                fill='#2196c9'
              />
              <path
                fill='#b1e7ff'
                d='M72.993 141.356v19.85l66.536-38.413-17.161-9.942z'
              />
            </g>
          </svg>
        </Rotater>
      </div>
      <div
        style={{
          paddingTop: '16px',
          textAlign: 'center',
          fontSize: '18px',
          fontWeight: 500,
        }}>
        {message}
      </div>
    </div>
  );
};
