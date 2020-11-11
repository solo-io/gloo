import styled from '@emotion/styled';
import { ReactComponent as EnvoyLogo } from 'assets/envoy-logo.svg';
import { ReactComponent as GatewayConfigLogo } from 'assets/gateway-config-icon.svg';
import { ReactComponent as HealthScoreIcon } from 'assets/health-score-icon.svg';
import { ReactComponent as ProxyConfigLogo } from 'assets/proxy-config-icon.svg';
import { ReactComponent as WatchedNamespaceIcon } from 'assets/watched-namespace-icon.svg';
import { ReactComponent as SecretsIcon } from 'assets/secrets-icon.svg';
import { ReactComponent as SecurityIcon } from 'assets/key-on-ring.svg';

import { GoodStateCongratulations } from 'Components/Common/DisplayOnly/GoodStateCongratulations';
import { StatusTile } from 'Components/Common/DisplayOnly/StatusTile';
import { TallyInformationDisplay } from 'Components/Common/DisplayOnly/TallyInformationDisplay';
import React from 'react';
import { RouteProps } from 'react-router';
import { envoyAPI } from 'store/envoy/api';
import { gatewayAPI } from 'store/gateway/api';
import { proxyAPI } from 'store/proxy/api';
import { colors, healthConstants, soloConstants } from 'Styles';
import { CardCSS } from 'Styles/CommonEmotions/card';
import useSWR from 'swr';
import { getResourceStatus } from 'utils/helpers';
import { css } from '@emotion/core';
import { configAPI } from 'store/config/api';
import { ErrorBoundary } from '../Errors/ErrorBoundary';
import { secretAPI } from 'store/secrets/api';

export const Container = styled.div`
  ${CardCSS};
  display: flex;
  flex-direction: column;
  background: white;
  width: 100%;
  padding: 30px ${soloConstants.buffer}px ${soloConstants.buffer}px;
`;

export const Header = styled.div`
  display: flex;
  justify-content: space-between;
  height: 50px;
  width: 100%;
  margin-bottom: ${soloConstants.smallBuffer}px;
  color: ${colors.novemberGrey};
`;

export const PageTitle = styled.div`
  font-size: 22px;
  line-height: 26px;
`;

export const PageSubtitle = styled.div`
  font-size: 18px;
  line-height: 22px;
`;

export const Row = styled.div`
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

type HealthScoreContainerProps = { health: number };
export const HealthScoreContainer = styled.div`
  display: flex;
  flex-direction: row;
  justify-content: center;
  align-items: center;
  font-size: 28px;
  font-weight: 500;

  .health-icon {
    margin: 0 12px;

    ${(props: HealthScoreContainerProps) =>
      props.health === healthConstants.Good.value
        ? '' //`fill: ${colors.forestGreen};`
        : props.health === healthConstants.Error.value
        ? `fill: ${colors.grapefruitOrange};`
        : `fill: ${colors.sunGold};`}
  }
