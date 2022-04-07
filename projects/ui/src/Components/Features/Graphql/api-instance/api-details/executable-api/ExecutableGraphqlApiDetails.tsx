import { ClusterObjectRef } from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';
import React, { useMemo } from 'react';
import ExecutableGraphqlUpstreamsTable from './upstreams/ExecutableGraphqlUpstreamsTable';
import SchemaDefinitions from '../schema/SchemaDefinitions';
import styled from '@emotion/styled';
import GraphqlDeleteApiButton from '../GraphqlDeleteApiButton';
import GraphqlEditApiButton from './GraphqlEditApiButton';
import { useGetConsoleOptions, useGetGraphqlApiDetails } from 'API/hooks';
import { getParsedExecutableApiSchema } from 'utils/graphql-helpers';

const ButtonContainer = styled.div`
  button {
    margin-right: 10px;
  }
`;

export const ExecutableGraphqlApiDetails: React.FC<{
  apiRef: ClusterObjectRef.AsObject;
}> = ({ apiRef }) => {
  const { readonly } = useGetConsoleOptions();

  const { data: graphqlApi } = useGetGraphqlApiDetails(apiRef);
  const schema = useMemo(
    () => getParsedExecutableApiSchema(graphqlApi),
    [graphqlApi]
  );

  return (
    <>
      <div className='mb-10'>
        <div className='text-lg mb-5'>Schema</div>
        <SchemaDefinitions schema={schema} apiRef={apiRef} isEditable={true} />
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
