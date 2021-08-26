import React, { useEffect } from 'react';
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
import { useParams, useNavigate } from 'react-router';
import { useListRouteTables, useIsGlooFedEnabled } from 'API/hooks';
import { RouteTable } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/gateway_resources_pb';
import { Loading } from 'Components/Common/Loading';
import { objectMetasAreEqual } from 'API/helpers';
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

const renderGlooInstanceList = (glooInstance: {
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
} & TableHolderProps;

export const RouteTablesTable = ({ routeTables, wholePage }: TableProps) => {
  const [tableData, setTableData] = React.useState<RouteTableTableFields[]>([]);

  const {
    data: glooFedCheckResponse,
    error: glooFedCheckError,
  } = useIsGlooFedEnabled();
  const isGlooFedEnabled = glooFedCheckResponse?.enabled;

  useEffect(() => {
    setTableData(
      routeTables.map(gwRoute => fillRowFields(gwRoute, isGlooFedEnabled))
    );
  }, [routeTables]);

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
            render: renderGlooInstanceList,
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
        <TableActions>
          <TableActionCircle onClick={() => onDownloadRouteTable(rt)}>
            <DownloadIcon />
          </TableActionCircle>
        </TableActions>
      ),
    },
  ];

  return (
    <RoutesTableHolder>
      <SoloTable
        columns={columns}
        dataSource={tableData}
        removePaging
        removeShadows
        curved={true}
      />
    </RoutesTableHolder>
  );
};

type Props = {
  statusFilter?: RouteTableStatus.StateMap[keyof RouteTableStatus.StateMap];
  nameFilter?: string;
  glooInstanceFilter?: {
    name: string;
    namespace: string;
  };
} & TableHolderProps;
export const RouteTablesPageCardContents = (props: Props) => {
  const { name, namespace } = useParams();

  const [tableData, setTableData] = React.useState<RouteTable.AsObject[]>([]);

  const { data: routeTables, error: routeTablesError } = useListRouteTables(
    !!name && !!namespace ? { name, namespace } : undefined
  );

  useEffect(() => {
    if (routeTables) {
      setTableData(
        routeTables
          .filter(
            vs =>
              vs.metadata?.name.includes(props.nameFilter ?? '') &&
              (props.statusFilter === undefined ||
                vs.status?.state === props.statusFilter) &&
              (!props.glooInstanceFilter ||
                objectMetasAreEqual(
                  {
                    name: vs.glooInstance?.name ?? '',
                    namespace: vs.glooInstance?.namespace ?? '',
                  },
                  props.glooInstanceFilter
                ))
          )
          .sort(
            (gA, gB) =>
              (gA.metadata?.name ?? '').localeCompare(
                gB.metadata?.name ?? ''
              ) ||
              (gA.glooInstance?.name ?? '').localeCompare(
                gB.glooInstance?.name ?? ''
              )
          )
      );
    } else {
      setTableData([]);
    }
  }, [
    routeTables,
    props.nameFilter,
    props.statusFilter,
    props.glooInstanceFilter,
  ]);

  if (!!routeTablesError) {
    return <DataError error={routeTablesError} />;
  } else if (!routeTables) {
    return <Loading message={'Retrieving route tables...'} />;
  }

  return (
    <TableHolder wholePage={props.wholePage}>
      <RouteTablesTable routeTables={tableData} wholePage={props.wholePage} />
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