`;

export const AdminLanding: React.FC<RouteProps> = props => {
  const { data: licenseData, error: licenseError } = useSWR(
    'hasValidLicense',
    configAPI.getIsLicenseValid,
    { refreshInterval: 0 }
  );
  const { data: settingsDetails, error: settingsError } = useSWR(
    'getSettings',
    configAPI.getSettings
  );

  const { data: secretsList, error: secretsError } = useSWR(
    'listSecrets',
    secretAPI.getSecretsList
  );

  const { data: watchedNamespacesList, error: watchedNamespacesError } = useSWR(
    'listNamespaces',
    configAPI.listNamespaces
  );
  if (!settingsDetails || !secretsList || !watchedNamespacesList) {
    return <div>Loading...</div>;
  }

  let glooEdition = licenseData?.isLicenseValid
    ? 'Gloo Edge Enterprise'
    : 'Gloo Edge';

  return (
    <ErrorBoundary
      fallback={<div>There was an error with the Admin section</div>}>
      <Container>
        <Header>
          <div>
            <PageTitle>{`${glooEdition} Administration`}</PageTitle>
            <PageSubtitle>
              {`Advanced Administration for your ${glooEdition} Configuration`}
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
        <div>
          <div className='mt-2 text-2xl text-gray-900'> Settings</div>
        </div>
        <div
          css={css`
            height: 250px;
          `}
          className='grid grid-cols-3 '>
          <StatusTile
            titleIcon={
              <span className='text-blue-500'>
                <SecurityIcon className='fill-current ' />
              </span>
            }
            titleText='Security'
            exploreMoreLink={{
              prompt: 'View Setttings',
              link: `/admin/settings/${settingsDetails.settings?.metadata?.namespace}/${settingsDetails.settings?.metadata?.name}`
            }}
            description={`Represents global settings for all of Gloo's components.`}></StatusTile>
          <StatusTile
            titleIcon={<WatchedNamespaceIcon className='w-8 h-8' />}
            exploreMoreLink={{
              prompt: 'View Watched Namespaces',
              link: '/admin/watched-namespaces'
            }}
            titleText='Watched Namespaces'
            description='Use this setting to restrict the namespaces that Gloo Edge controllers take into consideration when watching for resources.In a usual production scenario, RBAC policies will limit the namespaces that Gloo Edge has access to. '></StatusTile>
          <StatusTile
            titleIcon={
              <span className='text-blue-500'>
                <SecretsIcon className='fill-current ' />
              </span>
            }
            exploreMoreLink={{
              prompt: 'View Secrets',
              link: '/admin/secrets/'
            }}
            titleText='Secrets'
            description='Certain features such as the AWS Lambda option require the use of secrets for authentication, configuration of SSL Certificates, and other data that should not be stored in plaintext configuration. Gloo Edge runs an independent (goroutine) controller to monitor secrets. Secrets are stored in their own secret storage layer.'></StatusTile>
        </div>
      </Container>
    </ErrorBoundary>
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
  const { data: gatewaysList, error } = useSWR(
    'listGateways',
    gatewayAPI.listGateways
  );

  if (!gatewaysList) {
    return <div>Loading...</div>;
  }

  const gatewayErrorCount = gatewaysList.reduce((total, gateway) => {
    if (getResourceStatus(gateway!.gateway!.status?.state!) !== 'Rejected') {
      return total;
    }

    return total + 1;
  }, 0);

  return (
    <Gateway>
      <StatusTile
        titleText={'Gateway Configuration'}
        titleIcon={<GatewayConfigLogo />}
        description={
          'Gateways are used to configure the protocols and ports for Envoy. Optionally, gateways can be associated with a specific set of virtual services.'
        }
        exploreMoreLink={{
          prompt: 'View Gateways',
          link: '/admin/gateways/'
        }}
        healthStatus={
          !!gatewayErrorCount
            ? healthConstants.Error.value
            : healthConstants.Good.value
        }
        descriptionMinHeight={'95px'}>
        {!!gatewaysList.length ? (
          <>
            {!!gatewayErrorCount ? (
              <TallyInformationDisplay
                tallyCount={gatewayErrorCount}
                tallyDescription={`gateway error${
                  gatewayErrorCount === 1 ? '' : 's'
                }`}
                color='orange'
                moreInfoLink={{
                  prompt: 'View',
                  link: '/admin/gateways/?status=Rejected'
                }}
              />
            ) : (
              <GoodStateCongratulations typeOfItem={'gateway configurations'} />
            )}
            <TallyInformationDisplay
              tallyCount={gatewaysList.length}
              tallyDescription={`gateway configuration${
                gatewaysList.length === 1 ? '' : 's'
              } `}
              color='blue'
            />
          </>
        ) : (
          <div>You have no gateways configured yet.</div>
        )}
      </StatusTile>
    </Gateway>
  );
};

