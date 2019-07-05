import * as React from 'react';
/** @jsx jsx */
import { jsx } from '@emotion/core';

import styled from '@emotion/styled/macro';
import { withRouter, RouteComponentProps } from 'react-router';
import { colors } from 'Styles';
import { ReactComponent as Gloo } from 'assets/Gloo.svg';

import {
  ListingFilter,
  StringFilterProps,
  TypeFilterProps,
  CheckboxFilterProps,
  RadioFilterProps
} from '../../Common/ListingFilter';
import { CatalogTableToggle } from 'Components/Common/CatalogTableToggle';
import { Breadcrumb } from 'Components/Common/Breadcrumb';
import { ListUpstreamsRequest } from '../../../proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/upstream_pb';
import { useGetUpstreamsList } from 'Api';
import { SectionCard } from 'Components/Common/SectionCard';
import { CardsListing } from 'Components/Common/CardsListing';
import { SoloTable } from 'Components/Common/SoloTable';
import { CardType } from 'antd/lib/card';
import { Upstream } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/upstream_pb';
import { Status } from 'proto/github.com/solo-io/solo-kit/api/v1/status_pb';
import { getResourceStatus, getUpstreamType, groupBy } from 'utils/helpers';
import { NamespacesContext } from 'GlooIApp';
import { CreateUpstreamModal } from './Creation/CreateUpstreamModal';

const StringFilters: StringFilterProps[] = [
  {
    displayName: 'Filter By Name...',
    placeholder: 'Filter by name...',
    value: ''
  }
];

const TableColumns = [
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
    title: 'Routes',
    dataIndex: 'routes'
  },
  {
    title: 'Status',
    dataIndex: 'status'
  },
  {
    title: 'Use TLS',
    dataIndex: 'useTls'
  },

  {
    title: 'Actions',
    dataIndex: 'actions',
    render: (text: any) => <div>ACTION!</div>
  }
];

const CheckboxFilters: CheckboxFilterProps[] = [
  {
    displayName: 'AWS',
    value: false
  },
  {
    displayName: 'Azure',
    value: false
  },
  {
    displayName: 'REST',
    value: false
  },
  {
    displayName: 'gRPC',
    value: false
  },
  {
    displayName: 'Kubernetes',
    value: false
  },
  {
    displayName: 'Static',
    value: false
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
interface UpstreamCardData extends Upstream.AsObject {
  cardTitle: string;
  cardSubtitle: string;
  onRemovecard: (id: string) => void;
  onExpanded: () => void;
  details: { title: string; value: string }[];
}
interface Props extends RouteComponentProps {
  //... eg, virtualservice?: string
}

export const UpstreamsListing = (props: Props) => {
  const [catalogNotTable, setCatalogNotTable] = React.useState<boolean>(true);
  const namespaces = React.useContext(NamespacesContext);
  let request = new ListUpstreamsRequest();
  request.setNamespacesList(namespaces);
  const { data, loading, error } = useGetUpstreamsList(request);

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

    // group by type
    let upstreamsByType = groupBy(data.upstreamsList, u => getUpstreamType(u));
    let upstreamsByTypeArr = Array.from(upstreamsByType.entries());

    let checkboxesNotSet = checkboxes.every(c => !c.value!);
    return (
      <div>
        {catalogNotTable ? (
          upstreamsByTypeArr.map(
            ([type, upstreams]) =>
              // show section according to type filter
              (checkboxesNotSet ||
                checkboxes.find(c => c.displayName === type)!.value!) && (
                <SectionCard cardName={type} logoIcon={<Gloo />} key={type}>
                  <CardsListing
                    cardsData={getUsableCatalogData(nameFilterValue, upstreams)}
                  />
                </SectionCard>
              )
          )
        ) : (
          <SoloTable
            dataSource={getUsableTableData(
              nameFilterValue,
              data.upstreamsList,
              checkboxes
            )}
            columns={TableColumns}
          />
        )}
      </div>
    );
  };

  const getUsableCatalogData = (
    nameFilter: string,
    data: Upstream.AsObject[]
  ) => {
    const dataUsed: UpstreamCardData[] = data.map(upstream => {
      return {
        ...upstream,
        cardTitle: upstream.metadata!.name,
        cardSubtitle: upstream.metadata!.namespace,
        onRemovecard: (id: string): void => {},
        onExpanded: () => {},
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
            value: getResourceStatus(upstream)
          }
        ]
      };
    });

    return dataUsed.filter(row => row.cardTitle.includes(nameFilter));
  };

  const getUsableTableData = (
    nameFilter: string,
    data: Upstream.AsObject[],
    checkboxes: CheckboxFilterProps[]
  ) => {
    const dataUsed = data.map(upstream => {
      return {
        ...upstream,
        status: getResourceStatus(upstream),
        type: getUpstreamType(upstream),
        name: upstream.metadata!.name,
        key: `${upstream.metadata!.name}-${upstream.metadata!.namespace}`
      };
    });
    let checkboxesNotSet = checkboxes.every(c => !c.value!);

    return dataUsed
      .filter(row => row.name.includes(nameFilter))
      .filter(row => {
        return (
          checkboxes.find(c => c.displayName === row.type)!.value! ||
          checkboxesNotSet
        );
      });
  };

  return (
    <div>
      <Heading>
        <Breadcrumb />
        <Action>
          <CreateUpstreamModal />
          <CatalogTableToggle
            listIsSelected={!catalogNotTable}
            onToggle={() => {
              setCatalogNotTable(cNt => !cNt);
            }}
          />
        </Action>
      </Heading>
      <ListingFilter
        strings={StringFilters}
        checkboxes={CheckboxFilters}
        filterFunction={listDisplay}
      />
    </div>
  );
};
