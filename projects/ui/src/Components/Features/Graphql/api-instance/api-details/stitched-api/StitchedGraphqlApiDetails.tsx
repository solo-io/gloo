import {
  useGetConsoleOptions,
  useGetGraphqlApiDetails,
  useGetStitchedSchemaDefinition,
} from 'API/hooks';
import AreaHeader from 'Components/Common/AreaHeader';
import { ClusterObjectRef } from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';
import React, { useMemo } from 'react';
import { parseSchemaString } from 'utils/graphql-helpers';
import GraphqlDeleteApiButton from '../GraphqlDeleteApiButton';
import SchemaDefinitions from '../schema/SchemaDefinitions';
import StitchedGqlSubGraphs from './sub-graphs/StitchedGqlSubGraphs';

export const StitchedGraphqlApiDetails: React.FC<{
  apiRef: ClusterObjectRef.AsObject;
}> = ({ apiRef }) => {
  const { readonly } = useGetConsoleOptions();

  // -- SUBGRAPHS -- //
  const { data: graphqlApi } = useGetGraphqlApiDetails(apiRef);
  const subschemasList =
    graphqlApi?.spec?.stitchedSchema?.subschemasList ??
    [
      // * Uncomment this for fake data.
      // * The bookinfo-graphql sub-graph should have a
      // * working link to the bookinfo example.
      // {
      //   name: 'bookinfo-graphql',
      //   namespace: 'gloo-system',
      //   typeMergeMap: [],
      // },
      // {
      //   name: 'test-book',
      //   namespace: 'gloo-system',
      //   typeMergeMap: [],
      // },
    ];

  const { data: stitchedSchema } = useGetStitchedSchemaDefinition(apiRef);
  const parsedStitchedSchema = useMemo(() => {
    if (!stitchedSchema) return undefined;
    return parseSchemaString(stitchedSchema);
  }, [stitchedSchema]);

  return (
    <>
      <StitchedGqlSubGraphs apiRef={apiRef} subGraphs={subschemasList} />

      <hr className='mt-10 mb-10' />

      <div className='mb-10'>
        {/* <div className='text-lg mb-5'>Stitched Schema</div> */}
        <AreaHeader
          title='Stitched Schema'
          contentTitle={`${apiRef.namespace}--${apiRef.name}-gql.txt`}
          yaml={stitchedSchema}
          onLoadContent={() => new Promise((res, rej) => res('loaded!'))}
        />

        {parsedStitchedSchema && (
          <SchemaDefinitions
            isEditable={false}
            schema={parsedStitchedSchema}
            apiRef={apiRef}
          />
        )}
      </div>

      {!readonly && <GraphqlDeleteApiButton apiRef={apiRef} />}
    </>
  );
};
