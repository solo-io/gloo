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

const StringFilters: StringFilterProps[] = [
  {
    displayName: 'Filter By Name...',
    placeholder: 'Filter by name...',
    value: ''
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

export interface RouteParams {
  //... eg, virtualservice?: string
}

function UpstreamsListingC({
  history,
  match,
  location
}: RouteComponentProps<RouteParams>) {
  const [catalogNotTable, setCatalogNotTable] = React.useState<boolean>(true);

  const listDisplay = (
    strings: StringFilterProps[],
    types: TypeFilterProps[],
    checkboxes: CheckboxFilterProps[],
    radios: RadioFilterProps[]
  ) => {
    return <div />;
  };

  return (
    <ListingFilter
      strings={StringFilters}
      checkboxes={CheckboxFilters}
      filterFunction={listDisplay}
    />
  );
}

export const UpstreamsListing = withRouter(UpstreamsListingC);
