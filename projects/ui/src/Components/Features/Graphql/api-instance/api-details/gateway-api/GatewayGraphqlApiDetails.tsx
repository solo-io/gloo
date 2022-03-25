import { TabPanels, Tabs } from '@reach/tabs';
import {
  FolderTab,
  FolderTabContent,
  FolderTabList,
  StyledTabPanel,
} from 'Components/Common/Tabs';
import styled from '@emotion/styled/macro';
import { ClusterObjectRef } from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';
import React, { useState } from 'react';
import GatewayGraphqlMutationsTable from './schema/GatewayGraphqlMutationsTable';
import GatewayGraphqlObjectsTable from './schema/GatewayGraphqlObjectsTable';
import GatewayGraphqlQueriesTable from './schema/GatewayGraphqlQueriesTable';
import GatewayGraphqlSubGraphs from './sub-graphs/GatewayGraphqlSubGraphs';

export const GatewayGraphqlApiDetails: React.FC<{
  apiRef: ClusterObjectRef.AsObject;
}> = ({ apiRef }) => {
  const [schemaTabIndex, setSchemaTabIndex] = useState(0);

  return (
    <>
      <GatewayGraphqlSubGraphs apiRef={apiRef} />

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
                {schemaTabIndex === 0 && <GatewayGraphqlQueriesTable />}
              </FolderTabContent>
            </StyledTabPanel>
            <StyledTabPanel>
              <FolderTabContent>
                {schemaTabIndex === 1 && <GatewayGraphqlMutationsTable />}
              </FolderTabContent>
            </StyledTabPanel>
            <StyledTabPanel>
              <FolderTabContent>
                {schemaTabIndex === 2 && <GatewayGraphqlObjectsTable />}
              </FolderTabContent>
            </StyledTabPanel>
          </TabPanels>
        </Tabs>
      </div>
      {/* </Card> */}
    </>
  );
};
