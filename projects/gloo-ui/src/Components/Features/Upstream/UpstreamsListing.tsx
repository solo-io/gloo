import * as React from 'react';
/** @jsx jsx */
import { jsx } from '@emotion/core';

import styled from '@emotion/styled/macro';
import { RouteComponentProps, Route, Switch } from 'react-router-dom';
import {
  healthConstants,
  TableHealthCircleHolder,
  TableActionCircle,
  TableActions
} from 'Styles';

import {
  ListingFilter,
  StringFilterProps,
  TypeFilterProps,
  CheckboxFilterProps,
  RadioFilterProps
} from '../../Common/ListingFilter';
import { CatalogTableToggle } from 'Components/Common/CatalogTableToggle';
import { Breadcrumb } from 'Components/Common/Breadcrumb';
import { useGetUpstreamsListV2 } from 'Api/v2/useUpstreamClientV2';
import {
  ListUpstreamsRequest,
  DeleteUpstreamRequest
} from '../../../proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/upstream_pb';
import { useGetUpstreamsList, useDeleteUpstream } from 'Api';
import { SectionCard } from 'Components/Common/SectionCard';
import { CardsListing } from 'Components/Common/CardsListing';
import { SoloTable } from 'Components/Common/SoloTable';
import { Upstream } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/upstream_pb';
import { Status } from 'proto/github.com/solo-io/solo-kit/api/v1/status_pb';
import {
  getResourceStatus,
  getUpstreamType,
  groupBy,
  getIcon,
  getFunctionInfo,
  CheckboxFilters,
  RadioFilters
} from 'utils/helpers';
import { NamespacesContext } from 'GlooIApp';
import { CreateUpstreamModal } from './Creation/CreateUpstreamModal';
import { HealthInformation } from 'Components/Common/HealthInformation';
import { HealthIndicator } from 'Components/Common/HealthIndicator';
import { SoloModal } from 'Components/Common/SoloModal';
import { CreateRouteModal } from '../Route/CreateRouteModal';

import { ExtraInfo } from 'Components/Features/Upstream/ExtraInfo';
import _ from 'lodash';

import { SuccessModal } from 'Components/Common/DisplayOnly/SuccessModal';
import { Popconfirm } from 'antd';
import { upstreams } from 'Api/v2/UpstreamClient';

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
  startCreatingRoute: (upstream: Upstream.AsObject) => any,
  deleteUpstream: (name: string, namespace: string) => void
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
      render: (us: Upstream.AsObject) => {
        return (
          <TableActions>
            <Popconfirm
              onConfirm={() =>
                deleteUpstream(us.metadata!.name, us.metadata!.namespace)
              }
              title={'Are you sure you want to delete this upstream? '}
              okText='Yes'
              cancelText='No'>
              <TableActionCircle>x</TableActionCircle>
            </Popconfirm>
            <TableActionCircle onClick={() => startCreatingRoute(us)}>
              +
            </TableActionCircle>
          </TableActions>
        );
      }
    }
  ];
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
  const [showSuccessModal, setShowSuccessModal] = React.useState(false);
  let params = new URLSearchParams(props.location.search);

  const [catalogNotTable, setCatalogNotTable] = React.useState(false);
  const [
    upstreamForRouteCreation,
    setUpstreamForRouteCreation
  ] = React.useState<Upstream.AsObject | undefined>(undefined);
  const namespaces = React.useContext(NamespacesContext);

  const [upstreamsList, setUpstreamsList] = React.useState<Upstream.AsObject[]>(
    []
  );

  const { data, loading, error, setNewVariables } = useGetUpstreamsListV2({
    namespaces: namespaces.namespacesList
  });

  React.useEffect(() => {
    if (data && data.toObject().upstreamsList) {
      setUpstreamsList(data.toObject().upstreamsList);
    }
  }, [loading]);

  React.useEffect(() => {
    if (props.location.state && props.location.state.showSuccess) {
      props.location.state.showSuccess = false;
    }
  }, []);

  React.useEffect(() => {
    if (props.location.state && props.location.state.showSuccess) {
      setShowSuccessModal(true);
      setTimeout(() => {
        setNewVariables({ namespaces: namespaces.namespacesList });
      }, 1000);
    }
    return () => setShowSuccessModal(false);
  }, [props.location.state && props.location.state.showSuccess]);

  if (!data || loading) {
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
    const selectedRadio = radios[0].choice || params.get('status') || '';
    params.set('status', selectedRadio);
    // group by type

    let upstreamsByType = groupBy(upstreamsList, u => getUpstreamType(u));
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
                groupBy(upstreams, u => u.metadata!.namespace).entries()
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
              columns={getTableColumns(
                setUpstreamForRouteCreation,
                deleteUpstream
              )}
            />
          )}
        />
      </div>
    );
  };

  const getUsableCatalogData = (
    nameFilter: string,
    data: Upstream.AsObject[],
    radioFilter: string
  ) => {
    const dataUsed: UpstreamCardData[] = data.map(upstream => {
      return {
        healthStatus: upstream.status
          ? upstream.status.state
          : healthConstants.Pending.value,
        cardTitle: upstream.metadata!.name,
        cardSubtitle: upstream.metadata!.namespace,
        onRemoveCard: () =>
          deleteUpstream(upstream.metadata!.name, upstream.metadata!.namespace),
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
        ...(!!getFunctionInfo(upstream) && {
          extraInfoComponent: () => <ExtraInfo upstream={upstream} />
        }),
        onCreate: () => setUpstreamForRouteCreation(upstream)
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
    data: Upstream.AsObject[],
    checkboxes: CheckboxFilterProps[],
    radioFilter: string
  ) => {
    const dataUsed = data.map(upstream => {
      return {
        ...upstream,
        status: upstream.status,
        type: getUpstreamType(upstream),
        name: upstream.metadata!.name,
        key: `${upstream.metadata!.name}-${upstream.metadata!.namespace}`,
        actions: upstream
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

  async function deleteUpstream(name: string, namespace: string) {
    await upstreams.deleteUpstream({ name, namespace });
    setUpstreamsList(usList => usList.filter(us => us.metadata!.name !== name));
  }
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
      <Heading>
        <Breadcrumb />
        <Action>
          <CreateUpstreamModal toggleSuccessModal={setShowSuccessModal} />
          <SuccessModal
            visible={showSuccessModal}
            successMessage='Upstream added successfully'
          />
          <CatalogTableToggle
            listIsSelected={catalogNotTable}
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
        radios={RadioFilters}
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
