import { grpc } from '@improbable-eng/grpc-web';
import { Duration } from 'google-protobuf/google/protobuf/duration_pb';
import {
  GetIsLicenseValidRequest,
  GetIsLicenseValidResponse,
  GetOAuthEndpointRequest,
  GetOAuthEndpointResponse,
  GetPodNamespaceRequest,
  GetPodNamespaceResponse,
  GetSettingsRequest,
  GetSettingsResponse,
  GetVersionRequest,
  GetVersionResponse,
  ListNamespacesRequest,
  ListNamespacesResponse,
  UpdateSettingsRequest,
  UpdateSettingsResponse
} from 'proto/solo-projects/projects/grpcserver/api/v1/config_pb';
import { ConfigApiClient } from 'proto/solo-projects/projects/grpcserver/api/v1/config_pb_service';
import { host } from 'store';

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
function getSettings(): Promise<GetSettingsResponse> {
  return new Promise((resolve, reject) => {
    let request = new GetSettingsRequest();
    client.getSettings(request, (error, data) => {
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

function updateWatchNamespaces(updateWatchNamespacesRequest: {
  watchNamespacesList: string[];
}): Promise<UpdateSettingsResponse.AsObject> {
  return new Promise(async (resolve, reject) => {
    let currentSettingsReq = await config.getSettings();
    let settingsToUpdate = currentSettingsReq.getSettings();
    let request = new UpdateSettingsRequest();
    let { watchNamespacesList } = updateWatchNamespacesRequest;
    if (settingsToUpdate !== undefined) {
      settingsToUpdate.setWatchNamespacesList(watchNamespacesList);

      request.setSettings(settingsToUpdate);
    }
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

function updateRefreshRate(updateRefreshRateRequest: {
  refreshRate: Duration.AsObject;
}): Promise<UpdateSettingsResponse.AsObject> {
  return new Promise(async (resolve, reject) => {
    let currentSettingsReq = await config.getSettings();
    let settingsToUpdate = currentSettingsReq.getSettings();
    let newRefreshRate = new Duration();
    newRefreshRate.setNanos(updateRefreshRateRequest.refreshRate.nanos);
    newRefreshRate.setSeconds(updateRefreshRateRequest.refreshRate.seconds);

    let updateSettingsRequest = new UpdateSettingsRequest();
    if (settingsToUpdate !== undefined) {
      settingsToUpdate.setRefreshRate(newRefreshRate);
      updateSettingsRequest.setSettings(settingsToUpdate);
    }

    client.updateSettings(updateSettingsRequest, (error, data) => {
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
  updateSettingsRequest: UpdateSettingsRequest
): Promise<UpdateSettingsResponse.AsObject> {
  return new Promise(async (resolve, reject) => {
    client.updateSettings(updateSettingsRequest, (error, data) => {
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
  getPodNamespace,
  updateRefreshRate,
  updateWatchNamespaces
};
