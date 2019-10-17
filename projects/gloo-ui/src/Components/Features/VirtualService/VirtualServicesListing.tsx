import styled from '@emotion/styled';
import { Popconfirm } from 'antd';
import { ReactComponent as Gloo } from 'assets/Gloo.svg';
import { ReactComponent as RouteTableIcon } from 'assets/route-table-icon.svg';

import { Breadcrumb } from 'Components/Common/Breadcrumb';
import { CardsListing } from 'Components/Common/CardsListing';
import { CatalogTableToggle } from 'Components/Common/CatalogTableToggle';
import { FileDownloadActionCircle } from 'Components/Common/FileDownloadLink';
import { HealthIndicator } from 'Components/Common/HealthIndicator';
import { HealthInformation } from 'Components/Common/HealthInformation';
import {
  CheckboxFilterProps,
  ListingFilter,
  RadioFilterProps,
  StringFilterProps,
  TypeFilterProps
} from 'Components/Common/ListingFilter';
import { SectionCard } from 'Components/Common/SectionCard';
import { SoloModal } from 'Components/Common/SoloModal';
import { SoloTable } from 'Components/Common/SoloTable';
import { Status } from 'proto/github.com/solo-io/solo-kit/api/v1/status_pb';
import { VirtualServiceDetails } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/virtualservice_pb';
import * as React from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { Route, useHistory, useRouteMatch, useLocation } from 'react-router';
import { AppState } from 'store';
import {
  deleteVirtualService,
  listVirtualServices,
  createRoute
} from 'store/virtualServices/actions';
import { colors, healthConstants } from 'Styles';
import {
  TableActionCircle,
  TableActions,
  TableHealthCircleHolder
} from 'Styles/table';
import { getResourceStatus, getVSDomains, RadioFilters } from 'utils/helpers';
import { CreateVirtualServiceModal } from './Creation/CreateVirtualServiceModal';
import { listRouteTables, deleteRouteTable } from 'store/routeTables/actions';
import { RouteTableDetails } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/routetable_pb';
import { CardType } from 'Components/Common/Card';
import { RouteParent } from './RouteTableDetails';
import { CreateRouteTableModal } from './Creation/CreateRouteTableModal';
import {
  CreateRouteModal,
  CreateRouteValuesType
} from './Creation/CreateRouteModal';
import { routeTables } from 'store/routeTables/api';

const TableLink = styled.div`
  cursor: pointer;
  color: ${colors.seaBlue};
`;

const TableDomains = styled.div`
  max-width: 200px;
  max-height: 70px;
  overflow: hidden;
  text-overflow: ellipsis;
`;

