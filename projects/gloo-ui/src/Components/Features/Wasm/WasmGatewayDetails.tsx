import { ConfigDisplayer } from 'Components/Common/DisplayOnly/ConfigDisplayer';
import { SectionCard } from 'Components/Common/SectionCard';
import { GatewayDetails } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/gateway_pb';
import React from 'react';
import { useParams } from 'react-router';
import { gatewayAPI } from 'store/gateway/api';
import useSWR from 'swr';
import { ReactComponent as GatewayLogo } from 'assets/gateway-icon.svg';
import { healthConstants } from 'Styles/healthConstants';
import { Breadcrumb } from 'Components/Common/Breadcrumb';

export type WasmGatewayDetailsProps = {};

export const WasmGatewayDetails: React.FC<WasmGatewayDetailsProps> = ({}) => {
  let { gatewayname } = useParams<{
    gatewayname: string;
  }>();
  const { data: gatewaysList, error } = useSWR(
    'listGateways',
    gatewayAPI.listGateways
  );
  const [currentGateway, setCurrentGateway] = React.useState<
    GatewayDetails.AsObject
  >();

  React.useEffect(() => {
    if (!!gatewaysList) {
      setCurrentGateway(
        gatewaysList.find(
          gateway => gateway?.gateway?.metadata?.name === gatewayname
        )
      );
    }
  }, [gatewayname, gatewaysList?.length]);

  if (!gatewaysList && !error) {
    return <div>Loading...</div>;
  }
  return (
    <>
      <Breadcrumb />
      <SectionCard
        key={currentGateway?.gateway?.metadata?.name}
        cardName={currentGateway?.gateway?.metadata?.name ?? ''}
        logoIcon={<GatewayLogo />}
        headerSecondaryInformation={[
          {
            title: 'BindPort',
            value: currentGateway?.gateway?.bindPort.toString()!
          },
          {
            title: 'Namespace',
            value: currentGateway?.gateway?.metadata?.namespace!
          },
          {
            title: 'SSL',
            value: currentGateway?.gateway?.ssl ? 'True' : 'False'
          }
        ]}
        health={
          currentGateway?.gateway?.status
            ? currentGateway?.gateway?.status?.state
            : healthConstants.Pending.value
        }
        healthMessage={'Gateway Status'}>
        <ConfigDisplayer
          content={currentGateway?.raw?.content ?? ''}
          whiteBacked
        />
      </SectionCard>
    </>
  );
};
