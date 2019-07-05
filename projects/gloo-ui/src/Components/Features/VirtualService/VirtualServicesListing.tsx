import * as React from 'react';
/** @jsx jsx */
import { jsx } from '@emotion/core';

import styled from '@emotion/styled/macro';
import { withRouter, RouteComponentProps } from 'react-router';
import { colors } from 'Styles';
import {
  ListingFilter,
  StringFilterProps,
  TypeFilterProps,
  CheckboxFilterProps,
  RadioFilterProps
} from '../../Common/ListingFilter';
import { SoloTable } from 'Components/Common/SoloTable';
import { SectionCard } from 'Components/Common/SectionCard';
import { CatalogTableToggle } from 'Components/Common/CatalogTableToggle';
import { ReactComponent as Gloo } from 'assets/Gloo.svg';
import { Breadcrumb } from 'Components/Common/Breadcrumb';
import { CardsListing } from 'Components/Common/CardsListing';
import { CardType } from 'Components/Common/Card';
import { useListVirtualServices } from 'Api';
import { ListVirtualServicesRequest } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/virtualservice_pb';
import { VirtualService } from 'proto/github.com/solo-io/gloo/projects/gateway/api/v1/virtual_service_pb';
import { NamespacesContext } from 'GlooIApp';
import { getResourceStatus, getVSDomains } from 'utils/helpers';
import { CreateVirtualServiceModal } from './Creation/CreateVirtualServiceModal';

interface DataType {
  name: string;
  domain: string;
  namespace: string;
  version: string;
  status: number;
  routes: number;
  brLimit: {
    min: number;
    sec: number;
  };
  arLimit: {
    min: number;
    sec: number;
  };
  actions: null;
}
interface DataSourceType extends DataType {
  key: string;
}

const data: DataType[] = [
  {
    name: 'Jojoa.sdf',
    domain: 'abc.def',
    namespace: 'default',
    version: 'v012',
    status: 0,
    routes: 3,
    brLimit: {
      min: 5,
      sec: 0
    },
    arLimit: {
      min: 37,
      sec: 12
    },
    actions: null
  }
];

const TableColumns = [
  {
    title: 'Name',
    dataIndex: 'name'
  },
  {
    title: 'Domain',
    dataIndex: 'domains'
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
    dataIndex: 'status'
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
    render: (text: any) => <div>ACTION!</div>
  }
];

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
  let listVsRequest = React.useRef(new ListVirtualServicesRequest());
  const namespaces = React.useContext(NamespacesContext);

  listVsRequest.current.setNamespacesList(namespaces);
  const {
    data: vsListData,
    loading: vsLoading,
    error: vsError,
    refetch
  } = useListVirtualServices(listVsRequest.current);

  const [catalogNotTable, setCatalogNotTable] = React.useState<boolean>(true);
  const { history, match } = props;

  const getUsableCatalogData = (
    nameFilter: string,
    data: VirtualService.AsObject[]
  ) => {
    const dataUsed = data.map(virtualService => {
      return {
        ...virtualService,
        cardTitle: virtualService.metadata!.name,
        cardSubtitle: getVSDomains(virtualService),
        onRemovecard: (id: string): void => {},
        onExpanded: () => {},
        onClick: () => {
          history.push(
            `${match.path}${virtualService.metadata!.namespace}/${
              virtualService.metadata!.name
            }`
          );
        }
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
        name: virtualService.metadata!.name,
        domains: getVSDomains(virtualService),
        routes: virtualService.virtualHost!.routesList.length,
        status: getResourceStatus(virtualService),
        key: `${virtualService.metadata!.name}`
      };
    });

    return dataUsed.filter(row => row.name.includes(nameFilter));
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

    if (vsLoading) {
      return <div>Loading...</div>;
    }
    return (
      <div>
        {catalogNotTable ? (
          <SectionCard cardName={'Virtual Services'} logoIcon={<Gloo />}>
            <CardsListing
              cardsData={getUsableCatalogData(
                nameFilterValue,
                vsListData.virtualServicesList
              )}
            />
          </SectionCard>
        ) : (
          <SoloTable
            dataSource={getUsableTableData(
              nameFilterValue,
              vsListData.virtualServicesList
            )}
            columns={TableColumns}
          />
        )}
      </div>
    );
  };

  return (
    <div>
      <Heading>
        <Breadcrumb />
        <Action>
          <CreateVirtualServiceModal />
          <CatalogTableToggle
            listIsSelected={!catalogNotTable}
            onToggle={() => {
              setCatalogNotTable(cNt => !cNt);
            }}
          />
        </Action>
      </Heading>
      <ListingFilter strings={StringFilters} filterFunction={listDisplay} />
    </div>
  );
};
