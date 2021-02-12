import React from 'react';
import styled from '@emotion/styled';
import {
  Card,
  CardSubsectionWrapper,
  CardSubsectionContent,
} from 'Components/Common/Card';
import { ReactComponent as HealthIcon } from 'assets/health-icon.svg';
import { ReactComponent as GlooOverviewIcon } from 'assets/GlooFed-Specific/gloo-hub-overview-icon.svg';
import { ReactComponent as EnvoyIcon } from 'assets/envoy-logo.svg';
import {
  OverviewVirtualServiceBox,
  OverviewUpstreamsBox,
  OverviewGlooInstancesBox,
  OverviewClustersBox,
} from './OverviewBoxSummary';
import { useListClusterDetails, useListGlooInstances } from 'API/hooks';
import { Loading } from 'Components/Common/Loading';
import { UpstreamStatus } from 'proto/github.com/solo-io/solo-apis/api/gloo/gloo/v1/upstream_pb';
import { GlooInstance } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/glooinstance_pb';
import {
  getGlooInstanceListStatus,
  getGlooInstanceStatus,
} from 'utils/gloo-instance-helpers';
import { DataError } from 'Components/Common/DataError';

const OverviewCard = styled(Card)`
  margin-top: 20px;
`;

const Heading = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 500;
  margin-bottom: 25px;
`;

const HeadingTitle = styled.div`
  font-size: 22px;
  line-height: 26px;
`;
const HeadingSubtitle = styled.div`
  font-size: 18px;
  line-height: 22px;
`;
const HeadingLogo = styled.div``;

type BottomRowProps = {
  threeSlots: boolean;
};
const BottomRow = styled.div<BottomRowProps>`
  display: grid;
  grid-template-columns: 1fr 1fr ${(props: BottomRowProps) =>
      props.threeSlots ? ' 1fr' : ''};
  grid-gap: 20px;
  margin-top: 20px;
`;

export const OverviewLanding = () => {
  const { data: glooInstances, error: glooError } = useListGlooInstances();
  const { data: clusterDetailsList, error: cError } = useListClusterDetails();

  if (!!cError) {
    return <DataError error={cError} />;
  } else if (!!glooError) {
    return <DataError error={glooError} />;
  } else if (!glooInstances) {
    return <Loading message={'Retrieving instances of Gloo...'} />;
  } else if (!clusterDetailsList) {
    return <Loading message={'Retrieving clusters...'} />;
  }

  return (
    <OverviewCard>
      <Heading>
        <div>
          <HeadingTitle>Enterprise Gloo Edge Overview</HeadingTitle>
          <HeadingSubtitle>
            Your current configuration health at a glance
          </HeadingSubtitle>
        </div>
        <HeadingLogo>
          <HealthIcon />
        </HeadingLogo>
      </Heading>
      <CardSubsectionWrapper>
        <CardSubsectionContent>
          {clusterDetailsList.length <= 1 && glooInstances?.length <= 1 ? (
            <OverviewGlooInstancesBox
              title={'Envoy Configurations'}
              logo={<EnvoyIcon />}
              description="Gloo is responsible for configuring Envoy. Whenever Virtual
                Services or other configs change that affect the proxy, Gloo
                will immediately detect that change and update Envoy's
                configuration."
              descriptionTitle='Envoy Health Status'
              status={getGlooInstanceStatus(
                clusterDetailsList[0]?.glooInstancesList[0],
                'proxies'
              )}
              count={
                clusterDetailsList[0]?.glooInstancesList[0]?.spec?.proxiesList
                  .length ?? 0
              }
              countDescription={
                'envoy configurations currently deployed and configured'
              }
              link={
                !!clusterDetailsList[0]?.glooInstancesList[0]?.metadata
                  ? `gloo-instances/${clusterDetailsList[0].glooInstancesList[0].metadata.namespace}/${clusterDetailsList[0].glooInstancesList[0].metadata.name}/gloo-admin/envoy/`
                  : ''
              }
            />
          ) : (
            <OverviewGlooInstancesBox
              title={'Gloo Instances'}
              logo={<GlooOverviewIcon />}
              description='Gloo Edge observes relevant configuration and status events on all registered clusters and associates them with individual Gloo instances.'
              descriptionTitle='Overall Gloo Status'
              status={getGlooInstanceListStatus(glooInstances, 'proxies')}
              count={glooInstances?.length ?? 0}
              countDescription={
                'Gloo Instances currently deployed and configured'
              }
              link='/gloo-instances/'
            />
          )}
        </CardSubsectionContent>
      </CardSubsectionWrapper>
      <BottomRow
        threeSlots={clusterDetailsList.length > 1 || glooInstances?.length > 1}>
        {(clusterDetailsList.length > 1 || glooInstances?.length > 1) && (
          <CardSubsectionWrapper>
            <CardSubsectionContent>
              <OverviewClustersBox />
            </CardSubsectionContent>
          </CardSubsectionWrapper>
        )}
        <CardSubsectionWrapper>
          <CardSubsectionContent>
            <OverviewVirtualServiceBox />
          </CardSubsectionContent>
        </CardSubsectionWrapper>
        <CardSubsectionWrapper>
          <CardSubsectionContent>
            <OverviewUpstreamsBox />
          </CardSubsectionContent>
        </CardSubsectionWrapper>
      </BottomRow>
    </OverviewCard>
  );
};
