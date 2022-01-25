import styled from '@emotion/styled';
import { ReactComponent as CodeIcon } from 'assets/code-icon.svg';
import { ReactComponent as GraphQLIcon } from 'assets/graphql-icon.svg';
import { ReactComponent as RouteIcon } from 'assets/route-icon.svg';
import AreaHeader from 'Components/Common/AreaHeader';
import { SectionCard } from 'Components/Common/SectionCard';
import { SoloModal } from 'Components/Common/SoloModal';
import React from 'react';
import { useNavigate, useParams } from 'react-router';
import { colors } from 'Styles/colors';
import YAML from 'yaml';
import { bookInfoYsml } from './data/book-info-yaml';
import graphQLSchema from './data/book-info.json';
import { GraphqlApiExplorer } from './GraphqlApiExplorer';
import { GraphqlIconHolder } from './GraphqlTable';
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

type ResolverType = {
  restResolver: {
    request: {
      headers: {
        ':method': string;
        ':path': string;
      };
    };
    upstreamRef: {
      name: string;
      namespace: string;
    };
  };
};
export interface ResolverMapType {
  [resolverName: string]: {
    resolver: ResolverType;
  };
}

const ConfigArea = styled.div`
  margin-bottom: 20px;
`;

const YamlViewingSection = styled.div`
  top: -35px;
`;

export type GraphQLDetailsProps = {};

