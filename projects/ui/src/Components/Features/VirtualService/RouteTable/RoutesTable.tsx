import styled from '@emotion/styled/macro';
import { ColumnsType } from 'antd/lib/table';
import { ReactComponent as ArrowToggle } from 'assets/arrow-toggle.svg';
import { ReactComponent as KubeIcon } from 'assets/kubernetes-icon.svg';
import { ReactComponent as RouteIcon } from 'assets/route-icon.svg';
import { SoloTable } from 'Components/Common/SoloTable';
import { Route } from 'proto/github.com/solo-io/solo-apis/api/gloo/gateway/v1/virtual_service_pb';
import {
  HeaderMatcher,
  QueryParameterMatcher,
} from 'proto/github.com/solo-io/solo-apis/api/gloo/gloo/v1/core/matchers/matchers_pb';
import { Destination } from 'proto/github.com/solo-io/solo-apis/api/gloo/gloo/v1/proxy_pb';
import React, { useEffect } from 'react';
import { colors } from 'Styles/colors';
import { IconHolder } from 'Styles/StyledComponents/icons';
import {
  getRouteHeaders,
  getRouteMatcher,
  getRouteMethods,
  getRouteQueryParams,
  getRouteSingleUpstream,
} from 'utils/route-helpers';

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

  span {
    margin-right: 5px;
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

const getRouteDestinations = (rt: Route.AsObject): string => {
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

interface RouteTableFields {
  key: string;
  destination: {
    destination: string;
    depth: number;
  };
  matcher: string;
  matchType: string;
  methods: string;
  headers: HeaderMatcher.AsObject[];
  queryParams: QueryParameterMatcher.AsObject[];
  children?: RouteTableFields[];
}

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

function fillRowFields(rt: Route.AsObject): RouteTableFields {
  const upstreamName = getRouteSingleUpstream(rt) || '';
  const { matcher, matchType } = getRouteMatcher(rt);

  return {
    key: `${matcher}-${upstreamName}`,
    destination: {
      destination: getRouteDestinations(rt),
      depth: 0,
    },
    matcher: matcher,
    matchType: matchType,
    methods: getRouteMethods(rt),
    headers: getRouteHeaders(rt),
    queryParams: getRouteQueryParams(rt),
  };
}

const columns = [
  {
    title: 'Destination',
    dataIndex: 'destination',
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
    render: (value: any, record: RouteTableFields, index: number) => (
      <div key={record.key}>
        {record.headers.map(header => (
          <div key={`${header?.name ?? index}-${header?.value ?? ''}`}>
            {header.name} : {header.value}
          </div>
        ))}
      </div>
    ),
  },
  {
    title: 'Query Parameters',
    dataIndex: 'queryParams',
    render: (value: any, record: RouteTableFields, index: number) => (
      <div key={record.key}>
        {record.queryParams.map(queryParam => (
          <div key={`${queryParam?.name ?? index}-${queryParam?.value ?? ''}`}>
            {queryParam.name} : {queryParam.value}
          </div>
        ))}
      </div>
    ),
  },
] as ColumnsType<RouteTableFields>;

/* FULL COMPONENT PIECES */

type Props = {
  routes: Route.AsObject[];
};
export const RoutesTable = ({ routes }: Props) => {
  const [tableData, setTableData] = React.useState<RouteTableFields[]>([]);

  useEffect(() => {
    setTableData(routes.map(rt => fillRowFields(rt)));
  }, [routes]);

  return (
    <RouteTableContainer>
      <SoloTable
        columns={columns}
        dataSource={tableData}
        removePaging
        removeShadows
        curved
      />
    </RouteTableContainer>
  );
};
