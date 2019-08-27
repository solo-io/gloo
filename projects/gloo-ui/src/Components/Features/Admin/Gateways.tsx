import * as React from 'react';
/** @jsx jsx */
import { jsx } from '@emotion/core';
import styled from '@emotion/styled/macro';
import { useUpdateGateway } from 'Api/v2/useGatewayClientV2';
import { ReactComponent as GatewayLogo } from 'assets/gateway-icon.svg';
import { colors, soloConstants, healthConstants } from 'Styles';
import { SectionCard } from 'Components/Common/SectionCard';
import { FileDownloadLink } from 'Components/Common/FileDownloadLink';

import {
  GatewayDetails,
  UpdateGatewayRequest
} from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/gateway_pb';

import { UpdateGatewayHttpData } from 'Api/v2/GatewayClient';
import { useSelector, useDispatch } from 'react-redux';
import { AppState } from 'store';
import { HttpConnectionManagerSettings } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/hcm/hcm_pb';
import { GatewayForm, HttpConnectionManagerSettingsForm } from './GatewayForm';
import { updateGateway } from 'store/gateway/actions';
import _ from 'lodash/fp';

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
`;

interface Props {}

export const Gateways = (props: Props) => {
  const dispatch = useDispatch();
  const [gatewaysOpen, setGatewaysOpen] = React.useState<boolean[]>([]);
  const gatewaysList = useSelector(
    (state: AppState) => state.gateways.gatewaysList
  );

  const [allGateways, setAllGateways] = React.useState<
    GatewayDetails.AsObject[]
  >([]);

  React.useEffect(() => {
    if (!!gatewaysList.length) {
      const newGateways = gatewaysList.filter(gateway => !!gateway.gateway);
      setAllGateways(newGateways);
      setGatewaysOpen(Array.from({ length: newGateways.length }, () => false));
    }
  }, [gatewaysList.length]);

  if (!gatewaysList.length) {
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
    <React.Fragment>
      {gatewaysList.map((gateway, ind) => {
        return (
          <SectionCard
            key={gateway.gateway!.gatewayProxyName + ind}
            cardName={gateway.gateway!.gatewayProxyName}
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
              {!!gateway.raw && (
                <FileDownloadLink
                  fileName={gateway.raw.fileName}
                  fileContent={gateway.raw.content}
                />
              )}
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
            <Link onClick={() => toggleExpansion(ind)}>
              {gatewaysOpen[ind] ? 'Hide' : 'View'} Settings
            </Link>
          </SectionCard>
        );
      })}
    </React.Fragment>
  );
};
