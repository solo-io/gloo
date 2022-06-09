import React from 'react';
import styled from '@emotion/styled/macro';
import { Loading } from 'Components/Common/Loading';
import {
  useListGlooInstances,
  useListVirtualServices,
  useListUpstreams,
} from 'API/hooks';
import { GlooInstance } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/glooinstance_pb';
import { SectionCard } from 'Components/Common/SectionCard';
import { ReactComponent as GlooIcon } from 'assets/Gloo.svg';
import { ReactComponent as MeshIcon } from 'assets/mesh-icon.svg';
import { ReactComponent as ClusterIcon } from 'assets/cluster-icon.svg';
import { ReactComponent as NamespaceIcon } from 'assets/namespace-icon.svg';
import { ReactComponent as VersionsIcon } from 'assets/versions-icon.svg';
import { ReactComponent as EnvoyLogo } from 'assets/envoy-logo.svg';
import { ReactComponent as VirtualServiceIcon } from 'assets/virtualservice-icon.svg';
import { ReactComponent as UpstreamsIcon } from 'assets/upstreams-icon.svg';
import { ReactComponent as GearIcon } from 'assets/gear-icon.svg';
import { ReactComponent as SuccessCircle } from 'assets/big-successful-checkmark.svg';
import { colors } from 'Styles/colors';
import { CardWhiteSubsection } from 'Components/Common/Card';
import { HealthIndicator } from 'Components/Common/HealthIndicator';
import { CountBox } from 'Components/Common/CountBox';
import { UpstreamStatus } from 'proto/github.com/solo-io/solo-apis/api/gloo/gloo/v1/upstream_pb';
import { objectMetasAreEqual } from 'API/helpers';
import { getGlooInstanceStatus } from 'utils/gloo-instance-helpers';
import { GlooInstanceIssues } from './GlooInstanceIssues';
import { SoloLink } from 'Components/Common/SoloLink';
import { DataError } from 'Components/Common/DataError';
import { ServiceError } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/gateway_resources_pb_service';

const GlooIconHolder = styled.div`
  display: flex;
  align-items: center;
  justify-items: center;

  svg {
    width: 25px;
  }
`;

const OrangeIconHolder = styled.div`
  svg * {
    fill: ${colors.grapefruitOrange};
  }
`;
const BlueIconHolder = styled.div`
  svg * {
    fill: ${colors.seaBlue};
  }
`;

const CardContent = styled.div`
  padding: 20px;
`;

const QuickStats = styled.div`
  display: flex;
  align-items: center;

  > div {
    display: flex;
    align-items: center;
    margin-right: 45px;
  }

  svg {
    height: 30px;
    margin-right: 10px;
  }
`;
const QuickStatTitle = styled.div`
  font-weight: 500;
  margin-right: 4px;
`;

const Divider = styled.div`
  width: 1px;
  height: 42px;
  background: ${colors.marchGrey};
`;

const CardFooter = styled.div`
  border-radius: 0 0 10px 10px;
  background: ${colors.februaryGrey};
  padding: 0 22px;
  line-height: 40px;
`;

const StatusOuterBox = styled(CardWhiteSubsection)`
  margin: 20px 0;
  display: flex;
`;
const StatusInnerBox = styled.div`
  position: relative;
  display: flex;
  background: ${colors.januaryGrey};
  border-radius: 8px;
  padding: 18px;
  flex: 1;
  width: 100%;
  min-height: 140px;

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

const SuccessContainer = styled.div`
  display: flex;

  svg {
    width: 38px;
    min-width: 38px;
    margin-right: 14px;
  }
`;

const StatusTitle = styled.div`
  display: flex;
  align-items: center;
  font-size: 18px;
  font-weight: 500;
  margin-bottom: 5px;

  > div {
    margin-left: 8px;
  }
`;

const EnvoyLogoHolder = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  width: 140px;
  height: 140px;
  margin-right: 20px;

  svg {
    width: 85px;
  }
`;
const EnvoyCounts = styled.div`
  display: flex;
  flex-direction: column;
  justify-content: center;
  min-width: 550px;
  margin-left: 40px;

  > div {
    margin: 10px 0;
  }
`;

const BottomRow = styled.div`
  display: grid;
  grid-template-columns: 1fr 1fr 1fr;
  grid-gap: 22px;
  margin-top: 20px;
`;

