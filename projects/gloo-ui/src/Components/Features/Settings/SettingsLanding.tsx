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
import { SecretsPage } from './SecretsPage';
import { WatchedNamespacesPage } from './WatchedNamespacesPage';
import { SecurityPage } from './SecurityPage';

const PageChoiceFilter: TypeFilterProps = {
  id: 'pageChoice',
  options: [
    {
      displayName: 'Security'
    },
    {
      displayName: 'Watched Namespaces'
    },
    {
      displayName: 'Secrets'
    }
  ],
  choice: 'Security'
};

interface Props extends RouteComponentProps {}

export const SettingsLanding = (props: Props) => {
  const locationEnding = props.location.pathname.split('/settings/')[1];
  const startingChoice =
    locationEnding.length === 0
      ? 'Security'
      : locationEnding === 'namespaces'
      ? 'Watched Namespaces'
      : locationEnding.charAt(0).toUpperCase() + locationEnding.slice(1);

  const pageChanged = (
    strings: StringFilterProps[],
    types: TypeFilterProps[],
    checkboxes: CheckboxFilterProps[],
    radios: RadioFilterProps[]
  ) => {
    const newChoice = types.find(type => type.id === 'pageChoice')!.choice!;
    const newPageLocation =
      newChoice === 'Watched Namespaces'
        ? 'namespaces'
        : newChoice.toLowerCase();

    props.history.push({
      pathname: `/settings/${newPageLocation}`
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
          <Route path='/settings/security/' render={() => <SecurityPage />} />
          <Route
            path='/settings/namespaces/'
            render={() => <WatchedNamespacesPage />}
          />
          <Route path='/settings/secrets/' render={() => <SecretsPage />} />

          <Redirect exact from='/settings/' to='/settings/security/' />
        </Switch>
      </React.Fragment>
    );
  };

  return (
    <ListingFilter
      types={[{ ...PageChoiceFilter, choice: startingChoice }]}
      filterFunction={listDisplay}
      onChange={pageChanged}
    />
  );
};
