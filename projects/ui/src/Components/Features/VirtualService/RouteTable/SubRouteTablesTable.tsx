import React, { useEffect } from 'react';
import styled from '@emotion/styled/macro';
import { SoloTable } from 'Components/Common/SoloTable';
import { SectionCard } from 'Components/Common/SectionCard';
import { ReactComponent as KubeIcon } from 'assets/kubernetes-icon.svg';
import { ReactComponent as RouteIcon } from 'assets/route-icon.svg';
import { ReactComponent as ArrowToggle } from 'assets/arrow-toggle.svg';
import { colors } from 'Styles/colors';
import { useParams } from 'react-router';
import { useGetSubroutesForVirtualService } from 'API/hooks';
import { SubRouteTableRow } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/rt_selector_pb';
import { IconHolder } from 'Styles/StyledComponents/icons';
import { DataError } from 'Components/Common/DataError';
import { Loading } from 'Components/Common/Loading';
import { Destination } from 'proto/github.com/solo-io/solo-apis/api/gloo/gloo/v1/proxy_pb';

const RouteIconHolder = styled.div`
  display: flex;
  align-items: center;
  justify-items: start;
  width: 100%;

  svg {
    width: 28px;
  }
`;

// For inside table, blue back and white icon
const InverseRouteIconHolder = styled.div`
  display: flex;
  align-items: center;
  justify-items: start;
  width: 18px;
  height: 18px;
  background: ${colors.seaBlue};
  border-radius: 100%;

  svg {
    width: 15px;
    margin-left: -1px;
    transform: scaleX(1);

    * {
      fill: white;
    }
  }
`;

const RouteTableContainer = styled.div`
  table {
    thead.ant-table-thead tr th {
      background: ${colors.marchGrey};
    }

    .ant-table-tbody .ant-table-row > td.ant-table-row-expand-icon-cell,
    thead.ant-table-thead tr > th.ant-table-row-expand-icon-cell {
      position: relative;
      width: 10px;
      padding: 0;
      padding-left: 10px;

      button {
        position: absolute;
        right: -10px;
        top: 25px;
      }
    }

    .ant-table-row-expand-icon-spaced {
      background: transparent;
      border: 0;
      visibility: hidden;
      pointer-events: none;
    }

    .ant-table-tbody .ant-table-row > td:nth-of-type(2),
    thead.ant-table-thead tr > th:nth-of-type(2) {
      width: 250px;
    }

    .ant-table-tbody .ant-table-row > td:nth-of-type(3),
    thead.ant-table-thead tr > th:nth-of-type(3) {
      width: 225px;
    }

    .ant-table-tbody .ant-table-row > td:nth-of-type(4),
    thead.ant-table-thead tr > th:nth-of-type(4) {
      width: 200px;
    }

    .ant-table-tbody .ant-table-row > td:nth-of-type(5),
    thead.ant-table-thead tr > th:nth-of-type(5) {
      width: 150px;
    }
    .ant-table-tbody .ant-table-row > td:nth-of-type(6),
    thead.ant-table-thead tr > th:nth-of-type(6) {
      width: 150px;
    }
    .ant-table-tbody .ant-table-row > td:nth-of-type(7),
    thead.ant-table-thead tr > th:nth-of-type(7) {
      width: 150px;
    }
  }
`;

type ToggleHolderProps = {
  open: boolean;
};
const ToggleHolder = styled.div<ToggleHolderProps>`
  cursor: pointer;
  min-width: 14px;

  svg.toggle {
    width: 9px;
    margin-left: 5px;
    transition: transform 0.5s;
    transform: rotate(
      ${(props: ToggleHolderProps) => (props.open ? '0' : '180')}deg
    );
  }
`;

type DestinationHolderProps = {
  depth: number;
};

// FIXME: We duplicate some styles between the RoutesTable.tsx and SubRouteTablesTable.tsx

