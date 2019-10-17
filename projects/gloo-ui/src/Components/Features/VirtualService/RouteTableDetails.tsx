import styled from '@emotion/styled';
import { Breadcrumb, Popconfirm, Spin } from 'antd';
import { ReactComponent as EditPencil } from 'assets/edit-pencil.svg';
import { ReactComponent as RouteTableIcon } from 'assets/route-table-icon.svg';
import { ReactComponent as GreenPlus } from 'assets/small-green-plus.svg';
import { ConfigDisplayer } from 'Components/Common/DisplayOnly/ConfigDisplayer';
import { FileDownloadLink } from 'Components/Common/FileDownloadLink';
import { SectionCard } from 'Components/Common/SectionCard';
import { SoloDragSortableTable } from 'Components/Common/SoloDragSortableTable';
import { SoloModal } from 'Components/Common/SoloModal';
import { RouteTable } from 'proto/github.com/solo-io/gloo/projects/gateway/api/v1/route_table_pb';
import {
  Route,
  VirtualService
} from 'proto/github.com/solo-io/gloo/projects/gateway/api/v1/virtual_service_pb';
import React from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { useHistory, useParams } from 'react-router-dom';
import { AppState } from 'store';
import {
  listRouteTables,
  updateRouteTableYaml,
  updateRouteTable
} from 'store/routeTables/actions';
import {
  colors,
  healthConstants,
  TableActionCircle,
  TableActions
} from 'Styles';
import {
  getRouteHeaders,
  getRouteMatcher,
  getRouteMethods,
  getRouteQueryParams,
  getRouteSingleUpstream
} from 'utils/helpers';
import {
  ConfigurationToggle,
  DetailsContent,
  DetailsSection,
  DetailsSectionTitle,
  YamlLink
} from './Details/VirtualServiceDetails';
import {
  CreateRouteModal,
  CreateRouteValuesType
} from './Creation/CreateRouteModal';

const RouteMatch = styled.div`
  max-width: 200px;
  max-height: 70px;
  overflow: hidden;
  text-overflow: ellipsis;
`;

const RouteSectionTitle = styled.div`
  font-size: 18px;
  color: ${colors.novemberGrey};
  margin-top: 10px;
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

const YamlLinks = styled.div`
  display: flex;
  align-items: center;
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
export interface RouteParent
  extends VirtualService.AsObject,
    RouteTable.AsObject {}
