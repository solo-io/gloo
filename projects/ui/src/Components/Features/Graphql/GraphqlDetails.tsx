import styled from '@emotion/styled';
import { TabPanels, Tabs } from '@reach/tabs';
import { graphqlApi } from 'API/graphql';
import { useGetGraphqlSchemaDetails, useListUpstreams } from 'API/hooks';
import { ReactComponent as GraphQLIcon } from 'assets/graphql-icon.svg';
import AreaHeader from 'Components/Common/AreaHeader';
import { Loading } from 'Components/Common/Loading';
import { SectionCard } from 'Components/Common/SectionCard';
import { SoloModal } from 'Components/Common/SoloModal';
import {
  FolderTab,
  FolderTabContent,
  FolderTabList,
  StyledTabPanel,
} from 'Components/Common/Tabs';
import { Upstream } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/gloo_resources_pb';
import React from 'react';
import { useNavigate, useParams } from 'react-router';
import { colors } from 'Styles/colors';
import { GraphqlApiExplorer } from './GraphqlApiExplorer';
import { GraphqlIconHolder } from './GraphqlTable';
import ResolversTable from './ResolversTable';
import { ResolverWizard } from './ResolverWizard';

export const OperationDescription = styled('div')`
  grid-column: span 3 / span 3;
  /* Hide scrollbar for Chrome, Safari and Opera */
  &::-webkit-scrollbar {
    display: none !important;
  }

  /* Hide scrollbar for IE, Edge and Firefox */
  & {
    -ms-overflow-style: none !important; /* IE and Edge */
    scrollbar-width: none !important; /* Firefox */
  }
`;

type ArrowToggleProps = { active?: boolean };
export const ArrowToggle = styled('div')<ArrowToggleProps>`
  position: absolute;
  left: 1rem;

  &:before,
  &:after {
    position: absolute;
    content: '';
    display: block;
    width: 8px;
    height: 1px;
    background: ${colors.septemberGrey};
    transition: transform 0.3s;
  }

  &:before {
    right: 5px;
    border-top-left-radius: 10px;
    border-bottom-left-radius: 10px;
    transform: rotate(${props => (props.active ? '' : '-')}45deg);
  }

  &:after {
    right: 1px;
    transform: rotate(${props => (props.active ? '-' : '')}45deg);
  }
`;

const ConfigArea = styled.div`
  margin-bottom: 20px;
`;

export type GraphQLDetailsProps = {};

