import styled from '@emotion/styled';
import { Breadcrumb } from 'Components/Common/Breadcrumb';
import { SuccessModal } from 'Components/Common/DisplayOnly/SuccessModal';
import {
  CheckboxFilterProps,
  ListingFilter,
  RadioFilterProps,
  StringFilterProps,
  TypeFilterProps
} from 'Components/Common/ListingFilter';
import {
  AwsSecret,
  Secret,
  TlsSecret
} from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/secret_pb';
import * as React from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { Redirect, Route, Switch, useHistory, useLocation } from 'react-router';
import { AppState } from 'store';
import { createSecret, deleteSecret, listSecrets } from 'store/secrets/actions';
import { SecretValuesType } from './SecretForm';
import { SecretsPage } from './SecretsPage';
import { SecurityPage } from './SecurityPage';
import { WatchedNamespacesPage } from './WatchedNamespacesPage';

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

export const SettingsLanding = () => {
  let location = useLocation();
  let history = useHistory();
  const [awsSecrets, setAwsSecrets] = React.useState<Secret.AsObject[]>([]);
  const [azureSecrets, setAzureSecrets] = React.useState<Secret.AsObject[]>([]);
  const [tlsSecrets, setTlsSecrets] = React.useState<Secret.AsObject[]>([]);
  const [oAuthSecrets, setOAuthSecrets] = React.useState<Secret.AsObject[]>([]);

  const [showSuccessModal, setShowSuccessModal] = React.useState(false);

  // Redux
  const dispatch = useDispatch();
  const {
    secrets: { secretsList }
  } = useSelector((state: AppState) => state);
  const [isLoading, setIsLoading] = React.useState(false);
  const [allSecrets, setAllSecrets] = React.useState<Secret.AsObject[]>([]);

  React.useEffect(() => {
    if (secretsList.length) {
      setIsLoading(false);
      setAllSecrets(secretsList);
    } else {
      dispatch(listSecrets());
      setIsLoading(true);
    }
    return () => {
      setShowSuccessModal(false);
    };
  }, [secretsList.length, showSuccessModal]);

  React.useEffect(() => {
    if (secretsList && allSecrets) {
      setAwsSecrets(allSecrets.filter(s => !!s.aws));
      setAzureSecrets(allSecrets.filter(s => !!s.azure));
      setOAuthSecrets(allSecrets.filter(s => !!s.extension));
      setTlsSecrets(allSecrets.filter(s => !!s.tls));
    }
  }, [allSecrets.length, showSuccessModal]);

  if (!secretsList || (!secretsList && isLoading)) {
    return <div>Loading...</div>;
  }

  // Get subpage without the / at the end
  const locationEnding = location.pathname.split('/settings/')[1].slice(0, -1);

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

    history.replace({
      pathname: `/settings/${newPageLocation}/`
    });
  };

  const listDisplay = (
    strings: StringFilterProps[],
    types: TypeFilterProps[],
    checkboxes: CheckboxFilterProps[],
    radios: RadioFilterProps[]
  ): React.ReactNode => {
    return (
      <>
        <Switch>
          <Route
            path='/settings/security/'
            render={() => (
              <SecurityPage
                tlsSecrets={tlsSecrets}
                oAuthSecrets={oAuthSecrets}
                onCreateSecret={handleCreateSecret}
                onDeleteSecret={handleDeleteSecret}
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
                onCreateSecret={handleCreateSecret}
                onDeleteSecret={handleDeleteSecret}
              />
            )}
          />

          <Redirect exact from='/settings/' to='/settings/security/' />
        </Switch>
      </>
    );
  };

  async function handleCreateSecret(
    values: SecretValuesType,
    secretKind: Secret.KindCase
  ) {
    const {
      secretResourceRef: { name, namespace }
    } = values;

    let aws: AwsSecret.AsObject | undefined = undefined;
    if (secretKind === Secret.KindCase.AWS) {
      aws = values.awsSecret;
    }
    let tls: TlsSecret.AsObject | undefined = undefined;
    if (secretKind === Secret.KindCase.TLS) {
      tls = values.tlsSecret;
    }

    dispatch(createSecret({ ref: { name, namespace }, aws, tls }));
  }

  async function handleDeleteSecret(
    name: string,
    namespace: string,
    secretKind: Secret.KindCase
  ) {
    try {
      dispatch(deleteSecret({ ref: { name, namespace } }));
    } catch (error) {
      console.error(error);
    }
  }

  return (
    <div>
      <Heading>
        <Breadcrumb />
      </Heading>
      <SuccessModal
        visible={showSuccessModal}
        successMessage='Secret added successfully'
      />
      <ListingFilter
        types={[{ ...PageChoiceFilter, choice: startingChoice }]}
        onChange={pageChanged}>
        {() => (
          <>
            <Switch>
              <Route
                path='/settings/security/'
                render={() => (
                  <SecurityPage
                    tlsSecrets={tlsSecrets}
                    oAuthSecrets={oAuthSecrets}
                    onCreateSecret={handleCreateSecret}
                    onDeleteSecret={handleDeleteSecret}
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
                    onCreateSecret={handleCreateSecret}
                    onDeleteSecret={handleDeleteSecret}
                  />
                )}
              />

              <Redirect exact from='/settings/' to='/settings/security/' />
            </Switch>
          </>
        )}
      </ListingFilter>
    </div>
  );
};
