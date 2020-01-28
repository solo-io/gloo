import styled from '@emotion/styled';
import { Popconfirm } from 'antd';
import { ReactComponent as Gloo } from 'assets/Gloo.svg';
import { ReactComponent as RouteTableIcon } from 'assets/route-table-icon.svg';
import { Breadcrumb } from 'Components/Common/Breadcrumb';
import { CardsListing } from 'Components/Common/CardsListing';
import { ListIcon, TileIcon } from 'Components/Common/CatalogTableToggle';
import { FileDownloadActionCircle } from 'Components/Common/FileDownloadLink';
import { HealthIndicator } from 'Components/Common/HealthIndicator';
import { HealthInformation } from 'Components/Common/HealthInformation';
import { StyledHeader } from 'Components/Common/ListingFilter';
import { SectionCard } from 'Components/Common/SectionCard';
import { SoloCheckbox } from 'Components/Common/SoloCheckbox';
import { SoloInput } from 'Components/Common/SoloInput';
import { SoloModal } from 'Components/Common/SoloModal';
import { SoloRadioGroup } from 'Components/Common/SoloRadioGroup';
import { SoloTable } from 'Components/Common/SoloTable';
import { Status } from 'proto/solo-kit/api/v1/status_pb';
import { RouteTableDetails } from 'proto/solo-projects/projects/grpcserver/api/v1/routetable_pb';
import { VirtualServiceDetails } from 'proto/solo-projects/projects/grpcserver/api/v1/virtualservice_pb';
import React, { useState } from 'react';
import { useDispatch } from 'react-redux';
import {
  NavLink,
  useHistory,
  useLocation,
  useRouteMatch
} from 'react-router-dom';
import { deleteRouteTable } from 'store/routeTables/actions';
import { routeTableAPI } from 'store/routeTables/api';
import { deleteVirtualService } from 'store/virtualServices/actions';
import { virtualServiceAPI } from 'store/virtualServices/api';
import { colors, healthConstants } from 'Styles';
import {
  TableActionCircle,
  TableActions,
  TableHealthCircleHolder
} from 'Styles/table';
import useSWR from 'swr';
import { getResourceStatus, getVSDomains, RadioFilters } from 'utils/helpers';
import { CreateRouteModal } from './Creation/CreateRouteModal';
import { CreateRouteTableModal } from './Creation/CreateRouteTableModal';
import { CreateVirtualServiceModal } from './Creation/CreateVirtualServiceModal';
import { RouteParent } from './RouteTableDetails';
import { css } from '@emotion/core';

const FilterHeader = styled.div`
  ${StyledHeader};
  width: 185px;
`;

const VSListingContainer = styled.div`
  display: grid;
  grid-template-areas:
    'header header'
    'sidebar content';
  grid-template-columns: 200px 1fr;
  grid-template-rows: auto 1fr;
  grid-column-gap: 20px;
`;
const HeaderSection = styled.div`
  grid-area: header;
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

const SidebarSection = styled.div`
  grid-area: sidebar;
`;

const ContentSection = styled.div`
  grid-area: content;
