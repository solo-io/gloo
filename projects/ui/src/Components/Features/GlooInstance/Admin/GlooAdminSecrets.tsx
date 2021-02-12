import React, { useEffect } from 'react';
import { useParams, useNavigate, Routes, Route } from 'react-router';
import { colors } from 'Styles/colors';
import styled from '@emotion/styled';
import { SectionCard } from 'Components/Common/SectionCard';
import { GlooInstanceSpec } from 'proto/github.com/solo-io/solo-projects/projects/gloo-fed/api/fed/v1/instance_pb';
import { ProxyStatus } from 'proto/github.com/solo-io/solo-apis/api/gloo/gloo/v1/proxy_pb';
import { ReactComponent as LockIcon } from 'assets/lock-icon.svg';
import { ReactComponent as SecretsIcon } from 'assets/cloud-key-icon.svg';
import { ReactComponent as DocumentsIcon } from 'assets/document.svg';
import { Loading } from 'Components/Common/Loading';
import { Proxy } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/gloo_resources_pb';
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

const NamespaceBlock = styled.div`
  padding: 12px 20px;
  background: ${colors.marchGrey};
  border-radius: 20px;
`;

type TableProps = {
  key: string;
  name: string;
  namespace: string;
  type: any;
};

export const GlooAdminSecrets = () => {
  const { name, namespace } = useParams();

  const { data: settings, error: sError } = useListSettings({
    name,
    namespace,
  });

  const [tableData, setTableData] = React.useState<TableProps[]>([]);

  return (
    <SectionCard
      logoIcon={
        <IconHolder width={20}>
          <SecretsIcon />
        </IconHolder>
      }
      cardName={'Secrets'}></SectionCard>
  );
};
