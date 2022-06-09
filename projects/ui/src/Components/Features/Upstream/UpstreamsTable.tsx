import styled from '@emotion/styled/macro';
import Tooltip from 'antd/lib/tooltip';
import { glooResourceApi } from 'API/gloo-resource';
import {
  useIsGlooFedEnabled,
  useListClusterDetails,
  useListGlooInstances,
  useListUpstreams,
} from 'API/hooks';
import { ReactComponent as DownloadIcon } from 'assets/download-icon.svg';
import { ReactComponent as FailoverIcon } from 'assets/GlooFed-Specific/failover-icon.svg';
import { ReactComponent as UpstreamIcon } from 'assets/upstream-icon.svg';
import { DataError } from 'Components/Common/DataError';
import { SectionCard } from 'Components/Common/SectionCard';
import { RenderSimpleLink, SimpleLinkProps } from 'Components/Common/SoloLink';
import {
  RenderCluster,
  RenderStatus,
  SoloTable,
  TableActionCircle,
  TableActions,
} from 'Components/Common/SoloTable';
import { doDownload } from 'download-helper';
import { UpstreamStatus } from 'proto/github.com/solo-io/solo-apis/api/gloo/gloo/v1/upstream_pb';
import { Upstream } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/gloo_resources_pb';
import React, { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router';
import { colors } from 'Styles/colors';
import { IconHolder } from 'Styles/StyledComponents/icons';

const UpstreamIconHolder = styled.div`
  display: flex;
  align-items: center;
  justify-items: center;

  svg {
    width: 35px;
    max-width: none;
  }
`;

type TableHolderProps = { wholePage?: boolean };
const TableHolder = styled.div<TableHolderProps>`
  ${(props: TableHolderProps) =>
    props.wholePage
      ? ''
      : `
    table thead.ant-table-thead tr th {
      background: ${colors.marchGrey};
    }
  `};
`;

type Props = {
  statusFilter?: UpstreamStatus.StateMap[keyof UpstreamStatus.StateMap];
  nameFilter?: string;
  glooInstanceFilter?: {
    name: string;
    namespace: string;
  };
};
export const UpstreamsTable = (props: Props & TableHolderProps) => {
  const { name, namespace } = useParams();
  const [offset, setOffset] = React.useState(0);
  const limit = 10;
  const navigate = useNavigate();

  const [tableData, setTableData] = React.useState<
    {
      key: string;
      name: SimpleLinkProps;
      namespace: string;
      glooInstance?: { name: string; namespace: string };
      cluster: string;
      failover: boolean;
      status: number;
      actions: Upstream.AsObject;
    }[]
  >([]);

  const { data: upstreamsResponse, error: upstreamsResponseError } =
    useListUpstreams(
      { name, namespace },
      { limit, offset },
      props.nameFilter,
      props.statusFilter
    );
  const { data: glooInstances, error: glooError } = useListGlooInstances();
  const { data: clusterDetailsList, error: cError } = useListClusterDetails();
  const { data: glooFedCheckResponse, error: glooFedCheckError } =
    useIsGlooFedEnabled();

  const isGlooFedEnabled = glooFedCheckResponse?.enabled;
  const upstreams = upstreamsResponse?.upstreamsList;
  const total = upstreamsResponse?.total ?? 0;

  const multipleClustersOrInstances =
    (clusterDetailsList && clusterDetailsList.length > 1) ||
    (glooInstances && glooInstances.length > 1);

  const [page, setPage] = useState(1);
  useEffect(() => {
    setPage(1);
  }, [
    props.nameFilter,
    props.statusFilter,
    props.glooInstanceFilter,
    isGlooFedEnabled,
  ]);
  useEffect(() => {
    setOffset(limit * (page - 1));
  }, [page]);

  useEffect(() => {
    if (upstreams) {
      setTableData(
        upstreams.map(upstream => {
          const glooInstNamespace = upstream.glooInstance?.namespace;
          const glooInstName = upstream.glooInstance?.name;
          const upstreamCluster = upstream.metadata?.clusterName ?? '';
          const upstreamNamespace = upstream.metadata?.namespace ?? '';
          const upstreamName = upstream.metadata?.name ?? '';
          const link = isGlooFedEnabled
            ? `/gloo-instances/${glooInstNamespace}/${glooInstName}/upstreams/${upstreamCluster}/${upstreamNamespace}/${upstreamName}`
            : `/gloo-instances/${glooInstNamespace}/${glooInstName}/upstreams/${upstreamNamespace}/${upstreamName}`;
          return {
            key:
              upstream.metadata?.uid ?? 'An upstream was provided with no UID',
            name: {
              displayElement: upstreamName,
              link,
            },
            namespace: upstreamNamespace,
            glooInstance: upstream.glooInstance,
            cluster: upstreamCluster,
            failover: !!upstream.spec?.failover,
            status: upstream.status?.state ?? UpstreamStatus.State.PENDING,
            actions: upstream,
          };
        })
      );
    } else {
      setTableData([]);
    }
  }, [
    upstreams,
    props.nameFilter,
    props.statusFilter,
    props.glooInstanceFilter,
    isGlooFedEnabled,
    props.wholePage,
  ]);

  if (!!upstreamsResponseError) {
    return <DataError error={upstreamsResponseError} />;
  }

  const onDownloadUpstream = (upstream: Upstream.AsObject) => {
    if (upstream.metadata) {
      glooResourceApi
        .getUpstreamYAML({
          name: upstream.metadata.name,
          namespace: upstream.metadata.namespace,
          clusterName: upstream.metadata.clusterName,
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

  const renderGlooInstanceList = (glooInstance: {
    name: string;
    namespace: string;
  }) => {
    return (
      <div
        onClick={() =>
          navigate(
            `/gloo-instances/${glooInstance.namespace}/${glooInstance.name}/`
          )
        }>
        {glooInstance.name}
      </div>
    );
  };

  const renderFailover = (failoverExists: boolean) => {
    return failoverExists ? (
      <IconHolder>
        <FailoverIcon />
      </IconHolder>
    ) : (
      <React.Fragment />
    );
  };

  let columns: any = [
    {
      title: 'Name',
      dataIndex: 'name',
      width: 200,
      render: RenderSimpleLink,
    },
    {
      title: 'Namespace',
      dataIndex: 'namespace',
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
      render: (upstream: Upstream.AsObject) => (
        <Tooltip title='Download'>
          <TableActions>
            <TableActionCircle onClick={() => onDownloadUpstream(upstream)}>
              <DownloadIcon />
            </TableActionCircle>
          </TableActions>
        </Tooltip>
      ),
    },
  ];

  if (props.wholePage && multipleClustersOrInstances) {
    columns.splice(2, 0, {
      title: 'Cluster',
      dataIndex: 'cluster',
      render: RenderCluster,
    });
  }
  if (props.wholePage) {
    columns.splice(2, 0, {
      title: 'Gloo Instance',
      dataIndex: 'glooInstance',
      render: renderGlooInstanceList,
    });
  }

  return (
    <TableHolder wholePage={props.wholePage}>
      <SoloTable
        pagination={{
          total,
          pageSize: limit,
          current: page,
          onChange: newPage => setPage(newPage),
        }}
        removePaging={total <= limit}
        columns={columns}
        dataSource={tableData}
        removeShadows
        curved={props.wholePage}
      />
    </TableHolder>
  );
};

export const UpstreamsPageTable = (props: Props) => {
  return (
    <SectionCard
      cardName={'Upstreams'}
      logoIcon={
        <UpstreamIconHolder>
          <UpstreamIcon />
        </UpstreamIconHolder>
      }
      noPadding={true}>
      <UpstreamsTable {...props} wholePage={true} />
    </SectionCard>
  );
};
