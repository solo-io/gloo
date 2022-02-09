import React from 'react';
import styled from '@emotion/styled';
import { colors } from 'Styles/colors';
import {
  useListUpstreams,
  useListVirtualServices,
  useListClusterDetails,
} from 'API/hooks';
import { ReactComponent as VirtualServiceIcon } from 'assets/virtualservice-icon.svg';
import { ReactComponent as UpstreamIcon } from 'assets/upstream-icon.svg';
import { ReactComponent as ClusterIcon } from 'assets/cluster-icon.svg';
import { ReactComponent as SuccessCircle } from 'assets/big-successful-checkmark.svg';
import { ReactComponent as WarningExclamation } from 'assets/big-warning-exclamation.svg';
import { UpstreamStatus } from 'proto/github.com/solo-io/solo-apis/api/gloo/gloo/v1/upstream_pb';
import { Loading } from 'Components/Common/Loading';
import { HealthIndicator } from 'Components/Common/HealthIndicator';
import { VirtualServiceStatus } from 'proto/github.com/solo-io/solo-apis/api/gloo/gateway/v1/virtual_service_pb';
import { CountBox } from 'Components/Common/CountBox';
import { SoloLink } from 'Components/Common/SoloLink';
import { ServiceError } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/gateway_resources_pb_service';
import { DataError } from 'Components/Common/DataError';
import { ReactComponent as GraphQLIcon } from 'assets/graphql-icon.svg';
type BoxProps = {
  title: string;
  logo: React.ReactNode;
  description: string;
  descriptionTitle?: string;
  status: UpstreamStatus.StateMap[keyof UpstreamStatus.StateMap];
  count: number;
  countDescription: string;
  link: string;
};

export const WarningCircle = styled.div`
  display: inline-flex;
  justify-content: center;
  align-items: center;
  width: 38px;
  height: 38px;
  border-radius: 100%;
  background: ${colors.flashlightGold};
  border: 2px solid ${colors.sunGold};
  margin-right: 10px;

  svg {
    height: 20px !important;
    width: 3px;
    margin-right: 0 !important;
  }
`;

type StatusHealthProps = {
  isWarning: boolean;
};

export const StatusHealth = styled.div<StatusHealthProps>`
  display: flex;
  margin-bottom: 20px;
  ${(props: StatusHealthProps) =>
    props.isWarning ? `color: ${colors.sunGold};` : ''}

  svg {
    margin-right: 10px;
    height: 38px;
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

/**
 * LARGE BOX
 */

const LargeBoxContainer = styled.div`
  display: flex;
`;
const LargeBoxIconHolder = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  width: 135px;
  height: 135px;
  border-radius: 100%;
  border: 1px solid ${colors.marchGrey};
  margin-right: 30px;

  svg {
    height: 95px;
  }
`;
const LargeBoxContent = styled.div`
  position: relative;
  display: flex;
  background: ${colors.januaryGrey};
  border-radius: 8px;
  padding: 18px;
  flex: 1;
  width: 100%;

  &:before {
    content: '';

    position: absolute;
    left: -20px;
    top: 47px;
    width: 0;
    height: 0;
    border-top: 20px solid transparent;
    border-bottom: 20px solid transparent;
    border-right: 20px solid ${colors.januaryGrey};
  }
`;
const LargeBoxTitle = styled.div`
  display: flex;
  align-items: center;
  font-weight: 500;
  margin-bottom: 8px;
  font-size: 20px;
  line-height: 24px;

  > div {
    margin-left: 8px;
  }
`;
const LargeBoxDescription = styled.div`
  display: flex;
  flex-direction: column;
  justify-content: space-between;
`;

const LargeStatusBox = styled.div`
  width: 550px;
  margin-left: 50px;
`;

export const OverviewGlooInstancesBox = ({
  title,
  logo,
  description,
  status,
  count,
  countDescription,
  link,
  descriptionTitle,
}: BoxProps) => {
  return (
    <LargeBoxContainer>
      <LargeBoxIconHolder>{logo}</LargeBoxIconHolder>
      <LargeBoxContent>
        <LargeBoxDescription>
          <LargeBoxTitle>
            {descriptionTitle} <HealthIndicator healthStatus={status} />
          </LargeBoxTitle>
          {description}
          <SoloLink link={link} displayElement={`View ${title}`} />
        </LargeBoxDescription>
        <LargeStatusBox>
          {count > 0 && (
            <StatusHealth isWarning={status !== UpstreamStatus.State.ACCEPTED}>
              <div>
                {status === UpstreamStatus.State.ACCEPTED ? (
                  <SuccessCircle />
                ) : (
                  <WarningCircle>
                    <WarningExclamation />
                  </WarningCircle>
                )}
              </div>
              <div>
                {status === UpstreamStatus.State.ACCEPTED ? (
                  <>
                    <StatusTitle>Congratulations!</StatusTitle>
                    <div>
                      All of your <LowerCase>{title}</LowerCase> are configured
                      without any issues.
                    </div>
                  </>
                ) : (
                  <>
                    <StatusTitle>Warning</StatusTitle>
                    <div>
                      There are issues to address with your{' '}
                      <LowerCase>{title}</LowerCase>.
                    </div>
                  </>
                )}
              </div>
            </StatusHealth>
          )}
          <CountBox count={count} message={countDescription} healthy={true} />
        </LargeStatusBox>
      </LargeBoxContent>
    </LargeBoxContainer>
  );
};

/**
 * SMALL BOXES
 */
