import React from 'react';
import styled from '@emotion/styled/macro';
import { RouteProps, RouteComponentProps } from 'react-router-dom';
import { ReactComponent as EnvoyLogo } from 'assets/envoy-logo.svg';
import { ReactComponent as HealthScoreIcon } from 'assets/health-score-icon.svg';
import { ReactComponent as VSIcon } from 'assets/virtualservice-icon.svg';
import { ReactComponent as USIcon } from 'assets/upstream-icon.svg';
import { ReactComponent as EnvoyIcon } from 'assets/envoy-logo-title.svg';
import { colors } from 'Styles/colors';

import { TallyInformationDisplay } from 'Components/Common/DisplayOnly/TallyInformationDisplay';
import { GoodStateCongratulations } from 'Components/Common/DisplayOnly/GoodStateCongratulations';
import { StatusTile } from 'Components/Common/DisplayOnly/StatusTile';
import { soloConstants, healthConstants } from 'Styles';
import { CardCSS } from 'Styles/CommonEmotions/card';
import { HealthIndicator } from 'Components/Common/HealthIndicator';
import { VirtualService } from 'proto/github.com/solo-io/gloo/projects/gateway/api/v1/virtual_service_pb';
import { EnvoyDetails } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/envoy_pb';
import { useSelector, useDispatch } from 'react-redux';
import { AppState } from 'store';
import { listUpstreams } from 'store/upstreams/actions';
import { listVirtualServices } from 'store/virtualServices/actions';
import { groupBy, getUpstreamType, getIcon } from 'utils/helpers';
import { Upstream } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/upstream_pb';
import { getEnvoyHealth } from './Admin/Envoy';

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
  const envoysList = useSelector(
    (state: AppState) => state.envoy.envoyDetailsList
  );
  const [allEnvoy, setAllEnvoy] = React.useState<EnvoyDetails.AsObject[]>([]);

  React.useEffect(() => {
    if (!!envoysList.length) {
      setAllEnvoy(envoysList);
    }
  }, [envoysList.length]);

  const envoyErrorCount = allEnvoy.reduce((total, envoy) => {
    if (
      envoy.status &&
      getEnvoyHealth(envoy.status.code) === healthConstants.Error.value
    ) {
      return total + 1;
    }

    return total;
  }, 0);

  if (!envoysList.length) {
    return <div>Loading...</div>;
  }
  console.log(allEnvoy);

  const goToAdmin = (): void => {
    props.history.push('/admin/');
  };

  return (
    <EnvoyHealth>
      <StatusTile titleIcon={<EnvoyIcon />} horizontal>
        <EnvoyHealthContent>
          <EnvoyHealthHeader>
            <div>
              <EnvoyHealthTitle>
                <HealthIndicator
                  healthStatus={getEnvoyHealth(envoysList[0]!.status!.code)}
                />{' '}
                Envoy Health Status
              </EnvoyHealthTitle>
              <EnvoyHealthSubtitle>
                Gloo is responsible for configuring Envoy. Whenever Virtual
                Services or other configs change that affect the proxy, Gloo
                will immediately detect that change and update Envoy's
                configuration.
              </EnvoyHealthSubtitle>
            </div>
            <Link onClick={goToAdmin}>Go to Admin View</Link>
          </EnvoyHealthHeader>

          {!envoysList.length ? (
            <div>Loading...</div>
          ) : !!allEnvoy.length ? (
            <div>
              {!!envoyErrorCount ? (
                <TallyInformationDisplay
                  tallyCount={envoyErrorCount}
                  tallyDescription={'envoy configuration needs your attention'}
                  color='orange'
                  moreInfoLink={{
                    prompt: 'View envoy issues',
                    link: '/admin/envoy/?status=Rejected'
                  }}
                />
              ) : (
                <GoodStateCongratulations typeOfItem={'envoys'} />
              )}

              <TallyInformationDisplay
                tallyCount={allEnvoy.length}
                tallyDescription={'envoys configured'}
                color='blue'
              />
            </div>
          ) : (
            <div>You have no envoy configurations yet.</div>
          )}
        </EnvoyHealthContent>
      </StatusTile>
    </EnvoyHealth>
  );
};

const VirtualServicesOverview = () => {
  const dispatch = useDispatch();
  const {
    config: { namespacesList },
    virtualServices: { virtualServicesList }
  } = useSelector((state: AppState) => state);
  const [
    virtualServiceForRouteCreation,
    setVirtualServiceForRouteCreation
  ] = React.useState<VirtualService.AsObject | undefined>(undefined);
  React.useEffect(() => {
    if (!virtualServicesList.length) {
      dispatch(listVirtualServices({ namespacesList }));
    }
  }, [virtualServicesList.length]);

  const virtualServiceErrorCount = virtualServicesList.reduce(
    (total, vs) =>
      total +
      (!(
        vs.virtualService!.status &&
        vs.virtualService!.status!.state !== healthConstants.Error.value
      )
        ? 1
        : 0),
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
                'virtual services configured'
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

type UpstreamDetailsProps = { upstreamsList: Upstream.AsObject[] };
const UpstreamDetailsContainer = styled.div`
  display: grid;
  grid-template-rows: 1fr;
  grid-template-columns: repeat(auto-fill, minmax(130px, 1fr));
  padding: 0 10px 10px 10px;
`;

const UpstreamDetail = styled.div`
  display: flex;
  line-height: 1;
  align-items: center;
`;
const IconContainer = styled.div`
  padding: 0 3px;
`;
const UpstreamDetails: React.FC<UpstreamDetailsProps> = props => {
  let groupedUS = Array.from(
    groupBy(props.upstreamsList, us => getUpstreamType(us)).entries()
  );

  return (
    <UpstreamDetailsContainer>
      {groupedUS.map(([upstreamType, usList]) => (
        <UpstreamDetail key={upstreamType}>
          <IconContainer>{getIcon(upstreamType)}</IconContainer>
          <div>
            <b>{`${usList.length}`}</b>
            {`${upstreamType} 
              upstream${usList.length === 1 ? '' : 's'}`}
          </div>
        </UpstreamDetail>
      ))}
    </UpstreamDetailsContainer>
  );
};

const UpstreamsOverview = () => {
  const [isLoading, setIsLoading] = React.useState(false);
  const dispatch = useDispatch();
  const namespacesList = useSelector(
    (state: AppState) => state.config.namespacesList
  );

  const upstreamsList = useSelector((state: AppState) =>
    state.upstreams.upstreamsList.map(u => u.upstream!)
  );
  React.useEffect(() => {
    if (upstreamsList.length) {
      setIsLoading(false);
    } else {
      dispatch(listUpstreams({ namespacesList }));
    }
  }, [upstreamsList.length]);

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
        description={'Upstreams define destinations for routes.'}
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
              tallyDescription={'upstreams configured'}
              color='blue'
            />
            <UpstreamDetails upstreamsList={upstreamsList} />
          </React.Fragment>
        ) : (
          <div>You have no upstreams configured yet.</div>
        )}
      </StatusTile>
    </Upstreams>
  );
};
