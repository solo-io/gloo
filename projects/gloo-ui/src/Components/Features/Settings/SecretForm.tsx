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
import { Formik, Field } from 'formik';
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
  // oAuthSecret: unknown TODO: handle OAuth secret
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

    const resourceRef = new ResourceRef();
    resourceRef.setName(values.secretResourceRef.name);
    resourceRef.setNamespace(values.secretResourceRef.namespace);
    secretReq.setRef(resourceRef);

    const awsSecret = new AwsSecret();
    awsSecret.setAccessKey(values.awsSecret.accessKey);
    awsSecret.setSecretKey(values.awsSecret.secretKey);
    secretReq.setAws(awsSecret);

    // TODO: figure out correct way to input api keys map
    // https://docs.microsoft.com/en-us/azure/search/search-security-api-keys
    const azureSecret = new AzureSecret();
    const apiKeys = new Map<string, string>();

    azureSecret.getApiKeysMap().set('keyname', 'key');

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
              <Field
                name='secretResourceRef.name'
                placeholder='Name'
                component={SoloFormInput}
              />
              <Field
                name='secretResourceRef.namespace'
                placeholder='Namespace'
                defaultValue='gloo-system'
                presetOptions={namespaces}
                component={SoloFormTypeahead}
              />
            </TableFormWrapper>
            {secretKind === Secret.KindCase.AWS && <AwsSecretFields />}
            {secretKind === Secret.KindCase.AZURE && <AzureSecretFields />}
            {secretKind === Secret.KindCase.TLS && <TlsSecretFields />}
            {/* {secretKind === Secret.KindCase.EXTENSION && <OAuthSecretFields />} */}

            <td>
              <GreenPlus
                style={{ cursor: 'pointer' }}
                onClick={() => handleSubmit()}
              />
            </td>
          </React.Fragment>
        )}
      </Formik>
    </React.Fragment>
  );
};

const AwsSecretFields: React.FC = () => {
  return (
    <TableFormWrapper>
      <Field
        name='awsSecret.accessKey'
        placeholder='Access Key'
        component={SoloFormInput}
      />
      <Field
        name='awsSecret.secretKey'
        placeholder='Secret Key'
        component={SoloFormInput}
      />
    </TableFormWrapper>
  );
};

const TlsSecretFields: React.FC = () => {
  return (
    <TableFormWrapper>
      <Field
        name='tlsSecret.certChain'
        placeholder='Cert Chain'
        component={SoloFormInput}
      />
      <Field
        name='tlsSecret.privateKey'
        placeholder='Private Key'
        component={SoloFormInput}
      />
      <Field name='rootCa' placeholder='Root Ca' component={SoloFormInput} />
    </TableFormWrapper>
  );
};

const AzureSecretFields: React.FC = () => {
  return (
    <TableFormWrapper>
      <Field
        name='tlsSecret.certChain'
        placeholder='Cert Chain'
        component={SoloFormInput}
      />
    </TableFormWrapper>
  );
};

const OAuthSecretFields: React.FC = () => {
  return <div>oauth</div>;
};