export const RouteTableDetails = () => {
  let history = useHistory();
  let { routetablename, routetablenamespace } = useParams();

  const [showCreateRouteModal, setShowCreateRouteModal] = React.useState(false);
  const [showConfiguration, setShowConfiguration] = React.useState(false);
  const [routeBeingEdited, setRouteBeingEdited] = React.useState<
    Route.AsObject | undefined
  >(undefined);
  const routeTablesList = useSelector(
    (state: AppState) => state.routeTables.routeTablesList
  );

  const yamlError = useSelector(
    (state: AppState) => state.virtualServices.yamlParseError
  );

  const dispatch = useDispatch();

  React.useEffect(() => {
    if (routeTablesList.length || !routeTableDetails) {
      dispatch(listRouteTables());
    }
  }, [routeTablesList.length]);

  let routeTableDetails = routeTablesList.find(
    rtD =>
      !!rtD &&
      rtD.routeTable &&
      rtD.routeTable.metadata!.name === routetablename
  )!;

  const [routesList, setRoutesList] = React.useState<Route.AsObject[]>([]);

  React.useEffect(() => {
    if (
      routesList.length === 0 &&
      !!routeTableDetails &&
      !!routeTableDetails.routeTable
    ) {
      setRoutesList(routeTableDetails.routeTable.routesList);
    }
  }, [routeTablesList.length]);
  if (!routeTablesList || !routeTableDetails || !routeTableDetails.routeTable) {
    return (
      <>
        <Breadcrumb />
        <Spin size='large' />
      </>
    );
  }

  let { routeTable, raw } = routeTableDetails;
  const saveYamlChange = (newYaml: string) => {
    dispatch(
      updateRouteTableYaml({
        editedYamlData: {
          editedYaml: newYaml,
          ref: {
            name: routeTable!.metadata!.name,
            namespace: routeTable!.metadata!.namespace
          }
        }
      })
    );
  };

  const getRouteData = () => {
    const existingRoutes = routesList.map(route => {
      const upstreamName = getRouteSingleUpstream(route) || '';
      const { matcher, matchType } = getRouteMatcher(route);
      return {
        key: `${matcher}-${upstreamName}`,
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

  const handleDeleteRoute = (matcherToDelete: string) => {
    let index = routesList.findIndex(
      route => getRouteMatcher(route).matcher === matcherToDelete
    );
    const newList = routesList.filter(
      route => getRouteMatcher(route).matcher !== matcherToDelete
    );

    dispatch(
      updateRouteTable({
        routeTable: {
          routesList: newList
        }
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

  const handleCreateRoute = (values: CreateRouteValuesType) => {
    dispatch(
      updateRouteTable({
        routeTable: {
          routesList: [
            ...routeTable.routesList,
            {
              matchersList: [{
                prefix: values.matchType === 'PREFIX' ? values.path : '',
                exact: values.matchType === 'EXACT' ? values.path : '',
                regex: values.matchType === 'REGEX' ? values.path : '',
                methodsList: values.methods,
                headersList: values.headers,
                queryParametersList: values.queryParameters
              }],
              routeAction: {
                single: {
                  upstream: {
                    name: values.upstream!.metadata!.name,
                    namespace: values.upstream!.metadata!.namespace
                  },
                  destinationSpec: values.destinationSpec
                }
              }
            }
          ]
        }
      })
    );
  };

  const reorderRoutes = (dragIndex: number, hoverIndex: number) => {
    const movedRoute = routesList.splice(dragIndex, 1)[0];

    let newRoutesList = [...routesList];
    newRoutesList.splice(hoverIndex, 0, movedRoute);

    dispatch(
      updateRouteTable({
        routeTable: {
          routesList: newRoutesList
        }
      })
    );
    setRoutesList(newRoutesList);
  };

  const headerInfo = [
    {
      title: 'namespace',
      value: routetablenamespace!
    }
  ];
  return (
    <>
      <Breadcrumb />

      <SectionCard
        cardName={
          routeTable.metadata
            ? routeTable.metadata!.name
            : routetablename
            ? routetablename
            : 'Error'
        }
        logoIcon={<RouteTableIcon />}
        health={
          routeTable!.status
            ? routeTable!.status!.state
            : healthConstants.Pending.value
        }
        headerSecondaryInformation={headerInfo}
        healthMessage={
          routeTable!.status && routeTable!.status!.reason.length
            ? routeTable!.status!.reason
            : 'Service Status'
        }
        onClose={() => history.push(`/virtualservices/`)}>
        <RouteSectionTitle>
          <b>Routes</b>
          <YamlLinks>
            {!!raw && (
              <>
                <ConfigurationToggle
                  onClick={() => setShowConfiguration(s => !s)}>
                  {showConfiguration ? 'Hide' : 'View'} YAML Configuration
                </ConfigurationToggle>
                <FileDownloadLink
                  fileContent={raw.content}
                  fileName={raw.fileName}
                />
              </>
            )}
            <ModalTrigger
              data-testid='create-new-route-modal'
              onClick={() => setShowCreateRouteModal(true)}>
              <>
                <StyledGreenPlus />
                <b>Create Route</b>
              </>
            </ModalTrigger>
          </YamlLinks>
        </RouteSectionTitle>
        <DetailsContent
          configurationShowing={showConfiguration}
          style={{ gridTemplateRows: 'unset' }}>
          <DetailsSection></DetailsSection>
          <div style={{ height: '40px' }}></div>
          {showConfiguration && (
            <DetailsSection>
              <DetailsSectionTitle>YAML Configuration</DetailsSectionTitle>
              <ConfigDisplayer
                content={raw ? raw.content : ''}
                asEditor
                yamlError={yamlError}
                saveEdits={saveYamlChange}
              />
            </DetailsSection>
          )}
          <DetailsSection>
            <>
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
                  defaultRouteParent={routeTable! as RouteParent}
                  completeCreation={() => setShowCreateRouteModal(false)}
                  createRouteFn={handleCreateRoute}
                  lockVirtualService
                />
              </SoloModal>
              <SoloModal
                visible={!!routeBeingEdited}
                width={500}
                title={'Edit Route'}
                onClose={() => setRouteBeingEdited(undefined)}>
                <CreateRouteModal
                  defaultRouteParent={routeTable as RouteParent}
                  existingRoute={routeBeingEdited}
                  completeCreation={finishRouteEditiing}
                  createRouteFn={handleCreateRoute}
                />
              </SoloModal>
            </>
          </DetailsSection>
        </DetailsContent>
      </SectionCard>
    </>
  );
};
