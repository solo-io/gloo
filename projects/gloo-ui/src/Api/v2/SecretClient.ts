import { SecretApiClient } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/secret_pb_service';
import { host } from 'store';
import { grpc } from '@improbable-eng/grpc-web';
import {
  ListSecretsRequest,
  ListSecretsResponse,
  CreateSecretRequest,
  DeleteSecretRequest,
  GetSecretRequest,
  CreateSecretResponse,
  GetSecretResponse,
  DeleteSecretResponse
} from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/secret_pb';
import { SecretValuesType } from 'Components/Features/Settings/SecretForm';
import { ResourceRef } from 'proto/github.com/solo-io/solo-kit/api/v1/ref_pb';
import {
  Secret,
  AwsSecret,
  TlsSecret,
  AzureSecret
} from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/secret_pb';
import { getResourceRef } from './helpers';
import { guardByLicense } from 'store/config/actions';

const client = new SecretApiClient(host, {
  transport: grpc.CrossBrowserHttpTransport({ withCredentials: false }),
  debug: true
});

function getSecretsList(
  listSecretsRequest: ListSecretsRequest.AsObject
): Promise<ListSecretsResponse.AsObject> {
  return new Promise((resolve, reject) => {
    let req = new ListSecretsRequest();
    req.setNamespacesList(listSecretsRequest.namespacesList);
    client.listSecrets(req, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        resolve(data!.toObject());
      }
    });
  });
}

function getSecret(
  getSecretRequest: GetSecretRequest.AsObject
): Promise<GetSecretResponse.AsObject> {
  const { name, namespace } = getSecretRequest.ref!;
  let request = new GetSecretRequest();
  let ref = new ResourceRef();
  ref.setName(name);
  ref.setNamespace(namespace);
  request.setRef();
  return new Promise((resolve, reject) => {
    client.getSecret(request, (error, data) => {
      if (error !== null) {
        reject(error);
      } else {
        resolve(data!.toObject());
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
      throw new Error('Not supported');
  }
  return newSecret;
}

// TODO: Support other secrets
export function getCreateSecret(
  createSecretRequest: CreateSecretRequest.AsObject
): Promise<CreateSecretResponse.AsObject> {
  return new Promise((resolve, reject) => {
    let request = new CreateSecretRequest();
    let ref = new ResourceRef();
    ref.setName(createSecretRequest.ref!.name);
    ref.setNamespace(createSecretRequest.ref!.namespace);
    request.setRef(ref);
    let awsSecret = new AwsSecret();
    let tlsSecret = new TlsSecret();
    let azureSecret = new AzureSecret();

    if (createSecretRequest.aws) {
      const { accessKey, secretKey } = createSecretRequest.aws;
      awsSecret.setAccessKey(accessKey);
      awsSecret.setSecretKey(secretKey);
      request.setAws(awsSecret);
    }
    if (createSecretRequest.tls) {
      const { certChain, privateKey, rootCa } = createSecretRequest.tls;
      tlsSecret.setCertChain(certChain);
      tlsSecret.setPrivateKey(privateKey);
      tlsSecret.setRootCa(rootCa);
      request.setTls(tlsSecret);
    }

    guardByLicense()
    client.createSecret(request, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        resolve(data!.toObject());
      }
    });
  });
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
    guardByLicense()
    client.createSecret(newSecretReq, (error, data) => {
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

function deleteSecret(
  deleteSecretRequest: DeleteSecretRequest.AsObject
): Promise<DeleteSecretResponse.AsObject> {
  const { name, namespace } = deleteSecretRequest.ref!;
  let deleteSecretReq = new DeleteSecretRequest();
  let ref = getResourceRef(name, namespace);

  deleteSecretReq.setRef(ref);
  return new Promise((resolve, reject) => {
    guardByLicense()
    client.deleteSecret(deleteSecretReq, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        resolve(data!.toObject());
      }
    });
  });
}
export const secrets = {
  getSecretsList,
  getSecret,
  createSecret,
  deleteSecret
};
