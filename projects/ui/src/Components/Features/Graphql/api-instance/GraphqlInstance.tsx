import { TabPanels, Tabs } from '@reach/tabs';
import { useGetGraphqlApiDetails } from 'API/hooks';
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
import { ClusterObjectRef } from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';
import React from 'react';
import { useParams } from 'react-router';
import { GraphqlIconHolder } from '../GraphqlTable';
import { ExecutableGraphqlApiDetails } from './api-details/executable-api/ExecutableGraphqlApiDetails';
import { GraphqlApiExplorer } from './api-explorer/GraphqlApiExplorer';
import GraphqlDefineResolversPrompt from './GraphqlDefineResolversPrompt';
import GraphqlApiIntrospectionToggle from './GraphqlApiIntrospectionToggle';
import { GatewayGraphqlApiDetails } from './api-details/gateway-api/GatewayGraphqlApiDetails';

export const GraphqlInstance: React.FC = () => {
  // gets the graphql info from the URL
  const {
    graphqlApiName = '',
    graphqlApiNamespace = '',
    graphqlApiClusterName = '',
  } = useParams();
  const apiRef = {
    name: graphqlApiName,
    namespace: graphqlApiNamespace,
    clusterName: graphqlApiClusterName,
  } as ClusterObjectRef.AsObject;

  // gets the schema from the api
  const { data: graphqlApi } = useGetGraphqlApiDetails(apiRef);

  // tab logic
  const [tabIndex, setTabIndex] = React.useState(0);
  const handleTabsChange = (index: number) => {
    setTabIndex(index);
  };

  if (!graphqlApi) return <Loading />;
  return (
    <div className='w-full mx-auto '>
      <SectionCard
        cardName={graphqlApiName}
        logoIcon={<GraphqlIconHolder>{<GraphQLIcon />}</GraphqlIconHolder>}
        headerSecondaryInformation={[
          {
            title: 'Namespace',
            value: graphqlApiNamespace,
          },
        ]}>
        <GraphqlDefineResolversPrompt apiRef={apiRef} />

        <div className='float-right'>
          <GraphqlApiIntrospectionToggle apiRef={apiRef} />
        </div>

        <Tabs index={tabIndex} onChange={handleTabsChange}>
          <FolderTabList>
            <FolderTab>API Details</FolderTab>
            <FolderTab>Explore</FolderTab>
          </FolderTabList>

          <TabPanels>
            <StyledTabPanel>
              <FolderTabContent>
                {tabIndex === 0 && (
                  // ! Check for apiType here
                  // <GatewayGraphqlApiDetails apiRef={apiRef} />
                  <ExecutableGraphqlApiDetails apiRef={apiRef} />
                )}
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
