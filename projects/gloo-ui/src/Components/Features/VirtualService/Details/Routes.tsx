import styled from '@emotion/styled';
import { Popconfirm } from 'antd';
import { ReactComponent as EditPencil } from 'assets/edit-pencil.svg';
import { ReactComponent as GreenPlus } from 'assets/small-green-plus.svg';
import { SoloDragSortableTable } from 'Components/Common/SoloDragSortableTable';
import { SoloModal } from 'Components/Common/SoloModal';
import { Route } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/proxy_pb';
import * as React from 'react';
import { useDispatch } from 'react-redux';
import { shiftRoutes, deleteRoute } from 'store/virtualServices/actions';
import { colors, TableActionCircle, TableActions } from 'Styles';
import {
  getRouteHeaders,
  getRouteMatcher,
  getRouteMethods,
  getRouteQueryParams,
  getRouteSingleUpstream
} from 'utils/helpers';
import { RouteParent } from '../RouteTableDetails';
import { CreateRouteModal } from '../Creation/CreateRouteModal';

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
      dataIndex: 'destinationName',
      render: (destinationName: string) => <a>{destinationName}</a>
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
                onConfirm={() => deleteRoute(matcher)}
                title={'Are you sure you want to delete this route? '}
                okText='Yes'
                cancelText='No'>
                <TableActionCircle data-testid={`delete-route-${matcher}`}>
                  x
                </TableActionCircle>
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
  routeParent?: RouteParent;
}

export const Routes: React.FC<Props> = props => {
  const [routesList, setRoutesList] = React.useState<Route.AsObject[]>([]);
  const [routeBeingEdited, setRouteBeingEdited] = React.useState<
    Route.AsObject | undefined
  >(undefined);
  const [showCreateRouteModal, setShowCreateRouteModal] = React.useState(false);
  const dispatch = useDispatch();

  let routeParentRef = {
    name: props.routeParent ? props.routeParent.metadata!.name : '',
    namespace: props.routeParent ? props.routeParent.metadata!.namespace : ''
  };
  console.log('routeParentRef', routeParentRef);
  React.useEffect(() => {
    setRoutesList([...props.routes]);
  }, [props.routes]);

  const getRouteData = () => {
    const existingRoutes = routesList.map(route => {
      const upstreamName = getRouteSingleUpstream(route) || '';
      const destinationName = getRouteSingleUpstream(route);
      const { matcher, matchType } = getRouteMatcher(route);
      return {
        key: `${matcher}-${upstreamName}`,
        matcher: matcher,
        pathMatch: matchType,
        method: getRouteMethods(route),
        destinationName: destinationName,
        upstreamName: upstreamName,
        header: getRouteHeaders(route),
        queryParams: getRouteQueryParams(route),
        actions: matcher
      };
    });

    return existingRoutes;
  };

  const handleDeleteRoute = (matcherToDelete: string) => {
    let index = routesList.findIndex(
      route => getRouteMatcher(route).matcher === matcherToDelete
    );
    const newList = routesList.filter(
      route => getRouteMatcher(route).matcher !== matcherToDelete
    );
    console.log('index', index);
    console.log('newList', newList);

    dispatch(
      deleteRoute({
        virtualServiceRef: routeParentRef,
        index
      })
    );
    setRoutesList(newList);
  };

  const beginRouteEditing = (matcherToEdit: string) => {
    setRouteBeingEdited(
      routesList.find(route => getRouteMatcher(route).matcher === matcherToEdit)
    );
  };

  const finishRouteEditiing = () => {
    setRouteBeingEdited(undefined);
  };

  const reorderRoutes = (dragIndex: number, hoverIndex: number) => {
    const movedRoute = routesList.splice(dragIndex, 1)[0];

    let newRoutesList = [...routesList];
    newRoutesList.splice(hoverIndex, 0, movedRoute);

    dispatch(
      shiftRoutes({
        virtualServiceRef: routeParentRef,
        fromIndex: dragIndex,
        toIndex: hoverIndex
      })
    );
    setRoutesList(newRoutesList);
  };

  return (
    <>
      <RouteSectionTitle>
        Routes
        <ModalTrigger
          data-testid='create-new-route-modal'
          onClick={() => setShowCreateRouteModal(true)}>
          <>
            <StyledGreenPlus />
            Create Route
          </>
        </ModalTrigger>
      </RouteSectionTitle>

      <SoloDragSortableTable
        columns={getRouteColumns(beginRouteEditing, handleDeleteRoute)}
        dataSource={getRouteData()}
        moveRow={reorderRoutes}
      />

      <SoloModal
        visible={showCreateRouteModal}
        width={500}
        title={'Create Route'}
        onClose={() => setShowCreateRouteModal(false)}>
        <CreateRouteModal
          defaultRouteParent={props.routeParent!}
          completeCreation={() => setShowCreateRouteModal(false)}
        />
      </SoloModal>
      <SoloModal
        visible={!!routeBeingEdited}
        width={500}
        title={'Edit Route'}
        onClose={() => setRouteBeingEdited(undefined)}>
        <CreateRouteModal
          defaultRouteParent={props.routeParent!}
          existingRoute={routeBeingEdited}
          completeCreation={finishRouteEditiing}
        />
      </SoloModal>
    </>
  );
};