const ProxyOverview = () => {
  const { data: proxiesList, error } = useSWR(
    'listProxies',
    proxyAPI.getListProxies
  );

  if (!proxiesList) {
    return <div>Loading...</div>;
  }
  let proxyStatus = healthConstants.Pending.value;
  let proxyErrorCount = 0;
  proxiesList.forEach(proxy => {
    if (proxy.status!.code === 0) {
      proxyStatus = healthConstants.Error.value;
      proxyErrorCount += 1;
    } else if (proxy.status!.code === 2) {
      proxyStatus = healthConstants.Good.value;
    }
  });
  return (
    <Proxy>
      <StatusTile
        titleText={'Proxy Configuration'}
        titleIcon={<ProxyConfigLogo />}
        description={
          'Gloo Edge generates proxy configs from upstreams, virtual services, and gateways, and then transforms them directly into Envoy config. If a proxy config is rejected, it means Envoy will not receive configuration updates.'
        }
        exploreMoreLink={{
          prompt: 'View Proxy',
          link: '/admin/proxy/'
        }}
        healthStatus={proxyStatus}
        descriptionMinHeight={'95px'}>
        {!!proxiesList.length ? (
          <>
            {!!proxyErrorCount ? (
              <TallyInformationDisplay
                tallyCount={proxyErrorCount}
                tallyDescription={`proxy error${
                  proxyErrorCount === 1 ? '' : 's'
                }`}
                color='orange'
                moreInfoLink={{
                  prompt: 'View',
                  link: '/admin/proxy/?status=Rejected'
                }}
              />
            ) : (
              <GoodStateCongratulations typeOfItem={'proxy configurations'} />
            )}
            <TallyInformationDisplay
              tallyCount={proxiesList.length}
              tallyDescription={`proxy configuration${
                proxiesList.length === 1 ? '' : 's'
              } `}
              color='blue'
            />
          </>
        ) : (
          <div>You have no proxy configured yet.</div>
        )}
      </StatusTile>
    </Proxy>
  );
};

const EnvoyOverview = () => {
  const { data: envoysList, error } = useSWR(
    'listEnvoys',
    envoyAPI.getEnvoyList
  );

  if (!envoysList) {
    return <div>Loading...</div>;
  }

  let envoyStatus = healthConstants.Pending.value;
  let envoyErrorCount = 0;
  envoysList.forEach(envoy => {
    if (envoy.status!.code === 0) {
      envoyStatus = healthConstants.Pending.value;
      envoyErrorCount += 1;
    } else if (envoy.status!.code === 2) {
      envoyStatus = healthConstants.Good.value;
    }
  });

  return (
    <Envoy>
      <StatusTile
        titleText={'Envoy Configuration'}
        titleIcon={
          <EnvoyLogo
            css={css`
              height: 32px;
            `}
          />
        }
        description={
          'This is the live config dump from Envoy. This is translated directly from the proxy config and should be updated any time the proxy configuration changes.'
        }
        exploreMoreLink={{
          prompt: 'View Envoy',
          link: '/admin/envoy/'
        }}
        healthStatus={envoyStatus}
        descriptionMinHeight={'95px'}>
        {!envoysList.length ? (
          <div>Loading...</div>
        ) : !!envoysList.length ? (
          <>
            {envoyStatus === healthConstants.Error.value ? (
              <TallyInformationDisplay
                tallyCount={envoyErrorCount}
                tallyDescription={`envoy error${
                  envoyErrorCount === 1 ? '' : 's'
                }`}
                color='orange'
                moreInfoLink={{
                  prompt: 'View',
                  link: '/admin/envoy/?status=Rejected'
                }}
              />
            ) : (
              envoyStatus === healthConstants.Good.value && (
                <GoodStateCongratulations typeOfItem={'envoy configurations'} />
              )
            )}
            {envoyStatus === healthConstants.Pending.value && (
              <TallyInformationDisplay
                tallyCount={envoysList.length}
                tallyDescription={`envoy${
                  envoysList.length === 1 ? '' : 's'
                } configuration pending`}
                color='yellow'
                moreInfoLink={{
                  prompt: 'View issues',
                  link: '/admin/envoy/?status=Rejected'
                }}
              />
            )}
            {envoyStatus === healthConstants.Good.value && (
              <TallyInformationDisplay
                tallyCount={envoysList.length}
                tallyDescription={`envoy${
                  envoysList.length === 1 ? '' : 's'
                } configured`}
                color='blue'
              />
            )}
          </>
        ) : (
          <div>You have no envoy configured yet.</div>
        )}
      </StatusTile>
    </Envoy>
  );
};
