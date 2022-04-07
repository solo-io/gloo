import { TabPanels, Tabs } from '@reach/tabs';
import {
  useGetGraphqlApiDetails,
  usePageApiRef,
  useGetConsoleOptions,
} from 'API/hooks';
import { ReactComponent as GraphQLIcon } from 'assets/graphql-icon.svg';
import { Loading } from 'Components/Common/Loading';
import { SectionCard } from 'Components/Common/SectionCard';
import {
  FolderTab,
  FolderTabContent,
  FolderTabContentNoPadding,
  FolderTabList,
  StyledTabPanel,
} from 'Components/Common/Tabs';
import React from 'react';
import { isExecutableAPI } from 'utils/graphql-helpers';
import { GraphqlIconHolder } from '../GraphqlLanding.style';
import GraphqlApiDetails from './api-details/GraphqlApiDetails';
import { GraphqlApiExplorer } from './api-explorer/GraphqlApiExplorer';
import GraphqlApiIntrospectionToggle from './GraphqlApiIntrospectionToggle';
import GraphqlDefineResolversPrompt from './GraphqlDefineResolversPrompt';

export const GraphqlInstance: React.FC = () => {
  const apiRef = usePageApiRef();
  const { data: graphqlApi } = useGetGraphqlApiDetails(apiRef);
  const { apiExplorerEnabled } = useGetConsoleOptions();

  // tab logic
  const [tabIndex, setTabIndex] = React.useState(0);
  const handleTabsChange = (index: number) => {
    setTabIndex(index);
  };

  if (!graphqlApi) return <Loading />;
  return (
    <div className='w-full mx-auto '>
      <SectionCard
        cardName={apiRef.name}
        logoIcon={<GraphqlIconHolder>{<GraphQLIcon />}</GraphqlIconHolder>}
        headerSecondaryInformation={[
          {
            title: 'Namespace',
            value: apiRef.namespace,
          },
        ]}>
        {isExecutableAPI(graphqlApi) && (
          <>
            <GraphqlDefineResolversPrompt apiRef={apiRef} />
            <div className='float-right'>
              <GraphqlApiIntrospectionToggle apiRef={apiRef} />
            </div>
          </>
        )}

        <Tabs index={tabIndex} onChange={handleTabsChange}>
          <FolderTabList>
            <FolderTab>API Details</FolderTab>
            {apiExplorerEnabled && <FolderTab>Explore</FolderTab>}
          </FolderTabList>

          <TabPanels>
            <StyledTabPanel>
              <FolderTabContent>
                {tabIndex === 0 && <GraphqlApiDetails />}
              </FolderTabContent>
            </StyledTabPanel>
            <StyledTabPanel>
              <FolderTabContentNoPadding>
                {tabIndex === 1 && <GraphqlApiExplorer />}
              </FolderTabContentNoPadding>
            </StyledTabPanel>
          </TabPanels>
        </Tabs>
      </SectionCard>
    </div>
  );
};
