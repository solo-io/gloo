import {
  GetVersionRequest,
  GetVersionResponse,
  GetOAuthEndpointRequest,
  GetOAuthEndpointResponse,
  GetSettingsRequest,
  GetSettingsResponse,
  GetPodNamespaceRequest,
  GetPodNamespaceResponse,
  ListNamespacesRequest,
  ListNamespacesResponse,
  GetIsLicenseValidRequest,
  GetIsLicenseValidResponse,
  UpdateSettingsRequest,
  UpdateSettingsResponse
} from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/config_pb';
import { ConfigApiClient } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/config_pb_service';
import { host } from 'Api';
import { grpc } from '@improbable-eng/grpc-web';
import { ResourceRef } from 'proto/github.com/solo-io/solo-kit/api/v1/ref_pb';
import { Duration } from 'google-protobuf/google/protobuf/duration_pb';

const client = new ConfigApiClient(host, {
  transport: grpc.CrossBrowserHttpTransport({ withCredentials: false }),
  debug: true
});

function getVersion(): Promise<GetVersionResponse.AsObject> {
  return new Promise((resolve, reject) => {
    let request = new GetVersionRequest();
    client.getVersion(request, (error, data) => {
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

function getOAuthEndpoint(): Promise<GetOAuthEndpointResponse.AsObject> {
  return new Promise((resolve, reject) => {
    let request = new GetOAuthEndpointRequest();
    client.getOAuthEndpoint(request, (error, data) => {
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
function getSettings(): Promise<GetSettingsResponse.AsObject> {
  return new Promise((resolve, reject) => {
    let request = new GetOAuthEndpointRequest();
    client.getSettings(request, (error, data) => {
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

function updateSettings(
  updateSettingsRequest: UpdateSettingsRequest.AsObject
): Promise<UpdateSettingsResponse.AsObject> {
  return new Promise((resolve, reject) => {
    let request = new UpdateSettingsRequest();
    let settingsRef = new ResourceRef();
    settingsRef.setName(updateSettingsRequest.ref!.name);
    settingsRef.setNamespace(updateSettingsRequest.ref!.namespace);
    request.setRef(settingsRef);
    let duration = new Duration();
    duration.setSeconds(updateSettingsRequest.refreshRate!.seconds);
    duration.setNanos(updateSettingsRequest.refreshRate!.nanos);
    request.setRefreshRate(duration);
    request.setWatchNamespacesList(updateSettingsRequest.watchNamespacesList);
    client.updateSettings(request, (error, data) => {
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

function getIsLicenseValid(): Promise<GetIsLicenseValidResponse.AsObject> {
  return new Promise((resolve, reject) => {
    let request = new GetIsLicenseValidRequest();
    client.getIsLicenseValid(request, (error, data) => {
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

function listNamespaces(): Promise<ListNamespacesResponse.AsObject> {
  return new Promise((resolve, reject) => {
    let request = new ListNamespacesRequest();
    client.listNamespaces(request, (error, data) => {
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

function getPodNamespace(): Promise<GetPodNamespaceResponse.AsObject> {
  return new Promise((resolve, reject) => {
    let request = new GetPodNamespaceRequest();
    client.getPodNamespace(request, (error, data) => {
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

export const config = {
  getVersion,
  getOAuthEndpoint,
  getSettings,
  updateSettings,
  getIsLicenseValid,
  listNamespaces,
  getPodNamespace
};
