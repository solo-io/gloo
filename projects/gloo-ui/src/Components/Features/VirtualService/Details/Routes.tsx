import styled from '@emotion/styled';
import { Popconfirm, Popover } from 'antd';
import { ReactComponent as KubeLogo } from 'assets/kube-logo.svg';
import { ReactComponent as GreenPlus } from 'assets/small-green-plus.svg';
import { ReactComponent as RouteTableIcon } from 'assets/route-table-icon.svg';
import { ReactComponent as GlooIcon } from 'assets/GlooEE.svg';
import { SoloDragSortableTable } from 'Components/Common/SoloDragSortableTable';
import { SoloModal } from 'Components/Common/SoloModal';
import * as React from 'react';
import { useDispatch, shallowEqual, useSelector } from 'react-redux';
import { shiftRoutes, deleteRoute } from 'store/virtualServices/actions';
import { colors, TableActionCircle, TableActions } from 'Styles';
import {
  getRouteHeaders,
  getRouteMatcher,
  getRouteMethods,
  getRouteQueryParams,
  getRouteSingleUpstream,
  getIcon,
  getIconFromSpec
} from 'utils/helpers';
import { RouteParent } from '../RouteTableDetails';
import { CreateRouteModal } from '../Creation/CreateRouteModal';
import css from '@emotion/css';
import { Route } from 'proto/github.com/solo-io/gloo/projects/gateway/api/v1/virtual_service_pb';
import { AppState } from 'store';
import { NavLink } from 'react-router-dom';
import RT from 'assets/route-table-icon.png';
import { updateRouteTable } from 'store/routeTables/actions';
const RouteMatch = styled.div`
  max-width: 200px;
  max-height: 70px;
  overflow: hidden;
  font-weight: bold;
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

const RouteTablePopoverContainer = css`
  display: grid;
  grid-template-areas:
    'matcher destination'
    'content content';
  grid-template-columns: 1fr 3fr;
  grid-template-rows: 24px 1fr;
  grid-gap: 10px;
`;
const RouteTablePopoverRoutesContainer = css`
  grid-area: content;
  display: grid;
  grid-template-columns: 1fr 3fr;
  grid-column-gap: 10px;
