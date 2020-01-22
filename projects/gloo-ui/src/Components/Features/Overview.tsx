import styled from '@emotion/styled';
import { ReactComponent as EnvoyIcon } from 'assets/envoy-logo-title.svg';
import { ReactComponent as HealthScoreIcon } from 'assets/health-score-icon.svg';
import { ReactComponent as USIcon } from 'assets/upstream-icon.svg';
import { ReactComponent as VSIcon } from 'assets/virtualservice-icon.svg';
import { GoodStateCongratulations } from 'Components/Common/DisplayOnly/GoodStateCongratulations';
import { StatusTile } from 'Components/Common/DisplayOnly/StatusTile';
import { TallyInformationDisplay } from 'Components/Common/DisplayOnly/TallyInformationDisplay';
import { HealthIndicator } from 'Components/Common/HealthIndicator';
import { Upstream } from 'proto/gloo/projects/gloo/api/v1/upstream_pb';
import React from 'react';
import { useHistory } from 'react-router';
import { upstreamAPI } from 'store/upstreams/api';
import { virtualServiceAPI } from 'store/virtualServices/api';
import { healthConstants, soloConstants } from 'Styles';
import { colors } from 'Styles/colors';
import { CardCSS } from 'Styles/CommonEmotions/card';
import useSWR from 'swr';
import { getIcon, getUpstreamType, groupBy } from 'utils/helpers';
import { envoyAPI } from '../../store/envoy/api';

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
  display: grid;
  grid-template-columns: 1fr 1fr;
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

type HealthProps = { health: number };
const HealthScoreContainer = styled.div`
  display: flex;
  flex-direction: row;
  justify-content: center;
  align-items: center;
  font-size: 28px;
  font-weight: 500;

  .health-icon {
    margin: 0 12px;

    ${(props: HealthProps) =>
      props.health === healthConstants.Good.value
        ? '' //`fill: ${colors.forestGreen};`
        : props.health === healthConstants.Error.value
        ? `fill: ${colors.grapefruitOrange};`
        : `fill: ${colors.sunGold};`}
  }

  & span {
    ${(props: HealthProps) =>
      props.health === healthConstants.Good.value
        ? '' //`fill: ${colors.forestGreen};`
        : props.health === healthConstants.Error.value
        ? `color: ${colors.grapefruitOrange};`
        : `color: ${colors.sunGold};`}
    padding: 5px;
    font-weight: bold;
  }
`;

export const Overview = () => {
  return (
    <>
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
        <HealthStatus />
        <Row>
          <VirtualServicesOverview />
          <UpstreamsOverview />
        </Row>
      </Container>
    </>
  );
};

