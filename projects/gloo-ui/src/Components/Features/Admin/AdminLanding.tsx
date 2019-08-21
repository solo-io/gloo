import React from 'react';
import styled from '@emotion/styled/macro';
import { RouteProps, Route } from 'react-router';
import { Breadcrumb } from 'Components/Common/Breadcrumb';
import { healthConstants, colors, soloConstants } from 'Styles';
import { CardCSS } from 'Styles/CommonEmotions/card';
import { HealthIndicator } from 'Components/Common/HealthIndicator';
import { ReactComponent as HealthScoreIcon } from 'assets/health-score-icon.svg';
import { ReactComponent as GatewayConfigLogo } from 'assets/gateway-config-icon.svg';
import { ReactComponent as EnvoyLogo } from 'assets/envoy-logo.svg';
import { ReactComponent as ProxyConfigLogo } from 'assets/proxy-config-icon.svg';
import { StatusTile } from 'Components/Common/DisplayOnly/StatusTile';
import { TallyInformationDisplay } from 'Components/Common/DisplayOnly/TallyInformationDisplay';
import { GoodStateCongratulations } from 'Components/Common/DisplayOnly/GoodStateCongratulations';
import { useGetGatewayList } from 'Api/v2/useGatewayClientV2';
import { GatewayDetails } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/gateway_pb';
import { useGetProxiesList } from 'Api/v2/useProxyClientV2';
import { ProxyDetails } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/proxy_pb';
import { getResourceStatus } from 'utils/helpers';
import { AppState } from 'store';
import { useSelector } from 'react-redux';

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
  grid-template-columns: minmax(200px, 33.3%) minmax(200px, 33.3%) minmax(
      200px,
      33.3%
    );
  grid-gap: 17px;
  padding: 0 17px;
  border-radius: ${soloConstants.smallRadius}px;
  background: ${colors.januaryGrey};

  > div > div {
    /*Status Tiles*/
    padding: 17px 0;
  }
`;

const Link = styled.div`
  cursor: pointer;
  color: ${colors.seaBlue};
  font-size: 14px;
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

export const AdminLanding: React.FC<RouteProps> = props => {
  return (
    <React.Fragment>
      <Container>
        <Header>
          <div>
            <PageTitle>Enterprise Gloo Administration</PageTitle>
            <PageSubtitle>
              Advanced Administration for youur Gloo Configuration
            </PageSubtitle>
          </div>
          <HealthScoreContainer health={healthConstants.Good.value}>
            <HealthScoreIcon />
            {/*Health Score: <span>92</span>*/}
          </HealthScoreContainer>
        </Header>
        <Row>
          <GatewayOverview />
          <ProxyOverview />
          <EnvoyOverview />
        </Row>
      </Container>
    </React.Fragment>
  );
};

const Gateway = styled.div`
  width: 100%;
`;
const Proxy = styled.div`
  width: 100%;
`;
const Envoy = styled.div`
  width: 100%;
`;

const GatewayOverview = () => {
  const {
    config: { namespacesList }
  } = useSelector((state: AppState) => state);
  const { data, loading, error, setNewVariables } = useGetGatewayList({
    namespaces: namespacesList
  });
  const [allGateways, setAllGateways] = React.useState<
    GatewayDetails.AsObject[]
  >([]);

  React.useEffect(() => {
    if (!!data) {
      const newGateways = data
        .toObject()
        .gatewayDetailsList.filter(gateway => !!gateway.gateway);
      setAllGateways(newGateways);
    }
  }, [loading]);

  if (!data || (!data && loading)) {
    return <div>Loading...</div>;
  }

  const gatewayErrorCount = allGateways.reduce((total, gateway) => {
    if (getResourceStatus(gateway.gateway!) !== 'Rejected') {
      return total;
    }

    return total + 1;
  }, 0);

  return (
    <Gateway>
      <StatusTile
        titleText={'Gateway Configuration'}
        titleIcon={<GatewayConfigLogo />}
        description={'Gateways are used to configure the protocols and ports for Envoy. Optionally, gateways can be associated with a specific set of virtual services.'}
        exploreMoreLink={{
          prompt: 'View Gateways',
          link: '/admin/gateways/'
        }}
        healthStatus={
          !!gatewayErrorCount
            ? healthConstants.Error.value
            : healthConstants.Good.value
        }>
        {!!allGateways.length ? (
          <React.Fragment>
            {!!gatewayErrorCount ? (
              <TallyInformationDisplay
                tallyCount={gatewayErrorCount}
                tallyDescription={'gateway configurations need your attention'}
                color='orange'
                moreInfoLink={{
                  prompt: 'View gateway issues',
                  link: '/admin/gateways/?status=Rejected'
                }}
              />
            ) : (
              <GoodStateCongratulations typeOfItem={'gateway configurations'} />
            )}
            <TallyInformationDisplay
              tallyCount={allGateways.length}
              tallyDescription={
                'gateway configurations currently deployed'
              }
              color='blue'
            />
          </React.Fragment>
        ) : (
          <div>You have no gateways configured yet.</div>
        )}
      </StatusTile>
    </Gateway>
  );
};

