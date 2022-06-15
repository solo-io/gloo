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
import { di } from 'react-magnetic-di/macro';
import { SectionCard } from 'Components/Common/SectionCard';
import { ReactComponent as DownloadIcon } from 'assets/download-icon.svg';
import Tooltip from 'antd/lib/tooltip';
import { useNavigate } from 'react-router';
import { useListFederatedRateLimits } from 'API/hooks';
import { Loading } from 'Components/Common/Loading';
import { doDownload } from 'download-helper';
import { IconHolder } from 'Styles/StyledComponents/icons';
import { colors } from 'Styles/colors';
import { federatedEnterpriseGlooResourceApi } from 'API/federated-enterprise-gloo';
import { DataError } from 'Components/Common/DataError';
import { ReactComponent as ProxyLockIcon } from 'assets/lock-icon.svg';
import { FederatedRateLimitConfig } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/federated_ratelimit_resources_pb';

type RateLimitTableFields = {
  key: string;
  name: string;
  namespace: string;
  clusters: string[];
  inNamespaces: string[];
  status: number;
  actions: FederatedRateLimitConfig.AsObject;
};

export const FederatedRateLimits = () => {
  const [tableData, setTableData] = React.useState<RateLimitTableFields[]>([]);
  di(useListFederatedRateLimits);
  const { data: rateLimits, error: fedRLError } = useListFederatedRateLimits();

  useEffect(() => {
    if (rateLimits) {
      setTableData(
        rateLimits
          .sort(
            (gA, gB) =>
              (gA.metadata?.name ?? '').localeCompare(
                gB.metadata?.name ?? ''
              ) ||
              (gA.metadata?.namespace ?? '').localeCompare(
                gB.metadata?.namespace ?? ''
              )
          )
          .map(rateLimit => {
            return {
              key:
                rateLimit.metadata?.uid ??
                'A rate limit was provided with no UID',
              name: rateLimit.metadata?.name ?? '',
              namespace: rateLimit.metadata?.namespace ?? '',
              clusters: rateLimit.spec?.placement?.clustersList ?? [],
              inNamespaces: rateLimit.spec?.placement?.namespacesList ?? [],
              status: rateLimit.status?.placementStatus?.state ?? 0,
              actions: rateLimit,
            };
          })
      );
    } else {
      setTableData([]);
    }
  }, [rateLimits]);

  const onDownloadRateLimit = (
    rateLimit: FederatedRateLimitConfig.AsObject
  ) => {
    if (rateLimit.metadata) {
      federatedEnterpriseGlooResourceApi
        .getFederatedRateLimitYAML({
          name: rateLimit.metadata.name!,
          namespace: rateLimit.metadata.namespace!,
        })
        .then(rateLimitYaml => {
          doDownload(
            rateLimitYaml.yaml,
            rateLimit.metadata?.namespace +
              '--' +
              rateLimit.metadata?.name +
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
      render: (rateLimit: FederatedRateLimitConfig.AsObject) => (
        <Tooltip title='Download'>
          <TableActions>
            <TableActionCircle onClick={() => onDownloadRateLimit(rateLimit)}>
              <DownloadIcon />
            </TableActionCircle>
          </TableActions>
        </Tooltip>
      ),
    },
  ];

  return (
    <SectionCard
      cardName={'Federated Rate Limits'}
      logoIcon={
        <IconHolder width={18} applyColor={{ color: colors.seaBlue }}>
          <ProxyLockIcon />
        </IconHolder>
      }
      noPadding={true}>
      {!!fedRLError ? (
        <DataError error={fedRLError} />
      ) : !rateLimits ? (
        <Loading message={'Retrieving federated rate limits...'} />
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
