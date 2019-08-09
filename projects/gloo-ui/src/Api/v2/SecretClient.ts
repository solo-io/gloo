import { SecretApiClient } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/secret_pb_service';
import { host } from '../grpc-web-hooks';
import { grpc } from '@improbable-eng/grpc-web';
import {
  ListSecretsRequest,
  ListSecretsResponse,
  CreateSecretRequest
} from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/secret_pb';
import { SecretValuesType } from 'Components/Features/Settings/SecretForm';
import { ResourceRef } from 'proto/github.com/solo-io/solo-kit/api/v1/ref_pb';
import {
  Secret,
  AwsSecret,
  TlsSecret
} from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/secret_pb';

const client = new SecretApiClient(host, {
  transport: grpc.CrossBrowserHttpTransport({ withCredentials: false }),
  debug: true
});

function getSecretsList(params: {
  namespaces: string[];
}): Promise<ListSecretsResponse> {
  return new Promise((resolve, reject) => {
    let req = new ListSecretsRequest();
    req.setNamespacesList(params.namespaces);
    client.listSecrets(req, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        resolve(data!);
      }
    });
  });
}

function setSecretRequest(params: {
  secretKind: Secret.KindCase;
  values: SecretValuesType;
}) {
  const { secretKind, values } = params;
  let newSecret = new CreateSecretRequest();
  switch (secretKind) {
    case Secret.KindCase.AWS:
      const { accessKey, secretKey } = values.awsSecret;
      const awsSecret = new AwsSecret();
      awsSecret.setAccessKey(accessKey);
      awsSecret.setSecretKey(secretKey);
      newSecret.setAws(awsSecret);
      break;
    case Secret.KindCase.TLS:
      const tlsSecret = new TlsSecret();
      tlsSecret.setCertChain(values.tlsSecret.certChain);
      tlsSecret.setPrivateKey(values.tlsSecret.privateKey);
      tlsSecret.setRootCa(values.tlsSecret.rootCa);
      newSecret.setTls(tlsSecret);
    default:
      break;
  }
  return newSecret;
}
function getResourceRef(name: string, namespace: string): ResourceRef {
  let ref = new ResourceRef();
  ref.setName(name);
  ref.setNamespace(namespace);
  return ref;
}

function createSecret(params: {
  name: string;
  namespace: string;
  values: SecretValuesType;
  secretKind: Secret.KindCase;
}) {
  const { name, namespace, values, secretKind } = params;
  let newSecretReq = setSecretRequest({ secretKind, values });
  newSecretReq.setRef(getResourceRef(name, namespace));
  return new Promise((resolve, reject) => {
    client.createSecret(newSecretReq, (error, data) => {
      if (error !== null) {
        reject(error);
      } else {
        resolve(data!);
      }
    });
  });
}
export const secrets = {
  getSecretsList,
  createSecret
};
