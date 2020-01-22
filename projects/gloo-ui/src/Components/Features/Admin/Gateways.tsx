import styled from '@emotion/styled';
import { ReactComponent as GatewayLogo } from 'assets/gateway-icon.svg';
import { FileDownloadLink } from 'Components/Common/FileDownloadLink';
import { SectionCard } from 'Components/Common/SectionCard';
import _ from 'lodash/fp';
import * as React from 'react';
import { useDispatch } from 'react-redux';
import { updateGateway } from 'store/gateway/actions';
import { colors, healthConstants } from 'Styles';
import useSWR from 'swr';
import { gatewayAPI } from 'store/gateway/api';
import { GatewayForm, HttpConnectionManagerSettingsForm } from './GatewayForm';

const InsideHeader = styled.div`
  display: flex;
  justify-content: space-between;
  font-size: 18px;
  line-height: 22px;
  margin-bottom: 5px;
  color: ${colors.novemberGrey};
`;

const GatewayLogoFullSize = styled(GatewayLogo)`
  width: 33px !important;
  max-height: none !important;
`;

const Link = styled.div`
  cursor: pointer;
  color: ${colors.seaBlue};
  font-size: 14px;
  padding-right: 10px;
`;

export const Gateways = () => {
  const dispatch = useDispatch();
  const [gatewaysOpen, setGatewaysOpen] = React.useState<boolean[]>([]);

  const { data: gatewaysList, error } = useSWR(
    'listGateways',
    gatewayAPI.listGateways
  );

  React.useEffect(() => {
    if (!!gatewaysList && !!gatewaysList.length) {
      setGatewaysOpen(Array.from({ length: gatewaysList.length }, () => false));
    }
  }, [gatewaysList?.length]);

  if (!gatewaysList) {
    return <div>Loading...</div>;
  }

  const toggleExpansion = (indexToggled: number) => {
    setGatewaysOpen(
      gatewaysOpen.map((isOpen, ind) => {
        if (ind !== indexToggled) {
          return false;
        }
        return !isOpen;
      })
    );
  };
  return (
    <>
      {gatewaysList.map((gateway, ind) => {
        return (
          <SectionCard
            key={gateway.gateway!.proxyNamesList[0] + ind}
            cardName={gateway.gateway!.proxyNamesList[0]}
            logoIcon={<GatewayLogoFullSize />}
            headerSecondaryInformation={[
              {
                title: 'BindPort',
                value: gateway.gateway!.bindPort.toString()
              },
              {
                title: 'Namespace',
                value: gateway.gateway!.metadata!.namespace
              },
              { title: 'SSL', value: gateway.gateway!.ssl ? 'True' : 'False' }
            ]}
            health={
              gateway.gateway!.status
                ? gateway.gateway!.status!.state
                : healthConstants.Pending.value
            }
            healthMessage={'Gateway Status'}>
            <InsideHeader>
              <div>Configuration Settings</div>{' '}
              <div style={{ display: 'flex' }}>
                <Link onClick={() => toggleExpansion(ind)}>
                  {gatewaysOpen[ind] ? 'Hide' : 'View'} Settings
                </Link>
                {!!gateway.raw && (
                  <FileDownloadLink
                    fileName={gateway.raw.fileName}
                    fileContent={gateway.raw.content}
                  />
                )}
              </div>
            </InsideHeader>
            <GatewayForm
              doUpdate={(values: HttpConnectionManagerSettingsForm) => {
                let newGateway = _.set(
                  'gateway.httpGateway.plugins',
                  {
                    httpConnectionManagerSettings: values
                  },
                  gateway
                );
                dispatch(updateGateway({ gateway: newGateway.gateway! }));
              }}
              gatewayValues={gateway.gateway!}
              gatewayConfiguration={gateway.raw}
              isExpanded={gatewaysOpen[ind]}
            />
          </SectionCard>
        );
      })}
    </>
  );
};
