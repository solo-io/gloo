import styled from '@emotion/styled';
import { Breadcrumb } from 'Components/Common/Breadcrumb';
import {
  CheckboxFilterProps,
  ListingFilter,
  RadioFilterProps,
  StringFilterProps,
  TypeFilterProps
} from 'Components/Common/ListingFilter';
import * as React from 'react';
import { Route, Switch, useHistory, useParams } from 'react-router';
import { Envoy } from './Envoy';
import { Gateways } from './Gateways';
import { Proxys } from './Proxy';
import { ErrorBoundary } from '../Errors/ErrorBoundary';
import { SecretsPage } from '../Settings/SecretsPage';
import { WatchedNamespacesPage } from '../Settings/WatchedNamespacesPage';
import { Settings } from '../Settings/SettingsDetails';

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
    },
    {
      displayName: 'Settings'
    },
    {
      displayName: 'Watched Namespaces'
    },
    {
      displayName: 'Secrets'
    }
  ],
  choice: 'Gateways'
};

const Heading = styled.div`
  display: flex;
  justify-content: space-between;
  margin-bottom: 20px;
`;

export const AdminHub = () => {
  let history = useHistory();
  let { sublocation } = useParams();
  const [showSuccessModal, setShowSuccessModal] = React.useState(false);

  const locationChoice =
    sublocation!.charAt(0).toUpperCase() + sublocation!.slice(1);

  const pageChanged = (
    strings: StringFilterProps[],
    types: TypeFilterProps[],
    checkboxes: CheckboxFilterProps[],
    radios: RadioFilterProps[]
  ) => {
    let newChoice = types.find(type => type.id === 'pageChoice')!.choice!;
    if (newChoice === 'Watched Namespaces') {
      newChoice = 'watched-namespaces';
    }
    history.replace({
      pathname: `/admin/${newChoice.toLowerCase()}`
    });
  };

  return (
    <ErrorBoundary
      fallback={<div>There was an error with the Admin section</div>}>
      <div>
        <Heading>
          <Breadcrumb />
        </Heading>
        <ListingFilter
          types={[{ ...PageChoiceFilter, choice: locationChoice }]}
          onChange={pageChanged}>
          {() => (
            <>
              <Switch>
                <Route path='/admin/gateways/' component={Gateways} />
                <Route path='/admin/proxy/' component={Proxys} />
                <Route path='/admin/envoy/' component={Envoy} />
                <Route path='/admin/settings' component={Settings} />
                <Route
                  path='/admin/watched-namespaces/'
                  component={WatchedNamespacesPage}
                />
                <Route path='/admin/secrets/' component={SecretsPage} />
              </Switch>
            </>
          )}
        </ListingFilter>
      </div>
    </ErrorBoundary>
  );
};
