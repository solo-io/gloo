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
import { OauthSecret } from 'proto/gloo/projects/gloo/api/v1/enterprise/options/extauth/v1/extauth_pb';
import {
  AwsSecret,
  AzureSecret,
  Secret,
  TlsSecret
} from 'proto/gloo/projects/gloo/api/v1/secret_pb';
import * as React from 'react';
import { useDispatch } from 'react-redux';
import { Redirect, Route, Switch, useHistory, useLocation } from 'react-router';
import { createSecret, deleteSecret } from 'store/secrets/actions';
import { secretAPI } from 'store/secrets/api';
import useSWR from 'swr';
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

  const { data: secretsList, error } = useSWR(
    'listSecrets',
    secretAPI.getSecretsList
  );

  React.useEffect(() => {
    if (secretsList && secretsList.length) {
      setAwsSecrets(secretsList.filter(s => !!s.aws));
      setAzureSecrets(secretsList.filter(s => !!s.azure));
      setOAuthSecrets(secretsList.filter(s => !!s.oauth));
      setTlsSecrets(secretsList.filter(s => !!s.tls));
    }
  }, [secretsList?.length, showSuccessModal]);

  if (!secretsList) {
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
    let newSecret = new Secret().toObject();

    const {
      secretResourceRef: { name, namespace }
    } = values;

    let aws: AwsSecret.AsObject | undefined = undefined;
    if (secretKind === Secret.KindCase.AWS) {
      aws = values.awsSecret;
      newSecret = { ...newSecret, aws };
    }

    let azure: AzureSecret.AsObject | undefined = undefined;
    if (secretKind === Secret.KindCase.AZURE) {
      azure = values.azureSecret;
      newSecret = { ...newSecret, azure };
    }

    let tls: TlsSecret.AsObject | undefined = undefined;
    if (secretKind === Secret.KindCase.TLS) {
      tls = values.tlsSecret;
      newSecret = { ...newSecret, tls };
    }

    let oauth: OauthSecret.AsObject | undefined = undefined;
    if (secretKind === Secret.KindCase.OAUTH) {
      oauth = values.oAuthSecret;
      newSecret = { ...newSecret, oauth };
    }

    dispatch(
      createSecret({
        secret: {
          ...newSecret,
          metadata: {
            ...newSecret.metadata!,
            name,
            namespace
          }
        }
      })
    );
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
                data-testid='secrets'
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
