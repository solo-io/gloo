import React, { useState, useEffect } from 'react';
import styled from '@emotion/styled';
import { colors } from 'Styles/colors';
import { Card } from 'Components/Common/Card';
import { ReactComponent as HealthIcon } from 'assets/health-icon.svg';
import { useParams } from 'react-router';
import {
  GlooAdminGatewaysBox,
  GlooAdminProxiesBox,
  GlooAdminSettingsBox,
  GlooAdminEnvoyConfigurationsBox,
  GlooAdminWatchedNamespacesBox,
  GlooAdminSecretsBox,
} from './GlooInstanceAdminBoxSummary';
import { GlooInstance } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/glooinstance_pb';
import { Loading } from 'Components/Common/Loading';
import {
  useListGlooInstances,
  useListGateways,
  useListSettings,
} from 'API/hooks';
import { DataError } from 'Components/Common/DataError';

const Heading = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
`;

const HeadingTitle = styled.div`
  font-size: 22px;
  line-height: 26px;
`;
const HeadingSubtitle = styled.div`
  font-size: 18px;
  line-height: 22px;
  margin-top: 10px;
`;
const HeadingLogo = styled.div``;

const CardSection = styled.div`
  margin-top: 20px;
`;
const Row = styled.div`
  display: grid;
  grid-gap: 22px;
`;
const TopRow = styled(Row)`
  grid-template-columns: 1fr 1fr 1fr;
  margin-bottom: 22px;
`;
const BottomRow = styled(Row)`
  margin-top: 22px;
  grid-template-columns: 1fr 1fr;
`;

export const GlooInstanceAdministration = () => {
  const { name = '', namespace = '' } = useParams();

  const { data: glooInstances, error: instancesError } = useListGlooInstances();
  const { data: gateways, error: gatewaysError } = useListGateways({
    name,
    namespace,
  });

  const [glooInstance, setGlooInstance] = useState<GlooInstance.AsObject>();

  useEffect(() => {
    if (!!glooInstances) {
      setGlooInstance(
        glooInstances.find(
          instance =>
            instance.metadata?.name === name &&
            instance.metadata?.namespace === namespace
        )
      );
    } else {
      setGlooInstance(undefined);
    }
  }, [name, namespace, glooInstances]);

  if (!!instancesError) {
    return <DataError error={instancesError} />;
  } else if (!!gatewaysError) {
    return <DataError error={gatewaysError} />;
  } else if (!glooInstances) {
    return <Loading message={'Retrieving instances...'} />;
  } else if (!gateways) {
    return <Loading message={`Retrieving gateways for ${name}...`} />;
  }

  return (
    <Card>
      <Heading>
        <div>
          <HeadingTitle>{name} Administration</HeadingTitle>
          <HeadingSubtitle>
            Advanced Administration for your Gloo Edge Configuration
          </HeadingSubtitle>
        </div>
        <HeadingLogo>
          <HealthIcon />
        </HeadingLogo>
      </Heading>

      <CardSection>
        {!glooInstance ? (
          <Loading
            message={`Retrieving administration information for ${name} instance`}
          />
        ) : (
          <>
            <TopRow>
              <GlooAdminGatewaysBox gateways={gateways} error={gatewaysError} />
              <GlooAdminProxiesBox
                glooInstance={glooInstance}
                error={instancesError}
              />
              <GlooAdminEnvoyConfigurationsBox
                glooInstance={glooInstance}
                error={instancesError}
              />
            </TopRow>

            <Heading>
              <HeadingTitle>Security Settings</HeadingTitle>
            </Heading>

            <BottomRow>
              <GlooAdminSettingsBox />
              <GlooAdminWatchedNamespacesBox />
            </BottomRow>
          </>
        )}
      </CardSection>
    </Card>
  );
};