const SmallBoxContainer = styled.div`
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  height: 100%;
`;
const SmallBoxInfo = styled.div`
  flex: 1;
`;
const SmallBoxTitle = styled.div`
  display: flex;
  align-items: center;
  font-weight: 500;
  margin-bottom: 8px;
  font-size: 20px;
  line-height: 24px;

  svg {
    width: auto;
    height: 26px;
    max-height: 26px;
    margin-left: 8px;
  }
`;
const SmallBoxDescription = styled.div`
  font-size: 16px;
  line-height: 19px;
  padding-bottom: 18px;
  border-bottom: 1px solid ${colors.marchGrey};
  margin-bottom: 18px;
`;
const SmallStatusBox = styled.div`
  margin-bottom: 22px;
`;

export const OverviewSmallBoxSummary = ({
  title,
  logo,
  description,
  status,
  count,
  countDescription,
  link,
}: BoxProps) => {
  return (
    <SmallBoxContainer>
      <SmallBoxInfo>
        <SmallBoxTitle>
          {title} {logo}
        </SmallBoxTitle>
        <SmallBoxDescription>{description}</SmallBoxDescription>
        <SmallStatusBox>
          {count > 0 && (
            <StatusHealth isWarning={status !== UpstreamStatus.State.ACCEPTED}>
              <div>
                {status === UpstreamStatus.State.ACCEPTED ? (
                  <SuccessCircle />
                ) : (
                  <WarningCircle>
                    <WarningExclamation />
                  </WarningCircle>
                )}
              </div>
              <div>
                {status === UpstreamStatus.State.ACCEPTED ? (
                  <>
                    <StatusTitle>Congratulations!</StatusTitle>
                    <div>
                      All of your <LowerCase>{title}</LowerCase> are configured
                      without any issues.
                    </div>
                  </>
                ) : (
                  <>
                    <StatusTitle>Warning</StatusTitle>
                    <div>
                      There are issues to address with your{' '}
                      <LowerCase>{title}</LowerCase>.
                    </div>
                  </>
                )}
              </div>
            </StatusHealth>
          )}
          <CountBox count={count} message={countDescription} healthy={true} />
        </SmallStatusBox>
      </SmallBoxInfo>

      <SoloLink link={link} displayElement={`View ${title}`} />
    </SmallBoxContainer>
  );
};

export const OverviewVirtualServiceBox = () => {
  const { data: virtualServices, error: vsError } = useListVirtualServices();

  if (!!vsError) {
    return <DataError error={vsError} />;
  } else if (!virtualServices) {
    return <Loading message={'Retrieving virtual services...'} />;
  }

  const servicesStatus = virtualServices.some(
    vs => vs.status?.state !== VirtualServiceStatus.State.ACCEPTED
  )
    ? VirtualServiceStatus.State.WARNING
    : VirtualServiceStatus.State.ACCEPTED;

  return (
    <OverviewSmallBoxSummary
      title={'Virtual Services'}
      logo={<VirtualServiceIcon />}
      description='Virtual Services define a set of route rules for a given domain or set of domains.'
      status={servicesStatus}
      count={virtualServices?.length ?? 0}
      countDescription={
        'Virtual Services currently running across all of your Gloo instances'
      }
      link='/virtual-services/'
    />
  );
};
export const OverviewUpstreamsBox = () => {
  const { data: upstreams, error: uError } = useListUpstreams();

  if (!!uError) {
    return <DataError error={uError} />;
  } else if (!upstreams) {
    return <Loading message={'Retrieving upstreams...'} />;
  }

  const servicesStatus = upstreams.some(
    upstream => upstream.status?.state !== VirtualServiceStatus.State.ACCEPTED
  )
    ? UpstreamStatus.State.WARNING
    : UpstreamStatus.State.ACCEPTED;

  return (
    <OverviewSmallBoxSummary
      title={'Upstreams'}
      logo={<UpstreamIcon />}
      description='Upstreams define destinations for routes. Upstreams tell Gloo what to route to and how to route to them.'
      status={servicesStatus}
      count={upstreams?.length ?? 0}
      countDescription={
        'Upstreams currently running across all of your Gloo instances'
      }
      link='/upstreams/'
    />
  );
};

export const OverviewClustersBox = () => {
  const { data: clusterDetailsList, error: cError } = useListClusterDetails();

  if (!!cError) {
    return <DataError error={cError} />;
  } else if (!clusterDetailsList) {
    return <Loading message={'Retrieving clusters...'} />;
  }

  const servicesStatus = UpstreamStatus.State.ACCEPTED;

  return (
    <OverviewSmallBoxSummary
      title={'Clusters'}
      logo={<ClusterIcon />}
      description='Gloo Edge Federation identifies and manages Gloo Edge Instances on all registered clusters.'
      status={servicesStatus}
      count={clusterDetailsList?.length ?? 0}
      countDescription={'clusters are being managed within your environment'}
      link='/admin/clusters/'
    />
  );
};
export const OverviewGraphQLBox = () => {
  const { data: upstreams, error: uError } = useListUpstreams();

  if (!!uError) {
    return <DataError error={uError} />;
  } else if (!upstreams) {
    return <Loading message={'Retrieving upstreams...'} />;
  }

  const servicesStatus = upstreams.some(
    upstream => upstream.status?.state !== VirtualServiceStatus.State.ACCEPTED
  )
    ? UpstreamStatus.State.WARNING
    : UpstreamStatus.State.ACCEPTED;

  return (
    <OverviewSmallBoxSummary
      title={'APIs'}
      logo={<GraphQLIcon />}
      description='Generate graphql schema from other sources (e.g. openapi schema, grpc protos, etc.) and make those API available behind graphql'
      status={servicesStatus}
      count={1}
      countDescription={
        'APIs currently running across all of your Gloo instances'
      }
      link='/apis/'
    />
  );
};
