import React from 'react';
import styled from '@emotion/styled';
import { CountBox } from 'Components/Common/CountBox';
import {
  CardSubsectionWrapper,
  CardSubsectionContent,
} from 'Components/Common/Card';
import { HealthIndicator } from 'Components/Common/HealthIndicator';
import { colors } from 'Styles/colors';
import { ReactComponent as GatewayIcon } from 'assets/gateway.svg';
import { ReactComponent as ProxyIcon } from 'assets/proxy-icon.svg';
import { ReactComponent as SuccessCircle } from 'assets/big-successful-checkmark.svg';
import { ReactComponent as GearIcon } from 'assets/gear-icon.svg';
import { ReactComponent as EnvoyIcon } from 'assets/envoy-logo.svg';
import { ReactComponent as WatchedNamespacesIcon } from 'assets/watched-namespace-icon.svg';
import { ReactComponent as SecretsIcon } from 'assets/cloud-key-icon.svg';
import { Loading } from 'Components/Common/Loading';
import { PlacementStatus } from 'proto/github.com/solo-io/solo-projects/projects/gloo-fed/api/fed/core/v1/placement_pb';
import { GlooInstanceSpec } from 'proto/github.com/solo-io/solo-projects/projects/gloo-fed/api/fed/v1/instance_pb';
import { Gateway } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/gateway_resources_pb';
import { StatusType } from 'utils/health-status';
import { GatewayStatus } from 'proto/github.com/solo-io/solo-apis/api/gloo/gateway/v1/gateway_pb';
import { UpstreamStatus } from 'proto/github.com/solo-io/solo-apis/api/gloo/gloo/v1/upstream_pb';
import { SimpleLinkProps, SoloLink } from 'Components/Common/SoloLink';
import { ServiceError } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/gateway_resources_pb_service';
import { DataError } from 'Components/Common/DataError';

const ContentWrapper = styled(CardSubsectionContent)`
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
`;

const VerticalCenterer = styled.div`
  display: flex;
  align-items: center;
`;
const BoxSummaryTitle = styled(VerticalCenterer)`
  justify-content: space-between;
  font-weight: 500;
  margin-bottom: 8px;
  font-size: 20px;
  line-height: 24px;

  svg {
    max-height: 30px;
    margin-left: 8px;
  }
`;

const SpacingArea = styled.div`
  margin-bottom: 28px;
`;
const CountingArea = styled.div`
  margin-bottom: 15px;

  > div {
    margin-bottom: 10px;
  }
`;

const SuccessContainer = styled.div`
  display: flex;

  svg {
    width: 38px;
    min-width: 38px;
    margin-right: 14px;
  }
`;
const StatusTitle = styled.div`
  font-size: 18px;
  font-weight: 500;
  margin-bottom: 5px;
`;
const LowerCase = styled.span`
  text-transform: lowercase;
`;

type BoxProps = {
  title: string;
  logo: React.ReactNode;
  status?: PlacementStatus.StateMap[keyof PlacementStatus.StateMap];
  description: string;
  issuesCount?: number;
  count?: number;
  link: SimpleLinkProps;
  loading?: boolean;
  loadError?: ServiceError;
};

const AdminBoxSummary = ({
  title,
  logo,
  status,
  description,
  issuesCount,
  count,
  link,
  loading,
  loadError,
}: BoxProps) => {
  return (
    <CardSubsectionWrapper>
      <ContentWrapper>
        <div>
          <BoxSummaryTitle>
            <VerticalCenterer>
              {title} {logo}
            </VerticalCenterer>
            {status !== undefined && (
              <HealthIndicator
                healthStatus={status}
                statusType={StatusType.PLACEMENT}
              />
            )}
          </BoxSummaryTitle>

          <SpacingArea>{description}</SpacingArea>
        </div>
        <div>
          {!!loadError ? (
            <DataError error={loadError} />
          ) : (
            <CountingArea>
              {!!loading ? (
                <Loading message={`Retrieving ${title} information...`} />
              ) : (
                <>
                  {!!issuesCount ? (
                    <CountBox
                      count={issuesCount}
                      message={`${title}${issuesCount === 1 ? '' : 's'} need${
                        issuesCount !== 1 ? '' : 's'
                      } your attention`}
                      healthy={false}
                    />
                  ) : (
                    issuesCount !== undefined && (
                      <SuccessContainer>
                        <SuccessCircle />
                        <div>
                          <StatusTitle>Congratulations!</StatusTitle>
                          <div>
                            All of your <LowerCase>{title}</LowerCase> are
                            configured without any issues.
                          </div>
                        </div>
                      </SuccessContainer>
                    )
                  )}
                  {count !== undefined && (
                    <CountBox
                      count={count}
                      message={`${title}s are configured within Gloo Fed`}
                      healthy={true}
                    />
                  )}
                </>
              )}
            </CountingArea>
          )}
          <SoloLink {...link} />
        </div>
      </ContentWrapper>
    </CardSubsectionWrapper>
  );
};

type LogoHolderProps = {
  keepOriginal?: boolean;
};
const LogoHolder = styled.div`
  svg {
    height: 28px;

    ${(props: LogoHolderProps) =>
      !props.keepOriginal
        ? `* {
      fill: ${colors.seaBlue};
    }`
        : ''}
  }
`;
const LogoStrokeHolder = styled(LogoHolder)`
  svg * {
    stroke: ${colors.seaBlue};
  }
`;

