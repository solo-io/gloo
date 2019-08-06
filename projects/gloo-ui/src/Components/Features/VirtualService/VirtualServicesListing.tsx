import * as React from 'react';
/** @jsx jsx */
import { jsx } from '@emotion/core';

import styled from '@emotion/styled/macro';
import { RouteComponentProps } from 'react-router';
import { colors, healthConstants } from 'Styles';
import {
  TableActionCircle,
  TableHealthCircleHolder,
  TableActions
} from 'Styles/table';
import {
  ListingFilter,
  StringFilterProps,
  TypeFilterProps,
  CheckboxFilterProps,
  RadioFilterProps
} from 'Components/Common/ListingFilter';
import { SoloTable } from 'Components/Common/SoloTable';
import { SectionCard } from 'Components/Common/SectionCard';
import { CatalogTableToggle } from 'Components/Common/CatalogTableToggle';
import { ReactComponent as Gloo } from 'assets/Gloo.svg';
import { Breadcrumb } from 'Components/Common/Breadcrumb';
import { CardsListing } from 'Components/Common/CardsListing';
import { useListVirtualServices, useDeleteVirtualService } from 'Api';
import {
  ListVirtualServicesRequest,
  DeleteVirtualServiceRequest
} from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/virtualservice_pb';
import { VirtualService } from 'proto/github.com/solo-io/gloo/projects/gateway/api/v1/virtual_service_pb';
import { Status } from 'proto/github.com/solo-io/solo-kit/api/v1/status_pb';
import { NamespacesContext } from 'GlooIApp';
import { getResourceStatus, getVSDomains } from 'utils/helpers';
import { CreateVirtualServiceModal } from './Creation/CreateVirtualServiceModal';
import { HealthInformation } from 'Components/Common/HealthInformation';
import { HealthIndicator } from 'Components/Common/HealthIndicator';
import { SoloModal } from 'Components/Common/SoloModal';
import { CreateRouteModal } from 'Components/Features/Route/CreateRouteModal';
import { ResourceRef } from 'proto/github.com/solo-io/solo-kit/api/v1/ref_pb';
import { Popconfirm } from 'antd';

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

