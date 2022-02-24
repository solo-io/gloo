import React, { useEffect } from 'react';
import { useParams, useNavigate, Routes, Route } from 'react-router';
import { colors } from 'Styles/colors';
import styled from '@emotion/styled';
import { SectionCard } from 'Components/Common/SectionCard';
import { GlooInstance } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/glooinstance_pb';
import { ProxyStatus } from 'proto/github.com/solo-io/solo-apis/api/gloo/gloo/v1/proxy_pb';
import { ReactComponent as LockIcon } from 'assets/lock-icon.svg';
import { ReactComponent as WatchedNamespacesIcon } from 'assets/watched-namespace-icon.svg';
import { ReactComponent as DocumentsIcon } from 'assets/document.svg';
import { Loading } from 'Components/Common/Loading';
import { Proxy } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/gloo_resources_pb';
import { useListProxies, useListSettings } from 'API/hooks';
import { glooResourceApi } from 'API/gloo-resource';
import { doDownload } from 'download-helper';
import YamlDisplayer from 'Components/Common/YamlDisplayer';
import { IconHolder } from 'Styles/StyledComponents/icons';
import { HealthNotificationBox } from 'Components/Common/HealthNotificationBox';
import { DataError } from 'Components/Common/DataError';
import { getTimeAsSecondsString } from 'utils';

const Description = styled.div`
  background: ${colors.januaryGrey};
  margin-bottom: 20px;
`;

const TitleRow = styled.div`
  font-size: 18px;
  line-height: 22px;
  font-weight: 500;
  margin-bottom: 5px;
`;

const BlocksContainer = styled.div`
  display: flex;
`;

const NamespaceBlock = styled.div`
  padding: 12px 20px;
  background: ${colors.marchGrey};
  border-radius: 20px;
  margin-right: 8px;
`;

export const GlooAdminWatchedNamespaces = ({
  glooInstance,
}: {
  glooInstance: GlooInstance.AsObject;
}) => {
  const { name = '', namespace = '' } = useParams();

  const { data: settings, error: sError } = useListSettings({
    name,
    namespace,
  });

  let secondaryHeaderInfo: { title: string; value: string }[] = [];
  if (settings?.length && settings[0].spec?.refreshRate !== undefined) {
    secondaryHeaderInfo.push({
      title: 'Refresh Rate',
      value: getTimeAsSecondsString(settings[0].spec?.refreshRate),
    });
  }

  return (
    <SectionCard
      logoIcon={
        <IconHolder width={20}>
          <WatchedNamespacesIcon />
        </IconHolder>
      }
      cardName={'Watched Namespaces'}
      headerSecondaryInformation={secondaryHeaderInfo}>
      <Description>
        {
          'The namespaces that Gloo Edge controllers take into consideration when watching for resources. In a usual production scenario, RBAC policies will limit the namespaces that Gloo Edge has access to.'
        }
      </Description>
      <TitleRow>
        {!!glooInstance?.spec?.controlPlane?.watchedNamespacesList.length
          ? 'You are currently watching the following namespaces:'
          : 'You are currently watching all namespaces.'}
      </TitleRow>
      {!!glooInstance?.spec?.controlPlane?.watchedNamespacesList && (
        <BlocksContainer>
          {glooInstance?.spec?.controlPlane?.watchedNamespacesList.map(
            watchedNS => (
              <NamespaceBlock>{watchedNS}</NamespaceBlock>
            )
          )}
        </BlocksContainer>
      )}
    </SectionCard>
  );
};
