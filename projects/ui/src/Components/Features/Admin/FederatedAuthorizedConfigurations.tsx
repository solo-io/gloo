import React, { useEffect } from 'react';
import styled from '@emotion/styled/macro';
import {
  SoloTable,
  RenderStatus,
  TableActionCircle,
  TableActions,
  RenderClustersList,
  RenderExpandableNamesList,
} from 'Components/Common/SoloTable';
import { SectionCard } from 'Components/Common/SectionCard';
import { ReactComponent as CloudIcon } from 'assets/auth-cloud-icon.svg';
import { ReactComponent as DownloadIcon } from 'assets/download-icon.svg';
import { useNavigate } from 'react-router';
import { useListFederatedAuthConfigs } from 'API/hooks';
import { Loading } from 'Components/Common/Loading';
import { doDownload } from 'download-helper';
import { IconHolder } from 'Styles/StyledComponents/icons';
import { colors } from 'Styles/colors';
import { FederatedAuthConfig } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/federated_enterprise_gloo_resources_pb';
import { federatedEnterpriseGlooResourceApi } from 'API/federated-enterprise-gloo';
import { DataError } from 'Components/Common/DataError';

type AuthConfigTableFields = {
  key: string;
  name: string;
  namespace: string;
  clusters: string[];
  inNamespaces: string[];
  status: number;
  actions: FederatedAuthConfig.AsObject;
};

export const FederatedAuthorizedConfigurations = () => {
  const [tableData, setTableData] = React.useState<AuthConfigTableFields[]>([]);

  const {
    data: authConfigs,
    error: fedACError,
  } = useListFederatedAuthConfigs();

  useEffect(() => {
    if (authConfigs) {
      setTableData(
        authConfigs
        .sort((gA, gB) => (gA.metadata?.name ?? '').localeCompare(gB.metadata?.name ?? '') || (gA.metadata?.namespace ?? '').localeCompare(gB.metadata?.namespace ?? ''))
        .map(authConfig => {
          return {
            key:
              authConfig.metadata?.uid ??
              'A authorized configuration was provided with no UID',
            name: authConfig.metadata?.name ?? '',
            namespace: authConfig.metadata?.namespace ?? '',
            clusters: authConfig.spec?.placement?.clustersList ?? [],
            inNamespaces: authConfig.spec?.placement?.namespacesList ?? [],
            status: authConfig.status?.placementStatus?.state ?? 0,
            actions: authConfig,
          };
        })
      );
    } else {
      setTableData([]);
    }
  }, [authConfigs]);

  if (!!fedACError) {
    return <DataError error={fedACError} />;
  } else if (!authConfigs) {
    return <Loading message={`Retrieving configuration information...`} />;
  }

  const onDownloadAuthConfig = (authConfig: FederatedAuthConfig.AsObject) => {
    if (authConfig.metadata) {
      federatedEnterpriseGlooResourceApi
        .getFederatedAuthConfigYAML({
          name: authConfig.metadata.name!,
          namespace: authConfig.metadata.namespace!,
        })
        .then(authConfigYaml => {
          doDownload(
            authConfigYaml,
            authConfig.metadata?.namespace +
              '--' +
              authConfig.metadata?.name +
              '.yaml'
          );
        });
    }
  };

  let columns: any = [
    {
      title: 'Name',
      dataIndex: 'name',
      width: 200,
    },
    {
      title: 'Namespace',
      dataIndex: 'namespace',
    },
    {
      title: 'Placement Clusters',
      dataIndex: 'clusters',
      render: RenderClustersList,
    },
    {
      title: 'Placement Namespaces',
      dataIndex: 'inNamespaces',
      render: RenderExpandableNamesList,
    },
    {
      title: 'Status',
      dataIndex: 'status',
      render: RenderStatus,
    },
    {
      title: 'Actions',
      dataIndex: 'actions',
      render: (authConfig: FederatedAuthConfig.AsObject) => (
        <TableActions>
          <TableActionCircle onClick={() => onDownloadAuthConfig(authConfig)}>
            <DownloadIcon />
          </TableActionCircle>
        </TableActions>
      ),
    },
  ];

  return (
    <SectionCard
      cardName={'Federated Authorized Configurations'}
      logoIcon={
        <IconHolder width={30} applyColor={{ color: colors.seaBlue }}>
          <CloudIcon />
        </IconHolder>
      }
      noPadding={true}>
      {!!fedACError ? (
        <DataError error={fedACError} />
      ) : !authConfigs ? (
        <Loading
          message={'Retrieving federated authorized configurations...'}
        />
      ) : (
        <SoloTable
          columns={columns}
          dataSource={tableData}
          removePaging
          removeShadows
          curved={true}
        />
      )}
    </SectionCard>
  );
};
