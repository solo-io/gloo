import React from 'react';
import styled from '@emotion/styled/macro';
import { RouteProps, RouteComponentProps } from 'react-router-dom';
import { ReactComponent as EnvoyLogo } from 'assets/envoy-logo.svg';
import { ReactComponent as HealthScoreIcon } from 'assets/health-score-icon.svg';
import { ReactComponent as VSIcon } from 'assets/virtualservice-icon.svg';
import { ReactComponent as USIcon } from 'assets/upstream-icon.svg';
import { ReactComponent as EnvoyIcon } from 'assets/envoy-logo.svg';
import { colors } from 'Styles/colors';

import { TallyInformationDisplay } from 'Components/Common/DisplayOnly/TallyInformationDisplay';
import { GoodStateCongratulations } from 'Components/Common/DisplayOnly/GoodStateCongratulations';
import { StatusTile } from 'Components/Common/DisplayOnly/StatusTile';
import { soloConstants, healthConstants } from 'Styles';
import { CardCSS } from 'Styles/CommonEmotions/card';
import { HealthIndicator } from 'Components/Common/HealthIndicator';
import { useGetUpstreamsListV2 } from 'Api/v2/useUpstreamClientV2';
import { NamespacesContext } from 'GlooIApp';
import { Upstream } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/upstream_pb';
import { useListVirtualServices } from 'Api';
import { ListVirtualServicesRequest } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/virtualservice_pb';
import { VirtualService } from 'proto/github.com/solo-io/gloo/projects/gateway/api/v1/virtual_service_pb';

const Container = styled.div`
  ${CardCSS};
  display: flex;
  flex-direction: column;
  background: white;
  width: 100%;
  padding: 30px ${soloConstants.buffer}px ${soloConstants.buffer}px;
`;

const Header = styled.div`
  display: flex;
  justify-content: space-between;
  height: 50px;
  width: 100%;
  margin-bottom: ${soloConstants.smallBuffer}px;
  color: ${colors.novemberGrey};
`;
const PageTitle = styled.div`
  font-size: 22px;
  line-height: 26px;
`;
const PageSubtitle = styled.div`
  font-size: 18px;
  line-height: 22px;
`;

const Row = styled.div`
  display: grid;
  width: 100%;
  grid-template-columns: minmax(200px, 50%) minmax(200px, 50%);
  grid-gap: ${soloConstants.largeBuffer}px;
`;

const EnvoyHealth = styled.div`
  width: 100%;
  margin-bottom: ${soloConstants.largeBuffer}px;
`;
const EnvoyHealthContent = styled.div`
  display: flex;
  justify-content: space-between;
`;
const EnvoyHealthHeader = styled.div`
  max-width: 400px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  padding-right: 20px;
`;
const EnvoyHealthTitle = styled.div`
  display: flex;
  align-items: center;
  font-size: 20px;
  line-height: 24px;
  margin-bottom: 10px;

  > div {
    margin-left: 0;
    margin-right: 10px;
  }
`;
const EnvoyHealthSubtitle = styled.div`
  font-size: 16px;
  line-height: 19px;
`;

const Link = styled.div`
  cursor: pointer;
  color: ${colors.seaBlue};
  font-size: 14px;
`;

const VirtualServices = styled.div`
  width: 100%;
`;
const Upstreams = styled.div`
  width: 100%;
`;

const HealthScoreContainer = styled<'div', { health: number }>('div')`
  display: flex;
  flex-direction: row;
  justify-content: center;
  align-items: center;
  font-size: 28px;
  font-weight: 500;

  .health-icon {
    margin: 0 12px;

    ${props =>
      props.health === healthConstants.Good.value
        ? '' //`fill: ${colors.forestGreen};`
        : props.health === healthConstants.Error.value
        ? `fill: ${colors.grapefruitOrange};`
        : `fill: ${colors.sunGold};`}
  }

  & span {
    ${props =>
      props.health === healthConstants.Good.value
        ? '' //`fill: ${colors.forestGreen};`
        : props.health === healthConstants.Error.value
        ? `color: ${colors.grapefruitOrange};`
        : `color: ${colors.sunGold};`}
    padding: 5px;
    font-weight: bold;
  }
`;

interface Props extends RouteComponentProps {}

export const Overview = (props: Props) => {
  return (
    <React.Fragment>
      <Container>
        <Header>
          <div>
            <PageTitle>Enterprise Gloo Overview</PageTitle>
            <PageSubtitle>
              Your current configuration health at a glance
            </PageSubtitle>
          </div>
          <HealthScoreContainer health={healthConstants.Good.value}>
            <HealthScoreIcon />
            {/*Health Score: <span>92</span>*/}
          </HealthScoreContainer>
        </Header>
        <HealthStatus {...props} />
        <Row>
          <VirtualServicesOverview />
          <UpstreamsOverview />
        </Row>
      </Container>
    </React.Fragment>
  );
};

