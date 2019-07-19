import * as React from 'react';
import styled from '@emotion/styled/macro';

import { SoloTable } from 'Components/Common/SoloTable';
import { DetailsSectionTitle } from './VirtualServiceDetails';
import { Route } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/proxy_pb';
import {
  getRouteMethods,
  getRouteSingleUpstream,
  getRouteMatcher,
  getRouteHeaders,
  getRouteQueryParams
} from 'utils/helpers';
import { VirtualService } from 'proto/github.com/solo-io/gloo/projects/gateway/api/v1/virtual_service_pb';
import { TableActionCircle } from 'Styles';
import { SoloModal } from 'Components/Common/SoloModal';
import { CreateRouteModal } from 'Components/Features/Route/CreateRouteModal';

const RouteMatch = styled.div`
  max-width: 200px;
  max-height: 70px;
  overflow: hidden;
  text-overflow: ellipsis;
`;

const getRouteColumns = (showCreateRouteModal: (boolean: true) => any) => {
  return [
    {
      title: 'Matcher',
      dataIndex: 'matcher',
      render: (matcher: string) => {
        return <RouteMatch>{matcher}</RouteMatch>;
      }
    },
    {
      title: 'Path Match Type',
      dataIndex: 'pathMatch'
    },
    {
      title: 'Methods',
      dataIndex: 'method'
    },
    {
      title: 'Destination',
      dataIndex: 'upstreamName'
    },
    {
      title: 'Headers',
      dataIndex: 'header'
    },
    {
      title: 'Query Parameters',
      dataIndex: 'queryParams'
    },
    {
      title: 'Actions',
      dataIndex: 'actions',
      render: () => {
        return (
          <TableActionCircle onClick={() => showCreateRouteModal(true)}>
            +
          </TableActionCircle>
        );
      }
    }
  ];
};

interface Props {
  routes: Route.AsObject[];
  virtualService: VirtualService.AsObject;
}

export const Routes: React.FC<Props> = props => {
  const [createRoute, setCreateRoute] = React.useState<boolean>(false);

  const getRouteData = (routes: Route.AsObject[]) => {
    return routes.map(route => {
      const upstreamName = getRouteSingleUpstream(route).name || '';
      const { matcher, matchType } = getRouteMatcher(route);
      console.log(route);
      return {
        key: route.matcher!.prefix,
        matcher: matcher,
        pathMatch: matchType,
        method: getRouteMethods(route),
        upstreamName: upstreamName,
        header: getRouteHeaders(route),
        queryParams: getRouteQueryParams(route),
        actions: ''
      };
    });
  };

  return (
    <React.Fragment>
      <DetailsSectionTitle>Routes</DetailsSectionTitle>
      <SoloTable
        columns={getRouteColumns(setCreateRoute)}
        dataSource={getRouteData(props.routes)}
      />
      <SoloModal
        visible={createRoute}
        width={500}
        title={'Create Route'}
        onClose={() => setCreateRoute(false)}>
        <CreateRouteModal defaultVirtualService={props.virtualService} />
      </SoloModal>
    </React.Fragment>
  );
};
