import React from 'react';
import styled from '@emotion/styled';
import { SoloLink, SimpleLinkProps } from 'Components/Common/SoloLink';
import { CountBox } from 'Components/Common/CountBox';
import {
  CardSubsectionWrapper,
  CardSubsectionContent,
} from 'Components/Common/Card';
import { HealthIndicator } from 'Components/Common/HealthIndicator';
import { colors } from 'Styles/colors';
import { ReactComponent as CloudIcon } from 'assets/auth-cloud-icon.svg';
import { ReactComponent as VirtualServiceIcon } from 'assets/virtualservice-icon.svg';
import { ReactComponent as UpstreamsIcon } from 'assets/upstreams-icon.svg';
import { ReactComponent as UpstreamGroupIcon } from 'assets/upstream-group-icon.svg';
import { ReactComponent as FedResourcesIcon } from 'assets/GlooFed-Specific/federated-resources-icon-gray.svg';
import { ReactComponent as RoutesIcon } from 'assets/route-icon.svg';
import { ReactComponent as GatewayIcon } from 'assets/gateway.svg';
import { ReactComponent as KeyIcon } from 'assets/key-icon.svg';
import { ReactComponent as ClustersIcon } from 'assets/cluster-icon.svg';
import { ReactComponent as SuccessCircle } from 'assets/big-successful-checkmark.svg';
import { ReactComponent as WarningExclamation } from 'assets/big-warning-exclamation.svg';
import {
  useListFederatedVirtualServices,
  useListFederatedRouteTables,
  useListFederatedUpstreams,
  useListFederatedUpstreamGroups,
  useListFederatedAuthConfigs,
  useListFederatedGateways,
  useListFederatedSettings,
  useListClusterDetails,
} from 'API/hooks';
import { loadavg } from 'os';
import { Loading } from 'Components/Common/Loading';
import { VirtualServiceStatus } from 'proto/github.com/solo-io/solo-apis/api/gloo/gateway/v1/virtual_service_pb';
import { PlacementStatus } from 'proto/github.com/solo-io/solo-projects/projects/gloo-fed/api/fed/core/v1/placement_pb';
import { IconHolder } from 'Styles/StyledComponents/icons';
import { StatusType } from 'utils/health-status';
import { DataError } from 'Components/Common/DataError';
import { ServiceError } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/gateway_resources_pb_service';

const VerticalCenterer = styled.div`
  display: flex;
  align-items: space-between;
`;
const BoxSummaryTitle = styled(VerticalCenterer)`
  justify-content: space-between;
  font-weight: 500;
  margin-bottom: 13px;
  font-size: 20px;
  line-height: 24px;

  svg {
    max-height: 30px;
    margin-left: 8px;
  }
`;

const BoxDescription = styled.div`
  margin-bottom: 30px;
`;

const StatusHealth = styled.div`
  display: flex;
  margin-bottom: 20px;
  min-width: 250px;

  svg {
    margin-right: 10px;
    height: 38px;
    width: 38px;
  }
`;
const StatusTitle = styled.div`
  font-size: 18px;
  font-weight: 500;
  margin-bottom: 5px;
`;

const SmallContentSubsection = styled(CardSubsectionContent)`
  height: 100%;
`;

type BoxProps = {
  title: string;
  logo: React.ReactNode;
  status?: PlacementStatus.StateMap[keyof PlacementStatus.StateMap];
  issuesCount?: number;
  count: number;
  link: SimpleLinkProps;
  loadError?: ServiceError;
  loading: boolean;
  description: string;
  congratulationsText: string;
};

const AdminBoxSummary = ({
  title,
  logo,
  status,
  issuesCount,
  count,
  link,
  loadError,
  loading,
  description,
  congratulationsText,
}: BoxProps) => {
  return (
    <SmallContentSubsection>
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

      <BoxDescription>{description}</BoxDescription>

      {!!loadError ? (
        <DataError error={loadError} />
      ) : loading ? (
        <Loading message={`Retrieving federated ${title} information...`} />
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
            <StatusHealth>
              <SuccessCircle />
              <div>
                <StatusTitle>Congratulations!</StatusTitle>
                <div>{congratulationsText}</div>
              </div>
            </StatusHealth>
          )}
          <CountBox
            count={count}
            message={`${title}s are configured within Gloo Edge`}
            healthy={true}
          />

          <div
            style={{
              marginTop: '20px',
            }}>
            <SoloLink {...link} />
          </div>
        </>
      )}
    </SmallContentSubsection>
  );
};

const LogoHolder = styled.div`
  svg {
    height: 28px;
  }
`;
const LogoRecolorHolder = styled(LogoHolder)`
  svg {
    * {
      fill: ${colors.seaBlue};
    }
  }
`;

