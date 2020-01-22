/* eslint-disable */
import { EnvoyApiClient } from 'proto/solo-projects/projects/grpcserver/api/v1/envoy_pb_service';
import { host } from 'store';
import { grpc } from '@improbable-eng/grpc-web';
import {
  ListEnvoyDetailsResponse,
  ListEnvoyDetailsRequest,
  EnvoyDetails
} from 'proto/solo-projects/projects/grpcserver/api/v1/envoy_pb';

const client = new EnvoyApiClient(host, {
  transport: grpc.CrossBrowserHttpTransport({ withCredentials: false }),
  debug: true
});

function getEnvoyList(): Promise<EnvoyDetails.AsObject[]> {
  return new Promise((resolve, reject) => {
    let req = new ListEnvoyDetailsRequest();

    client.listEnvoyDetails(req, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        resolve(data!.toObject().envoyDetailsList);
      }
    });
  });
}

export const envoyAPI = {
  getEnvoyList
};