const getRouteTableColumns = (
  startCreatingRoute: (routeParent: RouteParent) => void,
  deleteRT: (name: string, namespace: string) => void
) => {
  return [
    {
      title: 'Name',
      dataIndex: 'name',
      render: (nameObject: {
        goToVirtualService: () => void;
        displayName: string;
      }) => {
        return (
          <TableLink onClick={nameObject.goToVirtualService}>
            {nameObject.displayName}
          </TableLink>
        );
      }
    },

    {
      title: 'Namespace',
      dataIndex: 'metadata.namespace'
    },
    {
      title: 'Version',
      dataIndex: 'metadata.resourceVersion'
    },
    {
      title: 'Status',
      dataIndex: 'status',
      render: (healthStatus: Status.AsObject) => (
        <div>
          <TableHealthCircleHolder>
            <HealthIndicator healthStatus={healthStatus.state} />
          </TableHealthCircleHolder>
          <HealthInformation healthStatus={healthStatus} />
        </div>
      )
    },
    {
      title: 'Routes',
      dataIndex: 'routes'
    },
    {
      title: 'BR Limit',
      dataIndex: 'brLimit'
    },
    {
      title: 'AR Limit',
      dataIndex: 'arLimit'
    },
    {
      title: 'Actions',
      dataIndex: 'actions',
      render: (routeTableDetails: RouteTableDetails.AsObject) => {
        const routeTable = routeTableDetails.routeTable!;
        return (
          <TableActions>
            <Popconfirm
              onConfirm={() =>
                deleteRT(
                  routeTable.metadata!.name,
                  routeTable.metadata!.namespace
                )
              }
              title={'Are you sure you want to delete this virtual service? '}
              okText='Yes'
              cancelText='No'>
              <TableActionCircle>x</TableActionCircle>
            </Popconfirm>
            {!!routeTableDetails.raw && (
              <FileDownloadActionCircle
                fileContent={routeTableDetails.raw.content}
                fileName={routeTableDetails.raw.fileName}
              />
            )}
            <TableActionCircle
              onClick={() => startCreatingRoute(routeTable as RouteParent)}>
              +
            </TableActionCircle>
          </TableActions>
        );
      }
    }
  ];
};
const getTableColumns = (
  startCreatingRoute: (routeParent: RouteParent) => void,
  deleteVirtualService: (name: string, namespace: string) => void
) => {
  return [
    {
      title: 'Name',
      dataIndex: 'name',
      render: (nameObject: {
        goToVirtualService: () => void;
        displayName: string;
      }) => {
        return (
          <TableLink onClick={nameObject.goToVirtualService}>
            {nameObject.displayName}
          </TableLink>
        );
      }
    },
    {
      title: 'Domain',
      dataIndex: 'domains',
      render: (domains: string) => {
        return <TableDomains>{domains}</TableDomains>;
      }
    },
    {
      title: 'Namespace',
      dataIndex: 'metadata.namespace'
    },
    {
      title: 'Version',
      dataIndex: 'metadata.resourceVersion'
    },
    {
      title: 'Status',
      dataIndex: 'status',
      render: (healthStatus: Status.AsObject) => (
        <div>
          <TableHealthCircleHolder>
            <HealthIndicator healthStatus={healthStatus.state} />
          </TableHealthCircleHolder>
          <HealthInformation healthStatus={healthStatus} />
        </div>
      )
    },
    {
      title: 'Routes',
      dataIndex: 'routes'
    },
    {
      title: 'BR Limit',
      dataIndex: 'brLimit'
    },
    {
      title: 'AR Limit',
      dataIndex: 'arLimit'
    },
    {
      title: 'Actions',
      dataIndex: 'actions',
      render: (virtualServiceDetails: VirtualServiceDetails.AsObject) => {
        const virtualService = virtualServiceDetails.virtualService!;
        return (
          <TableActions>
            <Popconfirm
              onConfirm={() =>
                deleteVirtualService(
                  virtualService.metadata!.name,
                  virtualService.metadata!.namespace
                )
              }
              title={'Are you sure you want to delete this virtual service? '}
              okText='Yes'
              cancelText='No'>
              <TableActionCircle>x</TableActionCircle>
            </Popconfirm>
            {!!virtualServiceDetails.raw && (
              <FileDownloadActionCircle
                fileContent={virtualServiceDetails.raw.content}
                fileName={virtualServiceDetails.raw.fileName}
              />
            )}
            <TableActionCircle
              onClick={() => startCreatingRoute(virtualService as RouteParent)}>
              +
            </TableActionCircle>
          </TableActions>
        );
      }
    }
  ];
};

const StringFilters: StringFilterProps[] = [
  {
    displayName: 'Filter By Name...',
    placeholder: 'Filter by name...',
    value: ''
  }
];

const Heading = styled.div`
  display: flex;
  justify-content: space-between;
  margin-bottom: 20px;
`;

const Action = styled.div`
  display: flex;
  flex-direction: row;
  align-items: center;
  align-items: baseline;
`;

const EmptyPrompt = styled.div`
  display: flex;
  align-items: center;
  font-size: 14px;
`;

