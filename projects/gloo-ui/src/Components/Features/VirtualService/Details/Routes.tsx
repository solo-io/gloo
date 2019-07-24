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
import { TableActionCircle, TableActions } from 'Styles';
import { SoloModal } from 'Components/Common/SoloModal';
import { NewRouteRowForm } from 'Components/Features/VirtualService/Details/NewRouteRowForm';
import { CreateRouteModal } from 'Components/Features/Route/CreateRouteModal';

const RouteMatch = styled.div`
  max-width: 200px;
  max-height: 70px;
  overflow: hidden;
  text-overflow: ellipsis;
`;

const getRouteColumns = (
  showCreateRouteModal: (boolean: true) => void,
  deleteRoute: (matcher: string) => any
) => {
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
      dataIndex: 'method',
      width: 150
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
      render: (matcher: string) => {
        return (
          <TableActions>
            {/* Unclear what edit should look like yet
            <TableActionCircle onClick={() => showCreateRouteModal(true)}>
              +
            </TableActionCircle>

            <div style={{ marginLeft: '5px' }}>*/}
            <div>
              <TableActionCircle onClick={() => deleteRoute(matcher)}>
                x
              </TableActionCircle>
            </div>
          </TableActions>
        );
      }
    }
  ];
};

interface Props {
  routes: Route.AsObject[];
  virtualService: VirtualService.AsObject;
  routesChanged: (newRoutes: Route.AsObject[]) => any;
  reloadVirtualService: (newVirtualService?: VirtualService.AsObject) => any;
}

export const Routes: React.FC<Props> = props => {
  const [editRoute, setEditRoute] = React.useState<boolean>(false);

  const getRouteData = (routes: Route.AsObject[]) => {
    const existingRoutes = routes.map(route => {
      const upstreamName = getRouteSingleUpstream(route).name || '';
      const { matcher, matchType } = getRouteMatcher(route);
      return {
        key: matcher,
        matcher: matcher,
        pathMatch: matchType,
        method: getRouteMethods(route),
        upstreamName: upstreamName,
        header: getRouteHeaders(route),
        queryParams: getRouteQueryParams(route),
        actions: matcher
      };
    });

    existingRoutes.push({
      key: 'creating-additional-route',
      matcher: '',
      pathMatch: '',
      method: '',
      upstreamName: '',
      header: '',
      queryParams: '',
      actions: ''
    });

    return existingRoutes;
  };

  const deleteRoute = (matcherToDelete: string) => {
    props.routesChanged(
      props.routes.filter(
        route => getRouteMatcher(route).matcher !== matcherToDelete
      )
    );
  };

  const finishRouteEditiing = (newVirtualService?: VirtualService.AsObject) => {
    setEditRoute(false);
    props.reloadVirtualService();
  };

  return (
    <React.Fragment>
      <DetailsSectionTitle>Routes</DetailsSectionTitle>
      <SoloTable
        columns={getRouteColumns(setEditRoute, deleteRoute)}
        dataSource={getRouteData(props.routes)}
        formComponent={() => (
          <NewRouteRowForm
            virtualService={props.virtualService}
            reloadVirtualService={props.reloadVirtualService}
          />
        )}
      />
    </React.Fragment>
  );
};
