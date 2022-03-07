import * as React from 'react';
import { Fetcher, GraphiQL } from 'graphiql';
import { buildSchema } from 'graphql';
import styled from '@emotion/styled';
import { colors } from 'Styles/colors';
import { useGetGraphqlSchemaDetails, useListVirtualServices } from 'API/hooks';
import { useParams } from 'react-router';
import { VirtualService } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/gateway_resources_pb';
import { StatusHealth, WarningCircle } from '../Overview/OverviewBoxSummary';
import { ReactComponent as WarningExclamation } from 'assets/big-warning-exclamation.svg';
import { QuestionCircleOutlined } from '@ant-design/icons';
import { SoloInput } from 'Components/Common/SoloInput';
import { createGraphiQLFetcher } from '@graphiql/toolkit';
import { Tooltip } from 'antd';
import { copyTextToClipboard } from 'utils';
import { ReactComponent as CopyIcon } from 'assets/document.svg';

const Wrapper = styled.div`
  background: white;
`;

const StyledContainer = styled.div`
  height: 70vh;
`;

const GqlInputContainer = styled.div`
  margin: 10px auto;
`;

const GqlInputWrapper = styled.div`
  display: flex;
  flex-direction: row;
`;

const LabelTextWrapper = styled.div<{ hasError: boolean; }>`
  width: 100%;
  label {
    width: 100%;
    margin-right: 10px;
    color: ${props => props.hasError ? colors.sunGold : 'black'};
  }
  input {
    width: 250px;
  }
`;

const StyledQuestionMark = styled(QuestionCircleOutlined)`
  margin-top: 3px;
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

const GQL_STORAGE_KEY = 'gqlStorageKey';

const StyledCopyIcon = styled(CopyIcon)`
  color: white;
  display: inline-block;
  margin-left: 10px;
  fill: white;
`;

const getGqlStorage = () => {
  return (
    localStorage.getItem(GQL_STORAGE_KEY) || 'http://localhost:8080/graphql'
  );
};

const setGqlStorage = (value: string) => {
  localStorage.setItem(GQL_STORAGE_KEY, value);
};

export const GraphqlApiExplorer = () => {
  const {
    graphqlSchemaName,
    graphqlSchemaNamespace,
    graphqlSchemaClusterName,
  } = useParams();
  const [gqlError, setGqlError] = React.useState('');
  const [refetch, setRefetch] = React.useState(false);
  const [url, setUrl] = React.useState(getGqlStorage());
  const [showTooltip, setShowTooltip] = React.useState(false);
  const [copyingKubectl, setCopyingKubectl] = React.useState(false);
  const [copyingProxy, setCopyingProxy] = React.useState(false);
  const [showUrlBar, setShowUrlBar] = React.useState(false);

  const {
    data: graphqlSchema,
    error: graphqlSchemaError,
    mutate,
  } = useGetGraphqlSchemaDetails({
    name: graphqlSchemaName,
    namespace: graphqlSchemaNamespace,
    clusterName: graphqlSchemaClusterName,
  });

  const changeUrl = (value: string) => {
    setGqlStorage(value);
    setUrl(value);
  };

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
  let gqlFetcher: Fetcher = createGraphiQLFetcher({
    url,
    schemaFetcher: async graphQLParams => {
      try {
        setRefetch(false);
        setGqlError('');
        const data = await fetch(url, {
          method: 'POST',
          headers: {
            Accept: 'application/json',
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(graphQLParams),
          credentials: 'same-origin',
        });
        return data.json().catch(() => data.text());
      } catch (error: any) {
        setGqlError(error.message);
      }
    },
  });

  const graphiqlRef = React.useRef<null | GraphiQL>(null);

  const { data: virtualServices, error: virtualServicesError } =
    useListVirtualServices();
  const [correspondingVirtualServices, setCorrespondingVirtualServices] =
    React.useState<VirtualService.AsObject[]>([]);

  React.useEffect(() => {
    let correspondingVs = virtualServices?.filter(vs =>
      vs.spec?.virtualHost?.routesList.some(
        route =>
          route?.graphqlSchemaRef?.name === graphqlSchemaName &&
          route?.graphqlSchemaRef?.namespace === graphqlSchemaNamespace
      )
    );

    if (!!correspondingVs) {
      setCorrespondingVirtualServices(correspondingVs);
    }
  }, [virtualServices, graphqlSchemaName, graphqlSchemaNamespace]);

  let executableSchema;

  if (graphqlSchema?.spec?.executableSchema?.schemaDefinition) {
    const schemaDef = graphqlSchema.spec.executableSchema.schemaDefinition;
    executableSchema = buildSchema(schemaDef, {
      assumeValidSDL: true,
    });
  }

  const handlePrettifyQuery = () => {
    graphiqlRef?.current?.handlePrettifyQuery();
  };
  const changeHost = (e: React.ChangeEvent<HTMLInputElement>) => {
    setRefetch(true);
    changeUrl(e.currentTarget.value);
  };

  const toggleUrlBar = () => {
    setShowUrlBar(!showUrlBar);
  }

  // TODO:  We can hide and show elements based on what we get back.
  //        The schema will only refetch if the executable schema is undefined.
  return !!correspondingVirtualServices?.length ? (
    <Wrapper>
      {Boolean(gqlError) || showUrlBar ? (
        <GqlInputContainer>
          <GqlInputWrapper>
            <LabelTextWrapper hasError={Boolean(gqlError)}>
              <SoloInput
                title={`${Boolean(gqlError) ? 'Failed to fetch Graphql service.  Update the host to attempt again.' : 'Current url'}`}
                value={url}
                onChange={changeHost}
              />
            </LabelTextWrapper>
            <Tooltip
              title={
                <CodeWrapper>
                  <p>
                    Endpoint URL for the gateway proxy. The default URL can be
                    used if you port forward with the following command:
                  </p>
                  <p
                    className='copy'
                    title='copy command'
                    onClick={copyKubectlCommand}>
                    <code>
                      <i>
                        kubectl port-forward -n gloo-system deploy/gateway-proxy
                        8080
                      </i>
                      {copyingKubectl ? (
                        <Copied>copied!</Copied>
                      ) : (
                        <StyledCopyIcon />
                      )}
                    </code>
                  </p>
                  <p>
                    Depending on your installation, you can also use the
                    following glooctl command:
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
          </GqlInputWrapper>
        </GqlInputContainer>
      ) : null}
      <StyledContainer>
        <GraphiQL
          ref={graphiqlRef}
          defaultQuery={''}
          variables={'{}'}
          schema={!refetch ? executableSchema : undefined}
          fetcher={gqlFetcher}>
          <GraphiQL.Toolbar>
            <GraphiQL.Button
              onClick={handlePrettifyQuery}
              label='Prettify'
              title='Prettify Query (Shift-Ctrl-P)'
            />
            <GraphiQL.Button
              onClick={toggleUrlBar}
              label={showUrlBar ? 'Hide Url Bar' : 'Show Url Bar'}
              title='Show/Hide Url Bar'
            />
          </GraphiQL.Toolbar>
        </GraphiQL>
      </StyledContainer>
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
