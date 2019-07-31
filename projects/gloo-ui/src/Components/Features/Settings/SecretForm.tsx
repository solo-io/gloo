import React from 'react';
import { useCreateSecret } from 'Api';
import { CreateSecretRequest } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/secret_pb';
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

interface SecretValuesType {
  secretResourceRef: ResourceRef.AsObject;
  awsSecret: AwsSecret.AsObject;
  azureSecret: AzureSecret.AsObject;
  tlsSecret: TlsSecret.AsObject;
  oAuthSecret: { clientSecret: string };
}

const initialValues: SecretValuesType = {
  secretResourceRef: {
    name: '',
    namespace: ''
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

interface Props {
  secretKind: Secret.KindCase;
}

export const SecretForm: React.FC<Props> = ({ secretKind }) => {
  const { refetch: makeRequest } = useCreateSecret(null);
  const namespaces = React.useContext(NamespacesContext);

  const createSecret = (values: typeof initialValues) => {
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

    makeRequest(secretReq);
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
                defaultValue='gloo-system'
                presetOptions={namespaces}
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
      <SoloFormInput name='rootCa' placeholder='Root Ca' />
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