const HealthStatus = () => {
  let history = useHistory();
  const { data: envoyList, error } = useSWR(
    'listEnvoys',
    envoyAPI.getEnvoyList
  );

  if (!envoyList) {
    return <div>Loading...</div>;
  }

  const goToAdmin = (): void => {
    history.push('/admin/');
  };

  let envoyStatus = healthConstants.Pending.value;
  let envoyErrorCount = 0;
  envoyList.forEach(envoy => {
    if (envoy.status!.code === 0) {
      envoyStatus = healthConstants.Pending.value;
      envoyErrorCount += 1;
    } else if (envoy.status!.code === 2) {
      envoyStatus = healthConstants.Good.value;
    }
  });

  return (
    <EnvoyHealth>
      <StatusTile titleIcon={<EnvoyIcon />} horizontal>
        <EnvoyHealthContent>
          <EnvoyHealthHeader>
            <div>
              <EnvoyHealthTitle>
                <HealthIndicator healthStatus={envoyStatus} /> Envoy Health
                Status
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

          {!envoyList.length ? (
            <div>Loading...</div>
          ) : !!envoyList.length ? (
            <div>
              {envoyStatus === healthConstants.Error.value ? (
                <TallyInformationDisplay
                  tallyCount={envoyErrorCount}
                  tallyDescription={`envoy configuration error${
                    envoyErrorCount === 1 ? '' : 's'
                  }`}
                  color='orange'
                  moreInfoLink={{
                    prompt: 'View envoy issues',
                    link: '/admin/envoy/?status=Rejected'
                  }}
                />
              ) : (
                envoyStatus === healthConstants.Good.value && (
                  <GoodStateCongratulations typeOfItem={'envoys'} />
                )
              )}
              {envoyStatus === healthConstants.Pending.value && (
                <TallyInformationDisplay
                  tallyCount={envoyList.length}
                  tallyDescription={`envoy${
                    envoyList.length === 1 ? '' : 's'
                  } configuration pending`}
                  color='yellow'
                  moreInfoLink={{
                    prompt: 'View envoy issues',
                    link: '/admin/envoy/?status=Rejected'
                  }}
                />
              )}
              {envoyStatus === healthConstants.Good.value && (
                <TallyInformationDisplay
                  tallyCount={envoyList.length}
                  tallyDescription={`envoy${
                    envoyList.length === 1 ? '' : 's'
                  } configured`}
                  color='blue'
                />
              )}
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
  const { data: virtualServicesList, error } = useSWR(
    'listVirtualServices',
    virtualServiceAPI.listVirtualServices
  );
  if (!virtualServicesList) {
    return <div>Loading...</div>;
  }

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
          link: '/virtualservices/',
          testId: 'view-virtual-services-link'
        }}>
        {!!virtualServicesList.length ? (
          <>
            {!!virtualServiceErrorCount ? (
              <TallyInformationDisplay
                tallyCount={virtualServiceErrorCount}
                tallyDescription={`virtual services error${
                  virtualServiceErrorCount === 1 ? '' : 's'
                }`}
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
              tallyDescription={`virtual service${
                virtualServicesList.length === 1 ? '' : 's'
              } configured`}
              color='blue'
            />
          </>
        ) : (
          <div>
            <TallyInformationDisplay
              tallyCount={null}
              tallyDescription={`Envoy will not receive configuration updates until a virtual service is defined.`}
              color={'yellow'}
            />
          </div>
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
const UpstreamDetails: React.FC<UpstreamDetailsProps> = ({
  upstreamsList = []
}) => {
  let groupedUS = Array.from(
    groupBy(upstreamsList, us => getUpstreamType(us)).entries()
  );

  return (
    <UpstreamDetailsContainer>
      {groupedUS.map(([upstreamType, usList]) => (
        <UpstreamDetail key={upstreamType}>
          <IconContainer>{getIcon(upstreamType)}</IconContainer>
          <div>
            <b>{`${usList.length}`}</b>
            {`  ${upstreamType} 
              upstream${usList.length === 1 ? '' : 's'}`}
          </div>
        </UpstreamDetail>
      ))}
    </UpstreamDetailsContainer>
  );
};

const UpstreamsOverview = () => {
  const { data: upstreamsList, error } = useSWR(
    'listUpstreams',
    upstreamAPI.listUpstreams
  );

  if (!upstreamsList) {
    return <div>Loading...</div>;
  }
  const upstreamErrorCount = upstreamsList.reduce(
    (total, upstream) =>
      total +
      (!(
        upstream?.upstream?.status &&
        upstream?.upstream.status.state !== healthConstants.Error.value
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
          link: '/upstreams/',
          testId: 'view-upstreams-link'
        }}>
        {!!upstreamsList.length ? (
          <>
            {!!upstreamErrorCount ? (
              <TallyInformationDisplay
                tallyCount={upstreamErrorCount}
                tallyDescription={`upstream error${
                  upstreamErrorCount === 1 ? '' : 's'
                } `}
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
              tallyDescription={`upstream${
                upstreamsList.length === 1 ? '' : 's'
              } configured`}
              color='blue'
            />
            <UpstreamDetails
              upstreamsList={upstreamsList.map(
                upstreamDetails => upstreamDetails.upstream!
              )}
            />
          </>
        ) : (
          <div>You have no upstreams configured yet.</div>
        )}
      </StatusTile>
    </Upstreams>
  );
};
