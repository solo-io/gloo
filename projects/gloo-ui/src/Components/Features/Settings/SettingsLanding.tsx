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
import {
  Secret,
  AwsSecret,
  AzureSecret,
  TlsSecret
} from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/secret_pb';
import { useListSecrets, useCreateSecret, useDeleteSecret } from 'Api';
import {
  ListSecretsRequest,
  CreateSecretRequest,
  DeleteSecretRequest
} from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/secret_pb';
import { NamespacesContext } from 'GlooIApp';
import { SuccessModal } from 'Components/Common/SuccessModal';
import { SecretValuesType } from './SecretForm';
import { ResourceRef } from 'proto/github.com/solo-io/solo-kit/api/v1/ref_pb';

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
  const [showSuccessModal, setShowSuccessModal] = React.useState(false);
  const { refetch: makeCreateRequest } = useCreateSecret(null);
  const { refetch: makeDeleteRequest } = useDeleteSecret(null);
  listSecretsReq.current.setNamespacesList(namespaces.namespacesList);
  const { data, loading, error, refetch: updateSecretsList } = useListSecrets(
    listSecretsReq.current
  );
  const [allSecrets, setAllSecrets] = React.useState<Secret.AsObject[]>([]);

  React.useEffect(() => {
    if (!!data) {
      setAllSecrets(data.secretsList);
    }
    return () => {
      setShowSuccessModal(false);
    };
  }, [data, showSuccessModal]);

  React.useEffect(() => {
    if (data && data.secretsList) {
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
  }, [allSecrets.length]);

  if (!data || (!data && loading)) {
    return <div>Loading...</div>;
  }
  if (data && data.secretsList.length > 0) {
    data.secretsList.forEach(secret => {
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
                onCreateSecret={createSecret}
                onDeleteSecret={deleteSecret}
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
                onCreateSecret={createSecret}
                onDeleteSecret={deleteSecret}
              />
            )}
          />

          <Redirect exact from='/settings/' to='/settings/security/' />
        </Switch>
      </React.Fragment>
    );
  };

  function createSecret(values: SecretValuesType, secretKind: Secret.KindCase) {
    const secretReq = new CreateSecretRequest();
    switch (secretKind) {
      case Secret.KindCase.AWS:
        const awsSecret = new AwsSecret();
        awsSecret.setAccessKey(values.awsSecret.accessKey);
        awsSecret.setSecretKey(values.awsSecret.secretKey);
        secretReq.setAws(awsSecret);
        break;
      case Secret.KindCase.AZURE:
        // TODO: figure out correct way to input api keys map
        // https://docs.microsoft.com/en-us/azure/search/search-security-api-keys
        const azureSecret = new AzureSecret();
        azureSecret.getApiKeysMap().set('keyname', 'key');
        const apiKeys = new Map<string, string>();
        break;
      case Secret.KindCase.TLS:
        const tlsSecret = new TlsSecret();
        tlsSecret.setCertChain(values.tlsSecret.certChain);
        tlsSecret.setPrivateKey(values.tlsSecret.privateKey);
        tlsSecret.setRootCa(values.tlsSecret.rootCa);
        secretReq.setTls(tlsSecret);
        break;
      default:
        break;
    }
    const resourceRef = new ResourceRef();
    resourceRef.setName(values.secretResourceRef.name);
    resourceRef.setNamespace(values.secretResourceRef.namespace);
    secretReq.setRef(resourceRef);

    makeCreateRequest(secretReq);
    setTimeout(() => {
      updateSecretsList(listSecretsReq.current);
    }, 500);
    setShowSuccessModal(true);
  }

  function deleteSecret(
    name: string,
    namespace: string,
    secretKind: Secret.KindCase
  ) {
    if (secretKind === Secret.KindCase.AWS) {
      awsSecrets.filter(s => s.metadata!.name !== name);
    }
    if (secretKind === Secret.KindCase.AZURE) {
      azureSecrets.filter(s => s.metadata!.name !== name);
    }
    if (secretKind === Secret.KindCase.TLS) {
      tlsSecrets.filter(s => s.metadata!.name !== name);
    }
    if (secretKind === Secret.KindCase.EXTENSION) {
      oAuthSecrets.filter(s => s.metadata!.name !== name);
    }
    let req = new DeleteSecretRequest();
    let ref = new ResourceRef();
    ref.setName(name);
    ref.setNamespace(namespace);
    req.setRef(ref);
    makeDeleteRequest(req);

    setTimeout(() => {
      updateSecretsList(listSecretsReq.current);
    }, 500);
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
