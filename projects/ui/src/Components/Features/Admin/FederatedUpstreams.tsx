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
import { ReactComponent as UpstreamIcon } from 'assets/upstream-icon.svg';
import { ReactComponent as DownloadIcon } from 'assets/download-icon.svg';
import { ReactComponent as FailoverIcon } from 'assets/GlooFed-Specific/failover-icon.svg';
import { useNavigate } from 'react-router';
import { useListFederatedUpstreams } from 'API/hooks';
import { Loading } from 'Components/Common/Loading';
import { FederatedUpstream } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/federated_gloo_resources_pb';
import { federatedGlooResourceApi } from 'API/federated-gloo';
import { doDownload } from 'download-helper';
import { IconHolder } from 'Styles/StyledComponents/icons';
import { DataError } from 'Components/Common/DataError';

const UpstreamIconHolder = styled.div`
  display: flex;
  align-items: center;
  justify-items: center;

  svg {
    width: 35px;
    max-width: none;
  }
`;

type UpstreamTableFields = {
  key: string;
  name: string;
  namespace: string;
  clusters: string[];
  inNamespaces: string[];
  failover: boolean;
  status: number;
  actions: FederatedUpstream.AsObject;
};

export const FederatedUpstreams = () => {
  const [tableData, setTableData] = React.useState<UpstreamTableFields[]>([]);

  const { data: upstreams, error: fedUError } = useListFederatedUpstreams();

  useEffect(() => {
    if (upstreams) {
      setTableData(
        upstreams.map(upstream => {
          return {
            key:
              upstream.metadata?.uid ?? 'An upstream was provided with no UID',
            name: upstream.metadata?.name ?? '',
            namespace: upstream.metadata?.namespace ?? '',
            clusters: upstream.spec?.placement?.clustersList ?? [],
            inNamespaces: upstream.spec?.placement?.namespacesList ?? [],
            status: upstream.status?.placementStatus?.state ?? 0,
            failover: !!upstream.spec?.template?.spec?.failover,
            actions: upstream,
          };
        })
      );
    } else {
      setTableData([]);
    }
  }, [upstreams]);

  if (!!fedUError) {
    return <DataError error={fedUError} />;
  } else if (!upstreams) {
    return <Loading message={`Retrieving upstreams...`} />;
  }

  const onDownloadUpstream = (upstream: FederatedUpstream.AsObject) => {
    if (upstream.metadata) {
      federatedGlooResourceApi
        .getFederatedUpstreamYAML({
          name: upstream.metadata.name!,
          namespace: upstream.metadata.namespace!,
        })
        .then(upstreamYaml => {
          doDownload(
            upstreamYaml,
            upstream.metadata?.namespace +
              '--' +
              upstream.metadata?.name +
              '.yaml'
          );
        });
    }
  };

  const renderFailover = (failoverExists: boolean) => {
    return failoverExists ? (
      <div>
        <FailoverIcon />
      </div>
    ) : (
      <React.Fragment />
    );
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
      title: 'Failover',
      dataIndex: 'failover',
      render: renderFailover,
    },
    {
      title: 'Status',
      dataIndex: 'status',
      render: RenderStatus,
    },

    {
      title: 'Actions',
      dataIndex: 'actions',
      render: (upstream: FederatedUpstream.AsObject) => (
        <TableActions>
          <TableActionCircle onClick={() => onDownloadUpstream(upstream)}>
            <DownloadIcon />
          </TableActionCircle>
        </TableActions>
      ),
    },
  ];

  return (
    <SectionCard
      cardName={'Federated Upstreams'}
      logoIcon={
        <UpstreamIconHolder>
          <UpstreamIcon />
        </UpstreamIconHolder>
      }
      noPadding={true}>
      {!!fedUError ? (
        <DataError error={fedUError} />
      ) : !upstreams ? (
        <Loading message={'Retrieving federated Upstreams...'} />
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
