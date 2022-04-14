import { useGetGraphqlApiDetails, usePageApiRef } from 'API/hooks';
import React from 'react';
import { isExecutableAPI } from 'utils/graphql-helpers';
import GraphqlDefineResolversPrompt from '../api-policies/GraphqlDefineResolversPrompt';
import { ExecutableGraphqlApiDetails } from './executable-api/ExecutableGraphqlApiDetails';
import GraphqlApiConfigurationHeader from './GraphqlApiConfigurationHeader';
import { StitchedGraphqlApiDetails } from './stitched-api/StitchedGraphqlApiDetails';

const GraphqlApiDetails = () => {
  const apiRef = usePageApiRef();
  const { data: graphqlApi } = useGetGraphqlApiDetails(apiRef);

  if (!graphqlApi) return null;
  return (
    <>
      <GraphqlApiConfigurationHeader apiRef={apiRef} />

      {isExecutableAPI(graphqlApi) ? (
        <>
          <GraphqlDefineResolversPrompt apiRef={apiRef} />
          <ExecutableGraphqlApiDetails apiRef={apiRef} />
        </>
      ) : (
        <StitchedGraphqlApiDetails apiRef={apiRef} />
      )}
    </>
  );
};

export default GraphqlApiDetails;