const BottomTitle = styled.div`
  display: flex;
  font-size: 20px;
  line-height: 24px;
  font-weight: 500;
  margin-bottom: 16px;

  svg {
    margin-left: 10px;
    height: 24px;

    * {
      stroke: ${colors.seaBlue};
    }
  }
`;

const HorizontalDivider = styled.div`
  height: 1px;
  width: 100%;
  background: ${colors.marchGrey};
  margin-bottom: 16px;
`;

const GlooInstanceCard = ({
  instance,
  upstreamsCount,
  upstreamsError,
  virtualServicesCount,
  virtualServicesError,
}: {
  instance: GlooInstance.AsObject;
  upstreamsCount?: number;
  upstreamsError?: ServiceError;
  virtualServicesCount?: number;
  virtualServicesError?: ServiceError;
}) => {
  const replicasCount = instance.spec // sort of the 'all possible' count
    ? instance.spec.proxiesList.reduce(
        (total, proxy) => total + proxy.replicas,
        0
      )
    : 0;
  const availableReplicasCount = instance.spec // sort of the 'all usable' count
    ? instance.spec.proxiesList.reduce(
        (total, proxy) => total + proxy.availableReplicas,
        0
      )
    : 0;

  return (
    <SectionCard
      cardName={
        instance.metadata?.name ?? 'An instance without a UID was found!'
      }
      logoIcon={
        <GlooIconHolder>
          <GlooIcon />
        </GlooIconHolder>
      }
      noPadding={true}
      health={{
        state: getGlooInstanceStatus(instance),
        title: 'Gloo Health',
        // reason: 'NO DATAA', TODO refactor gloo-instance-helpers to return most relevant message
      }}>
      <CardContent>
        <GlooInstanceIssues glooInstance={instance} />
        <CardWhiteSubsection>
          <QuickStats>
            <div>
              <OrangeIconHolder>
                <MeshIcon />
              </OrangeIconHolder>
              <QuickStatTitle>Region: </QuickStatTitle> {instance.spec?.region}
            </div>
            <div>
              <QuickStatTitle>Zones:</QuickStatTitle>{' '}
              {instance.spec?.proxiesList
                .map(prox => prox.zonesList.join(', '))
                .join(', ')}
            </div>

            <Divider />

            <div>
              <BlueIconHolder>
                <ClusterIcon />
              </BlueIconHolder>
              <QuickStatTitle>Cluster:</QuickStatTitle> {instance.spec?.cluster}
            </div>
            <div>
              <BlueIconHolder>
                <NamespaceIcon />
              </BlueIconHolder>
              <QuickStatTitle>Namespace:</QuickStatTitle>{' '}
              {instance.spec?.controlPlane?.namespace}
            </div>
            <div>
              <BlueIconHolder>
                <VersionsIcon />
              </BlueIconHolder>
              <QuickStatTitle>Version:</QuickStatTitle>{' '}
              {instance.spec?.controlPlane?.version}
            </div>
          </QuickStats>
        </CardWhiteSubsection>
        <StatusOuterBox>
          <EnvoyLogoHolder>
            <EnvoyLogo />
          </EnvoyLogoHolder>
          <StatusInnerBox>
            <div>
              <StatusTitle>
                Envoy Health Status
                {(instance.spec?.proxiesList.length ?? 0) > 0 && (
                  <HealthIndicator
                    healthStatus={
                      availableReplicasCount < replicasCount
                        ? UpstreamStatus.State.REJECTED
                        : UpstreamStatus.State.ACCEPTED
                    }
                  />
                )}
              </StatusTitle>
              <div style={{ marginBottom: '8px' }}>
                Gloo is responsible for configuring Envoy. Whenever Virtual
                Services or other configs change that affect the proxy, Gloo
                will immediately detect that change and update Envoy's
                configuration.
              </div>
              <SoloLink
                link={`${instance.metadata?.namespace}/${instance.metadata?.name}/gloo-admin/proxy/`}
                displayElement={'View Proxy Configurations'}
              />
            </div>
            <EnvoyCounts>
              {availableReplicasCount < replicasCount ? (
                <CountBox
                  count={replicasCount - availableReplicasCount}
                  message={
                    <>
                      Envoy instance
                      {replicasCount - availableReplicasCount === 1
                        ? ' is'
                        : 's are'}{' '}
                      unavailable |{' '}
                      <SoloLink
                        link={`${instance.metadata?.namespace}/${instance.metadata?.name}/`}
                        displayElement={`View Gloo Issues`}
                        inline={true}
                        stylingOverrides={`
                          color: ${colors.pumpkinOrange};
                          font-weight: 500;

                          &:hover,
                          &:focus {
                            color: ${colors.grapefruitOrange};
                          }`}
                      />
                    </>
                  }
                  healthy={false}
                />
              ) : (
                <SuccessContainer>
                  <SuccessCircle />
                  <div>
                    <StatusTitle>Congratulations!</StatusTitle>
                    <div>No configurations have issues.</div>
                  </div>
                </SuccessContainer>
              )}

              <CountBox
                count={replicasCount}
                message={`Envoy instances${
                  replicasCount === 1 ? '' : 's'
                } currently deployed and configured`}
                healthy={true}
              />
            </EnvoyCounts>
          </StatusInnerBox>
        </StatusOuterBox>
        <BottomRow>
          <CardWhiteSubsection>
            <BottomTitle>
              Virtual Services <VirtualServiceIcon />
            </BottomTitle>
            <HorizontalDivider />
            {virtualServicesCount === undefined ? (
              <Loading message={'Retrieving virtual services...'} />
            ) : (
              <CountBox
                count={virtualServicesCount}
                message={`Virtual Service${
                  /*successCount === 1*/ false ? '' : 's'
                } configured`}
                healthy={true}
              />
            )}
          </CardWhiteSubsection>
          <CardWhiteSubsection>
            <BottomTitle>
              Upstreams <UpstreamsIcon />
            </BottomTitle>
            <HorizontalDivider />
            {upstreamsCount === undefined ? (
              <Loading message={'Retrieving upstreams...'} />
            ) : (
              <CountBox
                count={upstreamsCount}
                message={`Upstream${
                  /*successCount === 1*/ false ? '' : 's'
                } configured`}
                healthy={true}
              />
            )}
          </CardWhiteSubsection>
          <CardWhiteSubsection>
            <BottomTitle>
              Admin Settings <GearIcon />
            </BottomTitle>
            <HorizontalDivider />

            <div>
              Advanced Administration for your Gloo Edge Configuration.{' '}
              <SoloLink
                inline={true}
                link={`${instance.metadata?.namespace}/${instance.metadata?.name}/gloo-admin/`}
                displayElement={'View Now.'}
              />
            </div>
          </CardWhiteSubsection>
        </BottomRow>
      </CardContent>
      <CardFooter>
        <SoloLink
          link={`${instance.metadata?.namespace}/${instance.metadata?.name}/`}
          displayElement={`View Gloo Details`}
        />
      </CardFooter>
    </SectionCard>
  );
};

