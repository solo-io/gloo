import * as React from 'react';
/** @jsx jsx */
import { jsx } from '@emotion/core';

import styled from '@emotion/styled/macro';
import { RouteComponentProps, Route, Switch, Redirect } from 'react-router';
import { colors } from 'Styles';
import {
  ListingFilter,
  TypeFilterProps,
  StringFilterProps,
  CheckboxFilterProps,
  RadioFilterProps
} from 'Components/Common/ListingFilter';
import { Envoy } from './Envoy';
import { Proxys } from './Proxy';
import { Gateways } from './Gateways';
import { Breadcrumb } from 'Components/Common/Breadcrumb';

const PageChoiceFilter: TypeFilterProps = {
  id: 'pageChoice',
  options: [
    {
      displayName: 'Gateways'
    },
    {
      displayName: 'Proxy'
    },
    {
      displayName: 'Envoy'
    }
  ],
  choice: 'Gateways'
};

const Heading = styled.div`
  display: flex;
  justify-content: space-between;
  margin-bottom: 20px;
`;

interface Props extends RouteComponentProps<{ sublocation: string }> {}

export const AdminHub = (props: Props) => {
  const [showSuccessModal, setShowSuccessModal] = React.useState(false);

  const locationChoice =
    props.match.params.sublocation.charAt(0).toUpperCase() +
    props.match.params.sublocation.slice(1);

  const pageChanged = (
    strings: StringFilterProps[],
    types: TypeFilterProps[],
    checkboxes: CheckboxFilterProps[],
    radios: RadioFilterProps[]
  ) => {
    const newChoice = types.find(type => type.id === 'pageChoice')!.choice!;
    props.history.replace({
      pathname: `/admin/${newChoice.toLowerCase()}`
    });
  };

  const listDisplay = (
    strings: StringFilterProps[],
    types: TypeFilterProps[],
    checkboxes: CheckboxFilterProps[],
    radios: RadioFilterProps[]
  ): React.ReactNode => {
    return (
      <React.Fragment>
        <Switch>
          <Route path='/admin/gateways/' render={() => <Gateways />} />
          <Route path='/admin/proxy/' render={() => <Proxys />} />
          <Route path='/admin/envoy/' render={() => <Envoy />} />
        </Switch>
      </React.Fragment>
    );
  };

  return (
    <div>
      <Heading>
        <Breadcrumb />
      </Heading>
      <ListingFilter
        types={[{ ...PageChoiceFilter, choice: locationChoice }]}
        filterFunction={listDisplay}
        onChange={pageChanged}
      />
    </div>
  );
};
