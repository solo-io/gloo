import styled from '@emotion/styled';
import { ReactComponent as ProxyLogo } from 'assets/proxy-icon.svg';
import { ConfigDisplayer } from 'Components/Common/DisplayOnly/ConfigDisplayer';
import { FileDownloadLink } from 'Components/Common/FileDownloadLink';
import { SectionCard } from 'Components/Common/SectionCard';
import * as React from 'react';
import { proxyAPI } from 'store/proxy/api';
import { colors, healthConstants } from 'Styles';
import useSWR from 'swr';

const InsideHeader = styled.div`
  display: flex;
  justify-content: space-between;
  font-size: 18px;
  line-height: 22px;
  margin-bottom: 18px;
  color: ${colors.novemberGrey};
`;

const ProxyLogoFullSize = styled(ProxyLogo)`
  width: 33px !important;
  max-height: none !important;
`;

interface Props {}

export const Proxys = (props: Props) => {
  const { data: proxiesList, error } = useSWR(
    'listProxies',
    proxyAPI.getListProxies
  );

  if (!proxiesList) {
    return <div>Loading...</div>;
  }
  if (!proxiesList.length) {
    return <div>You have no Proxy configurations.</div>;
  }

  return (
    <>
      {proxiesList.map((proxy, ind) => {
        return (
          <SectionCard
            key={proxy.proxy!.metadata!.name + ind}
            cardName={proxy.proxy!.metadata!.name}
            logoIcon={<ProxyLogoFullSize />}
            headerSecondaryInformation={[
              {
                title: 'Namespace',
                value: proxy.proxy!.metadata!.namespace
              },
              {
                title: 'Listener Ports',
                value:
                  proxy.proxy!.listenersList.map(l => l.bindPort).join(', ') ||
                  ''
              }
            ]}
            health={
              proxy.proxy!.status
                ? proxy.proxy!.status!.state
                : healthConstants.Pending.value
            }
            healthMessage={'Proxy Status'}>
            <InsideHeader>
              <div>Code Log (Read Only)</div>{' '}
              {!!proxy.raw && (
                <FileDownloadLink
                  fileName={proxy.raw.fileName}
                  fileContent={proxy.raw.content}
                />
              )}
            </InsideHeader>
            {!!proxy.raw && <ConfigDisplayer content={proxy.raw.content} />}
          </SectionCard>
        );
      })}
    </>
  );
};
