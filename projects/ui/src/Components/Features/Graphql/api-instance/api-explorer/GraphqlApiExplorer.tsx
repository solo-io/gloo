import { QuestionCircleOutlined } from '@ant-design/icons';
import { Global } from '@emotion/core';
import styled from '@emotion/styled';
import { createGraphiQLFetcher } from '@graphiql/toolkit';
import { Tooltip } from 'antd';
import {
  useGetGraphqlApiDetails,
  useListRouteTables,
  useListVirtualServices,
} from 'API/hooks';
import { ReactComponent as WarningExclamation } from 'assets/big-warning-exclamation.svg';
import { ReactComponent as CopyIcon } from 'assets/document.svg';
import { Loading } from 'Components/Common/Loading';
import { SoloInput } from 'Components/Common/SoloInput';
import { GraphiQL, Storage } from 'graphiql';
// @ts-ignore
import GraphiQLExplorer from 'graphiql-explorer';
import { DocumentNode } from 'graphql';
import { ResourceRef } from 'proto/github.com/solo-io/solo-kit/api/v1/ref_pb';
import { VirtualService } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/gateway_resources_pb';
import * as React from 'react';
import { useEffect, useMemo, useRef, useState } from 'react';
import { useParams } from 'react-router';
import { colors } from 'Styles/colors';
import { copyTextToClipboard } from 'utils';
import { useGetSchema } from 'utils/graphql-helpers';
import {
  StatusHealth,
  WarningCircle,
} from '../../../Overview/OverviewBoxSummary';
import graphiqlCustomStyles from './GraphqlApiExplorer.style';

type TabState = {
  id: string;
  hash: string;
  title: string;
  query: string | undefined;
  variables: string | undefined;
  headers: string | undefined;
  operationName: string | undefined;
  response: string | undefined;
};

type TabsState = {
  activeTabIndex: number;
  tabs: Array<TabState>;
};

const Wrapper = styled.div`
  background: white;
`;

const GqlInputContainer = styled.div`
  padding: 15px 10px;
  border-bottom: 1px solid ${colors.marchGrey};
  display: flex;
`;

const GqlInputWrapper = styled.div`
  flex-basis: min-content;
  display: flex;
  flex-direction: row;
`;

const LabelTextWrapper = styled.div<{ hasError: boolean }>`
  width: 100%;
  label {
    width: 100%;
    margin-right: 10px;
    color: ${props => (props.hasError ? colors.sunGold : 'black')};
  }
  input {
    width: 500px;
  }
`;

const StyledQuestionMark = styled(QuestionCircleOutlined)`
  margin-top: 3px;
  margin-left: 20px;
  display: inline-flex;
  &:hover {
    cursor: pointer;
  }
`;

const CodeWrapper = styled.div`
  code {
    &:hover {
      cursor: pointer;
      color: ${colors.aprilGrey};
      fill: ${colors.aprilGrey};
    }
  }
  p {
    padding: 10px 0;
  }
`;

const Copied = styled.span`
  display: inline-block;
  margin-left: 10px;
`;

const StyledCopyIcon = styled(CopyIcon)`
  color: white;
  display: inline-block;
  margin-left: 10px;
  fill: white;
`;

const defaultQuery = `query Example {
}

# Welcome to GraphiQL, an in-browser tool for
# writing, validating, and testing GraphQL queries.
#
# Type queries into this side of the screen, and you
# will see intelligent typeaheads aware of the current
# GraphQL type schema and live syntax and
# validation errors highlighted within the text.
#
# GraphQL queries typically start with a "{" character.
# Lines that start with a # are ignored.
# The name of the query on the first line of each tab
# is the title of that tab.
#
# An example GraphQL query might look like:
#     query Example {
#       field(arg: "value") {
#         subField
#       }
#     }
#
# Keyboard shortcuts:
#     Prettify Query:    Shift-Ctrl-P
#     Merge Query:     Shift-Ctrl-M
#     Run Query:        Ctrl-Enter
#     Auto Complete:  Ctrl-Space

`;

const refsMatch = (
  ref1: ResourceRef.AsObject | undefined,
  ref2: ResourceRef.AsObject | undefined
) => {
  if (!ref1 && !ref2 && ref1 === ref2) return true;
  if (!ref1 || !ref2) return false;
  return ref1.name === ref2.name && ref1.namespace === ref2.namespace;
};

