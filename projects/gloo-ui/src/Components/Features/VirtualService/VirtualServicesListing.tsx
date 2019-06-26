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

export interface RouteParams {
  //... eg, virtualservice?: string
}

function VirtualServicesListingC({
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

  return <ListingFilter strings={StringFilters} filterFunction={listDisplay} />;
}

export const VirtualServicesListing = withRouter(VirtualServicesListingC);
