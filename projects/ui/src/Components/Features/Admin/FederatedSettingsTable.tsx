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
import { ReactComponent as GearIcon } from 'assets/gear-icon.svg';
import { ReactComponent as DownloadIcon } from 'assets/download-icon.svg';
import Tooltip from 'antd/lib/tooltip';
import { useNavigate } from 'react-router';
import { useListFederatedSettings } from 'API/hooks';
import { Loading } from 'Components/Common/Loading';
import { federatedGlooResourceApi } from 'API/federated-gloo';
import { doDownload } from 'download-helper';
import { IconHolder } from 'Styles/StyledComponents/icons';
import { FederatedSettings } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/federated_gloo_resources_pb';
import { colors } from 'Styles/colors';
import { DataError } from 'Components/Common/DataError';

type SettingTableFields = {
  key: string;
  name: string;
  namespace: string;
  clusters: string[];
  inNamespaces: string[];
  status: number;
  actions: FederatedSettings.AsObject;
};

// Named uniquely from the others  to avoid conflict with GRPC Class
export const FederatedSettingsTable = () => {
  const [tableData, setTableData] = React.useState<SettingTableFields[]>([]);

  const { data: settings, error: sError } = useListFederatedSettings();

  useEffect(() => {
    if (settings) {
      setTableData(
        settings
          .sort(
            (gA, gB) =>
              (gA.metadata?.name ?? '').localeCompare(
                gB.metadata?.name ?? ''
              ) ||
              (gA.metadata?.namespace ?? '').localeCompare(
                gB.metadata?.namespace ?? ''
              )
          )
          .map(setting => {
            return {
              key:
                setting.metadata?.uid ?? 'An setting was provided with no UID',
              name: setting.metadata?.name ?? '',
              namespace: setting.metadata?.namespace ?? '',
              clusters: setting.spec?.placement?.clustersList ?? [],
              inNamespaces: setting.spec?.placement?.namespacesList ?? [],
              status: setting.status?.placementStatus?.state ?? 0,
              actions: setting,
            };
          })
      );
    } else {
      setTableData([]);
    }
  }, [settings]);

  if (!!sError) {
    return <DataError error={sError} />;
  } else if (!settings) {
    return <Loading message={`Retrieving settings...`} />;
  }

  const onDownloadSetting = (setting: FederatedSettings.AsObject) => {
    if (setting.metadata) {
      federatedGlooResourceApi
        .getFederatedSettingYAML({
          name: setting.metadata.name!,
          namespace: setting.metadata.namespace!,
        })
        .then(settingYaml => {
          doDownload(
            settingYaml,
            setting.metadata?.namespace +
              '--' +
              setting.metadata?.name +
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
      render: (setting: FederatedSettings.AsObject) => (
        <Tooltip title='Download'>
          <TableActions>
            <TableActionCircle onClick={() => onDownloadSetting(setting)}>
              <DownloadIcon />
            </TableActionCircle>
          </TableActions>
        </Tooltip>
      ),
    },
  ];

  return (
    <SectionCard
      cardName={'Federated Settings'}
      logoIcon={
        <IconHolder applyColor={{ color: colors.seaBlue, strokeNotFill: true }}>
          <GearIcon />
        </IconHolder>
      }
      noPadding={true}>
      {!!sError ? (
        <DataError error={sError} />
      ) : !settings ? (
        <Loading message={'Retrieving federated Settings...'} />
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
