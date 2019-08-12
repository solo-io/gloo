import * as React from 'react';
import styled from '@emotion/styled/macro';

import { SoloDragSortableTable } from 'Components/Common/SoloDragSortableTable';
import { Route } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/proxy_pb';
import {
  getRouteMethods,
  getRouteSingleUpstream,
  getRouteMatcher,
  getRouteHeaders,
  getRouteQueryParams
} from 'utils/helpers';
import { VirtualService } from 'proto/github.com/solo-io/gloo/projects/gateway/api/v1/virtual_service_pb';
import { TableActionCircle, TableActions, colors } from 'Styles';
import { SoloModal } from 'Components/Common/SoloModal';
import { ReactComponent as GreenPlus } from 'assets/small-green-plus.svg';
import { ReactComponent as EditPencil } from 'assets/edit-pencil.svg';
import { CreateRouteModal } from 'Components/Features/Route/CreateRouteModal';
import { Popconfirm } from 'antd';

const RouteMatch = styled.div`
  max-width: 200px;
  max-height: 70px;
  overflow: hidden;
  text-overflow: ellipsis;
`;

const RouteSectionTitle = styled.div`
  font-size: 18px;
  font-weight: bold;
  color: ${colors.novemberGrey};
  margin-top: 10px;
  margin-bottom: 10px;
  display: flex;
  justify-content: space-between;
`;

const StyledGreenPlus = styled(GreenPlus)`
  cursor: pointer;
  margin-right: 7px;
`;
const ModalTrigger = styled.div`
  cursor: pointer;
  display: flex;
  align-items: center;
  padding: 0 10px;
  font-size: 14px;
`;

const getRouteColumns = (
  showEditRouteModal: (matcher: string) => void,
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
            <TableActionCircle onClick={() => showEditRouteModal(matcher)}>
              <EditPencil />
            </TableActionCircle>

            <div style={{ marginLeft: '5px' }}>
              <Popconfirm
                onConfirm={() =>
                  deleteRoute(matcher)
                }
                title={'Are you sure you want to delete this route? '}
                okText='Yes'
                cancelText='No'>
                <TableActionCircle>x</TableActionCircle>
              </Popconfirm>
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
  const [routesList, setRoutesList] = React.useState<Route.AsObject[]>([]);
  const [routeBeingEdited, setRouteBeingEdited] = React.useState<
    Route.AsObject | undefined
  >(undefined);
  const [createNewRoute, setCreateNewRoute] = React.useState<boolean>(false);

  React.useEffect(() => {
    setRoutesList([...props.routes]);
  }, [props.routes]);

  const getRouteData = () => {
    const existingRoutes = routesList.map(route => {
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

    return existingRoutes;
  };

  const deleteRoute = (matcherToDelete: string) => {
    const newList = routesList.filter(
      route => getRouteMatcher(route).matcher !== matcherToDelete
    );

    setRoutesList(newList);
    props.routesChanged(newList);
  };

  const finishNewRouteCreation = () => {
    props.reloadVirtualService();
    setCreateNewRoute(false);
  };

  const beginRouteEditing = (matcherToEdit: string) => {
    setRouteBeingEdited(
      routesList.find(route => getRouteMatcher(route).matcher === matcherToEdit)
    );
  };

  const finishRouteEditiing = () => {
    props.reloadVirtualService();
    setRouteBeingEdited(undefined);
  };

  const reorderRoutes = (dragIndex: number, hoverIndex: number) => {
    const movedRoute = routesList.splice(dragIndex, 1)[0];

    let newRoutesList = [...routesList];
    newRoutesList.splice(hoverIndex, 0, movedRoute);

    setRoutesList(newRoutesList);
    props.routesChanged(newRoutesList);
  };

  return (
    <React.Fragment>
      <RouteSectionTitle>
        Routes
        <ModalTrigger onClick={() => setCreateNewRoute(true)}>
          <React.Fragment>
            <StyledGreenPlus />
            Create Route
          </React.Fragment>
        </ModalTrigger>
      </RouteSectionTitle>

      <SoloDragSortableTable
        columns={getRouteColumns(beginRouteEditing, deleteRoute)}
        dataSource={getRouteData()}
        moveRow={reorderRoutes}
      />

      <SoloModal
        visible={createNewRoute}
        width={500}
        title={'Create Route'}
        onClose={() => setCreateNewRoute(false)}>
        <CreateRouteModal
          defaultVirtualService={props.virtualService}
          completeCreation={finishNewRouteCreation}
          lockVirtualService
        />
      </SoloModal>
      <SoloModal
        visible={!!routeBeingEdited}
        width={500}
        title={'Edit Route'}
        onClose={() => setRouteBeingEdited(undefined)}>
        <CreateRouteModal
          defaultVirtualService={props.virtualService}
          existingRoute={routeBeingEdited}
          completeCreation={finishRouteEditiing}
        />
      </SoloModal>
    </React.Fragment>
  );
};