export const AdminFederatedResourcesBox = () => {
  const { data: fedVirtualServices, error: fedVsError } =
    useListFederatedVirtualServices();
  const { data: fedRouteTables, error: fedRtError } =
    useListFederatedRouteTables();
  const { data: fedUpstreams, error: fedUError } = useListFederatedUpstreams();
  const { data: fedUpstreamGroups, error: fedUGError } =
    useListFederatedUpstreamGroups();
  const { data: fedAuthConfigs, error: fedACError } =
    useListFederatedAuthConfigs();
  const { data: fedGateways, error: fedGError } = useListFederatedGateways();
  const { data: fedSettings, error: fedSError } = useListFederatedSettings();

  if (fedVsError) {
    return <DataError error={fedVsError} />;
  }
  if (fedRtError) {
    return <DataError error={fedRtError} />;
  }
  if (fedUError) {
    return <DataError error={fedUError} />;
  }
  if (fedUGError) {
    return <DataError error={fedUGError} />;
  }
  if (fedACError) {
    return <DataError error={fedACError} />;
  }
  if (fedGError) {
    return <DataError error={fedGError} />;
  }
  if (fedSError) {
    return <DataError error={fedSError} />;
  }

  const issueCount =
    (fedSettings?.reduce(
      (acc, setting) =>
        acc +
        (setting.status?.placementStatus?.state !== PlacementStatus.State.PLACED
          ? 1
          : 0),
      0
    ) ?? 0) +
    (fedVirtualServices?.reduce(
      (acc, vs) =>
        acc +
        (vs.status?.placementStatus?.state !== PlacementStatus.State.PLACED
          ? 1
          : 0),
      0
    ) ?? 0) +
    (fedRouteTables?.reduce(
      (acc, rt) =>
        acc +
        (rt.status?.placementStatus?.state !== PlacementStatus.State.PLACED
          ? 1
          : 0),
      0
    ) ?? 0) +
    (fedUpstreams?.reduce(
      (acc, upstream) =>
        acc +
        (upstream.status?.placementStatus?.state !==
        PlacementStatus.State.PLACED
          ? 1
          : 0),
      0
    ) ?? 0) +
    (fedUpstreamGroups?.reduce(
      (acc, ug) =>
        acc +
        (ug.status?.placementStatus?.state !== PlacementStatus.State.PLACED
          ? 1
          : 0),
      0
    ) ?? 0) +
    (fedAuthConfigs?.reduce(
      (acc, auth) =>
        acc +
        (auth.status?.placementStatus?.state !== PlacementStatus.State.PLACED
          ? 1
          : 0),
      0
    ) ?? 0) +
    (fedGateways?.reduce(
      (acc, gateway) =>
        acc +
        (gateway.status?.placementStatus?.state !== PlacementStatus.State.PLACED
          ? 1
          : 0),
      0
    ) ?? 0);

  return (
    <AdminBoxSummary
      title={'Federated Resources'}
      logo={
        <LogoRecolorHolder>
          <FedResourcesIcon />
        </LogoRecolorHolder>
      }
      loadError={
        fedSError ||
        fedVsError ||
        fedUGError ||
        fedRtError ||
        fedACError ||
        fedUError
      }
      loading={fedSettings === undefined}
      status={
        !!issueCount
          ? PlacementStatus.State.FAILED
          : PlacementStatus.State.PLACED
      }
      count={
        (fedSettings?.length ?? 0) +
        (fedVirtualServices?.length ?? 0) +
        (fedUpstreamGroups?.length ?? 0) +
        (fedUpstreams?.length ?? 0) +
        (fedAuthConfigs?.length ?? 0) +
        (fedRouteTables?.length ?? 0)
      }
      issuesCount={issueCount}
      link={{
        displayElement: 'View Federated Resources',
        link: 'federated-resources/',
      }}
      description={
        'Federated resources allow you to define your configuration in one place and have it distributed to one or many Gloo Instances.'
      }
      congratulationsText={
        'All of your resources are configured without any issues.'
      }
    />
  );
};

export const AdminClustersBox = () => {
  const { data: clusterDets, error: cError } = useListClusterDetails();

  if (!!cError) {
    return <DataError error={cError} />;
  } else if (!clusterDets) {
    return <Loading message={`Retrieving cluster information...`} />;
  }

  return (
    <AdminBoxSummary
      title={'Cluster'}
      logo={
        <LogoHolder>
          <ClustersIcon />
        </LogoHolder>
      }
      loadError={cError}
      loading={clusterDets === undefined}
      status={
        !!cError ? PlacementStatus.State.FAILED : PlacementStatus.State.PLACED
      }
      count={clusterDets?.length ?? 0}
      issuesCount={!!cError ? 1 : 0}
      link={{
        displayElement: 'View Clusters',
        link: 'clusters/',
      }}
      description={
        'Gloo Edge Federation identifies and manages Gloo Edge Instances on all registered clusters.'
      }
      congratulationsText={
        'All of your clusters are configured without any issues.'
      }
    />
  );
};
