import styled from '@emotion/styled/macro';
import Tooltip from 'antd/lib/tooltip';
import { glooResourceApi } from 'API/gloo-resource';
import { useIsGlooFedEnabled, useListUpstreamGroups } from 'API/hooks';
import { ReactComponent as DownloadIcon } from 'assets/download-icon.svg';
import { ReactComponent as UpstreamGroupIcon } from 'assets/upstream-group-icon.svg';
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
import { UpstreamGroup } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/gloo_resources_pb';
import React, { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router';
import { colors } from 'Styles/colors';

const UpstreamGroupIconHolder = styled.div`
  display: flex;
  align-items: center;
  justify-items: center;

  svg {
    width: 25px;
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
};
export const UpstreamGroupsTable = (props: Props & TableHolderProps) => {
  const { name, namespace } = useParams();
  const navigate = useNavigate();

  const [tableData, setTableData] = React.useState<
    {
      key: string;
      name: SimpleLinkProps;
      namespace: string;
      glooInstance?: { name: string; namespace: string };
      cluster: string;
      upstreamGroupsCount: number;
      status: number;
      actions: UpstreamGroup.AsObject;
    }[]
  >([]);

  const [offset, setOffset] = useState(0);
  const [limit, setLimit] = useState(5);
  const { data: ugResponse, error: upstreamGroupsError } =
    useListUpstreamGroups(
      !!name && !!namespace
        ? {
            name,
            namespace,
          }
        : undefined,
      { limit, offset },
      props.nameFilter,
      props.statusFilter
    );
  const upstreamGroups = ugResponse?.upstreamGroupsList;
  const total = ugResponse?.total ?? 0;

  const { data: glooFedCheckResponse, error: glooFedCheckError } =
    useIsGlooFedEnabled();
  const isGlooFedEnabled = glooFedCheckResponse?.enabled;

  const [page, setPage] = useState(1);
  useEffect(() => {
    setPage(1);
  }, [props.nameFilter, props.statusFilter, isGlooFedEnabled]);
  useEffect(() => {
    setOffset(limit * (page - 1));
  }, [page]);

  useEffect(() => {
    if (upstreamGroups) {
      setTableData(
        upstreamGroups.map(group => ({
          key: group.metadata?.uid ?? 'An group was provided with no UID',
          name: {
            displayElement: group.metadata?.name ?? '',
            link: group.metadata
              ? isGlooFedEnabled
                ? `/gloo-instances/${group.glooInstance?.namespace}/${group.glooInstance?.name}/upstream-groups/${group.metadata.clusterName}/${group.metadata.namespace}/${group.metadata.name}/`
                : `/gloo-instances/${group.glooInstance?.namespace}/${group.glooInstance?.name}/upstream-groups/${group.metadata.namespace}/${group.metadata.name}/`
              : '',
          },
          namespace: group.metadata?.namespace ?? '',
          glooInstance: group.glooInstance,
          cluster: group.metadata?.clusterName ?? '',
          upstreamGroupsCount: group.spec?.destinationsList.length ?? 0,
          status: group.status?.state ?? UpstreamStatus.State.PENDING,
          actions: group,
        }))
      );
    } else {
      setTableData([]);
    }
  }, [
    upstreamGroups,
    props.nameFilter,
    props.statusFilter,
    isGlooFedEnabled,
    props.wholePage,
  ]);

  if (!!upstreamGroupsError) {
    return <DataError error={upstreamGroupsError} />;
  }

  const onDownloadUpstream = (group: UpstreamGroup.AsObject) => {
    if (group.metadata) {
      glooResourceApi
        .getUpstreamGroupYAML({
          name: group.metadata.name,
          namespace: group.metadata.namespace,
          clusterName: group.metadata.clusterName,
        })
        .then(groupYaml => {
          doDownload(
            groupYaml,
            group.metadata?.namespace + '--' + group.metadata?.name + '.yaml'
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
      title: 'Status',
      dataIndex: 'status',
      render: RenderStatus,
    },

    {
      title: 'Actions',
      dataIndex: 'actions',
      render: (group: UpstreamGroup.AsObject) => (
        <Tooltip title='Download'>
          <TableActions>
            <TableActionCircle onClick={() => onDownloadUpstream(group)}>
              <DownloadIcon />
            </TableActionCircle>
          </TableActions>
        </Tooltip>
      ),
    },
  ];

  if (props.wholePage) {
    columns.splice(
      2,
      0,
      {
        title: 'Gloo Instance',
        dataIndex: 'glooInstance',
        render: renderGlooInstanceList,
      },
      {
        title: 'Cluster',
        dataIndex: 'cluster',
        render: RenderCluster,
      }
    );
  }

  return (
    <TableHolder wholePage={props.wholePage}>
      <SoloTable
        loading={ugResponse === undefined}
        pagination={{
          total,
          pageSize: limit,
          onShowSizeChange: (_page, size) => {
            setLimit(size);
            setPage(1);
          },
          current: page,
          onChange: newPage => setPage(newPage),
        }}
        columns={columns}
        dataSource={tableData}
        removeShadows
        curved={props.wholePage}
      />
    </TableHolder>
  );
};

export const UpstreamGroupsPageTable = (props: Props) => {
  return (
    <SectionCard
      cardName={'Upstream Groups'}
      logoIcon={
        <UpstreamGroupIconHolder>
          <UpstreamGroupIcon />
        </UpstreamGroupIconHolder>
      }
      noPadding={true}>
      <UpstreamGroupsTable {...props} wholePage={true} />
    </SectionCard>
  );
};
