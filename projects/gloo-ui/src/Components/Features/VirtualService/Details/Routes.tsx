import css from '@emotion/css';
import styled from '@emotion/styled';
import { Popconfirm, Popover, Tag } from 'antd';
import { ReactComponent as GlooIcon } from 'assets/GlooEE.svg';
import { ReactComponent as KubeLogo } from 'assets/kube-logo.svg';
import RT from 'assets/route-table-icon.png';
import { ReactComponent as GreenPlus } from 'assets/small-green-plus.svg';
import { ReactComponent as UpstreamGroupLogo } from 'assets/upstream-group-icon.svg';
import { SoloDragSortableTable } from 'Components/Common/SoloDragSortableTable';
import { SoloModal } from 'Components/Common/SoloModal';

import { Route } from 'proto/gloo/projects/gateway/api/v1/virtual_service_pb';
import * as React from 'react';
import { shallowEqual, useDispatch } from 'react-redux';
import { NavLink } from 'react-router-dom';
import { updateRouteTable } from 'store/routeTables/actions';
import { routeTableAPI } from 'store/routeTables/api';
import { upstreamGroupAPI } from 'store/upstreamGroups/api';
import { upstreamAPI } from 'store/upstreams/api';
import { deleteRoute, shiftRoutes } from 'store/virtualServices/actions';
import { colors, TableActionCircle, TableActions } from 'Styles';
import useSWR from 'swr';
import {
  getIconFromSpec,
  getRouteHeaders,
  getRouteMatcher,
  getRouteMethods,
  getRouteQueryParams,
  getRouteSingleUpstream
} from 'utils/helpers';
import { CreateRouteModal } from '../Creation/CreateRouteModal';
import { RouteParent } from '../RouteTableDetails';
import { RouteTable } from 'proto/gloo/projects/gateway/api/v1/route_table_pb';
import { RouteTableDetails } from 'proto/solo-projects/projects/grpcserver/api/v1/routetable_pb';

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
    'content content-2';
  grid-template-columns: 1fr 1fr;
  grid-template-rows: 24px 1fr;
  grid-gap: 5px;
`;

const RouteTablePopoverRoutesContainer = css`
  grid-area: content;
  grid-column-gap: 5px;
