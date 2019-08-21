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

import { SuccessModal } from 'Components/Common/DisplayOnly/SuccessModal';
import { SecretValuesType } from './SecretForm';

import { useDispatch, useSelector } from 'react-redux';
import { AppState } from 'store';
import { listSecrets, createSecret, deleteSecret } from 'store/secrets/actions';

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
  const [awsSecrets, setAwsSecrets] = React.useState<Secret.AsObject[]>([]);
  const [azureSecrets, setAzureSecrets] = React.useState<Secret.AsObject[]>([]);
  const [tlsSecrets, setTlsSecrets] = React.useState<Secret.AsObject[]>([]);
  const [oAuthSecrets, setOAuthSecrets] = React.useState<Secret.AsObject[]>([]);

  const [showSuccessModal, setShowSuccessModal] = React.useState(false);

  // Redux
  const dispatch = useDispatch();
  const {
    secrets: { secretsList },
    config: { namespacesList }
  } = useSelector((state: AppState) => state);
  const [isLoading, setIsLoading] = React.useState(false);
  const [allSecrets, setAllSecrets] = React.useState<Secret.AsObject[]>([]);

  React.useEffect(() => {
    if (secretsList.length) {
      setIsLoading(false);
      setAllSecrets(secretsList);
    } else {
      dispatch(listSecrets({ namespacesList }));
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
  const locationEnding = props.location.pathname
    .split('/settings/')[1]
    .slice(0, -1);

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

    props.history.replace({
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
      <React.Fragment>
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
      </React.Fragment>
    );
  };

  async function handleCreateSecret(
    values: SecretValuesType,
    secretKind: Secret.KindCase
  ) {
    const {
      secretResourceRef: { name, namespace }
    } = values;
    const { secretKey, accessKey } = values.awsSecret;
    dispatch(
      createSecret({ ref: { name, namespace }, aws: { accessKey, secretKey } })
    );
    // try {
    //   await secrets.createSecret({ name, namespace, values, secretKind });
    // } catch (error) {
    //   //   // TODO: show error modal
    //   //   console.error('error', error);
    //   // }
    //   // setNewVariables({ namespaces: namespaces.namespacesList });
    //   // setShowSuccessModal(true);
    // }
  }

  async function handleDeleteSecret(
    name: string,
    namespace: string,
    secretKind: Secret.KindCase
  ) {
    if (secretKind === Secret.KindCase.AWS) {
      setAwsSecrets(awsSecrets =>
        awsSecrets.filter(s => s.metadata!.name !== name)
      );
    }
    if (secretKind === Secret.KindCase.AZURE) {
      setAzureSecrets(azureSecrets =>
        azureSecrets.filter(s => s.metadata!.name !== name)
      );
    }
    if (secretKind === Secret.KindCase.TLS) {
      setTlsSecrets(tlsSecrets =>
        tlsSecrets.filter(s => s.metadata!.name !== name)
      );
    }
    if (secretKind === Secret.KindCase.EXTENSION) {
      setOAuthSecrets(oAuthSecrets =>
        oAuthSecrets.filter(s => s.metadata!.name !== name)
      );
    }
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
        filterFunction={listDisplay}
        onChange={pageChanged}
      />
    </div>
  );
};
