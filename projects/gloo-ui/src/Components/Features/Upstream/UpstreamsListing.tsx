import styled from '@emotion/styled';
import { Popconfirm } from 'antd';
import { Breadcrumb } from 'Components/Common/Breadcrumb';
import { CardsListing } from 'Components/Common/CardsListing';
import { CatalogTableToggle } from 'Components/Common/CatalogTableToggle';
import { SuccessModal } from 'Components/Common/DisplayOnly/SuccessModal';
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
import { Upstream } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/upstream_pb';
import { Status } from 'proto/github.com/solo-io/solo-kit/api/v1/status_pb';
import { UpstreamDetails } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/upstream_pb';
import * as React from 'react';
import { useDispatch, useSelector, shallowEqual } from 'react-redux';
import { LoadingBar } from 'react-redux-loading-bar';
import { Route, RouteComponentProps } from 'react-router-dom';
import { AppState } from 'store';
import { deleteUpstream, listUpstreams } from 'store/upstreams/actions';
import {
  healthConstants,
  TableActionCircle,
  TableActions,
  TableHealthCircleHolder
} from 'Styles';
import {
  CheckboxFilters,
  getFunctionInfo,
  getIcon,
  getResourceStatus,
  getUpstreamType,
  groupBy,
  RadioFilters
} from 'utils/helpers';
import { CreateRouteModal } from '../Route/CreateRouteModal';
import { CreateUpstreamModal } from './Creation/CreateUpstreamModal';

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

const getTableColumns = (
  startCreatingRoute: (upstream: Upstream.AsObject) => any
) => {
  return [
    {
      title: 'Name',
      dataIndex: 'name'
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
      title: 'Routes',
      dataIndex: 'routes'
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
  cardSubtitle: string;
  onRemoveCard?: () => any;
  onExpand: () => void;
  details: {
    title: string;
    value: string;
    valueDisplay?: React.ReactNode;
  }[];
  healthStatus: number;
}
interface Props extends RouteComponentProps {
  //... eg, virtualservice?: string
}

export const UpstreamsListing = (props: Props) => {
  const dispatch = useDispatch();
  const [isLoading, setIsLoading] = React.useState(false);
  const upstreamsList = useSelector(
    (state: AppState) => state.upstreams.upstreamsList
  );

  React.useEffect(() => {
    if (upstreamsList.length) {
      setIsLoading(false);
    } else {
      dispatch(listUpstreams());
    }
  }, [upstreamsList.length]);

  let params = new URLSearchParams(props.location.search);

  const [catalogNotTable, setCatalogNotTable] = React.useState(
    !props.location.pathname.includes('table')
  );
  const [
    upstreamForRouteCreation,
    setUpstreamForRouteCreation
  ] = React.useState<Upstream.AsObject | undefined>(undefined);

  React.useEffect(() => {
    if (props.location.state && props.location.state.showSuccess) {
      props.location.state.showSuccess = false;
    }
  }, []);

  if (isLoading) {
    return <div>Loading...</div>;
  }

  const listDisplay = (
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

    let upstreamsByType = groupBy(upstreamsList, u =>
      getUpstreamType(u.upstream!)
    );
    let upstreamsByTypeArr = Array.from(upstreamsByType.entries());
    let checkboxesNotSet = checkboxes.every(c => !c.value!);
    return (
      <div>
        <Route
          path={props.match.path}
          exact
          render={() =>
            upstreamsByTypeArr.map(([type, upstreams]) => {
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
            })
          }
        />
        <Route
          path={`${props.match.path}table`}
          render={() => (
            <SoloTable
              dataSource={getUsableTableData(
                nameFilterValue,
                upstreamsList,
                checkboxes,
                selectedRadio
              )}
              columns={getTableColumns(setUpstreamForRouteCreation)}
            />
          )}
        />
      </div>
    );
  };

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
        cardSubtitle: upstream.metadata!.namespace,
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

  function handleFilterChange(
    strings: StringFilterProps[],
    types: TypeFilterProps[],
    checkboxes: CheckboxFilterProps[],
    radios: RadioFilterProps[]
  ) {
    params.set('status', radios[0].choice || '');
    props.history.replace({
      pathname: `${props.location.pathname}`,
      search: radios[0].choice
        ? `?${'status'}=${radios[0].choice}`
        : params.get('status') || ''
    });
  }

  return (
    <div>
      <LoadingBar />
      <Heading>
        <Breadcrumb />
        <Action>
          <CreateUpstreamModal />
          <CatalogTableToggle
            listIsSelected={!catalogNotTable}
            onToggle={() => {
              const { location, match, history } = props;

              history.push({
                pathname: `${match.path}${
                  location.pathname.includes('table') ? '' : 'table'
                }`,
                search: location.search
              });
              setCatalogNotTable(cNt => !cNt);
            }}
          />
        </Action>
      </Heading>
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
        onChange={handleFilterChange}
        filterFunction={listDisplay}
      />
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
    </div>
  );
};
