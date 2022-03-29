import { ClusterObjectRef } from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';
import React from 'react';
import { ExecutableGraphqlApiDetails } from './executable-api/ExecutableGraphqlApiDetails';
import GraphqlApiConfigurationHeader from './GraphqlApiConfigurationHeader';
import GraphqlDeleteApiButton from './GraphqlDeleteApiButton';
import GraphqlEditApiButton from './GraphqlEditApiButton';
import styled from '@emotion/styled/macro';
import { useGetConsoleOptions } from 'API/hooks';

const ButtonContainer = styled.div`
  button {
    margin-right: 10px;
  }
`;

const GraphqlApiDetails: React.FC<{ apiRef: ClusterObjectRef.AsObject }> = ({
  apiRef,
}) => {
  const { readonly } = useGetConsoleOptions();

  return (
    <>
      <GraphqlApiConfigurationHeader apiRef={apiRef} />
      {/*
      // ! Check for api type here (gateway or executable)
      <GatewayGraphqlApiDetails apiRef={apiRef} />
      */}
      <ExecutableGraphqlApiDetails apiRef={apiRef} />
      {!readonly && (
        <ButtonContainer>
          <GraphqlEditApiButton apiRef={apiRef} />
          <GraphqlDeleteApiButton apiRef={apiRef} />
        </ButtonContainer>
      )}
    </>
  );
};

export default GraphqlApiDetails;