export const VirtualServicesListing = () => {
  let history = useHistory();
  let location = useLocation();
  let match = useRouteMatch({
    path: '/virtualservices/'
  })!;

  let params = new URLSearchParams(location.search);

  const [catalogNotTable, setCatalogNotTable] = React.useState(true);

  // redux
  const dispatch = useDispatch();
  const [isLoading, setIsLoading] = React.useState(false);
  const virtualServicesList = useSelector(
    (state: AppState) => state.virtualServices.virtualServicesList
  );

  const routeTablesList = useSelector(
    (state: AppState) => state.routeTables.routeTablesList
  );

  React.useEffect(() => {
    if (virtualServicesList.length) {
      setIsLoading(false);
    } else {
      dispatch(listVirtualServices());
    }
  }, [virtualServicesList.length]);

  React.useEffect(() => {
    if (routeTablesList.length) {
      setIsLoading(false);
    } else {
      dispatch(listRouteTables());
    }
  }, [routeTablesList.length]);

  const [
    routeParentForRouteCreation,
    setRouteParentForRouteCreation
  ] = React.useState<RouteParent | undefined>(undefined);

  interface VSorRTData
    extends VirtualServiceDetails.AsObject,
      RouteTableDetails.AsObject {}

  const getRouteTableCatalogData = (
    nameFilter: string,
    routeTableDetailsList: RouteTableDetails.AsObject[],
    radioFilter: string
  ) => {
    const catalogData = routeTableDetailsList.map(routeTableDetails => {
      let routeTable = routeTableDetails.routeTable!;
      let routeTablePrefix =
        routeTable.routesList.length > 0
          ? routeTable.routesList[0].matchersList[0]!.prefix
          : '';
      return {
        ...routeTable,
        healthStatus: routeTable.status
          ? routeTable.status.state
          : healthConstants.Pending.value,
        cardTitle: routeTable.metadata!.name,
        cardSubtitle: [
          {
            key: `Route${routeTable.routesList.length === 1 ? '' : 's'}`,
            value: `${routeTable.routesList.length}`
          }
        ],
        onRemoveCard: () =>
          dispatch(
            deleteRouteTable({
              ref: {
                name: routeTable.metadata!.name,
                namespace: routeTable.metadata!.namespace
              }
            })
          ),
        removeConfirmText: 'Are you sure you want to delete this Route Table?',
        onExpanded: () => {},
        onClick: () => {
          history.push({
            pathname: `/routetables/${routeTable.metadata!.namespace}/${
              routeTable.metadata!.name
            }`
          });
        },
        onCreate: () =>
          setRouteParentForRouteCreation(routeTable as RouteParent),
        downloadableContent: routeTableDetails.raw
      };
    });
    return catalogData
      .filter(row =>
        row.cardTitle.toLowerCase().includes(nameFilter.toLowerCase())
      )
      .filter(row => getResourceStatus(row).includes(radioFilter));
  };

  const getUsableCatalogData = (
    nameFilter: string,
    data: VSorRTData[],
    radioFilter: string
  ) => {
    const dataUsed = data.map(resource => {
      const virtualService = resource.virtualService!;
      let numRoutes = virtualService.virtualHost
        ? virtualService.virtualHost!.routesList.length
        : 0;
      return {
        ...virtualService,
        healthStatus: virtualService.status
          ? virtualService.status.state
          : healthConstants.Pending.value,
        cardTitle: virtualService.displayName || virtualService.metadata!.name,
        cardSubtitle: [
          { key: 'Domains', value: getVSDomains(virtualService) || '' },
          {
            key: `Route${
              virtualService.virtualHost
                ? virtualService.virtualHost!.routesList.length === 1
                  ? ''
                  : 's'
                : ''
            }`,
            value: `${numRoutes}`
          }
        ],
        onRemoveCard: () =>
          deleteVS(
            virtualService.metadata!.name,
            virtualService.metadata!.namespace
          ),
        removeConfirmText:
          'Are you sure you want to delete this virtual service?',
        onExpanded: () => {},
        onClick: () => {
          history.push({
            pathname: `${match.url}${virtualService.metadata!.namespace}/${
              virtualService.metadata!.name
            }`
          });
        },
        onCreate: () =>
          setRouteParentForRouteCreation(virtualService as RouteParent),
        downloadableContent: resource.raw
      };
    });

    return dataUsed
      .filter(row =>
        row.cardTitle.toLowerCase().includes(nameFilter.toLowerCase())
      )
      .filter(row => getResourceStatus(row).includes(radioFilter));
  };
  const getRouteTableData = (
    nameFilter: string,
    data: VSorRTData[],
    radioFilter: string
  ) => {
    const dataUsed = data.map(rtDetails => {
      const resource = rtDetails.routeTable!;

      return {
        ...resource,
        name: {
          displayName: resource.metadata!.name,
          goToVirtualService: () => {
            history.push({
              pathname: `/routetables/${resource.metadata!.namespace}/${
                resource.metadata!.name
              }`,
              search: location.search
            });
          }
        },
        domains: {},
        routes: resource.routesList.length,
        status: resource.status,
        key: `${resource.metadata!.name}`,
        actions: rtDetails
      };
    });

    return dataUsed
      .filter(row => row.name.displayName.includes(nameFilter))
      .filter(row => getResourceStatus(row).includes(radioFilter));
  };
  const getUsableTableData = (
    nameFilter: string,
    data: VSorRTData[],
    radioFilter: string
  ) => {
    const dataUsed = data.map(vsDetails => {
      const resource = vsDetails.virtualService! || vsDetails.routeTable!;

      return {
        ...resource,
        name: {
          displayName: resource.metadata!.name,
          goToVirtualService: () => {
            history.push({
              pathname: `${match.path}${resource.metadata!.namespace}/${
                resource.metadata!.name
              }`,
              search: location.search
            });
          }
        },
        domains: getVSDomains(resource),
        routes: resource.virtualHost!.routesList.length,
        status: resource.status,
        key: `${resource.metadata!.name}`,
        actions: vsDetails
      };
    });

    return dataUsed
      .filter(row => row.name.displayName.includes(nameFilter))
      .filter(row => getResourceStatus(row).includes(radioFilter));
  };

  function deleteVS(name: string, namespace: string) {
    dispatch(deleteVirtualService({ ref: { name, namespace } }));
  }

  function deleteRT(name: string, namespace: string) {
    dispatch(deleteRouteTable({ ref: { name, namespace } }));
  }

  function handleFilterChange(
    strings: StringFilterProps[],
    types: TypeFilterProps[],
    checkboxes: CheckboxFilterProps[],
    radios: RadioFilterProps[]
  ) {
    params.set('status', radios[0].choice || '');
    history.replace({
      pathname: `${location.pathname}`,
      search: radios[0].choice
        ? `?${'status'}=${radios[0].choice}`
        : params.get('status') || ''
    });
  }

  return (
    <div>
      <Heading>
        <Breadcrumb />
        <Action>
          <CreateVirtualServiceModal />
          <CreateRouteTableModal />
          <CatalogTableToggle
            listIsSelected={!catalogNotTable}
            onToggle={() => {
              history.push({
                pathname: `${match.path}${
                  location.pathname.includes('table') ? '' : 'table'
                }`
              });
              setCatalogNotTable(cNt => !cNt);
            }}
          />
        </Action>
      </Heading>
      <ListingFilter
        showLabels
        strings={StringFilters}
        checkboxes={[
          { displayName: 'Virtual Services', value: false },
          { displayName: 'Route Tables', value: false }
        ]}
        radios={[
          {
            ...RadioFilters,
            choice: params.has('status') ? params.get('status')! : undefined
          }
        ]}
        onChange={handleFilterChange}>
        {(
          strings: StringFilterProps[],
          types: TypeFilterProps[],
          checkboxes: CheckboxFilterProps[],
          radios: RadioFilterProps[]
        ) => {
          let checkboxesNotSet = checkboxes.every(c => !c.value!);
          const nameFilterValue: string = StringFilters.find(
            s => s.displayName === 'Filter By Name...'
          )!.value!;
          const radioFilter = radios[0].choice || params.get('status') || '';
          params.set('status', radioFilter);

          if (!virtualServicesList || isLoading) {
            return <div>Loading...</div>;
          }

          return (
            <div>
              <Route
                path={`${match.path}`}
                exact
                render={() => (
                  <>
                    {(checkboxesNotSet || checkboxes[0].value) && (
                      <SectionCard
                        data-testid='vs-listing-section'
                        cardName={'Virtual Services'}
                        logoIcon={<Gloo />}>
                        {!virtualServicesList.length && !isLoading ? (
                          <EmptyPrompt>
                            You don't have any virtual services.
                            <CreateVirtualServiceModal
                              promptText="Let's create one."
                              withoutDivider
                            />
                          </EmptyPrompt>
                        ) : (
                          <CardsListing
                            cardsData={getUsableCatalogData(
                              nameFilterValue,
                              virtualServicesList,
                              radioFilter
                            )}
                          />
                        )}
                      </SectionCard>
                    )}
                    {(checkboxesNotSet || checkboxes[1].value) && (
                      <SectionCard
                        data-testid='routetables-listing-section'
                        cardName={'Route Tables'}
                        logoIcon={<RouteTableIcon />}>
                        {!routeTablesList.length && !isLoading ? (
                          <EmptyPrompt>
                            You don't have any Route Tables.
                            <CreateRouteTableModal
                              promptText="Let's create one."
                              withoutDivider
                            />
                          </EmptyPrompt>
                        ) : (
                          <CardsListing
                            cardsData={getRouteTableCatalogData(
                              nameFilterValue,
                              routeTablesList,
                              radioFilter
                            )}
                          />
                        )}
                      </SectionCard>
                    )}
                  </>
                )}
              />
              <Route
                path={`${match.path}table`}
                exact
                render={() => (
                  <>
                    <SectionCard
                      noPadding
                      data-testid='vs-listing-section'
                      cardName={'Virtual Services'}
                      logoIcon={<Gloo />}>
                      <SoloTable
                        dataSource={getUsableTableData(
                          nameFilterValue,
                          virtualServicesList,
                          radioFilter
                        )}
                        columns={getTableColumns(
                          setRouteParentForRouteCreation,
                          deleteVS
                        )}
                      />
                    </SectionCard>
                    <SectionCard
                      noPadding
                      data-testid='routetables-listing-section'
                      cardName={'Route Tables'}
                      logoIcon={<RouteTableIcon />}>
                      <div style={{ padding: '-20px' }}>
                        <SoloTable
                          dataSource={getRouteTableData(
                            nameFilterValue,
                            routeTablesList,
                            radioFilter
                          )}
                          columns={getRouteTableColumns(
                            setRouteParentForRouteCreation,
                            deleteRT
                          )}
                        />
                      </div>
                    </SectionCard>
                  </>
                )}
              />
            </div>
          );
        }}
      </ListingFilter>
      <SoloModal
        visible={!!routeParentForRouteCreation}
        width={500}
        title={'Create Route'}
        onClose={() => setRouteParentForRouteCreation(undefined)}>
        <CreateRouteModal
          defaultRouteParent={routeParentForRouteCreation as RouteParent}
          completeCreation={() => setRouteParentForRouteCreation(undefined)}
        />
      </SoloModal>
    </div>
  );
};
