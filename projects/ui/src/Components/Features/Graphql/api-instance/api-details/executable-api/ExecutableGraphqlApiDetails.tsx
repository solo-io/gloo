import { ClusterObjectRef } from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';
import React from 'react';
import ExecutableGraphqlUpstreamsTable from './upstreams/ExecutableGraphqlUpstreamsTable';
import ExecutableGraphqlSchemaDefinitions from './schema/ExecutableGraphqlSchemaDefinitions';

export const ExecutableGraphqlApiDetails: React.FC<{
  apiRef: ClusterObjectRef.AsObject;
}> = ({ apiRef }) => {
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
    </>
  );
};
