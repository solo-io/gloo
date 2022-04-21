import React, { useEffect } from 'react';
import styled from '@emotion/styled/macro';
import {
  SoloTable,
  RenderStatus,
  TableActionCircle,
  TableActions,
  RenderCluster,
  RenderClustersList,
  RenderExpandableNamesList,
} from 'Components/Common/SoloTable';
import { SectionCard } from 'Components/Common/SectionCard';
import { ReactComponent as UpstreamGroupIcon } from 'assets/upstream-group-icon.svg';
import { ReactComponent as DownloadIcon } from 'assets/download-icon.svg';
import Tooltip from 'antd/lib/tooltip';
import { useNavigate } from 'react-router';
import { useListFederatedUpstreamGroups } from 'API/hooks';
import { Loading } from 'Components/Common/Loading';
import { federatedGlooResourceApi } from 'API/federated-gloo';
import { doDownload } from 'download-helper';
import { FederatedUpstreamGroup } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/federated_gloo_resources_pb';
import { IconHolder } from 'Styles/StyledComponents/icons';
import { DataError } from 'Components/Common/DataError';

type UpstreamGroupsTableFields = {
  key: string;
  name: string;
  namespace: string;
  clusters: string[];
  inNamespaces: string[];
  status: number;
  actions: FederatedUpstreamGroup.AsObject;
};

export const FederatedUpstreamGroups = () => {
  const [tableData, setTableData] = React.useState<UpstreamGroupsTableFields[]>(
    []
  );

  const { data: upstreamGroups, error: fedUGError } =
    useListFederatedUpstreamGroups();

  useEffect(() => {
    if (upstreamGroups) {
      setTableData(
        upstreamGroups
          .sort(
            (gA, gB) =>
              (gA.metadata?.name ?? '').localeCompare(
                gB.metadata?.name ?? ''
              ) ||
              (gA.metadata?.namespace ?? '').localeCompare(
                gB.metadata?.namespace ?? ''
              )
          )
          .map(upstreamGroup => {
            return {
              key:
                upstreamGroup.metadata?.uid ??
                'An upstream group was provided with no UID',
              name: upstreamGroup.metadata?.name ?? '',
              namespace: upstreamGroup.metadata?.namespace ?? '',
              clusters: upstreamGroup.spec?.placement?.clustersList ?? [],
              inNamespaces: upstreamGroup.spec?.placement?.namespacesList ?? [],
              status: upstreamGroup.status?.placementStatus?.state ?? 0,
              actions: upstreamGroup,
            };
          })
      );
    } else {
      setTableData([]);
    }
  }, [upstreamGroups]);

  if (!!fedUGError) {
    return <DataError error={fedUGError} />;
  } else if (!upstreamGroups) {
    return <Loading message={`Retrieving upstream groups...`} />;
  }

  const onDownloadUpstreamGroups = (
    upstreamGroup: FederatedUpstreamGroup.AsObject
  ) => {
    if (upstreamGroup.metadata) {
      federatedGlooResourceApi
        .getFederatedUpstreamGroupYAML({
          name: upstreamGroup.metadata.name!,
          namespace: upstreamGroup.metadata.namespace!,
        })
        .then(upstreamGroupYaml => {
          doDownload(
            upstreamGroupYaml,
            upstreamGroup.metadata?.namespace +
              '--' +
              upstreamGroup.metadata?.name +
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
      render: (upstreamGroup: FederatedUpstreamGroup.AsObject) => (
        <Tooltip title='Download'>
          <TableActions>
            <TableActionCircle
              onClick={() => onDownloadUpstreamGroups(upstreamGroup)}>
              <DownloadIcon />
            </TableActionCircle>
          </TableActions>
        </Tooltip>
      ),
    },
  ];

  return (
    <SectionCard
      cardName={'Federated Upstream Groups'}
      logoIcon={
        <IconHolder width={25}>
          <UpstreamGroupIcon />
        </IconHolder>
      }
      noPadding={true}>
      {!!fedUGError ? (
        <DataError error={fedUGError} />
      ) : !upstreamGroups ? (
        <Loading message={'Retrieving federated Upstream Groups...'} />
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
