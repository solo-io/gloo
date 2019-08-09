import React from 'react';
import { ResourceRef } from 'proto/github.com/solo-io/solo-kit/api/v1/ref_pb';
import {
  AwsSecret,
  AzureSecret,
  TlsSecret,
  Secret
} from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/secret_pb';
import { Formik } from 'formik';
import {
  SoloFormInput,
  TableFormWrapper,
  SoloFormTypeahead
} from 'Components/Common/Form/SoloFormField';
import { ReactComponent as GreenPlus } from 'assets/small-green-plus.svg';
import { NamespacesContext } from 'GlooIApp';

// TODO: modify for use outside a table
// TODO: set one source of truth for column names/order

export interface SecretValuesType {
  secretResourceRef: ResourceRef.AsObject;
  awsSecret: AwsSecret.AsObject;
  azureSecret: AzureSecret.AsObject;
  tlsSecret: TlsSecret.AsObject;
  oAuthSecret: { clientSecret: string };
}

interface Props {
  secretKind: Secret.KindCase;
  onCreateSecret: (
    values: SecretValuesType,
    secretKind: Secret.KindCase
  ) => void;
}

export const SecretForm: React.FC<Props> = props => {
  const { secretKind, onCreateSecret } = props;
  const namespaces = React.useContext(NamespacesContext);

  const initialValues: SecretValuesType = {
    secretResourceRef: {
      name: '',
      namespace: namespaces.defaultNamespace
    },
    awsSecret: {
      accessKey: '',
      secretKey: ''
    },
    azureSecret: {
      apiKeysMap: []
    },
    tlsSecret: {
      certChain: '',
      privateKey: '',
      rootCa: ''
    },
    oAuthSecret: {
      clientSecret: ''
    }
  };

  const createSecret = (values: SecretValuesType) => {
    onCreateSecret(values, secretKind);
  };

  return (
    <React.Fragment>
      <Formik<SecretValuesType>
        initialValues={initialValues}
        onSubmit={createSecret}>
        {({ handleSubmit }) => (
          <React.Fragment>
            <TableFormWrapper>
              <SoloFormInput name='secretResourceRef.name' placeholder='Name' />
              <SoloFormTypeahead
                name='secretResourceRef.namespace'
                placeholder='Namespace'
                defaultValue={namespaces.defaultNamespace}
                presetOptions={namespaces.namespacesList.map(ns => {
                  return { value: ns };
                })}
              />
            </TableFormWrapper>
            {secretKind === Secret.KindCase.AWS && <AwsSecretFields />}
            {secretKind === Secret.KindCase.AZURE && <AzureSecretFields />}
            {secretKind === Secret.KindCase.TLS && <TlsSecretFields />}
            {secretKind === Secret.KindCase.EXTENSION && <OAuthSecretFields />}
            <TableFormWrapper>
              <GreenPlus
                style={{ cursor: 'pointer' }}
                onClick={() => handleSubmit()}
              />
            </TableFormWrapper>
          </React.Fragment>
        )}
      </Formik>
    </React.Fragment>
  );
};

const AwsSecretFields: React.FC = () => {
  return (
    <TableFormWrapper>
      <SoloFormInput name='awsSecret.accessKey' placeholder='Access Key' />
      <SoloFormInput name='awsSecret.secretKey' placeholder='Secret Key' />
    </TableFormWrapper>
  );
};

const TlsSecretFields: React.FC = () => {
  return (
    <TableFormWrapper>
      <SoloFormInput name='tlsSecret.certChain' placeholder='Cert Chain' />
      <SoloFormInput name='tlsSecret.privateKey' placeholder='Private Key' />
      <SoloFormInput name='tlsSecret.rootCa' placeholder='Root Ca' />
    </TableFormWrapper>
  );
};

const AzureSecretFields: React.FC = () => {
  return (
    <TableFormWrapper>
      <SoloFormInput name='azureSecret.apiKeysMap' placeholder='Api Keys' />
    </TableFormWrapper>
  );
};

const OAuthSecretFields: React.FC = () => {
  return (
    <TableFormWrapper>
      <SoloFormInput
        name='oAuthSecret.clientSecret'
        placeholder='Client Secret'
      />
    </TableFormWrapper>
  );
};