const DestinationHolder = styled.div<DestinationHolderProps>`
  position: relative;
  display: flex;

  ${(props: DestinationHolderProps) => `
    padding-left: ${props.depth ? (props.depth - 1) * 25 + 25 : 0}px;

    ${
      props.depth > 0
        ? `
          &:before {
            content: '';
            position: absolute;
            border-left: 1px dotted ${colors.marchGrey};
            border-bottom: 1px dotted ${colors.marchGrey};
            left: ${(props.depth - 1) * 25 + 8}px;
            top: -20px;
            width: 8px;
            height: 29px;
          }
          &:after {
            content: '';
            position: absolute;
            left: ${(props.depth - 1) * 25 + 16}px;
            top: 6px;
            width: 5px;
            height: 5px;
            border-radius: 100%;
            border: 1px solid ${colors.marchGrey};
          }
          `
        : ''
    }
  }
  `}

  span:not(.ant-checkbox-inner) {
    margin-right: 5px;
  }

  svg.toggle {
    height: 15px;
  }
`;

const getSingleDestinationName = (dest: Destination.AsObject): string => {
  return (
    dest.consul?.serviceName ??
    dest.kube?.ref?.name ??
    dest.upstream?.name ??
    ''
  );
};

const getRouteDestinations = (rt: SubRouteTableRow.AsObject): string => {
  if (!!rt.routeAction?.single) {
    return getSingleDestinationName(rt.routeAction.single);
  } else if (rt.routeAction?.multi?.destinationsList) {
    return rt.routeAction.multi.destinationsList
      .map(dest =>
        dest.destination ? getSingleDestinationName(dest.destination) : ''
      )
      .join(', ');
  }

  return (
    rt.routeAction?.upstreamGroup?.name ??
    rt.delegateAction?.ref?.name ??
    rt.redirectAction?.pathRedirect ??
    ''
  );
};

type RouteTableFields = {
  key: string;
  destination: {
    destination: string;
    depth: number;
  };
  matcher: string;
  matchType: string;
  methods: string;
  headers: string;
  queryParams: string;
  subrouteRows?: React.ReactNode;
};

const renderDestinationCell = (cellInfo: {
  destination: string;
  depth: number;
}) => {
  return (
    <DestinationHolder depth={cellInfo.depth}>
      <span>
        {cellInfo.depth > 0 ? (
          <InverseRouteIconHolder>
            <RouteIcon />
          </InverseRouteIconHolder>
        ) : (
          <IconHolder width={18}>
            <KubeIcon />
          </IconHolder>
        )}
      </span>
      {cellInfo.destination}{' '}
      {cellInfo.depth > 0 ? (
        <ArrowToggle className='toggle' />
      ) : (
        <React.Fragment />
      )}
    </DestinationHolder>
  );
};

const HeaderlessSubrouteTable = ({
  subrouteRows,
  depth,
}: {
  subrouteRows: SubRouteTableRow.AsObject[];
  depth: number;
}) => {
  const columns = [
    {
      title: 'Destination',
      dataIndex: 'destination',
      width: 200,
      render: renderDestinationCell,
    },
    {
      title: 'Matcher',
      dataIndex: 'matcher',
    },
    {
      title: 'Path Match Type',
      dataIndex: 'matchType',
    },
    {
      title: 'Methods',
      dataIndex: 'methods',
    },
    {
      title: 'Headers',
      dataIndex: 'headers',
    },
    {
      title: 'Query Parameters',
      dataIndex: 'queryParams',
    },
  ];

  return (
    <RouteTableContainer>
      <SoloTable
        columns={columns}
        dataSource={subrouteRows.map(subRt => {
          return {
            key: getRouteDestinations(subRt) + subRt.matcher,
            destination: {
              destination: getRouteDestinations(subRt),
              depth,
            },
            matcher: subRt.matcher,
            matchType: subRt.matchType,
            methods: subRt.methodsList.join(',  '),
            headers: subRt.headersList
              .map(header => `${header.name}:${header.value}`)
              .join(', '),
            queryParams: subRt.queryParametersList
              .map(query => `${query.name}:${query.value}`)
              .join(', '),
            subrouteRtRows: subRt.rtRoutesList,
          };
        })}
        removePaging
        removeShadows
        expandable={{
          expandedRowRender: row => (
            <HeaderlessSubrouteTable
              subrouteRows={row.subrouteRtRows ?? []}
              depth={depth + 1}
            />
          ),
          expandIcon: ({ expanded, onExpand, record }) => (
            <ToggleHolder open={expanded} onClick={e => onExpand(record, e)}>
              {!!record.subrouteRtRows.length ? (
                <ArrowToggle className='toggle' />
              ) : null}
            </ToggleHolder>
          ),
        }}
        hideHeader={true}
      />
    </RouteTableContainer>
  );
};

