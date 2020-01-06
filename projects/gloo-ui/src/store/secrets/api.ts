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
  DeleteSecretResponse,
  UpdateSecretRequest,
  UpdateSecretResponse
} from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/secret_pb';
import { SecretValuesType } from 'Components/Features/Settings/SecretForm';
import { ResourceRef } from 'proto/github.com/solo-io/solo-kit/api/v1/ref_pb';
import {
  Secret,
  AwsSecret,
  TlsSecret,
  AzureSecret
} from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/secret_pb';
import { guardByLicense } from 'store/config/actions';
import { Metadata } from 'proto/github.com/solo-io/solo-kit/api/v1/metadata_pb';
import {
  OauthSecret,
  ApiKeySecret
} from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/enterprise/options/extauth/v1/extauth_pb';

const client = new SecretApiClient(host, {
  transport: grpc.CrossBrowserHttpTransport({ withCredentials: false }),
  debug: true
});

function getSecretsList(): Promise<ListSecretsResponse.AsObject> {
  return new Promise((resolve, reject) => {
    let req = new ListSecretsRequest();
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
): Promise<Secret> {
  return new Promise((resolve, reject) => {
    const { name, namespace } = getSecretRequest.ref!;
    let request = new GetSecretRequest();
    let ref = new ResourceRef();
    ref.setName(name);
    ref.setNamespace(namespace);
    request.setRef();

    client.getSecret(request, (error, data) => {
      if (error !== null) {
        reject(error);
      } else {
        resolve(data?.getSecret());
      }
    });
  });
}

function setSecretValues(
  secret: Secret.AsObject,
  secretToUpdate = new Secret()
): Secret {
  let { aws, azure, tls, oauth, apiKey, extensions, metadata } = secret;
  if (metadata !== undefined) {
    let { name, namespace } = metadata;
    let newMetadata = new Metadata();
    newMetadata.setName(name);
    newMetadata.setNamespace(namespace);
    secretToUpdate.setMetadata(newMetadata);
  }

  if (aws !== undefined) {
    let { accessKey, secretKey } = aws;
    let awsSecret = new AwsSecret();
    awsSecret.setAccessKey(accessKey);
    awsSecret.setSecretKey(secretKey);
    secretToUpdate.setAws(awsSecret);
  }

  if (azure !== undefined) {
    let { apiKeysMap } = azure;

    let azureSecret = new AzureSecret();
    apiKeysMap.forEach(([key, val]) => {
      azureSecret.getApiKeysMap().set(key, val);
    });
    secretToUpdate.setAzure(azureSecret);
  }
  if (tls !== undefined) {
    let { certChain, rootCa, privateKey } = tls;
    let tlsSecret = new TlsSecret();
    tlsSecret.setCertChain(certChain);
    tlsSecret.setRootCa(rootCa);
    tlsSecret.setPrivateKey(privateKey);
    secretToUpdate.setTls(tlsSecret);
  }

  if (oauth !== undefined) {
    let { clientSecret } = oauth;
    let oauthSecret = new OauthSecret();
    oauthSecret.setClientSecret(clientSecret);
    secretToUpdate.setOauth(oauthSecret);
  }
  if (apiKey !== undefined) {
    let { apiKey: apiKeyValue, labelsList, generateApiKey } = apiKey;
    let apiKeySecret = new ApiKeySecret();
    apiKeySecret.setApiKey(apiKeyValue);
    apiKeySecret.setLabelsList(labelsList);
    apiKeySecret.setGenerateApiKey(generateApiKey);
    secretToUpdate.setApiKey(apiKeySecret);
  }
  // TODO
  if (extensions !== undefined) {
    let { configsMap } = extensions;
  }

  return secretToUpdate;
}

function createSecret(
  createSecretRequest: CreateSecretRequest.AsObject
): Promise<CreateSecretResponse.AsObject> {
  return new Promise((resolve, reject) => {
    let request = new CreateSecretRequest();
    let { secret } = createSecretRequest;
    if (secret !== undefined) {
      let inputSecret = setSecretValues(secret);
      request.setSecret(inputSecret);
    }

    guardByLicense();
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

function updateSecret(
  updateSecretRequest: UpdateSecretRequest.AsObject
): Promise<UpdateSecretResponse.AsObject> {
  return new Promise(async (resolve, reject) => {
    let request = new UpdateSecretRequest();
    let { secret } = updateSecretRequest;
    if (secret !== undefined && secret.metadata !== undefined) {
      let { name, namespace } = secret.metadata;
      let secretToUpdate = await getSecret({ ref: { name, namespace } });
      let updatedSecret = setSecretValues(secret, secretToUpdate);
      request.setSecret(updatedSecret);
    }
    client.updateSecret(request, (error, data) => {
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

function deleteSecret(
  deleteSecretRequest: DeleteSecretRequest.AsObject
): Promise<DeleteSecretResponse.AsObject> {
  const { name, namespace } = deleteSecretRequest.ref!;
  let deleteSecretReq = new DeleteSecretRequest();

  let ref = new ResourceRef();
  ref.setName(name);
  ref.setNamespace(namespace);

  deleteSecretReq.setRef(ref);
  return new Promise((resolve, reject) => {
    guardByLicense();
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
  updateSecret,
  deleteSecret
};
