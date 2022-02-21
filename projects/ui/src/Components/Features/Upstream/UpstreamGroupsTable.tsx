import React, { useEffect } from 'react';
import styled from '@emotion/styled/macro';
import {
  SoloTable,
  RenderStatus,
  TableActionCircle,
  TableActions,
  RenderCluster,
} from 'Components/Common/SoloTable';
import { SectionCard } from 'Components/Common/SectionCard';
import { ReactComponent as UpstreamGroupIcon } from 'assets/upstream-group-icon.svg';
import { ReactComponent as DownloadIcon } from 'assets/download-icon.svg';
import { colors } from 'Styles/colors';
import { useParams, useNavigate } from 'react-router';
import { useListUpstreamGroups, useIsGlooFedEnabled } from 'API/hooks';
import { Loading } from 'Components/Common/Loading';
import { objectMetasAreEqual } from 'API/helpers';
import { UpstreamGroup } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/gloo_resources_pb';
import { UpstreamStatus } from 'proto/github.com/solo-io/solo-apis/api/gloo/gloo/v1/upstream_pb';
import { SimpleLinkProps, RenderSimpleLink } from 'Components/Common/SoloLink';
import { glooResourceApi } from 'API/gloo-resource';
import { doDownload } from 'download-helper';
import { DataError } from 'Components/Common/DataError';

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
  glooInstanceFilter?: {
    name: string;
    namespace: string;
  };
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

  const { data: upstreamGroups, error: upstreamGroupsError } =
    useListUpstreamGroups(
      !!name && !!namespace
        ? {
            name,
            namespace,
          }
        : undefined
    );

  const { data: glooFedCheckResponse, error: glooFedCheckError } =
    useIsGlooFedEnabled();
  const isGlooFedEnabled = glooFedCheckResponse?.enabled;

  useEffect(() => {
    if (upstreamGroups) {
      setTableData(
        upstreamGroups
          .filter(
            group =>
              group.metadata?.name.includes(props.nameFilter ?? '') &&
              (props.statusFilter === undefined ||
                group.status?.state === props.statusFilter) &&
              (!props.glooInstanceFilter ||
                objectMetasAreEqual(
                  {
                    name: group.glooInstance?.name ?? '',
                    namespace: group.glooInstance?.namespace ?? '',
                  },
                  props.glooInstanceFilter
                ))
          )
          .sort(
            (gA, gB) =>
              (gA.metadata?.name ?? '').localeCompare(
                gB.metadata?.name ?? ''
              ) ||
              (!props.wholePage
                ? 0
                : (gA.glooInstance?.name ?? '').localeCompare(
                    gB.glooInstance?.name ?? ''
                  ))
          )
          .map(group => {
            return {
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
            };
          })
      );
    } else {
      setTableData([]);
    }
    /* eslint-disable-next-line react-hooks/exhaustive-deps */
  }, [
    upstreamGroups,
    props.nameFilter,
    props.statusFilter,
    props.glooInstanceFilter,
  ]);

  if (!!upstreamGroupsError) {
    return <DataError error={upstreamGroupsError} />;
  } else if (!upstreamGroups) {
    return <Loading message={'Retrieving upstream groups...'} />;
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
        }
      >
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
        <TableActions>
          <TableActionCircle onClick={() => onDownloadUpstream(group)}>
            <DownloadIcon />
          </TableActionCircle>
        </TableActions>
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
        columns={columns}
        dataSource={tableData}
        removePaging
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
      noPadding={true}
    >
      <UpstreamGroupsTable {...props} wholePage={true} />
    </SectionCard>
  );
};
