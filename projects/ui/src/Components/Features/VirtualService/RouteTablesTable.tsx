import React, { useEffect, useState } from 'react';
import styled from '@emotion/styled/macro';
import { colors } from 'Styles/colors';
import {
  SoloTable,
  RenderStatus,
  TableActionCircle,
  TableActions,
  RenderCluster,
} from 'Components/Common/SoloTable';
import { SectionCard } from 'Components/Common/SectionCard';
import { ReactComponent as RouteTableIcon } from 'assets/route-icon.svg';
import { ReactComponent as DownloadIcon } from 'assets/download-icon.svg';
import Tooltip from 'antd/lib/tooltip';
import { useParams, useNavigate } from 'react-router';
import { useListRouteTables, useIsGlooFedEnabled } from 'API/hooks';
import { RouteTable } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/gateway_resources_pb';
import { Loading } from 'Components/Common/Loading';
import { SimpleLinkProps, RenderSimpleLink } from 'Components/Common/SoloLink';
import { gatewayResourceApi } from 'API/gateway-resources';
import { doDownload } from 'download-helper';
import { DataError } from 'Components/Common/DataError';
import { RouteTableStatus } from 'proto/github.com/solo-io/solo-apis/api/gloo/gateway/v1/route_table_pb';
import { ObjectRef } from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';

const RoutesTableHolder = styled.div`
  position: relative;
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

type RouteTableTableFields = {
  key: string;
  name: SimpleLinkProps;
  namespace: string;
  version: string;
  status?: RouteTableStatus.StateMap[keyof RouteTableStatus.StateMap];
  routes: number;

  glooInstance?: ObjectRef.AsObject;
  cluster: string;
  actions: RouteTable.AsObject;
};

function fillRowFields(
  rt: RouteTable.AsObject,
  isGlooFedEnabled?: boolean
): RouteTableTableFields {
  return {
    key: rt.metadata?.name ?? '',
    name: {
      displayElement: rt.metadata?.name ?? '',
      link: rt.metadata
        ? isGlooFedEnabled
          ? `/gloo-instances/${rt.glooInstance?.namespace}/${rt.glooInstance?.name}/route-tables/${rt.metadata.clusterName}/${rt.metadata.namespace}/${rt.metadata.name}/`
          : `/gloo-instances/${rt.glooInstance?.namespace}/${rt.glooInstance?.name}/route-tables/${rt.metadata.namespace}/${rt.metadata.name}/`
        : '',
    },
    namespace: rt.metadata?.namespace ?? '',
    version: rt.metadata?.resourceVersion ?? '',
    status: rt.status?.state,
    routes: rt.spec?.routesList.length ?? 0,
    glooInstance: rt.glooInstance,
    cluster: rt.metadata?.clusterName ?? '',
    actions: rt,
  };
}

const RenderGlooInstanceList = (glooInstance: {
  name: string;
  namespace: string;
}) => {
  const navigate = useNavigate();
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

const onDownloadRouteTable = (rt: RouteTable.AsObject) => {
  if (rt.metadata) {
    gatewayResourceApi
      .getRouteTableYAML({
        name: rt.metadata.name!,
        namespace: rt.metadata.namespace!,
        clusterName: rt.metadata.clusterName,
      })
      .then(yaml => {
        doDownload(
          yaml,
          `${rt.metadata?.namespace}--${rt.metadata?.name}.yaml`
        );
      });
  }
};

type TableProps = {
  routeTables: RouteTable.AsObject[];
  page: number;
  setPage(newPage: number): void;
  total: number;
  limit: number;
  setOffset: React.Dispatch<React.SetStateAction<number>>;
} & TableHolderProps;

export const RouteTablesTable = ({
  page,
  setPage,
  limit,
  setOffset,
  total,
  routeTables,
  wholePage,
}: TableProps) => {
  const [tableData, setTableData] = React.useState<RouteTableTableFields[]>([]);

  const { data: glooFedCheckResponse, error: glooFedCheckError } =
    useIsGlooFedEnabled();
  const isGlooFedEnabled = glooFedCheckResponse?.enabled;

  useEffect(() => {
    setTableData(
      routeTables.map(gwRoute => fillRowFields(gwRoute, isGlooFedEnabled))
    );
  }, [routeTables, isGlooFedEnabled]);

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

    ...(wholePage
      ? [
          {
            title: 'Gloo Instance',
            dataIndex: 'glooInstance',
            render: RenderGlooInstanceList,
          },
          {
            title: 'Cluster',
            dataIndex: 'cluster',
            render: RenderCluster,
          },
        ]
      : [{}]),

    {
      title: 'Version',
      dataIndex: 'version',
    },
    {
      title: 'Routes',
      dataIndex: 'routes',
    },
    {
      title: 'Status',
      dataIndex: 'status',
      render: RenderStatus,
    },

    {
      title: 'Actions',
      dataIndex: 'actions',
      render: (rt: RouteTable.AsObject) => (
        <Tooltip title='Download'>
          <TableActions>
            <TableActionCircle onClick={() => onDownloadRouteTable(rt)}>
              <DownloadIcon />
            </TableActionCircle>
          </TableActions>
        </Tooltip>
      ),
    },
  ];

  return (
    <RoutesTableHolder>
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
        curved={false}
      />
    </RoutesTableHolder>
  );
};

type Props = {
  statusFilter?: RouteTableStatus.StateMap[keyof RouteTableStatus.StateMap];
  nameFilter?: string;
} & TableHolderProps;
export const RouteTablesPageCardContents = (props: Props) => {
  const [tableData, setTableData] = React.useState<RouteTable.AsObject[]>([]);
  const [offset, setOffset] = React.useState(0);
  const [limit, _] = React.useState(10);
  const { name, namespace } = useParams();
  const { data: routeTablesResponse, error: routeTablesResponseError } =
    useListRouteTables(
      { name, namespace },
      { limit, offset },
      props.nameFilter,
      props.statusFilter
    );

  const routeTables = routeTablesResponse?.routeTablesList ?? [];
  const total = routeTablesResponse?.total ?? 0;

  const [page, setPage] = useState(1);
  useEffect(() => {
    setPage(1);
  }, [props.nameFilter, props.statusFilter]);
  useEffect(() => {
    setOffset(limit * (page - 1));
  }, [page]);

  useEffect(() => {
    if (routeTables) {
      setTableData(routeTables);
    } else {
      setTableData([]);
    }
  }, [routeTables, props.nameFilter, props.statusFilter]);

  if (!!routeTablesResponseError) {
    return <DataError error={routeTablesResponseError} />;
  }

  return (
    <TableHolder wholePage={props.wholePage}>
      <RouteTablesTable
        page={page}
        setPage={setPage}
        limit={limit}
        setOffset={setOffset}
        total={total}
        routeTables={tableData}
        wholePage={props.wholePage}
      />
    </TableHolder>
  );
};

const GlooIconHolder = styled.div`
  display: flex;
  align-items: center;
  justify-items: center;

  svg {
    width: 25px;
  }
`;

export const RouteTablesPageTable = (props: Props) => {
  return (
    <SectionCard
      cardName={'Route Tables'}
      logoIcon={
        <GlooIconHolder>
          <RouteTableIcon />
        </GlooIconHolder>
      }
      noPadding={true}>
      <RouteTablesPageCardContents {...props} wholePage={true} />
    </SectionCard>
  );
};
