import { ClusterObjectRef } from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';
import React from 'react';
import ExecutableGraphqlUpstreamsTable from './upstreams/ExecutableGraphqlUpstreamsTable';
import ExecutableGraphqlSchemaDefinitions from './schema/ExecutableGraphqlSchemaDefinitions';
import styled from '@emotion/styled';
import GraphqlDeleteApiButton from '../GraphqlDeleteApiButton';
import GraphqlEditApiButton from './GraphqlEditApiButton';
import { useGetConsoleOptions } from 'API/hooks';

const ButtonContainer = styled.div`
  button {
    margin-right: 10px;
  }
`;

export const ExecutableGraphqlApiDetails: React.FC<{
  apiRef: ClusterObjectRef.AsObject;
}> = ({ apiRef }) => {
  const { readonly } = useGetConsoleOptions();

  return (
    <>
      <div className='mb-10'>
        <div className='text-lg mb-5'>Schema</div>
        <ExecutableGraphqlSchemaDefinitions apiRef={apiRef} />
      </div>

      <div className='mb-10'>
        <div className='text-lg mb-5'>Upstreams</div>
        <ExecutableGraphqlUpstreamsTable apiRef={apiRef} />
      </div>

      {!readonly && (
        <ButtonContainer>
          <GraphqlEditApiButton apiRef={apiRef} />
          <GraphqlDeleteApiButton apiRef={apiRef} />
        </ButtonContainer>
      )}
    </>
  );
};
