import React from 'react';
import styled from '@emotion/styled';
import { colors } from 'Styles/colors';
import {
  Card,
  CardSubsectionContent,
  CardSubsectionWrapper,
} from 'Components/Common/Card';
import { useListGlooInstances, useListGateways } from 'API/hooks';
import { objectMetasAreEqual } from 'API/helpers';
import { SoloLink } from 'Components/Common/SoloLink';
import { HealthIndicator } from 'Components/Common/HealthIndicator';
import { getGlooInstanceStatus } from 'utils/gloo-instance-helpers';
import { ReactComponent as GlooIcon } from 'assets/Gloo.svg';
import { ReactComponent as EnvoyLogo } from 'assets/envoy-logo.svg';
import { ReactComponent as GatewayIcon } from 'assets/gateway-small-icon.svg';
import { ReactComponent as ProxyLockIcon } from 'assets/lock-icon.svg';
import { ReactComponent as GearIcon } from 'assets/gear-icon.svg';
import { DataError } from 'Components/Common/DataError';
import { Loading } from 'Components/Common/Loading';

const Title = styled.div`
  font-size: 20px;
  line-height: 24px;
  font-weight: 500;
`;

const IconHolder = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  width: 33px;
  min-width: 33px;
  max-width: 33px;
  height: 33px;
  border-radius: 33px;
  margin-right: 10px;
`;

const InstanceRow = styled(CardSubsectionContent)`
  display: flex;
  align-items: center;
  margin-top: 15px;
  border: 1px solid ${colors.marchGrey};
`;

const InstanceTitle = styled.div`
  display: flex;
  align-items: center;
  width: 250px;
  padding-right: 25px;
  border-right: 1px solid ${colors.marchGrey};
  margin-right: 40px;
  cursor: default;
`;
const InstanceTitleIconHolder = styled(IconHolder)`
  position: relative;
  background: ${colors.februaryGrey};

  svg {
    width: 24px;
    height: 24px;
  }
`;
const InstanceStatus = styled.div`
  position: absolute;
  top: 0;
  right: 0;
`;
const InstanceName = styled.div`
  flex: 1;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
`;

const ConfigBlock = styled.div`
  display: flex;
  align-items: center;
  margin-right: 35px;

  font-size: 14px;
  line-height: 17px;
`;

const GatewayIconHolder = styled(IconHolder)`
  justify-content: start;
  background: ${colors.neptuneBlue};

  svg {
    width: 30px;

    * {
      fill: white;
    }
  }
`;
const ProxyIconHolder = styled(IconHolder)`
  background: ${colors.seaBlue};

  svg {
    width: 16px;
    * {
      fill: white;
    }
  }
`;
const EnvoyIconHolder = styled(IconHolder)`
  background: ${colors.envoyPink};

  svg {
    height: 16px;
    * {
      fill: white;
    }
  }
`;

const ConfigCount = styled(Title)`
  margin-right: 7px;
`;

const AdminSettingsLinkHolder = styled.div`
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: flex-end;
`;

const AdminSettingsLink = styled.div`
  display: flex;
  align-items: center;

  svg {
    height: 24px;
    width: 24px;

    margin-left: 6px;

    * {
      stroke: ${colors.seaBlue};
      fill: ${colors.seaBlue};
    }
  }
`;

export const AdminGlooInstancesTable = () => {
  const { data: glooInstances, error: instancesError } = useListGlooInstances();
  const { data: gateways, error: gatewaysError } = useListGateways();

  if (!!instancesError) {
    return <DataError error={instancesError} />;
  } else if (!!gatewaysError) {
    return <DataError error={gatewaysError} />;
  } else if (!glooInstances) {
    return <Loading message={`Retrieving gloo instances...`} />;
  } else if (!gateways) {
    return <Loading message={`Retrieving gateway information...`} />;
  }

  return (
    <div>
      <Title>Gloo Instances</Title>
      {glooInstances.map(glooInstance => {
        const gatewaysForInstance = gateways?.filter(gateway =>
          objectMetasAreEqual(
            glooInstance.metadata
              ? {
                  name: glooInstance.metadata?.name,
                  namespace: glooInstance.metadata.namespace,
                }
              : undefined,
            gateway.glooInstance
          )
        );

        return (
          <InstanceRow
            key={
              glooInstance.metadata?.name ??
              '' + glooInstance.metadata?.namespace + glooInstance.spec?.cluster
            }>
            <InstanceTitle title={glooInstance.metadata?.name}>
              <InstanceTitleIconHolder>
                <GlooIcon />
                <InstanceStatus>
                  <HealthIndicator
                    healthStatus={getGlooInstanceStatus(glooInstance)}
                    small={true}
                  />
                </InstanceStatus>
              </InstanceTitleIconHolder>
              <InstanceName>{glooInstance.metadata?.name}</InstanceName>
            </InstanceTitle>
            <ConfigBlock>
              <GatewayIconHolder>
                <GatewayIcon />
              </GatewayIconHolder>
              <ConfigCount>{gatewaysForInstance?.length}</ConfigCount> Gateway
              Configuration{gatewaysForInstance?.length !== 1 && 's'}
            </ConfigBlock>
            <ConfigBlock>
              <ProxyIconHolder>
                <ProxyLockIcon />
              </ProxyIconHolder>
              <ConfigCount>{glooInstance.spec?.proxiesList.length}</ConfigCount>
              Proxy Configuration
              {glooInstance.spec?.proxiesList.length !== 1 && 's'}
            </ConfigBlock>
            <ConfigBlock>
              <EnvoyIconHolder>
                <EnvoyLogo />
              </EnvoyIconHolder>
              <ConfigCount>
                {glooInstance.spec
                  ? glooInstance.spec.proxiesList.reduce(
                      (total, proxy) => total + proxy.readyReplicas,
                      0
                    )
                  : 0}
              </ConfigCount>{' '}
              Envoy Configuration
              {glooInstances.length && 's'}
            </ConfigBlock>
            <AdminSettingsLinkHolder>
              <SoloLink
                displayElement={
                  <AdminSettingsLink>
                    Admin Settings <GearIcon />
                  </AdminSettingsLink>
                }
                link={`/gloo-instances/${glooInstance.metadata?.namespace}/${glooInstance.metadata?.name}/gloo-admin/`}
              />
            </AdminSettingsLinkHolder>
          </InstanceRow>
        );
      })}
    </div>
  );
};
