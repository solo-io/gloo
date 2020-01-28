import styled from '@emotion/styled';
import { Popconfirm } from 'antd';
import { Breadcrumb } from 'Components/Common/Breadcrumb';
import { CardsListing } from 'Components/Common/CardsListing';
import { ListIcon, TileIcon } from 'Components/Common/CatalogTableToggle';
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
import { ExtraInfo } from 'Components/Features/Upstream/ExtraInfo';
import { Upstream } from 'proto/gloo/projects/gloo/api/v1/upstream_pb';
import { Status } from 'proto/solo-kit/api/v1/status_pb';
import { UpstreamDetails } from 'proto/solo-projects/projects/grpcserver/api/v1/upstream_pb';
import * as React from 'react';
import { useDispatch } from 'react-redux';
import { useHistory, useLocation, useRouteMatch } from 'react-router';
import { NavLink, Route } from 'react-router-dom';
import { deleteUpstream } from 'store/upstreams/actions';
import { upstreamAPI } from 'store/upstreams/api';
import {
  healthConstants,
  TableActionCircle,
  TableActions,
  TableHealthCircleHolder,
  colors
} from 'Styles';
import useSWR, { mutate } from 'swr';
import {
  CheckboxFilters,
  getFunctionInfo,
  getIcon,
  getResourceStatus,
  getUpstreamType,
  groupBy,
  RadioFilters
} from 'utils/helpers';
import { CreateRouteModal } from '../VirtualService/Creation/CreateRouteModal';
import { CreateUpstreamGroupModal } from './Creation/CreateUpstreamGroupModal';
import { SoloInput } from 'Components/Common/SoloInput';
import { UpstreamGroupDetails } from 'proto/solo-projects/projects/grpcserver/api/v1/upstreamgroup_pb';
import { ReactComponent as UpstreamGroupIcon } from 'assets/upstream-group-icon.svg';
import { ReactComponent as UpstreamIcon } from 'assets/upstream-icon-circle.svg';

import { upstreamGroupAPI } from 'store/upstreamGroups/api';
import { CreateUpstreamModal } from './Creation/CreateUpstreamModal';
import { UpstreamGroup } from 'proto/gloo/projects/gloo/api/v1/proxy_pb';
import { css } from '@emotion/core';

const UpstreamsListingContainer = styled.div`
  display: grid;
  grid-template-areas:
    'header header'
    'content content';
  grid-template-columns: 200px 1fr;
  grid-template-rows: auto 1fr;
  grid-column-gap: 20px;
`;

const EmptyPrompt = styled.div`
  display: flex;
  align-items: center;
  font-size: 14px;
`;
const HeaderSection = styled.div`
  grid-area: header;
  display: flex;
  justify-content: space-between;
  margin-bottom: 20px;
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
const TypeHolder = styled.div`
  display: flex;
  align-items: center;

  svg {
    width: 20px;
    height: 20px;
  }