const getTableColumns = (
  startCreatingRoute: (vs: VirtualService.AsObject) => any,
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
      render: (vs: VirtualService.AsObject) => {
        return (
          <TableActions>
            <Popconfirm
              onConfirm={() =>
                deleteVirtualService(vs.metadata!.name, vs.metadata!.namespace)
              }
              title={'Are you sure you want to delete this upstream? '}
              okText='Yes'
              cancelText='No'>
              <TableActionCircle>x</TableActionCircle>
            </Popconfirm>
            <TableActionCircle onClick={() => startCreatingRoute(vs)}>
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

interface Props extends RouteComponentProps {}

export const VirtualServicesListing = (props: Props) => {
  const { history, match } = props;
  let listVsRequest = React.useRef(new ListVirtualServicesRequest());
  const namespaces = React.useContext(NamespacesContext);

  listVsRequest.current.setNamespacesList(namespaces.namespacesList);
  const {
    data: vsListData,
    loading: vsLoading,
    error: vsError,
    refetch
  } = useListVirtualServices(listVsRequest.current);
  const { refetch: makeRequest } = useDeleteVirtualService(null);
  const [catalogNotTable, setCatalogNotTable] = React.useState(true);
  const [virtualServices, setVirtualServices] = React.useState<
    VirtualService.AsObject[]
  >([]);

  const [
    virtualServiceForRouteCreation,
    setVirtualServiceForRouteCreation
  ] = React.useState<VirtualService.AsObject | undefined>(undefined);
  React.useEffect(() => {
    if (vsListData) {
      setVirtualServices(vsListData.virtualServicesList);
    }
  }, [vsLoading]);

  const getUsableCatalogData = (
    nameFilter: string,
    data: VirtualService.AsObject[]
  ) => {
    const dataUsed = data.map(virtualService => {
      return {
        ...virtualService,
        healthStatus: virtualService.status
          ? virtualService.status.state
          : healthConstants.Pending.value,
        cardTitle: virtualService.displayName || virtualService.metadata!.name,
        cardSubtitle: getVSDomains(virtualService),
        onRemoveCard: () =>
          deleteVS(
            virtualService.metadata!.name,
            virtualService.metadata!.namespace
          ),
        removeConfirmText:
          'Are you sure you want to delete this virtual service?',
        onExpanded: () => {},
        onClick: () => {
          history.push(
            `${match.path}${virtualService.metadata!.namespace}/${
              virtualService.metadata!.name
            }`
          );
        },
        onCreate: () => setVirtualServiceForRouteCreation(virtualService)
      };
    });

    return dataUsed.filter(row => row.cardTitle.includes(nameFilter));
  };

  const getUsableTableData = (
    nameFilter: string,
    data: VirtualService.AsObject[]
  ) => {
    const dataUsed = data.map(virtualService => {
      return {
        ...virtualService,
        name: {
          displayName: virtualService.metadata!.name,
          goToVirtualService: () => {
            history.push(
              `${match.path}${virtualService.metadata!.namespace}/${
                virtualService.metadata!.name
              }`
            );
          }
        },
        domains: getVSDomains(virtualService),
        routes: virtualService.virtualHost!.routesList.length,
        status: virtualService.status,
        key: `${virtualService.metadata!.name}`,
        actions: virtualService
      };
    });

    return dataUsed.filter(row => row.name.displayName.includes(nameFilter));
  };

  const listDisplay = (
    strings: StringFilterProps[],
    types: TypeFilterProps[],
    checkboxes: CheckboxFilterProps[],
    radios: RadioFilterProps[]
  ) => {
    const nameFilterValue: string = strings.find(
      s => s.displayName === 'Filter By Name...'
    )!.value!;

    if (!vsListData || vsLoading) {
      return <div>Loading...</div>;
    }
    return (
      <div>
        {catalogNotTable ? (
          <SectionCard cardName={'Virtual Services'} logoIcon={<Gloo />}>
            <CardsListing
              cardsData={getUsableCatalogData(nameFilterValue, virtualServices)}
            />
          </SectionCard>
        ) : (
          <SoloTable
            dataSource={getUsableTableData(nameFilterValue, virtualServices)}
            columns={getTableColumns(
              setVirtualServiceForRouteCreation,
              deleteVS
            )}
          />
        )}
      </div>
    );
  };

  const finishCreation = (succeeded?: {
    namespace: string;
    name: string;
  }): void => {
    //TODO : Proper way to do this is to be polling always and, once we see the VS that matches this exists, we then jump

    if (succeeded) {
      setTimeout(() => {
        history.push(`${match.path}${succeeded.namespace}/${succeeded.name}`);
      }, 500);
    }
  };
  function deleteVS(name: string, namespace: string) {
    setVirtualServices(vsList =>
      vsList.filter(vs => vs.metadata!.name !== name)
    );
    let deleteReq = new DeleteVirtualServiceRequest();
    let ref = new ResourceRef();
    ref.setName(name);
    ref.setNamespace(namespace);
    deleteReq.setRef(ref);
    makeRequest(deleteReq);
  }
  return (
    <div>
      <Heading>
        <Breadcrumb />
        <Action>
          <CreateVirtualServiceModal finishCreation={finishCreation} />
          <CatalogTableToggle
            listIsSelected={!catalogNotTable}
            onToggle={() => {
              setCatalogNotTable(cNt => !cNt);
            }}
          />
        </Action>
      </Heading>
      <ListingFilter strings={StringFilters} filterFunction={listDisplay} />
      <SoloModal
        visible={!!virtualServiceForRouteCreation}
        width={500}
        title={'Create Route'}
        onClose={() => setVirtualServiceForRouteCreation(undefined)}>
        <CreateRouteModal
          defaultVirtualService={virtualServiceForRouteCreation}
          completeCreation={() => setVirtualServiceForRouteCreation(undefined)}
        />
      </SoloModal>
    </div>
  );
};