export const GlooInstancesLanding = () => {
  const { data: glooInstances, error: instancesError } = useListGlooInstances();
  const { data: upstreamsResponse, error: upstreamsResponseError } =
    useListUpstreams(undefined, { limit: 0, offset: 0 });
  const { data: virtualServices, error: vsError } = useListVirtualServices(
    undefined,
    { limit: 0, offset: 0 }
  );
  const upstreamsCount = upstreamsResponse?.total ?? 0;
  const virtualServicesCount = virtualServices?.total ?? 0;

  const upstreams = upstreamsResponse?.upstreamsList;

  if (!!instancesError) {
    return <DataError error={instancesError} />;
  } else if (!!upstreamsResponseError) {
    return <DataError error={upstreamsResponseError} />;
  } else if (vsError) {
    return <DataError error={vsError} />;
  } else if (!glooInstances) {
    return (
      <Loading message={'Retrievng information on instances of Gloo...'} />
    );
  } else if (!upstreams) {
    return <Loading message={'Retrieving information on upstreams...'} />;
  } else if (!virtualServices) {
    return (
      <Loading message={'Retrieving information on virtual services...'} />
    );
  }

  return (
    <>
      {!!glooInstances.length ? (
        glooInstances.map(instance => (
          <GlooInstanceCard
            key={instance.metadata?.uid}
            instance={instance}
            virtualServicesCount={virtualServicesCount}
            virtualServicesError={vsError}
            upstreamsError={upstreamsResponseError}
            upstreamsCount={upstreamsCount}
          />
        ))
      ) : (
        <SectionCard
          cardName={'No Instances'}
          logoIcon={
            <GlooIconHolder>
              <GlooIcon />
            </GlooIconHolder>
          }
          health={{
            state: UpstreamStatus.State.PENDING,
            title: 'Unavailable',
            reason: 'No instances found',
          }}>
          No instances of Gloo were found.
        </SectionCard>
      )}
    </>
  );
};