export const ResolversTable: React.FC<{
  resolvers: typeof graphQLSchema.spec.executableSchema.executor.local.resolutions;
  handleResolverConfigModal: <T>(resolverName: string, resolver: T) => void;
  isQueryType?: boolean;
}> = props => {
  const { name, namespace } = useParams();
  const [isOpen, setIsOpen] = React.useState(true);
  const listRef = React.useRef<HTMLDivElement>(null);
  const navigate = useNavigate();

  const { resolvers, handleResolverConfigModal, isQueryType = false } = props;
  let [qResolverName = '', qResolver = ''] = Object.entries(resolvers).find(
    ([qResolverN, qResolverR]) => qResolverN.includes('Query|')
  )!;

  const queryType =
    Object.keys(resolvers).find(r => r.includes('Query|')) ?? '';

  let isQueryResolverConfigured = true;
  return (
    <div>
      <div className='relative flex flex-col w-full bg-gray-200 border h-28'>
        <div className='flex items-center justify-between gap-5 pt-4 my-2 ml-4 h-14 '>
          <div className='flex items-center mr-3'>
            <GraphqlIconHolder>
              <GraphQLIcon className='w-4 h-4 fill-current' />
            </GraphqlIconHolder>
            <span className='flex items-center font-medium text-gray-900 whitespace-nowrap'>
              {isQueryType ? 'Query' : 'Product'}
            </span>
          </div>
        </div>
        <div className='flex items-center justify-between w-full px-6 py-4 text-sm font-medium text-gray-900 whitespace-nowrap'>
          <div
            className='relative flex-wrap justify-between w-full h-full text-sm '
            style={{
              display: 'grid',
              flexWrap: 'wrap',
              gridTemplateColumns: '1fr 1fr  minmax(120px, 200px) 105px',
              gridTemplateRows: '1fr',
              gridAutoRows: 'min-content',
              columnGap: '15px',
            }}
          >
            <span className='flex items-center justify-start ml-6 font-medium text-gray-900 '>
              Field
            </span>
            <span className='flex items-center justify-start ml-8 font-medium text-gray-900 '>
              {isQueryType ? 'Type' : 'Path'}
            </span>

            <span className='flex items-center justify-start ml-8 font-medium text-gray-900 '>
              Resolver
            </span>
          </div>
        </div>
        <div
          className='absolute top-0 right-0 flex items-center w-10 h-10 p-4 mr-2 cursor-pointer '
          onClick={() => setIsOpen(!isOpen)}
        >
          <ArrowToggle active={isOpen} className='self-center m-4 ' />
        </div>
      </div>

      {isOpen && (
        <div>
          {isQueryType ? (
            <div
              key={`${namespace}-${name}-${queryType}`}
              className={`flex h-20 p-4 pl-0 border `}
            >
              <div className='flex items-center px-4 text-sm font-medium text-gray-900 whitespace-nowrap'>
                <CodeIcon className='w-4 h-4 ml-2 mr-3 fill-current text-blue-600gloo' />
              </div>
              <div className='relative flex items-center w-full text-sm text-gray-500 whitespace-nowrap'>
                <div
                  className='relative flex-wrap justify-between w-full h-full text-sm '
                  style={{
                    display: 'grid',
                    flexWrap: 'wrap',
                    gridTemplateColumns: '1fr 1fr  minmax(120px, 200px) 105px',
                    gridTemplateRows: '1fr',
                    gridAutoRows: 'min-content',
                    columnGap: '5px',
                  }}
                >
                  <span className='flex items-center font-medium text-gray-900 '>
                    {qResolverName.replace('Query|', '')}
                  </span>
                  <span className='flex items-center text-sm text-gray-700 '>
                    {`[Product]`}
                  </span>
                  <span className={`flex items-center justify-center`}>
                    {isQueryResolverConfigured ? (
                      <span
                        className={`inline-flex items-center min-w-max p-1 px-2 ${
                          !isQueryResolverConfigured
                            ? 'focus:ring-blue-500gloo text-blue-700gloo bg-blue-200gloo  border-blue-600gloo hover:bg-blue-300gloo'
                            : 'focus:ring-gray-500 text-gray-700 bg-gray-300  border-gray-600 hover:bg-gray-200'
                        }   border rounded-full shadow-sm cursor-pointer  focus:outline-none focus:ring-2 focus:ring-offset-2 `}
                        onClick={() => {
                          if (handleResolverConfigModal) {
                            let [qResolverName, qResolver] = Object.entries(
                              resolvers
                            ).find(([qResolverN, qResolverR]) =>
                              qResolverN.includes('Query|')
                            )!;
                            handleResolverConfigModal<typeof qResolver>(
                              qResolverName,
                              qResolver
                            );
                          }
                        }}
                      >
                        <RouteIcon className='w-6 h-6 mr-1 fill-current text-blue-600gloo' />

                        {isQueryResolverConfigured ? 'Configure' : 'Configure'}
                      </span>
                    ) : (
                      <div></div>
                    )}
                  </span>
                </div>
              </div>
            </div>
          ) : (
            Object.entries(resolvers)
              .filter(([resolverName, resolver]) => resolverName !== queryType)
              .map(([resolverName, resolver]) => {
                let isConfigured = false;

                return (
                  <div
                    key={`${namespace}-${name}-${resolverName}`}
                    className={`flex h-20 p-4 pl-0 border `}
                  >
                    <div className='flex items-center px-4 text-sm font-medium text-gray-900 whitespace-nowrap'>
                      <CodeIcon className='w-4 h-4 ml-2 mr-3 fill-current text-blue-600gloo' />
                    </div>
                    <div className='relative flex items-center w-full text-sm text-gray-500 whitespace-nowrap'>
                      <div
                        className='relative flex-wrap justify-between w-full h-full text-sm '
                        style={{
                          display: 'grid',
                          flexWrap: 'wrap',
                          gridTemplateColumns:
                            '1fr 1fr  minmax(120px, 200px) 105px',
                          gridTemplateRows: '1fr',
                          gridAutoRows: 'min-content',
                          columnGap: '5px',
                        }}
                      >
                        <span className='flex items-center font-medium text-gray-900 '>
                          {resolverName}
                        </span>
                        <span className='flex items-center text-sm text-gray-700 '>
                          {resolver.restResolver.request.headers[':path']}
                        </span>
                        <span className={`flex items-center justify-center`}>
                          {!isConfigured ? (
                            <span
                              className={`inline-flex items-center min-w-max p-1 px-2 ${
                                isConfigured
                                  ? 'focus:ring-blue-500gloo text-blue-700gloo bg-blue-200gloo  border-blue-600gloo hover:bg-blue-300gloo'
                                  : 'focus:ring-gray-500 text-gray-700 bg-gray-300  border-gray-600 hover:bg-gray-200'
                              }   border rounded-full shadow-sm cursor-pointer  focus:outline-none focus:ring-2 focus:ring-offset-2 `}
                              onClick={() => {
                                if (handleResolverConfigModal) {
                                  handleResolverConfigModal<typeof resolver>(
                                    resolverName,
                                    resolver
                                  );
                                }
                              }}
                            >
                              <RouteIcon className='w-6 h-6 mr-1 fill-current text-blue-600gloo' />

                              {isConfigured ? '' : 'Configure'}
                            </span>
                          ) : (
                            <div></div>
                          )}
                        </span>
                      </div>
                    </div>
                  </div>
                );
              })
          )}
        </div>
      )}
    </div>
  );
};
export const GraphQLDetails: React.FC<GraphQLDetailsProps> = props => {
  const { name, namespace } = useParams();
  const navigate = useNavigate();
  const [showSchemaExplorer, setShowSchemaExplorer] = React.useState(false);
  const [currentResolver, setCurrentResolver] = React.useState<any>(
    graphQLSchema.spec.executableSchema.executor.local.resolutions[
      'Query|productsForHome'
    ]
  );

  const [modalOpen, setModalOpen] = React.useState(false);
  const openModal = () => setModalOpen(true);
  const closeModal = () => setModalOpen(false);
  const loadYaml = async () => {
    if (!name || !namespace) {
      return '';
    }

    try {
      const yaml = YAML.stringify(bookInfoYsml);
      return yaml;
    } catch (error) {
      console.error(error);
    }
    return '';
  };

  function handleResolverConfigModal<T>(resolverName: string, resolver: T) {
    setCurrentResolver(resolver);
    openModal();
  }
  let hasResolver = true;
  let showYAML = true;
  return (
    <React.Fragment>
      <div className='relative w-full mx-auto '>
        <YamlViewingSection className='absolute right-0 flex items-center cursor-pointer'>
          <div
            className={
              'cursor-pointer mr-4 flex items-center text-blue-500gloo'
            }
            onClick={() => setShowSchemaExplorer(s => !s)}
          >
            Schema Explorer
          </div>
        </YamlViewingSection>
        <SectionCard
          cardName={name!}
          logoIcon={<GraphqlIconHolder>{<GraphQLIcon />}</GraphqlIconHolder>}
          headerSecondaryInformation={[
            {
              title: 'Namespace',
              value: namespace,
            },
            {
              title: 'Introspection',
              value: graphQLSchema.spec.executableSchema.executor.local
                .enableIntrospection
                ? 'Enabled'
                : 'Disabled',
            },
          ]}
        >
          {showSchemaExplorer ? (
              <GraphqlApiExplorer graphQLSchema={graphQLSchema} />
          ) : (
            <>
              {' '}
              <ConfigArea>
                <AreaHeader
                  title='Configuration'
                  contentTitle={`${namespace}--${name}.yaml`}
                  onLoadContent={loadYaml}
                />

                <div className='relative overflow-x-hidden overflow-y-scroll '>
                  <ResolversTable
                    isQueryType
                    resolvers={
                      graphQLSchema.spec.executableSchema.executor.local
                        .resolutions
                    }
                    handleResolverConfigModal={handleResolverConfigModal}
                  />
                </div>

                <div className='relative mt-4 overflow-x-hidden overflow-y-scroll'>
                  <ResolversTable
                    resolvers={
                      graphQLSchema.spec.executableSchema.executor.local
                        .resolutions
                    }
                    handleResolverConfigModal={handleResolverConfigModal}
                  />
                </div>
              </ConfigArea>
              <ConfigArea>
                {/* <AreaTitle>Upstreams</AreaTitle> */}
                <div className='flex p-4 mb-5 bg-gray-100 border border-gray-300 rounded-lg'>
                  <div className='w-1/5 mr-5'>
                    <div className='mb-2 text-lg font-medium'>Upstreams</div>
                    {Object.entries(
                      graphQLSchema.spec.executableSchema.executor.local
                        .resolutions
                    )
                      // remove duplicate upstreams
                      .filter(
                        ([rName, r], index, arr) =>
                          index ===
                          arr?.findIndex(
                            ([n, rr]) =>
                              rr?.restResolver?.upstreamRef?.name ===
                              r.restResolver.upstreamRef.name
                          )
                      )
                      ?.map(([resolverName, resolver]) => {
                        return (
                          <div
                            key={`/${resolverName}/${resolver.restResolver.upstreamRef.namespace}/${resolver.restResolver.upstreamRef.name}`}
                          >
                            <div
                              className={
                                'cursor-pointer text-blue-500gloo text-base'
                              }
                              onClick={() => {
                                navigate(
                                  `/upstreams/${resolver.restResolver.upstreamRef.namespace}/${resolver.restResolver.upstreamRef.name}`
                                );
                              }}
                            >
                              {resolver.restResolver.upstreamRef.name}
                            </div>
                          </div>
                        );
                      })}
                  </div>
                </div>
              </ConfigArea>
            </>
          )}
        </SectionCard>
      </div>
      <SoloModal visible={modalOpen} width={750} onClose={closeModal}>
        <ResolverWizard resolver={currentResolver} onClose={closeModal} />
      </SoloModal>
    </React.Fragment>
  );
};