const ProxyOverview = () => {
  const {
    config: { namespacesList }
  } = useSelector((state: AppState) => state);
  const { data, loading, error, setNewVariables } = useGetProxiesList({
    namespaces: namespacesList
  });
  const [allProxies, setAllProxies] = React.useState<ProxyDetails.AsObject[]>(
    []
  );

  React.useEffect(() => {
    if (!!data) {
      const newProxies = data
        .toObject()
        .proxyDetailsList.filter(proxy => proxy.proxy);
      setAllProxies(newProxies);
    }
  }, [loading]);

  if (!data || (!data && loading)) {
    return <div>Loading...</div>;
  }

  const proxyErrorCount = allProxies.reduce((total, proxy) => {
    if (getResourceStatus(proxy.proxy!) !== 'Rejected') {
      return total;
    }

    return total + 1;
  }, 0);

  return (
    <Proxy>
      <StatusTile
        titleText={'Proxy Configuration'}
        titleIcon={<ProxyConfigLogo />}
        description={'Gloo generates proxy configs from upstreams, virtual services, and gateways, and then transforms them directly into Envoy config. If a proxy config is rejected, it means Envoy will not receive configuration updates.'}
        exploreMoreLink={{
          prompt: 'View Proxy',
          link: '/admin/proxy/'
        }}
        healthStatus={
          !!proxyErrorCount
            ? healthConstants.Error.value
            : healthConstants.Good.value
        }>
        {!!allProxies.length ? (
          <React.Fragment>
            {!!proxyErrorCount ? (
              <TallyInformationDisplay
                tallyCount={proxyErrorCount}
                tallyDescription={'proxy configurations have errors'}
                color='orange'
                moreInfoLink={{
                  prompt: 'View proxy issues',
                  link: '/admin/proxy/?status=Rejected'
                }}
              />
            ) : (
              <GoodStateCongratulations typeOfItem={'proxy configurations'} />
            )}
            <TallyInformationDisplay
              tallyCount={allProxies.length}
              tallyDescription={
                'proxy configurations produced by Gloo'
              }
              color='blue'
            />
          </React.Fragment>
        ) : (
          <div>You have no proxy configured yet.</div>
        )}
      </StatusTile>
    </Proxy>
  );
};

const EnvoyOverview = () => {
  /*
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
    );*/
  return (
    <Envoy>
      <StatusTile
        titleText={'Envoy Configuration'}
        titleIcon={<EnvoyLogo />}
        description={'This is the live config dump from Envoy. This is translated directly from the proxy config and should be updated any time the proxy configuration changes.'}
        exploreMoreLink={{
          prompt: 'View Envoy',
          link: '/admin/envoy/'
        }}
        healthStatus={healthConstants.Error.value}>
        {!!3 ? (
          <React.Fragment>
            {!!3 ? (
              <TallyInformationDisplay
                tallyCount={1}
                tallyDescription={'envoy need your attention'}
                color='orange'
                moreInfoLink={{
                  prompt: 'View envoy issues',
                  link: '/admin/envoy/?status=Rejected'
                }}
              />
            ) : (
              <GoodStateCongratulations typeOfItem={'envoy configurations'} />
            )}
            <TallyInformationDisplay
              tallyCount={1}
              tallyDescription={
                'envoy configurations currently deployed'
              }
              color='blue'
            />
          </React.Fragment>
        ) : (
          <div>You have no envoy configured yet.</div>
        )}
      </StatusTile>
    </Envoy>
  );
};
