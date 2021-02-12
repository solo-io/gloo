import { FederatedGlooResourceApiClient } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/federated_gloo_resources_pb_service';
import { host, getObjectRefClassFromRefObj } from './helpers';
import { grpc } from '@improbable-eng/grpc-web';
import {
  ListFederatedSettingsRequest,
  ListFederatedUpstreamGroupsRequest,
  ListFederatedUpstreamsRequest,
  FederatedUpstream,
  FederatedSettings,
  FederatedUpstreamGroup,
  GetFederatedUpstreamYamlRequest,
  GetFederatedSettingsYamlRequest,
} from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/federated_gloo_resources_pb';
import { ObjectRef } from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';
const federatedGlooResourceApiClient = new FederatedGlooResourceApiClient(
  host,
  {
    transport: grpc.CrossBrowserHttpTransport({ withCredentials: false }),
    debug: true,
  }
);

export const federatedGlooResourceApi = {
  listFederatedUpstreams,
  listFederatedUpstreamGroups,
  listFederatedSettings,
  getFederatedUpstreamYAML,
  getFederatedUpstreamGroupYAML,
  getFederatedSettingYAML,
};

function listFederatedUpstreams(): Promise<FederatedUpstream.AsObject[]> {
  let request = new ListFederatedUpstreamsRequest();

  return new Promise((resolve, reject) => {
    federatedGlooResourceApiClient.listFederatedUpstreams(
      request,
      (error, data) => {
        if (error !== null) {
          console.error('Error:', error.message);
          console.error('Code:', error.code);
          console.error('Metadata:', error.metadata);
          reject(error);
        } else {
          resolve(data!.toObject().federatedUpstreamsList);
        }
      }
    );
  });
}

function listFederatedUpstreamGroups(): Promise<
  FederatedUpstreamGroup.AsObject[]
> {
  let request = new ListFederatedUpstreamGroupsRequest();

  return new Promise((resolve, reject) => {
    federatedGlooResourceApiClient.listFederatedUpstreamGroups(
      request,
      (error, data) => {
        if (error !== null) {
          console.error('Error:', error.message);
          console.error('Code:', error.code);
          console.error('Metadata:', error.metadata);
          reject(error);
        } else {
          resolve(data!.toObject().federatedUpstreamGroupsList);
        }
      }
    );
  });
}

function listFederatedSettings(): Promise<FederatedSettings.AsObject[]> {
  let request = new ListFederatedSettingsRequest();

  return new Promise((resolve, reject) => {
    federatedGlooResourceApiClient.listFederatedSettings(
      request,
      (error, data) => {
        if (error !== null) {
          console.error('Error:', error.message);
          console.error('Code:', error.code);
          console.error('Metadata:', error.metadata);
          reject(error);
        } else {
          resolve(data!.toObject().federatedSettingsList);
        }
      }
    );
  });
}

function getFederatedUpstreamYAML(
  fedUpstreamObjectRef: ObjectRef.AsObject
): Promise<string> {
  let request = new GetFederatedUpstreamYamlRequest();
  request.setFederatedUpstreamRef(
    getObjectRefClassFromRefObj(fedUpstreamObjectRef)
  );

  return new Promise((resolve, reject) => {
    federatedGlooResourceApiClient.getFederatedUpstreamYaml(
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

function getFederatedUpstreamGroupYAML(
  fedUpstreamObjectRef: ObjectRef.AsObject
): Promise<string> {
  let request = new GetFederatedUpstreamYamlRequest();
  request.setFederatedUpstreamRef(
    getObjectRefClassFromRefObj(fedUpstreamObjectRef)
  );

  return new Promise((resolve, reject) => {
    federatedGlooResourceApiClient.getFederatedUpstreamYaml(
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

function getFederatedSettingYAML(
  fedSettingObjectRef: ObjectRef.AsObject
): Promise<string> {
  let request = new GetFederatedSettingsYamlRequest();
  request.setFederatedSettingsRef(
    getObjectRefClassFromRefObj(fedSettingObjectRef)
  );

  return new Promise((resolve, reject) => {
    federatedGlooResourceApiClient.getFederatedSettingsYaml(
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
