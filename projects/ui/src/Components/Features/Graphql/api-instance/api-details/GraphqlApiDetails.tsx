import { ClusterObjectRef } from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';
import React from 'react';
import { ExecutableGraphqlApiDetails } from './executable-api/ExecutableGraphqlApiDetails';
import GraphqlApiConfigurationHeader from './GraphqlApiConfigurationHeader';
import GraphqlDeleteApiButton from './GraphqlDeleteApiButton';

const GraphqlApiDetails: React.FC<{ apiRef: ClusterObjectRef.AsObject }> = ({
  apiRef,
}) => {
  return (
    <>
      <GraphqlApiConfigurationHeader apiRef={apiRef} />
      {/* 
      // ! Check for api type here (gateway or executable)
      <GatewayGraphqlApiDetails apiRef={apiRef} /> 
      */}
      <ExecutableGraphqlApiDetails apiRef={apiRef} />

      <GraphqlDeleteApiButton apiRef={apiRef} />
    </>
  );
};

export default GraphqlApiDetails;