/* FULL COMPONENT PIECES */

type Props = {
  statusFilter?: 'accepted' | 'rejected';
  nameFilter?: string;
  onlyTable?: boolean;
};
// This is almost a dupe of VSTable, but as we are in early stages keeping it decoupled for now.
export const SubRouteTablesTable = (props: Props) => {
  const {
    virtualservicename = '',
    virtualservicenamespace = '',
    virtualserviceclustername = '',
  } = useParams();

  const [tableData, setTableData] = React.useState<RouteTableFields[]>([]);

  const { data: routeTables, error: rtError } =
    useGetSubroutesForVirtualService(
      !!virtualservicenamespace && !!virtualservicename
        ? {
            clusterName: virtualserviceclustername,
            name: virtualservicename,
            namespace: virtualservicenamespace,
          }
        : undefined
    );

  useEffect(() => {
    if (routeTables) {
      let newTableData: RouteTableFields[] = routeTables.map(rt => {
        return {
          key: getRouteDestinations(rt) + rt.matcher,
          destination: {
            destination: getRouteDestinations(rt),
            matcher: rt.matcher,
            depth: 0,
          },
          matcher: rt.matcher,
          matchType: rt.matchType,
          methods: rt.methodsList.join(',  '),
          headers: rt.headersList
            .map(header => `${header.name}:${header.value}`)
            .join(', '),
          queryParams: rt.queryParametersList
            .map(query => `${query.name}:${query.value}`)
            .join(', '),
          subrouteRows: rt.rtRoutesList,
        };
      });

      setTableData(newTableData);
    } else {
      setTableData([]);
    }
  }, [routeTables, props.nameFilter, props.statusFilter]);

  if (!!rtError) {
    return <DataError error={rtError} />;
  } else if (!routeTables) {
    return (
      <Loading
        message={`Retrieving route tables for ${virtualservicename}...`}
      />
    );
  }

  const columns = [
    {
      title: 'Destination',
      dataIndex: 'destination',
      width: 200,
      render: renderDestinationCell,
    },
    {
      title: 'Matcher',
      dataIndex: 'matcher',
    },
    {
      title: 'Path Match Type',
      dataIndex: 'matchType',
    },
    {
      title: 'Methods',
      dataIndex: 'methods',
    },
    {
      title: 'Headers',
      dataIndex: 'headers',
    },
    {
      title: 'Query Parameters',
      dataIndex: 'queryParams',
    },
  ];

  return (
    <RouteTableContainer>
      {props.onlyTable ? (
        <SoloTable
          columns={columns}
          dataSource={tableData}
          removePaging
          removeShadows
          curved
          expandable={{
            expandedRowRender: row => (
              <HeaderlessSubrouteTable
                subrouteRows={row.subrouteRows}
                depth={1}
              />
            ),
            expandIcon: ({ expanded, onExpand, record }) => (
              <ToggleHolder open={expanded} onClick={e => onExpand(record, e)}>
                {!!record.subrouteRows.length ? (
                  <ArrowToggle className='toggle' />
                ) : null}
              </ToggleHolder>
            ),
          }}
        />
      ) : (
        <SectionCard
          cardName={'Route Tables'}
          logoIcon={
            <RouteIconHolder>
              <RouteIcon />
            </RouteIconHolder>
          }
          noPadding={true}
        >
          <SoloTable
            columns={columns}
            dataSource={tableData}
            removePaging
            removeShadows
            curved
            expandable={{
              expandedRowRender: row => (
                <HeaderlessSubrouteTable
                  subrouteRows={row.subrouteRows}
                  depth={1}
                />
              ),
              expandIcon: ({ expanded, onExpand, record }) => (
                <ToggleHolder
                  open={expanded}
                  onClick={e => onExpand(record, e)}
                >
                  {!!record.subrouteRows.length ? (
                    <ArrowToggle className='toggle' />
                  ) : null}
                </ToggleHolder>
              ),
            }}
          />
        </SectionCard>
      )}
    </RouteTableContainer>
  );
};
