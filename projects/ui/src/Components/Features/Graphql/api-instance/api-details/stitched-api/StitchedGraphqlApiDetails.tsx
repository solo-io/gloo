import { TabPanels, Tabs } from '@reach/tabs';
import { useGetConsoleOptions, useGetGraphqlApiDetails } from 'API/hooks';
import {
  FolderTab,
  FolderTabContent,
  FolderTabList,
  StyledTabPanel,
} from 'Components/Common/Tabs';
import { ClusterObjectRef } from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';
import React, { useState } from 'react';
import GraphqlDeleteApiButton from '../GraphqlDeleteApiButton';
import StitchedGraphqlMutationsTable from './schema/StitchedGqlMutationsTable';
import StitchedGraphqlObjectsTable from './schema/StitchedGqlObjectsTable';
import StitchedGraphqlQueriesTable from './schema/StitchedGqlQueriesTable';
import StitchedGqlSubGraphs from './sub-graphs/StitchedGqlSubGraphs';

export const StitchedGraphqlApiDetails: React.FC<{
  apiRef: ClusterObjectRef.AsObject;
}> = ({ apiRef }) => {
  const [schemaTabIndex, setSchemaTabIndex] = useState(0);
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

  return (
    <>
      <StitchedGqlSubGraphs apiRef={apiRef} subGraphs={subschemasList} />

      <hr className='mt-10 mb-10' />

      {/* <Card className='shadow-[0_2px_5px_#8888] mb-5'> */}
      <div className='text-lg mb-5'>Schema</div>
      {/* <Collapse className='mb-5' defaultActiveKey={[0, 1, 2]}>
          <Collapse.Panel key={0} header='Mutations'>
            Mutations
          </Collapse.Panel>
          <Collapse.Panel key={1} header='Queries'>
            Queries
          </Collapse.Panel>
          <Collapse.Panel key={2} header='Objects'>
            Objects
          </Collapse.Panel>
        </Collapse> */}

      <div className='mb-5'>
        <Tabs index={schemaTabIndex} onChange={idx => setSchemaTabIndex(idx)}>
          <FolderTabList>
            <FolderTab>Queries</FolderTab>
            <FolderTab>Mutations</FolderTab>
            <FolderTab>Objects</FolderTab>
          </FolderTabList>

          <TabPanels>
            <StyledTabPanel>
              <FolderTabContent>
                {schemaTabIndex === 0 && <StitchedGraphqlQueriesTable />}
              </FolderTabContent>
            </StyledTabPanel>
            <StyledTabPanel>
              <FolderTabContent>
                {schemaTabIndex === 1 && <StitchedGraphqlMutationsTable />}
              </FolderTabContent>
            </StyledTabPanel>
            <StyledTabPanel>
              <FolderTabContent>
                {schemaTabIndex === 2 && <StitchedGraphqlObjectsTable />}
              </FolderTabContent>
            </StyledTabPanel>
          </TabPanels>
        </Tabs>
      </div>
      {/* </Card> */}

      {!readonly && <GraphqlDeleteApiButton apiRef={apiRef} />}
    </>
  );
};
