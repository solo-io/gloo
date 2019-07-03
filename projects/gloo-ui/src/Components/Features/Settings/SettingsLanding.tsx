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
import { Breadcrumb } from 'Components/Common/Breadcrumb';
import { Secret } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/secret_pb';
import { useListSecrets } from 'Api';
import { ListSecretsRequest } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/secret_pb';
import { NamespacesContext } from 'GlooIApp';

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

const Heading = styled.div`
  display: flex;
  justify-content: space-between;
  margin-bottom: 20px;
`;

interface Props extends RouteComponentProps {}

export const SettingsLanding = (props: Props) => {
  let awsSecrets = [] as Secret.AsObject[];
  let azureSecrets = [] as Secret.AsObject[];
  let tlsSecrets = [] as Secret.AsObject[];
  let oAuthSecrets = [] as Secret.AsObject[];
  const listSecretsReq = React.useRef(new ListSecretsRequest());
  const namespaces = React.useContext(NamespacesContext);

  listSecretsReq.current.setNamespacesList(namespaces);
  const { data, loading, error } = useListSecrets(listSecretsReq.current);

  if (!data || loading) {
    return <div>Loading...</div>;
  }
  if (data && data.secretsList.length > 0) {
    data.secretsList.map(secret => {
      if (!!secret.aws) {
        awsSecrets.push(secret);
      }
      if (!!secret.azure) {
        azureSecrets.push(secret);
      }
      if (!!secret.tls) {
        tlsSecrets.push(secret);
      }
      if (!!secret.extension) {
        oAuthSecrets.push(secret);
      }
    });
  }

  const locationEnding = props.location.pathname.split('/settings/')[1];
  const startingChoice =
    locationEnding && locationEnding.length
      ? locationEnding === 'namespaces'
        ? 'Watched Namespaces'
        : locationEnding.charAt(0).toUpperCase() + locationEnding.slice(1)
      : 'Security';

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
          <Route
            path='/settings/security/'
            render={() => (
              <SecurityPage
                tlsSecrets={tlsSecrets}
                oAuthSecrets={oAuthSecrets}
              />
            )}
          />
          <Route
            path='/settings/namespaces/'
            render={() => <WatchedNamespacesPage />}
          />
          <Route
            path='/settings/secrets/'
            render={() => (
              <SecretsPage
                awsSecrets={awsSecrets}
                azureSecrets={azureSecrets}
              />
            )}
          />

          <Redirect exact from='/settings/' to='/settings/security/' />
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
        types={[{ ...PageChoiceFilter, choice: startingChoice }]}
        filterFunction={listDisplay}
        onChange={pageChanged}
      />
    </div>
  );
};