export const GlooAdminGatewaysBox = ({
  gateways,
  error,
}: {
  gateways?: Gateway.AsObject[];
  error?: ServiceError;
}) => {
  const issueCount = gateways?.reduce(
    (acc, gateway) =>
      acc + (gateway.status?.state !== GatewayStatus.State.ACCEPTED ? 1 : 0),
    0
  );

  return (
    <AdminBoxSummary
      title={'Gateway Configuration'}
      logo={
        <LogoHolder>
          <GatewayIcon />
        </LogoHolder>
      }
      loadError={error}
      loading={gateways === undefined}
      status={
        !!issueCount
          ? GatewayStatus.State.REJECTED
          : GatewayStatus.State.ACCEPTED
      }
      description={
        'Gateways are used to configure the protocols and ports for Envoy. Optionally, gateways can be associated with a specific set of virtual services.'
      }
      count={gateways?.length ?? 0}
      issuesCount={issueCount}
      link={{
        displayElement: 'View Gateways',
        link: 'gateways/',
      }}
    />
  );
};

export const GlooAdminProxiesBox = ({
  spec,
  error,
}: {
  spec?: GlooInstanceSpec.AsObject;
  error?: ServiceError;
}) => {
  const issuesCount =
    (spec?.check?.proxies?.errorsList.length ?? 0) +
    (spec?.check?.proxies?.warningsList.length ?? 0);

  return (
    <AdminBoxSummary
      title={'Proxy Configuration'}
      logo={
        <LogoHolder>
          <ProxyIcon />
        </LogoHolder>
      }
      loadError={error}
      loading={spec === undefined}
      status={
        !!issuesCount
          ? UpstreamStatus.State.REJECTED
          : UpstreamStatus.State.ACCEPTED
      }
      description={
        'Gloo generates proxy configs from upstreams, virtual services, and gateways, and then transforms them directly into Envoy config. If a proxy config is rejected, it means Envoy will not receive configuration updates.'
      }
      issuesCount={issuesCount}
      count={spec?.proxiesList.length ?? 0}
      link={{
        displayElement: 'View Proxy',
        link: 'proxy/',
      }}
    />
  );
};

export const GlooAdminEnvoyConfigurationsBox = ({
  spec,
  error,
}: {
  spec?: GlooInstanceSpec.AsObject;
  error?: ServiceError;
}) => {
  const replicasCount = spec // sort of the 'all possible' count
    ? spec.proxiesList.reduce((total, proxy) => total + proxy.replicas, 0)
    : 0;
  const availableReplicasCount = spec // sort of the 'all usable' count
    ? spec.proxiesList.reduce(
        (total, proxy) => total + proxy.availableReplicas,
        0
      )
    : 0;

  return (
    <AdminBoxSummary
      title={'Envoy Configuration'}
      logo={
        <LogoHolder keepOriginal={true}>
          <EnvoyIcon />
        </LogoHolder>
      }
      loadError={error}
      loading={spec === undefined}
      status={
        spec?.proxiesList.length ?? 0 > 0
          ? availableReplicasCount < replicasCount
            ? UpstreamStatus.State.REJECTED
            : UpstreamStatus.State.ACCEPTED
          : UpstreamStatus.State.REJECTED
      }
      description={
        'This is the live config dump from Envoy. This is translated directly from the proxy config and should be updated any time the proxy configuration changes.'
      }
      issuesCount={replicasCount - availableReplicasCount}
      count={replicasCount}
      link={{
        displayElement: 'View Envoy',
        link: 'envoy/',
      }}
    />
  );
};

export const GlooAdminSettingsBox = () => {
  return (
    <AdminBoxSummary
      title={'Settings'}
      logo={
        <LogoStrokeHolder>
          <GearIcon />
        </LogoStrokeHolder>
      }
      description={
        "Settings expose configuration options for all components of a given Gloo Instance's control plane."
      }
      link={{
        displayElement: 'View Settings',
        link: 'settings/',
      }}
    />
  );
};

export const GlooAdminWatchedNamespacesBox = () => {
  return (
    <AdminBoxSummary
      title={'Watched Namespaces'}
      logo={
        <LogoHolder>
          <WatchedNamespacesIcon />
        </LogoHolder>
      }
      description={
        'This setting restricts the namespaces that Gloo Edge controllers take into consideration when watching for resources. In a typical production scenario, RBAC policies will limit the namespaces that Gloo Edge has access to.'
      }
      link={{
        displayElement: 'View Watched Namespaces',
        link: 'watched-namespaces/',
      }}
    />
  );
};

export const GlooAdminSecretsBox = () => {
  return (
    <AdminBoxSummary
      title={'Secrets'}
      logo={
        <LogoHolder>
          <SecretsIcon />
        </LogoHolder>
      }
      description={
        'Certain features such as the AWS Lambda option require the use of secrets for authentication, configuration of SSL Certificates, and other data that should not be stored in plaintext configuration. Gloo Edge runs an independent (goroutine) controller to monitor secrets. Secrets are stored in their own secret storage layer.'
      }
      link={{
        displayElement: 'View Secrets',
        link: 'secrets/',
      }}
    />
  );
};
