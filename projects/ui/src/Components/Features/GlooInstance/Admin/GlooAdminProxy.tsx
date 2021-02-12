import React, { useEffect } from 'react';
import { useParams, useNavigate, Routes, Route } from 'react-router';
import { colors } from 'Styles/colors';
import styled from '@emotion/styled';
import { SectionCard } from 'Components/Common/SectionCard';
import { GlooInstanceSpec } from 'proto/github.com/solo-io/solo-projects/projects/gloo-fed/api/fed/v1/instance_pb';
import { ProxyStatus } from 'proto/github.com/solo-io/solo-apis/api/gloo/gloo/v1/proxy_pb';
import { ReactComponent as LockIcon } from 'assets/lock-icon.svg';
import { ReactComponent as proxyIcon } from 'assets/proxy-small-icon.svg';
import { ReactComponent as DocumentsIcon } from 'assets/document.svg';
import { Loading } from 'Components/Common/Loading';
import { Proxy } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/gloo_resources_pb';
import { useListProxies } from 'API/hooks';
import { glooResourceApi } from 'API/gloo-resource';
import { doDownload } from 'download-helper';
import YamlDisplayer from 'Components/Common/YamlDisplayer';
import { IconHolder } from 'Styles/StyledComponents/icons';
import { HealthNotificationBox } from 'Components/Common/HealthNotificationBox';
import { DataError } from 'Components/Common/DataError';

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

export const GlooAdminProxy = () => {
  const { name, namespace } = useParams();

  const { data: proxies, error: pError } = useListProxies({ name, namespace });

  const [yamlsOpen, setYamlsOpen] = React.useState<{
    [key: string]: boolean; // proxy uid
  }>({});
  useEffect(() => {
    const yamlsListByProxyUid: { [key: string]: boolean } = {};
    if (proxies?.length) {
      proxies
        .filter(proxy => !!proxy.metadata)
        .forEach(proxy => {
          yamlsListByProxyUid[proxy.metadata!.uid] = false;
        });
    }

    setYamlsOpen(yamlsListByProxyUid);

    // expand yaml for first one by default
    if (proxies?.length) {
      toggleView(proxies[0]);
    }
  }, [proxies]);

  const [swaggerContentByUid, setSwaggerContentByUid] = React.useState<{
    [key: string]: string;
  }>({});

  if (!!pError) {
    return <DataError error={pError} />;
  } else if (!proxies) {
    return <Loading message={`Retrieving proxies for ${name}...`} />;
  }

  const toggleView = (proxy: Proxy.AsObject) => {
    let viewables = { ...yamlsOpen };
    viewables[proxy.metadata!.uid] = !viewables[proxy.metadata!.uid];
    setYamlsOpen(viewables);

    if (proxy.metadata && !swaggerContentByUid[proxy.metadata!.uid]) {
      glooResourceApi
        .getProxyYAML({
          name: proxy.metadata.name,
          namespace: proxy.metadata.namespace,
          clusterName: proxy.metadata.clusterName,
        })
        .then(proxyYaml => {
          let swaggers = { ...swaggerContentByUid };
          swaggers[proxy.metadata!.uid] = proxyYaml;
          setSwaggerContentByUid(swaggers);
        });
    }
  };

  const onDownloadProxy = (proxy: Proxy.AsObject) => {
    if (proxy.metadata && !swaggerContentByUid[proxy.metadata.uid]) {
      glooResourceApi
        .getProxyYAML({
          name: proxy.metadata.name,
          namespace: proxy.metadata.namespace,
          clusterName: proxy.metadata.clusterName,
        })
        .then(proxyYaml => {
          doDownload(proxyYaml, proxy.metadata?.name + '.yaml');

          let swaggers = { ...swaggerContentByUid };
          swaggers[proxy.metadata!.uid] = proxyYaml;
          setSwaggerContentByUid(swaggers);
        });
    } else {
      doDownload(
        swaggerContentByUid[proxy.metadata!.uid],
        proxy.metadata?.name + '.yaml'
      );
    }
  };

  return (
    <div>
      {proxies?.map(proxy => {
        let secondaryHeaderInfo = [
          {
            title: 'Namespace',
            value: proxy.metadata?.namespace,
          },
        ];
        if (proxy.spec?.listenersList[0].bindPort !== undefined) {
          secondaryHeaderInfo.push({
            title: 'Bind Port',
            value: proxy.spec?.listenersList[0].bindPort.toString(),
          });
        }

        return (
          <SectionCard
            key={proxy.metadata!.name}
            logoIcon={
              <IconHolder width={20} applyColor={{ color: colors.seaBlue }}>
                <LockIcon />
              </IconHolder>
            }
            cardName={proxy.metadata!.name}
            headerSecondaryInformation={secondaryHeaderInfo}
            health={{
              state: proxy.status?.state ?? ProxyStatus.State.PENDING,
              title: 'Proxy Status',
              reason: proxy.status?.reason,
            }}>
            <HealthNotificationBox
              state={proxy?.status?.state}
              reason={proxy?.status?.reason}
            />
            <TitleRow>
              <div>Code Log (Read Only)</div>
              <Actionables>
                <div onClick={() => toggleView(proxy)}>
                  {yamlsOpen[proxy.metadata!.uid] ? 'Hide' : 'View'} Raw Config
                </div>
                <div onClick={() => onDownloadProxy(proxy)}>
                  <DocumentsIcon /> {proxy.metadata?.name}.yaml
                </div>
              </Actionables>
            </TitleRow>
            {yamlsOpen[proxy.metadata!.uid] &&
              (!!swaggerContentByUid[proxy.metadata!.uid] ? (
                <YamlDisplayer
                  contentString={swaggerContentByUid[proxy.metadata!.uid]}
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