const HealthStatus = (props: Props) => {
  const goToEnvoys = (): void => {
    props.history.push('/admin/envoy/');
  };

  return (
    <EnvoyHealth>
      <StatusTile titleIcon={<EnvoyIcon />} horizontal>
        <EnvoyHealthContent>
          <EnvoyHealthHeader>
            <div>
              <EnvoyHealthTitle>
                <HealthIndicator healthStatus={healthConstants.Good.value} />{' '}
                Envoy Health Status
              </EnvoyHealthTitle>
              <EnvoyHealthSubtitle>
                Gloo is responsible for configuring Envoy. Whenever Virtual Services or other configs change that affect the proxy,
                Gloo will immediately detect that change and update Envoy's configuration. 
              </EnvoyHealthSubtitle>
            </div>
            <Link onClick={goToEnvoys}>View Envoy Configuration</Link>
          </EnvoyHealthHeader>

          <div>
            <TallyInformationDisplay
              tallyCount={10}
              tallyDescription={'envoy configuration needs your attention'}
              color='orange'
              moreInfoLink={{
                prompt: 'View envoy issues',
                link: '/admin/envoy/?status=Rejected'
              }}
            />

            <TallyInformationDisplay
              tallyCount={10}
              tallyDescription={
                'envoy configurations currently deployed and configured'
              }
              color='blue'
            />
          </div>
        </EnvoyHealthContent>
      </StatusTile>
    </EnvoyHealth>
  );
};

const VirtualServicesOverview = () => {
  const namespaces = React.useContext(NamespacesContext);
  const request = new ListVirtualServicesRequest();
  request.setNamespacesList(namespaces.namespacesList);
  const { data, loading } = useListVirtualServices(request);
  const [virtualServicesList, setVirtualServicesList] = React.useState<
    VirtualService.AsObject[]
  >([]);

  const [
    virtualServiceForRouteCreation,
    setVirtualServiceForRouteCreation
  ] = React.useState<VirtualService.AsObject | undefined>(undefined);
  React.useEffect(() => {
    if (data) {
      setVirtualServicesList(data.virtualServicesList);
    }
  }, [loading]);

  const virtualServiceErrorCount = virtualServicesList.reduce(
    (total, vs) =>
      total +
      (!(vs.status && vs.status.state !== healthConstants.Error.value) ? 1 : 0),
    0
  );

  return (
    <VirtualServices>
      <StatusTile
        titleText={'Virtual Services'}
        titleIcon={<VSIcon />}
        description={
          'Virtual Services define a set of route rules for a given set of domains.'
        }
        exploreMoreLink={{
          prompt: 'View Virtual Services',
          link: '/virtualservices/'
        }}>
        {!!virtualServicesList.length ? (
          <React.Fragment>
            {!!virtualServiceErrorCount ? (
              <TallyInformationDisplay
                tallyCount={virtualServiceErrorCount}
                tallyDescription={'virtual services need your attention'}
                color='orange'
                moreInfoLink={{
                  prompt: 'View virtual service issues',
                  link: '/virtualservices/table?status=Rejected'
                }}
              />
            ) : (
              <GoodStateCongratulations typeOfItem={'virtual services'} />
            )}
            <TallyInformationDisplay
              tallyCount={virtualServicesList.length}
              tallyDescription={
                'virtual services currently installed and configured'
              }
              color='blue'
            />
          </React.Fragment>
        ) : (
          <div>You have no virtual services configured yet.</div>
        )}
      </StatusTile>
    </VirtualServices>
  );
};

const UpstreamsOverview = () => {
  const namespaces = React.useContext(NamespacesContext);

  const [upstreamsList, setUpstreamsList] = React.useState<Upstream.AsObject[]>(
    []
  );
  const { data, loading, error } = useGetUpstreamsListV2({
    namespaces: namespaces.namespacesList
  });

  React.useEffect(() => {
    if (data && data.toObject().upstreamsList) {
      setUpstreamsList(data.toObject().upstreamsList);
    }
  }, [loading]);

  const upstreamErrorCount = upstreamsList.reduce(
    (total, upstream) =>
      total +
      (!(
        upstream.status && upstream.status.state !== healthConstants.Error.value
      )
        ? 1
        : 0),
    0
  );
  return (
    <Upstreams>
      <StatusTile
        titleText={'Upstreams'}
        titleIcon={<USIcon />}
        description={
          'Upstreams define destinations for routes.'
        }
        exploreMoreLink={{
          prompt: 'View Upstreams',
          link: '/upstreams/'
        }}>
        {!!upstreamsList.length ? (
          <React.Fragment>
            {!!upstreamErrorCount ? (
              <TallyInformationDisplay
                tallyCount={upstreamErrorCount}
                tallyDescription={'upstreams need your attention'}
                color='orange'
                moreInfoLink={{
                  prompt: 'View upstream issues',
                  link: '/upstreams/table?status=Rejected'
                }}
              />
            ) : (
              <GoodStateCongratulations typeOfItem={'upstreams'} />
            )}
            <TallyInformationDisplay
              tallyCount={upstreamsList.length}
              tallyDescription={'upstreams currently installed and configured'}
              color='blue'
            />
          </React.Fragment>
        ) : (
          <div>You have no upstreams configured yet.</div>
        )}
      </StatusTile>
    </Upstreams>
  );
};