`;
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

export const VirtualServicesListing = () => {
  const { data: virtualServicesList, error: listVirtualServicesError } = useSWR(
    'listVirtualServices',
    virtualServiceAPI.listVirtualServices
  );
  const { data: routeTablesList, error: listRouteTablesError } = useSWR(
    'listRouteTables',
    routeTableAPI.listRouteTables
  );
  let history = useHistory();
  let location = useLocation();
  let match = useRouteMatch({
    path: '/virtualservices/'
  })!;

  // filtering
  let params = new URLSearchParams(location.search);
  let currentStatusToShow = params.get('status') || 'Accepted';
  const [showVS, setShowVS] = useState(false);
  const [showRT, setShowRT] = useState(false);
  const [filterString, setFilterString] = useState('');

  // redux
  const dispatch = useDispatch();
  const [isLoading, setIsLoading] = useState(false);

  const [
    routeParentForRouteCreation,
    setRouteParentForRouteCreation
  ] = useState<RouteParent | undefined>(undefined);

  interface VSorRTData
    extends VirtualServiceDetails.AsObject,
      RouteTableDetails.AsObject {}

  const formatVSCatalogData = (data: VSorRTData[]) => {
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
    return (
      dataUsed
        // filter by search bar
        .filter(row =>
          row.cardTitle.toLowerCase().includes(filterString.toLowerCase())
        )
        // filter by status from query params
        .filter(row => {
          if (params.get('status') === null) {
            return true;
          }
          getResourceStatus(row).includes(currentStatusToShow);
        })
        // filter by checkbox
        .filter(row => {
          if ((!showVS && !showRT) || showVS) return row;
          else return;
        })
    );
  };

  const formatRTCatalogData = (
    routeTableDetailsList: RouteTableDetails.AsObject[]
  ) => {
    const catalogData = routeTableDetailsList.map(routeTableDetails => {
      let routeTable = routeTableDetails.routeTable!;

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
    return (
      catalogData
        // filter by search bar
        .filter(row =>
          row.cardTitle.toLowerCase().includes(filterString.toLowerCase())
        )
        // filter by status from query params
        .filter(row => {
          if (params.get('status') === null) {
            return true;
          }
          getResourceStatus(row).includes(currentStatusToShow);
        })
        .filter(row => {
          if ((!showVS && !showRT) || showRT) return row;
          else return;
        })
    );
  };

  const formatRTTableData = (data: VSorRTData[]) => {
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

    return (
      dataUsed
        .filter(row =>
          row.name.displayName.includes(filterString.toLowerCase())
        )
        // filter by status from query params
        .filter(row => {
          if (params.get('status') === null) {
            return true;
          } else {
            return getResourceStatus(row).includes(currentStatusToShow);
          }
        })
        .filter(row => {
          if ((!showVS && !showRT) || showRT) return row;
          else return;
        })
    );
  };

  const formatVSTableData = (data: VSorRTData[]) => {
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

    return (
      dataUsed
        .filter(row =>
          row.name.displayName.includes(filterString.toLowerCase())
        )
        // filter by status from query params
        .filter(row => {
          if (params.get('status') === null) {
            return true;
          } else {
            return getResourceStatus(row).includes(currentStatusToShow);
          }
        })
        .filter(row => {
          if ((!showVS && !showRT) || showVS) return row;
          else return;
        })
    );
  };

  function deleteVS(name: string, namespace: string) {
    dispatch(deleteVirtualService({ ref: { name, namespace } }));
  }

  function deleteRT(name: string, namespace: string) {
    dispatch(deleteRouteTable({ ref: { name, namespace } }));
  }

  let cardMatch = useRouteMatch({
    path: `${match.path}`,
    exact: true
  });

  let tableMatch = useRouteMatch({
    path: `${match.path}table`,
    exact: true
  });

  if (!virtualServicesList || !routeTablesList) {
    return <div>Loading...</div>;
  }

  return (
    <VSListingContainer>
      <HeaderSection>
        <Breadcrumb />
        <Action>
          <CreateVirtualServiceModal />
          <CreateRouteTableModal />
          <NavLink to={{ pathname: match.path, search: location.search }}>
            <TileIcon selected={!location.pathname.includes('table')} />
          </NavLink>
          <NavLink
            to={{ pathname: `${match.path}table`, search: location.search }}>
            <ListIcon selected={location.pathname.includes('table')} />
          </NavLink>
        </Action>
      </HeaderSection>
      <SidebarSection>
        <SoloInput
          value={filterString}
          placeholder={'Filter by name...'}
          onChange={({ target }) => {
            setFilterString(target.value);
          }}
        />
        <FilterHeader>Status Filter</FilterHeader>
        <SoloRadioGroup
          options={RadioFilters.options.map(option => {
            return {
              displayName: option.displayName,
              id: option.id || option.displayName
            };
          })}
          currentSelection={params.get('status')!}
          onChange={newValue => {
            if (newValue !== undefined) {
              history.push(`/virtualservices/?status=${newValue}`);
            } else {
              history.push('/virtualservices/');
            }
          }}
        />
        <FilterHeader>Types Filter</FilterHeader>
        <SoloCheckbox
          title={'Virtual Services'}
          checked={showVS}
          withWrapper={true}
          onChange={evt => {
            setShowVS(evt.target.checked);
          }}
        />
        <SoloCheckbox
          title={'Route Tables'}
          checked={showRT}
          withWrapper={true}
          onChange={evt => {
            setShowRT(evt.target.checked);
          }}
        />
      </SidebarSection>

      <ContentSection>
        {cardMatch && (
          <>
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
                  cardsData={formatVSCatalogData(virtualServicesList)}
                />
              )}
            </SectionCard>

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
                  cardsData={formatRTCatalogData(routeTablesList)}
                />
              )}
            </SectionCard>
          </>
        )}
        {tableMatch && (
          <>
            <SectionCard
              noPadding
              data-testid='vs-listing-section'
              cardName={'Virtual Services'}
              logoIcon={<Gloo />}>
              <SoloTable
                dataSource={formatVSTableData(virtualServicesList)}
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
                  dataSource={formatRTTableData(routeTablesList)}
                  columns={getRouteTableColumns(
                    setRouteParentForRouteCreation,
                    deleteRT
                  )}
                />
              </div>
            </SectionCard>
          </>
        )}
      </ContentSection>

      <SoloModal
        visible={!!routeParentForRouteCreation}
        width={500}
        title={'Create Route'}
        onClose={() => setRouteParentForRouteCreation(undefined)}>
        <CreateRouteModal
          defaultRouteParent={routeParentForRouteCreation}
          completeCreation={() => setRouteParentForRouteCreation(undefined)}
        />
      </SoloModal>
    </VSListingContainer>
  );
};