`;
interface Props {
  routes: Route.AsObject[];
  routeParent?: RouteParent;
  deleteRouteFromRouteTable?: (matcher: string) => void;
}

const DestinationIcon: React.FC<{
  route: Route.AsObject;
  routeTable?: RouteTable.AsObject;
  multipleRouteTables?: RouteTableDetails.AsObject[];
}> = props => {
  const { route, routeTable, multipleRouteTables } = props;
  let icon = <GlooIcon style={{ width: '20px', paddingRight: '5px' }} />;
  let destination = '';
  const { data: upstreamsList, error } = useSWR(
    'listUpstreams',
    upstreamAPI.listUpstreams
  );
  const { data: upstreamGroupsList, error: upstreamGroupError } = useSWR(
    'listUpstreamGroups',
    upstreamGroupAPI.listUpstreamGroups
  );
  if (!upstreamsList) {
    return <div>Loading...</div>;
  }
  if (!!route && route?.routeAction !== undefined) {
    if (route.routeAction?.single !== undefined) {
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

      if (route?.routeAction !== undefined) {
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
      } else if (route?.routeAction!.upstreamGroup !== undefined) {
        let upstreamGroupDestination = upstreamGroupsList?.find(
          upstreamGroupDetails =>
            upstreamGroupDetails?.upstreamGroup?.metadata?.name ===
            route?.routeAction?.upstreamGroup?.name
        );
        destination = upstreamGroupDestination?.upstreamGroup?.metadata?.name!;
        icon = (
          <UpstreamGroupLogo style={{ width: '25px', paddingRight: '5px' }} />
        );
      }
    }
  }
  if (
    route?.delegateAction !== undefined &&
    route?.delegateAction?.ref !== undefined
  ) {
    icon = <img src={RT} style={{ width: '25px', paddingRight: '5px' }} />;
    destination = route.delegateAction?.ref?.name;
  }
  if (!!routeTable && routeTable?.metadata?.name) {
    icon = <img src={RT} style={{ width: '25px', paddingRight: '5px' }} />;
    destination = routeTable?.metadata?.name;
  }
  if (!!multipleRouteTables) {
    if (multipleRouteTables.length === 1) {
      icon = <img src={RT} style={{ width: '25px', paddingRight: '5px' }} />;
      destination = multipleRouteTables[0].routeTable?.metadata?.name || '';
    } else {
      destination = `${multipleRouteTables.length} Route Tables`;
    }
  }

  return (
    <>
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
    </>
  );
};

function checkRoutesProps(
  oldProps: Readonly<Props>,
  newProps: Readonly<Props>
): boolean {
  return shallowEqual(oldProps.routes, newProps.routes);
}
export const Routes: React.FC<Props> = props => {
  const { data: upstreamsList, error: upstreamError } = useSWR(
    'listUpstreams',
    upstreamAPI.listUpstreams
  );
  const { data: routeTablesList, error } = useSWR(
    'listRouteTables',
    routeTableAPI.listRouteTables
  );

  const { data: upstreamGroupsList, error: upstreamGroupsError } = useSWR(
    'listUpstreamGroups',
    upstreamGroupAPI.listUpstreamGroups
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

  if (!routeTablesList || !upstreamsList) {
    return <div>Loading...</div>;
  }

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
    let upstreamDestination = route?.routeAction?.single?.upstream;
    let upstreamGroupDestination = route?.routeAction?.upstreamGroup;
    let routeTableDestination = route?.delegateAction;
    const currentUpstream = upstreamsList?.find(
      us =>
        us?.upstream?.metadata?.name === upstreamDestination?.name &&
        us?.upstream?.metadata?.namespace === upstreamDestination?.namespace
    );
    const currentUpstreamGroup = upstreamGroupsList?.find(
      usg =>
        usg.upstreamGroup?.metadata?.name === upstreamGroupDestination?.name &&
        usg.upstreamGroup?.metadata?.namespace ===
          upstreamGroupDestination?.namespace
    );
    // get the routes of the route table

    // TODO: replace this logic with a call to apiserver to avoid duplication
    let matchedRouteTables = routeTablesList?.filter(routeTable => {
      if (!!routeTableDestination?.name || !!routeTableDestination?.namespace) {
        // this means there is a single route table matched, so
        return (
          routeTable?.routeTable?.metadata?.name ===
            routeTableDestination?.name &&
          routeTable?.routeTable?.metadata?.namespace ===
            routeTableDestination?.namespace
        );
      } else if (!!routeTableDestination?.ref) {
        return (
          routeTable?.routeTable?.metadata?.name ===
            routeTableDestination?.ref?.name &&
          routeTable?.routeTable?.metadata?.namespace ===
            routeTableDestination?.ref?.namespace
        );
      } else if (!!routeTableDestination?.selector) {
        if (!!routeTableDestination?.selector?.labelsMap) {
          let labelsMatch = routeTableDestination?.selector?.labelsMap.every(
            ([key, val]) =>
              routeTable.routeTable?.metadata?.labelsMap.some(
                ([key2, val2]) => key === key2 && val === val2
              )
          );
          if (!labelsMatch) {
            return false;
          }
        }

        let nsList = routeTableDestination?.selector?.namespacesList;
        if (Array.isArray(nsList) && nsList.length) {
          let hasAllSelector = routeTableDestination.selector?.namespacesList.includes(
            '*'
          );
          let matchesNs = routeTableDestination.selector.namespacesList.includes(
            routeTable?.routeTable?.metadata?.namespace!
          );
          return matchesNs || hasAllSelector;
        } else {
          // if namespaces is [], only check route tables in the same ns as the route parent
          return (
            routeTable?.routeTable?.metadata?.namespace ===
            routeParentRef.namespace
          );
        }
      }

      return false;
    });

    // get the routes of the route table
    //can match either by name/namespace or by selector, so could match more than one rt
    let previewRouteTable = routeTablesList?.find(
      rt =>
        (rt?.routeTable?.metadata?.name === routeTableDestination?.name &&
          rt?.routeTable?.metadata?.namespace ===
            routeTableDestination?.namespace) ||
        (rt?.routeTable?.metadata?.name === routeTableDestination?.ref?.name &&
          rt?.routeTable?.metadata?.namespace ===
            routeTableDestination?.ref?.namespace)
    );

    // if its a route or an upstream
    // get type
    let type: 'Upstream' | 'Upstream Group' | 'Route Table';
    if (routeTableDestination !== undefined) {
      type = 'Route Table';
    } else if (upstreamGroupDestination !== undefined) {
      type = 'Upstream Group';
    } else {
      type = 'Upstream';
    }

    let content;
    if (type === 'Route Table') {
      content = !!routeTableDestination?.selector ? (
        <Popover
          placement='bottom'
          title='Selector'
          content={
            <>
              <div css={RouteTablePopoverContainer}>
                <div
                  css={css`
                    grid-area: matcher;
                  `}>
                  Labels
                </div>
                <div
                  css={css`
                    grid-area: destination;
                  `}>
                  Namespaces
                </div>
                <div css={RouteTablePopoverRoutesContainer}>
                  <>
                    <div>
                      {routeTableDestination?.selector?.labelsMap.map(
                        ([key, val]) => (
                          <Tag key={`${key}-${val}`} color='blue'>
                            {key}: {val}
                          </Tag>
                        )
                      )}
                    </div>
                  </>
                </div>
                <div
                  css={css`
                    grid-area: content-2;
                    grid-column-gap: 5px;
                  `}>
                  <>
                    <div>
                      {routeTableDestination?.selector?.namespacesList.map(
                        ns => (
                          <Tag key={`${ns}`} color='blue'>
                            {ns}
                          </Tag>
                        )
                      )}
                    </div>
                  </>
                </div>
              </div>
              <div
                css={css`
                  margin-top: 10px;
                `}>
                <div
                  css={css`
                    margin-bottom: 5px;
                  `}>
                  Matching Route Tables
                </div>
                <div>
                  {matchedRouteTables?.map(matchedRT => (
                    <NavLink
                      to={`/routetables/${matchedRT?.routeTable?.metadata?.namespace}/${matchedRT?.routeTable?.metadata?.name}`}>
                      <div
                        css={css`
                          color: #2196c9;
                          cursor: pointer;
                          font-weight: bold;
                        `}>
                        <DestinationIcon
                          route={route}
                          routeTable={matchedRT.routeTable!}
                        />
                      </div>
                    </NavLink>
                  ))}
                </div>
              </div>
            </>
          }>
          <div
            css={css`
              color: #2196c9;
              cursor: pointer;
              font-weight: bold;
            `}>
            <DestinationIcon
              route={route}
              multipleRouteTables={matchedRouteTables}
            />
          </div>
        </Popover>
      ) : (
        <Popover
          placement='bottom'
          content={
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
          }
          trigger='hover'>
          {!!previewRouteTable && (
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
          )}
        </Popover>
      );
    } else if (type === 'Upstream Group') {
      content = (
        <Popover
          placement='bottom'
          content={
            <>
              <div
                css={css`
                  display: grid;
                  grid-template-areas:
                    'upstream destination'
                    'content content';
                  grid-template-columns: 1fr 1fr;
                  grid-template-rows: 24px 1fr;
                  grid-gap: 10px;
                `}>
                <div
                  css={css`
                    grid-area: upstream;
                  `}>
                  Upstream
                </div>
                <div
                  css={css`
                    grid-area: destination;
                  `}>
                  Weight
                </div>
                <div
                  css={css`
                    grid-area: content;
                    display: grid;
                    grid-template-columns: 1fr 1fr;
                    grid-column-gap: 10px;
                  `}>
                  {currentUpstreamGroup?.upstreamGroup?.destinationsList?.map(
                    dest => (
                      <React.Fragment
                        key={`${dest.destination?.upstream?.name}`}>
                        <div
                          css={css`
                            display: flex;
                            flex-direction: row;
                            justify-content: flex-start;
                            align-items: center;
                          `}>
                          {dest.destination?.upstream?.name}
                        </div>
                        <div>{dest.weight}%</div>
                      </React.Fragment>
                    )
                  )}
                </div>
              </div>
              <div>
                {`${currentUpstreamGroup?.upstreamGroup?.destinationsList.length} upsteams`}
              </div>
            </>
          }
          trigger='hover'>
          <NavLink
            to={`/upstreams/upstreamgroups/${currentUpstreamGroup?.upstreamGroup?.metadata?.namespace}/${currentUpstreamGroup?.upstreamGroup?.metadata?.name}`}>
            <div
              css={css`
                color: #2196c9;
                cursor: pointer;
                font-weight: bold;
              `}>
              <DestinationIcon route={route} />
            </div>
          </NavLink>
        </Popover>
      );
    } else if (type === 'Upstream') {
      content = (
        <NavLink
          to={`/upstreams/${currentUpstream?.upstream?.metadata?.namespace}/${currentUpstream?.upstream?.metadata?.name}`}>
          <div
            css={css`
              color: #2196c9;
              cursor: pointer;
              font-weight: bold;
            `}>
            <DestinationIcon route={route} />
          </div>
        </NavLink>
      );
    }

    // style to match design
    return <>{content}</>;
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
};