`;

const StringFilters: StringFilterProps[] = [
  {
    displayName: 'Filter By Name...',
    placeholder: 'Filter by name...',
    value: ''
  }
];

const getUpstreamGroupTableColumns = () => {
  return [
    {
      title: 'Name',
      dataIndex: 'metadata.name',
      render: (name: string, resource: any) => (
        <>
          <NavLink
            css={css`
              cursor: pointer;
              color: ${colors.seaBlue};
            `}
            to={`/upstreams/upstreamgroups/${resource?.metadata?.namespace}/${resource?.metadata?.name}`}>
            {name}
          </NavLink>
        </>
      )
    },

    {
      title: 'Namespace',
      dataIndex: 'metadata.namespace'
    },
    {
      title: 'Upstreams',
      dataIndex: 'upstreamCount'
    },
    {
      title: 'Version',
      dataIndex: 'metadata.resourceVersion'
    },
    {
      title: 'Status',
      dataIndex: 'status',
      render: (status: any, resource: any) => {
        return (
          <div>
            <TableHealthCircleHolder>
              <HealthIndicator healthStatus={status.state} />
            </TableHealthCircleHolder>
            <HealthInformation healthStatus={status} />
          </div>
        );
      }
    },
    {
      title: 'Actions',
      dataIndex: 'actions',
      render: (upstreamGroupDetails: UpstreamGroupDetails.AsObject) => {
        const upstreamGroup = upstreamGroupDetails.upstreamGroup!;
        return (
          <TableActions>
            <Popconfirm
              onConfirm={() =>
                mutate(
                  [
                    'getUpstreamGroup',
                    upstreamGroup?.metadata?.name,
                    upstreamGroup?.metadata?.namespace
                  ],
                  upstreamGroupAPI.deleteUpstreamGroup({
                    ref: {
                      name: upstreamGroup?.metadata?.name!,
                      namespace: upstreamGroup?.metadata?.namespace!
                    }
                  })
                )
              }
              title={'Are you sure you want to delete this upstream group ? '}
              okText='Yes'
              cancelText='No'>
              <TableActionCircle>x</TableActionCircle>
            </Popconfirm>
            {!!upstreamGroupDetails.raw && (
              <FileDownloadActionCircle
                fileContent={upstreamGroupDetails.raw.content}
                fileName={upstreamGroupDetails.raw.fileName}
              />
            )}
            <TableActionCircle onClick={() => {}}>+</TableActionCircle>
          </TableActions>
        );
      }
    }
  ];
};

const getTableColumns = (
  startCreatingRoute: (upstream: Upstream.AsObject) => any
) => {
  return [
    {
      title: 'Name',
      dataIndex: 'name',
      render: (name: string, resource: any) => (
        <>
          <NavLink
            css={css`
              cursor: pointer;
              color: ${colors.seaBlue};
            `}
            to={`/upstreams/${resource?.metadata?.namespace}/${resource?.metadata?.name}`}>
            {name}
          </NavLink>
        </>
      )
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
      title: 'Type',
      dataIndex: 'type',
      render: (upstreamType: string) => (
        <TypeHolder>
          {getIcon(upstreamType)}
          <span style={{ marginLeft: '5px' }}>{upstreamType}</span>
        </TypeHolder>
      )
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
      title: 'Use TLS',
      dataIndex: 'useTls'
    },

    {
      title: 'Actions',
      dataIndex: 'actions',
      render: (upstreamDetails: UpstreamDetails.AsObject) => (
        <TableAction
          upstreamDetails={upstreamDetails}
          startCreatingRoute={startCreatingRoute}
        />
      )
    }
  ];
};

type TableActionProps = {
  upstreamDetails: UpstreamDetails.AsObject;
  startCreatingRoute: any;
};
const TableAction: React.FC<TableActionProps> = props => {
  const { upstreamDetails, startCreatingRoute } = props;
  const dispatch = useDispatch();
  return (
    <TableActions>
      <Popconfirm
        onConfirm={() =>
          dispatch(
            deleteUpstream({
              ref: {
                name: upstreamDetails.upstream!.metadata!.name,
                namespace: upstreamDetails.upstream!.metadata!.namespace
              }
            })
          )
        }
        title={'Are you sure you want to delete this upstream? '}
        okText='Yes'
        cancelText='No'>
        <TableActionCircle>x</TableActionCircle>
      </Popconfirm>
      {!!upstreamDetails.raw && (
        <FileDownloadActionCircle
          fileContent={upstreamDetails.raw.content}
          fileName={upstreamDetails.raw.fileName}
        />
      )}
      <TableActionCircle
        onClick={() => startCreatingRoute(upstreamDetails.upstream)}>
        +
      </TableActionCircle>
    </TableActions>
  );
};
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
interface UpstreamCardData {
  cardTitle: string;
  cardSubtitle: { key: string; value: string }[];
  onRemoveCard?: () => any;
  onExpand: () => void;
  details: {
    title: string;
    value: string;
    valueDisplay?: React.ReactNode;
  }[];
  healthStatus: number;
}

/* ----------------------------------------------------------------------------------------------- */

export const UpstreamsListing = () => {
  let location = useLocation();
  let match = useRouteMatch({
    path: '/upstreams/'
  })!;
  let history = useHistory();

  // filtering
  let params = new URLSearchParams(location.search);
  const [filterString, setFilterString] = React.useState('');
  let currentStatusToShow = params.get('status') || 'Accepted';

  const dispatch = useDispatch();
  const { data: upstreamsList, error } = useSWR(
    'listUpstreams',
    upstreamAPI.listUpstreams
  );
  const { data: upstreamGroupsList, error: upstreamGroupError } = useSWR(
    'listUpstreamGroups',
    upstreamGroupAPI.listUpstreamGroups
  );

  const [
    upstreamForRouteCreation,
    setUpstreamForRouteCreation
  ] = React.useState<Upstream.AsObject | UpstreamGroup.AsObject | undefined>(
    undefined
  );

  React.useEffect(() => {
    if (location.state && location.state.showSuccess) {
      location.state.showSuccess = false;
    }
  }, []);

  if (!upstreamsList || !upstreamGroupsList) {
    return <div>Loading...</div>;
  }

  const getUsableCatalogData = (
    nameFilter: string,
    data: UpstreamDetails.AsObject[],
    radioFilter: string
  ) => {
    const dataUsed: UpstreamCardData[] = data.map(upstreamDet => {
      const upstream = upstreamDet.upstream!;

      return {
        healthStatus: upstream.status
          ? upstream.status.state
          : healthConstants.Pending.value,
        cardTitle: upstream.metadata!.name,
        cardSubtitle: [
          { key: 'Namespace', value: upstream.metadata!.namespace }
        ],
        onRemoveCard: () =>
          dispatch(
            deleteUpstream({
              ref: {
                name: upstream.metadata!.name,
                namespace: upstream.metadata!.namespace
              }
            })
          ),
        removeConfirmText: 'Are you sure you want to delete this upstream?',
        onExpand: () => {},
        details: [
          {
            title: 'Name',
            value: upstream.metadata!.name
          },
          {
            title: 'Namespace',
            value: upstream.metadata!.namespace
          },
          {
            title: 'Version',
            value: upstream.metadata!.resourceVersion
          },
          {
            title: 'Type',
            value: getUpstreamType(upstream)
          },

          {
            title: 'Status',
            value: getResourceStatus(upstream),
            valueDisplay: <HealthInformation healthStatus={upstream.status} />
          },
          ...(!!getFunctionInfo(upstream)
            ? [
                {
                  title: 'Functions',
                  value: getFunctionInfo(upstream)
                }
              ]
            : [])
        ],
        onClick: () => {
          history.push({
            pathname: `${match.url}${upstream?.metadata?.namespace}/${upstream?.metadata?.name}`
          });
        },
        ExtraInfoComponent: <ExtraInfo upstream={upstream} />,
        onCreate: () => setUpstreamForRouteCreation(upstream),
        downloadableContent: upstreamDet.raw
      };
    });

    return dataUsed
      .filter(row => row.cardTitle.includes(nameFilter))
      .filter(row =>
        getResourceStatus(row.healthStatus)
          .toLowerCase()
          .includes(radioFilter.toLowerCase())
      );
  };

  const getUsableTableData = (
    nameFilter: string,
    data: UpstreamDetails.AsObject[],
    checkboxes: CheckboxFilterProps[],
    radioFilter: string
  ) => {
    const dataUsed = data.map(upstreamDet => {
      const upstream = upstreamDet.upstream!;

      return {
        ...upstream,
        status: upstream.status,
        type: getUpstreamType(upstream),
        name: upstream.metadata!.name,
        key: `${upstream.metadata!.name}-${upstream.metadata!.namespace}`,
        actions: upstreamDet
      };
    });
    let checkboxesNotSet = checkboxes.every(c => !c.value!);

    return dataUsed
      .filter(row => row.name.includes(nameFilter))
      .filter(row =>
        getResourceStatus(row)
          .toLowerCase()
          .includes(radioFilter.toLowerCase())
      )
      .filter(row => {
        return (
          checkboxes.find(c => c.displayName === row.type)!.value! ||
          checkboxesNotSet
        );
      });
  };

  function formatUpstreamGroupData(
    data: UpstreamGroupDetails.AsObject[],
    nameFilter: string,
    checkboxes: CheckboxFilterProps[],
    radioFilter: string
  ) {
    const dataUsed = data.map(upstreamGroupDetail => {
      let upstreamGroup = upstreamGroupDetail.upstreamGroup;

      return {
        ...upstreamGroup,
        removeConfirmText:
          'Are you sure you want to delete this Upstream Group?',
        onCreate: () => setUpstreamForRouteCreation(upstreamGroup),
        downloadableContent: upstreamGroupDetail.raw!,
        onRemoveCard: () =>
          mutate(
            [
              'getUpstreamGroup',
              upstreamGroup?.metadata?.name,
              upstreamGroup?.metadata?.namespace
            ],
            upstreamGroupAPI.deleteUpstreamGroup({
              ref: {
                name: upstreamGroup?.metadata?.name!,
                namespace: upstreamGroup?.metadata?.namespace!
              }
            })
          ),
        healthStatus: upstreamGroup?.status
          ? upstreamGroup.status.state
          : healthConstants.Pending.value,
        cardTitle: upstreamGroup?.metadata?.name!,
        cardSubtitle: [
          { key: 'Namespace', value: upstreamGroup?.metadata!.namespace }
        ],
        onClick: () => {
          history.push({
            pathname: `${match.url}upstreamgroups/${upstreamGroup?.metadata?.namespace}/${upstreamGroup?.metadata?.name}`
          });
        }
      };
    });
    let checkboxesNotSet = checkboxes.every(c => !c.value!);

    return dataUsed
      .filter(row => row.cardTitle.includes(nameFilter))
      .filter(row =>
        getResourceStatus(row.healthStatus)
          .toLowerCase()
          .includes(radioFilter.toLowerCase())
      )
      .filter(row => {
        return (
          checkboxes.find(c => c.displayName === row.cardTitle)?.value ||
          checkboxesNotSet
        );
      });
  }

  function formatUpstreamGroupTableData(
    data: UpstreamGroupDetails.AsObject[],
    nameFilter: string,
    checkboxes: CheckboxFilterProps[],
    radioFilter: string
  ) {
    const dataUsed = data.map(upstreamGroupDetail => {
      let upstreamGroup = upstreamGroupDetail.upstreamGroup;
      return {
        ...upstreamGroup,
        status: upstreamGroup?.status!,
        name: upstreamGroup?.metadata?.name!,
        namespace: upstreamGroup?.metadata?.namespace!,
        upstreamCount: upstreamGroup?.destinationsList.length,
        actions: upstreamGroupDetail
      };
    });
    let checkboxesNotSet = checkboxes.every(c => !c.value!);

    return dataUsed
      .filter(row => row.name.includes(nameFilter))
      .filter(row =>
        getResourceStatus(row.status.state)
          .toLowerCase()
          .includes(radioFilter.toLowerCase())
      )
      .filter(row => {
        return (
          checkboxes.find(c => c.displayName === row.name)?.value ||
          checkboxesNotSet
        );
      });
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

  let cardMatch = useRouteMatch({
    path: `${match.path}`,
    exact: true
  });

  let tableMatch = useRouteMatch({
    path: `${match.path}table`,
    exact: true
  });

  return (
    <UpstreamsListingContainer>
      <HeaderSection>
        <Breadcrumb />
        <Action>
          <CreateUpstreamModal />
          <CreateUpstreamGroupModal />

          <NavLink to={{ pathname: match.path, search: location.search }}>
            <TileIcon selected={!location.pathname.includes('table')} />
          </NavLink>
          <NavLink
            to={{ pathname: `${match.path}table`, search: location.search }}>
            <ListIcon selected={location.pathname.includes('table')} />
          </NavLink>
        </Action>
      </HeaderSection>
      <SidebarSection></SidebarSection>
      <ContentSection>
        <ListingFilter
          showLabels
          strings={StringFilters}
          checkboxes={CheckboxFilters}
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
            const nameFilterValue: string = strings.find(
              s => s.displayName === 'Filter By Name...'
            )!.value!;
            const selectedRadio = params.get('status') || '';
            params.set('status', selectedRadio);
            // group by type

            let upstreamsByType = React.useMemo(
              () => groupBy(upstreamsList, u => getUpstreamType(u.upstream!)),
              [upstreamsList.length]
            );
            let upstreamsByTypeArr = Array.from(upstreamsByType.entries());
            let checkboxesNotSet = checkboxes.every(c => !c.value!);
            return (
              <div>
                {cardMatch && (
                  <>
                    {!!upstreamGroupsList.length && (
                      <SectionCard
                        key='upstreamgroups-listing-section'
                        data-testid='upstreamgroups-listing-section'
                        cardName={'Upstream Groups'}
                        logoIcon={<UpstreamGroupIcon />}>
                        {!upstreamGroupsList.length ? (
                          <EmptyPrompt>
                            You don't have any Upstream Groups.
                            <CreateUpstreamGroupModal />
                          </EmptyPrompt>
                        ) : (
                          <CardsListing
                            cardsData={formatUpstreamGroupData(
                              upstreamGroupsList,
                              nameFilterValue,
                              checkboxes,
                              selectedRadio
                            )}
                          />
                        )}
                      </SectionCard>
                    )}
                    {upstreamsByTypeArr.map(([type, upstreams]) => {
                      // show section according to type filter
                      let groupedByNamespaces = Array.from(
                        groupBy(
                          upstreams,
                          u => u.upstream!.metadata!.namespace
                        ).entries()
                      );

                      if (
                        checkboxesNotSet ||
                        checkboxes.find(c => c.displayName === type)!.value!
                      ) {
                        const cardListingsData = groupedByNamespaces
                          .map(([namespace, upstreams]) => {
                            return {
                              namespace,
                              cardsData: getUsableCatalogData(
                                nameFilterValue,
                                upstreams,
                                selectedRadio
                              )
                            };
                          })
                          .filter(data => !!data.cardsData.length);

                        if (!cardListingsData.length) {
                          return null;
                        }

                        return (
                          <SectionCard
                            cardName={type}
                            logoIcon={getIcon(type)}
                            key={type}>
                            {cardListingsData.map(data => (
                              <CardsListing
                                key={data.namespace}
                                title={data.namespace}
                                cardsData={data.cardsData}
                              />
                            ))}
                          </SectionCard>
                        );
                      }
                    })}
                  </>
                )}
                {tableMatch && (
                  <>
                    {upstreamGroupsList.length > 0 && (
                      <SectionCard
                        noPadding
                        data-testid='upstreamgroups-listing-section'
                        cardName={'Upstream Groups'}
                        logoIcon={<UpstreamGroupIcon />}>
                        <SoloTable
                          dataSource={formatUpstreamGroupTableData(
                            upstreamGroupsList,
                            nameFilterValue,
                            checkboxes,
                            selectedRadio
                          )}
                          columns={getUpstreamGroupTableColumns()}
                        />
                      </SectionCard>
                    )}
                    <SectionCard
                      noPadding
                      data-testid='upstreams-listing-section'
                      cardName={'Upstreams'}
                      logoIcon={<UpstreamIcon />}>
                      <div style={{ padding: '-20px' }}>
                        <SoloTable
                          dataSource={getUsableTableData(
                            nameFilterValue,
                            upstreamsList,
                            checkboxes,
                            selectedRadio
                          )}
                          columns={getTableColumns(setUpstreamForRouteCreation)}
                        />
                      </div>
                    </SectionCard>
                  </>
                )}
              </div>
            );
          }}
        </ListingFilter>
      </ContentSection>
      <SoloModal
        visible={!!upstreamForRouteCreation}
        width={500}
        title={'Create Route'}
        onClose={() => setUpstreamForRouteCreation(undefined)}>
        <CreateRouteModal
          defaultUpstream={upstreamForRouteCreation}
          completeCreation={() => setUpstreamForRouteCreation(undefined)}
        />
      </SoloModal>
    </UpstreamsListingContainer>
  );
};
