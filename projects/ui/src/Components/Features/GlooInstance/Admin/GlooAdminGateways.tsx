import React, { useEffect } from 'react';
import { useParams, useNavigate, Routes, Route } from 'react-router';
import { colors } from 'Styles/colors';
import styled from '@emotion/styled';
import { SectionCard } from 'Components/Common/SectionCard';
import { ReactComponent as GatewayIcon } from 'assets/gateway-small-icon.svg';
import { ReactComponent as DocumentsIcon } from 'assets/document.svg';
import { useListGateways } from 'API/hooks';
import { Loading } from 'Components/Common/Loading';
import { GatewayStatus } from 'proto/github.com/solo-io/solo-apis/api/gloo/gateway/v1/gateway_pb';
import { Gateway } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/gateway_resources_pb';
import { gatewayResourceApi } from 'API/gateway-resources';
import YamlDisplayer from 'Components/Common/YamlDisplayer';
import { doDownload } from 'download-helper';
import { HealthNotificationBox } from 'Components/Common/HealthNotificationBox';
import { DataError } from 'Components/Common/DataError';

const GatewayIconHolder = styled.div`
  display: flex;
  align-items: center;
  justify-content: flex-start;

  svg {
    width: 33px;
    margin-left: -1px;

    * {
      fill: ${colors.seaBlue};
    }
  }
`;

const TitleRow = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
  font-size: 18px;
`;

const Actionables = styled.div`
  display: flex;
  align-items: center;
  color: ${colors.seaBlue};

  > div {
    display: flex;
    align-items: center;
    margin-left: 20px;
    cursor: pointer;
  }

  svg {
    margin-right: 8px;
  }
`;

export const GlooAdminGateways = () => {
  const { name, namespace } = useParams();

  const { data: gateways, error: gError } = useListGateways({
    name,
    namespace,
  });

  const [yamlsOpen, setYamlsOpen] = React.useState<{
    [key: string]: boolean; // gateway Uid
  }>({});
  useEffect(() => {
    const yamlsListByGatewayUid: { [key: string]: boolean } = {};
    if (gateways?.length) {
      gateways
        .filter(gateway => !!gateway.metadata)
        .forEach(gateway => {
          yamlsListByGatewayUid[gateway.metadata!.uid] = false;
        });
    }

    setYamlsOpen(yamlsListByGatewayUid);

    // expand yaml for first one by default
    if (gateways?.length) {
      toggleView(gateways[0]);
    }
  }, [gateways]);

  const [swaggerContentByUid, setSwaggerContentByUid] = React.useState<{
    [key: string]: string;
  }>({});

  if (!!gError) {
    return <DataError error={gError} />;
  } else if (!gateways) {
    return <Loading message={`Retrieving gateways for ${name}...`} />;
  }

  const toggleView = (gateway: Gateway.AsObject) => {
    let viewables = { ...yamlsOpen };
    viewables[gateway.metadata!.uid] = !viewables[gateway.metadata!.uid];
    setYamlsOpen(viewables);

    if (gateway.metadata && !swaggerContentByUid[gateway.metadata.uid]) {
      gatewayResourceApi
        .getGatewayYAML({
          name: gateway.metadata.name,
          namespace: gateway.metadata.namespace,
          clusterName: gateway.metadata.clusterName,
        })
        .then(gatewayYaml => {
          let swaggers = { ...swaggerContentByUid };
          swaggers[gateway.metadata!.uid] = gatewayYaml;
          setSwaggerContentByUid(swaggers);
        });
    }
  };

  const onDownloadGateway = (gateway: Gateway.AsObject) => {
    // meta should always be there, so really the test is for the 2nd check
    if (gateway.metadata && !swaggerContentByUid[gateway.metadata.uid]) {
      gatewayResourceApi
        .getGatewayYAML({
          name: gateway.metadata.name,
          namespace: gateway.metadata.namespace,
          clusterName: gateway.metadata.clusterName,
        })
        .then(gatewayYaml => {
          doDownload(gatewayYaml, gateway.metadata?.name + '.yaml');

          let swaggers = { ...swaggerContentByUid };
          swaggers[gateway.metadata!.uid] = gatewayYaml;
          setSwaggerContentByUid(swaggers);
        });
    } else {
      doDownload(
        swaggerContentByUid[gateway.metadata!.uid],
        gateway.metadata?.name + '.yaml'
      );
    }
  };

  return (
    <div>
      {gateways?.map(gateway => {
        let secondaryHeaderInfo: { title: string; value: string | number }[] = [
          {
            title: 'Namespace',
            value: gateway.metadata?.namespace ?? '',
          },
        ];

        if (!!gateway.spec) {
          secondaryHeaderInfo.splice(0, 0, {
            title: 'BindPort',
            value: gateway.spec?.bindPort,
          });
          secondaryHeaderInfo.push({
            title: 'SSL',
            value: gateway.spec?.ssl ? 'True' : 'False',
          });
        }

        return (
          <SectionCard
            key={gateway.metadata!.uid}
            logoIcon={
              <GatewayIconHolder>
                <GatewayIcon />
              </GatewayIconHolder>
            }
            cardName={gateway.metadata?.name ?? ''}
            headerSecondaryInformation={secondaryHeaderInfo}
            health={{
              state: gateway.status?.state ?? GatewayStatus.State.PENDING,
              title: 'Gateway Status',
              reason: gateway.status?.reason,
            }}>
            <HealthNotificationBox
              state={gateway?.status?.state}
              reason={gateway?.status?.reason}
            />
            <TitleRow>
              <div>Configuration Settings</div>{' '}
              <Actionables>
                <div onClick={() => toggleView(gateway)}>
                  {yamlsOpen[gateway.metadata!.uid] ? 'Hide' : 'View'} Raw
                  Config
                </div>
                <div onClick={() => onDownloadGateway(gateway)}>
                  <DocumentsIcon /> {gateway.metadata?.name}.yaml
                </div>
              </Actionables>
            </TitleRow>
            {yamlsOpen[gateway.metadata!.uid] &&
              (!!swaggerContentByUid[gateway.metadata!.uid] ? (
                <YamlDisplayer
                  description={
                    <div>
                      Below are gateway configuration settings for you to
                      review. For more information on these settings, please
                      visit our{' '}
                      <a
                        href='https://docs.solo.io/gloo/latest/reference/api/github.com/solo-io/gloo/projects/gateway/api/v1/gateway.proto.sk/'
                        target='_blank'>
                        hcm plugin documentation
                      </a>
                    </div>
                  }
                  contentString={swaggerContentByUid[gateway.metadata!.uid]}
                />
              ) : (
                <Loading message={'Retrieving configuration...'} />
              ))}
          </SectionCard>
        );
      })}
    </div>
  );
};