`;
interface Props {
  routes: Route.AsObject[];
  routeParent?: RouteParent;
  deleteRouteFromRouteTable?: (matcher: string) => void;
}

const DestinationIcon: React.FC<{ route: Route.AsObject }> = ({ route }) => {
  let icon = <GlooIcon style={{ width: '20px', paddingRight: '5px' }} />;
  let destination = '';
  const upstreamsList = useSelector(
    (state: AppState) => state.upstreams.upstreamsList
  );
  if (route.routeAction !== undefined) {
    let upstreamDestination = upstreamsList.find(
      upstreamDetails =>
        upstreamDetails?.upstream?.metadata?.name ===
        route?.routeAction?.single?.upstream?.name
    );
    let upstreamSpec = upstreamDestination?.upstream;
    if (upstreamSpec !== undefined) {
      icon = getIconFromSpec(upstreamSpec);
    } else {
      icon = <KubeLogo style={{ width: '20px', paddingRight: '5px' }} />;
    }

    destination = getRouteSingleUpstream(route);
  }
  if (route.delegateAction !== undefined) {
    icon = <img src={RT} style={{ width: '25px', paddingRight: '5px' }} />;
    destination = route.delegateAction.name;
  }
  return (
    <div
      css={css`
        display: flex;
        flex-direction: row;
        justify-content: flex-start;
        align-items: center;
      `}>
      {icon}
      {destination}
    </div>
  );
};

function checkRoutesProps(
  oldProps: Readonly<Props>,
  newProps: Readonly<Props>
): boolean {
  return shallowEqual(oldProps.routes, newProps.routes);
}
export const Routes: React.FC<Props> = React.memo(props => {
  const routeTablesList = useSelector(
    (state: AppState) => state.routeTables.routeTablesList,
    shallowEqual
  );

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

  React.useEffect(() => {
    setRoutesList([...props.routes]);
  }, [props.routes.length]);

  const getRouteData = () => {
    const existingRoutes = props.routes.map(route => {
      const upstreamName = getRouteSingleUpstream(route) || '';
      const { matcher, matchType } = getRouteMatcher(route);
      return {
        key: `${matcher}-${upstreamName}`,
        matcher: matcher,
        pathMatch: matchType,
        method: getRouteMethods(route),
        destinationName: getDestinationIcon(route),
        upstreamName: upstreamName,
        header: getRouteHeaders(route),
        queryParams: getRouteQueryParams(route),
        actions: matcher
      };
    });

    return existingRoutes;
  };

  function getDestinationIcon(route: Route.AsObject): React.ReactNode {
    let upstreamDestination = route?.routeAction?.single?.upstream?.name;
    let routeTableDestination = route?.delegateAction?.name;

    // get the routes of the route table
    const previewRouteTable = routeTablesList.find(
      rt => rt?.routeTable?.metadata?.name === routeTableDestination
    );
    // if its a route or an upstream
    // get type
    let type: 'Upstream' | 'Route Table';
    if (routeTableDestination !== undefined) {
      type = 'Route Table';
    } else {
      type = 'Upstream';
    }

    // style to match design
    return (
      <Popover
        placement='bottom'
        content={
          type === 'Route Table' ? (
            <>
              <div css={RouteTablePopoverContainer}>
                <div
                  css={css`
                    grid-area: matcher;
                  `}>
                  Matcher
                </div>
                <div
                  css={css`
                    grid-area: destination;
                  `}>
                  Destination
                </div>
                <div css={RouteTablePopoverRoutesContainer}>
                  {previewRouteTable?.routeTable?.routesList.map(route => (
                    <>
                      <div
                        css={css`
                          display: flex;
                          flex-direction: row;
                          justify-content: flex-start;
                          align-items: center;
                        `}>
                        {route.matchersList[0]?.prefix}
                      </div>

                      <DestinationIcon route={route} />
                    </>
                  ))}
                </div>
              </div>
              <div>
                {`${previewRouteTable?.routeTable?.routesList.length} total routes`}
              </div>
            </>
          ) : (
            <div
              css={css`
                display: flex;
                justify-content: center;
              `}>
              Upstream
            </div>
          )
        }
        trigger='hover'>
        {!!previewRouteTable ? (
          <NavLink
            to={`/routetables/${previewRouteTable?.routeTable?.metadata?.namespace}/${previewRouteTable?.routeTable?.metadata?.name}`}>
            <div
              css={css`
                color: #2196c9;
                cursor: pointer;
                font-weight: bold;
              `}>
              <DestinationIcon route={route} />
            </div>
          </NavLink>
        ) : (
          <div
            css={css`
              color: #2196c9;
              cursor: pointer;
              font-weight: bold;
            `}>
            <DestinationIcon route={route} />
          </div>
        )}
      </Popover>
    );
  }

  const handleDeleteRoute = (matcherToDelete: string, row: any) => {
    let isRouteTable = routeTablesList.find(
      rt =>
        rt?.routeTable?.metadata?.name === props.routeParent?.metadata?.name &&
        rt?.routeTable?.metadata?.namespace ===
          props?.routeParent?.metadata?.namespace
    );
    if (isRouteTable !== undefined && !('virtualHost' in props.routeParent!)) {
      const newList = routeTablesList
        .flatMap(rtd => rtd!.routeTable!.routesList)
        .filter(route => {
          return (
            getRouteMatcher(route).matcher !== matcherToDelete &&
            row.upstreamName !== route.delegateAction?.name
          );
        });

      dispatch(
        updateRouteTable({
          routeTable: {
            ...isRouteTable,
            metadata: {
              ...props.routeParent?.metadata!,
              name: isRouteTable?.routeTable?.metadata?.name!,
              namespace: isRouteTable?.routeTable?.metadata?.namespace!
            },
            routesList: newList
          }
        })
      );
    } else {
      let index = routesList.findIndex(
        route =>
          getRouteMatcher(route).matcher === matcherToDelete &&
          row.upstreamName === getRouteSingleUpstream(route)
      );
      const newList = routesList.filter(
        route =>
          getRouteMatcher(route).matcher !== matcherToDelete &&
          row.upstreamName !== getRouteSingleUpstream(route)
      );

      dispatch(
        deleteRoute({
          virtualServiceRef: routeParentRef,
          index
        })
      );
      setRoutesList(newList);
    }
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
    let isRouteTable = routeTablesList.find(
      rt =>
        rt?.routeTable?.metadata?.name === props.routeParent?.metadata?.name &&
        rt?.routeTable?.metadata?.namespace ===
          props?.routeParent?.metadata?.namespace
    );
    const movedRoute = routesList.splice(dragIndex, 1)[0];

    let newRoutesList = [...routesList];
    newRoutesList.splice(hoverIndex, 0, movedRoute);
    if ('virtualHost' in props.routeParent!) {
      dispatch(
        shiftRoutes({
          virtualServiceRef: routeParentRef,
          fromIndex: dragIndex,
          toIndex: hoverIndex
        })
      );
    } else {
      dispatch(
        updateRouteTable({
          routeTable: {
            ...isRouteTable,
            metadata: {
              ...props.routeParent?.metadata!,
              name: isRouteTable?.routeTable?.metadata?.name!,
              namespace: isRouteTable?.routeTable?.metadata?.namespace!
            },
            routesList: newRoutesList
          }
        })
      );
    }
    setRoutesList(newRoutesList);
  };

  const getRouteColumns = (
    showEditRouteModal: (matcher: string) => void,
    deleteRoute: (matcher: string, row: any) => any
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
        dataIndex: 'destinationName'
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
        render: (matcher: string, row: any) => {
          return (
            <TableActions>
              {/* disallowing edits until further notice TODO: (ascampos) */}
              {/* <TableActionCircle onClick={() => showEditRouteModal(matcher)}>
              <EditPencil />
            </TableActionCircle> */}

              <div style={{ marginLeft: '5px' }}>
                <Popconfirm
                  onConfirm={() => deleteRoute(matcher, row)}
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
          defaultRouteParent={props.routeParent}
          completeCreation={() => setShowCreateRouteModal(false)}
        />
      </SoloModal>
      {/* Temporarily removed edit route functionality */}
      {/* <SoloModal
        visible={!!routeBeingEdited}
        width={500}
        title={'Edit Route'}
        onClose={() => setRouteBeingEdited(undefined)}>
        <CreateRouteModal
          defaultRouteParent={props.routeParent!}
          existingRoute={routeBeingEdited}
          completeCreation={finishRouteEditiing}
        />
      </SoloModal> */}
    </>
  );
}, checkRoutesProps);