export const GraphqlApiExplorer = () => {
  const { graphqlApiName, graphqlApiNamespace, graphqlApiClusterName } =
    useParams();
  const [gqlError, setGqlError] = useState('');
  const [explorerOpen, setExplorerOpen] = useState(false);
  const [showTooltip, setShowTooltip] = useState(false);
  const [copyingKubectl, setCopyingKubectl] = useState(false);
  const [copyingProxy, setCopyingProxy] = useState(false);
  const [showUrlBar, setShowUrlBar] = useState(false);
  const [query, setQuery] = useState<string>();
  const graphiqlRef = useRef<null | GraphiQL>(null);
  const [hasTriedToFetch, setHasTriedToFetch] = useState(false);

  //
  // Including GraphQL API info in the localStorage keys, so they don't override each other.
  //
  const localStorageKey = `${graphqlApiName}/${graphqlApiNamespace}/${graphqlApiClusterName}:`;
  const localStorageUrlKey = localStorageKey + 'url';
  const getGqlStorage = () => {
    return (
      localStorage.getItem(localStorageUrlKey) ||
      'http://localhost:8080/graphql'
    );
  };
  const setGqlStorage = (value: string) => {
    localStorage.setItem(localStorageUrlKey, value);
  };
  const customStorage: Storage = {
    getItem: key => localStorage.getItem(localStorageKey + key),
    removeItem: key => localStorage.removeItem(localStorageKey + key),
    setItem: (key, value) => localStorage.setItem(localStorageKey + key, value),
    length: localStorage.length,
  };

  //
  // Schema URL
  //
  const urlDebounceMs = 1000;
  // `url` is updated after no input is registered for `urlDebounceMs`.
  const [url, setUrl] = useState(getGqlStorage());
  // The urlToDisplay is what the user edits, and updates instantly.
  const [urlToDisplay, setUrlToDisplay] = useState(url);
  useEffect(() => {
    // `urlTimeout` keeps track of the input delay.
    // Sets the real `url` after `urlDebounceMs`
    const urlTimeout = setTimeout(() => setUrl(urlToDisplay), urlDebounceMs);
    return () => {
      // Clears the timeout if the useEffect dependencies change.
      // (in which case a new timeout will be set)
      clearTimeout(urlTimeout);
    };
  }, [setUrl, urlToDisplay]);
  //
  // When the url changes, this updates localStorage.
  useEffect(() => {
    setGqlStorage(url);
  }, [url]);

  const graphqlApiRef: ResourceRef.AsObject = {
    name: graphqlApiName,
    namespace: graphqlApiNamespace,
  };
  const { data: graphqlApi } = useGetGraphqlApiDetails({
    ...graphqlApiRef,
    clusterName: graphqlApiClusterName,
  });
  const { parsedSchema } = useGetSchema(graphqlApi);

  const copyKubectlCommand = async () => {
    setCopyingKubectl(true);
    const text =
      'kubectl port-forward -n gloo-system deploy/gateway-proxy 8080';
    copyTextToClipboard(text).finally(() => {
      setTimeout(() => {
        setCopyingKubectl(false);
      }, 2000);
    });
  };

  const copyGlooctlCommand = async () => {
    setCopyingProxy(true);
    const text = 'glooctl proxy url';
    copyTextToClipboard(text).finally(() => {
      setTimeout(() => {
        setCopyingProxy(false);
      }, 2000);
    });
  };

  // If we need the custom fetcher, we can add `schemaFetcher` to the document.
  const gqlFetcher = useMemo(
    () =>
      createGraphiQLFetcher({
        url,
        schemaFetcher: async graphQLParams => {
          if (!graphQLParams.variables?.trim()) graphQLParams.variables = '{}';
          try {
            const data = await fetch(url, {
              method: 'POST',
              headers: {
                Accept: 'application/json',
                'Content-Type': 'application/json',
              },
              body: JSON.stringify(graphQLParams),
              credentials: 'same-origin',
            });
            setHasTriedToFetch(true);
            return data.json().catch(() => data.text());
          } catch (error: any) {}
        },
      }),
    [url, setHasTriedToFetch]
  );
  //
  // This checks if the URL is valid, separate from the fetcher
  // so we can control when it's called.
  useEffect(() => {
    (async () => {
      try {
        const res = await fetch(url, {
          method: 'POST',
          headers: {
            Accept: 'application/json',
            'Content-Type': 'application/json',
          },
          credentials: 'same-origin',
        });
        if (!res.ok) throw new Error(res.statusText);
        setGqlError('');
      } catch (error: any) {
        setGqlError(error.message);
      }
    })();
  }, [url, setGqlError]);

  //
  // Gets the corresponding VirtualServices, which can be either:
  // - a VirtualService with a matching graphqlApiRef, or
  // - a VirtualService with a delegateAction to a RouteTable,
  //   which has a matching graphqlApiRef.
  //
  const { data: vsResponse } = useListVirtualServices();
  const virtualServices = vsResponse?.virtualServicesList;
  const { data: rtResponse } = useListRouteTables();
  const routeTables = rtResponse?.routeTablesList;
  const correspondingVirtualServices = useMemo(() => {
    //
    // Gets the VirtualServices with a matching graphqlApiRef.
    const vssThatMatch = [] as VirtualService.AsObject[];
    const vssThatDontMatch = [] as VirtualService.AsObject[];
    if (virtualServices) {
      for (let vs of virtualServices) {
        if (
          vs.spec?.virtualHost?.routesList.some(r =>
            refsMatch(r.graphqlApiRef, graphqlApiRef)
          )
        )
          vssThatMatch.push(vs);
        else vssThatDontMatch.push(vs);
      }
    }
    //
    // Gets the VirtualServices with a delegateAction to a RouteTable,
    // which has a matching graphqlApiRef
    const rtsThatMatch = routeTables?.filter(rt =>
      rt.spec?.routesList.some(r => refsMatch(r.graphqlApiRef, graphqlApiRef))
    );
    // Only need to check VirtualServices that aren't already a correspondingVirtualService.
    const vssOfRtsThatMatch = vssThatDontMatch?.filter(vs =>
      vs.spec?.virtualHost?.routesList.some(r =>
        rtsThatMatch?.some(rt =>
          refsMatch(
            !!rt.metadata
              ? { name: rt.metadata.name, namespace: rt.metadata.namespace }
              : undefined,
            r.delegateAction?.ref
          )
        )
      )
    );
    //
    // Merge the VirtualServices (we made sure there aren't duplicates)
    return [...vssThatMatch, vssOfRtsThatMatch];
  }, [virtualServices, routeTables, graphqlApiName, graphqlApiNamespace]);

  const handlePrettifyQuery = () => {
    graphiqlRef?.current?.handlePrettifyQuery();
  };

  const toggleExplorer = () => {
    setExplorerOpen(!explorerOpen);
  };

  const handleQueryUpdate = (query?: string, documentAST?: DocumentNode) => {
    setQuery(query);
    // This queryChanged check is needed, since the explorer tab
    // might need to manually update the graphiqlRef in this way.
    const queryChanged =
      graphiqlRef.current?.getQueryEditor().getValue() !== query;
    if (queryChanged) graphiqlRef.current?.handleEditQuery(query ?? '');
  };

  // The operation name === the tab name.
  // This comes from the actual GraphQL operation.
  //   e.g. query Test {} will have the tab name "Test".
  //   The example has "query Example { }" in it.
  const [opName, setOpName] = useState('Example');

  // TODO:  We can hide and show elements based on what we get back.
  //        The schema will only refetch if the executable schema is undefined.
  if (correspondingVirtualServices === undefined) return null;
  if (parsedSchema === undefined) return <Loading />;
  return correspondingVirtualServices.length > 0 ? (
    <Wrapper>
      <Global styles={graphiqlCustomStyles} />
      {Boolean(gqlError) || showUrlBar ? (
        <GqlInputContainer>
          <GqlInputWrapper>
            <LabelTextWrapper hasError={Boolean(gqlError)}>
              <SoloInput
                label={
                  <>
                    <div className='ml-2'>
                      <span className='text-sm'>
                        {gqlError
                          ? 'Failed to fetch GraphQL service.  Update the host to attempt again.'
                          : 'Endpoint URL'}
                      </span>
                      <Tooltip
                        title={
                          <CodeWrapper>
                            <p>
                              Endpoint URL for the gateway proxy. The default
                              URL can be used if you port forward with the
                              following command:
                            </p>
                            <p
                              className='copy'
                              title='copy command'
                              onClick={copyKubectlCommand}>
                              <code>
                                <i>
                                  kubectl port-forward -n gloo-system
                                  deploy/gateway-proxy 8080
                                </i>
                                {copyingKubectl ? (
                                  <Copied>copied!</Copied>
                                ) : (
                                  <StyledCopyIcon />
                                )}
                              </code>
                            </p>
                            <p>
                              Depending on your installation, you can also use
                              the following glooctl command:
                            </p>
                            <p
                              className='copy'
                              title='copy command'
                              onClick={copyGlooctlCommand}>
                              <code>
                                <i>glooctl proxy url</i>
                                {copyingProxy ? (
                                  <Copied>copied!</Copied>
                                ) : (
                                  <StyledCopyIcon />
                                )}
                              </code>
                            </p>
                          </CodeWrapper>
                        }
                        trigger='hover'
                        visible={showTooltip}
                        onVisibleChange={() => {
                          setShowTooltip(!showTooltip);
                        }}>
                        <StyledQuestionMark />
                      </Tooltip>
                    </div>
                  </>
                }
                value={urlToDisplay}
                onChange={e => setUrlToDisplay(e.currentTarget.value)}
              />
            </LabelTextWrapper>
          </GqlInputWrapper>
        </GqlInputContainer>
      ) : null}
      <div className='graphiql-outer-container'>
        <GraphiQLExplorer
          schema={hasTriedToFetch ? parsedSchema : undefined}
          query={query}
          onEdit={handleQueryUpdate}
          onRunOperation={(operationName?: string) =>
            graphiqlRef.current?.handleRunQuery(operationName)
          }
          explorerIsOpen={explorerOpen}
          onToggleExplorer={toggleExplorer}
        />
        <GraphiQL
          ref={graphiqlRef}
          defaultQuery={defaultQuery}
          variables={'{}'}
          tabs={{
            onTabChange: (tabs: TabsState) => {
              /**
               * This is a little gnarly, but they don't have a way to update the tabsState
               * from within this method.  Here's the original PR on graphiql:
               *
               * https://github.com/graphql/graphiql/pull/2197/files#diff-26ce5690905d4057a50dc0071ebe62289aa386651901373ea48ca6a499f5639a
               *
               * Using the `graphiqlRef.current?.safeSetState` doesn't work because it uses a
               * reducer to calculate the state.
               *
               * So we have to fake manually entering a variable whenever the tab is changed.
               */
              const currentTab = tabs.tabs[tabs.activeTabIndex];
              const performChange = !Boolean(currentTab.variables?.trim());
              handleQueryUpdate(currentTab.query);
              if (performChange) {
                graphiqlRef.current?.handleEditVariables('{}');
              }
            },
          }}
          onEditQuery={handleQueryUpdate}
          storage={customStorage}
          query={query}
          operationName={opName}
          onEditOperationName={s => setOpName(s)}
          schema={hasTriedToFetch ? parsedSchema : undefined}
          fetcher={gqlFetcher}>
          <GraphiQL.Toolbar>
            <GraphiQL.Button
              onClick={toggleExplorer}
              label={explorerOpen ? 'Hide Explorer' : 'Show Explorer'}
              title='Show/Hide Explorer'
            />
            <GraphiQL.Button
              onClick={handlePrettifyQuery}
              label='Prettify'
              title='Prettify Query (Shift-Ctrl-P)'
            />
            <GraphiQL.Button
              onClick={() => setShowUrlBar(!showUrlBar)}
              label={showUrlBar ? 'Hide Url Bar' : 'Show Url Bar'}
              title='Show/Hide Url Bar'
            />
          </GraphiQL.Toolbar>
        </GraphiQL>
      </div>
    </Wrapper>
  ) : (
    <div className='overflow-hidden bg-white rounded-lg shadow'>
      <div className='px-4 py-5 sm:p-6'>
        <StatusHealth isWarning>
          <div>
            <WarningCircle>
              <WarningExclamation />
            </WarningCircle>
          </div>
          <div>
            <>
              <div className='text-xl '>Explorer unavailable</div>
              <div className='text-lg '>
                There is no Virtual Service that exposes this GraphQL endpoint
              </div>
            </>
          </div>
        </StatusHealth>
      </div>
    </div>
  );
};
