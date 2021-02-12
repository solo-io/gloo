import { host, getObjectRefClassFromRefObj } from './helpers';
import { grpc } from '@improbable-eng/grpc-web';
import { FederatedEnterpriseGlooResourceApiClient } from '../proto/github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/federated_enterprise_gloo_resources_pb_service';
import {
  FederatedAuthConfig,
  ListFederatedAuthConfigsRequest,
  GetFederatedAuthConfigYamlRequest,
} from '../proto/github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/federated_enterprise_gloo_resources_pb';
import { FederatedRatelimitResourceApiClient } from '../proto/github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/federated_ratelimit_resources_pb_service';
import {
  FederatedRateLimitConfig,
  ListFederatedRateLimitConfigsRequest,
  GetFederatedRateLimitConfigYamlRequest,
} from '../proto/github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/federated_ratelimit_resources_pb';
import { ObjectRef } from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';
import { ResourceYaml } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/common_pb';

const federatedGlooResourceApiClient = new FederatedEnterpriseGlooResourceApiClient(
  host,
  {
    transport: grpc.CrossBrowserHttpTransport({ withCredentials: false }),
    debug: true,
  }
);
const federatedRateLimitResourceApiClient = new FederatedRatelimitResourceApiClient(
  host,
  {
    transport: grpc.CrossBrowserHttpTransport({ withCredentials: false }),
    debug: true,
  }
);

export const federatedEnterpriseGlooResourceApi = {
  listFederatedAuthConfigs,
  getFederatedAuthConfigYAML,
  listFederatedRateLimitConfigs,
  getFederatedRateLimitYAML,
};

function listFederatedAuthConfigs(): Promise<FederatedAuthConfig.AsObject[]> {
  let request = new ListFederatedAuthConfigsRequest();

  return new Promise((resolve, reject) => {
    federatedGlooResourceApiClient.listFederatedAuthConfigs(
      request,
      (error, data) => {
        if (error !== null) {
          console.error('Error:', error.message);
          console.error('Code:', error.code);
          console.error('Metadata:', error.metadata);
          reject(error);
        } else {
          resolve(data!.toObject().federatedAuthConfigsList);
        }
      }
    );
  });
}

function getFederatedAuthConfigYAML(
  fedAuthConfigObjectRef: ObjectRef.AsObject
): Promise<string> {
  let request = new GetFederatedAuthConfigYamlRequest();
  request.setFederatedAuthConfigRef(
    getObjectRefClassFromRefObj(fedAuthConfigObjectRef)
  );

  return new Promise((resolve, reject) => {
    federatedGlooResourceApiClient.getFederatedAuthConfigYaml(
      request,
      (error, data) => {
        if (error !== null) {
          console.error('Error:', error.message);
          console.error('Code:', error.code);
          console.error('Metadata:', error.metadata);
          reject(error);
        } else {
          resolve(data!.toObject().yamlData?.yaml ?? 'None');
        }
      }
    );
  });
}

function listFederatedRateLimitConfigs(): Promise<
  FederatedRateLimitConfig.AsObject[]
> {
  let request = new ListFederatedRateLimitConfigsRequest();

  return new Promise((resolve, reject) => {
    federatedRateLimitResourceApiClient.listFederatedRateLimitConfigs(
      request,
      (error, data) => {
        if (error !== null) {
          console.error('Error:', error.message);
          console.error('Code:', error.code);
          console.error('Metadata:', error.metadata);
          reject(error);
        } else {
          resolve(data!.toObject().federatedRateLimitConfigsList);
        }
      }
    );
  });
}

function getFederatedRateLimitYAML(
  fedRateLimitConfigObjectRef: ObjectRef.AsObject
): Promise<ResourceYaml.AsObject> {
  let request = new GetFederatedRateLimitConfigYamlRequest();
  request.setFederatedRateLimitConfigRef(
    getObjectRefClassFromRefObj(fedRateLimitConfigObjectRef)
  );

  return new Promise((resolve, reject) => {
    federatedRateLimitResourceApiClient.getFederatedRateLimitConfigYaml(
      request,
      (error, data) => {
        if (error !== null) {
          console.error('Error:', error.message);
          console.error('Code:', error.code);
          console.error('Metadata:', error.metadata);
          reject(error);
        } else {
          resolve(data!.toObject().yamlData);
        }
      }
    );
  });
}
