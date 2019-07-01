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
    dataIndex: 'domain'
  },
  {
    title: 'Namespace',
    dataIndex: 'namespace'
  },
  {
    title: 'Version',
    dataIndex: 'version'
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

interface Props extends RouteComponentProps {
  //... eg, virtualservice?: string
}

export const VirtualServicesListing = (props: Props) => {
  const [catalogNotTable, setCatalogNotTable] = React.useState<boolean>(true);
  const { history, match } = props;
  console.log(history, match);

  const getUsableCatalogData = (nameFilter: string): CardType[] => {
    //REPLACE
    const dataUsed: CardType[] = [1, 2, 3, 4, 5, 6].map(num => {
      return {
        ...data[0],
        cardTitle: data[0].name + num,
        cardSubtitle: 'subtitle' + num,
        onRemovecard: (id: string): void => {},
        onExpanded: () => {},
        onClick: () => {
          history.push(`${match.path}${data[0].name + num}/details`);
        },
        details: [
          {
            title: 'sample',
            value: 'maybe',
            valueDisplay: <div>No, we'll display this instead</div>
          }
        ]
      };
    });

    return dataUsed.filter(row => row.cardTitle.includes(nameFilter));
  };

  const getUsableTableData = (nameFilter: string): DataSourceType[] => {
    //REPLACE
    const dataUsed: DataSourceType[] = [1, 2, 3, 4].map(num => {
      return {
        ...data[0],
        name: data[0].name + num,
        key: `${num}`
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

    return catalogNotTable ? (
      <SectionCard cardName={'Virtual Services'} logoIcon={<Gloo />}>
        <CardsListing cardsData={getUsableCatalogData(nameFilterValue)} />
      </SectionCard>
    ) : (
      <SoloTable
        dataSource={getUsableTableData(nameFilterValue)}
        columns={TableColumns}
      />
    );
  };

  return (
    <div>
      <Heading>
        <Breadcrumb />
        <CatalogTableToggle
          listIsSelected={!catalogNotTable}
          onToggle={() => {
            setCatalogNotTable(cNt => !cNt);
          }}
        />
      </Heading>
      <ListingFilter strings={StringFilters} filterFunction={listDisplay} />
    </div>
  );
};
