import React, { useState, useEffect } from 'react';
import styled from '@emotion/styled/macro';
import { ReactComponent as ErrorIcon } from 'assets/big-unsuccessful-x.svg';
import { ReactComponent as ConnectionIssueImage } from 'assets/connection-error-graphic.svg';
import { ServiceError } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/gateway_resources_pb_service';
import { colors } from 'Styles/colors';

type DataErrorContainerProps = {
  center: boolean;
  loadingComplete: boolean;
};
const DataErrorContainer = styled.div<DataErrorContainerProps>`
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 100%;
  padding: 20px 10px 35px;
  opacity: 0;
  transition: opacity 0.3s;

  ${(props: DataErrorContainerProps) =>
    props.loadingComplete ? 'opacity: 1;' : ''}

  ${(props: DataErrorContainerProps) =>
    props.center
      ? `
          height: 100%;
          justify-content: center;
        `
      : ''}
`;

type ImageHolderProps = {
  calm: boolean;
};
const ImageHolder = styled.div<ImageHolderProps>`
  font-weight: 600;
  font-size: 18px;
  text-align: center;

  ${(props: ImageHolderProps) =>
    props.calm
      ? `
      svg {
        display: block;
      width: auto;
      height: 55px;
      margin-bottom: 10px;
    }`
      : `
  svg {
    width: 30px;

    line {
      stroke: white;
      stroke-width: 6px;
    }

    circle {
      stroke: none;
      fill: ${colors.pumpkinOrange};
    }
  }`};
`;

const MessageBox = styled.div`
  font-weight: 600;
  font-size: 18px;
  margin-top: 12px;
  color: ${colors.grapefruitOrange};
`;
const CalmerMessageBox = styled.div`
  margin-top: 10px;
  font-size: 14px;
  color: ${colors.juneGrey};
`;
const CodeBox = styled(CalmerMessageBox)`
  font-style: italic;
`;

interface ErrorProps {
  error: any /* TODO ServiceError; */;
  center?: boolean;
}

export const DataError = ({ center = true, error }: ErrorProps) => {
  const [gracefulStartComplete, setGracefulStateComplete] = useState(false);

  useEffect(() => {
    const timer = setTimeout(() => {
      setGracefulStateComplete(true);
    }, 1000);

    return () => clearTimeout(timer);
  }, []);

  return (
    <DataErrorContainer center={center} loadingComplete={gracefulStartComplete}>
      {error.code === 2 || error.code === 14 || error.code === 15 ? (
        <>
          <ImageHolder calm={true}>
            <ConnectionIssueImage />
            Connection Error
          </ImageHolder>
          <CalmerMessageBox>{error.message}</CalmerMessageBox>
        </>
      ) : (
        <>
          <ImageHolder calm={false}>
            <ErrorIcon />
          </ImageHolder>
          <MessageBox>{error.message}</MessageBox>
        </>
      )}

      {error.code !== undefined && <CodeBox>Code: {error.code}</CodeBox>}
    </DataErrorContainer>
  );
};
