import React from 'react';
import styled from '@emotion/styled/macro';
import { RouteProps } from 'react-router-dom';
import { ReactComponent as EnvoyLogo } from 'assets/envoy-logo.svg';
import { ReactComponent as HealthScoreIcon } from 'assets/health-score-icon.svg';
import { ReactComponent as VSIcon } from 'assets/vs-icon.svg';
import { ReactComponent as USIcon } from 'assets/us-icon.svg';
import { colors } from 'Styles/colors';

const Container = styled.div`
  display: grid;
  background-color: white;
  grid-template-areas:
    'h h'
    'envoy envoy'
    'virtualservices upstreams';
  box-shadow: 0px 4px 9px #e5e5e5;
  grid-template-columns: minmax(200px, 600px) minmax(200px, 600px);
  grid-template-rows: 1fr 200px 350px;
  grid-gap: 15px;
  margin: 18px 88px;
  padding: 19px;
`;

const Header = styled.div`
  grid-area: h;
  height: 100px;
  background-color: white;
  display: flex;
  justify-content: space-between;
`;
const EnvoyHealth = styled.div`
  grid-area: envoy;
  background-color: #f7f7f7;
  padding: 18px;
`;
const VirtualServices = styled.div`
  grid-area: virtualservices;
  background-color: white;
`;

const Upstreams = styled.div`
  grid-area: upstreams;
  background-color: white;
`;

const HealthScoreContainer = styled.div`
  display: flex;
  flex-direction: row;
  justify-content: center;
  align-items: center;
  font-size: 28px;
  font-weight: 500;

  .health-icon {
    fill: ${colors.forestGreen};
  }
  & span {
    color: ${colors.forestGreen};
    padding: 5px;
    font-weight: bold;
  }
`;

export const Overview: React.FC<RouteProps> = props => {
  return (
    <React.Fragment>
      <Container>
        <Header>
          <div>
            <div style={{ fontSize: '22px' }}>Enterprise Gloo Overview</div>
            <div style={{ fontSize: '18px' }}>
              Your current configuration health at a glance
            </div>
          </div>
          <HealthScoreContainer>
            <HealthScoreIcon style={{ margin: '12px' }} />
            Health Score: <span>92</span>
          </HealthScoreContainer>
        </Header>
        <HealthStatus />
        <VirtualServicesOverview />
        <UpstreamsOverview />
      </Container>
    </React.Fragment>
  );
};
const HealthStatusContainer = styled.div`
  display: grid;
  grid-template-columns: 175px 1fr;
  background-color: white;
`;
const EnvoyStatus = styled.div`
  display: flex;
  margin: 18px;
  background-color: #f7f7f7;
`;

const HealthStatus = () => {
  return (
    <EnvoyHealth>
      <HealthStatusContainer>
        <EnvoyLogo />
        <EnvoyStatus>
          <div>
            Envoy Health Status
            <div>[healthIcon]</div>
            [description]
            <div>View Evnvoy Configuration</div>
          </div>
          <div>Health Summary</div>
        </EnvoyStatus>
      </HealthStatusContainer>
    </EnvoyHealth>
  );
};

const VirtualServicesOverview = () => {
  return (
    <VirtualServices>
      Virtual Services <VSIcon />
      <div>description</div>
      <div>summary</div>
      <div>view VirtuealServces</div>
    </VirtualServices>
  );
};
const UpstreamsOverview = () => {
  return (
    <Upstreams>
      Upstreams <USIcon />
      <div>description</div>
      <div>summary</div>
      <div>View Upstreams</div>
    </Upstreams>
  );
};