export const GraphQLDetails: React.FC<GraphQLDetailsProps> = props => {
  const {
    graphqlSchemaName = '',
    graphqlSchemaNamespace = '',
    graphqlSchemaClusterName = '',
  } = useParams();

  const navigate = useNavigate();
  const { data: graphqlSchema, error: graphqlSchemaError } =
    useGetGraphqlSchemaDetails({
      name: graphqlSchemaName,
      namespace: graphqlSchemaNamespace,
      clusterName: graphqlSchemaClusterName,
    });
  const [tabIndex, setTabIndex] = React.useState(0);
  const [currentResolver, setCurrentResolver] = React.useState<any>();
  const [modalOpen, setModalOpen] = React.useState(false);

  const { data: upstreams, error: upstreamsError } = useListUpstreams();
  const [resolverUpstreams, setResolverUpstreams] = React.useState<
    Upstream.AsObject[]
  >([]);

  React.useEffect(() => {
    let resolverUpstreams =
      graphqlSchema?.spec?.executableSchema?.executor?.local?.resolutionsMap
        .filter(
          ([rName, r], index, arr) =>
            index ===
            arr?.findIndex(
              ([n, rr]) =>
                rr?.restResolver?.upstreamRef?.name ===
                r.restResolver?.upstreamRef?.name
            )
        )
        .map(([resolveName, resolver]) => resolver.restResolver?.upstreamRef);

    let fullUpstreams = upstreams?.filter(
      upstream =>
        !!resolverUpstreams?.find(
          rU =>
            rU?.name === upstream.metadata?.name &&
            rU?.namespace === upstream.metadata?.namespace
        )
    );
    if (!!fullUpstreams) {
      setResolverUpstreams(fullUpstreams);
    }
  }, []);
  const handleTabsChange = (index: number) => {
    setTabIndex(index);
  };
  const closeModal = () => setModalOpen(false);
  const loadYaml = async () => {
    if (!graphqlSchemaName || !graphqlSchemaNamespace) {
      return '';
    }

    try {
      const yaml = await graphqlApi.getGraphqlSchemaYaml({
        name: graphqlSchemaName,
        namespace: graphqlSchemaNamespace,
        clusterName: graphqlSchemaClusterName,
      });
      return yaml;
    } catch (error) {
      console.error(error);
    }
    return '';
  };

  if (!graphqlSchema) return <Loading />;

  return (
    <React.Fragment>
      <div className='w-full mx-auto '>
        <SectionCard
          cardName={graphqlSchemaName}
          logoIcon={<GraphqlIconHolder>{<GraphQLIcon />}</GraphqlIconHolder>}
          headerSecondaryInformation={[
            {
              title: 'Namespace',
              value: graphqlSchemaNamespace,
            },
            {
              title: 'Introspection',
              value: graphqlSchema.spec?.executableSchema?.executor?.local
                ?.enableIntrospection
                ? 'Enabled'
                : 'Disabled',
            },
          ]}
        >
          <Tabs index={tabIndex} onChange={handleTabsChange}>
            <FolderTabList>
              <FolderTab>API Details</FolderTab>
              <FolderTab>Explore</FolderTab>
            </FolderTabList>

            <TabPanels>
              <StyledTabPanel>
                <FolderTabContent>
                  <>
                    <ConfigArea>
                      <AreaHeader
                        title='Configuration'
                        contentTitle={`${graphqlSchemaNamespace}--${graphqlSchemaName}.yaml`}
                        onLoadContent={loadYaml}
                      />

                      <ResolversTable
                        schemaRef={{
                          name: graphqlSchemaName,
                          namespace: graphqlSchemaNamespace,
                          clusterName: graphqlSchemaClusterName,
                        }}
                      />
                    </ConfigArea>
                    <ConfigArea>
                      <div className='flex p-4 mb-5 bg-gray-100 border border-gray-300 rounded-lg'>
                        <div className='w-1/5 mr-5'>
                          <div className='mb-2 text-lg font-medium'>
                            Upstreams
                          </div>
                          {resolverUpstreams?.map(resolverUpstream => {
                            const glooInstNamespace =
                              resolverUpstream.glooInstance?.namespace;
                            const glooInstName =
                              resolverUpstream.glooInstance?.name;
                            const upstreamCluster =
                              resolverUpstream.metadata?.clusterName ?? '';
                            const upstreamNamespace =
                              resolverUpstream.metadata?.namespace ?? '';
                            const upstreamName =
                              resolverUpstream.metadata?.name ?? '';
                            const link = !!upstreamCluster
                              ? `/gloo-instances/${glooInstNamespace}/${glooInstName}/upstreams/${upstreamCluster}/${upstreamNamespace}/${upstreamName}`
                              : `/gloo-instances/${glooInstNamespace}/${glooInstName}/upstreams/${upstreamNamespace}/${upstreamName}`;
                            return (
                              <div key={link}>
                                <div
                                  className={
                                    'cursor-pointer text-blue-500gloo text-base'
                                  }
                                  onClick={() => {
                                    navigate(link);
                                  }}
                                >
                                  {upstreamName}
                                </div>
                              </div>
                            );
                          })}
                        </div>
                      </div>
                    </ConfigArea>
                  </>
                </FolderTabContent>
              </StyledTabPanel>
              <StyledTabPanel>
                <div>
                  {tabIndex === 1 && (
                    <GraphqlApiExplorer graphQLSchema={graphqlSchema} />
                  )}
                </div>
              </StyledTabPanel>
            </TabPanels>
          </Tabs>
        </SectionCard>
      </div>
      <SoloModal visible={modalOpen} width={750} onClose={closeModal}>
        <ResolverWizard resolver={currentResolver} onClose={closeModal} />
      </SoloModal>
    </React.Fragment>
  );
};
